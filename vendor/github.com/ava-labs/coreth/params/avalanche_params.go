// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package params

import (
	"math/big"

	"github.com/ava-labs/avalanchego/utils/units"
)

// Minimum Gas Price
const (
	// MinGasPrice is the number of nAVAX required per gas unit for a
	// transaction to be valid, measured in wei
	LaunchMinGasPrice        int64 = 470_000_000_000
	ApricotPhase1MinGasPrice int64 = 225_000_000_000

	AvalancheAtomicTxFee = units.MilliAvax

	ApricotPhase1GasLimit uint64 = 8_000_000

	ApricotPhase3ExtraDataSize            uint64 = 80
	ApricotPhase3MinBaseFee               int64  = 75_000_000_000
	ApricotPhase3MaxBaseFee               int64  = 225_000_000_000
	ApricotPhase3InitialBaseFee           int64  = 225_000_000_000
	ApricotPhase3TargetGas                uint64 = 10_000_000
	ApricotPhase4MinBaseFee               int64  = 25_000_000_000
	ApricotPhase4MaxBaseFee               int64  = 1_000_000_000_000
	ApricotPhase4BaseFeeChangeDenominator uint64 = 12
	ApricotPhase5TargetGas                uint64 = 15_000_000
	ApricotPhase5BaseFeeChangeDenominator uint64 = 36

	// The base cost to charge per atomic transaction. Added in Apricot Phase 5.
	AtomicTxBaseCost uint64 = 10_000
)

// Constants for message sizes
const (
	MaxCodeHashesPerRequest = 5
)

var (
	// The atomic gas limit specifies the maximum amount of gas that can be consumed by the atomic
	// transactions included in a block and is enforced as of ApricotPhase5. Prior to ApricotPhase5,
	// a block included a single atomic transaction. As of ApricotPhase5, each block can include a set
	// of atomic transactions where the cumulative atomic gas consumed is capped by the atomic gas limit,
	// similar to the block gas limit.
	//
	// This value must always remain <= MaxUint64.
	AtomicGasLimit *big.Int = big.NewInt(100_000)
)
