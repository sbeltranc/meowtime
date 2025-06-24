package routes

import (
	"main/controllers"

	"github.com/gocql/gocql"
	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app *fiber.App, scylla *gocql.Session) {
	authController := controllers.NewAuthController(scylla)

	auth := app.Group("/auth-service")

	auth.Get("/slack/login", authController.SlackLogin)
	auth.Get("/slack/callback", authController.SlackCallback)
}
