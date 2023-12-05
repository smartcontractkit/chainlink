package upkeepstate

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTokenBuckets(t *testing.T) {
	tests := []struct {
		name         string
		rate         uint32
		rateInterval time.Duration
		calls        []string
		accepts      []string
	}{
		{
			name:         "accepts up to the rate limit",
			rate:         2,
			rateInterval: time.Second,
			calls:        []string{"a", "a", "a"},
			accepts:      []string{"a", "a"},
		},
		{
			name:         "manage multiple items",
			rate:         2,
			rateInterval: time.Second,
			calls:        []string{"a", "a", "a", "b", "c", "c", "a", "b", "c"},
			accepts:      []string{"a", "a", "b", "c", "c", "b"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tb := newTokenBuckets(tc.rate, tc.rateInterval)
			accepts := make([]string, 0)
			for _, call := range tc.calls {
				if tb.Accept(call, 1) {
					accepts = append(accepts, call)
				}
			}
			require.Equal(t, len(tc.accepts), len(accepts))
			for i, accept := range tc.accepts {
				require.Equal(t, accept, accepts[i])
			}
		})
	}
}

func TestTokenBuckets_Clean(t *testing.T) {
	tb := newTokenBuckets(3, time.Second)

	require.True(t, tb.Accept("a", 3))
	require.False(t, tb.Accept("a", 1))

	require.True(t, tb.Accept("b", 2))
	require.True(t, tb.Accept("b", 1))
	require.False(t, tb.Accept("b", 1))

	tb.clean()

	require.True(t, tb.Accept("a", 1))
	require.True(t, tb.Accept("b", 1))
}
