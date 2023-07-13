package cache

import (
	"sync"
	"time"

	"github.com/amar-jay/minicache/errors"
	"github.com/amar-jay/minicache/logger"
)

// Cache is a simple in-memory cache
type Cache interface {
	Set(k []byte, v []byte, ttl time.Duration) error
	Has(k []byte) bool
	Get(k []byte) ([]byte, error)
	Delete(k []byte) error
}

type Cacher struct {
	data map[string][]byte
	sync.RWMutex
}

var _ Cache = (*Cacher)(nil)

func New() *Cacher {
	return &Cacher{
		data: make(map[string][]byte),
	}
}

func (c *Cacher) Set(k []byte, v []byte, ttl time.Duration) error {
	c.Lock()
	c.data[string(k)] = v
	c.Unlock()

	if ttl > 0 {
		go func() {
			<-time.After(ttl)
			c.Lock()
			delete(c.data, string(k))
			c.Unlock()
		}()
	}
	return nil
}

func (c *Cacher) Has(k []byte) bool {
	c.Lock()
	_, ok := c.data[string(k)]
	c.Unlock()
	return ok
}

func (c *Cacher) Get(k []byte) ([]byte, error) {
	c.Lock()
	defer c.Unlock()
	v, ok := c.data[string(k)]
	if !ok {
		return nil, logger.Errorf(errors.NotFound, k)
	}

	return v, nil
}

func (c *Cacher) Delete(k []byte) error {
	c.Lock()
	defer c.Unlock()
	delete(c.data, string(k))
	return nil
}
