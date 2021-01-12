package worker

import (
	"time"

	"github.com/Pantani/batch/internal/db/database"
	"github.com/Pantani/batch/internal/model"
)

type (
	// Worker represents the batch transaction worker object.
	Worker struct {
		duration time.Duration
		db       database.IDatabase
		minValue int
	}
)

// New returns a new Worker object.
func New(db database.IDatabase, duration time.Duration, minValue int) *Worker {
	return &Worker{
		db:       db,
		duration: duration,
		minValue: minValue,
	}
}

// Run run the routine to verify pending kyc users.
// It returns an error if occurs.
func (w *Worker) Run() func() error {
	return func() error {
		ticker := time.NewTicker(w.duration)
		for {
			select {
			case <-ticker.C:
				err := w.sendTransactions()
				if err != nil {
					return err
				}
			}
		}
	}
}

// sendTransactions send all available transactions and create a new batch if needed.
// It returns an error if occurs.
func (w *Worker) sendTransactions() error {
	ok, err := w.db.IsEmpty()
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	b, err := w.db.GetPendingBatch()
	if err != nil {
		return err
	}

	if b.CheckAvailability() {
		// TODO send transaction

		b = model.NewBatch(b.ID, w.minValue, w.duration)
		return w.db.SaveBatch(b)
	}
	return nil
}
