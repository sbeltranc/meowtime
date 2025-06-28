package controllers

import (
	"encoding/json"
	"main/models"

	"github.com/gocql/gocql"
	"github.com/gofiber/fiber/v2"
)

type SonarController struct {
	scylla *gocql.Session
}

func NewSonarController(scylla *gocql.Session) *SonarController {
	return &SonarController{scylla: scylla}
}

func (ac *SonarController) SonarCall(c *fiber.Ctx) error {
	var req struct {
		ProjectName      string         `json:"project_name"`
		TotalTimeSeconds int            `json:"total_time_seconds"`
		Metadata         map[string]any `json:"metadata"`
	}

	// parsing dumb dumb body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "There was invalid data given to the Sonar call endpoint, please check the request data on what's wrong with it"},
		)
	}

	user := c.Locals("user").(*models.User)
	authenticated := c.Locals("authenticated").(bool)

	if !authenticated {
		return c.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{"error": "You are required to authenticate your session for being able to call the Sonar Service"},
		)
	}

	// we are gonna search for the project by it's name & the person's userid
	project, err := models.FindProjectByNameAndOwner(req.ProjectName, user.ID, ac.scylla)
	if err != nil {
		if err != gocql.ErrNotFound {
			return c.Status(fiber.StatusInternalServerError).JSON(
				fiber.Map{"error": "Something went wrong while trying to associate your Sonar call with a project"},
			)
		}

		project, err = models.CreateProject(req.ProjectName, user.ID, ac.scylla)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				fiber.Map{"error": "Something went wrong while trying to create a new project"},
			)
		}
	}

	sonar := &models.Sonar{
		ID:     gocql.TimeUUID().String(),
		UserID: user.ID,
		Metadata: func() json.RawMessage {
			data, _ := json.Marshal(req.Metadata)
			return json.RawMessage(data)
		}(),
		IPAddress: c.Locals("ip").(string),
		ProjectID: project.ID,
		TotalTime: int64(req.TotalTimeSeconds),
		Software:  c.Get("User-Agent"),
	}

	sonar, err = models.CreateSonar(sonar, ac.scylla)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "Something went wrong while trying to create a new sonar"},
		)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
	})
}
