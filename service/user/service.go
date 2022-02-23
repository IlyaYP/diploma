package user

import (
	"context"
	"crypto/hmac"
	"fmt"
	"github.com/IlyaYP/diploma/model"
	"github.com/IlyaYP/diploma/pkg"
	"github.com/IlyaYP/diploma/pkg/logging"
	"github.com/IlyaYP/diploma/storage"
	"github.com/rs/zerolog"
)

var _ Service = (*service)(nil)

const (
	serviceName = "user-service"
)

type (
	service struct {
		config      Config
		userStorage storage.UserStorage
	}

	Option func(svc *service) error
)

// WithUserStorage sets storage.UserStorage.
func WithUserStorage(st storage.UserStorage) Option {
	return func(svc *service) error {
		svc.userStorage = st
		return nil
	}
}

// WithConfig sets Config.
func WithConfig(cfg Config) Option {
	return func(svc *service) error {
		svc.config = cfg
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

	//if err := svc.config.validate(); err != nil {
	//	return nil, fmt.Errorf("config validation: %w", err)
	//}

	if svc.userStorage == nil {
		return nil, fmt.Errorf("userStorage: nil")
	}

	return svc, nil
}

// Register a new user.
func (svc *service) Register(ctx context.Context, login, password string) (model.User, error) {
	ctx, _ = logging.GetCtxLogger(ctx) // correlationID is created here

	logger := svc.Logger(ctx)

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

	user, err := svc.userStorage.CreateUser(ctx, model.User{Login: login, Password: pkg.Hash(password, login)})
	if err != nil {
		logger.Err(err).Msg("Error register user")
		return model.User{}, fmt.Errorf("register user: %w", err)
	}

	logger.Info().Msg("Successfully registered user")
	return user, nil
}

// Login Authenticates user
func (svc *service) Login(ctx context.Context, login, password string) (model.User, error) {
	ctx, _ = logging.GetCtxLogger(ctx) // correlationID is created here
	logger := svc.Logger(ctx)

	// Build input
	input := model.User{
		Login:    login,
		Password: password,
	}

	logger.UpdateContext(input.GetLoggerContext)

	user, err := svc.userStorage.GetUserByLogin(ctx, login)
	if err != nil {
		logger.Err(err).Msg("Login Unsuccessful")
		return model.User{}, pkg.ErrInvalidLogin
	}

	if !hmac.Equal([]byte(pkg.Hash(password, login)), []byte(user.Password)) {
		logger.Err(pkg.ErrInvalidPassword).Msg("Login Unsuccessful")
		return model.User{}, pkg.ErrInvalidLogin
	}

	logger.Info().Msg("Login Success")
	return *user, nil
}

// Logger returns logger with service field set.
func (svc *service) Logger(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, serviceName).Logger()

	return &logger
}
