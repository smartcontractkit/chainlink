// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package chain_specific_util_helper

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

var ChainSpecificUtilHelperMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"getBlockNumber\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"blockNumber\",\"type\":\"uint64\"}],\"name\":\"getBlockhash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"txCallData\",\"type\":\"string\"}],\"name\":\"getCurrentTxL1GasFees\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"calldataSize\",\"type\":\"uint256\"}],\"name\":\"getL1CalldataGasCost\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506109e5806100206000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c806342cbb15c1461005157806397e329a91461006b578063b778b1121461007e578063da9027ef14610091575b600080fd5b6100596100a4565b60405190815260200160405180910390f35b6100596100793660046107b1565b6100b3565b61005961008c36600461074e565b6100c4565b61005961009f36600461067f565b6100cf565b60006100ae6100da565b905090565b60006100be82610180565b92915050565b60006100be8261029c565b60006100be82610374565b60004661a4b18114806100ef575062066eed81145b1561017957606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b815260040160206040518083038186803b15801561013b57600080fd5b505afa15801561014f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101739190610666565b91505090565b4391505090565b60004661a4b1811480610195575062066eed81145b806101a2575062066eee81145b1561028c576101008367ffffffffffffffff166101bd6100da565b6101c791906108eb565b11806101e457506101d66100da565b8367ffffffffffffffff1610155b156101f25750600092915050565b6040517f2b407a8200000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84166004820152606490632b407a82906024015b60206040518083038186803b15801561024d57600080fd5b505afa158015610261573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102859190610666565b9392505050565b505067ffffffffffffffff164090565b6000466102a881610460565b15610354576000606c73ffffffffffffffffffffffffffffffffffffffff166341b247a86040518163ffffffff1660e01b815260040160c06040518083038186803b1580156102f657600080fd5b505afa15801561030a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061032e9190610767565b5050505091505083608c610342919061085b565b61034c90826108ae565b949350505050565b61035d81610483565b1561036b57610285836104bd565b50600092915050565b60004661038081610460565b156103cc57606c73ffffffffffffffffffffffffffffffffffffffff1663c6f7de0e6040518163ffffffff1660e01b815260040160206040518083038186803b15801561024d57600080fd5b6103d581610483565b1561036b5773420000000000000000000000000000000000000f73ffffffffffffffffffffffffffffffffffffffff166349948e0e84604051806080016040528060488152602001610991604891396040516020016104359291906107db565b6040516020818303038152906040526040518263ffffffff1660e01b8152600401610235919061080a565b600061a4b1821480610474575062066eed82145b806100be57505062066eee1490565b6000600a82148061049557506101a482145b806104a2575062aa37dc82145b806104ae575061210582145b806100be57505062014a331490565b60008073420000000000000000000000000000000000000f73ffffffffffffffffffffffffffffffffffffffff1663519b4bd36040518163ffffffff1660e01b815260040160206040518083038186803b15801561051a57600080fd5b505afa15801561052e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105529190610666565b90506000610561600585610873565b9050600061056f82866108eb565b9050600061057e8260106108ae565b6105898460046108ae565b610593919061085b565b9050600073420000000000000000000000000000000000000f73ffffffffffffffffffffffffffffffffffffffff16630c18c1626040518163ffffffff1660e01b815260040160206040518083038186803b1580156105f157600080fd5b505afa158015610605573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106299190610666565b905060006102ac61063a838561085b565b61064490886108ae565b610650906103e86108ae565b61065a9190610873565b98975050505050505050565b60006020828403121561067857600080fd5b5051919050565b60006020828403121561069157600080fd5b813567ffffffffffffffff808211156106a957600080fd5b818401915084601f8301126106bd57600080fd5b8135818111156106cf576106cf610961565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f0116810190838211818310171561071557610715610961565b8160405282815287602084870101111561072e57600080fd5b826020860160208301376000928101602001929092525095945050505050565b60006020828403121561076057600080fd5b5035919050565b60008060008060008060c0878903121561078057600080fd5b865195506020870151945060408701519350606087015192506080870151915060a087015190509295509295509295565b6000602082840312156107c357600080fd5b813567ffffffffffffffff8116811461028557600080fd5b600083516107ed818460208801610902565b835190830190610801818360208801610902565b01949350505050565b6020815260008251806020840152610829816040850160208701610902565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169190910160400192915050565b6000821982111561086e5761086e610932565b500190565b6000826108a9577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156108e6576108e6610932565b500290565b6000828210156108fd576108fd610932565b500390565b60005b8381101561091d578181015183820152602001610905565b8381111561092c576000848401525b50505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfe307866666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666a164736f6c6343000806000a",
}

var ChainSpecificUtilHelperABI = ChainSpecificUtilHelperMetaData.ABI

var ChainSpecificUtilHelperBin = ChainSpecificUtilHelperMetaData.Bin

func DeployChainSpecificUtilHelper(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ChainSpecificUtilHelper, error) {
	parsed, err := ChainSpecificUtilHelperMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ChainSpecificUtilHelperBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ChainSpecificUtilHelper{ChainSpecificUtilHelperCaller: ChainSpecificUtilHelperCaller{contract: contract}, ChainSpecificUtilHelperTransactor: ChainSpecificUtilHelperTransactor{contract: contract}, ChainSpecificUtilHelperFilterer: ChainSpecificUtilHelperFilterer{contract: contract}}, nil
}

type ChainSpecificUtilHelper struct {
	address common.Address
	abi     abi.ABI
	ChainSpecificUtilHelperCaller
	ChainSpecificUtilHelperTransactor
	ChainSpecificUtilHelperFilterer
}

type ChainSpecificUtilHelperCaller struct {
	contract *bind.BoundContract
}

type ChainSpecificUtilHelperTransactor struct {
	contract *bind.BoundContract
}

type ChainSpecificUtilHelperFilterer struct {
	contract *bind.BoundContract
}

type ChainSpecificUtilHelperSession struct {
	Contract     *ChainSpecificUtilHelper
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type ChainSpecificUtilHelperCallerSession struct {
	Contract *ChainSpecificUtilHelperCaller
	CallOpts bind.CallOpts
}

type ChainSpecificUtilHelperTransactorSession struct {
	Contract     *ChainSpecificUtilHelperTransactor
	TransactOpts bind.TransactOpts
}

type ChainSpecificUtilHelperRaw struct {
	Contract *ChainSpecificUtilHelper
}

type ChainSpecificUtilHelperCallerRaw struct {
	Contract *ChainSpecificUtilHelperCaller
}

type ChainSpecificUtilHelperTransactorRaw struct {
	Contract *ChainSpecificUtilHelperTransactor
}

func NewChainSpecificUtilHelper(address common.Address, backend bind.ContractBackend) (*ChainSpecificUtilHelper, error) {
	abi, err := abi.JSON(strings.NewReader(ChainSpecificUtilHelperABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindChainSpecificUtilHelper(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ChainSpecificUtilHelper{address: address, abi: abi, ChainSpecificUtilHelperCaller: ChainSpecificUtilHelperCaller{contract: contract}, ChainSpecificUtilHelperTransactor: ChainSpecificUtilHelperTransactor{contract: contract}, ChainSpecificUtilHelperFilterer: ChainSpecificUtilHelperFilterer{contract: contract}}, nil
}

func NewChainSpecificUtilHelperCaller(address common.Address, caller bind.ContractCaller) (*ChainSpecificUtilHelperCaller, error) {
	contract, err := bindChainSpecificUtilHelper(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ChainSpecificUtilHelperCaller{contract: contract}, nil
}

func NewChainSpecificUtilHelperTransactor(address common.Address, transactor bind.ContractTransactor) (*ChainSpecificUtilHelperTransactor, error) {
	contract, err := bindChainSpecificUtilHelper(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ChainSpecificUtilHelperTransactor{contract: contract}, nil
}

func NewChainSpecificUtilHelperFilterer(address common.Address, filterer bind.ContractFilterer) (*ChainSpecificUtilHelperFilterer, error) {
	contract, err := bindChainSpecificUtilHelper(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ChainSpecificUtilHelperFilterer{contract: contract}, nil
}

func bindChainSpecificUtilHelper(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ChainSpecificUtilHelperMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_ChainSpecificUtilHelper *ChainSpecificUtilHelperRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ChainSpecificUtilHelper.Contract.ChainSpecificUtilHelperCaller.contract.Call(opts, result, method, params...)
}

func (_ChainSpecificUtilHelper *ChainSpecificUtilHelperRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChainSpecificUtilHelper.Contract.ChainSpecificUtilHelperTransactor.contract.Transfer(opts)
}

func (_ChainSpecificUtilHelper *ChainSpecificUtilHelperRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ChainSpecificUtilHelper.Contract.ChainSpecificUtilHelperTransactor.contract.Transact(opts, method, params...)
}

func (_ChainSpecificUtilHelper *ChainSpecificUtilHelperCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ChainSpecificUtilHelper.Contract.contract.Call(opts, result, method, params...)
}

func (_ChainSpecificUtilHelper *ChainSpecificUtilHelperTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChainSpecificUtilHelper.Contract.contract.Transfer(opts)
}

func (_ChainSpecificUtilHelper *ChainSpecificUtilHelperTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ChainSpecificUtilHelper.Contract.contract.Transact(opts, method, params...)
}

func (_ChainSpecificUtilHelper *ChainSpecificUtilHelperCaller) GetBlockNumber(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ChainSpecificUtilHelper.contract.Call(opts, &out, "getBlockNumber")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_ChainSpecificUtilHelper *ChainSpecificUtilHelperSession) GetBlockNumber() (*big.Int, error) {
	return _ChainSpecificUtilHelper.Contract.GetBlockNumber(&_ChainSpecificUtilHelper.CallOpts)
}

func (_ChainSpecificUtilHelper *ChainSpecificUtilHelperCallerSession) GetBlockNumber() (*big.Int, error) {
	return _ChainSpecificUtilHelper.Contract.GetBlockNumber(&_ChainSpecificUtilHelper.CallOpts)
}

func (_ChainSpecificUtilHelper *ChainSpecificUtilHelperCaller) GetBlockhash(opts *bind.CallOpts, blockNumber uint64) ([32]byte, error) {
	var out []interface{}
	err := _ChainSpecificUtilHelper.contract.Call(opts, &out, "getBlockhash", blockNumber)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_ChainSpecificUtilHelper *ChainSpecificUtilHelperSession) GetBlockhash(blockNumber uint64) ([32]byte, error) {
	return _ChainSpecificUtilHelper.Contract.GetBlockhash(&_ChainSpecificUtilHelper.CallOpts, blockNumber)
}

func (_ChainSpecificUtilHelper *ChainSpecificUtilHelperCallerSession) GetBlockhash(blockNumber uint64) ([32]byte, error) {
	return _ChainSpecificUtilHelper.Contract.GetBlockhash(&_ChainSpecificUtilHelper.CallOpts, blockNumber)
}

func (_ChainSpecificUtilHelper *ChainSpecificUtilHelperCaller) GetCurrentTxL1GasFees(opts *bind.CallOpts, txCallData string) (*big.Int, error) {
	var out []interface{}
	err := _ChainSpecificUtilHelper.contract.Call(opts, &out, "getCurrentTxL1GasFees", txCallData)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_ChainSpecificUtilHelper *ChainSpecificUtilHelperSession) GetCurrentTxL1GasFees(txCallData string) (*big.Int, error) {
	return _ChainSpecificUtilHelper.Contract.GetCurrentTxL1GasFees(&_ChainSpecificUtilHelper.CallOpts, txCallData)
}

func (_ChainSpecificUtilHelper *ChainSpecificUtilHelperCallerSession) GetCurrentTxL1GasFees(txCallData string) (*big.Int, error) {
	return _ChainSpecificUtilHelper.Contract.GetCurrentTxL1GasFees(&_ChainSpecificUtilHelper.CallOpts, txCallData)
}

func (_ChainSpecificUtilHelper *ChainSpecificUtilHelperCaller) GetL1CalldataGasCost(opts *bind.CallOpts, calldataSize *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _ChainSpecificUtilHelper.contract.Call(opts, &out, "getL1CalldataGasCost", calldataSize)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_ChainSpecificUtilHelper *ChainSpecificUtilHelperSession) GetL1CalldataGasCost(calldataSize *big.Int) (*big.Int, error) {
	return _ChainSpecificUtilHelper.Contract.GetL1CalldataGasCost(&_ChainSpecificUtilHelper.CallOpts, calldataSize)
}

func (_ChainSpecificUtilHelper *ChainSpecificUtilHelperCallerSession) GetL1CalldataGasCost(calldataSize *big.Int) (*big.Int, error) {
	return _ChainSpecificUtilHelper.Contract.GetL1CalldataGasCost(&_ChainSpecificUtilHelper.CallOpts, calldataSize)
}

func (_ChainSpecificUtilHelper *ChainSpecificUtilHelper) Address() common.Address {
	return _ChainSpecificUtilHelper.address
}

type ChainSpecificUtilHelperInterface interface {
	GetBlockNumber(opts *bind.CallOpts) (*big.Int, error)

	GetBlockhash(opts *bind.CallOpts, blockNumber uint64) ([32]byte, error)

	GetCurrentTxL1GasFees(opts *bind.CallOpts, txCallData string) (*big.Int, error)

	GetL1CalldataGasCost(opts *bind.CallOpts, calldataSize *big.Int) (*big.Int, error)

	Address() common.Address
}
