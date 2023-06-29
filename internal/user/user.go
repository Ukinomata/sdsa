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

//регистрация пользователя

func (u *User) RegisterUser(db *sql.DB) error {
	data := `INSERT INTO users(login,hash_password) VALUES ($1,$2)`

	if _, err := db.Exec(data, u.Login, helper.Hashing(u.Password)); err != nil {
		log.Println("enter another username")
		return err
	}

	return nil
}

//проверка то что пользователь правильно ввел логин и пароль

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

//Получение user id

func (u *User) GetUserId(db *sql.DB) (int, error) {
	var id int
	err := db.QueryRow(`SELECT id FROM users WHERE login = $1 and hash_password = $2`, u.Login, helper.Hashing(u.Password)).Scan(&id)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return id, nil
}

//создание нового склада

func (u *User) CreateNewStock(db *sql.DB, str string) error {
	id, _ := u.GetUserId(db)
	if _, err := db.Exec(`INSERT INTO stocks_of_user(user_id, stock_name) VALUES ($1,$2)`, id, str); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//cписок складов

func (u *User) ListOfStock(db *sql.DB) ([]int, error) {
	id, _ := u.GetUserId(db)

	rows, err := db.Query(`SELECT stock_id FROM stocks_of_user WHERE user_id = $1`, id)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	list := []int{}
	for rows.Next() {
		var value int
		err = rows.Scan(&value)
		if err != nil {
			log.Fatal(err)
		}
		list = append(list, value)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return list, nil
}

//cписок id карточек товаров которые есть у пользователя

func (u *User) ListOfCarts(db *sql.DB) ([]int, error) {
	id, _ := u.GetUserId(db)

	rows, err := db.Query(`SELECT id FROM carts WHERE user_id = $1`, id)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	list := []int{}
	for rows.Next() {
		var value int
		err = rows.Scan(&value)
		if err != nil {
			log.Fatal(err)
		}
		list = append(list, value)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return list, nil
}
