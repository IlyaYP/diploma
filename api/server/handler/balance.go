package handler

import (
	"errors"
	"github.com/IlyaYP/diploma/model"
	"github.com/IlyaYP/diploma/pkg"
	"github.com/IlyaYP/diploma/pkg/logging"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"net/http"
)

func (h *Handler) balance(router chi.Router) {
	router.Use(jwtauth.Verifier(h.tokenAuth))
	router.Use(h.UserContext)
	router.Get("/", h.GetBalance)
	router.Post("/withdraw", h.Withdraw)
	router.Get("/withdrawals", h.GetWithdrawals)
}

// GetBalance Gets balance
func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	ctx, _ := logging.GetCtxLogger(r.Context()) // correlationID is created here
	logger := h.Logger(ctx)
	logger.Info().Msg("GetBalance")

	user, ok := model.UserFromContext(ctx)
	if !ok {
		logger.Err(pkg.ErrInvalidLogin).Msg("GetBalance: can't get user from context")
		render.Render(w, r, ErrInvalidLogin)
		return
	}

	balance, err := h.orderSvc.GetBalanceByUser(ctx, user.Login)
	if err != nil {
		//if errors.Is(err, pkg.ErrNoData) {
		//	render.Render(w, r, ErrNoData)
		//	return
		//}
		logger.Err(err).Msg("GetBalance: can't get balance from order service")
		render.Render(w, r, ErrServerError(err))
		return
	}

	render.Render(w, r, &balance)
}

// Withdraw Request withdraw
/*
POST /api/user/balance/withdraw HTTP/1.1
Content-Type: application/json

{
	"order": "2377225624",
    "sum": 751
}
Возможные коды ответа:

200 — успешная обработка запроса;
401 — пользователь не авторизован;
402 — на счету недостаточно средств;
422 — неверный номер заказа;
500 — внутренняя ошибка сервера.
*/
func (h *Handler) Withdraw(w http.ResponseWriter, r *http.Request) {
	ctx, _ := logging.GetCtxLogger(r.Context()) // correlationID is created here
	logger := h.Logger(ctx)
	logger.Info().Msg("Withdraw")

	user, ok := model.UserFromContext(ctx)
	if !ok {
		logger.Err(pkg.ErrInvalidLogin).Msg("Withdraw: can't get user from context")
		render.Render(w, r, ErrInvalidLogin)
		return
	}

	withdrawal := model.Withdrawal{}
	if err := render.Bind(r, &withdrawal); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		logger.Err(err).Msg("Error decode Withdrawal")
		return
	}
	withdrawal.User = user.Login

	if !pkg.ValidLuhn(withdrawal.Order) {
		render.Render(w, r, ErrInvalidOrderNum)
		logger.Err(pkg.ErrInvalidOrderNum).Msgf("NewOrder: wrong order number %v", withdrawal.Order)
		return
	}

	if err := h.orderSvc.NewWithdrawal(ctx, withdrawal); err != nil {
		if err == pkg.ErrInsufficientBalance {

			logger.Info().Msg("Withdraw: Insufficient Balance")
			render.Render(w, r, ErrInsufficientBalance)
			return
		}
		logger.Err(err).Msgf("Withdraw: order:%v sum:%v", withdrawal.Order, withdrawal.Sum)
		render.Render(w, r, ErrServerError(err))
		return
	}
}

// GetWithdrawals history
func (h *Handler) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	ctx, _ := logging.GetCtxLogger(r.Context()) // correlationID is created here
	logger := h.Logger(ctx)
	logger.Info().Msg("GetWithdrawals")

	user, ok := model.UserFromContext(ctx)
	if !ok {
		logger.Err(pkg.ErrInvalidLogin).Msg("GetWithdrawals: can't get user from context")
		render.Render(w, r, ErrInvalidLogin)
		return
	}

	withdrawals, err := h.orderSvc.GetWithdrawalsByUser(ctx, user.Login)
	if err != nil {
		if errors.Is(err, pkg.ErrNoData) {
			render.Render(w, r, ErrNoData)
			return
		}
		logger.Err(err).Msg("GetWithdrawals: can't get orders from DB")
		render.Render(w, r, ErrServerError(err))
		return
	}

	render.Render(w, r, withdrawals)
}
