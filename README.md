# butcher
[![CI Build](https://github.com/KawashiroNitori/butcher/actions/workflows/ci.yml/badge.svg)](https://github.com/KawashiroNitori/butcher/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/KawashiroNitori/butcher.svg)](https://pkg.go.dev/github.com/KawashiroNitori/butcher)
[![codecov](https://codecov.io/gh/KawashiroNitori/butcher/branch/master/graph/badge.svg?token=RB17B7IOMN)](https://codecov.io/gh/KawashiroNitori/butcher)
[![Go Report Card](https://goreportcard.com/badge/github.com/KawashiroNitori/butcher)](https://goreportcard.com/report/github.com/KawashiroNitori/butcher)

# Overview
Butcher is a library providing a simple way to execute some task concurrency. Integrates convenient features such as concurrency control and retry.

# Quick Start

```go
package main

import (
    "context"
    "fmt"
    "github.com/KawashiroNitori/butcher"
)

type executor struct{}

func (e *executor) GenerateJob(ctx context.Context, jobCh chan<- int) error {
    // generate your jobs here
    for i := 0; i < 100; i++ {
        // you can check canceled context or not, is up to you
        select {
        case <- ctx.Done():
            return nil
        default:
        }

        // push your job into channel
        jobCh <- i
    }
    // you may NOT close jobCh manually
    return nil
}

func (e *executor) Task(ctx context.Context, job int) error {
    // execute your job here
    fmt.Printf("job %v finished!\n", job)
    return nil
}

// OnFinish implement an optional func to check your job is finished if you want
func (e *executor) OnFinish(ctx context.Context, job int, err error) {
    if err != nil {
        fmt.Printf("job %v error: %v", job, err)
    }
}

func main() {
    ctx := context.Background()
    b, err := butcher.NewButcher[int](&executor{})  // you can add some options here
    if err != nil {
        panic(err)
    }

    err = b.Run(ctx)
    if err != nil {
        panic(err)
    }
}

```

