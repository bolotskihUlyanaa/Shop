package service

import "github.com/bolotskihUlyanaa/Shop/internal/models"

// Определение функций слоя сервис
type Service interface {
	// Аутентификация сотрудника user
	Auth(user models.User) (string, error)
	// Отправка монет, отправитель - id, получатель - data
	SendCoins(id int, data models.SendCoinRequest) error
	// Получение информации о сотруднике id (баланс, операции, инвентарь)
	GetInfo(id int) (map[string]interface{}, error)
	// Покупка мерча item сотрудником id
	Buy(id int, item string) error
}
