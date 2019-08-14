package utils

import (
	"github.com/pkg/errors"
	"time"
)

var (
	ERROR_RetryStopped = errors.New("stopped")
)

func Retry(f func() error, times int) (err error) {
	return Backoff(f, times, 0, 0, nil)
}

func Backoff(f func() error, times int, interval, max time.Duration, stopc <-chan struct{}) (err error) {
	current := interval
	for i := 0; i < times; i += 1 {
		if err = f(); err == nil {
			return
		} else if current > 0 {
			select {
			case <-stopc:
				return ERROR_RetryStopped
			case <-time.After(current):
				if current < max {
					current += interval
					if current >= max {
						current = max
					}
				}
			}
		}
	}
	return
}
