package order

import (
	"context"
	"github.com/IlyaYP/diploma/model"
)

type Service interface {
	// CreateOrder creates a new model.Order.
	CreateOrder(ctx context.Context, order model.Order) (model.Order, error)

	// GetOrder returns model.Order by its number if exists.
	GetOrder(ctx context.Context, orderNum uint64) (model.Order, error)

	// GetOrdersByUser returns *[]model.Order by its login if exists.
	GetOrdersByUser(ctx context.Context, login string) (*[]model.Order, error)
}
