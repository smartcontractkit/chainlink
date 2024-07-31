package testonlybridge

import (
	"bytes"
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
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/liquiditymanager"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/mock_l1_bridge_adapter"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/abiutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

const (
	// These correspond to the enumeration FinalizationAction in the MockL1BridgeAdapter contract.
	FinalizationActionProveWithdrawal    uint8 = 0
	FinalizationActionFinalizeWithdrawal uint8 = 1
)

var (
	adapterABI = abihelpers.MustParseABI(mock_l1_bridge_adapter.MockL1BridgeAdapterABI)

	// Emitted on both source and destination
	LiquidityTransferredTopic      = liquiditymanager.LiquidityManagerLiquidityTransferred{}.Topic()
	FinalizationStepCompletedTopic = liquiditymanager.LiquidityManagerFinalizationStepCompleted{}.Topic()
)

type testBridge struct {
	sourceSelector         models.NetworkSelector
	destSelector           models.NetworkSelector
	sourceLiquidityManager liquiditymanager.LiquidityManagerInterface
	destLiquidityManager   liquiditymanager.LiquidityManagerInterface
	sourceAdapter          mock_l1_bridge_adapter.MockL1BridgeAdapterInterface
	destAdapter            mock_l1_bridge_adapter.MockL1BridgeAdapterInterface
	sourceLogPoller        logpoller.LogPoller
	destLogPoller          logpoller.LogPoller
	sourceClient           client.Client
	destClient             client.Client
	lggr                   logger.Logger
}

func New(
	ctx context.Context,
	sourceSelector, destSelector models.NetworkSelector,
	sourceLiquidityManagerAddress, destLiquidityManagerAddress, sourceAdapter, destAdapter models.Address,
	sourceLogPoller, destLogPoller logpoller.LogPoller,
	sourceClient, destClient client.Client,
	lggr logger.Logger,
) (*testBridge, error) {
	err := sourceLogPoller.RegisterFilter(
		ctx,
		logpoller.Filter{
			Name: logpoller.FilterName("Local-LiquidityTransferred-FinalizationCompleted",
				sourceSelector, sourceLiquidityManagerAddress.String()),
			EventSigs: []common.Hash{
				LiquidityTransferredTopic,
				FinalizationStepCompletedTopic,
			},
			Addresses: []common.Address{
				common.Address(sourceLiquidityManagerAddress),
			},
		})
	if err != nil {
		return nil, fmt.Errorf("register filter for source log poller: %w", err)
	}

	err = destLogPoller.RegisterFilter(
		ctx,
		logpoller.Filter{
			Name: logpoller.FilterName("Remote-LiquidityTransferred-FinalizationCompleted",
				destSelector, destLiquidityManagerAddress.String()),
			EventSigs: []common.Hash{
				LiquidityTransferredTopic,
				FinalizationStepCompletedTopic,
			},
			Addresses: []common.Address{
				common.Address(destLiquidityManagerAddress),
			},
		})
	if err != nil {
		return nil, fmt.Errorf("register filter for dest log poller: %w", err)
	}

	lggr = lggr.Named("TestBridge").With(
		"sourceSelector", sourceSelector,
		"destSelector", destSelector,
		"sourceLiquidityManager", sourceLiquidityManagerAddress,
		"destLiquidityManager", destLiquidityManagerAddress,
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

	sourceLiquidityManager, err := liquiditymanager.NewLiquidityManager(common.Address(sourceLiquidityManagerAddress), sourceClient)
	if err != nil {
		return nil, fmt.Errorf("create source liquidityManager: %w", err)
	}

	destLiquidityManager, err := liquiditymanager.NewLiquidityManager(common.Address(destLiquidityManagerAddress), destClient)
	if err != nil {
		return nil, fmt.Errorf("create dest liquidityManager: %w", err)
	}

	return &testBridge{
		sourceSelector:         sourceSelector,
		destSelector:           destSelector,
		sourceLiquidityManager: sourceLiquidityManager,
		destLiquidityManager:   destLiquidityManager,
		sourceAdapter:          sourceAdapterWrapper,
		destAdapter:            destAdapterWrapper,
		sourceLogPoller:        sourceLogPoller,
		destLogPoller:          destLogPoller,
		sourceClient:           sourceClient,
		destClient:             destClient,
		lggr:                   lggr,
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
	latestSourceBlock, err := t.sourceLogPoller.LatestBlock(ctx)
	if err != nil {
		return nil, fmt.Errorf("get latest block: %w", err)
	}

	latestDestBlock, err := t.destLogPoller.LatestBlock(ctx)
	if err != nil {
		return nil, fmt.Errorf("get latest block: %w", err)
	}

	sendLogs, err := t.sourceLogPoller.LogsWithSigs(
		ctx,
		1,
		latestSourceBlock.BlockNumber,
		[]common.Hash{LiquidityTransferredTopic},
		t.sourceLiquidityManager.Address(),
	)
	if err != nil {
		return nil, fmt.Errorf("get source LiquidityTransferred logs: %w", err)
	}

	receiveLogs, err := t.destLogPoller.LogsWithSigs(
		ctx,
		1,
		latestDestBlock.BlockNumber,
		[]common.Hash{LiquidityTransferredTopic, FinalizationStepCompletedTopic},
		t.destLiquidityManager.Address(),
	)
	if err != nil {
		return nil, fmt.Errorf("get dest LiquidityTransferred logs: %w", err)
	}

	parsedSendLogs, parsedToLP, err := parseLiquidityTransferred(t.sourceLiquidityManager.ParseLiquidityTransferred, sendLogs)
	if err != nil {
		return nil, fmt.Errorf("parse source send logs: %w", err)
	}

	parsedStepCompleted, parsedFinalizeLogs, err := parseLiquidityTransferredAndFinalizationStepCompleted(
		t.destLiquidityManager.ParseLiquidityTransferred,
		t.destLiquidityManager.ParseFinalizationStepCompleted,
		receiveLogs)
	if err != nil {
		return nil, fmt.Errorf("parse dest finalize logs: %w", err)
	}

	readyToProve, readyToFinalize, err := filterAndGroupByStage(t.lggr, parsedSendLogs, parsedFinalizeLogs, parsedStepCompleted)
	if err != nil {
		return nil, fmt.Errorf("get ready to finalize: %w", err)
	}

	return t.toPendingTransfers(localToken, remoteToken, readyToProve, readyToFinalize, parsedToLP)
}

func (t *testBridge) toPendingTransfers(
	localToken, remoteToken models.Address,
	readyToProve,
	readyToFinalize []*liquiditymanager.LiquidityManagerLiquidityTransferred,
	parsedToLP map[logKey]logpoller.Log,
) ([]models.PendingTransfer, error) {
	var transfers []models.PendingTransfer

	for _, send := range readyToProve {
		lp := parsedToLP[logKey{txHash: send.Raw.TxHash, logIndex: int64(send.Raw.Index)}]
		sendNonce, err := UnpackBridgeSendReturnData(send.BridgeReturnData)
		if err != nil {
			return nil, fmt.Errorf("unpack send bridge data %x: %w", send.BridgeReturnData, err)
		}
		bridgeData, err := PackProveBridgePayload(sendNonce)
		if err != nil {
			return nil, fmt.Errorf("pack bridge data (%s): %w", sendNonce.String(), err)
		}
		transfers = append(transfers, models.PendingTransfer{
			Transfer: models.Transfer{
				From:               t.sourceSelector,
				To:                 t.destSelector,
				Sender:             models.Address(t.sourceAdapter.Address()),
				Receiver:           models.Address(t.destLiquidityManager.Address()),
				Amount:             ubig.New(send.Amount),
				LocalTokenAddress:  localToken,
				RemoteTokenAddress: remoteToken,
				Date:               lp.BlockTimestamp,
				BridgeData:         bridgeData,
				Stage:              1,
			},
			Status: models.TransferStatusReady,
			ID:     fmt.Sprintf("%s-%d-prove", send.Raw.TxHash.Hex(), send.Raw.Index),
		})
	}

	for _, send := range readyToFinalize {
		lp := parsedToLP[logKey{txHash: send.Raw.TxHash, logIndex: int64(send.Raw.Index)}]
		sendNonce, err := UnpackBridgeSendReturnData(send.BridgeReturnData)
		if err != nil {
			return nil, fmt.Errorf("unpack send bridge data %x: %w", send.BridgeReturnData, err)
		}
		bridgeData, err := PackFinalizeBridgePayload(send.Amount, sendNonce)
		if err != nil {
			return nil, fmt.Errorf("pack bridge data (%s): %w", sendNonce.String(), err)
		}
		transfers = append(transfers, models.PendingTransfer{
			Transfer: models.Transfer{
				From:               t.sourceSelector,
				To:                 t.destSelector,
				Sender:             models.Address(t.sourceAdapter.Address()),
				Receiver:           models.Address(t.destLiquidityManager.Address()),
				Amount:             ubig.New(send.Amount),
				LocalTokenAddress:  localToken,
				RemoteTokenAddress: remoteToken,
				Date:               lp.BlockTimestamp,
				BridgeData:         bridgeData,
				Stage:              2,
			},
			Status: models.TransferStatusReady,
			ID:     fmt.Sprintf("%s-%d-finalize", send.Raw.TxHash.Hex(), send.Raw.Index),
		})
	}

	if len(transfers) > 0 {
		t.lggr.Infow("produced pending transfers", "pendingTransfers", transfers)
	}

	return transfers, nil
}

// filterAndGroupByStage filters out sends that have already been finalized
// and groups the remaining sends into ready to prove and ready to finalize slices.
func filterAndGroupByStage(
	lggr logger.Logger,
	sends []*liquiditymanager.LiquidityManagerLiquidityTransferred,
	finalizes []*liquiditymanager.LiquidityManagerLiquidityTransferred,
	stepsCompleted []*liquiditymanager.LiquidityManagerFinalizationStepCompleted,
) (readyToProve, readyToFinalize []*liquiditymanager.LiquidityManagerLiquidityTransferred, err error) {
	lggr = lggr.With(
		"sendsLen", len(sends),
		"finalizesLen", len(finalizes),
		"stepsCompletedLen", len(stepsCompleted),
		"sends", sends,
		"finalizes", finalizes,
		"stepsCompleted", stepsCompleted)
	lggr.Debugw("Getting ready to finalize")

	// find sent events that don't have a matching finalized event
	unfinalized, err := filterFinalized(sends, finalizes)
	if err != nil {
		return nil, nil, fmt.Errorf("filter finalized: %w", err)
	}

	// group remaining unfinalized sends into ready to finalize and ready to prove.
	// ready to finalize sends will be finalized, while ready to prove will be proven.
	readyToProve, readyToFinalize, err = groupByStage(unfinalized, stepsCompleted)
	if err != nil {
		return nil, nil, fmt.Errorf("group by stage: %w", err)
	}

	if len(readyToFinalize) > 0 {
		lggr.Infow("found proven sends, ready to finalize",
			"provenSendsLen", len(readyToFinalize),
			"readyToFinalize", readyToFinalize)
	}
	if len(readyToProve) > 0 {
		lggr.Infow("found unproven sends, ready to prove",
			"unprovenSendsLen", len(readyToProve),
			"readyToProve", readyToProve)
	}

	if len(readyToFinalize) == 0 && len(readyToProve) == 0 {
		lggr.Debugw("No sends ready to finalize or prove",
			"finalizes", finalizes)
	}

	return
}

// groupByStage groups the unfinalized transfers into two categories: ready to prove and ready to finalize.
func groupByStage(
	unfinalized []*liquiditymanager.LiquidityManagerLiquidityTransferred,
	stepsCompleted []*liquiditymanager.LiquidityManagerFinalizationStepCompleted,
) (
	readyToProve,
	readyToFinalize []*liquiditymanager.LiquidityManagerLiquidityTransferred,
	err error,
) {
	for _, candidate := range unfinalized {
		proven, err := isCandidateProven(candidate, stepsCompleted)
		if err != nil {
			return nil, nil, fmt.Errorf("new function: %w", err)
		}

		if proven {
			readyToFinalize = append(readyToFinalize, candidate)
		} else {
			readyToProve = append(readyToProve, candidate)
		}
	}
	return
}

// isCandidateProven returns true if the candidate transfer has already been proven.
// it does this by checking if the candidate's nonce matches any of the nonces in the
// stepsCompleted logs.
// see contracts/src/v0.8/liquiditymanager/test/mocks/MockBridgeAdapter.sol for details on this.
func isCandidateProven(candidate *liquiditymanager.LiquidityManagerLiquidityTransferred, stepsCompleted []*liquiditymanager.LiquidityManagerFinalizationStepCompleted) (bool, error) {
	for _, stepCompleted := range stepsCompleted {
		sendNonce, err := UnpackBridgeSendReturnData(candidate.BridgeReturnData)
		if err != nil {
			return false, fmt.Errorf("unpack send bridge data: %w", err)
		}
		proveNonce, err := UnpackProveBridgePayload(stepCompleted.BridgeSpecificData)
		if err != nil {
			return false, fmt.Errorf("unpack prove bridge data: %w", err)
		}
		if proveNonce.Cmp(sendNonce) == 0 {
			return true, nil
		}
	}
	return false, nil
}

func filterFinalized(
	sends []*liquiditymanager.LiquidityManagerLiquidityTransferred,
	finalizes []*liquiditymanager.LiquidityManagerLiquidityTransferred) (
	[]*liquiditymanager.LiquidityManagerLiquidityTransferred,
	error,
) {
	var unfinalized []*liquiditymanager.LiquidityManagerLiquidityTransferred
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
			unfinalized = append(unfinalized, send)
		}
	}
	return unfinalized, nil
}

func PackProveBridgePayload(nonce *big.Int) ([]byte, error) {
	encodedProvePayload, err := adapterABI.Methods["encodeProvePayload"].Inputs.Pack(mock_l1_bridge_adapter.MockL1BridgeAdapterProvePayload{
		Nonce: nonce,
	})
	if err != nil {
		return nil, fmt.Errorf("pack prove bridge data: %w", err)
	}

	encodedPayload, err := adapterABI.Methods["encodePayload"].Inputs.Pack(
		mock_l1_bridge_adapter.MockL1BridgeAdapterPayload{
			Action: FinalizationActionProveWithdrawal,
			Data:   encodedProvePayload,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("pack bridge data: %w", err)
	}

	return encodedPayload, nil
}

func PackFinalizeBridgePayload(amount, nonce *big.Int) ([]byte, error) {
	encodedFinalizePayload, err := adapterABI.Methods["encodeFinalizePayload"].Inputs.Pack(mock_l1_bridge_adapter.MockL1BridgeAdapterFinalizePayload{
		Amount: amount,
		Nonce:  nonce,
	})
	if err != nil {
		return nil, fmt.Errorf("pack finalize bridge data: %w", err)
	}

	encodedPayload, err := adapterABI.Methods["encodePayload"].Inputs.Pack(
		mock_l1_bridge_adapter.MockL1BridgeAdapterPayload{
			Action: FinalizationActionFinalizeWithdrawal,
			Data:   encodedFinalizePayload,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("pack bridge data: %w", err)
	}

	return encodedPayload, nil
}

func UnpackProveBridgePayload(data []byte) (*big.Int, error) {
	ifaces, err := adapterABI.Methods["encodePayload"].Inputs.Unpack(data)
	if err != nil {
		return nil, fmt.Errorf("decode prove bridge data: %w", err)
	}

	if len(ifaces) != 1 {
		return nil, fmt.Errorf("decode payload: expected 1 argument, got %d", len(ifaces))
	}

	payload := *abi.ConvertType(ifaces[0], new(mock_l1_bridge_adapter.MockL1BridgeAdapterPayload)).(*mock_l1_bridge_adapter.MockL1BridgeAdapterPayload)

	// decode the prove payload from the payload
	proveIfaces, err := adapterABI.Methods["encodeProvePayload"].Inputs.Unpack(payload.Data)
	if err != nil {
		return nil, fmt.Errorf("decode prove payload: %w", err)
	}

	if len(proveIfaces) != 1 {
		return nil, fmt.Errorf("decode prove payload: expected 1 argument, got %d", len(proveIfaces))
	}

	provePayload := *abi.ConvertType(proveIfaces[0], new(mock_l1_bridge_adapter.MockL1BridgeAdapterProvePayload)).(*mock_l1_bridge_adapter.MockL1BridgeAdapterProvePayload)

	return provePayload.Nonce, nil
}

func UnpackFinalizeBridgePayload(data []byte) (*big.Int, *big.Int, error) {
	ifaces, err := adapterABI.Methods["encodePayload"].Inputs.Unpack(data)
	if err != nil {
		return nil, nil, fmt.Errorf("decode prove bridge data: %w", err)
	}

	if len(ifaces) != 1 {
		return nil, nil, fmt.Errorf("decode payload: expected 1 argument, got %d", len(ifaces))
	}

	payload := *abi.ConvertType(ifaces[0], new(mock_l1_bridge_adapter.MockL1BridgeAdapterPayload)).(*mock_l1_bridge_adapter.MockL1BridgeAdapterPayload)

	// decode the finalize payload from the payload
	finalizeIfaces, err := adapterABI.Methods["encodeFinalizePayload"].Inputs.Unpack(payload.Data)
	if err != nil {
		return nil, nil, fmt.Errorf("decode finalize payload: %w", err)
	}

	if len(finalizeIfaces) != 1 {
		return nil, nil, fmt.Errorf("decode finalize payload: expected 1 argument1, got %d", len(finalizeIfaces))
	}

	finalizePayload := *abi.ConvertType(finalizeIfaces[0], new(mock_l1_bridge_adapter.MockL1BridgeAdapterFinalizePayload)).(*mock_l1_bridge_adapter.MockL1BridgeAdapterFinalizePayload)

	return finalizePayload.Amount, finalizePayload.Nonce, nil
}

func UnpackBridgeSendReturnData(data []byte) (*big.Int, error) {
	return abiutils.UnpackUint256(data)
}

func PackBridgeSendReturnData(nonce *big.Int) ([]byte, error) {
	return utils.ABIEncode(`[{"type": "uint256"}]`, nonce)
}

type logKey struct {
	txHash   common.Hash
	logIndex int64
}

func parseLiquidityTransferred(parseFunc func(gethtypes.Log) (*liquiditymanager.LiquidityManagerLiquidityTransferred, error), lgs []logpoller.Log) ([]*liquiditymanager.LiquidityManagerLiquidityTransferred, map[logKey]logpoller.Log, error) {
	transferred := make([]*liquiditymanager.LiquidityManagerLiquidityTransferred, len(lgs))
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

func parseLiquidityTransferredAndFinalizationStepCompleted(
	transferredParse func(gethtypes.Log) (*liquiditymanager.LiquidityManagerLiquidityTransferred, error),
	finalizeParse func(gethtypes.Log) (*liquiditymanager.LiquidityManagerFinalizationStepCompleted, error),
	lgs []logpoller.Log) (
	[]*liquiditymanager.LiquidityManagerFinalizationStepCompleted,
	[]*liquiditymanager.LiquidityManagerLiquidityTransferred,
	error,
) {
	var finalizationStepCompletedLogs []*liquiditymanager.LiquidityManagerFinalizationStepCompleted
	var liquidityTransferredLogs []*liquiditymanager.LiquidityManagerLiquidityTransferred
	for _, lg := range lgs {
		if bytes.Equal(lg.Topics[0], LiquidityTransferredTopic.Bytes()) {
			parsed, err := transferredParse(lg.ToGethLog())
			if err != nil {
				// should never happen
				return nil, nil, fmt.Errorf("parse LiquidityTransferred log: %w", err)
			}
			liquidityTransferredLogs = append(liquidityTransferredLogs, parsed)
		} else if bytes.Equal(lg.Topics[0], FinalizationStepCompletedTopic.Bytes()) {
			parsed, err := finalizeParse(lg.ToGethLog())
			if err != nil {
				// should never happen
				return nil, nil, fmt.Errorf("parse FinalizationStepCompleted log: %w", err)
			}
			finalizationStepCompletedLogs = append(finalizationStepCompletedLogs, parsed)
		} else {
			return nil, nil, fmt.Errorf("unexpected topic: %x", lg.Topics[0])
		}
	}
	return finalizationStepCompletedLogs, liquidityTransferredLogs, nil
}
