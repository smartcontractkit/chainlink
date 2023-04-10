// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ocr2dr_registry

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

type FunctionsBillingRegistryCommitment struct {
	SubscriptionId uint64
	Client         common.Address
	GasLimit       uint32
	GasPrice       *big.Int
	Don            common.Address
	DonFee         *big.Int
	RegistryFee    *big.Int
	EstimatedCost  *big.Int
	Timestamp      *big.Int
}

type FunctionsBillingRegistryInterfaceRequestBilling struct {
	SubscriptionId uint64
	Client         common.Address
	GasLimit       uint32
	GasPrice       *big.Int
}

var OCR2DRRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"linkEthFeed\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"internalBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"externalBalance\",\"type\":\"uint256\"}],\"name\":\"BalanceInvariantViolated\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotSelfTransfer\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptySendersList\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"have\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"want\",\"type\":\"uint32\"}],\"name\":\"GasLimitTooBig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"InvalidConsumer\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"linkWei\",\"type\":\"int256\"}],\"name\":\"InvalidLinkWeiPrice\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSubscription\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposedOwner\",\"type\":\"address\"}],\"name\":\"MustBeRequestedOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"MustBeSubOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAllowedToSetSenders\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotProposedOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableFromLink\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OwnerMustBeSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentTooLarge\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingRequestExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Reentrant\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyConsumers\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedSender\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"senders\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"changedBy\",\"type\":\"address\"}],\"name\":\"AuthorizedSendersChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"signerPayment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"transmitterPayment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalCost\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"BillingEnd\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"don\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"donFee\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"registryFee\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"estimatedCost\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structFunctionsBillingRegistry.Commitment\",\"name\":\"commitment\",\"type\":\"tuple\"}],\"name\":\"BillingStart\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasOverhead\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"fee\",\"type\":\"uint96\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"RequestTimedOut\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"SubscriptionCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"SubscriptionCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"MAX_CONSUMERS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"acceptSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"addConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"cancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createSubscription\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"donFee\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"registryFee\",\"type\":\"uint96\"}],\"name\":\"estimateCost\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"response\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"internalType\":\"address[31]\",\"name\":\"signers\",\"type\":\"address[31]\"},{\"internalType\":\"uint8\",\"name\":\"signerCount\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"reportValidationGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"initialGas\",\"type\":\"uint256\"}],\"name\":\"fulfillAndBill\",\"outputs\":[{\"internalType\":\"enumFunctionsBillingRegistryInterface.FulfillResult\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAuthorizedSenders\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfig\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"internalType\":\"uint32\",\"name\":\"gasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"linkAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"linkPriceFeed\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCurrentsubscriptionId\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRequestConfig\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"},{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"}],\"internalType\":\"structFunctionsBillingRegistryInterface.RequestBilling\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"getRequiredFee\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"getSubscription\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"getSubscriptionOwner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTotalBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"linkEthFeed\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"isAuthorizedSender\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"oracleWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"ownerCancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"pendingRequestExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"removeConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"requestSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"senders\",\"type\":\"address[]\"}],\"name\":\"setAuthorizedSenders\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"internalType\":\"uint32\",\"name\":\"gasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"requestTimeoutSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"fee\",\"type\":\"uint96\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"}],\"internalType\":\"structFunctionsBillingRegistryInterface.RequestBilling\",\"name\":\"billing\",\"type\":\"tuple\"}],\"name\":\"startBilling\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"requestIdsToTimeout\",\"type\":\"bytes32[]\"}],\"name\":\"timeoutRequests\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162005dec38038062005dec83398101604081905262000034916200040c565b620000418383836200004a565b50505062000456565b600054610100900460ff16158080156200006b5750600054600160ff909116105b806200009b57506200008830620001c960201b62003c721760201c565b1580156200009b575060005460ff166001145b620001045760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b60648201526084015b60405180910390fd5b6000805460ff19166001179055801562000128576000805461ff0019166101001790555b62000132620001d8565b6200013f33600062000240565b606980546001600160a01b038087166001600160a01b031992831617909255606a8054868416908316179055606b8054928516929091169190911790558015620001c3576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b50505050565b6001600160a01b03163b151590565b600054610100900460ff16620002345760405162461bcd60e51b815260206004820152602b602482015260008051602062005dcc83398151915260448201526a6e697469616c697a696e6760a81b6064820152608401620000fb565b6200023e62000304565b565b600054610100900460ff166200029c5760405162461bcd60e51b815260206004820152602b602482015260008051602062005dcc83398151915260448201526a6e697469616c697a696e6760a81b6064820152608401620000fb565b6001600160a01b038216620002c457604051635b5a8afd60e11b815260040160405180910390fd5b600080546001600160a01b03808516620100000262010000600160b01b031990921691909117909155811615620003005762000300816200036c565b5050565b600054610100900460ff16620003605760405162461bcd60e51b815260206004820152602b602482015260008051602062005dcc83398151915260448201526a6e697469616c697a696e6760a81b6064820152608401620000fb565b6034805460ff19169055565b6001600160a01b038116331415620003975760405163282010c360e01b815260040160405180910390fd5b600180546001600160a01b0319166001600160a01b038381169182179092556000805460405192936201000090910416917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b03811681146200040757600080fd5b919050565b6000806000606084860312156200042257600080fd5b6200042d84620003ef565b92506200043d60208501620003ef565b91506200044d60408501620003ef565b90509250925092565b61596680620004666000396000f3fe608060405234801561001057600080fd5b50600436106102255760003560e01c80638da5cb5b1161012a578063c0c53b8b116100bd578063e82ad7d41161008c578063f1e14a2111610071578063f1e14a2114610561578063f2fde38b14610593578063fa00763a146105a657600080fd5b8063e82ad7d41461053b578063ee56997b1461054e57600080fd5b8063c0c53b8b1461048d578063c3f909d4146104a0578063d7ae1d3014610515578063e72f6e301461052857600080fd5b8063a47c7696116100f9578063a47c769614610432578063a4c0ed3614610454578063a9d03c0514610467578063b2a489ff1461047a57600080fd5b80638da5cb5b146103a25780639f87fad7146103e7578063a1a6d041146103fa578063a21a23e41461042a57600080fd5b80634dc14c7d116101bd578063665871ec1161018c57806379ba50971161017157806379ba50971461037f57806382359740146103875780638456cb591461039a57600080fd5b8063665871ec146103595780637341c10c1461036c57600080fd5b80634dc14c7d146103015780635c975abb1461031457806364d51a2a1461032b57806366316d8d1461034657600080fd5b806312b58349116101f957806312b58349146102915780632408afaa146102bd57806333652e3e146102d25780633f4ba83a146102f957600080fd5b80620122911461022a57806302bcc5b61461024957806304c357cb1461025e5780630739e4f114610271575b600080fd5b6102326105b9565b604051610240929190615655565b60405180910390f35b61025c61025736600461528d565b6105d8565b005b61025c61026c3660046152a8565b610655565b61028461027f366004614f7d565b610850565b6040516102409190615524565b606f546801000000000000000090046bffffffffffffffffffffffff165b604051908152602001610240565b6102c5610f12565b604051610240919061545a565b606f5467ffffffffffffffff165b60405167ffffffffffffffff9091168152602001610240565b61025c610f81565b61025c61030f366004615211565b610f93565b60345460ff165b6040519015158152602001610240565b610333606481565b60405161ffff9091168152602001610240565b61025c610354366004614ee2565b611169565b61025c610367366004614f19565b6113cf565b61025c61037a3660046152a8565b61169b565b61025c61192e565b61025c61039536600461528d565b611a21565b61025c611da2565b60005462010000900473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610240565b61025c6103f53660046152a8565b611db2565b61040d6104083660046151c0565b612239565b6040516bffffffffffffffffffffffff9091168152602001610240565b6102e061235d565b61044561044036600461528d565b6126f5565b6040516102409392919061567c565b61025c610462366004614e88565b612826565b6102af6104753660046150a7565b612a7f565b6103c261048836600461528d565b61327a565b61025c61049b366004614e45565b613313565b60a75460a85460725460a954606954606a546040805163ffffffff808916825265010000000000909804881660208201529081019590955260608501939093529316608083015273ffffffffffffffffffffffffffffffffffffffff92831660a08301529190911660c082015260e001610240565b61025c6105233660046152a8565b613515565b61025c610536366004614e2a565b61367c565b61031b61054936600461528d565b613899565b61025c61055c366004614f19565b613ad8565b61040d61056f366004615124565b505060a9546801000000000000000090046bffffffffffffffffffffffff16919050565b61025c6105a1366004614e2a565b613c4b565b61031b6105b4366004614e2a565b613c5f565b60a75460009060609063ffffffff166105d0610f12565b915091509091565b6105e0613c8e565b67ffffffffffffffff81166000908152606d602052604090205473ffffffffffffffffffffffffffffffffffffffff1680610647576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6106518282613ce5565b5050565b67ffffffffffffffff82166000908152606d6020526040902054829073ffffffffffffffffffffffffffffffffffffffff16806106be576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff82161461072a576040517fd8a3fb5200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821660048201526024015b60405180910390fd5b60a754640100000000900460ff161561076f576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6107776140ec565b67ffffffffffffffff84166000908152606d602052604090206001015473ffffffffffffffffffffffffffffffffffffffff84811691161461084a5767ffffffffffffffff84166000818152606d602090815260409182902060010180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff88169081179091558251338152918201527f69436ea6df009049404f564eff6622cd00522b0bd6a89efd9e52a355c4a879be91015b60405180910390a25b50505050565b600061085a614159565b60a754640100000000900460ff161561089f576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6108a76140ec565b60008b815260716020908152604091829020825161012081018452815467ffffffffffffffff8116825268010000000000000000810473ffffffffffffffffffffffffffffffffffffffff908116948301949094527c0100000000000000000000000000000000000000000000000000000000900463ffffffff169381019390935260018101546060840152600281015491821660808401819052740100000000000000000000000000000000000000009092046bffffffffffffffffffffffff90811660a0850152600382015480821660c08601526c0100000000000000000000000090041660e0840152600401546101008301526109ab576002915050610f04565b60008c81526071602052604080822082815560018101839055600281018390556003810180547fffffffffffffffff000000000000000000000000000000000000000000000000169055600401829055517f0ca761750000000000000000000000000000000000000000000000000000000090610a34908f908f908f908f908f9060240161546d565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152918152602080830180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009095169490941790935260a780547fffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffff1664010000000017905584015191840151909250600091610b019163ffffffff9091169084614198565b60a780547fffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffff16905560a85460a085015160c0860151929350600092610b4d92899290918c908c3a6141e4565b604080820151865167ffffffffffffffff166000908152606e60205291909120549192506bffffffffffffffffffffffff90811691161015610bbb576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b604080820151855167ffffffffffffffff166000908152606e602052918220805491929091610bf99084906bffffffffffffffffffffffff166157c6565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060005b8860ff16811015610cf3578151607060008c84601f8110610c4d57610c4d6158e1565b602002015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282829054906101000a90046bffffffffffffffffffffffff16610cb2919061570c565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055508080610ceb906157f3565b915050610c2a565b508360c0015160706000610d2360005473ffffffffffffffffffffffffffffffffffffffff620100009091041690565b73ffffffffffffffffffffffffffffffffffffffff168152602081019190915260400160009081208054909190610d699084906bffffffffffffffffffffffff1661570c565b82546101009290920a6bffffffffffffffffffffffff81810219909316918316021790915560208381015173ffffffffffffffffffffffffffffffffffffffff8e166000908152607090925260408220805491945092610dcb9185911661570c565b82546101009290920a6bffffffffffffffffffffffff81810219909316918316021790915560e0860151865167ffffffffffffffff166000908152606e60205260409020805491935091600c91610e349185916c010000000000000000000000009004166157c6565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055508e7fc8dc973332de19a5f71b6026983110e9c2e04b0c98b87eb771ccb78607fd114f856000015183600001518460200151856040015187604051610ee695949392919067ffffffffffffffff9590951685526bffffffffffffffffffffffff9384166020860152918316604085015290911660608301521515608082015260a00190565b60405180910390a281610efa576001610efd565b60005b9450505050505b9a9950505050505050505050565b60606068805480602002602001604051908101604052809291908181526020018280548015610f7757602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610f4c575b5050505050905090565b610f89613c8e565b610f91614364565b565b610f9b613c8e565b60008413610fd8576040517f43d4cf6600000000000000000000000000000000000000000000000000000000815260048101859052602401610721565b6040805160e08101825263ffffffff808a1680835260006020840152898216838501819052606084018a90528783166080850181905292871660a085018190526bffffffffffffffffffffffff871660c090950185905260a78054650100000000009093027fffffffffffffffffffffffffffffffffffffffffffffff0000000000000000009093169093179190911790915560a889905560a98054680100000000000000009094027fffffffffffffffffffffffff000000000000000000000000ffffffffffffffff6401000000009093027fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000090951690931793909317161790556072859055517f1f8b76dfe267648827eba3c1966a7149a45245c5af445ce68ba124d04b7520aa906111589089908990899089908990889063ffffffff968716815294861660208601526040850193909352606084019190915290921660808201526bffffffffffffffffffffffff9190911660a082015260c00190565b60405180910390a150505050505050565b60a754640100000000900460ff16156111ae576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6111b66140ec565b6bffffffffffffffffffffffff81166111e95750336000908152607060205260409020546bffffffffffffffffffffffff165b336000908152607060205260409020546bffffffffffffffffffffffff80831691161015611243576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b33600090815260706020526040812080548392906112709084906bffffffffffffffffffffffff166157c6565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555080606f60088282829054906101000a90046bffffffffffffffffffffffff166112c791906157c6565b82546101009290920a6bffffffffffffffffffffffff8181021990931691831602179091556069546040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff868116600483015292851660248201529116915063a9059cbb90604401602060405180830381600087803b15801561136157600080fd5b505af1158015611375573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906113999190614f5b565b610651576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6113d76140ec565b60005b818110156116965760008383838181106113f6576113f66158e1565b602090810292909201356000818152607184526040808220815161012081018352815467ffffffffffffffff811680835268010000000000000000820473ffffffffffffffffffffffffffffffffffffffff908116848b01527c010000000000000000000000000000000000000000000000000000000090920463ffffffff168386015260018401546060840152600284015480831660808501527401000000000000000000000000000000000000000090046bffffffffffffffffffffffff90811660a0850152600385015480821660c08601526c0100000000000000000000000090041660e0840152600490930154610100830152918452606d909652912054919450163314905061156b57805167ffffffffffffffff166000908152606d6020526040908190205490517fd8a3fb5200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602401610721565b60a954610100820151429161158f9164010000000090910463ffffffff16906156c8565b11156116815760e0810151815167ffffffffffffffff166000908152606e602052604090208054600c906115e29084906c0100000000000000000000000090046bffffffffffffffffffffffff166157c6565b82546bffffffffffffffffffffffff9182166101009390930a92830291909202199091161790555060008281526071602052604080822082815560018101839055600281018390556003810180547fffffffffffffffff0000000000000000000000000000000000000000000000001690556004018290555183917ff1ca1e9147be737b04a2b018a79405f687a97de8dd8a2559bbe62357343af41491a25b5050808061168e906157f3565b9150506113da565b505050565b67ffffffffffffffff82166000908152606d6020526040902054829073ffffffffffffffffffffffffffffffffffffffff1680611704576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff82161461176b576040517fd8a3fb5200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82166004820152602401610721565b60a754640100000000900460ff16156117b0576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6117b86140ec565b67ffffffffffffffff84166000908152606d60205260409020600201546064141561180f576040517f05a48e0f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff83166000908152606c6020908152604080832067ffffffffffffffff808916855292529091205416156118565761084a565b73ffffffffffffffffffffffffffffffffffffffff83166000818152606c6020908152604080832067ffffffffffffffff891680855290835281842080547fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000166001908117909155606d84528285206002018054918201815585529383902090930180547fffffffffffffffffffffffff000000000000000000000000000000000000000016851790555192835290917f43dc749a04ac8fb825cbd514f7c0e13f13bc6f2ee66043b76629d51776cff8e09101610841565b60015473ffffffffffffffffffffffffffffffffffffffff16331461197f576040517f0f22ca5f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000805433620100008181027fffffffffffffffffffff0000000000000000000000000000000000000000ffff8416178455600180547fffffffffffffffffffffffff000000000000000000000000000000000000000016905560405173ffffffffffffffffffffffffffffffffffffffff919093041692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60a754640100000000900460ff1615611a66576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611a6e6140ec565b606b60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16634b4fa0c16040518163ffffffff1660e01b815260040160206040518083038186803b158015611ad657600080fd5b505afa158015611aea573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611b0e9190614f5b565b8015611bb85750606b546040517ffa00763a00000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff9091169063fa00763a9060240160206040518083038186803b158015611b7e57600080fd5b505afa158015611b92573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611bb69190614f5b565b155b15611bef576040517f0809490800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff81166000908152606d602052604090205473ffffffffffffffffffffffffffffffffffffffff16611c55576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff81166000908152606d602052604090206001015473ffffffffffffffffffffffffffffffffffffffff163314611cf75767ffffffffffffffff81166000908152606d6020526040908190206001015490517fd084e97500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602401610721565b67ffffffffffffffff81166000818152606d60209081526040918290208054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560019093018054909316909255835173ffffffffffffffffffffffffffffffffffffffff909116808252928101919091529092917f6f1dc65165ffffedfd8e507b4a0f1fcfdada045ed11f6c26ba27cedfe87802f0910160405180910390a25050565b611daa613c8e565b610f916143e1565b67ffffffffffffffff82166000908152606d6020526040902054829073ffffffffffffffffffffffffffffffffffffffff1680611e1b576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff821614611e82576040517fd8a3fb5200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82166004820152602401610721565b60a754640100000000900460ff1615611ec7576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611ecf6140ec565b73ffffffffffffffffffffffffffffffffffffffff83166000908152606c6020908152604080832067ffffffffffffffff808916855292529091205416611f6a576040517ff0019fe600000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8516600482015273ffffffffffffffffffffffffffffffffffffffff84166024820152604401610721565b67ffffffffffffffff84166000908152606d6020908152604080832060020180548251818502810185019093528083529192909190830182828015611fe557602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311611fba575b50505050509050600060018251611ffc91906157af565b905060005b825181101561219b578573ffffffffffffffffffffffffffffffffffffffff16838281518110612033576120336158e1565b602002602001015173ffffffffffffffffffffffffffffffffffffffff16141561218957600083838151811061206b5761206b6158e1565b6020026020010151905080606d60008a67ffffffffffffffff1667ffffffffffffffff16815260200190815260200160002060020183815481106120b1576120b16158e1565b600091825260208083209190910180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff949094169390931790925567ffffffffffffffff8a168152606d9091526040902060020180548061212b5761212b6158b2565b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff00000000000000000000000000000000000000001690550190555061219b565b80612193816157f3565b915050612001565b5073ffffffffffffffffffffffffffffffffffffffff85166000818152606c6020908152604080832067ffffffffffffffff8b168085529083529281902080547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001690555192835290917f182bff9831466789164ca77075fffd84916d35a8180ba73c27e45634549b445b91015b60405180910390a2505050505050565b60008061224461443c565b905060008113612283576040517f43d4cf6600000000000000000000000000000000000000000000000000000000815260048101829052602401610721565b60a85460a95460009163ffffffff808a16926122a09291166156c8565b6122aa91906156c8565b9050600082826122c289670de0b6b3a7640000615772565b6122cc9190615772565b6122d69190615733565b905060006122f56bffffffffffffffffffffffff8088169089166156c8565b905061230d816b033b2e3c9fd0803ce80000006157af565b821115612346576040517fe80fa38100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61235081836156c8565b9998505050505050505050565b60a754600090640100000000900460ff16156123a5576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6123ad6140ec565b606b60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16634b4fa0c16040518163ffffffff1660e01b815260040160206040518083038186803b15801561241557600080fd5b505afa158015612429573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061244d9190614f5b565b80156124f75750606b546040517ffa00763a00000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff9091169063fa00763a9060240160206040518083038186803b1580156124bd57600080fd5b505afa1580156124d1573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906124f59190614f5b565b155b1561252e576040517f0809490800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b606f805467ffffffffffffffff169060006125488361582c565b82546101009290920a67ffffffffffffffff818102199093169183160217909155606f5416905060008060405190808252806020026020018201604052801561259b578160200160208202803683370190505b506040805180820182526000808252602080830182815267ffffffffffffffff8816808452606e83528584209451855492516bffffffffffffffffffffffff9081166c01000000000000000000000000027fffffffffffffffff000000000000000000000000000000000000000000000000909416911617919091179093558351606081018552338152808201838152818601878152948452606d8352949092208251815473ffffffffffffffffffffffffffffffffffffffff9182167fffffffffffffffffffffffff00000000000000000000000000000000000000009182161783559551600183018054919092169616959095179094559151805194955090936126ad9260028501920190614b76565b505060405133815267ffffffffffffffff841691507f464722b4166576d3dcbba877b999bc35cf911f4eaf434b7eba68fa113951d0bf9060200160405180910390a250905090565b67ffffffffffffffff81166000908152606d6020526040812054819060609073ffffffffffffffffffffffffffffffffffffffff16612760576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff84166000908152606e6020908152604080832054606d8352928190208054600290910180548351818602810186019094528084526bffffffffffffffffffffffff9095169473ffffffffffffffffffffffffffffffffffffffff90921693909291839183018282801561281257602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116127e7575b505050505090509250925092509193909250565b60a754640100000000900460ff161561286b576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6128736140ec565b60695473ffffffffffffffffffffffffffffffffffffffff1633146128c4576040517f44b0e3c300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b602081146128fe576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600061290c8284018461528d565b67ffffffffffffffff81166000908152606d602052604090205490915073ffffffffffffffffffffffffffffffffffffffff16612975576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff81166000908152606e6020526040812080546bffffffffffffffffffffffff16918691906129ac838561570c565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555084606f60088282829054906101000a90046bffffffffffffffffffffffff16612a03919061570c565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055508167ffffffffffffffff167fd39ec07f4e209f627a4c427971473820dc129761ba28de8906bd56f57101d4f8828784612a6a91906156c8565b60408051928352602083019190915201612229565b6000612a89614159565b60a754640100000000900460ff1615612ace576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b612ad66140ec565b6000606d81612ae8602086018661528d565b67ffffffffffffffff16815260208101919091526040016000205473ffffffffffffffffffffffffffffffffffffffff161415612b51576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000606c81612b666040860160208701614e2a565b73ffffffffffffffffffffffffffffffffffffffff168152602080820192909252604001600090812091612b9c9086018661528d565b67ffffffffffffffff908116825260208201929092526040016000205416905080612c3857612bce602084018461528d565b612bde6040850160208601614e2a565b6040517ff0019fe600000000000000000000000000000000000000000000000000000000815267ffffffffffffffff909216600483015273ffffffffffffffffffffffffffffffffffffffff166024820152604401610721565b60a75463ffffffff16612c5160608501604086016151a5565b63ffffffff161115612cb257612c6d60608401604085016151a5565b60a7546040517ff5d7e01e00000000000000000000000000000000000000000000000000000000815263ffffffff928316600482015291166024820152604401610721565b6040517ff1e14a21000000000000000000000000000000000000000000000000000000008152600090339063f1e14a2190612cf5908990899089906004016154a6565b60206040518083038186803b158015612d0d57600080fd5b505afa158015612d21573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612d45919061532b565b90506000612d5d878761056f36899003890189615170565b90506000612d80612d7460608801604089016151a5565b87606001358585612239565b90506000606e81612d9460208a018a61528d565b67ffffffffffffffff1681526020808201929092526040016000908120546c0100000000000000000000000090046bffffffffffffffffffffffff1691606e9190612de1908b018b61528d565b67ffffffffffffffff168152602081019190915260400160002054612e1491906bffffffffffffffffffffffff166157c6565b9050816bffffffffffffffffffffffff16816bffffffffffffffffffffffff161015612e6c576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000612e798660016156e0565b90506000612f0133612e9160408c0160208d01614e2a565b612e9e60208d018d61528d565b856040805173ffffffffffffffffffffffffffffffffffffffff958616602080830191909152949095168582015267ffffffffffffffff928316606086015291166080808501919091528151808503909101815260a09093019052815191012090565b60408051610120810190915290915060009080612f2160208d018d61528d565b67ffffffffffffffff1681526020018b6020016020810190612f439190614e2a565b73ffffffffffffffffffffffffffffffffffffffff168152602001612f6e60608d0160408e016151a5565b63ffffffff90811682526060808e0135602080850191909152336040808601919091526bffffffffffffffffffffffff808e16848701528c81166080808801919091528c821660a0808901919091524260c09889015260008b8152607186528481208a5181548c890151978d0151909a167c0100000000000000000000000000000000000000000000000000000000027bffffffffffffffffffffffffffffffffffffffffffffffffffffffff73ffffffffffffffffffffffffffffffffffffffff98891668010000000000000000027fffffffff00000000000000000000000000000000000000000000000000000000909c1667ffffffffffffffff909316929092179a909a171698909817885595890151600188015590880151908801518216740100000000000000000000000000000000000000000292169190911760028501559385015160038401805460e088015187166c01000000000000000000000000027fffffffffffffffff0000000000000000000000000000000000000000000000009091169290961691909117949094179093556101008401516004909201919091559192508691606e91613128908e018e61528d565b67ffffffffffffffff16815260208101919091526040016000208054600c906131709084906c0100000000000000000000000090046bffffffffffffffffffffffff1661570c565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550817f99f7f4e65b4b9fbabd4e357c47ed3099b36e57ecd3a43e84662f34c207d0ebe4826040516131ce9190615565565b60405180910390a282606c60006131eb60408e0160208f01614e2a565b73ffffffffffffffffffffffffffffffffffffffff168152602080820192909252604001600090812091613221908e018e61528d565b67ffffffffffffffff9081168252602082019290925260400160002080547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001692909116919091179055509a9950505050505050505050565b67ffffffffffffffff81166000908152606d602052604081205473ffffffffffffffffffffffffffffffffffffffff166132e0576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5067ffffffffffffffff166000908152606d602052604090205473ffffffffffffffffffffffffffffffffffffffff1690565b600054610100900460ff16158080156133335750600054600160ff909116105b8061334d5750303b15801561334d575060005460ff166001145b6133d9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a65640000000000000000000000000000000000006064820152608401610721565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055801561343757600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101001790555b61343f61452f565b61344a3360006145ce565b6069805473ffffffffffffffffffffffffffffffffffffffff8087167fffffffffffffffffffffffff000000000000000000000000000000000000000092831617909255606a8054868416908316179055606b805492851692909116919091179055801561084a57600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498906020015b60405180910390a150505050565b67ffffffffffffffff82166000908152606d6020526040902054829073ffffffffffffffffffffffffffffffffffffffff168061357e576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff8216146135e5576040517fd8a3fb5200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82166004820152602401610721565b60a754640100000000900460ff161561362a576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6136326140ec565b61363b84613899565b15613672576040517fb42f66e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61084a8484613ce5565b613684613c8e565b6069546040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015260009173ffffffffffffffffffffffffffffffffffffffff16906370a082319060240160206040518083038186803b1580156136ee57600080fd5b505afa158015613702573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613726919061518c565b606f549091506801000000000000000090046bffffffffffffffffffffffff168181111561378a576040517fa99da3020000000000000000000000000000000000000000000000000000000081526004810182905260248101839052604401610721565b8181101561169657600061379e82846157af565b6069546040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff87811660048301526024820184905292935091169063a9059cbb90604401602060405180830381600087803b15801561381457600080fd5b505af1158015613828573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061384c9190614f5b565b506040805173ffffffffffffffffffffffffffffffffffffffff86168152602081018390527f59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b4366009101613507565b67ffffffffffffffff81166000908152606d602090815260408083206002018054825181850281018501909352808352849383018282801561391157602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116138e6575b505050505090506000613922610f12565b905060005b8251811015613acd5760005b8251811015613aba576000613a6a848381518110613953576139536158e1565b602002602001015186858151811061396d5761396d6158e1565b602002602001015189606c60008a898151811061398c5761398c6158e1565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008c67ffffffffffffffff1667ffffffffffffffff16815260200190815260200160002060009054906101000a900467ffffffffffffffff166040805173ffffffffffffffffffffffffffffffffffffffff958616602080830191909152949095168582015267ffffffffffffffff928316606086015291166080808501919091528151808503909101815260a09093019052815191012090565b60008181526071602052604090206002015490915073ffffffffffffffffffffffffffffffffffffffff1615613aa7575060019695505050505050565b5080613ab2816157f3565b915050613933565b5080613ac5816157f3565b915050613927565b506000949350505050565b613ae061470e565b613b16576040517fad77f06100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80613b4d576040517f75158c3b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b606854811015613bad57613b9a60688281548110613b7057613b706158e1565b60009182526020909120015460669073ffffffffffffffffffffffffffffffffffffffff1661471e565b5080613ba5816157f3565b915050613b50565b5060005b81811015613bfe57613beb838383818110613bce57613bce6158e1565b9050602002016020810190613be39190614e2a565b606690614747565b5080613bf6816157f3565b915050613bb1565b50613c0b60688383614bfc565b507ff263cfb3e4298332e776194610cf9fdc09ccb3ada8b9aa39764d882e11fbf0a0828233604051613c3f939291906153e2565b60405180910390a15050565b613c53613c8e565b613c5c81614769565b50565b6000613c6c606683614836565b92915050565b73ffffffffffffffffffffffffffffffffffffffff163b151590565b60005462010000900473ffffffffffffffffffffffffffffffffffffffff163314610f91576040517f2b5c74de00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60a754640100000000900460ff1615613d2a576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff82166000908152606d602090815260408083208151606081018352815473ffffffffffffffffffffffffffffffffffffffff908116825260018301541681850152600282018054845181870281018701865281815292959394860193830182828015613dd557602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311613daa575b5050509190925250505067ffffffffffffffff84166000908152606e60205260408120549192506bffffffffffffffffffffffff909116905b826040015151811015613eb457606c600084604001518381518110613e3557613e356158e1565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040908101600090812067ffffffffffffffff89168252909252902080547fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000016905580613eac816157f3565b915050613e0e565b5067ffffffffffffffff84166000908152606d6020526040812080547fffffffffffffffffffffffff00000000000000000000000000000000000000009081168255600182018054909116905590613f0f6002830182614c74565b505067ffffffffffffffff84166000908152606e6020526040902080547fffffffffffffffff000000000000000000000000000000000000000000000000169055606f8054829190600890613f7f9084906801000000000000000090046bffffffffffffffffffffffff166157c6565b82546101009290920a6bffffffffffffffffffffffff8181021990931691831602179091556069546040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff878116600483015292851660248201529116915063a9059cbb90604401602060405180830381600087803b15801561401957600080fd5b505af115801561402d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906140519190614f5b565b614087576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805173ffffffffffffffffffffffffffffffffffffffff851681526bffffffffffffffffffffffff8316602082015267ffffffffffffffff8616917fe8ed5b475a5b5987aa9165e8731bb78043f39eee32ec5a1169a89e27fcd498159101610841565b60345460ff1615610f91576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f5061757361626c653a20706175736564000000000000000000000000000000006044820152606401610721565b61416233613c5f565b610f91576040517f0809490800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005a6113888110156141aa57600080fd5b6113888103905084604082048203116141c257600080fd5b50823b6141ce57600080fd5b60008083516020850160008789f1949350505050565b604080516060810182526000808252602082018190529181018290529061420961443c565b905060008113614248576040517f43d4cf6600000000000000000000000000000000000000000000000000000000815260048101829052602401610721565b6000815a8b6142578c896156c8565b61426191906156c8565b61426b91906157af565b61427d86670de0b6b3a7640000615772565b6142879190615772565b6142919190615733565b905060006142b06bffffffffffffffffffffffff808916908b166156c8565b90506142c8816b033b2e3c9fd0803ce80000006157af565b821115614301576040517fe80fa38100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600061431060ff8a168b615747565b905082600061432761432285846156c8565b614865565b604080516060810182526bffffffffffffffffffffffff958616815293851660208501529316928201929092529c9b505050505050505050505050565b61436c614907565b603480547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa335b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390a1565b6143e96140ec565b603480547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a2586143b73390565b60a754606a54604080517ffeaf968c000000000000000000000000000000000000000000000000000000008152905160009365010000000000900463ffffffff1692831515928592839273ffffffffffffffffffffffffffffffffffffffff169163feaf968c9160048083019260a0929190829003018186803b1580156144c257600080fd5b505afa1580156144d6573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906144fa91906152db565b5093505092505082801561451c575061451381426157af565b8463ffffffff16105b156145275760725491505b509392505050565b600054610100900460ff166145c6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201527f6e697469616c697a696e670000000000000000000000000000000000000000006064820152608401610721565b610f91614973565b600054610100900460ff16614665576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201527f6e697469616c697a696e670000000000000000000000000000000000000000006064820152608401610721565b73ffffffffffffffffffffffffffffffffffffffff82166146b2576040517fb6b515fa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000805473ffffffffffffffffffffffffffffffffffffffff80851662010000027fffffffffffffffffffff0000000000000000000000000000000000000000ffff909216919091179091558116156106515761065181614769565b6000614718613c8e565b50600190565b60006147408373ffffffffffffffffffffffffffffffffffffffff8416614a34565b9392505050565b60006147408373ffffffffffffffffffffffffffffffffffffffff8416614b27565b73ffffffffffffffffffffffffffffffffffffffff81163314156147b9576040517f282010c300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8381169182179092556000805460405192936201000090910416917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b73ffffffffffffffffffffffffffffffffffffffff811660009081526001830160205260408120541515614740565b60006bffffffffffffffffffffffff821115614903576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f53616665436173743a2076616c756520646f65736e27742066697420696e203960448201527f36206269747300000000000000000000000000000000000000000000000000006064820152608401610721565b5090565b60345460ff16610f91576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f5061757361626c653a206e6f74207061757365640000000000000000000000006044820152606401610721565b600054610100900460ff16614a0a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201527f6e697469616c697a696e670000000000000000000000000000000000000000006064820152608401610721565b603480547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055565b60008181526001830160205260408120548015614b1d576000614a586001836157af565b8554909150600090614a6c906001906157af565b9050818114614ad1576000866000018281548110614a8c57614a8c6158e1565b9060005260206000200154905080876000018481548110614aaf57614aaf6158e1565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080614ae257614ae26158b2565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050613c6c565b6000915050613c6c565b6000818152600183016020526040812054614b6e57508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155613c6c565b506000613c6c565b828054828255906000526020600020908101928215614bf0579160200282015b82811115614bf057825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190614b96565b50614903929150614c8e565b828054828255906000526020600020908101928215614bf0579160200282015b82811115614bf05781547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff843516178255602090920191600190910190614c1c565b5080546000825590600052602060002090810190613c5c91905b5b808211156149035760008155600101614c8f565b803573ffffffffffffffffffffffffffffffffffffffff81168114614cc757600080fd5b919050565b60008083601f840112614cde57600080fd5b50813567ffffffffffffffff811115614cf657600080fd5b6020830191508360208260051b8501011115614d1157600080fd5b9250929050565b60008083601f840112614d2a57600080fd5b50813567ffffffffffffffff811115614d4257600080fd5b602083019150836020828501011115614d1157600080fd5b600060808284031215614d6c57600080fd5b6040516080810181811067ffffffffffffffff82111715614d8f57614d8f615910565b604052905080614d9e83614de7565b8152614dac60208401614ca3565b6020820152614dbd60408401614dd3565b6040820152606083013560608201525092915050565b803563ffffffff81168114614cc757600080fd5b803567ffffffffffffffff81168114614cc757600080fd5b803560ff81168114614cc757600080fd5b805169ffffffffffffffffffff81168114614cc757600080fd5b600060208284031215614e3c57600080fd5b61474082614ca3565b600080600060608486031215614e5a57600080fd5b614e6384614ca3565b9250614e7160208501614ca3565b9150614e7f60408501614ca3565b90509250925092565b60008060008060608587031215614e9e57600080fd5b614ea785614ca3565b935060208501359250604085013567ffffffffffffffff811115614eca57600080fd5b614ed687828801614d18565b95989497509550505050565b60008060408385031215614ef557600080fd5b614efe83614ca3565b91506020830135614f0e8161593f565b809150509250929050565b60008060208385031215614f2c57600080fd5b823567ffffffffffffffff811115614f4357600080fd5b614f4f85828601614ccc565b90969095509350505050565b600060208284031215614f6d57600080fd5b8151801515811461474057600080fd5b6000806000806000806000806000806104c08b8d031215614f9d57600080fd5b8a35995060208b013567ffffffffffffffff80821115614fbc57600080fd5b614fc88e838f01614d18565b909b50995060408d0135915080821115614fe157600080fd5b614fed8e838f01614d18565b909950975087915061500160608e01614ca3565b96508d609f8e011261501257600080fd5b60405191506103e08201828110828211171561503057615030615910565b604052508060808d016104608e018f81111561504b57600080fd5b60005b601f8110156150755761506083614ca3565b8452602093840193929092019160010161504e565b5083975061508281614dff565b9650505050506104808b013591506104a08b013590509295989b9194979a5092959850565b600080600083850360a08112156150bd57600080fd5b843567ffffffffffffffff8111156150d457600080fd5b6150e087828801614d18565b90955093505060807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08201121561511657600080fd5b506020840190509250925092565b600080600060a0848603121561513957600080fd5b833567ffffffffffffffff81111561515057600080fd5b61515c86828701614d18565b9094509250614e7f90508560208601614d5a565b60006080828403121561518257600080fd5b6147408383614d5a565b60006020828403121561519e57600080fd5b5051919050565b6000602082840312156151b757600080fd5b61474082614dd3565b600080600080608085870312156151d657600080fd5b6151df85614dd3565b93506020850135925060408501356151f68161593f565b915060608501356152068161593f565b939692955090935050565b600080600080600080600060e0888a03121561522c57600080fd5b61523588614dd3565b965061524360208901614dd3565b9550604088013594506060880135935061525f60808901614dd3565b925061526d60a08901614dd3565b915060c088013561527d8161593f565b8091505092959891949750929550565b60006020828403121561529f57600080fd5b61474082614de7565b600080604083850312156152bb57600080fd5b6152c483614de7565b91506152d260208401614ca3565b90509250929050565b600080600080600060a086880312156152f357600080fd5b6152fc86614e10565b945060208601519350604086015192506060860151915061531f60808701614e10565b90509295509295909350565b60006020828403121561533d57600080fd5b81516147408161593f565b600081518084526020808501945080840160005b8381101561538e57815173ffffffffffffffffffffffffffffffffffffffff168752958201959082019060010161535c565b509495945050505050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b6040808252810183905260008460608301825b868110156154305773ffffffffffffffffffffffffffffffffffffffff61541b84614ca3565b168252602092830192909101906001016153f5565b50809250505073ffffffffffffffffffffffffffffffffffffffff83166020830152949350505050565b6020815260006147406020830184615348565b858152606060208201526000615487606083018688615399565b828103604084015261549a818587615399565b98975050505050505050565b60a0815260006154ba60a083018587615399565b905067ffffffffffffffff6154ce84614de7565b16602083015273ffffffffffffffffffffffffffffffffffffffff6154f560208501614ca3565b16604083015263ffffffff61550c60408501614dd3565b16606083015260608301356080830152949350505050565b602081016003831061555f577f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b91905290565b60006101208201905067ffffffffffffffff835116825273ffffffffffffffffffffffffffffffffffffffff602084015116602083015260408301516155b3604084018263ffffffff169052565b506060830151606083015260808301516155e5608084018273ffffffffffffffffffffffffffffffffffffffff169052565b5060a083015161560560a08401826bffffffffffffffffffffffff169052565b5060c083015161562560c08401826bffffffffffffffffffffffff169052565b5060e083015161564560e08401826bffffffffffffffffffffffff169052565b5061010092830151919092015290565b63ffffffff831681526040602082015260006156746040830184615348565b949350505050565b6bffffffffffffffffffffffff8416815273ffffffffffffffffffffffffffffffffffffffff831660208201526060604082015260006156bf6060830184615348565b95945050505050565b600082198211156156db576156db615854565b500190565b600067ffffffffffffffff80831681851680830382111561570357615703615854565b01949350505050565b60006bffffffffffffffffffffffff80831681851680830382111561570357615703615854565b60008261574257615742615883565b500490565b60006bffffffffffffffffffffffff8084168061576657615766615883565b92169190910492915050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156157aa576157aa615854565b500290565b6000828210156157c1576157c1615854565b500390565b60006bffffffffffffffffffffffff838116908316818110156157eb576157eb615854565b039392505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82141561582557615825615854565b5060010190565b600067ffffffffffffffff8083168181141561584a5761584a615854565b6001019392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6bffffffffffffffffffffffff81168114613c5c57600080fdfea164736f6c6343000806000a496e697469616c697a61626c653a20636f6e7472616374206973206e6f742069",
}

var OCR2DRRegistryABI = OCR2DRRegistryMetaData.ABI

var OCR2DRRegistryBin = OCR2DRRegistryMetaData.Bin

func DeployOCR2DRRegistry(auth *bind.TransactOpts, backend bind.ContractBackend, link common.Address, linkEthFeed common.Address, oracle common.Address) (common.Address, *types.Transaction, *OCR2DRRegistry, error) {
	parsed, err := OCR2DRRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OCR2DRRegistryBin), backend, link, linkEthFeed, oracle)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OCR2DRRegistry{OCR2DRRegistryCaller: OCR2DRRegistryCaller{contract: contract}, OCR2DRRegistryTransactor: OCR2DRRegistryTransactor{contract: contract}, OCR2DRRegistryFilterer: OCR2DRRegistryFilterer{contract: contract}}, nil
}

type OCR2DRRegistry struct {
	address common.Address
	abi     abi.ABI
	OCR2DRRegistryCaller
	OCR2DRRegistryTransactor
	OCR2DRRegistryFilterer
}

type OCR2DRRegistryCaller struct {
	contract *bind.BoundContract
}

type OCR2DRRegistryTransactor struct {
	contract *bind.BoundContract
}

type OCR2DRRegistryFilterer struct {
	contract *bind.BoundContract
}

type OCR2DRRegistrySession struct {
	Contract     *OCR2DRRegistry
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type OCR2DRRegistryCallerSession struct {
	Contract *OCR2DRRegistryCaller
	CallOpts bind.CallOpts
}

type OCR2DRRegistryTransactorSession struct {
	Contract     *OCR2DRRegistryTransactor
	TransactOpts bind.TransactOpts
}

type OCR2DRRegistryRaw struct {
	Contract *OCR2DRRegistry
}

type OCR2DRRegistryCallerRaw struct {
	Contract *OCR2DRRegistryCaller
}

type OCR2DRRegistryTransactorRaw struct {
	Contract *OCR2DRRegistryTransactor
}

func NewOCR2DRRegistry(address common.Address, backend bind.ContractBackend) (*OCR2DRRegistry, error) {
	abi, err := abi.JSON(strings.NewReader(OCR2DRRegistryABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindOCR2DRRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OCR2DRRegistry{address: address, abi: abi, OCR2DRRegistryCaller: OCR2DRRegistryCaller{contract: contract}, OCR2DRRegistryTransactor: OCR2DRRegistryTransactor{contract: contract}, OCR2DRRegistryFilterer: OCR2DRRegistryFilterer{contract: contract}}, nil
}

func NewOCR2DRRegistryCaller(address common.Address, caller bind.ContractCaller) (*OCR2DRRegistryCaller, error) {
	contract, err := bindOCR2DRRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OCR2DRRegistryCaller{contract: contract}, nil
}

func NewOCR2DRRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*OCR2DRRegistryTransactor, error) {
	contract, err := bindOCR2DRRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OCR2DRRegistryTransactor{contract: contract}, nil
}

func NewOCR2DRRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*OCR2DRRegistryFilterer, error) {
	contract, err := bindOCR2DRRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OCR2DRRegistryFilterer{contract: contract}, nil
}

func bindOCR2DRRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OCR2DRRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_OCR2DRRegistry *OCR2DRRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR2DRRegistry.Contract.OCR2DRRegistryCaller.contract.Call(opts, result, method, params...)
}

func (_OCR2DRRegistry *OCR2DRRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.OCR2DRRegistryTransactor.contract.Transfer(opts)
}

func (_OCR2DRRegistry *OCR2DRRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.OCR2DRRegistryTransactor.contract.Transact(opts, method, params...)
}

func (_OCR2DRRegistry *OCR2DRRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR2DRRegistry.Contract.contract.Call(opts, result, method, params...)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.contract.Transfer(opts)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.contract.Transact(opts, method, params...)
}

func (_OCR2DRRegistry *OCR2DRRegistryCaller) MAXCONSUMERS(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _OCR2DRRegistry.contract.Call(opts, &out, "MAX_CONSUMERS")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_OCR2DRRegistry *OCR2DRRegistrySession) MAXCONSUMERS() (uint16, error) {
	return _OCR2DRRegistry.Contract.MAXCONSUMERS(&_OCR2DRRegistry.CallOpts)
}

func (_OCR2DRRegistry *OCR2DRRegistryCallerSession) MAXCONSUMERS() (uint16, error) {
	return _OCR2DRRegistry.Contract.MAXCONSUMERS(&_OCR2DRRegistry.CallOpts)
}

func (_OCR2DRRegistry *OCR2DRRegistryCaller) EstimateCost(opts *bind.CallOpts, gasLimit uint32, gasPrice *big.Int, donFee *big.Int, registryFee *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _OCR2DRRegistry.contract.Call(opts, &out, "estimateCost", gasLimit, gasPrice, donFee, registryFee)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_OCR2DRRegistry *OCR2DRRegistrySession) EstimateCost(gasLimit uint32, gasPrice *big.Int, donFee *big.Int, registryFee *big.Int) (*big.Int, error) {
	return _OCR2DRRegistry.Contract.EstimateCost(&_OCR2DRRegistry.CallOpts, gasLimit, gasPrice, donFee, registryFee)
}

func (_OCR2DRRegistry *OCR2DRRegistryCallerSession) EstimateCost(gasLimit uint32, gasPrice *big.Int, donFee *big.Int, registryFee *big.Int) (*big.Int, error) {
	return _OCR2DRRegistry.Contract.EstimateCost(&_OCR2DRRegistry.CallOpts, gasLimit, gasPrice, donFee, registryFee)
}

func (_OCR2DRRegistry *OCR2DRRegistryCaller) GetAuthorizedSenders(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _OCR2DRRegistry.contract.Call(opts, &out, "getAuthorizedSenders")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_OCR2DRRegistry *OCR2DRRegistrySession) GetAuthorizedSenders() ([]common.Address, error) {
	return _OCR2DRRegistry.Contract.GetAuthorizedSenders(&_OCR2DRRegistry.CallOpts)
}

func (_OCR2DRRegistry *OCR2DRRegistryCallerSession) GetAuthorizedSenders() ([]common.Address, error) {
	return _OCR2DRRegistry.Contract.GetAuthorizedSenders(&_OCR2DRRegistry.CallOpts)
}

func (_OCR2DRRegistry *OCR2DRRegistryCaller) GetConfig(opts *bind.CallOpts) (GetConfig,

	error) {
	var out []interface{}
	err := _OCR2DRRegistry.contract.Call(opts, &out, "getConfig")

	outstruct := new(GetConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.MaxGasLimit = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.StalenessSeconds = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.GasAfterPaymentCalculation = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.FallbackWeiPerUnitLink = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.GasOverhead = *abi.ConvertType(out[4], new(uint32)).(*uint32)
	outstruct.LinkAddress = *abi.ConvertType(out[5], new(common.Address)).(*common.Address)
	outstruct.LinkPriceFeed = *abi.ConvertType(out[6], new(common.Address)).(*common.Address)

	return *outstruct, err

}

func (_OCR2DRRegistry *OCR2DRRegistrySession) GetConfig() (GetConfig,

	error) {
	return _OCR2DRRegistry.Contract.GetConfig(&_OCR2DRRegistry.CallOpts)
}

func (_OCR2DRRegistry *OCR2DRRegistryCallerSession) GetConfig() (GetConfig,

	error) {
	return _OCR2DRRegistry.Contract.GetConfig(&_OCR2DRRegistry.CallOpts)
}

func (_OCR2DRRegistry *OCR2DRRegistryCaller) GetCurrentsubscriptionId(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _OCR2DRRegistry.contract.Call(opts, &out, "getCurrentsubscriptionId")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_OCR2DRRegistry *OCR2DRRegistrySession) GetCurrentsubscriptionId() (uint64, error) {
	return _OCR2DRRegistry.Contract.GetCurrentsubscriptionId(&_OCR2DRRegistry.CallOpts)
}

func (_OCR2DRRegistry *OCR2DRRegistryCallerSession) GetCurrentsubscriptionId() (uint64, error) {
	return _OCR2DRRegistry.Contract.GetCurrentsubscriptionId(&_OCR2DRRegistry.CallOpts)
}

func (_OCR2DRRegistry *OCR2DRRegistryCaller) GetRequestConfig(opts *bind.CallOpts) (uint32, []common.Address, error) {
	var out []interface{}
	err := _OCR2DRRegistry.contract.Call(opts, &out, "getRequestConfig")

	if err != nil {
		return *new(uint32), *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)
	out1 := *abi.ConvertType(out[1], new([]common.Address)).(*[]common.Address)

	return out0, out1, err

}

func (_OCR2DRRegistry *OCR2DRRegistrySession) GetRequestConfig() (uint32, []common.Address, error) {
	return _OCR2DRRegistry.Contract.GetRequestConfig(&_OCR2DRRegistry.CallOpts)
}

func (_OCR2DRRegistry *OCR2DRRegistryCallerSession) GetRequestConfig() (uint32, []common.Address, error) {
	return _OCR2DRRegistry.Contract.GetRequestConfig(&_OCR2DRRegistry.CallOpts)
}

func (_OCR2DRRegistry *OCR2DRRegistryCaller) GetRequiredFee(opts *bind.CallOpts, arg0 []byte, arg1 FunctionsBillingRegistryInterfaceRequestBilling) (*big.Int, error) {
	var out []interface{}
	err := _OCR2DRRegistry.contract.Call(opts, &out, "getRequiredFee", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_OCR2DRRegistry *OCR2DRRegistrySession) GetRequiredFee(arg0 []byte, arg1 FunctionsBillingRegistryInterfaceRequestBilling) (*big.Int, error) {
	return _OCR2DRRegistry.Contract.GetRequiredFee(&_OCR2DRRegistry.CallOpts, arg0, arg1)
}

func (_OCR2DRRegistry *OCR2DRRegistryCallerSession) GetRequiredFee(arg0 []byte, arg1 FunctionsBillingRegistryInterfaceRequestBilling) (*big.Int, error) {
	return _OCR2DRRegistry.Contract.GetRequiredFee(&_OCR2DRRegistry.CallOpts, arg0, arg1)
}

func (_OCR2DRRegistry *OCR2DRRegistryCaller) GetSubscription(opts *bind.CallOpts, subscriptionId uint64) (GetSubscription,

	error) {
	var out []interface{}
	err := _OCR2DRRegistry.contract.Call(opts, &out, "getSubscription", subscriptionId)

	outstruct := new(GetSubscription)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Balance = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Owner = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.Consumers = *abi.ConvertType(out[2], new([]common.Address)).(*[]common.Address)

	return *outstruct, err

}

func (_OCR2DRRegistry *OCR2DRRegistrySession) GetSubscription(subscriptionId uint64) (GetSubscription,

	error) {
	return _OCR2DRRegistry.Contract.GetSubscription(&_OCR2DRRegistry.CallOpts, subscriptionId)
}

func (_OCR2DRRegistry *OCR2DRRegistryCallerSession) GetSubscription(subscriptionId uint64) (GetSubscription,

	error) {
	return _OCR2DRRegistry.Contract.GetSubscription(&_OCR2DRRegistry.CallOpts, subscriptionId)
}

func (_OCR2DRRegistry *OCR2DRRegistryCaller) GetSubscriptionOwner(opts *bind.CallOpts, subscriptionId uint64) (common.Address, error) {
	var out []interface{}
	err := _OCR2DRRegistry.contract.Call(opts, &out, "getSubscriptionOwner", subscriptionId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_OCR2DRRegistry *OCR2DRRegistrySession) GetSubscriptionOwner(subscriptionId uint64) (common.Address, error) {
	return _OCR2DRRegistry.Contract.GetSubscriptionOwner(&_OCR2DRRegistry.CallOpts, subscriptionId)
}

func (_OCR2DRRegistry *OCR2DRRegistryCallerSession) GetSubscriptionOwner(subscriptionId uint64) (common.Address, error) {
	return _OCR2DRRegistry.Contract.GetSubscriptionOwner(&_OCR2DRRegistry.CallOpts, subscriptionId)
}

func (_OCR2DRRegistry *OCR2DRRegistryCaller) GetTotalBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OCR2DRRegistry.contract.Call(opts, &out, "getTotalBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_OCR2DRRegistry *OCR2DRRegistrySession) GetTotalBalance() (*big.Int, error) {
	return _OCR2DRRegistry.Contract.GetTotalBalance(&_OCR2DRRegistry.CallOpts)
}

func (_OCR2DRRegistry *OCR2DRRegistryCallerSession) GetTotalBalance() (*big.Int, error) {
	return _OCR2DRRegistry.Contract.GetTotalBalance(&_OCR2DRRegistry.CallOpts)
}

func (_OCR2DRRegistry *OCR2DRRegistryCaller) IsAuthorizedSender(opts *bind.CallOpts, sender common.Address) (bool, error) {
	var out []interface{}
	err := _OCR2DRRegistry.contract.Call(opts, &out, "isAuthorizedSender", sender)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_OCR2DRRegistry *OCR2DRRegistrySession) IsAuthorizedSender(sender common.Address) (bool, error) {
	return _OCR2DRRegistry.Contract.IsAuthorizedSender(&_OCR2DRRegistry.CallOpts, sender)
}

func (_OCR2DRRegistry *OCR2DRRegistryCallerSession) IsAuthorizedSender(sender common.Address) (bool, error) {
	return _OCR2DRRegistry.Contract.IsAuthorizedSender(&_OCR2DRRegistry.CallOpts, sender)
}

func (_OCR2DRRegistry *OCR2DRRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OCR2DRRegistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_OCR2DRRegistry *OCR2DRRegistrySession) Owner() (common.Address, error) {
	return _OCR2DRRegistry.Contract.Owner(&_OCR2DRRegistry.CallOpts)
}

func (_OCR2DRRegistry *OCR2DRRegistryCallerSession) Owner() (common.Address, error) {
	return _OCR2DRRegistry.Contract.Owner(&_OCR2DRRegistry.CallOpts)
}

func (_OCR2DRRegistry *OCR2DRRegistryCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _OCR2DRRegistry.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_OCR2DRRegistry *OCR2DRRegistrySession) Paused() (bool, error) {
	return _OCR2DRRegistry.Contract.Paused(&_OCR2DRRegistry.CallOpts)
}

func (_OCR2DRRegistry *OCR2DRRegistryCallerSession) Paused() (bool, error) {
	return _OCR2DRRegistry.Contract.Paused(&_OCR2DRRegistry.CallOpts)
}

func (_OCR2DRRegistry *OCR2DRRegistryCaller) PendingRequestExists(opts *bind.CallOpts, subscriptionId uint64) (bool, error) {
	var out []interface{}
	err := _OCR2DRRegistry.contract.Call(opts, &out, "pendingRequestExists", subscriptionId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_OCR2DRRegistry *OCR2DRRegistrySession) PendingRequestExists(subscriptionId uint64) (bool, error) {
	return _OCR2DRRegistry.Contract.PendingRequestExists(&_OCR2DRRegistry.CallOpts, subscriptionId)
}

func (_OCR2DRRegistry *OCR2DRRegistryCallerSession) PendingRequestExists(subscriptionId uint64) (bool, error) {
	return _OCR2DRRegistry.Contract.PendingRequestExists(&_OCR2DRRegistry.CallOpts, subscriptionId)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2DRRegistry.contract.Transact(opts, "acceptOwnership")
}

func (_OCR2DRRegistry *OCR2DRRegistrySession) AcceptOwnership() (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.AcceptOwnership(&_OCR2DRRegistry.TransactOpts)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.AcceptOwnership(&_OCR2DRRegistry.TransactOpts)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactor) AcceptSubscriptionOwnerTransfer(opts *bind.TransactOpts, subscriptionId uint64) (*types.Transaction, error) {
	return _OCR2DRRegistry.contract.Transact(opts, "acceptSubscriptionOwnerTransfer", subscriptionId)
}

func (_OCR2DRRegistry *OCR2DRRegistrySession) AcceptSubscriptionOwnerTransfer(subscriptionId uint64) (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.AcceptSubscriptionOwnerTransfer(&_OCR2DRRegistry.TransactOpts, subscriptionId)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactorSession) AcceptSubscriptionOwnerTransfer(subscriptionId uint64) (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.AcceptSubscriptionOwnerTransfer(&_OCR2DRRegistry.TransactOpts, subscriptionId)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactor) AddConsumer(opts *bind.TransactOpts, subscriptionId uint64, consumer common.Address) (*types.Transaction, error) {
	return _OCR2DRRegistry.contract.Transact(opts, "addConsumer", subscriptionId, consumer)
}

func (_OCR2DRRegistry *OCR2DRRegistrySession) AddConsumer(subscriptionId uint64, consumer common.Address) (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.AddConsumer(&_OCR2DRRegistry.TransactOpts, subscriptionId, consumer)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactorSession) AddConsumer(subscriptionId uint64, consumer common.Address) (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.AddConsumer(&_OCR2DRRegistry.TransactOpts, subscriptionId, consumer)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactor) CancelSubscription(opts *bind.TransactOpts, subscriptionId uint64, to common.Address) (*types.Transaction, error) {
	return _OCR2DRRegistry.contract.Transact(opts, "cancelSubscription", subscriptionId, to)
}

func (_OCR2DRRegistry *OCR2DRRegistrySession) CancelSubscription(subscriptionId uint64, to common.Address) (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.CancelSubscription(&_OCR2DRRegistry.TransactOpts, subscriptionId, to)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactorSession) CancelSubscription(subscriptionId uint64, to common.Address) (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.CancelSubscription(&_OCR2DRRegistry.TransactOpts, subscriptionId, to)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactor) CreateSubscription(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2DRRegistry.contract.Transact(opts, "createSubscription")
}

func (_OCR2DRRegistry *OCR2DRRegistrySession) CreateSubscription() (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.CreateSubscription(&_OCR2DRRegistry.TransactOpts)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactorSession) CreateSubscription() (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.CreateSubscription(&_OCR2DRRegistry.TransactOpts)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactor) FulfillAndBill(opts *bind.TransactOpts, requestId [32]byte, response []byte, err []byte, transmitter common.Address, signers [31]common.Address, signerCount uint8, reportValidationGas *big.Int, initialGas *big.Int) (*types.Transaction, error) {
	return _OCR2DRRegistry.contract.Transact(opts, "fulfillAndBill", requestId, response, err, transmitter, signers, signerCount, reportValidationGas, initialGas)
}

func (_OCR2DRRegistry *OCR2DRRegistrySession) FulfillAndBill(requestId [32]byte, response []byte, err []byte, transmitter common.Address, signers [31]common.Address, signerCount uint8, reportValidationGas *big.Int, initialGas *big.Int) (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.FulfillAndBill(&_OCR2DRRegistry.TransactOpts, requestId, response, err, transmitter, signers, signerCount, reportValidationGas, initialGas)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactorSession) FulfillAndBill(requestId [32]byte, response []byte, err []byte, transmitter common.Address, signers [31]common.Address, signerCount uint8, reportValidationGas *big.Int, initialGas *big.Int) (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.FulfillAndBill(&_OCR2DRRegistry.TransactOpts, requestId, response, err, transmitter, signers, signerCount, reportValidationGas, initialGas)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactor) Initialize(opts *bind.TransactOpts, link common.Address, linkEthFeed common.Address, oracle common.Address) (*types.Transaction, error) {
	return _OCR2DRRegistry.contract.Transact(opts, "initialize", link, linkEthFeed, oracle)
}

func (_OCR2DRRegistry *OCR2DRRegistrySession) Initialize(link common.Address, linkEthFeed common.Address, oracle common.Address) (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.Initialize(&_OCR2DRRegistry.TransactOpts, link, linkEthFeed, oracle)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactorSession) Initialize(link common.Address, linkEthFeed common.Address, oracle common.Address) (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.Initialize(&_OCR2DRRegistry.TransactOpts, link, linkEthFeed, oracle)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactor) OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _OCR2DRRegistry.contract.Transact(opts, "onTokenTransfer", arg0, amount, data)
}

func (_OCR2DRRegistry *OCR2DRRegistrySession) OnTokenTransfer(arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.OnTokenTransfer(&_OCR2DRRegistry.TransactOpts, arg0, amount, data)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactorSession) OnTokenTransfer(arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.OnTokenTransfer(&_OCR2DRRegistry.TransactOpts, arg0, amount, data)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactor) OracleWithdraw(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _OCR2DRRegistry.contract.Transact(opts, "oracleWithdraw", recipient, amount)
}

func (_OCR2DRRegistry *OCR2DRRegistrySession) OracleWithdraw(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.OracleWithdraw(&_OCR2DRRegistry.TransactOpts, recipient, amount)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactorSession) OracleWithdraw(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.OracleWithdraw(&_OCR2DRRegistry.TransactOpts, recipient, amount)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactor) OwnerCancelSubscription(opts *bind.TransactOpts, subscriptionId uint64) (*types.Transaction, error) {
	return _OCR2DRRegistry.contract.Transact(opts, "ownerCancelSubscription", subscriptionId)
}

func (_OCR2DRRegistry *OCR2DRRegistrySession) OwnerCancelSubscription(subscriptionId uint64) (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.OwnerCancelSubscription(&_OCR2DRRegistry.TransactOpts, subscriptionId)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactorSession) OwnerCancelSubscription(subscriptionId uint64) (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.OwnerCancelSubscription(&_OCR2DRRegistry.TransactOpts, subscriptionId)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2DRRegistry.contract.Transact(opts, "pause")
}

func (_OCR2DRRegistry *OCR2DRRegistrySession) Pause() (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.Pause(&_OCR2DRRegistry.TransactOpts)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactorSession) Pause() (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.Pause(&_OCR2DRRegistry.TransactOpts)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactor) RecoverFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _OCR2DRRegistry.contract.Transact(opts, "recoverFunds", to)
}

func (_OCR2DRRegistry *OCR2DRRegistrySession) RecoverFunds(to common.Address) (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.RecoverFunds(&_OCR2DRRegistry.TransactOpts, to)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactorSession) RecoverFunds(to common.Address) (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.RecoverFunds(&_OCR2DRRegistry.TransactOpts, to)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactor) RemoveConsumer(opts *bind.TransactOpts, subscriptionId uint64, consumer common.Address) (*types.Transaction, error) {
	return _OCR2DRRegistry.contract.Transact(opts, "removeConsumer", subscriptionId, consumer)
}

func (_OCR2DRRegistry *OCR2DRRegistrySession) RemoveConsumer(subscriptionId uint64, consumer common.Address) (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.RemoveConsumer(&_OCR2DRRegistry.TransactOpts, subscriptionId, consumer)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactorSession) RemoveConsumer(subscriptionId uint64, consumer common.Address) (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.RemoveConsumer(&_OCR2DRRegistry.TransactOpts, subscriptionId, consumer)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactor) RequestSubscriptionOwnerTransfer(opts *bind.TransactOpts, subscriptionId uint64, newOwner common.Address) (*types.Transaction, error) {
	return _OCR2DRRegistry.contract.Transact(opts, "requestSubscriptionOwnerTransfer", subscriptionId, newOwner)
}

func (_OCR2DRRegistry *OCR2DRRegistrySession) RequestSubscriptionOwnerTransfer(subscriptionId uint64, newOwner common.Address) (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.RequestSubscriptionOwnerTransfer(&_OCR2DRRegistry.TransactOpts, subscriptionId, newOwner)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactorSession) RequestSubscriptionOwnerTransfer(subscriptionId uint64, newOwner common.Address) (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.RequestSubscriptionOwnerTransfer(&_OCR2DRRegistry.TransactOpts, subscriptionId, newOwner)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactor) SetAuthorizedSenders(opts *bind.TransactOpts, senders []common.Address) (*types.Transaction, error) {
	return _OCR2DRRegistry.contract.Transact(opts, "setAuthorizedSenders", senders)
}

func (_OCR2DRRegistry *OCR2DRRegistrySession) SetAuthorizedSenders(senders []common.Address) (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.SetAuthorizedSenders(&_OCR2DRRegistry.TransactOpts, senders)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactorSession) SetAuthorizedSenders(senders []common.Address) (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.SetAuthorizedSenders(&_OCR2DRRegistry.TransactOpts, senders)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactor) SetConfig(opts *bind.TransactOpts, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation *big.Int, fallbackWeiPerUnitLink *big.Int, gasOverhead uint32, requestTimeoutSeconds uint32, fee *big.Int) (*types.Transaction, error) {
	return _OCR2DRRegistry.contract.Transact(opts, "setConfig", maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, gasOverhead, requestTimeoutSeconds, fee)
}

func (_OCR2DRRegistry *OCR2DRRegistrySession) SetConfig(maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation *big.Int, fallbackWeiPerUnitLink *big.Int, gasOverhead uint32, requestTimeoutSeconds uint32, fee *big.Int) (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.SetConfig(&_OCR2DRRegistry.TransactOpts, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, gasOverhead, requestTimeoutSeconds, fee)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactorSession) SetConfig(maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation *big.Int, fallbackWeiPerUnitLink *big.Int, gasOverhead uint32, requestTimeoutSeconds uint32, fee *big.Int) (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.SetConfig(&_OCR2DRRegistry.TransactOpts, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, gasOverhead, requestTimeoutSeconds, fee)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactor) StartBilling(opts *bind.TransactOpts, data []byte, billing FunctionsBillingRegistryInterfaceRequestBilling) (*types.Transaction, error) {
	return _OCR2DRRegistry.contract.Transact(opts, "startBilling", data, billing)
}

func (_OCR2DRRegistry *OCR2DRRegistrySession) StartBilling(data []byte, billing FunctionsBillingRegistryInterfaceRequestBilling) (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.StartBilling(&_OCR2DRRegistry.TransactOpts, data, billing)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactorSession) StartBilling(data []byte, billing FunctionsBillingRegistryInterfaceRequestBilling) (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.StartBilling(&_OCR2DRRegistry.TransactOpts, data, billing)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactor) TimeoutRequests(opts *bind.TransactOpts, requestIdsToTimeout [][32]byte) (*types.Transaction, error) {
	return _OCR2DRRegistry.contract.Transact(opts, "timeoutRequests", requestIdsToTimeout)
}

func (_OCR2DRRegistry *OCR2DRRegistrySession) TimeoutRequests(requestIdsToTimeout [][32]byte) (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.TimeoutRequests(&_OCR2DRRegistry.TransactOpts, requestIdsToTimeout)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactorSession) TimeoutRequests(requestIdsToTimeout [][32]byte) (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.TimeoutRequests(&_OCR2DRRegistry.TransactOpts, requestIdsToTimeout)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _OCR2DRRegistry.contract.Transact(opts, "transferOwnership", to)
}

func (_OCR2DRRegistry *OCR2DRRegistrySession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.TransferOwnership(&_OCR2DRRegistry.TransactOpts, to)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.TransferOwnership(&_OCR2DRRegistry.TransactOpts, to)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2DRRegistry.contract.Transact(opts, "unpause")
}

func (_OCR2DRRegistry *OCR2DRRegistrySession) Unpause() (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.Unpause(&_OCR2DRRegistry.TransactOpts)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactorSession) Unpause() (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.Unpause(&_OCR2DRRegistry.TransactOpts)
}

type OCR2DRRegistryAuthorizedSendersChangedIterator struct {
	Event *OCR2DRRegistryAuthorizedSendersChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2DRRegistryAuthorizedSendersChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2DRRegistryAuthorizedSendersChanged)
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
		it.Event = new(OCR2DRRegistryAuthorizedSendersChanged)
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

func (it *OCR2DRRegistryAuthorizedSendersChangedIterator) Error() error {
	return it.fail
}

func (it *OCR2DRRegistryAuthorizedSendersChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2DRRegistryAuthorizedSendersChanged struct {
	Senders   []common.Address
	ChangedBy common.Address
	Raw       types.Log
}

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) FilterAuthorizedSendersChanged(opts *bind.FilterOpts) (*OCR2DRRegistryAuthorizedSendersChangedIterator, error) {

	logs, sub, err := _OCR2DRRegistry.contract.FilterLogs(opts, "AuthorizedSendersChanged")
	if err != nil {
		return nil, err
	}
	return &OCR2DRRegistryAuthorizedSendersChangedIterator{contract: _OCR2DRRegistry.contract, event: "AuthorizedSendersChanged", logs: logs, sub: sub}, nil
}

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) WatchAuthorizedSendersChanged(opts *bind.WatchOpts, sink chan<- *OCR2DRRegistryAuthorizedSendersChanged) (event.Subscription, error) {

	logs, sub, err := _OCR2DRRegistry.contract.WatchLogs(opts, "AuthorizedSendersChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2DRRegistryAuthorizedSendersChanged)
				if err := _OCR2DRRegistry.contract.UnpackLog(event, "AuthorizedSendersChanged", log); err != nil {
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

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) ParseAuthorizedSendersChanged(log types.Log) (*OCR2DRRegistryAuthorizedSendersChanged, error) {
	event := new(OCR2DRRegistryAuthorizedSendersChanged)
	if err := _OCR2DRRegistry.contract.UnpackLog(event, "AuthorizedSendersChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2DRRegistryBillingEndIterator struct {
	Event *OCR2DRRegistryBillingEnd

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2DRRegistryBillingEndIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2DRRegistryBillingEnd)
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
		it.Event = new(OCR2DRRegistryBillingEnd)
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

func (it *OCR2DRRegistryBillingEndIterator) Error() error {
	return it.fail
}

func (it *OCR2DRRegistryBillingEndIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2DRRegistryBillingEnd struct {
	RequestId          [32]byte
	SubscriptionId     uint64
	SignerPayment      *big.Int
	TransmitterPayment *big.Int
	TotalCost          *big.Int
	Success            bool
	Raw                types.Log
}

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) FilterBillingEnd(opts *bind.FilterOpts, requestId [][32]byte) (*OCR2DRRegistryBillingEndIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _OCR2DRRegistry.contract.FilterLogs(opts, "BillingEnd", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &OCR2DRRegistryBillingEndIterator{contract: _OCR2DRRegistry.contract, event: "BillingEnd", logs: logs, sub: sub}, nil
}

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) WatchBillingEnd(opts *bind.WatchOpts, sink chan<- *OCR2DRRegistryBillingEnd, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _OCR2DRRegistry.contract.WatchLogs(opts, "BillingEnd", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2DRRegistryBillingEnd)
				if err := _OCR2DRRegistry.contract.UnpackLog(event, "BillingEnd", log); err != nil {
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

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) ParseBillingEnd(log types.Log) (*OCR2DRRegistryBillingEnd, error) {
	event := new(OCR2DRRegistryBillingEnd)
	if err := _OCR2DRRegistry.contract.UnpackLog(event, "BillingEnd", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2DRRegistryBillingStartIterator struct {
	Event *OCR2DRRegistryBillingStart

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2DRRegistryBillingStartIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2DRRegistryBillingStart)
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
		it.Event = new(OCR2DRRegistryBillingStart)
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

func (it *OCR2DRRegistryBillingStartIterator) Error() error {
	return it.fail
}

func (it *OCR2DRRegistryBillingStartIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2DRRegistryBillingStart struct {
	RequestId  [32]byte
	Commitment FunctionsBillingRegistryCommitment
	Raw        types.Log
}

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) FilterBillingStart(opts *bind.FilterOpts, requestId [][32]byte) (*OCR2DRRegistryBillingStartIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _OCR2DRRegistry.contract.FilterLogs(opts, "BillingStart", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &OCR2DRRegistryBillingStartIterator{contract: _OCR2DRRegistry.contract, event: "BillingStart", logs: logs, sub: sub}, nil
}

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) WatchBillingStart(opts *bind.WatchOpts, sink chan<- *OCR2DRRegistryBillingStart, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _OCR2DRRegistry.contract.WatchLogs(opts, "BillingStart", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2DRRegistryBillingStart)
				if err := _OCR2DRRegistry.contract.UnpackLog(event, "BillingStart", log); err != nil {
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

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) ParseBillingStart(log types.Log) (*OCR2DRRegistryBillingStart, error) {
	event := new(OCR2DRRegistryBillingStart)
	if err := _OCR2DRRegistry.contract.UnpackLog(event, "BillingStart", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2DRRegistryConfigSetIterator struct {
	Event *OCR2DRRegistryConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2DRRegistryConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2DRRegistryConfigSet)
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
		it.Event = new(OCR2DRRegistryConfigSet)
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

func (it *OCR2DRRegistryConfigSetIterator) Error() error {
	return it.fail
}

func (it *OCR2DRRegistryConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2DRRegistryConfigSet struct {
	MaxGasLimit                uint32
	StalenessSeconds           uint32
	GasAfterPaymentCalculation *big.Int
	FallbackWeiPerUnitLink     *big.Int
	GasOverhead                uint32
	Fee                        *big.Int
	Raw                        types.Log
}

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) FilterConfigSet(opts *bind.FilterOpts) (*OCR2DRRegistryConfigSetIterator, error) {

	logs, sub, err := _OCR2DRRegistry.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &OCR2DRRegistryConfigSetIterator{contract: _OCR2DRRegistry.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *OCR2DRRegistryConfigSet) (event.Subscription, error) {

	logs, sub, err := _OCR2DRRegistry.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2DRRegistryConfigSet)
				if err := _OCR2DRRegistry.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) ParseConfigSet(log types.Log) (*OCR2DRRegistryConfigSet, error) {
	event := new(OCR2DRRegistryConfigSet)
	if err := _OCR2DRRegistry.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2DRRegistryFundsRecoveredIterator struct {
	Event *OCR2DRRegistryFundsRecovered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2DRRegistryFundsRecoveredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2DRRegistryFundsRecovered)
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
		it.Event = new(OCR2DRRegistryFundsRecovered)
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

func (it *OCR2DRRegistryFundsRecoveredIterator) Error() error {
	return it.fail
}

func (it *OCR2DRRegistryFundsRecoveredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2DRRegistryFundsRecovered struct {
	To     common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) FilterFundsRecovered(opts *bind.FilterOpts) (*OCR2DRRegistryFundsRecoveredIterator, error) {

	logs, sub, err := _OCR2DRRegistry.contract.FilterLogs(opts, "FundsRecovered")
	if err != nil {
		return nil, err
	}
	return &OCR2DRRegistryFundsRecoveredIterator{contract: _OCR2DRRegistry.contract, event: "FundsRecovered", logs: logs, sub: sub}, nil
}

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) WatchFundsRecovered(opts *bind.WatchOpts, sink chan<- *OCR2DRRegistryFundsRecovered) (event.Subscription, error) {

	logs, sub, err := _OCR2DRRegistry.contract.WatchLogs(opts, "FundsRecovered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2DRRegistryFundsRecovered)
				if err := _OCR2DRRegistry.contract.UnpackLog(event, "FundsRecovered", log); err != nil {
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

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) ParseFundsRecovered(log types.Log) (*OCR2DRRegistryFundsRecovered, error) {
	event := new(OCR2DRRegistryFundsRecovered)
	if err := _OCR2DRRegistry.contract.UnpackLog(event, "FundsRecovered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2DRRegistryInitializedIterator struct {
	Event *OCR2DRRegistryInitialized

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2DRRegistryInitializedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2DRRegistryInitialized)
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
		it.Event = new(OCR2DRRegistryInitialized)
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

func (it *OCR2DRRegistryInitializedIterator) Error() error {
	return it.fail
}

func (it *OCR2DRRegistryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2DRRegistryInitialized struct {
	Version uint8
	Raw     types.Log
}

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) FilterInitialized(opts *bind.FilterOpts) (*OCR2DRRegistryInitializedIterator, error) {

	logs, sub, err := _OCR2DRRegistry.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &OCR2DRRegistryInitializedIterator{contract: _OCR2DRRegistry.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *OCR2DRRegistryInitialized) (event.Subscription, error) {

	logs, sub, err := _OCR2DRRegistry.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2DRRegistryInitialized)
				if err := _OCR2DRRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
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

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) ParseInitialized(log types.Log) (*OCR2DRRegistryInitialized, error) {
	event := new(OCR2DRRegistryInitialized)
	if err := _OCR2DRRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2DRRegistryOwnershipTransferRequestedIterator struct {
	Event *OCR2DRRegistryOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2DRRegistryOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2DRRegistryOwnershipTransferRequested)
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
		it.Event = new(OCR2DRRegistryOwnershipTransferRequested)
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

func (it *OCR2DRRegistryOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *OCR2DRRegistryOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2DRRegistryOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OCR2DRRegistryOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR2DRRegistry.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OCR2DRRegistryOwnershipTransferRequestedIterator{contract: _OCR2DRRegistry.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OCR2DRRegistryOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR2DRRegistry.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2DRRegistryOwnershipTransferRequested)
				if err := _OCR2DRRegistry.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) ParseOwnershipTransferRequested(log types.Log) (*OCR2DRRegistryOwnershipTransferRequested, error) {
	event := new(OCR2DRRegistryOwnershipTransferRequested)
	if err := _OCR2DRRegistry.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2DRRegistryOwnershipTransferredIterator struct {
	Event *OCR2DRRegistryOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2DRRegistryOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2DRRegistryOwnershipTransferred)
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
		it.Event = new(OCR2DRRegistryOwnershipTransferred)
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

func (it *OCR2DRRegistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *OCR2DRRegistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2DRRegistryOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OCR2DRRegistryOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR2DRRegistry.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OCR2DRRegistryOwnershipTransferredIterator{contract: _OCR2DRRegistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OCR2DRRegistryOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR2DRRegistry.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2DRRegistryOwnershipTransferred)
				if err := _OCR2DRRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) ParseOwnershipTransferred(log types.Log) (*OCR2DRRegistryOwnershipTransferred, error) {
	event := new(OCR2DRRegistryOwnershipTransferred)
	if err := _OCR2DRRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2DRRegistryPausedIterator struct {
	Event *OCR2DRRegistryPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2DRRegistryPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2DRRegistryPaused)
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
		it.Event = new(OCR2DRRegistryPaused)
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

func (it *OCR2DRRegistryPausedIterator) Error() error {
	return it.fail
}

func (it *OCR2DRRegistryPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2DRRegistryPaused struct {
	Account common.Address
	Raw     types.Log
}

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) FilterPaused(opts *bind.FilterOpts) (*OCR2DRRegistryPausedIterator, error) {

	logs, sub, err := _OCR2DRRegistry.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &OCR2DRRegistryPausedIterator{contract: _OCR2DRRegistry.contract, event: "Paused", logs: logs, sub: sub}, nil
}

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *OCR2DRRegistryPaused) (event.Subscription, error) {

	logs, sub, err := _OCR2DRRegistry.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2DRRegistryPaused)
				if err := _OCR2DRRegistry.contract.UnpackLog(event, "Paused", log); err != nil {
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

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) ParsePaused(log types.Log) (*OCR2DRRegistryPaused, error) {
	event := new(OCR2DRRegistryPaused)
	if err := _OCR2DRRegistry.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2DRRegistryRequestTimedOutIterator struct {
	Event *OCR2DRRegistryRequestTimedOut

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2DRRegistryRequestTimedOutIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2DRRegistryRequestTimedOut)
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
		it.Event = new(OCR2DRRegistryRequestTimedOut)
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

func (it *OCR2DRRegistryRequestTimedOutIterator) Error() error {
	return it.fail
}

func (it *OCR2DRRegistryRequestTimedOutIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2DRRegistryRequestTimedOut struct {
	RequestId [32]byte
	Raw       types.Log
}

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) FilterRequestTimedOut(opts *bind.FilterOpts, requestId [][32]byte) (*OCR2DRRegistryRequestTimedOutIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _OCR2DRRegistry.contract.FilterLogs(opts, "RequestTimedOut", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &OCR2DRRegistryRequestTimedOutIterator{contract: _OCR2DRRegistry.contract, event: "RequestTimedOut", logs: logs, sub: sub}, nil
}

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) WatchRequestTimedOut(opts *bind.WatchOpts, sink chan<- *OCR2DRRegistryRequestTimedOut, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _OCR2DRRegistry.contract.WatchLogs(opts, "RequestTimedOut", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2DRRegistryRequestTimedOut)
				if err := _OCR2DRRegistry.contract.UnpackLog(event, "RequestTimedOut", log); err != nil {
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

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) ParseRequestTimedOut(log types.Log) (*OCR2DRRegistryRequestTimedOut, error) {
	event := new(OCR2DRRegistryRequestTimedOut)
	if err := _OCR2DRRegistry.contract.UnpackLog(event, "RequestTimedOut", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2DRRegistrySubscriptionCanceledIterator struct {
	Event *OCR2DRRegistrySubscriptionCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2DRRegistrySubscriptionCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2DRRegistrySubscriptionCanceled)
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
		it.Event = new(OCR2DRRegistrySubscriptionCanceled)
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

func (it *OCR2DRRegistrySubscriptionCanceledIterator) Error() error {
	return it.fail
}

func (it *OCR2DRRegistrySubscriptionCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2DRRegistrySubscriptionCanceled struct {
	SubscriptionId uint64
	To             common.Address
	Amount         *big.Int
	Raw            types.Log
}

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) FilterSubscriptionCanceled(opts *bind.FilterOpts, subscriptionId []uint64) (*OCR2DRRegistrySubscriptionCanceledIterator, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _OCR2DRRegistry.contract.FilterLogs(opts, "SubscriptionCanceled", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &OCR2DRRegistrySubscriptionCanceledIterator{contract: _OCR2DRRegistry.contract, event: "SubscriptionCanceled", logs: logs, sub: sub}, nil
}

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) WatchSubscriptionCanceled(opts *bind.WatchOpts, sink chan<- *OCR2DRRegistrySubscriptionCanceled, subscriptionId []uint64) (event.Subscription, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _OCR2DRRegistry.contract.WatchLogs(opts, "SubscriptionCanceled", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2DRRegistrySubscriptionCanceled)
				if err := _OCR2DRRegistry.contract.UnpackLog(event, "SubscriptionCanceled", log); err != nil {
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

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) ParseSubscriptionCanceled(log types.Log) (*OCR2DRRegistrySubscriptionCanceled, error) {
	event := new(OCR2DRRegistrySubscriptionCanceled)
	if err := _OCR2DRRegistry.contract.UnpackLog(event, "SubscriptionCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2DRRegistrySubscriptionConsumerAddedIterator struct {
	Event *OCR2DRRegistrySubscriptionConsumerAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2DRRegistrySubscriptionConsumerAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2DRRegistrySubscriptionConsumerAdded)
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
		it.Event = new(OCR2DRRegistrySubscriptionConsumerAdded)
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

func (it *OCR2DRRegistrySubscriptionConsumerAddedIterator) Error() error {
	return it.fail
}

func (it *OCR2DRRegistrySubscriptionConsumerAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2DRRegistrySubscriptionConsumerAdded struct {
	SubscriptionId uint64
	Consumer       common.Address
	Raw            types.Log
}

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) FilterSubscriptionConsumerAdded(opts *bind.FilterOpts, subscriptionId []uint64) (*OCR2DRRegistrySubscriptionConsumerAddedIterator, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _OCR2DRRegistry.contract.FilterLogs(opts, "SubscriptionConsumerAdded", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &OCR2DRRegistrySubscriptionConsumerAddedIterator{contract: _OCR2DRRegistry.contract, event: "SubscriptionConsumerAdded", logs: logs, sub: sub}, nil
}

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) WatchSubscriptionConsumerAdded(opts *bind.WatchOpts, sink chan<- *OCR2DRRegistrySubscriptionConsumerAdded, subscriptionId []uint64) (event.Subscription, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _OCR2DRRegistry.contract.WatchLogs(opts, "SubscriptionConsumerAdded", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2DRRegistrySubscriptionConsumerAdded)
				if err := _OCR2DRRegistry.contract.UnpackLog(event, "SubscriptionConsumerAdded", log); err != nil {
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

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) ParseSubscriptionConsumerAdded(log types.Log) (*OCR2DRRegistrySubscriptionConsumerAdded, error) {
	event := new(OCR2DRRegistrySubscriptionConsumerAdded)
	if err := _OCR2DRRegistry.contract.UnpackLog(event, "SubscriptionConsumerAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2DRRegistrySubscriptionConsumerRemovedIterator struct {
	Event *OCR2DRRegistrySubscriptionConsumerRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2DRRegistrySubscriptionConsumerRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2DRRegistrySubscriptionConsumerRemoved)
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
		it.Event = new(OCR2DRRegistrySubscriptionConsumerRemoved)
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

func (it *OCR2DRRegistrySubscriptionConsumerRemovedIterator) Error() error {
	return it.fail
}

func (it *OCR2DRRegistrySubscriptionConsumerRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2DRRegistrySubscriptionConsumerRemoved struct {
	SubscriptionId uint64
	Consumer       common.Address
	Raw            types.Log
}

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) FilterSubscriptionConsumerRemoved(opts *bind.FilterOpts, subscriptionId []uint64) (*OCR2DRRegistrySubscriptionConsumerRemovedIterator, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _OCR2DRRegistry.contract.FilterLogs(opts, "SubscriptionConsumerRemoved", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &OCR2DRRegistrySubscriptionConsumerRemovedIterator{contract: _OCR2DRRegistry.contract, event: "SubscriptionConsumerRemoved", logs: logs, sub: sub}, nil
}

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) WatchSubscriptionConsumerRemoved(opts *bind.WatchOpts, sink chan<- *OCR2DRRegistrySubscriptionConsumerRemoved, subscriptionId []uint64) (event.Subscription, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _OCR2DRRegistry.contract.WatchLogs(opts, "SubscriptionConsumerRemoved", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2DRRegistrySubscriptionConsumerRemoved)
				if err := _OCR2DRRegistry.contract.UnpackLog(event, "SubscriptionConsumerRemoved", log); err != nil {
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

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) ParseSubscriptionConsumerRemoved(log types.Log) (*OCR2DRRegistrySubscriptionConsumerRemoved, error) {
	event := new(OCR2DRRegistrySubscriptionConsumerRemoved)
	if err := _OCR2DRRegistry.contract.UnpackLog(event, "SubscriptionConsumerRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2DRRegistrySubscriptionCreatedIterator struct {
	Event *OCR2DRRegistrySubscriptionCreated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2DRRegistrySubscriptionCreatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2DRRegistrySubscriptionCreated)
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
		it.Event = new(OCR2DRRegistrySubscriptionCreated)
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

func (it *OCR2DRRegistrySubscriptionCreatedIterator) Error() error {
	return it.fail
}

func (it *OCR2DRRegistrySubscriptionCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2DRRegistrySubscriptionCreated struct {
	SubscriptionId uint64
	Owner          common.Address
	Raw            types.Log
}

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) FilterSubscriptionCreated(opts *bind.FilterOpts, subscriptionId []uint64) (*OCR2DRRegistrySubscriptionCreatedIterator, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _OCR2DRRegistry.contract.FilterLogs(opts, "SubscriptionCreated", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &OCR2DRRegistrySubscriptionCreatedIterator{contract: _OCR2DRRegistry.contract, event: "SubscriptionCreated", logs: logs, sub: sub}, nil
}

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) WatchSubscriptionCreated(opts *bind.WatchOpts, sink chan<- *OCR2DRRegistrySubscriptionCreated, subscriptionId []uint64) (event.Subscription, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _OCR2DRRegistry.contract.WatchLogs(opts, "SubscriptionCreated", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2DRRegistrySubscriptionCreated)
				if err := _OCR2DRRegistry.contract.UnpackLog(event, "SubscriptionCreated", log); err != nil {
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

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) ParseSubscriptionCreated(log types.Log) (*OCR2DRRegistrySubscriptionCreated, error) {
	event := new(OCR2DRRegistrySubscriptionCreated)
	if err := _OCR2DRRegistry.contract.UnpackLog(event, "SubscriptionCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2DRRegistrySubscriptionFundedIterator struct {
	Event *OCR2DRRegistrySubscriptionFunded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2DRRegistrySubscriptionFundedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2DRRegistrySubscriptionFunded)
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
		it.Event = new(OCR2DRRegistrySubscriptionFunded)
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

func (it *OCR2DRRegistrySubscriptionFundedIterator) Error() error {
	return it.fail
}

func (it *OCR2DRRegistrySubscriptionFundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2DRRegistrySubscriptionFunded struct {
	SubscriptionId uint64
	OldBalance     *big.Int
	NewBalance     *big.Int
	Raw            types.Log
}

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) FilterSubscriptionFunded(opts *bind.FilterOpts, subscriptionId []uint64) (*OCR2DRRegistrySubscriptionFundedIterator, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _OCR2DRRegistry.contract.FilterLogs(opts, "SubscriptionFunded", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &OCR2DRRegistrySubscriptionFundedIterator{contract: _OCR2DRRegistry.contract, event: "SubscriptionFunded", logs: logs, sub: sub}, nil
}

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) WatchSubscriptionFunded(opts *bind.WatchOpts, sink chan<- *OCR2DRRegistrySubscriptionFunded, subscriptionId []uint64) (event.Subscription, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _OCR2DRRegistry.contract.WatchLogs(opts, "SubscriptionFunded", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2DRRegistrySubscriptionFunded)
				if err := _OCR2DRRegistry.contract.UnpackLog(event, "SubscriptionFunded", log); err != nil {
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

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) ParseSubscriptionFunded(log types.Log) (*OCR2DRRegistrySubscriptionFunded, error) {
	event := new(OCR2DRRegistrySubscriptionFunded)
	if err := _OCR2DRRegistry.contract.UnpackLog(event, "SubscriptionFunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2DRRegistrySubscriptionOwnerTransferRequestedIterator struct {
	Event *OCR2DRRegistrySubscriptionOwnerTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2DRRegistrySubscriptionOwnerTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2DRRegistrySubscriptionOwnerTransferRequested)
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
		it.Event = new(OCR2DRRegistrySubscriptionOwnerTransferRequested)
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

func (it *OCR2DRRegistrySubscriptionOwnerTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *OCR2DRRegistrySubscriptionOwnerTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2DRRegistrySubscriptionOwnerTransferRequested struct {
	SubscriptionId uint64
	From           common.Address
	To             common.Address
	Raw            types.Log
}

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) FilterSubscriptionOwnerTransferRequested(opts *bind.FilterOpts, subscriptionId []uint64) (*OCR2DRRegistrySubscriptionOwnerTransferRequestedIterator, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _OCR2DRRegistry.contract.FilterLogs(opts, "SubscriptionOwnerTransferRequested", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &OCR2DRRegistrySubscriptionOwnerTransferRequestedIterator{contract: _OCR2DRRegistry.contract, event: "SubscriptionOwnerTransferRequested", logs: logs, sub: sub}, nil
}

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) WatchSubscriptionOwnerTransferRequested(opts *bind.WatchOpts, sink chan<- *OCR2DRRegistrySubscriptionOwnerTransferRequested, subscriptionId []uint64) (event.Subscription, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _OCR2DRRegistry.contract.WatchLogs(opts, "SubscriptionOwnerTransferRequested", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2DRRegistrySubscriptionOwnerTransferRequested)
				if err := _OCR2DRRegistry.contract.UnpackLog(event, "SubscriptionOwnerTransferRequested", log); err != nil {
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

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) ParseSubscriptionOwnerTransferRequested(log types.Log) (*OCR2DRRegistrySubscriptionOwnerTransferRequested, error) {
	event := new(OCR2DRRegistrySubscriptionOwnerTransferRequested)
	if err := _OCR2DRRegistry.contract.UnpackLog(event, "SubscriptionOwnerTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2DRRegistrySubscriptionOwnerTransferredIterator struct {
	Event *OCR2DRRegistrySubscriptionOwnerTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2DRRegistrySubscriptionOwnerTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2DRRegistrySubscriptionOwnerTransferred)
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
		it.Event = new(OCR2DRRegistrySubscriptionOwnerTransferred)
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

func (it *OCR2DRRegistrySubscriptionOwnerTransferredIterator) Error() error {
	return it.fail
}

func (it *OCR2DRRegistrySubscriptionOwnerTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2DRRegistrySubscriptionOwnerTransferred struct {
	SubscriptionId uint64
	From           common.Address
	To             common.Address
	Raw            types.Log
}

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) FilterSubscriptionOwnerTransferred(opts *bind.FilterOpts, subscriptionId []uint64) (*OCR2DRRegistrySubscriptionOwnerTransferredIterator, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _OCR2DRRegistry.contract.FilterLogs(opts, "SubscriptionOwnerTransferred", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &OCR2DRRegistrySubscriptionOwnerTransferredIterator{contract: _OCR2DRRegistry.contract, event: "SubscriptionOwnerTransferred", logs: logs, sub: sub}, nil
}

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) WatchSubscriptionOwnerTransferred(opts *bind.WatchOpts, sink chan<- *OCR2DRRegistrySubscriptionOwnerTransferred, subscriptionId []uint64) (event.Subscription, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _OCR2DRRegistry.contract.WatchLogs(opts, "SubscriptionOwnerTransferred", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2DRRegistrySubscriptionOwnerTransferred)
				if err := _OCR2DRRegistry.contract.UnpackLog(event, "SubscriptionOwnerTransferred", log); err != nil {
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

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) ParseSubscriptionOwnerTransferred(log types.Log) (*OCR2DRRegistrySubscriptionOwnerTransferred, error) {
	event := new(OCR2DRRegistrySubscriptionOwnerTransferred)
	if err := _OCR2DRRegistry.contract.UnpackLog(event, "SubscriptionOwnerTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2DRRegistryUnpausedIterator struct {
	Event *OCR2DRRegistryUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2DRRegistryUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2DRRegistryUnpaused)
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
		it.Event = new(OCR2DRRegistryUnpaused)
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

func (it *OCR2DRRegistryUnpausedIterator) Error() error {
	return it.fail
}

func (it *OCR2DRRegistryUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2DRRegistryUnpaused struct {
	Account common.Address
	Raw     types.Log
}

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) FilterUnpaused(opts *bind.FilterOpts) (*OCR2DRRegistryUnpausedIterator, error) {

	logs, sub, err := _OCR2DRRegistry.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &OCR2DRRegistryUnpausedIterator{contract: _OCR2DRRegistry.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *OCR2DRRegistryUnpaused) (event.Subscription, error) {

	logs, sub, err := _OCR2DRRegistry.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2DRRegistryUnpaused)
				if err := _OCR2DRRegistry.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

func (_OCR2DRRegistry *OCR2DRRegistryFilterer) ParseUnpaused(log types.Log) (*OCR2DRRegistryUnpaused, error) {
	event := new(OCR2DRRegistryUnpaused)
	if err := _OCR2DRRegistry.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetConfig struct {
	MaxGasLimit                uint32
	StalenessSeconds           uint32
	GasAfterPaymentCalculation *big.Int
	FallbackWeiPerUnitLink     *big.Int
	GasOverhead                uint32
	LinkAddress                common.Address
	LinkPriceFeed              common.Address
}
type GetSubscription struct {
	Balance   *big.Int
	Owner     common.Address
	Consumers []common.Address
}

func (_OCR2DRRegistry *OCR2DRRegistry) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _OCR2DRRegistry.abi.Events["AuthorizedSendersChanged"].ID:
		return _OCR2DRRegistry.ParseAuthorizedSendersChanged(log)
	case _OCR2DRRegistry.abi.Events["BillingEnd"].ID:
		return _OCR2DRRegistry.ParseBillingEnd(log)
	case _OCR2DRRegistry.abi.Events["BillingStart"].ID:
		return _OCR2DRRegistry.ParseBillingStart(log)
	case _OCR2DRRegistry.abi.Events["ConfigSet"].ID:
		return _OCR2DRRegistry.ParseConfigSet(log)
	case _OCR2DRRegistry.abi.Events["FundsRecovered"].ID:
		return _OCR2DRRegistry.ParseFundsRecovered(log)
	case _OCR2DRRegistry.abi.Events["Initialized"].ID:
		return _OCR2DRRegistry.ParseInitialized(log)
	case _OCR2DRRegistry.abi.Events["OwnershipTransferRequested"].ID:
		return _OCR2DRRegistry.ParseOwnershipTransferRequested(log)
	case _OCR2DRRegistry.abi.Events["OwnershipTransferred"].ID:
		return _OCR2DRRegistry.ParseOwnershipTransferred(log)
	case _OCR2DRRegistry.abi.Events["Paused"].ID:
		return _OCR2DRRegistry.ParsePaused(log)
	case _OCR2DRRegistry.abi.Events["RequestTimedOut"].ID:
		return _OCR2DRRegistry.ParseRequestTimedOut(log)
	case _OCR2DRRegistry.abi.Events["SubscriptionCanceled"].ID:
		return _OCR2DRRegistry.ParseSubscriptionCanceled(log)
	case _OCR2DRRegistry.abi.Events["SubscriptionConsumerAdded"].ID:
		return _OCR2DRRegistry.ParseSubscriptionConsumerAdded(log)
	case _OCR2DRRegistry.abi.Events["SubscriptionConsumerRemoved"].ID:
		return _OCR2DRRegistry.ParseSubscriptionConsumerRemoved(log)
	case _OCR2DRRegistry.abi.Events["SubscriptionCreated"].ID:
		return _OCR2DRRegistry.ParseSubscriptionCreated(log)
	case _OCR2DRRegistry.abi.Events["SubscriptionFunded"].ID:
		return _OCR2DRRegistry.ParseSubscriptionFunded(log)
	case _OCR2DRRegistry.abi.Events["SubscriptionOwnerTransferRequested"].ID:
		return _OCR2DRRegistry.ParseSubscriptionOwnerTransferRequested(log)
	case _OCR2DRRegistry.abi.Events["SubscriptionOwnerTransferred"].ID:
		return _OCR2DRRegistry.ParseSubscriptionOwnerTransferred(log)
	case _OCR2DRRegistry.abi.Events["Unpaused"].ID:
		return _OCR2DRRegistry.ParseUnpaused(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (OCR2DRRegistryAuthorizedSendersChanged) Topic() common.Hash {
	return common.HexToHash("0xf263cfb3e4298332e776194610cf9fdc09ccb3ada8b9aa39764d882e11fbf0a0")
}

func (OCR2DRRegistryBillingEnd) Topic() common.Hash {
	return common.HexToHash("0xc8dc973332de19a5f71b6026983110e9c2e04b0c98b87eb771ccb78607fd114f")
}

func (OCR2DRRegistryBillingStart) Topic() common.Hash {
	return common.HexToHash("0x99f7f4e65b4b9fbabd4e357c47ed3099b36e57ecd3a43e84662f34c207d0ebe4")
}

func (OCR2DRRegistryConfigSet) Topic() common.Hash {
	return common.HexToHash("0x1f8b76dfe267648827eba3c1966a7149a45245c5af445ce68ba124d04b7520aa")
}

func (OCR2DRRegistryFundsRecovered) Topic() common.Hash {
	return common.HexToHash("0x59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b436600")
}

func (OCR2DRRegistryInitialized) Topic() common.Hash {
	return common.HexToHash("0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498")
}

func (OCR2DRRegistryOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (OCR2DRRegistryOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (OCR2DRRegistryPaused) Topic() common.Hash {
	return common.HexToHash("0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258")
}

func (OCR2DRRegistryRequestTimedOut) Topic() common.Hash {
	return common.HexToHash("0xf1ca1e9147be737b04a2b018a79405f687a97de8dd8a2559bbe62357343af414")
}

func (OCR2DRRegistrySubscriptionCanceled) Topic() common.Hash {
	return common.HexToHash("0xe8ed5b475a5b5987aa9165e8731bb78043f39eee32ec5a1169a89e27fcd49815")
}

func (OCR2DRRegistrySubscriptionConsumerAdded) Topic() common.Hash {
	return common.HexToHash("0x43dc749a04ac8fb825cbd514f7c0e13f13bc6f2ee66043b76629d51776cff8e0")
}

func (OCR2DRRegistrySubscriptionConsumerRemoved) Topic() common.Hash {
	return common.HexToHash("0x182bff9831466789164ca77075fffd84916d35a8180ba73c27e45634549b445b")
}

func (OCR2DRRegistrySubscriptionCreated) Topic() common.Hash {
	return common.HexToHash("0x464722b4166576d3dcbba877b999bc35cf911f4eaf434b7eba68fa113951d0bf")
}

func (OCR2DRRegistrySubscriptionFunded) Topic() common.Hash {
	return common.HexToHash("0xd39ec07f4e209f627a4c427971473820dc129761ba28de8906bd56f57101d4f8")
}

func (OCR2DRRegistrySubscriptionOwnerTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x69436ea6df009049404f564eff6622cd00522b0bd6a89efd9e52a355c4a879be")
}

func (OCR2DRRegistrySubscriptionOwnerTransferred) Topic() common.Hash {
	return common.HexToHash("0x6f1dc65165ffffedfd8e507b4a0f1fcfdada045ed11f6c26ba27cedfe87802f0")
}

func (OCR2DRRegistryUnpaused) Topic() common.Hash {
	return common.HexToHash("0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa")
}

func (_OCR2DRRegistry *OCR2DRRegistry) Address() common.Address {
	return _OCR2DRRegistry.address
}

type OCR2DRRegistryInterface interface {
	MAXCONSUMERS(opts *bind.CallOpts) (uint16, error)

	EstimateCost(opts *bind.CallOpts, gasLimit uint32, gasPrice *big.Int, donFee *big.Int, registryFee *big.Int) (*big.Int, error)

	GetAuthorizedSenders(opts *bind.CallOpts) ([]common.Address, error)

	GetConfig(opts *bind.CallOpts) (GetConfig,

		error)

	GetCurrentsubscriptionId(opts *bind.CallOpts) (uint64, error)

	GetRequestConfig(opts *bind.CallOpts) (uint32, []common.Address, error)

	GetRequiredFee(opts *bind.CallOpts, arg0 []byte, arg1 FunctionsBillingRegistryInterfaceRequestBilling) (*big.Int, error)

	GetSubscription(opts *bind.CallOpts, subscriptionId uint64) (GetSubscription,

		error)

	GetSubscriptionOwner(opts *bind.CallOpts, subscriptionId uint64) (common.Address, error)

	GetTotalBalance(opts *bind.CallOpts) (*big.Int, error)

	IsAuthorizedSender(opts *bind.CallOpts, sender common.Address) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	Paused(opts *bind.CallOpts) (bool, error)

	PendingRequestExists(opts *bind.CallOpts, subscriptionId uint64) (bool, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AcceptSubscriptionOwnerTransfer(opts *bind.TransactOpts, subscriptionId uint64) (*types.Transaction, error)

	AddConsumer(opts *bind.TransactOpts, subscriptionId uint64, consumer common.Address) (*types.Transaction, error)

	CancelSubscription(opts *bind.TransactOpts, subscriptionId uint64, to common.Address) (*types.Transaction, error)

	CreateSubscription(opts *bind.TransactOpts) (*types.Transaction, error)

	FulfillAndBill(opts *bind.TransactOpts, requestId [32]byte, response []byte, err []byte, transmitter common.Address, signers [31]common.Address, signerCount uint8, reportValidationGas *big.Int, initialGas *big.Int) (*types.Transaction, error)

	Initialize(opts *bind.TransactOpts, link common.Address, linkEthFeed common.Address, oracle common.Address) (*types.Transaction, error)

	OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error)

	OracleWithdraw(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error)

	OwnerCancelSubscription(opts *bind.TransactOpts, subscriptionId uint64) (*types.Transaction, error)

	Pause(opts *bind.TransactOpts) (*types.Transaction, error)

	RecoverFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	RemoveConsumer(opts *bind.TransactOpts, subscriptionId uint64, consumer common.Address) (*types.Transaction, error)

	RequestSubscriptionOwnerTransfer(opts *bind.TransactOpts, subscriptionId uint64, newOwner common.Address) (*types.Transaction, error)

	SetAuthorizedSenders(opts *bind.TransactOpts, senders []common.Address) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation *big.Int, fallbackWeiPerUnitLink *big.Int, gasOverhead uint32, requestTimeoutSeconds uint32, fee *big.Int) (*types.Transaction, error)

	StartBilling(opts *bind.TransactOpts, data []byte, billing FunctionsBillingRegistryInterfaceRequestBilling) (*types.Transaction, error)

	TimeoutRequests(opts *bind.TransactOpts, requestIdsToTimeout [][32]byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Unpause(opts *bind.TransactOpts) (*types.Transaction, error)

	FilterAuthorizedSendersChanged(opts *bind.FilterOpts) (*OCR2DRRegistryAuthorizedSendersChangedIterator, error)

	WatchAuthorizedSendersChanged(opts *bind.WatchOpts, sink chan<- *OCR2DRRegistryAuthorizedSendersChanged) (event.Subscription, error)

	ParseAuthorizedSendersChanged(log types.Log) (*OCR2DRRegistryAuthorizedSendersChanged, error)

	FilterBillingEnd(opts *bind.FilterOpts, requestId [][32]byte) (*OCR2DRRegistryBillingEndIterator, error)

	WatchBillingEnd(opts *bind.WatchOpts, sink chan<- *OCR2DRRegistryBillingEnd, requestId [][32]byte) (event.Subscription, error)

	ParseBillingEnd(log types.Log) (*OCR2DRRegistryBillingEnd, error)

	FilterBillingStart(opts *bind.FilterOpts, requestId [][32]byte) (*OCR2DRRegistryBillingStartIterator, error)

	WatchBillingStart(opts *bind.WatchOpts, sink chan<- *OCR2DRRegistryBillingStart, requestId [][32]byte) (event.Subscription, error)

	ParseBillingStart(log types.Log) (*OCR2DRRegistryBillingStart, error)

	FilterConfigSet(opts *bind.FilterOpts) (*OCR2DRRegistryConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *OCR2DRRegistryConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*OCR2DRRegistryConfigSet, error)

	FilterFundsRecovered(opts *bind.FilterOpts) (*OCR2DRRegistryFundsRecoveredIterator, error)

	WatchFundsRecovered(opts *bind.WatchOpts, sink chan<- *OCR2DRRegistryFundsRecovered) (event.Subscription, error)

	ParseFundsRecovered(log types.Log) (*OCR2DRRegistryFundsRecovered, error)

	FilterInitialized(opts *bind.FilterOpts) (*OCR2DRRegistryInitializedIterator, error)

	WatchInitialized(opts *bind.WatchOpts, sink chan<- *OCR2DRRegistryInitialized) (event.Subscription, error)

	ParseInitialized(log types.Log) (*OCR2DRRegistryInitialized, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OCR2DRRegistryOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OCR2DRRegistryOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*OCR2DRRegistryOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OCR2DRRegistryOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OCR2DRRegistryOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*OCR2DRRegistryOwnershipTransferred, error)

	FilterPaused(opts *bind.FilterOpts) (*OCR2DRRegistryPausedIterator, error)

	WatchPaused(opts *bind.WatchOpts, sink chan<- *OCR2DRRegistryPaused) (event.Subscription, error)

	ParsePaused(log types.Log) (*OCR2DRRegistryPaused, error)

	FilterRequestTimedOut(opts *bind.FilterOpts, requestId [][32]byte) (*OCR2DRRegistryRequestTimedOutIterator, error)

	WatchRequestTimedOut(opts *bind.WatchOpts, sink chan<- *OCR2DRRegistryRequestTimedOut, requestId [][32]byte) (event.Subscription, error)

	ParseRequestTimedOut(log types.Log) (*OCR2DRRegistryRequestTimedOut, error)

	FilterSubscriptionCanceled(opts *bind.FilterOpts, subscriptionId []uint64) (*OCR2DRRegistrySubscriptionCanceledIterator, error)

	WatchSubscriptionCanceled(opts *bind.WatchOpts, sink chan<- *OCR2DRRegistrySubscriptionCanceled, subscriptionId []uint64) (event.Subscription, error)

	ParseSubscriptionCanceled(log types.Log) (*OCR2DRRegistrySubscriptionCanceled, error)

	FilterSubscriptionConsumerAdded(opts *bind.FilterOpts, subscriptionId []uint64) (*OCR2DRRegistrySubscriptionConsumerAddedIterator, error)

	WatchSubscriptionConsumerAdded(opts *bind.WatchOpts, sink chan<- *OCR2DRRegistrySubscriptionConsumerAdded, subscriptionId []uint64) (event.Subscription, error)

	ParseSubscriptionConsumerAdded(log types.Log) (*OCR2DRRegistrySubscriptionConsumerAdded, error)

	FilterSubscriptionConsumerRemoved(opts *bind.FilterOpts, subscriptionId []uint64) (*OCR2DRRegistrySubscriptionConsumerRemovedIterator, error)

	WatchSubscriptionConsumerRemoved(opts *bind.WatchOpts, sink chan<- *OCR2DRRegistrySubscriptionConsumerRemoved, subscriptionId []uint64) (event.Subscription, error)

	ParseSubscriptionConsumerRemoved(log types.Log) (*OCR2DRRegistrySubscriptionConsumerRemoved, error)

	FilterSubscriptionCreated(opts *bind.FilterOpts, subscriptionId []uint64) (*OCR2DRRegistrySubscriptionCreatedIterator, error)

	WatchSubscriptionCreated(opts *bind.WatchOpts, sink chan<- *OCR2DRRegistrySubscriptionCreated, subscriptionId []uint64) (event.Subscription, error)

	ParseSubscriptionCreated(log types.Log) (*OCR2DRRegistrySubscriptionCreated, error)

	FilterSubscriptionFunded(opts *bind.FilterOpts, subscriptionId []uint64) (*OCR2DRRegistrySubscriptionFundedIterator, error)

	WatchSubscriptionFunded(opts *bind.WatchOpts, sink chan<- *OCR2DRRegistrySubscriptionFunded, subscriptionId []uint64) (event.Subscription, error)

	ParseSubscriptionFunded(log types.Log) (*OCR2DRRegistrySubscriptionFunded, error)

	FilterSubscriptionOwnerTransferRequested(opts *bind.FilterOpts, subscriptionId []uint64) (*OCR2DRRegistrySubscriptionOwnerTransferRequestedIterator, error)

	WatchSubscriptionOwnerTransferRequested(opts *bind.WatchOpts, sink chan<- *OCR2DRRegistrySubscriptionOwnerTransferRequested, subscriptionId []uint64) (event.Subscription, error)

	ParseSubscriptionOwnerTransferRequested(log types.Log) (*OCR2DRRegistrySubscriptionOwnerTransferRequested, error)

	FilterSubscriptionOwnerTransferred(opts *bind.FilterOpts, subscriptionId []uint64) (*OCR2DRRegistrySubscriptionOwnerTransferredIterator, error)

	WatchSubscriptionOwnerTransferred(opts *bind.WatchOpts, sink chan<- *OCR2DRRegistrySubscriptionOwnerTransferred, subscriptionId []uint64) (event.Subscription, error)

	ParseSubscriptionOwnerTransferred(log types.Log) (*OCR2DRRegistrySubscriptionOwnerTransferred, error)

	FilterUnpaused(opts *bind.FilterOpts) (*OCR2DRRegistryUnpausedIterator, error)

	WatchUnpaused(opts *bind.WatchOpts, sink chan<- *OCR2DRRegistryUnpaused) (event.Subscription, error)

	ParseUnpaused(log types.Log) (*OCR2DRRegistryUnpaused, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
