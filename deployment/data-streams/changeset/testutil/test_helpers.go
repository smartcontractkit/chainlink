package testutil

import (
	"math/big"
	"testing"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/common/types"

	"github.com/smartcontractkit/chainlink/deployment"
	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	dsChangesets "github.com/smartcontractkit/chainlink/deployment/data-streams/changeset"
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
		_, err := changeset.Apply(t, env, nil,
			changeset.Configure(
				deployment.CreateLegacyChangeSet(changeset.DeployMCMSWithTimelockV2),
				map[uint64]types.MCMSWithTimelockConfigV2{
					chainSelector: config,
				},
			),
		)
		require.NoError(t, err)
	}

	return env
}

func DeployMCMS(
	t *testing.T,
	e deployment.Environment,
	addrs ...map[uint64][]common.Address,
) (env deployment.Environment, mcmsState *commonChangesets.MCMSWithTimelockState, timelocks map[uint64]*proposalutils.TimelockExecutionContracts) {
	t.Helper()

	chainSelector := TestChain.Selector
	config := proposalutils.SingleGroupMCMSV2(t)

	env, err := commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			deployment.CreateLegacyChangeSet(commonChangesets.DeployLinkToken),
			[]uint64{chainSelector},
		),
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

	if len(addrs) > 0 {
		env, err = commonChangesets.Apply(
			t, env, timelocks,
			commonChangesets.Configure(
				deployment.CreateLegacyChangeSet(commonChangesets.TransferToMCMSWithTimelockV2),
				commonChangesets.TransferToMCMSWithTimelockConfig{
					ContractsByChain: addrs[0],
					MinDelay:         0,
				},
			),
		)
		require.NoError(t, err)
	}

	return env, mcmsState, timelocks
}

func GetMCMSConfig(useMCMS bool) *dsChangesets.MCMSConfig {
	if useMCMS {
		return &dsChangesets.MCMSConfig{MinDelay: 0, OverrideRoot: true}
	}
	return nil
}
