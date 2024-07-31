package opstack

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"go.uber.org/multierr"

	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/optimism_standard_bridge"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/abiutils"
	bridgecommon "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/bridge/common"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/liquiditymanager"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/optimism_l1_bridge_adapter"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

type l1ToL2Bridge struct {
	localSelector      models.NetworkSelector
	remoteSelector     models.NetworkSelector
	l1LiquidityManager liquiditymanager.LiquidityManagerInterface
	l2LiquidityManager liquiditymanager.LiquidityManagerInterface
	l1BridgeAdapter    optimism_l1_bridge_adapter.OptimismL1BridgeAdapterInterface
	l1StandardBridge   optimism_standard_bridge.OptimismStandardBridgeInterface
	l2StandardBridge   optimism_standard_bridge.OptimismStandardBridgeInterface
	l1Client           client.Client
	l2Client           client.Client
	l1LogPoller        logpoller.LogPoller
	l2LogPoller        logpoller.LogPoller
	l1FilterName       string
	l2FilterName       string
	l1Token, l2Token   common.Address
	lggr               logger.Logger
}

func NewL1ToL2Bridge(
	ctx context.Context,
	lggr logger.Logger,
	localSelector,
	remoteSelector models.NetworkSelector,
	l1LiquidityManagerAddress,
	l2LiquidityManagerAddress,
	l1StandardBridgeProxyAddress,
	l2StandardBridgeAddress common.Address,
	l1Client,
	l2Client client.Client,
	l1LogPoller,
	l2LogPoller logpoller.LogPoller,
) (*l1ToL2Bridge, error) {
	localChain, ok := chainsel.ChainBySelector(uint64(localSelector))
	if !ok {
		return nil, fmt.Errorf("unknown chain selector for local chain: %d", localSelector)
	}
	remoteChain, ok := chainsel.ChainBySelector(uint64(remoteSelector))
	if !ok {
		return nil, fmt.Errorf("unknown chain selector for remote chain: %d", remoteSelector)
	}

	l1FilterName := bridgecommon.GetBridgeFilterName(
		"OptimismL1ToL2Bridge",
		"L1",
		l1LiquidityManagerAddress,
		localChain.Name,
		remoteChain.Name,
		"",
	)

	err := l1LogPoller.RegisterFilter(ctx, logpoller.Filter{
		Addresses: []common.Address{l1LiquidityManagerAddress}, // emits LiquidityTransferred
		Name:      l1FilterName,
		EventSigs: []common.Hash{
			bridgecommon.LiquidityTransferredTopic,
		},
		Retention: bridgecommon.DurationMonth,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to register L1 log filter: %w", err)
	}

	// TODO: confirm that we're able to use these L1 proxy addresses for listening to emitted events
	l1StandardBridge, err := optimism_standard_bridge.NewOptimismStandardBridge(l1StandardBridgeProxyAddress, l1Client)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate L1StandardBridge at %s: %w", l1StandardBridgeProxyAddress, err)
	}

	l1LiquidityManager, err := liquiditymanager.NewLiquidityManager(l1LiquidityManagerAddress, l1Client)
	if err != nil {
		return nil, fmt.Errorf("instantiate L1 liquidityManager at %s: %w", l1LiquidityManagerAddress, err)
	}

	xchainRebal, err := l1LiquidityManager.GetCrossChainRebalancer(nil, uint64(remoteSelector))
	if err != nil {
		return nil, fmt.Errorf("get cross chain liquidityManager for remote chain %s: %w", remoteChain.Name, err)
	}

	l1BridgeAdapter, err := optimism_l1_bridge_adapter.NewOptimismL1BridgeAdapter(xchainRebal.LocalBridge, l1Client)
	if err != nil {
		return nil, fmt.Errorf("instantiate L1 bridge adapter at %s: %w", xchainRebal.LocalBridge, err)
	}

	l1Token, err := l1LiquidityManager.ILocalToken(nil)
	if err != nil {
		return nil, fmt.Errorf("get local token from L1 LiquidityManager: %w", err)
	}

	l2StandardBridge, err := optimism_standard_bridge.NewOptimismStandardBridge(l2StandardBridgeAddress, l2Client)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate L2StandardBridge at %s: %w", l2StandardBridgeAddress, err)
	}

	l2LiquidityManager, err := liquiditymanager.NewLiquidityManager(l2LiquidityManagerAddress, l2Client)
	if err != nil {
		return nil, fmt.Errorf("instantiate L2 liquidityManager at %s: %w", l2LiquidityManagerAddress, err)
	}

	l2Token, err := l2LiquidityManager.ILocalToken(nil)
	if err != nil {
		return nil, fmt.Errorf("get local token from L2 LiquidityManager: %w", err)
	}

	l2FilterName := bridgecommon.GetBridgeFilterName(
		"OptimismL1ToL2Bridge",
		"L2",
		l2LiquidityManagerAddress,
		localChain.Name,
		remoteChain.Name,
		"",
	)
	err = l2LogPoller.RegisterFilter(ctx, logpoller.Filter{
		Addresses: []common.Address{
			l2StandardBridgeAddress,   // emits ERC20BridgeFinalized
			l2LiquidityManagerAddress, // emits LiquidityTransferred
		},
		Name: l2FilterName,
		EventSigs: []common.Hash{
			ERC20BridgeFinalizedTopic,              // emitted by the L2 StandardBridge
			bridgecommon.LiquidityTransferredTopic, // emitted by the L2 LiquidityManager
		},
		Retention: bridgecommon.DurationMonth,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to register L2 log filter: %w", err)
	}

	return &l1ToL2Bridge{
		localSelector:      localSelector,
		remoteSelector:     remoteSelector,
		l1LiquidityManager: l1LiquidityManager,
		l2LiquidityManager: l2LiquidityManager,
		l1BridgeAdapter:    l1BridgeAdapter,
		l1StandardBridge:   l1StandardBridge,
		l2StandardBridge:   l2StandardBridge,
		l1Client:           l1Client,
		l2Client:           l2Client,
		l1LogPoller:        l1LogPoller,
		l2LogPoller:        l2LogPoller,
		l1FilterName:       l1FilterName,
		l2FilterName:       l2FilterName,
		l1Token:            l1Token,
		l2Token:            l2Token,
		lggr:               lggr,
	}, nil
}

func (l *l1ToL2Bridge) GetTransfers(
	ctx context.Context,
	localToken,
	remoteToken models.Address,
) ([]models.PendingTransfer, error) {
	lggr := l.lggr.With(
		"localToken", localToken,
		"remoteToken", remoteToken,
	)
	lggr.Info("getting transfers from L1 -> L2")

	if l.l1Token.Cmp(common.Address(localToken)) != 0 {
		return nil, fmt.Errorf("local token mismatch: expected %s, got %s", l.l1Token, localToken)
	}
	if l.l2Token.Cmp(common.Address(remoteToken)) != 0 {
		return nil, fmt.Errorf("remote token mismatch: expected %s, got %s", l.l2Token, remoteToken)
	}

	// TODO: heavy query warning
	fromTs := time.Now().Add(-24 * time.Hour) // last day
	sendLogs, erc20BridgeFinalizedLogs, receiveLogs, err := l.getLogs(ctx, fromTs)
	if err != nil {
		return nil, err
	}

	lggr.Infow("got sorted logs",
		"sendLogs", sendLogs,
		"erc20BridgeFinalizedLogs", erc20BridgeFinalizedLogs,
		"receiveLogs", receiveLogs,
	)

	parsedSent, parsedToLP, err := bridgecommon.ParseLiquidityTransferred(l.l1LiquidityManager.ParseLiquidityTransferred, sendLogs)
	if err != nil {
		return nil, fmt.Errorf("parse L1 -> L2 LiquidityTransferred sent logs: %w", err)
	}

	parsedERC20BridgeFinalized, err := l.parseERC20BridgeFinalized(erc20BridgeFinalizedLogs)
	if err != nil {
		return nil, fmt.Errorf("parse ERC20BridgeFinalized logs: %w", err)
	}

	parseReceived, _, err := bridgecommon.ParseLiquidityTransferred(l.l2LiquidityManager.ParseLiquidityTransferred, receiveLogs)
	if err != nil {
		return nil, fmt.Errorf("parse L1 -> L2 LiquidityTransferred received logs: %w", err)
	}

	lggr.Infow("parsed logs",
		"parsedSent", len(parsedSent),
		"parsedERC20BridgeFinalized", len(parsedERC20BridgeFinalized),
		"parseReceived", len(parseReceived),
	)

	notReady, ready, missingSent, err := partitionTransfers(
		localToken,
		l.l1BridgeAdapter.Address(),
		l.l2StandardBridge.Address(),
		parsedSent,
		parsedERC20BridgeFinalized,
		parseReceived)
	if err != nil {
		return nil, fmt.Errorf("partition transfers: %w", err)
	}
	if len(missingSent) > 0 {
		lggr.Warnw("found L2 bridge finalization logs with no corresponding L1 LiquidityTransferred log", "missingSent", missingSent)
	}

	return l.toPendingTransfers(localToken, remoteToken, notReady, ready, parsedToLP)
}

func (l *l1ToL2Bridge) GetBridgePayloadAndFee(
	ctx context.Context,
	transfer models.Transfer,
) ([]byte, *big.Int, error) {
	// TODO: maybe add check if this is a native transfer or ERC20 transfer
	calldata, err := l1standardBridgeABI.Pack(
		// If we're sending WETH, the bridge adapter unwraps it and calls depositETHTo on the native bridge
		DepositETHToFunction,
		transfer.Receiver,   // 'to'
		uint32(0),           // 'l2Gas': hardcoded to 0 in the OptimismL1BridgeAdapter contract
		transfer.BridgeData, // 'data'
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to pack depositETHTo function call: %w", err)
	}

	// Estimate gas needed for the depositETHTo call issued from the L1 Bridge Adapter
	l1StandardBridgeAddress := l.l1StandardBridge.Address()
	gasPrice, err := l.l1Client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get suggested gas price: %w", err)
	}
	gasLimit, err := l.l1Client.EstimateGas(ctx, ethereum.CallMsg{
		From:     l.l1BridgeAdapter.Address(),
		To:       &l1StandardBridgeAddress,
		Data:     calldata,
		GasPrice: gasPrice,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to estimate gas: %w", err)
	}

	// Scale gas limit by recommended 20% buffer to account for gas burned for L2 txn:
	// https://docs.optimism.io/builders/app-developers/bridging/messaging#fees-for-sending-data-between-l1-and-l2
	// TODO: Applying the 1.2x gas limit bump here to the fee won't really have an effect on the actual gas burned by
	//   the OP bridge since the gas units are hardcoded to 1e6 in services/relay/evm/liquidity_manager.go. We should
	//   instead consider a better way to dynamically bump the gas used by the transmitter, or just hardcode an even
	//   higher gas limit in the transmitter. The OP team has confirmed that only gas up to the "market rate" will be
	//   burned, not all gas remaining in the limit.
	gasLimitBigInt := new(big.Int).SetUint64(gasLimit)
	gasLimitWithL2Buffer := new(big.Int).Mul(gasLimitBigInt, big.NewInt(120))
	gasLimitWithL2Buffer = new(big.Int).Div(gasLimitWithL2Buffer, big.NewInt(100))

	finalGasFee := new(big.Int).Mul(gasPrice, gasLimitWithL2Buffer)
	return transfer.BridgeData, finalGasFee, nil
}

func (l *l1ToL2Bridge) QuorumizedBridgePayload(payloads [][]byte, f int) ([]byte, error) {
	if len(payloads) <= f {
		return nil, fmt.Errorf("not enough payloads to reach quorum, need at least f+1=%d, got len(payloads)=%d", f+1, len(payloads))
	}

	var transferNonces []*big.Int
	for _, payload := range payloads {
		decodedTransferNonce, err := abiutils.UnpackUint256(payload)
		if err != nil {
			return nil, fmt.Errorf("unpack transfer nonce from bridge payload: %w", err)
		}
		transferNonces = append(transferNonces, decodedTransferNonce)
	}
	if len(transferNonces) == 0 {
		return nil, errors.New("no transfer nonces found in bridge payloads")
	}

	// TODO: reconfirm that all nonces should be the same across all payloads in a given round for OP
	for _, nonce := range transferNonces {
		if nonce.Cmp(transferNonces[0]) != 0 {
			return nil, fmt.Errorf("nonces in payloads do not match: %v", transferNonces)
		}
	}

	return payloads[0], nil
}

func (l *l1ToL2Bridge) getLogs(ctx context.Context, fromTs time.Time) (sendLogs []logpoller.Log, erc20BridgeFinalizedLogs []logpoller.Log, receiveLogs []logpoller.Log, err error) {
	// LiquidityTransferred events emitted by the L1 LiquidityManager. Represents transfers that have been initiated
	// from L1 to L2.
	sendLogs, err = l.l1LogPoller.IndexedLogsCreatedAfter(
		ctx,
		bridgecommon.LiquidityTransferredTopic,
		l.l1LiquidityManager.Address(),
		bridgecommon.LiquidityTransferredToChainSelectorTopicIndex,
		[]common.Hash{
			bridgecommon.NetworkSelectorToHash(l.remoteSelector),
		},
		fromTs,
		1,
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, nil, nil, fmt.Errorf("get LiquidityTransferred events from L1 LiquidityManager: %w", err)
	}

	// ERC20BridgeFinalized events emitted by Optimism's L2StandardBridge. Represents L1 to L2 transfers that have been
	// finalized on L2, but potentially not yet "received" by the L2 LiquidityManager.
	erc20BridgeFinalizedLogs, err = l.l2LogPoller.IndexedLogsCreatedAfter(
		ctx,
		ERC20BridgeFinalizedTopic,
		l.l2StandardBridge.Address(),
		// We register the filter on the "from" address in OP whereas Arb registers it on the "to" address. In OP only
		// the "to" address is indexed for this event. To be safe, we check the "to" address below in the partitioning
		// step anyway.
		ERC20BridgeFinalizedFromAddressTopicIndex,
		[]common.Hash{
			common.HexToHash(l.l1BridgeAdapter.Address().Hex()),
		},
		fromTs,
		evmtypes.Finalized,
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, nil, nil, fmt.Errorf("get DepositFinalized events from L2 gateway: %w", err)
	}

	// LiquidityTransferred events emitted by the L2 LiquidityManager. Represents transfers that have been received on L2.
	receiveLogs, err = l.l2LogPoller.IndexedLogsCreatedAfter(
		ctx,
		bridgecommon.LiquidityTransferredTopic,
		l.l2LiquidityManager.Address(),
		bridgecommon.LiquidityTransferredFromChainSelectorTopicIndex,
		[]common.Hash{
			bridgecommon.NetworkSelectorToHash(l.localSelector),
		},
		fromTs,
		1,
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, nil, nil, fmt.Errorf("get LiquidityTransferred events from L2 LiquidityManager: %w", err)
	}

	return sendLogs, erc20BridgeFinalizedLogs, receiveLogs, nil
}

func (l *l1ToL2Bridge) parseERC20BridgeFinalized(erc20BridgeFinalizedLogs []logpoller.Log) ([]*optimism_standard_bridge.OptimismStandardBridgeERC20BridgeFinalized, error) {
	finalized := make([]*optimism_standard_bridge.OptimismStandardBridgeERC20BridgeFinalized, len(erc20BridgeFinalizedLogs))
	for i, lg := range erc20BridgeFinalizedLogs {
		parsed, err := l.l2StandardBridge.ParseERC20BridgeFinalized(lg.ToGethLog())
		if err != nil {
			return nil, fmt.Errorf("parse ERC20BridgeFinalized log: %w", err)
		}
		finalized[i] = parsed
	}
	return finalized, nil
}

/**
 * This function partitions the transfer events into four categories:
 *   - notReady: 	sent (L1 LiquidityTransferred), but not finalized (L2 ERC20BridgeFinalized)
 *   - ready:	 	sent (L1 LiquidityTransferred), and finalized (L2 ERC20BridgeFinalized), but not received (L2 LiquidityTransferred)
 *   - done:	 	sent (L1 LiquidityTransferred), finalized (L2 ERC20BridgeFinalized), and received (L2 LiquidityTransferred)
 *   - missingSent: finalized (L2 ERC20BridgeFinalized), but no corresponding sent (L1 LiquidityTransferred)
 *
 * Since we only care about the pending transfers, this function only returns 'notReady', 'ready', and 'missingSent'.
 * The matching logic is performed based on the fact that a nonce is piped through all events:
 *   sent_LiquidityTransferred.bridgeReturnData == ERC20BridgeFinalized.extraData == received_LiquidityTransferred.bridgeSpecificData
 */
func partitionTransfers(
	localToken models.Address,
	l1BridgeAdapterAddress common.Address,
	l2LiquidityManagerAddress common.Address,
	sentLogs []*liquiditymanager.LiquidityManagerLiquidityTransferred,
	erc20BridgeFinalizedLogs []*optimism_standard_bridge.OptimismStandardBridgeERC20BridgeFinalized,
	receivedLogs []*liquiditymanager.LiquidityManagerLiquidityTransferred,
) (
	notReady,
	ready []*liquiditymanager.LiquidityManagerLiquidityTransferred,
	missingSent []*optimism_standard_bridge.OptimismStandardBridgeERC20BridgeFinalized,
	err error,
) {
	transferNonceToSentLogMap := make(map[string]*liquiditymanager.LiquidityManagerLiquidityTransferred)
	foundMatchingFinalizedLogMap := make(map[string]bool)
	for _, sentLog := range sentLogs {
		if sentLog.To != l2LiquidityManagerAddress {
			continue
		}
		var transferNonce *big.Int
		transferNonce, err = abiutils.UnpackUint256(sentLog.BridgeReturnData)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("unpack transfer nonce from L1 LiquidityTransferred log (%s): %w, data: %s",
				sentLog.Raw.TxHash, err, hexutil.Encode(sentLog.BridgeReturnData))
		}
		transferNonceToSentLogMap[transferNonce.String()] = sentLog
		foundMatchingFinalizedLogMap[transferNonce.String()] = false
	}

	// For each finalized log, check if it has a corresponding sent log. If there is no corresponding sent log, add it
	// to 'missingSent'
	for _, finalized := range erc20BridgeFinalizedLogs {
		if finalized.RemoteToken != common.Address(localToken) {
			continue
		}
		if finalized.From != l1BridgeAdapterAddress {
			continue
		}
		if finalized.To != l2LiquidityManagerAddress {
			continue
		}
		var transferNonce *big.Int
		transferNonce, err = abiutils.UnpackUint256(finalized.ExtraData)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("unpack transfer nonce from L2 ERC20BridgeFinalized log (%s): %w, data: %s",
				finalized.Raw.TxHash, err, hexutil.Encode(finalized.ExtraData))
		}
		if sentLog, exists := transferNonceToSentLogMap[transferNonce.String()]; exists {
			// If a corresponding sentLog exists for this finalized log, add it to ready
			ready = append(ready, sentLog)
			foundMatchingFinalizedLogMap[transferNonce.String()] = true
		} else if !exists {
			// Else, if a corresponding sentLog does not exist, add it to missingSent
			missingSent = append(missingSent, finalized)
		}
	}

	// Any entries in foundMatchingFinalizedLogMap that are 'false' were not found to have a matching finalized log
	// and are therefore not ready to be "received" yet.
	for transferID, found := range foundMatchingFinalizedLogMap {
		if !found {
			if sentLog, exists := transferNonceToSentLogMap[transferID]; exists {
				notReady = append(notReady, sentLog)
			}
		}
	}

	// Filter out from 'ready' any logs from transfers that have already been received by the L2 LM
	ready, err = filterExecuted(ready, receivedLogs)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("filter executed: %w", err)
	}
	return
}

func (l *l1ToL2Bridge) toPendingTransfers(
	localToken,
	remoteToken models.Address,
	notReady,
	ready []*liquiditymanager.LiquidityManagerLiquidityTransferred,
	parsedToLP map[bridgecommon.LogKey]logpoller.Log,
) ([]models.PendingTransfer, error) {
	var transfers []models.PendingTransfer
	for _, transfer := range notReady {
		transfers = append(transfers, models.PendingTransfer{
			Transfer: models.Transfer{
				From:               l.localSelector,
				To:                 l.remoteSelector,
				Sender:             models.Address(l.l1LiquidityManager.Address()),
				Receiver:           models.Address(l.l2LiquidityManager.Address()),
				LocalTokenAddress:  localToken,
				RemoteTokenAddress: remoteToken,
				Amount:             ubig.New(transfer.Amount),
				Date: parsedToLP[bridgecommon.LogKey{
					TxHash:   transfer.Raw.TxHash,
					LogIndex: int64(transfer.Raw.Index),
				}].BlockTimestamp,
				BridgeData:      transfer.BridgeReturnData, // unique nonce from the OP Bridge Adapter
				Stage:           bridgecommon.StageRebalanceConfirmed,
				NativeBridgeFee: ubig.NewI(0),
			},
			Status: models.TransferStatusNotReady,
			ID:     fmt.Sprintf("%s-%d", transfer.Raw.TxHash.Hex(), transfer.Raw.Index),
		})
	}
	for _, transfer := range ready {
		transfers = append(transfers, models.PendingTransfer{
			Transfer: models.Transfer{
				From:               l.localSelector,
				To:                 l.remoteSelector,
				Sender:             models.Address(l.l1LiquidityManager.Address()),
				Receiver:           models.Address(l.l2LiquidityManager.Address()),
				LocalTokenAddress:  localToken,
				RemoteTokenAddress: remoteToken,
				Amount:             ubig.New(transfer.Amount),
				Date: parsedToLP[bridgecommon.LogKey{
					TxHash:   transfer.Raw.TxHash,
					LogIndex: int64(transfer.Raw.Index),
				}].BlockTimestamp,
				BridgeData:      transfer.BridgeReturnData, // unique nonce from the OP Bridge Adapter
				Stage:           bridgecommon.StageFinalizeReady,
				NativeBridgeFee: ubig.NewI(0),
			},
			Status: models.TransferStatusReady, // ready == finalized for L1 -> L2 transfers due to auto-finalization by the native bridge
			ID:     fmt.Sprintf("%s-%d", transfer.Raw.TxHash.Hex(), transfer.Raw.Index),
		})
	}
	return transfers, nil
}

func (l *l1ToL2Bridge) Close(ctx context.Context) error {
	return multierr.Combine(
		l.l1LogPoller.UnregisterFilter(ctx, l.l1FilterName),
		l.l2LogPoller.UnregisterFilter(ctx, l.l2FilterName),
	)
}
