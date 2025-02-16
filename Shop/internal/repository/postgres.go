package repository

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

// Структура хранит даные необходимые для подключения к Postgres
type Postgres struct {
	port     int
	host     string
	user     string
	password string
	dbName   string
	sslMode  string
}

func NewPostgres() (*sql.DB, error) {
	postgres, err := initPostgres() // Инициализация
	if err != nil {
		return nil, err
	}
	db, err := postgres.Connection() // Подключение
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Функция для инициализации структуры Postgres переменными окружения
func initPostgres() (*Postgres, error) {
	postgres := &Postgres{
		host:     os.Getenv("HOST_POSTGRES"),
		user:     os.Getenv("USERNAME_POSTGRES"),
		password: os.Getenv("PASSWORD_POSTGRES"),
		dbName:   os.Getenv("DBNAME_POSTGRES"),
		sslMode:  os.Getenv("SSLMODE_POSTGRES"),
	}
	port, err := strconv.Atoi(os.Getenv("PORT_POSTGRES"))
	if err != nil {
		return nil, fmt.Errorf("Invalid port: %w", err)
	}
	postgres.port = port
	return postgres, nil
}

// Функция для подключения к базе данных
func (p *Postgres) Connection() (*sql.DB, error) {
	// Открытие соединения с базой данных
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s password=%s sslmode=%s",
		p.host, p.port, p.user,
		p.dbName, "password", p.sslMode))
	if err != nil {
		return nil, fmt.Errorf("Open error: %w", err)
	}
	err = db.Ping() // Проверка подключения к базе данных
	if err != nil {
		return nil, fmt.Errorf("Connection error: %w", err)
	}
	return db, nil
}
