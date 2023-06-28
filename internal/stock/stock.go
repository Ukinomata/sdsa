package stock

import (
	"database/sql"
	"log"
)

type Stock struct {
	Id     int
	CartId int
	Amount int
}

func (s *Stock) AddItems(db *sql.DB) error {
	data := `INSERT INTO stock VALUES($1,$2,$3)`

	if _, err := db.Exec(data, s.Id, s.CartId, s.Amount); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
