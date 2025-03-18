// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package arbitrum_module

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

var ArbitrumModuleMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"blockHash\",\"inputs\":[{\"name\":\"n\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"blockNumber\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCurrentL1Fee\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getGasOverhead\",\"inputs\":[],\"outputs\":[{\"name\":\"chainModuleFixedOverhead\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"chainModulePerByteOverhead\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMaxL1Fee\",\"inputs\":[{\"name\":\"dataSize\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"}]",
	Bin: "0x608060405234801561001057600080fd5b5061044a806100206000396000f3fe608060405234801561001057600080fd5b50600436106100675760003560e01c80637810d12a116100505780637810d12a1461009a57806385df51fd146100ad578063de9ee35e146100c057600080fd5b8063125441401461006c57806357e871e714610092575b600080fd5b61007f61007a366004610368565b6100d6565b6040519081526020015b60405180910390f35b61007f610163565b61007f6100a8366004610368565b6101da565b61007f6100bb366004610368565b610252565b6040805161138881526000602082015201610089565b600080606c73ffffffffffffffffffffffffffffffffffffffff166341b247a86040518163ffffffff1660e01b815260040160c060405180830381865afa158015610125573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101499190610381565b50505050915050828161015c91906103fa565b9392505050565b6000606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156101b1573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101d59190610411565b905090565b6000606c73ffffffffffffffffffffffffffffffffffffffff1663c6f7de0e6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610228573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061024c9190610411565b92915050565b600080606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156102a1573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102c59190610411565b905080831015806102e057506101006102de848361042a565b115b156102ee5750600092915050565b6040517f2b407a8200000000000000000000000000000000000000000000000000000000815260048101849052606490632b407a8290602401602060405180830381865afa158015610344573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061015c9190610411565b60006020828403121561037a57600080fd5b5035919050565b60008060008060008060c0878903121561039a57600080fd5b865195506020870151945060408701519350606087015192506080870151915060a087015190509295509295509295565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b808202811582820484141761024c5761024c6103cb565b60006020828403121561042357600080fd5b5051919050565b8181038181111561024c5761024c6103cb56fea164736f6c6343000813000a",
}

var ArbitrumModuleABI = ArbitrumModuleMetaData.ABI

var ArbitrumModuleBin = ArbitrumModuleMetaData.Bin

func DeployArbitrumModule(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ArbitrumModule, error) {
	parsed, err := ArbitrumModuleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ArbitrumModuleBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ArbitrumModule{address: address, abi: *parsed, ArbitrumModuleCaller: ArbitrumModuleCaller{contract: contract}, ArbitrumModuleTransactor: ArbitrumModuleTransactor{contract: contract}, ArbitrumModuleFilterer: ArbitrumModuleFilterer{contract: contract}}, nil
}

type ArbitrumModule struct {
	address common.Address
	abi     abi.ABI
	ArbitrumModuleCaller
	ArbitrumModuleTransactor
	ArbitrumModuleFilterer
}

type ArbitrumModuleCaller struct {
	contract *bind.BoundContract
}

type ArbitrumModuleTransactor struct {
	contract *bind.BoundContract
}

type ArbitrumModuleFilterer struct {
	contract *bind.BoundContract
}

type ArbitrumModuleSession struct {
	Contract     *ArbitrumModule
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type ArbitrumModuleCallerSession struct {
	Contract *ArbitrumModuleCaller
	CallOpts bind.CallOpts
}

type ArbitrumModuleTransactorSession struct {
	Contract     *ArbitrumModuleTransactor
	TransactOpts bind.TransactOpts
}

type ArbitrumModuleRaw struct {
	Contract *ArbitrumModule
}

type ArbitrumModuleCallerRaw struct {
	Contract *ArbitrumModuleCaller
}

type ArbitrumModuleTransactorRaw struct {
	Contract *ArbitrumModuleTransactor
}

func NewArbitrumModule(address common.Address, backend bind.ContractBackend) (*ArbitrumModule, error) {
	abi, err := abi.JSON(strings.NewReader(ArbitrumModuleABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindArbitrumModule(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ArbitrumModule{address: address, abi: abi, ArbitrumModuleCaller: ArbitrumModuleCaller{contract: contract}, ArbitrumModuleTransactor: ArbitrumModuleTransactor{contract: contract}, ArbitrumModuleFilterer: ArbitrumModuleFilterer{contract: contract}}, nil
}

func NewArbitrumModuleCaller(address common.Address, caller bind.ContractCaller) (*ArbitrumModuleCaller, error) {
	contract, err := bindArbitrumModule(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ArbitrumModuleCaller{contract: contract}, nil
}

func NewArbitrumModuleTransactor(address common.Address, transactor bind.ContractTransactor) (*ArbitrumModuleTransactor, error) {
	contract, err := bindArbitrumModule(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ArbitrumModuleTransactor{contract: contract}, nil
}

func NewArbitrumModuleFilterer(address common.Address, filterer bind.ContractFilterer) (*ArbitrumModuleFilterer, error) {
	contract, err := bindArbitrumModule(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ArbitrumModuleFilterer{contract: contract}, nil
}

func bindArbitrumModule(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ArbitrumModuleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_ArbitrumModule *ArbitrumModuleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ArbitrumModule.Contract.ArbitrumModuleCaller.contract.Call(opts, result, method, params...)
}

func (_ArbitrumModule *ArbitrumModuleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ArbitrumModule.Contract.ArbitrumModuleTransactor.contract.Transfer(opts)
}

func (_ArbitrumModule *ArbitrumModuleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ArbitrumModule.Contract.ArbitrumModuleTransactor.contract.Transact(opts, method, params...)
}

func (_ArbitrumModule *ArbitrumModuleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ArbitrumModule.Contract.contract.Call(opts, result, method, params...)
}

func (_ArbitrumModule *ArbitrumModuleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ArbitrumModule.Contract.contract.Transfer(opts)
}

func (_ArbitrumModule *ArbitrumModuleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ArbitrumModule.Contract.contract.Transact(opts, method, params...)
}

func (_ArbitrumModule *ArbitrumModuleCaller) BlockHash(opts *bind.CallOpts, n *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _ArbitrumModule.contract.Call(opts, &out, "blockHash", n)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_ArbitrumModule *ArbitrumModuleSession) BlockHash(n *big.Int) ([32]byte, error) {
	return _ArbitrumModule.Contract.BlockHash(&_ArbitrumModule.CallOpts, n)
}

func (_ArbitrumModule *ArbitrumModuleCallerSession) BlockHash(n *big.Int) ([32]byte, error) {
	return _ArbitrumModule.Contract.BlockHash(&_ArbitrumModule.CallOpts, n)
}

func (_ArbitrumModule *ArbitrumModuleCaller) BlockNumber(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ArbitrumModule.contract.Call(opts, &out, "blockNumber")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_ArbitrumModule *ArbitrumModuleSession) BlockNumber() (*big.Int, error) {
	return _ArbitrumModule.Contract.BlockNumber(&_ArbitrumModule.CallOpts)
}

func (_ArbitrumModule *ArbitrumModuleCallerSession) BlockNumber() (*big.Int, error) {
	return _ArbitrumModule.Contract.BlockNumber(&_ArbitrumModule.CallOpts)
}

func (_ArbitrumModule *ArbitrumModuleCaller) GetCurrentL1Fee(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _ArbitrumModule.contract.Call(opts, &out, "getCurrentL1Fee", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_ArbitrumModule *ArbitrumModuleSession) GetCurrentL1Fee(arg0 *big.Int) (*big.Int, error) {
	return _ArbitrumModule.Contract.GetCurrentL1Fee(&_ArbitrumModule.CallOpts, arg0)
}

func (_ArbitrumModule *ArbitrumModuleCallerSession) GetCurrentL1Fee(arg0 *big.Int) (*big.Int, error) {
	return _ArbitrumModule.Contract.GetCurrentL1Fee(&_ArbitrumModule.CallOpts, arg0)
}

func (_ArbitrumModule *ArbitrumModuleCaller) GetGasOverhead(opts *bind.CallOpts) (GetGasOverhead,

	error) {
	var out []interface{}
	err := _ArbitrumModule.contract.Call(opts, &out, "getGasOverhead")

	outstruct := new(GetGasOverhead)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ChainModuleFixedOverhead = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.ChainModulePerByteOverhead = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_ArbitrumModule *ArbitrumModuleSession) GetGasOverhead() (GetGasOverhead,

	error) {
	return _ArbitrumModule.Contract.GetGasOverhead(&_ArbitrumModule.CallOpts)
}

func (_ArbitrumModule *ArbitrumModuleCallerSession) GetGasOverhead() (GetGasOverhead,

	error) {
	return _ArbitrumModule.Contract.GetGasOverhead(&_ArbitrumModule.CallOpts)
}

func (_ArbitrumModule *ArbitrumModuleCaller) GetMaxL1Fee(opts *bind.CallOpts, dataSize *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _ArbitrumModule.contract.Call(opts, &out, "getMaxL1Fee", dataSize)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_ArbitrumModule *ArbitrumModuleSession) GetMaxL1Fee(dataSize *big.Int) (*big.Int, error) {
	return _ArbitrumModule.Contract.GetMaxL1Fee(&_ArbitrumModule.CallOpts, dataSize)
}

func (_ArbitrumModule *ArbitrumModuleCallerSession) GetMaxL1Fee(dataSize *big.Int) (*big.Int, error) {
	return _ArbitrumModule.Contract.GetMaxL1Fee(&_ArbitrumModule.CallOpts, dataSize)
}

type GetGasOverhead struct {
	ChainModuleFixedOverhead   *big.Int
	ChainModulePerByteOverhead *big.Int
}

func (_ArbitrumModule *ArbitrumModule) Address() common.Address {
	return _ArbitrumModule.address
}

type ArbitrumModuleInterface interface {
	BlockHash(opts *bind.CallOpts, n *big.Int) ([32]byte, error)

	BlockNumber(opts *bind.CallOpts) (*big.Int, error)

	GetCurrentL1Fee(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	GetGasOverhead(opts *bind.CallOpts) (GetGasOverhead,

		error)

	GetMaxL1Fee(opts *bind.CallOpts, dataSize *big.Int) (*big.Int, error)

	Address() common.Address
}
