package concurrentMap

import (
	"hash"
	"hash/fnv"
	"sync"
)

var SHARD_COUNT = 16

// A "thread" safe map of type string:Anything.
// To avoid lock bottlenecks this map is dived to several (SHARD_COUNT) map shards.
type ConcurrentMap struct {
	shards []*ConcurrentMapShared
	hasher func() hash.Hash32
}

// A "thread" safe string to anything map.
type ConcurrentMapShared struct {
	items        map[string]interface{}
	sync.RWMutex // Read Write mutex, guards access to internal map.
}

// Creates a new concurrent map.
func New() *ConcurrentMap {
	return AdvanceNew(fnv.New32, SHARD_COUNT)
}

func AdvanceNew(hasher func() hash.Hash32, sc int) *ConcurrentMap {
	if sc == 0 {
		sc = SHARD_COUNT
	}
	if hasher == nil {
		hasher = fnv.New32
	}
	m := &ConcurrentMap{make([]*ConcurrentMapShared, sc), hasher}
	for i := 0; i < sc; i++ {
		m.shards[i] = &ConcurrentMapShared{items: make(map[string]interface{})}
	}
	return m
}

// Returns shard under given key
func (m *ConcurrentMap) GetShard(key string) *ConcurrentMapShared {
	hasher := m.hasher()
	hasher.Write([]byte(key))
	return m.shards[uint(hasher.Sum32())%uint(SHARD_COUNT)]
}

func (m *ConcurrentMap) MSet(data map[string]interface{}) {
	for key, value := range data {
		shard := m.GetShard(key)
		shard.Lock()
		shard.items[key] = value
		shard.Unlock()
	}
}

// Sets the given value under the specified key.
func (m *ConcurrentMap) Set(key string, value interface{}) {
	// Get map shard.
	shard := m.GetShard(key)
	shard.Lock()
	shard.items[key] = value
	shard.Unlock()
}

// Sets the given value under the specified key if no value was associated with it.
func (m *ConcurrentMap) SetIfAbsent(key string, value interface{}) bool {
	// Get map shard.
	shard := m.GetShard(key)
	shard.Lock()
	_, ok := shard.items[key]
	if !ok {
		shard.items[key] = value
	}
	shard.Unlock()
	return !ok
}

// Retrieves an element from map under given key.
func (m *ConcurrentMap) Get(key string) (interface{}, bool) {
	// Get shard
	shard := m.GetShard(key)
	shard.RLock()
	// Get item from shard.
	val, ok := shard.items[key]
	shard.RUnlock()
	return val, ok
}

func (m *ConcurrentMap) GetOrSet(key string, value interface{}) (interface{}, bool) {
	shard := m.GetShard(key)
	var created bool
	shard.Lock()
	val, ok := shard.items[key]
	if !ok {
		shard.items[key] = value
		val = value
		created = true
	}
	shard.Unlock()
	return val, created
}

func (m *ConcurrentMap) GetOrBlock(key string, handler func() map[string]interface{}) (interface{}, bool) {
	shard := m.GetShard(key)
	shard.Lock()
	val, ok := shard.items[key]
	if !ok && handler != nil {
		for k, v := range handler() {
			shard.items[k] = v
		}
		val, ok = shard.items[key]
	}
	defer shard.Unlock()
	return val, ok
}

// Returns the number of elements within the map.
func (m *ConcurrentMap) Count() int {
	count := 0
	for i := 0; i < SHARD_COUNT; i++ {
		shard := m.shards[i]
		shard.RLock()
		count += len(shard.items)
		shard.RUnlock()
	}
	return count
}

// Looks up an item under specified key
func (m *ConcurrentMap) Has(key string) bool {
	// Get shard
	shard := m.GetShard(key)
	shard.RLock()
	// See if element is within shard.
	_, ok := shard.items[key]
	shard.RUnlock()
	return ok
}

// Removes an element from the map.
func (m *ConcurrentMap) Remove(key string) bool {
	// Try to get shard.
	shard := m.GetShard(key)
	shard.Lock()
	defer shard.Unlock()
	_, ok := shard.items[key]
	delete(shard.items, key)
	return ok
}

// Checks if map is empty.
func (m *ConcurrentMap) IsEmpty() bool {
	return m.Count() == 0
}

// Used by the Iter & IterBuffered functions to wrap two variables together over a channel,
type Tuple struct {
	Key string
	Val interface{}
}

// Returns a buffered iterator which could be used in a for range loop.
func (m *ConcurrentMap) IterBuffered() <-chan Tuple {
	ch := make(chan Tuple, m.Count())
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(SHARD_COUNT)
		// Foreach shard.
		for _, shard := range m.shards {
			go func(shard *ConcurrentMapShared) {
				// Foreach key, value pair.
				shard.RLock()
				for key, val := range shard.items {
					ch <- Tuple{key, val}
				}
				shard.RUnlock()
				wg.Done()
			}(shard)
		}
		wg.Wait()
		close(ch)
	}()
	return ch
}

// Return all keys as []string
func (m *ConcurrentMap) Keys() []string {
	count := m.Count()
	ch := make(chan string, count)
	go func() {
		// Foreach shard.
		wg := sync.WaitGroup{}
		wg.Add(SHARD_COUNT)
		for _, shard := range m.shards {
			go func(shard *ConcurrentMapShared) {
				// Foreach key, value pair.
				shard.RLock()
				for key := range shard.items {
					ch <- key
				}
				shard.RUnlock()
				wg.Done()
			}(shard)
		}
		wg.Wait()
		close(ch)
	}()

	// Generate keys
	keys := make([]string, count)
	for i := 0; i < count; i++ {
		keys[i] = <-ch
	}
	return keys
}
