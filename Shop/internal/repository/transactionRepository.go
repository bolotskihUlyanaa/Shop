package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/bolotskihUlyanaa/Shop/internal/models"
)

// Транзакция для отправки монет сотруднику
func (s *ShopRepository) SendCoins(idSender int, data models.SendCoinRequest) error {
	tx, err := s.db.Begin() // Начало транзакции
	if err != nil {
		log.Println(err)
		return fmt.Errorf("cant begin transaction %w", models.ErrInternal)
	}
	defer endTransaction(tx, &err)                       // Завершение транзакции
	coinsSender, err := getUserBalanceByID(tx, idSender) //Баланс отправителя
	if err != nil {
		log.Println(err)
		return fmt.Errorf("cant get user coins by id %w", models.ErrInternal)
	}
	coinsSender -= data.Amount // Баланс после отправки монет
	if coinsSender < 0 {       // Проверка достаточно ли средств
		return fmt.Errorf("insufficient funds %w", models.ErrBadRequest)
	}
	// Получить id и баланс получателя по имени
	idReceiver, coinsReceiver, err := getInfoAbout(tx, fmt.Sprintf(
		"SELECT id, coins FROM Employees WHERE name='%s';", data.ToUser))
	if err != nil {
		log.Println(err)
		return fmt.Errorf("user recipient not found %w", models.ErrBadRequest)
	}
	if idSender == idReceiver {
		return fmt.Errorf("can't translate for yourself %w", models.ErrBadRequest)
	}
	coinsReceiver += data.Amount // Баланс получателя
	// Обновление балансов отправителя и получателя
	_, err = tx.Exec("UPDATE Employees SET coins=$1 WHERE id=$2;", coinsSender, idSender)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("cant update coins user %w", models.ErrInternal)
	}
	_, err = tx.Exec("UPDATE Employees SET coins=$1 WHERE id=$2;", coinsReceiver, idReceiver)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("cant update coins user %w", models.ErrInternal)
	}
	// Сохранение транзакции
	_, err = tx.Exec("INSERT INTO Transaction (id_src, id_dest, coins) VALUES ($1, $2, $3);",
		idSender, idReceiver, data.Amount)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("cant insert transaction %w", models.ErrInternal)
	}
	return nil
}

// Функция для покупки мерча item
func (s *ShopRepository) Buy(id int, item string) error {
	tx, err := s.db.Begin() // Начало транзакции
	if err != nil {
		log.Println(err)
		return fmt.Errorf("cant begin transaction %w", models.ErrInternal)
	}
	defer endTransaction(tx, &err)             // Завершение транзакции
	balance, err := getUserBalanceByID(tx, id) // Узнать баланс
	if err != nil {
		log.Println(err)
		return fmt.Errorf("cant get user coins by id %w", models.ErrInternal)
	}
	if balance == 0 { // Недостаточно средств
		return fmt.Errorf("insufficient funds %w", models.ErrBadRequest)
	}
	idMerch, price, err := getInfoAbout(tx, fmt.Sprintf( // Получение id и цены мерча
		"SELECT id, price FROM Merch WHERE name='%s';", item))
	if err != nil {
		log.Println(err)
		return fmt.Errorf("merch not found %w", models.ErrBadRequest)
	}
	balance -= price // Отаток монет
	if balance < 0 { // Недостаточно средств
		return fmt.Errorf("insufficient funds %w", models.ErrBadRequest)
	}
	// Сохранение тразакции
	_, err = tx.Exec("INSERT INTO Buy (id_employees, id_merch) VALUES ($1, $2);", id, idMerch)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("cant insert buy %w", models.ErrInternal)
	}
	// Обновление баланса сотрудника после покупки
	_, err = tx.Exec("UPDATE Employees SET coins=$1 WHERE id=$2;", balance, id)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("cant update coins user %w", models.ErrInternal)
	}
	return nil
}

// Функция для завершения транзакции
func endTransaction(tx *sql.Tx, err *error) {
	if *err != nil { // Если произошла ошибка, то откат транзакции
		if rbErr := tx.Rollback(); rbErr != nil {
			log.Printf("failed to rollback transaction: %v", rbErr)
		}
		return
	}
	// Успешное завершение транзакции
	if commitErr := tx.Commit(); commitErr != nil {
		log.Println(commitErr)
		*err = fmt.Errorf("failed to commit transaction: %w", models.ErrInternal)
	}
}

// Функция получения баланса сотрудника по его id
func getUserBalanceByID(tx *sql.Tx, id int) (int, error) {
	row := tx.QueryRow("SELECT coins FROM Employees WHERE id=$1;", id)
	var coins int
	if err := row.Scan(&coins); err != nil {
		return 0, err
	}
	return coins, nil
}

// Функция для запросов request, которые возвращают 2 числа
func getInfoAbout(tx *sql.Tx, request string) (int, int, error) {
	row := tx.QueryRow(request)
	var id, coins int
	if err := row.Scan(&id, &coins); err != nil {
		return 0, 0, err
	}
	return id, coins, nil
}
