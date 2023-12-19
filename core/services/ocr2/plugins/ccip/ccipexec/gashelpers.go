package ccipexec

import (
	"math"
	"math/big"
	"time"
)

const (
	EVM_ADDRESS_LENGTH_BYTES = 20
	EVM_WORD_BYTES           = 32
	CALLDATA_GAS_PER_BYTE    = 16
	PER_TOKEN_OVERHEAD_GAS   = 2_100 + // COLD_SLOAD_COST for first reading the pool
		2_100 + // COLD_SLOAD_COST for pool to ensure allowed offramp calls it
		2_100 + // COLD_SLOAD_COST for accessing pool balance slot
		5_000 + // SSTORE_RESET_GAS for decreasing pool balance from non-zero to non-zero
		2_100 + // COLD_SLOAD_COST for accessing receiver balance
		20_000 + // SSTORE_SET_GAS for increasing receiver balance from zero to non-zero
		2_100 // COLD_SLOAD_COST for obtanining price of token to use for aggregate token bucket
	RATE_LIMITER_OVERHEAD_GAS = 2_100 + // COLD_SLOAD_COST for accessing token bucket
		5_000 // SSTORE_RESET_GAS for updating & decreasing token bucket
	EXTERNAL_CALL_OVERHEAD_GAS = 2600 + // because the receiver will be untouched initially
		30_000*3 // supportsInterface of ERC165Checker library performs 3 static-calls of 30k gas each
	FEE_BOOSTING_OVERHEAD_GAS               = 200_000
	CONSTANT_MESSAGE_PART_BYTES             = 10 * 32 // A message consists of 10 abi encoded fields 32B each (after encoding)
	EXECUTION_STATE_PROCESSING_OVERHEAD_GAS = 2_100 + // COLD_SLOAD_COST for first reading the state
		20_000 + // SSTORE_SET_GAS for writing from 0 (untouched) to non-zero (in-progress)
		100 //# SLOAD_GAS = WARM_STORAGE_READ_COST for rewriting from non-zero (in-progress) to non-zero (success/failure)
	EVM_MESSAGE_FIXED_BYTES     = 448 // Byte size of fixed-size fields in EVM2EVMMessage
	EVM_MESSAGE_BYTES_PER_TOKEN = 128 // Byte size of each token transfer, consisting of 1 EVMTokenAmount and 1 bytes, excl length of bytes
	DA_MULTIPLIER_BASE          = int64(10000)
)

// return the size of bytes for msg tokens
func bytesForMsgTokens(numTokens int) int {
	// token address (address) + token amount (uint256)
	return (EVM_ADDRESS_LENGTH_BYTES + EVM_WORD_BYTES) * numTokens
}

// Offchain: we compute the max overhead gas to determine msg executability.
func overheadGas(dataLength, numTokens int) uint64 {
	messageBytes := CONSTANT_MESSAGE_PART_BYTES +
		bytesForMsgTokens(numTokens) +
		dataLength

	messageCallDataGas := uint64(messageBytes * CALLDATA_GAS_PER_BYTE)

	// Rate limiter only limits value in tokens. It's not called if there are no
	// tokens in the message.
	rateLimiterOverhead := uint64(0)
	if numTokens >= 1 {
		rateLimiterOverhead = RATE_LIMITER_OVERHEAD_GAS
	}

	return messageCallDataGas +
		EXECUTION_STATE_PROCESSING_OVERHEAD_GAS +
		PER_TOKEN_OVERHEAD_GAS*uint64(numTokens) +
		rateLimiterOverhead +
		EXTERNAL_CALL_OVERHEAD_GAS
}

func maxGasOverHeadGas(numMsgs, dataLength, numTokens int) uint64 {
	merkleProofBytes := (math.Ceil(math.Log2(float64(numMsgs))))*32 + (1+2)*32 // only ever one outer root hash
	merkleGasShare := uint64(merkleProofBytes * CALLDATA_GAS_PER_BYTE)

	return overheadGas(dataLength, numTokens) + merkleGasShare
}

// waitBoostedFee boosts the given fee according to the time passed since the msg was sent.
// RelativeBoostPerWaitHour is used to normalize the time diff,
// it makes our loss taking "smooth" and gives us time to react without a hard deadline.
// At the same time, messages that are slightly underpaid will start going through after waiting for a little bit.
//
// wait_boosted_fee(m) = (1 + (now - m.send_time).hours * RELATIVE_BOOST_PER_WAIT_HOUR) * fee(m)
func waitBoostedFee(waitTime time.Duration, fee *big.Int, relativeBoostPerWaitHour float64) *big.Int {
	k := 1.0 + waitTime.Hours()*relativeBoostPerWaitHour

	boostedFee := big.NewFloat(0).Mul(big.NewFloat(k), new(big.Float).SetInt(fee))
	res, _ := boostedFee.Int(nil)

	return res
}
