package sdkmgo

import (
	"context"
	"fmt"
	"github.com/cesc1802/core-service/config"
	"github.com/cesc1802/core-service/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
	"sync"
	"time"
)

const retryCount = 10

type MgoOpt struct {
	name   string
	prefix string
	host   string
	port   int
}

type mongoDB struct {
	mgo       *mongo.Client
	isRunning bool
	logger    logger.Interface
	once      *sync.Once
	*MgoOpt
}

func NewMongoDB(name, prefix string, cfg config.NoSQLConfig) *mongoDB {
	return &mongoDB{
		MgoOpt: &MgoOpt{
			name:   name,
			prefix: prefix,
			host:   cfg.Host,
			port:   cfg.Port,
		},
		once:      new(sync.Once),
		isRunning: false,
		mgo:       nil,
	}
}

func (mdb *mongoDB) GetPrefix() string {
	return mdb.prefix
}

func (mdb *mongoDB) Name() string {
	return mdb.name
}

func (mdb *mongoDB) reconnectIfNeed() {
	for {
		if err := mdb.mgo.Ping(context.TODO(), nil); err != nil {
			_ = mdb.mgo.Disconnect(context.TODO())
			mdb.logger.Info("connect is gone, try to connect %s\n", mdb.name)
			mdb.isRunning = false
			mdb.once = new(sync.Once)
			_ = mdb.Get()
			return
		}
		time.Sleep(time.Second * time.Duration(5))
	}
}

func (mdb *mongoDB) Get() interface{} {
	mdb.once.Do(func() {
		if !mdb.isRunning {
			if mgo, err := mdb.getConnWithRetry(math.MaxInt32); err == nil {
				mdb.mgo = mgo
				mdb.isRunning = true
			} else {
				mdb.logger.Fatal("connection cannot reconnect\n", mdb.name, err)
			}
		}
	})

	if mdb.mgo == nil {
		return nil
	}

	//TODO: need setup logger
	//mdb.db.Logger =

	return mdb.mgo
}

func (mdb mongoDB) getDBConn() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d", mdb.host, mdb.port))
	mgoDB, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return nil, err
	}

	return mgoDB, nil
}

func (mdb *mongoDB) getConnWithRetry(retry int) (*mongo.Client, error) {

	mgoDB, err := mdb.getDBConn()

	if err != nil {
		for {
			time.Sleep(time.Second * 1)
			mgoDB, err = mdb.getDBConn()

			if err == nil {
				go mdb.reconnectIfNeed()
				break
			}
		}

	} else {
		go mdb.reconnectIfNeed()
	}
	return mgoDB, err
}

func (mdb *mongoDB) Configure() error {
	if mdb.isRunning {
		return nil
	}

	var err error
	mdb.mgo, err = mdb.getConnWithRetry(retryCount)
	if err != nil {
		return nil
	}

	mdb.isRunning = true
	return nil
}

func (mdb *mongoDB) Start() error {
	if err := mdb.Configure(); err != nil {
		return nil
	}
	return nil
}

func (mdb *mongoDB) Stop() error {
	if mdb.mgo != nil {
		if err := mdb.mgo.Disconnect(nil); err != nil {
			mdb.logger.Info("disconnect to server error", err)
		}
	}
	mdb.logger.Info("connection to mongodb closed")
	return nil
}
