package user

import (
	"warehouse-application/pkg/database"
	"warehouse-application/pkg/helper"
	"warehouse-application/pkg/logging"
)

type User struct {
	Id       uint
	Username string
	Password string
}

func (u *User) SignUpUser(logger logging.Logger) {
	db := database.ConnectToDB()
	defer db.Close()

	data := `INSERT INTO users(username, password) VALUES ($1,$2)`

	if _, err := db.Exec(data, u.Username, helper.Hashing(u.Password)); err != nil {
		logger.Info(err)
		return
	}
	logger.Info("new user is register")
}

func (u *User) CorrectData(logger logging.Logger) error {
	db := database.ConnectToDB()
	defer db.Close()

	data := `SELECT id FROM users WHERE username = $1 and password = $2`

	if err := db.QueryRow(data, u.Username, helper.Hashing(u.Password)).Scan(&u.Id); err != nil {
		logger.Info(err)
		return err
	}
	return nil
}

func (u *User) GetInfo(logger logging.Logger) {
	db := database.ConnectToDB()
	defer db.Close()

	data := `SELECT username FROM users WHERE id = $1`

	if err := db.QueryRow(data, u.Id).Scan(&u.Username); err != nil {
		logger.Info(err)
		return
	}
}
