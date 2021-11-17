package main

import (
	"fmt"
	"github.com/cesc1802/core-service"
	"github.com/cesc1802/core-service/config"
	"github.com/cesc1802/core-service/httpserver"
	"github.com/cesc1802/core-service/i18n"
	logger2 "github.com/cesc1802/core-service/logger"
	"github.com/cesc1802/core-service/plugin/storage/sdkgorm"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//type App struct {
//	config config.Config
//	router *router.Router
//	i18n   *i18n.I18n
//	//database *database.Database
//}

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type App struct {
	dbconn map[string]interface{}
}

func handlers(app *App) func(engine *gin.Engine) {
	return func(engine *gin.Engine) {

		user := engine.Group("/users")
		{
			user.POST("", func(c *gin.Context) {
				fmt.Println(app.dbconn)
				conn := app.dbconn["portal"].(*gorm.DB)

				var u User
				conn.Table("users").Find(&u)
				//if err := c.ShouldBind(&u); err != nil {
				//	panic(err)
				//	return
				//}
				c.JSON(200, u)
			})
		}
	}
}

type HasPrefix interface {
	Get() interface{}
	GetPrefix() string
}

type Storage interface {
	Get(prefix string) (interface{}, bool)
	MustGet(prefix string) interface{}
}

func InitApp(dbconn map[string]interface{}) *App {
	return &App{
		dbconn: dbconn,
	}
}

func main() {
	coreCfg, _ := config.LoadConfig()
	i18n, _ := i18n.NewI18n(coreCfg.I18nConfig)
	logger := logger2.Create(coreCfg.LogConfig)
	gormdb := sdkgorm.NewGormDB("portal", "portal", &coreCfg.DatabaseConfig)
	ginService, _ := httpserver.NewGinService(coreCfg, i18n, logger)

	appGroup := core_service.NewAppGroup(
		core_service.AppGroupOption{
			Name:     "Demo",
			Services: []core_service.Service{ginService, gormdb},
		})

	if err := appGroup.Run(); err != nil {
		panic(err)
	}
}
