package psql

import (
	"context"
	"fmt"
	"github.com/IlyaYP/diploma/model"
	"github.com/IlyaYP/diploma/pkg/logging"
	"github.com/IlyaYP/diploma/storage"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
)

var _ storage.UserStorage = (*service)(nil)

const (
	serviceName = "psql"

	dbTableLoggingKey     = "db-table"
	dbOperationLoggingKey = "db-operation"
)

type (
	service struct {
		config Config
		pool   *pgxpool.Pool
		ctx    context.Context
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

// WithConfig sets Config.
func WithContext(ctx context.Context) option {
	return func(svc *service) error {
		svc.ctx = ctx
		return nil
	}
}

// New creates a new service.
func New(opts ...option) (*service, error) {
	svc := &service{
		config: NewDefaultConfig(),
	}

	for _, opt := range opts {
		if err := opt(svc); err != nil {
			return nil, fmt.Errorf("initialising dependencies: %w", err)
		}
	}

	if err := svc.config.validate(); err != nil {
		return nil, fmt.Errorf("Config validation: %w", err)
	}

	pool, err := pgxpool.Connect(svc.ctx, svc.config.DSN)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %w", err)
	}
	svc.pool = pool

	if err := svc.Ping(svc.ctx); err != nil {
		return nil, fmt.Errorf("ping for DSN (%s) failed: %w", svc.config.DSN, err)
	}

	if err := svc.Migrate(svc.ctx); err != nil {
		return nil, fmt.Errorf("Unable to create table: %w", err)
	}

	return svc, nil
}

func (svc *service) Migrate(ctx context.Context) error {
	logger := svc.Logger(ctx)
	logger.Info().Msg("Creating Table")

	_, err := svc.pool.Exec(ctx,
		"CREATE TABLE IF NOT EXISTS users ( login varchar(40) primary key, password varchar(64));",
		//"CREATE TABLE IF NOT EXISTS counters ( id varchar(40) primary key, delta bigint);",
	)

	return err
}

// Ping checks db connection
func (svc *service) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, svc.config.timeout)
	defer cancel()

	return svc.pool.Ping(ctx)
}

// Close closes DB connection.
func (svc *service) Close() error {
	if svc.pool == nil {
		return nil
	}
	svc.pool.Close()
	return nil
}

// CreateUser creates a new model.User.
// Returns ErrAlreadyExists if user exists.
func (svc *service) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	logger := svc.Logger(ctx)
	logger.UpdateContext(user.GetLoggerContext)
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
