package execctxs

import (
	"maps"
	"slices"
	"testing"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/integration-tests/smoke/ccip/canonical/adapters"
	"github.com/smartcontractkit/chainlink/integration-tests/smoke/ccip/canonical/types"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
	"github.com/stretchr/testify/require"
)

func AllOneToOneExecContextNames() []string {
	return []string{
		"evm_2_evm",
		"evm_2_solana",
		"solana_2_evm",
	}
}

// AllOneToOneExecContexts returns a list of all one-source-to-one-dest execution contexts for testing.
func AllOneToOneExecContexts(t *testing.T) []types.ExecContext {
	return []types.ExecContext{
		NewEVM2EVMCtx(t),
		NewEVM2SolanaCtx(t),
		NewSolana2EVM(t),
	}
}

var _ types.ExecContext = &execCtx{}

type execCtx struct {
	name    string
	env     cldf.Environment
	state   stateview.CCIPOnChainState
	sources []types.Adapter
	dest    types.Adapter
}

func (e *execCtx) Name() string {
	return e.name
}

// ReplayLogs implements types.ExecContext.
func (e *execCtx) ReplayLogs(t *testing.T, selectorToBlockMap map[uint64]uint64) {
	testhelpers.ReplayLogs(t, e.env.Offchain, selectorToBlockMap)
}

// Env implements types.ExecContext.
func (e *execCtx) Env() cldf.Environment {
	return e.env
}

// OnchainState implements types.ExecContext.
func (e *execCtx) OnchainState() stateview.CCIPOnChainState {
	return e.state
}

func (e *execCtx) Sources() []types.Adapter {
	return e.sources
}

func (e *execCtx) Dest() types.Adapter {
	return e.dest
}

func NewEVM2EVMCtx(t *testing.T) types.ExecContext {
	e, _, _ := testsetups.NewIntegrationEnvironment(t)

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	allChainSelectors := slices.Collect(maps.Keys(e.Env.Chains))
	require.Len(t, allChainSelectors, 2)
	sourceChain := allChainSelectors[0]
	destChain := allChainSelectors[1]

	// connect a single lane, source to dest
	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)

	return &execCtx{
		name:  "evm_2_evm",
		env:   e.Env,
		state: state,
		sources: []types.Adapter{
			adapters.NewEVMAdapter(e.Env.Chains[sourceChain], state.Chains[sourceChain]),
		},
		dest: adapters.NewEVMAdapter(e.Env.Chains[destChain], state.Chains[destChain]),
	}
}

func NewEVM2SolanaCtx(t *testing.T) types.ExecContext {
	e, _, _ := testsetups.NewIntegrationEnvironment(t, testhelpers.WithSolChains(1))

	// TODO: do this as part of setup
	t.Logf("deploying solana ccip receiver")
	testhelpers.DeploySolanaCcipReceiver(t, e.Env)

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	allChainSelectors := slices.Collect(maps.Keys(e.Env.Chains))
	allSolChainSelectors := slices.Collect(maps.Keys(e.Env.SolChains))
	sourceChain := allChainSelectors[0]
	destChain := allSolChainSelectors[0]

	t.Logf("sourceChain: %d, destChain: %d", sourceChain, destChain)

	// connect a single lane, source to dest
	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)

	return &execCtx{
		name:  "evm_2_solana",
		env:   e.Env,
		state: state,
		sources: []types.Adapter{
			adapters.NewEVMAdapter(e.Env.Chains[sourceChain], state.Chains[sourceChain]),
		},
		dest: adapters.NewSVMAdapter(e.Env.SolChains[destChain], state.SolChains[destChain]),
	}
}

func NewSolana2EVM(t *testing.T) types.ExecContext {
	e, _, _ := testsetups.NewIntegrationEnvironment(t, testhelpers.WithSolChains(1))

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	allChainSelectors := slices.Collect(maps.Keys(e.Env.Chains))
	allSolChainSelectors := slices.Collect(maps.Keys(e.Env.SolChains))
	require.Len(t, allChainSelectors, 2)
	sourceChain := allSolChainSelectors[0]
	destChain := allChainSelectors[1]

	// connect a single lane, source to dest
	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)

	return &execCtx{
		name:  "solana_2_evm",
		env:   e.Env,
		state: state,
		sources: []types.Adapter{
			adapters.NewSVMAdapter(e.Env.SolChains[sourceChain], state.SolChains[sourceChain]),
		},
		dest: adapters.NewEVMAdapter(e.Env.Chains[destChain], state.Chains[destChain]),
	}
}
