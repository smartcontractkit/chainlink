package state_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"slices"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	solanaMCMS "github.com/smartcontractkit/chainlink/deployment/common/changeset/solana/mcms"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestMCMSWithTimelockState_GenerateMCMSWithTimelockViewSolana(t *testing.T) {
	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/DX-404")

	t.Parallel()
	envConfig := memory.MemoryEnvironmentConfig{SolChains: 1}
	env := memory.NewMemoryEnvironment(t, logger.TestLogger(t), zapcore.InfoLevel, envConfig)
	chainSelector := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilySolana))[0]
	chain := env.BlockChains.SolanaChains()[chainSelector]
	defaultState := func() *state.MCMSWithTimelockStateSolana {
		addressBook := cldf.NewMemoryAddressBook()
		chainState, err := solanaMCMS.DeployMCMSWithTimelockProgramsSolana(env, chain, addressBook,
			commontypes.MCMSWithTimelockConfigV2{
				Proposer: mcmstypes.Config{
					Quorum:  1,
					Signers: []common.Address{common.HexToAddress("0x0000000000000000000000000000000000000001")},
				},
				Canceller: mcmstypes.Config{
					Quorum:  1,
					Signers: []common.Address{common.HexToAddress("0x0000000000000000000000000000000000000002")},
				},
				Bypasser: mcmstypes.Config{
					Quorum:  1,
					Signers: []common.Address{common.HexToAddress("0x0000000000000000000000000000000000000002")},
				},
				TimelockMinDelay: big.NewInt(1),
			},
		)
		require.NoError(t, err)
		env.ExistingAddresses = addressBook
		return chainState
	}

	changeset.SetPreloadedSolanaAddresses(t, env, chainSelector)

	tests := []struct {
		name    string
		state   *state.MCMSWithTimelockStateSolana
		want    func(*state.MCMSWithTimelockStateSolana) string
		wantErr string
	}{
		{
			name:  "success",
			state: defaultState(),
			want: func(state *state.MCMSWithTimelockStateSolana) string {
				return fmt.Sprintf(`{
					"proposer": {
						"programID": "%s",
						"seed":      "%s",
						"owner":     "11111111111111111111111111111111",
						"config":    {
							"quorum":       1,
							"signers":      ["0x0000000000000000000000000000000000000001"],
							"groupSigners": []
						}
					},
					"canceller": {
						"programID": "%s",
						"seed":      "%s",
						"owner":     "11111111111111111111111111111111",
						"config":    {
							"quorum":       1,
							"signers":      ["0x0000000000000000000000000000000000000002" ],
							"groupSigners": []
						}
					},
					"bypasser": {
						"programID": "%s",
						"seed":      "%s",
						"owner":     "11111111111111111111111111111111",
						"config":    {
							"quorum": 1,
							"signers": ["0x0000000000000000000000000000000000000002"],
							"groupSigners": []
						}
					},
					"timelock": {
						"programID":  "%s",
						"seed":       "%s",
						"owner":      "11111111111111111111111111111111",
						"proposers":  ["%s"],
						"executors":  ["%s"],
						"bypassers":  ["%s"],
						"cancellers": %s
					}
				}`, state.McmProgram, state.ProposerMcmSeed, state.McmProgram, state.CancellerMcmSeed,
					state.McmProgram, state.BypasserMcmSeed, state.TimelockProgram, state.TimelockSeed,
					signerPDA(state.McmProgram, state.ProposerMcmSeed), chain.DeployerKey.PublicKey(),
					signerPDA(state.McmProgram, state.BypasserMcmSeed),
					toJSON(t, slices.Sorted(slices.Values([]string{
						signerPDA(state.McmProgram, state.CancellerMcmSeed),
						signerPDA(state.McmProgram, state.ProposerMcmSeed),
						signerPDA(state.McmProgram, state.BypasserMcmSeed),
					}))),
				)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.state.GenerateView(env.GetContext(), chain)

			if tt.wantErr == "" {
				require.NoError(t, err)
				require.JSONEq(t, tt.want(tt.state), toJSON(t, &got))
			} else {
				require.ErrorContains(t, err, tt.wantErr)
			}
		})
	}
}

// ----- helpers -----

func toJSON[T any](t *testing.T, value T) string {
	t.Helper()

	bytes, err := json.Marshal(value)
	require.NoError(t, err)

	return string(bytes)
}

func signerPDA(programID solana.PublicKey, seed state.PDASeed) string {
	return state.GetMCMSignerPDA(programID, seed).String()
}
