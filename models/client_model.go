package models

import "gorm.io/gorm"

type Client struct {
	gorm.Model
	IP        string
	Browser   string
	Version   string
	OS        string
	Device    string
	Origin    string
	API       string
}
