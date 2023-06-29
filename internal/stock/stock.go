package stock

import (
	"database/sql"
	"log"
)

type Stock struct {
	Id     int
	СartID int
}

//TODO Доделать stock!
//добавление предметов на склад

func (s *Stock) AddItems(db *sql.DB, amount int) error {
	data := `INSERT INTO stock VALUES($1,$2,$3)`

	if _, err := db.Exec(data, s.Id, s.СartID, amount); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

//функция которая изменяет количество товара на складе

func (s *Stock) EditAmount(db *sql.DB, newAmount int) error {
	data := `UPDATE stock SET amount = $1 WHERE cart_id = $2 and id = $3`

	if _, err := db.Exec(data, newAmount, s.СartID, s.Id); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

//функция которая плюсует количество товара на складе

func (s *Stock) PlusAmmount(db *sql.DB, plusAmount int) error {
	var amount int
	data := `SELECT amount FROM stock WHERE id = $1 and cart_id = $2`
	err := db.QueryRow(data, s.Id, s.СartID).Scan(&amount)
	if err != nil {
		log.Println(err)
		return err
	}
	amount = amount + plusAmount
	data = `UPDATE stock SET amount = $1 WHERE cart_id = $2 and id = $3`
	if _, err = db.Exec(data, amount, s.СartID, s.Id); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

//удаляет определенную карточку со склада

func (s *Stock) DeleteCartFromStock(db *sql.DB) error {
	data := `DELETE from stock WHERE id = $1 and cart_id = $2`

	if _, err := db.Exec(data, s.Id, s.СartID); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

//Список карточек которые есть в определенном складе

func (s *Stock) ListOfCartsInStock(db *sql.DB) (map[int]int, error) {
	data := `SELECT cart_id,amount FROM stock WHERE id = $1`

	rows, err := db.Query(data, s.Id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer rows.Close()

	list := map[int]int{}

	for rows.Next() {
		var cartId, amount int
		err = rows.Scan(&cartId, &amount)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		list[cartId] = amount
	}

	return list, nil
}
