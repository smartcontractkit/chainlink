package changeset

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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
	state, err := MaybeLoadMCMSWithTimelockChainState(e.Chains[chain1], addrs)
	require.NoError(t, err)
	link, err := MaybeLoadLinkTokenChainState(e.Chains[chain1], addrs)
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
	link, err = MaybeLoadLinkTokenChainState(e.Chains[chain1], addrs)
	require.NoError(t, err)
	o, err := link.LinkToken.Owner(nil)
	require.NoError(t, err)
	require.Equal(t, state.Timelock.Address(), o)

	// Try a rollback to the deployer.
	e, err = ApplyChangesets(t, e, nil, []ChangesetApplication{
		{
			Changeset: WrapChangeSet(TransferToDeployer),
			Config: TransferToDeployerConfig{
				ContractAddress: link.LinkToken.Address(),
				ChainSel:        chain1,
			},
		},
	})
	require.NoError(t, err)

	o, err = link.LinkToken.Owner(nil)
	require.NoError(t, err)
	require.Equal(t, e.Chains[chain1].DeployerKey.From, o)
}

func TestRevokeTimelockDeployer(t *testing.T) {
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, 0, memory.MemoryEnvironmentConfig{
		Chains: 1,
		Nodes:  1,
	})
	chain1 := e.AllChainSelectors()[0]
	e, err := ApplyChangesets(t, e, nil, []ChangesetApplication{
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

	state, err := MaybeLoadMCMSWithTimelockChainState(e.Chains[chain1], addrs)
	require.NoError(t, err)

	tl := state.Timelock
	require.NotNil(t, tl)

	// WHAT AM I DOING HERE? HOW DO I TEST?
	adminRole, err := tl.ADMINROLE(nil)
	require.NoError(t, err)

	r, err := tl.GetRoleAdmin(&bind.CallOpts{}, adminRole)
	require.NoError(t, err)
	h := hex.EncodeToString(r[:])
	t.Logf("admin role: %s", h)
	a := common.BytesToAddress(r[:])

	require.Equal(t, a, e.Chains[chain1].DeployerKey.From)

}
