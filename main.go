package main

import (
	"os"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/joho/godotenv"

	"github.com/gocql/gocql"
	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/fiber/v2/middleware/requestid"

	"main/models"
	"main/routes"
)

func main() {
	// loading the environment file
	err := godotenv.Load()

	if err != nil {
		panic(err)
	}

	// intializing the database
	var cluster = gocql.NewCluster(os.Getenv("SCYLLA_HOST_IP") + ":" + os.Getenv("SCYLLA_HOST_PORT"))
	cluster.Keyspace = os.Getenv("SCYLLA_KEYSPACE")

	scylla, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}

	defer scylla.Close()

	// intializing the application
	app := fiber.New(fiber.Config{
		Prefork:     true,
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	// adding the requests id
	app.Use(requestid.New(requestid.Config{
		Header: "Meowtime-Request-Id",
		Generator: func() string {
			return uuid.New().String()
		},
	}))

	// healthcheck route
	app.Get("/up", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("beep boop, meowtime OK!")
	})

	// configuring schemas
	models.InitUserSchema(scylla)
	models.InitSonarSchema(scylla)

	// authentication middleware

	// setup auth routes
	routes.SetupUserRoutes(app, scylla)
	routes.SetupAuthRoutes(app, scylla)

	// listening application
	app.Listen(":" + os.Getenv("SERVER_PORT"))
}
