package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	// 1. Map Go's "ID" to your SQL column "uid"
	ID uint `gorm:"column:uid;primaryKey;autoIncrement" json:"id"`

	FirstName string `gorm:"type:varchar(255);not null" json:"first_name"`
	LastName  string `gorm:"type:varchar(255);not null" json:"last_name"`
	Email     string `gorm:"type:varchar(255);unique;not null" json:"email"`
	Phone     string `gorm:"type:varchar(20)" json:"phone"`

	// Maps to password_hash automatically
	PasswordHash string `gorm:"type:varchar(255);not null" json:"password"`

	// 2. Default value handled by Go, Check constraint handled by DB
	Role string `gorm:"type:varchar(20);default:'traveler';not null" json:"role"`

	ProfileImage string `gorm:"type:varchar(255)" json:"profile_image"`
	Bio          string `gorm:"type:varchar(max)" json:"bio"` // Support MAX

	// 3. BIT in SQL maps to bool in Go
	IsVerified bool `gorm:"type:bit;default:0" json:"is_verified"`
	IsActive   bool `gorm:"type:bit;default:1" json:"is_active"`

	YearsOfExperience *int `json:"years_of_experience"`

	// 4. Decimal precision
	AverageRating *float64 `gorm:"type:decimal(3,2)" json:"average_rating"`
	TotalReviews  int      `gorm:"default:0" json:"total_reviews"`

	// 5. Use database default for creation, but Go handles updates usually
	CreatedAt time.Time      `gorm:"type:datetime2;default:SYSDATETIME()" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:datetime2;default:SYSDATETIME()" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"type:datetime2;index" json:"-"`
}

// TableName overrides the default pluralization if you want to be extra safe,
// though GORM usually guesses "users" correctly.
func (User) TableName() string {
	return "users"
}
