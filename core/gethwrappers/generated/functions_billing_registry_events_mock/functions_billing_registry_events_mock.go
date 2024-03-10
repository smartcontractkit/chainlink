// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package functions_billing_registry_events_mock

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

type FunctionsBillingRegistryEventsMockCommitment struct {
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

var FunctionsBillingRegistryEventsMockMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"senders\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"changedBy\",\"type\":\"address\"}],\"name\":\"AuthorizedSendersChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"signerPayment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"transmitterPayment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalCost\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"BillingEnd\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"don\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"donFee\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"registryFee\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"estimatedCost\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structFunctionsBillingRegistryEventsMock.Commitment\",\"name\":\"commitment\",\"type\":\"tuple\"}],\"name\":\"BillingStart\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasOverhead\",\"type\":\"uint32\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"RequestTimedOut\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"SubscriptionCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"SubscriptionCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"senders\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"changedBy\",\"type\":\"address\"}],\"name\":\"emitAuthorizedSendersChanged\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"uint96\",\"name\":\"signerPayment\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"transmitterPayment\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"totalCost\",\"type\":\"uint96\"},{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"emitBillingEnd\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"don\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"donFee\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"registryFee\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"estimatedCost\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"internalType\":\"structFunctionsBillingRegistryEventsMock.Commitment\",\"name\":\"commitment\",\"type\":\"tuple\"}],\"name\":\"emitBillingStart\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"internalType\":\"uint32\",\"name\":\"gasOverhead\",\"type\":\"uint32\"}],\"name\":\"emitConfigSet\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"emitFundsRecovered\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"emitInitialized\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"emitOwnershipTransferRequested\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"emitOwnershipTransferred\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"emitPaused\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"emitRequestTimedOut\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"emitSubscriptionCanceled\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"emitSubscriptionConsumerAdded\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"emitSubscriptionConsumerRemoved\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"emitSubscriptionCreated\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"emitSubscriptionFunded\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"emitSubscriptionOwnerTransferRequested\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"emitSubscriptionOwnerTransferred\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"emitUnpaused\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610eae806100206000396000f3fe608060405234801561001057600080fd5b50600436106101365760003560e01c8063a5257226116100b2578063dde69b3f11610081578063e2cab57b11610066578063e2cab57b1461025a578063e9c6260f1461026d578063f7420bc21461028057600080fd5b8063dde69b3f14610234578063e0f6eff11461024757600080fd5b8063a5257226146101e8578063b019b4e8146101fb578063bef9e1831461020e578063c1d2ad191461022157600080fd5b8063689300ea116101095780637be5c756116100ee5780637be5c756146101af5780637e1b44c0146101c25780639ec3ce4b146101d557600080fd5b8063689300ea14610189578063735bb0821461019c57600080fd5b80632d6d80b31461013b5780633f70afb6146101505780634bf6a80d14610163578063675b924414610176575b600080fd5b61014e61014936600461090f565b610293565b005b61014e61015e366004610ba1565b6102d0565b61014e610171366004610bbd565b61032a565b61014e610184366004610ba1565b61038c565b61014e6101973660046108e5565b6103de565b61014e6101aa366004610b4a565b61042a565b61014e6101bd366004610890565b610487565b61014e6101d03660046109d5565b6104d4565b61014e6101e3366004610890565b610502565b61014e6101f6366004610ba1565b610548565b61014e6102093660046108b2565b61059a565b61014e61021c366004610c6f565b6105f8565b61014e61022f366004610ad4565b61062b565b61014e610242366004610c00565b61069e565b61014e610255366004610bbd565b6106f6565b61014e610268366004610c3c565b61074f565b61014e61027b3660046109ee565b610791565b61014e61028e3660046108b2565b6107c1565b7ff263cfb3e4298332e776194610cf9fdc09ccb3ada8b9aa39764d882e11fbf0a082826040516102c4929190610c92565b60405180910390a15050565b60405173ffffffffffffffffffffffffffffffffffffffff8216815267ffffffffffffffff8316907f464722b4166576d3dcbba877b999bc35cf911f4eaf434b7eba68fa113951d0bf906020015b60405180910390a25050565b6040805173ffffffffffffffffffffffffffffffffffffffff80851682528316602082015267ffffffffffffffff8516917f6f1dc65165ffffedfd8e507b4a0f1fcfdada045ed11f6c26ba27cedfe87802f091015b60405180910390a2505050565b60405173ffffffffffffffffffffffffffffffffffffffff8216815267ffffffffffffffff8316907f43dc749a04ac8fb825cbd514f7c0e13f13bc6f2ee66043b76629d51776cff8e09060200161031e565b6040805173ffffffffffffffffffffffffffffffffffffffff84168152602081018390527f59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b43660091016102c4565b6040805163ffffffff87811682528681166020830152818301869052606082018590528316608082015290517f24d3d934adfef9b9029d6ffa463c07d0139ed47d26ee23506f85ece2879d2bd49181900360a00190a15050505050565b60405173ffffffffffffffffffffffffffffffffffffffff821681527f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258906020015b60405180910390a150565b60405181907ff1ca1e9147be737b04a2b018a79405f687a97de8dd8a2559bbe62357343af41490600090a250565b60405173ffffffffffffffffffffffffffffffffffffffff821681527f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa906020016104c9565b60405173ffffffffffffffffffffffffffffffffffffffff8216815267ffffffffffffffff8316907f182bff9831466789164ca77075fffd84916d35a8180ba73c27e45634549b445b9060200161031e565b8073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b60405160ff821681527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498906020016104c9565b6040805167ffffffffffffffff871681526bffffffffffffffffffffffff868116602083015285811682840152841660608201528215156080820152905187917fc8dc973332de19a5f71b6026983110e9c2e04b0c98b87eb771ccb78607fd114f919081900360a00190a2505050505050565b6040805173ffffffffffffffffffffffffffffffffffffffff841681526020810183905267ffffffffffffffff8516917fe8ed5b475a5b5987aa9165e8731bb78043f39eee32ec5a1169a89e27fcd49815910161037f565b6040805173ffffffffffffffffffffffffffffffffffffffff80851682528316602082015267ffffffffffffffff8516917f69436ea6df009049404f564eff6622cd00522b0bd6a89efd9e52a355c4a879be910161037f565b604080518381526020810183905267ffffffffffffffff8516917fd39ec07f4e209f627a4c427971473820dc129761ba28de8906bd56f57101d4f8910161037f565b817f99f7f4e65b4b9fbabd4e357c47ed3099b36e57ecd3a43e84662f34c207d0ebe48260405161031e9190610d09565b8073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127860405160405180910390a35050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461084357600080fd5b919050565b803563ffffffff8116811461084357600080fd5b803567ffffffffffffffff8116811461084357600080fd5b80356bffffffffffffffffffffffff8116811461084357600080fd5b6000602082840312156108a257600080fd5b6108ab8261081f565b9392505050565b600080604083850312156108c557600080fd5b6108ce8361081f565b91506108dc6020840161081f565b90509250929050565b600080604083850312156108f857600080fd5b6109018361081f565b946020939093013593505050565b6000806040838503121561092257600080fd5b823567ffffffffffffffff8082111561093a57600080fd5b818501915085601f83011261094e57600080fd5b813560208282111561096257610962610e72565b8160051b9250610973818401610e23565b8281528181019085830185870184018b101561098e57600080fd5b600096505b848710156109b8576109a48161081f565b835260019690960195918301918301610993565b5096506109c8905087820161081f565b9450505050509250929050565b6000602082840312156109e757600080fd5b5035919050565b600080828403610140811215610a0357600080fd5b83359250610120807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe083011215610a3957600080fd5b610a41610df9565b9150610a4f6020860161085c565b8252610a5d6040860161081f565b6020830152610a6e60608601610848565b604083015260808501356060830152610a8960a0860161081f565b6080830152610a9a60c08601610874565b60a0830152610aab60e08601610874565b60c0830152610100610abe818701610874565b60e0840152940135938101939093525092909150565b60008060008060008060c08789031215610aed57600080fd5b86359550610afd6020880161085c565b9450610b0b60408801610874565b9350610b1960608801610874565b9250610b2760808801610874565b915060a08701358015158114610b3c57600080fd5b809150509295509295509295565b600080600080600060a08688031215610b6257600080fd5b610b6b86610848565b9450610b7960208701610848565b93506040860135925060608601359150610b9560808701610848565b90509295509295909350565b60008060408385031215610bb457600080fd5b6108ce8361085c565b600080600060608486031215610bd257600080fd5b610bdb8461085c565b9250610be96020850161081f565b9150610bf76040850161081f565b90509250925092565b600080600060608486031215610c1557600080fd5b610c1e8461085c565b9250610c2c6020850161081f565b9150604084013590509250925092565b600080600060608486031215610c5157600080fd5b610c5a8461085c565b95602085013595506040909401359392505050565b600060208284031215610c8157600080fd5b813560ff811681146108ab57600080fd5b604080825283519082018190526000906020906060840190828701845b82811015610ce157815173ffffffffffffffffffffffffffffffffffffffff1684529284019290840190600101610caf565b50505073ffffffffffffffffffffffffffffffffffffffff9490941692019190915250919050565b60006101208201905067ffffffffffffffff835116825273ffffffffffffffffffffffffffffffffffffffff60208401511660208301526040830151610d57604084018263ffffffff169052565b50606083015160608301526080830151610d89608084018273ffffffffffffffffffffffffffffffffffffffff169052565b5060a0830151610da960a08401826bffffffffffffffffffffffff169052565b5060c0830151610dc960c08401826bffffffffffffffffffffffff169052565b5060e0830151610de960e08401826bffffffffffffffffffffffff169052565b5061010092830151919092015290565b604051610120810167ffffffffffffffff81118282101715610e1d57610e1d610e72565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610e6a57610e6a610e72565b604052919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var FunctionsBillingRegistryEventsMockABI = FunctionsBillingRegistryEventsMockMetaData.ABI

var FunctionsBillingRegistryEventsMockBin = FunctionsBillingRegistryEventsMockMetaData.Bin

func DeployFunctionsBillingRegistryEventsMock(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *FunctionsBillingRegistryEventsMock, error) {
	parsed, err := FunctionsBillingRegistryEventsMockMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(FunctionsBillingRegistryEventsMockBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &FunctionsBillingRegistryEventsMock{address: address, abi: *parsed, FunctionsBillingRegistryEventsMockCaller: FunctionsBillingRegistryEventsMockCaller{contract: contract}, FunctionsBillingRegistryEventsMockTransactor: FunctionsBillingRegistryEventsMockTransactor{contract: contract}, FunctionsBillingRegistryEventsMockFilterer: FunctionsBillingRegistryEventsMockFilterer{contract: contract}}, nil
}

type FunctionsBillingRegistryEventsMock struct {
	address common.Address
	abi     abi.ABI
	FunctionsBillingRegistryEventsMockCaller
	FunctionsBillingRegistryEventsMockTransactor
	FunctionsBillingRegistryEventsMockFilterer
}

type FunctionsBillingRegistryEventsMockCaller struct {
	contract *bind.BoundContract
}

type FunctionsBillingRegistryEventsMockTransactor struct {
	contract *bind.BoundContract
}

type FunctionsBillingRegistryEventsMockFilterer struct {
	contract *bind.BoundContract
}

type FunctionsBillingRegistryEventsMockSession struct {
	Contract     *FunctionsBillingRegistryEventsMock
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type FunctionsBillingRegistryEventsMockCallerSession struct {
	Contract *FunctionsBillingRegistryEventsMockCaller
	CallOpts bind.CallOpts
}

type FunctionsBillingRegistryEventsMockTransactorSession struct {
	Contract     *FunctionsBillingRegistryEventsMockTransactor
	TransactOpts bind.TransactOpts
}

type FunctionsBillingRegistryEventsMockRaw struct {
	Contract *FunctionsBillingRegistryEventsMock
}

type FunctionsBillingRegistryEventsMockCallerRaw struct {
	Contract *FunctionsBillingRegistryEventsMockCaller
}

type FunctionsBillingRegistryEventsMockTransactorRaw struct {
	Contract *FunctionsBillingRegistryEventsMockTransactor
}

func NewFunctionsBillingRegistryEventsMock(address common.Address, backend bind.ContractBackend) (*FunctionsBillingRegistryEventsMock, error) {
	abi, err := abi.JSON(strings.NewReader(FunctionsBillingRegistryEventsMockABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindFunctionsBillingRegistryEventsMock(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FunctionsBillingRegistryEventsMock{address: address, abi: abi, FunctionsBillingRegistryEventsMockCaller: FunctionsBillingRegistryEventsMockCaller{contract: contract}, FunctionsBillingRegistryEventsMockTransactor: FunctionsBillingRegistryEventsMockTransactor{contract: contract}, FunctionsBillingRegistryEventsMockFilterer: FunctionsBillingRegistryEventsMockFilterer{contract: contract}}, nil
}

func NewFunctionsBillingRegistryEventsMockCaller(address common.Address, caller bind.ContractCaller) (*FunctionsBillingRegistryEventsMockCaller, error) {
	contract, err := bindFunctionsBillingRegistryEventsMock(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FunctionsBillingRegistryEventsMockCaller{contract: contract}, nil
}

func NewFunctionsBillingRegistryEventsMockTransactor(address common.Address, transactor bind.ContractTransactor) (*FunctionsBillingRegistryEventsMockTransactor, error) {
	contract, err := bindFunctionsBillingRegistryEventsMock(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FunctionsBillingRegistryEventsMockTransactor{contract: contract}, nil
}

func NewFunctionsBillingRegistryEventsMockFilterer(address common.Address, filterer bind.ContractFilterer) (*FunctionsBillingRegistryEventsMockFilterer, error) {
	contract, err := bindFunctionsBillingRegistryEventsMock(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FunctionsBillingRegistryEventsMockFilterer{contract: contract}, nil
}

func bindFunctionsBillingRegistryEventsMock(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FunctionsBillingRegistryEventsMockMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FunctionsBillingRegistryEventsMock.Contract.FunctionsBillingRegistryEventsMockCaller.contract.Call(opts, result, method, params...)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.Contract.FunctionsBillingRegistryEventsMockTransactor.contract.Transfer(opts)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.Contract.FunctionsBillingRegistryEventsMockTransactor.contract.Transact(opts, method, params...)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FunctionsBillingRegistryEventsMock.Contract.contract.Call(opts, result, method, params...)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.Contract.contract.Transfer(opts)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.Contract.contract.Transact(opts, method, params...)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockTransactor) EmitAuthorizedSendersChanged(opts *bind.TransactOpts, senders []common.Address, changedBy common.Address) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.contract.Transact(opts, "emitAuthorizedSendersChanged", senders, changedBy)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockSession) EmitAuthorizedSendersChanged(senders []common.Address, changedBy common.Address) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.Contract.EmitAuthorizedSendersChanged(&_FunctionsBillingRegistryEventsMock.TransactOpts, senders, changedBy)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockTransactorSession) EmitAuthorizedSendersChanged(senders []common.Address, changedBy common.Address) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.Contract.EmitAuthorizedSendersChanged(&_FunctionsBillingRegistryEventsMock.TransactOpts, senders, changedBy)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockTransactor) EmitBillingEnd(opts *bind.TransactOpts, requestId [32]byte, subscriptionId uint64, signerPayment *big.Int, transmitterPayment *big.Int, totalCost *big.Int, success bool) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.contract.Transact(opts, "emitBillingEnd", requestId, subscriptionId, signerPayment, transmitterPayment, totalCost, success)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockSession) EmitBillingEnd(requestId [32]byte, subscriptionId uint64, signerPayment *big.Int, transmitterPayment *big.Int, totalCost *big.Int, success bool) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.Contract.EmitBillingEnd(&_FunctionsBillingRegistryEventsMock.TransactOpts, requestId, subscriptionId, signerPayment, transmitterPayment, totalCost, success)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockTransactorSession) EmitBillingEnd(requestId [32]byte, subscriptionId uint64, signerPayment *big.Int, transmitterPayment *big.Int, totalCost *big.Int, success bool) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.Contract.EmitBillingEnd(&_FunctionsBillingRegistryEventsMock.TransactOpts, requestId, subscriptionId, signerPayment, transmitterPayment, totalCost, success)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockTransactor) EmitBillingStart(opts *bind.TransactOpts, requestId [32]byte, commitment FunctionsBillingRegistryEventsMockCommitment) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.contract.Transact(opts, "emitBillingStart", requestId, commitment)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockSession) EmitBillingStart(requestId [32]byte, commitment FunctionsBillingRegistryEventsMockCommitment) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.Contract.EmitBillingStart(&_FunctionsBillingRegistryEventsMock.TransactOpts, requestId, commitment)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockTransactorSession) EmitBillingStart(requestId [32]byte, commitment FunctionsBillingRegistryEventsMockCommitment) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.Contract.EmitBillingStart(&_FunctionsBillingRegistryEventsMock.TransactOpts, requestId, commitment)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockTransactor) EmitConfigSet(opts *bind.TransactOpts, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation *big.Int, fallbackWeiPerUnitLink *big.Int, gasOverhead uint32) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.contract.Transact(opts, "emitConfigSet", maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, gasOverhead)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockSession) EmitConfigSet(maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation *big.Int, fallbackWeiPerUnitLink *big.Int, gasOverhead uint32) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.Contract.EmitConfigSet(&_FunctionsBillingRegistryEventsMock.TransactOpts, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, gasOverhead)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockTransactorSession) EmitConfigSet(maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation *big.Int, fallbackWeiPerUnitLink *big.Int, gasOverhead uint32) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.Contract.EmitConfigSet(&_FunctionsBillingRegistryEventsMock.TransactOpts, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, gasOverhead)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockTransactor) EmitFundsRecovered(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.contract.Transact(opts, "emitFundsRecovered", to, amount)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockSession) EmitFundsRecovered(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.Contract.EmitFundsRecovered(&_FunctionsBillingRegistryEventsMock.TransactOpts, to, amount)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockTransactorSession) EmitFundsRecovered(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.Contract.EmitFundsRecovered(&_FunctionsBillingRegistryEventsMock.TransactOpts, to, amount)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockTransactor) EmitInitialized(opts *bind.TransactOpts, version uint8) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.contract.Transact(opts, "emitInitialized", version)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockSession) EmitInitialized(version uint8) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.Contract.EmitInitialized(&_FunctionsBillingRegistryEventsMock.TransactOpts, version)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockTransactorSession) EmitInitialized(version uint8) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.Contract.EmitInitialized(&_FunctionsBillingRegistryEventsMock.TransactOpts, version)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockTransactor) EmitOwnershipTransferRequested(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.contract.Transact(opts, "emitOwnershipTransferRequested", from, to)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockSession) EmitOwnershipTransferRequested(from common.Address, to common.Address) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.Contract.EmitOwnershipTransferRequested(&_FunctionsBillingRegistryEventsMock.TransactOpts, from, to)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockTransactorSession) EmitOwnershipTransferRequested(from common.Address, to common.Address) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.Contract.EmitOwnershipTransferRequested(&_FunctionsBillingRegistryEventsMock.TransactOpts, from, to)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockTransactor) EmitOwnershipTransferred(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.contract.Transact(opts, "emitOwnershipTransferred", from, to)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockSession) EmitOwnershipTransferred(from common.Address, to common.Address) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.Contract.EmitOwnershipTransferred(&_FunctionsBillingRegistryEventsMock.TransactOpts, from, to)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockTransactorSession) EmitOwnershipTransferred(from common.Address, to common.Address) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.Contract.EmitOwnershipTransferred(&_FunctionsBillingRegistryEventsMock.TransactOpts, from, to)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockTransactor) EmitPaused(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.contract.Transact(opts, "emitPaused", account)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockSession) EmitPaused(account common.Address) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.Contract.EmitPaused(&_FunctionsBillingRegistryEventsMock.TransactOpts, account)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockTransactorSession) EmitPaused(account common.Address) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.Contract.EmitPaused(&_FunctionsBillingRegistryEventsMock.TransactOpts, account)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockTransactor) EmitRequestTimedOut(opts *bind.TransactOpts, requestId [32]byte) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.contract.Transact(opts, "emitRequestTimedOut", requestId)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockSession) EmitRequestTimedOut(requestId [32]byte) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.Contract.EmitRequestTimedOut(&_FunctionsBillingRegistryEventsMock.TransactOpts, requestId)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockTransactorSession) EmitRequestTimedOut(requestId [32]byte) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.Contract.EmitRequestTimedOut(&_FunctionsBillingRegistryEventsMock.TransactOpts, requestId)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockTransactor) EmitSubscriptionCanceled(opts *bind.TransactOpts, subscriptionId uint64, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.contract.Transact(opts, "emitSubscriptionCanceled", subscriptionId, to, amount)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockSession) EmitSubscriptionCanceled(subscriptionId uint64, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.Contract.EmitSubscriptionCanceled(&_FunctionsBillingRegistryEventsMock.TransactOpts, subscriptionId, to, amount)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockTransactorSession) EmitSubscriptionCanceled(subscriptionId uint64, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.Contract.EmitSubscriptionCanceled(&_FunctionsBillingRegistryEventsMock.TransactOpts, subscriptionId, to, amount)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockTransactor) EmitSubscriptionConsumerAdded(opts *bind.TransactOpts, subscriptionId uint64, consumer common.Address) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.contract.Transact(opts, "emitSubscriptionConsumerAdded", subscriptionId, consumer)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockSession) EmitSubscriptionConsumerAdded(subscriptionId uint64, consumer common.Address) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.Contract.EmitSubscriptionConsumerAdded(&_FunctionsBillingRegistryEventsMock.TransactOpts, subscriptionId, consumer)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockTransactorSession) EmitSubscriptionConsumerAdded(subscriptionId uint64, consumer common.Address) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.Contract.EmitSubscriptionConsumerAdded(&_FunctionsBillingRegistryEventsMock.TransactOpts, subscriptionId, consumer)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockTransactor) EmitSubscriptionConsumerRemoved(opts *bind.TransactOpts, subscriptionId uint64, consumer common.Address) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.contract.Transact(opts, "emitSubscriptionConsumerRemoved", subscriptionId, consumer)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockSession) EmitSubscriptionConsumerRemoved(subscriptionId uint64, consumer common.Address) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.Contract.EmitSubscriptionConsumerRemoved(&_FunctionsBillingRegistryEventsMock.TransactOpts, subscriptionId, consumer)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockTransactorSession) EmitSubscriptionConsumerRemoved(subscriptionId uint64, consumer common.Address) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.Contract.EmitSubscriptionConsumerRemoved(&_FunctionsBillingRegistryEventsMock.TransactOpts, subscriptionId, consumer)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockTransactor) EmitSubscriptionCreated(opts *bind.TransactOpts, subscriptionId uint64, owner common.Address) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.contract.Transact(opts, "emitSubscriptionCreated", subscriptionId, owner)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockSession) EmitSubscriptionCreated(subscriptionId uint64, owner common.Address) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.Contract.EmitSubscriptionCreated(&_FunctionsBillingRegistryEventsMock.TransactOpts, subscriptionId, owner)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockTransactorSession) EmitSubscriptionCreated(subscriptionId uint64, owner common.Address) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.Contract.EmitSubscriptionCreated(&_FunctionsBillingRegistryEventsMock.TransactOpts, subscriptionId, owner)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockTransactor) EmitSubscriptionFunded(opts *bind.TransactOpts, subscriptionId uint64, oldBalance *big.Int, newBalance *big.Int) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.contract.Transact(opts, "emitSubscriptionFunded", subscriptionId, oldBalance, newBalance)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockSession) EmitSubscriptionFunded(subscriptionId uint64, oldBalance *big.Int, newBalance *big.Int) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.Contract.EmitSubscriptionFunded(&_FunctionsBillingRegistryEventsMock.TransactOpts, subscriptionId, oldBalance, newBalance)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockTransactorSession) EmitSubscriptionFunded(subscriptionId uint64, oldBalance *big.Int, newBalance *big.Int) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.Contract.EmitSubscriptionFunded(&_FunctionsBillingRegistryEventsMock.TransactOpts, subscriptionId, oldBalance, newBalance)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockTransactor) EmitSubscriptionOwnerTransferRequested(opts *bind.TransactOpts, subscriptionId uint64, from common.Address, to common.Address) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.contract.Transact(opts, "emitSubscriptionOwnerTransferRequested", subscriptionId, from, to)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockSession) EmitSubscriptionOwnerTransferRequested(subscriptionId uint64, from common.Address, to common.Address) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.Contract.EmitSubscriptionOwnerTransferRequested(&_FunctionsBillingRegistryEventsMock.TransactOpts, subscriptionId, from, to)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockTransactorSession) EmitSubscriptionOwnerTransferRequested(subscriptionId uint64, from common.Address, to common.Address) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.Contract.EmitSubscriptionOwnerTransferRequested(&_FunctionsBillingRegistryEventsMock.TransactOpts, subscriptionId, from, to)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockTransactor) EmitSubscriptionOwnerTransferred(opts *bind.TransactOpts, subscriptionId uint64, from common.Address, to common.Address) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.contract.Transact(opts, "emitSubscriptionOwnerTransferred", subscriptionId, from, to)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockSession) EmitSubscriptionOwnerTransferred(subscriptionId uint64, from common.Address, to common.Address) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.Contract.EmitSubscriptionOwnerTransferred(&_FunctionsBillingRegistryEventsMock.TransactOpts, subscriptionId, from, to)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockTransactorSession) EmitSubscriptionOwnerTransferred(subscriptionId uint64, from common.Address, to common.Address) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.Contract.EmitSubscriptionOwnerTransferred(&_FunctionsBillingRegistryEventsMock.TransactOpts, subscriptionId, from, to)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockTransactor) EmitUnpaused(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.contract.Transact(opts, "emitUnpaused", account)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockSession) EmitUnpaused(account common.Address) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.Contract.EmitUnpaused(&_FunctionsBillingRegistryEventsMock.TransactOpts, account)
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockTransactorSession) EmitUnpaused(account common.Address) (*types.Transaction, error) {
	return _FunctionsBillingRegistryEventsMock.Contract.EmitUnpaused(&_FunctionsBillingRegistryEventsMock.TransactOpts, account)
}

type FunctionsBillingRegistryEventsMockAuthorizedSendersChangedIterator struct {
	Event *FunctionsBillingRegistryEventsMockAuthorizedSendersChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsBillingRegistryEventsMockAuthorizedSendersChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsBillingRegistryEventsMockAuthorizedSendersChanged)
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
		it.Event = new(FunctionsBillingRegistryEventsMockAuthorizedSendersChanged)
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

func (it *FunctionsBillingRegistryEventsMockAuthorizedSendersChangedIterator) Error() error {
	return it.fail
}

func (it *FunctionsBillingRegistryEventsMockAuthorizedSendersChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsBillingRegistryEventsMockAuthorizedSendersChanged struct {
	Senders   []common.Address
	ChangedBy common.Address
	Raw       types.Log
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) FilterAuthorizedSendersChanged(opts *bind.FilterOpts) (*FunctionsBillingRegistryEventsMockAuthorizedSendersChangedIterator, error) {

	logs, sub, err := _FunctionsBillingRegistryEventsMock.contract.FilterLogs(opts, "AuthorizedSendersChanged")
	if err != nil {
		return nil, err
	}
	return &FunctionsBillingRegistryEventsMockAuthorizedSendersChangedIterator{contract: _FunctionsBillingRegistryEventsMock.contract, event: "AuthorizedSendersChanged", logs: logs, sub: sub}, nil
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) WatchAuthorizedSendersChanged(opts *bind.WatchOpts, sink chan<- *FunctionsBillingRegistryEventsMockAuthorizedSendersChanged) (event.Subscription, error) {

	logs, sub, err := _FunctionsBillingRegistryEventsMock.contract.WatchLogs(opts, "AuthorizedSendersChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsBillingRegistryEventsMockAuthorizedSendersChanged)
				if err := _FunctionsBillingRegistryEventsMock.contract.UnpackLog(event, "AuthorizedSendersChanged", log); err != nil {
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

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) ParseAuthorizedSendersChanged(log types.Log) (*FunctionsBillingRegistryEventsMockAuthorizedSendersChanged, error) {
	event := new(FunctionsBillingRegistryEventsMockAuthorizedSendersChanged)
	if err := _FunctionsBillingRegistryEventsMock.contract.UnpackLog(event, "AuthorizedSendersChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsBillingRegistryEventsMockBillingEndIterator struct {
	Event *FunctionsBillingRegistryEventsMockBillingEnd

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsBillingRegistryEventsMockBillingEndIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsBillingRegistryEventsMockBillingEnd)
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
		it.Event = new(FunctionsBillingRegistryEventsMockBillingEnd)
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

func (it *FunctionsBillingRegistryEventsMockBillingEndIterator) Error() error {
	return it.fail
}

func (it *FunctionsBillingRegistryEventsMockBillingEndIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsBillingRegistryEventsMockBillingEnd struct {
	RequestId          [32]byte
	SubscriptionId     uint64
	SignerPayment      *big.Int
	TransmitterPayment *big.Int
	TotalCost          *big.Int
	Success            bool
	Raw                types.Log
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) FilterBillingEnd(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsBillingRegistryEventsMockBillingEndIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsBillingRegistryEventsMock.contract.FilterLogs(opts, "BillingEnd", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsBillingRegistryEventsMockBillingEndIterator{contract: _FunctionsBillingRegistryEventsMock.contract, event: "BillingEnd", logs: logs, sub: sub}, nil
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) WatchBillingEnd(opts *bind.WatchOpts, sink chan<- *FunctionsBillingRegistryEventsMockBillingEnd, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsBillingRegistryEventsMock.contract.WatchLogs(opts, "BillingEnd", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsBillingRegistryEventsMockBillingEnd)
				if err := _FunctionsBillingRegistryEventsMock.contract.UnpackLog(event, "BillingEnd", log); err != nil {
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

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) ParseBillingEnd(log types.Log) (*FunctionsBillingRegistryEventsMockBillingEnd, error) {
	event := new(FunctionsBillingRegistryEventsMockBillingEnd)
	if err := _FunctionsBillingRegistryEventsMock.contract.UnpackLog(event, "BillingEnd", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsBillingRegistryEventsMockBillingStartIterator struct {
	Event *FunctionsBillingRegistryEventsMockBillingStart

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsBillingRegistryEventsMockBillingStartIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsBillingRegistryEventsMockBillingStart)
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
		it.Event = new(FunctionsBillingRegistryEventsMockBillingStart)
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

func (it *FunctionsBillingRegistryEventsMockBillingStartIterator) Error() error {
	return it.fail
}

func (it *FunctionsBillingRegistryEventsMockBillingStartIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsBillingRegistryEventsMockBillingStart struct {
	RequestId  [32]byte
	Commitment FunctionsBillingRegistryEventsMockCommitment
	Raw        types.Log
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) FilterBillingStart(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsBillingRegistryEventsMockBillingStartIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsBillingRegistryEventsMock.contract.FilterLogs(opts, "BillingStart", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsBillingRegistryEventsMockBillingStartIterator{contract: _FunctionsBillingRegistryEventsMock.contract, event: "BillingStart", logs: logs, sub: sub}, nil
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) WatchBillingStart(opts *bind.WatchOpts, sink chan<- *FunctionsBillingRegistryEventsMockBillingStart, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsBillingRegistryEventsMock.contract.WatchLogs(opts, "BillingStart", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsBillingRegistryEventsMockBillingStart)
				if err := _FunctionsBillingRegistryEventsMock.contract.UnpackLog(event, "BillingStart", log); err != nil {
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

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) ParseBillingStart(log types.Log) (*FunctionsBillingRegistryEventsMockBillingStart, error) {
	event := new(FunctionsBillingRegistryEventsMockBillingStart)
	if err := _FunctionsBillingRegistryEventsMock.contract.UnpackLog(event, "BillingStart", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsBillingRegistryEventsMockConfigSetIterator struct {
	Event *FunctionsBillingRegistryEventsMockConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsBillingRegistryEventsMockConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsBillingRegistryEventsMockConfigSet)
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
		it.Event = new(FunctionsBillingRegistryEventsMockConfigSet)
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

func (it *FunctionsBillingRegistryEventsMockConfigSetIterator) Error() error {
	return it.fail
}

func (it *FunctionsBillingRegistryEventsMockConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsBillingRegistryEventsMockConfigSet struct {
	MaxGasLimit                uint32
	StalenessSeconds           uint32
	GasAfterPaymentCalculation *big.Int
	FallbackWeiPerUnitLink     *big.Int
	GasOverhead                uint32
	Raw                        types.Log
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) FilterConfigSet(opts *bind.FilterOpts) (*FunctionsBillingRegistryEventsMockConfigSetIterator, error) {

	logs, sub, err := _FunctionsBillingRegistryEventsMock.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &FunctionsBillingRegistryEventsMockConfigSetIterator{contract: _FunctionsBillingRegistryEventsMock.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *FunctionsBillingRegistryEventsMockConfigSet) (event.Subscription, error) {

	logs, sub, err := _FunctionsBillingRegistryEventsMock.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsBillingRegistryEventsMockConfigSet)
				if err := _FunctionsBillingRegistryEventsMock.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) ParseConfigSet(log types.Log) (*FunctionsBillingRegistryEventsMockConfigSet, error) {
	event := new(FunctionsBillingRegistryEventsMockConfigSet)
	if err := _FunctionsBillingRegistryEventsMock.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsBillingRegistryEventsMockFundsRecoveredIterator struct {
	Event *FunctionsBillingRegistryEventsMockFundsRecovered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsBillingRegistryEventsMockFundsRecoveredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsBillingRegistryEventsMockFundsRecovered)
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
		it.Event = new(FunctionsBillingRegistryEventsMockFundsRecovered)
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

func (it *FunctionsBillingRegistryEventsMockFundsRecoveredIterator) Error() error {
	return it.fail
}

func (it *FunctionsBillingRegistryEventsMockFundsRecoveredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsBillingRegistryEventsMockFundsRecovered struct {
	To     common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) FilterFundsRecovered(opts *bind.FilterOpts) (*FunctionsBillingRegistryEventsMockFundsRecoveredIterator, error) {

	logs, sub, err := _FunctionsBillingRegistryEventsMock.contract.FilterLogs(opts, "FundsRecovered")
	if err != nil {
		return nil, err
	}
	return &FunctionsBillingRegistryEventsMockFundsRecoveredIterator{contract: _FunctionsBillingRegistryEventsMock.contract, event: "FundsRecovered", logs: logs, sub: sub}, nil
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) WatchFundsRecovered(opts *bind.WatchOpts, sink chan<- *FunctionsBillingRegistryEventsMockFundsRecovered) (event.Subscription, error) {

	logs, sub, err := _FunctionsBillingRegistryEventsMock.contract.WatchLogs(opts, "FundsRecovered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsBillingRegistryEventsMockFundsRecovered)
				if err := _FunctionsBillingRegistryEventsMock.contract.UnpackLog(event, "FundsRecovered", log); err != nil {
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

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) ParseFundsRecovered(log types.Log) (*FunctionsBillingRegistryEventsMockFundsRecovered, error) {
	event := new(FunctionsBillingRegistryEventsMockFundsRecovered)
	if err := _FunctionsBillingRegistryEventsMock.contract.UnpackLog(event, "FundsRecovered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsBillingRegistryEventsMockInitializedIterator struct {
	Event *FunctionsBillingRegistryEventsMockInitialized

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsBillingRegistryEventsMockInitializedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsBillingRegistryEventsMockInitialized)
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
		it.Event = new(FunctionsBillingRegistryEventsMockInitialized)
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

func (it *FunctionsBillingRegistryEventsMockInitializedIterator) Error() error {
	return it.fail
}

func (it *FunctionsBillingRegistryEventsMockInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsBillingRegistryEventsMockInitialized struct {
	Version uint8
	Raw     types.Log
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) FilterInitialized(opts *bind.FilterOpts) (*FunctionsBillingRegistryEventsMockInitializedIterator, error) {

	logs, sub, err := _FunctionsBillingRegistryEventsMock.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &FunctionsBillingRegistryEventsMockInitializedIterator{contract: _FunctionsBillingRegistryEventsMock.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *FunctionsBillingRegistryEventsMockInitialized) (event.Subscription, error) {

	logs, sub, err := _FunctionsBillingRegistryEventsMock.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsBillingRegistryEventsMockInitialized)
				if err := _FunctionsBillingRegistryEventsMock.contract.UnpackLog(event, "Initialized", log); err != nil {
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

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) ParseInitialized(log types.Log) (*FunctionsBillingRegistryEventsMockInitialized, error) {
	event := new(FunctionsBillingRegistryEventsMockInitialized)
	if err := _FunctionsBillingRegistryEventsMock.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsBillingRegistryEventsMockOwnershipTransferRequestedIterator struct {
	Event *FunctionsBillingRegistryEventsMockOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsBillingRegistryEventsMockOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsBillingRegistryEventsMockOwnershipTransferRequested)
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
		it.Event = new(FunctionsBillingRegistryEventsMockOwnershipTransferRequested)
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

func (it *FunctionsBillingRegistryEventsMockOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *FunctionsBillingRegistryEventsMockOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsBillingRegistryEventsMockOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsBillingRegistryEventsMockOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsBillingRegistryEventsMock.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsBillingRegistryEventsMockOwnershipTransferRequestedIterator{contract: _FunctionsBillingRegistryEventsMock.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *FunctionsBillingRegistryEventsMockOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsBillingRegistryEventsMock.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsBillingRegistryEventsMockOwnershipTransferRequested)
				if err := _FunctionsBillingRegistryEventsMock.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) ParseOwnershipTransferRequested(log types.Log) (*FunctionsBillingRegistryEventsMockOwnershipTransferRequested, error) {
	event := new(FunctionsBillingRegistryEventsMockOwnershipTransferRequested)
	if err := _FunctionsBillingRegistryEventsMock.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsBillingRegistryEventsMockOwnershipTransferredIterator struct {
	Event *FunctionsBillingRegistryEventsMockOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsBillingRegistryEventsMockOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsBillingRegistryEventsMockOwnershipTransferred)
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
		it.Event = new(FunctionsBillingRegistryEventsMockOwnershipTransferred)
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

func (it *FunctionsBillingRegistryEventsMockOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *FunctionsBillingRegistryEventsMockOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsBillingRegistryEventsMockOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsBillingRegistryEventsMockOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsBillingRegistryEventsMock.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsBillingRegistryEventsMockOwnershipTransferredIterator{contract: _FunctionsBillingRegistryEventsMock.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FunctionsBillingRegistryEventsMockOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsBillingRegistryEventsMock.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsBillingRegistryEventsMockOwnershipTransferred)
				if err := _FunctionsBillingRegistryEventsMock.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) ParseOwnershipTransferred(log types.Log) (*FunctionsBillingRegistryEventsMockOwnershipTransferred, error) {
	event := new(FunctionsBillingRegistryEventsMockOwnershipTransferred)
	if err := _FunctionsBillingRegistryEventsMock.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsBillingRegistryEventsMockPausedIterator struct {
	Event *FunctionsBillingRegistryEventsMockPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsBillingRegistryEventsMockPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsBillingRegistryEventsMockPaused)
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
		it.Event = new(FunctionsBillingRegistryEventsMockPaused)
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

func (it *FunctionsBillingRegistryEventsMockPausedIterator) Error() error {
	return it.fail
}

func (it *FunctionsBillingRegistryEventsMockPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsBillingRegistryEventsMockPaused struct {
	Account common.Address
	Raw     types.Log
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) FilterPaused(opts *bind.FilterOpts) (*FunctionsBillingRegistryEventsMockPausedIterator, error) {

	logs, sub, err := _FunctionsBillingRegistryEventsMock.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &FunctionsBillingRegistryEventsMockPausedIterator{contract: _FunctionsBillingRegistryEventsMock.contract, event: "Paused", logs: logs, sub: sub}, nil
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *FunctionsBillingRegistryEventsMockPaused) (event.Subscription, error) {

	logs, sub, err := _FunctionsBillingRegistryEventsMock.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsBillingRegistryEventsMockPaused)
				if err := _FunctionsBillingRegistryEventsMock.contract.UnpackLog(event, "Paused", log); err != nil {
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

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) ParsePaused(log types.Log) (*FunctionsBillingRegistryEventsMockPaused, error) {
	event := new(FunctionsBillingRegistryEventsMockPaused)
	if err := _FunctionsBillingRegistryEventsMock.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsBillingRegistryEventsMockRequestTimedOutIterator struct {
	Event *FunctionsBillingRegistryEventsMockRequestTimedOut

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsBillingRegistryEventsMockRequestTimedOutIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsBillingRegistryEventsMockRequestTimedOut)
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
		it.Event = new(FunctionsBillingRegistryEventsMockRequestTimedOut)
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

func (it *FunctionsBillingRegistryEventsMockRequestTimedOutIterator) Error() error {
	return it.fail
}

func (it *FunctionsBillingRegistryEventsMockRequestTimedOutIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsBillingRegistryEventsMockRequestTimedOut struct {
	RequestId [32]byte
	Raw       types.Log
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) FilterRequestTimedOut(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsBillingRegistryEventsMockRequestTimedOutIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsBillingRegistryEventsMock.contract.FilterLogs(opts, "RequestTimedOut", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsBillingRegistryEventsMockRequestTimedOutIterator{contract: _FunctionsBillingRegistryEventsMock.contract, event: "RequestTimedOut", logs: logs, sub: sub}, nil
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) WatchRequestTimedOut(opts *bind.WatchOpts, sink chan<- *FunctionsBillingRegistryEventsMockRequestTimedOut, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsBillingRegistryEventsMock.contract.WatchLogs(opts, "RequestTimedOut", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsBillingRegistryEventsMockRequestTimedOut)
				if err := _FunctionsBillingRegistryEventsMock.contract.UnpackLog(event, "RequestTimedOut", log); err != nil {
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

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) ParseRequestTimedOut(log types.Log) (*FunctionsBillingRegistryEventsMockRequestTimedOut, error) {
	event := new(FunctionsBillingRegistryEventsMockRequestTimedOut)
	if err := _FunctionsBillingRegistryEventsMock.contract.UnpackLog(event, "RequestTimedOut", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsBillingRegistryEventsMockSubscriptionCanceledIterator struct {
	Event *FunctionsBillingRegistryEventsMockSubscriptionCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsBillingRegistryEventsMockSubscriptionCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsBillingRegistryEventsMockSubscriptionCanceled)
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
		it.Event = new(FunctionsBillingRegistryEventsMockSubscriptionCanceled)
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

func (it *FunctionsBillingRegistryEventsMockSubscriptionCanceledIterator) Error() error {
	return it.fail
}

func (it *FunctionsBillingRegistryEventsMockSubscriptionCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsBillingRegistryEventsMockSubscriptionCanceled struct {
	SubscriptionId uint64
	To             common.Address
	Amount         *big.Int
	Raw            types.Log
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) FilterSubscriptionCanceled(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsBillingRegistryEventsMockSubscriptionCanceledIterator, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsBillingRegistryEventsMock.contract.FilterLogs(opts, "SubscriptionCanceled", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsBillingRegistryEventsMockSubscriptionCanceledIterator{contract: _FunctionsBillingRegistryEventsMock.contract, event: "SubscriptionCanceled", logs: logs, sub: sub}, nil
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) WatchSubscriptionCanceled(opts *bind.WatchOpts, sink chan<- *FunctionsBillingRegistryEventsMockSubscriptionCanceled, subscriptionId []uint64) (event.Subscription, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsBillingRegistryEventsMock.contract.WatchLogs(opts, "SubscriptionCanceled", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsBillingRegistryEventsMockSubscriptionCanceled)
				if err := _FunctionsBillingRegistryEventsMock.contract.UnpackLog(event, "SubscriptionCanceled", log); err != nil {
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

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) ParseSubscriptionCanceled(log types.Log) (*FunctionsBillingRegistryEventsMockSubscriptionCanceled, error) {
	event := new(FunctionsBillingRegistryEventsMockSubscriptionCanceled)
	if err := _FunctionsBillingRegistryEventsMock.contract.UnpackLog(event, "SubscriptionCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsBillingRegistryEventsMockSubscriptionConsumerAddedIterator struct {
	Event *FunctionsBillingRegistryEventsMockSubscriptionConsumerAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsBillingRegistryEventsMockSubscriptionConsumerAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsBillingRegistryEventsMockSubscriptionConsumerAdded)
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
		it.Event = new(FunctionsBillingRegistryEventsMockSubscriptionConsumerAdded)
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

func (it *FunctionsBillingRegistryEventsMockSubscriptionConsumerAddedIterator) Error() error {
	return it.fail
}

func (it *FunctionsBillingRegistryEventsMockSubscriptionConsumerAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsBillingRegistryEventsMockSubscriptionConsumerAdded struct {
	SubscriptionId uint64
	Consumer       common.Address
	Raw            types.Log
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) FilterSubscriptionConsumerAdded(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsBillingRegistryEventsMockSubscriptionConsumerAddedIterator, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsBillingRegistryEventsMock.contract.FilterLogs(opts, "SubscriptionConsumerAdded", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsBillingRegistryEventsMockSubscriptionConsumerAddedIterator{contract: _FunctionsBillingRegistryEventsMock.contract, event: "SubscriptionConsumerAdded", logs: logs, sub: sub}, nil
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) WatchSubscriptionConsumerAdded(opts *bind.WatchOpts, sink chan<- *FunctionsBillingRegistryEventsMockSubscriptionConsumerAdded, subscriptionId []uint64) (event.Subscription, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsBillingRegistryEventsMock.contract.WatchLogs(opts, "SubscriptionConsumerAdded", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsBillingRegistryEventsMockSubscriptionConsumerAdded)
				if err := _FunctionsBillingRegistryEventsMock.contract.UnpackLog(event, "SubscriptionConsumerAdded", log); err != nil {
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

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) ParseSubscriptionConsumerAdded(log types.Log) (*FunctionsBillingRegistryEventsMockSubscriptionConsumerAdded, error) {
	event := new(FunctionsBillingRegistryEventsMockSubscriptionConsumerAdded)
	if err := _FunctionsBillingRegistryEventsMock.contract.UnpackLog(event, "SubscriptionConsumerAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsBillingRegistryEventsMockSubscriptionConsumerRemovedIterator struct {
	Event *FunctionsBillingRegistryEventsMockSubscriptionConsumerRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsBillingRegistryEventsMockSubscriptionConsumerRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsBillingRegistryEventsMockSubscriptionConsumerRemoved)
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
		it.Event = new(FunctionsBillingRegistryEventsMockSubscriptionConsumerRemoved)
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

func (it *FunctionsBillingRegistryEventsMockSubscriptionConsumerRemovedIterator) Error() error {
	return it.fail
}

func (it *FunctionsBillingRegistryEventsMockSubscriptionConsumerRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsBillingRegistryEventsMockSubscriptionConsumerRemoved struct {
	SubscriptionId uint64
	Consumer       common.Address
	Raw            types.Log
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) FilterSubscriptionConsumerRemoved(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsBillingRegistryEventsMockSubscriptionConsumerRemovedIterator, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsBillingRegistryEventsMock.contract.FilterLogs(opts, "SubscriptionConsumerRemoved", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsBillingRegistryEventsMockSubscriptionConsumerRemovedIterator{contract: _FunctionsBillingRegistryEventsMock.contract, event: "SubscriptionConsumerRemoved", logs: logs, sub: sub}, nil
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) WatchSubscriptionConsumerRemoved(opts *bind.WatchOpts, sink chan<- *FunctionsBillingRegistryEventsMockSubscriptionConsumerRemoved, subscriptionId []uint64) (event.Subscription, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsBillingRegistryEventsMock.contract.WatchLogs(opts, "SubscriptionConsumerRemoved", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsBillingRegistryEventsMockSubscriptionConsumerRemoved)
				if err := _FunctionsBillingRegistryEventsMock.contract.UnpackLog(event, "SubscriptionConsumerRemoved", log); err != nil {
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

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) ParseSubscriptionConsumerRemoved(log types.Log) (*FunctionsBillingRegistryEventsMockSubscriptionConsumerRemoved, error) {
	event := new(FunctionsBillingRegistryEventsMockSubscriptionConsumerRemoved)
	if err := _FunctionsBillingRegistryEventsMock.contract.UnpackLog(event, "SubscriptionConsumerRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsBillingRegistryEventsMockSubscriptionCreatedIterator struct {
	Event *FunctionsBillingRegistryEventsMockSubscriptionCreated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsBillingRegistryEventsMockSubscriptionCreatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsBillingRegistryEventsMockSubscriptionCreated)
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
		it.Event = new(FunctionsBillingRegistryEventsMockSubscriptionCreated)
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

func (it *FunctionsBillingRegistryEventsMockSubscriptionCreatedIterator) Error() error {
	return it.fail
}

func (it *FunctionsBillingRegistryEventsMockSubscriptionCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsBillingRegistryEventsMockSubscriptionCreated struct {
	SubscriptionId uint64
	Owner          common.Address
	Raw            types.Log
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) FilterSubscriptionCreated(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsBillingRegistryEventsMockSubscriptionCreatedIterator, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsBillingRegistryEventsMock.contract.FilterLogs(opts, "SubscriptionCreated", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsBillingRegistryEventsMockSubscriptionCreatedIterator{contract: _FunctionsBillingRegistryEventsMock.contract, event: "SubscriptionCreated", logs: logs, sub: sub}, nil
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) WatchSubscriptionCreated(opts *bind.WatchOpts, sink chan<- *FunctionsBillingRegistryEventsMockSubscriptionCreated, subscriptionId []uint64) (event.Subscription, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsBillingRegistryEventsMock.contract.WatchLogs(opts, "SubscriptionCreated", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsBillingRegistryEventsMockSubscriptionCreated)
				if err := _FunctionsBillingRegistryEventsMock.contract.UnpackLog(event, "SubscriptionCreated", log); err != nil {
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

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) ParseSubscriptionCreated(log types.Log) (*FunctionsBillingRegistryEventsMockSubscriptionCreated, error) {
	event := new(FunctionsBillingRegistryEventsMockSubscriptionCreated)
	if err := _FunctionsBillingRegistryEventsMock.contract.UnpackLog(event, "SubscriptionCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsBillingRegistryEventsMockSubscriptionFundedIterator struct {
	Event *FunctionsBillingRegistryEventsMockSubscriptionFunded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsBillingRegistryEventsMockSubscriptionFundedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsBillingRegistryEventsMockSubscriptionFunded)
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
		it.Event = new(FunctionsBillingRegistryEventsMockSubscriptionFunded)
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

func (it *FunctionsBillingRegistryEventsMockSubscriptionFundedIterator) Error() error {
	return it.fail
}

func (it *FunctionsBillingRegistryEventsMockSubscriptionFundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsBillingRegistryEventsMockSubscriptionFunded struct {
	SubscriptionId uint64
	OldBalance     *big.Int
	NewBalance     *big.Int
	Raw            types.Log
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) FilterSubscriptionFunded(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsBillingRegistryEventsMockSubscriptionFundedIterator, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsBillingRegistryEventsMock.contract.FilterLogs(opts, "SubscriptionFunded", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsBillingRegistryEventsMockSubscriptionFundedIterator{contract: _FunctionsBillingRegistryEventsMock.contract, event: "SubscriptionFunded", logs: logs, sub: sub}, nil
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) WatchSubscriptionFunded(opts *bind.WatchOpts, sink chan<- *FunctionsBillingRegistryEventsMockSubscriptionFunded, subscriptionId []uint64) (event.Subscription, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsBillingRegistryEventsMock.contract.WatchLogs(opts, "SubscriptionFunded", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsBillingRegistryEventsMockSubscriptionFunded)
				if err := _FunctionsBillingRegistryEventsMock.contract.UnpackLog(event, "SubscriptionFunded", log); err != nil {
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

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) ParseSubscriptionFunded(log types.Log) (*FunctionsBillingRegistryEventsMockSubscriptionFunded, error) {
	event := new(FunctionsBillingRegistryEventsMockSubscriptionFunded)
	if err := _FunctionsBillingRegistryEventsMock.contract.UnpackLog(event, "SubscriptionFunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsBillingRegistryEventsMockSubscriptionOwnerTransferRequestedIterator struct {
	Event *FunctionsBillingRegistryEventsMockSubscriptionOwnerTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsBillingRegistryEventsMockSubscriptionOwnerTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsBillingRegistryEventsMockSubscriptionOwnerTransferRequested)
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
		it.Event = new(FunctionsBillingRegistryEventsMockSubscriptionOwnerTransferRequested)
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

func (it *FunctionsBillingRegistryEventsMockSubscriptionOwnerTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *FunctionsBillingRegistryEventsMockSubscriptionOwnerTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsBillingRegistryEventsMockSubscriptionOwnerTransferRequested struct {
	SubscriptionId uint64
	From           common.Address
	To             common.Address
	Raw            types.Log
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) FilterSubscriptionOwnerTransferRequested(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsBillingRegistryEventsMockSubscriptionOwnerTransferRequestedIterator, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsBillingRegistryEventsMock.contract.FilterLogs(opts, "SubscriptionOwnerTransferRequested", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsBillingRegistryEventsMockSubscriptionOwnerTransferRequestedIterator{contract: _FunctionsBillingRegistryEventsMock.contract, event: "SubscriptionOwnerTransferRequested", logs: logs, sub: sub}, nil
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) WatchSubscriptionOwnerTransferRequested(opts *bind.WatchOpts, sink chan<- *FunctionsBillingRegistryEventsMockSubscriptionOwnerTransferRequested, subscriptionId []uint64) (event.Subscription, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsBillingRegistryEventsMock.contract.WatchLogs(opts, "SubscriptionOwnerTransferRequested", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsBillingRegistryEventsMockSubscriptionOwnerTransferRequested)
				if err := _FunctionsBillingRegistryEventsMock.contract.UnpackLog(event, "SubscriptionOwnerTransferRequested", log); err != nil {
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

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) ParseSubscriptionOwnerTransferRequested(log types.Log) (*FunctionsBillingRegistryEventsMockSubscriptionOwnerTransferRequested, error) {
	event := new(FunctionsBillingRegistryEventsMockSubscriptionOwnerTransferRequested)
	if err := _FunctionsBillingRegistryEventsMock.contract.UnpackLog(event, "SubscriptionOwnerTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsBillingRegistryEventsMockSubscriptionOwnerTransferredIterator struct {
	Event *FunctionsBillingRegistryEventsMockSubscriptionOwnerTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsBillingRegistryEventsMockSubscriptionOwnerTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsBillingRegistryEventsMockSubscriptionOwnerTransferred)
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
		it.Event = new(FunctionsBillingRegistryEventsMockSubscriptionOwnerTransferred)
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

func (it *FunctionsBillingRegistryEventsMockSubscriptionOwnerTransferredIterator) Error() error {
	return it.fail
}

func (it *FunctionsBillingRegistryEventsMockSubscriptionOwnerTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsBillingRegistryEventsMockSubscriptionOwnerTransferred struct {
	SubscriptionId uint64
	From           common.Address
	To             common.Address
	Raw            types.Log
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) FilterSubscriptionOwnerTransferred(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsBillingRegistryEventsMockSubscriptionOwnerTransferredIterator, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsBillingRegistryEventsMock.contract.FilterLogs(opts, "SubscriptionOwnerTransferred", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsBillingRegistryEventsMockSubscriptionOwnerTransferredIterator{contract: _FunctionsBillingRegistryEventsMock.contract, event: "SubscriptionOwnerTransferred", logs: logs, sub: sub}, nil
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) WatchSubscriptionOwnerTransferred(opts *bind.WatchOpts, sink chan<- *FunctionsBillingRegistryEventsMockSubscriptionOwnerTransferred, subscriptionId []uint64) (event.Subscription, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _FunctionsBillingRegistryEventsMock.contract.WatchLogs(opts, "SubscriptionOwnerTransferred", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsBillingRegistryEventsMockSubscriptionOwnerTransferred)
				if err := _FunctionsBillingRegistryEventsMock.contract.UnpackLog(event, "SubscriptionOwnerTransferred", log); err != nil {
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

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) ParseSubscriptionOwnerTransferred(log types.Log) (*FunctionsBillingRegistryEventsMockSubscriptionOwnerTransferred, error) {
	event := new(FunctionsBillingRegistryEventsMockSubscriptionOwnerTransferred)
	if err := _FunctionsBillingRegistryEventsMock.contract.UnpackLog(event, "SubscriptionOwnerTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsBillingRegistryEventsMockUnpausedIterator struct {
	Event *FunctionsBillingRegistryEventsMockUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsBillingRegistryEventsMockUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsBillingRegistryEventsMockUnpaused)
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
		it.Event = new(FunctionsBillingRegistryEventsMockUnpaused)
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

func (it *FunctionsBillingRegistryEventsMockUnpausedIterator) Error() error {
	return it.fail
}

func (it *FunctionsBillingRegistryEventsMockUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsBillingRegistryEventsMockUnpaused struct {
	Account common.Address
	Raw     types.Log
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) FilterUnpaused(opts *bind.FilterOpts) (*FunctionsBillingRegistryEventsMockUnpausedIterator, error) {

	logs, sub, err := _FunctionsBillingRegistryEventsMock.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &FunctionsBillingRegistryEventsMockUnpausedIterator{contract: _FunctionsBillingRegistryEventsMock.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *FunctionsBillingRegistryEventsMockUnpaused) (event.Subscription, error) {

	logs, sub, err := _FunctionsBillingRegistryEventsMock.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsBillingRegistryEventsMockUnpaused)
				if err := _FunctionsBillingRegistryEventsMock.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMockFilterer) ParseUnpaused(log types.Log) (*FunctionsBillingRegistryEventsMockUnpaused, error) {
	event := new(FunctionsBillingRegistryEventsMockUnpaused)
	if err := _FunctionsBillingRegistryEventsMock.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMock) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _FunctionsBillingRegistryEventsMock.abi.Events["AuthorizedSendersChanged"].ID:
		return _FunctionsBillingRegistryEventsMock.ParseAuthorizedSendersChanged(log)
	case _FunctionsBillingRegistryEventsMock.abi.Events["BillingEnd"].ID:
		return _FunctionsBillingRegistryEventsMock.ParseBillingEnd(log)
	case _FunctionsBillingRegistryEventsMock.abi.Events["BillingStart"].ID:
		return _FunctionsBillingRegistryEventsMock.ParseBillingStart(log)
	case _FunctionsBillingRegistryEventsMock.abi.Events["ConfigSet"].ID:
		return _FunctionsBillingRegistryEventsMock.ParseConfigSet(log)
	case _FunctionsBillingRegistryEventsMock.abi.Events["FundsRecovered"].ID:
		return _FunctionsBillingRegistryEventsMock.ParseFundsRecovered(log)
	case _FunctionsBillingRegistryEventsMock.abi.Events["Initialized"].ID:
		return _FunctionsBillingRegistryEventsMock.ParseInitialized(log)
	case _FunctionsBillingRegistryEventsMock.abi.Events["OwnershipTransferRequested"].ID:
		return _FunctionsBillingRegistryEventsMock.ParseOwnershipTransferRequested(log)
	case _FunctionsBillingRegistryEventsMock.abi.Events["OwnershipTransferred"].ID:
		return _FunctionsBillingRegistryEventsMock.ParseOwnershipTransferred(log)
	case _FunctionsBillingRegistryEventsMock.abi.Events["Paused"].ID:
		return _FunctionsBillingRegistryEventsMock.ParsePaused(log)
	case _FunctionsBillingRegistryEventsMock.abi.Events["RequestTimedOut"].ID:
		return _FunctionsBillingRegistryEventsMock.ParseRequestTimedOut(log)
	case _FunctionsBillingRegistryEventsMock.abi.Events["SubscriptionCanceled"].ID:
		return _FunctionsBillingRegistryEventsMock.ParseSubscriptionCanceled(log)
	case _FunctionsBillingRegistryEventsMock.abi.Events["SubscriptionConsumerAdded"].ID:
		return _FunctionsBillingRegistryEventsMock.ParseSubscriptionConsumerAdded(log)
	case _FunctionsBillingRegistryEventsMock.abi.Events["SubscriptionConsumerRemoved"].ID:
		return _FunctionsBillingRegistryEventsMock.ParseSubscriptionConsumerRemoved(log)
	case _FunctionsBillingRegistryEventsMock.abi.Events["SubscriptionCreated"].ID:
		return _FunctionsBillingRegistryEventsMock.ParseSubscriptionCreated(log)
	case _FunctionsBillingRegistryEventsMock.abi.Events["SubscriptionFunded"].ID:
		return _FunctionsBillingRegistryEventsMock.ParseSubscriptionFunded(log)
	case _FunctionsBillingRegistryEventsMock.abi.Events["SubscriptionOwnerTransferRequested"].ID:
		return _FunctionsBillingRegistryEventsMock.ParseSubscriptionOwnerTransferRequested(log)
	case _FunctionsBillingRegistryEventsMock.abi.Events["SubscriptionOwnerTransferred"].ID:
		return _FunctionsBillingRegistryEventsMock.ParseSubscriptionOwnerTransferred(log)
	case _FunctionsBillingRegistryEventsMock.abi.Events["Unpaused"].ID:
		return _FunctionsBillingRegistryEventsMock.ParseUnpaused(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (FunctionsBillingRegistryEventsMockAuthorizedSendersChanged) Topic() common.Hash {
	return common.HexToHash("0xf263cfb3e4298332e776194610cf9fdc09ccb3ada8b9aa39764d882e11fbf0a0")
}

func (FunctionsBillingRegistryEventsMockBillingEnd) Topic() common.Hash {
	return common.HexToHash("0xc8dc973332de19a5f71b6026983110e9c2e04b0c98b87eb771ccb78607fd114f")
}

func (FunctionsBillingRegistryEventsMockBillingStart) Topic() common.Hash {
	return common.HexToHash("0x99f7f4e65b4b9fbabd4e357c47ed3099b36e57ecd3a43e84662f34c207d0ebe4")
}

func (FunctionsBillingRegistryEventsMockConfigSet) Topic() common.Hash {
	return common.HexToHash("0x24d3d934adfef9b9029d6ffa463c07d0139ed47d26ee23506f85ece2879d2bd4")
}

func (FunctionsBillingRegistryEventsMockFundsRecovered) Topic() common.Hash {
	return common.HexToHash("0x59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b436600")
}

func (FunctionsBillingRegistryEventsMockInitialized) Topic() common.Hash {
	return common.HexToHash("0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498")
}

func (FunctionsBillingRegistryEventsMockOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (FunctionsBillingRegistryEventsMockOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (FunctionsBillingRegistryEventsMockPaused) Topic() common.Hash {
	return common.HexToHash("0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258")
}

func (FunctionsBillingRegistryEventsMockRequestTimedOut) Topic() common.Hash {
	return common.HexToHash("0xf1ca1e9147be737b04a2b018a79405f687a97de8dd8a2559bbe62357343af414")
}

func (FunctionsBillingRegistryEventsMockSubscriptionCanceled) Topic() common.Hash {
	return common.HexToHash("0xe8ed5b475a5b5987aa9165e8731bb78043f39eee32ec5a1169a89e27fcd49815")
}

func (FunctionsBillingRegistryEventsMockSubscriptionConsumerAdded) Topic() common.Hash {
	return common.HexToHash("0x43dc749a04ac8fb825cbd514f7c0e13f13bc6f2ee66043b76629d51776cff8e0")
}

func (FunctionsBillingRegistryEventsMockSubscriptionConsumerRemoved) Topic() common.Hash {
	return common.HexToHash("0x182bff9831466789164ca77075fffd84916d35a8180ba73c27e45634549b445b")
}

func (FunctionsBillingRegistryEventsMockSubscriptionCreated) Topic() common.Hash {
	return common.HexToHash("0x464722b4166576d3dcbba877b999bc35cf911f4eaf434b7eba68fa113951d0bf")
}

func (FunctionsBillingRegistryEventsMockSubscriptionFunded) Topic() common.Hash {
	return common.HexToHash("0xd39ec07f4e209f627a4c427971473820dc129761ba28de8906bd56f57101d4f8")
}

func (FunctionsBillingRegistryEventsMockSubscriptionOwnerTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x69436ea6df009049404f564eff6622cd00522b0bd6a89efd9e52a355c4a879be")
}

func (FunctionsBillingRegistryEventsMockSubscriptionOwnerTransferred) Topic() common.Hash {
	return common.HexToHash("0x6f1dc65165ffffedfd8e507b4a0f1fcfdada045ed11f6c26ba27cedfe87802f0")
}

func (FunctionsBillingRegistryEventsMockUnpaused) Topic() common.Hash {
	return common.HexToHash("0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa")
}

func (_FunctionsBillingRegistryEventsMock *FunctionsBillingRegistryEventsMock) Address() common.Address {
	return _FunctionsBillingRegistryEventsMock.address
}

type FunctionsBillingRegistryEventsMockInterface interface {
	EmitAuthorizedSendersChanged(opts *bind.TransactOpts, senders []common.Address, changedBy common.Address) (*types.Transaction, error)

	EmitBillingEnd(opts *bind.TransactOpts, requestId [32]byte, subscriptionId uint64, signerPayment *big.Int, transmitterPayment *big.Int, totalCost *big.Int, success bool) (*types.Transaction, error)

	EmitBillingStart(opts *bind.TransactOpts, requestId [32]byte, commitment FunctionsBillingRegistryEventsMockCommitment) (*types.Transaction, error)

	EmitConfigSet(opts *bind.TransactOpts, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation *big.Int, fallbackWeiPerUnitLink *big.Int, gasOverhead uint32) (*types.Transaction, error)

	EmitFundsRecovered(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error)

	EmitInitialized(opts *bind.TransactOpts, version uint8) (*types.Transaction, error)

	EmitOwnershipTransferRequested(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error)

	EmitOwnershipTransferred(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error)

	EmitPaused(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error)

	EmitRequestTimedOut(opts *bind.TransactOpts, requestId [32]byte) (*types.Transaction, error)

	EmitSubscriptionCanceled(opts *bind.TransactOpts, subscriptionId uint64, to common.Address, amount *big.Int) (*types.Transaction, error)

	EmitSubscriptionConsumerAdded(opts *bind.TransactOpts, subscriptionId uint64, consumer common.Address) (*types.Transaction, error)

	EmitSubscriptionConsumerRemoved(opts *bind.TransactOpts, subscriptionId uint64, consumer common.Address) (*types.Transaction, error)

	EmitSubscriptionCreated(opts *bind.TransactOpts, subscriptionId uint64, owner common.Address) (*types.Transaction, error)

	EmitSubscriptionFunded(opts *bind.TransactOpts, subscriptionId uint64, oldBalance *big.Int, newBalance *big.Int) (*types.Transaction, error)

	EmitSubscriptionOwnerTransferRequested(opts *bind.TransactOpts, subscriptionId uint64, from common.Address, to common.Address) (*types.Transaction, error)

	EmitSubscriptionOwnerTransferred(opts *bind.TransactOpts, subscriptionId uint64, from common.Address, to common.Address) (*types.Transaction, error)

	EmitUnpaused(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error)

	FilterAuthorizedSendersChanged(opts *bind.FilterOpts) (*FunctionsBillingRegistryEventsMockAuthorizedSendersChangedIterator, error)

	WatchAuthorizedSendersChanged(opts *bind.WatchOpts, sink chan<- *FunctionsBillingRegistryEventsMockAuthorizedSendersChanged) (event.Subscription, error)

	ParseAuthorizedSendersChanged(log types.Log) (*FunctionsBillingRegistryEventsMockAuthorizedSendersChanged, error)

	FilterBillingEnd(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsBillingRegistryEventsMockBillingEndIterator, error)

	WatchBillingEnd(opts *bind.WatchOpts, sink chan<- *FunctionsBillingRegistryEventsMockBillingEnd, requestId [][32]byte) (event.Subscription, error)

	ParseBillingEnd(log types.Log) (*FunctionsBillingRegistryEventsMockBillingEnd, error)

	FilterBillingStart(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsBillingRegistryEventsMockBillingStartIterator, error)

	WatchBillingStart(opts *bind.WatchOpts, sink chan<- *FunctionsBillingRegistryEventsMockBillingStart, requestId [][32]byte) (event.Subscription, error)

	ParseBillingStart(log types.Log) (*FunctionsBillingRegistryEventsMockBillingStart, error)

	FilterConfigSet(opts *bind.FilterOpts) (*FunctionsBillingRegistryEventsMockConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *FunctionsBillingRegistryEventsMockConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*FunctionsBillingRegistryEventsMockConfigSet, error)

	FilterFundsRecovered(opts *bind.FilterOpts) (*FunctionsBillingRegistryEventsMockFundsRecoveredIterator, error)

	WatchFundsRecovered(opts *bind.WatchOpts, sink chan<- *FunctionsBillingRegistryEventsMockFundsRecovered) (event.Subscription, error)

	ParseFundsRecovered(log types.Log) (*FunctionsBillingRegistryEventsMockFundsRecovered, error)

	FilterInitialized(opts *bind.FilterOpts) (*FunctionsBillingRegistryEventsMockInitializedIterator, error)

	WatchInitialized(opts *bind.WatchOpts, sink chan<- *FunctionsBillingRegistryEventsMockInitialized) (event.Subscription, error)

	ParseInitialized(log types.Log) (*FunctionsBillingRegistryEventsMockInitialized, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsBillingRegistryEventsMockOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *FunctionsBillingRegistryEventsMockOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*FunctionsBillingRegistryEventsMockOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsBillingRegistryEventsMockOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FunctionsBillingRegistryEventsMockOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*FunctionsBillingRegistryEventsMockOwnershipTransferred, error)

	FilterPaused(opts *bind.FilterOpts) (*FunctionsBillingRegistryEventsMockPausedIterator, error)

	WatchPaused(opts *bind.WatchOpts, sink chan<- *FunctionsBillingRegistryEventsMockPaused) (event.Subscription, error)

	ParsePaused(log types.Log) (*FunctionsBillingRegistryEventsMockPaused, error)

	FilterRequestTimedOut(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsBillingRegistryEventsMockRequestTimedOutIterator, error)

	WatchRequestTimedOut(opts *bind.WatchOpts, sink chan<- *FunctionsBillingRegistryEventsMockRequestTimedOut, requestId [][32]byte) (event.Subscription, error)

	ParseRequestTimedOut(log types.Log) (*FunctionsBillingRegistryEventsMockRequestTimedOut, error)

	FilterSubscriptionCanceled(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsBillingRegistryEventsMockSubscriptionCanceledIterator, error)

	WatchSubscriptionCanceled(opts *bind.WatchOpts, sink chan<- *FunctionsBillingRegistryEventsMockSubscriptionCanceled, subscriptionId []uint64) (event.Subscription, error)

	ParseSubscriptionCanceled(log types.Log) (*FunctionsBillingRegistryEventsMockSubscriptionCanceled, error)

	FilterSubscriptionConsumerAdded(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsBillingRegistryEventsMockSubscriptionConsumerAddedIterator, error)

	WatchSubscriptionConsumerAdded(opts *bind.WatchOpts, sink chan<- *FunctionsBillingRegistryEventsMockSubscriptionConsumerAdded, subscriptionId []uint64) (event.Subscription, error)

	ParseSubscriptionConsumerAdded(log types.Log) (*FunctionsBillingRegistryEventsMockSubscriptionConsumerAdded, error)

	FilterSubscriptionConsumerRemoved(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsBillingRegistryEventsMockSubscriptionConsumerRemovedIterator, error)

	WatchSubscriptionConsumerRemoved(opts *bind.WatchOpts, sink chan<- *FunctionsBillingRegistryEventsMockSubscriptionConsumerRemoved, subscriptionId []uint64) (event.Subscription, error)

	ParseSubscriptionConsumerRemoved(log types.Log) (*FunctionsBillingRegistryEventsMockSubscriptionConsumerRemoved, error)

	FilterSubscriptionCreated(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsBillingRegistryEventsMockSubscriptionCreatedIterator, error)

	WatchSubscriptionCreated(opts *bind.WatchOpts, sink chan<- *FunctionsBillingRegistryEventsMockSubscriptionCreated, subscriptionId []uint64) (event.Subscription, error)

	ParseSubscriptionCreated(log types.Log) (*FunctionsBillingRegistryEventsMockSubscriptionCreated, error)

	FilterSubscriptionFunded(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsBillingRegistryEventsMockSubscriptionFundedIterator, error)

	WatchSubscriptionFunded(opts *bind.WatchOpts, sink chan<- *FunctionsBillingRegistryEventsMockSubscriptionFunded, subscriptionId []uint64) (event.Subscription, error)

	ParseSubscriptionFunded(log types.Log) (*FunctionsBillingRegistryEventsMockSubscriptionFunded, error)

	FilterSubscriptionOwnerTransferRequested(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsBillingRegistryEventsMockSubscriptionOwnerTransferRequestedIterator, error)

	WatchSubscriptionOwnerTransferRequested(opts *bind.WatchOpts, sink chan<- *FunctionsBillingRegistryEventsMockSubscriptionOwnerTransferRequested, subscriptionId []uint64) (event.Subscription, error)

	ParseSubscriptionOwnerTransferRequested(log types.Log) (*FunctionsBillingRegistryEventsMockSubscriptionOwnerTransferRequested, error)

	FilterSubscriptionOwnerTransferred(opts *bind.FilterOpts, subscriptionId []uint64) (*FunctionsBillingRegistryEventsMockSubscriptionOwnerTransferredIterator, error)

	WatchSubscriptionOwnerTransferred(opts *bind.WatchOpts, sink chan<- *FunctionsBillingRegistryEventsMockSubscriptionOwnerTransferred, subscriptionId []uint64) (event.Subscription, error)

	ParseSubscriptionOwnerTransferred(log types.Log) (*FunctionsBillingRegistryEventsMockSubscriptionOwnerTransferred, error)

	FilterUnpaused(opts *bind.FilterOpts) (*FunctionsBillingRegistryEventsMockUnpausedIterator, error)

	WatchUnpaused(opts *bind.WatchOpts, sink chan<- *FunctionsBillingRegistryEventsMockUnpaused) (event.Subscription, error)

	ParseUnpaused(log types.Log) (*FunctionsBillingRegistryEventsMockUnpaused, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
