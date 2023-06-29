package cart

import (
	"database/sql"
	"log"
	"warehouse-application/internal/user"
)

type Cart struct {
	Name  string
	Price int64
}

//создание новой карточки товара

func (c *Cart) CreateNewCart(user *user.User, db *sql.DB) error {
	id, _ := user.GetUserId(db)
	if _, err := db.Exec(`INSERT INTO carts(user_id,name,price) VALUES($1,$2,$3)`, id, c.Name, c.Price); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
