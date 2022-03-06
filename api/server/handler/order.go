package handler

import (
	"encoding/json"
	"github.com/IlyaYP/diploma/model"
	"github.com/IlyaYP/diploma/pkg/logging"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
)

func (h *Handler) order(router chi.Router) {
	router.Use(jwtauth.Verifier(h.tokenAuth))
	router.Use(h.UserContext) //instead of router.Use(jwtauth.Authenticator)
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

	user, _ := ctx.Value("user").(*model.User)

	//_, claims, _ := jwtauth.FromContext(r.Context())
	////w.Write([]byte(fmt.Sprintf("protected area. hi %v", claims["login"])))
	//login := fmt.Sprintf("%v", claims["login"])
	res := struct {
		Name string `json:"name"`
	}{user.Login}
	resJson, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(resJson)
}
