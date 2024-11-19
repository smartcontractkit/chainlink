package smoke

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink/deployment"
	ccdeploy "github.com/smartcontractkit/chainlink/deployment/ccip"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testsetups"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/multicall3"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func Test_CCIPBatching(t *testing.T) {
	// Setup 3 chains, with 2 lanes going to the dest.
	lggr := logger.TestLogger(t)
	ctx := ccdeploy.Context(t)
	// Will load 3 chains when specified by the overrides.toml or env vars (E2E_TEST_SELECTED_NETWORK).
	// See e2e-tests.yml.
	e, _, _ := testsetups.NewLocalDevEnvironmentWithDefaultPrice(t, lggr)

	state, err := ccdeploy.LoadOnchainState(e.Env)
	require.NoError(t, err)

	allChainSelectors := maps.Keys(e.Env.Chains)
	require.Len(t, allChainSelectors, 3, "this test expects 3 chains")
	sourceChain1 := allChainSelectors[0]
	sourceChain2 := allChainSelectors[1]
	destChain := allChainSelectors[2]
	t.Log("All chain selectors:", allChainSelectors,
		", home chain selector:", e.HomeChainSel,
		", feed chain selector:", e.FeedChainSel,
		", source chain selector 1:", sourceChain1,
		", source chain selector 2:", sourceChain2,
		", dest chain selector:", destChain,
	)
	output, err := changeset.DeployPrerequisites(e.Env, changeset.DeployPrerequisiteConfig{
		ChainSelectors: e.Env.AllChainSelectors(),
	})
	require.NoError(t, err)
	require.NoError(t, e.Env.ExistingAddresses.Merge(output.AddressBook))

	tokenConfig := ccdeploy.NewTestTokenConfig(state.Chains[e.FeedChainSel].USDFeeds)
	// Apply migration
	output, err = changeset.InitialDeploy(e.Env, ccdeploy.DeployCCIPContractConfig{
		HomeChainSel:   e.HomeChainSel,
		FeedChainSel:   e.FeedChainSel,
		ChainsToDeploy: allChainSelectors,
		TokenConfig:    tokenConfig,
		MCMSConfig:     ccdeploy.NewTestMCMSConfig(t, e.Env),
		OCRSecrets:     deployment.XXXGenerateTestOCRSecrets(),
	})
	require.NoError(t, err)
	require.NoError(t, e.Env.ExistingAddresses.Merge(output.AddressBook))
	// Get new state after migration.
	state, err = ccdeploy.LoadOnchainState(e.Env)
	require.NoError(t, err)

	// Ensure capreg logs are up to date.
	ccdeploy.ReplayLogs(t, e.Env.Offchain, e.ReplayBlocks)

	// Apply the jobs.
	for nodeID, jobs := range output.JobSpecs {
		for _, job := range jobs {
			// Note these auto-accept
			_, err := e.Env.Offchain.ProposeJob(ctx,
				&jobv1.ProposeJobRequest{
					NodeId: nodeID,
					Spec:   job,
				})
			require.NoError(t, err)
		}
	}

	// connect sourceChain1 and sourceChain2 to destChain
	require.NoError(t, ccdeploy.AddLaneWithDefaultPrices(e.Env, state, sourceChain1, destChain))
	require.NoError(t, ccdeploy.AddLaneWithDefaultPrices(e.Env, state, sourceChain2, destChain))

	const (
		numMessages = 5
	)

	t.Run("batch data only messages from multiple sources", func(t *testing.T) {
		var wg sync.WaitGroup

		wg.Add(1)
		go func(sourceChainSelector uint64) {
			defer wg.Done()
			sendMessages(
				ctx,
				t,
				e.Env.Chains[sourceChainSelector],
				e.Env.Chains[sourceChainSelector].DeployerKey,
				state.Chains[sourceChainSelector].OnRamp,
				state.Chains[sourceChainSelector].Router,
				state.Chains[sourceChainSelector].Multicall3,
				destChain,
				numMessages,
				common.LeftPadBytes(state.Chains[destChain].Receiver.Address().Bytes(), 32),
			)
		}(sourceChain1)

		wg.Add(1)
		go func(sourceChainSelector uint64) {
			defer wg.Done()
			sendMessages(
				ctx,
				t,
				e.Env.Chains[sourceChainSelector],
				e.Env.Chains[sourceChainSelector].DeployerKey,
				state.Chains[sourceChainSelector].OnRamp,
				state.Chains[sourceChainSelector].Router,
				state.Chains[sourceChainSelector].Multicall3,
				destChain,
				numMessages,
				common.LeftPadBytes(state.Chains[destChain].Receiver.Address().Bytes(), 32),
			)
		}(sourceChain2)

		wg.Wait()

		// confirm the commit reports
		var (
			sourceChain1Report *offramp.OffRampCommitReportAccepted
			sourceChain2Report *offramp.OffRampCommitReportAccepted
		)
		wg.Add(1)
		go func() {
			defer wg.Done()
			var err error
			sourceChain1Report, err = ccdeploy.ConfirmCommitWithExpectedSeqNumRange(t,
				e.Env.Chains[sourceChain1],
				e.Env.Chains[destChain],
				state.Chains[destChain].OffRamp,
				nil,
				ccipocr3.NewSeqNumRange(ccipocr3.SeqNum(1), ccipocr3.SeqNum(numMessages)),
			)
			require.NoErrorf(t, err, "failed to confirm commit from chain %d", sourceChain1)
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			var err error
			sourceChain2Report, err = ccdeploy.ConfirmCommitWithExpectedSeqNumRange(t,
				e.Env.Chains[sourceChain2],
				e.Env.Chains[destChain],
				state.Chains[destChain].OffRamp,
				nil,
				ccipocr3.NewSeqNumRange(ccipocr3.SeqNum(1), ccipocr3.SeqNum(numMessages)),
			)
			require.NoErrorf(t, err, "failed to confirm commit from chain %d", sourceChain2)
		}()

		t.Log("waiting for commit report")
		wg.Wait()

		// the reports should be the same for both, since both roots should be batched within
		// that one report.
		require.Equal(t, sourceChain1Report, sourceChain2Report, "commit reports should be the same")
	})

	t.Run("batch mix of data only messages and token messages from multiple sources", func(t *testing.T) {
		// Generate some messages, on each source.
		// Send them from a multicall contract, i.e multiple ccip messages in the same tx.
		// assert they are committed in the same batch.
	})
}

func sendMessages(
	ctx context.Context,
	t *testing.T,
	sourceChain deployment.Chain,
	sourceTransactOpts *bind.TransactOpts,
	sourceOnRamp *onramp.OnRamp,
	sourceRouter *router.Router,
	sourceMulticall3 *multicall3.Multicall3,
	destChainSelector uint64,
	numMessages int,
	receiver []byte,
) {
	calls, totalValue := genMessages(
		ctx,
		t,
		sourceRouter,
		destChainSelector,
		numMessages,
		receiver,
	)

	// Send the tx with the messages through the multicall
	tx, err := sourceMulticall3.Aggregate3Value(
		&bind.TransactOpts{
			From:   sourceTransactOpts.From,
			Signer: sourceTransactOpts.Signer,
			Value:  totalValue,
		},
		calls,
	)
	_, err = deployment.ConfirmIfNoError(sourceChain, tx, err)
	require.NoError(t, err, "failed to confirm tx")

	// check that the message was emitted
	iter, err := sourceOnRamp.FilterCCIPMessageSent(
		nil, []uint64{destChainSelector}, nil,
	)
	require.NoError(t, err)

	// there should be numMessages messages emitted
	for i := 0; i < numMessages; i++ {
		require.Truef(t, iter.Next(), "expected %d messages, got %d", numMessages, i+1)
		t.Logf("Message id of msg %d: %x", i, iter.Event.Message.Header.MessageId[:])
	}
}

func genMessages(
	ctx context.Context,
	t *testing.T,
	sourceRouter *router.Router,
	destChainSelector uint64,
	count int,
	receiver []byte,
) (calls []multicall3.Multicall3Call3Value, totalValue *big.Int) {
	totalValue = big.NewInt(0)
	for i := 0; i < count; i++ {
		msg := router.ClientEVM2AnyMessage{
			Receiver:     receiver,
			Data:         []byte(fmt.Sprintf("hello world %d", i)),
			TokenAmounts: nil,
			FeeToken:     common.HexToAddress("0x0"),
			ExtraArgs:    nil,
		}

		fee, err := sourceRouter.GetFee(&bind.CallOpts{Context: ctx}, destChainSelector, msg)
		require.NoError(t, err)

		totalValue.Add(totalValue, fee)

		calls = append(calls, multicall3.Multicall3Call3Value{
			Target:       sourceRouter.Address(),
			AllowFailure: false,
			CallData:     ccdeploy.CCIPSendCalldata(t, destChainSelector, msg),
			Value:        fee,
		})
	}

	return calls, totalValue
}
