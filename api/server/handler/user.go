package handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func (h *Handler) user(router chi.Router) {
	router.Post("/register", h.UserRegister)
	router.Post("/login", h.UserLogin)
}

// UserRegister register new user
func (h *Handler) UserRegister(w http.ResponseWriter, r *http.Request) {

	if _, err := h.userSvc.Register(r.Context(), "vasya2", "God"); err != nil {
		log.Println(err)
	}

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
