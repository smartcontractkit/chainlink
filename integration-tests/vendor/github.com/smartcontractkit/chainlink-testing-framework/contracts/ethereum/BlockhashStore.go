// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ethereum

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

// Reference imports to suppress errors if they are not otherwise used.
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

// BlockhashStoreMetaData contains all meta data concerning the BlockhashStore contract.
var BlockhashStoreMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getBlockhash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"store\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"storeEarliest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"header\",\"type\":\"bytes\"}],\"name\":\"storeVerifyHeader\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506102e0806100206000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80636057361d1461005157806383b6d6b714610070578063e9413d3814610078578063fadff0e1146100a7575b600080fd5b61006e6004803603602081101561006757600080fd5b5035610152565b005b61006e6101ae565b6100956004803603602081101561008e57600080fd5b50356101bd565b60408051918252519081900360200190f35b61006e600480360360408110156100bd57600080fd5b81359190810190604081016020820135600160201b8111156100de57600080fd5b8201836020820111156100f057600080fd5b803590602001918460018302840111600160201b8311171561011157600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610224945050505050565b80408061019c576040805162461bcd60e51b8152602060048201526013602482015272189b1bd8dada185cda0a1b8a4819985a5b1959606a1b604482015290519081900360640190fd5b60009182526020829052604090912055565b6101bb6101004303610152565b565b6000818152602081905260408120548061021e576040805162461bcd60e51b815260206004820152601c60248201527f626c6f636b68617368206e6f7420666f756e6420696e2073746f726500000000604482015290519081900360640190fd5b92915050565b60008083600101815260200190815260200160002054818051906020012014610294576040805162461bcd60e51b815260206004820152601c60248201527f6865616465722068617320756e6b6e6f776e20626c6f636b6861736800000000604482015290519081900360640190fd5b602401516000918252602082905260409091205556fea2646970667358221220eb9f06d3c2357d51e1006d5b541cb3c45442fd8a3d3c58beedeb0c21361f0cdd64736f6c63430006060033",
}

// BlockhashStoreABI is the input ABI used to generate the binding from.
// Deprecated: Use BlockhashStoreMetaData.ABI instead.
var BlockhashStoreABI = BlockhashStoreMetaData.ABI

// BlockhashStoreBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BlockhashStoreMetaData.Bin instead.
var BlockhashStoreBin = BlockhashStoreMetaData.Bin

// DeployBlockhashStore deploys a new Ethereum contract, binding an instance of BlockhashStore to it.
func DeployBlockhashStore(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *BlockhashStore, error) {
	parsed, err := BlockhashStoreMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BlockhashStoreBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BlockhashStore{BlockhashStoreCaller: BlockhashStoreCaller{contract: contract}, BlockhashStoreTransactor: BlockhashStoreTransactor{contract: contract}, BlockhashStoreFilterer: BlockhashStoreFilterer{contract: contract}}, nil
}

// BlockhashStore is an auto generated Go binding around an Ethereum contract.
type BlockhashStore struct {
	BlockhashStoreCaller     // Read-only binding to the contract
	BlockhashStoreTransactor // Write-only binding to the contract
	BlockhashStoreFilterer   // Log filterer for contract events
}

// BlockhashStoreCaller is an auto generated read-only Go binding around an Ethereum contract.
type BlockhashStoreCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BlockhashStoreTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BlockhashStoreTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BlockhashStoreFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BlockhashStoreFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BlockhashStoreSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BlockhashStoreSession struct {
	Contract     *BlockhashStore   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BlockhashStoreCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BlockhashStoreCallerSession struct {
	Contract *BlockhashStoreCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// BlockhashStoreTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BlockhashStoreTransactorSession struct {
	Contract     *BlockhashStoreTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// BlockhashStoreRaw is an auto generated low-level Go binding around an Ethereum contract.
type BlockhashStoreRaw struct {
	Contract *BlockhashStore // Generic contract binding to access the raw methods on
}

// BlockhashStoreCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BlockhashStoreCallerRaw struct {
	Contract *BlockhashStoreCaller // Generic read-only contract binding to access the raw methods on
}

// BlockhashStoreTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BlockhashStoreTransactorRaw struct {
	Contract *BlockhashStoreTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBlockhashStore creates a new instance of BlockhashStore, bound to a specific deployed contract.
func NewBlockhashStore(address common.Address, backend bind.ContractBackend) (*BlockhashStore, error) {
	contract, err := bindBlockhashStore(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BlockhashStore{BlockhashStoreCaller: BlockhashStoreCaller{contract: contract}, BlockhashStoreTransactor: BlockhashStoreTransactor{contract: contract}, BlockhashStoreFilterer: BlockhashStoreFilterer{contract: contract}}, nil
}

// NewBlockhashStoreCaller creates a new read-only instance of BlockhashStore, bound to a specific deployed contract.
func NewBlockhashStoreCaller(address common.Address, caller bind.ContractCaller) (*BlockhashStoreCaller, error) {
	contract, err := bindBlockhashStore(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BlockhashStoreCaller{contract: contract}, nil
}

// NewBlockhashStoreTransactor creates a new write-only instance of BlockhashStore, bound to a specific deployed contract.
func NewBlockhashStoreTransactor(address common.Address, transactor bind.ContractTransactor) (*BlockhashStoreTransactor, error) {
	contract, err := bindBlockhashStore(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BlockhashStoreTransactor{contract: contract}, nil
}

// NewBlockhashStoreFilterer creates a new log filterer instance of BlockhashStore, bound to a specific deployed contract.
func NewBlockhashStoreFilterer(address common.Address, filterer bind.ContractFilterer) (*BlockhashStoreFilterer, error) {
	contract, err := bindBlockhashStore(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BlockhashStoreFilterer{contract: contract}, nil
}

// bindBlockhashStore binds a generic wrapper to an already deployed contract.
func bindBlockhashStore(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BlockhashStoreABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BlockhashStore *BlockhashStoreRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BlockhashStore.Contract.BlockhashStoreCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BlockhashStore *BlockhashStoreRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BlockhashStore.Contract.BlockhashStoreTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BlockhashStore *BlockhashStoreRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BlockhashStore.Contract.BlockhashStoreTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BlockhashStore *BlockhashStoreCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BlockhashStore.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BlockhashStore *BlockhashStoreTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BlockhashStore.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BlockhashStore *BlockhashStoreTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BlockhashStore.Contract.contract.Transact(opts, method, params...)
}

// GetBlockhash is a free data retrieval call binding the contract method 0xe9413d38.
//
// Solidity: function getBlockhash(uint256 n) view returns(bytes32)
func (_BlockhashStore *BlockhashStoreCaller) GetBlockhash(opts *bind.CallOpts, n *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _BlockhashStore.contract.Call(opts, &out, "getBlockhash", n)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetBlockhash is a free data retrieval call binding the contract method 0xe9413d38.
//
// Solidity: function getBlockhash(uint256 n) view returns(bytes32)
func (_BlockhashStore *BlockhashStoreSession) GetBlockhash(n *big.Int) ([32]byte, error) {
	return _BlockhashStore.Contract.GetBlockhash(&_BlockhashStore.CallOpts, n)
}

// GetBlockhash is a free data retrieval call binding the contract method 0xe9413d38.
//
// Solidity: function getBlockhash(uint256 n) view returns(bytes32)
func (_BlockhashStore *BlockhashStoreCallerSession) GetBlockhash(n *big.Int) ([32]byte, error) {
	return _BlockhashStore.Contract.GetBlockhash(&_BlockhashStore.CallOpts, n)
}

// Store is a paid mutator transaction binding the contract method 0x6057361d.
//
// Solidity: function store(uint256 n) returns()
func (_BlockhashStore *BlockhashStoreTransactor) Store(opts *bind.TransactOpts, n *big.Int) (*types.Transaction, error) {
	return _BlockhashStore.contract.Transact(opts, "store", n)
}

// Store is a paid mutator transaction binding the contract method 0x6057361d.
//
// Solidity: function store(uint256 n) returns()
func (_BlockhashStore *BlockhashStoreSession) Store(n *big.Int) (*types.Transaction, error) {
	return _BlockhashStore.Contract.Store(&_BlockhashStore.TransactOpts, n)
}

// Store is a paid mutator transaction binding the contract method 0x6057361d.
//
// Solidity: function store(uint256 n) returns()
func (_BlockhashStore *BlockhashStoreTransactorSession) Store(n *big.Int) (*types.Transaction, error) {
	return _BlockhashStore.Contract.Store(&_BlockhashStore.TransactOpts, n)
}

// StoreEarliest is a paid mutator transaction binding the contract method 0x83b6d6b7.
//
// Solidity: function storeEarliest() returns()
func (_BlockhashStore *BlockhashStoreTransactor) StoreEarliest(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BlockhashStore.contract.Transact(opts, "storeEarliest")
}

// StoreEarliest is a paid mutator transaction binding the contract method 0x83b6d6b7.
//
// Solidity: function storeEarliest() returns()
func (_BlockhashStore *BlockhashStoreSession) StoreEarliest() (*types.Transaction, error) {
	return _BlockhashStore.Contract.StoreEarliest(&_BlockhashStore.TransactOpts)
}

// StoreEarliest is a paid mutator transaction binding the contract method 0x83b6d6b7.
//
// Solidity: function storeEarliest() returns()
func (_BlockhashStore *BlockhashStoreTransactorSession) StoreEarliest() (*types.Transaction, error) {
	return _BlockhashStore.Contract.StoreEarliest(&_BlockhashStore.TransactOpts)
}

// StoreVerifyHeader is a paid mutator transaction binding the contract method 0xfadff0e1.
//
// Solidity: function storeVerifyHeader(uint256 n, bytes header) returns()
func (_BlockhashStore *BlockhashStoreTransactor) StoreVerifyHeader(opts *bind.TransactOpts, n *big.Int, header []byte) (*types.Transaction, error) {
	return _BlockhashStore.contract.Transact(opts, "storeVerifyHeader", n, header)
}

// StoreVerifyHeader is a paid mutator transaction binding the contract method 0xfadff0e1.
//
// Solidity: function storeVerifyHeader(uint256 n, bytes header) returns()
func (_BlockhashStore *BlockhashStoreSession) StoreVerifyHeader(n *big.Int, header []byte) (*types.Transaction, error) {
	return _BlockhashStore.Contract.StoreVerifyHeader(&_BlockhashStore.TransactOpts, n, header)
}

// StoreVerifyHeader is a paid mutator transaction binding the contract method 0xfadff0e1.
//
// Solidity: function storeVerifyHeader(uint256 n, bytes header) returns()
func (_BlockhashStore *BlockhashStoreTransactorSession) StoreVerifyHeader(n *big.Int, header []byte) (*types.Transaction, error) {
	return _BlockhashStore.Contract.StoreVerifyHeader(&_BlockhashStore.TransactOpts, n, header)
}
