package handlers

import (
	"awesomeProject5/internal/database"
	"awesomeProject5/internal/jwt"
	"awesomeProject5/pkg/models"
	"awesomeProject5/pkg/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

// RegisterUser обрабатывает регистрацию пользователя
func RegisterUser(c *fiber.Ctx) error {
	user := new(models.User)

	fmt.Println(user)

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	hashedPassword, err := utils.HashPassword(user.PasswordHash)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось хешировать пароль"})
	}
	user.PasswordHash = hashedPassword

	// Создаем пользователя в базе данных
	result := database.DB.Create(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось создать пользователя"})
	}

	// Генерируем JWT токен для нового пользователя
	token, err := jwt.GenerateJWT(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось создать токен"})
	}

	// Возвращаем данные пользователя и токен в ответе
	return c.JSON(fiber.Map{"user": user, "token": token})
}

// LoginUser обрабатывает вход пользователя
func LoginUser(c *fiber.Ctx) error {
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var foundUser models.User
	database.DB.Where("email = ?", user.Email).First(&foundUser)

	if !utils.CheckPasswordHash(user.PasswordHash, foundUser.PasswordHash) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Неправильный пароль"})
	}

	token, err := jwt.GenerateJWT(foundUser.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось создать токен"})
	}

	return c.JSON(fiber.Map{"token": token})
}

// GetMe обрабатывает получение информации о текущем пользователе
func GetMe(c *fiber.Ctx) error {
	// Извлекаем токен из запроса
	userToken := c.Get("Authorization")

	// Проверяем и декодируем токен
	claims, err := jwt.VerifyJWT(userToken)
	if err != nil {
		//return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Неверный токен"})
		fmt.Println("as")
	}

	var user models.User
	result := database.DB.First(&user, claims.UserID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "пользователь не найден"})
	}

	// Преобразуем данные пользователя, исключая хеш пароля
	userData := map[string]interface{}{
		"ID":        user.ID,
		"CreatedAt": user.CreatedAt,
		"UpdatedAt": user.UpdatedAt,
		"FullName":  user.FullName,
		"Email":     user.Email,
	}

	return c.JSON(userData)
}
