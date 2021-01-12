package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name    string
		dbType  Type
		wantErr bool
		err     string
	}{
		{
			name:    "test memory db init",
			dbType:  InMemory,
			wantErr: false,
			err:     "",
		}, {
			name:    "test redis init",
			dbType:  Redis,
			wantErr: true,
			err:     `{"error":"dial tcp [::1]:6379: connect: connection refused: Cannot connect to Redis","meta":{}}`,
		}, {
			name:    "test wrong init",
			dbType:  "wrong",
			wantErr: true,
			err:     `{"error":"invalid database type","meta":{"type":"wrong"}}`,
		}, {
			name:    "test wrong type",
			dbType:  Type("wrong"),
			wantErr: true,
			err:     `{"error":"invalid database type","meta":{"type":"wrong"}}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Init(tt.dbType)
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, tt.err, err.Error())
				return
			}
			assert.Nil(t, err)
			assert.NotNil(t, got)
		})
	}
}
