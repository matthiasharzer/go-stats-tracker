package inmemory

import (
	"sync"

	"github.com/matthiasharzer/go-stats-tracker/persistence"
)

type Database struct {
	mu *sync.RWMutex
	m  map[string]persistence.SheetContext
}

func NewDatabase() persistence.Database {
	return &Database{
		mu: &sync.RWMutex{},
		m:  make(map[string]persistence.SheetContext),
	}
}

func (d *Database) Lookup(accessKey string) (*persistence.SheetContext, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	sheetContext, ok := d.m[accessKey]
	if !ok {
		return nil, nil
	}
	return &sheetContext, nil
}

func (d *Database) Save(accessKey string, value persistence.SheetContext) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.m[accessKey] = value
	return nil
}

func (d *Database) Exists(accessKey string) (bool, error) {
	sheetContext, err := d.Lookup(accessKey)
	if err != nil {
		return false, err
	}
	return sheetContext != nil, nil
}
