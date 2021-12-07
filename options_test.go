package butcher

import (
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBufferSize(t *testing.T) {
	tests := []struct {
		name   string
		size   int
		errMsg string
	}{
		{
			"ValidSize",
			5,
			"",
		},
		{
			"InvalidSize",
			0,
			"buffer size cannot be less than 0",
		},
		{
			"InvalidSize2",
			-1,
			"buffer size cannot be less than 0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewButcher(nil, BufferSize(tt.size))
			if tt.errMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.errMsg)
			}
		})
	}
}

func TestMaxWorker(t *testing.T) {
	tests := []struct {
		name   string
		count  int
		errMsg string
	}{
		{
			"ValidCount",
			5,
			"",
		},
		{
			"InvalidCount",
			0,
			"worker count cannot be less than 0",
		},
		{
			"InvalidSize2",
			-1,
			"worker count cannot be less than 0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewButcher(nil, MaxWorker(tt.count))
			if tt.errMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.errMsg)
			}
		})
	}
}

func TestRateLimit(t *testing.T) {
	tests := []struct {
		name      string
		rateLimit float64
		errMsg    string
	}{
		{
			"ValidRate",
			5,
			"",
		},
		{
			"InvalidRate",
			0,
			"rate cannot be less than 0",
		},
		{
			"InvalidRate2",
			-1,
			"rate cannot be less than 0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewButcher(nil, RateLimit(tt.rateLimit))
			if tt.errMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.errMsg)
			}
		})
	}
}

func TestRetryOnError(t *testing.T) {
	tests := []struct {
		name     string
		maxTimes int
		errMsg   string
	}{
		{
			"ValidTime",
			5,
			"",
		},
		{
			"InvalidTime",
			0,
			"max retry times cannot be less than 0",
		},
		{
			"InvalidTime2",
			-1,
			"max retry times cannot be less than 0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewButcher(nil, RetryOnError(tt.maxTimes))
			if tt.errMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.errMsg)
			}
		})
	}
}

func TestTaskTimeout(t *testing.T) {
	tests := []struct {
		name    string
		timeout time.Duration
		errMsg  string
	}{
		{
			"ValidTimeout",
			5 * time.Second,
			"",
		},
		{
			"InvalidTimeout",
			0,
			"timeout cannot be less than 0",
		},
		{
			"InvalidTimeout2",
			-1,
			"timeout cannot be less than 0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewButcher(nil, TaskTimeout(tt.timeout))
			if tt.errMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.errMsg)
			}
		})
	}
}

func TestInterruptSignal(t *testing.T) {
	tests := []struct {
		name    string
		signals []os.Signal
		errMsg  string
	}{
		{
			"ValidSignals",
			[]os.Signal{syscall.SIGINT, syscall.SIGHUP},
			"",
		},
		{
			"EmptySignals",
			[]os.Signal{},
			"",
		},
		{
			"NilSignals",
			nil,
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewButcher(nil, InterruptSignal(tt.signals...))
			if tt.errMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.errMsg)
			}
		})
	}
}
