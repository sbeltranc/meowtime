package routes

import (
	"main/controllers"

	"github.com/gocql/gocql"
	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App, scylla *gocql.Session) {
	userController := controllers.NewUsersController(scylla)

	users := app.Group("/users-service")

	users.Get("/user/:id", userController.UserInfo)
	users.Get("/user/authenticated", userController.AuthenticatedUser)
}
