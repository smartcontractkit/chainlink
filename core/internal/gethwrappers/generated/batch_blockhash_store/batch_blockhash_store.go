// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package batch_blockhash_store

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

var BatchBlockhashStoreMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"blockhashStoreAddr\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"BHS\",\"outputs\":[{\"internalType\":\"contractBlockhashStore\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"blockNumbers\",\"type\":\"uint256[]\"}],\"name\":\"getBlockhashes\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"blockNumbers\",\"type\":\"uint256[]\"}],\"name\":\"store\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"blockNumbers\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes[]\",\"name\":\"headers\",\"type\":\"bytes[]\"}],\"name\":\"storeVerifyHeader\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561001057600080fd5b5060405161090c38038061090c83398101604081905261002f91610044565b60601b6001600160601b031916608052610074565b60006020828403121561005657600080fd5b81516001600160a01b038116811461006d57600080fd5b9392505050565b60805160601c6108676100a56000396000818160a70152818160fc0152818161020f015261038401526108676000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c806306bd010d146100515780631f600f86146100665780635d290e211461008f578063f745eafb146100a2575b600080fd5b61006461005f3660046104cb565b6100ee565b005b6100796100743660046104cb565b6101b7565b604051610086919061066a565b60405180910390f35b61006461009d366004610508565b610307565b6100c97f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610086565b60005b81518110156101b3577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16636057361d838381518110610148576101486107fc565b60200260200101516040518263ffffffff1660e01b815260040161016e91815260200190565b600060405180830381600087803b15801561018857600080fd5b505af115801561019c573d6000803e3d6000fd5b5050505080806101ab9061079c565b9150506100f1565b5050565b60606000825167ffffffffffffffff8111156101d5576101d561082b565b6040519080825280602002602001820160405280156101fe578160200160208202803683370190505b50905060005b8351811015610300577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663e9413d3885838151811061025b5761025b6107fc565b60200260200101516040518263ffffffff1660e01b815260040161028191815260200190565b60206040518083038186803b15801561029957600080fd5b505afa1580156102ad573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102d19190610651565b8282815181106102e3576102e36107fc565b6020908102919091010152806102f88161079c565b915050610204565b5092915050565b8051825114610376576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f696e70757420617272617920617267206c656e67746873206d69736d61746368604482015260640160405180910390fd5b60005b8251811015610454577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663fadff0e18483815181106103d0576103d06107fc565b60200260200101518484815181106103ea576103ea6107fc565b60200260200101516040518363ffffffff1660e01b815260040161040f9291906106ae565b600060405180830381600087803b15801561042957600080fd5b505af115801561043d573d6000803e3d6000fd5b50505050808061044c9061079c565b915050610379565b505050565b600082601f83011261046a57600080fd5b8135602061047f61047a83610778565b610729565b80838252828201915082860187848660051b890101111561049f57600080fd5b60005b858110156104be578135845292840192908401906001016104a2565b5090979650505050505050565b6000602082840312156104dd57600080fd5b813567ffffffffffffffff8111156104f457600080fd5b61050084828501610459565b949350505050565b600080604080848603121561051c57600080fd5b833567ffffffffffffffff8082111561053457600080fd5b61054087838801610459565b945060209150818601358181111561055757600080fd5b8601601f8101881361056857600080fd5b803561057661047a82610778565b8082825285820191508584018b878560051b870101111561059657600080fd5b60005b8481101561063f578135878111156105b057600080fd5b8601603f81018e136105c157600080fd5b88810135888111156105d5576105d561082b565b6106058a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601610729565b8181528f8c83850101111561061957600080fd5b818c84018c83013760009181018b01919091528552509287019290870190600101610599565b50989b909a5098505050505050505050565b60006020828403121561066357600080fd5b5051919050565b6020808252825182820181905260009190848201906040850190845b818110156106a257835183529284019291840191600101610686565b50909695505050505050565b82815260006020604081840152835180604085015260005b818110156106e2578581018301518582016060015282016106c6565b818111156106f4576000606083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01692909201606001949350505050565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156107705761077061082b565b604052919050565b600067ffffffffffffffff8211156107925761079261082b565b5060051b60200190565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8214156107f5577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var BatchBlockhashStoreABI = BatchBlockhashStoreMetaData.ABI

var BatchBlockhashStoreBin = BatchBlockhashStoreMetaData.Bin

func DeployBatchBlockhashStore(auth *bind.TransactOpts, backend bind.ContractBackend, blockhashStoreAddr common.Address) (common.Address, *types.Transaction, *BatchBlockhashStore, error) {
	parsed, err := BatchBlockhashStoreMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BatchBlockhashStoreBin), backend, blockhashStoreAddr)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BatchBlockhashStore{BatchBlockhashStoreCaller: BatchBlockhashStoreCaller{contract: contract}, BatchBlockhashStoreTransactor: BatchBlockhashStoreTransactor{contract: contract}, BatchBlockhashStoreFilterer: BatchBlockhashStoreFilterer{contract: contract}}, nil
}

type BatchBlockhashStore struct {
	address common.Address
	abi     abi.ABI
	BatchBlockhashStoreCaller
	BatchBlockhashStoreTransactor
	BatchBlockhashStoreFilterer
}

type BatchBlockhashStoreCaller struct {
	contract *bind.BoundContract
}

type BatchBlockhashStoreTransactor struct {
	contract *bind.BoundContract
}

type BatchBlockhashStoreFilterer struct {
	contract *bind.BoundContract
}

type BatchBlockhashStoreSession struct {
	Contract     *BatchBlockhashStore
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type BatchBlockhashStoreCallerSession struct {
	Contract *BatchBlockhashStoreCaller
	CallOpts bind.CallOpts
}

type BatchBlockhashStoreTransactorSession struct {
	Contract     *BatchBlockhashStoreTransactor
	TransactOpts bind.TransactOpts
}

type BatchBlockhashStoreRaw struct {
	Contract *BatchBlockhashStore
}

type BatchBlockhashStoreCallerRaw struct {
	Contract *BatchBlockhashStoreCaller
}

type BatchBlockhashStoreTransactorRaw struct {
	Contract *BatchBlockhashStoreTransactor
}

func NewBatchBlockhashStore(address common.Address, backend bind.ContractBackend) (*BatchBlockhashStore, error) {
	abi, err := abi.JSON(strings.NewReader(BatchBlockhashStoreABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindBatchBlockhashStore(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BatchBlockhashStore{address: address, abi: abi, BatchBlockhashStoreCaller: BatchBlockhashStoreCaller{contract: contract}, BatchBlockhashStoreTransactor: BatchBlockhashStoreTransactor{contract: contract}, BatchBlockhashStoreFilterer: BatchBlockhashStoreFilterer{contract: contract}}, nil
}

func NewBatchBlockhashStoreCaller(address common.Address, caller bind.ContractCaller) (*BatchBlockhashStoreCaller, error) {
	contract, err := bindBatchBlockhashStore(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BatchBlockhashStoreCaller{contract: contract}, nil
}

func NewBatchBlockhashStoreTransactor(address common.Address, transactor bind.ContractTransactor) (*BatchBlockhashStoreTransactor, error) {
	contract, err := bindBatchBlockhashStore(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BatchBlockhashStoreTransactor{contract: contract}, nil
}

func NewBatchBlockhashStoreFilterer(address common.Address, filterer bind.ContractFilterer) (*BatchBlockhashStoreFilterer, error) {
	contract, err := bindBatchBlockhashStore(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BatchBlockhashStoreFilterer{contract: contract}, nil
}

func bindBatchBlockhashStore(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BatchBlockhashStoreABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_BatchBlockhashStore *BatchBlockhashStoreRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BatchBlockhashStore.Contract.BatchBlockhashStoreCaller.contract.Call(opts, result, method, params...)
}

func (_BatchBlockhashStore *BatchBlockhashStoreRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BatchBlockhashStore.Contract.BatchBlockhashStoreTransactor.contract.Transfer(opts)
}

func (_BatchBlockhashStore *BatchBlockhashStoreRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BatchBlockhashStore.Contract.BatchBlockhashStoreTransactor.contract.Transact(opts, method, params...)
}

func (_BatchBlockhashStore *BatchBlockhashStoreCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BatchBlockhashStore.Contract.contract.Call(opts, result, method, params...)
}

func (_BatchBlockhashStore *BatchBlockhashStoreTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BatchBlockhashStore.Contract.contract.Transfer(opts)
}

func (_BatchBlockhashStore *BatchBlockhashStoreTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BatchBlockhashStore.Contract.contract.Transact(opts, method, params...)
}

func (_BatchBlockhashStore *BatchBlockhashStoreCaller) BHS(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BatchBlockhashStore.contract.Call(opts, &out, "BHS")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BatchBlockhashStore *BatchBlockhashStoreSession) BHS() (common.Address, error) {
	return _BatchBlockhashStore.Contract.BHS(&_BatchBlockhashStore.CallOpts)
}

func (_BatchBlockhashStore *BatchBlockhashStoreCallerSession) BHS() (common.Address, error) {
	return _BatchBlockhashStore.Contract.BHS(&_BatchBlockhashStore.CallOpts)
}

func (_BatchBlockhashStore *BatchBlockhashStoreCaller) GetBlockhashes(opts *bind.CallOpts, blockNumbers []*big.Int) ([][32]byte, error) {
	var out []interface{}
	err := _BatchBlockhashStore.contract.Call(opts, &out, "getBlockhashes", blockNumbers)

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

func (_BatchBlockhashStore *BatchBlockhashStoreSession) GetBlockhashes(blockNumbers []*big.Int) ([][32]byte, error) {
	return _BatchBlockhashStore.Contract.GetBlockhashes(&_BatchBlockhashStore.CallOpts, blockNumbers)
}

func (_BatchBlockhashStore *BatchBlockhashStoreCallerSession) GetBlockhashes(blockNumbers []*big.Int) ([][32]byte, error) {
	return _BatchBlockhashStore.Contract.GetBlockhashes(&_BatchBlockhashStore.CallOpts, blockNumbers)
}

func (_BatchBlockhashStore *BatchBlockhashStoreTransactor) Store(opts *bind.TransactOpts, blockNumbers []*big.Int) (*types.Transaction, error) {
	return _BatchBlockhashStore.contract.Transact(opts, "store", blockNumbers)
}

func (_BatchBlockhashStore *BatchBlockhashStoreSession) Store(blockNumbers []*big.Int) (*types.Transaction, error) {
	return _BatchBlockhashStore.Contract.Store(&_BatchBlockhashStore.TransactOpts, blockNumbers)
}

func (_BatchBlockhashStore *BatchBlockhashStoreTransactorSession) Store(blockNumbers []*big.Int) (*types.Transaction, error) {
	return _BatchBlockhashStore.Contract.Store(&_BatchBlockhashStore.TransactOpts, blockNumbers)
}

func (_BatchBlockhashStore *BatchBlockhashStoreTransactor) StoreVerifyHeader(opts *bind.TransactOpts, blockNumbers []*big.Int, headers [][]byte) (*types.Transaction, error) {
	return _BatchBlockhashStore.contract.Transact(opts, "storeVerifyHeader", blockNumbers, headers)
}

func (_BatchBlockhashStore *BatchBlockhashStoreSession) StoreVerifyHeader(blockNumbers []*big.Int, headers [][]byte) (*types.Transaction, error) {
	return _BatchBlockhashStore.Contract.StoreVerifyHeader(&_BatchBlockhashStore.TransactOpts, blockNumbers, headers)
}

func (_BatchBlockhashStore *BatchBlockhashStoreTransactorSession) StoreVerifyHeader(blockNumbers []*big.Int, headers [][]byte) (*types.Transaction, error) {
	return _BatchBlockhashStore.Contract.StoreVerifyHeader(&_BatchBlockhashStore.TransactOpts, blockNumbers, headers)
}

func (_BatchBlockhashStore *BatchBlockhashStore) Address() common.Address {
	return _BatchBlockhashStore.address
}

type BatchBlockhashStoreInterface interface {
	BHS(opts *bind.CallOpts) (common.Address, error)

	GetBlockhashes(opts *bind.CallOpts, blockNumbers []*big.Int) ([][32]byte, error)

	Store(opts *bind.TransactOpts, blockNumbers []*big.Int) (*types.Transaction, error)

	StoreVerifyHeader(opts *bind.TransactOpts, blockNumbers []*big.Int, headers [][]byte) (*types.Transaction, error)

	Address() common.Address
}
