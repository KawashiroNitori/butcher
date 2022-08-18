package butcher

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/time/rate"
)

type Option func(b *butcherCfg) error

// MaxWorker specifies maximum number of concurrent workers, default is 1.
func MaxWorker(count int) Option {
	return func(b *butcherCfg) error {
		if count <= 0 {
			return fmt.Errorf("worker count cannot be less than 0")
		}
		b.maxWorker = count
		return nil
	}
}

// RateLimit control task execute speed, default is Infinity.
func RateLimit(tasksPerSecond float64) Option {
	return func(b *butcherCfg) error {
		if tasksPerSecond <= 0 {
			return fmt.Errorf("rate cannot be less than 0")
		}
		b.rateLimit = rate.Limit(tasksPerSecond)
		return nil
	}
}

// BufferSize specifies job buffer size, recommended value is bigger than MaxWorker and RateLimit, default is 1.
func BufferSize(size int) Option {
	return func(b *butcherCfg) error {
		if size <= 0 {
			return fmt.Errorf("buffer size cannot be less than 0")
		}
		b.bufSize = size
		return nil
	}
}

// RetryOnError specifies retry times when task return error, default no retry.
func RetryOnError(maxTimes int) Option {
	return func(b *butcherCfg) error {
		if maxTimes <= 0 {
			return fmt.Errorf("max retry times cannot be less than 0")
		}
		b.maxRetryTimes = maxTimes
		return nil
	}
}

// TaskTimeout specifies task execute timeout, returns a context.TimeExceeded error if timeout is exceeded, default no timeout.
func TaskTimeout(timeout time.Duration) Option {
	return func(b *butcherCfg) error {
		if timeout <= 0 {
			return fmt.Errorf("timeout cannot be less than 0")
		}
		b.taskTimeout = timeout
		return nil
	}
}

// InterruptSignal specified signals can interrupt task running.
// Default is syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT.
func InterruptSignal(signals ...os.Signal) Option {
	return func(b *butcherCfg) error {
		b.interruptSignals = signals
		return nil
	}
}
