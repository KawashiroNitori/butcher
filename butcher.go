package butcher

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"golang.org/x/time/rate"
)

type Butcher interface {
	Run(context.Context) error
}

type butcherCfg struct {
	maxWorker        int
	bufSize          int
	maxRetryTimes    int
	rateLimit        rate.Limit
	taskTimeout      time.Duration
	interruptSignals []os.Signal
}

type butcher[T any] struct {
	*butcherCfg
	executor Executor[T]

	limiter *rate.Limiter

	jobCh         chan T
	readyCh       chan job[T]
	generateErrCh chan error
	completeCh    chan struct{}
	signalCh      chan os.Signal
}

// NewButcher returns a butcher object for execute task executor. It has some options to control execute behaviors.
// if no options given, it runs tasks serially.
func NewButcher[T any](executor Executor[T], opts ...Option) (Butcher, error) {
	b := &butcher[T]{
		executor: executor,
		butcherCfg: &butcherCfg{
			maxWorker:        1,
			bufSize:          1,
			maxRetryTimes:    0,
			rateLimit:        rate.Inf,
			taskTimeout:      0,
			interruptSignals: []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT},
		},
	}
	for _, opt := range opts {
		if err := opt(b.butcherCfg); err != nil {
			return nil, err
		}
	}

	b.limiter = rate.NewLimiter(b.rateLimit, b.bufSize)

	b.jobCh = make(chan T, b.bufSize)
	b.readyCh = make(chan job[T], b.bufSize)
	b.generateErrCh = make(chan error, 1)
	b.completeCh = make(chan struct{}, 1)
	b.signalCh = make(chan os.Signal, 1)

	return b, nil
}

// Run the task executor, return error if interrupted or GenerateJob return an error.
func (b *butcher[T]) Run(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel := context.WithCancel(ctx)

	b.generate(ctx)
	b.rectify(ctx)
	b.schedule(ctx)

	if len(b.interruptSignals) > 0 {
		signal.Notify(b.signalCh, b.interruptSignals...)
	}
	select {
	case <-ctx.Done():
		cancel()
		<-b.completeCh
		return fmt.Errorf("interrupted by context: %w", ctx.Err())
	case s := <-b.signalCh:
		cancel()
		<-b.completeCh
		return fmt.Errorf("interrupted by signal: %v", s)
	case err := <-b.generateErrCh:
		// waiting for scheduled jobs complete
		<-b.completeCh
		cancel()
		return fmt.Errorf("generator error occurred: %w", err)
	case <-b.completeCh:
		cancel()
		return nil
	}
}

func (b *butcher[T]) generate(ctx context.Context) {
	go func() {
		err := safelyRun(func() error {
			err := b.executor.GenerateJob(ctx, b.jobCh)
			return err
		})
		if err != nil {
			b.generateErrCh <- err
		}
		close(b.jobCh)
	}()
}

func (b *butcher[T]) rectify(ctx context.Context) {
	go func() {
		for payload := range b.jobCh {
			select {
			case <-ctx.Done():
				close(b.readyCh)
				return
			default:
			}
			_ = b.limiter.Wait(ctx)
			b.readyCh <- job[T]{Type: jobTypeJob, Payload: payload}
		}
		close(b.readyCh)
	}()
}

func (b *butcher[T]) schedule(ctx context.Context) {
	var wg sync.WaitGroup

	runFunc := func() {
		defer wg.Done()
		for j := range b.readyCh {
			select {
			case <-ctx.Done():
				return
			default:
			}
			var err error
			for j.RetryTime <= b.maxRetryTimes {
				if err = b.task(ctx, j); err != nil {
					j.RetryTime++
				} else {
					break
				}
			}
			b.onFinish(ctx, j, err)
		}
	}

	go func() {
		for i := 0; i < b.maxWorker; i++ {
			wg.Add(1)
			go runFunc()
		}
		wg.Wait()
		b.completeCh <- struct{}{}
	}()
}

func (b *butcher[T]) task(ctx context.Context, j job[T]) (err error) {
	if b.taskTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, b.taskTimeout)
		defer cancel()
	}
	errCh := make(chan error, 1)
	go func() {
		errCh <- safelyRun(func() error {
			return b.executor.Task(ctx, j.Payload)
		})
	}()

	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-errCh:
	}

	return err
}

func (b *butcher[T]) onFinish(ctx context.Context, j job[T], err error) {
	if watcher, ok := b.executor.(OnFinishWatcher[T]); ok {
		_ = safelyRun(func() error {
			watcher.OnFinish(ctx, j.Payload, err)
			return nil
		})
	}
}
