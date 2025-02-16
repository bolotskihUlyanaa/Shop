package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/bolotskihUlyanaa/Shop/internal/models"
	"github.com/bolotskihUlyanaa/Shop/internal/service"
	"github.com/gin-gonic/gin"
)

// Структура для обрабатки http запросов
type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service}
}

// Инициализация эндпоинтов
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New() // Создание маршрутизатора
	api := router.Group("/api")
	{
		api.POST("/auth", h.Auth)
		api.POST("/sendCoin", AuthMiddleware(), h.SendCoins)
		api.GET("/info", AuthMiddleware(), h.Info)
		api.GET("/buy/:item", AuthMiddleware(), h.Buy)
	}
	return router
}

// Функция покупки мерча
func (h *Handler) Buy(c *gin.Context) {
	idEmployee, ex := c.Get("idEmployee") // Получение из контекста после аутентификации id
	id := idEmployee.(int)
	if !ex {
		log.Println("context is missing")
		NewErrorResponse(c, fmt.Errorf("context is missing %w", models.ErrInternal))
	}
	item := c.Param("item") // Получение параметра item (название предмета мерча) из пути запроса
	err := h.service.Buy(id, item)
	if err != nil { // Если мерча с таким названием нет или недостаточно средств
		NewErrorResponse(c, err)
	} else {
		c.Status(http.StatusOK)
	}
}

// Функция для отправки монет сотруднику
func (h *Handler) SendCoins(c *gin.Context) {
	idEmployee, ex := c.Get("idEmployee")
	id := idEmployee.(int)
	if !ex {
		log.Println("context is missing")
		NewErrorResponse(c, fmt.Errorf("context is missing %w", models.ErrInternal))
	}
	var input models.SendCoinRequest
	if err := validate(c, &input); err != nil { // Функция для декодирования json из тела запроса
		log.Println(err)
		NewErrorResponse(c, fmt.Errorf("invalid input data %w", models.ErrBadRequest))
	}
	err := h.service.SendCoins(id, input)
	if err != nil { // Если сотрудника с таким именем нет или недостаточно средств
		NewErrorResponse(c, err)
	} else {
		c.Status(http.StatusOK)
	}
}

// Функция для аутентификации и получение JWT
func (h *Handler) Auth(c *gin.Context) {
	var input models.User
	if err := validate(c, &input); err != nil { // Функция для декодирования json из тела запроса
		log.Println(err)
		NewErrorResponse(c, fmt.Errorf("invalid input data %w", models.ErrBadRequest))
	}
	token, err := h.service.Auth(input)
	if err != nil {
		NewErrorResponse(c, err)
	} else {
		c.JSON(http.StatusOK, map[string]interface{}{ // Отправка токена клиенту
			"token": token,
		})
	}
}

// Функция для получения информации о сотруднике (баланс, инвентарь, переводы)
func (h *Handler) Info(c *gin.Context) {
	idEmployee, ex := c.Get("idEmployee")
	id := idEmployee.(int)
	if !ex {
		log.Println("context is missing")
		NewErrorResponse(c, fmt.Errorf("context is missing %w", models.ErrInternal))
	}
	data, err := h.service.GetInfo(id)
	if err != nil {
		NewErrorResponse(c, err)
	} else {
		c.JSON(http.StatusOK, data)
	}
}

// Функция для декодирования полученных данных в input
func validate(c *gin.Context, input interface{}) error {
	if err := c.BindJSON(&input); err != nil {
		return err
	}
	return nil
}

// Функция для отправки сообщения об ошибке в формате json
func NewErrorResponse(c *gin.Context, err error) {
	// Отправка ответа клиенту с кодом статуса и описанием ошибки в теле
	var code int
	if errors.Is(err, models.ErrAuth) { //401
		code = http.StatusUnauthorized
	}
	if errors.Is(err, models.ErrBadRequest) { //400
		code = http.StatusBadRequest
	}
	if errors.Is(err, models.ErrInternal) { //500
		code = http.StatusInternalServerError
	}
	c.AbortWithStatusJSON(code, map[string]string{
		"description": err.Error(),
	})
}
