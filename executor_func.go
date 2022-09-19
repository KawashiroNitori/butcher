package butcher

import (
	"context"
)

type JobFunc[T any] func(ctx context.Context, jobCh chan<- T) error
type TaskFunc[T any] func(ctx context.Context, job T) error

type funcExecutor[T any] struct {
	jobFunc  JobFunc[T]
	taskFunc TaskFunc[T]
}

var _ Executor[any] = new(funcExecutor[any])

func NewFuncExecutor[T any](jobFunc JobFunc[T], taskFunc TaskFunc[T]) Executor[T] {
	return &funcExecutor[T]{
		jobFunc:  jobFunc,
		taskFunc: taskFunc,
	}
}

func (f *funcExecutor[T]) GenerateJob(ctx context.Context, jobCh chan<- T) error {
	return f.jobFunc(ctx, jobCh)
}

func (f *funcExecutor[T]) Task(ctx context.Context, job T) error {
	return f.taskFunc(ctx, job)
}
