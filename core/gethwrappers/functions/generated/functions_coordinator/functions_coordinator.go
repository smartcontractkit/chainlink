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
	Bin: "0x60a06040523480156200001157600080fd5b506040516200578a3803806200578a83398101604081905262000034916200045c565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be8162000132565b5050506001600160a01b038316620000e957604051632530e88560e11b815260040160405180910390fd5b6001600160a01b03808416608052600b80549183166c01000000000000000000000000026001600160601b039092169190911790556200012982620001dd565b5050506200061b565b336001600160a01b038216036200018c5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b620001e76200033b565b80516008805460208401516040808601516060870151608088015160a089015160c08a015161ffff16600160f01b026001600160f01b0364ffffffffff909216600160c81b0264ffffffffff60c81b196001600160481b03909416600160801b0293909316600160801b600160f01b031963ffffffff9586166c010000000000000000000000000263ffffffff60601b19978716680100000000000000000297909716600160401b600160801b0319998716640100000000026001600160401b0319909b169c87169c909c1799909917979097169990991793909317959095169390931793909317929092169390931790915560e0830151610100840151909216600160e01b026001600160e01b0390921691909117600955517f5f32d06f5e83eda3a68e0e964ef2e6af5cb613e8117aa103c2d6bca5f5184862906200033090839062000565565b60405180910390a150565b6000546001600160a01b03163314620003975760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640162000082565b565b80516001600160a01b0381168114620003b157600080fd5b919050565b60405161012081016001600160401b0381118282101715620003e857634e487b7160e01b600052604160045260246000fd5b60405290565b805163ffffffff81168114620003b157600080fd5b80516001600160481b0381168114620003b157600080fd5b805164ffffffffff81168114620003b157600080fd5b805161ffff81168114620003b157600080fd5b80516001600160e01b0381168114620003b157600080fd5b60008060008385036101608112156200047457600080fd5b6200047f8562000399565b935061012080601f19830112156200049657600080fd5b620004a0620003b6565b9150620004b060208701620003ee565b8252620004c060408701620003ee565b6020830152620004d360608701620003ee565b6040830152620004e660808701620003ee565b6060830152620004f960a0870162000403565b60808301526200050c60c087016200041b565b60a08301526200051f60e0870162000431565b60c08301526101006200053481880162000444565b60e084015262000546828801620003ee565b908301525091506200055c610140850162000399565b90509250925092565b815163ffffffff908116825260208084015182169083015260408084015182169083015260608084015191821690830152610120820190506080830151620005b860808401826001600160481b03169052565b5060a0830151620005d260a084018264ffffffffff169052565b5060c0830151620005e960c084018261ffff169052565b5060e08301516200060560e08401826001600160e01b03169052565b506101009283015163ffffffff16919092015290565b6080516151296200066160003960008181610845015281816109d001528181610d0101528181610f920152818161109d01528181611824015261355301526151296000f3fe608060405234801561001057600080fd5b506004361061018d5760003560e01c806381ff7048116100e3578063c3f909d41161008c578063e3d0e71211610066578063e3d0e71214610560578063e4ddcea614610573578063f2fde38b1461058957600080fd5b8063c3f909d4146103b0578063d227d24514610528578063d328a91e1461055857600080fd5b8063a631571e116100bd578063a631571e1461035d578063afcb95d71461037d578063b1dc65a41461039d57600080fd5b806381ff7048146102b557806385b214cf146103225780638da5cb5b1461033557600080fd5b806366316d8d116101455780637f15e1661161011f5780637f15e16614610285578063814118341461029857806381f1b938146102ad57600080fd5b806366316d8d1461026257806379ba5097146102755780637d4807871461027d57600080fd5b8063181f5a7711610176578063181f5a77146101ba5780632a905ccc1461020c57806359b5b7ac1461022e57600080fd5b8063083a5466146101925780631112dadc146101a7575b600080fd5b6101a56101a0366004613a91565b61059c565b005b6101a56101b5366004613c3a565b6105f1565b6101f66040518060400160405280601c81526020017f46756e6374696f6e7320436f6f7264696e61746f722076312e312e300000000081525081565b6040516102039190613d5e565b60405180910390f35b610214610841565b60405168ffffffffffffffffff9091168152602001610203565b61021461023c366004613dff565b50600854700100000000000000000000000000000000900468ffffffffffffffffff1690565b6101a5610270366004613e8e565b6108d7565b6101a5610a8e565b6101a5610b90565b6101a5610293366004613a91565b610dea565b6102a0610e3a565b6040516102039190613f18565b6101f6610ea9565b6102ff60015460025463ffffffff74010000000000000000000000000000000000000000830481169378010000000000000000000000000000000000000000000000009093041691565b6040805163ffffffff948516815293909216602084015290820152606001610203565b6101a5610330366004613f2b565b610f7a565b60005460405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610203565b61037061036b366004613f44565b61102c565b6040516102039190614099565b604080516001815260006020820181905291810191909152606001610203565b6101a56103ab3660046140ed565b6111cd565b61051b6040805161012081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e0810182905261010081019190915250604080516101208101825260085463ffffffff80821683526401000000008204811660208401526801000000000000000082048116938301939093526c01000000000000000000000000810483166060830152700100000000000000000000000000000000810468ffffffffffffffffff166080830152790100000000000000000000000000000000000000000000000000810464ffffffffff1660a08301527e01000000000000000000000000000000000000000000000000000000000000900461ffff1660c08201526009547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff811660e08301527c0100000000000000000000000000000000000000000000000000000000900490911661010082015290565b60405161020391906141a4565b61053b610536366004614294565b6117dd565b6040516bffffffffffffffffffffffff9091168152602001610203565b6101f6611959565b6101a561056e3660046143ad565b6119b0565b61057b61252c565b604051908152602001610203565b6101a561059736600461447a565b612785565b6105a4612799565b60008190036105df576040517f4f42be3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600d6105ec828483614530565b505050565b6105f9612799565b80516008805460208401516040808601516060870151608088015160a089015160c08a015161ffff167e01000000000000000000000000000000000000000000000000000000000000027dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff64ffffffffff909216790100000000000000000000000000000000000000000000000000027fffff0000000000ffffffffffffffffffffffffffffffffffffffffffffffffff68ffffffffffffffffff90941670010000000000000000000000000000000002939093167fffff0000000000000000000000000000ffffffffffffffffffffffffffffffff63ffffffff9586166c01000000000000000000000000027fffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffffff9787166801000000000000000002979097167fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff998716640100000000027fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000909b169c87169c909c1799909917979097169990991793909317959095169390931793909317929092169390931790915560e08301516101008401519092167c0100000000000000000000000000000000000000000000000000000000027bffffffffffffffffffffffffffffffffffffffffffffffffffffffff90921691909117600955517f5f32d06f5e83eda3a68e0e964ef2e6af5cb613e8117aa103c2d6bca5f5184862906108369083906141a4565b60405180910390a150565b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16632a905ccc6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156108ae573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108d29190614656565b905090565b6108df61281c565b806bffffffffffffffffffffffff166000036109195750336000908152600a60205260409020546bffffffffffffffffffffffff16610973565b336000908152600a60205260409020546bffffffffffffffffffffffff80831691161015610973576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b336000908152600a6020526040812080548392906109a09084906bffffffffffffffffffffffff166146a2565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166366316d8d83836040518363ffffffff1660e01b8152600401610a5892919073ffffffffffffffffffffffffffffffffffffffff9290921682526bffffffffffffffffffffffff16602082015260400190565b600060405180830381600087803b158015610a7257600080fd5b505af1158015610a86573d6000803e3d6000fd5b505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610b14576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610b98612799565b610ba061281c565b60006006805480602002602001604051908101604052809291908181526020018280548015610c0557602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610bda575b5050505050905060005b8151811015610de6576000600a6000848481518110610c3057610c306146c7565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252810191909152604001600020546bffffffffffffffffffffffff1690508015610dd5576000600a6000858581518110610c8f57610c8f6146c7565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166366316d8d848481518110610d4d57610d4d6146c7565b6020026020010151836040518363ffffffff1660e01b8152600401610da292919073ffffffffffffffffffffffffffffffffffffffff9290921682526bffffffffffffffffffffffff16602082015260400190565b600060405180830381600087803b158015610dbc57600080fd5b505af1158015610dd0573d6000803e3d6000fd5b505050505b50610ddf816146f6565b9050610c0f565b5050565b610df2612799565b6000819003610e2d576040517f4f42be3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600c6105ec828483614530565b60606006805480602002602001604051908101604052809291908181526020018280548015610e9f57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610e74575b5050505050905090565b6060600d8054610eb890614497565b9050600003610ef3576040517f4f42be3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600d8054610f0090614497565b80601f0160208091040260200160405190810160405280929190818152602001828054610f2c90614497565b8015610e9f5780601f10610f4e57610100808354040283529160200191610e9f565b820191906000526020600020905b815481529060010190602001808311610f5c57509395945050505050565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610fe9576040517fc41a5b0900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008181526007602052604080822091909155517f8a4b97add3359bd6bcf5e82874363670eb5ad0f7615abddbd0ed0a3a98f0f416906108369083815260200190565b6040805161016081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e08101829052610100810182905261012081018290526101408101919091523373ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016146110f4576040517fc41a5b0900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6111056111008361472e565b612a2a565b9050611117606083016040840161447a565b815173ffffffffffffffffffffffffffffffffffffffff91909116907fbf50768ccf13bd0110ca6d53a9c4f1f3271abdd4c24a56878863ed25b20598ff3261116560c0870160a0880161481b565b6111776101608801610140890161447a565b6111818880614838565b6111936101208b016101008c0161489d565b60208b01356111a96101008d0160e08e016148b8565b8b6040516111bf999897969594939291906148d5565b60405180910390a35b919050565b60005a604080518b3580825262ffffff6020808f0135600881901c929092169084015293945092917fb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62910160405180910390a161122e8a8a8a8a8a8a612ef2565b60035460009060029061124c9060ff8082169161010090041661498b565b61125691906149d3565b61126190600161498b565b60ff1690508781146112cf576040517f660bd4ba00000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f77726f6e67206e756d626572206f66207369676e6174757265730000000000006044820152606401610b0b565b87861461135e576040517f660bd4ba00000000000000000000000000000000000000000000000000000000815260206004820152602860248201527f7265706f727420727320616e64207373206d757374206265206f66206571756160448201527f6c206c656e6774680000000000000000000000000000000000000000000000006064820152608401610b0b565b3360009081526004602090815260408083208151808301909252805460ff808216845292939192918401916101009091041660028111156113a1576113a16149f5565b60028111156113b2576113b26149f5565b90525090506002816020015160028111156113cf576113cf6149f5565b1415801561141857506006816000015160ff16815481106113f2576113f26146c7565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff163314155b1561147f576040517f660bd4ba00000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f756e617574686f72697a6564207472616e736d697474657200000000000000006044820152606401610b0b565b5050505061148b613a29565b6000808a8a60405161149e929190614a24565b6040519081900381206114b5918e90602001614a34565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181528282528051602091820120838301909252600080845290830152915060005b898110156117bf57600060018489846020811061151e5761151e6146c7565b61152b91901a601b61498b565b8e8e8681811061153d5761153d6146c7565b905060200201358d8d87818110611556576115566146c7565b9050602002013560405160008152602001604052604051611593949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa1580156115b5573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff811660009081526004602090815290849020838501909452835460ff80821685529296509294508401916101009004166002811115611635576116356149f5565b6002811115611646576116466149f5565b9052509250600183602001516002811115611663576116636149f5565b146116ca576040517f660bd4ba00000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f61646472657373206e6f7420617574686f72697a656420746f207369676e00006044820152606401610b0b565b8251600090879060ff16601f81106116e4576116e46146c7565b602002015173ffffffffffffffffffffffffffffffffffffffff1614611766576040517f660bd4ba00000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f6e6f6e2d756e69717565207369676e61747572650000000000000000000000006044820152606401610b0b565b8086846000015160ff16601f8110611780576117806146c7565b73ffffffffffffffffffffffffffffffffffffffff90921660209290920201526117ab60018661498b565b945050806117b8906146f6565b90506114ff565b5050506117d0833383858e8e612fa9565b5050505050505050505050565b6040517f10fc49c100000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8616600482015263ffffffff831660248201526000907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906310fc49c19060440160006040518083038186803b15801561187b57600080fd5b505afa15801561188f573d6000803e3d6000fd5b5050505066038d7ea4c680008211156118d4576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6118dc610841565b68ffffffffffffffffff1661192686868080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061023c92505050565b68ffffffffffffffffff1661193b85856131a8565b6119459190614a48565b61194f9190614a48565b9695505050505050565b6060600c805461196890614497565b90506000036119a3576040517f4f42be3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600c8054610f0090614497565b855185518560ff16601f831115611a23576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f746f6f206d616e79207369676e657273000000000000000000000000000000006044820152606401610b0b565b80600003611a8d576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f66206d75737420626520706f73697469766500000000000000000000000000006044820152606401610b0b565b818314611b1b576040517f89a61989000000000000000000000000000000000000000000000000000000008152602060048201526024808201527f6f7261636c6520616464726573736573206f7574206f6620726567697374726160448201527f74696f6e000000000000000000000000000000000000000000000000000000006064820152608401610b0b565b611b26816003614a6d565b8311611b8e576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f6661756c74792d6f7261636c65206620746f6f206869676800000000000000006044820152606401610b0b565b611b96612799565b6040805160c0810182528a8152602081018a905260ff89169181018290526060810188905267ffffffffffffffff8716608082015260a0810186905290611bdd90886132e9565b60055415611d9257600554600090611bf790600190614a84565b9050600060058281548110611c0e57611c0e6146c7565b60009182526020822001546006805473ffffffffffffffffffffffffffffffffffffffff90921693509084908110611c4857611c486146c7565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff85811684526004909252604080842080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000090811690915592909116808452922080549091169055600580549192509080611cc857611cc8614a97565b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff00000000000000000000000000000000000000001690550190556006805480611d3157611d31614a97565b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff000000000000000000000000000000000000000016905501905550611bdd915050565b60005b81515181101561234957815180516000919083908110611db757611db76146c7565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1603611e3c576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f7369676e6572206d757374206e6f7420626520656d70747900000000000000006044820152606401610b0b565b600073ffffffffffffffffffffffffffffffffffffffff1682602001518281518110611e6a57611e6a6146c7565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1603611eef576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f7472616e736d6974746572206d757374206e6f7420626520656d7074790000006044820152606401610b0b565b60006004600084600001518481518110611f0b57611f0b6146c7565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff16825281019190915260400160002054610100900460ff166002811115611f5557611f556149f5565b14611fbc576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265706561746564207369676e657220616464726573730000000000000000006044820152606401610b0b565b6040805180820190915260ff82168152600160208201528251805160049160009185908110611fed57611fed6146c7565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040016000208251815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082168117835592840151919283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000161761010083600281111561208e5761208e6149f5565b02179055506000915061209e9050565b60046000846020015184815181106120b8576120b86146c7565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff16825281019190915260400160002054610100900460ff166002811115612102576121026149f5565b14612169576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f7265706561746564207472616e736d69747465722061646472657373000000006044820152606401610b0b565b6040805180820190915260ff82168152602081016002815250600460008460200151848151811061219c5761219c6146c7565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040016000208251815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082168117835592840151919283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000161761010083600281111561223d5761223d6149f5565b02179055505082518051600592508390811061225b5761225b6146c7565b602090810291909101810151825460018101845560009384529282902090920180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff90931692909217909155820151805160069190839081106122d7576122d76146c7565b60209081029190910181015182546001810184556000938452919092200180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff90921691909117905580612341816146f6565b915050611d95565b506040810151600380547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff909216919091179055600180547fffffffff00000000ffffffffffffffffffffffffffffffffffffffffffffffff8116780100000000000000000000000000000000000000000000000063ffffffff438116820292909217808555920481169291829160149161240191849174010000000000000000000000000000000000000000900416614ac6565b92506101000a81548163ffffffff021916908363ffffffff1602179055506124604630600160149054906101000a900463ffffffff1663ffffffff16856000015186602001518760400151886060015189608001518a60a001516132f9565b600281905582518051600380547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1661010060ff9093169290920291909117905560015460208501516040808701516060880151608089015160a08a015193517f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e0598612517988b9891977401000000000000000000000000000000000000000090920463ffffffff16969095919491939192614ae3565b60405180910390a15050505050505050505050565b604080516101208101825260085463ffffffff80821683526401000000008204811660208401526801000000000000000082048116838501526c0100000000000000000000000080830482166060850152700100000000000000000000000000000000830468ffffffffffffffffff166080850152790100000000000000000000000000000000000000000000000000830464ffffffffff1660a0808601919091527e0100000000000000000000000000000000000000000000000000000000000090930461ffff1660c08501526009547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff811660e08601527c01000000000000000000000000000000000000000000000000000000009004909116610100840152600b5484517ffeaf968c00000000000000000000000000000000000000000000000000000000815294516000958694859490930473ffffffffffffffffffffffffffffffffffffffff169263feaf968c926004808401938290030181865afa1580156126ba573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906126de9190614b93565b5093505092505080426126f19190614a84565b836020015163ffffffff1610801561271357506000836020015163ffffffff16115b1561274157505060e001517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16919050565b6000821361277e576040517f43d4cf6600000000000000000000000000000000000000000000000000000000815260048101839052602401610b0b565b5092915050565b61278d612799565b612796816133a4565b50565b60005473ffffffffffffffffffffffffffffffffffffffff16331461281a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610b0b565b565b600b546bffffffffffffffffffffffff1660000361283657565b6000600680548060200260200160405190810160405280929190818152602001828054801561289b57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311612870575b50505050509050600081519050806000036128e2576040517f30274b3a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600b546000906129019083906bffffffffffffffffffffffff16614be3565b905060005b828110156129cc5781600a6000868481518110612925576129256146c7565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282829054906101000a90046bffffffffffffffffffffffff1661298d9190614a48565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550806129c5906146f6565b9050612906565b506129d78282614c0e565b600b80546000906129f79084906bffffffffffffffffffffffff166146a2565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550505050565b6040805161016081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e0810182905261010081018290526101208101829052610140810191909152604080516101208101825260085463ffffffff80821683526401000000008204811660208401526801000000000000000082048116938301939093526c0100000000000000000000000081048316606083015268ffffffffffffffffff700100000000000000000000000000000000820416608083015264ffffffffff79010000000000000000000000000000000000000000000000000082041660a083015261ffff7e01000000000000000000000000000000000000000000000000000000000000909104811660c083018190526009547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff811660e08501527c0100000000000000000000000000000000000000000000000000000000900490931661010080840191909152850151919291161115612be5576040517fdada758700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600854600090700100000000000000000000000000000000900468ffffffffffffffffff1690506000846080015168ffffffffffffffffff168268ffffffffffffffffff16612c388760e001513a6131a8565b612c429190614a48565b612c4c9190614a48565b9050806bffffffffffffffffffffffff1685606001516bffffffffffffffffffffffff161015612ca8576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600083610100015163ffffffff1642612cc19190614c36565b905060003087604001518860a001518960c001516001612ce19190614c49565b8a5180516020918201206101008d015160e08e0151604051612d9598979695948c918c9132910173ffffffffffffffffffffffffffffffffffffffff9a8b168152988a1660208a015267ffffffffffffffff97881660408a0152959096166060880152608087019390935261ffff9190911660a086015263ffffffff90811660c08601526bffffffffffffffffffffffff9190911660e0850152919091166101008301529091166101208201526101400190565b6040516020818303038152906040528051906020012090506040518061016001604052808281526020013073ffffffffffffffffffffffffffffffffffffffff168152602001846bffffffffffffffffffffffff168152602001886040015173ffffffffffffffffffffffffffffffffffffffff1681526020018860a0015167ffffffffffffffff1681526020018860e0015163ffffffff168152602001886080015168ffffffffffffffffff1681526020018568ffffffffffffffffff168152602001866040015163ffffffff1664ffffffffff168152602001866060015163ffffffff1664ffffffffff1681526020018363ffffffff16815250955085604051602001612ea49190614099565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152918152815160209283012060009384526007909252909120555092949350505050565b6000612eff826020614a6d565b612f0a856020614a6d565b612f1688610144614c36565b612f209190614c36565b612f2a9190614c36565b612f35906000614c36565b9050368114612fa0576040517f660bd4ba00000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f63616c6c64617461206c656e677468206d69736d6174636800000000000000006044820152606401610b0b565b50505050505050565b600080808080612fbb86880188614d45565b84519499509297509095509350915060ff16801580612fdb575084518114155b80612fe7575083518114155b80612ff3575082518114155b80612fff575081518114155b15613066576040517f660bd4ba00000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f4669656c6473206d75737420626520657175616c206c656e67746800000000006044820152606401610b0b565b60005b818110156131995760006130fe888381518110613088576130886146c7565b60200260200101518884815181106130a2576130a26146c7565b60200260200101518885815181106130bc576130bc6146c7565b60200260200101518886815181106130d6576130d66146c7565b60200260200101518887815181106130f0576130f06146c7565b602002602001015188613499565b90506000816006811115613114576131146149f5565b14806131315750600181600681111561312f5761312f6149f5565b145b1561318857878281518110613148576131486146c7565b60209081029190910181015160405133815290917fc708e0440951fd63499c0f7a73819b469ee5dd3ecc356c0ab4eb7f18389009d9910160405180910390a25b50613192816146f6565b9050613069565b50505050505050505050505050565b600854600090790100000000000000000000000000000000000000000000000000900464ffffffffff1682101561320357600854790100000000000000000000000000000000000000000000000000900464ffffffffff1691505b6008546000906127109061321d9063ffffffff1685614a6d565b6132279190614e17565b6132319084614c36565b600854909150600090859061326a9063ffffffff6c01000000000000000000000000820481169168010000000000000000900416614ac6565b6132749190614ac6565b63ffffffff16905060006132be6000368080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506137ad92505050565b90506132dd816132ce8486614a6d565b6132d89190614c36565b6138ef565b93505050505b92915050565b60065415610de657610de661281c565b6000808a8a8a8a8a8a8a8a8a60405160200161331d99989796959493929190614e2b565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179150509998505050505050505050565b3373ffffffffffffffffffffffffffffffffffffffff821603613423576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610b0b565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600080848060200190518101906134b09190614ef7565b905060003a8261012001518361010001516134cb9190614fbf565b64ffffffffff166134dc9190614a6d565b905060008460ff166135246000368080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506137ad92505050565b61352e9190614e17565b9050600061353f6132d88385614c36565b9050600061354c3a6138ef565b90506000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663330605298e8e868b60e0015168ffffffffffffffffff16896135ab9190614a48565b338d6040518763ffffffff1660e01b81526004016135ce96959493929190614fdd565b60408051808303816000875af11580156135ec573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906136109190615059565b90925090506000826006811115613629576136296149f5565b148061364657506001826006811115613644576136446149f5565b145b1561379c5760008e8152600760205260408120556136648185614a48565b336000908152600a6020526040812080547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff93841617905560e0890151600b805468ffffffffffffffffff909216939092916136d091859116614a48565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055508d7f90815c2e624694e8010bffad2bcefaf96af282ef1bc2ebc0042d1b89a585e0468487848b60c0015168ffffffffffffffffff168c60e0015168ffffffffffffffffff16878b61374f9190614a48565b6137599190614a48565b6137639190614a48565b604080516bffffffffffffffffffffffff9586168152602081019490945291841683830152909216606082015290519081900360800190a25b509c9b505050505050505050505050565b6000466137b98161391d565b1561383557606c73ffffffffffffffffffffffffffffffffffffffff1663c6f7de0e6040518163ffffffff1660e01b8152600401602060405180830381865afa15801561380a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061382e919061508c565b9392505050565b61383e81613940565b156138e65773420000000000000000000000000000000000000f73ffffffffffffffffffffffffffffffffffffffff166349948e0e846040518060800160405280604881526020016150d56048913960405160200161389e9291906150a5565b6040516020818303038152906040526040518263ffffffff1660e01b81526004016138c99190613d5e565b602060405180830381865afa15801561380a573d6000803e3d6000fd5b50600092915050565b60006132e36138fc61252c565b61390e84670de0b6b3a7640000614a6d565b6139189190614e17565b613987565b600061a4b1821480613931575062066eed82145b806132e357505062066eee1490565b6000600a82148061395257506101a482145b8061395f575062aa37dc82145b8061396b575061210582145b80613978575062014a3382145b806132e357505062014a341490565b60006bffffffffffffffffffffffff821115613a25576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f53616665436173743a2076616c756520646f65736e27742066697420696e203960448201527f36206269747300000000000000000000000000000000000000000000000000006064820152608401610b0b565b5090565b604051806103e00160405280601f906020820280368337509192915050565b60008083601f840112613a5a57600080fd5b50813567ffffffffffffffff811115613a7257600080fd5b602083019150836020828501011115613a8a57600080fd5b9250929050565b60008060208385031215613aa457600080fd5b823567ffffffffffffffff811115613abb57600080fd5b613ac785828601613a48565b90969095509350505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610120810167ffffffffffffffff81118282101715613b2657613b26613ad3565b60405290565b604051610160810167ffffffffffffffff81118282101715613b2657613b26613ad3565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715613b9757613b97613ad3565b604052919050565b63ffffffff8116811461279657600080fd5b80356111c881613b9f565b68ffffffffffffffffff8116811461279657600080fd5b80356111c881613bbc565b64ffffffffff8116811461279657600080fd5b80356111c881613bde565b803561ffff811681146111c857600080fd5b80357bffffffffffffffffffffffffffffffffffffffffffffffffffffffff811681146111c857600080fd5b60006101208284031215613c4d57600080fd5b613c55613b02565b613c5e83613bb1565b8152613c6c60208401613bb1565b6020820152613c7d60408401613bb1565b6040820152613c8e60608401613bb1565b6060820152613c9f60808401613bd3565b6080820152613cb060a08401613bf1565b60a0820152613cc160c08401613bfc565b60c0820152613cd260e08401613c0e565b60e0820152610100613ce5818501613bb1565b908201529392505050565b60005b83811015613d0b578181015183820152602001613cf3565b50506000910152565b60008151808452613d2c816020860160208601613cf0565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b60208152600061382e6020830184613d14565b600082601f830112613d8257600080fd5b813567ffffffffffffffff811115613d9c57613d9c613ad3565b613dcd60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601613b50565b818152846020838601011115613de257600080fd5b816020850160208301376000918101602001919091529392505050565b600060208284031215613e1157600080fd5b813567ffffffffffffffff811115613e2857600080fd5b613e3484828501613d71565b949350505050565b73ffffffffffffffffffffffffffffffffffffffff8116811461279657600080fd5b80356111c881613e3c565b6bffffffffffffffffffffffff8116811461279657600080fd5b80356111c881613e69565b60008060408385031215613ea157600080fd5b8235613eac81613e3c565b91506020830135613ebc81613e69565b809150509250929050565b600081518084526020808501945080840160005b83811015613f0d57815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101613edb565b509495945050505050565b60208152600061382e6020830184613ec7565b600060208284031215613f3d57600080fd5b5035919050565b600060208284031215613f5657600080fd5b813567ffffffffffffffff811115613f6d57600080fd5b8201610160818503121561382e57600080fd5b805182526020810151613fab602084018273ffffffffffffffffffffffffffffffffffffffff169052565b506040810151613fcb60408401826bffffffffffffffffffffffff169052565b506060810151613ff3606084018273ffffffffffffffffffffffffffffffffffffffff169052565b50608081015161400f608084018267ffffffffffffffff169052565b5060a081015161402760a084018263ffffffff169052565b5060c081015161404460c084018268ffffffffffffffffff169052565b5060e081015161406160e084018268ffffffffffffffffff169052565b506101008181015164ffffffffff9081169184019190915261012080830151909116908301526101409081015163ffffffff16910152565b61016081016132e38284613f80565b60008083601f8401126140ba57600080fd5b50813567ffffffffffffffff8111156140d257600080fd5b6020830191508360208260051b8501011115613a8a57600080fd5b60008060008060008060008060e0898b03121561410957600080fd5b606089018a81111561411a57600080fd5b8998503567ffffffffffffffff8082111561413457600080fd5b6141408c838d01613a48565b909950975060808b013591508082111561415957600080fd5b6141658c838d016140a8565b909750955060a08b013591508082111561417e57600080fd5b5061418b8b828c016140a8565b999c989b50969995989497949560c00135949350505050565b815163ffffffff9081168252602080840151821690830152604080840151821690830152606080840151918216908301526101208201905060808301516141f8608084018268ffffffffffffffffff169052565b5060a083015161421160a084018264ffffffffff169052565b5060c083015161422760c084018261ffff169052565b5060e083015161425760e08401827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff169052565b506101008381015163ffffffff8116848301525b505092915050565b67ffffffffffffffff8116811461279657600080fd5b80356111c881614273565b6000806000806000608086880312156142ac57600080fd5b85356142b781614273565b9450602086013567ffffffffffffffff8111156142d357600080fd5b6142df88828901613a48565b90955093505060408601356142f381613b9f565b949793965091946060013592915050565b600067ffffffffffffffff82111561431e5761431e613ad3565b5060051b60200190565b600082601f83011261433957600080fd5b8135602061434e61434983614304565b613b50565b82815260059290921b8401810191818101908684111561436d57600080fd5b8286015b8481101561439157803561438481613e3c565b8352918301918301614371565b509695505050505050565b803560ff811681146111c857600080fd5b60008060008060008060c087890312156143c657600080fd5b863567ffffffffffffffff808211156143de57600080fd5b6143ea8a838b01614328565b9750602089013591508082111561440057600080fd5b61440c8a838b01614328565b965061441a60408a0161439c565b9550606089013591508082111561443057600080fd5b61443c8a838b01613d71565b945061444a60808a01614289565b935060a089013591508082111561446057600080fd5b5061446d89828a01613d71565b9150509295509295509295565b60006020828403121561448c57600080fd5b813561382e81613e3c565b600181811c908216806144ab57607f821691505b6020821081036144e4577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f8211156105ec57600081815260208120601f850160051c810160208610156145115750805b601f850160051c820191505b81811015610a865782815560010161451d565b67ffffffffffffffff83111561454857614548613ad3565b61455c836145568354614497565b836144ea565b6000601f8411600181146145ae57600085156145785750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b178355614644565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b828110156145fd57868501358255602094850194600190920191016145dd565b5086821015614638577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b80516111c881613bbc565b60006020828403121561466857600080fd5b815161382e81613bbc565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6bffffffffffffffffffffffff82811682821603908082111561277e5761277e614673565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361472757614727614673565b5060010190565b6000610160823603121561474157600080fd5b614749613b2c565b823567ffffffffffffffff81111561476057600080fd5b61476c36828601613d71565b8252506020830135602082015261478560408401613e5e565b604082015261479660608401613e83565b60608201526147a760808401613bd3565b60808201526147b860a08401614289565b60a08201526147c960c08401614289565b60c08201526147da60e08401613bb1565b60e08201526101006147ed818501613bfc565b908201526101206147ff848201614289565b90820152610140614811848201613e5e565b9082015292915050565b60006020828403121561482d57600080fd5b813561382e81614273565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261486d57600080fd5b83018035915067ffffffffffffffff82111561488857600080fd5b602001915036819003821315613a8a57600080fd5b6000602082840312156148af57600080fd5b61382e82613bfc565b6000602082840312156148ca57600080fd5b813561382e81613b9f565b73ffffffffffffffffffffffffffffffffffffffff8a8116825267ffffffffffffffff8a166020830152881660408201526102406060820181905281018690526000610260878982850137600083890182015261ffff8716608084015260a0830186905263ffffffff851660c0840152601f88017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016830101905061497d60e0830184613f80565b9a9950505050505050505050565b60ff81811683821601908111156132e3576132e3614673565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b600060ff8316806149e6576149e66149a4565b8060ff84160491505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b8183823760009101908152919050565b828152606082602083013760800192915050565b6bffffffffffffffffffffffff81811683821601908082111561277e5761277e614673565b80820281158282048414176132e3576132e3614673565b818103818111156132e3576132e3614673565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b63ffffffff81811683821601908082111561277e5761277e614673565b600061012063ffffffff808d1684528b6020850152808b16604085015250806060840152614b138184018a613ec7565b90508281036080840152614b278189613ec7565b905060ff871660a084015282810360c0840152614b448187613d14565b905067ffffffffffffffff851660e0840152828103610100840152614b698185613d14565b9c9b505050505050505050505050565b805169ffffffffffffffffffff811681146111c857600080fd5b600080600080600060a08688031215614bab57600080fd5b614bb486614b79565b9450602086015193506040860151925060608601519150614bd760808701614b79565b90509295509295909350565b60006bffffffffffffffffffffffff80841680614c0257614c026149a4565b92169190910492915050565b6bffffffffffffffffffffffff81811683821602808216919082811461426b5761426b614673565b808201808211156132e3576132e3614673565b67ffffffffffffffff81811683821601908082111561277e5761277e614673565b600082601f830112614c7b57600080fd5b81356020614c8b61434983614304565b82815260059290921b84018101918181019086841115614caa57600080fd5b8286015b848110156143915780358352918301918301614cae565b600082601f830112614cd657600080fd5b81356020614ce661434983614304565b82815260059290921b84018101918181019086841115614d0557600080fd5b8286015b8481101561439157803567ffffffffffffffff811115614d295760008081fd5b614d378986838b0101613d71565b845250918301918301614d09565b600080600080600060a08688031215614d5d57600080fd5b853567ffffffffffffffff80821115614d7557600080fd5b614d8189838a01614c6a565b96506020880135915080821115614d9757600080fd5b614da389838a01614cc5565b95506040880135915080821115614db957600080fd5b614dc589838a01614cc5565b94506060880135915080821115614ddb57600080fd5b614de789838a01614cc5565b93506080880135915080821115614dfd57600080fd5b50614e0a88828901614cc5565b9150509295509295909350565b600082614e2657614e266149a4565b500490565b60006101208b835273ffffffffffffffffffffffffffffffffffffffff8b16602084015267ffffffffffffffff808b166040850152816060850152614e728285018b613ec7565b91508382036080850152614e86828a613ec7565b915060ff881660a085015283820360c0850152614ea38288613d14565b90861660e08501528381036101008501529050614b698185613d14565b80516111c881613e3c565b80516111c881613e69565b80516111c881614273565b80516111c881613b9f565b80516111c881613bde565b60006101608284031215614f0a57600080fd5b614f12613b2c565b82518152614f2260208401614ec0565b6020820152614f3360408401614ecb565b6040820152614f4460608401614ec0565b6060820152614f5560808401614ed6565b6080820152614f6660a08401614ee1565b60a0820152614f7760c0840161464b565b60c0820152614f8860e0840161464b565b60e0820152610100614f9b818501614eec565b90820152610120614fad848201614eec565b90820152610140613ce5848201614ee1565b64ffffffffff81811683821601908082111561277e5761277e614673565b6000610200808352614ff18184018a613d14565b905082810360208401526150058189613d14565b6bffffffffffffffffffffffff88811660408601528716606085015273ffffffffffffffffffffffffffffffffffffffff86166080850152915061504e905060a0830184613f80565b979650505050505050565b6000806040838503121561506c57600080fd5b82516007811061507b57600080fd5b6020840151909250613ebc81613e69565b60006020828403121561509e57600080fd5b5051919050565b600083516150b7818460208801613cf0565b8351908301906150cb818360208801613cf0565b0194935050505056fe307866666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666a164736f6c6343000813000a",
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
