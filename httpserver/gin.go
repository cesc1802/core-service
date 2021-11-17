package httpserver

import (
	"context"
	"fmt"
	"github.com/cesc1802/core-service/config"
	"github.com/cesc1802/core-service/httpserver/middleware"
	"github.com/cesc1802/core-service/i18n"
	baseValidator "github.com/cesc1802/core-service/validator"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"net"
	"net/http"
	"time"
)

type GinOpt struct {
	name string
	port string
	host string
}

type GinService struct {
	Engine        *gin.Engine
	graceFullServ *http.Server
	i18n          *i18n.I18n
	logger        *zerolog.Logger
	*GinOpt
}

func NewGinService(c config.Config, i18n *i18n.I18n, logger *zerolog.Logger) (*GinService, error) {
	return &GinService{
		i18n:   i18n,
		logger: logger,
		GinOpt: &GinOpt{
			name: "GIN-Service",
			port: c.ServerConfig.Port,
			host: c.ServerConfig.Host,
		},
	}, nil
}

func (r *GinService) Configure() error {
	r.Engine = gin.New()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(baseValidator.JsonTagNameFunc)
	}
	r.Engine.RedirectTrailingSlash = true
	r.Engine.RedirectFixedPath = true

	// Recovery
	r.Engine.Use(middleware.Recovery(r.i18n))

	//TODO: you can add more middleware here

	return nil

}

func (r *GinService) Name() string {
	return ""
}

func (r *GinService) Start() error {
	if err := r.Configure(); err != nil {
		return err
	}

	r.graceFullServ = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", r.host, r.port),
		Handler: r.Engine,
	}
	r.logger.Info().Msgf("Listening and serving HTTP on %v:%v", r.host, r.port)

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", r.host, r.port))
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

func (r *GinService) Stop() error {
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFn()

	if r.graceFullServ != nil {
		r.logger.Info().Msg("shutting down....")
		_ = r.graceFullServ.Shutdown(ctx)
	}
	return nil
}
