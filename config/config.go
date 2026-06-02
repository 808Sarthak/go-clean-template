package config

import (
	"fmt"
	"net/url"
	"time"

	"github.com/caarlos0/env/v11"
)

type (
	// Config -.
	Config struct {
		App     app
		HTTP    http
		Log     log
		DB      db
		GRPC    grpc
		RMQ     rmq
		NATS    nats
		JWT     jwt
		Metrics metrics
		Swagger swagger
		Redis   redis
	}

	// App -.
	app struct {
		Name    string `env:"APP_NAME,required"`
		Version string `env:"APP_VERSION,required"`
	}

	// HTTP -.
	http struct {
		Port           string `env:"HTTP_PORT,required"`
		UsePreforkMode bool   `env:"HTTP_USE_PREFORK_MODE" envDefault:"false"`
	}

	// Log -.
	log struct {
		Level string `env:"LOG_LEVEL,required"`
	}

	// DB -.
	db struct {
		PoolMax  int    `env:"DB_POOL_MAX" envDefault:"2"`
		User     string `env:"DB_USER" envDefault:"root"`
		Password string `env:"DB_PASSWORD" envDefault:"root"`
		Host     string `env:"DB_HOST" envDefault:"127.0.0.1"`
		Port     string `env:"DB_PORT" envDefault:"3306"`
		Name     string `env:"DB_NAME" envDefault:"myplex_sms"`
	}

	// GRPC -.
	grpc struct {
		Port string `env:"GRPC_PORT,required"`
	}

	// RMQ -.
	rmq struct {
		ServerExchange string `env:"RMQ_RPC_SERVER,required"`
		ClientExchange string `env:"RMQ_RPC_CLIENT,required"`
		URL            string `env:"RMQ_URL,required"`
	}

	// NATS -.
	nats struct {
		ServerExchange string `env:"NATS_RPC_SERVER,required"`
		URL            string `env:"NATS_URL,required"`
	}

	// JWT -.
	jwt struct {
		Secret      string        `env:"JWT_SECRET,required"`
		TokenExpiry time.Duration `env:"JWT_TOKEN_EXPIRY" envDefault:"24h"`
	}

	// Metrics -.
	metrics struct {
		Enabled bool `env:"METRICS_ENABLED" envDefault:"true"`
	}

	// Swagger -.
	swagger struct {
		Enabled bool `env:"SWAGGER_ENABLED" envDefault:"false"`
	}

	// Redis -.
	redis struct {
		URL string `env:"REDIS_URL,required"`
	}
)

// DSN builds a MySQL connection string from DB settings.
func (c *Config) DSN() string {
	user := url.QueryEscape(c.DB.User)
	pass := url.QueryEscape(c.DB.Password)

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local&charset=utf8mb4&multiStatements=true",
		user, pass, c.DB.Host, c.DB.Port, c.DB.Name)
}

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
