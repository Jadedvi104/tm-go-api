package handlers

import (
	"tm-go-api/database"
	"tm-go-api/models"

	"github.com/gofiber/fiber/v2"
)

// CreateUser
func CreateUser(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid payload"})
	}

	// TODO: In a real app, hash the password here before saving!
	// user.PasswordHash = hashPassword(user.PasswordHash)

	if err := database.Database.Db.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create user"})
	}

	// SECURITY: Clear the password before returning the data to the user
	user.PasswordHash = ""

	return c.Status(201).JSON(user)
}

// GetUsers
func GetUsers(c *fiber.Ctx) error {
	var users []models.User
	database.Database.Db.Find(&users)
	return c.Status(200).JSON(users)
}

// GetUser (by ID)
func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	result := database.Database.Db.Find(&user, id)
	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.Status(200).JSON(user)
}

// UpdateUser
func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	// Check if exists
	if result := database.Database.Db.First(&user, id); result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	// Parse body into a struct
	type UpdateUser struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Bio       string `json:"bio"`
	}

	var updateData UpdateUser
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid payload"})
	}

	// Update specific fields
	user.FirstName = updateData.FirstName
	user.LastName = updateData.LastName
	user.Bio = updateData.Bio

	database.Database.Db.Save(&user)

	return c.Status(200).JSON(user)
}

// DeleteUser
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	if result := database.Database.Db.First(&user, id); result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	database.Database.Db.Delete(&user)
	return c.Status(200).JSON(fiber.Map{"message": "User deleted successfully"})
}
