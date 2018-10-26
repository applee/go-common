package cmap

import (
	"encoding/json"
	"hash/fnv"
	"sync"
)

// DefaultShardCount represents the default count of shard.
var DefaultShardCount = 32

// ConcurrentMap represents a "thread" safe map of type string:Anything.
// To avoid lock bottlenecks this map is dived to several (DefaultShardCount) map shards.
type ConcurrentMap struct {
	shards []*ConcurrentMapShared
	cnt    uint
}

// ConcurrentMapShared represents a "thread" safe string to anything map.
type ConcurrentMapShared struct {
	items        map[string]interface{}
	sync.RWMutex // Read Write mutex, guards access to internal map.
}

// New creates a new concurrent map.
func New(shardCnt uint) ConcurrentMap {
	if shardCnt == 0 {
		shardCnt = uint(DefaultShardCount)
	}
	m := ConcurrentMap{
		shards: make([]*ConcurrentMapShared, shardCnt),
		cnt:    shardCnt,
	}
	for i := 0; i < DefaultShardCount; i++ {
		m.shards[i] = &ConcurrentMapShared{items: make(map[string]interface{})}
	}
	return m
}

// GetShard returns shard under given key
func (m ConcurrentMap) GetShard(key string) *ConcurrentMapShared {
	hasher := fnv.New32()
	hasher.Write([]byte(key))
	return m.shards[uint(hasher.Sum32())%m.cnt]
}

// Set sets the given value under the specified key.
func (m *ConcurrentMap) Set(key string, value interface{}) {
	// Get map shard.
	shard := m.GetShard(key)
	shard.Lock()
	defer shard.Unlock()
	shard.items[key] = value
}

// SetIfAbsent sets the given value under the specified key if no value was associated with it.
func (m *ConcurrentMap) SetIfAbsent(key string, value interface{}) bool {
	// Get map shard.
	shard := m.GetShard(key)
	shard.Lock()
	defer shard.Unlock()
	_, ok := shard.items[key]
	if !ok {
		shard.items[key] = value
	}
	return !ok
}

// Get retrieves an element from map under given key.
func (m ConcurrentMap) Get(key string) (interface{}, bool) {
	// Get shard
	shard := m.GetShard(key)
	shard.RLock()
	defer shard.RUnlock()

	// Get item from shard.
	val, ok := shard.items[key]
	return val, ok
}

// Count returns the number of elements within the map.
func (m ConcurrentMap) Count() int {
	count := 0
	for i := 0; i < int(m.cnt); i++ {
		shard := m.shards[i]
		shard.RLock()
		count += len(shard.items)
		shard.RUnlock()
	}
	return count
}

// Has looks up an item under specified key
func (m *ConcurrentMap) Has(key string) bool {
	// Get shard
	shard := m.GetShard(key)
	shard.RLock()
	defer shard.RUnlock()

	// See if element is within shard.
	_, ok := shard.items[key]
	return ok
}

// Remove removes an element from the map.
func (m *ConcurrentMap) Remove(key string) {
	// Try to get shard.
	shard := m.GetShard(key)
	shard.Lock()
	defer shard.Unlock()
	delete(shard.items, key)
}

// IsEmpty checks if map is empty.
func (m *ConcurrentMap) IsEmpty() bool {
	return m.Count() == 0
}

// Tuple used by the Iter & IterBuffered functions to wrap two variables
// together over a channel
type Tuple struct {
	Key string
	Val interface{}
}

// Iter returns an iterator which could be used in a for range loop.
func (m ConcurrentMap) Iter() <-chan Tuple {
	ch := make(chan Tuple)
	go func() {
		// Foreach shard.
		for _, shard := range m.shards {
			// Foreach key, value pair.
			shard.RLock()
			for key, val := range shard.items {
				ch <- Tuple{key, val}
			}
			shard.RUnlock()
		}
		close(ch)
	}()
	return ch
}

// IterBuffered returns a buffered iterator which could be used in a for range loop.
func (m ConcurrentMap) IterBuffered() <-chan Tuple {
	ch := make(chan Tuple, m.Count())
	go func() {
		// Foreach shard.
		for _, shard := range m.shards {
			// Foreach key, value pair.
			shard.RLock()
			for key, val := range shard.items {
				ch <- Tuple{key, val}
			}
			shard.RUnlock()
		}
		close(ch)
	}()
	return ch
}

// Items returns all items as map[string]interface{}
func (m ConcurrentMap) Items() map[string]interface{} {
	tmp := make(map[string]interface{})

	// Insert items to temporary map.
	for item := range m.Iter() {
		tmp[item.Key] = item.Val
	}

	return tmp
}

// MarshalJSON implements the Marshaler interface.
func (m ConcurrentMap) MarshalJSON() ([]byte, error) {
	// Create a temporary map, which will hold all item spread across shards.
	tmp := make(map[string]interface{})

	// Insert items to temporary map.
	for item := range m.Iter() {
		tmp[item.Key] = item.Val
	}
	return json.Marshal(tmp)
}
