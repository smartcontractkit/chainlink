package changeset

import (
	"math"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	"github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
)

func TestValidateEthBalMonSetWatchListConfig(t *testing.T) {
	t.Parallel()

	selector := chainselectors.TEST_90000001.Selector
	selectorOther := chainselectors.TEST_90000002.Selector

	env, err := environment.New(t.Context(),
		environment.WithEVMSimulated(t, []uint64{selector}),
	)
	require.NoError(t, err)

	addr := common.HexToAddress(testAddr1)
	minOne := big.NewInt(1)
	topOne := big.NewInt(1)
	neg := big.NewInt(-1)

	tests := []struct {
		name      string
		cfg       types.EthBalMonSetWatchListInput
		wantError bool
		errorMsg  string
	}{
		{
			name:      "empty chains",
			cfg:       types.EthBalMonSetWatchListInput{Chains: map[uint64]types.EthBalMonSetWatchListChainConfig{}},
			wantError: true,
			errorMsg:  "no chains provided",
		},
		{
			name: "unknown chain selector",
			cfg: types.EthBalMonSetWatchListInput{
				Chains: map[uint64]types.EthBalMonSetWatchListChainConfig{
					math.MaxUint64: {
						Addresses:       []common.Address{addr},
						MinBalancesWei:  []*big.Int{minOne},
						TopUpAmountsWei: []*big.Int{topOne},
					},
				},
			},
			wantError: true,
			errorMsg:  "not found in environment",
		},
		{
			name: "chain not in environment",
			cfg: types.EthBalMonSetWatchListInput{
				Chains: map[uint64]types.EthBalMonSetWatchListChainConfig{
					selectorOther: {
						Addresses:       []common.Address{addr},
						MinBalancesWei:  []*big.Int{minOne},
						TopUpAmountsWei: []*big.Int{topOne},
					},
				},
			},
			wantError: true,
			errorMsg:  "not found in environment",
		},
		{
			name: "empty addresses",
			cfg: types.EthBalMonSetWatchListInput{
				Chains: map[uint64]types.EthBalMonSetWatchListChainConfig{
					selector: {
						Addresses:       nil,
						MinBalancesWei:  nil,
						TopUpAmountsWei: nil,
					},
				},
			},
			wantError: true,
			errorMsg:  "addresses must not be empty",
		},
		{
			name: "slice length mismatch",
			cfg: types.EthBalMonSetWatchListInput{
				Chains: map[uint64]types.EthBalMonSetWatchListChainConfig{
					selector: {
						Addresses:       []common.Address{addr, addr},
						MinBalancesWei:  []*big.Int{minOne},
						TopUpAmountsWei: []*big.Int{topOne, topOne},
					},
				},
			},
			wantError: true,
			errorMsg:  "must have the same length",
		},
		{
			name: "negative min balance",
			cfg: types.EthBalMonSetWatchListInput{
				Chains: map[uint64]types.EthBalMonSetWatchListChainConfig{
					selector: {
						Addresses:       []common.Address{addr},
						MinBalancesWei:  []*big.Int{neg},
						TopUpAmountsWei: []*big.Int{topOne},
					},
				},
			},
			wantError: true,
			errorMsg:  "min_balance_wei at index 0 must be >= 0",
		},
		{
			name: "negative top-up",
			cfg: types.EthBalMonSetWatchListInput{
				Chains: map[uint64]types.EthBalMonSetWatchListChainConfig{
					selector: {
						Addresses:       []common.Address{addr},
						MinBalancesWei:  []*big.Int{minOne},
						TopUpAmountsWei: []*big.Int{neg},
					},
				},
			},
			wantError: true,
			errorMsg:  "topup_amounts_wei at index 0 must be >= 0",
		},
		{
			name: "valid",
			cfg: types.EthBalMonSetWatchListInput{
				Chains: map[uint64]types.EthBalMonSetWatchListChainConfig{
					selector: {
						Addresses:       []common.Address{addr},
						MinBalancesWei:  []*big.Int{minOne},
						TopUpAmountsWei: []*big.Int{topOne},
					},
				},
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := ValidateEthBalMonSetWatchListConfig(t.Context(), *env, tt.cfg)
			if tt.wantError {
				require.Error(t, err)
				if tt.errorMsg != "" {
					require.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestEthBalMonSetWatchListChangeset(t *testing.T) {
	t.Parallel()

	selector := chainselectors.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
	))
	require.NoError(t, err)

	setupMCMSInfrastructure(t, rt, []uint64{selector})
	fundDeployerAccounts(t, rt.Environment(), []uint64{selector})

	deployCfg := types.DeployEthBalMonInput{
		Chains: map[uint64]types.DeployEthBalMonChainConfig{
			selector: {SetKeeperRegistryAddress: testAddr1},
		},
	}
	deployTask := runtime.ChangesetTask(DeployEthBalMonChangeSet, deployCfg)
	require.NoError(t, rt.Exec(deployTask))

	watchAddr := common.HexToAddress(testAddr2)
	cfg := types.EthBalMonSetWatchListInput{
		Chains: map[uint64]types.EthBalMonSetWatchListChainConfig{
			selector: {
				Addresses:       []common.Address{watchAddr},
				MinBalancesWei:  []*big.Int{big.NewInt(1)},
				TopUpAmountsWei: []*big.Int{big.NewInt(2)},
			},
		},
	}
	require.NoError(t, EthBalMonSetWatchList.VerifyPreconditions(rt.Environment(), cfg))

	watchTask := runtime.ChangesetTask(EthBalMonSetWatchList, cfg)
	require.NoError(t, rt.Exec(watchTask))

	out := rt.State().Outputs[watchTask.ID()]
	require.NotEmpty(t, out.MCMSTimelockProposals)
	prop := out.MCMSTimelockProposals[0]
	require.Contains(t, prop.Description, "EthBalMon SetWatchList")
	require.Len(t, prop.Operations, 1)
	require.Len(t, prop.Operations[0].Transactions, 1)
	require.Contains(t, prop.Operations[0].Transactions[0].Tags, "setWatchList")
}

func TestEthBalMonSetWatchList_VerifyPreconditions_rejectsEmptyChains(t *testing.T) {
	t.Parallel()

	selector := chainselectors.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
	))
	require.NoError(t, err)

	err = EthBalMonSetWatchList.VerifyPreconditions(rt.Environment(), types.EthBalMonSetWatchListInput{
		Chains: map[uint64]types.EthBalMonSetWatchListChainConfig{},
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "no chains provided")
}

func TestEthBalMonSetWatchList_Apply_withoutEthBalMonInDatastore(t *testing.T) {
	t.Parallel()

	selector := chainselectors.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
	))
	require.NoError(t, err)

	setupMCMSInfrastructure(t, rt, []uint64{selector})
	fundDeployerAccounts(t, rt.Environment(), []uint64{selector})

	cfg := types.EthBalMonSetWatchListInput{
		Chains: map[uint64]types.EthBalMonSetWatchListChainConfig{
			selector: {
				Addresses:       []common.Address{common.HexToAddress(testAddr1)},
				MinBalancesWei:  []*big.Int{big.NewInt(1)},
				TopUpAmountsWei: []*big.Int{big.NewInt(1)},
			},
		},
	}
	require.NoError(t, EthBalMonSetWatchList.VerifyPreconditions(rt.Environment(), cfg))

	_, err = EthBalMonSetWatchList.Apply(rt.Environment(), cfg)
	require.Error(t, err)
}
