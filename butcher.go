package butcher

import (
	"context"
	"fmt"
	"golang.org/x/time/rate"
	"os"
	"os/signal"
	"sync"
	"time"
)

type Butcher interface {
	Run() error
}

type butcher struct {
	executor Executor

	maxWorker     int
	bufSize       int
	maxRetryTimes int
	rateLimit     rate.Limit
	taskTimeout   time.Duration

	limiter *rate.Limiter

	stopped bool

	jobCh      chan interface{}
	readyCh    chan job
	completeCh chan struct{}
}

func NewButcher(executor Executor, opts ...Option) (Butcher, error) {
	b := &butcher{
		executor: executor,

		maxWorker:     1,
		bufSize:       1,
		maxRetryTimes: 0,
		rateLimit:     rate.Inf,
		taskTimeout:   0,
	}
	for _, opt := range opts {
		if err := opt(b); err != nil {
			return nil, err
		}
	}

	return b, nil
}

func (b *butcher) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	b.limiter = rate.NewLimiter(b.rateLimit, b.bufSize)

	b.jobCh = make(chan interface{}, b.bufSize)
	b.readyCh = make(chan job, b.bufSize)
	b.completeCh = make(chan struct{})

	b.generate(ctx)
	b.rectify(ctx)
	b.schedule(ctx)

	signalCh := make(chan os.Signal)
	signal.Notify(signalCh)
	select {
	case s := <-signalCh:
		return fmt.Errorf("interrupted by signal: %v", s)
	case <-b.completeCh:
		return nil
	}
}

func (b *butcher) generate(ctx context.Context) {
	go func() {
		_ = safelyRun(func() error {
			err := b.executor.GenerateJob(ctx, b.jobCh)
			b.stopped = true
			close(b.jobCh)
			return err
		})
	}()
}

func (b *butcher) rectify(ctx context.Context) {
	go func() {
		for payload := range b.jobCh {
			_ = b.limiter.Wait(ctx)
			b.readyCh <- job{Type: jobTypeJob, Payload: payload}
		}
		close(b.readyCh)
	}()
}

func (b *butcher) schedule(ctx context.Context) {
	go func() {
		var wg sync.WaitGroup
		for i := 0; i < b.maxWorker; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := range b.readyCh {
					b.task(ctx, j)
				}
			}()
		}
		wg.Wait()
		b.completeCh <- struct{}{}
	}()
}

func (b *butcher) task(ctx context.Context, j job) {
	if b.taskTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, b.taskTimeout)
		defer cancel()
	}
	errCh := make(chan error)
	go func() {
		errCh <- safelyRun(func() error {
			return b.executor.Task(ctx, j.Payload)
		})
	}()

	select {
	case <-ctx.Done():
		b.onError(ctx, j, ctx.Err())
	case err := <-errCh:
		if err != nil {
			b.onError(ctx, j, err)
		} else {
			b.onFinish(ctx, j)
		}
	}
}

func (b *butcher) onError(ctx context.Context, j job, err error) {
	if j.RetryTime < b.maxRetryTimes && !b.stopped {
		j.RetryTime++
		b.jobCh <- j
		return
	}
	if watcher, ok := b.executor.(OnErrorWatcher); ok {
		_ = safelyRun(func() error {
			watcher.OnError(ctx, j.Payload, err)
			return nil
		})
	}
}

func (b *butcher) onFinish(ctx context.Context, j job) {
	if watcher, ok := b.executor.(OnFinishWatcher); ok {
		_ = safelyRun(func() error {
			watcher.OnFinish(ctx, j.Payload)
			return nil
		})
	}
}
