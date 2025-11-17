package messagelimitationstest

import (
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	aptos_feequoter "github.com/smartcontractkit/chainlink-aptos/bindings/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/onramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_3/fee_quoter"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	module_fee_quoter "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"

	ccipclient "github.com/smartcontractkit/chainlink/deployment/ccip/shared/client"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

// Expects WithDeployedEnv for ephemeral test environments or WithEnv for long-running test environments like staging.
func NewTestSetup(
	t *testing.T,
	onchainState stateview.CCIPOnChainState,
	sourceChain,
	destChain uint64,
	srctoken common.Address,
	srcFeeQuoterDestChainConfig any,
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

	family, err := chain_selectors.GetSelectorFamily(ts.SrcChain)
	require.NoError(ts.T, err)

	switch family {
	case chain_selectors.FamilyEVM:
		evmFeeQuoterDestChainConfig, ok := ts.SrcFeeQuoterDestChainConfig.(fee_quoter.FeeQuoterDestChainConfig)
		require.True(ts.T, ok, "expected Evm Fee quoter destination chain config type")
		ts.SrcFeeQuoterDestChainConfig = evmFeeQuoterDestChainConfig
	case chain_selectors.FamilyAptos:
		aptosFeeQuoterDestChainConfig, ok := ts.SrcFeeQuoterDestChainConfig.(aptos_feequoter.DestChainConfig)
		require.True(ts.T, ok, "expected Aptos Fee quoter destination chain config type")
		ts.SrcFeeQuoterDestChainConfig = aptosFeeQuoterDestChainConfig
	case chain_selectors.FamilySui:
		suiFeeQuoterDestChainConfig, ok := ts.SrcFeeQuoterDestChainConfig.(module_fee_quoter.DestChainConfig)
		require.True(ts.T, ok, "expected Sui Fee quoter destination chain config type")
		ts.SrcFeeQuoterDestChainConfig = suiFeeQuoterDestChainConfig
	default:
		ts.T.Fatalf("unsupported source chain family %v", family)
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

func WithEnv(env cldf.Environment) TestSetupOpts {
	return func(ts *TestSetup) {
		ts.Env = env
	}
}

type TestSetup struct {
	T                           *testing.T
	Env                         cldf.Environment
	DeployedEnv                 *testhelpers.DeployedEnv
	OnchainState                stateview.CCIPOnChainState
	SrcChain                    uint64
	DestChain                   uint64
	SrcToken                    common.Address
	SrcFeeQuoterDestChainConfig any
	TestRouter                  bool
	ValidateResp                bool
}

type TestCase struct {
	TestSetup
	Name      string
	Msg       any
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

	var msgOpt ccipclient.SendReqOpts

	family, err := chain_selectors.GetSelectorFamily(tc.SrcChain)
	require.NoError(tc.T, err)

	switch family {
	case chain_selectors.FamilyEVM:
		evmMsg, ok := tc.Msg.(router.ClientEVM2AnyMessage)
		require.True(tc.T, ok, "expected EVM message type")
		msgOpt = ccipclient.WithMessage(evmMsg)
	case chain_selectors.FamilyAptos:
		aptosMsg, ok := tc.Msg.(testhelpers.AptosSendRequest)
		require.True(tc.T, ok, "expected Aptos message type")
		msgOpt = ccipclient.WithMessage(aptosMsg)
	case chain_selectors.FamilySui:
		suiMsg, ok := tc.Msg.(testhelpers.SuiSendRequest)
		require.True(tc.T, ok, "expected Sui message type")
		msgOpt = ccipclient.WithMessage(suiMsg)
	default:
		tc.T.Fatalf("unsupported source chain family %v", family)
	}

	out, err := testhelpers.SendRequest(
		tc.Env, tc.OnchainState,
		ccipclient.WithSourceChain(tc.SrcChain),
		ccipclient.WithDestChain(tc.DestChain),
		ccipclient.WithTestRouter(tc.TestRouter),
		msgOpt)

	var errorMsg string

	if tc.ExpRevert {
		switch family {
		case chain_selectors.FamilyEVM:
			errorMsg = "execution reverted"
		case chain_selectors.FamilyAptos:
			errorMsg = "transaction reverted:"
		case chain_selectors.FamilySui:
			errorMsg = "sui" // dummy placeholder, handled below
		default:
			tc.T.Fatalf("unsupported source chain family %v", family)
		}

		tc.T.Logf("Message reverted as expected")
		require.Error(tc.T, err)

		if family == chain_selectors.FamilySui {
			require.True(tc.T,
				strings.Contains(err.Error(), "Transaction is rejected") ||
					strings.Contains(err.Error(), "Transaction validator signing failed") ||
					strings.Contains(err.Error(), "transaction failed with error"),
				"unexpected Sui revert: %v", err,
			)
		} else {
			require.Contains(tc.T, err.Error(), errorMsg)
		}
		return TestCaseOutput{}
	}
	require.NoError(tc.T, err)
	msgSentEvent := out.RawEvent.(*onramp.OnRampCCIPMessageSent)

	tc.T.Logf("Message not reverted as expected")

	return TestCaseOutput{
		MsgSentEvent: msgSentEvent,
	}
}
