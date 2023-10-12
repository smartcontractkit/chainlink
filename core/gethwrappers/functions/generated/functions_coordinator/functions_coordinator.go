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
	RequestTimeoutSeconds               uint32
	DonFee                              *big.Int
	MaxSupportedRequestDataVersion      uint16
	FallbackNativePerUnitLink           *big.Int
	MinimumEstimateGasPrice             uint32
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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"fulfillmentGasPriceOverEstimationBP\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"feedStalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasOverheadBeforeCallback\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasOverheadAfterCallback\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"requestTimeoutSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint72\",\"name\":\"donFee\",\"type\":\"uint72\"},{\"internalType\":\"uint16\",\"name\":\"maxSupportedRequestDataVersion\",\"type\":\"uint16\"},{\"internalType\":\"uint224\",\"name\":\"fallbackNativePerUnitLink\",\"type\":\"uint224\"},{\"internalType\":\"uint32\",\"name\":\"minimumEstimateGasPrice\",\"type\":\"uint32\"}],\"internalType\":\"structFunctionsBilling.Config\",\"name\":\"config\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"linkToNativeFeed\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"EmptyPublicKey\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InconsistentReportData\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"InvalidConfig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"linkWei\",\"type\":\"int256\"}],\"name\":\"InvalidLinkWeiPrice\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSubscription\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"MustBeSubOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoTransmittersSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByRouter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByRouterOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentTooLarge\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReportInvalid\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RouterMustBeSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedPublicKeyChange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedSender\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnsupportedRequestDataVersion\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"CommitmentDeleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"fulfillmentGasPriceOverEstimationBP\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"feedStalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasOverheadBeforeCallback\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasOverheadAfterCallback\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"requestTimeoutSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint72\",\"name\":\"donFee\",\"type\":\"uint72\"},{\"internalType\":\"uint16\",\"name\":\"maxSupportedRequestDataVersion\",\"type\":\"uint16\"},{\"internalType\":\"uint224\",\"name\":\"fallbackNativePerUnitLink\",\"type\":\"uint224\"},{\"internalType\":\"uint32\",\"name\":\"minimumEstimateGasPrice\",\"type\":\"uint32\"}],\"indexed\":false,\"internalType\":\"structFunctionsBilling.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"ConfigUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requestingContract\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requestInitiator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"subscriptionOwner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"dataVersion\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"flags\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"callbackGasLimit\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"estimatedTotalCostJuels\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint72\",\"name\":\"adminFee\",\"type\":\"uint72\"},{\"internalType\":\"uint72\",\"name\":\"donFee\",\"type\":\"uint72\"},{\"internalType\":\"uint40\",\"name\":\"gasOverheadBeforeCallback\",\"type\":\"uint40\"},{\"internalType\":\"uint40\",\"name\":\"gasOverheadAfterCallback\",\"type\":\"uint40\"},{\"internalType\":\"uint32\",\"name\":\"timeoutTimestamp\",\"type\":\"uint32\"}],\"indexed\":false,\"internalType\":\"structFunctionsResponse.Commitment\",\"name\":\"commitment\",\"type\":\"tuple\"}],\"name\":\"OracleRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"OracleResponse\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"deleteCommitment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasPriceWei\",\"type\":\"uint256\"}],\"name\":\"estimateCost\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAdminFee\",\"outputs\":[{\"internalType\":\"uint72\",\"name\":\"\",\"type\":\"uint72\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"fulfillmentGasPriceOverEstimationBP\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"feedStalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasOverheadBeforeCallback\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasOverheadAfterCallback\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"requestTimeoutSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint72\",\"name\":\"donFee\",\"type\":\"uint72\"},{\"internalType\":\"uint16\",\"name\":\"maxSupportedRequestDataVersion\",\"type\":\"uint16\"},{\"internalType\":\"uint224\",\"name\":\"fallbackNativePerUnitLink\",\"type\":\"uint224\"},{\"internalType\":\"uint32\",\"name\":\"minimumEstimateGasPrice\",\"type\":\"uint32\"}],\"internalType\":\"structFunctionsBilling.Config\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"getDONFee\",\"outputs\":[{\"internalType\":\"uint72\",\"name\":\"\",\"type\":\"uint72\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDONPublicKey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getThresholdPublicKey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getWeiPerUnitLink\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"oracleWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"oracleWithdrawAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"_f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"_onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"_offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"_offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"donPublicKey\",\"type\":\"bytes\"}],\"name\":\"setDONPublicKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"thresholdPublicKey\",\"type\":\"bytes\"}],\"name\":\"setThresholdPublicKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"flags\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"requestingContract\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"availableBalance\",\"type\":\"uint96\"},{\"internalType\":\"uint72\",\"name\":\"adminFee\",\"type\":\"uint72\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"initiatedRequests\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"dataVersion\",\"type\":\"uint16\"},{\"internalType\":\"uint64\",\"name\":\"completedRequests\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"subscriptionOwner\",\"type\":\"address\"}],\"internalType\":\"structFunctionsResponse.RequestMeta\",\"name\":\"request\",\"type\":\"tuple\"}],\"name\":\"startRequest\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"estimatedTotalCostJuels\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint72\",\"name\":\"adminFee\",\"type\":\"uint72\"},{\"internalType\":\"uint72\",\"name\":\"donFee\",\"type\":\"uint72\"},{\"internalType\":\"uint40\",\"name\":\"gasOverheadBeforeCallback\",\"type\":\"uint40\"},{\"internalType\":\"uint40\",\"name\":\"gasOverheadAfterCallback\",\"type\":\"uint40\"},{\"internalType\":\"uint32\",\"name\":\"timeoutTimestamp\",\"type\":\"uint32\"}],\"internalType\":\"structFunctionsResponse.Commitment\",\"name\":\"commitment\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"transmitters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"fulfillmentGasPriceOverEstimationBP\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"feedStalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasOverheadBeforeCallback\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasOverheadAfterCallback\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"requestTimeoutSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint72\",\"name\":\"donFee\",\"type\":\"uint72\"},{\"internalType\":\"uint16\",\"name\":\"maxSupportedRequestDataVersion\",\"type\":\"uint16\"},{\"internalType\":\"uint224\",\"name\":\"fallbackNativePerUnitLink\",\"type\":\"uint224\"},{\"internalType\":\"uint32\",\"name\":\"minimumEstimateGasPrice\",\"type\":\"uint32\"}],\"internalType\":\"structFunctionsBilling.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"updateConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60c06040523480156200001157600080fd5b506040516200523538038062005235833981016040819052620000349162000453565b8282828260013380600081620000915760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c457620000c48162000140565b50505015156080526001600160a01b038116620000f457604051632530e88560e11b815260040160405180910390fd5b6001600160a01b0390811660a052600b80549183166c01000000000000000000000000026001600160601b039092169190911790556200013482620001eb565b50505050505062000604565b336001600160a01b038216036200019a5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000088565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b620001f56200033e565b80516008805460208401516040808601516060870151608088015160a089015160c08a015161ffff16600160e81b0261ffff60e81b196001600160481b03909216600160a01b02600160a01b600160e81b031963ffffffff948516600160801b0216600160801b600160e81b03199585166c010000000000000000000000000263ffffffff60601b19978616680100000000000000000297909716600160401b600160801b0319998616640100000000026001600160401b0319909b169c86169c909c1799909917979097169990991793909317919091169490941792909217919091169390931790915560e0830151610100840151909116600160e01b026001600160e01b0390911617600955517fdfcfddac88b9a703c9943a696e4da4a2d76b5b5c285a821bbcc54052775ecd2b90620003339083906200055c565b60405180910390a150565b620003486200034a565b565b6000546001600160a01b03163314620003485760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640162000088565b80516001600160a01b0381168114620003be57600080fd5b919050565b60405161012081016001600160401b0381118282101715620003f557634e487b7160e01b600052604160045260246000fd5b60405290565b805163ffffffff81168114620003be57600080fd5b80516001600160481b0381168114620003be57600080fd5b805161ffff81168114620003be57600080fd5b80516001600160e01b0381168114620003be57600080fd5b60008060008385036101608112156200046b57600080fd5b6200047685620003a6565b935061012080601f19830112156200048d57600080fd5b62000497620003c3565b9150620004a760208701620003fb565b8252620004b760408701620003fb565b6020830152620004ca60608701620003fb565b6040830152620004dd60808701620003fb565b6060830152620004f060a08701620003fb565b60808301526200050360c0870162000410565b60a08301526200051660e0870162000428565b60c08301526101006200052b8188016200043b565b60e08401526200053d828801620003fb565b90830152509150620005536101408501620003a6565b90509250925092565b815163ffffffff9081168252602080840151821690830152604080840151821690830152606080840151821690830152608080840151918216908301526101208201905060a0830151620005bb60a08401826001600160481b03169052565b5060c0830151620005d260c084018261ffff169052565b5060e0830151620005ee60e08401826001600160e01b03169052565b506101009283015163ffffffff16919092015290565b60805160a051614be1620006546000396000818161082e015281816109bc01528181610c8f01528181610f230152818161102e015281816118190152613318015260006112570152614be16000f3fe608060405234801561001057600080fd5b506004361061018d5760003560e01c806381ff7048116100e3578063c3f909d41161008c578063e3d0e71211610066578063e3d0e71214610559578063e4ddcea61461056c578063f2fde38b1461058257600080fd5b8063c3f909d4146103b4578063d227d24514610521578063d328a91e1461055157600080fd5b8063a631571e116100bd578063a631571e14610361578063afcb95d714610381578063b1dc65a4146103a157600080fd5b806381ff7048146102b957806385b214cf146103265780638da5cb5b1461033957600080fd5b806366316d8d116101455780637f15e1661161011f5780637f15e16614610289578063814118341461029c57806381f1b938146102b157600080fd5b806366316d8d1461026657806379ba5097146102795780637d4807871461028157600080fd5b806326045ab31161017657806326045ab3146101f95780632a905ccc1461020c57806359b5b7ac1461022e57600080fd5b8063083a546614610192578063181f5a77146101a7575b600080fd5b6101a56101a036600461360b565b610595565b005b6101e36040518060400160405280601c81526020017f46756e6374696f6e7320436f6f7264696e61746f722076312e312e300000000081525081565b6040516101f091906136b1565b60405180910390f35b6101a5610207366004613814565b6105ea565b61021461082a565b60405168ffffffffffffffffff90911681526020016101f0565b61021461023c366004613958565b5060085474010000000000000000000000000000000000000000900468ffffffffffffffffff1690565b6101a56102743660046139e7565b6108c0565b6101a5610a79565b6101a5610b7b565b6101a561029736600461360b565b610d7b565b6102a4610dcb565b6040516101f09190613a71565b6101e3610e3a565b61030360015460025463ffffffff74010000000000000000000000000000000000000000830481169378010000000000000000000000000000000000000000000000009093041691565b6040805163ffffffff9485168152939092166020840152908201526060016101f0565b6101a5610334366004613a84565b610f0b565b60005460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101f0565b61037461036f366004613a9d565b610fbd565b6040516101f09190613bf2565b6040805160018152600060208201819052918101919091526060016101f0565b6101a56103af366004613c46565b61115e565b6105146040805161012081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e0810182905261010081019190915250604080516101208101825260085463ffffffff80821683526401000000008204811660208401526801000000000000000082048116938301939093526c0100000000000000000000000081048316606083015270010000000000000000000000000000000081048316608083015274010000000000000000000000000000000000000000810468ffffffffffffffffff1660a08301527d010000000000000000000000000000000000000000000000000000000000900461ffff1660c08201526009547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff811660e08301527c0100000000000000000000000000000000000000000000000000000000900490911661010082015290565b6040516101f09190613cfd565b61053461052f366004613de0565b611815565b6040516bffffffffffffffffffffffff90911681526020016101f0565b6101e3611977565b6101a5610567366004613ef9565b6119ce565b6105746123fa565b6040519081526020016101f0565b6101a5610590366004613fc6565b612648565b61059d61265c565b60008190036105d8576040517f4f42be3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600d6105e582848361407c565b505050565b6105f26126df565b80516008805460208401516040808601516060870151608088015160a089015160c08a015161ffff167d010000000000000000000000000000000000000000000000000000000000027fff0000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffff68ffffffffffffffffff90921674010000000000000000000000000000000000000000027fffffff000000000000000000ffffffffffffffffffffffffffffffffffffffff63ffffffff94851670010000000000000000000000000000000002167fffffff00000000000000000000000000ffffffffffffffffffffffffffffffff9585166c01000000000000000000000000027fffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffffff9786166801000000000000000002979097167fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff998616640100000000027fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000909b169c86169c909c1799909917979097169990991793909317919091169490941792909217919091169390931790915560e08301516101008401519091167c0100000000000000000000000000000000000000000000000000000000027bffffffffffffffffffffffffffffffffffffffffffffffffffffffff90911617600955517fdfcfddac88b9a703c9943a696e4da4a2d76b5b5c285a821bbcc54052775ecd2b9061081f908390613cfd565b60405180910390a150565b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16632a905ccc6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610897573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108bb91906141a2565b905090565b6108c86126e7565b806bffffffffffffffffffffffff166000036109025750336000908152600a60205260409020546bffffffffffffffffffffffff1661095c565b336000908152600a60205260409020546bffffffffffffffffffffffff8083169116101561095c576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b336000908152600a6020526040812080548392906109899084906bffffffffffffffffffffffff166141ee565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055506109de7f000000000000000000000000000000000000000000000000000000000000000090565b6040517f66316d8d00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff84811660048301526bffffffffffffffffffffffff8416602483015291909116906366316d8d90604401600060405180830381600087803b158015610a5d57600080fd5b505af1158015610a71573d6000803e3d6000fd5b505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610aff576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610b836126df565b610b8b6126e7565b6000610b95610dcb565b905060005b8151811015610d77576000600a6000848481518110610bbb57610bbb614213565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252810191909152604001600020546bffffffffffffffffffffffff1690508015610d66576000600a6000858581518110610c1a57610c1a614213565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550610cb17f000000000000000000000000000000000000000000000000000000000000000090565b73ffffffffffffffffffffffffffffffffffffffff166366316d8d848481518110610cde57610cde614213565b6020026020010151836040518363ffffffff1660e01b8152600401610d3392919073ffffffffffffffffffffffffffffffffffffffff9290921682526bffffffffffffffffffffffff16602082015260400190565b600060405180830381600087803b158015610d4d57600080fd5b505af1158015610d61573d6000803e3d6000fd5b505050505b50610d7081614242565b9050610b9a565b5050565b610d8361265c565b6000819003610dbe576040517f4f42be3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600c6105e582848361407c565b60606006805480602002602001604051908101604052809291908181526020018280548015610e3057602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610e05575b5050505050905090565b6060600d8054610e4990613fe3565b9050600003610e84576040517f4f42be3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600d8054610e9190613fe3565b80601f0160208091040260200160405190810160405280929190818152602001828054610ebd90613fe3565b8015610e305780601f10610edf57610100808354040283529160200191610e30565b820191906000526020600020905b815481529060010190602001808311610eed57509395945050505050565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610f7a576040517fc41a5b0900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008181526007602052604080822091909155517f8a4b97add3359bd6bcf5e82874363670eb5ad0f7615abddbd0ed0a3a98f0f4169061081f9083815260200190565b6040805161016081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e08101829052610100810182905261012081018290526101408101919091523373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614611085576040517fc41a5b0900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6110966110918361427a565b612892565b90506110a86060830160408401613fc6565b815173ffffffffffffffffffffffffffffffffffffffff91909116907fbf50768ccf13bd0110ca6d53a9c4f1f3271abdd4c24a56878863ed25b20598ff326110f660c0870160a08801614367565b61110861016088016101408901613fc6565b6111128880614384565b6111246101208b016101008c016143e9565b60208b013561113a6101008d0160e08e01614404565b8b60405161115099989796959493929190614421565b60405180910390a35b919050565b60005a604080518b3580825262ffffff6020808f0135600881901c929092169084015293945092917fb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62910160405180910390a16040805160608101825260025480825260035460ff80821660208501526101009091041692820192909252908314611245576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f636f6e666967446967657374206d69736d6174636800000000000000000000006044820152606401610af6565b6112538b8b8b8b8b8b612cbd565b60007f0000000000000000000000000000000000000000000000000000000000000000156112b05760028260200151836040015161129191906144c9565b61129b9190614511565b6112a69060016144c9565b60ff1690506112c6565b60208201516112c09060016144c9565b60ff1690505b88811461132f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f77726f6e67206e756d626572206f66207369676e6174757265730000000000006044820152606401610af6565b888714611398576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f7369676e617475726573206f7574206f6620726567697374726174696f6e00006044820152606401610af6565b3360009081526004602090815260408083208151808301909252805460ff808216845292939192918401916101009091041660028111156113db576113db614533565b60028111156113ec576113ec614533565b905250905060028160200151600281111561140957611409614533565b14801561145057506006816000015160ff168154811061142b5761142b614213565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff1633145b6114b6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f756e617574686f72697a6564207472616e736d697474657200000000000000006044820152606401610af6565b50505050506114c36135a3565b6000808a8a6040516114d6929190614562565b6040519081900381206114ed918e90602001614572565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181528282528051602091820120838301909252600080845290830152915060005b898110156117f757600060018489846020811061155657611556614213565b61156391901a601b6144c9565b8e8e8681811061157557611575614213565b905060200201358d8d8781811061158e5761158e614213565b90506020020135604051600081526020016040526040516115cb949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa1580156115ed573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff811660009081526004602090815290849020838501909452835460ff8082168552929650929450840191610100900416600281111561166d5761166d614533565b600281111561167e5761167e614533565b905250925060018360200151600281111561169b5761169b614533565b14611702576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f61646472657373206e6f7420617574686f72697a656420746f207369676e00006044820152606401610af6565b8251600090879060ff16601f811061171c5761171c614213565b602002015173ffffffffffffffffffffffffffffffffffffffff161461179e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f6e6f6e2d756e69717565207369676e61747572650000000000000000000000006044820152606401610af6565b8086846000015160ff16601f81106117b8576117b8614213565b73ffffffffffffffffffffffffffffffffffffffff90921660209290920201526117e36001866144c9565b945050806117f090614242565b9050611537565b505050611808833383858e8e612d74565b5050505050505050505050565b60007f00000000000000000000000000000000000000000000000000000000000000006040517f10fc49c100000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8816600482015263ffffffff8516602482015273ffffffffffffffffffffffffffffffffffffffff91909116906310fc49c19060440160006040518083038186803b1580156118b557600080fd5b505afa1580156118c9573d6000803e3d6000fd5b5050505066038d7ea4c6800082111561190e576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600061191861082a565b9050600061195b87878080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061023c92505050565b905061196985858385612f42565b925050505b95945050505050565b6060600c805461198690613fe3565b90506000036119c1576040517f4f42be3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600c8054610e9190613fe3565b855185518560ff16601f831115611a41576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f746f6f206d616e79207369676e657273000000000000000000000000000000006044820152606401610af6565b80600003611aab576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f66206d75737420626520706f73697469766500000000000000000000000000006044820152606401610af6565b818314611b39576040517f89a61989000000000000000000000000000000000000000000000000000000008152602060048201526024808201527f6f7261636c6520616464726573736573206f7574206f6620726567697374726160448201527f74696f6e000000000000000000000000000000000000000000000000000000006064820152608401610af6565b611b44816003614586565b8311611bac576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f6661756c74792d6f7261636c65206620746f6f206869676800000000000000006044820152606401610af6565b611bb461265c565b6040805160c0810182528a8152602081018a905260ff89169181018290526060810188905267ffffffffffffffff8716608082015260a0810186905290611bfb9088613088565b60055415611db057600554600090611c159060019061459d565b9050600060058281548110611c2c57611c2c614213565b60009182526020822001546006805473ffffffffffffffffffffffffffffffffffffffff90921693509084908110611c6657611c66614213565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff85811684526004909252604080842080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000090811690915592909116808452922080549091169055600580549192509080611ce657611ce66145b0565b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff00000000000000000000000000000000000000001690550190556006805480611d4f57611d4f6145b0565b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff000000000000000000000000000000000000000016905501905550611bfb915050565b60005b8151518110156122175760006004600084600001518481518110611dd957611dd9614213565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff16825281019190915260400160002054610100900460ff166002811115611e2357611e23614533565b14611e8a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265706561746564207369676e657220616464726573730000000000000000006044820152606401610af6565b6040805180820190915260ff82168152600160208201528251805160049160009185908110611ebb57611ebb614213565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040016000208251815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082168117835592840151919283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001617610100836002811115611f5c57611f5c614533565b021790555060009150611f6c9050565b6004600084602001518481518110611f8657611f86614213565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff16825281019190915260400160002054610100900460ff166002811115611fd057611fd0614533565b14612037576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f7265706561746564207472616e736d69747465722061646472657373000000006044820152606401610af6565b6040805180820190915260ff82168152602081016002815250600460008460200151848151811061206a5761206a614213565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040016000208251815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082168117835592840151919283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000161761010083600281111561210b5761210b614533565b02179055505082518051600592508390811061212957612129614213565b602090810291909101810151825460018101845560009384529282902090920180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff90931692909217909155820151805160069190839081106121a5576121a5614213565b60209081029190910181015182546001810184556000938452919092200180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff9092169190911790558061220f81614242565b915050611db3565b506040810151600380547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff909216919091179055600180547fffffffff00000000ffffffffffffffffffffffffffffffffffffffffffffffff8116780100000000000000000000000000000000000000000000000063ffffffff43811682029290921780855592048116929182916014916122cf918491740100000000000000000000000000000000000000009004166145df565b92506101000a81548163ffffffff021916908363ffffffff16021790555061232e4630600160149054906101000a900463ffffffff1663ffffffff16856000015186602001518760400151886060015189608001518a60a001516130a1565b600281905582518051600380547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1661010060ff9093169290920291909117905560015460208501516040808701516060880151608089015160a08a015193517f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05986123e5988b9891977401000000000000000000000000000000000000000090920463ffffffff169690959194919391926145fc565b60405180910390a15050505050505050505050565b604080516101208101825260085463ffffffff80821683526401000000008204811660208401526801000000000000000082048116838501526c010000000000000000000000008083048216606085015270010000000000000000000000000000000083048216608085015274010000000000000000000000000000000000000000830468ffffffffffffffffff1660a0808601919091527d01000000000000000000000000000000000000000000000000000000000090930461ffff1660c08501526009547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff811660e08601527c01000000000000000000000000000000000000000000000000000000009004909116610100840152600b5484517ffeaf968c00000000000000000000000000000000000000000000000000000000815294516000958694859490930473ffffffffffffffffffffffffffffffffffffffff169263feaf968c926004808401938290030181865afa15801561257d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906125a191906146ac565b5093505092505080426125b4919061459d565b836020015163ffffffff161080156125d657506000836020015163ffffffff16115b1561260457505060e001517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16919050565b60008213612641576040517f43d4cf6600000000000000000000000000000000000000000000000000000000815260048101839052602401610af6565b5092915050565b61265061265c565b6126598161314c565b50565b60005473ffffffffffffffffffffffffffffffffffffffff1633146126dd576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610af6565b565b6126dd61265c565b600b546bffffffffffffffffffffffff1660000361270157565b600061270b610dcb565b90508051600003612748576040517f30274b3a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8051600b54600091612767916bffffffffffffffffffffffff166146fc565b905060005b82518110156128335781600a600085848151811061278c5761278c614213565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282829054906101000a90046bffffffffffffffffffffffff166127f49190614727565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055508061282c90614242565b905061276c565b508151612840908261474c565b600b80546000906128609084906bffffffffffffffffffffffff166141ee565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055505050565b6040805161016081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e0810182905261010081018290526101208101829052610140810191909152604080516101208101825260085463ffffffff80821683526401000000008204811660208401526801000000000000000082048116938301939093526c0100000000000000000000000081048316606083015270010000000000000000000000000000000081048316608083015268ffffffffffffffffff7401000000000000000000000000000000000000000082041660a083015261ffff7d010000000000000000000000000000000000000000000000000000000000909104811660c083018190526009547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff811660e08501527c0100000000000000000000000000000000000000000000000000000000900490931661010080840191909152850151919291161115612a42576040517fdada758700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60085460009074010000000000000000000000000000000000000000900468ffffffffffffffffff1690506000612a838560e001513a848860800151612f42565b9050806bffffffffffffffffffffffff1685606001516bffffffffffffffffffffffff161015612adf576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000612b623087604001518860a001518960c001516001612b009190614774565b6040805173ffffffffffffffffffffffffffffffffffffffff958616602080830191909152949095168582015267ffffffffffffffff928316606086015291166080808501919091528151808503909101815260a09093019052815191012090565b90506040518061016001604052808281526020013073ffffffffffffffffffffffffffffffffffffffff168152602001836bffffffffffffffffffffffff168152602001876040015173ffffffffffffffffffffffffffffffffffffffff1681526020018760a0015167ffffffffffffffff1681526020018760e0015163ffffffff168152602001876080015168ffffffffffffffffff1681526020018468ffffffffffffffffff168152602001856040015163ffffffff1664ffffffffff168152602001856060015163ffffffff1664ffffffffff168152602001856080015163ffffffff1642612c549190614795565b63ffffffff16815250945084604051602001612c709190613bf2565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529181528151602092830120600093845260079092529091205550919392505050565b6000612cca826020614586565b612cd5856020614586565b612ce188610144614795565b612ceb9190614795565b612cf59190614795565b612d00906000614795565b9050368114612d6b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f63616c6c64617461206c656e677468206d69736d6174636800000000000000006044820152606401610af6565b50505050505050565b606080808080612d8686880188614883565b8451949950929750909550935091501580612da357508351855114155b80612db057508251855114155b80612dbd57508151855114155b80612dca57508051855114155b15612e01576040517f0be3632800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b8551811015612f34576000612e99878381518110612e2457612e24614213565b6020026020010151878481518110612e3e57612e3e614213565b6020026020010151878581518110612e5857612e58614213565b6020026020010151878681518110612e7257612e72614213565b6020026020010151878781518110612e8c57612e8c614213565b6020026020010151613241565b90506000816006811115612eaf57612eaf614533565b1480612ecc57506001816006811115612eca57612eca614533565b145b15612f2357868281518110612ee357612ee3614213565b60209081029190910181015160405133815290917fc708e0440951fd63499c0f7a73819b469ee5dd3ecc356c0ab4eb7f18389009d9910160405180910390a25b50612f2d81614242565b9050612e04565b505050505050505050505050565b60085460009081908690612f7a9063ffffffff6c010000000000000000000000008204811691680100000000000000009004166145df565b612f8491906145df565b60095463ffffffff91821692507c0100000000000000000000000000000000000000000000000000000000900416851015612fe5576009547c0100000000000000000000000000000000000000000000000000000000900463ffffffff1694505b60085460009061271090612fff9063ffffffff1688614586565b6130099190614955565b6130139087614795565b90506000613020826134d2565b9050600061303c846bffffffffffffffffffffffff8416614586565b9050600061305868ffffffffffffffffff808916908a16614727565b905061307a6130756bffffffffffffffffffffffff831684614795565b613501565b9a9950505050505050505050565b6000613092610dcb565b511115610d7757610d776126e7565b6000808a8a8a8a8a8a8a8a8a6040516020016130c599989796959493929190614969565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179150509998505050505050505050565b3373ffffffffffffffffffffffffffffffffffffffff8216036131cb576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610af6565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600080838060200190518101906132589190614a3f565b60008881526007602052604090205490915061327857600291505061196e565b806040516020016132899190613bf2565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152918152815160209283012060008a81526007909352912054146132db57600691505061196e565b60006132e63a6134d2565b905060008261012001518361010001516133009190614b07565b6133119064ffffffffff168361474c565b90506000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663330605298b8b878960e0015168ffffffffffffffffff16886133709190614727565b338b6040518763ffffffff1660e01b815260040161339396959493929190614b25565b60408051808303816000875af11580156133b1573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906133d59190614ba1565b909250905060008260068111156133ee576133ee614533565b148061340b5750600182600681111561340957613409614533565b145b156134c45760008b8152600760205260408120556134298184614727565b336000908152600a6020526040812080547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff93841617905560e0870151600b805468ffffffffffffffffff9092169390929161349591859116614727565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055505b509998505050505050505050565b60006134fb6134df6123fa565b6134f184670de0b6b3a7640000614586565b6130759190614955565b92915050565b60006bffffffffffffffffffffffff82111561359f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f53616665436173743a2076616c756520646f65736e27742066697420696e203960448201527f36206269747300000000000000000000000000000000000000000000000000006064820152608401610af6565b5090565b604051806103e00160405280601f906020820280368337509192915050565b60008083601f8401126135d457600080fd5b50813567ffffffffffffffff8111156135ec57600080fd5b60208301915083602082850101111561360457600080fd5b9250929050565b6000806020838503121561361e57600080fd5b823567ffffffffffffffff81111561363557600080fd5b613641858286016135c2565b90969095509350505050565b6000815180845260005b8181101561367357602081850181015186830182015201613657565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b6020815260006136c4602083018461364d565b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610120810167ffffffffffffffff8111828210171561371e5761371e6136cb565b60405290565b604051610160810167ffffffffffffffff8111828210171561371e5761371e6136cb565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561378f5761378f6136cb565b604052919050565b63ffffffff8116811461265957600080fd5b803561115981613797565b68ffffffffffffffffff8116811461265957600080fd5b8035611159816137b4565b803561ffff8116811461115957600080fd5b80357bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8116811461115957600080fd5b6000610120828403121561382757600080fd5b61382f6136fa565b613838836137a9565b8152613846602084016137a9565b6020820152613857604084016137a9565b6040820152613868606084016137a9565b6060820152613879608084016137a9565b608082015261388a60a084016137cb565b60a082015261389b60c084016137d6565b60c08201526138ac60e084016137e8565b60e08201526101006138bf8185016137a9565b908201529392505050565b600082601f8301126138db57600080fd5b813567ffffffffffffffff8111156138f5576138f56136cb565b61392660207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601613748565b81815284602083860101111561393b57600080fd5b816020850160208301376000918101602001919091529392505050565b60006020828403121561396a57600080fd5b813567ffffffffffffffff81111561398157600080fd5b61398d848285016138ca565b949350505050565b73ffffffffffffffffffffffffffffffffffffffff8116811461265957600080fd5b803561115981613995565b6bffffffffffffffffffffffff8116811461265957600080fd5b8035611159816139c2565b600080604083850312156139fa57600080fd5b8235613a0581613995565b91506020830135613a15816139c2565b809150509250929050565b600081518084526020808501945080840160005b83811015613a6657815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101613a34565b509495945050505050565b6020815260006136c46020830184613a20565b600060208284031215613a9657600080fd5b5035919050565b600060208284031215613aaf57600080fd5b813567ffffffffffffffff811115613ac657600080fd5b820161016081850312156136c457600080fd5b805182526020810151613b04602084018273ffffffffffffffffffffffffffffffffffffffff169052565b506040810151613b2460408401826bffffffffffffffffffffffff169052565b506060810151613b4c606084018273ffffffffffffffffffffffffffffffffffffffff169052565b506080810151613b68608084018267ffffffffffffffff169052565b5060a0810151613b8060a084018263ffffffff169052565b5060c0810151613b9d60c084018268ffffffffffffffffff169052565b5060e0810151613bba60e084018268ffffffffffffffffff169052565b506101008181015164ffffffffff9081169184019190915261012080830151909116908301526101409081015163ffffffff16910152565b61016081016134fb8284613ad9565b60008083601f840112613c1357600080fd5b50813567ffffffffffffffff811115613c2b57600080fd5b6020830191508360208260051b850101111561360457600080fd5b60008060008060008060008060e0898b031215613c6257600080fd5b606089018a811115613c7357600080fd5b8998503567ffffffffffffffff80821115613c8d57600080fd5b613c998c838d016135c2565b909950975060808b0135915080821115613cb257600080fd5b613cbe8c838d01613c01565b909750955060a08b0135915080821115613cd757600080fd5b50613ce48b828c01613c01565b999c989b50969995989497949560c00135949350505050565b815163ffffffff9081168252602080840151821690830152604080840151821690830152606080840151821690830152608080840151918216908301526101208201905060a0830151613d5d60a084018268ffffffffffffffffff169052565b5060c0830151613d7360c084018261ffff169052565b5060e0830151613da360e08401827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff169052565b506101008381015163ffffffff8116848301525b505092915050565b67ffffffffffffffff8116811461265957600080fd5b803561115981613dbf565b600080600080600060808688031215613df857600080fd5b8535613e0381613dbf565b9450602086013567ffffffffffffffff811115613e1f57600080fd5b613e2b888289016135c2565b9095509350506040860135613e3f81613797565b949793965091946060013592915050565b600067ffffffffffffffff821115613e6a57613e6a6136cb565b5060051b60200190565b600082601f830112613e8557600080fd5b81356020613e9a613e9583613e50565b613748565b82815260059290921b84018101918181019086841115613eb957600080fd5b8286015b84811015613edd578035613ed081613995565b8352918301918301613ebd565b509695505050505050565b803560ff8116811461115957600080fd5b60008060008060008060c08789031215613f1257600080fd5b863567ffffffffffffffff80821115613f2a57600080fd5b613f368a838b01613e74565b97506020890135915080821115613f4c57600080fd5b613f588a838b01613e74565b9650613f6660408a01613ee8565b95506060890135915080821115613f7c57600080fd5b613f888a838b016138ca565b9450613f9660808a01613dd5565b935060a0890135915080821115613fac57600080fd5b50613fb989828a016138ca565b9150509295509295509295565b600060208284031215613fd857600080fd5b81356136c481613995565b600181811c90821680613ff757607f821691505b602082108103614030577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f8211156105e557600081815260208120601f850160051c8101602086101561405d5750805b601f850160051c820191505b81811015610a7157828155600101614069565b67ffffffffffffffff831115614094576140946136cb565b6140a8836140a28354613fe3565b83614036565b6000601f8411600181146140fa57600085156140c45750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b178355614190565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b828110156141495786850135825560209485019460019092019101614129565b5086821015614184577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b8051611159816137b4565b6000602082840312156141b457600080fd5b81516136c4816137b4565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6bffffffffffffffffffffffff828116828216039080821115612641576126416141bf565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203614273576142736141bf565b5060010190565b6000610160823603121561428d57600080fd5b614295613724565b823567ffffffffffffffff8111156142ac57600080fd5b6142b8368286016138ca565b825250602083013560208201526142d1604084016139b7565b60408201526142e2606084016139dc565b60608201526142f3608084016137cb565b608082015261430460a08401613dd5565b60a082015261431560c08401613dd5565b60c082015261432660e084016137a9565b60e08201526101006143398185016137d6565b9082015261012061434b848201613dd5565b9082015261014061435d8482016139b7565b9082015292915050565b60006020828403121561437957600080fd5b81356136c481613dbf565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18436030181126143b957600080fd5b83018035915067ffffffffffffffff8211156143d457600080fd5b60200191503681900382131561360457600080fd5b6000602082840312156143fb57600080fd5b6136c4826137d6565b60006020828403121561441657600080fd5b81356136c481613797565b73ffffffffffffffffffffffffffffffffffffffff8a8116825267ffffffffffffffff8a166020830152881660408201526102406060820181905281018690526000610260878982850137600083890182015261ffff8716608084015260a0830186905263ffffffff851660c0840152601f88017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016830101905061307a60e0830184613ad9565b60ff81811683821601908111156134fb576134fb6141bf565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b600060ff831680614524576145246144e2565b8060ff84160491505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b8183823760009101908152919050565b828152606082602083013760800192915050565b80820281158282048414176134fb576134fb6141bf565b818103818111156134fb576134fb6141bf565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b63ffffffff818116838216019080821115612641576126416141bf565b600061012063ffffffff808d1684528b6020850152808b1660408501525080606084015261462c8184018a613a20565b905082810360808401526146408189613a20565b905060ff871660a084015282810360c084015261465d818761364d565b905067ffffffffffffffff851660e0840152828103610100840152614682818561364d565b9c9b505050505050505050505050565b805169ffffffffffffffffffff8116811461115957600080fd5b600080600080600060a086880312156146c457600080fd5b6146cd86614692565b94506020860151935060408601519250606086015191506146f060808701614692565b90509295509295909350565b60006bffffffffffffffffffffffff8084168061471b5761471b6144e2565b92169190910492915050565b6bffffffffffffffffffffffff818116838216019080821115612641576126416141bf565b6bffffffffffffffffffffffff818116838216028082169190828114613db757613db76141bf565b67ffffffffffffffff818116838216019080821115612641576126416141bf565b808201808211156134fb576134fb6141bf565b600082601f8301126147b957600080fd5b813560206147c9613e9583613e50565b82815260059290921b840181019181810190868411156147e857600080fd5b8286015b84811015613edd57803583529183019183016147ec565b600082601f83011261481457600080fd5b81356020614824613e9583613e50565b82815260059290921b8401810191818101908684111561484357600080fd5b8286015b84811015613edd57803567ffffffffffffffff8111156148675760008081fd5b6148758986838b01016138ca565b845250918301918301614847565b600080600080600060a0868803121561489b57600080fd5b853567ffffffffffffffff808211156148b357600080fd5b6148bf89838a016147a8565b965060208801359150808211156148d557600080fd5b6148e189838a01614803565b955060408801359150808211156148f757600080fd5b61490389838a01614803565b9450606088013591508082111561491957600080fd5b61492589838a01614803565b9350608088013591508082111561493b57600080fd5b5061494888828901614803565b9150509295509295909350565b600082614964576149646144e2565b500490565b60006101208b835273ffffffffffffffffffffffffffffffffffffffff8b16602084015267ffffffffffffffff808b1660408501528160608501526149b08285018b613a20565b915083820360808501526149c4828a613a20565b915060ff881660a085015283820360c08501526149e1828861364d565b90861660e08501528381036101008501529050614682818561364d565b805161115981613995565b8051611159816139c2565b805161115981613dbf565b805161115981613797565b805164ffffffffff8116811461115957600080fd5b60006101608284031215614a5257600080fd5b614a5a613724565b82518152614a6a602084016149fe565b6020820152614a7b60408401614a09565b6040820152614a8c606084016149fe565b6060820152614a9d60808401614a14565b6080820152614aae60a08401614a1f565b60a0820152614abf60c08401614197565b60c0820152614ad060e08401614197565b60e0820152610100614ae3818501614a2a565b90820152610120614af5848201614a2a565b908201526101406138bf848201614a1f565b64ffffffffff818116838216019080821115612641576126416141bf565b6000610200808352614b398184018a61364d565b90508281036020840152614b4d818961364d565b6bffffffffffffffffffffffff88811660408601528716606085015273ffffffffffffffffffffffffffffffffffffffff861660808501529150614b96905060a0830184613ad9565b979650505050505050565b60008060408385031215614bb457600080fd5b825160078110614bc357600080fd5b6020840151909250613a15816139c256fea164736f6c6343000813000a",
}

var FunctionsCoordinatorABI = FunctionsCoordinatorMetaData.ABI

var FunctionsCoordinatorBin = FunctionsCoordinatorMetaData.Bin

func DeployFunctionsCoordinator(auth *bind.TransactOpts, backend bind.ContractBackend, router common.Address, config FunctionsBillingConfig, linkToNativeFeed common.Address) (common.Address, *types.Transaction, *FunctionsCoordinator, error) {
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
	return common.HexToHash("0xdfcfddac88b9a703c9943a696e4da4a2d76b5b5c285a821bbcc54052775ecd2b")
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

func (FunctionsCoordinatorTransmitted) Topic() common.Hash {
	return common.HexToHash("0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62")
}

func (_FunctionsCoordinator *FunctionsCoordinator) Address() common.Address {
	return _FunctionsCoordinator.address
}

type FunctionsCoordinatorInterface interface {
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

	FilterTransmitted(opts *bind.FilterOpts) (*FunctionsCoordinatorTransmittedIterator, error)

	WatchTransmitted(opts *bind.WatchOpts, sink chan<- *FunctionsCoordinatorTransmitted) (event.Subscription, error)

	ParseTransmitted(log types.Log) (*FunctionsCoordinatorTransmitted, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
