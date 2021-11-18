package cmd

import (
	"context"
	"errors"
	sdkservice "github.com/cesc1802/core-service"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
)

type AppServiceOption func(service *AppService)
type AppService struct {
	name         string
	version      string
	ctx          context.Context
	cancel       func()
	signals      []os.Signal
	httpserver   sdkservice.HttpServer
	subServices  []sdkservice.Runnable
	initServices map[string]sdkservice.PrefixRunnable
}

func WithVersion(version string) AppServiceOption {
	return func(app *AppService) {
		app.version = version
	}
}

func WithName(name string) AppServiceOption {
	return func(app *AppService) {
		app.name = name
	}
}

func WithHttpServer(server sdkservice.HttpServer) AppServiceOption {
	return func(app *AppService) {
		app.httpserver = server
	}
}

func WithInitRunnable(runnable sdkservice.PrefixRunnable) AppServiceOption {
	return func(app *AppService) {
		app.subServices = append(app.subServices, runnable)
		app.initServices[runnable.GetPrefix()] = runnable
	}
}

func NewAppService(opts ...AppServiceOption) *AppService {
	sv := &AppService{
		signals:      []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		subServices:  []sdkservice.Runnable{},
		initServices: map[string]sdkservice.PrefixRunnable{},
	}

	if sv.ctx == nil {
		sv.ctx = context.Background()
	}

	// Setup cancelCtx to listen Ctrl + C
	cancelCtx, cancelFunc := context.WithCancel(sv.ctx)
	sv.ctx = cancelCtx
	sv.cancel = cancelFunc

	for _, opt := range opts {
		opt(sv)
	}
	sv.subServices = append(sv.subServices, sv.httpserver)

	return sv
}

func (s *AppService) Run() error {
	g, ctx := errgroup.WithContext(s.ctx)

	for _, s := range s.subServices {
		srv := s
		g.Go(func() error {
			<-ctx.Done()
			return srv.Stop()
		})
		g.Go(func() error {
			return srv.Start()
		})
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, s.signals...)
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-c:
				_ = s.Stop()
			}
		}
	})
	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

func (s *AppService) Stop() error {
	if s.cancel != nil {
		s.cancel()
	}
	return nil
}

func (s *AppService) HttpServer() sdkservice.HttpServer {
	return s.httpserver
}

func (s *AppService) Get(prefix string) (interface{}, bool) {
	if prefixRunnable, ok := s.initServices[prefix]; ok {
		return prefixRunnable.Get(), true
	}
	return nil, false
}

func (s *AppService) MustGet(prefix string) interface{} {
	return s.initServices[prefix].Get()
}

func (s *AppService) Name() string {
	return s.name
}

func (s *AppService) Version() string {
	return s.version
}
