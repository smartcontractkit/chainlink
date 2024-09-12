package ccipdeployment

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/devenv"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
)

// Context returns a context with the test's deadline, if available.
func Context(tb testing.TB) context.Context {
	ctx := context.Background()
	var cancel func()
	switch t := tb.(type) {
	case *testing.T:
		if d, ok := t.Deadline(); ok {
			ctx, cancel = context.WithDeadline(ctx, d)
		}
	}
	if cancel == nil {
		ctx, cancel = context.WithCancel(ctx)
	}
	tb.Cleanup(cancel)
	return ctx
}

type DeployedTestEnvironment struct {
	Ab           deployment.AddressBook
	Env          deployment.Environment
	HomeChainSel uint64
	Nodes        map[string]memory.Node
}

// NewDeployedEnvironment creates a new CCIP environment
// with capreg and nodes set up.
func NewDeployedTestEnvironment(t *testing.T, lggr logger.Logger) DeployedTestEnvironment {
	ctx := Context(t)
	chains := memory.NewMemoryChains(t, 3)
	homeChainSel := uint64(0)
	homeChainEVM := uint64(0)
	// Say first chain is home chain.
	for chainSel := range chains {
		homeChainEVM, _ = chainsel.ChainIdFromSelector(chainSel)
		homeChainSel = chainSel
		break
	}
	ab, capReg, err := DeployCapReg(lggr, chains, homeChainSel)
	require.NoError(t, err)

	nodes := memory.NewNodes(t, zapcore.InfoLevel, chains, 4, 1, deployment.RegistryConfig{
		EVMChainID: homeChainEVM,
		Contract:   capReg,
	})
	for _, node := range nodes {
		require.NoError(t, node.App.Start(ctx))
		t.Cleanup(func() {
			require.NoError(t, node.App.Stop())
		})
	}

	e := memory.NewMemoryEnvironmentFromChainsNodes(t, lggr, chains, nodes)
	return DeployedTestEnvironment{
		Ab:           ab,
		Env:          e,
		HomeChainSel: homeChainSel,
		Nodes:        nodes,
	}
}

type DeployedLocalDevEnvironment struct {
	Ab           deployment.AddressBook
	Env          deployment.Environment
	HomeChainSel uint64
	Nodes        []devenv.Node
}

func NewDeployedLocalDevEnvironment(t *testing.T, lggr logger.Logger) DeployedLocalDevEnvironment {
	ctx := Context(t)
	envConfig, testEnv, cfg, deployNodeFunc := devenv.DeployPrivateChains(t)
	require.NotNil(t, envConfig)
	require.NotEmpty(t, envConfig.Chains, "chainConfigs should not be empty")
	require.NotEmpty(t, envConfig.JDConfig, "jdUrl should not be empty")
	chains, err := devenv.NewChains(lggr, envConfig.Chains)
	require.NoError(t, err)
	homeChainSel := uint64(0)
	homeChainEVM := uint64(0)

	// Say first chain is home chain.
	for chainSel := range chains {
		homeChainEVM, _ = chainsel.ChainIdFromSelector(chainSel)
		homeChainSel = chainSel
		break
	}
	ab, capReg, err := DeployCapReg(lggr, chains, homeChainSel)
	require.NoError(t, err)

	err = deployNodeFunc(envConfig, deployment.RegistryConfig{
		EVMChainID: homeChainEVM,
		Contract:   capReg,
	},
		testEnv, cfg)
	require.NoError(t, err)

	e, don, err := devenv.NewEnvironment(ctx, lggr, *envConfig)
	require.NoError(t, err)
	require.NotNil(t, e)
	require.NotNil(t, don)

	// fund the nodes
	require.NoError(t, don.FundNodes(ctx, deployment.E18Mult(10), e.Chains))

	return DeployedLocalDevEnvironment{
		Ab:           ab,
		Env:          *e,
		HomeChainSel: homeChainSel,
		Nodes:        don.Nodes,
	}
}

func AddLanesForAll(e deployment.Environment, state CCIPOnChainState) error {
	for source := range e.Chains {
		for dest := range e.Chains {
			if source != dest {
				err := AddLane(e, state, source, dest)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func SendMessage(
	srcSelector, destSelector uint64,
	transactOpts *bind.TransactOpts,
	srcConfirm func(tx *types.Transaction) (uint64, error),
	state CCIPOnChainState,
) (uint64, error) {
	msg := router.ClientEVM2AnyMessage{
		Receiver:     common.LeftPadBytes(state.Chains[destSelector].Receiver.Address().Bytes(), 32),
		Data:         []byte("hello"),
		TokenAmounts: nil, // TODO: no tokens for now
		FeeToken:     state.Chains[srcSelector].Weth9.Address(),
		ExtraArgs:    nil, // TODO: no extra args for now, falls back to default
	}
	fee, err := state.Chains[srcSelector].Router.GetFee(
		&bind.CallOpts{Context: context.Background()}, destSelector, msg)
	if err != nil {
		return 0, deployment.MaybeDataErr(err)
	}
	tx, err := state.Chains[srcSelector].Weth9.Deposit(&bind.TransactOpts{
		From:   transactOpts.From,
		Signer: transactOpts.Signer,
		Value:  fee,
	})
	if err != nil {
		return 0, err
	}
	_, err = srcConfirm(tx)
	if err != nil {
		return 0, err
	}

	// TODO: should be able to avoid this by using native?
	tx, err = state.Chains[srcSelector].Weth9.Approve(
		transactOpts,
		state.Chains[srcSelector].Router.Address(), fee)
	if err != nil {
		return 0, err
	}
	_, err = srcConfirm(tx)
	if err != nil {
		return 0, err

	}
	tx, err = state.Chains[srcSelector].Router.CcipSend(transactOpts, destSelector, msg)
	if err != nil {
		return 0, err
	}
	return srcConfirm(tx)
}
