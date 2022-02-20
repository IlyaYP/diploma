package user

import (
	"context"
	"fmt"
	"github.com/IlyaYP/diploma/model"
	"github.com/IlyaYP/diploma/pkg"
	"github.com/IlyaYP/diploma/pkg/logging"
	"github.com/rs/zerolog"
)

var _ Service = (*service)(nil)

const (
	serviceName = "user-service"
)

type (
	service struct {
	}

	Option func(service *service) error
)

// New creates a new service.
func New(opts ...Option) (*service, error) {
	svc := &service{}
	for _, opt := range opts {
		if err := opt(svc); err != nil {
			return nil, fmt.Errorf("initialising dependencies: %w", err)
		}
	}

	//if err := svc.config.Validate(); err != nil {
	//	return nil, fmt.Errorf("config validation: %w", err)
	//}

	//if svc.userStorage == nil {
	//	return nil, fmt.Errorf("userStorage: nil")
	//}

	return svc, nil
}

// CreateUser creates a new user.
func (svc *service) CreateUser(ctx context.Context, login, password string) (model.User, error) {
	logger := svc.Logger(ctx)
	logger.Info().Msg("Creating user")

	// Input checks
	if login == "" {
		logger.Error().Msg("login: empty")
		return model.User{}, fmt.Errorf("%w: login: empty", pkg.ErrInvalidInput)
	}
	if password == "" {
		logger.Error().Msg("password: empty")
		return model.User{}, fmt.Errorf("%w: password: empty", pkg.ErrInvalidInput)
	}

	// Build input
	input := model.User{
		Login:    login,
		Password: password,
	}

	logger.UpdateContext(input.GetLoggerContext)
	logger.Info().Msg("Creating user")

	//logger.Warn().Err(err).Msg("Creating user")
	//
	//logger.UpdateContext(user.GetLoggerContext)
	return model.User{Login: login, Password: password}, nil
}

// Login Authenticates user
func (svc *service) Login(ctx context.Context, login, password string) (model.User, error) {
	return model.User{}, nil
}

// Logger returns logger with service field set.
func (svc *service) Logger(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, serviceName).Logger()

	return &logger
}
