package database

//
//import (
//	"fmt"
//	"github.com/cesc1802/core-service/config"
//	"gorm.io/driver/mysql"
//	"gorm.io/gorm"
//	"gorm.io/gorm/logger"
//	"log"
//	"os"
//	"time"
//)
//
//type Database struct {
//	DB *gorm.DB
//}
//
//func NewDatabase(c config.DatabaseConfig) (*Database, error) {
//	db, err := setupDatabase(c)
//	if err != nil {
//		return nil, err
//	}
//	return &Database{DB: db}, nil
//}
//
//func setupDatabase(c config.DatabaseConfig) (*gorm.DB, error) {
//
//	//&timeout=%s&readTimeout=%s&writeTimeout=%s
//
//	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local&charset=%s",
//		c.Username,
//		c.Password,
//		c.Host,
//		c.Port,
//		c.Dbname,
//		//c.Database.ConnTimeout,
//		//c.Database.ReadTimeout,
//		//c.Database.WriteTimeout,
//		"utf8mb4",
//	)
//	var err error
//	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
//		SkipDefaultTransaction: true,
//		Logger: logger.New(
//			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
//			logger.Config{
//				SlowThreshold: time.Second, // Slow SQL threshold
//				LogLevel:      logger.Info, // Log level
//				Colorful:      true,        // Enable color
//			},
//		),
//	})
//	if err != nil {
//		return nil, err
//	}
//
//	sqlDB, err := db.DB()
//	if err != nil {
//		return nil, err
//	}
//
//	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
//	sqlDB.SetMaxIdleConns(c.MaxIdleConns)
//
//	// SetMaxOpenConns sets the maximum number of open connections to the database.
//	sqlDB.SetMaxOpenConns(c.MaxOpenConns)
//
//	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
//	duration, _ := time.ParseDuration(c.ConnMaxLifetime)
//	sqlDB.SetConnMaxLifetime(duration)
//
//	return db, nil
//}
//
//func (r *Database) Close() error {
//	sqlDB, err := r.DB.DB()
//	if err != nil {
//		return err
//	}
//	if err = sqlDB.Close(); err != nil {
//		return err
//	}
//	return nil
//}
