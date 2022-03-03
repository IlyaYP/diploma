package handler

import (
	"encoding/json"
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

//func (h *Handler) user(router chi.Router) {
//	router.Post("/api/user/register", router.UserRegister)
//	router.Post("/api/user/login", router.UserLogin)
//	router.Post("/api/user/orders", router.PutOrder)
//	router.Get("/api/user/orders", router.GetOrders)
//	router.Get("/api/user/balance", router.GetBalance)
//	router.Post("/api/user/balance/withdraw", router.Withdraw)
//	router.Get("/api/user/balance/withdrawals", router.GetWithdrawals)
//
//}

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
