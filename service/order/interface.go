package order

import (
	"context"
	"github.com/IlyaYP/diploma/model"
)

type Service interface {
	// CreateOrder creates a new model.Order.
	CreateOrder(ctx context.Context, order model.Order) (model.Order, error)

	// GetOrder returns model.Order by its number if exists.
	GetOrder(ctx context.Context, orderNum string) (model.Order, error)

	// GetOrdersByUser returns *model.Orders by its login if exists.
	GetOrdersByUser(ctx context.Context, login string) (*model.Orders, error)

	// ProcessOrder do something with order
	ProcessOrder(ctx context.Context, order model.Order) error

	// ProcessOrders do something with order
	ProcessOrders(ctx context.Context) error

	// ProcessNewOrders do something with order
	ProcessNewOrders(ctx context.Context) error

	// GetBalanceByUser return model.Balance
	GetBalanceByUser(ctx context.Context, login string) (model.Balance, error)

	// GetWithdrawalsByUser returns *model.Withdrawals by login if exists.
	GetWithdrawalsByUser(ctx context.Context, login string) (*model.Withdrawals, error)

	// NewWithdrawal Do new Withdrawal for login if good Balance
	NewWithdrawal(ctx context.Context, withdrawal model.Withdrawal) error
}
