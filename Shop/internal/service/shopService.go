package service

import (
	"github.com/bolotskihUlyanaa/Shop/internal/models"
	"github.com/bolotskihUlyanaa/Shop/internal/repository"
)

// Реализация слоя сервис
type ShopService struct {
	Service
	repository repository.Repository
}

func NewShopService(repository repository.Repository) *ShopService {
	return &ShopService{repository: repository}
}

// Функция для отправки монет сотруднику data от id
func (s *ShopService) SendCoins(id int, data models.SendCoinRequest) error {
	return s.repository.SendCoins(id, data)
}

// Функция для получения информации о сотруднике id
func (s *ShopService) GetInfo(id int) (map[string]interface{}, error) {
	return s.repository.GetInfo(id)
}

// Функция для покупки мерча item сотрудником id
func (s *ShopService) Buy(id int, item string) error {
	return s.repository.Buy(id, item)
}
