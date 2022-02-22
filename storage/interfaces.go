package storage

import (
	"context"
	"github.com/IlyaYP/diploma/model"
	"io"
)

// UserStorage defines model.User create/update operations.
type UserStorage interface {
	io.Closer

	// CreateUser creates a new model.User.
	// Returns ErrAlreadyExists if user exists.
	CreateUser(ctx context.Context, user model.User) (model.User, error)

	// GetUserByLogin returns model.User by its login if exists.
	GetUserByLogin(ctx context.Context, login string) (*model.User, error)
}
