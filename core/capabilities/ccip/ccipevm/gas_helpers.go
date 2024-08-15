package ccipevm

import (
	"math"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
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

func NewGasEstimateProvider() EstimateProvider {
	return EstimateProvider{}
}

type EstimateProvider struct {
}

// CalculateMerkleTreeGas estimates the merkle tree gas based on number of requests
func (gp EstimateProvider) CalculateMerkleTreeGas(numRequests int) uint64 {
	if numRequests == 0 {
		return 0
	}
	merkleProofBytes := (math.Ceil(math.Log2(float64(numRequests))))*32 + (1+2)*32 // only ever one outer root hash
	return uint64(merkleProofBytes * CalldataGasPerByte)
}

// return the size of bytes for msg tokens
func bytesForMsgTokens(numTokens int) int {
	// token address (address) + token amount (uint256)
	return (EvmAddressLengthBytes + EvmWordBytes) * numTokens
}

// CalculateMessageMaxGas computes the maximum gas overhead for a message.
func (gp EstimateProvider) CalculateMessageMaxGas(msg cciptypes.Message) uint64 {
	numTokens := len(msg.TokenAmounts)
	var data []byte = msg.Data
	dataLength := len(data)

	// TODO: update interface to return error?
	// Although this decoding should never fail.
	messageGasLimit, err := decodeExtraArgsV1V2(msg.ExtraArgs)
	if err != nil {
		panic(err)
	}

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

	return messageGasLimit.Uint64() +
		messageCallDataGas +
		ExecutionStateProcessingOverheadGas +
		SupportsInterfaceCheck +
		adminRegistryOverhead +
		rateLimiterOverhead +
		PerTokenOverheadGas*uint64(numTokens)
}
