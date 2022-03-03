package handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (h *Handler) order(router chi.Router) {
	//router.Use(UserContext)
	router.Post("/", h.PutOrder)
	router.Get("/", h.GetOrders)

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
