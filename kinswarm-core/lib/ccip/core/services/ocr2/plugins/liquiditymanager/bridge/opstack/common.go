package opstack

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/liquiditymanager"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/optimism_cross_domain_messenger"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/optimism_l1_bridge_adapter_encoder"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/optimism_l1_standard_bridge"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/optimism_standard_bridge"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/abiutils"
)

const (
	// ERC20BridgeFinalizedFromAddressTopicIndex is the index of the topic in the ERC20BridgeFinalized event
	// that contains the "from" address. In the case of an L1 to L2 transfer, this event will be emitted by the OP
	// StandardBridge on L2 and the "from" address should be the L1 bridge adapter contract address.
	ERC20BridgeFinalizedFromAddressTopicIndex = 3

	// Optimism stages
	// StageRebalanceConfirmed is set as the transfer stage when the rebalanceLiquidity tx is confirmed onchain, but
	// when it has not yet been finalized.
	StageRebalanceConfirmed = 1
	// StageFinalizeReady is set as the transfer stage when the finalization is ready to execute onchain.
	StageFinalizeReady = 2
	// StageFinalizeConfirmed is set as the transfer stage when the finalization is confirmed onchain.
	// This is a terminal stage.
	StageFinalizeConfirmed = 3

	// Function calls
	DepositETHToFunction = "depositETHTo"
)

var (
	// Optimism events emitted on L2
	ERC20BridgeFinalizedTopic = optimism_standard_bridge.OptimismStandardBridgeERC20BridgeFinalized{}.Topic()

	// ABIs
	l1standardBridgeABI         = abihelpers.MustParseABI(optimism_l1_standard_bridge.OptimismL1StandardBridgeMetaData.ABI)
	l1OPBridgeAdapterEncoderABI = abihelpers.MustParseABI(optimism_l1_bridge_adapter_encoder.OptimismL1BridgeAdapterEncoderMetaData.ABI)
	opCrossDomainMessengerABI   = abihelpers.MustParseABI(optimism_cross_domain_messenger.OptimismCrossDomainMessengerMetaData.ABI)
	opStandardBridgeABI         = abihelpers.MustParseABI(optimism_standard_bridge.OptimismStandardBridgeMetaData.ABI)
)

/**
 * filterExecuted filters out the transfers that have already been executed on the destination chain.
 * @param readyCandidates The initiating transfer logs that are emitted on the source chain when a transfer is issued
 * @param receivedLogs The logs emitted on the destination chain when a transfer is received
 */
func filterExecuted(
	readyCandidates []*liquiditymanager.LiquidityManagerLiquidityTransferred,
	receivedLogs []*liquiditymanager.LiquidityManagerLiquidityTransferred,
) (
	ready []*liquiditymanager.LiquidityManagerLiquidityTransferred,
	err error,
) {
	for _, readyCandidate := range readyCandidates {
		exists, err := matchingExecutionExists(readyCandidate, receivedLogs)
		if err != nil {
			return nil, fmt.Errorf("error checking if ready candidate has been executed: %w", err)
		}
		if !exists {
			ready = append(ready, readyCandidate)
		}
	}
	return
}

// We encode the nonce (which is used as a unique ID for identifying a given transfer) in the bridgeSpecificData field
// of the receiving LiquidityTransferred log. We can use this to match the sent and received logs.
func matchingExecutionExists(
	readyCandidate *liquiditymanager.LiquidityManagerLiquidityTransferred,
	receivedLogs []*liquiditymanager.LiquidityManagerLiquidityTransferred,
) (bool, error) {
	// The nonce is included in the bridgeReturnData when it is emitted as a sent LiquidityTransferred event.
	sendLogNonceID, err := abiutils.UnpackUint256(readyCandidate.BridgeReturnData)
	if err != nil {
		return false, fmt.Errorf("unpack sendLogNonceID from send LiquidityTransferred log (%s): %w, BridgeReturnData: %s",
			readyCandidate.Raw.TxHash, err, hexutil.Encode(readyCandidate.BridgeReturnData))
	}
	// On the receiving side, the nonce is stored in the BridgeSpecificData field instead
	for _, receivedLog := range receivedLogs {
		receiveLogNonceID, err := abiutils.UnpackUint256(receivedLog.BridgeSpecificData)
		if err != nil {
			return false, fmt.Errorf("unpack receiveLogNonceID from receive LiquidityTransferred log (%s): %w, BridgeSpecificData: %s",
				receivedLog.Raw.TxHash, err, hexutil.Encode(receivedLog.BridgeSpecificData))
		}
		if sendLogNonceID.Cmp(receiveLogNonceID) == 0 {
			return true, nil
		}
	}
	return false, nil
}
