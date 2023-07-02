package user

import (
	"database/sql"
	"errors"
	"log"
)

type User struct {
	Login    string
	Password string
}

func (alreadyUser *User) RegisterUser(db *sql.DB) error {
	_, err := alreadyUser.getHashPassword(db)
	if err == nil {
		return errors.New("user already registered")
	}
	login := alreadyUser.Login
	password := alreadyUser.Password
	hashPassword := Hashing(password)
	query := `INSERT INTO users(login, hash_password) VALUES($1, $2)`
	_, err = db.Exec(query, login, hashPassword)
	return err
}

func (alreadyUser *User) getHashPassword(db *sql.DB) (hashPassword string, err error) {
	query := `SELECT hash_password FROM users WHERE login = $1`
	login := alreadyUser.Login
	err = db.QueryRow(query, login).Scan(&hashPassword)
	return
}

func (alreadyUser *User) LoginUser(db *sql.DB) (keySession string, err error) {
	validHashPassword, err := alreadyUser.getHashPassword(db)
	if err != nil {
		err = errors.New("user dont registered")
		return
	}
	hashPassword := Hashing(alreadyUser.Password)
	if validHashPassword != hashPassword {
		err = errors.New("wrong password or login")
		return
	}
	keySession = alreadyUser.StartSession(db)
	return
}

func (alreadyUser *User) ExitUser(db *sql.DB) (err error) {
	login := alreadyUser.Login
	_, err = getSessionKey(db, login)
	if err != nil {
		return
	}
	alreadyUser.CloseSession(db)
	return
}

func (alreadyUser *User) removeUser(db *sql.DB) (err error) {
	alreadyUser.ExitUser(db)
	login := alreadyUser.Login
	query := `DELETE FROM users WHERE login = $1`
	_, err = db.Exec(query, login)
	return
}

//Получение user id

func (alreadyUser *User) GetUserId(db *sql.DB) (int, error) {
	var id int
	err := db.QueryRow(`SELECT id FROM users WHERE login = $1`, alreadyUser.Login).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

//создание нового склада

func (alreadyUser *User) CreateNewStock(db *sql.DB, str string) error {
	id, _ := alreadyUser.GetUserId(db)
	if _, err := db.Exec(`INSERT INTO stocks_of_user(user_id, stock_name) VALUES ($1,$2)`, id, str); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//cписок складов

func (alreadyUser *User) ListOfStock(db *sql.DB) ([]int, error) {
	id, _ := alreadyUser.GetUserId(db)

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

func (alreadyUser *User) ListOfCarts(db *sql.DB) ([]int, error) {
	id, _ := alreadyUser.GetUserId(db)

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
