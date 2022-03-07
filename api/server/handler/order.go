package handler

import (
	"encoding/json"
	"errors"
	"github.com/IlyaYP/diploma/model"
	"github.com/IlyaYP/diploma/pkg"
	"github.com/IlyaYP/diploma/pkg/logging"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"io"
	"net/http"
	"strconv"
)

func (h *Handler) order(router chi.Router) {
	router.Use(jwtauth.Verifier(h.tokenAuth))
	router.Use(h.UserContext) // instead of jwtauth.Authenticator
	router.Post("/", h.NewOrder)
	router.Get("/", h.GetOrders)

}

// NewOrder Puts new order
//200 — номер заказа уже был загружен этим пользователем; // TODO: later if needed
//202 — новый номер заказа принят в обработку;
//400 — неверный формат запроса;
//401 — пользователь не аутентифицирован;
//409 — номер заказа уже был загружен другим пользователем;
//422 — неверный формат номера заказа;
//500 — внутренняя ошибка сервера.
func (h *Handler) NewOrder(w http.ResponseWriter, r *http.Request) {
	ctx, _ := logging.GetCtxLogger(r.Context()) // correlationID is created here
	logger := h.Logger(ctx)

	b, err := io.ReadAll(r.Body)
	if err != nil {
		render.Render(w, r, ErrServerError(err))
		return
	}

	ordernum, err := strconv.Atoi(string(b))
	if err != nil {
		render.Render(w, r, ErrServerError(err))
		logger.Err(err).Msgf("NewOrder: wrong number %s", string(b))
		return
	}

	user, ok := model.UserFromContext(ctx)
	if !ok {
		logger.Err(pkg.ErrInvalidLogin).Msg("GetOrders: can't get user from context")
		render.Render(w, r, ErrInvalidLogin)
		return
	}

	input := model.Order{Number: ordernum, User: user.Login, Status: model.OrderStatusNew}

	logger.UpdateContext(input.GetLoggerContext)

	if !pkg.ValidLuhn(ordernum) {
		render.Render(w, r, ErrInvalidOrderNum)
		logger.Err(pkg.ErrInvalidOrderNum).Msgf("NewOrder: wrong number %v", ordernum)
		return
	}

	order, err := h.orderSvc.CreateOrder(ctx, input)
	if err != nil {
		logger.Err(err).Msg("Error create order")
		if errors.Is(err, pkg.ErrAlreadyExists) {
			render.Render(w, r, ErrAlreadyExists)
			return
		}
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	logger.Info().Msgf("NewOrder:%v", order.Number)

	render.Render(w, r, NewOrderAccepted)
}

// GetOrders Gets order list
func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	ctx, _ := logging.GetCtxLogger(r.Context()) // correlationID is created here
	logger := h.Logger(ctx)
	logger.Info().Msg("GetOrders")

	//user, _ := ctx.Value("user").(*model.User) // Removed to model
	user, ok := model.UserFromContext(ctx)
	if !ok {
		logger.Err(pkg.ErrInvalidLogin).Msg("GetOrders: can't get user from context")
		render.Render(w, r, ErrInvalidLogin)
	}

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
