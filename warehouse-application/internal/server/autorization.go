package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"warehouse-application/internal/user"
	"warehouse-application/pkg/errors"
)

func HandlerRegistrationUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" && r.Header.Get("Content-Type") == "application/json" {
		newUser := getUser(r)
		err := newUser.RegisterUser(DB)
		writeAnswer(w, fmt.Sprint(err.Error()))
	}
}

func getUser(r *http.Request) (newUser user.User) {
	err := json.NewDecoder(r.Body).Decode(&newUser)
	errors.CheckWarning(err)
	return
}

func writeAnswer(w http.ResponseWriter, answer string) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(answer))
	errors.CheckWarning(err)
}

func HandlerLoginUser(w http.ResponseWriter, r *http.Request) {
	alreadyUser := getUser(r)
	keySession, err := alreadyUser.LoginUser(DB)
	if err == nil {
		setCookie(w, keySession)
	} else {
		writeAnswer(w, err.Error())
	}
}

func setCookie(w http.ResponseWriter, key string) {
	w.Header().Add("session", key)
}
