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
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
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
	NoOfReq                   int64 // no of Request fired - required for balance assertion at the end
	totalGEFee                *big.Int
	BalanceStats              BalanceStats  // balance assertion details
	CurrentMsgSerialNo        *atomic.Int64 // current msg serial number in the load sequence
	InitialSourceBlockNum     uint64
	InitialDestBlockNum       uint64        // blocknumber before the first message is fired in the load sequence
	CallTimeOut               time.Duration // max time to wait for various on-chain events
	reports                   *testreporters.CCIPLaneStats
	msg                       router.ClientEVM2AnyMessage
	MaxDataBytes              uint32
	SendMaxDataIntermittently bool
	LastFinalizedTxBlock      atomic.Uint64
	LastFinalizedTimestamp    atomic.Time
}
type BalanceStats struct {
	SourceBalanceReq        map[string]*big.Int
	SourceBalanceAssertions []testhelpers.BalanceAssertion
	DestBalanceReq          map[string]*big.Int
	DestBalanceAssertions   []testhelpers.BalanceAssertion
}

func NewCCIPLoad(t *testing.T, lane *actions.CCIPLane, timeout time.Duration, noOfReq int64, reporter *testreporters.CCIPLaneStats) *CCIPE2ELoad {
	return &CCIPE2ELoad{
		t:                         t,
		Lane:                      lane,
		CurrentMsgSerialNo:        atomic.NewInt64(1),
		CallTimeOut:               timeout,
		NoOfReq:                   noOfReq,
		reports:                   reporter,
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

	// save the current block numbers to use in various filter log requests
	currentBlockOnSource, err := sourceCCIP.Common.ChainClient.LatestBlockNumber(context.Background())
	require.NoError(c.t, err, "failed to fetch latest source block num")
	currentBlockOnDest, err := destCCIP.Common.ChainClient.LatestBlockNumber(context.Background())
	require.NoError(c.t, err, "failed to fetch latest dest block num")
	c.InitialDestBlockNum = currentBlockOnDest
	c.InitialSourceBlockNum = currentBlockOnSource
	// collect the balance requirement to verify balances after transfer
	sourceBalances, err := testhelpers.GetBalances(c.t, sourceCCIP.CollectBalanceRequirements())
	require.NoError(c.t, err, "fetching source balance")
	destBalances, err := testhelpers.GetBalances(c.t, destCCIP.CollectBalanceRequirements())
	require.NoError(c.t, err, "fetching dest balance")
	c.BalanceStats = BalanceStats{
		SourceBalanceReq: sourceBalances,
		DestBalanceReq:   destBalances,
	}
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
	// close all header subscriptions for dest chains
	queuedEvents := destCCIP.Common.ChainClient.GetHeaderSubscriptions()
	for subName := range queuedEvents {
		destCCIP.Common.ChainClient.DeleteHeaderEventSubscription(subName)
	}
	// close all header subscriptions for source chains except for finalized header
	queuedEvents = sourceCCIP.Common.ChainClient.GetHeaderSubscriptions()
	for subName := range queuedEvents {
		if subName == blockchain.FinalizedHeaderKey {
			continue
		}
		sourceCCIP.Common.ChainClient.DeleteHeaderEventSubscription(subName)
	}
}

func (c *CCIPE2ELoad) Call(_ *wasp.Generator) *wasp.CallResult {
	res := &wasp.CallResult{}
	sourceCCIP := c.Lane.Source
	msgSerialNo := c.CurrentMsgSerialNo.Load()
	c.CurrentMsgSerialNo.Inc()

	lggr := c.Lane.Logger.With().Int("msg Number", int(msgSerialNo)).Logger()
	stats := testreporters.NewCCIPRequestStats(msgSerialNo)
	defer c.reports.UpdatePhaseStatsForReq(stats)
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
	msgLog, sourceLogTime, err := c.Lane.Source.AssertEventCCIPSendRequested(lggr, sendTx.Hash().Hex(), c.CallTimeOut, txConfirmationTime, stats)
	if err != nil || msgLog == nil {
		res.Error = err.Error()
		res.Data = stats.StatusByPhase
		res.Failed = true
		return res
	}
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
			lggr, seqNum, msgLog, sourceLogTime, stats)
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
	err = c.Lane.Dest.AssertEventExecutionStateChanged(lggr, seqNum, c.CallTimeOut, blessedAt, stats)
	if err != nil {
		res.Error = err.Error()
		res.Data = stats.StatusByPhase
		res.Failed = true
		return res
	}

	res.Data = stats.StatusByPhase
	return res
}

func (c *CCIPE2ELoad) ReportAcceptedLog() {
	c.Lane.Logger.Info().Msg("Commit Report stats")
	it, err := c.Lane.Dest.CommitStore.Instance.FilterReportAccepted(&bind.FilterOpts{Start: c.InitialDestBlockNum})
	require.NoError(c.t, err, "report committed result")
	i := 1
	event := c.Lane.Logger.Info()
	for it.Next() {
		event.Interface(fmt.Sprintf("%d Report Intervals", i), it.Event.Report.Interval)
		i++
	}
	event.Msgf("CommitStore-Reports Accepted")
}
