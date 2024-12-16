package changeset_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/deployment"

	"github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

// setupSetConfigTestEnv deploys all required contracts for the setConfig MCMS contract call.
func setupSetConfigTestEnv(t *testing.T) deployment.Environment {

	lggr := logger.TestLogger(t)
	cfg := memory.MemoryEnvironmentConfig{
		Nodes:  1,
		Chains: 2,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)
	chainSelector := env.AllChainSelectors()[0]

	config := proposalutils.SingleGroupTimelockConfig(t)
	// Deploy MCMS and Timelock
	env, err := commonchangeset.ApplyChangesets(t, env, nil, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeployLinkToken),
			Config:    []uint64{chainSelector},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeployMCMSWithTimelock),
			Config: map[uint64]types.MCMSWithTimelockConfig{
				chainSelector: config,
			},
		},
	})
	require.NoError(t, err)
	return env
}

// TestSetConfigMCMS tests the SetConfigMCMS changeset by calling SetConfig and checking the config values.
func TestSetConfigMCMS(t *testing.T) {
	t.Parallel()
	ctx := tests.Context(t)

	env := setupSetConfigTestEnv(t)
	chainSelector := env.AllChainSelectors()[0]
	chain := env.Chains[chainSelector]
	addrs, err := env.ExistingAddresses.AddressesForChain(chainSelector)
	require.NoError(t, err)
	require.Len(t, addrs, 6)

	mcmsState, err := commonchangeset.MaybeLoadMCMSWithTimelockChainState(chain, addrs)
	require.NoError(t, err)

	timelockAddress := mcmsState.Timelock.Address()
	cfg := proposalutils.SingleGroupMCMS(t)
	// Add the timelock as a signer to check state changes
	cfg.Signers = append(cfg.Signers, timelockAddress)
	cfg.Quorum = 2 // quorum should change to 2 out of 2 signers

	// Set config on all 3 MCMS contracts
	_, err = commonchangeset.ApplyChangesets(t, env, nil, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.SetConfigMCMS),
			Config: commonchangeset.SetConfigParams{
				ConfigsPerChain: map[uint64]commonchangeset.ConfigPerRole{
					chainSelector: {
						Proposer:  cfg,
						Canceller: cfg,
						Bypasser:  cfg,
					},
				},
			},
		},
	})
	require.NoError(t, err)
	// Check new State
	expected := cfg.ToRawConfig()
	opts := &bind.CallOpts{Context: ctx}
	newConf, err := mcmsState.ProposerMcm.GetConfig(opts)
	require.NoError(t, err)
	require.Equal(t, expected, newConf)

	newConf, err = mcmsState.BypasserMcm.GetConfig(opts)
	require.NoError(t, err)
	require.Equal(t, expected, newConf)

	newConf, err = mcmsState.CancellerMcm.GetConfig(opts)
	require.NoError(t, err)
	require.Equal(t, expected, newConf)
}

// TestSetConfigMCMSProposal tests the SetConfigMCMS changeset proposal generation by calling SetConfig and checking the config values.
func TestSetConfigMCMSProposal(t *testing.T) {
	t.Parallel()
	ctx := tests.Context(t)

	env := setupSetConfigTestEnv(t)
	chainSelector := env.AllChainSelectors()[0]
	chain := env.Chains[chainSelector]
	addrs, err := env.ExistingAddresses.AddressesForChain(chainSelector)
	require.NoError(t, err)
	require.Len(t, addrs, 6)

	mcmsState, err := commonchangeset.MaybeLoadMCMSWithTimelockChainState(chain, addrs)
	require.NoError(t, err)

	timelockAddress := mcmsState.Timelock.Address()
	timelockMap := map[uint64]*proposalutils.TimelockExecutionContracts{
		chainSelector: {
			Timelock:  mcmsState.Timelock,
			CallProxy: mcmsState.CallProxy,
		},
	}
	cfg := proposalutils.SingleGroupMCMS(t)
	// Add the timelock as a signer to check state changes
	cfg.Signers = append(cfg.Signers, timelockAddress)
	cfg.Quorum = 2 // quorum should change to 2 out of 2 signers
	// Apply the changeset
	_, err = commonchangeset.ApplyChangesets(t, env, timelockMap, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.TransferToMCMSWithTimelock),
			Config: commonchangeset.TransferToMCMSWithTimelockConfig{
				ContractsByChain: map[uint64][]common.Address{
					chainSelector: {mcmsState.ProposerMcm.Address(), mcmsState.BypasserMcm.Address(), mcmsState.CancellerMcm.Address()},
				},
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.SetConfigMCMS),
			Config: commonchangeset.SetConfigParams{
				ProposalConfig: &commonchangeset.ProposalConfig{
					MinDelay: 0,
				},
				ConfigsPerChain: map[uint64]commonchangeset.ConfigPerRole{
					chainSelector: {
						Proposer:  cfg,
						Canceller: cfg,
						Bypasser:  cfg,
					},
				},
			},
		},
	})
	require.NoError(t, err)
	// Check new State
	expected := cfg.ToRawConfig()
	opts := &bind.CallOpts{Context: ctx}
	newConf, err := mcmsState.ProposerMcm.GetConfig(opts)
	require.NoError(t, err)
	require.Equal(t, expected, newConf)

	newConf, err = mcmsState.BypasserMcm.GetConfig(opts)
	require.NoError(t, err)
	require.Equal(t, expected, newConf)

	newConf, err = mcmsState.CancellerMcm.GetConfig(opts)
	require.NoError(t, err)
	require.Equal(t, expected, newConf)
}

func TestValidate(t *testing.T) {
	env := setupSetConfigTestEnv(t)

	chainSelector := env.AllChainSelectors()[0]
	chain := env.Chains[chainSelector]
	addrs, err := env.ExistingAddresses.AddressesForChain(chainSelector)
	require.NoError(t, err)
	require.Len(t, addrs, 6)
	mcmsState, err := commonchangeset.MaybeLoadMCMSWithTimelockChainState(chain, addrs)
	require.NoError(t, err)
	cfg := proposalutils.SingleGroupMCMS(t)
	timelockAddress := mcmsState.Timelock.Address()
	// Add the timelock as a signer to check state changes
	cfg.Signers = append(cfg.Signers, timelockAddress)
	cfg.Quorum = 2 // quorum

	cfgInvalid := proposalutils.SingleGroupMCMS(t)
	cfgInvalid.Quorum = 0
	require.NoError(t, err)
	tests := []struct {
		name     string
		cfg      commonchangeset.SetConfigParams
		errorMsg string
	}{
		{
			name: "valid config",
			cfg: commonchangeset.SetConfigParams{
				ProposalConfig: &commonchangeset.ProposalConfig{
					MinDelay: 0,
				},
				ConfigsPerChain: map[uint64]commonchangeset.ConfigPerRole{
					chainSelector: {
						Proposer:  cfg,
						Canceller: cfg,
						Bypasser:  cfg,
					},
				},
			},
		},
		{
			name: "valid non mcms config",
			cfg: commonchangeset.SetConfigParams{
				ConfigsPerChain: map[uint64]commonchangeset.ConfigPerRole{
					chainSelector: {
						Proposer:  cfg,
						Canceller: cfg,
						Bypasser:  cfg,
					},
				},
			},
		},
		{
			name: "no chain configurations",
			cfg: commonchangeset.SetConfigParams{
				ConfigsPerChain: map[uint64]commonchangeset.ConfigPerRole{},
			},
			errorMsg: "no chain configs provided",
		},
		{
			name: "non evm chain",
			cfg: commonchangeset.SetConfigParams{
				ConfigsPerChain: map[uint64]commonchangeset.ConfigPerRole{
					chain_selectors.APTOS_MAINNET.Selector: {
						Proposer:  cfg,
						Canceller: cfg,
						Bypasser:  cfg,
					},
				},
			},
			errorMsg: "chain selector: 4741433654826277614 is not an ethereum chain",
		},
		{
			name: "chain selector not found in environment",
			cfg: commonchangeset.SetConfigParams{
				ConfigsPerChain: map[uint64]commonchangeset.ConfigPerRole{
					123: {
						Proposer:  cfg,
						Canceller: cfg,
						Bypasser:  cfg,
					},
				},
			},
			errorMsg: "unknown chain selector 123",
		},
		{
			name: "invalid proposer config",
			cfg: commonchangeset.SetConfigParams{
				ProposalConfig: &commonchangeset.ProposalConfig{
					MinDelay: 0,
				},
				ConfigsPerChain: map[uint64]commonchangeset.ConfigPerRole{
					chainSelector: {
						Proposer:  cfgInvalid,
						Canceller: cfg,
						Bypasser:  cfg,
					},
				},
			},
			errorMsg: "invalid MCMS config: Quorum must be greater than 0",
		},
		{
			name: "invalid canceller config",
			cfg: commonchangeset.SetConfigParams{
				ProposalConfig: &commonchangeset.ProposalConfig{
					MinDelay: 0,
				},
				ConfigsPerChain: map[uint64]commonchangeset.ConfigPerRole{
					chainSelector: {
						Proposer:  cfg,
						Canceller: cfgInvalid,
						Bypasser:  cfg,
					},
				},
			},
			errorMsg: "invalid MCMS config: Quorum must be greater than 0",
		},
		{
			name: "invalid bypasser config",
			cfg: commonchangeset.SetConfigParams{
				ProposalConfig: &commonchangeset.ProposalConfig{
					MinDelay: 0,
				},
				ConfigsPerChain: map[uint64]commonchangeset.ConfigPerRole{
					chainSelector: {
						Proposer:  cfg,
						Canceller: cfg,
						Bypasser:  cfgInvalid,
					},
				},
			},
			errorMsg: "invalid MCMS config: Quorum must be greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			selectors := []uint64{chainSelector}

			err := tt.cfg.Validate(env, selectors)
			if tt.errorMsg != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
