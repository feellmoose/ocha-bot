package helper

import (
	"errors"
	"sync"
)

type GenID interface {
	NextID() (string, error)
	WithID(func(string) error) error
}

type GenRandomRepoShortID struct {
	locks  sync.Map
	repo   Repo
	MinLen int
	MaxLen int
	Retry  int
}

func NewGenRandomRepoShortID(minLen, maxLen, retry int, repo Repo) *GenRandomRepoShortID {
	return &GenRandomRepoShortID{
		locks:  sync.Map{},
		repo:   repo,
		MaxLen: maxLen,
		MinLen: minLen,
		Retry:  retry,
	}
}

func (g *GenRandomRepoShortID) NextID() (string, error) {
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

func (g *GenRandomRepoShortID) WithID(handle func(string) error) error {
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
