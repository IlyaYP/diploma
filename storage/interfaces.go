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

// OrderStorage defines model.Order create/update operations.
type OrderStorage interface {
	io.Closer

	// CreateOrder creates a new model.Order.
	CreateOrder(ctx context.Context, order model.Order) (model.Order, error)

	// GetOrder returns model.Order by its number if exists.
	GetOrder(ctx context.Context, orderNum uint64) (model.Order, error)

	// GetOrdersByUser returns *model.Orders by its login if exists.
	GetOrdersByUser(ctx context.Context, login string) (*model.Orders, error)

	// GetOrdersByStatus returns *model.Orders by its status if exists.
	GetOrdersByStatus(ctx context.Context, status model.OrderStatus) (model.Orders, error)

	// UpdateOrder updates model.Order.
	UpdateOrder(ctx context.Context, order model.Order) (model.Order, error)
}
