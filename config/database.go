package config

type DatabaseConfig struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	DBName          string `yaml:"dbname"`
	DBType          string `yaml:"dbtype"`
	MaxIdleConns    int    `yaml:"maxIdleConns"`
	MaxOpenConns    int    `yaml:"maxOpenConns"`
	ConnMaxLifetime string `yaml:"connMaxLifetime"`
	//ConnTimeout     time.Duration `yaml:"connTimeout"`
	//ReadTimeout     time.Duration `yaml:"readTimeout"`
	//WriteTimeout    time.Duration `yaml:"writeTimeout"`
}
