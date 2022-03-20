package accrual

import (
	"context"
	"github.com/IlyaYP/diploma/model"
)

type Provider interface {

	// ProcessOrder request accrual for order
	ProcessOrder(ctx context.Context, order model.Order) (model.Order, error)
}
