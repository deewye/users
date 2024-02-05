package config

import (
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pkg/errors"
	"github.com/vrischmann/envconfig"
)

type (
	Config struct {
		Env      string `envconfig:"default=PROD"`
		Postgres struct {
			Master *DatabaseConfig
			Slave  *DatabaseConfig
		}
	}

	DatabaseConfig struct {
		DSN         string `required:"true"`
		MaxIdleConn int    `envconfig:"default=2"`
		MaxOpenConn int    `envconfig:"default=5"`
	}
)

func InitConfig(prefix string) (*Config, error) {
	config := &Config{}
	if err := envconfig.InitWithPrefix(config, prefix); err != nil {
		return nil, errors.Wrap(err, "can't init envconfig")
	}

	return config, nil
}
