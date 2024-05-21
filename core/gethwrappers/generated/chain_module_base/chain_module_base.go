// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package chain_module_base

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

var ChainModuleBaseMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"blockHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"blockNumber\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCurrentL1Fee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getGasOverhead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"chainModuleFixedOverhead\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"chainModulePerByteOverhead\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"getMaxL1Fee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5061015c806100206000396000f3fe608060405234801561001057600080fd5b50600436106100675760003560e01c806357e871e71161005057806357e871e71461009a57806385df51fd146100a0578063de9ee35e146100b357600080fd5b8063125441401461006c57806318b8f61314610093575b600080fd5b61008061007a3660046100f6565b50600090565b6040519081526020015b60405180910390f35b6000610080565b43610080565b6100806100ae3660046100f6565b6100c9565b6040805161012c8152600060208201520161008a565b600043821015806100e457506101006100e2834361010f565b115b156100f157506000919050565b504090565b60006020828403121561010857600080fd5b5035919050565b81810381811115610149577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b9291505056fea164736f6c6343000813000a",
}

var ChainModuleBaseABI = ChainModuleBaseMetaData.ABI

var ChainModuleBaseBin = ChainModuleBaseMetaData.Bin

func DeployChainModuleBase(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ChainModuleBase, error) {
	parsed, err := ChainModuleBaseMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ChainModuleBaseBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ChainModuleBase{address: address, abi: *parsed, ChainModuleBaseCaller: ChainModuleBaseCaller{contract: contract}, ChainModuleBaseTransactor: ChainModuleBaseTransactor{contract: contract}, ChainModuleBaseFilterer: ChainModuleBaseFilterer{contract: contract}}, nil
}

type ChainModuleBase struct {
	address common.Address
	abi     abi.ABI
	ChainModuleBaseCaller
	ChainModuleBaseTransactor
	ChainModuleBaseFilterer
}

type ChainModuleBaseCaller struct {
	contract *bind.BoundContract
}

type ChainModuleBaseTransactor struct {
	contract *bind.BoundContract
}

type ChainModuleBaseFilterer struct {
	contract *bind.BoundContract
}

type ChainModuleBaseSession struct {
	Contract     *ChainModuleBase
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type ChainModuleBaseCallerSession struct {
	Contract *ChainModuleBaseCaller
	CallOpts bind.CallOpts
}

type ChainModuleBaseTransactorSession struct {
	Contract     *ChainModuleBaseTransactor
	TransactOpts bind.TransactOpts
}

type ChainModuleBaseRaw struct {
	Contract *ChainModuleBase
}

type ChainModuleBaseCallerRaw struct {
	Contract *ChainModuleBaseCaller
}

type ChainModuleBaseTransactorRaw struct {
	Contract *ChainModuleBaseTransactor
}

func NewChainModuleBase(address common.Address, backend bind.ContractBackend) (*ChainModuleBase, error) {
	abi, err := abi.JSON(strings.NewReader(ChainModuleBaseABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindChainModuleBase(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ChainModuleBase{address: address, abi: abi, ChainModuleBaseCaller: ChainModuleBaseCaller{contract: contract}, ChainModuleBaseTransactor: ChainModuleBaseTransactor{contract: contract}, ChainModuleBaseFilterer: ChainModuleBaseFilterer{contract: contract}}, nil
}

func NewChainModuleBaseCaller(address common.Address, caller bind.ContractCaller) (*ChainModuleBaseCaller, error) {
	contract, err := bindChainModuleBase(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ChainModuleBaseCaller{contract: contract}, nil
}

func NewChainModuleBaseTransactor(address common.Address, transactor bind.ContractTransactor) (*ChainModuleBaseTransactor, error) {
	contract, err := bindChainModuleBase(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ChainModuleBaseTransactor{contract: contract}, nil
}

func NewChainModuleBaseFilterer(address common.Address, filterer bind.ContractFilterer) (*ChainModuleBaseFilterer, error) {
	contract, err := bindChainModuleBase(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ChainModuleBaseFilterer{contract: contract}, nil
}

func bindChainModuleBase(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ChainModuleBaseMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_ChainModuleBase *ChainModuleBaseRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ChainModuleBase.Contract.ChainModuleBaseCaller.contract.Call(opts, result, method, params...)
}

func (_ChainModuleBase *ChainModuleBaseRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChainModuleBase.Contract.ChainModuleBaseTransactor.contract.Transfer(opts)
}

func (_ChainModuleBase *ChainModuleBaseRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ChainModuleBase.Contract.ChainModuleBaseTransactor.contract.Transact(opts, method, params...)
}

func (_ChainModuleBase *ChainModuleBaseCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ChainModuleBase.Contract.contract.Call(opts, result, method, params...)
}

func (_ChainModuleBase *ChainModuleBaseTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChainModuleBase.Contract.contract.Transfer(opts)
}

func (_ChainModuleBase *ChainModuleBaseTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ChainModuleBase.Contract.contract.Transact(opts, method, params...)
}

func (_ChainModuleBase *ChainModuleBaseCaller) BlockHash(opts *bind.CallOpts, n *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _ChainModuleBase.contract.Call(opts, &out, "blockHash", n)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_ChainModuleBase *ChainModuleBaseSession) BlockHash(n *big.Int) ([32]byte, error) {
	return _ChainModuleBase.Contract.BlockHash(&_ChainModuleBase.CallOpts, n)
}

func (_ChainModuleBase *ChainModuleBaseCallerSession) BlockHash(n *big.Int) ([32]byte, error) {
	return _ChainModuleBase.Contract.BlockHash(&_ChainModuleBase.CallOpts, n)
}

func (_ChainModuleBase *ChainModuleBaseCaller) BlockNumber(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ChainModuleBase.contract.Call(opts, &out, "blockNumber")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_ChainModuleBase *ChainModuleBaseSession) BlockNumber() (*big.Int, error) {
	return _ChainModuleBase.Contract.BlockNumber(&_ChainModuleBase.CallOpts)
}

func (_ChainModuleBase *ChainModuleBaseCallerSession) BlockNumber() (*big.Int, error) {
	return _ChainModuleBase.Contract.BlockNumber(&_ChainModuleBase.CallOpts)
}

func (_ChainModuleBase *ChainModuleBaseCaller) GetCurrentL1Fee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ChainModuleBase.contract.Call(opts, &out, "getCurrentL1Fee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_ChainModuleBase *ChainModuleBaseSession) GetCurrentL1Fee() (*big.Int, error) {
	return _ChainModuleBase.Contract.GetCurrentL1Fee(&_ChainModuleBase.CallOpts)
}

func (_ChainModuleBase *ChainModuleBaseCallerSession) GetCurrentL1Fee() (*big.Int, error) {
	return _ChainModuleBase.Contract.GetCurrentL1Fee(&_ChainModuleBase.CallOpts)
}

func (_ChainModuleBase *ChainModuleBaseCaller) GetGasOverhead(opts *bind.CallOpts) (GetGasOverhead,

	error) {
	var out []interface{}
	err := _ChainModuleBase.contract.Call(opts, &out, "getGasOverhead")

	outstruct := new(GetGasOverhead)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ChainModuleFixedOverhead = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.ChainModulePerByteOverhead = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_ChainModuleBase *ChainModuleBaseSession) GetGasOverhead() (GetGasOverhead,

	error) {
	return _ChainModuleBase.Contract.GetGasOverhead(&_ChainModuleBase.CallOpts)
}

func (_ChainModuleBase *ChainModuleBaseCallerSession) GetGasOverhead() (GetGasOverhead,

	error) {
	return _ChainModuleBase.Contract.GetGasOverhead(&_ChainModuleBase.CallOpts)
}

func (_ChainModuleBase *ChainModuleBaseCaller) GetMaxL1Fee(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _ChainModuleBase.contract.Call(opts, &out, "getMaxL1Fee", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_ChainModuleBase *ChainModuleBaseSession) GetMaxL1Fee(arg0 *big.Int) (*big.Int, error) {
	return _ChainModuleBase.Contract.GetMaxL1Fee(&_ChainModuleBase.CallOpts, arg0)
}

func (_ChainModuleBase *ChainModuleBaseCallerSession) GetMaxL1Fee(arg0 *big.Int) (*big.Int, error) {
	return _ChainModuleBase.Contract.GetMaxL1Fee(&_ChainModuleBase.CallOpts, arg0)
}

type GetGasOverhead struct {
	ChainModuleFixedOverhead   *big.Int
	ChainModulePerByteOverhead *big.Int
}

func (_ChainModuleBase *ChainModuleBase) Address() common.Address {
	return _ChainModuleBase.address
}

type ChainModuleBaseInterface interface {
	BlockHash(opts *bind.CallOpts, n *big.Int) ([32]byte, error)

	BlockNumber(opts *bind.CallOpts) (*big.Int, error)

	GetCurrentL1Fee(opts *bind.CallOpts) (*big.Int, error)

	GetGasOverhead(opts *bind.CallOpts) (GetGasOverhead,

		error)

	GetMaxL1Fee(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	Address() common.Address
}
