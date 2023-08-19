// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package evm_2_evm_onramp_1_0_0

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
	Router          common.Address
	MaxTokensLength uint16
	PriceRegistry   common.Address
	MaxDataSize     uint32
	MaxGasLimit     uint64
}

type EVM2EVMOnRampFeeTokenConfig struct {
	NetworkFeeAmountUSD   *big.Int
	GasMultiplier         uint64
	DestGasOverhead       uint32
	DestGasPerPayloadByte uint16
	Enabled               bool
}

type EVM2EVMOnRampFeeTokenConfigArgs struct {
	Token                 common.Address
	GasMultiplier         uint64
	NetworkFeeAmountUSD   *big.Int
	DestGasOverhead       uint32
	DestGasPerPayloadByte uint16
	Enabled               bool
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
	MinFee uint32
	MaxFee uint32
	Ratio  uint16
}

type EVM2EVMOnRampTokenTransferFeeConfigArgs struct {
	Token  common.Address
	MinFee uint32
	MaxFee uint32
	Ratio  uint16
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
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"linkToken\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"defaultTxGasLimit\",\"type\":\"uint64\"},{\"internalType\":\"uint96\",\"name\":\"maxNopFeesJuels\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"prevOnRamp\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"armProxy\",\"type\":\"address\"}],\"internalType\":\"structEVM2EVMOnRamp.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"maxTokensLength\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"priceRegistry\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"maxDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"maxGasLimit\",\"type\":\"uint64\"}],\"internalType\":\"structEVM2EVMOnRamp.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"}],\"internalType\":\"structInternal.PoolUpdate[]\",\"name\":\"tokensAndPools\",\"type\":\"tuple[]\"},{\"internalType\":\"address[]\",\"name\":\"allowlist\",\"type\":\"address[]\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"rateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"gasMultiplier\",\"type\":\"uint64\"},{\"internalType\":\"uint96\",\"name\":\"networkFeeAmountUSD\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerPayloadByte\",\"type\":\"uint16\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"internalType\":\"structEVM2EVMOnRamp.FeeTokenConfigArgs[]\",\"name\":\"feeTokenConfigs\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"minFee\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxFee\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"ratio\",\"type\":\"uint16\"}],\"internalType\":\"structEVM2EVMOnRamp.TokenTransferFeeConfigArgs[]\",\"name\":\"tokenTransferFeeConfigArgs\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"nop\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"weight\",\"type\":\"uint16\"}],\"internalType\":\"structEVM2EVMOnRamp.NopAndWeight[]\",\"name\":\"nopsAndWeights\",\"type\":\"tuple[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"}],\"name\":\"AggregateValueMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"}],\"name\":\"AggregateValueRateLimitReached\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BadARMSignal\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BucketOverfilled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedAddress\",\"type\":\"bytes\"}],\"name\":\"InvalidAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidConfig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidExtraArgsTag\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"nop\",\"type\":\"address\"}],\"name\":\"InvalidNopAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTokenPoolConfig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidWithdrawParams\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LinkBalanceNotSettled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxFeeBalanceReached\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MessageGasLimitTooHigh\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"maxSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actualSize\",\"type\":\"uint256\"}],\"name\":\"MessageTooLarge\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeCalledByRouter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoFeesToPay\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoNopsToPay\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"NotAFeeToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByAdminOrOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAdminOrNop\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PoolAlreadyAdded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"PoolDoesNotExist\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"PriceNotFoundForToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RouterMustSetOriginalSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"SenderNotAllowed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TokenPoolMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenRateLimitReached\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyNops\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnsupportedNumberOfTokens\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"UnsupportedToken\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AllowListAdd\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"name\":\"AllowListEnabledSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AllowListRemove\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"feeTokenAmount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"strict\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMTokenAmount[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"}],\"indexed\":false,\"internalType\":\"structInternal.EVM2EVMMessage\",\"name\":\"message\",\"type\":\"tuple\"}],\"name\":\"CCIPSendRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"linkToken\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"defaultTxGasLimit\",\"type\":\"uint64\"},{\"internalType\":\"uint96\",\"name\":\"maxNopFeesJuels\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"prevOnRamp\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"armProxy\",\"type\":\"address\"}],\"indexed\":false,\"internalType\":\"structEVM2EVMOnRamp.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"maxTokensLength\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"priceRegistry\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"maxDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"maxGasLimit\",\"type\":\"uint64\"}],\"indexed\":false,\"internalType\":\"structEVM2EVMOnRamp.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"gasMultiplier\",\"type\":\"uint64\"},{\"internalType\":\"uint96\",\"name\":\"networkFeeAmountUSD\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerPayloadByte\",\"type\":\"uint16\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"indexed\":false,\"internalType\":\"structEVM2EVMOnRamp.FeeTokenConfigArgs[]\",\"name\":\"feeConfig\",\"type\":\"tuple[]\"}],\"name\":\"FeeConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"nop\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"NopPaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nopWeightsTotal\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"nop\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"weight\",\"type\":\"uint16\"}],\"indexed\":false,\"internalType\":\"structEVM2EVMOnRamp.NopAndWeight[]\",\"name\":\"nopsAndWeights\",\"type\":\"tuple[]\"}],\"name\":\"NopsSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"}],\"name\":\"PoolAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"}],\"name\":\"PoolRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"minFee\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxFee\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"ratio\",\"type\":\"uint16\"}],\"indexed\":false,\"internalType\":\"structEVM2EVMOnRamp.TokenTransferFeeConfigArgs[]\",\"name\":\"transferFeeConfig\",\"type\":\"tuple[]\"}],\"name\":\"TokenTransferFeeConfigSet\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"removes\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"adds\",\"type\":\"address[]\"}],\"name\":\"applyAllowListUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"}],\"internalType\":\"structInternal.PoolUpdate[]\",\"name\":\"removes\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"}],\"internalType\":\"structInternal.PoolUpdate[]\",\"name\":\"adds\",\"type\":\"tuple[]\"}],\"name\":\"applyPoolUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMTokenAmount[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structClient.EVM2AnyMessage\",\"name\":\"message\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"feeTokenAmount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"originalSender\",\"type\":\"address\"}],\"name\":\"forwardFromRouter\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowList\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowListEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDynamicConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"maxTokensLength\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"priceRegistry\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"maxDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"maxGasLimit\",\"type\":\"uint64\"}],\"internalType\":\"structEVM2EVMOnRamp.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getExpectedNextSequenceNumber\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMTokenAmount[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structClient.EVM2AnyMessage\",\"name\":\"message\",\"type\":\"tuple\"}],\"name\":\"getFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getFeeTokenConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint96\",\"name\":\"networkFeeAmountUSD\",\"type\":\"uint96\"},{\"internalType\":\"uint64\",\"name\":\"gasMultiplier\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerPayloadByte\",\"type\":\"uint16\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"internalType\":\"structEVM2EVMOnRamp.FeeTokenConfig\",\"name\":\"feeTokenConfig\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNopFeesJuels\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNops\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"nop\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"weight\",\"type\":\"uint16\"}],\"internalType\":\"structEVM2EVMOnRamp.NopAndWeight[]\",\"name\":\"nopsAndWeights\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"weightsTotal\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"sourceToken\",\"type\":\"address\"}],\"name\":\"getPoolBySourceToken\",\"outputs\":[{\"internalType\":\"contractIPool\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"getSenderNonce\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStaticConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"linkToken\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"defaultTxGasLimit\",\"type\":\"uint64\"},{\"internalType\":\"uint96\",\"name\":\"maxNopFeesJuels\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"prevOnRamp\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"armProxy\",\"type\":\"address\"}],\"internalType\":\"structEVM2EVMOnRamp.StaticConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSupportedTokens\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTokenLimitAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getTokenTransferFeeConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"minFee\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxFee\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"ratio\",\"type\":\"uint16\"}],\"internalType\":\"structEVM2EVMOnRamp.TokenTransferFeeConfig\",\"name\":\"tokenTransferFeeConfig\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkAvailableForPayment\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"payNops\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"setAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"name\":\"setAllowListEnabled\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"maxTokensLength\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"priceRegistry\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"maxDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"maxGasLimit\",\"type\":\"uint64\"}],\"internalType\":\"structEVM2EVMOnRamp.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"name\":\"setDynamicConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"gasMultiplier\",\"type\":\"uint64\"},{\"internalType\":\"uint96\",\"name\":\"networkFeeAmountUSD\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerPayloadByte\",\"type\":\"uint16\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"internalType\":\"structEVM2EVMOnRamp.FeeTokenConfigArgs[]\",\"name\":\"feeTokenConfigArgs\",\"type\":\"tuple[]\"}],\"name\":\"setFeeTokenConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"nop\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"weight\",\"type\":\"uint16\"}],\"internalType\":\"structEVM2EVMOnRamp.NopAndWeight[]\",\"name\":\"nopsAndWeights\",\"type\":\"tuple[]\"}],\"name\":\"setNops\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"setRateLimiterConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"minFee\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxFee\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"ratio\",\"type\":\"uint16\"}],\"internalType\":\"structEVM2EVMOnRamp.TokenTransferFeeConfigArgs[]\",\"name\":\"tokenTransferFeeConfigArgs\",\"type\":\"tuple[]\"}],\"name\":\"setTokenTransferFeeConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdrawNonLinkFees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6101806040526012805460ff60c01b191690553480156200001f57600080fd5b50604051620081c9380380620081c9833981016040819052620000429162001eb3565b8333806000816200009a5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000cd57620000cd8162000369565b50506040805160a081018252602084810180516001600160801b039081168085524263ffffffff169385018490528751151585870181905292518216606086018190529790950151166080909301839052600380546001600160a01b031916909417600160801b9283021760ff60a01b1916600160a01b90910217909255029091176004555087516001600160a01b0316158062000176575060208801516001600160401b0316155b806200018d575060408801516001600160401b0316155b80620001a4575060608801516001600160401b0316155b80620001bb575060c08801516001600160a01b0316155b15620001da576040516306b7c75960e31b815260040160405180910390fd5b6020808901516040808b015181517fbdd59ac4dd1d82276c9a9c5d2656546346b9dcdb1f9b4204aed4ec15c23d7d3a948101949094526001600160401b039283169184019190915216606082015230608082015260a00160408051601f19818403018152918152815160209283012060809081528a516001600160a01b0390811660e052928b01516001600160401b0390811661010052918b015182166101205260608b015190911660a0908152908a01516001600160601b031660c0908152908a01518216610140528901511661016052620002b78762000414565b620002c28362000566565b620002cd82620006d8565b620002d881620007e5565b6040805160008082526020820190925262000323916200031b565b6040805180820190915260008082526020820152815260200190600190039081620002f35790505b508762000a0f565b8451156200035b576012805460ff60c81b1916600160c81b179055604080516000808252602082019092526200035b91508662000d12565b505050505050505062002405565b336001600160a01b03821603620003c35760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000091565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60408101516001600160a01b031662000440576040516306b7c75960e31b815260040160405180910390fd5b8051600580546020808501516001600160a01b039485166001600160b01b031990931692909217600160a01b61ffff909316830217909255604080850151600680546060808901516080808b0151958a166001600160c01b03199094169390931763ffffffff909116909602959095176001600160c01b0316600160c01b6001600160401b039485160217909155825160e080820185525187168152610100518316958101959095526101205182168584015260a080519092169385019390935260c080516001600160601b03169385019390935261014051851690840152610160519093169082015290517fdd226617d8d287f40a64c54741bbcdc492b3e096ef16bc5273a18cb6ab85f124916200055b91849062001fd6565b60405180910390a150565b60005b8151811015620006a65760008282815181106200058a576200058a62002086565b6020908102919091018101516040805160a08082018352828401516001600160601b039081168352848601516001600160401b0390811684880190815260608088015163ffffffff9081168789019081526080808b015161ffff908116948a01948552978b0151151590890190815299516001600160a01b03166000908152600f909b52979099209551865492519751915198511515600160d01b0260ff60d01b1999909616600160c01b0261ffff60c01b1992909a16600160a01b029190911665ffffffffffff60a01b19979093166c01000000000000000000000000026001600160a01b03199092169316929092179190911793909316929092179390931791909116179055506200069e81620020b2565b905062000569565b507ffba339fca97870ffdfaedbae3745db5e6de1a6909dfd0e0dbb56917469ffe236816040516200055b9190620020ce565b60005b8151811015620007b3576000828281518110620006fc57620006fc62002086565b60209081029190910181015160408051606080820183528385015163ffffffff90811683528385015181168387019081529185015161ffff90811684860190815295516001600160a01b031660009081526010909752939095209151825491519451909316680100000000000000000261ffff60401b19948616640100000000026001600160401b03199092169390951692909217919091179190911691909117905550620007ab81620020b2565b9050620006db565b507fcb0c5f472d325cf0c56953fc81870ddd80d0d3c9a3fbfe777002d75f380dfb81816040516200055b91906200216f565b805160408111156200080a57604051635ad0867d60e11b815260040160405180910390fd5b6012546c01000000000000000000000000900463ffffffff161580159062000854575060125463ffffffff6c010000000000000000000000008204166001600160601b0390911610155b1562000864576200086462000e5d565b60006200087260076200105d565b90505b8015620008be576000620008986200088f600184620021de565b60079062001070565b509050620008a86007826200108e565b505080620008b690620021f4565b905062000875565b506000805b82811015620009a6576000848281518110620008e357620008e362002086565b6020026020010151600001519050600085838151811062000908576200090862002086565b602002602001015160200151905060e0516001600160a01b0316826001600160a01b031614806200094057506001600160a01b038216155b156200096b57604051634de938d160e01b81526001600160a01b038316600482015260240162000091565b6200097d60078361ffff8416620010ac565b506200098e61ffff8216856200220e565b93505050806200099e90620020b2565b9050620008c3565b506012805463ffffffff60601b19166c0100000000000000000000000063ffffffff8416021790556040517f8c337bff38141c507abd25c547606bdde78fe8c12e941ab613f3a565fea6cd249062000a0290839086906200222e565b60405180910390a1505050565b60005b825181101562000b4b57600083828151811062000a335762000a3362002086565b6020026020010151600001519050600084838151811062000a585762000a5862002086565b6020908102919091018101510151905062000a75600a83620010cc565b62000a9f576040516373913ebd60e01b81526001600160a01b038316600482015260240162000091565b6001600160a01b03811662000ab6600a84620010e3565b6001600160a01b03161462000ade57604051630d98f73360e31b815260040160405180910390fd5b62000aeb600a83620010fa565b1562000b3557604080516001600160a01b038085168252831660208201527f987eb3c2f78454541205f72f34839b434c306c9eaf4922efd7c0c3060fdb2e4c910160405180910390a15b50508062000b4390620020b2565b905062000a12565b5060005b815181101562000d0d57600082828151811062000b705762000b7062002086565b6020026020010151600001519050600083838151811062000b955762000b9562002086565b602002602001015160200151905060006001600160a01b0316826001600160a01b0316148062000bcc57506001600160a01b038116155b1562000bea5760405162d8548360e71b815260040160405180910390fd5b806001600160a01b03166321df0da76040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000c29573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000c4f91906200229b565b6001600160a01b0316826001600160a01b03161462000c8157604051630d98f73360e31b815260040160405180910390fd5b62000c8f600a838362001111565b1562000cde57604080516001600160a01b038085168252831660208201527f95f865c2808f8b2a85eea2611db7843150ee7835ef1403f9755918a97d76933c910160405180910390a162000cf7565b604051633caf458560e01b815260040160405180910390fd5b50508062000d0590620020b2565b905062000b4f565b505050565b60005b825181101562000da757600083828151811062000d365762000d3662002086565b6020908102919091010151905062000d50600d8262001129565b1562000d93576040516001600160a01b03821681527f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf75669060200160405180910390a15b5062000d9f81620020b2565b905062000d15565b5060005b815181101562000d0d57600082828151811062000dcc5762000dcc62002086565b6020026020010151905060006001600160a01b0316816001600160a01b03160362000df8575062000e4a565b62000e05600d8262001140565b1562000e48576040516001600160a01b03821681527f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d89060200160405180910390a15b505b62000e5581620020b2565b905062000dab565b6000546001600160a01b0316331480159062000e8457506002546001600160a01b03163314155b801562000e9b575062000e9960073362001157565b155b1562000eba5760405163032bb72b60e31b815260040160405180910390fd5b6012546c01000000000000000000000000900463ffffffff16600081900362000ef65760405163990e30bf60e01b815260040160405180910390fd5b6012546001600160601b03168181101562000f24576040516311a1ee3b60e31b815260040160405180910390fd5b600062000f306200116e565b121562000f5057604051631e9acf1760e31b815260040160405180910390fd5b80600062000f5f60076200105d565b905060005b81811015620010375760008062000f7d60078462001070565b909250905060008762000f9a836001600160601b038a16620022bb565b62000fa69190620022d5565b905062000fb48187620022f8565b60e05190965062000fd9906001600160a01b0316846001600160601b038416620011fc565b6040516001600160601b03821681526001600160a01b038416907f55fdec2aab60a41fa5abb106670eb1006f5aeaee1ba7afea2bc89b5b3ec7678f9060200160405180910390a2505050806200102f90620020b2565b905062000f64565b5050601280546001600160601b0319166001600160601b03929092169190911790555050565b60006200106a8262001254565b92915050565b600080808062001081868662001261565b9097909650945050505050565b6000620010a5836001600160a01b0384166200128e565b9392505050565b6000620010c4846001600160a01b03851684620012ad565b949350505050565b6000620010a5836001600160a01b038416620012cc565b6000620010a5836001600160a01b038416620012da565b6000620010a5836001600160a01b038416620012e8565b6000620010c4846001600160a01b03851684620012f6565b6000620010a5836001600160a01b0384166200130e565b6000620010a5836001600160a01b03841662001419565b6000620010a5836001600160a01b0384166200146b565b60125460e0516040516370a0823160e01b81523060048201526000926001600160601b0316916001600160a01b0316906370a0823190602401602060405180830381865afa158015620011c5573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620011eb91906200231b565b620011f7919062002335565b905090565b604080516001600160a01b038416602482015260448082018490528251808303909101815260649091019091526020810180516001600160e01b0390811663a9059cbb60e01b1790915262000d0d9185916200147916565b60006200106a826200154a565b6000808062001271858562001555565b600081815260029690960160205260409095205494959350505050565b60008181526002830160205260408120819055620010a5838362001563565b60008281526002840160205260408120829055620010c4848462001571565b6000620010a583836200146b565b6000620010a583836200157f565b6000620010a583836200128e565b6000620010c484846001600160a01b038516620012ad565b600081815260018301602052604081205480156200140757600062001335600183620021de565b85549091506000906200134b90600190620021de565b9050818114620013b75760008660000182815481106200136f576200136f62002086565b906000526020600020015490508087600001848154811062001395576200139562002086565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080620013cb57620013cb62002358565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506200106a565b60009150506200106a565b5092915050565b600081815260018301602052604081205462001462575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556200106a565b5060006200106a565b6000620010a58383620015f4565b6040805180820190915260208082527f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c656490820152600090620014c8906001600160a01b0385169084906200160d565b80519091501562000d0d5780806020019051810190620014e991906200236e565b62000d0d5760405162461bcd60e51b815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e6044820152691bdd081cdd58d8d9595960b21b606482015260840162000091565b60006200106a825490565b6000620010a583836200161e565b6000620010a583836200130e565b6000620010a5838362001419565b600081815260028301602052604081205480151580620015a65750620015a684846200146b565b620010a55760405162461bcd60e51b815260206004820152601e60248201527f456e756d657261626c654d61703a206e6f6e6578697374656e74206b65790000604482015260640162000091565b60008181526001830160205260408120541515620010a5565b6060620010c484846000856200164b565b600082600001828154811062001638576200163862002086565b9060005260206000200154905092915050565b606082471015620016ae5760405162461bcd60e51b815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f6044820152651c8818d85b1b60d21b606482015260840162000091565b600080866001600160a01b03168587604051620016cc9190620023b2565b60006040518083038185875af1925050503d80600081146200170b576040519150601f19603f3d011682016040523d82523d6000602084013e62001710565b606091505b50909250905062001724878383876200172f565b979650505050505050565b60608315620017a35782516000036200179b576001600160a01b0385163b6200179b5760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000604482015260640162000091565b5081620010c4565b620010c48383815115620017ba5781518083602001fd5b8060405162461bcd60e51b8152600401620000919190620023d0565b634e487b7160e01b600052604160045260246000fd5b604080519081016001600160401b0381118282101715620018115762001811620017d6565b60405290565b60405160c081016001600160401b0381118282101715620018115762001811620017d6565b604051608081016001600160401b0381118282101715620018115762001811620017d6565b604051601f8201601f191681016001600160401b03811182821017156200188c576200188c620017d6565b604052919050565b6001600160a01b0381168114620018aa57600080fd5b50565b80516001600160401b0381168114620018c557600080fd5b919050565b80516001600160601b0381168114620018c557600080fd5b600060e08284031215620018f557600080fd5b60405160e081016001600160401b03811182821017156200191a576200191a620017d6565b806040525080915082516200192f8162001894565b81526200193f60208401620018ad565b60208201526200195260408401620018ad565b60408201526200196560608401620018ad565b60608201526200197860808401620018ca565b608082015260a08301516200198d8162001894565b60a082015260c0830151620019a28162001894565b60c0919091015292915050565b805161ffff81168114620018c557600080fd5b805163ffffffff81168114620018c557600080fd5b600060a08284031215620019ea57600080fd5b60405160a081016001600160401b038111828210171562001a0f5762001a0f620017d6565b8060405250809150825162001a248162001894565b815262001a3460208401620019af565b6020820152604083015162001a498162001894565b604082015262001a5c60608401620019c2565b606082015262001a6f60808401620018ad565b60808201525092915050565b60006001600160401b0382111562001a975762001a97620017d6565b5060051b60200190565b600082601f83011262001ab357600080fd5b8151602062001acc62001ac68362001a7b565b62001861565b82815260069290921b8401810191818101908684111562001aec57600080fd5b8286015b8481101562001b46576040818903121562001b0b5760008081fd5b62001b15620017ec565b815162001b228162001894565b81528185015162001b338162001894565b8186015283529183019160400162001af0565b509695505050505050565b600082601f83011262001b6357600080fd5b8151602062001b7662001ac68362001a7b565b82815260059290921b8401810191818101908684111562001b9657600080fd5b8286015b8481101562001b4657805162001bb08162001894565b835291830191830162001b9a565b80518015158114620018c557600080fd5b80516001600160801b0381168114620018c557600080fd5b60006060828403121562001bfa57600080fd5b604051606081016001600160401b038111828210171562001c1f5762001c1f620017d6565b60405290508062001c308362001bbe565b815262001c406020840162001bcf565b602082015262001c536040840162001bcf565b60408201525092915050565b600082601f83011262001c7157600080fd5b8151602062001c8462001ac68362001a7b565b82815260c0928302850182019282820191908785111562001ca457600080fd5b8387015b8581101562001d465781818a03121562001cc25760008081fd5b62001ccc62001817565b815162001cd98162001894565b815262001ce8828701620018ad565b86820152604062001cfb818401620018ca565b90820152606062001d0e838201620019c2565b90820152608062001d21838201620019af565b9082015260a062001d3483820162001bbe565b90820152845292840192810162001ca8565b5090979650505050505050565b600082601f83011262001d6557600080fd5b8151602062001d7862001ac68362001a7b565b82815260079290921b8401810191818101908684111562001d9857600080fd5b8286015b8481101562001b46576080818903121562001db75760008081fd5b62001dc16200183c565b815162001dce8162001894565b815262001ddd828601620019c2565b85820152604062001df0818401620019c2565b90820152606062001e03838201620019af565b9082015283529183019160800162001d9c565b600082601f83011262001e2857600080fd5b8151602062001e3b62001ac68362001a7b565b82815260069290921b8401810191818101908684111562001e5b57600080fd5b8286015b8481101562001b46576040818903121562001e7a5760008081fd5b62001e84620017ec565b815162001e918162001894565b815262001ea0828601620019af565b8186015283529183019160400162001e5f565b600080600080600080600080610280898b03121562001ed157600080fd5b62001edd8a8a620018e2565b975062001eee8a60e08b01620019d7565b6101808a01519097506001600160401b038082111562001f0d57600080fd5b62001f1b8c838d0162001aa1565b97506101a08b015191508082111562001f3357600080fd5b62001f418c838d0162001b51565b965062001f538c6101c08d0162001be7565b95506102208b015191508082111562001f6b57600080fd5b62001f798c838d0162001c5f565b94506102408b015191508082111562001f9157600080fd5b62001f9f8c838d0162001d53565b93506102608b015191508082111562001fb757600080fd5b5062001fc68b828c0162001e16565b9150509295985092959890939650565b82516001600160a01b0390811682526020808501516001600160401b03908116828501526040808701518216818601526060808801518316818701526080808901516001600160601b03168188015260a0808a015187169088015260c0808a01518716908801528751861660e08801529387015161ffff16610100870152908601519093166101208501529184015163ffffffff16610140840152830151166101608201526101808101620010a5565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b600060018201620020c757620020c76200209c565b5060010190565b602080825282518282018190526000919060409081850190868401855b828110156200216257815180516001600160a01b03168552868101516001600160401b031687860152858101516001600160601b03168686015260608082015163ffffffff169086015260808082015161ffff169086015260a09081015115159085015260c09093019290850190600101620020eb565b5091979650505050505050565b602080825282518282018190526000919060409081850190868401855b828110156200216257815180516001600160a01b031685528681015163ffffffff9081168887015286820151168686015260609081015161ffff1690850152608090930192908501906001016200218c565b818103818111156200106a576200106a6200209c565b6000816200220657620022066200209c565b506000190190565b63ffffffff8181168382160190808211156200141257620014126200209c565b6000604080830163ffffffff8616845260208281860152818651808452606087019150828801935060005b818110156200228d57845180516001600160a01b0316845284015161ffff1684840152938301939185019160010162002259565b509098975050505050505050565b600060208284031215620022ae57600080fd5b8151620010a58162001894565b80820281158282048414176200106a576200106a6200209c565b600082620022f357634e487b7160e01b600052601260045260246000fd5b500490565b6001600160601b038281168282160390808211156200141257620014126200209c565b6000602082840312156200232e57600080fd5b5051919050565b81810360008312801583831316838312821617156200141257620014126200209c565b634e487b7160e01b600052603160045260246000fd5b6000602082840312156200238157600080fd5b620010a58262001bbe565b60005b83811015620023a95781810151838201526020016200238f565b50506000910152565b60008251620023c68184602087016200238c565b9190910192915050565b6020815260008251806020840152620023f18160408501602087016200238c565b601f01601f19169190910160400192915050565b60805160a05160c05160e05161010051610120516101405161016051615cbd6200250c600039600081816103600152818161140a0152612c2001526000818161033101528181611320015281816113880152818161185a015281816118c20152612bf801526000818161029d01528181610af901528181611c520152612b7001526000818161026d0152818161198d0152612b4601526000818161023e01528181610dfa0152818161160d01528181611706015281816121a501528181612b2101528181612d6b015261323c0152600081816102fd015281816117d20152612bc20152600081816102cd015281816122bd0152612b9701526000611b6c0152615cbd6000f3fe608060405234801561001057600080fd5b50600436106101f05760003560e01c806376f6ae761161010f578063b06d41bc116100a2578063e0351e1311610071578063e0351e13146108ee578063efeadb6d14610921578063eff7cc4814610934578063f2fde38b1461093c57600080fd5b8063b06d41bc146108b5578063c92b2832146108cb578063d09dc339146108de578063d3c7c2c7146108e657600080fd5b80638da5cb5b116100de5780638da5cb5b146107105780639a113c3614610721578063a7cd63b71461088d578063a7d3e02f146108a257600080fd5b806376f6ae76146106cf578063799c3a67146106e257806379ba5097146106f5578063856c8247146106fd57600080fd5b8063549e946f116101875780635d86f141116101565780635d86f141146105d45780635ebbd9f8146105e7578063704b6c02146105fa5780637437ff9f1461060d57600080fd5b8063549e946f1461056957806354b714681461057c57806354c8a4f31461059c578063599f6431146105af57600080fd5b80633a87ac53116101c35780633a87ac53146104bc5780633a9bf949146104d15780634120fccd146104e4578063546719cd1461050557600080fd5b806306285c69146101f55780631772047e146103a6578063181f5a771461045257806338724a951461049b575b600080fd5b6103906040805160e081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c08101919091526040518060e001604052807f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031681526020017f000000000000000000000000000000000000000000000000000000000000000067ffffffffffffffff1681526020017f000000000000000000000000000000000000000000000000000000000000000067ffffffffffffffff1681526020017f000000000000000000000000000000000000000000000000000000000000000067ffffffffffffffff1681526020017f00000000000000000000000000000000000000000000000000000000000000006bffffffffffffffffffffffff1681526020017f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031681526020017f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316815250905090565b60405161039d9190614822565b60405180910390f35b6104216103b43660046148b3565b6040805160608082018352600080835260208084018290529284018190526001600160a01b039490941684526010825292829020825193840183525463ffffffff80821685526401000000008204169184019190915268010000000000000000900461ffff169082015290565b60408051825163ffffffff9081168252602080850151909116908201529181015161ffff169082015260600161039d565b61048e6040518060400160405280601381526020017f45564d3245564d4f6e52616d7020312e302e300000000000000000000000000081525081565b60405161039d919061493e565b6104ae6104a9366004614963565b61094f565b60405190815260200161039d565b6104cf6104ca366004614b4a565b610c8d565b005b6104cf6104df366004614bef565b610ca3565b6104ec610cb7565b60405167ffffffffffffffff909116815260200161039d565b61050d610ceb565b60405161039d919081516fffffffffffffffffffffffffffffffff908116825260208084015163ffffffff1690830152604080840151151590830152606080840151821690830152608092830151169181019190915260a00190565b6104cf610577366004614c83565b610d9b565b6012546040516bffffffffffffffffffffffff909116815260200161039d565b6104cf6105aa366004614d20565b610f50565b6002546001600160a01b03165b6040516001600160a01b03909116815260200161039d565b6105bc6105e23660046148b3565b610f62565b6104cf6105f5366004614d7a565b610fc1565b6104cf6106083660046148b3565b611027565b6106c26040805160a081018252600080825260208201819052918101829052606081018290526080810191909152506040805160a0810182526005546001600160a01b038082168352740100000000000000000000000000000000000000009182900461ffff16602084015260065490811693830193909352820463ffffffff166060820152780100000000000000000000000000000000000000000000000090910467ffffffffffffffff16608082015290565b60405161039d9190614e58565b6104cf6106dd366004614eb2565b6110f1565b6104cf6106f0366004614f40565b6111a9565b6104cf61120f565b6104ec61070b3660046148b3565b6112f2565b6000546001600160a01b03166105bc565b61082761072f3660046148b3565b6040805160a081018252600080825260208201819052918101829052606081018290526080810191909152506001600160a01b03166000908152600f6020908152604091829020825160a08101845290546bffffffffffffffffffffffff811682526c01000000000000000000000000810467ffffffffffffffff169282019290925274010000000000000000000000000000000000000000820463ffffffff16928101929092527801000000000000000000000000000000000000000000000000810461ffff1660608301527a010000000000000000000000000000000000000000000000000000900460ff161515608082015290565b60405161039d9190600060a0820190506bffffffffffffffffffffffff835116825267ffffffffffffffff602084015116602083015263ffffffff604084015116604083015261ffff606084015116606083015260808301511515608083015292915050565b6108956113fa565b60405161039d919061505e565b6104ae6108b03660046150ab565b611406565b6108bd611d39565b60405161039d929190615159565b6104cf6108d936600461519b565b611e3d565b6104ae611ea5565b610895611eaf565b601254790100000000000000000000000000000000000000000000000000900460ff16604051901515815260200161039d565b6104cf61092f366004615209565b611f60565b6104cf611fe6565b6104cf61094a3660046148b3565b61227d565b600080600f8161096560808601606087016148b3565b6001600160a01b031681526020808201929092526040908101600020815160a08101835290546bffffffffffffffffffffffff811682526c01000000000000000000000000810467ffffffffffffffff169382019390935274010000000000000000000000000000000000000000830463ffffffff16918101919091527801000000000000000000000000000000000000000000000000820461ffff1660608201527a01000000000000000000000000000000000000000000000000000090910460ff16151560808201819052909150610a8c57610a4960808401606085016148b3565b6040517fa7499d200000000000000000000000000000000000000000000000000000000081526001600160a01b0390911660048201526024015b60405180910390fd5b60065460009081906001600160a01b031663ffdb4b37610ab260808801606089016148b3565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e084901b1681526001600160a01b03909116600482015267ffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001660248201526044016040805180830381865afa158015610b3d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b61919061524e565b91509150600083600001516bffffffffffffffffffffffff16670de0b6b3a7640000856020015167ffffffffffffffff16866060015161ffff16898060200190610bab9190615281565b610bb6929150615315565b604088015163ffffffff16610bd6610bd160808d018d615281565b61228e565b51610be1919061532c565b610beb919061532c565b610bf59190615315565b610c199077ffffffffffffffffffffffffffffffffffffffffffffffff8616615315565b610c23919061533f565b610c2d919061532c565b9050610c55610c4260808801606089016148b3565b84610c5060408a018a61537a565b612382565b610c7977ffffffffffffffffffffffffffffffffffffffffffffffff8516836125ce565b610c83919061532c565b9695505050505050565b610c95612607565b610c9f828261267d565b5050565b610cab612607565b610cb4816129dd565b50565b601254600090610ce690700100000000000000000000000000000000900467ffffffffffffffff1660016153e2565b905090565b6040805160a0810182526000808252602082018190529181018290526060810182905260808101919091526040805160a0810182526003546fffffffffffffffffffffffffffffffff808216835270010000000000000000000000000000000080830463ffffffff1660208501527401000000000000000000000000000000000000000090920460ff161515938301939093526004548084166060840152049091166080820152610ce690612c76565b6000546001600160a01b03163314801590610dc157506002546001600160a01b03163314155b15610df8576040517ffbdb8e5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316826001600160a01b03161480610e3f57506001600160a01b038116155b15610e76576040517f232cb97f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000610e80612d28565b1215610eb8576040517f02075e0000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152610c9f9082906001600160a01b038516906370a0823190602401602060405180830381865afa158015610f1b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610f3f9190615403565b6001600160a01b0385169190612de8565b610f58612607565b610c9f8282612e68565b6000610f6f600a83612fa3565b610fb0576040517fbf16aab60000000000000000000000000000000000000000000000000000000081526001600160a01b0383166004820152602401610a83565b610fbb600a83612fb8565b92915050565b6000546001600160a01b03163314801590610fe757506002546001600160a01b03163314155b1561101e576040517ffbdb8e5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610cb481612fcd565b6000546001600160a01b0316331480159061104d57506002546001600160a01b03163314155b15611084576040517ff6cd562000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600280547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0383169081179091556040519081527f8fe72c3e0020beb3234e76ae6676fa576fbfcae600af1c4fea44784cf0db329c906020015b60405180910390a150565b6000546001600160a01b0316331480159061111757506002546001600160a01b03163314155b1561114e576040517ffbdb8e5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610c9f8282808060200260200160405190810160405280939291908181526020016000905b8282101561119f576111906040830286013681900381019061541c565b81526020019060010190611173565b5050505050613103565b6000546001600160a01b031633148015906111cf57506002546001600160a01b03163314155b15611206576040517ffbdb8e5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610cb481613376565b6001546001600160a01b03163314611283576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610a83565b60008054337fffffffffffffffffffffffff0000000000000000000000000000000000000000808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6001600160a01b03811660009081526011602052604081205467ffffffffffffffff168015801561134b57507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031615155b15610fbb576040517f856c82470000000000000000000000000000000000000000000000000000000081526001600160a01b0384811660048301527f0000000000000000000000000000000000000000000000000000000000000000169063856c824790602401602060405180830381865afa1580156113cf573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906113f3919061545b565b9392505050565b6060610ce6600d613585565b60007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663397796f76040518163ffffffff1660e01b8152600401602060405180830381865afa158015611466573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061148a9190615478565b156114c1576040517fc148371500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6114cb8480615281565b9050602014611512576114de8480615281565b6040517f370d875f000000000000000000000000000000000000000000000000000000008152600401610a839291906154de565b600061151e8580615281565b81019061152b91906154f2565b90506001600160a01b038111806115425750600a81105b15611551576114de8580615281565b6000611563610bd16080880188615281565b905061158f6115756020880188615281565b835190915061158760408a018a61537a565b905087613592565b61160361159f604088018861537a565b808060200260200160405190810160405280939291908181526020016000905b828210156115eb576115dc6040830286013681900381019061550b565b815260200190600101906115bf565b50506006546001600160a01b031692506137b5915050565b6001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001661163d60808801606089016148b3565b6001600160a01b0316036116a1576012805486919060009061166e9084906bffffffffffffffffffffffff16615545565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055506117c0565b6006546001600160a01b03166241e5be6116c16080890160608a016148b3565b60405160e083901b7fffffffff000000000000000000000000000000000000000000000000000000001681526001600160a01b039182166004820152602481018990527f00000000000000000000000000000000000000000000000000000000000000009091166044820152606401602060405180830381865afa15801561174d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906117719190615403565b601280546000906117919084906bffffffffffffffffffffffff16615545565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055505b6012546bffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000081169116111561182d576040517fe5c7a49100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600160a01b03841660009081526011602052604090205467ffffffffffffffff1615801561188557507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031615155b1561197d576040517f856c82470000000000000000000000000000000000000000000000000000000081526001600160a01b0385811660048301527f0000000000000000000000000000000000000000000000000000000000000000169063856c824790602401602060405180830381865afa158015611909573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061192d919061545b565b6001600160a01b038516600090815260116020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001667ffffffffffffffff929092169190911790555b60006040518061018001604052807f000000000000000000000000000000000000000000000000000000000000000067ffffffffffffffff1681526020016012601081819054906101000a900467ffffffffffffffff166119dd9061556a565b825467ffffffffffffffff9182166101009390930a8381029083021990911617909255825260208083018a90526001600160a01b03891660408085018290526000918252601190925290812080546060909401939092611a3d911661556a565b91906101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905567ffffffffffffffff16815260200183600001518152602001600015158152602001846001600160a01b03168152602001888060200190611aa39190615281565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250505090825250602001611aea60408a018a61537a565b808060200260200160405190810160405280939291908181526020016000905b82821015611b3657611b276040830286013681900381019061550b565b81526020019060010190611b0a565b5050509183525050602001611b5160808a0160608b016148b3565b6001600160a01b0316815260006020909101529050611b90817f000000000000000000000000000000000000000000000000000000000000000061396e565b61016082015260005b611ba6604089018961537a565b9050811015611cf2576000611bbe60408a018a61537a565b83818110611bce57611bce615591565b905060400201803603810190611be4919061550b565b9050611bf38160000151610f62565b6001600160a01b0316639687544588611c0c8c80615281565b60208087015160408051928301815260008352517fffffffff0000000000000000000000000000000000000000000000000000000060e088901b168152611c7a959493927f0000000000000000000000000000000000000000000000000000000000000000916004016155c0565b6000604051808303816000875af1158015611c99573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052611cdf9190810190615618565b505080611ceb906156ca565b9050611b99565b507faffc45517195d6499808c643bd4a7b0ffeedf95bea5852840d7bfcf63f59e82181604051611d229190615746565b60405180910390a161016001519695505050505050565b6060600080611d486007613a78565b90508067ffffffffffffffff811115611d6357611d63614998565b604051908082528060200260200182016040528015611da857816020015b6040805180820190915260008082526020820152815260200190600190039081611d815790505b50925060005b81811015611e1a57600080611dc4600784613a83565b915091506040518060400160405280836001600160a01b031681526020018261ffff16815250868481518110611dfc57611dfc615591565b6020026020010181905250505080611e13906156ca565b9050611dae565b505060125491926c0100000000000000000000000090920463ffffffff16919050565b6000546001600160a01b03163314801590611e6357506002546001600160a01b03163314155b15611e9a576040517ff6cd562000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610cb4600382613aa1565b6000610ce6612d28565b60606000611ebd600a613c79565b67ffffffffffffffff811115611ed557611ed5614998565b604051908082528060200260200182016040528015611efe578160200160208202803683370190505b50905060005b8151811015611f5a57611f18600a82613c84565b50828281518110611f2b57611f2b615591565b60200260200101816001600160a01b03166001600160a01b03168152505080611f53906156ca565b9050611f04565b50919050565b611f68612607565b60128054821515790100000000000000000000000000000000000000000000000000027fffffffffffff00ffffffffffffffffffffffffffffffffffffffffffffffffff9091161790556040517fccf4daf6ab6430389f26b970595dab82a5881ad454770907e415ede27c8df032906110e690831515815260200190565b6000546001600160a01b0316331480159061200c57506002546001600160a01b03163314155b8015612020575061201e600733613c93565b155b15612057576040517f195db95800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6012546c01000000000000000000000000900463ffffffff1660008190036120ab576040517f990e30bf00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6012546bffffffffffffffffffffffff16818110156120f6576040517f8d0f71d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000612100612d28565b1215612138576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8060006121456007613a78565b905060005b8181101561223a57600080612160600784613a83565b9092509050600087612180836bffffffffffffffffffffffff8a16615315565b61218a919061533f565b90506121968187615883565b95506121da6001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016846bffffffffffffffffffffffff8416612de8565b6040516bffffffffffffffffffffffff821681526001600160a01b038416907f55fdec2aab60a41fa5abb106670eb1006f5aeaee1ba7afea2bc89b5b3ec7678f9060200160405180910390a250505080612233906156ca565b905061214a565b5050601280547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff929092169190911790555050565b612285612607565b610cb481613ca8565b60408051602081019091526000815260008290036122e45750604080516020810190915267ffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000168152610fbb565b7f97a657c90000000000000000000000000000000000000000000000000000000061230f83856158a8565b7fffffffff000000000000000000000000000000000000000000000000000000001614612368576040517f5247fdce00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61237582600481866158f0565b8101906113f3919061591a565b6000818082036123965760009150506125c6565b60005b818110156125c35760008585838181106123b5576123b5615591565b9050604002018036038101906123cb919061550b565b80516001600160a01b031660009081526010602090815260408083208151606081018352905463ffffffff8082168352640100000000820416938201939093526801000000000000000090920461ffff16908201819052929350911561252257825189906001600160a01b038c81169116146124cc5760065484516040517f4ab35b0b0000000000000000000000000000000000000000000000000000000081526001600160a01b039182166004820152911690634ab35b0b90602401602060405180830381865afa1580156124a5573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906124c9919061595c565b90505b620186a0836040015161ffff1661250a86602001518477ffffffffffffffffffffffffffffffffffffffffffffffff16613d8390919063ffffffff16565b6125149190615315565b61251e919061533f565b9150505b815160009061253e9063ffffffff16662386f26fc10000615315565b90506000836020015163ffffffff16662386f26fc1000061255f9190615315565b9050818310156125715781925061257d565b8083111561257d578092505b6125a177ffffffffffffffffffffffffffffffffffffffffffffffff8c16846125ce565b6125ab908961532c565b97505050505050806125bc906156ca565b9050612399565b50505b949350505050565b600077ffffffffffffffffffffffffffffffffffffffffffffffff83166125fd83670de0b6b3a7640000615315565b6113f3919061533f565b6000546001600160a01b0316331461267b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610a83565b565b60005b82518110156127de57600083828151811061269d5761269d615591565b602002602001015160000151905060008483815181106126bf576126bf615591565b60200260200101516020015190506126e182600a612fa390919063ffffffff16565b612722576040517f73913ebd0000000000000000000000000000000000000000000000000000000081526001600160a01b0383166004820152602401610a83565b6001600160a01b038116612737600a84612fb8565b6001600160a01b031614612777576040517f6cc7b99800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b612782600a83613db2565b156127cb57604080516001600160a01b038085168252831660208201527f987eb3c2f78454541205f72f34839b434c306c9eaf4922efd7c0c3060fdb2e4c910160405180910390a15b5050806127d7906156ca565b9050612680565b5060005b81518110156129d85760008282815181106127ff576127ff615591565b6020026020010151600001519050600083838151811061282157612821615591565b602002602001015160200151905060006001600160a01b0316826001600160a01b0316148061285757506001600160a01b038116155b1561288e576040517f6c2a418000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b806001600160a01b03166321df0da76040518163ffffffff1660e01b8152600401602060405180830381865afa1580156128cc573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906128f09190615977565b6001600160a01b0316826001600160a01b03161461293a576040517f6cc7b99800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b612946600a8383613dc7565b1561299357604080516001600160a01b038085168252831660208201527f95f865c2808f8b2a85eea2611db7843150ee7835ef1403f9755918a97d76933c910160405180910390a16129c5565b6040517f3caf458500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050806129d1906156ca565b90506127e2565b505050565b60408101516001600160a01b0316612a21576040517f35be3ac800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8051600580546020808501516001600160a01b039485167fffffffffffffffffffff00000000000000000000000000000000000000000000909316929092177401000000000000000000000000000000000000000061ffff909316830217909255604080850151600680546060808901516080808b0151958a167fffffffffffffffff0000000000000000000000000000000000000000000000009094169390931763ffffffff9091169096029590951777ffffffffffffffffffffffffffffffffffffffffffffffff16780100000000000000000000000000000000000000000000000067ffffffffffffffff9485160217909155825160e0810184527f0000000000000000000000000000000000000000000000000000000000000000871681527f00000000000000000000000000000000000000000000000000000000000000008316958101959095527f00000000000000000000000000000000000000000000000000000000000000008216858401527f0000000000000000000000000000000000000000000000000000000000000000909116928401929092527f00000000000000000000000000000000000000000000000000000000000000006bffffffffffffffffffffffff16918301919091527f0000000000000000000000000000000000000000000000000000000000000000831660a08301527f000000000000000000000000000000000000000000000000000000000000000090921660c082015290517fdd226617d8d287f40a64c54741bbcdc492b3e096ef16bc5273a18cb6ab85f124916110e6918490615994565b6040805160a081018252600080825260208201819052918101829052606081018290526080810191909152612d0482606001516fffffffffffffffffffffffffffffffff1683600001516fffffffffffffffffffffffffffffffff16846020015163ffffffff1642612ce89190615a69565b85608001516fffffffffffffffffffffffffffffffff16613ddd565b6fffffffffffffffffffffffffffffffff1682525063ffffffff4216602082015290565b6012546040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526000916bffffffffffffffffffffffff16907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316906370a0823190602401602060405180830381865afa158015612dba573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612dde9190615403565b610ce69190615a7c565b604080516001600160a01b038416602482015260448082018490528251808303909101815260649091019091526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fa9059cbb000000000000000000000000000000000000000000000000000000001790526129d8908490613e05565b60005b8251811015612ef9576000838281518110612e8857612e88615591565b60200260200101519050612ea681600d613f0490919063ffffffff16565b15612ee8576040516001600160a01b03821681527f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf75669060200160405180910390a15b50612ef2816156ca565b9050612e6b565b5060005b81518110156129d8576000828281518110612f1a57612f1a615591565b6020026020010151905060006001600160a01b0316816001600160a01b031603612f445750612f93565b612f4f600d82613f19565b15612f91576040516001600160a01b03821681527f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d89060200160405180910390a15b505b612f9c816156ca565b9050612efd565b60006113f3836001600160a01b038416613f2e565b60006113f3836001600160a01b038416613f3a565b60005b81518110156130d3576000828281518110612fed57612fed615591565b60209081029190910181015160408051606080820183528385015163ffffffff90811683528385015181168387019081529185015161ffff90811684860190815295516001600160a01b03166000908152601090975293909520915182549151945190931668010000000000000000027fffffffffffffffffffffffffffffffffffffffffffff0000ffffffffffffffff948616640100000000027fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000090921693909516929092179190911791909116919091179055506130cc816156ca565b9050612fd0565b507fcb0c5f472d325cf0c56953fc81870ddd80d0d3c9a3fbfe777002d75f380dfb81816040516110e69190615a9c565b80516040811115613140576040517fb5a10cfa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6012546c01000000000000000000000000900463ffffffff161580159061318e575060125463ffffffff6c010000000000000000000000008204166bffffffffffffffffffffffff90911610155b1561319b5761319b611fe6565b60006131a76007613a78565b90505b80156131e95760006131c86131c0600184615a69565b600790613a83565b5090506131d6600782613f46565b5050806131e290615b16565b90506131aa565b506000805b828110156132f757600084828151811061320a5761320a615591565b6020026020010151600001519050600085838151811061322c5761322c615591565b60200260200101516020015190507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316826001600160a01b0316148061328157506001600160a01b038216155b156132c3576040517f4de938d10000000000000000000000000000000000000000000000000000000081526001600160a01b0383166004820152602401610a83565b6132d360078361ffff8416613f5b565b506132e261ffff821685615b4b565b93505050806132f0906156ca565b90506131ee565b50601280547fffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffffff166c0100000000000000000000000063ffffffff8416021790556040517f8c337bff38141c507abd25c547606bdde78fe8c12e941ab613f3a565fea6cd24906133699083908690615b68565b60405180910390a1505050565b60005b815181101561355557600082828151811061339657613396615591565b6020908102919091018101516040805160a08082018352828401516bffffffffffffffffffffffff90811683528486015167ffffffffffffffff90811684880190815260608088015163ffffffff9081168789019081526080808b015161ffff908116948a01948552978b0151151590890190815299516001600160a01b03166000908152600f909b529790992095518654925197519151985115157a010000000000000000000000000000000000000000000000000000027fffffffffff00ffffffffffffffffffffffffffffffffffffffffffffffffffff999096167801000000000000000000000000000000000000000000000000027fffffffffffff0000ffffffffffffffffffffffffffffffffffffffffffffffff92909a167401000000000000000000000000000000000000000002919091167fffffffffffff000000000000ffffffffffffffffffffffffffffffffffffffff979093166c01000000000000000000000000027fffffffffffffffffffffffff000000000000000000000000000000000000000090921693169290921791909117939093169290921793909317919091161790555061354e816156ca565b9050613379565b507ffba339fca97870ffdfaedbae3745db5e6de1a6909dfd0e0dbb56917469ffe236816040516110e69190615b87565b606060006113f383613f71565b6001600160a01b0381166135d2576040517fa4ec747900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6005546001600160a01b03163314613616576040517f1c0a352900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60065474010000000000000000000000000000000000000000900463ffffffff168085111561367b576040517f869337890000000000000000000000000000000000000000000000000000000081526004810182905260248101869052604401610a83565b6006547801000000000000000000000000000000000000000000000000900467ffffffffffffffff168411156136dd576040517f4c4fc93a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60055474010000000000000000000000000000000000000000900461ffff16831115613735576040517f4c056b6a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601254790100000000000000000000000000000000000000000000000000900460ff16801561376c575061376a600d83613fcd565b155b156137ae576040517fd0d259760000000000000000000000000000000000000000000000000000000081526001600160a01b0383166004820152602401610a83565b5050505050565b81516000805b8281101561395a576000846001600160a01b031663d02641a08784815181106137e6576137e6615591565b6020908102919091010151516040517fffffffff0000000000000000000000000000000000000000000000000000000060e084901b1681526001600160a01b0390911660048201526024016040805180830381865afa15801561384d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906138719190615c1f565b51905077ffffffffffffffffffffffffffffffffffffffffffffffff81166000036138f2578582815181106138a8576138a8615591565b6020908102919091010151516040517f9a655f7b0000000000000000000000000000000000000000000000000000000081526001600160a01b039091166004820152602401610a83565b61393c86838151811061390757613907615591565b6020026020010151602001518277ffffffffffffffffffffffffffffffffffffffffffffffff16613d8390919063ffffffff16565b613946908461532c565b92505080613953906156ca565b90506137bb565b506139686003826000613fef565b50505050565b60008060001b828460200151856080015186606001518760e00151886101000151805190602001208961012001516040516020016139ac9190615c52565b604051602081830303815290604052805190602001208a60a001518b60c001518c61014001518d60400151604051602001613a5a9c9b9a999897969594939291909b8c5260208c019a909a5267ffffffffffffffff98891660408c01529690971660608a01526001600160a01b0394851660808a015292841660a089015260c088019190915260e0870152610100860152911515610120850152166101408301526101608201526101800190565b60405160208183030381529060405280519060200120905092915050565b6000610fbb8261433e565b6000808080613a928686614349565b909450925050505b9250929050565b8154600090613aca90700100000000000000000000000000000000900463ffffffff1642615a69565b90508015613b6c5760018301548354613b12916fffffffffffffffffffffffffffffffff80821692811691859170010000000000000000000000000000000090910416613ddd565b83546fffffffffffffffffffffffffffffffff919091167fffffffffffffffffffffffff0000000000000000000000000000000000000000909116177001000000000000000000000000000000004263ffffffff16021783555b60208201518354613b92916fffffffffffffffffffffffffffffffff9081169116614374565b83548351151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00ffffffff000000000000000000000000000000009091166fffffffffffffffffffffffffffffffff92831617178455602083015160408085015183167001000000000000000000000000000000000291909216176001850155517f9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c19906133699084908151151581526020808301516fffffffffffffffffffffffffffffffff90811691830191909152604092830151169181019190915260600190565b6000610fbb82613a78565b6000808080613a928686613a83565b60006113f3836001600160a01b03841661438a565b336001600160a01b03821603613d1a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610a83565b600180547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000670de0b6b3a76400006125fd8377ffffffffffffffffffffffffffffffffffffffffffffffff8616615315565b60006113f3836001600160a01b038416614396565b60006125c6846001600160a01b038516846143a2565b6000613dfc85613ded8486615315565b613df7908761532c565b614374565b95945050505050565b6000613e5a826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c6564815250856001600160a01b03166143b89092919063ffffffff16565b8051909150156129d85780806020019051810190613e789190615478565b6129d8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f742073756363656564000000000000000000000000000000000000000000006064820152608401610a83565b60006113f3836001600160a01b0384166143c7565b60006113f3836001600160a01b0384166144c1565b60006113f3838361438a565b60006113f38383614510565b60006113f3836001600160a01b03841661459a565b60006125c6846001600160a01b038516846145b7565b606081600001805480602002602001604051908101604052809291908181526020018280548015613fc157602002820191906000526020600020905b815481526020019060010190808311613fad575b50505050509050919050565b6001600160a01b038116600090815260018301602052604081205415156113f3565b825474010000000000000000000000000000000000000000900460ff161580614016575081155b1561402057505050565b825460018401546fffffffffffffffffffffffffffffffff8083169291169060009061406690700100000000000000000000000000000000900463ffffffff1642615a69565b9050801561412657818311156140a8576040517f9725942a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60018601546140e29083908590849070010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff16613ddd565b86547fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff167001000000000000000000000000000000004263ffffffff160217875592505b848210156141c3576001600160a01b038416614178576040517ff94ebcd10000000000000000000000000000000000000000000000000000000081526004810183905260248101869052604401610a83565b6040517f1a76572a00000000000000000000000000000000000000000000000000000000815260048101839052602481018690526001600160a01b0385166044820152606401610a83565b848310156142bc5760018681015470010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff169060009082906142079082615a69565b614211878a615a69565b61421b919061532c565b614225919061533f565b90506001600160a01b038616614271576040517f15279c080000000000000000000000000000000000000000000000000000000081526004810182905260248101869052604401610a83565b6040517fd0c8d23a00000000000000000000000000000000000000000000000000000000815260048101829052602481018690526001600160a01b0387166044820152606401610a83565b6142c68584615a69565b86547fffffffffffffffffffffffffffffffff00000000000000000000000000000000166fffffffffffffffffffffffffffffffff82161787556040518681529093507f1871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690a9060200160405180910390a1505050505050565b6000610fbb826145d4565b6000808061435785856145de565b600081815260029690960160205260409095205494959350505050565b600081831061438357816113f3565b5090919050565b60006113f383836145ea565b60006113f3838361459a565b60006125c684846001600160a01b0385166145b7565b60606125c68484600085614602565b600081815260018301602052604081205480156144b05760006143eb600183615a69565b85549091506000906143ff90600190615a69565b905081811461446457600086600001828154811061441f5761441f615591565b906000526020600020015490508087600001848154811061444257614442615591565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061447557614475615c65565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050610fbb565b6000915050610fbb565b5092915050565b600081815260018301602052604081205461450857508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610fbb565b506000610fbb565b6000818152600283016020526040812054801515806145345750614534848461438a565b6113f3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f456e756d657261626c654d61703a206e6f6e6578697374656e74206b657900006044820152606401610a83565b600081815260028301602052604081208190556113f3838361470e565b600082815260028401602052604081208290556125c6848461471a565b6000610fbb825490565b60006113f38383614726565b600081815260018301602052604081205415156113f3565b606082471015614694576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f60448201527f722063616c6c00000000000000000000000000000000000000000000000000006064820152608401610a83565b600080866001600160a01b031685876040516146b09190615c94565b60006040518083038185875af1925050503d80600081146146ed576040519150601f19603f3d011682016040523d82523d6000602084013e6146f2565b606091505b509150915061470387838387614750565b979650505050505050565b60006113f383836143c7565b60006113f383836144c1565b600082600001828154811061473d5761473d615591565b9060005260206000200154905092915050565b606083156147d95782516000036147d2576001600160a01b0385163b6147d2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401610a83565b50816125c6565b6125c683838151156147ee5781518083602001fd5b806040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610a83919061493e565b60e08101610fbb82846001600160a01b03808251168352602082015167ffffffffffffffff808216602086015280604085015116604086015280606085015116606086015250506bffffffffffffffffffffffff60808301511660808401528060a08301511660a08401528060c08301511660c0840152505050565b6001600160a01b0381168114610cb457600080fd5b6000602082840312156148c557600080fd5b81356113f38161489e565b60005b838110156148eb5781810151838201526020016148d3565b50506000910152565b6000815180845261490c8160208601602086016148d0565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6020815260006113f360208301846148f4565b600060a08284031215611f5a57600080fd5b60006020828403121561497557600080fd5b813567ffffffffffffffff81111561498c57600080fd5b6125c684828501614951565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040805190810167ffffffffffffffff811182821017156149ea576149ea614998565b60405290565b6040516080810167ffffffffffffffff811182821017156149ea576149ea614998565b60405160c0810167ffffffffffffffff811182821017156149ea576149ea614998565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715614a7d57614a7d614998565b604052919050565b600067ffffffffffffffff821115614a9f57614a9f614998565b5060051b60200190565b600082601f830112614aba57600080fd5b81356020614acf614aca83614a85565b614a36565b82815260069290921b84018101918181019086841115614aee57600080fd5b8286015b84811015614b3f5760408189031215614b0b5760008081fd5b614b136149c7565b8135614b1e8161489e565b815281850135614b2d8161489e565b81860152835291830191604001614af2565b509695505050505050565b60008060408385031215614b5d57600080fd5b823567ffffffffffffffff80821115614b7557600080fd5b614b8186838701614aa9565b93506020850135915080821115614b9757600080fd5b50614ba485828601614aa9565b9150509250929050565b803561ffff81168114614bc057600080fd5b919050565b803563ffffffff81168114614bc057600080fd5b67ffffffffffffffff81168114610cb457600080fd5b600060a08284031215614c0157600080fd5b60405160a0810181811067ffffffffffffffff82111715614c2457614c24614998565b6040528235614c328161489e565b8152614c4060208401614bae565b60208201526040830135614c538161489e565b6040820152614c6460608401614bc5565b60608201526080830135614c7781614bd9565b60808201529392505050565b60008060408385031215614c9657600080fd5b8235614ca18161489e565b91506020830135614cb18161489e565b809150509250929050565b600082601f830112614ccd57600080fd5b81356020614cdd614aca83614a85565b82815260059290921b84018101918181019086841115614cfc57600080fd5b8286015b84811015614b3f578035614d138161489e565b8352918301918301614d00565b60008060408385031215614d3357600080fd5b823567ffffffffffffffff80821115614d4b57600080fd5b614d5786838701614cbc565b93506020850135915080821115614d6d57600080fd5b50614ba485828601614cbc565b60006020808385031215614d8d57600080fd5b823567ffffffffffffffff811115614da457600080fd5b8301601f81018513614db557600080fd5b8035614dc3614aca82614a85565b81815260079190911b82018301908381019087831115614de257600080fd5b928401925b828410156147035760808489031215614e005760008081fd5b614e086149f0565b8435614e138161489e565b8152614e20858701614bc5565b868201526040614e31818701614bc5565b908201526060614e42868201614bae565b9082015282526080939093019290840190614de7565b60a08101610fbb82846001600160a01b0380825116835261ffff60208301511660208401528060408301511660408401525063ffffffff606082015116606083015267ffffffffffffffff60808201511660808301525050565b60008060208385031215614ec557600080fd5b823567ffffffffffffffff80821115614edd57600080fd5b818501915085601f830112614ef157600080fd5b813581811115614f0057600080fd5b8660208260061b8501011115614f1557600080fd5b60209290920196919550909350505050565b8015158114610cb457600080fd5b8035614bc081614f27565b60006020808385031215614f5357600080fd5b823567ffffffffffffffff811115614f6a57600080fd5b8301601f81018513614f7b57600080fd5b8035614f89614aca82614a85565b81815260c09182028301840191848201919088841115614fa857600080fd5b938501935b838510156150525780858a031215614fc55760008081fd5b614fcd614a13565b8535614fd88161489e565b815285870135614fe781614bd9565b818801526040868101356bffffffffffffffffffffffff8116811461500c5760008081fd5b90820152606061501d878201614bc5565b90820152608061502e878201614bae565b9082015260a061503f878201614f35565b9082015283529384019391850191614fad565b50979650505050505050565b6020808252825182820181905260009190848201906040850190845b8181101561509f5783516001600160a01b03168352928401929184019160010161507a565b50909695505050505050565b6000806000606084860312156150c057600080fd5b833567ffffffffffffffff8111156150d757600080fd5b6150e386828701614951565b9350506020840135915060408401356150fb8161489e565b809150509250925092565b600081518084526020808501945080840160005b8381101561514e57815180516001600160a01b0316885283015161ffff16838801526040909601959082019060010161511a565b509495945050505050565b60408152600061516c6040830185615106565b90508260208301529392505050565b80356fffffffffffffffffffffffffffffffff81168114614bc057600080fd5b6000606082840312156151ad57600080fd5b6040516060810181811067ffffffffffffffff821117156151d0576151d0614998565b60405282356151de81614f27565b81526151ec6020840161517b565b60208201526151fd6040840161517b565b60408201529392505050565b60006020828403121561521b57600080fd5b81356113f381614f27565b805177ffffffffffffffffffffffffffffffffffffffffffffffff81168114614bc057600080fd5b6000806040838503121561526157600080fd5b61526a83615226565b915061527860208401615226565b90509250929050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18436030181126152b657600080fd5b83018035915067ffffffffffffffff8211156152d157600080fd5b602001915036819003821315613a9a57600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8082028115828204841417610fbb57610fbb6152e6565b80820180821115610fbb57610fbb6152e6565b600082615375577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18436030181126153af57600080fd5b83018035915067ffffffffffffffff8211156153ca57600080fd5b6020019150600681901b3603821315613a9a57600080fd5b67ffffffffffffffff8181168382160190808211156144ba576144ba6152e6565b60006020828403121561541557600080fd5b5051919050565b60006040828403121561542e57600080fd5b6154366149c7565b82356154418161489e565b815261544f60208401614bae565b60208201529392505050565b60006020828403121561546d57600080fd5b81516113f381614bd9565b60006020828403121561548a57600080fd5b81516113f381614f27565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b6020815260006125c6602083018486615495565b60006020828403121561550457600080fd5b5035919050565b60006040828403121561551d57600080fd5b6155256149c7565b82356155308161489e565b81526020928301359281019290925250919050565b6bffffffffffffffffffffffff8181168382160190808211156144ba576144ba6152e6565b600067ffffffffffffffff808316818103615587576155876152e6565b6001019392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b6001600160a01b038716815260a0602082015260006155e360a083018789615495565b85604084015267ffffffffffffffff85166060840152828103608084015261560b81856148f4565b9998505050505050505050565b60006020828403121561562a57600080fd5b815167ffffffffffffffff8082111561564257600080fd5b818401915084601f83011261565657600080fd5b81518181111561566857615668614998565b61569960207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601614a36565b91508082528560208285010111156156b057600080fd5b6156c18160208401602086016148d0565b50949350505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036156fb576156fb6152e6565b5060010190565b600081518084526020808501945080840160005b8381101561514e57815180516001600160a01b031688528301518388015260409096019590820190600101615716565b6020815261576160208201835167ffffffffffffffff169052565b6000602083015161577e604084018267ffffffffffffffff169052565b506040830151606083015260608301516157a360808401826001600160a01b03169052565b50608083015167ffffffffffffffff811660a08401525060a083015160c083015260c08301516157d760e084018215159052565b5060e08301516101006157f4818501836001600160a01b03169052565b8085015191505061018061012081818601526158146101a08601846148f4565b92508086015190506101407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe086850301818701526158528483615702565b935080870151915050610160615872818701836001600160a01b03169052565b959095015193019290925250919050565b6bffffffffffffffffffffffff8281168282160390808211156144ba576144ba6152e6565b7fffffffff0000000000000000000000000000000000000000000000000000000081358181169160048510156158e85780818660040360031b1b83161692505b505092915050565b6000808585111561590057600080fd5b8386111561590d57600080fd5b5050820193919092039150565b60006020828403121561592c57600080fd5b6040516020810181811067ffffffffffffffff8211171561594f5761594f614998565b6040529135825250919050565b60006020828403121561596e57600080fd5b6113f382615226565b60006020828403121561598957600080fd5b81516113f38161489e565b6101808101615a1182856001600160a01b03808251168352602082015167ffffffffffffffff808216602086015280604085015116604086015280606085015116606086015250506bffffffffffffffffffffffff60808301511660808401528060a08301511660a08401528060c08301511660c0840152505050565b82516001600160a01b0390811660e0840152602084015161ffff16610100840152604084015116610120830152606083015163ffffffff16610140830152608083015167ffffffffffffffff166101608301526113f3565b81810381811115610fbb57610fbb6152e6565b81810360008312801583831316838312821617156144ba576144ba6152e6565b602080825282518282018190526000919060409081850190868401855b82811015615b0957815180516001600160a01b031685528681015163ffffffff9081168887015286820151168686015260609081015161ffff169085015260809093019290850190600101615ab9565b5091979650505050505050565b600081615b2557615b256152e6565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190565b63ffffffff8181168382160190808211156144ba576144ba6152e6565b63ffffffff831681526040602082015260006125c66040830184615106565b602080825282518282018190526000919060409081850190868401855b82811015615b0957815180516001600160a01b031685528681015167ffffffffffffffff1687860152858101516bffffffffffffffffffffffff168686015260608082015163ffffffff169086015260808082015161ffff169086015260a09081015115159085015260c09093019290850190600101615ba4565b600060408284031215615c3157600080fd5b615c396149c7565b615c4283615226565b8152602083015161544f81614bd9565b6020815260006113f36020830184615702565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b60008251615ca68184602087016148d0565b919091019291505056fea164736f6c6343000813000a",
}

var EVM2EVMOnRampABI = EVM2EVMOnRampMetaData.ABI

var EVM2EVMOnRampBin = EVM2EVMOnRampMetaData.Bin

func DeployEVM2EVMOnRamp(auth *bind.TransactOpts, backend bind.ContractBackend, staticConfig EVM2EVMOnRampStaticConfig, dynamicConfig EVM2EVMOnRampDynamicConfig, tokensAndPools []InternalPoolUpdate, allowlist []common.Address, rateLimiterConfig RateLimiterConfig, feeTokenConfigs []EVM2EVMOnRampFeeTokenConfigArgs, tokenTransferFeeConfigArgs []EVM2EVMOnRampTokenTransferFeeConfigArgs, nopsAndWeights []EVM2EVMOnRampNopAndWeight) (common.Address, *types.Transaction, *EVM2EVMOnRamp, error) {
	parsed, err := EVM2EVMOnRampMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(EVM2EVMOnRampBin), backend, staticConfig, dynamicConfig, tokensAndPools, allowlist, rateLimiterConfig, feeTokenConfigs, tokenTransferFeeConfigArgs, nopsAndWeights)
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

func (_EVM2EVMOnRamp *EVM2EVMOnRampCaller) GetAllowList(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _EVM2EVMOnRamp.contract.Call(opts, &out, "getAllowList")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_EVM2EVMOnRamp *EVM2EVMOnRampSession) GetAllowList() ([]common.Address, error) {
	return _EVM2EVMOnRamp.Contract.GetAllowList(&_EVM2EVMOnRamp.CallOpts)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampCallerSession) GetAllowList() ([]common.Address, error) {
	return _EVM2EVMOnRamp.Contract.GetAllowList(&_EVM2EVMOnRamp.CallOpts)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampCaller) GetAllowListEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _EVM2EVMOnRamp.contract.Call(opts, &out, "getAllowListEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_EVM2EVMOnRamp *EVM2EVMOnRampSession) GetAllowListEnabled() (bool, error) {
	return _EVM2EVMOnRamp.Contract.GetAllowListEnabled(&_EVM2EVMOnRamp.CallOpts)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampCallerSession) GetAllowListEnabled() (bool, error) {
	return _EVM2EVMOnRamp.Contract.GetAllowListEnabled(&_EVM2EVMOnRamp.CallOpts)
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

func (_EVM2EVMOnRamp *EVM2EVMOnRampTransactor) ApplyAllowListUpdates(opts *bind.TransactOpts, removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _EVM2EVMOnRamp.contract.Transact(opts, "applyAllowListUpdates", removes, adds)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampSession) ApplyAllowListUpdates(removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _EVM2EVMOnRamp.Contract.ApplyAllowListUpdates(&_EVM2EVMOnRamp.TransactOpts, removes, adds)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampTransactorSession) ApplyAllowListUpdates(removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _EVM2EVMOnRamp.Contract.ApplyAllowListUpdates(&_EVM2EVMOnRamp.TransactOpts, removes, adds)
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

func (_EVM2EVMOnRamp *EVM2EVMOnRampTransactor) SetAllowListEnabled(opts *bind.TransactOpts, enabled bool) (*types.Transaction, error) {
	return _EVM2EVMOnRamp.contract.Transact(opts, "setAllowListEnabled", enabled)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampSession) SetAllowListEnabled(enabled bool) (*types.Transaction, error) {
	return _EVM2EVMOnRamp.Contract.SetAllowListEnabled(&_EVM2EVMOnRamp.TransactOpts, enabled)
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampTransactorSession) SetAllowListEnabled(enabled bool) (*types.Transaction, error) {
	return _EVM2EVMOnRamp.Contract.SetAllowListEnabled(&_EVM2EVMOnRamp.TransactOpts, enabled)
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

type EVM2EVMOnRampAllowListAddIterator struct {
	Event *EVM2EVMOnRampAllowListAdd

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMOnRampAllowListAddIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMOnRampAllowListAdd)
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
		it.Event = new(EVM2EVMOnRampAllowListAdd)
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

func (it *EVM2EVMOnRampAllowListAddIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMOnRampAllowListAddIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMOnRampAllowListAdd struct {
	Sender common.Address
	Raw    types.Log
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampFilterer) FilterAllowListAdd(opts *bind.FilterOpts) (*EVM2EVMOnRampAllowListAddIterator, error) {

	logs, sub, err := _EVM2EVMOnRamp.contract.FilterLogs(opts, "AllowListAdd")
	if err != nil {
		return nil, err
	}
	return &EVM2EVMOnRampAllowListAddIterator{contract: _EVM2EVMOnRamp.contract, event: "AllowListAdd", logs: logs, sub: sub}, nil
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampFilterer) WatchAllowListAdd(opts *bind.WatchOpts, sink chan<- *EVM2EVMOnRampAllowListAdd) (event.Subscription, error) {

	logs, sub, err := _EVM2EVMOnRamp.contract.WatchLogs(opts, "AllowListAdd")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMOnRampAllowListAdd)
				if err := _EVM2EVMOnRamp.contract.UnpackLog(event, "AllowListAdd", log); err != nil {
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

func (_EVM2EVMOnRamp *EVM2EVMOnRampFilterer) ParseAllowListAdd(log types.Log) (*EVM2EVMOnRampAllowListAdd, error) {
	event := new(EVM2EVMOnRampAllowListAdd)
	if err := _EVM2EVMOnRamp.contract.UnpackLog(event, "AllowListAdd", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMOnRampAllowListEnabledSetIterator struct {
	Event *EVM2EVMOnRampAllowListEnabledSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMOnRampAllowListEnabledSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMOnRampAllowListEnabledSet)
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
		it.Event = new(EVM2EVMOnRampAllowListEnabledSet)
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

func (it *EVM2EVMOnRampAllowListEnabledSetIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMOnRampAllowListEnabledSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMOnRampAllowListEnabledSet struct {
	Enabled bool
	Raw     types.Log
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampFilterer) FilterAllowListEnabledSet(opts *bind.FilterOpts) (*EVM2EVMOnRampAllowListEnabledSetIterator, error) {

	logs, sub, err := _EVM2EVMOnRamp.contract.FilterLogs(opts, "AllowListEnabledSet")
	if err != nil {
		return nil, err
	}
	return &EVM2EVMOnRampAllowListEnabledSetIterator{contract: _EVM2EVMOnRamp.contract, event: "AllowListEnabledSet", logs: logs, sub: sub}, nil
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampFilterer) WatchAllowListEnabledSet(opts *bind.WatchOpts, sink chan<- *EVM2EVMOnRampAllowListEnabledSet) (event.Subscription, error) {

	logs, sub, err := _EVM2EVMOnRamp.contract.WatchLogs(opts, "AllowListEnabledSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMOnRampAllowListEnabledSet)
				if err := _EVM2EVMOnRamp.contract.UnpackLog(event, "AllowListEnabledSet", log); err != nil {
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

func (_EVM2EVMOnRamp *EVM2EVMOnRampFilterer) ParseAllowListEnabledSet(log types.Log) (*EVM2EVMOnRampAllowListEnabledSet, error) {
	event := new(EVM2EVMOnRampAllowListEnabledSet)
	if err := _EVM2EVMOnRamp.contract.UnpackLog(event, "AllowListEnabledSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMOnRampAllowListRemoveIterator struct {
	Event *EVM2EVMOnRampAllowListRemove

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMOnRampAllowListRemoveIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMOnRampAllowListRemove)
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
		it.Event = new(EVM2EVMOnRampAllowListRemove)
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

func (it *EVM2EVMOnRampAllowListRemoveIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMOnRampAllowListRemoveIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMOnRampAllowListRemove struct {
	Sender common.Address
	Raw    types.Log
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampFilterer) FilterAllowListRemove(opts *bind.FilterOpts) (*EVM2EVMOnRampAllowListRemoveIterator, error) {

	logs, sub, err := _EVM2EVMOnRamp.contract.FilterLogs(opts, "AllowListRemove")
	if err != nil {
		return nil, err
	}
	return &EVM2EVMOnRampAllowListRemoveIterator{contract: _EVM2EVMOnRamp.contract, event: "AllowListRemove", logs: logs, sub: sub}, nil
}

func (_EVM2EVMOnRamp *EVM2EVMOnRampFilterer) WatchAllowListRemove(opts *bind.WatchOpts, sink chan<- *EVM2EVMOnRampAllowListRemove) (event.Subscription, error) {

	logs, sub, err := _EVM2EVMOnRamp.contract.WatchLogs(opts, "AllowListRemove")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMOnRampAllowListRemove)
				if err := _EVM2EVMOnRamp.contract.UnpackLog(event, "AllowListRemove", log); err != nil {
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

func (_EVM2EVMOnRamp *EVM2EVMOnRampFilterer) ParseAllowListRemove(log types.Log) (*EVM2EVMOnRampAllowListRemove, error) {
	event := new(EVM2EVMOnRampAllowListRemove)
	if err := _EVM2EVMOnRamp.contract.UnpackLog(event, "AllowListRemove", log); err != nil {
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
	case _EVM2EVMOnRamp.abi.Events["AllowListAdd"].ID:
		return _EVM2EVMOnRamp.ParseAllowListAdd(log)
	case _EVM2EVMOnRamp.abi.Events["AllowListEnabledSet"].ID:
		return _EVM2EVMOnRamp.ParseAllowListEnabledSet(log)
	case _EVM2EVMOnRamp.abi.Events["AllowListRemove"].ID:
		return _EVM2EVMOnRamp.ParseAllowListRemove(log)
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

func (EVM2EVMOnRampAllowListAdd) Topic() common.Hash {
	return common.HexToHash("0x2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d8")
}

func (EVM2EVMOnRampAllowListEnabledSet) Topic() common.Hash {
	return common.HexToHash("0xccf4daf6ab6430389f26b970595dab82a5881ad454770907e415ede27c8df032")
}

func (EVM2EVMOnRampAllowListRemove) Topic() common.Hash {
	return common.HexToHash("0x800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf7566")
}

func (EVM2EVMOnRampCCIPSendRequested) Topic() common.Hash {
	return common.HexToHash("0xaffc45517195d6499808c643bd4a7b0ffeedf95bea5852840d7bfcf63f59e821")
}

func (EVM2EVMOnRampConfigSet) Topic() common.Hash {
	return common.HexToHash("0xdd226617d8d287f40a64c54741bbcdc492b3e096ef16bc5273a18cb6ab85f124")
}

func (EVM2EVMOnRampFeeConfigSet) Topic() common.Hash {
	return common.HexToHash("0xfba339fca97870ffdfaedbae3745db5e6de1a6909dfd0e0dbb56917469ffe236")
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
	return common.HexToHash("0xcb0c5f472d325cf0c56953fc81870ddd80d0d3c9a3fbfe777002d75f380dfb81")
}

func (_EVM2EVMOnRamp *EVM2EVMOnRamp) Address() common.Address {
	return _EVM2EVMOnRamp.address
}

type EVM2EVMOnRampInterface interface {
	CurrentRateLimiterState(opts *bind.CallOpts) (RateLimiterTokenBucket, error)

	GetAllowList(opts *bind.CallOpts) ([]common.Address, error)

	GetAllowListEnabled(opts *bind.CallOpts) (bool, error)

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

	ApplyAllowListUpdates(opts *bind.TransactOpts, removes []common.Address, adds []common.Address) (*types.Transaction, error)

	ApplyPoolUpdates(opts *bind.TransactOpts, removes []InternalPoolUpdate, adds []InternalPoolUpdate) (*types.Transaction, error)

	ForwardFromRouter(opts *bind.TransactOpts, message ClientEVM2AnyMessage, feeTokenAmount *big.Int, originalSender common.Address) (*types.Transaction, error)

	PayNops(opts *bind.TransactOpts) (*types.Transaction, error)

	SetAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error)

	SetAllowListEnabled(opts *bind.TransactOpts, enabled bool) (*types.Transaction, error)

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

	FilterAllowListAdd(opts *bind.FilterOpts) (*EVM2EVMOnRampAllowListAddIterator, error)

	WatchAllowListAdd(opts *bind.WatchOpts, sink chan<- *EVM2EVMOnRampAllowListAdd) (event.Subscription, error)

	ParseAllowListAdd(log types.Log) (*EVM2EVMOnRampAllowListAdd, error)

	FilterAllowListEnabledSet(opts *bind.FilterOpts) (*EVM2EVMOnRampAllowListEnabledSetIterator, error)

	WatchAllowListEnabledSet(opts *bind.WatchOpts, sink chan<- *EVM2EVMOnRampAllowListEnabledSet) (event.Subscription, error)

	ParseAllowListEnabledSet(log types.Log) (*EVM2EVMOnRampAllowListEnabledSet, error)

	FilterAllowListRemove(opts *bind.FilterOpts) (*EVM2EVMOnRampAllowListRemoveIterator, error)

	WatchAllowListRemove(opts *bind.WatchOpts, sink chan<- *EVM2EVMOnRampAllowListRemove) (event.Subscription, error)

	ParseAllowListRemove(log types.Log) (*EVM2EVMOnRampAllowListRemove, error)

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
