package sdkgorm

import (
	"errors"
	"github.com/cesc1802/core-service/config"
	"github.com/cesc1802/core-service/plugin/storage/sdkgorm/gormdialects"
	"gorm.io/gorm"
	"log"
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

type GormOpt struct {
	Uri          string
	Prefix       string
	DBType       string
	PingInterval int //unit in second
}

type gormDB struct {
	prefix    string
	name      string
	db        *gorm.DB
	isRunning bool
	once      *sync.Once
	dbCfg     *config.DatabaseConfig
	*GormOpt
}

func NewGormDB(name, prefix string, cfg *config.DatabaseConfig) *gormDB {
	return &gormDB{
		name:      name,
		prefix:    prefix,
		isRunning: false,
		once:      new(sync.Once),
		dbCfg:     cfg,
	}
}

func (gdb *gormDB) GetPrefix() string {
	return gdb.Prefix
}

func (gdb *gormDB) Name() string {
	return gdb.name
}

func (gdb *gormDB) isDisable() bool {
	return gdb.Uri == ""
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
		return gormdialects.MssqlDB(gdb.Uri)
	case GormDBTypeSQLite:
		return gormdialects.SqliteDB(gdb.Uri)
	case GormDBTypePostgres:
		return gormdialects.PostgresDB(gdb.Uri)
	case GormDBTypeMySQL:
		return gormdialects.MySqlDB(gdb.Uri)
	}
	return nil, nil
}

func (gdb *gormDB) reconnectIfNeed() {
	for {
		conn, err := gdb.db.DB()
		if err = conn.Ping(); err != nil {
			_ = conn.Close()
			gdb.isRunning = false
			gdb.once = new(sync.Once)
			_ = gdb.Get()
		}
		time.Sleep(time.Second * time.Duration(gdb.PingInterval))
	}
}

func (gdb *gormDB) Get() interface{} {
	gdb.once.Do(func() {
		if !gdb.isRunning && gdb.isDisable() {
			if db, err := gdb.getConnWithRetry(getDBType(gdb.DBType), math.MaxInt32); err != nil {
				gdb.db = db
				gdb.isRunning = true
			} else {
				log.Printf("%v", err)
				log.Fatalf("%s connection cannot reconnect\n", gdb.name)
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
	if gdb.isDisable() || gdb.isRunning {
		return nil
	}

	dbType := getDBType(gdb.DBType)
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

func (gdb *gormDB) Run() error {
	if err := gdb.Configure(); err != nil {
		return nil
	}
	return nil
}
