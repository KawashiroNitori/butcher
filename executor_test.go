package butcher

import (
	"context"
	"fmt"
	"time"
)

type basicExecutor struct {
	Size    int
	Results []bool
}

func (b *basicExecutor) GenerateJob(ctx context.Context, jobCh chan<- int) error {
	for i := 0; i < b.Size; i++ {
		jobCh <- i
	}
	return nil
}

func (b *basicExecutor) Task(ctx context.Context, job int) error {
	i := job
	b.Results[i] = true
	return nil
}

type retryExecutor struct {
	Size     int
	Results  []int
	Finished []bool
	Errors   []error
}

func (r *retryExecutor) GenerateJob(ctx context.Context, jobCh chan<- int) error {
	for i := 0; i < r.Size; i++ {
		jobCh <- i
	}
	return nil
}

func (r *retryExecutor) Task(ctx context.Context, job int) error {
	i := job
	r.Results[i]++
	if v := r.Results[i]; v <= i {
		panic(fmt.Errorf("boom"))
	}
	return nil
}

func (r *retryExecutor) OnFinish(ctx context.Context, job int, err error) {
	i := job
	if err == nil {
		r.Finished[i] = true
	} else {
		r.Errors[i] = err
	}
}

type taskTimeoutExecutor struct {
	Size    int
	Results []bool
	Errors  []error
}

func (t *taskTimeoutExecutor) GenerateJob(ctx context.Context, jobCh chan<- int) error {
	for i := 0; i < t.Size; i++ {
		jobCh <- i
	}
	return nil
}

func (t *taskTimeoutExecutor) Task(ctx context.Context, job int) error {
	i := time.Duration(job)
	time.Sleep(i * 100 * time.Millisecond)
	t.Results[i] = true
	return nil
}

func (t *taskTimeoutExecutor) OnFinish(ctx context.Context, job int, err error) {
	i := job
	t.Errors[i] = err
}

type signalInterruptExecutor struct {
	Size    int
	Results []bool
}

func (s *signalInterruptExecutor) GenerateJob(ctx context.Context, jobCh chan<- int) error {
	for i := 0; i < s.Size; i++ {
		time.Sleep(100 * time.Millisecond)
		jobCh <- i
	}
	return nil
}

func (s *signalInterruptExecutor) Task(ctx context.Context, job int) error {
	i := job
	s.Results[i] = true
	return nil
}

type generatorErrorExecutor struct {
	Size    int
	Results []bool
}

func (g *generatorErrorExecutor) GenerateJob(ctx context.Context, jobCh chan<- int) error {
	for i := 0; i < g.Size; i++ {
		if i >= 3 {
			panic(fmt.Errorf("boom"))
		}
		jobCh <- i
	}
	return nil
}

func (g *generatorErrorExecutor) Task(ctx context.Context, job int) error {
	i := job
	g.Results[i] = true
	return nil
}
