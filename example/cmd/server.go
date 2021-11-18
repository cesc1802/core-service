package cmd

import (
	"example/cmd/handlers"
	"fmt"
	"github.com/cesc1802/core-service/config"
	"github.com/cesc1802/core-service/httpserver"
	"github.com/cesc1802/core-service/i18n"
	logger2 "github.com/cesc1802/core-service/logger"
	"github.com/cesc1802/core-service/plugin/storage/sdkgorm"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use: "server",
	RunE: func(cmd *cobra.Command, args []string) error {
		coreCfg, _ := config.LoadConfig()
		i18n, _ := i18n.NewI18n(coreCfg.I18nConfig)
		logger := logger2.NewLogger(coreCfg.LogConfig)
		gormdb := sdkgorm.NewGormDB("portal", "portal", &coreCfg.DatabaseConfig)
		gormdb1 := sdkgorm.NewGormDB("demo-portal", "demo-portal", &coreCfg.DatabaseConfig)
		ginService, _ := httpserver.New(coreCfg, i18n, logger)

		app := NewAppService(WithName("demo"),
			WithVersion("1.0.0"),
			WithHttpServer(ginService),
			WithInitRunnable(gormdb),
			WithInitRunnable(gormdb1),
		)

		ginService.AddHandler(handlers.PublicHandler(app))
		ginService.AddHandler(handlers.PrivateHandler(app))

		fmt.Println("This is a service name", app.Name(), app.Version())
		if err := app.Run(); err != nil {
			panic(err)
		}
		return nil
	},
}
