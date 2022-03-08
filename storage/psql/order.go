package psql

import (
	"context"
	"errors"
	"github.com/IlyaYP/diploma/model"
	"github.com/IlyaYP/diploma/pkg"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// CreateOrder creates a new model.Order.
func (svc *Storage) CreateOrder(ctx context.Context, order model.Order) (model.Order, error) {
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

// GetOrder returns model.Order by its number if exists.
func (svc *Storage) GetOrder(ctx context.Context, orderNum uint64) (model.Order, error) {
	logger := svc.Logger(ctx)
	var orderStatus int
	order := model.Order{}
	err := svc.pool.QueryRow(ctx, "select * from orders where num=$1", orderNum).Scan(
		&order.Number,
		&orderStatus,
		&order.Accrual,
		&order.UploadedAt,
		&order.User,
	)
	switch err {
	case nil:
		order.Status = model.NewOrderStatusFromInt(orderStatus)
		return order, nil
	case pgx.ErrNoRows:
		//logger.Err(err).Msg("GetOrder")
		return model.Order{}, pkg.ErrNotExists
	default:
		logger.Err(err).Msg("GetOrder")
		return model.Order{}, err
	}
}

// GetOrdersByUser returns *[]model.Order by its login if exists.
func (svc *Storage) GetOrdersByUser(ctx context.Context, login string) (*model.Orders, error) {
	logger := svc.Logger(ctx)
	var orders model.Orders
	ordersRows, _ := svc.pool.Query(ctx, "select * from orders where login=$1", login)
	defer ordersRows.Close()

	for ordersRows.Next() {
		order := model.Order{}
		var orderStatus int
		err := ordersRows.Scan(
			&order.Number,
			&orderStatus,
			&order.Accrual,
			&order.UploadedAt,
			&order.User,
		)
		if err != nil {
			logger.Err(err).Msg("GetOrdersByUser")
			continue
		}
		order.Status = model.NewOrderStatusFromInt(orderStatus)
		orders = append(orders, order)
	}

	return &orders, nil
}
