package helper

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type GenID interface {
	NextID() (string, error)
	WithID(func(string) error) error
}

type NanoTimeID struct {
}

func (n *NanoTimeID) NextID() (string, error) {
	return fmt.Sprintf("%d", time.Now().UnixNano()), nil
}

func (n *NanoTimeID) WithID(handle func(string) error) error {
	id, err := n.NextID()
	if err != nil {
		return err
	}
	return handle(id)
}

type GenRandomRepoShortID[T any] struct {
	locks  sync.Map
	repo   Repo[T]
	MinLen int
	MaxLen int
	Retry  int
}

func NewGenRandomRepoShortID[T any](minLen, maxLen, retry int, repo Repo[T]) *GenRandomRepoShortID[T] {
	return &GenRandomRepoShortID[T]{
		locks:  sync.Map{},
		repo:   repo,
		MaxLen: maxLen,
		MinLen: minLen,
		Retry:  retry,
	}
}

func (g *GenRandomRepoShortID[T]) NextID() (string, error) {
	for i := g.MinLen; i <= g.MaxLen; i++ {
		draft := RandomString(i)
		if _, ok := g.repo.Get(draft); !ok {
			return draft, nil
		}
	}
	for i := 0; i < g.Retry; i++ {
		draft := RandomString(i)
		if _, ok := g.repo.Get(draft); !ok {
			return draft, nil
		}
	}
	return "", errors.New("ID gen failed")
}

func (g *GenRandomRepoShortID[T]) WithID(handle func(string) error) error {
	id, err := g.NextID()
	if err != nil {
		return err
	}
	var mu *sync.Mutex
	val, _ := g.locks.LoadOrStore(id, &sync.Mutex{})
	mu = val.(*sync.Mutex)
	if mu.TryLock() {
		defer mu.Unlock()
		defer g.locks.Delete(id)
		return handle(id)
	}
	return errors.New("ID gen failed")
}
