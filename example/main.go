package main

import (
	core_service "github.com/cesc1802/core-service"
	"github.com/cesc1802/core-service/config"
	"github.com/cesc1802/core-service/database"
	"github.com/cesc1802/core-service/i18n"
	logger2 "github.com/cesc1802/core-service/logger"
	"github.com/cesc1802/core-service/router"
	"github.com/gin-gonic/gin"
)

type App struct {
	config   config.Config
	router   *router.Router
	i18n     *i18n.I18n
	database *database.Database
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
	configConfig, _ := config.LoadConfig()
	i18n, _ := i18n.NewI18n(configConfig)
	logger := logger2.Create(configConfig)
	router, _ := router.NewRouter(configConfig, i18n, logger)
	router.AddHandle(handlers())

	//database, _ := database.NewDatabase(configConfig)

	appGroup := core_service.NewAppGroup(
		core_service.AppGroupOption{
			Name:     "Gin Service",
			Services: []core_service.Service{router},
		})

	if err := appGroup.Run(); err != nil {
		logger.Fatal().Msg("stoping execute")
	}

	logger.Info().Msg("server started.")
}
