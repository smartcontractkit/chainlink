package messagelimitationstest

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/onramp"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
)

// Expects WithDeployedEnv for ephemeral test environments or WithEnv for long-running test environments like staging.
func NewTestSetup(
	t *testing.T,
	onchainState changeset.CCIPOnChainState,
	sourceChain,
	destChain uint64,
	srctoken common.Address,
	srcFeeQuoterDestChainConfig fee_quoter.FeeQuoterDestChainConfig,
	testRouter,
	validateResp bool,
	opts ...TestSetupOpts,
) TestSetup {
	ts := TestSetup{
		T:                           t,
		OnchainState:                onchainState,
		SrcChain:                    sourceChain,
		DestChain:                   destChain,
		SrcToken:                    srctoken,
		SrcFeeQuoterDestChainConfig: srcFeeQuoterDestChainConfig,
		TestRouter:                  testRouter,
		ValidateResp:                validateResp,
	}

	for _, opt := range opts {
		opt(&ts)
	}

	return ts
}

type TestSetupOpts func(*TestSetup)

func WithDeployedEnv(de testhelpers.DeployedEnv) TestSetupOpts {
	return func(ts *TestSetup) {
		ts.DeployedEnv = &de
		ts.Env = de.Env
	}
}

func WithEnv(env deployment.Environment) TestSetupOpts {
	return func(ts *TestSetup) {
		ts.Env = env
	}
}

type TestSetup struct {
	T                           *testing.T
	Env                         deployment.Environment
	DeployedEnv                 *testhelpers.DeployedEnv
	OnchainState                changeset.CCIPOnChainState
	SrcChain                    uint64
	DestChain                   uint64
	SrcToken                    common.Address
	SrcFeeQuoterDestChainConfig fee_quoter.FeeQuoterDestChainConfig
	TestRouter                  bool
	ValidateResp                bool
}

type TestCase struct {
	TestSetup
	Name      string
	Msg       router.ClientEVM2AnyMessage
	ExpRevert bool
}

type TestCaseOutput struct {
	MsgSentEvent *onramp.OnRampCCIPMessageSent
}

func Run(tc TestCase) TestCaseOutput {
	tc.T.Logf("Sending msg: %s", tc.Name)
	require.NotEqual(tc.T, tc.SrcChain, tc.DestChain, "fromChain and toChain cannot be the same")

	// Approve router to send token only on long-running environments
	if tc.DeployedEnv == nil && tc.SrcToken != (common.Address{}) {
		routerAddress := tc.OnchainState.Chains[tc.SrcChain].Router.Address()
		if tc.TestRouter {
			routerAddress = tc.OnchainState.Chains[tc.SrcChain].TestRouter.Address()
		}
		err := testhelpers.ApproveToken(tc.Env, tc.SrcChain, tc.SrcToken, routerAddress, testhelpers.OneCoin)
		require.NoError(tc.T, err)
	}

	msgSentEvent, err := testhelpers.SendRequest(
		tc.Env, tc.OnchainState,
		testhelpers.WithSourceChain(tc.SrcChain),
		testhelpers.WithDestChain(tc.DestChain),
		testhelpers.WithTestRouter(tc.TestRouter),
		testhelpers.WithEvm2AnyMessage(tc.Msg))

	if tc.ExpRevert {
		tc.T.Logf("Message reverted as expected")
		require.Error(tc.T, err)
		require.Contains(tc.T, err.Error(), "execution reverted")
		return TestCaseOutput{}
	}
	require.NoError(tc.T, err)

	tc.T.Logf("Message not reverted as expected")

	return TestCaseOutput{
		MsgSentEvent: msgSentEvent,
	}
}
