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
	ABI: "[{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"timelockBlocks\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"maximumTimelockBlocks\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"linkToken\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"internalBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"externalBalance\",\"type\":\"uint256\"}],\"name\":\"BalanceInvariantViolated\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConsumerRequestsInFlight\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"limit\",\"type\":\"uint32\"}],\"name\":\"GasLimitTooBig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"IdentifierIsReserved\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidConfigData\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"InvalidConsumer\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"value\",\"type\":\"uint8\"}],\"name\":\"InvalidGasFlagValue\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidProposal\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSubscription\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposedOwner\",\"type\":\"address\"}],\"name\":\"MustBeRequestedOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeSubscriptionOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByRoute\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableFromCoordinator\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableFromLink\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingRequestExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ProposedTimelockAboveMaximum\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Reentrant\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"RouteNotFound\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"SenderMustAcceptTermsOfService\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TimelockInEffect\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyConsumers\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"adminFee\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"bytes4\",\"name\":\"handleOracleFulfillmentSelector\",\"type\":\"bytes4\"},{\"indexed\":false,\"internalType\":\"uint32[]\",\"name\":\"maxCallbackGasLimits\",\"type\":\"uint32[]\"}],\"name\":\"ConfigChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"fromHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"toBytes\",\"type\":\"bytes\"}],\"name\":\"ConfigProposed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"fromHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"toBytes\",\"type\":\"bytes\"}],\"name\":\"ConfigUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"proposedContractSetId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"proposedContractSetFromAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"proposedContractSetToAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timelockEndBlock\",\"type\":\"uint256\"}],\"name\":\"ContractProposed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"proposedContractSetId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"proposedContractSetFromAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"proposedContractSetToAddress\",\"type\":\"address\"}],\"name\":\"ContractUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalCostJuels\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"resultCode\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"response\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"}],\"name\":\"RequestEnd\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"donId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"subscriptionOwner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requestingContract\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requestInitiator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"dataVersion\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"}],\"name\":\"RequestStart\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"RequestTimedOut\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"fundsRecipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fundsAmount\",\"type\":\"uint256\"}],\"name\":\"SubscriptionCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"SubscriptionCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"from\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"to\",\"type\":\"uint16\"}],\"name\":\"TimeLockProposed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"from\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"to\",\"type\":\"uint16\"}],\"name\":\"TimeLockUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"MAX_CALLBACK_RETURN_BYTES\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"acceptSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"addConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"cancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createSubscription\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"response\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"juelsPerGas\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"costWithoutFulfillment\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"fulfill\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"resultCode\",\"type\":\"uint8\"},{\"internalType\":\"uint96\",\"name\":\"callbackGasCostJuels\",\"type\":\"uint96\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAdminFee\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowListId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfigHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"config\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"getConsumer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"initiatedRequests\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"completedRequests\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"useProposed\",\"type\":\"bool\"}],\"name\":\"getContractById\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"routeDestination\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"getContractById\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"routeDestination\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"getFlags\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMaxConsumers\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getProposedContractSet\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"},{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"getSubscription\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"blockedBalance\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"requestedOwner\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSubscriptionCount\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTotalBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isPaused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"}],\"name\":\"isValidCallbackGasLimit\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"oracleWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"ownerCancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"ownerWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"pendingRequestExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"name\":\"proposeConfigUpdate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"proposedContractSetIds\",\"type\":\"bytes32[]\"},{\"internalType\":\"address[]\",\"name\":\"proposedContractSetAddresses\",\"type\":\"address[]\"}],\"name\":\"proposeContractsUpdate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"proposeSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"blocks\",\"type\":\"uint16\"}],\"name\":\"proposeTimelockBlocks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"removeConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint16\",\"name\":\"dataVersion\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"donId\",\"type\":\"bytes32\"}],\"name\":\"sendRequest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"flags\",\"type\":\"bytes32\"}],\"name\":\"setFlags\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"requestIdsToTimeout\",\"type\":\"bytes32[]\"}],\"name\":\"timeoutRequests\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"updateConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"updateContracts\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"updateTimelockBlocks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"validateProposedContracts\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b50604051620062b9380380620062b9833981016040819052620000349162000431565b6000805460ff191681558290339086908690859084908190816200009f5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0380851661010002610100600160a81b031990921691909117909155811615620000d957620000d98162000178565b50506009805461ffff80871661ffff19909216919091179091558316608052506000805260026020527fac33ff75c19e70fe83507db0d683fd3465c996598dc972688b7ace676c89077b80546001600160a01b031916301790556200013e8162000229565b80516020909101206007555050600a80546001600160a01b0319166001600160a01b03939093169290921790915550620006ab9350505050565b336001600160a01b03821603620001d25760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000096565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929361010090910416917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000806000838060200190518101906200024491906200052e565b604080516060810182526001600160601b0385168082526001600160e01b03198516602080840191909152928201849052600f80546001600160801b0319169091176c0100000000000000000000000060e087901c0217815583519598509396509194509092620002bc916010919086019062000303565b509050507fe6a1eda76d42a6d1d813f26765716562044db1e8bd8be7d088705e64afd301ca838383604051620002f5939291906200063a565b60405180910390a150505050565b82805482825590600052602060002090600701600890048101928215620003a75791602002820160005b838211156200037357835183826101000a81548163ffffffff021916908363ffffffff16021790555092602001926004016020816003010492830192600103026200032d565b8015620003a55782816101000a81549063ffffffff021916905560040160208160030104928301926001030262000373565b505b50620003b5929150620003b9565b5090565b5b80821115620003b55760008155600101620003ba565b805161ffff81168114620003e357600080fd5b919050565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f191681016001600160401b0381118282101715620004295762000429620003e8565b604052919050565b600080600080608085870312156200044857600080fd5b6200045385620003d0565b9350602062000464818701620003d0565b60408701519094506001600160a01b03811681146200048257600080fd5b60608701519093506001600160401b0380821115620004a057600080fd5b818801915088601f830112620004b557600080fd5b815181811115620004ca57620004ca620003e8565b620004de601f8201601f19168501620003fe565b91508082528984828501011115620004f557600080fd5b60005b8181101562000515578381018501518382018601528401620004f8565b5060008482840101525080935050505092959194509250565b6000806000606084860312156200054457600080fd5b83516001600160601b03811681146200055c57600080fd5b602085810151919450906001600160e01b0319811681146200057d57600080fd5b60408601519093506001600160401b03808211156200059b57600080fd5b818701915087601f830112620005b057600080fd5b815181811115620005c557620005c5620003e8565b8060051b9150620005d8848301620003fe565b818152918301840191848101908a841115620005f357600080fd5b938501935b838510156200062a578451925063ffffffff83168314620006195760008081fd5b8282529385019390850190620005f8565b8096505050505050509250925092565b6001600160601b03841681526001600160e01b031983166020808301919091526060604083018190528351908301819052600091848101916080850190845b818110156200069d57845163ffffffff168352938301939183019160010162000679565b509098975050505050505050565b608051615bf2620006c7600039600061292f0152615bf26000f3fe608060405234801561001057600080fd5b50600436106102ff5760003560e01c80637341c10c1161019c578063a9c9a918116100ee578063c4614d6111610097578063e82ad7d411610071578063e82ad7d414610766578063eb523d6c14610779578063f2fde38b1461078c57600080fd5b8063c4614d611461072d578063d7ae1d3014610740578063e72f6e301461075357600080fd5b8063b5643858116100c8578063b5643858146106fb578063b734c0f41461070e578063badc3eb61461071657600080fd5b8063a9c9a918146106b7578063aab396bd146106ca578063b187bd26146106f057600080fd5b80638fde531711610150578063a21a23e41161012a578063a21a23e414610678578063a47c769614610680578063a4c0ed36146106a457600080fd5b80638fde5317146106555780639883c10d1461065d5780639f87fad71461066557600080fd5b8063823597401161018157806382359740146106175780638456cb591461062a5780638da5cb5b1461063257600080fd5b80637341c10c146105fc57806379ba50971461060f57600080fd5b8063461d27621161025557806366316d8d11610209578063674603d0116101e3578063674603d0146105285780636a6df79b146105bd57806371ec28ac146105f557600080fd5b806366316d8d146104d457806366419970146104e7578063665871ec1461051557600080fd5b806355fedefa1161023a57806355fedefa146104975780635c975abb146104aa5780635ed6dfba146104c157600080fd5b8063461d2762146104635780634b8832d31461048457600080fd5b80631ded3b36116102b75780633e871e4d116102915780633e871e4d1461040f5780633f4ba83a146104225780633fd67e051461042a57600080fd5b80631ded3b36146103d35780632a905ccc146103e6578063385de9ae146103fc57600080fd5b806310fc49c1116102e857806310fc49c11461033957806312b583491461034c578063181f5a771461038a57600080fd5b806302bcc5b6146103045780630c5d49cb14610319575b600080fd5b610317610312366004614b5f565b61079f565b005b610321608481565b60405161ffff90911681526020015b60405180910390f35b610317610347366004614b8c565b61082c565b6009546b01000000000000000000000090046bffffffffffffffffffffffff165b6040516bffffffffffffffffffffffff9091168152602001610330565b6103c66040518060400160405280601781526020017f46756e6374696f6e7320526f757465722076312e302e3000000000000000000081525081565b6040516103309190614c27565b6103176103e1366004614c3a565b610836565b600f546bffffffffffffffffffffffff1661036d565b61031761040a366004614cad565b610867565b61031761041d366004614e31565b610a45565b610317610d40565b61043d610438366004614f92565b610d52565b6040805160ff90931683526bffffffffffffffffffffffff909116602083015201610330565b610476610471366004615046565b6110db565b604051908152602001610330565b6103176104923660046150c8565b61112f565b6104766104a5366004614b5f565b61122c565b60005460ff165b6040519015158152602001610330565b6103176104cf3660046150fb565b611250565b6103176104e23660046150fb565b611526565b6009546301000000900467ffffffffffffffff165b60405167ffffffffffffffff9091168152602001610330565b610317610523366004615127565b61174d565b61059561053636600461519c565b73ffffffffffffffffffffffffffffffffffffffff919091166000908152600c6020908152604080832067ffffffffffffffff948516845290915290205460ff8116926101008204831692690100000000000000000090920490911690565b60408051931515845267ffffffffffffffff9283166020850152911690820152606001610330565b6105d06105cb3660046151d4565b611b36565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610330565b6064610321565b61031761060a3660046150c8565b611b49565b610317611cf6565b610317610625366004614b5f565b611e18565b610317611f67565b600054610100900473ffffffffffffffffffffffffffffffffffffffff166105d0565b610317611f77565b600754610476565b6103176106733660046150c8565b61201f565b6104fc612427565b61069361068e366004614b5f565b6125bb565b60405161033095949392919061524a565b6103176106b236600461529f565b6126a7565b6105d06106c53660046152f9565b6128d6565b7fd8e0666292c202b1ce6a8ff0dd638652e662402ac53fbf9bd9d3bcc39d5eb097610476565b60005460ff166104b1565b610317610709366004615312565b6128e3565b610317612a44565b61071e612bb2565b6040516103309392919061532d565b6103c661073b366004614cad565b612c8d565b61031761074e3660046150c8565b612ca2565b61031761076136600461538e565b612cfb565b6104b1610774366004614b5f565b612f07565b6103176107873660046152f9565b612f12565b61031761079a36600461538e565b613130565b6107a7613141565b67ffffffffffffffff81166000908152600b60205260409020546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff168061081e576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6108288282613149565b5050565b61082882826134b7565b61083e613141565b610847826135d7565b67ffffffffffffffff9091166000908152600b6020526040902060030155565b61086f61364d565b600061087c8460006136d3565b905060003073ffffffffffffffffffffffffffffffffffffffff8316036108a6575060075461091a565b8173ffffffffffffffffffffffffffffffffffffffff16639883c10d6040518163ffffffff1660e01b81526004016020604051808303816000875af11580156108f3573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061091791906153a9565b90505b838360405161092a9291906153c2565b6040518091039020810361096a576040517fee03280800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b604051806060016040528082815260200185858080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152505050908252506009546020909101906109cb9061ffff1643615401565b90526000868152600660209081526040909120825181559082015160018201906109f590826154b5565b50604091820151600290910155517f0fcfd32a68209b42944376bc5f4bf72c41ba0f378cf60434f84390b82c9844cf90610a369087908490889088906155cf565b60405180910390a15050505050565b610a4d61364d565b8151815181141580610a5f5750600881115b15610a96576040517fee03280800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b81811015610bc7576000848281518110610ab557610ab561562a565b602002602001015190506000848381518110610ad357610ad361562a565b60200260200101519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161480610b3e575060008281526002602052604090205473ffffffffffffffffffffffffffffffffffffffff8281169116145b15610b75576040517fee03280800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b81610bb4576040517f4855c288000000000000000000000000000000000000000000000000000000008152600481018390526024015b60405180910390fd5b505080610bc090615659565b9050610a99565b50600954600090610bdc9061ffff1643615401565b604080516060810182528681526020808201879052918101839052865192935091600391610c0e9183918901906149a8565b506020828101518051610c2792600185019201906149ef565b506040820151816002015590505060005b8451811015610d39577f72a33d2f293a0a70fad221bb610d3d6b52aed2d840adae1fa721071fbd290cfd858281518110610c7457610c7461562a565b602002602001015160026000888581518110610c9257610c9261562a565b6020026020010151815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16868481518110610cdb57610cdb61562a565b602002602001015185604051610d21949392919093845273ffffffffffffffffffffffffffffffffffffffff928316602085015291166040830152606082015260800190565b60405180910390a1610d3281615659565b9050610c38565b5050505050565b610d4861364d565b610d506137ca565b565b600080610d5d613847565b6000888152600e602090815260409182902082516101008101845281546bffffffffffffffffffffffff80821683526c010000000000000000000000009182900473ffffffffffffffffffffffffffffffffffffffff90811695840186905260018501549081169684019690965274010000000000000000000000000000000000000000860467ffffffffffffffff1660608401527c010000000000000000000000000000000000000000000000000000000090950463ffffffff16608083015260029092015493841660a082015290830464ffffffffff1660c0820152710100000000000000000000000000000000009092046effffffffffffffffffffffffffffff1660e08301523314610e9f576040517f8bec23e700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b604081015173ffffffffffffffffffffffffffffffffffffffff16610ec85760029250506110d0565b8060e00151816080015163ffffffff16610ee29190615691565b6effffffffffffffffffffffffffffff165a1015610f045760039250506110d0565b606081015167ffffffffffffffff166000908152600b602052604090205460808201516bffffffffffffffffffffffff90911690610f479063ffffffff1661384f565b610f5190886156c0565b8251610f5e9088906156f0565b610f6891906156f0565b6bffffffffffffffffffffffff161115610f865760049250506110d0565b8060a001516bffffffffffffffffffffffff16610fac826080015163ffffffff1661384f565b610fb690886156c0565b8251610fc39088906156f0565b610fcd91906156f0565b6bffffffffffffffffffffffff161115610feb5760059250506110d0565b6000898152600e6020526040808220828155600181018390556002018290556080830151908301519194508491611027918c918c918c916138f1565b805190915061103757600161103a565b60005b9350600061106183606001518460a00151856040015186600001518c87602001518d613a6b565b9050826060015167ffffffffffffffff168b7f45bb48b6ec798595a260f114720360b95cc58c94c6ddd37a1acc3896ec94a23a8360200151898987600001516110aa578e6110ac565b8f5b88604001516040516110c2959493929190615715565b60405180910390a351925050505b965096945050505050565b60006111248260008989898080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508b92508a9150613cbd9050565b979650505050505050565b611137613847565b6111408261411e565b6111486141e4565b67ffffffffffffffff82166000908152600b602052604090206001015473ffffffffffffffffffffffffffffffffffffffff8281166c0100000000000000000000000090920416146108285767ffffffffffffffff82166000818152600b602090815260409182902060010180546bffffffffffffffffffffffff166c0100000000000000000000000073ffffffffffffffffffffffffffffffffffffffff8716908102919091179091558251338152918201527f69436ea6df009049404f564eff6622cd00522b0bd6a89efd9e52a355c4a879be91015b60405180910390a25050565b67ffffffffffffffff81166000908152600b60205260408120600301545b92915050565b611258613141565b806bffffffffffffffffffffffff1660000361128e5750306000908152600d60205260409020546bffffffffffffffffffffffff165b306000908152600d60205260409020546bffffffffffffffffffffffff808316911610156112e8576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b306000908152600d6020526040812080548392906113159084906bffffffffffffffffffffffff16615773565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550806009600b8282829054906101000a90046bffffffffffffffffffffffff1661136c9190615773565b82546101009290920a6bffffffffffffffffffffffff818102199093169183160217909155600a546040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff868116600483015292851660248201529116915063a9059cbb906044016020604051808303816000875af115801561140b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061142f9190615798565b61082857600a546040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015260009173ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa1580156114a2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906114c691906153a9565b6009546040517fa99da3020000000000000000000000000000000000000000000000000000000081526b0100000000000000000000009091046bffffffffffffffffffffffff166004820181905260248201839052919250604401610bab565b61152e613847565b806bffffffffffffffffffffffff16600003611576576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b336000908152600d60205260409020546bffffffffffffffffffffffff808316911610156115d0576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b336000908152600d6020526040812080548392906115fd9084906bffffffffffffffffffffffff16615773565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550806009600b8282829054906101000a90046bffffffffffffffffffffffff166116549190615773565b82546101009290920a6bffffffffffffffffffffffff818102199093169183160217909155600a546040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff868116600483015292851660248201529116915063a9059cbb906044016020604051808303816000875af11580156116f3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906117179190615798565b610828576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611755613847565b60005b81811015611b315760008383838181106117745761177461562a565b602090810292909201356000818152600e845260409081902081516101008101835281546bffffffffffffffffffffffff80821683526c010000000000000000000000009182900473ffffffffffffffffffffffffffffffffffffffff9081169884019890985260018401549788169483019490945274010000000000000000000000000000000000000000870467ffffffffffffffff16606083018190527c010000000000000000000000000000000000000000000000000000000090970463ffffffff16608083015260029092015492831660a082015290820464ffffffffff1660c0820152710100000000000000000000000000000000009091046effffffffffffffffffffffffffffff1660e0820152909350919050611897816135d7565b67ffffffffffffffff81166000908152600b60205260409020546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff16338114611910576040517f5a68151d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8260c0015164ffffffffff16421015611955576040517fbcc4005500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60208301516040517f85b214cf0000000000000000000000000000000000000000000000000000000081526004810186905273ffffffffffffffffffffffffffffffffffffffff8216906385b214cf906024016020604051808303816000875af11580156119c7573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906119eb9190615798565b5060a084015167ffffffffffffffff84166000908152600b602052604081206001018054909190611a2b9084906bffffffffffffffffffffffff16615773565b82546bffffffffffffffffffffffff9182166101009390930a92830291909202199091161790555060408085015173ffffffffffffffffffffffffffffffffffffffff166000908152600c602090815282822067ffffffffffffffff808816845291529190208054600192600991611ab291859169010000000000000000009004166157b5565b825467ffffffffffffffff9182166101009390930a9283029190920219909116179055506000858152600e6020526040808220828155600181018390556002018290555186917ff1ca1e9147be737b04a2b018a79405f687a97de8dd8a2559bbe62357343af41491a2505050505080611b2a90615659565b9050611758565b505050565b6000611b4283836136d3565b9392505050565b611b51613847565b611b5a8261411e565b611b626141e4565b67ffffffffffffffff82166000908152600b60205260409020600201547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9c01611bd7576040517f05a48e0f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff81166000908152600c6020908152604080832067ffffffffffffffff8616845290915290205460ff1615611c1e575050565b73ffffffffffffffffffffffffffffffffffffffff81166000818152600c6020908152604080832067ffffffffffffffff871680855290835281842080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001908117909155600b84528285206002018054918201815585529383902090930180547fffffffffffffffffffffffff000000000000000000000000000000000000000016851790555192835290917f43dc749a04ac8fb825cbd514f7c0e13f13bc6f2ee66043b76629d51776cff8e09101611220565b60015473ffffffffffffffffffffffffffffffffffffffff163314611d77576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610bab565b60008054336101008181027fffffffffffffffffffffff0000000000000000000000000000000000000000ff8416178455600180547fffffffffffffffffffffffff000000000000000000000000000000000000000016905560405173ffffffffffffffffffffffffffffffffffffffff919093041692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b611e20613847565b611e286141e4565b67ffffffffffffffff81166000908152600b60205260409020805460019091015473ffffffffffffffffffffffffffffffffffffffff6c010000000000000000000000009283900481169290910416338114611ec8576040517fd084e97500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82166004820152602401610bab565b67ffffffffffffffff83166000818152600b602090815260409182902080546c01000000000000000000000000339081026bffffffffffffffffffffffff928316178355600190920180549091169055825173ffffffffffffffffffffffffffffffffffffffff87168152918201527f6f1dc65165ffffedfd8e507b4a0f1fcfdada045ed11f6c26ba27cedfe87802f0910160405180910390a2505050565b611f6f61364d565b610d506142eb565b611f7f61364d565b60085464010000000090047bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16431015611fe1576040517fa93d035c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600854600980546201000090920461ffff167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909216919091179055565b612027613847565b6120308261411e565b6120386141e4565b73ffffffffffffffffffffffffffffffffffffffff81166000908152600c6020908152604080832067ffffffffffffffff8087168552908352928190208151606081018352905460ff8116151580835261010082048616948301949094526901000000000000000000900490931690830152612108576040517ff0019fe600000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8416600482015273ffffffffffffffffffffffffffffffffffffffff83166024820152604401610bab565b806040015167ffffffffffffffff16816020015167ffffffffffffffff161461215d576040517fbcc4005500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff83166000908152600b60209081526040808320600201805482518185028101850190935280835291929091908301828280156121d857602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116121ad575b505050505090506000600182516121ef91906157d6565b905060005b825181101561238b578473ffffffffffffffffffffffffffffffffffffffff168382815181106122265761222661562a565b602002602001015173ffffffffffffffffffffffffffffffffffffffff160361237b57600083838151811061225d5761225d61562a565b6020026020010151905080600b60008967ffffffffffffffff1667ffffffffffffffff16815260200190815260200160002060020183815481106122a3576122a361562a565b600091825260208083209190910180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff949094169390931790925567ffffffffffffffff89168152600b9091526040902060020180548061231d5761231d6157e9565b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff00000000000000000000000000000000000000001690550190555061238b565b61238481615659565b90506121f4565b5073ffffffffffffffffffffffffffffffffffffffff84166000818152600c6020908152604080832067ffffffffffffffff8a168085529083529281902080547fffffffffffffffffffffffffffffff00000000000000000000000000000000001690555192835290917f182bff9831466789164ca77075fffd84916d35a8180ba73c27e45634549b445b910160405180910390a25050505050565b6000612431613847565b6124396141e4565b6009805460039061245a906301000000900467ffffffffffffffff16615818565b825467ffffffffffffffff8083166101009490940a93840293021916919091179091556040805160c08101825260008082523360208301529181018290526060810182905291925060808201906040519080825280602002602001820160405280156124d0578160200160208202803683370190505b5081526000602091820181905267ffffffffffffffff84168152600b825260409081902083518484015173ffffffffffffffffffffffffffffffffffffffff9081166c010000000000000000000000009081026bffffffffffffffffffffffff938416178455938601516060870151909116909302921691909117600182015560808301518051919261256b926002850192909101906149ef565b5060a0919091015160039091015560405133815267ffffffffffffffff8216907f464722b4166576d3dcbba877b999bc35cf911f4eaf434b7eba68fa113951d0bf9060200160405180910390a290565b60008060008060606125cc866135d7565b67ffffffffffffffff86166000908152600b602090815260409182902080546001820154600290920180548551818602810186019096528086526bffffffffffffffffffffffff8084169b508416995073ffffffffffffffffffffffffffffffffffffffff6c0100000000000000000000000093849004811699509290930490911695509183018282801561269757602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff16815260019091019060200180831161266c575b5050505050905091939590929450565b6126af613847565b600a5473ffffffffffffffffffffffffffffffffffffffff163314612700576040517f44b0e3c300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6020811461273a576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600061274882840184614b5f565b67ffffffffffffffff81166000908152600b60205260409020549091506c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff166127c1576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff81166000908152600b6020526040812080546bffffffffffffffffffffffff16918691906127f883856156f0565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550846009600b8282829054906101000a90046bffffffffffffffffffffffff1661284f91906156f0565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055508167ffffffffffffffff167fd39ec07f4e209f627a4c427971473820dc129761ba28de8906bd56f57101d4f88287846128b69190615401565b6040805192835260208301919091520160405180910390a2505050505050565b600061124a8260006136d3565b6128eb61364d565b60095461ffff80831691160361292d576040517fee03280800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000061ffff168161ffff16111561298f576040517fe9a3062200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805160608101825260095461ffff908116808352908416602083015290918201906129bc9043615401565b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff90811690915281516008805460208501516040909501519093166401000000000263ffffffff61ffff95861662010000027fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000009095169590931694909417929092171691909117905550565b612a4c61364d565b600554431015612a88576040517fa93d035c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b600354811015612baf57600060036000018281548110612aad57612aad61562a565b6000918252602080832090910154808352600290915260408220546004805492945073ffffffffffffffffffffffffffffffffffffffff909116929185908110612af957612af961562a565b6000918252602080832091909101548583526002825260409283902080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92831690811790915583518781529186169282019290925291820181905291507ff8a6175bca1ba37d682089187edc5e20a859989727f10ca6bd9a5bc0de8caf949060600160405180910390a150505080612ba890615659565b9050612a8b565b50565b60006060806003600201546003600001600360010181805480602002602001604051908101604052809291908181526020018280548015612c1257602002820191906000526020600020905b815481526020019060010190808311612bfe575b5050505050915080805480602002602001604051908101604052809291908181526020018280548015612c7b57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311612c50575b50505050509050925092509250909192565b6060612c9a848484614346565b949350505050565b612caa613847565b612cb38261411e565b612cbb6141e4565b612cc482614413565b1561081e576040517fb42f66e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b612d03613141565b600a546040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015260009173ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa158015612d72573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612d9691906153a9565b6009549091506b01000000000000000000000090046bffffffffffffffffffffffff1681811115612dfd576040517fa99da3020000000000000000000000000000000000000000000000000000000081526004810182905260248101839052604401610bab565b81811015611b31576000612e1182846157d6565b600a546040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff87811660048301526024820184905292935091169063a9059cbb906044016020604051808303816000875af1158015612e8c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612eb09190615798565b506040805173ffffffffffffffffffffffffffffffffffffffff86168152602081018390527f59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b43660091015b60405180910390a150505050565b600061124a82614413565b612f1a61364d565b60006006600083815260200190815260200160002060405180606001604052908160008201548152602001600182018054612f5490615414565b80601f0160208091040260200160405190810160405280929190818152602001828054612f8090615414565b8015612fcd5780601f10612fa257610100808354040283529160200191612fcd565b820191906000526020600020905b815481529060010190602001808311612fb057829003601f168201915b5050505050815260200160028201548152505090508060400151431015613020576040517fa93d035c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b81613045576130328160200151614562565b60208082015180519101206007556130ed565b6130508260006136d3565b73ffffffffffffffffffffffffffffffffffffffff16638cc6acce82602001516040518263ffffffff1660e01b815260040161308c9190614c27565b600060405180830381600087803b1580156130a657600080fd5b505af19250505080156130b7575060015b6130ed576040517ffe680b2600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b805160208201516040517fc07626096b49b0462a576c9a7878cf0675e784bb4832584c1289c11664e55f479261312492869261583f565b60405180910390a15050565b61313861364d565b612baf8161465d565b610d5061364d565b67ffffffffffffffff82166000908152600b60209081526040808320815160c08101835281546bffffffffffffffffffffffff808216835273ffffffffffffffffffffffffffffffffffffffff6c010000000000000000000000009283900481168488015260018501549182168487015291900416606082015260028201805484518187028101870190955280855291949293608086019390929083018282801561322a57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116131ff575b505050918352505060039190910154602090910152805190915060005b8260800151518110156132eb57600c60008460800151838151811061326e5761326e61562a565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040908101600090812067ffffffffffffffff89168252909252902080547fffffffffffffffffffffffffffffff00000000000000000000000000000000001690556132e481615659565b9050613247565b5067ffffffffffffffff84166000908152600b60205260408120818155600181018290559061331d6002830182614a69565b60038201600090555050806009600b8282829054906101000a90046bffffffffffffffffffffffff166133509190615773565b82546101009290920a6bffffffffffffffffffffffff818102199093169183160217909155600a546040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff878116600483015292851660248201529116915063a9059cbb906044016020604051808303816000875af11580156133ef573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906134139190615798565b613449576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805173ffffffffffffffffffffffffffffffffffffffff851681526bffffffffffffffffffffffff8316602082015267ffffffffffffffff8616917fe8ed5b475a5b5987aa9165e8731bb78043f39eee32ec5a1169a89e27fcd49815910160405180910390a250505050565b67ffffffffffffffff82166000908152600b6020526040812060030154601054911a908110613517576040517f45c108ce00000000000000000000000000000000000000000000000000000000815260ff82166004820152602401610bab565b6010805460ff831690811061352e5761352e61562a565b90600052602060002090600891828204019190066004029054906101000a900463ffffffff1663ffffffff168263ffffffff161115611b31576010805460ff831690811061357e5761357e61562a565b600091825260209091206008820401546040517f1d70f87a000000000000000000000000000000000000000000000000000000008152600790921660049081026101000a90910463ffffffff1690820152602401610bab565b67ffffffffffffffff81166000908152600b60205260409020546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff16612baf576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600054610100900473ffffffffffffffffffffffffffffffffffffffff163314610d50576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610bab565b6000816137125760008381526002602052604090205473ffffffffffffffffffffffffffffffffffffffff16801561370c57905061124a565b50613795565b60005b6003548110156137935760038054829081106137335761373361562a565b906000526020600020015484036137835760048054829081106137585761375861562a565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff16915061124a9050565b61378c81615659565b9050613715565b505b6040517f80833e3300000000000000000000000000000000000000000000000000000000815260048101849052602401610bab565b6137d2614758565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa335b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390a1565b610d506147c4565b60006bffffffffffffffffffffffff8211156138ed576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f53616665436173743a2076616c756520646f65736e27742066697420696e203960448201527f36206269747300000000000000000000000000000000000000000000000000006064820152608401610bab565b5090565b60408051606080820183526000808352602083015291810191909152600f546040516000916c01000000000000000000000000900460e01b9061393c90899089908990602401615867565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00181529181526020820180517fffffffff00000000000000000000000000000000000000000000000000000000949094167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff909416939093179092528151608480825260c0820190935290925060009182918291602082018180368337019050509050853b6139ee57600080fd5b5a6113888110156139fe57600080fd5b611388810390508760408204820311613a1657600080fd5b60008086516020880160008b8df193505a900391503d6084811115613a39575060845b808252806000602084013e50604080516060810182529315158452602084019290925290820152979650505050505050565b60408051808201909152600080825260208201526000613a8a8461384f565b613a9490866156c0565b9050600081613aa388866156f0565b613aad91906156f0565b6040805180820182526bffffffffffffffffffffffff808616825280841660208084019190915267ffffffffffffffff8f166000908152600b909152928320805492975093945084939291613b0491859116615773565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055508184613b3e91906156f0565b336000908152600d602052604081208054909190613b6b9084906bffffffffffffffffffffffff166156f0565b82546101009290920a6bffffffffffffffffffffffff818102199093169183160217909155306000908152600d6020526040812080548b94509092613bb2918591166156f0565b82546101009290920a6bffffffffffffffffffffffff81810219909316918316021790915567ffffffffffffffff8c166000908152600b6020526040812060010180548d94509092613c0691859116615773565b82546bffffffffffffffffffffffff9182166101009390930a92830291909202199091161790555073ffffffffffffffffffffffffffffffffffffffff88166000908152600c6020908152604080832067ffffffffffffffff808f16855292529091208054600192600991613c8a91859169010000000000000000009004166157b5565b92506101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055505050979650505050505050565b6000613cc7613847565b613cd0856135d7565b613cda3386614831565b613ce485836134b7565b6000613cf088886136d3565b6040805160e08101825233815267ffffffffffffffff89166000818152600b602081815285832080546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff908116838801528688018e90526060870186905261ffff8d1660808801529484529190526003015460a084015263ffffffff881660c084015292517faecc12c50000000000000000000000000000000000000000000000000000000081529394509192839283929086169163aecc12c591613dbb9160040161589c565b6080604051808303816000875af1158015613dda573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613dfe919061592f565b92975090945092509050613e13338a856148cd565b6040805161010081018252600f546bffffffffffffffffffffffff908116825273ffffffffffffffffffffffffffffffffffffffff87166020830152339282019290925267ffffffffffffffff8b16606082015263ffffffff8816608082015290841660a082015260c08101613e898342615401565b64ffffffffff168152602001836effffffffffffffffffffffffffffff16815250600e600087815260200190815260200160002060008201518160000160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550602082015181600001600c6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060408201518160010160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060608201518160010160146101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550608082015181600101601c6101000a81548163ffffffff021916908363ffffffff16021790555060a08201518160020160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060c082015181600201600c6101000a81548164ffffffffff021916908364ffffffffff16021790555060e08201518160020160116101000a8154816effffffffffffffffffffffffffffff02191690836effffffffffffffffffffffffffffff1602179055509050508867ffffffffffffffff168b867f7c720ccd20069b8311a6be4ba1cf3294d09eb247aa5d73a8502054b6e68a2f54600b60008e67ffffffffffffffff1667ffffffffffffffff168152602001908152602001600020600001600c9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1633328e8e8e6040516141089695949392919061596d565b60405180910390a4505050509695505050505050565b67ffffffffffffffff81166000908152600b60205260409020546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1680614195576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff821614610828576040517f5a68151d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61420f7fd8e0666292c202b1ce6a8ff0dd638652e662402ac53fbf9bd9d3bcc39d5eb09760006136d3565b604080516000815260208101918290527f6b14daf80000000000000000000000000000000000000000000000000000000090915273ffffffffffffffffffffffffffffffffffffffff9190911690636b14daf890614272903390602481016159d1565b602060405180830381865afa15801561428f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906142b39190615798565b610d50576040517f22906263000000000000000000000000000000000000000000000000000000008152336004820152602401610bab565b6142f36147c4565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a25861381d3390565b6060600080808061435986880188615a00565b9350935093509350600061437289600187878787613cbd565b60408051602080825281830190925291925060208201818036833701905050955060005b6020811015614406578181602081106143b1576143b161562a565b1a60f81b8782815181106143c7576143c761562a565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053506143ff81615659565b9050614396565b5050505050509392505050565b67ffffffffffffffff81166000908152600b602090815260408083206002018054825181850281018501909352808352849383018282801561448b57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311614460575b5050505050905060005b8151811015614558576000600c60008484815181106144b6576144b661562a565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040908101600090812067ffffffffffffffff808a168352908452908290208251606081018452905460ff811615158252610100810483169482018590526901000000000000000000900490911691810182905292501461454757506001949350505050565b5061455181615659565b9050614495565b5060009392505050565b60008060008380602001905181019061457b9190615a71565b604080516060810182526bffffffffffffffffffffffff85168082527fffffffff000000000000000000000000000000000000000000000000000000008516602080840191909152928201849052600f80547fffffffffffffffffffffffffffffffff00000000000000000000000000000000169091176c0100000000000000000000000060e087901c02178155835195985093965091945090926146269160109190860190614a87565b509050507fe6a1eda76d42a6d1d813f26765716562044db1e8bd8be7d088705e64afd301ca838383604051612ef993929190615b5a565b3373ffffffffffffffffffffffffffffffffffffffff8216036146dc576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610bab565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929361010090910416917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60005460ff16610d50576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f5061757361626c653a206e6f74207061757365640000000000000000000000006044820152606401610bab565b60005460ff1615610d50576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f5061757361626c653a20706175736564000000000000000000000000000000006044820152606401610bab565b73ffffffffffffffffffffffffffffffffffffffff82166000908152600c6020908152604080832067ffffffffffffffff8516845290915290205460ff16610828576040517ff0019fe600000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8216600482015273ffffffffffffffffffffffffffffffffffffffff83166024820152604401610bab565b67ffffffffffffffff82166000908152600b6020526040812060010180548392906149079084906bffffffffffffffffffffffff166156f0565b82546bffffffffffffffffffffffff91821661010093840a908102920219161790915573ffffffffffffffffffffffffffffffffffffffff85166000908152600c6020908152604080832067ffffffffffffffff808916855292529091208054600194509092849261497d9284929004166157b5565b92506101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550505050565b8280548282559060005260206000209081019282156149e3579160200282015b828111156149e35782518255916020019190600101906149c8565b506138ed929150614b2d565b8280548282559060005260206000209081019282156149e3579160200282015b828111156149e357825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190614a0f565b5080546000825590600052602060002090810190612baf9190614b2d565b828054828255906000526020600020906007016008900481019282156149e35791602002820160005b83821115614af457835183826101000a81548163ffffffff021916908363ffffffff1602179055509260200192600401602081600301049283019260010302614ab0565b8015614b245782816101000a81549063ffffffff0219169055600401602081600301049283019260010302614af4565b50506138ed9291505b5b808211156138ed5760008155600101614b2e565b803567ffffffffffffffff81168114614b5a57600080fd5b919050565b600060208284031215614b7157600080fd5b611b4282614b42565b63ffffffff81168114612baf57600080fd5b60008060408385031215614b9f57600080fd5b614ba883614b42565b91506020830135614bb881614b7a565b809150509250929050565b6000815180845260005b81811015614be957602081850181015186830182015201614bcd565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b602081526000611b426020830184614bc3565b60008060408385031215614c4d57600080fd5b614c5683614b42565b946020939093013593505050565b60008083601f840112614c7657600080fd5b50813567ffffffffffffffff811115614c8e57600080fd5b602083019150836020828501011115614ca657600080fd5b9250929050565b600080600060408486031215614cc257600080fd5b83359250602084013567ffffffffffffffff811115614ce057600080fd5b614cec86828701614c64565b9497909650939450505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715614d6f57614d6f614cf9565b604052919050565b600067ffffffffffffffff821115614d9157614d91614cf9565b5060051b60200190565b803573ffffffffffffffffffffffffffffffffffffffff81168114614b5a57600080fd5b600082601f830112614dd057600080fd5b81356020614de5614de083614d77565b614d28565b82815260059290921b84018101918181019086841115614e0457600080fd5b8286015b84811015614e2657614e1981614d9b565b8352918301918301614e08565b509695505050505050565b60008060408385031215614e4457600080fd5b823567ffffffffffffffff80821115614e5c57600080fd5b818501915085601f830112614e7057600080fd5b81356020614e80614de083614d77565b82815260059290921b84018101918181019089841115614e9f57600080fd5b948201945b83861015614ebd57853582529482019490820190614ea4565b96505086013592505080821115614ed357600080fd5b50614ee085828601614dbf565b9150509250929050565b600082601f830112614efb57600080fd5b813567ffffffffffffffff811115614f1557614f15614cf9565b614f4660207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601614d28565b818152846020838601011115614f5b57600080fd5b816020850160208301376000918101602001919091529392505050565b6bffffffffffffffffffffffff81168114612baf57600080fd5b60008060008060008060c08789031215614fab57600080fd5b86359550602087013567ffffffffffffffff80821115614fca57600080fd5b614fd68a838b01614eea565b96506040890135915080821115614fec57600080fd5b50614ff989828a01614eea565b945050606087013561500a81614f78565b9250608087013561501a81614f78565b915061502860a08801614d9b565b90509295509295509295565b803561ffff81168114614b5a57600080fd5b60008060008060008060a0878903121561505f57600080fd5b61506887614b42565b9550602087013567ffffffffffffffff81111561508457600080fd5b61509089828a01614c64565b90965094506150a3905060408801615034565b925060608701356150b381614b7a565b80925050608087013590509295509295509295565b600080604083850312156150db57600080fd5b6150e483614b42565b91506150f260208401614d9b565b90509250929050565b6000806040838503121561510e57600080fd5b61511783614d9b565b91506020830135614bb881614f78565b6000806020838503121561513a57600080fd5b823567ffffffffffffffff8082111561515257600080fd5b818501915085601f83011261516657600080fd5b81358181111561517557600080fd5b8660208260051b850101111561518a57600080fd5b60209290920196919550909350505050565b600080604083850312156151af57600080fd5b6151b883614d9b565b91506150f260208401614b42565b8015158114612baf57600080fd5b600080604083850312156151e757600080fd5b823591506020830135614bb8816151c6565b600081518084526020808501945080840160005b8381101561523f57815173ffffffffffffffffffffffffffffffffffffffff168752958201959082019060010161520d565b509495945050505050565b6bffffffffffffffffffffffff86811682528516602082015273ffffffffffffffffffffffffffffffffffffffff84811660408301528316606082015260a060808201819052600090611124908301846151f9565b600080600080606085870312156152b557600080fd5b6152be85614d9b565b935060208501359250604085013567ffffffffffffffff8111156152e157600080fd5b6152ed87828801614c64565b95989497509550505050565b60006020828403121561530b57600080fd5b5035919050565b60006020828403121561532457600080fd5b611b4282615034565b6000606082018583526020606081850152818651808452608086019150828801935060005b8181101561536e57845183529383019391830191600101615352565b5050848103604086015261538281876151f9565b98975050505050505050565b6000602082840312156153a057600080fd5b611b4282614d9b565b6000602082840312156153bb57600080fd5b5051919050565b8183823760009101908152919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8082018082111561124a5761124a6153d2565b600181811c9082168061542857607f821691505b602082108103615461577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f821115611b3157600081815260208120601f850160051c8101602086101561548e5750805b601f850160051c820191505b818110156154ad5782815560010161549a565b505050505050565b815167ffffffffffffffff8111156154cf576154cf614cf9565b6154e3816154dd8454615414565b84615467565b602080601f83116001811461553657600084156155005750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b1785556154ad565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b8281101561558357888601518255948401946001909101908401615564565b50858210156155bf57878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b84815283602082015260606040820152816060820152818360808301376000818301608090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01601019392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361568a5761568a6153d2565b5060010190565b6effffffffffffffffffffffffffffff8181168382160190808211156156b9576156b96153d2565b5092915050565b6bffffffffffffffffffffffff8181168382160280821691908281146156e8576156e86153d2565b505092915050565b6bffffffffffffffffffffffff8181168382160190808211156156b9576156b96153d2565b6bffffffffffffffffffffffff8616815273ffffffffffffffffffffffffffffffffffffffff8516602082015260ff8416604082015260a06060820152600061576160a0830185614bc3565b82810360808401526153828185614bc3565b6bffffffffffffffffffffffff8281168282160390808211156156b9576156b96153d2565b6000602082840312156157aa57600080fd5b8151611b42816151c6565b67ffffffffffffffff8181168382160190808211156156b9576156b96153d2565b8181038181111561124a5761124a6153d2565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b600067ffffffffffffffff808316818103615835576158356153d2565b6001019392505050565b83815282602082015260606040820152600061585e6060830184614bc3565b95945050505050565b8381526060602082015260006158806060830185614bc3565b82810360408401526158928185614bc3565b9695505050505050565b60208152600073ffffffffffffffffffffffffffffffffffffffff80845116602084015280602085015116604084015250604083015160e060608401526158e7610100840182614bc3565b905067ffffffffffffffff606085015116608084015261ffff60808501511660a084015260a084015160c084015263ffffffff60c08501511660e08401528091505092915050565b6000806000806080858703121561594557600080fd5b84519350602085015161595781614f78565b6040860151606090960151949790965092505050565b600073ffffffffffffffffffffffffffffffffffffffff8089168352808816602084015280871660408401525060c060608301526159ae60c0830186614bc3565b905061ffff8416608083015263ffffffff831660a0830152979650505050505050565b73ffffffffffffffffffffffffffffffffffffffff83168152604060208201526000612c9a6040830184614bc3565b60008060008060808587031215615a1657600080fd5b615a1f85614b42565b9350602085013567ffffffffffffffff811115615a3b57600080fd5b615a4787828801614eea565b935050615a5660408601615034565b91506060850135615a6681614b7a565b939692955090935050565b600080600060608486031215615a8657600080fd5b8351615a9181614f78565b809350506020808501517fffffffff0000000000000000000000000000000000000000000000000000000081168114615ac957600080fd5b604086015190935067ffffffffffffffff811115615ae657600080fd5b8501601f81018713615af757600080fd5b8051615b05614de082614d77565b81815260059190911b82018301908381019089831115615b2457600080fd5b928401925b82841015615b4b578351615b3c81614b7a565b82529284019290840190615b29565b80955050505050509250925092565b6000606082016bffffffffffffffffffffffff8616835260207fffffffff0000000000000000000000000000000000000000000000000000000086168185015260606040850152818551808452608086019150828701935060005b81811015615bd757845163ffffffff1683529383019391830191600101615bb5565b50909897505050505050505056fea164736f6c6343000813000a",
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

func (_FunctionsRouter *FunctionsRouterCaller) MAXCALLBACKRETURNBYTES(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _FunctionsRouter.contract.Call(opts, &out, "MAX_CALLBACK_RETURN_BYTES")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_FunctionsRouter *FunctionsRouterSession) MAXCALLBACKRETURNBYTES() (uint16, error) {
	return _FunctionsRouter.Contract.MAXCALLBACKRETURNBYTES(&_FunctionsRouter.CallOpts)
}

func (_FunctionsRouter *FunctionsRouterCallerSession) MAXCALLBACKRETURNBYTES() (uint16, error) {
	return _FunctionsRouter.Contract.MAXCALLBACKRETURNBYTES(&_FunctionsRouter.CallOpts)
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

func (_FunctionsRouter *FunctionsRouterCaller) GetAllowListId(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _FunctionsRouter.contract.Call(opts, &out, "getAllowListId")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_FunctionsRouter *FunctionsRouterSession) GetAllowListId() ([32]byte, error) {
	return _FunctionsRouter.Contract.GetAllowListId(&_FunctionsRouter.CallOpts)
}

func (_FunctionsRouter *FunctionsRouterCallerSession) GetAllowListId() ([32]byte, error) {
	return _FunctionsRouter.Contract.GetAllowListId(&_FunctionsRouter.CallOpts)
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

func (_FunctionsRouter *FunctionsRouterCaller) GetFlags(opts *bind.CallOpts, subscriptionId uint64) ([32]byte, error) {
	var out []interface{}
	err := _FunctionsRouter.contract.Call(opts, &out, "getFlags", subscriptionId)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_FunctionsRouter *FunctionsRouterSession) GetFlags(subscriptionId uint64) ([32]byte, error) {
	return _FunctionsRouter.Contract.GetFlags(&_FunctionsRouter.CallOpts, subscriptionId)
}

func (_FunctionsRouter *FunctionsRouterCallerSession) GetFlags(subscriptionId uint64) ([32]byte, error) {
	return _FunctionsRouter.Contract.GetFlags(&_FunctionsRouter.CallOpts, subscriptionId)
}

func (_FunctionsRouter *FunctionsRouterCaller) GetMaxConsumers(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _FunctionsRouter.contract.Call(opts, &out, "getMaxConsumers")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_FunctionsRouter *FunctionsRouterSession) GetMaxConsumers() (uint16, error) {
	return _FunctionsRouter.Contract.GetMaxConsumers(&_FunctionsRouter.CallOpts)
}

func (_FunctionsRouter *FunctionsRouterCallerSession) GetMaxConsumers() (uint16, error) {
	return _FunctionsRouter.Contract.GetMaxConsumers(&_FunctionsRouter.CallOpts)
}

func (_FunctionsRouter *FunctionsRouterCaller) GetProposedContractSet(opts *bind.CallOpts) (*big.Int, [][32]byte, []common.Address, error) {
	var out []interface{}
	err := _FunctionsRouter.contract.Call(opts, &out, "getProposedContractSet")

	if err != nil {
		return *new(*big.Int), *new([][32]byte), *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new([][32]byte)).(*[][32]byte)
	out2 := *abi.ConvertType(out[2], new([]common.Address)).(*[]common.Address)

	return out0, out1, out2, err

}

func (_FunctionsRouter *FunctionsRouterSession) GetProposedContractSet() (*big.Int, [][32]byte, []common.Address, error) {
	return _FunctionsRouter.Contract.GetProposedContractSet(&_FunctionsRouter.CallOpts)
}

func (_FunctionsRouter *FunctionsRouterCallerSession) GetProposedContractSet() (*big.Int, [][32]byte, []common.Address, error) {
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

func (_FunctionsRouter *FunctionsRouterCaller) IsValidCallbackGasLimit(opts *bind.CallOpts, subscriptionId uint64, callbackGasLimit uint32) error {
	var out []interface{}
	err := _FunctionsRouter.contract.Call(opts, &out, "isValidCallbackGasLimit", subscriptionId, callbackGasLimit)

	if err != nil {
		return err
	}

	return err

}

func (_FunctionsRouter *FunctionsRouterSession) IsValidCallbackGasLimit(subscriptionId uint64, callbackGasLimit uint32) error {
	return _FunctionsRouter.Contract.IsValidCallbackGasLimit(&_FunctionsRouter.CallOpts, subscriptionId, callbackGasLimit)
}

func (_FunctionsRouter *FunctionsRouterCallerSession) IsValidCallbackGasLimit(subscriptionId uint64, callbackGasLimit uint32) error {
	return _FunctionsRouter.Contract.IsValidCallbackGasLimit(&_FunctionsRouter.CallOpts, subscriptionId, callbackGasLimit)
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

func (_FunctionsRouter *FunctionsRouterTransactor) OwnerWithdraw(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FunctionsRouter.contract.Transact(opts, "ownerWithdraw", recipient, amount)
}

func (_FunctionsRouter *FunctionsRouterSession) OwnerWithdraw(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.OwnerWithdraw(&_FunctionsRouter.TransactOpts, recipient, amount)
}

func (_FunctionsRouter *FunctionsRouterTransactorSession) OwnerWithdraw(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.OwnerWithdraw(&_FunctionsRouter.TransactOpts, recipient, amount)
}

func (_FunctionsRouter *FunctionsRouterTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsRouter.contract.Transact(opts, "pause")
}

func (_FunctionsRouter *FunctionsRouterSession) Pause() (*types.Transaction, error) {
	return _FunctionsRouter.Contract.Pause(&_FunctionsRouter.TransactOpts)
}

func (_FunctionsRouter *FunctionsRouterTransactorSession) Pause() (*types.Transaction, error) {
	return _FunctionsRouter.Contract.Pause(&_FunctionsRouter.TransactOpts)
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

func (_FunctionsRouter *FunctionsRouterTransactor) ProposeContractsUpdate(opts *bind.TransactOpts, proposedContractSetIds [][32]byte, proposedContractSetAddresses []common.Address) (*types.Transaction, error) {
	return _FunctionsRouter.contract.Transact(opts, "proposeContractsUpdate", proposedContractSetIds, proposedContractSetAddresses)
}

func (_FunctionsRouter *FunctionsRouterSession) ProposeContractsUpdate(proposedContractSetIds [][32]byte, proposedContractSetAddresses []common.Address) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.ProposeContractsUpdate(&_FunctionsRouter.TransactOpts, proposedContractSetIds, proposedContractSetAddresses)
}

func (_FunctionsRouter *FunctionsRouterTransactorSession) ProposeContractsUpdate(proposedContractSetIds [][32]byte, proposedContractSetAddresses []common.Address) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.ProposeContractsUpdate(&_FunctionsRouter.TransactOpts, proposedContractSetIds, proposedContractSetAddresses)
}

func (_FunctionsRouter *FunctionsRouterTransactor) ProposeSubscriptionOwnerTransfer(opts *bind.TransactOpts, subscriptionId uint64, newOwner common.Address) (*types.Transaction, error) {
	return _FunctionsRouter.contract.Transact(opts, "proposeSubscriptionOwnerTransfer", subscriptionId, newOwner)
}

func (_FunctionsRouter *FunctionsRouterSession) ProposeSubscriptionOwnerTransfer(subscriptionId uint64, newOwner common.Address) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.ProposeSubscriptionOwnerTransfer(&_FunctionsRouter.TransactOpts, subscriptionId, newOwner)
}

func (_FunctionsRouter *FunctionsRouterTransactorSession) ProposeSubscriptionOwnerTransfer(subscriptionId uint64, newOwner common.Address) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.ProposeSubscriptionOwnerTransfer(&_FunctionsRouter.TransactOpts, subscriptionId, newOwner)
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

func (_FunctionsRouter *FunctionsRouterTransactor) SendRequest(opts *bind.TransactOpts, subscriptionId uint64, data []byte, dataVersion uint16, callbackGasLimit uint32, donId [32]byte) (*types.Transaction, error) {
	return _FunctionsRouter.contract.Transact(opts, "sendRequest", subscriptionId, data, dataVersion, callbackGasLimit, donId)
}

func (_FunctionsRouter *FunctionsRouterSession) SendRequest(subscriptionId uint64, data []byte, dataVersion uint16, callbackGasLimit uint32, donId [32]byte) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.SendRequest(&_FunctionsRouter.TransactOpts, subscriptionId, data, dataVersion, callbackGasLimit, donId)
}

func (_FunctionsRouter *FunctionsRouterTransactorSession) SendRequest(subscriptionId uint64, data []byte, dataVersion uint16, callbackGasLimit uint32, donId [32]byte) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.SendRequest(&_FunctionsRouter.TransactOpts, subscriptionId, data, dataVersion, callbackGasLimit, donId)
}

func (_FunctionsRouter *FunctionsRouterTransactor) SetFlags(opts *bind.TransactOpts, subscriptionId uint64, flags [32]byte) (*types.Transaction, error) {
	return _FunctionsRouter.contract.Transact(opts, "setFlags", subscriptionId, flags)
}

func (_FunctionsRouter *FunctionsRouterSession) SetFlags(subscriptionId uint64, flags [32]byte) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.SetFlags(&_FunctionsRouter.TransactOpts, subscriptionId, flags)
}

func (_FunctionsRouter *FunctionsRouterTransactorSession) SetFlags(subscriptionId uint64, flags [32]byte) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.SetFlags(&_FunctionsRouter.TransactOpts, subscriptionId, flags)
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

func (_FunctionsRouter *FunctionsRouterTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _FunctionsRouter.contract.Transact(opts, "transferOwnership", to)
}

func (_FunctionsRouter *FunctionsRouterSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.TransferOwnership(&_FunctionsRouter.TransactOpts, to)
}

func (_FunctionsRouter *FunctionsRouterTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.TransferOwnership(&_FunctionsRouter.TransactOpts, to)
}

func (_FunctionsRouter *FunctionsRouterTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsRouter.contract.Transact(opts, "unpause")
}

func (_FunctionsRouter *FunctionsRouterSession) Unpause() (*types.Transaction, error) {
	return _FunctionsRouter.Contract.Unpause(&_FunctionsRouter.TransactOpts)
}

func (_FunctionsRouter *FunctionsRouterTransactorSession) Unpause() (*types.Transaction, error) {
	return _FunctionsRouter.Contract.Unpause(&_FunctionsRouter.TransactOpts)
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

type FunctionsRouterConfigChangedIterator struct {
	Event *FunctionsRouterConfigChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRouterConfigChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRouterConfigChanged)
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
		it.Event = new(FunctionsRouterConfigChanged)
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

func (it *FunctionsRouterConfigChangedIterator) Error() error {
	return it.fail
}

func (it *FunctionsRouterConfigChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRouterConfigChanged struct {
	AdminFee                        *big.Int
	HandleOracleFulfillmentSelector [4]byte
	MaxCallbackGasLimits            []uint32
	Raw                             types.Log
}

func (_FunctionsRouter *FunctionsRouterFilterer) FilterConfigChanged(opts *bind.FilterOpts) (*FunctionsRouterConfigChangedIterator, error) {

	logs, sub, err := _FunctionsRouter.contract.FilterLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return &FunctionsRouterConfigChangedIterator{contract: _FunctionsRouter.contract, event: "ConfigChanged", logs: logs, sub: sub}, nil
}

func (_FunctionsRouter *FunctionsRouterFilterer) WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *FunctionsRouterConfigChanged) (event.Subscription, error) {

	logs, sub, err := _FunctionsRouter.contract.WatchLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRouterConfigChanged)
				if err := _FunctionsRouter.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
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

func (_FunctionsRouter *FunctionsRouterFilterer) ParseConfigChanged(log types.Log) (*FunctionsRouterConfigChanged, error) {
	event := new(FunctionsRouterConfigChanged)
	if err := _FunctionsRouter.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
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
	ReturnData     []byte
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

type FunctionsRouterRequestTimedOutIterator struct {
	Event *FunctionsRouterRequestTimedOut

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRouterRequestTimedOutIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRouterRequestTimedOut)
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
		it.Event = new(FunctionsRouterRequestTimedOut)
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

func (it *FunctionsRouterRequestTimedOutIterator) Error() error {
	return it.fail
}

func (it *FunctionsRouterRequestTimedOutIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRouterRequestTimedOut struct {
	RequestId [32]byte
	Raw       types.Log
}

func (_FunctionsRouter *FunctionsRouterFilterer) FilterRequestTimedOut(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsRouterRequestTimedOutIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsRouter.contract.FilterLogs(opts, "RequestTimedOut", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsRouterRequestTimedOutIterator{contract: _FunctionsRouter.contract, event: "RequestTimedOut", logs: logs, sub: sub}, nil
}

func (_FunctionsRouter *FunctionsRouterFilterer) WatchRequestTimedOut(opts *bind.WatchOpts, sink chan<- *FunctionsRouterRequestTimedOut, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsRouter.contract.WatchLogs(opts, "RequestTimedOut", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRouterRequestTimedOut)
				if err := _FunctionsRouter.contract.UnpackLog(event, "RequestTimedOut", log); err != nil {
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

func (_FunctionsRouter *FunctionsRouterFilterer) ParseRequestTimedOut(log types.Log) (*FunctionsRouterRequestTimedOut, error) {
	event := new(FunctionsRouterRequestTimedOut)
	if err := _FunctionsRouter.contract.UnpackLog(event, "RequestTimedOut", log); err != nil {
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
	case _FunctionsRouter.abi.Events["ConfigChanged"].ID:
		return _FunctionsRouter.ParseConfigChanged(log)
	case _FunctionsRouter.abi.Events["ConfigProposed"].ID:
		return _FunctionsRouter.ParseConfigProposed(log)
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
	case _FunctionsRouter.abi.Events["RequestTimedOut"].ID:
		return _FunctionsRouter.ParseRequestTimedOut(log)
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

func (FunctionsRouterConfigChanged) Topic() common.Hash {
	return common.HexToHash("0xe6a1eda76d42a6d1d813f26765716562044db1e8bd8be7d088705e64afd301ca")
}

func (FunctionsRouterConfigProposed) Topic() common.Hash {
	return common.HexToHash("0x0fcfd32a68209b42944376bc5f4bf72c41ba0f378cf60434f84390b82c9844cf")
}

func (FunctionsRouterConfigUpdated) Topic() common.Hash {
	return common.HexToHash("0xc07626096b49b0462a576c9a7878cf0675e784bb4832584c1289c11664e55f47")
}

func (FunctionsRouterContractProposed) Topic() common.Hash {
	return common.HexToHash("0x72a33d2f293a0a70fad221bb610d3d6b52aed2d840adae1fa721071fbd290cfd")
}

func (FunctionsRouterContractUpdated) Topic() common.Hash {
	return common.HexToHash("0xf8a6175bca1ba37d682089187edc5e20a859989727f10ca6bd9a5bc0de8caf94")
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
	return common.HexToHash("0x45bb48b6ec798595a260f114720360b95cc58c94c6ddd37a1acc3896ec94a23a")
}

func (FunctionsRouterRequestStart) Topic() common.Hash {
	return common.HexToHash("0x7c720ccd20069b8311a6be4ba1cf3294d09eb247aa5d73a8502054b6e68a2f54")
}

func (FunctionsRouterRequestTimedOut) Topic() common.Hash {
	return common.HexToHash("0xf1ca1e9147be737b04a2b018a79405f687a97de8dd8a2559bbe62357343af414")
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
	MAXCALLBACKRETURNBYTES(opts *bind.CallOpts) (uint16, error)

	GetAdminFee(opts *bind.CallOpts) (*big.Int, error)

	GetAllowListId(opts *bind.CallOpts) ([32]byte, error)

	GetConfigHash(opts *bind.CallOpts) ([32]byte, error)

	GetConsumer(opts *bind.CallOpts, client common.Address, subscriptionId uint64) (GetConsumer,

		error)

	GetContractById(opts *bind.CallOpts, id [32]byte, useProposed bool) (common.Address, error)

	GetContractById0(opts *bind.CallOpts, id [32]byte) (common.Address, error)

	GetFlags(opts *bind.CallOpts, subscriptionId uint64) ([32]byte, error)

	GetMaxConsumers(opts *bind.CallOpts) (uint16, error)

	GetProposedContractSet(opts *bind.CallOpts) (*big.Int, [][32]byte, []common.Address, error)

	GetSubscription(opts *bind.CallOpts, subscriptionId uint64) (GetSubscription,

		error)

	GetSubscriptionCount(opts *bind.CallOpts) (uint64, error)

	GetTotalBalance(opts *bind.CallOpts) (*big.Int, error)

	IsPaused(opts *bind.CallOpts) (bool, error)

	IsValidCallbackGasLimit(opts *bind.CallOpts, subscriptionId uint64, callbackGasLimit uint32) error

	Owner(opts *bind.CallOpts) (common.Address, error)

	Paused(opts *bind.CallOpts) (bool, error)

	PendingRequestExists(opts *bind.CallOpts, subscriptionId uint64) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AcceptSubscriptionOwnerTransfer(opts *bind.TransactOpts, subscriptionId uint64) (*types.Transaction, error)

	AddConsumer(opts *bind.TransactOpts, subscriptionId uint64, consumer common.Address) (*types.Transaction, error)

	CancelSubscription(opts *bind.TransactOpts, subscriptionId uint64, to common.Address) (*types.Transaction, error)

	CreateSubscription(opts *bind.TransactOpts) (*types.Transaction, error)

	Fulfill(opts *bind.TransactOpts, requestId [32]byte, response []byte, err []byte, juelsPerGas *big.Int, costWithoutFulfillment *big.Int, transmitter common.Address) (*types.Transaction, error)

	OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error)

	OracleWithdraw(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error)

	OwnerCancelSubscription(opts *bind.TransactOpts, subscriptionId uint64) (*types.Transaction, error)

	OwnerWithdraw(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error)

	Pause(opts *bind.TransactOpts) (*types.Transaction, error)

	ProposeConfigUpdate(opts *bind.TransactOpts, id [32]byte, config []byte) (*types.Transaction, error)

	ProposeContractsUpdate(opts *bind.TransactOpts, proposedContractSetIds [][32]byte, proposedContractSetAddresses []common.Address) (*types.Transaction, error)

	ProposeSubscriptionOwnerTransfer(opts *bind.TransactOpts, subscriptionId uint64, newOwner common.Address) (*types.Transaction, error)

	ProposeTimelockBlocks(opts *bind.TransactOpts, blocks uint16) (*types.Transaction, error)

	RecoverFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	RemoveConsumer(opts *bind.TransactOpts, subscriptionId uint64, consumer common.Address) (*types.Transaction, error)

	SendRequest(opts *bind.TransactOpts, subscriptionId uint64, data []byte, dataVersion uint16, callbackGasLimit uint32, donId [32]byte) (*types.Transaction, error)

	SetFlags(opts *bind.TransactOpts, subscriptionId uint64, flags [32]byte) (*types.Transaction, error)

	TimeoutRequests(opts *bind.TransactOpts, requestIdsToTimeout [][32]byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Unpause(opts *bind.TransactOpts) (*types.Transaction, error)

	UpdateConfig(opts *bind.TransactOpts, id [32]byte) (*types.Transaction, error)

	UpdateContracts(opts *bind.TransactOpts) (*types.Transaction, error)

	UpdateTimelockBlocks(opts *bind.TransactOpts) (*types.Transaction, error)

	ValidateProposedContracts(opts *bind.TransactOpts, id [32]byte, data []byte) (*types.Transaction, error)

	FilterConfigChanged(opts *bind.FilterOpts) (*FunctionsRouterConfigChangedIterator, error)

	WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *FunctionsRouterConfigChanged) (event.Subscription, error)

	ParseConfigChanged(log types.Log) (*FunctionsRouterConfigChanged, error)

	FilterConfigProposed(opts *bind.FilterOpts) (*FunctionsRouterConfigProposedIterator, error)

	WatchConfigProposed(opts *bind.WatchOpts, sink chan<- *FunctionsRouterConfigProposed) (event.Subscription, error)

	ParseConfigProposed(log types.Log) (*FunctionsRouterConfigProposed, error)

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

	FilterRequestTimedOut(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsRouterRequestTimedOutIterator, error)

	WatchRequestTimedOut(opts *bind.WatchOpts, sink chan<- *FunctionsRouterRequestTimedOut, requestId [][32]byte) (event.Subscription, error)

	ParseRequestTimedOut(log types.Log) (*FunctionsRouterRequestTimedOut, error)

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
