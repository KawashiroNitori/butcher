package butcher

import (
	"context"
	"fmt"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestButcherBasic(t *testing.T) {
	size := 1000
	executor := &basicExecutor{
		Size:    size,
		Results: make([]bool, size),
	}
	b, err := NewButcher[int](executor, MaxWorker(3), BufferSize(1))
	expected := make([]bool, size)
	for i := range expected {
		expected[i] = true
	}
	assert.NoError(t, err)
	assert.NoError(t, b.Run(context.Background()))
	assert.EqualValues(t, expected, executor.Results)
}

func TestButcherRetry(t *testing.T) {
	size := 6
	executor := &retryExecutor{
		Size:     size,
		Results:  make([]int, size),
		Finished: make([]bool, size),
		Errors:   make([]error, size),
	}
	b, err := NewButcher[int](executor, MaxWorker(3), BufferSize(3), RetryOnError(3))
	assert.NoError(t, err)
	assert.NoError(t, b.Run(context.Background()))
	assert.EqualValues(t, []int{1, 2, 3, 4, 4, 4}, executor.Results)
	assert.EqualValues(t, []bool{true, true, true, true, false, false}, executor.Finished)
	assert.EqualValues(t, []error{nil, nil, nil, nil, fmt.Errorf("boom"), fmt.Errorf("boom")}, executor.Errors)
}

func TestButcherTaskTimeout(t *testing.T) {
	size := 5
	executor := &taskTimeoutExecutor{
		Size:    size,
		Errors:  make([]error, size),
		Results: make([]bool, size),
	}
	b, err := NewButcher[int](executor, MaxWorker(5), BufferSize(5), TaskTimeout(250*time.Millisecond))
	assert.NoError(t, err)
	assert.NoError(t, b.Run(nil))
	assert.EqualValues(t, []bool{true, true, true, false, false}, executor.Results)

	assert.NoError(t, executor.Errors[0])
	assert.NoError(t, executor.Errors[1])
	assert.NoError(t, executor.Errors[2])
	assert.ErrorIs(t, executor.Errors[3], context.DeadlineExceeded)
	assert.ErrorIs(t, executor.Errors[4], context.DeadlineExceeded)
}

func TestButcherInterruptedBySignal(t *testing.T) {
	size := 5
	executor := &signalInterruptExecutor{
		Size:    size,
		Results: make([]bool, size),
	}
	bb, err := NewButcher[int](executor, MaxWorker(5), BufferSize(5))
	b := bb.(*butcher[int])
	assert.NoError(t, err)
	go func() {
		time.Sleep(250 * time.Millisecond)
		b.signalCh <- syscall.SIGINT
	}()
	assert.EqualError(t, b.Run(context.Background()), "interrupted by signal: interrupt")
	assert.EqualValues(t, []bool{true, true, false, false, false}, executor.Results)
}

func TestButcherGeneratorError(t *testing.T) {
	size := 5
	executor := &generatorErrorExecutor{
		Size:    size,
		Results: make([]bool, size),
	}
	b, err := NewButcher[int](executor, MaxWorker(3), BufferSize(3))
	assert.NoError(t, err)
	assert.EqualError(t, b.Run(context.Background()), "generator error occurred: boom")
	assert.EqualValues(t, []bool{true, true, true, false, false}, executor.Results)
}

func TestButcherContextCancel(t *testing.T) {
	size := 5
	executor := &basicExecutor{
		Size:    size,
		Results: make([]bool, size),
	}
	b, err := NewButcher[int](executor, RateLimit(10))
	assert.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(150 * time.Millisecond)
		cancel()
	}()

	err = b.Run(ctx)
	assert.ErrorIs(t, err, context.Canceled)
	assert.EqualValues(t, []bool{true, true, false, false, false}, executor.Results)
}

func TestButcherFunc(t *testing.T) {
	var sum int
	var lock sync.Mutex

	executor := NewFuncExecutor(
		func(ctx context.Context, jobCh chan<- int) error {
			for i := 1; i <= 100; i++ {
				jobCh <- i
			}
			return nil
		}, func(ctx context.Context, job int) error {
			lock.Lock()
			defer lock.Unlock()
			sum += job
			return nil
		},
	)

	b, err := NewButcher(executor, MaxWorker(3), BufferSize(1))
	assert.NoError(t, err)
	assert.NoError(t, b.Run(context.Background()))
	assert.EqualValues(t, 5050, sum)
}
