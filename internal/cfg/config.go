package cfg

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/web-programming-fall-2022/airline-auth/internal/bootstrap"
	"github.com/web-programming-fall-2022/airline-auth/internal/storage"
)

type Config struct {
	Env string
	Log struct {
		Level string
	}

	bootstrap.GrpcServerRunnerConfig `mapstructure:",squash" yaml:",inline"`

	MainDB storage.DBConfig `mapstructure:"main_db"`

	JWT struct {
		Secret             string
		AuthTokenExpire    int64 `mapstructure:"auth_token_expire"`
		RefreshTokenExpire int64 `mapstructure:"refresh_token_expire"`
	}

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
