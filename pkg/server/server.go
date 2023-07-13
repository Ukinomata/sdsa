package server

import (
	"github.com/gorilla/sessions"
	"net/http"
	"strconv"
	"strings"
	"warehouse-application/internal/cart"
	"warehouse-application/internal/handlers"
	"warehouse-application/internal/stock"
	"warehouse-application/internal/user"
	"warehouse-application/pkg/helper"
	"warehouse-application/pkg/logging"
)

var store *sessions.CookieStore

func StartServer(logger logging.Logger) {
	store = sessions.NewCookieStore([]byte("ukinoshito-ukino"))

	newHandler := NewHandler(logger)
	newHandler.Register()

	logger.Fatal(http.ListenAndServe("localhost:8080", nil))
}

type handler struct {
	logger logging.Logger
}

func NewHandler(logger logging.Logger) handlers.Handler {
	return &handler{logger: logger}
}

func (h *handler) Register() {
	http.HandleFunc("/signup", h.signUpHandler)
	http.HandleFunc("/login", h.loginHandler)
	http.HandleFunc("/profile", h.profileHandler)
	http.HandleFunc("/logout", h.logountHandler)
	http.HandleFunc("/stocks", h.stocksHandler)
	http.HandleFunc("/carts", h.cartsHandler)
	http.HandleFunc("/stockinfo", h.stockInfoHandler)
}

// для регистрации пользователя
func (h *handler) signUpHandler(w http.ResponseWriter, r *http.Request) {
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
func (h *handler) logountHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		h.logger.Info(err)
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
func (h *handler) loginHandler(w http.ResponseWriter, r *http.Request) {
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
			h.logger.Info(err)
			return
		}

		session, err := store.Get(r, "session-name")

		if err != nil {
			h.logger.Info(err)
			return
		}

		session.Values["userID"] = usr.Id
		err = session.Save(r, w)

		if err != nil {
			h.logger.Info(err)
			return
		}

		http.Redirect(w, r, "/profile", http.StatusFound)
		return
	default:
		helper.LoadPage(w, "login", nil)
	}
}

// главная страница
func (h *handler) profileHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")

	if err != nil {
		h.logger.Info(err)
		return
	}

	userID, ok := session.Values["userID"].(uint)
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
func (h *handler) cartsHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		h.logger.Info(err)
		return
	}

	userID, ok := session.Values["userID"].(uint)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
	}

	switch r.Method {
	case http.MethodPatch:
		crt := &cart.Cart{
			UserID: userID,
		}

		helper.Unmarshal(r, crt)

		crt.CorrectPrice(h.logger)

	case http.MethodDelete:
		crt := &cart.Cart{
			UserID: userID,
		}

		helper.Unmarshal(r, crt)

		crt.DeleteCart(h.logger)

	case http.MethodPost:
		crt := &cart.Cart{
			UserID: userID,
		}

		helper.Unmarshal(r, crt)

		crt.AppendCartToDB(h.logger)
		return
	default:
		helper.LoadPage(w, "carts", cart.ShowCartsOfUser(userID, h.logger))
	}
}

// страница со складами +
func (h *handler) stocksHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		h.logger.Info(err)
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
			UserId: userID,
		}

		helper.Unmarshal(r, stck)
		stck.CreateNewStock(h.logger)
		return
	default:
		helper.LoadPage(w, "stocks", stock.ShowStocksOfUser(userID, h.logger))
	}

}

// добавление и изменение карточек на складах
func (h *handler) stockInfoHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		h.logger.Info(err)
		return
	}

	userID, ok := session.Values["userID"].(uint)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}

	st := strings.ReplaceAll(r.URL.String(), "/stockinfo?button=", "")
	stockID, _ := strconv.Atoi(st)

	switch r.Method {
	case http.MethodDelete:
		stk := &stock.StockInfo{
			StockID: uint(stockID),
		}

		helper.Unmarshal(r, stk)

		stk.DeleteFromStock(h.logger)
		return
	case http.MethodPatch:
		stk := &stock.StockInfo{
			StockID: uint(stockID),
			Amount:  0,
		}

		helper.Unmarshal(r, stk)
		stk.CorrectAmount(h.logger)
		return
	case http.MethodPost:
		stk := &stock.StockInfo{
			StockID: uint(stockID),
		}

		helper.Unmarshal(r, stk)
		stk.AppendCartToStock(userID, h.logger)
		return
	default:
		helper.LoadPage(w, "stockinfo", stock.ShowInfoAbountStock(uint(stockID), userID, h.logger))
	}
}
