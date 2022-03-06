package handler

import (
	"encoding/json"
	"github.com/IlyaYP/diploma/model"
	"github.com/IlyaYP/diploma/pkg/logging"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
	"io"
	"net/http"
)

func (h *Handler) order(router chi.Router) {
	router.Use(jwtauth.Verifier(h.tokenAuth))
	router.Use(h.UserContext) // instead of jwtauth.Authenticator
	router.Post("/", h.PutOrder)
	router.Get("/", h.GetOrders)

}

// PutOrder Puts new order
//200 — номер заказа уже был загружен этим пользователем;
//202 — новый номер заказа принят в обработку;
//400 — неверный формат запроса;
//401 — пользователь не аутентифицирован;
//409 — номер заказа уже был загружен другим пользователем;
//422 — неверный формат номера заказа;
//500 — внутренняя ошибка сервера.
func (h *Handler) PutOrder(w http.ResponseWriter, r *http.Request) {
	ctx, _ := logging.GetCtxLogger(r.Context()) // correlationID is created here
	logger := h.Logger(ctx)

	b, err := io.ReadAll(r.Body)
	if err != nil {
		render.Render(w, r, ErrServerError(err))
		return
	}

	// TODO: move to order model
	logger.UpdateContext(func(logCtx zerolog.Context) zerolog.Context {
		return logCtx.Str("ordernum", string(b))
	})
	//*logger = logger.With().Str("OrderNum", string(b)).Logger()
	logger.Info().Msg("PutOrder")

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
