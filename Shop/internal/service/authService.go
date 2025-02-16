package service

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bolotskihUlyanaa/Shop/internal/models"
	"github.com/dgrijalva/jwt-go"
)

const tokenTTL = 12 * time.Hour // Время жизни токена

// Функция для аутентификации сотрудника по имени и паролю
func (s *ShopService) Auth(user models.User) (string, error) {
	id, passwordHash, err := s.repository.GetUser(user.Username)
	if err != nil { // Проверка зарегистрирован ли сотрудник
		if !errors.Is(err, models.ErrBadRequest) {
			return "", err
		}
	}
	hash, err := generateHash(user.Password) // Вычисление хеша
	if err != nil {
		log.Println(err)
		return "", fmt.Errorf("error with generate hash %w", models.ErrInternal)
	}
	if passwordHash != "" { // Если сотрудник зарегистрирован
		if hash != passwordHash { //Проверка пароля
			return "", fmt.Errorf("incorrect password %w", models.ErrAuth)
		}
	} else { // Если сотрудник не зарегистрирован
		id, err = s.repository.CreateUser(user.Username, hash) // Регистрация
		if err != nil {
			return "", err
		}
	}
	// Если сотрудник только что зарегистрировался или успешно прошел аутентификацию
	// Выдача токена
	token, err := GenerateToken(user, id)
	if err != nil {
		log.Println(err)
		return "", fmt.Errorf("error with generate token %w", models.ErrInternal)
	}
	return token, nil
}

// Генерация хеша
func generateHash(password string) (string, error) {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash), nil
}

// Генерация токена
func GenerateToken(user models.User, id int) (string, error) {
	signingKey := os.Getenv("SIGNING_KEY") // Ключ подписи
	if signingKey == "" {
		return "", errors.New("signing key is empty")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &models.ShopClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(), // Время истечения срока действия токена
		},
		UserId: id, // К стандартной информации добавляется id сотрудника
	})
	return token.SignedString([]byte(signingKey)) // Подпись токена
}
