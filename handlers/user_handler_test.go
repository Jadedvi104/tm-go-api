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

// Setup test database
func setupTestDB() *gorm.DB {
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

// ============ USER HANDLER TESTS ============

func TestCreateUser(t *testing.T) {
	app := fiber.New()
	db := setupTestDB()
	database.Database.Db = db

	app.Post("/users/create", CreateUser)

	tests := []struct {
		name           string
		body           models.User
		expectedStatus int
		checkField     func(resp models.User) bool
	}{
		{
			name: "Create user successfully",
			body: models.User{
				FirstName:    "John",
				LastName:     "Doe",
				Email:        "john@example.com",
				Phone:        "1234567890",
				PasswordHash: "hashedpassword",
				Role:         "traveler",
			},
			expectedStatus: 201,
			checkField: func(resp models.User) bool {
				return resp.FirstName == "John" && resp.Email == "john@example.com"
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.body)
			req := httptest.NewRequest("POST", "/users/create", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := app.Test(req, 1)
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
		})
	}
}

func TestGetUsers(t *testing.T) {
	app := fiber.New()
	db := setupTestDB()
	database.Database.Db = db

	// Create test users
	db.Create(&models.User{
		FirstName:    "John",
		LastName:     "Doe",
		Email:        "john@example.com",
		PasswordHash: "password",
		Role:         "traveler",
	})

	app.Get("/users", GetUsers)

	req := httptest.NewRequest("GET", "/users", nil)
	resp, _ := app.Test(req, 1)

	assert.Equal(t, 200, resp.StatusCode)

	var users []models.User
	json.NewDecoder(resp.Body).Decode(&users)
	assert.Greater(t, len(users), 0)
}

func TestGetUserByID(t *testing.T) {
	app := fiber.New()
	db := setupTestDB()
	database.Database.Db = db

	// Create test user
	user := models.User{
		FirstName:    "Jane",
		LastName:     "Smith",
		Email:        "jane@example.com",
		PasswordHash: "password",
		Role:         "local-expert",
	}
	db.Create(&user)

	app.Get("/users/:id", GetUser)

	req := httptest.NewRequest("GET", "/users/1", nil)
	resp, _ := app.Test(req, 1)

	assert.Equal(t, 200, resp.StatusCode)

	var respUser models.User
	json.NewDecoder(resp.Body).Decode(&respUser)
	assert.Equal(t, "Jane", respUser.FirstName)
}

func TestGetUserNotFound(t *testing.T) {
	app := fiber.New()
	db := setupTestDB()
	database.Database.Db = db

	app.Get("/users/:id", GetUser)

	req := httptest.NewRequest("GET", "/users/999", nil)
	resp, _ := app.Test(req, 1)

	assert.Equal(t, 404, resp.StatusCode)
}

func TestUpdateUser(t *testing.T) {
	app := fiber.New()
	db := setupTestDB()
	database.Database.Db = db

	// Create test user
	user := models.User{
		FirstName:    "John",
		LastName:     "Doe",
		Email:        "john@example.com",
		PasswordHash: "password",
		Role:         "traveler",
	}
	db.Create(&user)

	app.Put("/users/:id", UpdateUser)

	updateData := models.User{
		FirstName: "Jonathan",
	}
	body, _ := json.Marshal(updateData)
	req := httptest.NewRequest("PUT", "/users/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, 1)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestDeleteUser(t *testing.T) {
	app := fiber.New()
	db := setupTestDB()
	database.Database.Db = db

	// Create test user
	user := models.User{
		FirstName:    "John",
		LastName:     "Doe",
		Email:        "john@example.com",
		PasswordHash: "password",
		Role:         "traveler",
	}
	db.Create(&user)

	app.Delete("/users/:id", DeleteUser)

	req := httptest.NewRequest("DELETE", "/users/1", nil)
	resp, _ := app.Test(req, 1)

	assert.Equal(t, 200, resp.StatusCode)
}

// ============ INVALID PAYLOAD TESTS ============

func TestCreateUserInvalidPayload(t *testing.T) {
	app := fiber.New()
	db := setupTestDB()
	database.Database.Db = db

	app.Post("/users/create", CreateUser)

	req := httptest.NewRequest("POST", "/users/create", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, 1)
	assert.Equal(t, 400, resp.StatusCode)
}

// ============ EDGE CASE TESTS ============

func TestCreateUserEmptyEmail(t *testing.T) {
	app := fiber.New()
	db := setupTestDB()
	database.Database.Db = db

	app.Post("/users/create", CreateUser)

	user := models.User{
		FirstName:    "John",
		LastName:     "Doe",
		Email:        "", // Empty email
		PasswordHash: "password",
		Role:         "traveler",
	}
	body, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/users/create", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, 1)
	// Should fail due to NOT NULL constraint
	assert.NotEqual(t, 201, resp.StatusCode)
}
