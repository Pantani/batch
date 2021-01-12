package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	tx1 = Transaction{
		UserID:      10,
		Value:       10,
		Description: "tx 1",
		CreatedAt:   time.Time{},
	}
	tx2 = Transaction{
		UserID:      10,
		Value:       50,
		Description: "tx 2",
		CreatedAt:   time.Time{},
	}
	tx3 = Transaction{
		UserID:      10,
		Value:       55,
		Description: "tx 3",
		CreatedAt:   time.Time{},
	}
	tx4 = Transaction{
		UserID:      10,
		Value:       99,
		Description: "tx 4",
		CreatedAt:   time.Time{},
	}
	tx5 = Transaction{
		UserID:      10,
		Value:       200,
		Description: "tx 5",
		CreatedAt:   time.Time{},
	}

	batch1 = Batch{
		ID:           1,
		Transactions: Transactions{tx1},
		MinValue:     100,
		Expires:      time.Time{},
		CreatedAt:    time.Time{}.Add(10),
	}
	batch2 = Batch{
		ID:           2,
		Transactions: Transactions{tx1, tx2},
		MinValue:     100,
		Expires:      time.Time{}.Add(10),
		CreatedAt:    time.Time{}.Add(10),
	}
	batch3 = Batch{
		ID:           3,
		Transactions: Transactions{tx2, tx3},
		MinValue:     100,
		Expires:      time.Unix(1<<63-1, 0),
		CreatedAt:    time.Time{}.Add(10),
	}
	batch4 = Batch{
		ID:           4,
		Transactions: Transactions{tx1, tx4},
		MinValue:     100,
		Expires:      time.Now().Add(1),
		CreatedAt:    time.Time{}.Add(100),
	}
	batch5 = Batch{
		ID:           5,
		Transactions: Transactions{tx5},
		MinValue:     100,
		Expires:      time.Now().Add(100000000),
		CreatedAt:    time.Time{}.Add(10),
	}
	batch6 = Batch{
		ID:           5,
		Transactions: Transactions{tx1},
		MinValue:     100,
		Expires:      time.Now().Add(100000000),
		CreatedAt:    time.Time{}.Add(10),
	}
)

func TestBatch_CheckAvailability(t *testing.T) {
	tests := []struct {
		name  string
		batch Batch
		want  bool
	}{
		{name: "test batch 1", batch: batch1, want: false},
		{name: "test batch 2", batch: batch2, want: true},
		{name: "test batch 3", batch: batch3, want: true},
		{name: "test batch 4", batch: batch4, want: true},
		{name: "test batch 5", batch: batch5, want: true},
		{name: "test batch 6", batch: batch6, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.batch.CheckAvailability()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBatch_CheckLimit(t *testing.T) {
	tests := []struct {
		name  string
		batch Batch
		want  bool
	}{
		{name: "test batch 1", batch: batch1, want: false},
		{name: "test batch 2", batch: batch2, want: false},
		{name: "test batch 3", batch: batch3, want: true},
		{name: "test batch 4", batch: batch4, want: true},
		{name: "test batch 5", batch: batch5, want: true},
		{name: "test batch 6", batch: batch6, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.batch.CheckLimit()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBatch_IsOutdated(t *testing.T) {
	tests := []struct {
		name  string
		batch Batch
		want  bool
	}{
		{name: "test batch 1", batch: batch1, want: false},
		{name: "test batch 2", batch: batch2, want: true},
		{name: "test batch 3", batch: batch3, want: true},
		{name: "test batch 4", batch: batch4, want: true},
		{name: "test batch 5", batch: batch5, want: false},
		{name: "test batch 6", batch: batch6, want: false},
		{name: "test empty batch", batch: Batch{Transactions: Transactions{}}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.batch.IsOutdated()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBatch_Value(t *testing.T) {
	tests := []struct {
		name  string
		batch Batch
		want  int
	}{
		{name: "test batch 1", batch: batch1, want: 10},
		{name: "test batch 2", batch: batch2, want: 60},
		{name: "test batch 3", batch: batch3, want: 105},
		{name: "test batch 4", batch: batch4, want: 109},
		{name: "test batch 5", batch: batch5, want: 200},
		{name: "test batch 6", batch: batch6, want: 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.batch.Value()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFirstBatch(t *testing.T) {
	type args struct {
		minValue int
		duration time.Duration
	}
	tests := []struct {
		name string
		args args
		want Batch
	}{
		{
			name: "test batch 1",
			args: args{minValue: 100, duration: 0},
			want: Batch{ID: 1, Transactions: Transactions{}, MinValue: 100},
		}, {
			name: "test batch 2",
			args: args{minValue: 100, duration: 1},
			want: Batch{ID: 1, Transactions: Transactions{}, MinValue: 100},
		}, {
			name: "test batch 3",
			args: args{minValue: 100, duration: 10},
			want: Batch{ID: 1, Transactions: Transactions{}, MinValue: 100},
		}, {
			name: "test batch 4",
			args: args{minValue: 100, duration: 100},
			want: Batch{ID: 1, Transactions: Transactions{}, MinValue: 100},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FirstBatch(tt.args.minValue, tt.args.duration)
			assert.EqualValues(t, tt.want.Transactions, got.Transactions)
			assert.EqualValues(t, tt.want.ID, got.ID)
			assert.EqualValues(t, tt.want.Value(), got.Value())
			assert.EqualValues(t, tt.want.CheckAvailability(), got.CheckAvailability())
			assert.EqualValues(t, tt.want.CheckLimit(), got.CheckLimit())
			assert.EqualValues(t, tt.want.IsOutdated(), got.IsOutdated())
		})
	}
}

func TestNewBatch(t *testing.T) {
	type args struct {
		lastID   uint
		minValue int
		duration time.Duration
	}
	tests := []struct {
		name string
		args args
		want Batch
	}{
		{name: "test batch 1", args: args{lastID: 0, minValue: 100, duration: 0}, want: Batch{ID: 1, Transactions: Transactions{}, MinValue: 100}},
		{name: "test batch 2", args: args{lastID: 1, minValue: 100, duration: 1}, want: Batch{ID: 2, Transactions: Transactions{}, MinValue: 100}},
		{name: "test batch 3", args: args{lastID: 2, minValue: 100, duration: 10}, want: Batch{ID: 3, Transactions: Transactions{}, MinValue: 100}},
		{name: "test batch 4", args: args{lastID: 3, minValue: 100, duration: 100}, want: Batch{ID: 4, Transactions: Transactions{}, MinValue: 100}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewBatch(tt.args.lastID, tt.args.minValue, tt.args.duration)
			assert.EqualValues(t, tt.want.Transactions, got.Transactions)
			assert.EqualValues(t, tt.want.ID, got.ID)
			assert.EqualValues(t, tt.want.Value(), got.Value())
			assert.EqualValues(t, tt.want.CheckAvailability(), got.CheckAvailability())
			assert.EqualValues(t, tt.want.CheckLimit(), got.CheckLimit())
			assert.EqualValues(t, tt.want.IsOutdated(), got.IsOutdated())
		})
	}
}
