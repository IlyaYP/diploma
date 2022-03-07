package psql

import (
	"context"
	"errors"
	"github.com/IlyaYP/diploma/model"
	"github.com/IlyaYP/diploma/pkg"
	"github.com/jackc/pgconn"
)

// CreateOrder creates a new model.Order.
func (svc *service) CreateOrder(ctx context.Context, order model.Order) (model.Order, error) {
	logger := svc.Logger(ctx)
	logger.UpdateContext(order.GetLoggerContext)

	_, err := svc.pool.Exec(ctx, `insert into orders(num, login, status) values ($1, $2, $3)`,
		order.Number, order.User, order.Status.Int())
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				logger.Err(err).Msg("Error creating order")
				return model.Order{}, pkg.ErrAlreadyExists
			}
		}

		logger.Err(err).Msg("Error creating order")
		return model.Order{}, err
	}

	logger.Info().Msg("Successfully created order")

	return order, nil
}

// GetOrdersByUser returns *[]model.Order by its login if exists.
func (svc *service) GetOrdersByUser(ctx context.Context, login string) (*[]model.Order, error) {
	return nil, nil
}
