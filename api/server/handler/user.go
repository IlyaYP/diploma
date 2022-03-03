package handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

/*
POST /api/user/register — регистрация пользователя;
POST /api/user/login — аутентификация пользователя;
POST /api/user/orders — загрузка пользователем номера заказа для расчёта;
GET /api/user/orders — получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях;
GET /api/user/balance — получение текущего баланса счёта баллов лояльности пользователя;
POST /api/user/balance/withdraw — запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа;
GET /api/user/balance/withdrawals — получение информации о выводе средств с накопительного счёта пользователем.

*/

func (h *Handler) user(router chi.Router) {
	router.Post("/register", h.UserRegister)
	router.Post("/login", h.UserLogin)
	router.Post("/orders", h.PutOrder)
	router.Get("/orders", h.GetOrders)
	router.Get("/balance", h.GetBalance)
	router.Post("/balance/withdraw", h.Withdraw)
	router.Get("/balance/withdrawals", h.GetWithdrawals)

}

// UserRegister register new user
func (h *Handler) UserRegister(w http.ResponseWriter, r *http.Request) {
	res := struct {
		Name string `json:"name"`
	}{"UserRegister"}
	resJson, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(resJson)
}

// UserLogin authenticates user
func (h *Handler) UserLogin(w http.ResponseWriter, r *http.Request) {
	res := struct {
		Name string `json:"name"`
	}{"UserLogin"}
	resJson, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(resJson)
}

// PutOrder Puts new order
func (h *Handler) PutOrder(w http.ResponseWriter, r *http.Request) {
	res := struct {
		Name string `json:"name"`
	}{"PutOrder"}
	resJson, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(resJson)
}

// GetOrders Gets order list
func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	res := struct {
		Name string `json:"name"`
	}{"GetOrders"}
	resJson, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(resJson)
}

// GetBalance Gets balance
func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	res := struct {
		Name string `json:"name"`
	}{"GetBalance"}
	resJson, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(resJson)
}

// Withdraw Request withdraw
func (h *Handler) Withdraw(w http.ResponseWriter, r *http.Request) {
	res := struct {
		Name string `json:"name"`
	}{"Withdraw"}
	resJson, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(resJson)
}

// GetWithdrawals history
func (h *Handler) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	res := struct {
		Name string `json:"name"`
	}{"GetWithdrawals"}
	resJson, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(resJson)
}
