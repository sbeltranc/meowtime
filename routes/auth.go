package routes

import (
	"main/controllers"

	"github.com/gocql/gocql"
	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app *fiber.App, scylla *gocql.Session) {
	authController := controllers.NewAuthController(scylla)

	auth := app.Group("/auth-service")
	auth.Get("/authenticated", authController.Authenticated)

	slack := auth.Group("/slack")

	slack.Get("/login", authController.SlackLogin)
	slack.Get("/callback", authController.SlackCallback)
}
