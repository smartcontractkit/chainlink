package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTxManager_updateLastConfirmedNonce(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		last      uint64
		submitted uint64
		want      uint64
	}{
		{"greater", 100, 101, 101},
		{"less", 100, 99, 100},
		{"equal", 100, 100, 100},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ma := &ManagedAccount{lastConfirmedNonce: test.last}
			ma.updateLastConfirmedNonce(test.submitted)

			assert.Equal(t, test.want, ma.lastConfirmedNonce)
		})
	}
}
