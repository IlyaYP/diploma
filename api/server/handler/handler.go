package handler

import (
	"fmt"
	"github.com/IlyaYP/diploma/service/user"
	"github.com/go-chi/chi/v5"
)

type (
	Handler struct {
		*chi.Mux
		userSvc user.Service
	}
	Option func(h *Handler) error
)

func WithUserService(user user.Service) Option {
	return func(h *Handler) error {
		h.userSvc = user
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

	h.Route("/api/user", h.user)

	return h, nil
}
