package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/bolotskihUlyanaa/Shop/internal/models"
)

const coins = 1000 // Первоначальный баланс сотрудника

// Структура реализует слой репозиторий
type ShopRepository struct {
	Repository
	db *sql.DB
}

func NewShopRepository(db *sql.DB) *ShopRepository {
	return &ShopRepository{db: db}
}

// Получение id и пароля сотрудника по имени
func (s *ShopRepository) GetUser(username string) (int, string, error) {
	row := s.db.QueryRow("SELECT id, password_hash FROM Employees "+
		"WHERE name=$1;", username)
	var id int
	var password string
	if err := row.Scan(&id, &password); err != nil {
		log.Println(err)
		return 0, "", fmt.Errorf("user with name: %s not found%w", username, models.ErrBadRequest)
	}
	return id, password, nil
}

// Регистрация сотрудника
func (s *ShopRepository) CreateUser(name, password string) (int, error) {
	row := s.db.QueryRow("INSERT INTO Employees (name, password_hash, coins) "+
		"VALUES ($1, $2, $3) RETURNING id;", name, password, coins)
	var id int
	if err := row.Scan(&id); err != nil {
		log.Println(err)
		return 0, fmt.Errorf("create user: %w", models.ErrInternal)
	}
	return id, nil
}
