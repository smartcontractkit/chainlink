// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package mock_v3_aggregator_contract

import (
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
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// MockV3AggregatorContractABI is the input ABI used to generate the binding from.
const MockV3AggregatorContractABI = "[{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"_decimals\",\"type\":\"uint8\"},{\"internalType\":\"int256\",\"name\":\"_initialAnswer\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"current\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"name\":\"AnswerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"getAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"_roundId\",\"type\":\"uint80\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"getTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"_answer\",\"type\":\"int256\"}],\"name\":\"updateAnswer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"_roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"_answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_startedAt\",\"type\":\"uint256\"}],\"name\":\"updateRoundData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// MockV3AggregatorContractBin is the compiled bytecode used for deploying new contracts.
var MockV3AggregatorContractBin = "0x608060405234801561001057600080fd5b506040516105113803806105118339818101604052604081101561003357600080fd5b5080516020909101516000805460ff191660ff84161790556100548161005b565b50506100a2565b600181815542600281905560038054909201808355600090815260046020908152604080832095909555835482526005815284822083905592548152600690925291902055565b610460806100b16000396000f3fe608060405234801561001057600080fd5b50600436106100d45760003560e01c80638205bf6a11610081578063b5ab58dc1161005b578063b5ab58dc14610273578063b633620c14610290578063feaf968c146102ad576100d4565b80638205bf6a146101db5780639a6fc8f5146101e3578063a87a20ce14610256576100d4565b806354fd4d50116100b257806354fd4d501461014e578063668a0f02146101565780637284e4161461015e576100d4565b8063313ce567146100d95780634aa2011f146100f757806350d25bcd14610134575b600080fd5b6100e16102b5565b6040805160ff9092168252519081900360200190f35b6101326004803603608081101561010d57600080fd5b5069ffffffffffffffffffff81351690602081013590604081013590606001356102be565b005b61013c61030b565b60408051918252519081900360200190f35b61013c610311565b61013c610316565b61016661031c565b6040805160208082528351818301528351919283929083019185019080838360005b838110156101a0578181015183820152602001610188565b50505050905090810190601f1680156101cd5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b61013c610353565b61020c600480360360208110156101f957600080fd5b503569ffffffffffffffffffff16610359565b604051808669ffffffffffffffffffff1681526020018581526020018481526020018381526020018269ffffffffffffffffffff1681526020019550505050505060405180910390f35b6101326004803603602081101561026c57600080fd5b5035610392565b61013c6004803603602081101561028957600080fd5b50356103d9565b61013c600480360360208110156102a657600080fd5b50356103eb565b61020c6103fd565b60005460ff1681565b69ffffffffffffffffffff90931660038181556001849055600283905560009182526004602090815260408084209590955581548352600581528483209390935554815260069091522055565b60015481565b600081565b60035481565b60408051808201909152601f81527f76302e362f74657374732f4d6f636b563341676772656761746f722e736f6c00602082015290565b60025481565b69ffffffffffffffffffff8116600090815260046020908152604080832054600683528184205460059093529220549293919290918490565b600181815542600281905560038054909201808355600090815260046020908152604080832095909555835482526005815284822083905592548152600690925291902055565b60046020526000908152604090205481565b60056020526000908152604090205481565b6003546000818152600460209081526040808320546006835281842054600590935292205483909192939456fea2646970667358221220ecf1c50e0f78cd131fb708022b7a4f2d2de0408537205a8d45c5a41fdbc0ad4d64736f6c63430007060033"

// DeployMockV3AggregatorContract deploys a new Ethereum contract, binding an instance of MockV3AggregatorContract to it.
func DeployMockV3AggregatorContract(auth *bind.TransactOpts, backend bind.ContractBackend, _decimals uint8, _initialAnswer *big.Int) (common.Address, *types.Transaction, *MockV3AggregatorContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MockV3AggregatorContractABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(MockV3AggregatorContractBin), backend, _decimals, _initialAnswer)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MockV3AggregatorContract{MockV3AggregatorContractCaller: MockV3AggregatorContractCaller{contract: contract}, MockV3AggregatorContractTransactor: MockV3AggregatorContractTransactor{contract: contract}, MockV3AggregatorContractFilterer: MockV3AggregatorContractFilterer{contract: contract}}, nil
}

// MockV3AggregatorContract is an auto generated Go binding around an Ethereum contract.
type MockV3AggregatorContract struct {
	MockV3AggregatorContractCaller     // Read-only binding to the contract
	MockV3AggregatorContractTransactor // Write-only binding to the contract
	MockV3AggregatorContractFilterer   // Log filterer for contract events
}

// MockV3AggregatorContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type MockV3AggregatorContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockV3AggregatorContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MockV3AggregatorContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockV3AggregatorContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MockV3AggregatorContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockV3AggregatorContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MockV3AggregatorContractSession struct {
	Contract     *MockV3AggregatorContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts             // Call options to use throughout this session
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// MockV3AggregatorContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MockV3AggregatorContractCallerSession struct {
	Contract *MockV3AggregatorContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                   // Call options to use throughout this session
}

// MockV3AggregatorContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MockV3AggregatorContractTransactorSession struct {
	Contract     *MockV3AggregatorContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                   // Transaction auth options to use throughout this session
}

// MockV3AggregatorContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type MockV3AggregatorContractRaw struct {
	Contract *MockV3AggregatorContract // Generic contract binding to access the raw methods on
}

// MockV3AggregatorContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MockV3AggregatorContractCallerRaw struct {
	Contract *MockV3AggregatorContractCaller // Generic read-only contract binding to access the raw methods on
}

// MockV3AggregatorContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MockV3AggregatorContractTransactorRaw struct {
	Contract *MockV3AggregatorContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMockV3AggregatorContract creates a new instance of MockV3AggregatorContract, bound to a specific deployed contract.
func NewMockV3AggregatorContract(address common.Address, backend bind.ContractBackend) (*MockV3AggregatorContract, error) {
	contract, err := bindMockV3AggregatorContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MockV3AggregatorContract{MockV3AggregatorContractCaller: MockV3AggregatorContractCaller{contract: contract}, MockV3AggregatorContractTransactor: MockV3AggregatorContractTransactor{contract: contract}, MockV3AggregatorContractFilterer: MockV3AggregatorContractFilterer{contract: contract}}, nil
}

// NewMockV3AggregatorContractCaller creates a new read-only instance of MockV3AggregatorContract, bound to a specific deployed contract.
func NewMockV3AggregatorContractCaller(address common.Address, caller bind.ContractCaller) (*MockV3AggregatorContractCaller, error) {
	contract, err := bindMockV3AggregatorContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MockV3AggregatorContractCaller{contract: contract}, nil
}

// NewMockV3AggregatorContractTransactor creates a new write-only instance of MockV3AggregatorContract, bound to a specific deployed contract.
func NewMockV3AggregatorContractTransactor(address common.Address, transactor bind.ContractTransactor) (*MockV3AggregatorContractTransactor, error) {
	contract, err := bindMockV3AggregatorContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MockV3AggregatorContractTransactor{contract: contract}, nil
}

// NewMockV3AggregatorContractFilterer creates a new log filterer instance of MockV3AggregatorContract, bound to a specific deployed contract.
func NewMockV3AggregatorContractFilterer(address common.Address, filterer bind.ContractFilterer) (*MockV3AggregatorContractFilterer, error) {
	contract, err := bindMockV3AggregatorContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MockV3AggregatorContractFilterer{contract: contract}, nil
}

// bindMockV3AggregatorContract binds a generic wrapper to an already deployed contract.
func bindMockV3AggregatorContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MockV3AggregatorContractABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MockV3AggregatorContract *MockV3AggregatorContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockV3AggregatorContract.Contract.MockV3AggregatorContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MockV3AggregatorContract *MockV3AggregatorContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockV3AggregatorContract.Contract.MockV3AggregatorContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MockV3AggregatorContract *MockV3AggregatorContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockV3AggregatorContract.Contract.MockV3AggregatorContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MockV3AggregatorContract *MockV3AggregatorContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockV3AggregatorContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MockV3AggregatorContract *MockV3AggregatorContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockV3AggregatorContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MockV3AggregatorContract *MockV3AggregatorContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockV3AggregatorContract.Contract.contract.Transact(opts, method, params...)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_MockV3AggregatorContract *MockV3AggregatorContractCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _MockV3AggregatorContract.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_MockV3AggregatorContract *MockV3AggregatorContractSession) Decimals() (uint8, error) {
	return _MockV3AggregatorContract.Contract.Decimals(&_MockV3AggregatorContract.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_MockV3AggregatorContract *MockV3AggregatorContractCallerSession) Decimals() (uint8, error) {
	return _MockV3AggregatorContract.Contract.Decimals(&_MockV3AggregatorContract.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_MockV3AggregatorContract *MockV3AggregatorContractCaller) Description(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MockV3AggregatorContract.contract.Call(opts, &out, "description")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_MockV3AggregatorContract *MockV3AggregatorContractSession) Description() (string, error) {
	return _MockV3AggregatorContract.Contract.Description(&_MockV3AggregatorContract.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_MockV3AggregatorContract *MockV3AggregatorContractCallerSession) Description() (string, error) {
	return _MockV3AggregatorContract.Contract.Description(&_MockV3AggregatorContract.CallOpts)
}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 ) view returns(int256)
func (_MockV3AggregatorContract *MockV3AggregatorContractCaller) GetAnswer(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _MockV3AggregatorContract.contract.Call(opts, &out, "getAnswer", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 ) view returns(int256)
func (_MockV3AggregatorContract *MockV3AggregatorContractSession) GetAnswer(arg0 *big.Int) (*big.Int, error) {
	return _MockV3AggregatorContract.Contract.GetAnswer(&_MockV3AggregatorContract.CallOpts, arg0)
}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 ) view returns(int256)
func (_MockV3AggregatorContract *MockV3AggregatorContractCallerSession) GetAnswer(arg0 *big.Int) (*big.Int, error) {
	return _MockV3AggregatorContract.Contract.GetAnswer(&_MockV3AggregatorContract.CallOpts, arg0)
}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_MockV3AggregatorContract *MockV3AggregatorContractCaller) GetRoundData(opts *bind.CallOpts, _roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	var out []interface{}
	err := _MockV3AggregatorContract.contract.Call(opts, &out, "getRoundData", _roundId)

	outstruct := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RoundId = out[0].(*big.Int)
	outstruct.Answer = out[1].(*big.Int)
	outstruct.StartedAt = out[2].(*big.Int)
	outstruct.UpdatedAt = out[3].(*big.Int)
	outstruct.AnsweredInRound = out[4].(*big.Int)

	return *outstruct, err

}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_MockV3AggregatorContract *MockV3AggregatorContractSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _MockV3AggregatorContract.Contract.GetRoundData(&_MockV3AggregatorContract.CallOpts, _roundId)
}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_MockV3AggregatorContract *MockV3AggregatorContractCallerSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _MockV3AggregatorContract.Contract.GetRoundData(&_MockV3AggregatorContract.CallOpts, _roundId)
}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 ) view returns(uint256)
func (_MockV3AggregatorContract *MockV3AggregatorContractCaller) GetTimestamp(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _MockV3AggregatorContract.contract.Call(opts, &out, "getTimestamp", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 ) view returns(uint256)
func (_MockV3AggregatorContract *MockV3AggregatorContractSession) GetTimestamp(arg0 *big.Int) (*big.Int, error) {
	return _MockV3AggregatorContract.Contract.GetTimestamp(&_MockV3AggregatorContract.CallOpts, arg0)
}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 ) view returns(uint256)
func (_MockV3AggregatorContract *MockV3AggregatorContractCallerSession) GetTimestamp(arg0 *big.Int) (*big.Int, error) {
	return _MockV3AggregatorContract.Contract.GetTimestamp(&_MockV3AggregatorContract.CallOpts, arg0)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_MockV3AggregatorContract *MockV3AggregatorContractCaller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MockV3AggregatorContract.contract.Call(opts, &out, "latestAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_MockV3AggregatorContract *MockV3AggregatorContractSession) LatestAnswer() (*big.Int, error) {
	return _MockV3AggregatorContract.Contract.LatestAnswer(&_MockV3AggregatorContract.CallOpts)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_MockV3AggregatorContract *MockV3AggregatorContractCallerSession) LatestAnswer() (*big.Int, error) {
	return _MockV3AggregatorContract.Contract.LatestAnswer(&_MockV3AggregatorContract.CallOpts)
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_MockV3AggregatorContract *MockV3AggregatorContractCaller) LatestRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MockV3AggregatorContract.contract.Call(opts, &out, "latestRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_MockV3AggregatorContract *MockV3AggregatorContractSession) LatestRound() (*big.Int, error) {
	return _MockV3AggregatorContract.Contract.LatestRound(&_MockV3AggregatorContract.CallOpts)
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_MockV3AggregatorContract *MockV3AggregatorContractCallerSession) LatestRound() (*big.Int, error) {
	return _MockV3AggregatorContract.Contract.LatestRound(&_MockV3AggregatorContract.CallOpts)
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_MockV3AggregatorContract *MockV3AggregatorContractCaller) LatestRoundData(opts *bind.CallOpts) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	var out []interface{}
	err := _MockV3AggregatorContract.contract.Call(opts, &out, "latestRoundData")

	outstruct := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RoundId = out[0].(*big.Int)
	outstruct.Answer = out[1].(*big.Int)
	outstruct.StartedAt = out[2].(*big.Int)
	outstruct.UpdatedAt = out[3].(*big.Int)
	outstruct.AnsweredInRound = out[4].(*big.Int)

	return *outstruct, err

}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_MockV3AggregatorContract *MockV3AggregatorContractSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _MockV3AggregatorContract.Contract.LatestRoundData(&_MockV3AggregatorContract.CallOpts)
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_MockV3AggregatorContract *MockV3AggregatorContractCallerSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _MockV3AggregatorContract.Contract.LatestRoundData(&_MockV3AggregatorContract.CallOpts)
}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256)
func (_MockV3AggregatorContract *MockV3AggregatorContractCaller) LatestTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MockV3AggregatorContract.contract.Call(opts, &out, "latestTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256)
func (_MockV3AggregatorContract *MockV3AggregatorContractSession) LatestTimestamp() (*big.Int, error) {
	return _MockV3AggregatorContract.Contract.LatestTimestamp(&_MockV3AggregatorContract.CallOpts)
}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256)
func (_MockV3AggregatorContract *MockV3AggregatorContractCallerSession) LatestTimestamp() (*big.Int, error) {
	return _MockV3AggregatorContract.Contract.LatestTimestamp(&_MockV3AggregatorContract.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_MockV3AggregatorContract *MockV3AggregatorContractCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MockV3AggregatorContract.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_MockV3AggregatorContract *MockV3AggregatorContractSession) Version() (*big.Int, error) {
	return _MockV3AggregatorContract.Contract.Version(&_MockV3AggregatorContract.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_MockV3AggregatorContract *MockV3AggregatorContractCallerSession) Version() (*big.Int, error) {
	return _MockV3AggregatorContract.Contract.Version(&_MockV3AggregatorContract.CallOpts)
}

// UpdateAnswer is a paid mutator transaction binding the contract method 0xa87a20ce.
//
// Solidity: function updateAnswer(int256 _answer) returns()
func (_MockV3AggregatorContract *MockV3AggregatorContractTransactor) UpdateAnswer(opts *bind.TransactOpts, _answer *big.Int) (*types.Transaction, error) {
	return _MockV3AggregatorContract.contract.Transact(opts, "updateAnswer", _answer)
}

// UpdateAnswer is a paid mutator transaction binding the contract method 0xa87a20ce.
//
// Solidity: function updateAnswer(int256 _answer) returns()
func (_MockV3AggregatorContract *MockV3AggregatorContractSession) UpdateAnswer(_answer *big.Int) (*types.Transaction, error) {
	return _MockV3AggregatorContract.Contract.UpdateAnswer(&_MockV3AggregatorContract.TransactOpts, _answer)
}

// UpdateAnswer is a paid mutator transaction binding the contract method 0xa87a20ce.
//
// Solidity: function updateAnswer(int256 _answer) returns()
func (_MockV3AggregatorContract *MockV3AggregatorContractTransactorSession) UpdateAnswer(_answer *big.Int) (*types.Transaction, error) {
	return _MockV3AggregatorContract.Contract.UpdateAnswer(&_MockV3AggregatorContract.TransactOpts, _answer)
}

// UpdateRoundData is a paid mutator transaction binding the contract method 0x4aa2011f.
//
// Solidity: function updateRoundData(uint80 _roundId, int256 _answer, uint256 _timestamp, uint256 _startedAt) returns()
func (_MockV3AggregatorContract *MockV3AggregatorContractTransactor) UpdateRoundData(opts *bind.TransactOpts, _roundId *big.Int, _answer *big.Int, _timestamp *big.Int, _startedAt *big.Int) (*types.Transaction, error) {
	return _MockV3AggregatorContract.contract.Transact(opts, "updateRoundData", _roundId, _answer, _timestamp, _startedAt)
}

// UpdateRoundData is a paid mutator transaction binding the contract method 0x4aa2011f.
//
// Solidity: function updateRoundData(uint80 _roundId, int256 _answer, uint256 _timestamp, uint256 _startedAt) returns()
func (_MockV3AggregatorContract *MockV3AggregatorContractSession) UpdateRoundData(_roundId *big.Int, _answer *big.Int, _timestamp *big.Int, _startedAt *big.Int) (*types.Transaction, error) {
	return _MockV3AggregatorContract.Contract.UpdateRoundData(&_MockV3AggregatorContract.TransactOpts, _roundId, _answer, _timestamp, _startedAt)
}

// UpdateRoundData is a paid mutator transaction binding the contract method 0x4aa2011f.
//
// Solidity: function updateRoundData(uint80 _roundId, int256 _answer, uint256 _timestamp, uint256 _startedAt) returns()
func (_MockV3AggregatorContract *MockV3AggregatorContractTransactorSession) UpdateRoundData(_roundId *big.Int, _answer *big.Int, _timestamp *big.Int, _startedAt *big.Int) (*types.Transaction, error) {
	return _MockV3AggregatorContract.Contract.UpdateRoundData(&_MockV3AggregatorContract.TransactOpts, _roundId, _answer, _timestamp, _startedAt)
}

// MockV3AggregatorContractAnswerUpdatedIterator is returned from FilterAnswerUpdated and is used to iterate over the raw logs and unpacked data for AnswerUpdated events raised by the MockV3AggregatorContract contract.
type MockV3AggregatorContractAnswerUpdatedIterator struct {
	Event *MockV3AggregatorContractAnswerUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MockV3AggregatorContractAnswerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockV3AggregatorContractAnswerUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MockV3AggregatorContractAnswerUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MockV3AggregatorContractAnswerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockV3AggregatorContractAnswerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockV3AggregatorContractAnswerUpdated represents a AnswerUpdated event raised by the MockV3AggregatorContract contract.
type MockV3AggregatorContractAnswerUpdated struct {
	Current   *big.Int
	RoundId   *big.Int
	UpdatedAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAnswerUpdated is a free log retrieval operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
func (_MockV3AggregatorContract *MockV3AggregatorContractFilterer) FilterAnswerUpdated(opts *bind.FilterOpts, current []*big.Int, roundId []*big.Int) (*MockV3AggregatorContractAnswerUpdatedIterator, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _MockV3AggregatorContract.contract.FilterLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return &MockV3AggregatorContractAnswerUpdatedIterator{contract: _MockV3AggregatorContract.contract, event: "AnswerUpdated", logs: logs, sub: sub}, nil
}

// WatchAnswerUpdated is a free log subscription operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
func (_MockV3AggregatorContract *MockV3AggregatorContractFilterer) WatchAnswerUpdated(opts *bind.WatchOpts, sink chan<- *MockV3AggregatorContractAnswerUpdated, current []*big.Int, roundId []*big.Int) (event.Subscription, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _MockV3AggregatorContract.contract.WatchLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockV3AggregatorContractAnswerUpdated)
				if err := _MockV3AggregatorContract.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAnswerUpdated is a log parse operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
func (_MockV3AggregatorContract *MockV3AggregatorContractFilterer) ParseAnswerUpdated(log types.Log) (*MockV3AggregatorContractAnswerUpdated, error) {
	event := new(MockV3AggregatorContractAnswerUpdated)
	if err := _MockV3AggregatorContract.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockV3AggregatorContractNewRoundIterator is returned from FilterNewRound and is used to iterate over the raw logs and unpacked data for NewRound events raised by the MockV3AggregatorContract contract.
type MockV3AggregatorContractNewRoundIterator struct {
	Event *MockV3AggregatorContractNewRound // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MockV3AggregatorContractNewRoundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockV3AggregatorContractNewRound)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MockV3AggregatorContractNewRound)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MockV3AggregatorContractNewRoundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockV3AggregatorContractNewRoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockV3AggregatorContractNewRound represents a NewRound event raised by the MockV3AggregatorContract contract.
type MockV3AggregatorContractNewRound struct {
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterNewRound is a free log retrieval operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_MockV3AggregatorContract *MockV3AggregatorContractFilterer) FilterNewRound(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*MockV3AggregatorContractNewRoundIterator, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _MockV3AggregatorContract.contract.FilterLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return &MockV3AggregatorContractNewRoundIterator{contract: _MockV3AggregatorContract.contract, event: "NewRound", logs: logs, sub: sub}, nil
}

// WatchNewRound is a free log subscription operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_MockV3AggregatorContract *MockV3AggregatorContractFilterer) WatchNewRound(opts *bind.WatchOpts, sink chan<- *MockV3AggregatorContractNewRound, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _MockV3AggregatorContract.contract.WatchLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockV3AggregatorContractNewRound)
				if err := _MockV3AggregatorContract.contract.UnpackLog(event, "NewRound", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNewRound is a log parse operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_MockV3AggregatorContract *MockV3AggregatorContractFilterer) ParseNewRound(log types.Log) (*MockV3AggregatorContractNewRound, error) {
	event := new(MockV3AggregatorContractNewRound)
	if err := _MockV3AggregatorContract.contract.UnpackLog(event, "NewRound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
