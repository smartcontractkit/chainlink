package testutil

import (
	"math/big"
	"testing"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	commonstate "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/common/types"

	"github.com/smartcontractkit/chainlink/deployment"
	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	dsTypes "github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// TestChain is the chain used by the in-memory environment.
var TestChain = chainselectors.Chain{
	EvmChainID: 90000001,
	Selector:   909606746561742123,
	Name:       "Test Chain",
	VarName:    "",
}

// NewMemoryEnv Deploys a memory environment with the provided number of nodes and optionally deploys MCMS and Timelock.
// Deprecated: use NewMemoryEnvV2 instead.
func NewMemoryEnv(t *testing.T, deployMCMS bool, optionalNumNodes ...int) deployment.Environment {
	lggr := logger.TestLogger(t)

	// Default to 0 if no extra argument is provided
	numNodes := 0
	if len(optionalNumNodes) > 0 {
		numNodes = optionalNumNodes[0]
	}

	memEnvConf := memory.MemoryEnvironmentConfig{
		Chains: 1,
		Nodes:  numNodes,
	}

	env := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memEnvConf)
	chainSelector := env.AllChainSelectors()[0]

	if deployMCMS {
		config := proposalutils.SingleGroupTimelockConfigV2(t)
		// Deploy MCMS and Timelock
		_, err := commonChangesets.Apply(t, env, nil,
			commonChangesets.Configure(
				deployment.CreateLegacyChangeSet(commonChangesets.DeployMCMSWithTimelockV2),
				map[uint64]types.MCMSWithTimelockConfigV2{
					chainSelector: config,
				},
			),
		)
		require.NoError(t, err)
	}

	return env
}

type MemoryEnvConfig struct {
	ShouldDeployMCMS      bool
	ShouldDeployLinkToken bool
	NumNodes              int
}

type MemoryEnv struct {
	Environment    deployment.Environment
	Timelocks      map[uint64]*proposalutils.TimelockExecutionContracts
	LinkTokenState *commonstate.LinkTokenState
}

// NewMemoryEnvV2 Deploys a memory environment with configuration and returns an environment wrapper with metadata
func NewMemoryEnvV2(t *testing.T, cfg MemoryEnvConfig) MemoryEnv {
	lggr := logger.TestLogger(t)

	memEnvConf := memory.MemoryEnvironmentConfig{
		Chains: 1,
		Nodes:  cfg.NumNodes,
	}

	env := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memEnvConf)
	chainSelector := env.AllChainSelectors()[0]
	chain := env.Chains[chainSelector]

	var linkTokenState *commonstate.LinkTokenState
	if cfg.ShouldDeployLinkToken {
		updatedEnv, err := commonChangesets.Apply(t, env, nil,
			commonChangesets.Configure(
				deployment.CreateLegacyChangeSet(commonChangesets.DeployLinkToken),
				[]uint64{chainSelector},
			),
		)
		require.NoError(t, err)
		addresses, err := updatedEnv.ExistingAddresses.AddressesForChain(chainSelector)
		require.NoError(t, err)
		env = updatedEnv
		linkState, err := commonstate.MaybeLoadLinkTokenChainState(chain, addresses)
		require.NoError(t, err)
		require.NotNil(t, linkState.LinkToken)
		linkTokenState = linkState
	}

	timelocks := make(map[uint64]*proposalutils.TimelockExecutionContracts)
	if cfg.ShouldDeployMCMS {
		config := proposalutils.SingleGroupTimelockConfigV2(t)
		// Deploy MCMS and Timelock
		updatedEnv, err := commonChangesets.Apply(t, env, nil,
			commonChangesets.Configure(
				deployment.CreateLegacyChangeSet(commonChangesets.DeployMCMSWithTimelockV2),
				map[uint64]types.MCMSWithTimelockConfigV2{
					chainSelector: config,
				},
			),
		)
		require.NoError(t, err)

		addresses, err := updatedEnv.ExistingAddresses.AddressesForChain(TestChain.Selector)
		require.NoError(t, err)

		mcmsState, err := commonstate.MaybeLoadMCMSWithTimelockChainState(chain, addresses)
		require.NoError(t, err)

		timelocks = map[uint64]*proposalutils.TimelockExecutionContracts{
			chainSelector: {
				Timelock:  mcmsState.Timelock,
				CallProxy: mcmsState.CallProxy,
			},
		}
		env = updatedEnv
	}

	return MemoryEnv{
		Environment:    env,
		Timelocks:      timelocks,
		LinkTokenState: linkTokenState,
	}
}

// Deploy MCMS and Timelock, optionally transferring ownership of the provided contracts to Timelock
func DeployMCMS(
	t *testing.T,
	e deployment.Environment,
	addressesToTransfer ...map[uint64][]common.Address,
) (env deployment.Environment, mcmsState *commonChangesets.MCMSWithTimelockState, timelocks map[uint64]*proposalutils.TimelockExecutionContracts) {
	t.Helper()

	chainSelector := TestChain.Selector
	config := proposalutils.SingleGroupMCMSV2(t)

	env, err := commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			deployment.CreateLegacyChangeSet(commonChangesets.DeployMCMSWithTimelockV2),
			map[uint64]types.MCMSWithTimelockConfigV2{
				chainSelector: {
					Canceller:        config,
					Bypasser:         config,
					Proposer:         config,
					TimelockMinDelay: big.NewInt(0),
				},
			},
		),
	)

	require.NoError(t, err)

	addresses, err := e.ExistingAddresses.AddressesForChain(TestChain.Selector)
	require.NoError(t, err)

	chain := e.Chains[chainSelector]

	mcmsState, err = commonChangesets.MaybeLoadMCMSWithTimelockChainState(chain, addresses)
	require.NoError(t, err)

	timelocks = map[uint64]*proposalutils.TimelockExecutionContracts{
		chainSelector: {
			Timelock:  mcmsState.Timelock,
			CallProxy: mcmsState.CallProxy,
		},
	}

	if len(addressesToTransfer) > 0 {
		env, err = commonChangesets.Apply(
			t, env, timelocks,
			commonChangesets.Configure(
				deployment.CreateLegacyChangeSet(commonChangesets.TransferToMCMSWithTimelockV2),
				commonChangesets.TransferToMCMSWithTimelockConfig{
					ContractsByChain: addressesToTransfer[0],
					MCMSConfig:       proposalutils.TimelockConfig{MinDelay: 0},
				},
			),
		)
		require.NoError(t, err)
	}

	return env, mcmsState, timelocks
}

func GetMCMSConfig(useMCMS bool) *dsTypes.MCMSConfig {
	if useMCMS {
		return &dsTypes.MCMSConfig{MinDelay: 0, OverrideRoot: true}
	}
	return nil
}
