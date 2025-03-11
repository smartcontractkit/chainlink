package messagelimitationstest

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/fee_quoter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/onramp"
	"github.com/stretchr/testify/require"
)

// Use this when testhelpers.DeployedEnv is not available (usually in long-running test environments like staging).
func NewTestSetup(
	t *testing.T,
	env deployment.Environment,
	onchainState changeset.CCIPOnChainState,
	sourceChain,
	destChain uint64,
	srctoken common.Address,
	srcFeeQuoterDestChainConfig fee_quoter.FeeQuoterDestChainConfig,
	testRouter,
	validateResp bool,
) TestSetup {
	return TestSetup{
		T:   t,
		Env: env,
		// no DeployedEnv
		OnchainState:                onchainState,
		SrcChain:                    sourceChain,
		DestChain:                   destChain,
		SrcToken:                    srctoken,
		SrcFeeQuoterDestChainConfig: srcFeeQuoterDestChainConfig,
		TestRouter:                  testRouter,
		ValidateResp:                validateResp,
	}
}

type TestSetup struct {
	T                           *testing.T
	Env                         deployment.Environment
	DeployedEnv                 testhelpers.DeployedEnv
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

	// Approve router to send token
	if tc.SrcToken != (common.Address{}) {
		routerAddress := tc.OnchainState.Chains[tc.SrcChain].TestRouter.Address()
		err := testhelpers.ApproveToken(tc.Env, tc.SrcChain, tc.SrcToken, routerAddress, testhelpers.OneCoin)
		require.NoError(tc.T, err)
	}

	msgSentEvent, err := testhelpers.DoSendRequest(
		tc.T, tc.Env, tc.OnchainState,
		testhelpers.WithSourceChain(tc.SrcChain),
		testhelpers.WithDestChain(tc.DestChain),
		testhelpers.WithTestRouter(true),
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
