package butcher

import (
	"fmt"
	"time"

	"golang.org/x/time/rate"
)

type Option func(b *butcher) error

func MaxWorker(count int) Option {
	return func(b *butcher) error {
		if count <= 0 {
			return fmt.Errorf("worker count cannot be less than 0")
		}
		b.maxWorker = count
		return nil
	}
}

func RateLimit(tasksPerSecond float64) Option {
	return func(b *butcher) error {
		if tasksPerSecond <= 0 {
			return fmt.Errorf("rate cannot be less than 0")
		}
		b.rateLimit = rate.Limit(tasksPerSecond)
		return nil
	}
}

func BufferSize(size int) Option {
	return func(b *butcher) error {
		if size <= 0 {
			return fmt.Errorf("buffer size cannot be less than 0")
		}
		b.bufSize = size
		return nil
	}
}

func RetryOnError(maxTimes int) Option {
	return func(b *butcher) error {
		if maxTimes <= 0 {
			return fmt.Errorf("max retry times cannot be less than 0")
		}
		b.maxRetryTimes = maxTimes
		return nil
	}
}

func TaskTimeout(timeout time.Duration) Option {
	return func(b *butcher) error {
		if timeout <= 0 {
			return fmt.Errorf("timeout cannot be less than 0")
		}
		b.taskTimeout = timeout
		return nil
	}
}
