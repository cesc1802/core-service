package gormdialects

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func MySqlDB(uri string) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(uri), &gorm.Config{})
}

func PostgresDB(uri string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(uri), &gorm.Config{})
}

func MssqlDB(uri string) (*gorm.DB, error) {
	return gorm.Open(sqlserver.Open(uri), &gorm.Config{})
}

func SqliteDB(uri string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(uri), &gorm.Config{})
}
