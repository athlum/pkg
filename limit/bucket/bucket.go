package limitBucket

import (
	"github.com/athlum/pkg/exitChan"
	"github.com/athlum/pkg/nextTicket"
	"github.com/athlum/pkg/utils"
	"sync"
)

type Engine struct {
	lock  *sync.Mutex
	ch    chan int
	stock chan int
	next  *Config
	stop  *exitChan.ExitChan
}

func New(cfg *Config) (*Engine, error) {
	cfg.changed = true
	e := &Engine{
		lock:  &sync.Mutex{},
		ch:    make(chan int),
		stock: make(chan int),
		stop:  exitChan.NewExitChan(),
		next:  cfg,
	}
	t, err := nextTicket.NewTicket(utils.Duration(e.next.Interval))
	if err != nil {
		return nil, err
	}
	go e.stockLoop()
	go e.loop(t)
	return e, nil
}

func (e *Engine) loop(t *nextTicket.Ticket) {
	e.stock <- e.next.Limit
	defer t.Stop()
	for {
		select {
		case <-e.stop.Chan():
			return
		case <-t.C:
			e.restock(t)
		}
	}
}

func (e *Engine) lockedNext() (Config, bool) {
	e.lock.Lock()
	defer e.lock.Unlock()
	changed := e.next.changed
	if changed {
		e.next.changed = false
	}
	return *e.next, changed
}

func (e *Engine) restock(t *nextTicket.Ticket) {
	n, c := e.lockedNext()
	if c {
		t.Next(utils.Duration(n.Interval))
	}
	e.stock <- n.Limit
}

func (e *Engine) stockLoop() {
	for {
		select {
		case <-e.stop.Chan():
			return
		case v := <-e.stock:
			for i := 1; i <= v; i += 1 {
				select {
				case <-e.stop.Chan():
					return
				case e.ch <- i:
					continue
				}
			}
		}
	}
}

func (e *Engine) Update(cfg *Config) {
	cfg.changed = true

	e.lock.Lock()
	defer e.lock.Unlock()
	e.next = cfg
}

func (e *Engine) Chan() chan int {
	return e.ch
}

func (e *Engine) Close() {
	if !e.stop.Exited() {
		close(e.ch)
		close(e.stock)
	}
	e.stop.Close()
}
