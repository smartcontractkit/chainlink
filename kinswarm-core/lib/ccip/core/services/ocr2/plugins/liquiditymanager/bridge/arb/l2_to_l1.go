package arb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/arb_node_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/arbitrum_l1_bridge_adapter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/arbitrum_rollup_core"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/arbsys"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/l2_arbitrum_gateway"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/l2_arbitrum_messenger"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/liquiditymanager"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/abiutils"
	bridgecommon "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/bridge/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

type l2ToL1Bridge struct {
	localSelector      models.NetworkSelector
	remoteSelector     models.NetworkSelector
	l1LiquidityManager liquiditymanager.LiquidityManagerInterface
	l2LiquidityManager liquiditymanager.LiquidityManagerInterface
	l2LogPoller        logpoller.LogPoller
	l1LogPoller        logpoller.LogPoller
	l2FilterName       string
	l1FilterName       string
	lggr               logger.Logger
	l2Client           client.Client
	arbSys             arbsys.ArbSysInterface
	l2ArbGateway       l2_arbitrum_gateway.L2ArbitrumGatewayInterface
	l2ArbMessenger     l2_arbitrum_messenger.L2ArbitrumMessengerInterface
	rollupCore         arbitrum_rollup_core.ArbRollupCoreInterface
	nodeInterface      arb_node_interface.NodeInterfaceInterface
	l1Token, l2Token   common.Address
}

func NewL2ToL1Bridge(
	ctx context.Context,
	lggr logger.Logger,
	localSelector,
	remoteSelector models.NetworkSelector,
	l1RollupAddress,
	l1LiquidityManagerAddress,
	l2LiquidityManagerAddress common.Address,
	l2LogPoller,
	l1LogPoller logpoller.LogPoller,
	l2Client,
	l1Client client.Client,
) (*l2ToL1Bridge, error) {
	localChain, ok := chainsel.ChainBySelector(uint64(localSelector))
	if !ok {
		return nil, fmt.Errorf("unknown chain selector for local chain: %d", localSelector)
	}
	remoteChain, ok := chainsel.ChainBySelector(uint64(remoteSelector))
	if !ok {
		return nil, fmt.Errorf("unknown chain selector for remote chain: %d", remoteSelector)
	}
	l2FilterName := bridgecommon.GetBridgeFilterName(
		"ArbitrumL2ToL1Bridge",
		"L2",
		l2LiquidityManagerAddress,
		localChain.Name,
		remoteChain.Name,
		"",
	)
	err := l2LogPoller.RegisterFilter(
		ctx,
		logpoller.Filter{
			Name: l2FilterName,
			EventSigs: []common.Hash{
				bridgecommon.LiquidityTransferredTopic,
			},
			Addresses: []common.Address{l2LiquidityManagerAddress},
			Retention: bridgecommon.DurationMonth,
		})
	if err != nil {
		return nil, fmt.Errorf("register filter for Arbitrum L2 to L1 bridge: %w", err)
	}

	l1FilterName := bridgecommon.GetBridgeFilterName(
		"ArbitrumL2ToL1Bridge",
		"L1",
		l1LiquidityManagerAddress,
		localChain.Name,
		remoteChain.Name,
		fmt.Sprintf("Rollup:%s", l1RollupAddress.Hex()),
	)
	err = l1LogPoller.RegisterFilter(
		ctx,
		logpoller.Filter{
			Name: l1FilterName,
			EventSigs: []common.Hash{
				NodeConfirmedTopic,                     // emitted by rollup
				bridgecommon.LiquidityTransferredTopic, // emitted by rebalancer
			},
			Addresses: []common.Address{
				l1RollupAddress,           // to get node confirmed logs
				l1LiquidityManagerAddress, // to get LiquidityTransferred logs
			},
			Retention: bridgecommon.DurationMonth,
		})
	if err != nil {
		return nil, fmt.Errorf("register filter for Arbitrum L1 to L2 bridge: %w", err)
	}

	l1LiquidityManager, err := liquiditymanager.NewLiquidityManager(l1LiquidityManagerAddress, l1Client)
	if err != nil {
		return nil, fmt.Errorf("instantiate L1 rebalancer: %w", err)
	}

	l1XchainRebal, err := l1LiquidityManager.GetCrossChainRebalancer(nil, uint64(localSelector))
	if err != nil {
		return nil, fmt.Errorf("get L1->L2 bridge adapter address: %w", err)
	}

	l2LiquidityManager, err := liquiditymanager.NewLiquidityManager(l2LiquidityManagerAddress, l2Client)
	if err != nil {
		return nil, fmt.Errorf("instantiate L2 rebalancer: %w", err)
	}

	l2XchainRebal, err := l2LiquidityManager.GetCrossChainRebalancer(nil, uint64(remoteSelector))
	if err != nil {
		return nil, fmt.Errorf("get L2->L1 bridge adapter address: %w", err)
	}

	arbSys, err := arbsys.NewArbSys(ArbSysAddress, l2Client)
	if err != nil {
		return nil, fmt.Errorf("instantiate ArbSys contract: %w", err)
	}

	// Addresses provided to the below wrappers don't matter,
	// we're just using them to parse the needed logs easily.
	l2ArbGateway, err := l2_arbitrum_gateway.NewL2ArbitrumGateway(
		utils.ZeroAddress,
		l2Client,
	)
	if err != nil {
		return nil, fmt.Errorf("instantiate L2ArbitrumGateway contract: %w", err)
	}

	l2ArbMessenger, err := l2_arbitrum_messenger.NewL2ArbitrumMessenger(
		utils.ZeroAddress,
		l2Client,
	)
	if err != nil {
		return nil, fmt.Errorf("instantiate L2ArbitrumMessenger contract: %w", err)
	}

	// have to use the correct address here
	rollupCore, err := arbitrum_rollup_core.NewArbRollupCore(l1RollupAddress, l1Client)
	if err != nil {
		return nil, fmt.Errorf("instantiate ArbRollupCore contract: %w", err)
	}

	// and here
	nodeInterface, err := arb_node_interface.NewNodeInterface(NodeInterfaceAddress, l2Client)
	if err != nil {
		return nil, fmt.Errorf("instantiate NodeInterface contract: %w", err)
	}

	l2Token, err := l2LiquidityManager.ILocalToken(nil)
	if err != nil {
		return nil, fmt.Errorf("get L2 token address: %w", err)
	}
	l1Token, err := l1LiquidityManager.ILocalToken(nil)
	if err != nil {
		return nil, fmt.Errorf("get L1 token address: %w", err)
	}

	lggr = lggr.Named("ArbitrumL2ToL1Bridge").With(
		"localSelector", localSelector,
		"remoteSelector", remoteSelector,
		"localChainID", localChain.EvmChainID,
		"remoteChainID", remoteChain.EvmChainID,
		"localChainName", localChain.Name,
		"remoteChainName", remoteChain.Name,
		"l1BridgeAdapter", l1XchainRebal.LocalBridge,
		"l2BridgeAdapter", l2XchainRebal.LocalBridge,
		"l1LiquidityManager", l1LiquidityManagerAddress,
		"l1Token", l1Token,
		"l2Token", l2Token,
	)
	lggr.Infow("Initialized arbitrum L2 -> L1 bridge")

	// TODO: replay log poller for any missed logs?
	return &l2ToL1Bridge{
		localSelector:      localSelector,
		remoteSelector:     remoteSelector,
		l2LogPoller:        l2LogPoller,
		l1LogPoller:        l1LogPoller,
		l2FilterName:       l2FilterName,
		l1FilterName:       l1FilterName,
		l1LiquidityManager: l1LiquidityManager,
		l2LiquidityManager: l2LiquidityManager,
		lggr:               lggr,
		l2Client:           l2Client,
		arbSys:             arbSys,
		l2ArbGateway:       l2ArbGateway,
		l2ArbMessenger:     l2ArbMessenger,
		rollupCore:         rollupCore,
		nodeInterface:      nodeInterface,
		l1Token:            l1Token,
		l2Token:            l2Token,
	}, nil
}

func (l *l2ToL1Bridge) QuorumizedBridgePayload(payloads [][]byte, f int) ([]byte, error) {
	// there's no payload for arbitrum L2 -> L1 transfers
	return []byte{}, nil
}

// GetBridgePayloadAndFee implements bridge.Bridge.
// Arbitrum L2 to L1 transfers require no bridge specific payload.
func (l *l2ToL1Bridge) GetBridgePayloadAndFee(ctx context.Context, transfer models.Transfer) ([]byte, *big.Int, error) {
	return []byte{}, big.NewInt(0), nil
}

// Close implements bridge.Bridge.
func (l *l2ToL1Bridge) Close(ctx context.Context) error {
	return multierr.Combine(
		l.l2LogPoller.UnregisterFilter(ctx, l.l2FilterName),
		l.l1LogPoller.UnregisterFilter(ctx, l.l1FilterName),
	)
}

// GetTransfers implements bridge.Bridge.
func (l *l2ToL1Bridge) GetTransfers(ctx context.Context, localToken models.Address, remoteToken models.Address) ([]models.PendingTransfer, error) {
	lggr := l.lggr.With("l2Token", localToken, "l1Token", remoteToken)
	if l.l2Token.Cmp(common.Address(localToken)) != 0 {
		return nil, fmt.Errorf("local token mismatch: expected %s, got %s", l.l2Token, localToken)
	}
	if l.l1Token.Cmp(common.Address(remoteToken)) != 0 {
		return nil, fmt.Errorf("remote token mismatch: expected %s, got %s", l.l1Token, remoteToken)
	}
	// get all the L2 -> L1 transfers in the past 14 days for the given l2Token.
	// that should be enough time to catch all the transfers that were potentially not finalized.
	// TODO: make more performant. Perhaps filter on more than just one topic here to avoid doing in-memory filtering.
	sendLogs, err := l.l2LogPoller.IndexedLogsCreatedAfter(
		ctx,
		bridgecommon.LiquidityTransferredTopic,
		l.l2LiquidityManager.Address(),
		bridgecommon.LiquidityTransferredToChainSelectorTopicIndex,
		[]common.Hash{
			bridgecommon.NetworkSelectorToHash(l.remoteSelector),
		},
		// todo: this should not be hardcoded
		// todo: heavy query warning
		time.Now().Add(-bridgecommon.DurationMonth/2),
		evmtypes.Finalized,
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("get L2 -> L1 transfers from log poller (on L2): %w", err)
	}

	// get all L2 -> L1 finalizations in the past 14 days
	// Note: we don't filter on finalized because we want to avoid marking a sent tx as
	// ready to finalize more than once, since that will cause reverts onchain.
	receiveLogs, err := l.l1LogPoller.IndexedLogsCreatedAfter(
		ctx,
		bridgecommon.LiquidityTransferredTopic,
		l.l1LiquidityManager.Address(),
		bridgecommon.LiquidityTransferredFromChainSelectorTopicIndex,
		[]common.Hash{
			bridgecommon.NetworkSelectorToHash(l.localSelector),
		},
		time.Now().Add(-bridgecommon.DurationMonth/2),
		1,
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("get L2 -> L1 finalizations from log poller (on L1): %w", err)
	}

	lggr.Infow("Got L2 -> L1 transfers and finalizations",
		"l2ToL1Transfers", sendLogs,
		"l2ToL1Finalizations", receiveLogs,
	)

	parsedSent, parsedToLP, err := bridgecommon.ParseLiquidityTransferred(l.l1LiquidityManager.ParseLiquidityTransferred, sendLogs)
	if err != nil {
		return nil, fmt.Errorf("parse L2 -> L1 transfers: %w", err)
	}

	parsedReceived, _, err := bridgecommon.ParseLiquidityTransferred(l.l1LiquidityManager.ParseLiquidityTransferred, receiveLogs)
	if err != nil {
		return nil, fmt.Errorf("parse L2 -> L1 finalizations: %w", err)
	}

	ready, readyData, notReady, err := l.partitionReadyTransfers(ctx, parsedSent, parsedReceived)
	if err != nil {
		return nil, fmt.Errorf("partition ready transfers: %w", err)
	}

	return l.toPendingTransfers(localToken, remoteToken, ready, readyData, notReady, parsedToLP)
}

func (l *l2ToL1Bridge) toPendingTransfers(
	localToken, remoteToken models.Address,
	ready []*liquiditymanager.LiquidityManagerLiquidityTransferred,
	readyData [][]byte,
	notReady []*liquiditymanager.LiquidityManagerLiquidityTransferred,
	parsedToLP map[bridgecommon.LogKey]logpoller.Log,
) ([]models.PendingTransfer, error) {
	if len(ready) != len(readyData) {
		return nil, fmt.Errorf("length of ready and readyData should be the same: len(ready) = %d, len(readyData) = %d",
			len(ready), len(readyData))
	}
	var transfers []models.PendingTransfer
	for i, transfer := range ready {
		transfers = append(transfers, models.PendingTransfer{
			Transfer: models.Transfer{
				From:               l.localSelector,
				To:                 l.remoteSelector,
				Sender:             models.Address(l.l2LiquidityManager.Address()),
				Receiver:           models.Address(l.l1LiquidityManager.Address()),
				LocalTokenAddress:  localToken,
				RemoteTokenAddress: remoteToken,
				Amount:             ubig.New(transfer.Amount),
				Date: parsedToLP[bridgecommon.LogKey{
					TxHash:   transfer.Raw.TxHash,
					LogIndex: int64(transfer.Raw.Index),
				}].BlockTimestamp,
				BridgeData:      readyData[i], // finalization data for withdrawals that are ready
				Stage:           bridgecommon.StageFinalizeReady,
				NativeBridgeFee: ubig.NewI(0),
			},
			Status: models.TransferStatusReady,
			ID:     fmt.Sprintf("%s-%d", transfer.Raw.TxHash.Hex(), transfer.Raw.Index),
		})
	}
	for _, transfer := range notReady {
		transfers = append(transfers, models.PendingTransfer{
			Transfer: models.Transfer{
				From:               l.localSelector,
				To:                 l.remoteSelector,
				Sender:             models.Address(l.l2LiquidityManager.Address()),
				Receiver:           models.Address(l.l1LiquidityManager.Address()),
				LocalTokenAddress:  localToken,
				RemoteTokenAddress: remoteToken,
				Amount:             ubig.New(transfer.Amount),
				Date: parsedToLP[bridgecommon.LogKey{
					TxHash:   transfer.Raw.TxHash,
					LogIndex: int64(transfer.Raw.Index),
				}].BlockTimestamp,
				BridgeData:      []byte{}, // No data since its not ready
				Stage:           bridgecommon.StageRebalanceConfirmed,
				NativeBridgeFee: ubig.NewI(0),
			},
			Status: models.TransferStatusNotReady,
			ID:     fmt.Sprintf("%s-%d", transfer.Raw.TxHash.Hex(), transfer.Raw.Index),
		})
	}
	// TODO: need to also return executed finalizations. See https://smartcontract-it.atlassian.net/browse/CCIP-1893.
	// Use stage StageFinalizeConfirmed for executed finalizations.
	return transfers, nil
}

func (l *l2ToL1Bridge) partitionReadyTransfers(
	ctx context.Context,
	sentLogs,
	receivedLogs []*liquiditymanager.LiquidityManagerLiquidityTransferred,
) (
	ready []*liquiditymanager.LiquidityManagerLiquidityTransferred,
	readyDatas [][]byte,
	notReady []*liquiditymanager.LiquidityManagerLiquidityTransferred,
	err error,
) {
	unfinalized, err := filterUnfinalizedTransfers(sentLogs, receivedLogs)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("filter unfinalized transfers: %w", err)
	}

	var errs error
	for _, transfer := range unfinalized {
		readyData, readyToFinalize, err := l.getFinalizationData(ctx, transfer)
		if err != nil {
			errs = multierr.Append(
				errs,
				fmt.Errorf("get finalization data for transfer %s: %w", transfer.Raw.TxHash, err),
			)
			continue
		}
		if readyToFinalize {
			l.lggr.Infow("transfer is ready to finalize!",
				"transfer", transfer.Raw.TxHash,
				"readyData", hexutil.Encode(readyData),
			)
			ready = append(ready, transfer)
			readyDatas = append(readyDatas, readyData)
		} else {
			notReady = append(notReady, transfer)
		}
	}
	return
}

func filterUnfinalizedTransfers(sentLogs, receivedLogs []*liquiditymanager.LiquidityManagerLiquidityTransferred) ([]*liquiditymanager.LiquidityManagerLiquidityTransferred, error) {
	var unfinalized []*liquiditymanager.LiquidityManagerLiquidityTransferred
	for _, sent := range sentLogs {
		var found bool
		for _, recv := range receivedLogs {
			finalizationPayload, err := unpackFinalizationPayload(recv.BridgeSpecificData)
			if err != nil {
				return nil, fmt.Errorf("unpack finalization payload (bridgeSpecificData) from recv event: %w", err)
			}
			l2ToL1Id, err := abiutils.UnpackUint256(sent.BridgeReturnData)
			if err != nil {
				return nil, fmt.Errorf("unpack l2ToL1TxId (bridgeReturnData) from send event: %w", err)
			}
			if finalizationPayload.Index.Cmp(l2ToL1Id) == 0 {
				found = true
				break
			}
		}
		if !found {
			unfinalized = append(unfinalized, sent)
		}
	}
	return unfinalized, nil
}

func (l *l2ToL1Bridge) getFinalizationData(
	ctx context.Context,
	transfer *liquiditymanager.LiquidityManagerLiquidityTransferred,
) (
	[]byte,
	bool,
	error,
) {
	l.lggr.Infow("Getting finalization data for transfer", "transfer", transfer)
	receipt, err := l.l2Client.TransactionReceipt(ctx, transfer.Raw.TxHash)
	if err != nil {
		// should be a transient error
		return nil, false, fmt.Errorf("get transaction receipt: %w", err)
	}
	var (
		l2ToL1TxLog, withdrawalInitiatedLog, txToL1Log *gethtypes.Log
	)
	for _, lg := range receipt.Logs {
		if len(lg.Topics) <= 0 {
			l.lggr.Warnw("getFinalizationData: log has no topics, skipping", "txHash", lg.TxHash)
			continue
		}
		if lg.Topics[0] == L2ToL1TxTopic {
			l2ToL1TxLog = lg
		} else if lg.Topics[0] == WithdrawalInitiatedTopic {
			withdrawalInitiatedLog = lg
		} else if lg.Topics[0] == TxToL1Topic {
			txToL1Log = lg
		}
	}
	if l2ToL1TxLog == nil || withdrawalInitiatedLog == nil || txToL1Log == nil {
		return nil, false, fmt.Errorf("missing one or more logs: l2ToL1TxLog: %+v, withdrawalInitiatedLog: %+v, txToL1Log: %+v",
			l2ToL1TxLog, withdrawalInitiatedLog, txToL1Log)
	}
	l2ToL1Tx, err := l.arbSys.ParseL2ToL1Tx(*l2ToL1TxLog)
	if err != nil {
		return nil, false, fmt.Errorf("parse L2ToL1Tx log in tx %s: %w", receipt.TxHash, err)
	}
	withdrawalInitiated, err := l.l2ArbGateway.ParseWithdrawalInitiated(*withdrawalInitiatedLog)
	if err != nil {
		return nil, false, fmt.Errorf("parse WithdrawalInitiated log in tx %s: %w", receipt.TxHash, err)
	}
	txToL1, err := l.l2ArbMessenger.ParseTxToL1(*txToL1Log)
	if err != nil {
		return nil, false, fmt.Errorf("parse TxToL1 log in tx %s: %w", receipt.TxHash, err)
	}
	l.lggr.Infow("Got logs for transfer, generating args", "l2ToL1Tx", l2ToL1Tx, "withdrawalInitiated", withdrawalInitiated, "txToL1", txToL1)
	// argument 0: proof
	arg0Proof, err := l.getProof(ctx, withdrawalInitiated.L2ToL1Id)
	if err != nil {
		return nil, false, fmt.Errorf("get proof: %w, l2tol1id: %s",
			err, withdrawalInitiated.L2ToL1Id)
	}
	if arg0Proof == nil {
		// if there's no proof, it means the transfer is not yet ready to finalize
		return nil, false, nil
	}
	// argument 1: index
	arg1Index := withdrawalInitiated.L2ToL1Id
	// argument 2: l2Sender
	arg2L2Sender := withdrawalInitiatedLog.Address
	// argument 3: to
	arg3To := txToL1.To
	// argument 4: l1Block
	arg4L1Block, err := l.getL1BlockFromRPC(ctx, receipt.TxHash)
	if err != nil {
		return nil, false, fmt.Errorf("get l1 block for tx (%s) from rpc: %w",
			receipt.TxHash, err)
	}
	// argument 5: l2Block
	arg5L2Block := receipt.BlockNumber
	// argument 6: l2Timestamp
	arg6L2Timestamp := l2ToL1Tx.Timestamp
	// argument 7: value
	arg7Value := withdrawalInitiated.Amount
	// argument 8: data
	arg8Data := txToL1.Data

	finalizationPayload, err := l1AdapterABI.Pack("exposeArbitrumFinalizationPayload", arbitrum_l1_bridge_adapter.ArbitrumL1BridgeAdapterArbitrumFinalizationPayload{
		Proof:       arg0Proof,
		Index:       arg1Index,
		L2Sender:    arg2L2Sender,
		To:          arg3To,
		L1Block:     arg4L1Block,
		L2Block:     arg5L2Block,
		L2Timestamp: arg6L2Timestamp,
		Value:       arg7Value,
		Data:        arg8Data,
	})
	if err != nil {
		return nil, false, fmt.Errorf("pack finalization payload: %w", err)
	}
	// trim the first four bytes (function signature)
	finalizationPayload = finalizationPayload[4:]
	return finalizationPayload, true, nil
}

func (l *l2ToL1Bridge) getL1BlockFromRPC(ctx context.Context, txHash common.Hash) (*big.Int, error) {
	type Response struct {
		L1BlockNumber hexutil.Big `json:"l1BlockNumber"`
	}
	response := new(Response)
	err := l.l2Client.CallContext(ctx, response, "eth_getTransactionReceipt", txHash.Hex())
	if err != nil {
		return nil, fmt.Errorf("call eth_getTransactionReceipt with tx hash %s: %w", txHash, err)
	}
	return response.L1BlockNumber.ToInt(), nil
}

func (l *l2ToL1Bridge) getProof(ctx context.Context, l2ToL1Id *big.Int) ([][32]byte, error) {
	l.lggr.Infow("Getting proof for l2ToL1Id", "l2ToL1Id", l2ToL1Id)
	// 1. Get the latest NodeConfirmed event on L1, which indicates the latest node that was confirmed by the rollup.
	latestNodeConfirmed, err := l.getLatestNodeConfirmed(ctx)
	if err != nil {
		return nil, fmt.Errorf("get latest node confirmed: %w", err)
	}
	// 2. Call eth_getBlockByHash on L2 specifying the L2 block hash in the NodeConfirmed event.
	sendCount, err := l.getSendCountForBlock(ctx, latestNodeConfirmed.BlockHash)
	if err != nil {
		return nil, fmt.Errorf("get send count for block: %w", err)
	}
	// 5. Call `constructOutboxProof` on the L2 node interface contract with the `sendCount` as the first argument and `l2ToL1Id` as the second argument.
	outboxProof, err := l.nodeInterface.ConstructOutboxProof(&bind.CallOpts{
		Context: ctx,
	}, sendCount, l2ToL1Id.Uint64())
	if err != nil {
		// if there's an error calling constructOutboxProof it means that the
		// transfer is not yet ready to finalize.
		l.lggr.Infow("construct outbox proof, transfer not ready to finalize",
			"l2ToL1Id", l2ToL1Id,
			"sendCount", sendCount,
			"err", err)
		return nil, nil
	}
	return outboxProof.Proof, nil
}

func (l *l2ToL1Bridge) getSendCountForBlock(ctx context.Context, blockHash [32]byte) (uint64, error) {
	type Response struct {
		SendCount hexutil.Big `json:"sendCount"`
	}
	response := new(Response)
	bhHex := hexutil.Encode(blockHash[:])
	err := l.l2Client.CallContext(ctx, response, "eth_getBlockByHash", bhHex, false)
	if err != nil {
		return 0, fmt.Errorf("call eth_getBlockByHash with blockhash %s: %w", bhHex, err)
	}
	return response.SendCount.ToInt().Uint64(), nil
}

func (l *l2ToL1Bridge) getLatestNodeConfirmed(ctx context.Context) (*arbitrum_rollup_core.ArbRollupCoreNodeConfirmed, error) {
	lg, err := l.l1LogPoller.LatestLogByEventSigWithConfs(
		ctx,
		NodeConfirmedTopic,
		l.rollupCore.Address(),
		evmtypes.Finalized,
	)
	if err != nil {
		return nil, fmt.Errorf("get latest node confirmed: %w, topic: %s, address: %s", err, NodeConfirmedTopic, l.rollupCore.Address())
	}

	parsed, err := l.rollupCore.ParseNodeConfirmed(lg.ToGethLog())
	if err != nil {
		return nil, fmt.Errorf("parse node confirmed log: %w", err)
	}

	return parsed, nil
}

func unpackFinalizationPayload(calldata []byte) (arbitrum_l1_bridge_adapter.ArbitrumL1BridgeAdapterArbitrumFinalizationPayload, error) {
	method, ok := l1AdapterABI.Methods["exposeArbitrumFinalizationPayload"]
	if !ok {
		return arbitrum_l1_bridge_adapter.ArbitrumL1BridgeAdapterArbitrumFinalizationPayload{}, fmt.Errorf("exposeArbitrumFinalizationPayload not found in ArbitrumL1BridgeAdapter ABI")
	}

	ifaces, err := method.Inputs.Unpack(calldata)
	if err != nil {
		return arbitrum_l1_bridge_adapter.ArbitrumL1BridgeAdapterArbitrumFinalizationPayload{}, fmt.Errorf("unpack exposeArbitrumFinalizationPayload: %w", err)
	}

	if len(ifaces) != 1 {
		return arbitrum_l1_bridge_adapter.ArbitrumL1BridgeAdapterArbitrumFinalizationPayload{}, fmt.Errorf("expected 1 argument, got %d", len(ifaces))
	}

	payload := *abi.ConvertType(ifaces[0], new(arbitrum_l1_bridge_adapter.ArbitrumL1BridgeAdapterArbitrumFinalizationPayload)).(*arbitrum_l1_bridge_adapter.ArbitrumL1BridgeAdapterArbitrumFinalizationPayload)

	return payload, nil
}
