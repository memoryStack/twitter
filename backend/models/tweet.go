package models

import "gorm.io/gorm"

type Tweet struct {
	Text     string      `json:"text" gorm:"type:text;not null"`
	MediaURL string      `json:"mediaUrl"`
	Likes    int64       `json:"likes" gorm:"default:0;not null"`
	UserID   uint        `json:"userId" gorm:"not null;index"` // FK -> users.id (enforced in repositories/postgres)
	User     User        `json:"-" gorm:"foreignKey:UserID"`
	Author   TweetAuthor `json:"user" gorm:"-"`
	gorm.Model
}

func (t *Tweet) AfterFind(_ *gorm.DB) error {
	if t.User.ID != 0 {
		t.Author = TweetAuthorFromUser(t.User)
	}
	return nil
}
