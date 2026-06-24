package models

import "gorm.io/gorm"

type User struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Email         string `json:"email" gorm:"unique"`
	PhoneNumber   string `json:"phone_number"`
	EmailVerified bool   `json:"email_verified"`
	Image         string `json:"image_url"`
	Auth0ID       string `json:"auth0_id" gorm:"not null;unique"`
	gorm.Model
}
