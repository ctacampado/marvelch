package kvstore

import (
	"errors"
	"log"
)

// Store is a key-value store
type Store struct {
	Data map[string]interface{} `json:"Data"`
}

// New returns a new store
func New() *Store {
	return &Store{Data: make(map[string]interface{})}
}

// Get returns element by key or keyword
func (s *Store) Get(k string) interface{} {
	// returns cached data
	if k == "ALL" {
		return s.Data
	}

	// return value with given key
	v, ok := s.Data[k]
	if ok {
		return v
	}
	return nil
}

// Delete removes elem from store
func (s *Store) Delete(k string) error {
	if s.Get(k) != nil {
		delete(s.Data, k)
		return nil
	}
	return errors.New("key cannot be blank")
}

// Set sets the content of store located in data[k]
func (s *Store) Set(k string, n interface{}) error {
	if k == "" {
		return errors.New("key cannot be blank")
	}
	s.Data[k] = n
	return nil
}

// Print displays the content of store
func (s *Store) Print() {
	for k, v := range s.Data {
		log.Printf("k: %s | v:", k)
		log.Printf("%+v\n", v)
	}
}
