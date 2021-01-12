package memory

import (
	"testing"
	"time"

	"github.com/Pantani/batch/internal/model"

	"github.com/stretchr/testify/assert"
)

var (
	tx1 = model.Transaction{
		UserID:      10,
		Value:       10,
		Description: "tx 1",
		CreatedAt:   time.Time{},
	}
	tx2 = model.Transaction{
		UserID:      10,
		Value:       50,
		Description: "tx 2",
		CreatedAt:   time.Time{},
	}
	tx3 = model.Transaction{
		UserID:      10,
		Value:       55,
		Description: "tx 3",
		CreatedAt:   time.Time{},
	}
)

func TestDatabase(t *testing.T) {
	batch1 := model.Batch{
		ID:           1,
		Transactions: model.Transactions{tx1},
		MinValue:     100,
		Expires:      time.Time{},
		CreatedAt:    time.Time{}.Add(10),
	}
	batch2 := model.Batch{
		ID:           2,
		Transactions: model.Transactions{tx1, tx2},
		MinValue:     100,
		Expires:      time.Time{}.Add(10),
		CreatedAt:    time.Time{}.Add(10),
	}
	batch3 := model.Batch{
		ID:           3,
		Transactions: model.Transactions{tx2, tx3},
		MinValue:     100,
		Expires:      time.Unix(1<<63-1, 0),
		CreatedAt:    time.Time{}.Add(10),
	}
	db := New()
	assert.NotNil(t, db)

	t.Run("verify if is empty", func(t *testing.T) {
		empty, err := db.IsEmpty()
		assert.Nil(t, err)
		assert.True(t, empty)
	})

	t.Run("save batch 1", func(t *testing.T) {
		err := db.SaveBatch(batch1)
		assert.Nil(t, err)

		empty, err := db.IsEmpty()
		assert.Nil(t, err)
		assert.False(t, empty)

		pending, err := db.GetPendingBatch()
		assert.Nil(t, err)
		assert.EqualValues(t, batch1, pending)
	})

	t.Run("save batch 2", func(t *testing.T) {
		err := db.SaveBatch(batch2)
		assert.Nil(t, err)

		pending, err := db.GetPendingBatch()
		assert.Nil(t, err)
		assert.EqualValues(t, batch2, pending)
	})

	t.Run("save batch 3", func(t *testing.T) {
		err := db.SaveBatch(batch3)
		assert.Nil(t, err)

		pending, err := db.GetPendingBatch()
		assert.Nil(t, err)
		assert.EqualValues(t, batch3, pending)
	})

	t.Run("get batch by id", func(t *testing.T) {
		got1, err := db.GetBatch("1")
		assert.Nil(t, err)
		assert.EqualValues(t, batch1, got1)
		got2, err := db.GetBatch("2")
		assert.Nil(t, err)
		assert.EqualValues(t, batch2, got2)
		got3, err := db.GetBatch("3")
		assert.Nil(t, err)
		assert.EqualValues(t, batch3, got3)
	})
}
