package changeset

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/deployment"
	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commonTypes "github.com/smartcontractkit/chainlink/deployment/common/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

func TestAcceptOwnership(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		Nodes:  1,
		Chains: 1,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)

	chainSelector := env.AllChainSelectors()[0]

	// Deploy and configure pre-requisite contracts
	cache, _ := DeployCache(env.Chains[chainSelector], []string{})
	tx, _ := cache.Contract.TransferOwnership(env.Chains[chainSelector].DeployerKey, common.HexToAddress("0x123"))
	_, err := env.Chains[chainSelector].Confirm(tx)
	require.NoError(t, err)

	ab, err := commonChangesets.DeployMCMSWithTimelockV2(env, map[uint64]commonTypes.MCMSWithTimelockConfigV2{
		chainSelector: proposalutils.SingleGroupTimelockConfigV2(t),
	})
	require.NoError(t, err)

	addresses, err := ab.AddressBook.Addresses()
	require.NoError(t, err)

	env.ExistingAddresses = deployment.NewMemoryAddressBookFromMap(addresses)
	// End of pre-requisite contracts

	resp, err := AcceptOwnershipChangeset(env, types.AcceptOwnershipConfig{
		ChainSelector:   chainSelector,
		ContractAddress: cache.Contract.Address(),
		McmsConfig: &types.MCMSConfig{
			MinDelay: 1,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.MCMSTimelockProposals, 1)
}
