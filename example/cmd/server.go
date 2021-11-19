package cmd

import (
	"example/cmd/handlers"
	"example/consumer"
	"fmt"
	"github.com/cesc1802/core-service/config"
	"github.com/cesc1802/core-service/httpserver"
	"github.com/cesc1802/core-service/i18n"
	"github.com/cesc1802/core-service/logger"
	"github.com/cesc1802/core-service/plugin/pubsub"
	"github.com/cesc1802/core-service/plugin/storage/sdkgorm"
	"github.com/spf13/cobra"
	"log"
)

var serverCmd = &cobra.Command{
	Use: "server",
	RunE: func(cmd *cobra.Command, args []string) error {
		coreCfg, _ := config.LoadConfig()
		i18n, _ := i18n.NewI18n(coreCfg.I18nConfig)
		baseLogger := logger.New(coreCfg.LogConfig.Level)

		app := NewAppService(
			WithName("demo"),
			WithVersion("1.0.0"),
			WithHttpServer(httpserver.New(coreCfg, i18n, *baseLogger)),
			WithInitRunnable(sdkgorm.NewGormDB("portal", "portal", &coreCfg.DatabaseConfig)),
			WithInitRunnable(sdkgorm.NewGormDB("demo-portal", "demo-portal", &coreCfg.DatabaseConfig)),
			WithInitRunnable(pubsub.New("pubsub", "pubsub")),
		)

		consumer.NewEngine(app).Start()
		app.httpserver.AddHandler(handlers.PrivateHandler(app))
		app.httpserver.AddHandler(handlers.PublicHandler(app))

		fmt.Println("This is a service name", app.Name(), app.Version())
		if err := app.Run(); err != nil {
			panic(err)
		} else {
			log.Println("application start")
		}
		return nil
	},
}
