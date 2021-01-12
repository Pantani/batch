package memory

import (
	"strconv"
	"sync"

	"github.com/Pantani/batch/internal/db/database"
	"github.com/Pantani/batch/internal/model"

	"github.com/Pantani/errors"
)

// Database represents the database object.
type Database struct {
	mu          sync.RWMutex
	cache       map[string]model.Batch
	lastBatchID uint
}

// New initialize the in memory database object.
// It returns a database object and an error if occurs.
func New() database.IDatabase {
	return &Database{
		cache:       make(map[string]model.Batch),
		lastBatchID: 0,
	}
}

// SaveBatch update or create a batch transactions.
// It returns an error if occurs.
func (d *Database) SaveBatch(batch model.Batch) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.cache[strconv.Itoa(int(batch.ID))] = batch
	if batch.ID > d.lastBatchID {
		d.lastBatchID = batch.ID
	}
	return nil
}

// GetBatch get the batch object by ID.
// It returns the batch model and an error if occurs.
func (d *Database) GetBatch(id string) (model.Batch, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.cache[id], nil
}

// GetPendingBatch get the current pending batch.
// It returns the batch model and an error if occurs.
func (d *Database) GetPendingBatch() (model.Batch, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	b, ok := d.cache[strconv.Itoa(len(d.cache))]
	if !ok {
		return model.Batch{}, errors.E("no pending batch available")
	}
	return b, nil
}

// IsEmpty verify if the database is empty.
// It returns the batch model and an error if occurs.
func (d *Database) IsEmpty() (bool, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return len(d.cache) == 0, nil
}
