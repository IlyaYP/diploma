package order

import (
	"context"
	"fmt"
	"github.com/IlyaYP/diploma/model"
	"github.com/IlyaYP/diploma/pkg/logging"
	"github.com/IlyaYP/diploma/storage"
	"github.com/rs/zerolog"
)

var _ Service = (*service)(nil)

const (
	serviceName = "order-service"
)

type (
	service struct {
		OrderStorage storage.OrderStorage
	}

	Option func(svc *service) error
)

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

	return svc, nil

}

// CreateOrder creates a new model.Order.
func (svc *service) CreateOrder(ctx context.Context, order model.Order) (model.Order, error) {
	return model.Order{}, nil
}

// GetOrdersByUser returns *[]model.Order by its login if exists.
func (svc *service) GetOrdersByUser(ctx context.Context, login string) (*[]model.Order, error) {
	return nil, nil
}

// Logger returns logger with service field set.
func (svc *service) Logger(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, serviceName).Logger()

	return &logger
}