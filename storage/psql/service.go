package psql

import (
	"context"
	"fmt"
	"github.com/IlyaYP/diploma/model"
	"github.com/IlyaYP/diploma/pkg/logging"
	"github.com/IlyaYP/diploma/storage"
	"github.com/rs/zerolog"
)

var _ storage.UserStorage = (*service)(nil)

const (
	serviceName = "psql"
)

type (
	service struct {
		config Config
	}

	option func(svc *service) error
)

// WithConfig sets Config.
func WithConfig(cfg Config) option {
	return func(svc *service) error {
		svc.config = cfg
		return nil
	}
}

// New creates a new service.
func New(opts ...option) (*service, error) {
	svc := &service{}
	for _, opt := range opts {
		if err := opt(svc); err != nil {
			return nil, fmt.Errorf("initialising dependencies: %w", err)
		}
	}

	//if err := svc.Config.Validate(); err != nil {
	//	return nil, fmt.Errorf("Config validation: %w", err)
	//}

	//if svc.userStorage == nil {
	//	return nil, fmt.Errorf("userStorage: nil")
	//}

	return svc, nil
}

// CreateUser creates a new model.User.
// Returns ErrAlreadyExists if user exists.
func (svc *service) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	logger := svc.Logger(ctx)
	//logger.UpdateContext(user.GetLoggerContext)
	logger.Info().Msg("Creating user")
	return user, nil
}

// GetUserByLogin returns model.User by its login if exists.
func (svc *service) GetUserByLogin(ctx context.Context, login string) (*model.User, error) {
	return &model.User{}, nil
}

// Logger returns logger with service field set.
func (svc *service) Logger(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, serviceName).Logger()

	return &logger
}
