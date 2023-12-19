package models

import (
	"gorm.io/gorm"
)

// User определяет структуру пользователя для вашего приложения
type User struct {
	gorm.Model
	FullName     string
	Role         string
	Email        string `gorm:"unique"`
	PasswordHash string
	AvatarURL    string
	Posts        []Post `gorm:"foreignKey:UserID"`
}
