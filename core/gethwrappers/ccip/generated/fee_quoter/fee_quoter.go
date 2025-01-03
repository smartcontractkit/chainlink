// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package fee_quoter

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

type FeeQuoterDestChainConfig struct {
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
	DefaultTxGasLimit                 uint32
	GasMultiplierWeiPerEth            uint64
	NetworkFeeUSDCents                uint32
	GasPriceStalenessThreshold        uint32
	EnforceOutOfOrder                 bool
	ChainFamilySelector               [4]byte
}

type FeeQuoterDestChainConfigArgs struct {
	DestChainSelector uint64
	DestChainConfig   FeeQuoterDestChainConfig
}

type FeeQuoterPremiumMultiplierWeiPerEthArgs struct {
	Token                      common.Address
	PremiumMultiplierWeiPerEth uint64
}

type FeeQuoterStaticConfig struct {
	MaxFeeJuelsPerMsg            *big.Int
	LinkToken                    common.Address
	TokenPriceStalenessThreshold uint32
}

type FeeQuoterTokenPriceFeedConfig struct {
	DataFeedAddress common.Address
	TokenDecimals   uint8
	IsEnabled       bool
}

type FeeQuoterTokenPriceFeedUpdate struct {
	SourceToken common.Address
	FeedConfig  FeeQuoterTokenPriceFeedConfig
}

type FeeQuoterTokenTransferFeeConfig struct {
	MinFeeUSDCents    uint32
	MaxFeeUSDCents    uint32
	DeciBps           uint16
	DestGasOverhead   uint32
	DestBytesOverhead uint32
	IsEnabled         bool
}

type FeeQuoterTokenTransferFeeConfigArgs struct {
	DestChainSelector       uint64
	TokenTransferFeeConfigs []FeeQuoterTokenTransferFeeConfigSingleTokenArgs
}

type FeeQuoterTokenTransferFeeConfigRemoveArgs struct {
	DestChainSelector uint64
	Token             common.Address
}

type FeeQuoterTokenTransferFeeConfigSingleTokenArgs struct {
	Token                  common.Address
	TokenTransferFeeConfig FeeQuoterTokenTransferFeeConfig
}

type InternalEVM2AnyTokenTransfer struct {
	SourcePoolAddress common.Address
	DestTokenAddress  []byte
	ExtraData         []byte
	Amount            *big.Int
	DestExecData      []byte
}

type InternalGasPriceUpdate struct {
	DestChainSelector uint64
	UsdPerUnitGas     *big.Int
}

type InternalPriceUpdates struct {
	TokenPriceUpdates []InternalTokenPriceUpdate
	GasPriceUpdates   []InternalGasPriceUpdate
}

type InternalTimestampedPackedUint224 struct {
	Value     *big.Int
	Timestamp uint32
}

type InternalTokenPriceUpdate struct {
	SourceToken common.Address
	UsdPerToken *big.Int
}

type KeystoneFeedsPermissionHandlerPermission struct {
	Forwarder     common.Address
	WorkflowName  [10]byte
	ReportName    [2]byte
	WorkflowOwner common.Address
	IsAllowed     bool
}

var FeeQuoterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"uint96\",\"name\":\"maxFeeJuelsPerMsg\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"linkToken\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"tokenPriceStalenessThreshold\",\"type\":\"uint32\"}],\"internalType\":\"structFeeQuoter.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"},{\"internalType\":\"address[]\",\"name\":\"priceUpdaters\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"feeTokens\",\"type\":\"address[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"dataFeedAddress\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"tokenDecimals\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"}],\"internalType\":\"structFeeQuoter.TokenPriceFeedConfig\",\"name\":\"feedConfig\",\"type\":\"tuple\"}],\"internalType\":\"structFeeQuoter.TokenPriceFeedUpdate[]\",\"name\":\"tokenPriceFeeds\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"minFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"deciBps\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"destBytesOverhead\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"}],\"internalType\":\"structFeeQuoter.TokenTransferFeeConfig\",\"name\":\"tokenTransferFeeConfig\",\"type\":\"tuple\"}],\"internalType\":\"structFeeQuoter.TokenTransferFeeConfigSingleTokenArgs[]\",\"name\":\"tokenTransferFeeConfigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structFeeQuoter.TokenTransferFeeConfigArgs[]\",\"name\":\"tokenTransferFeeConfigArgs\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"premiumMultiplierWeiPerEth\",\"type\":\"uint64\"}],\"internalType\":\"structFeeQuoter.PremiumMultiplierWeiPerEthArgs[]\",\"name\":\"premiumMultiplierWeiPerEthArgs\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint16\",\"name\":\"maxNumberOfTokensPerMsg\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxDataBytes\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerMsgGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerPayloadByte\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destDataAvailabilityOverheadGas\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerDataAvailabilityByte\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"destDataAvailabilityMultiplierBps\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"defaultTokenFeeUSDCents\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"defaultTokenDestGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"defaultTxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"gasMultiplierWeiPerEth\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"networkFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasPriceStalenessThreshold\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"enforceOutOfOrder\",\"type\":\"bool\"},{\"internalType\":\"bytes4\",\"name\":\"chainFamilySelector\",\"type\":\"bytes4\"}],\"internalType\":\"structFeeQuoter.DestChainConfig\",\"name\":\"destChainConfig\",\"type\":\"tuple\"}],\"internalType\":\"structFeeQuoter.DestChainConfigArgs[]\",\"name\":\"destChainConfigArgs\",\"type\":\"tuple[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"CannotTransferToSelf\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DataFeedValueOutOfUint224Range\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"}],\"name\":\"DestinationChainNotEnabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ExtraArgOutOfOrderExecutionMustBeTrue\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"FeeTokenNotSupported\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"destBytesOverhead\",\"type\":\"uint32\"}],\"name\":\"InvalidDestBytesOverhead\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"}],\"name\":\"InvalidDestChainConfig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedAddress\",\"type\":\"bytes\"}],\"name\":\"InvalidEVMAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidExtraArgsTag\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minFeeUSDCents\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeUSDCents\",\"type\":\"uint256\"}],\"name\":\"InvalidFeeRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidStaticConfig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"msgFeeJuels\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeJuelsPerMsg\",\"type\":\"uint256\"}],\"name\":\"MessageFeeTooHigh\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MessageGasLimitTooHigh\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"maxSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actualSize\",\"type\":\"uint256\"}],\"name\":\"MessageTooLarge\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeProposedOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OwnerCannotBeZero\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"workflowOwner\",\"type\":\"address\"},{\"internalType\":\"bytes10\",\"name\":\"workflowName\",\"type\":\"bytes10\"},{\"internalType\":\"bytes2\",\"name\":\"reportName\",\"type\":\"bytes2\"}],\"name\":\"ReportForwarderUnauthorized\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SourceTokenDataTooLarge\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"threshold\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timePassed\",\"type\":\"uint256\"}],\"name\":\"StaleGasPrice\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"TokenNotSupported\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"UnauthorizedCaller\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numberOfTokens\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxNumberOfTokensPerMsg\",\"type\":\"uint256\"}],\"name\":\"UnsupportedNumberOfTokens\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"AuthorizedCallerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"AuthorizedCallerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint16\",\"name\":\"maxNumberOfTokensPerMsg\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxDataBytes\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerMsgGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerPayloadByte\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destDataAvailabilityOverheadGas\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerDataAvailabilityByte\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"destDataAvailabilityMultiplierBps\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"defaultTokenFeeUSDCents\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"defaultTokenDestGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"defaultTxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"gasMultiplierWeiPerEth\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"networkFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasPriceStalenessThreshold\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"enforceOutOfOrder\",\"type\":\"bool\"},{\"internalType\":\"bytes4\",\"name\":\"chainFamilySelector\",\"type\":\"bytes4\"}],\"indexed\":false,\"internalType\":\"structFeeQuoter.DestChainConfig\",\"name\":\"destChainConfig\",\"type\":\"tuple\"}],\"name\":\"DestChainAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint16\",\"name\":\"maxNumberOfTokensPerMsg\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxDataBytes\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerMsgGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerPayloadByte\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destDataAvailabilityOverheadGas\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerDataAvailabilityByte\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"destDataAvailabilityMultiplierBps\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"defaultTokenFeeUSDCents\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"defaultTokenDestGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"defaultTxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"gasMultiplierWeiPerEth\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"networkFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasPriceStalenessThreshold\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"enforceOutOfOrder\",\"type\":\"bool\"},{\"internalType\":\"bytes4\",\"name\":\"chainFamilySelector\",\"type\":\"bytes4\"}],\"indexed\":false,\"internalType\":\"structFeeQuoter.DestChainConfig\",\"name\":\"destChainConfig\",\"type\":\"tuple\"}],\"name\":\"DestChainConfigUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"}],\"name\":\"FeeTokenAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"}],\"name\":\"FeeTokenRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"premiumMultiplierWeiPerEth\",\"type\":\"uint64\"}],\"name\":\"PremiumMultiplierWeiPerEthUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"dataFeedAddress\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"tokenDecimals\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"}],\"indexed\":false,\"internalType\":\"structFeeQuoter.TokenPriceFeedConfig\",\"name\":\"priceFeedConfig\",\"type\":\"tuple\"}],\"name\":\"PriceFeedPerTokenUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"reportId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"},{\"internalType\":\"bytes10\",\"name\":\"workflowName\",\"type\":\"bytes10\"},{\"internalType\":\"bytes2\",\"name\":\"reportName\",\"type\":\"bytes2\"},{\"internalType\":\"address\",\"name\":\"workflowOwner\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isAllowed\",\"type\":\"bool\"}],\"indexed\":false,\"internalType\":\"structKeystoneFeedsPermissionHandler.Permission\",\"name\":\"permission\",\"type\":\"tuple\"}],\"name\":\"ReportPermissionSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"TokenTransferFeeConfigDeleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"minFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"deciBps\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"destBytesOverhead\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"}],\"indexed\":false,\"internalType\":\"structFeeQuoter.TokenTransferFeeConfig\",\"name\":\"tokenTransferFeeConfig\",\"type\":\"tuple\"}],\"name\":\"TokenTransferFeeConfigUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"UsdPerTokenUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"destChain\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"UsdPerUnitGasUpdated\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"FEE_BASE_DECIMALS\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"KEYSTONE_PRICE_DECIMALS\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address[]\",\"name\":\"addedCallers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"removedCallers\",\"type\":\"address[]\"}],\"internalType\":\"structAuthorizedCallers.AuthorizedCallerArgs\",\"name\":\"authorizedCallerArgs\",\"type\":\"tuple\"}],\"name\":\"applyAuthorizedCallerUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint16\",\"name\":\"maxNumberOfTokensPerMsg\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxDataBytes\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerMsgGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerPayloadByte\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destDataAvailabilityOverheadGas\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerDataAvailabilityByte\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"destDataAvailabilityMultiplierBps\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"defaultTokenFeeUSDCents\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"defaultTokenDestGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"defaultTxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"gasMultiplierWeiPerEth\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"networkFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasPriceStalenessThreshold\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"enforceOutOfOrder\",\"type\":\"bool\"},{\"internalType\":\"bytes4\",\"name\":\"chainFamilySelector\",\"type\":\"bytes4\"}],\"internalType\":\"structFeeQuoter.DestChainConfig\",\"name\":\"destChainConfig\",\"type\":\"tuple\"}],\"internalType\":\"structFeeQuoter.DestChainConfigArgs[]\",\"name\":\"destChainConfigArgs\",\"type\":\"tuple[]\"}],\"name\":\"applyDestChainConfigUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"feeTokensToRemove\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"feeTokensToAdd\",\"type\":\"address[]\"}],\"name\":\"applyFeeTokensUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"premiumMultiplierWeiPerEth\",\"type\":\"uint64\"}],\"internalType\":\"structFeeQuoter.PremiumMultiplierWeiPerEthArgs[]\",\"name\":\"premiumMultiplierWeiPerEthArgs\",\"type\":\"tuple[]\"}],\"name\":\"applyPremiumMultiplierWeiPerEthUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"minFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"deciBps\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"destBytesOverhead\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"}],\"internalType\":\"structFeeQuoter.TokenTransferFeeConfig\",\"name\":\"tokenTransferFeeConfig\",\"type\":\"tuple\"}],\"internalType\":\"structFeeQuoter.TokenTransferFeeConfigSingleTokenArgs[]\",\"name\":\"tokenTransferFeeConfigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structFeeQuoter.TokenTransferFeeConfigArgs[]\",\"name\":\"tokenTransferFeeConfigArgs\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"internalType\":\"structFeeQuoter.TokenTransferFeeConfigRemoveArgs[]\",\"name\":\"tokensToUseDefaultFeeConfigs\",\"type\":\"tuple[]\"}],\"name\":\"applyTokenTransferFeeConfigUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"fromToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"fromTokenAmount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"toToken\",\"type\":\"address\"}],\"name\":\"convertTokenAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllAuthorizedCallers\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"}],\"name\":\"getDestChainConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint16\",\"name\":\"maxNumberOfTokensPerMsg\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxDataBytes\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerMsgGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerPayloadByte\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destDataAvailabilityOverheadGas\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerDataAvailabilityByte\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"destDataAvailabilityMultiplierBps\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"defaultTokenFeeUSDCents\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"defaultTokenDestGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"defaultTxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"gasMultiplierWeiPerEth\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"networkFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasPriceStalenessThreshold\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"enforceOutOfOrder\",\"type\":\"bool\"},{\"internalType\":\"bytes4\",\"name\":\"chainFamilySelector\",\"type\":\"bytes4\"}],\"internalType\":\"structFeeQuoter.DestChainConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"}],\"name\":\"getDestinationChainGasPrice\",\"outputs\":[{\"components\":[{\"internalType\":\"uint224\",\"name\":\"value\",\"type\":\"uint224\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"}],\"internalType\":\"structInternal.TimestampedPackedUint224\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFeeTokens\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getPremiumMultiplierWeiPerEth\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"premiumMultiplierWeiPerEth\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStaticConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint96\",\"name\":\"maxFeeJuelsPerMsg\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"linkToken\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"tokenPriceStalenessThreshold\",\"type\":\"uint32\"}],\"internalType\":\"structFeeQuoter.StaticConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"}],\"name\":\"getTokenAndGasPrices\",\"outputs\":[{\"internalType\":\"uint224\",\"name\":\"tokenPrice\",\"type\":\"uint224\"},{\"internalType\":\"uint224\",\"name\":\"gasPriceValue\",\"type\":\"uint224\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getTokenPrice\",\"outputs\":[{\"components\":[{\"internalType\":\"uint224\",\"name\":\"value\",\"type\":\"uint224\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"}],\"internalType\":\"structInternal.TimestampedPackedUint224\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getTokenPriceFeedConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"dataFeedAddress\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"tokenDecimals\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"}],\"internalType\":\"structFeeQuoter.TokenPriceFeedConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"tokens\",\"type\":\"address[]\"}],\"name\":\"getTokenPrices\",\"outputs\":[{\"components\":[{\"internalType\":\"uint224\",\"name\":\"value\",\"type\":\"uint224\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"}],\"internalType\":\"structInternal.TimestampedPackedUint224[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getTokenTransferFeeConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"minFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"deciBps\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"destBytesOverhead\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"}],\"internalType\":\"structFeeQuoter.TokenTransferFeeConfig\",\"name\":\"tokenTransferFeeConfig\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMTokenAmount[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structClient.EVM2AnyMessage\",\"name\":\"message\",\"type\":\"tuple\"}],\"name\":\"getValidatedFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"feeTokenAmount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getValidatedTokenPrice\",\"outputs\":[{\"internalType\":\"uint224\",\"name\":\"\",\"type\":\"uint224\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"}],\"name\":\"onReport\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"feeTokenAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"sourcePoolAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"destExecData\",\"type\":\"bytes\"}],\"internalType\":\"structInternal.EVM2AnyTokenTransfer[]\",\"name\":\"onRampTokenTransfers\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMTokenAmount[]\",\"name\":\"sourceTokenAmounts\",\"type\":\"tuple[]\"}],\"name\":\"processMessageArgs\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"msgFeeJuels\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isOutOfOrderExecution\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"convertedExtraArgs\",\"type\":\"bytes\"},{\"internalType\":\"bytes[]\",\"name\":\"destExecDataPerToken\",\"type\":\"bytes[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"},{\"internalType\":\"bytes10\",\"name\":\"workflowName\",\"type\":\"bytes10\"},{\"internalType\":\"bytes2\",\"name\":\"reportName\",\"type\":\"bytes2\"},{\"internalType\":\"address\",\"name\":\"workflowOwner\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isAllowed\",\"type\":\"bool\"}],\"internalType\":\"structKeystoneFeedsPermissionHandler.Permission[]\",\"name\":\"permissions\",\"type\":\"tuple[]\"}],\"name\":\"setReportPermissions\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"internalType\":\"uint224\",\"name\":\"usdPerToken\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.TokenPriceUpdate[]\",\"name\":\"tokenPriceUpdates\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint224\",\"name\":\"usdPerUnitGas\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.GasPriceUpdate[]\",\"name\":\"gasPriceUpdates\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.PriceUpdates\",\"name\":\"priceUpdates\",\"type\":\"tuple\"}],\"name\":\"updatePrices\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"dataFeedAddress\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"tokenDecimals\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"}],\"internalType\":\"structFeeQuoter.TokenPriceFeedConfig\",\"name\":\"feedConfig\",\"type\":\"tuple\"}],\"internalType\":\"structFeeQuoter.TokenPriceFeedUpdate[]\",\"name\":\"tokenPriceFeedUpdates\",\"type\":\"tuple[]\"}],\"name\":\"updateTokenPriceFeeds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60e060405234610fed57617375803803806100198161125f565b9283398101908082036101208112610fed57606013610fed5761003a611221565b81516001600160601b0381168103610fed57815261005a60208301611284565b906020810191825261006e60408401611298565b6040820190815260608401516001600160401b038111610fed57856100949186016112c0565b60808501519094906001600160401b038111610fed57866100b69183016112c0565b60a08201519096906001600160401b038111610fed5782019181601f84011215610fed5782516100ed6100e8826112a9565b61125f565b9360208086848152019260071b82010190848211610fed57602001915b8183106111a65750505060c08101516001600160401b038111610fed5781019782601f8a011215610fed578851986101446100e88b6112a9565b996020808c838152019160051b83010191858311610fed5760208101915b838310611044575050505060e08201516001600160401b038111610fed5782019183601f84011215610fed57825161019c6100e8826112a9565b9360208086848152019260061b82010190868211610fed57602001915b81831061100857505050610100810151906001600160401b038211610fed570183601f82011215610fed578051906101f36100e8836112a9565b946020610240818886815201940283010191818311610fed57602001925b828410610e5157505050503315610e4057600180546001600160a01b031916331790556020986102408a61125f565b97600089526000368137610252611240565b998a52888b8b015260005b89518110156102c4576001906001600160a01b0361027b828d61134b565b51168d61028782611537565b610294575b50500161025d565b7fc3803387881faad271c47728894e3e36fac830ffc8602ca6fc07733cbda7758091604051908152a1388d61028c565b508a985089519660005b885181101561033f576001600160a01b036102e9828b61134b565b511690811561032e577feb1b9b92e50b7f88f9ff25d56765095ac6e91540eee214906f4036a908ffbdef8c836103206001956114bf565b50604051908152a1016102ce565b6342bcdf7f60e11b60005260046000fd5b5089959697508860018060a01b03835116158015610e2e575b8015610e1f575b610e0e5791516001600160a01b031660a05290516001600160601b03166080525163ffffffff1660c0526103928461125f565b9360008552600036813760005b855181101561040e576001906103c76001600160a01b036103c0838a61134b565b51166113cc565b6103d2575b0161039f565b818060a01b036103e2828961134b565b51167f1795838dc8ab2ffc5f431a1729a6afa0b587f982f7b2be0b9d7187a1ef547f91600080a26103cc565b508694508560005b84518110156104855760019061043e6001600160a01b03610437838961134b565b51166114fe565b610449575b01610416565b818060a01b03610459828861134b565b51167fdf1b1bd32a69711488d71554706bb130b1fc63a5fa1a2cd85e8440f84065ba23600080a2610443565b508593508460005b855181101561054757806104a36001928861134b565b517fe6a7a17d710bf0b2cd05e5397dc6f97a5da4ee79e31e234bf5f965ee2bd9a5bf606087858060a01b03845116930151836000526007895260406000209060ff878060a01b038251169283898060a01b03198254161781558b8301908151604082549501948460a81b8651151560a81b16918560a01b9060a01b169061ffff60a01b19161717905560405193845251168a8301525115156040820152a20161048d565b50828460005b8451811015610ab957610560818661134b565b51836001600160401b03610574848961134b565b5151169101519080158015610aa6575b8015610a86575b8015610a68575b610a54579081600193926000526009865263ffffffff60e01b846040600020015460381b161560001461090957807f525e3d4e0c31cef19cf9426af8d2c0ddd2d576359ca26bed92aac5fadda4626560405180610702868291909161020061022082019380511515835261ffff602082015116602084015263ffffffff604082015116604084015263ffffffff606082015116606084015263ffffffff608082015116608084015261ffff60a08201511660a084015263ffffffff60c08201511660c084015261ffff60e08201511660e084015261ffff6101008201511661010084015261ffff6101208201511661012084015263ffffffff6101408201511661014084015263ffffffff6101608201511661016084015260018060401b036101808201511661018084015263ffffffff6101a0820151166101a084015263ffffffff6101c0820151166101c08401526101e081015115156101e08401528163ffffffff60e01b91015116910152565b0390a25b600052600985528260406000208251151581549062ffff008986015160081b169166ffffffff000000604087015160181b166affffffff00000000000000606088015160381b16906effffffff0000000000000000000000608089015160581b169361ffff60781b60a08a015160781b169163ffffffff60881b60c08b015160881b169061ffff60a81b60e08c015160a81b169361ffff60b81b6101008d015160b81b16968c61012061ffff60c81b91015160c81b16988d61014063ffffffff60d81b91015160d81b169a63ffffffff60d81b199861ffff60c81b199761ffff60b81b199661ffff60a81b199563ffffffff60881b199461ffff60781b19936effffffff0000000000000000000000199260ff6affffffff000000000000001992169066ffffffffffffff1916171617161716171617161716171617161717178155019063ffffffff80610160830151161663ffffffff198354161782556101808101519082546fffffffff0000000000000000000000006101a083015160601b169063ffffffff60801b6101c084015160801b16916101e0840151151560a01b9363ffffffff60a81b9061020063ffffffff60a81b91015160381c16169463ffffffff60a81b199260ff60a01b19916401000000008b60601b03908d1b16906401000000008b60a01b03191617161716179060ff60a01b16171790550161054d565b807f283b699f411baff8f1c29fe49f32a828c8151596244b8e7e4c164edd6569a83560405180610a4c868291909161020061022082019380511515835261ffff602082015116602084015263ffffffff604082015116604084015263ffffffff606082015116606084015263ffffffff608082015116608084015261ffff60a08201511660a084015263ffffffff60c08201511660c084015261ffff60e08201511660e084015261ffff6101008201511661010084015261ffff6101208201511661012084015263ffffffff6101408201511661014084015263ffffffff6101608201511661016084015260018060401b036101808201511661018084015263ffffffff6101a0820151166101a084015263ffffffff6101c0820151166101c08401526101e081015115156101e08401528163ffffffff60e01b91015116910152565b0390a2610706565b63c35aa79d60e01b60005260045260246000fd5b5063ffffffff6101608301511663ffffffff60608401511610610592565b506102008201516001600160e01b031916630a04b54b60e21b141561058b565b5063ffffffff6101608301511615610584565b508260005b8151811015610b3e576001906001600160a01b03610adc828561134b565b5151167fbb77da6f7210cdd16904228a9360133d1d7dfff99b1bc75f128da5b53e28f97d86848060401b0381610b12868961134b565b510151168360005260088252604060002081878060401b0319825416179055604051908152a201610abe565b50506001610b4b8361125f565b9160008352600091610e09575b81925b8151841015610d4757610b6e848361134b565b5180516001600160401b0316929086019190845b8784518051831015610d365782610b989161134b565b51015184516001600160a01b0390610bb190849061134b565b5151169063ffffffff8151168a82019063ffffffff8251169081811015610d21575050608082019063ffffffff8251168c8110610d0a57507f94967ae9ea7729ad4f54021c1981765d2b1d954f7c92fbec340aa0a54f46b8b560c08d8098979563ffffffff600198968f8f9761ffff6040838b879552600a89528181208f808f9160a01b0316825289522091838751169784891685198554161784558482519285549560408b019469ffff0000000000000000865160401b169b60608101986dffffffff000000000000000000008a5160501b169060a08d8860701b905160701b1693019e8f60ff60901b9051151560901b169460ff60901b199267ffffffff000000008a60701b19928c1b169060016401000000009060701b0319161716171617171790556040519a8b5251169089015251166040870152511660608501525116608083015251151560a0820152a301909150610b82565b6312766e0160e11b8b52600485905260245260448afd5b6305a7b3d160e11b8b52600452602452604489fd5b505050925093600191500192610b5b565b905083825b8251811015610dca576001906001600160401b03610d6a828661134b565b515116828060a01b0384610d7e848861134b565b5101511690808752600a855260408720848060a01b038316885285528660408120557f4de5b1bcbca6018c11303a2c3f4a4b4f22a1c741d8c4ba430d246ac06c5ddf8b8780a301610d4c565b604051615da990816115cc82396080518181816104d4015261312d015260a05181818161051701526130c4015260c05181818161053e0152613ddc0152f35b610b58565b63d794ef9560e01b60005260046000fd5b5063ffffffff8251161561035f565b5080516001600160601b031615610358565b639b15e16f60e01b60005260046000fd5b8382036102408112610fed57610220610e68611240565b91610e7287611328565b8352601f190112610fed576040519161022083016001600160401b03811184821017610ff257604052610ea76020870161131b565b8352610eb56040870161133c565b6020840152610ec660608701611298565b6040840152610ed760808701611298565b6060840152610ee860a08701611298565b6080840152610ef960c0870161133c565b60a0840152610f0a60e08701611298565b60c0840152610f1c610100870161133c565b60e0840152610f2e610120870161133c565b610100840152610f41610140870161133c565b610120840152610f546101608701611298565b610140840152610f676101808701611298565b610160840152610f7a6101a08701611328565b610180840152610f8d6101c08701611298565b6101a0840152610fa06101e08701611298565b6101c0840152610fb3610200870161131b565b6101e0840152610220860151916001600160e01b031983168303610fed578360209361020061024096015283820152815201930192610211565b600080fd5b634e487b7160e01b600052604160045260246000fd5b604083880312610fed57602060409161101f611240565b61102886611284565b8152611035838701611328565b838201528152019201916101b9565b82516001600160401b038111610fed5782016040818903601f190112610fed5761106c611240565b9061107960208201611328565b825260408101516001600160401b038111610fed57602091010188601f82011215610fed5780516110ac6100e8826112a9565b91602060e08185858152019302820101908b8211610fed57602001915b8183106110e85750505091816020938480940152815201920191610162565b828c0360e08112610fed5760c06110fd611240565b9161110786611284565b8352601f190112610fed576040519160c08301916001600160401b03831184841017610ff25760e093602093604052611141848801611298565b815261114f60408801611298565b8482015261115f6060880161133c565b604082015261117060808801611298565b606082015261118160a08801611298565b608082015261119260c0880161131b565b60a0820152838201528152019201916110c9565b82850360808112610fed5760606111bb611240565b916111c586611284565b8352601f190112610fed576111d8611221565b916111e560208601611284565b835260408501519160ff83168303610fed578360209384608096015261120d6060880161131b565b60408201528382015281520192019161010a565b60405190606082016001600160401b03811183821017610ff257604052565b60408051919082016001600160401b03811183821017610ff257604052565b6040519190601f01601f191682016001600160401b03811183821017610ff257604052565b51906001600160a01b0382168203610fed57565b519063ffffffff82168203610fed57565b6001600160401b038111610ff25760051b60200190565b9080601f83011215610fed5781516112da6100e8826112a9565b9260208085848152019260051b820101928311610fed57602001905b8282106113035750505090565b6020809161131084611284565b8152019101906112f6565b51908115158203610fed57565b51906001600160401b0382168203610fed57565b519061ffff82168203610fed57565b805182101561135f5760209160051b010190565b634e487b7160e01b600052603260045260246000fd5b805482101561135f5760005260206000200190600090565b805480156113b65760001901906113a48282611375565b8154906000199060031b1b1916905555565b634e487b7160e01b600052603160045260246000fd5b6000818152600c6020526040902054801561148d57600019810181811161147757600b5460001981019190821161147757818103611426575b505050611412600b61138d565b600052600c60205260006040812055600190565b61145f61143761144893600b611375565b90549060031b1c928392600b611375565b819391549060031b91821b91600019901b19161790565b9055600052600c602052604060002055388080611405565b634e487b7160e01b600052601160045260246000fd5b5050600090565b80549068010000000000000000821015610ff257816114489160016114bb94018155611375565b9055565b806000526003602052604060002054156000146114f8576114e1816002611494565b600254906000526003602052604060002055600190565b50600090565b80600052600c602052604060002054156000146114f85761152081600b611494565b600b5490600052600c602052604060002055600190565b600081815260036020526040902054801561148d5760001981018181116114775760025460001981019190821161147757808203611591575b50505061157d600261138d565b600052600360205260006040812055600190565b6115b36115a2611448936002611375565b90549060031b1c9283926002611375565b9055600052600360205260406000205538808061157056fe6080604052600436101561001257600080fd5b60003560e01c806241e5be1461020657806301ffc9a714610201578063061877e3146101fc57806306285c69146101f7578063181f5a77146101f25780632451a627146101ed578063325c868e146101e85780633937306f146101e357806341ed29e7146101de578063430d138c146101d957806345ac924d146101d45780634ab35b0b146101cf578063514e8cff146101ca5780636cb5f3dd146101c55780636def4ce7146101c0578063770e2dc4146101bb57806379ba5097146101b65780637afac322146101b1578063805f2132146101ac57806382b49eb0146101a75780638da5cb5b146101a257806391a2749a1461019d578063a69c64c014610198578063bf78e03f14610193578063cdc73d511461018e578063d02641a014610189578063d63d3af214610184578063d8694ccd1461017f578063f2fde38b1461017a578063fbe3f778146101755763ffdb4b371461017057600080fd5b612bf2565b612ab9565b6129ba565b61258d565b612553565b6124d7565b612442565b612357565b612280565b6121b0565b61215e565b611fcc565b611c8b565b611b1a565b6119ca565b61176f565b611600565b61129a565b6111ab565b611140565b61103b565b610edb565b610c44565b610885565b61084b565b6107aa565b6106d9565b61047a565b610407565b6102c5565b61023b565b73ffffffffffffffffffffffffffffffffffffffff81160361022957565b600080fd5b35906102398261020b565b565b346102295760607ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261022957602061029060043561027b8161020b565b6024356044359161028b8361020b565b612da7565b604051908152f35b35907fffffffff000000000000000000000000000000000000000000000000000000008216820361022957565b346102295760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610229576004357fffffffff0000000000000000000000000000000000000000000000000000000081168103610229577fffffffff00000000000000000000000000000000000000000000000000000000602091167f805f21320000000000000000000000000000000000000000000000000000000081149081156103dd575b81156103b3575b8115610389575b506040519015158152f35b7f01ffc9a7000000000000000000000000000000000000000000000000000000009150143861037e565b7f181f5a770000000000000000000000000000000000000000000000000000000081149150610377565b7f9b645f410000000000000000000000000000000000000000000000000000000081149150610370565b346102295760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102295773ffffffffffffffffffffffffffffffffffffffff6004356104578161020b565b166000526008602052602067ffffffffffffffff60406000205416604051908152f35b346102295760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610229576104b1612df0565b5060606040516104c0816105a5565b63ffffffff6bffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169182815273ffffffffffffffffffffffffffffffffffffffff60406020830192827f00000000000000000000000000000000000000000000000000000000000000001684520191837f00000000000000000000000000000000000000000000000000000000000000001683526040519485525116602084015251166040820152f35b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6060810190811067ffffffffffffffff8211176105c157604052565b610576565b60a0810190811067ffffffffffffffff8211176105c157604052565b6040810190811067ffffffffffffffff8211176105c157604052565b60c0810190811067ffffffffffffffff8211176105c157604052565b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff8211176105c157604052565b6040519061023960408361061a565b604051906102396102208361061a565b919082519283825260005b8481106106c45750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8460006020809697860101520116010190565b80602080928401015182828601015201610685565b346102295760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261022957610756604080519061071a818361061a565b601382527f46656551756f74657220312e362e302d6465760000000000000000000000000060208301525191829160208352602083019061067a565b0390f35b602060408183019282815284518094520192019060005b81811061077e5750505090565b825173ffffffffffffffffffffffffffffffffffffffff16845260209384019390920191600101610771565b346102295760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102295760405180602060025491828152019060026000527f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace9060005b81811061083557610756856108298187038261061a565b6040519182918261075a565b8254845260209093019260019283019201610812565b346102295760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261022957602060405160248152f35b346102295760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102295760043567ffffffffffffffff811161022957806004019060407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc8236030112610229576108ff614180565b6109098280612e0f565b4263ffffffff1692915060005b818110610ae35750506024019061092d8284612e0f565b92905060005b83811061093c57005b8061095b610956600193610950868a612e0f565b90612e92565b612f07565b7fdd84a3fa9ef9409f550d54d6affec7e9c480c878c6ab27b78912a03e1b371c6e67ffffffffffffffff610aaa610a876020850194610a796109b987517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1690565b6109e86109c461065b565b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff9092168252565b63ffffffff8c166020820152610a23610a09845167ffffffffffffffff1690565b67ffffffffffffffff166000526005602052604060002090565b815160209092015160e01b7fffffffff00000000000000000000000000000000000000000000000000000000167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff92909216919091179055565b5167ffffffffffffffff1690565b93517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1690565b604080517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff9290921682524260208301529190931692a201610933565b80610afc610af76001936109508980612e0f565b612ed0565b7f52f50aa6d1a95a4595361ecf953d095f125d442e4673716dede699e049de148a73ffffffffffffffffffffffffffffffffffffffff610bde610a876020850194610bc4610b6687517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1690565b610b716109c461065b565b63ffffffff8d166020820152610a23610b9e845173ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff166000526006602052604060002090565b5173ffffffffffffffffffffffffffffffffffffffff1690565b604080517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff9290921682524260208301529190931692a201610916565b67ffffffffffffffff81116105c15760051b60200190565b8015150361022957565b359061023982610c2f565b346102295760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102295760043567ffffffffffffffff8111610229573660238201121561022957806004013590610c9f82610c17565b90610cad604051928361061a565b828252602460a06020840194028201019036821161022957602401925b818410610cdc57610cda83612f2c565b005b60a0843603126102295760405190610cf3826105c6565b8435610cfe8161020b565b825260208501357fffffffffffffffffffff00000000000000000000000000000000000000000000811681036102295760208301526040850135907fffff000000000000000000000000000000000000000000000000000000000000821682036102295782602092604060a0950152610d796060880161022e565b6060820152610d8a60808801610c39565b6080820152815201930192610cca565b6004359067ffffffffffffffff8216820361022957565b6024359067ffffffffffffffff8216820361022957565b359067ffffffffffffffff8216820361022957565b9181601f840112156102295782359167ffffffffffffffff8311610229576020838186019501011161022957565b9181601f840112156102295782359167ffffffffffffffff8311610229576020808501948460051b01011161022957565b929091610e5d9284521515602084015260806040840152608083019061067a565b906060818303910152815180825260208201916020808360051b8301019401926000915b838310610e9057505050505090565b9091929394602080610ecc837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08660019603018752895161067a565b97019301930191939290610e81565b346102295760c07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261022957610f12610d9a565b60243590610f1f8261020b565b60443560643567ffffffffffffffff811161022957610f42903690600401610ddd565b60849391933567ffffffffffffffff811161022957610f65903690600401610e0b565b9160a4359567ffffffffffffffff871161022957366023880112156102295786600401359567ffffffffffffffff8711610229573660248860061b8a01011161022957610756986024610fb99901966130b8565b9060409492945194859485610e3c565b602060408183019282815284518094520192019060005b818110610fed5750505090565b9091926020604082611030600194885163ffffffff602080927bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8151168552015116910152565b019401929101610fe0565b346102295760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102295760043567ffffffffffffffff81116102295761108a903690600401610e0b565b61109381610c17565b916110a1604051938461061a565b8183527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe06110ce83610c17565b0160005b81811061112957505060005b8281101561111b576001906110ff6110fa8260051b8501613208565b613d7b565b61110982876130a4565b5261111481866130a4565b50016110de565b604051806107568682610fc9565b6020906111346131ef565b828288010152016110d2565b346102295760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102295760206111856004356111808161020b565b6140d1565b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff60405191168152f35b346102295760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102295767ffffffffffffffff6111eb610d9a565b6111f36131ef565b5016600052600560205260406000206040519061120f826105e2565b547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8116825260e01c6020820152604051809161075682604081019263ffffffff602080927bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8151168552015116910152565b359061ffff8216820361022957565b359063ffffffff8216820361022957565b346102295760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102295760043567ffffffffffffffff81116102295736602382011215610229578060040135906112f582610c17565b90611303604051928361061a565b82825260246102406020840194028201019036821161022957602401925b81841061133157610cda8361324c565b8336036102408112610229576102207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe06040519261136e846105e2565b61137788610dc8565b84520112610229576102409160209161138e61066a565b611399848901610c39565b81526113a76040890161127a565b848201526113b760608901611289565b60408201526113c860808901611289565b60608201526113d960a08901611289565b60808201526113ea60c0890161127a565b60a08201526113fb60e08901611289565b60c082015261140d610100890161127a565b60e082015261141f610120890161127a565b610100820152611432610140890161127a565b6101208201526114456101608901611289565b6101408201526114586101808901611289565b61016082015261146b6101a08901610dc8565b61018082015261147e6101c08901611289565b6101a08201526114916101e08901611289565b6101c08201526114a46102008901610c39565b6101e08201526114b76102208901610298565b61020082015283820152815201930192611321565b61023990929192610200806102208301956114e984825115159052565b60208181015161ffff169085015260408181015163ffffffff169085015260608181015163ffffffff169085015260808181015163ffffffff169085015260a08181015161ffff169085015260c08181015163ffffffff169085015260e08181015161ffff16908501526101008181015161ffff16908501526101208181015161ffff16908501526101408181015163ffffffff16908501526101608181015163ffffffff16908501526101808181015167ffffffffffffffff16908501526101a08181015163ffffffff16908501526101c08181015163ffffffff16908501526101e08181015115159085015201517fffffffff0000000000000000000000000000000000000000000000000000000016910152565b346102295760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610229576107566116d36116ce611640610d9a565b600061020061164d61066a565b8281528260208201528260408201528260608201528260808201528260a08201528260c08201528260e08201528261010082015282610120820152826101408201528261016082015282610180820152826101a0820152826101c0820152826101e0820152015267ffffffffffffffff166000526009602052604060002090565b6134bc565b604051918291826114cc565b81601f82011215610229578035906116f682610c17565b92611704604051948561061a565b82845260208085019360061b8301019181831161022957602001925b82841061172e575050505090565b6040848303126102295760206040918251611748816105e2565b61175187610dc8565b8152828701356117608161020b565b83820152815201930192611720565b346102295760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102295760043567ffffffffffffffff811161022957366023820112156102295780600401356117c981610c17565b916117d7604051938461061a565b8183526024602084019260051b820101903682116102295760248101925b828410611826576024358567ffffffffffffffff821161022957611820610cda9236906004016116df565b90613638565b833567ffffffffffffffff811161022957820160407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdc82360301126102295760405190611872826105e2565b61187e60248201610dc8565b8252604481013567ffffffffffffffff811161022957602491010136601f820112156102295780356118af81610c17565b916118bd604051938461061a565b818352602060e081850193028201019036821161022957602001915b8183106118f857505050918160209384809401528152019301926117f5565b82360360e081126102295760c07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe060405192611933846105e2565b863561193e8161020b565b845201126102295760e091602091604051611958816105fe565b611963848801611289565b815261197160408801611289565b848201526119816060880161127a565b604082015261199260808801611289565b60608201526119a360a08801611289565b608082015260c08701356119b681610c2f565b60a0820152838201528152019201916118d9565b346102295760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102295760005473ffffffffffffffffffffffffffffffffffffffff81163303611a89577fffffffffffffffffffffffff00000000000000000000000000000000000000006001549133828416176001551660005573ffffffffffffffffffffffffffffffffffffffff3391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0600080a3005b7f02b543c60000000000000000000000000000000000000000000000000000000060005260046000fd5b9080601f83011215610229578135611aca81610c17565b92611ad8604051948561061a565b81845260208085019260051b82010192831161022957602001905b828210611b005750505090565b602080918335611b0f8161020b565b815201910190611af3565b346102295760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102295760043567ffffffffffffffff811161022957611b69903690600401611ab3565b60243567ffffffffffffffff811161022957611b89903690600401611ab3565b90611b926141c4565b60005b8151811015611c0e5780611bb6611bb1610bc4600194866130a4565b61599e565b611bc1575b01611b95565b73ffffffffffffffffffffffffffffffffffffffff611be3610bc483866130a4565b167f1795838dc8ab2ffc5f431a1729a6afa0b587f982f7b2be0b9d7187a1ef547f91600080a2611bbb565b8260005b8151811015610cda5780611c33611c2e610bc4600194866130a4565b6159bf565b611c3e575b01611c12565b73ffffffffffffffffffffffffffffffffffffffff611c60610bc483866130a4565b167fdf1b1bd32a69711488d71554706bb130b1fc63a5fa1a2cd85e8440f84065ba23600080a2611c38565b346102295760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102295760043567ffffffffffffffff811161022957611cda903690600401610ddd565b6024359167ffffffffffffffff831161022957611d33611d2b611d11611d07611d3b963690600401610ddd565b9490953691613931565b90604082015190605e604a84015160601c93015191929190565b919033614e61565b810190613996565b60005b8151811015610cda57611da0611d9b611d75611d5a84866130a4565b515173ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff166000526007602052604060002090565b613a55565b611db4611db06040830151151590565b1590565b611f765790611e29611dcc6020600194015160ff1690565b611e23611e026020611dde86896130a4565b5101517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1690565b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1690565b90614f40565b611e446040611e3884876130a4565b51015163ffffffff1690565b63ffffffff611e6f611e66611e5f610b9e611d5a888b6130a4565b5460e01c90565b63ffffffff1690565b911610611f7057611ed2611e886040611e3885886130a4565b611ec2611e9361065b565b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff851681529163ffffffff166020830152565b610a23610b9e611d5a86896130a4565b7f52f50aa6d1a95a4595361ecf953d095f125d442e4673716dede699e049de148a73ffffffffffffffffffffffffffffffffffffffff611f15611d5a85886130a4565b611f66611f276040611e38888b6130a4565b60405193849316958390929163ffffffff6020917bffffffffffffffffffffffffffffffffffffffffffffffffffffffff604085019616845216910152565b0390a25b01611d3e565b50611f6a565b611fc8611f86611d5a84866130a4565b7f06439c6b0000000000000000000000000000000000000000000000000000000060005273ffffffffffffffffffffffffffffffffffffffff16600452602490565b6000fd5b346102295760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102295761075661207f612009610d9a565b67ffffffffffffffff6024359161201f8361020b565b600060a060405161202f816105fe565b828152826020820152826040820152826060820152826080820152015216600052600a60205260406000209073ffffffffffffffffffffffffffffffffffffffff16600052602052604060002090565b6120fb6120f260405192612092846105fe565b5463ffffffff8116845263ffffffff8160201c16602085015261ffff8160401c1660408501526120d96120cc8263ffffffff9060501c1690565b63ffffffff166060860152565b63ffffffff607082901c16608085015260901c60ff1690565b151560a0830152565b6040519182918291909160a08060c083019463ffffffff815116845263ffffffff602082015116602085015261ffff604082015116604085015263ffffffff606082015116606085015263ffffffff608082015116608085015201511515910152565b346102295760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261022957602073ffffffffffffffffffffffffffffffffffffffff60015416604051908152f35b346102295760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102295760043567ffffffffffffffff81116102295760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc823603011261022957604051612229816105e2565b816004013567ffffffffffffffff81116102295761224d9060043691850101611ab3565b8152602482013567ffffffffffffffff811161022957610cda9260046122769236920101611ab3565b6020820152613af2565b346102295760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102295760043567ffffffffffffffff81116102295736602382011215610229578060040135906122db82610c17565b906122e9604051928361061a565b8282526024602083019360061b8201019036821161022957602401925b81841061231657610cda83613c92565b6040843603126102295760206040918251612330816105e2565b863561233b8161020b565b8152612348838801610dc8565b83820152815201930192612306565b346102295760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102295773ffffffffffffffffffffffffffffffffffffffff6004356123a78161020b565b6123af612df0565b50166000526007602052610756604060002060ff604051916123d0836105a5565b5473ffffffffffffffffffffffffffffffffffffffff81168352818160a01c16602084015260a81c161515604082015260405191829182919091604080606083019473ffffffffffffffffffffffffffffffffffffffff815116845260ff602082015116602085015201511515910152565b346102295760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261022957604051806020600b54918281520190600b6000527f0175b7a638427703f0dbe7bb9bbf987a2551717b34e79f33b5b1008d1fa01db99060005b8181106124c157610756856108298187038261061a565b82548452602090930192600192830192016124aa565b346102295760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102295760406125176004356110fa8161020b565b6125518251809263ffffffff602080927bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8151168552015116910152565bf35b346102295760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261022957602060405160128152f35b346102295760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610229576125c4610d9a565b6024359067ffffffffffffffff821161022957816004019060a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc8436030112610229576126296116ce8267ffffffffffffffff166000526009602052604060002090565b91612637611db08451151590565b61298257606484019061267b611db061264f84613208565b73ffffffffffffffffffffffffffffffffffffffff166000526001600b01602052604060002054151590565b61293457604485019361268e8583612e0f565b94905060248701966126c16126a38986613ecb565b9050876126ba6126b38880613ecb565b3691613931565b918661526a565b6126cd61118086613208565b966126e96126e36101c086015163ffffffff1690565b8461532a565b92600080891561290c57505090612717918961270e6127078a613208565b9389612e0f565b93909288615477565b9591955b60009861ffff61273161010089015161ffff1690565b166128a5575b5088612888896128888a8f8f9861287a61285461284e8f6107569f6128959f6101807bffffffffffffffffffffffffffffffffffffffffffffffffffffffff9f61288d9f936dffffffffffffffffffffffffffff61282a61284195612823896128009f6128419f60846128156128319a63ffffffff6128009f85826128068a6128006127f960a06127ef61280d9861281d9f6127e589916127df60808a015163ffffffff1690565b9b613ecb565b9290501690613f1c565b93015161ffff1690565b61ffff1690565b90612d5b565b9116613f1c565b911690613f1c565b950190613ecb565b90615795565b5190613f1c565b9116612d5b565b93015167ffffffffffffffff1690565b67ffffffffffffffff1690565b96613208565b73ffffffffffffffffffffffffffffffffffffffff166000526008602052604060002090565b5467ffffffffffffffff1690565b613f1c565b911690612d6e565b6040519081529081906020820190565b8597929a9950928181999594889d9897936128d289956dffffffffffffffffffffffffffff9060701c1690565b6dffffffffffffffffffffffffffff16916128ed8b86613ecb565b90506128f993876156f4565b9a9b9398965091509293979a959a612737565b915095915061292e612929611e666101a088015163ffffffff1690565b612d14565b9561271b565b611fc861294083613208565b7f2502348c0000000000000000000000000000000000000000000000000000000060005273ffffffffffffffffffffffffffffffffffffffff16600452602490565b7f99ac52f20000000000000000000000000000000000000000000000000000000060005267ffffffffffffffff821660045260246000fd5b346102295760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102295773ffffffffffffffffffffffffffffffffffffffff600435612a0a8161020b565b612a126141c4565b16338114612a8457807fffffffffffffffffffffffff0000000000000000000000000000000000000000600054161760005573ffffffffffffffffffffffffffffffffffffffff600154167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278600080a3005b7fdad89dca0000000000000000000000000000000000000000000000000000000060005260046000fd5b60ff81160361022957565b346102295760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102295760043567ffffffffffffffff8111610229573660238201121561022957806004013590612b1482610c17565b90612b22604051928361061a565b8282526024602083019360071b8201019036821161022957602401925b818410612b4f57610cda83613f29565b833603608081126102295760607fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe060405192612b8a846105e2565b8735612b958161020b565b8452011261022957608091602091604051612baf816105a5565b83880135612bbc8161020b565b81526040880135612bcc81612aae565b848201526060880135612bde81610c2f565b604082015283820152815201930192612b3f565b346102295760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261022957600435612c2d8161020b565b612c35610db1565b9067ffffffffffffffff82169182600052600960205260ff6040600020541615612cb757612c65612c86926140d1565b92600052600960205263ffffffff60016040600020015460801c169061532a565b604080517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff9384168152919092166020820152f35b827f99ac52f20000000000000000000000000000000000000000000000000000000060005260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b90662386f26fc10000820291808304662386f26fc100001490151715612d3657565b612ce5565b90655af3107a4000820291808304655af3107a40001490151715612d3657565b81810292918115918404141715612d3657565b8115612d78570490565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b612de6612de0612ded94937bffffffffffffffffffffffffffffffffffffffffffffffffffffffff612dd981956140d1565b1690612d5b565b926140d1565b1690612d6e565b90565b60405190612dfd826105a5565b60006040838281528260208201520152565b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe181360301821215610229570180359067ffffffffffffffff821161022957602001918160061b3603831361022957565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b9190811015612ea25760061b0190565b612e63565b35907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8216820361022957565b60408136031261022957612eff602060405192612eec846105e2565b8035612ef78161020b565b845201612ea7565b602082015290565b60408136031261022957612eff602060405192612f23846105e2565b612ef781610dc8565b90612f356141c4565b60005b825181101561309f5780612f4e600192856130a4565b517f32a4ba3fa3351b11ad555d4c8ec70a744e8705607077a946807030d64b6ab1a360a073ffffffffffffffffffffffffffffffffffffffff83511692606081019373ffffffffffffffffffffffffffffffffffffffff80865116957fffff00000000000000000000000000000000000000000000000000000000000061300660208601947fffffffffffffffffffff00000000000000000000000000000000000000000000865116604088019a848c51169261583d565b977fffffffffffffffffffff000000000000000000000000000000000000000000006080870195613072875115158c600052600460205260406000209060ff7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0083541691151516179055565b8560405198511688525116602087015251166040850152511660608301525115156080820152a201612f38565b509050565b8051821015612ea25760209160051b010190565b989592919097949693977f00000000000000000000000000000000000000000000000000000000000000009073ffffffffffffffffffffffffffffffffffffffff821673ffffffffffffffffffffffffffffffffffffffff8216146000146131df575050965b6bffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000168089116131ae5750916131936131a194926131a7969463ffffffff8060016131898f67ffffffffffffffff166000526009602052604060002090565b01541616916142e8565b9660208801511515996144b0565b926146c4565b9293929190565b887f6a92a4830000000000000000000000000000000000000000000000000000000060005260045260245260446000fd5b916131e992612da7565b9661311e565b604051906131fc826105e2565b60006020838281520152565b35612ded8161020b565b9060405161321f816105e2565b91547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8116835260e01c6020830152565b906132556141c4565b60005b825181101561309f5761326b81846130a4565b51602061328a61327b84876130a4565b515167ffffffffffffffff1690565b9101519067ffffffffffffffff81168015801561349d575b8015613424575b80156133f6575b6133be5791613384826001959461333461330f886132e56133899867ffffffffffffffff166000526009602052604060002090565b015460381b7fffffffff000000000000000000000000000000000000000000000000000000001690565b7fffffffff000000000000000000000000000000000000000000000000000000001690565b61338f577f525e3d4e0c31cef19cf9426af8d2c0ddd2d576359ca26bed92aac5fadda462656040518061336787826114cc565b0390a267ffffffffffffffff166000526009602052604060002090565b614712565b01613258565b7f283b699f411baff8f1c29fe49f32a828c8151596244b8e7e4c164edd6569a8356040518061336787826114cc565b7fc35aa79d0000000000000000000000000000000000000000000000000000000060005267ffffffffffffffff821660045260246000fd5b5061016083015163ffffffff1663ffffffff61341c611e66606087015163ffffffff1690565b9116116132b0565b507f2812d52c000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000006134956102008601517fffffffff000000000000000000000000000000000000000000000000000000001690565b1614156132a9565b5063ffffffff6134b561016085015163ffffffff1690565b16156132a2565b9061023961360e60016134cd61066a565b9461358b61357d82546134e96134e38260ff1690565b15158a52565b61ffff600882901c1660208a015263ffffffff601882901c1660408a015263ffffffff603882901c1660608a015263ffffffff605882901c1660808a015261ffff607882901c1660a08a015263ffffffff608882901c1660c08a015261ffff60a882901c1660e08a015261ffff60b882901c166101008a015261ffff60c882901c166101208a015260d81c63ffffffff1690565b63ffffffff16610140880152565b015463ffffffff811661016086015267ffffffffffffffff602082901c1661018086015263ffffffff606082901c166101a086015263ffffffff608082901c166101c08601526135e660a082901c60ff1615156101e0870152565b60381b7fffffffff000000000000000000000000000000000000000000000000000000001690565b7fffffffff0000000000000000000000000000000000000000000000000000000016610200840152565b906136416141c4565b6000915b80518310156138725761365883826130a4565b519061366c825167ffffffffffffffff1690565b946020600093019367ffffffffffffffff8716935b8551805182101561385d57613698826020926130a4565b5101516136a9611d5a8389516130a4565b8151602083015163ffffffff908116911681811015613824575050608082015163ffffffff16602081106137d6575090867f94967ae9ea7729ad4f54021c1981765d2b1d954f7c92fbec340aa0a54f46b8b573ffffffffffffffffffffffffffffffffffffffff84613765858f6001999861373b6137609267ffffffffffffffff16600052600a602052604060002090565b9073ffffffffffffffffffffffffffffffffffffffff16600052602052604060002090565b614cd7565b6137cd60405192839216958291909160a08060c083019463ffffffff815116845263ffffffff602082015116602085015261ffff604082015116604085015263ffffffff606082015116606085015263ffffffff608082015116608085015201511515910152565b0390a301613681565b7f24ecdc020000000000000000000000000000000000000000000000000000000060005273ffffffffffffffffffffffffffffffffffffffff90911660045263ffffffff1660245260446000fd5b7f0b4f67a20000000000000000000000000000000000000000000000000000000060005263ffffffff9081166004521660245260446000fd5b50509550925092600191500191929092613645565b50905060005b815181101561392d578061389161327b600193856130a4565b67ffffffffffffffff73ffffffffffffffffffffffffffffffffffffffff6138da60206138be86896130a4565b51015173ffffffffffffffffffffffffffffffffffffffff1690565b60006138fe8261373b8767ffffffffffffffff16600052600a602052604060002090565b551691167f4de5b1bcbca6018c11303a2c3f4a4b4f22a1c741d8c4ba430d246ac06c5ddf8b600080a301613878565b5050565b92919267ffffffffffffffff82116105c15760405191613979601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0166020018461061a565b829481845281830111610229578281602093846000960137010152565b6020818303126102295780359067ffffffffffffffff8211610229570181601f82011215610229578035906139ca82610c17565b926139d8604051948561061a565b8284526020606081860194028301019181831161022957602001925b828410613a02575050505090565b606084830312610229576020606091604051613a1d816105a5565b8635613a288161020b565b8152613a35838801612ea7565b83820152613a4560408801611289565b60408201528152019301926139f4565b90604051613a62816105a5565b604060ff82945473ffffffffffffffffffffffffffffffffffffffff81168452818160a01c16602085015260a81c161515910152565b90610239604051613aa8816105fe565b925463ffffffff8082168552602082811c821690860152604082811c61ffff1690860152605082901c81166060860152607082901c16608085015260901c60ff16151560a0840152565b613afa6141c4565b60208101519160005b8351811015613bae5780613b1c610bc4600193876130a4565b613b58613b5373ffffffffffffffffffffffffffffffffffffffff83165b73ffffffffffffffffffffffffffffffffffffffff1690565b615cd5565b613b64575b5001613b03565b60405173ffffffffffffffffffffffffffffffffffffffff9190911681527fc3803387881faad271c47728894e3e36fac830ffc8602ca6fc07733cbda7758090602090a138613b5d565b5091505160005b815181101561392d57613bcb610bc482846130a4565b9073ffffffffffffffffffffffffffffffffffffffff821615613c68577feb1b9b92e50b7f88f9ff25d56765095ac6e91540eee214906f4036a908ffbdef613c5f83613c37613c32613b3a60019773ffffffffffffffffffffffffffffffffffffffff1690565b615c5c565b5060405173ffffffffffffffffffffffffffffffffffffffff90911681529081906020820190565b0390a101613bb5565b7f8579befe0000000000000000000000000000000000000000000000000000000060005260046000fd5b613c9a6141c4565b60005b815181101561392d578073ffffffffffffffffffffffffffffffffffffffff613cc8600193856130a4565b5151167fbb77da6f7210cdd16904228a9360133d1d7dfff99b1bc75f128da5b53e28f97d613d6567ffffffffffffffff6020613d0486896130a4565b51015116836000526008602052604060002067ffffffffffffffff82167fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000008254161790556040519182918291909167ffffffffffffffff6020820193169052565b0390a201613c9d565b91908203918211612d3657565b613d836131ef565b50613db6613db18273ffffffffffffffffffffffffffffffffffffffff166000526006602052604060002090565b613212565b6020810191613dd5613dcf611e66855163ffffffff1690565b42613d6e565b63ffffffff7f00000000000000000000000000000000000000000000000000000000000000001611613e8a57611d9b613e2e9173ffffffffffffffffffffffffffffffffffffffff166000526007602052604060002090565b613e3e611db06040830151151590565b8015613e90575b613e8a57613e52906150e1565b9163ffffffff613e7a611e66613e6f602087015163ffffffff1690565b935163ffffffff1690565b911610613e85575090565b905090565b50905090565b5073ffffffffffffffffffffffffffffffffffffffff613ec4825173ffffffffffffffffffffffffffffffffffffffff1690565b1615613e45565b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe181360301821215610229570180359067ffffffffffffffff82116102295760200191813603831361022957565b91908201809211612d3657565b90613f326141c4565b60005b825181101561309f5780613f4b600192856130a4565b517fe6a7a17d710bf0b2cd05e5397dc6f97a5da4ee79e31e234bf5f965ee2bd9a5bf6140c8602073ffffffffffffffffffffffffffffffffffffffff8451169301518360005260076020526040600020613ff773ffffffffffffffffffffffffffffffffffffffff835116829073ffffffffffffffffffffffffffffffffffffffff167fffffffffffffffffffffffff0000000000000000000000000000000000000000825416179055565b602082015181547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff1660a09190911b74ff000000000000000000000000000000000000000016177fffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffff75ff0000000000000000000000000000000000000000006040850151151560a81b16911617905560405191829182919091604080606083019473ffffffffffffffffffffffffffffffffffffffff815116845260ff602082015116602085015201511515910152565b0390a201613f35565b6140da81613d7b565b9063ffffffff602083015116158015614159575b6141155750517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff907f06439c6b000000000000000000000000000000000000000000000000000000006000521660045260246000fd5b507bffffffffffffffffffffffffffffffffffffffffffffffffffffffff825116156140ee565b3360005260036020526040600020541561419657565b7fd86ad9cf000000000000000000000000000000000000000000000000000000006000523360045260246000fd5b73ffffffffffffffffffffffffffffffffffffffff6001541633036141e557565b7f2b5c74de0000000000000000000000000000000000000000000000000000000060005260046000fd5b919091357fffffffff0000000000000000000000000000000000000000000000000000000081169260048110614243575050565b7fffffffff00000000000000000000000000000000000000000000000000000000929350829060040360031b1b161690565b909291928360041161022957831161022957600401917ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0190565b90816020910312610229575190565b90816040910312610229576020604051916142d9836105e2565b805183520151612eff81610c2f565b916142f16131ef565b5081156143e757506143326126b3828061432c7fffffffff00000000000000000000000000000000000000000000000000000000958761420f565b95614275565b91167f181dcf1000000000000000000000000000000000000000000000000000000000810361436f575080602080612ded935183010191016142bf565b7f97a657c900000000000000000000000000000000000000000000000000000000146143bf577f5247fdce0000000000000000000000000000000000000000000000000000000060005260046000fd5b806020806143d2935183010191016142b0565b6143da61065b565b9081526000602082015290565b91505067ffffffffffffffff6143fb61065b565b911681526000602082015290565b9061441382610c17565b614420604051918261061a565b8281527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe061444e8294610c17565b019060005b82811061445f57505050565b806060602080938501015201614453565b9190811015612ea25760051b810135907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6181360301821215610229570190565b9092916144d660016132e58467ffffffffffffffff166000526009602052604060002090565b906144e081614409565b9560005b8281106144f5575050505050505090565b614508614503828489612e92565b613208565b8388614522614518858484614470565b6040810190613ecb565b90506020811161463c575b50839261455c6145566126b361454c60019861457f9761457a97614470565b6020810190613ecb565b896158ce565b61373b8967ffffffffffffffff16600052600a602052604060002090565b613a98565b60a081015115614602576145e66145a060606145ba93015163ffffffff1690565b6040805163ffffffff909216602083015290928391820190565b037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0810183528261061a565b6145f0828b6130a4565b526145fb818a6130a4565b50016144e4565b506145ba6145e661463761462a8967ffffffffffffffff166000526009602052604060002090565b5460d81c63ffffffff1690565b6145a0565b915050614674611e666146678461373b8b67ffffffffffffffff16600052600a602052604060002090565b5460701c63ffffffff1690565b106146815783883861452d565b7f36f536ca0000000000000000000000000000000000000000000000000000000060005273ffffffffffffffffffffffffffffffffffffffff1660045260246000fd5b6020604051917f181dcf1000000000000000000000000000000000000000000000000000000000828401528051602484015201511515604482015260448152612ded60648261061a565b9055565b90614c8d61020060016102399461475d61472c8651151590565b829060ff7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0083541691151516179055565b6147a361476f602087015161ffff1690565b82547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000ff1660089190911b62ffff0016178255565b6147ef6147b7604087015163ffffffff1690565b82547fffffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffff1660189190911b66ffffffff00000016178255565b61483f614803606087015163ffffffff1690565b82547fffffffffffffffffffffffffffffffffffffffffff00000000ffffffffffffff1660389190911b6affffffff0000000000000016178255565b614893614853608087015163ffffffff1690565b82547fffffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffff1660589190911b6effffffff000000000000000000000016178255565b6148e76148a560a087015161ffff1690565b82547fffffffffffffffffffffffffffffff0000ffffffffffffffffffffffffffffff1660789190911b70ffff00000000000000000000000000000016178255565b6149416148fb60c087015163ffffffff1690565b82547fffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffffff1660889190911b74ffffffff000000000000000000000000000000000016178255565b61499b61495360e087015161ffff1690565b82547fffffffffffffffffff0000ffffffffffffffffffffffffffffffffffffffffff1660a89190911b76ffff00000000000000000000000000000000000000000016178255565b6149f86149ae61010087015161ffff1690565b82547fffffffffffffff0000ffffffffffffffffffffffffffffffffffffffffffffff1660b89190911b78ffff000000000000000000000000000000000000000000000016178255565b614a57614a0b61012087015161ffff1690565b82547fffffffffff0000ffffffffffffffffffffffffffffffffffffffffffffffffff1660c89190911b7affff0000000000000000000000000000000000000000000000000016178255565b614abc614a6c61014087015163ffffffff1690565b82547fff00000000ffffffffffffffffffffffffffffffffffffffffffffffffffffff1660d89190911b7effffffff00000000000000000000000000000000000000000000000000000016178255565b0192614b04614ad361016083015163ffffffff1690565b859063ffffffff167fffffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000825416179055565b614b5a614b1d61018083015167ffffffffffffffff1690565b85547fffffffffffffffffffffffffffffffffffffffff0000000000000000ffffffff1660209190911b6bffffffffffffffff0000000016178555565b614bb0614b6f6101a083015163ffffffff1690565b85547fffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffffff1660609190911b6fffffffff00000000000000000000000016178555565b614c0a614bc56101c083015163ffffffff1690565b85547fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff1660809190911b73ffffffff0000000000000000000000000000000016178555565b614c66614c1b6101e0830151151590565b85805490917fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff90911690151560a01b74ff000000000000000000000000000000000000000016179055565b01517fffffffff000000000000000000000000000000000000000000000000000000001690565b81547fffffffffffffff00000000ffffffffffffffffffffffffffffffffffffffffff1660389190911c78ffffffff00000000000000000000000000000000000000000016179055565b614e1d60a061023993614d1c63ffffffff825116859063ffffffff167fffffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000825416179055565b60208181015185546040808501517fffffffffffffffffffffffffffffffffffffffffffff000000000000ffffffff9092169290931b67ffffffff000000001691909117911b69ffff0000000000000000161784556060810151614dbf9063ffffffff1685547fffffffffffffffffffffffffffffffffffff00000000ffffffffffffffffffff1660509190911b6dffffffff0000000000000000000016178555565b614e16614dd3608083015163ffffffff1690565b85547fffffffffffffffffffffffffffff00000000ffffffffffffffffffffffffffff1660709190911b71ffffffff000000000000000000000000000016178555565b0151151590565b81547fffffffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffff1690151560901b72ff00000000000000000000000000000000000016179055565b91929092614e718282868661583d565b600052600460205260ff6040600020541615614e8d5750505050565b6040517f097e17ff00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff93841660048201529390921660248401527fffffffffffffffffffff0000000000000000000000000000000000000000000090911660448301527fffff000000000000000000000000000000000000000000000000000000000000166064820152608490fd5b0390fd5b604d8111612d3657600a0a90565b60ff1660120160ff8111612d365760ff16906024821115615005577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdc8201918211612d3657614f91614f9792614f32565b90612d6e565b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8111614fdb577bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1690565b7f10cb51d10000000000000000000000000000000000000000000000000000000060005260046000fd5b906024039060248211612d365761280061501e92614f32565b614f97565b9060ff80911691160160ff8111612d365760ff16906024821115615005577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdc8201918211612d3657614f91614f9792614f32565b519069ffffffffffffffffffff8216820361022957565b908160a0910312610229576150a281615077565b91602082015191604081015191612ded608060608401519301615077565b6040513d6000823e3d90fd5b908160209103126102295751612ded81612aae565b6150e96131ef565b5061510e613b3a613b3a835173ffffffffffffffffffffffffffffffffffffffff1690565b90604051907ffeaf968c00000000000000000000000000000000000000000000000000000000825260a082600481865afa92831561522b57600092600094615230575b5060008312614fdb576020600491604051928380927f313ce5670000000000000000000000000000000000000000000000000000000082525afa92831561522b57612ded9363ffffffff936151b7936000926151f5575b506020015160ff165b90615023565b926151e76151c361065b565b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff9095168552565b1663ffffffff166020830152565b6151b191925061521c602091823d8411615224575b615214818361061a565b8101906150cc565b9291506151a8565b503d61520a565b6150c0565b90935061525691925060a03d60a011615263575b61524e818361061a565b81019061508e565b5093925050919238615151565b503d615244565b919063ffffffff6040840151168082116152fa57505061ffff602083015116908181116152c4575050907fffffffff00000000000000000000000000000000000000000000000000000000610200610239930151166158ce565b61ffff92507fd88dddd6000000000000000000000000000000000000000000000000000000006000526004521660245260446000fd5b7f869337890000000000000000000000000000000000000000000000000000000060005260045260245260446000fd5b67ffffffffffffffff8116600052600560205260406000209160405192615350846105e2565b547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8116845260e01c9182602085015263ffffffff821692836153b4575b50505050612ded90517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1690565b63ffffffff164290810393908411612d365783116153d2578061538a565b7ff08bcb3e0000000000000000000000000000000000000000000000000000000060005267ffffffffffffffff1660045263ffffffff1660245260445260646000fd5b6040813603126102295760206040519161542e836105e2565b80356154398161020b565b83520135602082015290565b63ffffffff60209116019063ffffffff8211612d3657565b9063ffffffff8091169116019063ffffffff8211612d3657565b929593909491956000936000976000976000955b80871061549e5750505050505050929190565b9091929394959698996154ba6154b589848a612e92565b615415565b9961552261457a8c6154fe6154e38867ffffffffffffffff16600052600a602052604060002090565b915173ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff16600052602052604060002090565b91615533611db060a0850151151590565b6156a25760009b604084019061554e6127f9835161ffff1690565b615603575b5050606083015163ffffffff166155699161545d565b9b608083015161557c9063ffffffff1690565b6155859161545d565b9a82516155959063ffffffff1690565b63ffffffff166155a490612d14565b600193908083106155f75750612929611e6660206155c793015163ffffffff1690565b8082116155e657506155d891613f1c565b975b0195949392919061548b565b90506155f191613f1c565b976155da565b9150506155f191613f1c565b906128006127f9615693949f61568a9460208e61568193508d73ffffffffffffffffffffffffffffffffffffffff61564f855173ffffffffffffffffffffffffffffffffffffffff1690565b911673ffffffffffffffffffffffffffffffffffffffff82161461569b5761567791506140d1565b915b0151906159e0565b925161ffff1690565b620186a0900490565b9a3880615553565b5091615679565b989a50600191506156e86156cf6156ee926156c96129296127f96101208c015161ffff1690565b90613f1c565b996156e261014089015163ffffffff1690565b9061545d565b9b615445565b996155da565b91939093806101e00193846101e011612d36576101208102908082046101201490151715612d36576101e0910101809311612d36576127f96101006127ef612ded966dffffffffffffffffffffffffffff6157886157746157626157909a63ffffffff6128009a1690613f1c565b6128006127f960e08c015161ffff1690565b6156c9611e6660c08b015163ffffffff1690565b911690612d5b565b612d3b565b6157b4916157a16131ef565b5063ffffffff61016085015116916142e8565b90815163ffffffff806060840151161610615813576101e00151151580615807575b6157dd5790565b7fee433e990000000000000000000000000000000000000000000000000000000060005260046000fd5b506020810151156157d6565b7f4c4fc93a0000000000000000000000000000000000000000000000000000000060005260046000fd5b6040805173ffffffffffffffffffffffffffffffffffffffff9283166020820190815292909316908301527fffffffffffffffffffff0000000000000000000000000000000000000000000090921660608201527fffff0000000000000000000000000000000000000000000000000000000000009092166080830152906158c88160a081016145ba565b51902090565b7fffffffff00000000000000000000000000000000000000000000000000000000167f2812d52c00000000000000000000000000000000000000000000000000000000146159195750565b602081510361595c5761593560208251830101602083016142b0565b73ffffffffffffffffffffffffffffffffffffffff8111908115615992575b5061595c5750565b614f2e906040519182917f8d666f6000000000000000000000000000000000000000000000000000000000835260048301615a15565b61040091501038615954565b73ffffffffffffffffffffffffffffffffffffffff612ded9116600b615b09565b73ffffffffffffffffffffffffffffffffffffffff612ded9116600b615c97565b670de0b6b3a7640000917bffffffffffffffffffffffffffffffffffffffffffffffffffffffff615a119216612d5b565b0490565b906020612ded92818152019061067a565b8054821015612ea25760005260206000200190600090565b9161470e918354907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9060031b92831b921b19161790565b80548015615ada577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190615aab8282615a26565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82549160031b1b1916905555565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b6001810191806000528260205260406000205492831515600014615bf7577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8401848111612d36578354937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8501948511612d36576000958583615ba897615b999503615bae575b505050615a76565b90600052602052604060002090565b55600190565b615bde615bd891615bcf615bc5615bee9588615a26565b90549060031b1c90565b92839187615a26565b90615a3e565b8590600052602052604060002090565b55388080615b91565b50505050600090565b805490680100000000000000008210156105c15781615c2791600161470e94018155615a26565b81939154907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9060031b92831b921b19161790565b600081815260036020526040902054615c9157615c7a816002615c00565b600254906000526003602052604060002055600190565b50600090565b6000828152600182016020526040902054615cce5780615cb983600193615c00565b80549260005201602052604060002055600190565b5050600090565b600081815260036020526040902054908115615cce577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820190828211612d3657600254927fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8401938411612d36578383615ba89460009603615d71575b505050615d606002615a76565b600390600052602052604060002090565b615d60615bd891615d89615bc5615d93956002615a26565b9283916002615a26565b55388080615d5356fea164736f6c634300081a000a",
}

var FeeQuoterABI = FeeQuoterMetaData.ABI

var FeeQuoterBin = FeeQuoterMetaData.Bin

func DeployFeeQuoter(auth *bind.TransactOpts, backend bind.ContractBackend, staticConfig FeeQuoterStaticConfig, priceUpdaters []common.Address, feeTokens []common.Address, tokenPriceFeeds []FeeQuoterTokenPriceFeedUpdate, tokenTransferFeeConfigArgs []FeeQuoterTokenTransferFeeConfigArgs, premiumMultiplierWeiPerEthArgs []FeeQuoterPremiumMultiplierWeiPerEthArgs, destChainConfigArgs []FeeQuoterDestChainConfigArgs) (common.Address, *types.Transaction, *FeeQuoter, error) {
	parsed, err := FeeQuoterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(FeeQuoterBin), backend, staticConfig, priceUpdaters, feeTokens, tokenPriceFeeds, tokenTransferFeeConfigArgs, premiumMultiplierWeiPerEthArgs, destChainConfigArgs)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &FeeQuoter{address: address, abi: *parsed, FeeQuoterCaller: FeeQuoterCaller{contract: contract}, FeeQuoterTransactor: FeeQuoterTransactor{contract: contract}, FeeQuoterFilterer: FeeQuoterFilterer{contract: contract}}, nil
}

type FeeQuoter struct {
	address common.Address
	abi     abi.ABI
	FeeQuoterCaller
	FeeQuoterTransactor
	FeeQuoterFilterer
}

type FeeQuoterCaller struct {
	contract *bind.BoundContract
}

type FeeQuoterTransactor struct {
	contract *bind.BoundContract
}

type FeeQuoterFilterer struct {
	contract *bind.BoundContract
}

type FeeQuoterSession struct {
	Contract     *FeeQuoter
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type FeeQuoterCallerSession struct {
	Contract *FeeQuoterCaller
	CallOpts bind.CallOpts
}

type FeeQuoterTransactorSession struct {
	Contract     *FeeQuoterTransactor
	TransactOpts bind.TransactOpts
}

type FeeQuoterRaw struct {
	Contract *FeeQuoter
}

type FeeQuoterCallerRaw struct {
	Contract *FeeQuoterCaller
}

type FeeQuoterTransactorRaw struct {
	Contract *FeeQuoterTransactor
}

func NewFeeQuoter(address common.Address, backend bind.ContractBackend) (*FeeQuoter, error) {
	abi, err := abi.JSON(strings.NewReader(FeeQuoterABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindFeeQuoter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FeeQuoter{address: address, abi: abi, FeeQuoterCaller: FeeQuoterCaller{contract: contract}, FeeQuoterTransactor: FeeQuoterTransactor{contract: contract}, FeeQuoterFilterer: FeeQuoterFilterer{contract: contract}}, nil
}

func NewFeeQuoterCaller(address common.Address, caller bind.ContractCaller) (*FeeQuoterCaller, error) {
	contract, err := bindFeeQuoter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FeeQuoterCaller{contract: contract}, nil
}

func NewFeeQuoterTransactor(address common.Address, transactor bind.ContractTransactor) (*FeeQuoterTransactor, error) {
	contract, err := bindFeeQuoter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FeeQuoterTransactor{contract: contract}, nil
}

func NewFeeQuoterFilterer(address common.Address, filterer bind.ContractFilterer) (*FeeQuoterFilterer, error) {
	contract, err := bindFeeQuoter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FeeQuoterFilterer{contract: contract}, nil
}

func bindFeeQuoter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FeeQuoterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_FeeQuoter *FeeQuoterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FeeQuoter.Contract.FeeQuoterCaller.contract.Call(opts, result, method, params...)
}

func (_FeeQuoter *FeeQuoterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeeQuoter.Contract.FeeQuoterTransactor.contract.Transfer(opts)
}

func (_FeeQuoter *FeeQuoterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FeeQuoter.Contract.FeeQuoterTransactor.contract.Transact(opts, method, params...)
}

func (_FeeQuoter *FeeQuoterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FeeQuoter.Contract.contract.Call(opts, result, method, params...)
}

func (_FeeQuoter *FeeQuoterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeeQuoter.Contract.contract.Transfer(opts)
}

func (_FeeQuoter *FeeQuoterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FeeQuoter.Contract.contract.Transact(opts, method, params...)
}

func (_FeeQuoter *FeeQuoterCaller) FEEBASEDECIMALS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FeeQuoter.contract.Call(opts, &out, "FEE_BASE_DECIMALS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_FeeQuoter *FeeQuoterSession) FEEBASEDECIMALS() (*big.Int, error) {
	return _FeeQuoter.Contract.FEEBASEDECIMALS(&_FeeQuoter.CallOpts)
}

func (_FeeQuoter *FeeQuoterCallerSession) FEEBASEDECIMALS() (*big.Int, error) {
	return _FeeQuoter.Contract.FEEBASEDECIMALS(&_FeeQuoter.CallOpts)
}

func (_FeeQuoter *FeeQuoterCaller) KEYSTONEPRICEDECIMALS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FeeQuoter.contract.Call(opts, &out, "KEYSTONE_PRICE_DECIMALS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_FeeQuoter *FeeQuoterSession) KEYSTONEPRICEDECIMALS() (*big.Int, error) {
	return _FeeQuoter.Contract.KEYSTONEPRICEDECIMALS(&_FeeQuoter.CallOpts)
}

func (_FeeQuoter *FeeQuoterCallerSession) KEYSTONEPRICEDECIMALS() (*big.Int, error) {
	return _FeeQuoter.Contract.KEYSTONEPRICEDECIMALS(&_FeeQuoter.CallOpts)
}

func (_FeeQuoter *FeeQuoterCaller) ConvertTokenAmount(opts *bind.CallOpts, fromToken common.Address, fromTokenAmount *big.Int, toToken common.Address) (*big.Int, error) {
	var out []interface{}
	err := _FeeQuoter.contract.Call(opts, &out, "convertTokenAmount", fromToken, fromTokenAmount, toToken)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_FeeQuoter *FeeQuoterSession) ConvertTokenAmount(fromToken common.Address, fromTokenAmount *big.Int, toToken common.Address) (*big.Int, error) {
	return _FeeQuoter.Contract.ConvertTokenAmount(&_FeeQuoter.CallOpts, fromToken, fromTokenAmount, toToken)
}

func (_FeeQuoter *FeeQuoterCallerSession) ConvertTokenAmount(fromToken common.Address, fromTokenAmount *big.Int, toToken common.Address) (*big.Int, error) {
	return _FeeQuoter.Contract.ConvertTokenAmount(&_FeeQuoter.CallOpts, fromToken, fromTokenAmount, toToken)
}

func (_FeeQuoter *FeeQuoterCaller) GetAllAuthorizedCallers(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _FeeQuoter.contract.Call(opts, &out, "getAllAuthorizedCallers")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_FeeQuoter *FeeQuoterSession) GetAllAuthorizedCallers() ([]common.Address, error) {
	return _FeeQuoter.Contract.GetAllAuthorizedCallers(&_FeeQuoter.CallOpts)
}

func (_FeeQuoter *FeeQuoterCallerSession) GetAllAuthorizedCallers() ([]common.Address, error) {
	return _FeeQuoter.Contract.GetAllAuthorizedCallers(&_FeeQuoter.CallOpts)
}

func (_FeeQuoter *FeeQuoterCaller) GetDestChainConfig(opts *bind.CallOpts, destChainSelector uint64) (FeeQuoterDestChainConfig, error) {
	var out []interface{}
	err := _FeeQuoter.contract.Call(opts, &out, "getDestChainConfig", destChainSelector)

	if err != nil {
		return *new(FeeQuoterDestChainConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(FeeQuoterDestChainConfig)).(*FeeQuoterDestChainConfig)

	return out0, err

}

func (_FeeQuoter *FeeQuoterSession) GetDestChainConfig(destChainSelector uint64) (FeeQuoterDestChainConfig, error) {
	return _FeeQuoter.Contract.GetDestChainConfig(&_FeeQuoter.CallOpts, destChainSelector)
}

func (_FeeQuoter *FeeQuoterCallerSession) GetDestChainConfig(destChainSelector uint64) (FeeQuoterDestChainConfig, error) {
	return _FeeQuoter.Contract.GetDestChainConfig(&_FeeQuoter.CallOpts, destChainSelector)
}

func (_FeeQuoter *FeeQuoterCaller) GetDestinationChainGasPrice(opts *bind.CallOpts, destChainSelector uint64) (InternalTimestampedPackedUint224, error) {
	var out []interface{}
	err := _FeeQuoter.contract.Call(opts, &out, "getDestinationChainGasPrice", destChainSelector)

	if err != nil {
		return *new(InternalTimestampedPackedUint224), err
	}

	out0 := *abi.ConvertType(out[0], new(InternalTimestampedPackedUint224)).(*InternalTimestampedPackedUint224)

	return out0, err

}

func (_FeeQuoter *FeeQuoterSession) GetDestinationChainGasPrice(destChainSelector uint64) (InternalTimestampedPackedUint224, error) {
	return _FeeQuoter.Contract.GetDestinationChainGasPrice(&_FeeQuoter.CallOpts, destChainSelector)
}

func (_FeeQuoter *FeeQuoterCallerSession) GetDestinationChainGasPrice(destChainSelector uint64) (InternalTimestampedPackedUint224, error) {
	return _FeeQuoter.Contract.GetDestinationChainGasPrice(&_FeeQuoter.CallOpts, destChainSelector)
}

func (_FeeQuoter *FeeQuoterCaller) GetFeeTokens(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _FeeQuoter.contract.Call(opts, &out, "getFeeTokens")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_FeeQuoter *FeeQuoterSession) GetFeeTokens() ([]common.Address, error) {
	return _FeeQuoter.Contract.GetFeeTokens(&_FeeQuoter.CallOpts)
}

func (_FeeQuoter *FeeQuoterCallerSession) GetFeeTokens() ([]common.Address, error) {
	return _FeeQuoter.Contract.GetFeeTokens(&_FeeQuoter.CallOpts)
}

func (_FeeQuoter *FeeQuoterCaller) GetPremiumMultiplierWeiPerEth(opts *bind.CallOpts, token common.Address) (uint64, error) {
	var out []interface{}
	err := _FeeQuoter.contract.Call(opts, &out, "getPremiumMultiplierWeiPerEth", token)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_FeeQuoter *FeeQuoterSession) GetPremiumMultiplierWeiPerEth(token common.Address) (uint64, error) {
	return _FeeQuoter.Contract.GetPremiumMultiplierWeiPerEth(&_FeeQuoter.CallOpts, token)
}

func (_FeeQuoter *FeeQuoterCallerSession) GetPremiumMultiplierWeiPerEth(token common.Address) (uint64, error) {
	return _FeeQuoter.Contract.GetPremiumMultiplierWeiPerEth(&_FeeQuoter.CallOpts, token)
}

func (_FeeQuoter *FeeQuoterCaller) GetStaticConfig(opts *bind.CallOpts) (FeeQuoterStaticConfig, error) {
	var out []interface{}
	err := _FeeQuoter.contract.Call(opts, &out, "getStaticConfig")

	if err != nil {
		return *new(FeeQuoterStaticConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(FeeQuoterStaticConfig)).(*FeeQuoterStaticConfig)

	return out0, err

}

func (_FeeQuoter *FeeQuoterSession) GetStaticConfig() (FeeQuoterStaticConfig, error) {
	return _FeeQuoter.Contract.GetStaticConfig(&_FeeQuoter.CallOpts)
}

func (_FeeQuoter *FeeQuoterCallerSession) GetStaticConfig() (FeeQuoterStaticConfig, error) {
	return _FeeQuoter.Contract.GetStaticConfig(&_FeeQuoter.CallOpts)
}

func (_FeeQuoter *FeeQuoterCaller) GetTokenAndGasPrices(opts *bind.CallOpts, token common.Address, destChainSelector uint64) (GetTokenAndGasPrices,

	error) {
	var out []interface{}
	err := _FeeQuoter.contract.Call(opts, &out, "getTokenAndGasPrices", token, destChainSelector)

	outstruct := new(GetTokenAndGasPrices)
	if err != nil {
		return *outstruct, err
	}

	outstruct.TokenPrice = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.GasPriceValue = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_FeeQuoter *FeeQuoterSession) GetTokenAndGasPrices(token common.Address, destChainSelector uint64) (GetTokenAndGasPrices,

	error) {
	return _FeeQuoter.Contract.GetTokenAndGasPrices(&_FeeQuoter.CallOpts, token, destChainSelector)
}

func (_FeeQuoter *FeeQuoterCallerSession) GetTokenAndGasPrices(token common.Address, destChainSelector uint64) (GetTokenAndGasPrices,

	error) {
	return _FeeQuoter.Contract.GetTokenAndGasPrices(&_FeeQuoter.CallOpts, token, destChainSelector)
}

func (_FeeQuoter *FeeQuoterCaller) GetTokenPrice(opts *bind.CallOpts, token common.Address) (InternalTimestampedPackedUint224, error) {
	var out []interface{}
	err := _FeeQuoter.contract.Call(opts, &out, "getTokenPrice", token)

	if err != nil {
		return *new(InternalTimestampedPackedUint224), err
	}

	out0 := *abi.ConvertType(out[0], new(InternalTimestampedPackedUint224)).(*InternalTimestampedPackedUint224)

	return out0, err

}

func (_FeeQuoter *FeeQuoterSession) GetTokenPrice(token common.Address) (InternalTimestampedPackedUint224, error) {
	return _FeeQuoter.Contract.GetTokenPrice(&_FeeQuoter.CallOpts, token)
}

func (_FeeQuoter *FeeQuoterCallerSession) GetTokenPrice(token common.Address) (InternalTimestampedPackedUint224, error) {
	return _FeeQuoter.Contract.GetTokenPrice(&_FeeQuoter.CallOpts, token)
}

func (_FeeQuoter *FeeQuoterCaller) GetTokenPriceFeedConfig(opts *bind.CallOpts, token common.Address) (FeeQuoterTokenPriceFeedConfig, error) {
	var out []interface{}
	err := _FeeQuoter.contract.Call(opts, &out, "getTokenPriceFeedConfig", token)

	if err != nil {
		return *new(FeeQuoterTokenPriceFeedConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(FeeQuoterTokenPriceFeedConfig)).(*FeeQuoterTokenPriceFeedConfig)

	return out0, err

}

func (_FeeQuoter *FeeQuoterSession) GetTokenPriceFeedConfig(token common.Address) (FeeQuoterTokenPriceFeedConfig, error) {
	return _FeeQuoter.Contract.GetTokenPriceFeedConfig(&_FeeQuoter.CallOpts, token)
}

func (_FeeQuoter *FeeQuoterCallerSession) GetTokenPriceFeedConfig(token common.Address) (FeeQuoterTokenPriceFeedConfig, error) {
	return _FeeQuoter.Contract.GetTokenPriceFeedConfig(&_FeeQuoter.CallOpts, token)
}

func (_FeeQuoter *FeeQuoterCaller) GetTokenPrices(opts *bind.CallOpts, tokens []common.Address) ([]InternalTimestampedPackedUint224, error) {
	var out []interface{}
	err := _FeeQuoter.contract.Call(opts, &out, "getTokenPrices", tokens)

	if err != nil {
		return *new([]InternalTimestampedPackedUint224), err
	}

	out0 := *abi.ConvertType(out[0], new([]InternalTimestampedPackedUint224)).(*[]InternalTimestampedPackedUint224)

	return out0, err

}

func (_FeeQuoter *FeeQuoterSession) GetTokenPrices(tokens []common.Address) ([]InternalTimestampedPackedUint224, error) {
	return _FeeQuoter.Contract.GetTokenPrices(&_FeeQuoter.CallOpts, tokens)
}

func (_FeeQuoter *FeeQuoterCallerSession) GetTokenPrices(tokens []common.Address) ([]InternalTimestampedPackedUint224, error) {
	return _FeeQuoter.Contract.GetTokenPrices(&_FeeQuoter.CallOpts, tokens)
}

func (_FeeQuoter *FeeQuoterCaller) GetTokenTransferFeeConfig(opts *bind.CallOpts, destChainSelector uint64, token common.Address) (FeeQuoterTokenTransferFeeConfig, error) {
	var out []interface{}
	err := _FeeQuoter.contract.Call(opts, &out, "getTokenTransferFeeConfig", destChainSelector, token)

	if err != nil {
		return *new(FeeQuoterTokenTransferFeeConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(FeeQuoterTokenTransferFeeConfig)).(*FeeQuoterTokenTransferFeeConfig)

	return out0, err

}

func (_FeeQuoter *FeeQuoterSession) GetTokenTransferFeeConfig(destChainSelector uint64, token common.Address) (FeeQuoterTokenTransferFeeConfig, error) {
	return _FeeQuoter.Contract.GetTokenTransferFeeConfig(&_FeeQuoter.CallOpts, destChainSelector, token)
}

func (_FeeQuoter *FeeQuoterCallerSession) GetTokenTransferFeeConfig(destChainSelector uint64, token common.Address) (FeeQuoterTokenTransferFeeConfig, error) {
	return _FeeQuoter.Contract.GetTokenTransferFeeConfig(&_FeeQuoter.CallOpts, destChainSelector, token)
}

func (_FeeQuoter *FeeQuoterCaller) GetValidatedFee(opts *bind.CallOpts, destChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error) {
	var out []interface{}
	err := _FeeQuoter.contract.Call(opts, &out, "getValidatedFee", destChainSelector, message)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_FeeQuoter *FeeQuoterSession) GetValidatedFee(destChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error) {
	return _FeeQuoter.Contract.GetValidatedFee(&_FeeQuoter.CallOpts, destChainSelector, message)
}

func (_FeeQuoter *FeeQuoterCallerSession) GetValidatedFee(destChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error) {
	return _FeeQuoter.Contract.GetValidatedFee(&_FeeQuoter.CallOpts, destChainSelector, message)
}

func (_FeeQuoter *FeeQuoterCaller) GetValidatedTokenPrice(opts *bind.CallOpts, token common.Address) (*big.Int, error) {
	var out []interface{}
	err := _FeeQuoter.contract.Call(opts, &out, "getValidatedTokenPrice", token)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_FeeQuoter *FeeQuoterSession) GetValidatedTokenPrice(token common.Address) (*big.Int, error) {
	return _FeeQuoter.Contract.GetValidatedTokenPrice(&_FeeQuoter.CallOpts, token)
}

func (_FeeQuoter *FeeQuoterCallerSession) GetValidatedTokenPrice(token common.Address) (*big.Int, error) {
	return _FeeQuoter.Contract.GetValidatedTokenPrice(&_FeeQuoter.CallOpts, token)
}

func (_FeeQuoter *FeeQuoterCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FeeQuoter.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_FeeQuoter *FeeQuoterSession) Owner() (common.Address, error) {
	return _FeeQuoter.Contract.Owner(&_FeeQuoter.CallOpts)
}

func (_FeeQuoter *FeeQuoterCallerSession) Owner() (common.Address, error) {
	return _FeeQuoter.Contract.Owner(&_FeeQuoter.CallOpts)
}

func (_FeeQuoter *FeeQuoterCaller) ProcessMessageArgs(opts *bind.CallOpts, destChainSelector uint64, feeToken common.Address, feeTokenAmount *big.Int, extraArgs []byte, onRampTokenTransfers []InternalEVM2AnyTokenTransfer, sourceTokenAmounts []ClientEVMTokenAmount) (ProcessMessageArgs,

	error) {
	var out []interface{}
	err := _FeeQuoter.contract.Call(opts, &out, "processMessageArgs", destChainSelector, feeToken, feeTokenAmount, extraArgs, onRampTokenTransfers, sourceTokenAmounts)

	outstruct := new(ProcessMessageArgs)
	if err != nil {
		return *outstruct, err
	}

	outstruct.MsgFeeJuels = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.IsOutOfOrderExecution = *abi.ConvertType(out[1], new(bool)).(*bool)
	outstruct.ConvertedExtraArgs = *abi.ConvertType(out[2], new([]byte)).(*[]byte)
	outstruct.DestExecDataPerToken = *abi.ConvertType(out[3], new([][]byte)).(*[][]byte)

	return *outstruct, err

}

func (_FeeQuoter *FeeQuoterSession) ProcessMessageArgs(destChainSelector uint64, feeToken common.Address, feeTokenAmount *big.Int, extraArgs []byte, onRampTokenTransfers []InternalEVM2AnyTokenTransfer, sourceTokenAmounts []ClientEVMTokenAmount) (ProcessMessageArgs,

	error) {
	return _FeeQuoter.Contract.ProcessMessageArgs(&_FeeQuoter.CallOpts, destChainSelector, feeToken, feeTokenAmount, extraArgs, onRampTokenTransfers, sourceTokenAmounts)
}

func (_FeeQuoter *FeeQuoterCallerSession) ProcessMessageArgs(destChainSelector uint64, feeToken common.Address, feeTokenAmount *big.Int, extraArgs []byte, onRampTokenTransfers []InternalEVM2AnyTokenTransfer, sourceTokenAmounts []ClientEVMTokenAmount) (ProcessMessageArgs,

	error) {
	return _FeeQuoter.Contract.ProcessMessageArgs(&_FeeQuoter.CallOpts, destChainSelector, feeToken, feeTokenAmount, extraArgs, onRampTokenTransfers, sourceTokenAmounts)
}

func (_FeeQuoter *FeeQuoterCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _FeeQuoter.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_FeeQuoter *FeeQuoterSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _FeeQuoter.Contract.SupportsInterface(&_FeeQuoter.CallOpts, interfaceId)
}

func (_FeeQuoter *FeeQuoterCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _FeeQuoter.Contract.SupportsInterface(&_FeeQuoter.CallOpts, interfaceId)
}

func (_FeeQuoter *FeeQuoterCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _FeeQuoter.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_FeeQuoter *FeeQuoterSession) TypeAndVersion() (string, error) {
	return _FeeQuoter.Contract.TypeAndVersion(&_FeeQuoter.CallOpts)
}

func (_FeeQuoter *FeeQuoterCallerSession) TypeAndVersion() (string, error) {
	return _FeeQuoter.Contract.TypeAndVersion(&_FeeQuoter.CallOpts)
}

func (_FeeQuoter *FeeQuoterTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeeQuoter.contract.Transact(opts, "acceptOwnership")
}

func (_FeeQuoter *FeeQuoterSession) AcceptOwnership() (*types.Transaction, error) {
	return _FeeQuoter.Contract.AcceptOwnership(&_FeeQuoter.TransactOpts)
}

func (_FeeQuoter *FeeQuoterTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _FeeQuoter.Contract.AcceptOwnership(&_FeeQuoter.TransactOpts)
}

func (_FeeQuoter *FeeQuoterTransactor) ApplyAuthorizedCallerUpdates(opts *bind.TransactOpts, authorizedCallerArgs AuthorizedCallersAuthorizedCallerArgs) (*types.Transaction, error) {
	return _FeeQuoter.contract.Transact(opts, "applyAuthorizedCallerUpdates", authorizedCallerArgs)
}

func (_FeeQuoter *FeeQuoterSession) ApplyAuthorizedCallerUpdates(authorizedCallerArgs AuthorizedCallersAuthorizedCallerArgs) (*types.Transaction, error) {
	return _FeeQuoter.Contract.ApplyAuthorizedCallerUpdates(&_FeeQuoter.TransactOpts, authorizedCallerArgs)
}

func (_FeeQuoter *FeeQuoterTransactorSession) ApplyAuthorizedCallerUpdates(authorizedCallerArgs AuthorizedCallersAuthorizedCallerArgs) (*types.Transaction, error) {
	return _FeeQuoter.Contract.ApplyAuthorizedCallerUpdates(&_FeeQuoter.TransactOpts, authorizedCallerArgs)
}

func (_FeeQuoter *FeeQuoterTransactor) ApplyDestChainConfigUpdates(opts *bind.TransactOpts, destChainConfigArgs []FeeQuoterDestChainConfigArgs) (*types.Transaction, error) {
	return _FeeQuoter.contract.Transact(opts, "applyDestChainConfigUpdates", destChainConfigArgs)
}

func (_FeeQuoter *FeeQuoterSession) ApplyDestChainConfigUpdates(destChainConfigArgs []FeeQuoterDestChainConfigArgs) (*types.Transaction, error) {
	return _FeeQuoter.Contract.ApplyDestChainConfigUpdates(&_FeeQuoter.TransactOpts, destChainConfigArgs)
}

func (_FeeQuoter *FeeQuoterTransactorSession) ApplyDestChainConfigUpdates(destChainConfigArgs []FeeQuoterDestChainConfigArgs) (*types.Transaction, error) {
	return _FeeQuoter.Contract.ApplyDestChainConfigUpdates(&_FeeQuoter.TransactOpts, destChainConfigArgs)
}

func (_FeeQuoter *FeeQuoterTransactor) ApplyFeeTokensUpdates(opts *bind.TransactOpts, feeTokensToRemove []common.Address, feeTokensToAdd []common.Address) (*types.Transaction, error) {
	return _FeeQuoter.contract.Transact(opts, "applyFeeTokensUpdates", feeTokensToRemove, feeTokensToAdd)
}

func (_FeeQuoter *FeeQuoterSession) ApplyFeeTokensUpdates(feeTokensToRemove []common.Address, feeTokensToAdd []common.Address) (*types.Transaction, error) {
	return _FeeQuoter.Contract.ApplyFeeTokensUpdates(&_FeeQuoter.TransactOpts, feeTokensToRemove, feeTokensToAdd)
}

func (_FeeQuoter *FeeQuoterTransactorSession) ApplyFeeTokensUpdates(feeTokensToRemove []common.Address, feeTokensToAdd []common.Address) (*types.Transaction, error) {
	return _FeeQuoter.Contract.ApplyFeeTokensUpdates(&_FeeQuoter.TransactOpts, feeTokensToRemove, feeTokensToAdd)
}

func (_FeeQuoter *FeeQuoterTransactor) ApplyPremiumMultiplierWeiPerEthUpdates(opts *bind.TransactOpts, premiumMultiplierWeiPerEthArgs []FeeQuoterPremiumMultiplierWeiPerEthArgs) (*types.Transaction, error) {
	return _FeeQuoter.contract.Transact(opts, "applyPremiumMultiplierWeiPerEthUpdates", premiumMultiplierWeiPerEthArgs)
}

func (_FeeQuoter *FeeQuoterSession) ApplyPremiumMultiplierWeiPerEthUpdates(premiumMultiplierWeiPerEthArgs []FeeQuoterPremiumMultiplierWeiPerEthArgs) (*types.Transaction, error) {
	return _FeeQuoter.Contract.ApplyPremiumMultiplierWeiPerEthUpdates(&_FeeQuoter.TransactOpts, premiumMultiplierWeiPerEthArgs)
}

func (_FeeQuoter *FeeQuoterTransactorSession) ApplyPremiumMultiplierWeiPerEthUpdates(premiumMultiplierWeiPerEthArgs []FeeQuoterPremiumMultiplierWeiPerEthArgs) (*types.Transaction, error) {
	return _FeeQuoter.Contract.ApplyPremiumMultiplierWeiPerEthUpdates(&_FeeQuoter.TransactOpts, premiumMultiplierWeiPerEthArgs)
}

func (_FeeQuoter *FeeQuoterTransactor) ApplyTokenTransferFeeConfigUpdates(opts *bind.TransactOpts, tokenTransferFeeConfigArgs []FeeQuoterTokenTransferFeeConfigArgs, tokensToUseDefaultFeeConfigs []FeeQuoterTokenTransferFeeConfigRemoveArgs) (*types.Transaction, error) {
	return _FeeQuoter.contract.Transact(opts, "applyTokenTransferFeeConfigUpdates", tokenTransferFeeConfigArgs, tokensToUseDefaultFeeConfigs)
}

func (_FeeQuoter *FeeQuoterSession) ApplyTokenTransferFeeConfigUpdates(tokenTransferFeeConfigArgs []FeeQuoterTokenTransferFeeConfigArgs, tokensToUseDefaultFeeConfigs []FeeQuoterTokenTransferFeeConfigRemoveArgs) (*types.Transaction, error) {
	return _FeeQuoter.Contract.ApplyTokenTransferFeeConfigUpdates(&_FeeQuoter.TransactOpts, tokenTransferFeeConfigArgs, tokensToUseDefaultFeeConfigs)
}

func (_FeeQuoter *FeeQuoterTransactorSession) ApplyTokenTransferFeeConfigUpdates(tokenTransferFeeConfigArgs []FeeQuoterTokenTransferFeeConfigArgs, tokensToUseDefaultFeeConfigs []FeeQuoterTokenTransferFeeConfigRemoveArgs) (*types.Transaction, error) {
	return _FeeQuoter.Contract.ApplyTokenTransferFeeConfigUpdates(&_FeeQuoter.TransactOpts, tokenTransferFeeConfigArgs, tokensToUseDefaultFeeConfigs)
}

func (_FeeQuoter *FeeQuoterTransactor) OnReport(opts *bind.TransactOpts, metadata []byte, report []byte) (*types.Transaction, error) {
	return _FeeQuoter.contract.Transact(opts, "onReport", metadata, report)
}

func (_FeeQuoter *FeeQuoterSession) OnReport(metadata []byte, report []byte) (*types.Transaction, error) {
	return _FeeQuoter.Contract.OnReport(&_FeeQuoter.TransactOpts, metadata, report)
}

func (_FeeQuoter *FeeQuoterTransactorSession) OnReport(metadata []byte, report []byte) (*types.Transaction, error) {
	return _FeeQuoter.Contract.OnReport(&_FeeQuoter.TransactOpts, metadata, report)
}

func (_FeeQuoter *FeeQuoterTransactor) SetReportPermissions(opts *bind.TransactOpts, permissions []KeystoneFeedsPermissionHandlerPermission) (*types.Transaction, error) {
	return _FeeQuoter.contract.Transact(opts, "setReportPermissions", permissions)
}

func (_FeeQuoter *FeeQuoterSession) SetReportPermissions(permissions []KeystoneFeedsPermissionHandlerPermission) (*types.Transaction, error) {
	return _FeeQuoter.Contract.SetReportPermissions(&_FeeQuoter.TransactOpts, permissions)
}

func (_FeeQuoter *FeeQuoterTransactorSession) SetReportPermissions(permissions []KeystoneFeedsPermissionHandlerPermission) (*types.Transaction, error) {
	return _FeeQuoter.Contract.SetReportPermissions(&_FeeQuoter.TransactOpts, permissions)
}

func (_FeeQuoter *FeeQuoterTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _FeeQuoter.contract.Transact(opts, "transferOwnership", to)
}

func (_FeeQuoter *FeeQuoterSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _FeeQuoter.Contract.TransferOwnership(&_FeeQuoter.TransactOpts, to)
}

func (_FeeQuoter *FeeQuoterTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _FeeQuoter.Contract.TransferOwnership(&_FeeQuoter.TransactOpts, to)
}

func (_FeeQuoter *FeeQuoterTransactor) UpdatePrices(opts *bind.TransactOpts, priceUpdates InternalPriceUpdates) (*types.Transaction, error) {
	return _FeeQuoter.contract.Transact(opts, "updatePrices", priceUpdates)
}

func (_FeeQuoter *FeeQuoterSession) UpdatePrices(priceUpdates InternalPriceUpdates) (*types.Transaction, error) {
	return _FeeQuoter.Contract.UpdatePrices(&_FeeQuoter.TransactOpts, priceUpdates)
}

func (_FeeQuoter *FeeQuoterTransactorSession) UpdatePrices(priceUpdates InternalPriceUpdates) (*types.Transaction, error) {
	return _FeeQuoter.Contract.UpdatePrices(&_FeeQuoter.TransactOpts, priceUpdates)
}

func (_FeeQuoter *FeeQuoterTransactor) UpdateTokenPriceFeeds(opts *bind.TransactOpts, tokenPriceFeedUpdates []FeeQuoterTokenPriceFeedUpdate) (*types.Transaction, error) {
	return _FeeQuoter.contract.Transact(opts, "updateTokenPriceFeeds", tokenPriceFeedUpdates)
}

func (_FeeQuoter *FeeQuoterSession) UpdateTokenPriceFeeds(tokenPriceFeedUpdates []FeeQuoterTokenPriceFeedUpdate) (*types.Transaction, error) {
	return _FeeQuoter.Contract.UpdateTokenPriceFeeds(&_FeeQuoter.TransactOpts, tokenPriceFeedUpdates)
}

func (_FeeQuoter *FeeQuoterTransactorSession) UpdateTokenPriceFeeds(tokenPriceFeedUpdates []FeeQuoterTokenPriceFeedUpdate) (*types.Transaction, error) {
	return _FeeQuoter.Contract.UpdateTokenPriceFeeds(&_FeeQuoter.TransactOpts, tokenPriceFeedUpdates)
}

type FeeQuoterAuthorizedCallerAddedIterator struct {
	Event *FeeQuoterAuthorizedCallerAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FeeQuoterAuthorizedCallerAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeQuoterAuthorizedCallerAdded)
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
		it.Event = new(FeeQuoterAuthorizedCallerAdded)
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

func (it *FeeQuoterAuthorizedCallerAddedIterator) Error() error {
	return it.fail
}

func (it *FeeQuoterAuthorizedCallerAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FeeQuoterAuthorizedCallerAdded struct {
	Caller common.Address
	Raw    types.Log
}

func (_FeeQuoter *FeeQuoterFilterer) FilterAuthorizedCallerAdded(opts *bind.FilterOpts) (*FeeQuoterAuthorizedCallerAddedIterator, error) {

	logs, sub, err := _FeeQuoter.contract.FilterLogs(opts, "AuthorizedCallerAdded")
	if err != nil {
		return nil, err
	}
	return &FeeQuoterAuthorizedCallerAddedIterator{contract: _FeeQuoter.contract, event: "AuthorizedCallerAdded", logs: logs, sub: sub}, nil
}

func (_FeeQuoter *FeeQuoterFilterer) WatchAuthorizedCallerAdded(opts *bind.WatchOpts, sink chan<- *FeeQuoterAuthorizedCallerAdded) (event.Subscription, error) {

	logs, sub, err := _FeeQuoter.contract.WatchLogs(opts, "AuthorizedCallerAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FeeQuoterAuthorizedCallerAdded)
				if err := _FeeQuoter.contract.UnpackLog(event, "AuthorizedCallerAdded", log); err != nil {
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

func (_FeeQuoter *FeeQuoterFilterer) ParseAuthorizedCallerAdded(log types.Log) (*FeeQuoterAuthorizedCallerAdded, error) {
	event := new(FeeQuoterAuthorizedCallerAdded)
	if err := _FeeQuoter.contract.UnpackLog(event, "AuthorizedCallerAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FeeQuoterAuthorizedCallerRemovedIterator struct {
	Event *FeeQuoterAuthorizedCallerRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FeeQuoterAuthorizedCallerRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeQuoterAuthorizedCallerRemoved)
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
		it.Event = new(FeeQuoterAuthorizedCallerRemoved)
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

func (it *FeeQuoterAuthorizedCallerRemovedIterator) Error() error {
	return it.fail
}

func (it *FeeQuoterAuthorizedCallerRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FeeQuoterAuthorizedCallerRemoved struct {
	Caller common.Address
	Raw    types.Log
}

func (_FeeQuoter *FeeQuoterFilterer) FilterAuthorizedCallerRemoved(opts *bind.FilterOpts) (*FeeQuoterAuthorizedCallerRemovedIterator, error) {

	logs, sub, err := _FeeQuoter.contract.FilterLogs(opts, "AuthorizedCallerRemoved")
	if err != nil {
		return nil, err
	}
	return &FeeQuoterAuthorizedCallerRemovedIterator{contract: _FeeQuoter.contract, event: "AuthorizedCallerRemoved", logs: logs, sub: sub}, nil
}

func (_FeeQuoter *FeeQuoterFilterer) WatchAuthorizedCallerRemoved(opts *bind.WatchOpts, sink chan<- *FeeQuoterAuthorizedCallerRemoved) (event.Subscription, error) {

	logs, sub, err := _FeeQuoter.contract.WatchLogs(opts, "AuthorizedCallerRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FeeQuoterAuthorizedCallerRemoved)
				if err := _FeeQuoter.contract.UnpackLog(event, "AuthorizedCallerRemoved", log); err != nil {
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

func (_FeeQuoter *FeeQuoterFilterer) ParseAuthorizedCallerRemoved(log types.Log) (*FeeQuoterAuthorizedCallerRemoved, error) {
	event := new(FeeQuoterAuthorizedCallerRemoved)
	if err := _FeeQuoter.contract.UnpackLog(event, "AuthorizedCallerRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FeeQuoterDestChainAddedIterator struct {
	Event *FeeQuoterDestChainAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FeeQuoterDestChainAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeQuoterDestChainAdded)
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
		it.Event = new(FeeQuoterDestChainAdded)
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

func (it *FeeQuoterDestChainAddedIterator) Error() error {
	return it.fail
}

func (it *FeeQuoterDestChainAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FeeQuoterDestChainAdded struct {
	DestChainSelector uint64
	DestChainConfig   FeeQuoterDestChainConfig
	Raw               types.Log
}

func (_FeeQuoter *FeeQuoterFilterer) FilterDestChainAdded(opts *bind.FilterOpts, destChainSelector []uint64) (*FeeQuoterDestChainAddedIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _FeeQuoter.contract.FilterLogs(opts, "DestChainAdded", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &FeeQuoterDestChainAddedIterator{contract: _FeeQuoter.contract, event: "DestChainAdded", logs: logs, sub: sub}, nil
}

func (_FeeQuoter *FeeQuoterFilterer) WatchDestChainAdded(opts *bind.WatchOpts, sink chan<- *FeeQuoterDestChainAdded, destChainSelector []uint64) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _FeeQuoter.contract.WatchLogs(opts, "DestChainAdded", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FeeQuoterDestChainAdded)
				if err := _FeeQuoter.contract.UnpackLog(event, "DestChainAdded", log); err != nil {
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

func (_FeeQuoter *FeeQuoterFilterer) ParseDestChainAdded(log types.Log) (*FeeQuoterDestChainAdded, error) {
	event := new(FeeQuoterDestChainAdded)
	if err := _FeeQuoter.contract.UnpackLog(event, "DestChainAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FeeQuoterDestChainConfigUpdatedIterator struct {
	Event *FeeQuoterDestChainConfigUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FeeQuoterDestChainConfigUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeQuoterDestChainConfigUpdated)
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
		it.Event = new(FeeQuoterDestChainConfigUpdated)
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

func (it *FeeQuoterDestChainConfigUpdatedIterator) Error() error {
	return it.fail
}

func (it *FeeQuoterDestChainConfigUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FeeQuoterDestChainConfigUpdated struct {
	DestChainSelector uint64
	DestChainConfig   FeeQuoterDestChainConfig
	Raw               types.Log
}

func (_FeeQuoter *FeeQuoterFilterer) FilterDestChainConfigUpdated(opts *bind.FilterOpts, destChainSelector []uint64) (*FeeQuoterDestChainConfigUpdatedIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _FeeQuoter.contract.FilterLogs(opts, "DestChainConfigUpdated", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &FeeQuoterDestChainConfigUpdatedIterator{contract: _FeeQuoter.contract, event: "DestChainConfigUpdated", logs: logs, sub: sub}, nil
}

func (_FeeQuoter *FeeQuoterFilterer) WatchDestChainConfigUpdated(opts *bind.WatchOpts, sink chan<- *FeeQuoterDestChainConfigUpdated, destChainSelector []uint64) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _FeeQuoter.contract.WatchLogs(opts, "DestChainConfigUpdated", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FeeQuoterDestChainConfigUpdated)
				if err := _FeeQuoter.contract.UnpackLog(event, "DestChainConfigUpdated", log); err != nil {
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

func (_FeeQuoter *FeeQuoterFilterer) ParseDestChainConfigUpdated(log types.Log) (*FeeQuoterDestChainConfigUpdated, error) {
	event := new(FeeQuoterDestChainConfigUpdated)
	if err := _FeeQuoter.contract.UnpackLog(event, "DestChainConfigUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FeeQuoterFeeTokenAddedIterator struct {
	Event *FeeQuoterFeeTokenAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FeeQuoterFeeTokenAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeQuoterFeeTokenAdded)
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
		it.Event = new(FeeQuoterFeeTokenAdded)
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

func (it *FeeQuoterFeeTokenAddedIterator) Error() error {
	return it.fail
}

func (it *FeeQuoterFeeTokenAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FeeQuoterFeeTokenAdded struct {
	FeeToken common.Address
	Raw      types.Log
}

func (_FeeQuoter *FeeQuoterFilterer) FilterFeeTokenAdded(opts *bind.FilterOpts, feeToken []common.Address) (*FeeQuoterFeeTokenAddedIterator, error) {

	var feeTokenRule []interface{}
	for _, feeTokenItem := range feeToken {
		feeTokenRule = append(feeTokenRule, feeTokenItem)
	}

	logs, sub, err := _FeeQuoter.contract.FilterLogs(opts, "FeeTokenAdded", feeTokenRule)
	if err != nil {
		return nil, err
	}
	return &FeeQuoterFeeTokenAddedIterator{contract: _FeeQuoter.contract, event: "FeeTokenAdded", logs: logs, sub: sub}, nil
}

func (_FeeQuoter *FeeQuoterFilterer) WatchFeeTokenAdded(opts *bind.WatchOpts, sink chan<- *FeeQuoterFeeTokenAdded, feeToken []common.Address) (event.Subscription, error) {

	var feeTokenRule []interface{}
	for _, feeTokenItem := range feeToken {
		feeTokenRule = append(feeTokenRule, feeTokenItem)
	}

	logs, sub, err := _FeeQuoter.contract.WatchLogs(opts, "FeeTokenAdded", feeTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FeeQuoterFeeTokenAdded)
				if err := _FeeQuoter.contract.UnpackLog(event, "FeeTokenAdded", log); err != nil {
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

func (_FeeQuoter *FeeQuoterFilterer) ParseFeeTokenAdded(log types.Log) (*FeeQuoterFeeTokenAdded, error) {
	event := new(FeeQuoterFeeTokenAdded)
	if err := _FeeQuoter.contract.UnpackLog(event, "FeeTokenAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FeeQuoterFeeTokenRemovedIterator struct {
	Event *FeeQuoterFeeTokenRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FeeQuoterFeeTokenRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeQuoterFeeTokenRemoved)
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
		it.Event = new(FeeQuoterFeeTokenRemoved)
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

func (it *FeeQuoterFeeTokenRemovedIterator) Error() error {
	return it.fail
}

func (it *FeeQuoterFeeTokenRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FeeQuoterFeeTokenRemoved struct {
	FeeToken common.Address
	Raw      types.Log
}

func (_FeeQuoter *FeeQuoterFilterer) FilterFeeTokenRemoved(opts *bind.FilterOpts, feeToken []common.Address) (*FeeQuoterFeeTokenRemovedIterator, error) {

	var feeTokenRule []interface{}
	for _, feeTokenItem := range feeToken {
		feeTokenRule = append(feeTokenRule, feeTokenItem)
	}

	logs, sub, err := _FeeQuoter.contract.FilterLogs(opts, "FeeTokenRemoved", feeTokenRule)
	if err != nil {
		return nil, err
	}
	return &FeeQuoterFeeTokenRemovedIterator{contract: _FeeQuoter.contract, event: "FeeTokenRemoved", logs: logs, sub: sub}, nil
}

func (_FeeQuoter *FeeQuoterFilterer) WatchFeeTokenRemoved(opts *bind.WatchOpts, sink chan<- *FeeQuoterFeeTokenRemoved, feeToken []common.Address) (event.Subscription, error) {

	var feeTokenRule []interface{}
	for _, feeTokenItem := range feeToken {
		feeTokenRule = append(feeTokenRule, feeTokenItem)
	}

	logs, sub, err := _FeeQuoter.contract.WatchLogs(opts, "FeeTokenRemoved", feeTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FeeQuoterFeeTokenRemoved)
				if err := _FeeQuoter.contract.UnpackLog(event, "FeeTokenRemoved", log); err != nil {
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

func (_FeeQuoter *FeeQuoterFilterer) ParseFeeTokenRemoved(log types.Log) (*FeeQuoterFeeTokenRemoved, error) {
	event := new(FeeQuoterFeeTokenRemoved)
	if err := _FeeQuoter.contract.UnpackLog(event, "FeeTokenRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FeeQuoterOwnershipTransferRequestedIterator struct {
	Event *FeeQuoterOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FeeQuoterOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeQuoterOwnershipTransferRequested)
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
		it.Event = new(FeeQuoterOwnershipTransferRequested)
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

func (it *FeeQuoterOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *FeeQuoterOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FeeQuoterOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_FeeQuoter *FeeQuoterFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FeeQuoterOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FeeQuoter.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FeeQuoterOwnershipTransferRequestedIterator{contract: _FeeQuoter.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_FeeQuoter *FeeQuoterFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *FeeQuoterOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FeeQuoter.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FeeQuoterOwnershipTransferRequested)
				if err := _FeeQuoter.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_FeeQuoter *FeeQuoterFilterer) ParseOwnershipTransferRequested(log types.Log) (*FeeQuoterOwnershipTransferRequested, error) {
	event := new(FeeQuoterOwnershipTransferRequested)
	if err := _FeeQuoter.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FeeQuoterOwnershipTransferredIterator struct {
	Event *FeeQuoterOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FeeQuoterOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeQuoterOwnershipTransferred)
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
		it.Event = new(FeeQuoterOwnershipTransferred)
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

func (it *FeeQuoterOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *FeeQuoterOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FeeQuoterOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_FeeQuoter *FeeQuoterFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FeeQuoterOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FeeQuoter.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FeeQuoterOwnershipTransferredIterator{contract: _FeeQuoter.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_FeeQuoter *FeeQuoterFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FeeQuoterOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FeeQuoter.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FeeQuoterOwnershipTransferred)
				if err := _FeeQuoter.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_FeeQuoter *FeeQuoterFilterer) ParseOwnershipTransferred(log types.Log) (*FeeQuoterOwnershipTransferred, error) {
	event := new(FeeQuoterOwnershipTransferred)
	if err := _FeeQuoter.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FeeQuoterPremiumMultiplierWeiPerEthUpdatedIterator struct {
	Event *FeeQuoterPremiumMultiplierWeiPerEthUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FeeQuoterPremiumMultiplierWeiPerEthUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeQuoterPremiumMultiplierWeiPerEthUpdated)
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
		it.Event = new(FeeQuoterPremiumMultiplierWeiPerEthUpdated)
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

func (it *FeeQuoterPremiumMultiplierWeiPerEthUpdatedIterator) Error() error {
	return it.fail
}

func (it *FeeQuoterPremiumMultiplierWeiPerEthUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FeeQuoterPremiumMultiplierWeiPerEthUpdated struct {
	Token                      common.Address
	PremiumMultiplierWeiPerEth uint64
	Raw                        types.Log
}

func (_FeeQuoter *FeeQuoterFilterer) FilterPremiumMultiplierWeiPerEthUpdated(opts *bind.FilterOpts, token []common.Address) (*FeeQuoterPremiumMultiplierWeiPerEthUpdatedIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _FeeQuoter.contract.FilterLogs(opts, "PremiumMultiplierWeiPerEthUpdated", tokenRule)
	if err != nil {
		return nil, err
	}
	return &FeeQuoterPremiumMultiplierWeiPerEthUpdatedIterator{contract: _FeeQuoter.contract, event: "PremiumMultiplierWeiPerEthUpdated", logs: logs, sub: sub}, nil
}

func (_FeeQuoter *FeeQuoterFilterer) WatchPremiumMultiplierWeiPerEthUpdated(opts *bind.WatchOpts, sink chan<- *FeeQuoterPremiumMultiplierWeiPerEthUpdated, token []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _FeeQuoter.contract.WatchLogs(opts, "PremiumMultiplierWeiPerEthUpdated", tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FeeQuoterPremiumMultiplierWeiPerEthUpdated)
				if err := _FeeQuoter.contract.UnpackLog(event, "PremiumMultiplierWeiPerEthUpdated", log); err != nil {
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

func (_FeeQuoter *FeeQuoterFilterer) ParsePremiumMultiplierWeiPerEthUpdated(log types.Log) (*FeeQuoterPremiumMultiplierWeiPerEthUpdated, error) {
	event := new(FeeQuoterPremiumMultiplierWeiPerEthUpdated)
	if err := _FeeQuoter.contract.UnpackLog(event, "PremiumMultiplierWeiPerEthUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FeeQuoterPriceFeedPerTokenUpdatedIterator struct {
	Event *FeeQuoterPriceFeedPerTokenUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FeeQuoterPriceFeedPerTokenUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeQuoterPriceFeedPerTokenUpdated)
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
		it.Event = new(FeeQuoterPriceFeedPerTokenUpdated)
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

func (it *FeeQuoterPriceFeedPerTokenUpdatedIterator) Error() error {
	return it.fail
}

func (it *FeeQuoterPriceFeedPerTokenUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FeeQuoterPriceFeedPerTokenUpdated struct {
	Token           common.Address
	PriceFeedConfig FeeQuoterTokenPriceFeedConfig
	Raw             types.Log
}

func (_FeeQuoter *FeeQuoterFilterer) FilterPriceFeedPerTokenUpdated(opts *bind.FilterOpts, token []common.Address) (*FeeQuoterPriceFeedPerTokenUpdatedIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _FeeQuoter.contract.FilterLogs(opts, "PriceFeedPerTokenUpdated", tokenRule)
	if err != nil {
		return nil, err
	}
	return &FeeQuoterPriceFeedPerTokenUpdatedIterator{contract: _FeeQuoter.contract, event: "PriceFeedPerTokenUpdated", logs: logs, sub: sub}, nil
}

func (_FeeQuoter *FeeQuoterFilterer) WatchPriceFeedPerTokenUpdated(opts *bind.WatchOpts, sink chan<- *FeeQuoterPriceFeedPerTokenUpdated, token []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _FeeQuoter.contract.WatchLogs(opts, "PriceFeedPerTokenUpdated", tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FeeQuoterPriceFeedPerTokenUpdated)
				if err := _FeeQuoter.contract.UnpackLog(event, "PriceFeedPerTokenUpdated", log); err != nil {
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

func (_FeeQuoter *FeeQuoterFilterer) ParsePriceFeedPerTokenUpdated(log types.Log) (*FeeQuoterPriceFeedPerTokenUpdated, error) {
	event := new(FeeQuoterPriceFeedPerTokenUpdated)
	if err := _FeeQuoter.contract.UnpackLog(event, "PriceFeedPerTokenUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FeeQuoterReportPermissionSetIterator struct {
	Event *FeeQuoterReportPermissionSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FeeQuoterReportPermissionSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeQuoterReportPermissionSet)
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
		it.Event = new(FeeQuoterReportPermissionSet)
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

func (it *FeeQuoterReportPermissionSetIterator) Error() error {
	return it.fail
}

func (it *FeeQuoterReportPermissionSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FeeQuoterReportPermissionSet struct {
	ReportId   [32]byte
	Permission KeystoneFeedsPermissionHandlerPermission
	Raw        types.Log
}

func (_FeeQuoter *FeeQuoterFilterer) FilterReportPermissionSet(opts *bind.FilterOpts, reportId [][32]byte) (*FeeQuoterReportPermissionSetIterator, error) {

	var reportIdRule []interface{}
	for _, reportIdItem := range reportId {
		reportIdRule = append(reportIdRule, reportIdItem)
	}

	logs, sub, err := _FeeQuoter.contract.FilterLogs(opts, "ReportPermissionSet", reportIdRule)
	if err != nil {
		return nil, err
	}
	return &FeeQuoterReportPermissionSetIterator{contract: _FeeQuoter.contract, event: "ReportPermissionSet", logs: logs, sub: sub}, nil
}

func (_FeeQuoter *FeeQuoterFilterer) WatchReportPermissionSet(opts *bind.WatchOpts, sink chan<- *FeeQuoterReportPermissionSet, reportId [][32]byte) (event.Subscription, error) {

	var reportIdRule []interface{}
	for _, reportIdItem := range reportId {
		reportIdRule = append(reportIdRule, reportIdItem)
	}

	logs, sub, err := _FeeQuoter.contract.WatchLogs(opts, "ReportPermissionSet", reportIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FeeQuoterReportPermissionSet)
				if err := _FeeQuoter.contract.UnpackLog(event, "ReportPermissionSet", log); err != nil {
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

func (_FeeQuoter *FeeQuoterFilterer) ParseReportPermissionSet(log types.Log) (*FeeQuoterReportPermissionSet, error) {
	event := new(FeeQuoterReportPermissionSet)
	if err := _FeeQuoter.contract.UnpackLog(event, "ReportPermissionSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FeeQuoterTokenTransferFeeConfigDeletedIterator struct {
	Event *FeeQuoterTokenTransferFeeConfigDeleted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FeeQuoterTokenTransferFeeConfigDeletedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeQuoterTokenTransferFeeConfigDeleted)
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
		it.Event = new(FeeQuoterTokenTransferFeeConfigDeleted)
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

func (it *FeeQuoterTokenTransferFeeConfigDeletedIterator) Error() error {
	return it.fail
}

func (it *FeeQuoterTokenTransferFeeConfigDeletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FeeQuoterTokenTransferFeeConfigDeleted struct {
	DestChainSelector uint64
	Token             common.Address
	Raw               types.Log
}

func (_FeeQuoter *FeeQuoterFilterer) FilterTokenTransferFeeConfigDeleted(opts *bind.FilterOpts, destChainSelector []uint64, token []common.Address) (*FeeQuoterTokenTransferFeeConfigDeletedIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _FeeQuoter.contract.FilterLogs(opts, "TokenTransferFeeConfigDeleted", destChainSelectorRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return &FeeQuoterTokenTransferFeeConfigDeletedIterator{contract: _FeeQuoter.contract, event: "TokenTransferFeeConfigDeleted", logs: logs, sub: sub}, nil
}

func (_FeeQuoter *FeeQuoterFilterer) WatchTokenTransferFeeConfigDeleted(opts *bind.WatchOpts, sink chan<- *FeeQuoterTokenTransferFeeConfigDeleted, destChainSelector []uint64, token []common.Address) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _FeeQuoter.contract.WatchLogs(opts, "TokenTransferFeeConfigDeleted", destChainSelectorRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FeeQuoterTokenTransferFeeConfigDeleted)
				if err := _FeeQuoter.contract.UnpackLog(event, "TokenTransferFeeConfigDeleted", log); err != nil {
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

func (_FeeQuoter *FeeQuoterFilterer) ParseTokenTransferFeeConfigDeleted(log types.Log) (*FeeQuoterTokenTransferFeeConfigDeleted, error) {
	event := new(FeeQuoterTokenTransferFeeConfigDeleted)
	if err := _FeeQuoter.contract.UnpackLog(event, "TokenTransferFeeConfigDeleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FeeQuoterTokenTransferFeeConfigUpdatedIterator struct {
	Event *FeeQuoterTokenTransferFeeConfigUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FeeQuoterTokenTransferFeeConfigUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeQuoterTokenTransferFeeConfigUpdated)
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
		it.Event = new(FeeQuoterTokenTransferFeeConfigUpdated)
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

func (it *FeeQuoterTokenTransferFeeConfigUpdatedIterator) Error() error {
	return it.fail
}

func (it *FeeQuoterTokenTransferFeeConfigUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FeeQuoterTokenTransferFeeConfigUpdated struct {
	DestChainSelector      uint64
	Token                  common.Address
	TokenTransferFeeConfig FeeQuoterTokenTransferFeeConfig
	Raw                    types.Log
}

func (_FeeQuoter *FeeQuoterFilterer) FilterTokenTransferFeeConfigUpdated(opts *bind.FilterOpts, destChainSelector []uint64, token []common.Address) (*FeeQuoterTokenTransferFeeConfigUpdatedIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _FeeQuoter.contract.FilterLogs(opts, "TokenTransferFeeConfigUpdated", destChainSelectorRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return &FeeQuoterTokenTransferFeeConfigUpdatedIterator{contract: _FeeQuoter.contract, event: "TokenTransferFeeConfigUpdated", logs: logs, sub: sub}, nil
}

func (_FeeQuoter *FeeQuoterFilterer) WatchTokenTransferFeeConfigUpdated(opts *bind.WatchOpts, sink chan<- *FeeQuoterTokenTransferFeeConfigUpdated, destChainSelector []uint64, token []common.Address) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _FeeQuoter.contract.WatchLogs(opts, "TokenTransferFeeConfigUpdated", destChainSelectorRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FeeQuoterTokenTransferFeeConfigUpdated)
				if err := _FeeQuoter.contract.UnpackLog(event, "TokenTransferFeeConfigUpdated", log); err != nil {
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

func (_FeeQuoter *FeeQuoterFilterer) ParseTokenTransferFeeConfigUpdated(log types.Log) (*FeeQuoterTokenTransferFeeConfigUpdated, error) {
	event := new(FeeQuoterTokenTransferFeeConfigUpdated)
	if err := _FeeQuoter.contract.UnpackLog(event, "TokenTransferFeeConfigUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FeeQuoterUsdPerTokenUpdatedIterator struct {
	Event *FeeQuoterUsdPerTokenUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FeeQuoterUsdPerTokenUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeQuoterUsdPerTokenUpdated)
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
		it.Event = new(FeeQuoterUsdPerTokenUpdated)
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

func (it *FeeQuoterUsdPerTokenUpdatedIterator) Error() error {
	return it.fail
}

func (it *FeeQuoterUsdPerTokenUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FeeQuoterUsdPerTokenUpdated struct {
	Token     common.Address
	Value     *big.Int
	Timestamp *big.Int
	Raw       types.Log
}

func (_FeeQuoter *FeeQuoterFilterer) FilterUsdPerTokenUpdated(opts *bind.FilterOpts, token []common.Address) (*FeeQuoterUsdPerTokenUpdatedIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _FeeQuoter.contract.FilterLogs(opts, "UsdPerTokenUpdated", tokenRule)
	if err != nil {
		return nil, err
	}
	return &FeeQuoterUsdPerTokenUpdatedIterator{contract: _FeeQuoter.contract, event: "UsdPerTokenUpdated", logs: logs, sub: sub}, nil
}

func (_FeeQuoter *FeeQuoterFilterer) WatchUsdPerTokenUpdated(opts *bind.WatchOpts, sink chan<- *FeeQuoterUsdPerTokenUpdated, token []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _FeeQuoter.contract.WatchLogs(opts, "UsdPerTokenUpdated", tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FeeQuoterUsdPerTokenUpdated)
				if err := _FeeQuoter.contract.UnpackLog(event, "UsdPerTokenUpdated", log); err != nil {
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

func (_FeeQuoter *FeeQuoterFilterer) ParseUsdPerTokenUpdated(log types.Log) (*FeeQuoterUsdPerTokenUpdated, error) {
	event := new(FeeQuoterUsdPerTokenUpdated)
	if err := _FeeQuoter.contract.UnpackLog(event, "UsdPerTokenUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FeeQuoterUsdPerUnitGasUpdatedIterator struct {
	Event *FeeQuoterUsdPerUnitGasUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FeeQuoterUsdPerUnitGasUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeQuoterUsdPerUnitGasUpdated)
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
		it.Event = new(FeeQuoterUsdPerUnitGasUpdated)
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

func (it *FeeQuoterUsdPerUnitGasUpdatedIterator) Error() error {
	return it.fail
}

func (it *FeeQuoterUsdPerUnitGasUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FeeQuoterUsdPerUnitGasUpdated struct {
	DestChain uint64
	Value     *big.Int
	Timestamp *big.Int
	Raw       types.Log
}

func (_FeeQuoter *FeeQuoterFilterer) FilterUsdPerUnitGasUpdated(opts *bind.FilterOpts, destChain []uint64) (*FeeQuoterUsdPerUnitGasUpdatedIterator, error) {

	var destChainRule []interface{}
	for _, destChainItem := range destChain {
		destChainRule = append(destChainRule, destChainItem)
	}

	logs, sub, err := _FeeQuoter.contract.FilterLogs(opts, "UsdPerUnitGasUpdated", destChainRule)
	if err != nil {
		return nil, err
	}
	return &FeeQuoterUsdPerUnitGasUpdatedIterator{contract: _FeeQuoter.contract, event: "UsdPerUnitGasUpdated", logs: logs, sub: sub}, nil
}

func (_FeeQuoter *FeeQuoterFilterer) WatchUsdPerUnitGasUpdated(opts *bind.WatchOpts, sink chan<- *FeeQuoterUsdPerUnitGasUpdated, destChain []uint64) (event.Subscription, error) {

	var destChainRule []interface{}
	for _, destChainItem := range destChain {
		destChainRule = append(destChainRule, destChainItem)
	}

	logs, sub, err := _FeeQuoter.contract.WatchLogs(opts, "UsdPerUnitGasUpdated", destChainRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FeeQuoterUsdPerUnitGasUpdated)
				if err := _FeeQuoter.contract.UnpackLog(event, "UsdPerUnitGasUpdated", log); err != nil {
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

func (_FeeQuoter *FeeQuoterFilterer) ParseUsdPerUnitGasUpdated(log types.Log) (*FeeQuoterUsdPerUnitGasUpdated, error) {
	event := new(FeeQuoterUsdPerUnitGasUpdated)
	if err := _FeeQuoter.contract.UnpackLog(event, "UsdPerUnitGasUpdated", log); err != nil {
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
	DestExecDataPerToken  [][]byte
}

func (_FeeQuoter *FeeQuoter) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _FeeQuoter.abi.Events["AuthorizedCallerAdded"].ID:
		return _FeeQuoter.ParseAuthorizedCallerAdded(log)
	case _FeeQuoter.abi.Events["AuthorizedCallerRemoved"].ID:
		return _FeeQuoter.ParseAuthorizedCallerRemoved(log)
	case _FeeQuoter.abi.Events["DestChainAdded"].ID:
		return _FeeQuoter.ParseDestChainAdded(log)
	case _FeeQuoter.abi.Events["DestChainConfigUpdated"].ID:
		return _FeeQuoter.ParseDestChainConfigUpdated(log)
	case _FeeQuoter.abi.Events["FeeTokenAdded"].ID:
		return _FeeQuoter.ParseFeeTokenAdded(log)
	case _FeeQuoter.abi.Events["FeeTokenRemoved"].ID:
		return _FeeQuoter.ParseFeeTokenRemoved(log)
	case _FeeQuoter.abi.Events["OwnershipTransferRequested"].ID:
		return _FeeQuoter.ParseOwnershipTransferRequested(log)
	case _FeeQuoter.abi.Events["OwnershipTransferred"].ID:
		return _FeeQuoter.ParseOwnershipTransferred(log)
	case _FeeQuoter.abi.Events["PremiumMultiplierWeiPerEthUpdated"].ID:
		return _FeeQuoter.ParsePremiumMultiplierWeiPerEthUpdated(log)
	case _FeeQuoter.abi.Events["PriceFeedPerTokenUpdated"].ID:
		return _FeeQuoter.ParsePriceFeedPerTokenUpdated(log)
	case _FeeQuoter.abi.Events["ReportPermissionSet"].ID:
		return _FeeQuoter.ParseReportPermissionSet(log)
	case _FeeQuoter.abi.Events["TokenTransferFeeConfigDeleted"].ID:
		return _FeeQuoter.ParseTokenTransferFeeConfigDeleted(log)
	case _FeeQuoter.abi.Events["TokenTransferFeeConfigUpdated"].ID:
		return _FeeQuoter.ParseTokenTransferFeeConfigUpdated(log)
	case _FeeQuoter.abi.Events["UsdPerTokenUpdated"].ID:
		return _FeeQuoter.ParseUsdPerTokenUpdated(log)
	case _FeeQuoter.abi.Events["UsdPerUnitGasUpdated"].ID:
		return _FeeQuoter.ParseUsdPerUnitGasUpdated(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (FeeQuoterAuthorizedCallerAdded) Topic() common.Hash {
	return common.HexToHash("0xeb1b9b92e50b7f88f9ff25d56765095ac6e91540eee214906f4036a908ffbdef")
}

func (FeeQuoterAuthorizedCallerRemoved) Topic() common.Hash {
	return common.HexToHash("0xc3803387881faad271c47728894e3e36fac830ffc8602ca6fc07733cbda77580")
}

func (FeeQuoterDestChainAdded) Topic() common.Hash {
	return common.HexToHash("0x525e3d4e0c31cef19cf9426af8d2c0ddd2d576359ca26bed92aac5fadda46265")
}

func (FeeQuoterDestChainConfigUpdated) Topic() common.Hash {
	return common.HexToHash("0x283b699f411baff8f1c29fe49f32a828c8151596244b8e7e4c164edd6569a835")
}

func (FeeQuoterFeeTokenAdded) Topic() common.Hash {
	return common.HexToHash("0xdf1b1bd32a69711488d71554706bb130b1fc63a5fa1a2cd85e8440f84065ba23")
}

func (FeeQuoterFeeTokenRemoved) Topic() common.Hash {
	return common.HexToHash("0x1795838dc8ab2ffc5f431a1729a6afa0b587f982f7b2be0b9d7187a1ef547f91")
}

func (FeeQuoterOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (FeeQuoterOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (FeeQuoterPremiumMultiplierWeiPerEthUpdated) Topic() common.Hash {
	return common.HexToHash("0xbb77da6f7210cdd16904228a9360133d1d7dfff99b1bc75f128da5b53e28f97d")
}

func (FeeQuoterPriceFeedPerTokenUpdated) Topic() common.Hash {
	return common.HexToHash("0xe6a7a17d710bf0b2cd05e5397dc6f97a5da4ee79e31e234bf5f965ee2bd9a5bf")
}

func (FeeQuoterReportPermissionSet) Topic() common.Hash {
	return common.HexToHash("0x32a4ba3fa3351b11ad555d4c8ec70a744e8705607077a946807030d64b6ab1a3")
}

func (FeeQuoterTokenTransferFeeConfigDeleted) Topic() common.Hash {
	return common.HexToHash("0x4de5b1bcbca6018c11303a2c3f4a4b4f22a1c741d8c4ba430d246ac06c5ddf8b")
}

func (FeeQuoterTokenTransferFeeConfigUpdated) Topic() common.Hash {
	return common.HexToHash("0x94967ae9ea7729ad4f54021c1981765d2b1d954f7c92fbec340aa0a54f46b8b5")
}

func (FeeQuoterUsdPerTokenUpdated) Topic() common.Hash {
	return common.HexToHash("0x52f50aa6d1a95a4595361ecf953d095f125d442e4673716dede699e049de148a")
}

func (FeeQuoterUsdPerUnitGasUpdated) Topic() common.Hash {
	return common.HexToHash("0xdd84a3fa9ef9409f550d54d6affec7e9c480c878c6ab27b78912a03e1b371c6e")
}

func (_FeeQuoter *FeeQuoter) Address() common.Address {
	return _FeeQuoter.address
}

type FeeQuoterInterface interface {
	FEEBASEDECIMALS(opts *bind.CallOpts) (*big.Int, error)

	KEYSTONEPRICEDECIMALS(opts *bind.CallOpts) (*big.Int, error)

	ConvertTokenAmount(opts *bind.CallOpts, fromToken common.Address, fromTokenAmount *big.Int, toToken common.Address) (*big.Int, error)

	GetAllAuthorizedCallers(opts *bind.CallOpts) ([]common.Address, error)

	GetDestChainConfig(opts *bind.CallOpts, destChainSelector uint64) (FeeQuoterDestChainConfig, error)

	GetDestinationChainGasPrice(opts *bind.CallOpts, destChainSelector uint64) (InternalTimestampedPackedUint224, error)

	GetFeeTokens(opts *bind.CallOpts) ([]common.Address, error)

	GetPremiumMultiplierWeiPerEth(opts *bind.CallOpts, token common.Address) (uint64, error)

	GetStaticConfig(opts *bind.CallOpts) (FeeQuoterStaticConfig, error)

	GetTokenAndGasPrices(opts *bind.CallOpts, token common.Address, destChainSelector uint64) (GetTokenAndGasPrices,

		error)

	GetTokenPrice(opts *bind.CallOpts, token common.Address) (InternalTimestampedPackedUint224, error)

	GetTokenPriceFeedConfig(opts *bind.CallOpts, token common.Address) (FeeQuoterTokenPriceFeedConfig, error)

	GetTokenPrices(opts *bind.CallOpts, tokens []common.Address) ([]InternalTimestampedPackedUint224, error)

	GetTokenTransferFeeConfig(opts *bind.CallOpts, destChainSelector uint64, token common.Address) (FeeQuoterTokenTransferFeeConfig, error)

	GetValidatedFee(opts *bind.CallOpts, destChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error)

	GetValidatedTokenPrice(opts *bind.CallOpts, token common.Address) (*big.Int, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	ProcessMessageArgs(opts *bind.CallOpts, destChainSelector uint64, feeToken common.Address, feeTokenAmount *big.Int, extraArgs []byte, onRampTokenTransfers []InternalEVM2AnyTokenTransfer, sourceTokenAmounts []ClientEVMTokenAmount) (ProcessMessageArgs,

		error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	ApplyAuthorizedCallerUpdates(opts *bind.TransactOpts, authorizedCallerArgs AuthorizedCallersAuthorizedCallerArgs) (*types.Transaction, error)

	ApplyDestChainConfigUpdates(opts *bind.TransactOpts, destChainConfigArgs []FeeQuoterDestChainConfigArgs) (*types.Transaction, error)

	ApplyFeeTokensUpdates(opts *bind.TransactOpts, feeTokensToRemove []common.Address, feeTokensToAdd []common.Address) (*types.Transaction, error)

	ApplyPremiumMultiplierWeiPerEthUpdates(opts *bind.TransactOpts, premiumMultiplierWeiPerEthArgs []FeeQuoterPremiumMultiplierWeiPerEthArgs) (*types.Transaction, error)

	ApplyTokenTransferFeeConfigUpdates(opts *bind.TransactOpts, tokenTransferFeeConfigArgs []FeeQuoterTokenTransferFeeConfigArgs, tokensToUseDefaultFeeConfigs []FeeQuoterTokenTransferFeeConfigRemoveArgs) (*types.Transaction, error)

	OnReport(opts *bind.TransactOpts, metadata []byte, report []byte) (*types.Transaction, error)

	SetReportPermissions(opts *bind.TransactOpts, permissions []KeystoneFeedsPermissionHandlerPermission) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	UpdatePrices(opts *bind.TransactOpts, priceUpdates InternalPriceUpdates) (*types.Transaction, error)

	UpdateTokenPriceFeeds(opts *bind.TransactOpts, tokenPriceFeedUpdates []FeeQuoterTokenPriceFeedUpdate) (*types.Transaction, error)

	FilterAuthorizedCallerAdded(opts *bind.FilterOpts) (*FeeQuoterAuthorizedCallerAddedIterator, error)

	WatchAuthorizedCallerAdded(opts *bind.WatchOpts, sink chan<- *FeeQuoterAuthorizedCallerAdded) (event.Subscription, error)

	ParseAuthorizedCallerAdded(log types.Log) (*FeeQuoterAuthorizedCallerAdded, error)

	FilterAuthorizedCallerRemoved(opts *bind.FilterOpts) (*FeeQuoterAuthorizedCallerRemovedIterator, error)

	WatchAuthorizedCallerRemoved(opts *bind.WatchOpts, sink chan<- *FeeQuoterAuthorizedCallerRemoved) (event.Subscription, error)

	ParseAuthorizedCallerRemoved(log types.Log) (*FeeQuoterAuthorizedCallerRemoved, error)

	FilterDestChainAdded(opts *bind.FilterOpts, destChainSelector []uint64) (*FeeQuoterDestChainAddedIterator, error)

	WatchDestChainAdded(opts *bind.WatchOpts, sink chan<- *FeeQuoterDestChainAdded, destChainSelector []uint64) (event.Subscription, error)

	ParseDestChainAdded(log types.Log) (*FeeQuoterDestChainAdded, error)

	FilterDestChainConfigUpdated(opts *bind.FilterOpts, destChainSelector []uint64) (*FeeQuoterDestChainConfigUpdatedIterator, error)

	WatchDestChainConfigUpdated(opts *bind.WatchOpts, sink chan<- *FeeQuoterDestChainConfigUpdated, destChainSelector []uint64) (event.Subscription, error)

	ParseDestChainConfigUpdated(log types.Log) (*FeeQuoterDestChainConfigUpdated, error)

	FilterFeeTokenAdded(opts *bind.FilterOpts, feeToken []common.Address) (*FeeQuoterFeeTokenAddedIterator, error)

	WatchFeeTokenAdded(opts *bind.WatchOpts, sink chan<- *FeeQuoterFeeTokenAdded, feeToken []common.Address) (event.Subscription, error)

	ParseFeeTokenAdded(log types.Log) (*FeeQuoterFeeTokenAdded, error)

	FilterFeeTokenRemoved(opts *bind.FilterOpts, feeToken []common.Address) (*FeeQuoterFeeTokenRemovedIterator, error)

	WatchFeeTokenRemoved(opts *bind.WatchOpts, sink chan<- *FeeQuoterFeeTokenRemoved, feeToken []common.Address) (event.Subscription, error)

	ParseFeeTokenRemoved(log types.Log) (*FeeQuoterFeeTokenRemoved, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FeeQuoterOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *FeeQuoterOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*FeeQuoterOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FeeQuoterOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FeeQuoterOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*FeeQuoterOwnershipTransferred, error)

	FilterPremiumMultiplierWeiPerEthUpdated(opts *bind.FilterOpts, token []common.Address) (*FeeQuoterPremiumMultiplierWeiPerEthUpdatedIterator, error)

	WatchPremiumMultiplierWeiPerEthUpdated(opts *bind.WatchOpts, sink chan<- *FeeQuoterPremiumMultiplierWeiPerEthUpdated, token []common.Address) (event.Subscription, error)

	ParsePremiumMultiplierWeiPerEthUpdated(log types.Log) (*FeeQuoterPremiumMultiplierWeiPerEthUpdated, error)

	FilterPriceFeedPerTokenUpdated(opts *bind.FilterOpts, token []common.Address) (*FeeQuoterPriceFeedPerTokenUpdatedIterator, error)

	WatchPriceFeedPerTokenUpdated(opts *bind.WatchOpts, sink chan<- *FeeQuoterPriceFeedPerTokenUpdated, token []common.Address) (event.Subscription, error)

	ParsePriceFeedPerTokenUpdated(log types.Log) (*FeeQuoterPriceFeedPerTokenUpdated, error)

	FilterReportPermissionSet(opts *bind.FilterOpts, reportId [][32]byte) (*FeeQuoterReportPermissionSetIterator, error)

	WatchReportPermissionSet(opts *bind.WatchOpts, sink chan<- *FeeQuoterReportPermissionSet, reportId [][32]byte) (event.Subscription, error)

	ParseReportPermissionSet(log types.Log) (*FeeQuoterReportPermissionSet, error)

	FilterTokenTransferFeeConfigDeleted(opts *bind.FilterOpts, destChainSelector []uint64, token []common.Address) (*FeeQuoterTokenTransferFeeConfigDeletedIterator, error)

	WatchTokenTransferFeeConfigDeleted(opts *bind.WatchOpts, sink chan<- *FeeQuoterTokenTransferFeeConfigDeleted, destChainSelector []uint64, token []common.Address) (event.Subscription, error)

	ParseTokenTransferFeeConfigDeleted(log types.Log) (*FeeQuoterTokenTransferFeeConfigDeleted, error)

	FilterTokenTransferFeeConfigUpdated(opts *bind.FilterOpts, destChainSelector []uint64, token []common.Address) (*FeeQuoterTokenTransferFeeConfigUpdatedIterator, error)

	WatchTokenTransferFeeConfigUpdated(opts *bind.WatchOpts, sink chan<- *FeeQuoterTokenTransferFeeConfigUpdated, destChainSelector []uint64, token []common.Address) (event.Subscription, error)

	ParseTokenTransferFeeConfigUpdated(log types.Log) (*FeeQuoterTokenTransferFeeConfigUpdated, error)

	FilterUsdPerTokenUpdated(opts *bind.FilterOpts, token []common.Address) (*FeeQuoterUsdPerTokenUpdatedIterator, error)

	WatchUsdPerTokenUpdated(opts *bind.WatchOpts, sink chan<- *FeeQuoterUsdPerTokenUpdated, token []common.Address) (event.Subscription, error)

	ParseUsdPerTokenUpdated(log types.Log) (*FeeQuoterUsdPerTokenUpdated, error)

	FilterUsdPerUnitGasUpdated(opts *bind.FilterOpts, destChain []uint64) (*FeeQuoterUsdPerUnitGasUpdatedIterator, error)

	WatchUsdPerUnitGasUpdated(opts *bind.WatchOpts, sink chan<- *FeeQuoterUsdPerUnitGasUpdated, destChain []uint64) (event.Subscription, error)

	ParseUsdPerUnitGasUpdated(log types.Log) (*FeeQuoterUsdPerUnitGasUpdated, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
