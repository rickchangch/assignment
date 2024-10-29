package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Name string
		Port int
	}
	Log struct {
		Level string
	}
	Connection struct {
		PostgreSQL struct {
			Host         string
			DB           string
			User         string
			Password     string
			MaxConn      int
			MaxConnIdle  int
			UsePgBouncer bool
		}
		Redis struct {
			Host            string
			Password        string
			WriteTimeoutSec time.Duration
			ReadTimeoutSec  time.Duration
		}
	}
}

func Load() (*Config, error) {
	viper.SetConfigFile("./setup/config/config.yaml")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("viper.ReadInConfig failed: %w", err)
	}

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("unmarshal to config struct failed: %w", err)
	}

	return &c, nil
}
