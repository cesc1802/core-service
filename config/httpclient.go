package config

type HttpClientConfig struct {
	Debug   bool   `yaml:"debug"`
	Timeout string `yaml:"timeout"`
}
