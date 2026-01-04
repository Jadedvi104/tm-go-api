package handlers

import (
	"tm-go-api/database"
	"tm-go-api/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ============ TOUR HANDLERS ============

// CreateTour creates a new tour
func CreateTour(c *fiber.Ctx) error {
	tour := new(models.Tour)
	if err := c.BodyParser(tour); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid payload"})
	}

	if err := database.Database.Db.Create(&tour).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create tour"})
	}

	return c.Status(201).JSON(tour)
}

// GetTours retrieves all tours with filtering and pagination
func GetTours(c *fiber.Ctx) error {
	var tours []models.Tour
	query := database.Database.Db

	// Filtering
	if category := c.Query("category"); category != "" {
		query = query.Where("primary_category_id = ?", category)
	}
	if destination := c.Query("destination"); destination != "" {
		query = query.Where("destination_id = ?", destination)
	}
	if minPrice := c.Query("min_price"); minPrice != "" {
		query = query.Where("price_amount >= ?", minPrice)
	}
	if maxPrice := c.Query("max_price"); maxPrice != "" {
		query = query.Where("price_amount <= ?", maxPrice)
	}

	// Sorting
	if sortBy := c.Query("sort", "created_at"); sortBy != "" {
		order := c.Query("order", "DESC")
		query = query.Order(sortBy + " " + order)
	}

	// Pagination
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	query.Offset(offset).Limit(limit).
		Preload("Guide").
		Preload("Category").
		Preload("Destination").
		Preload("Images").
		Find(&tours)

	return c.Status(200).JSON(tours)
}

// GetTour retrieves a tour by ID
func GetTour(c *fiber.Ctx) error {
	id := c.Params("id")
	var tour models.Tour

	result := database.Database.Db.
		Preload("Guide").
		Preload("Category").
		Preload("Destination").
		Preload("Images").
		Preload("Tags").
		Preload("Includes").
		Preload("Excludes").
		Preload("Itineraries", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Activities").Preload("Meals")
		}).
		Find(&tour, id)

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Tour not found"})
	}

	return c.Status(200).JSON(tour)
}

// UpdateTour updates an existing tour
func UpdateTour(c *fiber.Ctx) error {
	id := c.Params("id")
	tour := new(models.Tour)

	if err := c.BodyParser(tour); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid payload"})
	}

	result := database.Database.Db.Model(&models.Tour{}).Where("id = ?", id).Updates(tour)

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Tour not found"})
	}

	if err := result.Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not update tour"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Tour updated successfully"})
}

// DeleteTour soft deletes a tour
func DeleteTour(c *fiber.Ctx) error {
	id := c.Params("id")

	result := database.Database.Db.Delete(&models.Tour{}, id)

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Tour not found"})
	}

	if err := result.Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not delete tour"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Tour deleted successfully"})
}

// ============ TOUR CATEGORY HANDLERS ============

// CreateTourCategory creates a new tour category
func CreateTourCategory(c *fiber.Ctx) error {
	category := new(models.TourCategory)
	if err := c.BodyParser(category); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid payload"})
	}

	if err := database.Database.Db.Create(&category).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create category"})
	}

	return c.Status(201).JSON(category)
}

// GetTourCategories retrieves all tour categories
func GetTourCategories(c *fiber.Ctx) error {
	var categories []models.TourCategory
	database.Database.Db.Find(&categories)
	return c.Status(200).JSON(categories)
}

// GetTourCategory retrieves a tour category by ID
func GetTourCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	var category models.TourCategory

	result := database.Database.Db.Find(&category, id)

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Category not found"})
	}

	return c.Status(200).JSON(category)
}

// UpdateTourCategory updates an existing tour category
func UpdateTourCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	category := new(models.TourCategory)

	if err := c.BodyParser(category); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid payload"})
	}

	result := database.Database.Db.Model(&models.TourCategory{}).Where("id = ?", id).Updates(category)

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Category not found"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Category updated successfully"})
}

// DeleteTourCategory deletes a tour category
func DeleteTourCategory(c *fiber.Ctx) error {
	id := c.Params("id")

	result := database.Database.Db.Delete(&models.TourCategory{}, id)

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Category not found"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Category deleted successfully"})
}

// ============ TOUR DESTINATION HANDLERS ============

// CreateTourDestination creates a new tour destination
func CreateTourDestination(c *fiber.Ctx) error {
	destination := new(models.TourDestination)
	if err := c.BodyParser(destination); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid payload"})
	}

	if err := database.Database.Db.Create(&destination).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create destination"})
	}

	return c.Status(201).JSON(destination)
}

// GetTourDestinations retrieves all tour destinations
func GetTourDestinations(c *fiber.Ctx) error {
	var destinations []models.TourDestination
	query := database.Database.Db

	// Filtering by country
	if country := c.Query("country"); country != "" {
		query = query.Where("country = ?", country)
	}

	// Filtering by city
	if city := c.Query("city"); city != "" {
		query = query.Where("city = ?", city)
	}

	query.Find(&destinations)
	return c.Status(200).JSON(destinations)
}

// GetTourDestination retrieves a tour destination by ID
func GetTourDestination(c *fiber.Ctx) error {
	id := c.Params("id")
	var destination models.TourDestination

	result := database.Database.Db.Preload("Tours").Find(&destination, id)

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Destination not found"})
	}

	return c.Status(200).JSON(destination)
}

// UpdateTourDestination updates an existing tour destination
func UpdateTourDestination(c *fiber.Ctx) error {
	id := c.Params("id")
	destination := new(models.TourDestination)

	if err := c.BodyParser(destination); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid payload"})
	}

	result := database.Database.Db.Model(&models.TourDestination{}).Where("id = ?", id).Updates(destination)

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Destination not found"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Destination updated successfully"})
}

// DeleteTourDestination deletes a tour destination
func DeleteTourDestination(c *fiber.Ctx) error {
	id := c.Params("id")

	result := database.Database.Db.Delete(&models.TourDestination{}, id)

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Destination not found"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Destination deleted successfully"})
}
