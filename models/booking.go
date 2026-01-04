package models

import (
	"time"

	"gorm.io/gorm"
)

// Booking stores booking information for tours
type Booking struct {
	ID               int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	BookingReference string `gorm:"type:varchar(50);not null;unique;index" json:"booking_reference"`
	Status           string `gorm:"type:varchar(20);default:'pending';index:idx_status" json:"status"` // 'pending', 'confirmed', 'completed', 'cancelled'

	// Traveler Information
	TravelerID    int64  `gorm:"not null;index:idx_traveler_id;constraint:OnDelete:RESTRICT" json:"traveler_id"`
	TravelerName  string `gorm:"type:varchar(255)" json:"traveler_name"`
	TravelerEmail string `gorm:"type:varchar(255)" json:"traveler_email"`
	TravelerPhone string `gorm:"type:varchar(20)" json:"traveler_phone"`

	// Tour Information
	TourID    int64  `gorm:"not null;index:idx_tour_id;constraint:OnDelete:RESTRICT" json:"tour_id"`
	TourTitle string `gorm:"type:varchar(255)" json:"tour_title"`
	GuideID   int64  `gorm:"not null;index:idx_guide_id;constraint:OnDelete:RESTRICT" json:"guide_id"`

	// Participants
	TotalParticipants int `gorm:"not null" json:"total_participants"`
	AdultCount        int `gorm:"not null" json:"adult_count"`
	ChildCount        int `gorm:"default:0" json:"child_count"`

	// Dates
	StartDate time.Time `gorm:"type:date;not null;index:idx_start_date" json:"start_date"`
	EndDate   time.Time `gorm:"type:date" json:"end_date"`

	// Total Pricing
	TotalPrice float64 `gorm:"type:decimal(12,2);not null" json:"total_price"`
	Currency   string  `gorm:"type:varchar(10);default:'USD'" json:"currency"`

	// Special Requirements
	SpecialRequests string `gorm:"type:text" json:"special_requests"`
	Notes           string `gorm:"type:text" json:"notes"`

	// Cancellation Info
	IsCancelled        bool       `gorm:"default:false" json:"is_cancelled"`
	CancelledAt        *time.Time `gorm:"type:datetime" json:"cancelled_at"`
	CancellationReason string     `gorm:"type:text" json:"cancellation_reason"`
	RefundAmount       *float64   `gorm:"type:decimal(12,2)" json:"refund_amount"`
	RefundStatus       string     `gorm:"type:varchar(50)" json:"refund_status"`

	CreatedAt time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP;index:idx_created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"type:datetime;index" json:"-"`

	// Relationships
	Traveler     User                 `gorm:"foreignKey:TravelerID" json:"traveler,omitempty"`
	Tour         Tour                 `gorm:"foreignKey:TourID" json:"tour,omitempty"`
	Guide        User                 `gorm:"foreignKey:GuideID" json:"guide,omitempty"`
	Participants []BookingParticipant `gorm:"foreignKey:BookingID" json:"participants,omitempty"`
	Pricing      *BookingPricing      `gorm:"foreignKey:BookingID" json:"pricing,omitempty"`
	HotelDetails *BookingHotelDetails `gorm:"foreignKey:BookingID" json:"hotel_details,omitempty"`
	Preferences  []BookingPreference  `gorm:"foreignKey:BookingID" json:"preferences,omitempty"`
	Payments     []Payment            `gorm:"foreignKey:BookingID" json:"payments,omitempty"`
	Review       *Review              `gorm:"foreignKey:BookingID" json:"review,omitempty"`
}

func (Booking) TableName() string {
	return "bookings"
}

// BookingParticipant stores detailed participant information for bookings
type BookingParticipant struct {
	ID             int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	BookingID      int64          `gorm:"not null;index:idx_booking_id" json:"booking_id"`
	Name           string         `gorm:"type:varchar(255);not null" json:"name"`
	Age            *int           `json:"age"`
	PassportNumber string         `gorm:"type:varchar(50)" json:"passport_number"`
	DeletedAt      gorm.DeletedAt `gorm:"type:datetime;index" json:"-"`

	Booking Booking `gorm:"foreignKey:BookingID" json:"booking,omitempty"`
}

func (BookingParticipant) TableName() string {
	return "booking_participants"
}

// BookingPricing stores detailed pricing breakdown for bookings
type BookingPricing struct {
	ID                 int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	BookingID          int64          `gorm:"not null;unique;index:idx_booking_id" json:"booking_id"`
	BasePrice          float64        `gorm:"type:decimal(10,2);not null" json:"base_price"`
	PricePerPerson     float64        `gorm:"type:decimal(10,2);not null" json:"price_per_person"`
	Subtotal           float64        `gorm:"type:decimal(12,2);not null" json:"subtotal"`
	Tax                *float64       `gorm:"type:decimal(12,2)" json:"tax"`
	DiscountAmount     *float64       `gorm:"type:decimal(12,2)" json:"discount_amount"`
	DiscountPercentage *float64       `gorm:"type:decimal(5,2)" json:"discount_percentage"`
	DiscountCode       string         `gorm:"type:varchar(50)" json:"discount_code"`
	TotalPrice         float64        `gorm:"type:decimal(12,2);not null" json:"total_price"`
	Currency           string         `gorm:"type:varchar(10);default:'USD'" json:"currency"`
	DeletedAt          gorm.DeletedAt `gorm:"type:datetime;index" json:"-"`

	Booking Booking `gorm:"foreignKey:BookingID" json:"booking,omitempty"`
}

func (BookingPricing) TableName() string {
	return "booking_pricing"
}

// BookingHotelDetails stores hotel accommodation details for bookings
type BookingHotelDetails struct {
	ID         int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	BookingID  int64          `gorm:"not null;unique;index:idx_booking_id" json:"booking_id"`
	HotelName  string         `gorm:"type:varchar(255)" json:"hotel_name"`
	RoomNumber string         `gorm:"type:varchar(50)" json:"room_number"`
	Address    string         `gorm:"type:text" json:"address"`
	DeletedAt  gorm.DeletedAt `gorm:"type:datetime;index" json:"-"`

	Booking Booking `gorm:"foreignKey:BookingID" json:"booking,omitempty"`
}

func (BookingHotelDetails) TableName() string {
	return "booking_hotel_details"
}

// BookingPreference stores preferences for bookings
type BookingPreference struct {
	ID         int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	BookingID  int64          `gorm:"not null;index:idx_booking_id" json:"booking_id"`
	Preference string         `gorm:"type:varchar(255);not null" json:"preference"`
	DeletedAt  gorm.DeletedAt `gorm:"type:datetime;index" json:"-"`

	Booking Booking `gorm:"foreignKey:BookingID" json:"booking,omitempty"`
}

func (BookingPreference) TableName() string {
	return "booking_preferences"
}

// Payment stores payment information for bookings
type Payment struct {
	ID            int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	BookingID     int64      `gorm:"not null;index:idx_booking_id;constraint:OnDelete:CASCADE" json:"booking_id"`
	PaymentMethod string     `gorm:"type:varchar(50);not null" json:"payment_method"`
	Status        string     `gorm:"type:varchar(20);default:'pending';index:idx_status" json:"status"` // 'pending', 'completed', 'failed', 'refunded'
	TransactionID string     `gorm:"type:varchar(255);unique;index" json:"transaction_id"`
	Amount        float64    `gorm:"type:decimal(12,2);not null" json:"amount"`
	Currency      string     `gorm:"type:varchar(10);default:'USD'" json:"currency"`
	PaidAt        *time.Time `gorm:"type:datetime;index:idx_paid_at" json:"paid_at"`

	CreatedAt time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"type:datetime;index" json:"-"`

	Booking Booking `gorm:"foreignKey:BookingID" json:"booking,omitempty"`
}

func (Payment) TableName() string {
	return "payments"
}

// Review stores reviews and ratings for tours and guides
type Review struct {
	ID            int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	TourID        int64  `gorm:"not null;index:idx_tour_id;constraint:OnDelete:RESTRICT" json:"tour_id"`
	GuideID       int64  `gorm:"not null;index:idx_guide_id;constraint:OnDelete:RESTRICT" json:"guide_id"`
	BookingID     *int64 `gorm:"index;constraint:OnDelete:SET NULL" json:"booking_id"`
	ReviewerID    int64  `gorm:"not null;index:idx_reviewer_id;constraint:OnDelete:RESTRICT" json:"reviewer_id"`
	ReviewerName  string `gorm:"type:varchar(255)" json:"reviewer_name"`
	ReviewerImage string `gorm:"type:varchar(255)" json:"reviewer_image"`

	Title   string `gorm:"type:varchar(255);not null" json:"title"`
	Comment string `gorm:"type:text;not null" json:"comment"`
	Rating  int    `gorm:"not null;index:idx_rating" json:"rating"` // 1-5

	HelpfulCount int  `gorm:"default:0" json:"helpful_count"`
	Verified     bool `gorm:"default:false;index:idx_verified" json:"verified"`

	GuideResponse     string     `gorm:"type:text" json:"guide_response"`
	GuideResponseDate *time.Time `gorm:"type:datetime" json:"guide_response_date"`

	CreatedAt time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP;index:idx_created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"type:datetime;index" json:"-"`

	// Relationships
	Tour            Tour                   `gorm:"foreignKey:TourID" json:"tour,omitempty"`
	Guide           User                   `gorm:"foreignKey:GuideID" json:"guide,omitempty"`
	Booking         *Booking               `gorm:"foreignKey:BookingID" json:"booking,omitempty"`
	Reviewer        User                   `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
	DetailedRatings *ReviewDetailedRatings `gorm:"foreignKey:ReviewID" json:"detailed_ratings,omitempty"`
	Images          []ReviewImage          `gorm:"foreignKey:ReviewID" json:"images,omitempty"`
}

func (Review) TableName() string {
	return "reviews"
}

// ReviewDetailedRatings stores breakdown of ratings for different aspects
type ReviewDetailedRatings struct {
	ID            int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	ReviewID      int64          `gorm:"not null;unique;index" json:"review_id"`
	Accuracy      *int           `json:"accuracy"` // 1-5
	Communication *int           `json:"communication"`
	Cleanliness   *int           `json:"cleanliness"`
	Location      *int           `json:"location"`
	Value         *int           `json:"value"`
	DeletedAt     gorm.DeletedAt `gorm:"type:datetime;index" json:"-"`

	Review Review `gorm:"foreignKey:ReviewID" json:"review,omitempty"`
}

func (ReviewDetailedRatings) TableName() string {
	return "review_detailed_ratings"
}

// ReviewImage stores images attached to reviews
type ReviewImage struct {
	ID        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	ReviewID  int64          `gorm:"not null;index:idx_review_id;constraint:OnDelete:CASCADE" json:"review_id"`
	ImageURL  string         `gorm:"type:varchar(500);not null" json:"image_url"`
	DeletedAt gorm.DeletedAt `gorm:"type:datetime;index" json:"-"`

	Review Review `gorm:"foreignKey:ReviewID" json:"review,omitempty"`
}

func (ReviewImage) TableName() string {
	return "review_images"
}
