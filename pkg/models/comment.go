package models

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	Text   string
	UserID uint
	User   User `gorm:"foreignKey:UserID"`
	PostID uint
	Post   Post `gorm:"foreignKey:PostID"`
}
