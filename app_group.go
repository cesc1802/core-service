package core_service

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
)

type AppGroup struct {
	option AppGroupOption
	ctx    context.Context
	cancel func()
}

type AppGroupOption struct {
	Name string

	// Context defaults to context.Background()
	Context context.Context
	// Signals default to [syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT]
	Signals  []os.Signal
	Services []Service
}

//type ServiceContext interface {
//	Get(prefix string) (interface{}, bool)
//	MustGet(prefix string) interface{}
//}

type Service interface {
	//ServiceContext
	Name() string
	Start() error
	Stop() error
}

func NewAppGroup(opt AppGroupOption) *AppGroup {
	if opt.Context == nil {
		opt.Context = context.Background()
	}
	if opt.Signals == nil {
		opt.Signals = []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT}
	}
	ctx, cancel := context.WithCancel(opt.Context)
	return &AppGroup{
		option: opt,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (appGroup *AppGroup) Run() error {
	g, ctx := errgroup.WithContext(appGroup.ctx)

	for _, s := range appGroup.option.Services {
		srv := s
		g.Go(func() error {
			<-ctx.Done()
			return srv.Stop()
		})
		g.Go(func() error {
			fmt.Println(srv.Name())

			return srv.Start()
		})
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, appGroup.option.Signals...)
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-c:
				_ = appGroup.Stop()
			}
		}
	})
	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

func (appGroup *AppGroup) Stop() error {
	if appGroup.cancel != nil {
		appGroup.cancel()
	}
	return nil
}
