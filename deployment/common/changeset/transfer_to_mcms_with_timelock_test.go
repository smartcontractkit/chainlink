package changeset

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/stretchr/testify/require"

	owner_helpers "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	evminternal "github.com/smartcontractkit/chainlink/deployment/common/changeset/internal/evm"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/internal/seqs"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestTransferToMCMSWithTimelockV2(t *testing.T) {
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, 0, memory.MemoryEnvironmentConfig{
		Chains: 1,
		Nodes:  1,
	})
	chain1 := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]
	evmChains := e.BlockChains.EVMChains()
	e, err := Apply(t, e, Configure(
		cldf.CreateLegacyChangeSet(DeployLinkToken),
		[]uint64{chain1},
	), Configure(
		cldf.CreateLegacyChangeSet(DeployMCMSWithTimelockV2),
		map[uint64]types.MCMSWithTimelockConfigV2{
			chain1: proposalutils.SingleGroupTimelockConfigV2(t),
		},
	))
	require.NoError(t, err)
	addrs, err := e.ExistingAddresses.AddressesForChain(chain1)
	require.NoError(t, err)
	state, err := MaybeLoadMCMSWithTimelockChainState(evmChains[chain1], addrs)
	require.NoError(t, err)
	link, err := MaybeLoadLinkTokenChainState(evmChains[chain1], addrs)
	require.NoError(t, err)
	e, err = Apply(t, e,
		Configure(
			cldf.CreateLegacyChangeSet(TransferToMCMSWithTimelockV2),
			TransferToMCMSWithTimelockConfig{
				ContractsByChain: map[uint64][]common.Address{
					chain1: {link.LinkToken.Address()},
				},
				MCMSConfig: proposalutils.TimelockConfig{
					MinDelay: 0,
				},
			},
		),
	)
	require.NoError(t, err)
	// We expect now that the link token is owned by the MCMS timelock.
	link, err = MaybeLoadLinkTokenChainState(evmChains[chain1], addrs)
	require.NoError(t, err)
	o, err := link.LinkToken.Owner(nil)
	require.NoError(t, err)
	require.Equal(t, state.Timelock.Address(), o)

	// Try a rollback to the deployer.
	e, err = Apply(t, e,
		Configure(
			cldf.CreateLegacyChangeSet(TransferToDeployer),
			TransferToDeployerConfig{
				ContractAddress: link.LinkToken.Address(),
				ChainSel:        chain1,
			},
		),
	)
	require.NoError(t, err)

	o, err = link.LinkToken.Owner(nil)
	require.NoError(t, err)
	require.Equal(t, evmChains[chain1].DeployerKey.From, o)
}

func TestRenounceTimelockDeployerConfigValidate(t *testing.T) {
	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/DX-724")
	t.Parallel()
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, 0, memory.MemoryEnvironmentConfig{
		Chains: 1,
		Nodes:  1,
	})
	chain1 := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]
	e, err := Apply(t, e,
		Configure(
			cldf.CreateLegacyChangeSet(DeployMCMSWithTimelockV2),
			map[uint64]types.MCMSWithTimelockConfigV2{
				chain1: proposalutils.SingleGroupTimelockConfigV2(t),
			},
		),
	)
	require.NoError(t, err)

	envWithNoMCMS := memory.NewMemoryEnvironment(t, lggr, 0, memory.MemoryEnvironmentConfig{
		Chains: 1,
		Nodes:  1,
	})
	chain2 := envWithNoMCMS.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]

	for _, test := range []struct {
		name   string
		config RenounceTimelockDeployerConfig
		env    cldf.Environment
		err    string
	}{
		{
			name: "valid config",
			env:  e,
			config: RenounceTimelockDeployerConfig{
				ChainSel: chain1,
			},
		},
		{
			name: "invalid chain selector",
			env:  e,
			config: RenounceTimelockDeployerConfig{
				ChainSel: 0,
			},
			err: "invalid chain selector: chain selector must be set",
		},
		{
			name: "chain does not exists on env",
			env:  e,
			config: RenounceTimelockDeployerConfig{
				ChainSel: chain_selectors.ETHEREUM_TESTNET_SEPOLIA.Selector,
			},
			err: "chain selector: 16015286601757825753 not found in environment",
		},
		{
			name: "no MCMS deployed",
			env:  envWithNoMCMS,
			config: RenounceTimelockDeployerConfig{
				ChainSel: chain2,
			},
			// chain does not match any existing addresses
			err: "chain selector 909606746561742123: chain not found",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := test.config.Validate(test.env)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRenounceTimelockDeployer(t *testing.T) {
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, 0, memory.MemoryEnvironmentConfig{
		Chains: 1,
		Nodes:  1,
	})
	chain1 := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]
	e, err := Apply(t, e,
		Configure(
			cldf.CreateLegacyChangeSet(DeployMCMSWithTimelockV2),
			map[uint64]types.MCMSWithTimelockConfigV2{
				chain1: proposalutils.SingleGroupTimelockConfigV2(t),
			},
		),
	)
	require.NoError(t, err)
	addrs, err := e.ExistingAddresses.AddressesForChain(chain1)
	require.NoError(t, err)

	state, err := MaybeLoadMCMSWithTimelockChainState(e.BlockChains.EVMChains()[chain1], addrs)
	require.NoError(t, err)

	tl := state.Timelock
	require.NotNil(t, tl)

	adminRole, err := tl.ADMINROLE(nil)
	require.NoError(t, err)

	r, err := tl.GetRoleMemberCount(&bind.CallOpts{}, adminRole)
	require.NoError(t, err)
	require.Equal(t, int64(2), r.Int64())

	// Revoke Deployer
	e, err = Apply(t, e,
		Configure(
			cldf.CreateLegacyChangeSet(RenounceTimelockDeployer),
			RenounceTimelockDeployerConfig{
				ChainSel: chain1,
			},
		),
	)
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

func TestGrantRolesForTimelockWhenDeployerisAdmin(t *testing.T) {
	lggr := logger.TestLogger(t)
	tenv := memory.NewMemoryEnvironment(t, lggr, 0, memory.MemoryEnvironmentConfig{
		Chains: 1,
		Nodes:  1,
	})
	chain1 := tenv.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]
	ab := cldf.NewMemoryAddressBook()
	mcmsTimelockContracts, err := evminternal.DeployMCMSWithTimelockContractsEVM(tenv,
		tenv.BlockChains.EVMChains()[chain1],
		ab, proposalutils.SingleGroupTimelockConfigV2(t), nil)
	require.NoError(t, err)

	addresses, err := ab.AddressesForChain(chain1)
	require.NoError(t, err)
	require.Len(t, addresses, 5)

	mcmsTxs, err := evminternal.GrantRolesForTimelock(tenv, tenv.BlockChains.EVMChains()[chain1], mcmsTimelockContracts, false)
	require.NoError(t, err)
	require.Empty(t, mcmsTxs, "expected 0 transactions to be created for granting roles as Deployer is the admin still")
}

func TestGrantRolesForTimelockWhenDeployerisNotAdmin(t *testing.T) {
	lggr := logger.TestLogger(t)
	tenv := memory.NewMemoryEnvironment(t, lggr, 0, memory.MemoryEnvironmentConfig{
		Chains: 1,
		Nodes:  1,
	})
	chain1 := tenv.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]
	// Deploy MCMS with Timelock V2
	tenv, err := Apply(t, tenv,
		Configure(
			cldf.CreateLegacyChangeSet(DeployMCMSWithTimelockV2),
			map[uint64]types.MCMSWithTimelockConfigV2{
				chain1: proposalutils.SingleGroupTimelockConfigV2(t),
			},
		),
	)
	require.NoError(t, err)

	addrs, err := tenv.ExistingAddresses.AddressesForChain(chain1)
	require.NoError(t, err)

	state, err := MaybeLoadMCMSWithTimelockChainState(tenv.BlockChains.EVMChains()[chain1], addrs)
	require.NoError(t, err)

	mcmsTimelockContracts := &state.MCMSWithTimelockContracts

	// Revoke Deployer
	tenv, err = Apply(t, tenv,
		Configure(
			cldf.CreateLegacyChangeSet(RenounceTimelockDeployer),
			RenounceTimelockDeployerConfig{
				ChainSel: chain1,
			},
		),
	)
	require.NoError(t, err)
	ab := cldf.NewMemoryAddressBook()
	// Deploy new Proposer MCMS
	report, err := operations.ExecuteSequence(
		tenv.OperationsBundle,
		seqs.SeqEVMDeployMCMWithConfig,
		seqs.SeqDeployMCMWithConfigDeps{
			Chain:    tenv.BlockChains.EVMChains()[chain1],
			AddrBook: ab,
		},
		seqs.SeqDeployMCMWithConfigInput{
			ContractType:  types.ProposerManyChainMultisig,
			MCMConfig:     proposalutils.SingleGroupTimelockConfigV2(t).Proposer,
			ChainSelector: chain1,
		},
	)
	require.NoError(t, err)
	require.NotEmpty(t, report.Output.Address, "expected address to be set in output")

	proposer, err := owner_helpers.NewManyChainMultiSig(report.Output.Address, tenv.BlockChains.EVMChains()[chain1].Client)
	require.NoError(t, err)

	// Check if the new proposer grant roles produces MCMS Output when the deployer is not the admin
	(*mcmsTimelockContracts).ProposerMcm = proposer
	mcmsTxs, err := evminternal.GrantRolesForTimelock(tenv, tenv.BlockChains.EVMChains()[chain1], *mcmsTimelockContracts, false)
	require.NoError(t, err)
	require.Len(t, mcmsTxs, 2, "expected 2 mcms transaction to be created for granting Proposer & canceller roles when deployer is not the admin")
}
