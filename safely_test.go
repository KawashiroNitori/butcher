package butcher

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSafelyRun(t *testing.T) {

	tests := []struct {
		name   string
		f      func() error
		errMsg string
	}{
		{
			"Normal",
			func() error {
				return nil
			},
			"",
		},
		{
			"ReturnError",
			func() error {
				return fmt.Errorf("boom")
			},
			"boom",
		},
		{
			"PanicError",
			func() error {
				panic(fmt.Errorf("panic boom"))
			},
			"panic boom",
		},
		{
			"PanicSomething",
			func() error {
				panic("wtf")
			},
			"unexpected panic occurred: wtf",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := safelyRun(tt.f)
			if tt.errMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.errMsg)
			}
		})
	}
}
