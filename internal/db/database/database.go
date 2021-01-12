package database

import (
	"github.com/Pantani/batch/internal/model"
)

// IDatabase represents the application database interfaces.
type IDatabase interface {
	SaveBatch(batch model.Batch) error
	GetBatch(id string) (model.Batch, error)
	GetPendingBatch() (model.Batch, error)
	IsEmpty() (bool, error)
}
