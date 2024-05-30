// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package evm_2_evm_multi_onramp

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated"
)

var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

type ClientEVM2AnyMessage struct {
	Receiver     []byte
	Data         []byte
	TokenAmounts []ClientEVMTokenAmount
	FeeToken     common.Address
	ExtraArgs    []byte
}

type ClientEVMTokenAmount struct {
	Token  common.Address
	Amount *big.Int
}

type EVM2EVMMultiOnRampDestChainConfig struct {
	DynamicConfig  EVM2EVMMultiOnRampDestChainDynamicConfig
	PrevOnRamp     common.Address
	SequenceNumber uint64
	MetadataHash   [32]byte
}

type EVM2EVMMultiOnRampDestChainConfigArgs struct {
	DestChainSelector uint64
	DynamicConfig     EVM2EVMMultiOnRampDestChainDynamicConfig
	PrevOnRamp        common.Address
}

type EVM2EVMMultiOnRampDestChainDynamicConfig struct {
	IsEnabled                         bool
	MaxNumberOfTokensPerMsg           uint16
	MaxDataBytes                      uint32
	MaxPerMsgGasLimit                 uint32
	DestGasOverhead                   uint32
	DestGasPerPayloadByte             uint16
	DestDataAvailabilityOverheadGas   uint32
	DestGasPerDataAvailabilityByte    uint16
	DestDataAvailabilityMultiplierBps uint16
	DefaultTokenFeeUSDCents           uint16
	DefaultTokenDestGasOverhead       uint32
	DefaultTokenDestBytesOverhead     uint32
	DefaultTxGasLimit                 uint64
	GasMultiplierWeiPerEth            uint64
	NetworkFeeUSDCents                uint32
}

type EVM2EVMMultiOnRampDynamicConfig struct {
	Router             common.Address
	PriceRegistry      common.Address
	TokenAdminRegistry common.Address
}

type EVM2EVMMultiOnRampNopAndWeight struct {
	Nop    common.Address
	Weight uint16
}

type EVM2EVMMultiOnRampPremiumMultiplierWeiPerEthArgs struct {
	Token                      common.Address
	PremiumMultiplierWeiPerEth uint64
}

type EVM2EVMMultiOnRampStaticConfig struct {
	LinkToken       common.Address
	ChainSelector   uint64
	MaxNopFeesJuels *big.Int
	RmnProxy        common.Address
}

type EVM2EVMMultiOnRampTokenTransferFeeConfig struct {
	MinFeeUSDCents            uint32
	MaxFeeUSDCents            uint32
	DeciBps                   uint16
	DestGasOverhead           uint32
	DestBytesOverhead         uint32
	AggregateRateLimitEnabled bool
	IsEnabled                 bool
}

type EVM2EVMMultiOnRampTokenTransferFeeConfigArgs struct {
	Token                     common.Address
	MinFeeUSDCents            uint32
	MaxFeeUSDCents            uint32
	DeciBps                   uint16
	DestGasOverhead           uint32
	DestBytesOverhead         uint32
	AggregateRateLimitEnabled bool
}

type InternalEVM2EVMMessage struct {
	SourceChainSelector uint64
	Sender              common.Address
	Receiver            common.Address
	SequenceNumber      uint64
	GasLimit            *big.Int
	Strict              bool
	Nonce               uint64
	FeeToken            common.Address
	FeeTokenAmount      *big.Int
	Data                []byte
	TokenAmounts        []ClientEVMTokenAmount
	SourceTokenData     [][]byte
	MessageId           [32]byte
}

type RateLimiterConfig struct {
	IsEnabled bool
	Capacity  *big.Int
	Rate      *big.Int
}

type RateLimiterTokenBucket struct {
	Tokens      *big.Int
	LastUpdated uint32
	IsEnabled   bool
	Capacity    *big.Int
	Rate        *big.Int
}

var EVM2EVMMultiOnRampMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"linkToken\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint96\",\"name\":\"maxNopFeesJuels\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"}],\"internalType\":\"structEVM2EVMMultiOnRamp.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"priceRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAdminRegistry\",\"type\":\"address\"}],\"internalType\":\"structEVM2EVMMultiOnRamp.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint16\",\"name\":\"maxNumberOfTokensPerMsg\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxDataBytes\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerMsgGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerPayloadByte\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destDataAvailabilityOverheadGas\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerDataAvailabilityByte\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"destDataAvailabilityMultiplierBps\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"defaultTokenFeeUSDCents\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"defaultTokenDestGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"defaultTokenDestBytesOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"defaultTxGasLimit\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"gasMultiplierWeiPerEth\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"networkFeeUSDCents\",\"type\":\"uint32\"}],\"internalType\":\"structEVM2EVMMultiOnRamp.DestChainDynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"prevOnRamp\",\"type\":\"address\"}],\"internalType\":\"structEVM2EVMMultiOnRamp.DestChainConfigArgs[]\",\"name\":\"destChainConfigArgs\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"rateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"premiumMultiplierWeiPerEth\",\"type\":\"uint64\"}],\"internalType\":\"structEVM2EVMMultiOnRamp.PremiumMultiplierWeiPerEthArgs[]\",\"name\":\"premiumMultiplierWeiPerEthArgs\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"minFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"deciBps\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"destBytesOverhead\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"aggregateRateLimitEnabled\",\"type\":\"bool\"}],\"internalType\":\"structEVM2EVMMultiOnRamp.TokenTransferFeeConfigArgs[]\",\"name\":\"tokenTransferFeeConfigArgs\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"nop\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"weight\",\"type\":\"uint16\"}],\"internalType\":\"structEVM2EVMMultiOnRamp.NopAndWeight[]\",\"name\":\"nopsAndWeights\",\"type\":\"tuple[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"}],\"name\":\"AggregateValueMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"}],\"name\":\"AggregateValueRateLimitReached\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BucketOverfilled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotSendZeroTokens\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"CursedByRMN\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"}],\"name\":\"DestinationChainNotEnabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GetSupportedTokensFunctionalityRemovedCheckAdminRegistry\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedAddress\",\"type\":\"bytes\"}],\"name\":\"InvalidAddress\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"}],\"name\":\"InvalidChainSelector\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidConfig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"destBytesOverhead\",\"type\":\"uint32\"}],\"name\":\"InvalidDestBytesOverhead\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"}],\"name\":\"InvalidDestChainConfig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedAddress\",\"type\":\"bytes\"}],\"name\":\"InvalidEVMAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidExtraArgsTag\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"nop\",\"type\":\"address\"}],\"name\":\"InvalidNopAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidWithdrawParams\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LinkBalanceNotSettled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxFeeBalanceReached\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MessageGasLimitTooHigh\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"maxSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actualSize\",\"type\":\"uint256\"}],\"name\":\"MessageTooLarge\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeCalledByRouter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoFeesToPay\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoNopsToPay\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"NotAFeeToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByAdminOrOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAdminOrNop\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"PriceNotFoundForToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RouterMustSetOriginalSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SourceTokenDataTooLarge\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenRateLimitReached\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyNops\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnsupportedNumberOfTokens\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"UnsupportedToken\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"strict\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"feeTokenAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMTokenAmount[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes[]\",\"name\":\"sourceTokenData\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"}],\"indexed\":false,\"internalType\":\"structInternal.EVM2EVMMessage\",\"name\":\"message\",\"type\":\"tuple\"}],\"name\":\"CCIPSendRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"ConfigChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"linkToken\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint96\",\"name\":\"maxNopFeesJuels\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"}],\"indexed\":false,\"internalType\":\"structEVM2EVMMultiOnRamp.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"priceRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAdminRegistry\",\"type\":\"address\"}],\"indexed\":false,\"internalType\":\"structEVM2EVMMultiOnRamp.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint16\",\"name\":\"maxNumberOfTokensPerMsg\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxDataBytes\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerMsgGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerPayloadByte\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destDataAvailabilityOverheadGas\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerDataAvailabilityByte\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"destDataAvailabilityMultiplierBps\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"defaultTokenFeeUSDCents\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"defaultTokenDestGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"defaultTokenDestBytesOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"defaultTxGasLimit\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"gasMultiplierWeiPerEth\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"networkFeeUSDCents\",\"type\":\"uint32\"}],\"internalType\":\"structEVM2EVMMultiOnRamp.DestChainDynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"prevOnRamp\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"metadataHash\",\"type\":\"bytes32\"}],\"indexed\":false,\"internalType\":\"structEVM2EVMMultiOnRamp.DestChainConfig\",\"name\":\"destChainConfig\",\"type\":\"tuple\"}],\"name\":\"DestChainAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint16\",\"name\":\"maxNumberOfTokensPerMsg\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxDataBytes\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerMsgGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerPayloadByte\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destDataAvailabilityOverheadGas\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerDataAvailabilityByte\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"destDataAvailabilityMultiplierBps\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"defaultTokenFeeUSDCents\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"defaultTokenDestGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"defaultTokenDestBytesOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"defaultTxGasLimit\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"gasMultiplierWeiPerEth\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"networkFeeUSDCents\",\"type\":\"uint32\"}],\"indexed\":false,\"internalType\":\"structEVM2EVMMultiOnRamp.DestChainDynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"name\":\"DestChainDynamicConfigUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"nop\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"NopPaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nopWeightsTotal\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"nop\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"weight\",\"type\":\"uint16\"}],\"indexed\":false,\"internalType\":\"structEVM2EVMMultiOnRamp.NopAndWeight[]\",\"name\":\"nopsAndWeights\",\"type\":\"tuple[]\"}],\"name\":\"NopsSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"premiumMultiplierWeiPerEth\",\"type\":\"uint64\"}],\"name\":\"PremiumMultiplierWeiPerEthUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"tokens\",\"type\":\"address[]\"}],\"name\":\"TokenTransferFeeConfigDeleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"minFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"deciBps\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"destBytesOverhead\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"aggregateRateLimitEnabled\",\"type\":\"bool\"}],\"indexed\":false,\"internalType\":\"structEVM2EVMMultiOnRamp.TokenTransferFeeConfigArgs[]\",\"name\":\"transferFeeConfig\",\"type\":\"tuple[]\"}],\"name\":\"TokenTransferFeeConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"TokensConsumed\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint16\",\"name\":\"maxNumberOfTokensPerMsg\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxDataBytes\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerMsgGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerPayloadByte\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destDataAvailabilityOverheadGas\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerDataAvailabilityByte\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"destDataAvailabilityMultiplierBps\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"defaultTokenFeeUSDCents\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"defaultTokenDestGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"defaultTokenDestBytesOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"defaultTxGasLimit\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"gasMultiplierWeiPerEth\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"networkFeeUSDCents\",\"type\":\"uint32\"}],\"internalType\":\"structEVM2EVMMultiOnRamp.DestChainDynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"prevOnRamp\",\"type\":\"address\"}],\"internalType\":\"structEVM2EVMMultiOnRamp.DestChainConfigArgs[]\",\"name\":\"destChainConfigArgs\",\"type\":\"tuple[]\"}],\"name\":\"applyDestChainConfigUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"premiumMultiplierWeiPerEth\",\"type\":\"uint64\"}],\"internalType\":\"structEVM2EVMMultiOnRamp.PremiumMultiplierWeiPerEthArgs[]\",\"name\":\"premiumMultiplierWeiPerEthArgs\",\"type\":\"tuple[]\"}],\"name\":\"applyPremiumMultiplierWeiPerEthUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMTokenAmount[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structClient.EVM2AnyMessage\",\"name\":\"message\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"feeTokenAmount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"originalSender\",\"type\":\"address\"}],\"name\":\"forwardFromRouter\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"}],\"name\":\"getDestChainConfig\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint16\",\"name\":\"maxNumberOfTokensPerMsg\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxDataBytes\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerMsgGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerPayloadByte\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destDataAvailabilityOverheadGas\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerDataAvailabilityByte\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"destDataAvailabilityMultiplierBps\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"defaultTokenFeeUSDCents\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"defaultTokenDestGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"defaultTokenDestBytesOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"defaultTxGasLimit\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"gasMultiplierWeiPerEth\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"networkFeeUSDCents\",\"type\":\"uint32\"}],\"internalType\":\"structEVM2EVMMultiOnRamp.DestChainDynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"prevOnRamp\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"metadataHash\",\"type\":\"bytes32\"}],\"internalType\":\"structEVM2EVMMultiOnRamp.DestChainConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDynamicConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"priceRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAdminRegistry\",\"type\":\"address\"}],\"internalType\":\"structEVM2EVMMultiOnRamp.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"}],\"name\":\"getExpectedNextSequenceNumber\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMTokenAmount[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structClient.EVM2AnyMessage\",\"name\":\"message\",\"type\":\"tuple\"}],\"name\":\"getFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"feeTokenAmount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNopFeesJuels\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNops\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"nop\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"weight\",\"type\":\"uint16\"}],\"internalType\":\"structEVM2EVMMultiOnRamp.NopAndWeight[]\",\"name\":\"nopsAndWeights\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"weightsTotal\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"contractIERC20\",\"name\":\"sourceToken\",\"type\":\"address\"}],\"name\":\"getPoolBySourceToken\",\"outputs\":[{\"internalType\":\"contractIPool\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getPremiumMultiplierWeiPerEth\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"premiumMultiplierWeiPerEth\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"getSenderNonce\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStaticConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"linkToken\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint96\",\"name\":\"maxNopFeesJuels\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"}],\"internalType\":\"structEVM2EVMMultiOnRamp.StaticConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"getSupportedTokens\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTokenLimitAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getTokenTransferFeeConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"minFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"deciBps\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"destBytesOverhead\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"aggregateRateLimitEnabled\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"}],\"internalType\":\"structEVM2EVMMultiOnRamp.TokenTransferFeeConfig\",\"name\":\"tokenTransferFeeConfig\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkAvailableForPayment\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"payNops\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"setAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"priceRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAdminRegistry\",\"type\":\"address\"}],\"internalType\":\"structEVM2EVMMultiOnRamp.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"name\":\"setDynamicConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"nop\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"weight\",\"type\":\"uint16\"}],\"internalType\":\"structEVM2EVMMultiOnRamp.NopAndWeight[]\",\"name\":\"nopsAndWeights\",\"type\":\"tuple[]\"}],\"name\":\"setNops\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"setRateLimiterConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"minFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"deciBps\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"destBytesOverhead\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"aggregateRateLimitEnabled\",\"type\":\"bool\"}],\"internalType\":\"structEVM2EVMMultiOnRamp.TokenTransferFeeConfigArgs[]\",\"name\":\"tokenTransferFeeConfigArgs\",\"type\":\"tuple[]\"},{\"internalType\":\"address[]\",\"name\":\"tokensToUseDefaultFeeConfigs\",\"type\":\"address[]\"}],\"name\":\"setTokenTransferFeeConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdrawNonLinkFees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6101006040523480156200001257600080fd5b506040516200874d3803806200874d833981016040819052620000359162001e36565b8333806000816200008d5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c057620000c0816200022f565b50506040805160a081018252602084810180516001600160801b039081168085524263ffffffff169385018490528751151585870181905292518216606086018190529790950151166080909301839052600380546001600160a01b031916909417600160801b9283021760ff60a01b1916600160a01b90910217909255029091176004555086516001600160a01b0316158062000169575060208701516001600160401b0316155b8062000180575060608701516001600160a01b0316155b156200019f576040516306b7c75960e31b815260040160405180910390fd5b86516001600160a01b0390811660a05260208801516001600160401b031660c05260408801516001600160601b031660805260608801511660e052620001e586620002da565b620001f0856200042c565b620001fb836200096c565b6040805160008152602081019091526200021790839062000a38565b620002228162000cd7565b505050505050506200251e565b336001600160a01b03821603620002895760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000084565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60208101516001600160a01b031662000306576040516306b7c75960e31b815260040160405180910390fd5b8051600580546001600160a01b039283166001600160a01b0319918216179091556020808401516006805491851691841691909117905560408085015160078054918616919094161790925581516080808201845260a0518516825260c0516001600160401b03169282019290925290516001600160601b03168183015260e0519092166060830152517f45d99abaa76ccb5c1a18f32b24b8a795ede9926a20056247319c426c1370fefa916200042191849082516001600160a01b0390811682526020808501516001600160401b0316818401526040808601516001600160601b031681850152606095860151831695840195909552835182166080840152830151811660a083015291909201511660c082015260e00190565b60405180910390a150565b60005b81518110156200096857600082828151811062000450576200045062001f30565b60200260200101519050600083838151811062000471576200047162001f30565b6020026020010151600001519050806001600160401b0316600003620004b65760405163c35aa79d60e01b81526001600160401b038216600482015260240162000084565b6000600b6000836001600160401b03166001600160401b0316815260200190815260200160002090506000836040015190506000604051806080016040528086602001518152602001836001600160a01b031681526020018460020160149054906101000a90046001600160401b03166001600160401b031681526020018460030154815250905080600001518360000160008201518160000160006101000a81548160ff02191690831515021790555060208201518160000160016101000a81548161ffff021916908361ffff16021790555060408201518160000160036101000a81548163ffffffff021916908363ffffffff16021790555060608201518160000160076101000a81548163ffffffff021916908363ffffffff160217905550608082015181600001600b6101000a81548163ffffffff021916908363ffffffff16021790555060a082015181600001600f6101000a81548161ffff021916908361ffff16021790555060c08201518160000160116101000a81548163ffffffff021916908363ffffffff16021790555060e08201518160000160156101000a81548161ffff021916908361ffff1602179055506101008201518160000160176101000a81548161ffff021916908361ffff1602179055506101208201518160000160196101000a81548161ffff021916908361ffff16021790555061014082015181600001601b6101000a81548163ffffffff021916908363ffffffff1602179055506101608201518160010160006101000a81548163ffffffff021916908363ffffffff1602179055506101808201518160010160046101000a8154816001600160401b0302191690836001600160401b031602179055506101a082015181600101600c6101000a8154816001600160401b0302191690836001600160401b031602179055506101c08201518160010160146101000a81548163ffffffff021916908363ffffffff16021790555090505082600301546000801b03620008855760c051604080517f8acd72527118c8324937b1a42e02cd246697c3b633f1742f3cae11de233722b360208201526001600160401b0392831691810191909152908516606082015230608082015260a00160408051601f1981840301815291905280516020909101206060820181905260038401556001600160a01b038216156200083c576002830180546001600160a01b0319166001600160a01b0384161790555b836001600160401b03167f7a70081ee29c1fc27898089ba2a5fc35ac0106b043c82ccecd24c6fd48f6ca868460405162000877919062001f46565b60405180910390a262000957565b60028301546001600160a01b03838116911614620008c25760405163c35aa79d60e01b81526001600160401b038516600482015260240162000084565b60208560200151610160015163ffffffff1610156200090f57602085015161016001516040516312766e0160e11b81526000600482015263ffffffff909116602482015260440162000084565b836001600160401b03167f944eb884a589931130671ee4a7379fbe5fe65ed605a048ba99c454582f2460b086602001516040516200094e9190620020da565b60405180910390a25b50505050508060010190506200042f565b5050565b60005b81518110156200096857600082828151811062000990576200099062001f30565b60200260200101516000015190506000838381518110620009b557620009b562001f30565b6020908102919091018101518101516001600160a01b0384166000818152600c845260409081902080546001600160401b0319166001600160401b0385169081179091559051908152919350917fbb77da6f7210cdd16904228a9360133d1d7dfff99b1bc75f128da5b53e28f97d910160405180910390a250506001016200096f565b60005b825181101562000bf757600083828151811062000a5c5762000a5c62001f30565b6020026020010151905060208160a0015163ffffffff16101562000ab257805160a08201516040516312766e0160e11b81526001600160a01b03909216600483015263ffffffff16602482015260440162000084565b6040805160e08101825260208381015163ffffffff908116835284840151811682840190815260608087015161ffff9081168688019081526080808a0151861693880193845260a0808b0151871691890191825260c0808c01511515918a019182526001908a018181529b516001600160a01b03166000908152600d9099529990972097518854955192519451915197519a519087166001600160401b031990961695909517640100000000928716929092029190911765ffffffffffff60401b191668010000000000000000939092169290920263ffffffff60501b1916176a0100000000000000000000918416919091021764ffffffffff60701b1916600160701b939092169290920260ff60901b191617600160901b941515949094029390931760ff60981b1916600160981b93151593909302929092179091550162000a3b565b507ff5791bc457b3bb990493cf5f655db46c25ccf5764c9b99b8969b4c72ea7df9d08260405162000c2991906200220d565b60405180910390a160005b815181101562000c9157600d600083838151811062000c575762000c5762001f30565b6020908102919091018101516001600160a01b0316825281019190915260400160002080546001600160a01b031916905560010162000c34565b5080511562000968577ffb95a0042158e60a33e7b5bec100f3d95407b1a71bee6633bd54b8887449750b8160405162000ccb9190620022af565b60405180910390a15050565b8051604081111562000cfc57604051635ad0867d60e11b815260040160405180910390fd5b600f546c01000000000000000000000000900463ffffffff161580159062000d465750600f5463ffffffff6c010000000000000000000000008204166001600160601b0390911610155b1562000d565762000d5662000ef9565b600062000d646008620010f1565b90505b801562000db057600062000d8a62000d8160018462002314565b60089062001104565b50905062000d9a60088262001122565b50508062000da8906200232a565b905062000d67565b506000805b8281101562000e9057600084828151811062000dd55762000dd562001f30565b6020026020010151600001519050600085838151811062000dfa5762000dfa62001f30565b602002602001015160200151905060a0516001600160a01b0316826001600160a01b0316148062000e3257506001600160a01b038216155b1562000e5d57604051634de938d160e01b81526001600160a01b038316600482015260240162000084565b62000e6f60088361ffff841662001140565b5062000e8061ffff82168562002344565b9350505080600101905062000db5565b50600f805463ffffffff60601b19166c0100000000000000000000000063ffffffff8416021790556040517f8c337bff38141c507abd25c547606bdde78fe8c12e941ab613f3a565fea6cd249062000eec908390869062002364565b60405180910390a1505050565b6000546001600160a01b0316331480159062000f2057506002546001600160a01b03163314155b801562000f37575062000f3560083362001160565b155b1562000f565760405163032bb72b60e31b815260040160405180910390fd5b600f546c01000000000000000000000000900463ffffffff16600081900362000f925760405163990e30bf60e01b815260040160405180910390fd5b600f546001600160601b03168181101562000fc0576040516311a1ee3b60e31b815260040160405180910390fd5b600062000fcc62001177565b121562000fec57604051631e9acf1760e31b815260040160405180910390fd5b80600062000ffb6008620010f1565b905060005b81811015620010cb576000806200101960088462001104565b909250905060008762001036836001600160601b038a16620023d4565b620010429190620023ee565b905062001050818762002411565b60a05190965062001075906001600160a01b0316846001600160601b03841662001205565b6040516001600160601b03821681526001600160a01b038416907f55fdec2aab60a41fa5abb106670eb1006f5aeaee1ba7afea2bc89b5b3ec7678f9060200160405180910390a250505080600101905062001000565b5050600f80546001600160601b0319166001600160601b03929092169190911790555050565b6000620010fe8262001262565b92915050565b60008080806200111586866200126f565b9097909650945050505050565b600062001139836001600160a01b0384166200129c565b9392505050565b600062001158846001600160a01b03851684620012bb565b949350505050565b600062001139836001600160a01b038416620012da565b600f5460a0516040516370a0823160e01b81523060048201526000926001600160601b0316916001600160a01b0316906370a0823190602401602060405180830381865afa158015620011ce573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620011f4919062002434565b6200120091906200244e565b905090565b604080516001600160a01b038416602482015260448082018490528251808303909101815260649091019091526020810180516001600160e01b0390811663a9059cbb60e01b179091526200125d918591620012e816565b505050565b6000620010fe82620013b9565b600080806200127f8585620013c4565b600081815260029690960160205260409095205494959350505050565b60008181526002830160205260408120819055620011398383620013d2565b60008281526002840160205260408120829055620011588484620013e0565b6000620011398383620013ee565b6040805180820190915260208082527f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65649082015260009062001337906001600160a01b03851690849062001407565b8051909150156200125d578080602001905181019062001358919062002471565b6200125d5760405162461bcd60e51b815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e6044820152691bdd081cdd58d8d9595960b21b606482015260840162000084565b6000620010fe825490565b600062001139838362001418565b600062001139838362001445565b600062001139838362001550565b6000818152600183016020526040812054151562001139565b6060620011588484600085620015a2565b600082600001828154811062001432576200143262001f30565b9060005260206000200154905092915050565b600081815260018301602052604081205480156200153e5760006200146c60018362002314565b8554909150600090620014829060019062002314565b9050818114620014ee576000866000018281548110620014a657620014a662001f30565b9060005260206000200154905080876000018481548110620014cc57620014cc62001f30565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806200150257620015026200248f565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050620010fe565b6000915050620010fe565b5092915050565b60008181526001830160205260408120546200159957508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155620010fe565b506000620010fe565b606082471015620016055760405162461bcd60e51b815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f6044820152651c8818d85b1b60d21b606482015260840162000084565b600080866001600160a01b03168587604051620016239190620024cb565b60006040518083038185875af1925050503d806000811462001662576040519150601f19603f3d011682016040523d82523d6000602084013e62001667565b606091505b5090925090506200167b8783838762001686565b979650505050505050565b60608315620016fa578251600003620016f2576001600160a01b0385163b620016f25760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000604482015260640162000084565b508162001158565b620011588383815115620017115781518083602001fd5b8060405162461bcd60e51b8152600401620000849190620024e9565b634e487b7160e01b600052604160045260246000fd5b604051606081016001600160401b03811182821017156200176857620017686200172d565b60405290565b6040516101e081016001600160401b03811182821017156200176857620017686200172d565b604080519081016001600160401b03811182821017156200176857620017686200172d565b60405160e081016001600160401b03811182821017156200176857620017686200172d565b604051601f8201601f191681016001600160401b03811182821017156200180957620018096200172d565b604052919050565b80516001600160a01b03811681146200182957600080fd5b919050565b80516001600160401b03811681146200182957600080fd5b6000608082840312156200185957600080fd5b604051608081016001600160401b03811182821017156200187e576200187e6200172d565b6040529050806200188f8362001811565b81526200189f602084016200182e565b602082015260408301516001600160601b0381168114620018bf57600080fd5b6040820152620018d26060840162001811565b60608201525092915050565b600060608284031215620018f157600080fd5b620018fb62001743565b9050620019088262001811565b8152620019186020830162001811565b60208201526200192b6040830162001811565b604082015292915050565b60006001600160401b038211156200195257620019526200172d565b5060051b60200190565b805180151581146200182957600080fd5b805161ffff811681146200182957600080fd5b805163ffffffff811681146200182957600080fd5b600082601f830112620019a757600080fd5b81516020620019c0620019ba8362001936565b620017de565b8281526102209283028501820192828201919087851115620019e157600080fd5b8387015b8581101562001b8b578089038281121562001a005760008081fd5b62001a0a62001743565b62001a15836200182e565b81526101e080601f198401121562001a2d5760008081fd5b62001a376200176e565b925062001a468885016200195c565b8352604062001a578186016200196d565b89850152606062001a6a81870162001980565b82860152608062001a7d81880162001980565b8287015260a0915062001a9282880162001980565b9086015260c062001aa58782016200196d565b8287015260e0915062001aba82880162001980565b9086015261010062001ace8782016200196d565b82870152610120915062001ae48288016200196d565b9086015261014062001af88782016200196d565b82870152610160915062001b0e82880162001980565b9086015261018062001b2287820162001980565b828701526101a0915062001b388288016200182e565b908601526101c062001b4c8782016200182e565b8287015262001b5d84880162001980565b818701525050838984015262001b77610200860162001811565b9083015250855250928401928101620019e5565b5090979650505050505050565b80516001600160801b03811681146200182957600080fd5b60006060828403121562001bc357600080fd5b62001bcd62001743565b905062001bda826200195c565b815262001bea6020830162001b98565b60208201526200192b6040830162001b98565b600082601f83011262001c0f57600080fd5b8151602062001c22620019ba8362001936565b82815260069290921b8401810191818101908684111562001c4257600080fd5b8286015b8481101562001c98576040818903121562001c615760008081fd5b62001c6b62001794565b62001c768262001811565b815262001c858583016200182e565b8186015283529183019160400162001c46565b509695505050505050565b600082601f83011262001cb557600080fd5b8151602062001cc8620019ba8362001936565b82815260e0928302850182019282820191908785111562001ce857600080fd5b8387015b8581101562001b8b5781818a03121562001d065760008081fd5b62001d10620017b9565b62001d1b8262001811565b815262001d2a86830162001980565b86820152604062001d3d81840162001980565b90820152606062001d508382016200196d565b90820152608062001d6383820162001980565b9082015260a062001d7683820162001980565b9082015260c062001d898382016200195c565b90820152845292840192810162001cec565b600082601f83011262001dad57600080fd5b8151602062001dc0620019ba8362001936565b82815260069290921b8401810191818101908684111562001de057600080fd5b8286015b8481101562001c98576040818903121562001dff5760008081fd5b62001e0962001794565b62001e148262001811565b815262001e238583016200196d565b8186015283529183019160400162001de4565b60008060008060008060006101c0888a03121562001e5357600080fd5b62001e5f898962001846565b965062001e708960808a01620018de565b60e08901519096506001600160401b038082111562001e8e57600080fd5b62001e9c8b838c0162001995565b965062001eae8b6101008c0162001bb0565b95506101608a015191508082111562001ec657600080fd5b62001ed48b838c0162001bfd565b94506101808a015191508082111562001eec57600080fd5b62001efa8b838c0162001ca3565b93506101a08a015191508082111562001f1257600080fd5b5062001f218a828b0162001d9b565b91505092959891949750929550565b634e487b7160e01b600052603260045260246000fd5b815460ff81161515825261024082019061ffff600882901c8116602085015263ffffffff601883901c8116604086015262001f8e60608601828560381c1663ffffffff169052565b62001fa660808601828560581c1663ffffffff169052565b62001fbc60a08601838560781c1661ffff169052565b62001fd460c08601828560881c1663ffffffff169052565b62001fea60e08601838560a81c1661ffff169052565b620020016101008601838560b81c1661ffff169052565b620020186101208601838560c81c1661ffff169052565b620020316101408601828560d81c1663ffffffff169052565b600186015463ffffffff8282161661016087015292506001600160401b03602084901c811661018087015291506200207a6101a08601838560601c166001600160401b03169052565b620020936101c08601828560a01c1663ffffffff169052565b5060028501546001600160a01b0381166101e08601529150620020c76102008501828460a01c166001600160401b03169052565b5050600383015461022083015292915050565b8151151581526101e081016020830151620020fb602084018261ffff169052565b50604083015162002114604084018263ffffffff169052565b5060608301516200212d606084018263ffffffff169052565b50608083015162002146608084018263ffffffff169052565b5060a08301516200215d60a084018261ffff169052565b5060c08301516200217660c084018263ffffffff169052565b5060e08301516200218d60e084018261ffff169052565b506101008381015161ffff9081169184019190915261012080850151909116908301526101408084015163ffffffff9081169184019190915261016080850151821690840152610180808501516001600160401b03908116918501919091526101a080860151909116908401526101c09384015116929091019190915290565b602080825282518282018190526000919060409081850190868401855b82811015620022a257815180516001600160a01b031685528681015163ffffffff908116888701528682015181168787015260608083015161ffff169087015260808083015182169087015260a0808301519091169086015260c09081015115159085015260e090930192908501906001016200222a565b5091979650505050505050565b6020808252825182820181905260009190848201906040850190845b81811015620022f25783516001600160a01b031683529284019291840191600101620022cb565b50909695505050505050565b634e487b7160e01b600052601160045260246000fd5b81810381811115620010fe57620010fe620022fe565b6000816200233c576200233c620022fe565b506000190190565b63ffffffff818116838216019080821115620015495762001549620022fe565b6000604080830163ffffffff8616845260206040602086015281865180845260608701915060208801935060005b81811015620023c657845180516001600160a01b0316845284015161ffff1684840152938301939185019160010162002392565b509098975050505050505050565b8082028115828204841417620010fe57620010fe620022fe565b6000826200240c57634e487b7160e01b600052601260045260246000fd5b500490565b6001600160601b03828116828216039080821115620015495762001549620022fe565b6000602082840312156200244757600080fd5b5051919050565b8181036000831280158383131683831282161715620015495762001549620022fe565b6000602082840312156200248457600080fd5b62001139826200195c565b634e487b7160e01b600052603160045260246000fd5b60005b83811015620024c2578181015183820152602001620024a8565b50506000910152565b60008251620024df818460208701620024a5565b9190910192915050565b60208152600082518060208401526200250a816040850160208701620024a5565b601f01601f19169190910160400192915050565b60805160a05160c05160e051616193620025ba600039600081816102ec015281816125b001526130870152600081816102880152818161254b015281816136490152613df001526000818161025901528181611029015281816115d601528181611c4f0152818161252601528181612c27015281816133cb01526134c40152600081816102b80152818161257d015261359001526161936000f3fe608060405234801561001057600080fd5b50600436106101cf5760003560e01c80637437ff9f11610104578063b06d41bc116100a2578063e080bcba11610071578063e080bcba146109b7578063eff7cc48146109ca578063f2fde38b146109d2578063fbca3b74146109e557600080fd5b8063b06d41bc14610973578063c92b283214610989578063d09dc3391461099c578063df0aa9e9146109a457600080fd5b80638b364334116100de5780638b364334146109295780638da5cb5b1461093c5780639041be3d1461094d578063a69c64c01461096057600080fd5b80637437ff9f146108ae57806376f6ae761461090e57806379ba50971461092157600080fd5b806348a98aa41161017157806354b714681161014b57806354b71468146105dd578063599f6431146105fd5780636def4ce71461060e578063704b6c021461089b57600080fd5b806348a98aa41461053b578063546719cd14610566578063549e946f146105ca57600080fd5b8063181f5a77116101ad578063181f5a77146104a957806320487ded146104f257806334d560e4146105135780634816f4f71461052857600080fd5b8063061877e3146101d457806306285c69146102255780631772047e14610329575b600080fd5b6102076101e2366004614ae7565b6001600160a01b03166000908152600c602052604090205467ffffffffffffffff1690565b60405167ffffffffffffffff90911681526020015b60405180910390f35b61031c60408051608081018252600080825260208201819052918101829052606081019190915260405180608001604052807f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031681526020017f000000000000000000000000000000000000000000000000000000000000000067ffffffffffffffff1681526020017f00000000000000000000000000000000000000000000000000000000000000006bffffffffffffffffffffffff1681526020017f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316815250905090565b60405161021c9190614b04565b61043d610337366004614ae7565b6040805160e081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810191909152506001600160a01b03166000908152600d6020908152604091829020825160e081018452905463ffffffff8082168352640100000000820481169383019390935261ffff68010000000000000000820416938201939093526a01000000000000000000008304821660608201526e0100000000000000000000000000008304909116608082015260ff720100000000000000000000000000000000000083048116151560a0830152730100000000000000000000000000000000000000909204909116151560c082015290565b60405161021c9190600060e08201905063ffffffff80845116835280602085015116602084015261ffff60408501511660408401528060608501511660608401528060808501511660808401525060a0830151151560a083015260c0830151151560c083015292915050565b6104e56040518060400160405280601c81526020017f45564d3245564d4d756c74694f6e52616d7020312e362e302d6465760000000081525081565b60405161021c9190614ba4565b610505610500366004614bf0565b610a05565b60405190815260200161021c565b610526610521366004614d1a565b610e37565b005b610526610536366004614e52565b610e4b565b61054e610549366004614f8f565b610e61565b6040516001600160a01b03909116815260200161021c565b61056e610ef0565b60405161021c919081516fffffffffffffffffffffffffffffffff908116825260208084015163ffffffff1690830152604080840151151590830152606080840151821690830152608092830151169181019190915260a00190565b6105266105d8366004614fc8565b610f98565b600f546040516bffffffffffffffffffffffff909116815260200161021c565b6002546001600160a01b031661054e565b61088e61061c366004614fe6565b604080516102608101825260006080820181815260a0830182905260c0830182905260e08301829052610100830182905261012083018290526101408301829052610160830182905261018083018290526101a083018290526101c083018290526101e0830182905261020083018290526102208301829052610240830182905282526020820181905291810182905260608101919091525067ffffffffffffffff9081166000908152600b6020908152604091829020825161026081018452815460ff811615156080830190815261ffff610100808404821660a086015263ffffffff63010000008504811660c08701526701000000000000008504811660e08701526b01000000000000000000000085048116918601919091526f01000000000000000000000000000000840482166101208601527101000000000000000000000000000000000084048116610140860152750100000000000000000000000000000000000000000084048216610160860152770100000000000000000000000000000000000000000000008404821661018086015279010000000000000000000000000000000000000000000000000084049091166101a08501527b0100000000000000000000000000000000000000000000000000000090920482166101c084015260018401548083166101e0850152640100000000810488166102008501526c01000000000000000000000000810488166102208501527401000000000000000000000000000000000000000090819004909216610240840152825260028301546001600160a01b0381169483019490945290920490931691810191909152600390910154606082015290565b60405161021c9190615125565b6105266108a9366004614ae7565b611111565b610901604080516060810182526000808252602082018190529181019190915250604080516060810182526005546001600160a01b03908116825260065481166020830152600754169181019190915290565b60405161021c9190615172565b61052661091c3660046151a2565b6111d0565b610526611233565b610207610937366004614f8f565b6112f1565b6000546001600160a01b031661054e565b61020761095b366004614fe6565b6113db565b61052661096e366004615217565b61141f565b61097b611430565b60405161021c929190615329565b61052661099736600461536b565b61152b565b610505611593565b6105056109b23660046153af565b611653565b6105266109c536600461541b565b611a7f565b610526611a90565b6105266109e0366004614ae7565b611d21565b6109f86109f3366004614fe6565b611d32565b60405161021c9190615615565b67ffffffffffffffff82166000908152600b60205260408120805460ff16610a6a576040517f99ac52f200000000000000000000000000000000000000000000000000000000815267ffffffffffffffff851660048201526024015b60405180910390fd5b6000610a796080850185615662565b159050610a9a57610a95610a906080860186615662565b611d66565b610ab2565b6001820154640100000000900467ffffffffffffffff165b9050610adc85610ac56020870187615662565b905083610ad560408901896156a9565b9050611e0e565b6000600c81610af16080880160608901614ae7565b6001600160a01b03168152602081019190915260400160009081205467ffffffffffffffff169150819003610b6e57610b306080860160608701614ae7565b6040517fa7499d200000000000000000000000000000000000000000000000000000000081526001600160a01b039091166004820152602401610a61565b60065460009081906001600160a01b031663ffdb4b37610b9460808a0160608b01614ae7565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e084901b1681526001600160a01b03909116600482015267ffffffffffffffff8b1660248201526044016040805180830381865afa158015610bff573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c23919061571f565b90925090506000808080610c3a60408c018c6156a9565b90501115610c7557610c698b610c5660808d0160608e01614ae7565b87610c6460408f018f6156a9565b611f1e565b91945092509050610cac565b6001880154610ca99074010000000000000000000000000000000000000000900463ffffffff16662386f26fc10000615768565b92505b875460009077010000000000000000000000000000000000000000000000900461ffff1615610d1857610d158c6dffffffffffffffffffffffffffff607088901c16610cfb60208f018f615662565b90508e8060400190610d0d91906156a9565b90508661231f565b90505b600089600101600c9054906101000a900467ffffffffffffffff1667ffffffffffffffff168463ffffffff168b600001600f9054906101000a900461ffff1661ffff168e8060200190610d6b9190615662565b610d76929150615768565b8c54610d97906b010000000000000000000000900463ffffffff168d61577f565b610da1919061577f565b610dab919061577f565b610dc5906dffffffffffffffffffffffffffff8916615768565b610dcf9190615768565b90507bffffffffffffffffffffffffffffffffffffffffffffffffffffffff87168282610e0667ffffffffffffffff8c1689615768565b610e10919061577f565b610e1a919061577f565b610e249190615792565b9a50505050505050505050505b92915050565b610e3f612425565b610e4881612481565b50565b610e53612605565b610e5d8282612662565b5050565b6007546040517fbbe4f6db0000000000000000000000000000000000000000000000000000000081526001600160a01b038381166004830152600092169063bbe4f6db90602401602060405180830381865afa158015610ec5573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ee991906157b4565b9392505050565b6040805160a0810182526000808252602082018190529181018290526060810182905260808101919091526040805160a0810182526003546fffffffffffffffffffffffffffffffff8082168352600160801b80830463ffffffff1660208501527401000000000000000000000000000000000000000090920460ff161515938301939093526004548084166060840152049091166080820152610f93906129bc565b905090565b610fa0612605565b6001600160a01b038116610fe0576040517f232cb97f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000610fea611593565b90506000811215611027576040517f02075e0000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316836001600160a01b031603611079576110746001600160a01b0384168383612a6e565b505050565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526110749083906001600160a01b038616906370a0823190602401602060405180830381865afa1580156110dc573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061110091906157d1565b6001600160a01b0386169190612a6e565b6000546001600160a01b0316331480159061113757506002546001600160a01b03163314155b1561116e576040517ff6cd562000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6002805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0383169081179091556040519081527f8fe72c3e0020beb3234e76ae6676fa576fbfcae600af1c4fea44784cf0db329c906020015b60405180910390a150565b6111d8612605565b610e5d8282808060200260200160405190810160405280939291908181526020016000905b828210156112295761121a604083028601368190038101906157ea565b815260200190600101906111fd565b5050505050612aee565b6001546001600160a01b0316331461128d5760405162461bcd60e51b815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610a61565b600080543373ffffffffffffffffffffffffffffffffffffffff19808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6001600160a01b0381166000908152600e602052604081205467ffffffffffffffff16808203610ee95767ffffffffffffffff84166000908152600b60205260409020600201546001600160a01b031680156113d3576040517f856c82470000000000000000000000000000000000000000000000000000000081526001600160a01b03858116600483015282169063856c824790602401602060405180830381865afa1580156113a6573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906113ca9190615829565b92505050610e31565b509392505050565b67ffffffffffffffff8082166000908152600b60205260408120600201549091610e3191740100000000000000000000000000000000000000009004166001615846565b611427612605565b610e4881612d5b565b606060008061143f6008612e21565b90508067ffffffffffffffff81111561145a5761145a614c40565b60405190808252806020026020018201604052801561149f57816020015b60408051808201909152600080825260208201528152602001906001900390816114785790505b50925060005b81811015611508576000806114bb600884612e2c565b915091506040518060400160405280836001600160a01b031681526020018261ffff168152508684815181106114f3576114f3615867565b602090810291909101015250506001016114a5565b5050600f5491926c0100000000000000000000000090920463ffffffff16919050565b6000546001600160a01b0316331480159061155157506002546001600160a01b03163314155b15611588576040517ff6cd562000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610e48600382612e4a565b600f546040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526000916bffffffffffffffffffffffff16907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316906370a0823190602401602060405180830381865afa158015611625573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061164991906157d1565b610f93919061587d565b67ffffffffffffffff84166000908152600b602052604081208161167a8288888888612fe3565b905060005b81610140015151811015611a2157600061169c60408901896156a9565b838181106116ac576116ac615867565b9050604002018036038101906116c2919061589d565b905060006116d48a8360000151610e61565b90506001600160a01b038116158061178a57506040517f01ffc9a70000000000000000000000000000000000000000000000000000000081527faff2afbf0000000000000000000000000000000000000000000000000000000060048201526001600160a01b038216906301ffc9a790602401602060405180830381865afa158015611764573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061178891906158d7565b155b156117cf5781516040517fbf16aab60000000000000000000000000000000000000000000000000000000081526001600160a01b039091166004820152602401610a61565b6000816001600160a01b0316639a4575b96040518060a001604052808d80600001906117fb9190615662565b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050509082525067ffffffffffffffff8f166020808301919091526001600160a01b03808e16604080850191909152918901516060840152885116608090920191909152517fffffffff0000000000000000000000000000000000000000000000000000000060e084901b1681526118a691906004016158f4565b6000604051808303816000875af11580156118c5573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526118ed91908101906159c1565b9050602081602001515111801561193c575082516001600160a01b03166000908152600d602090815260409091205490820151516e01000000000000000000000000000090910463ffffffff16105b156119815782516040517f36f536ca0000000000000000000000000000000000000000000000000000000081526001600160a01b039091166004820152602401610a61565b805161198c9061387e565b5060408051606081019091526001600160a01b03831660808201528060a0810160405160208183030381529060405281526020018260000151815260200182602001518152506040516020016119e29190615a52565b6040516020818303038152906040528561016001518581518110611a0857611a08615867565b602002602001018190525050505080600101905061167f565b50611a308183600301546138d9565b6101808201526040517fd0c3c799bf9e2639de44391e7f524d229b2b55f5b1ea94b2bf7da42f7243dddd90611a66908390615b49565b60405180910390a161018001519150505b949350505050565b611a87612605565b610e4881613a34565b6000546001600160a01b03163314801590611ab657506002546001600160a01b03163314155b8015611aca5750611ac8600833613fde565b155b15611b01576040517f195db95800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600f546c01000000000000000000000000900463ffffffff166000819003611b55576040517f990e30bf00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600f546bffffffffffffffffffffffff1681811015611ba0576040517f8d0f71d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000611baa611593565b1215611be2576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b806000611bef6008612e21565b905060005b81811015611cde57600080611c0a600884612e2c565b9092509050600087611c2a836bffffffffffffffffffffffff8a16615768565b611c349190615792565b9050611c408187615c7e565b9550611c846001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016846bffffffffffffffffffffffff8416612a6e565b6040516bffffffffffffffffffffffff821681526001600160a01b038416907f55fdec2aab60a41fa5abb106670eb1006f5aeaee1ba7afea2bc89b5b3ec7678f9060200160405180910390a2505050806001019050611bf4565b5050600f80547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff929092169190911790555050565b611d29612425565b610e4881613ff3565b60606040517f9e7177c800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60007f97a657c900000000000000000000000000000000000000000000000000000000611d938385615ca3565b7fffffffff000000000000000000000000000000000000000000000000000000001614611dec576040517f5247fdce00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611df98260048186615ceb565b810190611e069190615d15565b519392505050565b67ffffffffffffffff84166000908152600b6020526040902080546301000000900463ffffffff16841115611e875780546040517f86933789000000000000000000000000000000000000000000000000000000008152630100000090910463ffffffff16600482015260248101859052604401610a61565b8054670100000000000000900463ffffffff16831115611ed3576040517f4c4fc93a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8054610100900461ffff16821115611f17576040517f4c056b6a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050505050565b6000808083815b81811015612312576000878783818110611f4157611f41615867565b905060400201803603810190611f57919061589d565b905060006001600160a01b0316611f728c8360000151610e61565b6001600160a01b031603611fc05780516040517fbf16aab60000000000000000000000000000000000000000000000000000000081526001600160a01b039091166004820152602401610a61565b80516001600160a01b03166000908152600d6020908152604091829020825160e081018452905463ffffffff8082168352640100000000820481169383019390935261ffff68010000000000000000820416938201939093526a01000000000000000000008304821660608201526e0100000000000000000000000000008304909116608082015260ff720100000000000000000000000000000000000083048116151560a0830152730100000000000000000000000000000000000000909204909116151560c0820181905261213e5767ffffffffffffffff8c166000908152600b6020526040902080546120de90790100000000000000000000000000000000000000000000000000900461ffff16662386f26fc10000615768565b6120e8908961577f565b815490985061211c907b01000000000000000000000000000000000000000000000000000000900463ffffffff1688615d57565b60018201549097506121349063ffffffff1687615d57565b955050505061230a565b604081015160009061ffff161561225a5760008c6001600160a01b031684600001516001600160a01b0316146121fd5760065484516040517f4ab35b0b0000000000000000000000000000000000000000000000000000000081526001600160a01b039182166004820152911690634ab35b0b90602401602060405180830381865afa1580156121d2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906121f69190615d74565b9050612200565b508a5b620186a0836040015161ffff166122428660200151847bffffffffffffffffffffffffffffffffffffffffffffffffffffffff166140a990919063ffffffff16565b61224c9190615768565b6122569190615792565b9150505b60608201516122699088615d57565b965081608001518661227b9190615d57565b825190965060009061229a9063ffffffff16662386f26fc10000615768565b9050808210156122b9576122ae818a61577f565b98505050505061230a565b6000836020015163ffffffff16662386f26fc100006122d89190615768565b9050808311156122f8576122ec818b61577f565b9950505050505061230a565b612302838b61577f565b995050505050505b600101611f25565b5050955095509592505050565b60008063ffffffff8316612334608086615768565b6123408761022061577f565b61234a919061577f565b612354919061577f565b67ffffffffffffffff88166000908152600b6020526040812080549293509171010000000000000000000000000000000000810463ffffffff16906123b6907501000000000000000000000000000000000000000000900461ffff1685615768565b6123c0919061577f565b825490915077010000000000000000000000000000000000000000000000900461ffff166123fe6dffffffffffffffffffffffffffff8a1683615768565b6124089190615768565b61241890655af3107a4000615768565b9998505050505050505050565b6000546001600160a01b0316331461247f5760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610a61565b565b60208101516001600160a01b03166124c5576040517f35be3ac800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80516005805473ffffffffffffffffffffffffffffffffffffffff199081166001600160a01b0393841617909155602080840151600680548416918516919091179055604080850151600780549094169085161790925581516080810183527f0000000000000000000000000000000000000000000000000000000000000000841681527f000000000000000000000000000000000000000000000000000000000000000067ffffffffffffffff16918101919091527f00000000000000000000000000000000000000000000000000000000000000006bffffffffffffffffffffffff16818301527f00000000000000000000000000000000000000000000000000000000000000009092166060830152517f45d99abaa76ccb5c1a18f32b24b8a795ede9926a20056247319c426c1370fefa916111c5918490615d8f565b6000546001600160a01b0316331480159061262b57506002546001600160a01b03163314155b1561247f576040517ffbdb8e5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b82518110156128d957600083828151811061268257612682615867565b6020026020010151905060208160a0015163ffffffff1610156126ef57805160a08201516040517f24ecdc020000000000000000000000000000000000000000000000000000000081526001600160a01b03909216600483015263ffffffff166024820152604401610a61565b6040805160e08101825260208381015163ffffffff908116835284840151811682840190815260608087015161ffff9081168688019081526080808a0151861693880193845260a0808b0151871691890191825260c0808c01511515918a019182526001908a018181529b516001600160a01b03166000908152600d9099529990972097518854955192519451915197519a5190871667ffffffffffffffff199096169590951764010000000092871692909202919091177fffffffffffffffffffffffffffffffffffff000000000000ffffffffffffffff166801000000000000000093909216929092027fffffffffffffffffffffffffffffffffffff00000000ffffffffffffffffffff16176a010000000000000000000091841691909102177fffffffffffffffffffffffffff0000000000ffffffffffffffffffffffffffff166e01000000000000000000000000000093909216929092027fffffffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffff1617720100000000000000000000000000000000000094151594909402939093177fffffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffff16730100000000000000000000000000000000000000931515939093029290921790915501612665565b507ff5791bc457b3bb990493cf5f655db46c25ccf5764c9b99b8969b4c72ea7df9d0826040516129099190615e0d565b60405180910390a160005b815181101561297957600d600083838151811061293357612933615867565b6020908102919091018101516001600160a01b03168252810191909152604001600020805473ffffffffffffffffffffffffffffffffffffffff19169055600101612914565b50805115610e5d577ffb95a0042158e60a33e7b5bec100f3d95407b1a71bee6633bd54b8887449750b816040516129b09190615615565b60405180910390a15050565b6040805160a081018252600080825260208201819052918101829052606081018290526080810191909152612a4a82606001516fffffffffffffffffffffffffffffffff1683600001516fffffffffffffffffffffffffffffffff16846020015163ffffffff1642612a2e9190615ead565b85608001516fffffffffffffffffffffffffffffffff166140e6565b6fffffffffffffffffffffffffffffffff1682525063ffffffff4216602082015290565b604080516001600160a01b038416602482015260448082018490528251808303909101815260649091019091526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fa9059cbb0000000000000000000000000000000000000000000000000000000017905261107490849061410e565b80516040811115612b2b576040517fb5a10cfa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600f546c01000000000000000000000000900463ffffffff1615801590612b795750600f5463ffffffff6c010000000000000000000000008204166bffffffffffffffffffffffff90911610155b15612b8657612b86611a90565b6000612b926008612e21565b90505b8015612bd4576000612bb3612bab600184615ead565b600890612e2c565b509050612bc16008826141f3565b505080612bcd90615ec0565b9050612b95565b506000805b82811015612cdc576000848281518110612bf557612bf5615867565b60200260200101516000015190506000858381518110612c1757612c17615867565b60200260200101516020015190507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316826001600160a01b03161480612c6c57506001600160a01b038216155b15612cae576040517f4de938d10000000000000000000000000000000000000000000000000000000081526001600160a01b0383166004820152602401610a61565b612cbe60088361ffff8416614208565b50612ccd61ffff821685615d57565b93505050806001019050612bd9565b50600f80547fffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffffff166c0100000000000000000000000063ffffffff8416021790556040517f8c337bff38141c507abd25c547606bdde78fe8c12e941ab613f3a565fea6cd2490612d4e9083908690615ef5565b60405180910390a1505050565b60005b8151811015610e5d576000828281518110612d7b57612d7b615867565b60200260200101516000015190506000838381518110612d9d57612d9d615867565b6020908102919091018101518101516001600160a01b0384166000818152600c8452604090819020805467ffffffffffffffff191667ffffffffffffffff85169081179091559051908152919350917fbb77da6f7210cdd16904228a9360133d1d7dfff99b1bc75f128da5b53e28f97d910160405180910390a25050600101612d5e565b6000610e318261421e565b6000808080612e3b8686614229565b909450925050505b9250929050565b8154600090612e6690600160801b900463ffffffff1642615ead565b90508015612ee35760018301548354612ea1916fffffffffffffffffffffffffffffffff808216928116918591600160801b909104166140e6565b83546fffffffffffffffffffffffffffffffff9190911673ffffffffffffffffffffffffffffffffffffffff1990911617600160801b4263ffffffff16021783555b60208201518354612f09916fffffffffffffffffffffffffffffffff9081169116614254565b83548351151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00ffffffff000000000000000000000000000000009091166fffffffffffffffffffffffffffffffff9283161717845560208301516040808501518316600160801b0291909216176001850155517f9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c1990612d4e9084908151151581526020808301516fffffffffffffffffffffffffffffffff90811691830191909152604092830151169181019190915260600190565b604080516101a08101825260008082526020820181905281830181905260608083018290526080830182905260a0830182905260c0830182905260e0830182905261010083018290526101208301819052610140830181905261016083015261018082015290517f58babe3300000000000000000000000000000000000000000000000000000000815267ffffffffffffffff861660048201526001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016906358babe3390602401602060405180830381865afa1580156130ce573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906130f291906158d7565b15613135576040517ffdbd6a7200000000000000000000000000000000000000000000000000000000815267ffffffffffffffff86166004820152602401610a61565b6001600160a01b038216613175576040517fa4ec747900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6005546001600160a01b031633146131b9576040517f1c0a352900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b855460ff16613200576040517f99ac52f200000000000000000000000000000000000000000000000000000000815267ffffffffffffffff86166004820152602401610a61565b600061320f6080860186615662565b15905061322b57613226610a906080870187615662565b613243565b6001870154640100000000900467ffffffffffffffff165b9050600061325460408701876156a9565b91506132729050876132696020890189615662565b90508484611e0e565b80156133c1576000805b828110156133af5761329160408901896156a9565b828181106132a1576132a1615867565b905060400201602001356000036132e4576040517f5cf0444900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600d60006132f560408b018b6156a9565b8481811061330557613305615867565b61331b9260206040909202019081019150614ae7565b6001600160a01b031681526020810191909152604001600020547201000000000000000000000000000000000000900460ff16156133a75761339a61336360408a018a6156a9565b8381811061337357613373615867565b905060400201803603810190613389919061589d565b6006546001600160a01b031661426a565b6133a4908361577f565b91505b60010161327c565b5080156133bf576133bf8161438b565b505b6001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000166133fb6080880160608901614ae7565b6001600160a01b03160361345f57600f805486919060009061342c9084906bffffffffffffffffffffffff16615f14565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555061357e565b6006546001600160a01b03166241e5be61347f6080890160608a01614ae7565b60405160e083901b7fffffffff000000000000000000000000000000000000000000000000000000001681526001600160a01b039182166004820152602481018990527f00000000000000000000000000000000000000000000000000000000000000009091166044820152606401602060405180830381865afa15801561350b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061352f91906157d1565b600f805460009061354f9084906bffffffffffffffffffffffff16615f14565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055505b600f546bffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000008116911611156135eb576040517fe5c7a49100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006135f788866112f1565b613602906001615846565b6001600160a01b0386166000818152600e6020908152604091829020805467ffffffffffffffff191667ffffffffffffffff8681169190911790915582516101a0810184527f00000000000000000000000000000000000000000000000000000000000000009091168152908101929092529192509081016136c16136878a80615662565b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061387e92505050565b6001600160a01b031681526020018a600201601481819054906101000a900467ffffffffffffffff166136f390615f39565b91906101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905567ffffffffffffffff1681526020018481526020016000151581526020018267ffffffffffffffff1681526020018860600160208101906137599190614ae7565b6001600160a01b0316815260200187815260200188806020019061377d9190615662565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152505050908252506020016137c460408a018a6156a9565b808060200260200160405190810160405280939291908181526020016000905b82821015613810576138016040830286013681900381019061589d565b815260200190600101906137e4565b505050505081526020018367ffffffffffffffff81111561383357613833614c40565b60405190808252806020026020018201604052801561386657816020015b60608152602001906001900390816138515790505b50815260006020909101529998505050505050505050565b600081516020146138bd57816040517f8d666f60000000000000000000000000000000000000000000000000000000008152600401610a619190614ba4565b610e31828060200190518101906138d491906157d1565b614398565b60008060001b8284602001518560400151866060015187608001518860a001518960c001518a60e001518b610100015160405160200161396f9897969594939291906001600160a01b039889168152968816602088015267ffffffffffffffff95861660408801526060870194909452911515608086015290921660a0840152921660c082015260e08101919091526101000190565b60405160208183030381529060405280519060200120856101200151805190602001208661014001516040516020016139a89190615f60565b604051602081830303815290604052805190602001208761016001516040516020016139d49190615f73565b60408051601f198184030181528282528051602091820120908301979097528101949094526060840192909252608083015260a082015260c081019190915260e00160405160208183030381529060405280519060200120905092915050565b60005b8151811015610e5d576000828281518110613a5457613a54615867565b602002602001015190506000838381518110613a7257613a72615867565b60200260200101516000015190508067ffffffffffffffff16600003613ad0576040517fc35aa79d00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff82166004820152602401610a61565b6000600b60008367ffffffffffffffff1667ffffffffffffffff16815260200190815260200160002090506000836040015190506000604051806080016040528086602001518152602001836001600160a01b031681526020018460020160149054906101000a900467ffffffffffffffff1667ffffffffffffffff1681526020018460030154815250905080600001518360000160008201518160000160006101000a81548160ff02191690831515021790555060208201518160000160016101000a81548161ffff021916908361ffff16021790555060408201518160000160036101000a81548163ffffffff021916908363ffffffff16021790555060608201518160000160076101000a81548163ffffffff021916908363ffffffff160217905550608082015181600001600b6101000a81548163ffffffff021916908363ffffffff16021790555060a082015181600001600f6101000a81548161ffff021916908361ffff16021790555060c08201518160000160116101000a81548163ffffffff021916908363ffffffff16021790555060e08201518160000160156101000a81548161ffff021916908361ffff1602179055506101008201518160000160176101000a81548161ffff021916908361ffff1602179055506101208201518160000160196101000a81548161ffff021916908361ffff16021790555061014082015181600001601b6101000a81548163ffffffff021916908363ffffffff1602179055506101608201518160010160006101000a81548163ffffffff021916908363ffffffff1602179055506101808201518160010160046101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055506101a082015181600101600c6101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055506101c08201518160010160146101000a81548163ffffffff021916908363ffffffff16021790555090505082600301546000801b03613ece57604080517f8acd72527118c8324937b1a42e02cd246697c3b633f1742f3cae11de233722b3602082015267ffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000811692820192909252908516606082015230608082015260a00160408051601f1981840301815291905280516020909101206060820181905260038401556001600160a01b03821615613e875760028301805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0384161790555b8367ffffffffffffffff167f7a70081ee29c1fc27898089ba2a5fc35ac0106b043c82ccecd24c6fd48f6ca8684604051613ec19190615f86565b60405180910390a2613fce565b60028301546001600160a01b03838116911614613f23576040517fc35aa79d00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff85166004820152602401610a61565b60208560200151610160015163ffffffff161015613f8757602085015161016001516040517f24ecdc020000000000000000000000000000000000000000000000000000000081526000600482015263ffffffff9091166024820152604401610a61565b8367ffffffffffffffff167f944eb884a589931130671ee4a7379fbe5fe65ed605a048ba99c454582f2460b08660200151604051613fc59190616112565b60405180910390a25b5050505050806001019050613a37565b6000610ee9836001600160a01b038416614404565b336001600160a01b0382160361404b5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610a61565b6001805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000670de0b6b3a76400006140dc837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8616615768565b610ee99190615792565b6000614105856140f68486615768565b614100908761577f565b614254565b95945050505050565b6000614163826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c6564815250856001600160a01b03166144109092919063ffffffff16565b805190915015611074578080602001905181019061418191906158d7565b6110745760405162461bcd60e51b815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f742073756363656564000000000000000000000000000000000000000000006064820152608401610a61565b6000610ee9836001600160a01b03841661441f565b6000611a77846001600160a01b0385168461443c565b6000610e3182614459565b600080806142378585614463565b600081815260029690960160205260409095205494959350505050565b60008183106142635781610ee9565b5090919050565b81516040517fd02641a00000000000000000000000000000000000000000000000000000000081526001600160a01b03918216600482015260009182919084169063d02641a0906024016040805180830381865afa1580156142d0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906142f49190616121565b5190507bffffffffffffffffffffffffffffffffffffffffffffffffffffffff811660000361435d5783516040517f9a655f7b0000000000000000000000000000000000000000000000000000000081526001600160a01b039091166004820152602401610a61565b6020840151611a77907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8316906140a9565b610e48600382600061446f565b60006001600160a01b038211806143af5750600a82105b156144005760408051602081018490520160408051601f19818403018152908290527f8d666f60000000000000000000000000000000000000000000000000000000008252610a6191600401614ba4565b5090565b6000610ee9838361478a565b6060611a7784846000856147a2565b60008181526002830160205260408120819055610ee98383614894565b60008281526002840160205260408120829055611a7784846148a0565b6000610e31825490565b6000610ee983836148ac565b825474010000000000000000000000000000000000000000900460ff161580614496575081155b156144a057505050565b825460018401546fffffffffffffffffffffffffffffffff808316929116906000906144d990600160801b900463ffffffff1642615ead565b9050801561457f578183111561451b576040517f9725942a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600186015461454890839085908490600160801b90046fffffffffffffffffffffffffffffffff166140e6565b86547fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff16600160801b4263ffffffff160217875592505b8482101561461c576001600160a01b0384166145d1576040517ff94ebcd10000000000000000000000000000000000000000000000000000000081526004810183905260248101869052604401610a61565b6040517f1a76572a00000000000000000000000000000000000000000000000000000000815260048101839052602481018690526001600160a01b0385166044820152606401610a61565b8483101561470857600186810154600160801b90046fffffffffffffffffffffffffffffffff169060009082906146539082615ead565b61465d878a615ead565b614667919061577f565b6146719190615792565b90506001600160a01b0386166146bd576040517f15279c080000000000000000000000000000000000000000000000000000000081526004810182905260248101869052604401610a61565b6040517fd0c8d23a00000000000000000000000000000000000000000000000000000000815260048101829052602481018690526001600160a01b0387166044820152606401610a61565b6147128584615ead565b86547fffffffffffffffffffffffffffffffff00000000000000000000000000000000166fffffffffffffffffffffffffffffffff82161787556040518681529093507f1871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690a9060200160405180910390a1505050505050565b60008181526001830160205260408120541515610ee9565b60608247101561481a5760405162461bcd60e51b815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f60448201527f722063616c6c00000000000000000000000000000000000000000000000000006064820152608401610a61565b600080866001600160a01b031685876040516148369190616154565b60006040518083038185875af1925050503d8060008114614873576040519150601f19603f3d011682016040523d82523d6000602084013e614878565b606091505b5091509150614889878383876148d6565b979650505050505050565b6000610ee9838361494f565b6000610ee98383614a49565b60008260000182815481106148c3576148c3615867565b9060005260206000200154905092915050565b6060831561494557825160000361493e576001600160a01b0385163b61493e5760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401610a61565b5081611a77565b611a778383614a98565b60008181526001830160205260408120548015614a38576000614973600183615ead565b855490915060009061498790600190615ead565b90508181146149ec5760008660000182815481106149a7576149a7615867565b90600052602060002001549050808760000184815481106149ca576149ca615867565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806149fd576149fd616170565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050610e31565b6000915050610e31565b5092915050565b6000818152600183016020526040812054614a9057508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610e31565b506000610e31565b815115614aa85781518083602001fd5b8060405162461bcd60e51b8152600401610a619190614ba4565b6001600160a01b0381168114610e4857600080fd5b8035614ae281614ac2565b919050565b600060208284031215614af957600080fd5b8135610ee981614ac2565b60808101610e31828480516001600160a01b03908116835260208083015167ffffffffffffffff16908401526040808301516bffffffffffffffffffffffff169084015260609182015116910152565b60005b83811015614b6f578181015183820152602001614b57565b50506000910152565b60008151808452614b90816020860160208601614b54565b601f01601f19169290920160200192915050565b602081526000610ee96020830184614b78565b67ffffffffffffffff81168114610e4857600080fd5b8035614ae281614bb7565b600060a08284031215614bea57600080fd5b50919050565b60008060408385031215614c0357600080fd5b8235614c0e81614bb7565b9150602083013567ffffffffffffffff811115614c2a57600080fd5b614c3685828601614bd8565b9150509250929050565b634e487b7160e01b600052604160045260246000fd5b6040516060810167ffffffffffffffff81118282101715614c7957614c79614c40565b60405290565b60405160e0810167ffffffffffffffff81118282101715614c7957614c79614c40565b6040805190810167ffffffffffffffff81118282101715614c7957614c79614c40565b6040516101e0810167ffffffffffffffff81118282101715614c7957614c79614c40565b604051601f8201601f1916810167ffffffffffffffff81118282101715614d1257614d12614c40565b604052919050565b600060608284031215614d2c57600080fd5b614d34614c56565b8235614d3f81614ac2565b81526020830135614d4f81614ac2565b60208201526040830135614d6281614ac2565b60408201529392505050565b600067ffffffffffffffff821115614d8857614d88614c40565b5060051b60200190565b63ffffffff81168114610e4857600080fd5b8035614ae281614d92565b803561ffff81168114614ae257600080fd5b8015158114610e4857600080fd5b8035614ae281614dc1565b600082601f830112614deb57600080fd5b81356020614e00614dfb83614d6e565b614ce9565b8083825260208201915060208460051b870101935086841115614e2257600080fd5b602086015b84811015614e47578035614e3a81614ac2565b8352918301918301614e27565b509695505050505050565b6000806040808486031215614e6657600080fd5b833567ffffffffffffffff80821115614e7e57600080fd5b818601915086601f830112614e9257600080fd5b81356020614ea2614dfb83614d6e565b82815260e0928302850182019282820191908b851115614ec157600080fd5b958301955b84871015614f6a5780878d031215614ede5760008081fd5b614ee6614c7f565b8735614ef181614ac2565b815287850135614f0081614d92565b8186015287890135614f1181614d92565b818a01526060614f22898201614daf565b90820152608088810135614f3581614d92565b9082015260a0614f46898201614da4565b9082015260c0614f57898201614dcf565b9082015283529586019591830191614ec6565b5097505087013593505080831115614f8157600080fd5b5050614c3685828601614dda565b60008060408385031215614fa257600080fd5b8235614fad81614bb7565b91506020830135614fbd81614ac2565b809150509250929050565b60008060408385031215614fdb57600080fd5b8235614fad81614ac2565b600060208284031215614ff857600080fd5b8135610ee981614bb7565b805115158252602081015161501e602084018261ffff169052565b506040810151615036604084018263ffffffff169052565b50606081015161504e606084018263ffffffff169052565b506080810151615066608084018263ffffffff169052565b5060a081015161507c60a084018261ffff169052565b5060c081015161509460c084018263ffffffff169052565b5060e08101516150aa60e084018261ffff169052565b506101008181015161ffff9081169184019190915261012080830151909116908301526101408082015163ffffffff90811691840191909152610160808301518216908401526101808083015167ffffffffffffffff908116918501919091526101a080840151909116908401526101c09182015116910152565b600061024082019050615139828451615003565b60208301516001600160a01b03166101e0830152604083015167ffffffffffffffff166102008301526060909201516102209091015290565b60608101610e31828480516001600160a01b03908116835260208083015182169084015260409182015116910152565b600080602083850312156151b557600080fd5b823567ffffffffffffffff808211156151cd57600080fd5b818501915085601f8301126151e157600080fd5b8135818111156151f057600080fd5b8660208260061b850101111561520557600080fd5b60209290920196919550909350505050565b6000602080838503121561522a57600080fd5b823567ffffffffffffffff81111561524157600080fd5b8301601f8101851361525257600080fd5b8035615260614dfb82614d6e565b81815260069190911b8201830190838101908783111561527f57600080fd5b928401925b82841015614889576040848903121561529d5760008081fd5b6152a5614ca2565b84356152b081614ac2565b8152848601356152bf81614bb7565b8187015282526040939093019290840190615284565b60008151808452602080850194506020840160005b8381101561531e57815180516001600160a01b0316885283015161ffff1683880152604090960195908201906001016152ea565b509495945050505050565b60408152600061533c60408301856152d5565b90508260208301529392505050565b80356fffffffffffffffffffffffffffffffff81168114614ae257600080fd5b60006060828403121561537d57600080fd5b615385614c56565b823561539081614dc1565b815261539e6020840161534b565b6020820152614d626040840161534b565b600080600080608085870312156153c557600080fd5b84356153d081614bb7565b9350602085013567ffffffffffffffff8111156153ec57600080fd5b6153f887828801614bd8565b93505060408501359150606085013561541081614ac2565b939692955090935050565b6000602080838503121561542e57600080fd5b823567ffffffffffffffff81111561544557600080fd5b8301601f8101851361545657600080fd5b8035615464614dfb82614d6e565b818152610220918202830184019184820191908884111561548457600080fd5b938501935b8385101561560957848903818112156154a25760008081fd5b6154aa614c56565b86356154b581614bb7565b81526101e0601f1983018113156154cc5760008081fd5b6154d4614cc5565b92506154e1898901614dcf565b835260406154f0818a01614daf565b8a8501526060615501818b01614da4565b828601526080615512818c01614da4565b8287015260a09150615525828c01614da4565b9086015260c06155368b8201614daf565b8287015260e09150615549828c01614da4565b9086015261010061555b8b8201614daf565b82870152610120915061556f828c01614daf565b908601526101406155818b8201614daf565b828701526101609150615595828c01614da4565b908601526101806155a78b8201614da4565b828701526101a091506155bb828c01614bcd565b908601526101c06155cd8b8201614bcd565b828701526155dc848c01614da4565b818701525050838a8401526155f46102008a01614ad7565b90830152508452509384019391850191615489565b50979650505050505050565b6020808252825182820181905260009190848201906040850190845b818110156156565783516001600160a01b031683529284019291840191600101615631565b50909695505050505050565b6000808335601e1984360301811261567957600080fd5b83018035915067ffffffffffffffff82111561569457600080fd5b602001915036819003821315612e4357600080fd5b6000808335601e198436030181126156c057600080fd5b83018035915067ffffffffffffffff8211156156db57600080fd5b6020019150600681901b3603821315612e4357600080fd5b80517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff81168114614ae257600080fd5b6000806040838503121561573257600080fd5b61573b836156f3565b9150615749602084016156f3565b90509250929050565b634e487b7160e01b600052601160045260246000fd5b8082028115828204841417610e3157610e31615752565b80820180821115610e3157610e31615752565b6000826157af57634e487b7160e01b600052601260045260246000fd5b500490565b6000602082840312156157c657600080fd5b8151610ee981614ac2565b6000602082840312156157e357600080fd5b5051919050565b6000604082840312156157fc57600080fd5b615804614ca2565b823561580f81614ac2565b815261581d60208401614daf565b60208201529392505050565b60006020828403121561583b57600080fd5b8151610ee981614bb7565b67ffffffffffffffff818116838216019080821115614a4257614a42615752565b634e487b7160e01b600052603260045260246000fd5b8181036000831280158383131683831282161715614a4257614a42615752565b6000604082840312156158af57600080fd5b6158b7614ca2565b82356158c281614ac2565b81526020928301359281019290925250919050565b6000602082840312156158e957600080fd5b8151610ee981614dc1565b602081526000825160a0602084015261591060c0840182614b78565b905067ffffffffffffffff602085015116604084015260408401516001600160a01b038082166060860152606086015160808601528060808701511660a086015250508091505092915050565b600082601f83011261596e57600080fd5b815167ffffffffffffffff81111561598857615988614c40565b61599b6020601f19601f84011601614ce9565b8181528460208386010111156159b057600080fd5b611a77826020830160208701614b54565b6000602082840312156159d357600080fd5b815167ffffffffffffffff808211156159eb57600080fd5b90830190604082860312156159ff57600080fd5b615a07614ca2565b825182811115615a1657600080fd5b615a228782860161595d565b825250602083015182811115615a3757600080fd5b615a438782860161595d565b60208301525095945050505050565b602081526000825160606020840152615a6e6080840182614b78565b90506020840151601f1980858403016040860152615a8c8383614b78565b92506040860151915080858403016060860152506141058282614b78565b60008151808452602080850194506020840160005b8381101561531e57815180516001600160a01b031688528301518388015260409096019590820190600101615abf565b60008282518085526020808601955060208260051b8401016020860160005b84811015615b3c57601f19868403018952615b2a838351614b78565b98840198925090830190600101615b0e565b5090979650505050505050565b60208152615b6460208201835167ffffffffffffffff169052565b60006020830151615b8060408401826001600160a01b03169052565b5060408301516001600160a01b038116606084015250606083015167ffffffffffffffff8116608084015250608083015160a083015260a0830151615bc960c084018215159052565b5060c083015167ffffffffffffffff811660e08401525060e0830151610100615bfc818501836001600160a01b03169052565b840151610120848101919091528401516101a061014080860182905291925090615c2a6101c0860184614b78565b9250808601519050601f19610160818786030181880152615c4b8584615aaa565b945080880151925050610180818786030181880152615c6a8584615aef565b970151959092019490945250929392505050565b6bffffffffffffffffffffffff828116828216039080821115614a4257614a42615752565b7fffffffff000000000000000000000000000000000000000000000000000000008135818116916004851015615ce35780818660040360031b1b83161692505b505092915050565b60008085851115615cfb57600080fd5b83861115615d0857600080fd5b5050820193919092039150565b600060208284031215615d2757600080fd5b6040516020810181811067ffffffffffffffff82111715615d4a57615d4a614c40565b6040529135825250919050565b63ffffffff818116838216019080821115614a4257614a42615752565b600060208284031215615d8657600080fd5b610ee9826156f3565b60e08101615ddf828580516001600160a01b03908116835260208083015167ffffffffffffffff16908401526040808301516bffffffffffffffffffffffff169084015260609182015116910152565b82516001600160a01b0390811660808401526020840151811660a084015260408401511660c0830152610ee9565b602080825282518282018190526000919060409081850190868401855b82811015615ea057815180516001600160a01b031685528681015163ffffffff908116888701528682015181168787015260608083015161ffff169087015260808083015182169087015260a0808301519091169086015260c09081015115159085015260e09093019290850190600101615e2a565b5091979650505050505050565b81810381811115610e3157610e31615752565b600081615ecf57615ecf615752565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190565b63ffffffff83168152604060208201526000611a7760408301846152d5565b6bffffffffffffffffffffffff818116838216019080821115614a4257614a42615752565b600067ffffffffffffffff808316818103615f5657615f56615752565b6001019392505050565b602081526000610ee96020830184615aaa565b602081526000610ee96020830184615aef565b815460ff81161515825261024082019061ffff600882901c8116602085015263ffffffff601883901c81166040860152615fcd60608601828560381c1663ffffffff169052565b615fe460808601828560581c1663ffffffff169052565b615ff960a08601838560781c1661ffff169052565b61601060c08601828560881c1663ffffffff169052565b61602560e08601838560a81c1661ffff169052565b61603b6101008601838560b81c1661ffff169052565b6160516101208601838560c81c1661ffff169052565b6160696101408601828560d81c1663ffffffff169052565b600186015463ffffffff82821616610160870152925067ffffffffffffffff602084901c811661018087015291506160b36101a08601838560601c1667ffffffffffffffff169052565b6160cb6101c08601828560a01c1663ffffffff169052565b5060028501546001600160a01b0381166101e086015291506160ff6102008501828460a01c1667ffffffffffffffff169052565b5050600383015461022083015292915050565b6101e08101610e318284615003565b60006040828403121561613357600080fd5b61613b614ca2565b616144836156f3565b8152602083015161581d81614d92565b60008251616166818460208701614b54565b9190910192915050565b634e487b7160e01b600052603160045260246000fdfea164736f6c6343000818000a",
}

var EVM2EVMMultiOnRampABI = EVM2EVMMultiOnRampMetaData.ABI

var EVM2EVMMultiOnRampBin = EVM2EVMMultiOnRampMetaData.Bin

func DeployEVM2EVMMultiOnRamp(auth *bind.TransactOpts, backend bind.ContractBackend, staticConfig EVM2EVMMultiOnRampStaticConfig, dynamicConfig EVM2EVMMultiOnRampDynamicConfig, destChainConfigArgs []EVM2EVMMultiOnRampDestChainConfigArgs, rateLimiterConfig RateLimiterConfig, premiumMultiplierWeiPerEthArgs []EVM2EVMMultiOnRampPremiumMultiplierWeiPerEthArgs, tokenTransferFeeConfigArgs []EVM2EVMMultiOnRampTokenTransferFeeConfigArgs, nopsAndWeights []EVM2EVMMultiOnRampNopAndWeight) (common.Address, *types.Transaction, *EVM2EVMMultiOnRamp, error) {
	parsed, err := EVM2EVMMultiOnRampMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(EVM2EVMMultiOnRampBin), backend, staticConfig, dynamicConfig, destChainConfigArgs, rateLimiterConfig, premiumMultiplierWeiPerEthArgs, tokenTransferFeeConfigArgs, nopsAndWeights)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &EVM2EVMMultiOnRamp{address: address, abi: *parsed, EVM2EVMMultiOnRampCaller: EVM2EVMMultiOnRampCaller{contract: contract}, EVM2EVMMultiOnRampTransactor: EVM2EVMMultiOnRampTransactor{contract: contract}, EVM2EVMMultiOnRampFilterer: EVM2EVMMultiOnRampFilterer{contract: contract}}, nil
}

type EVM2EVMMultiOnRamp struct {
	address common.Address
	abi     abi.ABI
	EVM2EVMMultiOnRampCaller
	EVM2EVMMultiOnRampTransactor
	EVM2EVMMultiOnRampFilterer
}

type EVM2EVMMultiOnRampCaller struct {
	contract *bind.BoundContract
}

type EVM2EVMMultiOnRampTransactor struct {
	contract *bind.BoundContract
}

type EVM2EVMMultiOnRampFilterer struct {
	contract *bind.BoundContract
}

type EVM2EVMMultiOnRampSession struct {
	Contract     *EVM2EVMMultiOnRamp
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type EVM2EVMMultiOnRampCallerSession struct {
	Contract *EVM2EVMMultiOnRampCaller
	CallOpts bind.CallOpts
}

type EVM2EVMMultiOnRampTransactorSession struct {
	Contract     *EVM2EVMMultiOnRampTransactor
	TransactOpts bind.TransactOpts
}

type EVM2EVMMultiOnRampRaw struct {
	Contract *EVM2EVMMultiOnRamp
}

type EVM2EVMMultiOnRampCallerRaw struct {
	Contract *EVM2EVMMultiOnRampCaller
}

type EVM2EVMMultiOnRampTransactorRaw struct {
	Contract *EVM2EVMMultiOnRampTransactor
}

func NewEVM2EVMMultiOnRamp(address common.Address, backend bind.ContractBackend) (*EVM2EVMMultiOnRamp, error) {
	abi, err := abi.JSON(strings.NewReader(EVM2EVMMultiOnRampABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindEVM2EVMMultiOnRamp(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EVM2EVMMultiOnRamp{address: address, abi: abi, EVM2EVMMultiOnRampCaller: EVM2EVMMultiOnRampCaller{contract: contract}, EVM2EVMMultiOnRampTransactor: EVM2EVMMultiOnRampTransactor{contract: contract}, EVM2EVMMultiOnRampFilterer: EVM2EVMMultiOnRampFilterer{contract: contract}}, nil
}

func NewEVM2EVMMultiOnRampCaller(address common.Address, caller bind.ContractCaller) (*EVM2EVMMultiOnRampCaller, error) {
	contract, err := bindEVM2EVMMultiOnRamp(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EVM2EVMMultiOnRampCaller{contract: contract}, nil
}

func NewEVM2EVMMultiOnRampTransactor(address common.Address, transactor bind.ContractTransactor) (*EVM2EVMMultiOnRampTransactor, error) {
	contract, err := bindEVM2EVMMultiOnRamp(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EVM2EVMMultiOnRampTransactor{contract: contract}, nil
}

func NewEVM2EVMMultiOnRampFilterer(address common.Address, filterer bind.ContractFilterer) (*EVM2EVMMultiOnRampFilterer, error) {
	contract, err := bindEVM2EVMMultiOnRamp(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EVM2EVMMultiOnRampFilterer{contract: contract}, nil
}

func bindEVM2EVMMultiOnRamp(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := EVM2EVMMultiOnRampMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EVM2EVMMultiOnRamp.Contract.EVM2EVMMultiOnRampCaller.contract.Call(opts, result, method, params...)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EVM2EVMMultiOnRamp.Contract.EVM2EVMMultiOnRampTransactor.contract.Transfer(opts)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EVM2EVMMultiOnRamp.Contract.EVM2EVMMultiOnRampTransactor.contract.Transact(opts, method, params...)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EVM2EVMMultiOnRamp.Contract.contract.Call(opts, result, method, params...)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EVM2EVMMultiOnRamp.Contract.contract.Transfer(opts)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EVM2EVMMultiOnRamp.Contract.contract.Transact(opts, method, params...)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampCaller) CurrentRateLimiterState(opts *bind.CallOpts) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _EVM2EVMMultiOnRamp.contract.Call(opts, &out, "currentRateLimiterState")

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampSession) CurrentRateLimiterState() (RateLimiterTokenBucket, error) {
	return _EVM2EVMMultiOnRamp.Contract.CurrentRateLimiterState(&_EVM2EVMMultiOnRamp.CallOpts)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampCallerSession) CurrentRateLimiterState() (RateLimiterTokenBucket, error) {
	return _EVM2EVMMultiOnRamp.Contract.CurrentRateLimiterState(&_EVM2EVMMultiOnRamp.CallOpts)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampCaller) GetDestChainConfig(opts *bind.CallOpts, destChainSelector uint64) (EVM2EVMMultiOnRampDestChainConfig, error) {
	var out []interface{}
	err := _EVM2EVMMultiOnRamp.contract.Call(opts, &out, "getDestChainConfig", destChainSelector)

	if err != nil {
		return *new(EVM2EVMMultiOnRampDestChainConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(EVM2EVMMultiOnRampDestChainConfig)).(*EVM2EVMMultiOnRampDestChainConfig)

	return out0, err

}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampSession) GetDestChainConfig(destChainSelector uint64) (EVM2EVMMultiOnRampDestChainConfig, error) {
	return _EVM2EVMMultiOnRamp.Contract.GetDestChainConfig(&_EVM2EVMMultiOnRamp.CallOpts, destChainSelector)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampCallerSession) GetDestChainConfig(destChainSelector uint64) (EVM2EVMMultiOnRampDestChainConfig, error) {
	return _EVM2EVMMultiOnRamp.Contract.GetDestChainConfig(&_EVM2EVMMultiOnRamp.CallOpts, destChainSelector)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampCaller) GetDynamicConfig(opts *bind.CallOpts) (EVM2EVMMultiOnRampDynamicConfig, error) {
	var out []interface{}
	err := _EVM2EVMMultiOnRamp.contract.Call(opts, &out, "getDynamicConfig")

	if err != nil {
		return *new(EVM2EVMMultiOnRampDynamicConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(EVM2EVMMultiOnRampDynamicConfig)).(*EVM2EVMMultiOnRampDynamicConfig)

	return out0, err

}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampSession) GetDynamicConfig() (EVM2EVMMultiOnRampDynamicConfig, error) {
	return _EVM2EVMMultiOnRamp.Contract.GetDynamicConfig(&_EVM2EVMMultiOnRamp.CallOpts)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampCallerSession) GetDynamicConfig() (EVM2EVMMultiOnRampDynamicConfig, error) {
	return _EVM2EVMMultiOnRamp.Contract.GetDynamicConfig(&_EVM2EVMMultiOnRamp.CallOpts)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampCaller) GetExpectedNextSequenceNumber(opts *bind.CallOpts, destChainSelector uint64) (uint64, error) {
	var out []interface{}
	err := _EVM2EVMMultiOnRamp.contract.Call(opts, &out, "getExpectedNextSequenceNumber", destChainSelector)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampSession) GetExpectedNextSequenceNumber(destChainSelector uint64) (uint64, error) {
	return _EVM2EVMMultiOnRamp.Contract.GetExpectedNextSequenceNumber(&_EVM2EVMMultiOnRamp.CallOpts, destChainSelector)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampCallerSession) GetExpectedNextSequenceNumber(destChainSelector uint64) (uint64, error) {
	return _EVM2EVMMultiOnRamp.Contract.GetExpectedNextSequenceNumber(&_EVM2EVMMultiOnRamp.CallOpts, destChainSelector)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampCaller) GetFee(opts *bind.CallOpts, destChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error) {
	var out []interface{}
	err := _EVM2EVMMultiOnRamp.contract.Call(opts, &out, "getFee", destChainSelector, message)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampSession) GetFee(destChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error) {
	return _EVM2EVMMultiOnRamp.Contract.GetFee(&_EVM2EVMMultiOnRamp.CallOpts, destChainSelector, message)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampCallerSession) GetFee(destChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error) {
	return _EVM2EVMMultiOnRamp.Contract.GetFee(&_EVM2EVMMultiOnRamp.CallOpts, destChainSelector, message)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampCaller) GetNopFeesJuels(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _EVM2EVMMultiOnRamp.contract.Call(opts, &out, "getNopFeesJuels")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampSession) GetNopFeesJuels() (*big.Int, error) {
	return _EVM2EVMMultiOnRamp.Contract.GetNopFeesJuels(&_EVM2EVMMultiOnRamp.CallOpts)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampCallerSession) GetNopFeesJuels() (*big.Int, error) {
	return _EVM2EVMMultiOnRamp.Contract.GetNopFeesJuels(&_EVM2EVMMultiOnRamp.CallOpts)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampCaller) GetNops(opts *bind.CallOpts) (GetNops,

	error) {
	var out []interface{}
	err := _EVM2EVMMultiOnRamp.contract.Call(opts, &out, "getNops")

	outstruct := new(GetNops)
	if err != nil {
		return *outstruct, err
	}

	outstruct.NopsAndWeights = *abi.ConvertType(out[0], new([]EVM2EVMMultiOnRampNopAndWeight)).(*[]EVM2EVMMultiOnRampNopAndWeight)
	outstruct.WeightsTotal = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampSession) GetNops() (GetNops,

	error) {
	return _EVM2EVMMultiOnRamp.Contract.GetNops(&_EVM2EVMMultiOnRamp.CallOpts)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampCallerSession) GetNops() (GetNops,

	error) {
	return _EVM2EVMMultiOnRamp.Contract.GetNops(&_EVM2EVMMultiOnRamp.CallOpts)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampCaller) GetPoolBySourceToken(opts *bind.CallOpts, arg0 uint64, sourceToken common.Address) (common.Address, error) {
	var out []interface{}
	err := _EVM2EVMMultiOnRamp.contract.Call(opts, &out, "getPoolBySourceToken", arg0, sourceToken)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampSession) GetPoolBySourceToken(arg0 uint64, sourceToken common.Address) (common.Address, error) {
	return _EVM2EVMMultiOnRamp.Contract.GetPoolBySourceToken(&_EVM2EVMMultiOnRamp.CallOpts, arg0, sourceToken)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampCallerSession) GetPoolBySourceToken(arg0 uint64, sourceToken common.Address) (common.Address, error) {
	return _EVM2EVMMultiOnRamp.Contract.GetPoolBySourceToken(&_EVM2EVMMultiOnRamp.CallOpts, arg0, sourceToken)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampCaller) GetPremiumMultiplierWeiPerEth(opts *bind.CallOpts, token common.Address) (uint64, error) {
	var out []interface{}
	err := _EVM2EVMMultiOnRamp.contract.Call(opts, &out, "getPremiumMultiplierWeiPerEth", token)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampSession) GetPremiumMultiplierWeiPerEth(token common.Address) (uint64, error) {
	return _EVM2EVMMultiOnRamp.Contract.GetPremiumMultiplierWeiPerEth(&_EVM2EVMMultiOnRamp.CallOpts, token)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampCallerSession) GetPremiumMultiplierWeiPerEth(token common.Address) (uint64, error) {
	return _EVM2EVMMultiOnRamp.Contract.GetPremiumMultiplierWeiPerEth(&_EVM2EVMMultiOnRamp.CallOpts, token)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampCaller) GetSenderNonce(opts *bind.CallOpts, destChainSelector uint64, sender common.Address) (uint64, error) {
	var out []interface{}
	err := _EVM2EVMMultiOnRamp.contract.Call(opts, &out, "getSenderNonce", destChainSelector, sender)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampSession) GetSenderNonce(destChainSelector uint64, sender common.Address) (uint64, error) {
	return _EVM2EVMMultiOnRamp.Contract.GetSenderNonce(&_EVM2EVMMultiOnRamp.CallOpts, destChainSelector, sender)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampCallerSession) GetSenderNonce(destChainSelector uint64, sender common.Address) (uint64, error) {
	return _EVM2EVMMultiOnRamp.Contract.GetSenderNonce(&_EVM2EVMMultiOnRamp.CallOpts, destChainSelector, sender)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampCaller) GetStaticConfig(opts *bind.CallOpts) (EVM2EVMMultiOnRampStaticConfig, error) {
	var out []interface{}
	err := _EVM2EVMMultiOnRamp.contract.Call(opts, &out, "getStaticConfig")

	if err != nil {
		return *new(EVM2EVMMultiOnRampStaticConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(EVM2EVMMultiOnRampStaticConfig)).(*EVM2EVMMultiOnRampStaticConfig)

	return out0, err

}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampSession) GetStaticConfig() (EVM2EVMMultiOnRampStaticConfig, error) {
	return _EVM2EVMMultiOnRamp.Contract.GetStaticConfig(&_EVM2EVMMultiOnRamp.CallOpts)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampCallerSession) GetStaticConfig() (EVM2EVMMultiOnRampStaticConfig, error) {
	return _EVM2EVMMultiOnRamp.Contract.GetStaticConfig(&_EVM2EVMMultiOnRamp.CallOpts)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampCaller) GetSupportedTokens(opts *bind.CallOpts, arg0 uint64) ([]common.Address, error) {
	var out []interface{}
	err := _EVM2EVMMultiOnRamp.contract.Call(opts, &out, "getSupportedTokens", arg0)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampSession) GetSupportedTokens(arg0 uint64) ([]common.Address, error) {
	return _EVM2EVMMultiOnRamp.Contract.GetSupportedTokens(&_EVM2EVMMultiOnRamp.CallOpts, arg0)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampCallerSession) GetSupportedTokens(arg0 uint64) ([]common.Address, error) {
	return _EVM2EVMMultiOnRamp.Contract.GetSupportedTokens(&_EVM2EVMMultiOnRamp.CallOpts, arg0)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampCaller) GetTokenLimitAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EVM2EVMMultiOnRamp.contract.Call(opts, &out, "getTokenLimitAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampSession) GetTokenLimitAdmin() (common.Address, error) {
	return _EVM2EVMMultiOnRamp.Contract.GetTokenLimitAdmin(&_EVM2EVMMultiOnRamp.CallOpts)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampCallerSession) GetTokenLimitAdmin() (common.Address, error) {
	return _EVM2EVMMultiOnRamp.Contract.GetTokenLimitAdmin(&_EVM2EVMMultiOnRamp.CallOpts)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampCaller) GetTokenTransferFeeConfig(opts *bind.CallOpts, token common.Address) (EVM2EVMMultiOnRampTokenTransferFeeConfig, error) {
	var out []interface{}
	err := _EVM2EVMMultiOnRamp.contract.Call(opts, &out, "getTokenTransferFeeConfig", token)

	if err != nil {
		return *new(EVM2EVMMultiOnRampTokenTransferFeeConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(EVM2EVMMultiOnRampTokenTransferFeeConfig)).(*EVM2EVMMultiOnRampTokenTransferFeeConfig)

	return out0, err

}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampSession) GetTokenTransferFeeConfig(token common.Address) (EVM2EVMMultiOnRampTokenTransferFeeConfig, error) {
	return _EVM2EVMMultiOnRamp.Contract.GetTokenTransferFeeConfig(&_EVM2EVMMultiOnRamp.CallOpts, token)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampCallerSession) GetTokenTransferFeeConfig(token common.Address) (EVM2EVMMultiOnRampTokenTransferFeeConfig, error) {
	return _EVM2EVMMultiOnRamp.Contract.GetTokenTransferFeeConfig(&_EVM2EVMMultiOnRamp.CallOpts, token)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampCaller) LinkAvailableForPayment(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _EVM2EVMMultiOnRamp.contract.Call(opts, &out, "linkAvailableForPayment")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampSession) LinkAvailableForPayment() (*big.Int, error) {
	return _EVM2EVMMultiOnRamp.Contract.LinkAvailableForPayment(&_EVM2EVMMultiOnRamp.CallOpts)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampCallerSession) LinkAvailableForPayment() (*big.Int, error) {
	return _EVM2EVMMultiOnRamp.Contract.LinkAvailableForPayment(&_EVM2EVMMultiOnRamp.CallOpts)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EVM2EVMMultiOnRamp.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampSession) Owner() (common.Address, error) {
	return _EVM2EVMMultiOnRamp.Contract.Owner(&_EVM2EVMMultiOnRamp.CallOpts)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampCallerSession) Owner() (common.Address, error) {
	return _EVM2EVMMultiOnRamp.Contract.Owner(&_EVM2EVMMultiOnRamp.CallOpts)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _EVM2EVMMultiOnRamp.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampSession) TypeAndVersion() (string, error) {
	return _EVM2EVMMultiOnRamp.Contract.TypeAndVersion(&_EVM2EVMMultiOnRamp.CallOpts)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampCallerSession) TypeAndVersion() (string, error) {
	return _EVM2EVMMultiOnRamp.Contract.TypeAndVersion(&_EVM2EVMMultiOnRamp.CallOpts)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EVM2EVMMultiOnRamp.contract.Transact(opts, "acceptOwnership")
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampSession) AcceptOwnership() (*types.Transaction, error) {
	return _EVM2EVMMultiOnRamp.Contract.AcceptOwnership(&_EVM2EVMMultiOnRamp.TransactOpts)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _EVM2EVMMultiOnRamp.Contract.AcceptOwnership(&_EVM2EVMMultiOnRamp.TransactOpts)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampTransactor) ApplyDestChainConfigUpdates(opts *bind.TransactOpts, destChainConfigArgs []EVM2EVMMultiOnRampDestChainConfigArgs) (*types.Transaction, error) {
	return _EVM2EVMMultiOnRamp.contract.Transact(opts, "applyDestChainConfigUpdates", destChainConfigArgs)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampSession) ApplyDestChainConfigUpdates(destChainConfigArgs []EVM2EVMMultiOnRampDestChainConfigArgs) (*types.Transaction, error) {
	return _EVM2EVMMultiOnRamp.Contract.ApplyDestChainConfigUpdates(&_EVM2EVMMultiOnRamp.TransactOpts, destChainConfigArgs)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampTransactorSession) ApplyDestChainConfigUpdates(destChainConfigArgs []EVM2EVMMultiOnRampDestChainConfigArgs) (*types.Transaction, error) {
	return _EVM2EVMMultiOnRamp.Contract.ApplyDestChainConfigUpdates(&_EVM2EVMMultiOnRamp.TransactOpts, destChainConfigArgs)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampTransactor) ApplyPremiumMultiplierWeiPerEthUpdates(opts *bind.TransactOpts, premiumMultiplierWeiPerEthArgs []EVM2EVMMultiOnRampPremiumMultiplierWeiPerEthArgs) (*types.Transaction, error) {
	return _EVM2EVMMultiOnRamp.contract.Transact(opts, "applyPremiumMultiplierWeiPerEthUpdates", premiumMultiplierWeiPerEthArgs)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampSession) ApplyPremiumMultiplierWeiPerEthUpdates(premiumMultiplierWeiPerEthArgs []EVM2EVMMultiOnRampPremiumMultiplierWeiPerEthArgs) (*types.Transaction, error) {
	return _EVM2EVMMultiOnRamp.Contract.ApplyPremiumMultiplierWeiPerEthUpdates(&_EVM2EVMMultiOnRamp.TransactOpts, premiumMultiplierWeiPerEthArgs)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampTransactorSession) ApplyPremiumMultiplierWeiPerEthUpdates(premiumMultiplierWeiPerEthArgs []EVM2EVMMultiOnRampPremiumMultiplierWeiPerEthArgs) (*types.Transaction, error) {
	return _EVM2EVMMultiOnRamp.Contract.ApplyPremiumMultiplierWeiPerEthUpdates(&_EVM2EVMMultiOnRamp.TransactOpts, premiumMultiplierWeiPerEthArgs)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampTransactor) ForwardFromRouter(opts *bind.TransactOpts, destChainSelector uint64, message ClientEVM2AnyMessage, feeTokenAmount *big.Int, originalSender common.Address) (*types.Transaction, error) {
	return _EVM2EVMMultiOnRamp.contract.Transact(opts, "forwardFromRouter", destChainSelector, message, feeTokenAmount, originalSender)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampSession) ForwardFromRouter(destChainSelector uint64, message ClientEVM2AnyMessage, feeTokenAmount *big.Int, originalSender common.Address) (*types.Transaction, error) {
	return _EVM2EVMMultiOnRamp.Contract.ForwardFromRouter(&_EVM2EVMMultiOnRamp.TransactOpts, destChainSelector, message, feeTokenAmount, originalSender)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampTransactorSession) ForwardFromRouter(destChainSelector uint64, message ClientEVM2AnyMessage, feeTokenAmount *big.Int, originalSender common.Address) (*types.Transaction, error) {
	return _EVM2EVMMultiOnRamp.Contract.ForwardFromRouter(&_EVM2EVMMultiOnRamp.TransactOpts, destChainSelector, message, feeTokenAmount, originalSender)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampTransactor) PayNops(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EVM2EVMMultiOnRamp.contract.Transact(opts, "payNops")
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampSession) PayNops() (*types.Transaction, error) {
	return _EVM2EVMMultiOnRamp.Contract.PayNops(&_EVM2EVMMultiOnRamp.TransactOpts)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampTransactorSession) PayNops() (*types.Transaction, error) {
	return _EVM2EVMMultiOnRamp.Contract.PayNops(&_EVM2EVMMultiOnRamp.TransactOpts)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampTransactor) SetAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _EVM2EVMMultiOnRamp.contract.Transact(opts, "setAdmin", newAdmin)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampSession) SetAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _EVM2EVMMultiOnRamp.Contract.SetAdmin(&_EVM2EVMMultiOnRamp.TransactOpts, newAdmin)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampTransactorSession) SetAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _EVM2EVMMultiOnRamp.Contract.SetAdmin(&_EVM2EVMMultiOnRamp.TransactOpts, newAdmin)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampTransactor) SetDynamicConfig(opts *bind.TransactOpts, dynamicConfig EVM2EVMMultiOnRampDynamicConfig) (*types.Transaction, error) {
	return _EVM2EVMMultiOnRamp.contract.Transact(opts, "setDynamicConfig", dynamicConfig)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampSession) SetDynamicConfig(dynamicConfig EVM2EVMMultiOnRampDynamicConfig) (*types.Transaction, error) {
	return _EVM2EVMMultiOnRamp.Contract.SetDynamicConfig(&_EVM2EVMMultiOnRamp.TransactOpts, dynamicConfig)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampTransactorSession) SetDynamicConfig(dynamicConfig EVM2EVMMultiOnRampDynamicConfig) (*types.Transaction, error) {
	return _EVM2EVMMultiOnRamp.Contract.SetDynamicConfig(&_EVM2EVMMultiOnRamp.TransactOpts, dynamicConfig)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampTransactor) SetNops(opts *bind.TransactOpts, nopsAndWeights []EVM2EVMMultiOnRampNopAndWeight) (*types.Transaction, error) {
	return _EVM2EVMMultiOnRamp.contract.Transact(opts, "setNops", nopsAndWeights)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampSession) SetNops(nopsAndWeights []EVM2EVMMultiOnRampNopAndWeight) (*types.Transaction, error) {
	return _EVM2EVMMultiOnRamp.Contract.SetNops(&_EVM2EVMMultiOnRamp.TransactOpts, nopsAndWeights)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampTransactorSession) SetNops(nopsAndWeights []EVM2EVMMultiOnRampNopAndWeight) (*types.Transaction, error) {
	return _EVM2EVMMultiOnRamp.Contract.SetNops(&_EVM2EVMMultiOnRamp.TransactOpts, nopsAndWeights)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampTransactor) SetRateLimiterConfig(opts *bind.TransactOpts, config RateLimiterConfig) (*types.Transaction, error) {
	return _EVM2EVMMultiOnRamp.contract.Transact(opts, "setRateLimiterConfig", config)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampSession) SetRateLimiterConfig(config RateLimiterConfig) (*types.Transaction, error) {
	return _EVM2EVMMultiOnRamp.Contract.SetRateLimiterConfig(&_EVM2EVMMultiOnRamp.TransactOpts, config)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampTransactorSession) SetRateLimiterConfig(config RateLimiterConfig) (*types.Transaction, error) {
	return _EVM2EVMMultiOnRamp.Contract.SetRateLimiterConfig(&_EVM2EVMMultiOnRamp.TransactOpts, config)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampTransactor) SetTokenTransferFeeConfig(opts *bind.TransactOpts, tokenTransferFeeConfigArgs []EVM2EVMMultiOnRampTokenTransferFeeConfigArgs, tokensToUseDefaultFeeConfigs []common.Address) (*types.Transaction, error) {
	return _EVM2EVMMultiOnRamp.contract.Transact(opts, "setTokenTransferFeeConfig", tokenTransferFeeConfigArgs, tokensToUseDefaultFeeConfigs)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampSession) SetTokenTransferFeeConfig(tokenTransferFeeConfigArgs []EVM2EVMMultiOnRampTokenTransferFeeConfigArgs, tokensToUseDefaultFeeConfigs []common.Address) (*types.Transaction, error) {
	return _EVM2EVMMultiOnRamp.Contract.SetTokenTransferFeeConfig(&_EVM2EVMMultiOnRamp.TransactOpts, tokenTransferFeeConfigArgs, tokensToUseDefaultFeeConfigs)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampTransactorSession) SetTokenTransferFeeConfig(tokenTransferFeeConfigArgs []EVM2EVMMultiOnRampTokenTransferFeeConfigArgs, tokensToUseDefaultFeeConfigs []common.Address) (*types.Transaction, error) {
	return _EVM2EVMMultiOnRamp.Contract.SetTokenTransferFeeConfig(&_EVM2EVMMultiOnRamp.TransactOpts, tokenTransferFeeConfigArgs, tokensToUseDefaultFeeConfigs)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _EVM2EVMMultiOnRamp.contract.Transact(opts, "transferOwnership", to)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _EVM2EVMMultiOnRamp.Contract.TransferOwnership(&_EVM2EVMMultiOnRamp.TransactOpts, to)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _EVM2EVMMultiOnRamp.Contract.TransferOwnership(&_EVM2EVMMultiOnRamp.TransactOpts, to)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampTransactor) WithdrawNonLinkFees(opts *bind.TransactOpts, feeToken common.Address, to common.Address) (*types.Transaction, error) {
	return _EVM2EVMMultiOnRamp.contract.Transact(opts, "withdrawNonLinkFees", feeToken, to)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampSession) WithdrawNonLinkFees(feeToken common.Address, to common.Address) (*types.Transaction, error) {
	return _EVM2EVMMultiOnRamp.Contract.WithdrawNonLinkFees(&_EVM2EVMMultiOnRamp.TransactOpts, feeToken, to)
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampTransactorSession) WithdrawNonLinkFees(feeToken common.Address, to common.Address) (*types.Transaction, error) {
	return _EVM2EVMMultiOnRamp.Contract.WithdrawNonLinkFees(&_EVM2EVMMultiOnRamp.TransactOpts, feeToken, to)
}

type EVM2EVMMultiOnRampAdminSetIterator struct {
	Event *EVM2EVMMultiOnRampAdminSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMMultiOnRampAdminSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMMultiOnRampAdminSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(EVM2EVMMultiOnRampAdminSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *EVM2EVMMultiOnRampAdminSetIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMMultiOnRampAdminSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMMultiOnRampAdminSet struct {
	NewAdmin common.Address
	Raw      types.Log
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampFilterer) FilterAdminSet(opts *bind.FilterOpts) (*EVM2EVMMultiOnRampAdminSetIterator, error) {

	logs, sub, err := _EVM2EVMMultiOnRamp.contract.FilterLogs(opts, "AdminSet")
	if err != nil {
		return nil, err
	}
	return &EVM2EVMMultiOnRampAdminSetIterator{contract: _EVM2EVMMultiOnRamp.contract, event: "AdminSet", logs: logs, sub: sub}, nil
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampFilterer) WatchAdminSet(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOnRampAdminSet) (event.Subscription, error) {

	logs, sub, err := _EVM2EVMMultiOnRamp.contract.WatchLogs(opts, "AdminSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMMultiOnRampAdminSet)
				if err := _EVM2EVMMultiOnRamp.contract.UnpackLog(event, "AdminSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampFilterer) ParseAdminSet(log types.Log) (*EVM2EVMMultiOnRampAdminSet, error) {
	event := new(EVM2EVMMultiOnRampAdminSet)
	if err := _EVM2EVMMultiOnRamp.contract.UnpackLog(event, "AdminSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMMultiOnRampCCIPSendRequestedIterator struct {
	Event *EVM2EVMMultiOnRampCCIPSendRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMMultiOnRampCCIPSendRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMMultiOnRampCCIPSendRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(EVM2EVMMultiOnRampCCIPSendRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *EVM2EVMMultiOnRampCCIPSendRequestedIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMMultiOnRampCCIPSendRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMMultiOnRampCCIPSendRequested struct {
	Message InternalEVM2EVMMessage
	Raw     types.Log
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampFilterer) FilterCCIPSendRequested(opts *bind.FilterOpts) (*EVM2EVMMultiOnRampCCIPSendRequestedIterator, error) {

	logs, sub, err := _EVM2EVMMultiOnRamp.contract.FilterLogs(opts, "CCIPSendRequested")
	if err != nil {
		return nil, err
	}
	return &EVM2EVMMultiOnRampCCIPSendRequestedIterator{contract: _EVM2EVMMultiOnRamp.contract, event: "CCIPSendRequested", logs: logs, sub: sub}, nil
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampFilterer) WatchCCIPSendRequested(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOnRampCCIPSendRequested) (event.Subscription, error) {

	logs, sub, err := _EVM2EVMMultiOnRamp.contract.WatchLogs(opts, "CCIPSendRequested")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMMultiOnRampCCIPSendRequested)
				if err := _EVM2EVMMultiOnRamp.contract.UnpackLog(event, "CCIPSendRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampFilterer) ParseCCIPSendRequested(log types.Log) (*EVM2EVMMultiOnRampCCIPSendRequested, error) {
	event := new(EVM2EVMMultiOnRampCCIPSendRequested)
	if err := _EVM2EVMMultiOnRamp.contract.UnpackLog(event, "CCIPSendRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMMultiOnRampConfigChangedIterator struct {
	Event *EVM2EVMMultiOnRampConfigChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMMultiOnRampConfigChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMMultiOnRampConfigChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(EVM2EVMMultiOnRampConfigChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *EVM2EVMMultiOnRampConfigChangedIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMMultiOnRampConfigChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMMultiOnRampConfigChanged struct {
	Config RateLimiterConfig
	Raw    types.Log
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampFilterer) FilterConfigChanged(opts *bind.FilterOpts) (*EVM2EVMMultiOnRampConfigChangedIterator, error) {

	logs, sub, err := _EVM2EVMMultiOnRamp.contract.FilterLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return &EVM2EVMMultiOnRampConfigChangedIterator{contract: _EVM2EVMMultiOnRamp.contract, event: "ConfigChanged", logs: logs, sub: sub}, nil
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampFilterer) WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOnRampConfigChanged) (event.Subscription, error) {

	logs, sub, err := _EVM2EVMMultiOnRamp.contract.WatchLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMMultiOnRampConfigChanged)
				if err := _EVM2EVMMultiOnRamp.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampFilterer) ParseConfigChanged(log types.Log) (*EVM2EVMMultiOnRampConfigChanged, error) {
	event := new(EVM2EVMMultiOnRampConfigChanged)
	if err := _EVM2EVMMultiOnRamp.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMMultiOnRampConfigSetIterator struct {
	Event *EVM2EVMMultiOnRampConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMMultiOnRampConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMMultiOnRampConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(EVM2EVMMultiOnRampConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *EVM2EVMMultiOnRampConfigSetIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMMultiOnRampConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMMultiOnRampConfigSet struct {
	StaticConfig  EVM2EVMMultiOnRampStaticConfig
	DynamicConfig EVM2EVMMultiOnRampDynamicConfig
	Raw           types.Log
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampFilterer) FilterConfigSet(opts *bind.FilterOpts) (*EVM2EVMMultiOnRampConfigSetIterator, error) {

	logs, sub, err := _EVM2EVMMultiOnRamp.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &EVM2EVMMultiOnRampConfigSetIterator{contract: _EVM2EVMMultiOnRamp.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOnRampConfigSet) (event.Subscription, error) {

	logs, sub, err := _EVM2EVMMultiOnRamp.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMMultiOnRampConfigSet)
				if err := _EVM2EVMMultiOnRamp.contract.UnpackLog(event, "ConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampFilterer) ParseConfigSet(log types.Log) (*EVM2EVMMultiOnRampConfigSet, error) {
	event := new(EVM2EVMMultiOnRampConfigSet)
	if err := _EVM2EVMMultiOnRamp.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMMultiOnRampDestChainAddedIterator struct {
	Event *EVM2EVMMultiOnRampDestChainAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMMultiOnRampDestChainAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMMultiOnRampDestChainAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(EVM2EVMMultiOnRampDestChainAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *EVM2EVMMultiOnRampDestChainAddedIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMMultiOnRampDestChainAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMMultiOnRampDestChainAdded struct {
	DestChainSelector uint64
	DestChainConfig   EVM2EVMMultiOnRampDestChainConfig
	Raw               types.Log
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampFilterer) FilterDestChainAdded(opts *bind.FilterOpts, destChainSelector []uint64) (*EVM2EVMMultiOnRampDestChainAddedIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _EVM2EVMMultiOnRamp.contract.FilterLogs(opts, "DestChainAdded", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &EVM2EVMMultiOnRampDestChainAddedIterator{contract: _EVM2EVMMultiOnRamp.contract, event: "DestChainAdded", logs: logs, sub: sub}, nil
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampFilterer) WatchDestChainAdded(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOnRampDestChainAdded, destChainSelector []uint64) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _EVM2EVMMultiOnRamp.contract.WatchLogs(opts, "DestChainAdded", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMMultiOnRampDestChainAdded)
				if err := _EVM2EVMMultiOnRamp.contract.UnpackLog(event, "DestChainAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampFilterer) ParseDestChainAdded(log types.Log) (*EVM2EVMMultiOnRampDestChainAdded, error) {
	event := new(EVM2EVMMultiOnRampDestChainAdded)
	if err := _EVM2EVMMultiOnRamp.contract.UnpackLog(event, "DestChainAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMMultiOnRampDestChainDynamicConfigUpdatedIterator struct {
	Event *EVM2EVMMultiOnRampDestChainDynamicConfigUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMMultiOnRampDestChainDynamicConfigUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMMultiOnRampDestChainDynamicConfigUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(EVM2EVMMultiOnRampDestChainDynamicConfigUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *EVM2EVMMultiOnRampDestChainDynamicConfigUpdatedIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMMultiOnRampDestChainDynamicConfigUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMMultiOnRampDestChainDynamicConfigUpdated struct {
	DestChainSelector uint64
	DynamicConfig     EVM2EVMMultiOnRampDestChainDynamicConfig
	Raw               types.Log
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampFilterer) FilterDestChainDynamicConfigUpdated(opts *bind.FilterOpts, destChainSelector []uint64) (*EVM2EVMMultiOnRampDestChainDynamicConfigUpdatedIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _EVM2EVMMultiOnRamp.contract.FilterLogs(opts, "DestChainDynamicConfigUpdated", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &EVM2EVMMultiOnRampDestChainDynamicConfigUpdatedIterator{contract: _EVM2EVMMultiOnRamp.contract, event: "DestChainDynamicConfigUpdated", logs: logs, sub: sub}, nil
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampFilterer) WatchDestChainDynamicConfigUpdated(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOnRampDestChainDynamicConfigUpdated, destChainSelector []uint64) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _EVM2EVMMultiOnRamp.contract.WatchLogs(opts, "DestChainDynamicConfigUpdated", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMMultiOnRampDestChainDynamicConfigUpdated)
				if err := _EVM2EVMMultiOnRamp.contract.UnpackLog(event, "DestChainDynamicConfigUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampFilterer) ParseDestChainDynamicConfigUpdated(log types.Log) (*EVM2EVMMultiOnRampDestChainDynamicConfigUpdated, error) {
	event := new(EVM2EVMMultiOnRampDestChainDynamicConfigUpdated)
	if err := _EVM2EVMMultiOnRamp.contract.UnpackLog(event, "DestChainDynamicConfigUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMMultiOnRampNopPaidIterator struct {
	Event *EVM2EVMMultiOnRampNopPaid

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMMultiOnRampNopPaidIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMMultiOnRampNopPaid)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(EVM2EVMMultiOnRampNopPaid)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *EVM2EVMMultiOnRampNopPaidIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMMultiOnRampNopPaidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMMultiOnRampNopPaid struct {
	Nop    common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampFilterer) FilterNopPaid(opts *bind.FilterOpts, nop []common.Address) (*EVM2EVMMultiOnRampNopPaidIterator, error) {

	var nopRule []interface{}
	for _, nopItem := range nop {
		nopRule = append(nopRule, nopItem)
	}

	logs, sub, err := _EVM2EVMMultiOnRamp.contract.FilterLogs(opts, "NopPaid", nopRule)
	if err != nil {
		return nil, err
	}
	return &EVM2EVMMultiOnRampNopPaidIterator{contract: _EVM2EVMMultiOnRamp.contract, event: "NopPaid", logs: logs, sub: sub}, nil
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampFilterer) WatchNopPaid(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOnRampNopPaid, nop []common.Address) (event.Subscription, error) {

	var nopRule []interface{}
	for _, nopItem := range nop {
		nopRule = append(nopRule, nopItem)
	}

	logs, sub, err := _EVM2EVMMultiOnRamp.contract.WatchLogs(opts, "NopPaid", nopRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMMultiOnRampNopPaid)
				if err := _EVM2EVMMultiOnRamp.contract.UnpackLog(event, "NopPaid", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampFilterer) ParseNopPaid(log types.Log) (*EVM2EVMMultiOnRampNopPaid, error) {
	event := new(EVM2EVMMultiOnRampNopPaid)
	if err := _EVM2EVMMultiOnRamp.contract.UnpackLog(event, "NopPaid", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMMultiOnRampNopsSetIterator struct {
	Event *EVM2EVMMultiOnRampNopsSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMMultiOnRampNopsSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMMultiOnRampNopsSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(EVM2EVMMultiOnRampNopsSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *EVM2EVMMultiOnRampNopsSetIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMMultiOnRampNopsSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMMultiOnRampNopsSet struct {
	NopWeightsTotal *big.Int
	NopsAndWeights  []EVM2EVMMultiOnRampNopAndWeight
	Raw             types.Log
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampFilterer) FilterNopsSet(opts *bind.FilterOpts) (*EVM2EVMMultiOnRampNopsSetIterator, error) {

	logs, sub, err := _EVM2EVMMultiOnRamp.contract.FilterLogs(opts, "NopsSet")
	if err != nil {
		return nil, err
	}
	return &EVM2EVMMultiOnRampNopsSetIterator{contract: _EVM2EVMMultiOnRamp.contract, event: "NopsSet", logs: logs, sub: sub}, nil
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampFilterer) WatchNopsSet(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOnRampNopsSet) (event.Subscription, error) {

	logs, sub, err := _EVM2EVMMultiOnRamp.contract.WatchLogs(opts, "NopsSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMMultiOnRampNopsSet)
				if err := _EVM2EVMMultiOnRamp.contract.UnpackLog(event, "NopsSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampFilterer) ParseNopsSet(log types.Log) (*EVM2EVMMultiOnRampNopsSet, error) {
	event := new(EVM2EVMMultiOnRampNopsSet)
	if err := _EVM2EVMMultiOnRamp.contract.UnpackLog(event, "NopsSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMMultiOnRampOwnershipTransferRequestedIterator struct {
	Event *EVM2EVMMultiOnRampOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMMultiOnRampOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMMultiOnRampOwnershipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(EVM2EVMMultiOnRampOwnershipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *EVM2EVMMultiOnRampOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMMultiOnRampOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMMultiOnRampOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*EVM2EVMMultiOnRampOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _EVM2EVMMultiOnRamp.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &EVM2EVMMultiOnRampOwnershipTransferRequestedIterator{contract: _EVM2EVMMultiOnRamp.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOnRampOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _EVM2EVMMultiOnRamp.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMMultiOnRampOwnershipTransferRequested)
				if err := _EVM2EVMMultiOnRamp.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampFilterer) ParseOwnershipTransferRequested(log types.Log) (*EVM2EVMMultiOnRampOwnershipTransferRequested, error) {
	event := new(EVM2EVMMultiOnRampOwnershipTransferRequested)
	if err := _EVM2EVMMultiOnRamp.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMMultiOnRampOwnershipTransferredIterator struct {
	Event *EVM2EVMMultiOnRampOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMMultiOnRampOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMMultiOnRampOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(EVM2EVMMultiOnRampOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *EVM2EVMMultiOnRampOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMMultiOnRampOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMMultiOnRampOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*EVM2EVMMultiOnRampOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _EVM2EVMMultiOnRamp.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &EVM2EVMMultiOnRampOwnershipTransferredIterator{contract: _EVM2EVMMultiOnRamp.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOnRampOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _EVM2EVMMultiOnRamp.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMMultiOnRampOwnershipTransferred)
				if err := _EVM2EVMMultiOnRamp.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampFilterer) ParseOwnershipTransferred(log types.Log) (*EVM2EVMMultiOnRampOwnershipTransferred, error) {
	event := new(EVM2EVMMultiOnRampOwnershipTransferred)
	if err := _EVM2EVMMultiOnRamp.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMMultiOnRampPremiumMultiplierWeiPerEthUpdatedIterator struct {
	Event *EVM2EVMMultiOnRampPremiumMultiplierWeiPerEthUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMMultiOnRampPremiumMultiplierWeiPerEthUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMMultiOnRampPremiumMultiplierWeiPerEthUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(EVM2EVMMultiOnRampPremiumMultiplierWeiPerEthUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *EVM2EVMMultiOnRampPremiumMultiplierWeiPerEthUpdatedIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMMultiOnRampPremiumMultiplierWeiPerEthUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMMultiOnRampPremiumMultiplierWeiPerEthUpdated struct {
	Token                      common.Address
	PremiumMultiplierWeiPerEth uint64
	Raw                        types.Log
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampFilterer) FilterPremiumMultiplierWeiPerEthUpdated(opts *bind.FilterOpts, token []common.Address) (*EVM2EVMMultiOnRampPremiumMultiplierWeiPerEthUpdatedIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _EVM2EVMMultiOnRamp.contract.FilterLogs(opts, "PremiumMultiplierWeiPerEthUpdated", tokenRule)
	if err != nil {
		return nil, err
	}
	return &EVM2EVMMultiOnRampPremiumMultiplierWeiPerEthUpdatedIterator{contract: _EVM2EVMMultiOnRamp.contract, event: "PremiumMultiplierWeiPerEthUpdated", logs: logs, sub: sub}, nil
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampFilterer) WatchPremiumMultiplierWeiPerEthUpdated(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOnRampPremiumMultiplierWeiPerEthUpdated, token []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _EVM2EVMMultiOnRamp.contract.WatchLogs(opts, "PremiumMultiplierWeiPerEthUpdated", tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMMultiOnRampPremiumMultiplierWeiPerEthUpdated)
				if err := _EVM2EVMMultiOnRamp.contract.UnpackLog(event, "PremiumMultiplierWeiPerEthUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampFilterer) ParsePremiumMultiplierWeiPerEthUpdated(log types.Log) (*EVM2EVMMultiOnRampPremiumMultiplierWeiPerEthUpdated, error) {
	event := new(EVM2EVMMultiOnRampPremiumMultiplierWeiPerEthUpdated)
	if err := _EVM2EVMMultiOnRamp.contract.UnpackLog(event, "PremiumMultiplierWeiPerEthUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMMultiOnRampTokenTransferFeeConfigDeletedIterator struct {
	Event *EVM2EVMMultiOnRampTokenTransferFeeConfigDeleted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMMultiOnRampTokenTransferFeeConfigDeletedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMMultiOnRampTokenTransferFeeConfigDeleted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(EVM2EVMMultiOnRampTokenTransferFeeConfigDeleted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *EVM2EVMMultiOnRampTokenTransferFeeConfigDeletedIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMMultiOnRampTokenTransferFeeConfigDeletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMMultiOnRampTokenTransferFeeConfigDeleted struct {
	Tokens []common.Address
	Raw    types.Log
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampFilterer) FilterTokenTransferFeeConfigDeleted(opts *bind.FilterOpts) (*EVM2EVMMultiOnRampTokenTransferFeeConfigDeletedIterator, error) {

	logs, sub, err := _EVM2EVMMultiOnRamp.contract.FilterLogs(opts, "TokenTransferFeeConfigDeleted")
	if err != nil {
		return nil, err
	}
	return &EVM2EVMMultiOnRampTokenTransferFeeConfigDeletedIterator{contract: _EVM2EVMMultiOnRamp.contract, event: "TokenTransferFeeConfigDeleted", logs: logs, sub: sub}, nil
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampFilterer) WatchTokenTransferFeeConfigDeleted(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOnRampTokenTransferFeeConfigDeleted) (event.Subscription, error) {

	logs, sub, err := _EVM2EVMMultiOnRamp.contract.WatchLogs(opts, "TokenTransferFeeConfigDeleted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMMultiOnRampTokenTransferFeeConfigDeleted)
				if err := _EVM2EVMMultiOnRamp.contract.UnpackLog(event, "TokenTransferFeeConfigDeleted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampFilterer) ParseTokenTransferFeeConfigDeleted(log types.Log) (*EVM2EVMMultiOnRampTokenTransferFeeConfigDeleted, error) {
	event := new(EVM2EVMMultiOnRampTokenTransferFeeConfigDeleted)
	if err := _EVM2EVMMultiOnRamp.contract.UnpackLog(event, "TokenTransferFeeConfigDeleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMMultiOnRampTokenTransferFeeConfigSetIterator struct {
	Event *EVM2EVMMultiOnRampTokenTransferFeeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMMultiOnRampTokenTransferFeeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMMultiOnRampTokenTransferFeeConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(EVM2EVMMultiOnRampTokenTransferFeeConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *EVM2EVMMultiOnRampTokenTransferFeeConfigSetIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMMultiOnRampTokenTransferFeeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMMultiOnRampTokenTransferFeeConfigSet struct {
	TransferFeeConfig []EVM2EVMMultiOnRampTokenTransferFeeConfigArgs
	Raw               types.Log
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampFilterer) FilterTokenTransferFeeConfigSet(opts *bind.FilterOpts) (*EVM2EVMMultiOnRampTokenTransferFeeConfigSetIterator, error) {

	logs, sub, err := _EVM2EVMMultiOnRamp.contract.FilterLogs(opts, "TokenTransferFeeConfigSet")
	if err != nil {
		return nil, err
	}
	return &EVM2EVMMultiOnRampTokenTransferFeeConfigSetIterator{contract: _EVM2EVMMultiOnRamp.contract, event: "TokenTransferFeeConfigSet", logs: logs, sub: sub}, nil
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampFilterer) WatchTokenTransferFeeConfigSet(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOnRampTokenTransferFeeConfigSet) (event.Subscription, error) {

	logs, sub, err := _EVM2EVMMultiOnRamp.contract.WatchLogs(opts, "TokenTransferFeeConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMMultiOnRampTokenTransferFeeConfigSet)
				if err := _EVM2EVMMultiOnRamp.contract.UnpackLog(event, "TokenTransferFeeConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampFilterer) ParseTokenTransferFeeConfigSet(log types.Log) (*EVM2EVMMultiOnRampTokenTransferFeeConfigSet, error) {
	event := new(EVM2EVMMultiOnRampTokenTransferFeeConfigSet)
	if err := _EVM2EVMMultiOnRamp.contract.UnpackLog(event, "TokenTransferFeeConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMMultiOnRampTokensConsumedIterator struct {
	Event *EVM2EVMMultiOnRampTokensConsumed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMMultiOnRampTokensConsumedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMMultiOnRampTokensConsumed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(EVM2EVMMultiOnRampTokensConsumed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *EVM2EVMMultiOnRampTokensConsumedIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMMultiOnRampTokensConsumedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMMultiOnRampTokensConsumed struct {
	Tokens *big.Int
	Raw    types.Log
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampFilterer) FilterTokensConsumed(opts *bind.FilterOpts) (*EVM2EVMMultiOnRampTokensConsumedIterator, error) {

	logs, sub, err := _EVM2EVMMultiOnRamp.contract.FilterLogs(opts, "TokensConsumed")
	if err != nil {
		return nil, err
	}
	return &EVM2EVMMultiOnRampTokensConsumedIterator{contract: _EVM2EVMMultiOnRamp.contract, event: "TokensConsumed", logs: logs, sub: sub}, nil
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampFilterer) WatchTokensConsumed(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOnRampTokensConsumed) (event.Subscription, error) {

	logs, sub, err := _EVM2EVMMultiOnRamp.contract.WatchLogs(opts, "TokensConsumed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMMultiOnRampTokensConsumed)
				if err := _EVM2EVMMultiOnRamp.contract.UnpackLog(event, "TokensConsumed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRampFilterer) ParseTokensConsumed(log types.Log) (*EVM2EVMMultiOnRampTokensConsumed, error) {
	event := new(EVM2EVMMultiOnRampTokensConsumed)
	if err := _EVM2EVMMultiOnRamp.contract.UnpackLog(event, "TokensConsumed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetNops struct {
	NopsAndWeights []EVM2EVMMultiOnRampNopAndWeight
	WeightsTotal   *big.Int
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRamp) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _EVM2EVMMultiOnRamp.abi.Events["AdminSet"].ID:
		return _EVM2EVMMultiOnRamp.ParseAdminSet(log)
	case _EVM2EVMMultiOnRamp.abi.Events["CCIPSendRequested"].ID:
		return _EVM2EVMMultiOnRamp.ParseCCIPSendRequested(log)
	case _EVM2EVMMultiOnRamp.abi.Events["ConfigChanged"].ID:
		return _EVM2EVMMultiOnRamp.ParseConfigChanged(log)
	case _EVM2EVMMultiOnRamp.abi.Events["ConfigSet"].ID:
		return _EVM2EVMMultiOnRamp.ParseConfigSet(log)
	case _EVM2EVMMultiOnRamp.abi.Events["DestChainAdded"].ID:
		return _EVM2EVMMultiOnRamp.ParseDestChainAdded(log)
	case _EVM2EVMMultiOnRamp.abi.Events["DestChainDynamicConfigUpdated"].ID:
		return _EVM2EVMMultiOnRamp.ParseDestChainDynamicConfigUpdated(log)
	case _EVM2EVMMultiOnRamp.abi.Events["NopPaid"].ID:
		return _EVM2EVMMultiOnRamp.ParseNopPaid(log)
	case _EVM2EVMMultiOnRamp.abi.Events["NopsSet"].ID:
		return _EVM2EVMMultiOnRamp.ParseNopsSet(log)
	case _EVM2EVMMultiOnRamp.abi.Events["OwnershipTransferRequested"].ID:
		return _EVM2EVMMultiOnRamp.ParseOwnershipTransferRequested(log)
	case _EVM2EVMMultiOnRamp.abi.Events["OwnershipTransferred"].ID:
		return _EVM2EVMMultiOnRamp.ParseOwnershipTransferred(log)
	case _EVM2EVMMultiOnRamp.abi.Events["PremiumMultiplierWeiPerEthUpdated"].ID:
		return _EVM2EVMMultiOnRamp.ParsePremiumMultiplierWeiPerEthUpdated(log)
	case _EVM2EVMMultiOnRamp.abi.Events["TokenTransferFeeConfigDeleted"].ID:
		return _EVM2EVMMultiOnRamp.ParseTokenTransferFeeConfigDeleted(log)
	case _EVM2EVMMultiOnRamp.abi.Events["TokenTransferFeeConfigSet"].ID:
		return _EVM2EVMMultiOnRamp.ParseTokenTransferFeeConfigSet(log)
	case _EVM2EVMMultiOnRamp.abi.Events["TokensConsumed"].ID:
		return _EVM2EVMMultiOnRamp.ParseTokensConsumed(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (EVM2EVMMultiOnRampAdminSet) Topic() common.Hash {
	return common.HexToHash("0x8fe72c3e0020beb3234e76ae6676fa576fbfcae600af1c4fea44784cf0db329c")
}

func (EVM2EVMMultiOnRampCCIPSendRequested) Topic() common.Hash {
	return common.HexToHash("0xd0c3c799bf9e2639de44391e7f524d229b2b55f5b1ea94b2bf7da42f7243dddd")
}

func (EVM2EVMMultiOnRampConfigChanged) Topic() common.Hash {
	return common.HexToHash("0x9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c19")
}

func (EVM2EVMMultiOnRampConfigSet) Topic() common.Hash {
	return common.HexToHash("0x45d99abaa76ccb5c1a18f32b24b8a795ede9926a20056247319c426c1370fefa")
}

func (EVM2EVMMultiOnRampDestChainAdded) Topic() common.Hash {
	return common.HexToHash("0x7a70081ee29c1fc27898089ba2a5fc35ac0106b043c82ccecd24c6fd48f6ca86")
}

func (EVM2EVMMultiOnRampDestChainDynamicConfigUpdated) Topic() common.Hash {
	return common.HexToHash("0x944eb884a589931130671ee4a7379fbe5fe65ed605a048ba99c454582f2460b0")
}

func (EVM2EVMMultiOnRampNopPaid) Topic() common.Hash {
	return common.HexToHash("0x55fdec2aab60a41fa5abb106670eb1006f5aeaee1ba7afea2bc89b5b3ec7678f")
}

func (EVM2EVMMultiOnRampNopsSet) Topic() common.Hash {
	return common.HexToHash("0x8c337bff38141c507abd25c547606bdde78fe8c12e941ab613f3a565fea6cd24")
}

func (EVM2EVMMultiOnRampOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (EVM2EVMMultiOnRampOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (EVM2EVMMultiOnRampPremiumMultiplierWeiPerEthUpdated) Topic() common.Hash {
	return common.HexToHash("0xbb77da6f7210cdd16904228a9360133d1d7dfff99b1bc75f128da5b53e28f97d")
}

func (EVM2EVMMultiOnRampTokenTransferFeeConfigDeleted) Topic() common.Hash {
	return common.HexToHash("0xfb95a0042158e60a33e7b5bec100f3d95407b1a71bee6633bd54b8887449750b")
}

func (EVM2EVMMultiOnRampTokenTransferFeeConfigSet) Topic() common.Hash {
	return common.HexToHash("0xf5791bc457b3bb990493cf5f655db46c25ccf5764c9b99b8969b4c72ea7df9d0")
}

func (EVM2EVMMultiOnRampTokensConsumed) Topic() common.Hash {
	return common.HexToHash("0x1871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690a")
}

func (_EVM2EVMMultiOnRamp *EVM2EVMMultiOnRamp) Address() common.Address {
	return _EVM2EVMMultiOnRamp.address
}

type EVM2EVMMultiOnRampInterface interface {
	CurrentRateLimiterState(opts *bind.CallOpts) (RateLimiterTokenBucket, error)

	GetDestChainConfig(opts *bind.CallOpts, destChainSelector uint64) (EVM2EVMMultiOnRampDestChainConfig, error)

	GetDynamicConfig(opts *bind.CallOpts) (EVM2EVMMultiOnRampDynamicConfig, error)

	GetExpectedNextSequenceNumber(opts *bind.CallOpts, destChainSelector uint64) (uint64, error)

	GetFee(opts *bind.CallOpts, destChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error)

	GetNopFeesJuels(opts *bind.CallOpts) (*big.Int, error)

	GetNops(opts *bind.CallOpts) (GetNops,

		error)

	GetPoolBySourceToken(opts *bind.CallOpts, arg0 uint64, sourceToken common.Address) (common.Address, error)

	GetPremiumMultiplierWeiPerEth(opts *bind.CallOpts, token common.Address) (uint64, error)

	GetSenderNonce(opts *bind.CallOpts, destChainSelector uint64, sender common.Address) (uint64, error)

	GetStaticConfig(opts *bind.CallOpts) (EVM2EVMMultiOnRampStaticConfig, error)

	GetSupportedTokens(opts *bind.CallOpts, arg0 uint64) ([]common.Address, error)

	GetTokenLimitAdmin(opts *bind.CallOpts) (common.Address, error)

	GetTokenTransferFeeConfig(opts *bind.CallOpts, token common.Address) (EVM2EVMMultiOnRampTokenTransferFeeConfig, error)

	LinkAvailableForPayment(opts *bind.CallOpts) (*big.Int, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	ApplyDestChainConfigUpdates(opts *bind.TransactOpts, destChainConfigArgs []EVM2EVMMultiOnRampDestChainConfigArgs) (*types.Transaction, error)

	ApplyPremiumMultiplierWeiPerEthUpdates(opts *bind.TransactOpts, premiumMultiplierWeiPerEthArgs []EVM2EVMMultiOnRampPremiumMultiplierWeiPerEthArgs) (*types.Transaction, error)

	ForwardFromRouter(opts *bind.TransactOpts, destChainSelector uint64, message ClientEVM2AnyMessage, feeTokenAmount *big.Int, originalSender common.Address) (*types.Transaction, error)

	PayNops(opts *bind.TransactOpts) (*types.Transaction, error)

	SetAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error)

	SetDynamicConfig(opts *bind.TransactOpts, dynamicConfig EVM2EVMMultiOnRampDynamicConfig) (*types.Transaction, error)

	SetNops(opts *bind.TransactOpts, nopsAndWeights []EVM2EVMMultiOnRampNopAndWeight) (*types.Transaction, error)

	SetRateLimiterConfig(opts *bind.TransactOpts, config RateLimiterConfig) (*types.Transaction, error)

	SetTokenTransferFeeConfig(opts *bind.TransactOpts, tokenTransferFeeConfigArgs []EVM2EVMMultiOnRampTokenTransferFeeConfigArgs, tokensToUseDefaultFeeConfigs []common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	WithdrawNonLinkFees(opts *bind.TransactOpts, feeToken common.Address, to common.Address) (*types.Transaction, error)

	FilterAdminSet(opts *bind.FilterOpts) (*EVM2EVMMultiOnRampAdminSetIterator, error)

	WatchAdminSet(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOnRampAdminSet) (event.Subscription, error)

	ParseAdminSet(log types.Log) (*EVM2EVMMultiOnRampAdminSet, error)

	FilterCCIPSendRequested(opts *bind.FilterOpts) (*EVM2EVMMultiOnRampCCIPSendRequestedIterator, error)

	WatchCCIPSendRequested(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOnRampCCIPSendRequested) (event.Subscription, error)

	ParseCCIPSendRequested(log types.Log) (*EVM2EVMMultiOnRampCCIPSendRequested, error)

	FilterConfigChanged(opts *bind.FilterOpts) (*EVM2EVMMultiOnRampConfigChangedIterator, error)

	WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOnRampConfigChanged) (event.Subscription, error)

	ParseConfigChanged(log types.Log) (*EVM2EVMMultiOnRampConfigChanged, error)

	FilterConfigSet(opts *bind.FilterOpts) (*EVM2EVMMultiOnRampConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOnRampConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*EVM2EVMMultiOnRampConfigSet, error)

	FilterDestChainAdded(opts *bind.FilterOpts, destChainSelector []uint64) (*EVM2EVMMultiOnRampDestChainAddedIterator, error)

	WatchDestChainAdded(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOnRampDestChainAdded, destChainSelector []uint64) (event.Subscription, error)

	ParseDestChainAdded(log types.Log) (*EVM2EVMMultiOnRampDestChainAdded, error)

	FilterDestChainDynamicConfigUpdated(opts *bind.FilterOpts, destChainSelector []uint64) (*EVM2EVMMultiOnRampDestChainDynamicConfigUpdatedIterator, error)

	WatchDestChainDynamicConfigUpdated(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOnRampDestChainDynamicConfigUpdated, destChainSelector []uint64) (event.Subscription, error)

	ParseDestChainDynamicConfigUpdated(log types.Log) (*EVM2EVMMultiOnRampDestChainDynamicConfigUpdated, error)

	FilterNopPaid(opts *bind.FilterOpts, nop []common.Address) (*EVM2EVMMultiOnRampNopPaidIterator, error)

	WatchNopPaid(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOnRampNopPaid, nop []common.Address) (event.Subscription, error)

	ParseNopPaid(log types.Log) (*EVM2EVMMultiOnRampNopPaid, error)

	FilterNopsSet(opts *bind.FilterOpts) (*EVM2EVMMultiOnRampNopsSetIterator, error)

	WatchNopsSet(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOnRampNopsSet) (event.Subscription, error)

	ParseNopsSet(log types.Log) (*EVM2EVMMultiOnRampNopsSet, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*EVM2EVMMultiOnRampOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOnRampOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*EVM2EVMMultiOnRampOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*EVM2EVMMultiOnRampOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOnRampOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*EVM2EVMMultiOnRampOwnershipTransferred, error)

	FilterPremiumMultiplierWeiPerEthUpdated(opts *bind.FilterOpts, token []common.Address) (*EVM2EVMMultiOnRampPremiumMultiplierWeiPerEthUpdatedIterator, error)

	WatchPremiumMultiplierWeiPerEthUpdated(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOnRampPremiumMultiplierWeiPerEthUpdated, token []common.Address) (event.Subscription, error)

	ParsePremiumMultiplierWeiPerEthUpdated(log types.Log) (*EVM2EVMMultiOnRampPremiumMultiplierWeiPerEthUpdated, error)

	FilterTokenTransferFeeConfigDeleted(opts *bind.FilterOpts) (*EVM2EVMMultiOnRampTokenTransferFeeConfigDeletedIterator, error)

	WatchTokenTransferFeeConfigDeleted(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOnRampTokenTransferFeeConfigDeleted) (event.Subscription, error)

	ParseTokenTransferFeeConfigDeleted(log types.Log) (*EVM2EVMMultiOnRampTokenTransferFeeConfigDeleted, error)

	FilterTokenTransferFeeConfigSet(opts *bind.FilterOpts) (*EVM2EVMMultiOnRampTokenTransferFeeConfigSetIterator, error)

	WatchTokenTransferFeeConfigSet(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOnRampTokenTransferFeeConfigSet) (event.Subscription, error)

	ParseTokenTransferFeeConfigSet(log types.Log) (*EVM2EVMMultiOnRampTokenTransferFeeConfigSet, error)

	FilterTokensConsumed(opts *bind.FilterOpts) (*EVM2EVMMultiOnRampTokensConsumedIterator, error)

	WatchTokensConsumed(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOnRampTokensConsumed) (event.Subscription, error)

	ParseTokensConsumed(log types.Log) (*EVM2EVMMultiOnRampTokensConsumed, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
