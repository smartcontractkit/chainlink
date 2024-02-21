package arb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"go.uber.org/multierr"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/rebalancer/generated/abstract_arbitrum_token_gateway"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/rebalancer/generated/arbitrum_gateway_router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/rebalancer/generated/arbitrum_inbox"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/rebalancer/generated/arbitrum_l1_bridge_adapter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/rebalancer/generated/arbitrum_token_gateway"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/rebalancer/generated/l2_arbitrum_gateway"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/rebalancer/generated/rebalancer"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

type l1ToL2Bridge struct {
	localSelector       models.NetworkSelector
	remoteSelector      models.NetworkSelector
	l1Rebalancer        rebalancer.RebalancerInterface
	l2RebalancerAddress common.Address
	l1BridgeAdapter     arbitrum_l1_bridge_adapter.ArbitrumL1BridgeAdapterInterface
	l1GatewayRouter     arbitrum_gateway_router.ArbitrumGatewayRouterInterface
	l1Inbox             arbitrum_inbox.ArbitrumInboxInterface
	l2Gateway           l2_arbitrum_gateway.L2ArbitrumGatewayInterface
	l1Client            client.Client
	l2Client            client.Client
	l1LogPoller         logpoller.LogPoller
	l2LogPoller         logpoller.LogPoller
	l1FilterName        string
	l2FilterName        string
	lggr                logger.Logger
}

func NewL1ToL2Bridge(
	lggr logger.Logger,
	localSelector,
	remoteSelector models.NetworkSelector,
	l1RebalancerAddress,
	l2RebalancerAddress,
	l1GatewayRouterAddress,
	l1InboxAddress common.Address,
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

	l1GatewayRouter, err := arbitrum_gateway_router.NewArbitrumGatewayRouter(l1GatewayRouterAddress, l1Client)
	if err != nil {
		return nil, fmt.Errorf("instantiate L1 gateway router at %s: %w", l1GatewayRouterAddress, err)
	}

	l1Inbox, err := arbitrum_inbox.NewArbitrumInbox(l1InboxAddress, l1Client)
	if err != nil {
		return nil, fmt.Errorf("instantiate L1 inbox at %s: %w", l1InboxAddress, err)
	}

	l1FilterName := fmt.Sprintf("ArbitrumL2ToL1Bridge-L1-Rebalancer:%s-Local:%s-Remote:%s",
		l1RebalancerAddress.String(), localChain.Name, remoteChain.Name)
	err = l1LogPoller.RegisterFilter(logpoller.Filter{
		Addresses: []common.Address{l1RebalancerAddress},
		Name:      l1FilterName,
		EventSigs: []common.Hash{
			LiquidityTransferredTopic,
		},
		Retention: DurationMonth,
	})
	if err != nil {
		return nil, fmt.Errorf("register L1 log filter: %w", err)
	}

	// figure out which gateway to watch for the token on L2
	l1Rebalancer, err := rebalancer.NewRebalancer(l1RebalancerAddress, l1Client)
	if err != nil {
		return nil, fmt.Errorf("instantiate rebalancer at %s: %w", l1RebalancerAddress, err)
	}

	xchainRebal, err := l1Rebalancer.GetCrossChainRebalancer(nil, uint64(remoteSelector))
	if err != nil {
		return nil, fmt.Errorf("get cross chain rebalancer for remote chain %s: %w", remoteChain.Name, err)
	}

	l1BridgeAdapter, err := arbitrum_l1_bridge_adapter.NewArbitrumL1BridgeAdapter(xchainRebal.LocalBridge, l1Client)
	if err != nil {
		return nil, fmt.Errorf("instantiate L1 bridge adapter at %s: %w", xchainRebal.LocalBridge, err)
	}

	l1Token, err := l1Rebalancer.ILocalToken(nil)
	if err != nil {
		return nil, fmt.Errorf("get local token from rebalancer: %w", err)
	}

	// get the gateway on L1 and then it's counterpart gateway on L2
	// that's the one we need to watch
	l1TokenGateway, err := l1GatewayRouter.GetGateway(nil, l1Token)
	if err != nil {
		return nil, fmt.Errorf("get gateway for token %s: %w, gateway router: %s", l1Token, err, l1GatewayRouterAddress)
	}

	abstractGateway, err := abstract_arbitrum_token_gateway.NewAbstractArbitrumTokenGateway(l1TokenGateway, l1Client)
	if err != nil {
		return nil, fmt.Errorf("instantiate abstract gateway at %s: %w", l1TokenGateway, err)
	}

	l2Gateway, err := abstractGateway.CounterpartGateway(nil)
	if err != nil {
		return nil, fmt.Errorf("get counterpart gateway for gateway %s: %w", l1TokenGateway, err)
	}

	l2FilterName := fmt.Sprintf("ArbitrumL2ToL1Bridge-L2-L2Gateway:%s-Rebalancer:%s-Local:%s-Remote:%s",
		l2Gateway.Hex(), l2RebalancerAddress.Hex(), localChain.Name, remoteChain.Name)
	err = l2LogPoller.RegisterFilter(logpoller.Filter{
		Addresses: []common.Address{
			l2Gateway,           // emits DepositFinalized
			l2RebalancerAddress, // emits LiquidityTransferred
		},
		Name: l2FilterName,
		EventSigs: []common.Hash{
			DepositFinalizedTopic,     // emitted by the gateways
			LiquidityTransferredTopic, // emitted by the rebalancers
		},
		Retention: DurationMonth,
	})
	if err != nil {
		return nil, fmt.Errorf("register L2 log filter: %w", err)
	}

	l2GatewayWrapper, err := l2_arbitrum_gateway.NewL2ArbitrumGateway(l2Gateway, l2Client)
	if err != nil {
		return nil, fmt.Errorf("instantiate l2 arbitrum gateway at %s: %w", l2Gateway, err)
	}

	lggr = lggr.Named("ArbitrumL1ToL2Bridge").With(
		"localSelector", localSelector,
		"remoteSelector", remoteSelector,
		"localChainID", localChain.EvmChainID,
		"remoteChainID", remoteChain.EvmChainID,
		"l1Rebalancer", l1Rebalancer.Address(),
		"l2Rebalancer", l2RebalancerAddress,
		"l1BridgeAdapter", l1BridgeAdapter.Address(),
		"l1GatewayRouter", l1GatewayRouter.Address(),
		"l1Inbox", l1Inbox.Address(),
		"l2Gateway", l2Gateway,
	)
	lggr.Infow("successfully initialized arbitrum L1 -> L2 bridge")

	return &l1ToL2Bridge{
		localSelector:       localSelector,
		remoteSelector:      remoteSelector,
		l1Rebalancer:        l1Rebalancer,
		l2RebalancerAddress: l2RebalancerAddress,
		l1BridgeAdapter:     l1BridgeAdapter,
		l1GatewayRouter:     l1GatewayRouter,
		l1Inbox:             l1Inbox,
		l2Gateway:           l2GatewayWrapper,
		l1Client:            l1Client,
		l2Client:            l2Client,
		l1LogPoller:         l1LogPoller,
		l2LogPoller:         l2LogPoller,
		l1FilterName:        l1FilterName,
		l2FilterName:        l2FilterName,
		lggr:                lggr,
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
	lggr.Info("getting transfers")
	// TODO: check that l1Rebalancer token matches localToken
	// TODO: check that l2Rebalancer token matches remoteToken
	fromTs := time.Now().Add(-24 * time.Hour) // last day
	sendLogs, err := l.l1LogPoller.IndexedLogsCreatedAfter(
		LiquidityTransferredTopic,
		l.l1Rebalancer.Address(),
		3, // topic index 3: toChainSelector in event
		[]common.Hash{
			toHash(l.remoteSelector),
		},
		fromTs,
		logpoller.Finalized,
		pg.WithParentCtx(ctx),
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("get LiquidityTransferred events from L1 rebalancer: %w", err)
	}

	depositFinalizedLogs, err := l.l2LogPoller.IndexedLogsCreatedAfter(
		DepositFinalizedTopic,
		l.l2Gateway.Address(),
		3, // topic index 3: to address of deposit on L2
		[]common.Hash{
			common.HexToHash(l.l2RebalancerAddress.Hex()),
		},
		fromTs,
		logpoller.Finalized,
		pg.WithParentCtx(ctx),
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("get DepositFinalized events from L2 gateway: %w", err)
	}

	receiveLogs, err := l.l2LogPoller.IndexedLogsCreatedAfter(
		LiquidityTransferredTopic,
		l.l2RebalancerAddress,
		2, // topic index 2: fromChainSelector
		[]common.Hash{
			toHash(l.localSelector),
		},
		fromTs,
		logpoller.Finalized,
		pg.WithParentCtx(ctx),
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("get LiquidityTransferred events from L2 rebalancer: %w", err)
	}

	slices.SortFunc(sendLogs, func(a, b logpoller.Log) int {
		return a.BlockTimestamp.Compare(b.BlockTimestamp)
	})
	slices.SortFunc(depositFinalizedLogs, func(a, b logpoller.Log) int {
		return a.BlockTimestamp.Compare(b.BlockTimestamp)
	})
	slices.SortFunc(receiveLogs, func(a, b logpoller.Log) int {
		return a.BlockTimestamp.Compare(b.BlockTimestamp)
	})

	lggr.Infow("got logs",
		"sendLogs", len(sendLogs),
		"depositFinalizedLogs", len(depositFinalizedLogs),
		"receiveLogs", len(receiveLogs),
	)

	parsedSent, parsedToLP, err := parseLiquidityTransferred(l.l1Rebalancer.ParseLiquidityTransferred, sendLogs)
	if err != nil {
		return nil, fmt.Errorf("parse L1 -> L2 transfers: %w", err)
	}

	parsedDepositFinalized, err := l.parseDepositFinalized(depositFinalizedLogs)
	if err != nil {
		return nil, fmt.Errorf("parse DepositFinalized logs: %w", err)
	}

	parsedReceived, _, err := parseLiquidityTransferred(l.l1Rebalancer.ParseLiquidityTransferred, receiveLogs)
	if err != nil {
		return nil, fmt.Errorf("parse LiquidityTransferred logs: %w", err)
	}

	lggr.Infow("parsed logs",
		"parsedSent", len(parsedSent),
		"parsedDepositFinalized", len(parsedDepositFinalized),
		"parsedReceived", len(parsedReceived),
	)

	// Unfortunately its not easy to match DepositFinalized events with LiquidityTransferred events.
	// Reason being that arbitrum does not emit any identifying information as part of the DepositFinalized
	// event, such as the l1 to l2 tx id. This is only available as part of the calldata for when the L2 calls
	// submitRetryable on the ArbRetryableTx precompile.
	// e.g https://sepolia.arbiscan.io/tx/0xce0d0d7e74f184fa8cb264b6d9aab5ced159faf3d0d9ae54b67fd40ba9d965a7
	// therefore we're kind of relegated here to doing a simple count check - filter out all of the
	// LiquidityTransferred logs destined for the rebalancer on L2 and all the DepositFinalized logs that
	// pay out to the rebalancer on L2.
	// We can _probably_ assume that the earlier LiquidityTransferred logs on L1
	// are more likely to be finalizedNotExecuted than later ones.
	notReady, ready, readyData, err := l.partitionTransfers(localToken, parsedSent, parsedDepositFinalized, parsedReceived)
	if err != nil {
		return nil, fmt.Errorf("partition logs into not-ready and ready states: %w", err)
	}

	return l.toPendingTransfers(localToken, remoteToken, notReady, ready, readyData, parsedToLP)
}

func (l *l1ToL2Bridge) toPendingTransfers(
	localToken, remoteToken models.Address,
	notReady,
	ready []*rebalancer.RebalancerLiquidityTransferred,
	readyData [][]byte,
	parsedToLP map[logKey]logpoller.Log,
) ([]models.PendingTransfer, error) {
	if len(ready) != len(readyData) {
		return nil, fmt.Errorf("length of ready and readyData should be the same: len(ready) = %d, len(readyData) = %d",
			len(ready), len(readyData))
	}
	var transfers []models.PendingTransfer
	for _, transfer := range notReady {
		transfers = append(transfers, models.PendingTransfer{
			Transfer: models.Transfer{
				From:               l.localSelector,
				To:                 l.remoteSelector,
				Sender:             models.Address(l.l1Rebalancer.Address()),
				Receiver:           models.Address(l.l2RebalancerAddress),
				LocalTokenAddress:  localToken,
				RemoteTokenAddress: remoteToken,
				Amount:             ubig.New(transfer.Amount),
				Date: parsedToLP[logKey{
					txHash:   transfer.Raw.TxHash,
					logIndex: int64(transfer.Raw.Index),
				}].BlockTimestamp,
				BridgeData: []byte{}, // no finalization data, not ready
			},
			Status: models.TransferStatusNotReady,
			ID:     fmt.Sprintf("%s-%d", transfer.Raw.TxHash.Hex(), transfer.Raw.Index),
		})
	}
	for i, transfer := range ready {
		transfers = append(transfers, models.PendingTransfer{
			Transfer: models.Transfer{
				From:               l.localSelector,
				To:                 l.remoteSelector,
				Sender:             models.Address(l.l1Rebalancer.Address()),
				Receiver:           models.Address(l.l1Rebalancer.Address()),
				LocalTokenAddress:  localToken,
				RemoteTokenAddress: remoteToken,
				Amount:             ubig.New(transfer.Amount),
				Date: parsedToLP[logKey{
					txHash:   transfer.Raw.TxHash,
					logIndex: int64(transfer.Raw.Index),
				}].BlockTimestamp,
				BridgeData: readyData[i], // finalization data since its ready
			},
			Status: models.TransferStatusReady, // ready == finalized for L1 -> L2 transfers due to auto-finalization by the native bridge
			ID:     fmt.Sprintf("%s-%d", transfer.Raw.TxHash.Hex(), transfer.Raw.Index),
		})
	}
	return transfers, nil
}

// precondition: the input logs are already sorted in time-ascending order
func (l *l1ToL2Bridge) partitionTransfers(
	localToken models.Address,
	sentLogs []*rebalancer.RebalancerLiquidityTransferred,
	depositFinalizedLogs []*l2_arbitrum_gateway.L2ArbitrumGatewayDepositFinalized,
	receivedLogs []*rebalancer.RebalancerLiquidityTransferred,
) (
	notReady,
	ready []*rebalancer.RebalancerLiquidityTransferred,
	readyData [][]byte,
	err error,
) {
	effectiveDepositFinalized := l.getEffectiveEvents(localToken, depositFinalizedLogs)
	// determine ready and not ready first
	if len(sentLogs) > len(effectiveDepositFinalized) {
		// more sent than have been finalized
		for i := len(sentLogs) - len(effectiveDepositFinalized) + 1; i < len(sentLogs); i++ {
			notReady = append(notReady, sentLogs[i])
		}
		for i := 0; i < (len(sentLogs) - len(effectiveDepositFinalized)); i++ {
			ready = append(ready, sentLogs[i])
		}
	} else if len(sentLogs) < len(effectiveDepositFinalized) {
		// more finalized than have been sent - should be impossible
		return nil, nil, nil, fmt.Errorf("got more finalized logs than sent - should be impossible: len(sent) = %d, len(finalized) = %d",
			len(sentLogs), len(effectiveDepositFinalized))
	} else {
		ready = sentLogs
	}
	// figure out if any of the ready have been executed
	ready, err = l.filterExecuted(ready, receivedLogs)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("filter executed transfers: %w", err)
	}
	// get the readyData
	// this is just going to be the L1 to L2 tx id that is emitted in the L1 LiquidityTransferred.bridgeReturnData field.
	for _, r := range ready {
		readyData = append(readyData, r.BridgeReturnData)
	}
	return
}

func (l *l1ToL2Bridge) filterExecuted(
	readyCandidates []*rebalancer.RebalancerLiquidityTransferred,
	liquidityTransferredLogs []*rebalancer.RebalancerLiquidityTransferred,
) (
	ready []*rebalancer.RebalancerLiquidityTransferred,
	err error,
) {
	for _, readyCandidate := range readyCandidates {
		exists, err := l.matchingExecutionExists(readyCandidate, liquidityTransferredLogs)
		if err != nil {
			return nil, fmt.Errorf("error checking if ready candidate has been executed: %w", err)
		}
		if !exists {
			ready = append(ready, readyCandidate)
		}
	}
	return
}

// TODO: might be able to optimize this
// map[l2ToL1TxId]bool and check if the l2ToL1TxId exists in the map
func (l *l1ToL2Bridge) matchingExecutionExists(
	readyCandidate *rebalancer.RebalancerLiquidityTransferred,
	liquidityTransferredLogs []*rebalancer.RebalancerLiquidityTransferred,
) (bool, error) {
	for _, ltLog := range liquidityTransferredLogs {
		// decode the bridge specific data, which should be the l1 -> l2 tx id
		recvL1ToL2TxId, err := unpackUint256(ltLog.BridgeSpecificData)
		if err != nil {
			return false, fmt.Errorf("unpack bridge specific data from LiquidityTransferred log: %w, data: %s",
				err, hexutil.Encode(ltLog.BridgeSpecificData))
		}
		sendL1ToL2TxId, err := unpackUint256(readyCandidate.BridgeSpecificData)
		if err != nil {
			return false, fmt.Errorf("unpack outbound transfer result from LiquidityTransferred log: %w, data: %s",
				err, hexutil.Encode(readyCandidate.BridgeSpecificData))
		}
		if sendL1ToL2TxId.Cmp(recvL1ToL2TxId) == 0 {
			return true, nil
		}
	}
	return false, nil
}

func (l *l1ToL2Bridge) getEffectiveEvents(
	localToken models.Address,
	depositFinalizedLogs []*l2_arbitrum_gateway.L2ArbitrumGatewayDepositFinalized,
) (
	effectiveDepositFinalized []*l2_arbitrum_gateway.L2ArbitrumGatewayDepositFinalized,
) {
	// filter out DepositFinalized logs not coming from the l1 bridge adapter
	// and not matching the localToken provided.
	// in theory anyone can bridge any token to the rebalancer on L2 from L1
	// TODO: ideally this would be done in the log poller query but no such query exists
	// at the moment.
	// TODO: should we care about L1 -> L2 bridges not done by the bridge adapter?
	// in theory those are funds that can be injected into the pools.
	for _, depFinalized := range depositFinalizedLogs {
		if depFinalized.From == l.l1BridgeAdapter.Address() &&
			depFinalized.L1Token == common.Address(localToken) {
			effectiveDepositFinalized = append(effectiveDepositFinalized, depFinalized)
		}
	}
	return
}

func (l *l1ToL2Bridge) parseDepositFinalized(lgs []logpoller.Log) ([]*l2_arbitrum_gateway.L2ArbitrumGatewayDepositFinalized, error) {
	finalized := make([]*l2_arbitrum_gateway.L2ArbitrumGatewayDepositFinalized, len(lgs))
	for i, lg := range lgs {
		parsed, err := l.l2Gateway.ParseDepositFinalized(lg.ToGethLog())
		if err != nil {
			// should never happen
			return nil, fmt.Errorf("parse DepositFinalized log: %w", err)
		}
		finalized[i] = parsed
	}
	return finalized, nil
}

func (l *l1ToL2Bridge) QuorumizedBridgePayload(payloads [][]byte, f int) ([]byte, error) {
	if len(payloads) <= f {
		return nil, fmt.Errorf("not enough payloads to quorumize, need at least f+1: len(payloads) = %d, f = %d", len(payloads), f)
	}
	// TODO: decode and take top n-f index after decoding and sorting asc gasLimit/maxSubmissionCost/maxFeePerGas
	var (
		gasLimits          []*big.Int
		maxSubmissionCosts []*big.Int
		maxFeePerGases     []*big.Int
	)
	for _, payload := range payloads {
		params, err := UnpackL1ToL2SendBridgePayload(payload)
		if err != nil {
			return nil, fmt.Errorf("decode bridge payload: %w", err)
		}
		gasLimits = append(gasLimits, params.GasLimit)
		maxSubmissionCosts = append(maxSubmissionCosts, params.MaxSubmissionCost)
		maxFeePerGases = append(maxFeePerGases, params.MaxFeePerGas)
	}
	slices.SortFunc(gasLimits, func(i, j *big.Int) int {
		return i.Cmp(j)
	})
	slices.SortFunc(maxSubmissionCosts, func(i, j *big.Int) int {
		return i.Cmp(j)
	})
	slices.SortFunc(maxFeePerGases, func(i, j *big.Int) int {
		return i.Cmp(j)
	})
	// return f-th highest gasLimit/maxSubmissionCost/maxFeePerGas
	return PackL1ToL2SendBridgePayload(
		gasLimits[len(gasLimits)-f-1],
		maxSubmissionCosts[len(maxSubmissionCosts)-f-1],
		maxFeePerGases[len(maxFeePerGases)-f-1],
	)
}

// GetBridgePayloadAndFee implements bridge.Bridge
// For Arbitrum L1 -> L2 transfers, the bridge specific payload is a tuple of 3 numbers:
// 1. gasLimit
// 2. maxSubmissionCost
// 3. maxFeePerGas
func (l *l1ToL2Bridge) GetBridgePayloadAndFee(
	ctx context.Context,
	transfer models.Transfer,
) ([]byte, *big.Int, error) {
	// TODO: can this information be cached in the struct?
	// we already do this stuff in New() so it's unclear if we need to do this everytime.
	l1Gateway, err := l.l1GatewayRouter.GetGateway(&bind.CallOpts{
		Context: ctx,
	}, common.Address(transfer.LocalTokenAddress))
	if err != nil {
		return nil, nil, fmt.Errorf("get L1 gateway for local token %s: %w",
			transfer.LocalTokenAddress, err)
	}

	l1TokenGateway, err := arbitrum_token_gateway.NewArbitrumTokenGateway(l1Gateway, l.l1Client)
	if err != nil {
		return nil, nil, fmt.Errorf("instantiate L1 token gateway at %s: %w",
			l1Gateway, err)
	}

	// get the counterpart gateway on L2 from the L1 gateway
	// unfortunately we need to instantiate a new wrapper because the counterpartGateway field,
	// although it is public, is not accessible via a getter function on the token gateway interface
	abstractGateway, err := abstract_arbitrum_token_gateway.NewAbstractArbitrumTokenGateway(l1Gateway, l.l1Client)
	if err != nil {
		return nil, nil, fmt.Errorf("instantiate abstract gateway at %s: %w",
			l1Gateway, err)
	}

	l2Gateway, err := abstractGateway.CounterpartGateway(&bind.CallOpts{
		Context: ctx,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("get counterpart gateway for L1 gateway %s: %w",
			l1Gateway, err)
	}

	retryableData := RetryableData{
		From:                l1Gateway,
		To:                  l2Gateway,
		ExcessFeeRefundAddr: common.Address(transfer.Receiver),
		CallValueRefundAddr: common.Address(transfer.Sender),
		// this is the amount - see the arbitrum SDK.
		// https://github.com/OffchainLabs/arbitrum-sdk/blob/4c0d43abd5fcc5d219b20bc55e9d0ee152c01309/src/lib/assetBridger/ethBridger.ts#L318
		L2CallValue: transfer.Amount.ToInt(),
		// 3 seems to work, but not sure if it's the best value
		// you definitely need a non-nil deposit for the NodeInterface call to succeed
		Deposit: big.NewInt(3),
		// MaxSubmissionCost: , // To be filled in
		// GasLimit: , // To be filled in
		// MaxFeePerGas: , // To be filled in
		// Data: , // To be filled in
	}

	// determine the finalizeInboundTransfer calldata
	finalizeInboundTransferCalldata, err := l1TokenGateway.GetOutboundCalldata(
		nil,
		common.Address(transfer.LocalTokenAddress), // L1 token address
		l.l1BridgeAdapter.Address(),                // L1 sender address
		common.Address(transfer.Receiver),          // L2 recipient address
		transfer.Amount.ToInt(),                    // token amount
		[]byte{},                                   // extra data (unused here)
	)
	if err != nil {
		return nil, nil, fmt.Errorf("get finalizeInboundTransfer calldata: %w", err)
	}
	retryableData.Data = finalizeInboundTransferCalldata

	l.lggr.Infow("Constructed RetryableData",
		"from", retryableData.From,
		"to", retryableData.To,
		"excessFeeRefundAddr", retryableData.ExcessFeeRefundAddr,
		"callValueRefundAddr", retryableData.CallValueRefundAddr,
		"l2CallValue", retryableData.L2CallValue,
		"deposit", retryableData.Deposit,
		"data", hexutil.Encode(retryableData.Data))

	l1BaseFee, err := l.l1Client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("get L1 base fee: %w", err)
	}

	return l.estimateAll(ctx, retryableData, l1BaseFee)
}

func (l *l1ToL2Bridge) estimateAll(
	ctx context.Context,
	retryableData RetryableData,
	l1BaseFee *big.Int,
) ([]byte, *big.Int, error) {
	l2MaxFeePerGas, err := l.estimateMaxFeePerGasOnL2(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("estimate max fee per gas on L2: %w", err)
	}

	maxSubmissionFee, err := l.estimateMaxSubmissionFee(ctx, l1BaseFee, len(retryableData.Data))
	if err != nil {
		return nil, nil, fmt.Errorf("estimate max submission fee: %w", err)
	}

	gasLimit, err := l.estimateRetryableGasLimit(ctx, retryableData)
	if err != nil {
		return nil, nil, fmt.Errorf("estimate retryable gas limit: %w", err)
	}

	deposit := new(big.Int).Mul(gasLimit, l2MaxFeePerGas)
	deposit = deposit.Add(deposit, maxSubmissionFee)

	l.lggr.Infow("Estimated L1 -> L2 fees",
		"gasLimit", gasLimit,
		"maxSubmissionFee", maxSubmissionFee,
		"l2MaxFeePerGas", l2MaxFeePerGas,
		"deposit", deposit)

	bridgeCalldata, err := PackL1ToL2SendBridgePayload(gasLimit, maxSubmissionFee, l2MaxFeePerGas)
	if err != nil {
		return nil, nil, fmt.Errorf("pack bridge calldata for bridge adapter: %w", err)
	}

	return bridgeCalldata, deposit, nil
}

func (l *l1ToL2Bridge) estimateRetryableGasLimit(ctx context.Context, rd RetryableData) (*big.Int, error) {
	packed, err := nodeInterfaceABI.Pack("estimateRetryableTicket",
		rd.From,
		assets.Ether(1).ToInt(),
		rd.To,
		rd.L2CallValue,
		rd.ExcessFeeRefundAddr,
		rd.CallValueRefundAddr,
		rd.Data,
	)
	if err != nil {
		return nil, fmt.Errorf("pack estimateRetryableTicket call: %w", err)
	}

	gasLimit, err := l.l2Client.EstimateGas(ctx, ethereum.CallMsg{
		To:   &NodeInterfaceAddress,
		Data: packed,
	})
	if err != nil {
		return nil, fmt.Errorf("error esimtating gas on node interface for estimateRetryableTicket: %s, calldata: %s",
			err, hexutil.Encode(packed))
	}

	// no multiplier on gas limit
	// should be pretty accurate
	return big.NewInt(int64(gasLimit)), nil
}

func (l *l1ToL2Bridge) estimateMaxSubmissionFee(
	ctx context.Context,
	l1BaseFee *big.Int,
	dataLength int,
) (*big.Int, error) {
	submissionFee, err := l.l1Inbox.CalculateRetryableSubmissionFee(&bind.CallOpts{
		Context: ctx,
	}, big.NewInt(int64(dataLength)), l1BaseFee)
	if err != nil {
		return nil, fmt.Errorf("calculate retryable submission fee: %w", err)
	}

	submissionFee = submissionFee.Mul(submissionFee, submissionFeeMultiplier)
	return submissionFee, nil
}

func (l *l1ToL2Bridge) estimateMaxFeePerGasOnL2(ctx context.Context) (*big.Int, error) {
	l2BaseFee, err := l.l2Client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("suggest gas price on L2: %w", err)
	}

	l2BaseFee = l2BaseFee.Mul(l2BaseFee, l2BaseFeeMultiplier)
	return l2BaseFee, nil
}

func (l *l1ToL2Bridge) Close(ctx context.Context) error {
	return multierr.Combine(
		l.l2LogPoller.UnregisterFilter(l.l2FilterName),
		l.l1LogPoller.UnregisterFilter(l.l1FilterName),
	)
}

type RetryableData struct {
	// From is the gateway on L1 that will be sending the funds to the L2 gateway.
	From common.Address
	// To is the gateway on L2 that will be receiving the funds and eventually
	// sending them to the final recipient.
	To                common.Address
	L2CallValue       *big.Int
	Deposit           *big.Int
	MaxSubmissionCost *big.Int
	// ExcessFeeRefundAddr is an address on L2 that will be receiving excess fees
	ExcessFeeRefundAddr common.Address
	// CallValueRefundAddr is an address on L1 that will be receiving excess fees
	CallValueRefundAddr common.Address
	GasLimit            *big.Int
	MaxFeePerGas        *big.Int
	// Data is the calldata for the L2 gateway's `finalizeInboundTransfer` method.
	// The final recipient on L2 is specified in this calldata.
	Data []byte
}

func UnpackL1ToL2SendBridgePayload(payload []byte) (out arbitrum_l1_bridge_adapter.ArbitrumL1BridgeAdapterSendERC20Params, err error) {
	ifaces, err := l1AdapterABI.Methods["exposeSendERC20Params"].Inputs.UnpackValues(payload)
	if err != nil {
		return out, fmt.Errorf("unpack bridge payload: %w", err)
	}
	if len(ifaces) != 1 {
		return out, fmt.Errorf("expected 1 value, got %d", len(ifaces))
	}
	out = *abi.ConvertType(ifaces[0], new(arbitrum_l1_bridge_adapter.ArbitrumL1BridgeAdapterSendERC20Params)).(*arbitrum_l1_bridge_adapter.ArbitrumL1BridgeAdapterSendERC20Params)
	return out, nil
}

func PackL1ToL2SendBridgePayload(gasLimit, maxSubmissionCost, maxFeePerGas *big.Int) ([]byte, error) {
	return l1AdapterABI.Methods["exposeSendERC20Params"].Inputs.Pack(arbitrum_l1_bridge_adapter.ArbitrumL1BridgeAdapterSendERC20Params{
		GasLimit:          gasLimit,
		MaxSubmissionCost: maxSubmissionCost,
		MaxFeePerGas:      maxFeePerGas,
	})
}
