package user

import (
	"log"
	"warehouse-application/pkg/database"
	"warehouse-application/pkg/helper"
)

type User struct {
	Id       uint
	Username string
	Password string
}

func (u *User) SignUpUser() {
	db := database.ConnectToDB()
	defer db.Close()

	data := `INSERT INTO users(username, password) VALUES ($1,$2)`

	if _, err := db.Exec(data, u.Username, helper.Hashing(u.Password)); err != nil {
		log.Println(err)
		return
	}
}

func (u *User) CorrectData() error {
	db := database.ConnectToDB()
	defer db.Close()

	data := `SELECT id FROM users WHERE username = $1 and password = $2`

	if err := db.QueryRow(data, u.Username, helper.Hashing(u.Password)).Scan(&u.Id); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (u *User) GetInfo() {
	db := database.ConnectToDB()
	defer db.Close()

	data := `SELECT username FROM users WHERE id = $1`

	if err := db.QueryRow(data, u.Id).Scan(&u.Username); err != nil {
		log.Println(err)
		return
	}
}
