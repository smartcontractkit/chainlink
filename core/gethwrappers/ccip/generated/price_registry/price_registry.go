// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package price_registry

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

type AuthorizedCallersAuthorizedCallerArgs struct {
	AddedCallers   []common.Address
	RemovedCallers []common.Address
}

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

type IPriceRegistryTokenPriceFeedConfig struct {
	DataFeedAddress common.Address
	TokenDecimals   uint8
}

type InternalGasPriceUpdate struct {
	DestChainSelector uint64
	UsdPerUnitGas     *big.Int
}

type InternalPriceUpdates struct {
	TokenPriceUpdates []InternalTokenPriceUpdate
	GasPriceUpdates   []InternalGasPriceUpdate
}

type InternalRampTokenAmount struct {
	SourcePoolAddress []byte
	DestTokenAddress  []byte
	ExtraData         []byte
	Amount            *big.Int
}

type InternalTimestampedPackedUint224 struct {
	Value     *big.Int
	Timestamp uint32
}

type InternalTokenPriceUpdate struct {
	SourceToken common.Address
	UsdPerToken *big.Int
}

type PriceRegistryDestChainConfig struct {
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
	DefaultTxGasLimit                 uint32
	GasMultiplierWeiPerEth            uint64
	NetworkFeeUSDCents                uint32
	EnforceOutOfOrder                 bool
	ChainFamilySelector               [4]byte
}

type PriceRegistryDestChainConfigArgs struct {
	DestChainSelector uint64
	DestChainConfig   PriceRegistryDestChainConfig
}

type PriceRegistryPremiumMultiplierWeiPerEthArgs struct {
	Token                      common.Address
	PremiumMultiplierWeiPerEth uint64
}

type PriceRegistryStaticConfig struct {
	MaxFeeJuelsPerMsg  *big.Int
	LinkToken          common.Address
	StalenessThreshold uint32
}

type PriceRegistryTokenPriceFeedUpdate struct {
	SourceToken common.Address
	FeedConfig  IPriceRegistryTokenPriceFeedConfig
}

type PriceRegistryTokenTransferFeeConfig struct {
	MinFeeUSDCents    uint32
	MaxFeeUSDCents    uint32
	DeciBps           uint16
	DestGasOverhead   uint32
	DestBytesOverhead uint32
	IsEnabled         bool
}

type PriceRegistryTokenTransferFeeConfigArgs struct {
	DestChainSelector       uint64
	TokenTransferFeeConfigs []PriceRegistryTokenTransferFeeConfigSingleTokenArgs
}

type PriceRegistryTokenTransferFeeConfigRemoveArgs struct {
	DestChainSelector uint64
	Token             common.Address
}

type PriceRegistryTokenTransferFeeConfigSingleTokenArgs struct {
	Token                  common.Address
	TokenTransferFeeConfig PriceRegistryTokenTransferFeeConfig
}

var PriceRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"uint96\",\"name\":\"maxFeeJuelsPerMsg\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"linkToken\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"stalenessThreshold\",\"type\":\"uint32\"}],\"internalType\":\"structPriceRegistry.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"},{\"internalType\":\"address[]\",\"name\":\"priceUpdaters\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"feeTokens\",\"type\":\"address[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"dataFeedAddress\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"tokenDecimals\",\"type\":\"uint8\"}],\"internalType\":\"structIPriceRegistry.TokenPriceFeedConfig\",\"name\":\"feedConfig\",\"type\":\"tuple\"}],\"internalType\":\"structPriceRegistry.TokenPriceFeedUpdate[]\",\"name\":\"tokenPriceFeeds\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"minFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"deciBps\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"destBytesOverhead\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"}],\"internalType\":\"structPriceRegistry.TokenTransferFeeConfig\",\"name\":\"tokenTransferFeeConfig\",\"type\":\"tuple\"}],\"internalType\":\"structPriceRegistry.TokenTransferFeeConfigSingleTokenArgs[]\",\"name\":\"tokenTransferFeeConfigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structPriceRegistry.TokenTransferFeeConfigArgs[]\",\"name\":\"tokenTransferFeeConfigArgs\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"premiumMultiplierWeiPerEth\",\"type\":\"uint64\"}],\"internalType\":\"structPriceRegistry.PremiumMultiplierWeiPerEthArgs[]\",\"name\":\"premiumMultiplierWeiPerEthArgs\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint16\",\"name\":\"maxNumberOfTokensPerMsg\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxDataBytes\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerMsgGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerPayloadByte\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destDataAvailabilityOverheadGas\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerDataAvailabilityByte\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"destDataAvailabilityMultiplierBps\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"defaultTokenFeeUSDCents\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"defaultTokenDestGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"defaultTokenDestBytesOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"defaultTxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"gasMultiplierWeiPerEth\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"networkFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"enforceOutOfOrder\",\"type\":\"bool\"},{\"internalType\":\"bytes4\",\"name\":\"chainFamilySelector\",\"type\":\"bytes4\"}],\"internalType\":\"structPriceRegistry.DestChainConfig\",\"name\":\"destChainConfig\",\"type\":\"tuple\"}],\"internalType\":\"structPriceRegistry.DestChainConfigArgs[]\",\"name\":\"destChainConfigArgs\",\"type\":\"tuple[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chain\",\"type\":\"uint64\"}],\"name\":\"ChainNotSupported\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DataFeedValueOutOfUint224Range\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"}],\"name\":\"DestinationChainNotEnabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ExtraArgOutOfOrderExecutionMustBeTrue\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"destBytesOverhead\",\"type\":\"uint32\"}],\"name\":\"InvalidDestBytesOverhead\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"}],\"name\":\"InvalidDestChainConfig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedAddress\",\"type\":\"bytes\"}],\"name\":\"InvalidEVMAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidExtraArgsTag\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidStaticConfig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"msgFeeJuels\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeJuelsPerMsg\",\"type\":\"uint256\"}],\"name\":\"MessageFeeTooHigh\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MessageGasLimitTooHigh\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"maxSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actualSize\",\"type\":\"uint256\"}],\"name\":\"MessageTooLarge\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SourceTokenDataTooLarge\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"threshold\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timePassed\",\"type\":\"uint256\"}],\"name\":\"StaleGasPrice\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"TokenNotSupported\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"UnauthorizedCaller\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnsupportedNumberOfTokens\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"AuthorizedCallerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"AuthorizedCallerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint16\",\"name\":\"maxNumberOfTokensPerMsg\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxDataBytes\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerMsgGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerPayloadByte\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destDataAvailabilityOverheadGas\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerDataAvailabilityByte\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"destDataAvailabilityMultiplierBps\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"defaultTokenFeeUSDCents\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"defaultTokenDestGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"defaultTokenDestBytesOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"defaultTxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"gasMultiplierWeiPerEth\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"networkFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"enforceOutOfOrder\",\"type\":\"bool\"},{\"internalType\":\"bytes4\",\"name\":\"chainFamilySelector\",\"type\":\"bytes4\"}],\"indexed\":false,\"internalType\":\"structPriceRegistry.DestChainConfig\",\"name\":\"destChainConfig\",\"type\":\"tuple\"}],\"name\":\"DestChainAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint16\",\"name\":\"maxNumberOfTokensPerMsg\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxDataBytes\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerMsgGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerPayloadByte\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destDataAvailabilityOverheadGas\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerDataAvailabilityByte\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"destDataAvailabilityMultiplierBps\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"defaultTokenFeeUSDCents\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"defaultTokenDestGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"defaultTokenDestBytesOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"defaultTxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"gasMultiplierWeiPerEth\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"networkFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"enforceOutOfOrder\",\"type\":\"bool\"},{\"internalType\":\"bytes4\",\"name\":\"chainFamilySelector\",\"type\":\"bytes4\"}],\"indexed\":false,\"internalType\":\"structPriceRegistry.DestChainConfig\",\"name\":\"destChainConfig\",\"type\":\"tuple\"}],\"name\":\"DestChainConfigUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"}],\"name\":\"FeeTokenAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"}],\"name\":\"FeeTokenRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"premiumMultiplierWeiPerEth\",\"type\":\"uint64\"}],\"name\":\"PremiumMultiplierWeiPerEthUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"dataFeedAddress\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"tokenDecimals\",\"type\":\"uint8\"}],\"indexed\":false,\"internalType\":\"structIPriceRegistry.TokenPriceFeedConfig\",\"name\":\"priceFeedConfig\",\"type\":\"tuple\"}],\"name\":\"PriceFeedPerTokenUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"TokenTransferFeeConfigDeleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"minFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"deciBps\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"destBytesOverhead\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"}],\"indexed\":false,\"internalType\":\"structPriceRegistry.TokenTransferFeeConfig\",\"name\":\"tokenTransferFeeConfig\",\"type\":\"tuple\"}],\"name\":\"TokenTransferFeeConfigUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"UsdPerTokenUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"destChain\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"UsdPerUnitGasUpdated\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address[]\",\"name\":\"addedCallers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"removedCallers\",\"type\":\"address[]\"}],\"internalType\":\"structAuthorizedCallers.AuthorizedCallerArgs\",\"name\":\"authorizedCallerArgs\",\"type\":\"tuple\"}],\"name\":\"applyAuthorizedCallerUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint16\",\"name\":\"maxNumberOfTokensPerMsg\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxDataBytes\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerMsgGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerPayloadByte\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destDataAvailabilityOverheadGas\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerDataAvailabilityByte\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"destDataAvailabilityMultiplierBps\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"defaultTokenFeeUSDCents\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"defaultTokenDestGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"defaultTokenDestBytesOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"defaultTxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"gasMultiplierWeiPerEth\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"networkFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"enforceOutOfOrder\",\"type\":\"bool\"},{\"internalType\":\"bytes4\",\"name\":\"chainFamilySelector\",\"type\":\"bytes4\"}],\"internalType\":\"structPriceRegistry.DestChainConfig\",\"name\":\"destChainConfig\",\"type\":\"tuple\"}],\"internalType\":\"structPriceRegistry.DestChainConfigArgs[]\",\"name\":\"destChainConfigArgs\",\"type\":\"tuple[]\"}],\"name\":\"applyDestChainConfigUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"feeTokensToAdd\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"feeTokensToRemove\",\"type\":\"address[]\"}],\"name\":\"applyFeeTokensUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"premiumMultiplierWeiPerEth\",\"type\":\"uint64\"}],\"internalType\":\"structPriceRegistry.PremiumMultiplierWeiPerEthArgs[]\",\"name\":\"premiumMultiplierWeiPerEthArgs\",\"type\":\"tuple[]\"}],\"name\":\"applyPremiumMultiplierWeiPerEthUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"minFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"deciBps\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"destBytesOverhead\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"}],\"internalType\":\"structPriceRegistry.TokenTransferFeeConfig\",\"name\":\"tokenTransferFeeConfig\",\"type\":\"tuple\"}],\"internalType\":\"structPriceRegistry.TokenTransferFeeConfigSingleTokenArgs[]\",\"name\":\"tokenTransferFeeConfigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structPriceRegistry.TokenTransferFeeConfigArgs[]\",\"name\":\"tokenTransferFeeConfigArgs\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"internalType\":\"structPriceRegistry.TokenTransferFeeConfigRemoveArgs[]\",\"name\":\"tokensToUseDefaultFeeConfigs\",\"type\":\"tuple[]\"}],\"name\":\"applyTokenTransferFeeConfigUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"fromToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"fromTokenAmount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"toToken\",\"type\":\"address\"}],\"name\":\"convertTokenAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllAuthorizedCallers\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"}],\"name\":\"getDestChainConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint16\",\"name\":\"maxNumberOfTokensPerMsg\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxDataBytes\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerMsgGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerPayloadByte\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destDataAvailabilityOverheadGas\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerDataAvailabilityByte\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"destDataAvailabilityMultiplierBps\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"defaultTokenFeeUSDCents\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"defaultTokenDestGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"defaultTokenDestBytesOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"defaultTxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"gasMultiplierWeiPerEth\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"networkFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"enforceOutOfOrder\",\"type\":\"bool\"},{\"internalType\":\"bytes4\",\"name\":\"chainFamilySelector\",\"type\":\"bytes4\"}],\"internalType\":\"structPriceRegistry.DestChainConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"}],\"name\":\"getDestinationChainGasPrice\",\"outputs\":[{\"components\":[{\"internalType\":\"uint224\",\"name\":\"value\",\"type\":\"uint224\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"}],\"internalType\":\"structInternal.TimestampedPackedUint224\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFeeTokens\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getPremiumMultiplierWeiPerEth\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"premiumMultiplierWeiPerEth\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStaticConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint96\",\"name\":\"maxFeeJuelsPerMsg\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"linkToken\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"stalenessThreshold\",\"type\":\"uint32\"}],\"internalType\":\"structPriceRegistry.StaticConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"}],\"name\":\"getTokenAndGasPrices\",\"outputs\":[{\"internalType\":\"uint224\",\"name\":\"tokenPrice\",\"type\":\"uint224\"},{\"internalType\":\"uint224\",\"name\":\"gasPriceValue\",\"type\":\"uint224\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getTokenPrice\",\"outputs\":[{\"components\":[{\"internalType\":\"uint224\",\"name\":\"value\",\"type\":\"uint224\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"}],\"internalType\":\"structInternal.TimestampedPackedUint224\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getTokenPriceFeedConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"dataFeedAddress\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"tokenDecimals\",\"type\":\"uint8\"}],\"internalType\":\"structIPriceRegistry.TokenPriceFeedConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"tokens\",\"type\":\"address[]\"}],\"name\":\"getTokenPrices\",\"outputs\":[{\"components\":[{\"internalType\":\"uint224\",\"name\":\"value\",\"type\":\"uint224\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"}],\"internalType\":\"structInternal.TimestampedPackedUint224[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getTokenTransferFeeConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"minFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"deciBps\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"destBytesOverhead\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"}],\"internalType\":\"structPriceRegistry.TokenTransferFeeConfig\",\"name\":\"tokenTransferFeeConfig\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMTokenAmount[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structClient.EVM2AnyMessage\",\"name\":\"message\",\"type\":\"tuple\"}],\"name\":\"getValidatedFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"feeTokenAmount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getValidatedTokenPrice\",\"outputs\":[{\"internalType\":\"uint224\",\"name\":\"\",\"type\":\"uint224\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"feeTokenAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"name\":\"processMessageArgs\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"msgFeeJuels\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isOutOfOrderExecution\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"convertedExtraArgs\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"internalType\":\"uint224\",\"name\":\"usdPerToken\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.TokenPriceUpdate[]\",\"name\":\"tokenPriceUpdates\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint224\",\"name\":\"usdPerUnitGas\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.GasPriceUpdate[]\",\"name\":\"gasPriceUpdates\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.PriceUpdates\",\"name\":\"priceUpdates\",\"type\":\"tuple\"}],\"name\":\"updatePrices\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"dataFeedAddress\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"tokenDecimals\",\"type\":\"uint8\"}],\"internalType\":\"structIPriceRegistry.TokenPriceFeedConfig\",\"name\":\"feedConfig\",\"type\":\"tuple\"}],\"internalType\":\"structPriceRegistry.TokenPriceFeedUpdate[]\",\"name\":\"tokenPriceFeedUpdates\",\"type\":\"tuple[]\"}],\"name\":\"updateTokenPriceFeeds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structInternal.RampTokenAmount[]\",\"name\":\"rampTokenAmounts\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMTokenAmount[]\",\"name\":\"sourceTokenAmounts\",\"type\":\"tuple[]\"}],\"name\":\"validatePoolReturnData\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60e06040523480156200001157600080fd5b5060405162006a4238038062006a4283398101604081905262000034916200188e565b8533806000816200008c5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000bf57620000bf8162000207565b5050604080518082018252838152815160008152602080820190935291810191909152620000ee9150620002b2565b5060208701516001600160a01b0316158062000112575086516001600160601b0316155b80620001265750604087015163ffffffff16155b15620001455760405163d794ef9560e01b815260040160405180910390fd5b6020878101516001600160a01b031660a05287516001600160601b031660805260408089015163ffffffff1660c05280516000815291820190526200018c90869062000401565b620001978462000549565b620001a2816200061a565b620001ad8262000a9f565b60408051600080825260208201909252620001fa91859190620001f3565b6040805180820190915260008082526020820152815260200190600190039081620001cb5790505b5062000b6b565b5050505050505062001b4c565b336001600160a01b03821603620002615760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000083565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b602081015160005b815181101562000342576000828281518110620002db57620002db620019ad565b60209081029190910101519050620002f560028262000ea4565b1562000338576040516001600160a01b03821681527fc3803387881faad271c47728894e3e36fac830ffc8602ca6fc07733cbda775809060200160405180910390a15b50600101620002ba565b50815160005b8151811015620003fb576000828281518110620003695762000369620019ad565b6020026020010151905060006001600160a01b0316816001600160a01b031603620003a7576040516342bcdf7f60e11b815260040160405180910390fd5b620003b460028262000ec4565b506040516001600160a01b03821681527feb1b9b92e50b7f88f9ff25d56765095ac6e91540eee214906f4036a908ffbdef9060200160405180910390a15060010162000348565b50505050565b60005b8251811015620004a25762000440838281518110620004275762000427620019ad565b6020026020010151600a62000ec460201b90919060201c565b1562000499578281815181106200045b576200045b620019ad565b60200260200101516001600160a01b03167fdf1b1bd32a69711488d71554706bb130b1fc63a5fa1a2cd85e8440f84065ba2360405160405180910390a25b60010162000404565b5060005b81518110156200054457620004e2828281518110620004c957620004c9620019ad565b6020026020010151600a62000edb60201b90919060201c565b156200053b57818181518110620004fd57620004fd620019ad565b60200260200101516001600160a01b03167f1795838dc8ab2ffc5f431a1729a6afa0b587f982f7b2be0b9d7187a1ef547f9160405160405180910390a25b600101620004a6565b505050565b60005b8151811015620006165760008282815181106200056d576200056d620019ad565b6020908102919091018101518051818301516001600160a01b0380831660008181526006875260409081902084518154868a018051929096166001600160a81b03199091168117600160a01b60ff9384160217909255825191825293519093169683019690965293955091939092917f08a5f7f5bb38a81d8e43aca13ecd76431dbf8816ae4699affff7b00b2fc1c464910160405180910390a25050508060010190506200054c565b5050565b60005b8151811015620006165760008282815181106200063e576200063e620019ad565b6020026020010151905060008383815181106200065f576200065f620019ad565b6020026020010151600001519050600082602001519050816001600160401b03166000148062000698575061018081015163ffffffff16155b80620006ba57506102008101516001600160e01b031916630a04b54b60e21b14155b80620006d75750602063ffffffff1681610160015163ffffffff16105b80620006f75750806060015163ffffffff1681610180015163ffffffff16115b15620007225760405163c35aa79d60e01b81526001600160401b038316600482015260240162000083565b6001600160401b038216600090815260086020526040812060010154600160a81b900460e01b6001600160e01b0319169003620007a257816001600160401b03167fa937382a486d993de71c220bc8b559242deb4e286a353fa732330b4aa7d1357782604051620007949190620019c3565b60405180910390a2620007e6565b816001600160401b03167fa7b607fc10d28a1caf39ab7d27f4c94945db708a576d572781a455c5894fad9382604051620007dd9190620019c3565b60405180910390a25b8060086000846001600160401b03166001600160401b0316815260200190815260200160002060008201518160000160006101000a81548160ff02191690831515021790555060208201518160000160016101000a81548161ffff021916908361ffff16021790555060408201518160000160036101000a81548163ffffffff021916908363ffffffff16021790555060608201518160000160076101000a81548163ffffffff021916908363ffffffff160217905550608082015181600001600b6101000a81548163ffffffff021916908363ffffffff16021790555060a082015181600001600f6101000a81548161ffff021916908361ffff16021790555060c08201518160000160116101000a81548163ffffffff021916908363ffffffff16021790555060e08201518160000160156101000a81548161ffff021916908361ffff1602179055506101008201518160000160176101000a81548161ffff021916908361ffff1602179055506101208201518160000160196101000a81548161ffff021916908361ffff16021790555061014082015181600001601b6101000a81548163ffffffff021916908363ffffffff1602179055506101608201518160010160006101000a81548163ffffffff021916908363ffffffff1602179055506101808201518160010160046101000a81548163ffffffff021916908363ffffffff1602179055506101a08201518160010160086101000a8154816001600160401b0302191690836001600160401b031602179055506101c08201518160010160106101000a81548163ffffffff021916908363ffffffff1602179055506101e08201518160010160146101000a81548160ff0219169083151502179055506102008201518160010160156101000a81548163ffffffff021916908360e01c02179055509050505050508060010190506200061d565b60005b81518110156200061657600082828151811062000ac35762000ac3620019ad565b6020026020010151600001519050600083838151811062000ae85762000ae8620019ad565b6020908102919091018101518101516001600160a01b03841660008181526007845260409081902080546001600160401b0319166001600160401b0385169081179091559051908152919350917fbb77da6f7210cdd16904228a9360133d1d7dfff99b1bc75f128da5b53e28f97d910160405180910390a2505060010162000aa2565b60005b825181101562000dde57600083828151811062000b8f5762000b8f620019ad565b6020026020010151905060008160000151905060005b82602001515181101562000dcf5760008360200151828151811062000bce5762000bce620019ad565b602002602001015160200151905060008460200151838151811062000bf75762000bf7620019ad565b6020026020010151600001519050602063ffffffff16826080015163ffffffff16101562000c565760808201516040516312766e0160e11b81526001600160a01b038316600482015263ffffffff909116602482015260440162000083565b6001600160401b03841660008181526009602090815260408083206001600160a01b0386168085529083529281902086518154938801518389015160608a015160808b015160a08c01511515600160901b0260ff60901b1963ffffffff928316600160701b021664ffffffffff60701b199383166a01000000000000000000000263ffffffff60501b1961ffff90961668010000000000000000029590951665ffffffffffff60401b19968416640100000000026001600160401b0319909b16939097169290921798909817939093169390931717919091161792909217909155519091907f94967ae9ea7729ad4f54021c1981765d2b1d954f7c92fbec340aa0a54f46b8b59062000dbc908690600060c08201905063ffffffff80845116835280602085015116602084015261ffff60408501511660408401528060608501511660608401528060808501511660808401525060a0830151151560a083015292915050565b60405180910390a3505060010162000ba5565b50505080600101905062000b6e565b5060005b81518110156200054457600082828151811062000e035762000e03620019ad565b6020026020010151600001519050600083838151811062000e285762000e28620019ad565b6020908102919091018101518101516001600160401b03841660008181526009845260408082206001600160a01b038516808452955280822080546001600160981b03191690555192945090917f4de5b1bcbca6018c11303a2c3f4a4b4f22a1c741d8c4ba430d246ac06c5ddf8b9190a3505060010162000de2565b600062000ebb836001600160a01b03841662000ef2565b90505b92915050565b600062000ebb836001600160a01b03841662000ff6565b600062000ebb836001600160a01b03841662001048565b6000818152600183016020526040812054801562000feb57600062000f1960018362001b14565b855490915060009062000f2f9060019062001b14565b905081811462000f9b57600086600001828154811062000f535762000f53620019ad565b906000526020600020015490508087600001848154811062000f795762000f79620019ad565b6000918252602080832090910192909255918252600188019052604090208390555b855486908062000faf5762000faf62001b36565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505062000ebe565b600091505062000ebe565b60008181526001830160205260408120546200103f5750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915562000ebe565b50600062000ebe565b6000818152600183016020526040812054801562000feb5760006200106f60018362001b14565b8554909150600090620010859060019062001b14565b905080821462000f9b57600086600001828154811062000f535762000f53620019ad565b634e487b7160e01b600052604160045260246000fd5b604080519081016001600160401b0381118282101715620010e457620010e4620010a9565b60405290565b60405160c081016001600160401b0381118282101715620010e457620010e4620010a9565b60405161022081016001600160401b0381118282101715620010e457620010e4620010a9565b604051601f8201601f191681016001600160401b0381118282101715620011605762001160620010a9565b604052919050565b80516001600160a01b03811681146200118057600080fd5b919050565b805163ffffffff811681146200118057600080fd5b600060608284031215620011ad57600080fd5b604051606081016001600160401b0381118282101715620011d257620011d2620010a9565b604052825190915081906001600160601b0381168114620011f257600080fd5b8152620012026020840162001168565b6020820152620012156040840162001185565b60408201525092915050565b60006001600160401b038211156200123d576200123d620010a9565b5060051b60200190565b600082601f8301126200125957600080fd5b81516020620012726200126c8362001221565b62001135565b8083825260208201915060208460051b8701019350868411156200129557600080fd5b602086015b84811015620012bc57620012ae8162001168565b83529183019183016200129a565b509695505050505050565b600082601f830112620012d957600080fd5b81516020620012ec6200126c8362001221565b828152606092830285018201928282019190878511156200130c57600080fd5b8387015b858110156200139f57808903828112156200132b5760008081fd5b62001335620010bf565b620013408362001168565b8152604080601f1984011215620013575760008081fd5b62001361620010bf565b92506200137088850162001168565b835283015160ff81168114620013865760008081fd5b8288015280870191909152845292840192810162001310565b5090979650505050505050565b80516001600160401b03811681146200118057600080fd5b805161ffff811681146200118057600080fd5b805180151581146200118057600080fd5b600082601f830112620013fa57600080fd5b815160206200140d6200126c8362001221565b82815260059290921b840181019181810190868411156200142d57600080fd5b8286015b84811015620012bc5780516001600160401b03808211156200145257600080fd5b908801906040601f19838c0381018213156200146d57600080fd5b62001477620010bf565b62001484898601620013ac565b815282850151848111156200149857600080fd5b8086019550508c603f860112620014ae57600080fd5b888501519350620014c36200126c8562001221565b84815260e09094028501830193898101908e861115620014e257600080fd5b958401955b85871015620015bb57868f0360e08112156200150257600080fd5b6200150c620010bf565b620015178962001168565b815260c086830112156200152a57600080fd5b62001534620010ea565b9150620015438d8a0162001185565b825262001552878a0162001185565b8d8301526200156460608a01620013c4565b878301526200157660808a0162001185565b60608301526200158960a08a0162001185565b60808301526200159c60c08a01620013d7565b60a0830152808d0191909152825260e09690960195908a0190620014e7565b828b01525087525050509284019250830162001431565b600082601f830112620015e457600080fd5b81516020620015f76200126c8362001221565b82815260069290921b840181019181810190868411156200161757600080fd5b8286015b84811015620012bc5760408189031215620016365760008081fd5b62001640620010bf565b6200164b8262001168565b81526200165a858301620013ac565b818601528352918301916040016200161b565b80516001600160e01b0319811681146200118057600080fd5b600082601f8301126200169857600080fd5b81516020620016ab6200126c8362001221565b8281526102409283028501820192828201919087851115620016cc57600080fd5b8387015b858110156200139f5780890382811215620016eb5760008081fd5b620016f5620010bf565b6200170083620013ac565b815261022080601f1984011215620017185760008081fd5b620017226200110f565b925062001731888501620013d7565b8352604062001742818601620013c4565b8985015260606200175581870162001185565b82860152608091506200176a82870162001185565b9085015260a06200177d86820162001185565b8286015260c0915062001792828701620013c4565b9085015260e0620017a586820162001185565b828601526101009150620017bb828701620013c4565b90850152610120620017cf868201620013c4565b828601526101409150620017e5828701620013c4565b90850152610160620017f986820162001185565b8286015261018091506200180f82870162001185565b908501526101a06200182386820162001185565b828601526101c0915062001839828701620013ac565b908501526101e06200184d86820162001185565b82860152610200915062001863828701620013d7565b90850152620018748583016200166d565b9084015250808701919091528452928401928101620016d0565b6000806000806000806000610120888a031215620018ab57600080fd5b620018b789896200119a565b60608901519097506001600160401b0380821115620018d557600080fd5b620018e38b838c0162001247565b975060808a0151915080821115620018fa57600080fd5b620019088b838c0162001247565b965060a08a01519150808211156200191f57600080fd5b6200192d8b838c01620012c7565b955060c08a01519150808211156200194457600080fd5b620019528b838c01620013e8565b945060e08a01519150808211156200196957600080fd5b620019778b838c01620015d2565b93506101008a01519150808211156200198f57600080fd5b506200199e8a828b0162001686565b91505092959891949750929550565b634e487b7160e01b600052603260045260246000fd5b81511515815261022081016020830151620019e4602084018261ffff169052565b506040830151620019fd604084018263ffffffff169052565b50606083015162001a16606084018263ffffffff169052565b50608083015162001a2f608084018263ffffffff169052565b5060a083015162001a4660a084018261ffff169052565b5060c083015162001a5f60c084018263ffffffff169052565b5060e083015162001a7660e084018261ffff169052565b506101008381015161ffff9081169184019190915261012080850151909116908301526101408084015163ffffffff9081169184019190915261016080850151821690840152610180808501518216908401526101a0808501516001600160401b0316908401526101c080850151909116908301526101e080840151151590830152610200928301516001600160e01b031916929091019190915290565b8181038181111562000ebe57634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052603160045260246000fd5b60805160a05160c051614e9c62001ba6600039600081816102d901528181611ad30152611b3c01526000818161029d0152818161104e01526110ae015260008181610269015281816110d701526111470152614e9c6000f3fe608060405234801561001057600080fd5b50600436106101b85760003560e01c80637afac322116100f9578063cc88924c11610097578063d8694ccd11610071578063d8694ccd14610a5c578063f2fde38b14610a6f578063f700042a14610a82578063ffdb4b3714610a9557600080fd5b8063cc88924c14610a2e578063cdc73d5114610a41578063d02641a014610a4957600080fd5b806391a2749a116100d357806391a2749a14610939578063a69c64c01461094c578063bf78e03f1461095f578063c4276bfc14610a0c57600080fd5b80637afac3221461078e57806382b49eb0146107a15780638da5cb5b1461091157600080fd5b8063407e108611610166578063514e8cff11610140578063514e8cff146104385780636def4ce7146104db578063770e2dc41461077357806379ba50971461078657600080fd5b8063407e1086146103c557806345ac924d146103d85780634ab35b0b146103f857600080fd5b8063181f5a7711610197578063181f5a77146103525780632451a6271461039b5780633937306f146103b057600080fd5b806241e5be146101bd578063061877e3146101e357806306285c691461023c575b600080fd5b6101d06101cb366004613889565b610add565b6040519081526020015b60405180910390f35b6102236101f13660046138c5565b73ffffffffffffffffffffffffffffffffffffffff1660009081526007602052604090205467ffffffffffffffff1690565b60405167ffffffffffffffff90911681526020016101da565b610306604080516060810182526000808252602082018190529181019190915260405180606001604052807f00000000000000000000000000000000000000000000000000000000000000006bffffffffffffffffffffffff1681526020017f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1681526020017f000000000000000000000000000000000000000000000000000000000000000063ffffffff16815250905090565b6040805182516bffffffffffffffffffffffff16815260208084015173ffffffffffffffffffffffffffffffffffffffff16908201529181015163ffffffff16908201526060016101da565b61038e6040518060400160405280601781526020017f5072696365526567697374727920312e362e302d64657600000000000000000081525081565b6040516101da9190613944565b6103a3610b4b565b6040516101da9190613957565b6103c36103be3660046139b1565b610b5c565b005b6103c36103d3366004613b0d565b610e11565b6103eb6103e6366004613c6b565b610e25565b6040516101da9190613cad565b61040b6104063660046138c5565b610ef0565b6040517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff90911681526020016101da565b6104ce610446366004613d40565b60408051808201909152600080825260208201525067ffffffffffffffff166000908152600460209081526040918290208251808401909352547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff811683527c0100000000000000000000000000000000000000000000000000000000900463ffffffff169082015290565b6040516101da9190613d5b565b6107666104e9366004613d40565b6040805161022081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e08101829052610100810182905261012081018290526101408101829052610160810182905261018081018290526101a081018290526101c081018290526101e081018290526102008101919091525067ffffffffffffffff908116600090815260086020908152604091829020825161022081018452815460ff8082161515835261ffff61010080840482169685019690965263ffffffff630100000084048116978501979097526701000000000000008304871660608501526b0100000000000000000000008304871660808501526f010000000000000000000000000000008304811660a0850152710100000000000000000000000000000000008304871660c08501527501000000000000000000000000000000000000000000808404821660e08087019190915277010000000000000000000000000000000000000000000000850483169786019790975279010000000000000000000000000000000000000000000000000084049091166101208501527b010000000000000000000000000000000000000000000000000000009092048616610140840152600190930154808616610160840152640100000000810486166101808401526801000000000000000081049096166101a083015270010000000000000000000000000000000086049094166101c082015274010000000000000000000000000000000000000000850490911615156101e08201527fffffffff0000000000000000000000000000000000000000000000000000000092909304901b1661020082015290565b6040516101da9190613d96565b6103c3610781366004613fd3565b610efb565b6103c3610f11565b6103c361079c3660046142ed565b611013565b6108b16107af366004614351565b6040805160c081018252600080825260208201819052918101829052606081018290526080810182905260a08101919091525067ffffffffffffffff91909116600090815260096020908152604080832073ffffffffffffffffffffffffffffffffffffffff94909416835292815290829020825160c081018452905463ffffffff8082168352640100000000820481169383019390935268010000000000000000810461ffff16938201939093526a01000000000000000000008304821660608201526e01000000000000000000000000000083049091166080820152720100000000000000000000000000000000000090910460ff16151560a082015290565b6040516101da9190600060c08201905063ffffffff80845116835280602085015116602084015261ffff60408501511660408401528060608501511660608401528060808501511660808401525060a0830151151560a083015292915050565b60005460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101da565b6103c361094736600461437b565b611025565b6103c361095a36600461440c565b611036565b6109d861096d3660046138c5565b6040805180820182526000808252602091820181905273ffffffffffffffffffffffffffffffffffffffff93841681526006825282902082518084019093525492831682527401000000000000000000000000000000000000000090920460ff169181019190915290565b60408051825173ffffffffffffffffffffffffffffffffffffffff16815260209283015160ff1692810192909252016101da565b610a1f610a1a3660046144d1565b611047565b6040516101da9392919061456c565b6103c3610a3c366004614596565b611245565b6103a361141b565b6104ce610a573660046138c5565b611427565b6101d0610a6a366004614631565b611523565b6103c3610a7d3660046138c5565b6119dd565b6103c3610a903660046146b6565b6119ee565b610aa8610aa33660046148d6565b6119ff565b604080517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff9384168152929091166020830152016101da565b6000610ae882611b8a565b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16610b0f85611b8a565b610b37907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff168561492f565b610b419190614946565b90505b9392505050565b6060610b576002611c24565b905090565b610b64611c31565b6000610b708280614981565b9050905060005b81811015610cba576000610b8b8480614981565b83818110610b9b57610b9b6149e9565b905060400201803603810190610bb19190614a44565b604080518082018252602080840180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff908116845263ffffffff42818116858701908152885173ffffffffffffffffffffffffffffffffffffffff9081166000908152600590975295889020965190519092167c010000000000000000000000000000000000000000000000000000000002919092161790935584519051935194955016927f52f50aa6d1a95a4595361ecf953d095f125d442e4673716dede699e049de148a92610ca99290917bffffffffffffffffffffffffffffffffffffffffffffffffffffffff929092168252602082015260400190565b60405180910390a250600101610b77565b506000610cca6020840184614981565b9050905060005b81811015610e0b576000610ce86020860186614981565b83818110610cf857610cf86149e9565b905060400201803603810190610d0e9190614a81565b604080518082018252602080840180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff908116845263ffffffff42818116858701908152885167ffffffffffffffff9081166000908152600490975295889020965190519092167c010000000000000000000000000000000000000000000000000000000002919092161790935584519051935194955016927fdd84a3fa9ef9409f550d54d6affec7e9c480c878c6ab27b78912a03e1b371c6e92610dfa9290917bffffffffffffffffffffffffffffffffffffffffffffffffffffffff929092168252602082015260400190565b60405180910390a250600101610cd1565b50505050565b610e19611c76565b610e2281611cf7565b50565b60608160008167ffffffffffffffff811115610e4357610e436139ec565b604051908082528060200260200182016040528015610e8857816020015b6040805180820190915260008082526020820152815260200190600190039081610e615790505b50905060005b82811015610ee557610ec0868683818110610eab57610eab6149e9565b9050602002016020810190610a5791906138c5565b828281518110610ed257610ed26149e9565b6020908102919091010152600101610e8e565b509150505b92915050565b6000610eea82611b8a565b610f03611c76565b610f0d8282611df5565b5050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610f97576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b61101b611c76565b610f0d8282612207565b61102d611c76565b610e228161234e565b61103e611c76565b610e22816124da565b60008060607f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168773ffffffffffffffffffffffffffffffffffffffff16036110a7578592506110d5565b6110d287877f0000000000000000000000000000000000000000000000000000000000000000610add565b92505b7f00000000000000000000000000000000000000000000000000000000000000006bffffffffffffffffffffffff16831115611174576040517f6a92a483000000000000000000000000000000000000000000000000000000008152600481018490526bffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000166024820152604401610f8e565b67ffffffffffffffff8816600090815260086020526040812060010154640100000000900463ffffffff16906111ab8787846125c4565b9050806020015193508484611232836040805182516024820152602092830151151560448083019190915282518083039091018152606490910190915290810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f181dcf100000000000000000000000000000000000000000000000000000000017905290565b9450945094505050955095509592505050565b67ffffffffffffffff85166000908152600860205260408120600101547501000000000000000000000000000000000000000000900460e01b905b8481101561141257600084848381811061129c5761129c6149e9565b6112b292602060409092020190810191506138c5565b905060008787848181106112c8576112c86149e9565b90506020028101906112da9190614aa4565b6112e8906040810190614ae2565b91505060208111156113985767ffffffffffffffff8916600090815260096020908152604080832073ffffffffffffffffffffffffffffffffffffffff861684529091529020546e010000000000000000000000000000900463ffffffff16811115611398576040517f36f536ca00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff83166004820152602401610f8e565b611408848989868181106113ae576113ae6149e9565b90506020028101906113c09190614aa4565b6113ce906020810190614ae2565b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061276d92505050565b5050600101611280565b50505050505050565b6060610b57600a611c24565b604080518082019091526000808252602082015273ffffffffffffffffffffffffffffffffffffffff8281166000908152600660209081526040918290208251808401909352549283168083527401000000000000000000000000000000000000000090930460ff16908201529061151a57505073ffffffffffffffffffffffffffffffffffffffff166000908152600560209081526040918290208251808401909352547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff811683527c0100000000000000000000000000000000000000000000000000000000900463ffffffff169082015290565b610b44816127bf565b67ffffffffffffffff8083166000908152600860209081526040808320815161022081018352815460ff808216151580845261ffff61010080850482169886019890985263ffffffff630100000085048116978601979097526701000000000000008404871660608601526b0100000000000000000000008404871660808601526f010000000000000000000000000000008404811660a0860152710100000000000000000000000000000000008404871660c08601527501000000000000000000000000000000000000000000808504821660e08088019190915277010000000000000000000000000000000000000000000000860483169987019990995279010000000000000000000000000000000000000000000000000085049091166101208601527b010000000000000000000000000000000000000000000000000000009093048616610140850152600190940154808616610160850152640100000000810486166101808501526801000000000000000081049098166101a084015270010000000000000000000000000000000088049094166101c083015274010000000000000000000000000000000000000000870490931615156101e08201527fffffffff000000000000000000000000000000000000000000000000000000009290950490921b16610200840152909190611759576040517f99ac52f200000000000000000000000000000000000000000000000000000000815267ffffffffffffffff85166004820152602401610f8e565b60006117686040850185614981565b91506117c490508261177d6020870187614ae2565b90508361178a8880614ae2565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250612a0292505050565b60006007816117d960808801606089016138c5565b73ffffffffffffffffffffffffffffffffffffffff168152602081019190915260400160009081205467ffffffffffffffff169150806118286118226080890160608a016138c5565b896119ff565b909250905060008080861561186e57611862888c61184c60808e0160608f016138c5565b888e806040019061185d9190614981565b612aac565b9194509250905061188e565b6101c088015161188b9063ffffffff16662386f26fc1000061492f565b92505b61010088015160009061ffff16156118d2576118cf896dffffffffffffffffffffffffffff607088901c166118c660208f018f614ae2565b90508b86612d8a565b90505b6101a089015160009067ffffffffffffffff166118fb6118f560808f018f614ae2565b8d612e3a565b600001518563ffffffff168c60a0015161ffff168f806020019061191f9190614ae2565b61192a92915061492f565b8d6080015163ffffffff1661193f9190614b47565b6119499190614b47565b6119539190614b47565b61196d906dffffffffffffffffffffffffffff891661492f565b611977919061492f565b90507bffffffffffffffffffffffffffffffffffffffffffffffffffffffff871682826119ae67ffffffffffffffff8c168961492f565b6119b89190614b47565b6119c29190614b47565b6119cc9190614946565b9d9c50505050505050505050505050565b6119e5611c76565b610e2281612efb565b6119f6611c76565b610e2281612ff0565b67ffffffffffffffff811660009081526004602090815260408083208151808301909252547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff811682527c0100000000000000000000000000000000000000000000000000000000900463ffffffff1691810182905282918203611ab7576040517f2e59db3a00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff85166004820152602401610f8e565b6000816020015163ffffffff1642611acf9190614b5a565b90507f000000000000000000000000000000000000000000000000000000000000000063ffffffff16811115611b70576040517ff08bcb3e00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8616600482015263ffffffff7f000000000000000000000000000000000000000000000000000000000000000016602482015260448101829052606401610f8e565b611b7986611b8a565b9151919350909150505b9250929050565b600080611b9683611427565b9050806020015163ffffffff1660001480611bce575080517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16155b15611c1d576040517f06439c6b00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff84166004820152602401610f8e565b5192915050565b60606000610b44836134de565b611c3c60023361353a565b611c74576040517fd86ad9cf000000000000000000000000000000000000000000000000000000008152336004820152602401610f8e565b565b60005473ffffffffffffffffffffffffffffffffffffffff163314611c74576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610f8e565b60005b8151811015610f0d576000828281518110611d1757611d176149e9565b60209081029190910181015180518183015173ffffffffffffffffffffffffffffffffffffffff80831660008181526006875260409081902084518154868a018051929096167fffffffffffffffffffffff00000000000000000000000000000000000000000090911681177401000000000000000000000000000000000000000060ff9384160217909255825191825293519093169683019690965293955091939092917f08a5f7f5bb38a81d8e43aca13ecd76431dbf8816ae4699affff7b00b2fc1c464910160405180910390a2505050806001019050611cfa565b60005b825181101561211e576000838281518110611e1557611e156149e9565b6020026020010151905060008160000151905060005b82602001515181101561211057600083602001518281518110611e5057611e506149e9565b6020026020010151602001519050600084602001518381518110611e7657611e766149e9565b6020026020010151600001519050602063ffffffff16826080015163ffffffff161015611ef95760808201516040517f24ecdc0200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8316600482015263ffffffff9091166024820152604401610f8e565b67ffffffffffffffff8416600081815260096020908152604080832073ffffffffffffffffffffffffffffffffffffffff86168085529083529281902086518154938801518389015160608a015160808b015160a08c015115157201000000000000000000000000000000000000027fffffffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffff63ffffffff9283166e01000000000000000000000000000002167fffffffffffffffffffffffffff0000000000ffffffffffffffffffffffffffff9383166a0100000000000000000000027fffffffffffffffffffffffffffffffffffff00000000ffffffffffffffffffff61ffff9096166801000000000000000002959095167fffffffffffffffffffffffffffffffffffff000000000000ffffffffffffffff968416640100000000027fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000909b16939097169290921798909817939093169390931717919091161792909217909155519091907f94967ae9ea7729ad4f54021c1981765d2b1d954f7c92fbec340aa0a54f46b8b5906120fe908690600060c08201905063ffffffff80845116835280602085015116602084015261ffff60408501511660408401528060608501511660608401528060808501511660808401525060a0830151151560a083015292915050565b60405180910390a35050600101611e2b565b505050806001019050611df8565b5060005b815181101561220257600082828151811061213f5761213f6149e9565b60200260200101516000015190506000838381518110612161576121616149e9565b60209081029190910181015181015167ffffffffffffffff8416600081815260098452604080822073ffffffffffffffffffffffffffffffffffffffff8516808452955280822080547fffffffffffffffffffffffffff000000000000000000000000000000000000001690555192945090917f4de5b1bcbca6018c11303a2c3f4a4b4f22a1c741d8c4ba430d246ac06c5ddf8b9190a35050600101612122565b505050565b60005b82518110156122aa57612240838281518110612228576122286149e9565b6020026020010151600a61356990919063ffffffff16565b156122a257828181518110612257576122576149e9565b602002602001015173ffffffffffffffffffffffffffffffffffffffff167fdf1b1bd32a69711488d71554706bb130b1fc63a5fa1a2cd85e8440f84065ba2360405160405180910390a25b60010161220a565b5060005b8151811015612202576122e48282815181106122cc576122cc6149e9565b6020026020010151600a61358b90919063ffffffff16565b15612346578181815181106122fb576122fb6149e9565b602002602001015173ffffffffffffffffffffffffffffffffffffffff167f1795838dc8ab2ffc5f431a1729a6afa0b587f982f7b2be0b9d7187a1ef547f9160405160405180910390a25b6001016122ae565b602081015160005b81518110156123e9576000828281518110612373576123736149e9565b602002602001015190506123918160026135ad90919063ffffffff16565b156123e05760405173ffffffffffffffffffffffffffffffffffffffff821681527fc3803387881faad271c47728894e3e36fac830ffc8602ca6fc07733cbda775809060200160405180910390a15b50600101612356565b50815160005b8151811015610e0b57600082828151811061240c5761240c6149e9565b60200260200101519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff160361247c576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b612487600282613569565b5060405173ffffffffffffffffffffffffffffffffffffffff821681527feb1b9b92e50b7f88f9ff25d56765095ac6e91540eee214906f4036a908ffbdef9060200160405180910390a1506001016123ef565b60005b8151811015610f0d5760008282815181106124fa576124fa6149e9565b6020026020010151600001519050600083838151811061251c5761251c6149e9565b60209081029190910181015181015173ffffffffffffffffffffffffffffffffffffffff841660008181526007845260409081902080547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001667ffffffffffffffff85169081179091559051908152919350917fbb77da6f7210cdd16904228a9360133d1d7dfff99b1bc75f128da5b53e28f97d910160405180910390a250506001016124dd565b6040805180820190915260008082526020820152600083900361260557506040805180820190915267ffffffffffffffff8216815260006020820152610b44565b60006126118486614b6d565b905060006126228560048189614bb3565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509293505050507fffffffff0000000000000000000000000000000000000000000000000000000082167fe7e230f000000000000000000000000000000000000000000000000000000000016126bf57808060200190518101906126b69190614bdd565b92505050610b44565b7f6859a837000000000000000000000000000000000000000000000000000000007fffffffff0000000000000000000000000000000000000000000000000000000083160161273b576040518060400160405280828060200190518101906127279190614c09565b815260006020909101529250610b44915050565b6040517f5247fdce00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7fd7ed2ad4000000000000000000000000000000000000000000000000000000007fffffffff00000000000000000000000000000000000000000000000000000000831601610f0d57612202816135cf565b604080518082019091526000808252602082015260008260000151905060008173ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa158015612829573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061284d9190614c3c565b505050915050600081121561288e576040517f10cb51d100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000819050600085602001518473ffffffffffffffffffffffffffffffffffffffff1663313ce5676040518163ffffffff1660e01b8152600401602060405180830381865afa1580156128e5573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906129099190614c8c565b6129139190614ca9565b905060248160ff1611156129485761292c602482614cc2565b61293790600a614dfb565b6129419083614946565b915061296b565b612953816024614cc2565b61295e90600a614dfb565b612968908361492f565b91505b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8211156129c1576040517f10cb51d100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b50604080518082019091527bffffffffffffffffffffffffffffffffffffffffffffffffffffffff909116815263ffffffff42166020820152949350505050565b836040015163ffffffff16831115612a5b5760408085015190517f8693378900000000000000000000000000000000000000000000000000000000815263ffffffff909116600482015260248101849052604401610f8e565b836020015161ffff16821115612a9d576040517f4c056b6a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610e0b8461020001518261276d565b6000808083815b81811015612d7c576000878783818110612acf57612acf6149e9565b905060400201803603810190612ae59190614e0a565b67ffffffffffffffff8c166000908152600960209081526040808320845173ffffffffffffffffffffffffffffffffffffffff168452825291829020825160c081018452905463ffffffff8082168352640100000000820481169383019390935268010000000000000000810461ffff16938201939093526a01000000000000000000008304821660608201526e01000000000000000000000000000083049091166080820152720100000000000000000000000000000000000090910460ff16151560a0820181905291925090612c0b576101208d0151612bd29061ffff16662386f26fc1000061492f565b612bdc9088614b47565b96508c610140015186612bef9190614e43565b95508c610160015185612c029190614e43565b94505050612d74565b604081015160009061ffff1615612cc45760008c73ffffffffffffffffffffffffffffffffffffffff16846000015173ffffffffffffffffffffffffffffffffffffffff1614612c67578351612c6090611b8a565b9050612c6a565b508a5b620186a0836040015161ffff16612cac8660200151847bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1661368290919063ffffffff16565b612cb6919061492f565b612cc09190614946565b9150505b6060820151612cd39088614e43565b9650816080015186612ce59190614e43565b8251909650600090612d049063ffffffff16662386f26fc1000061492f565b905080821015612d2357612d18818a614b47565b985050505050612d74565b6000836020015163ffffffff16662386f26fc10000612d42919061492f565b905080831115612d6257612d56818b614b47565b99505050505050612d74565b612d6c838b614b47565b995050505050505b600101612ab3565b505096509650969350505050565b60008063ffffffff8316612da06101408661492f565b612dac876101c0614b47565b612db69190614b47565b612dc09190614b47565b905060008760c0015163ffffffff168860e0015161ffff1683612de3919061492f565b612ded9190614b47565b61010089015190915061ffff16612e146dffffffffffffffffffffffffffff89168361492f565b612e1e919061492f565b612e2e90655af3107a400061492f565b98975050505050505050565b60408051808201909152600080825260208201526000612e66858585610180015163ffffffff166125c4565b9050826060015163ffffffff1681600001511115612eb0576040517f4c4fc93a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b826101e001518015612ec457508060200151155b15610b41576040517fee433e9900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff821603612f7a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610f8e565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60005b8151811015610f0d576000828281518110613010576130106149e9565b60200260200101519050600083838151811061302e5761302e6149e9565b60200260200101516000015190506000826020015190508167ffffffffffffffff1660001480613067575061018081015163ffffffff16155b806130b957506102008101517fffffffff00000000000000000000000000000000000000000000000000000000167f2812d52c0000000000000000000000000000000000000000000000000000000014155b806130d55750602063ffffffff1681610160015163ffffffff16105b806130f45750806060015163ffffffff1681610180015163ffffffff16115b15613137576040517fc35aa79d00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff83166004820152602401610f8e565b67ffffffffffffffff82166000908152600860205260408120600101547501000000000000000000000000000000000000000000900460e01b7fffffffff000000000000000000000000000000000000000000000000000000001690036131df578167ffffffffffffffff167fa937382a486d993de71c220bc8b559242deb4e286a353fa732330b4aa7d13577826040516131d29190613d96565b60405180910390a2613222565b8167ffffffffffffffff167fa7b607fc10d28a1caf39ab7d27f4c94945db708a576d572781a455c5894fad93826040516132199190613d96565b60405180910390a25b80600860008467ffffffffffffffff1667ffffffffffffffff16815260200190815260200160002060008201518160000160006101000a81548160ff02191690831515021790555060208201518160000160016101000a81548161ffff021916908361ffff16021790555060408201518160000160036101000a81548163ffffffff021916908363ffffffff16021790555060608201518160000160076101000a81548163ffffffff021916908363ffffffff160217905550608082015181600001600b6101000a81548163ffffffff021916908363ffffffff16021790555060a082015181600001600f6101000a81548161ffff021916908361ffff16021790555060c08201518160000160116101000a81548163ffffffff021916908363ffffffff16021790555060e08201518160000160156101000a81548161ffff021916908361ffff1602179055506101008201518160000160176101000a81548161ffff021916908361ffff1602179055506101208201518160000160196101000a81548161ffff021916908361ffff16021790555061014082015181600001601b6101000a81548163ffffffff021916908363ffffffff1602179055506101608201518160010160006101000a81548163ffffffff021916908363ffffffff1602179055506101808201518160010160046101000a81548163ffffffff021916908363ffffffff1602179055506101a08201518160010160086101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055506101c08201518160010160106101000a81548163ffffffff021916908363ffffffff1602179055506101e08201518160010160146101000a81548160ff0219169083151502179055506102008201518160010160156101000a81548163ffffffff021916908360e01c0217905550905050505050806001019050612ff3565b60608160000180548060200260200160405190810160405280929190818152602001828054801561352e57602002820191906000526020600020905b81548152602001906001019080831161351a575b50505050509050919050565b73ffffffffffffffffffffffffffffffffffffffff811660009081526001830160205260408120541515610b44565b6000610b448373ffffffffffffffffffffffffffffffffffffffff84166136bf565b6000610b448373ffffffffffffffffffffffffffffffffffffffff841661370e565b6000610b448373ffffffffffffffffffffffffffffffffffffffff8416613808565b6000815160201461360e57816040517f8d666f60000000000000000000000000000000000000000000000000000000008152600401610f8e9190613944565b6000828060200190518101906136249190614c09565b905073ffffffffffffffffffffffffffffffffffffffff811180613649575061040081105b15610eea57826040517f8d666f60000000000000000000000000000000000000000000000000000000008152600401610f8e9190613944565b6000670de0b6b3a76400006136b5837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff861661492f565b610b449190614946565b600081815260018301602052604081205461370657508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610eea565b506000610eea565b600081815260018301602052604081205480156137f7576000613732600183614b5a565b855490915060009061374690600190614b5a565b90508082146137ab576000866000018281548110613766576137666149e9565b9060005260206000200154905080876000018481548110613789576137896149e9565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806137bc576137bc614e60565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050610eea565b6000915050610eea565b5092915050565b600081815260018301602052604081205480156137f757600061382c600183614b5a565b855490915060009061384090600190614b5a565b90508181146137ab576000866000018281548110613766576137666149e9565b803573ffffffffffffffffffffffffffffffffffffffff8116811461388457600080fd5b919050565b60008060006060848603121561389e57600080fd5b6138a784613860565b9250602084013591506138bc60408501613860565b90509250925092565b6000602082840312156138d757600080fd5b610b4482613860565b6000815180845260005b81811015613906576020818501810151868301820152016138ea565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b602081526000610b4460208301846138e0565b6020808252825182820181905260009190848201906040850190845b818110156139a557835173ffffffffffffffffffffffffffffffffffffffff1683529284019291840191600101613973565b50909695505050505050565b6000602082840312156139c357600080fd5b813567ffffffffffffffff8111156139da57600080fd5b820160408185031215610b4457600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040805190810167ffffffffffffffff81118282101715613a3e57613a3e6139ec565b60405290565b60405160c0810167ffffffffffffffff81118282101715613a3e57613a3e6139ec565b604051610220810167ffffffffffffffff81118282101715613a3e57613a3e6139ec565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715613ad257613ad26139ec565b604052919050565b600067ffffffffffffffff821115613af457613af46139ec565b5060051b60200190565b60ff81168114610e2257600080fd5b60006020808385031215613b2057600080fd5b823567ffffffffffffffff811115613b3757600080fd5b8301601f81018513613b4857600080fd5b8035613b5b613b5682613ada565b613a8b565b81815260609182028301840191848201919088841115613b7a57600080fd5b938501935b83851015613c1a5784890381811215613b985760008081fd5b613ba0613a1b565b613ba987613860565b81526040807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe084011215613bdd5760008081fd5b613be5613a1b565b9250613bf2898901613860565b8352870135613c0081613afe565b828901528088019190915283529384019391850191613b7f565b50979650505050505050565b60008083601f840112613c3857600080fd5b50813567ffffffffffffffff811115613c5057600080fd5b6020830191508360208260051b8501011115611b8357600080fd5b60008060208385031215613c7e57600080fd5b823567ffffffffffffffff811115613c9557600080fd5b613ca185828601613c26565b90969095509350505050565b602080825282518282018190526000919060409081850190868401855b82811015613d1b57613d0b84835180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16825260209081015163ffffffff16910152565b9284019290850190600101613cca565b5091979650505050505050565b803567ffffffffffffffff8116811461388457600080fd5b600060208284031215613d5257600080fd5b610b4482613d28565b81517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815260208083015163ffffffff169082015260408101610eea565b81511515815261022081016020830151613db6602084018261ffff169052565b506040830151613dce604084018263ffffffff169052565b506060830151613de6606084018263ffffffff169052565b506080830151613dfe608084018263ffffffff169052565b5060a0830151613e1460a084018261ffff169052565b5060c0830151613e2c60c084018263ffffffff169052565b5060e0830151613e4260e084018261ffff169052565b506101008381015161ffff9081169184019190915261012080850151909116908301526101408084015163ffffffff9081169184019190915261016080850151821690840152610180808501518216908401526101a08085015167ffffffffffffffff16908401526101c080850151909116908301526101e080840151151590830152610200808401517fffffffff000000000000000000000000000000000000000000000000000000008116828501525b505092915050565b803563ffffffff8116811461388457600080fd5b803561ffff8116811461388457600080fd5b8015158114610e2257600080fd5b803561388481613f22565b600082601f830112613f4c57600080fd5b81356020613f5c613b5683613ada565b82815260069290921b84018101918181019086841115613f7b57600080fd5b8286015b84811015613fc85760408189031215613f985760008081fd5b613fa0613a1b565b613fa982613d28565b8152613fb6858301613860565b81860152835291830191604001613f7f565b509695505050505050565b60008060408385031215613fe657600080fd5b67ffffffffffffffff83351115613ffc57600080fd5b83601f84358501011261400e57600080fd5b61401e613b568435850135613ada565b8335840180358083526020808401939260059290921b9091010186101561404457600080fd5b602085358601015b85358601803560051b016020018110156142515767ffffffffffffffff8135111561407657600080fd5b8035863587010160407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0828a030112156140af57600080fd5b6140b7613a1b565b6140c360208301613d28565b815267ffffffffffffffff604083013511156140de57600080fd5b88603f6040840135840101126140f357600080fd5b614109613b566020604085013585010135613ada565b6020604084810135850182810135808552928401939260e00201018b101561413057600080fd5b6040808501358501015b6040858101358601602081013560e00201018110156142325760e0818d03121561416357600080fd5b61416b613a1b565b61417482613860565b815260c07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0838f030112156141a857600080fd5b6141b0613a44565b6141bc60208401613efc565b81526141ca60408401613efc565b60208201526141db60608401613f10565b60408201526141ec60808401613efc565b60608201526141fd60a08401613efc565b608082015261420f60c0840135613f22565b60c083013560a0820152602082810191909152908452929092019160e00161413a565b508060208401525050808552505060208301925060208101905061404c565b5092505067ffffffffffffffff6020840135111561426e57600080fd5b61427e8460208501358501613f3b565b90509250929050565b600082601f83011261429857600080fd5b813560206142a8613b5683613ada565b8083825260208201915060208460051b8701019350868411156142ca57600080fd5b602086015b84811015613fc8576142e081613860565b83529183019183016142cf565b6000806040838503121561430057600080fd5b823567ffffffffffffffff8082111561431857600080fd5b61432486838701614287565b9350602085013591508082111561433a57600080fd5b5061434785828601614287565b9150509250929050565b6000806040838503121561436457600080fd5b61436d83613d28565b915061427e60208401613860565b60006020828403121561438d57600080fd5b813567ffffffffffffffff808211156143a557600080fd5b90830190604082860312156143b957600080fd5b6143c1613a1b565b8235828111156143d057600080fd5b6143dc87828601614287565b8252506020830135828111156143f157600080fd5b6143fd87828601614287565b60208301525095945050505050565b6000602080838503121561441f57600080fd5b823567ffffffffffffffff81111561443657600080fd5b8301601f8101851361444757600080fd5b8035614455613b5682613ada565b81815260069190911b8201830190838101908783111561447457600080fd5b928401925b828410156144c657604084890312156144925760008081fd5b61449a613a1b565b6144a385613860565b81526144b0868601613d28565b8187015282526040939093019290840190614479565b979650505050505050565b6000806000806000608086880312156144e957600080fd5b6144f286613d28565b945061450060208701613860565b935060408601359250606086013567ffffffffffffffff8082111561452457600080fd5b818801915088601f83011261453857600080fd5b81358181111561454757600080fd5b89602082850101111561455957600080fd5b9699959850939650602001949392505050565b838152821515602082015260606040820152600061458d60608301846138e0565b95945050505050565b6000806000806000606086880312156145ae57600080fd5b6145b786613d28565b9450602086013567ffffffffffffffff808211156145d457600080fd5b6145e089838a01613c26565b909650945060408801359150808211156145f957600080fd5b818801915088601f83011261460d57600080fd5b81358181111561461c57600080fd5b8960208260061b850101111561455957600080fd5b6000806040838503121561464457600080fd5b61464d83613d28565b9150602083013567ffffffffffffffff81111561466957600080fd5b830160a0818603121561467b57600080fd5b809150509250929050565b80357fffffffff000000000000000000000000000000000000000000000000000000008116811461388457600080fd5b600060208083850312156146c957600080fd5b823567ffffffffffffffff8111156146e057600080fd5b8301601f810185136146f157600080fd5b80356146ff613b5682613ada565b818152610240918202830184019184820191908884111561471f57600080fd5b938501935b83851015613c1a578489038181121561473d5760008081fd5b614745613a1b565b61474e87613d28565b8152610220807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0840112156147835760008081fd5b61478b613a67565b9250614798898901613f30565b835260406147a7818a01613f10565b8a85015260606147b8818b01613efc565b82860152608091506147cb828b01613efc565b9085015260a06147dc8a8201613efc565b8286015260c091506147ef828b01613f10565b9085015260e06148008a8201613efc565b828601526101009150614814828b01613f10565b908501526101206148268a8201613f10565b82860152610140915061483a828b01613f10565b9085015261016061484c8a8201613efc565b828601526101809150614860828b01613efc565b908501526101a06148728a8201613efc565b828601526101c09150614886828b01613d28565b908501526101e06148988a8201613efc565b8286015261020091506148ac828b01613f30565b908501526148bb898301614686565b90840152508088019190915283529384019391850191614724565b600080604083850312156148e957600080fd5b6148f283613860565b915061427e60208401613d28565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8082028115828204841417610eea57610eea614900565b60008261497c577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18436030181126149b657600080fd5b83018035915067ffffffffffffffff8211156149d157600080fd5b6020019150600681901b3603821315611b8357600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b80357bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8116811461388457600080fd5b600060408284031215614a5657600080fd5b614a5e613a1b565b614a6783613860565b8152614a7560208401614a18565b60208201529392505050565b600060408284031215614a9357600080fd5b614a9b613a1b565b614a6783613d28565b600082357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81833603018112614ad857600080fd5b9190910192915050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112614b1757600080fd5b83018035915067ffffffffffffffff821115614b3257600080fd5b602001915036819003821315611b8357600080fd5b80820180821115610eea57610eea614900565b81810381811115610eea57610eea614900565b7fffffffff000000000000000000000000000000000000000000000000000000008135818116916004851015613ef45760049490940360031b84901b1690921692915050565b60008085851115614bc357600080fd5b83861115614bd057600080fd5b5050820193919092039150565b600060408284031215614bef57600080fd5b614bf7613a1b565b825181526020830151614a7581613f22565b600060208284031215614c1b57600080fd5b5051919050565b805169ffffffffffffffffffff8116811461388457600080fd5b600080600080600060a08688031215614c5457600080fd5b614c5d86614c22565b9450602086015193506040860151925060608601519150614c8060808701614c22565b90509295509295909350565b600060208284031215614c9e57600080fd5b8151610b4481613afe565b60ff8181168382160190811115610eea57610eea614900565b60ff8281168282160390811115610eea57610eea614900565b600181815b80851115614d3457817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04821115614d1a57614d1a614900565b80851615614d2757918102915b93841c9390800290614ce0565b509250929050565b600082614d4b57506001610eea565b81614d5857506000610eea565b8160018114614d6e5760028114614d7857614d94565b6001915050610eea565b60ff841115614d8957614d89614900565b50506001821b610eea565b5060208310610133831016604e8410600b8410161715614db7575081810a610eea565b614dc18383614cdb565b807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04821115614df357614df3614900565b029392505050565b6000610b4460ff841683614d3c565b600060408284031215614e1c57600080fd5b614e24613a1b565b614e2d83613860565b8152602083013560208201528091505092915050565b63ffffffff81811683821601908082111561380157613801614900565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea164736f6c6343000818000a",
}

var PriceRegistryABI = PriceRegistryMetaData.ABI

var PriceRegistryBin = PriceRegistryMetaData.Bin

func DeployPriceRegistry(auth *bind.TransactOpts, backend bind.ContractBackend, staticConfig PriceRegistryStaticConfig, priceUpdaters []common.Address, feeTokens []common.Address, tokenPriceFeeds []PriceRegistryTokenPriceFeedUpdate, tokenTransferFeeConfigArgs []PriceRegistryTokenTransferFeeConfigArgs, premiumMultiplierWeiPerEthArgs []PriceRegistryPremiumMultiplierWeiPerEthArgs, destChainConfigArgs []PriceRegistryDestChainConfigArgs) (common.Address, *types.Transaction, *PriceRegistry, error) {
	parsed, err := PriceRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PriceRegistryBin), backend, staticConfig, priceUpdaters, feeTokens, tokenPriceFeeds, tokenTransferFeeConfigArgs, premiumMultiplierWeiPerEthArgs, destChainConfigArgs)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &PriceRegistry{address: address, abi: *parsed, PriceRegistryCaller: PriceRegistryCaller{contract: contract}, PriceRegistryTransactor: PriceRegistryTransactor{contract: contract}, PriceRegistryFilterer: PriceRegistryFilterer{contract: contract}}, nil
}

type PriceRegistry struct {
	address common.Address
	abi     abi.ABI
	PriceRegistryCaller
	PriceRegistryTransactor
	PriceRegistryFilterer
}

type PriceRegistryCaller struct {
	contract *bind.BoundContract
}

type PriceRegistryTransactor struct {
	contract *bind.BoundContract
}

type PriceRegistryFilterer struct {
	contract *bind.BoundContract
}

type PriceRegistrySession struct {
	Contract     *PriceRegistry
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type PriceRegistryCallerSession struct {
	Contract *PriceRegistryCaller
	CallOpts bind.CallOpts
}

type PriceRegistryTransactorSession struct {
	Contract     *PriceRegistryTransactor
	TransactOpts bind.TransactOpts
}

type PriceRegistryRaw struct {
	Contract *PriceRegistry
}

type PriceRegistryCallerRaw struct {
	Contract *PriceRegistryCaller
}

type PriceRegistryTransactorRaw struct {
	Contract *PriceRegistryTransactor
}

func NewPriceRegistry(address common.Address, backend bind.ContractBackend) (*PriceRegistry, error) {
	abi, err := abi.JSON(strings.NewReader(PriceRegistryABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindPriceRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PriceRegistry{address: address, abi: abi, PriceRegistryCaller: PriceRegistryCaller{contract: contract}, PriceRegistryTransactor: PriceRegistryTransactor{contract: contract}, PriceRegistryFilterer: PriceRegistryFilterer{contract: contract}}, nil
}

func NewPriceRegistryCaller(address common.Address, caller bind.ContractCaller) (*PriceRegistryCaller, error) {
	contract, err := bindPriceRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PriceRegistryCaller{contract: contract}, nil
}

func NewPriceRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*PriceRegistryTransactor, error) {
	contract, err := bindPriceRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PriceRegistryTransactor{contract: contract}, nil
}

func NewPriceRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*PriceRegistryFilterer, error) {
	contract, err := bindPriceRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PriceRegistryFilterer{contract: contract}, nil
}

func bindPriceRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PriceRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_PriceRegistry *PriceRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PriceRegistry.Contract.PriceRegistryCaller.contract.Call(opts, result, method, params...)
}

func (_PriceRegistry *PriceRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PriceRegistry.Contract.PriceRegistryTransactor.contract.Transfer(opts)
}

func (_PriceRegistry *PriceRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PriceRegistry.Contract.PriceRegistryTransactor.contract.Transact(opts, method, params...)
}

func (_PriceRegistry *PriceRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PriceRegistry.Contract.contract.Call(opts, result, method, params...)
}

func (_PriceRegistry *PriceRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PriceRegistry.Contract.contract.Transfer(opts)
}

func (_PriceRegistry *PriceRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PriceRegistry.Contract.contract.Transact(opts, method, params...)
}

func (_PriceRegistry *PriceRegistryCaller) ConvertTokenAmount(opts *bind.CallOpts, fromToken common.Address, fromTokenAmount *big.Int, toToken common.Address) (*big.Int, error) {
	var out []interface{}
	err := _PriceRegistry.contract.Call(opts, &out, "convertTokenAmount", fromToken, fromTokenAmount, toToken)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_PriceRegistry *PriceRegistrySession) ConvertTokenAmount(fromToken common.Address, fromTokenAmount *big.Int, toToken common.Address) (*big.Int, error) {
	return _PriceRegistry.Contract.ConvertTokenAmount(&_PriceRegistry.CallOpts, fromToken, fromTokenAmount, toToken)
}

func (_PriceRegistry *PriceRegistryCallerSession) ConvertTokenAmount(fromToken common.Address, fromTokenAmount *big.Int, toToken common.Address) (*big.Int, error) {
	return _PriceRegistry.Contract.ConvertTokenAmount(&_PriceRegistry.CallOpts, fromToken, fromTokenAmount, toToken)
}

func (_PriceRegistry *PriceRegistryCaller) GetAllAuthorizedCallers(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _PriceRegistry.contract.Call(opts, &out, "getAllAuthorizedCallers")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_PriceRegistry *PriceRegistrySession) GetAllAuthorizedCallers() ([]common.Address, error) {
	return _PriceRegistry.Contract.GetAllAuthorizedCallers(&_PriceRegistry.CallOpts)
}

func (_PriceRegistry *PriceRegistryCallerSession) GetAllAuthorizedCallers() ([]common.Address, error) {
	return _PriceRegistry.Contract.GetAllAuthorizedCallers(&_PriceRegistry.CallOpts)
}

func (_PriceRegistry *PriceRegistryCaller) GetDestChainConfig(opts *bind.CallOpts, destChainSelector uint64) (PriceRegistryDestChainConfig, error) {
	var out []interface{}
	err := _PriceRegistry.contract.Call(opts, &out, "getDestChainConfig", destChainSelector)

	if err != nil {
		return *new(PriceRegistryDestChainConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(PriceRegistryDestChainConfig)).(*PriceRegistryDestChainConfig)

	return out0, err

}

func (_PriceRegistry *PriceRegistrySession) GetDestChainConfig(destChainSelector uint64) (PriceRegistryDestChainConfig, error) {
	return _PriceRegistry.Contract.GetDestChainConfig(&_PriceRegistry.CallOpts, destChainSelector)
}

func (_PriceRegistry *PriceRegistryCallerSession) GetDestChainConfig(destChainSelector uint64) (PriceRegistryDestChainConfig, error) {
	return _PriceRegistry.Contract.GetDestChainConfig(&_PriceRegistry.CallOpts, destChainSelector)
}

func (_PriceRegistry *PriceRegistryCaller) GetDestinationChainGasPrice(opts *bind.CallOpts, destChainSelector uint64) (InternalTimestampedPackedUint224, error) {
	var out []interface{}
	err := _PriceRegistry.contract.Call(opts, &out, "getDestinationChainGasPrice", destChainSelector)

	if err != nil {
		return *new(InternalTimestampedPackedUint224), err
	}

	out0 := *abi.ConvertType(out[0], new(InternalTimestampedPackedUint224)).(*InternalTimestampedPackedUint224)

	return out0, err

}

func (_PriceRegistry *PriceRegistrySession) GetDestinationChainGasPrice(destChainSelector uint64) (InternalTimestampedPackedUint224, error) {
	return _PriceRegistry.Contract.GetDestinationChainGasPrice(&_PriceRegistry.CallOpts, destChainSelector)
}

func (_PriceRegistry *PriceRegistryCallerSession) GetDestinationChainGasPrice(destChainSelector uint64) (InternalTimestampedPackedUint224, error) {
	return _PriceRegistry.Contract.GetDestinationChainGasPrice(&_PriceRegistry.CallOpts, destChainSelector)
}

func (_PriceRegistry *PriceRegistryCaller) GetFeeTokens(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _PriceRegistry.contract.Call(opts, &out, "getFeeTokens")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_PriceRegistry *PriceRegistrySession) GetFeeTokens() ([]common.Address, error) {
	return _PriceRegistry.Contract.GetFeeTokens(&_PriceRegistry.CallOpts)
}

func (_PriceRegistry *PriceRegistryCallerSession) GetFeeTokens() ([]common.Address, error) {
	return _PriceRegistry.Contract.GetFeeTokens(&_PriceRegistry.CallOpts)
}

func (_PriceRegistry *PriceRegistryCaller) GetPremiumMultiplierWeiPerEth(opts *bind.CallOpts, token common.Address) (uint64, error) {
	var out []interface{}
	err := _PriceRegistry.contract.Call(opts, &out, "getPremiumMultiplierWeiPerEth", token)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_PriceRegistry *PriceRegistrySession) GetPremiumMultiplierWeiPerEth(token common.Address) (uint64, error) {
	return _PriceRegistry.Contract.GetPremiumMultiplierWeiPerEth(&_PriceRegistry.CallOpts, token)
}

func (_PriceRegistry *PriceRegistryCallerSession) GetPremiumMultiplierWeiPerEth(token common.Address) (uint64, error) {
	return _PriceRegistry.Contract.GetPremiumMultiplierWeiPerEth(&_PriceRegistry.CallOpts, token)
}

func (_PriceRegistry *PriceRegistryCaller) GetStaticConfig(opts *bind.CallOpts) (PriceRegistryStaticConfig, error) {
	var out []interface{}
	err := _PriceRegistry.contract.Call(opts, &out, "getStaticConfig")

	if err != nil {
		return *new(PriceRegistryStaticConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(PriceRegistryStaticConfig)).(*PriceRegistryStaticConfig)

	return out0, err

}

func (_PriceRegistry *PriceRegistrySession) GetStaticConfig() (PriceRegistryStaticConfig, error) {
	return _PriceRegistry.Contract.GetStaticConfig(&_PriceRegistry.CallOpts)
}

func (_PriceRegistry *PriceRegistryCallerSession) GetStaticConfig() (PriceRegistryStaticConfig, error) {
	return _PriceRegistry.Contract.GetStaticConfig(&_PriceRegistry.CallOpts)
}

func (_PriceRegistry *PriceRegistryCaller) GetTokenAndGasPrices(opts *bind.CallOpts, token common.Address, destChainSelector uint64) (GetTokenAndGasPrices,

	error) {
	var out []interface{}
	err := _PriceRegistry.contract.Call(opts, &out, "getTokenAndGasPrices", token, destChainSelector)

	outstruct := new(GetTokenAndGasPrices)
	if err != nil {
		return *outstruct, err
	}

	outstruct.TokenPrice = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.GasPriceValue = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_PriceRegistry *PriceRegistrySession) GetTokenAndGasPrices(token common.Address, destChainSelector uint64) (GetTokenAndGasPrices,

	error) {
	return _PriceRegistry.Contract.GetTokenAndGasPrices(&_PriceRegistry.CallOpts, token, destChainSelector)
}

func (_PriceRegistry *PriceRegistryCallerSession) GetTokenAndGasPrices(token common.Address, destChainSelector uint64) (GetTokenAndGasPrices,

	error) {
	return _PriceRegistry.Contract.GetTokenAndGasPrices(&_PriceRegistry.CallOpts, token, destChainSelector)
}

func (_PriceRegistry *PriceRegistryCaller) GetTokenPrice(opts *bind.CallOpts, token common.Address) (InternalTimestampedPackedUint224, error) {
	var out []interface{}
	err := _PriceRegistry.contract.Call(opts, &out, "getTokenPrice", token)

	if err != nil {
		return *new(InternalTimestampedPackedUint224), err
	}

	out0 := *abi.ConvertType(out[0], new(InternalTimestampedPackedUint224)).(*InternalTimestampedPackedUint224)

	return out0, err

}

func (_PriceRegistry *PriceRegistrySession) GetTokenPrice(token common.Address) (InternalTimestampedPackedUint224, error) {
	return _PriceRegistry.Contract.GetTokenPrice(&_PriceRegistry.CallOpts, token)
}

func (_PriceRegistry *PriceRegistryCallerSession) GetTokenPrice(token common.Address) (InternalTimestampedPackedUint224, error) {
	return _PriceRegistry.Contract.GetTokenPrice(&_PriceRegistry.CallOpts, token)
}

func (_PriceRegistry *PriceRegistryCaller) GetTokenPriceFeedConfig(opts *bind.CallOpts, token common.Address) (IPriceRegistryTokenPriceFeedConfig, error) {
	var out []interface{}
	err := _PriceRegistry.contract.Call(opts, &out, "getTokenPriceFeedConfig", token)

	if err != nil {
		return *new(IPriceRegistryTokenPriceFeedConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(IPriceRegistryTokenPriceFeedConfig)).(*IPriceRegistryTokenPriceFeedConfig)

	return out0, err

}

func (_PriceRegistry *PriceRegistrySession) GetTokenPriceFeedConfig(token common.Address) (IPriceRegistryTokenPriceFeedConfig, error) {
	return _PriceRegistry.Contract.GetTokenPriceFeedConfig(&_PriceRegistry.CallOpts, token)
}

func (_PriceRegistry *PriceRegistryCallerSession) GetTokenPriceFeedConfig(token common.Address) (IPriceRegistryTokenPriceFeedConfig, error) {
	return _PriceRegistry.Contract.GetTokenPriceFeedConfig(&_PriceRegistry.CallOpts, token)
}

func (_PriceRegistry *PriceRegistryCaller) GetTokenPrices(opts *bind.CallOpts, tokens []common.Address) ([]InternalTimestampedPackedUint224, error) {
	var out []interface{}
	err := _PriceRegistry.contract.Call(opts, &out, "getTokenPrices", tokens)

	if err != nil {
		return *new([]InternalTimestampedPackedUint224), err
	}

	out0 := *abi.ConvertType(out[0], new([]InternalTimestampedPackedUint224)).(*[]InternalTimestampedPackedUint224)

	return out0, err

}

func (_PriceRegistry *PriceRegistrySession) GetTokenPrices(tokens []common.Address) ([]InternalTimestampedPackedUint224, error) {
	return _PriceRegistry.Contract.GetTokenPrices(&_PriceRegistry.CallOpts, tokens)
}

func (_PriceRegistry *PriceRegistryCallerSession) GetTokenPrices(tokens []common.Address) ([]InternalTimestampedPackedUint224, error) {
	return _PriceRegistry.Contract.GetTokenPrices(&_PriceRegistry.CallOpts, tokens)
}

func (_PriceRegistry *PriceRegistryCaller) GetTokenTransferFeeConfig(opts *bind.CallOpts, destChainSelector uint64, token common.Address) (PriceRegistryTokenTransferFeeConfig, error) {
	var out []interface{}
	err := _PriceRegistry.contract.Call(opts, &out, "getTokenTransferFeeConfig", destChainSelector, token)

	if err != nil {
		return *new(PriceRegistryTokenTransferFeeConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(PriceRegistryTokenTransferFeeConfig)).(*PriceRegistryTokenTransferFeeConfig)

	return out0, err

}

func (_PriceRegistry *PriceRegistrySession) GetTokenTransferFeeConfig(destChainSelector uint64, token common.Address) (PriceRegistryTokenTransferFeeConfig, error) {
	return _PriceRegistry.Contract.GetTokenTransferFeeConfig(&_PriceRegistry.CallOpts, destChainSelector, token)
}

func (_PriceRegistry *PriceRegistryCallerSession) GetTokenTransferFeeConfig(destChainSelector uint64, token common.Address) (PriceRegistryTokenTransferFeeConfig, error) {
	return _PriceRegistry.Contract.GetTokenTransferFeeConfig(&_PriceRegistry.CallOpts, destChainSelector, token)
}

func (_PriceRegistry *PriceRegistryCaller) GetValidatedFee(opts *bind.CallOpts, destChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error) {
	var out []interface{}
	err := _PriceRegistry.contract.Call(opts, &out, "getValidatedFee", destChainSelector, message)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_PriceRegistry *PriceRegistrySession) GetValidatedFee(destChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error) {
	return _PriceRegistry.Contract.GetValidatedFee(&_PriceRegistry.CallOpts, destChainSelector, message)
}

func (_PriceRegistry *PriceRegistryCallerSession) GetValidatedFee(destChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error) {
	return _PriceRegistry.Contract.GetValidatedFee(&_PriceRegistry.CallOpts, destChainSelector, message)
}

func (_PriceRegistry *PriceRegistryCaller) GetValidatedTokenPrice(opts *bind.CallOpts, token common.Address) (*big.Int, error) {
	var out []interface{}
	err := _PriceRegistry.contract.Call(opts, &out, "getValidatedTokenPrice", token)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_PriceRegistry *PriceRegistrySession) GetValidatedTokenPrice(token common.Address) (*big.Int, error) {
	return _PriceRegistry.Contract.GetValidatedTokenPrice(&_PriceRegistry.CallOpts, token)
}

func (_PriceRegistry *PriceRegistryCallerSession) GetValidatedTokenPrice(token common.Address) (*big.Int, error) {
	return _PriceRegistry.Contract.GetValidatedTokenPrice(&_PriceRegistry.CallOpts, token)
}

func (_PriceRegistry *PriceRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PriceRegistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_PriceRegistry *PriceRegistrySession) Owner() (common.Address, error) {
	return _PriceRegistry.Contract.Owner(&_PriceRegistry.CallOpts)
}

func (_PriceRegistry *PriceRegistryCallerSession) Owner() (common.Address, error) {
	return _PriceRegistry.Contract.Owner(&_PriceRegistry.CallOpts)
}

func (_PriceRegistry *PriceRegistryCaller) ProcessMessageArgs(opts *bind.CallOpts, destChainSelector uint64, feeToken common.Address, feeTokenAmount *big.Int, extraArgs []byte) (ProcessMessageArgs,

	error) {
	var out []interface{}
	err := _PriceRegistry.contract.Call(opts, &out, "processMessageArgs", destChainSelector, feeToken, feeTokenAmount, extraArgs)

	outstruct := new(ProcessMessageArgs)
	if err != nil {
		return *outstruct, err
	}

	outstruct.MsgFeeJuels = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.IsOutOfOrderExecution = *abi.ConvertType(out[1], new(bool)).(*bool)
	outstruct.ConvertedExtraArgs = *abi.ConvertType(out[2], new([]byte)).(*[]byte)

	return *outstruct, err

}

func (_PriceRegistry *PriceRegistrySession) ProcessMessageArgs(destChainSelector uint64, feeToken common.Address, feeTokenAmount *big.Int, extraArgs []byte) (ProcessMessageArgs,

	error) {
	return _PriceRegistry.Contract.ProcessMessageArgs(&_PriceRegistry.CallOpts, destChainSelector, feeToken, feeTokenAmount, extraArgs)
}

func (_PriceRegistry *PriceRegistryCallerSession) ProcessMessageArgs(destChainSelector uint64, feeToken common.Address, feeTokenAmount *big.Int, extraArgs []byte) (ProcessMessageArgs,

	error) {
	return _PriceRegistry.Contract.ProcessMessageArgs(&_PriceRegistry.CallOpts, destChainSelector, feeToken, feeTokenAmount, extraArgs)
}

func (_PriceRegistry *PriceRegistryCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _PriceRegistry.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_PriceRegistry *PriceRegistrySession) TypeAndVersion() (string, error) {
	return _PriceRegistry.Contract.TypeAndVersion(&_PriceRegistry.CallOpts)
}

func (_PriceRegistry *PriceRegistryCallerSession) TypeAndVersion() (string, error) {
	return _PriceRegistry.Contract.TypeAndVersion(&_PriceRegistry.CallOpts)
}

func (_PriceRegistry *PriceRegistryCaller) ValidatePoolReturnData(opts *bind.CallOpts, destChainSelector uint64, rampTokenAmounts []InternalRampTokenAmount, sourceTokenAmounts []ClientEVMTokenAmount) error {
	var out []interface{}
	err := _PriceRegistry.contract.Call(opts, &out, "validatePoolReturnData", destChainSelector, rampTokenAmounts, sourceTokenAmounts)

	if err != nil {
		return err
	}

	return err

}

func (_PriceRegistry *PriceRegistrySession) ValidatePoolReturnData(destChainSelector uint64, rampTokenAmounts []InternalRampTokenAmount, sourceTokenAmounts []ClientEVMTokenAmount) error {
	return _PriceRegistry.Contract.ValidatePoolReturnData(&_PriceRegistry.CallOpts, destChainSelector, rampTokenAmounts, sourceTokenAmounts)
}

func (_PriceRegistry *PriceRegistryCallerSession) ValidatePoolReturnData(destChainSelector uint64, rampTokenAmounts []InternalRampTokenAmount, sourceTokenAmounts []ClientEVMTokenAmount) error {
	return _PriceRegistry.Contract.ValidatePoolReturnData(&_PriceRegistry.CallOpts, destChainSelector, rampTokenAmounts, sourceTokenAmounts)
}

func (_PriceRegistry *PriceRegistryTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PriceRegistry.contract.Transact(opts, "acceptOwnership")
}

func (_PriceRegistry *PriceRegistrySession) AcceptOwnership() (*types.Transaction, error) {
	return _PriceRegistry.Contract.AcceptOwnership(&_PriceRegistry.TransactOpts)
}

func (_PriceRegistry *PriceRegistryTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _PriceRegistry.Contract.AcceptOwnership(&_PriceRegistry.TransactOpts)
}

func (_PriceRegistry *PriceRegistryTransactor) ApplyAuthorizedCallerUpdates(opts *bind.TransactOpts, authorizedCallerArgs AuthorizedCallersAuthorizedCallerArgs) (*types.Transaction, error) {
	return _PriceRegistry.contract.Transact(opts, "applyAuthorizedCallerUpdates", authorizedCallerArgs)
}

func (_PriceRegistry *PriceRegistrySession) ApplyAuthorizedCallerUpdates(authorizedCallerArgs AuthorizedCallersAuthorizedCallerArgs) (*types.Transaction, error) {
	return _PriceRegistry.Contract.ApplyAuthorizedCallerUpdates(&_PriceRegistry.TransactOpts, authorizedCallerArgs)
}

func (_PriceRegistry *PriceRegistryTransactorSession) ApplyAuthorizedCallerUpdates(authorizedCallerArgs AuthorizedCallersAuthorizedCallerArgs) (*types.Transaction, error) {
	return _PriceRegistry.Contract.ApplyAuthorizedCallerUpdates(&_PriceRegistry.TransactOpts, authorizedCallerArgs)
}

func (_PriceRegistry *PriceRegistryTransactor) ApplyDestChainConfigUpdates(opts *bind.TransactOpts, destChainConfigArgs []PriceRegistryDestChainConfigArgs) (*types.Transaction, error) {
	return _PriceRegistry.contract.Transact(opts, "applyDestChainConfigUpdates", destChainConfigArgs)
}

func (_PriceRegistry *PriceRegistrySession) ApplyDestChainConfigUpdates(destChainConfigArgs []PriceRegistryDestChainConfigArgs) (*types.Transaction, error) {
	return _PriceRegistry.Contract.ApplyDestChainConfigUpdates(&_PriceRegistry.TransactOpts, destChainConfigArgs)
}

func (_PriceRegistry *PriceRegistryTransactorSession) ApplyDestChainConfigUpdates(destChainConfigArgs []PriceRegistryDestChainConfigArgs) (*types.Transaction, error) {
	return _PriceRegistry.Contract.ApplyDestChainConfigUpdates(&_PriceRegistry.TransactOpts, destChainConfigArgs)
}

func (_PriceRegistry *PriceRegistryTransactor) ApplyFeeTokensUpdates(opts *bind.TransactOpts, feeTokensToAdd []common.Address, feeTokensToRemove []common.Address) (*types.Transaction, error) {
	return _PriceRegistry.contract.Transact(opts, "applyFeeTokensUpdates", feeTokensToAdd, feeTokensToRemove)
}

func (_PriceRegistry *PriceRegistrySession) ApplyFeeTokensUpdates(feeTokensToAdd []common.Address, feeTokensToRemove []common.Address) (*types.Transaction, error) {
	return _PriceRegistry.Contract.ApplyFeeTokensUpdates(&_PriceRegistry.TransactOpts, feeTokensToAdd, feeTokensToRemove)
}

func (_PriceRegistry *PriceRegistryTransactorSession) ApplyFeeTokensUpdates(feeTokensToAdd []common.Address, feeTokensToRemove []common.Address) (*types.Transaction, error) {
	return _PriceRegistry.Contract.ApplyFeeTokensUpdates(&_PriceRegistry.TransactOpts, feeTokensToAdd, feeTokensToRemove)
}

func (_PriceRegistry *PriceRegistryTransactor) ApplyPremiumMultiplierWeiPerEthUpdates(opts *bind.TransactOpts, premiumMultiplierWeiPerEthArgs []PriceRegistryPremiumMultiplierWeiPerEthArgs) (*types.Transaction, error) {
	return _PriceRegistry.contract.Transact(opts, "applyPremiumMultiplierWeiPerEthUpdates", premiumMultiplierWeiPerEthArgs)
}

func (_PriceRegistry *PriceRegistrySession) ApplyPremiumMultiplierWeiPerEthUpdates(premiumMultiplierWeiPerEthArgs []PriceRegistryPremiumMultiplierWeiPerEthArgs) (*types.Transaction, error) {
	return _PriceRegistry.Contract.ApplyPremiumMultiplierWeiPerEthUpdates(&_PriceRegistry.TransactOpts, premiumMultiplierWeiPerEthArgs)
}

func (_PriceRegistry *PriceRegistryTransactorSession) ApplyPremiumMultiplierWeiPerEthUpdates(premiumMultiplierWeiPerEthArgs []PriceRegistryPremiumMultiplierWeiPerEthArgs) (*types.Transaction, error) {
	return _PriceRegistry.Contract.ApplyPremiumMultiplierWeiPerEthUpdates(&_PriceRegistry.TransactOpts, premiumMultiplierWeiPerEthArgs)
}

func (_PriceRegistry *PriceRegistryTransactor) ApplyTokenTransferFeeConfigUpdates(opts *bind.TransactOpts, tokenTransferFeeConfigArgs []PriceRegistryTokenTransferFeeConfigArgs, tokensToUseDefaultFeeConfigs []PriceRegistryTokenTransferFeeConfigRemoveArgs) (*types.Transaction, error) {
	return _PriceRegistry.contract.Transact(opts, "applyTokenTransferFeeConfigUpdates", tokenTransferFeeConfigArgs, tokensToUseDefaultFeeConfigs)
}

func (_PriceRegistry *PriceRegistrySession) ApplyTokenTransferFeeConfigUpdates(tokenTransferFeeConfigArgs []PriceRegistryTokenTransferFeeConfigArgs, tokensToUseDefaultFeeConfigs []PriceRegistryTokenTransferFeeConfigRemoveArgs) (*types.Transaction, error) {
	return _PriceRegistry.Contract.ApplyTokenTransferFeeConfigUpdates(&_PriceRegistry.TransactOpts, tokenTransferFeeConfigArgs, tokensToUseDefaultFeeConfigs)
}

func (_PriceRegistry *PriceRegistryTransactorSession) ApplyTokenTransferFeeConfigUpdates(tokenTransferFeeConfigArgs []PriceRegistryTokenTransferFeeConfigArgs, tokensToUseDefaultFeeConfigs []PriceRegistryTokenTransferFeeConfigRemoveArgs) (*types.Transaction, error) {
	return _PriceRegistry.Contract.ApplyTokenTransferFeeConfigUpdates(&_PriceRegistry.TransactOpts, tokenTransferFeeConfigArgs, tokensToUseDefaultFeeConfigs)
}

func (_PriceRegistry *PriceRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _PriceRegistry.contract.Transact(opts, "transferOwnership", to)
}

func (_PriceRegistry *PriceRegistrySession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _PriceRegistry.Contract.TransferOwnership(&_PriceRegistry.TransactOpts, to)
}

func (_PriceRegistry *PriceRegistryTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _PriceRegistry.Contract.TransferOwnership(&_PriceRegistry.TransactOpts, to)
}

func (_PriceRegistry *PriceRegistryTransactor) UpdatePrices(opts *bind.TransactOpts, priceUpdates InternalPriceUpdates) (*types.Transaction, error) {
	return _PriceRegistry.contract.Transact(opts, "updatePrices", priceUpdates)
}

func (_PriceRegistry *PriceRegistrySession) UpdatePrices(priceUpdates InternalPriceUpdates) (*types.Transaction, error) {
	return _PriceRegistry.Contract.UpdatePrices(&_PriceRegistry.TransactOpts, priceUpdates)
}

func (_PriceRegistry *PriceRegistryTransactorSession) UpdatePrices(priceUpdates InternalPriceUpdates) (*types.Transaction, error) {
	return _PriceRegistry.Contract.UpdatePrices(&_PriceRegistry.TransactOpts, priceUpdates)
}

func (_PriceRegistry *PriceRegistryTransactor) UpdateTokenPriceFeeds(opts *bind.TransactOpts, tokenPriceFeedUpdates []PriceRegistryTokenPriceFeedUpdate) (*types.Transaction, error) {
	return _PriceRegistry.contract.Transact(opts, "updateTokenPriceFeeds", tokenPriceFeedUpdates)
}

func (_PriceRegistry *PriceRegistrySession) UpdateTokenPriceFeeds(tokenPriceFeedUpdates []PriceRegistryTokenPriceFeedUpdate) (*types.Transaction, error) {
	return _PriceRegistry.Contract.UpdateTokenPriceFeeds(&_PriceRegistry.TransactOpts, tokenPriceFeedUpdates)
}

func (_PriceRegistry *PriceRegistryTransactorSession) UpdateTokenPriceFeeds(tokenPriceFeedUpdates []PriceRegistryTokenPriceFeedUpdate) (*types.Transaction, error) {
	return _PriceRegistry.Contract.UpdateTokenPriceFeeds(&_PriceRegistry.TransactOpts, tokenPriceFeedUpdates)
}

type PriceRegistryAuthorizedCallerAddedIterator struct {
	Event *PriceRegistryAuthorizedCallerAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PriceRegistryAuthorizedCallerAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PriceRegistryAuthorizedCallerAdded)
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
		it.Event = new(PriceRegistryAuthorizedCallerAdded)
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

func (it *PriceRegistryAuthorizedCallerAddedIterator) Error() error {
	return it.fail
}

func (it *PriceRegistryAuthorizedCallerAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PriceRegistryAuthorizedCallerAdded struct {
	Caller common.Address
	Raw    types.Log
}

func (_PriceRegistry *PriceRegistryFilterer) FilterAuthorizedCallerAdded(opts *bind.FilterOpts) (*PriceRegistryAuthorizedCallerAddedIterator, error) {

	logs, sub, err := _PriceRegistry.contract.FilterLogs(opts, "AuthorizedCallerAdded")
	if err != nil {
		return nil, err
	}
	return &PriceRegistryAuthorizedCallerAddedIterator{contract: _PriceRegistry.contract, event: "AuthorizedCallerAdded", logs: logs, sub: sub}, nil
}

func (_PriceRegistry *PriceRegistryFilterer) WatchAuthorizedCallerAdded(opts *bind.WatchOpts, sink chan<- *PriceRegistryAuthorizedCallerAdded) (event.Subscription, error) {

	logs, sub, err := _PriceRegistry.contract.WatchLogs(opts, "AuthorizedCallerAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PriceRegistryAuthorizedCallerAdded)
				if err := _PriceRegistry.contract.UnpackLog(event, "AuthorizedCallerAdded", log); err != nil {
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

func (_PriceRegistry *PriceRegistryFilterer) ParseAuthorizedCallerAdded(log types.Log) (*PriceRegistryAuthorizedCallerAdded, error) {
	event := new(PriceRegistryAuthorizedCallerAdded)
	if err := _PriceRegistry.contract.UnpackLog(event, "AuthorizedCallerAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PriceRegistryAuthorizedCallerRemovedIterator struct {
	Event *PriceRegistryAuthorizedCallerRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PriceRegistryAuthorizedCallerRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PriceRegistryAuthorizedCallerRemoved)
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
		it.Event = new(PriceRegistryAuthorizedCallerRemoved)
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

func (it *PriceRegistryAuthorizedCallerRemovedIterator) Error() error {
	return it.fail
}

func (it *PriceRegistryAuthorizedCallerRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PriceRegistryAuthorizedCallerRemoved struct {
	Caller common.Address
	Raw    types.Log
}

func (_PriceRegistry *PriceRegistryFilterer) FilterAuthorizedCallerRemoved(opts *bind.FilterOpts) (*PriceRegistryAuthorizedCallerRemovedIterator, error) {

	logs, sub, err := _PriceRegistry.contract.FilterLogs(opts, "AuthorizedCallerRemoved")
	if err != nil {
		return nil, err
	}
	return &PriceRegistryAuthorizedCallerRemovedIterator{contract: _PriceRegistry.contract, event: "AuthorizedCallerRemoved", logs: logs, sub: sub}, nil
}

func (_PriceRegistry *PriceRegistryFilterer) WatchAuthorizedCallerRemoved(opts *bind.WatchOpts, sink chan<- *PriceRegistryAuthorizedCallerRemoved) (event.Subscription, error) {

	logs, sub, err := _PriceRegistry.contract.WatchLogs(opts, "AuthorizedCallerRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PriceRegistryAuthorizedCallerRemoved)
				if err := _PriceRegistry.contract.UnpackLog(event, "AuthorizedCallerRemoved", log); err != nil {
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

func (_PriceRegistry *PriceRegistryFilterer) ParseAuthorizedCallerRemoved(log types.Log) (*PriceRegistryAuthorizedCallerRemoved, error) {
	event := new(PriceRegistryAuthorizedCallerRemoved)
	if err := _PriceRegistry.contract.UnpackLog(event, "AuthorizedCallerRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PriceRegistryDestChainAddedIterator struct {
	Event *PriceRegistryDestChainAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PriceRegistryDestChainAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PriceRegistryDestChainAdded)
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
		it.Event = new(PriceRegistryDestChainAdded)
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

func (it *PriceRegistryDestChainAddedIterator) Error() error {
	return it.fail
}

func (it *PriceRegistryDestChainAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PriceRegistryDestChainAdded struct {
	DestChainSelector uint64
	DestChainConfig   PriceRegistryDestChainConfig
	Raw               types.Log
}

func (_PriceRegistry *PriceRegistryFilterer) FilterDestChainAdded(opts *bind.FilterOpts, destChainSelector []uint64) (*PriceRegistryDestChainAddedIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _PriceRegistry.contract.FilterLogs(opts, "DestChainAdded", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &PriceRegistryDestChainAddedIterator{contract: _PriceRegistry.contract, event: "DestChainAdded", logs: logs, sub: sub}, nil
}

func (_PriceRegistry *PriceRegistryFilterer) WatchDestChainAdded(opts *bind.WatchOpts, sink chan<- *PriceRegistryDestChainAdded, destChainSelector []uint64) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _PriceRegistry.contract.WatchLogs(opts, "DestChainAdded", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PriceRegistryDestChainAdded)
				if err := _PriceRegistry.contract.UnpackLog(event, "DestChainAdded", log); err != nil {
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

func (_PriceRegistry *PriceRegistryFilterer) ParseDestChainAdded(log types.Log) (*PriceRegistryDestChainAdded, error) {
	event := new(PriceRegistryDestChainAdded)
	if err := _PriceRegistry.contract.UnpackLog(event, "DestChainAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PriceRegistryDestChainConfigUpdatedIterator struct {
	Event *PriceRegistryDestChainConfigUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PriceRegistryDestChainConfigUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PriceRegistryDestChainConfigUpdated)
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
		it.Event = new(PriceRegistryDestChainConfigUpdated)
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

func (it *PriceRegistryDestChainConfigUpdatedIterator) Error() error {
	return it.fail
}

func (it *PriceRegistryDestChainConfigUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PriceRegistryDestChainConfigUpdated struct {
	DestChainSelector uint64
	DestChainConfig   PriceRegistryDestChainConfig
	Raw               types.Log
}

func (_PriceRegistry *PriceRegistryFilterer) FilterDestChainConfigUpdated(opts *bind.FilterOpts, destChainSelector []uint64) (*PriceRegistryDestChainConfigUpdatedIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _PriceRegistry.contract.FilterLogs(opts, "DestChainConfigUpdated", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &PriceRegistryDestChainConfigUpdatedIterator{contract: _PriceRegistry.contract, event: "DestChainConfigUpdated", logs: logs, sub: sub}, nil
}

func (_PriceRegistry *PriceRegistryFilterer) WatchDestChainConfigUpdated(opts *bind.WatchOpts, sink chan<- *PriceRegistryDestChainConfigUpdated, destChainSelector []uint64) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _PriceRegistry.contract.WatchLogs(opts, "DestChainConfigUpdated", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PriceRegistryDestChainConfigUpdated)
				if err := _PriceRegistry.contract.UnpackLog(event, "DestChainConfigUpdated", log); err != nil {
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

func (_PriceRegistry *PriceRegistryFilterer) ParseDestChainConfigUpdated(log types.Log) (*PriceRegistryDestChainConfigUpdated, error) {
	event := new(PriceRegistryDestChainConfigUpdated)
	if err := _PriceRegistry.contract.UnpackLog(event, "DestChainConfigUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PriceRegistryFeeTokenAddedIterator struct {
	Event *PriceRegistryFeeTokenAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PriceRegistryFeeTokenAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PriceRegistryFeeTokenAdded)
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
		it.Event = new(PriceRegistryFeeTokenAdded)
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

func (it *PriceRegistryFeeTokenAddedIterator) Error() error {
	return it.fail
}

func (it *PriceRegistryFeeTokenAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PriceRegistryFeeTokenAdded struct {
	FeeToken common.Address
	Raw      types.Log
}

func (_PriceRegistry *PriceRegistryFilterer) FilterFeeTokenAdded(opts *bind.FilterOpts, feeToken []common.Address) (*PriceRegistryFeeTokenAddedIterator, error) {

	var feeTokenRule []interface{}
	for _, feeTokenItem := range feeToken {
		feeTokenRule = append(feeTokenRule, feeTokenItem)
	}

	logs, sub, err := _PriceRegistry.contract.FilterLogs(opts, "FeeTokenAdded", feeTokenRule)
	if err != nil {
		return nil, err
	}
	return &PriceRegistryFeeTokenAddedIterator{contract: _PriceRegistry.contract, event: "FeeTokenAdded", logs: logs, sub: sub}, nil
}

func (_PriceRegistry *PriceRegistryFilterer) WatchFeeTokenAdded(opts *bind.WatchOpts, sink chan<- *PriceRegistryFeeTokenAdded, feeToken []common.Address) (event.Subscription, error) {

	var feeTokenRule []interface{}
	for _, feeTokenItem := range feeToken {
		feeTokenRule = append(feeTokenRule, feeTokenItem)
	}

	logs, sub, err := _PriceRegistry.contract.WatchLogs(opts, "FeeTokenAdded", feeTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PriceRegistryFeeTokenAdded)
				if err := _PriceRegistry.contract.UnpackLog(event, "FeeTokenAdded", log); err != nil {
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

func (_PriceRegistry *PriceRegistryFilterer) ParseFeeTokenAdded(log types.Log) (*PriceRegistryFeeTokenAdded, error) {
	event := new(PriceRegistryFeeTokenAdded)
	if err := _PriceRegistry.contract.UnpackLog(event, "FeeTokenAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PriceRegistryFeeTokenRemovedIterator struct {
	Event *PriceRegistryFeeTokenRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PriceRegistryFeeTokenRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PriceRegistryFeeTokenRemoved)
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
		it.Event = new(PriceRegistryFeeTokenRemoved)
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

func (it *PriceRegistryFeeTokenRemovedIterator) Error() error {
	return it.fail
}

func (it *PriceRegistryFeeTokenRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PriceRegistryFeeTokenRemoved struct {
	FeeToken common.Address
	Raw      types.Log
}

func (_PriceRegistry *PriceRegistryFilterer) FilterFeeTokenRemoved(opts *bind.FilterOpts, feeToken []common.Address) (*PriceRegistryFeeTokenRemovedIterator, error) {

	var feeTokenRule []interface{}
	for _, feeTokenItem := range feeToken {
		feeTokenRule = append(feeTokenRule, feeTokenItem)
	}

	logs, sub, err := _PriceRegistry.contract.FilterLogs(opts, "FeeTokenRemoved", feeTokenRule)
	if err != nil {
		return nil, err
	}
	return &PriceRegistryFeeTokenRemovedIterator{contract: _PriceRegistry.contract, event: "FeeTokenRemoved", logs: logs, sub: sub}, nil
}

func (_PriceRegistry *PriceRegistryFilterer) WatchFeeTokenRemoved(opts *bind.WatchOpts, sink chan<- *PriceRegistryFeeTokenRemoved, feeToken []common.Address) (event.Subscription, error) {

	var feeTokenRule []interface{}
	for _, feeTokenItem := range feeToken {
		feeTokenRule = append(feeTokenRule, feeTokenItem)
	}

	logs, sub, err := _PriceRegistry.contract.WatchLogs(opts, "FeeTokenRemoved", feeTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PriceRegistryFeeTokenRemoved)
				if err := _PriceRegistry.contract.UnpackLog(event, "FeeTokenRemoved", log); err != nil {
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

func (_PriceRegistry *PriceRegistryFilterer) ParseFeeTokenRemoved(log types.Log) (*PriceRegistryFeeTokenRemoved, error) {
	event := new(PriceRegistryFeeTokenRemoved)
	if err := _PriceRegistry.contract.UnpackLog(event, "FeeTokenRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PriceRegistryOwnershipTransferRequestedIterator struct {
	Event *PriceRegistryOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PriceRegistryOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PriceRegistryOwnershipTransferRequested)
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
		it.Event = new(PriceRegistryOwnershipTransferRequested)
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

func (it *PriceRegistryOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *PriceRegistryOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PriceRegistryOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_PriceRegistry *PriceRegistryFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*PriceRegistryOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _PriceRegistry.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &PriceRegistryOwnershipTransferRequestedIterator{contract: _PriceRegistry.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_PriceRegistry *PriceRegistryFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *PriceRegistryOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _PriceRegistry.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PriceRegistryOwnershipTransferRequested)
				if err := _PriceRegistry.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_PriceRegistry *PriceRegistryFilterer) ParseOwnershipTransferRequested(log types.Log) (*PriceRegistryOwnershipTransferRequested, error) {
	event := new(PriceRegistryOwnershipTransferRequested)
	if err := _PriceRegistry.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PriceRegistryOwnershipTransferredIterator struct {
	Event *PriceRegistryOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PriceRegistryOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PriceRegistryOwnershipTransferred)
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
		it.Event = new(PriceRegistryOwnershipTransferred)
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

func (it *PriceRegistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *PriceRegistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PriceRegistryOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_PriceRegistry *PriceRegistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*PriceRegistryOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _PriceRegistry.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &PriceRegistryOwnershipTransferredIterator{contract: _PriceRegistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_PriceRegistry *PriceRegistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *PriceRegistryOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _PriceRegistry.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PriceRegistryOwnershipTransferred)
				if err := _PriceRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_PriceRegistry *PriceRegistryFilterer) ParseOwnershipTransferred(log types.Log) (*PriceRegistryOwnershipTransferred, error) {
	event := new(PriceRegistryOwnershipTransferred)
	if err := _PriceRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PriceRegistryPremiumMultiplierWeiPerEthUpdatedIterator struct {
	Event *PriceRegistryPremiumMultiplierWeiPerEthUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PriceRegistryPremiumMultiplierWeiPerEthUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PriceRegistryPremiumMultiplierWeiPerEthUpdated)
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
		it.Event = new(PriceRegistryPremiumMultiplierWeiPerEthUpdated)
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

func (it *PriceRegistryPremiumMultiplierWeiPerEthUpdatedIterator) Error() error {
	return it.fail
}

func (it *PriceRegistryPremiumMultiplierWeiPerEthUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PriceRegistryPremiumMultiplierWeiPerEthUpdated struct {
	Token                      common.Address
	PremiumMultiplierWeiPerEth uint64
	Raw                        types.Log
}

func (_PriceRegistry *PriceRegistryFilterer) FilterPremiumMultiplierWeiPerEthUpdated(opts *bind.FilterOpts, token []common.Address) (*PriceRegistryPremiumMultiplierWeiPerEthUpdatedIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _PriceRegistry.contract.FilterLogs(opts, "PremiumMultiplierWeiPerEthUpdated", tokenRule)
	if err != nil {
		return nil, err
	}
	return &PriceRegistryPremiumMultiplierWeiPerEthUpdatedIterator{contract: _PriceRegistry.contract, event: "PremiumMultiplierWeiPerEthUpdated", logs: logs, sub: sub}, nil
}

func (_PriceRegistry *PriceRegistryFilterer) WatchPremiumMultiplierWeiPerEthUpdated(opts *bind.WatchOpts, sink chan<- *PriceRegistryPremiumMultiplierWeiPerEthUpdated, token []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _PriceRegistry.contract.WatchLogs(opts, "PremiumMultiplierWeiPerEthUpdated", tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PriceRegistryPremiumMultiplierWeiPerEthUpdated)
				if err := _PriceRegistry.contract.UnpackLog(event, "PremiumMultiplierWeiPerEthUpdated", log); err != nil {
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

func (_PriceRegistry *PriceRegistryFilterer) ParsePremiumMultiplierWeiPerEthUpdated(log types.Log) (*PriceRegistryPremiumMultiplierWeiPerEthUpdated, error) {
	event := new(PriceRegistryPremiumMultiplierWeiPerEthUpdated)
	if err := _PriceRegistry.contract.UnpackLog(event, "PremiumMultiplierWeiPerEthUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PriceRegistryPriceFeedPerTokenUpdatedIterator struct {
	Event *PriceRegistryPriceFeedPerTokenUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PriceRegistryPriceFeedPerTokenUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PriceRegistryPriceFeedPerTokenUpdated)
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
		it.Event = new(PriceRegistryPriceFeedPerTokenUpdated)
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

func (it *PriceRegistryPriceFeedPerTokenUpdatedIterator) Error() error {
	return it.fail
}

func (it *PriceRegistryPriceFeedPerTokenUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PriceRegistryPriceFeedPerTokenUpdated struct {
	Token           common.Address
	PriceFeedConfig IPriceRegistryTokenPriceFeedConfig
	Raw             types.Log
}

func (_PriceRegistry *PriceRegistryFilterer) FilterPriceFeedPerTokenUpdated(opts *bind.FilterOpts, token []common.Address) (*PriceRegistryPriceFeedPerTokenUpdatedIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _PriceRegistry.contract.FilterLogs(opts, "PriceFeedPerTokenUpdated", tokenRule)
	if err != nil {
		return nil, err
	}
	return &PriceRegistryPriceFeedPerTokenUpdatedIterator{contract: _PriceRegistry.contract, event: "PriceFeedPerTokenUpdated", logs: logs, sub: sub}, nil
}

func (_PriceRegistry *PriceRegistryFilterer) WatchPriceFeedPerTokenUpdated(opts *bind.WatchOpts, sink chan<- *PriceRegistryPriceFeedPerTokenUpdated, token []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _PriceRegistry.contract.WatchLogs(opts, "PriceFeedPerTokenUpdated", tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PriceRegistryPriceFeedPerTokenUpdated)
				if err := _PriceRegistry.contract.UnpackLog(event, "PriceFeedPerTokenUpdated", log); err != nil {
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

func (_PriceRegistry *PriceRegistryFilterer) ParsePriceFeedPerTokenUpdated(log types.Log) (*PriceRegistryPriceFeedPerTokenUpdated, error) {
	event := new(PriceRegistryPriceFeedPerTokenUpdated)
	if err := _PriceRegistry.contract.UnpackLog(event, "PriceFeedPerTokenUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PriceRegistryTokenTransferFeeConfigDeletedIterator struct {
	Event *PriceRegistryTokenTransferFeeConfigDeleted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PriceRegistryTokenTransferFeeConfigDeletedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PriceRegistryTokenTransferFeeConfigDeleted)
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
		it.Event = new(PriceRegistryTokenTransferFeeConfigDeleted)
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

func (it *PriceRegistryTokenTransferFeeConfigDeletedIterator) Error() error {
	return it.fail
}

func (it *PriceRegistryTokenTransferFeeConfigDeletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PriceRegistryTokenTransferFeeConfigDeleted struct {
	DestChainSelector uint64
	Token             common.Address
	Raw               types.Log
}

func (_PriceRegistry *PriceRegistryFilterer) FilterTokenTransferFeeConfigDeleted(opts *bind.FilterOpts, destChainSelector []uint64, token []common.Address) (*PriceRegistryTokenTransferFeeConfigDeletedIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _PriceRegistry.contract.FilterLogs(opts, "TokenTransferFeeConfigDeleted", destChainSelectorRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return &PriceRegistryTokenTransferFeeConfigDeletedIterator{contract: _PriceRegistry.contract, event: "TokenTransferFeeConfigDeleted", logs: logs, sub: sub}, nil
}

func (_PriceRegistry *PriceRegistryFilterer) WatchTokenTransferFeeConfigDeleted(opts *bind.WatchOpts, sink chan<- *PriceRegistryTokenTransferFeeConfigDeleted, destChainSelector []uint64, token []common.Address) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _PriceRegistry.contract.WatchLogs(opts, "TokenTransferFeeConfigDeleted", destChainSelectorRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PriceRegistryTokenTransferFeeConfigDeleted)
				if err := _PriceRegistry.contract.UnpackLog(event, "TokenTransferFeeConfigDeleted", log); err != nil {
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

func (_PriceRegistry *PriceRegistryFilterer) ParseTokenTransferFeeConfigDeleted(log types.Log) (*PriceRegistryTokenTransferFeeConfigDeleted, error) {
	event := new(PriceRegistryTokenTransferFeeConfigDeleted)
	if err := _PriceRegistry.contract.UnpackLog(event, "TokenTransferFeeConfigDeleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PriceRegistryTokenTransferFeeConfigUpdatedIterator struct {
	Event *PriceRegistryTokenTransferFeeConfigUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PriceRegistryTokenTransferFeeConfigUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PriceRegistryTokenTransferFeeConfigUpdated)
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
		it.Event = new(PriceRegistryTokenTransferFeeConfigUpdated)
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

func (it *PriceRegistryTokenTransferFeeConfigUpdatedIterator) Error() error {
	return it.fail
}

func (it *PriceRegistryTokenTransferFeeConfigUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PriceRegistryTokenTransferFeeConfigUpdated struct {
	DestChainSelector      uint64
	Token                  common.Address
	TokenTransferFeeConfig PriceRegistryTokenTransferFeeConfig
	Raw                    types.Log
}

func (_PriceRegistry *PriceRegistryFilterer) FilterTokenTransferFeeConfigUpdated(opts *bind.FilterOpts, destChainSelector []uint64, token []common.Address) (*PriceRegistryTokenTransferFeeConfigUpdatedIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _PriceRegistry.contract.FilterLogs(opts, "TokenTransferFeeConfigUpdated", destChainSelectorRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return &PriceRegistryTokenTransferFeeConfigUpdatedIterator{contract: _PriceRegistry.contract, event: "TokenTransferFeeConfigUpdated", logs: logs, sub: sub}, nil
}

func (_PriceRegistry *PriceRegistryFilterer) WatchTokenTransferFeeConfigUpdated(opts *bind.WatchOpts, sink chan<- *PriceRegistryTokenTransferFeeConfigUpdated, destChainSelector []uint64, token []common.Address) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _PriceRegistry.contract.WatchLogs(opts, "TokenTransferFeeConfigUpdated", destChainSelectorRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PriceRegistryTokenTransferFeeConfigUpdated)
				if err := _PriceRegistry.contract.UnpackLog(event, "TokenTransferFeeConfigUpdated", log); err != nil {
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

func (_PriceRegistry *PriceRegistryFilterer) ParseTokenTransferFeeConfigUpdated(log types.Log) (*PriceRegistryTokenTransferFeeConfigUpdated, error) {
	event := new(PriceRegistryTokenTransferFeeConfigUpdated)
	if err := _PriceRegistry.contract.UnpackLog(event, "TokenTransferFeeConfigUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PriceRegistryUsdPerTokenUpdatedIterator struct {
	Event *PriceRegistryUsdPerTokenUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PriceRegistryUsdPerTokenUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PriceRegistryUsdPerTokenUpdated)
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
		it.Event = new(PriceRegistryUsdPerTokenUpdated)
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

func (it *PriceRegistryUsdPerTokenUpdatedIterator) Error() error {
	return it.fail
}

func (it *PriceRegistryUsdPerTokenUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PriceRegistryUsdPerTokenUpdated struct {
	Token     common.Address
	Value     *big.Int
	Timestamp *big.Int
	Raw       types.Log
}

func (_PriceRegistry *PriceRegistryFilterer) FilterUsdPerTokenUpdated(opts *bind.FilterOpts, token []common.Address) (*PriceRegistryUsdPerTokenUpdatedIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _PriceRegistry.contract.FilterLogs(opts, "UsdPerTokenUpdated", tokenRule)
	if err != nil {
		return nil, err
	}
	return &PriceRegistryUsdPerTokenUpdatedIterator{contract: _PriceRegistry.contract, event: "UsdPerTokenUpdated", logs: logs, sub: sub}, nil
}

func (_PriceRegistry *PriceRegistryFilterer) WatchUsdPerTokenUpdated(opts *bind.WatchOpts, sink chan<- *PriceRegistryUsdPerTokenUpdated, token []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _PriceRegistry.contract.WatchLogs(opts, "UsdPerTokenUpdated", tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PriceRegistryUsdPerTokenUpdated)
				if err := _PriceRegistry.contract.UnpackLog(event, "UsdPerTokenUpdated", log); err != nil {
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

func (_PriceRegistry *PriceRegistryFilterer) ParseUsdPerTokenUpdated(log types.Log) (*PriceRegistryUsdPerTokenUpdated, error) {
	event := new(PriceRegistryUsdPerTokenUpdated)
	if err := _PriceRegistry.contract.UnpackLog(event, "UsdPerTokenUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PriceRegistryUsdPerUnitGasUpdatedIterator struct {
	Event *PriceRegistryUsdPerUnitGasUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PriceRegistryUsdPerUnitGasUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PriceRegistryUsdPerUnitGasUpdated)
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
		it.Event = new(PriceRegistryUsdPerUnitGasUpdated)
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

func (it *PriceRegistryUsdPerUnitGasUpdatedIterator) Error() error {
	return it.fail
}

func (it *PriceRegistryUsdPerUnitGasUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PriceRegistryUsdPerUnitGasUpdated struct {
	DestChain uint64
	Value     *big.Int
	Timestamp *big.Int
	Raw       types.Log
}

func (_PriceRegistry *PriceRegistryFilterer) FilterUsdPerUnitGasUpdated(opts *bind.FilterOpts, destChain []uint64) (*PriceRegistryUsdPerUnitGasUpdatedIterator, error) {

	var destChainRule []interface{}
	for _, destChainItem := range destChain {
		destChainRule = append(destChainRule, destChainItem)
	}

	logs, sub, err := _PriceRegistry.contract.FilterLogs(opts, "UsdPerUnitGasUpdated", destChainRule)
	if err != nil {
		return nil, err
	}
	return &PriceRegistryUsdPerUnitGasUpdatedIterator{contract: _PriceRegistry.contract, event: "UsdPerUnitGasUpdated", logs: logs, sub: sub}, nil
}

func (_PriceRegistry *PriceRegistryFilterer) WatchUsdPerUnitGasUpdated(opts *bind.WatchOpts, sink chan<- *PriceRegistryUsdPerUnitGasUpdated, destChain []uint64) (event.Subscription, error) {

	var destChainRule []interface{}
	for _, destChainItem := range destChain {
		destChainRule = append(destChainRule, destChainItem)
	}

	logs, sub, err := _PriceRegistry.contract.WatchLogs(opts, "UsdPerUnitGasUpdated", destChainRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PriceRegistryUsdPerUnitGasUpdated)
				if err := _PriceRegistry.contract.UnpackLog(event, "UsdPerUnitGasUpdated", log); err != nil {
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

func (_PriceRegistry *PriceRegistryFilterer) ParseUsdPerUnitGasUpdated(log types.Log) (*PriceRegistryUsdPerUnitGasUpdated, error) {
	event := new(PriceRegistryUsdPerUnitGasUpdated)
	if err := _PriceRegistry.contract.UnpackLog(event, "UsdPerUnitGasUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetTokenAndGasPrices struct {
	TokenPrice    *big.Int
	GasPriceValue *big.Int
}
type ProcessMessageArgs struct {
	MsgFeeJuels           *big.Int
	IsOutOfOrderExecution bool
	ConvertedExtraArgs    []byte
}

func (_PriceRegistry *PriceRegistry) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _PriceRegistry.abi.Events["AuthorizedCallerAdded"].ID:
		return _PriceRegistry.ParseAuthorizedCallerAdded(log)
	case _PriceRegistry.abi.Events["AuthorizedCallerRemoved"].ID:
		return _PriceRegistry.ParseAuthorizedCallerRemoved(log)
	case _PriceRegistry.abi.Events["DestChainAdded"].ID:
		return _PriceRegistry.ParseDestChainAdded(log)
	case _PriceRegistry.abi.Events["DestChainConfigUpdated"].ID:
		return _PriceRegistry.ParseDestChainConfigUpdated(log)
	case _PriceRegistry.abi.Events["FeeTokenAdded"].ID:
		return _PriceRegistry.ParseFeeTokenAdded(log)
	case _PriceRegistry.abi.Events["FeeTokenRemoved"].ID:
		return _PriceRegistry.ParseFeeTokenRemoved(log)
	case _PriceRegistry.abi.Events["OwnershipTransferRequested"].ID:
		return _PriceRegistry.ParseOwnershipTransferRequested(log)
	case _PriceRegistry.abi.Events["OwnershipTransferred"].ID:
		return _PriceRegistry.ParseOwnershipTransferred(log)
	case _PriceRegistry.abi.Events["PremiumMultiplierWeiPerEthUpdated"].ID:
		return _PriceRegistry.ParsePremiumMultiplierWeiPerEthUpdated(log)
	case _PriceRegistry.abi.Events["PriceFeedPerTokenUpdated"].ID:
		return _PriceRegistry.ParsePriceFeedPerTokenUpdated(log)
	case _PriceRegistry.abi.Events["TokenTransferFeeConfigDeleted"].ID:
		return _PriceRegistry.ParseTokenTransferFeeConfigDeleted(log)
	case _PriceRegistry.abi.Events["TokenTransferFeeConfigUpdated"].ID:
		return _PriceRegistry.ParseTokenTransferFeeConfigUpdated(log)
	case _PriceRegistry.abi.Events["UsdPerTokenUpdated"].ID:
		return _PriceRegistry.ParseUsdPerTokenUpdated(log)
	case _PriceRegistry.abi.Events["UsdPerUnitGasUpdated"].ID:
		return _PriceRegistry.ParseUsdPerUnitGasUpdated(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (PriceRegistryAuthorizedCallerAdded) Topic() common.Hash {
	return common.HexToHash("0xeb1b9b92e50b7f88f9ff25d56765095ac6e91540eee214906f4036a908ffbdef")
}

func (PriceRegistryAuthorizedCallerRemoved) Topic() common.Hash {
	return common.HexToHash("0xc3803387881faad271c47728894e3e36fac830ffc8602ca6fc07733cbda77580")
}

func (PriceRegistryDestChainAdded) Topic() common.Hash {
	return common.HexToHash("0xa937382a486d993de71c220bc8b559242deb4e286a353fa732330b4aa7d13577")
}

func (PriceRegistryDestChainConfigUpdated) Topic() common.Hash {
	return common.HexToHash("0xa7b607fc10d28a1caf39ab7d27f4c94945db708a576d572781a455c5894fad93")
}

func (PriceRegistryFeeTokenAdded) Topic() common.Hash {
	return common.HexToHash("0xdf1b1bd32a69711488d71554706bb130b1fc63a5fa1a2cd85e8440f84065ba23")
}

func (PriceRegistryFeeTokenRemoved) Topic() common.Hash {
	return common.HexToHash("0x1795838dc8ab2ffc5f431a1729a6afa0b587f982f7b2be0b9d7187a1ef547f91")
}

func (PriceRegistryOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (PriceRegistryOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (PriceRegistryPremiumMultiplierWeiPerEthUpdated) Topic() common.Hash {
	return common.HexToHash("0xbb77da6f7210cdd16904228a9360133d1d7dfff99b1bc75f128da5b53e28f97d")
}

func (PriceRegistryPriceFeedPerTokenUpdated) Topic() common.Hash {
	return common.HexToHash("0x08a5f7f5bb38a81d8e43aca13ecd76431dbf8816ae4699affff7b00b2fc1c464")
}

func (PriceRegistryTokenTransferFeeConfigDeleted) Topic() common.Hash {
	return common.HexToHash("0x4de5b1bcbca6018c11303a2c3f4a4b4f22a1c741d8c4ba430d246ac06c5ddf8b")
}

func (PriceRegistryTokenTransferFeeConfigUpdated) Topic() common.Hash {
	return common.HexToHash("0x94967ae9ea7729ad4f54021c1981765d2b1d954f7c92fbec340aa0a54f46b8b5")
}

func (PriceRegistryUsdPerTokenUpdated) Topic() common.Hash {
	return common.HexToHash("0x52f50aa6d1a95a4595361ecf953d095f125d442e4673716dede699e049de148a")
}

func (PriceRegistryUsdPerUnitGasUpdated) Topic() common.Hash {
	return common.HexToHash("0xdd84a3fa9ef9409f550d54d6affec7e9c480c878c6ab27b78912a03e1b371c6e")
}

func (_PriceRegistry *PriceRegistry) Address() common.Address {
	return _PriceRegistry.address
}

type PriceRegistryInterface interface {
	ConvertTokenAmount(opts *bind.CallOpts, fromToken common.Address, fromTokenAmount *big.Int, toToken common.Address) (*big.Int, error)

	GetAllAuthorizedCallers(opts *bind.CallOpts) ([]common.Address, error)

	GetDestChainConfig(opts *bind.CallOpts, destChainSelector uint64) (PriceRegistryDestChainConfig, error)

	GetDestinationChainGasPrice(opts *bind.CallOpts, destChainSelector uint64) (InternalTimestampedPackedUint224, error)

	GetFeeTokens(opts *bind.CallOpts) ([]common.Address, error)

	GetPremiumMultiplierWeiPerEth(opts *bind.CallOpts, token common.Address) (uint64, error)

	GetStaticConfig(opts *bind.CallOpts) (PriceRegistryStaticConfig, error)

	GetTokenAndGasPrices(opts *bind.CallOpts, token common.Address, destChainSelector uint64) (GetTokenAndGasPrices,

		error)

	GetTokenPrice(opts *bind.CallOpts, token common.Address) (InternalTimestampedPackedUint224, error)

	GetTokenPriceFeedConfig(opts *bind.CallOpts, token common.Address) (IPriceRegistryTokenPriceFeedConfig, error)

	GetTokenPrices(opts *bind.CallOpts, tokens []common.Address) ([]InternalTimestampedPackedUint224, error)

	GetTokenTransferFeeConfig(opts *bind.CallOpts, destChainSelector uint64, token common.Address) (PriceRegistryTokenTransferFeeConfig, error)

	GetValidatedFee(opts *bind.CallOpts, destChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error)

	GetValidatedTokenPrice(opts *bind.CallOpts, token common.Address) (*big.Int, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	ProcessMessageArgs(opts *bind.CallOpts, destChainSelector uint64, feeToken common.Address, feeTokenAmount *big.Int, extraArgs []byte) (ProcessMessageArgs,

		error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	ValidatePoolReturnData(opts *bind.CallOpts, destChainSelector uint64, rampTokenAmounts []InternalRampTokenAmount, sourceTokenAmounts []ClientEVMTokenAmount) error

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	ApplyAuthorizedCallerUpdates(opts *bind.TransactOpts, authorizedCallerArgs AuthorizedCallersAuthorizedCallerArgs) (*types.Transaction, error)

	ApplyDestChainConfigUpdates(opts *bind.TransactOpts, destChainConfigArgs []PriceRegistryDestChainConfigArgs) (*types.Transaction, error)

	ApplyFeeTokensUpdates(opts *bind.TransactOpts, feeTokensToAdd []common.Address, feeTokensToRemove []common.Address) (*types.Transaction, error)

	ApplyPremiumMultiplierWeiPerEthUpdates(opts *bind.TransactOpts, premiumMultiplierWeiPerEthArgs []PriceRegistryPremiumMultiplierWeiPerEthArgs) (*types.Transaction, error)

	ApplyTokenTransferFeeConfigUpdates(opts *bind.TransactOpts, tokenTransferFeeConfigArgs []PriceRegistryTokenTransferFeeConfigArgs, tokensToUseDefaultFeeConfigs []PriceRegistryTokenTransferFeeConfigRemoveArgs) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	UpdatePrices(opts *bind.TransactOpts, priceUpdates InternalPriceUpdates) (*types.Transaction, error)

	UpdateTokenPriceFeeds(opts *bind.TransactOpts, tokenPriceFeedUpdates []PriceRegistryTokenPriceFeedUpdate) (*types.Transaction, error)

	FilterAuthorizedCallerAdded(opts *bind.FilterOpts) (*PriceRegistryAuthorizedCallerAddedIterator, error)

	WatchAuthorizedCallerAdded(opts *bind.WatchOpts, sink chan<- *PriceRegistryAuthorizedCallerAdded) (event.Subscription, error)

	ParseAuthorizedCallerAdded(log types.Log) (*PriceRegistryAuthorizedCallerAdded, error)

	FilterAuthorizedCallerRemoved(opts *bind.FilterOpts) (*PriceRegistryAuthorizedCallerRemovedIterator, error)

	WatchAuthorizedCallerRemoved(opts *bind.WatchOpts, sink chan<- *PriceRegistryAuthorizedCallerRemoved) (event.Subscription, error)

	ParseAuthorizedCallerRemoved(log types.Log) (*PriceRegistryAuthorizedCallerRemoved, error)

	FilterDestChainAdded(opts *bind.FilterOpts, destChainSelector []uint64) (*PriceRegistryDestChainAddedIterator, error)

	WatchDestChainAdded(opts *bind.WatchOpts, sink chan<- *PriceRegistryDestChainAdded, destChainSelector []uint64) (event.Subscription, error)

	ParseDestChainAdded(log types.Log) (*PriceRegistryDestChainAdded, error)

	FilterDestChainConfigUpdated(opts *bind.FilterOpts, destChainSelector []uint64) (*PriceRegistryDestChainConfigUpdatedIterator, error)

	WatchDestChainConfigUpdated(opts *bind.WatchOpts, sink chan<- *PriceRegistryDestChainConfigUpdated, destChainSelector []uint64) (event.Subscription, error)

	ParseDestChainConfigUpdated(log types.Log) (*PriceRegistryDestChainConfigUpdated, error)

	FilterFeeTokenAdded(opts *bind.FilterOpts, feeToken []common.Address) (*PriceRegistryFeeTokenAddedIterator, error)

	WatchFeeTokenAdded(opts *bind.WatchOpts, sink chan<- *PriceRegistryFeeTokenAdded, feeToken []common.Address) (event.Subscription, error)

	ParseFeeTokenAdded(log types.Log) (*PriceRegistryFeeTokenAdded, error)

	FilterFeeTokenRemoved(opts *bind.FilterOpts, feeToken []common.Address) (*PriceRegistryFeeTokenRemovedIterator, error)

	WatchFeeTokenRemoved(opts *bind.WatchOpts, sink chan<- *PriceRegistryFeeTokenRemoved, feeToken []common.Address) (event.Subscription, error)

	ParseFeeTokenRemoved(log types.Log) (*PriceRegistryFeeTokenRemoved, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*PriceRegistryOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *PriceRegistryOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*PriceRegistryOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*PriceRegistryOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *PriceRegistryOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*PriceRegistryOwnershipTransferred, error)

	FilterPremiumMultiplierWeiPerEthUpdated(opts *bind.FilterOpts, token []common.Address) (*PriceRegistryPremiumMultiplierWeiPerEthUpdatedIterator, error)

	WatchPremiumMultiplierWeiPerEthUpdated(opts *bind.WatchOpts, sink chan<- *PriceRegistryPremiumMultiplierWeiPerEthUpdated, token []common.Address) (event.Subscription, error)

	ParsePremiumMultiplierWeiPerEthUpdated(log types.Log) (*PriceRegistryPremiumMultiplierWeiPerEthUpdated, error)

	FilterPriceFeedPerTokenUpdated(opts *bind.FilterOpts, token []common.Address) (*PriceRegistryPriceFeedPerTokenUpdatedIterator, error)

	WatchPriceFeedPerTokenUpdated(opts *bind.WatchOpts, sink chan<- *PriceRegistryPriceFeedPerTokenUpdated, token []common.Address) (event.Subscription, error)

	ParsePriceFeedPerTokenUpdated(log types.Log) (*PriceRegistryPriceFeedPerTokenUpdated, error)

	FilterTokenTransferFeeConfigDeleted(opts *bind.FilterOpts, destChainSelector []uint64, token []common.Address) (*PriceRegistryTokenTransferFeeConfigDeletedIterator, error)

	WatchTokenTransferFeeConfigDeleted(opts *bind.WatchOpts, sink chan<- *PriceRegistryTokenTransferFeeConfigDeleted, destChainSelector []uint64, token []common.Address) (event.Subscription, error)

	ParseTokenTransferFeeConfigDeleted(log types.Log) (*PriceRegistryTokenTransferFeeConfigDeleted, error)

	FilterTokenTransferFeeConfigUpdated(opts *bind.FilterOpts, destChainSelector []uint64, token []common.Address) (*PriceRegistryTokenTransferFeeConfigUpdatedIterator, error)

	WatchTokenTransferFeeConfigUpdated(opts *bind.WatchOpts, sink chan<- *PriceRegistryTokenTransferFeeConfigUpdated, destChainSelector []uint64, token []common.Address) (event.Subscription, error)

	ParseTokenTransferFeeConfigUpdated(log types.Log) (*PriceRegistryTokenTransferFeeConfigUpdated, error)

	FilterUsdPerTokenUpdated(opts *bind.FilterOpts, token []common.Address) (*PriceRegistryUsdPerTokenUpdatedIterator, error)

	WatchUsdPerTokenUpdated(opts *bind.WatchOpts, sink chan<- *PriceRegistryUsdPerTokenUpdated, token []common.Address) (event.Subscription, error)

	ParseUsdPerTokenUpdated(log types.Log) (*PriceRegistryUsdPerTokenUpdated, error)

	FilterUsdPerUnitGasUpdated(opts *bind.FilterOpts, destChain []uint64) (*PriceRegistryUsdPerUnitGasUpdatedIterator, error)

	WatchUsdPerUnitGasUpdated(opts *bind.WatchOpts, sink chan<- *PriceRegistryUsdPerUnitGasUpdated, destChain []uint64) (event.Subscription, error)

	ParseUsdPerUnitGasUpdated(log types.Log) (*PriceRegistryUsdPerUnitGasUpdated, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
