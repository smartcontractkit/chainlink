package v1_0_0

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/price_registry_1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/cciptypes"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
)

// ApplyPriceRegistryUpdate is a helper function used in tests only.
func ApplyPriceRegistryUpdate(t *testing.T, user *bind.TransactOpts, addr common.Address, ec client.Client, gasPrice []cciptypes.GasPrice, tokenPrices []cciptypes.TokenPrice) {
	require.True(t, len(gasPrice) <= 1)
	pr, err := price_registry_1_0_0.NewPriceRegistry(addr, ec)
	require.NoError(t, err)
	var tps []price_registry_1_0_0.InternalTokenPriceUpdate
	for _, tp := range tokenPrices {
		evmAddrs, err1 := ccipcalc.GenericAddrsToEvm(tp.Token)
		assert.NoError(t, err1)
		tps = append(tps, price_registry_1_0_0.InternalTokenPriceUpdate{
			SourceToken: evmAddrs[0],
			UsdPerToken: tp.Value,
		})
	}
	dest := uint64(0)
	gas := big.NewInt(0)
	if len(gasPrice) == 1 {
		dest = gasPrice[0].DestChainSelector
		gas = gasPrice[0].Value
	}
	_, err = pr.UpdatePrices(user, price_registry_1_0_0.InternalPriceUpdates{
		TokenPriceUpdates: tps,
		DestChainSelector: dest,
		UsdPerUnitGas:     gas,
	})
	require.NoError(t, err)
}
