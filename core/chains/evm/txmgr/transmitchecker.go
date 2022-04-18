package txmgr

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	v1 "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/solidity_vrf_coordinator_interface"
	v2 "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
	bigmath "github.com/smartcontractkit/chainlink/core/utils/big_math"
)

var (
	// NoChecker is a TransmitChecker that always determines a transaction should be submitted.
	NoChecker TransmitChecker = noChecker{}

	_ TransmitCheckerFactory = &CheckerFactory{}
	_ TransmitChecker        = &SimulateChecker{}
	_ TransmitChecker        = &VRFV1Checker{}
	_ TransmitChecker        = &VRFV2Checker{}
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
	case TransmitCheckerTypeVRFV1:
		coord, err := v1.NewVRFCoordinator(spec.VRFCoordinatorAddress, c.Client)
		if err != nil {
			return nil, errors.Wrapf(err,
				"failed to create VRF V1 coordinator at address %v", spec.VRFCoordinatorAddress)
		}
		return &VRFV1Checker{coord.Callbacks}, nil
	case TransmitCheckerTypeVRFV2:
		coord, err := v2.NewVRFCoordinatorV2(spec.VRFCoordinatorAddress, c.Client)
		if err != nil {
			return nil, errors.Wrapf(err,
				"failed to create VRF V2 coordinator at address %v", spec.VRFCoordinatorAddress)
		}
		if spec.VRFRequestBlockNumber == nil {
			return nil, errors.New("VRFRequestBlockNumber parameter must be non-nil")
		}
		return &VRFV2Checker{
			GetCommitment:      coord.GetCommitment,
			HeaderByNumber:     c.Client.HeaderByNumber,
			RequestBlockNumber: spec.VRFRequestBlockNumber,
		}, nil
	case "":
		return NoChecker, nil
	default:
		return nil, errors.Errorf("unrecognized checker type: %s", spec.CheckerType)
	}
}

type noChecker struct{}

// Check satisfies the TransmitChecker interface.
func (noChecker) Check(
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
			l.Criticalw("Transaction reverted during simulation",
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

// VRFV1Checker is an implementation of TransmitChecker that checks whether a VRF V1 fulfillment
// has already been fulfilled.
type VRFV1Checker struct {

	// Callbacks checks whether a VRF V1 request has already been fulfilled on the VRFCoordinator
	// Solidity contract
	Callbacks func(opts *bind.CallOpts, reqID [32]byte) (v1.Callbacks, error)
}

// Check satisfies the TransmitChecker interface.
func (v *VRFV1Checker) Check(
	ctx context.Context,
	l logger.Logger,
	tx EthTx,
	_ EthTxAttempt,
) error {
	meta, err := tx.GetMeta()
	if err != nil {
		l.Errorw("Failed to parse transaction meta. Attempting to transmit anyway.",
			"err", err,
			"ethTxID", tx.ID,
			"meta", tx.Meta)
		return nil
	}

	if meta == nil {
		l.Errorw("Expected a non-nil meta for a VRF transaction. Attempting to transmit anyway.",
			"err", err,
			"ethTxID", tx.ID,
			"meta", tx.Meta)
		return nil
	}

	if len(meta.RequestID.Bytes()) != 32 {
		l.Errorw("Unexpected request ID. Attempting to transmit anyway.",
			"err", err,
			"ethTxID", tx.ID,
			"meta", tx.Meta)
		return nil
	}

	var reqID [32]byte
	copy(reqID[:], meta.RequestID.Bytes())
	callback, err := v.Callbacks(&bind.CallOpts{Context: ctx}, reqID)
	if err != nil {
		l.Errorw("Unable to check if already fulfilled. Attempting to transmit anyway.",
			"err", err,
			"ethTxID", tx.ID,
			"meta", tx.Meta,
			"reqID", reqID)
		return nil
	} else if utils.IsEmpty(callback.SeedAndBlockNum[:]) {
		// Request already fulfilled
		l.Infow("Request already fulfilled",
			"err", err,
			"ethTxID", tx.ID,
			"meta", tx.Meta,
			"reqID", reqID)
		return errors.New("request already fulfilled")
	} else {
		// Request not fulfilled
		return nil
	}
}

// VRFV2Checker is an implementation of TransmitChecker that checks whether a VRF V2 fulfillment
// has already been fulfilled.
type VRFV2Checker struct {

	// GetCommitment checks whether a VRF V2 request has been fulfilled on the VRFCoordinatorV2
	// Solidity contract.
	GetCommitment func(opts *bind.CallOpts, requestID *big.Int) ([32]byte, error)

	// HeaderByNumber fetches the header given the number. If nil is provided,
	// the latest header is fetched.
	HeaderByNumber func(ctx context.Context, n *big.Int) (*gethtypes.Header, error)

	// RequestBlockNumber is the block number of the VRFV2 request.
	RequestBlockNumber *big.Int
}

// Check satisfies the TransmitChecker interface.
func (v *VRFV2Checker) Check(
	ctx context.Context,
	l logger.Logger,
	tx EthTx,
	_ EthTxAttempt,
) error {
	meta, err := tx.GetMeta()
	if err != nil {
		l.Errorw("Failed to parse transaction meta. Attempting to transmit anyway.",
			"err", err,
			"ethTxID", tx.ID,
			"meta", tx.Meta)
		return nil
	}

	if meta == nil {
		l.Errorw("Expected a non-nil meta for a VRF transaction. Attempting to transmit anyway.",
			"err", err,
			"ethTxID", tx.ID,
			"meta", tx.Meta)
		return nil
	}

	h, err := v.HeaderByNumber(ctx, nil)
	if err != nil {
		l.Errorw("Failed to fetch latest header. Attempting to transmit anyway.",
			"err", err,
			"ethTxID", tx.ID,
			"meta", tx.Meta,
		)
		return nil
	}

	// If the request block number is not provided, transmit anyway just to be safe.
	// Worst we can do is revert due to the request already being fulfilled.
	if v.RequestBlockNumber == nil {
		l.Errorw("Was provided with a nil request block number. Attempting to transmit anyway.",
			"ethTxID", tx.ID,
			"meta", tx.Meta,
		)
		return nil
	}

	vrfRequestID := meta.RequestID.Big()

	// Subtract 5 since the newest block likely isn't indexed yet and will cause "header not found"
	// errors.
	latest := new(big.Int).Sub(h.Number, big.NewInt(5))
	blockNumber := bigmath.Max(latest, v.RequestBlockNumber)
	callback, err := v.GetCommitment(&bind.CallOpts{
		Context:     ctx,
		BlockNumber: blockNumber,
	}, vrfRequestID)
	if err != nil {
		l.Errorw("Failed to check request fulfillment status, error calling GetCommitment. Attempting to transmit anyway.",
			"err", err,
			"ethTxID", tx.ID,
			"meta", tx.Meta,
			"vrfRequestId", vrfRequestID,
			"blockNumber", h.Number,
		)
		return nil
	} else if utils.IsEmpty(callback[:]) {
		// If seedAndBlockNumber is zero then the response has been fulfilled and we should skip it.
		l.Infow("Request already fulfilled.",
			"ethTxID", tx.ID,
			"meta", tx.Meta,
			"vrfRequestId", vrfRequestID)
		return errors.New("request already fulfilled")
	} else {
		l.Debugw("Request not yet fulfilled",
			"ethTxID", tx.ID,
			"meta", tx.Meta,
			"vrfRequestId", vrfRequestID)
		return nil
	}
}
