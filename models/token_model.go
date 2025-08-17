package models

import (
	"time"

	"gorm.io/gorm"
)

type RefreshToken struct {
	gorm.Model

	ID        uint      `gorm:"primaryKey"`
	UserUUID  string    `gorm:"index;not null"`
	Token     string    `gorm:"not null;unique"`
	ExpiresAt time.Time `gorm:"not null"`

	CreatedAt time.Time  `gorm:"not null"`
	UpdatedAt time.Time  `gorm:"not null"`
	DeletedAt *time.Time `gorm:"index"`
}

type BlacklistedToken struct {
	gorm.Model

	ID        uint      `gorm:"primaryKey"`
	Token     string    `gorm:"not null;unique;index"`
	ExpiresAt time.Time `gorm:"not null"`

	CreatedAt time.Time  `gorm:"not null"`
	UpdatedAt time.Time  `gorm:"not null"`
	DeletedAt *time.Time `gorm:"index"`
}

type VerificationToken struct {
	gorm.Model

	ID        uint      `gorm:"primaryKey"`
	UserUUID  string    `gorm:"index;not null"`
	Token     string    `gorm:"not null;index"`
	ExpiresAt time.Time `gorm:"not null"`

	CreatedAt time.Time  `gorm:"not null"`
	UpdatedAt time.Time  `gorm:"not null"`
	DeletedAt *time.Time `gorm:"index"`
}
