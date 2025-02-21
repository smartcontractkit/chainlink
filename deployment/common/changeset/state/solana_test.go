package state_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"slices"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	mcmstypes "github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/deployment"
	solanainternal "github.com/smartcontractkit/chainlink/deployment/common/changeset/internal/solana"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestMCMSWithTimelockState_GenerateMCMSWithTimelockViewSolana(t *testing.T) {
	envConfig := memory.MemoryEnvironmentConfig{SolChains: 1}
	env := memory.NewMemoryEnvironment(t, logger.TestLogger(t), zapcore.InfoLevel, envConfig)
	chainSelector := env.AllChainSelectorsSolana()[0]
	chain := env.SolChains[chainSelector]
	defaultState := func() *state.MCMSWithTimelockStateSolana {
		addressBook := deployment.NewMemoryAddressBook()
		chainState, err := solanainternal.DeployMCMSWithTimelockProgramsSolana(env, chain, addressBook,
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

	setPreloadedSolanaAddresses(t, env, chainSelector)

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

func setPreloadedSolanaAddresses(t *testing.T, env deployment.Environment, selector uint64) {
	typeAndVersion := deployment.NewTypeAndVersion(commontypes.ManyChainMultisigProgram, deployment.Version1_0_0)
	err := env.ExistingAddresses.Save(selector, memory.SolanaProgramIDs["mcm"], typeAndVersion)
	require.NoError(t, err)

	typeAndVersion = deployment.NewTypeAndVersion(commontypes.AccessControllerProgram, deployment.Version1_0_0)
	err = env.ExistingAddresses.Save(selector, memory.SolanaProgramIDs["access_controller"], typeAndVersion)
	require.NoError(t, err)

	typeAndVersion = deployment.NewTypeAndVersion(commontypes.RBACTimelockProgram, deployment.Version1_0_0)
	err = env.ExistingAddresses.Save(selector, memory.SolanaProgramIDs["timelock"], typeAndVersion)
	require.NoError(t, err)
}

func signerPDA(programID solana.PublicKey, seed state.PDASeed) string {
	return state.GetMCMSignerPDA(programID, seed).String()
}
