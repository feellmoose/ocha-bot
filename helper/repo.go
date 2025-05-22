package helper

import (
	"sync"
	"time"
)

type Repo interface {
	Get(key string) (interface{}, bool)
	Put(key string, value interface{}) bool
	Del(key string) bool
}

type Item struct {
	value      interface{}
	expiration time.Time
}

type MemRepo struct {
	data    sync.Map
	cleanup *time.Ticker
	ttl     time.Duration
	stop    chan struct{}
}

func NewMemRepo() *MemRepo {
	repo := &MemRepo{
		stop:    make(chan struct{}),
		cleanup: time.NewTicker(1 * time.Second),
		ttl:     time.Hour,
	}

	go repo.cleanupExpired()

	return repo
}

func (r *MemRepo) cleanupExpired() {
	for {
		select {
		case <-r.cleanup.C:
			r.data.Range(func(key, value interface{}) bool {
				item, ok := value.(Item)
				if !ok {
					return true
				}
				if time.Now().After(item.expiration) {
					r.Del(key.(string))
				}
				return true
			})
		case <-r.stop:
			return
		}
	}
}

func (r *MemRepo) Stop() {
	close(r.stop)
	r.cleanup.Stop()
}

func (r *MemRepo) Get(key string) (interface{}, bool) {
	item, exists := r.data.Load(key)
	if !exists {
		return nil, false
	}

	castedItem, ok := item.(Item)
	if !ok {
		return nil, false
	}

	if time.Now().After(castedItem.expiration) {
		r.Del(key)
		return nil, false
	}

	r.Put(key, castedItem.value)

	return castedItem.value, true
}

func (r *MemRepo) Put(key string, value interface{}) bool {
	expirationTime := time.Now().Add(r.ttl)

	r.data.Store(key, Item{
		value:      value,
		expiration: expirationTime,
	})

	return true
}

func (r *MemRepo) Del(key string) bool {
	_, exists := r.data.Load(key)
	if exists {
		r.data.Delete(key)
		return true
	}
	return false
}
