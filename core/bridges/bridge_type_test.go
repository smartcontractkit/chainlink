package bridges_test

import (
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/smartcontractkit/chainlink/core/bridges"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBridgeType_Authenticate(t *testing.T) {
	t.Parallel()

	bta, bt := cltest.NewBridgeType(t, cltest.BridgeOpts{})
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
			ok, err := bridges.AuthenticateBridgeType(bt, test.token)
			require.NoError(t, err)

			if test.wantError {
				assert.False(t, ok)
			} else {
				assert.True(t, ok)
			}
		})
	}
}

func BenchmarkParseBridgeName(b *testing.B) {
	const valid = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_`
	for _, l := range []int{1, 10, 20, 50, 100, 1000, 10000} {
		b.Run(strconv.Itoa(l), func(b *testing.B) {
			var sb strings.Builder
			for i := 0; i < l; i++ {
				sb.WriteByte(valid[rand.Intn(len(valid))])
			}
			name := sb.String()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := bridges.ParseBridgeName(name)
				if err != nil {
					b.Fatalf("failed to parse %q: %v\n", name, err)
				}
			}
		})
	}
}
