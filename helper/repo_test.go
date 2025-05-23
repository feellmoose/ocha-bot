package helper

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestMemRepo(t *testing.T) {

}

func TestFileRepo(t *testing.T) {
	repo := NewFileRepo("", "test")
	defer repo.Stop()
	repo.Put("test", "test")
	repo.Put("1", 1)
	inter, ok := repo.Get("1")
	assert.Equal(t, ok, true)
	i := inter.(int)
	assert.Equal(t, i, 1)
	repo.Sync()
}
