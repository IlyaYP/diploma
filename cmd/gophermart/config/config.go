package config

import (
	"flag"
	"fmt"
	"github.com/IlyaYP/diploma/pkg/logging"
	"github.com/IlyaYP/diploma/service/user"
	"github.com/IlyaYP/diploma/storage"
	"github.com/IlyaYP/diploma/storage/psql"
	"github.com/caarlos0/env/v6"
)

// Config combines sub-configs for all services, storages and providers.
type Config struct {
	UserService    user.Config
	PSQLStorage    psql.Config
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

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validating config: %w", err)
	}

	cfg.PSQLStorage.DSN = cfg.DSN

	return &cfg, nil
}

// Validate performs a basic validation.
func (c Config) Validate() error {
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

	logger.Info().Msg("initialized")
	return nil
}

// BuildPsqlStorage builds psql.Storage dependency.
func (c Config) BuildPsqlStorage() (storage.UserStorage, error) {
	st, err := psql.New(
		psql.WithConfig(c.PSQLStorage),
	)
	if err != nil {
		return nil, fmt.Errorf("building psql storage: %w", err)
	}

	return st, nil
}

// BuildUserService builds user.Processor dependency.
func (c Config) BuildUserService() (user.Service, error) {
	st, err := c.BuildPsqlStorage()
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
