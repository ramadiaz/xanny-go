package models

import (
	"time"
)

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey"`
	UserUUID  string    `gorm:"index;not null"`
	Token     string    `gorm:"not null;unique"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}

type BlacklistedToken struct {
	ID        uint      `gorm:"primaryKey"`
	Token     string    `gorm:"not null;unique;index"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}
