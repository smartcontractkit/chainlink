package loop

import (
	"context"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-relay/pkg/config"
)

func TestContextValues(t *testing.T) {
	for _, tt := range []struct {
		name string
		vals ContextValues
		len  int
	}{
		{name: "full", vals: ContextValues{
			JobID:         42,
			JobName:       "name",
			ContractID:    config.MustParseURL("http://example.com"),
			FeedID:        big.NewInt(1234567890987654321),
			TransmitterID: "0xfake-test-id",
		}, len: 10},
		{name: "feedless", vals: ContextValues{
			JobID:      42,
			JobName:    "name",
			ContractID: config.MustParseURL("http://example.com"),
		}, len: 6},
		{name: "empty", vals: ContextValues{}, len: 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.vals.ContextWithValues(context.Background())
			var vals ContextValues
			vals.SetValues(ctx)
			require.Equal(t, tt.vals, vals)
			exp, got := tt.vals.Args(), vals.Args()
			require.Len(t, exp, tt.len)
			require.Equal(t, exp, got)
		})
	}
}
