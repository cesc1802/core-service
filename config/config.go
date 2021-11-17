package config

import (
	"fmt"
	"github.com/cesc1802/core-service/constant"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type Config struct {
	Env string

	RedisConfig      `mapstructure:"redis"`
	DatabaseConfig   `mapstructure:"database"`
	ServerConfig     `mapstructure:"server"`
	HttpClientConfig `mapstructure:"client"`
	I18nConfig       `mapstructure:"i18n"`
	CORSConfig       `mapstructure:"cors"`
	LogConfig        `mapstructure:"log"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig() (Config, error) {
	viper.AddConfigPath("./config")
	env := extractEnv()
	viper.SetConfigName(fmt.Sprintf("app-%v", strings.ToLower(env)))
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return Config{
			Env: env,
		}, err
	}

	config := Config{
		Env: env,
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return Config{
			Env: env,
		}, err
	}
	return config, nil
}

func extractEnv() string {
	env := os.Getenv("ENVIRONMENT")
	if len(env) == 0 {
		env = os.Getenv("ENV")
	}
	if len(env) == 0 {
		env = constant.DefaultEnv
	}
	return env
}
