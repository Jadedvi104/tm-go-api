package handlers

import (
	"strconv"
	"time"
	"tm-go-api/database"
	"tm-go-api/models"

	"github.com/gofiber/fiber/v2"
)

// ============ BOOKING HANDLERS ============

// CreateBooking creates a new booking
func CreateBooking(c *fiber.Ctx) error {
	booking := new(models.Booking)
	if err := c.BodyParser(booking); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid payload"})
	}

	// Validate tour exists
	var tour models.Tour
	if err := database.Database.Db.First(&tour, booking.TourID).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Tour not found"})
	}

	// Validate traveler exists
	var traveler models.User
	if err := database.Database.Db.First(&traveler, booking.TravelerID).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Traveler not found"})
	}

	if err := database.Database.Db.Create(&booking).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create booking"})
	}

	return c.Status(201).JSON(booking)
}

// GetBookings retrieves all bookings with filtering
func GetBookings(c *fiber.Ctx) error {
	var bookings []models.Booking
	query := database.Database.Db

	// Filtering
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if travelerID := c.Query("traveler_id"); travelerID != "" {
		query = query.Where("traveler_id = ?", travelerID)
	}
	if tourID := c.Query("tour_id"); tourID != "" {
		query = query.Where("tour_id = ?", tourID)
	}

	// Pagination
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	query.Offset(offset).Limit(limit).
		Preload("Traveler").
		Preload("Tour").
		Preload("Guide").
		Preload("Participants").
		Preload("Pricing").
		Preload("HotelDetails").
		Find(&bookings)

	return c.Status(200).JSON(bookings)
}

// GetBooking retrieves a booking by ID
func GetBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	var booking models.Booking

	result := database.Database.Db.
		Preload("Traveler").
		Preload("Tour").
		Preload("Guide").
		Preload("Participants").
		Preload("Pricing").
		Preload("HotelDetails").
		Preload("Preferences").
		Preload("Payments").
		Preload("Review").
		Find(&booking, id)

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Booking not found"})
	}

	return c.Status(200).JSON(booking)
}

// GetBookingByReference retrieves a booking by booking reference
func GetBookingByReference(c *fiber.Ctx) error {
	reference := c.Params("reference")
	var booking models.Booking

	result := database.Database.Db.
		Where("booking_reference = ?", reference).
		Preload("Traveler").
		Preload("Tour").
		Preload("Guide").
		Preload("Participants").
		Preload("Pricing").
		Preload("HotelDetails").
		Preload("Preferences").
		Preload("Payments").
		Find(&booking)

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Booking not found"})
	}

	return c.Status(200).JSON(booking)
}

// UpdateBooking updates an existing booking
func UpdateBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking := new(models.Booking)

	if err := c.BodyParser(booking); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid payload"})
	}

	result := database.Database.Db.Model(&models.Booking{}).Where("id = ?", id).Updates(booking)

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Booking not found"})
	}

	if err := result.Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not update booking"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Booking updated successfully"})
}

// CancelBooking cancels a booking
func CancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	cancelRequest := struct {
		Reason       string  `json:"reason"`
		RefundAmount float64 `json:"refund_amount"`
	}{}

	if err := c.BodyParser(&cancelRequest); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid payload"})
	}

	now := time.Now()
	result := database.Database.Db.Model(&models.Booking{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_cancelled":        true,
			"cancelled_at":        now,
			"cancellation_reason": cancelRequest.Reason,
			"refund_amount":       cancelRequest.RefundAmount,
			"refund_status":       "pending",
			"status":              "cancelled",
		})

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Booking not found"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Booking cancelled successfully"})
}

// DeleteBooking soft deletes a booking
func DeleteBooking(c *fiber.Ctx) error {
	id := c.Params("id")

	result := database.Database.Db.Delete(&models.Booking{}, id)

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Booking not found"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Booking deleted successfully"})
}

// ============ BOOKING PARTICIPANT HANDLERS ============

// AddBookingParticipant adds a participant to a booking
func AddBookingParticipant(c *fiber.Ctx) error {
	bookingID := c.Params("booking_id")
	participant := new(models.BookingParticipant)

	if err := c.BodyParser(participant); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid payload"})
	}

	participant.BookingID, _ = strconv.ParseInt(bookingID, 10, 64)

	if err := database.Database.Db.Create(&participant).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not add participant"})
	}

	return c.Status(201).JSON(participant)
}

// GetBookingParticipants retrieves all participants for a booking
func GetBookingParticipants(c *fiber.Ctx) error {
	bookingID := c.Params("booking_id")
	var participants []models.BookingParticipant

	database.Database.Db.Where("booking_id = ?", bookingID).Find(&participants)
	return c.Status(200).JSON(participants)
}

// DeleteBookingParticipant removes a participant from a booking
func DeleteBookingParticipant(c *fiber.Ctx) error {
	id := c.Params("id")

	result := database.Database.Db.Delete(&models.BookingParticipant{}, id)

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Participant not found"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Participant removed successfully"})
}

// ============ BOOKING PRICING HANDLERS ============

// CreateBookingPricing creates pricing details for a booking
func CreateBookingPricing(c *fiber.Ctx) error {
	pricing := new(models.BookingPricing)
	if err := c.BodyParser(pricing); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid payload"})
	}

	if err := database.Database.Db.Create(&pricing).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create pricing"})
	}

	return c.Status(201).JSON(pricing)
}

// GetBookingPricing retrieves pricing for a booking
func GetBookingPricing(c *fiber.Ctx) error {
	bookingID := c.Params("booking_id")
	var pricing models.BookingPricing

	result := database.Database.Db.Where("booking_id = ?", bookingID).First(&pricing)

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Pricing not found"})
	}

	return c.Status(200).JSON(pricing)
}

// UpdateBookingPricing updates pricing details for a booking
func UpdateBookingPricing(c *fiber.Ctx) error {
	bookingID := c.Params("booking_id")
	pricing := new(models.BookingPricing)

	if err := c.BodyParser(pricing); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid payload"})
	}

	result := database.Database.Db.Model(&models.BookingPricing{}).
		Where("booking_id = ?", bookingID).
		Updates(pricing)

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Pricing not found"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Pricing updated successfully"})
}

// ============ PAYMENT HANDLERS ============

// CreatePayment creates a payment record for a booking
func CreatePayment(c *fiber.Ctx) error {
	payment := new(models.Payment)
	if err := c.BodyParser(payment); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid payload"})
	}

	if err := database.Database.Db.Create(&payment).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create payment"})
	}

	return c.Status(201).JSON(payment)
}

// GetPayments retrieves all payments for a booking
func GetPayments(c *fiber.Ctx) error {
	bookingID := c.Params("booking_id")
	var payments []models.Payment

	database.Database.Db.Where("booking_id = ?", bookingID).Find(&payments)
	return c.Status(200).JSON(payments)
}

// GetPayment retrieves a specific payment
func GetPayment(c *fiber.Ctx) error {
	id := c.Params("id")
	var payment models.Payment

	result := database.Database.Db.Find(&payment, id)

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Payment not found"})
	}

	return c.Status(200).JSON(payment)
}

// UpdatePaymentStatus updates the status of a payment
func UpdatePaymentStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	statusUpdate := struct {
		Status string `json:"status"`
	}{}

	if err := c.BodyParser(&statusUpdate); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid payload"})
	}

	result := database.Database.Db.Model(&models.Payment{}).
		Where("id = ?", id).
		Update("status", statusUpdate.Status)

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Payment not found"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Payment status updated successfully"})
}

// ============ BOOKING HOTEL DETAILS HANDLERS ============

// CreateBookingHotelDetails creates hotel details for a booking
func CreateBookingHotelDetails(c *fiber.Ctx) error {
	hotelDetails := new(models.BookingHotelDetails)
	if err := c.BodyParser(hotelDetails); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid payload"})
	}

	if err := database.Database.Db.Create(&hotelDetails).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create hotel details"})
	}

	return c.Status(201).JSON(hotelDetails)
}

// GetBookingHotelDetails retrieves hotel details for a booking
func GetBookingHotelDetails(c *fiber.Ctx) error {
	bookingID := c.Params("booking_id")
	var hotelDetails models.BookingHotelDetails

	result := database.Database.Db.Where("booking_id = ?", bookingID).First(&hotelDetails)

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Hotel details not found"})
	}

	return c.Status(200).JSON(hotelDetails)
}

// UpdateBookingHotelDetails updates hotel details for a booking
func UpdateBookingHotelDetails(c *fiber.Ctx) error {
	bookingID := c.Params("booking_id")
	hotelDetails := new(models.BookingHotelDetails)

	if err := c.BodyParser(hotelDetails); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid payload"})
	}

	result := database.Database.Db.Model(&models.BookingHotelDetails{}).
		Where("booking_id = ?", bookingID).
		Updates(hotelDetails)

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Hotel details not found"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Hotel details updated successfully"})
}

// ============ REVIEW HANDLERS ============

// CreateReview creates a review for a tour
func CreateReview(c *fiber.Ctx) error {
	review := new(models.Review)
	if err := c.BodyParser(review); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid payload"})
	}

	if err := database.Database.Db.Create(&review).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create review"})
	}

	return c.Status(201).JSON(review)
}

// GetReviews retrieves all reviews with filtering
func GetReviews(c *fiber.Ctx) error {
	var reviews []models.Review
	query := database.Database.Db

	// Filtering
	if tourID := c.Query("tour_id"); tourID != "" {
		query = query.Where("tour_id = ?", tourID)
	}
	if guideID := c.Query("guide_id"); guideID != "" {
		query = query.Where("guide_id = ?", guideID)
	}
	if verified := c.Query("verified"); verified != "" {
		query = query.Where("verified = ?", verified == "true")
	}

	// Pagination
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	query.Offset(offset).Limit(limit).
		Preload("Tour").
		Preload("Guide").
		Preload("Reviewer").
		Preload("DetailedRatings").
		Preload("Images").
		Order("created_at DESC").
		Find(&reviews)

	return c.Status(200).JSON(reviews)
}

// GetReview retrieves a review by ID
func GetReview(c *fiber.Ctx) error {
	id := c.Params("id")
	var review models.Review

	result := database.Database.Db.
		Preload("Tour").
		Preload("Guide").
		Preload("Reviewer").
		Preload("DetailedRatings").
		Preload("Images").
		Find(&review, id)

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Review not found"})
	}

	return c.Status(200).JSON(review)
}

// UpdateReview updates a review
func UpdateReview(c *fiber.Ctx) error {
	id := c.Params("id")
	review := new(models.Review)

	if err := c.BodyParser(review); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid payload"})
	}

	result := database.Database.Db.Model(&models.Review{}).Where("id = ?", id).Updates(review)

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Review not found"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Review updated successfully"})
}

// DeleteReview soft deletes a review
func DeleteReview(c *fiber.Ctx) error {
	id := c.Params("id")

	result := database.Database.Db.Delete(&models.Review{}, id)

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Review not found"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Review deleted successfully"})
}
