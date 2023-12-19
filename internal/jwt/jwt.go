package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var jwtKey = []byte("your_secret_key") // Замените на ваш секретный ключ

type Claims struct {
	UserID uint `json:"userId"`
	jwt.RegisteredClaims
}

// GenerateJWT генерирует новый JWT токен
func GenerateJWT(userId uint) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// VerifyJWT проверяет JWT токен
func VerifyJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		// Логирование ошибки
		fmt.Printf("Error parsing token: %v\n", err)
		return nil, err
	}

	if !token.Valid {
		// Логирование невалидности токена
		fmt.Println("Token is invalid")
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
