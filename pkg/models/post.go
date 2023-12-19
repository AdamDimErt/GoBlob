package models

import (
	"gorm.io/gorm"
)

// Post определяет структуру поста для блога
type Post struct {
	gorm.Model
	Title      string
	Text       string `gorm:"unique"`
	ViewsCount int
	ImageURL   string
	Comments   []Comment `gorm:"foreignKey:PostID"`
	UserID     uint
	User       User `gorm:"foreignKey:UserID"`
}
