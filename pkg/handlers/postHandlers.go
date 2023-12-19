package handlers

import (
	"awesomeProject5/internal/database"
	"awesomeProject5/internal/jwt"
	"awesomeProject5/pkg/models"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

// GetTags обрабатывает получение уникальных тегов
func GetTags(c *fiber.Ctx) error {
	return nil
}

// GetAllPosts обрабатывает получение всех постов
func GetAllPosts(c *fiber.Ctx) error {
	var posts []models.Post
	result := database.DB.Preload("User").Find(&posts)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "не удалось получить статьи"})
	}

	return c.JSON(posts)
}
func UpdatePost(c *fiber.Ctx) error {
	// Получите ID поста из параметров маршрута и преобразуйте его в uint
	postID := c.Params("id")
	postIDUint, err := strconv.ParseUint(postID, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Неверный формат ID"})
	}

	// Найдите пост по его ID
	post, err := FindPostByID(database.DB, uint(postIDUint))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Пост не найден"})
	}

	// Распарсите данные для обновления из запроса
	updateData := new(models.Post)
	if err := c.BodyParser(updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Ошибка обработки данных"})
	}

	// Обновите поля поста с новыми данными
	post.Title = updateData.Title
	post.Text = updateData.Text
	post.ImageURL = updateData.ImageURL

	// Сохраните обновленный пост в базе данных
	result := database.DB.Save(&post)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Не удалось обновить пост"})
	}

	// Верните обновленный пост в ответе
	return c.JSON(post)
}

// GetPost обрабатывает получение одного поста по ID
func GetPost(c *fiber.Ctx) error {
	// Получите ID поста из параметров маршрута и преобразуйте его в uint
	postID := c.Params("id")
	postIDUint, err := strconv.ParseUint(postID, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Неверный формат ID"})
	}

	// Найдите пост по его ID
	post, err := FindPostByID(database.DB, uint(postIDUint))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Пост не найден"})
	}

	// Верните пост в ответе
	return c.JSON(post)
}

func FindPostByID(db *gorm.DB, postID uint) (*models.Post, error) {
	var post models.Post
	result := db.Preload("User").Where("id = ?", postID).First(&post)
	if result.Error != nil {
		return nil, result.Error
	}
	return &post, nil
}

// CreatePost обрабатывает создание нового поста
func CreatePost(c *fiber.Ctx) error {
	post := new(models.Post)
	if err := c.BodyParser(post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Ошибка обработки данных"})
	}

	// userID должен быть получен из JWT токена (если используется аутентификация)
	userID, err := getUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Неавторизованный доступ"})
	}
	post.UserID = userID

	// Генерация уникального текста, например, добавление текущей даты и времени
	post.Text = fmt.Sprintf("Содержимое вашего поста (%s)", time.Now())

	result := database.DB.Create(&post)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "не удалось создать пост"})
	}

	return c.JSON(post)
}

func getUserIDFromToken(c *fiber.Ctx) (uint, error) {
	authHeader := c.Get("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		return 0, errors.New("формат токена не соответствует стандарту 'Bearer [token]'")
	}
	userToken := splitToken[1]

	claims, err := jwt.VerifyJWT(userToken)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

// Другие обработчики (RemovePost, UpdatePost)
