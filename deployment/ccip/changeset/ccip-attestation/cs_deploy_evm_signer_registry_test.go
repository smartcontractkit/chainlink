package ccip_attestation_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	"github.com/smartcontractkit/chainlink/deployment"
	ccip_attestation "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/ccip-attestation"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	signer_registry "github.com/smartcontractkit/chainlink/deployment/ccip/shared/bindings/signer_registry"
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
	for i := range n {
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
	e, terr := environment.New(t.Context(),
		environment.WithLogger(logger.Test(t)),
	)
	require.NoError(t, terr)

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
			t.Parallel()

			terr = ccip_attestation.EVMSignerRegistryDeploymentChangeset.VerifyPreconditions(*e, tt.config)

			if tt.expectedErr != "" {
				require.ErrorContains(t, terr, tt.expectedErr)
			} else {
				require.NoError(t, terr)
			}
		})
	}
}

func TestEVMSignerRegistry_DeploysOnlyOnBaseChains(t *testing.T) {
	t.Parallel()

	// Create environment with Base Mainnet and Base Sepolia chain IDs
	baseMainnetSelector := uint64(ccip_attestation.BaseMainnetSelector)
	baseSepoliaSelector := uint64(ccip_attestation.BaseSepoliaSelector)
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{baseMainnetSelector, baseSepoliaSelector}),
		environment.WithLogger(logger.Test(t)),
	))
	require.NoError(t, err)

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
	err = rt.Exec(
		runtime.ChangesetTask(ccip_attestation.EVMSignerRegistryDeploymentChangeset, config),
	)
	require.NoError(t, err)

	// Verify deployment on Base Mainnet
	baseMainnetAddr, found := findSignerRegistryAddress(rt.Environment(), baseMainnetSelector)
	require.True(t, found, "signer registry should be deployed on Base Mainnet")
	require.NotEqual(t, common.Address{}, baseMainnetAddr)

	// Verify deployment on Base Sepolia
	baseSepoliaAddr, found := findSignerRegistryAddress(rt.Environment(), baseSepoliaSelector)
	require.True(t, found, "signer registry should be deployed on Base Sepolia")
	require.NotEqual(t, common.Address{}, baseSepoliaAddr)

	// Verify contract state on Base Mainnet
	baseMainnetChain := rt.Environment().BlockChains.EVMChains()[baseMainnetSelector]
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
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulatedN(t, 2),
		environment.WithLogger(logger.Test(t)),
	))
	require.NoError(t, err)

	config := ccip_attestation.SignerRegistryChangesetConfig{
		MaxSigners: ccip_attestation.MaxSigners,
		Signers: []signer_registry.ISignerRegistrySigner{
			{EvmAddress: utils.RandomAddress(), NewEVMAddress: utils.ZeroAddress},
		},
	}

	// Apply changeset - should skip all non-Base chains
	err = rt.Exec(
		runtime.ChangesetTask(ccip_attestation.EVMSignerRegistryDeploymentChangeset, config),
	)
	require.NoError(t, err)

	// Verify no deployment on any chain
	for selector := range rt.Environment().BlockChains.EVMChains() {
		_, found := findSignerRegistryAddress(rt.Environment(), selector)
		require.False(t, found, "signer registry should not be deployed on non-Base chain %d", selector)
	}
}
