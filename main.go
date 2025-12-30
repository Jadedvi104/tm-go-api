package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Create a new Fiber instance
	app := fiber.New(fiber.Config{
		AppName: "TM Go API v1.0.0",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to TM Go API",
			"status":  "running",
		})
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	// API v1 routes
	v1 := app.Group("/api/v1")

	v1.Get("/users", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"users": []fiber.Map{
				{"id": 1, "name": "John Doe"},
				{"id": 2, "name": "Jane Smith"},
			},
		})
	})

	v1.Get("/users/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		return c.JSON(fiber.Map{
			"user": fiber.Map{
				"id":   id,
				"name": "Sample User",
			},
		})
	})

	v1.Post("/users", func(c *fiber.Ctx) error {
		type User struct {
			Name string `json:"name"`
		}

		user := new(User)
		if err := c.BodyParser(user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "User created successfully",
			"user": fiber.Map{
				"id":   3,
				"name": user.Name,
			},
		})
	})

	// Start server
	port := ":3000"
	log.Printf("Server starting on http://localhost%s", port)
	if err := app.Listen(port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
