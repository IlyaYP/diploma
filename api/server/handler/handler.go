package handler

import (
	"fmt"
	"github.com/IlyaYP/diploma/service/user"
	"github.com/go-chi/chi/v5"
)

type (
	Handler struct {
		*chi.Mux
		user user.Service
	}
	Option func(h *Handler) error
)

func WithUserService(user user.Service) Option {
	return func(h *Handler) error {
		h.user = user
		return nil
	}
}

func NewHandler(opts ...Option) (*Handler, error) {
	h := &Handler{
		Mux: chi.NewMux(),
	}

	for _, opt := range opts {
		if err := opt(h); err != nil {
			return nil, fmt.Errorf("initialising dependencies: %w", err)
		}
	}

	//h.Route("/items", user)

	h.Post("/api/user/register", h.UserRegister)
	h.Post("/api/user/login", h.UserLogin)
	h.Post("/api/user/orders", h.PutOrder)
	h.Get("/api/user/orders", h.GetOrders)
	h.Get("/api/user/balance", h.GetBalance)
	h.Post("/api/user/balance/withdraw", h.Withdraw)
	h.Get("/api/user/balance/withdrawals", h.GetWithdrawals)

	return h, nil
}
