package api

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_errorResponse(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want errResponse
	}{
		{
			name: "test 1",
			err:  fmt.Errorf("test 1"),
			want: errResponse{Error: errorDetails{Message: "test 1"}},
		}, {
			name: "test 2",
			err:  fmt.Errorf("test 2"),
			want: errResponse{Error: errorDetails{Message: "test 2"}},
		}, {
			name: "test nil",
			err:  nil,
			want: errResponse{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := errorResponse(tt.err)
			assert.Equal(t, tt.want, got)
		})
	}
}
