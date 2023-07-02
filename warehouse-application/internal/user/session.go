package user

import (
	"database/sql"
	"math/rand"
	"time"
	"warehouse-application/pkg/errors"
)

func (alreadyUser *User) StartSession(db *sql.DB) string {
	login := alreadyUser.Login
	sessionKey, err := getSessionKey(db, login)
	if err == nil {
		return sessionKey
	}
	return newSession(db, login, 0)
}

func getSessionKey(db *sql.DB, login string) (sessionKey string, err error) {
	query := `SELECT key FROM sessions WHERE login = $1`
	err = db.QueryRow(query, login).Scan(&sessionKey)
	return sessionKey, err
}

func newSession(db *sql.DB, login string, nReload int) string {
	if nReload >= 3 {
		return ""
	}
	sessionKey := generateSessionKey()
	query := `INSERT INTO sessions VALUES($1, $2)`
	_, err := db.Exec(query, login, sessionKey)
	if err != nil {
		return newSession(db, login, nReload+1)
	}
	return sessionKey
}

func generateSessionKey() (sessionKey string) {
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	maxChars := len(chars)
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 32; i++ {
		sessionKey += string(chars[random.Intn(maxChars)])
	}
	return
}

func (alreadyUser *User) CloseSession(db *sql.DB) {
	login := alreadyUser.Login
	_, err := getSessionKey(db, login)
	if err == nil {
		query := `DELETE FROM sessions WHERE login = $1`
		_, err = db.Exec(query, login)
		errors.CheckWarning(err)
	}
}
