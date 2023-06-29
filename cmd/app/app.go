package main

import (
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
}
