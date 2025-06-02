// internal/storage/kvstore/kvstore.go
package kvstore

import (
	"sync"
	"time"
)

type Store struct {
	Mu      sync.RWMutex
	Data    map[string]string
	Expires map[string]time.Time
}

func NewStore() *Store {
	s := &Store{
		Data:    make(map[string]string),
		Expires: make(map[string]time.Time),
	}
	// 封装清理过期键的goroutine
	s.StartCleanup(10 * time.Second)
	return s
}

func (s *Store) Set(key, value string) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.Data[key] = value
}

func (s *Store) Get(key string) (string, bool) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	// Get时检查是否过期  （ 惰性删除 ）
	if expire, ok := s.Expires[key]; ok && time.Now().After(expire) {
		delete(s.Data, key)
		delete(s.Expires, key)
		return "", false
	}

	val, ok := s.Data[key]
	return val, ok
}

func (s *Store) SetWithExpire(key, value string, ttl time.Duration) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.Data[key] = value
	if ttl > 0 {
		s.Expires[key] = time.Now().Add(ttl)
	} else {
		delete(s.Expires, key)
	}
}

func (s *Store) Delete(key string) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	delete(s.Data, key)
	delete(s.Expires, key)
}

func (s *Store) Keys() []string {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	keys := make([]string, 0, len(s.Data))
	for k := range s.Data {
		keys = append(keys, k)
	}
	return keys
}

// 后台清理过期键的goroutine
func (s *Store) StartCleanup(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			s.cleanupExpired()
		}
	}()
}

// 最小堆过期键清理方案 ROI低 搁置
func (s *Store) cleanupExpired() {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	now := time.Now()
	for key, expire := range s.Expires {
		if now.After(expire) {
			delete(s.Data, key)
			delete(s.Expires, key)
		}
	}
}
