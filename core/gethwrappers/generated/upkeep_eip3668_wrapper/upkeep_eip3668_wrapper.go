// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package upkeep_eip3668_wrapper

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

var UpkeepEIP3668MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"string[]\",\"name\":\"urls\",\"type\":\"string[]\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"},{\"internalType\":\"bytes4\",\"name\":\"callbackFunction\",\"type\":\"bytes4\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"OffchainLookup\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"initialBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"lastBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"previousBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"counter\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"resp\",\"type\":\"bytes\"}],\"name\":\"PerformingUpkeep\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"resp\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extra\",\"type\":\"bytes\"}],\"name\":\"callback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"interval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"previousPerformBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"input\",\"type\":\"string[]\"}],\"name\":\"setURLs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"testRange\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"urls\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162000f8b38038062000f8b833981016040819052620000349162000212565b6000828155600182905560038190554360025560048190556005556040805160a081018252602c918101828152909182919062000f5f606084013981526020016040518060600160405280602f815260200162000f30602f91399052620000a0906006906002620000a9565b50505062000274565b828054828255906000526020600020908101928215620000fb579160200282015b82811115620000fb5782518051620000ea9184916020909101906200010d565b5091602001919060010190620000ca565b506200010992915062000198565b5090565b8280546200011b9062000237565b90600052602060002090601f0160209004810192826200013f57600085556200018a565b82601f106200015a57805160ff19168380011785556200018a565b828001600101855582156200018a579182015b828111156200018a5782518255916020019190600101906200016d565b5062000109929150620001b9565b8082111562000109576000620001af8282620001d0565b5060010162000198565b5b80821115620001095760008155600101620001ba565b508054620001de9062000237565b6000825580601f10620001ef575050565b601f0160209004906000526020600020908101906200020f9190620001b9565b50565b600080604083850312156200022657600080fd5b505080516020909101519092909150565b600181811c908216806200024c57607f821691505b602082108114156200026e57634e487b7160e01b600052602260045260246000fd5b50919050565b610cac80620002846000396000f3fe608060405234801561001057600080fd5b50600436106100df5760003560e01c80636250a13a1161008c578063806b984f11610066578063806b984f146101b2578063917d895f146101bb578063947a36fb146101c4578063d832d92f146101cd57600080fd5b80636250a13a146101765780636e04ff0d1461017f578063796676be1461019257600080fd5b806339bc8e14116100bd57806339bc8e14146101475780634585e33b1461015a57806361bc221a1461016d57600080fd5b806318aa5141146100e45780631e34c5851461010e5780632cb1586414610130575b600080fd5b6100f76100f23660046107ea565b6101e5565b604051610105929190610af9565b60405180910390f35b61012e61011c36600461086f565b60009182556001556004819055600555565b005b61013960045481565b604051908152602001610105565b61012e610155366004610645565b610242565b61012e6101683660046107a8565b610259565b61013960055481565b61013960005481565b6100f761018d3660046107a8565b6102e3565b6101a56101a0366004610856565b6103c1565b6040516101059190610b30565b61013960025481565b61013960035481565b61013960015481565b6101d561046d565b6040519015158152602001610105565b60006060816101f68486018661077f565b905080878781818080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250959d929c50919a5050505050505050505050565b80516102559060069060208401906104b0565b5050565b60045461026557436004555b600554610273906001610bbe565b6005819055503273ffffffffffffffffffffffffffffffffffffffff167f955d4d39a5ccb8ff2d5a3d8508cb233a43a4f67441f1fe3b7d3a20007978811f60045460025460035460055487876040516102d196959493929190610b43565b60405180910390a25050600254600355565b6000606030600685856040516020016102fd929190610b1c565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815282825260016020840152917f18aa5141000000000000000000000000000000000000000000000000000000009101604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290527f556f18300000000000000000000000000000000000000000000000000000000082526103b89594939291600401610945565b60405180910390fd5b600681815481106103d157600080fd5b9060005260206000200160009150905080546103ec90610bed565b80601f016020809104026020016040519081016040528092919081815260200182805461041890610bed565b80156104655780601f1061043a57610100808354040283529160200191610465565b820191906000526020600020905b81548152906001019060200180831161044857829003601f168201915b505050505081565b6000600454600014156104805750600190565b6000546004546104909043610bd6565b1080156104ab57506001546002546104a89043610bd6565b10155b905090565b8280548282559060005260206000209081019282156104fd579160200282015b828111156104fd57825180516104ed91849160209091019061050d565b50916020019190600101906104d0565b5061050992915061058d565b5090565b82805461051990610bed565b90600052602060002090601f01602090048101928261053b5760008555610581565b82601f1061055457805160ff1916838001178555610581565b82800160010185558215610581579182015b82811115610581578251825591602001919060010190610566565b506105099291506105aa565b808211156105095760006105a182826105bf565b5060010161058d565b5b8082111561050957600081556001016105ab565b5080546105cb90610bed565b6000825580601f106105db575050565b601f0160209004906000526020600020908101906105f991906105aa565b50565b60008083601f84011261060e57600080fd5b50813567ffffffffffffffff81111561062657600080fd5b60208301915083602082850101111561063e57600080fd5b9250929050565b6000602080838503121561065857600080fd5b823567ffffffffffffffff8082111561067057600080fd5b8185019150601f868184011261068557600080fd5b82358281111561069757610697610c70565b8060051b6106a6868201610b6f565b8281528681019086880183880189018c10156106c157600080fd5b600093505b84841015610770578035878111156106dd57600080fd5b8801603f81018d136106ee57600080fd5b8981013560408982111561070457610704610c70565b6107338c7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08b85011601610b6f565b8281528f8284860101111561074757600080fd5b828285018e83013760009281018d019290925250845250600193909301929188019188016106c6565b509a9950505050505050505050565b60006020828403121561079157600080fd5b813580151581146107a157600080fd5b9392505050565b600080602083850312156107bb57600080fd5b823567ffffffffffffffff8111156107d257600080fd5b6107de858286016105fc565b90969095509350505050565b6000806000806040858703121561080057600080fd5b843567ffffffffffffffff8082111561081857600080fd5b610824888389016105fc565b9096509450602087013591508082111561083d57600080fd5b5061084a878288016105fc565b95989497509550505050565b60006020828403121561086857600080fd5b5035919050565b6000806040838503121561088257600080fd5b50508035926020909101359150565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b6000815180845260005b81811015610900576020818501810151868301820152016108e4565b81811115610912576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b600060a0820173ffffffffffffffffffffffffffffffffffffffff88168352602060a08185015281885480845260c08601915060c08160051b870101935060008a8152838120815b83811015610a99578887037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff4001855281548390600181811c90808316806109d557607f831692505b8a8310811415610a0c577f4e487b710000000000000000000000000000000000000000000000000000000088526022600452602488fd5b828c5260208c01818015610a275760018114610a5657610a80565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00861682528c82019650610a80565b6000898152602090208a5b86811015610a7a57815484820152908501908e01610a61565b83019750505b50949b505097890197949094019350505060010161098d565b5050505050508281036040840152610ab181876108da565b7fffffffff000000000000000000000000000000000000000000000000000000008616606085015290508281036080840152610aed81856108da565b98975050505050505050565b8215158152604060208201526000610b1460408301846108da565b949350505050565b602081526000610b14602083018486610891565b6020815260006107a160208301846108da565b86815285602082015284604082015283606082015260a060808201526000610aed60a083018486610891565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610bb657610bb6610c70565b604052919050565b60008219821115610bd157610bd1610c41565b500190565b600082821015610be857610be8610c41565b500390565b600181811c90821680610c0157607f821691505b60208210811415610c3b577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a68747470733a2f2f7777772e676f6f676c652e636f6d2f7365617263683f713d7b73656e6465727d2b7b646174617d68747470733a2f2f636174666163742e6e696e6a612f666163743f713d7b73656e6465727d2b7b646174617d",
}

var UpkeepEIP3668ABI = UpkeepEIP3668MetaData.ABI

var UpkeepEIP3668Bin = UpkeepEIP3668MetaData.Bin

func DeployUpkeepEIP3668(auth *bind.TransactOpts, backend bind.ContractBackend, _testRange *big.Int, _interval *big.Int) (common.Address, *types.Transaction, *UpkeepEIP3668, error) {
	parsed, err := UpkeepEIP3668MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(UpkeepEIP3668Bin), backend, _testRange, _interval)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &UpkeepEIP3668{UpkeepEIP3668Caller: UpkeepEIP3668Caller{contract: contract}, UpkeepEIP3668Transactor: UpkeepEIP3668Transactor{contract: contract}, UpkeepEIP3668Filterer: UpkeepEIP3668Filterer{contract: contract}}, nil
}

type UpkeepEIP3668 struct {
	address common.Address
	abi     abi.ABI
	UpkeepEIP3668Caller
	UpkeepEIP3668Transactor
	UpkeepEIP3668Filterer
}

type UpkeepEIP3668Caller struct {
	contract *bind.BoundContract
}

type UpkeepEIP3668Transactor struct {
	contract *bind.BoundContract
}

type UpkeepEIP3668Filterer struct {
	contract *bind.BoundContract
}

type UpkeepEIP3668Session struct {
	Contract     *UpkeepEIP3668
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type UpkeepEIP3668CallerSession struct {
	Contract *UpkeepEIP3668Caller
	CallOpts bind.CallOpts
}

type UpkeepEIP3668TransactorSession struct {
	Contract     *UpkeepEIP3668Transactor
	TransactOpts bind.TransactOpts
}

type UpkeepEIP3668Raw struct {
	Contract *UpkeepEIP3668
}

type UpkeepEIP3668CallerRaw struct {
	Contract *UpkeepEIP3668Caller
}

type UpkeepEIP3668TransactorRaw struct {
	Contract *UpkeepEIP3668Transactor
}

func NewUpkeepEIP3668(address common.Address, backend bind.ContractBackend) (*UpkeepEIP3668, error) {
	abi, err := abi.JSON(strings.NewReader(UpkeepEIP3668ABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindUpkeepEIP3668(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UpkeepEIP3668{address: address, abi: abi, UpkeepEIP3668Caller: UpkeepEIP3668Caller{contract: contract}, UpkeepEIP3668Transactor: UpkeepEIP3668Transactor{contract: contract}, UpkeepEIP3668Filterer: UpkeepEIP3668Filterer{contract: contract}}, nil
}

func NewUpkeepEIP3668Caller(address common.Address, caller bind.ContractCaller) (*UpkeepEIP3668Caller, error) {
	contract, err := bindUpkeepEIP3668(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UpkeepEIP3668Caller{contract: contract}, nil
}

func NewUpkeepEIP3668Transactor(address common.Address, transactor bind.ContractTransactor) (*UpkeepEIP3668Transactor, error) {
	contract, err := bindUpkeepEIP3668(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UpkeepEIP3668Transactor{contract: contract}, nil
}

func NewUpkeepEIP3668Filterer(address common.Address, filterer bind.ContractFilterer) (*UpkeepEIP3668Filterer, error) {
	contract, err := bindUpkeepEIP3668(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UpkeepEIP3668Filterer{contract: contract}, nil
}

func bindUpkeepEIP3668(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(UpkeepEIP3668ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_UpkeepEIP3668 *UpkeepEIP3668Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UpkeepEIP3668.Contract.UpkeepEIP3668Caller.contract.Call(opts, result, method, params...)
}

func (_UpkeepEIP3668 *UpkeepEIP3668Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpkeepEIP3668.Contract.UpkeepEIP3668Transactor.contract.Transfer(opts)
}

func (_UpkeepEIP3668 *UpkeepEIP3668Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UpkeepEIP3668.Contract.UpkeepEIP3668Transactor.contract.Transact(opts, method, params...)
}

func (_UpkeepEIP3668 *UpkeepEIP3668CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UpkeepEIP3668.Contract.contract.Call(opts, result, method, params...)
}

func (_UpkeepEIP3668 *UpkeepEIP3668TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpkeepEIP3668.Contract.contract.Transfer(opts)
}

func (_UpkeepEIP3668 *UpkeepEIP3668TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UpkeepEIP3668.Contract.contract.Transact(opts, method, params...)
}

func (_UpkeepEIP3668 *UpkeepEIP3668Caller) Callback(opts *bind.CallOpts, resp []byte, extra []byte) (bool, []byte, error) {
	var out []interface{}
	err := _UpkeepEIP3668.contract.Call(opts, &out, "callback", resp, extra)

	if err != nil {
		return *new(bool), *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return out0, out1, err

}

func (_UpkeepEIP3668 *UpkeepEIP3668Session) Callback(resp []byte, extra []byte) (bool, []byte, error) {
	return _UpkeepEIP3668.Contract.Callback(&_UpkeepEIP3668.CallOpts, resp, extra)
}

func (_UpkeepEIP3668 *UpkeepEIP3668CallerSession) Callback(resp []byte, extra []byte) (bool, []byte, error) {
	return _UpkeepEIP3668.Contract.Callback(&_UpkeepEIP3668.CallOpts, resp, extra)
}

func (_UpkeepEIP3668 *UpkeepEIP3668Caller) CheckUpkeep(opts *bind.CallOpts, data []byte) (bool, []byte, error) {
	var out []interface{}
	err := _UpkeepEIP3668.contract.Call(opts, &out, "checkUpkeep", data)

	if err != nil {
		return *new(bool), *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return out0, out1, err

}

func (_UpkeepEIP3668 *UpkeepEIP3668Session) CheckUpkeep(data []byte) (bool, []byte, error) {
	return _UpkeepEIP3668.Contract.CheckUpkeep(&_UpkeepEIP3668.CallOpts, data)
}

func (_UpkeepEIP3668 *UpkeepEIP3668CallerSession) CheckUpkeep(data []byte) (bool, []byte, error) {
	return _UpkeepEIP3668.Contract.CheckUpkeep(&_UpkeepEIP3668.CallOpts, data)
}

func (_UpkeepEIP3668 *UpkeepEIP3668Caller) Counter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepEIP3668.contract.Call(opts, &out, "counter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_UpkeepEIP3668 *UpkeepEIP3668Session) Counter() (*big.Int, error) {
	return _UpkeepEIP3668.Contract.Counter(&_UpkeepEIP3668.CallOpts)
}

func (_UpkeepEIP3668 *UpkeepEIP3668CallerSession) Counter() (*big.Int, error) {
	return _UpkeepEIP3668.Contract.Counter(&_UpkeepEIP3668.CallOpts)
}

func (_UpkeepEIP3668 *UpkeepEIP3668Caller) Eligible(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _UpkeepEIP3668.contract.Call(opts, &out, "eligible")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_UpkeepEIP3668 *UpkeepEIP3668Session) Eligible() (bool, error) {
	return _UpkeepEIP3668.Contract.Eligible(&_UpkeepEIP3668.CallOpts)
}

func (_UpkeepEIP3668 *UpkeepEIP3668CallerSession) Eligible() (bool, error) {
	return _UpkeepEIP3668.Contract.Eligible(&_UpkeepEIP3668.CallOpts)
}

func (_UpkeepEIP3668 *UpkeepEIP3668Caller) InitialBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepEIP3668.contract.Call(opts, &out, "initialBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_UpkeepEIP3668 *UpkeepEIP3668Session) InitialBlock() (*big.Int, error) {
	return _UpkeepEIP3668.Contract.InitialBlock(&_UpkeepEIP3668.CallOpts)
}

func (_UpkeepEIP3668 *UpkeepEIP3668CallerSession) InitialBlock() (*big.Int, error) {
	return _UpkeepEIP3668.Contract.InitialBlock(&_UpkeepEIP3668.CallOpts)
}

func (_UpkeepEIP3668 *UpkeepEIP3668Caller) Interval(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepEIP3668.contract.Call(opts, &out, "interval")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_UpkeepEIP3668 *UpkeepEIP3668Session) Interval() (*big.Int, error) {
	return _UpkeepEIP3668.Contract.Interval(&_UpkeepEIP3668.CallOpts)
}

func (_UpkeepEIP3668 *UpkeepEIP3668CallerSession) Interval() (*big.Int, error) {
	return _UpkeepEIP3668.Contract.Interval(&_UpkeepEIP3668.CallOpts)
}

func (_UpkeepEIP3668 *UpkeepEIP3668Caller) LastBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepEIP3668.contract.Call(opts, &out, "lastBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_UpkeepEIP3668 *UpkeepEIP3668Session) LastBlock() (*big.Int, error) {
	return _UpkeepEIP3668.Contract.LastBlock(&_UpkeepEIP3668.CallOpts)
}

func (_UpkeepEIP3668 *UpkeepEIP3668CallerSession) LastBlock() (*big.Int, error) {
	return _UpkeepEIP3668.Contract.LastBlock(&_UpkeepEIP3668.CallOpts)
}

func (_UpkeepEIP3668 *UpkeepEIP3668Caller) PreviousPerformBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepEIP3668.contract.Call(opts, &out, "previousPerformBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_UpkeepEIP3668 *UpkeepEIP3668Session) PreviousPerformBlock() (*big.Int, error) {
	return _UpkeepEIP3668.Contract.PreviousPerformBlock(&_UpkeepEIP3668.CallOpts)
}

func (_UpkeepEIP3668 *UpkeepEIP3668CallerSession) PreviousPerformBlock() (*big.Int, error) {
	return _UpkeepEIP3668.Contract.PreviousPerformBlock(&_UpkeepEIP3668.CallOpts)
}

func (_UpkeepEIP3668 *UpkeepEIP3668Caller) TestRange(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepEIP3668.contract.Call(opts, &out, "testRange")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_UpkeepEIP3668 *UpkeepEIP3668Session) TestRange() (*big.Int, error) {
	return _UpkeepEIP3668.Contract.TestRange(&_UpkeepEIP3668.CallOpts)
}

func (_UpkeepEIP3668 *UpkeepEIP3668CallerSession) TestRange() (*big.Int, error) {
	return _UpkeepEIP3668.Contract.TestRange(&_UpkeepEIP3668.CallOpts)
}

func (_UpkeepEIP3668 *UpkeepEIP3668Caller) Urls(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	var out []interface{}
	err := _UpkeepEIP3668.contract.Call(opts, &out, "urls", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_UpkeepEIP3668 *UpkeepEIP3668Session) Urls(arg0 *big.Int) (string, error) {
	return _UpkeepEIP3668.Contract.Urls(&_UpkeepEIP3668.CallOpts, arg0)
}

func (_UpkeepEIP3668 *UpkeepEIP3668CallerSession) Urls(arg0 *big.Int) (string, error) {
	return _UpkeepEIP3668.Contract.Urls(&_UpkeepEIP3668.CallOpts, arg0)
}

func (_UpkeepEIP3668 *UpkeepEIP3668Transactor) PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error) {
	return _UpkeepEIP3668.contract.Transact(opts, "performUpkeep", performData)
}

func (_UpkeepEIP3668 *UpkeepEIP3668Session) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _UpkeepEIP3668.Contract.PerformUpkeep(&_UpkeepEIP3668.TransactOpts, performData)
}

func (_UpkeepEIP3668 *UpkeepEIP3668TransactorSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _UpkeepEIP3668.Contract.PerformUpkeep(&_UpkeepEIP3668.TransactOpts, performData)
}

func (_UpkeepEIP3668 *UpkeepEIP3668Transactor) SetConfig(opts *bind.TransactOpts, _testRange *big.Int, _interval *big.Int) (*types.Transaction, error) {
	return _UpkeepEIP3668.contract.Transact(opts, "setConfig", _testRange, _interval)
}

func (_UpkeepEIP3668 *UpkeepEIP3668Session) SetConfig(_testRange *big.Int, _interval *big.Int) (*types.Transaction, error) {
	return _UpkeepEIP3668.Contract.SetConfig(&_UpkeepEIP3668.TransactOpts, _testRange, _interval)
}

func (_UpkeepEIP3668 *UpkeepEIP3668TransactorSession) SetConfig(_testRange *big.Int, _interval *big.Int) (*types.Transaction, error) {
	return _UpkeepEIP3668.Contract.SetConfig(&_UpkeepEIP3668.TransactOpts, _testRange, _interval)
}

func (_UpkeepEIP3668 *UpkeepEIP3668Transactor) SetURLs(opts *bind.TransactOpts, input []string) (*types.Transaction, error) {
	return _UpkeepEIP3668.contract.Transact(opts, "setURLs", input)
}

func (_UpkeepEIP3668 *UpkeepEIP3668Session) SetURLs(input []string) (*types.Transaction, error) {
	return _UpkeepEIP3668.Contract.SetURLs(&_UpkeepEIP3668.TransactOpts, input)
}

func (_UpkeepEIP3668 *UpkeepEIP3668TransactorSession) SetURLs(input []string) (*types.Transaction, error) {
	return _UpkeepEIP3668.Contract.SetURLs(&_UpkeepEIP3668.TransactOpts, input)
}

type UpkeepEIP3668PerformingUpkeepIterator struct {
	Event *UpkeepEIP3668PerformingUpkeep

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *UpkeepEIP3668PerformingUpkeepIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpkeepEIP3668PerformingUpkeep)
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
		it.Event = new(UpkeepEIP3668PerformingUpkeep)
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

func (it *UpkeepEIP3668PerformingUpkeepIterator) Error() error {
	return it.fail
}

func (it *UpkeepEIP3668PerformingUpkeepIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type UpkeepEIP3668PerformingUpkeep struct {
	From          common.Address
	InitialBlock  *big.Int
	LastBlock     *big.Int
	PreviousBlock *big.Int
	Counter       *big.Int
	Resp          []byte
	Raw           types.Log
}

func (_UpkeepEIP3668 *UpkeepEIP3668Filterer) FilterPerformingUpkeep(opts *bind.FilterOpts, from []common.Address) (*UpkeepEIP3668PerformingUpkeepIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _UpkeepEIP3668.contract.FilterLogs(opts, "PerformingUpkeep", fromRule)
	if err != nil {
		return nil, err
	}
	return &UpkeepEIP3668PerformingUpkeepIterator{contract: _UpkeepEIP3668.contract, event: "PerformingUpkeep", logs: logs, sub: sub}, nil
}

func (_UpkeepEIP3668 *UpkeepEIP3668Filterer) WatchPerformingUpkeep(opts *bind.WatchOpts, sink chan<- *UpkeepEIP3668PerformingUpkeep, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _UpkeepEIP3668.contract.WatchLogs(opts, "PerformingUpkeep", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(UpkeepEIP3668PerformingUpkeep)
				if err := _UpkeepEIP3668.contract.UnpackLog(event, "PerformingUpkeep", log); err != nil {
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

func (_UpkeepEIP3668 *UpkeepEIP3668Filterer) ParsePerformingUpkeep(log types.Log) (*UpkeepEIP3668PerformingUpkeep, error) {
	event := new(UpkeepEIP3668PerformingUpkeep)
	if err := _UpkeepEIP3668.contract.UnpackLog(event, "PerformingUpkeep", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_UpkeepEIP3668 *UpkeepEIP3668) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _UpkeepEIP3668.abi.Events["PerformingUpkeep"].ID:
		return _UpkeepEIP3668.ParsePerformingUpkeep(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (UpkeepEIP3668PerformingUpkeep) Topic() common.Hash {
	return common.HexToHash("0x955d4d39a5ccb8ff2d5a3d8508cb233a43a4f67441f1fe3b7d3a20007978811f")
}

func (_UpkeepEIP3668 *UpkeepEIP3668) Address() common.Address {
	return _UpkeepEIP3668.address
}

type UpkeepEIP3668Interface interface {
	Callback(opts *bind.CallOpts, resp []byte, extra []byte) (bool, []byte, error)

	CheckUpkeep(opts *bind.CallOpts, data []byte) (bool, []byte, error)

	Counter(opts *bind.CallOpts) (*big.Int, error)

	Eligible(opts *bind.CallOpts) (bool, error)

	InitialBlock(opts *bind.CallOpts) (*big.Int, error)

	Interval(opts *bind.CallOpts) (*big.Int, error)

	LastBlock(opts *bind.CallOpts) (*big.Int, error)

	PreviousPerformBlock(opts *bind.CallOpts) (*big.Int, error)

	TestRange(opts *bind.CallOpts) (*big.Int, error)

	Urls(opts *bind.CallOpts, arg0 *big.Int) (string, error)

	PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, _testRange *big.Int, _interval *big.Int) (*types.Transaction, error)

	SetURLs(opts *bind.TransactOpts, input []string) (*types.Transaction, error)

	FilterPerformingUpkeep(opts *bind.FilterOpts, from []common.Address) (*UpkeepEIP3668PerformingUpkeepIterator, error)

	WatchPerformingUpkeep(opts *bind.WatchOpts, sink chan<- *UpkeepEIP3668PerformingUpkeep, from []common.Address) (event.Subscription, error)

	ParsePerformingUpkeep(log types.Log) (*UpkeepEIP3668PerformingUpkeep, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
