package cart

type Cart struct {
	Id     int
	Name   string
	Price  int64
	Amount int64
}

type Warehouse struct {
	Id    int
	Carts []Cart
}
