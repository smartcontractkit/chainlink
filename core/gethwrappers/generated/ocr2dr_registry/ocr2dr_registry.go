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
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated"
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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"linkEthFeed\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"internalBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"externalBalance\",\"type\":\"uint256\"}],\"name\":\"BalanceInvariantViolated\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptySendersList\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"have\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"want\",\"type\":\"uint32\"}],\"name\":\"GasLimitTooBig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectRequestID\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"InvalidConsumer\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"linkWei\",\"type\":\"int256\"}],\"name\":\"InvalidLinkWeiPrice\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSubscription\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposedOwner\",\"type\":\"address\"}],\"name\":\"MustBeRequestedOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"MustBeSubOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAllowedToSetSenders\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableFromLink\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentTooLarge\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingRequestExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Reentrant\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyConsumers\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedSender\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"senders\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"changedBy\",\"type\":\"address\"}],\"name\":\"AuthorizedSendersChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"signerPayment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"transmitterPayment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalCost\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"BillingEnd\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"don\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"donFee\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"registryFee\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"estimatedCost\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structFunctionsBillingRegistry.Commitment\",\"name\":\"commitment\",\"type\":\"tuple\"}],\"name\":\"BillingStart\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasOverhead\",\"type\":\"uint32\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"RequestTimedOut\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"SubscriptionCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"SubscriptionCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK_ETH_FEED\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_CONSUMERS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"acceptSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"addConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"cancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createSubscription\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"donFee\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"registryFee\",\"type\":\"uint96\"}],\"name\":\"estimateCost\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"response\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"internalType\":\"address[31]\",\"name\":\"signers\",\"type\":\"address[31]\"},{\"internalType\":\"uint8\",\"name\":\"signerCount\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"reportValidationGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"initialGas\",\"type\":\"uint256\"}],\"name\":\"fulfillAndBill\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAuthorizedSenders\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfig\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"internalType\":\"uint32\",\"name\":\"gasOverhead\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCurrentsubscriptionId\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRequestConfig\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"},{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"}],\"internalType\":\"structFunctionsBillingRegistryInterface.RequestBilling\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"getRequiredFee\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"getSubscription\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"getSubscriptionOwner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTotalBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"isAuthorizedSender\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"oracleWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"ownerCancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"pendingRequestExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"removeConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"requestSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"senders\",\"type\":\"address[]\"}],\"name\":\"setAuthorizedSenders\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"internalType\":\"uint32\",\"name\":\"gasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"requestTimeoutSeconds\",\"type\":\"uint32\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"}],\"internalType\":\"structFunctionsBillingRegistryInterface.RequestBilling\",\"name\":\"billing\",\"type\":\"tuple\"}],\"name\":\"startBilling\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"requestIdsToTimeout\",\"type\":\"bytes32[]\"}],\"name\":\"timeoutRequests\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60e06040523480156200001157600080fd5b5060405162005856380380620058568339810160408190526200003491620001bf565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000f6565b50506001805460ff60a01b19169055506001600160601b0319606093841b811660805291831b821660a05290911b1660c05262000209565b6001600160a01b038116331415620001515760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b0381168114620001ba57600080fd5b919050565b600080600060608486031215620001d557600080fd5b620001e084620001a2565b9250620001f060208501620001a2565b91506200020060408501620001a2565b90509250925092565b60805160601c60a05160601c60c05160601c6155d56200028160003960008181611a9e01528181611b6f015281816123f701526124c80152600081816104ce01526143c60152600081816102e0015281816112e3015281816128ed01528181613532015281816136880152613e7f01526155d56000f3fe608060405234801561001057600080fd5b50600436106102405760003560e01c80638456cb5911610145578063b2a489ff116100bd578063e82ad7d41161008c578063f1e14a2111610071578063f1e14a211461059a578063f2fde38b146105b1578063fa00763a146105c457600080fd5b8063e82ad7d414610574578063ee56997b1461058757600080fd5b8063b2a489ff146104f0578063c3f909d414610503578063d7ae1d301461054e578063e72f6e301461056157600080fd5b8063a21a23e411610114578063a4c0ed36116100f9578063a4c0ed36146104a3578063a9d03c05146104b6578063ad178361146104c957600080fd5b8063a21a23e414610479578063a47c76961461048157600080fd5b80638456cb59146104105780638da5cb5b146104185780639f87fad714610436578063a1a6d0411461044957600080fd5b806333652e3e116101d857806366316d8d116101a75780637341c10c1161018c5780637341c10c146103e257806379ba5097146103f557806382359740146103fd57600080fd5b806366316d8d146103bc578063665871ec146103cf57600080fd5b806333652e3e1461034f5780633f4ba83a146103765780635c975abb1461037e57806364d51a2a146103a157600080fd5b806312b583491161021457806312b58349146102af5780631b6b6d23146102db5780632408afaa1461032757806327923e411461033c57600080fd5b80620122911461024557806302bcc5b61461026457806304c357cb146102795780630739e4f11461028c575b600080fd5b61024d6105d7565b60405161025b9291906152c4565b60405180910390f35b610277610272366004614f3d565b6105f6565b005b610277610287366004614f58565b610673565b61029f61029a366004614c3a565b61086e565b604051901515815260200161025b565b6008546801000000000000000090046bffffffffffffffffffffffff165b60405190815260200161025b565b6103027f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161025b565b61032f610f82565b60405161025b919061510a565b61027761034a366004614ed7565b610ff1565b60085467ffffffffffffffff165b60405167ffffffffffffffff909116815260200161025b565b610277611144565b60015474010000000000000000000000000000000000000000900460ff1661029f565b6103a9606481565b60405161ffff909116815260200161025b565b6102776103ca366004614b9f565b611156565b6102776103dd366004614bd6565b6113f3565b6102776103f0366004614f58565b6116bf565b610277611952565b61027761040b366004614f3d565b611a4f565b610277611dea565b60005473ffffffffffffffffffffffffffffffffffffffff16610302565b610277610444366004614f58565b611dfa565b61045c610457366004614e86565b612281565b6040516bffffffffffffffffffffffff909116815260200161025b565b61035d6123a5565b61049461048f366004614f3d565b612757565b60405161025b939291906152eb565b6102776104b1366004614b45565b612888565b6102cd6104c4366004614d64565b612afe565b6103027f000000000000000000000000000000000000000000000000000000000000000081565b6103026104fe366004614f3d565b6132f9565b600c54600d54600b54600e546040805163ffffffff8087168252650100000000009096048616602082015290810193909352606083019190915291909116608082015260a00161025b565b61027761055c366004614f58565b613392565b61027761056f366004614b2a565b6134f9565b61029f610582366004614f3d565b61375c565b610277610595366004614bd6565b61399b565b61045c6105a8366004614de1565b60009392505050565b6102776105bf366004614b2a565b613b0e565b61029f6105d2366004614b2a565b613b22565b600c5460009060609063ffffffff166105ee610f82565b915091509091565b6105fe613b35565b67ffffffffffffffff811660009081526006602052604090205473ffffffffffffffffffffffffffffffffffffffff1680610665576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61066f8282613bb6565b5050565b67ffffffffffffffff8216600090815260066020526040902054829073ffffffffffffffffffffffffffffffffffffffff16806106dc576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff821614610748576040517fd8a3fb5200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821660048201526024015b60405180910390fd5b600c54640100000000900460ff161561078d576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610795613ff4565b67ffffffffffffffff841660009081526006602052604090206001015473ffffffffffffffffffffffffffffffffffffffff8481169116146108685767ffffffffffffffff841660008181526006602090815260409182902060010180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff88169081179091558251338152918201527f69436ea6df009049404f564eff6622cd00522b0bd6a89efd9e52a355c4a879be91015b60405180910390a25b50505050565b6000610878614079565b600c54640100000000900460ff16156108bd576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6108c5613ff4565b60008b8152600a6020908152604091829020825161012081018452815467ffffffffffffffff8116825268010000000000000000810473ffffffffffffffffffffffffffffffffffffffff908116948301949094527c0100000000000000000000000000000000000000000000000000000000900463ffffffff169381019390935260018101546060840152600281015491821660808401819052740100000000000000000000000000000000000000009092046bffffffffffffffffffffffff90811660a0850152600382015480821660c08601526c0100000000000000000000000090041660e0840152600401546101008301526109f1576040517fda7aa3e100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008c8152600a602052604080822082815560018101839055600281018390556003810180547fffffffffffffffff000000000000000000000000000000000000000000000000169055600401829055517f0ca761750000000000000000000000000000000000000000000000000000000090610a7a908f908f908f908f908f9060240161511d565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152918152602080830180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090951694909417909352600c80547fffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffff1664010000000017905584015191840151909250610b429163ffffffff1690836140b8565b600c80547fffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffff169055600d5460a084015160c0850151929550600092610b8e92889290918b908b3a614104565b604080820151855167ffffffffffffffff166000908152600760205291909120549192506bffffffffffffffffffffffff90811691161015610bfc576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b604080820151845167ffffffffffffffff1660009081526007602052918220805491929091610c3a9084906bffffffffffffffffffffffff16615435565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060005b8760ff16811015610d7d578973ffffffffffffffffffffffffffffffffffffffff168982601f8110610c9f57610c9f615550565b602002015173ffffffffffffffffffffffffffffffffffffffff1614610d6b578151600960008b84601f8110610cd757610cd7615550565b602002015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282829054906101000a90046bffffffffffffffffffffffff16610d3c919061537b565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055505b80610d7581615462565b915050610c6b565b508260c0015160096000610da660005473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff168152602081019190915260400160009081208054909190610dec9084906bffffffffffffffffffffffff1661537b565b82546101009290920a6bffffffffffffffffffffffff81810219909316918316021790915560208381015173ffffffffffffffffffffffffffffffffffffffff8d166000908152600990925260408220805491945092610e4e9185911661537b565b82546101009290920a6bffffffffffffffffffffffff81810219909316918316021790915560e0850151855167ffffffffffffffff166000908152600760205260409020805491935091600c91610eb79185916c01000000000000000000000000900416615435565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055508d7fc8dc973332de19a5f71b6026983110e9c2e04b0c98b87eb771ccb78607fd114f846000015183600001518460200151856040015189604051610f6995949392919067ffffffffffffffff9590951685526bffffffffffffffffffffffff9384166020860152918316604085015290911660608301521515608082015260a00190565b60405180910390a25050509a9950505050505050505050565b60606004805480602002602001604051908101604052809291908181526020018280548015610fe757602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610fbc575b5050505050905090565b610ff9613b35565b60008313611036576040517f43d4cf660000000000000000000000000000000000000000000000000000000081526004810184905260240161073f565b6040805160c08101825263ffffffff888116808352600060208085019190915289831684860181905260608086018b9052888516608080880182905295891660a0978801819052600c80547fffffffffffffffffffffffffffffffffffffffffffffff000000000000000000168717650100000000008602179055600d8d9055600e80547fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000168317640100000000909202919091179055600b8b9055875194855292840191909152948201899052938101879052908101929092527f24d3d934adfef9b9029d6ffa463c07d0139ed47d26ee23506f85ece2879d2bd4910160405180910390a1505050505050565b61114c613b35565b611154614291565b565b600c54640100000000900460ff161561119b576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6111a3613ff4565b6bffffffffffffffffffffffff81166111d65750336000908152600960205260409020546bffffffffffffffffffffffff165b336000908152600960205260409020546bffffffffffffffffffffffff80831691161015611230576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b336000908152600960205260408120805483929061125d9084906bffffffffffffffffffffffff16615435565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550806008808282829054906101000a90046bffffffffffffffffffffffff166112b39190615435565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb83836040518363ffffffff1660e01b815260040161136b92919073ffffffffffffffffffffffffffffffffffffffff9290921682526bffffffffffffffffffffffff16602082015260400190565b602060405180830381600087803b15801561138557600080fd5b505af1158015611399573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906113bd9190614c18565b61066f576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6113fb613ff4565b60005b818110156116ba57600083838381811061141a5761141a615550565b602090810292909201356000818152600a84526040808220815161012081018352815467ffffffffffffffff811680835268010000000000000000820473ffffffffffffffffffffffffffffffffffffffff908116848b01527c010000000000000000000000000000000000000000000000000000000090920463ffffffff168386015260018401546060840152600284015480831660808501527401000000000000000000000000000000000000000090046bffffffffffffffffffffffff90811660a0850152600385015480821660c08601526c0100000000000000000000000090041660e08401526004909301546101008301529184526006909652912054919450163314905061158f57805167ffffffffffffffff16600090815260066020526040908190205490517fd8a3fb5200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff909116600482015260240161073f565b600e5461010082015142916115b39164010000000090910463ffffffff1690615337565b11156116a55760e0810151815167ffffffffffffffff1660009081526007602052604090208054600c906116069084906c0100000000000000000000000090046bffffffffffffffffffffffff16615435565b82546bffffffffffffffffffffffff9182166101009390930a9283029190920219909116179055506000828152600a602052604080822082815560018101839055600281018390556003810180547fffffffffffffffff0000000000000000000000000000000000000000000000001690556004018290555183917ff1ca1e9147be737b04a2b018a79405f687a97de8dd8a2559bbe62357343af41491a25b505080806116b290615462565b9150506113fe565b505050565b67ffffffffffffffff8216600090815260066020526040902054829073ffffffffffffffffffffffffffffffffffffffff1680611728576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff82161461178f576040517fd8a3fb5200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8216600482015260240161073f565b600c54640100000000900460ff16156117d4576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6117dc613ff4565b67ffffffffffffffff841660009081526006602052604090206002015460641415611833576040517f05a48e0f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8316600090815260056020908152604080832067ffffffffffffffff8089168552925290912054161561187a57610868565b73ffffffffffffffffffffffffffffffffffffffff8316600081815260056020908152604080832067ffffffffffffffff891680855290835281842080547fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000166001908117909155600684528285206002018054918201815585529383902090930180547fffffffffffffffffffffffff000000000000000000000000000000000000000016851790555192835290917f43dc749a04ac8fb825cbd514f7c0e13f13bc6f2ee66043b76629d51776cff8e0910161085f565b60015473ffffffffffffffffffffffffffffffffffffffff1633146119d3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015260640161073f565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b600c54640100000000900460ff1615611a94576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611a9c613ff4565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16634b4fa0c16040518163ffffffff1660e01b815260040160206040518083038186803b158015611b0257600080fd5b505afa158015611b16573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611b3a9190614c18565b8015611c0057506040517ffa00763a0000000000000000000000000000000000000000000000000000000081523360048201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff169063fa00763a9060240160206040518083038186803b158015611bc657600080fd5b505afa158015611bda573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611bfe9190614c18565b155b15611c37576040517f0809490800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff811660009081526006602052604090205473ffffffffffffffffffffffffffffffffffffffff16611c9d576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff811660009081526006602052604090206001015473ffffffffffffffffffffffffffffffffffffffff163314611d3f5767ffffffffffffffff8116600090815260066020526040908190206001015490517fd084e97500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff909116600482015260240161073f565b67ffffffffffffffff81166000818152600660209081526040918290208054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560019093018054909316909255835173ffffffffffffffffffffffffffffffffffffffff909116808252928101919091529092917f6f1dc65165ffffedfd8e507b4a0f1fcfdada045ed11f6c26ba27cedfe87802f0910160405180910390a25050565b611df2613b35565b61115461430e565b67ffffffffffffffff8216600090815260066020526040902054829073ffffffffffffffffffffffffffffffffffffffff1680611e63576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff821614611eca576040517fd8a3fb5200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8216600482015260240161073f565b600c54640100000000900460ff1615611f0f576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611f17613ff4565b73ffffffffffffffffffffffffffffffffffffffff8316600090815260056020908152604080832067ffffffffffffffff808916855292529091205416611fb2576040517ff0019fe600000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8516600482015273ffffffffffffffffffffffffffffffffffffffff8416602482015260440161073f565b67ffffffffffffffff841660009081526006602090815260408083206002018054825181850281018501909352808352919290919083018282801561202d57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311612002575b50505050509050600060018251612044919061541e565b905060005b82518110156121e3578573ffffffffffffffffffffffffffffffffffffffff1683828151811061207b5761207b615550565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1614156121d15760008383815181106120b3576120b3615550565b6020026020010151905080600660008a67ffffffffffffffff1667ffffffffffffffff16815260200190815260200160002060020183815481106120f9576120f9615550565b600091825260208083209190910180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff949094169390931790925567ffffffffffffffff8a16815260069091526040902060020180548061217357612173615521565b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff0000000000000000000000000000000000000000169055019055506121e3565b806121db81615462565b915050612049565b5073ffffffffffffffffffffffffffffffffffffffff8516600081815260056020908152604080832067ffffffffffffffff8b168085529083529281902080547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001690555192835290917f182bff9831466789164ca77075fffd84916d35a8180ba73c27e45634549b445b91015b60405180910390a2505050505050565b60008061228c61437d565b9050600081136122cb576040517f43d4cf660000000000000000000000000000000000000000000000000000000081526004810182905260240161073f565b600d54600e5460009163ffffffff808a16926122e8929116615337565b6122f29190615337565b90506000828261230a89670de0b6b3a76400006153e1565b61231491906153e1565b61231e91906153a2565b9050600061233d6bffffffffffffffffffffffff808816908916615337565b9050612355816b033b2e3c9fd0803ce800000061541e565b82111561238e576040517fe80fa38100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6123988183615337565b9998505050505050505050565b600c54600090640100000000900460ff16156123ed576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6123f5613ff4565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16634b4fa0c16040518163ffffffff1660e01b815260040160206040518083038186803b15801561245b57600080fd5b505afa15801561246f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906124939190614c18565b801561255957506040517ffa00763a0000000000000000000000000000000000000000000000000000000081523360048201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff169063fa00763a9060240160206040518083038186803b15801561251f57600080fd5b505afa158015612533573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906125579190614c18565b155b15612590576040517f0809490800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6008805467ffffffffffffffff169060006125aa8361549b565b82546101009290920a67ffffffffffffffff8181021990931691831602179091556008541690506000806040519080825280602002602001820160405280156125fd578160200160208202803683370190505b506040805180820182526000808252602080830182815267ffffffffffffffff8816808452600783528584209451855492516bffffffffffffffffffffffff9081166c01000000000000000000000000027fffffffffffffffff00000000000000000000000000000000000000000000000090941691161791909117909355835160608101855233815280820183815281860187815294845260068352949092208251815473ffffffffffffffffffffffffffffffffffffffff9182167fffffffffffffffffffffffff000000000000000000000000000000000000000091821617835595516001830180549190921696169590951790945591518051949550909361270f9260028501920190614876565b505060405133815267ffffffffffffffff841691507f464722b4166576d3dcbba877b999bc35cf911f4eaf434b7eba68fa113951d0bf9060200160405180910390a250905090565b67ffffffffffffffff8116600090815260066020526040812054819060609073ffffffffffffffffffffffffffffffffffffffff166127c2576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff841660009081526007602090815260408083205460068352928190208054600290910180548351818602810186019094528084526bffffffffffffffffffffffff9095169473ffffffffffffffffffffffffffffffffffffffff90921693909291839183018282801561287457602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311612849575b505050505090509250925092509193909250565b600c54640100000000900460ff16156128cd576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6128d5613ff4565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614612944576040517f44b0e3c300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6020811461297e576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600061298c82840184614f3d565b67ffffffffffffffff811660009081526006602052604090205490915073ffffffffffffffffffffffffffffffffffffffff166129f5576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff8116600090815260076020526040812080546bffffffffffffffffffffffff1691869190612a2c838561537b565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550846008808282829054906101000a90046bffffffffffffffffffffffff16612a82919061537b565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055508167ffffffffffffffff167fd39ec07f4e209f627a4c427971473820dc129761ba28de8906bd56f57101d4f8828784612ae99190615337565b60408051928352602083019190915201612271565b6000612b08614079565b600c54640100000000900460ff1615612b4d576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b612b55613ff4565b6000600681612b676020860186614f3d565b67ffffffffffffffff16815260208101919091526040016000205473ffffffffffffffffffffffffffffffffffffffff161415612bd0576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000600581612be56040860160208701614b2a565b73ffffffffffffffffffffffffffffffffffffffff168152602080820192909252604001600090812091612c1b90860186614f3d565b67ffffffffffffffff908116825260208201929092526040016000205416905080612cb757612c4d6020840184614f3d565b612c5d6040850160208601614b2a565b6040517ff0019fe600000000000000000000000000000000000000000000000000000000815267ffffffffffffffff909216600483015273ffffffffffffffffffffffffffffffffffffffff16602482015260440161073f565b600c5463ffffffff16612cd06060850160408601614e6b565b63ffffffff161115612d3157612cec6060840160408501614e6b565b600c546040517ff5d7e01e00000000000000000000000000000000000000000000000000000000815263ffffffff92831660048201529116602482015260440161073f565b6040517ff1e14a21000000000000000000000000000000000000000000000000000000008152600090339063f1e14a2190612d7490899089908990600401615156565b60206040518083038186803b158015612d8c57600080fd5b505afa158015612da0573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612dc49190614fdb565b90506000612ddc87876105a836899003890189614e36565b90506000612dff612df36060880160408901614e6b565b87606001358585612281565b90506000600781612e1360208a018a614f3d565b67ffffffffffffffff1681526020808201929092526040016000908120546c0100000000000000000000000090046bffffffffffffffffffffffff169160079190612e60908b018b614f3d565b67ffffffffffffffff168152602081019190915260400160002054612e9391906bffffffffffffffffffffffff16615435565b9050816bffffffffffffffffffffffff16816bffffffffffffffffffffffff161015612eeb576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000612ef886600161534f565b90506000612f8033612f1060408c0160208d01614b2a565b612f1d60208d018d614f3d565b856040805173ffffffffffffffffffffffffffffffffffffffff958616602080830191909152949095168582015267ffffffffffffffff928316606086015291166080808501919091528151808503909101815260a09093019052815191012090565b60408051610120810190915290915060009080612fa060208d018d614f3d565b67ffffffffffffffff1681526020018b6020016020810190612fc29190614b2a565b73ffffffffffffffffffffffffffffffffffffffff168152602001612fed60608d0160408e01614e6b565b63ffffffff90811682526060808e0135602080850191909152336040808601919091526bffffffffffffffffffffffff808e16848701528c81166080808801919091528c821660a0808901919091524260c09889015260008b8152600a86528481208a5181548c890151978d0151909a167c0100000000000000000000000000000000000000000000000000000000027bffffffffffffffffffffffffffffffffffffffffffffffffffffffff73ffffffffffffffffffffffffffffffffffffffff98891668010000000000000000027fffffffff00000000000000000000000000000000000000000000000000000000909c1667ffffffffffffffff909316929092179a909a171698909817885595890151600188015590880151908801518216740100000000000000000000000000000000000000000292169190911760028501559385015160038401805460e088015187166c01000000000000000000000000027fffffffffffffffff00000000000000000000000000000000000000000000000090911692909616919091179490941790935561010084015160049092019190915591925086916007916131a7908e018e614f3d565b67ffffffffffffffff16815260208101919091526040016000208054600c906131ef9084906c0100000000000000000000000090046bffffffffffffffffffffffff1661537b565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550817f99f7f4e65b4b9fbabd4e357c47ed3099b36e57ecd3a43e84662f34c207d0ebe48260405161324d91906151d4565b60405180910390a2826005600061326a60408e0160208f01614b2a565b73ffffffffffffffffffffffffffffffffffffffff1681526020808201929092526040016000908120916132a0908e018e614f3d565b67ffffffffffffffff9081168252602082019290925260400160002080547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001692909116919091179055509a9950505050505050505050565b67ffffffffffffffff811660009081526006602052604081205473ffffffffffffffffffffffffffffffffffffffff1661335f576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5067ffffffffffffffff1660009081526006602052604090205473ffffffffffffffffffffffffffffffffffffffff1690565b67ffffffffffffffff8216600090815260066020526040902054829073ffffffffffffffffffffffffffffffffffffffff16806133fb576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff821614613462576040517fd8a3fb5200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8216600482015260240161073f565b600c54640100000000900460ff16156134a7576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6134af613ff4565b6134b88461375c565b156134ef576040517fb42f66e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6108688484613bb6565b613501613b35565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526000907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906370a082319060240160206040518083038186803b15801561358957600080fd5b505afa15801561359d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906135c19190614e52565b6008549091506801000000000000000090046bffffffffffffffffffffffff1681811115613625576040517fa99da302000000000000000000000000000000000000000000000000000000008152600481018290526024810183905260440161073f565b818110156116ba576000613639828461541e565b6040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8681166004830152602482018390529192507f00000000000000000000000000000000000000000000000000000000000000009091169063a9059cbb90604401602060405180830381600087803b1580156136ce57600080fd5b505af11580156136e2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906137069190614c18565b506040805173ffffffffffffffffffffffffffffffffffffffff86168152602081018390527f59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b436600910160405180910390a150505050565b67ffffffffffffffff81166000908152600660209081526040808320600201805482518185028101850190935280835284938301828280156137d457602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116137a9575b5050505050905060006137e5610f82565b905060005b82518110156139905760005b825181101561397d57600061392d84838151811061381657613816615550565b602002602001015186858151811061383057613830615550565b602002602001015189600560008a898151811061384f5761384f615550565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008c67ffffffffffffffff1667ffffffffffffffff16815260200190815260200160002060009054906101000a900467ffffffffffffffff166040805173ffffffffffffffffffffffffffffffffffffffff958616602080830191909152949095168582015267ffffffffffffffff928316606086015291166080808501919091528151808503909101815260a09093019052815191012090565b6000818152600a602052604090206002015490915073ffffffffffffffffffffffffffffffffffffffff161561396a575060019695505050505050565b508061397581615462565b9150506137f6565b508061398881615462565b9150506137ea565b506000949350505050565b6139a361448e565b6139d9576040517fad77f06100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80613a10576040517f75158c3b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b600454811015613a7057613a5d60048281548110613a3357613a33615550565b60009182526020909120015460029073ffffffffffffffffffffffffffffffffffffffff1661449e565b5080613a6881615462565b915050613a13565b5060005b81811015613ac157613aae838383818110613a9157613a91615550565b9050602002016020810190613aa69190614b2a565b6002906144c7565b5080613ab981615462565b915050613a74565b50613ace600483836148fc565b507ff263cfb3e4298332e776194610cf9fdc09ccb3ada8b9aa39764d882e11fbf0a0828233604051613b0293929190615092565b60405180910390a15050565b613b16613b35565b613b1f816144e9565b50565b6000613b2f6002836145df565b92915050565b60005473ffffffffffffffffffffffffffffffffffffffff163314611154576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161073f565b600c54640100000000900460ff1615613bfb576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff821660009081526006602090815260408083208151606081018352815473ffffffffffffffffffffffffffffffffffffffff908116825260018301541681850152600282018054845181870281018701865281815292959394860193830182828015613ca657602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311613c7b575b5050509190925250505067ffffffffffffffff84166000908152600760205260408120549192506bffffffffffffffffffffffff909116905b826040015151811015613d85576005600084604001518381518110613d0657613d06615550565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040908101600090812067ffffffffffffffff89168252909252902080547fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000016905580613d7d81615462565b915050613cdf565b5067ffffffffffffffff8416600090815260066020526040812080547fffffffffffffffffffffffff00000000000000000000000000000000000000009081168255600182018054909116905590613de06002830182614974565b505067ffffffffffffffff8416600090815260076020526040902080547fffffffffffffffff000000000000000000000000000000000000000000000000169055600880548291908190613e4f9084906801000000000000000090046bffffffffffffffffffffffff16615435565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb84836bffffffffffffffffffffffff166040518363ffffffff1660e01b8152600401613f0792919073ffffffffffffffffffffffffffffffffffffffff929092168252602082015260400190565b602060405180830381600087803b158015613f2157600080fd5b505af1158015613f35573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613f599190614c18565b613f8f576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805173ffffffffffffffffffffffffffffffffffffffff851681526bffffffffffffffffffffffff8316602082015267ffffffffffffffff8616917fe8ed5b475a5b5987aa9165e8731bb78043f39eee32ec5a1169a89e27fcd49815910161085f565b60015474010000000000000000000000000000000000000000900460ff1615611154576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f5061757361626c653a2070617573656400000000000000000000000000000000604482015260640161073f565b61408233613b22565b611154576040517f0809490800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005a6113888110156140ca57600080fd5b6113888103905084604082048203116140e257600080fd5b50823b6140ee57600080fd5b60008083516020850160008789f1949350505050565b604080516060810182526000808252602082018190529181018290529061412961437d565b905060008113614168576040517f43d4cf660000000000000000000000000000000000000000000000000000000081526004810182905260240161073f565b6000815a8b6141778c89615337565b6141819190615337565b61418b919061541e565b61419d86670de0b6b3a76400006153e1565b6141a791906153e1565b6141b191906153a2565b905060006141d06bffffffffffffffffffffffff808916908b16615337565b90506141e8816b033b2e3c9fd0803ce800000061541e565b821115614221576040517fe80fa38100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600061423060ff8a168b6153b6565b9050600061423e828561537b565b9050600061425461424f8587615337565b61460e565b604080516060810182526bffffffffffffffffffffffff958616815293851660208501529316928201929092529c9b505050505050505050505050565b6142996146b0565b600180547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff1690557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa335b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390a1565b614316613ff4565b600180547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff16740100000000000000000000000000000000000000001790557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a2586142e43390565b600c54604080517ffeaf968c000000000000000000000000000000000000000000000000000000008152905160009265010000000000900463ffffffff169182151591849182917f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff169163feaf968c9160048083019260a0929190829003018186803b15801561442157600080fd5b505afa158015614435573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906144599190614f8b565b5093505092505082801561447b5750614472814261541e565b8463ffffffff16105b1561448657600b5491505b509392505050565b6000614498613b35565b50600190565b60006144c08373ffffffffffffffffffffffffffffffffffffffff8416614734565b9392505050565b60006144c08373ffffffffffffffffffffffffffffffffffffffff8416614827565b73ffffffffffffffffffffffffffffffffffffffff8116331415614569576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161073f565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b73ffffffffffffffffffffffffffffffffffffffff8116600090815260018301602052604081205415156144c0565b60006bffffffffffffffffffffffff8211156146ac576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f53616665436173743a2076616c756520646f65736e27742066697420696e203960448201527f3620626974730000000000000000000000000000000000000000000000000000606482015260840161073f565b5090565b60015474010000000000000000000000000000000000000000900460ff16611154576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f5061757361626c653a206e6f7420706175736564000000000000000000000000604482015260640161073f565b6000818152600183016020526040812054801561481d57600061475860018361541e565b855490915060009061476c9060019061541e565b90508181146147d157600086600001828154811061478c5761478c615550565b90600052602060002001549050808760000184815481106147af576147af615550565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806147e2576147e2615521565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050613b2f565b6000915050613b2f565b600081815260018301602052604081205461486e57508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155613b2f565b506000613b2f565b8280548282559060005260206000209081019282156148f0579160200282015b828111156148f057825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190614896565b506146ac92915061498e565b8280548282559060005260206000209081019282156148f0579160200282015b828111156148f05781547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff84351617825560209092019160019091019061491c565b5080546000825590600052602060002090810190613b1f91905b5b808211156146ac576000815560010161498f565b803573ffffffffffffffffffffffffffffffffffffffff811681146149c757600080fd5b919050565b60008083601f8401126149de57600080fd5b50813567ffffffffffffffff8111156149f657600080fd5b6020830191508360208260051b8501011115614a1157600080fd5b9250929050565b60008083601f840112614a2a57600080fd5b50813567ffffffffffffffff811115614a4257600080fd5b602083019150836020828501011115614a1157600080fd5b600060808284031215614a6c57600080fd5b6040516080810181811067ffffffffffffffff82111715614a8f57614a8f61557f565b604052905080614a9e83614ae7565b8152614aac602084016149a3565b6020820152614abd60408401614ad3565b6040820152606083013560608201525092915050565b803563ffffffff811681146149c757600080fd5b803567ffffffffffffffff811681146149c757600080fd5b803560ff811681146149c757600080fd5b805169ffffffffffffffffffff811681146149c757600080fd5b600060208284031215614b3c57600080fd5b6144c0826149a3565b60008060008060608587031215614b5b57600080fd5b614b64856149a3565b935060208501359250604085013567ffffffffffffffff811115614b8757600080fd5b614b9387828801614a18565b95989497509550505050565b60008060408385031215614bb257600080fd5b614bbb836149a3565b91506020830135614bcb816155ae565b809150509250929050565b60008060208385031215614be957600080fd5b823567ffffffffffffffff811115614c0057600080fd5b614c0c858286016149cc565b90969095509350505050565b600060208284031215614c2a57600080fd5b815180151581146144c057600080fd5b6000806000806000806000806000806104c08b8d031215614c5a57600080fd5b8a35995060208b013567ffffffffffffffff80821115614c7957600080fd5b614c858e838f01614a18565b909b50995060408d0135915080821115614c9e57600080fd5b614caa8e838f01614a18565b9099509750879150614cbe60608e016149a3565b96508d609f8e0112614ccf57600080fd5b60405191506103e082018281108282111715614ced57614ced61557f565b604052508060808d016104608e018f811115614d0857600080fd5b60005b601f811015614d3257614d1d836149a3565b84526020938401939290920191600101614d0b565b50839750614d3f81614aff565b9650505050506104808b013591506104a08b013590509295989b9194979a5092959850565b600080600083850360a0811215614d7a57600080fd5b843567ffffffffffffffff811115614d9157600080fd5b614d9d87828801614a18565b90955093505060807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe082011215614dd357600080fd5b506020840190509250925092565b600080600060a08486031215614df657600080fd5b833567ffffffffffffffff811115614e0d57600080fd5b614e1986828701614a18565b9094509250614e2d90508560208601614a5a565b90509250925092565b600060808284031215614e4857600080fd5b6144c08383614a5a565b600060208284031215614e6457600080fd5b5051919050565b600060208284031215614e7d57600080fd5b6144c082614ad3565b60008060008060808587031215614e9c57600080fd5b614ea585614ad3565b9350602085013592506040850135614ebc816155ae565b91506060850135614ecc816155ae565b939692955090935050565b60008060008060008060c08789031215614ef057600080fd5b614ef987614ad3565b9550614f0760208801614ad3565b94506040870135935060608701359250614f2360808801614ad3565b9150614f3160a08801614ad3565b90509295509295509295565b600060208284031215614f4f57600080fd5b6144c082614ae7565b60008060408385031215614f6b57600080fd5b614f7483614ae7565b9150614f82602084016149a3565b90509250929050565b600080600080600060a08688031215614fa357600080fd5b614fac86614b10565b9450602086015193506040860151925060608601519150614fcf60808701614b10565b90509295509295909350565b600060208284031215614fed57600080fd5b81516144c0816155ae565b600081518084526020808501945080840160005b8381101561503e57815173ffffffffffffffffffffffffffffffffffffffff168752958201959082019060010161500c565b509495945050505050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b6040808252810183905260008460608301825b868110156150e05773ffffffffffffffffffffffffffffffffffffffff6150cb846149a3565b168252602092830192909101906001016150a5565b50809250505073ffffffffffffffffffffffffffffffffffffffff83166020830152949350505050565b6020815260006144c06020830184614ff8565b858152606060208201526000615137606083018688615049565b828103604084015261514a818587615049565b98975050505050505050565b60a08152600061516a60a083018587615049565b905067ffffffffffffffff61517e84614ae7565b16602083015273ffffffffffffffffffffffffffffffffffffffff6151a5602085016149a3565b16604083015263ffffffff6151bc60408501614ad3565b16606083015260608301356080830152949350505050565b60006101208201905067ffffffffffffffff835116825273ffffffffffffffffffffffffffffffffffffffff60208401511660208301526040830151615222604084018263ffffffff169052565b50606083015160608301526080830151615254608084018273ffffffffffffffffffffffffffffffffffffffff169052565b5060a083015161527460a08401826bffffffffffffffffffffffff169052565b5060c083015161529460c08401826bffffffffffffffffffffffff169052565b5060e08301516152b460e08401826bffffffffffffffffffffffff169052565b5061010092830151919092015290565b63ffffffff831681526040602082015260006152e36040830184614ff8565b949350505050565b6bffffffffffffffffffffffff8416815273ffffffffffffffffffffffffffffffffffffffff8316602082015260606040820152600061532e6060830184614ff8565b95945050505050565b6000821982111561534a5761534a6154c3565b500190565b600067ffffffffffffffff808316818516808303821115615372576153726154c3565b01949350505050565b60006bffffffffffffffffffffffff808316818516808303821115615372576153726154c3565b6000826153b1576153b16154f2565b500490565b60006bffffffffffffffffffffffff808416806153d5576153d56154f2565b92169190910492915050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615615419576154196154c3565b500290565b600082821015615430576154306154c3565b500390565b60006bffffffffffffffffffffffff8381169083168181101561545a5761545a6154c3565b039392505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415615494576154946154c3565b5060010190565b600067ffffffffffffffff808316818114156154b9576154b96154c3565b6001019392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6bffffffffffffffffffffffff81168114613b1f57600080fdfea164736f6c6343000806000a",
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
	parsed, err := abi.JSON(strings.NewReader(OCR2DRRegistryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
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

func (_OCR2DRRegistry *OCR2DRRegistryCaller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OCR2DRRegistry.contract.Call(opts, &out, "LINK")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_OCR2DRRegistry *OCR2DRRegistrySession) LINK() (common.Address, error) {
	return _OCR2DRRegistry.Contract.LINK(&_OCR2DRRegistry.CallOpts)
}

func (_OCR2DRRegistry *OCR2DRRegistryCallerSession) LINK() (common.Address, error) {
	return _OCR2DRRegistry.Contract.LINK(&_OCR2DRRegistry.CallOpts)
}

func (_OCR2DRRegistry *OCR2DRRegistryCaller) LINKETHFEED(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OCR2DRRegistry.contract.Call(opts, &out, "LINK_ETH_FEED")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_OCR2DRRegistry *OCR2DRRegistrySession) LINKETHFEED() (common.Address, error) {
	return _OCR2DRRegistry.Contract.LINKETHFEED(&_OCR2DRRegistry.CallOpts)
}

func (_OCR2DRRegistry *OCR2DRRegistryCallerSession) LINKETHFEED() (common.Address, error) {
	return _OCR2DRRegistry.Contract.LINKETHFEED(&_OCR2DRRegistry.CallOpts)
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

func (_OCR2DRRegistry *OCR2DRRegistryTransactor) SetConfig(opts *bind.TransactOpts, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation *big.Int, fallbackWeiPerUnitLink *big.Int, gasOverhead uint32, requestTimeoutSeconds uint32) (*types.Transaction, error) {
	return _OCR2DRRegistry.contract.Transact(opts, "setConfig", maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, gasOverhead, requestTimeoutSeconds)
}

func (_OCR2DRRegistry *OCR2DRRegistrySession) SetConfig(maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation *big.Int, fallbackWeiPerUnitLink *big.Int, gasOverhead uint32, requestTimeoutSeconds uint32) (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.SetConfig(&_OCR2DRRegistry.TransactOpts, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, gasOverhead, requestTimeoutSeconds)
}

func (_OCR2DRRegistry *OCR2DRRegistryTransactorSession) SetConfig(maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation *big.Int, fallbackWeiPerUnitLink *big.Int, gasOverhead uint32, requestTimeoutSeconds uint32) (*types.Transaction, error) {
	return _OCR2DRRegistry.Contract.SetConfig(&_OCR2DRRegistry.TransactOpts, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, gasOverhead, requestTimeoutSeconds)
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
	return common.HexToHash("0x24d3d934adfef9b9029d6ffa463c07d0139ed47d26ee23506f85ece2879d2bd4")
}

func (OCR2DRRegistryFundsRecovered) Topic() common.Hash {
	return common.HexToHash("0x59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b436600")
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
	LINK(opts *bind.CallOpts) (common.Address, error)

	LINKETHFEED(opts *bind.CallOpts) (common.Address, error)

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

	OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error)

	OracleWithdraw(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error)

	OwnerCancelSubscription(opts *bind.TransactOpts, subscriptionId uint64) (*types.Transaction, error)

	Pause(opts *bind.TransactOpts) (*types.Transaction, error)

	RecoverFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	RemoveConsumer(opts *bind.TransactOpts, subscriptionId uint64, consumer common.Address) (*types.Transaction, error)

	RequestSubscriptionOwnerTransfer(opts *bind.TransactOpts, subscriptionId uint64, newOwner common.Address) (*types.Transaction, error)

	SetAuthorizedSenders(opts *bind.TransactOpts, senders []common.Address) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation *big.Int, fallbackWeiPerUnitLink *big.Int, gasOverhead uint32, requestTimeoutSeconds uint32) (*types.Transaction, error)

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
