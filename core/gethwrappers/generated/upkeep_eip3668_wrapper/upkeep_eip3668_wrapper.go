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
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"string[]\",\"name\":\"urls\",\"type\":\"string[]\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"},{\"internalType\":\"bytes4\",\"name\":\"callbackFunction\",\"type\":\"bytes4\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"OffchainLookup\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"initialBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"lastBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"previousBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"counter\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"resp\",\"type\":\"uint256\"}],\"name\":\"PerformingUpkeep\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"resp\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extra\",\"type\":\"bytes\"}],\"name\":\"callback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"interval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"previousPerformBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"input\",\"type\":\"string[]\"}],\"name\":\"setURLs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"testRange\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"urls\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162000f6438038062000f648339810160408190526200003491620001f4565b60008281556001829055600381905543600255600481905560055560408051608081018252602f6020820181815291928392919062000f35908401399052620000829060069060016200008b565b50505062000256565b828054828255906000526020600020908101928215620000dd579160200282015b82811115620000dd5782518051620000cc918491602090910190620000ef565b5091602001919060010190620000ac565b50620000eb9291506200017a565b5090565b828054620000fd9062000219565b90600052602060002090601f0160209004810192826200012157600085556200016c565b82601f106200013c57805160ff19168380011785556200016c565b828001600101855582156200016c579182015b828111156200016c5782518255916020019190600101906200014f565b50620000eb9291506200019b565b80821115620000eb576000620001918282620001b2565b506001016200017a565b5b80821115620000eb57600081556001016200019c565b508054620001c09062000219565b6000825580601f10620001d1575050565b601f016020900490600052602060002090810190620001f191906200019b565b50565b600080604083850312156200020857600080fd5b505080516020909101519092909150565b600181811c908216806200022e57607f821691505b602082108114156200025057634e487b7160e01b600052602260045260246000fd5b50919050565b610ccf80620002666000396000f3fe608060405234801561001057600080fd5b50600436106100df5760003560e01c80636250a13a1161008c578063806b984f11610066578063806b984f146101b2578063917d895f146101bb578063947a36fb146101c4578063d832d92f146101cd57600080fd5b80636250a13a146101765780636e04ff0d1461017f578063796676be1461019257600080fd5b806339bc8e14116100bd57806339bc8e14146101475780634585e33b1461015a57806361bc221a1461016d57600080fd5b806318aa5141146100e45780631e34c5851461010e5780632cb1586414610130575b600080fd5b6100f76100f2366004610849565b6101e5565b604051610105929190610b0f565b60405180910390f35b61012e61011c3660046108ce565b60009182556001556004819055600555565b005b61013960045481565b604051908152602001610105565b61012e6101553660046106a4565b610242565b61012e610168366004610807565b610259565b61013960055481565b61013960005481565b6100f761018d366004610807565b6102f2565b6101a56101a03660046108b5565b610427565b6040516101059190610b7f565b61013960025481565b61013960035481565b61013960015481565b6101d56104d3565b6040519015158152602001610105565b60006060816101f6848601866107de565b905080878781818080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250959d929c50919a5050505050505050505050565b8051610255906006906020840190610516565b5050565b60045461026557436004555b6000610273828401846108b5565b43600255600554909150610288906001610be1565b60058190556004546002546003546040805193845260208401929092529082015260608101919091526080810182905232907f4874b8dd61a40fe23599b4360a9a824d7081742fca9f555bcee3d389c4f4bd659060a00160405180910390a2505060025460035550565b600060606102fe6104d3565b156103dd573060068585604051602001610319929190610b32565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815282825260016020840152917f18aa5141000000000000000000000000000000000000000000000000000000009101604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290527f556f18300000000000000000000000000000000000000000000000000000000082526103d4959493929160040161095b565b60405180910390fd5b6000848481818080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525095975091955050505050505b9250929050565b6006818154811061043757600080fd5b90600052602060002001600091509050805461045290610c10565b80601f016020809104026020016040519081016040528092919081815260200182805461047e90610c10565b80156104cb5780601f106104a0576101008083540402835291602001916104cb565b820191906000526020600020905b8154815290600101906020018083116104ae57829003601f168201915b505050505081565b6000600454600014156104e65750600190565b6000546004546104f69043610bf9565b108015610511575060015460025461050e9043610bf9565b10155b905090565b828054828255906000526020600020908101928215610563579160200282015b828111156105635782518051610553918491602090910190610573565b5091602001919060010190610536565b5061056f9291506105f3565b5090565b82805461057f90610c10565b90600052602060002090601f0160209004810192826105a157600085556105e7565b82601f106105ba57805160ff19168380011785556105e7565b828001600101855582156105e7579182015b828111156105e75782518255916020019190600101906105cc565b5061056f929150610610565b8082111561056f5760006106078282610625565b506001016105f3565b5b8082111561056f5760008155600101610611565b50805461063190610c10565b6000825580601f10610641575050565b601f01602090049060005260206000209081019061065f9190610610565b50565b60008083601f84011261067457600080fd5b50813567ffffffffffffffff81111561068c57600080fd5b60208301915083602082850101111561042057600080fd5b600060208083850312156106b757600080fd5b823567ffffffffffffffff808211156106cf57600080fd5b8185019150601f86818401126106e457600080fd5b8235828111156106f6576106f6610c93565b8060051b610705868201610b92565b8281528681019086880183880189018c101561072057600080fd5b600093505b848410156107cf5780358781111561073c57600080fd5b8801603f81018d1361074d57600080fd5b8981013560408982111561076357610763610c93565b6107928c7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08b85011601610b92565b8281528f828486010111156107a657600080fd5b828285018e83013760009281018d01929092525084525060019390930192918801918801610725565b509a9950505050505050505050565b6000602082840312156107f057600080fd5b8135801515811461080057600080fd5b9392505050565b6000806020838503121561081a57600080fd5b823567ffffffffffffffff81111561083157600080fd5b61083d85828601610662565b90969095509350505050565b6000806000806040858703121561085f57600080fd5b843567ffffffffffffffff8082111561087757600080fd5b61088388838901610662565b9096509450602087013591508082111561089c57600080fd5b506108a987828801610662565b95989497509550505050565b6000602082840312156108c757600080fd5b5035919050565b600080604083850312156108e157600080fd5b50508035926020909101359150565b6000815180845260005b81811015610916576020818501810151868301820152016108fa565b81811115610928576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b600060a0820173ffffffffffffffffffffffffffffffffffffffff88168352602060a08185015281885480845260c08601915060c08160051b870101935060008a8152838120815b83811015610aaf578887037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff4001855281548390600181811c90808316806109eb57607f831692505b8a8310811415610a22577f4e487b710000000000000000000000000000000000000000000000000000000088526022600452602488fd5b828c5260208c01818015610a3d5760018114610a6c57610a96565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00861682528c82019650610a96565b6000898152602090208a5b86811015610a9057815484820152908501908e01610a77565b83019750505b50949b50509789019794909401935050506001016109a3565b5050505050508281036040840152610ac781876108f0565b7fffffffff000000000000000000000000000000000000000000000000000000008616606085015290508281036080840152610b0381856108f0565b98975050505050505050565b8215158152604060208201526000610b2a60408301846108f0565b949350505050565b60208152816020820152818360408301376000818301604090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0160101919050565b60208152600061080060208301846108f0565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610bd957610bd9610c93565b604052919050565b60008219821115610bf457610bf4610c64565b500190565b600082821015610c0b57610c0b610c64565b500390565b600181811c90821680610c2457607f821691505b60208210811415610c5e577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a68747470733a2f2f7777772e676f6f676c652e636f6d2f7365617263683f713d7b73656e6465727d2b7b646174617d",
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
	Resp          *big.Int
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
	return common.HexToHash("0x4874b8dd61a40fe23599b4360a9a824d7081742fca9f555bcee3d389c4f4bd65")
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
