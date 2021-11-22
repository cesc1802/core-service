package config

import (
	"fmt"
	"github.com/cesc1802/core-service/constant"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Env string

	RedisConfig      `mapstructure:"redis"`
	SQLDBConfigs     `mapstructure:"databases"`
	NoSQLConfigs     `mapstructure:"nosqldatabases"`
	ServerConfig     `mapstructure:"server"`
	HttpClientConfig `mapstructure:"client"`
	I18nConfig       `mapstructure:"i18n"`
	CORSConfig       `mapstructure:"cors"`
	LogConfig        `mapstructure:"log"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig() (c Config, err error) {
	env := extractEnv()
	defer func() {
		if err != nil {
			c = Config{
				Env: env,
			}
		}
	}()
	// get current path
	pwd, err := os.Getwd()
	if err != nil {
		return
	}

	path, err := filepath.Abs(pwd)
	if err != nil {
		return
	}

	// load config from config directory
	if path == "/" {
		viper.AddConfigPath("/config")
	} else {
		viper.AddConfigPath(fmt.Sprintf("%v/config", path))
	}
	viper.SetConfigName(fmt.Sprintf("app-%v", strings.ToLower(env)))
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}
	if err = viper.Unmarshal(&c); err != nil {
		return
	}
	return c, nil
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

func getAbsPath(dir string) (string, error) {
	path, err := filepath.Abs(dir)
	if err != nil {
		panic(err)
	}
	return path, nil
}
