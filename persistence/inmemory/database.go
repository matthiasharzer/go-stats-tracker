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

func (d *Database) Lookup(userID string) (*persistence.UserContext, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	userContext, ok := d.m[userID]
	if !ok {
		return nil, nil
	}
	return &userContext, nil
}

func (d *Database) Save(userID string, value persistence.UserContext) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.m[userID] = value
	return nil
}

func (d *Database) Exists(userID string) (bool, error) {
	userContext, err := d.Lookup(userID)
	if err != nil {
		return false, err
	}
	return userContext != nil, nil
}
