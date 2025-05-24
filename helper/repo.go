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

type FileRepo struct {
	dir     string
	name    string
	data    *sync.Map
	sticker *time.Ticker
	stop    chan struct{}
}

func NewFileRepo(dir string, name string) *FileRepo {
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
		if err = encoder.Encode(make(map[string]interface{})); err != nil {
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

		var temp map[string]interface{}
		decoder := json.NewDecoder(file)
		if err = decoder.Decode(&temp); err != nil {
			log.Panicf("Error create file repo: decode failed error=%v", err)
		}

		for key, value := range temp {
			data.Store(key, value)
		}
		log.Printf("Load File(filename=%s) data success", filename)
	}

	repo := &FileRepo{
		dir:     dir,
		name:    fn,
		data:    &data,
		stop:    make(chan struct{}),
		sticker: time.NewTicker(5 * time.Second),
	}

	go repo.loop()

	return repo
}

func (r *FileRepo) loop() {
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

func (r *FileRepo) Sync() error {
	temp, err := os.CreateTemp(r.dir, "repo_*.json")
	if err != nil {
		return err
	}
	defer temp.Close()
	defer os.Remove(temp.Name())
	encoder := json.NewEncoder(temp)

	tempMap := make(map[string]interface{})
	r.data.Range(func(key, value interface{}) bool {
		strKey, ok := key.(string)
		if !ok {
			return true
		}
		tempMap[strKey] = value
		return true
	})

	if err = encoder.Encode(tempMap); err != nil {
		return err
	}
	if err = temp.Sync(); err != nil {
		return err
	}

	dst, err := os.Create(filepath.Join(r.dir, r.name))
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

func (r *FileRepo) Stop() {
	close(r.stop)
	err := r.Sync()
	if err != nil {
		log.Printf("Sync repo into file failed %v", err)
	}
	if err != nil {
		log.Printf("Sync repo into file failed %v", err)
	}
}

func (r *FileRepo) Get(key string) (interface{}, bool) {
	return r.data.Load(key)
}

func (r *FileRepo) Put(key string, value interface{}) bool {
	r.data.Store(key, value)
	return true
}

func (r *FileRepo) Del(key string) bool {
	_, exists := r.data.Load(key)
	if exists {
		r.data.Delete(key)
		return true
	}
	return false
}
