package controllers

import (
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

	// now, we want to get the user authenticatio

	return c.JSON(fiber.Map{
		"status": "success",
	})
}
