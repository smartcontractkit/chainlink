// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package functions_router

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

var FunctionsRouterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"timelockBlocks\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"maximumTimelockBlocks\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"linkToken\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"internalBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"externalBalance\",\"type\":\"uint256\"}],\"name\":\"BalanceInvariantViolated\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConsumerRequestsInFlight\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidConfigData\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"InvalidConsumer\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidProposal\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSubscription\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposedOwner\",\"type\":\"address\"}],\"name\":\"MustBeRequestedOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"MustBeSubOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByRoute\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableFromCoordinator\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableFromLink\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingRequestExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ProposedTimelockAboveMaximum\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Reentrant\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"ReservedLabel\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"RouteNotFound\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TimelockInEffect\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyConsumers\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"fromHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"toBytes\",\"type\":\"bytes\"}],\"name\":\"ConfigProposed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"adminFee\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"bytes4\",\"name\":\"handleOracleFulfillmentSelector\",\"type\":\"bytes4\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"fromHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"toBytes\",\"type\":\"bytes\"}],\"name\":\"ConfigUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"proposedContractSetId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"proposedContractSetFromAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"proposedContractSetToAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timelockEndBlock\",\"type\":\"uint256\"}],\"name\":\"ContractProposed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"proposedContractSetId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"proposedContractSetFromAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"proposedContractSetToAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"major\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"minor\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"patch\",\"type\":\"uint16\"}],\"name\":\"ContractUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalCostJuels\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"resultCode\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"response\",\"type\":\"bytes\"}],\"name\":\"RequestEnd\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"donId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"subscriptionOwner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requestingContract\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requestInitiator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"dataVersion\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"}],\"name\":\"RequestStart\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"fundsRecipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fundsAmount\",\"type\":\"uint256\"}],\"name\":\"SubscriptionCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"SubscriptionCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"from\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"to\",\"type\":\"uint16\"}],\"name\":\"TimeLockProposed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"from\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"to\",\"type\":\"uint16\"}],\"name\":\"TimeLockUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"MAX_CONSUMERS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"acceptSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"addConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"cancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createSubscription\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"response\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"juelsPerGas\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"costWithoutFulfillment\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"fulfill\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"resultCode\",\"type\":\"uint8\"},{\"internalType\":\"uint96\",\"name\":\"callbackGasCostJuels\",\"type\":\"uint96\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAdminFee\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfigHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"config\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"getConsumer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"initiatedRequests\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"completedRequests\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"useProposed\",\"type\":\"bool\"}],\"name\":\"getContractById\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"routeDestination\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"getContractById\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"routeDestination\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getProposedContractSet\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"},{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"getSubscription\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"blockedBalance\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"requestedOwner\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSubscriptionCount\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTotalBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isPaused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"oracleWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"ownerCancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"pendingRequestExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"name\":\"proposeConfigUpdate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"proposedContractSetIds\",\"type\":\"bytes32[]\"},{\"internalType\":\"address[]\",\"name\":\"proposedContractSetFromAddresses\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"proposedContractSetToAddresses\",\"type\":\"address[]\"}],\"name\":\"proposeContractsUpdate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"blocks\",\"type\":\"uint16\"}],\"name\":\"proposeTimelockBlocks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"removeConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"requestSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint16\",\"name\":\"dataVersion\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"donId\",\"type\":\"bytes32\"}],\"name\":\"sendRequest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"requestIdsToTimeout\",\"type\":\"bytes32[]\"}],\"name\":\"timeoutRequests\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"togglePaused\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"updateConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"updateContracts\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"updateTimelockBlocks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"validateProposedContracts\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60806040526001805463ffffffff60a01b191690553480156200002157600080fd5b50604051620063f0380380620063f083398101604081905262000044916200030d565b6000805460ff1916815582903390869086908590849081620000ad5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0380851661010002610100600160a81b031990921691909117909155811615620000e757620000e78162000198565b50506008805461ffff84811661ffff1991871662010000029190911663ffffffff19909216919091171790556000805260026020527fac33ff75c19e70fe83507db0d683fd3465c996598dc972688b7ace676c89077b80546001600160a01b0319163017905562000158816200024a565b8051602090910120600b555050600c80546001600160a01b0390931661010002610100600160a81b031990931692909217909155506200049c9350505050565b6001600160a01b038116331415620001f35760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401620000a4565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929361010090910416917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000808280602001905181019062000263919062000430565b6040805180820182526001600160601b0384168082526001600160e01b031984166020928301819052601280546001600160801b03191683176c0100000000000000000000000060e088901c021790558351918252918101919091529294509092507fb2b9bfb736b0bfef266985352cfca2a114226362209f5a6e5a2a2ed14383ba57910160405180910390a1505050565b805161ffff811681146200030857600080fd5b919050565b600080600080608085870312156200032457600080fd5b6200032f85620002f5565b9350602062000340818701620002f5565b60408701519094506001600160a01b03811681146200035e57600080fd5b60608701519093506001600160401b03808211156200037c57600080fd5b818801915088601f8301126200039157600080fd5b815181811115620003a657620003a662000486565b604051601f8201601f19908116603f01168101908382118183101715620003d157620003d162000486565b816040528281528b86848701011115620003ea57600080fd5b600093505b828410156200040e5784840186015181850187015292850192620003ef565b82841115620004205760008684830101525b989b979a50959850505050505050565b600080604083850312156200044457600080fd5b82516001600160601b03811681146200045c57600080fd5b60208401519092506001600160e01b0319811681146200047b57600080fd5b809150509250929050565b634e487b7160e01b600052604160045260246000fd5b615f4480620004ac6000396000f3fe608060405234801561001057600080fd5b50600436106102ad5760003560e01c806379ba50971161017b578063b187bd26116100d8578063d7ae1d301161008c578063e82ad7d411610071578063e82ad7d4146106e7578063eb523d6c146106fa578063f2fde38b1461070d57600080fd5b8063d7ae1d30146106c1578063e72f6e30146106d457600080fd5b8063b734c0f4116100bd578063b734c0f41461068e578063badc3eb614610696578063c4614d61146106ae57600080fd5b8063b187bd2614610670578063b56438581461067b57600080fd5b80639f87fad71161012f578063a47c769611610114578063a47c769614610626578063a4c0ed361461064a578063a9c9a9181461065d57600080fd5b80639f87fad71461060b578063a21a23e41461061e57600080fd5b80638da5cb5b116101605780638da5cb5b146105d85780638fde5317146105fb5780639883c10d1461060357600080fd5b806379ba5097146105bd57806382359740146105c557600080fd5b8063461d27621161022957806366419970116101dd578063674603d0116101c2578063674603d0146104dd5780636a6df79b146105725780637341c10c146105aa57600080fd5b8063664199701461048a578063665871ec146104ca57600080fd5b80635c975abb1161020e5780635c975abb1461044557806364d51a2a1461045c57806366316d8d1461047757600080fd5b8063461d2762146103cd57806354fd4d50146103ee57600080fd5b80631e9fb8cd1161028057806336566f061161026557806336566f0614610379578063385de9ae146103815780633fd67e051461039457600080fd5b80631e9fb8cd146103505780632a905ccc1461036357600080fd5b806302bcc5b6146102b257806304c357cb146102c757806312b58349146102da578063181f5a771461030e575b600080fd5b6102c56102c03660046155dc565b610720565b005b6102c56102d53660046155f7565b6107a0565b600d546bffffffffffffffffffffffff165b6040516bffffffffffffffffffffffff90911681526020015b60405180910390f35b60408051808201909152601381527f46756e6374696f6e7320526f757465722076310000000000000000000000000060208201525b6040516103059190615a35565b6102c561035e366004615339565b61098e565b6012546bffffffffffffffffffffffff166102ec565b6102c5610b1d565b6102c561038f366004615491565b610b42565b6103a76103a23660046154dd565b610d8e565b6040805160ff90931683526bffffffffffffffffffffffff909116602083015201610305565b6103e06103db366004615623565b61106f565b604051908152602001610305565b600180546040805192835261ffff7401000000000000000000000000000000000000000083048116602085015276010000000000000000000000000000000000000000000090920490911690820152606001610305565b60005460ff165b6040519015158152602001610305565b610464606481565b60405161ffff9091168152602001610305565b6102c561048536600461528b565b61110b565b600c547501000000000000000000000000000000000000000000900467ffffffffffffffff165b60405167ffffffffffffffff9091168152602001610305565b6102c56104d83660046152c4565b611378565b61054a6104eb366004615256565b73ffffffffffffffffffffffffffffffffffffffff919091166000908152600f6020908152604080832067ffffffffffffffff948516845290915290205460ff8116926101008204831692690100000000000000000090920490911690565b60408051931515845267ffffffffffffffff9283166020850152911690820152606001610305565b61058561058036600461546c565b611861565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610305565b6102c56105b83660046155f7565b611874565b6102c5611b00565b6102c56105d33660046155dc565b611c22565b600054610100900473ffffffffffffffffffffffffffffffffffffffff16610585565b6102c5611e18565b600b546103e0565b6102c56106193660046155f7565b611e9c565b6104b16123a0565b6106396106343660046155dc565b6125d3565b604051610305959493929190615b94565b6102c56106583660046151fa565b6126c1565b61058561066b36600461543a565b61291d565b60005460ff1661044c565b6102c56106893660046155c1565b612930565b6102c5612a4a565b61069e612d3f565b6040516103059493929190615ad2565b6103436106bc366004615491565b612e8c565b6102c56106cf3660046155f7565b612ea1565b6102c56106e23660046151dd565b613083565b61044c6106f53660046155dc565b6132a9565b6102c561070836600461543a565b6133fa565b6102c561071b3660046151dd565b6136ae565b6107286136bf565b67ffffffffffffffff81166000908152600e602052604090206001015473ffffffffffffffffffffffffffffffffffffffff1680610792576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61079c8282613745565b5050565b67ffffffffffffffff82166000908152600e6020526040902060010154829073ffffffffffffffffffffffffffffffffffffffff168061080c576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff821614610878576040517fd8a3fb5200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821660048201526024015b60405180910390fd5b600c5460ff16156108b5576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff84166000908152600e602052604090206002015473ffffffffffffffffffffffffffffffffffffffff8481169116146109885767ffffffffffffffff84166000818152600e602090815260409182902060020180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff88169081179091558251338152918201527f69436ea6df009049404f564eff6622cd00522b0bd6a89efd9e52a355c4a879be91015b60405180910390a25b50505050565b6109966136bf565b6109a1838383613b60565b6008546000906109bb9062010000900461ffff1643615c82565b604080516080810182528681526020808201879052918101859052606081018390528651929350916003916109f4918391890190614ee5565b506020828101518051610a0d9260018501920190614f2c565b5060408201518051610a29916002840191602090910190614f2c565b506060820151816003015590505060005b84518160ff161015610b16577f72a33d2f293a0a70fad221bb610d3d6b52aed2d840adae1fa721071fbd290cfd858260ff1681518110610a7c57610a7c615e8f565b6020026020010151858360ff1681518110610a9957610a99615e8f565b6020026020010151858460ff1681518110610ab657610ab6615e8f565b602002602001015185604051610afc949392919093845273ffffffffffffffffffffffffffffffffffffffff928316602085015291166040830152606082015260800190565b60405180910390a180610b0e81615e11565b915050610a3a565b5050505050565b610b256136bf565b60005460ff1615610b3a57610b38613e4f565b565b610b38613ecc565b610b4a6136bf565b60008381526002602052604090205473ffffffffffffffffffffffffffffffffffffffff1680610ba9576040517f80833e330000000000000000000000000000000000000000000000000000000081526004810185905260240161086f565b600073ffffffffffffffffffffffffffffffffffffffff8216301415610bd25750600b54610c55565b8173ffffffffffffffffffffffffffffffffffffffff16639883c10d6040518163ffffffff1660e01b8152600401602060405180830381600087803b158015610c1a57600080fd5b505af1158015610c2e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c529190615453565b90505b8383604051610c65929190615913565b6040518091039020811415610ca6576040517fee03280800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b604051806060016040528082815260200185858080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250505090825250600854602090910190610d0d9062010000900461ffff1643615c82565b90526000868152600760209081526040909120825181558282015180519192610d3e92600185019290910190614fa6565b50604091820151600290910155517f0fcfd32a68209b42944376bc5f4bf72c41ba0f378cf60434f84390b82c9844cf90610d7f908790849088908890615987565b60405180910390a15050505050565b600c54600090819060ff1615610dd0576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600088815260116020908152604091829020825161010081018452815473ffffffffffffffffffffffffffffffffffffffff9081168252600183015490811693820184905274010000000000000000000000000000000000000000810467ffffffffffffffff16948201949094527c010000000000000000000000000000000000000000000000000000000090930463ffffffff16606084015260028101546bffffffffffffffffffffffff9081166080850152600382015460a0850152600482015460c08501526005909101541660e0830152610eb2576002925050611064565b805173ffffffffffffffffffffffffffffffffffffffff163314610f02576040517f8bec23e700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610f2081608001518260c001518360e0015184606001518a8a613f27565b60008a8152601160209081526040822080547fffffffffffffffffffffffff0000000000000000000000000000000000000000168155600181018390556002810180547fffffffffffffffffffffffffffffffffffffffff00000000000000000000000090811690915560038201849055600482018490556005909101805490911690556060840151908401519295509091610fc1918c918c918c91613fad565b8051909150610fd1576001610fd4565b60005b93506000610ffb8360400151846080015185602001518660e001518c87602001518d61412b565b9050826040015167ffffffffffffffff168b7f526019d853b0b88e828dbdc124746afde5a47053256d2a0872a025ab145b0520836020015189898760000151611044578e611046565b8f5b6040516110569493929190615b48565b60405180910390a351925050505b965096945050505050565b600c5460009060ff16156110af576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6110b761439a565b6110fe8260008989898080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508b92508a91506144079050565b90505b9695505050505050565b600c5460ff1615611148576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6bffffffffffffffffffffffff811661118d576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b336000908152601060205260409020546bffffffffffffffffffffffff808316911610156111e7576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b33600090815260106020526040812080548392906112149084906bffffffffffffffffffffffff16615d2f565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555080600d60008282829054906101000a90046bffffffffffffffffffffffff1661126b9190615d2f565b82546bffffffffffffffffffffffff91821661010093840a9081029083021990911617909255600c546040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff87811660048301529386166024820152919004909116915063a9059cbb90604401602060405180830381600087803b15801561130a57600080fd5b505af115801561131e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611342919061541d565b61079c576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600c5460ff16156113b5576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b8181101561185c5760008383838181106113d4576113d4615e8f565b602090810292909201356000818152601184526040808220815161010081018352815473ffffffffffffffffffffffffffffffffffffffff908116825260018301549081169782019790975274010000000000000000000000000000000000000000870467ffffffffffffffff168184018190527c010000000000000000000000000000000000000000000000000000000090970463ffffffff16606082015260028201546bffffffffffffffffffffffff9081166080830152600383015460a083015260048084015460c08401526005909301541660e082015291517fa47c7696000000000000000000000000000000000000000000000000000000008152908101959095529194509092909150309063a47c76969060240160006040518083038186803b15801561150657600080fd5b505afa15801561151a573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526115609190810190615764565b5050925050508073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146115e3576040517fd8a3fb5200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8216600482015260240161086f565b8160a00151421015611621576040517fbcc4005500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b81516040517f85b214cf0000000000000000000000000000000000000000000000000000000081526004810185905273ffffffffffffffffffffffffffffffffffffffff8216906385b214cf90602401602060405180830381600087803b15801561168b57600080fd5b505af115801561169f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906116c3919061541d565b1561184557608083015160408085015167ffffffffffffffff166000908152600e602052208054600c906117169084906c0100000000000000000000000090046bffffffffffffffffffffffff16615d2f565b82546bffffffffffffffffffffffff9182166101009390930a92830291909202199091161790555060208084015173ffffffffffffffffffffffffffffffffffffffff166000908152600f825260408082208187015167ffffffffffffffff908116845293529020805460019260099161179f9185916901000000000000000000900416615c9a565b825467ffffffffffffffff9182166101009390930a928302919092021990911617905550600084815260116020526040812080547fffffffffffffffffffffffff0000000000000000000000000000000000000000168155600181018290556002810180547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000908116909155600382018390556004820192909255600501805490911690555b50505050808061185490615db0565b9150506113b8565b505050565b600061186d8383614816565b9392505050565b67ffffffffffffffff82166000908152600e6020526040902060010154829073ffffffffffffffffffffffffffffffffffffffff16806118e0576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff821614611947576040517fd8a3fb5200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8216600482015260240161086f565b600c5460ff1615611984576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff84166000908152600e6020526040902060030154606414156119db576040517f05a48e0f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff83166000908152600f6020908152604080832067ffffffffffffffff8816845290915290205460ff16151560011415611a2857610988565b73ffffffffffffffffffffffffffffffffffffffff83166000818152600f6020908152604080832067ffffffffffffffff891680855290835281842080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001908117909155600e84528285206003018054918201815585529383902090930180547fffffffffffffffffffffffff000000000000000000000000000000000000000016851790555192835290917f43dc749a04ac8fb825cbd514f7c0e13f13bc6f2ee66043b76629d51776cff8e0910161097f565b60015473ffffffffffffffffffffffffffffffffffffffff163314611b81576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015260640161086f565b60008054336101008181027fffffffffffffffffffffff0000000000000000000000000000000000000000ff8416178455600180547fffffffffffffffffffffffff000000000000000000000000000000000000000016905560405173ffffffffffffffffffffffffffffffffffffffff919093041692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b600c5460ff1615611c5f576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff81166000908152600e602052604090206001015473ffffffffffffffffffffffffffffffffffffffff16611cc8576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff81166000908152600e602052604090206002015473ffffffffffffffffffffffffffffffffffffffff163314611d6a5767ffffffffffffffff81166000908152600e6020526040908190206002015490517fd084e97500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff909116600482015260240161086f565b67ffffffffffffffff81166000818152600e60209081526040918290206001810180547fffffffffffffffffffffffff00000000000000000000000000000000000000008082163390811790935560029093018054909316909255835173ffffffffffffffffffffffffffffffffffffffff909216808352928201529092917f6f1dc65165ffffedfd8e507b4a0f1fcfdada045ed11f6c26ba27cedfe87802f0910160405180910390a25050565b611e206136bf565b600a54431015611e5c576040517fa93d035c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600954600880547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000ffff16620100009283900461ffff16909202919091179055565b67ffffffffffffffff82166000908152600e6020526040902060010154829073ffffffffffffffffffffffffffffffffffffffff1680611f08576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff821614611f6f576040517fd8a3fb5200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8216600482015260240161086f565b600c5460ff1615611fac576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff83166000908152600f6020908152604080832067ffffffffffffffff8089168552908352928190208151606081018352905460ff811615158083526101008204861694830194909452690100000000000000000090049093169083015261207c576040517ff0019fe600000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8616600482015273ffffffffffffffffffffffffffffffffffffffff8516602482015260440161086f565b806040015167ffffffffffffffff16816020015167ffffffffffffffff16146120d1576040517fbcc4005500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff85166000908152600e602090815260408083206003018054825181850281018501909352808352919290919083018282801561214c57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311612121575b505050505090506000600182516121639190615d18565b905060005b8251811015612302578673ffffffffffffffffffffffffffffffffffffffff1683828151811061219a5761219a615e8f565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1614156122f05760008383815181106121d2576121d2615e8f565b6020026020010151905080600e60008b67ffffffffffffffff1667ffffffffffffffff168152602001908152602001600020600301838154811061221857612218615e8f565b600091825260208083209190910180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff949094169390931790925567ffffffffffffffff8b168152600e9091526040902060030180548061229257612292615e60565b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff000000000000000000000000000000000000000016905501905550612302565b806122fa81615db0565b915050612168565b5073ffffffffffffffffffffffffffffffffffffffff86166000818152600f6020908152604080832067ffffffffffffffff8c168085529083529281902080547fffffffffffffffffffffffffffffff00000000000000000000000000000000001690555192835290917f182bff9831466789164ca77075fffd84916d35a8180ba73c27e45634549b445b910160405180910390a250505050505050565b600c5460009060ff16156123e0576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600c80547501000000000000000000000000000000000000000000900467ffffffffffffffff1690601561241383615de9565b82546101009290920a67ffffffffffffffff818102199093169183160217909155600c547501000000000000000000000000000000000000000000900416905060008060405190808252806020026020018201604052801561247f578160200160208202803683370190505b506040805160a0810182526000808252602080830182815233848601908152606085018481526080860188815267ffffffffffffffff8b168652600e8552969094208551815493516bffffffffffffffffffffffff9081166c01000000000000000000000000027fffffffffffffffff0000000000000000000000000000000000000000000000009095169116179290921782555160018201805473ffffffffffffffffffffffffffffffffffffffff9283167fffffffffffffffffffffffff000000000000000000000000000000000000000091821617909155935160028301805491909216941693909317909255925180519495509193909261258b926003850192910190614f2c565b505060405133815267ffffffffffffffff841691507f464722b4166576d3dcbba877b999bc35cf911f4eaf434b7eba68fa113951d0bf9060200160405180910390a250905090565b60008060008060606125e486614946565b67ffffffffffffffff86166000908152600e6020908152604091829020805460018201546002830154600390930180548651818702810187019097528087526bffffffffffffffffffffffff8085169c506c01000000000000000000000000909404909316995073ffffffffffffffffffffffffffffffffffffffff91821698509216955090918301828280156126b157602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311612686575b5050505050905091939590929450565b600c5460ff16156126fe576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600c54610100900473ffffffffffffffffffffffffffffffffffffffff163314612754576040517f44b0e3c300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6020811461278e576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600061279c828401846155dc565b67ffffffffffffffff81166000908152600e602052604090206001015490915073ffffffffffffffffffffffffffffffffffffffff16612808576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff81166000908152600e6020526040812080546bffffffffffffffffffffffff169186919061283f8385615cbd565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555084600d60008282829054906101000a90046bffffffffffffffffffffffff166128969190615cbd565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055508167ffffffffffffffff167fd39ec07f4e209f627a4c427971473820dc129761ba28de8906bd56f57101d4f88287846128fd9190615c82565b6040805192835260208301919091520160405180910390a2505050505050565b600061292a826000614816565b92915050565b6129386136bf565b60085461ffff82811662010000909204161415612981576040517fee03280800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60085461ffff90811690821611156129c5576040517fe9a3062200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805160608101825260085461ffff620100009091048116808352908416602083015290918201906129f89043615c82565b9052805160098054602084015161ffff90811662010000027fffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000090921693169290921791909117905560400151600a5550565b612a526136bf565b600654431015612a8e576040517fa93d035c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600160149054906101000a900461ffff166001612aab9190615c5c565b600180547fffffffffffffffffffff0000ffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000061ffff93841602179081905576010000000000000000000000000000000000000000000090041615612b3c57600180547fffffffffffffffff0000ffffffffffffffffffffffffffffffffffffffffffff1690555b60005b60035460ff82161015612d3c576005805460ff8316908110612b6357612b63615e8f565b60009182526020822001546003805473ffffffffffffffffffffffffffffffffffffffff9092169260029290919060ff8616908110612ba457612ba4615e8f565b9060005260206000200154815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055507ff17ce6b5bbeda6598ac50e505071ddbd27c188c5ecaf9126b42c951c9c4b61e960036000018260ff1681548110612c3557612c35615e8f565b906000526020600020015460036001018360ff1681548110612c5957612c59615e8f565b6000918252602090912001546005805473ffffffffffffffffffffffffffffffffffffffff9092169160ff8616908110612c9557612c95615e8f565b60009182526020918290200154600180546040805196875273ffffffffffffffffffffffffffffffffffffffff958616948701949094529390911691840191909152606083015261ffff740100000000000000000000000000000000000000008204811660808401527601000000000000000000000000000000000000000000009091041660a082015260c00160405180910390a180612d3481615e11565b915050612b3f565b50565b60006060806060600380015460036000016003600101600360020182805480602002602001604051908101604052809291908181526020018280548015612da557602002820191906000526020600020905b815481526020019060010190808311612d91575b5050505050925081805480602002602001604051908101604052809291908181526020018280548015612e0e57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311612de3575b5050505050915080805480602002602001604051908101604052809291908181526020018280548015612e7757602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311612e4c575b50505050509050935093509350935090919293565b6060612e998484846149af565b949350505050565b67ffffffffffffffff82166000908152600e6020526040902060010154829073ffffffffffffffffffffffffffffffffffffffff1680612f0d576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff821614612f74576040517fd8a3fb5200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8216600482015260240161086f565b600c5460ff1615612fb1576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040517fe82ad7d400000000000000000000000000000000000000000000000000000000815267ffffffffffffffff85166004820152309063e82ad7d49060240160206040518083038186803b15801561300a57600080fd5b505afa15801561301e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613042919061541d565b15613079576040517fb42f66e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6109888484613745565b61308b6136bf565b600c546040517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152600091610100900473ffffffffffffffffffffffffffffffffffffffff16906370a082319060240160206040518083038186803b1580156130fa57600080fd5b505afa15801561310e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906131329190615453565b600d549091506bffffffffffffffffffffffff168181111561318a576040517fa99da302000000000000000000000000000000000000000000000000000000008152600481018290526024810183905260440161086f565b8181101561185c57600061319e8284615d18565b600c546040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8781166004830152602482018490529293506101009091049091169063a9059cbb90604401602060405180830381600087803b15801561321b57600080fd5b505af115801561322f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613253919061541d565b506040805173ffffffffffffffffffffffffffffffffffffffff86168152602081018390527f59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b436600910160405180910390a150505050565b67ffffffffffffffff81166000908152600e602090815260408083206003018054825181850281018501909352808352849383018282801561332157602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116132f6575b5050505050905060005b81518110156133f0576000600f600084848151811061334c5761334c615e8f565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040908101600090812067ffffffffffffffff808a168352908452908290208251606081018452905460ff81161515825261010081048316948201859052690100000000000000000090049091169181018290529250146133dd57506001949350505050565b50806133e881615db0565b91505061332b565b5060009392505050565b6134026136bf565b600081815260026020908152604080832054600783528184208251606081019093528054835260018101805473ffffffffffffffffffffffffffffffffffffffff909316959491929184019161345790615d5c565b80601f016020809104026020016040519081016040528092919081815260200182805461348390615d5c565b80156134d05780601f106134a5576101008083540402835291602001916134d0565b820191906000526020600020905b8154815290600101906020018083116134b357829003601f168201915b5050505050815260200160028201548152505090508060400151431015613523576040517fa93d035c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b82613548576135358160200151614a7e565b6020808201518051910120600b556135ff565b60208101516040517f354497aa00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff84169163354497aa9161359e9190600401615a35565b600060405180830381600087803b1580156135b857600080fd5b505af19250505080156135c9575060015b6135ff576040517ffe680b2600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600160169054906101000a900461ffff16600161361c9190615c5c565b6001805461ffff92909216760100000000000000000000000000000000000000000000027fffffffffffffffff0000ffffffffffffffffffffffffffffffffffffffffffff909216919091179055805160208201516040517fc07626096b49b0462a576c9a7878cf0675e784bb4832584c1289c11664e55f47926136a19287926159e2565b60405180910390a1505050565b6136b66136bf565b612d3c81614b54565b600054610100900473ffffffffffffffffffffffffffffffffffffffff163314610b38576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161086f565b600c5460ff1615613782576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff82166000908152600e60209081526040808320815160a08101835281546bffffffffffffffffffffffff80821683526c010000000000000000000000009091041681850152600182015473ffffffffffffffffffffffffffffffffffffffff90811682850152600283015416606082015260038201805484518187028101870190955280855291949293608086019390929083018282801561386357602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311613838575b5050509190925250505067ffffffffffffffff84166000908152600e60205260408120549192506bffffffffffffffffffffffff909116905b82608001515181101561394257600f6000846080015183815181106138c3576138c3615e8f565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040908101600090812067ffffffffffffffff89168252909252902080547fffffffffffffffffffffffffffffff00000000000000000000000000000000001690558061393a81615db0565b91505061389c565b5067ffffffffffffffff84166000908152600e6020526040812080547fffffffffffffffff0000000000000000000000000000000000000000000000001681556001810180547fffffffffffffffffffffffff00000000000000000000000000000000000000009081169091556002820180549091169055906139c86003830182615019565b5050600d80548291906000906139ed9084906bffffffffffffffffffffffff16615d2f565b82546bffffffffffffffffffffffff91821661010093840a9081029083021990911617909255600c546040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8881166004830152938616602482015291900491909116915063a9059cbb90604401602060405180830381600087803b158015613a8d57600080fd5b505af1158015613aa1573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613ac5919061541d565b613afb576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805173ffffffffffffffffffffffffffffffffffffffff851681526bffffffffffffffffffffffff8316602082015267ffffffffffffffff8616917fe8ed5b475a5b5987aa9165e8731bb78043f39eee32ec5a1169a89e27fcd49815910161097f565b81518351141580613b7357508051835114155b15613baa576040517fee03280800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b825160081015613be6576040517fee03280800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b83518160ff16101561098857600073ffffffffffffffffffffffffffffffffffffffff16828260ff1681518110613c2257613c22615e8f565b602002602001015173ffffffffffffffffffffffffffffffffffffffff161415613c78576040517fee03280800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b818160ff1681518110613c8d57613c8d615e8f565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1660026000868460ff1681518110613cc457613cc4615e8f565b60209081029190910181015182528101919091526040016000205473ffffffffffffffffffffffffffffffffffffffff161415613d2d576040517fee03280800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b828160ff1681518110613d4257613d42615e8f565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1660026000868460ff1681518110613d7957613d79615e8f565b60209081029190910181015182528101919091526040016000205473ffffffffffffffffffffffffffffffffffffffff1614613de1576040517fee03280800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000801b848260ff1681518110613dfa57613dfa615e8f565b60200260200101511415613e3d576040517f7ee392220000000000000000000000000000000000000000000000000000000081526000600482015260240161086f565b80613e4781615e11565b915050613be9565b613e57614c50565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa335b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390a1565b613ed461439a565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258613ea23390565b6000613f398663ffffffff8616615c82565b5a1015613f4857506003611101565b866bffffffffffffffffffffffff16613f668563ffffffff16614cbc565b613f709085615ce4565b613f7a8488615cbd565b613f849190615cbd565b6bffffffffffffffffffffffff161115613fa057506004611101565b5060009695505050505050565b60408051808201909152600080825260208201526012546040516000916c01000000000000000000000000900460e01b90613ff090899089908990602401615a0a565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090931692909217909152600c80547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016600117905590506000805a6113888110156140ad57600080fd5b6113888103905086604082048203116140c557600080fd5b853b6140d057600080fd5b60008085516020870160008a8cf192505a600c80547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690556040805180820190915293151584529003602083015250979650505050505050565b6040805180820190915260008082526020820152600061414a84614cbc565b6141549086615ce4565b905060006141628482615cbd565b905061416e8782615cbd565b905061417a8282615cbd565b6040805180820182526bffffffffffffffffffffffff808616825280841660208084019190915267ffffffffffffffff8f166000908152600e9091529283208054929750939450849392916141d191859116615d2f565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550818461420b9190615cbd565b33600090815260106020526040812080549091906142389084906bffffffffffffffffffffffff16615cbd565b82546101009290920a6bffffffffffffffffffffffff81810219909316918316021790915530600090815260106020526040812080548b9450909261427f91859116615cbd565b82546101009290920a6bffffffffffffffffffffffff81810219909316918316021790915567ffffffffffffffff8c166000908152600e6020526040902080548c93509091600c916142e39185916c01000000000000000000000000900416615d2f565b82546bffffffffffffffffffffffff9182166101009390930a92830291909202199091161790555073ffffffffffffffffffffffffffffffffffffffff88166000908152600f6020908152604080832067ffffffffffffffff808f168552925290912080546001926009916143679185916901000000000000000000900416615c9a565b92506101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055505050979650505050505050565b60005460ff1615610b38576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f5061757361626c653a2070617573656400000000000000000000000000000000604482015260640161086f565b600061441285614946565b61441c3386614d5e565b60006144288888614816565b6040805160c08101825267ffffffffffffffff891680825260208083018a905261ffff89168385015263ffffffff881660608401523360808401526000918252600e90528281206001015473ffffffffffffffffffffffffffffffffffffffff90811660a084015292517f0a33e0a5000000000000000000000000000000000000000000000000000000008152939450928392839290861691630a33e0a5916144d391600401615a48565b608060405180830381600087803b1580156144ed57600080fd5b505af1158015614501573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906145259190615583565b9297509094509250905061453a338a85614dfa565b6040518061010001604052808573ffffffffffffffffffffffffffffffffffffffff1681526020013373ffffffffffffffffffffffffffffffffffffffff1681526020018a67ffffffffffffffff1681526020018763ffffffff168152602001846bffffffffffffffffffffffff16815260200182426145ba9190615c82565b8152602001838152602001601260000160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff168152506011600087815260200190815260200160002060008201518160000160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060208201518160010160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060408201518160010160146101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550606082015181600101601c6101000a81548163ffffffff021916908363ffffffff16021790555060808201518160020160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060a0820151816003015560c0820151816004015560e08201518160050160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055509050508867ffffffffffffffff168b867f7c720ccd20069b8311a6be4ba1cf3294d09eb247aa5d73a8502054b6e68a2f54600e60008e67ffffffffffffffff1667ffffffffffffffff16815260200190815260200160002060010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1633328e8e8e60405161480096959493929190615923565b60405180910390a4505050509695505050505050565b6000816148835760008381526002602052604090205473ffffffffffffffffffffffffffffffffffffffff168061487c576040517f80833e330000000000000000000000000000000000000000000000000000000081526004810185905260240161086f565b905061292a565b60005b60035460ff82161015614910576003805460ff83169081106148aa576148aa615e8f565b90600052602060002001548414156148fe576005805460ff83169081106148d3576148d3615e8f565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff16915061292a9050565b8061490881615e11565b915050614886565b506040517f80833e330000000000000000000000000000000000000000000000000000000081526004810184905260240161086f565b67ffffffffffffffff81166000908152600e602052604090206001015473ffffffffffffffffffffffffffffffffffffffff16612d3c576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b606060008080806149c2868801886156a1565b935093509350935060006149db89600187878787614407565b60408051602080825281830190925291925060208201818036833701905050955060005b6020811015614a7157818160208110614a1a57614a1a615e8f565b1a60f81b878281518110614a3057614a30615e8f565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535080614a6981615db0565b9150506149ff565b5050505050509392505050565b60008082806020019051810190614a959190615710565b6040805180820182526bffffffffffffffffffffffff84168082527fffffffff0000000000000000000000000000000000000000000000000000000084166020928301819052601280547fffffffffffffffffffffffffffffffff000000000000000000000000000000001683176c0100000000000000000000000060e088901c021790558351918252918101919091529294509092507fb2b9bfb736b0bfef266985352cfca2a114226362209f5a6e5a2a2ed14383ba5791016136a1565b73ffffffffffffffffffffffffffffffffffffffff8116331415614bd4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161086f565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929361010090910416917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60005460ff16610b38576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f5061757361626c653a206e6f7420706175736564000000000000000000000000604482015260640161086f565b60006bffffffffffffffffffffffff821115614d5a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f53616665436173743a2076616c756520646f65736e27742066697420696e203960448201527f3620626974730000000000000000000000000000000000000000000000000000606482015260840161086f565b5090565b73ffffffffffffffffffffffffffffffffffffffff82166000908152600f6020908152604080832067ffffffffffffffff8516845290915290205460ff1661079c576040517ff0019fe600000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8216600482015273ffffffffffffffffffffffffffffffffffffffff8316602482015260440161086f565b67ffffffffffffffff82166000908152600e602052604090208054829190600c90614e449084906c0100000000000000000000000090046bffffffffffffffffffffffff16615cbd565b82546bffffffffffffffffffffffff91821661010093840a908102920219161790915573ffffffffffffffffffffffffffffffffffffffff85166000908152600f6020908152604080832067ffffffffffffffff8089168552925290912080546001945090928492614eba928492900416615c9a565b92506101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550505050565b828054828255906000526020600020908101928215614f20579160200282015b82811115614f20578251825591602001919060010190614f05565b50614d5a929150615033565b828054828255906000526020600020908101928215614f20579160200282015b82811115614f2057825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190614f4c565b828054614fb290615d5c565b90600052602060002090601f016020900481019282614fd45760008555614f20565b82601f10614fed57805160ff1916838001178555614f20565b82800160010185558215614f205791820182811115614f20578251825591602001919060010190614f05565b5080546000825590600052602060002090810190612d3c91905b5b80821115614d5a5760008155600101615034565b600082601f83011261505957600080fd5b8135602061506e61506983615c38565b615be9565b80838252828201915082860187848660051b890101111561508e57600080fd5b60005b858110156150b65781356150a481615eed565b84529284019290840190600101615091565b5090979650505050505050565b60008083601f8401126150d557600080fd5b50813567ffffffffffffffff8111156150ed57600080fd5b60208301915083602082850101111561510557600080fd5b9250929050565b600082601f83011261511d57600080fd5b813567ffffffffffffffff81111561513757615137615ebe565b61516860207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601615be9565b81815284602083860101111561517d57600080fd5b816020850160208301376000918101602001919091529392505050565b803561ffff811681146151ac57600080fd5b919050565b803563ffffffff811681146151ac57600080fd5b803567ffffffffffffffff811681146151ac57600080fd5b6000602082840312156151ef57600080fd5b813561186d81615eed565b6000806000806060858703121561521057600080fd5b843561521b81615eed565b935060208501359250604085013567ffffffffffffffff81111561523e57600080fd5b61524a878288016150c3565b95989497509550505050565b6000806040838503121561526957600080fd5b823561527481615eed565b9150615282602084016151c5565b90509250929050565b6000806040838503121561529e57600080fd5b82356152a981615eed565b915060208301356152b981615f1d565b809150509250929050565b600080602083850312156152d757600080fd5b823567ffffffffffffffff808211156152ef57600080fd5b818501915085601f83011261530357600080fd5b81358181111561531257600080fd5b8660208260051b850101111561532757600080fd5b60209290920196919550909350505050565b60008060006060848603121561534e57600080fd5b833567ffffffffffffffff8082111561536657600080fd5b818601915086601f83011261537a57600080fd5b8135602061538a61506983615c38565b8083825282820191508286018b848660051b89010111156153aa57600080fd5b600096505b848710156153cd5780358352600196909601959183019183016153af565b50975050870135925050808211156153e457600080fd5b6153f087838801615048565b9350604086013591508082111561540657600080fd5b5061541386828701615048565b9150509250925092565b60006020828403121561542f57600080fd5b815161186d81615f0f565b60006020828403121561544c57600080fd5b5035919050565b60006020828403121561546557600080fd5b5051919050565b6000806040838503121561547f57600080fd5b8235915060208301356152b981615f0f565b6000806000604084860312156154a657600080fd5b83359250602084013567ffffffffffffffff8111156154c457600080fd5b6154d0868287016150c3565b9497909650939450505050565b60008060008060008060c087890312156154f657600080fd5b86359550602087013567ffffffffffffffff8082111561551557600080fd5b6155218a838b0161510c565b9650604089013591508082111561553757600080fd5b5061554489828a0161510c565b945050606087013561555581615f1d565b9250608087013561556581615f1d565b915060a087013561557581615eed565b809150509295509295509295565b6000806000806080858703121561559957600080fd5b8451935060208501516155ab81615f1d565b6040860151606090960151949790965092505050565b6000602082840312156155d357600080fd5b61186d8261519a565b6000602082840312156155ee57600080fd5b61186d826151c5565b6000806040838503121561560a57600080fd5b615613836151c5565b915060208301356152b981615eed565b60008060008060008060a0878903121561563c57600080fd5b615645876151c5565b9550602087013567ffffffffffffffff81111561566157600080fd5b61566d89828a016150c3565b909650945061568090506040880161519a565b925061568e606088016151b1565b9150608087013590509295509295509295565b600080600080608085870312156156b757600080fd5b6156c0856151c5565b9350602085013567ffffffffffffffff8111156156dc57600080fd5b6156e88782880161510c565b9350506156f76040860161519a565b9150615705606086016151b1565b905092959194509250565b6000806040838503121561572357600080fd5b825161572e81615f1d565b60208401519092507fffffffff00000000000000000000000000000000000000000000000000000000811681146152b957600080fd5b600080600080600060a0868803121561577c57600080fd5b855161578781615f1d565b8095505060208087015161579a81615f1d565b60408801519095506157ab81615eed565b60608801519094506157bc81615eed565b608088015190935067ffffffffffffffff8111156157d957600080fd5b8701601f810189136157ea57600080fd5b80516157f861506982615c38565b8082825284820191508484018c868560051b870101111561581857600080fd5b600094505b8385101561584457805161583081615eed565b83526001949094019391850191850161581d565b5080955050505050509295509295909350565b600081518084526020808501945080840160005b8381101561589d57815173ffffffffffffffffffffffffffffffffffffffff168752958201959082019060010161586b565b509495945050505050565b6000815180845260005b818110156158ce576020818501810151868301820152016158b2565b818111156158e0576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b8183823760009101908152919050565b600073ffffffffffffffffffffffffffffffffffffffff8089168352808816602084015280871660408401525060c0606083015261596460c08301866158a8565b905061ffff8416608083015263ffffffff831660a0830152979650505050505050565b84815283602082015260606040820152816060820152818360808301376000818301608090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01601019392505050565b838152826020820152606060408201526000615a0160608301846158a8565b95945050505050565b838152606060208201526000615a2360608301856158a8565b828103604084015261110181856158a8565b60208152600061186d60208301846158a8565b6020815267ffffffffffffffff82511660208201526000602083015160c06040840152615a7860e08401826158a8565b905061ffff604085015116606084015263ffffffff6060850151166080840152608084015173ffffffffffffffffffffffffffffffffffffffff80821660a08601528060a08701511660c086015250508091505092915050565b600060808201868352602060808185015281875180845260a086019150828901935060005b81811015615b1357845183529383019391830191600101615af7565b50508481036040860152615b278188615857565b925050508281036060840152615b3d8185615857565b979650505050505050565b6bffffffffffffffffffffffff8516815273ffffffffffffffffffffffffffffffffffffffff8416602082015260ff8316604082015260806060820152600061110160808301846158a8565b6bffffffffffffffffffffffff86811682528516602082015273ffffffffffffffffffffffffffffffffffffffff84811660408301528316606082015260a0608082018190526000906110fe90830184615857565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715615c3057615c30615ebe565b604052919050565b600067ffffffffffffffff821115615c5257615c52615ebe565b5060051b60200190565b600061ffff808316818516808303821115615c7957615c79615e31565b01949350505050565b60008219821115615c9557615c95615e31565b500190565b600067ffffffffffffffff808316818516808303821115615c7957615c79615e31565b60006bffffffffffffffffffffffff808316818516808303821115615c7957615c79615e31565b60006bffffffffffffffffffffffff80831681851681830481118215151615615d0f57615d0f615e31565b02949350505050565b600082821015615d2a57615d2a615e31565b500390565b60006bffffffffffffffffffffffff83811690831681811015615d5457615d54615e31565b039392505050565b600181811c90821680615d7057607f821691505b60208210811415615daa577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415615de257615de2615e31565b5060010190565b600067ffffffffffffffff80831681811415615e0757615e07615e31565b6001019392505050565b600060ff821660ff811415615e2857615e28615e31565b60010192915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b73ffffffffffffffffffffffffffffffffffffffff81168114612d3c57600080fd5b8015158114612d3c57600080fd5b6bffffffffffffffffffffffff81168114612d3c57600080fdfea164736f6c6343000806000a",
}

var FunctionsRouterABI = FunctionsRouterMetaData.ABI

var FunctionsRouterBin = FunctionsRouterMetaData.Bin

func DeployFunctionsRouter(auth *bind.TransactOpts, backend bind.ContractBackend, timelockBlocks uint16, maximumTimelockBlocks uint16, linkToken common.Address, config []byte) (common.Address, *types.Transaction, *FunctionsRouter, error) {
	parsed, err := FunctionsRouterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(FunctionsRouterBin), backend, timelockBlocks, maximumTimelockBlocks, linkToken, config)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &FunctionsRouter{FunctionsRouterCaller: FunctionsRouterCaller{contract: contract}, FunctionsRouterTransactor: FunctionsRouterTransactor{contract: contract}, FunctionsRouterFilterer: FunctionsRouterFilterer{contract: contract}}, nil
}

type FunctionsRouter struct {
	address common.Address
	abi     abi.ABI
	FunctionsRouterCaller
	FunctionsRouterTransactor
	FunctionsRouterFilterer
}

type FunctionsRouterCaller struct {
	contract *bind.BoundContract
}

type FunctionsRouterTransactor struct {
	contract *bind.BoundContract
}

type FunctionsRouterFilterer struct {
	contract *bind.BoundContract
}

type FunctionsRouterSession struct {
	Contract     *FunctionsRouter
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type FunctionsRouterCallerSession struct {
	Contract *FunctionsRouterCaller
	CallOpts bind.CallOpts
}

type FunctionsRouterTransactorSession struct {
	Contract     *FunctionsRouterTransactor
	TransactOpts bind.TransactOpts
}

type FunctionsRouterRaw struct {
	Contract *FunctionsRouter
}

type FunctionsRouterCallerRaw struct {
	Contract *FunctionsRouterCaller
}

type FunctionsRouterTransactorRaw struct {
	Contract *FunctionsRouterTransactor
}

func NewFunctionsRouter(address common.Address, backend bind.ContractBackend) (*FunctionsRouter, error) {
	abi, err := abi.JSON(strings.NewReader(FunctionsRouterABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindFunctionsRouter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FunctionsRouter{address: address, abi: abi, FunctionsRouterCaller: FunctionsRouterCaller{contract: contract}, FunctionsRouterTransactor: FunctionsRouterTransactor{contract: contract}, FunctionsRouterFilterer: FunctionsRouterFilterer{contract: contract}}, nil
}

func NewFunctionsRouterCaller(address common.Address, caller bind.ContractCaller) (*FunctionsRouterCaller, error) {
	contract, err := bindFunctionsRouter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FunctionsRouterCaller{contract: contract}, nil
}

func NewFunctionsRouterTransactor(address common.Address, transactor bind.ContractTransactor) (*FunctionsRouterTransactor, error) {
	contract, err := bindFunctionsRouter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FunctionsRouterTransactor{contract: contract}, nil
}

func NewFunctionsRouterFilterer(address common.Address, filterer bind.ContractFilterer) (*FunctionsRouterFilterer, error) {
	contract, err := bindFunctionsRouter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FunctionsRouterFilterer{contract: contract}, nil
}

func bindFunctionsRouter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FunctionsRouterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_FunctionsRouter *FunctionsRouterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FunctionsRouter.Contract.FunctionsRouterCaller.contract.Call(opts, result, method, params...)
}

func (_FunctionsRouter *FunctionsRouterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.FunctionsRouterTransactor.contract.Transfer(opts)
}

func (_FunctionsRouter *FunctionsRouterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.FunctionsRouterTransactor.contract.Transact(opts, method, params...)
}

func (_FunctionsRouter *FunctionsRouterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FunctionsRouter.Contract.contract.Call(opts, result, method, params...)
}

func (_FunctionsRouter *FunctionsRouterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.contract.Transfer(opts)
}

func (_FunctionsRouter *FunctionsRouterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.contract.Transact(opts, method, params...)
}

func (_FunctionsRouter *FunctionsRouterCaller) MAXCONSUMERS(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _FunctionsRouter.contract.Call(opts, &out, "MAX_CONSUMERS")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_FunctionsRouter *FunctionsRouterSession) MAXCONSUMERS() (uint16, error) {
	return _FunctionsRouter.Contract.MAXCONSUMERS(&_FunctionsRouter.CallOpts)
}

func (_FunctionsRouter *FunctionsRouterCallerSession) MAXCONSUMERS() (uint16, error) {
	return _FunctionsRouter.Contract.MAXCONSUMERS(&_FunctionsRouter.CallOpts)
}

func (_FunctionsRouter *FunctionsRouterCaller) GetAdminFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FunctionsRouter.contract.Call(opts, &out, "getAdminFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_FunctionsRouter *FunctionsRouterSession) GetAdminFee() (*big.Int, error) {
	return _FunctionsRouter.Contract.GetAdminFee(&_FunctionsRouter.CallOpts)
}

func (_FunctionsRouter *FunctionsRouterCallerSession) GetAdminFee() (*big.Int, error) {
	return _FunctionsRouter.Contract.GetAdminFee(&_FunctionsRouter.CallOpts)
}

func (_FunctionsRouter *FunctionsRouterCaller) GetConfigHash(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _FunctionsRouter.contract.Call(opts, &out, "getConfigHash")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_FunctionsRouter *FunctionsRouterSession) GetConfigHash() ([32]byte, error) {
	return _FunctionsRouter.Contract.GetConfigHash(&_FunctionsRouter.CallOpts)
}

func (_FunctionsRouter *FunctionsRouterCallerSession) GetConfigHash() ([32]byte, error) {
	return _FunctionsRouter.Contract.GetConfigHash(&_FunctionsRouter.CallOpts)
}

func (_FunctionsRouter *FunctionsRouterCaller) GetConsumer(opts *bind.CallOpts, client common.Address, subscriptionId uint64) (GetConsumer,

	error) {
	var out []interface{}
	err := _FunctionsRouter.contract.Call(opts, &out, "getConsumer", client, subscriptionId)

	outstruct := new(GetConsumer)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Allowed = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.InitiatedRequests = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	outstruct.CompletedRequests = *abi.ConvertType(out[2], new(uint64)).(*uint64)

	return *outstruct, err

}

func (_FunctionsRouter *FunctionsRouterSession) GetConsumer(client common.Address, subscriptionId uint64) (GetConsumer,

	error) {
	return _FunctionsRouter.Contract.GetConsumer(&_FunctionsRouter.CallOpts, client, subscriptionId)
}

func (_FunctionsRouter *FunctionsRouterCallerSession) GetConsumer(client common.Address, subscriptionId uint64) (GetConsumer,

	error) {
	return _FunctionsRouter.Contract.GetConsumer(&_FunctionsRouter.CallOpts, client, subscriptionId)
}

func (_FunctionsRouter *FunctionsRouterCaller) GetContractById(opts *bind.CallOpts, id [32]byte, useProposed bool) (common.Address, error) {
	var out []interface{}
	err := _FunctionsRouter.contract.Call(opts, &out, "getContractById", id, useProposed)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_FunctionsRouter *FunctionsRouterSession) GetContractById(id [32]byte, useProposed bool) (common.Address, error) {
	return _FunctionsRouter.Contract.GetContractById(&_FunctionsRouter.CallOpts, id, useProposed)
}

func (_FunctionsRouter *FunctionsRouterCallerSession) GetContractById(id [32]byte, useProposed bool) (common.Address, error) {
	return _FunctionsRouter.Contract.GetContractById(&_FunctionsRouter.CallOpts, id, useProposed)
}

func (_FunctionsRouter *FunctionsRouterCaller) GetContractById0(opts *bind.CallOpts, id [32]byte) (common.Address, error) {
	var out []interface{}
	err := _FunctionsRouter.contract.Call(opts, &out, "getContractById0", id)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_FunctionsRouter *FunctionsRouterSession) GetContractById0(id [32]byte) (common.Address, error) {
	return _FunctionsRouter.Contract.GetContractById0(&_FunctionsRouter.CallOpts, id)
}

func (_FunctionsRouter *FunctionsRouterCallerSession) GetContractById0(id [32]byte) (common.Address, error) {
	return _FunctionsRouter.Contract.GetContractById0(&_FunctionsRouter.CallOpts, id)
}

func (_FunctionsRouter *FunctionsRouterCaller) GetProposedContractSet(opts *bind.CallOpts) (*big.Int, [][32]byte, []common.Address, []common.Address, error) {
	var out []interface{}
	err := _FunctionsRouter.contract.Call(opts, &out, "getProposedContractSet")

	if err != nil {
		return *new(*big.Int), *new([][32]byte), *new([]common.Address), *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new([][32]byte)).(*[][32]byte)
	out2 := *abi.ConvertType(out[2], new([]common.Address)).(*[]common.Address)
	out3 := *abi.ConvertType(out[3], new([]common.Address)).(*[]common.Address)

	return out0, out1, out2, out3, err

}

func (_FunctionsRouter *FunctionsRouterSession) GetProposedContractSet() (*big.Int, [][32]byte, []common.Address, []common.Address, error) {
	return _FunctionsRouter.Contract.GetProposedContractSet(&_FunctionsRouter.CallOpts)
}

func (_FunctionsRouter *FunctionsRouterCallerSession) GetProposedContractSet() (*big.Int, [][32]byte, []common.Address, []common.Address, error) {
	return _FunctionsRouter.Contract.GetProposedContractSet(&_FunctionsRouter.CallOpts)
}

func (_FunctionsRouter *FunctionsRouterCaller) GetSubscription(opts *bind.CallOpts, subscriptionId uint64) (GetSubscription,

	error) {
	var out []interface{}
	err := _FunctionsRouter.contract.Call(opts, &out, "getSubscription", subscriptionId)

	outstruct := new(GetSubscription)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Balance = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.BlockedBalance = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Owner = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.RequestedOwner = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)
	outstruct.Consumers = *abi.ConvertType(out[4], new([]common.Address)).(*[]common.Address)

	return *outstruct, err

}

func (_FunctionsRouter *FunctionsRouterSession) GetSubscription(subscriptionId uint64) (GetSubscription,

	error) {
	return _FunctionsRouter.Contract.GetSubscription(&_FunctionsRouter.CallOpts, subscriptionId)
}

func (_FunctionsRouter *FunctionsRouterCallerSession) GetSubscription(subscriptionId uint64) (GetSubscription,

	error) {
	return _FunctionsRouter.Contract.GetSubscription(&_FunctionsRouter.CallOpts, subscriptionId)
}

func (_FunctionsRouter *FunctionsRouterCaller) GetSubscriptionCount(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _FunctionsRouter.contract.Call(opts, &out, "getSubscriptionCount")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_FunctionsRouter *FunctionsRouterSession) GetSubscriptionCount() (uint64, error) {
	return _FunctionsRouter.Contract.GetSubscriptionCount(&_FunctionsRouter.CallOpts)
}

func (_FunctionsRouter *FunctionsRouterCallerSession) GetSubscriptionCount() (uint64, error) {
	return _FunctionsRouter.Contract.GetSubscriptionCount(&_FunctionsRouter.CallOpts)
}

func (_FunctionsRouter *FunctionsRouterCaller) GetTotalBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FunctionsRouter.contract.Call(opts, &out, "getTotalBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_FunctionsRouter *FunctionsRouterSession) GetTotalBalance() (*big.Int, error) {
	return _FunctionsRouter.Contract.GetTotalBalance(&_FunctionsRouter.CallOpts)
}

func (_FunctionsRouter *FunctionsRouterCallerSession) GetTotalBalance() (*big.Int, error) {
	return _FunctionsRouter.Contract.GetTotalBalance(&_FunctionsRouter.CallOpts)
}

func (_FunctionsRouter *FunctionsRouterCaller) IsPaused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _FunctionsRouter.contract.Call(opts, &out, "isPaused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_FunctionsRouter *FunctionsRouterSession) IsPaused() (bool, error) {
	return _FunctionsRouter.Contract.IsPaused(&_FunctionsRouter.CallOpts)
}

func (_FunctionsRouter *FunctionsRouterCallerSession) IsPaused() (bool, error) {
	return _FunctionsRouter.Contract.IsPaused(&_FunctionsRouter.CallOpts)
}

func (_FunctionsRouter *FunctionsRouterCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FunctionsRouter.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_FunctionsRouter *FunctionsRouterSession) Owner() (common.Address, error) {
	return _FunctionsRouter.Contract.Owner(&_FunctionsRouter.CallOpts)
}

func (_FunctionsRouter *FunctionsRouterCallerSession) Owner() (common.Address, error) {
	return _FunctionsRouter.Contract.Owner(&_FunctionsRouter.CallOpts)
}

func (_FunctionsRouter *FunctionsRouterCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _FunctionsRouter.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_FunctionsRouter *FunctionsRouterSession) Paused() (bool, error) {
	return _FunctionsRouter.Contract.Paused(&_FunctionsRouter.CallOpts)
}

func (_FunctionsRouter *FunctionsRouterCallerSession) Paused() (bool, error) {
	return _FunctionsRouter.Contract.Paused(&_FunctionsRouter.CallOpts)
}

func (_FunctionsRouter *FunctionsRouterCaller) PendingRequestExists(opts *bind.CallOpts, subscriptionId uint64) (bool, error) {
	var out []interface{}
	err := _FunctionsRouter.contract.Call(opts, &out, "pendingRequestExists", subscriptionId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_FunctionsRouter *FunctionsRouterSession) PendingRequestExists(subscriptionId uint64) (bool, error) {
	return _FunctionsRouter.Contract.PendingRequestExists(&_FunctionsRouter.CallOpts, subscriptionId)
}

func (_FunctionsRouter *FunctionsRouterCallerSession) PendingRequestExists(subscriptionId uint64) (bool, error) {
	return _FunctionsRouter.Contract.PendingRequestExists(&_FunctionsRouter.CallOpts, subscriptionId)
}

func (_FunctionsRouter *FunctionsRouterCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _FunctionsRouter.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_FunctionsRouter *FunctionsRouterSession) TypeAndVersion() (string, error) {
	return _FunctionsRouter.Contract.TypeAndVersion(&_FunctionsRouter.CallOpts)
}

func (_FunctionsRouter *FunctionsRouterCallerSession) TypeAndVersion() (string, error) {
	return _FunctionsRouter.Contract.TypeAndVersion(&_FunctionsRouter.CallOpts)
}

func (_FunctionsRouter *FunctionsRouterCaller) Version(opts *bind.CallOpts) (uint16, uint16, uint16, error) {
	var out []interface{}
	err := _FunctionsRouter.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(uint16), *new(uint16), *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)
	out1 := *abi.ConvertType(out[1], new(uint16)).(*uint16)
	out2 := *abi.ConvertType(out[2], new(uint16)).(*uint16)

	return out0, out1, out2, err

}

func (_FunctionsRouter *FunctionsRouterSession) Version() (uint16, uint16, uint16, error) {
	return _FunctionsRouter.Contract.Version(&_FunctionsRouter.CallOpts)
}

func (_FunctionsRouter *FunctionsRouterCallerSession) Version() (uint16, uint16, uint16, error) {
	return _FunctionsRouter.Contract.Version(&_FunctionsRouter.CallOpts)
}

func (_FunctionsRouter *FunctionsRouterTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsRouter.contract.Transact(opts, "acceptOwnership")
}

func (_FunctionsRouter *FunctionsRouterSession) AcceptOwnership() (*types.Transaction, error) {
	return _FunctionsRouter.Contract.AcceptOwnership(&_FunctionsRouter.TransactOpts)
}

func (_FunctionsRouter *FunctionsRouterTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _FunctionsRouter.Contract.AcceptOwnership(&_FunctionsRouter.TransactOpts)
}

func (_FunctionsRouter *FunctionsRouterTransactor) AcceptSubscriptionOwnerTransfer(opts *bind.TransactOpts, subscriptionId uint64) (*types.Transaction, error) {
	return _FunctionsRouter.contract.Transact(opts, "acceptSubscriptionOwnerTransfer", subscriptionId)
}

func (_FunctionsRouter *FunctionsRouterSession) AcceptSubscriptionOwnerTransfer(subscriptionId uint64) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.AcceptSubscriptionOwnerTransfer(&_FunctionsRouter.TransactOpts, subscriptionId)
}

func (_FunctionsRouter *FunctionsRouterTransactorSession) AcceptSubscriptionOwnerTransfer(subscriptionId uint64) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.AcceptSubscriptionOwnerTransfer(&_FunctionsRouter.TransactOpts, subscriptionId)
}

func (_FunctionsRouter *FunctionsRouterTransactor) AddConsumer(opts *bind.TransactOpts, subscriptionId uint64, consumer common.Address) (*types.Transaction, error) {
	return _FunctionsRouter.contract.Transact(opts, "addConsumer", subscriptionId, consumer)
}

func (_FunctionsRouter *FunctionsRouterSession) AddConsumer(subscriptionId uint64, consumer common.Address) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.AddConsumer(&_FunctionsRouter.TransactOpts, subscriptionId, consumer)
}

func (_FunctionsRouter *FunctionsRouterTransactorSession) AddConsumer(subscriptionId uint64, consumer common.Address) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.AddConsumer(&_FunctionsRouter.TransactOpts, subscriptionId, consumer)
}

func (_FunctionsRouter *FunctionsRouterTransactor) CancelSubscription(opts *bind.TransactOpts, subscriptionId uint64, to common.Address) (*types.Transaction, error) {
	return _FunctionsRouter.contract.Transact(opts, "cancelSubscription", subscriptionId, to)
}

func (_FunctionsRouter *FunctionsRouterSession) CancelSubscription(subscriptionId uint64, to common.Address) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.CancelSubscription(&_FunctionsRouter.TransactOpts, subscriptionId, to)
}

func (_FunctionsRouter *FunctionsRouterTransactorSession) CancelSubscription(subscriptionId uint64, to common.Address) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.CancelSubscription(&_FunctionsRouter.TransactOpts, subscriptionId, to)
}

func (_FunctionsRouter *FunctionsRouterTransactor) CreateSubscription(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsRouter.contract.Transact(opts, "createSubscription")
}

func (_FunctionsRouter *FunctionsRouterSession) CreateSubscription() (*types.Transaction, error) {
	return _FunctionsRouter.Contract.CreateSubscription(&_FunctionsRouter.TransactOpts)
}

func (_FunctionsRouter *FunctionsRouterTransactorSession) CreateSubscription() (*types.Transaction, error) {
	return _FunctionsRouter.Contract.CreateSubscription(&_FunctionsRouter.TransactOpts)
}

func (_FunctionsRouter *FunctionsRouterTransactor) Fulfill(opts *bind.TransactOpts, requestId [32]byte, response []byte, err []byte, juelsPerGas *big.Int, costWithoutFulfillment *big.Int, transmitter common.Address) (*types.Transaction, error) {
	return _FunctionsRouter.contract.Transact(opts, "fulfill", requestId, response, err, juelsPerGas, costWithoutFulfillment, transmitter)
}

func (_FunctionsRouter *FunctionsRouterSession) Fulfill(requestId [32]byte, response []byte, err []byte, juelsPerGas *big.Int, costWithoutFulfillment *big.Int, transmitter common.Address) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.Fulfill(&_FunctionsRouter.TransactOpts, requestId, response, err, juelsPerGas, costWithoutFulfillment, transmitter)
}

func (_FunctionsRouter *FunctionsRouterTransactorSession) Fulfill(requestId [32]byte, response []byte, err []byte, juelsPerGas *big.Int, costWithoutFulfillment *big.Int, transmitter common.Address) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.Fulfill(&_FunctionsRouter.TransactOpts, requestId, response, err, juelsPerGas, costWithoutFulfillment, transmitter)
}

func (_FunctionsRouter *FunctionsRouterTransactor) OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _FunctionsRouter.contract.Transact(opts, "onTokenTransfer", arg0, amount, data)
}

func (_FunctionsRouter *FunctionsRouterSession) OnTokenTransfer(arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.OnTokenTransfer(&_FunctionsRouter.TransactOpts, arg0, amount, data)
}

func (_FunctionsRouter *FunctionsRouterTransactorSession) OnTokenTransfer(arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.OnTokenTransfer(&_FunctionsRouter.TransactOpts, arg0, amount, data)
}

func (_FunctionsRouter *FunctionsRouterTransactor) OracleWithdraw(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FunctionsRouter.contract.Transact(opts, "oracleWithdraw", recipient, amount)
}

func (_FunctionsRouter *FunctionsRouterSession) OracleWithdraw(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.OracleWithdraw(&_FunctionsRouter.TransactOpts, recipient, amount)
}

func (_FunctionsRouter *FunctionsRouterTransactorSession) OracleWithdraw(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.OracleWithdraw(&_FunctionsRouter.TransactOpts, recipient, amount)
}

func (_FunctionsRouter *FunctionsRouterTransactor) OwnerCancelSubscription(opts *bind.TransactOpts, subscriptionId uint64) (*types.Transaction, error) {
	return _FunctionsRouter.contract.Transact(opts, "ownerCancelSubscription", subscriptionId)
}

func (_FunctionsRouter *FunctionsRouterSession) OwnerCancelSubscription(subscriptionId uint64) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.OwnerCancelSubscription(&_FunctionsRouter.TransactOpts, subscriptionId)
}

func (_FunctionsRouter *FunctionsRouterTransactorSession) OwnerCancelSubscription(subscriptionId uint64) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.OwnerCancelSubscription(&_FunctionsRouter.TransactOpts, subscriptionId)
}

func (_FunctionsRouter *FunctionsRouterTransactor) ProposeConfigUpdate(opts *bind.TransactOpts, id [32]byte, config []byte) (*types.Transaction, error) {
	return _FunctionsRouter.contract.Transact(opts, "proposeConfigUpdate", id, config)
}

func (_FunctionsRouter *FunctionsRouterSession) ProposeConfigUpdate(id [32]byte, config []byte) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.ProposeConfigUpdate(&_FunctionsRouter.TransactOpts, id, config)
}

func (_FunctionsRouter *FunctionsRouterTransactorSession) ProposeConfigUpdate(id [32]byte, config []byte) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.ProposeConfigUpdate(&_FunctionsRouter.TransactOpts, id, config)
}

func (_FunctionsRouter *FunctionsRouterTransactor) ProposeContractsUpdate(opts *bind.TransactOpts, proposedContractSetIds [][32]byte, proposedContractSetFromAddresses []common.Address, proposedContractSetToAddresses []common.Address) (*types.Transaction, error) {
	return _FunctionsRouter.contract.Transact(opts, "proposeContractsUpdate", proposedContractSetIds, proposedContractSetFromAddresses, proposedContractSetToAddresses)
}

func (_FunctionsRouter *FunctionsRouterSession) ProposeContractsUpdate(proposedContractSetIds [][32]byte, proposedContractSetFromAddresses []common.Address, proposedContractSetToAddresses []common.Address) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.ProposeContractsUpdate(&_FunctionsRouter.TransactOpts, proposedContractSetIds, proposedContractSetFromAddresses, proposedContractSetToAddresses)
}

func (_FunctionsRouter *FunctionsRouterTransactorSession) ProposeContractsUpdate(proposedContractSetIds [][32]byte, proposedContractSetFromAddresses []common.Address, proposedContractSetToAddresses []common.Address) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.ProposeContractsUpdate(&_FunctionsRouter.TransactOpts, proposedContractSetIds, proposedContractSetFromAddresses, proposedContractSetToAddresses)
}

func (_FunctionsRouter *FunctionsRouterTransactor) ProposeTimelockBlocks(opts *bind.TransactOpts, blocks uint16) (*types.Transaction, error) {
	return _FunctionsRouter.contract.Transact(opts, "proposeTimelockBlocks", blocks)
}

func (_FunctionsRouter *FunctionsRouterSession) ProposeTimelockBlocks(blocks uint16) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.ProposeTimelockBlocks(&_FunctionsRouter.TransactOpts, blocks)
}

func (_FunctionsRouter *FunctionsRouterTransactorSession) ProposeTimelockBlocks(blocks uint16) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.ProposeTimelockBlocks(&_FunctionsRouter.TransactOpts, blocks)
}

func (_FunctionsRouter *FunctionsRouterTransactor) RecoverFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _FunctionsRouter.contract.Transact(opts, "recoverFunds", to)
}

func (_FunctionsRouter *FunctionsRouterSession) RecoverFunds(to common.Address) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.RecoverFunds(&_FunctionsRouter.TransactOpts, to)
}

func (_FunctionsRouter *FunctionsRouterTransactorSession) RecoverFunds(to common.Address) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.RecoverFunds(&_FunctionsRouter.TransactOpts, to)
}

func (_FunctionsRouter *FunctionsRouterTransactor) RemoveConsumer(opts *bind.TransactOpts, subscriptionId uint64, consumer common.Address) (*types.Transaction, error) {
	return _FunctionsRouter.contract.Transact(opts, "removeConsumer", subscriptionId, consumer)
}

func (_FunctionsRouter *FunctionsRouterSession) RemoveConsumer(subscriptionId uint64, consumer common.Address) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.RemoveConsumer(&_FunctionsRouter.TransactOpts, subscriptionId, consumer)
}

func (_FunctionsRouter *FunctionsRouterTransactorSession) RemoveConsumer(subscriptionId uint64, consumer common.Address) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.RemoveConsumer(&_FunctionsRouter.TransactOpts, subscriptionId, consumer)
}

func (_FunctionsRouter *FunctionsRouterTransactor) RequestSubscriptionOwnerTransfer(opts *bind.TransactOpts, subscriptionId uint64, newOwner common.Address) (*types.Transaction, error) {
	return _FunctionsRouter.contract.Transact(opts, "requestSubscriptionOwnerTransfer", subscriptionId, newOwner)
}

func (_FunctionsRouter *FunctionsRouterSession) RequestSubscriptionOwnerTransfer(subscriptionId uint64, newOwner common.Address) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.RequestSubscriptionOwnerTransfer(&_FunctionsRouter.TransactOpts, subscriptionId, newOwner)
}

func (_FunctionsRouter *FunctionsRouterTransactorSession) RequestSubscriptionOwnerTransfer(subscriptionId uint64, newOwner common.Address) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.RequestSubscriptionOwnerTransfer(&_FunctionsRouter.TransactOpts, subscriptionId, newOwner)
}

func (_FunctionsRouter *FunctionsRouterTransactor) SendRequest(opts *bind.TransactOpts, subscriptionId uint64, data []byte, dataVersion uint16, callbackGasLimit uint32, donId [32]byte) (*types.Transaction, error) {
	return _FunctionsRouter.contract.Transact(opts, "sendRequest", subscriptionId, data, dataVersion, callbackGasLimit, donId)
}

func (_FunctionsRouter *FunctionsRouterSession) SendRequest(subscriptionId uint64, data []byte, dataVersion uint16, callbackGasLimit uint32, donId [32]byte) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.SendRequest(&_FunctionsRouter.TransactOpts, subscriptionId, data, dataVersion, callbackGasLimit, donId)
}

func (_FunctionsRouter *FunctionsRouterTransactorSession) SendRequest(subscriptionId uint64, data []byte, dataVersion uint16, callbackGasLimit uint32, donId [32]byte) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.SendRequest(&_FunctionsRouter.TransactOpts, subscriptionId, data, dataVersion, callbackGasLimit, donId)
}

func (_FunctionsRouter *FunctionsRouterTransactor) TimeoutRequests(opts *bind.TransactOpts, requestIdsToTimeout [][32]byte) (*types.Transaction, error) {
	return _FunctionsRouter.contract.Transact(opts, "timeoutRequests", requestIdsToTimeout)
}

func (_FunctionsRouter *FunctionsRouterSession) TimeoutRequests(requestIdsToTimeout [][32]byte) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.TimeoutRequests(&_FunctionsRouter.TransactOpts, requestIdsToTimeout)
}

func (_FunctionsRouter *FunctionsRouterTransactorSession) TimeoutRequests(requestIdsToTimeout [][32]byte) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.TimeoutRequests(&_FunctionsRouter.TransactOpts, requestIdsToTimeout)
}

func (_FunctionsRouter *FunctionsRouterTransactor) TogglePaused(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsRouter.contract.Transact(opts, "togglePaused")
}

func (_FunctionsRouter *FunctionsRouterSession) TogglePaused() (*types.Transaction, error) {
	return _FunctionsRouter.Contract.TogglePaused(&_FunctionsRouter.TransactOpts)
}

func (_FunctionsRouter *FunctionsRouterTransactorSession) TogglePaused() (*types.Transaction, error) {
	return _FunctionsRouter.Contract.TogglePaused(&_FunctionsRouter.TransactOpts)
}

func (_FunctionsRouter *FunctionsRouterTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _FunctionsRouter.contract.Transact(opts, "transferOwnership", to)
}

func (_FunctionsRouter *FunctionsRouterSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.TransferOwnership(&_FunctionsRouter.TransactOpts, to)
}

func (_FunctionsRouter *FunctionsRouterTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.TransferOwnership(&_FunctionsRouter.TransactOpts, to)
}

func (_FunctionsRouter *FunctionsRouterTransactor) UpdateConfig(opts *bind.TransactOpts, id [32]byte) (*types.Transaction, error) {
	return _FunctionsRouter.contract.Transact(opts, "updateConfig", id)
}

func (_FunctionsRouter *FunctionsRouterSession) UpdateConfig(id [32]byte) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.UpdateConfig(&_FunctionsRouter.TransactOpts, id)
}

func (_FunctionsRouter *FunctionsRouterTransactorSession) UpdateConfig(id [32]byte) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.UpdateConfig(&_FunctionsRouter.TransactOpts, id)
}

func (_FunctionsRouter *FunctionsRouterTransactor) UpdateContracts(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsRouter.contract.Transact(opts, "updateContracts")
}

func (_FunctionsRouter *FunctionsRouterSession) UpdateContracts() (*types.Transaction, error) {
	return _FunctionsRouter.Contract.UpdateContracts(&_FunctionsRouter.TransactOpts)
}

func (_FunctionsRouter *FunctionsRouterTransactorSession) UpdateContracts() (*types.Transaction, error) {
	return _FunctionsRouter.Contract.UpdateContracts(&_FunctionsRouter.TransactOpts)
}

func (_FunctionsRouter *FunctionsRouterTransactor) UpdateTimelockBlocks(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsRouter.contract.Transact(opts, "updateTimelockBlocks")
}

func (_FunctionsRouter *FunctionsRouterSession) UpdateTimelockBlocks() (*types.Transaction, error) {
	return _FunctionsRouter.Contract.UpdateTimelockBlocks(&_FunctionsRouter.TransactOpts)
}

func (_FunctionsRouter *FunctionsRouterTransactorSession) UpdateTimelockBlocks() (*types.Transaction, error) {
	return _FunctionsRouter.Contract.UpdateTimelockBlocks(&_FunctionsRouter.TransactOpts)
}

func (_FunctionsRouter *FunctionsRouterTransactor) ValidateProposedContracts(opts *bind.TransactOpts, id [32]byte, data []byte) (*types.Transaction, error) {
	return _FunctionsRouter.contract.Transact(opts, "validateProposedContracts", id, data)
}

func (_FunctionsRouter *FunctionsRouterSession) ValidateProposedContracts(id [32]byte, data []byte) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.ValidateProposedContracts(&_FunctionsRouter.TransactOpts, id, data)
}

func (_FunctionsRouter *FunctionsRouterTransactorSession) ValidateProposedContracts(id [32]byte, data []byte) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.ValidateProposedContracts(&_FunctionsRouter.TransactOpts, id, data)
}

type FunctionsRouterConfigProposedIterator struct {
	Event *FunctionsRouterConfigProposed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRouterConfigProposedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRouterConfigProposed)
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
		it.Event = new(FunctionsRouterConfigProposed)
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

func (it *FunctionsRouterConfigProposedIterator) Error() error {
	return it.fail
}

func (it *FunctionsRouterConfigProposedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRouterConfigProposed struct {
	Id       [32]byte
	FromHash [32]byte
	ToBytes  []byte
	Raw      types.Log
}

func (_FunctionsRouter *FunctionsRouterFilterer) FilterConfigProposed(opts *bind.FilterOpts) (*FunctionsRouterConfigProposedIterator, error) {

	logs, sub, err := _FunctionsRouter.contract.FilterLogs(opts, "ConfigProposed")
	if err != nil {
		return nil, err
	}
	return &FunctionsRouterConfigProposedIterator{contract: _FunctionsRouter.contract, event: "ConfigProposed", logs: logs, sub: sub}, nil
}

func (_FunctionsRouter *FunctionsRouterFilterer) WatchConfigProposed(opts *bind.WatchOpts, sink chan<- *FunctionsRouterConfigProposed) (event.Subscription, error) {

	logs, sub, err := _FunctionsRouter.contract.WatchLogs(opts, "ConfigProposed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRouterConfigProposed)
				if err := _FunctionsRouter.contract.UnpackLog(event, "ConfigProposed", log); err != nil {
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

func (_FunctionsRouter *FunctionsRouterFilterer) ParseConfigProposed(log types.Log) (*FunctionsRouterConfigProposed, error) {
	event := new(FunctionsRouterConfigProposed)
	if err := _FunctionsRouter.contract.UnpackLog(event, "ConfigProposed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsRouterConfigSetIterator struct {
	Event *FunctionsRouterConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRouterConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRouterConfigSet)
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
		it.Event = new(FunctionsRouterConfigSet)
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

func (it *FunctionsRouterConfigSetIterator) Error() error {
	return it.fail
}

func (it *FunctionsRouterConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRouterConfigSet struct {
	AdminFee                        *big.Int
	HandleOracleFulfillmentSelector [4]byte
	Raw                             types.Log
}

func (_FunctionsRouter *FunctionsRouterFilterer) FilterConfigSet(opts *bind.FilterOpts) (*FunctionsRouterConfigSetIterator, error) {

	logs, sub, err := _FunctionsRouter.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &FunctionsRouterConfigSetIterator{contract: _FunctionsRouter.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_FunctionsRouter *FunctionsRouterFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *FunctionsRouterConfigSet) (event.Subscription, error) {

	logs, sub, err := _FunctionsRouter.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRouterConfigSet)
				if err := _FunctionsRouter.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_FunctionsRouter *FunctionsRouterFilterer) ParseConfigSet(log types.Log) (*FunctionsRouterConfigSet, error) {
	event := new(FunctionsRouterConfigSet)
	if err := _FunctionsRouter.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsRouterConfigUpdatedIterator struct {
	Event *FunctionsRouterConfigUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRouterConfigUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRouterConfigUpdated)
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
		it.Event = new(FunctionsRouterConfigUpdated)
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

func (it *FunctionsRouterConfigUpdatedIterator) Error() error {
	return it.fail
}

func (it *FunctionsRouterConfigUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRouterConfigUpdated struct {
	Id       [32]byte
	FromHash [32]byte
	ToBytes  []byte
	Raw      types.Log
}

func (_FunctionsRouter *FunctionsRouterFilterer) FilterConfigUpdated(opts *bind.FilterOpts) (*FunctionsRouterConfigUpdatedIterator, error) {

	logs, sub, err := _FunctionsRouter.contract.FilterLogs(opts, "ConfigUpdated")
	if err != nil {
		return nil, err
	}
	return &FunctionsRouterConfigUpdatedIterator{contract: _FunctionsRouter.contract, event: "ConfigUpdated", logs: logs, sub: sub}, nil
}

func (_FunctionsRouter *FunctionsRouterFilterer) WatchConfigUpdated(opts *bind.WatchOpts, sink chan<- *FunctionsRouterConfigUpdated) (event.Subscription, error) {

	logs, sub, err := _FunctionsRouter.contract.WatchLogs(opts, "ConfigUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRouterConfigUpdated)
				if err := _FunctionsRouter.contract.UnpackLog(event, "ConfigUpdated", log); err != nil {
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

func (_FunctionsRouter *FunctionsRouterFilterer) ParseConfigUpdated(log types.Log) (*FunctionsRouterConfigUpdated, error) {
	event := new(FunctionsRouterConfigUpdated)
	if err := _FunctionsRouter.contract.UnpackLog(event, "ConfigUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsRouterContractProposedIterator struct {
	Event *FunctionsRouterContractProposed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRouterContractProposedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRouterContractProposed)
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
		it.Event = new(FunctionsRouterContractProposed)
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

func (it *FunctionsRouterContractProposedIterator) Error() error {
	return it.fail
}

func (it *FunctionsRouterContractProposedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRouterContractProposed struct {
	ProposedContractSetId          [32]byte
	ProposedContractSetFromAddress common.Address
	ProposedContractSetToAddress   common.Address
	TimelockEndBlock               *big.Int
	Raw                            types.Log
}

func (_FunctionsRouter *FunctionsRouterFilterer) FilterContractProposed(opts *bind.FilterOpts) (*FunctionsRouterContractProposedIterator, error) {

	logs, sub, err := _FunctionsRouter.contract.FilterLogs(opts, "ContractProposed")
	if err != nil {
		return nil, err
	}
	return &FunctionsRouterContractProposedIterator{contract: _FunctionsRouter.contract, event: "ContractProposed", logs: logs, sub: sub}, nil
}

func (_FunctionsRouter *FunctionsRouterFilterer) WatchContractProposed(opts *bind.WatchOpts, sink chan<- *FunctionsRouterContractProposed) (event.Subscription, error) {

	logs, sub, err := _FunctionsRouter.contract.WatchLogs(opts, "ContractProposed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRouterContractProposed)
				if err := _FunctionsRouter.contract.UnpackLog(event, "ContractProposed", log); err != nil {
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

func (_FunctionsRouter *FunctionsRouterFilterer) ParseContractProposed(log types.Log) (*FunctionsRouterContractProposed, error) {
	event := new(FunctionsRouterContractProposed)
	if err := _FunctionsRouter.contract.UnpackLog(event, "ContractProposed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsRouterContractUpdatedIterator struct {
	Event *FunctionsRouterContractUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRouterContractUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRouterContractUpdated)
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
		it.Event = new(FunctionsRouterContractUpdated)
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

func (it *FunctionsRouterContractUpdatedIterator) Error() error {
	return it.fail
}

func (it *FunctionsRouterContractUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRouterContractUpdated struct {
	ProposedContractSetId          [32]byte
	ProposedContractSetFromAddress common.Address
	ProposedContractSetToAddress   common.Address
	Major                          uint16
	Minor                          uint16
	Patch                          uint16
	Raw                            types.Log
}

func (_FunctionsRouter *FunctionsRouterFilterer) FilterContractUpdated(opts *bind.FilterOpts) (*FunctionsRouterContractUpdatedIterator, error) {

	logs, sub, err := _FunctionsRouter.contract.FilterLogs(opts, "ContractUpdated")
	if err != nil {
		return nil, err
	}
	return &FunctionsRouterContractUpdatedIterator{contract: _FunctionsRouter.contract, event: "ContractUpdated", logs: logs, sub: sub}, nil
}

func (_FunctionsRouter *FunctionsRouterFilterer) WatchContractUpdated(opts *bind.WatchOpts, sink chan<- *FunctionsRouterContractUpdated) (event.Subscription, error) {

	logs, sub, err := _FunctionsRouter.contract.WatchLogs(opts, "ContractUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRouterContractUpdated)
				if err := _FunctionsRouter.contract.UnpackLog(event, "ContractUpdated", log); err != nil {
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

func (_FunctionsRouter *FunctionsRouterFilterer) ParseContractUpdated(log types.Log) (*FunctionsRouterContractUpdated, error) {
	event := new(FunctionsRouterContractUpdated)
	if err := _FunctionsRouter.contract.UnpackLog(event, "ContractUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsRouterFundsRecoveredIterator struct {
	Event *FunctionsRouterFundsRecovered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRouterFundsRecoveredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRouterFundsRecovered)
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
		it.Event = new(FunctionsRouterFundsRecovered)
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

func (it *FunctionsRouterFundsRecoveredIterator) Error() error {
	return it.fail
}

func (it *FunctionsRouterFundsRecoveredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRouterFundsRecovered struct {
	To     common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_FunctionsRouter *FunctionsRouterFilterer) FilterFundsRecovered(opts *bind.FilterOpts) (*FunctionsRouterFundsRecoveredIterator, error) {

	logs, sub, err := _FunctionsRouter.contract.FilterLogs(opts, "FundsRecovered")
	if err != nil {
		return nil, err
	}
	return &FunctionsRouterFundsRecoveredIterator{contract: _FunctionsRouter.contract, event: "FundsRecovered", logs: logs, sub: sub}, nil
}

func (_FunctionsRouter *FunctionsRouterFilterer) WatchFundsRecovered(opts *bind.WatchOpts, sink chan<- *FunctionsRouterFundsRecovered) (event.Subscription, error) {

	logs, sub, err := _FunctionsRouter.contract.WatchLogs(opts, "FundsRecovered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRouterFundsRecovered)
				if err := _FunctionsRouter.contract.UnpackLog(event, "FundsRecovered", log); err != nil {
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

func (_FunctionsRouter *FunctionsRouterFilterer) ParseFundsRecovered(log types.Log) (*FunctionsRouterFundsRecovered, error) {
	event := new(FunctionsRouterFundsRecovered)
	if err := _FunctionsRouter.contract.UnpackLog(event, "FundsRecovered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsRouterOwnershipTransferRequestedIterator struct {
	Event *FunctionsRouterOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRouterOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRouterOwnershipTransferRequested)
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
		it.Event = new(FunctionsRouterOwnershipTransferRequested)
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

func (it *FunctionsRouterOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *FunctionsRouterOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRouterOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_FunctionsRouter *FunctionsRouterFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsRouterOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsRouter.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsRouterOwnershipTransferRequestedIterator{contract: _FunctionsRouter.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_FunctionsRouter *FunctionsRouterFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *FunctionsRouterOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsRouter.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRouterOwnershipTransferRequested)
				if err := _FunctionsRouter.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_FunctionsRouter *FunctionsRouterFilterer) ParseOwnershipTransferRequested(log types.Log) (*FunctionsRouterOwnershipTransferRequested, error) {
	event := new(FunctionsRouterOwnershipTransferRequested)
	if err := _FunctionsRouter.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsRouterOwnershipTransferredIterator struct {
	Event *FunctionsRouterOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRouterOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRouterOwnershipTransferred)
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
		it.Event = new(FunctionsRouterOwnershipTransferred)
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

func (it *FunctionsRouterOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *FunctionsRouterOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRouterOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_FunctionsRouter *FunctionsRouterFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsRouterOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsRouter.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsRouterOwnershipTransferredIterator{contract: _FunctionsRouter.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_FunctionsRouter *FunctionsRouterFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FunctionsRouterOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsRouter.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRouterOwnershipTransferred)
				if err := _FunctionsRouter.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_FunctionsRouter *FunctionsRouterFilterer) ParseOwnershipTransferred(log types.Log) (*FunctionsRouterOwnershipTransferred, error) {
	event := new(FunctionsRouterOwnershipTransferred)
	if err := _FunctionsRouter.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsRouterPausedIterator struct {
	Event *FunctionsRouterPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRouterPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRouterPaused)
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
		it.Event = new(FunctionsRouterPaused)
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

func (it *FunctionsRouterPausedIterator) Error() error {
	return it.fail
}

func (it *FunctionsRouterPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRouterPaused struct {
	Account common.Address
	Raw     types.Log
}

func (_FunctionsRouter *FunctionsRouterFilterer) FilterPaused(opts *bind.FilterOpts) (*FunctionsRouterPausedIterator, error) {

	logs, sub, err := _FunctionsRouter.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &FunctionsRouterPausedIterator{contract: _FunctionsRouter.contract, event: "Paused", logs: logs, sub: sub}, nil
}

func (_FunctionsRouter *FunctionsRouterFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *FunctionsRouterPaused) (event.Subscription, error) {

	logs, sub, err := _FunctionsRouter.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRouterPaused)
				if err := _FunctionsRouter.contract.UnpackLog(event, "Paused", log); err != nil {
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

func (_FunctionsRouter *FunctionsRouterFilterer) ParsePaused(log types.Log) (*FunctionsRouterPaused, error) {
	event := new(FunctionsRouterPaused)
	if err := _FunctionsRouter.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsRouterRequestEndIterator struct {
	Event *FunctionsRouterRequestEnd

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRouterRequestEndIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRouterRequestEnd)
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
		it.Event = new(FunctionsRouterRequestEnd)
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

func (it *FunctionsRouterRequestEndIterator) Error() error {
	return it.fail
}

func (it *FunctionsRouterRequestEndIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRouterRequestEnd struct {
	RequestId      [32]byte
	SubscriptionId uint64
	TotalCostJuels *big.Int
	Transmitter    common.Address
	ResultCode     uint8
	Response       []byte
	Raw            types.Log
}

func (_FunctionsRouter *FunctionsRouterFilterer) FilterRequestEnd(opts *bind.FilterOpts, requestId [][32]byte, subscriptionId []uint64) (*FunctionsRouterRequestEndIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsRouter.contract.FilterLogs(opts, "RequestEnd", requestIdRule, subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsRouterRequestEndIterator{contract: _FunctionsRouter.contract, event: "RequestEnd", logs: logs, sub: sub}, nil
}

func (_FunctionsRouter *FunctionsRouterFilterer) WatchRequestEnd(opts *bind.WatchOpts, sink chan<- *FunctionsRouterRequestEnd, requestId [][32]byte, subscriptionId []uint64) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsRouter.contract.WatchLogs(opts, "RequestEnd", requestIdRule, subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRouterRequestEnd)
				if err := _FunctionsRouter.contract.UnpackLog(event, "RequestEnd", log); err != nil {
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

func (_FunctionsRouter *FunctionsRouterFilterer) ParseRequestEnd(log types.Log) (*FunctionsRouterRequestEnd, error) {
	event := new(FunctionsRouterRequestEnd)
	if err := _FunctionsRouter.contract.UnpackLog(event, "RequestEnd", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsRouterRequestStartIterator struct {
	Event *FunctionsRouterRequestStart

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRouterRequestStartIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRouterRequestStart)
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
		it.Event = new(FunctionsRouterRequestStart)
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

func (it *FunctionsRouterRequestStartIterator) Error() error {
	return it.fail
}

func (it *FunctionsRouterRequestStartIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRouterRequestStart struct {
	RequestId          [32]byte
	DonId              [32]byte
	SubscriptionId     uint64
	SubscriptionOwner  common.Address
	RequestingContract common.Address
	RequestInitiator   common.Address
	Data               []byte
	DataVersion        uint16
	CallbackGasLimit   uint32
	Raw                types.Log
}

func (_FunctionsRouter *FunctionsRouterFilterer) FilterRequestStart(opts *bind.FilterOpts, requestId [][32]byte, donId [][32]byte, subscriptionId []uint64) (*FunctionsRouterRequestStartIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var donIdRule []interface{}
	for _, donIdItem := range donId {
		donIdRule = append(donIdRule, donIdItem)
	}
	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsRouter.contract.FilterLogs(opts, "RequestStart", requestIdRule, donIdRule, subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsRouterRequestStartIterator{contract: _FunctionsRouter.contract, event: "RequestStart", logs: logs, sub: sub}, nil
}

func (_FunctionsRouter *FunctionsRouterFilterer) WatchRequestStart(opts *bind.WatchOpts, sink chan<- *FunctionsRouterRequestStart, requestId [][32]byte, donId [][32]byte, subscriptionId []uint64) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var donIdRule []interface{}
	for _, donIdItem := range donId {
		donIdRule = append(donIdRule, donIdItem)
	}
	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsRouter.contract.WatchLogs(opts, "RequestStart", requestIdRule, donIdRule, subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRouterRequestStart)
				if err := _FunctionsRouter.contract.UnpackLog(event, "RequestStart", log); err != nil {
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

func (_FunctionsRouter *FunctionsRouterFilterer) ParseRequestStart(log types.Log) (*FunctionsRouterRequestStart, error) {
	event := new(FunctionsRouterRequestStart)
	if err := _FunctionsRouter.contract.UnpackLog(event, "RequestStart", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsRouterSubscriptionCanceledIterator struct {
	Event *FunctionsRouterSubscriptionCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRouterSubscriptionCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRouterSubscriptionCanceled)
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
		it.Event = new(FunctionsRouterSubscriptionCanceled)
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

func (it *FunctionsRouterSubscriptionCanceledIterator) Error() error {
	return it.fail
}

func (it *FunctionsRouterSubscriptionCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRouterSubscriptionCanceled struct {
	SubscriptionId uint64
	FundsRecipient common.Address
	FundsAmount    *big.Int
	Raw            types.Log
}

func (_FunctionsRouter *FunctionsRouterFilterer) FilterSubscriptionCanceled(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsRouterSubscriptionCanceledIterator, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsRouter.contract.FilterLogs(opts, "SubscriptionCanceled", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsRouterSubscriptionCanceledIterator{contract: _FunctionsRouter.contract, event: "SubscriptionCanceled", logs: logs, sub: sub}, nil
}

func (_FunctionsRouter *FunctionsRouterFilterer) WatchSubscriptionCanceled(opts *bind.WatchOpts, sink chan<- *FunctionsRouterSubscriptionCanceled, subscriptionId []uint64) (event.Subscription, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsRouter.contract.WatchLogs(opts, "SubscriptionCanceled", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRouterSubscriptionCanceled)
				if err := _FunctionsRouter.contract.UnpackLog(event, "SubscriptionCanceled", log); err != nil {
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

func (_FunctionsRouter *FunctionsRouterFilterer) ParseSubscriptionCanceled(log types.Log) (*FunctionsRouterSubscriptionCanceled, error) {
	event := new(FunctionsRouterSubscriptionCanceled)
	if err := _FunctionsRouter.contract.UnpackLog(event, "SubscriptionCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsRouterSubscriptionConsumerAddedIterator struct {
	Event *FunctionsRouterSubscriptionConsumerAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRouterSubscriptionConsumerAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRouterSubscriptionConsumerAdded)
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
		it.Event = new(FunctionsRouterSubscriptionConsumerAdded)
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

func (it *FunctionsRouterSubscriptionConsumerAddedIterator) Error() error {
	return it.fail
}

func (it *FunctionsRouterSubscriptionConsumerAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRouterSubscriptionConsumerAdded struct {
	SubscriptionId uint64
	Consumer       common.Address
	Raw            types.Log
}

func (_FunctionsRouter *FunctionsRouterFilterer) FilterSubscriptionConsumerAdded(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsRouterSubscriptionConsumerAddedIterator, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsRouter.contract.FilterLogs(opts, "SubscriptionConsumerAdded", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsRouterSubscriptionConsumerAddedIterator{contract: _FunctionsRouter.contract, event: "SubscriptionConsumerAdded", logs: logs, sub: sub}, nil
}

func (_FunctionsRouter *FunctionsRouterFilterer) WatchSubscriptionConsumerAdded(opts *bind.WatchOpts, sink chan<- *FunctionsRouterSubscriptionConsumerAdded, subscriptionId []uint64) (event.Subscription, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsRouter.contract.WatchLogs(opts, "SubscriptionConsumerAdded", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRouterSubscriptionConsumerAdded)
				if err := _FunctionsRouter.contract.UnpackLog(event, "SubscriptionConsumerAdded", log); err != nil {
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

func (_FunctionsRouter *FunctionsRouterFilterer) ParseSubscriptionConsumerAdded(log types.Log) (*FunctionsRouterSubscriptionConsumerAdded, error) {
	event := new(FunctionsRouterSubscriptionConsumerAdded)
	if err := _FunctionsRouter.contract.UnpackLog(event, "SubscriptionConsumerAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsRouterSubscriptionConsumerRemovedIterator struct {
	Event *FunctionsRouterSubscriptionConsumerRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRouterSubscriptionConsumerRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRouterSubscriptionConsumerRemoved)
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
		it.Event = new(FunctionsRouterSubscriptionConsumerRemoved)
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

func (it *FunctionsRouterSubscriptionConsumerRemovedIterator) Error() error {
	return it.fail
}

func (it *FunctionsRouterSubscriptionConsumerRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRouterSubscriptionConsumerRemoved struct {
	SubscriptionId uint64
	Consumer       common.Address
	Raw            types.Log
}

func (_FunctionsRouter *FunctionsRouterFilterer) FilterSubscriptionConsumerRemoved(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsRouterSubscriptionConsumerRemovedIterator, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsRouter.contract.FilterLogs(opts, "SubscriptionConsumerRemoved", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsRouterSubscriptionConsumerRemovedIterator{contract: _FunctionsRouter.contract, event: "SubscriptionConsumerRemoved", logs: logs, sub: sub}, nil
}

func (_FunctionsRouter *FunctionsRouterFilterer) WatchSubscriptionConsumerRemoved(opts *bind.WatchOpts, sink chan<- *FunctionsRouterSubscriptionConsumerRemoved, subscriptionId []uint64) (event.Subscription, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsRouter.contract.WatchLogs(opts, "SubscriptionConsumerRemoved", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRouterSubscriptionConsumerRemoved)
				if err := _FunctionsRouter.contract.UnpackLog(event, "SubscriptionConsumerRemoved", log); err != nil {
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

func (_FunctionsRouter *FunctionsRouterFilterer) ParseSubscriptionConsumerRemoved(log types.Log) (*FunctionsRouterSubscriptionConsumerRemoved, error) {
	event := new(FunctionsRouterSubscriptionConsumerRemoved)
	if err := _FunctionsRouter.contract.UnpackLog(event, "SubscriptionConsumerRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsRouterSubscriptionCreatedIterator struct {
	Event *FunctionsRouterSubscriptionCreated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRouterSubscriptionCreatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRouterSubscriptionCreated)
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
		it.Event = new(FunctionsRouterSubscriptionCreated)
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

func (it *FunctionsRouterSubscriptionCreatedIterator) Error() error {
	return it.fail
}

func (it *FunctionsRouterSubscriptionCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRouterSubscriptionCreated struct {
	SubscriptionId uint64
	Owner          common.Address
	Raw            types.Log
}

func (_FunctionsRouter *FunctionsRouterFilterer) FilterSubscriptionCreated(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsRouterSubscriptionCreatedIterator, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsRouter.contract.FilterLogs(opts, "SubscriptionCreated", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsRouterSubscriptionCreatedIterator{contract: _FunctionsRouter.contract, event: "SubscriptionCreated", logs: logs, sub: sub}, nil
}

func (_FunctionsRouter *FunctionsRouterFilterer) WatchSubscriptionCreated(opts *bind.WatchOpts, sink chan<- *FunctionsRouterSubscriptionCreated, subscriptionId []uint64) (event.Subscription, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsRouter.contract.WatchLogs(opts, "SubscriptionCreated", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRouterSubscriptionCreated)
				if err := _FunctionsRouter.contract.UnpackLog(event, "SubscriptionCreated", log); err != nil {
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

func (_FunctionsRouter *FunctionsRouterFilterer) ParseSubscriptionCreated(log types.Log) (*FunctionsRouterSubscriptionCreated, error) {
	event := new(FunctionsRouterSubscriptionCreated)
	if err := _FunctionsRouter.contract.UnpackLog(event, "SubscriptionCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsRouterSubscriptionFundedIterator struct {
	Event *FunctionsRouterSubscriptionFunded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRouterSubscriptionFundedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRouterSubscriptionFunded)
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
		it.Event = new(FunctionsRouterSubscriptionFunded)
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

func (it *FunctionsRouterSubscriptionFundedIterator) Error() error {
	return it.fail
}

func (it *FunctionsRouterSubscriptionFundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRouterSubscriptionFunded struct {
	SubscriptionId uint64
	OldBalance     *big.Int
	NewBalance     *big.Int
	Raw            types.Log
}

func (_FunctionsRouter *FunctionsRouterFilterer) FilterSubscriptionFunded(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsRouterSubscriptionFundedIterator, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsRouter.contract.FilterLogs(opts, "SubscriptionFunded", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsRouterSubscriptionFundedIterator{contract: _FunctionsRouter.contract, event: "SubscriptionFunded", logs: logs, sub: sub}, nil
}

func (_FunctionsRouter *FunctionsRouterFilterer) WatchSubscriptionFunded(opts *bind.WatchOpts, sink chan<- *FunctionsRouterSubscriptionFunded, subscriptionId []uint64) (event.Subscription, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsRouter.contract.WatchLogs(opts, "SubscriptionFunded", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRouterSubscriptionFunded)
				if err := _FunctionsRouter.contract.UnpackLog(event, "SubscriptionFunded", log); err != nil {
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

func (_FunctionsRouter *FunctionsRouterFilterer) ParseSubscriptionFunded(log types.Log) (*FunctionsRouterSubscriptionFunded, error) {
	event := new(FunctionsRouterSubscriptionFunded)
	if err := _FunctionsRouter.contract.UnpackLog(event, "SubscriptionFunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsRouterSubscriptionOwnerTransferRequestedIterator struct {
	Event *FunctionsRouterSubscriptionOwnerTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRouterSubscriptionOwnerTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRouterSubscriptionOwnerTransferRequested)
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
		it.Event = new(FunctionsRouterSubscriptionOwnerTransferRequested)
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

func (it *FunctionsRouterSubscriptionOwnerTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *FunctionsRouterSubscriptionOwnerTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRouterSubscriptionOwnerTransferRequested struct {
	SubscriptionId uint64
	From           common.Address
	To             common.Address
	Raw            types.Log
}

func (_FunctionsRouter *FunctionsRouterFilterer) FilterSubscriptionOwnerTransferRequested(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsRouterSubscriptionOwnerTransferRequestedIterator, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsRouter.contract.FilterLogs(opts, "SubscriptionOwnerTransferRequested", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsRouterSubscriptionOwnerTransferRequestedIterator{contract: _FunctionsRouter.contract, event: "SubscriptionOwnerTransferRequested", logs: logs, sub: sub}, nil
}

func (_FunctionsRouter *FunctionsRouterFilterer) WatchSubscriptionOwnerTransferRequested(opts *bind.WatchOpts, sink chan<- *FunctionsRouterSubscriptionOwnerTransferRequested, subscriptionId []uint64) (event.Subscription, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsRouter.contract.WatchLogs(opts, "SubscriptionOwnerTransferRequested", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRouterSubscriptionOwnerTransferRequested)
				if err := _FunctionsRouter.contract.UnpackLog(event, "SubscriptionOwnerTransferRequested", log); err != nil {
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

func (_FunctionsRouter *FunctionsRouterFilterer) ParseSubscriptionOwnerTransferRequested(log types.Log) (*FunctionsRouterSubscriptionOwnerTransferRequested, error) {
	event := new(FunctionsRouterSubscriptionOwnerTransferRequested)
	if err := _FunctionsRouter.contract.UnpackLog(event, "SubscriptionOwnerTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsRouterSubscriptionOwnerTransferredIterator struct {
	Event *FunctionsRouterSubscriptionOwnerTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRouterSubscriptionOwnerTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRouterSubscriptionOwnerTransferred)
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
		it.Event = new(FunctionsRouterSubscriptionOwnerTransferred)
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

func (it *FunctionsRouterSubscriptionOwnerTransferredIterator) Error() error {
	return it.fail
}

func (it *FunctionsRouterSubscriptionOwnerTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRouterSubscriptionOwnerTransferred struct {
	SubscriptionId uint64
	From           common.Address
	To             common.Address
	Raw            types.Log
}

func (_FunctionsRouter *FunctionsRouterFilterer) FilterSubscriptionOwnerTransferred(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsRouterSubscriptionOwnerTransferredIterator, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsRouter.contract.FilterLogs(opts, "SubscriptionOwnerTransferred", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsRouterSubscriptionOwnerTransferredIterator{contract: _FunctionsRouter.contract, event: "SubscriptionOwnerTransferred", logs: logs, sub: sub}, nil
}

func (_FunctionsRouter *FunctionsRouterFilterer) WatchSubscriptionOwnerTransferred(opts *bind.WatchOpts, sink chan<- *FunctionsRouterSubscriptionOwnerTransferred, subscriptionId []uint64) (event.Subscription, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsRouter.contract.WatchLogs(opts, "SubscriptionOwnerTransferred", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRouterSubscriptionOwnerTransferred)
				if err := _FunctionsRouter.contract.UnpackLog(event, "SubscriptionOwnerTransferred", log); err != nil {
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

func (_FunctionsRouter *FunctionsRouterFilterer) ParseSubscriptionOwnerTransferred(log types.Log) (*FunctionsRouterSubscriptionOwnerTransferred, error) {
	event := new(FunctionsRouterSubscriptionOwnerTransferred)
	if err := _FunctionsRouter.contract.UnpackLog(event, "SubscriptionOwnerTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsRouterTimeLockProposedIterator struct {
	Event *FunctionsRouterTimeLockProposed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRouterTimeLockProposedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRouterTimeLockProposed)
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
		it.Event = new(FunctionsRouterTimeLockProposed)
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

func (it *FunctionsRouterTimeLockProposedIterator) Error() error {
	return it.fail
}

func (it *FunctionsRouterTimeLockProposedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRouterTimeLockProposed struct {
	From uint16
	To   uint16
	Raw  types.Log
}

func (_FunctionsRouter *FunctionsRouterFilterer) FilterTimeLockProposed(opts *bind.FilterOpts) (*FunctionsRouterTimeLockProposedIterator, error) {

	logs, sub, err := _FunctionsRouter.contract.FilterLogs(opts, "TimeLockProposed")
	if err != nil {
		return nil, err
	}
	return &FunctionsRouterTimeLockProposedIterator{contract: _FunctionsRouter.contract, event: "TimeLockProposed", logs: logs, sub: sub}, nil
}

func (_FunctionsRouter *FunctionsRouterFilterer) WatchTimeLockProposed(opts *bind.WatchOpts, sink chan<- *FunctionsRouterTimeLockProposed) (event.Subscription, error) {

	logs, sub, err := _FunctionsRouter.contract.WatchLogs(opts, "TimeLockProposed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRouterTimeLockProposed)
				if err := _FunctionsRouter.contract.UnpackLog(event, "TimeLockProposed", log); err != nil {
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

func (_FunctionsRouter *FunctionsRouterFilterer) ParseTimeLockProposed(log types.Log) (*FunctionsRouterTimeLockProposed, error) {
	event := new(FunctionsRouterTimeLockProposed)
	if err := _FunctionsRouter.contract.UnpackLog(event, "TimeLockProposed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsRouterTimeLockUpdatedIterator struct {
	Event *FunctionsRouterTimeLockUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRouterTimeLockUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRouterTimeLockUpdated)
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
		it.Event = new(FunctionsRouterTimeLockUpdated)
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

func (it *FunctionsRouterTimeLockUpdatedIterator) Error() error {
	return it.fail
}

func (it *FunctionsRouterTimeLockUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRouterTimeLockUpdated struct {
	From uint16
	To   uint16
	Raw  types.Log
}

func (_FunctionsRouter *FunctionsRouterFilterer) FilterTimeLockUpdated(opts *bind.FilterOpts) (*FunctionsRouterTimeLockUpdatedIterator, error) {

	logs, sub, err := _FunctionsRouter.contract.FilterLogs(opts, "TimeLockUpdated")
	if err != nil {
		return nil, err
	}
	return &FunctionsRouterTimeLockUpdatedIterator{contract: _FunctionsRouter.contract, event: "TimeLockUpdated", logs: logs, sub: sub}, nil
}

func (_FunctionsRouter *FunctionsRouterFilterer) WatchTimeLockUpdated(opts *bind.WatchOpts, sink chan<- *FunctionsRouterTimeLockUpdated) (event.Subscription, error) {

	logs, sub, err := _FunctionsRouter.contract.WatchLogs(opts, "TimeLockUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRouterTimeLockUpdated)
				if err := _FunctionsRouter.contract.UnpackLog(event, "TimeLockUpdated", log); err != nil {
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

func (_FunctionsRouter *FunctionsRouterFilterer) ParseTimeLockUpdated(log types.Log) (*FunctionsRouterTimeLockUpdated, error) {
	event := new(FunctionsRouterTimeLockUpdated)
	if err := _FunctionsRouter.contract.UnpackLog(event, "TimeLockUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsRouterUnpausedIterator struct {
	Event *FunctionsRouterUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRouterUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRouterUnpaused)
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
		it.Event = new(FunctionsRouterUnpaused)
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

func (it *FunctionsRouterUnpausedIterator) Error() error {
	return it.fail
}

func (it *FunctionsRouterUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRouterUnpaused struct {
	Account common.Address
	Raw     types.Log
}

func (_FunctionsRouter *FunctionsRouterFilterer) FilterUnpaused(opts *bind.FilterOpts) (*FunctionsRouterUnpausedIterator, error) {

	logs, sub, err := _FunctionsRouter.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &FunctionsRouterUnpausedIterator{contract: _FunctionsRouter.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

func (_FunctionsRouter *FunctionsRouterFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *FunctionsRouterUnpaused) (event.Subscription, error) {

	logs, sub, err := _FunctionsRouter.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRouterUnpaused)
				if err := _FunctionsRouter.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

func (_FunctionsRouter *FunctionsRouterFilterer) ParseUnpaused(log types.Log) (*FunctionsRouterUnpaused, error) {
	event := new(FunctionsRouterUnpaused)
	if err := _FunctionsRouter.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetConsumer struct {
	Allowed           bool
	InitiatedRequests uint64
	CompletedRequests uint64
}
type GetSubscription struct {
	Balance        *big.Int
	BlockedBalance *big.Int
	Owner          common.Address
	RequestedOwner common.Address
	Consumers      []common.Address
}

func (_FunctionsRouter *FunctionsRouter) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _FunctionsRouter.abi.Events["ConfigProposed"].ID:
		return _FunctionsRouter.ParseConfigProposed(log)
	case _FunctionsRouter.abi.Events["ConfigSet"].ID:
		return _FunctionsRouter.ParseConfigSet(log)
	case _FunctionsRouter.abi.Events["ConfigUpdated"].ID:
		return _FunctionsRouter.ParseConfigUpdated(log)
	case _FunctionsRouter.abi.Events["ContractProposed"].ID:
		return _FunctionsRouter.ParseContractProposed(log)
	case _FunctionsRouter.abi.Events["ContractUpdated"].ID:
		return _FunctionsRouter.ParseContractUpdated(log)
	case _FunctionsRouter.abi.Events["FundsRecovered"].ID:
		return _FunctionsRouter.ParseFundsRecovered(log)
	case _FunctionsRouter.abi.Events["OwnershipTransferRequested"].ID:
		return _FunctionsRouter.ParseOwnershipTransferRequested(log)
	case _FunctionsRouter.abi.Events["OwnershipTransferred"].ID:
		return _FunctionsRouter.ParseOwnershipTransferred(log)
	case _FunctionsRouter.abi.Events["Paused"].ID:
		return _FunctionsRouter.ParsePaused(log)
	case _FunctionsRouter.abi.Events["RequestEnd"].ID:
		return _FunctionsRouter.ParseRequestEnd(log)
	case _FunctionsRouter.abi.Events["RequestStart"].ID:
		return _FunctionsRouter.ParseRequestStart(log)
	case _FunctionsRouter.abi.Events["SubscriptionCanceled"].ID:
		return _FunctionsRouter.ParseSubscriptionCanceled(log)
	case _FunctionsRouter.abi.Events["SubscriptionConsumerAdded"].ID:
		return _FunctionsRouter.ParseSubscriptionConsumerAdded(log)
	case _FunctionsRouter.abi.Events["SubscriptionConsumerRemoved"].ID:
		return _FunctionsRouter.ParseSubscriptionConsumerRemoved(log)
	case _FunctionsRouter.abi.Events["SubscriptionCreated"].ID:
		return _FunctionsRouter.ParseSubscriptionCreated(log)
	case _FunctionsRouter.abi.Events["SubscriptionFunded"].ID:
		return _FunctionsRouter.ParseSubscriptionFunded(log)
	case _FunctionsRouter.abi.Events["SubscriptionOwnerTransferRequested"].ID:
		return _FunctionsRouter.ParseSubscriptionOwnerTransferRequested(log)
	case _FunctionsRouter.abi.Events["SubscriptionOwnerTransferred"].ID:
		return _FunctionsRouter.ParseSubscriptionOwnerTransferred(log)
	case _FunctionsRouter.abi.Events["TimeLockProposed"].ID:
		return _FunctionsRouter.ParseTimeLockProposed(log)
	case _FunctionsRouter.abi.Events["TimeLockUpdated"].ID:
		return _FunctionsRouter.ParseTimeLockUpdated(log)
	case _FunctionsRouter.abi.Events["Unpaused"].ID:
		return _FunctionsRouter.ParseUnpaused(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (FunctionsRouterConfigProposed) Topic() common.Hash {
	return common.HexToHash("0x0fcfd32a68209b42944376bc5f4bf72c41ba0f378cf60434f84390b82c9844cf")
}

func (FunctionsRouterConfigSet) Topic() common.Hash {
	return common.HexToHash("0xb2b9bfb736b0bfef266985352cfca2a114226362209f5a6e5a2a2ed14383ba57")
}

func (FunctionsRouterConfigUpdated) Topic() common.Hash {
	return common.HexToHash("0xc07626096b49b0462a576c9a7878cf0675e784bb4832584c1289c11664e55f47")
}

func (FunctionsRouterContractProposed) Topic() common.Hash {
	return common.HexToHash("0x72a33d2f293a0a70fad221bb610d3d6b52aed2d840adae1fa721071fbd290cfd")
}

func (FunctionsRouterContractUpdated) Topic() common.Hash {
	return common.HexToHash("0xf17ce6b5bbeda6598ac50e505071ddbd27c188c5ecaf9126b42c951c9c4b61e9")
}

func (FunctionsRouterFundsRecovered) Topic() common.Hash {
	return common.HexToHash("0x59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b436600")
}

func (FunctionsRouterOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (FunctionsRouterOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (FunctionsRouterPaused) Topic() common.Hash {
	return common.HexToHash("0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258")
}

func (FunctionsRouterRequestEnd) Topic() common.Hash {
	return common.HexToHash("0x526019d853b0b88e828dbdc124746afde5a47053256d2a0872a025ab145b0520")
}

func (FunctionsRouterRequestStart) Topic() common.Hash {
	return common.HexToHash("0x7c720ccd20069b8311a6be4ba1cf3294d09eb247aa5d73a8502054b6e68a2f54")
}

func (FunctionsRouterSubscriptionCanceled) Topic() common.Hash {
	return common.HexToHash("0xe8ed5b475a5b5987aa9165e8731bb78043f39eee32ec5a1169a89e27fcd49815")
}

func (FunctionsRouterSubscriptionConsumerAdded) Topic() common.Hash {
	return common.HexToHash("0x43dc749a04ac8fb825cbd514f7c0e13f13bc6f2ee66043b76629d51776cff8e0")
}

func (FunctionsRouterSubscriptionConsumerRemoved) Topic() common.Hash {
	return common.HexToHash("0x182bff9831466789164ca77075fffd84916d35a8180ba73c27e45634549b445b")
}

func (FunctionsRouterSubscriptionCreated) Topic() common.Hash {
	return common.HexToHash("0x464722b4166576d3dcbba877b999bc35cf911f4eaf434b7eba68fa113951d0bf")
}

func (FunctionsRouterSubscriptionFunded) Topic() common.Hash {
	return common.HexToHash("0xd39ec07f4e209f627a4c427971473820dc129761ba28de8906bd56f57101d4f8")
}

func (FunctionsRouterSubscriptionOwnerTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x69436ea6df009049404f564eff6622cd00522b0bd6a89efd9e52a355c4a879be")
}

func (FunctionsRouterSubscriptionOwnerTransferred) Topic() common.Hash {
	return common.HexToHash("0x6f1dc65165ffffedfd8e507b4a0f1fcfdada045ed11f6c26ba27cedfe87802f0")
}

func (FunctionsRouterTimeLockProposed) Topic() common.Hash {
	return common.HexToHash("0x391ced87fc414ac6e67620b6432a8df0e27da63c8cc6782941fb535ac9e06c0a")
}

func (FunctionsRouterTimeLockUpdated) Topic() common.Hash {
	return common.HexToHash("0xbda3a518f5d5dbadc2e16a279c572f58eeddb114eb084747e012080e64c54d1d")
}

func (FunctionsRouterUnpaused) Topic() common.Hash {
	return common.HexToHash("0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa")
}

func (_FunctionsRouter *FunctionsRouter) Address() common.Address {
	return _FunctionsRouter.address
}

type FunctionsRouterInterface interface {
	MAXCONSUMERS(opts *bind.CallOpts) (uint16, error)

	GetAdminFee(opts *bind.CallOpts) (*big.Int, error)

	GetConfigHash(opts *bind.CallOpts) ([32]byte, error)

	GetConsumer(opts *bind.CallOpts, client common.Address, subscriptionId uint64) (GetConsumer,

		error)

	GetContractById(opts *bind.CallOpts, id [32]byte, useProposed bool) (common.Address, error)

	GetContractById0(opts *bind.CallOpts, id [32]byte) (common.Address, error)

	GetProposedContractSet(opts *bind.CallOpts) (*big.Int, [][32]byte, []common.Address, []common.Address, error)

	GetSubscription(opts *bind.CallOpts, subscriptionId uint64) (GetSubscription,

		error)

	GetSubscriptionCount(opts *bind.CallOpts) (uint64, error)

	GetTotalBalance(opts *bind.CallOpts) (*big.Int, error)

	IsPaused(opts *bind.CallOpts) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	Paused(opts *bind.CallOpts) (bool, error)

	PendingRequestExists(opts *bind.CallOpts, subscriptionId uint64) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	Version(opts *bind.CallOpts) (uint16, uint16, uint16, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AcceptSubscriptionOwnerTransfer(opts *bind.TransactOpts, subscriptionId uint64) (*types.Transaction, error)

	AddConsumer(opts *bind.TransactOpts, subscriptionId uint64, consumer common.Address) (*types.Transaction, error)

	CancelSubscription(opts *bind.TransactOpts, subscriptionId uint64, to common.Address) (*types.Transaction, error)

	CreateSubscription(opts *bind.TransactOpts) (*types.Transaction, error)

	Fulfill(opts *bind.TransactOpts, requestId [32]byte, response []byte, err []byte, juelsPerGas *big.Int, costWithoutFulfillment *big.Int, transmitter common.Address) (*types.Transaction, error)

	OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error)

	OracleWithdraw(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error)

	OwnerCancelSubscription(opts *bind.TransactOpts, subscriptionId uint64) (*types.Transaction, error)

	ProposeConfigUpdate(opts *bind.TransactOpts, id [32]byte, config []byte) (*types.Transaction, error)

	ProposeContractsUpdate(opts *bind.TransactOpts, proposedContractSetIds [][32]byte, proposedContractSetFromAddresses []common.Address, proposedContractSetToAddresses []common.Address) (*types.Transaction, error)

	ProposeTimelockBlocks(opts *bind.TransactOpts, blocks uint16) (*types.Transaction, error)

	RecoverFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	RemoveConsumer(opts *bind.TransactOpts, subscriptionId uint64, consumer common.Address) (*types.Transaction, error)

	RequestSubscriptionOwnerTransfer(opts *bind.TransactOpts, subscriptionId uint64, newOwner common.Address) (*types.Transaction, error)

	SendRequest(opts *bind.TransactOpts, subscriptionId uint64, data []byte, dataVersion uint16, callbackGasLimit uint32, donId [32]byte) (*types.Transaction, error)

	TimeoutRequests(opts *bind.TransactOpts, requestIdsToTimeout [][32]byte) (*types.Transaction, error)

	TogglePaused(opts *bind.TransactOpts) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	UpdateConfig(opts *bind.TransactOpts, id [32]byte) (*types.Transaction, error)

	UpdateContracts(opts *bind.TransactOpts) (*types.Transaction, error)

	UpdateTimelockBlocks(opts *bind.TransactOpts) (*types.Transaction, error)

	ValidateProposedContracts(opts *bind.TransactOpts, id [32]byte, data []byte) (*types.Transaction, error)

	FilterConfigProposed(opts *bind.FilterOpts) (*FunctionsRouterConfigProposedIterator, error)

	WatchConfigProposed(opts *bind.WatchOpts, sink chan<- *FunctionsRouterConfigProposed) (event.Subscription, error)

	ParseConfigProposed(log types.Log) (*FunctionsRouterConfigProposed, error)

	FilterConfigSet(opts *bind.FilterOpts) (*FunctionsRouterConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *FunctionsRouterConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*FunctionsRouterConfigSet, error)

	FilterConfigUpdated(opts *bind.FilterOpts) (*FunctionsRouterConfigUpdatedIterator, error)

	WatchConfigUpdated(opts *bind.WatchOpts, sink chan<- *FunctionsRouterConfigUpdated) (event.Subscription, error)

	ParseConfigUpdated(log types.Log) (*FunctionsRouterConfigUpdated, error)

	FilterContractProposed(opts *bind.FilterOpts) (*FunctionsRouterContractProposedIterator, error)

	WatchContractProposed(opts *bind.WatchOpts, sink chan<- *FunctionsRouterContractProposed) (event.Subscription, error)

	ParseContractProposed(log types.Log) (*FunctionsRouterContractProposed, error)

	FilterContractUpdated(opts *bind.FilterOpts) (*FunctionsRouterContractUpdatedIterator, error)

	WatchContractUpdated(opts *bind.WatchOpts, sink chan<- *FunctionsRouterContractUpdated) (event.Subscription, error)

	ParseContractUpdated(log types.Log) (*FunctionsRouterContractUpdated, error)

	FilterFundsRecovered(opts *bind.FilterOpts) (*FunctionsRouterFundsRecoveredIterator, error)

	WatchFundsRecovered(opts *bind.WatchOpts, sink chan<- *FunctionsRouterFundsRecovered) (event.Subscription, error)

	ParseFundsRecovered(log types.Log) (*FunctionsRouterFundsRecovered, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsRouterOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *FunctionsRouterOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*FunctionsRouterOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsRouterOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FunctionsRouterOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*FunctionsRouterOwnershipTransferred, error)

	FilterPaused(opts *bind.FilterOpts) (*FunctionsRouterPausedIterator, error)

	WatchPaused(opts *bind.WatchOpts, sink chan<- *FunctionsRouterPaused) (event.Subscription, error)

	ParsePaused(log types.Log) (*FunctionsRouterPaused, error)

	FilterRequestEnd(opts *bind.FilterOpts, requestId [][32]byte, subscriptionId []uint64) (*FunctionsRouterRequestEndIterator, error)

	WatchRequestEnd(opts *bind.WatchOpts, sink chan<- *FunctionsRouterRequestEnd, requestId [][32]byte, subscriptionId []uint64) (event.Subscription, error)

	ParseRequestEnd(log types.Log) (*FunctionsRouterRequestEnd, error)

	FilterRequestStart(opts *bind.FilterOpts, requestId [][32]byte, donId [][32]byte, subscriptionId []uint64) (*FunctionsRouterRequestStartIterator, error)

	WatchRequestStart(opts *bind.WatchOpts, sink chan<- *FunctionsRouterRequestStart, requestId [][32]byte, donId [][32]byte, subscriptionId []uint64) (event.Subscription, error)

	ParseRequestStart(log types.Log) (*FunctionsRouterRequestStart, error)

	FilterSubscriptionCanceled(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsRouterSubscriptionCanceledIterator, error)

	WatchSubscriptionCanceled(opts *bind.WatchOpts, sink chan<- *FunctionsRouterSubscriptionCanceled, subscriptionId []uint64) (event.Subscription, error)

	ParseSubscriptionCanceled(log types.Log) (*FunctionsRouterSubscriptionCanceled, error)

	FilterSubscriptionConsumerAdded(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsRouterSubscriptionConsumerAddedIterator, error)

	WatchSubscriptionConsumerAdded(opts *bind.WatchOpts, sink chan<- *FunctionsRouterSubscriptionConsumerAdded, subscriptionId []uint64) (event.Subscription, error)

	ParseSubscriptionConsumerAdded(log types.Log) (*FunctionsRouterSubscriptionConsumerAdded, error)

	FilterSubscriptionConsumerRemoved(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsRouterSubscriptionConsumerRemovedIterator, error)

	WatchSubscriptionConsumerRemoved(opts *bind.WatchOpts, sink chan<- *FunctionsRouterSubscriptionConsumerRemoved, subscriptionId []uint64) (event.Subscription, error)

	ParseSubscriptionConsumerRemoved(log types.Log) (*FunctionsRouterSubscriptionConsumerRemoved, error)

	FilterSubscriptionCreated(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsRouterSubscriptionCreatedIterator, error)

	WatchSubscriptionCreated(opts *bind.WatchOpts, sink chan<- *FunctionsRouterSubscriptionCreated, subscriptionId []uint64) (event.Subscription, error)

	ParseSubscriptionCreated(log types.Log) (*FunctionsRouterSubscriptionCreated, error)

	FilterSubscriptionFunded(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsRouterSubscriptionFundedIterator, error)

	WatchSubscriptionFunded(opts *bind.WatchOpts, sink chan<- *FunctionsRouterSubscriptionFunded, subscriptionId []uint64) (event.Subscription, error)

	ParseSubscriptionFunded(log types.Log) (*FunctionsRouterSubscriptionFunded, error)

	FilterSubscriptionOwnerTransferRequested(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsRouterSubscriptionOwnerTransferRequestedIterator, error)

	WatchSubscriptionOwnerTransferRequested(opts *bind.WatchOpts, sink chan<- *FunctionsRouterSubscriptionOwnerTransferRequested, subscriptionId []uint64) (event.Subscription, error)

	ParseSubscriptionOwnerTransferRequested(log types.Log) (*FunctionsRouterSubscriptionOwnerTransferRequested, error)

	FilterSubscriptionOwnerTransferred(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsRouterSubscriptionOwnerTransferredIterator, error)

	WatchSubscriptionOwnerTransferred(opts *bind.WatchOpts, sink chan<- *FunctionsRouterSubscriptionOwnerTransferred, subscriptionId []uint64) (event.Subscription, error)

	ParseSubscriptionOwnerTransferred(log types.Log) (*FunctionsRouterSubscriptionOwnerTransferred, error)

	FilterTimeLockProposed(opts *bind.FilterOpts) (*FunctionsRouterTimeLockProposedIterator, error)

	WatchTimeLockProposed(opts *bind.WatchOpts, sink chan<- *FunctionsRouterTimeLockProposed) (event.Subscription, error)

	ParseTimeLockProposed(log types.Log) (*FunctionsRouterTimeLockProposed, error)

	FilterTimeLockUpdated(opts *bind.FilterOpts) (*FunctionsRouterTimeLockUpdatedIterator, error)

	WatchTimeLockUpdated(opts *bind.WatchOpts, sink chan<- *FunctionsRouterTimeLockUpdated) (event.Subscription, error)

	ParseTimeLockUpdated(log types.Log) (*FunctionsRouterTimeLockUpdated, error)

	FilterUnpaused(opts *bind.FilterOpts) (*FunctionsRouterUnpausedIterator, error)

	WatchUnpaused(opts *bind.WatchOpts, sink chan<- *FunctionsRouterUnpaused) (event.Subscription, error)

	ParseUnpaused(log types.Log) (*FunctionsRouterUnpaused, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
