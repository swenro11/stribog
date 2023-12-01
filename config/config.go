package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App        `yaml:"app"`
		HTTP       `yaml:"http"`
		Log        `yaml:"logger"`
		PG         `yaml:"postgres"`
		RMQ        `yaml:"rabbitmq"`
		Redis      `yaml:"redis"`
		Mongo      `yaml:"mongo"`
		PARAM      `yaml:"param"`
	}

	// App -.
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" yaml:"log_level" env:"LOG_LEVEL"`
	}

	// PG -.
	PG struct {
		PoolMax int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
		URL     string `env-required:"true" yaml:"pg_url" env:"PG_URL"`
	}

	// RMQ -.
	RMQ struct {
		ServerExchange string `env-required:"true" yaml:"rpc_server_exchange" env:"RMQ_RPC_SERVER"`
		ClientExchange string `env-required:"true" yaml:"rpc_client_exchange" env:"RMQ_RPC_CLIENT"`
		URL            string `env-required:"false" yaml:"rpc_url" env:"RMQ_URL"`
	}

	// Redis -.
	Redis struct {
		Addr string `env-required:"true"  yaml:"redis_addr" env:"REDIS_ADDR"`
	}

	// Mongo -.
	Mongo struct {
		URI string `env-required:"true" yaml:"mongo_uri" env:"MONGO_URI"`
	}

	PARAM struct {
		DiasbleSwaggerHttpHandler string `env-required:"true" yaml:"disable_swagger_http_handler" env:"DISABLE_SWAGGER_HTTP_HANDLER"`
		GinMode                   string `env-required:"true" yaml:"gin_mode" env:"GIN_MODE"`
		TgBotApi                  string `env-required:"false" yaml:"tg_bot_api" env:"TG_BOT_API"`
		TgChatId                  string `env-required:"false" yaml:"tg_chat_id" env:"TG_CHAT_ID"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	// all configuration
	//../../config/config.yml
	err := cleanenv.ReadConfig("./config/config.yml", cfg) //./config/config.yml 4 migrate
	if err != nil {
		return nil, fmt.Errorf("yml config error: %w", err)
	}

	/* doesn't read from env in project root*/
	/* err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("env config error: %w", err)
	}
	*/

	return cfg, nil
}
