package server

import (
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

// страница с карточками
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
	case http.MethodPost:
		crt := &cart.Cart{
			UserID:      userID,
			ProductName: r.FormValue("productName"),
		}
		product, _ := strconv.Atoi(r.FormValue("price"))
		crt.Price = uint(product)
		crt.AppendCartToDB()
		helper.LoadPage(w, "carts", cart.ShowCartsOfUser(userID))
	default:
		helper.LoadPage(w, "carts", cart.ShowCartsOfUser(userID))
	}
}

// страница со складами
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
		stck := &stock.Stock{
			UserId:    userID,
			StockName: r.FormValue("stockName"),
		}
		stck.CreateNewStock()
		helper.LoadPage(w, "stocks", stock.ShowStocksOfUser(userID))
	default:
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
	case http.MethodGet:
		helper.LoadPage(w, "stockinfo", stock.ShowInfoAbountStock(uint(stockID), userID))
	case http.MethodPost:
		productID, _ := strconv.Atoi(r.FormValue("productID"))
		amount, _ := strconv.Atoi(r.FormValue("amount"))

		crt := cart.Cart{
			ProductID:   uint(productID),
			UserID:      userID,
			ProductName: "",
			Price:       0,
		}

		crt.AppendCartToStock(uint(stockID), uint(amount))
		helper.LoadPage(w, "stockinfo", stock.ShowInfoAbountStock(uint(stockID), userID))
	default:
		http.Error(w, "invalid http method", http.StatusMethodNotAllowed)
	}
}
