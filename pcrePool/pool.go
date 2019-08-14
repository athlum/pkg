package pcrePool

import (
	"github.com/athlum/golang-pkg-pcre/src/pkg/pcre"
	"github.com/pkg/errors"
	"hash/fnv"
	"sync"
)

var PCRE_ERROR_MATCHLIMIT = pcre.PCRE_ERROR_MATCHLIMIT

var (
	_pool *pool
	once  = &sync.Once{}
)

var ERROR_NeedInit = errors.New("pcre pool need init first.")

func Init() {
	once.Do(_init)
}

func _init() {
	if _pool == nil {
		_pool = New(16)
	}
}

func Compile(ps string) (*statedRegex, error) {
	if _pool == nil {
		return nil, ERROR_NeedInit
	}
	return _pool.Compile(ps)
}

type statedRegex struct {
	*pcre.Regexp
	counter *counter
	using   bool
}

func (s *statedRegex) Collect() {
	s.counter.Lock()
	defer s.counter.Unlock()

	s.using = false
}

type counter struct {
	sync.Mutex
	index int
	pool  []*statedRegex
}

func (c *counter) getOrCompile(p string) (*statedRegex, error) {
	c.Lock()
	defer c.Unlock()

	length := len(c.pool)
	for i := 0; i < length; i += 1 {
		ii := i + c.index
		if ii >= length {
			ii -= length
		}
		s := c.pool[ii]
		if !s.using {
			s.using = true
			c.index = ii
			return s, nil
		}
	}

	re, err := pcre.Compile(p, 0)
	if err != nil {
		return nil, err
	}
	s := &statedRegex{
		Regexp:  &re,
		counter: c,
		using:   true,
	}
	if c.index == length-1 {
		c.index += 1
	}
	c.pool = append(c.pool, s)
	return s, nil
}

type poolShard struct {
	res map[string]*counter
	sync.Mutex
}

func (s *poolShard) getOrCompile(p string) (*statedRegex, error) {
	s.Lock()
	defer s.Unlock()

	c, ok := s.res[p]
	if !ok {
		c = &counter{}
		s.res[p] = c
	}
	return c.getOrCompile(p)
}

type pool struct {
	shards []*poolShard
	count  int
}

func New(c int) *pool {
	p := &pool{make([]*poolShard, c), c}
	for i := 0; i < p.count; i += 1 {
		p.shards[i] = &poolShard{res: make(map[string]*counter)}
	}
	return p
}

func (p *pool) getShard(key string) *poolShard {
	hasher := fnv.New32()
	hasher.Write([]byte(key))
	return p.shards[uint(hasher.Sum32())%uint(p.count)]
}

func (p *pool) Compile(ps string) (*statedRegex, error) {
	return p.getShard(ps).getOrCompile(ps)
}
