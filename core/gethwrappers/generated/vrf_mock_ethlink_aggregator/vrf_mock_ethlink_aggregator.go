// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_mock_ethlink_aggregator

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

var VRFMockETHLINKAggregatorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"_answer\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"answer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"\",\"type\":\"uint80\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"ans\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"ans\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_blockTimestampDeduction\",\"type\":\"uint256\"}],\"name\":\"setBlockTimestampDeduction\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x6080604052600060015534801561001557600080fd5b506040516103383803806103388339810160408190526100349161003c565b600055610055565b60006020828403121561004e57600080fd5b5051919050565b6102d4806100646000396000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c806385bb7d691161005b57806385bb7d69146100e65780639a6fc8f5146100ef578063f0ad37df14610139578063feaf968c1461014e57600080fd5b8063313ce5671461008257806354fd4d50146100965780637284e416146100a7575b600080fd5b604051601281526020015b60405180910390f35b60015b60405190815260200161008d565b604080518082018252601881527f5652464d6f636b4554484c494e4b41676772656761746f7200000000000000006020820152905161008d9190610216565b61009960005481565b6101026100fd3660046101e3565b610156565b6040805169ffffffffffffffffffff968716815260208101959095528401929092526060830152909116608082015260a00161008d565b61014c6101473660046101ca565b600155565b005b610102610186565b6000806000806000600160005461016b6101b5565b6101736101b5565b9299919850965090945060019350915050565b6000806000806000600160005461019b6101b5565b6101a36101b5565b92989197509550909350600192509050565b6000600154426101c59190610289565b905090565b6000602082840312156101dc57600080fd5b5035919050565b6000602082840312156101f557600080fd5b813569ffffffffffffffffffff8116811461020f57600080fd5b9392505050565b600060208083528351808285015260005b8181101561024357858101830151858201604001528201610227565b81811115610255576000604083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016929092016040019392505050565b6000828210156102c2577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b50039056fea164736f6c6343000806000a",
}

var VRFMockETHLINKAggregatorABI = VRFMockETHLINKAggregatorMetaData.ABI

var VRFMockETHLINKAggregatorBin = VRFMockETHLINKAggregatorMetaData.Bin

func DeployVRFMockETHLINKAggregator(auth *bind.TransactOpts, backend bind.ContractBackend, _answer *big.Int) (common.Address, *types.Transaction, *VRFMockETHLINKAggregator, error) {
	parsed, err := VRFMockETHLINKAggregatorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFMockETHLINKAggregatorBin), backend, _answer)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFMockETHLINKAggregator{address: address, abi: *parsed, VRFMockETHLINKAggregatorCaller: VRFMockETHLINKAggregatorCaller{contract: contract}, VRFMockETHLINKAggregatorTransactor: VRFMockETHLINKAggregatorTransactor{contract: contract}, VRFMockETHLINKAggregatorFilterer: VRFMockETHLINKAggregatorFilterer{contract: contract}}, nil
}

type VRFMockETHLINKAggregator struct {
	address common.Address
	abi     abi.ABI
	VRFMockETHLINKAggregatorCaller
	VRFMockETHLINKAggregatorTransactor
	VRFMockETHLINKAggregatorFilterer
}

type VRFMockETHLINKAggregatorCaller struct {
	contract *bind.BoundContract
}

type VRFMockETHLINKAggregatorTransactor struct {
	contract *bind.BoundContract
}

type VRFMockETHLINKAggregatorFilterer struct {
	contract *bind.BoundContract
}

type VRFMockETHLINKAggregatorSession struct {
	Contract     *VRFMockETHLINKAggregator
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFMockETHLINKAggregatorCallerSession struct {
	Contract *VRFMockETHLINKAggregatorCaller
	CallOpts bind.CallOpts
}

type VRFMockETHLINKAggregatorTransactorSession struct {
	Contract     *VRFMockETHLINKAggregatorTransactor
	TransactOpts bind.TransactOpts
}

type VRFMockETHLINKAggregatorRaw struct {
	Contract *VRFMockETHLINKAggregator
}

type VRFMockETHLINKAggregatorCallerRaw struct {
	Contract *VRFMockETHLINKAggregatorCaller
}

type VRFMockETHLINKAggregatorTransactorRaw struct {
	Contract *VRFMockETHLINKAggregatorTransactor
}

func NewVRFMockETHLINKAggregator(address common.Address, backend bind.ContractBackend) (*VRFMockETHLINKAggregator, error) {
	abi, err := abi.JSON(strings.NewReader(VRFMockETHLINKAggregatorABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFMockETHLINKAggregator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFMockETHLINKAggregator{address: address, abi: abi, VRFMockETHLINKAggregatorCaller: VRFMockETHLINKAggregatorCaller{contract: contract}, VRFMockETHLINKAggregatorTransactor: VRFMockETHLINKAggregatorTransactor{contract: contract}, VRFMockETHLINKAggregatorFilterer: VRFMockETHLINKAggregatorFilterer{contract: contract}}, nil
}

func NewVRFMockETHLINKAggregatorCaller(address common.Address, caller bind.ContractCaller) (*VRFMockETHLINKAggregatorCaller, error) {
	contract, err := bindVRFMockETHLINKAggregator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFMockETHLINKAggregatorCaller{contract: contract}, nil
}

func NewVRFMockETHLINKAggregatorTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFMockETHLINKAggregatorTransactor, error) {
	contract, err := bindVRFMockETHLINKAggregator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFMockETHLINKAggregatorTransactor{contract: contract}, nil
}

func NewVRFMockETHLINKAggregatorFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFMockETHLINKAggregatorFilterer, error) {
	contract, err := bindVRFMockETHLINKAggregator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFMockETHLINKAggregatorFilterer{contract: contract}, nil
}

func bindVRFMockETHLINKAggregator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFMockETHLINKAggregatorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFMockETHLINKAggregator *VRFMockETHLINKAggregatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFMockETHLINKAggregator.Contract.VRFMockETHLINKAggregatorCaller.contract.Call(opts, result, method, params...)
}

func (_VRFMockETHLINKAggregator *VRFMockETHLINKAggregatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFMockETHLINKAggregator.Contract.VRFMockETHLINKAggregatorTransactor.contract.Transfer(opts)
}

func (_VRFMockETHLINKAggregator *VRFMockETHLINKAggregatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFMockETHLINKAggregator.Contract.VRFMockETHLINKAggregatorTransactor.contract.Transact(opts, method, params...)
}

func (_VRFMockETHLINKAggregator *VRFMockETHLINKAggregatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFMockETHLINKAggregator.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFMockETHLINKAggregator *VRFMockETHLINKAggregatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFMockETHLINKAggregator.Contract.contract.Transfer(opts)
}

func (_VRFMockETHLINKAggregator *VRFMockETHLINKAggregatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFMockETHLINKAggregator.Contract.contract.Transact(opts, method, params...)
}

func (_VRFMockETHLINKAggregator *VRFMockETHLINKAggregatorCaller) Answer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFMockETHLINKAggregator.contract.Call(opts, &out, "answer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFMockETHLINKAggregator *VRFMockETHLINKAggregatorSession) Answer() (*big.Int, error) {
	return _VRFMockETHLINKAggregator.Contract.Answer(&_VRFMockETHLINKAggregator.CallOpts)
}

func (_VRFMockETHLINKAggregator *VRFMockETHLINKAggregatorCallerSession) Answer() (*big.Int, error) {
	return _VRFMockETHLINKAggregator.Contract.Answer(&_VRFMockETHLINKAggregator.CallOpts)
}

func (_VRFMockETHLINKAggregator *VRFMockETHLINKAggregatorCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _VRFMockETHLINKAggregator.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_VRFMockETHLINKAggregator *VRFMockETHLINKAggregatorSession) Decimals() (uint8, error) {
	return _VRFMockETHLINKAggregator.Contract.Decimals(&_VRFMockETHLINKAggregator.CallOpts)
}

func (_VRFMockETHLINKAggregator *VRFMockETHLINKAggregatorCallerSession) Decimals() (uint8, error) {
	return _VRFMockETHLINKAggregator.Contract.Decimals(&_VRFMockETHLINKAggregator.CallOpts)
}

func (_VRFMockETHLINKAggregator *VRFMockETHLINKAggregatorCaller) Description(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _VRFMockETHLINKAggregator.contract.Call(opts, &out, "description")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_VRFMockETHLINKAggregator *VRFMockETHLINKAggregatorSession) Description() (string, error) {
	return _VRFMockETHLINKAggregator.Contract.Description(&_VRFMockETHLINKAggregator.CallOpts)
}

func (_VRFMockETHLINKAggregator *VRFMockETHLINKAggregatorCallerSession) Description() (string, error) {
	return _VRFMockETHLINKAggregator.Contract.Description(&_VRFMockETHLINKAggregator.CallOpts)
}

func (_VRFMockETHLINKAggregator *VRFMockETHLINKAggregatorCaller) GetRoundData(opts *bind.CallOpts, arg0 *big.Int) (GetRoundData,

	error) {
	var out []interface{}
	err := _VRFMockETHLINKAggregator.contract.Call(opts, &out, "getRoundData", arg0)

	outstruct := new(GetRoundData)
	if err != nil {
		return *outstruct, err
	}

	outstruct.RoundId = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Ans = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.StartedAt = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.UpdatedAt = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.AnsweredInRound = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_VRFMockETHLINKAggregator *VRFMockETHLINKAggregatorSession) GetRoundData(arg0 *big.Int) (GetRoundData,

	error) {
	return _VRFMockETHLINKAggregator.Contract.GetRoundData(&_VRFMockETHLINKAggregator.CallOpts, arg0)
}

func (_VRFMockETHLINKAggregator *VRFMockETHLINKAggregatorCallerSession) GetRoundData(arg0 *big.Int) (GetRoundData,

	error) {
	return _VRFMockETHLINKAggregator.Contract.GetRoundData(&_VRFMockETHLINKAggregator.CallOpts, arg0)
}

func (_VRFMockETHLINKAggregator *VRFMockETHLINKAggregatorCaller) LatestRoundData(opts *bind.CallOpts) (LatestRoundData,

	error) {
	var out []interface{}
	err := _VRFMockETHLINKAggregator.contract.Call(opts, &out, "latestRoundData")

	outstruct := new(LatestRoundData)
	if err != nil {
		return *outstruct, err
	}

	outstruct.RoundId = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Ans = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.StartedAt = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.UpdatedAt = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.AnsweredInRound = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_VRFMockETHLINKAggregator *VRFMockETHLINKAggregatorSession) LatestRoundData() (LatestRoundData,

	error) {
	return _VRFMockETHLINKAggregator.Contract.LatestRoundData(&_VRFMockETHLINKAggregator.CallOpts)
}

func (_VRFMockETHLINKAggregator *VRFMockETHLINKAggregatorCallerSession) LatestRoundData() (LatestRoundData,

	error) {
	return _VRFMockETHLINKAggregator.Contract.LatestRoundData(&_VRFMockETHLINKAggregator.CallOpts)
}

func (_VRFMockETHLINKAggregator *VRFMockETHLINKAggregatorCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFMockETHLINKAggregator.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFMockETHLINKAggregator *VRFMockETHLINKAggregatorSession) Version() (*big.Int, error) {
	return _VRFMockETHLINKAggregator.Contract.Version(&_VRFMockETHLINKAggregator.CallOpts)
}

func (_VRFMockETHLINKAggregator *VRFMockETHLINKAggregatorCallerSession) Version() (*big.Int, error) {
	return _VRFMockETHLINKAggregator.Contract.Version(&_VRFMockETHLINKAggregator.CallOpts)
}

func (_VRFMockETHLINKAggregator *VRFMockETHLINKAggregatorTransactor) SetBlockTimestampDeduction(opts *bind.TransactOpts, _blockTimestampDeduction *big.Int) (*types.Transaction, error) {
	return _VRFMockETHLINKAggregator.contract.Transact(opts, "setBlockTimestampDeduction", _blockTimestampDeduction)
}

func (_VRFMockETHLINKAggregator *VRFMockETHLINKAggregatorSession) SetBlockTimestampDeduction(_blockTimestampDeduction *big.Int) (*types.Transaction, error) {
	return _VRFMockETHLINKAggregator.Contract.SetBlockTimestampDeduction(&_VRFMockETHLINKAggregator.TransactOpts, _blockTimestampDeduction)
}

func (_VRFMockETHLINKAggregator *VRFMockETHLINKAggregatorTransactorSession) SetBlockTimestampDeduction(_blockTimestampDeduction *big.Int) (*types.Transaction, error) {
	return _VRFMockETHLINKAggregator.Contract.SetBlockTimestampDeduction(&_VRFMockETHLINKAggregator.TransactOpts, _blockTimestampDeduction)
}

type GetRoundData struct {
	RoundId         *big.Int
	Ans             *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}
type LatestRoundData struct {
	RoundId         *big.Int
	Ans             *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}

func (_VRFMockETHLINKAggregator *VRFMockETHLINKAggregator) Address() common.Address {
	return _VRFMockETHLINKAggregator.address
}

type VRFMockETHLINKAggregatorInterface interface {
	Answer(opts *bind.CallOpts) (*big.Int, error)

	Decimals(opts *bind.CallOpts) (uint8, error)

	Description(opts *bind.CallOpts) (string, error)

	GetRoundData(opts *bind.CallOpts, arg0 *big.Int) (GetRoundData,

		error)

	LatestRoundData(opts *bind.CallOpts) (LatestRoundData,

		error)

	Version(opts *bind.CallOpts) (*big.Int, error)

	SetBlockTimestampDeduction(opts *bind.TransactOpts, _blockTimestampDeduction *big.Int) (*types.Transaction, error)

	Address() common.Address
}
