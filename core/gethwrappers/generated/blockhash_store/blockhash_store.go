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
	_ = abi.ConvertType
)

var BlockhashStoreMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getBlockhash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"store\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"storeEarliest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"header\",\"type\":\"bytes\"}],\"name\":\"storeVerifyHeader\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506105d3806100206000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80636057361d1461005157806383b6d6b714610066578063e9413d381461006e578063fadff0e114610093575b600080fd5b61006461005f366004610447565b6100a6565b005b610064610131565b61008161007c366004610447565b61014b565b60405190815260200160405180910390f35b6100646100a1366004610460565b6101c7565b60006100b182610269565b90508061011f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f626c6f636b68617368286e29206661696c65640000000000000000000000000060448201526064015b60405180910390fd5b60009182526020829052604090912055565b61014961010061013f61036e565b61005f9190610551565b565b600081815260208190526040812054806101c1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f626c6f636b68617368206e6f7420666f756e6420696e2073746f7265000000006044820152606401610116565b92915050565b6000806101d5846001610539565b815260200190815260200160002054818051906020012014610253576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f6865616465722068617320756e6b6e6f776e20626c6f636b68617368000000006044820152606401610116565b6024015160009182526020829052604090912055565b6000466102758161040b565b1561035e576101008367ffffffffffffffff1661029061036e565b61029a9190610551565b11806102b757506102a961036e565b8367ffffffffffffffff1610155b156102c55750600092915050565b6040517f2b407a8200000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84166004820152606490632b407a829060240160206040518083038186803b15801561031f57600080fd5b505afa158015610333573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610357919061042e565b9392505050565b505067ffffffffffffffff164090565b60004661037a8161040b565b1561040457606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b815260040160206040518083038186803b1580156103c657600080fd5b505afa1580156103da573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103fe919061042e565b91505090565b4391505090565b600061a4b182148061041f575062066eed82145b806101c157505062066eee1490565b60006020828403121561044057600080fd5b5051919050565b60006020828403121561045957600080fd5b5035919050565b6000806040838503121561047357600080fd5b82359150602083013567ffffffffffffffff8082111561049257600080fd5b818501915085601f8301126104a657600080fd5b8135818111156104b8576104b8610597565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f011681019083821181831017156104fe576104fe610597565b8160405282815288602084870101111561051757600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b6000821982111561054c5761054c610568565b500190565b60008282101561056357610563610568565b500390565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
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
	parsed, err := BlockhashStoreMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
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
