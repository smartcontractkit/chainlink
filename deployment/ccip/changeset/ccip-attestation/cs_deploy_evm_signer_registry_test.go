package ccip_attestation_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	ccip_attestation "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/ccip-attestation"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	signer_registry "github.com/smartcontractkit/chainlink/deployment/ccip/shared/bindings/signer_registry"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// Helper function to find signer registry address in address book
func findSignerRegistryAddress(e cldf.Environment, selector uint64) (common.Address, bool) {
	addresses, err := e.ExistingAddresses.AddressesForChain(selector)
	if err != nil {
		return common.Address{}, false
	}

	for addr, tv := range addresses {
		if tv.Type == shared.EVMSignerRegistry && tv.Version == deployment.Version1_0_0 {
			return common.HexToAddress(addr), true
		}
	}
	return common.Address{}, false
}

// Helper function to create test signers
func makeSigners(n int) []signer_registry.ISignerRegistrySigner {
	signers := make([]signer_registry.ISignerRegistrySigner, n)
	for i := 0; i < n; i++ {
		signers[i] = signer_registry.ISignerRegistrySigner{
			EvmAddress: utils.RandomAddress(),
			// Alternate between zero and non-zero NewEVMAddress
			NewEVMAddress: func() common.Address {
				if i%2 == 0 {
					return utils.ZeroAddress
				}
				return utils.RandomAddress()
			}(),
		}
	}
	return signers
}

func TestEVMSignerRegistry_Preconditions(t *testing.T) {
	t.Parallel()

	// Create a minimal environment for precondition tests
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 1,
	})

	tests := []struct {
		name        string
		config      ccip_attestation.SignerRegistryChangesetConfig
		expectedErr string
	}{
		{
			name: "Base case",
			config: ccip_attestation.SignerRegistryChangesetConfig{
				MaxSigners: ccip_attestation.MaxSigners,
				Signers: []signer_registry.ISignerRegistrySigner{
					{EvmAddress: utils.RandomAddress(), NewEVMAddress: utils.ZeroAddress},
				},
			},
			expectedErr: "",
		},
		{
			name: "MaxSigners mismatch",
			config: ccip_attestation.SignerRegistryChangesetConfig{
				MaxSigners: ccip_attestation.MaxSigners - 1,
				Signers:    []signer_registry.ISignerRegistrySigner{},
			},
			expectedErr: "max signers must be",
		},
		{
			name: "Too many signers",
			config: ccip_attestation.SignerRegistryChangesetConfig{
				MaxSigners: ccip_attestation.MaxSigners,
				Signers:    makeSigners(ccip_attestation.MaxSigners + 1),
			},
			expectedErr: "too many signers",
		},
		{
			name: "Zero evm address",
			config: ccip_attestation.SignerRegistryChangesetConfig{
				MaxSigners: ccip_attestation.MaxSigners,
				Signers: []signer_registry.ISignerRegistrySigner{
					{EvmAddress: utils.ZeroAddress, NewEVMAddress: utils.RandomAddress()},
				},
			},
			expectedErr: "has zero evm address",
		},
		{
			name: "Same evm and new address",
			config: func() ccip_attestation.SignerRegistryChangesetConfig {
				addr := utils.RandomAddress()
				return ccip_attestation.SignerRegistryChangesetConfig{
					MaxSigners: ccip_attestation.MaxSigners,
					Signers: []signer_registry.ISignerRegistrySigner{
						{EvmAddress: addr, NewEVMAddress: addr},
					},
				}
			}(),
			expectedErr: "has the same evm address and new evm address",
		},
		{
			name: "Duplicate EvmAddress",
			config: func() ccip_attestation.SignerRegistryChangesetConfig {
				addr := utils.RandomAddress()
				return ccip_attestation.SignerRegistryChangesetConfig{
					MaxSigners: ccip_attestation.MaxSigners,
					Signers: []signer_registry.ISignerRegistrySigner{
						{EvmAddress: addr, NewEVMAddress: utils.RandomAddress()},
						{EvmAddress: addr, NewEVMAddress: utils.ZeroAddress},
					},
				}
			}(),
			expectedErr: "duplicate signer evm address",
		},
		{
			name: "Duplicate non-zero NewEVMAddress",
			config: func() ccip_attestation.SignerRegistryChangesetConfig {
				newAddr := utils.RandomAddress()
				return ccip_attestation.SignerRegistryChangesetConfig{
					MaxSigners: ccip_attestation.MaxSigners,
					Signers: []signer_registry.ISignerRegistrySigner{
						{EvmAddress: utils.RandomAddress(), NewEVMAddress: newAddr},
						{EvmAddress: utils.RandomAddress(), NewEVMAddress: newAddr},
					},
				}
			}(),
			expectedErr: "duplicate signer new EVM address",
		},
		{
			name: "EvmAddress equals another's NewEVMAddress",
			config: func() ccip_attestation.SignerRegistryChangesetConfig {
				addrB := utils.RandomAddress()
				return ccip_attestation.SignerRegistryChangesetConfig{
					MaxSigners: ccip_attestation.MaxSigners,
					Signers: []signer_registry.ISignerRegistrySigner{
						{EvmAddress: utils.RandomAddress(), NewEVMAddress: addrB},
						{EvmAddress: addrB, NewEVMAddress: utils.RandomAddress()},
					},
				}
			}(),
			expectedErr: "duplicate",
		},
		{
			name: "Valid config with multiple zero new addresses",
			config: ccip_attestation.SignerRegistryChangesetConfig{
				MaxSigners: ccip_attestation.MaxSigners,
				Signers: []signer_registry.ISignerRegistrySigner{
					{EvmAddress: utils.RandomAddress(), NewEVMAddress: utils.ZeroAddress},
					{EvmAddress: utils.RandomAddress(), NewEVMAddress: utils.ZeroAddress},
				},
			},
			expectedErr: "",
		},
		{
			name: "Valid config with max signers",
			config: ccip_attestation.SignerRegistryChangesetConfig{
				MaxSigners: ccip_attestation.MaxSigners,
				Signers:    makeSigners(ccip_attestation.MaxSigners),
			},
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := commonchangeset.Apply(t, e,
				commonchangeset.Configure(ccip_attestation.EVMSignerRegistryDeploymentChangeset, tt.config))

			if tt.expectedErr != "" {
				require.ErrorContains(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestEVMSignerRegistry_DeploysOnlyOnBaseChains(t *testing.T) {
	t.Parallel()

	// Create environment with Base Mainnet and Base Sepolia chain IDs
	lggr := logger.TestLogger(t)
	evmChains := memory.NewMemoryChainsEVMWithChainIDs(t, []uint64{8453, 84532}, 1) // Base Mainnet: 8453, Base Sepolia: 84532
	chains := cldf_chain.NewBlockChainsFromSlice(evmChains)
	e := memory.NewMemoryEnvironmentFromChainsNodes(t.Context, lggr, chains, map[string]memory.Node{})

	// Create config with test signers
	signer1 := utils.RandomAddress()
	signer2 := utils.RandomAddress()
	config := ccip_attestation.SignerRegistryChangesetConfig{
		MaxSigners: ccip_attestation.MaxSigners,
		Signers: []signer_registry.ISignerRegistrySigner{
			{EvmAddress: signer1, NewEVMAddress: utils.ZeroAddress},
			{EvmAddress: signer2, NewEVMAddress: utils.RandomAddress()},
		},
	}

	// Apply changeset - should deploy to both Base chains
	e, err := commonchangeset.Apply(t, e,
		commonchangeset.Configure(ccip_attestation.EVMSignerRegistryDeploymentChangeset, config))
	require.NoError(t, err)

	// Verify deployment on Base Mainnet
	baseMainnetAddr, found := findSignerRegistryAddress(e, ccip_attestation.BaseMainnetSelector)
	require.True(t, found, "signer registry should be deployed on Base Mainnet")
	require.NotEqual(t, common.Address{}, baseMainnetAddr)

	// Verify deployment on Base Sepolia
	baseSepoliaAddr, found := findSignerRegistryAddress(e, ccip_attestation.BaseSepoliaSelector)
	require.True(t, found, "signer registry should be deployed on Base Sepolia")
	require.NotEqual(t, common.Address{}, baseSepoliaAddr)

	// Verify contract state on Base Mainnet
	baseMainnetChain := e.BlockChains.EVMChains()[ccip_attestation.BaseMainnetSelector]
	registry, err := signer_registry.NewSignerRegistry(baseMainnetAddr, baseMainnetChain.Client)
	require.NoError(t, err)

	maxSigners, err := registry.GetMaxSigners(nil)
	require.NoError(t, err)
	require.Equal(t, uint64(config.MaxSigners), maxSigners.Uint64())

	signerCount, err := registry.GetSignerCount(nil)
	require.NoError(t, err)
	require.Equal(t, uint64(len(config.Signers)), signerCount.Uint64())

	// Verify signers
	signers, err := registry.GetSigners(nil)
	require.NoError(t, err)
	require.Len(t, signers, len(config.Signers))
	require.Equal(t, signer1, signers[0].EvmAddress)
	require.Equal(t, signer2, signers[1].EvmAddress)
}

func TestEVMSignerRegistry_SkipsNonBaseChains(t *testing.T) {
	t.Parallel()

	// Create environment with non-Base chains
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 2, // These will be test chains, not Base chains
	})

	config := ccip_attestation.SignerRegistryChangesetConfig{
		MaxSigners: ccip_attestation.MaxSigners,
		Signers: []signer_registry.ISignerRegistrySigner{
			{EvmAddress: utils.RandomAddress(), NewEVMAddress: utils.ZeroAddress},
		},
	}

	// Apply changeset - should skip all non-Base chains
	e, err := commonchangeset.Apply(t, e,
		commonchangeset.Configure(ccip_attestation.EVMSignerRegistryDeploymentChangeset, config))
	require.NoError(t, err)

	// Verify no deployment on any chain
	for selector := range e.BlockChains.EVMChains() {
		_, found := findSignerRegistryAddress(e, selector)
		require.False(t, found, "signer registry should not be deployed on non-Base chain %d", selector)
	}
}
