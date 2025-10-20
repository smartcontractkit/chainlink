package ccip_attestation_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	"github.com/smartcontractkit/chainlink/deployment"
	ccip_attestation "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/ccip-attestation"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	signer_registry "github.com/smartcontractkit/chainlink/deployment/ccip/shared/bindings/signer_registry"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
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

	selector := uint64(ccip_attestation.BaseMainnetSelector)
	e, err := environment.New(t.Context(),
		environment.WithEVMSimulated(t, []uint64{selector}),
	)
	require.NoError(t, err)

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
			_, err := commonchangeset.Apply(t, *e,
				commonchangeset.Configure(ccip_attestation.EVMSignerRegistrySetNewSignerAddressesChangeset, tt.config))
			require.ErrorContains(t, err, tt.expectedErr)
		})
	}
}

func TestEVMSignerRegistryConfiguration_StateValidation(t *testing.T) {
	t.Parallel()

	selector := uint64(ccip_attestation.BaseMainnetSelector)
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
	))
	require.NoError(t, err)

	// Deploy registry with known signers
	signer1 := utils.RandomAddress()
	signer2 := utils.RandomAddress()
	initialSigners := []signer_registry.ISignerRegistrySigner{
		{EvmAddress: signer1, NewEVMAddress: utils.ZeroAddress},
		{EvmAddress: signer2, NewEVMAddress: utils.ZeroAddress},
	}
	deployTestSignerRegistry(t, rt.Environment(), selector, initialSigners)

	// Test updating non-existent signer
	nonExistent := utils.RandomAddress()
	config := ccip_attestation.SetNewSignerAddressesConfig{
		UpdatesByChain: map[uint64]map[common.Address]common.Address{
			selector: {
				nonExistent: utils.RandomAddress(),
			},
		},
	}

	err = rt.Exec(
		runtime.ChangesetTask(ccip_attestation.EVMSignerRegistrySetNewSignerAddressesChangeset, config),
	)
	require.ErrorContains(t, err, "is not a registered signer")

	// Test new address conflicts with existing signer
	config = ccip_attestation.SetNewSignerAddressesConfig{
		UpdatesByChain: map[uint64]map[common.Address]common.Address{
			selector: {
				signer1: signer2, // signer2 already exists
			},
		},
	}

	// Test new address conflicts with existing signer
	err = rt.Exec(
		runtime.ChangesetTask(ccip_attestation.EVMSignerRegistrySetNewSignerAddressesChangeset, config),
	)
	require.ErrorContains(t, err, "is already a signer")
}

func TestEVMSignerRegistryConfiguration_DirectExecution(t *testing.T) {
	t.Parallel()

	selector := uint64(ccip_attestation.BaseMainnetSelector)
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
		environment.WithLogger(logger.Test(t)),
	))
	require.NoError(t, err)

	chain := rt.Environment().BlockChains.EVMChains()[selector]

	// Deploy registry with signers
	signer1 := utils.RandomAddress()
	signer2 := utils.RandomAddress()
	initialSigners := []signer_registry.ISignerRegistrySigner{
		{EvmAddress: signer1, NewEVMAddress: utils.ZeroAddress},
		{EvmAddress: signer2, NewEVMAddress: utils.ZeroAddress},
	}
	registryAddr := deployTestSignerRegistry(t, rt.Environment(), selector, initialSigners)

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
	err = rt.Exec(
		runtime.ChangesetTask(ccip_attestation.EVMSignerRegistrySetNewSignerAddressesChangeset, config),
	)
	require.NoError(t, err)
	// Verify no MCMS proposal (direct execution)
	require.Empty(t, rt.State().Proposals, 0)

	// Verify registry exists and was updated
	registry, err := signer_registry.NewSignerRegistry(registryAddr, chain.Client)
	require.NoError(t, err)

	// Check signer count is still correct
	count, err := registry.GetSignerCount(nil)
	require.NoError(t, err)
	require.Equal(t, uint64(2), count.Uint64())
}

func TestEVMSignerRegistryConfiguration_NoRegistries(t *testing.T) {
	t.Parallel()

	selector := uint64(ccip_attestation.BaseMainnetSelector)
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
		environment.WithLogger(logger.Test(t)),
	))
	require.NoError(t, err)

	// No registries deployed
	config := ccip_attestation.SetNewSignerAddressesConfig{
		UpdatesByChain: map[uint64]map[common.Address]common.Address{
			selector: {
				utils.RandomAddress(): utils.RandomAddress(),
			},
		},
	}

	// Should fail with error
	err = rt.Exec(
		runtime.ChangesetTask(ccip_attestation.EVMSignerRegistrySetNewSignerAddressesChangeset, config),
	)
	require.Error(t, err, "no signer registry found on chain selector %d", selector)
}
