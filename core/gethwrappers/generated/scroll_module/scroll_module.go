// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package scroll_module

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
	_ = abi.ConvertType
)

var ScrollModuleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"blockHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"blockNumber\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCurrentL1Fee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getGasOverhead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"chainModuleFixedOverhead\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"chainModulePerByteOverhead\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"dataSize\",\"type\":\"uint256\"}],\"name\":\"getMaxL1Fee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610403806100206000396000f3fe608060405234801561001057600080fd5b50600436106100675760003560e01c806357e871e71161005057806357e871e71461009a57806385df51fd146100a0578063de9ee35e146100b357600080fd5b8063125441401461006c57806318b8f61314610092575b600080fd5b61007f61007a36600461027d565b6100c9565b6040519081526020015b60405180910390f35b61007f6101b5565b4361007f565b61007f6100ae36600461027d565b610250565b60408051613a9881526014602082015201610089565b6000806100d78360046102c5565b67ffffffffffffffff8111156100ef576100ef6102e2565b6040519080825280601f01601f191660200182016040528015610119576020820181803683370190505b506040517f49948e0e000000000000000000000000000000000000000000000000000000008152909150735300000000000000000000000000000000000002906349948e0e9061016d908490600401610311565b602060405180830381865afa15801561018a573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101ae919061037d565b9392505050565b6040517f49948e0e000000000000000000000000000000000000000000000000000000008152600090735300000000000000000000000000000000000002906349948e0e9061020a9084903690600401610396565b602060405180830381865afa158015610227573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061024b919061037d565b905090565b6000438210158061026b575061010061026983436103e3565b115b1561027857506000919050565b504090565b60006020828403121561028f57600080fd5b5035919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b80820281158282048414176102dc576102dc610296565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600060208083528351808285015260005b8181101561033e57858101830151858201604001528201610322565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b60006020828403121561038f57600080fd5b5051919050565b60208152816020820152818360408301376000818301604090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0160101919050565b818103818111156102dc576102dc61029656fea164736f6c6343000813000a",
}

var ScrollModuleABI = ScrollModuleMetaData.ABI

var ScrollModuleBin = ScrollModuleMetaData.Bin

func DeployScrollModule(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ScrollModule, error) {
	parsed, err := ScrollModuleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ScrollModuleBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ScrollModule{address: address, abi: *parsed, ScrollModuleCaller: ScrollModuleCaller{contract: contract}, ScrollModuleTransactor: ScrollModuleTransactor{contract: contract}, ScrollModuleFilterer: ScrollModuleFilterer{contract: contract}}, nil
}

type ScrollModule struct {
	address common.Address
	abi     abi.ABI
	ScrollModuleCaller
	ScrollModuleTransactor
	ScrollModuleFilterer
}

type ScrollModuleCaller struct {
	contract *bind.BoundContract
}

type ScrollModuleTransactor struct {
	contract *bind.BoundContract
}

type ScrollModuleFilterer struct {
	contract *bind.BoundContract
}

type ScrollModuleSession struct {
	Contract     *ScrollModule
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type ScrollModuleCallerSession struct {
	Contract *ScrollModuleCaller
	CallOpts bind.CallOpts
}

type ScrollModuleTransactorSession struct {
	Contract     *ScrollModuleTransactor
	TransactOpts bind.TransactOpts
}

type ScrollModuleRaw struct {
	Contract *ScrollModule
}

type ScrollModuleCallerRaw struct {
	Contract *ScrollModuleCaller
}

type ScrollModuleTransactorRaw struct {
	Contract *ScrollModuleTransactor
}

func NewScrollModule(address common.Address, backend bind.ContractBackend) (*ScrollModule, error) {
	abi, err := abi.JSON(strings.NewReader(ScrollModuleABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindScrollModule(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ScrollModule{address: address, abi: abi, ScrollModuleCaller: ScrollModuleCaller{contract: contract}, ScrollModuleTransactor: ScrollModuleTransactor{contract: contract}, ScrollModuleFilterer: ScrollModuleFilterer{contract: contract}}, nil
}

func NewScrollModuleCaller(address common.Address, caller bind.ContractCaller) (*ScrollModuleCaller, error) {
	contract, err := bindScrollModule(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ScrollModuleCaller{contract: contract}, nil
}

func NewScrollModuleTransactor(address common.Address, transactor bind.ContractTransactor) (*ScrollModuleTransactor, error) {
	contract, err := bindScrollModule(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ScrollModuleTransactor{contract: contract}, nil
}

func NewScrollModuleFilterer(address common.Address, filterer bind.ContractFilterer) (*ScrollModuleFilterer, error) {
	contract, err := bindScrollModule(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ScrollModuleFilterer{contract: contract}, nil
}

func bindScrollModule(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ScrollModuleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_ScrollModule *ScrollModuleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ScrollModule.Contract.ScrollModuleCaller.contract.Call(opts, result, method, params...)
}

func (_ScrollModule *ScrollModuleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ScrollModule.Contract.ScrollModuleTransactor.contract.Transfer(opts)
}

func (_ScrollModule *ScrollModuleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ScrollModule.Contract.ScrollModuleTransactor.contract.Transact(opts, method, params...)
}

func (_ScrollModule *ScrollModuleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ScrollModule.Contract.contract.Call(opts, result, method, params...)
}

func (_ScrollModule *ScrollModuleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ScrollModule.Contract.contract.Transfer(opts)
}

func (_ScrollModule *ScrollModuleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ScrollModule.Contract.contract.Transact(opts, method, params...)
}

func (_ScrollModule *ScrollModuleCaller) BlockHash(opts *bind.CallOpts, n *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _ScrollModule.contract.Call(opts, &out, "blockHash", n)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_ScrollModule *ScrollModuleSession) BlockHash(n *big.Int) ([32]byte, error) {
	return _ScrollModule.Contract.BlockHash(&_ScrollModule.CallOpts, n)
}

func (_ScrollModule *ScrollModuleCallerSession) BlockHash(n *big.Int) ([32]byte, error) {
	return _ScrollModule.Contract.BlockHash(&_ScrollModule.CallOpts, n)
}

func (_ScrollModule *ScrollModuleCaller) BlockNumber(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ScrollModule.contract.Call(opts, &out, "blockNumber")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_ScrollModule *ScrollModuleSession) BlockNumber() (*big.Int, error) {
	return _ScrollModule.Contract.BlockNumber(&_ScrollModule.CallOpts)
}

func (_ScrollModule *ScrollModuleCallerSession) BlockNumber() (*big.Int, error) {
	return _ScrollModule.Contract.BlockNumber(&_ScrollModule.CallOpts)
}

func (_ScrollModule *ScrollModuleCaller) GetCurrentL1Fee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ScrollModule.contract.Call(opts, &out, "getCurrentL1Fee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_ScrollModule *ScrollModuleSession) GetCurrentL1Fee() (*big.Int, error) {
	return _ScrollModule.Contract.GetCurrentL1Fee(&_ScrollModule.CallOpts)
}

func (_ScrollModule *ScrollModuleCallerSession) GetCurrentL1Fee() (*big.Int, error) {
	return _ScrollModule.Contract.GetCurrentL1Fee(&_ScrollModule.CallOpts)
}

func (_ScrollModule *ScrollModuleCaller) GetGasOverhead(opts *bind.CallOpts) (GetGasOverhead,

	error) {
	var out []interface{}
	err := _ScrollModule.contract.Call(opts, &out, "getGasOverhead")

	outstruct := new(GetGasOverhead)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ChainModuleFixedOverhead = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.ChainModulePerByteOverhead = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_ScrollModule *ScrollModuleSession) GetGasOverhead() (GetGasOverhead,

	error) {
	return _ScrollModule.Contract.GetGasOverhead(&_ScrollModule.CallOpts)
}

func (_ScrollModule *ScrollModuleCallerSession) GetGasOverhead() (GetGasOverhead,

	error) {
	return _ScrollModule.Contract.GetGasOverhead(&_ScrollModule.CallOpts)
}

func (_ScrollModule *ScrollModuleCaller) GetMaxL1Fee(opts *bind.CallOpts, dataSize *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _ScrollModule.contract.Call(opts, &out, "getMaxL1Fee", dataSize)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_ScrollModule *ScrollModuleSession) GetMaxL1Fee(dataSize *big.Int) (*big.Int, error) {
	return _ScrollModule.Contract.GetMaxL1Fee(&_ScrollModule.CallOpts, dataSize)
}

func (_ScrollModule *ScrollModuleCallerSession) GetMaxL1Fee(dataSize *big.Int) (*big.Int, error) {
	return _ScrollModule.Contract.GetMaxL1Fee(&_ScrollModule.CallOpts, dataSize)
}

type GetGasOverhead struct {
	ChainModuleFixedOverhead   *big.Int
	ChainModulePerByteOverhead *big.Int
}

func (_ScrollModule *ScrollModule) Address() common.Address {
	return _ScrollModule.address
}

type ScrollModuleInterface interface {
	BlockHash(opts *bind.CallOpts, n *big.Int) ([32]byte, error)

	BlockNumber(opts *bind.CallOpts) (*big.Int, error)

	GetCurrentL1Fee(opts *bind.CallOpts) (*big.Int, error)

	GetGasOverhead(opts *bind.CallOpts) (GetGasOverhead,

		error)

	GetMaxL1Fee(opts *bind.CallOpts, dataSize *big.Int) (*big.Int, error)

	Address() common.Address
}
