// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package functions_v1_events_mock

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

type FunctionsV1EventsMockConfig struct {
	MaxConsumersPerSubscription     uint16
	AdminFee                        *big.Int
	HandleOracleFulfillmentSelector [4]byte
	GasForCallExactCheck            uint16
	MaxCallbackGasLimits            []uint32
}

var FunctionsV1EventsMockMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint16\",\"name\":\"maxConsumersPerSubscription\",\"type\":\"uint16\"},{\"internalType\":\"uint72\",\"name\":\"adminFee\",\"type\":\"uint72\"},{\"internalType\":\"bytes4\",\"name\":\"handleOracleFulfillmentSelector\",\"type\":\"bytes4\"},{\"internalType\":\"uint16\",\"name\":\"gasForCallExactCheck\",\"type\":\"uint16\"},{\"internalType\":\"uint32[]\",\"name\":\"maxCallbackGasLimits\",\"type\":\"uint32[]\"}],\"indexed\":false,\"internalType\":\"structFunctionsV1EventsMock.Config\",\"name\":\"param1\",\"type\":\"tuple\"}],\"name\":\"ConfigUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"proposedContractSetId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"proposedContractSetFromAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"proposedContractSetToAddress\",\"type\":\"address\"}],\"name\":\"ContractProposed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"ContractUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"resultCode\",\"type\":\"uint8\"}],\"name\":\"RequestNotProcessed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalCostJuels\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"resultCode\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"response\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"callbackReturnData\",\"type\":\"bytes\"}],\"name\":\"RequestProcessed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"donId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"subscriptionOwner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requestingContract\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requestInitiator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"dataVersion\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"estimatedTotalCostJuels\",\"type\":\"uint96\"}],\"name\":\"RequestStart\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"RequestTimedOut\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"fundsRecipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fundsAmount\",\"type\":\"uint256\"}],\"name\":\"SubscriptionCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"SubscriptionCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint16\",\"name\":\"maxConsumersPerSubscription\",\"type\":\"uint16\"},{\"internalType\":\"uint72\",\"name\":\"adminFee\",\"type\":\"uint72\"},{\"internalType\":\"bytes4\",\"name\":\"handleOracleFulfillmentSelector\",\"type\":\"bytes4\"},{\"internalType\":\"uint16\",\"name\":\"gasForCallExactCheck\",\"type\":\"uint16\"},{\"internalType\":\"uint32[]\",\"name\":\"maxCallbackGasLimits\",\"type\":\"uint32[]\"}],\"internalType\":\"structFunctionsV1EventsMock.Config\",\"name\":\"param1\",\"type\":\"tuple\"}],\"name\":\"emitConfigUpdated\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"proposedContractSetId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"proposedContractSetFromAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"proposedContractSetToAddress\",\"type\":\"address\"}],\"name\":\"emitContractProposed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"emitContractUpdated\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"emitFundsRecovered\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"emitOwnershipTransferRequested\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"emitOwnershipTransferred\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"emitPaused\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"resultCode\",\"type\":\"uint8\"}],\"name\":\"emitRequestNotProcessed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"uint96\",\"name\":\"totalCostJuels\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"resultCode\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"response\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"callbackReturnData\",\"type\":\"bytes\"}],\"name\":\"emitRequestProcessed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"donId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"subscriptionOwner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"requestingContract\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"requestInitiator\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint16\",\"name\":\"dataVersion\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"estimatedTotalCostJuels\",\"type\":\"uint96\"}],\"name\":\"emitRequestStart\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"emitRequestTimedOut\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"fundsRecipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"fundsAmount\",\"type\":\"uint256\"}],\"name\":\"emitSubscriptionCanceled\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"emitSubscriptionConsumerAdded\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"emitSubscriptionConsumerRemoved\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"emitSubscriptionCreated\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"emitSubscriptionFunded\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"emitSubscriptionOwnerTransferRequested\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"emitSubscriptionOwnerTransferred\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"emitUnpaused\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506111a8806100206000396000f3fe608060405234801561001057600080fd5b50600436106101515760003560e01c8063a5257226116100cd578063e0f6eff111610081578063e9bfcd1811610066578063e9bfcd1814610288578063f7420bc21461029b578063fa7dd96b146102ae57600080fd5b8063e0f6eff114610262578063e2cab57b1461027557600080fd5b8063b24a02cb116100b2578063b24a02cb14610229578063ce150ef11461023c578063dde69b3f1461024f57600080fd5b8063a525722614610203578063b019b4e81461021657600080fd5b8063689300ea116101245780637e1b44c0116101095780637e1b44c0146101ca57806389d38eb4146101dd5780639ec3ce4b146101f057600080fd5b8063689300ea146101a45780637be5c756146101b757600080fd5b8063027d7d22146101565780633f70afb61461016b5780634bf6a80d1461017e578063675b924414610191575b600080fd5b610169610164366004610919565b6102c1565b005b61016961017936600461097e565b610323565b61016961018c3660046109b1565b61037d565b61016961019f36600461097e565b6103df565b6101696101b23660046109f4565b610431565b6101696101c5366004610a1e565b610484565b6101696101d8366004610a40565b6104d1565b6101696101eb366004610bd0565b6104ff565b6101696101fe366004610a1e565b61055b565b61016961021136600461097e565b6105a1565b610169610224366004610c97565b6105f3565b610169610237366004610cb3565b610651565b61016961024a366004610cd8565b6106b1565b61016961025d366004610dad565b610708565b6101696102703660046109b1565b610760565b610169610283366004610de9565b6107b9565b610169610296366004610cb3565b6107fb565b6101696102a9366004610c97565b610852565b6101696102bc366004610ea3565b6108b0565b6040805173ffffffffffffffffffffffffffffffffffffffff85811682528416602082015260ff831681830152905185917f1a90e9a50793db2e394cf581e7c522e10c358a81e70acf6b5a0edd620c08dee1919081900360600190a250505050565b60405173ffffffffffffffffffffffffffffffffffffffff8216815267ffffffffffffffff8316907f464722b4166576d3dcbba877b999bc35cf911f4eaf434b7eba68fa113951d0bf906020015b60405180910390a25050565b6040805173ffffffffffffffffffffffffffffffffffffffff80851682528316602082015267ffffffffffffffff8516917f6f1dc65165ffffedfd8e507b4a0f1fcfdada045ed11f6c26ba27cedfe87802f091015b60405180910390a2505050565b60405173ffffffffffffffffffffffffffffffffffffffff8216815267ffffffffffffffff8316907f43dc749a04ac8fb825cbd514f7c0e13f13bc6f2ee66043b76629d51776cff8e090602001610371565b6040805173ffffffffffffffffffffffffffffffffffffffff84168152602081018390527f59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b436600910160405180910390a15050565b60405173ffffffffffffffffffffffffffffffffffffffff821681527f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258906020015b60405180910390a150565b60405181907ff1ca1e9147be737b04a2b018a79405f687a97de8dd8a2559bbe62357343af41490600090a250565b8767ffffffffffffffff16898b7ff67aec45c9a7ede407974a3e0c3a743dffeab99ee3f2d4c9a8144c2ebf2c7ec98a8a8a8a8a8a8a6040516105479796959493929190610fef565b60405180910390a450505050505050505050565b60405173ffffffffffffffffffffffffffffffffffffffff821681527f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa906020016104c6565b60405173ffffffffffffffffffffffffffffffffffffffff8216815267ffffffffffffffff8316907f182bff9831466789164ca77075fffd84916d35a8180ba73c27e45634549b445b90602001610371565b8073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b6040805184815273ffffffffffffffffffffffffffffffffffffffff80851660208301528316918101919091527f8b052f0f4bf82fede7daffea71592b29d5ef86af1f3c7daaa0345dbb2f52f481906060015b60405180910390a1505050565b8667ffffffffffffffff16887f64778f26c70b60a8d7e29e2451b3844302d959448401c0535b768ed88c6b505e8888888888886040516106f696959493929190611067565b60405180910390a35050505050505050565b6040805173ffffffffffffffffffffffffffffffffffffffff841681526020810183905267ffffffffffffffff8516917fe8ed5b475a5b5987aa9165e8731bb78043f39eee32ec5a1169a89e27fcd4981591016103d2565b6040805173ffffffffffffffffffffffffffffffffffffffff80851682528316602082015267ffffffffffffffff8516917f69436ea6df009049404f564eff6622cd00522b0bd6a89efd9e52a355c4a879be91016103d2565b604080518381526020810183905267ffffffffffffffff8516917fd39ec07f4e209f627a4c427971473820dc129761ba28de8906bd56f57101d4f891016103d2565b6040805184815273ffffffffffffffffffffffffffffffffffffffff80851660208301528316918101919091527ff8a6175bca1ba37d682089187edc5e20a859989727f10ca6bd9a5bc0de8caf94906060016106a4565b8073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127860405160405180910390a35050565b7f049ce2e6e1420eb4b07b425e90129186833eb346bda40b37d5d921aad482f71c816040516104c691906110e6565b803573ffffffffffffffffffffffffffffffffffffffff8116811461090357600080fd5b919050565b803560ff8116811461090357600080fd5b6000806000806080858703121561092f57600080fd5b8435935061093f602086016108df565b925061094d604086016108df565b915061095b60608601610908565b905092959194509250565b803567ffffffffffffffff8116811461090357600080fd5b6000806040838503121561099157600080fd5b61099a83610966565b91506109a8602084016108df565b90509250929050565b6000806000606084860312156109c657600080fd5b6109cf84610966565b92506109dd602085016108df565b91506109eb604085016108df565b90509250925092565b60008060408385031215610a0757600080fd5b610a10836108df565b946020939093013593505050565b600060208284031215610a3057600080fd5b610a39826108df565b9392505050565b600060208284031215610a5257600080fd5b5035919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160a0810167ffffffffffffffff81118282101715610aab57610aab610a59565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610af857610af8610a59565b604052919050565b600082601f830112610b1157600080fd5b813567ffffffffffffffff811115610b2b57610b2b610a59565b610b5c60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601610ab1565b818152846020838601011115610b7157600080fd5b816020850160208301376000918101602001919091529392505050565b803561ffff8116811461090357600080fd5b803563ffffffff8116811461090357600080fd5b80356bffffffffffffffffffffffff8116811461090357600080fd5b6000806000806000806000806000806101408b8d031215610bf057600080fd5b8a35995060208b01359850610c0760408c01610966565b9750610c1560608c016108df565b9650610c2360808c016108df565b9550610c3160a08c016108df565b945060c08b013567ffffffffffffffff811115610c4d57600080fd5b610c598d828e01610b00565b945050610c6860e08c01610b8e565b9250610c776101008c01610ba0565b9150610c866101208c01610bb4565b90509295989b9194979a5092959850565b60008060408385031215610caa57600080fd5b61099a836108df565b600080600060608486031215610cc857600080fd5b833592506109dd602085016108df565b600080600080600080600080610100898b031215610cf557600080fd5b88359750610d0560208a01610966565b9650610d1360408a01610bb4565b9550610d2160608a016108df565b9450610d2f60808a01610908565b935060a089013567ffffffffffffffff80821115610d4c57600080fd5b610d588c838d01610b00565b945060c08b0135915080821115610d6e57600080fd5b610d7a8c838d01610b00565b935060e08b0135915080821115610d9057600080fd5b50610d9d8b828c01610b00565b9150509295985092959890939650565b600080600060608486031215610dc257600080fd5b610dcb84610966565b9250610dd9602085016108df565b9150604084013590509250925092565b600080600060608486031215610dfe57600080fd5b610e0784610966565b95602085013595506040909401359392505050565b600082601f830112610e2d57600080fd5b8135602067ffffffffffffffff821115610e4957610e49610a59565b8160051b610e58828201610ab1565b9283528481018201928281019087851115610e7257600080fd5b83870192505b84831015610e9857610e8983610ba0565b82529183019190830190610e78565b979650505050505050565b600060208284031215610eb557600080fd5b813567ffffffffffffffff80821115610ecd57600080fd5b9083019060a08286031215610ee157600080fd5b610ee9610a88565b610ef283610b8e565b8152602083013568ffffffffffffffffff81168114610f1057600080fd5b602082015260408301357fffffffff0000000000000000000000000000000000000000000000000000000081168114610f4857600080fd5b6040820152610f5960608401610b8e565b6060820152608083013582811115610f7057600080fd5b610f7c87828601610e1c565b60808301525095945050505050565b6000815180845260005b81811015610fb157602081850181015186830182015201610f95565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b600073ffffffffffffffffffffffffffffffffffffffff808a168352808916602084015280881660408401525060e0606083015261103060e0830187610f8b565b61ffff9590951660808301525063ffffffff9290921660a08301526bffffffffffffffffffffffff1660c090910152949350505050565b6bffffffffffffffffffffffff8716815273ffffffffffffffffffffffffffffffffffffffff8616602082015260ff8516604082015260c0606082015260006110b360c0830186610f8b565b82810360808401526110c58186610f8b565b905082810360a08401526110d98185610f8b565b9998505050505050505050565b6000602080835260c0830161ffff808651168386015268ffffffffffffffffff838701511660408601527fffffffff00000000000000000000000000000000000000000000000000000000604087015116606086015280606087015116608086015250608085015160a08086015281815180845260e0870191508483019350600092505b8083101561119057835163ffffffff16825292840192600192909201919084019061116a565b50969550505050505056fea164736f6c6343000813000a",
}

var FunctionsV1EventsMockABI = FunctionsV1EventsMockMetaData.ABI

var FunctionsV1EventsMockBin = FunctionsV1EventsMockMetaData.Bin

func DeployFunctionsV1EventsMock(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *FunctionsV1EventsMock, error) {
	parsed, err := FunctionsV1EventsMockMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(FunctionsV1EventsMockBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &FunctionsV1EventsMock{address: address, abi: *parsed, FunctionsV1EventsMockCaller: FunctionsV1EventsMockCaller{contract: contract}, FunctionsV1EventsMockTransactor: FunctionsV1EventsMockTransactor{contract: contract}, FunctionsV1EventsMockFilterer: FunctionsV1EventsMockFilterer{contract: contract}}, nil
}

type FunctionsV1EventsMock struct {
	address common.Address
	abi     abi.ABI
	FunctionsV1EventsMockCaller
	FunctionsV1EventsMockTransactor
	FunctionsV1EventsMockFilterer
}

type FunctionsV1EventsMockCaller struct {
	contract *bind.BoundContract
}

type FunctionsV1EventsMockTransactor struct {
	contract *bind.BoundContract
}

type FunctionsV1EventsMockFilterer struct {
	contract *bind.BoundContract
}

type FunctionsV1EventsMockSession struct {
	Contract     *FunctionsV1EventsMock
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type FunctionsV1EventsMockCallerSession struct {
	Contract *FunctionsV1EventsMockCaller
	CallOpts bind.CallOpts
}

type FunctionsV1EventsMockTransactorSession struct {
	Contract     *FunctionsV1EventsMockTransactor
	TransactOpts bind.TransactOpts
}

type FunctionsV1EventsMockRaw struct {
	Contract *FunctionsV1EventsMock
}

type FunctionsV1EventsMockCallerRaw struct {
	Contract *FunctionsV1EventsMockCaller
}

type FunctionsV1EventsMockTransactorRaw struct {
	Contract *FunctionsV1EventsMockTransactor
}

func NewFunctionsV1EventsMock(address common.Address, backend bind.ContractBackend) (*FunctionsV1EventsMock, error) {
	abi, err := abi.JSON(strings.NewReader(FunctionsV1EventsMockABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindFunctionsV1EventsMock(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FunctionsV1EventsMock{address: address, abi: abi, FunctionsV1EventsMockCaller: FunctionsV1EventsMockCaller{contract: contract}, FunctionsV1EventsMockTransactor: FunctionsV1EventsMockTransactor{contract: contract}, FunctionsV1EventsMockFilterer: FunctionsV1EventsMockFilterer{contract: contract}}, nil
}

func NewFunctionsV1EventsMockCaller(address common.Address, caller bind.ContractCaller) (*FunctionsV1EventsMockCaller, error) {
	contract, err := bindFunctionsV1EventsMock(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FunctionsV1EventsMockCaller{contract: contract}, nil
}

func NewFunctionsV1EventsMockTransactor(address common.Address, transactor bind.ContractTransactor) (*FunctionsV1EventsMockTransactor, error) {
	contract, err := bindFunctionsV1EventsMock(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FunctionsV1EventsMockTransactor{contract: contract}, nil
}

func NewFunctionsV1EventsMockFilterer(address common.Address, filterer bind.ContractFilterer) (*FunctionsV1EventsMockFilterer, error) {
	contract, err := bindFunctionsV1EventsMock(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FunctionsV1EventsMockFilterer{contract: contract}, nil
}

func bindFunctionsV1EventsMock(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FunctionsV1EventsMockMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FunctionsV1EventsMock.Contract.FunctionsV1EventsMockCaller.contract.Call(opts, result, method, params...)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.Contract.FunctionsV1EventsMockTransactor.contract.Transfer(opts)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.Contract.FunctionsV1EventsMockTransactor.contract.Transact(opts, method, params...)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FunctionsV1EventsMock.Contract.contract.Call(opts, result, method, params...)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.Contract.contract.Transfer(opts)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.Contract.contract.Transact(opts, method, params...)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockTransactor) EmitConfigUpdated(opts *bind.TransactOpts, param1 FunctionsV1EventsMockConfig) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.contract.Transact(opts, "emitConfigUpdated", param1)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockSession) EmitConfigUpdated(param1 FunctionsV1EventsMockConfig) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.Contract.EmitConfigUpdated(&_FunctionsV1EventsMock.TransactOpts, param1)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockTransactorSession) EmitConfigUpdated(param1 FunctionsV1EventsMockConfig) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.Contract.EmitConfigUpdated(&_FunctionsV1EventsMock.TransactOpts, param1)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockTransactor) EmitContractProposed(opts *bind.TransactOpts, proposedContractSetId [32]byte, proposedContractSetFromAddress common.Address, proposedContractSetToAddress common.Address) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.contract.Transact(opts, "emitContractProposed", proposedContractSetId, proposedContractSetFromAddress, proposedContractSetToAddress)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockSession) EmitContractProposed(proposedContractSetId [32]byte, proposedContractSetFromAddress common.Address, proposedContractSetToAddress common.Address) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.Contract.EmitContractProposed(&_FunctionsV1EventsMock.TransactOpts, proposedContractSetId, proposedContractSetFromAddress, proposedContractSetToAddress)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockTransactorSession) EmitContractProposed(proposedContractSetId [32]byte, proposedContractSetFromAddress common.Address, proposedContractSetToAddress common.Address) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.Contract.EmitContractProposed(&_FunctionsV1EventsMock.TransactOpts, proposedContractSetId, proposedContractSetFromAddress, proposedContractSetToAddress)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockTransactor) EmitContractUpdated(opts *bind.TransactOpts, id [32]byte, from common.Address, to common.Address) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.contract.Transact(opts, "emitContractUpdated", id, from, to)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockSession) EmitContractUpdated(id [32]byte, from common.Address, to common.Address) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.Contract.EmitContractUpdated(&_FunctionsV1EventsMock.TransactOpts, id, from, to)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockTransactorSession) EmitContractUpdated(id [32]byte, from common.Address, to common.Address) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.Contract.EmitContractUpdated(&_FunctionsV1EventsMock.TransactOpts, id, from, to)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockTransactor) EmitFundsRecovered(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.contract.Transact(opts, "emitFundsRecovered", to, amount)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockSession) EmitFundsRecovered(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.Contract.EmitFundsRecovered(&_FunctionsV1EventsMock.TransactOpts, to, amount)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockTransactorSession) EmitFundsRecovered(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.Contract.EmitFundsRecovered(&_FunctionsV1EventsMock.TransactOpts, to, amount)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockTransactor) EmitOwnershipTransferRequested(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.contract.Transact(opts, "emitOwnershipTransferRequested", from, to)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockSession) EmitOwnershipTransferRequested(from common.Address, to common.Address) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.Contract.EmitOwnershipTransferRequested(&_FunctionsV1EventsMock.TransactOpts, from, to)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockTransactorSession) EmitOwnershipTransferRequested(from common.Address, to common.Address) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.Contract.EmitOwnershipTransferRequested(&_FunctionsV1EventsMock.TransactOpts, from, to)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockTransactor) EmitOwnershipTransferred(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.contract.Transact(opts, "emitOwnershipTransferred", from, to)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockSession) EmitOwnershipTransferred(from common.Address, to common.Address) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.Contract.EmitOwnershipTransferred(&_FunctionsV1EventsMock.TransactOpts, from, to)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockTransactorSession) EmitOwnershipTransferred(from common.Address, to common.Address) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.Contract.EmitOwnershipTransferred(&_FunctionsV1EventsMock.TransactOpts, from, to)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockTransactor) EmitPaused(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.contract.Transact(opts, "emitPaused", account)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockSession) EmitPaused(account common.Address) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.Contract.EmitPaused(&_FunctionsV1EventsMock.TransactOpts, account)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockTransactorSession) EmitPaused(account common.Address) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.Contract.EmitPaused(&_FunctionsV1EventsMock.TransactOpts, account)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockTransactor) EmitRequestNotProcessed(opts *bind.TransactOpts, requestId [32]byte, coordinator common.Address, transmitter common.Address, resultCode uint8) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.contract.Transact(opts, "emitRequestNotProcessed", requestId, coordinator, transmitter, resultCode)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockSession) EmitRequestNotProcessed(requestId [32]byte, coordinator common.Address, transmitter common.Address, resultCode uint8) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.Contract.EmitRequestNotProcessed(&_FunctionsV1EventsMock.TransactOpts, requestId, coordinator, transmitter, resultCode)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockTransactorSession) EmitRequestNotProcessed(requestId [32]byte, coordinator common.Address, transmitter common.Address, resultCode uint8) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.Contract.EmitRequestNotProcessed(&_FunctionsV1EventsMock.TransactOpts, requestId, coordinator, transmitter, resultCode)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockTransactor) EmitRequestProcessed(opts *bind.TransactOpts, requestId [32]byte, subscriptionId uint64, totalCostJuels *big.Int, transmitter common.Address, resultCode uint8, response []byte, err []byte, callbackReturnData []byte) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.contract.Transact(opts, "emitRequestProcessed", requestId, subscriptionId, totalCostJuels, transmitter, resultCode, response, err, callbackReturnData)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockSession) EmitRequestProcessed(requestId [32]byte, subscriptionId uint64, totalCostJuels *big.Int, transmitter common.Address, resultCode uint8, response []byte, err []byte, callbackReturnData []byte) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.Contract.EmitRequestProcessed(&_FunctionsV1EventsMock.TransactOpts, requestId, subscriptionId, totalCostJuels, transmitter, resultCode, response, err, callbackReturnData)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockTransactorSession) EmitRequestProcessed(requestId [32]byte, subscriptionId uint64, totalCostJuels *big.Int, transmitter common.Address, resultCode uint8, response []byte, err []byte, callbackReturnData []byte) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.Contract.EmitRequestProcessed(&_FunctionsV1EventsMock.TransactOpts, requestId, subscriptionId, totalCostJuels, transmitter, resultCode, response, err, callbackReturnData)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockTransactor) EmitRequestStart(opts *bind.TransactOpts, requestId [32]byte, donId [32]byte, subscriptionId uint64, subscriptionOwner common.Address, requestingContract common.Address, requestInitiator common.Address, data []byte, dataVersion uint16, callbackGasLimit uint32, estimatedTotalCostJuels *big.Int) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.contract.Transact(opts, "emitRequestStart", requestId, donId, subscriptionId, subscriptionOwner, requestingContract, requestInitiator, data, dataVersion, callbackGasLimit, estimatedTotalCostJuels)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockSession) EmitRequestStart(requestId [32]byte, donId [32]byte, subscriptionId uint64, subscriptionOwner common.Address, requestingContract common.Address, requestInitiator common.Address, data []byte, dataVersion uint16, callbackGasLimit uint32, estimatedTotalCostJuels *big.Int) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.Contract.EmitRequestStart(&_FunctionsV1EventsMock.TransactOpts, requestId, donId, subscriptionId, subscriptionOwner, requestingContract, requestInitiator, data, dataVersion, callbackGasLimit, estimatedTotalCostJuels)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockTransactorSession) EmitRequestStart(requestId [32]byte, donId [32]byte, subscriptionId uint64, subscriptionOwner common.Address, requestingContract common.Address, requestInitiator common.Address, data []byte, dataVersion uint16, callbackGasLimit uint32, estimatedTotalCostJuels *big.Int) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.Contract.EmitRequestStart(&_FunctionsV1EventsMock.TransactOpts, requestId, donId, subscriptionId, subscriptionOwner, requestingContract, requestInitiator, data, dataVersion, callbackGasLimit, estimatedTotalCostJuels)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockTransactor) EmitRequestTimedOut(opts *bind.TransactOpts, requestId [32]byte) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.contract.Transact(opts, "emitRequestTimedOut", requestId)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockSession) EmitRequestTimedOut(requestId [32]byte) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.Contract.EmitRequestTimedOut(&_FunctionsV1EventsMock.TransactOpts, requestId)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockTransactorSession) EmitRequestTimedOut(requestId [32]byte) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.Contract.EmitRequestTimedOut(&_FunctionsV1EventsMock.TransactOpts, requestId)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockTransactor) EmitSubscriptionCanceled(opts *bind.TransactOpts, subscriptionId uint64, fundsRecipient common.Address, fundsAmount *big.Int) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.contract.Transact(opts, "emitSubscriptionCanceled", subscriptionId, fundsRecipient, fundsAmount)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockSession) EmitSubscriptionCanceled(subscriptionId uint64, fundsRecipient common.Address, fundsAmount *big.Int) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.Contract.EmitSubscriptionCanceled(&_FunctionsV1EventsMock.TransactOpts, subscriptionId, fundsRecipient, fundsAmount)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockTransactorSession) EmitSubscriptionCanceled(subscriptionId uint64, fundsRecipient common.Address, fundsAmount *big.Int) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.Contract.EmitSubscriptionCanceled(&_FunctionsV1EventsMock.TransactOpts, subscriptionId, fundsRecipient, fundsAmount)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockTransactor) EmitSubscriptionConsumerAdded(opts *bind.TransactOpts, subscriptionId uint64, consumer common.Address) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.contract.Transact(opts, "emitSubscriptionConsumerAdded", subscriptionId, consumer)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockSession) EmitSubscriptionConsumerAdded(subscriptionId uint64, consumer common.Address) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.Contract.EmitSubscriptionConsumerAdded(&_FunctionsV1EventsMock.TransactOpts, subscriptionId, consumer)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockTransactorSession) EmitSubscriptionConsumerAdded(subscriptionId uint64, consumer common.Address) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.Contract.EmitSubscriptionConsumerAdded(&_FunctionsV1EventsMock.TransactOpts, subscriptionId, consumer)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockTransactor) EmitSubscriptionConsumerRemoved(opts *bind.TransactOpts, subscriptionId uint64, consumer common.Address) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.contract.Transact(opts, "emitSubscriptionConsumerRemoved", subscriptionId, consumer)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockSession) EmitSubscriptionConsumerRemoved(subscriptionId uint64, consumer common.Address) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.Contract.EmitSubscriptionConsumerRemoved(&_FunctionsV1EventsMock.TransactOpts, subscriptionId, consumer)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockTransactorSession) EmitSubscriptionConsumerRemoved(subscriptionId uint64, consumer common.Address) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.Contract.EmitSubscriptionConsumerRemoved(&_FunctionsV1EventsMock.TransactOpts, subscriptionId, consumer)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockTransactor) EmitSubscriptionCreated(opts *bind.TransactOpts, subscriptionId uint64, owner common.Address) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.contract.Transact(opts, "emitSubscriptionCreated", subscriptionId, owner)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockSession) EmitSubscriptionCreated(subscriptionId uint64, owner common.Address) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.Contract.EmitSubscriptionCreated(&_FunctionsV1EventsMock.TransactOpts, subscriptionId, owner)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockTransactorSession) EmitSubscriptionCreated(subscriptionId uint64, owner common.Address) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.Contract.EmitSubscriptionCreated(&_FunctionsV1EventsMock.TransactOpts, subscriptionId, owner)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockTransactor) EmitSubscriptionFunded(opts *bind.TransactOpts, subscriptionId uint64, oldBalance *big.Int, newBalance *big.Int) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.contract.Transact(opts, "emitSubscriptionFunded", subscriptionId, oldBalance, newBalance)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockSession) EmitSubscriptionFunded(subscriptionId uint64, oldBalance *big.Int, newBalance *big.Int) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.Contract.EmitSubscriptionFunded(&_FunctionsV1EventsMock.TransactOpts, subscriptionId, oldBalance, newBalance)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockTransactorSession) EmitSubscriptionFunded(subscriptionId uint64, oldBalance *big.Int, newBalance *big.Int) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.Contract.EmitSubscriptionFunded(&_FunctionsV1EventsMock.TransactOpts, subscriptionId, oldBalance, newBalance)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockTransactor) EmitSubscriptionOwnerTransferRequested(opts *bind.TransactOpts, subscriptionId uint64, from common.Address, to common.Address) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.contract.Transact(opts, "emitSubscriptionOwnerTransferRequested", subscriptionId, from, to)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockSession) EmitSubscriptionOwnerTransferRequested(subscriptionId uint64, from common.Address, to common.Address) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.Contract.EmitSubscriptionOwnerTransferRequested(&_FunctionsV1EventsMock.TransactOpts, subscriptionId, from, to)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockTransactorSession) EmitSubscriptionOwnerTransferRequested(subscriptionId uint64, from common.Address, to common.Address) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.Contract.EmitSubscriptionOwnerTransferRequested(&_FunctionsV1EventsMock.TransactOpts, subscriptionId, from, to)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockTransactor) EmitSubscriptionOwnerTransferred(opts *bind.TransactOpts, subscriptionId uint64, from common.Address, to common.Address) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.contract.Transact(opts, "emitSubscriptionOwnerTransferred", subscriptionId, from, to)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockSession) EmitSubscriptionOwnerTransferred(subscriptionId uint64, from common.Address, to common.Address) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.Contract.EmitSubscriptionOwnerTransferred(&_FunctionsV1EventsMock.TransactOpts, subscriptionId, from, to)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockTransactorSession) EmitSubscriptionOwnerTransferred(subscriptionId uint64, from common.Address, to common.Address) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.Contract.EmitSubscriptionOwnerTransferred(&_FunctionsV1EventsMock.TransactOpts, subscriptionId, from, to)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockTransactor) EmitUnpaused(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.contract.Transact(opts, "emitUnpaused", account)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockSession) EmitUnpaused(account common.Address) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.Contract.EmitUnpaused(&_FunctionsV1EventsMock.TransactOpts, account)
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockTransactorSession) EmitUnpaused(account common.Address) (*types.Transaction, error) {
	return _FunctionsV1EventsMock.Contract.EmitUnpaused(&_FunctionsV1EventsMock.TransactOpts, account)
}

type FunctionsV1EventsMockConfigUpdatedIterator struct {
	Event *FunctionsV1EventsMockConfigUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsV1EventsMockConfigUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsV1EventsMockConfigUpdated)
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
		it.Event = new(FunctionsV1EventsMockConfigUpdated)
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

func (it *FunctionsV1EventsMockConfigUpdatedIterator) Error() error {
	return it.fail
}

func (it *FunctionsV1EventsMockConfigUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsV1EventsMockConfigUpdated struct {
	Param1 FunctionsV1EventsMockConfig
	Raw    types.Log
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) FilterConfigUpdated(opts *bind.FilterOpts) (*FunctionsV1EventsMockConfigUpdatedIterator, error) {

	logs, sub, err := _FunctionsV1EventsMock.contract.FilterLogs(opts, "ConfigUpdated")
	if err != nil {
		return nil, err
	}
	return &FunctionsV1EventsMockConfigUpdatedIterator{contract: _FunctionsV1EventsMock.contract, event: "ConfigUpdated", logs: logs, sub: sub}, nil
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) WatchConfigUpdated(opts *bind.WatchOpts, sink chan<- *FunctionsV1EventsMockConfigUpdated) (event.Subscription, error) {

	logs, sub, err := _FunctionsV1EventsMock.contract.WatchLogs(opts, "ConfigUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsV1EventsMockConfigUpdated)
				if err := _FunctionsV1EventsMock.contract.UnpackLog(event, "ConfigUpdated", log); err != nil {
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

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) ParseConfigUpdated(log types.Log) (*FunctionsV1EventsMockConfigUpdated, error) {
	event := new(FunctionsV1EventsMockConfigUpdated)
	if err := _FunctionsV1EventsMock.contract.UnpackLog(event, "ConfigUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsV1EventsMockContractProposedIterator struct {
	Event *FunctionsV1EventsMockContractProposed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsV1EventsMockContractProposedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsV1EventsMockContractProposed)
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
		it.Event = new(FunctionsV1EventsMockContractProposed)
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

func (it *FunctionsV1EventsMockContractProposedIterator) Error() error {
	return it.fail
}

func (it *FunctionsV1EventsMockContractProposedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsV1EventsMockContractProposed struct {
	ProposedContractSetId          [32]byte
	ProposedContractSetFromAddress common.Address
	ProposedContractSetToAddress   common.Address
	Raw                            types.Log
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) FilterContractProposed(opts *bind.FilterOpts) (*FunctionsV1EventsMockContractProposedIterator, error) {

	logs, sub, err := _FunctionsV1EventsMock.contract.FilterLogs(opts, "ContractProposed")
	if err != nil {
		return nil, err
	}
	return &FunctionsV1EventsMockContractProposedIterator{contract: _FunctionsV1EventsMock.contract, event: "ContractProposed", logs: logs, sub: sub}, nil
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) WatchContractProposed(opts *bind.WatchOpts, sink chan<- *FunctionsV1EventsMockContractProposed) (event.Subscription, error) {

	logs, sub, err := _FunctionsV1EventsMock.contract.WatchLogs(opts, "ContractProposed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsV1EventsMockContractProposed)
				if err := _FunctionsV1EventsMock.contract.UnpackLog(event, "ContractProposed", log); err != nil {
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

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) ParseContractProposed(log types.Log) (*FunctionsV1EventsMockContractProposed, error) {
	event := new(FunctionsV1EventsMockContractProposed)
	if err := _FunctionsV1EventsMock.contract.UnpackLog(event, "ContractProposed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsV1EventsMockContractUpdatedIterator struct {
	Event *FunctionsV1EventsMockContractUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsV1EventsMockContractUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsV1EventsMockContractUpdated)
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
		it.Event = new(FunctionsV1EventsMockContractUpdated)
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

func (it *FunctionsV1EventsMockContractUpdatedIterator) Error() error {
	return it.fail
}

func (it *FunctionsV1EventsMockContractUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsV1EventsMockContractUpdated struct {
	Id   [32]byte
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) FilterContractUpdated(opts *bind.FilterOpts) (*FunctionsV1EventsMockContractUpdatedIterator, error) {

	logs, sub, err := _FunctionsV1EventsMock.contract.FilterLogs(opts, "ContractUpdated")
	if err != nil {
		return nil, err
	}
	return &FunctionsV1EventsMockContractUpdatedIterator{contract: _FunctionsV1EventsMock.contract, event: "ContractUpdated", logs: logs, sub: sub}, nil
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) WatchContractUpdated(opts *bind.WatchOpts, sink chan<- *FunctionsV1EventsMockContractUpdated) (event.Subscription, error) {

	logs, sub, err := _FunctionsV1EventsMock.contract.WatchLogs(opts, "ContractUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsV1EventsMockContractUpdated)
				if err := _FunctionsV1EventsMock.contract.UnpackLog(event, "ContractUpdated", log); err != nil {
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

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) ParseContractUpdated(log types.Log) (*FunctionsV1EventsMockContractUpdated, error) {
	event := new(FunctionsV1EventsMockContractUpdated)
	if err := _FunctionsV1EventsMock.contract.UnpackLog(event, "ContractUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsV1EventsMockFundsRecoveredIterator struct {
	Event *FunctionsV1EventsMockFundsRecovered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsV1EventsMockFundsRecoveredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsV1EventsMockFundsRecovered)
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
		it.Event = new(FunctionsV1EventsMockFundsRecovered)
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

func (it *FunctionsV1EventsMockFundsRecoveredIterator) Error() error {
	return it.fail
}

func (it *FunctionsV1EventsMockFundsRecoveredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsV1EventsMockFundsRecovered struct {
	To     common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) FilterFundsRecovered(opts *bind.FilterOpts) (*FunctionsV1EventsMockFundsRecoveredIterator, error) {

	logs, sub, err := _FunctionsV1EventsMock.contract.FilterLogs(opts, "FundsRecovered")
	if err != nil {
		return nil, err
	}
	return &FunctionsV1EventsMockFundsRecoveredIterator{contract: _FunctionsV1EventsMock.contract, event: "FundsRecovered", logs: logs, sub: sub}, nil
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) WatchFundsRecovered(opts *bind.WatchOpts, sink chan<- *FunctionsV1EventsMockFundsRecovered) (event.Subscription, error) {

	logs, sub, err := _FunctionsV1EventsMock.contract.WatchLogs(opts, "FundsRecovered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsV1EventsMockFundsRecovered)
				if err := _FunctionsV1EventsMock.contract.UnpackLog(event, "FundsRecovered", log); err != nil {
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

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) ParseFundsRecovered(log types.Log) (*FunctionsV1EventsMockFundsRecovered, error) {
	event := new(FunctionsV1EventsMockFundsRecovered)
	if err := _FunctionsV1EventsMock.contract.UnpackLog(event, "FundsRecovered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsV1EventsMockOwnershipTransferRequestedIterator struct {
	Event *FunctionsV1EventsMockOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsV1EventsMockOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsV1EventsMockOwnershipTransferRequested)
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
		it.Event = new(FunctionsV1EventsMockOwnershipTransferRequested)
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

func (it *FunctionsV1EventsMockOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *FunctionsV1EventsMockOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsV1EventsMockOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsV1EventsMockOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsV1EventsMock.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsV1EventsMockOwnershipTransferRequestedIterator{contract: _FunctionsV1EventsMock.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *FunctionsV1EventsMockOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsV1EventsMock.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsV1EventsMockOwnershipTransferRequested)
				if err := _FunctionsV1EventsMock.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) ParseOwnershipTransferRequested(log types.Log) (*FunctionsV1EventsMockOwnershipTransferRequested, error) {
	event := new(FunctionsV1EventsMockOwnershipTransferRequested)
	if err := _FunctionsV1EventsMock.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsV1EventsMockOwnershipTransferredIterator struct {
	Event *FunctionsV1EventsMockOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsV1EventsMockOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsV1EventsMockOwnershipTransferred)
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
		it.Event = new(FunctionsV1EventsMockOwnershipTransferred)
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

func (it *FunctionsV1EventsMockOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *FunctionsV1EventsMockOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsV1EventsMockOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsV1EventsMockOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsV1EventsMock.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsV1EventsMockOwnershipTransferredIterator{contract: _FunctionsV1EventsMock.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FunctionsV1EventsMockOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsV1EventsMock.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsV1EventsMockOwnershipTransferred)
				if err := _FunctionsV1EventsMock.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) ParseOwnershipTransferred(log types.Log) (*FunctionsV1EventsMockOwnershipTransferred, error) {
	event := new(FunctionsV1EventsMockOwnershipTransferred)
	if err := _FunctionsV1EventsMock.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsV1EventsMockPausedIterator struct {
	Event *FunctionsV1EventsMockPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsV1EventsMockPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsV1EventsMockPaused)
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
		it.Event = new(FunctionsV1EventsMockPaused)
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

func (it *FunctionsV1EventsMockPausedIterator) Error() error {
	return it.fail
}

func (it *FunctionsV1EventsMockPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsV1EventsMockPaused struct {
	Account common.Address
	Raw     types.Log
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) FilterPaused(opts *bind.FilterOpts) (*FunctionsV1EventsMockPausedIterator, error) {

	logs, sub, err := _FunctionsV1EventsMock.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &FunctionsV1EventsMockPausedIterator{contract: _FunctionsV1EventsMock.contract, event: "Paused", logs: logs, sub: sub}, nil
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *FunctionsV1EventsMockPaused) (event.Subscription, error) {

	logs, sub, err := _FunctionsV1EventsMock.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsV1EventsMockPaused)
				if err := _FunctionsV1EventsMock.contract.UnpackLog(event, "Paused", log); err != nil {
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

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) ParsePaused(log types.Log) (*FunctionsV1EventsMockPaused, error) {
	event := new(FunctionsV1EventsMockPaused)
	if err := _FunctionsV1EventsMock.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsV1EventsMockRequestNotProcessedIterator struct {
	Event *FunctionsV1EventsMockRequestNotProcessed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsV1EventsMockRequestNotProcessedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsV1EventsMockRequestNotProcessed)
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
		it.Event = new(FunctionsV1EventsMockRequestNotProcessed)
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

func (it *FunctionsV1EventsMockRequestNotProcessedIterator) Error() error {
	return it.fail
}

func (it *FunctionsV1EventsMockRequestNotProcessedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsV1EventsMockRequestNotProcessed struct {
	RequestId   [32]byte
	Coordinator common.Address
	Transmitter common.Address
	ResultCode  uint8
	Raw         types.Log
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) FilterRequestNotProcessed(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsV1EventsMockRequestNotProcessedIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsV1EventsMock.contract.FilterLogs(opts, "RequestNotProcessed", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsV1EventsMockRequestNotProcessedIterator{contract: _FunctionsV1EventsMock.contract, event: "RequestNotProcessed", logs: logs, sub: sub}, nil
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) WatchRequestNotProcessed(opts *bind.WatchOpts, sink chan<- *FunctionsV1EventsMockRequestNotProcessed, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsV1EventsMock.contract.WatchLogs(opts, "RequestNotProcessed", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsV1EventsMockRequestNotProcessed)
				if err := _FunctionsV1EventsMock.contract.UnpackLog(event, "RequestNotProcessed", log); err != nil {
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

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) ParseRequestNotProcessed(log types.Log) (*FunctionsV1EventsMockRequestNotProcessed, error) {
	event := new(FunctionsV1EventsMockRequestNotProcessed)
	if err := _FunctionsV1EventsMock.contract.UnpackLog(event, "RequestNotProcessed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsV1EventsMockRequestProcessedIterator struct {
	Event *FunctionsV1EventsMockRequestProcessed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsV1EventsMockRequestProcessedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsV1EventsMockRequestProcessed)
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
		it.Event = new(FunctionsV1EventsMockRequestProcessed)
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

func (it *FunctionsV1EventsMockRequestProcessedIterator) Error() error {
	return it.fail
}

func (it *FunctionsV1EventsMockRequestProcessedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsV1EventsMockRequestProcessed struct {
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

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) FilterRequestProcessed(opts *bind.FilterOpts, requestId [][32]byte, subscriptionId []uint64) (*FunctionsV1EventsMockRequestProcessedIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsV1EventsMock.contract.FilterLogs(opts, "RequestProcessed", requestIdRule, subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsV1EventsMockRequestProcessedIterator{contract: _FunctionsV1EventsMock.contract, event: "RequestProcessed", logs: logs, sub: sub}, nil
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) WatchRequestProcessed(opts *bind.WatchOpts, sink chan<- *FunctionsV1EventsMockRequestProcessed, requestId [][32]byte, subscriptionId []uint64) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsV1EventsMock.contract.WatchLogs(opts, "RequestProcessed", requestIdRule, subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsV1EventsMockRequestProcessed)
				if err := _FunctionsV1EventsMock.contract.UnpackLog(event, "RequestProcessed", log); err != nil {
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

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) ParseRequestProcessed(log types.Log) (*FunctionsV1EventsMockRequestProcessed, error) {
	event := new(FunctionsV1EventsMockRequestProcessed)
	if err := _FunctionsV1EventsMock.contract.UnpackLog(event, "RequestProcessed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsV1EventsMockRequestStartIterator struct {
	Event *FunctionsV1EventsMockRequestStart

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsV1EventsMockRequestStartIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsV1EventsMockRequestStart)
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
		it.Event = new(FunctionsV1EventsMockRequestStart)
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

func (it *FunctionsV1EventsMockRequestStartIterator) Error() error {
	return it.fail
}

func (it *FunctionsV1EventsMockRequestStartIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsV1EventsMockRequestStart struct {
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

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) FilterRequestStart(opts *bind.FilterOpts, requestId [][32]byte, donId [][32]byte, subscriptionId []uint64) (*FunctionsV1EventsMockRequestStartIterator, error) {

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

	logs, sub, err := _FunctionsV1EventsMock.contract.FilterLogs(opts, "RequestStart", requestIdRule, donIdRule, subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsV1EventsMockRequestStartIterator{contract: _FunctionsV1EventsMock.contract, event: "RequestStart", logs: logs, sub: sub}, nil
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) WatchRequestStart(opts *bind.WatchOpts, sink chan<- *FunctionsV1EventsMockRequestStart, requestId [][32]byte, donId [][32]byte, subscriptionId []uint64) (event.Subscription, error) {

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

	logs, sub, err := _FunctionsV1EventsMock.contract.WatchLogs(opts, "RequestStart", requestIdRule, donIdRule, subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsV1EventsMockRequestStart)
				if err := _FunctionsV1EventsMock.contract.UnpackLog(event, "RequestStart", log); err != nil {
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

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) ParseRequestStart(log types.Log) (*FunctionsV1EventsMockRequestStart, error) {
	event := new(FunctionsV1EventsMockRequestStart)
	if err := _FunctionsV1EventsMock.contract.UnpackLog(event, "RequestStart", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsV1EventsMockRequestTimedOutIterator struct {
	Event *FunctionsV1EventsMockRequestTimedOut

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsV1EventsMockRequestTimedOutIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsV1EventsMockRequestTimedOut)
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
		it.Event = new(FunctionsV1EventsMockRequestTimedOut)
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

func (it *FunctionsV1EventsMockRequestTimedOutIterator) Error() error {
	return it.fail
}

func (it *FunctionsV1EventsMockRequestTimedOutIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsV1EventsMockRequestTimedOut struct {
	RequestId [32]byte
	Raw       types.Log
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) FilterRequestTimedOut(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsV1EventsMockRequestTimedOutIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsV1EventsMock.contract.FilterLogs(opts, "RequestTimedOut", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsV1EventsMockRequestTimedOutIterator{contract: _FunctionsV1EventsMock.contract, event: "RequestTimedOut", logs: logs, sub: sub}, nil
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) WatchRequestTimedOut(opts *bind.WatchOpts, sink chan<- *FunctionsV1EventsMockRequestTimedOut, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsV1EventsMock.contract.WatchLogs(opts, "RequestTimedOut", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsV1EventsMockRequestTimedOut)
				if err := _FunctionsV1EventsMock.contract.UnpackLog(event, "RequestTimedOut", log); err != nil {
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

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) ParseRequestTimedOut(log types.Log) (*FunctionsV1EventsMockRequestTimedOut, error) {
	event := new(FunctionsV1EventsMockRequestTimedOut)
	if err := _FunctionsV1EventsMock.contract.UnpackLog(event, "RequestTimedOut", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsV1EventsMockSubscriptionCanceledIterator struct {
	Event *FunctionsV1EventsMockSubscriptionCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsV1EventsMockSubscriptionCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsV1EventsMockSubscriptionCanceled)
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
		it.Event = new(FunctionsV1EventsMockSubscriptionCanceled)
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

func (it *FunctionsV1EventsMockSubscriptionCanceledIterator) Error() error {
	return it.fail
}

func (it *FunctionsV1EventsMockSubscriptionCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsV1EventsMockSubscriptionCanceled struct {
	SubscriptionId uint64
	FundsRecipient common.Address
	FundsAmount    *big.Int
	Raw            types.Log
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) FilterSubscriptionCanceled(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsV1EventsMockSubscriptionCanceledIterator, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsV1EventsMock.contract.FilterLogs(opts, "SubscriptionCanceled", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsV1EventsMockSubscriptionCanceledIterator{contract: _FunctionsV1EventsMock.contract, event: "SubscriptionCanceled", logs: logs, sub: sub}, nil
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) WatchSubscriptionCanceled(opts *bind.WatchOpts, sink chan<- *FunctionsV1EventsMockSubscriptionCanceled, subscriptionId []uint64) (event.Subscription, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsV1EventsMock.contract.WatchLogs(opts, "SubscriptionCanceled", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsV1EventsMockSubscriptionCanceled)
				if err := _FunctionsV1EventsMock.contract.UnpackLog(event, "SubscriptionCanceled", log); err != nil {
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

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) ParseSubscriptionCanceled(log types.Log) (*FunctionsV1EventsMockSubscriptionCanceled, error) {
	event := new(FunctionsV1EventsMockSubscriptionCanceled)
	if err := _FunctionsV1EventsMock.contract.UnpackLog(event, "SubscriptionCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsV1EventsMockSubscriptionConsumerAddedIterator struct {
	Event *FunctionsV1EventsMockSubscriptionConsumerAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsV1EventsMockSubscriptionConsumerAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsV1EventsMockSubscriptionConsumerAdded)
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
		it.Event = new(FunctionsV1EventsMockSubscriptionConsumerAdded)
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

func (it *FunctionsV1EventsMockSubscriptionConsumerAddedIterator) Error() error {
	return it.fail
}

func (it *FunctionsV1EventsMockSubscriptionConsumerAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsV1EventsMockSubscriptionConsumerAdded struct {
	SubscriptionId uint64
	Consumer       common.Address
	Raw            types.Log
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) FilterSubscriptionConsumerAdded(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsV1EventsMockSubscriptionConsumerAddedIterator, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsV1EventsMock.contract.FilterLogs(opts, "SubscriptionConsumerAdded", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsV1EventsMockSubscriptionConsumerAddedIterator{contract: _FunctionsV1EventsMock.contract, event: "SubscriptionConsumerAdded", logs: logs, sub: sub}, nil
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) WatchSubscriptionConsumerAdded(opts *bind.WatchOpts, sink chan<- *FunctionsV1EventsMockSubscriptionConsumerAdded, subscriptionId []uint64) (event.Subscription, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsV1EventsMock.contract.WatchLogs(opts, "SubscriptionConsumerAdded", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsV1EventsMockSubscriptionConsumerAdded)
				if err := _FunctionsV1EventsMock.contract.UnpackLog(event, "SubscriptionConsumerAdded", log); err != nil {
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

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) ParseSubscriptionConsumerAdded(log types.Log) (*FunctionsV1EventsMockSubscriptionConsumerAdded, error) {
	event := new(FunctionsV1EventsMockSubscriptionConsumerAdded)
	if err := _FunctionsV1EventsMock.contract.UnpackLog(event, "SubscriptionConsumerAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsV1EventsMockSubscriptionConsumerRemovedIterator struct {
	Event *FunctionsV1EventsMockSubscriptionConsumerRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsV1EventsMockSubscriptionConsumerRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsV1EventsMockSubscriptionConsumerRemoved)
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
		it.Event = new(FunctionsV1EventsMockSubscriptionConsumerRemoved)
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

func (it *FunctionsV1EventsMockSubscriptionConsumerRemovedIterator) Error() error {
	return it.fail
}

func (it *FunctionsV1EventsMockSubscriptionConsumerRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsV1EventsMockSubscriptionConsumerRemoved struct {
	SubscriptionId uint64
	Consumer       common.Address
	Raw            types.Log
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) FilterSubscriptionConsumerRemoved(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsV1EventsMockSubscriptionConsumerRemovedIterator, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsV1EventsMock.contract.FilterLogs(opts, "SubscriptionConsumerRemoved", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsV1EventsMockSubscriptionConsumerRemovedIterator{contract: _FunctionsV1EventsMock.contract, event: "SubscriptionConsumerRemoved", logs: logs, sub: sub}, nil
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) WatchSubscriptionConsumerRemoved(opts *bind.WatchOpts, sink chan<- *FunctionsV1EventsMockSubscriptionConsumerRemoved, subscriptionId []uint64) (event.Subscription, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsV1EventsMock.contract.WatchLogs(opts, "SubscriptionConsumerRemoved", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsV1EventsMockSubscriptionConsumerRemoved)
				if err := _FunctionsV1EventsMock.contract.UnpackLog(event, "SubscriptionConsumerRemoved", log); err != nil {
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

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) ParseSubscriptionConsumerRemoved(log types.Log) (*FunctionsV1EventsMockSubscriptionConsumerRemoved, error) {
	event := new(FunctionsV1EventsMockSubscriptionConsumerRemoved)
	if err := _FunctionsV1EventsMock.contract.UnpackLog(event, "SubscriptionConsumerRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsV1EventsMockSubscriptionCreatedIterator struct {
	Event *FunctionsV1EventsMockSubscriptionCreated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsV1EventsMockSubscriptionCreatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsV1EventsMockSubscriptionCreated)
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
		it.Event = new(FunctionsV1EventsMockSubscriptionCreated)
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

func (it *FunctionsV1EventsMockSubscriptionCreatedIterator) Error() error {
	return it.fail
}

func (it *FunctionsV1EventsMockSubscriptionCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsV1EventsMockSubscriptionCreated struct {
	SubscriptionId uint64
	Owner          common.Address
	Raw            types.Log
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) FilterSubscriptionCreated(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsV1EventsMockSubscriptionCreatedIterator, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsV1EventsMock.contract.FilterLogs(opts, "SubscriptionCreated", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsV1EventsMockSubscriptionCreatedIterator{contract: _FunctionsV1EventsMock.contract, event: "SubscriptionCreated", logs: logs, sub: sub}, nil
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) WatchSubscriptionCreated(opts *bind.WatchOpts, sink chan<- *FunctionsV1EventsMockSubscriptionCreated, subscriptionId []uint64) (event.Subscription, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsV1EventsMock.contract.WatchLogs(opts, "SubscriptionCreated", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsV1EventsMockSubscriptionCreated)
				if err := _FunctionsV1EventsMock.contract.UnpackLog(event, "SubscriptionCreated", log); err != nil {
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

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) ParseSubscriptionCreated(log types.Log) (*FunctionsV1EventsMockSubscriptionCreated, error) {
	event := new(FunctionsV1EventsMockSubscriptionCreated)
	if err := _FunctionsV1EventsMock.contract.UnpackLog(event, "SubscriptionCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsV1EventsMockSubscriptionFundedIterator struct {
	Event *FunctionsV1EventsMockSubscriptionFunded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsV1EventsMockSubscriptionFundedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsV1EventsMockSubscriptionFunded)
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
		it.Event = new(FunctionsV1EventsMockSubscriptionFunded)
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

func (it *FunctionsV1EventsMockSubscriptionFundedIterator) Error() error {
	return it.fail
}

func (it *FunctionsV1EventsMockSubscriptionFundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsV1EventsMockSubscriptionFunded struct {
	SubscriptionId uint64
	OldBalance     *big.Int
	NewBalance     *big.Int
	Raw            types.Log
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) FilterSubscriptionFunded(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsV1EventsMockSubscriptionFundedIterator, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsV1EventsMock.contract.FilterLogs(opts, "SubscriptionFunded", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsV1EventsMockSubscriptionFundedIterator{contract: _FunctionsV1EventsMock.contract, event: "SubscriptionFunded", logs: logs, sub: sub}, nil
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) WatchSubscriptionFunded(opts *bind.WatchOpts, sink chan<- *FunctionsV1EventsMockSubscriptionFunded, subscriptionId []uint64) (event.Subscription, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsV1EventsMock.contract.WatchLogs(opts, "SubscriptionFunded", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsV1EventsMockSubscriptionFunded)
				if err := _FunctionsV1EventsMock.contract.UnpackLog(event, "SubscriptionFunded", log); err != nil {
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

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) ParseSubscriptionFunded(log types.Log) (*FunctionsV1EventsMockSubscriptionFunded, error) {
	event := new(FunctionsV1EventsMockSubscriptionFunded)
	if err := _FunctionsV1EventsMock.contract.UnpackLog(event, "SubscriptionFunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsV1EventsMockSubscriptionOwnerTransferRequestedIterator struct {
	Event *FunctionsV1EventsMockSubscriptionOwnerTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsV1EventsMockSubscriptionOwnerTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsV1EventsMockSubscriptionOwnerTransferRequested)
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
		it.Event = new(FunctionsV1EventsMockSubscriptionOwnerTransferRequested)
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

func (it *FunctionsV1EventsMockSubscriptionOwnerTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *FunctionsV1EventsMockSubscriptionOwnerTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsV1EventsMockSubscriptionOwnerTransferRequested struct {
	SubscriptionId uint64
	From           common.Address
	To             common.Address
	Raw            types.Log
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) FilterSubscriptionOwnerTransferRequested(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsV1EventsMockSubscriptionOwnerTransferRequestedIterator, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsV1EventsMock.contract.FilterLogs(opts, "SubscriptionOwnerTransferRequested", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsV1EventsMockSubscriptionOwnerTransferRequestedIterator{contract: _FunctionsV1EventsMock.contract, event: "SubscriptionOwnerTransferRequested", logs: logs, sub: sub}, nil
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) WatchSubscriptionOwnerTransferRequested(opts *bind.WatchOpts, sink chan<- *FunctionsV1EventsMockSubscriptionOwnerTransferRequested, subscriptionId []uint64) (event.Subscription, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsV1EventsMock.contract.WatchLogs(opts, "SubscriptionOwnerTransferRequested", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsV1EventsMockSubscriptionOwnerTransferRequested)
				if err := _FunctionsV1EventsMock.contract.UnpackLog(event, "SubscriptionOwnerTransferRequested", log); err != nil {
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

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) ParseSubscriptionOwnerTransferRequested(log types.Log) (*FunctionsV1EventsMockSubscriptionOwnerTransferRequested, error) {
	event := new(FunctionsV1EventsMockSubscriptionOwnerTransferRequested)
	if err := _FunctionsV1EventsMock.contract.UnpackLog(event, "SubscriptionOwnerTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsV1EventsMockSubscriptionOwnerTransferredIterator struct {
	Event *FunctionsV1EventsMockSubscriptionOwnerTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsV1EventsMockSubscriptionOwnerTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsV1EventsMockSubscriptionOwnerTransferred)
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
		it.Event = new(FunctionsV1EventsMockSubscriptionOwnerTransferred)
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

func (it *FunctionsV1EventsMockSubscriptionOwnerTransferredIterator) Error() error {
	return it.fail
}

func (it *FunctionsV1EventsMockSubscriptionOwnerTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsV1EventsMockSubscriptionOwnerTransferred struct {
	SubscriptionId uint64
	From           common.Address
	To             common.Address
	Raw            types.Log
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) FilterSubscriptionOwnerTransferred(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsV1EventsMockSubscriptionOwnerTransferredIterator, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsV1EventsMock.contract.FilterLogs(opts, "SubscriptionOwnerTransferred", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsV1EventsMockSubscriptionOwnerTransferredIterator{contract: _FunctionsV1EventsMock.contract, event: "SubscriptionOwnerTransferred", logs: logs, sub: sub}, nil
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) WatchSubscriptionOwnerTransferred(opts *bind.WatchOpts, sink chan<- *FunctionsV1EventsMockSubscriptionOwnerTransferred, subscriptionId []uint64) (event.Subscription, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsV1EventsMock.contract.WatchLogs(opts, "SubscriptionOwnerTransferred", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsV1EventsMockSubscriptionOwnerTransferred)
				if err := _FunctionsV1EventsMock.contract.UnpackLog(event, "SubscriptionOwnerTransferred", log); err != nil {
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

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) ParseSubscriptionOwnerTransferred(log types.Log) (*FunctionsV1EventsMockSubscriptionOwnerTransferred, error) {
	event := new(FunctionsV1EventsMockSubscriptionOwnerTransferred)
	if err := _FunctionsV1EventsMock.contract.UnpackLog(event, "SubscriptionOwnerTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsV1EventsMockUnpausedIterator struct {
	Event *FunctionsV1EventsMockUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsV1EventsMockUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsV1EventsMockUnpaused)
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
		it.Event = new(FunctionsV1EventsMockUnpaused)
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

func (it *FunctionsV1EventsMockUnpausedIterator) Error() error {
	return it.fail
}

func (it *FunctionsV1EventsMockUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsV1EventsMockUnpaused struct {
	Account common.Address
	Raw     types.Log
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) FilterUnpaused(opts *bind.FilterOpts) (*FunctionsV1EventsMockUnpausedIterator, error) {

	logs, sub, err := _FunctionsV1EventsMock.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &FunctionsV1EventsMockUnpausedIterator{contract: _FunctionsV1EventsMock.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *FunctionsV1EventsMockUnpaused) (event.Subscription, error) {

	logs, sub, err := _FunctionsV1EventsMock.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsV1EventsMockUnpaused)
				if err := _FunctionsV1EventsMock.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

func (_FunctionsV1EventsMock *FunctionsV1EventsMockFilterer) ParseUnpaused(log types.Log) (*FunctionsV1EventsMockUnpaused, error) {
	event := new(FunctionsV1EventsMockUnpaused)
	if err := _FunctionsV1EventsMock.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMock) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _FunctionsV1EventsMock.abi.Events["ConfigUpdated"].ID:
		return _FunctionsV1EventsMock.ParseConfigUpdated(log)
	case _FunctionsV1EventsMock.abi.Events["ContractProposed"].ID:
		return _FunctionsV1EventsMock.ParseContractProposed(log)
	case _FunctionsV1EventsMock.abi.Events["ContractUpdated"].ID:
		return _FunctionsV1EventsMock.ParseContractUpdated(log)
	case _FunctionsV1EventsMock.abi.Events["FundsRecovered"].ID:
		return _FunctionsV1EventsMock.ParseFundsRecovered(log)
	case _FunctionsV1EventsMock.abi.Events["OwnershipTransferRequested"].ID:
		return _FunctionsV1EventsMock.ParseOwnershipTransferRequested(log)
	case _FunctionsV1EventsMock.abi.Events["OwnershipTransferred"].ID:
		return _FunctionsV1EventsMock.ParseOwnershipTransferred(log)
	case _FunctionsV1EventsMock.abi.Events["Paused"].ID:
		return _FunctionsV1EventsMock.ParsePaused(log)
	case _FunctionsV1EventsMock.abi.Events["RequestNotProcessed"].ID:
		return _FunctionsV1EventsMock.ParseRequestNotProcessed(log)
	case _FunctionsV1EventsMock.abi.Events["RequestProcessed"].ID:
		return _FunctionsV1EventsMock.ParseRequestProcessed(log)
	case _FunctionsV1EventsMock.abi.Events["RequestStart"].ID:
		return _FunctionsV1EventsMock.ParseRequestStart(log)
	case _FunctionsV1EventsMock.abi.Events["RequestTimedOut"].ID:
		return _FunctionsV1EventsMock.ParseRequestTimedOut(log)
	case _FunctionsV1EventsMock.abi.Events["SubscriptionCanceled"].ID:
		return _FunctionsV1EventsMock.ParseSubscriptionCanceled(log)
	case _FunctionsV1EventsMock.abi.Events["SubscriptionConsumerAdded"].ID:
		return _FunctionsV1EventsMock.ParseSubscriptionConsumerAdded(log)
	case _FunctionsV1EventsMock.abi.Events["SubscriptionConsumerRemoved"].ID:
		return _FunctionsV1EventsMock.ParseSubscriptionConsumerRemoved(log)
	case _FunctionsV1EventsMock.abi.Events["SubscriptionCreated"].ID:
		return _FunctionsV1EventsMock.ParseSubscriptionCreated(log)
	case _FunctionsV1EventsMock.abi.Events["SubscriptionFunded"].ID:
		return _FunctionsV1EventsMock.ParseSubscriptionFunded(log)
	case _FunctionsV1EventsMock.abi.Events["SubscriptionOwnerTransferRequested"].ID:
		return _FunctionsV1EventsMock.ParseSubscriptionOwnerTransferRequested(log)
	case _FunctionsV1EventsMock.abi.Events["SubscriptionOwnerTransferred"].ID:
		return _FunctionsV1EventsMock.ParseSubscriptionOwnerTransferred(log)
	case _FunctionsV1EventsMock.abi.Events["Unpaused"].ID:
		return _FunctionsV1EventsMock.ParseUnpaused(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (FunctionsV1EventsMockConfigUpdated) Topic() common.Hash {
	return common.HexToHash("0x049ce2e6e1420eb4b07b425e90129186833eb346bda40b37d5d921aad482f71c")
}

func (FunctionsV1EventsMockContractProposed) Topic() common.Hash {
	return common.HexToHash("0x8b052f0f4bf82fede7daffea71592b29d5ef86af1f3c7daaa0345dbb2f52f481")
}

func (FunctionsV1EventsMockContractUpdated) Topic() common.Hash {
	return common.HexToHash("0xf8a6175bca1ba37d682089187edc5e20a859989727f10ca6bd9a5bc0de8caf94")
}

func (FunctionsV1EventsMockFundsRecovered) Topic() common.Hash {
	return common.HexToHash("0x59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b436600")
}

func (FunctionsV1EventsMockOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (FunctionsV1EventsMockOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (FunctionsV1EventsMockPaused) Topic() common.Hash {
	return common.HexToHash("0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258")
}

func (FunctionsV1EventsMockRequestNotProcessed) Topic() common.Hash {
	return common.HexToHash("0x1a90e9a50793db2e394cf581e7c522e10c358a81e70acf6b5a0edd620c08dee1")
}

func (FunctionsV1EventsMockRequestProcessed) Topic() common.Hash {
	return common.HexToHash("0x64778f26c70b60a8d7e29e2451b3844302d959448401c0535b768ed88c6b505e")
}

func (FunctionsV1EventsMockRequestStart) Topic() common.Hash {
	return common.HexToHash("0xf67aec45c9a7ede407974a3e0c3a743dffeab99ee3f2d4c9a8144c2ebf2c7ec9")
}

func (FunctionsV1EventsMockRequestTimedOut) Topic() common.Hash {
	return common.HexToHash("0xf1ca1e9147be737b04a2b018a79405f687a97de8dd8a2559bbe62357343af414")
}

func (FunctionsV1EventsMockSubscriptionCanceled) Topic() common.Hash {
	return common.HexToHash("0xe8ed5b475a5b5987aa9165e8731bb78043f39eee32ec5a1169a89e27fcd49815")
}

func (FunctionsV1EventsMockSubscriptionConsumerAdded) Topic() common.Hash {
	return common.HexToHash("0x43dc749a04ac8fb825cbd514f7c0e13f13bc6f2ee66043b76629d51776cff8e0")
}

func (FunctionsV1EventsMockSubscriptionConsumerRemoved) Topic() common.Hash {
	return common.HexToHash("0x182bff9831466789164ca77075fffd84916d35a8180ba73c27e45634549b445b")
}

func (FunctionsV1EventsMockSubscriptionCreated) Topic() common.Hash {
	return common.HexToHash("0x464722b4166576d3dcbba877b999bc35cf911f4eaf434b7eba68fa113951d0bf")
}

func (FunctionsV1EventsMockSubscriptionFunded) Topic() common.Hash {
	return common.HexToHash("0xd39ec07f4e209f627a4c427971473820dc129761ba28de8906bd56f57101d4f8")
}

func (FunctionsV1EventsMockSubscriptionOwnerTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x69436ea6df009049404f564eff6622cd00522b0bd6a89efd9e52a355c4a879be")
}

func (FunctionsV1EventsMockSubscriptionOwnerTransferred) Topic() common.Hash {
	return common.HexToHash("0x6f1dc65165ffffedfd8e507b4a0f1fcfdada045ed11f6c26ba27cedfe87802f0")
}

func (FunctionsV1EventsMockUnpaused) Topic() common.Hash {
	return common.HexToHash("0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa")
}

func (_FunctionsV1EventsMock *FunctionsV1EventsMock) Address() common.Address {
	return _FunctionsV1EventsMock.address
}

type FunctionsV1EventsMockInterface interface {
	EmitConfigUpdated(opts *bind.TransactOpts, param1 FunctionsV1EventsMockConfig) (*types.Transaction, error)

	EmitContractProposed(opts *bind.TransactOpts, proposedContractSetId [32]byte, proposedContractSetFromAddress common.Address, proposedContractSetToAddress common.Address) (*types.Transaction, error)

	EmitContractUpdated(opts *bind.TransactOpts, id [32]byte, from common.Address, to common.Address) (*types.Transaction, error)

	EmitFundsRecovered(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error)

	EmitOwnershipTransferRequested(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error)

	EmitOwnershipTransferred(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error)

	EmitPaused(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error)

	EmitRequestNotProcessed(opts *bind.TransactOpts, requestId [32]byte, coordinator common.Address, transmitter common.Address, resultCode uint8) (*types.Transaction, error)

	EmitRequestProcessed(opts *bind.TransactOpts, requestId [32]byte, subscriptionId uint64, totalCostJuels *big.Int, transmitter common.Address, resultCode uint8, response []byte, err []byte, callbackReturnData []byte) (*types.Transaction, error)

	EmitRequestStart(opts *bind.TransactOpts, requestId [32]byte, donId [32]byte, subscriptionId uint64, subscriptionOwner common.Address, requestingContract common.Address, requestInitiator common.Address, data []byte, dataVersion uint16, callbackGasLimit uint32, estimatedTotalCostJuels *big.Int) (*types.Transaction, error)

	EmitRequestTimedOut(opts *bind.TransactOpts, requestId [32]byte) (*types.Transaction, error)

	EmitSubscriptionCanceled(opts *bind.TransactOpts, subscriptionId uint64, fundsRecipient common.Address, fundsAmount *big.Int) (*types.Transaction, error)

	EmitSubscriptionConsumerAdded(opts *bind.TransactOpts, subscriptionId uint64, consumer common.Address) (*types.Transaction, error)

	EmitSubscriptionConsumerRemoved(opts *bind.TransactOpts, subscriptionId uint64, consumer common.Address) (*types.Transaction, error)

	EmitSubscriptionCreated(opts *bind.TransactOpts, subscriptionId uint64, owner common.Address) (*types.Transaction, error)

	EmitSubscriptionFunded(opts *bind.TransactOpts, subscriptionId uint64, oldBalance *big.Int, newBalance *big.Int) (*types.Transaction, error)

	EmitSubscriptionOwnerTransferRequested(opts *bind.TransactOpts, subscriptionId uint64, from common.Address, to common.Address) (*types.Transaction, error)

	EmitSubscriptionOwnerTransferred(opts *bind.TransactOpts, subscriptionId uint64, from common.Address, to common.Address) (*types.Transaction, error)

	EmitUnpaused(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error)

	FilterConfigUpdated(opts *bind.FilterOpts) (*FunctionsV1EventsMockConfigUpdatedIterator, error)

	WatchConfigUpdated(opts *bind.WatchOpts, sink chan<- *FunctionsV1EventsMockConfigUpdated) (event.Subscription, error)

	ParseConfigUpdated(log types.Log) (*FunctionsV1EventsMockConfigUpdated, error)

	FilterContractProposed(opts *bind.FilterOpts) (*FunctionsV1EventsMockContractProposedIterator, error)

	WatchContractProposed(opts *bind.WatchOpts, sink chan<- *FunctionsV1EventsMockContractProposed) (event.Subscription, error)

	ParseContractProposed(log types.Log) (*FunctionsV1EventsMockContractProposed, error)

	FilterContractUpdated(opts *bind.FilterOpts) (*FunctionsV1EventsMockContractUpdatedIterator, error)

	WatchContractUpdated(opts *bind.WatchOpts, sink chan<- *FunctionsV1EventsMockContractUpdated) (event.Subscription, error)

	ParseContractUpdated(log types.Log) (*FunctionsV1EventsMockContractUpdated, error)

	FilterFundsRecovered(opts *bind.FilterOpts) (*FunctionsV1EventsMockFundsRecoveredIterator, error)

	WatchFundsRecovered(opts *bind.WatchOpts, sink chan<- *FunctionsV1EventsMockFundsRecovered) (event.Subscription, error)

	ParseFundsRecovered(log types.Log) (*FunctionsV1EventsMockFundsRecovered, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsV1EventsMockOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *FunctionsV1EventsMockOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*FunctionsV1EventsMockOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsV1EventsMockOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FunctionsV1EventsMockOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*FunctionsV1EventsMockOwnershipTransferred, error)

	FilterPaused(opts *bind.FilterOpts) (*FunctionsV1EventsMockPausedIterator, error)

	WatchPaused(opts *bind.WatchOpts, sink chan<- *FunctionsV1EventsMockPaused) (event.Subscription, error)

	ParsePaused(log types.Log) (*FunctionsV1EventsMockPaused, error)

	FilterRequestNotProcessed(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsV1EventsMockRequestNotProcessedIterator, error)

	WatchRequestNotProcessed(opts *bind.WatchOpts, sink chan<- *FunctionsV1EventsMockRequestNotProcessed, requestId [][32]byte) (event.Subscription, error)

	ParseRequestNotProcessed(log types.Log) (*FunctionsV1EventsMockRequestNotProcessed, error)

	FilterRequestProcessed(opts *bind.FilterOpts, requestId [][32]byte, subscriptionId []uint64) (*FunctionsV1EventsMockRequestProcessedIterator, error)

	WatchRequestProcessed(opts *bind.WatchOpts, sink chan<- *FunctionsV1EventsMockRequestProcessed, requestId [][32]byte, subscriptionId []uint64) (event.Subscription, error)

	ParseRequestProcessed(log types.Log) (*FunctionsV1EventsMockRequestProcessed, error)

	FilterRequestStart(opts *bind.FilterOpts, requestId [][32]byte, donId [][32]byte, subscriptionId []uint64) (*FunctionsV1EventsMockRequestStartIterator, error)

	WatchRequestStart(opts *bind.WatchOpts, sink chan<- *FunctionsV1EventsMockRequestStart, requestId [][32]byte, donId [][32]byte, subscriptionId []uint64) (event.Subscription, error)

	ParseRequestStart(log types.Log) (*FunctionsV1EventsMockRequestStart, error)

	FilterRequestTimedOut(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsV1EventsMockRequestTimedOutIterator, error)

	WatchRequestTimedOut(opts *bind.WatchOpts, sink chan<- *FunctionsV1EventsMockRequestTimedOut, requestId [][32]byte) (event.Subscription, error)

	ParseRequestTimedOut(log types.Log) (*FunctionsV1EventsMockRequestTimedOut, error)

	FilterSubscriptionCanceled(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsV1EventsMockSubscriptionCanceledIterator, error)

	WatchSubscriptionCanceled(opts *bind.WatchOpts, sink chan<- *FunctionsV1EventsMockSubscriptionCanceled, subscriptionId []uint64) (event.Subscription, error)

	ParseSubscriptionCanceled(log types.Log) (*FunctionsV1EventsMockSubscriptionCanceled, error)

	FilterSubscriptionConsumerAdded(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsV1EventsMockSubscriptionConsumerAddedIterator, error)

	WatchSubscriptionConsumerAdded(opts *bind.WatchOpts, sink chan<- *FunctionsV1EventsMockSubscriptionConsumerAdded, subscriptionId []uint64) (event.Subscription, error)

	ParseSubscriptionConsumerAdded(log types.Log) (*FunctionsV1EventsMockSubscriptionConsumerAdded, error)

	FilterSubscriptionConsumerRemoved(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsV1EventsMockSubscriptionConsumerRemovedIterator, error)

	WatchSubscriptionConsumerRemoved(opts *bind.WatchOpts, sink chan<- *FunctionsV1EventsMockSubscriptionConsumerRemoved, subscriptionId []uint64) (event.Subscription, error)

	ParseSubscriptionConsumerRemoved(log types.Log) (*FunctionsV1EventsMockSubscriptionConsumerRemoved, error)

	FilterSubscriptionCreated(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsV1EventsMockSubscriptionCreatedIterator, error)

	WatchSubscriptionCreated(opts *bind.WatchOpts, sink chan<- *FunctionsV1EventsMockSubscriptionCreated, subscriptionId []uint64) (event.Subscription, error)

	ParseSubscriptionCreated(log types.Log) (*FunctionsV1EventsMockSubscriptionCreated, error)

	FilterSubscriptionFunded(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsV1EventsMockSubscriptionFundedIterator, error)

	WatchSubscriptionFunded(opts *bind.WatchOpts, sink chan<- *FunctionsV1EventsMockSubscriptionFunded, subscriptionId []uint64) (event.Subscription, error)

	ParseSubscriptionFunded(log types.Log) (*FunctionsV1EventsMockSubscriptionFunded, error)

	FilterSubscriptionOwnerTransferRequested(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsV1EventsMockSubscriptionOwnerTransferRequestedIterator, error)

	WatchSubscriptionOwnerTransferRequested(opts *bind.WatchOpts, sink chan<- *FunctionsV1EventsMockSubscriptionOwnerTransferRequested, subscriptionId []uint64) (event.Subscription, error)

	ParseSubscriptionOwnerTransferRequested(log types.Log) (*FunctionsV1EventsMockSubscriptionOwnerTransferRequested, error)

	FilterSubscriptionOwnerTransferred(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsV1EventsMockSubscriptionOwnerTransferredIterator, error)

	WatchSubscriptionOwnerTransferred(opts *bind.WatchOpts, sink chan<- *FunctionsV1EventsMockSubscriptionOwnerTransferred, subscriptionId []uint64) (event.Subscription, error)

	ParseSubscriptionOwnerTransferred(log types.Log) (*FunctionsV1EventsMockSubscriptionOwnerTransferred, error)

	FilterUnpaused(opts *bind.FilterOpts) (*FunctionsV1EventsMockUnpausedIterator, error)

	WatchUnpaused(opts *bind.WatchOpts, sink chan<- *FunctionsV1EventsMockUnpaused) (event.Subscription, error)

	ParseUnpaused(log types.Log) (*FunctionsV1EventsMockUnpaused, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
