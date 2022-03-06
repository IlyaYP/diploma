package psql

import (
	"context"
	"github.com/IlyaYP/diploma/model"
)

// CreateOrder creates a new model.Order.
func (svc *service) CreateOrder(ctx context.Context, order model.Order) (model.Order, error) {
	return model.Order{}, nil
}

// GetOrdersByUser returns *[]model.Order by its login if exists.
func (svc *service) GetOrdersByUser(ctx context.Context, login string) (*[]model.Order, error) {
	return nil, nil
}
