package main

import (
	"fmt"
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
	app.Use(func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		connectingIP := c.Get("CF-Connecting-IP")

		c.Locals("user", nil)
		c.Locals("authenticated", false)

		if authHeader == "" {
			return c.Next()
		}

		if connectingIP == "" {
			c.Locals("ip", c.IP())
		}

		session, err := models.FindSession(authHeader, scylla)
		if err != nil {
			fmt.Printf("Something went wrong while trying to get session '%s': %v", authHeader, err)
			return c.Next()
		}

		c.Locals("user", session) // no worries, the find session function returns the user
		c.Locals("authenticated", true)

		return c.Next()
	})

	// setup auth routes
	routes.SetupUserRoutes(app, scylla)
	routes.SetupAuthRoutes(app, scylla)
	routes.SetupProjectRoutes(app, scylla)

	// listening application
	app.Listen(":" + os.Getenv("SERVER_PORT"))
}
