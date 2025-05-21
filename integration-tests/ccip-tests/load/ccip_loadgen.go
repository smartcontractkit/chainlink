package load

import (
	"context"
	crypto_rand "crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"

	"github.com/smartcontractkit/chainlink-common/pkg/config"

	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testconfig"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testreporters"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/testhelpers"
)

// CCIPLaneOptimized is a light-weight version of CCIPLane, It only contains elements which are used during load triggering and validation
type CCIPLaneOptimized struct {
	Logger            *zerolog.Logger
	SourceNetworkName string
	DestNetworkName   string
	Source            *actions.SourceCCIPModule
	Dest              *actions.DestCCIPModule
	Reports           *testreporters.CCIPLaneStats
}

type CCIPE2ELoad struct {
	t                                          *testing.T
	Lane                                       *CCIPLaneOptimized
	NoOfReq                                    int64         // approx no of Request fired
	CurrentMsgSerialNo                         *atomic.Int64 // current msg serial number in the load sequence
	CallTimeOut                                time.Duration // max time to wait for various on-chain events
	msg                                        router.ClientEVM2AnyMessage
	MaxDataBytes                               uint32
	SendMaxDataIntermittentlyInMsgCount        int64
	SkipRequestIfAnotherRequestTriggeredWithin *config.Duration
	LastFinalizedTxBlock                       atomic.Uint64
	LastFinalizedTimestamp                     atomic.Time
	MsgProfiles                                *testconfig.MsgProfile
	EOAReceiver                                []byte
}

func NewCCIPLoad(
	t *testing.T,
	lane *actions.CCIPLane,
	timeout time.Duration,
	noOfReq int64,
	m *testconfig.MsgProfile,
	sendMaxDataIntermittentlyInEveryMsgCount int64,
	SkipRequestIfAnotherRequestTriggeredWithin *config.Duration,
) *CCIPE2ELoad {
	// to avoid holding extra data
	loadLane := &CCIPLaneOptimized{
		Logger:            lane.Logger,
		SourceNetworkName: lane.SourceNetworkName,
		DestNetworkName:   lane.DestNetworkName,
		Source:            lane.Source,
		Dest:              lane.Dest,
		Reports:           lane.Reports,
	}

	return &CCIPE2ELoad{
		t:                                   t,
		Lane:                                loadLane,
		CurrentMsgSerialNo:                  atomic.NewInt64(1),
		CallTimeOut:                         timeout,
		NoOfReq:                             noOfReq,
		SendMaxDataIntermittentlyInMsgCount: sendMaxDataIntermittentlyInEveryMsgCount,
		SkipRequestIfAnotherRequestTriggeredWithin: SkipRequestIfAnotherRequestTriggeredWithin,
		MsgProfiles: m,
	}
}

// BeforeAllCall funds subscription, approves the token transfer amount.
// Needs to be called before load sequence is started.
// Needs to approve and fund for the entire sequence.
func (c *CCIPE2ELoad) BeforeAllCall() {
	sourceCCIP := c.Lane.Source
	destCCIP := c.Lane.Dest

	receiver, err := utils.ABIEncode(`[{"type":"address"}]`, destCCIP.ReceiverDapp.EthAddress)
	require.NoError(c.t, err, "Failed encoding the receiver address")
	c.msg = router.ClientEVM2AnyMessage{
		Receiver: receiver,
		FeeToken: common.HexToAddress(sourceCCIP.Common.FeeToken.Address()),
		Data:     []byte("message with Id 1"),
	}
	var tokenAndAmounts []router.ClientEVMTokenAmount
	if len(c.Lane.Source.Common.BridgeTokens) > 0 {
		for i := range c.Lane.Source.TransferAmount {
			// if length of sourceCCIP.TransferAmount is more than available bridge token use first bridge token
			token := sourceCCIP.Common.BridgeTokens[0]
			if i < len(sourceCCIP.Common.BridgeTokens) {
				token = sourceCCIP.Common.BridgeTokens[i]
			}
			tokenAndAmounts = append(tokenAndAmounts, router.ClientEVMTokenAmount{
				Token: common.HexToAddress(token.Address()), Amount: c.Lane.Source.TransferAmount[i],
			})
		}
		c.msg.TokenAmounts = tokenAndAmounts
	}
	// we might need to change the receiver to the default wallet of destination based on the gaslimit of msg
	// Get the receiver's bytecode to check if it's a contract or EOA
	bytecode, err := c.Lane.Dest.Common.ChainClient.Backend().CodeAt(context.Background(), c.Lane.Dest.ReceiverDapp.EthAddress, nil)
	require.NoError(c.t, err, "Failed to get bytecode of the receiver contract")
	// if the bytecode is empty, it's an EOA,
	// In that case save the receiver address as EOA to be used in the message
	// Otherwise save destination's default wallet address as EOA
	// so that it can be used later for msgs with gaslimit 0
	if len(bytecode) > 0 {
		receiver, err := utils.ABIEncode(`[{"type":"address"}]`, common.HexToAddress(c.Lane.Dest.Common.ChainClient.GetDefaultWallet().Address()))
		require.NoError(c.t, err, "Failed encoding the receiver address")
		c.EOAReceiver = receiver
	} else {
		c.EOAReceiver = c.msg.Receiver
	}
	if c.SendMaxDataIntermittentlyInMsgCount > 0 {
		cfg, err := sourceCCIP.OnRamp.Instance.GetDynamicConfig(nil)
		require.NoError(c.t, err, "failed to fetch dynamic config")
		c.MaxDataBytes = cfg.MaxDataBytes
	}
	// if the msg is sent via multicall, transfer the token transfer amount to multicall contract
	if sourceCCIP.Common.MulticallEnabled &&
		sourceCCIP.Common.MulticallContract != (common.Address{}) &&
		len(c.Lane.Source.Common.BridgeTokens) > 0 {
		for i, amount := range sourceCCIP.TransferAmount {
			// if length of sourceCCIP.TransferAmount is more than available bridge token use first bridge token
			token := sourceCCIP.Common.BridgeTokens[0]
			if i < len(sourceCCIP.Common.BridgeTokens) {
				token = sourceCCIP.Common.BridgeTokens[i]
			}
			amountToApprove := new(big.Int).Mul(amount, big.NewInt(c.NoOfReq))
			bal, err := token.BalanceOf(context.Background(), sourceCCIP.Common.MulticallContract.Hex())
			require.NoError(c.t, err, "Failed to get token balance")
			if bal.Cmp(amountToApprove) < 0 {
				err := token.Transfer(token.OwnerWallet, sourceCCIP.Common.MulticallContract.Hex(), amountToApprove)
				require.NoError(c.t, err, "Failed to approve token transfer amount")
			}
		}
	}

	c.LastFinalizedTxBlock.Store(c.Lane.Source.NewFinalizedBlockNum.Load())
	c.LastFinalizedTimestamp.Store(c.Lane.Source.NewFinalizedBlockTimestamp.Load())

	sourceCCIP.Common.ChainClient.ParallelTransactions(false)
	destCCIP.Common.ChainClient.ParallelTransactions(false)
}

func (c *CCIPE2ELoad) CCIPMsg() (router.ClientEVM2AnyMessage, *testreporters.RequestStat, error) {
	msgSerialNo := c.CurrentMsgSerialNo.Load()
	c.CurrentMsgSerialNo.Inc()
	msgDetails := c.MsgProfiles.MsgDetailsForIteration(msgSerialNo)
	stats := testreporters.NewCCIPRequestStats(msgSerialNo, c.Lane.SourceNetworkName, c.Lane.DestNetworkName)
	// form the message for transfer
	msgLength := pointer.GetInt64(msgDetails.DataLength)
	gasLimit := pointer.GetInt64(msgDetails.DestGasLimit)
	msg := c.msg
	if msgLength > 0 && msgDetails.IsDataTransfer() {
		if c.SendMaxDataIntermittentlyInMsgCount > 0 {
			// every SendMaxDataIntermittentlyInMsgCount message will have extra data with almost MaxDataBytes
			if msgSerialNo%c.SendMaxDataIntermittentlyInMsgCount == 0 {
				msgLength = int64(c.MaxDataBytes - 1)
			}
		}
		b := make([]byte, msgLength)
		_, err := crypto_rand.Read(b)
		if err != nil {
			return router.ClientEVM2AnyMessage{}, stats, fmt.Errorf("failed to generate random string %w", err)
		}
		randomString := base64.URLEncoding.EncodeToString(b)
		msg.Data = []byte(randomString[:msgLength])
	}
	if !msgDetails.IsTokenTransfer() {
		msg.TokenAmounts = []router.ClientEVMTokenAmount{}
	}

	var (
		extraArgs []byte
		err       error
	)
	// v1.5.0 and later starts using V2 extra args
	matchErr := contracts.MatchContractVersionsOrAbove(map[contracts.Name]contracts.Version{
		contracts.OnRampContract: contracts.V1_5_0,
	})
	// TODO: The CCIP Offchain upgrade tests make the AllowOutOfOrder flag tricky in this case.
	// It runs with out of date contract versions at first, then upgrades them. So transactions will assume that the new contracts are there
	// before being deployed. So setting v2 args will break the test. This is a bit of a hack to get around that.
	// The test will soon be deprecated, so a temporary solution is fine.
	if matchErr != nil || !c.Lane.Source.Common.AllowOutOfOrder {
		extraArgs, err = testhelpers.GetEVMExtraArgsV1(big.NewInt(gasLimit), false)
	} else {
		extraArgs, err = testhelpers.GetEVMExtraArgsV2(big.NewInt(gasLimit), c.Lane.Source.Common.AllowOutOfOrder)
	}
	if err != nil {
		return router.ClientEVM2AnyMessage{}, stats, err
	}
	msg.ExtraArgs = extraArgs
	// if gaslimit is 0, set the receiver to EOA
	if gasLimit == 0 {
		msg.Receiver = c.EOAReceiver
	}
	return msg, stats, nil
}

func (c *CCIPE2ELoad) Call(_ *wasp.Generator) *wasp.Response {
	res := &wasp.Response{}
	sourceCCIP := c.Lane.Source
	var recentRequestFoundAt *time.Time
	var latestEvent *types.Log
	var err error
	// Use IsPastRequestTriggeredWithinTimeframe to check for any historical CCIP send request events
	// within the specified timeframe for the first message. Subsequently, use the watcher method to monitor
	// and detect any new events as they occur.
	if c.CurrentMsgSerialNo.Load() == int64(1) {
		latestEvent, err = sourceCCIP.IsPastRequestTriggeredWithinTimeframe(testcontext.Get(c.t), c.SkipRequestIfAnotherRequestTriggeredWithin)
		require.NoError(c.t, err, "error while filtering past requests")
		if latestEvent != nil {
			//nolint:gosec // safe to cast
			hdr, err := sourceCCIP.Common.ChainClient.HeaderByNumber(context.Background(), big.NewInt(int64(latestEvent.BlockNumber)))
			require.NoError(c.t, err, "error while getting header by block number")
			recentRequestFoundAt = pointer.ToTime(hdr.Timestamp)
		}
	} else {
		recentRequestFoundAt = sourceCCIP.IsRequestTriggeredWithinTimeframe(c.SkipRequestIfAnotherRequestTriggeredWithin)
	}

	if recentRequestFoundAt != nil {
		c.Lane.Logger.
			Info().
			Str("Found At=", recentRequestFoundAt.String()).
			Msgf("Skipping ...Another Request found within given timeframe %s", c.SkipRequestIfAnotherRequestTriggeredWithin.String())
		return res
	}
	// if there is an connection error , we will skip sending the request
	// this is to avoid sending the request when the connection is not restored yet
	if sourceCCIP.Common.IsConnectionRestoredRecently != nil {
		if !sourceCCIP.Common.IsConnectionRestoredRecently.Load() {
			c.Lane.Logger.Info().Msg("RPC Connection Error.. skipping this request")
			res.Failed = true
			res.Error = "RPC Connection error .. this request was skipped"
			return res
		}
		c.Lane.Logger.Info().Msg("Connection is restored, Resuming load")
	}
	msg, stats, err := c.CCIPMsg()
	if err != nil {
		res.Error = err.Error()
		res.Failed = true
		return res
	}
	msgSerialNo := stats.ReqNo
	// create a sub-logger for the request
	lggr := c.Lane.Logger.With().Int64("msg Number", stats.ReqNo).Logger()

	feeToken := sourceCCIP.Common.FeeToken.EthAddress
	// initiate the transfer
	lggr.Debug().Str("triggeredAt", time.Now().GoString()).Msg("triggering transfer")
	var sendTx *types.Transaction

	destChainSelector, err := chain_selectors.SelectorFromChainId(sourceCCIP.DestinationChainId)
	if err != nil {
		res.Error = fmt.Sprintf("reqNo %d err %s - while getting selector from chainid", msgSerialNo, err.Error())
		res.Failed = true
		return res
	}

	// initiate the transfer
	// if the token address is 0x0 it will use Native as fee token and the fee amount should be mentioned in bind.TransactOpts's value
	fee, err := sourceCCIP.Common.Router.GetFee(destChainSelector, msg)
	if err != nil {
		res.Error = fmt.Sprintf("reqNo %d err %s - while getting fee from router", msgSerialNo, err.Error())
		res.Failed = true
		return res
	}
	startTime := time.Now().UTC()
	if feeToken != common.HexToAddress("0x0") {
		sendTx, err = sourceCCIP.Common.Router.CCIPSendAndProcessTx(destChainSelector, msg, nil)
	} else {
		// add a bit buffer to fee
		sendTx, err = sourceCCIP.Common.Router.CCIPSendAndProcessTx(destChainSelector, msg, new(big.Int).Add(big.NewInt(1e5), fee))
	}
	if err != nil {
		stats.UpdateState(&lggr, 0, testreporters.TX, time.Since(startTime), testreporters.Failure, nil)
		res.Error = fmt.Sprintf("ccip-send tx error %s for reqNo %d", err.Error(), msgSerialNo)
		res.Data = stats.StatusByPhase
		res.Failed = true
		return res
	}

	// the msg is no longer needed, so we can clear it to avoid holding extra data during load
	//nolint:ineffassign,staticcheck
	msg = router.ClientEVM2AnyMessage{}

	txConfirmationTime := time.Now().UTC()
	lggr = lggr.With().Str("Msg Tx", sendTx.Hash().String()).Logger()

	stats.UpdateState(&lggr, 0, testreporters.TX, txConfirmationTime.Sub(startTime), testreporters.Success, nil)
	err = c.Validate(lggr, sendTx, txConfirmationTime, []*testreporters.RequestStat{stats})
	if err != nil {
		res.Error = err.Error()
		res.Failed = true
		res.Data = stats.StatusByPhase
		return res
	}
	res.Data = stats.StatusByPhase
	return res
}

func (c *CCIPE2ELoad) Validate(lggr zerolog.Logger, sendTx *types.Transaction, txConfirmationTime time.Time, stats []*testreporters.RequestStat) error {
	// wait for
	// - CCIPSendRequested Event log to be generated,
	msgLogs, sourceLogTime, err := c.Lane.Source.AssertEventCCIPSendRequested(&lggr, sendTx.Hash().Hex(), c.CallTimeOut, txConfirmationTime, stats)
	if err != nil {
		return err
	}

	lstFinalizedBlock := c.LastFinalizedTxBlock.Load()
	var sourceLogFinalizedAt time.Time
	// if the finality tag is enabled and the last finalized block is greater than the block number of the message
	// consider the message finalized
	if c.Lane.Source.Common.ChainClient.GetNetworkConfig().FinalityDepth == 0 &&
		lstFinalizedBlock != 0 && lstFinalizedBlock > msgLogs[0].LogInfo.BlockNumber {
		sourceLogFinalizedAt = c.LastFinalizedTimestamp.Load()
		for i, stat := range stats {
			stat.UpdateState(&lggr, stat.SeqNum, testreporters.SourceLogFinalized,
				sourceLogFinalizedAt.Sub(sourceLogTime), testreporters.Success,
				&testreporters.TransactionStats{
					TxHash:             msgLogs[i].LogInfo.TxHash.Hex(),
					FinalizedByBlock:   strconv.FormatUint(lstFinalizedBlock, 10),
					FinalizedAt:        sourceLogFinalizedAt.String(),
					Fee:                msgLogs[i].Fee.String(),
					NoOfTokensSent:     msgLogs[i].NoOfTokens,
					MessageBytesLength: int64(msgLogs[i].DataLength),
					MsgID:              fmt.Sprintf("0x%x", msgLogs[i].MessageId[:]),
				})
		}
	} else {
		var finalizingBlock uint64
		sourceLogFinalizedAt, finalizingBlock, err = c.Lane.Source.AssertSendRequestedLogFinalized(
			&lggr, msgLogs[0].LogInfo.TxHash, msgLogs, sourceLogTime, stats)
		if err != nil {
			return err
		}
		c.LastFinalizedTxBlock.Store(finalizingBlock)
		c.LastFinalizedTimestamp.Store(sourceLogFinalizedAt)
	}

	for _, msgLog := range msgLogs {
		seqNum := msgLog.SequenceNumber
		var reqStat *testreporters.RequestStat
		lggr = lggr.With().Str("MsgID", fmt.Sprintf("0x%x", msgLog.MessageId[:])).Logger()
		for _, stat := range stats {
			if stat.SeqNum == seqNum {
				reqStat = stat
				break
			}
		}
		if reqStat == nil {
			return fmt.Errorf("could not find request stat for seq number %d", seqNum)
		}
		// wait for
		// - CommitStore to increase the seq number,
		err = c.Lane.Dest.AssertSeqNumberExecuted(&lggr, seqNum, c.CallTimeOut, sourceLogFinalizedAt, reqStat)
		if err != nil {
			return err
		}
		// wait for ReportAccepted event
		commitReport, reportAcceptedAt, err := c.Lane.Dest.AssertEventReportAccepted(&lggr, seqNum, c.CallTimeOut, sourceLogFinalizedAt, reqStat)
		if err != nil || commitReport == nil {
			return err
		}
		blessedAt, err := c.Lane.Dest.AssertReportBlessed(&lggr, seqNum, c.CallTimeOut, *commitReport, reportAcceptedAt, reqStat)
		if err != nil {
			return err
		}
		_, err = c.Lane.Dest.AssertEventExecutionStateChanged(&lggr, seqNum, c.CallTimeOut, blessedAt, reqStat, testhelpers.ExecutionStateSuccess)
		if err != nil {
			return err
		}
	}

	return nil
}
