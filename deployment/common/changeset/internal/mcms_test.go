package internal_test

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/internal"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestDeployMCMSWithConfig(t *testing.T) {
	lggr := logger.TestLogger(t)
	chains := memory.NewMemoryChainsWithChainIDs(t, []uint64{
		chainsel.TEST_90000001.EvmChainID,
	})
	ab := deployment.NewMemoryAddressBook()
	_, err := internal.DeployMCMSWithConfig(types.ProposerManyChainMultisig,
		lggr, chains[chainsel.TEST_90000001.Selector], ab, changeset.SingleGroupMCMS(t))
	require.NoError(t, err)
}

func TestDeployMCMSWithTimelockContracts(t *testing.T) {
	lggr := logger.TestLogger(t)
	chains := memory.NewMemoryChainsWithChainIDs(t, []uint64{
		chainsel.TEST_90000001.EvmChainID,
	})
	ab := deployment.NewMemoryAddressBook()
	_, err := internal.DeployMCMSWithTimelockContracts(lggr,
		chains[chainsel.TEST_90000001.Selector],
		ab, types.MCMSWithTimelockConfig{
			Canceller: changeset.SingleGroupMCMS(t),
			Bypasser:  changeset.SingleGroupMCMS(t),
			Proposer:  changeset.SingleGroupMCMS(t),
			TimelockExecutors: []common.Address{
				chains[chainsel.TEST_90000001.Selector].DeployerKey.From,
			},
			TimelockMinDelay: big.NewInt(0),
		})
	require.NoError(t, err)
	addresses, err := ab.AddressesForChain(chainsel.TEST_90000001.Selector)
	require.NoError(t, err)
	require.Len(t, addresses, 4)
	mcmsState, err := changeset.LoadMCMSWithTimelockState(chains[chainsel.TEST_90000001.Selector], addresses)
	require.NoError(t, err)
	v, err := mcmsState.GenerateMCMSWithTimelockView()
	b, err := json.MarshalIndent(v, "", "  ")
	require.NoError(t, err)
	t.Log(string(b))
}
