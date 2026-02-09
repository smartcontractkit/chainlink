package ccip

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/onramp"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	ccipclient "github.com/smartcontractkit/chainlink/deployment/ccip/shared/client"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
)

func Test_CCIP_with_RMN_enabled(t *testing.T) {
	tenv, _, _ := testsetups.NewIntegrationEnvironment(t,
		testhelpers.WithNumOfChains(3),
		testhelpers.WithNumOfUsersPerChain(2),
	)

	e := tenv.Env
	state, err := stateview.LoadOnchainState(e)
	require.NoError(t, err)

	// add all lanes
	testhelpers.AddLanesForAll(t, &tenv, state)

	var (
		chains                 = e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilyEVM))
		chainA, chainB, chainC = chains[0], chains[1], chains[2]
		evmChains              = e.BlockChains.EVMChains()
		expectedSeqNumExec     = make(map[testhelpers.SourceDestPair][]uint64)
		startBlocks            = make(map[uint64]*uint64)
		sendmessage            = func(src, dest uint64, deployer *bind.TransactOpts) (*onramp.OnRampCCIPMessageSent, error) {
			out, err := testhelpers.SendRequest(
				e,
				state,
				ccipclient.WithSender(deployer),
				ccipclient.WithSourceChain(src),
				ccipclient.WithDestChain(dest),
				ccipclient.WithTestRouter(false),
				ccipclient.WithMessage(router.ClientEVM2AnyMessage{
					Receiver:     common.LeftPadBytes(state.MustGetEVMChainState(dest).Receiver.Address().Bytes(), 32),
					Data:         []byte("hello"),
					TokenAmounts: nil,
					FeeToken:     common.HexToAddress("0x0"),
					ExtraArgs:    nil,
				}))
			if err != nil {
				return nil, err
			}
			return out.RawEvent.(*onramp.OnRampCCIPMessageSent), nil
		}
	)
	latestDestHeader, err := evmChains[chainC].Client.HeaderByNumber(t.Context(), nil)
	require.NoError(t, err)
	latestDestBlock := latestDestHeader.Number.Uint64()
	startBlocks[chainC] = &latestDestBlock

	// send message from chainB to chainC
	messageSentEvent, err := sendmessage(chainB, chainC, evmChains[chainB].DeployerKey)
	require.NoError(t, err)
	expectedSeqNumExec[testhelpers.SourceDestPair{
		SourceChainSelector: chainB,
		DestChainSelector:   chainC,
	}] = []uint64{messageSentEvent.SequenceNumber}

	// change the chainA to chainC config to enable RMN
	cs := []commoncs.ConfiguredChangeSet{commoncs.Configure(
		cldf.CreateLegacyChangeSet(v1_6.UpdateOffRampSourcesChangeset),
		v1_6.UpdateOffRampSourcesConfig{
			UpdatesByChain: map[uint64]map[uint64]v1_6.OffRampSourceUpdate{
				chainC: {
					chainA: {
						IsEnabled:                 true,
						IsRMNVerificationDisabled: false,
					},
				},
			},
		})}
	e, _, err = commoncs.ApplyChangesets(t, e, cs)
	require.NoError(t, err)

	// send message from chainB to chainC again, make sure it succeeds
	messageSentEvent, err = sendmessage(chainB, chainC, evmChains[chainB].DeployerKey)
	require.NoError(t, err)
	expectedSeqNumExec[testhelpers.SourceDestPair{
		SourceChainSelector: chainB,
		DestChainSelector:   chainC,
	}] = []uint64{messageSentEvent.SequenceNumber}
	testhelpers.ConfirmExecWithSeqNrsForAll(t, e, state, expectedSeqNumExec, startBlocks)

	// send message from chainA to chainC
	messageSentEvent, err = sendmessage(chainA, chainC, evmChains[chainA].DeployerKey)
	// make sure there's no onchain revert because of source config
	require.NoError(t, err)
	expectedSeqNumExec[testhelpers.SourceDestPair{
		SourceChainSelector: chainA,
		DestChainSelector:   chainC,
	}] = []uint64{messageSentEvent.SequenceNumber}

	// change the chainA to chainC config to disable RMN
	cs = []commoncs.ConfiguredChangeSet{commoncs.Configure(
		cldf.CreateLegacyChangeSet(v1_6.UpdateOffRampSourcesChangeset),
		v1_6.UpdateOffRampSourcesConfig{
			UpdatesByChain: map[uint64]map[uint64]v1_6.OffRampSourceUpdate{
				chainC: {
					chainA: {
						IsEnabled:                 true,
						IsRMNVerificationDisabled: true,
					},
				},
			},
		})}
	e, _, err = commoncs.ApplyChangesets(t, e, cs)
	require.NoError(t, err)

	// send message from chainA to chainC, make sure both messages succeed
	messageSentEvent, err = sendmessage(chainA, chainC, evmChains[chainA].DeployerKey)
	require.NoError(t, err)
	expectedSeqNumExec[testhelpers.SourceDestPair{
		SourceChainSelector: chainA,
		DestChainSelector:   chainC,
	}] = []uint64{messageSentEvent.SequenceNumber}
	testhelpers.ConfirmExecWithSeqNrsForAll(t, e, state, expectedSeqNumExec, startBlocks)
}
