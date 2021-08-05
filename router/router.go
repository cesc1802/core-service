package router

import (
	"github.com/cesc1802/core-service/config"
	"github.com/cesc1802/core-service/i18n"
	"github.com/cesc1802/core-service/middleware"
	appValidator "github.com/cesc1802/core-service/validator"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Router struct {
	Engine *gin.Engine
	I18n   *i18n.I18n
}

func NewRouter(c config.Config, i18n *i18n.I18n) (*Router, error) {
	e := gin.New()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(appValidator.JsonTagNameFunc)
	}
	e.RedirectTrailingSlash = true
	e.RedirectFixedPath = true

	e.Use(middleware.Recovery(i18n))

	// CORS
	e.Use(middleware.Cors(c))

	return &Router{
		Engine: e,
		I18n:   i18n,
	}, nil
}
