package stock

import (
	"database/sql"
	"log"
)

type Stock struct {
	Id     int
	cartID int
}

//TODO Доделать stock!
//добавление предметов на склад

func (s *Stock) AddItems(db *sql.DB, amount int) error {
	data := `INSERT INTO stock VALUES($1,$2,$3)`

	if _, err := db.Exec(data, s.Id, s.cartID, amount); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

//функция которая изменяет количество товара на складе

func (s *Stock) EditAmount(db *sql.DB, newAmount int) error {
	data := `UPDATE stock SET amount = $1 WHERE cart_id = $2 and id = $3`

	if _, err := db.Exec(data, newAmount, s.cartID, s.Id); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
