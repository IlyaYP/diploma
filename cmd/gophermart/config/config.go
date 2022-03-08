package config

import (
	"context"
	"flag"
	"fmt"
	"github.com/IlyaYP/diploma/api/server"
	"github.com/IlyaYP/diploma/api/server/handler"
	"github.com/IlyaYP/diploma/pkg/logging"
	"github.com/IlyaYP/diploma/service/order"
	"github.com/IlyaYP/diploma/service/user"
	"github.com/IlyaYP/diploma/storage/psql"
	"github.com/caarlos0/env/v6"
)

// Config combines sub-configs for all services, storages and providers.
type Config struct {
	UserService    user.Config
	PSQLStorage    psql.Config
	APISever       server.Config
	Address        string `env:"RUN_ADDRESS"`
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	DSN            string `env:"DATABASE_URI"`
}

// New initializes a new config.
func New() (*Config, error) {
	cfg := Config{}
	flag.StringVar(&cfg.DSN, "d", "postgres://postgres:postgres@localhost:5432/gmart", "DATABASE_URI")
	flag.StringVar(&cfg.Address, "a", "localhost:8081", "RUN_ADDRESS")
	flag.StringVar(&cfg.AccrualAddress, "r", "localhost:8080", "ACCRUAL_SYSTEM_ADDRESS")
	flag.Parse()

	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("initializing config: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("validating config: %w", err)
	}

	cfg.PSQLStorage = psql.NewDefaultConfig()
	cfg.PSQLStorage.DSN = cfg.DSN
	cfg.APISever.Address = cfg.Address

	return &cfg, nil
}

// validate performs a basic validation.
func (c Config) validate() error {
	logger := logging.NewLogger()
	if c.Address == "" {
		return fmt.Errorf("%s field: empty", "RUN_ADDRESS")
	}
	logger = logger.With().Str("RUN_ADDRESS", c.Address).Logger()

	if c.DSN == "" {
		return fmt.Errorf("%s field: empty", "DATABASE_URI")
	}
	logger = logger.With().Str("DATABASE_URI", c.DSN).Logger()

	if c.AccrualAddress == "" {
		return fmt.Errorf("%s field: empty", "ACCRUAL_SYSTEM_ADDRESS")
	}
	logger = logger.With().Str("ACCRUAL_SYSTEM_ADDRESS", c.AccrualAddress).Logger()

	logger.Info().Msg("Initialized with args:")
	return nil
}

// BuildPsqlStorage builds psql.Storage dependency.
func (c Config) BuildPsqlStorage(ctx context.Context) (*psql.Storage, error) {
	st, err := psql.New(
		psql.WithConfig(c.PSQLStorage),
		psql.WithContext(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("building psql storage: %w", err)
	}

	return st, nil
}

// BuildUserService builds user.Service dependency.
func (c Config) BuildUserService(ctx context.Context) (user.Service, error) {
	st, err := c.BuildPsqlStorage(ctx)
	if err != nil {
		return nil, err
	}

	svc, err := user.New(
		user.WithConfig(c.UserService),
		user.WithUserStorage(st),
	)

	if err != nil {
		return nil, fmt.Errorf("building user service: %w", err)
	}

	return svc, nil

}

// BuildOrderService builds order.Service dependency.
func (c Config) BuildOrderService(ctx context.Context) (order.Service, error) {
	st, err := c.BuildPsqlStorage(ctx)
	if err != nil {
		return nil, err
	}

	svc, err := order.New(
		order.WithOrderStorage(st),
	)
	if err != nil {
		return nil, fmt.Errorf("building test service: %w", err)
	}

	return svc, nil
}

// BuildServer builds REST API Server dependency.
func (c Config) BuildServer(ctx context.Context) (*server.Server, error) {
	//userSvc, err := c.BuildUserService(ctx)
	//if err != nil {
	//	return nil, fmt.Errorf("building server: %w", err)
	//}
	//
	//orderSvc, err := c.BuildOrderService(ctx)
	//if err != nil {
	//	return nil, fmt.Errorf("building server: %w", err)
	//}

	// Build Storage
	st, err := c.BuildPsqlStorage(ctx)
	if err != nil {
		return nil, err
	}

	// Build User Service
	userSvc, err := user.New(
		user.WithConfig(c.UserService),
		user.WithUserStorage(st),
	)

	if err != nil {
		return nil, fmt.Errorf("building user service: %w", err)
	}

	// Build Order Service
	orderSvc, err := order.New(
		order.WithOrderStorage(st),
	)
	if err != nil {
		return nil, fmt.Errorf("building order service: %w", err)
	}

	// Build REST API Service
	r, err := handler.NewHandler(
		handler.WithUserService(userSvc),
		handler.WithOrderService(orderSvc),
	)
	if err != nil {
		return nil, fmt.Errorf("building server: %w", err)
	}

	s, err := server.New(
		server.WithConfig(&c.APISever),
		server.WithRouter(r),
	)
	if err != nil {
		return nil, fmt.Errorf("building server: %w", err)
	}
	return s, nil
}
