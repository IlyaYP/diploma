package config

import (
	"fmt"
	"github.com/IlyaYP/diploma/service/user"
	"github.com/IlyaYP/diploma/storage"
	"github.com/IlyaYP/diploma/storage/psql"
)

// Config combines sub-configs for all services, storages and providers.
type Config struct {
	UserService user.Config
	PSQLStorage psql.Config
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
