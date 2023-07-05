package server

import (
	"fmt"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"warehouse-application/internal/user"
	"warehouse-application/pkg/helper"
)

var store *sessions.CookieStore

func StartServer() {
	store = sessions.NewCookieStore([]byte("ukinoshito-ukino"))
	_ = store
	http.HandleFunc("/signup", signUpHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/profile", profileHandler)
	http.HandleFunc("/logout", logountHandler)

	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func signUpHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		username := r.FormValue("username")
		password := r.FormValue("password")
		usr := &user.User{
			Username: username,
			Password: password,
		}
		usr.SignUpUser()
		http.Redirect(w, r, "/profile", http.StatusFound)
		return
	default:
		helper.LoadPage(w, "signup", nil)
	}
}

func logountHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		log.Println(err)
		return
	}

	session.Options.MaxAge = -1
	err = session.Save(r, w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusFound)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		username := r.FormValue("username")
		password := r.FormValue("password")
		usr := &user.User{
			Username: username,
			Password: password,
		}
		err := usr.CorrectData()

		if err != nil {
			log.Println(err)
			return
		}

		session, err := store.Get(r, "session-name")

		if err != nil {
			log.Println(err)
			return
		}

		session.Values["userID"] = usr.Id
		fmt.Println(session.Values["userID"])
		err = session.Save(r, w)

		if err != nil {
			log.Println(err)
			return
		}

		http.Redirect(w, r, "/profile", http.StatusFound)
		return
	default:
		helper.LoadPage(w, "login", nil)
	}
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")

	if err != nil {
		log.Println(err)
		return
	}

	userID, ok := session.Values["userID"].(uint)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	usr := &user.User{
		Id: userID,
	}
	usr.GetInfo()
	helper.LoadPage(w, "profile", usr)
}
