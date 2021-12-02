package butcher

import (
	"fmt"

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

func RateLimit(rateLimit float64) Option {
	return func(b *butcher) error {
		if rateLimit <= 0 {
			return fmt.Errorf("rate cannot be less than 0")
		}
		b.rateLimit = rate.Limit(rateLimit)
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
