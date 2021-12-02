package butcher

import "context"

type Executor interface {
	GenerateJob(ctx context.Context, jobCh chan<- interface{}) error
	Task(ctx context.Context, job interface{}) error
}

type OnFinishWatcher interface {
	OnFinish(ctx context.Context, job interface{})
}

type OnErrorWatcher interface {
	OnError(ctx context.Context, job interface{}, err error)
}
