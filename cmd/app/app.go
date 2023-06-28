package main

import (
	"warehouse-application/internal/cart"
	"warehouse-application/internal/stock"
	"warehouse-application/internal/user"
	"warehouse-application/pkg/database"
)

func main() {
	db := database.ConnectToDB()

	j := &user.User{
		Login:    "ukino",
		Password: "951753554455Al",
	}
	j.LoginToDatabase(db)

	c := &cart.Cart{
		Name:  "Converse All Stars",
		Price: 31000,
	}

	c.CreateNewCart(j, db)

	k := &stock.Stock{Id: 1, CartId: 2, Amount: 35}

	k.AddItems(db)

}
