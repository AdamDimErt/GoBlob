package database

import (
	"awesomeProject5/pkg/models" // Импортируйте ваш пакет с моделями
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	dsn := "host=localhost user=postgres password=09012008 dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	// Автоматическая миграция для создания таблицы users
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		panic("failed to migrate database")
	}
	err = DB.AutoMigrate(&models.Post{})
	if err != nil {
		panic("failed to migrate database")
	}
	err = DB.AutoMigrate(&models.Comment{})
	if err != nil {
		panic("failed to migrate database")
	}
}
