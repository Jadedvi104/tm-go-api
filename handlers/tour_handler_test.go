package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"tm-go-api/database"
	"tm-go-api/models"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Setup tour test database
func setupTourTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(
		&models.User{},
		&models.TourCategory{},
		&models.TourDestination{},
		&models.Tour{},
		&models.TourImage{},
		&models.TourTag{},
		&models.TourInclude{},
		&models.TourExclude{},
		&models.TourItinerary{},
		&models.TourItineraryActivity{},
		&models.TourItineraryMeal{},
	)
	return db
}

// ============ TOUR CATEGORY TESTS ============

func TestCreateTourCategory(t *testing.T) {
	app := fiber.New()
	db := setupTourTestDB()
	database.Database.Db = db

	app.Post("/tour-categories", CreateTourCategory)

	tests := []struct {
		name           string
		body           models.TourCategory
		expectedStatus int
	}{
		{
			name: "Create category successfully",
			body: models.TourCategory{
				Name:        "Adventure",
				Description: "Adventure tours",
			},
			expectedStatus: 201,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.body)
			req := httptest.NewRequest("POST", "/tour-categories", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := app.Test(req, 1)
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
		})
	}
}

func TestGetTourCategories(t *testing.T) {
	app := fiber.New()
	db := setupTourTestDB()
	database.Database.Db = db

	// Create test category
	db.Create(&models.TourCategory{
		Name:        "Food & Culture",
		Description: "Food and culture tours",
	})

	app.Get("/tour-categories", GetTourCategories)

	req := httptest.NewRequest("GET", "/tour-categories", nil)
	resp, _ := app.Test(req, 1)

	assert.Equal(t, 200, resp.StatusCode)

	var categories []models.TourCategory
	json.NewDecoder(resp.Body).Decode(&categories)
	assert.Greater(t, len(categories), 0)
}

func TestGetTourCategoryByID(t *testing.T) {
	app := fiber.New()
	db := setupTourTestDB()
	database.Database.Db = db

	// Create test category
	category := models.TourCategory{
		Name:        "Nature Trails",
		Description: "Nature trail tours",
	}
	db.Create(&category)

	app.Get("/tour-categories/:id", GetTourCategory)

	req := httptest.NewRequest("GET", "/tour-categories/1", nil)
	resp, _ := app.Test(req, 1)

	assert.Equal(t, 200, resp.StatusCode)

	var respCategory models.TourCategory
	json.NewDecoder(resp.Body).Decode(&respCategory)
	assert.Equal(t, "Nature Trails", respCategory.Name)
}

func TestGetTourCategoryNotFound(t *testing.T) {
	app := fiber.New()
	db := setupTourTestDB()
	database.Database.Db = db

	app.Get("/tour-categories/:id", GetTourCategory)

	req := httptest.NewRequest("GET", "/tour-categories/999", nil)
	resp, _ := app.Test(req, 1)

	assert.Equal(t, 404, resp.StatusCode)
}

// ============ TOUR DESTINATION TESTS ============

func TestCreateTourDestination(t *testing.T) {
	app := fiber.New()
	db := setupTourTestDB()
	database.Database.Db = db

	app.Post("/tour-destinations", CreateTourDestination)

	destination := models.TourDestination{
		City:      "Paris",
		Country:   "France",
		Latitude:  48.8566,
		Longitude: 2.3522,
	}
	body, _ := json.Marshal(destination)
	req := httptest.NewRequest("POST", "/tour-destinations", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, 1)
	assert.Equal(t, 201, resp.StatusCode)
}

func TestGetTourDestinations(t *testing.T) {
	app := fiber.New()
	db := setupTourTestDB()
	database.Database.Db = db

	// Create test destination
	db.Create(&models.TourDestination{
		City:      "Tokyo",
		Country:   "Japan",
		Latitude:  35.6762,
		Longitude: 139.6503,
	})

	app.Get("/tour-destinations", GetTourDestinations)

	req := httptest.NewRequest("GET", "/tour-destinations", nil)
	resp, _ := app.Test(req, 1)

	assert.Equal(t, 200, resp.StatusCode)

	var destinations []models.TourDestination
	json.NewDecoder(resp.Body).Decode(&destinations)
	assert.Greater(t, len(destinations), 0)
}

func TestGetTourDestinationWithFiltering(t *testing.T) {
	app := fiber.New()
	db := setupTourTestDB()
	database.Database.Db = db

	// Create test destinations
	db.Create(&models.TourDestination{
		City:      "Rome",
		Country:   "Italy",
		Latitude:  41.9028,
		Longitude: 12.4964,
	})

	app.Get("/tour-destinations", GetTourDestinations)

	req := httptest.NewRequest("GET", "/tour-destinations?country=Italy", nil)
	resp, _ := app.Test(req, 1)

	assert.Equal(t, 200, resp.StatusCode)
}

func TestGetTourDestinationByID(t *testing.T) {
	app := fiber.New()
	db := setupTourTestDB()
	database.Database.Db = db

	// Create test destination
	destination := models.TourDestination{
		City:      "Barcelona",
		Country:   "Spain",
		Latitude:  41.3851,
		Longitude: 2.1734,
	}
	db.Create(&destination)

	app.Get("/tour-destinations/:id", GetTourDestination)

	req := httptest.NewRequest("GET", "/tour-destinations/1", nil)
	resp, _ := app.Test(req, 1)

	assert.Equal(t, 200, resp.StatusCode)

	var respDestination models.TourDestination
	json.NewDecoder(resp.Body).Decode(&respDestination)
	assert.Equal(t, "Barcelona", respDestination.City)
}

// ============ TOUR TESTS ============

func TestCreateTour(t *testing.T) {
	app := fiber.New()
	db := setupTourTestDB()
	database.Database.Db = db

	// Create guide user first
	guide := models.User{
		FirstName:    "John",
		LastName:     "Guide",
		Email:        "guide@example.com",
		PasswordHash: "password",
		Role:         "local-expert",
	}
	db.Create(&guide)

	// Create category
	category := models.TourCategory{
		Name:        "Adventure",
		Description: "Adventure tours",
	}
	db.Create(&category)

	// Create destination
	destination := models.TourDestination{
		City:      "Bangkok",
		Country:   "Thailand",
		Latitude:  13.7563,
		Longitude: 100.5018,
	}
	db.Create(&destination)

	app.Post("/tours", CreateTour)

	price := 99.99
	tour := models.Tour{
		Title:             "Bangkok City Tour",
		Description:       "Explore Bangkok",
		ShortDescription:  "City exploration",
		Slug:              "bangkok-city-tour",
		PrimaryCategoryID: category.ID,
		DestinationID:     destination.ID,
		PriceAmount:       price,
		PriceCurrency:     "USD",
		PricePerPerson:    true,
		DurationValue:     3,
		DurationUnit:      "days",
		MaxParticipants:   20,
		MinParticipants:   2,
		GuideID:           guide.ID,
		GuideName:         "John Guide",
		IsActive:          true,
		IsListed:          true,
	}
	body, _ := json.Marshal(tour)
	req := httptest.NewRequest("POST", "/tours", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, 1)
	assert.Equal(t, 201, resp.StatusCode)
}

func TestGetTours(t *testing.T) {
	app := fiber.New()
	db := setupTourTestDB()
	database.Database.Db = db

	// Setup: Create guide, category, and destination
	guide := models.User{
		FirstName:    "Jane",
		LastName:     "Expert",
		Email:        "expert@example.com",
		PasswordHash: "password",
		Role:         "local-expert",
	}
	db.Create(&guide)

	category := models.TourCategory{
		Name:        "Food",
		Description: "Food tours",
	}
	db.Create(&category)

	destination := models.TourDestination{
		City:      "Bangkok",
		Country:   "Thailand",
		Latitude:  13.7563,
		Longitude: 100.5018,
	}
	db.Create(&destination)

	// Create tour
	price := 150.00
	db.Create(&models.Tour{
		Title:             "Food Tour",
		Description:       "Thai food experience",
		Slug:              "food-tour",
		PrimaryCategoryID: category.ID,
		DestinationID:     destination.ID,
		PriceAmount:       price,
		DurationValue:     2,
		DurationUnit:      "days",
		GuideID:           guide.ID,
		IsActive:          true,
		IsListed:          true,
	})

	app.Get("/tours", GetTours)

	req := httptest.NewRequest("GET", "/tours", nil)
	resp, _ := app.Test(req, 1)

	assert.Equal(t, 200, resp.StatusCode)

	var tours []models.Tour
	json.NewDecoder(resp.Body).Decode(&tours)
	assert.Greater(t, len(tours), 0)
}

func TestGetToursWithPagination(t *testing.T) {
	app := fiber.New()
	db := setupTourTestDB()
	database.Database.Db = db

	app.Get("/tours", GetTours)

	req := httptest.NewRequest("GET", "/tours?page=1&limit=5", nil)
	resp, _ := app.Test(req, 1)

	assert.Equal(t, 200, resp.StatusCode)
}

func TestGetTourByID(t *testing.T) {
	app := fiber.New()
	db := setupTourTestDB()
	database.Database.Db = db

	// Setup
	guide := models.User{
		FirstName:    "Bob",
		LastName:     "Guide",
		Email:        "bob@example.com",
		PasswordHash: "password",
		Role:         "local-expert",
	}
	db.Create(&guide)

	category := models.TourCategory{
		Name:        "Water Sports",
		Description: "Water sports tours",
	}
	db.Create(&category)

	destination := models.TourDestination{
		City:      "Phuket",
		Country:   "Thailand",
		Latitude:  8.1031,
		Longitude: 98.2867,
	}
	db.Create(&destination)

	price := 120.00
	tour := models.Tour{
		Title:             "Scuba Diving",
		Description:       "Scuba diving experience",
		Slug:              "scuba-diving",
		PrimaryCategoryID: category.ID,
		DestinationID:     destination.ID,
		PriceAmount:       price,
		DurationValue:     1,
		DurationUnit:      "days",
		GuideID:           guide.ID,
		IsActive:          true,
		IsListed:          true,
	}
	db.Create(&tour)

	app.Get("/tours/:id", GetTour)

	req := httptest.NewRequest("GET", "/tours/1", nil)
	resp, _ := app.Test(req, 1)

	assert.Equal(t, 200, resp.StatusCode)

	var respTour models.Tour
	json.NewDecoder(resp.Body).Decode(&respTour)
	assert.Equal(t, "Scuba Diving", respTour.Title)
}

func TestGetTourNotFound(t *testing.T) {
	app := fiber.New()
	db := setupTourTestDB()
	database.Database.Db = db

	app.Get("/tours/:id", GetTour)

	req := httptest.NewRequest("GET", "/tours/999", nil)
	resp, _ := app.Test(req, 1)

	assert.Equal(t, 404, resp.StatusCode)
}

func TestUpdateTour(t *testing.T) {
	app := fiber.New()
	db := setupTourTestDB()
	database.Database.Db = db

	// Setup
	guide := models.User{
		FirstName:    "Alice",
		LastName:     "Guide",
		Email:        "alice@example.com",
		PasswordHash: "password",
		Role:         "local-expert",
	}
	db.Create(&guide)

	category := models.TourCategory{
		Name:        "Cultural",
		Description: "Cultural tours",
	}
	db.Create(&category)

	destination := models.TourDestination{
		City:      "Chiang Mai",
		Country:   "Thailand",
		Latitude:  18.7883,
		Longitude: 98.9853,
	}
	db.Create(&destination)

	price := 75.00
	tour := models.Tour{
		Title:             "Temple Tour",
		Description:       "Visit temples",
		Slug:              "temple-tour",
		PrimaryCategoryID: category.ID,
		DestinationID:     destination.ID,
		PriceAmount:       price,
		DurationValue:     1,
		DurationUnit:      "days",
		GuideID:           guide.ID,
		IsActive:          true,
		IsListed:          true,
	}
	db.Create(&tour)

	app.Put("/tours/:id", UpdateTour)

	updateData := models.Tour{
		Title: "Ancient Temple Tour",
	}
	body, _ := json.Marshal(updateData)
	req := httptest.NewRequest("PUT", "/tours/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, 1)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestDeleteTour(t *testing.T) {
	app := fiber.New()
	db := setupTourTestDB()
	database.Database.Db = db

	// Setup
	guide := models.User{
		FirstName:    "Charlie",
		LastName:     "Guide",
		Email:        "charlie@example.com",
		PasswordHash: "password",
		Role:         "local-expert",
	}
	db.Create(&guide)

	category := models.TourCategory{
		Name:        "Adventure",
		Description: "Adventure tours",
	}
	db.Create(&category)

	destination := models.TourDestination{
		City:      "Krabi",
		Country:   "Thailand",
		Latitude:  8.3569,
		Longitude: 98.9322,
	}
	db.Create(&destination)

	price := 200.00
	tour := models.Tour{
		Title:             "Rock Climbing",
		Description:       "Rock climbing adventure",
		Slug:              "rock-climbing",
		PrimaryCategoryID: category.ID,
		DestinationID:     destination.ID,
		PriceAmount:       price,
		DurationValue:     2,
		DurationUnit:      "days",
		GuideID:           guide.ID,
		IsActive:          true,
		IsListed:          true,
	}
	db.Create(&tour)

	app.Delete("/tours/:id", DeleteTour)

	req := httptest.NewRequest("DELETE", "/tours/1", nil)
	resp, _ := app.Test(req, 1)

	assert.Equal(t, 200, resp.StatusCode)
}

// ============ EDGE CASE TESTS ============

func TestCreateTourInvalidPayload(t *testing.T) {
	app := fiber.New()
	db := setupTourTestDB()
	database.Database.Db = db

	app.Post("/tours", CreateTour)

	req := httptest.NewRequest("POST", "/tours", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, 1)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestGetToursWithPriceFiltering(t *testing.T) {
	app := fiber.New()
	db := setupTourTestDB()
	database.Database.Db = db

	app.Get("/tours", GetTours)

	req := httptest.NewRequest("GET", "/tours?min_price=50&max_price=150", nil)
	resp, _ := app.Test(req, 1)

	assert.Equal(t, 200, resp.StatusCode)
}
