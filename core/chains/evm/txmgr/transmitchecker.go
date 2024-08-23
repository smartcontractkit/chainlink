package txmgr

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	pkgerrors "github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	bigmath "github.com/smartcontractkit/chainlink-common/pkg/utils/big_math"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/bytes"

	"github.com/smartcontractkit/chainlink/v2/common/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	v1 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/solidity_vrf_coordinator_interface"
	v2 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2plus_interface"
)

type (
	TransmitChecker     = txmgr.TransmitChecker[*big.Int, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee]
	TransmitCheckerSpec = txmgrtypes.TransmitCheckerSpec[common.Address]
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
		if spec.VRFCoordinatorAddress == nil {
			return nil, pkgerrors.Errorf("malformed checker, expected non-nil VRFCoordinatorAddress, got: %v", spec)
		}
		coord, err := v1.NewVRFCoordinator(*spec.VRFCoordinatorAddress, c.Client)
		if err != nil {
			return nil, pkgerrors.Wrapf(err,
				"failed to create VRF V1 coordinator at address %v", spec.VRFCoordinatorAddress)
		}
		return &VRFV1Checker{
			Callbacks: coord.Callbacks,
			Client:    c.Client,
		}, nil
	case TransmitCheckerTypeVRFV2:
		if spec.VRFCoordinatorAddress == nil {
			return nil, pkgerrors.Errorf("malformed checker, expected non-nil VRFCoordinatorAddress, got: %v", spec)
		}
		coord, err := v2.NewVRFCoordinatorV2(*spec.VRFCoordinatorAddress, c.Client)
		if err != nil {
			return nil, pkgerrors.Wrapf(err,
				"failed to create VRF V2 coordinator at address %v", spec.VRFCoordinatorAddress)
		}
		if spec.VRFRequestBlockNumber == nil {
			return nil, pkgerrors.New("VRFRequestBlockNumber parameter must be non-nil")
		}
		return &VRFV2Checker{
			GetCommitment:      coord.GetCommitment,
			HeadByNumber:       c.Client.HeadByNumber,
			RequestBlockNumber: spec.VRFRequestBlockNumber,
		}, nil
	case TransmitCheckerTypeVRFV2Plus:
		if spec.VRFCoordinatorAddress == nil {
			return nil, pkgerrors.Errorf("malformed checker, expected non-nil VRFCoordinatorAddress, got: %v", spec)
		}
		coord, err := vrf_coordinator_v2plus_interface.NewIVRFCoordinatorV2PlusInternal(*spec.VRFCoordinatorAddress, c.Client)
		if err != nil {
			return nil, pkgerrors.Wrapf(err,
				"failed to create VRF V2 coordinator plus at address %v", spec.VRFCoordinatorAddress)
		}
		if spec.VRFRequestBlockNumber == nil {
			return nil, pkgerrors.New("VRFRequestBlockNumber parameter must be non-nil")
		}
		return &VRFV2Checker{
			GetCommitment:      coord.SRequestCommitments,
			HeadByNumber:       c.Client.HeadByNumber,
			RequestBlockNumber: spec.VRFRequestBlockNumber,
		}, nil
	case "":
		return NoChecker, nil
	default:
		return nil, pkgerrors.Errorf("unrecognized checker type: %s", spec.CheckerType)
	}
}

type noChecker struct{}

// Check satisfies the TransmitChecker interface.
func (noChecker) Check(
	_ context.Context,
	_ logger.SugaredLogger,
	_ Tx,
	_ TxAttempt,
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
	l logger.SugaredLogger,
	tx Tx,
	a TxAttempt,
) error {
	// See: https://github.com/ethereum/go-ethereum/blob/acdf9238fb03d79c9b1c20c2fa476a7e6f4ac2ac/ethclient/gethclient/gethclient.go#L193
	callArg := map[string]interface{}{
		"from": tx.FromAddress,
		"to":   &tx.ToAddress,
		"gas":  hexutil.Uint64(a.ChainSpecificFeeLimit),
		// NOTE: Deliberately do not include gas prices. We never want to fatally error a
		// transaction just because the wallet has insufficient eth.
		// Relevant info regarding EIP1559 transactions: https://github.com/ethereum/go-ethereum/pull/23027
		"gasPrice":             nil,
		"maxFeePerGas":         nil,
		"maxPriorityFeePerGas": nil,
		"value":                (*hexutil.Big)(&tx.Value),
		"data":                 hexutil.Bytes(tx.EncodedPayload),
	}
	var b hexutil.Bytes
	// always run simulation on "latest" block
	err := s.Client.CallContext(ctx, &b, "eth_call", callArg, evmclient.ToBlockNumArg(nil))
	if err != nil {
		if jErr := evmclient.ExtractRPCErrorOrNil(err); jErr != nil {
			l.Criticalw("Transaction reverted during simulation",
				"ethTxAttemptID", a.ID, "txHash", a.Hash, "err", err, "rpcErr", jErr.String(), "returnValue", b.String())
			return pkgerrors.Errorf("transaction reverted during simulation: %s", jErr.String())
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

	Client evmclient.Client
}

// Check satisfies the TransmitChecker interface.
func (v *VRFV1Checker) Check(
	ctx context.Context,
	l logger.SugaredLogger,
	tx Tx,
	_ TxAttempt,
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

	if meta.RequestTxHash == nil {
		l.Errorw("Request tx hash is nil. Attempting to transmit anyway.",
			"err", err,
			"ethTxID", tx.ID,
			"meta", tx.Meta)
		return nil
	}

	// Construct and execute batch call to retrieve most the recent block number and the
	// block number of the request transaction.
	mostRecentHead := &types.Head{}
	requestTransactionReceipt := &gethtypes.Receipt{}
	batch := []rpc.BatchElem{{
		Method: "eth_getBlockByNumber",
		Args:   []interface{}{"latest", false},
		Result: mostRecentHead,
	}, {
		Method: "eth_getTransactionReceipt",
		Args:   []interface{}{*meta.RequestTxHash},
		Result: requestTransactionReceipt,
	}}
	err = v.Client.BatchCallContext(ctx, batch)
	if err != nil {
		l.Errorw("Failed to fetch latest header and transaction receipt. Attempting to transmit anyway.",
			"err", err,
			"ethTxID", tx.ID,
			"meta", tx.Meta,
		)
		return nil
	}

	// Subtract 5 since the newest block likely isn't indexed yet and will cause "header not found"
	// errors.
	latest := new(big.Int).Sub(big.NewInt(mostRecentHead.Number), big.NewInt(5))
	blockNumber := bigmath.Max(latest, requestTransactionReceipt.BlockNumber)
	var reqID [32]byte
	copy(reqID[:], meta.RequestID.Bytes())
	callback, err := v.Callbacks(&bind.CallOpts{
		Context:     ctx,
		BlockNumber: blockNumber,
	}, reqID)
	if err != nil {
		l.Errorw("Unable to check if already fulfilled. Attempting to transmit anyway.",
			"err", err,
			"ethTxID", tx.ID,
			"meta", tx.Meta,
			"reqID", reqID)
		return nil
	} else if bytes.IsEmpty(callback.SeedAndBlockNum[:]) {
		// Request already fulfilled
		l.Infow("Request already fulfilled",
			"err", err,
			"ethTxID", tx.ID,
			"meta", tx.Meta,
			"reqID", reqID)
		return pkgerrors.New("request already fulfilled")
	}
	// Request not fulfilled
	return nil
}

// VRFV2Checker is an implementation of TransmitChecker that checks whether a VRF V2 fulfillment
// has already been fulfilled.
type VRFV2Checker struct {

	// GetCommitment checks whether a VRF V2 request has been fulfilled on the VRFCoordinatorV2
	// Solidity contract.
	GetCommitment func(opts *bind.CallOpts, requestID *big.Int) ([32]byte, error)

	// HeadByNumber fetches the head given the number. If nil is provided,
	// the latest header is fetched.
	HeadByNumber func(ctx context.Context, n *big.Int) (*types.Head, error)

	// RequestBlockNumber is the block number of the VRFV2 request.
	RequestBlockNumber *big.Int
}

// Check satisfies the TransmitChecker interface.
func (v *VRFV2Checker) Check(
	ctx context.Context,
	l logger.SugaredLogger,
	tx Tx,
	_ TxAttempt,
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

	h, err := v.HeadByNumber(ctx, nil)
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
	latest := new(big.Int).Sub(big.NewInt(h.Number), big.NewInt(5))
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
	} else if bytes.IsEmpty(callback[:]) {
		// If seedAndBlockNumber is zero then the response has been fulfilled and we should skip it.
		l.Infow("Request already fulfilled.",
			"ethTxID", tx.ID,
			"meta", tx.Meta,
			"vrfRequestId", vrfRequestID)
		return pkgerrors.New("request already fulfilled")
	}
	l.Debugw("Request not yet fulfilled",
		"ethTxID", tx.ID,
		"meta", tx.Meta,
		"vrfRequestId", vrfRequestID)
	return nil
}
