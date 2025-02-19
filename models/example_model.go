package models

import (
	"time"

	"gorm.io/gorm"
)

type Example struct {
	gorm.Model
	
	ID      int64  `gorm:"primaryKey"`
	UUID    string `gorm:"not null;unique;index"`
	Message string `gorm:"not null"`

	CreatedAt time.Time  `gorm:"not null"`
	UpdatedAt time.Time  `gorm:"not null"`
	DeletedAt *time.Time `gorm:"index"`
}
