package limitLock

import (
	"fmt"
	"testing"
	"time"
)

func emit(l *Lock, ch chan int, v int) {
	l.Lock(1)
	defer l.Unlock()
	fmt.Println(time.Now(), "emit", v)
	ch <- v
}

func read(l *Lock, ch chan int) {
	fmt.Println(time.Now(), "read", <-ch)
	l.Release(1)
}

func TestLock(t *testing.T) {
	l := New(5)
	ch := make(chan int, 10)
	go func() {
		for i := 0; i < 10; i += 1 {
			time.Sleep(time.Millisecond)
			emit(l, ch, i)
		}
	}()

	go func() {
		for i := 0; i < 5; i += 1 {
			time.Sleep(time.Second * 2)
			read(l, ch)
		}
	}()

	for i := 0; i < 5; i += 1 {
		time.Sleep(time.Second * 2)
		read(l, ch)
	}
}
