// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package functions_oracle_events_mock

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

var FunctionsOracleEventsMockMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"AuthorizedSendersActive\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"senders\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"changedBy\",\"type\":\"address\"}],\"name\":\"AuthorizedSendersChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"AuthorizedSendersDeactive\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"InvalidRequestID\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requestingContract\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requestInitiator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"subscriptionOwner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"OracleRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"OracleResponse\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"ResponseTransmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"reason\",\"type\":\"string\"}],\"name\":\"UserCallbackError\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"lowLevelData\",\"type\":\"bytes\"}],\"name\":\"UserCallbackRawError\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"emitAuthorizedSendersActive\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"senders\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"changedBy\",\"type\":\"address\"}],\"name\":\"emitAuthorizedSendersChanged\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"emitAuthorizedSendersDeactive\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"emitConfigSet\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"emitInitialized\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"emitInvalidRequestID\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"requestingContract\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"requestInitiator\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"subscriptionOwner\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"emitOracleRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"emitOracleResponse\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"emitOwnershipTransferRequested\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"emitOwnershipTransferred\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"emitResponseTransmitted\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"name\":\"emitTransmitted\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"reason\",\"type\":\"string\"}],\"name\":\"emitUserCallbackError\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"lowLevelData\",\"type\":\"bytes\"}],\"name\":\"emitUserCallbackRawError\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610cd9806100206000396000f3fe608060405234801561001057600080fd5b50600436106100ea5760003560e01c8063b019b4e81161008c578063ddd5603c11610066578063ddd5603c146101af578063df4a9fe0146101c2578063e055cff0146101d5578063f7420bc2146101e857600080fd5b8063b019b4e814610176578063bef9e18314610189578063c9d3d0811461019c57600080fd5b80632d6d80b3116100c85780632d6d80b31461012a5780636446fe921461013d5780638eb62eb6146101505780639784f15a1461016357600080fd5b806317472dac146100ef57806327a88d59146101045780632a7f477b14610117575b600080fd5b6101026100fd36600461072f565b6101fb565b005b6101026101123660046107c2565b610248565b610102610125366004610886565b610276565b61010261013836600461077d565b6102b2565b61010261014b3660046107fe565b6102ef565b61010261015e3660046107c2565b610337565b61010261017136600461072f565b610365565b61010261018436600461074a565b6103ab565b610102610197366004610a3b565b610409565b6101026101aa3660046108cd565b61043c565b6101026101bd36600461091e565b61046c565b6101026101d03660046107db565b6104a7565b6101026101e3366004610941565b6104ef565b6101026101f636600461074a565b610541565b60405173ffffffffffffffffffffffffffffffffffffffff821681527fae51766a982895b0c444fc99fc1a560762b464d709e6c78376c85617f7eeb5ce906020015b60405180910390a150565b60405181907f9e9bc7616d42c2835d05ae617e508454e63b30b934be8aa932ebc125e0e58a6490600090a250565b817fe0b838ffe6ee22a0d3acf19a85db6a41b34a1ab739e2d6c759a2e42d95bdccb2826040516102a69190610ba6565b60405180910390a25050565b7ff263cfb3e4298332e776194610cf9fdc09ccb3ada8b9aa39764d882e11fbf0a082826040516102e3929190610b6e565b60405180910390a15050565b857fa1ec73989d79578cd6f67d4f593ac3e0a4d1020e5c0164db52108d7ff785406c8686868686604051610327959493929190610b12565b60405180910390a2505050505050565b60405181907fa1c120e327c9ad8b075793878c88d59b8934b97ae37117faa3bb21616237f7be90600090a250565b60405173ffffffffffffffffffffffffffffffffffffffff821681527fea3828816a323b8d7ff49d755efd105e7719166d6c76fad97a28eee5eccc3d9a9060200161023d565b8073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b60405160ff821681527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200161023d565b817fb2931868c372fe17a25643458add467d60ec5c51125a99b7309f41f5bcd2da6c826040516102a69190610ba6565b6040805183815263ffffffff831660208201527fb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a6291016102e3565b60405173ffffffffffffffffffffffffffffffffffffffff8216815282907fdc941eddab34a6109ab77798299c6b1f035b125fd6f774d266ecbf9541d630a6906020016102a6565b7f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e0589898989898989898960405161052e99989796959493929190610bb9565b60405180910390a1505050505050505050565b8073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127860405160405180910390a35050565b600067ffffffffffffffff8311156105b9576105b9610c9d565b6105ea60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f86011601610c4e565b90508281528383830111156105fe57600080fd5b828260208301376000602084830101529392505050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461063957600080fd5b919050565b600082601f83011261064f57600080fd5b8135602067ffffffffffffffff82111561066b5761066b610c9d565b8160051b61067a828201610c4e565b83815282810190868401838801850189101561069557600080fd5b600093505b858410156106bf576106ab81610615565b83526001939093019291840191840161069a565b50979650505050505050565b600082601f8301126106dc57600080fd5b6106eb8383356020850161059f565b9392505050565b803563ffffffff8116811461063957600080fd5b803567ffffffffffffffff8116811461063957600080fd5b803560ff8116811461063957600080fd5b60006020828403121561074157600080fd5b6106eb82610615565b6000806040838503121561075d57600080fd5b61076683610615565b915061077460208401610615565b90509250929050565b6000806040838503121561079057600080fd5b823567ffffffffffffffff8111156107a757600080fd5b6107b38582860161063e565b92505061077460208401610615565b6000602082840312156107d457600080fd5b5035919050565b600080604083850312156107ee57600080fd5b8235915061077460208401610615565b60008060008060008060c0878903121561081757600080fd5b8635955061082760208801610615565b945061083560408801610615565b935061084360608801610706565b925061085160808801610615565b915060a087013567ffffffffffffffff81111561086d57600080fd5b61087989828a016106cb565b9150509295509295509295565b6000806040838503121561089957600080fd5b82359150602083013567ffffffffffffffff8111156108b757600080fd5b6108c3858286016106cb565b9150509250929050565b600080604083850312156108e057600080fd5b82359150602083013567ffffffffffffffff8111156108fe57600080fd5b8301601f8101851361090f57600080fd5b6108c38582356020840161059f565b6000806040838503121561093157600080fd5b82359150610774602084016106f2565b60008060008060008060008060006101208a8c03121561096057600080fd5b6109698a6106f2565b985060208a0135975061097e60408b01610706565b965060608a013567ffffffffffffffff8082111561099b57600080fd5b6109a78d838e0161063e565b975060808c01359150808211156109bd57600080fd5b6109c98d838e0161063e565b96506109d760a08d0161071e565b955060c08c01359150808211156109ed57600080fd5b6109f98d838e016106cb565b9450610a0760e08d01610706565b93506101008c0135915080821115610a1e57600080fd5b50610a2b8c828d016106cb565b9150509295985092959850929598565b600060208284031215610a4d57600080fd5b6106eb8261071e565b600081518084526020808501945080840160005b83811015610a9c57815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101610a6a565b509495945050505050565b6000815180845260005b81811015610acd57602081850181015186830182015201610ab1565b81811115610adf576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b600073ffffffffffffffffffffffffffffffffffffffff8088168352808716602084015267ffffffffffffffff8616604084015280851660608401525060a06080830152610b6360a0830184610aa7565b979650505050505050565b604081526000610b816040830185610a56565b905073ffffffffffffffffffffffffffffffffffffffff831660208301529392505050565b6020815260006106eb6020830184610aa7565b600061012063ffffffff8c1683528a602084015267ffffffffffffffff808b166040850152816060850152610bf08285018b610a56565b91508382036080850152610c04828a610a56565b915060ff881660a085015283820360c0850152610c218288610aa7565b90861660e08501528381036101008501529050610c3e8185610aa7565b9c9b505050505050505050505050565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610c9557610c95610c9d565b604052919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var FunctionsOracleEventsMockABI = FunctionsOracleEventsMockMetaData.ABI

var FunctionsOracleEventsMockBin = FunctionsOracleEventsMockMetaData.Bin

func DeployFunctionsOracleEventsMock(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *FunctionsOracleEventsMock, error) {
	parsed, err := FunctionsOracleEventsMockMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(FunctionsOracleEventsMockBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &FunctionsOracleEventsMock{address: address, abi: *parsed, FunctionsOracleEventsMockCaller: FunctionsOracleEventsMockCaller{contract: contract}, FunctionsOracleEventsMockTransactor: FunctionsOracleEventsMockTransactor{contract: contract}, FunctionsOracleEventsMockFilterer: FunctionsOracleEventsMockFilterer{contract: contract}}, nil
}

type FunctionsOracleEventsMock struct {
	address common.Address
	abi     abi.ABI
	FunctionsOracleEventsMockCaller
	FunctionsOracleEventsMockTransactor
	FunctionsOracleEventsMockFilterer
}

type FunctionsOracleEventsMockCaller struct {
	contract *bind.BoundContract
}

type FunctionsOracleEventsMockTransactor struct {
	contract *bind.BoundContract
}

type FunctionsOracleEventsMockFilterer struct {
	contract *bind.BoundContract
}

type FunctionsOracleEventsMockSession struct {
	Contract     *FunctionsOracleEventsMock
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type FunctionsOracleEventsMockCallerSession struct {
	Contract *FunctionsOracleEventsMockCaller
	CallOpts bind.CallOpts
}

type FunctionsOracleEventsMockTransactorSession struct {
	Contract     *FunctionsOracleEventsMockTransactor
	TransactOpts bind.TransactOpts
}

type FunctionsOracleEventsMockRaw struct {
	Contract *FunctionsOracleEventsMock
}

type FunctionsOracleEventsMockCallerRaw struct {
	Contract *FunctionsOracleEventsMockCaller
}

type FunctionsOracleEventsMockTransactorRaw struct {
	Contract *FunctionsOracleEventsMockTransactor
}

func NewFunctionsOracleEventsMock(address common.Address, backend bind.ContractBackend) (*FunctionsOracleEventsMock, error) {
	abi, err := abi.JSON(strings.NewReader(FunctionsOracleEventsMockABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindFunctionsOracleEventsMock(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FunctionsOracleEventsMock{address: address, abi: abi, FunctionsOracleEventsMockCaller: FunctionsOracleEventsMockCaller{contract: contract}, FunctionsOracleEventsMockTransactor: FunctionsOracleEventsMockTransactor{contract: contract}, FunctionsOracleEventsMockFilterer: FunctionsOracleEventsMockFilterer{contract: contract}}, nil
}

func NewFunctionsOracleEventsMockCaller(address common.Address, caller bind.ContractCaller) (*FunctionsOracleEventsMockCaller, error) {
	contract, err := bindFunctionsOracleEventsMock(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FunctionsOracleEventsMockCaller{contract: contract}, nil
}

func NewFunctionsOracleEventsMockTransactor(address common.Address, transactor bind.ContractTransactor) (*FunctionsOracleEventsMockTransactor, error) {
	contract, err := bindFunctionsOracleEventsMock(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FunctionsOracleEventsMockTransactor{contract: contract}, nil
}

func NewFunctionsOracleEventsMockFilterer(address common.Address, filterer bind.ContractFilterer) (*FunctionsOracleEventsMockFilterer, error) {
	contract, err := bindFunctionsOracleEventsMock(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FunctionsOracleEventsMockFilterer{contract: contract}, nil
}

func bindFunctionsOracleEventsMock(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FunctionsOracleEventsMockMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FunctionsOracleEventsMock.Contract.FunctionsOracleEventsMockCaller.contract.Call(opts, result, method, params...)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.Contract.FunctionsOracleEventsMockTransactor.contract.Transfer(opts)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.Contract.FunctionsOracleEventsMockTransactor.contract.Transact(opts, method, params...)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FunctionsOracleEventsMock.Contract.contract.Call(opts, result, method, params...)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.Contract.contract.Transfer(opts)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.Contract.contract.Transact(opts, method, params...)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockTransactor) EmitAuthorizedSendersActive(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.contract.Transact(opts, "emitAuthorizedSendersActive", account)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockSession) EmitAuthorizedSendersActive(account common.Address) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.Contract.EmitAuthorizedSendersActive(&_FunctionsOracleEventsMock.TransactOpts, account)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockTransactorSession) EmitAuthorizedSendersActive(account common.Address) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.Contract.EmitAuthorizedSendersActive(&_FunctionsOracleEventsMock.TransactOpts, account)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockTransactor) EmitAuthorizedSendersChanged(opts *bind.TransactOpts, senders []common.Address, changedBy common.Address) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.contract.Transact(opts, "emitAuthorizedSendersChanged", senders, changedBy)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockSession) EmitAuthorizedSendersChanged(senders []common.Address, changedBy common.Address) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.Contract.EmitAuthorizedSendersChanged(&_FunctionsOracleEventsMock.TransactOpts, senders, changedBy)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockTransactorSession) EmitAuthorizedSendersChanged(senders []common.Address, changedBy common.Address) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.Contract.EmitAuthorizedSendersChanged(&_FunctionsOracleEventsMock.TransactOpts, senders, changedBy)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockTransactor) EmitAuthorizedSendersDeactive(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.contract.Transact(opts, "emitAuthorizedSendersDeactive", account)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockSession) EmitAuthorizedSendersDeactive(account common.Address) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.Contract.EmitAuthorizedSendersDeactive(&_FunctionsOracleEventsMock.TransactOpts, account)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockTransactorSession) EmitAuthorizedSendersDeactive(account common.Address) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.Contract.EmitAuthorizedSendersDeactive(&_FunctionsOracleEventsMock.TransactOpts, account)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockTransactor) EmitConfigSet(opts *bind.TransactOpts, previousConfigBlockNumber uint32, configDigest [32]byte, configCount uint64, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.contract.Transact(opts, "emitConfigSet", previousConfigBlockNumber, configDigest, configCount, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockSession) EmitConfigSet(previousConfigBlockNumber uint32, configDigest [32]byte, configCount uint64, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.Contract.EmitConfigSet(&_FunctionsOracleEventsMock.TransactOpts, previousConfigBlockNumber, configDigest, configCount, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockTransactorSession) EmitConfigSet(previousConfigBlockNumber uint32, configDigest [32]byte, configCount uint64, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.Contract.EmitConfigSet(&_FunctionsOracleEventsMock.TransactOpts, previousConfigBlockNumber, configDigest, configCount, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockTransactor) EmitInitialized(opts *bind.TransactOpts, version uint8) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.contract.Transact(opts, "emitInitialized", version)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockSession) EmitInitialized(version uint8) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.Contract.EmitInitialized(&_FunctionsOracleEventsMock.TransactOpts, version)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockTransactorSession) EmitInitialized(version uint8) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.Contract.EmitInitialized(&_FunctionsOracleEventsMock.TransactOpts, version)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockTransactor) EmitInvalidRequestID(opts *bind.TransactOpts, requestId [32]byte) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.contract.Transact(opts, "emitInvalidRequestID", requestId)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockSession) EmitInvalidRequestID(requestId [32]byte) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.Contract.EmitInvalidRequestID(&_FunctionsOracleEventsMock.TransactOpts, requestId)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockTransactorSession) EmitInvalidRequestID(requestId [32]byte) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.Contract.EmitInvalidRequestID(&_FunctionsOracleEventsMock.TransactOpts, requestId)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockTransactor) EmitOracleRequest(opts *bind.TransactOpts, requestId [32]byte, requestingContract common.Address, requestInitiator common.Address, subscriptionId uint64, subscriptionOwner common.Address, data []byte) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.contract.Transact(opts, "emitOracleRequest", requestId, requestingContract, requestInitiator, subscriptionId, subscriptionOwner, data)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockSession) EmitOracleRequest(requestId [32]byte, requestingContract common.Address, requestInitiator common.Address, subscriptionId uint64, subscriptionOwner common.Address, data []byte) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.Contract.EmitOracleRequest(&_FunctionsOracleEventsMock.TransactOpts, requestId, requestingContract, requestInitiator, subscriptionId, subscriptionOwner, data)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockTransactorSession) EmitOracleRequest(requestId [32]byte, requestingContract common.Address, requestInitiator common.Address, subscriptionId uint64, subscriptionOwner common.Address, data []byte) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.Contract.EmitOracleRequest(&_FunctionsOracleEventsMock.TransactOpts, requestId, requestingContract, requestInitiator, subscriptionId, subscriptionOwner, data)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockTransactor) EmitOracleResponse(opts *bind.TransactOpts, requestId [32]byte) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.contract.Transact(opts, "emitOracleResponse", requestId)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockSession) EmitOracleResponse(requestId [32]byte) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.Contract.EmitOracleResponse(&_FunctionsOracleEventsMock.TransactOpts, requestId)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockTransactorSession) EmitOracleResponse(requestId [32]byte) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.Contract.EmitOracleResponse(&_FunctionsOracleEventsMock.TransactOpts, requestId)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockTransactor) EmitOwnershipTransferRequested(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.contract.Transact(opts, "emitOwnershipTransferRequested", from, to)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockSession) EmitOwnershipTransferRequested(from common.Address, to common.Address) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.Contract.EmitOwnershipTransferRequested(&_FunctionsOracleEventsMock.TransactOpts, from, to)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockTransactorSession) EmitOwnershipTransferRequested(from common.Address, to common.Address) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.Contract.EmitOwnershipTransferRequested(&_FunctionsOracleEventsMock.TransactOpts, from, to)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockTransactor) EmitOwnershipTransferred(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.contract.Transact(opts, "emitOwnershipTransferred", from, to)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockSession) EmitOwnershipTransferred(from common.Address, to common.Address) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.Contract.EmitOwnershipTransferred(&_FunctionsOracleEventsMock.TransactOpts, from, to)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockTransactorSession) EmitOwnershipTransferred(from common.Address, to common.Address) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.Contract.EmitOwnershipTransferred(&_FunctionsOracleEventsMock.TransactOpts, from, to)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockTransactor) EmitResponseTransmitted(opts *bind.TransactOpts, requestId [32]byte, transmitter common.Address) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.contract.Transact(opts, "emitResponseTransmitted", requestId, transmitter)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockSession) EmitResponseTransmitted(requestId [32]byte, transmitter common.Address) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.Contract.EmitResponseTransmitted(&_FunctionsOracleEventsMock.TransactOpts, requestId, transmitter)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockTransactorSession) EmitResponseTransmitted(requestId [32]byte, transmitter common.Address) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.Contract.EmitResponseTransmitted(&_FunctionsOracleEventsMock.TransactOpts, requestId, transmitter)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockTransactor) EmitTransmitted(opts *bind.TransactOpts, configDigest [32]byte, epoch uint32) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.contract.Transact(opts, "emitTransmitted", configDigest, epoch)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockSession) EmitTransmitted(configDigest [32]byte, epoch uint32) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.Contract.EmitTransmitted(&_FunctionsOracleEventsMock.TransactOpts, configDigest, epoch)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockTransactorSession) EmitTransmitted(configDigest [32]byte, epoch uint32) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.Contract.EmitTransmitted(&_FunctionsOracleEventsMock.TransactOpts, configDigest, epoch)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockTransactor) EmitUserCallbackError(opts *bind.TransactOpts, requestId [32]byte, reason string) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.contract.Transact(opts, "emitUserCallbackError", requestId, reason)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockSession) EmitUserCallbackError(requestId [32]byte, reason string) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.Contract.EmitUserCallbackError(&_FunctionsOracleEventsMock.TransactOpts, requestId, reason)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockTransactorSession) EmitUserCallbackError(requestId [32]byte, reason string) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.Contract.EmitUserCallbackError(&_FunctionsOracleEventsMock.TransactOpts, requestId, reason)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockTransactor) EmitUserCallbackRawError(opts *bind.TransactOpts, requestId [32]byte, lowLevelData []byte) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.contract.Transact(opts, "emitUserCallbackRawError", requestId, lowLevelData)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockSession) EmitUserCallbackRawError(requestId [32]byte, lowLevelData []byte) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.Contract.EmitUserCallbackRawError(&_FunctionsOracleEventsMock.TransactOpts, requestId, lowLevelData)
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockTransactorSession) EmitUserCallbackRawError(requestId [32]byte, lowLevelData []byte) (*types.Transaction, error) {
	return _FunctionsOracleEventsMock.Contract.EmitUserCallbackRawError(&_FunctionsOracleEventsMock.TransactOpts, requestId, lowLevelData)
}

type FunctionsOracleEventsMockAuthorizedSendersActiveIterator struct {
	Event *FunctionsOracleEventsMockAuthorizedSendersActive

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsOracleEventsMockAuthorizedSendersActiveIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsOracleEventsMockAuthorizedSendersActive)
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
		it.Event = new(FunctionsOracleEventsMockAuthorizedSendersActive)
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

func (it *FunctionsOracleEventsMockAuthorizedSendersActiveIterator) Error() error {
	return it.fail
}

func (it *FunctionsOracleEventsMockAuthorizedSendersActiveIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsOracleEventsMockAuthorizedSendersActive struct {
	Account common.Address
	Raw     types.Log
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockFilterer) FilterAuthorizedSendersActive(opts *bind.FilterOpts) (*FunctionsOracleEventsMockAuthorizedSendersActiveIterator, error) {

	logs, sub, err := _FunctionsOracleEventsMock.contract.FilterLogs(opts, "AuthorizedSendersActive")
	if err != nil {
		return nil, err
	}
	return &FunctionsOracleEventsMockAuthorizedSendersActiveIterator{contract: _FunctionsOracleEventsMock.contract, event: "AuthorizedSendersActive", logs: logs, sub: sub}, nil
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockFilterer) WatchAuthorizedSendersActive(opts *bind.WatchOpts, sink chan<- *FunctionsOracleEventsMockAuthorizedSendersActive) (event.Subscription, error) {

	logs, sub, err := _FunctionsOracleEventsMock.contract.WatchLogs(opts, "AuthorizedSendersActive")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsOracleEventsMockAuthorizedSendersActive)
				if err := _FunctionsOracleEventsMock.contract.UnpackLog(event, "AuthorizedSendersActive", log); err != nil {
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

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockFilterer) ParseAuthorizedSendersActive(log types.Log) (*FunctionsOracleEventsMockAuthorizedSendersActive, error) {
	event := new(FunctionsOracleEventsMockAuthorizedSendersActive)
	if err := _FunctionsOracleEventsMock.contract.UnpackLog(event, "AuthorizedSendersActive", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsOracleEventsMockAuthorizedSendersChangedIterator struct {
	Event *FunctionsOracleEventsMockAuthorizedSendersChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsOracleEventsMockAuthorizedSendersChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsOracleEventsMockAuthorizedSendersChanged)
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
		it.Event = new(FunctionsOracleEventsMockAuthorizedSendersChanged)
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

func (it *FunctionsOracleEventsMockAuthorizedSendersChangedIterator) Error() error {
	return it.fail
}

func (it *FunctionsOracleEventsMockAuthorizedSendersChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsOracleEventsMockAuthorizedSendersChanged struct {
	Senders   []common.Address
	ChangedBy common.Address
	Raw       types.Log
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockFilterer) FilterAuthorizedSendersChanged(opts *bind.FilterOpts) (*FunctionsOracleEventsMockAuthorizedSendersChangedIterator, error) {

	logs, sub, err := _FunctionsOracleEventsMock.contract.FilterLogs(opts, "AuthorizedSendersChanged")
	if err != nil {
		return nil, err
	}
	return &FunctionsOracleEventsMockAuthorizedSendersChangedIterator{contract: _FunctionsOracleEventsMock.contract, event: "AuthorizedSendersChanged", logs: logs, sub: sub}, nil
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockFilterer) WatchAuthorizedSendersChanged(opts *bind.WatchOpts, sink chan<- *FunctionsOracleEventsMockAuthorizedSendersChanged) (event.Subscription, error) {

	logs, sub, err := _FunctionsOracleEventsMock.contract.WatchLogs(opts, "AuthorizedSendersChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsOracleEventsMockAuthorizedSendersChanged)
				if err := _FunctionsOracleEventsMock.contract.UnpackLog(event, "AuthorizedSendersChanged", log); err != nil {
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

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockFilterer) ParseAuthorizedSendersChanged(log types.Log) (*FunctionsOracleEventsMockAuthorizedSendersChanged, error) {
	event := new(FunctionsOracleEventsMockAuthorizedSendersChanged)
	if err := _FunctionsOracleEventsMock.contract.UnpackLog(event, "AuthorizedSendersChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsOracleEventsMockAuthorizedSendersDeactiveIterator struct {
	Event *FunctionsOracleEventsMockAuthorizedSendersDeactive

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsOracleEventsMockAuthorizedSendersDeactiveIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsOracleEventsMockAuthorizedSendersDeactive)
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
		it.Event = new(FunctionsOracleEventsMockAuthorizedSendersDeactive)
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

func (it *FunctionsOracleEventsMockAuthorizedSendersDeactiveIterator) Error() error {
	return it.fail
}

func (it *FunctionsOracleEventsMockAuthorizedSendersDeactiveIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsOracleEventsMockAuthorizedSendersDeactive struct {
	Account common.Address
	Raw     types.Log
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockFilterer) FilterAuthorizedSendersDeactive(opts *bind.FilterOpts) (*FunctionsOracleEventsMockAuthorizedSendersDeactiveIterator, error) {

	logs, sub, err := _FunctionsOracleEventsMock.contract.FilterLogs(opts, "AuthorizedSendersDeactive")
	if err != nil {
		return nil, err
	}
	return &FunctionsOracleEventsMockAuthorizedSendersDeactiveIterator{contract: _FunctionsOracleEventsMock.contract, event: "AuthorizedSendersDeactive", logs: logs, sub: sub}, nil
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockFilterer) WatchAuthorizedSendersDeactive(opts *bind.WatchOpts, sink chan<- *FunctionsOracleEventsMockAuthorizedSendersDeactive) (event.Subscription, error) {

	logs, sub, err := _FunctionsOracleEventsMock.contract.WatchLogs(opts, "AuthorizedSendersDeactive")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsOracleEventsMockAuthorizedSendersDeactive)
				if err := _FunctionsOracleEventsMock.contract.UnpackLog(event, "AuthorizedSendersDeactive", log); err != nil {
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

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockFilterer) ParseAuthorizedSendersDeactive(log types.Log) (*FunctionsOracleEventsMockAuthorizedSendersDeactive, error) {
	event := new(FunctionsOracleEventsMockAuthorizedSendersDeactive)
	if err := _FunctionsOracleEventsMock.contract.UnpackLog(event, "AuthorizedSendersDeactive", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsOracleEventsMockConfigSetIterator struct {
	Event *FunctionsOracleEventsMockConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsOracleEventsMockConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsOracleEventsMockConfigSet)
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
		it.Event = new(FunctionsOracleEventsMockConfigSet)
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

func (it *FunctionsOracleEventsMockConfigSetIterator) Error() error {
	return it.fail
}

func (it *FunctionsOracleEventsMockConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsOracleEventsMockConfigSet struct {
	PreviousConfigBlockNumber uint32
	ConfigDigest              [32]byte
	ConfigCount               uint64
	Signers                   []common.Address
	Transmitters              []common.Address
	F                         uint8
	OnchainConfig             []byte
	OffchainConfigVersion     uint64
	OffchainConfig            []byte
	Raw                       types.Log
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockFilterer) FilterConfigSet(opts *bind.FilterOpts) (*FunctionsOracleEventsMockConfigSetIterator, error) {

	logs, sub, err := _FunctionsOracleEventsMock.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &FunctionsOracleEventsMockConfigSetIterator{contract: _FunctionsOracleEventsMock.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *FunctionsOracleEventsMockConfigSet) (event.Subscription, error) {

	logs, sub, err := _FunctionsOracleEventsMock.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsOracleEventsMockConfigSet)
				if err := _FunctionsOracleEventsMock.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockFilterer) ParseConfigSet(log types.Log) (*FunctionsOracleEventsMockConfigSet, error) {
	event := new(FunctionsOracleEventsMockConfigSet)
	if err := _FunctionsOracleEventsMock.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsOracleEventsMockInitializedIterator struct {
	Event *FunctionsOracleEventsMockInitialized

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsOracleEventsMockInitializedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsOracleEventsMockInitialized)
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
		it.Event = new(FunctionsOracleEventsMockInitialized)
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

func (it *FunctionsOracleEventsMockInitializedIterator) Error() error {
	return it.fail
}

func (it *FunctionsOracleEventsMockInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsOracleEventsMockInitialized struct {
	Version uint8
	Raw     types.Log
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockFilterer) FilterInitialized(opts *bind.FilterOpts) (*FunctionsOracleEventsMockInitializedIterator, error) {

	logs, sub, err := _FunctionsOracleEventsMock.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &FunctionsOracleEventsMockInitializedIterator{contract: _FunctionsOracleEventsMock.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *FunctionsOracleEventsMockInitialized) (event.Subscription, error) {

	logs, sub, err := _FunctionsOracleEventsMock.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsOracleEventsMockInitialized)
				if err := _FunctionsOracleEventsMock.contract.UnpackLog(event, "Initialized", log); err != nil {
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

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockFilterer) ParseInitialized(log types.Log) (*FunctionsOracleEventsMockInitialized, error) {
	event := new(FunctionsOracleEventsMockInitialized)
	if err := _FunctionsOracleEventsMock.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsOracleEventsMockInvalidRequestIDIterator struct {
	Event *FunctionsOracleEventsMockInvalidRequestID

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsOracleEventsMockInvalidRequestIDIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsOracleEventsMockInvalidRequestID)
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
		it.Event = new(FunctionsOracleEventsMockInvalidRequestID)
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

func (it *FunctionsOracleEventsMockInvalidRequestIDIterator) Error() error {
	return it.fail
}

func (it *FunctionsOracleEventsMockInvalidRequestIDIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsOracleEventsMockInvalidRequestID struct {
	RequestId [32]byte
	Raw       types.Log
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockFilterer) FilterInvalidRequestID(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsOracleEventsMockInvalidRequestIDIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsOracleEventsMock.contract.FilterLogs(opts, "InvalidRequestID", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsOracleEventsMockInvalidRequestIDIterator{contract: _FunctionsOracleEventsMock.contract, event: "InvalidRequestID", logs: logs, sub: sub}, nil
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockFilterer) WatchInvalidRequestID(opts *bind.WatchOpts, sink chan<- *FunctionsOracleEventsMockInvalidRequestID, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsOracleEventsMock.contract.WatchLogs(opts, "InvalidRequestID", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsOracleEventsMockInvalidRequestID)
				if err := _FunctionsOracleEventsMock.contract.UnpackLog(event, "InvalidRequestID", log); err != nil {
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

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockFilterer) ParseInvalidRequestID(log types.Log) (*FunctionsOracleEventsMockInvalidRequestID, error) {
	event := new(FunctionsOracleEventsMockInvalidRequestID)
	if err := _FunctionsOracleEventsMock.contract.UnpackLog(event, "InvalidRequestID", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsOracleEventsMockOracleRequestIterator struct {
	Event *FunctionsOracleEventsMockOracleRequest

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsOracleEventsMockOracleRequestIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsOracleEventsMockOracleRequest)
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
		it.Event = new(FunctionsOracleEventsMockOracleRequest)
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

func (it *FunctionsOracleEventsMockOracleRequestIterator) Error() error {
	return it.fail
}

func (it *FunctionsOracleEventsMockOracleRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsOracleEventsMockOracleRequest struct {
	RequestId          [32]byte
	RequestingContract common.Address
	RequestInitiator   common.Address
	SubscriptionId     uint64
	SubscriptionOwner  common.Address
	Data               []byte
	Raw                types.Log
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockFilterer) FilterOracleRequest(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsOracleEventsMockOracleRequestIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsOracleEventsMock.contract.FilterLogs(opts, "OracleRequest", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsOracleEventsMockOracleRequestIterator{contract: _FunctionsOracleEventsMock.contract, event: "OracleRequest", logs: logs, sub: sub}, nil
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockFilterer) WatchOracleRequest(opts *bind.WatchOpts, sink chan<- *FunctionsOracleEventsMockOracleRequest, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsOracleEventsMock.contract.WatchLogs(opts, "OracleRequest", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsOracleEventsMockOracleRequest)
				if err := _FunctionsOracleEventsMock.contract.UnpackLog(event, "OracleRequest", log); err != nil {
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

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockFilterer) ParseOracleRequest(log types.Log) (*FunctionsOracleEventsMockOracleRequest, error) {
	event := new(FunctionsOracleEventsMockOracleRequest)
	if err := _FunctionsOracleEventsMock.contract.UnpackLog(event, "OracleRequest", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsOracleEventsMockOracleResponseIterator struct {
	Event *FunctionsOracleEventsMockOracleResponse

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsOracleEventsMockOracleResponseIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsOracleEventsMockOracleResponse)
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
		it.Event = new(FunctionsOracleEventsMockOracleResponse)
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

func (it *FunctionsOracleEventsMockOracleResponseIterator) Error() error {
	return it.fail
}

func (it *FunctionsOracleEventsMockOracleResponseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsOracleEventsMockOracleResponse struct {
	RequestId [32]byte
	Raw       types.Log
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockFilterer) FilterOracleResponse(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsOracleEventsMockOracleResponseIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsOracleEventsMock.contract.FilterLogs(opts, "OracleResponse", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsOracleEventsMockOracleResponseIterator{contract: _FunctionsOracleEventsMock.contract, event: "OracleResponse", logs: logs, sub: sub}, nil
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockFilterer) WatchOracleResponse(opts *bind.WatchOpts, sink chan<- *FunctionsOracleEventsMockOracleResponse, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsOracleEventsMock.contract.WatchLogs(opts, "OracleResponse", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsOracleEventsMockOracleResponse)
				if err := _FunctionsOracleEventsMock.contract.UnpackLog(event, "OracleResponse", log); err != nil {
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

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockFilterer) ParseOracleResponse(log types.Log) (*FunctionsOracleEventsMockOracleResponse, error) {
	event := new(FunctionsOracleEventsMockOracleResponse)
	if err := _FunctionsOracleEventsMock.contract.UnpackLog(event, "OracleResponse", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsOracleEventsMockOwnershipTransferRequestedIterator struct {
	Event *FunctionsOracleEventsMockOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsOracleEventsMockOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsOracleEventsMockOwnershipTransferRequested)
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
		it.Event = new(FunctionsOracleEventsMockOwnershipTransferRequested)
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

func (it *FunctionsOracleEventsMockOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *FunctionsOracleEventsMockOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsOracleEventsMockOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsOracleEventsMockOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsOracleEventsMock.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsOracleEventsMockOwnershipTransferRequestedIterator{contract: _FunctionsOracleEventsMock.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *FunctionsOracleEventsMockOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsOracleEventsMock.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsOracleEventsMockOwnershipTransferRequested)
				if err := _FunctionsOracleEventsMock.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockFilterer) ParseOwnershipTransferRequested(log types.Log) (*FunctionsOracleEventsMockOwnershipTransferRequested, error) {
	event := new(FunctionsOracleEventsMockOwnershipTransferRequested)
	if err := _FunctionsOracleEventsMock.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsOracleEventsMockOwnershipTransferredIterator struct {
	Event *FunctionsOracleEventsMockOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsOracleEventsMockOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsOracleEventsMockOwnershipTransferred)
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
		it.Event = new(FunctionsOracleEventsMockOwnershipTransferred)
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

func (it *FunctionsOracleEventsMockOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *FunctionsOracleEventsMockOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsOracleEventsMockOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsOracleEventsMockOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsOracleEventsMock.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsOracleEventsMockOwnershipTransferredIterator{contract: _FunctionsOracleEventsMock.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FunctionsOracleEventsMockOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsOracleEventsMock.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsOracleEventsMockOwnershipTransferred)
				if err := _FunctionsOracleEventsMock.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockFilterer) ParseOwnershipTransferred(log types.Log) (*FunctionsOracleEventsMockOwnershipTransferred, error) {
	event := new(FunctionsOracleEventsMockOwnershipTransferred)
	if err := _FunctionsOracleEventsMock.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsOracleEventsMockResponseTransmittedIterator struct {
	Event *FunctionsOracleEventsMockResponseTransmitted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsOracleEventsMockResponseTransmittedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsOracleEventsMockResponseTransmitted)
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
		it.Event = new(FunctionsOracleEventsMockResponseTransmitted)
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

func (it *FunctionsOracleEventsMockResponseTransmittedIterator) Error() error {
	return it.fail
}

func (it *FunctionsOracleEventsMockResponseTransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsOracleEventsMockResponseTransmitted struct {
	RequestId   [32]byte
	Transmitter common.Address
	Raw         types.Log
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockFilterer) FilterResponseTransmitted(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsOracleEventsMockResponseTransmittedIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsOracleEventsMock.contract.FilterLogs(opts, "ResponseTransmitted", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsOracleEventsMockResponseTransmittedIterator{contract: _FunctionsOracleEventsMock.contract, event: "ResponseTransmitted", logs: logs, sub: sub}, nil
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockFilterer) WatchResponseTransmitted(opts *bind.WatchOpts, sink chan<- *FunctionsOracleEventsMockResponseTransmitted, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsOracleEventsMock.contract.WatchLogs(opts, "ResponseTransmitted", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsOracleEventsMockResponseTransmitted)
				if err := _FunctionsOracleEventsMock.contract.UnpackLog(event, "ResponseTransmitted", log); err != nil {
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

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockFilterer) ParseResponseTransmitted(log types.Log) (*FunctionsOracleEventsMockResponseTransmitted, error) {
	event := new(FunctionsOracleEventsMockResponseTransmitted)
	if err := _FunctionsOracleEventsMock.contract.UnpackLog(event, "ResponseTransmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsOracleEventsMockTransmittedIterator struct {
	Event *FunctionsOracleEventsMockTransmitted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsOracleEventsMockTransmittedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsOracleEventsMockTransmitted)
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
		it.Event = new(FunctionsOracleEventsMockTransmitted)
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

func (it *FunctionsOracleEventsMockTransmittedIterator) Error() error {
	return it.fail
}

func (it *FunctionsOracleEventsMockTransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsOracleEventsMockTransmitted struct {
	ConfigDigest [32]byte
	Epoch        uint32
	Raw          types.Log
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockFilterer) FilterTransmitted(opts *bind.FilterOpts) (*FunctionsOracleEventsMockTransmittedIterator, error) {

	logs, sub, err := _FunctionsOracleEventsMock.contract.FilterLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return &FunctionsOracleEventsMockTransmittedIterator{contract: _FunctionsOracleEventsMock.contract, event: "Transmitted", logs: logs, sub: sub}, nil
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockFilterer) WatchTransmitted(opts *bind.WatchOpts, sink chan<- *FunctionsOracleEventsMockTransmitted) (event.Subscription, error) {

	logs, sub, err := _FunctionsOracleEventsMock.contract.WatchLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsOracleEventsMockTransmitted)
				if err := _FunctionsOracleEventsMock.contract.UnpackLog(event, "Transmitted", log); err != nil {
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

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockFilterer) ParseTransmitted(log types.Log) (*FunctionsOracleEventsMockTransmitted, error) {
	event := new(FunctionsOracleEventsMockTransmitted)
	if err := _FunctionsOracleEventsMock.contract.UnpackLog(event, "Transmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsOracleEventsMockUserCallbackErrorIterator struct {
	Event *FunctionsOracleEventsMockUserCallbackError

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsOracleEventsMockUserCallbackErrorIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsOracleEventsMockUserCallbackError)
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
		it.Event = new(FunctionsOracleEventsMockUserCallbackError)
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

func (it *FunctionsOracleEventsMockUserCallbackErrorIterator) Error() error {
	return it.fail
}

func (it *FunctionsOracleEventsMockUserCallbackErrorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsOracleEventsMockUserCallbackError struct {
	RequestId [32]byte
	Reason    string
	Raw       types.Log
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockFilterer) FilterUserCallbackError(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsOracleEventsMockUserCallbackErrorIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsOracleEventsMock.contract.FilterLogs(opts, "UserCallbackError", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsOracleEventsMockUserCallbackErrorIterator{contract: _FunctionsOracleEventsMock.contract, event: "UserCallbackError", logs: logs, sub: sub}, nil
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockFilterer) WatchUserCallbackError(opts *bind.WatchOpts, sink chan<- *FunctionsOracleEventsMockUserCallbackError, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsOracleEventsMock.contract.WatchLogs(opts, "UserCallbackError", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsOracleEventsMockUserCallbackError)
				if err := _FunctionsOracleEventsMock.contract.UnpackLog(event, "UserCallbackError", log); err != nil {
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

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockFilterer) ParseUserCallbackError(log types.Log) (*FunctionsOracleEventsMockUserCallbackError, error) {
	event := new(FunctionsOracleEventsMockUserCallbackError)
	if err := _FunctionsOracleEventsMock.contract.UnpackLog(event, "UserCallbackError", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsOracleEventsMockUserCallbackRawErrorIterator struct {
	Event *FunctionsOracleEventsMockUserCallbackRawError

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsOracleEventsMockUserCallbackRawErrorIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsOracleEventsMockUserCallbackRawError)
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
		it.Event = new(FunctionsOracleEventsMockUserCallbackRawError)
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

func (it *FunctionsOracleEventsMockUserCallbackRawErrorIterator) Error() error {
	return it.fail
}

func (it *FunctionsOracleEventsMockUserCallbackRawErrorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsOracleEventsMockUserCallbackRawError struct {
	RequestId    [32]byte
	LowLevelData []byte
	Raw          types.Log
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockFilterer) FilterUserCallbackRawError(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsOracleEventsMockUserCallbackRawErrorIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsOracleEventsMock.contract.FilterLogs(opts, "UserCallbackRawError", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsOracleEventsMockUserCallbackRawErrorIterator{contract: _FunctionsOracleEventsMock.contract, event: "UserCallbackRawError", logs: logs, sub: sub}, nil
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockFilterer) WatchUserCallbackRawError(opts *bind.WatchOpts, sink chan<- *FunctionsOracleEventsMockUserCallbackRawError, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _FunctionsOracleEventsMock.contract.WatchLogs(opts, "UserCallbackRawError", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsOracleEventsMockUserCallbackRawError)
				if err := _FunctionsOracleEventsMock.contract.UnpackLog(event, "UserCallbackRawError", log); err != nil {
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

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMockFilterer) ParseUserCallbackRawError(log types.Log) (*FunctionsOracleEventsMockUserCallbackRawError, error) {
	event := new(FunctionsOracleEventsMockUserCallbackRawError)
	if err := _FunctionsOracleEventsMock.contract.UnpackLog(event, "UserCallbackRawError", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMock) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _FunctionsOracleEventsMock.abi.Events["AuthorizedSendersActive"].ID:
		return _FunctionsOracleEventsMock.ParseAuthorizedSendersActive(log)
	case _FunctionsOracleEventsMock.abi.Events["AuthorizedSendersChanged"].ID:
		return _FunctionsOracleEventsMock.ParseAuthorizedSendersChanged(log)
	case _FunctionsOracleEventsMock.abi.Events["AuthorizedSendersDeactive"].ID:
		return _FunctionsOracleEventsMock.ParseAuthorizedSendersDeactive(log)
	case _FunctionsOracleEventsMock.abi.Events["ConfigSet"].ID:
		return _FunctionsOracleEventsMock.ParseConfigSet(log)
	case _FunctionsOracleEventsMock.abi.Events["Initialized"].ID:
		return _FunctionsOracleEventsMock.ParseInitialized(log)
	case _FunctionsOracleEventsMock.abi.Events["InvalidRequestID"].ID:
		return _FunctionsOracleEventsMock.ParseInvalidRequestID(log)
	case _FunctionsOracleEventsMock.abi.Events["OracleRequest"].ID:
		return _FunctionsOracleEventsMock.ParseOracleRequest(log)
	case _FunctionsOracleEventsMock.abi.Events["OracleResponse"].ID:
		return _FunctionsOracleEventsMock.ParseOracleResponse(log)
	case _FunctionsOracleEventsMock.abi.Events["OwnershipTransferRequested"].ID:
		return _FunctionsOracleEventsMock.ParseOwnershipTransferRequested(log)
	case _FunctionsOracleEventsMock.abi.Events["OwnershipTransferred"].ID:
		return _FunctionsOracleEventsMock.ParseOwnershipTransferred(log)
	case _FunctionsOracleEventsMock.abi.Events["ResponseTransmitted"].ID:
		return _FunctionsOracleEventsMock.ParseResponseTransmitted(log)
	case _FunctionsOracleEventsMock.abi.Events["Transmitted"].ID:
		return _FunctionsOracleEventsMock.ParseTransmitted(log)
	case _FunctionsOracleEventsMock.abi.Events["UserCallbackError"].ID:
		return _FunctionsOracleEventsMock.ParseUserCallbackError(log)
	case _FunctionsOracleEventsMock.abi.Events["UserCallbackRawError"].ID:
		return _FunctionsOracleEventsMock.ParseUserCallbackRawError(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (FunctionsOracleEventsMockAuthorizedSendersActive) Topic() common.Hash {
	return common.HexToHash("0xae51766a982895b0c444fc99fc1a560762b464d709e6c78376c85617f7eeb5ce")
}

func (FunctionsOracleEventsMockAuthorizedSendersChanged) Topic() common.Hash {
	return common.HexToHash("0xf263cfb3e4298332e776194610cf9fdc09ccb3ada8b9aa39764d882e11fbf0a0")
}

func (FunctionsOracleEventsMockAuthorizedSendersDeactive) Topic() common.Hash {
	return common.HexToHash("0xea3828816a323b8d7ff49d755efd105e7719166d6c76fad97a28eee5eccc3d9a")
}

func (FunctionsOracleEventsMockConfigSet) Topic() common.Hash {
	return common.HexToHash("0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05")
}

func (FunctionsOracleEventsMockInitialized) Topic() common.Hash {
	return common.HexToHash("0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498")
}

func (FunctionsOracleEventsMockInvalidRequestID) Topic() common.Hash {
	return common.HexToHash("0xa1c120e327c9ad8b075793878c88d59b8934b97ae37117faa3bb21616237f7be")
}

func (FunctionsOracleEventsMockOracleRequest) Topic() common.Hash {
	return common.HexToHash("0xa1ec73989d79578cd6f67d4f593ac3e0a4d1020e5c0164db52108d7ff785406c")
}

func (FunctionsOracleEventsMockOracleResponse) Topic() common.Hash {
	return common.HexToHash("0x9e9bc7616d42c2835d05ae617e508454e63b30b934be8aa932ebc125e0e58a64")
}

func (FunctionsOracleEventsMockOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (FunctionsOracleEventsMockOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (FunctionsOracleEventsMockResponseTransmitted) Topic() common.Hash {
	return common.HexToHash("0xdc941eddab34a6109ab77798299c6b1f035b125fd6f774d266ecbf9541d630a6")
}

func (FunctionsOracleEventsMockTransmitted) Topic() common.Hash {
	return common.HexToHash("0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62")
}

func (FunctionsOracleEventsMockUserCallbackError) Topic() common.Hash {
	return common.HexToHash("0xb2931868c372fe17a25643458add467d60ec5c51125a99b7309f41f5bcd2da6c")
}

func (FunctionsOracleEventsMockUserCallbackRawError) Topic() common.Hash {
	return common.HexToHash("0xe0b838ffe6ee22a0d3acf19a85db6a41b34a1ab739e2d6c759a2e42d95bdccb2")
}

func (_FunctionsOracleEventsMock *FunctionsOracleEventsMock) Address() common.Address {
	return _FunctionsOracleEventsMock.address
}

type FunctionsOracleEventsMockInterface interface {
	EmitAuthorizedSendersActive(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error)

	EmitAuthorizedSendersChanged(opts *bind.TransactOpts, senders []common.Address, changedBy common.Address) (*types.Transaction, error)

	EmitAuthorizedSendersDeactive(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error)

	EmitConfigSet(opts *bind.TransactOpts, previousConfigBlockNumber uint32, configDigest [32]byte, configCount uint64, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error)

	EmitInitialized(opts *bind.TransactOpts, version uint8) (*types.Transaction, error)

	EmitInvalidRequestID(opts *bind.TransactOpts, requestId [32]byte) (*types.Transaction, error)

	EmitOracleRequest(opts *bind.TransactOpts, requestId [32]byte, requestingContract common.Address, requestInitiator common.Address, subscriptionId uint64, subscriptionOwner common.Address, data []byte) (*types.Transaction, error)

	EmitOracleResponse(opts *bind.TransactOpts, requestId [32]byte) (*types.Transaction, error)

	EmitOwnershipTransferRequested(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error)

	EmitOwnershipTransferred(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error)

	EmitResponseTransmitted(opts *bind.TransactOpts, requestId [32]byte, transmitter common.Address) (*types.Transaction, error)

	EmitTransmitted(opts *bind.TransactOpts, configDigest [32]byte, epoch uint32) (*types.Transaction, error)

	EmitUserCallbackError(opts *bind.TransactOpts, requestId [32]byte, reason string) (*types.Transaction, error)

	EmitUserCallbackRawError(opts *bind.TransactOpts, requestId [32]byte, lowLevelData []byte) (*types.Transaction, error)

	FilterAuthorizedSendersActive(opts *bind.FilterOpts) (*FunctionsOracleEventsMockAuthorizedSendersActiveIterator, error)

	WatchAuthorizedSendersActive(opts *bind.WatchOpts, sink chan<- *FunctionsOracleEventsMockAuthorizedSendersActive) (event.Subscription, error)

	ParseAuthorizedSendersActive(log types.Log) (*FunctionsOracleEventsMockAuthorizedSendersActive, error)

	FilterAuthorizedSendersChanged(opts *bind.FilterOpts) (*FunctionsOracleEventsMockAuthorizedSendersChangedIterator, error)

	WatchAuthorizedSendersChanged(opts *bind.WatchOpts, sink chan<- *FunctionsOracleEventsMockAuthorizedSendersChanged) (event.Subscription, error)

	ParseAuthorizedSendersChanged(log types.Log) (*FunctionsOracleEventsMockAuthorizedSendersChanged, error)

	FilterAuthorizedSendersDeactive(opts *bind.FilterOpts) (*FunctionsOracleEventsMockAuthorizedSendersDeactiveIterator, error)

	WatchAuthorizedSendersDeactive(opts *bind.WatchOpts, sink chan<- *FunctionsOracleEventsMockAuthorizedSendersDeactive) (event.Subscription, error)

	ParseAuthorizedSendersDeactive(log types.Log) (*FunctionsOracleEventsMockAuthorizedSendersDeactive, error)

	FilterConfigSet(opts *bind.FilterOpts) (*FunctionsOracleEventsMockConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *FunctionsOracleEventsMockConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*FunctionsOracleEventsMockConfigSet, error)

	FilterInitialized(opts *bind.FilterOpts) (*FunctionsOracleEventsMockInitializedIterator, error)

	WatchInitialized(opts *bind.WatchOpts, sink chan<- *FunctionsOracleEventsMockInitialized) (event.Subscription, error)

	ParseInitialized(log types.Log) (*FunctionsOracleEventsMockInitialized, error)

	FilterInvalidRequestID(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsOracleEventsMockInvalidRequestIDIterator, error)

	WatchInvalidRequestID(opts *bind.WatchOpts, sink chan<- *FunctionsOracleEventsMockInvalidRequestID, requestId [][32]byte) (event.Subscription, error)

	ParseInvalidRequestID(log types.Log) (*FunctionsOracleEventsMockInvalidRequestID, error)

	FilterOracleRequest(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsOracleEventsMockOracleRequestIterator, error)

	WatchOracleRequest(opts *bind.WatchOpts, sink chan<- *FunctionsOracleEventsMockOracleRequest, requestId [][32]byte) (event.Subscription, error)

	ParseOracleRequest(log types.Log) (*FunctionsOracleEventsMockOracleRequest, error)

	FilterOracleResponse(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsOracleEventsMockOracleResponseIterator, error)

	WatchOracleResponse(opts *bind.WatchOpts, sink chan<- *FunctionsOracleEventsMockOracleResponse, requestId [][32]byte) (event.Subscription, error)

	ParseOracleResponse(log types.Log) (*FunctionsOracleEventsMockOracleResponse, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsOracleEventsMockOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *FunctionsOracleEventsMockOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*FunctionsOracleEventsMockOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsOracleEventsMockOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FunctionsOracleEventsMockOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*FunctionsOracleEventsMockOwnershipTransferred, error)

	FilterResponseTransmitted(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsOracleEventsMockResponseTransmittedIterator, error)

	WatchResponseTransmitted(opts *bind.WatchOpts, sink chan<- *FunctionsOracleEventsMockResponseTransmitted, requestId [][32]byte) (event.Subscription, error)

	ParseResponseTransmitted(log types.Log) (*FunctionsOracleEventsMockResponseTransmitted, error)

	FilterTransmitted(opts *bind.FilterOpts) (*FunctionsOracleEventsMockTransmittedIterator, error)

	WatchTransmitted(opts *bind.WatchOpts, sink chan<- *FunctionsOracleEventsMockTransmitted) (event.Subscription, error)

	ParseTransmitted(log types.Log) (*FunctionsOracleEventsMockTransmitted, error)

	FilterUserCallbackError(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsOracleEventsMockUserCallbackErrorIterator, error)

	WatchUserCallbackError(opts *bind.WatchOpts, sink chan<- *FunctionsOracleEventsMockUserCallbackError, requestId [][32]byte) (event.Subscription, error)

	ParseUserCallbackError(log types.Log) (*FunctionsOracleEventsMockUserCallbackError, error)

	FilterUserCallbackRawError(opts *bind.FilterOpts, requestId [][32]byte) (*FunctionsOracleEventsMockUserCallbackRawErrorIterator, error)

	WatchUserCallbackRawError(opts *bind.WatchOpts, sink chan<- *FunctionsOracleEventsMockUserCallbackRawError, requestId [][32]byte) (event.Subscription, error)

	ParseUserCallbackRawError(log types.Log) (*FunctionsOracleEventsMockUserCallbackRawError, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
