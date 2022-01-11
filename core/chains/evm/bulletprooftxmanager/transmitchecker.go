package bulletprooftxmanager

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var (
	_ TransmitCheckerFactory = &CheckerFactory{}
	_ TransmitChecker        = NoChecker{}
	_ TransmitChecker        = &SimulateChecker{}
	_ TransmitChecker        = VRFV2Checker{}
)

// CheckerFactory is a real implementation of TransmitCheckerFactory.
type CheckerFactory struct {
	Client evmclient.Client
}

// BuildChecker satisfies the TransmitCheckerFactory interface.
func (c *CheckerFactory) BuildChecker(spec TransmitCheckerSpec) (TransmitChecker, error) {
	switch spec.CheckerType {
	case TransmitCheckerTypeSimulate:
		return &SimulateChecker{c.Client}, nil
	case TransmitCheckerTypeVRFV2:
		coord, err := vrf_coordinator_v2.NewVRFCoordinatorV2(spec.VRFCoordinatorAddress, c.Client)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create coordinator at address %v", spec.VRFCoordinatorAddress)
		}
		return VRFV2Checker{coord.GetCommitment}, nil
	case "":
		return NoChecker{}, nil
	default:
		return nil, errors.Errorf("unrecognized checker type: %s", spec.CheckerType)
	}
}

// NoChecker is a TransmitChecker that always determines a transaction should be submitted.
type NoChecker struct{}

// Check satisfies the TransmitChecker interface.
func (NoChecker) Check(
	_ context.Context,
	_ logger.Logger,
	_ EthTx,
	_ EthTxAttempt,
) error {
	return nil
}

// SimulateChecker simulates transactions, producing an error if they revert on chain.
type SimulateChecker struct {
	Client evmclient.Client
}

// Check satisfies the TransmitChecker interface.
func (s *SimulateChecker) Check(
	ctx context.Context,
	l logger.Logger,
	tx EthTx,
	a EthTxAttempt,
) error {
	// See: https://github.com/ethereum/go-ethereum/blob/acdf9238fb03d79c9b1c20c2fa476a7e6f4ac2ac/ethclient/gethclient/gethclient.go#L193
	callArg := map[string]interface{}{
		"from": tx.FromAddress,
		"to":   &tx.ToAddress,
		"gas":  hexutil.Uint64(a.ChainSpecificGasLimit),
		// NOTE: Deliberately do not include gas prices. We never want to fatally error a
		// transaction just because the wallet has insufficient eth.
		// Relevant info regarding EIP1559 transactions: https://github.com/ethereum/go-ethereum/pull/23027
		"gasPrice":             nil,
		"maxFeePerGas":         nil,
		"maxPriorityFeePerGas": nil,
		"value":                (*hexutil.Big)(tx.Value.ToInt()),
		"data":                 hexutil.Bytes(tx.EncodedPayload),
	}
	var b hexutil.Bytes
	// always run simulation on "latest" block
	err := s.Client.CallContext(ctx, &b, "eth_call", callArg, evmclient.ToBlockNumArg(nil))
	if err != nil {
		if jErr := evmclient.ExtractRPCError(err); jErr != nil {
			l.CriticalW("Transaction reverted during simulation",
				"ethTxAttemptID", a.ID, "txHash", a.Hash, "err", err, "rpcErr", jErr.String(), "returnValue", b.String())
			return errors.Errorf("transaction reverted during simulation: %s", jErr.String())
		}
		l.Warnw("Transaction simulation failed, will attempt to send anyway",
			"ethTxAttemptID", a.ID, "txHash", a.Hash, "err", err, "returnValue", b.String())
	} else {
		l.Debugw("Transaction simulation succeeded",
			"ethTxAttemptID", a.ID, "txHash", a.Hash, "returnValue", b.String())
	}
	return nil
}

// VRFV2Checker is an implementation of TransmitChecker that checks whether a VRF V2 fulfillment
// has already been fulfilled.
type VRFV2Checker struct {

	// GetCommitment checks whether a VRF V2 request has been fulfilled on the VRFCoordinatorV2
	// Solidity contract.
	GetCommitment func(opts *bind.CallOpts, requestID *big.Int) ([32]byte, error)
}

// Check satisfies the TransmitChecker interface.
func (v VRFV2Checker) Check(
	ctx context.Context,
	l logger.Logger,
	tx EthTx,
	_ EthTxAttempt,
) error {
	meta, err := tx.GetMeta()
	if err != nil {
		l.Errorw("Failed to parse transaction meta. Attempting to transmit anyway.",
			"err", err,
			"ethTxId", tx.ID,
			"meta", tx.Meta)
		return nil
	}

	if meta == nil {
		l.Errorw("Expected a non-nil meta for a VRF transaction. Attempting to transmit anyway.",
			"err", err,
			"ethTxId", tx.ID,
			"meta", tx.Meta)
		return nil
	}

	vrfRequestID := meta.RequestID.Big()
	callback, err := v.GetCommitment(&bind.CallOpts{Context: ctx}, vrfRequestID)
	if err != nil {
		l.Errorw("Failed to check request fulfillment status, error calling GetCommitment. Attempting to transmit anyway.",
			"err", err,
			"ethTxId", tx.ID,
			"meta", tx.Meta,
			"vrfRequestId", vrfRequestID)
		return nil
	} else if utils.IsEmpty(callback[:]) {
		// If seedAndBlockNumber is zero then the response has been fulfilled and we should skip it.
		l.Infow("Request already fulfilled.",
			"ethTxId", tx.ID,
			"meta", tx.Meta,
			"vrfRequestId", vrfRequestID)
		return errors.New("request already fulfilled")
	} else {
		l.Debugw("Request not yet fulfilled",
			"ethTxId", tx.ID,
			"meta", tx.Meta,
			"vrfRequestId", vrfRequestID)
		return nil
	}
}
