package butcher

import "fmt"

func safelyRun(f func() error) (err error) {
	defer func() {
		if v := recover(); v != nil {
			if e, ok := v.(error); ok {
				err = e
			}
		} else {
			err = fmt.Errorf("error occurred: %v", v)
		}
	}()

	return f()
}
