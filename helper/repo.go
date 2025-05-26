package helper

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type RepoInfo interface {
	Size() int
	Name() string
	Type() string
}

type FileRepoInfo interface {
	DataSize() int64
}

type Repo[T any] interface {
	Get(key string) (T, bool)
	Put(key string, value T) bool
	Del(key string) bool
	Range(f func(key string, value T) bool)
	Stop()
}

type RepoType string

const (
	Memory RepoType = "Memory"
	File   RepoType = "File"
	None   RepoType = "None"
)

type Item[T any] struct {
	value      T
	expiration time.Time
}

type MemRepo[T any] struct {
	name    string
	data    sync.Map
	cleanup *time.Ticker
	ttl     time.Duration
	stop    chan struct{}
}

func NewMemRepo[T any](name string) *MemRepo[T] {
	repo := &MemRepo[T]{
		name:    name,
		stop:    make(chan struct{}),
		cleanup: time.NewTicker(1 * time.Second),
		ttl:     time.Hour,
	}

	go repo.cleanupExpired()

	return repo
}

func (r *MemRepo[T]) cleanupExpired() {
	for {
		select {
		case <-r.cleanup.C:
			r.data.Range(func(key, value any) bool {
				item, ok := value.(Item[T])
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

func (r *MemRepo[T]) Stop() {
	close(r.stop)
	r.cleanup.Stop()
}

func (r *MemRepo[T]) Type() string {
	return string(Memory)
}

func (r *MemRepo[T]) Name() string {
	return r.name
}

func (r *MemRepo[T]) Size() int {
	count := 0
	r.data.Range(func(key, value any) bool {
		item, ok := value.(Item[T])
		if !ok {
			return true
		}
		if time.Now().After(item.expiration) {
			r.Del(key.(string))
		} else {
			count++
		}
		return true
	})
	return count
}

func (r *MemRepo[T]) Range(f func(key string, value T) bool) {
	r.data.Range(func(k, v any) bool {
		item, ok := v.(Item[T])
		if !ok {
			return true
		}
		if time.Now().After(item.expiration) {
			r.Del(k.(string))
		} else {
			return f(k.(string), item.value)
		}
		return true
	})
}

func (r *MemRepo[T]) Get(key string) (T, bool) {
	item, exists := r.data.Load(key)
	if !exists {
		var zero T
		return zero, false
	}

	castedItem, ok := item.(Item[T])
	if !ok {
		var zero T
		return zero, false
	}

	if time.Now().After(castedItem.expiration) {
		r.Del(key)
		var zero T
		return zero, false
	}

	r.Put(key, castedItem.value)

	return castedItem.value, true
}

func (r *MemRepo[T]) Put(key string, value T) bool {
	expirationTime := time.Now().Add(r.ttl)

	r.data.Store(key, Item[T]{
		value:      value,
		expiration: expirationTime,
	})

	return true
}

func (r *MemRepo[T]) Del(key string) bool {
	_, exists := r.data.Load(key)
	if exists {
		r.data.Delete(key)
		return true
	}
	return false
}

type FileRepo[T any] struct {
	dir      string
	name     string
	filename string
	data     *sync.Map
	sticker  *time.Ticker
	stop     chan struct{}
}

func NewFileRepo[T any](dir string, name string) *FileRepo[T] {
	var (
		fn       = name + "_" + strconv.FormatInt(BotID, 10) + "_data.json"
		filename = filepath.Join(dir, fn)
		file     *os.File
		err      error
		data     sync.Map
	)

	if _, err = os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		file, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			log.Panicf("Error create file repo: file error=%v", err)
		}

		encoder := json.NewEncoder(file)
		if err = encoder.Encode(make(map[string]T)); err != nil {
			log.Panicf("Error create file repo: encode {} in create file error=%v", err)
		}

		if err = file.Sync(); err != nil {
			log.Panicf("Error create file repo: sync file error=%v", err)
		}

		log.Printf("Create File(filename=%s) with {} success", filename)
	} else {
		file, err = os.OpenFile(filename, os.O_RDWR, 0666)
		if err != nil {
			log.Panicf("Error create file repo: file error=%v", err)
		}

		var temp map[string]T
		decoder := json.NewDecoder(file)
		if err = decoder.Decode(&temp); err != nil {
			log.Panicf("Error create file repo: decode failed error=%v", err)
		}

		for key, value := range temp {
			data.Store(key, value)
		}
		log.Printf("Load File(filename=%s) data success", filename)
	}

	repo := &FileRepo[T]{
		dir:      dir,
		name:     name,
		filename: fn,
		data:     &data,
		stop:     make(chan struct{}),
		sticker:  time.NewTicker(5 * time.Second),
	}

	go repo.loop()

	return repo
}

func (r *FileRepo[T]) loop() {
	for {
		select {
		case <-r.sticker.C:
			err := r.Sync()
			if err != nil {
				log.Printf("Sync repo into file failed %v", err)
			}
		case <-r.stop:
			return
		}
	}
}

func (r *FileRepo[T]) Sync() error {
	temp, err := os.CreateTemp(r.dir, "repo_*.json")
	if err != nil {
		return err
	}
	defer temp.Close()
	defer os.Remove(temp.Name())
	encoder := json.NewEncoder(temp)

	tempMap := make(map[string]T)
	r.data.Range(func(key, value any) bool {
		strKey, ok := key.(string)
		if !ok {
			return true
		}
		tv, ok := value.(T)
		if !ok {
			return true
		}
		tempMap[strKey] = tv
		return true
	})

	if err = encoder.Encode(tempMap); err != nil {
		return err
	}
	if err = temp.Sync(); err != nil {
		return err
	}

	dst, err := os.Create(filepath.Join(r.dir, r.filename))
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = temp.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	_, err = io.Copy(dst, temp)
	if err != nil {
		return err
	}

	return nil
}

func (r *FileRepo[T]) Type() string {
	return string(File)
}

func (r *FileRepo[T]) Name() string {
	return r.name
}

func (r *FileRepo[T]) DataSize() int64 {
	fileInfo, err := os.Stat(filepath.Join(r.dir, r.filename))
	if err != nil {
		return 0
	}
	return fileInfo.Size()
}

func (r *FileRepo[T]) Size() int {
	count := 0
	r.data.Range(func(_, _ any) bool {
		count++
		return true
	})
	return count
}

func (r *FileRepo[T]) Range(f func(key string, value T) bool) {
	r.data.Range(func(k, v any) bool {
		return f(k.(string), v.(T))
	})
}

func (r *FileRepo[T]) Stop() {
	close(r.stop)
	err := r.Sync()
	if err != nil {
		log.Printf("Sync repo into file failed %v", err)
	}
	if err != nil {
		log.Printf("Sync repo into file failed %v", err)
	}
}

func (r *FileRepo[T]) Get(key string) (T, bool) {
	v, ok := r.data.Load(key)
	return v.(T), ok
}

func (r *FileRepo[T]) Put(key string, value T) bool {
	r.data.Store(key, value)
	return true
}

func (r *FileRepo[T]) Del(key string) bool {
	_, exists := r.data.Load(key)
	if exists {
		r.data.Delete(key)
		return true
	}
	return false
}
