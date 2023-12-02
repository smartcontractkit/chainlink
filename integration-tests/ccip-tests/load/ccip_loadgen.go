package load

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/wasp"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/testhelpers"
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testreporters"
)

type CCIPE2ELoad struct {
	t                         *testing.T
	Lane                      *actions.CCIPLane
	NoOfReq                   int64         // approx no of Request fired
	CurrentMsgSerialNo        *atomic.Int64 // current msg serial number in the load sequence
	CallTimeOut               time.Duration // max time to wait for various on-chain events
	msg                       router.ClientEVM2AnyMessage
	MaxDataBytes              uint32
	SendMaxDataIntermittently bool
	LastFinalizedTxBlock      atomic.Uint64
	LastFinalizedTimestamp    atomic.Time
}

func NewCCIPLoad(t *testing.T, lane *actions.CCIPLane, timeout time.Duration, noOfReq int64) *CCIPE2ELoad {
	return &CCIPE2ELoad{
		t:                         t,
		Lane:                      lane,
		CurrentMsgSerialNo:        atomic.NewInt64(1),
		CallTimeOut:               timeout,
		NoOfReq:                   noOfReq,
		SendMaxDataIntermittently: false,
	}
}

// BeforeAllCall funds subscription, approves the token transfer amount.
// Needs to be called before load sequence is started.
// Needs to approve and fund for the entire sequence.
func (c *CCIPE2ELoad) BeforeAllCall(msgType string) {
	sourceCCIP := c.Lane.Source
	destCCIP := c.Lane.Dest
	var tokenAndAmounts []router.ClientEVMTokenAmount
	for i := range c.Lane.Source.TransferAmount {
		token := sourceCCIP.Common.BridgeTokens[i]
		tokenAndAmounts = append(tokenAndAmounts, router.ClientEVMTokenAmount{
			Token: common.HexToAddress(token.Address()), Amount: c.Lane.Source.TransferAmount[i],
		})
	}

	err := sourceCCIP.Common.ChainClient.WaitForEvents()
	require.NoError(c.t, err, "Failed to wait for events")

	extraArgsV1, err := testhelpers.GetEVMExtraArgsV1(big.NewInt(100_000), false)
	require.NoError(c.t, err, "Failed encoding the options field")

	receiver, err := utils.ABIEncode(`[{"type":"address"}]`, destCCIP.ReceiverDapp.EthAddress)
	require.NoError(c.t, err, "Failed encoding the receiver address")
	c.msg = router.ClientEVM2AnyMessage{
		Receiver:  receiver,
		ExtraArgs: extraArgsV1,
		FeeToken:  common.HexToAddress(sourceCCIP.Common.FeeToken.Address()),
		Data:      []byte("message with Id 1"),
	}
	if msgType == actions.TokenTransfer {
		c.msg.TokenAmounts = tokenAndAmounts
	}
	if c.SendMaxDataIntermittently {
		dCfg, err := sourceCCIP.OnRamp.Instance.GetDynamicConfig(nil)
		require.NoError(c.t, err, "failed to fetch dynamic config")
		c.MaxDataBytes = dCfg.MaxDataBytes
	}

	// wait for any pending txs before moving on
	err = sourceCCIP.Common.ChainClient.WaitForEvents()
	require.NoError(c.t, err, "Failed to wait for events")
	err = destCCIP.Common.ChainClient.WaitForEvents()
	require.NoError(c.t, err, "Failed to wait for events")
	c.LastFinalizedTxBlock.Store(c.Lane.Source.NewFinalizedBlockNum.Load())
	c.LastFinalizedTimestamp.Store(c.Lane.Source.NewFinalizedBlockTimestamp.Load())

	sourceCCIP.Common.ChainClient.ParallelTransactions(false)
	destCCIP.Common.ChainClient.ParallelTransactions(false)
}

func (c *CCIPE2ELoad) Call(_ *wasp.Generator) *wasp.CallResult {
	res := &wasp.CallResult{}
	sourceCCIP := c.Lane.Source
	msgSerialNo := c.CurrentMsgSerialNo.Load()
	c.CurrentMsgSerialNo.Inc()

	lggr := c.Lane.Logger.With().Int("msg Number", int(msgSerialNo)).Logger()
	stats := testreporters.NewCCIPRequestStats(msgSerialNo)
	defer c.Lane.Reports.UpdatePhaseStatsForReq(stats)
	// form the message for transfer
	msgStr := fmt.Sprintf("new message with Id %d", msgSerialNo)

	if c.SendMaxDataIntermittently {
		lggr.Info().Msg("sending max data intermittently")
		// every 10th message will have extra data with almost MaxDataBytes
		if msgSerialNo%10 == 0 {
			length := c.MaxDataBytes - 1
			b := make([]byte, c.MaxDataBytes-1)
			_, err := rand.Read(b)
			if err != nil {
				res.Error = err.Error()
				res.Failed = true
				return res
			}
			randomString := base64.URLEncoding.EncodeToString(b)
			msgStr = randomString[:length]
		}
	}
	msg := c.msg
	// if msg contains more than 2 tokens, selectively choose random 2 tokens
	if len(msg.TokenAmounts) > 2 {
		// randomize the order of elements in the slice
		rand.Shuffle(len(msg.TokenAmounts), func(i, j int) {
			msg.TokenAmounts[i], msg.TokenAmounts[j] = msg.TokenAmounts[j], msg.TokenAmounts[i]
		})
		// select first 2 tokens
		msg.TokenAmounts = msg.TokenAmounts[:2]
	}
	msg.Data = []byte(msgStr)

	feeToken := sourceCCIP.Common.FeeToken.EthAddress
	// initiate the transfer
	lggr.Debug().Str("triggeredAt", time.Now().GoString()).Msg("triggering transfer")
	var sendTx *types.Transaction
	var err error

	destChainSelector, err := chain_selectors.SelectorFromChainId(sourceCCIP.DestinationChainId)
	if err != nil {
		res.Error = err.Error()
		res.Failed = true
		return res
	}
	// initiate the transfer
	// if the token address is 0x0 it will use Native as fee token and the fee amount should be mentioned in bind.TransactOpts's value

	fee, err := sourceCCIP.Common.Router.GetFee(destChainSelector, msg)
	if err != nil {
		res.Error = err.Error()
		res.Failed = true
		return res
	}
	startTime := time.Now()
	if feeToken != common.HexToAddress("0x0") {
		sendTx, err = sourceCCIP.Common.Router.CCIPSend(destChainSelector, msg, nil)
	} else {
		sendTx, err = sourceCCIP.Common.Router.CCIPSend(destChainSelector, msg, fee)
	}

	if err != nil {
		stats.UpdateState(lggr, 0, testreporters.TX, time.Since(startTime), testreporters.Failure)
		res.Error = fmt.Sprintf("ccip-send tx error %+v for msg ID %d", err, msgSerialNo)
		res.Data = stats.StatusByPhase
		res.Failed = true
		return res
	}
	lggr = lggr.With().Str("Msg Tx", sendTx.Hash().String()).Logger()
	txConfirmationTime := time.Now().UTC()
	rcpt, err1 := bind.WaitMined(context.Background(), sourceCCIP.Common.ChainClient.DeployBackend(), sendTx)
	if err1 == nil {
		hdr, err1 := c.Lane.Source.Common.ChainClient.HeaderByNumber(context.Background(), rcpt.BlockNumber)
		if err1 == nil {
			txConfirmationTime = hdr.Timestamp
		}
	}
	var gasUsed uint64
	if rcpt != nil {
		gasUsed = rcpt.GasUsed
	}
	stats.UpdateState(lggr, 0, testreporters.TX, startTime.Sub(txConfirmationTime), testreporters.Success,
		testreporters.TransactionStats{
			Fee:                fee.String(),
			GasUsed:            gasUsed,
			TxHash:             sendTx.Hash().Hex(),
			NoOfTokensSent:     len(msg.TokenAmounts),
			MessageBytesLength: len(msg.Data),
		})
	// wait for
	// - CCIPSendRequested Event log to be generated,
	msgLogs, sourceLogTime, err := c.Lane.Source.AssertEventCCIPSendRequested(lggr, sendTx.Hash().Hex(), c.CallTimeOut, txConfirmationTime, []*testreporters.RequestStat{stats})

	if err != nil || msgLogs == nil || len(msgLogs) == 0 {
		res.Error = err.Error()
		res.Data = stats.StatusByPhase
		res.Failed = true
		return res
	}
	msgLog := msgLogs[0]
	sentMsg := msgLog.Message
	seqNum := sentMsg.SequenceNumber
	lggr = lggr.With().Str("msgId ", fmt.Sprintf("0x%x", sentMsg.MessageId[:])).Logger()

	if bytes.Compare(sentMsg.Data, []byte(msgStr)) != 0 {
		res.Error = fmt.Sprintf("the message byte didnot match expected %s received %s msg ID %d", msgStr, string(sentMsg.Data), msgSerialNo)
		res.Data = stats.StatusByPhase
		res.Failed = true
		return res
	}

	lstFinalizedBlock := c.LastFinalizedTxBlock.Load()
	var sourceLogFinalizedAt time.Time
	// if the finality tag is enabled and the last finalized block is greater than the block number of the message
	// consider the message finalized
	if c.Lane.Source.Common.ChainClient.GetNetworkConfig().FinalityDepth == 0 &&
		lstFinalizedBlock != 0 && lstFinalizedBlock > msgLog.Raw.BlockNumber {
		sourceLogFinalizedAt = c.LastFinalizedTimestamp.Load()
		stats.UpdateState(lggr, seqNum, testreporters.SourceLogFinalized,
			sourceLogFinalizedAt.Sub(sourceLogTime), testreporters.Success,
			testreporters.TransactionStats{
				TxHash:           msgLog.Raw.TxHash.String(),
				FinalizedByBlock: strconv.FormatUint(lstFinalizedBlock, 10),
				FinalizedAt:      sourceLogFinalizedAt.String(),
			})
	} else {
		var finalizingBlock uint64
		sourceLogFinalizedAt, finalizingBlock, err = c.Lane.Source.AssertSendRequestedLogFinalized(
			lggr, sendTx.Hash(), sourceLogTime, []*testreporters.RequestStat{stats})
		if err != nil {
			res.Error = err.Error()
			res.Data = stats.StatusByPhase
			res.Failed = true
			return res
		}
		c.LastFinalizedTxBlock.Store(finalizingBlock)
		c.LastFinalizedTimestamp.Store(sourceLogFinalizedAt)
	}

	// wait for
	// - CommitStore to increase the seq number,
	err = c.Lane.Dest.AssertSeqNumberExecuted(lggr, seqNum, c.CallTimeOut, sourceLogFinalizedAt, stats)
	if err != nil {
		res.Error = err.Error()
		res.Data = stats.StatusByPhase
		res.Failed = true
		return res
	}
	// wait for ReportAccepted event
	commitReport, reportAcceptedAt, err := c.Lane.Dest.AssertEventReportAccepted(lggr, seqNum, c.CallTimeOut, sourceLogFinalizedAt, stats)
	if err != nil || commitReport == nil {
		res.Error = err.Error()
		res.Data = stats.StatusByPhase
		res.Failed = true
		return res
	}
	blessedAt, err := c.Lane.Dest.AssertReportBlessed(lggr, seqNum, c.CallTimeOut, *commitReport, reportAcceptedAt, stats)
	if err != nil {
		res.Error = err.Error()
		res.Data = stats.StatusByPhase
		res.Failed = true
		return res
	}
	_, err = c.Lane.Dest.AssertEventExecutionStateChanged(lggr, seqNum, c.CallTimeOut, blessedAt, stats, testhelpers.ExecutionStateSuccess)
	if err != nil {
		res.Error = err.Error()
		res.Data = stats.StatusByPhase
		res.Failed = true
		return res
	}

	res.Data = stats.StatusByPhase
	return res
}
