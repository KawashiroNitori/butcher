package butcher

import (
	"golang.org/x/time/rate"
)

type Butcher interface {
	Run() error
}

type butcher struct {
	executor Executor

	maxWorker int
	bufSize   int
	rateLimit rate.Limit

	limiter *rate.Limiter
}

func NewButcher(executor Executor, opts ...Option) (Butcher, error) {
	b := &butcher{
		executor: executor,

		maxWorker: 1,
		bufSize:   1,
		rateLimit: rate.Inf,
	}
	for _, opt := range opts {
		if err := opt(b); err != nil {
			return nil, err
		}
	}
	b.limiter = rate.NewLimiter(b.rateLimit, b.bufSize)
	return b, nil
}

func (b *butcher) Run() error {
	panic("implement me")
}
