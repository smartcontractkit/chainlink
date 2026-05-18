package prices

import (
	"math/big"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
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

// GasPriceEstimatorCommit provides gasPriceEstimatorCommon + features needed in commit plugin, e.g. price deviation check.
type GasPriceEstimatorCommit interface {
	cciptypes.GasPriceEstimatorCommit
}

// GasPriceEstimatorExec provides gasPriceEstimatorCommon + features needed in exec plugin, e.g. message cost estimation.
type GasPriceEstimatorExec interface {
	cciptypes.GasPriceEstimatorExec
}

// GasPriceEstimator provides complete gas price estimator functions.
type GasPriceEstimator interface {
	cciptypes.GasPriceEstimator
}

func NewGasPriceEstimatorForCommitPlugin(
	commitStoreVersion semver.Version,
	estimator gas.EvmFeeEstimator,
	maxExecGasPrice *big.Int,
	daDeviationPPB int64,
	execDeviationPPB int64,
	feeEstimatorConfig ccipdata.FeeEstimatorConfigReader,
) (GasPriceEstimatorCommit, error) {
	switch commitStoreVersion.String() {
	case "1.0.0", "1.1.0":
		return NewExecGasPriceEstimator(estimator, maxExecGasPrice, execDeviationPPB), nil
	case "1.2.0":
		return NewDAGasPriceEstimator(estimator, maxExecGasPrice, execDeviationPPB, daDeviationPPB, feeEstimatorConfig), nil
	default:
		return nil, errors.Errorf("Invalid commitStore version: %s", commitStoreVersion)
	}
}
