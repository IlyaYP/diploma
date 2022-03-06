package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/IlyaYP/diploma/model"
	"github.com/IlyaYP/diploma/pkg"
	"github.com/IlyaYP/diploma/pkg/logging"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"github.com/lestrrat-go/jwx/jwt"
	"net/http"
)

func (h *Handler) user(router chi.Router) {
	router.Post("/register", h.UserRegister)
	router.Post("/login", h.UserLogin)
}

// UserRegister register new user
func (h *Handler) UserRegister(w http.ResponseWriter, r *http.Request) {
	ctx, _ := logging.GetCtxLogger(r.Context()) // correlationID is created here
	logger := h.Logger(ctx)

	input := &model.User{}
	if err := render.Bind(r, input); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		logger.Err(err).Msg("Error register user")
		return
	}

	logger.UpdateContext(input.GetLoggerContext)

	user, err := h.userSvc.Register(ctx, input.Login, input.Password)
	if err != nil {
		logger.Err(err).Msg("Error register user")
		if errors.Is(err, pkg.ErrAlreadyExists) {
			render.Render(w, r, ErrAlreadyExists)
			return
		}
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	_, tokenString, _ := h.tokenAuth.Encode(map[string]interface{}{"login": user.Login})
	w.Header().Set("Authorization", "Bearer "+tokenString)

	logger.Info().Msg("Successfully registered user")
	//render.Render(w, r, &user) // TODO: Remove

}

// UserLogin authenticates user
func (h *Handler) UserLogin(w http.ResponseWriter, r *http.Request) {
	ctx, _ := logging.GetCtxLogger(r.Context()) // correlationID is created here
	logger := h.Logger(ctx)

	input := &model.User{}
	if err := render.Bind(r, input); err != nil {
		logger.Err(err).Msg("Error login user")
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	logger.UpdateContext(input.GetLoggerContext)

	user, err := h.userSvc.Login(ctx, input.Login, input.Password)
	if err != nil {
		logger.Err(err).Msg("Login Unsuccessful")
		if errors.Is(err, pkg.ErrInvalidLogin) {
			render.Render(w, r, ErrInvalidLogin)
			return
		}
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	_, tokenString, _ := h.tokenAuth.Encode(map[string]interface{}{"login": user.Login})
	w.Header().Set("Authorization", "Bearer "+tokenString)

	logger.Info().Msg("Login Success")
}

func (h *Handler) UserContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, _ := logging.GetCtxLogger(r.Context()) // correlationID is created here
		logger := h.Logger(ctx)

		token, claims, err := jwtauth.FromContext(ctx)

		if err != nil {
			http.Error(w, err.Error(), 401)
			return
		}

		if token == nil || jwt.Validate(token) != nil {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		login := fmt.Sprintf("%v", claims["login"])
		input := &model.User{Login: login}
		logger.UpdateContext(input.GetLoggerContext)
		ctx = logging.SetCtxLogger(ctx, *logger)

		user, err := h.userSvc.GetUserByLogin(ctx, login)
		if err != nil {
			http.Error(w, err.Error(), 401)
			logger.Err(err).Msg("GetUserByLogin")
			return
		}

		ctx = context.WithValue(ctx, "user", user)
		logger.Info().Msg("UserContext")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
