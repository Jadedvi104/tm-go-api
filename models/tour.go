package models

import (
	"time"

	"gorm.io/gorm"
)

// TourCategory represents tour categories (lookup table)
type TourCategory struct {
	ID          int8   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string `gorm:"type:varchar(100);not null;unique" json:"name"`
	Description string `gorm:"type:text" json:"description"`

	// Relationships
	Tours []Tour `gorm:"many2many:tour_category_mapping;" json:"tours,omitempty"`
}

func (TourCategory) TableName() string {
	return "tour_categories"
}

// TourDestination represents destination information
type TourDestination struct {
	ID        int64   `gorm:"primaryKey;autoIncrement" json:"id"`
	City      string  `gorm:"type:varchar(100);not null" json:"city"`
	Country   string  `gorm:"type:varchar(100);not null" json:"country"`
	Latitude  float64 `gorm:"type:decimal(10,8)" json:"latitude"`
	Longitude float64 `gorm:"type:decimal(11,8)" json:"longitude"`

	// Relationships
	Tours []Tour `gorm:"foreignKey:DestinationID" json:"tours,omitempty"`
}

func (TourDestination) TableName() string {
	return "tour_destinations"
}

// Tour represents tour information created by local experts
type Tour struct {
	ID                int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Title             string `gorm:"type:varchar(255);not null;index" json:"title"`
	Description       string `gorm:"type:text;not null" json:"description"`
	ShortDescription  string `gorm:"type:varchar(500)" json:"short_description"`
	Slug              string `gorm:"type:varchar(255);not null;unique;index" json:"slug"`
	PrimaryCategoryID int8   `gorm:"not null;index:idx_category" json:"primary_category_id"`
	DestinationID     int64  `gorm:"not null;index:idx_destination" json:"destination_id"`

	// Pricing
	PriceAmount    float64 `gorm:"type:decimal(10,2);not null" json:"price_amount"`
	PriceCurrency  string  `gorm:"type:varchar(10);default:'USD'" json:"price_currency"`
	PricePerPerson bool    `gorm:"default:true" json:"price_per_person"`

	// Duration
	DurationValue int    `gorm:"not null" json:"duration_value"`
	DurationUnit  string `gorm:"type:varchar(20);default:'days'" json:"duration_unit"` // 'days' or 'hours'

	// Participants
	MaxParticipants int `gorm:"default:null" json:"max_participants"`
	MinParticipants int `gorm:"default:null" json:"min_participants"`

	// Guide Information
	GuideID    int64  `gorm:"not null;index:idx_guide_id;constraint:OnDelete:RESTRICT" json:"guide_id"`
	GuideName  string `gorm:"type:varchar(255)" json:"guide_name"`
	GuideImage string `gorm:"type:varchar(255)" json:"guide_image"`

	// Ratings
	AverageRating *float64 `gorm:"type:decimal(3,2);default:0" json:"average_rating"`
	TotalReviews  int      `gorm:"default:0" json:"total_reviews"`

	// Status
	IsActive bool `gorm:"default:true" json:"is_active"`
	IsListed bool `gorm:"default:true" json:"is_listed"`

	CreatedAt time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"type:datetime;index" json:"-"`

	// Relationships
	Guide       User            `gorm:"foreignKey:GuideID" json:"guide,omitempty"`
	Category    TourCategory    `gorm:"foreignKey:PrimaryCategoryID" json:"category,omitempty"`
	Destination TourDestination `gorm:"foreignKey:DestinationID" json:"destination,omitempty"`
	Images      []TourImage     `gorm:"foreignKey:TourID" json:"images,omitempty"`
	Itineraries []TourItinerary `gorm:"foreignKey:TourID" json:"itineraries,omitempty"`
	Tags        []TourTag       `gorm:"foreignKey:TourID" json:"tags,omitempty"`
	Includes    []TourInclude   `gorm:"foreignKey:TourID" json:"includes,omitempty"`
	Excludes    []TourExclude   `gorm:"foreignKey:TourID" json:"excludes,omitempty"`
	Categories  []TourCategory  `gorm:"many2many:tour_category_mapping;" json:"categories,omitempty"`
	Bookings    []Booking       `gorm:"foreignKey:TourID" json:"bookings,omitempty"`
	Reviews     []Review        `gorm:"foreignKey:TourID" json:"reviews,omitempty"`
}

func (Tour) TableName() string {
	return "tours"
}

// TourImage stores images for tours
type TourImage struct {
	ID        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	TourID    int64          `gorm:"not null;index:idx_tour_id" json:"tour_id"`
	URL       string         `gorm:"type:varchar(500);not null" json:"url"`
	Alt       string         `gorm:"type:varchar(255)" json:"alt"`
	Caption   string         `gorm:"type:text" json:"caption"`
	IsMain    bool           `gorm:"default:false;index:idx_is_main" json:"is_main"`
	DeletedAt gorm.DeletedAt `gorm:"type:datetime;index" json:"-"`

	Tour Tour `gorm:"foreignKey:TourID" json:"tour,omitempty"`
}

func (TourImage) TableName() string {
	return "tour_images"
}

// TourTag stores tags for tours
type TourTag struct {
	ID        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	TourID    int64          `gorm:"not null;index:idx_tour_id" json:"tour_id"`
	Tag       string         `gorm:"type:varchar(100);not null;index:idx_tag" json:"tag"`
	DeletedAt gorm.DeletedAt `gorm:"type:datetime;index" json:"-"`

	Tour Tour `gorm:"foreignKey:TourID" json:"tour,omitempty"`
}

func (TourTag) TableName() string {
	return "tour_tags"
}

// TourInclude stores what is included in the tour
type TourInclude struct {
	ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	TourID      int64          `gorm:"not null;index:idx_tour_id" json:"tour_id"`
	IncludeItem string         `gorm:"type:varchar(255);not null" json:"include_item"`
	DeletedAt   gorm.DeletedAt `gorm:"type:datetime;index" json:"-"`

	Tour Tour `gorm:"foreignKey:TourID" json:"tour,omitempty"`
}

func (TourInclude) TableName() string {
	return "tour_includes"
}

// TourExclude stores what is excluded from the tour
type TourExclude struct {
	ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	TourID      int64          `gorm:"not null;index:idx_tour_id" json:"tour_id"`
	ExcludeItem string         `gorm:"type:varchar(255);not null" json:"exclude_item"`
	DeletedAt   gorm.DeletedAt `gorm:"type:datetime;index" json:"-"`

	Tour Tour `gorm:"foreignKey:TourID" json:"tour,omitempty"`
}

func (TourExclude) TableName() string {
	return "tour_excludes"
}

// TourItinerary stores daily itinerary information for tours
type TourItinerary struct {
	ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	TourID      int64          `gorm:"not null;index:idx_tour_id" json:"tour_id"`
	Day         int            `gorm:"not null" json:"day"`
	Title       string         `gorm:"type:varchar(255);not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	DeletedAt   gorm.DeletedAt `gorm:"type:datetime;index" json:"-"`

	Tour       Tour                    `gorm:"foreignKey:TourID" json:"tour,omitempty"`
	Activities []TourItineraryActivity `gorm:"foreignKey:ItineraryID" json:"activities,omitempty"`
	Meals      []TourItineraryMeal     `gorm:"foreignKey:ItineraryID" json:"meals,omitempty"`
}

func (TourItinerary) TableName() string {
	return "tour_itineraries"
}

// TourItineraryActivity stores activities for each itinerary day
type TourItineraryActivity struct {
	ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	ItineraryID int64          `gorm:"not null;index:idx_itinerary_id" json:"itinerary_id"`
	Activity    string         `gorm:"type:varchar(255);not null" json:"activity"`
	DeletedAt   gorm.DeletedAt `gorm:"type:datetime;index" json:"-"`

	Itinerary TourItinerary `gorm:"foreignKey:ItineraryID" json:"itinerary,omitempty"`
}

func (TourItineraryActivity) TableName() string {
	return "tour_itinerary_activities"
}

// TourItineraryMeal stores meal information for each itinerary day
type TourItineraryMeal struct {
	ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	ItineraryID int64          `gorm:"not null;index:idx_itinerary_id" json:"itinerary_id"`
	Meal        string         `gorm:"type:varchar(255);not null" json:"meal"`
	DeletedAt   gorm.DeletedAt `gorm:"type:datetime;index" json:"-"`

	Itinerary TourItinerary `gorm:"foreignKey:ItineraryID" json:"itinerary,omitempty"`
}

func (TourItineraryMeal) TableName() string {
	return "tour_itinerary_meals"
}
