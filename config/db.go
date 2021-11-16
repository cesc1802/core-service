package config

type Database struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	Dbname          string `yaml:"dbname"`
	MaxIdleConns    int    `yaml:"maxIdleConns"`
	MaxOpenConns    int    `yaml:"maxOpenConns"`
	ConnMaxLifetime string `yaml:"connMaxLifetime"`
	//ConnTimeout     time.Duration `yaml:"connTimeout"`
	//ReadTimeout     time.Duration `yaml:"readTimeout"`
	//WriteTimeout    time.Duration `yaml:"writeTimeout"`
}
