package model

import (
	"time"
)

type (
	// Batch represents the transactions batch object.
	Batch struct {
		ID           uint         `json:"id"`
		Transactions Transactions `json:"transactions"`
		MinValue     int          `json:"minimal_value"`
		Expires      time.Time    `json:"expires"`
		CreatedAt    time.Time    `json:"created_at"`
	}
	// Transactions represents a list of transaction objects.
	Transactions []Transaction
	// Transaction represents the transaction object.
	Transaction struct {
		UserID      uint      `json:"user_id"`
		Value       int       `json:"value"`
		Description string    `json:"description,omitempty"`
		CreatedAt   time.Time `json:"created_at"`
	}
)

// FirstBatch returns the first one batch.
func FirstBatch(minValue int, duration time.Duration) Batch {
	return NewBatch(0, minValue, duration)
}

// NewBatch returns a new batch object based in the last batch ID.
func NewBatch(lastID uint, minValue int, duration time.Duration) Batch {
	var expires time.Time
	if duration > 0 {
		expires = time.Now().Add(duration)
	}
	return Batch{
		ID:           lastID + 1,
		Transactions: make(Transactions, 0),
		MinValue:     minValue,
		Expires:      expires,
		CreatedAt:    time.Now(),
	}
}

// IsOutdated verify if the batch is outdated.
func (b *Batch) IsOutdated() bool {
	if (b.Expires == time.Time{} || len(b.Transactions) == 0) {
		return false
	}
	return b.Expires.Before(time.Now())
}

// Value sum all transaction values from batch.
func (b *Batch) Value() int {
	value := 0
	for _, tx := range b.Transactions {
		value += tx.Value
	}
	return value
}

// CheckLimit verify if the batch reach the minimal limit.
func (b *Batch) CheckLimit() bool {
	return b.Value() > b.MinValue
}

// CheckAvailability verify if the batch must be broadcasted.
func (b *Batch) CheckAvailability() bool {
	return b.CheckLimit() || b.IsOutdated()
}
