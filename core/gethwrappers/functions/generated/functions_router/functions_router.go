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

type FunctionsRouterConfig struct {
	MaxConsumersPerSubscription          uint16
	AdminFee                             *big.Int
	HandleOracleFulfillmentSelector      [4]byte
	GasForCallExactCheck                 uint16
	MaxCallbackGasLimits                 []uint32
	SubscriptionDepositCompletedRequests uint16
	SubscriptionDepositJuels             *big.Int
}

type IFunctionsSubscriptionsConsumer struct {
	Allowed           bool
	InitiatedRequests uint64
	CompletedRequests uint64
}

type IFunctionsSubscriptionsSubscription struct {
	Balance        *big.Int
	Owner          common.Address
	BlockedBalance *big.Int
	ProposedOwner  common.Address
	Consumers      []common.Address
	Flags          [32]byte
}

var FunctionsRouterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"linkToken\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint16\",\"name\":\"maxConsumersPerSubscription\",\"type\":\"uint16\"},{\"internalType\":\"uint72\",\"name\":\"adminFee\",\"type\":\"uint72\"},{\"internalType\":\"bytes4\",\"name\":\"handleOracleFulfillmentSelector\",\"type\":\"bytes4\"},{\"internalType\":\"uint16\",\"name\":\"gasForCallExactCheck\",\"type\":\"uint16\"},{\"internalType\":\"uint32[]\",\"name\":\"maxCallbackGasLimits\",\"type\":\"uint32[]\"}],\"internalType\":\"structFunctionsRouter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"CannotRemoveWithPendingRequests\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"DuplicateRequestId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptyRequestData\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"limit\",\"type\":\"uint32\"}],\"name\":\"GasLimitTooBig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"IdentifierIsReserved\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"currentBalanceJuels\",\"type\":\"uint96\"}],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidConsumer\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"value\",\"type\":\"uint8\"}],\"name\":\"InvalidGasFlagValue\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidProposal\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSubscription\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposedOwner\",\"type\":\"address\"}],\"name\":\"MustBeProposedOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeSubscriptionOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableFromCoordinator\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableFromLink\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"RouteNotFound\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"SenderMustAcceptTermsOfService\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TimeoutNotExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"maximumConsumers\",\"type\":\"uint16\"}],\"name\":\"TooManyConsumers\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"totalBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deductionAttempt\",\"type\":\"uint256\"}],\"name\":\"TotalBalanceInvariantViolated\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint16\",\"name\":\"maxConsumersPerSubscription\",\"type\":\"uint16\"},{\"internalType\":\"uint72\",\"name\":\"adminFee\",\"type\":\"uint72\"},{\"internalType\":\"bytes4\",\"name\":\"handleOracleFulfillmentSelector\",\"type\":\"bytes4\"},{\"internalType\":\"uint16\",\"name\":\"gasForCallExactCheck\",\"type\":\"uint16\"},{\"internalType\":\"uint32[]\",\"name\":\"maxCallbackGasLimits\",\"type\":\"uint32[]\"}],\"indexed\":false,\"internalType\":\"structFunctionsRouter.Config\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"ConfigUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"proposedContractSetId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"proposedContractSetFromAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"proposedContractSetToAddress\",\"type\":\"address\"}],\"name\":\"ContractProposed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"ContractUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"enumFunctionsResponse.FulfillResult\",\"name\":\"resultCode\",\"type\":\"uint8\"}],\"name\":\"RequestNotProcessed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalCostJuels\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"enumFunctionsResponse.FulfillResult\",\"name\":\"resultCode\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"response\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"callbackReturnData\",\"type\":\"bytes\"}],\"name\":\"RequestProcessed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"donId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"subscriptionOwner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requestingContract\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requestInitiator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"dataVersion\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"estimatedTotalCostJuels\",\"type\":\"uint96\"}],\"name\":\"RequestStart\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"RequestTimedOut\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"fundsRecipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fundsAmount\",\"type\":\"uint256\"}],\"name\":\"SubscriptionCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"SubscriptionCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"MAX_CALLBACK_RETURN_BYTES\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"acceptSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"addConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"cancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createSubscription\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"createSubscriptionWithConsumer\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"response\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"juelsPerGas\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"costWithoutCallback\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"estimatedTotalCostJuels\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint72\",\"name\":\"adminFee\",\"type\":\"uint72\"},{\"internalType\":\"uint72\",\"name\":\"donFee\",\"type\":\"uint72\"},{\"internalType\":\"uint40\",\"name\":\"gasOverheadBeforeCallback\",\"type\":\"uint40\"},{\"internalType\":\"uint40\",\"name\":\"gasOverheadAfterCallback\",\"type\":\"uint40\"},{\"internalType\":\"uint32\",\"name\":\"timeoutTimestamp\",\"type\":\"uint32\"}],\"internalType\":\"structFunctionsResponse.Commitment\",\"name\":\"commitment\",\"type\":\"tuple\"}],\"name\":\"fulfill\",\"outputs\":[{\"internalType\":\"enumFunctionsResponse.FulfillResult\",\"name\":\"resultCode\",\"type\":\"uint8\"},{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAdminFee\",\"outputs\":[{\"internalType\":\"uint72\",\"name\":\"\",\"type\":\"uint72\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowListId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint16\",\"name\":\"maxConsumersPerSubscription\",\"type\":\"uint16\"},{\"internalType\":\"uint72\",\"name\":\"adminFee\",\"type\":\"uint72\"},{\"internalType\":\"bytes4\",\"name\":\"handleOracleFulfillmentSelector\",\"type\":\"bytes4\"},{\"internalType\":\"uint16\",\"name\":\"gasForCallExactCheck\",\"type\":\"uint16\"},{\"internalType\":\"uint32[]\",\"name\":\"maxCallbackGasLimits\",\"type\":\"uint32[]\"}],\"internalType\":\"structFunctionsRouter.Config\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"getConsumer\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"initiatedRequests\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"completedRequests\",\"type\":\"uint64\"}],\"internalType\":\"structIFunctionsSubscriptions.Consumer\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"getContractById\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"getFlags\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"getProposedContractById\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getProposedContractSet\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"},{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"getSubscription\",\"outputs\":[{\"components\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"blockedBalance\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"proposedOwner\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"},{\"internalType\":\"bytes32\",\"name\":\"flags\",\"type\":\"bytes32\"}],\"internalType\":\"structIFunctionsSubscriptions.Subscription\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSubscriptionCount\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTotalBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"}],\"name\":\"isValidCallbackGasLimit\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"oracleWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"ownerCancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"ownerWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"pendingRequestExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"proposedContractSetIds\",\"type\":\"bytes32[]\"},{\"internalType\":\"address[]\",\"name\":\"proposedContractSetAddresses\",\"type\":\"address[]\"}],\"name\":\"proposeContractsUpdate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"proposeSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"removeConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint16\",\"name\":\"dataVersion\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"donId\",\"type\":\"bytes32\"}],\"name\":\"sendRequest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint16\",\"name\":\"dataVersion\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"donId\",\"type\":\"bytes32\"}],\"name\":\"sendRequestToProposed\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"allowListId\",\"type\":\"bytes32\"}],\"name\":\"setAllowListId\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"flags\",\"type\":\"bytes32\"}],\"name\":\"setFlags\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"estimatedTotalCostJuels\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint72\",\"name\":\"adminFee\",\"type\":\"uint72\"},{\"internalType\":\"uint72\",\"name\":\"donFee\",\"type\":\"uint72\"},{\"internalType\":\"uint40\",\"name\":\"gasOverheadBeforeCallback\",\"type\":\"uint40\"},{\"internalType\":\"uint40\",\"name\":\"gasOverheadAfterCallback\",\"type\":\"uint40\"},{\"internalType\":\"uint32\",\"name\":\"timeoutTimestamp\",\"type\":\"uint32\"}],\"internalType\":\"structFunctionsResponse.Commitment[]\",\"name\":\"requestsToTimeoutByCommitment\",\"type\":\"tuple[]\"}],\"name\":\"timeoutRequests\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint16\",\"name\":\"maxConsumersPerSubscription\",\"type\":\"uint16\"},{\"internalType\":\"uint72\",\"name\":\"adminFee\",\"type\":\"uint72\"},{\"internalType\":\"bytes4\",\"name\":\"handleOracleFulfillmentSelector\",\"type\":\"bytes4\"},{\"internalType\":\"uint16\",\"name\":\"gasForCallExactCheck\",\"type\":\"uint16\"},{\"internalType\":\"uint32[]\",\"name\":\"maxCallbackGasLimits\",\"type\":\"uint32[]\"}],\"internalType\":\"structFunctionsRouter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"updateConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"updateContracts\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b50604051620062d1380380620062d18339810160408190526200003491620004df565b6001600160a01b0382166080526006805460ff191690553380600081620000a25760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600680546001600160a01b0380851661010002610100600160a81b031990921691909117909155811615620000dc57620000dc81620000f8565b505050620000f081620001aa60201b60201c565b50506200067b565b336001600160a01b03821603620001525760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000099565b600780546001600160a01b0319166001600160a01b03838116918217909255600654604051919261010090910416907fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127890600090a350565b620001b462000287565b8051600a80546020808501516040860151606087015161ffff908116600160781b0261ffff60781b1960e09390931c6b010000000000000000000000029290921665ffffffffffff60581b196001600160481b0390941662010000026001600160581b031990961691909716179390931716939093171781556080830151805184936200024792600b92910190620002ea565b509050507f049ce2e6e1420eb4b07b425e90129186833eb346bda40b37d5d921aad482f71c816040516200027c9190620005dd565b60405180910390a150565b60065461010090046001600160a01b03163314620002e85760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640162000099565b565b828054828255906000526020600020906007016008900481019282156200038e5791602002820160005b838211156200035a57835183826101000a81548163ffffffff021916908363ffffffff160217905550926020019260040160208160030104928301926001030262000314565b80156200038c5782816101000a81549063ffffffff02191690556004016020816003010492830192600103026200035a565b505b506200039c929150620003a0565b5090565b5b808211156200039c5760008155600101620003a1565b634e487b7160e01b600052604160045260246000fd5b60405160a081016001600160401b0381118282101715620003f257620003f2620003b7565b60405290565b604051601f8201601f191681016001600160401b0381118282101715620004235762000423620003b7565b604052919050565b805161ffff811681146200043e57600080fd5b919050565b600082601f8301126200045557600080fd5b815160206001600160401b03821115620004735762000473620003b7565b8160051b62000484828201620003f8565b92835284810182019282810190878511156200049f57600080fd5b83870192505b84831015620004d457825163ffffffff81168114620004c45760008081fd5b82529183019190830190620004a5565b979650505050505050565b60008060408385031215620004f357600080fd5b82516001600160a01b03811681146200050b57600080fd5b60208401519092506001600160401b03808211156200052957600080fd5b9084019060a082870312156200053e57600080fd5b62000548620003cd565b62000553836200042b565b815260208301516001600160481b03811681146200057057600080fd5b602082015260408301516001600160e01b0319811681146200059157600080fd5b6040820152620005a4606084016200042b565b6060820152608083015182811115620005bc57600080fd5b620005ca8882860162000443565b6080830152508093505050509250929050565b6020808252825161ffff90811683830152838201516001600160481b03166040808501919091528401516001600160e01b0319166060808501919091528401511660808084019190915283015160a080840152805160c0840181905260009291820190839060e08601905b808310156200067057835163ffffffff16825292840192600192909201919084019062000648565b509695505050505050565b608051615c1e620006b3600039600081816113b40152818161218601528181612a7e01528181612b4201526132bd0152615c1e6000f3fe608060405234801561001057600080fd5b50600436106102de5760003560e01c80637341c10c11610186578063aab396bd116100e3578063d7ae1d3011610097578063e82ad7d411610071578063e82ad7d4146106bf578063ea320e0b146106d2578063f2fde38b146106e557600080fd5b8063d7ae1d3014610686578063e72f6e3014610699578063e82622aa146106ac57600080fd5b8063badc3eb6116100c8578063badc3eb614610648578063c3f909d41461065e578063cc77470a1461067357600080fd5b8063aab396bd14610638578063b734c0f41461064057600080fd5b80639f87fad71161013a578063a47c76961161011f578063a47c7696146105f2578063a4c0ed3614610612578063a9c9a9181461062557600080fd5b80639f87fad7146105d7578063a21a23e4146105ea57600080fd5b8063823597401161016b57806382359740146105995780638456cb59146105ac5780638da5cb5b146105b457600080fd5b80637341c10c1461057e57806379ba50971461059157600080fd5b80633f4ba83a1161023f5780635c975abb116101f357806366419970116101cd57806366419970146104d6578063674603d0146104fd5780636a2215de1461054657600080fd5b80635c975abb146104995780635ed6dfba146104b057806366316d8d146104c357600080fd5b8063461d276211610224578063461d2762146104455780634b8832d31461045857806355fedefa1461046b57600080fd5b80633f4ba83a1461041c57806341db4ca31461042457600080fd5b80631ded3b36116102965780632a905ccc1161027b5780632a905ccc146103ba57806333060529146103e85780633e871e4d1461040957600080fd5b80631ded3b361461039457806321b60e7f146103a757600080fd5b806310fc49c1116102c757806310fc49c11461031857806312b583491461032b578063181f5a771461034b57600080fd5b806302bcc5b6146102e35780630c5d49cb146102f8575b600080fd5b6102f66102f13660046148f7565b6106f8565b005b610300608481565b60405161ffff90911681526020015b60405180910390f35b6102f6610326366004614938565b610757565b6000546040516bffffffffffffffffffffffff909116815260200161030f565b6103876040518060400160405280601781526020017f46756e6374696f6e7320526f757465722076312e302e3000000000000000000081525081565b60405161030f91906149df565b6102f66103a23660046149f2565b610853565b6102f66103b5366004614bb5565b610885565b600a5462010000900468ffffffffffffffffff1660405168ffffffffffffffffff909116815260200161030f565b6103fb6103f6366004614e61565b6109b2565b60405161030f929190614f49565b6102f6610417366004614fd6565b610d8a565b6102f661100f565b6104376104323660046150d8565b611021565b60405190815260200161030f565b6104376104533660046150d8565b611081565b6102f661046636600461515c565b61108d565b6104376104793660046148f7565b67ffffffffffffffff166000908152600360208190526040909120015490565b60065460ff165b604051901515815260200161030f565b6102f66104be36600461518a565b6111db565b6102f66104d136600461518a565b6113fd565b60025467ffffffffffffffff165b60405167ffffffffffffffff909116815260200161030f565b61051061050b3660046151b8565b61154f565b6040805182511515815260208084015167ffffffffffffffff90811691830191909152928201519092169082015260600161030f565b6105596105543660046151e6565b6115df565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161030f565b6102f661058c36600461515c565b61169e565b6102f6611851565b6102f66105a73660046148f7565b611978565b6102f6611abf565b600654610100900473ffffffffffffffffffffffffffffffffffffffff16610559565b6102f66105e536600461515c565b611acf565b6104e4611ea4565b6106056106003660046148f7565b612031565b60405161030f9190615250565b6102f66106203660046152d8565b612166565b6105596106333660046151e6565b6123b2565b600954610437565b6102f6612411565b61065061255d565b60405161030f929190615334565b61066661262d565b60405161030f919061538b565b6104e4610681366004615435565b612762565b6102f661069436600461515c565b6129e2565b6102f66106a7366004615435565b612a45565b6102f66106ba366004615452565b612bbe565b6104a06106cd3660046148f7565b612e8f565b6102f66106e03660046151e6565b612fde565b6102f66106f3366004615435565b612feb565b610700612ffc565b61070981613004565b67ffffffffffffffff81166000908152600360205260409020546107549082906c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1661307a565b50565b67ffffffffffffffff8216600090815260036020819052604082200154600b54911a9081106107bc576040517f45c108ce00000000000000000000000000000000000000000000000000000000815260ff821660048201526024015b60405180910390fd5b6000600a6001018260ff16815481106107d7576107d76154c8565b90600052602060002090600891828204019190066004029054906101000a900463ffffffff1690508063ffffffff168363ffffffff16111561084d576040517f1d70f87a00000000000000000000000000000000000000000000000000000000815263ffffffff821660048201526024016107b3565b50505050565b61085b612ffc565b61086482613004565b67ffffffffffffffff90911660009081526003602081905260409091200155565b61088d613366565b8051600a80546020808501516040860151606087015161ffff9081166f01000000000000000000000000000000027fffffffffffffffffffffffffffffff0000ffffffffffffffffffffffffffffff60e09390931c6b01000000000000000000000002929092167fffffffffffffffffffffffffffffff000000000000ffffffffffffffffffffff68ffffffffffffffffff90941662010000027fffffffffffffffffffffffffffffffffffffffffff0000000000000000000000909616919097161793909317169390931717815560808301518051849361097492600b92910190614742565b509050507f049ce2e6e1420eb4b07b425e90129186833eb346bda40b37d5d921aad482f71c816040516109a7919061538b565b60405180910390a150565b6000806109bd6133ec565b826020015173ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610a26576040517f8bec23e700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8251600090815260056020526040902054610a885782516020840151604051600294507f1a90e9a50793db2e394cf581e7c522e10c358a81e70acf6b5a0edd620c08dee191610a7891889087906154f7565b60405180910390a2506000610d7f565b8251600090815260056020908152604091829020549151610aab91869101615529565b6040516020818303038152906040528051906020012014610b035782516020840151604051600694507f1a90e9a50793db2e394cf581e7c522e10c358a81e70acf6b5a0edd620c08dee191610a7891889087906154f7565b8261012001518360a0015163ffffffff16610b1e9190615685565b64ffffffffff165a1015610b695782516020840151604051600494507f1a90e9a50793db2e394cf581e7c522e10c358a81e70acf6b5a0edd620c08dee191610a7891889087906154f7565b6000610b7e8460a0015163ffffffff166133f4565b610b8890886156aa565b9050600081878660c0015168ffffffffffffffffff16610ba891906156d2565b610bb291906156d2565b9050610bc18560800151612031565b600001516bffffffffffffffffffffffff16816bffffffffffffffffffffffff161115610c395784516020860151604051600596507f1a90e9a50793db2e394cf581e7c522e10c358a81e70acf6b5a0edd620c08dee191610c25918a9089906154f7565b60405180910390a25060009150610d7f9050565b84604001516bffffffffffffffffffffffff16816bffffffffffffffffffffffff161115610c9e5784516020860151604051600396507f1a90e9a50793db2e394cf581e7c522e10c358a81e70acf6b5a0edd620c08dee191610c25918a9089906154f7565b505082516000908152600560205260408120819055835160a08501516060860151610cce92918c918c9190613496565b8051909150610cde576001610ce1565b60005b92506000610d1b8560800151866040015187606001518860c0015168ffffffffffffffffff168c610d1588602001516133f4565b8d613654565b9050846080015167ffffffffffffffff1685600001517f64778f26c70b60a8d7e29e2451b3844302d959448401c0535b768ed88c6b505e836020015189888f8f8960400151604051610d72969594939291906156f7565b60405180910390a3519150505b965096945050505050565b610d92613366565b8151815181141580610da45750600881115b15610ddb576040517fee03280800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b81811015610ecd576000848281518110610dfa57610dfa6154c8565b602002602001015190506000848381518110610e1857610e186154c8565b60200260200101519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161480610e83575060008281526008602052604090205473ffffffffffffffffffffffffffffffffffffffff8281169116145b15610eba576040517fee03280800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b505080610ec69061577a565b9050610dde565b506040805180820190915283815260208082018490528451600c91610ef69183918801906147ed565b506020828101518051610f0f9260018501920190614828565b5090505060005b835181101561084d577f8b052f0f4bf82fede7daffea71592b29d5ef86af1f3c7daaa0345dbb2f52f481848281518110610f5257610f526154c8565b602002602001015160086000878581518110610f7057610f706154c8565b6020026020010151815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16858481518110610fb957610fb96154c8565b6020026020010151604051610ff79392919092835273ffffffffffffffffffffffffffffffffffffffff918216602084015216604082015260600190565b60405180910390a16110088161577a565b9050610f16565b611017613366565b61101f6139d6565b565b60008061102d836115df565b905061107583828a8a8a8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508c92508b9150613a539050565b98975050505050505050565b60008061102d836123b2565b6110956133ec565b61109e82613e4b565b6110a6613f11565b73ffffffffffffffffffffffffffffffffffffffff8116158061110d575067ffffffffffffffff821660009081526003602052604090206001015473ffffffffffffffffffffffffffffffffffffffff8281166c0100000000000000000000000090920416145b15611144576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff821660008181526003602090815260409182902060010180546bffffffffffffffffffffffff166c0100000000000000000000000073ffffffffffffffffffffffffffffffffffffffff8716908102919091179091558251338152918201527f69436ea6df009049404f564eff6622cd00522b0bd6a89efd9e52a355c4a879be910160405180910390a25050565b6111e3612ffc565b806bffffffffffffffffffffffff166000036112195750306000908152600160205260409020546bffffffffffffffffffffffff165b306000908152600160205260409020546bffffffffffffffffffffffff908116908216811015611285576040517f6b0fe56f0000000000000000000000000000000000000000000000000000000081526bffffffffffffffffffffffff821660048201526024016107b3565b6000546bffffffffffffffffffffffff808416911610156112ee576000546040517fdda2b2160000000000000000000000000000000000000000000000000000000081526bffffffffffffffffffffffff918216600482015290831660248201526044016107b3565b306000908152600160205260408120805484929061131b9084906bffffffffffffffffffffffff166157b2565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550816000808282829054906101000a90046bffffffffffffffffffffffff1661137191906157b2565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055506113f883836bffffffffffffffffffffffff167f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1661401b9092919063ffffffff16565b505050565b6114056133ec565b806bffffffffffffffffffffffff1660000361144d576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b336000908152600160205260409020546bffffffffffffffffffffffff9081169082168110156114b9576040517f6b0fe56f0000000000000000000000000000000000000000000000000000000081526bffffffffffffffffffffffff821660048201526024016107b3565b6000546bffffffffffffffffffffffff80841691161015611522576000546040517fdda2b2160000000000000000000000000000000000000000000000000000000081526bffffffffffffffffffffffff918216600482015290831660248201526044016107b3565b336000908152600160205260408120805484929061131b9084906bffffffffffffffffffffffff166157b2565b60408051606080820183526000808352602080840182905292840181905273ffffffffffffffffffffffffffffffffffffffff861681526004835283812067ffffffffffffffff868116835290845290849020845192830185525460ff81161515835261010081048216938301939093526901000000000000000000909204909116918101919091525b92915050565b6000805b600c5460ff8216101561166857600c805460ff8316908110611607576116076154c8565b9060005260206000200154830361165857600d805460ff831690811061162f5761162f6154c8565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff169392505050565b611661816157d7565b90506115e3565b506040517f80833e33000000000000000000000000000000000000000000000000000000008152600481018390526024016107b3565b6116a66133ec565b6116af82613e4b565b6116b7613f11565b60006116c6600a5461ffff1690565b67ffffffffffffffff841660009081526003602052604090206002015490915061ffff821611611728576040517fb72bc70300000000000000000000000000000000000000000000000000000000815261ffff821660048201526024016107b3565b73ffffffffffffffffffffffffffffffffffffffff8216600090815260046020908152604080832067ffffffffffffffff8716845290915290205460ff161561177057505050565b73ffffffffffffffffffffffffffffffffffffffff8216600081815260046020908152604080832067ffffffffffffffff881680855290835281842080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001908117909155600384528285206002018054918201815585529383902090930180547fffffffffffffffffffffffff000000000000000000000000000000000000000016851790555192835290917f43dc749a04ac8fb825cbd514f7c0e13f13bc6f2ee66043b76629d51776cff8e091015b60405180910390a2505050565b60075473ffffffffffffffffffffffffffffffffffffffff1633146118d2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016107b3565b600680547fffffffffffffffffffffff0000000000000000000000000000000000000000ff81166101003381810292909217909355600780547fffffffffffffffffffffffff00000000000000000000000000000000000000001690556040519290910473ffffffffffffffffffffffffffffffffffffffff169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a350565b6119806133ec565b611988613f11565b67ffffffffffffffff81166000908152600360205260409020805460019091015473ffffffffffffffffffffffffffffffffffffffff6c010000000000000000000000009283900481169290910416338114611a28576040517f4e1d9f1800000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821660048201526024016107b3565b67ffffffffffffffff831660008181526003602090815260409182902080546c01000000000000000000000000339081026bffffffffffffffffffffffff928316178355600190920180549091169055825173ffffffffffffffffffffffffffffffffffffffff87168152918201527f6f1dc65165ffffedfd8e507b4a0f1fcfdada045ed11f6c26ba27cedfe87802f09101611844565b611ac7613366565b61101f6140a8565b611ad76133ec565b611ae082613e4b565b611ae8613f11565b73ffffffffffffffffffffffffffffffffffffffff8116600090815260046020908152604080832067ffffffffffffffff8087168552908352928190208151606081018352905460ff8116151580835261010082048616948301949094526901000000000000000000900490931690830152611b90576040517f71e8313700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b806040015167ffffffffffffffff16816020015167ffffffffffffffff1614611be5576040517f06eb10c800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff8316600090815260036020908152604080832060020180548251818502810185019093528083529192909190830182828015611c6057602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311611c35575b5050505050905060005b8151811015611e08578373ffffffffffffffffffffffffffffffffffffffff16828281518110611c9c57611c9c6154c8565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1603611df8578160018351611cce91906157f6565b81518110611cde57611cde6154c8565b6020026020010151600360008767ffffffffffffffff1667ffffffffffffffff1681526020019081526020016000206002018281548110611d2157611d216154c8565b600091825260208083209190910180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff949094169390931790925567ffffffffffffffff87168152600390915260409020600201805480611d9b57611d9b615809565b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff0000000000000000000000000000000000000000169055019055611e08565b611e018161577a565b9050611c6a565b5073ffffffffffffffffffffffffffffffffffffffff8316600081815260046020908152604080832067ffffffffffffffff89168085529083529281902080547fffffffffffffffffffffffffffffff00000000000000000000000000000000001690555192835290917f182bff9831466789164ca77075fffd84916d35a8180ba73c27e45634549b445b91015b60405180910390a250505050565b6000611eae6133ec565b611eb6613f11565b60028054600090611ed09067ffffffffffffffff16615838565b825467ffffffffffffffff8083166101009490940a93840293021916919091179091556040805160c0810182526000808252336020830152918101829052606081018290529192506080820190604051908082528060200260200182016040528015611f46578160200160208202803683370190505b5081526000602091820181905267ffffffffffffffff841681526003825260409081902083518484015173ffffffffffffffffffffffffffffffffffffffff9081166c010000000000000000000000009081026bffffffffffffffffffffffff9384161784559386015160608701519091169093029216919091176001820155608083015180519192611fe192600285019290910190614828565b5060a0919091015160039091015560405133815267ffffffffffffffff8216907f464722b4166576d3dcbba877b999bc35cf911f4eaf434b7eba68fa113951d0bf9060200160405180910390a290565b6040805160c0810182526000808252602082018190529181018290526060808201839052608082015260a081019190915261206b82613004565b67ffffffffffffffff8216600090815260036020908152604091829020825160c08101845281546bffffffffffffffffffffffff808216835273ffffffffffffffffffffffffffffffffffffffff6c0100000000000000000000000092839004811684870152600185015491821684880152919004166060820152600282018054855181860281018601909652808652919492936080860193929083018282801561214c57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311612121575b505050505081526020016003820154815250509050919050565b61216e6133ec565b3373ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016146121dd576040517f44b0e3c300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60208114612217576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000612225828401846148f7565b67ffffffffffffffff81166000908152600360205260409020549091506c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1661229e576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff8116600090815260036020526040812080546bffffffffffffffffffffffff16918691906122d583856156d2565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550846000808282829054906101000a90046bffffffffffffffffffffffff1661232b91906156d2565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055508167ffffffffffffffff167fd39ec07f4e209f627a4c427971473820dc129761ba28de8906bd56f57101d4f8828784612392919061585f565b6040805192835260208301919091520160405180910390a2505050505050565b60008181526008602052604081205473ffffffffffffffffffffffffffffffffffffffff16806115d9576040517f80833e33000000000000000000000000000000000000000000000000000000008152600481018490526024016107b3565b612419613366565b60005b600c5481101561253c576000600c600001828154811061243e5761243e6154c8565b906000526020600020015490506000600c6001018381548110612463576124636154c8565b6000918252602080832091909101548483526008825260409283902054835186815273ffffffffffffffffffffffffffffffffffffffff91821693810193909352169181018290529091507ff8a6175bca1ba37d682089187edc5e20a859989727f10ca6bd9a5bc0de8caf949060600160405180910390a160009182526008602052604090912080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff9092169190911790556125358161577a565b905061241c565b50600c600061254b82826148a2565b6125596001830160006148a2565b5050565b606080600c600001600c600101818054806020026020016040519081016040528092919081815260200182805480156125b557602002820191906000526020600020905b8154815260200190600101908083116125a1575b505050505091508080548060200260200160405190810160405280929190818152602001828054801561261e57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116125f3575b50505050509050915091509091565b6040805160a08101825260008082526020820181905291810182905260608082019290925260808101919091526040805160a081018252600a805461ffff808216845268ffffffffffffffffff620100008304166020808601919091527fffffffff000000000000000000000000000000000000000000000000000000006b010000000000000000000000840460e01b16858701526f01000000000000000000000000000000909204166060840152600b80548551818402810184019096528086529394929360808601939283018282801561275457602002820191906000526020600020906000905b82829054906101000a900463ffffffff1663ffffffff16815260200190600401906020826003010492830192600103820291508084116127175790505b505050505081525050905090565b600061276c6133ec565b612774613f11565b6002805460009061278e9067ffffffffffffffff16615838565b825467ffffffffffffffff8083166101009490940a93840293021916919091179091556040805160c0810182526000808252336020830152918101829052606081018290529192506080820190604051908082528060200260200182016040528015612804578160200160208202803683370190505b5081526000602091820181905267ffffffffffffffff841681526003825260409081902083518484015173ffffffffffffffffffffffffffffffffffffffff9081166c010000000000000000000000009081026bffffffffffffffffffffffff938416178455938601516060870151909116909302921691909117600182015560808301518051919261289f92600285019290910190614828565b5060a0919091015160039182015567ffffffffffffffff82166000818152602092835260408082206002018054600180820183559184528584200180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff891690811790915583526004855281832084845285529181902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169092179091555133815290917f464722b4166576d3dcbba877b999bc35cf911f4eaf434b7eba68fa113951d0bf910160405180910390a260405173ffffffffffffffffffffffffffffffffffffffff8316815267ffffffffffffffff8216907f43dc749a04ac8fb825cbd514f7c0e13f13bc6f2ee66043b76629d51776cff8e09060200160405180910390a2919050565b6129ea6133ec565b6129f382613e4b565b6129fb613f11565b612a0482612e8f565b15612a3b576040517f06eb10c800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b612559828261307a565b612a4d612ffc565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526000907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa158015612ada573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612afe9190615872565b6000549091506bffffffffffffffffffffffff16818110156113f8576000612b2682846157f6565b9050612b6973ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016858361401b565b6040805173ffffffffffffffffffffffffffffffffffffffff86168152602081018390527f59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b436600910160405180910390a150505050565b612bc66133ec565b60005b818110156113f8576000838383818110612be557612be56154c8565b90506101600201803603810190612bfc919061588b565b80516080820151600082815260056020908152604091829020549151949550929391929091612c2d91869101615529565b6040516020818303038152906040528051906020012014612c7a576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b82610140015163ffffffff16421015612cbf576040517fa2376fe800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60208301516040517f85b214cf0000000000000000000000000000000000000000000000000000000081526004810184905273ffffffffffffffffffffffffffffffffffffffff909116906385b214cf906024016020604051808303816000875af1158015612d32573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612d5691906158a8565b5060408084015167ffffffffffffffff831660009081526003602052918220600101805491929091612d979084906bffffffffffffffffffffffff166157b2565b82546bffffffffffffffffffffffff9182166101009390930a928302919092021990911617905550606083015173ffffffffffffffffffffffffffffffffffffffff16600090815260046020908152604080832067ffffffffffffffff808616855292529091208054600192600991612e1f91859169010000000000000000009004166158ca565b825467ffffffffffffffff9182166101009390930a9283029190920219909116179055506000828152600560205260408082208290555183917ff1ca1e9147be737b04a2b018a79405f687a97de8dd8a2559bbe62357343af41491a250505080612e889061577a565b9050612bc9565b67ffffffffffffffff8116600090815260036020908152604080832060020180548251818502810185019093528083528493830182828015612f0757602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311612edc575b5050505050905060005b8151811015612fd457600060046000848481518110612f3257612f326154c8565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040908101600090812067ffffffffffffffff808a168352908452908290208251606081018452905460ff8116151582526101008104831694820185905269010000000000000000009004909116918101829052925014612fc357506001949350505050565b50612fcd8161577a565b9050612f11565b5060009392505050565b612fe6613366565b600955565b612ff3613366565b61075481614103565b61101f613366565b67ffffffffffffffff81166000908152600360205260409020546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff16610754576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff82166000908152600360209081526040808320815160c08101835281546bffffffffffffffffffffffff808216835273ffffffffffffffffffffffffffffffffffffffff6c010000000000000000000000009283900481168488015260018501549182168487015291900416606082015260028201805484518187028101870190955280855291949293608086019390929083018282801561315b57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311613130575b505050918352505060039190910154602090910152805190915060005b82608001515181101561321c57600460008460800151838151811061319f5761319f6154c8565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040908101600090812067ffffffffffffffff89168252909252902080547fffffffffffffffffffffffffffffff00000000000000000000000000000000001690556132158161577a565b9050613178565b5067ffffffffffffffff84166000908152600360205260408120818155600181018290559061324e60028301826148a2565b5060006003919091018190558054829190819061327a9084906bffffffffffffffffffffffff166157b2565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555061330183826bffffffffffffffffffffffff167f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1661401b9092919063ffffffff16565b6040805173ffffffffffffffffffffffffffffffffffffffff851681526bffffffffffffffffffffffff8316602082015267ffffffffffffffff8616917fe8ed5b475a5b5987aa9165e8731bb78043f39eee32ec5a1169a89e27fcd498159101611e96565b600654610100900473ffffffffffffffffffffffffffffffffffffffff16331461101f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016107b3565b61101f6141ff565b60006bffffffffffffffffffffffff821115613492576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f53616665436173743a2076616c756520646f65736e27742066697420696e203960448201527f362062697473000000000000000000000000000000000000000000000000000060648201526084016107b3565b5090565b60408051606080820183526000808352602083015291810191909152813b1580156134e957505060408051606081018252600080825260208083018290528351918252810183529181019190915261364b565b600a546040516000916b010000000000000000000000900460e01b90613517908a908a908a906024016158eb565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152918152602080830180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009590951694909417909352600a548151608480825260c0820190935292945061ffff6f01000000000000000000000000000000909104169260009283928392820181803683370190505090505a848110156135e557600080fd5b8490036040810481038a106135f957600080fd5b505a60008087516020890160008d8ff193505a900391503d608481111561361e575060845b808252806000602084013e5060408051606081018252931515845260208401929092529082015293505050505b95945050505050565b6040805180820190915260008082526020820152600061367484866156aa565b905060008161368388866156d2565b61368d91906156d2565b67ffffffffffffffff8b166000908152600360205260409020549091506bffffffffffffffffffffffff808316911610156137205767ffffffffffffffff8a16600090815260036020526040908190205490517f6b0fe56f0000000000000000000000000000000000000000000000000000000081526bffffffffffffffffffffffff90911660048201526024016107b3565b67ffffffffffffffff8a16600090815260036020526040812080548392906137579084906bffffffffffffffffffffffff166157b2565b82546101009290920a6bffffffffffffffffffffffff81810219909316918316021790915567ffffffffffffffff8c166000908152600360205260409020600101548b82169116101590506138045767ffffffffffffffff8a16600090815260036020526040908190205490517f6b0fe56f0000000000000000000000000000000000000000000000000000000081526bffffffffffffffffffffffff90911660048201526024016107b3565b67ffffffffffffffff8a16600090815260036020526040812060010180548b929061383e9084906bffffffffffffffffffffffff166157b2565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550818461387891906156d2565b33600090815260016020526040812080549091906138a59084906bffffffffffffffffffffffff166156d2565b82546101009290920a6bffffffffffffffffffffffff81810219909316918316021790915530600090815260016020526040812080548b945090926138ec918591166156d2565b82546bffffffffffffffffffffffff9182166101009390930a92830291909202199091161790555073ffffffffffffffffffffffffffffffffffffffff8816600090815260046020908152604080832067ffffffffffffffff808f1685529252909120805460019260099161397091859169010000000000000000009004166158ca565b92506101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055506040518060400160405280836bffffffffffffffffffffffff168152602001826bffffffffffffffffffffffff1681525092505050979650505050505050565b6139de61426c565b600680547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa335b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390a1565b6000613a5d6133ec565b613a6685613004565b613a7033866142d8565b613a7a8583610757565b8351600003613ab4576040517ec1cfc000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000613abf86612031565b90506000613acd338861154f565b905060008873ffffffffffffffffffffffffffffffffffffffff1663a631571e6040518061016001604052808a8152602001613b238c67ffffffffffffffff166000908152600360208190526040909120015490565b815233602082015260408781015188519190920191613b41916157b2565b6bffffffffffffffffffffffff168152602001600a60000160029054906101000a900468ffffffffffffffffff1668ffffffffffffffffff1681526020018b67ffffffffffffffff168152602001856020015167ffffffffffffffff1681526020018863ffffffff1681526020018961ffff168152602001856040015167ffffffffffffffff168152602001866020015173ffffffffffffffffffffffffffffffffffffffff168152506040518263ffffffff1660e01b8152600401613c079190615916565b610160604051808303816000875af1158015613c27573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613c4b9190615a7b565b805160009081526005602052604090205490915015613c9c5780516040517f304f32e800000000000000000000000000000000000000000000000000000000815260048101919091526024016107b3565b604051806101600160405280826000015181526020018a73ffffffffffffffffffffffffffffffffffffffff16815260200182604001516bffffffffffffffffffffffff1681526020013373ffffffffffffffffffffffffffffffffffffffff1681526020018967ffffffffffffffff1681526020018663ffffffff168152602001600a60000160029054906101000a900468ffffffffffffffffff1668ffffffffffffffffff1681526020018260e0015168ffffffffffffffffff16815260200182610100015164ffffffffff16815260200182610120015164ffffffffff16815260200182610140015163ffffffff16815250604051602001613da19190615529565b60405160208183030381529060405280519060200120600560008360000151815260200190815260200160002081905550613de13389836040015161434c565b8767ffffffffffffffff168a82600001517ff67aec45c9a7ede407974a3e0c3a743dffeab99ee3f2d4c9a8144c2ebf2c7ec9866020015133328d8d8d8a60400151604051613e359796959493929190615b4e565b60405180910390a4519998505050505050505050565b67ffffffffffffffff81166000908152600360205260409020546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1680613ec2576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff821614612559576040517f5a68151d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60095460009081526008602052604090205473ffffffffffffffffffffffffffffffffffffffff1680613f415750565b604080516000815260208101918290527f6b14daf80000000000000000000000000000000000000000000000000000000090915273ffffffffffffffffffffffffffffffffffffffff821690636b14daf890613fa290339060248101615bc6565b602060405180830381865afa158015613fbf573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613fe391906158a8565b610754576040517f229062630000000000000000000000000000000000000000000000000000000081523360048201526024016107b3565b6040805173ffffffffffffffffffffffffffffffffffffffff8416602482015260448082018490528251808303909101815260649091019091526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fa9059cbb000000000000000000000000000000000000000000000000000000001790526113f8908490614427565b6140b06141ff565b600680547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258613a293390565b3373ffffffffffffffffffffffffffffffffffffffff821603614182576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016107b3565b600780547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff838116918217909255600654604051919261010090910416907fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127890600090a350565b60065460ff161561101f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f5061757361626c653a207061757365640000000000000000000000000000000060448201526064016107b3565b60065460ff1661101f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f5061757361626c653a206e6f742070617573656400000000000000000000000060448201526064016107b3565b73ffffffffffffffffffffffffffffffffffffffff8216600090815260046020908152604080832067ffffffffffffffff8516845290915290205460ff16612559576040517f71e8313700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff8216600090815260036020526040812060010180548392906143869084906bffffffffffffffffffffffff166156d2565b82546bffffffffffffffffffffffff91821661010093840a908102920219161790915573ffffffffffffffffffffffffffffffffffffffff8516600090815260046020908152604080832067ffffffffffffffff80891685529252909120805460019450909284926143fc9284929004166158ca565b92506101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550505050565b6000614489826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff166145339092919063ffffffff16565b8051909150156113f857808060200190518101906144a791906158a8565b6113f8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f7420737563636565640000000000000000000000000000000000000000000060648201526084016107b3565b6060614542848460008561454a565b949350505050565b6060824710156145dc576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f60448201527f722063616c6c000000000000000000000000000000000000000000000000000060648201526084016107b3565b6000808673ffffffffffffffffffffffffffffffffffffffff1685876040516146059190615bf5565b60006040518083038185875af1925050503d8060008114614642576040519150601f19603f3d011682016040523d82523d6000602084013e614647565b606091505b509150915061465887838387614663565b979650505050505050565b606083156146f95782516000036146f25773ffffffffffffffffffffffffffffffffffffffff85163b6146f2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e747261637400000060448201526064016107b3565b5081614542565b614542838381511561470e5781518083602001fd5b806040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107b391906149df565b828054828255906000526020600020906007016008900481019282156147e15791602002820160005b838211156147af57835183826101000a81548163ffffffff021916908363ffffffff160217905550926020019260040160208160030104928301926001030261476b565b80156147df5782816101000a81549063ffffffff02191690556004016020816003010492830192600103026147af565b505b506134929291506148bc565b8280548282559060005260206000209081019282156147e1579160200282015b828111156147e157825182559160200191906001019061480d565b8280548282559060005260206000209081019282156147e1579160200282015b828111156147e157825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190614848565b508054600082559060005260206000209081019061075491905b5b8082111561349257600081556001016148bd565b67ffffffffffffffff8116811461075457600080fd5b80356148f2816148d1565b919050565b60006020828403121561490957600080fd5b8135614914816148d1565b9392505050565b63ffffffff8116811461075457600080fd5b80356148f28161491b565b6000806040838503121561494b57600080fd5b8235614956816148d1565b915060208301356149668161491b565b809150509250929050565b60005b8381101561498c578181015183820152602001614974565b50506000910152565b600081518084526149ad816020860160208601614971565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6020815260006149146020830184614995565b60008060408385031215614a0557600080fd5b8235614a10816148d1565b946020939093013593505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160a0810167ffffffffffffffff81118282101715614a7057614a70614a1e565b60405290565b604051610160810167ffffffffffffffff81118282101715614a7057614a70614a1e565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715614ae157614ae1614a1e565b604052919050565b803561ffff811681146148f257600080fd5b68ffffffffffffffffff8116811461075457600080fd5b80356148f281614afb565b600067ffffffffffffffff821115614b3757614b37614a1e565b5060051b60200190565b600082601f830112614b5257600080fd5b81356020614b67614b6283614b1d565b614a9a565b82815260059290921b84018101918181019086841115614b8657600080fd5b8286015b84811015614baa578035614b9d8161491b565b8352918301918301614b8a565b509695505050505050565b600060208284031215614bc757600080fd5b813567ffffffffffffffff80821115614bdf57600080fd5b9083019060a08286031215614bf357600080fd5b614bfb614a4d565b614c0483614ae9565b81526020830135614c1481614afb565b602082015260408301357fffffffff0000000000000000000000000000000000000000000000000000000081168114614c4c57600080fd5b6040820152614c5d60608401614ae9565b6060820152608083013582811115614c7457600080fd5b614c8087828601614b41565b60808301525095945050505050565b600082601f830112614ca057600080fd5b813567ffffffffffffffff811115614cba57614cba614a1e565b614ceb60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601614a9a565b818152846020838601011115614d0057600080fd5b816020850160208301376000918101602001919091529392505050565b6bffffffffffffffffffffffff8116811461075457600080fd5b80356148f281614d1d565b73ffffffffffffffffffffffffffffffffffffffff8116811461075457600080fd5b80356148f281614d42565b64ffffffffff8116811461075457600080fd5b80356148f281614d6f565b60006101608284031215614da057600080fd5b614da8614a76565b905081358152614dba60208301614d64565b6020820152614dcb60408301614d37565b6040820152614ddc60608301614d64565b6060820152614ded608083016148e7565b6080820152614dfe60a0830161492d565b60a0820152614e0f60c08301614b12565b60c0820152614e2060e08301614b12565b60e0820152610100614e33818401614d82565b90820152610120614e45838201614d82565b90820152610140614e5783820161492d565b9082015292915050565b6000806000806000806102008789031215614e7b57600080fd5b863567ffffffffffffffff80821115614e9357600080fd5b614e9f8a838b01614c8f565b97506020890135915080821115614eb557600080fd5b50614ec289828a01614c8f565b9550506040870135614ed381614d1d565b93506060870135614ee381614d1d565b92506080870135614ef381614d42565b9150614f028860a08901614d8d565b90509295509295509295565b60078110614f45577f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b9052565b60408101614f578285614f0e565b6bffffffffffffffffffffffff831660208301529392505050565b600082601f830112614f8357600080fd5b81356020614f93614b6283614b1d565b82815260059290921b84018101918181019086841115614fb257600080fd5b8286015b84811015614baa578035614fc981614d42565b8352918301918301614fb6565b60008060408385031215614fe957600080fd5b823567ffffffffffffffff8082111561500157600080fd5b818501915085601f83011261501557600080fd5b81356020615025614b6283614b1d565b82815260059290921b8401810191818101908984111561504457600080fd5b948201945b8386101561506257853582529482019490820190615049565b9650508601359250508082111561507857600080fd5b5061508585828601614f72565b9150509250929050565b60008083601f8401126150a157600080fd5b50813567ffffffffffffffff8111156150b957600080fd5b6020830191508360208285010111156150d157600080fd5b9250929050565b60008060008060008060a087890312156150f157600080fd5b86356150fc816148d1565b9550602087013567ffffffffffffffff81111561511857600080fd5b61512489828a0161508f565b9096509450615137905060408801614ae9565b925060608701356151478161491b565b80925050608087013590509295509295509295565b6000806040838503121561516f57600080fd5b823561517a816148d1565b9150602083013561496681614d42565b6000806040838503121561519d57600080fd5b82356151a881614d42565b9150602083013561496681614d1d565b600080604083850312156151cb57600080fd5b82356151d681614d42565b91506020830135614966816148d1565b6000602082840312156151f857600080fd5b5035919050565b600081518084526020808501945080840160005b8381101561524557815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101615213565b509495945050505050565b6020815260006bffffffffffffffffffffffff808451166020840152602084015173ffffffffffffffffffffffffffffffffffffffff8082166040860152826040870151166060860152806060870151166080860152505050608083015160c060a08401526152c260e08401826151ff565b905060a084015160c08401528091505092915050565b600080600080606085870312156152ee57600080fd5b84356152f981614d42565b935060208501359250604085013567ffffffffffffffff81111561531c57600080fd5b6153288782880161508f565b95989497509550505050565b604080825283519082018190526000906020906060840190828701845b8281101561536d57815184529284019290840190600101615351565b5050508381038285015261538181866151ff565b9695505050505050565b6000602080835260c0830161ffff808651168386015268ffffffffffffffffff838701511660408601527fffffffff00000000000000000000000000000000000000000000000000000000604087015116606086015280606087015116608086015250608085015160a08086015281815180845260e0870191508483019350600092505b80831015614baa57835163ffffffff16825292840192600192909201919084019061540f565b60006020828403121561544757600080fd5b813561491481614d42565b6000806020838503121561546557600080fd5b823567ffffffffffffffff8082111561547d57600080fd5b818501915085601f83011261549157600080fd5b8135818111156154a057600080fd5b866020610160830285010111156154b657600080fd5b60209290920196919550909350505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b73ffffffffffffffffffffffffffffffffffffffff848116825283166020820152606081016145426040830184614f0e565b8151815260208083015161016083019161555a9084018273ffffffffffffffffffffffffffffffffffffffff169052565b50604083015161557a60408401826bffffffffffffffffffffffff169052565b5060608301516155a2606084018273ffffffffffffffffffffffffffffffffffffffff169052565b5060808301516155be608084018267ffffffffffffffff169052565b5060a08301516155d660a084018263ffffffff169052565b5060c08301516155f360c084018268ffffffffffffffffff169052565b5060e083015161561060e084018268ffffffffffffffffff169052565b506101008381015164ffffffffff81168483015250506101208381015164ffffffffff81168483015250506101408381015163ffffffff8116848301525b505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b64ffffffffff8181168382160190808211156156a3576156a3615656565b5092915050565b6bffffffffffffffffffffffff81811683821602808216919082811461564e5761564e615656565b6bffffffffffffffffffffffff8181168382160190808211156156a3576156a3615656565b6bffffffffffffffffffffffff8716815273ffffffffffffffffffffffffffffffffffffffff861660208201526157316040820186614f0e565b60c06060820152600061574760c0830186614995565b82810360808401526157598186614995565b905082810360a084015261576d8185614995565b9998505050505050505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036157ab576157ab615656565b5060010190565b6bffffffffffffffffffffffff8281168282160390808211156156a3576156a3615656565b600060ff821660ff81036157ed576157ed615656565b60010192915050565b818103818111156115d9576115d9615656565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b600067ffffffffffffffff80831681810361585557615855615656565b6001019392505050565b808201808211156115d9576115d9615656565b60006020828403121561588457600080fd5b5051919050565b6000610160828403121561589e57600080fd5b6149148383614d8d565b6000602082840312156158ba57600080fd5b8151801515811461491457600080fd5b67ffffffffffffffff8181168382160190808211156156a3576156a3615656565b8381526060602082015260006159046060830185614995565b82810360408401526153818185614995565b6020815260008251610160806020850152615935610180850183614995565b9150602085015160408501526040850151615968606086018273ffffffffffffffffffffffffffffffffffffffff169052565b5060608501516bffffffffffffffffffffffff8116608086015250608085015168ffffffffffffffffff811660a08601525060a085015167ffffffffffffffff811660c08601525060c085015167ffffffffffffffff811660e08601525060e08501516101006159df8187018363ffffffff169052565b86015190506101206159f68682018361ffff169052565b8601519050610140615a138682018367ffffffffffffffff169052565b9095015173ffffffffffffffffffffffffffffffffffffffff1693019290925250919050565b80516148f281614d42565b80516148f281614d1d565b80516148f2816148d1565b80516148f28161491b565b80516148f281614afb565b80516148f281614d6f565b60006101608284031215615a8e57600080fd5b615a96614a76565b82518152615aa660208401615a39565b6020820152615ab760408401615a44565b6040820152615ac860608401615a39565b6060820152615ad960808401615a4f565b6080820152615aea60a08401615a5a565b60a0820152615afb60c08401615a65565b60c0820152615b0c60e08401615a65565b60e0820152610100615b1f818501615a70565b90820152610120615b31848201615a70565b90820152610140615b43848201615a5a565b908201529392505050565b600073ffffffffffffffffffffffffffffffffffffffff808a168352808916602084015280881660408401525060e06060830152615b8f60e0830187614995565b61ffff9590951660808301525063ffffffff9290921660a08301526bffffffffffffffffffffffff1660c090910152949350505050565b73ffffffffffffffffffffffffffffffffffffffff831681526040602082015260006145426040830184614995565b60008251615c07818460208701614971565b919091019291505056fea164736f6c6343000813000a",
}

var FunctionsRouterABI = FunctionsRouterMetaData.ABI

var FunctionsRouterBin = FunctionsRouterMetaData.Bin

func DeployFunctionsRouter(auth *bind.TransactOpts, backend bind.ContractBackend, linkToken common.Address, config FunctionsRouterConfig) (common.Address, *types.Transaction, *FunctionsRouter, error) {
	parsed, err := FunctionsRouterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(FunctionsRouterBin), backend, linkToken, config)
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

func (_FunctionsRouter *FunctionsRouterCaller) GetConfig(opts *bind.CallOpts) (FunctionsRouterConfig, error) {
	var out []interface{}
	err := _FunctionsRouter.contract.Call(opts, &out, "getConfig")

	if err != nil {
		return *new(FunctionsRouterConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(FunctionsRouterConfig)).(*FunctionsRouterConfig)

	return out0, err

}

func (_FunctionsRouter *FunctionsRouterSession) GetConfig() (FunctionsRouterConfig, error) {
	return _FunctionsRouter.Contract.GetConfig(&_FunctionsRouter.CallOpts)
}

func (_FunctionsRouter *FunctionsRouterCallerSession) GetConfig() (FunctionsRouterConfig, error) {
	return _FunctionsRouter.Contract.GetConfig(&_FunctionsRouter.CallOpts)
}

func (_FunctionsRouter *FunctionsRouterCaller) GetConsumer(opts *bind.CallOpts, client common.Address, subscriptionId uint64) (IFunctionsSubscriptionsConsumer, error) {
	var out []interface{}
	err := _FunctionsRouter.contract.Call(opts, &out, "getConsumer", client, subscriptionId)

	if err != nil {
		return *new(IFunctionsSubscriptionsConsumer), err
	}

	out0 := *abi.ConvertType(out[0], new(IFunctionsSubscriptionsConsumer)).(*IFunctionsSubscriptionsConsumer)

	return out0, err

}

func (_FunctionsRouter *FunctionsRouterSession) GetConsumer(client common.Address, subscriptionId uint64) (IFunctionsSubscriptionsConsumer, error) {
	return _FunctionsRouter.Contract.GetConsumer(&_FunctionsRouter.CallOpts, client, subscriptionId)
}

func (_FunctionsRouter *FunctionsRouterCallerSession) GetConsumer(client common.Address, subscriptionId uint64) (IFunctionsSubscriptionsConsumer, error) {
	return _FunctionsRouter.Contract.GetConsumer(&_FunctionsRouter.CallOpts, client, subscriptionId)
}

func (_FunctionsRouter *FunctionsRouterCaller) GetContractById(opts *bind.CallOpts, id [32]byte) (common.Address, error) {
	var out []interface{}
	err := _FunctionsRouter.contract.Call(opts, &out, "getContractById", id)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_FunctionsRouter *FunctionsRouterSession) GetContractById(id [32]byte) (common.Address, error) {
	return _FunctionsRouter.Contract.GetContractById(&_FunctionsRouter.CallOpts, id)
}

func (_FunctionsRouter *FunctionsRouterCallerSession) GetContractById(id [32]byte) (common.Address, error) {
	return _FunctionsRouter.Contract.GetContractById(&_FunctionsRouter.CallOpts, id)
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

func (_FunctionsRouter *FunctionsRouterCaller) GetProposedContractById(opts *bind.CallOpts, id [32]byte) (common.Address, error) {
	var out []interface{}
	err := _FunctionsRouter.contract.Call(opts, &out, "getProposedContractById", id)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_FunctionsRouter *FunctionsRouterSession) GetProposedContractById(id [32]byte) (common.Address, error) {
	return _FunctionsRouter.Contract.GetProposedContractById(&_FunctionsRouter.CallOpts, id)
}

func (_FunctionsRouter *FunctionsRouterCallerSession) GetProposedContractById(id [32]byte) (common.Address, error) {
	return _FunctionsRouter.Contract.GetProposedContractById(&_FunctionsRouter.CallOpts, id)
}

func (_FunctionsRouter *FunctionsRouterCaller) GetProposedContractSet(opts *bind.CallOpts) ([][32]byte, []common.Address, error) {
	var out []interface{}
	err := _FunctionsRouter.contract.Call(opts, &out, "getProposedContractSet")

	if err != nil {
		return *new([][32]byte), *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)
	out1 := *abi.ConvertType(out[1], new([]common.Address)).(*[]common.Address)

	return out0, out1, err

}

func (_FunctionsRouter *FunctionsRouterSession) GetProposedContractSet() ([][32]byte, []common.Address, error) {
	return _FunctionsRouter.Contract.GetProposedContractSet(&_FunctionsRouter.CallOpts)
}

func (_FunctionsRouter *FunctionsRouterCallerSession) GetProposedContractSet() ([][32]byte, []common.Address, error) {
	return _FunctionsRouter.Contract.GetProposedContractSet(&_FunctionsRouter.CallOpts)
}

func (_FunctionsRouter *FunctionsRouterCaller) GetSubscription(opts *bind.CallOpts, subscriptionId uint64) (IFunctionsSubscriptionsSubscription, error) {
	var out []interface{}
	err := _FunctionsRouter.contract.Call(opts, &out, "getSubscription", subscriptionId)

	if err != nil {
		return *new(IFunctionsSubscriptionsSubscription), err
	}

	out0 := *abi.ConvertType(out[0], new(IFunctionsSubscriptionsSubscription)).(*IFunctionsSubscriptionsSubscription)

	return out0, err

}

func (_FunctionsRouter *FunctionsRouterSession) GetSubscription(subscriptionId uint64) (IFunctionsSubscriptionsSubscription, error) {
	return _FunctionsRouter.Contract.GetSubscription(&_FunctionsRouter.CallOpts, subscriptionId)
}

func (_FunctionsRouter *FunctionsRouterCallerSession) GetSubscription(subscriptionId uint64) (IFunctionsSubscriptionsSubscription, error) {
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

func (_FunctionsRouter *FunctionsRouterTransactor) CreateSubscriptionWithConsumer(opts *bind.TransactOpts, consumer common.Address) (*types.Transaction, error) {
	return _FunctionsRouter.contract.Transact(opts, "createSubscriptionWithConsumer", consumer)
}

func (_FunctionsRouter *FunctionsRouterSession) CreateSubscriptionWithConsumer(consumer common.Address) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.CreateSubscriptionWithConsumer(&_FunctionsRouter.TransactOpts, consumer)
}

func (_FunctionsRouter *FunctionsRouterTransactorSession) CreateSubscriptionWithConsumer(consumer common.Address) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.CreateSubscriptionWithConsumer(&_FunctionsRouter.TransactOpts, consumer)
}

func (_FunctionsRouter *FunctionsRouterTransactor) Fulfill(opts *bind.TransactOpts, response []byte, err []byte, juelsPerGas *big.Int, costWithoutCallback *big.Int, transmitter common.Address, commitment FunctionsResponseCommitment) (*types.Transaction, error) {
	return _FunctionsRouter.contract.Transact(opts, "fulfill", response, err, juelsPerGas, costWithoutCallback, transmitter, commitment)
}

func (_FunctionsRouter *FunctionsRouterSession) Fulfill(response []byte, err []byte, juelsPerGas *big.Int, costWithoutCallback *big.Int, transmitter common.Address, commitment FunctionsResponseCommitment) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.Fulfill(&_FunctionsRouter.TransactOpts, response, err, juelsPerGas, costWithoutCallback, transmitter, commitment)
}

func (_FunctionsRouter *FunctionsRouterTransactorSession) Fulfill(response []byte, err []byte, juelsPerGas *big.Int, costWithoutCallback *big.Int, transmitter common.Address, commitment FunctionsResponseCommitment) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.Fulfill(&_FunctionsRouter.TransactOpts, response, err, juelsPerGas, costWithoutCallback, transmitter, commitment)
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

func (_FunctionsRouter *FunctionsRouterTransactor) SendRequestToProposed(opts *bind.TransactOpts, subscriptionId uint64, data []byte, dataVersion uint16, callbackGasLimit uint32, donId [32]byte) (*types.Transaction, error) {
	return _FunctionsRouter.contract.Transact(opts, "sendRequestToProposed", subscriptionId, data, dataVersion, callbackGasLimit, donId)
}

func (_FunctionsRouter *FunctionsRouterSession) SendRequestToProposed(subscriptionId uint64, data []byte, dataVersion uint16, callbackGasLimit uint32, donId [32]byte) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.SendRequestToProposed(&_FunctionsRouter.TransactOpts, subscriptionId, data, dataVersion, callbackGasLimit, donId)
}

func (_FunctionsRouter *FunctionsRouterTransactorSession) SendRequestToProposed(subscriptionId uint64, data []byte, dataVersion uint16, callbackGasLimit uint32, donId [32]byte) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.SendRequestToProposed(&_FunctionsRouter.TransactOpts, subscriptionId, data, dataVersion, callbackGasLimit, donId)
}

func (_FunctionsRouter *FunctionsRouterTransactor) SetAllowListId(opts *bind.TransactOpts, allowListId [32]byte) (*types.Transaction, error) {
	return _FunctionsRouter.contract.Transact(opts, "setAllowListId", allowListId)
}

func (_FunctionsRouter *FunctionsRouterSession) SetAllowListId(allowListId [32]byte) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.SetAllowListId(&_FunctionsRouter.TransactOpts, allowListId)
}

func (_FunctionsRouter *FunctionsRouterTransactorSession) SetAllowListId(allowListId [32]byte) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.SetAllowListId(&_FunctionsRouter.TransactOpts, allowListId)
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

func (_FunctionsRouter *FunctionsRouterTransactor) TimeoutRequests(opts *bind.TransactOpts, requestsToTimeoutByCommitment []FunctionsResponseCommitment) (*types.Transaction, error) {
	return _FunctionsRouter.contract.Transact(opts, "timeoutRequests", requestsToTimeoutByCommitment)
}

func (_FunctionsRouter *FunctionsRouterSession) TimeoutRequests(requestsToTimeoutByCommitment []FunctionsResponseCommitment) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.TimeoutRequests(&_FunctionsRouter.TransactOpts, requestsToTimeoutByCommitment)
}

func (_FunctionsRouter *FunctionsRouterTransactorSession) TimeoutRequests(requestsToTimeoutByCommitment []FunctionsResponseCommitment) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.TimeoutRequests(&_FunctionsRouter.TransactOpts, requestsToTimeoutByCommitment)
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

func (_FunctionsRouter *FunctionsRouterTransactor) UpdateConfig(opts *bind.TransactOpts, config FunctionsRouterConfig) (*types.Transaction, error) {
	return _FunctionsRouter.contract.Transact(opts, "updateConfig", config)
}

func (_FunctionsRouter *FunctionsRouterSession) UpdateConfig(config FunctionsRouterConfig) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.UpdateConfig(&_FunctionsRouter.TransactOpts, config)
}

func (_FunctionsRouter *FunctionsRouterTransactorSession) UpdateConfig(config FunctionsRouterConfig) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.UpdateConfig(&_FunctionsRouter.TransactOpts, config)
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
	Arg0 FunctionsRouterConfig
	Raw  types.Log
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
	Id   [32]byte
	From common.Address
	To   common.Address
	Raw  types.Log
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

type FunctionsRouterRequestNotProcessedIterator struct {
	Event *FunctionsRouterRequestNotProcessed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRouterRequestNotProcessedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRouterRequestNotProcessed)
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
		it.Event = new(FunctionsRouterRequestNotProcessed)
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

func (it *FunctionsRouterRequestNotProcessedIterator) Error() error {
	return it.fail
}

func (it *FunctionsRouterRequestNotProcessedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRouterRequestNotProcessed struct {
	RequestId   [32]byte
	Coordinator common.Address
	Transmitter common.Address
	ResultCode  uint8
	Raw         types.Log
}

func (_FunctionsRouter *FunctionsRouterFilterer) FilterRequestNotProcessed(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsRouterRequestNotProcessedIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsRouter.contract.FilterLogs(opts, "RequestNotProcessed", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsRouterRequestNotProcessedIterator{contract: _FunctionsRouter.contract, event: "RequestNotProcessed", logs: logs, sub: sub}, nil
}

func (_FunctionsRouter *FunctionsRouterFilterer) WatchRequestNotProcessed(opts *bind.WatchOpts, sink chan<- *FunctionsRouterRequestNotProcessed, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsRouter.contract.WatchLogs(opts, "RequestNotProcessed", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRouterRequestNotProcessed)
				if err := _FunctionsRouter.contract.UnpackLog(event, "RequestNotProcessed", log); err != nil {
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

func (_FunctionsRouter *FunctionsRouterFilterer) ParseRequestNotProcessed(log types.Log) (*FunctionsRouterRequestNotProcessed, error) {
	event := new(FunctionsRouterRequestNotProcessed)
	if err := _FunctionsRouter.contract.UnpackLog(event, "RequestNotProcessed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsRouterRequestProcessedIterator struct {
	Event *FunctionsRouterRequestProcessed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRouterRequestProcessedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRouterRequestProcessed)
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
		it.Event = new(FunctionsRouterRequestProcessed)
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

func (it *FunctionsRouterRequestProcessedIterator) Error() error {
	return it.fail
}

func (it *FunctionsRouterRequestProcessedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRouterRequestProcessed struct {
	RequestId          [32]byte
	SubscriptionId     uint64
	TotalCostJuels     *big.Int
	Transmitter        common.Address
	ResultCode         uint8
	Response           []byte
	Err                []byte
	CallbackReturnData []byte
	Raw                types.Log
}

func (_FunctionsRouter *FunctionsRouterFilterer) FilterRequestProcessed(opts *bind.FilterOpts, requestId [][32]byte, subscriptionId []uint64) (*FunctionsRouterRequestProcessedIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsRouter.contract.FilterLogs(opts, "RequestProcessed", requestIdRule, subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsRouterRequestProcessedIterator{contract: _FunctionsRouter.contract, event: "RequestProcessed", logs: logs, sub: sub}, nil
}

func (_FunctionsRouter *FunctionsRouterFilterer) WatchRequestProcessed(opts *bind.WatchOpts, sink chan<- *FunctionsRouterRequestProcessed, requestId [][32]byte, subscriptionId []uint64) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsRouter.contract.WatchLogs(opts, "RequestProcessed", requestIdRule, subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRouterRequestProcessed)
				if err := _FunctionsRouter.contract.UnpackLog(event, "RequestProcessed", log); err != nil {
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

func (_FunctionsRouter *FunctionsRouterFilterer) ParseRequestProcessed(log types.Log) (*FunctionsRouterRequestProcessed, error) {
	event := new(FunctionsRouterRequestProcessed)
	if err := _FunctionsRouter.contract.UnpackLog(event, "RequestProcessed", log); err != nil {
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
	RequestId               [32]byte
	DonId                   [32]byte
	SubscriptionId          uint64
	SubscriptionOwner       common.Address
	RequestingContract      common.Address
	RequestInitiator        common.Address
	Data                    []byte
	DataVersion             uint16
	CallbackGasLimit        uint32
	EstimatedTotalCostJuels *big.Int
	Raw                     types.Log
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

func (_FunctionsRouter *FunctionsRouter) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
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
	case _FunctionsRouter.abi.Events["RequestNotProcessed"].ID:
		return _FunctionsRouter.ParseRequestNotProcessed(log)
	case _FunctionsRouter.abi.Events["RequestProcessed"].ID:
		return _FunctionsRouter.ParseRequestProcessed(log)
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
	case _FunctionsRouter.abi.Events["Unpaused"].ID:
		return _FunctionsRouter.ParseUnpaused(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (FunctionsRouterConfigUpdated) Topic() common.Hash {
	return common.HexToHash("0x00a5832bf95f66c7814294cc4db681f20ee79608bfb8912a5321d66cfed5e985")
}

func (FunctionsRouterContractProposed) Topic() common.Hash {
	return common.HexToHash("0x8b052f0f4bf82fede7daffea71592b29d5ef86af1f3c7daaa0345dbb2f52f481")
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

func (FunctionsRouterRequestNotProcessed) Topic() common.Hash {
	return common.HexToHash("0x1a90e9a50793db2e394cf581e7c522e10c358a81e70acf6b5a0edd620c08dee1")
}

func (FunctionsRouterRequestProcessed) Topic() common.Hash {
	return common.HexToHash("0x64778f26c70b60a8d7e29e2451b3844302d959448401c0535b768ed88c6b505e")
}

func (FunctionsRouterRequestStart) Topic() common.Hash {
	return common.HexToHash("0xf67aec45c9a7ede407974a3e0c3a743dffeab99ee3f2d4c9a8144c2ebf2c7ec9")
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

	GetConfig(opts *bind.CallOpts) (FunctionsRouterConfig, error)

	GetConsumer(opts *bind.CallOpts, client common.Address, subscriptionId uint64) (IFunctionsSubscriptionsConsumer, error)

	GetContractById(opts *bind.CallOpts, id [32]byte) (common.Address, error)

	GetFlags(opts *bind.CallOpts, subscriptionId uint64) ([32]byte, error)

	GetProposedContractById(opts *bind.CallOpts, id [32]byte) (common.Address, error)

	GetProposedContractSet(opts *bind.CallOpts) ([][32]byte, []common.Address, error)

	GetSubscription(opts *bind.CallOpts, subscriptionId uint64) (IFunctionsSubscriptionsSubscription, error)

	GetSubscriptionCount(opts *bind.CallOpts) (uint64, error)

	GetTotalBalance(opts *bind.CallOpts) (*big.Int, error)

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

	CreateSubscriptionWithConsumer(opts *bind.TransactOpts, consumer common.Address) (*types.Transaction, error)

	Fulfill(opts *bind.TransactOpts, response []byte, err []byte, juelsPerGas *big.Int, costWithoutCallback *big.Int, transmitter common.Address, commitment FunctionsResponseCommitment) (*types.Transaction, error)

	OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error)

	OracleWithdraw(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error)

	OwnerCancelSubscription(opts *bind.TransactOpts, subscriptionId uint64) (*types.Transaction, error)

	OwnerWithdraw(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error)

	Pause(opts *bind.TransactOpts) (*types.Transaction, error)

	ProposeContractsUpdate(opts *bind.TransactOpts, proposedContractSetIds [][32]byte, proposedContractSetAddresses []common.Address) (*types.Transaction, error)

	ProposeSubscriptionOwnerTransfer(opts *bind.TransactOpts, subscriptionId uint64, newOwner common.Address) (*types.Transaction, error)

	RecoverFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	RemoveConsumer(opts *bind.TransactOpts, subscriptionId uint64, consumer common.Address) (*types.Transaction, error)

	SendRequest(opts *bind.TransactOpts, subscriptionId uint64, data []byte, dataVersion uint16, callbackGasLimit uint32, donId [32]byte) (*types.Transaction, error)

	SendRequestToProposed(opts *bind.TransactOpts, subscriptionId uint64, data []byte, dataVersion uint16, callbackGasLimit uint32, donId [32]byte) (*types.Transaction, error)

	SetAllowListId(opts *bind.TransactOpts, allowListId [32]byte) (*types.Transaction, error)

	SetFlags(opts *bind.TransactOpts, subscriptionId uint64, flags [32]byte) (*types.Transaction, error)

	TimeoutRequests(opts *bind.TransactOpts, requestsToTimeoutByCommitment []FunctionsResponseCommitment) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Unpause(opts *bind.TransactOpts) (*types.Transaction, error)

	UpdateConfig(opts *bind.TransactOpts, config FunctionsRouterConfig) (*types.Transaction, error)

	UpdateContracts(opts *bind.TransactOpts) (*types.Transaction, error)

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

	FilterRequestNotProcessed(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsRouterRequestNotProcessedIterator, error)

	WatchRequestNotProcessed(opts *bind.WatchOpts, sink chan<- *FunctionsRouterRequestNotProcessed, requestId [][32]byte) (event.Subscription, error)

	ParseRequestNotProcessed(log types.Log) (*FunctionsRouterRequestNotProcessed, error)

	FilterRequestProcessed(opts *bind.FilterOpts, requestId [][32]byte, subscriptionId []uint64) (*FunctionsRouterRequestProcessedIterator, error)

	WatchRequestProcessed(opts *bind.WatchOpts, sink chan<- *FunctionsRouterRequestProcessed, requestId [][32]byte, subscriptionId []uint64) (event.Subscription, error)

	ParseRequestProcessed(log types.Log) (*FunctionsRouterRequestProcessed, error)

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

	FilterUnpaused(opts *bind.FilterOpts) (*FunctionsRouterUnpausedIterator, error)

	WatchUnpaused(opts *bind.WatchOpts, sink chan<- *FunctionsRouterUnpaused) (event.Subscription, error)

	ParseUnpaused(log types.Log) (*FunctionsRouterUnpaused, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
