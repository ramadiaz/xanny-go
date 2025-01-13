package models

import "gorm.io/gorm"

type Example struct {
	gorm.Model
	ID      int64  `gorm:"primaryKey"`
	Message string `gorm:"not null"`
}
