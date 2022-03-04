package handler

import (
	"context"
	"fmt"
	"github.com/IlyaYP/diploma/pkg/logging"
	"github.com/IlyaYP/diploma/service/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
	"net/http"
)

const (
	serviceName = "handler"
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

/*
POST /api/user/register — регистрация пользователя;
POST /api/user/login — аутентификация пользователя;
POST /api/user/orders — загрузка пользователем номера заказа для расчёта;
GET /api/user/orders — получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях;
GET /api/user/balance — получение текущего баланса счёта баллов лояльности пользователя;
POST /api/user/balance/withdraw — запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа;
GET /api/user/balance/withdrawals — получение информации о выводе средств с накопительного счёта пользователем.

*/

func NewHandler(opts ...Option) (*Handler, error) {
	h := &Handler{
		Mux: chi.NewMux(),
	}

	for _, opt := range opts {
		if err := opt(h); err != nil {
			return nil, fmt.Errorf("initialising dependencies: %w", err)
		}
	}

	h.Use(render.SetContentType(render.ContentTypeJSON))
	h.MethodNotAllowed(methodNotAllowedHandler)
	h.NotFound(notFoundHandler)
	h.Route("/api/user", h.user)
	h.Route("/api/user/orders", h.order)
	h.Route("/api/user/balance", h.balance)
	h.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	return h, nil
}

func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(405)
	render.Render(w, r, ErrMethodNotAllowed)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(400)
	render.Render(w, r, ErrNotFound)
}

func (h *Handler) Logger(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, serviceName).Logger()

	return &logger
}
