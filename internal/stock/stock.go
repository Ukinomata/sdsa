package stock

import (
	"log"
	"warehouse-application/internal/cart"
	"warehouse-application/pkg/database"
)

//структура для создания нового склада

type Stock struct {
	UserId    uint
	StockName string
}

//функция которая создает новый склад

func (s *Stock) CreateNewStock() {
	db := database.ConnectToDB()
	defer db.Close()

	data := `INSERT INTO stock_of_user(user_id, stock_name) VALUES ($1,$2)`

	if _, err := db.Exec(data, s.UserId, s.StockName); err != nil {
		log.Println(err)
		return
	}
}

//структура для удобного вывода складов которые есть у пользователя

type superStock struct {
	MapVariable map[uint]string
}

//склады которые есть у пользователя

func ShowStocksOfUser(userID uint) superStock {
	db := database.ConnectToDB()
	defer db.Close()

	data := `SELECT id,stock_name FROM stock_of_user WHERE user_id = $1`

	query, err := db.Query(data, userID)
	if err != nil {
		log.Println(err)
		return superStock{nil}
	}

	defer query.Close()

	m := map[uint]string{}

	for query.Next() {

		var (
			stockId   uint
			stockName string
		)

		err = query.Scan(&stockId, &stockName)
		if err != nil {
			log.Println(err)
			return superStock{nil}
		}

		m[stockId] = stockName

	}
	if query.Err() != nil {
		log.Println(err)
		return superStock{nil}
	}

	return superStock{m}
}

//информация о товаре внутри склада и его количестве

type StockInfo struct {
	ProductID   uint
	ProductName string
	Price       uint
	Amount      uint
}

//заполнить структуру stockInfo при наличии product_id

func (sI *StockInfo) FillStock() {
	db := database.ConnectToDB()
	defer db.Close()
	data := `SELECT name,price FROM product WHERE id = $1`

	err := db.QueryRow(data, sI.ProductID).Scan(&sI.ProductName, &sI.Price)
	if err != nil {
		log.Println(err)
		return
	}
}

type AllInfo struct {
	Stock []StockInfo
	Carts cart.SuperCart
}

func ShowInfoAbountStock(stockID, userID uint) AllInfo {
	db := database.ConnectToDB()
	defer db.Close()

	data := `SELECT product_id,amount FROM stock WHERE id = $1`
	query, err := db.Query(data, stockID)
	if err != nil {
		log.Println(err)
		return AllInfo{}
	}

	defer query.Close()

	var result AllInfo

	for query.Next() {
		var stk StockInfo

		err = query.Scan(&stk.ProductID, &stk.Amount)
		if err != nil {
			return AllInfo{}
		}
		stk.FillStock()

		result.Stock = append(result.Stock, stk)
	}
	if query.Err() != nil {
		log.Println(err)
		return AllInfo{}
	}
	result.Carts = cart.ShowCartsOfUser(userID)
	return result
}
