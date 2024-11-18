package internal

import (
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment"
	common2 "github.com/smartcontractkit/chainlink/deployment/common"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestDeployMCMSWithConfig(t *testing.T) {
	lggr := logger.TestLogger(t)
	chains := memory.NewMemoryChainsWithChainIDs(t, []uint64{
		chainsel.TEST_90000001.EvmChainID,
	})
	ab := deployment.NewMemoryAddressBook()
	_, err := DeployMCMSWithConfig(ProposerManyChainMultisig,
		lggr, chains[chainsel.TEST_90000001.Selector], ab, common2.SingleGroupMCMS(t))
	require.NoError(t, err)
}

func TestDeployMCMSWithTimelockContracts(t *testing.T) {
	lggr := logger.TestLogger(t)
	chains := memory.NewMemoryChainsWithChainIDs(t, []uint64{
		chainsel.TEST_90000001.EvmChainID,
	})
	ab := deployment.NewMemoryAddressBook()
	_, err := DeployMCMSWithTimelockContracts(lggr,
		chains[chainsel.TEST_90000001.Selector],
		ab, MCMSConfig{
			Canceller: common2.SingleGroupMCMS(t),
			Bypasser:  common2.SingleGroupMCMS(t),
			Proposer:  common2.SingleGroupMCMS(t),
			Executors: []common.Address{
				chains[chainsel.TEST_90000001.Selector].DeployerKey.From,
			},
		})
	require.NoError(t, err)
	addresses, err := ab.AddressesForChain(chainsel.TEST_90000001.Selector)
	require.NoError(t, err)
	require.Len(t, addresses, 4)
	mcmsState, err := common2.LoadMCMSWithTimelockState(chains[chainsel.TEST_90000001.Selector], addresses)
	require.NoError(t, err)
	v, err := mcmsState.GenerateMCMSWithTimelockView()
	b, err := json.MarshalIndent(v, "", "  ")
	require.NoError(t, err)
	t.Log(string(b))
}
