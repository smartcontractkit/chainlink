package v1_0_0

import (
	"encoding/binary"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/price_registry_1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
)

// ApplyPriceRegistryUpdate is a helper function used in tests only.
func ApplyPriceRegistryUpdate(t *testing.T, user *bind.TransactOpts, addr common.Address, ec client.Client, gasPrice []cciptypes.GasPrice, tokenPrices []cciptypes.TokenPrice) {
	require.True(t, len(gasPrice) <= 2)
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
	if len(gasPrice) >= 1 {
		dest = gasPrice[0].DestChainSelector
		gas = gasPrice[0].Value
	}
	_, err = pr.UpdatePrices(user, price_registry_1_0_0.InternalPriceUpdates{
		TokenPriceUpdates: tps,
		DestChainSelector: dest,
		UsdPerUnitGas:     gas,
	})
	require.NoError(t, err)

	for i := 1; i < len(gasPrice); i++ {
		dest = gasPrice[i].DestChainSelector
		gas = gasPrice[i].Value
		_, err = pr.UpdatePrices(user, price_registry_1_0_0.InternalPriceUpdates{
			TokenPriceUpdates: []price_registry_1_0_0.InternalTokenPriceUpdate{},
			DestChainSelector: dest,
			UsdPerUnitGas:     gas,
		})
		require.NoError(t, err)
	}
}

func CreateExecutionStateChangeEventLog(t *testing.T, seqNr uint64, blockNumber int64, messageID common.Hash) logpoller.Log {
	tAbi, err := evm_2_evm_offramp.EVM2EVMOffRampMetaData.GetAbi()
	require.NoError(t, err)
	eseEvent, ok := tAbi.Events["ExecutionStateChanged"]
	require.True(t, ok)

	logData, err := eseEvent.Inputs.NonIndexed().Pack(uint8(1), []byte("some return data"))
	require.NoError(t, err)
	seqNrBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(seqNrBytes, seqNr)
	seqNrTopic := common.BytesToHash(seqNrBytes)
	topic0 := evm_2_evm_offramp.EVM2EVMOffRampExecutionStateChanged{}.Topic()

	return logpoller.Log{
		Topics: [][]byte{
			topic0[:],
			seqNrTopic[:],
			messageID[:],
		},
		Data:        logData,
		LogIndex:    1,
		BlockHash:   utils.RandomBytes32(),
		BlockNumber: blockNumber,
		EventSig:    topic0,
		Address:     testutils.NewAddress(),
		TxHash:      utils.RandomBytes32(),
	}
}
