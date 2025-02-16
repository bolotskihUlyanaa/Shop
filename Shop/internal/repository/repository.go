package repository

import "github.com/bolotskihUlyanaa/Shop/internal/models"

// Контракт слоя репозиторий
type Repository interface {
	// Регистрация сотрудника
	CreateUser(name, password string) (int, error)
	// Получение информации о сотруднике (по имени получить id и пароль)
	GetUser(username string) (int, string, error)
	// Отправка монет, отправитель - id, получатель - data
	SendCoins(id int, data models.SendCoinRequest) error
	// Получение информации о сотруднике по id (баланс, операции, инвентарь)
	GetInfo(id int) (map[string]interface{}, error)
	// Покупка мерча item сотрудником id
	Buy(id int, item string) error
}
