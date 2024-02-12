// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package functions_coordinator_1_1_0

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

var FunctionsCoordinator110MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"fulfillmentGasPriceOverEstimationBP\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"feedStalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasOverheadBeforeCallback\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasOverheadAfterCallback\",\"type\":\"uint32\"},{\"internalType\":\"uint72\",\"name\":\"donFee\",\"type\":\"uint72\"},{\"internalType\":\"uint40\",\"name\":\"minimumEstimateGasPriceWei\",\"type\":\"uint40\"},{\"internalType\":\"uint16\",\"name\":\"maxSupportedRequestDataVersion\",\"type\":\"uint16\"},{\"internalType\":\"uint224\",\"name\":\"fallbackNativePerUnitLink\",\"type\":\"uint224\"},{\"internalType\":\"uint32\",\"name\":\"requestTimeoutSeconds\",\"type\":\"uint32\"}],\"internalType\":\"structFunctionsBilling.Config\",\"name\":\"config\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"linkToNativeFeed\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"EmptyPublicKey\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InconsistentReportData\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"InvalidConfig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"linkWei\",\"type\":\"int256\"}],\"name\":\"InvalidLinkWeiPrice\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSubscription\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"MustBeSubOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoTransmittersSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByRouter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByRouterOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentTooLarge\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"ReportInvalid\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RouterMustBeSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedPublicKeyChange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedSender\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnsupportedRequestDataVersion\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"CommitmentDeleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"fulfillmentGasPriceOverEstimationBP\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"feedStalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasOverheadBeforeCallback\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasOverheadAfterCallback\",\"type\":\"uint32\"},{\"internalType\":\"uint72\",\"name\":\"donFee\",\"type\":\"uint72\"},{\"internalType\":\"uint40\",\"name\":\"minimumEstimateGasPriceWei\",\"type\":\"uint40\"},{\"internalType\":\"uint16\",\"name\":\"maxSupportedRequestDataVersion\",\"type\":\"uint16\"},{\"internalType\":\"uint224\",\"name\":\"fallbackNativePerUnitLink\",\"type\":\"uint224\"},{\"internalType\":\"uint32\",\"name\":\"requestTimeoutSeconds\",\"type\":\"uint32\"}],\"indexed\":false,\"internalType\":\"structFunctionsBilling.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"ConfigUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requestingContract\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requestInitiator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"subscriptionOwner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"dataVersion\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"flags\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"callbackGasLimit\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"estimatedTotalCostJuels\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint72\",\"name\":\"adminFee\",\"type\":\"uint72\"},{\"internalType\":\"uint72\",\"name\":\"donFee\",\"type\":\"uint72\"},{\"internalType\":\"uint40\",\"name\":\"gasOverheadBeforeCallback\",\"type\":\"uint40\"},{\"internalType\":\"uint40\",\"name\":\"gasOverheadAfterCallback\",\"type\":\"uint40\"},{\"internalType\":\"uint32\",\"name\":\"timeoutTimestamp\",\"type\":\"uint32\"}],\"indexed\":false,\"internalType\":\"structFunctionsResponse.Commitment\",\"name\":\"commitment\",\"type\":\"tuple\"}],\"name\":\"OracleRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"OracleResponse\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"juelsPerGas\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"l1FeeShareWei\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"callbackCostJuels\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalCostJuels\",\"type\":\"uint96\"}],\"name\":\"RequestBilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"deleteCommitment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasPriceWei\",\"type\":\"uint256\"}],\"name\":\"estimateCost\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAdminFee\",\"outputs\":[{\"internalType\":\"uint72\",\"name\":\"\",\"type\":\"uint72\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"fulfillmentGasPriceOverEstimationBP\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"feedStalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasOverheadBeforeCallback\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasOverheadAfterCallback\",\"type\":\"uint32\"},{\"internalType\":\"uint72\",\"name\":\"donFee\",\"type\":\"uint72\"},{\"internalType\":\"uint40\",\"name\":\"minimumEstimateGasPriceWei\",\"type\":\"uint40\"},{\"internalType\":\"uint16\",\"name\":\"maxSupportedRequestDataVersion\",\"type\":\"uint16\"},{\"internalType\":\"uint224\",\"name\":\"fallbackNativePerUnitLink\",\"type\":\"uint224\"},{\"internalType\":\"uint32\",\"name\":\"requestTimeoutSeconds\",\"type\":\"uint32\"}],\"internalType\":\"structFunctionsBilling.Config\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"getDONFee\",\"outputs\":[{\"internalType\":\"uint72\",\"name\":\"\",\"type\":\"uint72\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDONPublicKey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getThresholdPublicKey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getWeiPerUnitLink\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"oracleWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"oracleWithdrawAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"_f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"_onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"_offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"_offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"donPublicKey\",\"type\":\"bytes\"}],\"name\":\"setDONPublicKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"thresholdPublicKey\",\"type\":\"bytes\"}],\"name\":\"setThresholdPublicKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"flags\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"requestingContract\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"availableBalance\",\"type\":\"uint96\"},{\"internalType\":\"uint72\",\"name\":\"adminFee\",\"type\":\"uint72\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"initiatedRequests\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"dataVersion\",\"type\":\"uint16\"},{\"internalType\":\"uint64\",\"name\":\"completedRequests\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"subscriptionOwner\",\"type\":\"address\"}],\"internalType\":\"structFunctionsResponse.RequestMeta\",\"name\":\"request\",\"type\":\"tuple\"}],\"name\":\"startRequest\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"estimatedTotalCostJuels\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint72\",\"name\":\"adminFee\",\"type\":\"uint72\"},{\"internalType\":\"uint72\",\"name\":\"donFee\",\"type\":\"uint72\"},{\"internalType\":\"uint40\",\"name\":\"gasOverheadBeforeCallback\",\"type\":\"uint40\"},{\"internalType\":\"uint40\",\"name\":\"gasOverheadAfterCallback\",\"type\":\"uint40\"},{\"internalType\":\"uint32\",\"name\":\"timeoutTimestamp\",\"type\":\"uint32\"}],\"internalType\":\"structFunctionsResponse.Commitment\",\"name\":\"commitment\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"transmitters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"fulfillmentGasPriceOverEstimationBP\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"feedStalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasOverheadBeforeCallback\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasOverheadAfterCallback\",\"type\":\"uint32\"},{\"internalType\":\"uint72\",\"name\":\"donFee\",\"type\":\"uint72\"},{\"internalType\":\"uint40\",\"name\":\"minimumEstimateGasPriceWei\",\"type\":\"uint40\"},{\"internalType\":\"uint16\",\"name\":\"maxSupportedRequestDataVersion\",\"type\":\"uint16\"},{\"internalType\":\"uint224\",\"name\":\"fallbackNativePerUnitLink\",\"type\":\"uint224\"},{\"internalType\":\"uint32\",\"name\":\"requestTimeoutSeconds\",\"type\":\"uint32\"}],\"internalType\":\"structFunctionsBilling.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"updateConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b50604051620056f7380380620056f783398101604081905262000034916200046d565b8282828233806000816200008f5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c257620000c28162000139565b5050506001600160a01b038116620000ed57604051632530e88560e11b815260040160405180910390fd5b6001600160a01b03908116608052600b80549183166c01000000000000000000000000026001600160601b039092169190911790556200012d82620001e4565b5050505050506200062c565b336001600160a01b03821603620001935760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000086565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b620001ee62000342565b80516008805460208401516040808601516060870151608088015160a089015160c08a015161ffff16600160f01b026001600160f01b0364ffffffffff909216600160c81b0264ffffffffff60c81b196001600160481b03909416600160801b0293909316600160801b600160f01b031963ffffffff9586166c010000000000000000000000000263ffffffff60601b19978716680100000000000000000297909716600160401b600160801b0319998716640100000000026001600160401b0319909b169c87169c909c1799909917979097169990991793909317959095169390931793909317929092169390931790915560e0830151610100840151909216600160e01b026001600160e01b0390921691909117600955517f5f32d06f5e83eda3a68e0e964ef2e6af5cb613e8117aa103c2d6bca5f5184862906200033790839062000576565b60405180910390a150565b6200034c6200034e565b565b6000546001600160a01b031633146200034c5760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640162000086565b80516001600160a01b0381168114620003c257600080fd5b919050565b60405161012081016001600160401b0381118282101715620003f957634e487b7160e01b600052604160045260246000fd5b60405290565b805163ffffffff81168114620003c257600080fd5b80516001600160481b0381168114620003c257600080fd5b805164ffffffffff81168114620003c257600080fd5b805161ffff81168114620003c257600080fd5b80516001600160e01b0381168114620003c257600080fd5b60008060008385036101608112156200048557600080fd5b6200049085620003aa565b935061012080601f1983011215620004a757600080fd5b620004b1620003c7565b9150620004c160208701620003ff565b8252620004d160408701620003ff565b6020830152620004e460608701620003ff565b6040830152620004f760808701620003ff565b60608301526200050a60a0870162000414565b60808301526200051d60c087016200042c565b60a08301526200053060e0870162000442565b60c08301526101006200054581880162000455565b60e084015262000557828801620003ff565b908301525091506200056d6101408501620003aa565b90509250925092565b815163ffffffff908116825260208084015182169083015260408084015182169083015260608084015191821690830152610120820190506080830151620005c960808401826001600160481b03169052565b5060a0830151620005e360a084018264ffffffffff169052565b5060c0830151620005fa60c084018261ffff169052565b5060e08301516200061660e08401826001600160e01b03169052565b506101009283015163ffffffff16919092015290565b6080516150856200067260003960008181610845015281816109d301528181610ca601528181610f3a0152818161104501528181611789015261349001526150856000f3fe608060405234801561001057600080fd5b506004361061018d5760003560e01c806381ff7048116100e3578063c3f909d41161008c578063e3d0e71211610066578063e3d0e71214610560578063e4ddcea614610573578063f2fde38b1461058957600080fd5b8063c3f909d4146103b0578063d227d24514610528578063d328a91e1461055857600080fd5b8063a631571e116100bd578063a631571e1461035d578063afcb95d71461037d578063b1dc65a41461039d57600080fd5b806381ff7048146102b557806385b214cf146103225780638da5cb5b1461033557600080fd5b806366316d8d116101455780637f15e1661161011f5780637f15e16614610285578063814118341461029857806381f1b938146102ad57600080fd5b806366316d8d1461026257806379ba5097146102755780637d4807871461027d57600080fd5b8063181f5a7711610176578063181f5a77146101ba5780632a905ccc1461020c57806359b5b7ac1461022e57600080fd5b8063083a5466146101925780631112dadc146101a7575b600080fd5b6101a56101a03660046139fb565b61059c565b005b6101a56101b5366004613ba4565b6105f1565b6101f66040518060400160405280601c81526020017f46756e6374696f6e7320436f6f7264696e61746f722076312e312e300000000081525081565b6040516102039190613cc8565b60405180910390f35b610214610841565b60405168ffffffffffffffffff9091168152602001610203565b61021461023c366004613d69565b50600854700100000000000000000000000000000000900468ffffffffffffffffff1690565b6101a5610270366004613df8565b6108d7565b6101a5610a90565b6101a5610b92565b6101a56102933660046139fb565b610d92565b6102a0610de2565b6040516102039190613e82565b6101f6610e51565b6102ff60015460025463ffffffff74010000000000000000000000000000000000000000830481169378010000000000000000000000000000000000000000000000009093041691565b6040805163ffffffff948516815293909216602084015290820152606001610203565b6101a5610330366004613e95565b610f22565b60005460405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610203565b61037061036b366004613eae565b610fd4565b6040516102039190614003565b604080516001815260006020820181905291810191909152606001610203565b6101a56103ab366004614057565b611175565b61051b6040805161012081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e0810182905261010081019190915250604080516101208101825260085463ffffffff80821683526401000000008204811660208401526801000000000000000082048116938301939093526c01000000000000000000000000810483166060830152700100000000000000000000000000000000810468ffffffffffffffffff166080830152790100000000000000000000000000000000000000000000000000810464ffffffffff1660a08301527e01000000000000000000000000000000000000000000000000000000000000900461ffff1660c08201526009547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff811660e08301527c0100000000000000000000000000000000000000000000000000000000900490911661010082015290565b604051610203919061410e565b61053b6105363660046141fe565b611785565b6040516bffffffffffffffffffffffff9091168152602001610203565b6101f66118e5565b6101a561056e366004614317565b61193c565b61057b6124b8565b604051908152602001610203565b6101a56105973660046143e4565b612711565b6105a4612725565b60008190036105df576040517f4f42be3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600d6105ec82848361449a565b505050565b6105f96127a8565b80516008805460208401516040808601516060870151608088015160a089015160c08a015161ffff167e01000000000000000000000000000000000000000000000000000000000000027dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff64ffffffffff909216790100000000000000000000000000000000000000000000000000027fffff0000000000ffffffffffffffffffffffffffffffffffffffffffffffffff68ffffffffffffffffff90941670010000000000000000000000000000000002939093167fffff0000000000000000000000000000ffffffffffffffffffffffffffffffff63ffffffff9586166c01000000000000000000000000027fffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffffff9787166801000000000000000002979097167fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff998716640100000000027fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000909b169c87169c909c1799909917979097169990991793909317959095169390931793909317929092169390931790915560e08301516101008401519092167c0100000000000000000000000000000000000000000000000000000000027bffffffffffffffffffffffffffffffffffffffffffffffffffffffff90921691909117600955517f5f32d06f5e83eda3a68e0e964ef2e6af5cb613e8117aa103c2d6bca5f51848629061083690839061410e565b60405180910390a150565b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16632a905ccc6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156108ae573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108d291906145c0565b905090565b6108df6127b0565b806bffffffffffffffffffffffff166000036109195750336000908152600a60205260409020546bffffffffffffffffffffffff16610973565b336000908152600a60205260409020546bffffffffffffffffffffffff80831691161015610973576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b336000908152600a6020526040812080548392906109a09084906bffffffffffffffffffffffff1661460c565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055506109f57f000000000000000000000000000000000000000000000000000000000000000090565b6040517f66316d8d00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff84811660048301526bffffffffffffffffffffffff8416602483015291909116906366316d8d90604401600060405180830381600087803b158015610a7457600080fd5b505af1158015610a88573d6000803e3d6000fd5b505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610b16576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610b9a6127a8565b610ba26127b0565b6000610bac610de2565b905060005b8151811015610d8e576000600a6000848481518110610bd257610bd2614631565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252810191909152604001600020546bffffffffffffffffffffffff1690508015610d7d576000600a6000858581518110610c3157610c31614631565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550610cc87f000000000000000000000000000000000000000000000000000000000000000090565b73ffffffffffffffffffffffffffffffffffffffff166366316d8d848481518110610cf557610cf5614631565b6020026020010151836040518363ffffffff1660e01b8152600401610d4a92919073ffffffffffffffffffffffffffffffffffffffff9290921682526bffffffffffffffffffffffff16602082015260400190565b600060405180830381600087803b158015610d6457600080fd5b505af1158015610d78573d6000803e3d6000fd5b505050505b50610d8781614660565b9050610bb1565b5050565b610d9a612725565b6000819003610dd5576040517f4f42be3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600c6105ec82848361449a565b60606006805480602002602001604051908101604052809291908181526020018280548015610e4757602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610e1c575b5050505050905090565b6060600d8054610e6090614401565b9050600003610e9b576040517f4f42be3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600d8054610ea890614401565b80601f0160208091040260200160405190810160405280929190818152602001828054610ed490614401565b8015610e475780601f10610ef657610100808354040283529160200191610e47565b820191906000526020600020905b815481529060010190602001808311610f0457509395945050505050565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610f91576040517fc41a5b0900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008181526007602052604080822091909155517f8a4b97add3359bd6bcf5e82874363670eb5ad0f7615abddbd0ed0a3a98f0f416906108369083815260200190565b6040805161016081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e08101829052610100810182905261012081018290526101408101919091523373ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161461109c576040517fc41a5b0900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6110ad6110a883614698565b61295c565b90506110bf60608301604084016143e4565b815173ffffffffffffffffffffffffffffffffffffffff91909116907fbf50768ccf13bd0110ca6d53a9c4f1f3271abdd4c24a56878863ed25b20598ff3261110d60c0870160a08801614785565b61111f610160880161014089016143e4565b61112988806147a2565b61113b6101208b016101008c01614807565b60208b01356111516101008d0160e08e01614822565b8b6040516111679998979695949392919061483f565b60405180910390a35b919050565b60005a604080518b3580825262ffffff6020808f0135600881901c929092169084015293945092917fb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62910160405180910390a16111d68a8a8a8a8a8a612dfa565b6003546000906002906111f49060ff808216916101009004166148e7565b6111fe919061492f565b6112099060016148e7565b60ff169050878114611277576040517f660bd4ba00000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f77726f6e67206e756d626572206f66207369676e6174757265730000000000006044820152606401610b0d565b878614611306576040517f660bd4ba00000000000000000000000000000000000000000000000000000000815260206004820152602860248201527f7265706f727420727320616e64207373206d757374206265206f66206571756160448201527f6c206c656e6774680000000000000000000000000000000000000000000000006064820152608401610b0d565b3360009081526004602090815260408083208151808301909252805460ff8082168452929391929184019161010090910416600281111561134957611349614951565b600281111561135a5761135a614951565b905250905060028160200151600281111561137757611377614951565b141580156113c057506006816000015160ff168154811061139a5761139a614631565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff163314155b15611427576040517f660bd4ba00000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f756e617574686f72697a6564207472616e736d697474657200000000000000006044820152606401610b0d565b50505050611433613993565b6000808a8a604051611446929190614980565b60405190819003812061145d918e90602001614990565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181528282528051602091820120838301909252600080845290830152915060005b898110156117675760006001848984602081106114c6576114c6614631565b6114d391901a601b6148e7565b8e8e868181106114e5576114e5614631565b905060200201358d8d878181106114fe576114fe614631565b905060200201356040516000815260200160405260405161153b949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa15801561155d573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff811660009081526004602090815290849020838501909452835460ff808216855292965092945084019161010090041660028111156115dd576115dd614951565b60028111156115ee576115ee614951565b905250925060018360200151600281111561160b5761160b614951565b14611672576040517f660bd4ba00000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f61646472657373206e6f7420617574686f72697a656420746f207369676e00006044820152606401610b0d565b8251600090879060ff16601f811061168c5761168c614631565b602002015173ffffffffffffffffffffffffffffffffffffffff161461170e576040517f660bd4ba00000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f6e6f6e2d756e69717565207369676e61747572650000000000000000000000006044820152606401610b0d565b8086846000015160ff16601f811061172857611728614631565b73ffffffffffffffffffffffffffffffffffffffff90921660209290920201526117536001866148e7565b9450508061176090614660565b90506114a7565b505050611778833383858e8e612eb1565b5050505050505050505050565b60007f00000000000000000000000000000000000000000000000000000000000000006040517f10fc49c100000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8816600482015263ffffffff8516602482015273ffffffffffffffffffffffffffffffffffffffff91909116906310fc49c19060440160006040518083038186803b15801561182557600080fd5b505afa158015611839573d6000803e3d6000fd5b5050505066038d7ea4c6800082111561187e576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000611888610841565b905060006118cb87878080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061023c92505050565b90506118d9858583856130b0565b98975050505050505050565b6060600c80546118f490614401565b905060000361192f576040517f4f42be3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600c8054610ea890614401565b855185518560ff16601f8311156119af576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f746f6f206d616e79207369676e657273000000000000000000000000000000006044820152606401610b0d565b80600003611a19576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f66206d75737420626520706f73697469766500000000000000000000000000006044820152606401610b0d565b818314611aa7576040517f89a61989000000000000000000000000000000000000000000000000000000008152602060048201526024808201527f6f7261636c6520616464726573736573206f7574206f6620726567697374726160448201527f74696f6e000000000000000000000000000000000000000000000000000000006064820152608401610b0d565b611ab28160036149a4565b8311611b1a576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f6661756c74792d6f7261636c65206620746f6f206869676800000000000000006044820152606401610b0d565b611b22612725565b6040805160c0810182528a8152602081018a905260ff89169181018290526060810188905267ffffffffffffffff8716608082015260a0810186905290611b69908861321d565b60055415611d1e57600554600090611b83906001906149bb565b9050600060058281548110611b9a57611b9a614631565b60009182526020822001546006805473ffffffffffffffffffffffffffffffffffffffff90921693509084908110611bd457611bd4614631565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff85811684526004909252604080842080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000090811690915592909116808452922080549091169055600580549192509080611c5457611c546149ce565b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff00000000000000000000000000000000000000001690550190556006805480611cbd57611cbd6149ce565b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff000000000000000000000000000000000000000016905501905550611b69915050565b60005b8151518110156122d557815180516000919083908110611d4357611d43614631565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1603611dc8576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f7369676e6572206d757374206e6f7420626520656d70747900000000000000006044820152606401610b0d565b600073ffffffffffffffffffffffffffffffffffffffff1682602001518281518110611df657611df6614631565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1603611e7b576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f7472616e736d6974746572206d757374206e6f7420626520656d7074790000006044820152606401610b0d565b60006004600084600001518481518110611e9757611e97614631565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff16825281019190915260400160002054610100900460ff166002811115611ee157611ee1614951565b14611f48576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265706561746564207369676e657220616464726573730000000000000000006044820152606401610b0d565b6040805180820190915260ff82168152600160208201528251805160049160009185908110611f7957611f79614631565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040016000208251815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082168117835592840151919283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000161761010083600281111561201a5761201a614951565b02179055506000915061202a9050565b600460008460200151848151811061204457612044614631565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff16825281019190915260400160002054610100900460ff16600281111561208e5761208e614951565b146120f5576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f7265706561746564207472616e736d69747465722061646472657373000000006044820152606401610b0d565b6040805180820190915260ff82168152602081016002815250600460008460200151848151811061212857612128614631565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040016000208251815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082168117835592840151919283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000016176101008360028111156121c9576121c9614951565b0217905550508251805160059250839081106121e7576121e7614631565b602090810291909101810151825460018101845560009384529282902090920180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909316929092179091558201518051600691908390811061226357612263614631565b60209081029190910181015182546001810184556000938452919092200180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909216919091179055806122cd81614660565b915050611d21565b506040810151600380547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff909216919091179055600180547fffffffff00000000ffffffffffffffffffffffffffffffffffffffffffffffff8116780100000000000000000000000000000000000000000000000063ffffffff438116820292909217808555920481169291829160149161238d918491740100000000000000000000000000000000000000009004166149fd565b92506101000a81548163ffffffff021916908363ffffffff1602179055506123ec4630600160149054906101000a900463ffffffff1663ffffffff16856000015186602001518760400151886060015189608001518a60a00151613236565b600281905582518051600380547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1661010060ff9093169290920291909117905560015460208501516040808701516060880151608089015160a08a015193517f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05986124a3988b9891977401000000000000000000000000000000000000000090920463ffffffff16969095919491939192614a1a565b60405180910390a15050505050505050505050565b604080516101208101825260085463ffffffff80821683526401000000008204811660208401526801000000000000000082048116838501526c0100000000000000000000000080830482166060850152700100000000000000000000000000000000830468ffffffffffffffffff166080850152790100000000000000000000000000000000000000000000000000830464ffffffffff1660a0808601919091527e0100000000000000000000000000000000000000000000000000000000000090930461ffff1660c08501526009547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff811660e08601527c01000000000000000000000000000000000000000000000000000000009004909116610100840152600b5484517ffeaf968c00000000000000000000000000000000000000000000000000000000815294516000958694859490930473ffffffffffffffffffffffffffffffffffffffff169263feaf968c926004808401938290030181865afa158015612646573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061266a9190614aca565b50935050925050804261267d91906149bb565b836020015163ffffffff1610801561269f57506000836020015163ffffffff16115b156126cd57505060e001517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16919050565b6000821361270a576040517f43d4cf6600000000000000000000000000000000000000000000000000000000815260048101839052602401610b0d565b5092915050565b612719612725565b612722816132e1565b50565b60005473ffffffffffffffffffffffffffffffffffffffff1633146127a6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610b0d565b565b6127a6612725565b600b546bffffffffffffffffffffffff166000036127ca57565b60006127d4610de2565b80519091506000819003612814576040517f30274b3a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600b546000906128339083906bffffffffffffffffffffffff16614b1a565b905060005b828110156128fe5781600a600086848151811061285757612857614631565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282829054906101000a90046bffffffffffffffffffffffff166128bf9190614b45565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550806128f790614660565b9050612838565b506129098282614b6a565b600b80546000906129299084906bffffffffffffffffffffffff1661460c565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550505050565b6040805161016081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e0810182905261010081018290526101208101829052610140810191909152604080516101208101825260085463ffffffff80821683526401000000008204811660208401526801000000000000000082048116938301939093526c0100000000000000000000000081048316606083015268ffffffffffffffffff700100000000000000000000000000000000820416608083015264ffffffffff79010000000000000000000000000000000000000000000000000082041660a083015261ffff7e01000000000000000000000000000000000000000000000000000000000000909104811660c083018190526009547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff811660e08501527c0100000000000000000000000000000000000000000000000000000000900490931661010080840191909152850151919291161115612b17576040517fdada758700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600854600090700100000000000000000000000000000000900468ffffffffffffffffff1690506000612b548560e001513a8488608001516130b0565b9050806bffffffffffffffffffffffff1685606001516bffffffffffffffffffffffff161015612bb0576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600083610100015163ffffffff1642612bc99190614b92565b905060003087604001518860a001518960c001516001612be99190614ba5565b8a5180516020918201206101008d015160e08e0151604051612c9d98979695948c918c9132910173ffffffffffffffffffffffffffffffffffffffff9a8b168152988a1660208a015267ffffffffffffffff97881660408a0152959096166060880152608087019390935261ffff9190911660a086015263ffffffff90811660c08601526bffffffffffffffffffffffff9190911660e0850152919091166101008301529091166101208201526101400190565b6040516020818303038152906040528051906020012090506040518061016001604052808281526020013073ffffffffffffffffffffffffffffffffffffffff168152602001846bffffffffffffffffffffffff168152602001886040015173ffffffffffffffffffffffffffffffffffffffff1681526020018860a0015167ffffffffffffffff1681526020018860e0015163ffffffff168152602001886080015168ffffffffffffffffff1681526020018568ffffffffffffffffff168152602001866040015163ffffffff1664ffffffffff168152602001866060015163ffffffff1664ffffffffff1681526020018363ffffffff16815250955085604051602001612dac9190614003565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152918152815160209283012060009384526007909252909120555092949350505050565b6000612e078260206149a4565b612e128560206149a4565b612e1e88610144614b92565b612e289190614b92565b612e329190614b92565b612e3d906000614b92565b9050368114612ea8576040517f660bd4ba00000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f63616c6c64617461206c656e677468206d69736d6174636800000000000000006044820152606401610b0d565b50505050505050565b600080808080612ec386880188614ca1565b84519499509297509095509350915060ff16801580612ee3575084518114155b80612eef575083518114155b80612efb575082518114155b80612f07575081518114155b15612f6e576040517f660bd4ba00000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f4669656c6473206d75737420626520657175616c206c656e67746800000000006044820152606401610b0d565b60005b818110156130a1576000613006888381518110612f9057612f90614631565b6020026020010151888481518110612faa57612faa614631565b6020026020010151888581518110612fc457612fc4614631565b6020026020010151888681518110612fde57612fde614631565b6020026020010151888781518110612ff857612ff8614631565b6020026020010151886133d6565b9050600081600681111561301c5761301c614951565b14806130395750600181600681111561303757613037614951565b145b156130905787828151811061305057613050614631565b60209081029190910181015160405133815290917fc708e0440951fd63499c0f7a73819b469ee5dd3ecc356c0ab4eb7f18389009d9910160405180910390a25b5061309a81614660565b9050612f71565b50505050505050505050505050565b600854600090790100000000000000000000000000000000000000000000000000900464ffffffffff1684101561310b57600854790100000000000000000000000000000000000000000000000000900464ffffffffff1693505b600854600090612710906131259063ffffffff16876149a4565b61312f9190614d73565b6131399086614b92565b60085490915060009087906131729063ffffffff6c010000000000000000000000008204811691680100000000000000009004166149fd565b61317c91906149fd565b63ffffffff16905060006131c66000368080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061371192505050565b905060006131e7826131d885876149a4565b6131e29190614b92565b613853565b9050600061320368ffffffffffffffffff808916908a16614b45565b905061320f8183614b45565b9a9950505050505050505050565b6000613227610de2565b511115610d8e57610d8e6127b0565b6000808a8a8a8a8a8a8a8a8a60405160200161325a99989796959493929190614d87565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179150509998505050505050505050565b3373ffffffffffffffffffffffffffffffffffffffff821603613360576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610b0d565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600080848060200190518101906133ed9190614e53565b905060003a8261012001518361010001516134089190614f1b565b64ffffffffff1661341991906149a4565b905060008460ff166134616000368080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061371192505050565b61346b9190614d73565b9050600061347c6131e28385614b92565b905060006134893a613853565b90506000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663330605298e8e868b60e0015168ffffffffffffffffff16896134e89190614b45565b338d6040518763ffffffff1660e01b815260040161350b96959493929190614f39565b60408051808303816000875af1158015613529573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061354d9190614fb5565b9092509050600082600681111561356657613566614951565b14806135835750600182600681111561358157613581614951565b145b156137005760008e8152600760205260408120556135a18185614b45565b336000908152600a6020526040812080549091906135ce9084906bffffffffffffffffffffffff16614b45565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055508660e0015168ffffffffffffffffff16600b60008282829054906101000a90046bffffffffffffffffffffffff166136349190614b45565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055508d7f90815c2e624694e8010bffad2bcefaf96af282ef1bc2ebc0042d1b89a585e0468487848b60c0015168ffffffffffffffffff168c60e0015168ffffffffffffffffff16878b6136b39190614b45565b6136bd9190614b45565b6136c79190614b45565b604080516bffffffffffffffffffffffff9586168152602081019490945291841683830152909216606082015290519081900360800190a25b509c9b505050505050505050505050565b60004661371d81613887565b1561379957606c73ffffffffffffffffffffffffffffffffffffffff1663c6f7de0e6040518163ffffffff1660e01b8152600401602060405180830381865afa15801561376e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906137929190614fe8565b9392505050565b6137a2816138aa565b1561384a5773420000000000000000000000000000000000000f73ffffffffffffffffffffffffffffffffffffffff166349948e0e8460405180608001604052806048815260200161503160489139604051602001613802929190615001565b6040516020818303038152906040526040518263ffffffff1660e01b815260040161382d9190613cc8565b602060405180830381865afa15801561376e573d6000803e3d6000fd5b50600092915050565b60006138816138606124b8565b61387284670de0b6b3a76400006149a4565b61387c9190614d73565b6138f1565b92915050565b600061a4b182148061389b575062066eed82145b8061388157505062066eee1490565b6000600a8214806138bc57506101a482145b806138c9575062aa37dc82145b806138d5575061210582145b806138e2575062014a3382145b8061388157505062014a341490565b60006bffffffffffffffffffffffff82111561398f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f53616665436173743a2076616c756520646f65736e27742066697420696e203960448201527f36206269747300000000000000000000000000000000000000000000000000006064820152608401610b0d565b5090565b604051806103e00160405280601f906020820280368337509192915050565b60008083601f8401126139c457600080fd5b50813567ffffffffffffffff8111156139dc57600080fd5b6020830191508360208285010111156139f457600080fd5b9250929050565b60008060208385031215613a0e57600080fd5b823567ffffffffffffffff811115613a2557600080fd5b613a31858286016139b2565b90969095509350505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610120810167ffffffffffffffff81118282101715613a9057613a90613a3d565b60405290565b604051610160810167ffffffffffffffff81118282101715613a9057613a90613a3d565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715613b0157613b01613a3d565b604052919050565b63ffffffff8116811461272257600080fd5b803561117081613b09565b68ffffffffffffffffff8116811461272257600080fd5b803561117081613b26565b64ffffffffff8116811461272257600080fd5b803561117081613b48565b803561ffff8116811461117057600080fd5b80357bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8116811461117057600080fd5b60006101208284031215613bb757600080fd5b613bbf613a6c565b613bc883613b1b565b8152613bd660208401613b1b565b6020820152613be760408401613b1b565b6040820152613bf860608401613b1b565b6060820152613c0960808401613b3d565b6080820152613c1a60a08401613b5b565b60a0820152613c2b60c08401613b66565b60c0820152613c3c60e08401613b78565b60e0820152610100613c4f818501613b1b565b908201529392505050565b60005b83811015613c75578181015183820152602001613c5d565b50506000910152565b60008151808452613c96816020860160208601613c5a565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6020815260006137926020830184613c7e565b600082601f830112613cec57600080fd5b813567ffffffffffffffff811115613d0657613d06613a3d565b613d3760207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601613aba565b818152846020838601011115613d4c57600080fd5b816020850160208301376000918101602001919091529392505050565b600060208284031215613d7b57600080fd5b813567ffffffffffffffff811115613d9257600080fd5b613d9e84828501613cdb565b949350505050565b73ffffffffffffffffffffffffffffffffffffffff8116811461272257600080fd5b803561117081613da6565b6bffffffffffffffffffffffff8116811461272257600080fd5b803561117081613dd3565b60008060408385031215613e0b57600080fd5b8235613e1681613da6565b91506020830135613e2681613dd3565b809150509250929050565b600081518084526020808501945080840160005b83811015613e7757815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101613e45565b509495945050505050565b6020815260006137926020830184613e31565b600060208284031215613ea757600080fd5b5035919050565b600060208284031215613ec057600080fd5b813567ffffffffffffffff811115613ed757600080fd5b8201610160818503121561379257600080fd5b805182526020810151613f15602084018273ffffffffffffffffffffffffffffffffffffffff169052565b506040810151613f3560408401826bffffffffffffffffffffffff169052565b506060810151613f5d606084018273ffffffffffffffffffffffffffffffffffffffff169052565b506080810151613f79608084018267ffffffffffffffff169052565b5060a0810151613f9160a084018263ffffffff169052565b5060c0810151613fae60c084018268ffffffffffffffffff169052565b5060e0810151613fcb60e084018268ffffffffffffffffff169052565b506101008181015164ffffffffff9081169184019190915261012080830151909116908301526101409081015163ffffffff16910152565b61016081016138818284613eea565b60008083601f84011261402457600080fd5b50813567ffffffffffffffff81111561403c57600080fd5b6020830191508360208260051b85010111156139f457600080fd5b60008060008060008060008060e0898b03121561407357600080fd5b606089018a81111561408457600080fd5b8998503567ffffffffffffffff8082111561409e57600080fd5b6140aa8c838d016139b2565b909950975060808b01359150808211156140c357600080fd5b6140cf8c838d01614012565b909750955060a08b01359150808211156140e857600080fd5b506140f58b828c01614012565b999c989b50969995989497949560c00135949350505050565b815163ffffffff908116825260208084015182169083015260408084015182169083015260608084015191821690830152610120820190506080830151614162608084018268ffffffffffffffffff169052565b5060a083015161417b60a084018264ffffffffff169052565b5060c083015161419160c084018261ffff169052565b5060e08301516141c160e08401827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff169052565b506101008381015163ffffffff8116848301525b505092915050565b67ffffffffffffffff8116811461272257600080fd5b8035611170816141dd565b60008060008060006080868803121561421657600080fd5b8535614221816141dd565b9450602086013567ffffffffffffffff81111561423d57600080fd5b614249888289016139b2565b909550935050604086013561425d81613b09565b949793965091946060013592915050565b600067ffffffffffffffff82111561428857614288613a3d565b5060051b60200190565b600082601f8301126142a357600080fd5b813560206142b86142b38361426e565b613aba565b82815260059290921b840181019181810190868411156142d757600080fd5b8286015b848110156142fb5780356142ee81613da6565b83529183019183016142db565b509695505050505050565b803560ff8116811461117057600080fd5b60008060008060008060c0878903121561433057600080fd5b863567ffffffffffffffff8082111561434857600080fd5b6143548a838b01614292565b9750602089013591508082111561436a57600080fd5b6143768a838b01614292565b965061438460408a01614306565b9550606089013591508082111561439a57600080fd5b6143a68a838b01613cdb565b94506143b460808a016141f3565b935060a08901359150808211156143ca57600080fd5b506143d789828a01613cdb565b9150509295509295509295565b6000602082840312156143f657600080fd5b813561379281613da6565b600181811c9082168061441557607f821691505b60208210810361444e577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f8211156105ec57600081815260208120601f850160051c8101602086101561447b5750805b601f850160051c820191505b81811015610a8857828155600101614487565b67ffffffffffffffff8311156144b2576144b2613a3d565b6144c6836144c08354614401565b83614454565b6000601f84116001811461451857600085156144e25750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b1783556145ae565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b828110156145675786850135825560209485019460019092019101614547565b50868210156145a2577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b805161117081613b26565b6000602082840312156145d257600080fd5b815161379281613b26565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6bffffffffffffffffffffffff82811682821603908082111561270a5761270a6145dd565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203614691576146916145dd565b5060010190565b600061016082360312156146ab57600080fd5b6146b3613a96565b823567ffffffffffffffff8111156146ca57600080fd5b6146d636828601613cdb565b825250602083013560208201526146ef60408401613dc8565b604082015261470060608401613ded565b606082015261471160808401613b3d565b608082015261472260a084016141f3565b60a082015261473360c084016141f3565b60c082015261474460e08401613b1b565b60e0820152610100614757818501613b66565b908201526101206147698482016141f3565b9082015261014061477b848201613dc8565b9082015292915050565b60006020828403121561479757600080fd5b8135613792816141dd565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18436030181126147d757600080fd5b83018035915067ffffffffffffffff8211156147f257600080fd5b6020019150368190038213156139f457600080fd5b60006020828403121561481957600080fd5b61379282613b66565b60006020828403121561483457600080fd5b813561379281613b09565b73ffffffffffffffffffffffffffffffffffffffff8a8116825267ffffffffffffffff8a166020830152881660408201526102406060820181905281018690526000610260878982850137600083890182015261ffff8716608084015260a0830186905263ffffffff851660c0840152601f88017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016830101905061320f60e0830184613eea565b60ff8181168382160190811115613881576138816145dd565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b600060ff83168061494257614942614900565b8060ff84160491505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b8183823760009101908152919050565b828152606082602083013760800192915050565b8082028115828204841417613881576138816145dd565b81810381811115613881576138816145dd565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b63ffffffff81811683821601908082111561270a5761270a6145dd565b600061012063ffffffff808d1684528b6020850152808b16604085015250806060840152614a4a8184018a613e31565b90508281036080840152614a5e8189613e31565b905060ff871660a084015282810360c0840152614a7b8187613c7e565b905067ffffffffffffffff851660e0840152828103610100840152614aa08185613c7e565b9c9b505050505050505050505050565b805169ffffffffffffffffffff8116811461117057600080fd5b600080600080600060a08688031215614ae257600080fd5b614aeb86614ab0565b9450602086015193506040860151925060608601519150614b0e60808701614ab0565b90509295509295909350565b60006bffffffffffffffffffffffff80841680614b3957614b39614900565b92169190910492915050565b6bffffffffffffffffffffffff81811683821601908082111561270a5761270a6145dd565b6bffffffffffffffffffffffff8181168382160280821691908281146141d5576141d56145dd565b80820180821115613881576138816145dd565b67ffffffffffffffff81811683821601908082111561270a5761270a6145dd565b600082601f830112614bd757600080fd5b81356020614be76142b38361426e565b82815260059290921b84018101918181019086841115614c0657600080fd5b8286015b848110156142fb5780358352918301918301614c0a565b600082601f830112614c3257600080fd5b81356020614c426142b38361426e565b82815260059290921b84018101918181019086841115614c6157600080fd5b8286015b848110156142fb57803567ffffffffffffffff811115614c855760008081fd5b614c938986838b0101613cdb565b845250918301918301614c65565b600080600080600060a08688031215614cb957600080fd5b853567ffffffffffffffff80821115614cd157600080fd5b614cdd89838a01614bc6565b96506020880135915080821115614cf357600080fd5b614cff89838a01614c21565b95506040880135915080821115614d1557600080fd5b614d2189838a01614c21565b94506060880135915080821115614d3757600080fd5b614d4389838a01614c21565b93506080880135915080821115614d5957600080fd5b50614d6688828901614c21565b9150509295509295909350565b600082614d8257614d82614900565b500490565b60006101208b835273ffffffffffffffffffffffffffffffffffffffff8b16602084015267ffffffffffffffff808b166040850152816060850152614dce8285018b613e31565b91508382036080850152614de2828a613e31565b915060ff881660a085015283820360c0850152614dff8288613c7e565b90861660e08501528381036101008501529050614aa08185613c7e565b805161117081613da6565b805161117081613dd3565b8051611170816141dd565b805161117081613b09565b805161117081613b48565b60006101608284031215614e6657600080fd5b614e6e613a96565b82518152614e7e60208401614e1c565b6020820152614e8f60408401614e27565b6040820152614ea060608401614e1c565b6060820152614eb160808401614e32565b6080820152614ec260a08401614e3d565b60a0820152614ed360c084016145b5565b60c0820152614ee460e084016145b5565b60e0820152610100614ef7818501614e48565b90820152610120614f09848201614e48565b90820152610140613c4f848201614e3d565b64ffffffffff81811683821601908082111561270a5761270a6145dd565b6000610200808352614f4d8184018a613c7e565b90508281036020840152614f618189613c7e565b6bffffffffffffffffffffffff88811660408601528716606085015273ffffffffffffffffffffffffffffffffffffffff861660808501529150614faa905060a0830184613eea565b979650505050505050565b60008060408385031215614fc857600080fd5b825160078110614fd757600080fd5b6020840151909250613e2681613dd3565b600060208284031215614ffa57600080fd5b5051919050565b60008351615013818460208801613c5a565b835190830190615027818360208801613c5a565b0194935050505056fe307866666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666a164736f6c6343000813000a",
}

var FunctionsCoordinator110ABI = FunctionsCoordinator110MetaData.ABI

var FunctionsCoordinator110Bin = FunctionsCoordinator110MetaData.Bin

func DeployFunctionsCoordinator110(auth *bind.TransactOpts, backend bind.ContractBackend, router common.Address, config FunctionsBillingConfig, linkToNativeFeed common.Address) (common.Address, *types.Transaction, *FunctionsCoordinator110, error) {
	parsed, err := FunctionsCoordinator110MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(FunctionsCoordinator110Bin), backend, router, config, linkToNativeFeed)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &FunctionsCoordinator110{address: address, abi: *parsed, FunctionsCoordinator110Caller: FunctionsCoordinator110Caller{contract: contract}, FunctionsCoordinator110Transactor: FunctionsCoordinator110Transactor{contract: contract}, FunctionsCoordinator110Filterer: FunctionsCoordinator110Filterer{contract: contract}}, nil
}

type FunctionsCoordinator110 struct {
	address common.Address
	abi     abi.ABI
	FunctionsCoordinator110Caller
	FunctionsCoordinator110Transactor
	FunctionsCoordinator110Filterer
}

type FunctionsCoordinator110Caller struct {
	contract *bind.BoundContract
}

type FunctionsCoordinator110Transactor struct {
	contract *bind.BoundContract
}

type FunctionsCoordinator110Filterer struct {
	contract *bind.BoundContract
}

type FunctionsCoordinator110Session struct {
	Contract     *FunctionsCoordinator110
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type FunctionsCoordinator110CallerSession struct {
	Contract *FunctionsCoordinator110Caller
	CallOpts bind.CallOpts
}

type FunctionsCoordinator110TransactorSession struct {
	Contract     *FunctionsCoordinator110Transactor
	TransactOpts bind.TransactOpts
}

type FunctionsCoordinator110Raw struct {
	Contract *FunctionsCoordinator110
}

type FunctionsCoordinator110CallerRaw struct {
	Contract *FunctionsCoordinator110Caller
}

type FunctionsCoordinator110TransactorRaw struct {
	Contract *FunctionsCoordinator110Transactor
}

func NewFunctionsCoordinator110(address common.Address, backend bind.ContractBackend) (*FunctionsCoordinator110, error) {
	abi, err := abi.JSON(strings.NewReader(FunctionsCoordinator110ABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindFunctionsCoordinator110(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FunctionsCoordinator110{address: address, abi: abi, FunctionsCoordinator110Caller: FunctionsCoordinator110Caller{contract: contract}, FunctionsCoordinator110Transactor: FunctionsCoordinator110Transactor{contract: contract}, FunctionsCoordinator110Filterer: FunctionsCoordinator110Filterer{contract: contract}}, nil
}

func NewFunctionsCoordinator110Caller(address common.Address, caller bind.ContractCaller) (*FunctionsCoordinator110Caller, error) {
	contract, err := bindFunctionsCoordinator110(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FunctionsCoordinator110Caller{contract: contract}, nil
}

func NewFunctionsCoordinator110Transactor(address common.Address, transactor bind.ContractTransactor) (*FunctionsCoordinator110Transactor, error) {
	contract, err := bindFunctionsCoordinator110(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FunctionsCoordinator110Transactor{contract: contract}, nil
}

func NewFunctionsCoordinator110Filterer(address common.Address, filterer bind.ContractFilterer) (*FunctionsCoordinator110Filterer, error) {
	contract, err := bindFunctionsCoordinator110(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FunctionsCoordinator110Filterer{contract: contract}, nil
}

func bindFunctionsCoordinator110(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FunctionsCoordinator110MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FunctionsCoordinator110.Contract.FunctionsCoordinator110Caller.contract.Call(opts, result, method, params...)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsCoordinator110.Contract.FunctionsCoordinator110Transactor.contract.Transfer(opts)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FunctionsCoordinator110.Contract.FunctionsCoordinator110Transactor.contract.Transact(opts, method, params...)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FunctionsCoordinator110.Contract.contract.Call(opts, result, method, params...)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsCoordinator110.Contract.contract.Transfer(opts)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FunctionsCoordinator110.Contract.contract.Transact(opts, method, params...)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Caller) EstimateCost(opts *bind.CallOpts, subscriptionId uint64, data []byte, callbackGasLimit uint32, gasPriceWei *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _FunctionsCoordinator110.contract.Call(opts, &out, "estimateCost", subscriptionId, data, callbackGasLimit, gasPriceWei)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Session) EstimateCost(subscriptionId uint64, data []byte, callbackGasLimit uint32, gasPriceWei *big.Int) (*big.Int, error) {
	return _FunctionsCoordinator110.Contract.EstimateCost(&_FunctionsCoordinator110.CallOpts, subscriptionId, data, callbackGasLimit, gasPriceWei)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110CallerSession) EstimateCost(subscriptionId uint64, data []byte, callbackGasLimit uint32, gasPriceWei *big.Int) (*big.Int, error) {
	return _FunctionsCoordinator110.Contract.EstimateCost(&_FunctionsCoordinator110.CallOpts, subscriptionId, data, callbackGasLimit, gasPriceWei)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Caller) GetAdminFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FunctionsCoordinator110.contract.Call(opts, &out, "getAdminFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Session) GetAdminFee() (*big.Int, error) {
	return _FunctionsCoordinator110.Contract.GetAdminFee(&_FunctionsCoordinator110.CallOpts)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110CallerSession) GetAdminFee() (*big.Int, error) {
	return _FunctionsCoordinator110.Contract.GetAdminFee(&_FunctionsCoordinator110.CallOpts)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Caller) GetConfig(opts *bind.CallOpts) (FunctionsBillingConfig, error) {
	var out []interface{}
	err := _FunctionsCoordinator110.contract.Call(opts, &out, "getConfig")

	if err != nil {
		return *new(FunctionsBillingConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(FunctionsBillingConfig)).(*FunctionsBillingConfig)

	return out0, err

}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Session) GetConfig() (FunctionsBillingConfig, error) {
	return _FunctionsCoordinator110.Contract.GetConfig(&_FunctionsCoordinator110.CallOpts)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110CallerSession) GetConfig() (FunctionsBillingConfig, error) {
	return _FunctionsCoordinator110.Contract.GetConfig(&_FunctionsCoordinator110.CallOpts)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Caller) GetDONFee(opts *bind.CallOpts, arg0 []byte) (*big.Int, error) {
	var out []interface{}
	err := _FunctionsCoordinator110.contract.Call(opts, &out, "getDONFee", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Session) GetDONFee(arg0 []byte) (*big.Int, error) {
	return _FunctionsCoordinator110.Contract.GetDONFee(&_FunctionsCoordinator110.CallOpts, arg0)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110CallerSession) GetDONFee(arg0 []byte) (*big.Int, error) {
	return _FunctionsCoordinator110.Contract.GetDONFee(&_FunctionsCoordinator110.CallOpts, arg0)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Caller) GetDONPublicKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _FunctionsCoordinator110.contract.Call(opts, &out, "getDONPublicKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Session) GetDONPublicKey() ([]byte, error) {
	return _FunctionsCoordinator110.Contract.GetDONPublicKey(&_FunctionsCoordinator110.CallOpts)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110CallerSession) GetDONPublicKey() ([]byte, error) {
	return _FunctionsCoordinator110.Contract.GetDONPublicKey(&_FunctionsCoordinator110.CallOpts)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Caller) GetThresholdPublicKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _FunctionsCoordinator110.contract.Call(opts, &out, "getThresholdPublicKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Session) GetThresholdPublicKey() ([]byte, error) {
	return _FunctionsCoordinator110.Contract.GetThresholdPublicKey(&_FunctionsCoordinator110.CallOpts)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110CallerSession) GetThresholdPublicKey() ([]byte, error) {
	return _FunctionsCoordinator110.Contract.GetThresholdPublicKey(&_FunctionsCoordinator110.CallOpts)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Caller) GetWeiPerUnitLink(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FunctionsCoordinator110.contract.Call(opts, &out, "getWeiPerUnitLink")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Session) GetWeiPerUnitLink() (*big.Int, error) {
	return _FunctionsCoordinator110.Contract.GetWeiPerUnitLink(&_FunctionsCoordinator110.CallOpts)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110CallerSession) GetWeiPerUnitLink() (*big.Int, error) {
	return _FunctionsCoordinator110.Contract.GetWeiPerUnitLink(&_FunctionsCoordinator110.CallOpts)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Caller) LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

	error) {
	var out []interface{}
	err := _FunctionsCoordinator110.contract.Call(opts, &out, "latestConfigDetails")

	outstruct := new(LatestConfigDetails)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ConfigDigest = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Session) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _FunctionsCoordinator110.Contract.LatestConfigDetails(&_FunctionsCoordinator110.CallOpts)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110CallerSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _FunctionsCoordinator110.Contract.LatestConfigDetails(&_FunctionsCoordinator110.CallOpts)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Caller) LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

	error) {
	var out []interface{}
	err := _FunctionsCoordinator110.contract.Call(opts, &out, "latestConfigDigestAndEpoch")

	outstruct := new(LatestConfigDigestAndEpoch)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ScanLogs = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfigDigest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Epoch = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Session) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _FunctionsCoordinator110.Contract.LatestConfigDigestAndEpoch(&_FunctionsCoordinator110.CallOpts)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110CallerSession) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _FunctionsCoordinator110.Contract.LatestConfigDigestAndEpoch(&_FunctionsCoordinator110.CallOpts)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FunctionsCoordinator110.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Session) Owner() (common.Address, error) {
	return _FunctionsCoordinator110.Contract.Owner(&_FunctionsCoordinator110.CallOpts)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110CallerSession) Owner() (common.Address, error) {
	return _FunctionsCoordinator110.Contract.Owner(&_FunctionsCoordinator110.CallOpts)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Caller) Transmitters(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _FunctionsCoordinator110.contract.Call(opts, &out, "transmitters")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Session) Transmitters() ([]common.Address, error) {
	return _FunctionsCoordinator110.Contract.Transmitters(&_FunctionsCoordinator110.CallOpts)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110CallerSession) Transmitters() ([]common.Address, error) {
	return _FunctionsCoordinator110.Contract.Transmitters(&_FunctionsCoordinator110.CallOpts)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Caller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _FunctionsCoordinator110.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Session) TypeAndVersion() (string, error) {
	return _FunctionsCoordinator110.Contract.TypeAndVersion(&_FunctionsCoordinator110.CallOpts)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110CallerSession) TypeAndVersion() (string, error) {
	return _FunctionsCoordinator110.Contract.TypeAndVersion(&_FunctionsCoordinator110.CallOpts)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Transactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsCoordinator110.contract.Transact(opts, "acceptOwnership")
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Session) AcceptOwnership() (*types.Transaction, error) {
	return _FunctionsCoordinator110.Contract.AcceptOwnership(&_FunctionsCoordinator110.TransactOpts)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110TransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _FunctionsCoordinator110.Contract.AcceptOwnership(&_FunctionsCoordinator110.TransactOpts)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Transactor) DeleteCommitment(opts *bind.TransactOpts, requestId [32]byte) (*types.Transaction, error) {
	return _FunctionsCoordinator110.contract.Transact(opts, "deleteCommitment", requestId)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Session) DeleteCommitment(requestId [32]byte) (*types.Transaction, error) {
	return _FunctionsCoordinator110.Contract.DeleteCommitment(&_FunctionsCoordinator110.TransactOpts, requestId)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110TransactorSession) DeleteCommitment(requestId [32]byte) (*types.Transaction, error) {
	return _FunctionsCoordinator110.Contract.DeleteCommitment(&_FunctionsCoordinator110.TransactOpts, requestId)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Transactor) OracleWithdraw(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FunctionsCoordinator110.contract.Transact(opts, "oracleWithdraw", recipient, amount)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Session) OracleWithdraw(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FunctionsCoordinator110.Contract.OracleWithdraw(&_FunctionsCoordinator110.TransactOpts, recipient, amount)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110TransactorSession) OracleWithdraw(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FunctionsCoordinator110.Contract.OracleWithdraw(&_FunctionsCoordinator110.TransactOpts, recipient, amount)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Transactor) OracleWithdrawAll(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsCoordinator110.contract.Transact(opts, "oracleWithdrawAll")
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Session) OracleWithdrawAll() (*types.Transaction, error) {
	return _FunctionsCoordinator110.Contract.OracleWithdrawAll(&_FunctionsCoordinator110.TransactOpts)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110TransactorSession) OracleWithdrawAll() (*types.Transaction, error) {
	return _FunctionsCoordinator110.Contract.OracleWithdrawAll(&_FunctionsCoordinator110.TransactOpts)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Transactor) SetConfig(opts *bind.TransactOpts, _signers []common.Address, _transmitters []common.Address, _f uint8, _onchainConfig []byte, _offchainConfigVersion uint64, _offchainConfig []byte) (*types.Transaction, error) {
	return _FunctionsCoordinator110.contract.Transact(opts, "setConfig", _signers, _transmitters, _f, _onchainConfig, _offchainConfigVersion, _offchainConfig)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Session) SetConfig(_signers []common.Address, _transmitters []common.Address, _f uint8, _onchainConfig []byte, _offchainConfigVersion uint64, _offchainConfig []byte) (*types.Transaction, error) {
	return _FunctionsCoordinator110.Contract.SetConfig(&_FunctionsCoordinator110.TransactOpts, _signers, _transmitters, _f, _onchainConfig, _offchainConfigVersion, _offchainConfig)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110TransactorSession) SetConfig(_signers []common.Address, _transmitters []common.Address, _f uint8, _onchainConfig []byte, _offchainConfigVersion uint64, _offchainConfig []byte) (*types.Transaction, error) {
	return _FunctionsCoordinator110.Contract.SetConfig(&_FunctionsCoordinator110.TransactOpts, _signers, _transmitters, _f, _onchainConfig, _offchainConfigVersion, _offchainConfig)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Transactor) SetDONPublicKey(opts *bind.TransactOpts, donPublicKey []byte) (*types.Transaction, error) {
	return _FunctionsCoordinator110.contract.Transact(opts, "setDONPublicKey", donPublicKey)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Session) SetDONPublicKey(donPublicKey []byte) (*types.Transaction, error) {
	return _FunctionsCoordinator110.Contract.SetDONPublicKey(&_FunctionsCoordinator110.TransactOpts, donPublicKey)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110TransactorSession) SetDONPublicKey(donPublicKey []byte) (*types.Transaction, error) {
	return _FunctionsCoordinator110.Contract.SetDONPublicKey(&_FunctionsCoordinator110.TransactOpts, donPublicKey)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Transactor) SetThresholdPublicKey(opts *bind.TransactOpts, thresholdPublicKey []byte) (*types.Transaction, error) {
	return _FunctionsCoordinator110.contract.Transact(opts, "setThresholdPublicKey", thresholdPublicKey)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Session) SetThresholdPublicKey(thresholdPublicKey []byte) (*types.Transaction, error) {
	return _FunctionsCoordinator110.Contract.SetThresholdPublicKey(&_FunctionsCoordinator110.TransactOpts, thresholdPublicKey)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110TransactorSession) SetThresholdPublicKey(thresholdPublicKey []byte) (*types.Transaction, error) {
	return _FunctionsCoordinator110.Contract.SetThresholdPublicKey(&_FunctionsCoordinator110.TransactOpts, thresholdPublicKey)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Transactor) StartRequest(opts *bind.TransactOpts, request FunctionsResponseRequestMeta) (*types.Transaction, error) {
	return _FunctionsCoordinator110.contract.Transact(opts, "startRequest", request)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Session) StartRequest(request FunctionsResponseRequestMeta) (*types.Transaction, error) {
	return _FunctionsCoordinator110.Contract.StartRequest(&_FunctionsCoordinator110.TransactOpts, request)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110TransactorSession) StartRequest(request FunctionsResponseRequestMeta) (*types.Transaction, error) {
	return _FunctionsCoordinator110.Contract.StartRequest(&_FunctionsCoordinator110.TransactOpts, request)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Transactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _FunctionsCoordinator110.contract.Transact(opts, "transferOwnership", to)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Session) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _FunctionsCoordinator110.Contract.TransferOwnership(&_FunctionsCoordinator110.TransactOpts, to)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110TransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _FunctionsCoordinator110.Contract.TransferOwnership(&_FunctionsCoordinator110.TransactOpts, to)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Transactor) Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _FunctionsCoordinator110.contract.Transact(opts, "transmit", reportContext, report, rs, ss, rawVs)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Session) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _FunctionsCoordinator110.Contract.Transmit(&_FunctionsCoordinator110.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110TransactorSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _FunctionsCoordinator110.Contract.Transmit(&_FunctionsCoordinator110.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Transactor) UpdateConfig(opts *bind.TransactOpts, config FunctionsBillingConfig) (*types.Transaction, error) {
	return _FunctionsCoordinator110.contract.Transact(opts, "updateConfig", config)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Session) UpdateConfig(config FunctionsBillingConfig) (*types.Transaction, error) {
	return _FunctionsCoordinator110.Contract.UpdateConfig(&_FunctionsCoordinator110.TransactOpts, config)
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110TransactorSession) UpdateConfig(config FunctionsBillingConfig) (*types.Transaction, error) {
	return _FunctionsCoordinator110.Contract.UpdateConfig(&_FunctionsCoordinator110.TransactOpts, config)
}

type FunctionsCoordinator110CommitmentDeletedIterator struct {
	Event *FunctionsCoordinator110CommitmentDeleted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsCoordinator110CommitmentDeletedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsCoordinator110CommitmentDeleted)
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
		it.Event = new(FunctionsCoordinator110CommitmentDeleted)
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

func (it *FunctionsCoordinator110CommitmentDeletedIterator) Error() error {
	return it.fail
}

func (it *FunctionsCoordinator110CommitmentDeletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsCoordinator110CommitmentDeleted struct {
	RequestId [32]byte
	Raw       types.Log
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Filterer) FilterCommitmentDeleted(opts *bind.FilterOpts) (*FunctionsCoordinator110CommitmentDeletedIterator, error) {

	logs, sub, err := _FunctionsCoordinator110.contract.FilterLogs(opts, "CommitmentDeleted")
	if err != nil {
		return nil, err
	}
	return &FunctionsCoordinator110CommitmentDeletedIterator{contract: _FunctionsCoordinator110.contract, event: "CommitmentDeleted", logs: logs, sub: sub}, nil
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Filterer) WatchCommitmentDeleted(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinator110CommitmentDeleted) (event.Subscription, error) {

	logs, sub, err := _FunctionsCoordinator110.contract.WatchLogs(opts, "CommitmentDeleted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsCoordinator110CommitmentDeleted)
				if err := _FunctionsCoordinator110.contract.UnpackLog(event, "CommitmentDeleted", log); err != nil {
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

func (_FunctionsCoordinator110 *FunctionsCoordinator110Filterer) ParseCommitmentDeleted(log types.Log) (*FunctionsCoordinator110CommitmentDeleted, error) {
	event := new(FunctionsCoordinator110CommitmentDeleted)
	if err := _FunctionsCoordinator110.contract.UnpackLog(event, "CommitmentDeleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsCoordinator110ConfigSetIterator struct {
	Event *FunctionsCoordinator110ConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsCoordinator110ConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsCoordinator110ConfigSet)
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
		it.Event = new(FunctionsCoordinator110ConfigSet)
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

func (it *FunctionsCoordinator110ConfigSetIterator) Error() error {
	return it.fail
}

func (it *FunctionsCoordinator110ConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsCoordinator110ConfigSet struct {
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

func (_FunctionsCoordinator110 *FunctionsCoordinator110Filterer) FilterConfigSet(opts *bind.FilterOpts) (*FunctionsCoordinator110ConfigSetIterator, error) {

	logs, sub, err := _FunctionsCoordinator110.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &FunctionsCoordinator110ConfigSetIterator{contract: _FunctionsCoordinator110.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Filterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinator110ConfigSet) (event.Subscription, error) {

	logs, sub, err := _FunctionsCoordinator110.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsCoordinator110ConfigSet)
				if err := _FunctionsCoordinator110.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_FunctionsCoordinator110 *FunctionsCoordinator110Filterer) ParseConfigSet(log types.Log) (*FunctionsCoordinator110ConfigSet, error) {
	event := new(FunctionsCoordinator110ConfigSet)
	if err := _FunctionsCoordinator110.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsCoordinator110ConfigUpdatedIterator struct {
	Event *FunctionsCoordinator110ConfigUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsCoordinator110ConfigUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsCoordinator110ConfigUpdated)
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
		it.Event = new(FunctionsCoordinator110ConfigUpdated)
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

func (it *FunctionsCoordinator110ConfigUpdatedIterator) Error() error {
	return it.fail
}

func (it *FunctionsCoordinator110ConfigUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsCoordinator110ConfigUpdated struct {
	Config FunctionsBillingConfig
	Raw    types.Log
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Filterer) FilterConfigUpdated(opts *bind.FilterOpts) (*FunctionsCoordinator110ConfigUpdatedIterator, error) {

	logs, sub, err := _FunctionsCoordinator110.contract.FilterLogs(opts, "ConfigUpdated")
	if err != nil {
		return nil, err
	}
	return &FunctionsCoordinator110ConfigUpdatedIterator{contract: _FunctionsCoordinator110.contract, event: "ConfigUpdated", logs: logs, sub: sub}, nil
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Filterer) WatchConfigUpdated(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinator110ConfigUpdated) (event.Subscription, error) {

	logs, sub, err := _FunctionsCoordinator110.contract.WatchLogs(opts, "ConfigUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsCoordinator110ConfigUpdated)
				if err := _FunctionsCoordinator110.contract.UnpackLog(event, "ConfigUpdated", log); err != nil {
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

func (_FunctionsCoordinator110 *FunctionsCoordinator110Filterer) ParseConfigUpdated(log types.Log) (*FunctionsCoordinator110ConfigUpdated, error) {
	event := new(FunctionsCoordinator110ConfigUpdated)
	if err := _FunctionsCoordinator110.contract.UnpackLog(event, "ConfigUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsCoordinator110OracleRequestIterator struct {
	Event *FunctionsCoordinator110OracleRequest

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsCoordinator110OracleRequestIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsCoordinator110OracleRequest)
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
		it.Event = new(FunctionsCoordinator110OracleRequest)
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

func (it *FunctionsCoordinator110OracleRequestIterator) Error() error {
	return it.fail
}

func (it *FunctionsCoordinator110OracleRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsCoordinator110OracleRequest struct {
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

func (_FunctionsCoordinator110 *FunctionsCoordinator110Filterer) FilterOracleRequest(opts *bind.FilterOpts, requestId [][32]byte, requestingContract []common.Address) (*FunctionsCoordinator110OracleRequestIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var requestingContractRule []interface{}
	for _, requestingContractItem := range requestingContract {
		requestingContractRule = append(requestingContractRule, requestingContractItem)
	}

	logs, sub, err := _FunctionsCoordinator110.contract.FilterLogs(opts, "OracleRequest", requestIdRule, requestingContractRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsCoordinator110OracleRequestIterator{contract: _FunctionsCoordinator110.contract, event: "OracleRequest", logs: logs, sub: sub}, nil
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Filterer) WatchOracleRequest(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinator110OracleRequest, requestId [][32]byte, requestingContract []common.Address) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var requestingContractRule []interface{}
	for _, requestingContractItem := range requestingContract {
		requestingContractRule = append(requestingContractRule, requestingContractItem)
	}

	logs, sub, err := _FunctionsCoordinator110.contract.WatchLogs(opts, "OracleRequest", requestIdRule, requestingContractRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsCoordinator110OracleRequest)
				if err := _FunctionsCoordinator110.contract.UnpackLog(event, "OracleRequest", log); err != nil {
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

func (_FunctionsCoordinator110 *FunctionsCoordinator110Filterer) ParseOracleRequest(log types.Log) (*FunctionsCoordinator110OracleRequest, error) {
	event := new(FunctionsCoordinator110OracleRequest)
	if err := _FunctionsCoordinator110.contract.UnpackLog(event, "OracleRequest", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsCoordinator110OracleResponseIterator struct {
	Event *FunctionsCoordinator110OracleResponse

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsCoordinator110OracleResponseIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsCoordinator110OracleResponse)
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
		it.Event = new(FunctionsCoordinator110OracleResponse)
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

func (it *FunctionsCoordinator110OracleResponseIterator) Error() error {
	return it.fail
}

func (it *FunctionsCoordinator110OracleResponseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsCoordinator110OracleResponse struct {
	RequestId   [32]byte
	Transmitter common.Address
	Raw         types.Log
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Filterer) FilterOracleResponse(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsCoordinator110OracleResponseIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsCoordinator110.contract.FilterLogs(opts, "OracleResponse", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsCoordinator110OracleResponseIterator{contract: _FunctionsCoordinator110.contract, event: "OracleResponse", logs: logs, sub: sub}, nil
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Filterer) WatchOracleResponse(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinator110OracleResponse, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsCoordinator110.contract.WatchLogs(opts, "OracleResponse", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsCoordinator110OracleResponse)
				if err := _FunctionsCoordinator110.contract.UnpackLog(event, "OracleResponse", log); err != nil {
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

func (_FunctionsCoordinator110 *FunctionsCoordinator110Filterer) ParseOracleResponse(log types.Log) (*FunctionsCoordinator110OracleResponse, error) {
	event := new(FunctionsCoordinator110OracleResponse)
	if err := _FunctionsCoordinator110.contract.UnpackLog(event, "OracleResponse", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsCoordinator110OwnershipTransferRequestedIterator struct {
	Event *FunctionsCoordinator110OwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsCoordinator110OwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsCoordinator110OwnershipTransferRequested)
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
		it.Event = new(FunctionsCoordinator110OwnershipTransferRequested)
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

func (it *FunctionsCoordinator110OwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *FunctionsCoordinator110OwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsCoordinator110OwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Filterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsCoordinator110OwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsCoordinator110.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsCoordinator110OwnershipTransferRequestedIterator{contract: _FunctionsCoordinator110.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Filterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinator110OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsCoordinator110.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsCoordinator110OwnershipTransferRequested)
				if err := _FunctionsCoordinator110.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_FunctionsCoordinator110 *FunctionsCoordinator110Filterer) ParseOwnershipTransferRequested(log types.Log) (*FunctionsCoordinator110OwnershipTransferRequested, error) {
	event := new(FunctionsCoordinator110OwnershipTransferRequested)
	if err := _FunctionsCoordinator110.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsCoordinator110OwnershipTransferredIterator struct {
	Event *FunctionsCoordinator110OwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsCoordinator110OwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsCoordinator110OwnershipTransferred)
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
		it.Event = new(FunctionsCoordinator110OwnershipTransferred)
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

func (it *FunctionsCoordinator110OwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *FunctionsCoordinator110OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsCoordinator110OwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsCoordinator110OwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsCoordinator110.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsCoordinator110OwnershipTransferredIterator{contract: _FunctionsCoordinator110.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinator110OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsCoordinator110.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsCoordinator110OwnershipTransferred)
				if err := _FunctionsCoordinator110.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_FunctionsCoordinator110 *FunctionsCoordinator110Filterer) ParseOwnershipTransferred(log types.Log) (*FunctionsCoordinator110OwnershipTransferred, error) {
	event := new(FunctionsCoordinator110OwnershipTransferred)
	if err := _FunctionsCoordinator110.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsCoordinator110RequestBilledIterator struct {
	Event *FunctionsCoordinator110RequestBilled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsCoordinator110RequestBilledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsCoordinator110RequestBilled)
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
		it.Event = new(FunctionsCoordinator110RequestBilled)
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

func (it *FunctionsCoordinator110RequestBilledIterator) Error() error {
	return it.fail
}

func (it *FunctionsCoordinator110RequestBilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsCoordinator110RequestBilled struct {
	RequestId         [32]byte
	JuelsPerGas       *big.Int
	L1FeeShareWei     *big.Int
	CallbackCostJuels *big.Int
	TotalCostJuels    *big.Int
	Raw               types.Log
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Filterer) FilterRequestBilled(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsCoordinator110RequestBilledIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsCoordinator110.contract.FilterLogs(opts, "RequestBilled", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsCoordinator110RequestBilledIterator{contract: _FunctionsCoordinator110.contract, event: "RequestBilled", logs: logs, sub: sub}, nil
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Filterer) WatchRequestBilled(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinator110RequestBilled, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsCoordinator110.contract.WatchLogs(opts, "RequestBilled", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsCoordinator110RequestBilled)
				if err := _FunctionsCoordinator110.contract.UnpackLog(event, "RequestBilled", log); err != nil {
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

func (_FunctionsCoordinator110 *FunctionsCoordinator110Filterer) ParseRequestBilled(log types.Log) (*FunctionsCoordinator110RequestBilled, error) {
	event := new(FunctionsCoordinator110RequestBilled)
	if err := _FunctionsCoordinator110.contract.UnpackLog(event, "RequestBilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsCoordinator110TransmittedIterator struct {
	Event *FunctionsCoordinator110Transmitted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsCoordinator110TransmittedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsCoordinator110Transmitted)
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
		it.Event = new(FunctionsCoordinator110Transmitted)
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

func (it *FunctionsCoordinator110TransmittedIterator) Error() error {
	return it.fail
}

func (it *FunctionsCoordinator110TransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsCoordinator110Transmitted struct {
	ConfigDigest [32]byte
	Epoch        uint32
	Raw          types.Log
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Filterer) FilterTransmitted(opts *bind.FilterOpts) (*FunctionsCoordinator110TransmittedIterator, error) {

	logs, sub, err := _FunctionsCoordinator110.contract.FilterLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return &FunctionsCoordinator110TransmittedIterator{contract: _FunctionsCoordinator110.contract, event: "Transmitted", logs: logs, sub: sub}, nil
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110Filterer) WatchTransmitted(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinator110Transmitted) (event.Subscription, error) {

	logs, sub, err := _FunctionsCoordinator110.contract.WatchLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsCoordinator110Transmitted)
				if err := _FunctionsCoordinator110.contract.UnpackLog(event, "Transmitted", log); err != nil {
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

func (_FunctionsCoordinator110 *FunctionsCoordinator110Filterer) ParseTransmitted(log types.Log) (*FunctionsCoordinator110Transmitted, error) {
	event := new(FunctionsCoordinator110Transmitted)
	if err := _FunctionsCoordinator110.contract.UnpackLog(event, "Transmitted", log); err != nil {
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

func (_FunctionsCoordinator110 *FunctionsCoordinator110) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _FunctionsCoordinator110.abi.Events["CommitmentDeleted"].ID:
		return _FunctionsCoordinator110.ParseCommitmentDeleted(log)
	case _FunctionsCoordinator110.abi.Events["ConfigSet"].ID:
		return _FunctionsCoordinator110.ParseConfigSet(log)
	case _FunctionsCoordinator110.abi.Events["ConfigUpdated"].ID:
		return _FunctionsCoordinator110.ParseConfigUpdated(log)
	case _FunctionsCoordinator110.abi.Events["OracleRequest"].ID:
		return _FunctionsCoordinator110.ParseOracleRequest(log)
	case _FunctionsCoordinator110.abi.Events["OracleResponse"].ID:
		return _FunctionsCoordinator110.ParseOracleResponse(log)
	case _FunctionsCoordinator110.abi.Events["OwnershipTransferRequested"].ID:
		return _FunctionsCoordinator110.ParseOwnershipTransferRequested(log)
	case _FunctionsCoordinator110.abi.Events["OwnershipTransferred"].ID:
		return _FunctionsCoordinator110.ParseOwnershipTransferred(log)
	case _FunctionsCoordinator110.abi.Events["RequestBilled"].ID:
		return _FunctionsCoordinator110.ParseRequestBilled(log)
	case _FunctionsCoordinator110.abi.Events["Transmitted"].ID:
		return _FunctionsCoordinator110.ParseTransmitted(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (FunctionsCoordinator110CommitmentDeleted) Topic() common.Hash {
	return common.HexToHash("0x8a4b97add3359bd6bcf5e82874363670eb5ad0f7615abddbd0ed0a3a98f0f416")
}

func (FunctionsCoordinator110ConfigSet) Topic() common.Hash {
	return common.HexToHash("0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05")
}

func (FunctionsCoordinator110ConfigUpdated) Topic() common.Hash {
	return common.HexToHash("0x5f32d06f5e83eda3a68e0e964ef2e6af5cb613e8117aa103c2d6bca5f5184862")
}

func (FunctionsCoordinator110OracleRequest) Topic() common.Hash {
	return common.HexToHash("0xbf50768ccf13bd0110ca6d53a9c4f1f3271abdd4c24a56878863ed25b20598ff")
}

func (FunctionsCoordinator110OracleResponse) Topic() common.Hash {
	return common.HexToHash("0xc708e0440951fd63499c0f7a73819b469ee5dd3ecc356c0ab4eb7f18389009d9")
}

func (FunctionsCoordinator110OwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (FunctionsCoordinator110OwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (FunctionsCoordinator110RequestBilled) Topic() common.Hash {
	return common.HexToHash("0x90815c2e624694e8010bffad2bcefaf96af282ef1bc2ebc0042d1b89a585e046")
}

func (FunctionsCoordinator110Transmitted) Topic() common.Hash {
	return common.HexToHash("0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62")
}

func (_FunctionsCoordinator110 *FunctionsCoordinator110) Address() common.Address {
	return _FunctionsCoordinator110.address
}

type FunctionsCoordinator110Interface interface {
	EstimateCost(opts *bind.CallOpts, subscriptionId uint64, data []byte, callbackGasLimit uint32, gasPriceWei *big.Int) (*big.Int, error)

	GetAdminFee(opts *bind.CallOpts) (*big.Int, error)

	GetConfig(opts *bind.CallOpts) (FunctionsBillingConfig, error)

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

	UpdateConfig(opts *bind.TransactOpts, config FunctionsBillingConfig) (*types.Transaction, error)

	FilterCommitmentDeleted(opts *bind.FilterOpts) (*FunctionsCoordinator110CommitmentDeletedIterator, error)

	WatchCommitmentDeleted(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinator110CommitmentDeleted) (event.Subscription, error)

	ParseCommitmentDeleted(log types.Log) (*FunctionsCoordinator110CommitmentDeleted, error)

	FilterConfigSet(opts *bind.FilterOpts) (*FunctionsCoordinator110ConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinator110ConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*FunctionsCoordinator110ConfigSet, error)

	FilterConfigUpdated(opts *bind.FilterOpts) (*FunctionsCoordinator110ConfigUpdatedIterator, error)

	WatchConfigUpdated(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinator110ConfigUpdated) (event.Subscription, error)

	ParseConfigUpdated(log types.Log) (*FunctionsCoordinator110ConfigUpdated, error)

	FilterOracleRequest(opts *bind.FilterOpts, requestId [][32]byte, requestingContract []common.Address) (*FunctionsCoordinator110OracleRequestIterator, error)

	WatchOracleRequest(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinator110OracleRequest, requestId [][32]byte, requestingContract []common.Address) (event.Subscription, error)

	ParseOracleRequest(log types.Log) (*FunctionsCoordinator110OracleRequest, error)

	FilterOracleResponse(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsCoordinator110OracleResponseIterator, error)

	WatchOracleResponse(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinator110OracleResponse, requestId [][32]byte) (event.Subscription, error)

	ParseOracleResponse(log types.Log) (*FunctionsCoordinator110OracleResponse, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsCoordinator110OwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinator110OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*FunctionsCoordinator110OwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsCoordinator110OwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinator110OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*FunctionsCoordinator110OwnershipTransferred, error)

	FilterRequestBilled(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsCoordinator110RequestBilledIterator, error)

	WatchRequestBilled(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinator110RequestBilled, requestId [][32]byte) (event.Subscription, error)

	ParseRequestBilled(log types.Log) (*FunctionsCoordinator110RequestBilled, error)

	FilterTransmitted(opts *bind.FilterOpts) (*FunctionsCoordinator110TransmittedIterator, error)

	WatchTransmitted(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinator110Transmitted) (event.Subscription, error)

	ParseTransmitted(log types.Log) (*FunctionsCoordinator110Transmitted, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
