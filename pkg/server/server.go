package server

import (
	"fmt"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"strconv"
	"warehouse-application/internal/cart"
	"warehouse-application/internal/stock"
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
	http.HandleFunc("/stocks", stocksHandler)
	http.HandleFunc("/carts", cartsHandler)
	http.HandleFunc("/stockinfo", stockInfoHandler)

	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

// для регистрации пользователя
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

// выйти из профиля
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

// войти в профиль
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

// главная страница
func profileHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")

	if err != nil {
		log.Println(err)
		return
	}

	userID, ok := session.Values["userID"].(uint)
	//fmt.Println(ok)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}
	usr := &user.User{
		Id: userID,
	}
	usr.GetInfo()
	helper.LoadPage(w, "profile", usr)
}

// страница с карточками +
func cartsHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		log.Println(err)
		return
	}

	userID, ok := session.Values["userID"].(uint)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
	}

	switch r.Method {
	case http.MethodPatch:
		fmt.Println("IT IS METHOD PATCH")
		crt := &cart.Cart{
			UserID: userID,
		}

		helper.Unmarshal(r, crt)

		crt.CorrectPrice()

		fmt.Println(crt)
	case http.MethodDelete:
		fmt.Println("IT IS METHOD DELETE")
		crt := &cart.Cart{
			UserID: userID,
		}

		helper.Unmarshal(r, crt)

		crt.DeleteCart()

		fmt.Println(crt)
	case http.MethodPost:
		fmt.Println("IT IS POST")
		crt := &cart.Cart{
			UserID: userID,
		}

		helper.Unmarshal(r, crt)

		fmt.Println(crt)
		crt.AppendCartToDB()
		return
	default:
		fmt.Println("IT IS DEFAULT")
		helper.LoadPage(w, "carts", cart.ShowCartsOfUser(userID))
	}
}

// страница со складами +
func stocksHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		log.Println(err)
		return
	}

	userID, ok := session.Values["userID"].(uint)
	_ = userID
	if !ok {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodPost:
		fmt.Println("IT IS POST")
		stck := &stock.Stock{
			UserId: userID,
		}

		helper.Unmarshal(r, stck)
		fmt.Println(stck)
		stck.CreateNewStock()
		return
	default:
		fmt.Println("IT IS DEFAULT")
		helper.LoadPage(w, "stocks", stock.ShowStocksOfUser(userID))
	}

}

// добавление и изменение карточек на складах
func stockInfoHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		log.Println(err)
		return
	}

	userID, ok := session.Values["userID"].(uint)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}

	st := r.URL.String()
	stockID, _ := strconv.Atoi(string(st[len(st)-1]))

	switch r.Method {
	case http.MethodDelete:
		fmt.Println("IT IS DELETE METHOD")
		stk := &stock.StockInfo{
			StockID: uint(stockID),
		}

		helper.Unmarshal(r, stk)

		fmt.Println(stk)
		stk.DeleteFromStock()

	case http.MethodPatch:
		fmt.Println("IT IS PATCH")
		stk := &stock.StockInfo{
			StockID: uint(stockID),
			Amount:  0,
		}

		helper.Unmarshal(r, stk)
		fmt.Println(stk)
		stk.CorrectAmount()
		return
	case http.MethodPost:
		fmt.Println("IT IS POST METHOD")
		stk := &stock.StockInfo{
			StockID: uint(stockID),
		}

		helper.Unmarshal(r, stk)
		fmt.Println(stk)
		stk.AppendCartToStock(userID)
		return
	default:
		fmt.Println("IT IS DEFAULT METHOD")
		helper.LoadPage(w, "stockinfo", stock.ShowInfoAbountStock(uint(stockID), userID))
	}
}

//todo добавить новый логгер, создать отдельную директорию под ошибки,  и удалять карточки

//todo после добавления товара и сразу изменения его количества добавляется новый товар.Думаю нужно переписать передачу данных на java script вместо html

//todo реализовать получение данных с бд с использованием java script
