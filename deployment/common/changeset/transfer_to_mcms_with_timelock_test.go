package changeset

import (
	"math/big"
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

func TestRenounceTimelockDeployer(t *testing.T) {
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

	adminRole, err := tl.ADMINROLE(nil)
	require.NoError(t, err)

	r, err := tl.GetRoleMemberCount(&bind.CallOpts{}, adminRole)
	require.NoError(t, err)
	require.Equal(t, int64(2), r.Int64())

	// Revoke Deployer
	e, err = ApplyChangesets(t, e, nil, []ChangesetApplication{
		{
			Changeset: WrapChangeSet(RenounceTimelockDeployer),
			Config: RenounceTimelockDeployerConfig{
				ChainSel: chain1,
			},
		},
	})
	require.NoError(t, err)

	// Check that the deployer is no longer an admin
	r, err = tl.GetRoleMemberCount(&bind.CallOpts{}, adminRole)
	require.NoError(t, err)
	require.Equal(t, int64(1), r.Int64())

	// Retrieve the admin address
	admin, err := tl.GetRoleMember(&bind.CallOpts{}, adminRole, big.NewInt(0))
	require.NoError(t, err)

	// Check that the admin is the timelock
	require.Equal(t, tl.Address(), admin)
}
