package routes

import (
	"tm-go-api/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Grouping routes is cleaner
	api := app.Group("/api/v1")

	// User routes
	user := api.Group("/users")
	user.Post("/create", handlers.CreateUser)
	user.Get("/", handlers.GetUsers)
	user.Get("/:id", handlers.GetUser)
	user.Put("/:id", handlers.UpdateUser)
	user.Delete("/:id", handlers.DeleteUser)

	// Tour routes
	tours := api.Group("/tours")
	tours.Post("/", handlers.CreateTour)
	tours.Get("/", handlers.GetTours)
	tours.Get("/:id", handlers.GetTour)
	tours.Put("/:id", handlers.UpdateTour)
	tours.Delete("/:id", handlers.DeleteTour)

	// Tour Category routes
	categories := api.Group("/tour-categories")
	categories.Post("/", handlers.CreateTourCategory)
	categories.Get("/", handlers.GetTourCategories)
	categories.Get("/:id", handlers.GetTourCategory)
	categories.Put("/:id", handlers.UpdateTourCategory)
	categories.Delete("/:id", handlers.DeleteTourCategory)

	// Tour Destination routes
	destinations := api.Group("/tour-destinations")
	destinations.Post("/", handlers.CreateTourDestination)
	destinations.Get("/", handlers.GetTourDestinations)
	destinations.Get("/:id", handlers.GetTourDestination)
	destinations.Put("/:id", handlers.UpdateTourDestination)
	destinations.Delete("/:id", handlers.DeleteTourDestination)

	// Booking routes
	bookings := api.Group("/bookings")
	bookings.Post("/", handlers.CreateBooking)
	bookings.Get("/", handlers.GetBookings)
	bookings.Get("/:id", handlers.GetBooking)
	bookings.Get("/reference/:reference", handlers.GetBookingByReference)
	bookings.Put("/:id", handlers.UpdateBooking)
	bookings.Put("/:id/cancel", handlers.CancelBooking)
	bookings.Delete("/:id", handlers.DeleteBooking)

	// Booking Participant routes
	participants := api.Group("/bookings/:booking_id/participants")
	participants.Post("/", handlers.AddBookingParticipant)
	participants.Get("/", handlers.GetBookingParticipants)
	participants.Delete("/:id", handlers.DeleteBookingParticipant)

	// Booking Pricing routes
	pricing := api.Group("/bookings/:booking_id/pricing")
	pricing.Post("/", handlers.CreateBookingPricing)
	pricing.Get("/", handlers.GetBookingPricing)
	pricing.Put("/", handlers.UpdateBookingPricing)

	// Payment routes
	payments := api.Group("/bookings/:booking_id/payments")
	payments.Post("/", handlers.CreatePayment)
	payments.Get("/", handlers.GetPayments)
	payments.Get("/:id", handlers.GetPayment)
	payments.Put("/:id/status", handlers.UpdatePaymentStatus)

	// Booking Hotel Details routes
	hotelDetails := api.Group("/bookings/:booking_id/hotel-details")
	hotelDetails.Post("/", handlers.CreateBookingHotelDetails)
	hotelDetails.Get("/", handlers.GetBookingHotelDetails)
	hotelDetails.Put("/", handlers.UpdateBookingHotelDetails)

	// Review routes
	reviews := api.Group("/reviews")
	reviews.Post("/", handlers.CreateReview)
	reviews.Get("/", handlers.GetReviews)
	reviews.Get("/:id", handlers.GetReview)
	reviews.Put("/:id", handlers.UpdateReview)
	reviews.Delete("/:id", handlers.DeleteReview)
}
