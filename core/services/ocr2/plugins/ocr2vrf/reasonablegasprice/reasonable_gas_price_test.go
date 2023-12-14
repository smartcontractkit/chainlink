package reasonablegasprice

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
)

func Test_ReasonableGasPrice(t *testing.T) {
	t.Parallel()

	t.Run("returns reasonable gas price", func(t *testing.T) {
		r := NewReasonableGasPriceProvider(nil, 1*time.Second, assets.GWei(100), true)
		g, err := r.ReasonableGasPrice()
		require.NoError(t, err)

		require.Equal(t, int64(0), g.Int64())
	})
}
