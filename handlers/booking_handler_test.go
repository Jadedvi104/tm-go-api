package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"tm-go-api/database"
	"tm-go-api/models"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Setup booking test database
func setupBookingTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(
		&models.User{},
		&models.TourCategory{},
		&models.TourDestination{},
		&models.Tour{},
		&models.Booking{},
		&models.BookingParticipant{},
		&models.BookingPricing{},
		&models.BookingHotelDetails{},
		&models.BookingPreference{},
		&models.Payment{},
		&models.Review{},
		&models.ReviewDetailedRatings{},
		&models.ReviewImage{},
	)
	return db
}

// ============ BOOKING TESTS ============

func TestCreateBooking(t *testing.T) {
	app := fiber.New()
	db := setupBookingTestDB()
	database.Database.Db = db

	// Setup: Create users
	traveler := models.User{
		FirstName:    "John",
		LastName:     "Traveler",
		Email:        "traveler@example.com",
		PasswordHash: "password",
		Role:         "traveler",
	}
	db.Create(&traveler)

	guide := models.User{
		FirstName:    "Jane",
		LastName:     "Guide",
		Email:        "guide@example.com",
		PasswordHash: "password",
		Role:         "local-expert",
	}
	db.Create(&guide)

	// Create category and destination
	category := models.TourCategory{Name: "Adventure", Description: "Adventure tours"}
	db.Create(&category)

	destination := models.TourDestination{
		City:    "Bangkok",
		Country: "Thailand",
	}
	db.Create(&destination)

	// Create tour
	price := 100.00
	tour := models.Tour{
		Title:             "Bangkok Tour",
		Description:       "Explore Bangkok",
		Slug:              "bangkok-tour",
		PrimaryCategoryID: category.ID,
		DestinationID:     destination.ID,
		PriceAmount:       price,
		DurationValue:     3,
		DurationUnit:      "days",
		GuideID:           guide.ID,
		IsActive:          true,
		IsListed:          true,
	}
	db.Create(&tour)

	app.Post("/bookings", CreateBooking)

	booking := models.Booking{
		BookingReference:  "BK001",
		Status:            "pending",
		TravelerID:        traveler.ID,
		TravelerName:      "John Traveler",
		TravelerEmail:     "traveler@example.com",
		TourID:            tour.ID,
		TourTitle:         "Bangkok Tour",
		GuideID:           guide.ID,
		TotalParticipants: 2,
		AdultCount:        2,
		ChildCount:        0,
		StartDate:         time.Now(),
		TotalPrice:        200.00,
		Currency:          "USD",
	}
	body, _ := json.Marshal(booking)
	req := httptest.NewRequest("POST", "/bookings", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, 1)
	assert.Equal(t, 201, resp.StatusCode)
}

func TestGetBookings(t *testing.T) {
	app := fiber.New()
	db := setupBookingTestDB()
	database.Database.Db = db

	// Setup
	traveler := models.User{
		FirstName:    "Alice",
		LastName:     "Traveler",
		Email:        "alice@example.com",
		PasswordHash: "password",
		Role:         "traveler",
	}
	db.Create(&traveler)

	guide := models.User{
		FirstName:    "Bob",
		LastName:     "Guide",
		Email:        "bob@example.com",
		PasswordHash: "password",
		Role:         "local-expert",
	}
	db.Create(&guide)

	category := models.TourCategory{Name: "Food", Description: "Food tours"}
	db.Create(&category)

	destination := models.TourDestination{City: "Paris", Country: "France"}
	db.Create(&destination)

	tour := models.Tour{
		Title:             "Paris Food Tour",
		Description:       "Food experience",
		Slug:              "paris-food",
		PrimaryCategoryID: category.ID,
		DestinationID:     destination.ID,
		PriceAmount:       150.00,
		DurationValue:     2,
		DurationUnit:      "days",
		GuideID:           guide.ID,
		IsActive:          true,
		IsListed:          true,
	}
	db.Create(&tour)

	booking := models.Booking{
		BookingReference:  "BK002",
		Status:            "confirmed",
		TravelerID:        traveler.ID,
		TravelerName:      "Alice Traveler",
		TourID:            tour.ID,
		GuideID:           guide.ID,
		TotalParticipants: 2,
		AdultCount:        2,
		StartDate:         time.Now(),
		TotalPrice:        300.00,
	}
	db.Create(&booking)

	app.Get("/bookings", GetBookings)

	req := httptest.NewRequest("GET", "/bookings", nil)
	resp, _ := app.Test(req, 1)

	assert.Equal(t, 200, resp.StatusCode)

	var bookings []models.Booking
	json.NewDecoder(resp.Body).Decode(&bookings)
	assert.Greater(t, len(bookings), 0)
}

func TestGetBookingByID(t *testing.T) {
	app := fiber.New()
	db := setupBookingTestDB()
	database.Database.Db = db

	// Setup
	traveler := models.User{
		FirstName:    "Charlie",
		LastName:     "Traveler",
		Email:        "charlie@example.com",
		PasswordHash: "password",
		Role:         "traveler",
	}
	db.Create(&traveler)

	guide := models.User{
		FirstName:    "Diana",
		LastName:     "Guide",
		Email:        "diana@example.com",
		PasswordHash: "password",
		Role:         "local-expert",
	}
	db.Create(&guide)

	category := models.TourCategory{Name: "Culture", Description: "Cultural tours"}
	db.Create(&category)

	destination := models.TourDestination{City: "Rome", Country: "Italy"}
	db.Create(&destination)

	tour := models.Tour{
		Title:             "Rome History Tour",
		Description:       "Explore Rome",
		Slug:              "rome-history",
		PrimaryCategoryID: category.ID,
		DestinationID:     destination.ID,
		PriceAmount:       120.00,
		DurationValue:     3,
		DurationUnit:      "days",
		GuideID:           guide.ID,
		IsActive:          true,
		IsListed:          true,
	}
	db.Create(&tour)

	booking := models.Booking{
		BookingReference:  "BK003",
		Status:            "pending",
		TravelerID:        traveler.ID,
		TravelerName:      "Charlie Traveler",
		TourID:            tour.ID,
		GuideID:           guide.ID,
		TotalParticipants: 1,
		AdultCount:        1,
		StartDate:         time.Now(),
		TotalPrice:        120.00,
	}
	db.Create(&booking)

	app.Get("/bookings/:id", GetBooking)

	req := httptest.NewRequest("GET", "/bookings/1", nil)
	resp, _ := app.Test(req, 1)

	assert.Equal(t, 200, resp.StatusCode)

	var respBooking models.Booking
	json.NewDecoder(resp.Body).Decode(&respBooking)
	assert.Equal(t, "BK003", respBooking.BookingReference)
}

func TestGetBookingByReference(t *testing.T) {
	app := fiber.New()
	db := setupBookingTestDB()
	database.Database.Db = db

	// Setup
	traveler := models.User{
		FirstName:    "Eve",
		LastName:     "Traveler",
		Email:        "eve@example.com",
		PasswordHash: "password",
		Role:         "traveler",
	}
	db.Create(&traveler)

	guide := models.User{
		FirstName:    "Frank",
		LastName:     "Guide",
		Email:        "frank@example.com",
		PasswordHash: "password",
		Role:         "local-expert",
	}
	db.Create(&guide)

	category := models.TourCategory{Name: "Nature", Description: "Nature tours"}
	db.Create(&category)

	destination := models.TourDestination{City: "Tokyo", Country: "Japan"}
	db.Create(&destination)

	tour := models.Tour{
		Title:             "Tokyo Nature",
		Description:       "Nature in Tokyo",
		Slug:              "tokyo-nature",
		PrimaryCategoryID: category.ID,
		DestinationID:     destination.ID,
		PriceAmount:       200.00,
		DurationValue:     2,
		DurationUnit:      "days",
		GuideID:           guide.ID,
		IsActive:          true,
		IsListed:          true,
	}
	db.Create(&tour)

	booking := models.Booking{
		BookingReference:  "BK-UNIQUE-2024",
		Status:            "confirmed",
		TravelerID:        traveler.ID,
		TravelerName:      "Eve Traveler",
		TourID:            tour.ID,
		GuideID:           guide.ID,
		TotalParticipants: 2,
		AdultCount:        2,
		StartDate:         time.Now(),
		TotalPrice:        400.00,
	}
	db.Create(&booking)

	app.Get("/bookings/reference/:reference", GetBookingByReference)

	req := httptest.NewRequest("GET", "/bookings/reference/BK-UNIQUE-2024", nil)
	resp, _ := app.Test(req, 1)

	assert.Equal(t, 200, resp.StatusCode)

	var respBooking models.Booking
	json.NewDecoder(resp.Body).Decode(&respBooking)
	assert.Equal(t, "BK-UNIQUE-2024", respBooking.BookingReference)
}

func TestCancelBooking(t *testing.T) {
	app := fiber.New()
	db := setupBookingTestDB()
	database.Database.Db = db

	// Setup
	traveler := models.User{
		FirstName:    "George",
		LastName:     "Traveler",
		Email:        "george@example.com",
		PasswordHash: "password",
		Role:         "traveler",
	}
	db.Create(&traveler)

	guide := models.User{
		FirstName:    "Grace",
		LastName:     "Guide",
		Email:        "grace@example.com",
		PasswordHash: "password",
		Role:         "local-expert",
	}
	db.Create(&guide)

	category := models.TourCategory{Name: "Water Sports", Description: "Water sports tours"}
	db.Create(&category)

	destination := models.TourDestination{City: "Bali", Country: "Indonesia"}
	db.Create(&destination)

	tour := models.Tour{
		Title:             "Bali Water Sports",
		Description:       "Water sports",
		Slug:              "bali-water",
		PrimaryCategoryID: category.ID,
		DestinationID:     destination.ID,
		PriceAmount:       250.00,
		DurationValue:     1,
		DurationUnit:      "days",
		GuideID:           guide.ID,
		IsActive:          true,
		IsListed:          true,
	}
	db.Create(&tour)

	booking := models.Booking{
		BookingReference:  "BK004",
		Status:            "pending",
		TravelerID:        traveler.ID,
		TravelerName:      "George Traveler",
		TourID:            tour.ID,
		GuideID:           guide.ID,
		TotalParticipants: 1,
		AdultCount:        1,
		StartDate:         time.Now(),
		TotalPrice:        250.00,
	}
	db.Create(&booking)

	app.Put("/bookings/:id/cancel", CancelBooking)

	cancelReq := struct {
		Reason       string  `json:"reason"`
		RefundAmount float64 `json:"refund_amount"`
	}{
		Reason:       "Changed plans",
		RefundAmount: 250.00,
	}
	body, _ := json.Marshal(cancelReq)
	req := httptest.NewRequest("PUT", "/bookings/1/cancel", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, 1)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestDeleteBooking(t *testing.T) {
	app := fiber.New()
	db := setupBookingTestDB()
	database.Database.Db = db

	// Setup
	traveler := models.User{
		FirstName:    "Henry",
		LastName:     "Traveler",
		Email:        "henry@example.com",
		PasswordHash: "password",
		Role:         "traveler",
	}
	db.Create(&traveler)

	guide := models.User{
		FirstName:    "Hannah",
		LastName:     "Guide",
		Email:        "hannah@example.com",
		PasswordHash: "password",
		Role:         "local-expert",
	}
	db.Create(&guide)

	category := models.TourCategory{Name: "Beach", Description: "Beach tours"}
	db.Create(&category)

	destination := models.TourDestination{City: "Phuket", Country: "Thailand"}
	db.Create(&destination)

	tour := models.Tour{
		Title:             "Phuket Beach",
		Description:       "Beach tour",
		Slug:              "phuket-beach",
		PrimaryCategoryID: category.ID,
		DestinationID:     destination.ID,
		PriceAmount:       180.00,
		DurationValue:     1,
		DurationUnit:      "days",
		GuideID:           guide.ID,
		IsActive:          true,
		IsListed:          true,
	}
	db.Create(&tour)

	booking := models.Booking{
		BookingReference:  "BK005",
		Status:            "pending",
		TravelerID:        traveler.ID,
		TravelerName:      "Henry Traveler",
		TourID:            tour.ID,
		GuideID:           guide.ID,
		TotalParticipants: 2,
		AdultCount:        2,
		StartDate:         time.Now(),
		TotalPrice:        360.00,
	}
	db.Create(&booking)

	app.Delete("/bookings/:id", DeleteBooking)

	req := httptest.NewRequest("DELETE", "/bookings/1", nil)
	resp, _ := app.Test(req, 1)

	assert.Equal(t, 200, resp.StatusCode)
}

// ============ BOOKING PARTICIPANT TESTS ============

func TestAddBookingParticipant(t *testing.T) {
	app := fiber.New()
	db := setupBookingTestDB()
	database.Database.Db = db

	// Setup - Create booking first
	traveler := models.User{
		FirstName:    "Iris",
		LastName:     "Traveler",
		Email:        "iris@example.com",
		PasswordHash: "password",
		Role:         "traveler",
	}
	db.Create(&traveler)

	guide := models.User{
		FirstName:    "Ian",
		LastName:     "Guide",
		Email:        "ian@example.com",
		PasswordHash: "password",
		Role:         "local-expert",
	}
	db.Create(&guide)

	category := models.TourCategory{Name: "Mixed", Description: "Mixed tours"}
	db.Create(&category)

	destination := models.TourDestination{City: "London", Country: "UK"}
	db.Create(&destination)

	tour := models.Tour{
		Title:             "London Tour",
		Description:       "Explore London",
		Slug:              "london-tour",
		PrimaryCategoryID: category.ID,
		DestinationID:     destination.ID,
		PriceAmount:       200.00,
		DurationValue:     2,
		DurationUnit:      "days",
		GuideID:           guide.ID,
		IsActive:          true,
		IsListed:          true,
	}
	db.Create(&tour)

	booking := models.Booking{
		BookingReference:  "BK006",
		Status:            "pending",
		TravelerID:        traveler.ID,
		TravelerName:      "Iris Traveler",
		TourID:            tour.ID,
		GuideID:           guide.ID,
		TotalParticipants: 2,
		AdultCount:        2,
		StartDate:         time.Now(),
		TotalPrice:        400.00,
	}
	db.Create(&booking)

	app.Post("/bookings/:booking_id/participants", AddBookingParticipant)

	participant := models.BookingParticipant{
		BookingID: booking.ID,
		Name:      "Iris Traveler",
		Age:       30,
	}
	body, _ := json.Marshal(participant)
	req := httptest.NewRequest("POST", "/bookings/1/participants", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, 1)
	assert.Equal(t, 201, resp.StatusCode)
}

func TestGetBookingParticipants(t *testing.T) {
	app := fiber.New()
	db := setupBookingTestDB()
	database.Database.Db = db

	// Setup
	traveler := models.User{
		FirstName:    "Jack",
		LastName:     "Traveler",
		Email:        "jack@example.com",
		PasswordHash: "password",
		Role:         "traveler",
	}
	db.Create(&traveler)

	guide := models.User{
		FirstName:    "Jill",
		LastName:     "Guide",
		Email:        "jill@example.com",
		PasswordHash: "password",
		Role:         "local-expert",
	}
	db.Create(&guide)

	category := models.TourCategory{Name: "Group", Description: "Group tours"}
	db.Create(&category)

	destination := models.TourDestination{City: "New York", Country: "USA"}
	db.Create(&destination)

	tour := models.Tour{
		Title:             "NYC Tour",
		Description:       "NYC experience",
		Slug:              "nyc-tour",
		PrimaryCategoryID: category.ID,
		DestinationID:     destination.ID,
		PriceAmount:       300.00,
		DurationValue:     3,
		DurationUnit:      "days",
		GuideID:           guide.ID,
		IsActive:          true,
		IsListed:          true,
	}
	db.Create(&tour)

	booking := models.Booking{
		BookingReference:  "BK007",
		Status:            "confirmed",
		TravelerID:        traveler.ID,
		TravelerName:      "Jack Traveler",
		TourID:            tour.ID,
		GuideID:           guide.ID,
		TotalParticipants: 3,
		AdultCount:        2,
		ChildCount:        1,
		StartDate:         time.Now(),
		TotalPrice:        600.00,
	}
	db.Create(&booking)

	// Add participants
	participant1 := models.BookingParticipant{
		BookingID: booking.ID,
		Name:      "Jack Traveler",
		Age:       35,
	}
	participant2 := models.BookingParticipant{
		BookingID: booking.ID,
		Name:      "Jane Traveler",
		Age:       33,
	}
	db.Create(&participant1)
	db.Create(&participant2)

	app.Get("/bookings/:booking_id/participants", GetBookingParticipants)

	req := httptest.NewRequest("GET", "/bookings/1/participants", nil)
	resp, _ := app.Test(req, 1)

	assert.Equal(t, 200, resp.StatusCode)

	var participants []models.BookingParticipant
	json.NewDecoder(resp.Body).Decode(&participants)
	assert.Greater(t, len(participants), 0)
}

// ============ PAYMENT TESTS ============

func TestCreatePayment(t *testing.T) {
	app := fiber.New()
	db := setupBookingTestDB()
	database.Database.Db = db

	// Setup
	traveler := models.User{
		FirstName:    "Kevin",
		LastName:     "Traveler",
		Email:        "kevin@example.com",
		PasswordHash: "password",
		Role:         "traveler",
	}
	db.Create(&traveler)

	guide := models.User{
		FirstName:    "Kate",
		LastName:     "Guide",
		Email:        "kate@example.com",
		PasswordHash: "password",
		Role:         "local-expert",
	}
	db.Create(&guide)

	category := models.TourCategory{Name: "Luxury", Description: "Luxury tours"}
	db.Create(&category)

	destination := models.TourDestination{City: "Dubai", Country: "UAE"}
	db.Create(&destination)

	tour := models.Tour{
		Title:             "Dubai Luxury",
		Description:       "Luxury experience",
		Slug:              "dubai-luxury",
		PrimaryCategoryID: category.ID,
		DestinationID:     destination.ID,
		PriceAmount:       500.00,
		DurationValue:     3,
		DurationUnit:      "days",
		GuideID:           guide.ID,
		IsActive:          true,
		IsListed:          true,
	}
	db.Create(&tour)

	booking := models.Booking{
		BookingReference:  "BK008",
		Status:            "pending",
		TravelerID:        traveler.ID,
		TravelerName:      "Kevin Traveler",
		TourID:            tour.ID,
		GuideID:           guide.ID,
		TotalParticipants: 1,
		AdultCount:        1,
		StartDate:         time.Now(),
		TotalPrice:        500.00,
	}
	db.Create(&booking)

	app.Post("/bookings/:booking_id/payments", CreatePayment)

	payment := models.Payment{
		BookingID:     booking.ID,
		PaymentMethod: "credit_card",
		Status:        "pending",
		TransactionID: "TXN123456",
		Amount:        500.00,
		Currency:      "USD",
	}
	body, _ := json.Marshal(payment)
	req := httptest.NewRequest("POST", "/bookings/1/payments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, 1)
	assert.Equal(t, 201, resp.StatusCode)
}

func TestGetPayments(t *testing.T) {
	app := fiber.New()
	db := setupBookingTestDB()
	database.Database.Db = db

	// Setup
	traveler := models.User{
		FirstName:    "Laura",
		LastName:     "Traveler",
		Email:        "laura@example.com",
		PasswordHash: "password",
		Role:         "traveler",
	}
	db.Create(&traveler)

	guide := models.User{
		FirstName:    "Leo",
		LastName:     "Guide",
		Email:        "leo@example.com",
		PasswordHash: "password",
		Role:         "local-expert",
	}
	db.Create(&guide)

	category := models.TourCategory{Name: "Budget", Description: "Budget tours"}
	db.Create(&category)

	destination := models.TourDestination{City: "Bangkok", Country: "Thailand"}
	db.Create(&destination)

	tour := models.Tour{
		Title:             "Budget Bangkok",
		Description:       "Budget tour",
		Slug:              "budget-bangkok",
		PrimaryCategoryID: category.ID,
		DestinationID:     destination.ID,
		PriceAmount:       50.00,
		DurationValue:     1,
		DurationUnit:      "days",
		GuideID:           guide.ID,
		IsActive:          true,
		IsListed:          true,
	}
	db.Create(&tour)

	booking := models.Booking{
		BookingReference:  "BK009",
		Status:            "confirmed",
		TravelerID:        traveler.ID,
		TravelerName:      "Laura Traveler",
		TourID:            tour.ID,
		GuideID:           guide.ID,
		TotalParticipants: 1,
		AdultCount:        1,
		StartDate:         time.Now(),
		TotalPrice:        50.00,
	}
	db.Create(&booking)

	payment := models.Payment{
		BookingID:     booking.ID,
		PaymentMethod: "bank_transfer",
		Status:        "completed",
		TransactionID: "TXN789012",
		Amount:        50.00,
		Currency:      "USD",
		PaidAt:        func() *time.Time { t := time.Now(); return &t }(),
	}
	db.Create(&payment)

	app.Get("/bookings/:booking_id/payments", GetPayments)

	req := httptest.NewRequest("GET", "/bookings/1/payments", nil)
	resp, _ := app.Test(req, 1)

	assert.Equal(t, 200, resp.StatusCode)

	var payments []models.Payment
	json.NewDecoder(resp.Body).Decode(&payments)
	assert.Greater(t, len(payments), 0)
}

// ============ REVIEW TESTS ============

func TestCreateReview(t *testing.T) {
	app := fiber.New()
	db := setupBookingTestDB()
	database.Database.Db = db

	// Setup
	reviewer := models.User{
		FirstName:    "Mark",
		LastName:     "Reviewer",
		Email:        "mark@example.com",
		PasswordHash: "password",
		Role:         "traveler",
	}
	db.Create(&reviewer)

	guide := models.User{
		FirstName:    "Monica",
		LastName:     "Guide",
		Email:        "monica@example.com",
		PasswordHash: "password",
		Role:         "local-expert",
	}
	db.Create(&guide)

	category := models.TourCategory{Name: "Review Test", Description: "Review test"}
	db.Create(&category)

	destination := models.TourDestination{City: "Barcelona", Country: "Spain"}
	db.Create(&destination)

	tour := models.Tour{
		Title:             "Barcelona Review",
		Description:       "Review tour",
		Slug:              "barcelona-review",
		PrimaryCategoryID: category.ID,
		DestinationID:     destination.ID,
		PriceAmount:       180.00,
		DurationValue:     2,
		DurationUnit:      "days",
		GuideID:           guide.ID,
		IsActive:          true,
		IsListed:          true,
	}
	db.Create(&tour)

	app.Post("/reviews", CreateReview)

	review := models.Review{
		TourID:       tour.ID,
		GuideID:      guide.ID,
		ReviewerID:   reviewer.ID,
		ReviewerName: "Mark Reviewer",
		Title:        "Great Tour!",
		Comment:      "Had a wonderful experience",
		Rating:       5,
		Verified:     true,
	}
	body, _ := json.Marshal(review)
	req := httptest.NewRequest("POST", "/reviews", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, 1)
	assert.Equal(t, 201, resp.StatusCode)
}

func TestGetReviews(t *testing.T) {
	app := fiber.New()
	db := setupBookingTestDB()
	database.Database.Db = db

	// Setup
	reviewer := models.User{
		FirstName:    "Nancy",
		LastName:     "Reviewer",
		Email:        "nancy@example.com",
		PasswordHash: "password",
		Role:         "traveler",
	}
	db.Create(&reviewer)

	guide := models.User{
		FirstName:    "Nathan",
		LastName:     "Guide",
		Email:        "nathan@example.com",
		PasswordHash: "password",
		Role:         "local-expert",
	}
	db.Create(&guide)

	category := models.TourCategory{Name: "Review Multiple", Description: "Review multiple"}
	db.Create(&category)

	destination := models.TourDestination{City: "Milan", Country: "Italy"}
	db.Create(&destination)

	tour := models.Tour{
		Title:             "Milan Review",
		Description:       "Review tour",
		Slug:              "milan-review",
		PrimaryCategoryID: category.ID,
		DestinationID:     destination.ID,
		PriceAmount:       200.00,
		DurationValue:     2,
		DurationUnit:      "days",
		GuideID:           guide.ID,
		IsActive:          true,
		IsListed:          true,
	}
	db.Create(&tour)

	review := models.Review{
		TourID:       tour.ID,
		GuideID:      guide.ID,
		ReviewerID:   reviewer.ID,
		ReviewerName: "Nancy Reviewer",
		Title:        "Excellent!",
		Comment:      "Perfect tour experience",
		Rating:       5,
		Verified:     true,
	}
	db.Create(&review)

	app.Get("/reviews", GetReviews)

	req := httptest.NewRequest("GET", "/reviews", nil)
	resp, _ := app.Test(req, 1)

	assert.Equal(t, 200, resp.StatusCode)

	var reviews []models.Review
	json.NewDecoder(resp.Body).Decode(&reviews)
	assert.Greater(t, len(reviews), 0)
}

// ============ EDGE CASE TESTS ============

func TestCreateBookingInvalidPayload(t *testing.T) {
	app := fiber.New()
	db := setupBookingTestDB()
	database.Database.Db = db

	app.Post("/bookings", CreateBooking)

	req := httptest.NewRequest("POST", "/bookings", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, 1)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestGetBookingsWithFiltering(t *testing.T) {
	app := fiber.New()
	db := setupBookingTestDB()
	database.Database.Db = db

	app.Get("/bookings", GetBookings)

	req := httptest.NewRequest("GET", "/bookings?status=confirmed&page=1&limit=5", nil)
	resp, _ := app.Test(req, 1)

	assert.Equal(t, 200, resp.StatusCode)
}
