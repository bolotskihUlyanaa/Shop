package repository

import (
	"fmt"
	"log"

	"github.com/bolotskihUlyanaa/Shop/internal/models"
)

// Функция для получения информации о сотруднике (баланс, операции, инвентарь)
func (s *ShopRepository) GetInfo(id int) (map[string]interface{}, error) {
	response := make(map[string]interface{})
	coins, err := s.getCoins(id) // Баланс
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("cant get user coins %w", models.ErrInternal)
	}
	response["coins"] = coins
	inv, err := s.getItem(fmt.Sprintf("SELECT name, NUM FROM "+ // Инвертарь
		"(SELECT COUNT(*) AS NUM, id_merch "+
		"FROM Buy WHERE id_employees=%d GROUP BY id_merch) AS Subquery "+
		"INNER JOIN Merch ON Subquery.id_merch = Merch.id;", id))
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("cant get user inventory %w", models.ErrInternal)
	}
	response["inventory"] = inv
	history, err := s.getCoinHistory(id) // История переводов
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("cant get user history %w", models.ErrInternal)
	}
	response["coinHistory"] = history
	return response, nil
}

// Функция чтобы узнать баланс по id сотрудника
func (s *ShopRepository) getCoins(id int) (int, error) {
	row := s.db.QueryRow("SELECT coins FROM Employees WHERE id=$1;", id)
	var coins int
	if err := row.Scan(&coins); err != nil {
		return 0, err
	}
	return coins, nil
}

// Вспомогательная структура
type Item struct {
	Name string
	Num  int
}

// Функция для запросов, которые возвращают массив пар (строка и число)
func (s *ShopRepository) getItem(request string) ([]Item, error) {
	rows, err := s.db.Query(request)
	if err != nil {
		return nil, err
	}
	arr := make([]Item, 0)
	var i Item
	for rows.Next() {
		if err := rows.Scan(&i.Name, &i.Num); err != nil {
			return nil, err
		}
		arr = append(arr, i)
	}
	return arr, nil
}

// Вспомогательная структура для представления истории операций
type CoinHistory struct {
	Received []Item // Полученные переводы
	Send     []Item // Отправленные переводы
}

// Функция для получения истории операций переводов
func (s *ShopRepository) getCoinHistory(id int) (CoinHistory, error) {
	send, err := s.getItem(fmt.Sprintf("SELECT e.name, t.coins FROM Transaction AS t "+
		"INNER JOIN Employees AS e ON t.id_dest = e.id WHERE t.id_src = %d;", id))
	if err != nil {
		return CoinHistory{}, err
	}
	received, err := s.getItem(fmt.Sprintf("SELECT e.name, t.coins FROM Transaction AS t "+
		"INNER JOIN Employees AS e ON t.id_src = e.id WHERE t.id_dest = %d;", id))
	if err != nil {
		return CoinHistory{}, err
	}
	return CoinHistory{Received: received, Send: send}, nil
}
