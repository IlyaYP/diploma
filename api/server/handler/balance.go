package handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (h *Handler) balance(router chi.Router) {
	//router.Use(UserContext)
	router.Get("/", h.GetBalance)
	router.Post("/withdraw", h.Withdraw)
	router.Get("/withdrawals", h.GetWithdrawals)
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
