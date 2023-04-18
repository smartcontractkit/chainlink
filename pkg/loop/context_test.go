package loop

import (
	"context"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-relay/pkg/utils"
)

func TestContextValues(t *testing.T) {
	exp := ContextValues{
		JobID:      42,
		JobName:    "name",
		ContractID: utils.MustParseURL("http://example.com"),
		FeedID:     big.NewInt(1234567890987654321),
	}
	ctx := exp.ContextWithValues(context.Background())
	var got ContextValues
	got.SetValues(ctx)
	require.Equal(t, exp, got)
}
