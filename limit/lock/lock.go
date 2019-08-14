package limitLock

import (
	"sync"
)

type Lock struct {
	_lock *sync.Mutex

	limit  int64
	val    int64
	locked int32
	lock   *sync.Mutex
}

func New(limit int64) *Lock {
	return &Lock{
		_lock: &sync.Mutex{},
		limit: limit,
		lock:  &sync.Mutex{},
	}
}

func (l *Lock) Lock(v int64) {
	l.lock.Lock()
	l._lock.Lock()
	defer l._lock.Unlock()
	l.val += v
	l.locked = 1
}

func (l *Lock) Release(v int64) {
	l._lock.Lock()
	defer l._lock.Unlock()
	l.val -= v
	l.unlock()
}

func (l *Lock) Unlock() {
	l._lock.Lock()
	defer l._lock.Unlock()
	l.unlock()
}

func (l *Lock) unlock() {
	if l.locked != 1 {
		return
	}
	if l.val < l.limit {
		l.locked = 0
		l.lock.Unlock()
	}
}
