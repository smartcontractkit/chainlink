package ccip_attestation_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"

	"github.com/smartcontractkit/chainlink/deployment"
	ccip_attestation "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/ccip-attestation"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	signer_registry "github.com/smartcontractkit/chainlink/deployment/ccip/shared/bindings/signer_registry"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

// deployTestSignerRegistryWithMax mirrors deployTestSignerRegistry but lets
// callers pick a custom max signers value (needed to exercise the capacity
// precondition without having to seed a huge number of signers).
func deployTestSignerRegistryWithMax(t *testing.T, env cldf.Environment, selector uint64, maxSigners int64, initialSigners []signer_registry.ISignerRegistrySigner) common.Address {
	chain := env.BlockChains.EVMChains()[selector]

	signerRegistry, err := cldf.DeployContract(env.Logger, chain, env.ExistingAddresses,
		func(chain cldf_evm.Chain) cldf.ContractDeploy[*signer_registry.SignerRegistry] {
			address, tx, signerRegistry, err := signer_registry.DeploySignerRegistry(
				chain.DeployerKey,
				chain.Client,
				big.NewInt(maxSigners),
				initialSigners,
			)
			return cldf.ContractDeploy[*signer_registry.SignerRegistry]{
				Address:  address,
				Contract: signerRegistry,
				Tx:       tx,
				Tv:       cldf.NewTypeAndVersion(shared.EVMSignerRegistry, deployment.Version1_0_0),
				Err:      err,
			}
		},
	)
	require.NoError(t, err)
	return signerRegistry.Address
}

func newBaseMainnetEnv(t *testing.T) cldf.Environment {
	e, err := environment.New(t.Context(),
		environment.WithLogger(logger.Test(t)),
		environment.WithEVMSimulated(t, []uint64{uint64(ccip_attestation.BaseMainnetSelector)}),
	)
	require.NoError(t, err)
	return *e
}

func TestEVMSignerRegistryAddSigners_Preconditions(t *testing.T) {
	t.Parallel()

	e := newBaseMainnetEnv(t)
	selector := uint64(ccip_attestation.BaseMainnetSelector)

	t.Run("Non-Base chain selector", func(t *testing.T) {
		nonBaseSelector := chain_selectors.ETHEREUM_MAINNET.Selector
		nonBaseEnv, envErr := environment.New(t.Context(),
			environment.WithLogger(logger.Test(t)),
			environment.WithEVMSimulated(t, []uint64{nonBaseSelector}),
		)
		require.NoError(t, envErr)

		config := ccip_attestation.AddSignersConfig{
			SignersByChain: map[uint64][]signer_registry.ISignerRegistrySigner{
				nonBaseSelector: {{EvmAddress: utils.RandomAddress()}},
			},
		}
		_, err := commonchangeset.Apply(t, *nonBaseEnv,
			commonchangeset.Configure(ccip_attestation.EVMSignerRegistryAddSignersChangeset, config))
		require.ErrorContains(t, err, "is not a Base chain")
	})

	tests := []struct {
		name        string
		config      ccip_attestation.AddSignersConfig
		expectedErr string
	}{
		{
			name: "Empty SignersByChain",
			config: ccip_attestation.AddSignersConfig{
				SignersByChain: map[uint64][]signer_registry.ISignerRegistrySigner{},
			},
			expectedErr: "no signers provided",
		},
		{
			name: "Empty signer list for chain",
			config: ccip_attestation.AddSignersConfig{
				SignersByChain: map[uint64][]signer_registry.ISignerRegistrySigner{
					selector: {},
				},
			},
			expectedErr: "no signers provided for chain selector",
		},
		{
			name: "Zero evm address",
			config: ccip_attestation.AddSignersConfig{
				SignersByChain: map[uint64][]signer_registry.ISignerRegistrySigner{
					selector: {{EvmAddress: utils.ZeroAddress}},
				},
			},
			expectedErr: "zero evm address",
		},
		{
			name: "Non-zero new evm address",
			config: func() ccip_attestation.AddSignersConfig {
				addr := utils.RandomAddress()
				return ccip_attestation.AddSignersConfig{
					SignersByChain: map[uint64][]signer_registry.ISignerRegistrySigner{
						selector: {{EvmAddress: addr, NewEVMAddress: addr}},
					},
				}
			}(),
			expectedErr: "non-zero new evm address",
		},
		{
			name: "Duplicate EvmAddress in batch",
			config: func() ccip_attestation.AddSignersConfig {
				addr := utils.RandomAddress()
				return ccip_attestation.AddSignersConfig{
					SignersByChain: map[uint64][]signer_registry.ISignerRegistrySigner{
						selector: {
							{EvmAddress: addr},
							{EvmAddress: addr},
						},
					},
				}
			}(),
			expectedErr: "duplicate signer evm address",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := commonchangeset.Apply(t, e,
				commonchangeset.Configure(ccip_attestation.EVMSignerRegistryAddSignersChangeset, tt.config))
			require.ErrorContains(t, err, tt.expectedErr)
		})
	}
}

func TestEVMSignerRegistryAddSigners_StateValidation(t *testing.T) {
	t.Parallel()

	e := newBaseMainnetEnv(t)
	selector := uint64(ccip_attestation.BaseMainnetSelector)

	// Deploy registry with two signers: signer1 (no pending) and signer2
	// (with a pending green key).
	signer1 := utils.RandomAddress()
	signer2 := utils.RandomAddress()
	pendingGreen := utils.RandomAddress()
	initialSigners := []signer_registry.ISignerRegistrySigner{
		{EvmAddress: signer1, NewEVMAddress: utils.ZeroAddress},
		{EvmAddress: signer2, NewEVMAddress: pendingGreen},
	}
	deployTestSignerRegistryWithMax(t, e, selector, ccip_attestation.MaxSigners, initialSigners)

	t.Run("No registry deployed on chain", func(t *testing.T) {
		eFresh := newBaseMainnetEnv(t)
		config := ccip_attestation.AddSignersConfig{
			SignersByChain: map[uint64][]signer_registry.ISignerRegistrySigner{
				selector: {{EvmAddress: utils.RandomAddress()}},
			},
		}
		_, err := commonchangeset.Apply(t, eFresh,
			commonchangeset.Configure(ccip_attestation.EVMSignerRegistryAddSignersChangeset, config))
		require.ErrorContains(t, err, "no signer registry found")
	})

	t.Run("New address collides with existing EvmAddress", func(t *testing.T) {
		config := ccip_attestation.AddSignersConfig{
			SignersByChain: map[uint64][]signer_registry.ISignerRegistrySigner{
				selector: {{EvmAddress: signer1}},
			},
		}
		_, err := commonchangeset.Apply(t, e,
			commonchangeset.Configure(ccip_attestation.EVMSignerRegistryAddSignersChangeset, config))
		require.ErrorContains(t, err, "is already registered")
	})

	t.Run("New address collides with existing pending NewEVMAddress", func(t *testing.T) {
		config := ccip_attestation.AddSignersConfig{
			SignersByChain: map[uint64][]signer_registry.ISignerRegistrySigner{
				selector: {{EvmAddress: pendingGreen}},
			},
		}
		_, err := commonchangeset.Apply(t, e,
			commonchangeset.Configure(ccip_attestation.EVMSignerRegistryAddSignersChangeset, config))
		require.ErrorContains(t, err, "is already registered")
	})

}

func TestEVMSignerRegistryAddSigners_CapacityCheck(t *testing.T) {
	t.Parallel()

	e := newBaseMainnetEnv(t)
	selector := uint64(ccip_attestation.BaseMainnetSelector)

	// Deploy a small registry (max=3) already holding 2 signers. Adding 2
	// more would overflow capacity (2 + 2 > 3).
	initialSigners := []signer_registry.ISignerRegistrySigner{
		{EvmAddress: utils.RandomAddress()},
		{EvmAddress: utils.RandomAddress()},
	}
	deployTestSignerRegistryWithMax(t, e, selector, 3, initialSigners)

	config := ccip_attestation.AddSignersConfig{
		SignersByChain: map[uint64][]signer_registry.ISignerRegistrySigner{
			selector: {
				{EvmAddress: utils.RandomAddress()},
				{EvmAddress: utils.RandomAddress()},
			},
		},
	}
	_, err := commonchangeset.Apply(t, e,
		commonchangeset.Configure(ccip_attestation.EVMSignerRegistryAddSignersChangeset, config))
	require.ErrorContains(t, err, "would exceed max signers")
}

func TestEVMSignerRegistryAddSigners_DirectExecution(t *testing.T) {
	t.Parallel()

	e := newBaseMainnetEnv(t)
	selector := uint64(ccip_attestation.BaseMainnetSelector)

	// Start at 5 signers; the production expansion goes from 5 to 10.
	initialSigners := make([]signer_registry.ISignerRegistrySigner, 5)
	for i := range initialSigners {
		initialSigners[i] = signer_registry.ISignerRegistrySigner{
			EvmAddress: utils.RandomAddress(),
		}
	}
	registryAddr := deployTestSignerRegistryWithMax(t, e, selector, ccip_attestation.MaxSigners, initialSigners)

	newSigners := make([]signer_registry.ISignerRegistrySigner, 5)
	for i := range newSigners {
		newSigners[i] = signer_registry.ISignerRegistrySigner{
			EvmAddress: utils.RandomAddress(),
		}
	}

	config := ccip_attestation.AddSignersConfig{
		SignersByChain: map[uint64][]signer_registry.ISignerRegistrySigner{
			selector: newSigners,
		},
	}

	_, outputs, err := commonchangeset.ApplyChangesets(t, e,
		[]commonchangeset.ConfiguredChangeSet{
			commonchangeset.Configure(ccip_attestation.EVMSignerRegistryAddSignersChangeset, config),
		})
	require.NoError(t, err)
	require.Len(t, outputs, 1)
	require.Empty(t, outputs[0].MCMSTimelockProposals, "direct execution should not produce MCMS proposals")

	chain := e.BlockChains.EVMChains()[selector]
	registry, err := signer_registry.NewSignerRegistry(registryAddr, chain.Client)
	require.NoError(t, err)

	count, err := registry.GetSignerCount(nil)
	require.NoError(t, err)
	require.Equal(t, uint64(10), count.Uint64(), "registry should have 5 initial + 5 new signers")

	registered, err := registry.GetSigners(nil)
	require.NoError(t, err)
	require.Len(t, registered, 10)

	// The union of initial + new addresses must all be present.
	present := make(map[common.Address]struct{}, len(registered))
	for _, s := range registered {
		present[s.EvmAddress] = struct{}{}
	}
	for _, s := range initialSigners {
		require.Contains(t, present, s.EvmAddress, "initial signer should remain")
	}
	for _, s := range newSigners {
		require.Contains(t, present, s.EvmAddress, "new signer should be registered")
	}
}

func TestEVMSignerRegistryAddSigners_MCMSProposal(t *testing.T) {
	t.Parallel()

	e := newBaseMainnetEnv(t)
	selector := uint64(ccip_attestation.BaseMainnetSelector)

	initialSigners := []signer_registry.ISignerRegistrySigner{
		{EvmAddress: utils.RandomAddress()},
	}
	registryAddr := deployTestSignerRegistryWithMax(t, e, selector, ccip_attestation.MaxSigners, initialSigners)

	e, err := commonchangeset.Apply(t, e,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2),
			map[uint64]commontypes.MCMSWithTimelockConfigV2{
				selector: proposalutils.SingleGroupTimelockConfigV2(t),
			},
		),
	)
	require.NoError(t, err)

	e, err = commonchangeset.Apply(t, e,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(commonchangeset.TransferToMCMSWithTimelockV2),
			commonchangeset.TransferToMCMSWithTimelockConfig{
				ContractsByChain: map[uint64][]common.Address{
					selector: {registryAddr},
				},
				MCMSConfig: proposalutils.TimelockConfig{MinDelay: 0},
			},
		),
	)
	require.NoError(t, err)

	chain := e.BlockChains.EVMChains()[selector]
	registry, err := signer_registry.NewSignerRegistry(registryAddr, chain.Client)
	require.NoError(t, err)

	state, err := stateview.LoadOnchainState(e)
	require.NoError(t, err)

	owner, err := registry.Owner(nil)
	require.NoError(t, err)
	require.Equal(t, state.Chains[selector].Timelock.Address(), owner)

	newSigners := []signer_registry.ISignerRegistrySigner{
		{EvmAddress: utils.RandomAddress()},
		{EvmAddress: utils.RandomAddress()},
	}
	config := ccip_attestation.AddSignersConfig{
		MCMS: &proposalutils.TimelockConfig{MinDelay: 0},
		SignersByChain: map[uint64][]signer_registry.ISignerRegistrySigner{
			selector: newSigners,
		},
	}

	output, err := commonchangeset.Configure(ccip_attestation.EVMSignerRegistryAddSignersChangeset, config).Apply(e)
	require.NoError(t, err)
	require.Len(t, output.MCMSTimelockProposals, 1)
	require.Equal(t, "add signers to EVM SignerRegistry", output.MCMSTimelockProposals[0].Description)
	require.Len(t, output.MCMSTimelockProposals[0].Operations, 1)
	require.Len(t, output.MCMSTimelockProposals[0].Operations[0].Transactions, 1)

	tx := output.MCMSTimelockProposals[0].Operations[0].Transactions[0]
	require.Equal(t, registryAddr.Hex(), tx.To)
	parsedABI, err := signer_registry.SignerRegistryMetaData.GetAbi()
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(tx.Data), 4)
	require.Equal(t, parsedABI.Methods["addSigners"].ID, tx.Data[:4])

	count, err := registry.GetSignerCount(nil)
	require.NoError(t, err)
	require.Equal(t, uint64(1), count.Uint64(), "proposal generation should not execute the signer update")

	_, outputs, err := commonchangeset.ApplyChangesets(t, e,
		[]commonchangeset.ConfiguredChangeSet{
			commonchangeset.Configure(ccip_attestation.EVMSignerRegistryAddSignersChangeset, config),
		})
	require.NoError(t, err)
	require.Len(t, outputs, 1)
	require.Len(t, outputs[0].MCMSTimelockProposals, 1)

	count, err = registry.GetSignerCount(nil)
	require.NoError(t, err)
	require.Equal(t, uint64(3), count.Uint64(), "executed MCMS proposal should add the new signers")

	registered, err := registry.GetSigners(nil)
	require.NoError(t, err)
	present := make(map[common.Address]struct{}, len(registered))
	for _, signer := range registered {
		present[signer.EvmAddress] = struct{}{}
	}
	for _, signer := range newSigners {
		require.Contains(t, present, signer.EvmAddress, "new signer should be registered after MCMS execution")
	}
}
