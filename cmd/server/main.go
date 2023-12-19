package main

import (
	"awesomeProject5/internal/database"
	"awesomeProject5/pkg/handlers"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())
	database.ConnectDB()
	setupRoutes(app)
	log.Fatal(app.Listen(":1234"))
}

func setupRoutes(app *fiber.App) {
	app.Post("/auth/register", handlers.RegisterUser)
	app.Post("/auth/login", handlers.LoginUser)
	app.Get("/auth/me", handlers.GetMe) // Новый маршрут для получения информации о текущем пользователе

	app.Get("/posts", handlers.GetAllPosts)
	app.Get("/posts/:id", handlers.GetPost)
	app.Post("/posts", handlers.CreatePost)
	app.Put("/posts/:id", handlers.UpdatePost)

	// Добавьте маршруты для работы с комментариями
	app.Post("/posts/:postID/comments", handlers.CreateComment)
	app.Get("/posts/:postID/comments", handlers.GetAllCommentsByPost)
	app.Put("/comments/:id", handlers.UpdateComment)
	app.Delete("/comments/:id", handlers.DeleteComment)

	// Маршрут для загрузки изображений
	// app.Post("/upload", handlers.UploadImage)
}
