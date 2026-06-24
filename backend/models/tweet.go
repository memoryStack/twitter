package models

import "gorm.io/gorm"

type Tweet struct {
	Text     string `json:"text" gorm:"type:text;not null"`
	MediaURL string `json:"mediaUrl"`
	Likes    int64  `json:"likes" gorm:"default:0;not null"`
	UserID   uint   `json:"userId" gorm:"not null;index"` // FK -> users.id (enforced in repositories/postgres)
	gorm.Model
}
