// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_single_consumer_example

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
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

var VRFSingleConsumerExampleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"requestRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_randomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestConfig\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"subscribe\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"topUpSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unsubscribe\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561001057600080fd5b50604051610bfa380380610bfa83398101604081905261002f91610137565b606086811b6001600160601b0319166080908152600080546001600160a01b03998a166001600160a01b03199182161782556001805490911698909916979097179097556040805160a08101825296875263ffffffff9586166020880181905261ffff959095169087018190529290941693850184905293909401839052600280546001600160701b0319166801000000000000000090920261ffff60601b1916919091176c010000000000000000000000009094029390931763ffffffff60701b1916600160701b909102179091556003556101ad565b80516001600160a01b038116811461011e57600080fd5b919050565b805163ffffffff8116811461011e57600080fd5b60008060008060008060c0878903121561015057600080fd5b61015987610107565b955061016760208801610107565b945061017560408801610123565b9350606087015161ffff8116811461018c57600080fd5b925061019a60808801610123565b915060a087015190509295509295509295565b60805160601c610a286101d2600039600081816101af01526102170152610a286000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c8063e0c862891161005b578063e0c862891461015d578063e89e106a14610165578063f6eaffc81461017c578063fcae44841461018f57600080fd5b80631fe543e31461008d5780637db9263f146100a257806386850e93146101425780638f449a0514610155575b600080fd5b6100a061009b36600461080c565b610197565b005b6002546003546100fb9167ffffffffffffffff81169163ffffffff68010000000000000000830481169261ffff6c01000000000000000000000000820416926e0100000000000000000000000000009091049091169085565b6040805167ffffffffffffffff909616865263ffffffff948516602087015261ffff90931692850192909252919091166060830152608082015260a0015b60405180910390f35b6100a06101503660046107da565b610256565b6100a0610317565b6100a061051a565b61016e60055481565b604051908152602001610139565b61016e61018a3660046107da565b610657565b6100a0610678565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610248576040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016602482015260440160405180910390fd5b6102528282610739565b5050565b6001546000546002546040805167ffffffffffffffff909216602083015273ffffffffffffffffffffffffffffffffffffffff93841693634000aea09316918591016040516020818303038152906040526040518463ffffffff1660e01b81526004016102c593929190610925565b602060405180830381600087803b1580156102df57600080fd5b505af11580156102f3573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061025291906107b1565b60408051600180825281830190925260009160208083019080368337019050509050308160008151811061034d5761034d6109bd565b602002602001019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff168152505060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a21a23e46040518163ffffffff1660e01b8152600401602060405180830381600087803b1580156103ef57600080fd5b505af1158015610403573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061042791906108fb565b600280547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001667ffffffffffffffff92909216918217905560008054835173ffffffffffffffffffffffffffffffffffffffff90911692637341c10c929091859190610495576104956109bd565b60200260200101516040518363ffffffff1660e01b81526004016104e592919067ffffffffffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b600060405180830381600087803b1580156104ff57600080fd5b505af1158015610513573d6000803e3d6000fd5b5050505050565b6040805160a08101825260025467ffffffffffffffff811680835263ffffffff68010000000000000000830481166020850181905261ffff6c010000000000000000000000008504168587018190526e010000000000000000000000000000909404909116606085018190526003546080860181905260005496517f5d3b1d3000000000000000000000000000000000000000000000000000000000815260048101919091526024810193909352604483019390935260648201526084810191909152909173ffffffffffffffffffffffffffffffffffffffff1690635d3b1d309060a401602060405180830381600087803b15801561061957600080fd5b505af115801561062d573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061065191906107f3565b60055550565b6004818154811061066757600080fd5b600091825260209091200154905081565b6000546002546040517fd7ae1d3000000000000000000000000000000000000000000000000000000000815267ffffffffffffffff909116600482015230602482015273ffffffffffffffffffffffffffffffffffffffff9091169063d7ae1d3090604401600060405180830381600087803b1580156106f757600080fd5b505af115801561070b573d6000803e3d6000fd5b5050600280547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001690555050565b805161074c906004906020840190610751565b505050565b82805482825590600052602060002090810192821561078c579160200282015b8281111561078c578251825591602001919060010190610771565b5061079892915061079c565b5090565b5b80821115610798576000815560010161079d565b6000602082840312156107c357600080fd5b815180151581146107d357600080fd5b9392505050565b6000602082840312156107ec57600080fd5b5035919050565b60006020828403121561080557600080fd5b5051919050565b6000806040838503121561081f57600080fd5b8235915060208084013567ffffffffffffffff8082111561083f57600080fd5b818601915086601f83011261085357600080fd5b813581811115610865576108656109ec565b8060051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f830116810181811085821117156108a8576108a86109ec565b604052828152858101935084860182860187018b10156108c757600080fd5b600095505b838610156108ea5780358552600195909501949386019386016108cc565b508096505050505050509250929050565b60006020828403121561090d57600080fd5b815167ffffffffffffffff811681146107d357600080fd5b73ffffffffffffffffffffffffffffffffffffffff8416815260006020848184015260606040840152835180606085015260005b8181101561097557858101830151858201608001528201610959565b81811115610987576000608083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160800195945050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var VRFSingleConsumerExampleABI = VRFSingleConsumerExampleMetaData.ABI

var VRFSingleConsumerExampleBin = VRFSingleConsumerExampleMetaData.Bin

func DeployVRFSingleConsumerExample(auth *bind.TransactOpts, backend bind.ContractBackend, vrfCoordinator common.Address, link common.Address, callbackGasLimit uint32, requestConfirmations uint16, numWords uint32, keyHash [32]byte) (common.Address, *types.Transaction, *VRFSingleConsumerExample, error) {
	parsed, err := VRFSingleConsumerExampleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFSingleConsumerExampleBin), backend, vrfCoordinator, link, callbackGasLimit, requestConfirmations, numWords, keyHash)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFSingleConsumerExample{VRFSingleConsumerExampleCaller: VRFSingleConsumerExampleCaller{contract: contract}, VRFSingleConsumerExampleTransactor: VRFSingleConsumerExampleTransactor{contract: contract}, VRFSingleConsumerExampleFilterer: VRFSingleConsumerExampleFilterer{contract: contract}}, nil
}

type VRFSingleConsumerExample struct {
	address common.Address
	abi     abi.ABI
	VRFSingleConsumerExampleCaller
	VRFSingleConsumerExampleTransactor
	VRFSingleConsumerExampleFilterer
}

type VRFSingleConsumerExampleCaller struct {
	contract *bind.BoundContract
}

type VRFSingleConsumerExampleTransactor struct {
	contract *bind.BoundContract
}

type VRFSingleConsumerExampleFilterer struct {
	contract *bind.BoundContract
}

type VRFSingleConsumerExampleSession struct {
	Contract     *VRFSingleConsumerExample
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFSingleConsumerExampleCallerSession struct {
	Contract *VRFSingleConsumerExampleCaller
	CallOpts bind.CallOpts
}

type VRFSingleConsumerExampleTransactorSession struct {
	Contract     *VRFSingleConsumerExampleTransactor
	TransactOpts bind.TransactOpts
}

type VRFSingleConsumerExampleRaw struct {
	Contract *VRFSingleConsumerExample
}

type VRFSingleConsumerExampleCallerRaw struct {
	Contract *VRFSingleConsumerExampleCaller
}

type VRFSingleConsumerExampleTransactorRaw struct {
	Contract *VRFSingleConsumerExampleTransactor
}

func NewVRFSingleConsumerExample(address common.Address, backend bind.ContractBackend) (*VRFSingleConsumerExample, error) {
	abi, err := abi.JSON(strings.NewReader(VRFSingleConsumerExampleABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFSingleConsumerExample(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFSingleConsumerExample{address: address, abi: abi, VRFSingleConsumerExampleCaller: VRFSingleConsumerExampleCaller{contract: contract}, VRFSingleConsumerExampleTransactor: VRFSingleConsumerExampleTransactor{contract: contract}, VRFSingleConsumerExampleFilterer: VRFSingleConsumerExampleFilterer{contract: contract}}, nil
}

func NewVRFSingleConsumerExampleCaller(address common.Address, caller bind.ContractCaller) (*VRFSingleConsumerExampleCaller, error) {
	contract, err := bindVRFSingleConsumerExample(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFSingleConsumerExampleCaller{contract: contract}, nil
}

func NewVRFSingleConsumerExampleTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFSingleConsumerExampleTransactor, error) {
	contract, err := bindVRFSingleConsumerExample(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFSingleConsumerExampleTransactor{contract: contract}, nil
}

func NewVRFSingleConsumerExampleFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFSingleConsumerExampleFilterer, error) {
	contract, err := bindVRFSingleConsumerExample(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFSingleConsumerExampleFilterer{contract: contract}, nil
}

func bindVRFSingleConsumerExample(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFSingleConsumerExampleABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFSingleConsumerExample.Contract.VRFSingleConsumerExampleCaller.contract.Call(opts, result, method, params...)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFSingleConsumerExample.Contract.VRFSingleConsumerExampleTransactor.contract.Transfer(opts)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFSingleConsumerExample.Contract.VRFSingleConsumerExampleTransactor.contract.Transact(opts, method, params...)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFSingleConsumerExample.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFSingleConsumerExample.Contract.contract.Transfer(opts)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFSingleConsumerExample.Contract.contract.Transact(opts, method, params...)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleCaller) SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VRFSingleConsumerExample.contract.Call(opts, &out, "s_randomWords", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFSingleConsumerExample.Contract.SRandomWords(&_VRFSingleConsumerExample.CallOpts, arg0)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleCallerSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFSingleConsumerExample.Contract.SRandomWords(&_VRFSingleConsumerExample.CallOpts, arg0)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleCaller) SRequestConfig(opts *bind.CallOpts) (SRequestConfig,

	error) {
	var out []interface{}
	err := _VRFSingleConsumerExample.contract.Call(opts, &out, "s_requestConfig")

	outstruct := new(SRequestConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.SubId = *abi.ConvertType(out[0], new(uint64)).(*uint64)
	outstruct.CallbackGasLimit = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.RequestConfirmations = *abi.ConvertType(out[2], new(uint16)).(*uint16)
	outstruct.NumWords = *abi.ConvertType(out[3], new(uint32)).(*uint32)
	outstruct.KeyHash = *abi.ConvertType(out[4], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleSession) SRequestConfig() (SRequestConfig,

	error) {
	return _VRFSingleConsumerExample.Contract.SRequestConfig(&_VRFSingleConsumerExample.CallOpts)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleCallerSession) SRequestConfig() (SRequestConfig,

	error) {
	return _VRFSingleConsumerExample.Contract.SRequestConfig(&_VRFSingleConsumerExample.CallOpts)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleCaller) SRequestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFSingleConsumerExample.contract.Call(opts, &out, "s_requestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleSession) SRequestId() (*big.Int, error) {
	return _VRFSingleConsumerExample.Contract.SRequestId(&_VRFSingleConsumerExample.CallOpts)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleCallerSession) SRequestId() (*big.Int, error) {
	return _VRFSingleConsumerExample.Contract.SRequestId(&_VRFSingleConsumerExample.CallOpts)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFSingleConsumerExample.contract.Transact(opts, "rawFulfillRandomWords", requestId, randomWords)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFSingleConsumerExample.Contract.RawFulfillRandomWords(&_VRFSingleConsumerExample.TransactOpts, requestId, randomWords)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleTransactorSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFSingleConsumerExample.Contract.RawFulfillRandomWords(&_VRFSingleConsumerExample.TransactOpts, requestId, randomWords)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleTransactor) RequestRandomWords(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFSingleConsumerExample.contract.Transact(opts, "requestRandomWords")
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleSession) RequestRandomWords() (*types.Transaction, error) {
	return _VRFSingleConsumerExample.Contract.RequestRandomWords(&_VRFSingleConsumerExample.TransactOpts)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleTransactorSession) RequestRandomWords() (*types.Transaction, error) {
	return _VRFSingleConsumerExample.Contract.RequestRandomWords(&_VRFSingleConsumerExample.TransactOpts)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleTransactor) Subscribe(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFSingleConsumerExample.contract.Transact(opts, "subscribe")
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleSession) Subscribe() (*types.Transaction, error) {
	return _VRFSingleConsumerExample.Contract.Subscribe(&_VRFSingleConsumerExample.TransactOpts)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleTransactorSession) Subscribe() (*types.Transaction, error) {
	return _VRFSingleConsumerExample.Contract.Subscribe(&_VRFSingleConsumerExample.TransactOpts)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleTransactor) TopUpSubscription(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VRFSingleConsumerExample.contract.Transact(opts, "topUpSubscription", amount)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleSession) TopUpSubscription(amount *big.Int) (*types.Transaction, error) {
	return _VRFSingleConsumerExample.Contract.TopUpSubscription(&_VRFSingleConsumerExample.TransactOpts, amount)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleTransactorSession) TopUpSubscription(amount *big.Int) (*types.Transaction, error) {
	return _VRFSingleConsumerExample.Contract.TopUpSubscription(&_VRFSingleConsumerExample.TransactOpts, amount)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleTransactor) Unsubscribe(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFSingleConsumerExample.contract.Transact(opts, "unsubscribe")
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleSession) Unsubscribe() (*types.Transaction, error) {
	return _VRFSingleConsumerExample.Contract.Unsubscribe(&_VRFSingleConsumerExample.TransactOpts)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleTransactorSession) Unsubscribe() (*types.Transaction, error) {
	return _VRFSingleConsumerExample.Contract.Unsubscribe(&_VRFSingleConsumerExample.TransactOpts)
}

type SRequestConfig struct {
	SubId                uint64
	CallbackGasLimit     uint32
	RequestConfirmations uint16
	NumWords             uint32
	KeyHash              [32]byte
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExample) Address() common.Address {
	return _VRFSingleConsumerExample.address
}

type VRFSingleConsumerExampleInterface interface {
	SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	SRequestConfig(opts *bind.CallOpts) (SRequestConfig,

		error)

	SRequestId(opts *bind.CallOpts) (*big.Int, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	RequestRandomWords(opts *bind.TransactOpts) (*types.Transaction, error)

	Subscribe(opts *bind.TransactOpts) (*types.Transaction, error)

	TopUpSubscription(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	Unsubscribe(opts *bind.TransactOpts) (*types.Transaction, error)

	Address() common.Address
}
