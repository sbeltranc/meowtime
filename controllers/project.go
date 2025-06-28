package controllers

import (
	"main/models"

	"github.com/gocql/gocql"
	"github.com/gofiber/fiber/v2"
)

type ProjectController struct {
	scylla *gocql.Session
}

func NewProjectController(scylla *gocql.Session) *ProjectController {
	return &ProjectController{scylla: scylla}
}

func (pc *ProjectController) ObtainProjectById(c *fiber.Ctx) error {
	id := c.Query("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "No Project ID was provided on the request"},
		)
	}

	project, err := models.FindProjectByID(id, pc.scylla)
	if err != nil {
		if err != gocql.ErrNotFound {
			return c.Status(fiber.StatusInternalServerError).JSON(
				fiber.Map{"error": "Something went wrong while trying to retrieve the project"},
			)
		}

		return c.Status(fiber.StatusNotFound).JSON(
			fiber.Map{"error": "There was no project found with the provided ID"},
		)
	}

	return c.Status(fiber.StatusOK).JSON(project)
}
