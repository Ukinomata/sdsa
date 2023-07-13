package stock

import (
	"warehouse-application/internal/cart"
	"warehouse-application/pkg/database"
	"warehouse-application/pkg/logging"
)

//структура для создания нового склада

type Stock struct {
	UserId    uint
	StockName string `json:"stockName"`
}

//функция которая создает новый склад

func (s *Stock) CreateNewStock(logger logging.Logger) {
	db := database.ConnectToDB()
	defer db.Close()

	data := `INSERT INTO stock_of_user(user_id, stock_name) VALUES ($1,$2)`

	if _, err := db.Exec(data, s.UserId, s.StockName); err != nil {
		logger.Info(err)
		return
	}
	logger.Info("New stock is created")
}

//структура для удобного вывода складов которые есть у пользователя

type superStock struct {
	MapVariable map[uint]string
}

//склады которые есть у пользователя

func ShowStocksOfUser(userID uint, logger logging.Logger) superStock {
	db := database.ConnectToDB()
	defer db.Close()

	data := `SELECT id,stock_name FROM stock_of_user WHERE user_id = $1`

	query, err := db.Query(data, userID)
	if err != nil {
		logger.Info(err)
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
			logger.Info(err)
			return superStock{nil}
		}

		m[stockId] = stockName

	}
	if query.Err() != nil {
		logger.Info(err)
		return superStock{nil}
	}

	return superStock{m}
}

//информация о товаре внутри склада и его количестве

type StockInfo struct {
	StockID     uint
	ProductID   uint   `json:"productID"`
	ProductName string `json:"productName"`
	Price       uint   `json:"price"`
	Amount      uint   `json:"amount"`
}

//заполнить структуру stockInfo при наличии product_id

func (sI *StockInfo) FillStock(logger logging.Logger) {
	db := database.ConnectToDB()
	defer db.Close()
	data := `SELECT name,price FROM product WHERE id = $1`

	err := db.QueryRow(data, sI.ProductID).Scan(&sI.ProductName, &sI.Price)
	if err != nil {
		logger.Info(err)
		return
	}
}

type AllInfo struct {
	Stock []StockInfo
	Carts cart.SuperCart
}

func ShowInfoAbountStock(stockID, userID uint, logger logging.Logger) AllInfo {
	db := database.ConnectToDB()
	defer db.Close()

	data := `SELECT product_id,amount FROM stock WHERE id = $1`
	query, err := db.Query(data, stockID)
	if err != nil {
		logger.Info(err)
		return AllInfo{}
	}

	defer query.Close()

	var result AllInfo

	for query.Next() {
		var stk StockInfo

		err = query.Scan(&stk.ProductID, &stk.Amount)
		if err != nil {
			logger.Info(err)
			return AllInfo{}
		}
		stk.FillStock(logger)
		stk.StockID = stockID

		result.Stock = append(result.Stock, stk)
	}
	if query.Err() != nil {
		logger.Info(err)
		return AllInfo{}
	}
	result.Carts = cart.ShowCartsOfUser(userID, logger)
	return result
}

func (s *StockInfo) CorrectAmount(logger logging.Logger) {
	db := database.ConnectToDB()
	defer db.Close()

	data := `UPDATE stock SET amount = $1 WHERE id = $2 and product_id = $3`

	if _, err := db.Exec(data, s.Amount, s.StockID, s.ProductID); err != nil {
		logger.Info(err)
		return
	}
	logger.Info("amount is corrected")
}

func (s *StockInfo) AppendCartToStock(userID uint, logger logging.Logger) {
	db := database.ConnectToDB()
	defer db.Close()

	data := `INSERT INTO stock(id, user_id,product_id, amount)
SELECT $1,$2,$3,$4
WHERE NOT EXISTS(
    SELECT * FROM stock
    WHERE id = $1 and user_id = $2 and product_id = $3
) `

	if _, err := db.Exec(data, s.StockID, userID, s.ProductID, s.Amount); err != nil {
		logger.Info(err)
		return
	}
	logger.Info("cart has been added to stock")
}

func (s *StockInfo) DeleteFromStock(logger logging.Logger) {
	db := database.ConnectToDB()
	defer db.Close()

	data := `DELETE FROM stock WHERE id = $1 and product_id = $2`

	if _, err := db.Exec(data, s.StockID, s.ProductID); err != nil {
		logger.Info(err)
		return
	}
	logger.Info("cart has been deleted from stock")
}
