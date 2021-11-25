package consumer

import (
	"context"
	"example/asyncjob"
	"example/common"
	core_service "github.com/cesc1802/core-service"
	"github.com/cesc1802/core-service/events"
	"log"
)

type engine struct {
	sc core_service.Service
}

type consumerJob struct {
	Title   string
	Handler func(ctx context.Context, message interface{}) error
}

func NewEngine(sc core_service.Service) *engine {
	return &engine{
		sc: sc,
	}
}

func (e *engine) Start() error {
	e.subscribe("test", true, SendEmailOTP(e.sc), SendSMSOTP(e.sc))
	return nil
}

func (e *engine) subscribe(topic string, isParallel bool, jobs ...consumerJob) error {
	event, err := e.sc.MustGet(common.KeyPubSub).(events.Stream).Consume(topic)

	for _, job := range jobs {
		log.Println("Setup consumer job for:", job.Title)
	}

	getJobHandler := func(job *consumerJob, message interface{}) asyncjob.JobHandler {
		return func(ctx context.Context) error {
			return job.Handler(ctx, message)
		}
	}
	if err != nil {
		log.Println(err)
	}

	go func(event <-chan events.Event) {
		for {
			msg := <-event

			jobHdlArr := make([]asyncjob.Job, len(jobs))

			for i := range jobs {
				jobHdlArr[i] = asyncjob.NewJob(getJobHandler(&jobs[i], msg))
			}

			group := asyncjob.NewManager(isParallel, jobHdlArr...)

			if err := group.Run(context.Background()); err != nil {
				log.Println(err)
			}
		}
	}(event)

	return nil
}
