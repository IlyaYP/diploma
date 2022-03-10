package handler

import (
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
//200 — номер заказа уже был загружен этим пользователем;
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

	orderNum, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		render.Render(w, r, ErrServerError(err))
		logger.Err(err).Msgf("NewOrder: wrong number %s", string(b))
		return
	}

	orderNumSt := strconv.FormatUint(orderNum, 10)

	user, ok := model.UserFromContext(ctx)
	if !ok {
		logger.Err(pkg.ErrInvalidLogin).Msg("GetOrders: can't get user from context")
		render.Render(w, r, ErrInvalidLogin)
		return
	}

	input := model.Order{Number: orderNumSt,
		User: user.Login, Status: model.OrderStatusNew}

	logger.UpdateContext(input.GetLoggerContext)

	if !pkg.ValidLuhn(orderNum) {
		render.Render(w, r, ErrInvalidOrderNum)
		logger.Err(pkg.ErrInvalidOrderNum).Msgf("NewOrder: wrong number %v", orderNum)
		return
	}

	// check if order exists   TODO: maybe move to order service
	if order, err := h.orderSvc.GetOrder(ctx, orderNumSt); err == nil {
		//fmt.Println("\n\n", order, "\n\n")
		if order.User == user.Login { // 200 order already loaded by this user
			return
		}
		render.Render(w, r, ErrAlreadyExists) // 409 order already loaded by other user
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

	render.Render(w, r, NewOrderAccepted) // 202
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
		return
	}

	orders, err := h.orderSvc.GetOrdersByUser(ctx, user.Login)
	if err != nil {
		if errors.Is(err, pkg.ErrNoData) {
			render.Render(w, r, ErrNoData)
			return
		}
		logger.Err(err).Msg("GetOrders: can't get orders from DB")
		render.Render(w, r, ErrServerError(err))
		return
	}

	render.Render(w, r, orders)
}
