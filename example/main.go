package main

import (
	core_service "github.com/cesc1802/core-service"
	"github.com/cesc1802/core-service/config"
	"github.com/cesc1802/core-service/database/mysql"
	"github.com/cesc1802/core-service/i18n"
	logger2 "github.com/cesc1802/core-service/logger"
	"github.com/cesc1802/core-service/router"
	"github.com/gin-gonic/gin"
)

type App struct {
	config   config.Config
	router   *router.Router
	i18n     *i18n.I18n
	database *mysql.Database
}

type User struct {
	//Id       int    `json:"id"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

//type User struct {
//	Id      int    `json:"id" gorm:"id"`
//	LoginId string `json:"loginId" gorm:"login_id"`
//}

func handlers(database *mysql.Database) func(engine *gin.Engine) {
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

			user.GET("", func(c *gin.Context) {
				db := database.MustGet()
				var u User
				query := db.Table("users").Where("status = 1")
				if err := query.Find(&u).Error; err != nil {
					panic(err)
				}

				c.JSON(200, gin.H{"user": u})
			})
		}
	}
}

func main() {
	configConfig, _ := config.LoadConfig()
	i18n, _ := i18n.NewI18n(configConfig)
	logger := logger2.Create(configConfig)
	router, _ := router.NewRouter(configConfig, i18n, logger)

	database, err := mysql.NewDatabase(configConfig)

	if err != nil {
		logger.Info().Msgf("cannot init db %v", err)
	}

	appGroup := core_service.NewAppGroup(
		core_service.AppGroupOption{
			Name:     "Gin Service",
			Services: []core_service.Service{router, database},
		})
	router.AddHandle(handlers(database))
	if err := appGroup.Run(); err != nil {
		panic(err)
	}
}
