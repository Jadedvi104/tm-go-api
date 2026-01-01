package routes

import (
	"tm-go-api/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Grouping routes is cleaner
	api := app.Group("/api/v1")

	user := api.Group("/users")
	user.Post("/", handlers.CreateUser)
	user.Get("/", handlers.GetUsers)
	user.Get("/:id", handlers.GetUser)
	user.Put("/:id", handlers.UpdateUser)
	user.Delete("/:id", handlers.DeleteUser)
}
