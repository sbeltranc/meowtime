package controllers

import (
	"main/models"

	"github.com/gocql/gocql"
	"github.com/gofiber/fiber/v2"
)

type UsersController struct {
	scylla *gocql.Session
}

func NewUsersController(scylla *gocql.Session) *UsersController {
	return &UsersController{scylla: scylla}
}

func (ac *UsersController) UserInfo(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing user id"})
	}
	user, err := models.GetUserByID(ac.scylla, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	return c.JSON(fiber.Map{
		"id":        user.ID,
		"email":     user.Email,
		"name":      user.Name,
		"picture":   user.Picture,
		"team_id":   user.TeamID,
		"team_name": user.TeamName,
	})
}
