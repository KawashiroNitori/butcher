package butcher

import (
	"context"
	"fmt"
	"syscall"
	"time"
)

func Example_newButcher_Options() {
	ctx := context.Background()
	executor := &basicExecutor{}
	b, err := NewButcher(
		executor,
		MaxWorker(3),               // specifies maximum number of concurrent workers, default is 1.
		BufferSize(20),             // specifies job buffer size, recommended value is bigger than MaxWorker and RateLimit, default is 1.
		RateLimit(10.0),            // control task execute speed, default is Infinity.
		TaskTimeout(1*time.Second), // specifies task execute timeout, returns a context.TimeExceeded error if timeout is exceeded, default no timeout.
		RetryOnError(3),            // specifies retry times when task return error, default no retry.
		InterruptSignal(syscall.SIGINT, syscall.SIGTERM), // specified signals can interrupt task running, default is syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT.
	)
	if err != nil {
		panic(err)
	}
	err = b.Run(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("all done")
}
