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

	Database `yaml:"database"`

	Redis struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Password string `yaml:"password"`
		Database int    `yaml:"database"`
		Ttl      struct {
			Common string `yaml:"common"`
		} `mapstructure:"ttls"`
	} `yaml:"redis"`

	Server `yaml:"server"`

	//Jwt struct {
	//	Realm              string `yaml:"realm"`
	//	SigningAlg         string `yaml:"signAlg"`
	//	Secret             string `yaml:"secret"`
	//	ExpiredTime        string `yaml:"expiredTime"`
	//	RefreshExpTime     string `yaml:"refreshExpTime"`
	//	LongRefreshExpTime string `yaml:"longRefreshExpTime"`
	//} `yaml:"jwt"`

	Resty struct {
		Debug   bool   `yaml:"debug"`
		Timeout string `yaml:"timeout"`
	} `yaml:"resty"`

	I18n struct {
		Langs []string `yaml:"langs"`
	} `yaml:"i18n"`

	CORS `yaml:"cors"`

	//HostUrl map[string]string `yaml:"hostUrl"`
	//
	//Security struct {
	//	AuthorizedRequests []ConfigAuthorizedRequests `yaml:"authorizedRequests"`
	//} `yaml:"security"`

	Log struct {
		Level                 string `yaml:"level"`
		ConsoleLoggingEnabled bool   `yaml:"consoleLoggingEnabled"`
		EncodeLogsAsJson      bool   `yaml:"encodeLogsAsJson"`
		FileLoggingEnabled    bool   `yaml:"fileLoggingEnabled"`
		Directory             string `yaml:"directory"`
		Filename              string `yaml:"filename"`
		MaxSize               int    `yaml:"maxSize"`
		MaxBackups            int    `yaml:"maxBackups"`
		MaxAge                int    `yaml:"maxAge"`
	} `yaml:"log"`
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
