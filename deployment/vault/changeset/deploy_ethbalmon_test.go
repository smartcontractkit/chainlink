package changeset

import (
	"fmt"
	"math"
	"testing"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
	"github.com/stretchr/testify/require"
)

func TestDeployEthBalMonValidation(t *testing.T) {
	t.Parallel()

	selector := chainselectors.TEST_90000001.Selector

	env, err := environment.New(t.Context(),
		environment.WithEVMSimulated(t, []uint64{selector}),
	)
	require.NoError(t, err)

	tests := []struct {
		name      string
		config    types.DeployEthBalMonInput
		wantError bool
		errorMsg  string
	}{
		{
			name: "empty chains",
			config: types.DeployEthBalMonInput{
				Chains: map[uint64]types.DeployEthBalMonChainConfig{},
			},
			wantError: true,
			errorMsg:  "chains must not be empty",
		},
		{
			name: "unknown chain selector",
			config: types.DeployEthBalMonInput{
				Chains: map[uint64]types.DeployEthBalMonChainConfig{
					math.MaxUint64: {
						SetKeeperRegistryAddress: "0x1234567890123456789012345678901234567890",
					},
				},
			},
			wantError: true,
			errorMsg:  fmt.Sprintf("unknown chain selector %d", uint64(math.MaxUint64)),
		},
		{
			name: "empty setKeeperRegistryAddress",
			config: types.DeployEthBalMonInput{
				Chains: map[uint64]types.DeployEthBalMonChainConfig{
					selector: {
						SetKeeperRegistryAddress: "",
					},
				},
			},
			wantError: true,
			errorMsg:  "setKeeperRegistryAddress must not be empty",
		},
		{
			name: "invalid setKeeperRegistryAddress",
			config: types.DeployEthBalMonInput{
				Chains: map[uint64]types.DeployEthBalMonChainConfig{
					selector: {
						SetKeeperRegistryAddress: "not-a-valid-address",
					},
				},
			},
			wantError: true,
			errorMsg:  fmt.Sprintf("chain %d: setKeeperRegistryAddress is not a valid hex address: not-a-valid-address", selector),
		},
		{
			name: "valid config",
			config: types.DeployEthBalMonInput{
				Chains: map[uint64]types.DeployEthBalMonChainConfig{
					selector: {
						SetKeeperRegistryAddress: "0x1234567890123456789012345678901234567890",
					},
				},
			},
			wantError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateDeployEthBalMonConfig(env.GetContext(), *env, test.config)

			if test.wantError {
				require.Error(t, err)
				if test.errorMsg != "" {
					require.Contains(t, err.Error(), test.errorMsg)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDeployEthBalMonChangeset(t *testing.T) {
	t.Parallel()
	// rt, err := runtime.New(t.Context(), runtime.WithEnvOpts())
	// require.NoError(t, err)

	t.Run("single chain", func(t *testing.T) {
		t.Parallel()

	})
}
