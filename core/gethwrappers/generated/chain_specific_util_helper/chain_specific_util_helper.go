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
	Bin: "0x608060405234801561001057600080fd5b50610c1a806100206000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c806342cbb15c1461005157806397e329a91461006b578063b778b1121461007e578063da9027ef14610091575b600080fd5b6100596100a4565b60405190815260200160405180910390f35b6100596100793660046108bf565b6100b3565b61005961008c36600461085c565b6100c4565b61005961009f36600461078d565b6100cf565b60006100ae6100da565b905090565b60006100be82610177565b92915050565b60006100be8261027d565b60006100be82610355565b6000466100e681610441565b1561017057606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b815260040160206040518083038186803b15801561013257600080fd5b505afa158015610146573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061016a9190610774565b91505090565b4391505090565b60004661018381610441565b1561026d576101008367ffffffffffffffff1661019e6100da565b6101a89190610b20565b11806101c557506101b76100da565b8367ffffffffffffffff1610155b156101d35750600092915050565b6040517f2b407a8200000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84166004820152606490632b407a82906024015b60206040518083038186803b15801561022e57600080fd5b505afa158015610242573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102669190610774565b9392505050565b505067ffffffffffffffff164090565b60004661028981610441565b15610335576000606c73ffffffffffffffffffffffffffffffffffffffff166341b247a86040518163ffffffff1660e01b815260040160c06040518083038186803b1580156102d757600080fd5b505afa1580156102eb573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061030f9190610875565b5050505091505083608c6103239190610969565b61032d9082610ae3565b949350505050565b61033e81610464565b1561034c576102668361049e565b50600092915050565b60004661036181610441565b156103ad57606c73ffffffffffffffffffffffffffffffffffffffff1663c6f7de0e6040518163ffffffff1660e01b815260040160206040518083038186803b15801561022e57600080fd5b6103b681610464565b1561034c5773420000000000000000000000000000000000000f73ffffffffffffffffffffffffffffffffffffffff166349948e0e84604051806080016040528060488152602001610bc6604891396040516020016104169291906108e9565b6040516020818303038152906040526040518263ffffffff1660e01b81526004016102169190610918565b600061a4b1821480610455575062066eed82145b806100be57505062066eee1490565b6000600a82148061047657506101a482145b80610483575062aa37dc82145b8061048f575061210582145b806100be57505062014a331490565b60008073420000000000000000000000000000000000000f73ffffffffffffffffffffffffffffffffffffffff1663519b4bd36040518163ffffffff1660e01b815260040160206040518083038186803b1580156104fb57600080fd5b505afa15801561050f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105339190610774565b90506000806105428186610b20565b90506000610551826010610ae3565b61055c846004610ae3565b6105669190610969565b9050600073420000000000000000000000000000000000000f73ffffffffffffffffffffffffffffffffffffffff16630c18c1626040518163ffffffff1660e01b815260040160206040518083038186803b1580156105c457600080fd5b505afa1580156105d8573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105fc9190610774565b9050600073420000000000000000000000000000000000000f73ffffffffffffffffffffffffffffffffffffffff1663f45e65d86040518163ffffffff1660e01b815260040160206040518083038186803b15801561065a57600080fd5b505afa15801561066e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106929190610774565b9050600073420000000000000000000000000000000000000f73ffffffffffffffffffffffffffffffffffffffff1663313ce5676040518163ffffffff1660e01b815260040160206040518083038186803b1580156106f057600080fd5b505afa158015610704573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107289190610774565b9050600061073782600a610a1d565b9050600081846107478789610969565b610751908c610ae3565b61075b9190610ae3565b6107659190610981565b9b9a5050505050505050505050565b60006020828403121561078657600080fd5b5051919050565b60006020828403121561079f57600080fd5b813567ffffffffffffffff808211156107b757600080fd5b818401915084601f8301126107cb57600080fd5b8135818111156107dd576107dd610b96565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f0116810190838211818310171561082357610823610b96565b8160405282815287602084870101111561083c57600080fd5b826020860160208301376000928101602001929092525095945050505050565b60006020828403121561086e57600080fd5b5035919050565b60008060008060008060c0878903121561088e57600080fd5b865195506020870151945060408701519350606087015192506080870151915060a087015190509295509295509295565b6000602082840312156108d157600080fd5b813567ffffffffffffffff8116811461026657600080fd5b600083516108fb818460208801610b37565b83519083019061090f818360208801610b37565b01949350505050565b6020815260008251806020840152610937816040850160208701610b37565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169190910160400192915050565b6000821982111561097c5761097c610b67565b500190565b6000826109b7577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b600181815b80851115610a1557817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048211156109fb576109fb610b67565b80851615610a0857918102915b93841c93908002906109c1565b509250929050565b60006102668383600082610a33575060016100be565b81610a40575060006100be565b8160018114610a565760028114610a6057610a7c565b60019150506100be565b60ff841115610a7157610a71610b67565b50506001821b6100be565b5060208310610133831016604e8410600b8410161715610a9f575081810a6100be565b610aa983836109bc565b807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04821115610adb57610adb610b67565b029392505050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615610b1b57610b1b610b67565b500290565b600082821015610b3257610b32610b67565b500390565b60005b83811015610b52578181015183820152602001610b3a565b83811115610b61576000848401525b50505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfe307866666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666a164736f6c6343000806000a",
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
	return address, tx, &ChainSpecificUtilHelper{address: address, abi: *parsed, ChainSpecificUtilHelperCaller: ChainSpecificUtilHelperCaller{contract: contract}, ChainSpecificUtilHelperTransactor: ChainSpecificUtilHelperTransactor{contract: contract}, ChainSpecificUtilHelperFilterer: ChainSpecificUtilHelperFilterer{contract: contract}}, nil
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
