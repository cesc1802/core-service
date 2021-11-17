package main

import (
	"fmt"
	"github.com/cesc1802/core-service/config"
	"github.com/cesc1802/core-service/i18n"
	"github.com/cesc1802/core-service/router"
	"github.com/gin-gonic/gin"
)

type App struct {
	config config.Config
	router *router.Router
	i18n   *i18n.I18n
	//database *database.Database
}

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func handlers() func(engine *gin.Engine) {
	return func(engine *gin.Engine) {
		user := engine.Group("/users")
		{
			user.POST("", func(c *gin.Context) {
				var u User
				if err := c.ShouldBind(&u); err != nil {
					panic(err)
					return
				}

			})
		}
	}
}

//func (r App) Start() error {
//
//	r.router.Engine.Run(fmt.Sprintf(":%s", r.config.Server.Port))
//
//	gracefulShutdown(&http.Server{
//		Addr:    fmt.Sprintf(":%s", r.config.Server.Port),
//		Handler: r.router.Engine,
//	})
//
//	return nil
//}
//
//func (r App) Stop() {
//	if err := r.database.Close(); err != nil {
//		panic(err)
//	}
//
//}

func main() {
	coreCfg, _ := config.LoadConfig()
	fmt.Println(coreCfg)
	//i18n, _ := i18n.NewI18n(coreCfg.I18nConfig)
	//logger := logger2.Create(coreCfg.LogConfig)
	//router, _ := router.NewRouter(coreCfg, i18n, logger)
	////router.AddHandle(handlers())
	//
	//appGroup := core_service.NewAppGroup(
	//	core_service.AppGroupOption{
	//		Name:     "Gin Service",
	//		Services: []core_service.Service{router},
	//	})
	//
	//if err := appGroup.Run(); err != nil {
	//	panic(err)
	//}
}
