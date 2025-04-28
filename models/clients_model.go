package models

import "gorm.io/gorm"

type Clients struct {
	gorm.Model
	IP        string
	Browser   string
	Version   string
	OS        string
	Device    string
	Origin    string
	API       string
}
