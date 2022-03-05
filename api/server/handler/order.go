package handler

import (
	"encoding/json"
	"fmt"
	"github.com/IlyaYP/diploma/pkg/logging"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
)

func (h *Handler) order(router chi.Router) {
	// use the Bearer Authentication middleware
	//router.Use(oauth.Authorize("GMartSuperSecret", nil))

	// Seek, verify and validate JWT tokens
	router.Use(jwtauth.Verifier(h.tokenAuth))

	// Handle valid / invalid tokens. In this example, we use
	// the provided authenticator middleware, but you can write your
	// own very easily, look at the Authenticator method in jwtauth.go
	// and tweak it, its not scary.
	router.Use(jwtauth.Authenticator)

	router.Use(h.UserContext)
	router.Post("/", h.PutOrder)
	router.Get("/", h.GetOrders)

}

// PutOrder Puts new order
func (h *Handler) PutOrder(w http.ResponseWriter, r *http.Request) {
	ctx, _ := logging.GetCtxLogger(r.Context()) // correlationID is created here
	logger := h.Logger(ctx)
	logger.Info().Msg("PutOrder")

	//b, err := io.ReadAll(r.Body)
	//if err != nil {
	//	//log.Fatal(err)
	//}

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
	ctx, _ := logging.GetCtxLogger(r.Context()) // correlationID is created here
	logger := h.Logger(ctx)
	logger.Info().Msg("GetOrders")

	_, claims, _ := jwtauth.FromContext(r.Context())
	//w.Write([]byte(fmt.Sprintf("protected area. hi %v", claims["login"])))
	login := fmt.Sprintf("%v", claims["login"])
	res := struct {
		Name string `json:"name"`
	}{login}
	resJson, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(resJson)
}
