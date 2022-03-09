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
func (h *Handler) Withdraw(w http.ResponseWriter, r *http.Request) {
	res := struct {
		Name string `json:"name"`
	}{"Withdrawal"}
	resJson, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(resJson)
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
