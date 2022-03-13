package order

import (
	"context"
	"fmt"
	"github.com/IlyaYP/diploma/model"
	"github.com/IlyaYP/diploma/pkg"
	"github.com/IlyaYP/diploma/pkg/logging"
	"github.com/IlyaYP/diploma/provider/accrual"
	"github.com/IlyaYP/diploma/storage"
	"github.com/rs/zerolog"
)

var _ Service = (*service)(nil)

const (
	serviceName = "order-service"
)

type (
	service struct {
		OrderStorage    storage.OrderStorage
		AccrualProvider accrual.Provider
	}

	Option func(svc *service) error
)

// WithAccrualProvider sets accrual.Provider.
func WithAccrualProvider(pr accrual.Provider) Option {
	return func(svc *service) error {
		svc.AccrualProvider = pr
		return nil
	}
}

// WithOrderStorage sets storage.UserStorage.
func WithOrderStorage(st storage.OrderStorage) Option {
	return func(svc *service) error {
		svc.OrderStorage = st
		return nil
	}
}

// New creates a new service.
func New(opts ...Option) (*service, error) {
	svc := &service{}
	for _, opt := range opts {
		if err := opt(svc); err != nil {
			return nil, fmt.Errorf("initialising dependencies: %w", err)
		}
	}

	if svc.OrderStorage == nil {
		return nil, fmt.Errorf("OrderStorage: nil")
	}

	if svc.AccrualProvider == nil {
		return nil, fmt.Errorf("AccrualProvider: nil")
	}
	return svc, nil

}

// CreateOrder creates a new model.Order.
func (svc *service) CreateOrder(ctx context.Context, order model.Order) (model.Order, error) {
	order, err := svc.OrderStorage.CreateOrder(ctx, order)
	if err != nil {
		return model.Order{}, err
	}

	go svc.ProcessOrder(context.Background(), order)
	return order, nil
}

// GetOrder returns model.Order by its number if exists.
func (svc *service) GetOrder(ctx context.Context, orderNum string) (model.Order, error) {
	return svc.OrderStorage.GetOrder(ctx, orderNum)
}

// GetOrdersByUser returns *[]model.Order by its login if exists.
func (svc *service) GetOrdersByUser(ctx context.Context, login string) (*model.Orders, error) {
	return svc.OrderStorage.GetOrdersByUser(ctx, login)
}

// ProcessOrder do something with order
func (svc *service) ProcessOrder(ctx context.Context, order model.Order) error {
	//ctx, _ = logging.GetCtxLogger(ctx) // correlationID is created here
	logger := svc.Logger(ctx)

	rOrder, err := svc.AccrualProvider.ProcessOrder(ctx, order)
	if err != nil {
		logger.Err(err).Msg("ProcessOrder error.")
		return err
	}

	if _, err := svc.OrderStorage.UpdateOrder(ctx, rOrder); err != nil {
		logger.Err(err).Msg("ProcessOrder: UpdateOrder error.")
		return err
	}

	logger.Info().Msgf("ProcessOrder success, order status:%s", rOrder.Status)
	return nil
}

func (svc *service) ProcessOrders(ctx context.Context) error {
	orders, err := svc.OrderStorage.GetOrdersByStatus(ctx, model.OrderStatusProcessing)
	if err != nil {
		return err
	}

	for i, _ := range orders {
		//// TODO: request from accrual
		//orders[i].Accrual = float64(int(100000*rand.Float32())) / 100
		//orders[i].Status = model.OrderStatusProcessed
		//if _, err := svc.OrderStorage.UpdateOrder(ctx, orders[i]); err != nil {
		//	return err
		//}
		svc.ProcessOrder(ctx, orders[i])
	}

	return nil
}

func (svc *service) ProcessNewOrders(ctx context.Context) error {
	orders, err := svc.OrderStorage.GetOrdersByStatus(ctx, model.OrderStatusNew)
	if err != nil {
		return err
	}

	for i, _ := range orders {
		//// TODO: send to accrual
		//orders[i].Status = model.OrderStatusProcessing
		//if _, err := svc.OrderStorage.UpdateOrder(ctx, orders[i]); err != nil {
		//	return err
		//}
		svc.ProcessOrder(ctx, orders[i])
	}

	return nil
}

// Logger returns logger with service field set.
func (svc *service) Logger(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, serviceName).Logger()

	return &logger
}

// GetBalanceByUser return model.Balance
func (svc *service) GetBalanceByUser(ctx context.Context, login string) (model.Balance, error) {
	ctx, _ = logging.GetCtxLogger(ctx) // correlationID is created here
	logger := svc.Logger(ctx)

	// **** TODO: Remove
	// Temp, only for Debug
	if err := svc.ProcessNewOrders(ctx); err != nil {
		if err != pkg.ErrNoData {
			logger.Err(err).Msg("ProcessNewOrders:")
		}
	}

	if err := svc.ProcessOrders(ctx); err != nil {
		if err != pkg.ErrNoData {
			logger.Err(err).Msg("ProcessOrders:")
		}
	}
	// ****

	return svc.OrderStorage.GetBalanceByUser(ctx, login)
}

// GetWithdrawalsByUser returns *model.Withdrawals by login if exists.
func (svc *service) GetWithdrawalsByUser(ctx context.Context, login string) (*model.Withdrawals, error) {
	return svc.OrderStorage.GetWithdrawalsByUser(ctx, login)
}

// NewWithdrawal Do new Withdrawal for login if good balance
func (svc *service) NewWithdrawal(ctx context.Context, withdrawal model.Withdrawal) error {
	ctx, _ = logging.GetCtxLogger(ctx) // correlationID is created here
	logger := svc.Logger(ctx)

	balance, err := svc.GetBalanceByUser(ctx, withdrawal.User)
	if err != nil {
		logger.Err(err).Msg("NewWithdrawal: can't get balance")
		//return err
	}

	// check if enough balance
	if withdrawal.Sum+balance.Withdrawn > balance.Current {
		logger.Err(pkg.ErrInsufficientBalance).
			Msgf("NewWithdrawal: current balance:%v less then withdrawal sum: %v",
				balance.Current-balance.Withdrawn,
				withdrawal.Sum,
			)
		return pkg.ErrInsufficientBalance
	}

	// add withdrawal to db
	if err := svc.OrderStorage.NewWithdrawal(ctx, withdrawal); err != nil {
		logger.Err(err).Msg("NewWithdrawal: can't apply withdrawal to DB")
		return err
	}

	logger.Info().Msgf("Success NewWithdrawal sum: %v by %v",
		withdrawal.Sum,
		withdrawal.User,
	)
	return nil
}
