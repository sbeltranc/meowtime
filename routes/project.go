package routes

import (
	"main/controllers"

	"github.com/gocql/gocql"
	"github.com/gofiber/fiber/v2"
)

func SetupProjectRoutes(app *fiber.App, scylla *gocql.Session) {
	projectController := controllers.NewProjectController(scylla)

	project := app.Group("/project-service")

	project.Get("/project", projectController.ObtainProjectById)
}
