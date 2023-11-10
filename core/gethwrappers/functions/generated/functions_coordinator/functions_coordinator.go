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

type FunctionsCoordinatorConfig struct {
	FulfillmentGasPriceOverEstimationBP uint32
	FeedStalenessSeconds                uint32
	GasOverheadBeforeCallback           uint32
	GasOverheadAfterCallback            uint32
	DonFee                              *big.Int
	MinimumEstimateGasPriceWei          *big.Int
	MaxSupportedRequestDataVersion      uint16
	FallbackNativePerUnitLink           *big.Int
	RequestTimeoutSeconds               uint32
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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"fulfillmentGasPriceOverEstimationBP\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"feedStalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasOverheadBeforeCallback\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasOverheadAfterCallback\",\"type\":\"uint32\"},{\"internalType\":\"uint72\",\"name\":\"donFee\",\"type\":\"uint72\"},{\"internalType\":\"uint40\",\"name\":\"minimumEstimateGasPriceWei\",\"type\":\"uint40\"},{\"internalType\":\"uint16\",\"name\":\"maxSupportedRequestDataVersion\",\"type\":\"uint16\"},{\"internalType\":\"uint224\",\"name\":\"fallbackNativePerUnitLink\",\"type\":\"uint224\"},{\"internalType\":\"uint32\",\"name\":\"requestTimeoutSeconds\",\"type\":\"uint32\"}],\"internalType\":\"structFunctionsCoordinator.Config\",\"name\":\"config\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"linkToNativeFeed\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"EmptyPublicKey\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InconsistentReportData\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"InvalidConfig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"linkWei\",\"type\":\"int256\"}],\"name\":\"InvalidLinkWeiPrice\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSubscription\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"MustBeSubOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoTransmittersSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByRouter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByRouterOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentTooLarge\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"ReportInvalid\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RouterMustBeSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedPublicKeyChange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedSender\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnsupportedRequestDataVersion\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"CommitmentDeleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"fulfillmentGasPriceOverEstimationBP\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"feedStalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasOverheadBeforeCallback\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasOverheadAfterCallback\",\"type\":\"uint32\"},{\"internalType\":\"uint72\",\"name\":\"donFee\",\"type\":\"uint72\"},{\"internalType\":\"uint40\",\"name\":\"minimumEstimateGasPriceWei\",\"type\":\"uint40\"},{\"internalType\":\"uint16\",\"name\":\"maxSupportedRequestDataVersion\",\"type\":\"uint16\"},{\"internalType\":\"uint224\",\"name\":\"fallbackNativePerUnitLink\",\"type\":\"uint224\"},{\"internalType\":\"uint32\",\"name\":\"requestTimeoutSeconds\",\"type\":\"uint32\"}],\"indexed\":false,\"internalType\":\"structFunctionsCoordinator.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"ConfigUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requestingContract\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requestInitiator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"subscriptionOwner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"dataVersion\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"flags\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"callbackGasLimit\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"estimatedTotalCostJuels\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint72\",\"name\":\"adminFee\",\"type\":\"uint72\"},{\"internalType\":\"uint72\",\"name\":\"donFee\",\"type\":\"uint72\"},{\"internalType\":\"uint40\",\"name\":\"gasOverheadBeforeCallback\",\"type\":\"uint40\"},{\"internalType\":\"uint40\",\"name\":\"gasOverheadAfterCallback\",\"type\":\"uint40\"},{\"internalType\":\"uint32\",\"name\":\"timeoutTimestamp\",\"type\":\"uint32\"}],\"indexed\":false,\"internalType\":\"structFunctionsResponse.Commitment\",\"name\":\"commitment\",\"type\":\"tuple\"}],\"name\":\"OracleRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"OracleResponse\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"juelsPerGas\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"l1FeeShareWei\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"callbackCostJuels\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalCostJuels\",\"type\":\"uint96\"}],\"name\":\"RequestBilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"deleteCommitment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasPriceWei\",\"type\":\"uint256\"}],\"name\":\"estimateCost\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAdminFee\",\"outputs\":[{\"internalType\":\"uint72\",\"name\":\"\",\"type\":\"uint72\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"fulfillmentGasPriceOverEstimationBP\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"feedStalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasOverheadBeforeCallback\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasOverheadAfterCallback\",\"type\":\"uint32\"},{\"internalType\":\"uint72\",\"name\":\"donFee\",\"type\":\"uint72\"},{\"internalType\":\"uint40\",\"name\":\"minimumEstimateGasPriceWei\",\"type\":\"uint40\"},{\"internalType\":\"uint16\",\"name\":\"maxSupportedRequestDataVersion\",\"type\":\"uint16\"},{\"internalType\":\"uint224\",\"name\":\"fallbackNativePerUnitLink\",\"type\":\"uint224\"},{\"internalType\":\"uint32\",\"name\":\"requestTimeoutSeconds\",\"type\":\"uint32\"}],\"internalType\":\"structFunctionsCoordinator.Config\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"getDONFee\",\"outputs\":[{\"internalType\":\"uint72\",\"name\":\"\",\"type\":\"uint72\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDONPublicKey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getThresholdPublicKey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getWeiPerUnitLink\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"oracleWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"oracleWithdrawAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"_f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"_onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"_offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"_offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"donPublicKey\",\"type\":\"bytes\"}],\"name\":\"setDONPublicKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"thresholdPublicKey\",\"type\":\"bytes\"}],\"name\":\"setThresholdPublicKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"flags\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"requestingContract\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"availableBalance\",\"type\":\"uint96\"},{\"internalType\":\"uint72\",\"name\":\"adminFee\",\"type\":\"uint72\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"initiatedRequests\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"dataVersion\",\"type\":\"uint16\"},{\"internalType\":\"uint64\",\"name\":\"completedRequests\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"subscriptionOwner\",\"type\":\"address\"}],\"internalType\":\"structFunctionsResponse.RequestMeta\",\"name\":\"request\",\"type\":\"tuple\"}],\"name\":\"startRequest\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"estimatedTotalCostJuels\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint72\",\"name\":\"adminFee\",\"type\":\"uint72\"},{\"internalType\":\"uint72\",\"name\":\"donFee\",\"type\":\"uint72\"},{\"internalType\":\"uint40\",\"name\":\"gasOverheadBeforeCallback\",\"type\":\"uint40\"},{\"internalType\":\"uint40\",\"name\":\"gasOverheadAfterCallback\",\"type\":\"uint40\"},{\"internalType\":\"uint32\",\"name\":\"timeoutTimestamp\",\"type\":\"uint32\"}],\"internalType\":\"structFunctionsResponse.Commitment\",\"name\":\"commitment\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"transmitters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"fulfillmentGasPriceOverEstimationBP\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"feedStalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasOverheadBeforeCallback\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasOverheadAfterCallback\",\"type\":\"uint32\"},{\"internalType\":\"uint72\",\"name\":\"donFee\",\"type\":\"uint72\"},{\"internalType\":\"uint40\",\"name\":\"minimumEstimateGasPriceWei\",\"type\":\"uint40\"},{\"internalType\":\"uint16\",\"name\":\"maxSupportedRequestDataVersion\",\"type\":\"uint16\"},{\"internalType\":\"uint224\",\"name\":\"fallbackNativePerUnitLink\",\"type\":\"uint224\"},{\"internalType\":\"uint32\",\"name\":\"requestTimeoutSeconds\",\"type\":\"uint32\"}],\"internalType\":\"structFunctionsCoordinator.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"updateConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b50604051620057953803806200579583398101604081905262000034916200045d565b8233806000816200008c5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000bf57620000bf8162000133565b5050506001600160a01b038116620000ea57604051632530e88560e11b815260040160405180910390fd5b6001600160a01b03908116608052600b80549183166c01000000000000000000000000026001600160601b039092169190911790556200012a82620001de565b5050506200061c565b336001600160a01b038216036200018d5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000083565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b620001e86200033c565b80516008805460208401516040808601516060870151608088015160a089015160c08a015161ffff16600160f01b026001600160f01b0364ffffffffff909216600160c81b0264ffffffffff60c81b196001600160481b03909416600160801b0293909316600160801b600160f01b031963ffffffff9586166c010000000000000000000000000263ffffffff60601b19978716680100000000000000000297909716600160401b600160801b0319998716640100000000026001600160401b0319909b169c87169c909c1799909917979097169990991793909317959095169390931793909317929092169390931790915560e0830151610100840151909216600160e01b026001600160e01b0390921691909117600955517f5f32d06f5e83eda3a68e0e964ef2e6af5cb613e8117aa103c2d6bca5f5184862906200033190839062000566565b60405180910390a150565b6000546001600160a01b03163314620003985760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640162000083565b565b80516001600160a01b0381168114620003b257600080fd5b919050565b60405161012081016001600160401b0381118282101715620003e957634e487b7160e01b600052604160045260246000fd5b60405290565b805163ffffffff81168114620003b257600080fd5b80516001600160481b0381168114620003b257600080fd5b805164ffffffffff81168114620003b257600080fd5b805161ffff81168114620003b257600080fd5b80516001600160e01b0381168114620003b257600080fd5b60008060008385036101608112156200047557600080fd5b62000480856200039a565b935061012080601f19830112156200049757600080fd5b620004a1620003b7565b9150620004b160208701620003ef565b8252620004c160408701620003ef565b6020830152620004d460608701620003ef565b6040830152620004e760808701620003ef565b6060830152620004fa60a0870162000404565b60808301526200050d60c087016200041c565b60a08301526200052060e0870162000432565b60c08301526101006200053581880162000445565b60e084015262000547828801620003ef565b908301525091506200055d61014085016200039a565b90509250925092565b815163ffffffff908116825260208084015182169083015260408084015182169083015260608084015191821690830152610120820190506080830151620005b960808401826001600160481b03169052565b5060a0830151620005d360a084018264ffffffffff169052565b5060c0830151620005ea60c084018261ffff169052565b5060e08301516200060660e08401826001600160e01b03169052565b506101009283015163ffffffff16919092015290565b6080516151336200066260003960008181610845015281816109d301528181610d0601528181610f9a015281816110a5015281816117e9015261355d01526151336000f3fe608060405234801561001057600080fd5b506004361061018d5760003560e01c806381ff7048116100e3578063c3f909d41161008c578063e3d0e71211610066578063e3d0e71214610560578063e4ddcea614610573578063f2fde38b1461058957600080fd5b8063c3f909d4146103b0578063d227d24514610528578063d328a91e1461055857600080fd5b8063a631571e116100bd578063a631571e1461035d578063afcb95d71461037d578063b1dc65a41461039d57600080fd5b806381ff7048146102b557806385b214cf146103225780638da5cb5b1461033557600080fd5b806366316d8d116101455780637f15e1661161011f5780637f15e16614610285578063814118341461029857806381f1b938146102ad57600080fd5b806366316d8d1461026257806379ba5097146102755780637d4807871461027d57600080fd5b8063181f5a7711610176578063181f5a77146101ba5780632a905ccc1461020c57806359b5b7ac1461022e57600080fd5b8063083a5466146101925780631112dadc146101a7575b600080fd5b6101a56101a0366004613a9b565b61059c565b005b6101a56101b5366004613c44565b6105f1565b6101f66040518060400160405280601c81526020017f46756e6374696f6e7320436f6f7264696e61746f722076312e312e300000000081525081565b6040516102039190613d68565b60405180910390f35b610214610841565b60405168ffffffffffffffffff9091168152602001610203565b61021461023c366004613e09565b50600854700100000000000000000000000000000000900468ffffffffffffffffff1690565b6101a5610270366004613e98565b6108d7565b6101a5610a90565b6101a5610b92565b6101a5610293366004613a9b565b610df2565b6102a0610e42565b6040516102039190613f22565b6101f6610eb1565b6102ff60015460025463ffffffff74010000000000000000000000000000000000000000830481169378010000000000000000000000000000000000000000000000009093041691565b6040805163ffffffff948516815293909216602084015290820152606001610203565b6101a5610330366004613f35565b610f82565b60005460405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610203565b61037061036b366004613f4e565b611034565b60405161020391906140a3565b604080516001815260006020820181905291810191909152606001610203565b6101a56103ab3660046140f7565b6111d5565b61051b6040805161012081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e0810182905261010081019190915250604080516101208101825260085463ffffffff80821683526401000000008204811660208401526801000000000000000082048116938301939093526c01000000000000000000000000810483166060830152700100000000000000000000000000000000810468ffffffffffffffffff166080830152790100000000000000000000000000000000000000000000000000810464ffffffffff1660a08301527e01000000000000000000000000000000000000000000000000000000000000900461ffff1660c08201526009547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff811660e08301527c0100000000000000000000000000000000000000000000000000000000900490911661010082015290565b60405161020391906141ae565b61053b61053636600461429e565b6117e5565b6040516bffffffffffffffffffffffff9091168152602001610203565b6101f6611963565b6101a561056e3660046143b7565b6119ba565b61057b612536565b604051908152602001610203565b6101a5610597366004614484565b61278f565b6105a46127a3565b60008190036105df576040517f4f42be3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600d6105ec82848361453a565b505050565b6105f96127a3565b80516008805460208401516040808601516060870151608088015160a089015160c08a015161ffff167e01000000000000000000000000000000000000000000000000000000000000027dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff64ffffffffff909216790100000000000000000000000000000000000000000000000000027fffff0000000000ffffffffffffffffffffffffffffffffffffffffffffffffff68ffffffffffffffffff90941670010000000000000000000000000000000002939093167fffff0000000000000000000000000000ffffffffffffffffffffffffffffffff63ffffffff9586166c01000000000000000000000000027fffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffffff9787166801000000000000000002979097167fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff998716640100000000027fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000909b169c87169c909c1799909917979097169990991793909317959095169390931793909317929092169390931790915560e08301516101008401519092167c0100000000000000000000000000000000000000000000000000000000027bffffffffffffffffffffffffffffffffffffffffffffffffffffffff90921691909117600955517f5f32d06f5e83eda3a68e0e964ef2e6af5cb613e8117aa103c2d6bca5f5184862906108369083906141ae565b60405180910390a150565b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16632a905ccc6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156108ae573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108d29190614660565b905090565b6108df612826565b806bffffffffffffffffffffffff166000036109195750336000908152600a60205260409020546bffffffffffffffffffffffff16610973565b336000908152600a60205260409020546bffffffffffffffffffffffff80831691161015610973576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b336000908152600a6020526040812080548392906109a09084906bffffffffffffffffffffffff166146ac565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055506109f57f000000000000000000000000000000000000000000000000000000000000000090565b6040517f66316d8d00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff84811660048301526bffffffffffffffffffffffff8416602483015291909116906366316d8d90604401600060405180830381600087803b158015610a7457600080fd5b505af1158015610a88573d6000803e3d6000fd5b505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610b16576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610b9a6127a3565b610ba2612826565b60006006805480602002602001604051908101604052809291908181526020018280548015610c0757602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610bdc575b5050505050905060005b8151811015610dee576000600a6000848481518110610c3257610c326146d1565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252810191909152604001600020546bffffffffffffffffffffffff1690508015610ddd576000600a6000858581518110610c9157610c916146d1565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550610d287f000000000000000000000000000000000000000000000000000000000000000090565b73ffffffffffffffffffffffffffffffffffffffff166366316d8d848481518110610d5557610d556146d1565b6020026020010151836040518363ffffffff1660e01b8152600401610daa92919073ffffffffffffffffffffffffffffffffffffffff9290921682526bffffffffffffffffffffffff16602082015260400190565b600060405180830381600087803b158015610dc457600080fd5b505af1158015610dd8573d6000803e3d6000fd5b505050505b50610de781614700565b9050610c11565b5050565b610dfa6127a3565b6000819003610e35576040517f4f42be3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600c6105ec82848361453a565b60606006805480602002602001604051908101604052809291908181526020018280548015610ea757602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610e7c575b5050505050905090565b6060600d8054610ec0906144a1565b9050600003610efb576040517f4f42be3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600d8054610f08906144a1565b80601f0160208091040260200160405190810160405280929190818152602001828054610f34906144a1565b8015610ea75780601f10610f5657610100808354040283529160200191610ea7565b820191906000526020600020905b815481529060010190602001808311610f6457509395945050505050565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610ff1576040517fc41a5b0900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008181526007602052604080822091909155517f8a4b97add3359bd6bcf5e82874363670eb5ad0f7615abddbd0ed0a3a98f0f416906108369083815260200190565b6040805161016081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e08101829052610100810182905261012081018290526101408101919091523373ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016146110fc576040517fc41a5b0900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61110d61110883614738565b612a34565b905061111f6060830160408401614484565b815173ffffffffffffffffffffffffffffffffffffffff91909116907fbf50768ccf13bd0110ca6d53a9c4f1f3271abdd4c24a56878863ed25b20598ff3261116d60c0870160a08801614825565b61117f61016088016101408901614484565b6111898880614842565b61119b6101208b016101008c016148a7565b60208b01356111b16101008d0160e08e016148c2565b8b6040516111c7999897969594939291906148df565b60405180910390a35b919050565b60005a604080518b3580825262ffffff6020808f0135600881901c929092169084015293945092917fb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62910160405180910390a16112368a8a8a8a8a8a612efc565b6003546000906002906112549060ff80821691610100900416614995565b61125e91906149dd565b611269906001614995565b60ff1690508781146112d7576040517f660bd4ba00000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f77726f6e67206e756d626572206f66207369676e6174757265730000000000006044820152606401610b0d565b878614611366576040517f660bd4ba00000000000000000000000000000000000000000000000000000000815260206004820152602860248201527f7265706f727420727320616e64207373206d757374206265206f66206571756160448201527f6c206c656e6774680000000000000000000000000000000000000000000000006064820152608401610b0d565b3360009081526004602090815260408083208151808301909252805460ff808216845292939192918401916101009091041660028111156113a9576113a96149ff565b60028111156113ba576113ba6149ff565b90525090506002816020015160028111156113d7576113d76149ff565b1415801561142057506006816000015160ff16815481106113fa576113fa6146d1565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff163314155b15611487576040517f660bd4ba00000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f756e617574686f72697a6564207472616e736d697474657200000000000000006044820152606401610b0d565b50505050611493613a33565b6000808a8a6040516114a6929190614a2e565b6040519081900381206114bd918e90602001614a3e565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181528282528051602091820120838301909252600080845290830152915060005b898110156117c7576000600184898460208110611526576115266146d1565b61153391901a601b614995565b8e8e86818110611545576115456146d1565b905060200201358d8d8781811061155e5761155e6146d1565b905060200201356040516000815260200160405260405161159b949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa1580156115bd573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff811660009081526004602090815290849020838501909452835460ff8082168552929650929450840191610100900416600281111561163d5761163d6149ff565b600281111561164e5761164e6149ff565b905250925060018360200151600281111561166b5761166b6149ff565b146116d2576040517f660bd4ba00000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f61646472657373206e6f7420617574686f72697a656420746f207369676e00006044820152606401610b0d565b8251600090879060ff16601f81106116ec576116ec6146d1565b602002015173ffffffffffffffffffffffffffffffffffffffff161461176e576040517f660bd4ba00000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f6e6f6e2d756e69717565207369676e61747572650000000000000000000000006044820152606401610b0d565b8086846000015160ff16601f8110611788576117886146d1565b73ffffffffffffffffffffffffffffffffffffffff90921660209290920201526117b3600186614995565b945050806117c090614700565b9050611507565b5050506117d8833383858e8e612fb3565b5050505050505050505050565b60007f00000000000000000000000000000000000000000000000000000000000000006040517f10fc49c100000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8816600482015263ffffffff8516602482015273ffffffffffffffffffffffffffffffffffffffff91909116906310fc49c19060440160006040518083038186803b15801561188557600080fd5b505afa158015611899573d6000803e3d6000fd5b5050505066038d7ea4c680008211156118de576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6118e6610841565b68ffffffffffffffffff1661193086868080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061023c92505050565b68ffffffffffffffffff1661194585856131b2565b61194f9190614a52565b6119599190614a52565b9695505050505050565b6060600c8054611972906144a1565b90506000036119ad576040517f4f42be3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600c8054610f08906144a1565b855185518560ff16601f831115611a2d576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f746f6f206d616e79207369676e657273000000000000000000000000000000006044820152606401610b0d565b80600003611a97576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f66206d75737420626520706f73697469766500000000000000000000000000006044820152606401610b0d565b818314611b25576040517f89a61989000000000000000000000000000000000000000000000000000000008152602060048201526024808201527f6f7261636c6520616464726573736573206f7574206f6620726567697374726160448201527f74696f6e000000000000000000000000000000000000000000000000000000006064820152608401610b0d565b611b30816003614a77565b8311611b98576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f6661756c74792d6f7261636c65206620746f6f206869676800000000000000006044820152606401610b0d565b611ba06127a3565b6040805160c0810182528a8152602081018a905260ff89169181018290526060810188905267ffffffffffffffff8716608082015260a0810186905290611be790886132f3565b60055415611d9c57600554600090611c0190600190614a8e565b9050600060058281548110611c1857611c186146d1565b60009182526020822001546006805473ffffffffffffffffffffffffffffffffffffffff90921693509084908110611c5257611c526146d1565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff85811684526004909252604080842080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000090811690915592909116808452922080549091169055600580549192509080611cd257611cd2614aa1565b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff00000000000000000000000000000000000000001690550190556006805480611d3b57611d3b614aa1565b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff000000000000000000000000000000000000000016905501905550611be7915050565b60005b81515181101561235357815180516000919083908110611dc157611dc16146d1565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1603611e46576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f7369676e6572206d757374206e6f7420626520656d70747900000000000000006044820152606401610b0d565b600073ffffffffffffffffffffffffffffffffffffffff1682602001518281518110611e7457611e746146d1565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1603611ef9576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f7472616e736d6974746572206d757374206e6f7420626520656d7074790000006044820152606401610b0d565b60006004600084600001518481518110611f1557611f156146d1565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff16825281019190915260400160002054610100900460ff166002811115611f5f57611f5f6149ff565b14611fc6576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265706561746564207369676e657220616464726573730000000000000000006044820152606401610b0d565b6040805180820190915260ff82168152600160208201528251805160049160009185908110611ff757611ff76146d1565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040016000208251815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082168117835592840151919283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001617610100836002811115612098576120986149ff565b0217905550600091506120a89050565b60046000846020015184815181106120c2576120c26146d1565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff16825281019190915260400160002054610100900460ff16600281111561210c5761210c6149ff565b14612173576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f7265706561746564207472616e736d69747465722061646472657373000000006044820152606401610b0d565b6040805180820190915260ff8216815260208101600281525060046000846020015184815181106121a6576121a66146d1565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040016000208251815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082168117835592840151919283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001617610100836002811115612247576122476149ff565b021790555050825180516005925083908110612265576122656146d1565b602090810291909101810151825460018101845560009384529282902090920180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff90931692909217909155820151805160069190839081106122e1576122e16146d1565b60209081029190910181015182546001810184556000938452919092200180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff9092169190911790558061234b81614700565b915050611d9f565b506040810151600380547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff909216919091179055600180547fffffffff00000000ffffffffffffffffffffffffffffffffffffffffffffffff8116780100000000000000000000000000000000000000000000000063ffffffff438116820292909217808555920481169291829160149161240b91849174010000000000000000000000000000000000000000900416614ad0565b92506101000a81548163ffffffff021916908363ffffffff16021790555061246a4630600160149054906101000a900463ffffffff1663ffffffff16856000015186602001518760400151886060015189608001518a60a00151613303565b600281905582518051600380547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1661010060ff9093169290920291909117905560015460208501516040808701516060880151608089015160a08a015193517f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e0598612521988b9891977401000000000000000000000000000000000000000090920463ffffffff16969095919491939192614aed565b60405180910390a15050505050505050505050565b604080516101208101825260085463ffffffff80821683526401000000008204811660208401526801000000000000000082048116838501526c0100000000000000000000000080830482166060850152700100000000000000000000000000000000830468ffffffffffffffffff166080850152790100000000000000000000000000000000000000000000000000830464ffffffffff1660a0808601919091527e0100000000000000000000000000000000000000000000000000000000000090930461ffff1660c08501526009547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff811660e08601527c01000000000000000000000000000000000000000000000000000000009004909116610100840152600b5484517ffeaf968c00000000000000000000000000000000000000000000000000000000815294516000958694859490930473ffffffffffffffffffffffffffffffffffffffff169263feaf968c926004808401938290030181865afa1580156126c4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906126e89190614b9d565b5093505092505080426126fb9190614a8e565b836020015163ffffffff1610801561271d57506000836020015163ffffffff16115b1561274b57505060e001517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16919050565b60008213612788576040517f43d4cf6600000000000000000000000000000000000000000000000000000000815260048101839052602401610b0d565b5092915050565b6127976127a3565b6127a0816133ae565b50565b60005473ffffffffffffffffffffffffffffffffffffffff163314612824576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610b0d565b565b600b546bffffffffffffffffffffffff1660000361284057565b600060068054806020026020016040519081016040528092919081815260200182805480156128a557602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff16815260019091019060200180831161287a575b50505050509050600081519050806000036128ec576040517f30274b3a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600b5460009061290b9083906bffffffffffffffffffffffff16614bed565b905060005b828110156129d65781600a600086848151811061292f5761292f6146d1565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282829054906101000a90046bffffffffffffffffffffffff166129979190614a52565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550806129cf90614700565b9050612910565b506129e18282614c18565b600b8054600090612a019084906bffffffffffffffffffffffff166146ac565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550505050565b6040805161016081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e0810182905261010081018290526101208101829052610140810191909152604080516101208101825260085463ffffffff80821683526401000000008204811660208401526801000000000000000082048116938301939093526c0100000000000000000000000081048316606083015268ffffffffffffffffff700100000000000000000000000000000000820416608083015264ffffffffff79010000000000000000000000000000000000000000000000000082041660a083015261ffff7e01000000000000000000000000000000000000000000000000000000000000909104811660c083018190526009547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff811660e08501527c0100000000000000000000000000000000000000000000000000000000900490931661010080840191909152850151919291161115612bef576040517fdada758700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600854600090700100000000000000000000000000000000900468ffffffffffffffffff1690506000846080015168ffffffffffffffffff168268ffffffffffffffffff16612c428760e001513a6131b2565b612c4c9190614a52565b612c569190614a52565b9050806bffffffffffffffffffffffff1685606001516bffffffffffffffffffffffff161015612cb2576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600083610100015163ffffffff1642612ccb9190614c40565b905060003087604001518860a001518960c001516001612ceb9190614c53565b8a5180516020918201206101008d015160e08e0151604051612d9f98979695948c918c9132910173ffffffffffffffffffffffffffffffffffffffff9a8b168152988a1660208a015267ffffffffffffffff97881660408a0152959096166060880152608087019390935261ffff9190911660a086015263ffffffff90811660c08601526bffffffffffffffffffffffff9190911660e0850152919091166101008301529091166101208201526101400190565b6040516020818303038152906040528051906020012090506040518061016001604052808281526020013073ffffffffffffffffffffffffffffffffffffffff168152602001846bffffffffffffffffffffffff168152602001886040015173ffffffffffffffffffffffffffffffffffffffff1681526020018860a0015167ffffffffffffffff1681526020018860e0015163ffffffff168152602001886080015168ffffffffffffffffff1681526020018568ffffffffffffffffff168152602001866040015163ffffffff1664ffffffffff168152602001866060015163ffffffff1664ffffffffff1681526020018363ffffffff16815250955085604051602001612eae91906140a3565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152918152815160209283012060009384526007909252909120555092949350505050565b6000612f09826020614a77565b612f14856020614a77565b612f2088610144614c40565b612f2a9190614c40565b612f349190614c40565b612f3f906000614c40565b9050368114612faa576040517f660bd4ba00000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f63616c6c64617461206c656e677468206d69736d6174636800000000000000006044820152606401610b0d565b50505050505050565b600080808080612fc586880188614d4f565b84519499509297509095509350915060ff16801580612fe5575084518114155b80612ff1575083518114155b80612ffd575082518114155b80613009575081518114155b15613070576040517f660bd4ba00000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f4669656c6473206d75737420626520657175616c206c656e67746800000000006044820152606401610b0d565b60005b818110156131a3576000613108888381518110613092576130926146d1565b60200260200101518884815181106130ac576130ac6146d1565b60200260200101518885815181106130c6576130c66146d1565b60200260200101518886815181106130e0576130e06146d1565b60200260200101518887815181106130fa576130fa6146d1565b6020026020010151886134a3565b9050600081600681111561311e5761311e6149ff565b148061313b57506001816006811115613139576131396149ff565b145b1561319257878281518110613152576131526146d1565b60209081029190910181015160405133815290917fc708e0440951fd63499c0f7a73819b469ee5dd3ecc356c0ab4eb7f18389009d9910160405180910390a25b5061319c81614700565b9050613073565b50505050505050505050505050565b600854600090790100000000000000000000000000000000000000000000000000900464ffffffffff1682101561320d57600854790100000000000000000000000000000000000000000000000000900464ffffffffff1691505b600854600090612710906132279063ffffffff1685614a77565b6132319190614e21565b61323b9084614c40565b60085490915060009085906132749063ffffffff6c01000000000000000000000000820481169168010000000000000000900416614ad0565b61327e9190614ad0565b63ffffffff16905060006132c86000368080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506137b792505050565b90506132e7816132d88486614a77565b6132e29190614c40565b6138f9565b93505050505b92915050565b60065415610dee57610dee612826565b6000808a8a8a8a8a8a8a8a8a60405160200161332799989796959493929190614e35565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179150509998505050505050505050565b3373ffffffffffffffffffffffffffffffffffffffff82160361342d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610b0d565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600080848060200190518101906134ba9190614f01565b905060003a8261012001518361010001516134d59190614fc9565b64ffffffffff166134e69190614a77565b905060008460ff1661352e6000368080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506137b792505050565b6135389190614e21565b905060006135496132e28385614c40565b905060006135563a6138f9565b90506000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663330605298e8e868b60e0015168ffffffffffffffffff16896135b59190614a52565b338d6040518763ffffffff1660e01b81526004016135d896959493929190614fe7565b60408051808303816000875af11580156135f6573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061361a9190615063565b90925090506000826006811115613633576136336149ff565b14806136505750600182600681111561364e5761364e6149ff565b145b156137a65760008e81526007602052604081205561366e8185614a52565b336000908152600a6020526040812080547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff93841617905560e0890151600b805468ffffffffffffffffff909216939092916136da91859116614a52565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055508d7f90815c2e624694e8010bffad2bcefaf96af282ef1bc2ebc0042d1b89a585e0468487848b60c0015168ffffffffffffffffff168c60e0015168ffffffffffffffffff16878b6137599190614a52565b6137639190614a52565b61376d9190614a52565b604080516bffffffffffffffffffffffff9586168152602081019490945291841683830152909216606082015290519081900360800190a25b509c9b505050505050505050505050565b6000466137c381613927565b1561383f57606c73ffffffffffffffffffffffffffffffffffffffff1663c6f7de0e6040518163ffffffff1660e01b8152600401602060405180830381865afa158015613814573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906138389190615096565b9392505050565b6138488161394a565b156138f05773420000000000000000000000000000000000000f73ffffffffffffffffffffffffffffffffffffffff166349948e0e846040518060800160405280604881526020016150df604891396040516020016138a89291906150af565b6040516020818303038152906040526040518263ffffffff1660e01b81526004016138d39190613d68565b602060405180830381865afa158015613814573d6000803e3d6000fd5b50600092915050565b60006132ed613906612536565b61391884670de0b6b3a7640000614a77565b6139229190614e21565b613991565b600061a4b182148061393b575062066eed82145b806132ed57505062066eee1490565b6000600a82148061395c57506101a482145b80613969575062aa37dc82145b80613975575061210582145b80613982575062014a3382145b806132ed57505062014a341490565b60006bffffffffffffffffffffffff821115613a2f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f53616665436173743a2076616c756520646f65736e27742066697420696e203960448201527f36206269747300000000000000000000000000000000000000000000000000006064820152608401610b0d565b5090565b604051806103e00160405280601f906020820280368337509192915050565b60008083601f840112613a6457600080fd5b50813567ffffffffffffffff811115613a7c57600080fd5b602083019150836020828501011115613a9457600080fd5b9250929050565b60008060208385031215613aae57600080fd5b823567ffffffffffffffff811115613ac557600080fd5b613ad185828601613a52565b90969095509350505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610120810167ffffffffffffffff81118282101715613b3057613b30613add565b60405290565b604051610160810167ffffffffffffffff81118282101715613b3057613b30613add565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715613ba157613ba1613add565b604052919050565b63ffffffff811681146127a057600080fd5b80356111d081613ba9565b68ffffffffffffffffff811681146127a057600080fd5b80356111d081613bc6565b64ffffffffff811681146127a057600080fd5b80356111d081613be8565b803561ffff811681146111d057600080fd5b80357bffffffffffffffffffffffffffffffffffffffffffffffffffffffff811681146111d057600080fd5b60006101208284031215613c5757600080fd5b613c5f613b0c565b613c6883613bbb565b8152613c7660208401613bbb565b6020820152613c8760408401613bbb565b6040820152613c9860608401613bbb565b6060820152613ca960808401613bdd565b6080820152613cba60a08401613bfb565b60a0820152613ccb60c08401613c06565b60c0820152613cdc60e08401613c18565b60e0820152610100613cef818501613bbb565b908201529392505050565b60005b83811015613d15578181015183820152602001613cfd565b50506000910152565b60008151808452613d36816020860160208601613cfa565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6020815260006138386020830184613d1e565b600082601f830112613d8c57600080fd5b813567ffffffffffffffff811115613da657613da6613add565b613dd760207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601613b5a565b818152846020838601011115613dec57600080fd5b816020850160208301376000918101602001919091529392505050565b600060208284031215613e1b57600080fd5b813567ffffffffffffffff811115613e3257600080fd5b613e3e84828501613d7b565b949350505050565b73ffffffffffffffffffffffffffffffffffffffff811681146127a057600080fd5b80356111d081613e46565b6bffffffffffffffffffffffff811681146127a057600080fd5b80356111d081613e73565b60008060408385031215613eab57600080fd5b8235613eb681613e46565b91506020830135613ec681613e73565b809150509250929050565b600081518084526020808501945080840160005b83811015613f1757815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101613ee5565b509495945050505050565b6020815260006138386020830184613ed1565b600060208284031215613f4757600080fd5b5035919050565b600060208284031215613f6057600080fd5b813567ffffffffffffffff811115613f7757600080fd5b8201610160818503121561383857600080fd5b805182526020810151613fb5602084018273ffffffffffffffffffffffffffffffffffffffff169052565b506040810151613fd560408401826bffffffffffffffffffffffff169052565b506060810151613ffd606084018273ffffffffffffffffffffffffffffffffffffffff169052565b506080810151614019608084018267ffffffffffffffff169052565b5060a081015161403160a084018263ffffffff169052565b5060c081015161404e60c084018268ffffffffffffffffff169052565b5060e081015161406b60e084018268ffffffffffffffffff169052565b506101008181015164ffffffffff9081169184019190915261012080830151909116908301526101409081015163ffffffff16910152565b61016081016132ed8284613f8a565b60008083601f8401126140c457600080fd5b50813567ffffffffffffffff8111156140dc57600080fd5b6020830191508360208260051b8501011115613a9457600080fd5b60008060008060008060008060e0898b03121561411357600080fd5b606089018a81111561412457600080fd5b8998503567ffffffffffffffff8082111561413e57600080fd5b61414a8c838d01613a52565b909950975060808b013591508082111561416357600080fd5b61416f8c838d016140b2565b909750955060a08b013591508082111561418857600080fd5b506141958b828c016140b2565b999c989b50969995989497949560c00135949350505050565b815163ffffffff908116825260208084015182169083015260408084015182169083015260608084015191821690830152610120820190506080830151614202608084018268ffffffffffffffffff169052565b5060a083015161421b60a084018264ffffffffff169052565b5060c083015161423160c084018261ffff169052565b5060e083015161426160e08401827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff169052565b506101008381015163ffffffff8116848301525b505092915050565b67ffffffffffffffff811681146127a057600080fd5b80356111d08161427d565b6000806000806000608086880312156142b657600080fd5b85356142c18161427d565b9450602086013567ffffffffffffffff8111156142dd57600080fd5b6142e988828901613a52565b90955093505060408601356142fd81613ba9565b949793965091946060013592915050565b600067ffffffffffffffff82111561432857614328613add565b5060051b60200190565b600082601f83011261434357600080fd5b813560206143586143538361430e565b613b5a565b82815260059290921b8401810191818101908684111561437757600080fd5b8286015b8481101561439b57803561438e81613e46565b835291830191830161437b565b509695505050505050565b803560ff811681146111d057600080fd5b60008060008060008060c087890312156143d057600080fd5b863567ffffffffffffffff808211156143e857600080fd5b6143f48a838b01614332565b9750602089013591508082111561440a57600080fd5b6144168a838b01614332565b965061442460408a016143a6565b9550606089013591508082111561443a57600080fd5b6144468a838b01613d7b565b945061445460808a01614293565b935060a089013591508082111561446a57600080fd5b5061447789828a01613d7b565b9150509295509295509295565b60006020828403121561449657600080fd5b813561383881613e46565b600181811c908216806144b557607f821691505b6020821081036144ee577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f8211156105ec57600081815260208120601f850160051c8101602086101561451b5750805b601f850160051c820191505b81811015610a8857828155600101614527565b67ffffffffffffffff83111561455257614552613add565b6145668361456083546144a1565b836144f4565b6000601f8411600181146145b857600085156145825750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b17835561464e565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b8281101561460757868501358255602094850194600190920191016145e7565b5086821015614642577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b80516111d081613bc6565b60006020828403121561467257600080fd5b815161383881613bc6565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6bffffffffffffffffffffffff8281168282160390808211156127885761278861467d565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036147315761473161467d565b5060010190565b6000610160823603121561474b57600080fd5b614753613b36565b823567ffffffffffffffff81111561476a57600080fd5b61477636828601613d7b565b8252506020830135602082015261478f60408401613e68565b60408201526147a060608401613e8d565b60608201526147b160808401613bdd565b60808201526147c260a08401614293565b60a08201526147d360c08401614293565b60c08201526147e460e08401613bbb565b60e08201526101006147f7818501613c06565b90820152610120614809848201614293565b9082015261014061481b848201613e68565b9082015292915050565b60006020828403121561483757600080fd5b81356138388161427d565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261487757600080fd5b83018035915067ffffffffffffffff82111561489257600080fd5b602001915036819003821315613a9457600080fd5b6000602082840312156148b957600080fd5b61383882613c06565b6000602082840312156148d457600080fd5b813561383881613ba9565b73ffffffffffffffffffffffffffffffffffffffff8a8116825267ffffffffffffffff8a166020830152881660408201526102406060820181905281018690526000610260878982850137600083890182015261ffff8716608084015260a0830186905263ffffffff851660c0840152601f88017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016830101905061498760e0830184613f8a565b9a9950505050505050505050565b60ff81811683821601908111156132ed576132ed61467d565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b600060ff8316806149f0576149f06149ae565b8060ff84160491505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b8183823760009101908152919050565b828152606082602083013760800192915050565b6bffffffffffffffffffffffff8181168382160190808211156127885761278861467d565b80820281158282048414176132ed576132ed61467d565b818103818111156132ed576132ed61467d565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b63ffffffff8181168382160190808211156127885761278861467d565b600061012063ffffffff808d1684528b6020850152808b16604085015250806060840152614b1d8184018a613ed1565b90508281036080840152614b318189613ed1565b905060ff871660a084015282810360c0840152614b4e8187613d1e565b905067ffffffffffffffff851660e0840152828103610100840152614b738185613d1e565b9c9b505050505050505050505050565b805169ffffffffffffffffffff811681146111d057600080fd5b600080600080600060a08688031215614bb557600080fd5b614bbe86614b83565b9450602086015193506040860151925060608601519150614be160808701614b83565b90509295509295909350565b60006bffffffffffffffffffffffff80841680614c0c57614c0c6149ae565b92169190910492915050565b6bffffffffffffffffffffffff8181168382160280821691908281146142755761427561467d565b808201808211156132ed576132ed61467d565b67ffffffffffffffff8181168382160190808211156127885761278861467d565b600082601f830112614c8557600080fd5b81356020614c956143538361430e565b82815260059290921b84018101918181019086841115614cb457600080fd5b8286015b8481101561439b5780358352918301918301614cb8565b600082601f830112614ce057600080fd5b81356020614cf06143538361430e565b82815260059290921b84018101918181019086841115614d0f57600080fd5b8286015b8481101561439b57803567ffffffffffffffff811115614d335760008081fd5b614d418986838b0101613d7b565b845250918301918301614d13565b600080600080600060a08688031215614d6757600080fd5b853567ffffffffffffffff80821115614d7f57600080fd5b614d8b89838a01614c74565b96506020880135915080821115614da157600080fd5b614dad89838a01614ccf565b95506040880135915080821115614dc357600080fd5b614dcf89838a01614ccf565b94506060880135915080821115614de557600080fd5b614df189838a01614ccf565b93506080880135915080821115614e0757600080fd5b50614e1488828901614ccf565b9150509295509295909350565b600082614e3057614e306149ae565b500490565b60006101208b835273ffffffffffffffffffffffffffffffffffffffff8b16602084015267ffffffffffffffff808b166040850152816060850152614e7c8285018b613ed1565b91508382036080850152614e90828a613ed1565b915060ff881660a085015283820360c0850152614ead8288613d1e565b90861660e08501528381036101008501529050614b738185613d1e565b80516111d081613e46565b80516111d081613e73565b80516111d08161427d565b80516111d081613ba9565b80516111d081613be8565b60006101608284031215614f1457600080fd5b614f1c613b36565b82518152614f2c60208401614eca565b6020820152614f3d60408401614ed5565b6040820152614f4e60608401614eca565b6060820152614f5f60808401614ee0565b6080820152614f7060a08401614eeb565b60a0820152614f8160c08401614655565b60c0820152614f9260e08401614655565b60e0820152610100614fa5818501614ef6565b90820152610120614fb7848201614ef6565b90820152610140613cef848201614eeb565b64ffffffffff8181168382160190808211156127885761278861467d565b6000610200808352614ffb8184018a613d1e565b9050828103602084015261500f8189613d1e565b6bffffffffffffffffffffffff88811660408601528716606085015273ffffffffffffffffffffffffffffffffffffffff861660808501529150615058905060a0830184613f8a565b979650505050505050565b6000806040838503121561507657600080fd5b82516007811061508557600080fd5b6020840151909250613ec681613e73565b6000602082840312156150a857600080fd5b5051919050565b600083516150c1818460208801613cfa565b8351908301906150d5818360208801613cfa565b0194935050505056fe307866666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666a164736f6c6343000813000a",
}

var FunctionsCoordinatorABI = FunctionsCoordinatorMetaData.ABI

var FunctionsCoordinatorBin = FunctionsCoordinatorMetaData.Bin

func DeployFunctionsCoordinator(auth *bind.TransactOpts, backend bind.ContractBackend, router common.Address, config FunctionsCoordinatorConfig, linkToNativeFeed common.Address) (common.Address, *types.Transaction, *FunctionsCoordinator, error) {
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

func (_FunctionsCoordinator *FunctionsCoordinatorCaller) GetAdminFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FunctionsCoordinator.contract.Call(opts, &out, "getAdminFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_FunctionsCoordinator *FunctionsCoordinatorSession) GetAdminFee() (*big.Int, error) {
	return _FunctionsCoordinator.Contract.GetAdminFee(&_FunctionsCoordinator.CallOpts)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCallerSession) GetAdminFee() (*big.Int, error) {
	return _FunctionsCoordinator.Contract.GetAdminFee(&_FunctionsCoordinator.CallOpts)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCaller) GetConfig(opts *bind.CallOpts) (FunctionsCoordinatorConfig, error) {
	var out []interface{}
	err := _FunctionsCoordinator.contract.Call(opts, &out, "getConfig")

	if err != nil {
		return *new(FunctionsCoordinatorConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(FunctionsCoordinatorConfig)).(*FunctionsCoordinatorConfig)

	return out0, err

}

func (_FunctionsCoordinator *FunctionsCoordinatorSession) GetConfig() (FunctionsCoordinatorConfig, error) {
	return _FunctionsCoordinator.Contract.GetConfig(&_FunctionsCoordinator.CallOpts)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCallerSession) GetConfig() (FunctionsCoordinatorConfig, error) {
	return _FunctionsCoordinator.Contract.GetConfig(&_FunctionsCoordinator.CallOpts)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCaller) GetDONFee(opts *bind.CallOpts, arg0 []byte) (*big.Int, error) {
	var out []interface{}
	err := _FunctionsCoordinator.contract.Call(opts, &out, "getDONFee", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_FunctionsCoordinator *FunctionsCoordinatorSession) GetDONFee(arg0 []byte) (*big.Int, error) {
	return _FunctionsCoordinator.Contract.GetDONFee(&_FunctionsCoordinator.CallOpts, arg0)
}

func (_FunctionsCoordinator *FunctionsCoordinatorCallerSession) GetDONFee(arg0 []byte) (*big.Int, error) {
	return _FunctionsCoordinator.Contract.GetDONFee(&_FunctionsCoordinator.CallOpts, arg0)
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

func (_FunctionsCoordinator *FunctionsCoordinatorTransactor) UpdateConfig(opts *bind.TransactOpts, config FunctionsCoordinatorConfig) (*types.Transaction, error) {
	return _FunctionsCoordinator.contract.Transact(opts, "updateConfig", config)
}

func (_FunctionsCoordinator *FunctionsCoordinatorSession) UpdateConfig(config FunctionsCoordinatorConfig) (*types.Transaction, error) {
	return _FunctionsCoordinator.Contract.UpdateConfig(&_FunctionsCoordinator.TransactOpts, config)
}

func (_FunctionsCoordinator *FunctionsCoordinatorTransactorSession) UpdateConfig(config FunctionsCoordinatorConfig) (*types.Transaction, error) {
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
	Config FunctionsCoordinatorConfig
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
	TotalCostJuels    *big.Int
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
	return common.HexToHash("0x5f32d06f5e83eda3a68e0e964ef2e6af5cb613e8117aa103c2d6bca5f5184862")
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
	return common.HexToHash("0x90815c2e624694e8010bffad2bcefaf96af282ef1bc2ebc0042d1b89a585e046")
}

func (FunctionsCoordinatorTransmitted) Topic() common.Hash {
	return common.HexToHash("0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62")
}

func (_FunctionsCoordinator *FunctionsCoordinator) Address() common.Address {
	return _FunctionsCoordinator.address
}

type FunctionsCoordinatorInterface interface {
	EstimateCost(opts *bind.CallOpts, subscriptionId uint64, data []byte, callbackGasLimit uint32, gasPriceWei *big.Int) (*big.Int, error)

	GetAdminFee(opts *bind.CallOpts) (*big.Int, error)

	GetConfig(opts *bind.CallOpts) (FunctionsCoordinatorConfig, error)

	GetDONFee(opts *bind.CallOpts, arg0 []byte) (*big.Int, error)

	GetDONPublicKey(opts *bind.CallOpts) ([]byte, error)

	GetThresholdPublicKey(opts *bind.CallOpts) ([]byte, error)

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

	UpdateConfig(opts *bind.TransactOpts, config FunctionsCoordinatorConfig) (*types.Transaction, error)

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
