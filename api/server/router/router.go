package router

import (
	"encoding/json"
	"fmt"
	"github.com/IlyaYP/diploma/service/user"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type (
	Handler struct {
		*chi.Mux
		user user.Service
	}
	Option func(h *Handler) error
)

func WithUserService(user user.Service) Option {
	return func(h *Handler) error {
		h.user = user
		return nil
	}
}

func NewHandler(opts ...Option) (*Handler, error) {
	h := &Handler{
		Mux: chi.NewMux(),
	}

	for _, opt := range opts {
		if err := opt(h); err != nil {
			return nil, fmt.Errorf("initialising dependencies: %w", err)
		}
	}

	h.Post("/api/user/register", h.UserRegister())
	h.Post("/api/user/login", h.UserLogin())
	h.Post("/api/user/orders", h.PutOrder())
	h.Get("/api/user/orders", h.GetOrders())
	h.Get("/api/user/balance", h.GetBalance())
	h.Post("/api/user/balance/withdraw", h.Withdraw())
	h.Get("/api/user/balance/withdrawals", h.GetWithdrawals())

	return h, nil
}

/*
POST /api/user/register — регистрация пользователя;
POST /api/user/login — аутентификация пользователя;
POST /api/user/orders — загрузка пользователем номера заказа для расчёта;
GET /api/user/orders — получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях;
GET /api/user/balance — получение текущего баланса счёта баллов лояльности пользователя;
POST /api/user/balance/withdraw — запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа;
GET /api/user/balance/withdrawals — получение информации о выводе средств с накопительного счёта пользователем.

*/

// UserRegister register new user
func (h *Handler) UserRegister() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}

// UserLogin authenticates user
func (h *Handler) UserLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}

// PutOrder Puts new order
func (h *Handler) PutOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}

// GetOrders Gets order list
func (h *Handler) GetOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}

// GetBalance Gets balance
func (h *Handler) GetBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}

// Withdraw Request withdraw
func (h *Handler) Withdraw() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}

// GetWithdrawals history
func (h *Handler) GetWithdrawals() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}
