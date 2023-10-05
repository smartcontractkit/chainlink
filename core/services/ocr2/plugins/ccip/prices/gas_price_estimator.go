package prices

import (
	"context"
	"math/big"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal"
)

const (
	feeBoostingOverheadGas = 200_000
	// execGasPerToken is lower-bound estimation of ERC20 releaseOrMint gas cost (Mint with static minter).
	// Use this in per-token gas cost calc as heuristic to simplify estimation logic.
	execGasPerToken = 10_000
	// execGasPerPayloadByte is gas charged for passing each byte of `data` payload to CCIP receiver, ignores 4 gas per 0-byte rule.
	// This can be a constant as it is part of EVM spec. Changes should be rare.
	execGasPerPayloadByte = 16
	// evmMessageFixedBytes is byte size of fixed-size fields in EVM2EVMMessage
	// Updating EVM2EVMMessage involves an offchain upgrade, safe to keep this as constant in code.
	evmMessageFixedBytes     = 448
	evmMessageBytesPerToken  = 128          // Byte size of each token transfer, consisting of 1 EVMTokenAmount and 1 bytes, excl length of bytes
	daMultiplierBase         = int64(10000) // DA multiplier is in multiples of 0.0001, i.e. 1/daMultiplierBase
	daGasPriceEncodingLength = 112          // Each gas price takes up at most GasPriceEncodingLength number of bits
)

// GasPrice represents gas price as a single big.Int, same as gas price representation onchain.
// (multi-component gas prices are encoded into the int)
type GasPrice *big.Int

// gasPriceEstimatorCommon is abstraction over multi-component gas prices.
type gasPriceEstimatorCommon interface {
	// GetGasPrice fetches the current gas price.
	GetGasPrice(ctx context.Context) (GasPrice, error)
	// DenoteInUSD converts the gas price to be in units of USD. Input prices should not be nil.
	DenoteInUSD(p GasPrice, wrappedNativePrice *big.Int) (GasPrice, error)
	// Median finds the median gas price in slice. If gas price has multiple components, median of each individual component should be taken. Input prices should not contain nil.
	Median(gasPrices []GasPrice) (GasPrice, error)
	// String converts the gas price to string.
	String(p GasPrice) string
}

// GasPriceEstimatorCommit provides gasPriceEstimatorCommon + features needed in commit plugin, e.g. price deviation check.
//
//go:generate mockery --quiet --name GasPriceEstimatorCommit --output . --filename gas_price_estimator_commit_mock.go --inpackage --case=underscore
type GasPriceEstimatorCommit interface {
	gasPriceEstimatorCommon
	// Deviates checks if p1 gas price diffs from p2 by deviation options. Input prices should not be nil.
	Deviates(p1 GasPrice, p2 GasPrice) (bool, error)
}

// GasPriceEstimatorExec provides gasPriceEstimatorCommon + features needed in exec plugin, e.g. message cost estimation.
//
//go:generate mockery --quiet --name GasPriceEstimatorExec --output . --filename gas_price_estimator_exec_mock.go --inpackage --case=underscore
type GasPriceEstimatorExec interface {
	gasPriceEstimatorCommon
	// EstimateMsgCostUSD estimates the costs for msg execution, and converts to USD value scaled by 1e18 (e.g. 5$ = 5e18).
	EstimateMsgCostUSD(p GasPrice, wrappedNativePrice *big.Int, msg internal.EVM2EVMOnRampCCIPSendRequestedWithMeta) (*big.Int, error)
}

// GasPriceEstimator provides complete gas price estimator functions.
//
//go:generate mockery --quiet --name GasPriceEstimator --output . --filename gas_price_estimator_mock.go --inpackage --case=underscore
type GasPriceEstimator interface {
	GasPriceEstimatorCommit
	GasPriceEstimatorExec
}

func NewGasPriceEstimatorForCommitPlugin(
	commitStoreVersion semver.Version,
	estimator gas.EvmFeeEstimator,
	maxExecGasPrice *big.Int,
	daDeviationPPB int64,
	execDeviationPPB int64,
) (GasPriceEstimatorCommit, error) {
	execEstimator := ExecGasPriceEstimator{
		estimator:    estimator,
		maxGasPrice:  maxExecGasPrice,
		deviationPPB: execDeviationPPB,
	}

	switch commitStoreVersion.String() {
	case "1.0.0", "1.1.0":
		return execEstimator, nil
	case "1.2.0":
		return DAGasPriceEstimator{
			execEstimator:       execEstimator,
			l1Oracle:            estimator.L1Oracle(),
			priceEncodingLength: daGasPriceEncodingLength,
			daDeviationPPB:      daDeviationPPB,
			daOverheadGas:       0,
			gasPerDAByte:        0,
			daMultiplier:        0,
		}, nil
	default:
		return nil, errors.Errorf("Invalid commitStore version: %s", commitStoreVersion)
	}
}
