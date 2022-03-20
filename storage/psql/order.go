package psql

import (
	"context"
	"errors"
	"github.com/IlyaYP/diploma/model"
	"github.com/IlyaYP/diploma/pkg"
	"github.com/IlyaYP/diploma/pkg/logging"
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
func (svc *Storage) GetOrder(ctx context.Context, orderNum string) (model.Order, error) {
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

// GetOrdersByUser returns *model.Orders by its login if exists.
func (svc *Storage) GetOrdersByUser(ctx context.Context, login string) (*model.Orders, error) {
	logger := svc.Logger(ctx)
	var orders model.Orders
	ordersRows, err := svc.pool.Query(
		ctx,
		"select * from orders where login=$1 ORDER BY uploaded_at ASC",
		login,
	)
	if err != nil {
		logger.Err(err).Msg("GetOrdersByUser")
		return nil, err //pgx.ErrNoRows
	}
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

	if len(orders) == 0 {
		return nil, pkg.ErrNoData
	}

	return &orders, nil
}

// GetOrdersByStatus returns *model.Orders by its status if exists.
func (svc *Storage) GetOrdersByStatus(ctx context.Context, status model.OrderStatus) (model.Orders, error) {
	logger := svc.Logger(ctx)
	var orders model.Orders
	ordersRows, err := svc.pool.Query(
		ctx,
		"select * from orders where status=$1 ORDER BY uploaded_at ASC",
		status.Int(),
	)
	if err != nil {
		logger.Err(err).Msg("GetOrdersByStatus")
		return nil, err //pgx.ErrNoRows
	}
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
			logger.Err(err).Msg("GetOrdersByStatus")
			continue
		}
		order.Status = model.NewOrderStatusFromInt(orderStatus)
		orders = append(orders, order)
	}

	if len(orders) == 0 {
		return nil, pkg.ErrNoData
	}

	return orders, nil

}

// UpdateOrder updates  model.Order.
func (svc *Storage) UpdateOrder(ctx context.Context, order model.Order) (model.Order, error) {
	logger := svc.Logger(ctx)
	logger.UpdateContext(order.GetLoggerContext)

	_, err := svc.pool.Exec(ctx,
		`update orders set status=$1, accrual=$2 where num=$3`,
		order.Status.Int(),
		order.Accrual,
		order.Number,
	)

	if err != nil {
		//var pgErr *pgconn.PgError
		//if errors.As(err, &pgErr) {
		//	if pgErr.Code == "23505" {
		//		logger.Err(err).Msg("Error updating order")
		//		return model.Order{}, pkg.ErrAlreadyExists
		//	}
		//}

		logger.Err(err).Msg("Error updating order")
		return model.Order{}, err
	}

	return order, nil
}

// GetBalanceByUser return model.Balance
//SELECT t1.current,  t2.withdrawn
//FROM (select SUM(accrual) as current from orders where login='vasya' and status=4) as t1,
//(select SUM(sum) as withdrawn from withdrawals where login='vasya') as t2;
func (svc *Storage) GetBalanceByUser(ctx context.Context, login string) (model.Balance, error) {
	logger := svc.Logger(ctx)
	balance := model.Balance{}

	err := svc.pool.QueryRow(ctx,
		`SELECT * FROM 
			(select COALESCE(SUM(accrual), 0) as current from orders where login=$1 and status=4) as t1,
			(select COALESCE(SUM(sum), 0) as withdrawn from withdrawals where login=$1) as t2;`,
		login,
	).Scan(
		&balance.Current,
		&balance.Withdrawn,
	)
	switch err {
	case nil:
		balance.Current = balance.Current - balance.Withdrawn
		return balance, nil
	case pgx.ErrNoRows:
		logger.Err(err).Msg("GetBalanceByUser")
		return model.Balance{}, pkg.ErrNotExists
	default:
		logger.Err(err).Msg("GetBalanceByUser")
		return model.Balance{}, err
	}
}

// GetWithdrawalsByUser returns *model.Withdrawals by login if exists.
func (svc *Storage) GetWithdrawalsByUser(ctx context.Context, login string) (*model.Withdrawals, error) {
	logger := svc.Logger(ctx)
	var withdrawals model.Withdrawals
	withdrawalsRows, err := svc.pool.Query(
		ctx,
		"select * from withdrawals where login=$1 ORDER BY processed_at ASC",
		login,
	)
	if err != nil {
		logger.Err(err).Msg("GetWithdrawalsByUser")
		return nil, err //pgx.ErrNoRows
	}
	defer withdrawalsRows.Close()

	for withdrawalsRows.Next() {
		var withdrawal model.Withdrawal
		err := withdrawalsRows.Scan( //ordernum | sum | processed_at | login
			&withdrawal.Order,
			&withdrawal.Sum,
			&withdrawal.ProcessedAt,
			&withdrawal.User,
		)
		if err != nil {
			logger.Err(err).Msg("GetWithdrawalsByUser")
			continue
		}
		withdrawals = append(withdrawals, withdrawal)
	}

	if len(withdrawals) == 0 {
		return nil, pkg.ErrNoData
	}

	return &withdrawals, nil
}

// NewWithdrawal Apply new Withdrawal to DB
// insert into withdrawals(ordernum, sum, login) values (22345678902, 35, 'vasya');
func (svc *Storage) NewWithdrawal(ctx context.Context, withdrawal model.Withdrawal) error {
	ctx, _ = logging.GetCtxLogger(ctx) // correlationID is created here
	logger := svc.Logger(ctx)
	logger.UpdateContext(withdrawal.GetLoggerContext)

	_, err := svc.pool.Exec(ctx, `insert into withdrawals(ordernum, sum, login) values ($1, $2, $3)`,
		withdrawal.Order, withdrawal.Sum, withdrawal.User)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				logger.Err(pgErr).Msg("Error insert into withdrawals")
				return pkg.ErrAlreadyExists
			}
		}

		logger.Err(err).Msg("Error insert into withdrawals")
		return err
	}

	logger.Info().Msg("Successfully created order")

	return nil
}
