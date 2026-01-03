package main

import (
	"log"
	"tm-go-api/database"
	"tm-go-api/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on System Environment Variables")
	}

	// Connect to DB
	database.ConnectDb()

	app := fiber.New()

	// Middleware
	app.Use(logger.New())

	// Health Check
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("I am alive!") // Azure loves 200 OK responses
	})

	// Setup Routes
	routes.SetupRoutes(app)

	// Start server
	log.Fatal(app.Listen(":80"))
}
