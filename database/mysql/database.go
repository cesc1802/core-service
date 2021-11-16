package mysql

import (
	"fmt"
	"github.com/cesc1802/core-service/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type Database struct {
	config *config.Database
	DB     *gorm.DB
}

func NewDatabase(c config.Config) (*Database, error) {
	return &Database{
		config: &c.Database,
		DB:     nil,
	}, nil
}

func (d *Database) MustGet() *gorm.DB {
	return d.DB
}

func (d *Database) Configure() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local&charset=%s",
		d.config.Username,
		d.config.Password,
		d.config.Host,
		d.config.Port,
		d.config.Dbname,
		//d.config.ConnTimeout,
		//d.config.ReadTimeout,
		//d.config.WriteTimeout,
		"utf8mb4",
	)
	var err error
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Second, // Slow SQL threshold
				LogLevel:      logger.Info, // Log level
				Colorful:      true,        // Enable color
			},
		),
	})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(d.config.MaxIdleConns)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(d.config.MaxOpenConns)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	duration, _ := time.ParseDuration(d.config.ConnMaxLifetime)
	sqlDB.SetConnMaxLifetime(duration)

	d.DB = db
	return nil
}

func (d *Database) Start() error {
	if err := d.Configure(); err != nil {
		return err
	}

	sqlDB, err := d.DB.DB()

	if err != nil {
		return err
	}

	if err = sqlDB.Ping(); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
func (d *Database) Stop() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	if err = sqlDB.Close(); err != nil {
		return err
	}
	return nil
}
