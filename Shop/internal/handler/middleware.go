package handler

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bolotskihUlyanaa/Shop/internal/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Промежуточная функция для аутентификации на основе JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := getToken(c) // Извлечь токен
		if err != nil {
			log.Println(err)
			NewErrorResponse(c, fmt.Errorf("%w %w", err, models.ErrAuth))
		}
		id, err := parceToken(token) // Проверка токена и извлечение из него id сотрудника
		if err != nil {
			NewErrorResponse(c, fmt.Errorf("%w %w", err, models.ErrAuth))
		}
		c.Set("idEmployee", id) // Добавить idEmployee в контекст
		c.Next()                // Переход к следующему обработчику
	}
}

// Функция для извлечения токена из заголовка
func getToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization") // Чтение заголовка Authorization
	if authHeader == "" {
		return "", errors.New("authorization header is empty")
	}
	parts := strings.Split(authHeader, " ") // В bearer находится токен
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("invalid authorization header format")
	}
	return parts[1], nil // Извлечь токен
}

// Проверка токена и извлечение id
func parceToken(accessToken string) (int, error) {
	signingKey := os.Getenv("SIGNING_KEY") // Ключ подписи
	if signingKey == "" {
		return 0, errors.New("signing key is empty")
	}
	// Проверка и разбор токена
	token, err := jwt.ParseWithClaims(accessToken, &models.ShopClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok { // Проверка подписи
				return nil, errors.New("invalid signing method")
			}
			return []byte(signingKey), nil
		})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*models.ShopClaims) // Извлечь claims и привести к ShopClaims
	if !ok {
		return 0, errors.New("token claims are not of type *shopClaims")
	}
	return claims.UserId, nil
}
