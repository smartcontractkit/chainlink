// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package blockhash_store

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

var BlockhashStoreMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getBlockhash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"store\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"storeEarliest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"header\",\"type\":\"bytes\"}],\"name\":\"storeVerifyHeader\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610311806100206000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80636057361d1461005157806383b6d6b714610070578063e9413d3814610078578063fadff0e1146100a7575b600080fd5b61006e6004803603602081101561006757600080fd5b5035610154565b005b61006e6101d4565b6100956004803603602081101561008e57600080fd5b50356101e3565b60408051918252519081900360200190f35b61006e600480360360408110156100bd57600080fd5b813591908101906040810160208201356401000000008111156100df57600080fd5b8201836020820111156100f157600080fd5b8035906020019184600183028401116401000000008311171561011357600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610264945050505050565b8040806101c257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f626c6f636b68617368286e29206661696c656400000000000000000000000000604482015290519081900360640190fd5b60009182526020829052604090912055565b6101e16101004303610154565b565b6000818152602081905260408120548061025e57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f626c6f636b68617368206e6f7420666f756e6420696e2073746f726500000000604482015290519081900360640190fd5b92915050565b600080836001018152602001908152602001600020548180519060200120146102ee57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f6865616465722068617320756e6b6e6f776e20626c6f636b6861736800000000604482015290519081900360640190fd5b602401516000918252602082905260409091205556fea164736f6c6343000606000a",
}

var BlockhashStoreABI = BlockhashStoreMetaData.ABI

var BlockhashStoreBin = BlockhashStoreMetaData.Bin

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

type BlockhashStore struct {
	address common.Address
	abi     abi.ABI
	BlockhashStoreCaller
	BlockhashStoreTransactor
	BlockhashStoreFilterer
}

type BlockhashStoreCaller struct {
	contract *bind.BoundContract
}

type BlockhashStoreTransactor struct {
	contract *bind.BoundContract
}

type BlockhashStoreFilterer struct {
	contract *bind.BoundContract
}

type BlockhashStoreSession struct {
	Contract     *BlockhashStore
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type BlockhashStoreCallerSession struct {
	Contract *BlockhashStoreCaller
	CallOpts bind.CallOpts
}

type BlockhashStoreTransactorSession struct {
	Contract     *BlockhashStoreTransactor
	TransactOpts bind.TransactOpts
}

type BlockhashStoreRaw struct {
	Contract *BlockhashStore
}

type BlockhashStoreCallerRaw struct {
	Contract *BlockhashStoreCaller
}

type BlockhashStoreTransactorRaw struct {
	Contract *BlockhashStoreTransactor
}

func NewBlockhashStore(address common.Address, backend bind.ContractBackend) (*BlockhashStore, error) {
	abi, err := abi.JSON(strings.NewReader(BlockhashStoreABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindBlockhashStore(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BlockhashStore{address: address, abi: abi, BlockhashStoreCaller: BlockhashStoreCaller{contract: contract}, BlockhashStoreTransactor: BlockhashStoreTransactor{contract: contract}, BlockhashStoreFilterer: BlockhashStoreFilterer{contract: contract}}, nil
}

func NewBlockhashStoreCaller(address common.Address, caller bind.ContractCaller) (*BlockhashStoreCaller, error) {
	contract, err := bindBlockhashStore(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BlockhashStoreCaller{contract: contract}, nil
}

func NewBlockhashStoreTransactor(address common.Address, transactor bind.ContractTransactor) (*BlockhashStoreTransactor, error) {
	contract, err := bindBlockhashStore(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BlockhashStoreTransactor{contract: contract}, nil
}

func NewBlockhashStoreFilterer(address common.Address, filterer bind.ContractFilterer) (*BlockhashStoreFilterer, error) {
	contract, err := bindBlockhashStore(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BlockhashStoreFilterer{contract: contract}, nil
}

func bindBlockhashStore(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BlockhashStoreABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_BlockhashStore *BlockhashStoreRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BlockhashStore.Contract.BlockhashStoreCaller.contract.Call(opts, result, method, params...)
}

func (_BlockhashStore *BlockhashStoreRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BlockhashStore.Contract.BlockhashStoreTransactor.contract.Transfer(opts)
}

func (_BlockhashStore *BlockhashStoreRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BlockhashStore.Contract.BlockhashStoreTransactor.contract.Transact(opts, method, params...)
}

func (_BlockhashStore *BlockhashStoreCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BlockhashStore.Contract.contract.Call(opts, result, method, params...)
}

func (_BlockhashStore *BlockhashStoreTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BlockhashStore.Contract.contract.Transfer(opts)
}

func (_BlockhashStore *BlockhashStoreTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BlockhashStore.Contract.contract.Transact(opts, method, params...)
}

func (_BlockhashStore *BlockhashStoreCaller) GetBlockhash(opts *bind.CallOpts, n *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _BlockhashStore.contract.Call(opts, &out, "getBlockhash", n)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_BlockhashStore *BlockhashStoreSession) GetBlockhash(n *big.Int) ([32]byte, error) {
	return _BlockhashStore.Contract.GetBlockhash(&_BlockhashStore.CallOpts, n)
}

func (_BlockhashStore *BlockhashStoreCallerSession) GetBlockhash(n *big.Int) ([32]byte, error) {
	return _BlockhashStore.Contract.GetBlockhash(&_BlockhashStore.CallOpts, n)
}

func (_BlockhashStore *BlockhashStoreTransactor) Store(opts *bind.TransactOpts, n *big.Int) (*types.Transaction, error) {
	return _BlockhashStore.contract.Transact(opts, "store", n)
}

func (_BlockhashStore *BlockhashStoreSession) Store(n *big.Int) (*types.Transaction, error) {
	return _BlockhashStore.Contract.Store(&_BlockhashStore.TransactOpts, n)
}

func (_BlockhashStore *BlockhashStoreTransactorSession) Store(n *big.Int) (*types.Transaction, error) {
	return _BlockhashStore.Contract.Store(&_BlockhashStore.TransactOpts, n)
}

func (_BlockhashStore *BlockhashStoreTransactor) StoreEarliest(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BlockhashStore.contract.Transact(opts, "storeEarliest")
}

func (_BlockhashStore *BlockhashStoreSession) StoreEarliest() (*types.Transaction, error) {
	return _BlockhashStore.Contract.StoreEarliest(&_BlockhashStore.TransactOpts)
}

func (_BlockhashStore *BlockhashStoreTransactorSession) StoreEarliest() (*types.Transaction, error) {
	return _BlockhashStore.Contract.StoreEarliest(&_BlockhashStore.TransactOpts)
}

func (_BlockhashStore *BlockhashStoreTransactor) StoreVerifyHeader(opts *bind.TransactOpts, n *big.Int, header []byte) (*types.Transaction, error) {
	return _BlockhashStore.contract.Transact(opts, "storeVerifyHeader", n, header)
}

func (_BlockhashStore *BlockhashStoreSession) StoreVerifyHeader(n *big.Int, header []byte) (*types.Transaction, error) {
	return _BlockhashStore.Contract.StoreVerifyHeader(&_BlockhashStore.TransactOpts, n, header)
}

func (_BlockhashStore *BlockhashStoreTransactorSession) StoreVerifyHeader(n *big.Int, header []byte) (*types.Transaction, error) {
	return _BlockhashStore.Contract.StoreVerifyHeader(&_BlockhashStore.TransactOpts, n, header)
}

func (_BlockhashStore *BlockhashStore) Address() common.Address {
	return _BlockhashStore.address
}

type BlockhashStoreInterface interface {
	GetBlockhash(opts *bind.CallOpts, n *big.Int) ([32]byte, error)

	Store(opts *bind.TransactOpts, n *big.Int) (*types.Transaction, error)

	StoreEarliest(opts *bind.TransactOpts) (*types.Transaction, error)

	StoreVerifyHeader(opts *bind.TransactOpts, n *big.Int, header []byte) (*types.Transaction, error)

	Address() common.Address
}
