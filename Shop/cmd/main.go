package main

import (
	"log"

	"github.com/bolotskihUlyanaa/Shop/internal/handler"
	"github.com/bolotskihUlyanaa/Shop/internal/repository"
	"github.com/bolotskihUlyanaa/Shop/internal/service"
)

func main() {
	// Внедрение зависимостей
	db, err := repository.NewPostgres()
	if err != nil {
		log.Fatal("Connection to db: ", err)
	}
	defer db.Close()
	repository := repository.NewShopRepository(db)
	service := service.NewShopService(repository)
	// Создаём экземпляр Handler
	handler := handler.NewHandler(service)
	// Инициализируем маршруты
	router := handler.InitRoutes()
	// Запускаем сервер
	router.Run(":8080")
}
