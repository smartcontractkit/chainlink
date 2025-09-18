package ccip_attestation_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	ccip_attestation "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/ccip-attestation"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	signer_registry "github.com/smartcontractkit/chainlink/deployment/ccip/shared/bindings/signer_registry"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const (
	BaseMainnetID = 8453
)

// Helper to deploy signer registry directly for testing
func deployTestSignerRegistry(t *testing.T, env cldf.Environment, selector uint64, initialSigners []signer_registry.ISignerRegistrySigner) common.Address {
	chain := env.BlockChains.EVMChains()[selector]

	signerRegistry, err := cldf.DeployContract(env.Logger, chain, env.ExistingAddresses,
		func(chain cldf_evm.Chain) cldf.ContractDeploy[*signer_registry.SignerRegistry] {
			address, tx, signerRegistry, err := signer_registry.DeploySignerRegistry(
				chain.DeployerKey,
				chain.Client,
				big.NewInt(ccip_attestation.MaxSigners),
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

func TestEVMSignerRegistryConfiguration_Preconditions(t *testing.T) {
	t.Parallel()

	e := memory.NewMemoryEnvironment(t, logger.TestLogger(t), zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 1,
	})
	selector := uint64(ccip_attestation.BaseMainnetSelector)
	e.BlockChains = cldf_chain.NewBlockChainsFromSlice(
		memory.NewMemoryChainsEVMWithChainIDs(t, []uint64{BaseMainnetID}, 1),
	)

	tests := []struct {
		name        string
		config      ccip_attestation.SetNewSignerAddressesConfig
		expectedErr string
	}{
		{
			name: "Empty updates",
			config: ccip_attestation.SetNewSignerAddressesConfig{
				UpdatesByChain: map[uint64]map[common.Address]common.Address{},
			},

			expectedErr: "no signer updates provided",
		},
		{
			name: "Zero existing address",
			config: ccip_attestation.SetNewSignerAddressesConfig{
				UpdatesByChain: map[uint64]map[common.Address]common.Address{
					selector: {
						utils.ZeroAddress: utils.RandomAddress(),
					},
				},
			},
			expectedErr: "existing signer address cannot be zero address",
		},
		{
			name: "Zero new address",
			config: ccip_attestation.SetNewSignerAddressesConfig{
				UpdatesByChain: map[uint64]map[common.Address]common.Address{
					selector: {
						utils.RandomAddress(): utils.ZeroAddress,
					},
				},
			},
			expectedErr: "cannot be zero address",
		},
		{
			name: "Same existing and new",
			config: func() ccip_attestation.SetNewSignerAddressesConfig {
				addr := utils.RandomAddress()
				return ccip_attestation.SetNewSignerAddressesConfig{
					UpdatesByChain: map[uint64]map[common.Address]common.Address{
						selector: {
							addr: addr,
						},
					},
				}
			}(),
			expectedErr: "and new address are the same",
		},
		{
			name: "Duplicate new addresses",
			config: func() ccip_attestation.SetNewSignerAddressesConfig {
				newAddr := utils.RandomAddress()
				return ccip_attestation.SetNewSignerAddressesConfig{
					UpdatesByChain: map[uint64]map[common.Address]common.Address{
						selector: {
							utils.RandomAddress(): newAddr,
							utils.RandomAddress(): newAddr,
						},
					},
				}
			}(),
			expectedErr: "duplicate new address",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := commonchangeset.Apply(t, e,
				commonchangeset.Configure(ccip_attestation.EVMSignerRegistrySetNewSignerAddressesChangeset, tt.config))
			require.ErrorContains(t, err, tt.expectedErr)
		})
	}
}

func TestEVMSignerRegistryConfiguration_StateValidation(t *testing.T) {
	t.Parallel()

	e := memory.NewMemoryEnvironment(t, logger.TestLogger(t), zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 1,
	})

	selector := uint64(ccip_attestation.BaseMainnetSelector)
	e.BlockChains = cldf_chain.NewBlockChainsFromSlice(
		memory.NewMemoryChainsEVMWithChainIDs(t, []uint64{BaseMainnetID}, 1),
	)

	// Deploy registry with known signers
	signer1 := utils.RandomAddress()
	signer2 := utils.RandomAddress()
	initialSigners := []signer_registry.ISignerRegistrySigner{
		{EvmAddress: signer1, NewEVMAddress: utils.ZeroAddress},
		{EvmAddress: signer2, NewEVMAddress: utils.ZeroAddress},
	}
	deployTestSignerRegistry(t, e, selector, initialSigners)

	// Test updating non-existent signer
	nonExistent := utils.RandomAddress()
	config := ccip_attestation.SetNewSignerAddressesConfig{
		UpdatesByChain: map[uint64]map[common.Address]common.Address{
			selector: {
				nonExistent: utils.RandomAddress(),
			},
		},
	}

	_, err := commonchangeset.Apply(t, e,
		commonchangeset.Configure(ccip_attestation.EVMSignerRegistrySetNewSignerAddressesChangeset, config))
	require.ErrorContains(t, err, "is not a registered signer")

	// Test new address conflicts with existing signer
	config = ccip_attestation.SetNewSignerAddressesConfig{
		UpdatesByChain: map[uint64]map[common.Address]common.Address{
			selector: {
				signer1: signer2, // signer2 already exists
			},
		},
	}

	_, err = commonchangeset.Apply(t, e,
		commonchangeset.Configure(ccip_attestation.EVMSignerRegistrySetNewSignerAddressesChangeset, config))
	require.ErrorContains(t, err, "is already a signer")
}

func TestEVMSignerRegistryConfiguration_DirectExecution(t *testing.T) {
	t.Parallel()

	e := memory.NewMemoryEnvironment(t, logger.TestLogger(t), zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 1,
	})

	selector := uint64(ccip_attestation.BaseMainnetSelector)
	e.BlockChains = cldf_chain.NewBlockChainsFromSlice(
		memory.NewMemoryChainsEVMWithChainIDs(t, []uint64{BaseMainnetID}, 1),
	)

	// Deploy registry with signers
	signer1 := utils.RandomAddress()
	signer2 := utils.RandomAddress()
	initialSigners := []signer_registry.ISignerRegistrySigner{
		{EvmAddress: signer1, NewEVMAddress: utils.ZeroAddress},
		{EvmAddress: signer2, NewEVMAddress: utils.ZeroAddress},
	}
	registryAddr := deployTestSignerRegistry(t, e, selector, initialSigners)

	// Configure valid updates
	config := ccip_attestation.SetNewSignerAddressesConfig{
		UpdatesByChain: map[uint64]map[common.Address]common.Address{
			selector: {
				signer1: utils.RandomAddress(),
				signer2: utils.RandomAddress(),
			},
		},
	}

	// Execute changeset
	_, outputs, err := commonchangeset.ApplyChangesets(t, e,
		[]commonchangeset.ConfiguredChangeSet{
			commonchangeset.Configure(ccip_attestation.EVMSignerRegistrySetNewSignerAddressesChangeset, config),
		})
	require.NoError(t, err)
	require.Len(t, outputs, 1)

	// Verify no MCMS proposal (direct execution)
	require.Empty(t, outputs[0].MCMSTimelockProposals)

	// Verify registry exists and was updated
	chain := e.BlockChains.EVMChains()[selector]
	registry, err := signer_registry.NewSignerRegistry(registryAddr, chain.Client)
	require.NoError(t, err)

	// Check signer count is still correct
	count, err := registry.GetSignerCount(nil)
	require.NoError(t, err)
	require.Equal(t, uint64(2), count.Uint64())
}

func TestEVMSignerRegistryConfiguration_NoRegistries(t *testing.T) {
	t.Parallel()

	e := memory.NewMemoryEnvironment(t, logger.TestLogger(t), zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 1,
	})

	selector := uint64(ccip_attestation.BaseMainnetSelector)
	e.BlockChains = cldf_chain.NewBlockChainsFromSlice(
		memory.NewMemoryChainsEVMWithChainIDs(t, []uint64{BaseMainnetID}, 1),
	)

	// No registries deployed
	config := ccip_attestation.SetNewSignerAddressesConfig{
		UpdatesByChain: map[uint64]map[common.Address]common.Address{
			selector: {
				utils.RandomAddress(): utils.RandomAddress(),
			},
		},
	}

	// Should fail with error
	_, outputs, err := commonchangeset.ApplyChangesets(t, e,
		[]commonchangeset.ConfiguredChangeSet{
			commonchangeset.Configure(ccip_attestation.EVMSignerRegistrySetNewSignerAddressesChangeset, config),
		})
	require.Error(t, err, "no signer registry found on chain selector %d", selector)
	require.Empty(t, outputs)
}
