package redis

import (
	"context"
	"strconv"

	"github.com/Pantani/batch/internal/db/database"
	"github.com/Pantani/batch/internal/model"

	"github.com/Pantani/redis"
)

// Database represents the database object.
type Database struct {
	client      *redis.Redis
	lastBatchID uint
}

// New initialize the Redis database object.
// It returns a database object and an error if occurs.
func New(host, password string, dbIndex int) (database.IDatabase, error) {
	client, err := redis.New(host, password, dbIndex)
	if err != nil {
		return nil, err
	}
	return &Database{client: client}, nil
}

// SaveBatch update or create a batch transactions.
// It returns an error if occurs.
func (d *Database) SaveBatch(batch model.Batch) error {
	err := d.client.AddObject(context.TODO(), strconv.Itoa(int(batch.ID)), batch, 0)
	if err != nil {
		return err
	}
	if batch.ID > d.lastBatchID {
		d.lastBatchID = batch.ID
	}
	return nil
}

// GetBatch get the batch object by ID.
// It returns the batch model and an error if occurs.
func (d *Database) GetBatch(id string) (b model.Batch, err error) {
	err = d.client.GetObject(context.TODO(), id, &b)
	return
}

// GetPendingBatch get the current pending batch.
// It returns the batch model and an error if occurs.
func (d *Database) GetPendingBatch() (b model.Batch, err error) {
	err = d.client.GetObject(context.TODO(), strconv.Itoa(int(d.lastBatchID)), &b)
	return
}

// IsEmpty verify if the database is empty.
// It returns the batch model and an error if occurs.
func (d *Database) IsEmpty() (bool, error) {
	var b *model.Batch
	err := d.client.GetObject(context.TODO(), strconv.Itoa(int(d.lastBatchID)), &b)
	return b == nil, err
}
