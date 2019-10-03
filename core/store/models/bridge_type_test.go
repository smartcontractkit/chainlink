package models_test

import (
	"testing"

	"chainlink/core/internal/cltest"
	"chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBridgeType_Authenticate(t *testing.T) {
	t.Parallel()

	bta, bt := cltest.NewBridgeType(t)
	tests := []struct {
		name, token string
		wantError   bool
	}{
		{"correct", bta.IncomingToken, false},
		{"incorrect", "gibberish", true},
		{"empty incorrect", "", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ok, err := models.AuthenticateBridgeType(bt, test.token)
			require.NoError(t, err)

			if test.wantError {
				assert.False(t, ok)
			} else {
				assert.True(t, ok)
			}
		})
	}
}
