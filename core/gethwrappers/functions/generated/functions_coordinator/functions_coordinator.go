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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"linkToNativeFeed\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"EmptyPublicKey\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptyRequestData\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InconsistentReportData\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"InvalidConfig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"linkWei\",\"type\":\"int256\"}],\"name\":\"InvalidLinkWeiPrice\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSubscription\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"MustBeSubOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoTransmittersSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByRouter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByRouterOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentTooLarge\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReportInvalid\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RouterMustBeSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedPublicKeyChange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedSender\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnsupportedRequestDataVersion\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"signerPayment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"transmitterPayment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalCost\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"enumFulfillResult\",\"name\":\"result\",\"type\":\"uint8\"}],\"name\":\"BillingEnd\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"don\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"adminFee\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"estimatedTotalCostJuels\",\"type\":\"uint96\"},{\"internalType\":\"uint80\",\"name\":\"donFee\",\"type\":\"uint80\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint40\",\"name\":\"gasOverhead\",\"type\":\"uint40\"},{\"internalType\":\"uint256\",\"name\":\"expectedGasPrice\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structFunctionsBilling.Commitment\",\"name\":\"commitment\",\"type\":\"tuple\"}],\"name\":\"BillingStart\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maxCallbackGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"feedStalenessSeconds\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasOverheadBeforeCallback\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasOverheadAfterCallback\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"requestTimeoutSeconds\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint80\",\"name\":\"donFee\",\"type\":\"uint80\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"maxSupportedRequestDataVersion\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fulfillmentGasPriceOverEstimationBP\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"fallbackNativePerUnitLink\",\"type\":\"int256\"}],\"name\":\"ConfigChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"CostExceedsCommitment\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"InsufficientGasProvided\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"InsufficientSubscriptionBalance\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"InvalidRequestID\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requestingContract\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requestInitiator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"subscriptionOwner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"dataVersion\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"flags\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"callbackGasLimit\",\"type\":\"uint64\"}],\"name\":\"OracleRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"OracleResponse\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"RequestTimedOut\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"deleteCommitment\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"node\",\"type\":\"address\"}],\"name\":\"deleteNodePublicKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"}],\"name\":\"estimateCost\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"expectedGasPrice\",\"type\":\"uint256\"}],\"internalType\":\"structIFunctionsBilling.RequestBilling\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"getAdminFee\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllNodePublicKeys\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"},{\"internalType\":\"bytes[]\",\"name\":\"\",\"type\":\"bytes[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfig\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"maxCallbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"feedStalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasOverheadBeforeCallback\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasOverheadAfterCallback\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"requestTimeoutSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint80\",\"name\":\"donFee\",\"type\":\"uint80\"},{\"internalType\":\"uint16\",\"name\":\"maxSupportedRequestDataVersion\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"fulfillmentGasPriceOverEstimationBP\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"fallbackNativePerUnitLink\",\"type\":\"int256\"},{\"internalType\":\"address\",\"name\":\"linkPriceFeed\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfigHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"config\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"expectedGasPrice\",\"type\":\"uint256\"}],\"internalType\":\"structIFunctionsBilling.RequestBilling\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"getDONFee\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDONPublicKey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFeedData\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getThresholdPublicKey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"oracleWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"requestingContract\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"subscriptionOwner\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"dataVersion\",\"type\":\"uint16\"},{\"internalType\":\"bytes32\",\"name\":\"flags\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"}],\"internalType\":\"structIFunctionsCoordinator.Request\",\"name\":\"request\",\"type\":\"tuple\"}],\"name\":\"sendRequest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint96\",\"name\":\"estimatedCost\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requestTimeoutSeconds\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"_f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"_onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"_offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"_offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"donPublicKey\",\"type\":\"bytes\"}],\"name\":\"setDONPublicKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"node\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"}],\"name\":\"setNodePublicKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"thresholdPublicKey\",\"type\":\"bytes\"}],\"name\":\"setThresholdPublicKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"transmitters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"name\":\"updateConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60c06040523480156200001157600080fd5b5060405162005d1f38038062005d1f833981016040819052620000349162000462565b828282828260013380600081620000925760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c557620000c58162000154565b50505015156080526001600160a01b038216620000f557604051632530e88560e11b815260040160405180910390fd5b6001600160a01b03821660a0526200010d81620001ff565b805160209091012060075550600d80546001600160a01b039092166c01000000000000000000000000026001600160601b0390921691909117905550620006349350505050565b336001600160a01b03821603620001ae5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000089565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60008060008060008060008060008980602001905181019062000223919062000572565b985098509850985098509850985098509850600081136200025b576040516321ea67b360e11b81526004810182905260240162000089565b604080516101208101825263ffffffff808c168083528b8216602084018190528b83168486018190528b841660608601819052938b16608086018190526001600160501b038b1660a0870181905261ffff8b1660c0880181905260e088018b905261010090970189905260098054600160f01b9098026001600160f01b03600160a01b909302600160a01b600160f01b0319600160801b90950294909416600160801b600160f01b03196c0100000000000000000000000090990263ffffffff60601b196801000000000000000090970296909616600160401b600160801b03196401000000009098026001600160401b0319909b169098179990991795909516959095179290921794909416949094179290921792909216179055600a839055600b829055517f3332571032ee658b30d567e63b2443e9d921162625e0f3fe86968a586c1a090f906200041b908b908b908b908b908b908b908b908b908b9063ffffffff998a1681529789166020890152958816604088015293871660608701529190951660808501526001600160501b039490941660a084015261ffff9390931660c083015260e08201929092526101008101919091526101200190565b60405180910390a150505050505050505050565b80516001600160a01b03811681146200044757600080fd5b919050565b634e487b7160e01b600052604160045260246000fd5b6000806000606084860312156200047857600080fd5b62000483846200042f565b602085810151919450906001600160401b0380821115620004a357600080fd5b818701915087601f830112620004b857600080fd5b815181811115620004cd57620004cd6200044c565b604051601f8201601f19908116603f01168101908382118183101715620004f857620004f86200044c565b816040528281528a868487010111156200051157600080fd5b600093505b8284101562000535578484018601518185018701529285019262000516565b600086848301015280975050505050505062000554604085016200042f565b90509250925092565b805163ffffffff811681146200044757600080fd5b60008060008060008060008060006101208a8c0312156200059257600080fd5b6200059d8a6200055d565b9850620005ad60208b016200055d565b9750620005bd60408b016200055d565b9650620005cd60608b016200055d565b9550620005dd60808b016200055d565b60a08b01519095506001600160501b0381168114620005fb57600080fd5b60c08b015190945061ffff811681146200061457600080fd5b8093505060e08a015191506101008a015190509295985092959850929598565b60805160a0516156946200068b60003960008181610a6f01528181610e3701528181611074015281816110fe01528181611ba801528181611dac01528181612f6c0152613e73015260006115a301526156946000f3fe608060405234801561001057600080fd5b50600436106101ae5760003560e01c80638cc6acce116100ee578063b1dc65a411610097578063d328a91e11610071578063d328a91e1461054e578063de9dfa4614610556578063e3d0e71214610569578063f2fde38b1461057c57600080fd5b8063b1dc65a41461040d578063c3f909d414610420578063d227d2451461051e57600080fd5b8063aecc12c5116100c8578063aecc12c5146103a5578063af1296d3146103e5578063afcb95d7146103ed57600080fd5b80638cc6acce146103585780638da5cb5b1461036b5780639883c10d1461039357600080fd5b806379ba50971161015b578063814118341161013557806381411834146102ab57806381f1b938146102c057806381ff7048146102c857806385b214cf1461033557600080fd5b806379ba50971461027d5780637f15e16614610285578063807560311461029857600080fd5b8063466557181161018c5780634665571814610226578063533989871461025457806366316d8d1461026a57600080fd5b8063083a5466146101b3578063181f5a77146101c857806326ceabac14610213575b600080fd5b6101c66101c1366004614134565b61058f565b005b60408051808201909152601881527f46756e6374696f6e7320436f6f7264696e61746f72207631000000000000000060208201525b60405161020a91906141da565b60405180910390f35b6101c661022136600461420f565b6105e4565b61023961023436600461436b565b610686565b60405169ffffffffffffffffffff909116815260200161020a565b61025c6106b4565b60405161020a929190614492565b6101c661027836600461453c565b610938565b6101c6610ad4565b6101c6610293366004614134565b610bd6565b6101c66102a6366004614575565b610c26565b6102b3610cdd565b60405161020a91906145ca565b6101fd610d4c565b61031260015460025463ffffffff74010000000000000000000000000000000000000000830481169378010000000000000000000000000000000000000000000000009093041691565b6040805163ffffffff94851681529390921660208401529082015260600161020a565b6103486103433660046145dd565b610e1d565b604051901515815260200161020a565b6101c66103663660046145f6565b61105c565b60005460405173ffffffffffffffffffffffffffffffffffffffff909116815260200161020a565b6007545b60405190815260200161020a565b6103b86103b3366004614633565b6110e1565b604080519485526bffffffffffffffffffffffff909316602085015291830152606082015260800161020a565b61039761133d565b60408051600181526000602082018190529181019190915260600161020a565b6101c661041b3660046146b3565b611430565b600954600a54600b54600d546040805163ffffffff80871682526401000000008704811660208301526801000000000000000087048116928201929092526c01000000000000000000000000808704831660608301527001000000000000000000000000000000008704909216608082015274010000000000000000000000000000000000000000860469ffffffffffffffffffff1660a08201527e0100000000000000000000000000000000000000000000000000000000000090950461ffff1660c086015260e08501939093526101008401919091520473ffffffffffffffffffffffffffffffffffffffff166101208201526101400161020a565b61053161052c36600461476a565b611b61565b6040516bffffffffffffffffffffffff909116815260200161020a565b6101fd611d51565b61053161056436600461436b565b611da8565b6101c661057736600461488c565b611e40565b6101c661058a36600461420f565b612864565b610597612875565b60008190036105d2576040517f4f42be3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60106105df8284836149fa565b505050565b60005473ffffffffffffffffffffffffffffffffffffffff1633148061061f57503373ffffffffffffffffffffffffffffffffffffffff8216145b610655576040517fed6dd19b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff81166000908152600f6020526040812061068391614079565b50565b60095474010000000000000000000000000000000000000000900469ffffffffffffffffffff165b92915050565b6060806000600680548060200260200160405190810160405280929190818152602001828054801561071c57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116106f1575b505050505090506000815167ffffffffffffffff81111561073f5761073f61422c565b60405190808252806020026020018201604052801561077257816020015b606081526020019060019003908161075d5790505b50905060005b825181101561092e57600f600084838151811061079757610797614b15565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002080546107e490614959565b905060000361081f576040517f4f42be3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600f600084838151811061083557610835614b15565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020805461088290614959565b80601f01602080910402602001604051908101604052809291908181526020018280546108ae90614959565b80156108fb5780601f106108d0576101008083540402835291602001916108fb565b820191906000526020600020905b8154815290600101906020018083116108de57829003601f168201915b505050505082828151811061091257610912614b15565b60200260200101819052508061092790614b73565b9050610778565b5090939092509050565b6109406128f8565b806bffffffffffffffffffffffff166000036109765750336000908152600c60205260409020546bffffffffffffffffffffffff165b336000908152600c60205260409020546bffffffffffffffffffffffff808316911610156109d0576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b336000908152600c6020526040812080548392906109fd9084906bffffffffffffffffffffffff16614bab565b82546101009290920a6bffffffffffffffffffffffff8181021990931691831602179091556040517f66316d8d00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff858116600483015291841660248201527f00000000000000000000000000000000000000000000000000000000000000009250908216906366316d8d90604401600060405180830381600087803b158015610ab757600080fd5b505af1158015610acb573d6000803e3d6000fd5b50505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610b5a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610bde612875565b6000819003610c19576040517f4f42be3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600e6105df8284836149fa565b60005473ffffffffffffffffffffffffffffffffffffffff16331480610c715750610c5033612a8f565b8015610c7157503373ffffffffffffffffffffffffffffffffffffffff8416145b610ca7576040517fed6dd19b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff83166000908152600f60205260409020610cd78284836149fa565b50505050565b60606006805480602002602001604051908101604052809291908181526020018280548015610d4257602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610d17575b5050505050905090565b606060108054610d5b90614959565b9050600003610d96576040517f4f42be3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60108054610da390614959565b80601f0160208091040260200160405190810160405280929190818152602001828054610dcf90614959565b8015610d425780601f10610df157610100808354040283529160200191610d42565b820191906000526020600020905b815481529060010190602001808311610dff57509395945050505050565b60003373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610e8e576040517fc41a5b0900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260086020908152604091829020825161014081018452815467ffffffffffffffff8116825268010000000000000000810473ffffffffffffffffffffffffffffffffffffffff908116948301949094527c0100000000000000000000000000000000000000000000000000000000900463ffffffff90811694820194909452600182015492831660608201819052740100000000000000000000000000000000000000009093046bffffffffffffffffffffffff9081166080830152600283015490811660a08301526c01000000000000000000000000810469ffffffffffffffffffff1660c0830152760100000000000000000000000000000000000000000000810490941660e08201527a01000000000000000000000000000000000000000000000000000090930464ffffffffff1661010084015260030154610120830152610fe25750600092915050565b600083815260086020526040808220828155600181018390556002810180547fff000000000000000000000000000000000000000000000000000000000000001690556003018290555184917ff1ca1e9147be737b04a2b018a79405f687a97de8dd8a2559bbe62357343af41491a260019150505b919050565b3373ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016146110cb576040517fc41a5b0900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6110d481612b78565b8051602090910120600755565b60008080803373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614611155576040517fc41a5b0900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6111626040860186614bd7565b905060000361119c576040517ec1cfc000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600060405180608001604052808760600160208101906111bc9190614c3c565b67ffffffffffffffff1681526020908101906111da9089018961420f565b73ffffffffffffffffffffffffffffffffffffffff16815260200161120560e0890160c08a01614c59565b63ffffffff1681523a60209091015290506112716112266040880188614bd7565b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061126b9250505060a0890160808a01614c86565b83612e81565b92975090955093509150611288602087018761420f565b73ffffffffffffffffffffffffffffffffffffffff16857fd6bf2929a5458fc22a90821275e4f54e2566c579e4bb3d38184880bf1701762c326112d160808b0160608c01614c3c565b6112e160408c0160208d0161420f565b6112ee60408d018d614bd7565b8d60800160208101906113019190614c86565b8e60a001358f60c00160208101906113199190614c59565b60405161132d989796959493929190614ca3565b60405180910390a3509193509193565b600954600d54604080517ffeaf968c0000000000000000000000000000000000000000000000000000000081529051600093640100000000900463ffffffff169283151592859283926c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff169163feaf968c9160048083019260a09291908290030181865afa1580156113d7573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906113fb9190614d63565b5093505092505082801561141d57506114148142614db3565b8463ffffffff16105b1561142857600b5491505b509392505050565b60005a604080516020601f8b018190048102820181019092528981529192508a3591818c01359161148691849163ffffffff851691908e908e908190840183828082843760009201919091525061352c92505050565b6114bc576040517f0be3632800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805183815262ffffff600884901c1660208201527fb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62910160405180910390a16040805160608101825260025480825260035460ff80821660208501526101009091041692820192909252908314611591576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f636f6e666967446967657374206d69736d6174636800000000000000000000006044820152606401610b51565b61159f8b8b8b8b8b8b613535565b60007f0000000000000000000000000000000000000000000000000000000000000000156115fc576002826020015183604001516115dd9190614dc6565b6115e79190614e0e565b6115f2906001614dc6565b60ff169050611612565b602082015161160c906001614dc6565b60ff1690505b88811461167b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f77726f6e67206e756d626572206f66207369676e6174757265730000000000006044820152606401610b51565b8887146116e4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f7369676e617475726573206f7574206f6620726567697374726174696f6e00006044820152606401610b51565b3360009081526004602090815260408083208151808301909252805460ff8082168452929391929184019161010090910416600281111561172757611727614e30565b600281111561173857611738614e30565b905250905060028160200151600281111561175557611755614e30565b14801561179c57506006816000015160ff168154811061177757611777614b15565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff1633145b611802576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f756e617574686f72697a6564207472616e736d697474657200000000000000006044820152606401610b51565b505050505061180f6140b3565b6000808a8a604051611822929190614e5f565b604051908190038120611839918e90602001614e6f565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181528282528051602091820120838301909252600080845290830152915060005b89811015611b435760006001848984602081106118a2576118a2614b15565b6118af91901a601b614dc6565b8e8e868181106118c1576118c1614b15565b905060200201358d8d878181106118da576118da614b15565b9050602002013560405160008152602001604052604051611917949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa158015611939573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff811660009081526004602090815290849020838501909452835460ff808216855292965092945084019161010090041660028111156119b9576119b9614e30565b60028111156119ca576119ca614e30565b90525092506001836020015160028111156119e7576119e7614e30565b14611a4e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f61646472657373206e6f7420617574686f72697a656420746f207369676e00006044820152606401610b51565b8251600090879060ff16601f8110611a6857611a68614b15565b602002015173ffffffffffffffffffffffffffffffffffffffff1614611aea576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f6e6f6e2d756e69717565207369676e61747572650000000000000000000000006044820152606401610b51565b8086846000015160ff16601f8110611b0457611b04614b15565b73ffffffffffffffffffffffffffffffffffffffff9092166020929092020152611b2f600186614dc6565b94505080611b3c90614b73565b9050611883565b505050611b54833383858e8e6135e3565b5050505050505050505050565b6040517f10fc49c100000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8616600482015263ffffffff831660248201526000907f00000000000000000000000000000000000000000000000000000000000000009073ffffffffffffffffffffffffffffffffffffffff8216906310fc49c19060440160006040518083038186803b158015611c0157600080fd5b505afa158015611c15573d6000803e3d6000fd5b50505050620f4240831115611c56576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600060405180608001604052808967ffffffffffffffff1681526020013373ffffffffffffffffffffffffffffffffffffffff1681526020018663ffffffff1681526020018581525090506000611ce488888080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250869250610686915050565b69ffffffffffffffffffff1690506000611d3589898080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250879250611da8915050565b9050611d43878784846138ea565b9a9950505050505050505050565b6060600e8054611d6090614959565b9050600003611d9b576040517f4f42be3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600e8054610da390614959565b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16632a905ccc6040518163ffffffff1660e01b8152600401602060405180830381865afa158015611e15573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611e399190614e83565b9392505050565b855185518560ff16601f831115611eb3576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f746f6f206d616e79207369676e657273000000000000000000000000000000006044820152606401610b51565b80600003611f1d576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f66206d75737420626520706f73697469766500000000000000000000000000006044820152606401610b51565b818314611fab576040517f89a61989000000000000000000000000000000000000000000000000000000008152602060048201526024808201527f6f7261636c6520616464726573736573206f7574206f6620726567697374726160448201527f74696f6e000000000000000000000000000000000000000000000000000000006064820152608401610b51565b611fb6816003614ea0565b831161201e576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f6661756c74792d6f7261636c65206620746f6f206869676800000000000000006044820152606401610b51565b612026612875565b6040805160c0810182528a8152602081018a905260ff89169181018290526060810188905267ffffffffffffffff8716608082015260a081018690529061206d9088613a4c565b600554156122225760055460009061208790600190614db3565b905060006005828154811061209e5761209e614b15565b60009182526020822001546006805473ffffffffffffffffffffffffffffffffffffffff909216935090849081106120d8576120d8614b15565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff85811684526004909252604080842080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00009081169091559290911680845292208054909116905560058054919250908061215857612158614eb7565b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff000000000000000000000000000000000000000016905501905560068054806121c1576121c1614eb7565b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff00000000000000000000000000000000000000001690550190555061206d915050565b60005b815151811015612689576000600460008460000151848151811061224b5761224b614b15565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff16825281019190915260400160002054610100900460ff16600281111561229557612295614e30565b146122fc576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265706561746564207369676e657220616464726573730000000000000000006044820152606401610b51565b6040805180820190915260ff8216815260016020820152825180516004916000918590811061232d5761232d614b15565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040016000208251815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082168117835592840151919283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000016176101008360028111156123ce576123ce614e30565b0217905550600091506123de9050565b60046000846020015184815181106123f8576123f8614b15565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff16825281019190915260400160002054610100900460ff16600281111561244257612442614e30565b146124a9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f7265706561746564207472616e736d69747465722061646472657373000000006044820152606401610b51565b6040805180820190915260ff8216815260208101600281525060046000846020015184815181106124dc576124dc614b15565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040016000208251815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082168117835592840151919283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000161761010083600281111561257d5761257d614e30565b02179055505082518051600592508390811061259b5761259b614b15565b602090810291909101810151825460018101845560009384529282902090920180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909316929092179091558201518051600691908390811061261757612617614b15565b60209081029190910181015182546001810184556000938452919092200180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff9092169190911790558061268181614b73565b915050612225565b506040810151600380547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff909216919091179055600180547fffffffff00000000ffffffffffffffffffffffffffffffffffffffffffffffff8116780100000000000000000000000000000000000000000000000063ffffffff438116820292909217808555920481169291829160149161274191849174010000000000000000000000000000000000000000900416614ee6565b92506101000a81548163ffffffff021916908363ffffffff1602179055506127a04630600160149054906101000a900463ffffffff1663ffffffff16856000015186602001518760400151886060015189608001518a60a00151613a69565b600281905582518051600380547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1661010060ff9093169290920291909117905560015460208501516040808701516060880151608089015160a08a015193517f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e0598612857988b9891977401000000000000000000000000000000000000000090920463ffffffff16969095919491939192614f03565b60405180910390a1611b54565b61286c612875565b61068381613b14565b60005473ffffffffffffffffffffffffffffffffffffffff1633146128f6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610b51565b565b6000612902610cdd565b9050805160000361293f576040517f30274b3a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8051600d5460009161295e916bffffffffffffffffffffffff16614f89565b905060005b82518160ff161015612a305781600c6000858460ff168151811061298957612989614b15565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282829054906101000a90046bffffffffffffffffffffffff166129f19190614fb4565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555080612a2990614fd9565b9050612963565b508151612a3d9082614ff8565b600d8054600090612a5d9084906bffffffffffffffffffffffff16614bab565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055505050565b6000806006805480602002602001604051908101604052809291908181526020018280548015612af557602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311612aca575b5050505050905060005b8151811015612b6e578373ffffffffffffffffffffffffffffffffffffffff16828281518110612b3157612b31614b15565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1603612b5e575060019392505050565b612b6781614b73565b9050612aff565b5060009392505050565b600080600080600080600080600089806020019051810190612b9a9190615028565b98509850985098509850985098509850985060008113612be9576040517f43d4cf6600000000000000000000000000000000000000000000000000000000815260048101829052602401610b51565b604080516101208101825263ffffffff808c168083528b8216602084018190528b83168486018190528b841660608601819052938b166080860181905269ffffffffffffffffffff8b1660a0870181905261ffff8b1660c0880181905260e088018b9052610100909701899052600980547e010000000000000000000000000000000000000000000000000000000000009098027dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff740100000000000000000000000000000000000000009093027fffff00000000000000000000ffffffffffffffffffffffffffffffffffffffff700100000000000000000000000000000000909502949094167fffff0000000000000000000000000000ffffffffffffffffffffffffffffffff6c010000000000000000000000009099027fffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffffff68010000000000000000909702969096167fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff6401000000009098027fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000909b169098179990991795909516959095179290921794909416949094179290921792909216179055600a839055600b829055517f3332571032ee658b30d567e63b2443e9d921162625e0f3fe86968a586c1a090f90612e6d908b908b908b908b908b908b908b908b908b9063ffffffff998a16815297891660208901529588166040880152938716606087015291909516608085015269ffffffffffffffffffff9490941660a084015261ffff9390931660c083015260e08201929092526101008101919091526101200190565b60405180910390a150505050505050505050565b60095460009081908190819061ffff7e0100000000000000000000000000000000000000000000000000000000000090910481169087161115612ef0576040517fdada758700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000612efc8887610686565b90506000612f0a8988611da8565b90506000612f2e886040015189606001518569ffffffffffffffffffff16856138ea565b88516040517fa47c769600000000000000000000000000000000000000000000000000000000815267ffffffffffffffff90911660048201529091507f000000000000000000000000000000000000000000000000000000000000000090600090819073ffffffffffffffffffffffffffffffffffffffff84169063a47c769690602401600060405180830381865afa158015612fcf573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405261301591908101906150d4565b50505060208d01518d516040517f674603d000000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff928316600482015267ffffffffffffffff90911660248201529294509092506000919085169063674603d090604401606060405180830381865afa1580156130a6573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906130ca91906151c0565b509150506bffffffffffffffffffffffff85166130e78385614bab565b6bffffffffffffffffffffffff16101561312d576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006131ac308e602001518f6000015185600161314a9190615212565b6040805173ffffffffffffffffffffffffffffffffffffffff958616602080830191909152949095168582015267ffffffffffffffff928316606086015291166080808501919091528151808503909101815260a09093019052815191012090565b905060006040518061014001604052808f6000015167ffffffffffffffff1681526020018f6020015173ffffffffffffffffffffffffffffffffffffffff1681526020018f6040015163ffffffff1681526020013073ffffffffffffffffffffffffffffffffffffffff168152602001896bffffffffffffffffffffffff168152602001886bffffffffffffffffffffffff1681526020018a69ffffffffffffffffffff1681526020014263ffffffff1681526020016009600001600c9054906101000a900463ffffffff16600960000160089054906101000a900463ffffffff166132989190614ee6565b63ffffffff1664ffffffffff1681526020018f606001518152509050806008600084815260200190815260200160002060008201518160000160006101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555060208201518160000160086101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550604082015181600001601c6101000a81548163ffffffff021916908363ffffffff16021790555060608201518160010160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060808201518160010160146101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060a08201518160020160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060c082015181600201600c6101000a81548169ffffffffffffffffffff021916908369ffffffffffffffffffff16021790555060e08201518160020160166101000a81548163ffffffff021916908363ffffffff16021790555061010082015181600201601a6101000a81548164ffffffffff021916908364ffffffffff1602179055506101208201518160030155905050817fe2cc993cebc49ce369b0174671d90dbae2201decdae368805de999b437336dcd826040516134dc9190615233565b60405180910390a250600954909f959e5063ffffffff6c01000000000000000000000000820481169e50700100000000000000000000000000000000909104169b50939950505050505050505050565b60019392505050565b6000613542826020614ea0565b61354d856020614ea0565b61355988610144615348565b6135639190615348565b61356d9190615348565b613578906000615348565b9050368114610acb576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f63616c6c64617461206c656e677468206d69736d6174636800000000000000006044820152606401610b51565b606080806135f3848601866153db565b82519295509093509150158061360b57508151835114155b8061361857508051835114155b1561364f576040517f0be3632800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b83518110156138de5760006136b385838151811061367257613672614b15565b602002602001015185848151811061368c5761368c614b15565b60200260200101518585815181106136a6576136a6614b15565b6020026020010151613c09565b905060008160068111156136c9576136c9614e30565b14806136e6575060018160068111156136e4576136e4614e30565b145b15613741578482815181106136fd576136fd614b15565b60209081029190910181015160405133815290917fc708e0440951fd63499c0f7a73819b469ee5dd3ecc356c0ab4eb7f18389009d9910160405180910390a26138cd565b600281600681111561375557613755614e30565b036137a55784828151811061376c5761376c614b15565b60200260200101517fa1c120e327c9ad8b075793878c88d59b8934b97ae37117faa3bb21616237f7be60405160405180910390a26138cd565b60038160068111156137b9576137b9614e30565b03613809578482815181106137d0576137d0614b15565b60200260200101517f357a60574f143e144d62ddc89636fee45ba5e2f8c9a30b9f91e28f8dfe5a56e060405160405180910390a26138cd565b600581600681111561381d5761381d614e30565b0361386d5784828151811061383457613834614b15565b60200260200101517f689f49eacf6db85082542f4b7873f4a519744a4e093fd002310545bfb4354cc860405160405180910390a26138cd565b600481600681111561388157613881614e30565b036138cd5784828151811061389857613898614b15565b60200260200101517fb83d9204606ee3780daa106ad026c4c2853f5af126edbb8fe15e813e62dbd94d60405160405180910390a25b506138d781614b73565b9050613652565b50505050505050505050565b6000806138f561133d565b905060008113613934576040517f43d4cf6600000000000000000000000000000000000000000000000000000000815260048101829052602401610b51565b600954600090879061396a9063ffffffff6c01000000000000000000000000820481169168010000000000000000900416614ee6565b6139749190614ee6565b63ffffffff1690506000612710600960010154886139929190614ea0565b61399c91906154b8565b6139a69088615348565b9050600083836139be84670de0b6b3a7640000614ea0565b6139c89190614ea0565b6139d291906154b8565b905060006139f16bffffffffffffffffffffffff808916908a16615348565b9050613a09816b033b2e3c9fd0803ce8000000614db3565b821115613a42576040517fe80fa38100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611d438183615348565b6000613a56610cdd565b511115613a6557613a656128f8565b5050565b6000808a8a8a8a8a8a8a8a8a604051602001613a8d999897969594939291906154cc565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179150509998505050505050505050565b3373ffffffffffffffffffffffffffffffffffffffff821603613b93576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610b51565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000838152600860209081526040808320815161014081018352815467ffffffffffffffff8116825273ffffffffffffffffffffffffffffffffffffffff68010000000000000000820481169583019590955263ffffffff7c01000000000000000000000000000000000000000000000000000000009091048116938201939093526001820154938416606082018190526bffffffffffffffffffffffff7401000000000000000000000000000000000000000090950485166080830152600283015494851660a083015269ffffffffffffffffffff6c0100000000000000000000000086041660c0830152760100000000000000000000000000000000000000000000850490931660e082015264ffffffffff7a0100000000000000000000000000000000000000000000000000009094049390931661010084015260030154610120830152613d5e576002915050611e39565b6000858152600860205260408120818155600181018290556002810180547fff00000000000000000000000000000000000000000000000000000000000000169055600301819055613dae61133d565b905060008113613ded576040517f43d4cf6600000000000000000000000000000000000000000000000000000000815260048101829052602401610b51565b600081613e023a670de0b6b3a7640000614ea0565b613e0c91906154b8565b9050600083610100015164ffffffffff1682613e289190614ea0565b905060008460c0015169ffffffffffffffffffff1682613e489190614fb4565b6040517f3fd67e050000000000000000000000000000000000000000000000000000000081529091507f000000000000000000000000000000000000000000000000000000000000000090600090819073ffffffffffffffffffffffffffffffffffffffff841690633fd67e0590613ece908f908f908f908c908b903390600401615561565b60408051808303816000875af1158015613eec573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613f1091906155d2565b9092509050613f1f8186614fb4565b336000908152600c6020526040812080547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff93841617905560c08a0151600d805469ffffffffffffffffffff90921693909291613f8c91859116614fb4565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055508b7f30e1e299f390e4f3fa75e4489fed477e86183937a4a21a0ae7b137a4dff0353589600001518a60c001518489613ff29190614fb4565b60808d015160c08e015169ffffffffffffffffffff16614012888d614fb4565b61401c9190614fb4565b6140269190614fb4565b8760ff16600681111561403b5761403b614e30565b60405161404c959493929190615601565b60405180910390a28160ff16600681111561406957614069614e30565b9c9b505050505050505050505050565b50805461408590614959565b6000825580601f10614095575050565b601f01602090049060005260206000209081019061068391906140d2565b604051806103e00160405280601f906020820280368337509192915050565b5b808211156140e757600081556001016140d3565b5090565b60008083601f8401126140fd57600080fd5b50813567ffffffffffffffff81111561411557600080fd5b60208301915083602082850101111561412d57600080fd5b9250929050565b6000806020838503121561414757600080fd5b823567ffffffffffffffff81111561415e57600080fd5b61416a858286016140eb565b90969095509350505050565b6000815180845260005b8181101561419c57602081850181015186830182015201614180565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b602081526000611e396020830184614176565b73ffffffffffffffffffffffffffffffffffffffff8116811461068357600080fd5b60006020828403121561422157600080fd5b8135611e39816141ed565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156142a2576142a261422c565b604052919050565b600082601f8301126142bb57600080fd5b813567ffffffffffffffff8111156142d5576142d561422c565b61430660207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8401160161425b565b81815284602083860101111561431b57600080fd5b816020850160208301376000918101602001919091529392505050565b67ffffffffffffffff8116811461068357600080fd5b803561105781614338565b63ffffffff8116811461068357600080fd5b60008082840360a081121561437f57600080fd5b833567ffffffffffffffff8082111561439757600080fd5b6143a3878388016142aa565b945060807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0840112156143d557600080fd5b604051925060808301915082821081831117156143f4576143f461422c565b50604052602084013561440681614338565b81526040840135614416816141ed565b6020820152606084013561442981614359565b60408201526080939093013560608401525092909150565b600081518084526020808501945080840160005b8381101561448757815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101614455565b509495945050505050565b6040815260006144a56040830185614441565b6020838203818501528185518084528284019150828160051b85010183880160005b83811015614513577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0878403018552614501838351614176565b948601949250908501906001016144c7565b50909998505050505050505050565b6bffffffffffffffffffffffff8116811461068357600080fd5b6000806040838503121561454f57600080fd5b823561455a816141ed565b9150602083013561456a81614522565b809150509250929050565b60008060006040848603121561458a57600080fd5b8335614595816141ed565b9250602084013567ffffffffffffffff8111156145b157600080fd5b6145bd868287016140eb565b9497909650939450505050565b602081526000611e396020830184614441565b6000602082840312156145ef57600080fd5b5035919050565b60006020828403121561460857600080fd5b813567ffffffffffffffff81111561461f57600080fd5b61462b848285016142aa565b949350505050565b60006020828403121561464557600080fd5b813567ffffffffffffffff81111561465c57600080fd5b820160e08185031215611e3957600080fd5b60008083601f84011261468057600080fd5b50813567ffffffffffffffff81111561469857600080fd5b6020830191508360208260051b850101111561412d57600080fd5b60008060008060008060008060e0898b0312156146cf57600080fd5b606089018a8111156146e057600080fd5b8998503567ffffffffffffffff808211156146fa57600080fd5b6147068c838d016140eb565b909950975060808b013591508082111561471f57600080fd5b61472b8c838d0161466e565b909750955060a08b013591508082111561474457600080fd5b506147518b828c0161466e565b999c989b50969995989497949560c00135949350505050565b60008060008060006080868803121561478257600080fd5b853561478d81614338565b9450602086013567ffffffffffffffff8111156147a957600080fd5b6147b5888289016140eb565b90955093505060408601356147c981614359565b949793965091946060013592915050565b600067ffffffffffffffff8211156147f4576147f461422c565b5060051b60200190565b600082601f83011261480f57600080fd5b8135602061482461481f836147da565b61425b565b82815260059290921b8401810191818101908684111561484357600080fd5b8286015b8481101561486757803561485a816141ed565b8352918301918301614847565b509695505050505050565b60ff8116811461068357600080fd5b803561105781614872565b60008060008060008060c087890312156148a557600080fd5b863567ffffffffffffffff808211156148bd57600080fd5b6148c98a838b016147fe565b975060208901359150808211156148df57600080fd5b6148eb8a838b016147fe565b96506148f960408a01614881565b9550606089013591508082111561490f57600080fd5b61491b8a838b016142aa565b945061492960808a0161434e565b935060a089013591508082111561493f57600080fd5b5061494c89828a016142aa565b9150509295509295509295565b600181811c9082168061496d57607f821691505b6020821081036149a6577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f8211156105df57600081815260208120601f850160051c810160208610156149d35750805b601f850160051c820191505b818110156149f2578281556001016149df565b505050505050565b67ffffffffffffffff831115614a1257614a1261422c565b614a2683614a208354614959565b836149ac565b6000601f841160018114614a785760008515614a425750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b178355614b0e565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b82811015614ac75786850135825560209485019460019092019101614aa7565b5086821015614b02577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203614ba457614ba4614b44565b5060010190565b6bffffffffffffffffffffffff828116828216039080821115614bd057614bd0614b44565b5092915050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112614c0c57600080fd5b83018035915067ffffffffffffffff821115614c2757600080fd5b60200191503681900382131561412d57600080fd5b600060208284031215614c4e57600080fd5b8135611e3981614338565b600060208284031215614c6b57600080fd5b8135611e3981614359565b61ffff8116811461068357600080fd5b600060208284031215614c9857600080fd5b8135611e3981614c76565b600073ffffffffffffffffffffffffffffffffffffffff808b16835267ffffffffffffffff8a16602084015280891660408401525060e060608301528560e0830152610100868882850137600083880182015261ffff9590951660808301525060a081019290925263ffffffff1660c0820152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01690910101949350505050565b805169ffffffffffffffffffff8116811461105757600080fd5b600080600080600060a08688031215614d7b57600080fd5b614d8486614d49565b9450602086015193506040860151925060608601519150614da760808701614d49565b90509295509295909350565b818103818111156106ae576106ae614b44565b60ff81811683821601908111156106ae576106ae614b44565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b600060ff831680614e2157614e21614ddf565b8060ff84160491505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b8183823760009101908152919050565b828152606082602083013760800192915050565b600060208284031215614e9557600080fd5b8151611e3981614522565b80820281158282048414176106ae576106ae614b44565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b63ffffffff818116838216019080821115614bd057614bd0614b44565b600061012063ffffffff808d1684528b6020850152808b16604085015250806060840152614f338184018a614441565b90508281036080840152614f478189614441565b905060ff871660a084015282810360c0840152614f648187614176565b905067ffffffffffffffff851660e08401528281036101008401526140698185614176565b60006bffffffffffffffffffffffff80841680614fa857614fa8614ddf565b92169190910492915050565b6bffffffffffffffffffffffff818116838216019080821115614bd057614bd0614b44565b600060ff821660ff8103614fef57614fef614b44565b60010192915050565b6bffffffffffffffffffffffff81811683821602808216919082811461502057615020614b44565b505092915050565b60008060008060008060008060006101208a8c03121561504757600080fd5b895161505281614359565b60208b015190995061506381614359565b60408b015190985061507481614359565b60608b015190975061508581614359565b60808b015190965061509681614359565b94506150a460a08b01614d49565b935060c08a01516150b481614c76565b8093505060e08a015191506101008a015190509295985092959850929598565b600080600080600060a086880312156150ec57600080fd5b85516150f781614522565b8095505060208087015161510a81614522565b604088015190955061511b816141ed565b606088015190945061512c816141ed565b608088015190935067ffffffffffffffff81111561514957600080fd5b8701601f8101891361515a57600080fd5b805161516861481f826147da565b81815260059190911b8201830190838101908b83111561518757600080fd5b928401925b828410156151ae57835161519f816141ed565b8252928401929084019061518c565b80955050505050509295509295909350565b6000806000606084860312156151d557600080fd5b835180151581146151e557600080fd5b60208501519093506151f681614338565b604085015190925061520781614338565b809150509250925092565b67ffffffffffffffff818116838216019080821115614bd057614bd0614b44565b815167ffffffffffffffff1681526101408101602083015161526d602084018273ffffffffffffffffffffffffffffffffffffffff169052565b506040830151615285604084018263ffffffff169052565b5060608301516152ad606084018273ffffffffffffffffffffffffffffffffffffffff169052565b5060808301516152cd60808401826bffffffffffffffffffffffff169052565b5060a08301516152ed60a08401826bffffffffffffffffffffffff169052565b5060c083015161530b60c084018269ffffffffffffffffffff169052565b5060e083015161532360e084018263ffffffff169052565b506101008381015164ffffffffff811684830152505061012092830151919092015290565b808201808211156106ae576106ae614b44565b600082601f83011261536c57600080fd5b8135602061537c61481f836147da565b82815260059290921b8401810191818101908684111561539b57600080fd5b8286015b8481101561486757803567ffffffffffffffff8111156153bf5760008081fd5b6153cd8986838b01016142aa565b84525091830191830161539f565b6000806000606084860312156153f057600080fd5b833567ffffffffffffffff8082111561540857600080fd5b818601915086601f83011261541c57600080fd5b8135602061542c61481f836147da565b82815260059290921b8401810191818101908a84111561544b57600080fd5b948201945b8386101561546957853582529482019490820190615450565b9750508701359250508082111561547f57600080fd5b61548b8783880161535b565b935060408601359150808211156154a157600080fd5b506154ae8682870161535b565b9150509250925092565b6000826154c7576154c7614ddf565b500490565b60006101208b835273ffffffffffffffffffffffffffffffffffffffff8b16602084015267ffffffffffffffff808b1660408501528160608501526155138285018b614441565b91508382036080850152615527828a614441565b915060ff881660a085015283820360c08501526155448288614176565b90861660e085015283810361010085015290506140698185614176565b86815260c06020820152600061557a60c0830188614176565b828103604084015261558c8188614176565b6bffffffffffffffffffffffff96871660608501529490951660808301525073ffffffffffffffffffffffffffffffffffffffff9190911660a090910152949350505050565b600080604083850312156155e557600080fd5b82516155f081614872565b602084015190925061456a81614522565b67ffffffffffffffff8616815269ffffffffffffffffffff851660208201526bffffffffffffffffffffffff84811660408301528316606082015260a0810160078310615677577f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b826080830152969550505050505056fea164736f6c6343000813000a",
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
