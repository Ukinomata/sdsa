package main

import (
	"warehouse-application/pkg/database"
)

func main() {
	db := database.ConnectToDB()

	_ = db

}
