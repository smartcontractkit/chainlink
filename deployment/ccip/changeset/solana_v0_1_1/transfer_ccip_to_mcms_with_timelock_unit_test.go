package solana_test

import (
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"

	ccipChangesetSolana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana_v0_1_1"
	solanastateview "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/solana"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

func TestValidateContracts(t *testing.T) {
	validPubkey := solana.NewWallet().PublicKey()

	zeroPubkey := solana.PublicKey{} // Zero public key

	makeState := func(router, feeQuoter solana.PublicKey) solanastateview.CCIPChainState {
		return solanastateview.CCIPChainState{
			Router:    router,
			FeeQuoter: feeQuoter,
		}
	}

	tests := []struct {
		name          string
		state         solanastateview.CCIPChainState
		contracts     ccipChangesetSolana.CCIPContractsToTransfer
		chainSelector uint64
		expectedError string
	}{
		{
			name:          "All required contracts present",
			state:         makeState(validPubkey, validPubkey),
			contracts:     ccipChangesetSolana.CCIPContractsToTransfer{Router: true},
			chainSelector: 12345,
		},
		{
			name:          "Missing Router contract",
			state:         makeState(zeroPubkey, validPubkey),
			contracts:     ccipChangesetSolana.CCIPContractsToTransfer{Router: true},
			chainSelector: 12345,
			expectedError: "missing required contract Router on chain 12345",
		},
		{
			name:          "Missing FeeQuoter contract",
			state:         makeState(validPubkey, zeroPubkey),
			contracts:     ccipChangesetSolana.CCIPContractsToTransfer{Router: true, FeeQuoter: true},
			chainSelector: 12345,
			expectedError: "missing required contract FeeQuoter on chain 12345",
		},
		{
			name:          "invalid pub key",
			state:         makeState(validPubkey, zeroPubkey),
			contracts:     ccipChangesetSolana.CCIPContractsToTransfer{Router: true, FeeQuoter: true},
			chainSelector: 12345,
			expectedError: "missing required contract FeeQuoter on chain 12345",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ccipChangesetSolana.ValidateContracts(tt.state, tt.chainSelector, tt.contracts)

			if tt.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Equal(t, tt.expectedError, err.Error())
			}
		})
	}
}

func TestValidate(t *testing.T) {
	selector := chainselectors.TEST_22222222222222222222222222222222222222222222.Selector

	tests := []struct {
		name             string
		env              func(t *testing.T) cldf.Environment
		contractsByChain map[uint64]ccipChangesetSolana.CCIPContractsToTransfer
		expectedError    string
	}{
		{
			name: "No chains found in environment",
			env: func(t *testing.T) cldf.Environment {
				t.Helper()

				e, err := environment.New(t.Context())
				require.NoError(t, err)

				return *e
			},
			expectedError: "no chains found",
		},
		{
			name: "Chain selector not found in environment",
			env: func(t *testing.T) cldf.Environment {
				t.Helper()

				e, err := environment.New(t.Context())
				require.NoError(t, err)

				e.BlockChains = cldf_chain.NewBlockChainsFromSlice([]cldf_chain.BlockChain{
					cldf_solana.Chain{
						Selector: selector,
					},
				})

				return *e
			},
			contractsByChain: map[uint64]ccipChangesetSolana.CCIPContractsToTransfer{
				99999: {Router: true, FeeQuoter: true},
			},
			expectedError: "chain 99999 not found in environment",
		},
		{
			name: "Invalid chain family",
			env: func(t *testing.T) cldf.Environment {
				t.Helper()

				e, err := environment.New(t.Context())
				require.NoError(t, err)

				e.BlockChains = cldf_chain.NewBlockChainsFromSlice([]cldf_chain.BlockChain{
					cldf_solana.Chain{
						Selector: selector,
					},
				})

				return *e
			},
			contractsByChain: map[uint64]ccipChangesetSolana.CCIPContractsToTransfer{
				selector: {Router: true, FeeQuoter: true},
			},
			expectedError: "failed to load addresses for chain 12463857294658392847: chain selector 12463857294658392847: chain not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := ccipChangesetSolana.TransferCCIPToMCMSWithTimelockSolanaConfig{
				ContractsByChain: tt.contractsByChain,
				MCMSCfg: proposalutils.TimelockConfig{
					MinDelay: 0 * time.Second,
				},
			}

			err := cfg.Validate(tt.env(t))

			if tt.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedError)
			}
		})
	}
}
