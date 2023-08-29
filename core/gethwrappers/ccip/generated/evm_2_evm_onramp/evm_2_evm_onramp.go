// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package evm_2_evm_onramp

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

type EVM2EVMOnRampDynamicConfig struct {
	Router                common.Address
	MaxTokensLength       uint16
	DestGasOverhead       uint32
	DestGasPerPayloadByte uint16
	PriceRegistry         common.Address
	MaxDataSize           uint32
	MaxGasLimit           uint64
}

type EVM2EVMOnRampFeeTokenConfig struct {
	NetworkFeeUSD          uint32
	MinTokenTransferFeeUSD uint32
	MaxTokenTransferFeeUSD uint32
	GasMultiplier          uint64
	PremiumMultiplier      uint64
	Enabled                bool
}

type EVM2EVMOnRampFeeTokenConfigArgs struct {
	Token                  common.Address
	NetworkFeeUSD          uint32
	MinTokenTransferFeeUSD uint32
	MaxTokenTransferFeeUSD uint32
	GasMultiplier          uint64
	PremiumMultiplier      uint64
	Enabled                bool
}

type EVM2EVMOnRampNopAndWeight struct {
	Nop    common.Address
	Weight uint16
}

type EVM2EVMOnRampStaticConfig struct {
	LinkToken         common.Address
	ChainSelector     uint64
	DestChainSelector uint64
	DefaultTxGasLimit uint64
	MaxNopFeesJuels   *big.Int
	PrevOnRamp        common.Address
	ArmProxy          common.Address
}

type EVM2EVMOnRampTokenTransferFeeConfig struct {
	Ratio           uint16
	DestGasOverhead uint32
}

type EVM2EVMOnRampTokenTransferFeeConfigArgs struct {
	Token           common.Address
	Ratio           uint16
	DestGasOverhead uint32
}

type InternalEVM2EVMMessage struct {
	SourceChainSelector uint64
	SequenceNumber      uint64
	FeeTokenAmount      *big.Int
	Sender              common.Address
	Nonce               uint64
	GasLimit            *big.Int
	Strict              bool
	Receiver            common.Address
	Data                []byte
	TokenAmounts        []ClientEVMTokenAmount
	FeeToken            common.Address
	MessageId           [32]byte
}

type InternalPoolUpdate struct {
	Token common.Address
	Pool  common.Address
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

var EVM2EVMOnRampMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"linkToken\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"defaultTxGasLimit\",\"type\":\"uint64\"},{\"internalType\":\"uint96\",\"name\":\"maxNopFeesJuels\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"prevOnRamp\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"armProxy\",\"type\":\"address\"}],\"internalType\":\"structEVM2EVMOnRamp.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"maxTokensLength\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerPayloadByte\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"priceRegistry\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"maxDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"maxGasLimit\",\"type\":\"uint64\"}],\"internalType\":\"structEVM2EVMOnRamp.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"}],\"internalType\":\"structInternal.PoolUpdate[]\",\"name\":\"tokensAndPools\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"rateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"networkFeeUSD\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"minTokenTransferFeeUSD\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxTokenTransferFeeUSD\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"gasMultiplier\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"premiumMultiplier\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"internalType\":\"structEVM2EVMOnRamp.FeeTokenConfigArgs[]\",\"name\":\"feeTokenConfigs\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"ratio\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"}],\"internalType\":\"structEVM2EVMOnRamp.TokenTransferFeeConfigArgs[]\",\"name\":\"tokenTransferFeeConfigArgs\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"nop\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"weight\",\"type\":\"uint16\"}],\"internalType\":\"structEVM2EVMOnRamp.NopAndWeight[]\",\"name\":\"nopsAndWeights\",\"type\":\"tuple[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"}],\"name\":\"AggregateValueMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"}],\"name\":\"AggregateValueRateLimitReached\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BadARMSignal\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BucketOverfilled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedAddress\",\"type\":\"bytes\"}],\"name\":\"InvalidAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidConfig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidExtraArgsTag\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"nop\",\"type\":\"address\"}],\"name\":\"InvalidNopAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTokenPoolConfig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidWithdrawParams\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LinkBalanceNotSettled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxFeeBalanceReached\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MessageGasLimitTooHigh\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"maxSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actualSize\",\"type\":\"uint256\"}],\"name\":\"MessageTooLarge\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeCalledByRouter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoFeesToPay\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoNopsToPay\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"NotAFeeToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByAdminOrOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAdminOrNop\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PoolAlreadyAdded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"PoolDoesNotExist\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"PriceNotFoundForToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RouterMustSetOriginalSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TokenPoolMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenRateLimitReached\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyNops\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnsupportedNumberOfTokens\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"UnsupportedToken\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"feeTokenAmount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"strict\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMTokenAmount[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"}],\"indexed\":false,\"internalType\":\"structInternal.EVM2EVMMessage\",\"name\":\"message\",\"type\":\"tuple\"}],\"name\":\"CCIPSendRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"linkToken\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"defaultTxGasLimit\",\"type\":\"uint64\"},{\"internalType\":\"uint96\",\"name\":\"maxNopFeesJuels\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"prevOnRamp\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"armProxy\",\"type\":\"address\"}],\"indexed\":false,\"internalType\":\"structEVM2EVMOnRamp.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"maxTokensLength\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerPayloadByte\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"priceRegistry\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"maxDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"maxGasLimit\",\"type\":\"uint64\"}],\"indexed\":false,\"internalType\":\"structEVM2EVMOnRamp.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"networkFeeUSD\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"minTokenTransferFeeUSD\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxTokenTransferFeeUSD\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"gasMultiplier\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"premiumMultiplier\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"indexed\":false,\"internalType\":\"structEVM2EVMOnRamp.FeeTokenConfigArgs[]\",\"name\":\"feeConfig\",\"type\":\"tuple[]\"}],\"name\":\"FeeConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"nop\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"NopPaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nopWeightsTotal\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"nop\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"weight\",\"type\":\"uint16\"}],\"indexed\":false,\"internalType\":\"structEVM2EVMOnRamp.NopAndWeight[]\",\"name\":\"nopsAndWeights\",\"type\":\"tuple[]\"}],\"name\":\"NopsSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"}],\"name\":\"PoolAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"}],\"name\":\"PoolRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"ratio\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"}],\"indexed\":false,\"internalType\":\"structEVM2EVMOnRamp.TokenTransferFeeConfigArgs[]\",\"name\":\"transferFeeConfig\",\"type\":\"tuple[]\"}],\"name\":\"TokenTransferFeeConfigSet\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"}],\"internalType\":\"structInternal.PoolUpdate[]\",\"name\":\"removes\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"}],\"internalType\":\"structInternal.PoolUpdate[]\",\"name\":\"adds\",\"type\":\"tuple[]\"}],\"name\":\"applyPoolUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMTokenAmount[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structClient.EVM2AnyMessage\",\"name\":\"message\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"feeTokenAmount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"originalSender\",\"type\":\"address\"}],\"name\":\"forwardFromRouter\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDynamicConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"maxTokensLength\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerPayloadByte\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"priceRegistry\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"maxDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"maxGasLimit\",\"type\":\"uint64\"}],\"internalType\":\"structEVM2EVMOnRamp.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getExpectedNextSequenceNumber\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMTokenAmount[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structClient.EVM2AnyMessage\",\"name\":\"message\",\"type\":\"tuple\"}],\"name\":\"getFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getFeeTokenConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"networkFeeUSD\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"minTokenTransferFeeUSD\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxTokenTransferFeeUSD\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"gasMultiplier\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"premiumMultiplier\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"internalType\":\"structEVM2EVMOnRamp.FeeTokenConfig\",\"name\":\"feeTokenConfig\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNopFeesJuels\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNops\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"nop\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"weight\",\"type\":\"uint16\"}],\"internalType\":\"structEVM2EVMOnRamp.NopAndWeight[]\",\"name\":\"nopsAndWeights\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"weightsTotal\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"sourceToken\",\"type\":\"address\"}],\"name\":\"getPoolBySourceToken\",\"outputs\":[{\"internalType\":\"contractIPool\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"getSenderNonce\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStaticConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"linkToken\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"defaultTxGasLimit\",\"type\":\"uint64\"},{\"internalType\":\"uint96\",\"name\":\"maxNopFeesJuels\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"prevOnRamp\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"armProxy\",\"type\":\"address\"}],\"internalType\":\"structEVM2EVMOnRamp.StaticConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSupportedTokens\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTokenLimitAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getTokenTransferFeeConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint16\",\"name\":\"ratio\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"}],\"internalType\":\"structEVM2EVMOnRamp.TokenTransferFeeConfig\",\"name\":\"tokenTransferFeeConfig\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkAvailableForPayment\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"payNops\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"setAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"maxTokensLength\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerPayloadByte\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"priceRegistry\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"maxDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"maxGasLimit\",\"type\":\"uint64\"}],\"internalType\":\"structEVM2EVMOnRamp.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"name\":\"setDynamicConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"networkFeeUSD\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"minTokenTransferFeeUSD\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxTokenTransferFeeUSD\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"gasMultiplier\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"premiumMultiplier\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"internalType\":\"structEVM2EVMOnRamp.FeeTokenConfigArgs[]\",\"name\":\"feeTokenConfigArgs\",\"type\":\"tuple[]\"}],\"name\":\"setFeeTokenConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"nop\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"weight\",\"type\":\"uint16\"}],\"internalType\":\"structEVM2EVMOnRamp.NopAndWeight[]\",\"name\":\"nopsAndWeights\",\"type\":\"tuple[]\"}],\"name\":\"setNops\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"setRateLimiterConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"ratio\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"}],\"internalType\":\"structEVM2EVMOnRamp.TokenTransferFeeConfigArgs[]\",\"name\":\"tokenTransferFeeConfigArgs\",\"type\":\"tuple[]\"}],\"name\":\"setTokenTransferFeeConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdrawNonLinkFees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6101806040526010805460ff60c01b191690553480156200001f57600080fd5b5060405162007c8d38038062007c8d833981016040819052620000429162001c71565b8333806000816200009a5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000cd57620000cd8162000330565b50506040805160a081018252602084810180516001600160801b039081168085524263ffffffff169385018490528751151585870181905292518216606086018190529790950151166080909301839052600380546001600160a01b031916909417600160801b9283021760ff60a01b1916600160a01b90910217909255029091176004555086516001600160a01b0316158062000176575060208701516001600160401b0316155b806200018d575060408701516001600160401b0316155b80620001a4575060608701516001600160401b0316155b80620001bb575060c08701516001600160a01b0316155b15620001da576040516306b7c75960e31b815260040160405180910390fd5b6020808801516040808a015181517fbdd59ac4dd1d82276c9a9c5d2656546346b9dcdb1f9b4204aed4ec15c23d7d3a948101949094526001600160401b039283169184019190915216606082015230608082015260a00160408051601f198184030181529181528151602092830120608090815289516001600160a01b0390811660e052928a01516001600160401b0390811661010052918a015182166101205260608a015190911660a0908152908901516001600160601b031660c0908152908901518216610140528801511661016052620002b786620003db565b620002c28362000565565b620002cd8262000700565b620002d881620007d7565b6040805160008082526020820190925262000323916200031b565b6040805180820190915260008082526020820152815260200190600190039081620002f35790505b508662000a01565b50505050505050620021c0565b336001600160a01b038216036200038a5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000091565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60808101516001600160a01b031662000407576040516306b7c75960e31b815260040160405180910390fd5b8051600580546020808501516040808701516060808901516001600160a01b039889166001600160b01b031990971696909617600160a01b61ffff95861681029190911765ffffffffffff60b01b1916600160b01b63ffffffff9485160261ffff60d01b191617600160d01b9590971694909402959095179095556080808801516006805460a0808c015160c0808e0151958d166001600160c01b0319909416939093179a16909602989098176001600160c01b0316600160c01b6001600160401b0393841602179055825160e08082018552518916815261010051821695810195909552610120518116858401528351169484019490945284516001600160601b031693830193909352610140518516908201526101605190931691830191909152517f72c6aaba4dde02f77d291123a76185c418ba63f8c217a2d56b08aec84e9bbfb8916200055a91849062001d6c565b60405180910390a150565b60005b8151811015620006ce57600082828151811062000589576200058962001e38565b6020908102919091018101516040805160c080820183528385015163ffffffff908116835283850151811683870190815260608087015183168587019081526080808901516001600160401b0390811693880193845260a0808b01518216928901928352968a0151151596880196875298516001600160a01b03166000908152600d909a529690982094518554925198519151965194511515600160e01b0260ff60e01b19958916600160a01b0295909516600160a01b600160e81b0319979098166c0100000000000000000000000002600160601b600160a01b0319928516680100000000000000000292909216600160401b600160a01b0319998516640100000000026001600160401b0319909416919094161791909117969096161794909417919091169190911791909117905550620006c68162001e64565b905062000568565b507f2386f61ab5cafc3fed44f9f614f721ab53479ef64067fd16c1a2491b63ddf1a8816040516200055a919062001e80565b60005b8151811015620007a557600082828151811062000724576200072462001e38565b6020908102919091018101516040805180820182528284015161ffff90811682528284015163ffffffff90811683870190815294516001600160a01b03166000908152600e90965292909420905181549351909216620100000265ffffffffffff199093169190931617179055506200079d8162001e64565b905062000703565b507f4230c60a9725eb5fb992cf6a215398b4e81b4606d4a1e6be8dfe0b60dc172ec1816040516200055a919062001f38565b80516040811115620007fc57604051635ad0867d60e11b815260040160405180910390fd5b6010546c01000000000000000000000000900463ffffffff161580159062000846575060105463ffffffff6c010000000000000000000000008204166001600160601b0390911610155b1562000856576200085662000d04565b600062000864600762000f04565b90505b8015620008b05760006200088a6200088160018462001f99565b60079062000f17565b5090506200089a60078262000f35565b505080620008a89062001faf565b905062000867565b506000805b8281101562000998576000848281518110620008d557620008d562001e38565b60200260200101516000015190506000858381518110620008fa57620008fa62001e38565b602002602001015160200151905060e0516001600160a01b0316826001600160a01b031614806200093257506001600160a01b038216155b156200095d57604051634de938d160e01b81526001600160a01b038316600482015260240162000091565b6200096f60078361ffff841662000f53565b506200098061ffff82168562001fc9565b9350505080620009909062001e64565b9050620008b5565b506010805463ffffffff60601b19166c0100000000000000000000000063ffffffff8416021790556040517f8c337bff38141c507abd25c547606bdde78fe8c12e941ab613f3a565fea6cd2490620009f4908390869062001fe9565b60405180910390a1505050565b60005b825181101562000b3d57600083828151811062000a255762000a2562001e38565b6020026020010151600001519050600084838151811062000a4a5762000a4a62001e38565b6020908102919091018101510151905062000a67600a8362000f73565b62000a91576040516373913ebd60e01b81526001600160a01b038316600482015260240162000091565b6001600160a01b03811662000aa8600a8462000f8a565b6001600160a01b03161462000ad057604051630d98f73360e31b815260040160405180910390fd5b62000add600a8362000fa1565b1562000b2757604080516001600160a01b038085168252831660208201527f987eb3c2f78454541205f72f34839b434c306c9eaf4922efd7c0c3060fdb2e4c910160405180910390a15b50508062000b359062001e64565b905062000a04565b5060005b815181101562000cff57600082828151811062000b625762000b6262001e38565b6020026020010151600001519050600083838151811062000b875762000b8762001e38565b602002602001015160200151905060006001600160a01b0316826001600160a01b0316148062000bbe57506001600160a01b038116155b1562000bdc5760405162d8548360e71b815260040160405180910390fd5b806001600160a01b03166321df0da76040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000c1b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000c41919062002056565b6001600160a01b0316826001600160a01b03161462000c7357604051630d98f73360e31b815260040160405180910390fd5b62000c81600a838362000fb8565b1562000cd057604080516001600160a01b038085168252831660208201527f95f865c2808f8b2a85eea2611db7843150ee7835ef1403f9755918a97d76933c910160405180910390a162000ce9565b604051633caf458560e01b815260040160405180910390fd5b50508062000cf79062001e64565b905062000b41565b505050565b6000546001600160a01b0316331480159062000d2b57506002546001600160a01b03163314155b801562000d42575062000d4060073362000fd0565b155b1562000d615760405163032bb72b60e31b815260040160405180910390fd5b6010546c01000000000000000000000000900463ffffffff16600081900362000d9d5760405163990e30bf60e01b815260040160405180910390fd5b6010546001600160601b03168181101562000dcb576040516311a1ee3b60e31b815260040160405180910390fd5b600062000dd762000fe7565b121562000df757604051631e9acf1760e31b815260040160405180910390fd5b80600062000e06600762000f04565b905060005b8181101562000ede5760008062000e2460078462000f17565b909250905060008762000e41836001600160601b038a1662002076565b62000e4d919062002090565b905062000e5b8187620020b3565b60e05190965062000e80906001600160a01b0316846001600160601b03841662001075565b6040516001600160601b03821681526001600160a01b038416907f55fdec2aab60a41fa5abb106670eb1006f5aeaee1ba7afea2bc89b5b3ec7678f9060200160405180910390a25050508062000ed69062001e64565b905062000e0b565b5050601080546001600160601b0319166001600160601b03929092169190911790555050565b600062000f1182620010cd565b92915050565b600080808062000f288686620010da565b9097909650945050505050565b600062000f4c836001600160a01b03841662001107565b9392505050565b600062000f6b846001600160a01b0385168462001126565b949350505050565b600062000f4c836001600160a01b03841662001145565b600062000f4c836001600160a01b03841662001153565b600062000f4c836001600160a01b03841662001161565b600062000f6b846001600160a01b038516846200116f565b600062000f4c836001600160a01b03841662001187565b60105460e0516040516370a0823160e01b81523060048201526000926001600160601b0316916001600160a01b0316906370a0823190602401602060405180830381865afa1580156200103e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620010649190620020d6565b620010709190620020f0565b905090565b604080516001600160a01b038416602482015260448082018490528251808303909101815260649091019091526020810180516001600160e01b0390811663a9059cbb60e01b1790915262000cff9185916200119516565b600062000f118262001266565b60008080620010ea858562001271565b600081815260029690960160205260409095205494959350505050565b6000818152600283016020526040812081905562000f4c83836200127f565b6000828152600284016020526040812082905562000f6b84846200128d565b600062000f4c838362001187565b600062000f4c83836200129b565b600062000f4c838362001107565b600062000f6b84846001600160a01b03851662001126565b600062000f4c838362001310565b6040805180820190915260208082527f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c656490820152600090620011e4906001600160a01b03851690849062001329565b80519091501562000cff578080602001905181019062001205919062002113565b62000cff5760405162461bcd60e51b815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e6044820152691bdd081cdd58d8d9595960b21b606482015260840162000091565b600062000f11825490565b600062000f4c83836200133a565b600062000f4c838362001367565b600062000f4c838362001472565b600081815260028301602052604081205480151580620012c25750620012c2848462001187565b62000f4c5760405162461bcd60e51b815260206004820152601e60248201527f456e756d657261626c654d61703a206e6f6e6578697374656e74206b65790000604482015260640162000091565b6000818152600183016020526040812054151562000f4c565b606062000f6b8484600085620014c4565b600082600001828154811062001354576200135462001e38565b9060005260206000200154905092915050565b60008181526001830160205260408120548015620014605760006200138e60018362001f99565b8554909150600090620013a49060019062001f99565b905081811462001410576000866000018281548110620013c857620013c862001e38565b9060005260206000200154905080876000018481548110620013ee57620013ee62001e38565b6000918252602080832090910192909255918252600188019052604090208390555b855486908062001424576200142462002131565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505062000f11565b600091505062000f11565b5092915050565b6000818152600183016020526040812054620014bb5750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915562000f11565b50600062000f11565b606082471015620015275760405162461bcd60e51b815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f6044820152651c8818d85b1b60d21b606482015260840162000091565b600080866001600160a01b031685876040516200154591906200216d565b60006040518083038185875af1925050503d806000811462001584576040519150601f19603f3d011682016040523d82523d6000602084013e62001589565b606091505b5090925090506200159d87838387620015a8565b979650505050505050565b606083156200161c57825160000362001614576001600160a01b0385163b620016145760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000604482015260640162000091565b508162000f6b565b62000f6b8383815115620016335781518083602001fd5b8060405162461bcd60e51b81526004016200009191906200218b565b634e487b7160e01b600052604160045260246000fd5b60405160e081016001600160401b03811182821017156200168a576200168a6200164f565b60405290565b604080519081016001600160401b03811182821017156200168a576200168a6200164f565b604051606081016001600160401b03811182821017156200168a576200168a6200164f565b604051601f8201601f191681016001600160401b03811182821017156200170557620017056200164f565b604052919050565b6001600160a01b03811681146200172357600080fd5b50565b805162001733816200170d565b919050565b80516001600160401b03811681146200173357600080fd5b600060e082840312156200176357600080fd5b6200176d62001665565b905081516200177c816200170d565b81526200178c6020830162001738565b60208201526200179f6040830162001738565b6040820152620017b26060830162001738565b606082015260808201516001600160601b0381168114620017d257600080fd5b6080820152620017e560a0830162001726565b60a0820152620017f860c0830162001726565b60c082015292915050565b805161ffff811681146200173357600080fd5b805163ffffffff811681146200173357600080fd5b600060e082840312156200183e57600080fd5b6200184862001665565b9050815162001857816200170d565b8152620018676020830162001803565b60208201526200187a6040830162001816565b60408201526200188d6060830162001803565b60608201526080820151620018a2816200170d565b6080820152620018b560a0830162001816565b60a0820152620017f860c0830162001738565b60006001600160401b03821115620018e457620018e46200164f565b5060051b60200190565b600082601f8301126200190057600080fd5b81516020620019196200191383620018c8565b620016da565b82815260069290921b840181019181810190868411156200193957600080fd5b8286015b84811015620019935760408189031215620019585760008081fd5b6200196262001690565b81516200196f816200170d565b81528185015162001980816200170d565b818601528352918301916040016200193d565b509695505050505050565b805180151581146200173357600080fd5b80516001600160801b03811681146200173357600080fd5b600060608284031215620019da57600080fd5b620019e4620016b5565b9050620019f1826200199e565b815262001a0160208301620019af565b602082015262001a1460408301620019af565b604082015292915050565b600082601f83011262001a3157600080fd5b8151602062001a446200191383620018c8565b82815260e0928302850182019282820191908785111562001a6457600080fd5b8387015b8581101562001b195781818a03121562001a825760008081fd5b62001a8c62001665565b815162001a99816200170d565b815262001aa882870162001816565b86820152604062001abb81840162001816565b90820152606062001ace83820162001816565b90820152608062001ae183820162001738565b9082015260a062001af483820162001738565b9082015260c062001b078382016200199e565b90820152845292840192810162001a68565b5090979650505050505050565b600082601f83011262001b3857600080fd5b8151602062001b4b6200191383620018c8565b8281526060928302850182019282820191908785111562001b6b57600080fd5b8387015b8581101562001b195781818a03121562001b895760008081fd5b62001b93620016b5565b815162001ba0816200170d565b815262001baf82870162001803565b86820152604062001bc281840162001816565b90820152845292840192810162001b6f565b600082601f83011262001be657600080fd5b8151602062001bf96200191383620018c8565b82815260069290921b8401810191818101908684111562001c1957600080fd5b8286015b8481101562001993576040818903121562001c385760008081fd5b62001c4262001690565b815162001c4f816200170d565b815262001c5e82860162001803565b8186015283529183019160400162001c1d565b60008060008060008060006102a0888a03121562001c8e57600080fd5b62001c9a898962001750565b965062001cab8960e08a016200182b565b6101c08901519096506001600160401b038082111562001cca57600080fd5b62001cd88b838c01620018ee565b965062001cea8b6101e08c01620019c7565b95506102408a015191508082111562001d0257600080fd5b62001d108b838c0162001a1f565b94506102608a015191508082111562001d2857600080fd5b62001d368b838c0162001b26565b93506102808a015191508082111562001d4e57600080fd5b5062001d5d8a828b0162001bd4565b91505092959891949750929550565b82516001600160a01b0390811682526020808501516001600160401b03908116828501526040808701518216818601526060808801518316818701526080808901516001600160601b03168188015260a0808a015187168189015260c0808b01518816818a01528951881660e08a01529589015161ffff9081166101008a01529389015163ffffffff9081166101208a015292890151909316610140880152870151909416610160860152850151909216610180840152830151166101a08201526101c0810162000f4c565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b60006001820162001e795762001e7962001e4e565b5060010190565b602080825282518282018190526000919060409081850190868401855b8281101562001f2b57815180516001600160a01b031685528681015163ffffffff9081168887015286820151811687870152606080830151909116908601526080808201516001600160401b03169086015260a08082015162001f0a878301826001600160401b03169052565b505060c09081015115159085015260e0909301929085019060010162001e9d565b5091979650505050505050565b602080825282518282018190526000919060409081850190868401855b8281101562001f2b57815180516001600160a01b031685528681015161ffff168786015285015163ffffffff16858501526060909301929085019060010162001f55565b8181038181111562000f115762000f1162001e4e565b60008162001fc15762001fc162001e4e565b506000190190565b63ffffffff8181168382160190808211156200146b576200146b62001e4e565b6000604080830163ffffffff8616845260208281860152818651808452606087019150828801935060005b818110156200204857845180516001600160a01b0316845284015161ffff1684840152938301939185019160010162002014565b509098975050505050505050565b6000602082840312156200206957600080fd5b815162000f4c816200170d565b808202811582820484141762000f115762000f1162001e4e565b600082620020ae57634e487b7160e01b600052601260045260246000fd5b500490565b6001600160601b038281168282160390808211156200146b576200146b62001e4e565b600060208284031215620020e957600080fd5b5051919050565b81810360008312801583831316838312821617156200146b576200146b62001e4e565b6000602082840312156200212657600080fd5b62000f4c826200199e565b634e487b7160e01b600052603160045260246000fd5b60005b83811015620021645781810151838201526020016200214a565b50506000910152565b600082516200218181846020870162002147565b9190910192915050565b6020815260008251806020840152620021ac81604085016020870162002147565b601f01601f19169190910160400192915050565b60805160a05160c05160e051610100516101205161014051610160516159c6620022c760003960008181610334015281816113dc015261328a015260008181610305015281816112fe01528181611366015281816118a801528181611910015261326301526000818161027101528181610b4201528181611c9f01526131dd015260008181610241015281816119db01526131b301526000818161021201528181610ea50152818161165b01528181611754015281816121d201528181612e8d0152818161318e015261341a0152600081816102d101528181611820015261322d0152600081816102a101528181612519015261320401526000611bb901526159c66000f3fe608060405234801561001057600080fd5b50600436106101c45760003560e01c80637437ff9f116100f9578063b06d41bc11610097578063d3c7c2c711610071578063d3c7c2c7146108b4578063e76c0b06146108c9578063eff7cc48146108dc578063f2fde38b146108e457600080fd5b8063b06d41bc14610883578063c92b283214610899578063d09dc339146108ac57600080fd5b8063856c8247116100d3578063856c8247146106df5780638da5cb5b146106f25780639a113c3614610703578063a7d3e02f1461087057600080fd5b80637437ff9f146105a857806376f6ae76146106c457806379ba5097146106d757600080fd5b8063546719cd11610166578063599f643111610140578063599f64311461054a5780635d86f1411461056f578063704b6c02146105825780637132721a1461059557600080fd5b8063546719cd146104b3578063549e946f1461051757806354b714681461052a57600080fd5b8063352e4bc8116101a2578063352e4bc81461044957806338724a951461045e5780633a87ac531461047f5780634120fccd1461049257600080fd5b806306285c69146101c95780631772047e1461037a578063181f5a7714610400575b600080fd5b6103646040805160e081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c08101919091526040518060e001604052807f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031681526020017f000000000000000000000000000000000000000000000000000000000000000067ffffffffffffffff1681526020017f000000000000000000000000000000000000000000000000000000000000000067ffffffffffffffff1681526020017f000000000000000000000000000000000000000000000000000000000000000067ffffffffffffffff1681526020017f00000000000000000000000000000000000000000000000000000000000000006bffffffffffffffffffffffff1681526020017f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031681526020017f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316815250905090565b6040516103719190614607565b60405180910390f35b6103db610388366004614698565b6040805180820190915260008082526020820152506001600160a01b03166000908152600e602090815260409182902082518084019093525461ffff8116835262010000900463ffffffff169082015290565b60408051825161ffff16815260209283015163ffffffff169281019290925201610371565b61043c6040518060400160405280601381526020017f45564d3245564d4f6e52616d7020312e322e300000000000000000000000000081525081565b6040516103719190614723565b61045c610457366004614884565b6108f7565b005b61047161046c3660046149ba565b610960565b604051908152602001610371565b61045c61048d366004614a8b565b610d4c565b61049a610d62565b60405167ffffffffffffffff9091168152602001610371565b6104bb610d96565b604051610371919081516fffffffffffffffffffffffffffffffff908116825260208084015163ffffffff1690830152604080840151151590830152606080840151821690830152608092830151169181019190915260a00190565b61045c610525366004614aef565b610e46565b6010546040516bffffffffffffffffffffffff9091168152602001610371565b6002546001600160a01b03165b6040516001600160a01b039091168152602001610371565b61055761057d366004614698565b610ffb565b61045c610590366004614698565b61105a565b61045c6105a3366004614b3a565b611124565b6106b76040805160e081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810191909152506040805160e0810182526005546001600160a01b03808216835261ffff740100000000000000000000000000000000000000008084048216602086015263ffffffff76010000000000000000000000000000000000000000000085048116968601969096527a010000000000000000000000000000000000000000000000000000909304166060840152600654908116608084015290810490921660a082015267ffffffffffffffff78010000000000000000000000000000000000000000000000009092049190911660c082015290565b6040516103719190614bd2565b61045c6106d2366004614c4d565b611135565b61045c6111ed565b61049a6106ed366004614698565b6112d0565b6000546001600160a01b0316610557565b610806610711366004614698565b6040805160c081018252600080825260208201819052918101829052606081018290526080810182905260a0810191909152506001600160a01b03166000908152600d6020908152604091829020825160c081018452905463ffffffff808216835264010000000082048116938301939093526801000000000000000081049092169281019290925267ffffffffffffffff6c0100000000000000000000000082048116606084015274010000000000000000000000000000000000000000820416608083015260ff7c010000000000000000000000000000000000000000000000000000000090910416151560a082015290565b6040516103719190600060c08201905063ffffffff80845116835280602085015116602084015280604085015116604084015250606083015167ffffffffffffffff8082166060850152806080860151166080850152505060a0830151151560a083015292915050565b61047161087e366004614cc2565b6113d8565b61088b611d86565b604051610371929190614d70565b61045c6108a7366004614db2565b611e8a565b610471611ef2565b6108bc611efc565b6040516103719190614e02565b61045c6108d7366004614e4f565b611fad565b61045c612013565b61045c6108f2366004614698565b6122aa565b6000546001600160a01b0316331480159061091d57506002546001600160a01b03163314155b15610954576040517ffbdb8e5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61095d816122bb565b50565b6000806109786109736080850185614f18565b6124ea565b90506109a361098a6020850185614f18565b835190915061099c6040870187614f7d565b90506125de565b6000600d816109b86080870160608801614698565b6001600160a01b031681526020808201929092526040908101600020815160c081018352905463ffffffff80821683526401000000008204811694830194909452680100000000000000008104909316918101919091526c01000000000000000000000000820467ffffffffffffffff90811660608301527401000000000000000000000000000000000000000083041660808201527c010000000000000000000000000000000000000000000000000000000090910460ff16151560a08201819052909150610ad557610a926080850160608601614698565b6040517fa7499d200000000000000000000000000000000000000000000000000000000081526001600160a01b0390911660048201526024015b60405180910390fd5b60065460009081906001600160a01b031663ffdb4b37610afb6080890160608a01614698565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e084901b1681526001600160a01b03909116600482015267ffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001660248201526044016040805180830381865afa158015610b86573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610baa919061500d565b909250905060008080610bc060408a018a614f7d565b90501115610bf957610bef610bdb60808a0160608b01614698565b85610be960408c018c614f7d565b89612703565b9092509050610c15565b8451610c129063ffffffff16662386f26fc1000061506f565b91505b6080850151610c2e9067ffffffffffffffff168361506f565b606086015160055491935060009167ffffffffffffffff9091169063ffffffff8416907a010000000000000000000000000000000000000000000000000000900461ffff16610c8060208d018d614f18565b610c8b92915061506f565b6005548a51610cba91760100000000000000000000000000000000000000000000900463ffffffff1690615086565b610cc49190615086565b610cce9190615086565b610cd8919061506f565b610cfc9077ffffffffffffffffffffffffffffffffffffffffffffffff861661506f565b9050610d3f670de0b6b3a7640000610d148386615086565b610d1e9190615099565b77ffffffffffffffffffffffffffffffffffffffffffffffff871690612989565b9998505050505050505050565b610d546129c2565b610d5e8282612a38565b5050565b601054600090610d9190700100000000000000000000000000000000900467ffffffffffffffff1660016150d4565b905090565b6040805160a0810182526000808252602082018190529181018290526060810182905260808101919091526040805160a0810182526003546fffffffffffffffffffffffffffffffff808216835270010000000000000000000000000000000080830463ffffffff1660208501527401000000000000000000000000000000000000000090920460ff161515938301939093526004548084166060840152049091166080820152610d9190612d98565b6000546001600160a01b03163314801590610e6c57506002546001600160a01b03163314155b15610ea3576040517ffbdb8e5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316826001600160a01b03161480610eea57506001600160a01b038116155b15610f21576040517f232cb97f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000610f2b612e4a565b1215610f63576040517f02075e0000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152610d5e9082906001600160a01b038516906370a0823190602401602060405180830381865afa158015610fc6573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610fea91906150f5565b6001600160a01b0385169190612f0a565b6000611008600a83612f8a565b611049576040517fbf16aab60000000000000000000000000000000000000000000000000000000081526001600160a01b0383166004820152602401610acc565b611054600a83612f9f565b92915050565b6000546001600160a01b0316331480159061108057506002546001600160a01b03163314155b156110b7576040517ff6cd562000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600280547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0383169081179091556040519081527f8fe72c3e0020beb3234e76ae6676fa576fbfcae600af1c4fea44784cf0db329c906020015b60405180910390a150565b61112c6129c2565b61095d81612fb4565b6000546001600160a01b0316331480159061115b57506002546001600160a01b03163314155b15611192576040517ffbdb8e5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610d5e8282808060200260200160405190810160405280939291908181526020016000905b828210156111e3576111d46040830286013681900381019061510e565b815260200190600101906111b7565b50505050506132e1565b6001546001600160a01b03163314611261576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610acc565b60008054337fffffffffffffffffffffffff0000000000000000000000000000000000000000808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6001600160a01b0381166000908152600f602052604081205467ffffffffffffffff168015801561132957507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031615155b15611054576040517f856c82470000000000000000000000000000000000000000000000000000000081526001600160a01b0384811660048301527f0000000000000000000000000000000000000000000000000000000000000000169063856c824790602401602060405180830381865afa1580156113ad573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906113d1919061514d565b9392505050565b60007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663397796f76040518163ffffffff1660e01b8152600401602060405180830381865afa158015611438573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061145c919061516a565b15611493576040517fc148371500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600160a01b0382166114d3576040517fa4ec747900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6005546001600160a01b03163314611517576040517f1c0a352900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6115218480614f18565b9050602014611568576115348480614f18565b6040517f370d875f000000000000000000000000000000000000000000000000000000008152600401610acc9291906151d0565b60006115748580614f18565b81019061158191906151e4565b90506001600160a01b038111806115985750600a81105b156115a7576115348580614f18565b60006115b96109736080880188614f18565b90506115dd6115cb6020880188614f18565b835190915061099c60408a018a614f7d565b6116516115ed6040880188614f7d565b808060200260200160405190810160405280939291908181526020016000905b828210156116395761162a604083028601368190038101906151fd565b8152602001906001019061160d565b50506006546001600160a01b03169250613554915050565b6001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001661168b6080880160608901614698565b6001600160a01b0316036116ef57601080548691906000906116bc9084906bffffffffffffffffffffffff16615237565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555061180e565b6006546001600160a01b03166241e5be61170f6080890160608a01614698565b60405160e083901b7fffffffff000000000000000000000000000000000000000000000000000000001681526001600160a01b039182166004820152602481018990527f00000000000000000000000000000000000000000000000000000000000000009091166044820152606401602060405180830381865afa15801561179b573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906117bf91906150f5565b601080546000906117df9084906bffffffffffffffffffffffff16615237565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055505b6010546bffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000081169116111561187b576040517fe5c7a49100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600160a01b0384166000908152600f602052604090205467ffffffffffffffff161580156118d357507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031615155b156119cb576040517f856c82470000000000000000000000000000000000000000000000000000000081526001600160a01b0385811660048301527f0000000000000000000000000000000000000000000000000000000000000000169063856c824790602401602060405180830381865afa158015611957573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061197b919061514d565b6001600160a01b0385166000908152600f6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001667ffffffffffffffff929092169190911790555b60006040518061018001604052807f000000000000000000000000000000000000000000000000000000000000000067ffffffffffffffff16815260200160108081819054906101000a900467ffffffffffffffff16611a2a9061525c565b825467ffffffffffffffff9182166101009390930a8381029083021990911617909255825260208083018a90526001600160a01b03891660408085018290526000918252600f90925290812080546060909401939092611a8a911661525c565b91906101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905567ffffffffffffffff16815260200183600001518152602001600015158152602001846001600160a01b03168152602001888060200190611af09190614f18565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250505090825250602001611b3760408a018a614f7d565b808060200260200160405190810160405280939291908181526020016000905b82821015611b8357611b74604083028601368190038101906151fd565b81526020019060010190611b57565b5050509183525050602001611b9e60808a0160608b01614698565b6001600160a01b0316815260006020909101529050611bdd817f0000000000000000000000000000000000000000000000000000000000000000613707565b61016082015260005b611bf36040890189614f7d565b9050811015611d3f576000611c0b60408a018a614f7d565b83818110611c1b57611c1b615283565b905060400201803603810190611c3191906151fd565b9050611c408160000151610ffb565b6001600160a01b0316639687544588611c598c80614f18565b60208087015160408051928301815260008352517fffffffff0000000000000000000000000000000000000000000000000000000060e088901b168152611cc7959493927f0000000000000000000000000000000000000000000000000000000000000000916004016152b2565b6000604051808303816000875af1158015611ce6573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052611d2c91908101906152fd565b505080611d38906153af565b9050611be6565b507faffc45517195d6499808c643bd4a7b0ffeedf95bea5852840d7bfcf63f59e82181604051611d6f919061542b565b60405180910390a161016001519695505050505050565b6060600080611d956007613811565b90508067ffffffffffffffff811115611db057611db0614736565b604051908082528060200260200182016040528015611df557816020015b6040805180820190915260008082526020820152815260200190600190039081611dce5790505b50925060005b81811015611e6757600080611e1160078461381c565b915091506040518060400160405280836001600160a01b031681526020018261ffff16815250868481518110611e4957611e49615283565b6020026020010181905250505080611e60906153af565b9050611dfb565b505060105491926c0100000000000000000000000090920463ffffffff16919050565b6000546001600160a01b03163314801590611eb057506002546001600160a01b03163314155b15611ee7576040517ff6cd562000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61095d60038261383a565b6000610d91612e4a565b60606000611f0a600a613a12565b67ffffffffffffffff811115611f2257611f22614736565b604051908082528060200260200182016040528015611f4b578160200160208202803683370190505b50905060005b8151811015611fa757611f65600a82613a1d565b50828281518110611f7857611f78615283565b60200260200101816001600160a01b03166001600160a01b03168152505080611fa0906153af565b9050611f51565b50919050565b6000546001600160a01b03163314801590611fd357506002546001600160a01b03163314155b1561200a576040517ffbdb8e5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61095d81613a2c565b6000546001600160a01b0316331480159061203957506002546001600160a01b03163314155b801561204d575061204b600733613b13565b155b15612084576040517f195db95800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6010546c01000000000000000000000000900463ffffffff1660008190036120d8576040517f990e30bf00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6010546bffffffffffffffffffffffff1681811015612123576040517f8d0f71d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600061212d612e4a565b1215612165576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8060006121726007613811565b905060005b818110156122675760008061218d60078461381c565b90925090506000876121ad836bffffffffffffffffffffffff8a1661506f565b6121b79190615099565b90506121c38187615568565b95506122076001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016846bffffffffffffffffffffffff8416612f0a565b6040516bffffffffffffffffffffffff821681526001600160a01b038416907f55fdec2aab60a41fa5abb106670eb1006f5aeaee1ba7afea2bc89b5b3ec7678f9060200160405180910390a250505080612260906153af565b9050612177565b5050601080547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff929092169190911790555050565b6122b26129c2565b61095d81613b28565b60005b81518110156124ba5760008282815181106122db576122db615283565b6020908102919091018101516040805160c080820183528385015163ffffffff9081168352838501518116838701908152606080870151831685870190815260808089015167ffffffffffffffff90811693880193845260a0808b01518216928901928352968a0151151596880196875298516001600160a01b03166000908152600d909a5296909820945185549251985191519651945115157c0100000000000000000000000000000000000000000000000000000000027fffffff00ffffffffffffffffffffffffffffffffffffffffffffffffffffffff9589167401000000000000000000000000000000000000000002959095167fffffff000000000000000000ffffffffffffffffffffffffffffffffffffffff979098166c01000000000000000000000000027fffffffffffffffffffffffff0000000000000000ffffffffffffffffffffffff9285166801000000000000000002929092167fffffffffffffffffffffffff000000000000000000000000ffffffffffffffff998516640100000000027fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000009094169190941617919091179690961617949094179190911691909117919091179055506124b3816153af565b90506122be565b507f2386f61ab5cafc3fed44f9f614f721ab53479ef64067fd16c1a2491b63ddf1a881604051611119919061558d565b60408051602081019091526000815260008290036125405750604080516020810190915267ffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000168152611054565b7f97a657c90000000000000000000000000000000000000000000000000000000061256b8385615644565b7fffffffff0000000000000000000000000000000000000000000000000000000016146125c4576040517f5247fdce00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6125d1826004818661568c565b8101906113d191906156b6565b60065474010000000000000000000000000000000000000000900463ffffffff1680841115612643576040517f869337890000000000000000000000000000000000000000000000000000000081526004810182905260248101859052604401610acc565b6006547801000000000000000000000000000000000000000000000000900467ffffffffffffffff168311156126a5576040517f4c4fc93a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60055474010000000000000000000000000000000000000000900461ffff168211156126fd576040517f4c056b6a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b50505050565b60008083815b8181101561291557600087878381811061272557612725615283565b90506040020180360381019061273b91906151fd565b80516001600160a01b03166000908152600e602090815260409182902082518084019093525461ffff8116835263ffffffff620100009091048116918301919091528251929350909161279291600a9190612f8a16565b6127d65781516040517fbf16aab60000000000000000000000000000000000000000000000000000000081526001600160a01b039091166004820152602401610acc565b805160009061ffff16156128e35760008c6001600160a01b031684600001516001600160a01b0316146128925760065484516040517f4ab35b0b0000000000000000000000000000000000000000000000000000000081526001600160a01b039182166004820152911690634ab35b0b90602401602060405180830381865afa158015612867573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061288b91906156f8565b9050612895565b508a5b82516020850151620186a09161ffff16906128cb9077ffffffffffffffffffffffffffffffffffffffffffffffff851690613c03565b6128d5919061506f565b6128df9190615099565b9150505b6128ed8188615086565b96508160200151866128ff9190615713565b95505050508061290e906153af565b9050612709565b506000846020015163ffffffff16662386f26fc10000612935919061506f565b90508084101561294857925061297f9050565b6000856040015163ffffffff16662386f26fc10000612967919061506f565b90508085111561297b57935061297f915050565b5050505b9550959350505050565b600077ffffffffffffffffffffffffffffffffffffffffffffffff83166129b883670de0b6b3a764000061506f565b6113d19190615099565b6000546001600160a01b03163314612a36576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610acc565b565b60005b8251811015612b99576000838281518110612a5857612a58615283565b60200260200101516000015190506000848381518110612a7a57612a7a615283565b6020026020010151602001519050612a9c82600a612f8a90919063ffffffff16565b612add576040517f73913ebd0000000000000000000000000000000000000000000000000000000081526001600160a01b0383166004820152602401610acc565b6001600160a01b038116612af2600a84612f9f565b6001600160a01b031614612b32576040517f6cc7b99800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b612b3d600a83613c32565b15612b8657604080516001600160a01b038085168252831660208201527f987eb3c2f78454541205f72f34839b434c306c9eaf4922efd7c0c3060fdb2e4c910160405180910390a15b505080612b92906153af565b9050612a3b565b5060005b8151811015612d93576000828281518110612bba57612bba615283565b60200260200101516000015190506000838381518110612bdc57612bdc615283565b602002602001015160200151905060006001600160a01b0316826001600160a01b03161480612c1257506001600160a01b038116155b15612c49576040517f6c2a418000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b806001600160a01b03166321df0da76040518163ffffffff1660e01b8152600401602060405180830381865afa158015612c87573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612cab9190615730565b6001600160a01b0316826001600160a01b031614612cf5576040517f6cc7b99800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b612d01600a8383613c47565b15612d4e57604080516001600160a01b038085168252831660208201527f95f865c2808f8b2a85eea2611db7843150ee7835ef1403f9755918a97d76933c910160405180910390a1612d80565b6040517f3caf458500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b505080612d8c906153af565b9050612b9d565b505050565b6040805160a081018252600080825260208201819052918101829052606081018290526080810191909152612e2682606001516fffffffffffffffffffffffffffffffff1683600001516fffffffffffffffffffffffffffffffff16846020015163ffffffff1642612e0a919061574d565b85608001516fffffffffffffffffffffffffffffffff16613c65565b6fffffffffffffffffffffffffffffffff1682525063ffffffff4216602082015290565b6010546040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526000916bffffffffffffffffffffffff16907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316906370a0823190602401602060405180830381865afa158015612edc573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612f0091906150f5565b610d919190615760565b604080516001600160a01b038416602482015260448082018490528251808303909101815260649091019091526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fa9059cbb00000000000000000000000000000000000000000000000000000000179052612d93908490613c8d565b60006113d1836001600160a01b038416613d8c565b60006113d1836001600160a01b038416613d98565b60808101516001600160a01b0316612ff8576040517f35be3ac800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8051600580546020808501516040808701516060808901516001600160a01b039889167fffffffffffffffffffff00000000000000000000000000000000000000000000909716969096177401000000000000000000000000000000000000000061ffff9586168102919091177fffffffff000000000000ffffffffffffffffffffffffffffffffffffffffffff1676010000000000000000000000000000000000000000000063ffffffff948516027fffffffff0000ffffffffffffffffffffffffffffffffffffffffffffffffffff16177a0100000000000000000000000000000000000000000000000000009590971694909402959095179095556080808801516006805460a0808c015160c0808e0151958d167fffffffffffffffff000000000000000000000000000000000000000000000000909416939093179a169096029890981777ffffffffffffffffffffffffffffffffffffffffffffffff16780100000000000000000000000000000000000000000000000067ffffffffffffffff93841602179055825160e0810184527f0000000000000000000000000000000000000000000000000000000000000000891681527f00000000000000000000000000000000000000000000000000000000000000008216958101959095527f00000000000000000000000000000000000000000000000000000000000000008116858401527f000000000000000000000000000000000000000000000000000000000000000016948401949094527f00000000000000000000000000000000000000000000000000000000000000006bffffffffffffffffffffffff16938301939093527f00000000000000000000000000000000000000000000000000000000000000008516908201527f000000000000000000000000000000000000000000000000000000000000000090931691830191909152517f72c6aaba4dde02f77d291123a76185c418ba63f8c217a2d56b08aec84e9bbfb891611119918490615780565b8051604081111561331e576040517fb5a10cfa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6010546c01000000000000000000000000900463ffffffff161580159061336c575060105463ffffffff6c010000000000000000000000008204166bffffffffffffffffffffffff90911610155b1561337957613379612013565b60006133856007613811565b90505b80156133c75760006133a661339e60018461574d565b60079061381c565b5090506133b4600782613da4565b5050806133c090615875565b9050613388565b506000805b828110156134d55760008482815181106133e8576133e8615283565b6020026020010151600001519050600085838151811061340a5761340a615283565b60200260200101516020015190507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316826001600160a01b0316148061345f57506001600160a01b038216155b156134a1576040517f4de938d10000000000000000000000000000000000000000000000000000000081526001600160a01b0383166004820152602401610acc565b6134b160078361ffff8416613db9565b506134c061ffff821685615713565b93505050806134ce906153af565b90506133cc565b50601080547fffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffffff166c0100000000000000000000000063ffffffff8416021790556040517f8c337bff38141c507abd25c547606bdde78fe8c12e941ab613f3a565fea6cd249061354790839086906158aa565b60405180910390a1505050565b81516000805b828110156136f9576000846001600160a01b031663d02641a087848151811061358557613585615283565b6020908102919091010151516040517fffffffff0000000000000000000000000000000000000000000000000000000060e084901b1681526001600160a01b0390911660048201526024016040805180830381865afa1580156135ec573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061361091906158c9565b51905077ffffffffffffffffffffffffffffffffffffffffffffffff81166000036136915785828151811061364757613647615283565b6020908102919091010151516040517f9a655f7b0000000000000000000000000000000000000000000000000000000081526001600160a01b039091166004820152602401610acc565b6136db8683815181106136a6576136a6615283565b6020026020010151602001518277ffffffffffffffffffffffffffffffffffffffffffffffff16613c0390919063ffffffff16565b6136e59084615086565b925050806136f2906153af565b905061355a565b506126fd6003826000613dcf565b60008060001b828460200151856080015186606001518760e001518861010001518051906020012089610120015160405160200161374591906158fc565b604051602081830303815290604052805190602001208a60a001518b60c001518c61014001518d604001516040516020016137f39c9b9a999897969594939291909b8c5260208c019a909a5267ffffffffffffffff98891660408c01529690971660608a01526001600160a01b0394851660808a015292841660a089015260c088019190915260e0870152610100860152911515610120850152166101408301526101608201526101800190565b60405160208183030381529060405280519060200120905092915050565b60006110548261411e565b600080808061382b8686614129565b909450925050505b9250929050565b815460009061386390700100000000000000000000000000000000900463ffffffff164261574d565b9050801561390557600183015483546138ab916fffffffffffffffffffffffffffffffff80821692811691859170010000000000000000000000000000000090910416613c65565b83546fffffffffffffffffffffffffffffffff919091167fffffffffffffffffffffffff0000000000000000000000000000000000000000909116177001000000000000000000000000000000004263ffffffff16021783555b6020820151835461392b916fffffffffffffffffffffffffffffffff9081169116614154565b83548351151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00ffffffff000000000000000000000000000000009091166fffffffffffffffffffffffffffffffff92831617178455602083015160408085015183167001000000000000000000000000000000000291909216176001850155517f9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c19906135479084908151151581526020808301516fffffffffffffffffffffffffffffffff90811691830191909152604092830151169181019190915260600190565b600061105482613811565b600080808061382b868661381c565b60005b8151811015613ae3576000828281518110613a4c57613a4c615283565b6020908102919091018101516040805180820182528284015161ffff90811682528284015163ffffffff90811683870190815294516001600160a01b03166000908152600e9096529290942090518154935190921662010000027fffffffffffffffffffffffffffffffffffffffffffffffffffff000000000000909316919093161717905550613adc816153af565b9050613a2f565b507f4230c60a9725eb5fb992cf6a215398b4e81b4606d4a1e6be8dfe0b60dc172ec181604051611119919061590f565b60006113d1836001600160a01b03841661416a565b336001600160a01b03821603613b9a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610acc565b600180547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000670de0b6b3a76400006129b88377ffffffffffffffffffffffffffffffffffffffffffffffff861661506f565b60006113d1836001600160a01b038416614176565b6000613c5d846001600160a01b03851684614182565b949350505050565b6000613c8485613c75848661506f565b613c7f9087615086565b614154565b95945050505050565b6000613ce2826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c6564815250856001600160a01b03166141989092919063ffffffff16565b805190915015612d935780806020019051810190613d00919061516a565b612d93576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f742073756363656564000000000000000000000000000000000000000000006064820152608401610acc565b60006113d1838361416a565b60006113d183836141a7565b60006113d1836001600160a01b038416614231565b6000613c5d846001600160a01b0385168461424e565b825474010000000000000000000000000000000000000000900460ff161580613df6575081155b15613e0057505050565b825460018401546fffffffffffffffffffffffffffffffff80831692911690600090613e4690700100000000000000000000000000000000900463ffffffff164261574d565b90508015613f065781831115613e88576040517f9725942a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001860154613ec29083908590849070010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff16613c65565b86547fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff167001000000000000000000000000000000004263ffffffff160217875592505b84821015613fa3576001600160a01b038416613f58576040517ff94ebcd10000000000000000000000000000000000000000000000000000000081526004810183905260248101869052604401610acc565b6040517f1a76572a00000000000000000000000000000000000000000000000000000000815260048101839052602481018690526001600160a01b0385166044820152606401610acc565b8483101561409c5760018681015470010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff16906000908290613fe7908261574d565b613ff1878a61574d565b613ffb9190615086565b6140059190615099565b90506001600160a01b038616614051576040517f15279c080000000000000000000000000000000000000000000000000000000081526004810182905260248101869052604401610acc565b6040517fd0c8d23a00000000000000000000000000000000000000000000000000000000815260048101829052602481018690526001600160a01b0387166044820152606401610acc565b6140a6858461574d565b86547fffffffffffffffffffffffffffffffff00000000000000000000000000000000166fffffffffffffffffffffffffffffffff82161787556040518681529093507f1871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690a9060200160405180910390a1505050505050565b60006110548261426b565b600080806141378585614275565b600081815260029690960160205260409095205494959350505050565b600081831061416357816113d1565b5090919050565b60006113d18383614281565b60006113d18383614231565b6000613c5d84846001600160a01b03851661424e565b6060613c5d8484600085614299565b6000818152600283016020526040812054801515806141cb57506141cb848461416a565b6113d1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f456e756d657261626c654d61703a206e6f6e6578697374656e74206b657900006044820152606401610acc565b600081815260028301602052604081208190556113d183836143a5565b60008281526002840160205260408120829055613c5d84846143b1565b6000611054825490565b60006113d183836143bd565b600081815260018301602052604081205415156113d1565b60608247101561432b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f60448201527f722063616c6c00000000000000000000000000000000000000000000000000006064820152608401610acc565b600080866001600160a01b03168587604051614347919061596e565b60006040518083038185875af1925050503d8060008114614384576040519150601f19603f3d011682016040523d82523d6000602084013e614389565b606091505b509150915061439a878383876143e7565b979650505050505050565b60006113d1838361447a565b60006113d18383614574565b60008260000182815481106143d4576143d4615283565b9060005260206000200154905092915050565b60608315614470578251600003614469576001600160a01b0385163b614469576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401610acc565b5081613c5d565b613c5d83836145c3565b6000818152600183016020526040812054801561456357600061449e60018361574d565b85549091506000906144b29060019061574d565b90508181146145175760008660000182815481106144d2576144d2615283565b90600052602060002001549050808760000184815481106144f5576144f5615283565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806145285761452861598a565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050611054565b6000915050611054565b5092915050565b60008181526001830160205260408120546145bb57508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155611054565b506000611054565b8151156145d35781518083602001fd5b806040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610acc9190614723565b60e0810161105482846001600160a01b03808251168352602082015167ffffffffffffffff808216602086015280604085015116604086015280606085015116606086015250506bffffffffffffffffffffffff60808301511660808401528060a08301511660a08401528060c08301511660c0840152505050565b6001600160a01b038116811461095d57600080fd5b6000602082840312156146aa57600080fd5b81356113d181614683565b60005b838110156146d05781810151838201526020016146b8565b50506000910152565b600081518084526146f18160208601602086016146b5565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6020815260006113d160208301846146d9565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160e0810167ffffffffffffffff8111828210171561478857614788614736565b60405290565b6040805190810167ffffffffffffffff8111828210171561478857614788614736565b6040516060810167ffffffffffffffff8111828210171561478857614788614736565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561481b5761481b614736565b604052919050565b600067ffffffffffffffff82111561483d5761483d614736565b5060051b60200190565b803563ffffffff8116811461485b57600080fd5b919050565b67ffffffffffffffff8116811461095d57600080fd5b801515811461095d57600080fd5b6000602080838503121561489757600080fd5b823567ffffffffffffffff8111156148ae57600080fd5b8301601f810185136148bf57600080fd5b80356148d26148cd82614823565b6147d4565b81815260e091820283018401918482019190888411156148f157600080fd5b938501935b8385101561499c5780858a03121561490e5760008081fd5b614916614765565b853561492181614683565b815261492e868801614847565b87820152604061493f818801614847565b908201526060614950878201614847565b9082015260808681013561496381614860565b9082015260a08681013561497681614860565b9082015260c08681013561498981614876565b90820152835293840193918501916148f6565b50979650505050505050565b600060a08284031215611fa757600080fd5b6000602082840312156149cc57600080fd5b813567ffffffffffffffff8111156149e357600080fd5b613c5d848285016149a8565b600082601f830112614a0057600080fd5b81356020614a106148cd83614823565b82815260069290921b84018101918181019086841115614a2f57600080fd5b8286015b84811015614a805760408189031215614a4c5760008081fd5b614a5461478e565b8135614a5f81614683565b815281850135614a6e81614683565b81860152835291830191604001614a33565b509695505050505050565b60008060408385031215614a9e57600080fd5b823567ffffffffffffffff80821115614ab657600080fd5b614ac2868387016149ef565b93506020850135915080821115614ad857600080fd5b50614ae5858286016149ef565b9150509250929050565b60008060408385031215614b0257600080fd5b8235614b0d81614683565b91506020830135614b1d81614683565b809150509250929050565b803561ffff8116811461485b57600080fd5b600060e08284031215614b4c57600080fd5b614b54614765565b8235614b5f81614683565b8152614b6d60208401614b28565b6020820152614b7e60408401614847565b6040820152614b8f60608401614b28565b60608201526080830135614ba281614683565b6080820152614bb360a08401614847565b60a082015260c0830135614bc681614860565b60c08201529392505050565b60e0810161105482846001600160a01b03808251168352602082015161ffff80821660208601526040840151915063ffffffff80831660408701528160608601511660608701528360808601511660808701528060a08601511660a08701525050505067ffffffffffffffff60c08201511660c08301525050565b60008060208385031215614c6057600080fd5b823567ffffffffffffffff80821115614c7857600080fd5b818501915085601f830112614c8c57600080fd5b813581811115614c9b57600080fd5b8660208260061b8501011115614cb057600080fd5b60209290920196919550909350505050565b600080600060608486031215614cd757600080fd5b833567ffffffffffffffff811115614cee57600080fd5b614cfa868287016149a8565b935050602084013591506040840135614d1281614683565b809150509250925092565b600081518084526020808501945080840160005b83811015614d6557815180516001600160a01b0316885283015161ffff168388015260409096019590820190600101614d31565b509495945050505050565b604081526000614d836040830185614d1d565b90508260208301529392505050565b80356fffffffffffffffffffffffffffffffff8116811461485b57600080fd5b600060608284031215614dc457600080fd5b614dcc6147b1565b8235614dd781614876565b8152614de560208401614d92565b6020820152614df660408401614d92565b60408201529392505050565b6020808252825182820181905260009190848201906040850190845b81811015614e435783516001600160a01b031683529284019291840191600101614e1e565b50909695505050505050565b60006020808385031215614e6257600080fd5b823567ffffffffffffffff811115614e7957600080fd5b8301601f81018513614e8a57600080fd5b8035614e986148cd82614823565b81815260609182028301840191848201919088841115614eb757600080fd5b938501935b8385101561499c5780858a031215614ed45760008081fd5b614edc6147b1565b8535614ee781614683565b8152614ef4868801614b28565b878201526040614f05818801614847565b9082015283529384019391850191614ebc565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112614f4d57600080fd5b83018035915067ffffffffffffffff821115614f6857600080fd5b60200191503681900382131561383357600080fd5b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112614fb257600080fd5b83018035915067ffffffffffffffff821115614fcd57600080fd5b6020019150600681901b360382131561383357600080fd5b805177ffffffffffffffffffffffffffffffffffffffffffffffff8116811461485b57600080fd5b6000806040838503121561502057600080fd5b61502983614fe5565b915061503760208401614fe5565b90509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b808202811582820484141761105457611054615040565b8082018082111561105457611054615040565b6000826150cf577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b67ffffffffffffffff81811683821601908082111561456d5761456d615040565b60006020828403121561510757600080fd5b5051919050565b60006040828403121561512057600080fd5b61512861478e565b823561513381614683565b815261514160208401614b28565b60208201529392505050565b60006020828403121561515f57600080fd5b81516113d181614860565b60006020828403121561517c57600080fd5b81516113d181614876565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b602081526000613c5d602083018486615187565b6000602082840312156151f657600080fd5b5035919050565b60006040828403121561520f57600080fd5b61521761478e565b823561522281614683565b81526020928301359281019290925250919050565b6bffffffffffffffffffffffff81811683821601908082111561456d5761456d615040565b600067ffffffffffffffff80831681810361527957615279615040565b6001019392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b6001600160a01b038716815260a0602082015260006152d560a083018789615187565b85604084015267ffffffffffffffff851660608401528281036080840152610d3f81856146d9565b60006020828403121561530f57600080fd5b815167ffffffffffffffff8082111561532757600080fd5b818401915084601f83011261533b57600080fd5b81518181111561534d5761534d614736565b61537e60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116016147d4565b915080825285602082850101111561539557600080fd5b6153a68160208401602086016146b5565b50949350505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036153e0576153e0615040565b5060010190565b600081518084526020808501945080840160005b83811015614d6557815180516001600160a01b0316885283015183880152604090960195908201906001016153fb565b6020815261544660208201835167ffffffffffffffff169052565b60006020830151615463604084018267ffffffffffffffff169052565b5060408301516060830152606083015161548860808401826001600160a01b03169052565b50608083015167ffffffffffffffff811660a08401525060a083015160c083015260c08301516154bc60e084018215159052565b5060e08301516101006154d9818501836001600160a01b03169052565b8085015191505061018061012081818601526154f96101a08601846146d9565b92508086015190506101407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0868503018187015261553784836153e7565b935080870151915050610160615557818701836001600160a01b03169052565b959095015193019290925250919050565b6bffffffffffffffffffffffff82811682821603908082111561456d5761456d615040565b602080825282518282018190526000919060409081850190868401855b8281101561563757815180516001600160a01b031685528681015163ffffffff90811688870152868201518116878701526060808301519091169086015260808082015167ffffffffffffffff169086015260a0808201516156178288018267ffffffffffffffff169052565b505060c09081015115159085015260e090930192908501906001016155aa565b5091979650505050505050565b7fffffffff0000000000000000000000000000000000000000000000000000000081358181169160048510156156845780818660040360031b1b83161692505b505092915050565b6000808585111561569c57600080fd5b838611156156a957600080fd5b5050820193919092039150565b6000602082840312156156c857600080fd5b6040516020810181811067ffffffffffffffff821117156156eb576156eb614736565b6040529135825250919050565b60006020828403121561570a57600080fd5b6113d182614fe5565b63ffffffff81811683821601908082111561456d5761456d615040565b60006020828403121561574257600080fd5b81516113d181614683565b8181038181111561105457611054615040565b818103600083128015838313168383128216171561456d5761456d615040565b6101c081016157fd82856001600160a01b03808251168352602082015167ffffffffffffffff808216602086015280604085015116604086015280606085015116606086015250506bffffffffffffffffffffffff60808301511660808401528060a08301511660a08401528060c08301511660c0840152505050565b82516001600160a01b0390811660e0840152602084015161ffff908116610100850152604085015163ffffffff9081166101208601526060860151909116610140850152608085015190911661016084015260a08401511661018083015260c083015167ffffffffffffffff166101a08301526113d1565b60008161588457615884615040565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190565b63ffffffff83168152604060208201526000613c5d6040830184614d1d565b6000604082840312156158db57600080fd5b6158e361478e565b6158ec83614fe5565b8152602083015161514181614860565b6020815260006113d160208301846153e7565b602080825282518282018190526000919060409081850190868401855b8281101561563757815180516001600160a01b031685528681015161ffff168786015285015163ffffffff16858501526060909301929085019060010161592c565b600082516159808184602087016146b5565b9190910192915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea164736f6c6343000813000a",
}

var EVM2EVMOnRampABI = EVM2EVMOnRampMetaData.ABI

var EVM2EVMOnRampBin = EVM2EVMOnRampMetaData.Bin

func DeployEVM2EVMOnRamp(auth *bind.TransactOpts, backend bind.ContractBackend, staticConfig EVM2EVMOnRampStaticConfig, dynamicConfig EVM2EVMOnRampDynamicConfig, tokensAndPools []InternalPoolUpdate, rateLimiterConfig RateLimiterConfig, feeTokenConfigs []EVM2EVMOnRampFeeTokenConfigArgs, tokenTransferFeeConfigArgs []EVM2EVMOnRampTokenTransferFeeConfigArgs, nopsAndWeights []EVM2EVMOnRampNopAndWeight) (common.Address, *types.Transaction, *EVM2EVMOnRamp, error) {
	parsed, err := EVM2EVMOnRampMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(EVM2EVMOnRampBin), backend, staticConfig, dynamicConfig, tokensAndPools, rateLimiterConfig, feeTokenConfigs, tokenTransferFeeConfigArgs, nopsAndWeights)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &EVM2EVMOnRamp{EVM2EVMOnRampCaller: EVM2EVMOnRampCaller{contract: contract}, EVM2EVMOnRampTransactor: EVM2EVMOnRampTransactor{contract: contract}, EVM2EVMOnRampFilterer: EVM2EVMOnRampFilterer{contract: contract}}, nil
}

type EVM2EVMOnRamp struct {
	address common.Address
	abi     abi.ABI
	EVM2EVMOnRampCaller
	EVM2EVMOnRampTransactor
	EVM2EVMOnRampFilterer
}

type EVM2EVMOnRampCaller struct {
	contract *bind.BoundContract
}

type EVM2EVMOnRampTransactor struct {
	contract *bind.BoundContract
}

type EVM2EVMOnRampFilterer struct {
	contract *bind.BoundContract
}

type EVM2EVMOnRampSession struct {
	Contract     *EVM2EVMOnRamp
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type EVM2EVMOnRampCallerSession struct {
	Contract *EVM2EVMOnRampCaller
	CallOpts bind.CallOpts
}

type EVM2EVMOnRampTransactorSession struct {
	Contract     *EVM2EVMOnRampTransactor
	TransactOpts bind.TransactOpts
}

type EVM2EVMOnRampRaw struct {
	Contract *EVM2EVMOnRamp
}

type EVM2EVMOnRampCallerRaw struct {
	Contract *EVM2EVMOnRampCaller
}

type EVM2EVMOnRampTransactorRaw struct {
	Contract *EVM2EVMOnRampTransactor
}

func NewEVM2EVMOnRamp(address common.Address, backend bind.ContractBackend) (*EVM2EVMOnRamp, error) {
	abi, err := abi.JSON(strings.NewReader(EVM2EVMOnRampABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindEVM2EVMOnRamp(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EVM2EVMOnRamp{address: address, abi: abi, EVM2EVMOnRampCaller: EVM2EVMOnRampCaller{contract: contract}, EVM2EVMOnRampTransactor: EVM2EVMOnRampTransactor{contract: contract}, EVM2EVMOnRampFilterer: EVM2EVMOnRampFilterer{contract: contract}}, nil
}

func NewEVM2EVMOnRampCaller(address common.Address, caller bind.ContractCaller) (*EVM2EVMOnRampCaller, error) {
	contract, err := bindEVM2EVMOnRamp(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EVM2EVMOnRampCaller{contract: contract}, nil
}

func NewEVM2EVMOnRampTransactor(address common.Address, transactor bind.ContractTransactor) (*EVM2EVMOnRampTransactor, error) {
	contract, err := bindEVM2EVMOnRamp(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EVM2EVMOnRampTransactor{contract: contract}, nil
}

func NewEVM2EVMOnRampFilterer(address common.Address, filterer bind.ContractFilterer) (*EVM2EVMOnRampFilterer, error) {
	contract, err := bindEVM2EVMOnRamp(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EVM2EVMOnRampFilterer{contract: contract}, nil
}

func bindEVM2EVMOnRamp(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := EVM2EVMOnRampMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EVM2EVMOnRamp.Contract.EVM2EVMOnRampCaller.contract.Call(opts, result, method, params...)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EVM2EVMOnRamp.Contract.EVM2EVMOnRampTransactor.contract.Transfer(opts)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EVM2EVMOnRamp.Contract.EVM2EVMOnRampTransactor.contract.Transact(opts, method, params...)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EVM2EVMOnRamp.Contract.contract.Call(opts, result, method, params...)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EVM2EVMOnRamp.Contract.contract.Transfer(opts)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EVM2EVMOnRamp.Contract.contract.Transact(opts, method, params...)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampCaller) CurrentRateLimiterState(opts *bind.CallOpts) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _EVM2EVMOnRamp.contract.Call(opts, &out, "currentRateLimiterState")

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_EVM2EVMOnRamp *EVM2EVMOnRampSession) CurrentRateLimiterState() (RateLimiterTokenBucket, error) {
	return _EVM2EVMOnRamp.Contract.CurrentRateLimiterState(&_EVM2EVMOnRamp.CallOpts)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampCallerSession) CurrentRateLimiterState() (RateLimiterTokenBucket, error) {
	return _EVM2EVMOnRamp.Contract.CurrentRateLimiterState(&_EVM2EVMOnRamp.CallOpts)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampCaller) GetDynamicConfig(opts *bind.CallOpts) (EVM2EVMOnRampDynamicConfig, error) {
	var out []interface{}
	err := _EVM2EVMOnRamp.contract.Call(opts, &out, "getDynamicConfig")

	if err != nil {
		return *new(EVM2EVMOnRampDynamicConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(EVM2EVMOnRampDynamicConfig)).(*EVM2EVMOnRampDynamicConfig)

	return out0, err

}

func (_EVM2EVMOnRamp *EVM2EVMOnRampSession) GetDynamicConfig() (EVM2EVMOnRampDynamicConfig, error) {
	return _EVM2EVMOnRamp.Contract.GetDynamicConfig(&_EVM2EVMOnRamp.CallOpts)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampCallerSession) GetDynamicConfig() (EVM2EVMOnRampDynamicConfig, error) {
	return _EVM2EVMOnRamp.Contract.GetDynamicConfig(&_EVM2EVMOnRamp.CallOpts)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampCaller) GetExpectedNextSequenceNumber(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _EVM2EVMOnRamp.contract.Call(opts, &out, "getExpectedNextSequenceNumber")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_EVM2EVMOnRamp *EVM2EVMOnRampSession) GetExpectedNextSequenceNumber() (uint64, error) {
	return _EVM2EVMOnRamp.Contract.GetExpectedNextSequenceNumber(&_EVM2EVMOnRamp.CallOpts)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampCallerSession) GetExpectedNextSequenceNumber() (uint64, error) {
	return _EVM2EVMOnRamp.Contract.GetExpectedNextSequenceNumber(&_EVM2EVMOnRamp.CallOpts)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampCaller) GetFee(opts *bind.CallOpts, message ClientEVM2AnyMessage) (*big.Int, error) {
	var out []interface{}
	err := _EVM2EVMOnRamp.contract.Call(opts, &out, "getFee", message)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_EVM2EVMOnRamp *EVM2EVMOnRampSession) GetFee(message ClientEVM2AnyMessage) (*big.Int, error) {
	return _EVM2EVMOnRamp.Contract.GetFee(&_EVM2EVMOnRamp.CallOpts, message)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampCallerSession) GetFee(message ClientEVM2AnyMessage) (*big.Int, error) {
	return _EVM2EVMOnRamp.Contract.GetFee(&_EVM2EVMOnRamp.CallOpts, message)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampCaller) GetFeeTokenConfig(opts *bind.CallOpts, token common.Address) (EVM2EVMOnRampFeeTokenConfig, error) {
	var out []interface{}
	err := _EVM2EVMOnRamp.contract.Call(opts, &out, "getFeeTokenConfig", token)

	if err != nil {
		return *new(EVM2EVMOnRampFeeTokenConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(EVM2EVMOnRampFeeTokenConfig)).(*EVM2EVMOnRampFeeTokenConfig)

	return out0, err

}

func (_EVM2EVMOnRamp *EVM2EVMOnRampSession) GetFeeTokenConfig(token common.Address) (EVM2EVMOnRampFeeTokenConfig, error) {
	return _EVM2EVMOnRamp.Contract.GetFeeTokenConfig(&_EVM2EVMOnRamp.CallOpts, token)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampCallerSession) GetFeeTokenConfig(token common.Address) (EVM2EVMOnRampFeeTokenConfig, error) {
	return _EVM2EVMOnRamp.Contract.GetFeeTokenConfig(&_EVM2EVMOnRamp.CallOpts, token)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampCaller) GetNopFeesJuels(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _EVM2EVMOnRamp.contract.Call(opts, &out, "getNopFeesJuels")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_EVM2EVMOnRamp *EVM2EVMOnRampSession) GetNopFeesJuels() (*big.Int, error) {
	return _EVM2EVMOnRamp.Contract.GetNopFeesJuels(&_EVM2EVMOnRamp.CallOpts)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampCallerSession) GetNopFeesJuels() (*big.Int, error) {
	return _EVM2EVMOnRamp.Contract.GetNopFeesJuels(&_EVM2EVMOnRamp.CallOpts)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampCaller) GetNops(opts *bind.CallOpts) (GetNops,

	error) {
	var out []interface{}
	err := _EVM2EVMOnRamp.contract.Call(opts, &out, "getNops")

	outstruct := new(GetNops)
	if err != nil {
		return *outstruct, err
	}

	outstruct.NopsAndWeights = *abi.ConvertType(out[0], new([]EVM2EVMOnRampNopAndWeight)).(*[]EVM2EVMOnRampNopAndWeight)
	outstruct.WeightsTotal = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_EVM2EVMOnRamp *EVM2EVMOnRampSession) GetNops() (GetNops,

	error) {
	return _EVM2EVMOnRamp.Contract.GetNops(&_EVM2EVMOnRamp.CallOpts)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampCallerSession) GetNops() (GetNops,

	error) {
	return _EVM2EVMOnRamp.Contract.GetNops(&_EVM2EVMOnRamp.CallOpts)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampCaller) GetPoolBySourceToken(opts *bind.CallOpts, sourceToken common.Address) (common.Address, error) {
	var out []interface{}
	err := _EVM2EVMOnRamp.contract.Call(opts, &out, "getPoolBySourceToken", sourceToken)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_EVM2EVMOnRamp *EVM2EVMOnRampSession) GetPoolBySourceToken(sourceToken common.Address) (common.Address, error) {
	return _EVM2EVMOnRamp.Contract.GetPoolBySourceToken(&_EVM2EVMOnRamp.CallOpts, sourceToken)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampCallerSession) GetPoolBySourceToken(sourceToken common.Address) (common.Address, error) {
	return _EVM2EVMOnRamp.Contract.GetPoolBySourceToken(&_EVM2EVMOnRamp.CallOpts, sourceToken)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampCaller) GetSenderNonce(opts *bind.CallOpts, sender common.Address) (uint64, error) {
	var out []interface{}
	err := _EVM2EVMOnRamp.contract.Call(opts, &out, "getSenderNonce", sender)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_EVM2EVMOnRamp *EVM2EVMOnRampSession) GetSenderNonce(sender common.Address) (uint64, error) {
	return _EVM2EVMOnRamp.Contract.GetSenderNonce(&_EVM2EVMOnRamp.CallOpts, sender)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampCallerSession) GetSenderNonce(sender common.Address) (uint64, error) {
	return _EVM2EVMOnRamp.Contract.GetSenderNonce(&_EVM2EVMOnRamp.CallOpts, sender)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampCaller) GetStaticConfig(opts *bind.CallOpts) (EVM2EVMOnRampStaticConfig, error) {
	var out []interface{}
	err := _EVM2EVMOnRamp.contract.Call(opts, &out, "getStaticConfig")

	if err != nil {
		return *new(EVM2EVMOnRampStaticConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(EVM2EVMOnRampStaticConfig)).(*EVM2EVMOnRampStaticConfig)

	return out0, err

}

func (_EVM2EVMOnRamp *EVM2EVMOnRampSession) GetStaticConfig() (EVM2EVMOnRampStaticConfig, error) {
	return _EVM2EVMOnRamp.Contract.GetStaticConfig(&_EVM2EVMOnRamp.CallOpts)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampCallerSession) GetStaticConfig() (EVM2EVMOnRampStaticConfig, error) {
	return _EVM2EVMOnRamp.Contract.GetStaticConfig(&_EVM2EVMOnRamp.CallOpts)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampCaller) GetSupportedTokens(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _EVM2EVMOnRamp.contract.Call(opts, &out, "getSupportedTokens")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_EVM2EVMOnRamp *EVM2EVMOnRampSession) GetSupportedTokens() ([]common.Address, error) {
	return _EVM2EVMOnRamp.Contract.GetSupportedTokens(&_EVM2EVMOnRamp.CallOpts)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampCallerSession) GetSupportedTokens() ([]common.Address, error) {
	return _EVM2EVMOnRamp.Contract.GetSupportedTokens(&_EVM2EVMOnRamp.CallOpts)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampCaller) GetTokenLimitAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EVM2EVMOnRamp.contract.Call(opts, &out, "getTokenLimitAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_EVM2EVMOnRamp *EVM2EVMOnRampSession) GetTokenLimitAdmin() (common.Address, error) {
	return _EVM2EVMOnRamp.Contract.GetTokenLimitAdmin(&_EVM2EVMOnRamp.CallOpts)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampCallerSession) GetTokenLimitAdmin() (common.Address, error) {
	return _EVM2EVMOnRamp.Contract.GetTokenLimitAdmin(&_EVM2EVMOnRamp.CallOpts)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampCaller) GetTokenTransferFeeConfig(opts *bind.CallOpts, token common.Address) (EVM2EVMOnRampTokenTransferFeeConfig, error) {
	var out []interface{}
	err := _EVM2EVMOnRamp.contract.Call(opts, &out, "getTokenTransferFeeConfig", token)

	if err != nil {
		return *new(EVM2EVMOnRampTokenTransferFeeConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(EVM2EVMOnRampTokenTransferFeeConfig)).(*EVM2EVMOnRampTokenTransferFeeConfig)

	return out0, err

}

func (_EVM2EVMOnRamp *EVM2EVMOnRampSession) GetTokenTransferFeeConfig(token common.Address) (EVM2EVMOnRampTokenTransferFeeConfig, error) {
	return _EVM2EVMOnRamp.Contract.GetTokenTransferFeeConfig(&_EVM2EVMOnRamp.CallOpts, token)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampCallerSession) GetTokenTransferFeeConfig(token common.Address) (EVM2EVMOnRampTokenTransferFeeConfig, error) {
	return _EVM2EVMOnRamp.Contract.GetTokenTransferFeeConfig(&_EVM2EVMOnRamp.CallOpts, token)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampCaller) LinkAvailableForPayment(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _EVM2EVMOnRamp.contract.Call(opts, &out, "linkAvailableForPayment")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_EVM2EVMOnRamp *EVM2EVMOnRampSession) LinkAvailableForPayment() (*big.Int, error) {
	return _EVM2EVMOnRamp.Contract.LinkAvailableForPayment(&_EVM2EVMOnRamp.CallOpts)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampCallerSession) LinkAvailableForPayment() (*big.Int, error) {
	return _EVM2EVMOnRamp.Contract.LinkAvailableForPayment(&_EVM2EVMOnRamp.CallOpts)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EVM2EVMOnRamp.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_EVM2EVMOnRamp *EVM2EVMOnRampSession) Owner() (common.Address, error) {
	return _EVM2EVMOnRamp.Contract.Owner(&_EVM2EVMOnRamp.CallOpts)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampCallerSession) Owner() (common.Address, error) {
	return _EVM2EVMOnRamp.Contract.Owner(&_EVM2EVMOnRamp.CallOpts)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _EVM2EVMOnRamp.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_EVM2EVMOnRamp *EVM2EVMOnRampSession) TypeAndVersion() (string, error) {
	return _EVM2EVMOnRamp.Contract.TypeAndVersion(&_EVM2EVMOnRamp.CallOpts)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampCallerSession) TypeAndVersion() (string, error) {
	return _EVM2EVMOnRamp.Contract.TypeAndVersion(&_EVM2EVMOnRamp.CallOpts)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EVM2EVMOnRamp.contract.Transact(opts, "acceptOwnership")
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampSession) AcceptOwnership() (*types.Transaction, error) {
	return _EVM2EVMOnRamp.Contract.AcceptOwnership(&_EVM2EVMOnRamp.TransactOpts)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _EVM2EVMOnRamp.Contract.AcceptOwnership(&_EVM2EVMOnRamp.TransactOpts)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampTransactor) ApplyPoolUpdates(opts *bind.TransactOpts, removes []InternalPoolUpdate, adds []InternalPoolUpdate) (*types.Transaction, error) {
	return _EVM2EVMOnRamp.contract.Transact(opts, "applyPoolUpdates", removes, adds)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampSession) ApplyPoolUpdates(removes []InternalPoolUpdate, adds []InternalPoolUpdate) (*types.Transaction, error) {
	return _EVM2EVMOnRamp.Contract.ApplyPoolUpdates(&_EVM2EVMOnRamp.TransactOpts, removes, adds)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampTransactorSession) ApplyPoolUpdates(removes []InternalPoolUpdate, adds []InternalPoolUpdate) (*types.Transaction, error) {
	return _EVM2EVMOnRamp.Contract.ApplyPoolUpdates(&_EVM2EVMOnRamp.TransactOpts, removes, adds)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampTransactor) ForwardFromRouter(opts *bind.TransactOpts, message ClientEVM2AnyMessage, feeTokenAmount *big.Int, originalSender common.Address) (*types.Transaction, error) {
	return _EVM2EVMOnRamp.contract.Transact(opts, "forwardFromRouter", message, feeTokenAmount, originalSender)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampSession) ForwardFromRouter(message ClientEVM2AnyMessage, feeTokenAmount *big.Int, originalSender common.Address) (*types.Transaction, error) {
	return _EVM2EVMOnRamp.Contract.ForwardFromRouter(&_EVM2EVMOnRamp.TransactOpts, message, feeTokenAmount, originalSender)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampTransactorSession) ForwardFromRouter(message ClientEVM2AnyMessage, feeTokenAmount *big.Int, originalSender common.Address) (*types.Transaction, error) {
	return _EVM2EVMOnRamp.Contract.ForwardFromRouter(&_EVM2EVMOnRamp.TransactOpts, message, feeTokenAmount, originalSender)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampTransactor) PayNops(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EVM2EVMOnRamp.contract.Transact(opts, "payNops")
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampSession) PayNops() (*types.Transaction, error) {
	return _EVM2EVMOnRamp.Contract.PayNops(&_EVM2EVMOnRamp.TransactOpts)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampTransactorSession) PayNops() (*types.Transaction, error) {
	return _EVM2EVMOnRamp.Contract.PayNops(&_EVM2EVMOnRamp.TransactOpts)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampTransactor) SetAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _EVM2EVMOnRamp.contract.Transact(opts, "setAdmin", newAdmin)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampSession) SetAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _EVM2EVMOnRamp.Contract.SetAdmin(&_EVM2EVMOnRamp.TransactOpts, newAdmin)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampTransactorSession) SetAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _EVM2EVMOnRamp.Contract.SetAdmin(&_EVM2EVMOnRamp.TransactOpts, newAdmin)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampTransactor) SetDynamicConfig(opts *bind.TransactOpts, dynamicConfig EVM2EVMOnRampDynamicConfig) (*types.Transaction, error) {
	return _EVM2EVMOnRamp.contract.Transact(opts, "setDynamicConfig", dynamicConfig)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampSession) SetDynamicConfig(dynamicConfig EVM2EVMOnRampDynamicConfig) (*types.Transaction, error) {
	return _EVM2EVMOnRamp.Contract.SetDynamicConfig(&_EVM2EVMOnRamp.TransactOpts, dynamicConfig)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampTransactorSession) SetDynamicConfig(dynamicConfig EVM2EVMOnRampDynamicConfig) (*types.Transaction, error) {
	return _EVM2EVMOnRamp.Contract.SetDynamicConfig(&_EVM2EVMOnRamp.TransactOpts, dynamicConfig)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampTransactor) SetFeeTokenConfig(opts *bind.TransactOpts, feeTokenConfigArgs []EVM2EVMOnRampFeeTokenConfigArgs) (*types.Transaction, error) {
	return _EVM2EVMOnRamp.contract.Transact(opts, "setFeeTokenConfig", feeTokenConfigArgs)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampSession) SetFeeTokenConfig(feeTokenConfigArgs []EVM2EVMOnRampFeeTokenConfigArgs) (*types.Transaction, error) {
	return _EVM2EVMOnRamp.Contract.SetFeeTokenConfig(&_EVM2EVMOnRamp.TransactOpts, feeTokenConfigArgs)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampTransactorSession) SetFeeTokenConfig(feeTokenConfigArgs []EVM2EVMOnRampFeeTokenConfigArgs) (*types.Transaction, error) {
	return _EVM2EVMOnRamp.Contract.SetFeeTokenConfig(&_EVM2EVMOnRamp.TransactOpts, feeTokenConfigArgs)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampTransactor) SetNops(opts *bind.TransactOpts, nopsAndWeights []EVM2EVMOnRampNopAndWeight) (*types.Transaction, error) {
	return _EVM2EVMOnRamp.contract.Transact(opts, "setNops", nopsAndWeights)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampSession) SetNops(nopsAndWeights []EVM2EVMOnRampNopAndWeight) (*types.Transaction, error) {
	return _EVM2EVMOnRamp.Contract.SetNops(&_EVM2EVMOnRamp.TransactOpts, nopsAndWeights)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampTransactorSession) SetNops(nopsAndWeights []EVM2EVMOnRampNopAndWeight) (*types.Transaction, error) {
	return _EVM2EVMOnRamp.Contract.SetNops(&_EVM2EVMOnRamp.TransactOpts, nopsAndWeights)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampTransactor) SetRateLimiterConfig(opts *bind.TransactOpts, config RateLimiterConfig) (*types.Transaction, error) {
	return _EVM2EVMOnRamp.contract.Transact(opts, "setRateLimiterConfig", config)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampSession) SetRateLimiterConfig(config RateLimiterConfig) (*types.Transaction, error) {
	return _EVM2EVMOnRamp.Contract.SetRateLimiterConfig(&_EVM2EVMOnRamp.TransactOpts, config)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampTransactorSession) SetRateLimiterConfig(config RateLimiterConfig) (*types.Transaction, error) {
	return _EVM2EVMOnRamp.Contract.SetRateLimiterConfig(&_EVM2EVMOnRamp.TransactOpts, config)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampTransactor) SetTokenTransferFeeConfig(opts *bind.TransactOpts, tokenTransferFeeConfigArgs []EVM2EVMOnRampTokenTransferFeeConfigArgs) (*types.Transaction, error) {
	return _EVM2EVMOnRamp.contract.Transact(opts, "setTokenTransferFeeConfig", tokenTransferFeeConfigArgs)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampSession) SetTokenTransferFeeConfig(tokenTransferFeeConfigArgs []EVM2EVMOnRampTokenTransferFeeConfigArgs) (*types.Transaction, error) {
	return _EVM2EVMOnRamp.Contract.SetTokenTransferFeeConfig(&_EVM2EVMOnRamp.TransactOpts, tokenTransferFeeConfigArgs)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampTransactorSession) SetTokenTransferFeeConfig(tokenTransferFeeConfigArgs []EVM2EVMOnRampTokenTransferFeeConfigArgs) (*types.Transaction, error) {
	return _EVM2EVMOnRamp.Contract.SetTokenTransferFeeConfig(&_EVM2EVMOnRamp.TransactOpts, tokenTransferFeeConfigArgs)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _EVM2EVMOnRamp.contract.Transact(opts, "transferOwnership", to)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _EVM2EVMOnRamp.Contract.TransferOwnership(&_EVM2EVMOnRamp.TransactOpts, to)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _EVM2EVMOnRamp.Contract.TransferOwnership(&_EVM2EVMOnRamp.TransactOpts, to)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampTransactor) WithdrawNonLinkFees(opts *bind.TransactOpts, feeToken common.Address, to common.Address) (*types.Transaction, error) {
	return _EVM2EVMOnRamp.contract.Transact(opts, "withdrawNonLinkFees", feeToken, to)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampSession) WithdrawNonLinkFees(feeToken common.Address, to common.Address) (*types.Transaction, error) {
	return _EVM2EVMOnRamp.Contract.WithdrawNonLinkFees(&_EVM2EVMOnRamp.TransactOpts, feeToken, to)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampTransactorSession) WithdrawNonLinkFees(feeToken common.Address, to common.Address) (*types.Transaction, error) {
	return _EVM2EVMOnRamp.Contract.WithdrawNonLinkFees(&_EVM2EVMOnRamp.TransactOpts, feeToken, to)
}

type EVM2EVMOnRampAdminSetIterator struct {
	Event *EVM2EVMOnRampAdminSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMOnRampAdminSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMOnRampAdminSet)
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
		it.Event = new(EVM2EVMOnRampAdminSet)
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

func (it *EVM2EVMOnRampAdminSetIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMOnRampAdminSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMOnRampAdminSet struct {
	NewAdmin common.Address
	Raw      types.Log
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampFilterer) FilterAdminSet(opts *bind.FilterOpts) (*EVM2EVMOnRampAdminSetIterator, error) {

	logs, sub, err := _EVM2EVMOnRamp.contract.FilterLogs(opts, "AdminSet")
	if err != nil {
		return nil, err
	}
	return &EVM2EVMOnRampAdminSetIterator{contract: _EVM2EVMOnRamp.contract, event: "AdminSet", logs: logs, sub: sub}, nil
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampFilterer) WatchAdminSet(opts *bind.WatchOpts, sink chan<- *EVM2EVMOnRampAdminSet) (event.Subscription, error) {

	logs, sub, err := _EVM2EVMOnRamp.contract.WatchLogs(opts, "AdminSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMOnRampAdminSet)
				if err := _EVM2EVMOnRamp.contract.UnpackLog(event, "AdminSet", log); err != nil {
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

func (_EVM2EVMOnRamp *EVM2EVMOnRampFilterer) ParseAdminSet(log types.Log) (*EVM2EVMOnRampAdminSet, error) {
	event := new(EVM2EVMOnRampAdminSet)
	if err := _EVM2EVMOnRamp.contract.UnpackLog(event, "AdminSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMOnRampCCIPSendRequestedIterator struct {
	Event *EVM2EVMOnRampCCIPSendRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMOnRampCCIPSendRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMOnRampCCIPSendRequested)
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
		it.Event = new(EVM2EVMOnRampCCIPSendRequested)
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

func (it *EVM2EVMOnRampCCIPSendRequestedIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMOnRampCCIPSendRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMOnRampCCIPSendRequested struct {
	Message InternalEVM2EVMMessage
	Raw     types.Log
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampFilterer) FilterCCIPSendRequested(opts *bind.FilterOpts) (*EVM2EVMOnRampCCIPSendRequestedIterator, error) {

	logs, sub, err := _EVM2EVMOnRamp.contract.FilterLogs(opts, "CCIPSendRequested")
	if err != nil {
		return nil, err
	}
	return &EVM2EVMOnRampCCIPSendRequestedIterator{contract: _EVM2EVMOnRamp.contract, event: "CCIPSendRequested", logs: logs, sub: sub}, nil
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampFilterer) WatchCCIPSendRequested(opts *bind.WatchOpts, sink chan<- *EVM2EVMOnRampCCIPSendRequested) (event.Subscription, error) {

	logs, sub, err := _EVM2EVMOnRamp.contract.WatchLogs(opts, "CCIPSendRequested")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMOnRampCCIPSendRequested)
				if err := _EVM2EVMOnRamp.contract.UnpackLog(event, "CCIPSendRequested", log); err != nil {
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

func (_EVM2EVMOnRamp *EVM2EVMOnRampFilterer) ParseCCIPSendRequested(log types.Log) (*EVM2EVMOnRampCCIPSendRequested, error) {
	event := new(EVM2EVMOnRampCCIPSendRequested)
	if err := _EVM2EVMOnRamp.contract.UnpackLog(event, "CCIPSendRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMOnRampConfigSetIterator struct {
	Event *EVM2EVMOnRampConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMOnRampConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMOnRampConfigSet)
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
		it.Event = new(EVM2EVMOnRampConfigSet)
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

func (it *EVM2EVMOnRampConfigSetIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMOnRampConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMOnRampConfigSet struct {
	StaticConfig  EVM2EVMOnRampStaticConfig
	DynamicConfig EVM2EVMOnRampDynamicConfig
	Raw           types.Log
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampFilterer) FilterConfigSet(opts *bind.FilterOpts) (*EVM2EVMOnRampConfigSetIterator, error) {

	logs, sub, err := _EVM2EVMOnRamp.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &EVM2EVMOnRampConfigSetIterator{contract: _EVM2EVMOnRamp.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *EVM2EVMOnRampConfigSet) (event.Subscription, error) {

	logs, sub, err := _EVM2EVMOnRamp.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMOnRampConfigSet)
				if err := _EVM2EVMOnRamp.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_EVM2EVMOnRamp *EVM2EVMOnRampFilterer) ParseConfigSet(log types.Log) (*EVM2EVMOnRampConfigSet, error) {
	event := new(EVM2EVMOnRampConfigSet)
	if err := _EVM2EVMOnRamp.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMOnRampFeeConfigSetIterator struct {
	Event *EVM2EVMOnRampFeeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMOnRampFeeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMOnRampFeeConfigSet)
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
		it.Event = new(EVM2EVMOnRampFeeConfigSet)
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

func (it *EVM2EVMOnRampFeeConfigSetIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMOnRampFeeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMOnRampFeeConfigSet struct {
	FeeConfig []EVM2EVMOnRampFeeTokenConfigArgs
	Raw       types.Log
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampFilterer) FilterFeeConfigSet(opts *bind.FilterOpts) (*EVM2EVMOnRampFeeConfigSetIterator, error) {

	logs, sub, err := _EVM2EVMOnRamp.contract.FilterLogs(opts, "FeeConfigSet")
	if err != nil {
		return nil, err
	}
	return &EVM2EVMOnRampFeeConfigSetIterator{contract: _EVM2EVMOnRamp.contract, event: "FeeConfigSet", logs: logs, sub: sub}, nil
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampFilterer) WatchFeeConfigSet(opts *bind.WatchOpts, sink chan<- *EVM2EVMOnRampFeeConfigSet) (event.Subscription, error) {

	logs, sub, err := _EVM2EVMOnRamp.contract.WatchLogs(opts, "FeeConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMOnRampFeeConfigSet)
				if err := _EVM2EVMOnRamp.contract.UnpackLog(event, "FeeConfigSet", log); err != nil {
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

func (_EVM2EVMOnRamp *EVM2EVMOnRampFilterer) ParseFeeConfigSet(log types.Log) (*EVM2EVMOnRampFeeConfigSet, error) {
	event := new(EVM2EVMOnRampFeeConfigSet)
	if err := _EVM2EVMOnRamp.contract.UnpackLog(event, "FeeConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMOnRampNopPaidIterator struct {
	Event *EVM2EVMOnRampNopPaid

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMOnRampNopPaidIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMOnRampNopPaid)
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
		it.Event = new(EVM2EVMOnRampNopPaid)
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

func (it *EVM2EVMOnRampNopPaidIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMOnRampNopPaidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMOnRampNopPaid struct {
	Nop    common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampFilterer) FilterNopPaid(opts *bind.FilterOpts, nop []common.Address) (*EVM2EVMOnRampNopPaidIterator, error) {

	var nopRule []interface{}
	for _, nopItem := range nop {
		nopRule = append(nopRule, nopItem)
	}

	logs, sub, err := _EVM2EVMOnRamp.contract.FilterLogs(opts, "NopPaid", nopRule)
	if err != nil {
		return nil, err
	}
	return &EVM2EVMOnRampNopPaidIterator{contract: _EVM2EVMOnRamp.contract, event: "NopPaid", logs: logs, sub: sub}, nil
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampFilterer) WatchNopPaid(opts *bind.WatchOpts, sink chan<- *EVM2EVMOnRampNopPaid, nop []common.Address) (event.Subscription, error) {

	var nopRule []interface{}
	for _, nopItem := range nop {
		nopRule = append(nopRule, nopItem)
	}

	logs, sub, err := _EVM2EVMOnRamp.contract.WatchLogs(opts, "NopPaid", nopRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMOnRampNopPaid)
				if err := _EVM2EVMOnRamp.contract.UnpackLog(event, "NopPaid", log); err != nil {
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

func (_EVM2EVMOnRamp *EVM2EVMOnRampFilterer) ParseNopPaid(log types.Log) (*EVM2EVMOnRampNopPaid, error) {
	event := new(EVM2EVMOnRampNopPaid)
	if err := _EVM2EVMOnRamp.contract.UnpackLog(event, "NopPaid", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMOnRampNopsSetIterator struct {
	Event *EVM2EVMOnRampNopsSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMOnRampNopsSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMOnRampNopsSet)
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
		it.Event = new(EVM2EVMOnRampNopsSet)
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

func (it *EVM2EVMOnRampNopsSetIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMOnRampNopsSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMOnRampNopsSet struct {
	NopWeightsTotal *big.Int
	NopsAndWeights  []EVM2EVMOnRampNopAndWeight
	Raw             types.Log
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampFilterer) FilterNopsSet(opts *bind.FilterOpts) (*EVM2EVMOnRampNopsSetIterator, error) {

	logs, sub, err := _EVM2EVMOnRamp.contract.FilterLogs(opts, "NopsSet")
	if err != nil {
		return nil, err
	}
	return &EVM2EVMOnRampNopsSetIterator{contract: _EVM2EVMOnRamp.contract, event: "NopsSet", logs: logs, sub: sub}, nil
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampFilterer) WatchNopsSet(opts *bind.WatchOpts, sink chan<- *EVM2EVMOnRampNopsSet) (event.Subscription, error) {

	logs, sub, err := _EVM2EVMOnRamp.contract.WatchLogs(opts, "NopsSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMOnRampNopsSet)
				if err := _EVM2EVMOnRamp.contract.UnpackLog(event, "NopsSet", log); err != nil {
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

func (_EVM2EVMOnRamp *EVM2EVMOnRampFilterer) ParseNopsSet(log types.Log) (*EVM2EVMOnRampNopsSet, error) {
	event := new(EVM2EVMOnRampNopsSet)
	if err := _EVM2EVMOnRamp.contract.UnpackLog(event, "NopsSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMOnRampOwnershipTransferRequestedIterator struct {
	Event *EVM2EVMOnRampOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMOnRampOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMOnRampOwnershipTransferRequested)
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
		it.Event = new(EVM2EVMOnRampOwnershipTransferRequested)
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

func (it *EVM2EVMOnRampOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMOnRampOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMOnRampOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*EVM2EVMOnRampOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _EVM2EVMOnRamp.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &EVM2EVMOnRampOwnershipTransferRequestedIterator{contract: _EVM2EVMOnRamp.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *EVM2EVMOnRampOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _EVM2EVMOnRamp.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMOnRampOwnershipTransferRequested)
				if err := _EVM2EVMOnRamp.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_EVM2EVMOnRamp *EVM2EVMOnRampFilterer) ParseOwnershipTransferRequested(log types.Log) (*EVM2EVMOnRampOwnershipTransferRequested, error) {
	event := new(EVM2EVMOnRampOwnershipTransferRequested)
	if err := _EVM2EVMOnRamp.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMOnRampOwnershipTransferredIterator struct {
	Event *EVM2EVMOnRampOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMOnRampOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMOnRampOwnershipTransferred)
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
		it.Event = new(EVM2EVMOnRampOwnershipTransferred)
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

func (it *EVM2EVMOnRampOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMOnRampOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMOnRampOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*EVM2EVMOnRampOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _EVM2EVMOnRamp.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &EVM2EVMOnRampOwnershipTransferredIterator{contract: _EVM2EVMOnRamp.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *EVM2EVMOnRampOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _EVM2EVMOnRamp.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMOnRampOwnershipTransferred)
				if err := _EVM2EVMOnRamp.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_EVM2EVMOnRamp *EVM2EVMOnRampFilterer) ParseOwnershipTransferred(log types.Log) (*EVM2EVMOnRampOwnershipTransferred, error) {
	event := new(EVM2EVMOnRampOwnershipTransferred)
	if err := _EVM2EVMOnRamp.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMOnRampPoolAddedIterator struct {
	Event *EVM2EVMOnRampPoolAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMOnRampPoolAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMOnRampPoolAdded)
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
		it.Event = new(EVM2EVMOnRampPoolAdded)
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

func (it *EVM2EVMOnRampPoolAddedIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMOnRampPoolAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMOnRampPoolAdded struct {
	Token common.Address
	Pool  common.Address
	Raw   types.Log
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampFilterer) FilterPoolAdded(opts *bind.FilterOpts) (*EVM2EVMOnRampPoolAddedIterator, error) {

	logs, sub, err := _EVM2EVMOnRamp.contract.FilterLogs(opts, "PoolAdded")
	if err != nil {
		return nil, err
	}
	return &EVM2EVMOnRampPoolAddedIterator{contract: _EVM2EVMOnRamp.contract, event: "PoolAdded", logs: logs, sub: sub}, nil
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampFilterer) WatchPoolAdded(opts *bind.WatchOpts, sink chan<- *EVM2EVMOnRampPoolAdded) (event.Subscription, error) {

	logs, sub, err := _EVM2EVMOnRamp.contract.WatchLogs(opts, "PoolAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMOnRampPoolAdded)
				if err := _EVM2EVMOnRamp.contract.UnpackLog(event, "PoolAdded", log); err != nil {
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

func (_EVM2EVMOnRamp *EVM2EVMOnRampFilterer) ParsePoolAdded(log types.Log) (*EVM2EVMOnRampPoolAdded, error) {
	event := new(EVM2EVMOnRampPoolAdded)
	if err := _EVM2EVMOnRamp.contract.UnpackLog(event, "PoolAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMOnRampPoolRemovedIterator struct {
	Event *EVM2EVMOnRampPoolRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMOnRampPoolRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMOnRampPoolRemoved)
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
		it.Event = new(EVM2EVMOnRampPoolRemoved)
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

func (it *EVM2EVMOnRampPoolRemovedIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMOnRampPoolRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMOnRampPoolRemoved struct {
	Token common.Address
	Pool  common.Address
	Raw   types.Log
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampFilterer) FilterPoolRemoved(opts *bind.FilterOpts) (*EVM2EVMOnRampPoolRemovedIterator, error) {

	logs, sub, err := _EVM2EVMOnRamp.contract.FilterLogs(opts, "PoolRemoved")
	if err != nil {
		return nil, err
	}
	return &EVM2EVMOnRampPoolRemovedIterator{contract: _EVM2EVMOnRamp.contract, event: "PoolRemoved", logs: logs, sub: sub}, nil
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampFilterer) WatchPoolRemoved(opts *bind.WatchOpts, sink chan<- *EVM2EVMOnRampPoolRemoved) (event.Subscription, error) {

	logs, sub, err := _EVM2EVMOnRamp.contract.WatchLogs(opts, "PoolRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMOnRampPoolRemoved)
				if err := _EVM2EVMOnRamp.contract.UnpackLog(event, "PoolRemoved", log); err != nil {
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

func (_EVM2EVMOnRamp *EVM2EVMOnRampFilterer) ParsePoolRemoved(log types.Log) (*EVM2EVMOnRampPoolRemoved, error) {
	event := new(EVM2EVMOnRampPoolRemoved)
	if err := _EVM2EVMOnRamp.contract.UnpackLog(event, "PoolRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMOnRampTokenTransferFeeConfigSetIterator struct {
	Event *EVM2EVMOnRampTokenTransferFeeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMOnRampTokenTransferFeeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMOnRampTokenTransferFeeConfigSet)
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
		it.Event = new(EVM2EVMOnRampTokenTransferFeeConfigSet)
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

func (it *EVM2EVMOnRampTokenTransferFeeConfigSetIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMOnRampTokenTransferFeeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMOnRampTokenTransferFeeConfigSet struct {
	TransferFeeConfig []EVM2EVMOnRampTokenTransferFeeConfigArgs
	Raw               types.Log
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampFilterer) FilterTokenTransferFeeConfigSet(opts *bind.FilterOpts) (*EVM2EVMOnRampTokenTransferFeeConfigSetIterator, error) {

	logs, sub, err := _EVM2EVMOnRamp.contract.FilterLogs(opts, "TokenTransferFeeConfigSet")
	if err != nil {
		return nil, err
	}
	return &EVM2EVMOnRampTokenTransferFeeConfigSetIterator{contract: _EVM2EVMOnRamp.contract, event: "TokenTransferFeeConfigSet", logs: logs, sub: sub}, nil
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampFilterer) WatchTokenTransferFeeConfigSet(opts *bind.WatchOpts, sink chan<- *EVM2EVMOnRampTokenTransferFeeConfigSet) (event.Subscription, error) {

	logs, sub, err := _EVM2EVMOnRamp.contract.WatchLogs(opts, "TokenTransferFeeConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMOnRampTokenTransferFeeConfigSet)
				if err := _EVM2EVMOnRamp.contract.UnpackLog(event, "TokenTransferFeeConfigSet", log); err != nil {
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

func (_EVM2EVMOnRamp *EVM2EVMOnRampFilterer) ParseTokenTransferFeeConfigSet(log types.Log) (*EVM2EVMOnRampTokenTransferFeeConfigSet, error) {
	event := new(EVM2EVMOnRampTokenTransferFeeConfigSet)
	if err := _EVM2EVMOnRamp.contract.UnpackLog(event, "TokenTransferFeeConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetNops struct {
	NopsAndWeights []EVM2EVMOnRampNopAndWeight
	WeightsTotal   *big.Int
}

func (_EVM2EVMOnRamp *EVM2EVMOnRamp) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _EVM2EVMOnRamp.abi.Events["AdminSet"].ID:
		return _EVM2EVMOnRamp.ParseAdminSet(log)
	case _EVM2EVMOnRamp.abi.Events["CCIPSendRequested"].ID:
		return _EVM2EVMOnRamp.ParseCCIPSendRequested(log)
	case _EVM2EVMOnRamp.abi.Events["ConfigSet"].ID:
		return _EVM2EVMOnRamp.ParseConfigSet(log)
	case _EVM2EVMOnRamp.abi.Events["FeeConfigSet"].ID:
		return _EVM2EVMOnRamp.ParseFeeConfigSet(log)
	case _EVM2EVMOnRamp.abi.Events["NopPaid"].ID:
		return _EVM2EVMOnRamp.ParseNopPaid(log)
	case _EVM2EVMOnRamp.abi.Events["NopsSet"].ID:
		return _EVM2EVMOnRamp.ParseNopsSet(log)
	case _EVM2EVMOnRamp.abi.Events["OwnershipTransferRequested"].ID:
		return _EVM2EVMOnRamp.ParseOwnershipTransferRequested(log)
	case _EVM2EVMOnRamp.abi.Events["OwnershipTransferred"].ID:
		return _EVM2EVMOnRamp.ParseOwnershipTransferred(log)
	case _EVM2EVMOnRamp.abi.Events["PoolAdded"].ID:
		return _EVM2EVMOnRamp.ParsePoolAdded(log)
	case _EVM2EVMOnRamp.abi.Events["PoolRemoved"].ID:
		return _EVM2EVMOnRamp.ParsePoolRemoved(log)
	case _EVM2EVMOnRamp.abi.Events["TokenTransferFeeConfigSet"].ID:
		return _EVM2EVMOnRamp.ParseTokenTransferFeeConfigSet(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (EVM2EVMOnRampAdminSet) Topic() common.Hash {
	return common.HexToHash("0x8fe72c3e0020beb3234e76ae6676fa576fbfcae600af1c4fea44784cf0db329c")
}

func (EVM2EVMOnRampCCIPSendRequested) Topic() common.Hash {
	return common.HexToHash("0xaffc45517195d6499808c643bd4a7b0ffeedf95bea5852840d7bfcf63f59e821")
}

func (EVM2EVMOnRampConfigSet) Topic() common.Hash {
	return common.HexToHash("0x72c6aaba4dde02f77d291123a76185c418ba63f8c217a2d56b08aec84e9bbfb8")
}

func (EVM2EVMOnRampFeeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2386f61ab5cafc3fed44f9f614f721ab53479ef64067fd16c1a2491b63ddf1a8")
}

func (EVM2EVMOnRampNopPaid) Topic() common.Hash {
	return common.HexToHash("0x55fdec2aab60a41fa5abb106670eb1006f5aeaee1ba7afea2bc89b5b3ec7678f")
}

func (EVM2EVMOnRampNopsSet) Topic() common.Hash {
	return common.HexToHash("0x8c337bff38141c507abd25c547606bdde78fe8c12e941ab613f3a565fea6cd24")
}

func (EVM2EVMOnRampOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (EVM2EVMOnRampOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (EVM2EVMOnRampPoolAdded) Topic() common.Hash {
	return common.HexToHash("0x95f865c2808f8b2a85eea2611db7843150ee7835ef1403f9755918a97d76933c")
}

func (EVM2EVMOnRampPoolRemoved) Topic() common.Hash {
	return common.HexToHash("0x987eb3c2f78454541205f72f34839b434c306c9eaf4922efd7c0c3060fdb2e4c")
}

func (EVM2EVMOnRampTokenTransferFeeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x4230c60a9725eb5fb992cf6a215398b4e81b4606d4a1e6be8dfe0b60dc172ec1")
}

func (_EVM2EVMOnRamp *EVM2EVMOnRamp) Address() common.Address {
	return _EVM2EVMOnRamp.address
}

type EVM2EVMOnRampInterface interface {
	CurrentRateLimiterState(opts *bind.CallOpts) (RateLimiterTokenBucket, error)

	GetDynamicConfig(opts *bind.CallOpts) (EVM2EVMOnRampDynamicConfig, error)

	GetExpectedNextSequenceNumber(opts *bind.CallOpts) (uint64, error)

	GetFee(opts *bind.CallOpts, message ClientEVM2AnyMessage) (*big.Int, error)

	GetFeeTokenConfig(opts *bind.CallOpts, token common.Address) (EVM2EVMOnRampFeeTokenConfig, error)

	GetNopFeesJuels(opts *bind.CallOpts) (*big.Int, error)

	GetNops(opts *bind.CallOpts) (GetNops,

		error)

	GetPoolBySourceToken(opts *bind.CallOpts, sourceToken common.Address) (common.Address, error)

	GetSenderNonce(opts *bind.CallOpts, sender common.Address) (uint64, error)

	GetStaticConfig(opts *bind.CallOpts) (EVM2EVMOnRampStaticConfig, error)

	GetSupportedTokens(opts *bind.CallOpts) ([]common.Address, error)

	GetTokenLimitAdmin(opts *bind.CallOpts) (common.Address, error)

	GetTokenTransferFeeConfig(opts *bind.CallOpts, token common.Address) (EVM2EVMOnRampTokenTransferFeeConfig, error)

	LinkAvailableForPayment(opts *bind.CallOpts) (*big.Int, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	ApplyPoolUpdates(opts *bind.TransactOpts, removes []InternalPoolUpdate, adds []InternalPoolUpdate) (*types.Transaction, error)

	ForwardFromRouter(opts *bind.TransactOpts, message ClientEVM2AnyMessage, feeTokenAmount *big.Int, originalSender common.Address) (*types.Transaction, error)

	PayNops(opts *bind.TransactOpts) (*types.Transaction, error)

	SetAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error)

	SetDynamicConfig(opts *bind.TransactOpts, dynamicConfig EVM2EVMOnRampDynamicConfig) (*types.Transaction, error)

	SetFeeTokenConfig(opts *bind.TransactOpts, feeTokenConfigArgs []EVM2EVMOnRampFeeTokenConfigArgs) (*types.Transaction, error)

	SetNops(opts *bind.TransactOpts, nopsAndWeights []EVM2EVMOnRampNopAndWeight) (*types.Transaction, error)

	SetRateLimiterConfig(opts *bind.TransactOpts, config RateLimiterConfig) (*types.Transaction, error)

	SetTokenTransferFeeConfig(opts *bind.TransactOpts, tokenTransferFeeConfigArgs []EVM2EVMOnRampTokenTransferFeeConfigArgs) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	WithdrawNonLinkFees(opts *bind.TransactOpts, feeToken common.Address, to common.Address) (*types.Transaction, error)

	FilterAdminSet(opts *bind.FilterOpts) (*EVM2EVMOnRampAdminSetIterator, error)

	WatchAdminSet(opts *bind.WatchOpts, sink chan<- *EVM2EVMOnRampAdminSet) (event.Subscription, error)

	ParseAdminSet(log types.Log) (*EVM2EVMOnRampAdminSet, error)

	FilterCCIPSendRequested(opts *bind.FilterOpts) (*EVM2EVMOnRampCCIPSendRequestedIterator, error)

	WatchCCIPSendRequested(opts *bind.WatchOpts, sink chan<- *EVM2EVMOnRampCCIPSendRequested) (event.Subscription, error)

	ParseCCIPSendRequested(log types.Log) (*EVM2EVMOnRampCCIPSendRequested, error)

	FilterConfigSet(opts *bind.FilterOpts) (*EVM2EVMOnRampConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *EVM2EVMOnRampConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*EVM2EVMOnRampConfigSet, error)

	FilterFeeConfigSet(opts *bind.FilterOpts) (*EVM2EVMOnRampFeeConfigSetIterator, error)

	WatchFeeConfigSet(opts *bind.WatchOpts, sink chan<- *EVM2EVMOnRampFeeConfigSet) (event.Subscription, error)

	ParseFeeConfigSet(log types.Log) (*EVM2EVMOnRampFeeConfigSet, error)

	FilterNopPaid(opts *bind.FilterOpts, nop []common.Address) (*EVM2EVMOnRampNopPaidIterator, error)

	WatchNopPaid(opts *bind.WatchOpts, sink chan<- *EVM2EVMOnRampNopPaid, nop []common.Address) (event.Subscription, error)

	ParseNopPaid(log types.Log) (*EVM2EVMOnRampNopPaid, error)

	FilterNopsSet(opts *bind.FilterOpts) (*EVM2EVMOnRampNopsSetIterator, error)

	WatchNopsSet(opts *bind.WatchOpts, sink chan<- *EVM2EVMOnRampNopsSet) (event.Subscription, error)

	ParseNopsSet(log types.Log) (*EVM2EVMOnRampNopsSet, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*EVM2EVMOnRampOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *EVM2EVMOnRampOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*EVM2EVMOnRampOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*EVM2EVMOnRampOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *EVM2EVMOnRampOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*EVM2EVMOnRampOwnershipTransferred, error)

	FilterPoolAdded(opts *bind.FilterOpts) (*EVM2EVMOnRampPoolAddedIterator, error)

	WatchPoolAdded(opts *bind.WatchOpts, sink chan<- *EVM2EVMOnRampPoolAdded) (event.Subscription, error)

	ParsePoolAdded(log types.Log) (*EVM2EVMOnRampPoolAdded, error)

	FilterPoolRemoved(opts *bind.FilterOpts) (*EVM2EVMOnRampPoolRemovedIterator, error)

	WatchPoolRemoved(opts *bind.WatchOpts, sink chan<- *EVM2EVMOnRampPoolRemoved) (event.Subscription, error)

	ParsePoolRemoved(log types.Log) (*EVM2EVMOnRampPoolRemoved, error)

	FilterTokenTransferFeeConfigSet(opts *bind.FilterOpts) (*EVM2EVMOnRampTokenTransferFeeConfigSetIterator, error)

	WatchTokenTransferFeeConfigSet(opts *bind.WatchOpts, sink chan<- *EVM2EVMOnRampTokenTransferFeeConfigSet) (event.Subscription, error)

	ParseTokenTransferFeeConfigSet(log types.Log) (*EVM2EVMOnRampTokenTransferFeeConfigSet, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
