package v1_2_0

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/fee_quoter"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
)

// ApplyPriceRegistryUpdate is a helper function used in tests only.
func ApplyPriceRegistryUpdate(t *testing.T, user *bind.TransactOpts, addr common.Address, ec client.Client, gasPrices []cciptypes.GasPrice, tokenPrices []cciptypes.TokenPrice) common.Hash {
	require.True(t, len(gasPrices) <= 2)
	pr, err := fee_quoter.NewFeeQuoter(addr, ec)
	require.NoError(t, err)
	o, err := pr.Owner(nil)
	require.NoError(t, err)
	require.Equal(t, user.From, o)
	var tps []fee_quoter.InternalTokenPriceUpdate
	for _, tp := range tokenPrices {
		evmAddrs, err1 := ccipcalc.GenericAddrsToEvm(tp.Token)
		assert.NoError(t, err1)
		tps = append(tps, fee_quoter.InternalTokenPriceUpdate{
			SourceToken: evmAddrs[0],
			UsdPerToken: tp.Value,
		})
	}
	var gps []fee_quoter.InternalGasPriceUpdate
	for _, gp := range gasPrices {
		gps = append(gps, fee_quoter.InternalGasPriceUpdate{
			DestChainSelector: gp.DestChainSelector,
			UsdPerUnitGas:     gp.Value,
		})
	}
	tx, err := pr.UpdatePrices(user, fee_quoter.InternalPriceUpdates{
		TokenPriceUpdates: tps,
		GasPriceUpdates:   gps,
	})
	require.NoError(t, err)
	return tx.Hash()
}
