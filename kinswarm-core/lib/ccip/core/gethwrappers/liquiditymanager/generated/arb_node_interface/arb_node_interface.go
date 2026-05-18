// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package arb_node_interface

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

var NodeInterfaceMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"l2BlockNum\",\"type\":\"uint64\"}],\"name\":\"blockL1Num\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"l1BlockNum\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"size\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"leaf\",\"type\":\"uint64\"}],\"name\":\"constructOutboxProof\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"send\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"proof\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deposit\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"l2CallValue\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"excessFeeRefundAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"callValueRefundAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"estimateRetryableTicket\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"blockNum\",\"type\":\"uint64\"}],\"name\":\"findBatchContainingBlock\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"batch\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"contractCreation\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"gasEstimateComponents\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"gasEstimate\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"gasEstimateForL1\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"baseFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"l1BaseFeeEstimate\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"contractCreation\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"gasEstimateL1Component\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"gasEstimateForL1\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"baseFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"l1BaseFeeEstimate\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"}],\"name\":\"getL1Confirmations\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"confirmations\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"blockNum\",\"type\":\"uint64\"}],\"name\":\"l2BlockRangeForL1\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"firstBlock\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"lastBlock\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchNum\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"index\",\"type\":\"uint64\"}],\"name\":\"legacyLookupMessageBatchProof\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"proof\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256\",\"name\":\"path\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"l2Sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"l1Dest\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"l2Block\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"l1Block\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"calldataForL1\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nitroGenesisBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"number\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
}

var NodeInterfaceABI = NodeInterfaceMetaData.ABI

type NodeInterface struct {
	address common.Address
	abi     abi.ABI
	NodeInterfaceCaller
	NodeInterfaceTransactor
	NodeInterfaceFilterer
}

type NodeInterfaceCaller struct {
	contract *bind.BoundContract
}

type NodeInterfaceTransactor struct {
	contract *bind.BoundContract
}

type NodeInterfaceFilterer struct {
	contract *bind.BoundContract
}

type NodeInterfaceSession struct {
	Contract     *NodeInterface
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type NodeInterfaceCallerSession struct {
	Contract *NodeInterfaceCaller
	CallOpts bind.CallOpts
}

type NodeInterfaceTransactorSession struct {
	Contract     *NodeInterfaceTransactor
	TransactOpts bind.TransactOpts
}

type NodeInterfaceRaw struct {
	Contract *NodeInterface
}

type NodeInterfaceCallerRaw struct {
	Contract *NodeInterfaceCaller
}

type NodeInterfaceTransactorRaw struct {
	Contract *NodeInterfaceTransactor
}

func NewNodeInterface(address common.Address, backend bind.ContractBackend) (*NodeInterface, error) {
	abi, err := abi.JSON(strings.NewReader(NodeInterfaceABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindNodeInterface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NodeInterface{address: address, abi: abi, NodeInterfaceCaller: NodeInterfaceCaller{contract: contract}, NodeInterfaceTransactor: NodeInterfaceTransactor{contract: contract}, NodeInterfaceFilterer: NodeInterfaceFilterer{contract: contract}}, nil
}

func NewNodeInterfaceCaller(address common.Address, caller bind.ContractCaller) (*NodeInterfaceCaller, error) {
	contract, err := bindNodeInterface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NodeInterfaceCaller{contract: contract}, nil
}

func NewNodeInterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*NodeInterfaceTransactor, error) {
	contract, err := bindNodeInterface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NodeInterfaceTransactor{contract: contract}, nil
}

func NewNodeInterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*NodeInterfaceFilterer, error) {
	contract, err := bindNodeInterface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NodeInterfaceFilterer{contract: contract}, nil
}

func bindNodeInterface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := NodeInterfaceMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_NodeInterface *NodeInterfaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NodeInterface.Contract.NodeInterfaceCaller.contract.Call(opts, result, method, params...)
}

func (_NodeInterface *NodeInterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodeInterface.Contract.NodeInterfaceTransactor.contract.Transfer(opts)
}

func (_NodeInterface *NodeInterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NodeInterface.Contract.NodeInterfaceTransactor.contract.Transact(opts, method, params...)
}

func (_NodeInterface *NodeInterfaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NodeInterface.Contract.contract.Call(opts, result, method, params...)
}

func (_NodeInterface *NodeInterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodeInterface.Contract.contract.Transfer(opts)
}

func (_NodeInterface *NodeInterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NodeInterface.Contract.contract.Transact(opts, method, params...)
}

func (_NodeInterface *NodeInterfaceCaller) BlockL1Num(opts *bind.CallOpts, l2BlockNum uint64) (uint64, error) {
	var out []interface{}
	err := _NodeInterface.contract.Call(opts, &out, "blockL1Num", l2BlockNum)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_NodeInterface *NodeInterfaceSession) BlockL1Num(l2BlockNum uint64) (uint64, error) {
	return _NodeInterface.Contract.BlockL1Num(&_NodeInterface.CallOpts, l2BlockNum)
}

func (_NodeInterface *NodeInterfaceCallerSession) BlockL1Num(l2BlockNum uint64) (uint64, error) {
	return _NodeInterface.Contract.BlockL1Num(&_NodeInterface.CallOpts, l2BlockNum)
}

func (_NodeInterface *NodeInterfaceCaller) ConstructOutboxProof(opts *bind.CallOpts, size uint64, leaf uint64) (ConstructOutboxProof,

	error) {
	var out []interface{}
	err := _NodeInterface.contract.Call(opts, &out, "constructOutboxProof", size, leaf)

	outstruct := new(ConstructOutboxProof)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Send = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.Root = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Proof = *abi.ConvertType(out[2], new([][32]byte)).(*[][32]byte)

	return *outstruct, err

}

func (_NodeInterface *NodeInterfaceSession) ConstructOutboxProof(size uint64, leaf uint64) (ConstructOutboxProof,

	error) {
	return _NodeInterface.Contract.ConstructOutboxProof(&_NodeInterface.CallOpts, size, leaf)
}

func (_NodeInterface *NodeInterfaceCallerSession) ConstructOutboxProof(size uint64, leaf uint64) (ConstructOutboxProof,

	error) {
	return _NodeInterface.Contract.ConstructOutboxProof(&_NodeInterface.CallOpts, size, leaf)
}

func (_NodeInterface *NodeInterfaceCaller) FindBatchContainingBlock(opts *bind.CallOpts, blockNum uint64) (uint64, error) {
	var out []interface{}
	err := _NodeInterface.contract.Call(opts, &out, "findBatchContainingBlock", blockNum)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_NodeInterface *NodeInterfaceSession) FindBatchContainingBlock(blockNum uint64) (uint64, error) {
	return _NodeInterface.Contract.FindBatchContainingBlock(&_NodeInterface.CallOpts, blockNum)
}

func (_NodeInterface *NodeInterfaceCallerSession) FindBatchContainingBlock(blockNum uint64) (uint64, error) {
	return _NodeInterface.Contract.FindBatchContainingBlock(&_NodeInterface.CallOpts, blockNum)
}

func (_NodeInterface *NodeInterfaceCaller) GetL1Confirmations(opts *bind.CallOpts, blockHash [32]byte) (uint64, error) {
	var out []interface{}
	err := _NodeInterface.contract.Call(opts, &out, "getL1Confirmations", blockHash)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_NodeInterface *NodeInterfaceSession) GetL1Confirmations(blockHash [32]byte) (uint64, error) {
	return _NodeInterface.Contract.GetL1Confirmations(&_NodeInterface.CallOpts, blockHash)
}

func (_NodeInterface *NodeInterfaceCallerSession) GetL1Confirmations(blockHash [32]byte) (uint64, error) {
	return _NodeInterface.Contract.GetL1Confirmations(&_NodeInterface.CallOpts, blockHash)
}

func (_NodeInterface *NodeInterfaceCaller) L2BlockRangeForL1(opts *bind.CallOpts, blockNum uint64) (L2BlockRangeForL1,

	error) {
	var out []interface{}
	err := _NodeInterface.contract.Call(opts, &out, "l2BlockRangeForL1", blockNum)

	outstruct := new(L2BlockRangeForL1)
	if err != nil {
		return *outstruct, err
	}

	outstruct.FirstBlock = *abi.ConvertType(out[0], new(uint64)).(*uint64)
	outstruct.LastBlock = *abi.ConvertType(out[1], new(uint64)).(*uint64)

	return *outstruct, err

}

func (_NodeInterface *NodeInterfaceSession) L2BlockRangeForL1(blockNum uint64) (L2BlockRangeForL1,

	error) {
	return _NodeInterface.Contract.L2BlockRangeForL1(&_NodeInterface.CallOpts, blockNum)
}

func (_NodeInterface *NodeInterfaceCallerSession) L2BlockRangeForL1(blockNum uint64) (L2BlockRangeForL1,

	error) {
	return _NodeInterface.Contract.L2BlockRangeForL1(&_NodeInterface.CallOpts, blockNum)
}

func (_NodeInterface *NodeInterfaceCaller) LegacyLookupMessageBatchProof(opts *bind.CallOpts, batchNum *big.Int, index uint64) (LegacyLookupMessageBatchProof,

	error) {
	var out []interface{}
	err := _NodeInterface.contract.Call(opts, &out, "legacyLookupMessageBatchProof", batchNum, index)

	outstruct := new(LegacyLookupMessageBatchProof)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Proof = *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)
	outstruct.Path = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.L2Sender = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.L1Dest = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)
	outstruct.L2Block = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.L1Block = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.Timestamp = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)
	outstruct.Amount = *abi.ConvertType(out[7], new(*big.Int)).(**big.Int)
	outstruct.CalldataForL1 = *abi.ConvertType(out[8], new([]byte)).(*[]byte)

	return *outstruct, err

}

func (_NodeInterface *NodeInterfaceSession) LegacyLookupMessageBatchProof(batchNum *big.Int, index uint64) (LegacyLookupMessageBatchProof,

	error) {
	return _NodeInterface.Contract.LegacyLookupMessageBatchProof(&_NodeInterface.CallOpts, batchNum, index)
}

func (_NodeInterface *NodeInterfaceCallerSession) LegacyLookupMessageBatchProof(batchNum *big.Int, index uint64) (LegacyLookupMessageBatchProof,

	error) {
	return _NodeInterface.Contract.LegacyLookupMessageBatchProof(&_NodeInterface.CallOpts, batchNum, index)
}

func (_NodeInterface *NodeInterfaceCaller) NitroGenesisBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NodeInterface.contract.Call(opts, &out, "nitroGenesisBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_NodeInterface *NodeInterfaceSession) NitroGenesisBlock() (*big.Int, error) {
	return _NodeInterface.Contract.NitroGenesisBlock(&_NodeInterface.CallOpts)
}

func (_NodeInterface *NodeInterfaceCallerSession) NitroGenesisBlock() (*big.Int, error) {
	return _NodeInterface.Contract.NitroGenesisBlock(&_NodeInterface.CallOpts)
}

func (_NodeInterface *NodeInterfaceTransactor) EstimateRetryableTicket(opts *bind.TransactOpts, sender common.Address, deposit *big.Int, to common.Address, l2CallValue *big.Int, excessFeeRefundAddress common.Address, callValueRefundAddress common.Address, data []byte) (*types.Transaction, error) {
	return _NodeInterface.contract.Transact(opts, "estimateRetryableTicket", sender, deposit, to, l2CallValue, excessFeeRefundAddress, callValueRefundAddress, data)
}

func (_NodeInterface *NodeInterfaceSession) EstimateRetryableTicket(sender common.Address, deposit *big.Int, to common.Address, l2CallValue *big.Int, excessFeeRefundAddress common.Address, callValueRefundAddress common.Address, data []byte) (*types.Transaction, error) {
	return _NodeInterface.Contract.EstimateRetryableTicket(&_NodeInterface.TransactOpts, sender, deposit, to, l2CallValue, excessFeeRefundAddress, callValueRefundAddress, data)
}

func (_NodeInterface *NodeInterfaceTransactorSession) EstimateRetryableTicket(sender common.Address, deposit *big.Int, to common.Address, l2CallValue *big.Int, excessFeeRefundAddress common.Address, callValueRefundAddress common.Address, data []byte) (*types.Transaction, error) {
	return _NodeInterface.Contract.EstimateRetryableTicket(&_NodeInterface.TransactOpts, sender, deposit, to, l2CallValue, excessFeeRefundAddress, callValueRefundAddress, data)
}

func (_NodeInterface *NodeInterfaceTransactor) GasEstimateComponents(opts *bind.TransactOpts, to common.Address, contractCreation bool, data []byte) (*types.Transaction, error) {
	return _NodeInterface.contract.Transact(opts, "gasEstimateComponents", to, contractCreation, data)
}

func (_NodeInterface *NodeInterfaceSession) GasEstimateComponents(to common.Address, contractCreation bool, data []byte) (*types.Transaction, error) {
	return _NodeInterface.Contract.GasEstimateComponents(&_NodeInterface.TransactOpts, to, contractCreation, data)
}

func (_NodeInterface *NodeInterfaceTransactorSession) GasEstimateComponents(to common.Address, contractCreation bool, data []byte) (*types.Transaction, error) {
	return _NodeInterface.Contract.GasEstimateComponents(&_NodeInterface.TransactOpts, to, contractCreation, data)
}

func (_NodeInterface *NodeInterfaceTransactor) GasEstimateL1Component(opts *bind.TransactOpts, to common.Address, contractCreation bool, data []byte) (*types.Transaction, error) {
	return _NodeInterface.contract.Transact(opts, "gasEstimateL1Component", to, contractCreation, data)
}

func (_NodeInterface *NodeInterfaceSession) GasEstimateL1Component(to common.Address, contractCreation bool, data []byte) (*types.Transaction, error) {
	return _NodeInterface.Contract.GasEstimateL1Component(&_NodeInterface.TransactOpts, to, contractCreation, data)
}

func (_NodeInterface *NodeInterfaceTransactorSession) GasEstimateL1Component(to common.Address, contractCreation bool, data []byte) (*types.Transaction, error) {
	return _NodeInterface.Contract.GasEstimateL1Component(&_NodeInterface.TransactOpts, to, contractCreation, data)
}

type ConstructOutboxProof struct {
	Send  [32]byte
	Root  [32]byte
	Proof [][32]byte
}
type L2BlockRangeForL1 struct {
	FirstBlock uint64
	LastBlock  uint64
}
type LegacyLookupMessageBatchProof struct {
	Proof         [][32]byte
	Path          *big.Int
	L2Sender      common.Address
	L1Dest        common.Address
	L2Block       *big.Int
	L1Block       *big.Int
	Timestamp     *big.Int
	Amount        *big.Int
	CalldataForL1 []byte
}

func (_NodeInterface *NodeInterface) Address() common.Address {
	return _NodeInterface.address
}

type NodeInterfaceInterface interface {
	BlockL1Num(opts *bind.CallOpts, l2BlockNum uint64) (uint64, error)

	ConstructOutboxProof(opts *bind.CallOpts, size uint64, leaf uint64) (ConstructOutboxProof,

		error)

	FindBatchContainingBlock(opts *bind.CallOpts, blockNum uint64) (uint64, error)

	GetL1Confirmations(opts *bind.CallOpts, blockHash [32]byte) (uint64, error)

	L2BlockRangeForL1(opts *bind.CallOpts, blockNum uint64) (L2BlockRangeForL1,

		error)

	LegacyLookupMessageBatchProof(opts *bind.CallOpts, batchNum *big.Int, index uint64) (LegacyLookupMessageBatchProof,

		error)

	NitroGenesisBlock(opts *bind.CallOpts) (*big.Int, error)

	EstimateRetryableTicket(opts *bind.TransactOpts, sender common.Address, deposit *big.Int, to common.Address, l2CallValue *big.Int, excessFeeRefundAddress common.Address, callValueRefundAddress common.Address, data []byte) (*types.Transaction, error)

	GasEstimateComponents(opts *bind.TransactOpts, to common.Address, contractCreation bool, data []byte) (*types.Transaction, error)

	GasEstimateL1Component(opts *bind.TransactOpts, to common.Address, contractCreation bool, data []byte) (*types.Transaction, error)

	Address() common.Address
}
