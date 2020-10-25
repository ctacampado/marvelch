package service

import (
	"log"
	"marvel-chars/internal/kvstore"
)

// Cache interface
type Cache interface {
	Print()
	Set(string, interface{}) error
	Get(string) interface{}
	Delete(string) error
}

// SvcCache is a cache implementation
type SvcCache struct {
	Cache Cache
}

// InitCache is a wrapper method for the InitFunc
func (c *SvcCache) InitCache(f func(Cache) error) {
	c.Cache = kvstore.New()
	err := f(c.Cache)
	if err != nil {
		log.Fatal(err)
	}
}
