package changeset

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
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
				chain1: proposalutils.SingleGroupTimelockConfig(t),
			},
		},
	})
	require.NoError(t, err)
	addrs, err := e.ExistingAddresses.AddressesForChain(chain1)
	require.NoError(t, err)
	state, err := MaybeLoadMCMSWithTimelockState(e.Chains[chain1], addrs)
	require.NoError(t, err)
	link, err := MaybeLoadLinkTokenState(e.Chains[chain1], addrs)
	require.NoError(t, err)
	e, err = ApplyChangesets(t, e, map[uint64]*proposalutils.TimelockExecutionContracts{
		chain1: {
			Timelock:  state.Timelock,
			CallProxy: state.CallProxy,
		},
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
	link, err = MaybeLoadLinkTokenState(e.Chains[chain1], addrs)
	require.NoError(t, err)
	o, err := link.LinkToken.Owner(nil)
	require.NoError(t, err)
	require.Equal(t, state.Timelock.Address(), o)
}
