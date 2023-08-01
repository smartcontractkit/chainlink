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

var FunctionsRouterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"timelockBlocks\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"maximumTimelockBlocks\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"linkToken\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"internalBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"externalBalance\",\"type\":\"uint256\"}],\"name\":\"BalanceInvariantViolated\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConsumerRequestsInFlight\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"limit\",\"type\":\"uint32\"}],\"name\":\"GasLimitTooBig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"IdentifierIsReserved\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidConfigData\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"InvalidConsumer\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"value\",\"type\":\"uint8\"}],\"name\":\"InvalidGasFlagValue\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidProposal\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSubscription\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposedOwner\",\"type\":\"address\"}],\"name\":\"MustBeRequestedOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeSubscriptionOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByRoute\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableFromCoordinator\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableFromLink\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingRequestExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ProposedTimelockAboveMaximum\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Reentrant\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"RouteNotFound\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"SenderMustAcceptTermsOfService\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TimelockInEffect\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyConsumers\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"adminFee\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"bytes4\",\"name\":\"handleOracleFulfillmentSelector\",\"type\":\"bytes4\"},{\"indexed\":false,\"internalType\":\"uint32[]\",\"name\":\"maxCallbackGasLimits\",\"type\":\"uint32[]\"}],\"name\":\"ConfigChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"fromHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"toBytes\",\"type\":\"bytes\"}],\"name\":\"ConfigProposed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"fromHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"toBytes\",\"type\":\"bytes\"}],\"name\":\"ConfigUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"proposedContractSetId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"proposedContractSetFromAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"proposedContractSetToAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timelockEndBlock\",\"type\":\"uint256\"}],\"name\":\"ContractProposed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"proposedContractSetId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"proposedContractSetFromAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"proposedContractSetToAddress\",\"type\":\"address\"}],\"name\":\"ContractUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalCostJuels\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"resultCode\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"response\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"}],\"name\":\"RequestEnd\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"donId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"subscriptionOwner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requestingContract\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requestInitiator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"dataVersion\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"}],\"name\":\"RequestStart\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"RequestTimedOut\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"fundsRecipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fundsAmount\",\"type\":\"uint256\"}],\"name\":\"SubscriptionCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"SubscriptionCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"from\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"to\",\"type\":\"uint16\"}],\"name\":\"TimeLockProposed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"from\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"to\",\"type\":\"uint16\"}],\"name\":\"TimeLockUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"MAX_CALLBACK_RETURN_BYTES\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"acceptSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"addConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"cancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createSubscription\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"response\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"juelsPerGas\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"costWithoutFulfillment\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint96\",\"name\":\"adminFee\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"estimatedTotalCostJuels\",\"type\":\"uint96\"},{\"internalType\":\"uint40\",\"name\":\"timeoutTimestamp\",\"type\":\"uint40\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint80\",\"name\":\"donFee\",\"type\":\"uint80\"},{\"internalType\":\"uint40\",\"name\":\"gasOverheadBeforeCallback\",\"type\":\"uint40\"},{\"internalType\":\"uint40\",\"name\":\"gasOverheadAfterCallback\",\"type\":\"uint40\"}],\"internalType\":\"structIFunctionsRequest.Commitment\",\"name\":\"commitment\",\"type\":\"tuple\"}],\"name\":\"fulfill\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"resultCode\",\"type\":\"uint8\"},{\"internalType\":\"uint96\",\"name\":\"callbackGasCostJuels\",\"type\":\"uint96\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAdminFee\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowListId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfigHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"config\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"getConsumer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"initiatedRequests\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"completedRequests\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"useProposed\",\"type\":\"bool\"}],\"name\":\"getContractById\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"routeDestination\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"getContractById\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"routeDestination\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"getFlags\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMaxConsumers\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getProposedContractSet\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"},{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"getSubscription\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"blockedBalance\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"requestedOwner\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSubscriptionCount\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTotalBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isPaused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"}],\"name\":\"isValidCallbackGasLimit\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"oracleWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"ownerCancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"ownerWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"pendingRequestExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"name\":\"proposeConfigUpdate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"proposedContractSetIds\",\"type\":\"bytes32[]\"},{\"internalType\":\"address[]\",\"name\":\"proposedContractSetAddresses\",\"type\":\"address[]\"}],\"name\":\"proposeContractsUpdate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"proposeSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"blocks\",\"type\":\"uint16\"}],\"name\":\"proposeTimelockBlocks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"removeConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint16\",\"name\":\"dataVersion\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"donId\",\"type\":\"bytes32\"}],\"name\":\"sendRequest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"flags\",\"type\":\"bytes32\"}],\"name\":\"setFlags\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint96\",\"name\":\"adminFee\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"estimatedTotalCostJuels\",\"type\":\"uint96\"},{\"internalType\":\"uint40\",\"name\":\"timeoutTimestamp\",\"type\":\"uint40\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint80\",\"name\":\"donFee\",\"type\":\"uint80\"},{\"internalType\":\"uint40\",\"name\":\"gasOverheadBeforeCallback\",\"type\":\"uint40\"},{\"internalType\":\"uint40\",\"name\":\"gasOverheadAfterCallback\",\"type\":\"uint40\"}],\"internalType\":\"structIFunctionsRequest.Commitment[]\",\"name\":\"requestsToTimeoutByCommitment\",\"type\":\"tuple[]\"}],\"name\":\"timeoutRequests\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"updateConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"updateContracts\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"updateTimelockBlocks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"validateProposedContracts\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b506040516200636f3803806200636f833981016040819052620000349162000431565b6000805460ff191681558290339086908690859084908190816200009f5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0380851661010002610100600160a81b031990921691909117909155811615620000d957620000d98162000178565b50506009805461ffff80871661ffff19909216919091179091558316608052506000805260026020527fac33ff75c19e70fe83507db0d683fd3465c996598dc972688b7ace676c89077b80546001600160a01b031916301790556200013e8162000229565b80516020909101206007555050600a80546001600160a01b0319166001600160a01b03939093169290921790915550620006ab9350505050565b336001600160a01b03821603620001d25760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000096565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929361010090910416917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000806000838060200190518101906200024491906200052e565b604080516060810182526001600160601b0385168082526001600160e01b03198516602080840191909152928201849052600f80546001600160801b0319169091176c0100000000000000000000000060e087901c0217815583519598509396509194509092620002bc916010919086019062000303565b509050507fe6a1eda76d42a6d1d813f26765716562044db1e8bd8be7d088705e64afd301ca838383604051620002f5939291906200063a565b60405180910390a150505050565b82805482825590600052602060002090600701600890048101928215620003a75791602002820160005b838211156200037357835183826101000a81548163ffffffff021916908363ffffffff16021790555092602001926004016020816003010492830192600103026200032d565b8015620003a55782816101000a81549063ffffffff021916905560040160208160030104928301926001030262000373565b505b50620003b5929150620003b9565b5090565b5b80821115620003b55760008155600101620003ba565b805161ffff81168114620003e357600080fd5b919050565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f191681016001600160401b0381118282101715620004295762000429620003e8565b604052919050565b600080600080608085870312156200044857600080fd5b6200045385620003d0565b9350602062000464818701620003d0565b60408701519094506001600160a01b03811681146200048257600080fd5b60608701519093506001600160401b0380821115620004a057600080fd5b818801915088601f830112620004b557600080fd5b815181811115620004ca57620004ca620003e8565b620004de601f8201601f19168501620003fe565b91508082528984828501011115620004f557600080fd5b60005b8181101562000515578381018501518382018601528401620004f8565b5060008482840101525080935050505092959194509250565b6000806000606084860312156200054457600080fd5b83516001600160601b03811681146200055c57600080fd5b602085810151919450906001600160e01b0319811681146200057d57600080fd5b60408601519093506001600160401b03808211156200059b57600080fd5b818701915087601f830112620005b057600080fd5b815181811115620005c557620005c5620003e8565b8060051b9150620005d8848301620003fe565b818152918301840191848101908a841115620005f357600080fd5b938501935b838510156200062a578451925063ffffffff83168314620006195760008081fd5b8282529385019390850190620005f8565b8096505050505050509250925092565b6001600160601b03841681526001600160e01b031983166020808301919091526060604083018190528351908301819052600091848101916080850190845b818110156200069d57845163ffffffff168352938301939183019160010162000679565b509098975050505050505050565b608051615ca8620006c760003960006128980152615ca86000f3fe608060405234801561001057600080fd5b50600436106102ff5760003560e01c806379ba50971161019c578063a9c9a918116100ee578063c4614d6111610097578063e82ad7d411610071578063e82ad7d414610766578063eb523d6c14610779578063f2fde38b1461078c57600080fd5b8063c4614d611461072d578063d7ae1d3014610740578063e72f6e301461075357600080fd5b8063b5643858116100c8578063b5643858146106fb578063b734c0f41461070e578063badc3eb61461071657600080fd5b8063a9c9a918146106b7578063aab396bd146106ca578063b187bd26146106f057600080fd5b80639883c10d11610150578063a21a23e41161012a578063a21a23e414610678578063a47c769614610680578063a4c0ed36146106a457600080fd5b80639883c10d146106245780639f87fad71461062c578063a1d4c8281461063f57600080fd5b80638456cb59116101815780638456cb59146105f15780638da5cb5b146105f95780638fde53171461061c57600080fd5b806379ba5097146105d657806382359740146105de57600080fd5b80634b8832d31161025557806366419970116102095780636e3b3323116101e35780636e3b3323146105a957806371ec28ac146105bc5780637341c10c146105c357600080fd5b806366419970146104ae578063674603d0146104dc5780636a6df79b1461057157600080fd5b80635c975abb1161023a5780635c975abb146104715780635ed6dfba1461048857806366316d8d1461049b57600080fd5b80634b8832d31461044b57806355fedefa1461045e57600080fd5b80631ded3b36116102b75780633e871e4d116102915780633e871e4d1461040f5780633f4ba83a14610422578063461d27621461042a57600080fd5b80631ded3b36146103d35780632a905ccc146103e6578063385de9ae146103fc57600080fd5b806310fc49c1116102e857806310fc49c11461033957806312b583491461034c578063181f5a771461038a57600080fd5b806302bcc5b6146103045780630c5d49cb14610319575b600080fd5b610317610312366004614889565b61079f565b005b610321608481565b60405161ffff90911681526020015b60405180910390f35b6103176103473660046148c3565b61082c565b6009546b01000000000000000000000090046bffffffffffffffffffffffff165b6040516bffffffffffffffffffffffff9091168152602001610330565b6103c66040518060400160405280601781526020017f46756e6374696f6e7320526f757465722076312e302e3000000000000000000081525081565b6040516103309190614960565b6103176103e1366004614973565b610956565b600f546bffffffffffffffffffffffff1661036d565b61031761040a3660046149e8565b610987565b61031761041d366004614ba1565b610b65565b610317610e5b565b61043d610438366004614c6c565b610e6d565b604051908152602001610330565b610317610459366004614cf0565b610ec1565b61043d61046c366004614889565b610fbe565b60005460ff165b6040519015158152602001610330565b610317610496366004614d43565b610fe2565b6103176104a9366004614d43565b6112b8565b6009546301000000900467ffffffffffffffff165b60405167ffffffffffffffff9091168152602001610330565b6105496104ea366004614d71565b73ffffffffffffffffffffffffffffffffffffffff919091166000908152600c6020908152604080832067ffffffffffffffff948516845290915290205460ff8116926101008204831692690100000000000000000090920490911690565b60408051931515845267ffffffffffffffff9283166020850152911690820152606001610330565b61058461057f366004614dad565b6114df565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610330565b6103176105b7366004614dd2565b6114f2565b6064610321565b6103176105d1366004614cf0565b6117c8565b610317611975565b6103176105ec366004614889565b611a97565b610317611be6565b600054610100900473ffffffffffffffffffffffffffffffffffffffff16610584565b610317611bf6565b60075461043d565b61031761063a366004614cf0565b611c9e565b61065261064d366004614feb565b6120a6565b6040805160ff90931683526bffffffffffffffffffffffff909116602083015201610330565b6104c3612390565b61069361068e366004614889565b612524565b6040516103309594939291906150e9565b6103176106b236600461513e565b612610565b6105846106c536600461519a565b61283f565b7fd8e0666292c202b1ce6a8ff0dd638652e662402ac53fbf9bd9d3bcc39d5eb09761043d565b60005460ff16610478565b6103176107093660046151b3565b61284c565b6103176129ad565b61071e612b1b565b604051610330939291906151ce565b6103c661073b3660046149e8565b612bf6565b61031761074e366004614cf0565b612c0b565b61031761076136600461522f565b612c64565b610478610774366004614889565b612e70565b61031761078736600461519a565b612e7b565b61031761079a36600461522f565b613099565b6107a76130aa565b67ffffffffffffffff81166000908152600b60205260409020546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff168061081e576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61082882826130b2565b5050565b67ffffffffffffffff82166000908152600b6020526040812060030154601054911a908110610891576040517f45c108ce00000000000000000000000000000000000000000000000000000000815260ff821660048201526024015b60405180910390fd5b6010805460ff83169081106108a8576108a861524c565b90600052602060002090600891828204019190066004029054906101000a900463ffffffff1663ffffffff168263ffffffff161115610951576010805460ff83169081106108f8576108f861524c565b600091825260209091206008820401546040517f1d70f87a000000000000000000000000000000000000000000000000000000008152600790921660049081026101000a90910463ffffffff1690820152602401610888565b505050565b61095e6130aa565b61096782613420565b67ffffffffffffffff9091166000908152600b6020526040902060030155565b61098f613496565b600061099c84600061351c565b905060003073ffffffffffffffffffffffffffffffffffffffff8316036109c65750600754610a3a565b8173ffffffffffffffffffffffffffffffffffffffff16639883c10d6040518163ffffffff1660e01b81526004016020604051808303816000875af1158015610a13573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a37919061527b565b90505b8383604051610a4a929190615294565b60405180910390208103610a8a576040517fee03280800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b604051806060016040528082815260200185858080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250505090825250600954602090910190610aeb9061ffff16436152d3565b9052600086815260066020908152604090912082518155908201516001820190610b159082615387565b50604091820151600290910155517f0fcfd32a68209b42944376bc5f4bf72c41ba0f378cf60434f84390b82c9844cf90610b569087908490889088906154a1565b60405180910390a15050505050565b610b6d613496565b8151815181141580610b7f5750600881115b15610bb6576040517fee03280800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b81811015610ce2576000848281518110610bd557610bd561524c565b602002602001015190506000848381518110610bf357610bf361524c565b60200260200101519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161480610c5e575060008281526002602052604090205473ffffffffffffffffffffffffffffffffffffffff8281169116145b15610c95576040517fee03280800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b81610ccf576040517f4855c28800000000000000000000000000000000000000000000000000000000815260048101839052602401610888565b505080610cdb906154fc565b9050610bb9565b50600954600090610cf79061ffff16436152d3565b604080516060810182528681526020808201879052918101839052865192935091600391610d299183918901906146c9565b506020828101518051610d429260018501920190614710565b506040820151816002015590505060005b8451811015610e54577f72a33d2f293a0a70fad221bb610d3d6b52aed2d840adae1fa721071fbd290cfd858281518110610d8f57610d8f61524c565b602002602001015160026000888581518110610dad57610dad61524c565b6020026020010151815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16868481518110610df657610df661524c565b602002602001015185604051610e3c949392919093845273ffffffffffffffffffffffffffffffffffffffff928316602085015291166040830152606082015260800190565b60405180910390a1610e4d816154fc565b9050610d53565b5050505050565b610e63613496565b610e6b613613565b565b6000610eb68260008989898080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508b92508a91506136909050565b979650505050505050565b610ec96139c9565b610ed2826139d1565b610eda613a97565b67ffffffffffffffff82166000908152600b602052604090206001015473ffffffffffffffffffffffffffffffffffffffff8281166c0100000000000000000000000090920416146108285767ffffffffffffffff82166000818152600b602090815260409182902060010180546bffffffffffffffffffffffff166c0100000000000000000000000073ffffffffffffffffffffffffffffffffffffffff8716908102919091179091558251338152918201527f69436ea6df009049404f564eff6622cd00522b0bd6a89efd9e52a355c4a879be91015b60405180910390a25050565b67ffffffffffffffff81166000908152600b60205260408120600301545b92915050565b610fea6130aa565b806bffffffffffffffffffffffff166000036110205750306000908152600d60205260409020546bffffffffffffffffffffffff165b306000908152600d60205260409020546bffffffffffffffffffffffff8083169116101561107a576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b306000908152600d6020526040812080548392906110a79084906bffffffffffffffffffffffff16615534565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550806009600b8282829054906101000a90046bffffffffffffffffffffffff166110fe9190615534565b82546101009290920a6bffffffffffffffffffffffff818102199093169183160217909155600a546040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff868116600483015292851660248201529116915063a9059cbb906044016020604051808303816000875af115801561119d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906111c19190615560565b61082857600a546040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015260009173ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa158015611234573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611258919061527b565b6009546040517fa99da3020000000000000000000000000000000000000000000000000000000081526b0100000000000000000000009091046bffffffffffffffffffffffff166004820181905260248201839052919250604401610888565b6112c06139c9565b806bffffffffffffffffffffffff16600003611308576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b336000908152600d60205260409020546bffffffffffffffffffffffff80831691161015611362576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b336000908152600d60205260408120805483929061138f9084906bffffffffffffffffffffffff16615534565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550806009600b8282829054906101000a90046bffffffffffffffffffffffff166113e69190615534565b82546101009290920a6bffffffffffffffffffffffff818102199093169183160217909155600a546040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff868116600483015292851660248201529116915063a9059cbb906044016020604051808303816000875af1158015611485573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906114a99190615560565b610828576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006114eb838361351c565b9392505050565b6114fa6139c9565b60005b818110156109515760008383838181106115195761151961524c565b90506101600201803603810190611530919061557d565b905060008160e001519050600e60008281526020019081526020016000205482604051602001611560919061559a565b60405160208183030381529060405280519060200120146115ad576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8160c0015164ffffffffff164210156115f2576040517fbcc4005500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60208201516040517f85b214cf0000000000000000000000000000000000000000000000000000000081526004810183905273ffffffffffffffffffffffffffffffffffffffff8216906385b214cf906024016020604051808303816000875af1158015611664573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906116889190615560565b5060a0830151606084015167ffffffffffffffff166000908152600b6020526040812060010180549091906116cc9084906bffffffffffffffffffffffff16615534565b82546bffffffffffffffffffffffff9182166101009390930a92830291909202199091161790555060408084015173ffffffffffffffffffffffffffffffffffffffff166000908152600c6020908152828220606087015167ffffffffffffffff90811684529152919020805460019260099161175891859169010000000000000000009004166156c9565b825467ffffffffffffffff9182166101009390930a9283029190920219909116179055506000828152600e60205260408082208290555183917ff1ca1e9147be737b04a2b018a79405f687a97de8dd8a2559bbe62357343af41491a2505050806117c1906154fc565b90506114fd565b6117d06139c9565b6117d9826139d1565b6117e1613a97565b67ffffffffffffffff82166000908152600b60205260409020600201547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9c01611856576040517f05a48e0f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff81166000908152600c6020908152604080832067ffffffffffffffff8616845290915290205460ff161561189d575050565b73ffffffffffffffffffffffffffffffffffffffff81166000818152600c6020908152604080832067ffffffffffffffff871680855290835281842080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001908117909155600b84528285206002018054918201815585529383902090930180547fffffffffffffffffffffffff000000000000000000000000000000000000000016851790555192835290917f43dc749a04ac8fb825cbd514f7c0e13f13bc6f2ee66043b76629d51776cff8e09101610fb2565b60015473ffffffffffffffffffffffffffffffffffffffff1633146119f6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610888565b60008054336101008181027fffffffffffffffffffffff0000000000000000000000000000000000000000ff8416178455600180547fffffffffffffffffffffffff000000000000000000000000000000000000000016905560405173ffffffffffffffffffffffffffffffffffffffff919093041692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b611a9f6139c9565b611aa7613a97565b67ffffffffffffffff81166000908152600b60205260409020805460019091015473ffffffffffffffffffffffffffffffffffffffff6c010000000000000000000000009283900481169290910416338114611b47576040517fd084e97500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82166004820152602401610888565b67ffffffffffffffff83166000818152600b602090815260409182902080546c01000000000000000000000000339081026bffffffffffffffffffffffff928316178355600190920180549091169055825173ffffffffffffffffffffffffffffffffffffffff87168152918201527f6f1dc65165ffffedfd8e507b4a0f1fcfdada045ed11f6c26ba27cedfe87802f0910160405180910390a2505050565b611bee613496565b610e6b613b9e565b611bfe613496565b60085464010000000090047bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16431015611c60576040517fa93d035c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600854600980546201000090920461ffff167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909216919091179055565b611ca66139c9565b611caf826139d1565b611cb7613a97565b73ffffffffffffffffffffffffffffffffffffffff81166000908152600c6020908152604080832067ffffffffffffffff8087168552908352928190208151606081018352905460ff8116151580835261010082048616948301949094526901000000000000000000900490931690830152611d87576040517ff0019fe600000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8416600482015273ffffffffffffffffffffffffffffffffffffffff83166024820152604401610888565b806040015167ffffffffffffffff16816020015167ffffffffffffffff1614611ddc576040517fbcc4005500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff83166000908152600b6020908152604080832060020180548251818502810185019093528083529192909190830182828015611e5757602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311611e2c575b50505050509050600060018251611e6e91906156ea565b905060005b825181101561200a578473ffffffffffffffffffffffffffffffffffffffff16838281518110611ea557611ea561524c565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1603611ffa576000838381518110611edc57611edc61524c565b6020026020010151905080600b60008967ffffffffffffffff1667ffffffffffffffff1681526020019081526020016000206002018381548110611f2257611f2261524c565b600091825260208083209190910180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff949094169390931790925567ffffffffffffffff89168152600b90915260409020600201805480611f9c57611f9c6156fd565b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff00000000000000000000000000000000000000001690550190555061200a565b612003816154fc565b9050611e73565b5073ffffffffffffffffffffffffffffffffffffffff84166000818152600c6020908152604080832067ffffffffffffffff8a168085529083529281902080547fffffffffffffffffffffffffffffff00000000000000000000000000000000001690555192835290917f182bff9831466789164ca77075fffd84916d35a8180ba73c27e45634549b445b910160405180910390a25050505050565b6000806120b16139c9565b826020015173ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461211a576040517f8bec23e700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60e08301516000908152600e602052604090205461213b5760029150612385565b600e60008460e0015181526020019081526020016000205483604051602001612164919061559a565b60405160208183030381529060405280519060200120146121885760079150612385565b826101400151836080015163ffffffff166121a3919061572c565b64ffffffffff165a10156121ba5760039150612385565b606083015167ffffffffffffffff166000908152600b602052604090205460808401516bffffffffffffffffffffffff909116906121fd9063ffffffff16613bf9565b612207908861574a565b8451612214908890615772565b61221e9190615772565b6bffffffffffffffffffffffff16111561223b5760049150612385565b8260a001516bffffffffffffffffffffffff16612261846080015163ffffffff16613bf9565b61226b908861574a565b8451612278908890615772565b6122829190615772565b6bffffffffffffffffffffffff16111561229f5760059150612385565b60009150600e60008460e0015181526020019081526020016000206000905560006122d98460e001518a8a87608001518860400151613c9b565b80519091506122e95760016122ec565b60005b9250600061231385606001518660a00151876040015188600001518c87602001518d613e15565b9050846060015167ffffffffffffffff168560e001517f45bb48b6ec798595a260f114720360b95cc58c94c6ddd37a1acc3896ec94a23a836020015189888760000151612360578e612362565b8f5b8860400151604051612378959493929190615797565b60405180910390a3519150505b965096945050505050565b600061239a6139c9565b6123a2613a97565b600980546003906123c3906301000000900467ffffffffffffffff166157f5565b825467ffffffffffffffff8083166101009490940a93840293021916919091179091556040805160c0810182526000808252336020830152918101829052606081018290529192506080820190604051908082528060200260200182016040528015612439578160200160208202803683370190505b5081526000602091820181905267ffffffffffffffff84168152600b825260409081902083518484015173ffffffffffffffffffffffffffffffffffffffff9081166c010000000000000000000000009081026bffffffffffffffffffffffff93841617845593860151606087015190911690930292169190911760018201556080830151805191926124d492600285019290910190614710565b5060a0919091015160039091015560405133815267ffffffffffffffff8216907f464722b4166576d3dcbba877b999bc35cf911f4eaf434b7eba68fa113951d0bf9060200160405180910390a290565b600080600080606061253586613420565b67ffffffffffffffff86166000908152600b602090815260409182902080546001820154600290920180548551818602810186019096528086526bffffffffffffffffffffffff8084169b508416995073ffffffffffffffffffffffffffffffffffffffff6c0100000000000000000000000093849004811699509290930490911695509183018282801561260057602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116125d5575b5050505050905091939590929450565b6126186139c9565b600a5473ffffffffffffffffffffffffffffffffffffffff163314612669576040517f44b0e3c300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b602081146126a3576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006126b182840184614889565b67ffffffffffffffff81166000908152600b60205260409020549091506c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1661272a576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff81166000908152600b6020526040812080546bffffffffffffffffffffffff16918691906127618385615772565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550846009600b8282829054906101000a90046bffffffffffffffffffffffff166127b89190615772565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055508167ffffffffffffffff167fd39ec07f4e209f627a4c427971473820dc129761ba28de8906bd56f57101d4f882878461281f91906152d3565b6040805192835260208301919091520160405180910390a2505050505050565b6000610fdc82600061351c565b612854613496565b60095461ffff808316911603612896576040517fee03280800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000061ffff168161ffff1611156128f8576040517fe9a3062200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805160608101825260095461ffff9081168083529084166020830152909182019061292590436152d3565b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff90811690915281516008805460208501516040909501519093166401000000000263ffffffff61ffff95861662010000027fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000009095169590931694909417929092171691909117905550565b6129b5613496565b6005544310156129f1576040517fa93d035c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b600354811015612b1857600060036000018281548110612a1657612a1661524c565b6000918252602080832090910154808352600290915260408220546004805492945073ffffffffffffffffffffffffffffffffffffffff909116929185908110612a6257612a6261524c565b6000918252602080832091909101548583526002825260409283902080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92831690811790915583518781529186169282019290925291820181905291507ff8a6175bca1ba37d682089187edc5e20a859989727f10ca6bd9a5bc0de8caf949060600160405180910390a150505080612b11906154fc565b90506129f4565b50565b60006060806003600201546003600001600360010181805480602002602001604051908101604052809291908181526020018280548015612b7b57602002820191906000526020600020905b815481526020019060010190808311612b67575b5050505050915080805480602002602001604051908101604052809291908181526020018280548015612be457602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311612bb9575b50505050509050925092509250909192565b6060612c03848484614067565b949350505050565b612c136139c9565b612c1c826139d1565b612c24613a97565b612c2d82614134565b1561081e576040517fb42f66e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b612c6c6130aa565b600a546040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015260009173ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa158015612cdb573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612cff919061527b565b6009549091506b01000000000000000000000090046bffffffffffffffffffffffff1681811115612d66576040517fa99da3020000000000000000000000000000000000000000000000000000000081526004810182905260248101839052604401610888565b81811015610951576000612d7a82846156ea565b600a546040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff87811660048301526024820184905292935091169063a9059cbb906044016020604051808303816000875af1158015612df5573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612e199190615560565b506040805173ffffffffffffffffffffffffffffffffffffffff86168152602081018390527f59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b43660091015b60405180910390a150505050565b6000610fdc82614134565b612e83613496565b60006006600083815260200190815260200160002060405180606001604052908160008201548152602001600182018054612ebd906152e6565b80601f0160208091040260200160405190810160405280929190818152602001828054612ee9906152e6565b8015612f365780601f10612f0b57610100808354040283529160200191612f36565b820191906000526020600020905b815481529060010190602001808311612f1957829003601f168201915b5050505050815260200160028201548152505090508060400151431015612f89576040517fa93d035c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b81612fae57612f9b8160200151614283565b6020808201518051910120600755613056565b612fb982600061351c565b73ffffffffffffffffffffffffffffffffffffffff16638cc6acce82602001516040518263ffffffff1660e01b8152600401612ff59190614960565b600060405180830381600087803b15801561300f57600080fd5b505af1925050508015613020575060015b613056576040517ffe680b2600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b805160208201516040517fc07626096b49b0462a576c9a7878cf0675e784bb4832584c1289c11664e55f479261308d92869261581c565b60405180910390a15050565b6130a1613496565b612b188161437e565b610e6b613496565b67ffffffffffffffff82166000908152600b60209081526040808320815160c08101835281546bffffffffffffffffffffffff808216835273ffffffffffffffffffffffffffffffffffffffff6c010000000000000000000000009283900481168488015260018501549182168487015291900416606082015260028201805484518187028101870190955280855291949293608086019390929083018282801561319357602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311613168575b505050918352505060039190910154602090910152805190915060005b82608001515181101561325457600c6000846080015183815181106131d7576131d761524c565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040908101600090812067ffffffffffffffff89168252909252902080547fffffffffffffffffffffffffffffff000000000000000000000000000000000016905561324d816154fc565b90506131b0565b5067ffffffffffffffff84166000908152600b602052604081208181556001810182905590613286600283018261478a565b60038201600090555050806009600b8282829054906101000a90046bffffffffffffffffffffffff166132b99190615534565b82546101009290920a6bffffffffffffffffffffffff818102199093169183160217909155600a546040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff878116600483015292851660248201529116915063a9059cbb906044016020604051808303816000875af1158015613358573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061337c9190615560565b6133b2576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805173ffffffffffffffffffffffffffffffffffffffff851681526bffffffffffffffffffffffff8316602082015267ffffffffffffffff8616917fe8ed5b475a5b5987aa9165e8731bb78043f39eee32ec5a1169a89e27fcd49815910160405180910390a250505050565b67ffffffffffffffff81166000908152600b60205260409020546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff16612b18576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600054610100900473ffffffffffffffffffffffffffffffffffffffff163314610e6b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610888565b60008161355b5760008381526002602052604090205473ffffffffffffffffffffffffffffffffffffffff168015613555579050610fdc565b506135de565b60005b6003548110156135dc57600380548290811061357c5761357c61524c565b906000526020600020015484036135cc5760048054829081106135a1576135a161524c565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff169150610fdc9050565b6135d5816154fc565b905061355e565b505b6040517f80833e3300000000000000000000000000000000000000000000000000000000815260048101849052602401610888565b61361b614479565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa335b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390a1565b600061369a6139c9565b6136a385613420565b6136ad33866144e5565b6136b7858361082c565b60006136c3888861351c565b6040805160e08101825233815267ffffffffffffffff89166000818152600b602081815285832080546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff908116838801528688018e90526060870186905261ffff8d1660808801529484529190526003015460a084015263ffffffff881660c084015292517faecc12c500000000000000000000000000000000000000000000000000000000815293945091929184169163aecc12c59161378a91600401615844565b610160604051808303816000875af11580156137aa573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906137ce9190615919565b90506137df33888360a00151614581565b604051806101600160405280600f60000160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff1681526020013373ffffffffffffffffffffffffffffffffffffffff1681526020018867ffffffffffffffff1681526020018563ffffffff1681526020018260a001516bffffffffffffffffffffffff1681526020018260c0015164ffffffffff1681526020018260e00151815260200182610100015169ffffffffffffffffffff16815260200182610120015164ffffffffff16815260200182610140015164ffffffffff168152506040516020016138ec919061559a565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0018152828252805160209182012060e0850180516000908152600e8452848120929092555167ffffffffffffffff8c16808352600b9093529290205490928c92917f7c720ccd20069b8311a6be4ba1cf3294d09eb247aa5d73a8502054b6e68a2f54916139b1916c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1690339032908e908e908e906159ec565b60405180910390a460e0015198975050505050505050565b610e6b61465c565b67ffffffffffffffff81166000908152600b60205260409020546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1680613a48576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff821614610828576040517f5a68151d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b613ac27fd8e0666292c202b1ce6a8ff0dd638652e662402ac53fbf9bd9d3bcc39d5eb097600061351c565b604080516000815260208101918290527f6b14daf80000000000000000000000000000000000000000000000000000000090915273ffffffffffffffffffffffffffffffffffffffff9190911690636b14daf890613b2590339060248101615a50565b602060405180830381865afa158015613b42573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613b669190615560565b610e6b576040517f22906263000000000000000000000000000000000000000000000000000000008152336004820152602401610888565b613ba661465c565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a2586136663390565b60006bffffffffffffffffffffffff821115613c97576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f53616665436173743a2076616c756520646f65736e27742066697420696e203960448201527f36206269747300000000000000000000000000000000000000000000000000006064820152608401610888565b5090565b60408051606080820183526000808352602083015291810191909152600f546040516000916c01000000000000000000000000900460e01b90613ce690899089908990602401615a7f565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00181529181526020820180517fffffffff00000000000000000000000000000000000000000000000000000000949094167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff909416939093179092528151608480825260c0820190935290925060009182918291602082018180368337019050509050853b613d9857600080fd5b5a611388811015613da857600080fd5b611388810390508760408204820311613dc057600080fd5b60008086516020880160008b8df193505a900391503d6084811115613de3575060845b808252806000602084013e50604080516060810182529315158452602084019290925290820152979650505050505050565b60408051808201909152600080825260208201526000613e3484613bf9565b613e3e908661574a565b9050600081613e4d8886615772565b613e579190615772565b6040805180820182526bffffffffffffffffffffffff808616825280841660208084019190915267ffffffffffffffff8f166000908152600b909152928320805492975093945084939291613eae91859116615534565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055508184613ee89190615772565b336000908152600d602052604081208054909190613f159084906bffffffffffffffffffffffff16615772565b82546101009290920a6bffffffffffffffffffffffff818102199093169183160217909155306000908152600d6020526040812080548b94509092613f5c91859116615772565b82546101009290920a6bffffffffffffffffffffffff81810219909316918316021790915567ffffffffffffffff8c166000908152600b6020526040812060010180548d94509092613fb091859116615534565b82546bffffffffffffffffffffffff9182166101009390930a92830291909202199091161790555073ffffffffffffffffffffffffffffffffffffffff88166000908152600c6020908152604080832067ffffffffffffffff808f1685529252909120805460019260099161403491859169010000000000000000009004166156c9565b92506101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055505050979650505050505050565b6060600080808061407a86880188615ab4565b9350935093509350600061409389600187878787613690565b60408051602080825281830190925291925060208201818036833701905050955060005b6020811015614127578181602081106140d2576140d261524c565b1a60f81b8782815181106140e8576140e861524c565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350614120816154fc565b90506140b7565b5050505050509392505050565b67ffffffffffffffff81166000908152600b60209081526040808320600201805482518185028101850190935280835284938301828280156141ac57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311614181575b5050505050905060005b8151811015614279576000600c60008484815181106141d7576141d761524c565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040908101600090812067ffffffffffffffff808a168352908452908290208251606081018452905460ff811615158252610100810483169482018590526901000000000000000000900490911691810182905292501461426857506001949350505050565b50614272816154fc565b90506141b6565b5060009392505050565b60008060008380602001905181019061429c9190615b27565b604080516060810182526bffffffffffffffffffffffff85168082527fffffffff000000000000000000000000000000000000000000000000000000008516602080840191909152928201849052600f80547fffffffffffffffffffffffffffffffff00000000000000000000000000000000169091176c0100000000000000000000000060e087901c021781558351959850939650919450909261434791601091908601906147a8565b509050507fe6a1eda76d42a6d1d813f26765716562044db1e8bd8be7d088705e64afd301ca838383604051612e6293929190615c10565b3373ffffffffffffffffffffffffffffffffffffffff8216036143fd576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610888565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929361010090910416917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60005460ff16610e6b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f5061757361626c653a206e6f74207061757365640000000000000000000000006044820152606401610888565b73ffffffffffffffffffffffffffffffffffffffff82166000908152600c6020908152604080832067ffffffffffffffff8516845290915290205460ff16610828576040517ff0019fe600000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8216600482015273ffffffffffffffffffffffffffffffffffffffff83166024820152604401610888565b67ffffffffffffffff82166000908152600b6020526040812060010180548392906145bb9084906bffffffffffffffffffffffff16615772565b82546bffffffffffffffffffffffff91821661010093840a908102920219161790915573ffffffffffffffffffffffffffffffffffffffff85166000908152600c6020908152604080832067ffffffffffffffff80891685529252909120805460019450909284926146319284929004166156c9565b92506101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550505050565b60005460ff1615610e6b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f5061757361626c653a20706175736564000000000000000000000000000000006044820152606401610888565b828054828255906000526020600020908101928215614704579160200282015b828111156147045782518255916020019190600101906146e9565b50613c9792915061484e565b828054828255906000526020600020908101928215614704579160200282015b8281111561470457825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190614730565b5080546000825590600052602060002090810190612b18919061484e565b828054828255906000526020600020906007016008900481019282156147045791602002820160005b8382111561481557835183826101000a81548163ffffffff021916908363ffffffff16021790555092602001926004016020816003010492830192600103026147d1565b80156148455782816101000a81549063ffffffff0219169055600401602081600301049283019260010302614815565b5050613c979291505b5b80821115613c97576000815560010161484f565b67ffffffffffffffff81168114612b1857600080fd5b803561488481614863565b919050565b60006020828403121561489b57600080fd5b81356114eb81614863565b63ffffffff81168114612b1857600080fd5b8035614884816148a6565b600080604083850312156148d657600080fd5b82356148e181614863565b915060208301356148f1816148a6565b809150509250929050565b6000815180845260005b8181101561492257602081850181015186830182015201614906565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b6020815260006114eb60208301846148fc565b6000806040838503121561498657600080fd5b823561499181614863565b946020939093013593505050565b60008083601f8401126149b157600080fd5b50813567ffffffffffffffff8111156149c957600080fd5b6020830191508360208285010111156149e157600080fd5b9250929050565b6000806000604084860312156149fd57600080fd5b83359250602084013567ffffffffffffffff811115614a1b57600080fd5b614a278682870161499f565b9497909650939450505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610160810167ffffffffffffffff81118282101715614a8757614a87614a34565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715614ad457614ad4614a34565b604052919050565b600067ffffffffffffffff821115614af657614af6614a34565b5060051b60200190565b73ffffffffffffffffffffffffffffffffffffffff81168114612b1857600080fd5b803561488481614b00565b600082601f830112614b3e57600080fd5b81356020614b53614b4e83614adc565b614a8d565b82815260059290921b84018101918181019086841115614b7257600080fd5b8286015b84811015614b96578035614b8981614b00565b8352918301918301614b76565b509695505050505050565b60008060408385031215614bb457600080fd5b823567ffffffffffffffff80821115614bcc57600080fd5b818501915085601f830112614be057600080fd5b81356020614bf0614b4e83614adc565b82815260059290921b84018101918181019089841115614c0f57600080fd5b948201945b83861015614c2d57853582529482019490820190614c14565b96505086013592505080821115614c4357600080fd5b50614c5085828601614b2d565b9150509250929050565b803561ffff8116811461488457600080fd5b60008060008060008060a08789031215614c8557600080fd5b8635614c9081614863565b9550602087013567ffffffffffffffff811115614cac57600080fd5b614cb889828a0161499f565b9096509450614ccb905060408801614c5a565b92506060870135614cdb816148a6565b80925050608087013590509295509295509295565b60008060408385031215614d0357600080fd5b8235614d0e81614863565b915060208301356148f181614b00565b6bffffffffffffffffffffffff81168114612b1857600080fd5b803561488481614d1e565b60008060408385031215614d5657600080fd5b8235614d6181614b00565b915060208301356148f181614d1e565b60008060408385031215614d8457600080fd5b8235614d8f81614b00565b915060208301356148f181614863565b8015158114612b1857600080fd5b60008060408385031215614dc057600080fd5b8235915060208301356148f181614d9f565b60008060208385031215614de557600080fd5b823567ffffffffffffffff80821115614dfd57600080fd5b818501915085601f830112614e1157600080fd5b813581811115614e2057600080fd5b86602061016083028501011115614e3657600080fd5b60209290920196919550909350505050565b600082601f830112614e5957600080fd5b813567ffffffffffffffff811115614e7357614e73614a34565b614ea460207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601614a8d565b818152846020838601011115614eb957600080fd5b816020850160208301376000918101602001919091529392505050565b64ffffffffff81168114612b1857600080fd5b803561488481614ed6565b69ffffffffffffffffffff81168114612b1857600080fd5b803561488481614ef4565b60006101608284031215614f2a57600080fd5b614f32614a63565b9050614f3d82614d38565b8152614f4b60208301614b22565b6020820152614f5c60408301614b22565b6040820152614f6d60608301614879565b6060820152614f7e608083016148b8565b6080820152614f8f60a08301614d38565b60a0820152614fa060c08301614ee9565b60c082015260e082013560e0820152610100614fbd818401614f0c565b90820152610120614fcf838201614ee9565b90820152610140614fe1838201614ee9565b9082015292915050565b600080600080600080610200878903121561500557600080fd5b863567ffffffffffffffff8082111561501d57600080fd5b6150298a838b01614e48565b9750602089013591508082111561503f57600080fd5b5061504c89828a01614e48565b955050604087013561505d81614d1e565b9350606087013561506d81614d1e565b9250608087013561507d81614b00565b915061508c8860a08901614f17565b90509295509295509295565b600081518084526020808501945080840160005b838110156150de57815173ffffffffffffffffffffffffffffffffffffffff16875295820195908201906001016150ac565b509495945050505050565b6bffffffffffffffffffffffff86811682528516602082015273ffffffffffffffffffffffffffffffffffffffff84811660408301528316606082015260a060808201819052600090610eb690830184615098565b6000806000806060858703121561515457600080fd5b843561515f81614b00565b935060208501359250604085013567ffffffffffffffff81111561518257600080fd5b61518e8782880161499f565b95989497509550505050565b6000602082840312156151ac57600080fd5b5035919050565b6000602082840312156151c557600080fd5b6114eb82614c5a565b6000606082018583526020606081850152818651808452608086019150828801935060005b8181101561520f578451835293830193918301916001016151f3565b505084810360408601526152238187615098565b98975050505050505050565b60006020828403121561524157600080fd5b81356114eb81614b00565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60006020828403121561528d57600080fd5b5051919050565b8183823760009101908152919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b80820180821115610fdc57610fdc6152a4565b600181811c908216806152fa57607f821691505b602082108103615333577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f82111561095157600081815260208120601f850160051c810160208610156153605750805b601f850160051c820191505b8181101561537f5782815560010161536c565b505050505050565b815167ffffffffffffffff8111156153a1576153a1614a34565b6153b5816153af84546152e6565b84615339565b602080601f83116001811461540857600084156153d25750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b17855561537f565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b8281101561545557888601518255948401946001909101908401615436565b508582101561549157878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b84815283602082015260606040820152816060820152818360808301376000818301608090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01601019392505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361552d5761552d6152a4565b5060010190565b6bffffffffffffffffffffffff828116828216039080821115615559576155596152a4565b5092915050565b60006020828403121561557257600080fd5b81516114eb81614d9f565b6000610160828403121561559057600080fd5b6114eb8383614f17565b81516bffffffffffffffffffffffff168152610160810160208301516155d8602084018273ffffffffffffffffffffffffffffffffffffffff169052565b506040830151615600604084018273ffffffffffffffffffffffffffffffffffffffff169052565b50606083015161561c606084018267ffffffffffffffff169052565b506080830151615634608084018263ffffffff169052565b5060a083015161565460a08401826bffffffffffffffffffffffff169052565b5060c083015161566d60c084018264ffffffffff169052565b5060e083015160e0830152610100808401516156968285018269ffffffffffffffffffff169052565b50506101208381015164ffffffffff81168483015250506101408381015164ffffffffff8116848301525b505092915050565b67ffffffffffffffff818116838216019080821115615559576155596152a4565b81810381811115610fdc57610fdc6152a4565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b64ffffffffff818116838216019080821115615559576155596152a4565b6bffffffffffffffffffffffff8181168382160280821691908281146156c1576156c16152a4565b6bffffffffffffffffffffffff818116838216019080821115615559576155596152a4565b6bffffffffffffffffffffffff8616815273ffffffffffffffffffffffffffffffffffffffff8516602082015260ff8416604082015260a0606082015260006157e360a08301856148fc565b828103608084015261522381856148fc565b600067ffffffffffffffff808316818103615812576158126152a4565b6001019392505050565b83815282602082015260606040820152600061583b60608301846148fc565b95945050505050565b60208152600073ffffffffffffffffffffffffffffffffffffffff80845116602084015280602085015116604084015250604083015160e0606084015261588f6101008401826148fc565b905067ffffffffffffffff606085015116608084015261ffff60808501511660a084015260a084015160c084015263ffffffff60c08501511660e08401528091505092915050565b805161488481614d1e565b805161488481614b00565b805161488481614863565b8051614884816148a6565b805161488481614ed6565b805161488481614ef4565b6000610160828403121561592c57600080fd5b615934614a63565b61593d836158d7565b815261594b602084016158e2565b602082015261595c604084016158e2565b604082015261596d606084016158ed565b606082015261597e608084016158f8565b608082015261598f60a084016158d7565b60a08201526159a060c08401615903565b60c082015260e083015160e08201526101006159bd81850161590e565b908201526101206159cf848201615903565b908201526101406159e1848201615903565b908201529392505050565b600073ffffffffffffffffffffffffffffffffffffffff8089168352808816602084015280871660408401525060c06060830152615a2d60c08301866148fc565b905061ffff8416608083015263ffffffff831660a0830152979650505050505050565b73ffffffffffffffffffffffffffffffffffffffff83168152604060208201526000612c0360408301846148fc565b838152606060208201526000615a9860608301856148fc565b8281036040840152615aaa81856148fc565b9695505050505050565b60008060008060808587031215615aca57600080fd5b8435615ad581614863565b9350602085013567ffffffffffffffff811115615af157600080fd5b615afd87828801614e48565b935050615b0c60408601614c5a565b91506060850135615b1c816148a6565b939692955090935050565b600080600060608486031215615b3c57600080fd5b8351615b4781614d1e565b809350506020808501517fffffffff0000000000000000000000000000000000000000000000000000000081168114615b7f57600080fd5b604086015190935067ffffffffffffffff811115615b9c57600080fd5b8501601f81018713615bad57600080fd5b8051615bbb614b4e82614adc565b81815260059190911b82018301908381019089831115615bda57600080fd5b928401925b82841015615c01578351615bf2816148a6565b82529284019290840190615bdf565b80955050505050509250925092565b6000606082016bffffffffffffffffffffffff8616835260207fffffffff0000000000000000000000000000000000000000000000000000000086168185015260606040850152818551808452608086019150828701935060005b81811015615c8d57845163ffffffff1683529383019391830191600101615c6b565b50909897505050505050505056fea164736f6c6343000813000a",
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

func (_FunctionsRouter *FunctionsRouterCaller) GetConfig(opts *bind.CallOpts) (IFunctionsRouterConfig, error) {
	var out []interface{}
	err := _FunctionsRouter.contract.Call(opts, &out, "getConfig")

	if err != nil {
		return *new(IFunctionsRouterConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(IFunctionsRouterConfig)).(*IFunctionsRouterConfig)

	return out0, err

}

func (_FunctionsRouter *FunctionsRouterSession) GetConfig() (IFunctionsRouterConfig, error) {
	return _FunctionsRouter.Contract.GetConfig(&_FunctionsRouter.CallOpts)
}

func (_FunctionsRouter *FunctionsRouterCallerSession) GetConfig() (IFunctionsRouterConfig, error) {
	return _FunctionsRouter.Contract.GetConfig(&_FunctionsRouter.CallOpts)
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

func (_FunctionsRouter *FunctionsRouterTransactor) Fulfill(opts *bind.TransactOpts, response []byte, err []byte, juelsPerGas *big.Int, costWithoutFulfillment *big.Int, transmitter common.Address, commitment IFunctionsRequestCommitment) (*types.Transaction, error) {
	return _FunctionsRouter.contract.Transact(opts, "fulfill", response, err, juelsPerGas, costWithoutFulfillment, transmitter, commitment)
}

func (_FunctionsRouter *FunctionsRouterSession) Fulfill(response []byte, err []byte, juelsPerGas *big.Int, costWithoutFulfillment *big.Int, transmitter common.Address, commitment IFunctionsRequestCommitment) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.Fulfill(&_FunctionsRouter.TransactOpts, response, err, juelsPerGas, costWithoutFulfillment, transmitter, commitment)
}

func (_FunctionsRouter *FunctionsRouterTransactorSession) Fulfill(response []byte, err []byte, juelsPerGas *big.Int, costWithoutFulfillment *big.Int, transmitter common.Address, commitment IFunctionsRequestCommitment) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.Fulfill(&_FunctionsRouter.TransactOpts, response, err, juelsPerGas, costWithoutFulfillment, transmitter, commitment)
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

func (_FunctionsRouter *FunctionsRouterTransactor) TimeoutRequests(opts *bind.TransactOpts, requestsToTimeoutByCommitment []IFunctionsRequestCommitment) (*types.Transaction, error) {
	return _FunctionsRouter.contract.Transact(opts, "timeoutRequests", requestsToTimeoutByCommitment)
}

func (_FunctionsRouter *FunctionsRouterSession) TimeoutRequests(requestsToTimeoutByCommitment []IFunctionsRequestCommitment) (*types.Transaction, error) {
	return _FunctionsRouter.Contract.TimeoutRequests(&_FunctionsRouter.TransactOpts, requestsToTimeoutByCommitment)
}

func (_FunctionsRouter *FunctionsRouterTransactorSession) TimeoutRequests(requestsToTimeoutByCommitment []IFunctionsRequestCommitment) (*types.Transaction, error) {
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

	GetAllowListId(opts *bind.CallOpts) ([32]byte, error)

	GetConfig(opts *bind.CallOpts) (IFunctionsRouterConfig, error)

	GetConfigHash(opts *bind.CallOpts) ([32]byte, error)

	GetConsumer(opts *bind.CallOpts, client common.Address, subscriptionId uint64) (GetConsumer,

		error)

	GetContractById(opts *bind.CallOpts, id [32]byte, useProposed bool) (common.Address, error)

	GetContractById0(opts *bind.CallOpts, id [32]byte) (common.Address, error)

	GetFlags(opts *bind.CallOpts, subscriptionId uint64) ([32]byte, error)

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

	Fulfill(opts *bind.TransactOpts, response []byte, err []byte, juelsPerGas *big.Int, costWithoutFulfillment *big.Int, transmitter common.Address, commitment IFunctionsRequestCommitment) (*types.Transaction, error)

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

	TimeoutRequests(opts *bind.TransactOpts, requestsToTimeoutByCommitment []IFunctionsRequestCommitment) (*types.Transaction, error)

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
