package config

type NoSQLConfig struct {
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	DBType string `yaml:"dbtype"`
}

type NoSQLConfigs []NoSQLConfig
