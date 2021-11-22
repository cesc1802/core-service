package gormdialects

import (
	"fmt"
	"github.com/cesc1802/core-service/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"strconv"
)

func MySqlDB(cfg *config.SQLDBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username, cfg.Password, cfg.Host, strconv.Itoa(cfg.Port), cfg.DBName)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func PostgresDB(cfg *config.SQLDBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username, cfg.Password, cfg.Host, strconv.Itoa(cfg.Port), cfg.DBName)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func MssqlDB(cfg *config.SQLDBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username, cfg.Password, cfg.Host, strconv.Itoa(cfg.Port), cfg.DBName)
	return gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
}

func SqliteDB(cfg *config.SQLDBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username, cfg.Password, cfg.Host, strconv.Itoa(cfg.Port), cfg.DBName)
	return gorm.Open(sqlite.Open(dsn), &gorm.Config{})
}
