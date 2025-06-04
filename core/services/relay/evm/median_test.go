package evm

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
	"github.com/smartcontractkit/chainlink/v2/common/chains/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

func TestNewMedianProvider(t *testing.T) {
	lggr := logger.Test(t)

	chain := mocks.NewChain(t)
	chainID := testutils.NewRandomEVMChainID()
	chain.On("ID").Return(chainID)
	contractID := testutils.NewAddress()
	relayer := Relayer{lggr: logger.Sugared(lggr), chain: chain}

	pargs := commontypes.PluginArgs{}

	t.Run("wrong chainID", func(t *testing.T) {
		relayConfigBadChainID := evmtypes.RelayConfig{}
		rc, err2 := json.Marshal(&relayConfigBadChainID)
		rargs2 := commontypes.RelayArgs{ContractID: contractID.String(), RelayConfig: rc}
		require.NoError(t, err2)
		_, err2 = relayer.NewMedianProvider(testutils.Context(t), rargs2, pargs)
		assert.ErrorContains(t, err2, "chain id in spec does not match")
	})

	t.Run("invalid contractID", func(t *testing.T) {
		relayConfig := evmtypes.RelayConfig{ChainID: big.New(chainID)}
		rc, err2 := json.Marshal(&relayConfig)
		require.NoError(t, err2)
		rargsBadContractID := commontypes.RelayArgs{ContractID: "NotAContractID", RelayConfig: rc}
		_, err2 = relayer.NewMedianProvider(testutils.Context(t), rargsBadContractID, pargs)
		assert.ErrorContains(t, err2, "invalid contractID")
	})
}
