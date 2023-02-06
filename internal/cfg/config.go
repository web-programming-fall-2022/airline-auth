package cfg

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/web-programming-fall-2022/airline-auth/internal/bootstrap"
)

type Config struct {
	Env string
	Log struct {
		Level string
	}

	bootstrap.GrpcServerRunnerConfig `mapstructure:",squash" yaml:",inline"`

	HttpServer struct {
		Port int
	}

	Redis struct {
		Addr string
	}
}

func (c *Config) Validate() error {
	return validation.Errors{
		"env": validation.Validate(c.Env, validation.Required),
		"log.level": validation.Validate(c.Log.Level, validation.Required, validation.In(
			"panic", "fatal", "error", "warn", "info", "debug", "trace",
		)),
		"server.host": validation.Validate(c.Server.Host, validation.Required),
		"server.port": validation.Validate(c.Server.Port, validation.Required),
	}.Filter()
}
