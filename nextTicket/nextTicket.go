package nextTicket

import (
	"github.com/pkg/errors"
	"time"
)

var (
	Error_InvalidDuration = errors.New("Invalid duration.")
	Error_NextTimeout     = errors.New("Push next timeout.")
)

type Ticket struct {
	C    <-chan time.Time
	t    chan time.Duration
	stop chan int
	d    time.Duration
}

func NewTicket(ds ...time.Duration) (*Ticket, error) {
	var d time.Duration
	if len(ds) > 0 {
		d = ds[0]
	}
	if d < 0 {
		return nil, Error_InvalidDuration
	}
	c := make(chan time.Time, 1)
	t := &Ticket{
		C:    c,
		t:    make(chan time.Duration, 1),
		stop: make(chan int),
	}
	if d > 0 {
		t.Next(d)
	}
	go t.loop(c)
	return t, nil
}

func (t *Ticket) Next(d time.Duration) error {
	if d < 0 {
		return Error_InvalidDuration
	}
	select {
	case t.t <- d:
		break
	case <-time.After(time.Millisecond * 10):
		return Error_NextTimeout
	}
	return nil
}

func (t *Ticket) loop(c chan time.Time) {
	for {
		select {
		case <-t.stop:
			return
		case d := <-t.t:
			t.d = d
		case v := <-time.After(t.d):
			select {
			case <-t.stop:
				return
			case c <- v:
				break
			}
		}
	}
}

func (t *Ticket) Stop() {
	close(t.stop)
	close(t.t)
}
