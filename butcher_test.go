package butcher

import (
	"context"
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
	b, err := NewButcher(executor, MaxWorker(3), BufferSize(1))
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
		T:          t,
		Size:       size,
		Results:    make([]int, size),
		RetryTimes: make([]int, size),
		Finished:   make([]bool, size),
	}
	b, err := NewButcher(executor, MaxWorker(3), BufferSize(3), RetryOnError(3))
	assert.NoError(t, err)
	assert.NoError(t, b.Run(context.Background()))
	assert.EqualValues(t, []int{1, 2, 3, 4, 4, 4}, executor.Results)
	assert.EqualValues(t, []int{0, 0, 1, 2, 3, 3}, executor.RetryTimes)
	assert.EqualValues(t, []bool{true, true, true, true, false, false}, executor.Finished)
}

func TestButcherTaskTimeout(t *testing.T) {
	size := 5
	executor := &taskTimeoutExecutor{
		Size:    size,
		Errors:  make([]error, size),
		Results: make([]bool, size),
	}
	b, err := NewButcher(executor, MaxWorker(5), BufferSize(5), TaskTimeout(250*time.Millisecond))
	assert.NoError(t, err)
	assert.NoError(t, b.Run(context.Background()))
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
	bb, err := NewButcher(executor, MaxWorker(5), BufferSize(5))
	b := bb.(*butcher)
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
	b, err := NewButcher(executor, MaxWorker(3), BufferSize(3))
	assert.NoError(t, err)
	assert.EqualError(t, b.Run(context.Background()), "generator error occurred: boom")
	assert.EqualValues(t, []bool{true, true, true, false, false}, executor.Results)
}
