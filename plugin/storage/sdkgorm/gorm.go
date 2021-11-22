package sdkgorm

import (
	"errors"
	"github.com/cesc1802/core-service/config"
	"github.com/cesc1802/core-service/logger"
	"github.com/cesc1802/core-service/plugin/storage/sdkgorm/gormdialects"
	"gorm.io/gorm"
	"math"
	"strings"
	"sync"
	"time"
)

type GormDBType int

const (
	GormDBTypeMySQL GormDBType = iota + 1
	GormDBTypePostgres
	GormDBTypeSQLite
	GormDBTypeMsSQL
	GormDBTypeNotSupported
)

const (
	retryCount = 10
)

type gormDB struct {
	prefix    string
	name      string
	db        *gorm.DB
	isRunning bool
	once      *sync.Once
	logger    logger.Interface
	cfg       *config.SQLDBConfig
}

func NewGormDB(name, prefix string, cfg *config.SQLDBConfig) *gormDB {
	return &gormDB{
		name:      name,
		prefix:    prefix,
		isRunning: false,
		once:      new(sync.Once),
		cfg:       cfg,
	}
}

func (gdb *gormDB) GetPrefix() string {
	return gdb.prefix
}

func (gdb *gormDB) Name() string {
	return gdb.name
}

func getDBType(dbType string) GormDBType {
	switch strings.ToLower(dbType) {
	case "mysql":
		return GormDBTypeMySQL
	case "postgres":
		return GormDBTypePostgres
	case "mssql":
		return GormDBTypeMsSQL
	case "sqlite":
		return GormDBTypeSQLite
	default:
		return GormDBTypeNotSupported
	}
}

func (gdb *gormDB) getDBConn(t GormDBType) (*gorm.DB, error) {
	switch t {
	case GormDBTypeMsSQL:
		return gormdialects.MssqlDB(gdb.cfg)
	case GormDBTypeSQLite:
		return gormdialects.SqliteDB(gdb.cfg)
	case GormDBTypePostgres:
		return gormdialects.PostgresDB(gdb.cfg)
	case GormDBTypeMySQL:
		return gormdialects.MySqlDB(gdb.cfg)
	}
	return nil, nil
}

func (gdb *gormDB) reconnectIfNeed() {
	for {
		conn, err := gdb.db.DB()
		if err = conn.Ping(); err != nil {
			_ = conn.Close()
			gdb.logger.Info("connect is gone, try to connect %s\n", gdb.name)
			gdb.isRunning = false
			gdb.once = new(sync.Once)
			_ = gdb.Get()
			return
		}
		time.Sleep(time.Second * time.Duration(5))
	}
}

func (gdb *gormDB) Get() interface{} {
	gdb.once.Do(func() {
		if !gdb.isRunning {
			if db, err := gdb.getConnWithRetry(getDBType(gdb.cfg.DBType), math.MaxInt32); err == nil {
				gdb.db = db
				gdb.isRunning = true
			} else {
				gdb.logger.Fatal("connection cannot reconnect\n", gdb.name, err)
			}
		}
	})

	if gdb.db == nil {
		return nil
	}

	//TODO: need setup logger
	//gdb.db.Logger =

	return gdb.db
}
func (gdb *gormDB) getConnWithRetry(dbType GormDBType, retry int) (*gorm.DB, error) {

	db, err := gdb.getDBConn(dbType)

	if err != nil {
		for {
			time.Sleep(time.Second * 1)
			db, err = gdb.getDBConn(dbType)

			if err == nil {
				go gdb.reconnectIfNeed()
				break
			}
		}

	} else {
		go gdb.reconnectIfNeed()
	}
	return db, err
}
func (gdb *gormDB) Configure() error {
	if gdb.isRunning {
		return nil
	}

	dbType := getDBType(gdb.cfg.DBType)
	if dbType == GormDBTypeNotSupported {
		return errors.New("gorm database type is not supported")
	}

	var err error
	gdb.db, err = gdb.getConnWithRetry(dbType, retryCount)
	if err != nil {
		return nil
	}

	gdb.isRunning = true
	return nil
}

func (gdb *gormDB) Start() error {
	if err := gdb.Configure(); err != nil {
		return nil
	}
	return nil
}

func (gdb *gormDB) Stop() error {
	if gdb.db != nil {
		if conn, err := gdb.db.DB(); err != nil {
			conn.Close()
		}
	}
	return nil
}
