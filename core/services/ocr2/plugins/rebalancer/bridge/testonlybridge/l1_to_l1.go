package testonlybridge

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/rebalancer/generated/mock_l1_bridge_adapter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/rebalancer/generated/rebalancer"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

var (
	// Emitted on both source and destination
	LiquidityTransferredTopic = rebalancer.RebalancerLiquidityTransferred{}.Topic()
)

type testBridge struct {
	sourceSelector   models.NetworkSelector
	destSelector     models.NetworkSelector
	sourceRebalancer rebalancer.RebalancerInterface
	destRebalancer   rebalancer.RebalancerInterface
	sourceAdapter    *mock_l1_bridge_adapter.MockL1BridgeAdapter
	destAdapter      *mock_l1_bridge_adapter.MockL1BridgeAdapter
	sourceLogPoller  logpoller.LogPoller
	destLogPoller    logpoller.LogPoller
	sourceClient     client.Client
	destClient       client.Client
	lggr             logger.Logger
}

func New(
	sourceSelector, destSelector models.NetworkSelector,
	sourceRebalancerAddress, destRebalancerAddress, sourceAdapter, destAdapter models.Address,
	sourceLogPoller, destLogPoller logpoller.LogPoller,
	sourceClient, destClient client.Client,
	lggr logger.Logger,
) (*testBridge, error) {
	err := sourceLogPoller.RegisterFilter(logpoller.Filter{
		Name: logpoller.FilterName("L1-LiquidityTransferred", sourceSelector),
		EventSigs: []common.Hash{
			LiquidityTransferredTopic,
		},
		Addresses: []common.Address{
			common.Address(sourceRebalancerAddress),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("register filter for source log poller: %w", err)
	}

	err = destLogPoller.RegisterFilter(logpoller.Filter{
		Name: logpoller.FilterName("L2-LiquidityTransferred", destSelector),
		EventSigs: []common.Hash{
			LiquidityTransferredTopic,
		},
		Addresses: []common.Address{
			common.Address(destRebalancerAddress),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("register filter for dest log poller: %w", err)
	}

	lggr = lggr.Named("TestBridge").With(
		"sourceSelector", sourceSelector,
		"destSelector", destSelector,
		"sourceRebalancer", sourceRebalancerAddress,
		"destRebalancer", destRebalancerAddress,
		"sourceAdapter", sourceAdapter,
		"destAdapter", destAdapter,
	)
	lggr.Infow("TestBridge created")

	sourceAdapterWrapper, err := mock_l1_bridge_adapter.NewMockL1BridgeAdapter(common.Address(sourceAdapter), sourceClient)
	if err != nil {
		return nil, fmt.Errorf("create source adapter wrapper: %w", err)
	}

	destAdapterWrapper, err := mock_l1_bridge_adapter.NewMockL1BridgeAdapter(common.Address(destAdapter), destClient)
	if err != nil {
		return nil, fmt.Errorf("create dest adapter wrapper: %w", err)
	}

	sourceRebalancer, err := rebalancer.NewRebalancer(common.Address(sourceRebalancerAddress), sourceClient)
	if err != nil {
		return nil, fmt.Errorf("create source rebalancer: %w", err)
	}

	destRebalancer, err := rebalancer.NewRebalancer(common.Address(destRebalancerAddress), destClient)
	if err != nil {
		return nil, fmt.Errorf("create dest rebalancer: %w", err)
	}

	return &testBridge{
		sourceSelector:   sourceSelector,
		destSelector:     destSelector,
		sourceRebalancer: sourceRebalancer,
		destRebalancer:   destRebalancer,
		sourceAdapter:    sourceAdapterWrapper,
		destAdapter:      destAdapterWrapper,
		sourceLogPoller:  sourceLogPoller,
		destLogPoller:    destLogPoller,
		sourceClient:     sourceClient,
		destClient:       destClient,
		lggr:             lggr,
	}, nil
}

// Close implements bridge.Bridge.
func (t *testBridge) Close(ctx context.Context) error {
	return nil
}

// QuorumizedBridgePayload implements bridge.Bridge.
func (t *testBridge) QuorumizedBridgePayload(payloads [][]byte, f int) ([]byte, error) {
	// TODO: implement, should just return Amount and they should all be the same
	return payloads[0], nil
}

// GetBridgePayloadAndFee implements bridge.Bridge.
func (t *testBridge) GetBridgePayloadAndFee(ctx context.Context, transfer models.Transfer) ([]byte, *big.Int, error) {
	return []byte{}, big.NewInt(0), nil
}

// GetTransfers implements bridge.Bridge.
func (t *testBridge) GetTransfers(ctx context.Context, localToken models.Address, remoteToken models.Address) ([]models.PendingTransfer, error) {
	latestSourceBlock, err := t.sourceLogPoller.LatestBlock(pg.WithParentCtx(ctx))
	if err != nil {
		return nil, fmt.Errorf("get latest block: %w", err)
	}

	latestDestBlock, err := t.destLogPoller.LatestBlock(pg.WithParentCtx(ctx))
	if err != nil {
		return nil, fmt.Errorf("get latest block: %w", err)
	}

	sendLogs, err := t.sourceLogPoller.LogsWithSigs(
		1,
		latestSourceBlock.BlockNumber,
		[]common.Hash{LiquidityTransferredTopic},
		t.sourceRebalancer.Address(),
		pg.WithParentCtx(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("get source LiquidityTransferred logs: %w", err)
	}

	receiveLogs, err := t.destLogPoller.LogsWithSigs(
		1,
		latestDestBlock.BlockNumber,
		[]common.Hash{LiquidityTransferredTopic},
		t.destRebalancer.Address(),
		pg.WithParentCtx(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("get dest LiquidityTransferred logs: %w", err)
	}

	parsedSendLogs, parsedToLP, err := parseLiquidityTransferred(t.sourceRebalancer.ParseLiquidityTransferred, sendLogs)
	if err != nil {
		return nil, fmt.Errorf("parse source send logs: %w", err)
	}

	parsedFinalizeLogs, _, err := parseLiquidityTransferred(t.destRebalancer.ParseLiquidityTransferred, receiveLogs)
	if err != nil {
		return nil, fmt.Errorf("parse dest finalize logs: %w", err)
	}

	ready, err := t.getReadyToFinalize(parsedSendLogs, parsedFinalizeLogs)
	if err != nil {
		return nil, fmt.Errorf("get ready to finalize: %w", err)
	}

	return t.toPendingTransfers(localToken, remoteToken, ready, parsedToLP), nil
}

func (t *testBridge) toPendingTransfers(
	localToken, remoteToken models.Address,
	ready []*rebalancer.RebalancerLiquidityTransferred,
	parsedToLP map[logKey]logpoller.Log,
) []models.PendingTransfer {
	var transfers []models.PendingTransfer
	for _, send := range ready {
		lp := parsedToLP[logKey{txHash: send.Raw.TxHash, logIndex: int64(send.Raw.Index)}]
		sendNonce, err := UnpackBridgeSendReturnData(send.BridgeReturnData)
		if err != nil {
			t.lggr.Errorw("unpack send bridge data", "err", err)
			continue
		}
		bridgeData, err := PackFinalizeBridgePayload(send.Amount, sendNonce)
		if err != nil {
			t.lggr.Errorw("pack bridge data", "err", err)
			continue
		}
		transfers = append(transfers, models.PendingTransfer{
			Transfer: models.Transfer{
				From:               t.sourceSelector,
				To:                 t.destSelector,
				Sender:             models.Address(t.sourceAdapter.Address()),
				Receiver:           models.Address(t.destRebalancer.Address()),
				Amount:             ubig.New(send.Amount),
				LocalTokenAddress:  localToken,
				RemoteTokenAddress: remoteToken,
				Date:               lp.BlockTimestamp,
				BridgeData:         bridgeData,
			},
			Status: models.TransferStatusReady,
			ID:     fmt.Sprintf("%s-%d", send.Raw.TxHash.Hex(), send.Raw.Index),
		})
	}

	if len(transfers) > 0 {
		t.lggr.Infow("produced pending transfers", "pendingTransfers", transfers)
	}

	return transfers
}

func (t *testBridge) getReadyToFinalize(
	sends []*rebalancer.RebalancerLiquidityTransferred,
	finalizes []*rebalancer.RebalancerLiquidityTransferred,
) ([]*rebalancer.RebalancerLiquidityTransferred, error) {
	t.lggr.Debugw("Getting ready to finalize",
		"sendsLen", len(sends),
		"finalizesLen", len(finalizes),
		"sends", sends,
		"finalizes", finalizes)

	// find sent events that don't have a matching finalized event
	var ready []*rebalancer.RebalancerLiquidityTransferred
	for _, send := range sends {
		var finalized bool
		for _, finalize := range finalizes {
			sendNonce, err := UnpackBridgeSendReturnData(send.BridgeReturnData)
			if err != nil {
				return nil, fmt.Errorf("unpack send bridge data: %w", err)
			}
			_, finalizeNonce, err := UnpackFinalizeBridgePayload(finalize.BridgeSpecificData)
			if err != nil {
				return nil, fmt.Errorf("unpack finalize bridge data: %w", err)
			}
			if sendNonce.Cmp(finalizeNonce) == 0 {
				finalized = true
				break
			}
		}
		if !finalized {
			ready = append(ready, send)
		}
	}

	if len(ready) > 0 {
		t.lggr.Infow("found ready to finalize", "sendsLen", len(sends),
			"finalizesLen", len(finalizes),
			"sends", sends,
			"finalizes", finalizes,
			"ready", ready)
	} else {
		t.lggr.Debugw("no requests ready to finalize", "sendsLen", len(sends),
			"finalizesLen", len(finalizes),
			"sends", sends,
			"finalizes", finalizes)
	}

	return ready, nil
}

func PackFinalizeBridgePayload(amount, nonce *big.Int) ([]byte, error) {
	return utils.ABIEncode(`[{"type": "uint256"}, {"type": "uint256"}]`, amount, nonce)
}

func UnpackFinalizeBridgePayload(data []byte) (*big.Int, *big.Int, error) {
	ifaces, err := utils.ABIDecode(`[{"type": "uint256"}, {"type": "uint256"}]`, data)
	if err != nil {
		return nil, nil, fmt.Errorf("decode bridge data: %w", err)
	}
	if len(ifaces) != 2 {
		return nil, nil, fmt.Errorf("expected 2 arguments, got %d", len(ifaces))
	}
	val1 := *abi.ConvertType(ifaces[0], new(*big.Int)).(**big.Int)
	val2 := *abi.ConvertType(ifaces[1], new(*big.Int)).(**big.Int)
	return val1, val2, nil
}

func UnpackBridgeSendReturnData(data []byte) (*big.Int, error) {
	ifaces, err := utils.ABIDecode(`[{"type": "uint256"}]`, data)
	if err != nil {
		return nil, fmt.Errorf("decode bridge data: %w", err)
	}
	if len(ifaces) != 1 {
		return nil, fmt.Errorf("expected 1 argument, got %d", len(ifaces))
	}
	val := *abi.ConvertType(ifaces[0], new(*big.Int)).(**big.Int)
	return val, nil
}

type logKey struct {
	txHash   common.Hash
	logIndex int64
}

func parseLiquidityTransferred(parseFunc func(gethtypes.Log) (*rebalancer.RebalancerLiquidityTransferred, error), lgs []logpoller.Log) ([]*rebalancer.RebalancerLiquidityTransferred, map[logKey]logpoller.Log, error) {
	transferred := make([]*rebalancer.RebalancerLiquidityTransferred, len(lgs))
	toLP := make(map[logKey]logpoller.Log)
	for i, lg := range lgs {
		parsed, err := parseFunc(lg.ToGethLog())
		if err != nil {
			// should never happen
			return nil, nil, fmt.Errorf("parse LiquidityTransferred log: %w", err)
		}
		transferred[i] = parsed
		toLP[logKey{
			txHash:   lg.TxHash,
			logIndex: lg.LogIndex,
		}] = lg
	}
	return transferred, toLP, nil
}
