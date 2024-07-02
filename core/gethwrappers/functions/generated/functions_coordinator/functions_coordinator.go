// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package functions_coordinator

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

type FunctionsBillingConfig struct {
	FulfillmentGasPriceOverEstimationBP uint32
	FeedStalenessSeconds                uint32
	GasOverheadBeforeCallback           uint32
	GasOverheadAfterCallback            uint32
	MinimumEstimateGasPriceWei          *big.Int
	MaxSupportedRequestDataVersion      uint16
	FallbackUsdPerUnitLink              uint64
	FallbackUsdPerUnitLinkDecimals      uint8
	FallbackNativePerUnitLink           *big.Int
	RequestTimeoutSeconds               uint32
	DonFeeCentsUsd                      uint16
	OperationFeeCentsUsd                uint16
	TransmitTxSizeBytes                 uint16
}

type FunctionsResponseCommitment struct {
	RequestId                 [32]byte
	Coordinator               common.Address
	EstimatedTotalCostJuels   *big.Int
	Client                    common.Address
	SubscriptionId            uint64
	CallbackGasLimit          uint32
	AdminFee                  *big.Int
	DonFee                    *big.Int
	GasOverheadBeforeCallback *big.Int
	GasOverheadAfterCallback  *big.Int
	TimeoutTimestamp          uint32
}

type FunctionsResponseRequestMeta struct {
	Data               []byte
	Flags              [32]byte
	RequestingContract common.Address
	AvailableBalance   *big.Int
	AdminFee           *big.Int
	SubscriptionId     uint64
	InitiatedRequests  uint64
	CallbackGasLimit   uint32
	DataVersion        uint16
	CompletedRequests  uint64
	SubscriptionOwner  common.Address
}

var FunctionsCoordinatorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"fulfillmentGasPriceOverEstimationBP\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"feedStalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasOverheadBeforeCallback\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasOverheadAfterCallback\",\"type\":\"uint32\"},{\"internalType\":\"uint40\",\"name\":\"minimumEstimateGasPriceWei\",\"type\":\"uint40\"},{\"internalType\":\"uint16\",\"name\":\"maxSupportedRequestDataVersion\",\"type\":\"uint16\"},{\"internalType\":\"uint64\",\"name\":\"fallbackUsdPerUnitLink\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"fallbackUsdPerUnitLinkDecimals\",\"type\":\"uint8\"},{\"internalType\":\"uint224\",\"name\":\"fallbackNativePerUnitLink\",\"type\":\"uint224\"},{\"internalType\":\"uint32\",\"name\":\"requestTimeoutSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"donFeeCentsUsd\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"operationFeeCentsUsd\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"transmitTxSizeBytes\",\"type\":\"uint16\"}],\"internalType\":\"structFunctionsBillingConfig\",\"name\":\"config\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"linkToNativeFeed\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"linkToUsdFeed\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"EmptyPublicKey\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InconsistentReportData\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"InvalidConfig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"linkWei\",\"type\":\"int256\"}],\"name\":\"InvalidLinkWeiPrice\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSubscription\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"usdLink\",\"type\":\"int256\"}],\"name\":\"InvalidUsdLinkPrice\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"MustBeSubOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoTransmittersSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByRouter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByRouterOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentTooLarge\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"ReportInvalid\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RouterMustBeSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedPublicKeyChange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedSender\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnsupportedRequestDataVersion\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"CommitmentDeleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"fulfillmentGasPriceOverEstimationBP\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"feedStalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasOverheadBeforeCallback\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasOverheadAfterCallback\",\"type\":\"uint32\"},{\"internalType\":\"uint40\",\"name\":\"minimumEstimateGasPriceWei\",\"type\":\"uint40\"},{\"internalType\":\"uint16\",\"name\":\"maxSupportedRequestDataVersion\",\"type\":\"uint16\"},{\"internalType\":\"uint64\",\"name\":\"fallbackUsdPerUnitLink\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"fallbackUsdPerUnitLinkDecimals\",\"type\":\"uint8\"},{\"internalType\":\"uint224\",\"name\":\"fallbackNativePerUnitLink\",\"type\":\"uint224\"},{\"internalType\":\"uint32\",\"name\":\"requestTimeoutSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"donFeeCentsUsd\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"operationFeeCentsUsd\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"transmitTxSizeBytes\",\"type\":\"uint16\"}],\"indexed\":false,\"internalType\":\"structFunctionsBillingConfig\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"ConfigUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requestingContract\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requestInitiator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"subscriptionOwner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"dataVersion\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"flags\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"callbackGasLimit\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"estimatedTotalCostJuels\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint72\",\"name\":\"adminFee\",\"type\":\"uint72\"},{\"internalType\":\"uint72\",\"name\":\"donFee\",\"type\":\"uint72\"},{\"internalType\":\"uint40\",\"name\":\"gasOverheadBeforeCallback\",\"type\":\"uint40\"},{\"internalType\":\"uint40\",\"name\":\"gasOverheadAfterCallback\",\"type\":\"uint40\"},{\"internalType\":\"uint32\",\"name\":\"timeoutTimestamp\",\"type\":\"uint32\"}],\"indexed\":false,\"internalType\":\"structFunctionsResponse.Commitment\",\"name\":\"commitment\",\"type\":\"tuple\"}],\"name\":\"OracleRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"OracleResponse\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"juelsPerGas\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"l1FeeShareWei\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"callbackCostJuels\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint72\",\"name\":\"donFeeJuels\",\"type\":\"uint72\"},{\"indexed\":false,\"internalType\":\"uint72\",\"name\":\"adminFeeJuels\",\"type\":\"uint72\"},{\"indexed\":false,\"internalType\":\"uint72\",\"name\":\"operationFeeJuels\",\"type\":\"uint72\"}],\"name\":\"RequestBilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"deleteCommitment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasPriceWei\",\"type\":\"uint256\"}],\"name\":\"estimateCost\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAdminFeeJuels\",\"outputs\":[{\"internalType\":\"uint72\",\"name\":\"\",\"type\":\"uint72\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"fulfillmentGasPriceOverEstimationBP\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"feedStalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasOverheadBeforeCallback\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasOverheadAfterCallback\",\"type\":\"uint32\"},{\"internalType\":\"uint40\",\"name\":\"minimumEstimateGasPriceWei\",\"type\":\"uint40\"},{\"internalType\":\"uint16\",\"name\":\"maxSupportedRequestDataVersion\",\"type\":\"uint16\"},{\"internalType\":\"uint64\",\"name\":\"fallbackUsdPerUnitLink\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"fallbackUsdPerUnitLinkDecimals\",\"type\":\"uint8\"},{\"internalType\":\"uint224\",\"name\":\"fallbackNativePerUnitLink\",\"type\":\"uint224\"},{\"internalType\":\"uint32\",\"name\":\"requestTimeoutSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"donFeeCentsUsd\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"operationFeeCentsUsd\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"transmitTxSizeBytes\",\"type\":\"uint16\"}],\"internalType\":\"structFunctionsBillingConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"getDONFeeJuels\",\"outputs\":[{\"internalType\":\"uint72\",\"name\":\"\",\"type\":\"uint72\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDONPublicKey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOperationFeeJuels\",\"outputs\":[{\"internalType\":\"uint72\",\"name\":\"\",\"type\":\"uint72\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getThresholdPublicKey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getUsdPerUnitLink\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getWeiPerUnitLink\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"oracleWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"oracleWithdrawAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"_f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"_onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"_offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"_offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"donPublicKey\",\"type\":\"bytes\"}],\"name\":\"setDONPublicKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"thresholdPublicKey\",\"type\":\"bytes\"}],\"name\":\"setThresholdPublicKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"flags\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"requestingContract\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"availableBalance\",\"type\":\"uint96\"},{\"internalType\":\"uint72\",\"name\":\"adminFee\",\"type\":\"uint72\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"initiatedRequests\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"dataVersion\",\"type\":\"uint16\"},{\"internalType\":\"uint64\",\"name\":\"completedRequests\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"subscriptionOwner\",\"type\":\"address\"}],\"internalType\":\"structFunctionsResponse.RequestMeta\",\"name\":\"request\",\"type\":\"tuple\"}],\"name\":\"startRequest\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"estimatedTotalCostJuels\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint72\",\"name\":\"adminFee\",\"type\":\"uint72\"},{\"internalType\":\"uint72\",\"name\":\"donFee\",\"type\":\"uint72\"},{\"internalType\":\"uint40\",\"name\":\"gasOverheadBeforeCallback\",\"type\":\"uint40\"},{\"internalType\":\"uint40\",\"name\":\"gasOverheadAfterCallback\",\"type\":\"uint40\"},{\"internalType\":\"uint32\",\"name\":\"timeoutTimestamp\",\"type\":\"uint32\"}],\"internalType\":\"structFunctionsResponse.Commitment\",\"name\":\"commitment\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"transmitters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"fulfillmentGasPriceOverEstimationBP\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"feedStalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasOverheadBeforeCallback\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasOverheadAfterCallback\",\"type\":\"uint32\"},{\"internalType\":\"uint40\",\"name\":\"minimumEstimateGasPriceWei\",\"type\":\"uint40\"},{\"internalType\":\"uint16\",\"name\":\"maxSupportedRequestDataVersion\",\"type\":\"uint16\"},{\"internalType\":\"uint64\",\"name\":\"fallbackUsdPerUnitLink\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"fallbackUsdPerUnitLinkDecimals\",\"type\":\"uint8\"},{\"internalType\":\"uint224\",\"name\":\"fallbackNativePerUnitLink\",\"type\":\"uint224\"},{\"internalType\":\"uint32\",\"name\":\"requestTimeoutSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"donFeeCentsUsd\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"operationFeeCentsUsd\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"transmitTxSizeBytes\",\"type\":\"uint16\"}],\"internalType\":\"structFunctionsBillingConfig\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"updateConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b5060405162006047380380620060478339810160408190526200003491620004f2565b83838383833380600081620000905760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c357620000c3816200014e565b5050506001600160a01b038116620000ee57604051632530e88560e11b815260040160405180910390fd5b6001600160a01b03908116608052600c80546001600160601b03166c0100000000000000000000000085841602179055600d80546001600160a01b0319169183169190911790556200014083620001f9565b505050505050505062000779565b336001600160a01b03821603620001a85760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000087565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b62000203620003b5565b80516008805460208401516040808601516060870151608088015160a089015160c08a015160e08b015163ffffffff9a8b166001600160401b031990991698909817640100000000978b16880217600160401b600160801b03191668010000000000000000958b169590950263ffffffff60601b1916949094176c01000000000000000000000000938a16939093029290921766ffffffffffffff60801b1916600160801b64ffffffffff9092169190910261ffff60a81b191617600160a81b61ffff92831602176001600160b81b0316600160b81b6001600160401b03909316929092026001600160f81b031691909117600160f81b60ff90951694909402939093179093556101008501516101208601516001600160e01b03909116600160e01b919095160293909317600955610140840151600a805461016087015161018088015193851663ffffffff199092169190911762010000918516919091021761ffff60201b19169190921690930292909217909155517f9832846a11905a128363408ace5f9803f449fcb017d3cf7e22a50c6e332931b290620003aa90839062000663565b60405180910390a150565b620003bf620003c1565b565b6000546001600160a01b03163314620003bf5760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640162000087565b80516001600160a01b03811681146200043557600080fd5b919050565b6040516101a081016001600160401b03811182821017156200046c57634e487b7160e01b600052604160045260246000fd5b60405290565b805163ffffffff811681146200043557600080fd5b805164ffffffffff811681146200043557600080fd5b805161ffff811681146200043557600080fd5b80516001600160401b03811681146200043557600080fd5b805160ff811681146200043557600080fd5b80516001600160e01b03811681146200043557600080fd5b6000806000808486036102008112156200050b57600080fd5b62000516866200041d565b94506101a080601f19830112156200052d57600080fd5b620005376200043a565b9150620005476020880162000472565b8252620005576040880162000472565b60208301526200056a6060880162000472565b60408301526200057d6080880162000472565b60608301526200059060a0880162000487565b6080830152620005a360c088016200049d565b60a0830152620005b660e08801620004b0565b60c0830152610100620005cb818901620004c8565b60e0840152610120620005e0818a01620004da565b828501526101409150620005f6828a0162000472565b908401526101606200060a8982016200049d565b82850152610180915062000620828a016200049d565b90840152620006318883016200049d565b90830152509250620006476101c086016200041d565b9150620006586101e086016200041d565b905092959194509250565b815163ffffffff1681526101a0810160208301516200068a602084018263ffffffff169052565b506040830151620006a3604084018263ffffffff169052565b506060830151620006bc606084018263ffffffff169052565b506080830151620006d6608084018264ffffffffff169052565b5060a0830151620006ed60a084018261ffff169052565b5060c08301516200070960c08401826001600160401b03169052565b5060e08301516200071f60e084018260ff169052565b50610100838101516001600160e01b0316908301526101208084015163ffffffff16908301526101408084015161ffff90811691840191909152610160808501518216908401526101809384015116929091019190915290565b608051615888620007bf600039600081816107af01528181610c9001528181610f240152818161103a01528181611b7101528181612a2f015261397801526158886000f3fe608060405234801561001057600080fd5b50600436106101a35760003560e01c80638da5cb5b116100ee578063d227d24511610097578063e4ddcea611610071578063e4ddcea6146105e8578063f2f22ef1146105fe578063f2fde38b14610606578063f6ea41f61461061957600080fd5b8063d227d2451461059d578063d328a91e146105cd578063e3d0e712146105d557600080fd5b8063b1dc65a4116100c8578063b1dc65a414610396578063c074ef21146103a9578063c3f909d4146103bc57600080fd5b80638da5cb5b1461032e578063a631571e14610356578063afcb95d71461037657600080fd5b80637d4807871161015057806381f1b9381161012a57806381f1b938146102a657806381ff7048146102ae57806385b214cf1461031b57600080fd5b80637d480787146102765780637f15e1661461027e578063814118341461029157600080fd5b806366316d8d1161018157806366316d8d1461023c5780637212762f1461024f57806379ba50971461026e57600080fd5b8063083a5466146101a8578063181f5a77146101bd578063626f458c1461020f575b600080fd5b6101bb6101b6366004614029565b610621565b005b6101f96040518060400160405280601c81526020017f46756e6374696f6e7320436f6f7264696e61746f722076312e332e310000000081525081565b60405161020691906140cf565b60405180910390f35b61022261021d36600461423c565b610676565b60405168ffffffffffffffffff9091168152602001610206565b6101bb61024a3660046142c3565b6106b3565b61025761086c565b6040805192835260ff909116602083015201610206565b6101bb610a7f565b6101bb610b7c565b6101bb61028c366004614029565b610d7c565b610299610dcc565b604051610206919061434d565b6101f9610e3b565b6102f860015460025463ffffffff74010000000000000000000000000000000000000000830481169378010000000000000000000000000000000000000000000000009093041691565b6040805163ffffffff948516815293909216602084015290820152606001610206565b6101bb610329366004614360565b610f0c565b60005460405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610206565b610369610364366004614379565b610fc9565b60405161020691906144ce565b604080516001815260006020820181905291810191909152606001610206565b6101bb6103a4366004614522565b61125a565b6101bb6103b736600461468d565b611871565b610590604080516101a081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e08101829052610100810182905261012081018290526101408101829052610160810182905261018081019190915250604080516101a08101825260085463ffffffff8082168352640100000000808304821660208501526801000000000000000083048216948401949094526c01000000000000000000000000820481166060840152700100000000000000000000000000000000820464ffffffffff1660808401527501000000000000000000000000000000000000000000820461ffff90811660a085015277010000000000000000000000000000000000000000000000830467ffffffffffffffff1660c08501527f010000000000000000000000000000000000000000000000000000000000000090920460ff1660e08401526009547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff81166101008501527c0100000000000000000000000000000000000000000000000000000000900416610120830152600a5480821661014084015262010000810482166101608401529290920490911661018082015290565b604051610206919061478b565b6105b06105ab3660046148b5565b611b6d565b6040516bffffffffffffffffffffffff9091168152602001610206565b6101f9611cdb565b6101bb6105e33660046149bd565b611d32565b6105f06128ae565b604051908152602001610206565b6102226129f2565b6101bb610614366004614a8a565b612a17565b610222612a2b565b610629612abc565b6000819003610664576040517f4f42be3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600f610671828483614b3a565b505050565b600a546000906106ad906064906106909061ffff16612b3f565b61069a9190614cb3565b6bffffffffffffffffffffffff16612b8c565b92915050565b6106bb612c2b565b806bffffffffffffffffffffffff166000036106f55750336000908152600b60205260409020546bffffffffffffffffffffffff1661074f565b336000908152600b60205260409020546bffffffffffffffffffffffff8083169116101561074f576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b336000908152600b60205260408120805483929061077c9084906bffffffffffffffffffffffff16614cde565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055506107d17f000000000000000000000000000000000000000000000000000000000000000090565b6040517f66316d8d00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff84811660048301526bffffffffffffffffffffffff8416602483015291909116906366316d8d90604401600060405180830381600087803b15801561085057600080fd5b505af1158015610864573d6000803e3d6000fd5b505050505050565b600080600080600d60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa1580156108df573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906109039190614d24565b5093505092505080426109169190614d74565b600854640100000000900463ffffffff161080156109435750600854640100000000900463ffffffff1615155b156109a057505060085477010000000000000000000000000000000000000000000000810467ffffffffffffffff16937f010000000000000000000000000000000000000000000000000000000000000090910460ff1692509050565b600082136109e2576040517f56b22ab8000000000000000000000000000000000000000000000000000000008152600481018390526024015b60405180910390fd5b600d54604080517f313ce5670000000000000000000000000000000000000000000000000000000081529051849273ffffffffffffffffffffffffffffffffffffffff169163313ce5679160048083019260209291908290030181865afa158015610a51573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a759190614d87565b9350935050509091565b60015473ffffffffffffffffffffffffffffffffffffffff163314610b00576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016109d9565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610b84612df2565b610b8c612c2b565b6000610b96610dcc565b905060005b8151811015610d78576000600b6000848481518110610bbc57610bbc614da4565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252810191909152604001600020546bffffffffffffffffffffffff1690508015610d67576000600b6000858581518110610c1b57610c1b614da4565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550610cb27f000000000000000000000000000000000000000000000000000000000000000090565b73ffffffffffffffffffffffffffffffffffffffff166366316d8d848481518110610cdf57610cdf614da4565b6020026020010151836040518363ffffffff1660e01b8152600401610d3492919073ffffffffffffffffffffffffffffffffffffffff9290921682526bffffffffffffffffffffffff16602082015260400190565b600060405180830381600087803b158015610d4e57600080fd5b505af1158015610d62573d6000803e3d6000fd5b505050505b50610d7181614dd3565b9050610b9b565b5050565b610d84612abc565b6000819003610dbf576040517f4f42be3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600e610671828483614b3a565b60606006805480602002602001604051908101604052809291908181526020018280548015610e3157602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610e06575b5050505050905090565b6060600f8054610e4a90614aa7565b9050600003610e85576040517f4f42be3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600f8054610e9290614aa7565b80601f0160208091040260200160405190810160405280929190818152602001828054610ebe90614aa7565b8015610e315780601f10610ee057610100808354040283529160200191610e31565b820191906000526020600020905b815481529060010190602001808311610eee57509395945050505050565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610f7b576040517fc41a5b0900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008181526007602052604080822091909155517f8a4b97add3359bd6bcf5e82874363670eb5ad0f7615abddbd0ed0a3a98f0f41690610fbe9083815260200190565b60405180910390a150565b6040805161016081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e08101829052610100810182905261012081018290526101408101919091523373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614611091576040517fc41a5b0900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006110a461109f84614e2d565b612dfa565b90925090506110b96060840160408501614a8a565b825173ffffffffffffffffffffffffffffffffffffffff91909116907fbf50768ccf13bd0110ca6d53a9c4f1f3271abdd4c24a56878863ed25b20598ff3261110760c0880160a08901614f1a565b61111961016089016101408a01614a8a565b6111238980614f37565b6111356101208c016101008d01614f9c565b60208c013561114b6101008e0160e08f01614fb7565b6040518061016001604052808e6000015181526020018e6020015173ffffffffffffffffffffffffffffffffffffffff1681526020018e604001516bffffffffffffffffffffffff1681526020018e6060015173ffffffffffffffffffffffffffffffffffffffff1681526020018e6080015167ffffffffffffffff1681526020018e60a0015163ffffffff1681526020018d68ffffffffffffffffff1681526020018e60e0015168ffffffffffffffffff1681526020018e610100015164ffffffffff1681526020018e610120015164ffffffffff1681526020018e610140015163ffffffff1681525060405161124b99989796959493929190614fd4565b60405180910390a3505b919050565b60008061126789896131ac565b915091508115611278575050611867565b604080518b3580825262ffffff6020808f0135600881901c9290921690840152909290917fb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62910160405180910390a16112d58b8b8b8b8b8b613335565b6003546000906002906112f39060ff8082169161010090041661508a565b6112fd91906150a3565b61130890600161508a565b60ff169050888114611376576040517f660bd4ba00000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f77726f6e67206e756d626572206f66207369676e61747572657300000000000060448201526064016109d9565b888714611405576040517f660bd4ba00000000000000000000000000000000000000000000000000000000815260206004820152602860248201527f7265706f727420727320616e64207373206d757374206265206f66206571756160448201527f6c206c656e67746800000000000000000000000000000000000000000000000060648201526084016109d9565b3360009081526004602090815260408083208151808301909252805460ff80821684529293919291840191610100909104166002811115611448576114486150c5565b6002811115611459576114596150c5565b9052509050600281602001516002811115611476576114766150c5565b141580156114bf57506006816000015160ff168154811061149957611499614da4565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff163314155b15611526576040517f660bd4ba00000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f756e617574686f72697a6564207472616e736d6974746572000000000000000060448201526064016109d9565b50505050611532613fc8565b60008a8a6040516115449291906150f4565b60405190819003812061155b918e90602001615104565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181528282528051602091820120838301909252600080845290830152915060005b898110156118575760006001848984602081106115c4576115c4614da4565b6115d191901a601b61508a565b8e8e868181106115e3576115e3614da4565b905060200201358d8d878181106115fc576115fc614da4565b9050602002013560405160008152602001604052604051611639949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa15801561165b573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff811660009081526004602090815290849020838501909452835460ff808216855292965092945084019161010090041660028111156116db576116db6150c5565b60028111156116ec576116ec6150c5565b9052509250600183602001516002811115611709576117096150c5565b14611770576040517f660bd4ba00000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f61646472657373206e6f7420617574686f72697a656420746f207369676e000060448201526064016109d9565b8251600090869060ff16601f811061178a5761178a614da4565b602002015173ffffffffffffffffffffffffffffffffffffffff161461180c576040517f660bd4ba00000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f6e6f6e2d756e69717565207369676e617475726500000000000000000000000060448201526064016109d9565b8085846000015160ff16601f811061182657611826614da4565b73ffffffffffffffffffffffffffffffffffffffff90921660209290920201525061185081614dd3565b90506115a5565b505050611863826133ec565b5050505b5050505050505050565b611879612df2565b80516008805460208401516040808601516060870151608088015160a089015160c08a015160e08b015163ffffffff9a8b167fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000090991698909817640100000000978b168802177fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff1668010000000000000000958b16959095027fffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffffff16949094176c01000000000000000000000000938a1693909302929092177fffffffffffffffffff00000000000000ffffffffffffffffffffffffffffffff1670010000000000000000000000000000000064ffffffffff909216919091027fffffffffffffffffff0000ffffffffffffffffffffffffffffffffffffffffff1617750100000000000000000000000000000000000000000061ffff928316021776ffffffffffffffffffffffffffffffffffffffffffffff167701000000000000000000000000000000000000000000000067ffffffffffffffff909316929092027effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff16919091177f010000000000000000000000000000000000000000000000000000000000000060ff90951694909402939093179093556101008501516101208601517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff9091167c0100000000000000000000000000000000000000000000000000000000919095160293909317600955610140840151600a80546101608701516101808801519385167fffffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000909216919091176201000091851691909102177fffffffffffffffffffffffffffffffffffffffffffffffffffff0000ffffffff169190921690930292909217909155517f9832846a11905a128363408ace5f9803f449fcb017d3cf7e22a50c6e332931b290610fbe90839061478b565b60007f00000000000000000000000000000000000000000000000000000000000000006040517f10fc49c100000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8816600482015263ffffffff8516602482015273ffffffffffffffffffffffffffffffffffffffff91909116906310fc49c19060440160006040518083038186803b158015611c0d57600080fd5b505afa158015611c21573d6000803e3d6000fd5b5050505066038d7ea4c68000821115611c66576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000611c70612a2b565b90506000611cb387878080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061067692505050565b90506000611cbf6129f2565b9050611cce868684868561353b565b9998505050505050505050565b6060600e8054611cea90614aa7565b9050600003611d25576040517f4f42be3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600e8054610e9290614aa7565b855185518560ff16601f831115611da5576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f746f6f206d616e79207369676e6572730000000000000000000000000000000060448201526064016109d9565b80600003611e0f576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f66206d75737420626520706f736974697665000000000000000000000000000060448201526064016109d9565b818314611e9d576040517f89a61989000000000000000000000000000000000000000000000000000000008152602060048201526024808201527f6f7261636c6520616464726573736573206f7574206f6620726567697374726160448201527f74696f6e0000000000000000000000000000000000000000000000000000000060648201526084016109d9565b611ea8816003615118565b8311611f10576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f6661756c74792d6f7261636c65206620746f6f2068696768000000000000000060448201526064016109d9565b611f18612abc565b6040805160c0810182528a8152602081018a905260ff89169181018290526060810188905267ffffffffffffffff8716608082015260a0810186905290611f5f9088613699565b6005541561211457600554600090611f7990600190614d74565b9050600060058281548110611f9057611f90614da4565b60009182526020822001546006805473ffffffffffffffffffffffffffffffffffffffff90921693509084908110611fca57611fca614da4565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff85811684526004909252604080842080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00009081169091559290911680845292208054909116905560058054919250908061204a5761204a61512f565b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff000000000000000000000000000000000000000016905501905560068054806120b3576120b361512f565b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff000000000000000000000000000000000000000016905501905550611f5f915050565b60005b8151518110156126cb5781518051600091908390811061213957612139614da4565b602002602001015173ffffffffffffffffffffffffffffffffffffffff16036121be576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f7369676e6572206d757374206e6f7420626520656d707479000000000000000060448201526064016109d9565b600073ffffffffffffffffffffffffffffffffffffffff16826020015182815181106121ec576121ec614da4565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1603612271576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f7472616e736d6974746572206d757374206e6f7420626520656d70747900000060448201526064016109d9565b6000600460008460000151848151811061228d5761228d614da4565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff16825281019190915260400160002054610100900460ff1660028111156122d7576122d76150c5565b1461233e576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265706561746564207369676e6572206164647265737300000000000000000060448201526064016109d9565b6040805180820190915260ff8216815260016020820152825180516004916000918590811061236f5761236f614da4565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040016000208251815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082168117835592840151919283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001617610100836002811115612410576124106150c5565b0217905550600091506124209050565b600460008460200151848151811061243a5761243a614da4565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff16825281019190915260400160002054610100900460ff166002811115612484576124846150c5565b146124eb576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f7265706561746564207472616e736d697474657220616464726573730000000060448201526064016109d9565b6040805180820190915260ff82168152602081016002815250600460008460200151848151811061251e5761251e614da4565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040016000208251815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082168117835592840151919283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000016176101008360028111156125bf576125bf6150c5565b0217905550508251805160059250839081106125dd576125dd614da4565b602090810291909101810151825460018101845560009384529282902090920180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909316929092179091558201518051600691908390811061265957612659614da4565b60209081029190910181015182546001810184556000938452919092200180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909216919091179055806126c381614dd3565b915050612117565b506040810151600380547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff909216919091179055600180547fffffffff00000000ffffffffffffffffffffffffffffffffffffffffffffffff8116780100000000000000000000000000000000000000000000000063ffffffff43811682029290921780855592048116929182916014916127839184917401000000000000000000000000000000000000000090041661515e565b92506101000a81548163ffffffff021916908363ffffffff1602179055506127e24630600160149054906101000a900463ffffffff1663ffffffff16856000015186602001518760400151886060015189608001518a60a001516136b2565b600281905582518051600380547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1661010060ff9093169290920291909117905560015460208501516040808701516060880151608089015160a08a015193517f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e0598612899988b9891977401000000000000000000000000000000000000000090920463ffffffff1696909591949193919261517b565b60405180910390a15050505050505050505050565b6000806000600c8054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa15801561291e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906129429190614d24565b5093505092505080426129559190614d74565b600854640100000000900463ffffffff161080156129825750600854640100000000900463ffffffff1615155b156129af5750506009547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16919050565b600082136129ec576040517f43d4cf66000000000000000000000000000000000000000000000000000000008152600481018390526024016109d9565b50919050565b600a54600090612a12906064906106909062010000900461ffff16612b3f565b905090565b612a1f612abc565b612a288161375d565b50565b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16632a905ccc6040518163ffffffff1660e01b8152600401602060405180830381865afa158015612a98573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612a12919061520c565b60005473ffffffffffffffffffffffffffffffffffffffff163314612b3d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016109d9565b565b6000806000612b4c61086c565b9092509050612b8482612b6083601261508a565b612b6b90600a615349565b612b759087615118565b612b7f9190615358565b613852565b949350505050565b600068ffffffffffffffffff821115612c27576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f53616665436173743a2076616c756520646f65736e27742066697420696e203760448201527f322062697473000000000000000000000000000000000000000000000000000060648201526084016109d9565b5090565b600c546bffffffffffffffffffffffff16600003612c4557565b6000612c4f610dcc565b80519091506000819003612c8f576040517f30274b3a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600c54600090612cae9083906bffffffffffffffffffffffff16614cb3565b9050806bffffffffffffffffffffffff16600003612ccb57505050565b60005b82811015612d945781600b6000868481518110612ced57612ced614da4565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282829054906101000a90046bffffffffffffffffffffffff16612d55919061536c565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555080612d8d90614dd3565b9050612cce565b50612d9f8282615391565b600c8054600090612dbf9084906bffffffffffffffffffffffff16614cde565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550505050565b612b3d612abc565b6040805161016081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e081018290526101008101829052610120810182905261014081019190915260085461010083015160009161ffff7501000000000000000000000000000000000000000000909104811691161115612eb8576040517fdada758700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000612ec78460000151610676565b9050612ed16129f2565b91506000612eea8560e001513a8488608001518761353b565b9050806bffffffffffffffffffffffff1685606001516bffffffffffffffffffffffff161015612f46576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600954600090612f7c907c0100000000000000000000000000000000000000000000000000000000900463ffffffff16426153b9565b905060003087604001518860a001518960c001516001612f9c91906153cc565b8a5180516020918201206101008d015160e08e015160405161305098979695948c918c9132910173ffffffffffffffffffffffffffffffffffffffff9a8b168152988a1660208a015267ffffffffffffffff97881660408a0152959096166060880152608087019390935261ffff9190911660a086015263ffffffff90811660c08601526bffffffffffffffffffffffff9190911660e0850152919091166101008301529091166101208201526101400190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815282825280516020918201206101608401835280845230848301526bffffffffffffffffffffffff8716848401528a83015173ffffffffffffffffffffffffffffffffffffffff16606085015260a0808c015167ffffffffffffffff1660808087019190915260e0808e015163ffffffff90811693880193909352908d015168ffffffffffffffffff90811660c08801528a169086015260085468010000000000000000810482166101008701526c01000000000000000000000000900481166101208601528616610140850152915192985090925061315c918891016144ce565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529181528151602092830120600093845260079092529091205550929491935090915050565b60006131e06040518060a0016040528060608152602001606081526020016060815260200160608152602001606081525090565b6000808080806131f2888a018a6154c8565b84519499509297509095509350915060ff16801580613212575084518114155b8061321e575083518114155b8061322a575082518114155b80613236575081518114155b1561329d576040517f660bd4ba00000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f4669656c6473206d75737420626520657175616c206c656e677468000000000060448201526064016109d9565b60005b81811015613303576132d98782815181106132bd576132bd614da4565b6020026020010151600090815260076020526040902054151590565b613303576132e8600183614d74565b81036132f357600198505b6132fc81614dd3565b90506132a0565b50506040805160a0810182529586526020860194909452928401919091526060830152608082015290505b9250929050565b6000613342826020615118565b61334d856020615118565b613359886101446153b9565b61336391906153b9565b61336d91906153b9565b6133789060006153b9565b90503681146133e3576040517f660bd4ba00000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f63616c6c64617461206c656e677468206d69736d61746368000000000000000060448201526064016109d9565b50505050505050565b80515160ff1660005b8181101561067157600061349e8460000151838151811061341857613418614da4565b60200260200101518560200151848151811061343657613436614da4565b60200260200101518660400151858151811061345457613454614da4565b60200260200101518760600151868151811061347257613472614da4565b60200260200101518860800151878151811061349057613490614da4565b6020026020010151886138f0565b905060008160068111156134b4576134b46150c5565b14806134d1575060018160068111156134cf576134cf6150c5565b145b1561352a5783518051839081106134ea576134ea614da4565b60209081029190910181015160405133815290917fc708e0440951fd63499c0f7a73819b469ee5dd3ecc356c0ab4eb7f18389009d9910160405180910390a25b5061353481614dd3565b90506133f5565b600854600090700100000000000000000000000000000000900464ffffffffff1685101561358457600854700100000000000000000000000000000000900464ffffffffff1694505b60085460009087906135ba9063ffffffff6c0100000000000000000000000082048116916801000000000000000090041661515e565b6135c4919061515e565b600a5463ffffffff9190911691506000906135ea90640100000000900461ffff16613d82565b90506000816135f9848a615118565b61360391906153b9565b600854909150600090612710906136209063ffffffff1684615118565b61362a9190615358565b61363490836153b9565b9050600061364182613ece565b905060008768ffffffffffffffffff168968ffffffffffffffffff168b68ffffffffffffffffff16613673919061536c565b61367d919061536c565b9050613689818361536c565b9c9b505050505050505050505050565b60006136a3610dcc565b511115610d7857610d78612c2b565b6000808a8a8a8a8a8a8a8a8a6040516020016136d69998979695949392919061559a565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179150509998505050505050505050565b3373ffffffffffffffffffffffffffffffffffffffff8216036137dc576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016109d9565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60006bffffffffffffffffffffffff821115612c27576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f53616665436173743a2076616c756520646f65736e27742066697420696e203960448201527f362062697473000000000000000000000000000000000000000000000000000060648201526084016109d9565b600080848060200190518101906139079190615666565b905060003a826101200151836101000151613922919061572e565b64ffffffffff166139339190615118565b9050600060ff851661394436613d82565b61394e9190615358565b9050600061396461395f83856153b9565b613ece565b905060006139713a613ece565b90506000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663330605298e8e868b60c0015168ffffffffffffffffff168c60e0015168ffffffffffffffffff168a6139e0919061536c565b6139ea919061536c565b336040518061016001604052808f6000015181526020018f6020015173ffffffffffffffffffffffffffffffffffffffff1681526020018f604001516bffffffffffffffffffffffff1681526020018f6060015173ffffffffffffffffffffffffffffffffffffffff1681526020018f6080015167ffffffffffffffff1681526020018f60a0015163ffffffff168152602001600068ffffffffffffffffff1681526020018f60e0015168ffffffffffffffffff1681526020018f610100015164ffffffffff1681526020018f610120015164ffffffffff1681526020018f610140015163ffffffff168152506040518763ffffffff1660e01b8152600401613af89695949392919061574c565b60408051808303816000875af1158015613b16573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613b3a91906157c8565b90925090506000826006811115613b5357613b536150c5565b1480613b7057506001826006811115613b6e57613b6e6150c5565b145b15613d715760008e815260076020526040812055613b8e818561536c565b336000908152600b602052604081208054909190613bbb9084906bffffffffffffffffffffffff1661536c565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055508660e0015168ffffffffffffffffff16600c60008282829054906101000a90046bffffffffffffffffffffffff16613c21919061536c565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055508660c0015168ffffffffffffffffff16600b6000613c6b613eed565b73ffffffffffffffffffffffffffffffffffffffff168152602081019190915260400160009081208054909190613cb19084906bffffffffffffffffffffffff1661536c565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055508d7f08a4a0761e3c98d288cb4af9342660f49550d83139fb3b762b70d34bed6273688487848b60e0015160008d60c00151604051613d68969594939291906bffffffffffffffffffffffff9687168152602081019590955292909416604084015268ffffffffffffffffff9081166060840152928316608083015290911660a082015260c00190565b60405180910390a25b509c9b505050505050505050505050565b600046613d8e81613f5e565b15613e23576000606c73ffffffffffffffffffffffffffffffffffffffff166341b247a86040518163ffffffff1660e01b815260040160c060405180830381865afa158015613de1573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613e0591906157fb565b50505050915050608c84613e1991906153b9565b612b849082615118565b613e2c81613f81565b15613ec5576040517ff1c7a58b0000000000000000000000000000000000000000000000000000000081526004810184905273420000000000000000000000000000000000000f9063f1c7a58b90602401602060405180830381865afa158015613e9a573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613ebe9190615845565b9392505050565b50600092915050565b60006106ad613edb6128ae565b612b7584670de0b6b3a7640000615118565b60003073ffffffffffffffffffffffffffffffffffffffff16638da5cb5b6040518163ffffffff1660e01b8152600401602060405180830381865afa158015613f3a573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612a12919061585e565b600061a4b1821480613f72575062066eed82145b806106ad57505062066eee1490565b6000600a821480613f9357506101a482145b80613fa0575062aa37dc82145b80613fac575061210582145b80613fb9575062014a3382145b806106ad57505062014a341490565b604051806103e00160405280601f906020820280368337509192915050565b60008083601f840112613ff957600080fd5b50813567ffffffffffffffff81111561401157600080fd5b60208301915083602082850101111561332e57600080fd5b6000806020838503121561403c57600080fd5b823567ffffffffffffffff81111561405357600080fd5b61405f85828601613fe7565b90969095509350505050565b6000815180845260005b8181101561409157602081850181015186830182015201614075565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b602081526000613ebe602083018461406b565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040516101a0810167ffffffffffffffff81118282101715614135576141356140e2565b60405290565b604051610160810167ffffffffffffffff81118282101715614135576141356140e2565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156141a6576141a66140e2565b604052919050565b600082601f8301126141bf57600080fd5b813567ffffffffffffffff8111156141d9576141d96140e2565b61420a60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8401160161415f565b81815284602083860101111561421f57600080fd5b816020850160208301376000918101602001919091529392505050565b60006020828403121561424e57600080fd5b813567ffffffffffffffff81111561426557600080fd5b612b84848285016141ae565b73ffffffffffffffffffffffffffffffffffffffff81168114612a2857600080fd5b803561125581614271565b6bffffffffffffffffffffffff81168114612a2857600080fd5b80356112558161429e565b600080604083850312156142d657600080fd5b82356142e181614271565b915060208301356142f18161429e565b809150509250929050565b600081518084526020808501945080840160005b8381101561434257815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101614310565b509495945050505050565b602081526000613ebe60208301846142fc565b60006020828403121561437257600080fd5b5035919050565b60006020828403121561438b57600080fd5b813567ffffffffffffffff8111156143a257600080fd5b82016101608185031215613ebe57600080fd5b8051825260208101516143e0602084018273ffffffffffffffffffffffffffffffffffffffff169052565b50604081015161440060408401826bffffffffffffffffffffffff169052565b506060810151614428606084018273ffffffffffffffffffffffffffffffffffffffff169052565b506080810151614444608084018267ffffffffffffffff169052565b5060a081015161445c60a084018263ffffffff169052565b5060c081015161447960c084018268ffffffffffffffffff169052565b5060e081015161449660e084018268ffffffffffffffffff169052565b506101008181015164ffffffffff9081169184019190915261012080830151909116908301526101409081015163ffffffff16910152565b61016081016106ad82846143b5565b60008083601f8401126144ef57600080fd5b50813567ffffffffffffffff81111561450757600080fd5b6020830191508360208260051b850101111561332e57600080fd5b60008060008060008060008060e0898b03121561453e57600080fd5b606089018a81111561454f57600080fd5b8998503567ffffffffffffffff8082111561456957600080fd5b6145758c838d01613fe7565b909950975060808b013591508082111561458e57600080fd5b61459a8c838d016144dd565b909750955060a08b01359150808211156145b357600080fd5b506145c08b828c016144dd565b999c989b50969995989497949560c00135949350505050565b63ffffffff81168114612a2857600080fd5b8035611255816145d9565b64ffffffffff81168114612a2857600080fd5b8035611255816145f6565b803561ffff8116811461125557600080fd5b67ffffffffffffffff81168114612a2857600080fd5b803561125581614626565b60ff81168114612a2857600080fd5b803561125581614647565b80357bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8116811461125557600080fd5b60006101a082840312156146a057600080fd5b6146a8614111565b6146b1836145eb565b81526146bf602084016145eb565b60208201526146d0604084016145eb565b60408201526146e1606084016145eb565b60608201526146f260808401614609565b608082015261470360a08401614614565b60a082015261471460c0840161463c565b60c082015261472560e08401614656565b60e0820152610100614738818501614661565b9082015261012061474a8482016145eb565b9082015261014061475c848201614614565b9082015261016061476e848201614614565b90820152610180614780848201614614565b908201529392505050565b815163ffffffff1681526101a0810160208301516147b1602084018263ffffffff169052565b5060408301516147c9604084018263ffffffff169052565b5060608301516147e1606084018263ffffffff169052565b5060808301516147fa608084018264ffffffffff169052565b5060a083015161481060a084018261ffff169052565b5060c083015161482c60c084018267ffffffffffffffff169052565b5060e083015161484160e084018260ff169052565b50610100838101517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16908301526101208084015163ffffffff16908301526101408084015161ffff90811691840191909152610160808501518216908401526101808085015191821681850152905b505092915050565b6000806000806000608086880312156148cd57600080fd5b85356148d881614626565b9450602086013567ffffffffffffffff8111156148f457600080fd5b61490088828901613fe7565b9095509350506040860135614914816145d9565b949793965091946060013592915050565b600067ffffffffffffffff82111561493f5761493f6140e2565b5060051b60200190565b600082601f83011261495a57600080fd5b8135602061496f61496a83614925565b61415f565b82815260059290921b8401810191818101908684111561498e57600080fd5b8286015b848110156149b25780356149a581614271565b8352918301918301614992565b509695505050505050565b60008060008060008060c087890312156149d657600080fd5b863567ffffffffffffffff808211156149ee57600080fd5b6149fa8a838b01614949565b97506020890135915080821115614a1057600080fd5b614a1c8a838b01614949565b9650614a2a60408a01614656565b95506060890135915080821115614a4057600080fd5b614a4c8a838b016141ae565b9450614a5a60808a0161463c565b935060a0890135915080821115614a7057600080fd5b50614a7d89828a016141ae565b9150509295509295509295565b600060208284031215614a9c57600080fd5b8135613ebe81614271565b600181811c90821680614abb57607f821691505b6020821081036129ec577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b601f82111561067157600081815260208120601f850160051c81016020861015614b1b5750805b601f850160051c820191505b8181101561086457828155600101614b27565b67ffffffffffffffff831115614b5257614b526140e2565b614b6683614b608354614aa7565b83614af4565b6000601f841160018114614bb85760008515614b825750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b178355614c4e565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b82811015614c075786850135825560209485019460019092019101614be7565b5086821015614c42577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006bffffffffffffffffffffffff80841680614cd257614cd2614c55565b92169190910492915050565b6bffffffffffffffffffffffff828116828216039080821115614d0357614d03614c84565b5092915050565b805169ffffffffffffffffffff8116811461125557600080fd5b600080600080600060a08688031215614d3c57600080fd5b614d4586614d0a565b9450602086015193506040860151925060608601519150614d6860808701614d0a565b90509295509295909350565b818103818111156106ad576106ad614c84565b600060208284031215614d9957600080fd5b8151613ebe81614647565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203614e0457614e04614c84565b5060010190565b68ffffffffffffffffff81168114612a2857600080fd5b803561125581614e0b565b60006101608236031215614e4057600080fd5b614e4861413b565b823567ffffffffffffffff811115614e5f57600080fd5b614e6b368286016141ae565b82525060208301356020820152614e8460408401614293565b6040820152614e95606084016142b8565b6060820152614ea660808401614e22565b6080820152614eb760a0840161463c565b60a0820152614ec860c0840161463c565b60c0820152614ed960e084016145eb565b60e0820152610100614eec818501614614565b90820152610120614efe84820161463c565b90820152610140614f10848201614293565b9082015292915050565b600060208284031215614f2c57600080fd5b8135613ebe81614626565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112614f6c57600080fd5b83018035915067ffffffffffffffff821115614f8757600080fd5b60200191503681900382131561332e57600080fd5b600060208284031215614fae57600080fd5b613ebe82614614565b600060208284031215614fc957600080fd5b8135613ebe816145d9565b73ffffffffffffffffffffffffffffffffffffffff8a8116825267ffffffffffffffff8a166020830152881660408201526102406060820181905281018690526000610260878982850137600083890182015261ffff8716608084015260a0830186905263ffffffff851660c0840152601f88017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016830101905061507c60e08301846143b5565b9a9950505050505050505050565b60ff81811683821601908111156106ad576106ad614c84565b600060ff8316806150b6576150b6614c55565b8060ff84160491505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b8183823760009101908152919050565b828152606082602083013760800192915050565b80820281158282048414176106ad576106ad614c84565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b63ffffffff818116838216019080821115614d0357614d03614c84565b600061012063ffffffff808d1684528b6020850152808b166040850152508060608401526151ab8184018a6142fc565b905082810360808401526151bf81896142fc565b905060ff871660a084015282810360c08401526151dc818761406b565b905067ffffffffffffffff851660e0840152828103610100840152613689818561406b565b805161125581614e0b565b60006020828403121561521e57600080fd5b8151613ebe81614e0b565b600181815b8085111561528257817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0482111561526857615268614c84565b8085161561527557918102915b93841c939080029061522e565b509250929050565b600082615299575060016106ad565b816152a6575060006106ad565b81600181146152bc57600281146152c6576152e2565b60019150506106ad565b60ff8411156152d7576152d7614c84565b50506001821b6106ad565b5060208310610133831016604e8410600b8410161715615305575081810a6106ad565b61530f8383615229565b807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0482111561534157615341614c84565b029392505050565b6000613ebe60ff84168361528a565b60008261536757615367614c55565b500490565b6bffffffffffffffffffffffff818116838216019080821115614d0357614d03614c84565b6bffffffffffffffffffffffff8181168382160280821691908281146148ad576148ad614c84565b808201808211156106ad576106ad614c84565b67ffffffffffffffff818116838216019080821115614d0357614d03614c84565b600082601f8301126153fe57600080fd5b8135602061540e61496a83614925565b82815260059290921b8401810191818101908684111561542d57600080fd5b8286015b848110156149b25780358352918301918301615431565b600082601f83011261545957600080fd5b8135602061546961496a83614925565b82815260059290921b8401810191818101908684111561548857600080fd5b8286015b848110156149b257803567ffffffffffffffff8111156154ac5760008081fd5b6154ba8986838b01016141ae565b84525091830191830161548c565b600080600080600060a086880312156154e057600080fd5b853567ffffffffffffffff808211156154f857600080fd5b61550489838a016153ed565b9650602088013591508082111561551a57600080fd5b61552689838a01615448565b9550604088013591508082111561553c57600080fd5b61554889838a01615448565b9450606088013591508082111561555e57600080fd5b61556a89838a01615448565b9350608088013591508082111561558057600080fd5b5061558d88828901615448565b9150509295509295909350565b60006101208b835273ffffffffffffffffffffffffffffffffffffffff8b16602084015267ffffffffffffffff808b1660408501528160608501526155e18285018b6142fc565b915083820360808501526155f5828a6142fc565b915060ff881660a085015283820360c0850152615612828861406b565b90861660e08501528381036101008501529050613689818561406b565b805161125581614271565b80516112558161429e565b805161125581614626565b8051611255816145d9565b8051611255816145f6565b6000610160828403121561567957600080fd5b61568161413b565b825181526156916020840161562f565b60208201526156a26040840161563a565b60408201526156b36060840161562f565b60608201526156c460808401615645565b60808201526156d560a08401615650565b60a08201526156e660c08401615201565b60c08201526156f760e08401615201565b60e082015261010061570a81850161565b565b9082015261012061571c84820161565b565b90820152610140614780848201615650565b64ffffffffff818116838216019080821115614d0357614d03614c84565b60006102008083526157608184018a61406b565b90508281036020840152615774818961406b565b6bffffffffffffffffffffffff88811660408601528716606085015273ffffffffffffffffffffffffffffffffffffffff8616608085015291506157bd905060a08301846143b5565b979650505050505050565b600080604083850312156157db57600080fd5b8251600781106157ea57600080fd5b60208401519092506142f18161429e565b60008060008060008060c0878903121561581457600080fd5b865195506020870151945060408701519350606087015192506080870151915060a087015190509295509295509295565b60006020828403121561585757600080fd5b5051919050565b60006020828403121561587057600080fd5b8151613ebe8161427156fea164736f6c6343000813000a",
}

var FunctionsCoordinatorABI = FunctionsCoordinatorMetaData.ABI

var FunctionsCoordinatorBin = FunctionsCoordinatorMetaData.Bin

func DeployFunctionsCoordinator(auth *bind.TransactOpts, backend bind.ContractBackend, router common.Address, config FunctionsBillingConfig, linkToNativeFeed common.Address, linkToUsdFeed common.Address) (common.Address, *types.Transaction, *FunctionsCoordinator, error) {
	parsed, err := FunctionsCoordinatorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(FunctionsCoordinatorBin), backend, router, config, linkToNativeFeed, linkToUsdFeed)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &FunctionsCoordinator{address: address, abi: *parsed, FunctionsCoordinatorCaller: FunctionsCoordinatorCaller{contract: contract}, FunctionsCoordinatorTransactor: FunctionsCoordinatorTransactor{contract: contract}, FunctionsCoordinatorFilterer: FunctionsCoordinatorFilterer{contract: contract}}, nil
}

type FunctionsCoordinator struct {
	address common.Address
	abi     abi.ABI
	FunctionsCoordinatorCaller
	FunctionsCoordinatorTransactor
	FunctionsCoordinatorFilterer
}

type FunctionsCoordinatorCaller struct {
	contract *bind.BoundContract
}

type FunctionsCoordinatorTransactor struct {
	contract *bind.BoundContract
}

type FunctionsCoordinatorFilterer struct {
	contract *bind.BoundContract
}

type FunctionsCoordinatorSession struct {
	Contract     *FunctionsCoordinator
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type FunctionsCoordinatorCallerSession struct {
	Contract *FunctionsCoordinatorCaller
	CallOpts bind.CallOpts
}

type FunctionsCoordinatorTransactorSession struct {
	Contract     *FunctionsCoordinatorTransactor
	TransactOpts bind.TransactOpts
}

type FunctionsCoordinatorRaw struct {
	Contract *FunctionsCoordinator
}

type FunctionsCoordinatorCallerRaw struct {
	Contract *FunctionsCoordinatorCaller
}

type FunctionsCoordinatorTransactorRaw struct {
	Contract *FunctionsCoordinatorTransactor
}

func NewFunctionsCoordinator(address common.Address, backend bind.ContractBackend) (*FunctionsCoordinator, error) {
	abi, err := abi.JSON(strings.NewReader(FunctionsCoordinatorABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindFunctionsCoordinator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FunctionsCoordinator{address: address, abi: abi, FunctionsCoordinatorCaller: FunctionsCoordinatorCaller{contract: contract}, FunctionsCoordinatorTransactor: FunctionsCoordinatorTransactor{contract: contract}, FunctionsCoordinatorFilterer: FunctionsCoordinatorFilterer{contract: contract}}, nil
}

func NewFunctionsCoordinatorCaller(address common.Address, caller bind.ContractCaller) (*FunctionsCoordinatorCaller, error) {
	contract, err := bindFunctionsCoordinator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FunctionsCoordinatorCaller{contract: contract}, nil
}

func NewFunctionsCoordinatorTransactor(address common.Address, transactor bind.ContractTransactor) (*FunctionsCoordinatorTransactor, error) {
	contract, err := bindFunctionsCoordinator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FunctionsCoordinatorTransactor{contract: contract}, nil
}

func NewFunctionsCoordinatorFilterer(address common.Address, filterer bind.ContractFilterer) (*FunctionsCoordinatorFilterer, error) {
	contract, err := bindFunctionsCoordinator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FunctionsCoordinatorFilterer{contract: contract}, nil
}

func bindFunctionsCoordinator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FunctionsCoordinatorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_FunctionsCoordinator *FunctionsCoordinatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FunctionsCoordinator.Contract.FunctionsCoordinatorCaller.contract.Call(opts, result, method, params...)
}

func (_FunctionsCoordinator *FunctionsCoordinatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsCoordinator.Contract.FunctionsCoordinatorTransactor.contract.Transfer(opts)
}

func (_FunctionsCoordinator *FunctionsCoordinatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FunctionsCoordinator.Contract.FunctionsCoordinatorTransactor.contract.Transact(opts, method, params...)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FunctionsCoordinator.Contract.contract.Call(opts, result, method, params...)
}

func (_FunctionsCoordinator *FunctionsCoordinatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsCoordinator.Contract.contract.Transfer(opts)
}

func (_FunctionsCoordinator *FunctionsCoordinatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FunctionsCoordinator.Contract.contract.Transact(opts, method, params...)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCaller) EstimateCost(opts *bind.CallOpts, subscriptionId uint64, data []byte, callbackGasLimit uint32, gasPriceWei *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _FunctionsCoordinator.contract.Call(opts, &out, "estimateCost", subscriptionId, data, callbackGasLimit, gasPriceWei)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_FunctionsCoordinator *FunctionsCoordinatorSession) EstimateCost(subscriptionId uint64, data []byte, callbackGasLimit uint32, gasPriceWei *big.Int) (*big.Int, error) {
	return _FunctionsCoordinator.Contract.EstimateCost(&_FunctionsCoordinator.CallOpts, subscriptionId, data, callbackGasLimit, gasPriceWei)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCallerSession) EstimateCost(subscriptionId uint64, data []byte, callbackGasLimit uint32, gasPriceWei *big.Int) (*big.Int, error) {
	return _FunctionsCoordinator.Contract.EstimateCost(&_FunctionsCoordinator.CallOpts, subscriptionId, data, callbackGasLimit, gasPriceWei)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCaller) GetAdminFeeJuels(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FunctionsCoordinator.contract.Call(opts, &out, "getAdminFeeJuels")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_FunctionsCoordinator *FunctionsCoordinatorSession) GetAdminFeeJuels() (*big.Int, error) {
	return _FunctionsCoordinator.Contract.GetAdminFeeJuels(&_FunctionsCoordinator.CallOpts)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCallerSession) GetAdminFeeJuels() (*big.Int, error) {
	return _FunctionsCoordinator.Contract.GetAdminFeeJuels(&_FunctionsCoordinator.CallOpts)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCaller) GetConfig(opts *bind.CallOpts) (FunctionsBillingConfig, error) {
	var out []interface{}
	err := _FunctionsCoordinator.contract.Call(opts, &out, "getConfig")

	if err != nil {
		return *new(FunctionsBillingConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(FunctionsBillingConfig)).(*FunctionsBillingConfig)

	return out0, err

}

func (_FunctionsCoordinator *FunctionsCoordinatorSession) GetConfig() (FunctionsBillingConfig, error) {
	return _FunctionsCoordinator.Contract.GetConfig(&_FunctionsCoordinator.CallOpts)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCallerSession) GetConfig() (FunctionsBillingConfig, error) {
	return _FunctionsCoordinator.Contract.GetConfig(&_FunctionsCoordinator.CallOpts)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCaller) GetDONFeeJuels(opts *bind.CallOpts, arg0 []byte) (*big.Int, error) {
	var out []interface{}
	err := _FunctionsCoordinator.contract.Call(opts, &out, "getDONFeeJuels", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_FunctionsCoordinator *FunctionsCoordinatorSession) GetDONFeeJuels(arg0 []byte) (*big.Int, error) {
	return _FunctionsCoordinator.Contract.GetDONFeeJuels(&_FunctionsCoordinator.CallOpts, arg0)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCallerSession) GetDONFeeJuels(arg0 []byte) (*big.Int, error) {
	return _FunctionsCoordinator.Contract.GetDONFeeJuels(&_FunctionsCoordinator.CallOpts, arg0)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCaller) GetDONPublicKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _FunctionsCoordinator.contract.Call(opts, &out, "getDONPublicKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_FunctionsCoordinator *FunctionsCoordinatorSession) GetDONPublicKey() ([]byte, error) {
	return _FunctionsCoordinator.Contract.GetDONPublicKey(&_FunctionsCoordinator.CallOpts)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCallerSession) GetDONPublicKey() ([]byte, error) {
	return _FunctionsCoordinator.Contract.GetDONPublicKey(&_FunctionsCoordinator.CallOpts)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCaller) GetOperationFeeJuels(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FunctionsCoordinator.contract.Call(opts, &out, "getOperationFeeJuels")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_FunctionsCoordinator *FunctionsCoordinatorSession) GetOperationFeeJuels() (*big.Int, error) {
	return _FunctionsCoordinator.Contract.GetOperationFeeJuels(&_FunctionsCoordinator.CallOpts)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCallerSession) GetOperationFeeJuels() (*big.Int, error) {
	return _FunctionsCoordinator.Contract.GetOperationFeeJuels(&_FunctionsCoordinator.CallOpts)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCaller) GetThresholdPublicKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _FunctionsCoordinator.contract.Call(opts, &out, "getThresholdPublicKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_FunctionsCoordinator *FunctionsCoordinatorSession) GetThresholdPublicKey() ([]byte, error) {
	return _FunctionsCoordinator.Contract.GetThresholdPublicKey(&_FunctionsCoordinator.CallOpts)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCallerSession) GetThresholdPublicKey() ([]byte, error) {
	return _FunctionsCoordinator.Contract.GetThresholdPublicKey(&_FunctionsCoordinator.CallOpts)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCaller) GetUsdPerUnitLink(opts *bind.CallOpts) (*big.Int, uint8, error) {
	var out []interface{}
	err := _FunctionsCoordinator.contract.Call(opts, &out, "getUsdPerUnitLink")

	if err != nil {
		return *new(*big.Int), *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(uint8)).(*uint8)

	return out0, out1, err

}

func (_FunctionsCoordinator *FunctionsCoordinatorSession) GetUsdPerUnitLink() (*big.Int, uint8, error) {
	return _FunctionsCoordinator.Contract.GetUsdPerUnitLink(&_FunctionsCoordinator.CallOpts)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCallerSession) GetUsdPerUnitLink() (*big.Int, uint8, error) {
	return _FunctionsCoordinator.Contract.GetUsdPerUnitLink(&_FunctionsCoordinator.CallOpts)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCaller) GetWeiPerUnitLink(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FunctionsCoordinator.contract.Call(opts, &out, "getWeiPerUnitLink")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_FunctionsCoordinator *FunctionsCoordinatorSession) GetWeiPerUnitLink() (*big.Int, error) {
	return _FunctionsCoordinator.Contract.GetWeiPerUnitLink(&_FunctionsCoordinator.CallOpts)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCallerSession) GetWeiPerUnitLink() (*big.Int, error) {
	return _FunctionsCoordinator.Contract.GetWeiPerUnitLink(&_FunctionsCoordinator.CallOpts)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCaller) LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

	error) {
	var out []interface{}
	err := _FunctionsCoordinator.contract.Call(opts, &out, "latestConfigDetails")

	outstruct := new(LatestConfigDetails)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ConfigDigest = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

func (_FunctionsCoordinator *FunctionsCoordinatorSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _FunctionsCoordinator.Contract.LatestConfigDetails(&_FunctionsCoordinator.CallOpts)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCallerSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _FunctionsCoordinator.Contract.LatestConfigDetails(&_FunctionsCoordinator.CallOpts)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCaller) LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

	error) {
	var out []interface{}
	err := _FunctionsCoordinator.contract.Call(opts, &out, "latestConfigDigestAndEpoch")

	outstruct := new(LatestConfigDigestAndEpoch)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ScanLogs = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfigDigest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Epoch = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

func (_FunctionsCoordinator *FunctionsCoordinatorSession) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _FunctionsCoordinator.Contract.LatestConfigDigestAndEpoch(&_FunctionsCoordinator.CallOpts)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCallerSession) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _FunctionsCoordinator.Contract.LatestConfigDigestAndEpoch(&_FunctionsCoordinator.CallOpts)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FunctionsCoordinator.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_FunctionsCoordinator *FunctionsCoordinatorSession) Owner() (common.Address, error) {
	return _FunctionsCoordinator.Contract.Owner(&_FunctionsCoordinator.CallOpts)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCallerSession) Owner() (common.Address, error) {
	return _FunctionsCoordinator.Contract.Owner(&_FunctionsCoordinator.CallOpts)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCaller) Transmitters(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _FunctionsCoordinator.contract.Call(opts, &out, "transmitters")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_FunctionsCoordinator *FunctionsCoordinatorSession) Transmitters() ([]common.Address, error) {
	return _FunctionsCoordinator.Contract.Transmitters(&_FunctionsCoordinator.CallOpts)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCallerSession) Transmitters() ([]common.Address, error) {
	return _FunctionsCoordinator.Contract.Transmitters(&_FunctionsCoordinator.CallOpts)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _FunctionsCoordinator.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_FunctionsCoordinator *FunctionsCoordinatorSession) TypeAndVersion() (string, error) {
	return _FunctionsCoordinator.Contract.TypeAndVersion(&_FunctionsCoordinator.CallOpts)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCallerSession) TypeAndVersion() (string, error) {
	return _FunctionsCoordinator.Contract.TypeAndVersion(&_FunctionsCoordinator.CallOpts)
}

func (_FunctionsCoordinator *FunctionsCoordinatorTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsCoordinator.contract.Transact(opts, "acceptOwnership")
}

func (_FunctionsCoordinator *FunctionsCoordinatorSession) AcceptOwnership() (*types.Transaction, error) {
	return _FunctionsCoordinator.Contract.AcceptOwnership(&_FunctionsCoordinator.TransactOpts)
}

func (_FunctionsCoordinator *FunctionsCoordinatorTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _FunctionsCoordinator.Contract.AcceptOwnership(&_FunctionsCoordinator.TransactOpts)
}

func (_FunctionsCoordinator *FunctionsCoordinatorTransactor) DeleteCommitment(opts *bind.TransactOpts, requestId [32]byte) (*types.Transaction, error) {
	return _FunctionsCoordinator.contract.Transact(opts, "deleteCommitment", requestId)
}

func (_FunctionsCoordinator *FunctionsCoordinatorSession) DeleteCommitment(requestId [32]byte) (*types.Transaction, error) {
	return _FunctionsCoordinator.Contract.DeleteCommitment(&_FunctionsCoordinator.TransactOpts, requestId)
}

func (_FunctionsCoordinator *FunctionsCoordinatorTransactorSession) DeleteCommitment(requestId [32]byte) (*types.Transaction, error) {
	return _FunctionsCoordinator.Contract.DeleteCommitment(&_FunctionsCoordinator.TransactOpts, requestId)
}

func (_FunctionsCoordinator *FunctionsCoordinatorTransactor) OracleWithdraw(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FunctionsCoordinator.contract.Transact(opts, "oracleWithdraw", recipient, amount)
}

func (_FunctionsCoordinator *FunctionsCoordinatorSession) OracleWithdraw(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FunctionsCoordinator.Contract.OracleWithdraw(&_FunctionsCoordinator.TransactOpts, recipient, amount)
}

func (_FunctionsCoordinator *FunctionsCoordinatorTransactorSession) OracleWithdraw(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FunctionsCoordinator.Contract.OracleWithdraw(&_FunctionsCoordinator.TransactOpts, recipient, amount)
}

func (_FunctionsCoordinator *FunctionsCoordinatorTransactor) OracleWithdrawAll(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsCoordinator.contract.Transact(opts, "oracleWithdrawAll")
}

func (_FunctionsCoordinator *FunctionsCoordinatorSession) OracleWithdrawAll() (*types.Transaction, error) {
	return _FunctionsCoordinator.Contract.OracleWithdrawAll(&_FunctionsCoordinator.TransactOpts)
}

func (_FunctionsCoordinator *FunctionsCoordinatorTransactorSession) OracleWithdrawAll() (*types.Transaction, error) {
	return _FunctionsCoordinator.Contract.OracleWithdrawAll(&_FunctionsCoordinator.TransactOpts)
}

func (_FunctionsCoordinator *FunctionsCoordinatorTransactor) SetConfig(opts *bind.TransactOpts, _signers []common.Address, _transmitters []common.Address, _f uint8, _onchainConfig []byte, _offchainConfigVersion uint64, _offchainConfig []byte) (*types.Transaction, error) {
	return _FunctionsCoordinator.contract.Transact(opts, "setConfig", _signers, _transmitters, _f, _onchainConfig, _offchainConfigVersion, _offchainConfig)
}

func (_FunctionsCoordinator *FunctionsCoordinatorSession) SetConfig(_signers []common.Address, _transmitters []common.Address, _f uint8, _onchainConfig []byte, _offchainConfigVersion uint64, _offchainConfig []byte) (*types.Transaction, error) {
	return _FunctionsCoordinator.Contract.SetConfig(&_FunctionsCoordinator.TransactOpts, _signers, _transmitters, _f, _onchainConfig, _offchainConfigVersion, _offchainConfig)
}

func (_FunctionsCoordinator *FunctionsCoordinatorTransactorSession) SetConfig(_signers []common.Address, _transmitters []common.Address, _f uint8, _onchainConfig []byte, _offchainConfigVersion uint64, _offchainConfig []byte) (*types.Transaction, error) {
	return _FunctionsCoordinator.Contract.SetConfig(&_FunctionsCoordinator.TransactOpts, _signers, _transmitters, _f, _onchainConfig, _offchainConfigVersion, _offchainConfig)
}

func (_FunctionsCoordinator *FunctionsCoordinatorTransactor) SetDONPublicKey(opts *bind.TransactOpts, donPublicKey []byte) (*types.Transaction, error) {
	return _FunctionsCoordinator.contract.Transact(opts, "setDONPublicKey", donPublicKey)
}

func (_FunctionsCoordinator *FunctionsCoordinatorSession) SetDONPublicKey(donPublicKey []byte) (*types.Transaction, error) {
	return _FunctionsCoordinator.Contract.SetDONPublicKey(&_FunctionsCoordinator.TransactOpts, donPublicKey)
}

func (_FunctionsCoordinator *FunctionsCoordinatorTransactorSession) SetDONPublicKey(donPublicKey []byte) (*types.Transaction, error) {
	return _FunctionsCoordinator.Contract.SetDONPublicKey(&_FunctionsCoordinator.TransactOpts, donPublicKey)
}

func (_FunctionsCoordinator *FunctionsCoordinatorTransactor) SetThresholdPublicKey(opts *bind.TransactOpts, thresholdPublicKey []byte) (*types.Transaction, error) {
	return _FunctionsCoordinator.contract.Transact(opts, "setThresholdPublicKey", thresholdPublicKey)
}

func (_FunctionsCoordinator *FunctionsCoordinatorSession) SetThresholdPublicKey(thresholdPublicKey []byte) (*types.Transaction, error) {
	return _FunctionsCoordinator.Contract.SetThresholdPublicKey(&_FunctionsCoordinator.TransactOpts, thresholdPublicKey)
}

func (_FunctionsCoordinator *FunctionsCoordinatorTransactorSession) SetThresholdPublicKey(thresholdPublicKey []byte) (*types.Transaction, error) {
	return _FunctionsCoordinator.Contract.SetThresholdPublicKey(&_FunctionsCoordinator.TransactOpts, thresholdPublicKey)
}

func (_FunctionsCoordinator *FunctionsCoordinatorTransactor) StartRequest(opts *bind.TransactOpts, request FunctionsResponseRequestMeta) (*types.Transaction, error) {
	return _FunctionsCoordinator.contract.Transact(opts, "startRequest", request)
}

func (_FunctionsCoordinator *FunctionsCoordinatorSession) StartRequest(request FunctionsResponseRequestMeta) (*types.Transaction, error) {
	return _FunctionsCoordinator.Contract.StartRequest(&_FunctionsCoordinator.TransactOpts, request)
}

func (_FunctionsCoordinator *FunctionsCoordinatorTransactorSession) StartRequest(request FunctionsResponseRequestMeta) (*types.Transaction, error) {
	return _FunctionsCoordinator.Contract.StartRequest(&_FunctionsCoordinator.TransactOpts, request)
}

func (_FunctionsCoordinator *FunctionsCoordinatorTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _FunctionsCoordinator.contract.Transact(opts, "transferOwnership", to)
}

func (_FunctionsCoordinator *FunctionsCoordinatorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _FunctionsCoordinator.Contract.TransferOwnership(&_FunctionsCoordinator.TransactOpts, to)
}

func (_FunctionsCoordinator *FunctionsCoordinatorTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _FunctionsCoordinator.Contract.TransferOwnership(&_FunctionsCoordinator.TransactOpts, to)
}

func (_FunctionsCoordinator *FunctionsCoordinatorTransactor) Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _FunctionsCoordinator.contract.Transact(opts, "transmit", reportContext, report, rs, ss, rawVs)
}

func (_FunctionsCoordinator *FunctionsCoordinatorSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _FunctionsCoordinator.Contract.Transmit(&_FunctionsCoordinator.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_FunctionsCoordinator *FunctionsCoordinatorTransactorSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _FunctionsCoordinator.Contract.Transmit(&_FunctionsCoordinator.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_FunctionsCoordinator *FunctionsCoordinatorTransactor) UpdateConfig(opts *bind.TransactOpts, config FunctionsBillingConfig) (*types.Transaction, error) {
	return _FunctionsCoordinator.contract.Transact(opts, "updateConfig", config)
}

func (_FunctionsCoordinator *FunctionsCoordinatorSession) UpdateConfig(config FunctionsBillingConfig) (*types.Transaction, error) {
	return _FunctionsCoordinator.Contract.UpdateConfig(&_FunctionsCoordinator.TransactOpts, config)
}

func (_FunctionsCoordinator *FunctionsCoordinatorTransactorSession) UpdateConfig(config FunctionsBillingConfig) (*types.Transaction, error) {
	return _FunctionsCoordinator.Contract.UpdateConfig(&_FunctionsCoordinator.TransactOpts, config)
}

type FunctionsCoordinatorCommitmentDeletedIterator struct {
	Event *FunctionsCoordinatorCommitmentDeleted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsCoordinatorCommitmentDeletedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsCoordinatorCommitmentDeleted)
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
		it.Event = new(FunctionsCoordinatorCommitmentDeleted)
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

func (it *FunctionsCoordinatorCommitmentDeletedIterator) Error() error {
	return it.fail
}

func (it *FunctionsCoordinatorCommitmentDeletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsCoordinatorCommitmentDeleted struct {
	RequestId [32]byte
	Raw       types.Log
}

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) FilterCommitmentDeleted(opts *bind.FilterOpts) (*FunctionsCoordinatorCommitmentDeletedIterator, error) {

	logs, sub, err := _FunctionsCoordinator.contract.FilterLogs(opts, "CommitmentDeleted")
	if err != nil {
		return nil, err
	}
	return &FunctionsCoordinatorCommitmentDeletedIterator{contract: _FunctionsCoordinator.contract, event: "CommitmentDeleted", logs: logs, sub: sub}, nil
}

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) WatchCommitmentDeleted(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorCommitmentDeleted) (event.Subscription, error) {

	logs, sub, err := _FunctionsCoordinator.contract.WatchLogs(opts, "CommitmentDeleted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsCoordinatorCommitmentDeleted)
				if err := _FunctionsCoordinator.contract.UnpackLog(event, "CommitmentDeleted", log); err != nil {
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

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) ParseCommitmentDeleted(log types.Log) (*FunctionsCoordinatorCommitmentDeleted, error) {
	event := new(FunctionsCoordinatorCommitmentDeleted)
	if err := _FunctionsCoordinator.contract.UnpackLog(event, "CommitmentDeleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsCoordinatorConfigSetIterator struct {
	Event *FunctionsCoordinatorConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsCoordinatorConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsCoordinatorConfigSet)
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
		it.Event = new(FunctionsCoordinatorConfigSet)
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

func (it *FunctionsCoordinatorConfigSetIterator) Error() error {
	return it.fail
}

func (it *FunctionsCoordinatorConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsCoordinatorConfigSet struct {
	PreviousConfigBlockNumber uint32
	ConfigDigest              [32]byte
	ConfigCount               uint64
	Signers                   []common.Address
	Transmitters              []common.Address
	F                         uint8
	OnchainConfig             []byte
	OffchainConfigVersion     uint64
	OffchainConfig            []byte
	Raw                       types.Log
}

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) FilterConfigSet(opts *bind.FilterOpts) (*FunctionsCoordinatorConfigSetIterator, error) {

	logs, sub, err := _FunctionsCoordinator.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &FunctionsCoordinatorConfigSetIterator{contract: _FunctionsCoordinator.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorConfigSet) (event.Subscription, error) {

	logs, sub, err := _FunctionsCoordinator.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsCoordinatorConfigSet)
				if err := _FunctionsCoordinator.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) ParseConfigSet(log types.Log) (*FunctionsCoordinatorConfigSet, error) {
	event := new(FunctionsCoordinatorConfigSet)
	if err := _FunctionsCoordinator.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsCoordinatorConfigUpdatedIterator struct {
	Event *FunctionsCoordinatorConfigUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsCoordinatorConfigUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsCoordinatorConfigUpdated)
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
		it.Event = new(FunctionsCoordinatorConfigUpdated)
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

func (it *FunctionsCoordinatorConfigUpdatedIterator) Error() error {
	return it.fail
}

func (it *FunctionsCoordinatorConfigUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsCoordinatorConfigUpdated struct {
	Config FunctionsBillingConfig
	Raw    types.Log
}

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) FilterConfigUpdated(opts *bind.FilterOpts) (*FunctionsCoordinatorConfigUpdatedIterator, error) {

	logs, sub, err := _FunctionsCoordinator.contract.FilterLogs(opts, "ConfigUpdated")
	if err != nil {
		return nil, err
	}
	return &FunctionsCoordinatorConfigUpdatedIterator{contract: _FunctionsCoordinator.contract, event: "ConfigUpdated", logs: logs, sub: sub}, nil
}

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) WatchConfigUpdated(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorConfigUpdated) (event.Subscription, error) {

	logs, sub, err := _FunctionsCoordinator.contract.WatchLogs(opts, "ConfigUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsCoordinatorConfigUpdated)
				if err := _FunctionsCoordinator.contract.UnpackLog(event, "ConfigUpdated", log); err != nil {
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

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) ParseConfigUpdated(log types.Log) (*FunctionsCoordinatorConfigUpdated, error) {
	event := new(FunctionsCoordinatorConfigUpdated)
	if err := _FunctionsCoordinator.contract.UnpackLog(event, "ConfigUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsCoordinatorOracleRequestIterator struct {
	Event *FunctionsCoordinatorOracleRequest

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsCoordinatorOracleRequestIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsCoordinatorOracleRequest)
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
		it.Event = new(FunctionsCoordinatorOracleRequest)
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

func (it *FunctionsCoordinatorOracleRequestIterator) Error() error {
	return it.fail
}

func (it *FunctionsCoordinatorOracleRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsCoordinatorOracleRequest struct {
	RequestId          [32]byte
	RequestingContract common.Address
	RequestInitiator   common.Address
	SubscriptionId     uint64
	SubscriptionOwner  common.Address
	Data               []byte
	DataVersion        uint16
	Flags              [32]byte
	CallbackGasLimit   uint64
	Commitment         FunctionsResponseCommitment
	Raw                types.Log
}

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) FilterOracleRequest(opts *bind.FilterOpts, requestId [][32]byte, requestingContract []common.Address) (*FunctionsCoordinatorOracleRequestIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var requestingContractRule []interface{}
	for _, requestingContractItem := range requestingContract {
		requestingContractRule = append(requestingContractRule, requestingContractItem)
	}

	logs, sub, err := _FunctionsCoordinator.contract.FilterLogs(opts, "OracleRequest", requestIdRule, requestingContractRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsCoordinatorOracleRequestIterator{contract: _FunctionsCoordinator.contract, event: "OracleRequest", logs: logs, sub: sub}, nil
}

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) WatchOracleRequest(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorOracleRequest, requestId [][32]byte, requestingContract []common.Address) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var requestingContractRule []interface{}
	for _, requestingContractItem := range requestingContract {
		requestingContractRule = append(requestingContractRule, requestingContractItem)
	}

	logs, sub, err := _FunctionsCoordinator.contract.WatchLogs(opts, "OracleRequest", requestIdRule, requestingContractRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsCoordinatorOracleRequest)
				if err := _FunctionsCoordinator.contract.UnpackLog(event, "OracleRequest", log); err != nil {
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

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) ParseOracleRequest(log types.Log) (*FunctionsCoordinatorOracleRequest, error) {
	event := new(FunctionsCoordinatorOracleRequest)
	if err := _FunctionsCoordinator.contract.UnpackLog(event, "OracleRequest", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsCoordinatorOracleResponseIterator struct {
	Event *FunctionsCoordinatorOracleResponse

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsCoordinatorOracleResponseIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsCoordinatorOracleResponse)
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
		it.Event = new(FunctionsCoordinatorOracleResponse)
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

func (it *FunctionsCoordinatorOracleResponseIterator) Error() error {
	return it.fail
}

func (it *FunctionsCoordinatorOracleResponseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsCoordinatorOracleResponse struct {
	RequestId   [32]byte
	Transmitter common.Address
	Raw         types.Log
}

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) FilterOracleResponse(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsCoordinatorOracleResponseIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsCoordinator.contract.FilterLogs(opts, "OracleResponse", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsCoordinatorOracleResponseIterator{contract: _FunctionsCoordinator.contract, event: "OracleResponse", logs: logs, sub: sub}, nil
}

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) WatchOracleResponse(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorOracleResponse, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsCoordinator.contract.WatchLogs(opts, "OracleResponse", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsCoordinatorOracleResponse)
				if err := _FunctionsCoordinator.contract.UnpackLog(event, "OracleResponse", log); err != nil {
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

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) ParseOracleResponse(log types.Log) (*FunctionsCoordinatorOracleResponse, error) {
	event := new(FunctionsCoordinatorOracleResponse)
	if err := _FunctionsCoordinator.contract.UnpackLog(event, "OracleResponse", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsCoordinatorOwnershipTransferRequestedIterator struct {
	Event *FunctionsCoordinatorOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsCoordinatorOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsCoordinatorOwnershipTransferRequested)
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
		it.Event = new(FunctionsCoordinatorOwnershipTransferRequested)
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

func (it *FunctionsCoordinatorOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *FunctionsCoordinatorOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsCoordinatorOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsCoordinatorOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsCoordinator.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsCoordinatorOwnershipTransferRequestedIterator{contract: _FunctionsCoordinator.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsCoordinator.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsCoordinatorOwnershipTransferRequested)
				if err := _FunctionsCoordinator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) ParseOwnershipTransferRequested(log types.Log) (*FunctionsCoordinatorOwnershipTransferRequested, error) {
	event := new(FunctionsCoordinatorOwnershipTransferRequested)
	if err := _FunctionsCoordinator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsCoordinatorOwnershipTransferredIterator struct {
	Event *FunctionsCoordinatorOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsCoordinatorOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsCoordinatorOwnershipTransferred)
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
		it.Event = new(FunctionsCoordinatorOwnershipTransferred)
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

func (it *FunctionsCoordinatorOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *FunctionsCoordinatorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsCoordinatorOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsCoordinatorOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsCoordinator.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsCoordinatorOwnershipTransferredIterator{contract: _FunctionsCoordinator.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsCoordinator.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsCoordinatorOwnershipTransferred)
				if err := _FunctionsCoordinator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) ParseOwnershipTransferred(log types.Log) (*FunctionsCoordinatorOwnershipTransferred, error) {
	event := new(FunctionsCoordinatorOwnershipTransferred)
	if err := _FunctionsCoordinator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsCoordinatorRequestBilledIterator struct {
	Event *FunctionsCoordinatorRequestBilled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsCoordinatorRequestBilledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsCoordinatorRequestBilled)
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
		it.Event = new(FunctionsCoordinatorRequestBilled)
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

func (it *FunctionsCoordinatorRequestBilledIterator) Error() error {
	return it.fail
}

func (it *FunctionsCoordinatorRequestBilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsCoordinatorRequestBilled struct {
	RequestId         [32]byte
	JuelsPerGas       *big.Int
	L1FeeShareWei     *big.Int
	CallbackCostJuels *big.Int
	DonFeeJuels       *big.Int
	AdminFeeJuels     *big.Int
	OperationFeeJuels *big.Int
	Raw               types.Log
}

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) FilterRequestBilled(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsCoordinatorRequestBilledIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsCoordinator.contract.FilterLogs(opts, "RequestBilled", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsCoordinatorRequestBilledIterator{contract: _FunctionsCoordinator.contract, event: "RequestBilled", logs: logs, sub: sub}, nil
}

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) WatchRequestBilled(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorRequestBilled, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsCoordinator.contract.WatchLogs(opts, "RequestBilled", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsCoordinatorRequestBilled)
				if err := _FunctionsCoordinator.contract.UnpackLog(event, "RequestBilled", log); err != nil {
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

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) ParseRequestBilled(log types.Log) (*FunctionsCoordinatorRequestBilled, error) {
	event := new(FunctionsCoordinatorRequestBilled)
	if err := _FunctionsCoordinator.contract.UnpackLog(event, "RequestBilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsCoordinatorTransmittedIterator struct {
	Event *FunctionsCoordinatorTransmitted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsCoordinatorTransmittedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsCoordinatorTransmitted)
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
		it.Event = new(FunctionsCoordinatorTransmitted)
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

func (it *FunctionsCoordinatorTransmittedIterator) Error() error {
	return it.fail
}

func (it *FunctionsCoordinatorTransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsCoordinatorTransmitted struct {
	ConfigDigest [32]byte
	Epoch        uint32
	Raw          types.Log
}

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) FilterTransmitted(opts *bind.FilterOpts) (*FunctionsCoordinatorTransmittedIterator, error) {

	logs, sub, err := _FunctionsCoordinator.contract.FilterLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return &FunctionsCoordinatorTransmittedIterator{contract: _FunctionsCoordinator.contract, event: "Transmitted", logs: logs, sub: sub}, nil
}

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) WatchTransmitted(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorTransmitted) (event.Subscription, error) {

	logs, sub, err := _FunctionsCoordinator.contract.WatchLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsCoordinatorTransmitted)
				if err := _FunctionsCoordinator.contract.UnpackLog(event, "Transmitted", log); err != nil {
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

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) ParseTransmitted(log types.Log) (*FunctionsCoordinatorTransmitted, error) {
	event := new(FunctionsCoordinatorTransmitted)
	if err := _FunctionsCoordinator.contract.UnpackLog(event, "Transmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LatestConfigDetails struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}
type LatestConfigDigestAndEpoch struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}

func (_FunctionsCoordinator *FunctionsCoordinator) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _FunctionsCoordinator.abi.Events["CommitmentDeleted"].ID:
		return _FunctionsCoordinator.ParseCommitmentDeleted(log)
	case _FunctionsCoordinator.abi.Events["ConfigSet"].ID:
		return _FunctionsCoordinator.ParseConfigSet(log)
	case _FunctionsCoordinator.abi.Events["ConfigUpdated"].ID:
		return _FunctionsCoordinator.ParseConfigUpdated(log)
	case _FunctionsCoordinator.abi.Events["OracleRequest"].ID:
		return _FunctionsCoordinator.ParseOracleRequest(log)
	case _FunctionsCoordinator.abi.Events["OracleResponse"].ID:
		return _FunctionsCoordinator.ParseOracleResponse(log)
	case _FunctionsCoordinator.abi.Events["OwnershipTransferRequested"].ID:
		return _FunctionsCoordinator.ParseOwnershipTransferRequested(log)
	case _FunctionsCoordinator.abi.Events["OwnershipTransferred"].ID:
		return _FunctionsCoordinator.ParseOwnershipTransferred(log)
	case _FunctionsCoordinator.abi.Events["RequestBilled"].ID:
		return _FunctionsCoordinator.ParseRequestBilled(log)
	case _FunctionsCoordinator.abi.Events["Transmitted"].ID:
		return _FunctionsCoordinator.ParseTransmitted(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (FunctionsCoordinatorCommitmentDeleted) Topic() common.Hash {
	return common.HexToHash("0x8a4b97add3359bd6bcf5e82874363670eb5ad0f7615abddbd0ed0a3a98f0f416")
}

func (FunctionsCoordinatorConfigSet) Topic() common.Hash {
	return common.HexToHash("0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05")
}

func (FunctionsCoordinatorConfigUpdated) Topic() common.Hash {
	return common.HexToHash("0x9832846a11905a128363408ace5f9803f449fcb017d3cf7e22a50c6e332931b2")
}

func (FunctionsCoordinatorOracleRequest) Topic() common.Hash {
	return common.HexToHash("0xbf50768ccf13bd0110ca6d53a9c4f1f3271abdd4c24a56878863ed25b20598ff")
}

func (FunctionsCoordinatorOracleResponse) Topic() common.Hash {
	return common.HexToHash("0xc708e0440951fd63499c0f7a73819b469ee5dd3ecc356c0ab4eb7f18389009d9")
}

func (FunctionsCoordinatorOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (FunctionsCoordinatorOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (FunctionsCoordinatorRequestBilled) Topic() common.Hash {
	return common.HexToHash("0x08a4a0761e3c98d288cb4af9342660f49550d83139fb3b762b70d34bed627368")
}

func (FunctionsCoordinatorTransmitted) Topic() common.Hash {
	return common.HexToHash("0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62")
}

func (_FunctionsCoordinator *FunctionsCoordinator) Address() common.Address {
	return _FunctionsCoordinator.address
}

type FunctionsCoordinatorInterface interface {
	EstimateCost(opts *bind.CallOpts, subscriptionId uint64, data []byte, callbackGasLimit uint32, gasPriceWei *big.Int) (*big.Int, error)

	GetAdminFeeJuels(opts *bind.CallOpts) (*big.Int, error)

	GetConfig(opts *bind.CallOpts) (FunctionsBillingConfig, error)

	GetDONFeeJuels(opts *bind.CallOpts, arg0 []byte) (*big.Int, error)

	GetDONPublicKey(opts *bind.CallOpts) ([]byte, error)

	GetOperationFeeJuels(opts *bind.CallOpts) (*big.Int, error)

	GetThresholdPublicKey(opts *bind.CallOpts) ([]byte, error)

	GetUsdPerUnitLink(opts *bind.CallOpts) (*big.Int, uint8, error)

	GetWeiPerUnitLink(opts *bind.CallOpts) (*big.Int, error)

	LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

		error)

	LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

		error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	Transmitters(opts *bind.CallOpts) ([]common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	DeleteCommitment(opts *bind.TransactOpts, requestId [32]byte) (*types.Transaction, error)

	OracleWithdraw(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error)

	OracleWithdrawAll(opts *bind.TransactOpts) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, _signers []common.Address, _transmitters []common.Address, _f uint8, _onchainConfig []byte, _offchainConfigVersion uint64, _offchainConfig []byte) (*types.Transaction, error)

	SetDONPublicKey(opts *bind.TransactOpts, donPublicKey []byte) (*types.Transaction, error)

	SetThresholdPublicKey(opts *bind.TransactOpts, thresholdPublicKey []byte) (*types.Transaction, error)

	StartRequest(opts *bind.TransactOpts, request FunctionsResponseRequestMeta) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error)

	UpdateConfig(opts *bind.TransactOpts, config FunctionsBillingConfig) (*types.Transaction, error)

	FilterCommitmentDeleted(opts *bind.FilterOpts) (*FunctionsCoordinatorCommitmentDeletedIterator, error)

	WatchCommitmentDeleted(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorCommitmentDeleted) (event.Subscription, error)

	ParseCommitmentDeleted(log types.Log) (*FunctionsCoordinatorCommitmentDeleted, error)

	FilterConfigSet(opts *bind.FilterOpts) (*FunctionsCoordinatorConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*FunctionsCoordinatorConfigSet, error)

	FilterConfigUpdated(opts *bind.FilterOpts) (*FunctionsCoordinatorConfigUpdatedIterator, error)

	WatchConfigUpdated(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorConfigUpdated) (event.Subscription, error)

	ParseConfigUpdated(log types.Log) (*FunctionsCoordinatorConfigUpdated, error)

	FilterOracleRequest(opts *bind.FilterOpts, requestId [][32]byte, requestingContract []common.Address) (*FunctionsCoordinatorOracleRequestIterator, error)

	WatchOracleRequest(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorOracleRequest, requestId [][32]byte, requestingContract []common.Address) (event.Subscription, error)

	ParseOracleRequest(log types.Log) (*FunctionsCoordinatorOracleRequest, error)

	FilterOracleResponse(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsCoordinatorOracleResponseIterator, error)

	WatchOracleResponse(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorOracleResponse, requestId [][32]byte) (event.Subscription, error)

	ParseOracleResponse(log types.Log) (*FunctionsCoordinatorOracleResponse, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsCoordinatorOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*FunctionsCoordinatorOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsCoordinatorOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*FunctionsCoordinatorOwnershipTransferred, error)

	FilterRequestBilled(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsCoordinatorRequestBilledIterator, error)

	WatchRequestBilled(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorRequestBilled, requestId [][32]byte) (event.Subscription, error)

	ParseRequestBilled(log types.Log) (*FunctionsCoordinatorRequestBilled, error)

	FilterTransmitted(opts *bind.FilterOpts) (*FunctionsCoordinatorTransmittedIterator, error)

	WatchTransmitted(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorTransmitted) (event.Subscription, error)

	ParseTransmitted(log types.Log) (*FunctionsCoordinatorTransmitted, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
