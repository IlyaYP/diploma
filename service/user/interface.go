package user

import (
	"context"
	"github.com/IlyaYP/diploma/model"
)

type Service interface {
	// CreateUser creates a new user.
	Register(ctx context.Context, login, password string) (model.User, error)
	// Login Authenticates user
	Login(ctx context.Context, login, password string) (model.User, error)
}
