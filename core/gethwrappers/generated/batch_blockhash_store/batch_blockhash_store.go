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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"blockhashStoreAddr\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"BHS\",\"outputs\":[{\"internalType\":\"contractBlockhashStore\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"blockNumbers\",\"type\":\"uint256[]\"}],\"name\":\"getBlockhashes\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"blockNumbers\",\"type\":\"uint256[]\"}],\"name\":\"store\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"blockNumbers\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes[]\",\"name\":\"headers\",\"type\":\"bytes[]\"}],\"name\":\"storeVerifyHeader\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561001057600080fd5b50604051610b81380380610b8183398101604081905261002f91610044565b60601b6001600160601b031916608052610074565b60006020828403121561005657600080fd5b81516001600160a01b038116811461006d57600080fd5b9392505050565b60805160601c610adb6100a66000396000818160a7015281816101270152818161023a01526104290152610adb6000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c806306bd010d146100515780631f600f86146100665780635d290e211461008f578063f745eafb146100a2575b600080fd5b61006461005f366004610654565b6100ee565b005b610079610074366004610654565b6101e2565b60405161008691906107ff565b60405180910390f35b61006461009d366004610691565b6103ac565b6100c97f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610086565b60005b81518110156101de5761011c82828151811061010f5761010f6109ac565b60200260200101516104fe565b610125576101cc565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16636057361d838381518110610173576101736109ac565b60200260200101516040518263ffffffff1660e01b815260040161019991815260200190565b600060405180830381600087803b1580156101b357600080fd5b505af11580156101c7573d6000803e3d6000fd5b505050505b806101d681610944565b9150506100f1565b5050565b60606000825167ffffffffffffffff811115610200576102006109db565b604051908082528060200260200182016040528015610229578160200160208202803683370190505b50905060005b83518110156103a5577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663e9413d38858381518110610286576102866109ac565b60200260200101516040518263ffffffff1660e01b81526004016102ac91815260200190565b60206040518083038186803b1580156102c457600080fd5b505afa925050508015610312575060408051601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820190925261030f918101906107e6565b60015b6103725761031e610a0a565b806308c379a014156103665750610333610a26565b8061033e5750610368565b6000801b838381518110610354576103546109ac565b60200260200101818152505050610393565b505b3d6000803e3d6000fd5b80838381518110610385576103856109ac565b602002602001018181525050505b8061039d81610944565b91505061022f565b5092915050565b805182511461041b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f696e70757420617272617920617267206c656e67746873206d69736d61746368604482015260640160405180910390fd5b60005b82518110156104f9577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663fadff0e1848381518110610475576104756109ac565b602002602001015184848151811061048f5761048f6109ac565b60200260200101516040518363ffffffff1660e01b81526004016104b4929190610843565b600060405180830381600087803b1580156104ce57600080fd5b505af11580156104e2573d6000803e3d6000fd5b5050505080806104f190610944565b91505061041e565b505050565b600061010061050b610537565b111561052e5761010061051c610537565b61052691906108e2565b821015610531565b60015b92915050565b60004661a4b181148061054c575062066eed81145b156105d657606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b815260040160206040518083038186803b15801561059857600080fd5b505afa1580156105ac573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105d091906107e6565b91505090565b4391505090565b600082601f8301126105ee57600080fd5b813560206105fb826108be565b60405161060882826108f9565b8381528281019150858301600585901b8701840188101561062857600080fd5b60005b858110156106475781358452928401929084019060010161062b565b5090979650505050505050565b60006020828403121561066657600080fd5b813567ffffffffffffffff81111561067d57600080fd5b610689848285016105dd565b949350505050565b60008060408084860312156106a557600080fd5b833567ffffffffffffffff808211156106bd57600080fd5b6106c9878388016105dd565b94506020915081860135818111156106e057600080fd5b8601601f810188136106f157600080fd5b80356106fc816108be565b855161070882826108f9565b8281528581019150838601600584901b850187018c101561072857600080fd5b60005b848110156107d45781358781111561074257600080fd5b8601603f81018e1361075357600080fd5b8881013588811115610767576107676109db565b8a5161079a8b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f85011601826108f9565b8181528f8c8385010111156107ae57600080fd5b818c84018c83013760009181018b0191909152855250928701929087019060010161072b565b50989b909a5098505050505050505050565b6000602082840312156107f857600080fd5b5051919050565b6020808252825182820181905260009190848201906040850190845b818110156108375783518352928401929184019160010161081b565b50909695505050505050565b82815260006020604081840152835180604085015260005b818110156108775785810183015185820160600152820161085b565b81811115610889576000606083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01692909201606001949350505050565b600067ffffffffffffffff8211156108d8576108d86109db565b5060051b60200190565b6000828210156108f4576108f461097d565b500390565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f830116810181811067ffffffffffffffff8211171561093d5761093d6109db565b6040525050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8214156109765761097661097d565b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600060033d1115610a235760046000803e5060005160e01c5b90565b600060443d1015610a345790565b6040517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc803d016004833e81513d67ffffffffffffffff8160248401118184111715610a8257505050505090565b8285019150815181811115610a9a5750505050505090565b843d8701016020828501011115610ab45750505050505090565b610ac3602082860101876108f9565b50909594505050505056fea164736f6c6343000806000a",
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
