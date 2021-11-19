package asyncjob

import (
	"context"
	"log"
	"sync"
)

type manager struct {
	isConcurrent bool
	jobs         []Job
	wg           *sync.WaitGroup
}

func NewManager(isParallel bool, jobs ...Job) *manager {
	g := &manager{
		isConcurrent: isParallel,
		jobs:         jobs,
		wg:           new(sync.WaitGroup),
	}

	return g
}

func (g *manager) Run(ctx context.Context) error {
	g.wg.Add(len(g.jobs))

	errChan := make(chan error, len(g.jobs))

	for i, _ := range g.jobs {
		if g.isConcurrent {
			// Do this instead
			go func(aj Job) {
				errChan <- g.runJob(ctx, aj)
				g.wg.Done()
			}(g.jobs[i])

			continue
		}

		job := g.jobs[i]
		errChan <- g.runJob(ctx, job)
		g.wg.Done()
	}

	var err error

	for i := 1; i <= len(g.jobs); i++ {
		if v := <-errChan; v != nil {
			err = v
			//break
		}
	}

	g.wg.Wait()
	return err
}

// Retry if needed
func (g *manager) runJob(ctx context.Context, j Job) error {
	if err := j.Execute(ctx); err != nil {
		for {
			log.Println(err)
			if j.State() == StateRetryFailed {
				return err
			}

			if j.Retry(ctx) == nil {
				return nil
			}
		}
	}
	return nil
}
