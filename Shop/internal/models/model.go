package models

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

// Структура, которая передается в теле запроса для аутентификации
type User struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Структура, которая передается в теле запроса для перевода монет сотруднику
type SendCoinRequest struct {
	ToUser string `json:"toUser" binding:"required"`
	Amount int    `json:"amount" binding:"required"`
}

// Информация которая передается в токене
type ShopClaims struct {
	jwt.StandardClaims
	UserId int `json:"id_user"` // id пользователя
}

// Ошибки
var (
	ErrBadRequest = errors.New("") // 500
	ErrAuth       = errors.New("") // 401
	ErrInternal   = errors.New("") // 400
)
