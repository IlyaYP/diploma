package mock

import (
	"fmt"
	"time"
)

const (
	defaultAccrualAddress = "localhost:8080"
	defaultConfigTimeOut  = 10
)

type Config struct {
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	timeout        time.Duration
}

// validate performs a basic validation.
func (c Config) validate() error {
	if c.AccrualAddress == "" {
		return fmt.Errorf("%s field: empty", "ACCRUAL_SYSTEM_ADDRESS")
	}
	if c.timeout == 0 {
		return fmt.Errorf("%s field: empty", "timeout")
	}

	return nil
}

// NewDefaultConfig builds a Config with default values.
func NewDefaultConfig() Config {
	return Config{
		AccrualAddress: defaultAccrualAddress,
		timeout:        time.Duration(defaultConfigTimeOut) * time.Second,
	}
}
