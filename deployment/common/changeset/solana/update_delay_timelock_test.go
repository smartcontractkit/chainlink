package solana

import (
	"fmt"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	mcmsSolana "github.com/smartcontractkit/mcms/sdk/solana"
	"github.com/stretchr/testify/require"

	timelockBindings "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/timelock"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

func TestUpdateTimelockDelaySolana_VerifyPreconditions(t *testing.T) {
	t.Parallel()

	selector1 := chainselectors.TEST_22222222222222222222222222222222222222222222.Selector
	selector2 := chainselectors.TEST_33333333333333333333333333333333333333333333.Selector

	env, err := environment.New(t.Context(),
		environment.WithSolanaContainer(t, []uint64{selector1, selector2}, t.TempDir(), map[string]string{}),
	)
	require.NoError(t, err)

	// Setup selector1 to have a chain where the timelock program has been deployed.
	timelockID := mcmsSolana.ContractAddress(
		solana.NewWallet().PublicKey(),
		[32]byte{'t', 'e', 's', 't'},
	)

	err = env.ExistingAddresses.Save(selector1, timelockID, cldf.TypeAndVersion{
		Type:    types.RBACTimelock,
		Version: deployment.Version1_0_0,
	})
	require.NoError(t, err)

	// Setup selector2 to have a chain where the timelock program has not been deployed.
	// e.g. missing the required addresses so that the state loader returns empty seeds.
	mcmDummyProgram := solana.NewWallet().PublicKey()
	mcmsProposerIDEmpty := mcmsSolana.ContractAddress(
		mcmDummyProgram,
		[32]byte{},
	)

	err = env.ExistingAddresses.Save(selector2, mcmsProposerIDEmpty, cldf.TypeAndVersion{
		Type:    types.BypasserManyChainMultisig,
		Version: deployment.Version1_0_0,
	})
	require.NoError(t, err)

	// Create an environment with no solana chains
	emptyEnv, err := environment.New(t.Context())
	require.NoError(t, err)

	tests := []struct {
		name          string
		env           cldf.Environment
		config        UpdateTimelockDelaySolanaCfg
		expectedError string
	}{
		{
			name: "All preconditions satisfied",
			env:  *env,
			config: UpdateTimelockDelaySolanaCfg{
				DelayPerChain: map[uint64]time.Duration{selector1: 5 * time.Minute},
			},
			expectedError: "",
		},
		{
			name: "No Solana chains found in environment",
			env:  *emptyEnv,
			config: UpdateTimelockDelaySolanaCfg{
				DelayPerChain: map[uint64]time.Duration{selector1: 5 * time.Minute},
			},
			expectedError: "no solana chains provided",
		},
		{
			name: "Chain selector not found in environment",
			env:  *env,
			config: UpdateTimelockDelaySolanaCfg{
				DelayPerChain: map[uint64]time.Duration{9999: 5 * time.Minute},
			},
			expectedError: "solana chain not found for selector 9999",
		},
		{
			name: "Timelock not deployed (empty seeds)",
			env:  *env,
			config: UpdateTimelockDelaySolanaCfg{
				DelayPerChain: map[uint64]time.Duration{selector2: 5 * time.Minute},
			},
			expectedError: fmt.Sprintf("timelock program not deployed for chain %d", selector2),
		},
		{
			name:          "empty config provided",
			env:           *env,
			config:        UpdateTimelockDelaySolanaCfg{},
			expectedError: "no delay configs provided",
		},
	}

	cs := UpdateTimelockDelaySolana{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cs.VerifyPreconditions(tt.env, tt.config)
			if tt.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedError)
			}
		})
	}
}

func TestUpdateTimelockDelaySolana_Apply(t *testing.T) {
	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/DX-762")
	t.Parallel()

	rt, selector := setupTest(t)
	require.Len(t, rt.Environment().BlockChains.SolanaChains(), 1)

	newDelayDuration := 5 * time.Minute

	// Run the UpdateTimelockDelaySolana changeset
	err := rt.Exec(
		runtime.ChangesetTask(UpdateTimelockDelaySolana{}, UpdateTimelockDelaySolanaCfg{
			DelayPerChain: map[uint64]time.Duration{
				selector: newDelayDuration,
			},
		}),
	)
	require.NoError(t, err)

	addresses, err := rt.State().AddressBook.AddressesForChain(selector)
	require.NoError(t, err)

	chain := rt.Environment().BlockChains.SolanaChains()[selector]

	// Check new delay config value
	mcmState, err := state.MaybeLoadMCMSWithTimelockChainStateSolana(chain, addresses)
	require.NoError(t, err)

	timelockConfigPDA := state.GetTimelockConfigPDA(mcmState.TimelockProgram, mcmState.TimelockSeed)

	var timelockConfig timelockBindings.Config
	err = chain.GetAccountDataBorshInto(t.Context(), timelockConfigPDA, &timelockConfig)
	require.NoError(t, err)
	require.Equal(t, timelockConfig.MinDelay, uint64(newDelayDuration.Seconds()))
}
