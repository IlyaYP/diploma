package psql

type Config struct {
	DSN string `env:"DATABASE_URI"`
}

func (c Config) validate() error {
	return nil
}
