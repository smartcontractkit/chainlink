// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_load_test_with_metrics

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated"
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

var VRFV2LoadTestWithMetricsMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlyCoordinatorCanFulfill\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"COORDINATOR\",\"outputs\":[{\"internalType\":\"contractVRFCoordinatorV2Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"_subId\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"_requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bytes32\",\"name\":\"_keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint16\",\"name\":\"_requestCount\",\"type\":\"uint16\"}],\"name\":\"requestRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_averageFulfillmentInMillions\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_fastestFulfillment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_responseCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_slowestFulfillment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60e0604052600060045560006005556103e760065534801561002057600080fd5b50604051610cb6380380610cb683398101604081905261003f916101be565b6001600160601b0319606083901b1660805233806000816100a75760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156100d7576100d7816100f8565b5050506001600160601b0319606092831b811660a052911b1660c0526101f1565b6001600160a01b0381163314156101515760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161009e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b03811681146101b957600080fd5b919050565b600080604083850312156101d157600080fd5b6101da836101a2565b91506101e8602084016101a2565b90509250929050565b60805160601c60a05160601c60c05160601c610a7f610237600039600061011a015260008181610179015261028e01526000818161037801526103e00152610a7f6000f3fe608060405234801561001057600080fd5b50600436106100df5760003560e01c8063737144bc1161008c5780638da5cb5b116100665780638da5cb5b146101be578063d826f88f146101dc578063dc1670db146101fb578063f2fde38b1461020457600080fd5b8063737144bc146101a457806374dba124146101ad57806379ba5097146101b657600080fd5b80631fe543e3116100bd5780631fe543e3146101615780633b2bcbf114610174578063557d2e921461019b57600080fd5b8063096cb17b146100e45780631757f11c146100f95780631b6b6d2314610115575b600080fd5b6100f76100f23660046108b6565b610217565b005b61010260055481565b6040519081526020015b60405180910390f35b61013c7f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161010c565b6100f761016f3660046107c7565b610360565b61013c7f000000000000000000000000000000000000000000000000000000000000000081565b61010260035481565b61010260045481565b61010260065481565b6100f7610420565b60005473ffffffffffffffffffffffffffffffffffffffff1661013c565b6100f76000600481905560058190556103e76006556003819055600255565b61010260025481565b6100f7610212366004610771565b61051d565b61021f610531565b60005b8161ffff168161ffff161015610359576040517f5d3b1d300000000000000000000000000000000000000000000000000000000081526004810184905267ffffffffffffffff8616602482015261ffff85166044820152620493e06064820152600160848201526000907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690635d3b1d309060a401602060405180830381600087803b1580156102e757600080fd5b505af11580156102fb573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061031f91906107ae565b600380549192506000610331836109db565b9091555050600090815260076020526040902043905580610351816109b9565b915050610222565b5050505050565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610412576040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001660248201526044015b60405180910390fd5b61041c82826105b4565b5050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146104a1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610409565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610525610531565b61052e81610664565b50565b60005473ffffffffffffffffffffffffffffffffffffffff1633146105b2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610409565b565b6000828152600760205260408120546105cd90436109a2565b905060006105de82620f4240610965565b90506005548211156105f05760058290555b600654821061060157600654610603565b815b6006556002546106135780610646565b600254610621906001610912565b816002546004546106329190610965565b61063c9190610912565b610646919061092a565b60045560028054906000610659836109db565b919050555050505050565b73ffffffffffffffffffffffffffffffffffffffff81163314156106e4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610409565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b803561ffff8116811461076c57600080fd5b919050565b60006020828403121561078357600080fd5b813573ffffffffffffffffffffffffffffffffffffffff811681146107a757600080fd5b9392505050565b6000602082840312156107c057600080fd5b5051919050565b600080604083850312156107da57600080fd5b8235915060208084013567ffffffffffffffff808211156107fa57600080fd5b818601915086601f83011261080e57600080fd5b81358181111561082057610820610a43565b8060051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f8301168101818110858211171561086357610863610a43565b604052828152858101935084860182860187018b101561088257600080fd5b600095505b838610156108a5578035855260019590950194938601938601610887565b508096505050505050509250929050565b600080600080608085870312156108cc57600080fd5b843567ffffffffffffffff811681146108e457600080fd5b93506108f26020860161075a565b9250604085013591506109076060860161075a565b905092959194509250565b6000821982111561092557610925610a14565b500190565b600082610960577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048311821515161561099d5761099d610a14565b500290565b6000828210156109b4576109b4610a14565b500390565b600061ffff808316818114156109d1576109d1610a14565b6001019392505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415610a0d57610a0d610a14565b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var VRFV2LoadTestWithMetricsABI = VRFV2LoadTestWithMetricsMetaData.ABI

var VRFV2LoadTestWithMetricsBin = VRFV2LoadTestWithMetricsMetaData.Bin

func DeployVRFV2LoadTestWithMetrics(auth *bind.TransactOpts, backend bind.ContractBackend, _vrfCoordinator common.Address, _link common.Address) (common.Address, *types.Transaction, *VRFV2LoadTestWithMetrics, error) {
	parsed, err := VRFV2LoadTestWithMetricsMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFV2LoadTestWithMetricsBin), backend, _vrfCoordinator, _link)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFV2LoadTestWithMetrics{VRFV2LoadTestWithMetricsCaller: VRFV2LoadTestWithMetricsCaller{contract: contract}, VRFV2LoadTestWithMetricsTransactor: VRFV2LoadTestWithMetricsTransactor{contract: contract}, VRFV2LoadTestWithMetricsFilterer: VRFV2LoadTestWithMetricsFilterer{contract: contract}}, nil
}

type VRFV2LoadTestWithMetrics struct {
	address common.Address
	abi     abi.ABI
	VRFV2LoadTestWithMetricsCaller
	VRFV2LoadTestWithMetricsTransactor
	VRFV2LoadTestWithMetricsFilterer
}

type VRFV2LoadTestWithMetricsCaller struct {
	contract *bind.BoundContract
}

type VRFV2LoadTestWithMetricsTransactor struct {
	contract *bind.BoundContract
}

type VRFV2LoadTestWithMetricsFilterer struct {
	contract *bind.BoundContract
}

type VRFV2LoadTestWithMetricsSession struct {
	Contract     *VRFV2LoadTestWithMetrics
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFV2LoadTestWithMetricsCallerSession struct {
	Contract *VRFV2LoadTestWithMetricsCaller
	CallOpts bind.CallOpts
}

type VRFV2LoadTestWithMetricsTransactorSession struct {
	Contract     *VRFV2LoadTestWithMetricsTransactor
	TransactOpts bind.TransactOpts
}

type VRFV2LoadTestWithMetricsRaw struct {
	Contract *VRFV2LoadTestWithMetrics
}

type VRFV2LoadTestWithMetricsCallerRaw struct {
	Contract *VRFV2LoadTestWithMetricsCaller
}

type VRFV2LoadTestWithMetricsTransactorRaw struct {
	Contract *VRFV2LoadTestWithMetricsTransactor
}

func NewVRFV2LoadTestWithMetrics(address common.Address, backend bind.ContractBackend) (*VRFV2LoadTestWithMetrics, error) {
	abi, err := abi.JSON(strings.NewReader(VRFV2LoadTestWithMetricsABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFV2LoadTestWithMetrics(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFV2LoadTestWithMetrics{address: address, abi: abi, VRFV2LoadTestWithMetricsCaller: VRFV2LoadTestWithMetricsCaller{contract: contract}, VRFV2LoadTestWithMetricsTransactor: VRFV2LoadTestWithMetricsTransactor{contract: contract}, VRFV2LoadTestWithMetricsFilterer: VRFV2LoadTestWithMetricsFilterer{contract: contract}}, nil
}

func NewVRFV2LoadTestWithMetricsCaller(address common.Address, caller bind.ContractCaller) (*VRFV2LoadTestWithMetricsCaller, error) {
	contract, err := bindVRFV2LoadTestWithMetrics(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2LoadTestWithMetricsCaller{contract: contract}, nil
}

func NewVRFV2LoadTestWithMetricsTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFV2LoadTestWithMetricsTransactor, error) {
	contract, err := bindVRFV2LoadTestWithMetrics(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2LoadTestWithMetricsTransactor{contract: contract}, nil
}

func NewVRFV2LoadTestWithMetricsFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFV2LoadTestWithMetricsFilterer, error) {
	contract, err := bindVRFV2LoadTestWithMetrics(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFV2LoadTestWithMetricsFilterer{contract: contract}, nil
}

func bindVRFV2LoadTestWithMetrics(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFV2LoadTestWithMetricsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2LoadTestWithMetrics.Contract.VRFV2LoadTestWithMetricsCaller.contract.Call(opts, result, method, params...)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.Contract.VRFV2LoadTestWithMetricsTransactor.contract.Transfer(opts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.Contract.VRFV2LoadTestWithMetricsTransactor.contract.Transact(opts, method, params...)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2LoadTestWithMetrics.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.Contract.contract.Transfer(opts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.Contract.contract.Transact(opts, method, params...)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsCaller) COORDINATOR(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2LoadTestWithMetrics.contract.Call(opts, &out, "COORDINATOR")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsSession) COORDINATOR() (common.Address, error) {
	return _VRFV2LoadTestWithMetrics.Contract.COORDINATOR(&_VRFV2LoadTestWithMetrics.CallOpts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsCallerSession) COORDINATOR() (common.Address, error) {
	return _VRFV2LoadTestWithMetrics.Contract.COORDINATOR(&_VRFV2LoadTestWithMetrics.CallOpts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsCaller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2LoadTestWithMetrics.contract.Call(opts, &out, "LINK")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsSession) LINK() (common.Address, error) {
	return _VRFV2LoadTestWithMetrics.Contract.LINK(&_VRFV2LoadTestWithMetrics.CallOpts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsCallerSession) LINK() (common.Address, error) {
	return _VRFV2LoadTestWithMetrics.Contract.LINK(&_VRFV2LoadTestWithMetrics.CallOpts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2LoadTestWithMetrics.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsSession) Owner() (common.Address, error) {
	return _VRFV2LoadTestWithMetrics.Contract.Owner(&_VRFV2LoadTestWithMetrics.CallOpts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsCallerSession) Owner() (common.Address, error) {
	return _VRFV2LoadTestWithMetrics.Contract.Owner(&_VRFV2LoadTestWithMetrics.CallOpts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsCaller) SAverageFulfillmentInMillions(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2LoadTestWithMetrics.contract.Call(opts, &out, "s_averageFulfillmentInMillions")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsSession) SAverageFulfillmentInMillions() (*big.Int, error) {
	return _VRFV2LoadTestWithMetrics.Contract.SAverageFulfillmentInMillions(&_VRFV2LoadTestWithMetrics.CallOpts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsCallerSession) SAverageFulfillmentInMillions() (*big.Int, error) {
	return _VRFV2LoadTestWithMetrics.Contract.SAverageFulfillmentInMillions(&_VRFV2LoadTestWithMetrics.CallOpts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsCaller) SFastestFulfillment(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2LoadTestWithMetrics.contract.Call(opts, &out, "s_fastestFulfillment")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsSession) SFastestFulfillment() (*big.Int, error) {
	return _VRFV2LoadTestWithMetrics.Contract.SFastestFulfillment(&_VRFV2LoadTestWithMetrics.CallOpts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsCallerSession) SFastestFulfillment() (*big.Int, error) {
	return _VRFV2LoadTestWithMetrics.Contract.SFastestFulfillment(&_VRFV2LoadTestWithMetrics.CallOpts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsCaller) SRequestCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2LoadTestWithMetrics.contract.Call(opts, &out, "s_requestCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsSession) SRequestCount() (*big.Int, error) {
	return _VRFV2LoadTestWithMetrics.Contract.SRequestCount(&_VRFV2LoadTestWithMetrics.CallOpts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsCallerSession) SRequestCount() (*big.Int, error) {
	return _VRFV2LoadTestWithMetrics.Contract.SRequestCount(&_VRFV2LoadTestWithMetrics.CallOpts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsCaller) SResponseCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2LoadTestWithMetrics.contract.Call(opts, &out, "s_responseCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsSession) SResponseCount() (*big.Int, error) {
	return _VRFV2LoadTestWithMetrics.Contract.SResponseCount(&_VRFV2LoadTestWithMetrics.CallOpts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsCallerSession) SResponseCount() (*big.Int, error) {
	return _VRFV2LoadTestWithMetrics.Contract.SResponseCount(&_VRFV2LoadTestWithMetrics.CallOpts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsCaller) SSlowestFulfillment(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2LoadTestWithMetrics.contract.Call(opts, &out, "s_slowestFulfillment")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsSession) SSlowestFulfillment() (*big.Int, error) {
	return _VRFV2LoadTestWithMetrics.Contract.SSlowestFulfillment(&_VRFV2LoadTestWithMetrics.CallOpts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsCallerSession) SSlowestFulfillment() (*big.Int, error) {
	return _VRFV2LoadTestWithMetrics.Contract.SSlowestFulfillment(&_VRFV2LoadTestWithMetrics.CallOpts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.contract.Transact(opts, "acceptOwnership")
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.Contract.AcceptOwnership(&_VRFV2LoadTestWithMetrics.TransactOpts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.Contract.AcceptOwnership(&_VRFV2LoadTestWithMetrics.TransactOpts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.contract.Transact(opts, "rawFulfillRandomWords", requestId, randomWords)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.Contract.RawFulfillRandomWords(&_VRFV2LoadTestWithMetrics.TransactOpts, requestId, randomWords)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsTransactorSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.Contract.RawFulfillRandomWords(&_VRFV2LoadTestWithMetrics.TransactOpts, requestId, randomWords)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsTransactor) RequestRandomWords(opts *bind.TransactOpts, _subId uint64, _requestConfirmations uint16, _keyHash [32]byte, _requestCount uint16) (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.contract.Transact(opts, "requestRandomWords", _subId, _requestConfirmations, _keyHash, _requestCount)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsSession) RequestRandomWords(_subId uint64, _requestConfirmations uint16, _keyHash [32]byte, _requestCount uint16) (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.Contract.RequestRandomWords(&_VRFV2LoadTestWithMetrics.TransactOpts, _subId, _requestConfirmations, _keyHash, _requestCount)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsTransactorSession) RequestRandomWords(_subId uint64, _requestConfirmations uint16, _keyHash [32]byte, _requestCount uint16) (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.Contract.RequestRandomWords(&_VRFV2LoadTestWithMetrics.TransactOpts, _subId, _requestConfirmations, _keyHash, _requestCount)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsTransactor) Reset(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.contract.Transact(opts, "reset")
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsSession) Reset() (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.Contract.Reset(&_VRFV2LoadTestWithMetrics.TransactOpts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsTransactorSession) Reset() (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.Contract.Reset(&_VRFV2LoadTestWithMetrics.TransactOpts)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.contract.Transact(opts, "transferOwnership", to)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.Contract.TransferOwnership(&_VRFV2LoadTestWithMetrics.TransactOpts, to)
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFV2LoadTestWithMetrics.Contract.TransferOwnership(&_VRFV2LoadTestWithMetrics.TransactOpts, to)
}

type VRFV2LoadTestWithMetricsOwnershipTransferRequestedIterator struct {
	Event *VRFV2LoadTestWithMetricsOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2LoadTestWithMetricsOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2LoadTestWithMetricsOwnershipTransferRequested)
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

	select {
	case log := <-it.logs:
		it.Event = new(VRFV2LoadTestWithMetricsOwnershipTransferRequested)
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

func (it *VRFV2LoadTestWithMetricsOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFV2LoadTestWithMetricsOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2LoadTestWithMetricsOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2LoadTestWithMetricsOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2LoadTestWithMetrics.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2LoadTestWithMetricsOwnershipTransferRequestedIterator{contract: _VRFV2LoadTestWithMetrics.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFV2LoadTestWithMetricsOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2LoadTestWithMetrics.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2LoadTestWithMetricsOwnershipTransferRequested)
				if err := _VRFV2LoadTestWithMetrics.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsFilterer) ParseOwnershipTransferRequested(log types.Log) (*VRFV2LoadTestWithMetricsOwnershipTransferRequested, error) {
	event := new(VRFV2LoadTestWithMetricsOwnershipTransferRequested)
	if err := _VRFV2LoadTestWithMetrics.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFV2LoadTestWithMetricsOwnershipTransferredIterator struct {
	Event *VRFV2LoadTestWithMetricsOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2LoadTestWithMetricsOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2LoadTestWithMetricsOwnershipTransferred)
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

	select {
	case log := <-it.logs:
		it.Event = new(VRFV2LoadTestWithMetricsOwnershipTransferred)
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

func (it *VRFV2LoadTestWithMetricsOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFV2LoadTestWithMetricsOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2LoadTestWithMetricsOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2LoadTestWithMetricsOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2LoadTestWithMetrics.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2LoadTestWithMetricsOwnershipTransferredIterator{contract: _VRFV2LoadTestWithMetrics.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFV2LoadTestWithMetricsOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2LoadTestWithMetrics.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2LoadTestWithMetricsOwnershipTransferred)
				if err := _VRFV2LoadTestWithMetrics.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetricsFilterer) ParseOwnershipTransferred(log types.Log) (*VRFV2LoadTestWithMetricsOwnershipTransferred, error) {
	event := new(VRFV2LoadTestWithMetricsOwnershipTransferred)
	if err := _VRFV2LoadTestWithMetrics.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetrics) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFV2LoadTestWithMetrics.abi.Events["OwnershipTransferRequested"].ID:
		return _VRFV2LoadTestWithMetrics.ParseOwnershipTransferRequested(log)
	case _VRFV2LoadTestWithMetrics.abi.Events["OwnershipTransferred"].ID:
		return _VRFV2LoadTestWithMetrics.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFV2LoadTestWithMetricsOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VRFV2LoadTestWithMetricsOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_VRFV2LoadTestWithMetrics *VRFV2LoadTestWithMetrics) Address() common.Address {
	return _VRFV2LoadTestWithMetrics.address
}

type VRFV2LoadTestWithMetricsInterface interface {
	COORDINATOR(opts *bind.CallOpts) (common.Address, error)

	LINK(opts *bind.CallOpts) (common.Address, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SAverageFulfillmentInMillions(opts *bind.CallOpts) (*big.Int, error)

	SFastestFulfillment(opts *bind.CallOpts) (*big.Int, error)

	SRequestCount(opts *bind.CallOpts) (*big.Int, error)

	SResponseCount(opts *bind.CallOpts) (*big.Int, error)

	SSlowestFulfillment(opts *bind.CallOpts) (*big.Int, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	RequestRandomWords(opts *bind.TransactOpts, _subId uint64, _requestConfirmations uint16, _keyHash [32]byte, _requestCount uint16) (*types.Transaction, error)

	Reset(opts *bind.TransactOpts) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2LoadTestWithMetricsOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFV2LoadTestWithMetricsOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VRFV2LoadTestWithMetricsOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2LoadTestWithMetricsOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFV2LoadTestWithMetricsOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VRFV2LoadTestWithMetricsOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
