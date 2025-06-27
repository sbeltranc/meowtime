package controllers

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"main/models"

	"github.com/gocql/gocql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/slack"
)

type AuthController struct {
	scylla *gocql.Session
}

func NewAuthController(scylla *gocql.Session) *AuthController {
	return &AuthController{scylla: scylla}
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*()_+"
	b := make([]byte, length)
	rand.Read(b)
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return string(b)
}

var (
	oauthConf = &oauth2.Config{
		ClientID:     os.Getenv("SLACK_CLIENT_ID"),
		ClientSecret: os.Getenv("SLACK_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("SLACK_REDIRECT_URL"),
		Scopes:       []string{"email", "profile", "openid"},
		Endpoint:     slack.Endpoint,
	}
	sessionStore = session.New()
)

func (ac *AuthController) SlackLogin(c *fiber.Ctx) error {
	url := oauthConf.AuthCodeURL("", oauth2.AccessTypeOffline)
	return c.Redirect(url, http.StatusTemporaryRedirect)
}

func (ac *AuthController) SlackCallback(c *fiber.Ctx) error {
	code := c.Query("code")

	// verifying if the code is present
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "The Authorization code on the callback URL is missing"},
		)
	}

	// exchanging the code for a token from slack
	token, err := oauthConf.Exchange(c.Context(), code)
	if err != nil {
		// for some reason the exchange failed, we return an error for the user 2 see what happened
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": fmt.Sprintf("The Token exchange between Slack and meowtime has failed due to '%s', try again later or contact santiago [at] hackclub [dot] app if this issue persists", err)},
		)
	}

	sess, err := sessionStore.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "Failed to create a session for the user, try again later or contact santiago [at] hackclub [dot] app if this issue persists"},
		)
	}

	sess.Set("slack_token", token.AccessToken)
	sess.Save()

	// fetching the user information from Slack with the token we obtained
	client := oauthConf.Client(c.Context(), token)
	resp, err := client.Get("https://slack.com/api/openid.connect.userInfo")
	if err != nil {
		// ts so dumb, if this happens im gonna yell at slack for it
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "We failed to obtain your information from your Slack account, try again later or contact santiago [at] hackclub [dot] app if this issue persists"},
		)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "Slack did not return a valid response, try again later or contact santiago [at] hackclub [dot] app if this issue persists"},
		)
	}

	// parsing all the response data from Slack
	var userData map[string]interface{ any }
	err = json.NewDecoder(resp.Body).Decode(&userData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "Failed to decode the response from Slack, try again later or contact santiago [at] hackclub [dot] app if this issue persists"},
		)
	}

	if userData["ok"] != true {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "Something went wrong on Slack's side while obtaining your Slack account information, try again later or contact santiago [at] hackclub [dot] app if this issue persists"},
		)
	}

	if userData["email_verified"] != true {
		return c.Status(fiber.StatusConflict).JSON(
			fiber.Map{"error": "Your Slack account email is not verified, please verify your email on Slack before logging in to meowtime with it."},
		)
	}

	// now that we are sure the user data is correct, let's check if the user already exists in our database
	// we are gonna search for the sub (which is the user ID) in the database
	userID := userData["sub"].(string)
	user, err := models.GetUserBySUB(ac.scylla, userID)

	if err != nil {
		// let's check if the error isnt because the user was not found
		if err != gocql.ErrNotFound {
			return c.Status(fiber.StatusInternalServerError).JSON(
				fiber.Map{"error": "An internal error occurred while searching for your account, try again later or contact santiago [at] hackclub [dot] app if this issue persists"},
			)
		}

		newUser := &models.User{
			ID:       gocql.TimeUUID().String(),
			SUB:      userID,
			Email:    userData["email"].(string),
			Name:     userData["name"].(string),
			Picture:  userData["picture"].(string),
			TeamID:   userData["https://slack.com/team_id"].(string),
			TeamName: userData["https://slack.com/team_name"].(string),
		}
		err = models.CreateUser(ac.scylla, newUser)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				fiber.Map{"error": "Your account was not created in our database for some internal error, try again later or contact santiago [at] hackclub [dot] app if this issue persists"},
			)
		}

		// creating the session for the user
		newSession := &models.Session{
			ID:           gocql.TimeUUID().String(),
			UserID:       newUser.ID,
			SessionToken: generateRandomString(64),
			ExpiresAt:    time.Now().Add(time.Hour * 24).Unix(),
		}

		session, err := models.CreateSession(newSession, ac.scylla)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				fiber.Map{"error": "Your session was not created in our database for some internal error, try again later or contact santiago [at] hackclub [dot] app if this issue persists"},
			)
		}

		c.Set("Authorization", fmt.Sprintf("Bearer %s", session.SessionToken))

		return c.Status(fiber.StatusCreated).JSON(
			fiber.Map{
				"id":        newUser.ID,
				"name":      newUser.Name,
				"email":     newUser.Email,
				"picture":   newUser.Picture,
				"team_id":   newUser.TeamID,
				"team_name": newUser.TeamName,
			},
		)
	}

	// creating the session and adding the token to the headers
	newSession := &models.Session{
		ID:           gocql.TimeUUID().String(),
		UserID:       user.ID,
		SessionToken: generateRandomString(64),
		ExpiresAt:    time.Now().Add(time.Hour * 24).Unix(),
	}

	session, err := models.CreateSession(newSession, ac.scylla)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "Your session was not created in our database for some internal error, try again later or contact santiago [at] hackclub [dot] app if this issue persists"},
		)
	}

	c.Set("Authorization", fmt.Sprintf("Bearer %s", session.SessionToken))

	// just for being safe, we are gonna update the user information in case it changed
	user.Email = userData["email"].(string)
	user.Name = userData["name"].(string)
	user.Picture = userData["picture"].(string)
	user.TeamID = userData["https://slack.com/team_id"].(string)
	user.TeamName = userData["https://slack.com/team_name"].(string)

	err = models.UpdateUser(ac.scylla, user)

	if err != nil {
		// eh it's whatever, let's say it was not updated and continue ig
		return c.Status(fiber.StatusNotModified).JSON(
			fiber.Map{
				"id":        user.ID,
				"name":      user.Name,
				"email":     user.Email,
				"picture":   user.Picture,
				"team_id":   user.TeamID,
				"team_name": user.TeamName,
			},
		)
	}

	return c.Status(fiber.StatusFound).JSON(
		fiber.Map{
			"id":        user.ID,
			"name":      user.Name,
			"email":     user.Email,
			"picture":   user.Picture,
			"team_id":   user.TeamID,
			"team_name": user.TeamName,
		},
	)
}
