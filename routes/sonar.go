package routes

import (
	"main/controllers"

	"github.com/gocql/gocql"
	"github.com/gofiber/fiber/v2"
)

func SetupSonarRoutes(app *fiber.App, scylla *gocql.Session) {
	sonarController := controllers.NewSonarController(scylla)

	sonar := app.Group("/sonar-service")

	sonar.Post("/call", sonarController.SonarCall)
}
