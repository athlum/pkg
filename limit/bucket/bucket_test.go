package limitBucket

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestBucket(t *testing.T) {
	e, err := New(&Config{
		Limit:    5,
		Interval: 10.0,
	})
	if err != nil {
		t.Error(err)
	}
	defer e.Close()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	for i := 0; i < 2; i += 1 {
		go func(i int, e *Engine) {
			for v := range e.Chan() {
				fmt.Printf("%v echo %v: %v\n", time.Now(), i, v)
			}
		}(i, e)
	}

	time.Sleep(time.Second * 20)
	e.Update(&Config{
		Limit:    1,
		Interval: 5,
	})
	wg.Wait()
}

func TestBlockBucket(t *testing.T) {
	e, err := New(&Config{
		Limit:    2400000000,
		Interval: 1.0,
	})
	if err != nil {
		t.Error(err)
	}
	defer e.Close()

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func(e *Engine) {
		for v := range e.Chan() {
			if v == 2400000000 {
				fmt.Printf("%v echo: %v\n", time.Now(), v)
			}
		}
	}(e)

	wg.Wait()
}
