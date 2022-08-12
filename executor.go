package butcher

import "context"

// Executor implements task executor interface.
type Executor[T any] interface {
	// GenerateJob generate your jobs here. put your job into jobCh, don't close jobCh manually.
	GenerateJob(ctx context.Context, jobCh chan<- T) error
	// Task execute your job here. It will be scheduled by butcher.
	Task(ctx context.Context, job T) error
}

// OnFinishWatcher implements optional OnFinish function if you want to watch the result of job.
type OnFinishWatcher[T any] interface {
	OnFinish(ctx context.Context, job T, err error)
}
