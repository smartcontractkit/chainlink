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

type IFunctionsRouterConfig struct {
	MaxConsumers                    uint16
	AdminFee                        *big.Int
	HandleOracleFulfillmentSelector [4]byte
	MaxCallbackGasLimits            []uint32
}

var FunctionsRouterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"timelockBlocks\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"maximumTimelockBlocks\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"linkToken\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"internalBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"externalBalance\",\"type\":\"uint256\"}],\"name\":\"BalanceInvariantViolated\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConsumerRequestsInFlight\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"limit\",\"type\":\"uint32\"}],\"name\":\"GasLimitTooBig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"IdentifierIsReserved\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidConfigData\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"InvalidConsumer\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"value\",\"type\":\"uint8\"}],\"name\":\"InvalidGasFlagValue\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidProposal\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSubscription\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposedOwner\",\"type\":\"address\"}],\"name\":\"MustBeRequestedOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeSubscriptionOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByRoute\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableFromCoordinator\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableFromLink\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingRequestExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ProposedTimelockAboveMaximum\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Reentrant\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"RouteNotFound\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"SenderMustAcceptTermsOfService\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TimelockInEffect\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyConsumers\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"adminFee\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"bytes4\",\"name\":\"handleOracleFulfillmentSelector\",\"type\":\"bytes4\"},{\"indexed\":false,\"internalType\":\"uint32[]\",\"name\":\"maxCallbackGasLimits\",\"type\":\"uint32[]\"}],\"name\":\"ConfigChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"fromHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"toBytes\",\"type\":\"bytes\"}],\"name\":\"ConfigProposed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"fromHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"toBytes\",\"type\":\"bytes\"}],\"name\":\"ConfigUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"proposedContractSetId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"proposedContractSetFromAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"proposedContractSetToAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timelockEndBlock\",\"type\":\"uint256\"}],\"name\":\"ContractProposed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"proposedContractSetId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"proposedContractSetFromAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"proposedContractSetToAddress\",\"type\":\"address\"}],\"name\":\"ContractUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"resultCode\",\"type\":\"uint8\"}],\"name\":\"RequestNotProcessed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalCostJuels\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"resultCode\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"response\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"}],\"name\":\"RequestProcessed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"donId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"subscriptionOwner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requestingContract\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requestInitiator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"dataVersion\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"}],\"name\":\"RequestStart\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"RequestTimedOut\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"fundsRecipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fundsAmount\",\"type\":\"uint256\"}],\"name\":\"SubscriptionCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"SubscriptionCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"from\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"to\",\"type\":\"uint16\"}],\"name\":\"TimeLockProposed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"from\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"to\",\"type\":\"uint16\"}],\"name\":\"TimeLockUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"MAX_CALLBACK_RETURN_BYTES\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"acceptSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"addConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"cancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createSubscription\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"response\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"juelsPerGas\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"costWithoutFulfillment\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint96\",\"name\":\"adminFee\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"estimatedTotalCostJuels\",\"type\":\"uint96\"},{\"internalType\":\"uint40\",\"name\":\"timeoutTimestamp\",\"type\":\"uint40\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint80\",\"name\":\"donFee\",\"type\":\"uint80\"},{\"internalType\":\"uint40\",\"name\":\"gasOverheadBeforeCallback\",\"type\":\"uint40\"},{\"internalType\":\"uint40\",\"name\":\"gasOverheadAfterCallback\",\"type\":\"uint40\"}],\"internalType\":\"structIFunctionsRequest.Commitment\",\"name\":\"commitment\",\"type\":\"tuple\"}],\"name\":\"fulfill\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"resultCode\",\"type\":\"uint8\"},{\"internalType\":\"uint96\",\"name\":\"callbackGasCostJuels\",\"type\":\"uint96\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowListId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint16\",\"name\":\"maxConsumers\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"adminFee\",\"type\":\"uint96\"},{\"internalType\":\"bytes4\",\"name\":\"handleOracleFulfillmentSelector\",\"type\":\"bytes4\"},{\"internalType\":\"uint32[]\",\"name\":\"maxCallbackGasLimits\",\"type\":\"uint32[]\"}],\"internalType\":\"structIFunctionsRouter.Config\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfigHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"config\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"getConsumer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"initiatedRequests\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"completedRequests\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"useProposed\",\"type\":\"bool\"}],\"name\":\"getContractById\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"routeDestination\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"getContractById\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"routeDestination\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"getFlags\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getProposedContractSet\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"},{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"getSubscription\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"blockedBalance\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"requestedOwner\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSubscriptionCount\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTotalBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isPaused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"}],\"name\":\"isValidCallbackGasLimit\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"oracleWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"ownerCancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"ownerWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"pendingRequestExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"name\":\"proposeConfigUpdate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"proposedContractSetIds\",\"type\":\"bytes32[]\"},{\"internalType\":\"address[]\",\"name\":\"proposedContractSetAddresses\",\"type\":\"address[]\"}],\"name\":\"proposeContractsUpdate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"proposeSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"blocks\",\"type\":\"uint16\"}],\"name\":\"proposeTimelockBlocks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"removeConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint16\",\"name\":\"dataVersion\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"donId\",\"type\":\"bytes32\"}],\"name\":\"sendRequest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"flags\",\"type\":\"bytes32\"}],\"name\":\"setFlags\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint96\",\"name\":\"adminFee\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"estimatedTotalCostJuels\",\"type\":\"uint96\"},{\"internalType\":\"uint40\",\"name\":\"timeoutTimestamp\",\"type\":\"uint40\"},{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint80\",\"name\":\"donFee\",\"type\":\"uint80\"},{\"internalType\":\"uint40\",\"name\":\"gasOverheadBeforeCallback\",\"type\":\"uint40\"},{\"internalType\":\"uint40\",\"name\":\"gasOverheadAfterCallback\",\"type\":\"uint40\"}],\"internalType\":\"structIFunctionsRequest.Commitment[]\",\"name\":\"requestsToTimeoutByCommitment\",\"type\":\"tuple[]\"}],\"name\":\"timeoutRequests\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"updateConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"updateContracts\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"updateTimelockBlocks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"validateProposedContracts\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b506040516200666b3803806200666b83398101604081905262000034916200044b565b6000805460ff191681558290339086908690859084908190816200009f5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0380851661010002610100600160a81b031990921691909117909155811615620000d957620000d98162000178565b50506009805461ffff80871661ffff19909216919091179091558316608052506000805260026020527fac33ff75c19e70fe83507db0d683fd3465c996598dc972688b7ace676c89077b80546001600160a01b031916301790556200013e8162000229565b80516020909101206007555050600a80546001600160a01b0319166001600160a01b03939093169290921790915550620006d79350505050565b336001600160a01b03821603620001d25760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000096565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929361010090910416917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000806000808480602001905181019062000245919062000548565b6040805160808101825261ffff86168082526001600160601b03861660208084018290526001600160e01b031987169484019490945260608301859052600f80546001600160701b031916909217620100009091021763ffffffff60701b1916600160701b60e087901c021781558351969a509498509296509094509092620002d591601091908601906200031d565b509050507fe6a1eda76d42a6d1d813f26765716562044db1e8bd8be7d088705e64afd301ca8383836040516200030e9392919062000666565b60405180910390a15050505050565b82805482825590600052602060002090600701600890048101928215620003c15791602002820160005b838211156200038d57835183826101000a81548163ffffffff021916908363ffffffff160217905550926020019260040160208160030104928301926001030262000347565b8015620003bf5782816101000a81549063ffffffff02191690556004016020816003010492830192600103026200038d565b505b50620003cf929150620003d3565b5090565b5b80821115620003cf5760008155600101620003d4565b805161ffff81168114620003fd57600080fd5b919050565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f191681016001600160401b038111828210171562000443576200044362000402565b604052919050565b600080600080608085870312156200046257600080fd5b6200046d85620003ea565b935060206200047e818701620003ea565b60408701519094506001600160a01b03811681146200049c57600080fd5b60608701519093506001600160401b0380821115620004ba57600080fd5b818801915088601f830112620004cf57600080fd5b815181811115620004e457620004e462000402565b620004f8601f8201601f1916850162000418565b915080825289848285010111156200050f57600080fd5b60005b818110156200052f57838101850151838201860152840162000512565b5060008482840101525080935050505092959194509250565b600080600080608085870312156200055f57600080fd5b6200056a85620003ea565b602086810151919550906001600160601b03811681146200058a57600080fd5b60408701519094506001600160e01b031981168114620005a957600080fd5b60608701519093506001600160401b0380821115620005c757600080fd5b818801915088601f830112620005dc57600080fd5b815181811115620005f157620005f162000402565b8060051b91506200060484830162000418565b818152918301840191848101908b8411156200061f57600080fd5b938501935b8385101562000656578451925063ffffffff83168314620006455760008081fd5b828252938501939085019062000624565b989b979a50959850505050505050565b6001600160601b03841681526001600160e01b031983166020808301919091526060604083018190528351908301819052600091848101916080850190845b81811015620006c957845163ffffffff1683529383019391830191600101620006a5565b509098975050505050505050565b608051615f78620006f3600039600061294f0152615f786000f3fe608060405234801561001057600080fd5b50600436106102f45760003560e01c80638456cb5911610191578063b187bd26116100e3578063c4614d6111610097578063e82ad7d411610071578063e82ad7d414610752578063eb523d6c14610765578063f2fde38b1461077857600080fd5b8063c4614d6114610719578063d7ae1d301461072c578063e72f6e301461073f57600080fd5b8063b734c0f4116100c8578063b734c0f4146106e5578063badc3eb6146106ed578063c3f909d41461070457600080fd5b8063b187bd26146106c7578063b5643858146106d257600080fd5b8063a1d4c82811610145578063a4c0ed361161011f578063a4c0ed361461067b578063a9c9a9181461068e578063aab396bd146106a157600080fd5b8063a1d4c82814610616578063a21a23e41461064f578063a47c76961461065757600080fd5b80638fde5317116101765780638fde5317146105f35780639883c10d146105fb5780639f87fad71461060357600080fd5b80638456cb59146105c85780638da5cb5b146105d057600080fd5b806355fedefa1161024a578063674603d0116101fe5780637341c10c116101d85780637341c10c1461059a57806379ba5097146105ad57806382359740146105b557600080fd5b8063674603d0146104ba5780636a6df79b1461054f5780636e3b33231461058757600080fd5b80635ed6dfba1161022f5780635ed6dfba1461046657806366316d8d14610479578063664199701461048c57600080fd5b806355fedefa1461043c5780635c975abb1461044f57600080fd5b80631ded3b36116102ac5780633f4ba83a116102865780633f4ba83a14610400578063461d2762146104085780634b8832d31461042957600080fd5b80631ded3b36146103c7578063385de9ae146103da5780633e871e4d146103ed57600080fd5b806310fc49c1116102dd57806310fc49c11461032e57806312b5834914610341578063181f5a771461037e57600080fd5b806302bcc5b6146102f95780630c5d49cb1461030e575b600080fd5b61030c610307366004614ab1565b61078b565b005b610316608481565b60405161ffff90911681526020015b60405180910390f35b61030c61033c366004614aeb565b610818565b6009546b01000000000000000000000090046bffffffffffffffffffffffff166040516bffffffffffffffffffffffff9091168152602001610325565b6103ba6040518060400160405280601781526020017f46756e6374696f6e7320526f757465722076312e302e3000000000000000000081525081565b6040516103259190614b88565b61030c6103d5366004614b9b565b610942565b61030c6103e8366004614c10565b610973565b61030c6103fb366004614dc9565b610b51565b61030c610e47565b61041b610416366004614e92565b610e59565b604051908152602001610325565b61030c610437366004614f17565b610ead565b61041b61044a366004614ab1565b610faa565b60005460ff165b6040519015158152602001610325565b61030c610474366004614f6a565b610fce565b61030c610487366004614f6a565b6112a4565b6009546301000000900467ffffffffffffffff165b60405167ffffffffffffffff9091168152602001610325565b6105276104c8366004614f98565b73ffffffffffffffffffffffffffffffffffffffff919091166000908152600c6020908152604080832067ffffffffffffffff948516845290915290205460ff8116926101008204831692690100000000000000000090920490911690565b60408051931515845267ffffffffffffffff9283166020850152911690820152606001610325565b61056261055d366004614fd4565b6114cb565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610325565b61030c610595366004614ff9565b6114de565b61030c6105a8366004614f17565b6117b4565b61030c61194a565b61030c6105c3366004614ab1565b611a6c565b61030c611bbb565b600054610100900473ffffffffffffffffffffffffffffffffffffffff16610562565b61030c611bcb565b60075461041b565b61030c610611366004614f17565b611c73565b610629610624366004615212565b61207b565b6040805160ff90931683526bffffffffffffffffffffffff909116602083015201610325565b6104a1612447565b61066a610665366004614ab1565b6125db565b604051610325959493929190615310565b61030c610689366004615365565b6126c7565b61056261069c3660046153c1565b6128f6565b7fd8e0666292c202b1ce6a8ff0dd638652e662402ac53fbf9bd9d3bcc39d5eb09761041b565b60005460ff16610456565b61030c6106e03660046153da565b612903565b61030c612a64565b6106f5612bd2565b604051610325939291906153f7565b61070c612cad565b604051610325919061548e565b6103ba610727366004614c10565b612dc4565b61030c61073a366004614f17565b612dd9565b61030c61074d3660046154fb565b612e32565b610456610760366004614ab1565b61303d565b61030c6107733660046153c1565b613048565b61030c6107863660046154fb565b613266565b610793613277565b67ffffffffffffffff81166000908152600b60205260409020546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff168061080a576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610814828261327f565b5050565b67ffffffffffffffff82166000908152600b6020526040812060030154601054911a90811061087d576040517f45c108ce00000000000000000000000000000000000000000000000000000000815260ff821660048201526024015b60405180910390fd5b6010805460ff831690811061089457610894615518565b90600052602060002090600891828204019190066004029054906101000a900463ffffffff1663ffffffff168263ffffffff16111561093d576010805460ff83169081106108e4576108e4615518565b600091825260209091206008820401546040517f1d70f87a000000000000000000000000000000000000000000000000000000008152600790921660049081026101000a90910463ffffffff1690820152602401610874565b505050565b61094a613277565b610953826135ed565b67ffffffffffffffff9091166000908152600b6020526040902060030155565b61097b613663565b60006109888460006136e9565b905060003073ffffffffffffffffffffffffffffffffffffffff8316036109b25750600754610a26565b8173ffffffffffffffffffffffffffffffffffffffff16639883c10d6040518163ffffffff1660e01b81526004016020604051808303816000875af11580156109ff573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a239190615547565b90505b8383604051610a36929190615560565b60405180910390208103610a76576040517fee03280800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b604051806060016040528082815260200185858080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250505090825250600954602090910190610ad79061ffff164361559f565b9052600086815260066020908152604090912082518155908201516001820190610b019082615653565b50604091820151600290910155517f0fcfd32a68209b42944376bc5f4bf72c41ba0f378cf60434f84390b82c9844cf90610b4290879084908890889061576d565b60405180910390a15050505050565b610b59613663565b8151815181141580610b6b5750600881115b15610ba2576040517fee03280800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b81811015610cce576000848281518110610bc157610bc1615518565b602002602001015190506000848381518110610bdf57610bdf615518565b60200260200101519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161480610c4a575060008281526002602052604090205473ffffffffffffffffffffffffffffffffffffffff8281169116145b15610c81576040517fee03280800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b81610cbb576040517f4855c28800000000000000000000000000000000000000000000000000000000815260048101839052602401610874565b505080610cc7906157c8565b9050610ba5565b50600954600090610ce39061ffff164361559f565b604080516060810182528681526020808201879052918101839052865192935091600391610d159183918901906148f1565b506020828101518051610d2e9260018501920190614938565b506040820151816002015590505060005b8451811015610e40577f72a33d2f293a0a70fad221bb610d3d6b52aed2d840adae1fa721071fbd290cfd858281518110610d7b57610d7b615518565b602002602001015160026000888581518110610d9957610d99615518565b6020026020010151815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16868481518110610de257610de2615518565b602002602001015185604051610e28949392919093845273ffffffffffffffffffffffffffffffffffffffff928316602085015291166040830152606082015260800190565b60405180910390a1610e39816157c8565b9050610d3f565b5050505050565b610e4f613663565b610e576137e0565b565b6000610ea28260008989898080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508b92508a915061385d9050565b979650505050505050565b610eb5613bb3565b610ebe82613bbb565b610ec6613c81565b67ffffffffffffffff82166000908152600b602052604090206001015473ffffffffffffffffffffffffffffffffffffffff8281166c0100000000000000000000000090920416146108145767ffffffffffffffff82166000818152600b602090815260409182902060010180546bffffffffffffffffffffffff166c0100000000000000000000000073ffffffffffffffffffffffffffffffffffffffff8716908102919091179091558251338152918201527f69436ea6df009049404f564eff6622cd00522b0bd6a89efd9e52a355c4a879be91015b60405180910390a25050565b67ffffffffffffffff81166000908152600b60205260408120600301545b92915050565b610fd6613277565b806bffffffffffffffffffffffff1660000361100c5750306000908152600d60205260409020546bffffffffffffffffffffffff165b306000908152600d60205260409020546bffffffffffffffffffffffff80831691161015611066576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b306000908152600d6020526040812080548392906110939084906bffffffffffffffffffffffff16615800565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550806009600b8282829054906101000a90046bffffffffffffffffffffffff166110ea9190615800565b82546101009290920a6bffffffffffffffffffffffff818102199093169183160217909155600a546040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff868116600483015292851660248201529116915063a9059cbb906044016020604051808303816000875af1158015611189573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906111ad919061582c565b61081457600a546040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015260009173ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa158015611220573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112449190615547565b6009546040517fa99da3020000000000000000000000000000000000000000000000000000000081526b0100000000000000000000009091046bffffffffffffffffffffffff166004820181905260248201839052919250604401610874565b6112ac613bb3565b806bffffffffffffffffffffffff166000036112f4576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b336000908152600d60205260409020546bffffffffffffffffffffffff8083169116101561134e576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b336000908152600d60205260408120805483929061137b9084906bffffffffffffffffffffffff16615800565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550806009600b8282829054906101000a90046bffffffffffffffffffffffff166113d29190615800565b82546101009290920a6bffffffffffffffffffffffff818102199093169183160217909155600a546040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff868116600483015292851660248201529116915063a9059cbb906044016020604051808303816000875af1158015611471573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611495919061582c565b610814576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006114d783836136e9565b9392505050565b6114e6613bb3565b60005b8181101561093d57600083838381811061150557611505615518565b9050610160020180360381019061151c9190615849565b905060008160e001519050600e6000828152602001908152602001600020548260405160200161154c9190615866565b6040516020818303038152906040528051906020012014611599576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8160c0015164ffffffffff164210156115de576040517fbcc4005500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60208201516040517f85b214cf0000000000000000000000000000000000000000000000000000000081526004810183905273ffffffffffffffffffffffffffffffffffffffff8216906385b214cf906024016020604051808303816000875af1158015611650573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611674919061582c565b5060a0830151606084015167ffffffffffffffff166000908152600b6020526040812060010180549091906116b89084906bffffffffffffffffffffffff16615800565b82546bffffffffffffffffffffffff9182166101009390930a92830291909202199091161790555060408084015173ffffffffffffffffffffffffffffffffffffffff166000908152600c6020908152828220606087015167ffffffffffffffff9081168452915291902080546001926009916117449185916901000000000000000000900416615995565b825467ffffffffffffffff9182166101009390930a9283029190920219909116179055506000828152600e60205260408082208290555183917ff1ca1e9147be737b04a2b018a79405f687a97de8dd8a2559bbe62357343af41491a2505050806117ad906157c8565b90506114e9565b6117bc613bb3565b6117c582613bbb565b6117cd613c81565b600f5467ffffffffffffffff83166000908152600b602052604090206002015461ffff909116900361182b576040517f05a48e0f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff81166000908152600c6020908152604080832067ffffffffffffffff8616845290915290205460ff1615611872575050565b73ffffffffffffffffffffffffffffffffffffffff81166000818152600c6020908152604080832067ffffffffffffffff871680855290835281842080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001908117909155600b84528285206002018054918201815585529383902090930180547fffffffffffffffffffffffff000000000000000000000000000000000000000016851790555192835290917f43dc749a04ac8fb825cbd514f7c0e13f13bc6f2ee66043b76629d51776cff8e09101610f9e565b60015473ffffffffffffffffffffffffffffffffffffffff1633146119cb576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610874565b60008054336101008181027fffffffffffffffffffffff0000000000000000000000000000000000000000ff8416178455600180547fffffffffffffffffffffffff000000000000000000000000000000000000000016905560405173ffffffffffffffffffffffffffffffffffffffff919093041692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b611a74613bb3565b611a7c613c81565b67ffffffffffffffff81166000908152600b60205260409020805460019091015473ffffffffffffffffffffffffffffffffffffffff6c010000000000000000000000009283900481169290910416338114611b1c576040517fd084e97500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82166004820152602401610874565b67ffffffffffffffff83166000818152600b602090815260409182902080546c01000000000000000000000000339081026bffffffffffffffffffffffff928316178355600190920180549091169055825173ffffffffffffffffffffffffffffffffffffffff87168152918201527f6f1dc65165ffffedfd8e507b4a0f1fcfdada045ed11f6c26ba27cedfe87802f0910160405180910390a2505050565b611bc3613663565b610e57613d88565b611bd3613663565b60085464010000000090047bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16431015611c35576040517fa93d035c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600854600980546201000090920461ffff167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909216919091179055565b611c7b613bb3565b611c8482613bbb565b611c8c613c81565b73ffffffffffffffffffffffffffffffffffffffff81166000908152600c6020908152604080832067ffffffffffffffff8087168552908352928190208151606081018352905460ff8116151580835261010082048616948301949094526901000000000000000000900490931690830152611d5c576040517ff0019fe600000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8416600482015273ffffffffffffffffffffffffffffffffffffffff83166024820152604401610874565b806040015167ffffffffffffffff16816020015167ffffffffffffffff1614611db1576040517fbcc4005500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff83166000908152600b6020908152604080832060020180548251818502810185019093528083529192909190830182828015611e2c57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311611e01575b50505050509050600060018251611e4391906159b6565b905060005b8251811015611fdf578473ffffffffffffffffffffffffffffffffffffffff16838281518110611e7a57611e7a615518565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1603611fcf576000838381518110611eb157611eb1615518565b6020026020010151905080600b60008967ffffffffffffffff1667ffffffffffffffff1681526020019081526020016000206002018381548110611ef757611ef7615518565b600091825260208083209190910180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff949094169390931790925567ffffffffffffffff89168152600b90915260409020600201805480611f7157611f716159c9565b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff000000000000000000000000000000000000000016905501905550611fdf565b611fd8816157c8565b9050611e48565b5073ffffffffffffffffffffffffffffffffffffffff84166000818152600c6020908152604080832067ffffffffffffffff8a168085529083529281902080547fffffffffffffffffffffffffffffff00000000000000000000000000000000001690555192835290917f182bff9831466789164ca77075fffd84916d35a8180ba73c27e45634549b445b910160405180910390a25050505050565b600080612086613bb3565b826020015173ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146120ef576040517f8bec23e700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60e08301516000908152600e60205260409020546121845760025b91508260e001517f1a90e9a50793db2e394cf581e7c522e10c358a81e70acf6b5a0edd620c08dee1846020015186856040516121779392919073ffffffffffffffffffffffffffffffffffffffff938416815291909216602082015260ff91909116604082015260600190565b60405180910390a261243c565b600e60008460e00151815260200190815260200160002054836040516020016121ad9190615866565b60405160208183030381529060405280519060200120146121cf57600661210a565b826101400151836080015163ffffffff166121ea91906159f8565b64ffffffffff165a10156121ff57600461210a565b606083015167ffffffffffffffff166000908152600b602052604090205460808401516bffffffffffffffffffffffff909116906122429063ffffffff16613de3565b61224c9088615a16565b8451612259908890615a3e565b6122639190615a3e565b6bffffffffffffffffffffffff1611156122f45760e0830180516000908152600e6020908152604080832092909255915182860151825173ffffffffffffffffffffffffffffffffffffffff9182168152908816938101939093526005918301829052909350907f1a90e9a50793db2e394cf581e7c522e10c358a81e70acf6b5a0edd620c08dee190606001612177565b8260a001516bffffffffffffffffffffffff1661231a846080015163ffffffff16613de3565b6123249088615a16565b8451612331908890615a3e565b61233b9190615a3e565b6bffffffffffffffffffffffff16111561235657600361210a565b60009150600e60008460e0015181526020019081526020016000206000905560006123908460e001518a8a87608001518860400151613e85565b80519091506123a05760016123a3565b60005b925060006123ca85606001518660a00151876040015188600001518c87602001518d614001565b9050846060015167ffffffffffffffff168560e001517f47ffcaa55fde21cc7135c65541826c9e65dda59c29dc109aae964989e8fc664b836020015189888760000151612417578e612419565b8f5b886040015160405161242f959493929190615a63565b60405180910390a3519150505b965096945050505050565b6000612451613bb3565b612459613c81565b6009805460039061247a906301000000900467ffffffffffffffff16615ac1565b825467ffffffffffffffff8083166101009490940a93840293021916919091179091556040805160c08101825260008082523360208301529181018290526060810182905291925060808201906040519080825280602002602001820160405280156124f0578160200160208202803683370190505b5081526000602091820181905267ffffffffffffffff84168152600b825260409081902083518484015173ffffffffffffffffffffffffffffffffffffffff9081166c010000000000000000000000009081026bffffffffffffffffffffffff938416178455938601516060870151909116909302921691909117600182015560808301518051919261258b92600285019290910190614938565b5060a0919091015160039091015560405133815267ffffffffffffffff8216907f464722b4166576d3dcbba877b999bc35cf911f4eaf434b7eba68fa113951d0bf9060200160405180910390a290565b60008060008060606125ec866135ed565b67ffffffffffffffff86166000908152600b602090815260409182902080546001820154600290920180548551818602810186019096528086526bffffffffffffffffffffffff8084169b508416995073ffffffffffffffffffffffffffffffffffffffff6c010000000000000000000000009384900481169950929093049091169550918301828280156126b757602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff16815260019091019060200180831161268c575b5050505050905091939590929450565b6126cf613bb3565b600a5473ffffffffffffffffffffffffffffffffffffffff163314612720576040517f44b0e3c300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6020811461275a576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600061276882840184614ab1565b67ffffffffffffffff81166000908152600b60205260409020549091506c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff166127e1576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff81166000908152600b6020526040812080546bffffffffffffffffffffffff16918691906128188385615a3e565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550846009600b8282829054906101000a90046bffffffffffffffffffffffff1661286f9190615a3e565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055508167ffffffffffffffff167fd39ec07f4e209f627a4c427971473820dc129761ba28de8906bd56f57101d4f88287846128d6919061559f565b6040805192835260208301919091520160405180910390a2505050505050565b6000610fc88260006136e9565b61290b613663565b60095461ffff80831691160361294d576040517fee03280800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000061ffff168161ffff1611156129af576040517fe9a3062200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805160608101825260095461ffff908116808352908416602083015290918201906129dc904361559f565b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff90811690915281516008805460208501516040909501519093166401000000000263ffffffff61ffff95861662010000027fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000009095169590931694909417929092171691909117905550565b612a6c613663565b600554431015612aa8576040517fa93d035c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b600354811015612bcf57600060036000018281548110612acd57612acd615518565b6000918252602080832090910154808352600290915260408220546004805492945073ffffffffffffffffffffffffffffffffffffffff909116929185908110612b1957612b19615518565b6000918252602080832091909101548583526002825260409283902080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92831690811790915583518781529186169282019290925291820181905291507ff8a6175bca1ba37d682089187edc5e20a859989727f10ca6bd9a5bc0de8caf949060600160405180910390a150505080612bc8906157c8565b9050612aab565b50565b60006060806003600201546003600001600360010181805480602002602001604051908101604052809291908181526020018280548015612c3257602002820191906000526020600020905b815481526020019060010190808311612c1e575b5050505050915080805480602002602001604051908101604052809291908181526020018280548015612c9b57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311612c70575b50505050509050925092509250909192565b604080516080810182526000808252602082018190529181019190915260608082015260408051608081018252600f805461ffff811683526201000081046bffffffffffffffffffffffff166020808501919091526e01000000000000000000000000000090910460e01b7fffffffff0000000000000000000000000000000000000000000000000000000016838501526010805485518184028101840190965280865293949293606086019392830182828015612db657602002820191906000526020600020906000905b82829054906101000a900463ffffffff1663ffffffff1681526020019060040190602082600301049283019260010382029150808411612d795790505b505050505081525050905090565b6060612dd1848484614253565b949350505050565b612de1613bb3565b612dea82613bbb565b612df2613c81565b612dfb82614320565b1561080a576040517fb42f66e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b612e3a613277565b600a546040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015260009173ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa158015612ea9573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612ecd9190615547565b6009549091506b01000000000000000000000090046bffffffffffffffffffffffff1681811115612f34576040517fa99da3020000000000000000000000000000000000000000000000000000000081526004810182905260248101839052604401610874565b8181101561093d576000612f4882846159b6565b600a546040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff87811660048301526024820184905292935091169063a9059cbb906044016020604051808303816000875af1158015612fc3573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612fe7919061582c565b506040805173ffffffffffffffffffffffffffffffffffffffff86168152602081018390527f59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b436600910160405180910390a150505050565b6000610fc882614320565b613050613663565b6000600660008381526020019081526020016000206040518060600160405290816000820154815260200160018201805461308a906155b2565b80601f01602080910402602001604051908101604052809291908181526020018280546130b6906155b2565b80156131035780601f106130d857610100808354040283529160200191613103565b820191906000526020600020905b8154815290600101906020018083116130e657829003601f168201915b5050505050815260200160028201548152505090508060400151431015613156576040517fa93d035c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8161317b57613168816020015161446f565b6020808201518051910120600755613223565b6131868260006136e9565b73ffffffffffffffffffffffffffffffffffffffff16638cc6acce82602001516040518263ffffffff1660e01b81526004016131c29190614b88565b600060405180830381600087803b1580156131dc57600080fd5b505af19250505080156131ed575060015b613223576040517ffe680b2600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b805160208201516040517fc07626096b49b0462a576c9a7878cf0675e784bb4832584c1289c11664e55f479261325a928692615ae8565b60405180910390a15050565b61326e613663565b612bcf816145a6565b610e57613663565b67ffffffffffffffff82166000908152600b60209081526040808320815160c08101835281546bffffffffffffffffffffffff808216835273ffffffffffffffffffffffffffffffffffffffff6c010000000000000000000000009283900481168488015260018501549182168487015291900416606082015260028201805484518187028101870190955280855291949293608086019390929083018282801561336057602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311613335575b505050918352505060039190910154602090910152805190915060005b82608001515181101561342157600c6000846080015183815181106133a4576133a4615518565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040908101600090812067ffffffffffffffff89168252909252902080547fffffffffffffffffffffffffffffff000000000000000000000000000000000016905561341a816157c8565b905061337d565b5067ffffffffffffffff84166000908152600b60205260408120818155600181018290559061345360028301826149b2565b60038201600090555050806009600b8282829054906101000a90046bffffffffffffffffffffffff166134869190615800565b82546101009290920a6bffffffffffffffffffffffff818102199093169183160217909155600a546040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff878116600483015292851660248201529116915063a9059cbb906044016020604051808303816000875af1158015613525573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613549919061582c565b61357f576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805173ffffffffffffffffffffffffffffffffffffffff851681526bffffffffffffffffffffffff8316602082015267ffffffffffffffff8616917fe8ed5b475a5b5987aa9165e8731bb78043f39eee32ec5a1169a89e27fcd49815910160405180910390a250505050565b67ffffffffffffffff81166000908152600b60205260409020546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff16612bcf576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600054610100900473ffffffffffffffffffffffffffffffffffffffff163314610e57576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610874565b6000816137285760008381526002602052604090205473ffffffffffffffffffffffffffffffffffffffff168015613722579050610fc8565b506137ab565b60005b6003548110156137a957600380548290811061374957613749615518565b9060005260206000200154840361379957600480548290811061376e5761376e615518565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff169150610fc89050565b6137a2816157c8565b905061372b565b505b6040517f80833e3300000000000000000000000000000000000000000000000000000000815260048101849052602401610874565b6137e86146a1565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa335b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390a1565b6000613867613bb3565b613870856135ed565b61387a338661470d565b6138848583610818565b600061389088886136e9565b604080516101008101825233815267ffffffffffffffff89166000818152600b6020818152858320805473ffffffffffffffffffffffffffffffffffffffff6c010000000000000000000000009091048116838801528688018e90526060870186905261ffff8d1660808801529484529190526003015460a084015263ffffffff881660c0840152600f546201000090046bffffffffffffffffffffffff1660e084015292517fbdd7e88000000000000000000000000000000000000000000000000000000000815293945091929184169163bdd7e8809161397491600401615b10565b610160604051808303816000875af1158015613994573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906139b89190615c11565b90506139c933888360a001516147a9565b604051806101600160405280600f60000160029054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff1681526020013373ffffffffffffffffffffffffffffffffffffffff1681526020018867ffffffffffffffff1681526020018563ffffffff1681526020018260a001516bffffffffffffffffffffffff1681526020018260c0015164ffffffffff1681526020018260e00151815260200182610100015169ffffffffffffffffffff16815260200182610120015164ffffffffff16815260200182610140015164ffffffffff16815250604051602001613ad69190615866565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0018152828252805160209182012060e0850180516000908152600e8452848120929092555167ffffffffffffffff8c16808352600b9093529290205490928c92917f7c720ccd20069b8311a6be4ba1cf3294d09eb247aa5d73a8502054b6e68a2f5491613b9b916c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1690339032908e908e908e90615ce4565b60405180910390a460e0015198975050505050505050565b610e57614884565b67ffffffffffffffff81166000908152600b60205260409020546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1680613c32576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff821614610814576040517f5a68151d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b613cac7fd8e0666292c202b1ce6a8ff0dd638652e662402ac53fbf9bd9d3bcc39d5eb09760006136e9565b604080516000815260208101918290527f6b14daf80000000000000000000000000000000000000000000000000000000090915273ffffffffffffffffffffffffffffffffffffffff9190911690636b14daf890613d0f90339060248101615d48565b602060405180830381865afa158015613d2c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613d50919061582c565b610e57576040517f22906263000000000000000000000000000000000000000000000000000000008152336004820152602401610874565b613d90614884565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a2586138333390565b60006bffffffffffffffffffffffff821115613e81576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f53616665436173743a2076616c756520646f65736e27742066697420696e203960448201527f36206269747300000000000000000000000000000000000000000000000000006064820152608401610874565b5090565b60408051606080820183526000808352602083015291810191909152600f546040516000916e010000000000000000000000000000900460e01b90613ed290899089908990602401615d77565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00181529181526020820180517fffffffff00000000000000000000000000000000000000000000000000000000949094167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff909416939093179092528151608480825260c0820190935290925060009182918291602082018180368337019050509050853b613f8457600080fd5b5a611388811015613f9457600080fd5b611388810390508760408204820311613fac57600080fd5b60008086516020880160008b8df193505a900391503d6084811115613fcf575060845b808252806000602084013e50604080516060810182529315158452602084019290925290820152979650505050505050565b6040805180820190915260008082526020820152600061402084613de3565b61402a9086615a16565b90506000816140398886615a3e565b6140439190615a3e565b6040805180820182526bffffffffffffffffffffffff808616825280841660208084019190915267ffffffffffffffff8f166000908152600b90915292832080549297509394508493929161409a91859116615800565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555081846140d49190615a3e565b336000908152600d6020526040812080549091906141019084906bffffffffffffffffffffffff16615a3e565b82546101009290920a6bffffffffffffffffffffffff818102199093169183160217909155306000908152600d6020526040812080548b9450909261414891859116615a3e565b82546101009290920a6bffffffffffffffffffffffff81810219909316918316021790915567ffffffffffffffff8c166000908152600b6020526040812060010180548d9450909261419c91859116615800565b82546bffffffffffffffffffffffff9182166101009390930a92830291909202199091161790555073ffffffffffffffffffffffffffffffffffffffff88166000908152600c6020908152604080832067ffffffffffffffff808f168552925290912080546001926009916142209185916901000000000000000000900416615995565b92506101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055505050979650505050505050565b6060600080808061426686880188615dac565b9350935093509350600061427f8960018787878761385d565b60408051602080825281830190925291925060208201818036833701905050955060005b6020811015614313578181602081106142be576142be615518565b1a60f81b8782815181106142d4576142d4615518565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535061430c816157c8565b90506142a3565b5050505050509392505050565b67ffffffffffffffff81166000908152600b602090815260408083206002018054825181850281018501909352808352849383018282801561439857602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff16815260019091019060200180831161436d575b5050505050905060005b8151811015614465576000600c60008484815181106143c3576143c3615518565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040908101600090812067ffffffffffffffff808a168352908452908290208251606081018452905460ff811615158252610100810483169482018590526901000000000000000000900490911691810182905292501461445457506001949350505050565b5061445e816157c8565b90506143a2565b5060009392505050565b600080600080848060200190518101906144899190615e21565b6040805160808101825261ffff86168082526bffffffffffffffffffffffff861660208084018290527fffffffff0000000000000000000000000000000000000000000000000000000087169484019490945260608301859052600f80547fffffffffffffffffffffffffffffffffffff00000000000000000000000000001690921762010000909102177fffffffffffffffffffffffffffff00000000ffffffffffffffffffffffffffff166e01000000000000000000000000000060e087901c021781558351969a50949850929650909450909261456f91601091908601906149d0565b509050507fe6a1eda76d42a6d1d813f26765716562044db1e8bd8be7d088705e64afd301ca838383604051610b4293929190615f1c565b3373ffffffffffffffffffffffffffffffffffffffff821603614625576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610874565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929361010090910416917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60005460ff16610e57576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f5061757361626c653a206e6f74207061757365640000000000000000000000006044820152606401610874565b73ffffffffffffffffffffffffffffffffffffffff82166000908152600c6020908152604080832067ffffffffffffffff8516845290915290205460ff16610814576040517ff0019fe600000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8216600482015273ffffffffffffffffffffffffffffffffffffffff83166024820152604401610874565b67ffffffffffffffff82166000908152600b6020526040812060010180548392906147e39084906bffffffffffffffffffffffff16615a3e565b82546bffffffffffffffffffffffff91821661010093840a908102920219161790915573ffffffffffffffffffffffffffffffffffffffff85166000908152600c6020908152604080832067ffffffffffffffff8089168552925290912080546001945090928492614859928492900416615995565b92506101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550505050565b60005460ff1615610e57576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f5061757361626c653a20706175736564000000000000000000000000000000006044820152606401610874565b82805482825590600052602060002090810192821561492c579160200282015b8281111561492c578251825591602001919060010190614911565b50613e81929150614a76565b82805482825590600052602060002090810192821561492c579160200282015b8281111561492c57825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190614958565b5080546000825590600052602060002090810190612bcf9190614a76565b8280548282559060005260206000209060070160089004810192821561492c5791602002820160005b83821115614a3d57835183826101000a81548163ffffffff021916908363ffffffff16021790555092602001926004016020816003010492830192600103026149f9565b8015614a6d5782816101000a81549063ffffffff0219169055600401602081600301049283019260010302614a3d565b5050613e819291505b5b80821115613e815760008155600101614a77565b67ffffffffffffffff81168114612bcf57600080fd5b8035614aac81614a8b565b919050565b600060208284031215614ac357600080fd5b81356114d781614a8b565b63ffffffff81168114612bcf57600080fd5b8035614aac81614ace565b60008060408385031215614afe57600080fd5b8235614b0981614a8b565b91506020830135614b1981614ace565b809150509250929050565b6000815180845260005b81811015614b4a57602081850181015186830182015201614b2e565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b6020815260006114d76020830184614b24565b60008060408385031215614bae57600080fd5b8235614bb981614a8b565b946020939093013593505050565b60008083601f840112614bd957600080fd5b50813567ffffffffffffffff811115614bf157600080fd5b602083019150836020828501011115614c0957600080fd5b9250929050565b600080600060408486031215614c2557600080fd5b83359250602084013567ffffffffffffffff811115614c4357600080fd5b614c4f86828701614bc7565b9497909650939450505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610160810167ffffffffffffffff81118282101715614caf57614caf614c5c565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715614cfc57614cfc614c5c565b604052919050565b600067ffffffffffffffff821115614d1e57614d1e614c5c565b5060051b60200190565b73ffffffffffffffffffffffffffffffffffffffff81168114612bcf57600080fd5b8035614aac81614d28565b600082601f830112614d6657600080fd5b81356020614d7b614d7683614d04565b614cb5565b82815260059290921b84018101918181019086841115614d9a57600080fd5b8286015b84811015614dbe578035614db181614d28565b8352918301918301614d9e565b509695505050505050565b60008060408385031215614ddc57600080fd5b823567ffffffffffffffff80821115614df457600080fd5b818501915085601f830112614e0857600080fd5b81356020614e18614d7683614d04565b82815260059290921b84018101918181019089841115614e3757600080fd5b948201945b83861015614e5557853582529482019490820190614e3c565b96505086013592505080821115614e6b57600080fd5b50614e7885828601614d55565b9150509250929050565b61ffff81168114612bcf57600080fd5b60008060008060008060a08789031215614eab57600080fd5b8635614eb681614a8b565b9550602087013567ffffffffffffffff811115614ed257600080fd5b614ede89828a01614bc7565b9096509450506040870135614ef281614e82565b92506060870135614f0281614ace565b80925050608087013590509295509295509295565b60008060408385031215614f2a57600080fd5b8235614f3581614a8b565b91506020830135614b1981614d28565b6bffffffffffffffffffffffff81168114612bcf57600080fd5b8035614aac81614f45565b60008060408385031215614f7d57600080fd5b8235614f8881614d28565b91506020830135614b1981614f45565b60008060408385031215614fab57600080fd5b8235614fb681614d28565b91506020830135614b1981614a8b565b8015158114612bcf57600080fd5b60008060408385031215614fe757600080fd5b823591506020830135614b1981614fc6565b6000806020838503121561500c57600080fd5b823567ffffffffffffffff8082111561502457600080fd5b818501915085601f83011261503857600080fd5b81358181111561504757600080fd5b8660206101608302850101111561505d57600080fd5b60209290920196919550909350505050565b600082601f83011261508057600080fd5b813567ffffffffffffffff81111561509a5761509a614c5c565b6150cb60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601614cb5565b8181528460208386010111156150e057600080fd5b816020850160208301376000918101602001919091529392505050565b64ffffffffff81168114612bcf57600080fd5b8035614aac816150fd565b69ffffffffffffffffffff81168114612bcf57600080fd5b8035614aac8161511b565b6000610160828403121561515157600080fd5b615159614c8b565b905061516482614f5f565b815261517260208301614d4a565b602082015261518360408301614d4a565b604082015261519460608301614aa1565b60608201526151a560808301614ae0565b60808201526151b660a08301614f5f565b60a08201526151c760c08301615110565b60c082015260e082013560e08201526101006151e4818401615133565b908201526101206151f6838201615110565b90820152610140615208838201615110565b9082015292915050565b600080600080600080610200878903121561522c57600080fd5b863567ffffffffffffffff8082111561524457600080fd5b6152508a838b0161506f565b9750602089013591508082111561526657600080fd5b5061527389828a0161506f565b955050604087013561528481614f45565b9350606087013561529481614f45565b925060808701356152a481614d28565b91506152b38860a0890161513e565b90509295509295509295565b600081518084526020808501945080840160005b8381101561530557815173ffffffffffffffffffffffffffffffffffffffff16875295820195908201906001016152d3565b509495945050505050565b6bffffffffffffffffffffffff86811682528516602082015273ffffffffffffffffffffffffffffffffffffffff84811660408301528316606082015260a060808201819052600090610ea2908301846152bf565b6000806000806060858703121561537b57600080fd5b843561538681614d28565b935060208501359250604085013567ffffffffffffffff8111156153a957600080fd5b6153b587828801614bc7565b95989497509550505050565b6000602082840312156153d357600080fd5b5035919050565b6000602082840312156153ec57600080fd5b81356114d781614e82565b6000606082018583526020606081850152818651808452608086019150828801935060005b818110156154385784518352938301939183019160010161541c565b5050848103604086015261544c81876152bf565b98975050505050505050565b600081518084526020808501945080840160005b8381101561530557815163ffffffff168752958201959082019060010161546c565b6020815261ffff82511660208201526bffffffffffffffffffffffff60208301511660408201527fffffffff00000000000000000000000000000000000000000000000000000000604083015116606082015260006060830151608080840152612dd160a0840182615458565b60006020828403121561550d57600080fd5b81356114d781614d28565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60006020828403121561555957600080fd5b5051919050565b8183823760009101908152919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b80820180821115610fc857610fc8615570565b600181811c908216806155c657607f821691505b6020821081036155ff577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f82111561093d57600081815260208120601f850160051c8101602086101561562c5750805b601f850160051c820191505b8181101561564b57828155600101615638565b505050505050565b815167ffffffffffffffff81111561566d5761566d614c5c565b6156818161567b84546155b2565b84615605565b602080601f8311600181146156d4576000841561569e5750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b17855561564b565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b8281101561572157888601518255948401946001909101908401615702565b508582101561575d57878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b84815283602082015260606040820152816060820152818360808301376000818301608090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01601019392505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036157f9576157f9615570565b5060010190565b6bffffffffffffffffffffffff82811682821603908082111561582557615825615570565b5092915050565b60006020828403121561583e57600080fd5b81516114d781614fc6565b6000610160828403121561585c57600080fd5b6114d7838361513e565b81516bffffffffffffffffffffffff168152610160810160208301516158a4602084018273ffffffffffffffffffffffffffffffffffffffff169052565b5060408301516158cc604084018273ffffffffffffffffffffffffffffffffffffffff169052565b5060608301516158e8606084018267ffffffffffffffff169052565b506080830151615900608084018263ffffffff169052565b5060a083015161592060a08401826bffffffffffffffffffffffff169052565b5060c083015161593960c084018264ffffffffff169052565b5060e083015160e0830152610100808401516159628285018269ffffffffffffffffffff169052565b50506101208381015164ffffffffff81168483015250506101408381015164ffffffffff8116848301525b505092915050565b67ffffffffffffffff81811683821601908082111561582557615825615570565b81810381811115610fc857610fc8615570565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b64ffffffffff81811683821601908082111561582557615825615570565b6bffffffffffffffffffffffff81811683821602808216919082811461598d5761598d615570565b6bffffffffffffffffffffffff81811683821601908082111561582557615825615570565b6bffffffffffffffffffffffff8616815273ffffffffffffffffffffffffffffffffffffffff8516602082015260ff8416604082015260a060608201526000615aaf60a0830185614b24565b828103608084015261544c8185614b24565b600067ffffffffffffffff808316818103615ade57615ade615570565b6001019392505050565b838152826020820152606060408201526000615b076060830184614b24565b95945050505050565b60208152600073ffffffffffffffffffffffffffffffffffffffff808451166020840152806020850151166040840152506040830151610100806060850152615b5d610120850183614b24565b91506060850151615b7a608086018267ffffffffffffffff169052565b50608085015161ffff811660a08601525060a085015160c085015260c0850151615bac60e086018263ffffffff169052565b5060e08501516bffffffffffffffffffffffff8116858301525090949350505050565b8051614aac81614f45565b8051614aac81614d28565b8051614aac81614a8b565b8051614aac81614ace565b8051614aac816150fd565b8051614aac8161511b565b60006101608284031215615c2457600080fd5b615c2c614c8b565b615c3583615bcf565b8152615c4360208401615bda565b6020820152615c5460408401615bda565b6040820152615c6560608401615be5565b6060820152615c7660808401615bf0565b6080820152615c8760a08401615bcf565b60a0820152615c9860c08401615bfb565b60c082015260e083015160e0820152610100615cb5818501615c06565b90820152610120615cc7848201615bfb565b90820152610140615cd9848201615bfb565b908201529392505050565b600073ffffffffffffffffffffffffffffffffffffffff8089168352808816602084015280871660408401525060c06060830152615d2560c0830186614b24565b905061ffff8416608083015263ffffffff831660a0830152979650505050505050565b73ffffffffffffffffffffffffffffffffffffffff83168152604060208201526000612dd16040830184614b24565b838152606060208201526000615d906060830185614b24565b8281036040840152615da28185614b24565b9695505050505050565b60008060008060808587031215615dc257600080fd5b8435615dcd81614a8b565b9350602085013567ffffffffffffffff811115615de957600080fd5b615df58782880161506f565b9350506040850135615e0681614e82565b91506060850135615e1681614ace565b939692955090935050565b60008060008060808587031215615e3757600080fd5b8451615e4281614e82565b80945050602080860151615e5581614f45565b60408701519094507fffffffff0000000000000000000000000000000000000000000000000000000081168114615e8b57600080fd5b606087015190935067ffffffffffffffff811115615ea857600080fd5b8601601f81018813615eb957600080fd5b8051615ec7614d7682614d04565b81815260059190911b8201830190838101908a831115615ee657600080fd5b928401925b82841015615f0d578351615efe81614ace565b82529284019290840190615eeb565b979a9699509497505050505050565b6bffffffffffffffffffffffff841681527fffffffff0000000000000000000000000000000000000000000000000000000083166020820152606060408201526000615b07606083018461545856fea164736f6c6343000813000a",
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
	RequestId      [32]byte
	SubscriptionId uint64
	TotalCostJuels *big.Int
	Transmitter    common.Address
	ResultCode     uint8
	Response       []byte
	ReturnData     []byte
	Raw            types.Log
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

func (FunctionsRouterRequestNotProcessed) Topic() common.Hash {
	return common.HexToHash("0x1a90e9a50793db2e394cf581e7c522e10c358a81e70acf6b5a0edd620c08dee1")
}

func (FunctionsRouterRequestProcessed) Topic() common.Hash {
	return common.HexToHash("0x47ffcaa55fde21cc7135c65541826c9e65dda59c29dc109aae964989e8fc664b")
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
