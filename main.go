package main

import (
	"log"
	"my-fiber-app/database"
	"my-fiber-app/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to DB
	database.ConnectDb()

	app := fiber.New()

	// Middleware
	app.Use(logger.New())

	// Setup Routes
	routes.SetupRoutes(app)

	// Start server
	log.Fatal(app.Listen(":3000"))
}
