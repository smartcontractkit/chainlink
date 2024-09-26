package ccipexec

import (
	"math"
	"math/big"
	"time"
)

const (
	EvmAddressLengthBytes           = 20
	EvmWordBytes                    = 32
	CalldataGasPerByte              = 16
	TokenAdminRegistryWarmupCost    = 2_500
	TokenAdminRegistryPoolLookupGas = 100 + // WARM_ACCESS_COST TokenAdminRegistry
		700 + // CALL cost for TokenAdminRegistry
		2_100 // COLD_SLOAD_COST loading the pool address
	SupportsInterfaceCheck = 2600 + // because the receiver will be untouched initially
		30_000*3 // supportsInterface of ERC165Checker library performs 3 static-calls of 30k gas each
	PerTokenOverheadGas = TokenAdminRegistryPoolLookupGas +
		SupportsInterfaceCheck +
		200_000 + // releaseOrMint using callWithExactGas
		50_000 // transfer using callWithExactGas
	RateLimiterOverheadGas = 2_100 + // COLD_SLOAD_COST for accessing token bucket
		5_000 // SSTORE_RESET_GAS for updating & decreasing token bucket
	ConstantMessagePartBytes            = 10 * 32 // A message consists of 10 abi encoded fields 32B each (after encoding)
	ExecutionStateProcessingOverheadGas = 2_100 + // COLD_SLOAD_COST for first reading the state
		20_000 + // SSTORE_SET_GAS for writing from 0 (untouched) to non-zero (in-progress)
		100 //# SLOAD_GAS = WARM_STORAGE_READ_COST for rewriting from non-zero (in-progress) to non-zero (success/failure)
)

// return the size of bytes for msg tokens
func bytesForMsgTokens(numTokens int) int {
	// token address (address) + token amount (uint256)
	return (EvmAddressLengthBytes + EvmWordBytes) * numTokens
}

// Offchain: we compute the max overhead gas to determine msg executability.
func overheadGas(dataLength, numTokens int) uint64 {
	messageBytes := ConstantMessagePartBytes +
		bytesForMsgTokens(numTokens) +
		dataLength

	messageCallDataGas := uint64(messageBytes * CalldataGasPerByte)

	// Rate limiter only limits value in tokens. It's not called if there are no
	// tokens in the message. The same goes for the admin registry, it's only loaded
	// if there are tokens, and it's only loaded once.
	rateLimiterOverhead := uint64(0)
	adminRegistryOverhead := uint64(0)
	if numTokens >= 1 {
		rateLimiterOverhead = RateLimiterOverheadGas
		adminRegistryOverhead = TokenAdminRegistryWarmupCost
	}

	return messageCallDataGas +
		ExecutionStateProcessingOverheadGas +
		SupportsInterfaceCheck +
		adminRegistryOverhead +
		rateLimiterOverhead +
		PerTokenOverheadGas*uint64(numTokens)
}

func maxGasOverHeadGas(numMsgs, dataLength, numTokens int) uint64 {
	merkleProofBytes := (math.Ceil(math.Log2(float64(numMsgs))))*32 + (1+2)*32 // only ever one outer root hash
	merkleGasShare := uint64(merkleProofBytes * CalldataGasPerByte)

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
