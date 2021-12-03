package butcher

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBufferSize(t *testing.T) {
	type args struct {
		size int
	}
	tests := []struct {
		name   string
		size   int
		hasErr bool
	}{
		{
			"valid size",
			5,
			false,
		},
		{
			"invalid size",
			0,
			true,
		},
		{
			"invalid size 2",
			-1,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := NewButcher(nil, BufferSize(tt.size))
			assert.Equal(t, err != nil, tt.hasErr)
			assert.Equal(t, b == nil, tt.hasErr)
		})
	}
}

func TestMaxWorker(t *testing.T) {
	type args struct {
		count int
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaxWorker(tt.args.count); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MaxWorker() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRateLimit(t *testing.T) {
	type args struct {
		rateLimit float64
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RateLimit(tt.args.rateLimit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RateLimit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRetryOnError(t *testing.T) {
	type args struct {
		maxTimes int
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RetryOnError(tt.args.maxTimes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RetryOnError() = %v, want %v", got, tt.want)
			}
		})
	}
}
