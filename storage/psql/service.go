package psql

import (
	"context"
	"fmt"
	"github.com/IlyaYP/diploma/pkg/logging"
	"github.com/IlyaYP/diploma/storage"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
)

var _ storage.UserStorage = (*Storage)(nil)
var _ storage.OrderStorage = (*Storage)(nil)

const (
	serviceName = "psql"

	dbTableLoggingKey     = "db-table"
	dbOperationLoggingKey = "db-operation"
)

type (
	Storage struct {
		config Config
		pool   *pgxpool.Pool
		ctx    context.Context
	}

	option func(svc *Storage) error
)

// WithConfig sets Config.
func WithConfig(cfg Config) option {
	return func(svc *Storage) error {
		svc.config = cfg
		return nil
	}
}

// WithContext sets Context.
func WithContext(ctx context.Context) option {
	return func(svc *Storage) error {
		svc.ctx = ctx
		return nil
	}
}

// New creates a new Storage.
func New(opts ...option) (*Storage, error) {
	svc := &Storage{
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

func (svc *Storage) Migrate(ctx context.Context) error {
	logger := svc.Logger(ctx)
	logger.Info().Msg("Creating Tables")

	_, err := svc.pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS users
		(
			login varchar(40) not null,
			password varchar(64) not null,
			primary key (login)
		);
		CREATE TABLE IF NOT EXISTS orders
		(
			num int not null,
			status int not null,
			accrual int,
			uploaded_at timestamp    not null default now(),
			login varchar(64) not null,
			primary key (num),
			foreign key (login) references users (login)
		);
	`)

	return err
}

// Ping checks db connection
func (svc *Storage) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, svc.config.timeout)
	defer cancel()

	return svc.pool.Ping(ctx)
}

// Close closes DB connection.
func (svc *Storage) Close() error {
	if svc.pool == nil {
		return nil
	}
	svc.pool.Close()
	return nil
}

// Logger returns logger with Storage field set.
func (svc *Storage) Logger(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, serviceName).Logger()

	return &logger
}
