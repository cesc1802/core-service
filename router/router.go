package router

import (
	"context"
	"errors"
	"fmt"
	"github.com/cesc1802/core-service/config"
	"github.com/cesc1802/core-service/i18n"
	"github.com/cesc1802/core-service/middleware"
	appValidator "github.com/cesc1802/core-service/validator"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"net"
	"net/http"
	"time"
)

type Router struct {
	serverCfg     *config.Server
	corsCfg       *config.CORS
	Engine        *gin.Engine
	handlers      []func(*gin.Engine)
	graceFullServ *http.Server
	I18n          *i18n.I18n
	logger        *zerolog.Logger
}

func NewRouter(c config.Config, i18n *i18n.I18n, logger *zerolog.Logger) (*Router, error) {
	return &Router{
		I18n:      i18n,
		serverCfg: &c.Server,
		corsCfg:   &c.CORS,
		logger:    logger,
		handlers:  []func(*gin.Engine){},
	}, nil
}

func (r *Router) AddHandle(hdl func(*gin.Engine)) {
	r.handlers = append(r.handlers, hdl)
}

func (r *Router) Configure() error {
	r.Engine = gin.New()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(appValidator.JsonTagNameFunc)
	}
	r.Engine.RedirectTrailingSlash = true
	r.Engine.RedirectFixedPath = true

	// Recovery
	r.Engine.Use(middleware.Recovery(r.I18n))

	// CORS
	r.Engine.Use(middleware.Cors(*r.corsCfg))

	//TODO: you can add more configure here

	return nil

}

func (r *Router) Start() error {
	if err := r.Configure(); err != nil {
		return err
	}

	for _, hdl := range r.handlers {
		hdl(r.Engine)
	}

	r.graceFullServ = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", r.serverCfg.Host, r.serverCfg.Port),
		Handler: r.Engine,
	}
	r.logger.Info().Msgf("Listening and serving HTTP on %v:%v", r.serverCfg.Host, r.serverCfg.Port)

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", r.serverCfg.Host, r.serverCfg.Port))
	if err != nil {
		r.logger.Info().Msgf("Listening error: %v", err)
		return err
	}

	err = r.graceFullServ.Serve(lis)

	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (r *Router) Stop() error {
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFn()

	if r.graceFullServ != nil {
		r.logger.Info().Msg("shutting down....")
		_ = r.graceFullServ.Shutdown(ctx)
	}
	return nil
}
