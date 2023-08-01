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

type IFunctionsBillingRequestBilling struct {
	SubscriptionId   uint64
	Client           common.Address
	CallbackGasLimit uint32
	ExpectedGasPrice *big.Int
}

type IFunctionsCoordinatorRequest struct {
	RequestingContract common.Address
	SubscriptionOwner  common.Address
	Data               []byte
	SubscriptionId     uint64
	DataVersion        uint16
	Flags              [32]byte
	CallbackGasLimit   uint32
}

type IFunctionsRequestCommitment struct {
	AdminFee                  *big.Int
	Coordinator               common.Address
	Client                    common.Address
	SubscriptionId            uint64
	CallbackGasLimit          uint32
	EstimatedTotalCostJuels   *big.Int
	TimeoutTimestamp          *big.Int
	RequestId                 [32]byte
	DonFee                    *big.Int
	GasOverheadBeforeCallback *big.Int
	GasOverheadAfterCallback  *big.Int
}

var FunctionsCoordinatorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"linkToNativeFeed\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"EmptyPublicKey\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptyRequestData\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InconsistentReportData\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"InvalidConfig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"linkWei\",\"type\":\"int256\"}],\"name\":\"InvalidLinkWeiPrice\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSubscription\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"MustBeSubOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoTransmittersSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByRouter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByRouterOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentTooLarge\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReportInvalid\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RouterMustBeSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedPublicKeyChange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedSender\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnsupportedRequestDataVersion\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"signerPayment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"transmitterPayment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalCost\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"enumFulfillResult\",\"name\":\"result\",\"type\":\"uint8\"}],\"name\":\"BillingEnd\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"uint96\",\"name\":\"adminFee\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"estimatedTotalCostJuels\",\"type\":\"uint96\"},{\"internalType\":\"uint40\",\"name\":\"timeoutTimestamp\",\"type\":\"uint40\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint80\",\"name\":\"donFee\",\"type\":\"uint80\"},{\"internalType\":\"uint40\",\"name\":\"gasOverheadBeforeCallback\",\"type\":\"uint40\"},{\"internalType\":\"uint40\",\"name\":\"gasOverheadAfterCallback\",\"type\":\"uint40\"}],\"indexed\":false,\"internalType\":\"structIFunctionsRequest.Commitment\",\"name\":\"commitment\",\"type\":\"tuple\"}],\"name\":\"BillingStart\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"CommitmentDeleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maxCallbackGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"feedStalenessSeconds\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasOverheadBeforeCallback\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasOverheadAfterCallback\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"requestTimeoutSeconds\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint80\",\"name\":\"donFee\",\"type\":\"uint80\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"maxSupportedRequestDataVersion\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fulfillmentGasPriceOverEstimationBP\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"fallbackNativePerUnitLink\",\"type\":\"int256\"}],\"name\":\"ConfigChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"CostExceedsCommitment\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"InsufficientGasProvided\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"InsufficientSubscriptionBalance\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"InvalidRequestID\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requestingContract\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requestInitiator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"subscriptionOwner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"dataVersion\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"flags\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"callbackGasLimit\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"commitment\",\"type\":\"bytes\"}],\"name\":\"OracleRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"OracleResponse\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"deleteCommitment\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"node\",\"type\":\"address\"}],\"name\":\"deleteNodePublicKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"}],\"name\":\"estimateCost\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"expectedGasPrice\",\"type\":\"uint256\"}],\"internalType\":\"structIFunctionsBilling.RequestBilling\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"getAdminFee\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllNodePublicKeys\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"},{\"internalType\":\"bytes[]\",\"name\":\"\",\"type\":\"bytes[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfig\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"maxCallbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"feedStalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasOverheadBeforeCallback\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasOverheadAfterCallback\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"requestTimeoutSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint80\",\"name\":\"donFee\",\"type\":\"uint80\"},{\"internalType\":\"uint16\",\"name\":\"maxSupportedRequestDataVersion\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"fulfillmentGasPriceOverEstimationBP\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"fallbackNativePerUnitLink\",\"type\":\"int256\"},{\"internalType\":\"address\",\"name\":\"linkPriceFeed\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfigHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"config\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"expectedGasPrice\",\"type\":\"uint256\"}],\"internalType\":\"structIFunctionsBilling.RequestBilling\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"getDONFee\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDONPublicKey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFeedData\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getThresholdPublicKey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"oracleWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"requestingContract\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"subscriptionOwner\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"dataVersion\",\"type\":\"uint16\"},{\"internalType\":\"bytes32\",\"name\":\"flags\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"}],\"internalType\":\"structIFunctionsCoordinator.Request\",\"name\":\"request\",\"type\":\"tuple\"}],\"name\":\"sendRequest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint96\",\"name\":\"estimatedCost\",\"type\":\"uint96\"},{\"components\":[{\"internalType\":\"uint96\",\"name\":\"adminFee\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"estimatedTotalCostJuels\",\"type\":\"uint96\"},{\"internalType\":\"uint40\",\"name\":\"timeoutTimestamp\",\"type\":\"uint40\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint80\",\"name\":\"donFee\",\"type\":\"uint80\"},{\"internalType\":\"uint40\",\"name\":\"gasOverheadBeforeCallback\",\"type\":\"uint40\"},{\"internalType\":\"uint40\",\"name\":\"gasOverheadAfterCallback\",\"type\":\"uint40\"}],\"internalType\":\"structIFunctionsRequest.Commitment\",\"name\":\"commitment\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"_f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"_onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"_offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"_offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"donPublicKey\",\"type\":\"bytes\"}],\"name\":\"setDONPublicKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"node\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"}],\"name\":\"setNodePublicKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"thresholdPublicKey\",\"type\":\"bytes\"}],\"name\":\"setThresholdPublicKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"transmitters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"name\":\"updateConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60c06040523480156200001157600080fd5b5060405162005c2238038062005c22833981016040819052620000349162000462565b828282828260013380600081620000925760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c557620000c58162000154565b50505015156080526001600160a01b038216620000f557604051632530e88560e11b815260040160405180910390fd5b6001600160a01b03821660a0526200010d81620001ff565b805160209091012060075550600d80546001600160a01b039092166c01000000000000000000000000026001600160601b0390921691909117905550620006349350505050565b336001600160a01b03821603620001ae5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000089565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60008060008060008060008060008980602001905181019062000223919062000572565b985098509850985098509850985098509850600081136200025b576040516321ea67b360e11b81526004810182905260240162000089565b604080516101208101825263ffffffff808c168083528b8216602084018190528b83168486018190528b841660608601819052938b16608086018190526001600160501b038b1660a0870181905261ffff8b1660c0880181905260e088018b905261010090970189905260098054600160f01b9098026001600160f01b03600160a01b909302600160a01b600160f01b0319600160801b90950294909416600160801b600160f01b03196c0100000000000000000000000090990263ffffffff60601b196801000000000000000090970296909616600160401b600160801b03196401000000009098026001600160401b0319909b169098179990991795909516959095179290921794909416949094179290921792909216179055600a839055600b829055517f3332571032ee658b30d567e63b2443e9d921162625e0f3fe86968a586c1a090f906200041b908b908b908b908b908b908b908b908b908b9063ffffffff998a1681529789166020890152958816604088015293871660608701529190951660808501526001600160501b039490941660a084015261ffff9390931660c083015260e08201929092526101008101919091526101200190565b60405180910390a150505050505050505050565b80516001600160a01b03811681146200044757600080fd5b919050565b634e487b7160e01b600052604160045260246000fd5b6000806000606084860312156200047857600080fd5b62000483846200042f565b602085810151919450906001600160401b0380821115620004a357600080fd5b818701915087601f830112620004b857600080fd5b815181811115620004cd57620004cd6200044c565b604051601f8201601f19908116603f01168101908382118183101715620004f857620004f86200044c565b816040528281528a868487010111156200051157600080fd5b600093505b8284101562000535578484018601518185018701529285019262000516565b600086848301015280975050505050505062000554604085016200042f565b90509250925092565b805163ffffffff811681146200044757600080fd5b60008060008060008060008060006101208a8c0312156200059257600080fd5b6200059d8a6200055d565b9850620005ad60208b016200055d565b9750620005bd60408b016200055d565b9650620005cd60608b016200055d565b9550620005dd60808b016200055d565b60a08b01519095506001600160501b0381168114620005fb57600080fd5b60c08b015190945061ffff811681146200061457600080fd5b8093505060e08a015191506101008a015190509295985092959850929598565b60805160a0516155976200068b60003960008181610a1b01528181610e3e01528181610f1c01528181610ffb01528181611a9a01528181611ce401528181612eb70152613b57015260006114d701526155976000f3fe608060405234801561001057600080fd5b50600436106101ae5760003560e01c80638cc6acce116100ee578063b1dc65a411610097578063d328a91e11610071578063d328a91e14610537578063de9dfa461461053f578063e3d0e71214610552578063f2fde38b1461056557600080fd5b8063b1dc65a4146103f6578063c3f909d414610409578063d227d2451461050757600080fd5b8063aecc12c5116100c8578063aecc12c5146103ac578063af1296d3146103ce578063afcb95d7146103d657600080fd5b80638cc6acce1461035f5780638da5cb5b146103725780639883c10d1461039a57600080fd5b806379ba50971161015b578063814118341161013557806381411834146102b257806381f1b938146102c757806381ff7048146102cf57806385b214cf1461033c57600080fd5b806379ba5097146102845780637f15e1661461028c578063807560311461029f57600080fd5b8063466557181161018c578063466557181461022d578063533989871461025b57806366316d8d1461027157600080fd5b8063083a5466146101b3578063181f5a77146101c857806326ceabac1461021a575b600080fd5b6101c66101c1366004613e41565b610578565b005b6102046040518060400160405280601c81526020017f46756e6374696f6e7320436f6f7264696e61746f722076312e302e300000000081525081565b6040516102119190613ee7565b60405180910390f35b6101c6610228366004613f1c565b6105cd565b61024061023b3660046140a2565b61066f565b60405169ffffffffffffffffffff9091168152602001610211565b61026361069d565b6040516102119291906141c9565b6101c661027f366004614273565b610921565b6101c6610adb565b6101c661029a366004613e41565b610bdd565b6101c66102ad3660046142ac565b610c2d565b6102ba610ce4565b6040516102119190614301565b610204610d53565b61031960015460025463ffffffff74010000000000000000000000000000000000000000830481169378010000000000000000000000000000000000000000000000009093041691565b6040805163ffffffff948516815293909216602084015290820152606001610211565b61034f61034a366004614314565b610e24565b6040519015158152602001610211565b6101c661036d36600461432d565b610f04565b60005460405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610211565b6007545b604051908152602001610211565b6103bf6103ba366004614362565b610f89565b604051610211939291906144c3565b61039e611271565b604080516001815260006020820181905291810191909152606001610211565b6101c6610404366004614531565b611364565b600954600a54600b54600d546040805163ffffffff80871682526401000000008704811660208301526801000000000000000087048116928201929092526c01000000000000000000000000808704831660608301527001000000000000000000000000000000008704909216608082015274010000000000000000000000000000000000000000860469ffffffffffffffffffff1660a08201527e0100000000000000000000000000000000000000000000000000000000000090950461ffff1660c086015260e08501939093526101008401919091520473ffffffffffffffffffffffffffffffffffffffff1661012082015261014001610211565b61051a6105153660046145e8565b611a95565b6040516bffffffffffffffffffffffff9091168152602001610211565b610204611c89565b61051a61054d3660046140a2565b611ce0565b6101c661056036600461470a565b611d78565b6101c6610573366004613f1c565b61279c565b6105806127ad565b60008190036105bb576040517f4f42be3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60106105c8828483614878565b505050565b60005473ffffffffffffffffffffffffffffffffffffffff1633148061060857503373ffffffffffffffffffffffffffffffffffffffff8216145b61063e576040517fed6dd19b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff81166000908152600f6020526040812061066c91613d86565b50565b60095474010000000000000000000000000000000000000000900469ffffffffffffffffffff165b92915050565b6060806000600680548060200260200160405190810160405280929190818152602001828054801561070557602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116106da575b505050505090506000815167ffffffffffffffff81111561072857610728613f39565b60405190808252806020026020018201604052801561075b57816020015b60608152602001906001900390816107465790505b50905060005b825181101561091757600f600084838151811061078057610780614993565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002080546107cd906147d7565b9050600003610808576040517f4f42be3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600f600084838151811061081e5761081e614993565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020805461086b906147d7565b80601f0160208091040260200160405190810160405280929190818152602001828054610897906147d7565b80156108e45780601f106108b9576101008083540402835291602001916108e4565b820191906000526020600020905b8154815290600101906020018083116108c757829003601f168201915b50505050508282815181106108fb576108fb614993565b602002602001018190525080610910906149f1565b9050610761565b5090939092509050565b610929612830565b806bffffffffffffffffffffffff1660000361095f5750336000908152600c60205260409020546bffffffffffffffffffffffff165b336000908152600c60205260409020546bffffffffffffffffffffffff808316911610156109b9576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b336000908152600c6020526040812080548392906109e69084906bffffffffffffffffffffffff16614a29565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055506000610a3d7f000000000000000000000000000000000000000000000000000000000000000090565b6040517f66316d8d00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff85811660048301526bffffffffffffffffffffffff85166024830152919250908216906366316d8d90604401600060405180830381600087803b158015610abe57600080fd5b505af1158015610ad2573d6000803e3d6000fd5b50505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610b61576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610be56127ad565b6000819003610c20576040517f4f42be3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600e6105c8828483614878565b60005473ffffffffffffffffffffffffffffffffffffffff16331480610c785750610c57336129c1565b8015610c7857503373ffffffffffffffffffffffffffffffffffffffff8416145b610cae576040517fed6dd19b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff83166000908152600f60205260409020610cde828483614878565b50505050565b60606006805480602002602001604051908101604052809291908181526020018280548015610d4957602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610d1e575b5050505050905090565b606060108054610d62906147d7565b9050600003610d9d576040517f4f42be3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60108054610daa906147d7565b80601f0160208091040260200160405190810160405280929190818152602001828054610dd6906147d7565b8015610d495780601f10610df857610100808354040283529160200191610d49565b820191906000526020600020905b815481529060010190602001808311610e0657509395945050505050565b60003373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610e95576040517fc41a5b0900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260086020526040902054610eb057506000919050565b60008281526008602052604080822091909155517f8a4b97add3359bd6bcf5e82874363670eb5ad0f7615abddbd0ed0a3a98f0f41690610ef39084815260200190565b60405180910390a15060015b919050565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610f73576040517fc41a5b0900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610f7c81612aaa565b8051602090910120600755565b6040805161016081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e0810182905261010081018290526101208101829052610140810182905281903373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614611052576040517fc41a5b0900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61105f6040850185614a55565b9050600003611099576040517ec1cfc000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600060405180608001604052808660600160208101906110b99190614aba565b67ffffffffffffffff1681526020908101906110d790880188613f1c565b73ffffffffffffffffffffffffffffffffffffffff16815260200161110260e0880160c08901614ad7565b63ffffffff1681523a602090910152905061116e6111236040870187614a55565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506111689250505060a0880160808901614b04565b83612db3565b919550935091506111826020860186613f1c565b73ffffffffffffffffffffffffffffffffffffffff16847fd88ef982f3624222a23ef83ed2192cbdfba8e1e474a6cb1d315736d08479d2e3326111cb60808a0160608b01614aba565b6111db60408b0160208c01613f1c565b6111e860408c018c614a55565b6111f860a08e0160808f01614b04565b8d60a001358e60c00160208101906112109190614ad7565b8c6040516020016112219190614b21565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815290829052611261999897969594939291614b30565b60405180910390a3509193909250565b600954600d54604080517ffeaf968c0000000000000000000000000000000000000000000000000000000081529051600093640100000000900463ffffffff169283151592859283926c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff169163feaf968c9160048083019260a09291908290030181865afa15801561130b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061132f9190614bfd565b5093505092505082801561135157506113488142614c4d565b8463ffffffff16105b1561135c57600b5491505b509392505050565b60005a604080516020601f8b018190048102820181019092528981529192508a3591818c0135916113ba91849163ffffffff851691908e908e90819084018382808284376000920191909152506132ef92505050565b6113f0576040517f0be3632800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805183815262ffffff600884901c1660208201527fb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62910160405180910390a16040805160608101825260025480825260035460ff808216602085015261010090910416928201929092529083146114c5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f636f6e666967446967657374206d69736d6174636800000000000000000000006044820152606401610b58565b6114d38b8b8b8b8b8b6132f8565b60007f000000000000000000000000000000000000000000000000000000000000000015611530576002826020015183604001516115119190614c60565b61151b9190614ca8565b611526906001614c60565b60ff169050611546565b6020820151611540906001614c60565b60ff1690505b8881146115af576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f77726f6e67206e756d626572206f66207369676e6174757265730000000000006044820152606401610b58565b888714611618576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f7369676e617475726573206f7574206f6620726567697374726174696f6e00006044820152606401610b58565b3360009081526004602090815260408083208151808301909252805460ff8082168452929391929184019161010090910416600281111561165b5761165b614cca565b600281111561166c5761166c614cca565b905250905060028160200151600281111561168957611689614cca565b1480156116d057506006816000015160ff16815481106116ab576116ab614993565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff1633145b611736576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f756e617574686f72697a6564207472616e736d697474657200000000000000006044820152606401610b58565b5050505050611743613dc0565b6000808a8a604051611756929190614cf9565b60405190819003812061176d918e90602001614d09565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181528282528051602091820120838301909252600080845290830152915060005b89811015611a775760006001848984602081106117d6576117d6614993565b6117e391901a601b614c60565b8e8e868181106117f5576117f5614993565b905060200201358d8d8781811061180e5761180e614993565b905060200201356040516000815260200160405260405161184b949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa15801561186d573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff811660009081526004602090815290849020838501909452835460ff808216855292965092945084019161010090041660028111156118ed576118ed614cca565b60028111156118fe576118fe614cca565b905250925060018360200151600281111561191b5761191b614cca565b14611982576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f61646472657373206e6f7420617574686f72697a656420746f207369676e00006044820152606401610b58565b8251600090879060ff16601f811061199c5761199c614993565b602002015173ffffffffffffffffffffffffffffffffffffffff1614611a1e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f6e6f6e2d756e69717565207369676e61747572650000000000000000000000006044820152606401610b58565b8086846000015160ff16601f8110611a3857611a38614993565b73ffffffffffffffffffffffffffffffffffffffff9092166020929092020152611a63600186614c60565b94505080611a70906149f1565b90506117b7565b505050611a88833383858e8e6133a6565b5050505050505050505050565b6000807f00000000000000000000000000000000000000000000000000000000000000006040517f10fc49c100000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8916600482015263ffffffff8616602482015290915073ffffffffffffffffffffffffffffffffffffffff8216906310fc49c19060440160006040518083038186803b158015611b3757600080fd5b505afa158015611b4b573d6000803e3d6000fd5b50505050620f4240831115611b8c576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600060405180608001604052808967ffffffffffffffff1681526020013373ffffffffffffffffffffffffffffffffffffffff1681526020018663ffffffff1681526020018581525090506000611c1a88888080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525086925061066f915050565b69ffffffffffffffffffff1690506000611c6b89898080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250879250611ce0915050565b9050611c7987878484613704565b9450505050505b95945050505050565b6060600e8054611c98906147d7565b9050600003611cd3576040517f4f42be3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600e8054610daa906147d7565b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16632a905ccc6040518163ffffffff1660e01b8152600401602060405180830381865afa158015611d4d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611d719190614d28565b9392505050565b855185518560ff16601f831115611deb576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f746f6f206d616e79207369676e657273000000000000000000000000000000006044820152606401610b58565b80600003611e55576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f66206d75737420626520706f73697469766500000000000000000000000000006044820152606401610b58565b818314611ee3576040517f89a61989000000000000000000000000000000000000000000000000000000008152602060048201526024808201527f6f7261636c6520616464726573736573206f7574206f6620726567697374726160448201527f74696f6e000000000000000000000000000000000000000000000000000000006064820152608401610b58565b611eee816003614d45565b8311611f56576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f6661756c74792d6f7261636c65206620746f6f206869676800000000000000006044820152606401610b58565b611f5e6127ad565b6040805160c0810182528a8152602081018a905260ff89169181018290526060810188905267ffffffffffffffff8716608082015260a0810186905290611fa59088613876565b6005541561215a57600554600090611fbf90600190614c4d565b9050600060058281548110611fd657611fd6614993565b60009182526020822001546006805473ffffffffffffffffffffffffffffffffffffffff9092169350908490811061201057612010614993565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff85811684526004909252604080842080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00009081169091559290911680845292208054909116905560058054919250908061209057612090614d5c565b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff000000000000000000000000000000000000000016905501905560068054806120f9576120f9614d5c565b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff000000000000000000000000000000000000000016905501905550611fa5915050565b60005b8151518110156125c1576000600460008460000151848151811061218357612183614993565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff16825281019190915260400160002054610100900460ff1660028111156121cd576121cd614cca565b14612234576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265706561746564207369676e657220616464726573730000000000000000006044820152606401610b58565b6040805180820190915260ff8216815260016020820152825180516004916000918590811061226557612265614993565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040016000208251815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082168117835592840151919283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000161761010083600281111561230657612306614cca565b0217905550600091506123169050565b600460008460200151848151811061233057612330614993565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff16825281019190915260400160002054610100900460ff16600281111561237a5761237a614cca565b146123e1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f7265706561746564207472616e736d69747465722061646472657373000000006044820152606401610b58565b6040805180820190915260ff82168152602081016002815250600460008460200151848151811061241457612414614993565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040016000208251815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082168117835592840151919283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000016176101008360028111156124b5576124b5614cca565b0217905550508251805160059250839081106124d3576124d3614993565b602090810291909101810151825460018101845560009384529282902090920180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909316929092179091558201518051600691908390811061254f5761254f614993565b60209081029190910181015182546001810184556000938452919092200180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909216919091179055806125b9816149f1565b91505061215d565b506040810151600380547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff909216919091179055600180547fffffffff00000000ffffffffffffffffffffffffffffffffffffffffffffffff8116780100000000000000000000000000000000000000000000000063ffffffff438116820292909217808555920481169291829160149161267991849174010000000000000000000000000000000000000000900416614d8b565b92506101000a81548163ffffffff021916908363ffffffff1602179055506126d84630600160149054906101000a900463ffffffff1663ffffffff16856000015186602001518760400151886060015189608001518a60a00151613893565b600281905582518051600380547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1661010060ff9093169290920291909117905560015460208501516040808701516060880151608089015160a08a015193517f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e059861278f988b9891977401000000000000000000000000000000000000000090920463ffffffff16969095919491939192614da8565b60405180910390a1611a88565b6127a46127ad565b61066c8161393e565b60005473ffffffffffffffffffffffffffffffffffffffff16331461282e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610b58565b565b600061283a610ce4565b90508051600003612877576040517f30274b3a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8051600d54600091612896916bffffffffffffffffffffffff16614e3e565b905060005b82518110156129625781600c60008584815181106128bb576128bb614993565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282829054906101000a90046bffffffffffffffffffffffff166129239190614e69565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055508061295b906149f1565b905061289b565b50815161296f9082614e8e565b600d805460009061298f9084906bffffffffffffffffffffffff16614a29565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055505050565b6000806006805480602002602001604051908101604052809291908181526020018280548015612a2757602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116129fc575b5050505050905060005b8151811015612aa0578373ffffffffffffffffffffffffffffffffffffffff16828281518110612a6357612a63614993565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1603612a90575060019392505050565b612a99816149f1565b9050612a31565b5060009392505050565b600080600080600080600080600089806020019051810190612acc9190614ec9565b98509850985098509850985098509850985060008113612b1b576040517f43d4cf6600000000000000000000000000000000000000000000000000000000815260048101829052602401610b58565b604080516101208101825263ffffffff808c168083528b8216602084018190528b83168486018190528b841660608601819052938b166080860181905269ffffffffffffffffffff8b1660a0870181905261ffff8b1660c0880181905260e088018b9052610100909701899052600980547e010000000000000000000000000000000000000000000000000000000000009098027dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff740100000000000000000000000000000000000000009093027fffff00000000000000000000ffffffffffffffffffffffffffffffffffffffff700100000000000000000000000000000000909502949094167fffff0000000000000000000000000000ffffffffffffffffffffffffffffffff6c010000000000000000000000009099027fffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffffff68010000000000000000909702969096167fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff6401000000009098027fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000909b169098179990991795909516959095179290921794909416949094179290921792909216179055600a839055600b829055517f3332571032ee658b30d567e63b2443e9d921162625e0f3fe86968a586c1a090f90612d9f908b908b908b908b908b908b908b908b908b9063ffffffff998a16815297891660208901529588166040880152938716606087015291909516608085015269ffffffffffffffffffff9490941660a084015261ffff9390931660c083015260e08201929092526101008101919091526101200190565b60405180910390a150505050505050505050565b6040805161016081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e08101829052610100810182905261012081018290526101408101829052819060095461ffff7e0100000000000000000000000000000000000000000000000000000000000090910481169086161115612e73576040517fdada758700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000612e7f878661066f565b90506000612e8d8887611ce0565b90506000612eb1876040015188606001518569ffffffffffffffffffff1685613704565b905060007f000000000000000000000000000000000000000000000000000000000000000088516040517fa47c769600000000000000000000000000000000000000000000000000000000815267ffffffffffffffff9091166004820152909150600090819073ffffffffffffffffffffffffffffffffffffffff84169063a47c769690602401600060405180830381865afa158015612f55573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052612f9b9190810190614f80565b50505060208c01518c516040517f674603d000000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff928316600482015267ffffffffffffffff90911660248201529294509092506000919085169063674603d090604401606060405180830381865afa15801561302c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906130509190615077565b509150506bffffffffffffffffffffffff851661306d8385614a29565b6bffffffffffffffffffffffff1610156130b3576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000613132308d602001518e600001518560016130d091906150c9565b6040805173ffffffffffffffffffffffffffffffffffffffff958616602080830191909152949095168582015267ffffffffffffffff928316606086015291166080808501919091528151808503909101815260a09093019052815191012090565b90506000604051806101600160405280896bffffffffffffffffffffffff1681526020013073ffffffffffffffffffffffffffffffffffffffff1681526020018e6020015173ffffffffffffffffffffffffffffffffffffffff1681526020018e6000015167ffffffffffffffff1681526020018e6040015163ffffffff168152602001886bffffffffffffffffffffffff168152602001600960000160109054906101000a900463ffffffff1663ffffffff16426131f191906150ea565b64ffffffffff168152602080820185905269ffffffffffffffffffff8c1660408084019190915260095463ffffffff680100000000000000008204811660608601526c0100000000000000000000000090910416608090930192909252905191925061325f91839101614b21565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181528282528051602091820120600086815260089092529190205582907febc94a6fea57cb363e5d876da6bbc9072372f6839c0b5f86e890ee07aa8db6be906132d2908490614b21565b60405180910390a2909e959d509b50939950505050505050505050565b60019392505050565b6000613305826020614d45565b613310856020614d45565b61331c886101446150ea565b61332691906150ea565b61333091906150ea565b61333b9060006150ea565b9050368114610ad2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f63616c6c64617461206c656e677468206d69736d6174636800000000000000006044820152606401610b58565b6060808080806133b8868801886151d8565b84519499509297509095509350915015806133d557508351855114155b806133e257508251855114155b806133ef57508151855114155b806133fc57508051855114155b15613433576040517f0be3632800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b85518110156136f65760006134cb87838151811061345657613456614993565b602002602001015187848151811061347057613470614993565b602002602001015187858151811061348a5761348a614993565b60200260200101518786815181106134a4576134a4614993565b60200260200101518787815181106134be576134be614993565b6020026020010151613a33565b905060008160068111156134e1576134e1614cca565b14806134fe575060018160068111156134fc576134fc614cca565b145b156135595786828151811061351557613515614993565b60209081029190910181015160405133815290917fc708e0440951fd63499c0f7a73819b469ee5dd3ecc356c0ab4eb7f18389009d9910160405180910390a26136e5565b600281600681111561356d5761356d614cca565b036135bd5786828151811061358457613584614993565b60200260200101517fa1c120e327c9ad8b075793878c88d59b8934b97ae37117faa3bb21616237f7be60405160405180910390a26136e5565b60038160068111156135d1576135d1614cca565b03613621578682815181106135e8576135e8614993565b60200260200101517f357a60574f143e144d62ddc89636fee45ba5e2f8c9a30b9f91e28f8dfe5a56e060405160405180910390a26136e5565b600581600681111561363557613635614cca565b036136855786828151811061364c5761364c614993565b60200260200101517f689f49eacf6db85082542f4b7873f4a519744a4e093fd002310545bfb4354cc860405160405180910390a26136e5565b600481600681111561369957613699614cca565b036136e5578682815181106136b0576136b0614993565b60200260200101517fb83d9204606ee3780daa106ad026c4c2853f5af126edbb8fe15e813e62dbd94d60405160405180910390a25b506136ef816149f1565b9050613436565b505050505050505050505050565b60008061370f611271565b90506000811361374e576040517f43d4cf6600000000000000000000000000000000000000000000000000000000815260048101829052602401610b58565b60095460009087906137849063ffffffff6c01000000000000000000000000820481169168010000000000000000900416614d8b565b61378e9190614d8b565b63ffffffff1690506000612710600960010154886137ac9190614d45565b6137b691906152aa565b6137c090886150ea565b9050600083836137d884670de0b6b3a7640000614d45565b6137e29190614d45565b6137ec91906152aa565b9050600061380b6bffffffffffffffffffffffff808916908a166150ea565b9050613823816b033b2e3c9fd0803ce8000000614c4d565b82111561385c576040517fe80fa38100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61386681836150ea565b955050505050505b949350505050565b6000613880610ce4565b51111561388f5761388f612830565b5050565b6000808a8a8a8a8a8a8a8a8a6040516020016138b7999897969594939291906152be565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179150509998505050505050505050565b3373ffffffffffffffffffffffffffffffffffffffff8216036139bd576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610b58565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60008083806020019051810190613a4a9190615368565b905080604051602001613a5d9190614b21565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152918152815160209283012060008a8152600890935291205414613aaf576002915050611c80565b6000878152600860205260408120819055613ac8611271565b905060008113613b07576040517f43d4cf6600000000000000000000000000000000000000000000000000000000815260048101829052602401610b58565b600081613b1c3a670de0b6b3a7640000614d45565b613b2691906152aa565b90506000836101400151846101200151613b40919061543b565b613b519064ffffffffff1683614d45565b905060007f000000000000000000000000000000000000000000000000000000000000000090506000808273ffffffffffffffffffffffffffffffffffffffff1663a1d4c8288d8d888b610100015169ffffffffffffffffffff1689613bb79190614e69565b338d6040518763ffffffff1660e01b8152600401613bda96959493929190615459565b60408051808303816000875af1158015613bf8573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613c1c91906154d5565b9092509050613c2b8185614e69565b336000908152600c6020526040812080547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff938416179055610100890151600d805469ffffffffffffffffffff90921693909291613c9991859116614e69565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055508c7f30e1e299f390e4f3fa75e4489fed477e86183937a4a21a0ae7b137a4dff0353588606001518961010001518488613d009190614e69565b8b516101008d015169ffffffffffffffffffff16613d1e888c614e69565b613d289190614e69565b613d329190614e69565b8760ff166006811115613d4757613d47614cca565b604051613d58959493929190615504565b60405180910390a28160ff166006811115613d7557613d75614cca565b9d9c50505050505050505050505050565b508054613d92906147d7565b6000825580601f10613da2575050565b601f01602090049060005260206000209081019061066c9190613ddf565b604051806103e00160405280601f906020820280368337509192915050565b5b80821115613df45760008155600101613de0565b5090565b60008083601f840112613e0a57600080fd5b50813567ffffffffffffffff811115613e2257600080fd5b602083019150836020828501011115613e3a57600080fd5b9250929050565b60008060208385031215613e5457600080fd5b823567ffffffffffffffff811115613e6b57600080fd5b613e7785828601613df8565b90969095509350505050565b6000815180845260005b81811015613ea957602081850181015186830182015201613e8d565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b602081526000611d716020830184613e83565b73ffffffffffffffffffffffffffffffffffffffff8116811461066c57600080fd5b600060208284031215613f2e57600080fd5b8135611d7181613efa565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610160810167ffffffffffffffff81118282101715613f8c57613f8c613f39565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715613fd957613fd9613f39565b604052919050565b600082601f830112613ff257600080fd5b813567ffffffffffffffff81111561400c5761400c613f39565b61403d60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601613f92565b81815284602083860101111561405257600080fd5b816020850160208301376000918101602001919091529392505050565b67ffffffffffffffff8116811461066c57600080fd5b8035610eff8161406f565b63ffffffff8116811461066c57600080fd5b60008082840360a08112156140b657600080fd5b833567ffffffffffffffff808211156140ce57600080fd5b6140da87838801613fe1565b945060807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08401121561410c57600080fd5b6040519250608083019150828210818311171561412b5761412b613f39565b50604052602084013561413d8161406f565b8152604084013561414d81613efa565b6020820152606084013561416081614090565b60408201526080939093013560608401525092909150565b600081518084526020808501945080840160005b838110156141be57815173ffffffffffffffffffffffffffffffffffffffff168752958201959082019060010161418c565b509495945050505050565b6040815260006141dc6040830185614178565b6020838203818501528185518084528284019150828160051b85010183880160005b8381101561424a577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0878403018552614238838351613e83565b948601949250908501906001016141fe565b50909998505050505050505050565b6bffffffffffffffffffffffff8116811461066c57600080fd5b6000806040838503121561428657600080fd5b823561429181613efa565b915060208301356142a181614259565b809150509250929050565b6000806000604084860312156142c157600080fd5b83356142cc81613efa565b9250602084013567ffffffffffffffff8111156142e857600080fd5b6142f486828701613df8565b9497909650939450505050565b602081526000611d716020830184614178565b60006020828403121561432657600080fd5b5035919050565b60006020828403121561433f57600080fd5b813567ffffffffffffffff81111561435657600080fd5b61386e84828501613fe1565b60006020828403121561437457600080fd5b813567ffffffffffffffff81111561438b57600080fd5b820160e08185031215611d7157600080fd5b80516bffffffffffffffffffffffff16825260208101516143d6602084018273ffffffffffffffffffffffffffffffffffffffff169052565b5060408101516143fe604084018273ffffffffffffffffffffffffffffffffffffffff169052565b50606081015161441a606084018267ffffffffffffffff169052565b506080810151614432608084018263ffffffff169052565b5060a081015161445260a08401826bffffffffffffffffffffffff169052565b5060c081015161446b60c084018264ffffffffff169052565b5060e081015160e0830152610100808201516144948285018269ffffffffffffffffffff169052565b50506101208181015164ffffffffff81168483015250506101408181015164ffffffffff811684830152610cde565b8381526bffffffffffffffffffffffff831660208201526101a0810161386e604083018461439d565b60008083601f8401126144fe57600080fd5b50813567ffffffffffffffff81111561451657600080fd5b6020830191508360208260051b8501011115613e3a57600080fd5b60008060008060008060008060e0898b03121561454d57600080fd5b606089018a81111561455e57600080fd5b8998503567ffffffffffffffff8082111561457857600080fd5b6145848c838d01613df8565b909950975060808b013591508082111561459d57600080fd5b6145a98c838d016144ec565b909750955060a08b01359150808211156145c257600080fd5b506145cf8b828c016144ec565b999c989b50969995989497949560c00135949350505050565b60008060008060006080868803121561460057600080fd5b853561460b8161406f565b9450602086013567ffffffffffffffff81111561462757600080fd5b61463388828901613df8565b909550935050604086013561464781614090565b949793965091946060013592915050565b600067ffffffffffffffff82111561467257614672613f39565b5060051b60200190565b600082601f83011261468d57600080fd5b813560206146a261469d83614658565b613f92565b82815260059290921b840181019181810190868411156146c157600080fd5b8286015b848110156146e55780356146d881613efa565b83529183019183016146c5565b509695505050505050565b60ff8116811461066c57600080fd5b8035610eff816146f0565b60008060008060008060c0878903121561472357600080fd5b863567ffffffffffffffff8082111561473b57600080fd5b6147478a838b0161467c565b9750602089013591508082111561475d57600080fd5b6147698a838b0161467c565b965061477760408a016146ff565b9550606089013591508082111561478d57600080fd5b6147998a838b01613fe1565b94506147a760808a01614085565b935060a08901359150808211156147bd57600080fd5b506147ca89828a01613fe1565b9150509295509295509295565b600181811c908216806147eb57607f821691505b602082108103614824577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f8211156105c857600081815260208120601f850160051c810160208610156148515750805b601f850160051c820191505b818110156148705782815560010161485d565b505050505050565b67ffffffffffffffff83111561489057614890613f39565b6148a48361489e83546147d7565b8361482a565b6000601f8411600181146148f657600085156148c05750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b17835561498c565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b828110156149455786850135825560209485019460019092019101614925565b5086821015614980577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203614a2257614a226149c2565b5060010190565b6bffffffffffffffffffffffff828116828216039080821115614a4e57614a4e6149c2565b5092915050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112614a8a57600080fd5b83018035915067ffffffffffffffff821115614aa557600080fd5b602001915036819003821315613e3a57600080fd5b600060208284031215614acc57600080fd5b8135611d718161406f565b600060208284031215614ae957600080fd5b8135611d7181614090565b61ffff8116811461066c57600080fd5b600060208284031215614b1657600080fd5b8135611d7181614af4565b6101608101610697828461439d565b73ffffffffffffffffffffffffffffffffffffffff8a8116825267ffffffffffffffff8a16602083015288166040820152610100606082018190528101869052600061012087898285013760008189850101527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f890116830161ffff881660808501528660a0850152614bcd60c085018763ffffffff169052565b818482030160e0850152613d7582820186613e83565b805169ffffffffffffffffffff81168114610eff57600080fd5b600080600080600060a08688031215614c1557600080fd5b614c1e86614be3565b9450602086015193506040860151925060608601519150614c4160808701614be3565b90509295509295909350565b81810381811115610697576106976149c2565b60ff8181168382160190811115610697576106976149c2565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b600060ff831680614cbb57614cbb614c79565b8060ff84160491505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b8183823760009101908152919050565b828152606082602083013760800192915050565b8051610eff81614259565b600060208284031215614d3a57600080fd5b8151611d7181614259565b8082028115828204841417610697576106976149c2565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b63ffffffff818116838216019080821115614a4e57614a4e6149c2565b600061012063ffffffff808d1684528b6020850152808b16604085015250806060840152614dd88184018a614178565b90508281036080840152614dec8189614178565b905060ff871660a084015282810360c0840152614e098187613e83565b905067ffffffffffffffff851660e0840152828103610100840152614e2e8185613e83565b9c9b505050505050505050505050565b60006bffffffffffffffffffffffff80841680614e5d57614e5d614c79565b92169190910492915050565b6bffffffffffffffffffffffff818116838216019080821115614a4e57614a4e6149c2565b6bffffffffffffffffffffffff818116838216028082169190828114614eb657614eb66149c2565b505092915050565b8051610eff81614090565b60008060008060008060008060006101208a8c031215614ee857600080fd5b8951614ef381614090565b60208b0151909950614f0481614090565b60408b0151909850614f1581614090565b60608b0151909750614f2681614090565b60808b0151909650614f3781614090565b9450614f4560a08b01614be3565b935060c08a0151614f5581614af4565b8093505060e08a015191506101008a015190509295985092959850929598565b8051610eff81613efa565b600080600080600060a08688031215614f9857600080fd5b8551614fa381614259565b80955050602080870151614fb681614259565b6040880151909550614fc781613efa565b6060880151909450614fd881613efa565b608088015190935067ffffffffffffffff811115614ff557600080fd5b8701601f8101891361500657600080fd5b805161501461469d82614658565b81815260059190911b8201830190838101908b83111561503357600080fd5b928401925b8284101561505a57835161504b81613efa565b82529284019290840190615038565b80955050505050509295509295909350565b8051610eff8161406f565b60008060006060848603121561508c57600080fd5b8351801515811461509c57600080fd5b60208501519093506150ad8161406f565b60408501519092506150be8161406f565b809150509250925092565b67ffffffffffffffff818116838216019080821115614a4e57614a4e6149c2565b80820180821115610697576106976149c2565b600082601f83011261510e57600080fd5b8135602061511e61469d83614658565b82815260059290921b8401810191818101908684111561513d57600080fd5b8286015b848110156146e55780358352918301918301615141565b600082601f83011261516957600080fd5b8135602061517961469d83614658565b82815260059290921b8401810191818101908684111561519857600080fd5b8286015b848110156146e557803567ffffffffffffffff8111156151bc5760008081fd5b6151ca8986838b0101613fe1565b84525091830191830161519c565b600080600080600060a086880312156151f057600080fd5b853567ffffffffffffffff8082111561520857600080fd5b61521489838a016150fd565b9650602088013591508082111561522a57600080fd5b61523689838a01615158565b9550604088013591508082111561524c57600080fd5b61525889838a01615158565b9450606088013591508082111561526e57600080fd5b61527a89838a01615158565b9350608088013591508082111561529057600080fd5b5061529d88828901615158565b9150509295509295909350565b6000826152b9576152b9614c79565b500490565b60006101208b835273ffffffffffffffffffffffffffffffffffffffff8b16602084015267ffffffffffffffff808b1660408501528160608501526153058285018b614178565b91508382036080850152615319828a614178565b915060ff881660a085015283820360c08501526153368288613e83565b90861660e08501528381036101008501529050614e2e8185613e83565b805164ffffffffff81168114610eff57600080fd5b6000610160828403121561537b57600080fd5b615383613f68565b61538c83614d1d565b815261539a60208401614f75565b60208201526153ab60408401614f75565b60408201526153bc6060840161506c565b60608201526153cd60808401614ebe565b60808201526153de60a08401614d1d565b60a08201526153ef60c08401615353565b60c082015260e083015160e082015261010061540c818501614be3565b9082015261012061541e848201615353565b90820152610140615430848201615353565b908201529392505050565b64ffffffffff818116838216019080821115614a4e57614a4e6149c2565b600061020080835261546d8184018a613e83565b905082810360208401526154818189613e83565b6bffffffffffffffffffffffff88811660408601528716606085015273ffffffffffffffffffffffffffffffffffffffff8616608085015291506154ca905060a083018461439d565b979650505050505050565b600080604083850312156154e857600080fd5b82516154f3816146f0565b60208401519092506142a181614259565b67ffffffffffffffff8616815269ffffffffffffffffffff851660208201526bffffffffffffffffffffffff84811660408301528316606082015260a081016007831061557a577f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b826080830152969550505050505056fea164736f6c6343000813000a",
}

var FunctionsCoordinatorABI = FunctionsCoordinatorMetaData.ABI

var FunctionsCoordinatorBin = FunctionsCoordinatorMetaData.Bin

func DeployFunctionsCoordinator(auth *bind.TransactOpts, backend bind.ContractBackend, router common.Address, config []byte, linkToNativeFeed common.Address) (common.Address, *types.Transaction, *FunctionsCoordinator, error) {
	parsed, err := FunctionsCoordinatorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(FunctionsCoordinatorBin), backend, router, config, linkToNativeFeed)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &FunctionsCoordinator{FunctionsCoordinatorCaller: FunctionsCoordinatorCaller{contract: contract}, FunctionsCoordinatorTransactor: FunctionsCoordinatorTransactor{contract: contract}, FunctionsCoordinatorFilterer: FunctionsCoordinatorFilterer{contract: contract}}, nil
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

func (_FunctionsCoordinator *FunctionsCoordinatorCaller) EstimateCost(opts *bind.CallOpts, subscriptionId uint64, data []byte, callbackGasLimit uint32, gasPrice *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _FunctionsCoordinator.contract.Call(opts, &out, "estimateCost", subscriptionId, data, callbackGasLimit, gasPrice)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_FunctionsCoordinator *FunctionsCoordinatorSession) EstimateCost(subscriptionId uint64, data []byte, callbackGasLimit uint32, gasPrice *big.Int) (*big.Int, error) {
	return _FunctionsCoordinator.Contract.EstimateCost(&_FunctionsCoordinator.CallOpts, subscriptionId, data, callbackGasLimit, gasPrice)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCallerSession) EstimateCost(subscriptionId uint64, data []byte, callbackGasLimit uint32, gasPrice *big.Int) (*big.Int, error) {
	return _FunctionsCoordinator.Contract.EstimateCost(&_FunctionsCoordinator.CallOpts, subscriptionId, data, callbackGasLimit, gasPrice)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCaller) GetAdminFee(opts *bind.CallOpts, arg0 []byte, arg1 IFunctionsBillingRequestBilling) (*big.Int, error) {
	var out []interface{}
	err := _FunctionsCoordinator.contract.Call(opts, &out, "getAdminFee", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_FunctionsCoordinator *FunctionsCoordinatorSession) GetAdminFee(arg0 []byte, arg1 IFunctionsBillingRequestBilling) (*big.Int, error) {
	return _FunctionsCoordinator.Contract.GetAdminFee(&_FunctionsCoordinator.CallOpts, arg0, arg1)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCallerSession) GetAdminFee(arg0 []byte, arg1 IFunctionsBillingRequestBilling) (*big.Int, error) {
	return _FunctionsCoordinator.Contract.GetAdminFee(&_FunctionsCoordinator.CallOpts, arg0, arg1)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCaller) GetAllNodePublicKeys(opts *bind.CallOpts) ([]common.Address, [][]byte, error) {
	var out []interface{}
	err := _FunctionsCoordinator.contract.Call(opts, &out, "getAllNodePublicKeys")

	if err != nil {
		return *new([]common.Address), *new([][]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	out1 := *abi.ConvertType(out[1], new([][]byte)).(*[][]byte)

	return out0, out1, err

}

func (_FunctionsCoordinator *FunctionsCoordinatorSession) GetAllNodePublicKeys() ([]common.Address, [][]byte, error) {
	return _FunctionsCoordinator.Contract.GetAllNodePublicKeys(&_FunctionsCoordinator.CallOpts)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCallerSession) GetAllNodePublicKeys() ([]common.Address, [][]byte, error) {
	return _FunctionsCoordinator.Contract.GetAllNodePublicKeys(&_FunctionsCoordinator.CallOpts)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCaller) GetConfig(opts *bind.CallOpts) (GetConfig,

	error) {
	var out []interface{}
	err := _FunctionsCoordinator.contract.Call(opts, &out, "getConfig")

	outstruct := new(GetConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.MaxCallbackGasLimit = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.FeedStalenessSeconds = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.GasOverheadBeforeCallback = *abi.ConvertType(out[2], new(uint32)).(*uint32)
	outstruct.GasOverheadAfterCallback = *abi.ConvertType(out[3], new(uint32)).(*uint32)
	outstruct.RequestTimeoutSeconds = *abi.ConvertType(out[4], new(uint32)).(*uint32)
	outstruct.DonFee = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.MaxSupportedRequestDataVersion = *abi.ConvertType(out[6], new(uint16)).(*uint16)
	outstruct.FulfillmentGasPriceOverEstimationBP = *abi.ConvertType(out[7], new(*big.Int)).(**big.Int)
	outstruct.FallbackNativePerUnitLink = *abi.ConvertType(out[8], new(*big.Int)).(**big.Int)
	outstruct.LinkPriceFeed = *abi.ConvertType(out[9], new(common.Address)).(*common.Address)

	return *outstruct, err

}

func (_FunctionsCoordinator *FunctionsCoordinatorSession) GetConfig() (GetConfig,

	error) {
	return _FunctionsCoordinator.Contract.GetConfig(&_FunctionsCoordinator.CallOpts)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCallerSession) GetConfig() (GetConfig,

	error) {
	return _FunctionsCoordinator.Contract.GetConfig(&_FunctionsCoordinator.CallOpts)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCaller) GetConfigHash(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _FunctionsCoordinator.contract.Call(opts, &out, "getConfigHash")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_FunctionsCoordinator *FunctionsCoordinatorSession) GetConfigHash() ([32]byte, error) {
	return _FunctionsCoordinator.Contract.GetConfigHash(&_FunctionsCoordinator.CallOpts)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCallerSession) GetConfigHash() ([32]byte, error) {
	return _FunctionsCoordinator.Contract.GetConfigHash(&_FunctionsCoordinator.CallOpts)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCaller) GetDONFee(opts *bind.CallOpts, arg0 []byte, arg1 IFunctionsBillingRequestBilling) (*big.Int, error) {
	var out []interface{}
	err := _FunctionsCoordinator.contract.Call(opts, &out, "getDONFee", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_FunctionsCoordinator *FunctionsCoordinatorSession) GetDONFee(arg0 []byte, arg1 IFunctionsBillingRequestBilling) (*big.Int, error) {
	return _FunctionsCoordinator.Contract.GetDONFee(&_FunctionsCoordinator.CallOpts, arg0, arg1)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCallerSession) GetDONFee(arg0 []byte, arg1 IFunctionsBillingRequestBilling) (*big.Int, error) {
	return _FunctionsCoordinator.Contract.GetDONFee(&_FunctionsCoordinator.CallOpts, arg0, arg1)
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

func (_FunctionsCoordinator *FunctionsCoordinatorCaller) GetFeedData(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FunctionsCoordinator.contract.Call(opts, &out, "getFeedData")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_FunctionsCoordinator *FunctionsCoordinatorSession) GetFeedData() (*big.Int, error) {
	return _FunctionsCoordinator.Contract.GetFeedData(&_FunctionsCoordinator.CallOpts)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCallerSession) GetFeedData() (*big.Int, error) {
	return _FunctionsCoordinator.Contract.GetFeedData(&_FunctionsCoordinator.CallOpts)
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

func (_FunctionsCoordinator *FunctionsCoordinatorTransactor) DeleteNodePublicKey(opts *bind.TransactOpts, node common.Address) (*types.Transaction, error) {
	return _FunctionsCoordinator.contract.Transact(opts, "deleteNodePublicKey", node)
}

func (_FunctionsCoordinator *FunctionsCoordinatorSession) DeleteNodePublicKey(node common.Address) (*types.Transaction, error) {
	return _FunctionsCoordinator.Contract.DeleteNodePublicKey(&_FunctionsCoordinator.TransactOpts, node)
}

func (_FunctionsCoordinator *FunctionsCoordinatorTransactorSession) DeleteNodePublicKey(node common.Address) (*types.Transaction, error) {
	return _FunctionsCoordinator.Contract.DeleteNodePublicKey(&_FunctionsCoordinator.TransactOpts, node)
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

func (_FunctionsCoordinator *FunctionsCoordinatorTransactor) SendRequest(opts *bind.TransactOpts, request IFunctionsCoordinatorRequest) (*types.Transaction, error) {
	return _FunctionsCoordinator.contract.Transact(opts, "sendRequest", request)
}

func (_FunctionsCoordinator *FunctionsCoordinatorSession) SendRequest(request IFunctionsCoordinatorRequest) (*types.Transaction, error) {
	return _FunctionsCoordinator.Contract.SendRequest(&_FunctionsCoordinator.TransactOpts, request)
}

func (_FunctionsCoordinator *FunctionsCoordinatorTransactorSession) SendRequest(request IFunctionsCoordinatorRequest) (*types.Transaction, error) {
	return _FunctionsCoordinator.Contract.SendRequest(&_FunctionsCoordinator.TransactOpts, request)
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

func (_FunctionsCoordinator *FunctionsCoordinatorTransactor) SetNodePublicKey(opts *bind.TransactOpts, node common.Address, publicKey []byte) (*types.Transaction, error) {
	return _FunctionsCoordinator.contract.Transact(opts, "setNodePublicKey", node, publicKey)
}

func (_FunctionsCoordinator *FunctionsCoordinatorSession) SetNodePublicKey(node common.Address, publicKey []byte) (*types.Transaction, error) {
	return _FunctionsCoordinator.Contract.SetNodePublicKey(&_FunctionsCoordinator.TransactOpts, node, publicKey)
}

func (_FunctionsCoordinator *FunctionsCoordinatorTransactorSession) SetNodePublicKey(node common.Address, publicKey []byte) (*types.Transaction, error) {
	return _FunctionsCoordinator.Contract.SetNodePublicKey(&_FunctionsCoordinator.TransactOpts, node, publicKey)
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

func (_FunctionsCoordinator *FunctionsCoordinatorTransactor) UpdateConfig(opts *bind.TransactOpts, config []byte) (*types.Transaction, error) {
	return _FunctionsCoordinator.contract.Transact(opts, "updateConfig", config)
}

func (_FunctionsCoordinator *FunctionsCoordinatorSession) UpdateConfig(config []byte) (*types.Transaction, error) {
	return _FunctionsCoordinator.Contract.UpdateConfig(&_FunctionsCoordinator.TransactOpts, config)
}

func (_FunctionsCoordinator *FunctionsCoordinatorTransactorSession) UpdateConfig(config []byte) (*types.Transaction, error) {
	return _FunctionsCoordinator.Contract.UpdateConfig(&_FunctionsCoordinator.TransactOpts, config)
}

type FunctionsCoordinatorBillingEndIterator struct {
	Event *FunctionsCoordinatorBillingEnd

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsCoordinatorBillingEndIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsCoordinatorBillingEnd)
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
		it.Event = new(FunctionsCoordinatorBillingEnd)
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

func (it *FunctionsCoordinatorBillingEndIterator) Error() error {
	return it.fail
}

func (it *FunctionsCoordinatorBillingEndIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsCoordinatorBillingEnd struct {
	RequestId          [32]byte
	SubscriptionId     uint64
	SignerPayment      *big.Int
	TransmitterPayment *big.Int
	TotalCost          *big.Int
	Result             uint8
	Raw                types.Log
}

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) FilterBillingEnd(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsCoordinatorBillingEndIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsCoordinator.contract.FilterLogs(opts, "BillingEnd", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsCoordinatorBillingEndIterator{contract: _FunctionsCoordinator.contract, event: "BillingEnd", logs: logs, sub: sub}, nil
}

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) WatchBillingEnd(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorBillingEnd, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsCoordinator.contract.WatchLogs(opts, "BillingEnd", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsCoordinatorBillingEnd)
				if err := _FunctionsCoordinator.contract.UnpackLog(event, "BillingEnd", log); err != nil {
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

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) ParseBillingEnd(log types.Log) (*FunctionsCoordinatorBillingEnd, error) {
	event := new(FunctionsCoordinatorBillingEnd)
	if err := _FunctionsCoordinator.contract.UnpackLog(event, "BillingEnd", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsCoordinatorBillingStartIterator struct {
	Event *FunctionsCoordinatorBillingStart

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsCoordinatorBillingStartIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsCoordinatorBillingStart)
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
		it.Event = new(FunctionsCoordinatorBillingStart)
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

func (it *FunctionsCoordinatorBillingStartIterator) Error() error {
	return it.fail
}

func (it *FunctionsCoordinatorBillingStartIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsCoordinatorBillingStart struct {
	RequestId  [32]byte
	Commitment IFunctionsRequestCommitment
	Raw        types.Log
}

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) FilterBillingStart(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsCoordinatorBillingStartIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsCoordinator.contract.FilterLogs(opts, "BillingStart", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsCoordinatorBillingStartIterator{contract: _FunctionsCoordinator.contract, event: "BillingStart", logs: logs, sub: sub}, nil
}

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) WatchBillingStart(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorBillingStart, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsCoordinator.contract.WatchLogs(opts, "BillingStart", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsCoordinatorBillingStart)
				if err := _FunctionsCoordinator.contract.UnpackLog(event, "BillingStart", log); err != nil {
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

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) ParseBillingStart(log types.Log) (*FunctionsCoordinatorBillingStart, error) {
	event := new(FunctionsCoordinatorBillingStart)
	if err := _FunctionsCoordinator.contract.UnpackLog(event, "BillingStart", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
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

type FunctionsCoordinatorConfigChangedIterator struct {
	Event *FunctionsCoordinatorConfigChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsCoordinatorConfigChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsCoordinatorConfigChanged)
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
		it.Event = new(FunctionsCoordinatorConfigChanged)
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

func (it *FunctionsCoordinatorConfigChangedIterator) Error() error {
	return it.fail
}

func (it *FunctionsCoordinatorConfigChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsCoordinatorConfigChanged struct {
	MaxCallbackGasLimit                 uint32
	FeedStalenessSeconds                uint32
	GasOverheadBeforeCallback           uint32
	GasOverheadAfterCallback            uint32
	RequestTimeoutSeconds               uint32
	DonFee                              *big.Int
	MaxSupportedRequestDataVersion      uint16
	FulfillmentGasPriceOverEstimationBP *big.Int
	FallbackNativePerUnitLink           *big.Int
	Raw                                 types.Log
}

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) FilterConfigChanged(opts *bind.FilterOpts) (*FunctionsCoordinatorConfigChangedIterator, error) {

	logs, sub, err := _FunctionsCoordinator.contract.FilterLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return &FunctionsCoordinatorConfigChangedIterator{contract: _FunctionsCoordinator.contract, event: "ConfigChanged", logs: logs, sub: sub}, nil
}

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorConfigChanged) (event.Subscription, error) {

	logs, sub, err := _FunctionsCoordinator.contract.WatchLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsCoordinatorConfigChanged)
				if err := _FunctionsCoordinator.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
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

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) ParseConfigChanged(log types.Log) (*FunctionsCoordinatorConfigChanged, error) {
	event := new(FunctionsCoordinatorConfigChanged)
	if err := _FunctionsCoordinator.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
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

type FunctionsCoordinatorCostExceedsCommitmentIterator struct {
	Event *FunctionsCoordinatorCostExceedsCommitment

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsCoordinatorCostExceedsCommitmentIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsCoordinatorCostExceedsCommitment)
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
		it.Event = new(FunctionsCoordinatorCostExceedsCommitment)
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

func (it *FunctionsCoordinatorCostExceedsCommitmentIterator) Error() error {
	return it.fail
}

func (it *FunctionsCoordinatorCostExceedsCommitmentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsCoordinatorCostExceedsCommitment struct {
	RequestId [32]byte
	Raw       types.Log
}

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) FilterCostExceedsCommitment(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsCoordinatorCostExceedsCommitmentIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsCoordinator.contract.FilterLogs(opts, "CostExceedsCommitment", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsCoordinatorCostExceedsCommitmentIterator{contract: _FunctionsCoordinator.contract, event: "CostExceedsCommitment", logs: logs, sub: sub}, nil
}

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) WatchCostExceedsCommitment(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorCostExceedsCommitment, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsCoordinator.contract.WatchLogs(opts, "CostExceedsCommitment", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsCoordinatorCostExceedsCommitment)
				if err := _FunctionsCoordinator.contract.UnpackLog(event, "CostExceedsCommitment", log); err != nil {
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

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) ParseCostExceedsCommitment(log types.Log) (*FunctionsCoordinatorCostExceedsCommitment, error) {
	event := new(FunctionsCoordinatorCostExceedsCommitment)
	if err := _FunctionsCoordinator.contract.UnpackLog(event, "CostExceedsCommitment", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsCoordinatorInsufficientGasProvidedIterator struct {
	Event *FunctionsCoordinatorInsufficientGasProvided

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsCoordinatorInsufficientGasProvidedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsCoordinatorInsufficientGasProvided)
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
		it.Event = new(FunctionsCoordinatorInsufficientGasProvided)
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

func (it *FunctionsCoordinatorInsufficientGasProvidedIterator) Error() error {
	return it.fail
}

func (it *FunctionsCoordinatorInsufficientGasProvidedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsCoordinatorInsufficientGasProvided struct {
	RequestId [32]byte
	Raw       types.Log
}

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) FilterInsufficientGasProvided(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsCoordinatorInsufficientGasProvidedIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsCoordinator.contract.FilterLogs(opts, "InsufficientGasProvided", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsCoordinatorInsufficientGasProvidedIterator{contract: _FunctionsCoordinator.contract, event: "InsufficientGasProvided", logs: logs, sub: sub}, nil
}

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) WatchInsufficientGasProvided(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorInsufficientGasProvided, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsCoordinator.contract.WatchLogs(opts, "InsufficientGasProvided", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsCoordinatorInsufficientGasProvided)
				if err := _FunctionsCoordinator.contract.UnpackLog(event, "InsufficientGasProvided", log); err != nil {
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

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) ParseInsufficientGasProvided(log types.Log) (*FunctionsCoordinatorInsufficientGasProvided, error) {
	event := new(FunctionsCoordinatorInsufficientGasProvided)
	if err := _FunctionsCoordinator.contract.UnpackLog(event, "InsufficientGasProvided", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsCoordinatorInsufficientSubscriptionBalanceIterator struct {
	Event *FunctionsCoordinatorInsufficientSubscriptionBalance

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsCoordinatorInsufficientSubscriptionBalanceIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsCoordinatorInsufficientSubscriptionBalance)
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
		it.Event = new(FunctionsCoordinatorInsufficientSubscriptionBalance)
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

func (it *FunctionsCoordinatorInsufficientSubscriptionBalanceIterator) Error() error {
	return it.fail
}

func (it *FunctionsCoordinatorInsufficientSubscriptionBalanceIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsCoordinatorInsufficientSubscriptionBalance struct {
	RequestId [32]byte
	Raw       types.Log
}

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) FilterInsufficientSubscriptionBalance(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsCoordinatorInsufficientSubscriptionBalanceIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsCoordinator.contract.FilterLogs(opts, "InsufficientSubscriptionBalance", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsCoordinatorInsufficientSubscriptionBalanceIterator{contract: _FunctionsCoordinator.contract, event: "InsufficientSubscriptionBalance", logs: logs, sub: sub}, nil
}

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) WatchInsufficientSubscriptionBalance(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorInsufficientSubscriptionBalance, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsCoordinator.contract.WatchLogs(opts, "InsufficientSubscriptionBalance", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsCoordinatorInsufficientSubscriptionBalance)
				if err := _FunctionsCoordinator.contract.UnpackLog(event, "InsufficientSubscriptionBalance", log); err != nil {
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

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) ParseInsufficientSubscriptionBalance(log types.Log) (*FunctionsCoordinatorInsufficientSubscriptionBalance, error) {
	event := new(FunctionsCoordinatorInsufficientSubscriptionBalance)
	if err := _FunctionsCoordinator.contract.UnpackLog(event, "InsufficientSubscriptionBalance", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsCoordinatorInvalidRequestIDIterator struct {
	Event *FunctionsCoordinatorInvalidRequestID

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsCoordinatorInvalidRequestIDIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsCoordinatorInvalidRequestID)
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
		it.Event = new(FunctionsCoordinatorInvalidRequestID)
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

func (it *FunctionsCoordinatorInvalidRequestIDIterator) Error() error {
	return it.fail
}

func (it *FunctionsCoordinatorInvalidRequestIDIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsCoordinatorInvalidRequestID struct {
	RequestId [32]byte
	Raw       types.Log
}

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) FilterInvalidRequestID(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsCoordinatorInvalidRequestIDIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsCoordinator.contract.FilterLogs(opts, "InvalidRequestID", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsCoordinatorInvalidRequestIDIterator{contract: _FunctionsCoordinator.contract, event: "InvalidRequestID", logs: logs, sub: sub}, nil
}

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) WatchInvalidRequestID(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorInvalidRequestID, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsCoordinator.contract.WatchLogs(opts, "InvalidRequestID", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsCoordinatorInvalidRequestID)
				if err := _FunctionsCoordinator.contract.UnpackLog(event, "InvalidRequestID", log); err != nil {
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

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) ParseInvalidRequestID(log types.Log) (*FunctionsCoordinatorInvalidRequestID, error) {
	event := new(FunctionsCoordinatorInvalidRequestID)
	if err := _FunctionsCoordinator.contract.UnpackLog(event, "InvalidRequestID", log); err != nil {
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
	Commitment         []byte
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

type GetConfig struct {
	MaxCallbackGasLimit                 uint32
	FeedStalenessSeconds                uint32
	GasOverheadBeforeCallback           uint32
	GasOverheadAfterCallback            uint32
	RequestTimeoutSeconds               uint32
	DonFee                              *big.Int
	MaxSupportedRequestDataVersion      uint16
	FulfillmentGasPriceOverEstimationBP *big.Int
	FallbackNativePerUnitLink           *big.Int
	LinkPriceFeed                       common.Address
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
	case _FunctionsCoordinator.abi.Events["BillingEnd"].ID:
		return _FunctionsCoordinator.ParseBillingEnd(log)
	case _FunctionsCoordinator.abi.Events["BillingStart"].ID:
		return _FunctionsCoordinator.ParseBillingStart(log)
	case _FunctionsCoordinator.abi.Events["CommitmentDeleted"].ID:
		return _FunctionsCoordinator.ParseCommitmentDeleted(log)
	case _FunctionsCoordinator.abi.Events["ConfigChanged"].ID:
		return _FunctionsCoordinator.ParseConfigChanged(log)
	case _FunctionsCoordinator.abi.Events["ConfigSet"].ID:
		return _FunctionsCoordinator.ParseConfigSet(log)
	case _FunctionsCoordinator.abi.Events["CostExceedsCommitment"].ID:
		return _FunctionsCoordinator.ParseCostExceedsCommitment(log)
	case _FunctionsCoordinator.abi.Events["InsufficientGasProvided"].ID:
		return _FunctionsCoordinator.ParseInsufficientGasProvided(log)
	case _FunctionsCoordinator.abi.Events["InsufficientSubscriptionBalance"].ID:
		return _FunctionsCoordinator.ParseInsufficientSubscriptionBalance(log)
	case _FunctionsCoordinator.abi.Events["InvalidRequestID"].ID:
		return _FunctionsCoordinator.ParseInvalidRequestID(log)
	case _FunctionsCoordinator.abi.Events["OracleRequest"].ID:
		return _FunctionsCoordinator.ParseOracleRequest(log)
	case _FunctionsCoordinator.abi.Events["OracleResponse"].ID:
		return _FunctionsCoordinator.ParseOracleResponse(log)
	case _FunctionsCoordinator.abi.Events["OwnershipTransferRequested"].ID:
		return _FunctionsCoordinator.ParseOwnershipTransferRequested(log)
	case _FunctionsCoordinator.abi.Events["OwnershipTransferred"].ID:
		return _FunctionsCoordinator.ParseOwnershipTransferred(log)
	case _FunctionsCoordinator.abi.Events["Transmitted"].ID:
		return _FunctionsCoordinator.ParseTransmitted(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (FunctionsCoordinatorBillingEnd) Topic() common.Hash {
	return common.HexToHash("0x30e1e299f390e4f3fa75e4489fed477e86183937a4a21a0ae7b137a4dff03535")
}

func (FunctionsCoordinatorBillingStart) Topic() common.Hash {
	return common.HexToHash("0xebc94a6fea57cb363e5d876da6bbc9072372f6839c0b5f86e890ee07aa8db6be")
}

func (FunctionsCoordinatorCommitmentDeleted) Topic() common.Hash {
	return common.HexToHash("0x8a4b97add3359bd6bcf5e82874363670eb5ad0f7615abddbd0ed0a3a98f0f416")
}

func (FunctionsCoordinatorConfigChanged) Topic() common.Hash {
	return common.HexToHash("0x3332571032ee658b30d567e63b2443e9d921162625e0f3fe86968a586c1a090f")
}

func (FunctionsCoordinatorConfigSet) Topic() common.Hash {
	return common.HexToHash("0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05")
}

func (FunctionsCoordinatorCostExceedsCommitment) Topic() common.Hash {
	return common.HexToHash("0x689f49eacf6db85082542f4b7873f4a519744a4e093fd002310545bfb4354cc8")
}

func (FunctionsCoordinatorInsufficientGasProvided) Topic() common.Hash {
	return common.HexToHash("0x357a60574f143e144d62ddc89636fee45ba5e2f8c9a30b9f91e28f8dfe5a56e0")
}

func (FunctionsCoordinatorInsufficientSubscriptionBalance) Topic() common.Hash {
	return common.HexToHash("0xb83d9204606ee3780daa106ad026c4c2853f5af126edbb8fe15e813e62dbd94d")
}

func (FunctionsCoordinatorInvalidRequestID) Topic() common.Hash {
	return common.HexToHash("0xa1c120e327c9ad8b075793878c88d59b8934b97ae37117faa3bb21616237f7be")
}

func (FunctionsCoordinatorOracleRequest) Topic() common.Hash {
	return common.HexToHash("0xd88ef982f3624222a23ef83ed2192cbdfba8e1e474a6cb1d315736d08479d2e3")
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

func (FunctionsCoordinatorTransmitted) Topic() common.Hash {
	return common.HexToHash("0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62")
}

func (_FunctionsCoordinator *FunctionsCoordinator) Address() common.Address {
	return _FunctionsCoordinator.address
}

type FunctionsCoordinatorInterface interface {
	EstimateCost(opts *bind.CallOpts, subscriptionId uint64, data []byte, callbackGasLimit uint32, gasPrice *big.Int) (*big.Int, error)

	GetAdminFee(opts *bind.CallOpts, arg0 []byte, arg1 IFunctionsBillingRequestBilling) (*big.Int, error)

	GetAllNodePublicKeys(opts *bind.CallOpts) ([]common.Address, [][]byte, error)

	GetConfig(opts *bind.CallOpts) (GetConfig,

		error)

	GetConfigHash(opts *bind.CallOpts) ([32]byte, error)

	GetDONFee(opts *bind.CallOpts, arg0 []byte, arg1 IFunctionsBillingRequestBilling) (*big.Int, error)

	GetDONPublicKey(opts *bind.CallOpts) ([]byte, error)

	GetFeedData(opts *bind.CallOpts) (*big.Int, error)

	GetThresholdPublicKey(opts *bind.CallOpts) ([]byte, error)

	LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

		error)

	LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

		error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	Transmitters(opts *bind.CallOpts) ([]common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	DeleteCommitment(opts *bind.TransactOpts, requestId [32]byte) (*types.Transaction, error)

	DeleteNodePublicKey(opts *bind.TransactOpts, node common.Address) (*types.Transaction, error)

	OracleWithdraw(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error)

	SendRequest(opts *bind.TransactOpts, request IFunctionsCoordinatorRequest) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, _signers []common.Address, _transmitters []common.Address, _f uint8, _onchainConfig []byte, _offchainConfigVersion uint64, _offchainConfig []byte) (*types.Transaction, error)

	SetDONPublicKey(opts *bind.TransactOpts, donPublicKey []byte) (*types.Transaction, error)

	SetNodePublicKey(opts *bind.TransactOpts, node common.Address, publicKey []byte) (*types.Transaction, error)

	SetThresholdPublicKey(opts *bind.TransactOpts, thresholdPublicKey []byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error)

	UpdateConfig(opts *bind.TransactOpts, config []byte) (*types.Transaction, error)

	FilterBillingEnd(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsCoordinatorBillingEndIterator, error)

	WatchBillingEnd(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorBillingEnd, requestId [][32]byte) (event.Subscription, error)

	ParseBillingEnd(log types.Log) (*FunctionsCoordinatorBillingEnd, error)

	FilterBillingStart(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsCoordinatorBillingStartIterator, error)

	WatchBillingStart(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorBillingStart, requestId [][32]byte) (event.Subscription, error)

	ParseBillingStart(log types.Log) (*FunctionsCoordinatorBillingStart, error)

	FilterCommitmentDeleted(opts *bind.FilterOpts) (*FunctionsCoordinatorCommitmentDeletedIterator, error)

	WatchCommitmentDeleted(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorCommitmentDeleted) (event.Subscription, error)

	ParseCommitmentDeleted(log types.Log) (*FunctionsCoordinatorCommitmentDeleted, error)

	FilterConfigChanged(opts *bind.FilterOpts) (*FunctionsCoordinatorConfigChangedIterator, error)

	WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorConfigChanged) (event.Subscription, error)

	ParseConfigChanged(log types.Log) (*FunctionsCoordinatorConfigChanged, error)

	FilterConfigSet(opts *bind.FilterOpts) (*FunctionsCoordinatorConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*FunctionsCoordinatorConfigSet, error)

	FilterCostExceedsCommitment(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsCoordinatorCostExceedsCommitmentIterator, error)

	WatchCostExceedsCommitment(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorCostExceedsCommitment, requestId [][32]byte) (event.Subscription, error)

	ParseCostExceedsCommitment(log types.Log) (*FunctionsCoordinatorCostExceedsCommitment, error)

	FilterInsufficientGasProvided(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsCoordinatorInsufficientGasProvidedIterator, error)

	WatchInsufficientGasProvided(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorInsufficientGasProvided, requestId [][32]byte) (event.Subscription, error)

	ParseInsufficientGasProvided(log types.Log) (*FunctionsCoordinatorInsufficientGasProvided, error)

	FilterInsufficientSubscriptionBalance(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsCoordinatorInsufficientSubscriptionBalanceIterator, error)

	WatchInsufficientSubscriptionBalance(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorInsufficientSubscriptionBalance, requestId [][32]byte) (event.Subscription, error)

	ParseInsufficientSubscriptionBalance(log types.Log) (*FunctionsCoordinatorInsufficientSubscriptionBalance, error)

	FilterInvalidRequestID(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsCoordinatorInvalidRequestIDIterator, error)

	WatchInvalidRequestID(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorInvalidRequestID, requestId [][32]byte) (event.Subscription, error)

	ParseInvalidRequestID(log types.Log) (*FunctionsCoordinatorInvalidRequestID, error)

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

	FilterTransmitted(opts *bind.FilterOpts) (*FunctionsCoordinatorTransmittedIterator, error)

	WatchTransmitted(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorTransmitted) (event.Subscription, error)

	ParseTransmitted(log types.Log) (*FunctionsCoordinatorTransmitted, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
