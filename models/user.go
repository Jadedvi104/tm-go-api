package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // Soft Delete

	FirstName string `gorm:"not null" json:"first_name"`
	LastName  string `gorm:"not null" json:"last_name"`
	Email     string `gorm:"uniqueIndex;not null" json:"email"`
	Phone     string `json:"phone"`

	// Ensure we never return the password hash in JSON responses
	PasswordHash string `gorm:"not null" json:"-"`

	Role       string `gorm:"default:'traveler'" json:"role"`
	Bio        string `json:"bio"`
	IsVerified bool   `gorm:"default:false" json:"is_verified"`
	IsActive   bool   `gorm:"default:true" json:"is_active"`

	// Pointers allow these to be NULL in the database
	YearsOfExperience *int     `json:"years_of_experience"`
	AverageRating     *float64 `json:"average_rating"`
	TotalReviews      int      `gorm:"default:0" json:"total_reviews"`
}
