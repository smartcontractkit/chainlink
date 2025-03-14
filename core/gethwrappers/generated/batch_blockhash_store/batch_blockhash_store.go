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
	_ = abi.ConvertType
)

var BatchBlockhashStoreMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"blockhashStoreAddr\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"BHS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractBlockhashStore\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBlockhashes\",\"inputs\":[{\"name\":\"blockNumbers\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"store\",\"inputs\":[{\"name\":\"blockNumbers\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"storeVerifyHeader\",\"inputs\":[{\"name\":\"blockNumbers\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"headers\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"}]",
	Bin: "0x60a060405234801561001057600080fd5b50604051610b5f380380610b5f83398101604081905261002f91610040565b6001600160a01b0316608052610070565b60006020828403121561005257600080fd5b81516001600160a01b038116811461006957600080fd5b9392505050565b608051610ac061009f6000396000818160a7015281816101230152818161023601526104150152610ac06000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c806306bd010d146100515780631f600f86146100665780635d290e211461008f578063f745eafb146100a2575b600080fd5b61006461005f3660046106e3565b6100ee565b005b6100796100743660046106e3565b6101de565b6040516100869190610720565b60405180910390f35b61006461009d366004610764565b610398565b6100c97f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610086565b60005b81518110156101da5761011c82828151811061010f5761010f6108b9565b60200260200101516104ea565b156101c8577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16636057361d83838151811061016f5761016f6108b9565b60200260200101516040518263ffffffff1660e01b815260040161019591815260200190565b600060405180830381600087803b1580156101af57600080fd5b505af11580156101c3573d6000803e3d6000fd5b505050505b806101d281610917565b9150506100f1565b5050565b60606000825167ffffffffffffffff8111156101fc576101fc6105d4565b604051908082528060200260200182016040528015610225578160200160208202803683370190505b50905060005b8351811015610391577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663e9413d38858381518110610282576102826108b9565b60200260200101516040518263ffffffff1660e01b81526004016102a891815260200190565b602060405180830381865afa9250505080156102ff575060408051601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682019092526102fc9181019061094f565b60015b61035e5761030b610968565b806308c379a003610352575061031f610984565b8061032a5750610354565b6000801b838381518110610340576103406108b9565b6020026020010181815250505061037f565b505b3d6000803e3d6000fd5b80838381518110610371576103716108b9565b602002602001018181525050505b8061038981610917565b91505061022b565b5092915050565b8051825114610407576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f696e70757420617272617920617267206c656e67746873206d69736d61746368604482015260640160405180910390fd5b60005b82518110156104e5577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663fadff0e1848381518110610461576104616108b9565b602002602001015184848151811061047b5761047b6108b9565b60200260200101516040518363ffffffff1660e01b81526004016104a0929190610a2c565b600060405180830381600087803b1580156104ba57600080fd5b505af11580156104ce573d6000803e3d6000fd5b5050505080806104dd90610917565b91505061040a565b505050565b60006101006104f7610523565b111561051a57610100610508610523565b6105129190610aa0565b82101561051d565b60015b92915050565b60004661052f816105b1565b156105aa57606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610580573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105a4919061094f565b91505090565b4391505090565b600061a4b18214806105c5575062066eed82145b8061051d57505062066eee1490565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f830116810181811067ffffffffffffffff82111715610647576106476105d4565b6040525050565b600067ffffffffffffffff821115610668576106686105d4565b5060051b60200190565b600082601f83011261068357600080fd5b813560206106908261064e565b60405161069d8282610603565b83815260059390931b85018201928281019150868411156106bd57600080fd5b8286015b848110156106d857803583529183019183016106c1565b509695505050505050565b6000602082840312156106f557600080fd5b813567ffffffffffffffff81111561070c57600080fd5b61071884828501610672565b949350505050565b6020808252825182820181905260009190848201906040850190845b818110156107585783518352928401929184019160010161073c565b50909695505050505050565b600080604080848603121561077857600080fd5b833567ffffffffffffffff8082111561079057600080fd5b61079c87838801610672565b94506020915081860135818111156107b357600080fd5b8601601f80820189136107c557600080fd5b81356107d08161064e565b86516107dc8282610603565b82815260059290921b840186019186810191508b8311156107fc57600080fd5b8685015b838110156108a6578035878111156108185760008081fd5b8601603f81018e1361082a5760008081fd5b888101358881111561083e5761083e6105d4565b8a516108708b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08a8501160182610603565b8181528f8c8385010111156108855760008081fd5b818c84018c83013760009181018b0191909152845250918701918701610800565b5080985050505050505050509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203610948576109486108e8565b5060010190565b60006020828403121561096157600080fd5b5051919050565b600060033d11156109815760046000803e5060005160e01c5b90565b600060443d10156109925790565b6040517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc803d016004833e81513d67ffffffffffffffff81602484011181841117156109e057505050505090565b82850191508151818111156109f85750505050505090565b843d8701016020828501011115610a125750505050505090565b610a2160208286010187610603565b509095945050505050565b82815260006020604081840152835180604085015260005b81811015610a6057858101830151858201606001528201610a44565b5060006060828601015260607fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f830116850101925050509392505050565b8181038181111561051d5761051d6108e856fea164736f6c6343000813000a",
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
	return address, tx, &BatchBlockhashStore{address: address, abi: *parsed, BatchBlockhashStoreCaller: BatchBlockhashStoreCaller{contract: contract}, BatchBlockhashStoreTransactor: BatchBlockhashStoreTransactor{contract: contract}, BatchBlockhashStoreFilterer: BatchBlockhashStoreFilterer{contract: contract}}, nil
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
	parsed, err := BatchBlockhashStoreMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
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
