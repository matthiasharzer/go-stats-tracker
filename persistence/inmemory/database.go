package inmemory

import (
	"sync"

	"github.com/matthiasharzer/go-stats-tracker/persistence"
)

type Database struct {
	mu *sync.RWMutex
	m  map[string]persistence.UserContext
}

func NewDatabase() persistence.Database {
	return &Database{
		mu: &sync.RWMutex{},
		m:  make(map[string]persistence.UserContext),
	}
}

func (d *Database) Lookup(key string) *persistence.UserContext {
	d.mu.RLock()
	defer d.mu.RUnlock()
	userContext, ok := d.m[key]
	if !ok {
		return nil
	}
	return &userContext
}

func (d *Database) Save(key string, value persistence.UserContext) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.m[key] = value
	return nil
}
