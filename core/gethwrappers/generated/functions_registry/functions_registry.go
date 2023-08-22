// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package functions_registry

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

type IFunctionsBillingRegistryRequestBilling struct {
	SubscriptionId uint64
	Client         common.Address
	GasLimit       uint32
	GasPrice       *big.Int
}

var FunctionsRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"linkEthFeed\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"internalBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"externalBalance\",\"type\":\"uint256\"}],\"name\":\"BalanceInvariantViolated\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotSelfTransfer\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptySendersList\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"have\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"want\",\"type\":\"uint32\"}],\"name\":\"GasLimitTooBig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"InvalidConsumer\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"linkWei\",\"type\":\"int256\"}],\"name\":\"InvalidLinkWeiPrice\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSubscription\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposedOwner\",\"type\":\"address\"}],\"name\":\"MustBeRequestedOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"MustBeSubOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAllowedToSetSenders\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotProposedOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableFromLink\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OwnerMustBeSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentTooLarge\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingRequestExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Reentrant\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyConsumers\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedSender\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"senders\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"changedBy\",\"type\":\"address\"}],\"name\":\"AuthorizedSendersChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"signerPayment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"transmitterPayment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalCost\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"BillingEnd\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"don\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"donFee\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"registryFee\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"estimatedCost\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structFunctionsBillingRegistry.Commitment\",\"name\":\"commitment\",\"type\":\"tuple\"}],\"name\":\"BillingStart\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasOverhead\",\"type\":\"uint32\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"RequestTimedOut\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"SubscriptionCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"SubscriptionCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"MAX_CONSUMERS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"acceptSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"addConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"cancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createSubscription\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"donFee\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"registryFee\",\"type\":\"uint96\"}],\"name\":\"estimateCost\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"response\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"internalType\":\"address[31]\",\"name\":\"signers\",\"type\":\"address[31]\"},{\"internalType\":\"uint8\",\"name\":\"signerCount\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"reportValidationGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"initialGas\",\"type\":\"uint256\"}],\"name\":\"fulfillAndBill\",\"outputs\":[{\"internalType\":\"enumIFunctionsBillingRegistry.FulfillResult\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAuthorizedSenders\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfig\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"internalType\":\"uint32\",\"name\":\"gasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"linkAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"linkPriceFeed\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCurrentsubscriptionId\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRequestConfig\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"},{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"}],\"internalType\":\"structIFunctionsBillingRegistry.RequestBilling\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"getRequiredFee\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"getSubscription\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"getSubscriptionOwner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTotalBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"linkEthFeed\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"isAuthorizedSender\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"oracleWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"ownerCancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"pendingRequestExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"removeConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"requestSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"senders\",\"type\":\"address[]\"}],\"name\":\"setAuthorizedSenders\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"internalType\":\"uint32\",\"name\":\"gasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"requestTimeoutSeconds\",\"type\":\"uint32\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"}],\"internalType\":\"structIFunctionsBillingRegistry.RequestBilling\",\"name\":\"billing\",\"type\":\"tuple\"}],\"name\":\"startBilling\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"requestIdsToTimeout\",\"type\":\"bytes32[]\"}],\"name\":\"timeoutRequests\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162005d3838038062005d3883398101604081905262000034916200040c565b620000418383836200004a565b50505062000456565b600054610100900460ff16158080156200006b5750600054600160ff909116105b806200009b57506200008830620001c960201b62003bd41760201c565b1580156200009b575060005460ff166001145b620001045760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b60648201526084015b60405180910390fd5b6000805460ff19166001179055801562000128576000805461ff0019166101001790555b62000132620001d8565b6200013f33600062000240565b606980546001600160a01b038087166001600160a01b031992831617909255606a8054868416908316179055606b8054928516929091169190911790558015620001c3576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b50505050565b6001600160a01b03163b151590565b600054610100900460ff16620002345760405162461bcd60e51b815260206004820152602b602482015260008051602062005d1883398151915260448201526a6e697469616c697a696e6760a81b6064820152608401620000fb565b6200023e62000304565b565b600054610100900460ff166200029c5760405162461bcd60e51b815260206004820152602b602482015260008051602062005d1883398151915260448201526a6e697469616c697a696e6760a81b6064820152608401620000fb565b6001600160a01b038216620002c457604051635b5a8afd60e11b815260040160405180910390fd5b600080546001600160a01b03808516620100000262010000600160b01b031990921691909117909155811615620003005762000300816200036c565b5050565b600054610100900460ff16620003605760405162461bcd60e51b815260206004820152602b602482015260008051602062005d1883398151915260448201526a6e697469616c697a696e6760a81b6064820152608401620000fb565b6034805460ff19169055565b6001600160a01b038116331415620003975760405163282010c360e01b815260040160405180910390fd5b600180546001600160a01b0319166001600160a01b038381169182179092556000805460405192936201000090910416917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b03811681146200040757600080fd5b919050565b6000806000606084860312156200042257600080fd5b6200042d84620003ef565b92506200043d60208501620003ef565b91506200044d60408501620003ef565b90509250925092565b6158b280620004666000396000f3fe608060405234801561001057600080fd5b50600436106102255760003560e01c80638da5cb5b1161012a578063c0c53b8b116100bd578063e82ad7d41161008c578063f1e14a2111610071578063f1e14a2114610561578063f2fde38b14610578578063fa00763a1461058b57600080fd5b8063e82ad7d41461053b578063ee56997b1461054e57600080fd5b8063c0c53b8b1461048d578063c3f909d4146104a0578063d7ae1d3014610515578063e72f6e301461052857600080fd5b8063a47c7696116100f9578063a47c769614610432578063a4c0ed3614610454578063a9d03c0514610467578063b2a489ff1461047a57600080fd5b80638da5cb5b146103a25780639f87fad7146103e7578063a1a6d041146103fa578063a21a23e41461042a57600080fd5b80633f4ba83a116101bd578063665871ec1161018c57806379ba50971161017157806379ba50971461037f57806382359740146103875780638456cb591461039a57600080fd5b8063665871ec146103595780637341c10c1461036c57600080fd5b80633f4ba83a1461030c5780635c975abb1461031457806364d51a2a1461032b57806366316d8d1461034657600080fd5b806312b58349116101f957806312b58349146102915780632408afaa146102bd57806327923e41146102d257806333652e3e146102e557600080fd5b80620122911461022a57806302bcc5b61461024957806304c357cb1461025e5780630739e4f114610271575b600080fd5b61023261059e565b6040516102409291906155a1565b60405180910390f35b61025c6102573660046151d9565b6105bd565b005b61025c61026c3660046151f4565b61063a565b61028461027f366004614edf565b610835565b6040516102409190615470565b606f546801000000000000000090046bffffffffffffffffffffffff165b604051908152602001610240565b6102c5610ef7565b60405161024091906153a6565b61025c6102e0366004615173565b610f66565b606f5467ffffffffffffffff165b60405167ffffffffffffffff9091168152602001610240565b61025c6110b9565b60345460ff165b6040519015158152602001610240565b610333606481565b60405161ffff9091168152602001610240565b61025c610354366004614e44565b6110cb565b61025c610367366004614e7b565b611331565b61025c61037a3660046151f4565b6115fd565b61025c611890565b61025c6103953660046151d9565b611983565b61025c611d04565b60005462010000900473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610240565b61025c6103f53660046151f4565b611d14565b61040d610408366004615122565b61219b565b6040516bffffffffffffffffffffffff9091168152602001610240565b6102f36122bf565b6104456104403660046151d9565b612657565b604051610240939291906155c8565b61025c610462366004614dea565b612788565b6102af610475366004615009565b6129e1565b6103c26104883660046151d9565b6131dc565b61025c61049b366004614da7565b613275565b607354607454607254607554606954606a546040805163ffffffff808916825265010000000000909804881660208201529081019590955260608501939093529316608083015273ffffffffffffffffffffffffffffffffffffffff92831660a08301529190911660c082015260e001610240565b61025c6105233660046151f4565b613477565b61025c610536366004614d8c565b6135de565b61031b6105493660046151d9565b6137fb565b61025c61055c366004614e7b565b613a3a565b61040d61056f366004615086565b60009392505050565b61025c610586366004614d8c565b613bad565b61031b610599366004614d8c565b613bc1565b60735460009060609063ffffffff166105b5610ef7565b915091509091565b6105c5613bf0565b67ffffffffffffffff81166000908152606d602052604090205473ffffffffffffffffffffffffffffffffffffffff168061062c576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6106368282613c47565b5050565b67ffffffffffffffff82166000908152606d6020526040902054829073ffffffffffffffffffffffffffffffffffffffff16806106a3576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff82161461070f576040517fd8a3fb5200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821660048201526024015b60405180910390fd5b607354640100000000900460ff1615610754576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61075c61404e565b67ffffffffffffffff84166000908152606d602052604090206001015473ffffffffffffffffffffffffffffffffffffffff84811691161461082f5767ffffffffffffffff84166000818152606d602090815260409182902060010180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff88169081179091558251338152918201527f69436ea6df009049404f564eff6622cd00522b0bd6a89efd9e52a355c4a879be91015b60405180910390a25b50505050565b600061083f6140bb565b607354640100000000900460ff1615610884576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61088c61404e565b60008b815260716020908152604091829020825161012081018452815467ffffffffffffffff8116825268010000000000000000810473ffffffffffffffffffffffffffffffffffffffff908116948301949094527c0100000000000000000000000000000000000000000000000000000000900463ffffffff169381019390935260018101546060840152600281015491821660808401819052740100000000000000000000000000000000000000009092046bffffffffffffffffffffffff90811660a0850152600382015480821660c08601526c0100000000000000000000000090041660e084015260040154610100830152610990576002915050610ee9565b60008c81526071602052604080822082815560018101839055600281018390556003810180547fffffffffffffffff000000000000000000000000000000000000000000000000169055600401829055517f0ca761750000000000000000000000000000000000000000000000000000000090610a19908f908f908f908f908f906024016153b9565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152918152602080830180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090951694909417909352607380547fffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffff1664010000000017905584015191840151909250600091610ae69163ffffffff90911690846140fa565b607380547fffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffff16905560745460a085015160c0860151929350600092610b3292899290918c908c3a614146565b604080820151865167ffffffffffffffff166000908152606e60205291909120549192506bffffffffffffffffffffffff90811691161015610ba0576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b604080820151855167ffffffffffffffff166000908152606e602052918220805491929091610bde9084906bffffffffffffffffffffffff16615712565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060005b8860ff16811015610cd8578151607060008c84601f8110610c3257610c3261582d565b602002015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282829054906101000a90046bffffffffffffffffffffffff16610c979190615658565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055508080610cd09061573f565b915050610c0f565b508360c0015160706000610d0860005473ffffffffffffffffffffffffffffffffffffffff620100009091041690565b73ffffffffffffffffffffffffffffffffffffffff168152602081019190915260400160009081208054909190610d4e9084906bffffffffffffffffffffffff16615658565b82546101009290920a6bffffffffffffffffffffffff81810219909316918316021790915560208381015173ffffffffffffffffffffffffffffffffffffffff8e166000908152607090925260408220805491945092610db091859116615658565b82546101009290920a6bffffffffffffffffffffffff81810219909316918316021790915560e0860151865167ffffffffffffffff166000908152606e60205260409020805491935091600c91610e199185916c01000000000000000000000000900416615712565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055508e7fc8dc973332de19a5f71b6026983110e9c2e04b0c98b87eb771ccb78607fd114f856000015183600001518460200151856040015187604051610ecb95949392919067ffffffffffffffff9590951685526bffffffffffffffffffffffff9384166020860152918316604085015290911660608301521515608082015260a00190565b60405180910390a281610edf576001610ee2565b60005b9450505050505b9a9950505050505050505050565b60606068805480602002602001604051908101604052809291908181526020018280548015610f5c57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610f31575b5050505050905090565b610f6e613bf0565b60008313610fab576040517f43d4cf6600000000000000000000000000000000000000000000000000000000815260048101849052602401610706565b6040805160c08101825263ffffffff888116808352600060208085019190915289831684860181905260608086018b9052888516608080880182905295891660a0978801819052607380547fffffffffffffffffffffffffffffffffffffffffffffff00000000000000000016871765010000000000860217905560748d9055607580547fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000016831764010000000090920291909117905560728b9055875194855292840191909152948201899052938101879052908101929092527f24d3d934adfef9b9029d6ffa463c07d0139ed47d26ee23506f85ece2879d2bd4910160405180910390a1505050505050565b6110c1613bf0565b6110c96142c6565b565b607354640100000000900460ff1615611110576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61111861404e565b6bffffffffffffffffffffffff811661114b5750336000908152607060205260409020546bffffffffffffffffffffffff165b336000908152607060205260409020546bffffffffffffffffffffffff808316911610156111a5576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b33600090815260706020526040812080548392906111d29084906bffffffffffffffffffffffff16615712565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555080606f60088282829054906101000a90046bffffffffffffffffffffffff166112299190615712565b82546101009290920a6bffffffffffffffffffffffff8181021990931691831602179091556069546040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff868116600483015292851660248201529116915063a9059cbb90604401602060405180830381600087803b1580156112c357600080fd5b505af11580156112d7573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112fb9190614ebd565b610636576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61133961404e565b60005b818110156115f85760008383838181106113585761135861582d565b602090810292909201356000818152607184526040808220815161012081018352815467ffffffffffffffff811680835268010000000000000000820473ffffffffffffffffffffffffffffffffffffffff908116848b01527c010000000000000000000000000000000000000000000000000000000090920463ffffffff168386015260018401546060840152600284015480831660808501527401000000000000000000000000000000000000000090046bffffffffffffffffffffffff90811660a0850152600385015480821660c08601526c0100000000000000000000000090041660e0840152600490930154610100830152918452606d90965291205491945016331490506114cd57805167ffffffffffffffff166000908152606d6020526040908190205490517fd8a3fb5200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602401610706565b60755461010082015142916114f19164010000000090910463ffffffff1690615614565b11156115e35760e0810151815167ffffffffffffffff166000908152606e602052604090208054600c906115449084906c0100000000000000000000000090046bffffffffffffffffffffffff16615712565b82546bffffffffffffffffffffffff9182166101009390930a92830291909202199091161790555060008281526071602052604080822082815560018101839055600281018390556003810180547fffffffffffffffff0000000000000000000000000000000000000000000000001690556004018290555183917ff1ca1e9147be737b04a2b018a79405f687a97de8dd8a2559bbe62357343af41491a25b505080806115f09061573f565b91505061133c565b505050565b67ffffffffffffffff82166000908152606d6020526040902054829073ffffffffffffffffffffffffffffffffffffffff1680611666576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff8216146116cd576040517fd8a3fb5200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82166004820152602401610706565b607354640100000000900460ff1615611712576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61171a61404e565b67ffffffffffffffff84166000908152606d602052604090206002015460641415611771576040517f05a48e0f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff83166000908152606c6020908152604080832067ffffffffffffffff808916855292529091205416156117b85761082f565b73ffffffffffffffffffffffffffffffffffffffff83166000818152606c6020908152604080832067ffffffffffffffff891680855290835281842080547fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000166001908117909155606d84528285206002018054918201815585529383902090930180547fffffffffffffffffffffffff000000000000000000000000000000000000000016851790555192835290917f43dc749a04ac8fb825cbd514f7c0e13f13bc6f2ee66043b76629d51776cff8e09101610826565b60015473ffffffffffffffffffffffffffffffffffffffff1633146118e1576040517f0f22ca5f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000805433620100008181027fffffffffffffffffffff0000000000000000000000000000000000000000ffff8416178455600180547fffffffffffffffffffffffff000000000000000000000000000000000000000016905560405173ffffffffffffffffffffffffffffffffffffffff919093041692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b607354640100000000900460ff16156119c8576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6119d061404e565b606b60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16634b4fa0c16040518163ffffffff1660e01b815260040160206040518083038186803b158015611a3857600080fd5b505afa158015611a4c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611a709190614ebd565b8015611b1a5750606b546040517ffa00763a00000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff9091169063fa00763a9060240160206040518083038186803b158015611ae057600080fd5b505afa158015611af4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611b189190614ebd565b155b15611b51576040517f0809490800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff81166000908152606d602052604090205473ffffffffffffffffffffffffffffffffffffffff16611bb7576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff81166000908152606d602052604090206001015473ffffffffffffffffffffffffffffffffffffffff163314611c595767ffffffffffffffff81166000908152606d6020526040908190206001015490517fd084e97500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602401610706565b67ffffffffffffffff81166000818152606d60209081526040918290208054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560019093018054909316909255835173ffffffffffffffffffffffffffffffffffffffff909116808252928101919091529092917f6f1dc65165ffffedfd8e507b4a0f1fcfdada045ed11f6c26ba27cedfe87802f0910160405180910390a25050565b611d0c613bf0565b6110c9614343565b67ffffffffffffffff82166000908152606d6020526040902054829073ffffffffffffffffffffffffffffffffffffffff1680611d7d576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff821614611de4576040517fd8a3fb5200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82166004820152602401610706565b607354640100000000900460ff1615611e29576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611e3161404e565b73ffffffffffffffffffffffffffffffffffffffff83166000908152606c6020908152604080832067ffffffffffffffff808916855292529091205416611ecc576040517ff0019fe600000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8516600482015273ffffffffffffffffffffffffffffffffffffffff84166024820152604401610706565b67ffffffffffffffff84166000908152606d6020908152604080832060020180548251818502810185019093528083529192909190830182828015611f4757602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311611f1c575b50505050509050600060018251611f5e91906156fb565b905060005b82518110156120fd578573ffffffffffffffffffffffffffffffffffffffff16838281518110611f9557611f9561582d565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1614156120eb576000838381518110611fcd57611fcd61582d565b6020026020010151905080606d60008a67ffffffffffffffff1667ffffffffffffffff16815260200190815260200160002060020183815481106120135761201361582d565b600091825260208083209190910180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff949094169390931790925567ffffffffffffffff8a168152606d9091526040902060020180548061208d5761208d6157fe565b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff0000000000000000000000000000000000000000169055019055506120fd565b806120f58161573f565b915050611f63565b5073ffffffffffffffffffffffffffffffffffffffff85166000818152606c6020908152604080832067ffffffffffffffff8b168085529083529281902080547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001690555192835290917f182bff9831466789164ca77075fffd84916d35a8180ba73c27e45634549b445b91015b60405180910390a2505050505050565b6000806121a661439e565b9050600081136121e5576040517f43d4cf6600000000000000000000000000000000000000000000000000000000815260048101829052602401610706565b60745460755460009163ffffffff808a1692612202929116615614565b61220c9190615614565b90506000828261222489670de0b6b3a76400006156be565b61222e91906156be565b612238919061567f565b905060006122576bffffffffffffffffffffffff808816908916615614565b905061226f816b033b2e3c9fd0803ce80000006156fb565b8211156122a8576040517fe80fa38100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6122b28183615614565b9998505050505050505050565b607354600090640100000000900460ff1615612307576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61230f61404e565b606b60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16634b4fa0c16040518163ffffffff1660e01b815260040160206040518083038186803b15801561237757600080fd5b505afa15801561238b573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906123af9190614ebd565b80156124595750606b546040517ffa00763a00000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff9091169063fa00763a9060240160206040518083038186803b15801561241f57600080fd5b505afa158015612433573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906124579190614ebd565b155b15612490576040517f0809490800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b606f805467ffffffffffffffff169060006124aa83615778565b82546101009290920a67ffffffffffffffff818102199093169183160217909155606f541690506000806040519080825280602002602001820160405280156124fd578160200160208202803683370190505b506040805180820182526000808252602080830182815267ffffffffffffffff8816808452606e83528584209451855492516bffffffffffffffffffffffff9081166c01000000000000000000000000027fffffffffffffffff000000000000000000000000000000000000000000000000909416911617919091179093558351606081018552338152808201838152818601878152948452606d8352949092208251815473ffffffffffffffffffffffffffffffffffffffff9182167fffffffffffffffffffffffff000000000000000000000000000000000000000091821617835595516001830180549190921696169590951790945591518051949550909361260f9260028501920190614ad8565b505060405133815267ffffffffffffffff841691507f464722b4166576d3dcbba877b999bc35cf911f4eaf434b7eba68fa113951d0bf9060200160405180910390a250905090565b67ffffffffffffffff81166000908152606d6020526040812054819060609073ffffffffffffffffffffffffffffffffffffffff166126c2576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff84166000908152606e6020908152604080832054606d8352928190208054600290910180548351818602810186019094528084526bffffffffffffffffffffffff9095169473ffffffffffffffffffffffffffffffffffffffff90921693909291839183018282801561277457602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311612749575b505050505090509250925092509193909250565b607354640100000000900460ff16156127cd576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6127d561404e565b60695473ffffffffffffffffffffffffffffffffffffffff163314612826576040517f44b0e3c300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60208114612860576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600061286e828401846151d9565b67ffffffffffffffff81166000908152606d602052604090205490915073ffffffffffffffffffffffffffffffffffffffff166128d7576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff81166000908152606e6020526040812080546bffffffffffffffffffffffff169186919061290e8385615658565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555084606f60088282829054906101000a90046bffffffffffffffffffffffff166129659190615658565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055508167ffffffffffffffff167fd39ec07f4e209f627a4c427971473820dc129761ba28de8906bd56f57101d4f88287846129cc9190615614565b6040805192835260208301919091520161218b565b60006129eb6140bb565b607354640100000000900460ff1615612a30576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b612a3861404e565b6000606d81612a4a60208601866151d9565b67ffffffffffffffff16815260208101919091526040016000205473ffffffffffffffffffffffffffffffffffffffff161415612ab3576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000606c81612ac86040860160208701614d8c565b73ffffffffffffffffffffffffffffffffffffffff168152602080820192909252604001600090812091612afe908601866151d9565b67ffffffffffffffff908116825260208201929092526040016000205416905080612b9a57612b3060208401846151d9565b612b406040850160208601614d8c565b6040517ff0019fe600000000000000000000000000000000000000000000000000000000815267ffffffffffffffff909216600483015273ffffffffffffffffffffffffffffffffffffffff166024820152604401610706565b60735463ffffffff16612bb36060850160408601615107565b63ffffffff161115612c1457612bcf6060840160408501615107565b6073546040517ff5d7e01e00000000000000000000000000000000000000000000000000000000815263ffffffff928316600482015291166024820152604401610706565b6040517ff1e14a21000000000000000000000000000000000000000000000000000000008152600090339063f1e14a2190612c57908990899089906004016153f2565b60206040518083038186803b158015612c6f57600080fd5b505afa158015612c83573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612ca79190615277565b90506000612cbf878761056f368990038901896150d2565b90506000612ce2612cd66060880160408901615107565b8760600135858561219b565b90506000606e81612cf660208a018a6151d9565b67ffffffffffffffff1681526020808201929092526040016000908120546c0100000000000000000000000090046bffffffffffffffffffffffff1691606e9190612d43908b018b6151d9565b67ffffffffffffffff168152602081019190915260400160002054612d7691906bffffffffffffffffffffffff16615712565b9050816bffffffffffffffffffffffff16816bffffffffffffffffffffffff161015612dce576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000612ddb86600161562c565b90506000612e6333612df360408c0160208d01614d8c565b612e0060208d018d6151d9565b856040805173ffffffffffffffffffffffffffffffffffffffff958616602080830191909152949095168582015267ffffffffffffffff928316606086015291166080808501919091528151808503909101815260a09093019052815191012090565b60408051610120810190915290915060009080612e8360208d018d6151d9565b67ffffffffffffffff1681526020018b6020016020810190612ea59190614d8c565b73ffffffffffffffffffffffffffffffffffffffff168152602001612ed060608d0160408e01615107565b63ffffffff90811682526060808e0135602080850191909152336040808601919091526bffffffffffffffffffffffff808e16848701528c81166080808801919091528c821660a0808901919091524260c09889015260008b8152607186528481208a5181548c890151978d0151909a167c0100000000000000000000000000000000000000000000000000000000027bffffffffffffffffffffffffffffffffffffffffffffffffffffffff73ffffffffffffffffffffffffffffffffffffffff98891668010000000000000000027fffffffff00000000000000000000000000000000000000000000000000000000909c1667ffffffffffffffff909316929092179a909a171698909817885595890151600188015590880151908801518216740100000000000000000000000000000000000000000292169190911760028501559385015160038401805460e088015187166c01000000000000000000000000027fffffffffffffffff0000000000000000000000000000000000000000000000009091169290961691909117949094179093556101008401516004909201919091559192508691606e9161308a908e018e6151d9565b67ffffffffffffffff16815260208101919091526040016000208054600c906130d29084906c0100000000000000000000000090046bffffffffffffffffffffffff16615658565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550817f99f7f4e65b4b9fbabd4e357c47ed3099b36e57ecd3a43e84662f34c207d0ebe48260405161313091906154b1565b60405180910390a282606c600061314d60408e0160208f01614d8c565b73ffffffffffffffffffffffffffffffffffffffff168152602080820192909252604001600090812091613183908e018e6151d9565b67ffffffffffffffff9081168252602082019290925260400160002080547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001692909116919091179055509a9950505050505050505050565b67ffffffffffffffff81166000908152606d602052604081205473ffffffffffffffffffffffffffffffffffffffff16613242576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5067ffffffffffffffff166000908152606d602052604090205473ffffffffffffffffffffffffffffffffffffffff1690565b600054610100900460ff16158080156132955750600054600160ff909116105b806132af5750303b1580156132af575060005460ff166001145b61333b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a65640000000000000000000000000000000000006064820152608401610706565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055801561339957600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101001790555b6133a1614491565b6133ac336000614530565b6069805473ffffffffffffffffffffffffffffffffffffffff8087167fffffffffffffffffffffffff000000000000000000000000000000000000000092831617909255606a8054868416908316179055606b805492851692909116919091179055801561082f57600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498906020015b60405180910390a150505050565b67ffffffffffffffff82166000908152606d6020526040902054829073ffffffffffffffffffffffffffffffffffffffff16806134e0576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff821614613547576040517fd8a3fb5200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82166004820152602401610706565b607354640100000000900460ff161561358c576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61359461404e565b61359d846137fb565b156135d4576040517fb42f66e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61082f8484613c47565b6135e6613bf0565b6069546040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015260009173ffffffffffffffffffffffffffffffffffffffff16906370a082319060240160206040518083038186803b15801561365057600080fd5b505afa158015613664573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061368891906150ee565b606f549091506801000000000000000090046bffffffffffffffffffffffff16818111156136ec576040517fa99da3020000000000000000000000000000000000000000000000000000000081526004810182905260248101839052604401610706565b818110156115f857600061370082846156fb565b6069546040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff87811660048301526024820184905292935091169063a9059cbb90604401602060405180830381600087803b15801561377657600080fd5b505af115801561378a573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906137ae9190614ebd565b506040805173ffffffffffffffffffffffffffffffffffffffff86168152602081018390527f59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b4366009101613469565b67ffffffffffffffff81166000908152606d602090815260408083206002018054825181850281018501909352808352849383018282801561387357602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311613848575b505050505090506000613884610ef7565b905060005b8251811015613a2f5760005b8251811015613a1c5760006139cc8483815181106138b5576138b561582d565b60200260200101518685815181106138cf576138cf61582d565b602002602001015189606c60008a89815181106138ee576138ee61582d565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008c67ffffffffffffffff1667ffffffffffffffff16815260200190815260200160002060009054906101000a900467ffffffffffffffff166040805173ffffffffffffffffffffffffffffffffffffffff958616602080830191909152949095168582015267ffffffffffffffff928316606086015291166080808501919091528151808503909101815260a09093019052815191012090565b60008181526071602052604090206002015490915073ffffffffffffffffffffffffffffffffffffffff1615613a09575060019695505050505050565b5080613a148161573f565b915050613895565b5080613a278161573f565b915050613889565b506000949350505050565b613a42614670565b613a78576040517fad77f06100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80613aaf576040517f75158c3b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b606854811015613b0f57613afc60688281548110613ad257613ad261582d565b60009182526020909120015460669073ffffffffffffffffffffffffffffffffffffffff16614680565b5080613b078161573f565b915050613ab2565b5060005b81811015613b6057613b4d838383818110613b3057613b3061582d565b9050602002016020810190613b459190614d8c565b6066906146a9565b5080613b588161573f565b915050613b13565b50613b6d60688383614b5e565b507ff263cfb3e4298332e776194610cf9fdc09ccb3ada8b9aa39764d882e11fbf0a0828233604051613ba19392919061532e565b60405180910390a15050565b613bb5613bf0565b613bbe816146cb565b50565b6000613bce606683614798565b92915050565b73ffffffffffffffffffffffffffffffffffffffff163b151590565b60005462010000900473ffffffffffffffffffffffffffffffffffffffff1633146110c9576040517f2b5c74de00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b607354640100000000900460ff1615613c8c576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff82166000908152606d602090815260408083208151606081018352815473ffffffffffffffffffffffffffffffffffffffff908116825260018301541681850152600282018054845181870281018701865281815292959394860193830182828015613d3757602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311613d0c575b5050509190925250505067ffffffffffffffff84166000908152606e60205260408120549192506bffffffffffffffffffffffff909116905b826040015151811015613e1657606c600084604001518381518110613d9757613d9761582d565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040908101600090812067ffffffffffffffff89168252909252902080547fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000016905580613e0e8161573f565b915050613d70565b5067ffffffffffffffff84166000908152606d6020526040812080547fffffffffffffffffffffffff00000000000000000000000000000000000000009081168255600182018054909116905590613e716002830182614bd6565b505067ffffffffffffffff84166000908152606e6020526040902080547fffffffffffffffff000000000000000000000000000000000000000000000000169055606f8054829190600890613ee19084906801000000000000000090046bffffffffffffffffffffffff16615712565b82546101009290920a6bffffffffffffffffffffffff8181021990931691831602179091556069546040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff878116600483015292851660248201529116915063a9059cbb90604401602060405180830381600087803b158015613f7b57600080fd5b505af1158015613f8f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613fb39190614ebd565b613fe9576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805173ffffffffffffffffffffffffffffffffffffffff851681526bffffffffffffffffffffffff8316602082015267ffffffffffffffff8616917fe8ed5b475a5b5987aa9165e8731bb78043f39eee32ec5a1169a89e27fcd498159101610826565b60345460ff16156110c9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f5061757361626c653a20706175736564000000000000000000000000000000006044820152606401610706565b6140c433613bc1565b6110c9576040517f0809490800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005a61138881101561410c57600080fd5b61138881039050846040820482031161412457600080fd5b50823b61413057600080fd5b60008083516020850160008789f1949350505050565b604080516060810182526000808252602082018190529181018290529061416b61439e565b9050600081136141aa576040517f43d4cf6600000000000000000000000000000000000000000000000000000000815260048101829052602401610706565b6000815a8b6141b98c89615614565b6141c39190615614565b6141cd91906156fb565b6141df86670de0b6b3a76400006156be565b6141e991906156be565b6141f3919061567f565b905060006142126bffffffffffffffffffffffff808916908b16615614565b905061422a816b033b2e3c9fd0803ce80000006156fb565b821115614263576040517fe80fa38100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600061427260ff8a168b615693565b90508260006142896142848584615614565b6147c7565b604080516060810182526bffffffffffffffffffffffff958616815293851660208501529316928201929092529c9b505050505050505050505050565b6142ce614869565b603480547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa335b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390a1565b61434b61404e565b603480547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a2586143193390565b607354606a54604080517ffeaf968c000000000000000000000000000000000000000000000000000000008152905160009365010000000000900463ffffffff1692831515928592839273ffffffffffffffffffffffffffffffffffffffff169163feaf968c9160048083019260a0929190829003018186803b15801561442457600080fd5b505afa158015614438573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061445c9190615227565b5093505092505082801561447e575061447581426156fb565b8463ffffffff16105b156144895760725491505b509392505050565b600054610100900460ff16614528576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201527f6e697469616c697a696e670000000000000000000000000000000000000000006064820152608401610706565b6110c96148d5565b600054610100900460ff166145c7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201527f6e697469616c697a696e670000000000000000000000000000000000000000006064820152608401610706565b73ffffffffffffffffffffffffffffffffffffffff8216614614576040517fb6b515fa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000805473ffffffffffffffffffffffffffffffffffffffff80851662010000027fffffffffffffffffffff0000000000000000000000000000000000000000ffff9092169190911790915581161561063657610636816146cb565b600061467a613bf0565b50600190565b60006146a28373ffffffffffffffffffffffffffffffffffffffff8416614996565b9392505050565b60006146a28373ffffffffffffffffffffffffffffffffffffffff8416614a89565b73ffffffffffffffffffffffffffffffffffffffff811633141561471b576040517f282010c300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8381169182179092556000805460405192936201000090910416917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b73ffffffffffffffffffffffffffffffffffffffff8116600090815260018301602052604081205415156146a2565b60006bffffffffffffffffffffffff821115614865576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f53616665436173743a2076616c756520646f65736e27742066697420696e203960448201527f36206269747300000000000000000000000000000000000000000000000000006064820152608401610706565b5090565b60345460ff166110c9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f5061757361626c653a206e6f74207061757365640000000000000000000000006044820152606401610706565b600054610100900460ff1661496c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201527f6e697469616c697a696e670000000000000000000000000000000000000000006064820152608401610706565b603480547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055565b60008181526001830160205260408120548015614a7f5760006149ba6001836156fb565b85549091506000906149ce906001906156fb565b9050818114614a335760008660000182815481106149ee576149ee61582d565b9060005260206000200154905080876000018481548110614a1157614a1161582d565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080614a4457614a446157fe565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050613bce565b6000915050613bce565b6000818152600183016020526040812054614ad057508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155613bce565b506000613bce565b828054828255906000526020600020908101928215614b52579160200282015b82811115614b5257825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190614af8565b50614865929150614bf0565b828054828255906000526020600020908101928215614b52579160200282015b82811115614b525781547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff843516178255602090920191600190910190614b7e565b5080546000825590600052602060002090810190613bbe91905b5b808211156148655760008155600101614bf1565b803573ffffffffffffffffffffffffffffffffffffffff81168114614c2957600080fd5b919050565b60008083601f840112614c4057600080fd5b50813567ffffffffffffffff811115614c5857600080fd5b6020830191508360208260051b8501011115614c7357600080fd5b9250929050565b60008083601f840112614c8c57600080fd5b50813567ffffffffffffffff811115614ca457600080fd5b602083019150836020828501011115614c7357600080fd5b600060808284031215614cce57600080fd5b6040516080810181811067ffffffffffffffff82111715614cf157614cf161585c565b604052905080614d0083614d49565b8152614d0e60208401614c05565b6020820152614d1f60408401614d35565b6040820152606083013560608201525092915050565b803563ffffffff81168114614c2957600080fd5b803567ffffffffffffffff81168114614c2957600080fd5b803560ff81168114614c2957600080fd5b805169ffffffffffffffffffff81168114614c2957600080fd5b600060208284031215614d9e57600080fd5b6146a282614c05565b600080600060608486031215614dbc57600080fd5b614dc584614c05565b9250614dd360208501614c05565b9150614de160408501614c05565b90509250925092565b60008060008060608587031215614e0057600080fd5b614e0985614c05565b935060208501359250604085013567ffffffffffffffff811115614e2c57600080fd5b614e3887828801614c7a565b95989497509550505050565b60008060408385031215614e5757600080fd5b614e6083614c05565b91506020830135614e708161588b565b809150509250929050565b60008060208385031215614e8e57600080fd5b823567ffffffffffffffff811115614ea557600080fd5b614eb185828601614c2e565b90969095509350505050565b600060208284031215614ecf57600080fd5b815180151581146146a257600080fd5b6000806000806000806000806000806104c08b8d031215614eff57600080fd5b8a35995060208b013567ffffffffffffffff80821115614f1e57600080fd5b614f2a8e838f01614c7a565b909b50995060408d0135915080821115614f4357600080fd5b614f4f8e838f01614c7a565b9099509750879150614f6360608e01614c05565b96508d609f8e0112614f7457600080fd5b60405191506103e082018281108282111715614f9257614f9261585c565b604052508060808d016104608e018f811115614fad57600080fd5b60005b601f811015614fd757614fc283614c05565b84526020938401939290920191600101614fb0565b50839750614fe481614d61565b9650505050506104808b013591506104a08b013590509295989b9194979a5092959850565b600080600083850360a081121561501f57600080fd5b843567ffffffffffffffff81111561503657600080fd5b61504287828801614c7a565b90955093505060807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08201121561507857600080fd5b506020840190509250925092565b600080600060a0848603121561509b57600080fd5b833567ffffffffffffffff8111156150b257600080fd5b6150be86828701614c7a565b9094509250614de190508560208601614cbc565b6000608082840312156150e457600080fd5b6146a28383614cbc565b60006020828403121561510057600080fd5b5051919050565b60006020828403121561511957600080fd5b6146a282614d35565b6000806000806080858703121561513857600080fd5b61514185614d35565b93506020850135925060408501356151588161588b565b915060608501356151688161588b565b939692955090935050565b60008060008060008060c0878903121561518c57600080fd5b61519587614d35565b95506151a360208801614d35565b945060408701359350606087013592506151bf60808801614d35565b91506151cd60a08801614d35565b90509295509295509295565b6000602082840312156151eb57600080fd5b6146a282614d49565b6000806040838503121561520757600080fd5b61521083614d49565b915061521e60208401614c05565b90509250929050565b600080600080600060a0868803121561523f57600080fd5b61524886614d72565b945060208601519350604086015192506060860151915061526b60808701614d72565b90509295509295909350565b60006020828403121561528957600080fd5b81516146a28161588b565b600081518084526020808501945080840160005b838110156152da57815173ffffffffffffffffffffffffffffffffffffffff16875295820195908201906001016152a8565b509495945050505050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b6040808252810183905260008460608301825b8681101561537c5773ffffffffffffffffffffffffffffffffffffffff61536784614c05565b16825260209283019290910190600101615341565b50809250505073ffffffffffffffffffffffffffffffffffffffff83166020830152949350505050565b6020815260006146a26020830184615294565b8581526060602082015260006153d36060830186886152e5565b82810360408401526153e68185876152e5565b98975050505050505050565b60a08152600061540660a0830185876152e5565b905067ffffffffffffffff61541a84614d49565b16602083015273ffffffffffffffffffffffffffffffffffffffff61544160208501614c05565b16604083015263ffffffff61545860408501614d35565b16606083015260608301356080830152949350505050565b60208101600383106154ab577f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b91905290565b60006101208201905067ffffffffffffffff835116825273ffffffffffffffffffffffffffffffffffffffff602084015116602083015260408301516154ff604084018263ffffffff169052565b50606083015160608301526080830151615531608084018273ffffffffffffffffffffffffffffffffffffffff169052565b5060a083015161555160a08401826bffffffffffffffffffffffff169052565b5060c083015161557160c08401826bffffffffffffffffffffffff169052565b5060e083015161559160e08401826bffffffffffffffffffffffff169052565b5061010092830151919092015290565b63ffffffff831681526040602082015260006155c06040830184615294565b949350505050565b6bffffffffffffffffffffffff8416815273ffffffffffffffffffffffffffffffffffffffff8316602082015260606040820152600061560b6060830184615294565b95945050505050565b60008219821115615627576156276157a0565b500190565b600067ffffffffffffffff80831681851680830382111561564f5761564f6157a0565b01949350505050565b60006bffffffffffffffffffffffff80831681851680830382111561564f5761564f6157a0565b60008261568e5761568e6157cf565b500490565b60006bffffffffffffffffffffffff808416806156b2576156b26157cf565b92169190910492915050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156156f6576156f66157a0565b500290565b60008282101561570d5761570d6157a0565b500390565b60006bffffffffffffffffffffffff83811690831681811015615737576157376157a0565b039392505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415615771576157716157a0565b5060010190565b600067ffffffffffffffff80831681811415615796576157966157a0565b6001019392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6bffffffffffffffffffffffff81168114613bbe57600080fdfea164736f6c6343000806000a496e697469616c697a61626c653a20636f6e7472616374206973206e6f742069",
}

var FunctionsRegistryABI = FunctionsRegistryMetaData.ABI

var FunctionsRegistryBin = FunctionsRegistryMetaData.Bin

func DeployFunctionsRegistry(auth *bind.TransactOpts, backend bind.ContractBackend, link common.Address, linkEthFeed common.Address, oracle common.Address) (common.Address, *types.Transaction, *FunctionsRegistry, error) {
	parsed, err := FunctionsRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(FunctionsRegistryBin), backend, link, linkEthFeed, oracle)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &FunctionsRegistry{FunctionsRegistryCaller: FunctionsRegistryCaller{contract: contract}, FunctionsRegistryTransactor: FunctionsRegistryTransactor{contract: contract}, FunctionsRegistryFilterer: FunctionsRegistryFilterer{contract: contract}}, nil
}

type FunctionsRegistry struct {
	address common.Address
	abi     abi.ABI
	FunctionsRegistryCaller
	FunctionsRegistryTransactor
	FunctionsRegistryFilterer
}

type FunctionsRegistryCaller struct {
	contract *bind.BoundContract
}

type FunctionsRegistryTransactor struct {
	contract *bind.BoundContract
}

type FunctionsRegistryFilterer struct {
	contract *bind.BoundContract
}

type FunctionsRegistrySession struct {
	Contract     *FunctionsRegistry
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type FunctionsRegistryCallerSession struct {
	Contract *FunctionsRegistryCaller
	CallOpts bind.CallOpts
}

type FunctionsRegistryTransactorSession struct {
	Contract     *FunctionsRegistryTransactor
	TransactOpts bind.TransactOpts
}

type FunctionsRegistryRaw struct {
	Contract *FunctionsRegistry
}

type FunctionsRegistryCallerRaw struct {
	Contract *FunctionsRegistryCaller
}

type FunctionsRegistryTransactorRaw struct {
	Contract *FunctionsRegistryTransactor
}

func NewFunctionsRegistry(address common.Address, backend bind.ContractBackend) (*FunctionsRegistry, error) {
	abi, err := abi.JSON(strings.NewReader(FunctionsRegistryABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindFunctionsRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FunctionsRegistry{address: address, abi: abi, FunctionsRegistryCaller: FunctionsRegistryCaller{contract: contract}, FunctionsRegistryTransactor: FunctionsRegistryTransactor{contract: contract}, FunctionsRegistryFilterer: FunctionsRegistryFilterer{contract: contract}}, nil
}

func NewFunctionsRegistryCaller(address common.Address, caller bind.ContractCaller) (*FunctionsRegistryCaller, error) {
	contract, err := bindFunctionsRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FunctionsRegistryCaller{contract: contract}, nil
}

func NewFunctionsRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*FunctionsRegistryTransactor, error) {
	contract, err := bindFunctionsRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FunctionsRegistryTransactor{contract: contract}, nil
}

func NewFunctionsRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*FunctionsRegistryFilterer, error) {
	contract, err := bindFunctionsRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FunctionsRegistryFilterer{contract: contract}, nil
}

func bindFunctionsRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FunctionsRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_FunctionsRegistry *FunctionsRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FunctionsRegistry.Contract.FunctionsRegistryCaller.contract.Call(opts, result, method, params...)
}

func (_FunctionsRegistry *FunctionsRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.FunctionsRegistryTransactor.contract.Transfer(opts)
}

func (_FunctionsRegistry *FunctionsRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.FunctionsRegistryTransactor.contract.Transact(opts, method, params...)
}

func (_FunctionsRegistry *FunctionsRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FunctionsRegistry.Contract.contract.Call(opts, result, method, params...)
}

func (_FunctionsRegistry *FunctionsRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.contract.Transfer(opts)
}

func (_FunctionsRegistry *FunctionsRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.contract.Transact(opts, method, params...)
}

func (_FunctionsRegistry *FunctionsRegistryCaller) MAXCONSUMERS(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _FunctionsRegistry.contract.Call(opts, &out, "MAX_CONSUMERS")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_FunctionsRegistry *FunctionsRegistrySession) MAXCONSUMERS() (uint16, error) {
	return _FunctionsRegistry.Contract.MAXCONSUMERS(&_FunctionsRegistry.CallOpts)
}

func (_FunctionsRegistry *FunctionsRegistryCallerSession) MAXCONSUMERS() (uint16, error) {
	return _FunctionsRegistry.Contract.MAXCONSUMERS(&_FunctionsRegistry.CallOpts)
}

func (_FunctionsRegistry *FunctionsRegistryCaller) EstimateCost(opts *bind.CallOpts, gasLimit uint32, gasPrice *big.Int, donFee *big.Int, registryFee *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _FunctionsRegistry.contract.Call(opts, &out, "estimateCost", gasLimit, gasPrice, donFee, registryFee)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_FunctionsRegistry *FunctionsRegistrySession) EstimateCost(gasLimit uint32, gasPrice *big.Int, donFee *big.Int, registryFee *big.Int) (*big.Int, error) {
	return _FunctionsRegistry.Contract.EstimateCost(&_FunctionsRegistry.CallOpts, gasLimit, gasPrice, donFee, registryFee)
}

func (_FunctionsRegistry *FunctionsRegistryCallerSession) EstimateCost(gasLimit uint32, gasPrice *big.Int, donFee *big.Int, registryFee *big.Int) (*big.Int, error) {
	return _FunctionsRegistry.Contract.EstimateCost(&_FunctionsRegistry.CallOpts, gasLimit, gasPrice, donFee, registryFee)
}

func (_FunctionsRegistry *FunctionsRegistryCaller) GetAuthorizedSenders(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _FunctionsRegistry.contract.Call(opts, &out, "getAuthorizedSenders")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_FunctionsRegistry *FunctionsRegistrySession) GetAuthorizedSenders() ([]common.Address, error) {
	return _FunctionsRegistry.Contract.GetAuthorizedSenders(&_FunctionsRegistry.CallOpts)
}

func (_FunctionsRegistry *FunctionsRegistryCallerSession) GetAuthorizedSenders() ([]common.Address, error) {
	return _FunctionsRegistry.Contract.GetAuthorizedSenders(&_FunctionsRegistry.CallOpts)
}

func (_FunctionsRegistry *FunctionsRegistryCaller) GetConfig(opts *bind.CallOpts) (GetConfig,

	error) {
	var out []interface{}
	err := _FunctionsRegistry.contract.Call(opts, &out, "getConfig")

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

func (_FunctionsRegistry *FunctionsRegistrySession) GetConfig() (GetConfig,

	error) {
	return _FunctionsRegistry.Contract.GetConfig(&_FunctionsRegistry.CallOpts)
}

func (_FunctionsRegistry *FunctionsRegistryCallerSession) GetConfig() (GetConfig,

	error) {
	return _FunctionsRegistry.Contract.GetConfig(&_FunctionsRegistry.CallOpts)
}

func (_FunctionsRegistry *FunctionsRegistryCaller) GetCurrentsubscriptionId(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _FunctionsRegistry.contract.Call(opts, &out, "getCurrentsubscriptionId")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_FunctionsRegistry *FunctionsRegistrySession) GetCurrentsubscriptionId() (uint64, error) {
	return _FunctionsRegistry.Contract.GetCurrentsubscriptionId(&_FunctionsRegistry.CallOpts)
}

func (_FunctionsRegistry *FunctionsRegistryCallerSession) GetCurrentsubscriptionId() (uint64, error) {
	return _FunctionsRegistry.Contract.GetCurrentsubscriptionId(&_FunctionsRegistry.CallOpts)
}

func (_FunctionsRegistry *FunctionsRegistryCaller) GetRequestConfig(opts *bind.CallOpts) (uint32, []common.Address, error) {
	var out []interface{}
	err := _FunctionsRegistry.contract.Call(opts, &out, "getRequestConfig")

	if err != nil {
		return *new(uint32), *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)
	out1 := *abi.ConvertType(out[1], new([]common.Address)).(*[]common.Address)

	return out0, out1, err

}

func (_FunctionsRegistry *FunctionsRegistrySession) GetRequestConfig() (uint32, []common.Address, error) {
	return _FunctionsRegistry.Contract.GetRequestConfig(&_FunctionsRegistry.CallOpts)
}

func (_FunctionsRegistry *FunctionsRegistryCallerSession) GetRequestConfig() (uint32, []common.Address, error) {
	return _FunctionsRegistry.Contract.GetRequestConfig(&_FunctionsRegistry.CallOpts)
}

func (_FunctionsRegistry *FunctionsRegistryCaller) GetRequiredFee(opts *bind.CallOpts, arg0 []byte, arg1 IFunctionsBillingRegistryRequestBilling) (*big.Int, error) {
	var out []interface{}
	err := _FunctionsRegistry.contract.Call(opts, &out, "getRequiredFee", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_FunctionsRegistry *FunctionsRegistrySession) GetRequiredFee(arg0 []byte, arg1 IFunctionsBillingRegistryRequestBilling) (*big.Int, error) {
	return _FunctionsRegistry.Contract.GetRequiredFee(&_FunctionsRegistry.CallOpts, arg0, arg1)
}

func (_FunctionsRegistry *FunctionsRegistryCallerSession) GetRequiredFee(arg0 []byte, arg1 IFunctionsBillingRegistryRequestBilling) (*big.Int, error) {
	return _FunctionsRegistry.Contract.GetRequiredFee(&_FunctionsRegistry.CallOpts, arg0, arg1)
}

func (_FunctionsRegistry *FunctionsRegistryCaller) GetSubscription(opts *bind.CallOpts, subscriptionId uint64) (GetSubscription,

	error) {
	var out []interface{}
	err := _FunctionsRegistry.contract.Call(opts, &out, "getSubscription", subscriptionId)

	outstruct := new(GetSubscription)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Balance = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Owner = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.Consumers = *abi.ConvertType(out[2], new([]common.Address)).(*[]common.Address)

	return *outstruct, err

}

func (_FunctionsRegistry *FunctionsRegistrySession) GetSubscription(subscriptionId uint64) (GetSubscription,

	error) {
	return _FunctionsRegistry.Contract.GetSubscription(&_FunctionsRegistry.CallOpts, subscriptionId)
}

func (_FunctionsRegistry *FunctionsRegistryCallerSession) GetSubscription(subscriptionId uint64) (GetSubscription,

	error) {
	return _FunctionsRegistry.Contract.GetSubscription(&_FunctionsRegistry.CallOpts, subscriptionId)
}

func (_FunctionsRegistry *FunctionsRegistryCaller) GetSubscriptionOwner(opts *bind.CallOpts, subscriptionId uint64) (common.Address, error) {
	var out []interface{}
	err := _FunctionsRegistry.contract.Call(opts, &out, "getSubscriptionOwner", subscriptionId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_FunctionsRegistry *FunctionsRegistrySession) GetSubscriptionOwner(subscriptionId uint64) (common.Address, error) {
	return _FunctionsRegistry.Contract.GetSubscriptionOwner(&_FunctionsRegistry.CallOpts, subscriptionId)
}

func (_FunctionsRegistry *FunctionsRegistryCallerSession) GetSubscriptionOwner(subscriptionId uint64) (common.Address, error) {
	return _FunctionsRegistry.Contract.GetSubscriptionOwner(&_FunctionsRegistry.CallOpts, subscriptionId)
}

func (_FunctionsRegistry *FunctionsRegistryCaller) GetTotalBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FunctionsRegistry.contract.Call(opts, &out, "getTotalBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_FunctionsRegistry *FunctionsRegistrySession) GetTotalBalance() (*big.Int, error) {
	return _FunctionsRegistry.Contract.GetTotalBalance(&_FunctionsRegistry.CallOpts)
}

func (_FunctionsRegistry *FunctionsRegistryCallerSession) GetTotalBalance() (*big.Int, error) {
	return _FunctionsRegistry.Contract.GetTotalBalance(&_FunctionsRegistry.CallOpts)
}

func (_FunctionsRegistry *FunctionsRegistryCaller) IsAuthorizedSender(opts *bind.CallOpts, sender common.Address) (bool, error) {
	var out []interface{}
	err := _FunctionsRegistry.contract.Call(opts, &out, "isAuthorizedSender", sender)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_FunctionsRegistry *FunctionsRegistrySession) IsAuthorizedSender(sender common.Address) (bool, error) {
	return _FunctionsRegistry.Contract.IsAuthorizedSender(&_FunctionsRegistry.CallOpts, sender)
}

func (_FunctionsRegistry *FunctionsRegistryCallerSession) IsAuthorizedSender(sender common.Address) (bool, error) {
	return _FunctionsRegistry.Contract.IsAuthorizedSender(&_FunctionsRegistry.CallOpts, sender)
}

func (_FunctionsRegistry *FunctionsRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FunctionsRegistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_FunctionsRegistry *FunctionsRegistrySession) Owner() (common.Address, error) {
	return _FunctionsRegistry.Contract.Owner(&_FunctionsRegistry.CallOpts)
}

func (_FunctionsRegistry *FunctionsRegistryCallerSession) Owner() (common.Address, error) {
	return _FunctionsRegistry.Contract.Owner(&_FunctionsRegistry.CallOpts)
}

func (_FunctionsRegistry *FunctionsRegistryCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _FunctionsRegistry.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_FunctionsRegistry *FunctionsRegistrySession) Paused() (bool, error) {
	return _FunctionsRegistry.Contract.Paused(&_FunctionsRegistry.CallOpts)
}

func (_FunctionsRegistry *FunctionsRegistryCallerSession) Paused() (bool, error) {
	return _FunctionsRegistry.Contract.Paused(&_FunctionsRegistry.CallOpts)
}

func (_FunctionsRegistry *FunctionsRegistryCaller) PendingRequestExists(opts *bind.CallOpts, subscriptionId uint64) (bool, error) {
	var out []interface{}
	err := _FunctionsRegistry.contract.Call(opts, &out, "pendingRequestExists", subscriptionId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_FunctionsRegistry *FunctionsRegistrySession) PendingRequestExists(subscriptionId uint64) (bool, error) {
	return _FunctionsRegistry.Contract.PendingRequestExists(&_FunctionsRegistry.CallOpts, subscriptionId)
}

func (_FunctionsRegistry *FunctionsRegistryCallerSession) PendingRequestExists(subscriptionId uint64) (bool, error) {
	return _FunctionsRegistry.Contract.PendingRequestExists(&_FunctionsRegistry.CallOpts, subscriptionId)
}

func (_FunctionsRegistry *FunctionsRegistryTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsRegistry.contract.Transact(opts, "acceptOwnership")
}

func (_FunctionsRegistry *FunctionsRegistrySession) AcceptOwnership() (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.AcceptOwnership(&_FunctionsRegistry.TransactOpts)
}

func (_FunctionsRegistry *FunctionsRegistryTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.AcceptOwnership(&_FunctionsRegistry.TransactOpts)
}

func (_FunctionsRegistry *FunctionsRegistryTransactor) AcceptSubscriptionOwnerTransfer(opts *bind.TransactOpts, subscriptionId uint64) (*types.Transaction, error) {
	return _FunctionsRegistry.contract.Transact(opts, "acceptSubscriptionOwnerTransfer", subscriptionId)
}

func (_FunctionsRegistry *FunctionsRegistrySession) AcceptSubscriptionOwnerTransfer(subscriptionId uint64) (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.AcceptSubscriptionOwnerTransfer(&_FunctionsRegistry.TransactOpts, subscriptionId)
}

func (_FunctionsRegistry *FunctionsRegistryTransactorSession) AcceptSubscriptionOwnerTransfer(subscriptionId uint64) (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.AcceptSubscriptionOwnerTransfer(&_FunctionsRegistry.TransactOpts, subscriptionId)
}

func (_FunctionsRegistry *FunctionsRegistryTransactor) AddConsumer(opts *bind.TransactOpts, subscriptionId uint64, consumer common.Address) (*types.Transaction, error) {
	return _FunctionsRegistry.contract.Transact(opts, "addConsumer", subscriptionId, consumer)
}

func (_FunctionsRegistry *FunctionsRegistrySession) AddConsumer(subscriptionId uint64, consumer common.Address) (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.AddConsumer(&_FunctionsRegistry.TransactOpts, subscriptionId, consumer)
}

func (_FunctionsRegistry *FunctionsRegistryTransactorSession) AddConsumer(subscriptionId uint64, consumer common.Address) (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.AddConsumer(&_FunctionsRegistry.TransactOpts, subscriptionId, consumer)
}

func (_FunctionsRegistry *FunctionsRegistryTransactor) CancelSubscription(opts *bind.TransactOpts, subscriptionId uint64, to common.Address) (*types.Transaction, error) {
	return _FunctionsRegistry.contract.Transact(opts, "cancelSubscription", subscriptionId, to)
}

func (_FunctionsRegistry *FunctionsRegistrySession) CancelSubscription(subscriptionId uint64, to common.Address) (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.CancelSubscription(&_FunctionsRegistry.TransactOpts, subscriptionId, to)
}

func (_FunctionsRegistry *FunctionsRegistryTransactorSession) CancelSubscription(subscriptionId uint64, to common.Address) (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.CancelSubscription(&_FunctionsRegistry.TransactOpts, subscriptionId, to)
}

func (_FunctionsRegistry *FunctionsRegistryTransactor) CreateSubscription(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsRegistry.contract.Transact(opts, "createSubscription")
}

func (_FunctionsRegistry *FunctionsRegistrySession) CreateSubscription() (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.CreateSubscription(&_FunctionsRegistry.TransactOpts)
}

func (_FunctionsRegistry *FunctionsRegistryTransactorSession) CreateSubscription() (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.CreateSubscription(&_FunctionsRegistry.TransactOpts)
}

func (_FunctionsRegistry *FunctionsRegistryTransactor) FulfillAndBill(opts *bind.TransactOpts, requestId [32]byte, response []byte, err []byte, transmitter common.Address, signers [31]common.Address, signerCount uint8, reportValidationGas *big.Int, initialGas *big.Int) (*types.Transaction, error) {
	return _FunctionsRegistry.contract.Transact(opts, "fulfillAndBill", requestId, response, err, transmitter, signers, signerCount, reportValidationGas, initialGas)
}

func (_FunctionsRegistry *FunctionsRegistrySession) FulfillAndBill(requestId [32]byte, response []byte, err []byte, transmitter common.Address, signers [31]common.Address, signerCount uint8, reportValidationGas *big.Int, initialGas *big.Int) (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.FulfillAndBill(&_FunctionsRegistry.TransactOpts, requestId, response, err, transmitter, signers, signerCount, reportValidationGas, initialGas)
}

func (_FunctionsRegistry *FunctionsRegistryTransactorSession) FulfillAndBill(requestId [32]byte, response []byte, err []byte, transmitter common.Address, signers [31]common.Address, signerCount uint8, reportValidationGas *big.Int, initialGas *big.Int) (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.FulfillAndBill(&_FunctionsRegistry.TransactOpts, requestId, response, err, transmitter, signers, signerCount, reportValidationGas, initialGas)
}

func (_FunctionsRegistry *FunctionsRegistryTransactor) Initialize(opts *bind.TransactOpts, link common.Address, linkEthFeed common.Address, oracle common.Address) (*types.Transaction, error) {
	return _FunctionsRegistry.contract.Transact(opts, "initialize", link, linkEthFeed, oracle)
}

func (_FunctionsRegistry *FunctionsRegistrySession) Initialize(link common.Address, linkEthFeed common.Address, oracle common.Address) (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.Initialize(&_FunctionsRegistry.TransactOpts, link, linkEthFeed, oracle)
}

func (_FunctionsRegistry *FunctionsRegistryTransactorSession) Initialize(link common.Address, linkEthFeed common.Address, oracle common.Address) (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.Initialize(&_FunctionsRegistry.TransactOpts, link, linkEthFeed, oracle)
}

func (_FunctionsRegistry *FunctionsRegistryTransactor) OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _FunctionsRegistry.contract.Transact(opts, "onTokenTransfer", arg0, amount, data)
}

func (_FunctionsRegistry *FunctionsRegistrySession) OnTokenTransfer(arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.OnTokenTransfer(&_FunctionsRegistry.TransactOpts, arg0, amount, data)
}

func (_FunctionsRegistry *FunctionsRegistryTransactorSession) OnTokenTransfer(arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.OnTokenTransfer(&_FunctionsRegistry.TransactOpts, arg0, amount, data)
}

func (_FunctionsRegistry *FunctionsRegistryTransactor) OracleWithdraw(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FunctionsRegistry.contract.Transact(opts, "oracleWithdraw", recipient, amount)
}

func (_FunctionsRegistry *FunctionsRegistrySession) OracleWithdraw(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.OracleWithdraw(&_FunctionsRegistry.TransactOpts, recipient, amount)
}

func (_FunctionsRegistry *FunctionsRegistryTransactorSession) OracleWithdraw(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.OracleWithdraw(&_FunctionsRegistry.TransactOpts, recipient, amount)
}

func (_FunctionsRegistry *FunctionsRegistryTransactor) OwnerCancelSubscription(opts *bind.TransactOpts, subscriptionId uint64) (*types.Transaction, error) {
	return _FunctionsRegistry.contract.Transact(opts, "ownerCancelSubscription", subscriptionId)
}

func (_FunctionsRegistry *FunctionsRegistrySession) OwnerCancelSubscription(subscriptionId uint64) (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.OwnerCancelSubscription(&_FunctionsRegistry.TransactOpts, subscriptionId)
}

func (_FunctionsRegistry *FunctionsRegistryTransactorSession) OwnerCancelSubscription(subscriptionId uint64) (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.OwnerCancelSubscription(&_FunctionsRegistry.TransactOpts, subscriptionId)
}

func (_FunctionsRegistry *FunctionsRegistryTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsRegistry.contract.Transact(opts, "pause")
}

func (_FunctionsRegistry *FunctionsRegistrySession) Pause() (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.Pause(&_FunctionsRegistry.TransactOpts)
}

func (_FunctionsRegistry *FunctionsRegistryTransactorSession) Pause() (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.Pause(&_FunctionsRegistry.TransactOpts)
}

func (_FunctionsRegistry *FunctionsRegistryTransactor) RecoverFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _FunctionsRegistry.contract.Transact(opts, "recoverFunds", to)
}

func (_FunctionsRegistry *FunctionsRegistrySession) RecoverFunds(to common.Address) (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.RecoverFunds(&_FunctionsRegistry.TransactOpts, to)
}

func (_FunctionsRegistry *FunctionsRegistryTransactorSession) RecoverFunds(to common.Address) (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.RecoverFunds(&_FunctionsRegistry.TransactOpts, to)
}

func (_FunctionsRegistry *FunctionsRegistryTransactor) RemoveConsumer(opts *bind.TransactOpts, subscriptionId uint64, consumer common.Address) (*types.Transaction, error) {
	return _FunctionsRegistry.contract.Transact(opts, "removeConsumer", subscriptionId, consumer)
}

func (_FunctionsRegistry *FunctionsRegistrySession) RemoveConsumer(subscriptionId uint64, consumer common.Address) (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.RemoveConsumer(&_FunctionsRegistry.TransactOpts, subscriptionId, consumer)
}

func (_FunctionsRegistry *FunctionsRegistryTransactorSession) RemoveConsumer(subscriptionId uint64, consumer common.Address) (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.RemoveConsumer(&_FunctionsRegistry.TransactOpts, subscriptionId, consumer)
}

func (_FunctionsRegistry *FunctionsRegistryTransactor) RequestSubscriptionOwnerTransfer(opts *bind.TransactOpts, subscriptionId uint64, newOwner common.Address) (*types.Transaction, error) {
	return _FunctionsRegistry.contract.Transact(opts, "requestSubscriptionOwnerTransfer", subscriptionId, newOwner)
}

func (_FunctionsRegistry *FunctionsRegistrySession) RequestSubscriptionOwnerTransfer(subscriptionId uint64, newOwner common.Address) (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.RequestSubscriptionOwnerTransfer(&_FunctionsRegistry.TransactOpts, subscriptionId, newOwner)
}

func (_FunctionsRegistry *FunctionsRegistryTransactorSession) RequestSubscriptionOwnerTransfer(subscriptionId uint64, newOwner common.Address) (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.RequestSubscriptionOwnerTransfer(&_FunctionsRegistry.TransactOpts, subscriptionId, newOwner)
}

func (_FunctionsRegistry *FunctionsRegistryTransactor) SetAuthorizedSenders(opts *bind.TransactOpts, senders []common.Address) (*types.Transaction, error) {
	return _FunctionsRegistry.contract.Transact(opts, "setAuthorizedSenders", senders)
}

func (_FunctionsRegistry *FunctionsRegistrySession) SetAuthorizedSenders(senders []common.Address) (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.SetAuthorizedSenders(&_FunctionsRegistry.TransactOpts, senders)
}

func (_FunctionsRegistry *FunctionsRegistryTransactorSession) SetAuthorizedSenders(senders []common.Address) (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.SetAuthorizedSenders(&_FunctionsRegistry.TransactOpts, senders)
}

func (_FunctionsRegistry *FunctionsRegistryTransactor) SetConfig(opts *bind.TransactOpts, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation *big.Int, fallbackWeiPerUnitLink *big.Int, gasOverhead uint32, requestTimeoutSeconds uint32) (*types.Transaction, error) {
	return _FunctionsRegistry.contract.Transact(opts, "setConfig", maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, gasOverhead, requestTimeoutSeconds)
}

func (_FunctionsRegistry *FunctionsRegistrySession) SetConfig(maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation *big.Int, fallbackWeiPerUnitLink *big.Int, gasOverhead uint32, requestTimeoutSeconds uint32) (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.SetConfig(&_FunctionsRegistry.TransactOpts, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, gasOverhead, requestTimeoutSeconds)
}

func (_FunctionsRegistry *FunctionsRegistryTransactorSession) SetConfig(maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation *big.Int, fallbackWeiPerUnitLink *big.Int, gasOverhead uint32, requestTimeoutSeconds uint32) (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.SetConfig(&_FunctionsRegistry.TransactOpts, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, gasOverhead, requestTimeoutSeconds)
}

func (_FunctionsRegistry *FunctionsRegistryTransactor) StartBilling(opts *bind.TransactOpts, data []byte, billing IFunctionsBillingRegistryRequestBilling) (*types.Transaction, error) {
	return _FunctionsRegistry.contract.Transact(opts, "startBilling", data, billing)
}

func (_FunctionsRegistry *FunctionsRegistrySession) StartBilling(data []byte, billing IFunctionsBillingRegistryRequestBilling) (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.StartBilling(&_FunctionsRegistry.TransactOpts, data, billing)
}

func (_FunctionsRegistry *FunctionsRegistryTransactorSession) StartBilling(data []byte, billing IFunctionsBillingRegistryRequestBilling) (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.StartBilling(&_FunctionsRegistry.TransactOpts, data, billing)
}

func (_FunctionsRegistry *FunctionsRegistryTransactor) TimeoutRequests(opts *bind.TransactOpts, requestIdsToTimeout [][32]byte) (*types.Transaction, error) {
	return _FunctionsRegistry.contract.Transact(opts, "timeoutRequests", requestIdsToTimeout)
}

func (_FunctionsRegistry *FunctionsRegistrySession) TimeoutRequests(requestIdsToTimeout [][32]byte) (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.TimeoutRequests(&_FunctionsRegistry.TransactOpts, requestIdsToTimeout)
}

func (_FunctionsRegistry *FunctionsRegistryTransactorSession) TimeoutRequests(requestIdsToTimeout [][32]byte) (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.TimeoutRequests(&_FunctionsRegistry.TransactOpts, requestIdsToTimeout)
}

func (_FunctionsRegistry *FunctionsRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _FunctionsRegistry.contract.Transact(opts, "transferOwnership", to)
}

func (_FunctionsRegistry *FunctionsRegistrySession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.TransferOwnership(&_FunctionsRegistry.TransactOpts, to)
}

func (_FunctionsRegistry *FunctionsRegistryTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.TransferOwnership(&_FunctionsRegistry.TransactOpts, to)
}

func (_FunctionsRegistry *FunctionsRegistryTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsRegistry.contract.Transact(opts, "unpause")
}

func (_FunctionsRegistry *FunctionsRegistrySession) Unpause() (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.Unpause(&_FunctionsRegistry.TransactOpts)
}

func (_FunctionsRegistry *FunctionsRegistryTransactorSession) Unpause() (*types.Transaction, error) {
	return _FunctionsRegistry.Contract.Unpause(&_FunctionsRegistry.TransactOpts)
}

type FunctionsRegistryAuthorizedSendersChangedIterator struct {
	Event *FunctionsRegistryAuthorizedSendersChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRegistryAuthorizedSendersChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRegistryAuthorizedSendersChanged)
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
		it.Event = new(FunctionsRegistryAuthorizedSendersChanged)
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

func (it *FunctionsRegistryAuthorizedSendersChangedIterator) Error() error {
	return it.fail
}

func (it *FunctionsRegistryAuthorizedSendersChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRegistryAuthorizedSendersChanged struct {
	Senders   []common.Address
	ChangedBy common.Address
	Raw       types.Log
}

func (_FunctionsRegistry *FunctionsRegistryFilterer) FilterAuthorizedSendersChanged(opts *bind.FilterOpts) (*FunctionsRegistryAuthorizedSendersChangedIterator, error) {

	logs, sub, err := _FunctionsRegistry.contract.FilterLogs(opts, "AuthorizedSendersChanged")
	if err != nil {
		return nil, err
	}
	return &FunctionsRegistryAuthorizedSendersChangedIterator{contract: _FunctionsRegistry.contract, event: "AuthorizedSendersChanged", logs: logs, sub: sub}, nil
}

func (_FunctionsRegistry *FunctionsRegistryFilterer) WatchAuthorizedSendersChanged(opts *bind.WatchOpts, sink chan<- *FunctionsRegistryAuthorizedSendersChanged) (event.Subscription, error) {

	logs, sub, err := _FunctionsRegistry.contract.WatchLogs(opts, "AuthorizedSendersChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRegistryAuthorizedSendersChanged)
				if err := _FunctionsRegistry.contract.UnpackLog(event, "AuthorizedSendersChanged", log); err != nil {
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

func (_FunctionsRegistry *FunctionsRegistryFilterer) ParseAuthorizedSendersChanged(log types.Log) (*FunctionsRegistryAuthorizedSendersChanged, error) {
	event := new(FunctionsRegistryAuthorizedSendersChanged)
	if err := _FunctionsRegistry.contract.UnpackLog(event, "AuthorizedSendersChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsRegistryBillingEndIterator struct {
	Event *FunctionsRegistryBillingEnd

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRegistryBillingEndIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRegistryBillingEnd)
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
		it.Event = new(FunctionsRegistryBillingEnd)
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

func (it *FunctionsRegistryBillingEndIterator) Error() error {
	return it.fail
}

func (it *FunctionsRegistryBillingEndIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRegistryBillingEnd struct {
	RequestId          [32]byte
	SubscriptionId     uint64
	SignerPayment      *big.Int
	TransmitterPayment *big.Int
	TotalCost          *big.Int
	Success            bool
	Raw                types.Log
}

func (_FunctionsRegistry *FunctionsRegistryFilterer) FilterBillingEnd(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsRegistryBillingEndIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsRegistry.contract.FilterLogs(opts, "BillingEnd", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsRegistryBillingEndIterator{contract: _FunctionsRegistry.contract, event: "BillingEnd", logs: logs, sub: sub}, nil
}

func (_FunctionsRegistry *FunctionsRegistryFilterer) WatchBillingEnd(opts *bind.WatchOpts, sink chan<- *FunctionsRegistryBillingEnd, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsRegistry.contract.WatchLogs(opts, "BillingEnd", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRegistryBillingEnd)
				if err := _FunctionsRegistry.contract.UnpackLog(event, "BillingEnd", log); err != nil {
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

func (_FunctionsRegistry *FunctionsRegistryFilterer) ParseBillingEnd(log types.Log) (*FunctionsRegistryBillingEnd, error) {
	event := new(FunctionsRegistryBillingEnd)
	if err := _FunctionsRegistry.contract.UnpackLog(event, "BillingEnd", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsRegistryBillingStartIterator struct {
	Event *FunctionsRegistryBillingStart

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRegistryBillingStartIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRegistryBillingStart)
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
		it.Event = new(FunctionsRegistryBillingStart)
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

func (it *FunctionsRegistryBillingStartIterator) Error() error {
	return it.fail
}

func (it *FunctionsRegistryBillingStartIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRegistryBillingStart struct {
	RequestId  [32]byte
	Commitment FunctionsBillingRegistryCommitment
	Raw        types.Log
}

func (_FunctionsRegistry *FunctionsRegistryFilterer) FilterBillingStart(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsRegistryBillingStartIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsRegistry.contract.FilterLogs(opts, "BillingStart", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsRegistryBillingStartIterator{contract: _FunctionsRegistry.contract, event: "BillingStart", logs: logs, sub: sub}, nil
}

func (_FunctionsRegistry *FunctionsRegistryFilterer) WatchBillingStart(opts *bind.WatchOpts, sink chan<- *FunctionsRegistryBillingStart, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsRegistry.contract.WatchLogs(opts, "BillingStart", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRegistryBillingStart)
				if err := _FunctionsRegistry.contract.UnpackLog(event, "BillingStart", log); err != nil {
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

func (_FunctionsRegistry *FunctionsRegistryFilterer) ParseBillingStart(log types.Log) (*FunctionsRegistryBillingStart, error) {
	event := new(FunctionsRegistryBillingStart)
	if err := _FunctionsRegistry.contract.UnpackLog(event, "BillingStart", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsRegistryConfigSetIterator struct {
	Event *FunctionsRegistryConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRegistryConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRegistryConfigSet)
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
		it.Event = new(FunctionsRegistryConfigSet)
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

func (it *FunctionsRegistryConfigSetIterator) Error() error {
	return it.fail
}

func (it *FunctionsRegistryConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRegistryConfigSet struct {
	MaxGasLimit                uint32
	StalenessSeconds           uint32
	GasAfterPaymentCalculation *big.Int
	FallbackWeiPerUnitLink     *big.Int
	GasOverhead                uint32
	Raw                        types.Log
}

func (_FunctionsRegistry *FunctionsRegistryFilterer) FilterConfigSet(opts *bind.FilterOpts) (*FunctionsRegistryConfigSetIterator, error) {

	logs, sub, err := _FunctionsRegistry.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &FunctionsRegistryConfigSetIterator{contract: _FunctionsRegistry.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_FunctionsRegistry *FunctionsRegistryFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *FunctionsRegistryConfigSet) (event.Subscription, error) {

	logs, sub, err := _FunctionsRegistry.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRegistryConfigSet)
				if err := _FunctionsRegistry.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_FunctionsRegistry *FunctionsRegistryFilterer) ParseConfigSet(log types.Log) (*FunctionsRegistryConfigSet, error) {
	event := new(FunctionsRegistryConfigSet)
	if err := _FunctionsRegistry.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsRegistryFundsRecoveredIterator struct {
	Event *FunctionsRegistryFundsRecovered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRegistryFundsRecoveredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRegistryFundsRecovered)
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
		it.Event = new(FunctionsRegistryFundsRecovered)
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

func (it *FunctionsRegistryFundsRecoveredIterator) Error() error {
	return it.fail
}

func (it *FunctionsRegistryFundsRecoveredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRegistryFundsRecovered struct {
	To     common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_FunctionsRegistry *FunctionsRegistryFilterer) FilterFundsRecovered(opts *bind.FilterOpts) (*FunctionsRegistryFundsRecoveredIterator, error) {

	logs, sub, err := _FunctionsRegistry.contract.FilterLogs(opts, "FundsRecovered")
	if err != nil {
		return nil, err
	}
	return &FunctionsRegistryFundsRecoveredIterator{contract: _FunctionsRegistry.contract, event: "FundsRecovered", logs: logs, sub: sub}, nil
}

func (_FunctionsRegistry *FunctionsRegistryFilterer) WatchFundsRecovered(opts *bind.WatchOpts, sink chan<- *FunctionsRegistryFundsRecovered) (event.Subscription, error) {

	logs, sub, err := _FunctionsRegistry.contract.WatchLogs(opts, "FundsRecovered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRegistryFundsRecovered)
				if err := _FunctionsRegistry.contract.UnpackLog(event, "FundsRecovered", log); err != nil {
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

func (_FunctionsRegistry *FunctionsRegistryFilterer) ParseFundsRecovered(log types.Log) (*FunctionsRegistryFundsRecovered, error) {
	event := new(FunctionsRegistryFundsRecovered)
	if err := _FunctionsRegistry.contract.UnpackLog(event, "FundsRecovered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsRegistryInitializedIterator struct {
	Event *FunctionsRegistryInitialized

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRegistryInitializedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRegistryInitialized)
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
		it.Event = new(FunctionsRegistryInitialized)
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

func (it *FunctionsRegistryInitializedIterator) Error() error {
	return it.fail
}

func (it *FunctionsRegistryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRegistryInitialized struct {
	Version uint8
	Raw     types.Log
}

func (_FunctionsRegistry *FunctionsRegistryFilterer) FilterInitialized(opts *bind.FilterOpts) (*FunctionsRegistryInitializedIterator, error) {

	logs, sub, err := _FunctionsRegistry.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &FunctionsRegistryInitializedIterator{contract: _FunctionsRegistry.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

func (_FunctionsRegistry *FunctionsRegistryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *FunctionsRegistryInitialized) (event.Subscription, error) {

	logs, sub, err := _FunctionsRegistry.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRegistryInitialized)
				if err := _FunctionsRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
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

func (_FunctionsRegistry *FunctionsRegistryFilterer) ParseInitialized(log types.Log) (*FunctionsRegistryInitialized, error) {
	event := new(FunctionsRegistryInitialized)
	if err := _FunctionsRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsRegistryOwnershipTransferRequestedIterator struct {
	Event *FunctionsRegistryOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRegistryOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRegistryOwnershipTransferRequested)
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
		it.Event = new(FunctionsRegistryOwnershipTransferRequested)
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

func (it *FunctionsRegistryOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *FunctionsRegistryOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRegistryOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_FunctionsRegistry *FunctionsRegistryFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsRegistryOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsRegistry.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsRegistryOwnershipTransferRequestedIterator{contract: _FunctionsRegistry.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_FunctionsRegistry *FunctionsRegistryFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *FunctionsRegistryOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsRegistry.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRegistryOwnershipTransferRequested)
				if err := _FunctionsRegistry.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_FunctionsRegistry *FunctionsRegistryFilterer) ParseOwnershipTransferRequested(log types.Log) (*FunctionsRegistryOwnershipTransferRequested, error) {
	event := new(FunctionsRegistryOwnershipTransferRequested)
	if err := _FunctionsRegistry.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsRegistryOwnershipTransferredIterator struct {
	Event *FunctionsRegistryOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRegistryOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRegistryOwnershipTransferred)
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
		it.Event = new(FunctionsRegistryOwnershipTransferred)
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

func (it *FunctionsRegistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *FunctionsRegistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRegistryOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_FunctionsRegistry *FunctionsRegistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsRegistryOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsRegistry.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsRegistryOwnershipTransferredIterator{contract: _FunctionsRegistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_FunctionsRegistry *FunctionsRegistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FunctionsRegistryOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsRegistry.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRegistryOwnershipTransferred)
				if err := _FunctionsRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_FunctionsRegistry *FunctionsRegistryFilterer) ParseOwnershipTransferred(log types.Log) (*FunctionsRegistryOwnershipTransferred, error) {
	event := new(FunctionsRegistryOwnershipTransferred)
	if err := _FunctionsRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsRegistryPausedIterator struct {
	Event *FunctionsRegistryPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRegistryPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRegistryPaused)
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
		it.Event = new(FunctionsRegistryPaused)
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

func (it *FunctionsRegistryPausedIterator) Error() error {
	return it.fail
}

func (it *FunctionsRegistryPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRegistryPaused struct {
	Account common.Address
	Raw     types.Log
}

func (_FunctionsRegistry *FunctionsRegistryFilterer) FilterPaused(opts *bind.FilterOpts) (*FunctionsRegistryPausedIterator, error) {

	logs, sub, err := _FunctionsRegistry.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &FunctionsRegistryPausedIterator{contract: _FunctionsRegistry.contract, event: "Paused", logs: logs, sub: sub}, nil
}

func (_FunctionsRegistry *FunctionsRegistryFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *FunctionsRegistryPaused) (event.Subscription, error) {

	logs, sub, err := _FunctionsRegistry.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRegistryPaused)
				if err := _FunctionsRegistry.contract.UnpackLog(event, "Paused", log); err != nil {
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

func (_FunctionsRegistry *FunctionsRegistryFilterer) ParsePaused(log types.Log) (*FunctionsRegistryPaused, error) {
	event := new(FunctionsRegistryPaused)
	if err := _FunctionsRegistry.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsRegistryRequestTimedOutIterator struct {
	Event *FunctionsRegistryRequestTimedOut

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRegistryRequestTimedOutIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRegistryRequestTimedOut)
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
		it.Event = new(FunctionsRegistryRequestTimedOut)
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

func (it *FunctionsRegistryRequestTimedOutIterator) Error() error {
	return it.fail
}

func (it *FunctionsRegistryRequestTimedOutIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRegistryRequestTimedOut struct {
	RequestId [32]byte
	Raw       types.Log
}

func (_FunctionsRegistry *FunctionsRegistryFilterer) FilterRequestTimedOut(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsRegistryRequestTimedOutIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsRegistry.contract.FilterLogs(opts, "RequestTimedOut", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsRegistryRequestTimedOutIterator{contract: _FunctionsRegistry.contract, event: "RequestTimedOut", logs: logs, sub: sub}, nil
}

func (_FunctionsRegistry *FunctionsRegistryFilterer) WatchRequestTimedOut(opts *bind.WatchOpts, sink chan<- *FunctionsRegistryRequestTimedOut, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsRegistry.contract.WatchLogs(opts, "RequestTimedOut", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRegistryRequestTimedOut)
				if err := _FunctionsRegistry.contract.UnpackLog(event, "RequestTimedOut", log); err != nil {
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

func (_FunctionsRegistry *FunctionsRegistryFilterer) ParseRequestTimedOut(log types.Log) (*FunctionsRegistryRequestTimedOut, error) {
	event := new(FunctionsRegistryRequestTimedOut)
	if err := _FunctionsRegistry.contract.UnpackLog(event, "RequestTimedOut", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsRegistrySubscriptionCanceledIterator struct {
	Event *FunctionsRegistrySubscriptionCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRegistrySubscriptionCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRegistrySubscriptionCanceled)
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
		it.Event = new(FunctionsRegistrySubscriptionCanceled)
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

func (it *FunctionsRegistrySubscriptionCanceledIterator) Error() error {
	return it.fail
}

func (it *FunctionsRegistrySubscriptionCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRegistrySubscriptionCanceled struct {
	SubscriptionId uint64
	To             common.Address
	Amount         *big.Int
	Raw            types.Log
}

func (_FunctionsRegistry *FunctionsRegistryFilterer) FilterSubscriptionCanceled(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsRegistrySubscriptionCanceledIterator, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsRegistry.contract.FilterLogs(opts, "SubscriptionCanceled", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsRegistrySubscriptionCanceledIterator{contract: _FunctionsRegistry.contract, event: "SubscriptionCanceled", logs: logs, sub: sub}, nil
}

func (_FunctionsRegistry *FunctionsRegistryFilterer) WatchSubscriptionCanceled(opts *bind.WatchOpts, sink chan<- *FunctionsRegistrySubscriptionCanceled, subscriptionId []uint64) (event.Subscription, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsRegistry.contract.WatchLogs(opts, "SubscriptionCanceled", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRegistrySubscriptionCanceled)
				if err := _FunctionsRegistry.contract.UnpackLog(event, "SubscriptionCanceled", log); err != nil {
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

func (_FunctionsRegistry *FunctionsRegistryFilterer) ParseSubscriptionCanceled(log types.Log) (*FunctionsRegistrySubscriptionCanceled, error) {
	event := new(FunctionsRegistrySubscriptionCanceled)
	if err := _FunctionsRegistry.contract.UnpackLog(event, "SubscriptionCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsRegistrySubscriptionConsumerAddedIterator struct {
	Event *FunctionsRegistrySubscriptionConsumerAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRegistrySubscriptionConsumerAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRegistrySubscriptionConsumerAdded)
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
		it.Event = new(FunctionsRegistrySubscriptionConsumerAdded)
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

func (it *FunctionsRegistrySubscriptionConsumerAddedIterator) Error() error {
	return it.fail
}

func (it *FunctionsRegistrySubscriptionConsumerAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRegistrySubscriptionConsumerAdded struct {
	SubscriptionId uint64
	Consumer       common.Address
	Raw            types.Log
}

func (_FunctionsRegistry *FunctionsRegistryFilterer) FilterSubscriptionConsumerAdded(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsRegistrySubscriptionConsumerAddedIterator, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsRegistry.contract.FilterLogs(opts, "SubscriptionConsumerAdded", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsRegistrySubscriptionConsumerAddedIterator{contract: _FunctionsRegistry.contract, event: "SubscriptionConsumerAdded", logs: logs, sub: sub}, nil
}

func (_FunctionsRegistry *FunctionsRegistryFilterer) WatchSubscriptionConsumerAdded(opts *bind.WatchOpts, sink chan<- *FunctionsRegistrySubscriptionConsumerAdded, subscriptionId []uint64) (event.Subscription, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsRegistry.contract.WatchLogs(opts, "SubscriptionConsumerAdded", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRegistrySubscriptionConsumerAdded)
				if err := _FunctionsRegistry.contract.UnpackLog(event, "SubscriptionConsumerAdded", log); err != nil {
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

func (_FunctionsRegistry *FunctionsRegistryFilterer) ParseSubscriptionConsumerAdded(log types.Log) (*FunctionsRegistrySubscriptionConsumerAdded, error) {
	event := new(FunctionsRegistrySubscriptionConsumerAdded)
	if err := _FunctionsRegistry.contract.UnpackLog(event, "SubscriptionConsumerAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsRegistrySubscriptionConsumerRemovedIterator struct {
	Event *FunctionsRegistrySubscriptionConsumerRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRegistrySubscriptionConsumerRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRegistrySubscriptionConsumerRemoved)
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
		it.Event = new(FunctionsRegistrySubscriptionConsumerRemoved)
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

func (it *FunctionsRegistrySubscriptionConsumerRemovedIterator) Error() error {
	return it.fail
}

func (it *FunctionsRegistrySubscriptionConsumerRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRegistrySubscriptionConsumerRemoved struct {
	SubscriptionId uint64
	Consumer       common.Address
	Raw            types.Log
}

func (_FunctionsRegistry *FunctionsRegistryFilterer) FilterSubscriptionConsumerRemoved(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsRegistrySubscriptionConsumerRemovedIterator, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsRegistry.contract.FilterLogs(opts, "SubscriptionConsumerRemoved", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsRegistrySubscriptionConsumerRemovedIterator{contract: _FunctionsRegistry.contract, event: "SubscriptionConsumerRemoved", logs: logs, sub: sub}, nil
}

func (_FunctionsRegistry *FunctionsRegistryFilterer) WatchSubscriptionConsumerRemoved(opts *bind.WatchOpts, sink chan<- *FunctionsRegistrySubscriptionConsumerRemoved, subscriptionId []uint64) (event.Subscription, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsRegistry.contract.WatchLogs(opts, "SubscriptionConsumerRemoved", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRegistrySubscriptionConsumerRemoved)
				if err := _FunctionsRegistry.contract.UnpackLog(event, "SubscriptionConsumerRemoved", log); err != nil {
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

func (_FunctionsRegistry *FunctionsRegistryFilterer) ParseSubscriptionConsumerRemoved(log types.Log) (*FunctionsRegistrySubscriptionConsumerRemoved, error) {
	event := new(FunctionsRegistrySubscriptionConsumerRemoved)
	if err := _FunctionsRegistry.contract.UnpackLog(event, "SubscriptionConsumerRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsRegistrySubscriptionCreatedIterator struct {
	Event *FunctionsRegistrySubscriptionCreated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRegistrySubscriptionCreatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRegistrySubscriptionCreated)
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
		it.Event = new(FunctionsRegistrySubscriptionCreated)
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

func (it *FunctionsRegistrySubscriptionCreatedIterator) Error() error {
	return it.fail
}

func (it *FunctionsRegistrySubscriptionCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRegistrySubscriptionCreated struct {
	SubscriptionId uint64
	Owner          common.Address
	Raw            types.Log
}

func (_FunctionsRegistry *FunctionsRegistryFilterer) FilterSubscriptionCreated(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsRegistrySubscriptionCreatedIterator, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsRegistry.contract.FilterLogs(opts, "SubscriptionCreated", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsRegistrySubscriptionCreatedIterator{contract: _FunctionsRegistry.contract, event: "SubscriptionCreated", logs: logs, sub: sub}, nil
}

func (_FunctionsRegistry *FunctionsRegistryFilterer) WatchSubscriptionCreated(opts *bind.WatchOpts, sink chan<- *FunctionsRegistrySubscriptionCreated, subscriptionId []uint64) (event.Subscription, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsRegistry.contract.WatchLogs(opts, "SubscriptionCreated", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRegistrySubscriptionCreated)
				if err := _FunctionsRegistry.contract.UnpackLog(event, "SubscriptionCreated", log); err != nil {
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

func (_FunctionsRegistry *FunctionsRegistryFilterer) ParseSubscriptionCreated(log types.Log) (*FunctionsRegistrySubscriptionCreated, error) {
	event := new(FunctionsRegistrySubscriptionCreated)
	if err := _FunctionsRegistry.contract.UnpackLog(event, "SubscriptionCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsRegistrySubscriptionFundedIterator struct {
	Event *FunctionsRegistrySubscriptionFunded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRegistrySubscriptionFundedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRegistrySubscriptionFunded)
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
		it.Event = new(FunctionsRegistrySubscriptionFunded)
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

func (it *FunctionsRegistrySubscriptionFundedIterator) Error() error {
	return it.fail
}

func (it *FunctionsRegistrySubscriptionFundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRegistrySubscriptionFunded struct {
	SubscriptionId uint64
	OldBalance     *big.Int
	NewBalance     *big.Int
	Raw            types.Log
}

func (_FunctionsRegistry *FunctionsRegistryFilterer) FilterSubscriptionFunded(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsRegistrySubscriptionFundedIterator, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsRegistry.contract.FilterLogs(opts, "SubscriptionFunded", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsRegistrySubscriptionFundedIterator{contract: _FunctionsRegistry.contract, event: "SubscriptionFunded", logs: logs, sub: sub}, nil
}

func (_FunctionsRegistry *FunctionsRegistryFilterer) WatchSubscriptionFunded(opts *bind.WatchOpts, sink chan<- *FunctionsRegistrySubscriptionFunded, subscriptionId []uint64) (event.Subscription, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsRegistry.contract.WatchLogs(opts, "SubscriptionFunded", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRegistrySubscriptionFunded)
				if err := _FunctionsRegistry.contract.UnpackLog(event, "SubscriptionFunded", log); err != nil {
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

func (_FunctionsRegistry *FunctionsRegistryFilterer) ParseSubscriptionFunded(log types.Log) (*FunctionsRegistrySubscriptionFunded, error) {
	event := new(FunctionsRegistrySubscriptionFunded)
	if err := _FunctionsRegistry.contract.UnpackLog(event, "SubscriptionFunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsRegistrySubscriptionOwnerTransferRequestedIterator struct {
	Event *FunctionsRegistrySubscriptionOwnerTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRegistrySubscriptionOwnerTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRegistrySubscriptionOwnerTransferRequested)
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
		it.Event = new(FunctionsRegistrySubscriptionOwnerTransferRequested)
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

func (it *FunctionsRegistrySubscriptionOwnerTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *FunctionsRegistrySubscriptionOwnerTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRegistrySubscriptionOwnerTransferRequested struct {
	SubscriptionId uint64
	From           common.Address
	To             common.Address
	Raw            types.Log
}

func (_FunctionsRegistry *FunctionsRegistryFilterer) FilterSubscriptionOwnerTransferRequested(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsRegistrySubscriptionOwnerTransferRequestedIterator, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsRegistry.contract.FilterLogs(opts, "SubscriptionOwnerTransferRequested", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsRegistrySubscriptionOwnerTransferRequestedIterator{contract: _FunctionsRegistry.contract, event: "SubscriptionOwnerTransferRequested", logs: logs, sub: sub}, nil
}

func (_FunctionsRegistry *FunctionsRegistryFilterer) WatchSubscriptionOwnerTransferRequested(opts *bind.WatchOpts, sink chan<- *FunctionsRegistrySubscriptionOwnerTransferRequested, subscriptionId []uint64) (event.Subscription, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsRegistry.contract.WatchLogs(opts, "SubscriptionOwnerTransferRequested", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRegistrySubscriptionOwnerTransferRequested)
				if err := _FunctionsRegistry.contract.UnpackLog(event, "SubscriptionOwnerTransferRequested", log); err != nil {
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

func (_FunctionsRegistry *FunctionsRegistryFilterer) ParseSubscriptionOwnerTransferRequested(log types.Log) (*FunctionsRegistrySubscriptionOwnerTransferRequested, error) {
	event := new(FunctionsRegistrySubscriptionOwnerTransferRequested)
	if err := _FunctionsRegistry.contract.UnpackLog(event, "SubscriptionOwnerTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsRegistrySubscriptionOwnerTransferredIterator struct {
	Event *FunctionsRegistrySubscriptionOwnerTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRegistrySubscriptionOwnerTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRegistrySubscriptionOwnerTransferred)
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
		it.Event = new(FunctionsRegistrySubscriptionOwnerTransferred)
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

func (it *FunctionsRegistrySubscriptionOwnerTransferredIterator) Error() error {
	return it.fail
}

func (it *FunctionsRegistrySubscriptionOwnerTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRegistrySubscriptionOwnerTransferred struct {
	SubscriptionId uint64
	From           common.Address
	To             common.Address
	Raw            types.Log
}

func (_FunctionsRegistry *FunctionsRegistryFilterer) FilterSubscriptionOwnerTransferred(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsRegistrySubscriptionOwnerTransferredIterator, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsRegistry.contract.FilterLogs(opts, "SubscriptionOwnerTransferred", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsRegistrySubscriptionOwnerTransferredIterator{contract: _FunctionsRegistry.contract, event: "SubscriptionOwnerTransferred", logs: logs, sub: sub}, nil
}

func (_FunctionsRegistry *FunctionsRegistryFilterer) WatchSubscriptionOwnerTransferred(opts *bind.WatchOpts, sink chan<- *FunctionsRegistrySubscriptionOwnerTransferred, subscriptionId []uint64) (event.Subscription, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsRegistry.contract.WatchLogs(opts, "SubscriptionOwnerTransferred", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRegistrySubscriptionOwnerTransferred)
				if err := _FunctionsRegistry.contract.UnpackLog(event, "SubscriptionOwnerTransferred", log); err != nil {
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

func (_FunctionsRegistry *FunctionsRegistryFilterer) ParseSubscriptionOwnerTransferred(log types.Log) (*FunctionsRegistrySubscriptionOwnerTransferred, error) {
	event := new(FunctionsRegistrySubscriptionOwnerTransferred)
	if err := _FunctionsRegistry.contract.UnpackLog(event, "SubscriptionOwnerTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsRegistryUnpausedIterator struct {
	Event *FunctionsRegistryUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsRegistryUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsRegistryUnpaused)
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
		it.Event = new(FunctionsRegistryUnpaused)
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

func (it *FunctionsRegistryUnpausedIterator) Error() error {
	return it.fail
}

func (it *FunctionsRegistryUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsRegistryUnpaused struct {
	Account common.Address
	Raw     types.Log
}

func (_FunctionsRegistry *FunctionsRegistryFilterer) FilterUnpaused(opts *bind.FilterOpts) (*FunctionsRegistryUnpausedIterator, error) {

	logs, sub, err := _FunctionsRegistry.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &FunctionsRegistryUnpausedIterator{contract: _FunctionsRegistry.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

func (_FunctionsRegistry *FunctionsRegistryFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *FunctionsRegistryUnpaused) (event.Subscription, error) {

	logs, sub, err := _FunctionsRegistry.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsRegistryUnpaused)
				if err := _FunctionsRegistry.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

func (_FunctionsRegistry *FunctionsRegistryFilterer) ParseUnpaused(log types.Log) (*FunctionsRegistryUnpaused, error) {
	event := new(FunctionsRegistryUnpaused)
	if err := _FunctionsRegistry.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

func (_FunctionsRegistry *FunctionsRegistry) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _FunctionsRegistry.abi.Events["AuthorizedSendersChanged"].ID:
		return _FunctionsRegistry.ParseAuthorizedSendersChanged(log)
	case _FunctionsRegistry.abi.Events["BillingEnd"].ID:
		return _FunctionsRegistry.ParseBillingEnd(log)
	case _FunctionsRegistry.abi.Events["BillingStart"].ID:
		return _FunctionsRegistry.ParseBillingStart(log)
	case _FunctionsRegistry.abi.Events["ConfigSet"].ID:
		return _FunctionsRegistry.ParseConfigSet(log)
	case _FunctionsRegistry.abi.Events["FundsRecovered"].ID:
		return _FunctionsRegistry.ParseFundsRecovered(log)
	case _FunctionsRegistry.abi.Events["Initialized"].ID:
		return _FunctionsRegistry.ParseInitialized(log)
	case _FunctionsRegistry.abi.Events["OwnershipTransferRequested"].ID:
		return _FunctionsRegistry.ParseOwnershipTransferRequested(log)
	case _FunctionsRegistry.abi.Events["OwnershipTransferred"].ID:
		return _FunctionsRegistry.ParseOwnershipTransferred(log)
	case _FunctionsRegistry.abi.Events["Paused"].ID:
		return _FunctionsRegistry.ParsePaused(log)
	case _FunctionsRegistry.abi.Events["RequestTimedOut"].ID:
		return _FunctionsRegistry.ParseRequestTimedOut(log)
	case _FunctionsRegistry.abi.Events["SubscriptionCanceled"].ID:
		return _FunctionsRegistry.ParseSubscriptionCanceled(log)
	case _FunctionsRegistry.abi.Events["SubscriptionConsumerAdded"].ID:
		return _FunctionsRegistry.ParseSubscriptionConsumerAdded(log)
	case _FunctionsRegistry.abi.Events["SubscriptionConsumerRemoved"].ID:
		return _FunctionsRegistry.ParseSubscriptionConsumerRemoved(log)
	case _FunctionsRegistry.abi.Events["SubscriptionCreated"].ID:
		return _FunctionsRegistry.ParseSubscriptionCreated(log)
	case _FunctionsRegistry.abi.Events["SubscriptionFunded"].ID:
		return _FunctionsRegistry.ParseSubscriptionFunded(log)
	case _FunctionsRegistry.abi.Events["SubscriptionOwnerTransferRequested"].ID:
		return _FunctionsRegistry.ParseSubscriptionOwnerTransferRequested(log)
	case _FunctionsRegistry.abi.Events["SubscriptionOwnerTransferred"].ID:
		return _FunctionsRegistry.ParseSubscriptionOwnerTransferred(log)
	case _FunctionsRegistry.abi.Events["Unpaused"].ID:
		return _FunctionsRegistry.ParseUnpaused(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (FunctionsRegistryAuthorizedSendersChanged) Topic() common.Hash {
	return common.HexToHash("0xf263cfb3e4298332e776194610cf9fdc09ccb3ada8b9aa39764d882e11fbf0a0")
}

func (FunctionsRegistryBillingEnd) Topic() common.Hash {
	return common.HexToHash("0xc8dc973332de19a5f71b6026983110e9c2e04b0c98b87eb771ccb78607fd114f")
}

func (FunctionsRegistryBillingStart) Topic() common.Hash {
	return common.HexToHash("0x99f7f4e65b4b9fbabd4e357c47ed3099b36e57ecd3a43e84662f34c207d0ebe4")
}

func (FunctionsRegistryConfigSet) Topic() common.Hash {
	return common.HexToHash("0x24d3d934adfef9b9029d6ffa463c07d0139ed47d26ee23506f85ece2879d2bd4")
}

func (FunctionsRegistryFundsRecovered) Topic() common.Hash {
	return common.HexToHash("0x59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b436600")
}

func (FunctionsRegistryInitialized) Topic() common.Hash {
	return common.HexToHash("0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498")
}

func (FunctionsRegistryOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (FunctionsRegistryOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (FunctionsRegistryPaused) Topic() common.Hash {
	return common.HexToHash("0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258")
}

func (FunctionsRegistryRequestTimedOut) Topic() common.Hash {
	return common.HexToHash("0xf1ca1e9147be737b04a2b018a79405f687a97de8dd8a2559bbe62357343af414")
}

func (FunctionsRegistrySubscriptionCanceled) Topic() common.Hash {
	return common.HexToHash("0xe8ed5b475a5b5987aa9165e8731bb78043f39eee32ec5a1169a89e27fcd49815")
}

func (FunctionsRegistrySubscriptionConsumerAdded) Topic() common.Hash {
	return common.HexToHash("0x43dc749a04ac8fb825cbd514f7c0e13f13bc6f2ee66043b76629d51776cff8e0")
}

func (FunctionsRegistrySubscriptionConsumerRemoved) Topic() common.Hash {
	return common.HexToHash("0x182bff9831466789164ca77075fffd84916d35a8180ba73c27e45634549b445b")
}

func (FunctionsRegistrySubscriptionCreated) Topic() common.Hash {
	return common.HexToHash("0x464722b4166576d3dcbba877b999bc35cf911f4eaf434b7eba68fa113951d0bf")
}

func (FunctionsRegistrySubscriptionFunded) Topic() common.Hash {
	return common.HexToHash("0xd39ec07f4e209f627a4c427971473820dc129761ba28de8906bd56f57101d4f8")
}

func (FunctionsRegistrySubscriptionOwnerTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x69436ea6df009049404f564eff6622cd00522b0bd6a89efd9e52a355c4a879be")
}

func (FunctionsRegistrySubscriptionOwnerTransferred) Topic() common.Hash {
	return common.HexToHash("0x6f1dc65165ffffedfd8e507b4a0f1fcfdada045ed11f6c26ba27cedfe87802f0")
}

func (FunctionsRegistryUnpaused) Topic() common.Hash {
	return common.HexToHash("0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa")
}

func (_FunctionsRegistry *FunctionsRegistry) Address() common.Address {
	return _FunctionsRegistry.address
}

type FunctionsRegistryInterface interface {
	MAXCONSUMERS(opts *bind.CallOpts) (uint16, error)

	EstimateCost(opts *bind.CallOpts, gasLimit uint32, gasPrice *big.Int, donFee *big.Int, registryFee *big.Int) (*big.Int, error)

	GetAuthorizedSenders(opts *bind.CallOpts) ([]common.Address, error)

	GetConfig(opts *bind.CallOpts) (GetConfig,

		error)

	GetCurrentsubscriptionId(opts *bind.CallOpts) (uint64, error)

	GetRequestConfig(opts *bind.CallOpts) (uint32, []common.Address, error)

	GetRequiredFee(opts *bind.CallOpts, arg0 []byte, arg1 IFunctionsBillingRegistryRequestBilling) (*big.Int, error)

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

	SetConfig(opts *bind.TransactOpts, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation *big.Int, fallbackWeiPerUnitLink *big.Int, gasOverhead uint32, requestTimeoutSeconds uint32) (*types.Transaction, error)

	StartBilling(opts *bind.TransactOpts, data []byte, billing IFunctionsBillingRegistryRequestBilling) (*types.Transaction, error)

	TimeoutRequests(opts *bind.TransactOpts, requestIdsToTimeout [][32]byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Unpause(opts *bind.TransactOpts) (*types.Transaction, error)

	FilterAuthorizedSendersChanged(opts *bind.FilterOpts) (*FunctionsRegistryAuthorizedSendersChangedIterator, error)

	WatchAuthorizedSendersChanged(opts *bind.WatchOpts, sink chan<- *FunctionsRegistryAuthorizedSendersChanged) (event.Subscription, error)

	ParseAuthorizedSendersChanged(log types.Log) (*FunctionsRegistryAuthorizedSendersChanged, error)

	FilterBillingEnd(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsRegistryBillingEndIterator, error)

	WatchBillingEnd(opts *bind.WatchOpts, sink chan<- *FunctionsRegistryBillingEnd, requestId [][32]byte) (event.Subscription, error)

	ParseBillingEnd(log types.Log) (*FunctionsRegistryBillingEnd, error)

	FilterBillingStart(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsRegistryBillingStartIterator, error)

	WatchBillingStart(opts *bind.WatchOpts, sink chan<- *FunctionsRegistryBillingStart, requestId [][32]byte) (event.Subscription, error)

	ParseBillingStart(log types.Log) (*FunctionsRegistryBillingStart, error)

	FilterConfigSet(opts *bind.FilterOpts) (*FunctionsRegistryConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *FunctionsRegistryConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*FunctionsRegistryConfigSet, error)

	FilterFundsRecovered(opts *bind.FilterOpts) (*FunctionsRegistryFundsRecoveredIterator, error)

	WatchFundsRecovered(opts *bind.WatchOpts, sink chan<- *FunctionsRegistryFundsRecovered) (event.Subscription, error)

	ParseFundsRecovered(log types.Log) (*FunctionsRegistryFundsRecovered, error)

	FilterInitialized(opts *bind.FilterOpts) (*FunctionsRegistryInitializedIterator, error)

	WatchInitialized(opts *bind.WatchOpts, sink chan<- *FunctionsRegistryInitialized) (event.Subscription, error)

	ParseInitialized(log types.Log) (*FunctionsRegistryInitialized, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsRegistryOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *FunctionsRegistryOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*FunctionsRegistryOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsRegistryOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FunctionsRegistryOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*FunctionsRegistryOwnershipTransferred, error)

	FilterPaused(opts *bind.FilterOpts) (*FunctionsRegistryPausedIterator, error)

	WatchPaused(opts *bind.WatchOpts, sink chan<- *FunctionsRegistryPaused) (event.Subscription, error)

	ParsePaused(log types.Log) (*FunctionsRegistryPaused, error)

	FilterRequestTimedOut(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsRegistryRequestTimedOutIterator, error)

	WatchRequestTimedOut(opts *bind.WatchOpts, sink chan<- *FunctionsRegistryRequestTimedOut, requestId [][32]byte) (event.Subscription, error)

	ParseRequestTimedOut(log types.Log) (*FunctionsRegistryRequestTimedOut, error)

	FilterSubscriptionCanceled(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsRegistrySubscriptionCanceledIterator, error)

	WatchSubscriptionCanceled(opts *bind.WatchOpts, sink chan<- *FunctionsRegistrySubscriptionCanceled, subscriptionId []uint64) (event.Subscription, error)

	ParseSubscriptionCanceled(log types.Log) (*FunctionsRegistrySubscriptionCanceled, error)

	FilterSubscriptionConsumerAdded(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsRegistrySubscriptionConsumerAddedIterator, error)

	WatchSubscriptionConsumerAdded(opts *bind.WatchOpts, sink chan<- *FunctionsRegistrySubscriptionConsumerAdded, subscriptionId []uint64) (event.Subscription, error)

	ParseSubscriptionConsumerAdded(log types.Log) (*FunctionsRegistrySubscriptionConsumerAdded, error)

	FilterSubscriptionConsumerRemoved(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsRegistrySubscriptionConsumerRemovedIterator, error)

	WatchSubscriptionConsumerRemoved(opts *bind.WatchOpts, sink chan<- *FunctionsRegistrySubscriptionConsumerRemoved, subscriptionId []uint64) (event.Subscription, error)

	ParseSubscriptionConsumerRemoved(log types.Log) (*FunctionsRegistrySubscriptionConsumerRemoved, error)

	FilterSubscriptionCreated(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsRegistrySubscriptionCreatedIterator, error)

	WatchSubscriptionCreated(opts *bind.WatchOpts, sink chan<- *FunctionsRegistrySubscriptionCreated, subscriptionId []uint64) (event.Subscription, error)

	ParseSubscriptionCreated(log types.Log) (*FunctionsRegistrySubscriptionCreated, error)

	FilterSubscriptionFunded(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsRegistrySubscriptionFundedIterator, error)

	WatchSubscriptionFunded(opts *bind.WatchOpts, sink chan<- *FunctionsRegistrySubscriptionFunded, subscriptionId []uint64) (event.Subscription, error)

	ParseSubscriptionFunded(log types.Log) (*FunctionsRegistrySubscriptionFunded, error)

	FilterSubscriptionOwnerTransferRequested(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsRegistrySubscriptionOwnerTransferRequestedIterator, error)

	WatchSubscriptionOwnerTransferRequested(opts *bind.WatchOpts, sink chan<- *FunctionsRegistrySubscriptionOwnerTransferRequested, subscriptionId []uint64) (event.Subscription, error)

	ParseSubscriptionOwnerTransferRequested(log types.Log) (*FunctionsRegistrySubscriptionOwnerTransferRequested, error)

	FilterSubscriptionOwnerTransferred(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsRegistrySubscriptionOwnerTransferredIterator, error)

	WatchSubscriptionOwnerTransferred(opts *bind.WatchOpts, sink chan<- *FunctionsRegistrySubscriptionOwnerTransferred, subscriptionId []uint64) (event.Subscription, error)

	ParseSubscriptionOwnerTransferred(log types.Log) (*FunctionsRegistrySubscriptionOwnerTransferred, error)

	FilterUnpaused(opts *bind.FilterOpts) (*FunctionsRegistryUnpausedIterator, error)

	WatchUnpaused(opts *bind.WatchOpts, sink chan<- *FunctionsRegistryUnpaused) (event.Subscription, error)

	ParseUnpaused(log types.Log) (*FunctionsRegistryUnpaused, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
