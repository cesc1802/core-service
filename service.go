package core_service

import "github.com/gin-gonic/gin"

type HasPrefix interface {
	Get() interface{}
	GetPrefix() string
}

type Storage interface {
	Get(prefix string) (interface{}, bool)
	MustGet(prefix string) interface{}
}

type Runnable interface {
	Name() string
	Configure() error
	Start() error
	Stop() error
}
type PrefixRunnable interface {
	HasPrefix
	Runnable
}

type HTTPServerHandler = func(*gin.Engine)
type HttpServer interface {
	Runnable
	AddHandler(HTTPServerHandler)
}

type ServiceContext interface {
	Get(prefix string) (interface{}, bool)
	MustGet(prefix string) interface{}
}

type Service interface {
	ServiceContext
	Version() string
	HttpServer() HttpServer
	Run() error
	Stop() error
}
