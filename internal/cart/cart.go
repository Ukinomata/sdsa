package cart

import (
	"log"
	"warehouse-application/pkg/database"
)

//todo переделать показ карточек для пользователя и сделать так чтобы функция передавала структуру где будет указан product id карточки

//структура карточки

type Cart struct {
	ProductID   uint
	UserID      uint
	ProductName string
	Price       uint
}

//функция добавления новой карточки

func (c *Cart) AppendCartToDB() {
	db := database.ConnectToDB()
	defer db.Close()

	data := `INSERT INTO product(user_id, name, price) VALUES($1,$2,$3)`

	if _, err := db.Exec(data, c.UserID, c.ProductName, c.Price); err != nil {
		log.Println(err)
		return
	}
}

//структура для удобного вывода всех карточек пользователя

type SuperCart struct {
	Carts []struct {
		ProductID   uint
		ProductName string
		Price       uint
	}
}

//функция для показа всех карточек которые есть у пользователя

func ShowCartsOfUser(userId uint) SuperCart {
	db := database.ConnectToDB()
	defer db.Close()

	var result SuperCart

	data := `SELECT id,name,price FROM product WHERE user_id = $1`

	query, err := db.Query(data, userId)

	if err != nil {
		log.Println(err)
		return SuperCart{nil}
	}

	defer query.Close()

	for query.Next() {

		cart := struct {
			ProductID   uint
			ProductName string
			Price       uint
		}{}

		err = query.Scan(&cart.ProductID, &cart.ProductName, &cart.Price)
		if err != nil {
			log.Println(err)
			return SuperCart{nil}
		}
		result.Carts = append(result.Carts, cart)
	}
	if query.Err() != nil {
		log.Println(err)
		return SuperCart{nil}
	}

	return result
}

func (c *Cart) AppendCartToStock(stockID, amount uint) {
	db := database.ConnectToDB()
	defer db.Close()

	data := `INSERT INTO stock(id, user_id,product_id, amount) VALUES($1,$2,$3,$4)`

	if _, err := db.Exec(data, stockID, c.UserID, c.ProductID, amount); err != nil {
		log.Println(err)
		return
	}
}
