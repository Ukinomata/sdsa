package user

import (
	"database/sql"
	"errors"
	"log"
	"warehouse-application/pkg/helper"
)

type User struct {
	Login    string
	Password string
}

func (u *User) RegisterUser(db *sql.DB) error {
	data := `INSERT INTO users(login,hash_password) VALUES ($1,$2)`

	if _, err := db.Exec(data, u.Login, helper.Hashing(u.Password)); err != nil {
		log.Println("enter another username")
		return err
	}

	return nil
}

func (u *User) LoginToDatabase(db *sql.DB) error {
	var (
		logd string
		pas  string
	)

	err := db.QueryRow(`SELECT login,hash_password FROM users`).Scan(&logd, &pas)
	if err != nil {
		log.Println(err)
		return err
	}

	if logd == u.Login && pas == helper.Hashing(u.Password) {
		log.Println("you entry in system")
		return nil
	}
	log.Println("Wrong password or login!")
	return errors.New("Wrong password or login!")
}
