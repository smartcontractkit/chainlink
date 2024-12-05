package changeset

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	owner_helpers "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/stretchr/testify/require"

	"math/big"

	"github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestTransferToMCMSWithTimelock(t *testing.T) {
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, 0, memory.MemoryEnvironmentConfig{
		Chains: 1,
		Nodes:  1,
	})
	chain1 := e.AllChainSelectors()[0]
	e, err := ApplyChangesets(t, e, nil, []ChangesetApplication{
		{
			Changeset: WrapChangeSet(DeployLinkToken),
			Config:    []uint64{chain1},
		},
		{
			Changeset: WrapChangeSet(DeployMCMSWithTimelock),
			Config: map[uint64]types.MCMSWithTimelockConfig{
				chain1: {
					Canceller:         SingleGroupMCMS(t),
					Bypasser:          SingleGroupMCMS(t),
					Proposer:          SingleGroupMCMS(t),
					TimelockExecutors: e.AllDeployerKeys(),
					TimelockMinDelay:  big.NewInt(0),
				},
			},
		},
	})
	require.NoError(t, err)
	addrs, err := e.ExistingAddresses.AddressesForChain(chain1)
	require.NoError(t, err)
	state, err := LoadMCMSWithTimelockState(e.Chains[chain1], addrs)
	require.NoError(t, err)
	link, err := LoadLinkTokenState(e.Chains[chain1], addrs)
	require.NoError(t, err)
	e, err = ApplyChangesets(t, e, map[uint64]*owner_helpers.RBACTimelock{
		chain1: state.Timelock,
	}, []ChangesetApplication{
		{
			Changeset: WrapChangeSet(TransferToMCMSWithTimelock),
			Config: TransferToMCMSWithTimelockConfig{
				ContractsByChain: map[uint64][]common.Address{
					chain1: {link.LinkToken.Address()},
				},
				MinDelay: 0,
			},
		},
	})
	require.NoError(t, err)
	// We expect now that the link token is owned by the MCMS timelock.
	link, err = LoadLinkTokenState(e.Chains[chain1], addrs)
	require.NoError(t, err)
	o, err := link.LinkToken.Owner(nil)
	require.NoError(t, err)
	require.Equal(t, state.Timelock.Address(), o)
}
