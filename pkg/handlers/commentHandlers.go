package handlers

import (
	"awesomeProject5/internal/database"
	"awesomeProject5/pkg/models"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

// CreateComment обрабатывает добавление комментария к посту
func CreateComment(c *fiber.Ctx) error {
	// Проверяем авторизацию пользователя
	userID, err := getUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Неавторизованный доступ"})
	}

	// Получаем ID поста
	postID := c.Params("postID")
	postIDUint, err := strconv.ParseUint(postID, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Неверный формат ID поста"})
	}

	// Создаем новый комментарий
	comment := new(models.Comment)
	if err := c.BodyParser(comment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Ошибка обработки данных"})
	}

	comment.UserID = userID
	comment.PostID = uint(postIDUint)

	// Сохраняем комментарий в базе данных
	result := database.DB.Create(&comment)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Не удалось добавить комментарий"})
	}

	return c.JSON(comment)
}

// UpdateComment обновляет комментарий
func UpdateComment(c *fiber.Ctx) error {
	commentID := c.Params("id")
	commentIDUint, err := strconv.ParseUint(commentID, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Неверный формат ID комментария"})
	}

	// Проверяем права пользователя

	var comment models.Comment
	result := database.DB.First(&comment, commentIDUint)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Комментарий не найден"})
	}

	if err := c.BodyParser(&comment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Ошибка обработки данных"})
	}

	database.DB.Save(&comment)

	return c.JSON(comment)
}
func GetAllCommentsByPost(c *fiber.Ctx) error {
	postID := c.Params("postID")
	postIDUint, err := strconv.ParseUint(postID, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Неверный формат ID поста"})
	}

	var comments []models.Comment
	result := database.DB.Where("post_id = ?", postIDUint).Find(&comments)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Ошибка при получении комментариев"})
	}

	return c.JSON(comments)
}

// DeleteComment удаляет комментарий
func DeleteComment(c *fiber.Ctx) error {
	commentID := c.Params("id")
	commentIDUint, err := strconv.ParseUint(commentID, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Неверный формат ID комментария"})
	}

	// Проверяем права пользователя
	// ...

	result := database.DB.Delete(&models.Comment{}, commentIDUint)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Комментарий не найден или ошибка при удалении"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
