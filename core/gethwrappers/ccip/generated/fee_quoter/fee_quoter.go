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
	DestGasPerPayloadByteBase         uint8
	DestGasPerPayloadByteHigh         uint8
	DestGasPerPayloadByteThreshold    uint16
	DestDataAvailabilityOverheadGas   uint32
	DestGasPerDataAvailabilityByte    uint16
	DestDataAvailabilityMultiplierBps uint16
	ChainFamilySelector               [4]byte
	EnforceOutOfOrder                 bool
	DefaultTokenFeeUSDCents           uint16
	DefaultTokenDestGasOverhead       uint32
	DefaultTxGasLimit                 uint32
	GasMultiplierWeiPerEth            uint64
	GasPriceStalenessThreshold        uint32
	NetworkFeeUSDCents                uint32
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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"staticConfig\",\"type\":\"tuple\",\"internalType\":\"structFeeQuoter.StaticConfig\",\"components\":[{\"name\":\"maxFeeJuelsPerMsg\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"linkToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenPriceStalenessThreshold\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"priceUpdaters\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"feeTokens\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"tokenPriceFeeds\",\"type\":\"tuple[]\",\"internalType\":\"structFeeQuoter.TokenPriceFeedUpdate[]\",\"components\":[{\"name\":\"sourceToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"feedConfig\",\"type\":\"tuple\",\"internalType\":\"structFeeQuoter.TokenPriceFeedConfig\",\"components\":[{\"name\":\"dataFeedAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenDecimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"}]}]},{\"name\":\"tokenTransferFeeConfigArgs\",\"type\":\"tuple[]\",\"internalType\":\"structFeeQuoter.TokenTransferFeeConfigArgs[]\",\"components\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"tokenTransferFeeConfigs\",\"type\":\"tuple[]\",\"internalType\":\"structFeeQuoter.TokenTransferFeeConfigSingleTokenArgs[]\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenTransferFeeConfig\",\"type\":\"tuple\",\"internalType\":\"structFeeQuoter.TokenTransferFeeConfig\",\"components\":[{\"name\":\"minFeeUSDCents\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxFeeUSDCents\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"deciBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"destGasOverhead\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"destBytesOverhead\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"}]}]}]},{\"name\":\"premiumMultiplierWeiPerEthArgs\",\"type\":\"tuple[]\",\"internalType\":\"structFeeQuoter.PremiumMultiplierWeiPerEthArgs[]\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"premiumMultiplierWeiPerEth\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"destChainConfigArgs\",\"type\":\"tuple[]\",\"internalType\":\"structFeeQuoter.DestChainConfigArgs[]\",\"components\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"destChainConfig\",\"type\":\"tuple\",\"internalType\":\"structFeeQuoter.DestChainConfig\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"maxNumberOfTokensPerMsg\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"maxDataBytes\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxPerMsgGasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"destGasOverhead\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"destGasPerPayloadByteBase\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"destGasPerPayloadByteHigh\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"destGasPerPayloadByteThreshold\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"destDataAvailabilityOverheadGas\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"destGasPerDataAvailabilityByte\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"destDataAvailabilityMultiplierBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"chainFamilySelector\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"},{\"name\":\"enforceOutOfOrder\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"defaultTokenFeeUSDCents\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"defaultTokenDestGasOverhead\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"defaultTxGasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"gasMultiplierWeiPerEth\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"gasPriceStalenessThreshold\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"networkFeeUSDCents\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"FEE_BASE_DECIMALS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"KEYSTONE_PRICE_DECIMALS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"applyAuthorizedCallerUpdates\",\"inputs\":[{\"name\":\"authorizedCallerArgs\",\"type\":\"tuple\",\"internalType\":\"structAuthorizedCallers.AuthorizedCallerArgs\",\"components\":[{\"name\":\"addedCallers\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"removedCallers\",\"type\":\"address[]\",\"internalType\":\"address[]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"applyDestChainConfigUpdates\",\"inputs\":[{\"name\":\"destChainConfigArgs\",\"type\":\"tuple[]\",\"internalType\":\"structFeeQuoter.DestChainConfigArgs[]\",\"components\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"destChainConfig\",\"type\":\"tuple\",\"internalType\":\"structFeeQuoter.DestChainConfig\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"maxNumberOfTokensPerMsg\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"maxDataBytes\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxPerMsgGasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"destGasOverhead\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"destGasPerPayloadByteBase\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"destGasPerPayloadByteHigh\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"destGasPerPayloadByteThreshold\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"destDataAvailabilityOverheadGas\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"destGasPerDataAvailabilityByte\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"destDataAvailabilityMultiplierBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"chainFamilySelector\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"},{\"name\":\"enforceOutOfOrder\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"defaultTokenFeeUSDCents\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"defaultTokenDestGasOverhead\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"defaultTxGasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"gasMultiplierWeiPerEth\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"gasPriceStalenessThreshold\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"networkFeeUSDCents\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"applyFeeTokensUpdates\",\"inputs\":[{\"name\":\"feeTokensToRemove\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"feeTokensToAdd\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"applyPremiumMultiplierWeiPerEthUpdates\",\"inputs\":[{\"name\":\"premiumMultiplierWeiPerEthArgs\",\"type\":\"tuple[]\",\"internalType\":\"structFeeQuoter.PremiumMultiplierWeiPerEthArgs[]\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"premiumMultiplierWeiPerEth\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"applyTokenTransferFeeConfigUpdates\",\"inputs\":[{\"name\":\"tokenTransferFeeConfigArgs\",\"type\":\"tuple[]\",\"internalType\":\"structFeeQuoter.TokenTransferFeeConfigArgs[]\",\"components\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"tokenTransferFeeConfigs\",\"type\":\"tuple[]\",\"internalType\":\"structFeeQuoter.TokenTransferFeeConfigSingleTokenArgs[]\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenTransferFeeConfig\",\"type\":\"tuple\",\"internalType\":\"structFeeQuoter.TokenTransferFeeConfig\",\"components\":[{\"name\":\"minFeeUSDCents\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxFeeUSDCents\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"deciBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"destGasOverhead\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"destBytesOverhead\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"}]}]}]},{\"name\":\"tokensToUseDefaultFeeConfigs\",\"type\":\"tuple[]\",\"internalType\":\"structFeeQuoter.TokenTransferFeeConfigRemoveArgs[]\",\"components\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"convertTokenAmount\",\"inputs\":[{\"name\":\"fromToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"fromTokenAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"toToken\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllAuthorizedCallers\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDestChainConfig\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structFeeQuoter.DestChainConfig\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"maxNumberOfTokensPerMsg\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"maxDataBytes\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxPerMsgGasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"destGasOverhead\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"destGasPerPayloadByteBase\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"destGasPerPayloadByteHigh\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"destGasPerPayloadByteThreshold\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"destDataAvailabilityOverheadGas\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"destGasPerDataAvailabilityByte\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"destDataAvailabilityMultiplierBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"chainFamilySelector\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"},{\"name\":\"enforceOutOfOrder\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"defaultTokenFeeUSDCents\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"defaultTokenDestGasOverhead\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"defaultTxGasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"gasMultiplierWeiPerEth\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"gasPriceStalenessThreshold\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"networkFeeUSDCents\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDestinationChainGasPrice\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structInternal.TimestampedPackedUint224\",\"components\":[{\"name\":\"value\",\"type\":\"uint224\",\"internalType\":\"uint224\"},{\"name\":\"timestamp\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFeeTokens\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPremiumMultiplierWeiPerEth\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"premiumMultiplierWeiPerEth\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStaticConfig\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structFeeQuoter.StaticConfig\",\"components\":[{\"name\":\"maxFeeJuelsPerMsg\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"linkToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenPriceStalenessThreshold\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTokenAndGasPrices\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"tokenPrice\",\"type\":\"uint224\",\"internalType\":\"uint224\"},{\"name\":\"gasPriceValue\",\"type\":\"uint224\",\"internalType\":\"uint224\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTokenPrice\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structInternal.TimestampedPackedUint224\",\"components\":[{\"name\":\"value\",\"type\":\"uint224\",\"internalType\":\"uint224\"},{\"name\":\"timestamp\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTokenPriceFeedConfig\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structFeeQuoter.TokenPriceFeedConfig\",\"components\":[{\"name\":\"dataFeedAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenDecimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTokenPrices\",\"inputs\":[{\"name\":\"tokens\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.TimestampedPackedUint224[]\",\"components\":[{\"name\":\"value\",\"type\":\"uint224\",\"internalType\":\"uint224\"},{\"name\":\"timestamp\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTokenTransferFeeConfig\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"tokenTransferFeeConfig\",\"type\":\"tuple\",\"internalType\":\"structFeeQuoter.TokenTransferFeeConfig\",\"components\":[{\"name\":\"minFeeUSDCents\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxFeeUSDCents\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"deciBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"destGasOverhead\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"destBytesOverhead\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getValidatedFee\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"message\",\"type\":\"tuple\",\"internalType\":\"structClient.EVM2AnyMessage\",\"components\":[{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"tokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structClient.EVMTokenAmount[]\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"feeToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"extraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"feeTokenAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getValidatedTokenPrice\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint224\",\"internalType\":\"uint224\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"onReport\",\"inputs\":[{\"name\":\"metadata\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"report\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"processMessageArgs\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"feeToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"feeTokenAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"extraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"onRampTokenTransfers\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.EVM2AnyTokenTransfer[]\",\"components\":[{\"name\":\"sourcePoolAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"destTokenAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"extraData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"destExecData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"sourceTokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structClient.EVMTokenAmount[]\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[{\"name\":\"msgFeeJuels\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"isOutOfOrderExecution\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"convertedExtraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destExecDataPerToken\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setReportPermissions\",\"inputs\":[{\"name\":\"permissions\",\"type\":\"tuple[]\",\"internalType\":\"structKeystoneFeedsPermissionHandler.Permission[]\",\"components\":[{\"name\":\"forwarder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"workflowName\",\"type\":\"bytes10\",\"internalType\":\"bytes10\"},{\"name\":\"reportName\",\"type\":\"bytes2\",\"internalType\":\"bytes2\"},{\"name\":\"workflowOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"isAllowed\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"typeAndVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"updatePrices\",\"inputs\":[{\"name\":\"priceUpdates\",\"type\":\"tuple\",\"internalType\":\"structInternal.PriceUpdates\",\"components\":[{\"name\":\"tokenPriceUpdates\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.TokenPriceUpdate[]\",\"components\":[{\"name\":\"sourceToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"usdPerToken\",\"type\":\"uint224\",\"internalType\":\"uint224\"}]},{\"name\":\"gasPriceUpdates\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.GasPriceUpdate[]\",\"components\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"usdPerUnitGas\",\"type\":\"uint224\",\"internalType\":\"uint224\"}]}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateTokenPriceFeeds\",\"inputs\":[{\"name\":\"tokenPriceFeedUpdates\",\"type\":\"tuple[]\",\"internalType\":\"structFeeQuoter.TokenPriceFeedUpdate[]\",\"components\":[{\"name\":\"sourceToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"feedConfig\",\"type\":\"tuple\",\"internalType\":\"structFeeQuoter.TokenPriceFeedConfig\",\"components\":[{\"name\":\"dataFeedAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenDecimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"}]}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AuthorizedCallerAdded\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AuthorizedCallerRemoved\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DestChainAdded\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"destChainConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structFeeQuoter.DestChainConfig\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"maxNumberOfTokensPerMsg\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"maxDataBytes\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxPerMsgGasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"destGasOverhead\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"destGasPerPayloadByteBase\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"destGasPerPayloadByteHigh\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"destGasPerPayloadByteThreshold\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"destDataAvailabilityOverheadGas\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"destGasPerDataAvailabilityByte\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"destDataAvailabilityMultiplierBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"chainFamilySelector\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"},{\"name\":\"enforceOutOfOrder\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"defaultTokenFeeUSDCents\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"defaultTokenDestGasOverhead\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"defaultTxGasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"gasMultiplierWeiPerEth\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"gasPriceStalenessThreshold\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"networkFeeUSDCents\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DestChainConfigUpdated\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"destChainConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structFeeQuoter.DestChainConfig\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"maxNumberOfTokensPerMsg\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"maxDataBytes\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxPerMsgGasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"destGasOverhead\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"destGasPerPayloadByteBase\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"destGasPerPayloadByteHigh\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"destGasPerPayloadByteThreshold\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"destDataAvailabilityOverheadGas\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"destGasPerDataAvailabilityByte\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"destDataAvailabilityMultiplierBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"chainFamilySelector\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"},{\"name\":\"enforceOutOfOrder\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"defaultTokenFeeUSDCents\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"defaultTokenDestGasOverhead\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"defaultTxGasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"gasMultiplierWeiPerEth\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"gasPriceStalenessThreshold\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"networkFeeUSDCents\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FeeTokenAdded\",\"inputs\":[{\"name\":\"feeToken\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FeeTokenRemoved\",\"inputs\":[{\"name\":\"feeToken\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PremiumMultiplierWeiPerEthUpdated\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"premiumMultiplierWeiPerEth\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PriceFeedPerTokenUpdated\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"priceFeedConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structFeeQuoter.TokenPriceFeedConfig\",\"components\":[{\"name\":\"dataFeedAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenDecimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReportPermissionSet\",\"inputs\":[{\"name\":\"reportId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"permission\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structKeystoneFeedsPermissionHandler.Permission\",\"components\":[{\"name\":\"forwarder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"workflowName\",\"type\":\"bytes10\",\"internalType\":\"bytes10\"},{\"name\":\"reportName\",\"type\":\"bytes2\",\"internalType\":\"bytes2\"},{\"name\":\"workflowOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"isAllowed\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TokenTransferFeeConfigDeleted\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TokenTransferFeeConfigUpdated\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenTransferFeeConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structFeeQuoter.TokenTransferFeeConfig\",\"components\":[{\"name\":\"minFeeUSDCents\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxFeeUSDCents\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"deciBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"destGasOverhead\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"destBytesOverhead\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UsdPerTokenUpdated\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UsdPerUnitGasUpdated\",\"inputs\":[{\"name\":\"destChain\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"CannotTransferToSelf\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DataFeedValueOutOfUint224Range\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DestinationChainNotEnabled\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ExtraArgOutOfOrderExecutionMustBeTrue\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FeeTokenNotSupported\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InvalidChainFamilySelector\",\"inputs\":[{\"name\":\"chainFamilySelector\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}]},{\"type\":\"error\",\"name\":\"InvalidDestBytesOverhead\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"destBytesOverhead\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"InvalidDestChainConfig\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidEVMAddress\",\"inputs\":[{\"name\":\"encodedAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"InvalidExtraArgsData\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidExtraArgsTag\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidFeeRange\",\"inputs\":[{\"name\":\"minFeeUSDCents\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxFeeUSDCents\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidSVMAddress\",\"inputs\":[{\"name\":\"SVMAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"InvalidStaticConfig\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidTokenReceiver\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MessageComputeUnitLimitTooHigh\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MessageFeeTooHigh\",\"inputs\":[{\"name\":\"msgFeeJuels\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxFeeJuelsPerMsg\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"MessageGasLimitTooHigh\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MessageTooLarge\",\"inputs\":[{\"name\":\"maxSize\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actualSize\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"MustBeProposedOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnerCannotBeZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReportForwarderUnauthorized\",\"inputs\":[{\"name\":\"forwarder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"workflowOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"workflowName\",\"type\":\"bytes10\",\"internalType\":\"bytes10\"},{\"name\":\"reportName\",\"type\":\"bytes2\",\"internalType\":\"bytes2\"}]},{\"type\":\"error\",\"name\":\"SourceTokenDataTooLarge\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"StaleGasPrice\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"threshold\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"timePassed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"TokenNotSupported\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UnauthorizedCaller\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UnsupportedNumberOfTokens\",\"inputs\":[{\"name\":\"numberOfTokens\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxNumberOfTokensPerMsg\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ZeroAddressNotAllowed\",\"inputs\":[]}]",
	Bin: "0x60e060405234611029576171018038038061001981611295565b9283398101908082036101208112611029576060136110295761003a611257565b81516001600160601b038116810361102957815261005a602083016112ba565b906020810191825261006e604084016112ce565b6040820190815260608401516001600160401b03811161102957856100949186016112f6565b60808501519094906001600160401b03811161102957866100b69183016112f6565b60a08201519096906001600160401b0381116110295782019080601f830112156110295781516100ed6100e8826112df565b611295565b9260208085848152019260071b8201019083821161102957602001915b8183106111e25750505060c08301516001600160401b0381116110295783019781601f8a011215611029578851986101446100e88b6112df565b996020808c838152019160051b830101918483116110295760208101915b838310611080575050505060e08401516001600160401b0381116110295784019382601f8601121561102957845161019c6100e8826112df565b9560208088848152019260061b8201019085821161102957602001915b81831061104457505050610100810151906001600160401b038211611029570182601f82011215611029578051906101f36100e8836112df565b93602061028081878681520194028301019181831161102957602001925b828410610e6757505050503315610e5657600180546001600160a01b031916331790556020986102408a611295565b97600089526000368137610252611276565b998a52888b8b015260005b89518110156102c4576001906001600160a01b0361027b828d61138f565b51168d6102878261157b565b610294575b50500161025d565b7fc3803387881faad271c47728894e3e36fac830ffc8602ca6fc07733cbda7758091604051908152a1388d61028c565b508a985089519660005b885181101561033f576001600160a01b036102e9828b61138f565b511690811561032e577feb1b9b92e50b7f88f9ff25d56765095ac6e91540eee214906f4036a908ffbdef8c83610320600195611503565b50604051908152a1016102ce565b6342bcdf7f60e11b60005260046000fd5b5081518a985089906001600160a01b0316158015610e44575b8015610e35575b610e245791516001600160a01b031660a05290516001600160601b03166080525163ffffffff1660c05261039286611295565b9360008552600036813760005b855181101561040e576001906103c76001600160a01b036103c0838a61138f565b5116611410565b6103d2575b0161039f565b818060a01b036103e2828961138f565b51167f1795838dc8ab2ffc5f431a1729a6afa0b587f982f7b2be0b9d7187a1ef547f91600080a26103cc565b508694508560005b84518110156104855760019061043e6001600160a01b03610437838961138f565b5116611542565b610449575b01610416565b818060a01b03610459828861138f565b51167fdf1b1bd32a69711488d71554706bb130b1fc63a5fa1a2cd85e8440f84065ba23600080a2610443565b508593508460005b835181101561054757806104a36001928661138f565b517fe6a7a17d710bf0b2cd05e5397dc6f97a5da4ee79e31e234bf5f965ee2bd9a5bf606089858060a01b038451169301518360005260078b5260406000209060ff878060a01b038251169283898060a01b03198254161781558d8301908151604082549501948460a81b8651151560a81b16918560a01b9060a01b169061ffff60a01b19161717905560405193845251168c8301525115156040820152a20161048d565b5091509160005b8251811015610abc57610561818461138f565b51856001600160401b03610575848761138f565b5151169101519080158015610aa9575b8015610a8b575b610a7757600081815260098852604090205460019392919060081b6001600160e01b03191661092f57807f71e9302ab4e912a9678ae7f5a8542856706806f2817e1bf2a20b171e265cb4ad604051806106f5868291909161024063ffffffff8161026084019580511515855261ffff602082015116602086015282604082015116604086015282606082015116606086015282608082015116608086015260ff60a08201511660a086015260ff60c08201511660c086015261ffff60e08201511660e0860152826101008201511661010086015261ffff6101208201511661012086015261ffff610140820151166101408601528260e01b61016082015116610160860152610180810151151561018086015261ffff6101a0820151166101a0860152826101c0820151166101c0860152826101e0820151166101e086015260018060401b03610200820151166102008601528261022082015116610220860152015116910152565b0390a25b60005260098752826040600020825115158382549162ffff008c83015160081b169066ffffffff000000604084015160181b166affffffff00000000000000606085015160381b16926effffffff0000000000000000000000608086015160581b169260ff60781b60a087015160781b169460ff60801b60c088015160801b169161ffff60881b60e089015160881b169063ffffffff60981b6101008a015160981b169361ffff60b81b6101208b015160b81b169661ffff60c81b6101408c015160c81b169963ffffffff60d81b6101608d015160081c169b61018060ff60f81b910151151560f81b169c8f8060f81b039a63ffffffff60d81b199961ffff60c81b199861ffff60b81b199763ffffffff60981b199661ffff60881b199560ff60801b199460ff60781b19936effffffff0000000000000000000000199260ff6affffffff000000000000001992169066ffffffffffffff19161716171617161716171617161716171617161716179063ffffffff60d81b1617178155019061ffff6101a0820151169082549165ffffffff00006101c083015160101b169269ffffffff0000000000006101e084015160301b166a01000000000000000000008860901b0361020085015160501b169263ffffffff60901b61022086015160901b169461024063ffffffff60b01b91015160b01b169563ffffffff60b01b199363ffffffff60901b19926a01000000000000000000008c60901b0319918c8060501b03191617161716171617171790550161054e565b807f2431cc0363f2f66b21782c7e3d54dd9085927981a21bd0cc6be45a51b19689e360405180610a6f868291909161024063ffffffff8161026084019580511515855261ffff602082015116602086015282604082015116604086015282606082015116606086015282608082015116608086015260ff60a08201511660a086015260ff60c08201511660c086015261ffff60e08201511660e0860152826101008201511661010086015261ffff6101208201511661012086015261ffff610140820151166101408601528260e01b61016082015116610160860152610180810151151561018086015261ffff6101a0820151166101a0860152826101c0820151166101c0860152826101e0820151166101e086015260018060401b03610200820151166102008601528261022082015116610220860152015116910152565b0390a26106f9565b63c35aa79d60e01b60005260045260246000fd5b5063ffffffff6101e08301511663ffffffff6060840151161061058c565b5063ffffffff6101e08301511615610585565b84828560005b8151811015610b42576001906001600160a01b03610ae0828561138f565b5151167fbb77da6f7210cdd16904228a9360133d1d7dfff99b1bc75f128da5b53e28f97d86848060401b0381610b16868961138f565b510151168360005260088252604060002081878060401b0319825416179055604051908152a201610ac2565b83600184610b4f83611295565b9060008252600092610e1f575b909282935b8251851015610d5e57610b74858461138f565b5180516001600160401b0316939083019190855b83518051821015610d4d57610b9e82879261138f565b51015184516001600160a01b0390610bb790849061138f565b5151169063ffffffff815116908781019163ffffffff8351169081811015610d385750506080810163ffffffff815116898110610d21575090899291838c52600a8a5260408c20600160a01b6001900386168d528a5260408c2092825163ffffffff169380549180518d1b67ffffffff0000000016916040860192835160401b69ffff000000000000000016966060810195865160501b6dffffffff00000000000000000000169063ffffffff60701b895160701b169260a001998b60ff60901b8c51151560901b169560ff60901b199363ffffffff60701b19926dffffffff000000000000000000001991600160501b60019003191617161716171617171790556040519586525163ffffffff168c8601525161ffff1660408501525163ffffffff1660608401525163ffffffff16608083015251151560a082015260c07f94967ae9ea7729ad4f54021c1981765d2b1d954f7c92fbec340aa0a54f46b8b591a3600101610b88565b6312766e0160e11b8c52600485905260245260448bfd5b6305a7b3d160e11b8c5260045260245260448afd5b505060019096019593509050610b61565b9150825b8251811015610de0576001906001600160401b03610d80828661138f565b515116828060a01b0384610d94848861138f565b5101511690808752600a855260408720848060a01b038316885285528660408120557f4de5b1bcbca6018c11303a2c3f4a4b4f22a1c741d8c4ba430d246ac06c5ddf8b8780a301610d62565b604051615af1908161161082396080518181816104420152612a0e015260a05181818161047801526129bf015260c05181818161049f01526136e50152f35b610b5c565b63d794ef9560e01b60005260046000fd5b5063ffffffff8251161561035f565b5080516001600160601b031615610358565b639b15e16f60e01b60005260046000fd5b838203610280811261102957610260610e7e611276565b91610e888761136c565b8352601f190112611029576040519161026083016001600160401b0381118482101761102e57604052610ebd6020870161135f565b8352610ecb60408701611380565b6020840152610edc606087016112ce565b6040840152610eed608087016112ce565b6060840152610efe60a087016112ce565b6080840152610f0f60c08701611351565b60a0840152610f2060e08701611351565b60c0840152610f326101008701611380565b60e0840152610f4461012087016112ce565b610100840152610f576101408701611380565b610120840152610f6a6101608701611380565b610140840152610180860151916001600160e01b0319831683036110295783602093610160610280960152610fa26101a0890161135f565b610180820152610fb56101c08901611380565b6101a0820152610fc86101e089016112ce565b6101c0820152610fdb61020089016112ce565b6101e0820152610fee610220890161136c565b61020082015261100161024089016112ce565b61022082015261101461026089016112ce565b61024082015283820152815201930192610211565b600080fd5b634e487b7160e01b600052604160045260246000fd5b60408387031261102957602060409161105b611276565b611064866112ba565b815261107183870161136c565b838201528152019201916101b9565b82516001600160401b0381116110295782016040818803601f190112611029576110a8611276565b906110b56020820161136c565b825260408101516001600160401b03811161102957602091010187601f820112156110295780516110e86100e8826112df565b91602060e08185858152019302820101908a821161102957602001915b8183106111245750505091816020938480940152815201920191610162565b828b0360e081126110295760c0611139611276565b91611143866112ba565b8352601f190112611029576040519160c08301916001600160401b0383118484101761102e5760e09360209360405261117d8488016112ce565b815261118b604088016112ce565b8482015261119b60608801611380565b60408201526111ac608088016112ce565b60608201526111bd60a088016112ce565b60808201526111ce60c0880161135f565b60a082015283820152815201920191611105565b828403608081126110295760606111f7611276565b91611201866112ba565b8352601f1901126110295760809160209161121a611257565b6112258488016112ba565b815261123360408801611351565b848201526112436060880161135f565b60408201528382015281520192019161010a565b60405190606082016001600160401b0381118382101761102e57604052565b60408051919082016001600160401b0381118382101761102e57604052565b6040519190601f01601f191682016001600160401b0381118382101761102e57604052565b51906001600160a01b038216820361102957565b519063ffffffff8216820361102957565b6001600160401b03811161102e5760051b60200190565b9080601f830112156110295781516113106100e8826112df565b9260208085848152019260051b82010192831161102957602001905b8282106113395750505090565b60208091611346846112ba565b81520191019061132c565b519060ff8216820361102957565b5190811515820361102957565b51906001600160401b038216820361102957565b519061ffff8216820361102957565b80518210156113a35760209160051b010190565b634e487b7160e01b600052603260045260246000fd5b80548210156113a35760005260206000200190600090565b805480156113fa5760001901906113e882826113b9565b8154906000199060031b1b1916905555565b634e487b7160e01b600052603160045260246000fd5b6000818152600c602052604090205480156114d15760001981018181116114bb57600b546000198101919082116114bb5781810361146a575b505050611456600b6113d1565b600052600c60205260006040812055600190565b6114a361147b61148c93600b6113b9565b90549060031b1c928392600b6113b9565b819391549060031b91821b91600019901b19161790565b9055600052600c602052604060002055388080611449565b634e487b7160e01b600052601160045260246000fd5b5050600090565b8054906801000000000000000082101561102e578161148c9160016114ff940181556113b9565b9055565b8060005260036020526040600020541560001461153c576115258160026114d8565b600254906000526003602052604060002055600190565b50600090565b80600052600c6020526040600020541560001461153c5761156481600b6114d8565b600b5490600052600c602052604060002055600190565b60008181526003602052604090205480156114d15760001981018181116114bb576002546000198101919082116114bb578082036115d5575b5050506115c160026113d1565b600052600360205260006040812055600190565b6115f76115e661148c9360026113b9565b90549060031b1c92839260026113b9565b905560005260036020526040600020553880806115b456fe6080604052600436101561001257600080fd5b60003560e01c806241e5be1461020657806301ffc9a714610201578063061877e3146101fc57806306285c69146101f7578063181f5a77146101f25780632451a627146101ed578063325c868e146101e85780633937306f146101e357806341ed29e7146101de578063430d138c146101d957806345ac924d146101d45780634ab35b0b146101cf578063514e8cff146101ca5780636def4ce7146101c5578063770e2dc4146101c057806379ba5097146101bb5780637afac322146101b6578063805f2132146101b157806382b49eb0146101ac57806387b8d879146101a75780638da5cb5b146101a257806391a2749a1461019d578063a69c64c014610198578063bf78e03f14610193578063cdc73d511461018e578063d02641a014610189578063d63d3af214610184578063d8694ccd1461017f578063f2fde38b1461017a578063fbe3f778146101755763ffdb4b371461017057600080fd5b612565565b612468565b6123ac565b611f87565b611f6b565b611f22565b611eab565b611e05565b611d4c565b611cb8565b611c91565b611a75565b6118f8565b61165d565b611524565b61140c565b6111ed565b61106e565b610e97565b610e5f565b610d96565b610c69565b610a0e565b610752565b610736565b6106b3565b61060d565b610406565b6103be565b61029a565b61022e565b6001600160a01b0381160361021c57565b600080fd5b359061022c8261020b565b565b3461021c57606060031936011261021c5760206102656004356102508161020b565b602435604435916102608361020b565b6126e7565b604051908152f35b35907fffffffff000000000000000000000000000000000000000000000000000000008216820361021c57565b3461021c57602060031936011261021c576004357fffffffff000000000000000000000000000000000000000000000000000000008116810361021c577fffffffff00000000000000000000000000000000000000000000000000000000602091167f805f2132000000000000000000000000000000000000000000000000000000008114908115610394575b811561036a575b8115610340575b506040519015158152f35b7f01ffc9a70000000000000000000000000000000000000000000000000000000091501438610335565b7f181f5a77000000000000000000000000000000000000000000000000000000008114915061032e565b7f9b645f410000000000000000000000000000000000000000000000000000000081149150610327565b3461021c57602060031936011261021c576001600160a01b036004356103e38161020b565b166000526008602052602067ffffffffffffffff60406000205416604051908152f35b3461021c57600060031936011261021c5761041f61271b565b50606060405161042e81610506565b63ffffffff6bffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016918281526001600160a01b0360406020830192827f00000000000000000000000000000000000000000000000000000000000000001684520191837f00000000000000000000000000000000000000000000000000000000000000001683526040519485525116602084015251166040820152f35b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6060810190811067ffffffffffffffff82111761052257604052565b6104d7565b60a0810190811067ffffffffffffffff82111761052257604052565b6040810190811067ffffffffffffffff82111761052257604052565b60c0810190811067ffffffffffffffff82111761052257604052565b90601f601f19910116810190811067ffffffffffffffff82111761052257604052565b6040519061022c60408361057b565b6040519061022c6102608361057b565b6040519061022c60a08361057b565b919082519283825260005b8481106105f8575050601f19601f8460006020809697860101520116010190565b806020809284010151828286010152016105d7565b3461021c57600060031936011261021c5761066c6040805190610630818361057b565b601382527f46656551756f74657220312e362e302d646576000000000000000000000000006020830152519182916020835260208301906105cc565b0390f35b602060408183019282815284518094520192019060005b8181106106945750505090565b82516001600160a01b0316845260209384019390920191600101610687565b3461021c57600060031936011261021c5760405180602060025491828152019060026000527f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace9060005b8181106107205761066c856107148187038261057b565b60405191829182610670565b82548452602090930192600192830192016106fd565b3461021c57600060031936011261021c57602060405160248152f35b3461021c57602060031936011261021c5760043567ffffffffffffffff811161021c5780600401906040600319823603011261021c57610790613a1e565b61079a828061273a565b4263ffffffff1692915060005b81811061090b575050602401906107be828461273a565b92905060005b8381106107cd57005b806107ec6107e76001936107e1868a61273a565b906127bd565b61281d565b7fdd84a3fa9ef9409f550d54d6affec7e9c480c878c6ab27b78912a03e1b371c6e67ffffffffffffffff6108e76108d960208501946108cb61083587516001600160e01b031690565b61084f61084061059e565b6001600160e01b039092168252565b63ffffffff8c16602082015261088a610870845167ffffffffffffffff1690565b67ffffffffffffffff166000526005602052604060002090565b815160209092015160e01b7fffffffff00000000000000000000000000000000000000000000000000000000166001600160e01b0392909216919091179055565b5167ffffffffffffffff1690565b93516001600160e01b031690565b604080516001600160e01b039290921682524260208301529190931692a2016107c4565b8061092461091f6001936107e1898061273a565b6127e6565b7f52f50aa6d1a95a4595361ecf953d095f125d442e4673716dede699e049de148a6001600160a01b036109bd6108d960208501946109b061096c87516001600160e01b031690565b61097761084061059e565b63ffffffff8d16602082015261088a61099784516001600160a01b031690565b6001600160a01b03166000526006602052604060002090565b516001600160a01b031690565b604080516001600160e01b039290921682524260208301529190931692a2016107a7565b67ffffffffffffffff81116105225760051b60200190565b8015150361021c57565b359061022c826109f9565b3461021c57602060031936011261021c5760043567ffffffffffffffff811161021c573660238201121561021c57806004013590610a4b826109e1565b90610a59604051928361057b565b828252602460a06020840194028201019036821161021c57602401925b818410610a8857610a8683612842565b005b60a08436031261021c5760405190610a9f82610527565b8435610aaa8161020b565b825260208501357fffffffffffffffffffff000000000000000000000000000000000000000000008116810361021c5760208301526040850135907fffff0000000000000000000000000000000000000000000000000000000000008216820361021c5782602092604060a0950152610b2560608801610221565b6060820152610b3660808801610a03565b6080820152815201930192610a76565b6004359067ffffffffffffffff8216820361021c57565b6024359067ffffffffffffffff8216820361021c57565b359067ffffffffffffffff8216820361021c57565b9181601f8401121561021c5782359167ffffffffffffffff831161021c576020838186019501011161021c57565b9181601f8401121561021c5782359167ffffffffffffffff831161021c576020808501948460051b01011161021c57565b929091610c09928452151560208401526080604084015260808301906105cc565b906060818303910152815180825260208201916020808360051b8301019401926000915b838310610c3c57505050505090565b9091929394602080610c5a83601f19866001960301875289516105cc565b97019301930191939290610c2d565b3461021c5760c060031936011261021c57610c82610b46565b60243590610c8f8261020b565b60443560643567ffffffffffffffff811161021c57610cb2903690600401610b89565b60849391933567ffffffffffffffff811161021c57610cd5903690600401610bb7565b9160a4359567ffffffffffffffff871161021c573660238801121561021c5786600401359567ffffffffffffffff871161021c573660248860061b8a01011161021c5761066c986024610d299901966129b4565b9060409492945194859485610be8565b602060408183019282815284518094520192019060005b818110610d5d5750505090565b9091926020604082610d8b600194885163ffffffff602080926001600160e01b038151168552015116910152565b019401929101610d50565b3461021c57602060031936011261021c5760043567ffffffffffffffff811161021c57610dc7903690600401610bb7565b610dd0816109e1565b91610dde604051938461057b565b818352601f19610ded836109e1565b0160005b818110610e4857505060005b82811015610e3a57600190610e1e610e198260051b8501612ccb565b613691565b610e2882876129a0565b52610e3381866129a0565b5001610dfd565b6040518061066c8682610d39565b602090610e53612cb2565b82828801015201610df1565b3461021c57602060031936011261021c576020610e86600435610e818161020b565b6139a6565b6001600160e01b0360405191168152f35b3461021c57602060031936011261021c5767ffffffffffffffff610eb9610b46565b610ec1612cb2565b50166000526005602052604060002060405190610edd82610543565b546001600160e01b038116825260e01c6020820152604051809161066c82604081019263ffffffff602080926001600160e01b038151168552015116910152565b61022c9092919261024080610260830195610f3b84825115159052565b60208181015161ffff169085015260408181015163ffffffff169085015260608181015163ffffffff169085015260808181015163ffffffff169085015260a08181015160ff169085015260c08181015160ff169085015260e08181015161ffff16908501526101008181015163ffffffff16908501526101208181015161ffff16908501526101408181015161ffff1690850152610160818101517fffffffff000000000000000000000000000000000000000000000000000000001690850152610180818101511515908501526101a08181015161ffff16908501526101c08181015163ffffffff16908501526101e08181015163ffffffff16908501526102008181015167ffffffffffffffff16908501526102208181015163ffffffff1690850152015163ffffffff16910152565b3461021c57602060031936011261021c5761066c61113161112c611090610b46565b600061024061109d6105ad565b8281528260208201528260408201528260608201528260808201528260a08201528260c08201528260e08201528261010082015282610120820152826101408201528261016082015282610180820152826101a0820152826101c0820152826101e08201528261020082015282610220820152015267ffffffffffffffff166000526009602052604060002090565b612cfa565b60405191829182610f1e565b359063ffffffff8216820361021c57565b359061ffff8216820361021c57565b81601f8201121561021c57803590611174826109e1565b92611182604051948561057b565b82845260208085019360061b8301019181831161021c57602001925b8284106111ac575050505090565b60408483031261021c57602060409182516111c681610543565b6111cf87610b74565b8152828701356111de8161020b565b8382015281520193019261119e565b3461021c57604060031936011261021c5760043567ffffffffffffffff811161021c573660238201121561021c578060040135611229816109e1565b91611237604051938461057b565b8183526024602084019260051b8201019036821161021c5760248101925b828410611286576024358567ffffffffffffffff821161021c57611280610a8692369060040161115d565b90612e68565b833567ffffffffffffffff811161021c57820160407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdc823603011261021c57604051906112d282610543565b6112de60248201610b74565b8252604481013567ffffffffffffffff811161021c57602491010136601f8201121561021c57803561130f816109e1565b9161131d604051938461057b565b818352602060e081850193028201019036821161021c57602001915b8183106113585750505091816020938480940152815201930192611255565b82360360e0811261021c5760c0601f196040519261137584610543565b86356113808161020b565b8452011261021c5760e09160209160405161139a8161055f565b6113a584880161113d565b81526113b36040880161113d565b848201526113c36060880161114e565b60408201526113d46080880161113d565b60608201526113e560a0880161113d565b608082015260c08701356113f8816109f9565b60a082015283820152815201920191611339565b3461021c57600060031936011261021c576000546001600160a01b0381163303611493577fffffffffffffffffffffffff0000000000000000000000000000000000000000600154913382841617600155166000556001600160a01b033391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0600080a3005b7f02b543c60000000000000000000000000000000000000000000000000000000060005260046000fd5b9080601f8301121561021c5781356114d4816109e1565b926114e2604051948561057b565b81845260208085019260051b82010192831161021c57602001905b82821061150a5750505090565b6020809183356115198161020b565b8152019101906114fd565b3461021c57604060031936011261021c5760043567ffffffffffffffff811161021c576115559036906004016114bd565b60243567ffffffffffffffff811161021c576115759036906004016114bd565b9061157e613a62565b60005b81518110156115ed57806115a261159d6109b0600194866129a0565b6156f6565b6115ad575b01611581565b6001600160a01b036115c26109b083866129a0565b167f1795838dc8ab2ffc5f431a1729a6afa0b587f982f7b2be0b9d7187a1ef547f91600080a26115a7565b8260005b8151811015610a86578061161261160d6109b0600194866129a0565b61570a565b61161d575b016115f1565b6001600160a01b036116326109b083866129a0565b167fdf1b1bd32a69711488d71554706bb130b1fc63a5fa1a2cd85e8440f84065ba23600080a2611617565b3461021c57604060031936011261021c5760043567ffffffffffffffff811161021c5761168e903690600401610b89565b6024359167ffffffffffffffff831161021c576116e76116df6116c56116bb6116ef963690600401610b89565b9490953691613117565b90604082015190605e604a84015160601c93015191929190565b919033613ee1565b81019061315e565b60005b8151811015610a865761173a61173561171c61170e84866129a0565b51516001600160a01b031690565b6001600160a01b03166000526007602052604060002090565b61321d565b61174e61174a6040830151151590565b1590565b6118af57906117996117666020600194015160ff1690565b611793611787602061177886896129a0565b5101516001600160e01b031690565b6001600160e01b031690565b90613fb3565b6117b460406117a884876129a0565b51015163ffffffff1690565b63ffffffff6117df6117d66117cf61099761170e888b6129a0565b5460e01c90565b63ffffffff1690565b9116106118a95761182d6117f860406117a885886129a0565b61181d61180361059e565b6001600160e01b03851681529163ffffffff166020830152565b61088a61099761170e86896129a0565b7f52f50aa6d1a95a4595361ecf953d095f125d442e4673716dede699e049de148a6001600160a01b0361186361170e85886129a0565b61189f61187560406117a8888b6129a0565b60405193849316958390929163ffffffff6020916001600160e01b03604085019616845216910152565b0390a25b016116f2565b506118a3565b6118f46118bf61170e84866129a0565b7f06439c6b000000000000000000000000000000000000000000000000000000006000526001600160a01b0316600452602490565b6000fd5b3461021c57604060031936011261021c5761066c611980611917610b46565b67ffffffffffffffff6024359161192d8361020b565b600060a060405161193d8161055f565b828152826020820152826040820152826060820152826080820152015216600052600a6020526040600020906001600160a01b0316600052602052604060002090565b6119fc6119f3604051926119938461055f565b5463ffffffff8116845263ffffffff8160201c16602085015261ffff8160401c1660408501526119da6119cd8263ffffffff9060501c1690565b63ffffffff166060860152565b63ffffffff607082901c16608085015260901c60ff1690565b151560a0830152565b6040519182918291909160a08060c083019463ffffffff815116845263ffffffff602082015116602085015261ffff604082015116604085015263ffffffff606082015116606085015263ffffffff608082015116608085015201511515910152565b60ff81160361021c57565b359061022c82611a5f565b3461021c57602060031936011261021c5760043567ffffffffffffffff811161021c573660238201121561021c57806004013590611ab2826109e1565b90611ac0604051928361057b565b82825260246102806020840194028201019036821161021c57602401925b818410611aee57610a86836132ad565b833603610280811261021c57610260601f1960405192611b0d84610543565b611b1688610b74565b8452011261021c5761028091602091611b2d6105ad565b611b38848901610a03565b8152611b466040890161114e565b84820152611b566060890161113d565b6040820152611b676080890161113d565b6060820152611b7860a0890161113d565b6080820152611b8960c08901611a6a565b60a0820152611b9a60e08901611a6a565b60c0820152611bac610100890161114e565b60e0820152611bbe610120890161113d565b610100820152611bd1610140890161114e565b610120820152611be4610160890161114e565b610140820152611bf7610180890161026d565b610160820152611c0a6101a08901610a03565b610180820152611c1d6101c0890161114e565b6101a0820152611c306101e0890161113d565b6101c0820152611c43610200890161113d565b6101e0820152611c566102208901610b74565b610200820152611c69610240890161113d565b610220820152611c7c610260890161113d565b61024082015283820152815201930192611ade565b3461021c57600060031936011261021c5760206001600160a01b0360015416604051908152f35b3461021c57602060031936011261021c5760043567ffffffffffffffff811161021c576040600319823603011261021c57604051611cf581610543565b816004013567ffffffffffffffff811161021c57611d1990600436918501016114bd565b8152602482013567ffffffffffffffff811161021c57610a86926004611d4292369201016114bd565b6020820152613463565b3461021c57602060031936011261021c5760043567ffffffffffffffff811161021c573660238201121561021c57806004013590611d89826109e1565b90611d97604051928361057b565b8282526024602083019360061b8201019036821161021c57602401925b818410611dc457610a86836135b5565b60408436031261021c5760206040918251611dde81610543565b8635611de98161020b565b8152611df6838801610b74565b83820152815201930192611db4565b3461021c57602060031936011261021c576001600160a01b03600435611e2a8161020b565b611e3261271b565b5016600052600760205261066c604060002060ff60405191611e5383610506565b546001600160a01b0381168352818160a01c16602084015260a81c16151560408201526040519182918291909160408060608301946001600160a01b03815116845260ff602082015116602085015201511515910152565b3461021c57600060031936011261021c57604051806020600b54918281520190600b6000527f0175b7a638427703f0dbe7bb9bbf987a2551717b34e79f33b5b1008d1fa01db99060005b818110611f0c5761066c856107148187038261057b565b8254845260209093019260019283019201611ef5565b3461021c57602060031936011261021c576040611f44600435610e198161020b565b611f698251809263ffffffff602080926001600160e01b038151168552015116910152565bf35b3461021c57600060031936011261021c57602060405160128152f35b3461021c57604060031936011261021c57611fa0610b46565b60243567ffffffffffffffff811161021c57806004019060a0600319823603011261021c57611fe661112c8467ffffffffffffffff166000526009602052604060002090565b606482019161201961174a611ffa85612ccb565b6001600160a01b03166000526001600b01602052604060002054151590565b61236b57908160448593019561202f878561273a565b959050602461204b8561204560848701896137ad565b906148f5565b93019061207a61205b83886137ad565b9050858961207361206c8b806137ad565b3691613117565b9289614aa9565b84612087610e8183612ccb565b9889946120a561209f61022085015163ffffffff1690565b82614b6b565b9b6000808c1561233157505061210b61ffff856121309961211799989661214b9661213e966121026120f26101c06120e66101a06121519f015161ffff1690565b97015163ffffffff1690565b916120fc8c612ccb565b9461273a565b96909516614c5c565b98919897909894612ccb565b6001600160a01b03166000526008602052604060002090565b5467ffffffffffffffff1690565b67ffffffffffffffff1690565b9061269b565b9560009761ffff61216861014089015161ffff1690565b166122d6575b509461214b61213e61020061223c61066c9d6dffffffffffffffffffffffffffff6122346122549f9e9b61222f6001600160e01b039f9b9c61224c9f61222f9e63ffffffff6121c361222f9f6121cd946137ad565b92905016906137fe565b908b60a081016121f06121ea6121e4835160ff1690565b60ff1690565b8561269b565b9360e0830191612202835161ffff1690565b9061ffff82168311612264575b505050506080015161222f916117d69163ffffffff1661383c565b61383c565b6137fe565b91169061269b565b93015167ffffffffffffffff1690565b9116906126ae565b6040519081529081906020820190565b6117d694965061222f959361ffff6122c56122b461222a966122ae6122a761229e60809960ff6122986122cc9b5160ff1690565b1661380b565b965161ffff1690565b61ffff1690565b90613684565b61214b6121e460c08d015160ff1690565b91166137fe565b959383955061220f565b909594989750826122fc8b989495986dffffffffffffffffffffffffffff9060701c1690565b6dffffffffffffffffffffffffffff169161231784896137ad565b90506123239388614e7e565b96979394389693929661216e565b969593509650505061214b61213e6121306121176123656123606117d661024061215199015163ffffffff1690565b612654565b94612ccb565b6118f461237784612ccb565b7f2502348c000000000000000000000000000000000000000000000000000000006000526001600160a01b0316600452602490565b3461021c57602060031936011261021c576001600160a01b036004356123d18161020b565b6123d9613a62565b1633811461243e57807fffffffffffffffffffffffff000000000000000000000000000000000000000060005416176000556001600160a01b03600154167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278600080a3005b7fdad89dca0000000000000000000000000000000000000000000000000000000060005260046000fd5b3461021c57602060031936011261021c5760043567ffffffffffffffff811161021c573660238201121561021c578060040135906124a5826109e1565b906124b3604051928361057b565b8282526024602083019360071b8201019036821161021c57602401925b8184106124e057610a8683613856565b8336036080811261021c576060601f19604051926124fd84610543565b87356125088161020b565b8452011261021c5760809160209160405161252281610506565b8388013561252f8161020b565b8152604088013561253f81611a5f565b848201526060880135612551816109f9565b6040820152838201528152019301926124d0565b3461021c57604060031936011261021c576004356125828161020b565b61258a610b5d565b9067ffffffffffffffff82169182600052600960205260ff60406000205416156125f7576125ba6125db926139a6565b92600052600960205263ffffffff60016040600020015460901c1690614b6b565b604080516001600160e01b039384168152919092166020820152f35b827f99ac52f20000000000000000000000000000000000000000000000000000000060005260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b90662386f26fc10000820291808304662386f26fc10000149015171561267657565b612625565b90655af3107a4000820291808304655af3107a4000149015171561267657565b8181029291811591840414171561267657565b81156126b8570490565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b61271161270b61271894936001600160e01b0361270481956139a6565b169061269b565b926139a6565b16906126ae565b90565b6040519061272882610506565b60006040838281528260208201520152565b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18136030182121561021c570180359067ffffffffffffffff821161021c57602001918160061b3603831361021c57565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b91908110156127cd5760061b0190565b61278e565b35906001600160e01b038216820361021c57565b60408136031261021c5761281560206040519261280284610543565b803561280d8161020b565b8452016127d2565b602082015290565b60408136031261021c5761281560206040519261283984610543565b61280d81610b74565b9061284b613a62565b60005b825181101561299b5780612864600192856129a0565b517f32a4ba3fa3351b11ad555d4c8ec70a744e8705607077a946807030d64b6ab1a360a06001600160a01b038351169260608101936001600160a01b0380865116957fffff00000000000000000000000000000000000000000000000000000000000061290260208601947fffffffffffffffffffff00000000000000000000000000000000000000000000865116604088019a848c511692614f23565b977fffffffffffffffffffff00000000000000000000000000000000000000000000608087019561296e875115158c600052600460205260406000209060ff7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0083541691151516179055565b8560405198511688525116602087015251166040850152511660608301525115156080820152a20161284e565b509050565b80518210156127cd5760209160051b010190565b9895979692949391907f0000000000000000000000000000000000000000000000000000000000000000906001600160a01b0382166001600160a01b03821614600014612ca2575050935b6bffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016808611612c71575090612a44918715158a613aa0565b979096612a91612a688367ffffffffffffffff166000526009602052604060002090565b5460081b7fffffffff000000000000000000000000000000000000000000000000000000001690565b90612a9b81613d04565b9760005b828110612ab3575050505050505093929190565b612ac6612ac18284896127bd565b612ccb565b8388612ae0612ad6858484613d4d565b60408101906137ad565b905060208111612bf6575b508392612b1a612b1461206c612b0a600198612b5597612b5097613d4d565b60208101906137ad565b896154bf565b612b388967ffffffffffffffff16600052600a602052604060002090565b906001600160a01b0316600052602052604060002090565b613253565b60a081015115612bba57612b9e612b766060612b9093015163ffffffff1690565b6040805163ffffffff909216602083015290928391820190565b03601f19810183528261057b565b612ba8828d6129a0565b52612bb3818c6129a0565b5001612a9f565b50612b90612b9e612bf184612be38a67ffffffffffffffff166000526009602052604060002090565b015460101c63ffffffff1690565b612b76565b915050612c2e6117d6612c2184612b388b67ffffffffffffffff16600052600a602052604060002090565b5460701c63ffffffff1690565b10612c3b57838838612aeb565b7f36f536ca000000000000000000000000000000000000000000000000000000006000526001600160a01b031660045260246000fd5b857f6a92a4830000000000000000000000000000000000000000000000000000000060005260045260245260446000fd5b91612cac926126e7565b936129ff565b60405190612cbf82610543565b60006020838281520152565b356127188161020b565b90604051612ce281610543565b91546001600160e01b038116835260e01c6020830152565b9061022c612e5a6001612d0b6105ad565b94612df9612def8254612d27612d218260ff1690565b15158a52565b61ffff600882901c1660208a015263ffffffff601882901c1660408a015263ffffffff603882901c1660608a015263ffffffff605882901c1660808a015260ff607882901c1660a08a015260ff608082901c1660c08a015261ffff608882901c1660e08a015263ffffffff609882901c166101008a015261ffff60b882901c166101208a015261ffff60c882901c166101408a01527fffffffff00000000000000000000000000000000000000000000000000000000600882901b166101608a015260f81c90565b1515610180880152565b015461ffff81166101a086015263ffffffff601082901c166101c086015263ffffffff603082901c166101e086015267ffffffffffffffff605082901c1661020086015263ffffffff609082901c1661022086015260b01c63ffffffff1690565b63ffffffff16610240840152565b90612e71613a62565b6000915b805183101561306357612e8883826129a0565b5190612e9c825167ffffffffffffffff1690565b946020600093019367ffffffffffffffff8716935b8551805182101561304e57612ec8826020926129a0565b510151612ed961170e8389516129a0565b8151602083015163ffffffff908116911681811015613015575050608082015163ffffffff1660208110612fd4575090867f94967ae9ea7729ad4f54021c1981765d2b1d954f7c92fbec340aa0a54f46b8b56001600160a01b0384612f63858f60019998612b38612f5e9267ffffffffffffffff16600052600a602052604060002090565b613d8d565b612fcb60405192839216958291909160a08060c083019463ffffffff815116845263ffffffff602082015116602085015261ffff604082015116604085015263ffffffff606082015116606085015263ffffffff608082015116608085015201511515910152565b0390a301612eb1565b7f24ecdc02000000000000000000000000000000000000000000000000000000006000526001600160a01b0390911660045263ffffffff1660245260446000fd5b7f0b4f67a20000000000000000000000000000000000000000000000000000000060005263ffffffff9081166004521660245260446000fd5b50509550925092600191500191929092612e75565b50905060005b81518110156131135780613091613082600193856129a0565b515167ffffffffffffffff1690565b67ffffffffffffffff6001600160a01b036130c060206130b186896129a0565b5101516001600160a01b031690565b60006130e482612b388767ffffffffffffffff16600052600a602052604060002090565b551691167f4de5b1bcbca6018c11303a2c3f4a4b4f22a1c741d8c4ba430d246ac06c5ddf8b600080a301613069565b5050565b92919267ffffffffffffffff82116105225760405191613141601f8201601f19166020018461057b565b82948184528183011161021c578281602093846000960137010152565b60208183031261021c5780359067ffffffffffffffff821161021c570181601f8201121561021c57803590613192826109e1565b926131a0604051948561057b565b8284526020606081860194028301019181831161021c57602001925b8284106131ca575050505090565b60608483031261021c5760206060916040516131e581610506565b86356131f08161020b565b81526131fd8388016127d2565b8382015261320d6040880161113d565b60408201528152019301926131bc565b9060405161322a81610506565b604060ff8294546001600160a01b0381168452818160a01c16602085015260a81c161515910152565b9061022c6040516132638161055f565b925463ffffffff8082168552602082811c821690860152604082811c61ffff1690860152605082901c81166060860152607082901c16608085015260901c60ff16151560a0840152565b906132b6613a62565b60005b825181101561299b576132cc81846129a0565b5160206132dc61308284876129a0565b9101519067ffffffffffffffff811680158015613444575b8015613416575b6133de57916133a4826001959461335461332f612a686133a99767ffffffffffffffff166000526009602052604060002090565b7fffffffff000000000000000000000000000000000000000000000000000000001690565b6133af577f71e9302ab4e912a9678ae7f5a8542856706806f2817e1bf2a20b171e265cb4ad604051806133878782610f1e565b0390a267ffffffffffffffff166000526009602052604060002090565b6140c0565b016132b9565b7f2431cc0363f2f66b21782c7e3d54dd9085927981a21bd0cc6be45a51b19689e3604051806133878782610f1e565b7fc35aa79d0000000000000000000000000000000000000000000000000000000060005267ffffffffffffffff821660045260246000fd5b506101e083015163ffffffff1663ffffffff61343c6117d6606087015163ffffffff1690565b9116116132fb565b5063ffffffff61345c6101e085015163ffffffff1690565b16156132f4565b61346b613a62565b60208101519160005b83518110156134f8578061348d6109b0600193876129a0565b6134af6134aa6001600160a01b0383165b6001600160a01b031690565b615a59565b6134bb575b5001613474565b6040516001600160a01b039190911681527fc3803387881faad271c47728894e3e36fac830ffc8602ca6fc07733cbda7758090602090a1386134b4565b5091505160005b8151811015613113576135156109b082846129a0565b906001600160a01b0382161561358b577feb1b9b92e50b7f88f9ff25d56765095ac6e91540eee214906f4036a908ffbdef6135828361356761356261349e6001976001600160a01b031690565b6159e0565b506040516001600160a01b0390911681529081906020820190565b0390a1016134ff565b7f8579befe0000000000000000000000000000000000000000000000000000000060005260046000fd5b6135bd613a62565b60005b815181101561311357806001600160a01b036135de600193856129a0565b5151167fbb77da6f7210cdd16904228a9360133d1d7dfff99b1bc75f128da5b53e28f97d61367b67ffffffffffffffff602061361a86896129a0565b51015116836000526008602052604060002067ffffffffffffffff82167fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000008254161790556040519182918291909167ffffffffffffffff6020820193169052565b0390a2016135c0565b9190820391821161267657565b613699612cb2565b506136bf6136ba826001600160a01b03166000526006602052604060002090565b612cd5565b60208101916136de6136d86117d6855163ffffffff1690565b42613684565b63ffffffff7f000000000000000000000000000000000000000000000000000000000000000016116137865761173561372a916001600160a01b03166000526007602052604060002090565b61373a61174a6040830151151590565b801561378c575b6137865761374e9061478e565b9163ffffffff6137766117d661376b602087015163ffffffff1690565b935163ffffffff1690565b911610613781575090565b905090565b50905090565b506001600160a01b036137a682516001600160a01b031690565b1615613741565b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18136030182121561021c570180359067ffffffffffffffff821161021c5760200191813603831361021c57565b9190820180921161267657565b9061ffff8091169116029061ffff821691820361267657565b63ffffffff60209116019063ffffffff821161267657565b9063ffffffff8091169116019063ffffffff821161267657565b9061385f613a62565b60005b825181101561299b5780613878600192856129a0565b517fe6a7a17d710bf0b2cd05e5397dc6f97a5da4ee79e31e234bf5f965ee2bd9a5bf61399d60206001600160a01b0384511693015183600052600760205260406000206138fd6001600160a01b0383511682906001600160a01b03167fffffffffffffffffffffffff0000000000000000000000000000000000000000825416179055565b602082015181547fffffffffffffffffffff0000ffffffffffffffffffffffffffffffffffffffff74ff000000000000000000000000000000000000000075ff0000000000000000000000000000000000000000006040870151151560a81b169360a01b169116171790556040519182918291909160408060608301946001600160a01b03815116845260ff602082015116602085015201511515910152565b0390a201613862565b6139af81613691565b9063ffffffff602083015116158015613a0c575b6139d55750516001600160e01b031690565b6001600160a01b03907f06439c6b000000000000000000000000000000000000000000000000000000006000521660045260246000fd5b506001600160e01b03825116156139c3565b33600052600360205260406000205415613a3457565b7fd86ad9cf000000000000000000000000000000000000000000000000000000006000523360045260246000fd5b6001600160a01b03600154163303613a7657565b7f2b5c74de0000000000000000000000000000000000000000000000000000000060005260046000fd5b61112c613ac691959493929567ffffffffffffffff166000526009602052604060002090565b936101608501947f2812d52c000000000000000000000000000000000000000000000000000000007fffffffff00000000000000000000000000000000000000000000000000000000613b3988517fffffffff000000000000000000000000000000000000000000000000000000001690565b1614613cc8577f1e10bdc4000000000000000000000000000000000000000000000000000000007fffffffff00000000000000000000000000000000000000000000000000000000613bab88517fffffffff000000000000000000000000000000000000000000000000000000001690565b1614613c2c576118f4613bde87517fffffffff000000000000000000000000000000000000000000000000000000001690565b7f2ee82075000000000000000000000000000000000000000000000000000000006000527fffffffff0000000000000000000000000000000000000000000000000000000016600452602490565b606091929394955080613c4a6101e0613c7793015163ffffffff1690565b63ffffffff613c6d610180613c658786015163ffffffff1690565b940151151590565b921690878761530b565b015181613cbf575b50613c9557613c8f913691613117565b90600190565b7f5bed51920000000000000000000000000000000000000000000000000000000060005260046000fd5b90501538613c7f565b6101e00151939450613ced93919291613ce7915063ffffffff166117d6565b91615070565b906127186020613cfc84615191565b930151151590565b90613d0e826109e1565b613d1b604051918261057b565b828152601f19613d2b82946109e1565b019060005b828110613d3c57505050565b806060602080938501015201613d30565b91908110156127cd5760051b810135907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff618136030182121561021c570190565b815181546020808501516040808701517fffffffffffffffffffffffffffffffffffffffffffff0000000000000000000090941663ffffffff958616179190921b67ffffffff00000000161791901b69ffff000000000000000016178255606083015161022c93613e9d9260a092613e3f911685547fffffffffffffffffffffffffffffffffffff00000000ffffffffffffffffffff1660509190911b6dffffffff0000000000000000000016178555565b613e96613e53608083015163ffffffff1690565b85547fffffffffffffffffffffffffffff00000000ffffffffffffffffffffffffffff1660709190911b71ffffffff000000000000000000000000000016178555565b0151151590565b81547fffffffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffff1690151560901b72ff00000000000000000000000000000000000016179055565b91929092613ef182828686614f23565b600052600460205260ff6040600020541615613f0d5750505050565b6040517f097e17ff0000000000000000000000000000000000000000000000000000000081526001600160a01b0393841660048201529390921660248401527fffffffffffffffffffff0000000000000000000000000000000000000000000090911660448301527fffff000000000000000000000000000000000000000000000000000000000000166064820152608490fd5b0390fd5b604d811161267657600a0a90565b60ff1660120160ff81116126765760ff1690602482111561404e577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdc82019182116126765761400461400a92613fa5565b906126ae565b6001600160e01b038111614024576001600160e01b031690565b7f10cb51d10000000000000000000000000000000000000000000000000000000060005260046000fd5b9060240390602482116126765761214b61406792613fa5565b61400a565b9060ff80911691160160ff81116126765760ff1690602482111561404e577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdc82019182116126765761400461400a92613fa5565b906146da610240600161022c9461410b6140da8651151590565b829060ff7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0083541691151516179055565b61415161411d602087015161ffff1690565b82547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000ff1660089190911b62ffff0016178255565b61419d614165604087015163ffffffff1690565b82547fffffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffff1660189190911b66ffffffff00000016178255565b6141ed6141b1606087015163ffffffff1690565b82547fffffffffffffffffffffffffffffffffffffffffff00000000ffffffffffffff1660389190911b6affffffff0000000000000016178255565b614241614201608087015163ffffffff1690565b82547fffffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffff1660589190911b6effffffff000000000000000000000016178255565b61429361425260a087015160ff1690565b82547fffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffff1660789190911b6fff00000000000000000000000000000016178255565b6142e66142a460c087015160ff1690565b82547fffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffffff1660809190911b70ff0000000000000000000000000000000016178255565b61433c6142f860e087015161ffff1690565b82547fffffffffffffffffffffffffff0000ffffffffffffffffffffffffffffffffff1660889190911b72ffff000000000000000000000000000000000016178255565b61439961435161010087015163ffffffff1690565b82547fffffffffffffffffff00000000ffffffffffffffffffffffffffffffffffffff1660989190911b76ffffffff0000000000000000000000000000000000000016178255565b6143f66143ac61012087015161ffff1690565b82547fffffffffffffff0000ffffffffffffffffffffffffffffffffffffffffffffff1660b89190911b78ffff000000000000000000000000000000000000000000000016178255565b61445561440961014087015161ffff1690565b82547fffffffffff0000ffffffffffffffffffffffffffffffffffffffffffffffffff1660c89190911b7affff0000000000000000000000000000000000000000000000000016178255565b6144d66144866101608701517fffffffff000000000000000000000000000000000000000000000000000000001690565b82547fff00000000ffffffffffffffffffffffffffffffffffffffffffffffffffffff1660089190911c7effffffff00000000000000000000000000000000000000000000000000000016178255565b6145376144e7610180870151151590565b82547effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1690151560f81b7fff0000000000000000000000000000000000000000000000000000000000000016178255565b019261457b61454c6101a083015161ffff1690565b859061ffff167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000825416179055565b6145c76145906101c083015163ffffffff1690565b85547fffffffffffffffffffffffffffffffffffffffffffffffffffff00000000ffff1660109190911b65ffffffff000016178555565b6146176145dc6101e083015163ffffffff1690565b85547fffffffffffffffffffffffffffffffffffffffffffff00000000ffffffffffff1660309190911b69ffffffff00000000000016178555565b61467361463061020083015167ffffffffffffffff1690565b85547fffffffffffffffffffffffffffff0000000000000000ffffffffffffffffffff1660509190911b71ffffffffffffffff0000000000000000000016178555565b6146cf61468861022083015163ffffffff1690565b85547fffffffffffffffffffff00000000ffffffffffffffffffffffffffffffffffff1660909190911b75ffffffff00000000000000000000000000000000000016178555565b015163ffffffff1690565b7fffffffffffff00000000ffffffffffffffffffffffffffffffffffffffffffff79ffffffff0000000000000000000000000000000000000000000083549260b01b169116179055565b519069ffffffffffffffffffff8216820361021c57565b908160a091031261021c5761474f81614724565b91602082015191604081015191612718608060608401519301614724565b6040513d6000823e3d90fd5b9081602091031261021c575161271881611a5f565b614796612cb2565b506147ae61349e61349e83516001600160a01b031690565b90604051907ffeaf968c00000000000000000000000000000000000000000000000000000000825260a082600481865afa9283156148b6576000926000946148bb575b5060008312614024576020600491604051928380927f313ce5670000000000000000000000000000000000000000000000000000000082525afa9283156148b6576127189363ffffffff9361485793600092614880575b506020015160ff165b9061406c565b9261487261486361059e565b6001600160e01b039095168552565b1663ffffffff166020830152565b6148519192506148a7602091823d84116148af575b61489f818361057b565b810190614779565b929150614848565b503d614895565b61476d565b9093506148e191925060a03d60a0116148ee575b6148d9818361057b565b81019061473b565b50939250509192386147f1565b503d6148cf565b929190926101608201937f2812d52c000000000000000000000000000000000000000000000000000000007fffffffff0000000000000000000000000000000000000000000000000000000061496b87517fffffffff000000000000000000000000000000000000000000000000000000001690565b1614614a68577f1e10bdc4000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000006149dd87517fffffffff000000000000000000000000000000000000000000000000000000001690565b1614614a10576118f4613bde86517fffffffff000000000000000000000000000000000000000000000000000000001690565b61271893945091614a5e916117d693614a316101e084015163ffffffff1690565b9063ffffffff614a56610180614a4e606088015163ffffffff1690565b960151151590565b94169261530b565b5163ffffffff1690565b614aa5939450614a806101e084015163ffffffff1690565b9063ffffffff614a9d610180614a4e606088015163ffffffff1690565b94169261571e565b5190565b9493919063ffffffff604087015116808211614b3b57505061ffff60208601511690818111614b055750507fffffffff0000000000000000000000000000000000000000000000000000000061016061022c9495015116615608565b61ffff92507fd88dddd6000000000000000000000000000000000000000000000000000000006000526004521660245260446000fd5b7f869337890000000000000000000000000000000000000000000000000000000060005260045260245260446000fd5b67ffffffffffffffff8116600052600560205260406000209160405192614b9184610543565b546001600160e01b038116845260e01c9182602085015263ffffffff82169283614bcb575b5050505061271890516001600160e01b031690565b63ffffffff164290810393908411612676578311614be95780614bb6565b7ff08bcb3e0000000000000000000000000000000000000000000000000000000060005267ffffffffffffffff1660045263ffffffff1660245260445260646000fd5b60408136031261021c57602060405191614c4583610543565b8035614c508161020b565b83520135602082015290565b9694919695929390956000946000986000986000965b808810614c86575050505050505050929190565b9091929394959697999a614ca3614c9e8a848b6127bd565b614c2c565b9a614cf1612b508d614cda614ccc8967ffffffffffffffff16600052600a602052604060002090565b91516001600160a01b031690565b6001600160a01b0316600052602052604060002090565b91614d0261174a60a0850151151590565b614e4b5760009c6040840190614d1d6122a7835161ffff1690565b614dd3575b5050606083015163ffffffff16614d389161383c565b9c6080830151614d4b9063ffffffff1690565b614d549161383c565b9b8251614d649063ffffffff1690565b63ffffffff16614d7390612654565b60019390808310614dc757506123606117d66020614d9693015163ffffffff1690565b808211614db65750614da7916137fe565b985b0196959493929190614c72565b9050614dc1916137fe565b98614da9565b915050614dc1916137fe565b9061214b614e3c939f614e2a614e339460208f8e6122a795506001600160a01b03614e0585516001600160a01b031690565b91166001600160a01b03821614614e4457614e2091506139a6565b915b015190615786565b925161ffff1690565b620186a0900490565b9b3880614d22565b5091614e22565b999b5060019150614e7284614e6c614e7893614e668b612654565b906137fe565b9b61383c565b9c613824565b9a614da9565b91939093806101e00193846101e011612676576101208102908082046101201490151715612676576101e0910101809311612676576122a7610140614f14612718966dffffffffffffffffffffffffffff612234614eff614eec614f1e9a63ffffffff61214b9a16906137fe565b61214b6122a76101208c015161ffff1690565b614e666117d66101008b015163ffffffff1690565b93015161ffff1690565b61267b565b604080516001600160a01b039283166020820190815292909316908301527fffffffffffffffffffff0000000000000000000000000000000000000000000090921660608201527fffff000000000000000000000000000000000000000000000000000000000000909216608083015290614fa18160a08101612b90565b51902090565b919091357fffffffff0000000000000000000000000000000000000000000000000000000081169260048110614fdb575050565b7fffffffff00000000000000000000000000000000000000000000000000000000929350829060040360031b1b161690565b909291928360041161021c57831161021c57600401916003190190565b9060041161021c5790600490565b9081602091031261021c575190565b9081604091031261021c5760206040519161506183610543565b805183520151612815816109f9565b91615079612cb2565b50811561516f57506150ba61206c82806150b47fffffffff000000000000000000000000000000000000000000000000000000009587614fa7565b9561500d565b91167f181dcf100000000000000000000000000000000000000000000000000000000081036150f757508060208061271893518301019101615047565b7f97a657c90000000000000000000000000000000000000000000000000000000014615147577f5247fdce0000000000000000000000000000000000000000000000000000000060005260046000fd5b8060208061515a93518301019101615038565b61516261059e565b9081526000602082015290565b91505067ffffffffffffffff61518361059e565b911681526000602082015290565b6020604051917f181dcf100000000000000000000000000000000000000000000000000000000082840152805160248401520151151560448201526044815261271860648261057b565b604051906151e882610527565b60606080836000815260006020820152600060408201526000838201520152565b6040519061521860208361057b565b6000808352366020840137565b60208183031261021c5780359067ffffffffffffffff821161021c57019060a08282031261021c576040519161525a83610527565b6152638161113d565b835261527160208201610b74565b60208401526040810135615284816109f9565b60408401526060810135606084015260808101359067ffffffffffffffff821161021c57019080601f8301121561021c5781356152c0816109e1565b926152ce604051948561057b565b81845260208085019260051b82010192831161021c57602001905b8282106152fb57505050608082015290565b81358152602091820191016152e9565b90916153156151db565b5060048310615495577f1f3b3aba000000000000000000000000000000000000000000000000000000007fffffffff0000000000000000000000000000000000000000000000000000000061537361536d868661502a565b90614fa7565b160361546b57600483111561542a575081615399926153919261500d565b810190615225565b9180615414575b6153ea5763ffffffff6153b7835163ffffffff1690565b16116153c05790565b7f2e2b0c290000000000000000000000000000000000000000000000000000000060005260046000fd5b7fee433e990000000000000000000000000000000000000000000000000000000060005260046000fd5b5061542561174a6040840151151590565b6153a0565b9350505050615437615209565b61544e6154426105bd565b63ffffffff9093168352565b600060208301526001604083015260006060830152608082015290565b7f5247fdce0000000000000000000000000000000000000000000000000000000060005260046000fd5b7fb00b53dc0000000000000000000000000000000000000000000000000000000060005260046000fd5b907fffffffff0000000000000000000000000000000000000000000000000000000082167f2812d52c00000000000000000000000000000000000000000000000000000000811461558d577f1e10bdc40000000000000000000000000000000000000000000000000000000014615580577f2ee82075000000000000000000000000000000000000000000000000000000006000527fffffffff00000000000000000000000000000000000000000000000000000000821660045260246000fd5b61022c91506001906157b7565b50905060208151036155c6576155ac6020825183010160208301615038565b6001600160a01b0381119081156155fc575b506155c65750565b613fa1906040519182917f8d666f60000000000000000000000000000000000000000000000000000000008352600483016157a6565b610400915010386155be565b917fffffffff0000000000000000000000000000000000000000000000000000000083167f2812d52c0000000000000000000000000000000000000000000000000000000081146156d6577f1e10bdc400000000000000000000000000000000000000000000000000000000146156c9577f2ee82075000000000000000000000000000000000000000000000000000000006000527fffffffff00000000000000000000000000000000000000000000000000000000831660045260246000fd5b61022c92501515906157b7565b5050905060208151036155c6576155ac6020825183010160208301615038565b6001600160a01b036127189116600b6158e7565b6001600160a01b036127189116600b615a1b565b9063ffffffff61573b93959495615733612cb2565b501691615070565b9182511161575c5780615750575b6153ea5790565b50602081015115615749565b7f4c4fc93a0000000000000000000000000000000000000000000000000000000060005260046000fd5b670de0b6b3a7640000916001600160e01b036157a2921661269b565b0490565b9060206127189281815201906105cc565b90602082510361581f576157c85750565b60208180518101031261021c576020810151156157e25750565b613fa1906040519182917fff828faa00000000000000000000000000000000000000000000000000000000835260206004840181815201906105cc565b6040517fff828faa0000000000000000000000000000000000000000000000000000000081526020600482015280613fa160248201856105cc565b80548210156127cd5760005260206000200190600090565b9161588c918354906000199060031b92831b921b19161790565b9055565b805480156158b85760001901906158a7828261585a565b60001982549160031b1b1916905555565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b600181019180600052826020526040600020549283151560001461599957600019840184811161267657835493600019850194851161267657600095858361594a9761593b9503615950575b505050615890565b90600052602052604060002090565b55600190565b61598061597a91615971615967615990958861585a565b90549060031b1c90565b9283918761585a565b90615872565b8590600052602052604060002090565b55388080615933565b50505050600090565b8054906801000000000000000082101561052257816159c991600161588c9401815561585a565b81939154906000199060031b92831b921b19161790565b600081815260036020526040902054615a15576159fe8160026159a2565b600254906000526003602052604060002055600190565b50600090565b6000828152600182016020526040902054615a525780615a3d836001936159a2565b80549260005201602052604060002055600190565b5050600090565b600081815260036020526040902054908115615a52576000198201908282116126765760025492600019840193841161267657838361594a9460009603615ab9575b505050615aa86002615890565b600390600052602052604060002090565b615aa861597a91615ad1615967615adb95600261585a565b928391600261585a565b55388080615a9b56fea164736f6c634300081a000a",
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
	return common.HexToHash("0x71e9302ab4e912a9678ae7f5a8542856706806f2817e1bf2a20b171e265cb4ad")
}

func (FeeQuoterDestChainConfigUpdated) Topic() common.Hash {
	return common.HexToHash("0x2431cc0363f2f66b21782c7e3d54dd9085927981a21bd0cc6be45a51b19689e3")
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
