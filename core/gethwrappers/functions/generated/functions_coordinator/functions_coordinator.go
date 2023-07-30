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

type FunctionsBillingCommitment struct {
	SubscriptionId          uint64
	Client                  common.Address
	CallbackGasLimit        uint32
	Don                     common.Address
	AdminFee                *big.Int
	EstimatedTotalCostJuels *big.Int
	DonFee                  *big.Int
	Timestamp               uint32
	GasOverhead             *big.Int
	ExpectedGasPrice        *big.Int
}

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

var FunctionsCoordinatorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"linkToNativeFeed\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"EmptyPublicKey\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptyRequestData\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InconsistentReportData\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"InvalidConfig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"linkWei\",\"type\":\"int256\"}],\"name\":\"InvalidLinkWeiPrice\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSubscription\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"MustBeSubOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoTransmittersSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByRouter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByRouterOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentTooLarge\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReportInvalid\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RouterMustBeSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedPublicKeyChange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedSender\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnsupportedRequestDataVersion\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"signerPayment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"transmitterPayment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalCost\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"enumFulfillResult\",\"name\":\"result\",\"type\":\"uint8\"}],\"name\":\"BillingEnd\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"don\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"adminFee\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"estimatedTotalCostJuels\",\"type\":\"uint96\"},{\"internalType\":\"uint80\",\"name\":\"donFee\",\"type\":\"uint80\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint40\",\"name\":\"gasOverhead\",\"type\":\"uint40\"},{\"internalType\":\"uint256\",\"name\":\"expectedGasPrice\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structFunctionsBilling.Commitment\",\"name\":\"commitment\",\"type\":\"tuple\"}],\"name\":\"BillingStart\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maxCallbackGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"feedStalenessSeconds\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasOverheadBeforeCallback\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasOverheadAfterCallback\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"fallbackNativePerUnitLink\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"uint80\",\"name\":\"donFee\",\"type\":\"uint80\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"maxSupportedRequestDataVersion\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fulfillmentGasPriceOverEstimationBP\",\"type\":\"uint256\"}],\"name\":\"ConfigChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"CostExceedsCommitment\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"InsufficientGasProvided\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"InsufficientSubscriptionBalance\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"reason\",\"type\":\"string\"}],\"name\":\"InternalError\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"InternalErrorBytes\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"InvalidRequestID\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requestingContract\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requestInitiator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"subscriptionOwner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"dataVersion\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"flags\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"callbackGasLimit\",\"type\":\"uint64\"}],\"name\":\"OracleRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"OracleResponse\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"RequestTimedOut\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"deleteCommitment\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"node\",\"type\":\"address\"}],\"name\":\"deleteNodePublicKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"}],\"name\":\"estimateCost\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"expectedGasPrice\",\"type\":\"uint256\"}],\"internalType\":\"structIFunctionsBilling.RequestBilling\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"getAdminFee\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllNodePublicKeys\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"},{\"internalType\":\"bytes[]\",\"name\":\"\",\"type\":\"bytes[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfig\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"maxCallbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"feedStalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasOverheadAfterCallback\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"fallbackNativePerUnitLink\",\"type\":\"int256\"},{\"internalType\":\"uint32\",\"name\":\"gasOverheadBeforeCallback\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"linkPriceFeed\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"maxSupportedRequestDataVersion\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"fulfillmentGasPriceOverEstimationBP\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfigHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"config\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"expectedGasPrice\",\"type\":\"uint256\"}],\"internalType\":\"structIFunctionsBilling.RequestBilling\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"getDONFee\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDONPublicKey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFeedData\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getThresholdPublicKey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"oracleWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"requestingContract\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"subscriptionOwner\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"dataVersion\",\"type\":\"uint16\"},{\"internalType\":\"bytes32\",\"name\":\"flags\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"}],\"internalType\":\"structIFunctionsCoordinator.Request\",\"name\":\"request\",\"type\":\"tuple\"}],\"name\":\"sendRequest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint96\",\"name\":\"estimatedCost\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requestTimeoutSeconds\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"_f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"_onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"_offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"_offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"donPublicKey\",\"type\":\"bytes\"}],\"name\":\"setDONPublicKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"node\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"}],\"name\":\"setNodePublicKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"thresholdPublicKey\",\"type\":\"bytes\"}],\"name\":\"setThresholdPublicKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"transmitters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"name\":\"updateConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b5060405162005dd938038062005dd9833981016040819052620000349162000462565b828282828260013380600081620000925760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c557620000c58162000162565b50505015156080526001600160a01b038216620000f557604051632530e88560e11b815260040160405180910390fd5b600880546001600160a01b0319166001600160a01b0384161790556200011b816200020d565b805160209091012060075550600e80546001600160a01b039092166c01000000000000000000000000026001600160601b0390921691909117905550620006349350505050565b336001600160a01b03821603620001bc5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000089565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60008060008060008060008060008980602001905181019062000231919062000572565b9850985098509850985098509850985098506000851362000269576040516321ea67b360e11b81526004810186905260240162000089565b604080516101208101825263ffffffff808c168083528b8216602084018190528a83168486018190528c841660608601819052938a16608086018190526001600160501b038a1660a0870181905261ffff8a1660c0880181905260e088018a90526101009097018d9052600a8054600160f01b9098026001600160f01b03600160a01b909302600160a01b600160f01b0319600160801b90950294909416600160801b600160f01b03196c0100000000000000000000000090990263ffffffff60601b196801000000000000000090970296909616600160401b600160801b03196401000000009098026001600160401b0319909b169098179990991795909516959095179290921794909416949094179290921792909216179055600b829055600c869055517f9136fb3a39fbe82d27ef8048354b8fcc93cf828f72e5c877d8135fde8dc7fd73906200041b908b908b908a908c908b908a908a908a9063ffffffff98891681529688166020880152948716604087015292909516606085015260808401526001600160501b039390931660a083015261ffff9290921660c082015260e08101919091526101000190565b60405180910390a150505050505050505050565b80516001600160a01b03811681146200044757600080fd5b919050565b634e487b7160e01b600052604160045260246000fd5b6000806000606084860312156200047857600080fd5b62000483846200042f565b602085810151919450906001600160401b0380821115620004a357600080fd5b818701915087601f830112620004b857600080fd5b815181811115620004cd57620004cd6200044c565b604051601f8201601f19908116603f01168101908382118183101715620004f857620004f86200044c565b816040528281528a868487010111156200051157600080fd5b600093505b8284101562000535578484018601518185018701529285019262000516565b600086848301015280975050505050505062000554604085016200042f565b90509250925092565b805163ffffffff811681146200044757600080fd5b60008060008060008060008060006101208a8c0312156200059257600080fd5b6200059d8a6200055d565b9850620005ad60208b016200055d565b9750620005bd60408b016200055d565b9650620005cd60608b016200055d565b955060808a01519450620005e460a08b016200055d565b60c08b01519094506001600160501b03811681146200060257600080fd5b60e08b015190935061ffff811681146200061b57600080fd5b809250506101008a015190509295985092959850929598565b60805161578962000650600039600061152101526157896000f3fe608060405234801561001057600080fd5b50600436106101ae5760003560e01c80638cc6acce116100ee578063b1dc65a411610097578063d328a91e11610071578063d328a91e1461053f578063de9dfa4614610547578063e3d0e7121461055a578063f2fde38b1461056d57600080fd5b8063b1dc65a41461040d578063c3f909d414610420578063d227d2451461050f57600080fd5b8063aecc12c5116100c8578063aecc12c5146103a5578063af1296d3146103e5578063afcb95d7146103ed57600080fd5b80638cc6acce146103585780638da5cb5b1461036b5780639883c10d1461039357600080fd5b806379ba50971161015b578063814118341161013557806381411834146102ab57806381f1b938146102c057806381ff7048146102c857806385b214cf1461033557600080fd5b806379ba50971461027d5780637f15e16614610285578063807560311461029857600080fd5b8063466557181161018c5780634665571814610226578063533989871461025457806366316d8d1461026a57600080fd5b8063083a5466146101b3578063181f5a77146101c857806326ceabac14610213575b600080fd5b6101c66101c136600461412e565b610580565b005b60408051808201909152601881527f46756e6374696f6e7320436f6f7264696e61746f72207631000000000000000060208201525b60405161020a91906141d4565b60405180910390f35b6101c6610221366004614209565b6105d5565b610239610234366004614365565b610677565b60405169ffffffffffffffffffff909116815260200161020a565b61025c6106a5565b60405161020a92919061448c565b6101c6610278366004614536565b610929565b6101c6610aa7565b6101c661029336600461412e565b610ba9565b6101c66102a636600461456f565b610bf9565b6102b3610cb0565b60405161020a91906145c4565b6101fd610d1f565b61031260015460025463ffffffff74010000000000000000000000000000000000000000830481169378010000000000000000000000000000000000000000000000009093041691565b6040805163ffffffff94851681529390921660208401529082015260600161020a565b6103486103433660046145d7565b610df0565b604051901515815260200161020a565b6101c66103663660046145f0565b611012565b60005460405173ffffffffffffffffffffffffffffffffffffffff909116815260200161020a565b6007545b60405190815260200161020a565b6103b86103b336600461462d565b611079565b604080519485526bffffffffffffffffffffffff909316602085015291830152606082015260800161020a565b6103976112bb565b60408051600181526000602082018190529181019190915260600161020a565b6101c661041b3660046146ad565b6113ae565b6104ab600a54600c54600e54600b5463ffffffff8085169564010000000086048216956c010000000000000000000000008082048416969568010000000000000000830490941694930473ffffffffffffffffffffffffffffffffffffffff16927e0100000000000000000000000000000000000000000000000000000000000090910461ffff1691565b6040805163ffffffff998a168152978916602089015287019590955260608601939093529416608084015273ffffffffffffffffffffffffffffffffffffffff90931660a083015261ffff90921660c082015260e08101919091526101000161020a565b61052261051d366004614764565b611adf565b6040516bffffffffffffffffffffffff909116815260200161020a565b6101fd611cb1565b610522610555366004614365565b611d08565b6101c661056836600461488c565b611da3565b6101c661057b366004614209565b6127c7565b6105886127d8565b60008190036105c3576040517f4f42be3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60116105d08284836149fa565b505050565b60005473ffffffffffffffffffffffffffffffffffffffff1633148061061057503373ffffffffffffffffffffffffffffffffffffffff8216145b610646576040517fed6dd19b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8116600090815260106020526040812061067491614073565b50565b600a5474010000000000000000000000000000000000000000900469ffffffffffffffffffff165b92915050565b6060806000600680548060200260200160405190810160405280929190818152602001828054801561070d57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116106e2575b505050505090506000815167ffffffffffffffff81111561073057610730614226565b60405190808252806020026020018201604052801561076357816020015b606081526020019060019003908161074e5790505b50905060005b825181101561091f576010600084838151811061078857610788614b15565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002080546107d590614959565b9050600003610810576040517f4f42be3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6010600084838151811061082657610826614b15565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020805461087390614959565b80601f016020809104026020016040519081016040528092919081815260200182805461089f90614959565b80156108ec5780601f106108c1576101008083540402835291602001916108ec565b820191906000526020600020905b8154815290600101906020018083116108cf57829003601f168201915b505050505082828151811061090357610903614b15565b60200260200101819052508061091890614b73565b9050610769565b5090939092509050565b61093161285b565b806bffffffffffffffffffffffff166000036109675750336000908152600d60205260409020546bffffffffffffffffffffffff165b336000908152600d60205260409020546bffffffffffffffffffffffff808316911610156109c1576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b336000908152600d6020526040812080548392906109ee9084906bffffffffffffffffffffffff16614bab565b82546101009290920a6bffffffffffffffffffffffff8181021990931691831602179091556008546040517f66316d8d00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff868116600483015292851660248201529116915081906366316d8d90604401600060405180830381600087803b158015610a8a57600080fd5b505af1158015610a9e573d6000803e3d6000fd5b50505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610b2d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610bb16127d8565b6000819003610bec576040517f4f42be3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600f6105d08284836149fa565b60005473ffffffffffffffffffffffffffffffffffffffff16331480610c445750610c23336129f2565b8015610c4457503373ffffffffffffffffffffffffffffffffffffffff8416145b610c7a576040517fed6dd19b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff83166000908152601060205260409020610caa8284836149fa565b50505050565b60606006805480602002602001604051908101604052809291908181526020018280548015610d1557602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610cea575b5050505050905090565b606060118054610d2e90614959565b9050600003610d69576040517f4f42be3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60118054610d7690614959565b80601f0160208091040260200160405190810160405280929190818152602001828054610da290614959565b8015610d155780601f10610dc457610100808354040283529160200191610d15565b820191906000526020600020905b815481529060010190602001808311610dd257509395945050505050565b60085460009073ffffffffffffffffffffffffffffffffffffffff163314610e44576040517fc41a5b0900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260096020908152604091829020825161014081018452815467ffffffffffffffff8116825268010000000000000000810473ffffffffffffffffffffffffffffffffffffffff908116948301949094527c0100000000000000000000000000000000000000000000000000000000900463ffffffff90811694820194909452600182015492831660608201819052740100000000000000000000000000000000000000009093046bffffffffffffffffffffffff9081166080830152600283015490811660a08301526c01000000000000000000000000810469ffffffffffffffffffff1660c0830152760100000000000000000000000000000000000000000000810490941660e08201527a01000000000000000000000000000000000000000000000000000090930464ffffffffff1661010084015260030154610120830152610f985750600092915050565b600083815260096020526040808220828155600181018390556002810180547fff000000000000000000000000000000000000000000000000000000000000001690556003018290555184917ff1ca1e9147be737b04a2b018a79405f687a97de8dd8a2559bbe62357343af41491a260019150505b919050565b60085473ffffffffffffffffffffffffffffffffffffffff163314611063576040517fc41a5b0900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61106c81612adb565b8051602090910120600755565b60085460009081908190819073ffffffffffffffffffffffffffffffffffffffff1633146110d3576040517fc41a5b0900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6110e06040860186614bd7565b905060000361111a576040517ec1cfc000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000604051806080016040528087606001602081019061113a9190614c3c565b67ffffffffffffffff16815260209081019061115890890189614209565b73ffffffffffffffffffffffffffffffffffffffff16815260200161118360e0890160c08a01614c59565b63ffffffff1681523a60209091015290506111ef6111a46040880188614bd7565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506111e99250505060a0890160808a01614c86565b83612dd6565b929750909550935091506112066020870187614209565b73ffffffffffffffffffffffffffffffffffffffff16857fd6bf2929a5458fc22a90821275e4f54e2566c579e4bb3d38184880bf1701762c3261124f60808b0160608c01614c3c565b61125f60408c0160208d01614209565b61126c60408d018d614bd7565b8d608001602081019061127f9190614c86565b8e60a001358f60c00160208101906112979190614c59565b6040516112ab989796959493929190614ca3565b60405180910390a3509193509193565b600a54600e54604080517ffeaf968c0000000000000000000000000000000000000000000000000000000081529051600093640100000000900463ffffffff169283151592859283926c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff169163feaf968c9160048083019260a09291908290030181865afa158015611355573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906113799190614d63565b5093505092505082801561139b57506113928142614db3565b8463ffffffff16105b156113a657600c5491505b509392505050565b60005a604080516020601f8b018190048102820181019092528981529192508a3591818c01359161140491849163ffffffff851691908e908e908190840183828082843760009201919091525061346392505050565b61143a576040517f0be3632800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805183815262ffffff600884901c1660208201527fb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62910160405180910390a16040805160608101825260025480825260035460ff8082166020850152610100909104169282019290925290831461150f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f636f6e666967446967657374206d69736d6174636800000000000000000000006044820152606401610b24565b61151d8b8b8b8b8b8b61346c565b60007f00000000000000000000000000000000000000000000000000000000000000001561157a5760028260200151836040015161155b9190614dc6565b6115659190614e0e565b611570906001614dc6565b60ff169050611590565b602082015161158a906001614dc6565b60ff1690505b8881146115f9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f77726f6e67206e756d626572206f66207369676e6174757265730000000000006044820152606401610b24565b888714611662576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f7369676e617475726573206f7574206f6620726567697374726174696f6e00006044820152606401610b24565b3360009081526004602090815260408083208151808301909252805460ff808216845292939192918401916101009091041660028111156116a5576116a5614e30565b60028111156116b6576116b6614e30565b90525090506002816020015160028111156116d3576116d3614e30565b14801561171a57506006816000015160ff16815481106116f5576116f5614b15565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff1633145b611780576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f756e617574686f72697a6564207472616e736d697474657200000000000000006044820152606401610b24565b505050505061178d6140ad565b6000808a8a6040516117a0929190614e5f565b6040519081900381206117b7918e90602001614e6f565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181528282528051602091820120838301909252600080845290830152915060005b89811015611ac157600060018489846020811061182057611820614b15565b61182d91901a601b614dc6565b8e8e8681811061183f5761183f614b15565b905060200201358d8d8781811061185857611858614b15565b9050602002013560405160008152602001604052604051611895949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa1580156118b7573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff811660009081526004602090815290849020838501909452835460ff8082168552929650929450840191610100900416600281111561193757611937614e30565b600281111561194857611948614e30565b905250925060018360200151600281111561196557611965614e30565b146119cc576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f61646472657373206e6f7420617574686f72697a656420746f207369676e00006044820152606401610b24565b8251600090879060ff16601f81106119e6576119e6614b15565b602002015173ffffffffffffffffffffffffffffffffffffffff1614611a68576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f6e6f6e2d756e69717565207369676e61747572650000000000000000000000006044820152606401610b24565b8086846000015160ff16601f8110611a8257611a82614b15565b73ffffffffffffffffffffffffffffffffffffffff9092166020929092020152611aad600186614dc6565b94505080611aba90614b73565b9050611801565b505050611ad2833383858e8e61351a565b5050505050505050505050565b6008546040517f10fc49c100000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8716600482015263ffffffff8416602482015260009173ffffffffffffffffffffffffffffffffffffffff169081906310fc49c19060440160006040518083038186803b158015611b6157600080fd5b505afa158015611b75573d6000803e3d6000fd5b50505050620f4240831115611bb6576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600060405180608001604052808967ffffffffffffffff1681526020013373ffffffffffffffffffffffffffffffffffffffff1681526020018663ffffffff1681526020018581525090506000611c4488888080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250869250610677915050565b69ffffffffffffffffffff1690506000611c9589898080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250879250611d08915050565b9050611ca3878784846137e6565b9a9950505050505050505050565b6060600f8054611cc090614959565b9050600003611cfb576040517f4f42be3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600f8054610d7690614959565b600854604080517f2a905ccc000000000000000000000000000000000000000000000000000000008152905160009273ffffffffffffffffffffffffffffffffffffffff1691632a905ccc9160048083019260209291908290030181865afa158015611d78573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611d9c9190614e83565b9392505050565b855185518560ff16601f831115611e16576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f746f6f206d616e79207369676e657273000000000000000000000000000000006044820152606401610b24565b80600003611e80576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f66206d75737420626520706f73697469766500000000000000000000000000006044820152606401610b24565b818314611f0e576040517f89a61989000000000000000000000000000000000000000000000000000000008152602060048201526024808201527f6f7261636c6520616464726573736573206f7574206f6620726567697374726160448201527f74696f6e000000000000000000000000000000000000000000000000000000006064820152608401610b24565b611f19816003614ea0565b8311611f81576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f6661756c74792d6f7261636c65206620746f6f206869676800000000000000006044820152606401610b24565b611f896127d8565b6040805160c0810182528a8152602081018a905260ff89169181018290526060810188905267ffffffffffffffff8716608082015260a0810186905290611fd09088613948565b6005541561218557600554600090611fea90600190614db3565b905060006005828154811061200157612001614b15565b60009182526020822001546006805473ffffffffffffffffffffffffffffffffffffffff9092169350908490811061203b5761203b614b15565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff85811684526004909252604080842080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000908116909155929091168084529220805490911690556005805491925090806120bb576120bb614eb7565b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff0000000000000000000000000000000000000000169055019055600680548061212457612124614eb7565b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff000000000000000000000000000000000000000016905501905550611fd0915050565b60005b8151518110156125ec57600060046000846000015184815181106121ae576121ae614b15565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff16825281019190915260400160002054610100900460ff1660028111156121f8576121f8614e30565b1461225f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265706561746564207369676e657220616464726573730000000000000000006044820152606401610b24565b6040805180820190915260ff8216815260016020820152825180516004916000918590811061229057612290614b15565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040016000208251815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082168117835592840151919283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000161761010083600281111561233157612331614e30565b0217905550600091506123419050565b600460008460200151848151811061235b5761235b614b15565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff16825281019190915260400160002054610100900460ff1660028111156123a5576123a5614e30565b1461240c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f7265706561746564207472616e736d69747465722061646472657373000000006044820152606401610b24565b6040805180820190915260ff82168152602081016002815250600460008460200151848151811061243f5761243f614b15565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040016000208251815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082168117835592840151919283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000016176101008360028111156124e0576124e0614e30565b0217905550508251805160059250839081106124fe576124fe614b15565b602090810291909101810151825460018101845560009384529282902090920180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909316929092179091558201518051600691908390811061257a5761257a614b15565b60209081029190910181015182546001810184556000938452919092200180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909216919091179055806125e481614b73565b915050612188565b506040810151600380547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff909216919091179055600180547fffffffff00000000ffffffffffffffffffffffffffffffffffffffffffffffff8116780100000000000000000000000000000000000000000000000063ffffffff43811682029290921780855592048116929182916014916126a491849174010000000000000000000000000000000000000000900416614ee6565b92506101000a81548163ffffffff021916908363ffffffff1602179055506127034630600160149054906101000a900463ffffffff1663ffffffff16856000015186602001518760400151886060015189608001518a60a00151613965565b600281905582518051600380547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1661010060ff9093169290920291909117905560015460208501516040808701516060880151608089015160a08a015193517f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05986127ba988b9891977401000000000000000000000000000000000000000090920463ffffffff16969095919491939192614f03565b60405180910390a1611ad2565b6127cf6127d8565b61067481613a10565b60005473ffffffffffffffffffffffffffffffffffffffff163314612859576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610b24565b565b6000612865610cb0565b905080516000036128a2576040517f30274b3a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8051600e546000916128c1916bffffffffffffffffffffffff16614f99565b905060005b82518160ff1610156129935781600d6000858460ff16815181106128ec576128ec614b15565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282829054906101000a90046bffffffffffffffffffffffff166129549190614fc4565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055508061298c90614fe9565b90506128c6565b5081516129a09082615008565b600e80546000906129c09084906bffffffffffffffffffffffff16614bab565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055505050565b6000806006805480602002602001604051908101604052809291908181526020018280548015612a5857602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311612a2d575b5050505050905060005b8151811015612ad1578373ffffffffffffffffffffffffffffffffffffffff16828281518110612a9457612a94614b15565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1603612ac1575060019392505050565b612aca81614b73565b9050612a62565b5060009392505050565b600080600080600080600080600089806020019051810190612afd9190615038565b98509850985098509850985098509850985060008513612b4c576040517f43d4cf6600000000000000000000000000000000000000000000000000000000815260048101869052602401610b24565b604080516101208101825263ffffffff808c168083528b8216602084018190528a83168486018190528c841660608601819052938a166080860181905269ffffffffffffffffffff8a1660a0870181905261ffff8a1660c0880181905260e088018a90526101009097018d9052600a80547e010000000000000000000000000000000000000000000000000000000000009098027dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff740100000000000000000000000000000000000000009093027fffff00000000000000000000ffffffffffffffffffffffffffffffffffffffff700100000000000000000000000000000000909502949094167fffff0000000000000000000000000000ffffffffffffffffffffffffffffffff6c010000000000000000000000009099027fffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffffff68010000000000000000909702969096167fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff6401000000009098027fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000909b169098179990991795909516959095179290921794909416949094179290921792909216179055600b829055600c869055517f9136fb3a39fbe82d27ef8048354b8fcc93cf828f72e5c877d8135fde8dc7fd7390612dc2908b908b908a908c908b908a908a908a9063ffffffff988916815296881660208801529487166040870152929095166060850152608084015269ffffffffffffffffffff9390931660a083015261ffff9290921660c082015260e08101919091526101000190565b60405180910390a150505050505050505050565b600a5460009081908190819061ffff7e0100000000000000000000000000000000000000000000000000000000000090910481169087161115612e45576040517fdada758700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000612e518887610677565b90506000612e5f8988611d08565b90506000612e83886040015189606001518569ffffffffffffffffffff16856137e6565b60085489516040517fa47c769600000000000000000000000000000000000000000000000000000000815267ffffffffffffffff909116600482015291925073ffffffffffffffffffffffffffffffffffffffff16906000908190839063a47c769690602401600060405180830381865afa158015612f06573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052612f4c91908101906150e4565b50505060208d01518d516040517f674603d000000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff928316600482015267ffffffffffffffff90911660248201529294509092506000919085169063674603d090604401606060405180830381865afa158015612fdd573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061300191906151db565b509150506bffffffffffffffffffffffff851661301e8385614bab565b6bffffffffffffffffffffffff161015613064576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006130e3308e602001518f60000151856001613081919061522d565b6040805173ffffffffffffffffffffffffffffffffffffffff958616602080830191909152949095168582015267ffffffffffffffff928316606086015291166080808501919091528151808503909101815260a09093019052815191012090565b905060006040518061014001604052808f6000015167ffffffffffffffff1681526020018f6020015173ffffffffffffffffffffffffffffffffffffffff1681526020018f6040015163ffffffff1681526020013073ffffffffffffffffffffffffffffffffffffffff168152602001896bffffffffffffffffffffffff168152602001886bffffffffffffffffffffffff1681526020018a69ffffffffffffffffffff1681526020014263ffffffff168152602001600a600001600c9054906101000a900463ffffffff16600a60000160089054906101000a900463ffffffff166131cf9190614ee6565b63ffffffff1664ffffffffff1681526020018f606001518152509050806009600084815260200190815260200160002060008201518160000160006101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555060208201518160000160086101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550604082015181600001601c6101000a81548163ffffffff021916908363ffffffff16021790555060608201518160010160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060808201518160010160146101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060a08201518160020160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060c082015181600201600c6101000a81548169ffffffffffffffffffff021916908369ffffffffffffffffffff16021790555060e08201518160020160166101000a81548163ffffffff021916908363ffffffff16021790555061010082015181600201601a6101000a81548164ffffffffff021916908364ffffffffff1602179055506101208201518160030155905050817fe2cc993cebc49ce369b0174671d90dbae2201decdae368805de999b437336dcd82604051613413919061524e565b60405180910390a250600a54909f959e5063ffffffff6c01000000000000000000000000820481169e50700100000000000000000000000000000000909104169b50939950505050505050505050565b60019392505050565b6000613479826020614ea0565b613484856020614ea0565b61349088610144615363565b61349a9190615363565b6134a49190615363565b6134af906000615363565b9050368114610a9e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f63616c6c64617461206c656e677468206d69736d6174636800000000000000006044820152606401610b24565b6060808061352a84860186615401565b82519295509093509150158061354257508151835114155b8061354f57508051835114155b15613586576040517f0be3632800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b83518110156137da5760006135ea8583815181106135a9576135a9614b15565b60200260200101518584815181106135c3576135c3614b15565b60200260200101518585815181106135dd576135dd614b15565b6020026020010151613b05565b9050600281600681111561360057613600614e30565b036136505784828151811061361757613617614b15565b60200260200101517fa1c120e327c9ad8b075793878c88d59b8934b97ae37117faa3bb21616237f7be60405160405180910390a2613778565b600381600681111561366457613664614e30565b036136b45784828151811061367b5761367b614b15565b60200260200101517f357a60574f143e144d62ddc89636fee45ba5e2f8c9a30b9f91e28f8dfe5a56e060405160405180910390a2613778565b60058160068111156136c8576136c8614e30565b03613718578482815181106136df576136df614b15565b60200260200101517f689f49eacf6db85082542f4b7873f4a519744a4e093fd002310545bfb4354cc860405160405180910390a2613778565b600481600681111561372c5761372c614e30565b036137785784828151811061374357613743614b15565b60200260200101517fb83d9204606ee3780daa106ad026c4c2853f5af126edbb8fe15e813e62dbd94d60405160405180910390a25b84828151811061378a5761378a614b15565b60209081029190910181015160405133815290917fc708e0440951fd63499c0f7a73819b469ee5dd3ecc356c0ab4eb7f18389009d9910160405180910390a2506137d381614b73565b9050613589565b50505050505050505050565b6000806137f16112bb565b905060008113613830576040517f43d4cf6600000000000000000000000000000000000000000000000000000000815260048101829052602401610b24565b600a5460009087906138669063ffffffff6c01000000000000000000000000820481169168010000000000000000900416614ee6565b6138709190614ee6565b63ffffffff1690506000612710600a600101548861388e9190614ea0565b61389891906154e9565b6138a29088615363565b9050600083836138ba84670de0b6b3a7640000614ea0565b6138c49190614ea0565b6138ce91906154e9565b905060006138ed6bffffffffffffffffffffffff808916908a16615363565b9050613905816b033b2e3c9fd0803ce8000000614db3565b82111561393e576040517fe80fa38100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611ca38183615363565b6000613952610cb0565b5111156139615761396161285b565b5050565b6000808a8a8a8a8a8a8a8a8a604051602001613989999897969594939291906154fd565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179150509998505050505050505050565b3373ffffffffffffffffffffffffffffffffffffffff821603613a8f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610b24565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000838152600960209081526040808320815161014081018352815467ffffffffffffffff8116825273ffffffffffffffffffffffffffffffffffffffff68010000000000000000820481169583019590955263ffffffff7c01000000000000000000000000000000000000000000000000000000009091048116938201939093526001820154938416606082018190526bffffffffffffffffffffffff7401000000000000000000000000000000000000000090950485166080830152600283015494851660a083015269ffffffffffffffffffff6c0100000000000000000000000086041660c0830152760100000000000000000000000000000000000000000000850490931660e082015264ffffffffff7a0100000000000000000000000000000000000000000000000000009094049390931661010084015260030154610120830152613c5a576002915050611d9c565b6000858152600960205260408120818155600181018290556002810180547fff00000000000000000000000000000000000000000000000000000000000000169055600301819055613caa6112bb565b905060008113613ce9576040517f43d4cf6600000000000000000000000000000000000000000000000000000000815260048101829052602401610b24565b600081613cfe3a670de0b6b3a7640000614ea0565b613d0891906154e9565b9050600083610100015164ffffffffff1682613d249190614ea0565b905060008460c0015169ffffffffffffffffffff1682613d449190614fc4565b6008546040517f3fd67e0500000000000000000000000000000000000000000000000000000000815291925073ffffffffffffffffffffffffffffffffffffffff16908190633fd67e0590613da7908d908d908d908a9089903390600401615592565b60408051808303816000875af1925050508015613dff575060408051601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201909252613dfc91810190615603565b60015b613ed457613e0b615632565b806308c379a003613e725750613e1f61564e565b80613e2a5750613e74565b8a7ffb78cb471c8bc73de4ca3325ff527cfac9462f08cf259c153f25b3a51afabd7982604051613e5a91906141d4565b60405180910390a26006975050505050505050611d9c565b505b3d808015613e9e576040519150601f19603f3d011682016040523d82523d6000602084013e613ea3565b606091505b508a7f0efeb0e5bf272b1b81d86c1dc2b6e73cbcd571830561eb4646d8195ed5cf6df182604051613e5a91906141d4565b60008260ff166006811115613eeb57613eeb614e30565b90506000816006811115613f0157613f01614e30565b1480613f1e57506001816006811115613f1c57613f1c614e30565b145b1561406357613f2d8287614fc4565b336000908152600d6020526040812080547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff93841617905560c08b0151600e805469ffffffffffffffffffff90921693909291613f9a91859116614fc4565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055508c7f30e1e299f390e4f3fa75e4489fed477e86183937a4a21a0ae7b137a4dff035358a600001518b60c00151858a6140009190614fc4565b60808e015160c08f015169ffffffffffffffffffff16614020898e614fc4565b61402a9190614fc4565b6140349190614fc4565b8860ff16600681111561404957614049614e30565b60405161405a9594939291906156f6565b60405180910390a25b9850611d9c975050505050505050565b50805461407f90614959565b6000825580601f1061408f575050565b601f01602090049060005260206000209081019061067491906140cc565b604051806103e00160405280601f906020820280368337509192915050565b5b808211156140e157600081556001016140cd565b5090565b60008083601f8401126140f757600080fd5b50813567ffffffffffffffff81111561410f57600080fd5b60208301915083602082850101111561412757600080fd5b9250929050565b6000806020838503121561414157600080fd5b823567ffffffffffffffff81111561415857600080fd5b614164858286016140e5565b90969095509350505050565b6000815180845260005b818110156141965760208185018101518683018201520161417a565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b602081526000611d9c6020830184614170565b73ffffffffffffffffffffffffffffffffffffffff8116811461067457600080fd5b60006020828403121561421b57600080fd5b8135611d9c816141e7565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f830116810181811067ffffffffffffffff8211171561429957614299614226565b6040525050565b600082601f8301126142b157600080fd5b813567ffffffffffffffff8111156142cb576142cb614226565b60405161430060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8501160182614255565b81815284602083860101111561431557600080fd5b816020850160208301376000918101602001919091529392505050565b67ffffffffffffffff8116811461067457600080fd5b803561100d81614332565b63ffffffff8116811461067457600080fd5b60008082840360a081121561437957600080fd5b833567ffffffffffffffff8082111561439157600080fd5b61439d878388016142a0565b945060807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0840112156143cf57600080fd5b604051925060808301915082821081831117156143ee576143ee614226565b50604052602084013561440081614332565b81526040840135614410816141e7565b6020820152606084013561442381614353565b60408201526080939093013560608401525092909150565b600081518084526020808501945080840160005b8381101561448157815173ffffffffffffffffffffffffffffffffffffffff168752958201959082019060010161444f565b509495945050505050565b60408152600061449f604083018561443b565b6020838203818501528185518084528284019150828160051b85010183880160005b8381101561450d577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08784030185526144fb838351614170565b948601949250908501906001016144c1565b50909998505050505050505050565b6bffffffffffffffffffffffff8116811461067457600080fd5b6000806040838503121561454957600080fd5b8235614554816141e7565b915060208301356145648161451c565b809150509250929050565b60008060006040848603121561458457600080fd5b833561458f816141e7565b9250602084013567ffffffffffffffff8111156145ab57600080fd5b6145b7868287016140e5565b9497909650939450505050565b602081526000611d9c602083018461443b565b6000602082840312156145e957600080fd5b5035919050565b60006020828403121561460257600080fd5b813567ffffffffffffffff81111561461957600080fd5b614625848285016142a0565b949350505050565b60006020828403121561463f57600080fd5b813567ffffffffffffffff81111561465657600080fd5b820160e08185031215611d9c57600080fd5b60008083601f84011261467a57600080fd5b50813567ffffffffffffffff81111561469257600080fd5b6020830191508360208260051b850101111561412757600080fd5b60008060008060008060008060e0898b0312156146c957600080fd5b606089018a8111156146da57600080fd5b8998503567ffffffffffffffff808211156146f457600080fd5b6147008c838d016140e5565b909950975060808b013591508082111561471957600080fd5b6147258c838d01614668565b909750955060a08b013591508082111561473e57600080fd5b5061474b8b828c01614668565b999c989b50969995989497949560c00135949350505050565b60008060008060006080868803121561477c57600080fd5b853561478781614332565b9450602086013567ffffffffffffffff8111156147a357600080fd5b6147af888289016140e5565b90955093505060408601356147c381614353565b949793965091946060013592915050565b600067ffffffffffffffff8211156147ee576147ee614226565b5060051b60200190565b600082601f83011261480957600080fd5b81356020614816826147d4565b6040516148238282614255565b83815260059390931b850182019282810191508684111561484357600080fd5b8286015b8481101561486757803561485a816141e7565b8352918301918301614847565b509695505050505050565b60ff8116811461067457600080fd5b803561100d81614872565b60008060008060008060c087890312156148a557600080fd5b863567ffffffffffffffff808211156148bd57600080fd5b6148c98a838b016147f8565b975060208901359150808211156148df57600080fd5b6148eb8a838b016147f8565b96506148f960408a01614881565b9550606089013591508082111561490f57600080fd5b61491b8a838b016142a0565b945061492960808a01614348565b935060a089013591508082111561493f57600080fd5b5061494c89828a016142a0565b9150509295509295509295565b600181811c9082168061496d57607f821691505b6020821081036149a6577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f8211156105d057600081815260208120601f850160051c810160208610156149d35750805b601f850160051c820191505b818110156149f2578281556001016149df565b505050505050565b67ffffffffffffffff831115614a1257614a12614226565b614a2683614a208354614959565b836149ac565b6000601f841160018114614a785760008515614a425750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b178355614b0e565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b82811015614ac75786850135825560209485019460019092019101614aa7565b5086821015614b02577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203614ba457614ba4614b44565b5060010190565b6bffffffffffffffffffffffff828116828216039080821115614bd057614bd0614b44565b5092915050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112614c0c57600080fd5b83018035915067ffffffffffffffff821115614c2757600080fd5b60200191503681900382131561412757600080fd5b600060208284031215614c4e57600080fd5b8135611d9c81614332565b600060208284031215614c6b57600080fd5b8135611d9c81614353565b61ffff8116811461067457600080fd5b600060208284031215614c9857600080fd5b8135611d9c81614c76565b600073ffffffffffffffffffffffffffffffffffffffff808b16835267ffffffffffffffff8a16602084015280891660408401525060e060608301528560e0830152610100868882850137600083880182015261ffff9590951660808301525060a081019290925263ffffffff1660c0820152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01690910101949350505050565b805169ffffffffffffffffffff8116811461100d57600080fd5b600080600080600060a08688031215614d7b57600080fd5b614d8486614d49565b9450602086015193506040860151925060608601519150614da760808701614d49565b90509295509295909350565b8181038181111561069f5761069f614b44565b60ff818116838216019081111561069f5761069f614b44565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b600060ff831680614e2157614e21614ddf565b8060ff84160491505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b8183823760009101908152919050565b828152606082602083013760800192915050565b600060208284031215614e9557600080fd5b8151611d9c8161451c565b808202811582820484141761069f5761069f614b44565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b63ffffffff818116838216019080821115614bd057614bd0614b44565b600061012063ffffffff808d1684528b6020850152808b16604085015250806060840152614f338184018a61443b565b90508281036080840152614f47818961443b565b905060ff871660a084015282810360c0840152614f648187614170565b905067ffffffffffffffff851660e0840152828103610100840152614f898185614170565b9c9b505050505050505050505050565b60006bffffffffffffffffffffffff80841680614fb857614fb8614ddf565b92169190910492915050565b6bffffffffffffffffffffffff818116838216019080821115614bd057614bd0614b44565b600060ff821660ff8103614fff57614fff614b44565b60010192915050565b6bffffffffffffffffffffffff81811683821602808216919082811461503057615030614b44565b505092915050565b60008060008060008060008060006101208a8c03121561505757600080fd5b895161506281614353565b60208b015190995061507381614353565b60408b015190985061508481614353565b60608b015190975061509581614353565b60808b015160a08c015191975095506150ad81614353565b93506150bb60c08b01614d49565b925060e08a01516150cb81614c76565b809250506101008a015190509295985092959850929598565b600080600080600060a086880312156150fc57600080fd5b85516151078161451c565b8095505060208087015161511a8161451c565b604088015190955061512b816141e7565b606088015190945061513c816141e7565b608088015190935067ffffffffffffffff81111561515957600080fd5b8701601f8101891361516a57600080fd5b8051615175816147d4565b6040516151828282614255565b82815260059290921b830184019184810191508b8311156151a257600080fd5b928401925b828410156151c95783516151ba816141e7565b825292840192908401906151a7565b80955050505050509295509295909350565b6000806000606084860312156151f057600080fd5b8351801515811461520057600080fd5b602085015190935061521181614332565b604085015190925061522281614332565b809150509250925092565b67ffffffffffffffff818116838216019080821115614bd057614bd0614b44565b815167ffffffffffffffff16815261014081016020830151615288602084018273ffffffffffffffffffffffffffffffffffffffff169052565b5060408301516152a0604084018263ffffffff169052565b5060608301516152c8606084018273ffffffffffffffffffffffffffffffffffffffff169052565b5060808301516152e860808401826bffffffffffffffffffffffff169052565b5060a083015161530860a08401826bffffffffffffffffffffffff169052565b5060c083015161532660c084018269ffffffffffffffffffff169052565b5060e083015161533e60e084018263ffffffff169052565b506101008381015164ffffffffff811684830152505061012092830151919092015290565b8082018082111561069f5761069f614b44565b600082601f83011261538757600080fd5b81356020615394826147d4565b6040516153a18282614255565b83815260059390931b85018201928281019150868411156153c157600080fd5b8286015b8481101561486757803567ffffffffffffffff8111156153e55760008081fd5b6153f38986838b01016142a0565b8452509183019183016153c5565b60008060006060848603121561541657600080fd5b833567ffffffffffffffff8082111561542e57600080fd5b818601915086601f83011261544257600080fd5b8135602061544f826147d4565b60405161545c8282614255565b83815260059390931b850182019282810191508a84111561547c57600080fd5b948201945b8386101561549a57853582529482019490820190615481565b975050870135925050808211156154b057600080fd5b6154bc87838801615376565b935060408601359150808211156154d257600080fd5b506154df86828701615376565b9150509250925092565b6000826154f8576154f8614ddf565b500490565b60006101208b835273ffffffffffffffffffffffffffffffffffffffff8b16602084015267ffffffffffffffff808b1660408501528160608501526155448285018b61443b565b91508382036080850152615558828a61443b565b915060ff881660a085015283820360c08501526155758288614170565b90861660e08501528381036101008501529050614f898185614170565b86815260c0602082015260006155ab60c0830188614170565b82810360408401526155bd8188614170565b6bffffffffffffffffffffffff96871660608501529490951660808301525073ffffffffffffffffffffffffffffffffffffffff9190911660a090910152949350505050565b6000806040838503121561561657600080fd5b825161562181614872565b60208401519092506145648161451c565b600060033d111561564b5760046000803e5060005160e01c5b90565b600060443d101561565c5790565b6040517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc803d016004833e81513d67ffffffffffffffff81602484011181841117156156aa57505050505090565b82850191508151818111156156c25750505050505090565b843d87010160208285010111156156dc5750505050505090565b6156eb60208286010187614255565b509095945050505050565b67ffffffffffffffff8616815269ffffffffffffffffffff851660208201526bffffffffffffffffffffffff84811660408301528316606082015260a081016007831061576c577f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b826080830152969550505050505056fea164736f6c6343000813000a",
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
	outstruct.GasOverheadAfterCallback = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.FallbackNativePerUnitLink = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.GasOverheadBeforeCallback = *abi.ConvertType(out[4], new(uint32)).(*uint32)
	outstruct.LinkPriceFeed = *abi.ConvertType(out[5], new(common.Address)).(*common.Address)
	outstruct.MaxSupportedRequestDataVersion = *abi.ConvertType(out[6], new(uint16)).(*uint16)
	outstruct.FulfillmentGasPriceOverEstimationBP = *abi.ConvertType(out[7], new(*big.Int)).(**big.Int)

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
	Commitment FunctionsBillingCommitment
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
	FallbackNativePerUnitLink           *big.Int
	DonFee                              *big.Int
	MaxSupportedRequestDataVersion      uint16
	FulfillmentGasPriceOverEstimationBP *big.Int
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

type FunctionsCoordinatorInternalErrorIterator struct {
	Event *FunctionsCoordinatorInternalError

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsCoordinatorInternalErrorIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsCoordinatorInternalError)
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
		it.Event = new(FunctionsCoordinatorInternalError)
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

func (it *FunctionsCoordinatorInternalErrorIterator) Error() error {
	return it.fail
}

func (it *FunctionsCoordinatorInternalErrorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsCoordinatorInternalError struct {
	RequestId [32]byte
	Reason    string
	Raw       types.Log
}

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) FilterInternalError(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsCoordinatorInternalErrorIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsCoordinator.contract.FilterLogs(opts, "InternalError", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsCoordinatorInternalErrorIterator{contract: _FunctionsCoordinator.contract, event: "InternalError", logs: logs, sub: sub}, nil
}

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) WatchInternalError(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorInternalError, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsCoordinator.contract.WatchLogs(opts, "InternalError", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsCoordinatorInternalError)
				if err := _FunctionsCoordinator.contract.UnpackLog(event, "InternalError", log); err != nil {
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

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) ParseInternalError(log types.Log) (*FunctionsCoordinatorInternalError, error) {
	event := new(FunctionsCoordinatorInternalError)
	if err := _FunctionsCoordinator.contract.UnpackLog(event, "InternalError", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsCoordinatorInternalErrorBytesIterator struct {
	Event *FunctionsCoordinatorInternalErrorBytes

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsCoordinatorInternalErrorBytesIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsCoordinatorInternalErrorBytes)
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
		it.Event = new(FunctionsCoordinatorInternalErrorBytes)
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

func (it *FunctionsCoordinatorInternalErrorBytesIterator) Error() error {
	return it.fail
}

func (it *FunctionsCoordinatorInternalErrorBytesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsCoordinatorInternalErrorBytes struct {
	RequestId [32]byte
	Reason    []byte
	Raw       types.Log
}

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) FilterInternalErrorBytes(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsCoordinatorInternalErrorBytesIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsCoordinator.contract.FilterLogs(opts, "InternalErrorBytes", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsCoordinatorInternalErrorBytesIterator{contract: _FunctionsCoordinator.contract, event: "InternalErrorBytes", logs: logs, sub: sub}, nil
}

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) WatchInternalErrorBytes(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorInternalErrorBytes, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsCoordinator.contract.WatchLogs(opts, "InternalErrorBytes", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsCoordinatorInternalErrorBytes)
				if err := _FunctionsCoordinator.contract.UnpackLog(event, "InternalErrorBytes", log); err != nil {
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

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) ParseInternalErrorBytes(log types.Log) (*FunctionsCoordinatorInternalErrorBytes, error) {
	event := new(FunctionsCoordinatorInternalErrorBytes)
	if err := _FunctionsCoordinator.contract.UnpackLog(event, "InternalErrorBytes", log); err != nil {
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

type FunctionsCoordinatorRequestTimedOutIterator struct {
	Event *FunctionsCoordinatorRequestTimedOut

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsCoordinatorRequestTimedOutIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsCoordinatorRequestTimedOut)
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
		it.Event = new(FunctionsCoordinatorRequestTimedOut)
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

func (it *FunctionsCoordinatorRequestTimedOutIterator) Error() error {
	return it.fail
}

func (it *FunctionsCoordinatorRequestTimedOutIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsCoordinatorRequestTimedOut struct {
	RequestId [32]byte
	Raw       types.Log
}

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) FilterRequestTimedOut(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsCoordinatorRequestTimedOutIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsCoordinator.contract.FilterLogs(opts, "RequestTimedOut", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsCoordinatorRequestTimedOutIterator{contract: _FunctionsCoordinator.contract, event: "RequestTimedOut", logs: logs, sub: sub}, nil
}

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) WatchRequestTimedOut(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorRequestTimedOut, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsCoordinator.contract.WatchLogs(opts, "RequestTimedOut", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsCoordinatorRequestTimedOut)
				if err := _FunctionsCoordinator.contract.UnpackLog(event, "RequestTimedOut", log); err != nil {
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

func (_FunctionsCoordinator *FunctionsCoordinatorFilterer) ParseRequestTimedOut(log types.Log) (*FunctionsCoordinatorRequestTimedOut, error) {
	event := new(FunctionsCoordinatorRequestTimedOut)
	if err := _FunctionsCoordinator.contract.UnpackLog(event, "RequestTimedOut", log); err != nil {
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
	GasOverheadAfterCallback            *big.Int
	FallbackNativePerUnitLink           *big.Int
	GasOverheadBeforeCallback           uint32
	LinkPriceFeed                       common.Address
	MaxSupportedRequestDataVersion      uint16
	FulfillmentGasPriceOverEstimationBP *big.Int
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
	case _FunctionsCoordinator.abi.Events["InternalError"].ID:
		return _FunctionsCoordinator.ParseInternalError(log)
	case _FunctionsCoordinator.abi.Events["InternalErrorBytes"].ID:
		return _FunctionsCoordinator.ParseInternalErrorBytes(log)
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
	case _FunctionsCoordinator.abi.Events["RequestTimedOut"].ID:
		return _FunctionsCoordinator.ParseRequestTimedOut(log)
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
	return common.HexToHash("0xe2cc993cebc49ce369b0174671d90dbae2201decdae368805de999b437336dcd")
}

func (FunctionsCoordinatorConfigChanged) Topic() common.Hash {
	return common.HexToHash("0x9136fb3a39fbe82d27ef8048354b8fcc93cf828f72e5c877d8135fde8dc7fd73")
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

func (FunctionsCoordinatorInternalError) Topic() common.Hash {
	return common.HexToHash("0xfb78cb471c8bc73de4ca3325ff527cfac9462f08cf259c153f25b3a51afabd79")
}

func (FunctionsCoordinatorInternalErrorBytes) Topic() common.Hash {
	return common.HexToHash("0x0efeb0e5bf272b1b81d86c1dc2b6e73cbcd571830561eb4646d8195ed5cf6df1")
}

func (FunctionsCoordinatorInvalidRequestID) Topic() common.Hash {
	return common.HexToHash("0xa1c120e327c9ad8b075793878c88d59b8934b97ae37117faa3bb21616237f7be")
}

func (FunctionsCoordinatorOracleRequest) Topic() common.Hash {
	return common.HexToHash("0xd6bf2929a5458fc22a90821275e4f54e2566c579e4bb3d38184880bf1701762c")
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

func (FunctionsCoordinatorRequestTimedOut) Topic() common.Hash {
	return common.HexToHash("0xf1ca1e9147be737b04a2b018a79405f687a97de8dd8a2559bbe62357343af414")
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

	FilterInternalError(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsCoordinatorInternalErrorIterator, error)

	WatchInternalError(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorInternalError, requestId [][32]byte) (event.Subscription, error)

	ParseInternalError(log types.Log) (*FunctionsCoordinatorInternalError, error)

	FilterInternalErrorBytes(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsCoordinatorInternalErrorBytesIterator, error)

	WatchInternalErrorBytes(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorInternalErrorBytes, requestId [][32]byte) (event.Subscription, error)

	ParseInternalErrorBytes(log types.Log) (*FunctionsCoordinatorInternalErrorBytes, error)

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

	FilterRequestTimedOut(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsCoordinatorRequestTimedOutIterator, error)

	WatchRequestTimedOut(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorRequestTimedOut, requestId [][32]byte) (event.Subscription, error)

	ParseRequestTimedOut(log types.Log) (*FunctionsCoordinatorRequestTimedOut, error)

	FilterTransmitted(opts *bind.FilterOpts) (*FunctionsCoordinatorTransmittedIterator, error)

	WatchTransmitted(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorTransmitted) (event.Subscription, error)

	ParseTransmitted(log types.Log) (*FunctionsCoordinatorTransmitted, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
