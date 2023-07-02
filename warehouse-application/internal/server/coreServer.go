package server

import (
	"database/sql"
	"net/http"
	"warehouse-application/pkg/errors"
	"warehouse-application/pkg/page"
)

var DB *sql.DB

func StartWebServer(db *sql.DB) {
	DB = db

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		page.ShowPage(w, r, "hub")
	})

	http.HandleFunc("/registration", func(w http.ResponseWriter, r *http.Request) {
		page.ShowPage(w, r, "registration")
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		page.ShowPage(w, r, "login")
	})

	http.HandleFunc("/RegistrationUser", HandlerRegistrationUser)

	http.HandleFunc("/LoginUser", HandlerLoginUser)

	http.HandleFunc("/mainMenu", func(w http.ResponseWriter, r *http.Request) {
		page.ShowPage(w, r, "mainMenu")
	})

	//http.HandleFunc("/CreateStock", handlerCreateStock)

	err := http.ListenAndServe("localhost:8080", nil)
	errors.CheckError(err)
}
