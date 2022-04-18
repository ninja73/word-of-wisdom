package cache

import (
	"sync"
	"time"
)

type item struct {
	key       uint64
	timestamp int64
}

type inMemoryCache struct {
	queue []item
	data  map[uint64]struct{}
	ttl   int64
	lock  sync.Mutex
}

func NewInMemoryCache(cleanupInterval time.Duration, ttl int64) *inMemoryCache {
	queue := make([]item, 0)
	data := make(map[uint64]struct{})

	cache := &inMemoryCache{
		queue: queue,
		data:  data,
		ttl:   ttl,
	}

	go cache.cleanupLoop(cleanupInterval)

	return cache
}

func (s *inMemoryCache) ContainsOrAdd(val uint64) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	_, ok := s.data[val]
	if ok {
		return ok
	}

	s.queue = append(s.queue, item{
		key:       val,
		timestamp: time.Now().Unix(),
	})

	s.data[val] = struct{}{}

	return false
}

func (s *inMemoryCache) cleanupLoop(cleanupInterval time.Duration) {
	ticker := time.NewTicker(cleanupInterval)
	for range ticker.C {
		s.cleanup()
	}
}

func (s *inMemoryCache) cleanup() {
	now := time.Now().Unix()
	for {
		s.lock.Lock()
		if len(s.queue) == 0 {
			s.lock.Unlock()
			return
		}

		i := s.queue[0]
		if i.timestamp+s.ttl > now {
			s.lock.Unlock()
			return
		}

		delete(s.data, i.key)
		s.queue = s.queue[1:]
		s.lock.Unlock()
	}
}
