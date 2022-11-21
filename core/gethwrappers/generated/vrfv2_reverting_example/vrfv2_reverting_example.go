// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrfv2_reverting_example

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

var VRFV2RevertingExampleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlyCoordinatorCanFulfill\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_gasAvailable\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_randomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_subId\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"testCreateSubscriptionAndFund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"minReqConfs\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"}],\"name\":\"testRequestRandomness\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"topUpSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"name\":\"updateSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561001057600080fd5b50604051610da6380380610da683398101604081905261002f9161008e565b6001600160601b0319606083901b16608052600280546001600160a01b03199081166001600160a01b0394851617909155600380549290931691161790556100c1565b80516001600160a01b038116811461008957600080fd5b919050565b600080604083850312156100a157600080fd5b6100aa83610072565b91506100b860208401610072565b90509250929050565b60805160601c610cc06100e66000396000818161019e01526102060152610cc06000f3fe608060405234801561001057600080fd5b50600436106100a35760003560e01c80636802f72611610076578063e89e106a1161005b578063e89e106a14610161578063f08c5daa1461016a578063f6eaffc81461017357600080fd5b80636802f72614610109578063706da1ca1461011c57600080fd5b80631fe543e3146100a857806327784fad146100bd5780632fa4e442146100e357806336bfffed146100f6575b600080fd5b6100bb6100b63660046109d4565b610186565b005b6100d06100cb366004610939565b610246565b6040519081526020015b60405180910390f35b6100bb6100f1366004610a95565b610328565b6100bb610104366004610851565b610488565b6100bb610117366004610a95565b610610565b6003546101489074010000000000000000000000000000000000000000900467ffffffffffffffff1681565b60405167ffffffffffffffff90911681526020016100da565b6100d060015481565b6100d060045481565b6100d06101813660046109a2565b610817565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610238576040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001660248201526044015b60405180910390fd5b6102428282600080fd5b5050565b6002546040517f5d3b1d300000000000000000000000000000000000000000000000000000000081526004810187905267ffffffffffffffff8616602482015261ffff8516604482015263ffffffff80851660648301528316608482015260009173ffffffffffffffffffffffffffffffffffffffff1690635d3b1d309060a401602060405180830381600087803b1580156102e157600080fd5b505af11580156102f5573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061031991906109bb565b60018190559695505050505050565b60035474010000000000000000000000000000000000000000900467ffffffffffffffff166103b3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f737562206e6f7420736574000000000000000000000000000000000000000000604482015260640161022f565b6003546002546040805174010000000000000000000000000000000000000000840467ffffffffffffffff16602082015273ffffffffffffffffffffffffffffffffffffffff93841693634000aea09316918591015b6040516020818303038152906040526040518463ffffffff1660e01b815260040161043693929190610ac3565b602060405180830381600087803b15801561045057600080fd5b505af1158015610464573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102429190610910565b60035474010000000000000000000000000000000000000000900467ffffffffffffffff16610513576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f7375624944206e6f742073657400000000000000000000000000000000000000604482015260640161022f565b60005b815181101561024257600254600354835173ffffffffffffffffffffffffffffffffffffffff90921691637341c10c9174010000000000000000000000000000000000000000900467ffffffffffffffff169085908590811061057b5761057b610c3c565b60200260200101516040518363ffffffff1660e01b81526004016105cb92919067ffffffffffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b600060405180830381600087803b1580156105e557600080fd5b505af11580156105f9573d6000803e3d6000fd5b50505050808061060890610bdc565b915050610516565b60035474010000000000000000000000000000000000000000900467ffffffffffffffff166103b357600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a21a23e46040518163ffffffff1660e01b8152600401602060405180830381600087803b1580156106a357600080fd5b505af11580156106b7573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106db9190610a78565b600380547fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000067ffffffffffffffff938416810291909117918290556002546040517f7341c10c00000000000000000000000000000000000000000000000000000000815291909204909216600483015230602483015273ffffffffffffffffffffffffffffffffffffffff1690637341c10c90604401600060405180830381600087803b1580156107a257600080fd5b505af11580156107b6573d6000803e3d6000fd5b50506003546002546040805174010000000000000000000000000000000000000000840467ffffffffffffffff16602082015273ffffffffffffffffffffffffffffffffffffffff9384169550634000aea094509290911691859101610409565b6000818154811061082757600080fd5b600091825260209091200154905081565b803563ffffffff8116811461084c57600080fd5b919050565b6000602080838503121561086457600080fd5b823567ffffffffffffffff81111561087b57600080fd5b8301601f8101851361088c57600080fd5b803561089f61089a82610bb8565b610b69565b80828252848201915084840188868560051b87010111156108bf57600080fd5b60009450845b8481101561090257813573ffffffffffffffffffffffffffffffffffffffff811681146108f0578687fd5b845292860192908601906001016108c5565b509098975050505050505050565b60006020828403121561092257600080fd5b8151801515811461093257600080fd5b9392505050565b600080600080600060a0868803121561095157600080fd5b85359450602086013561096381610c9a565b9350604086013561ffff8116811461097a57600080fd5b925061098860608701610838565b915061099660808701610838565b90509295509295909350565b6000602082840312156109b457600080fd5b5035919050565b6000602082840312156109cd57600080fd5b5051919050565b600080604083850312156109e757600080fd5b8235915060208084013567ffffffffffffffff811115610a0657600080fd5b8401601f81018613610a1757600080fd5b8035610a2561089a82610bb8565b80828252848201915084840189868560051b8701011115610a4557600080fd5b600094505b83851015610a68578035835260019490940193918501918501610a4a565b5080955050505050509250929050565b600060208284031215610a8a57600080fd5b815161093281610c9a565b600060208284031215610aa757600080fd5b81356bffffffffffffffffffffffff8116811461093257600080fd5b73ffffffffffffffffffffffffffffffffffffffff84168152600060206bffffffffffffffffffffffff85168184015260606040840152835180606085015260005b81811015610b2157858101830151858201608001528201610b05565b81811115610b33576000608083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160800195945050505050565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610bb057610bb0610c6b565b604052919050565b600067ffffffffffffffff821115610bd257610bd2610c6b565b5060051b60200190565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415610c35577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b67ffffffffffffffff81168114610cb057600080fd5b5056fea164736f6c6343000806000a",
}

var VRFV2RevertingExampleABI = VRFV2RevertingExampleMetaData.ABI

var VRFV2RevertingExampleBin = VRFV2RevertingExampleMetaData.Bin

func DeployVRFV2RevertingExample(auth *bind.TransactOpts, backend bind.ContractBackend, vrfCoordinator common.Address, link common.Address) (common.Address, *types.Transaction, *VRFV2RevertingExample, error) {
	parsed, err := VRFV2RevertingExampleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFV2RevertingExampleBin), backend, vrfCoordinator, link)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFV2RevertingExample{VRFV2RevertingExampleCaller: VRFV2RevertingExampleCaller{contract: contract}, VRFV2RevertingExampleTransactor: VRFV2RevertingExampleTransactor{contract: contract}, VRFV2RevertingExampleFilterer: VRFV2RevertingExampleFilterer{contract: contract}}, nil
}

type VRFV2RevertingExample struct {
	address common.Address
	abi     abi.ABI
	VRFV2RevertingExampleCaller
	VRFV2RevertingExampleTransactor
	VRFV2RevertingExampleFilterer
}

type VRFV2RevertingExampleCaller struct {
	contract *bind.BoundContract
}

type VRFV2RevertingExampleTransactor struct {
	contract *bind.BoundContract
}

type VRFV2RevertingExampleFilterer struct {
	contract *bind.BoundContract
}

type VRFV2RevertingExampleSession struct {
	Contract     *VRFV2RevertingExample
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFV2RevertingExampleCallerSession struct {
	Contract *VRFV2RevertingExampleCaller
	CallOpts bind.CallOpts
}

type VRFV2RevertingExampleTransactorSession struct {
	Contract     *VRFV2RevertingExampleTransactor
	TransactOpts bind.TransactOpts
}

type VRFV2RevertingExampleRaw struct {
	Contract *VRFV2RevertingExample
}

type VRFV2RevertingExampleCallerRaw struct {
	Contract *VRFV2RevertingExampleCaller
}

type VRFV2RevertingExampleTransactorRaw struct {
	Contract *VRFV2RevertingExampleTransactor
}

func NewVRFV2RevertingExample(address common.Address, backend bind.ContractBackend) (*VRFV2RevertingExample, error) {
	abi, err := abi.JSON(strings.NewReader(VRFV2RevertingExampleABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFV2RevertingExample(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFV2RevertingExample{address: address, abi: abi, VRFV2RevertingExampleCaller: VRFV2RevertingExampleCaller{contract: contract}, VRFV2RevertingExampleTransactor: VRFV2RevertingExampleTransactor{contract: contract}, VRFV2RevertingExampleFilterer: VRFV2RevertingExampleFilterer{contract: contract}}, nil
}

func NewVRFV2RevertingExampleCaller(address common.Address, caller bind.ContractCaller) (*VRFV2RevertingExampleCaller, error) {
	contract, err := bindVRFV2RevertingExample(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2RevertingExampleCaller{contract: contract}, nil
}

func NewVRFV2RevertingExampleTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFV2RevertingExampleTransactor, error) {
	contract, err := bindVRFV2RevertingExample(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2RevertingExampleTransactor{contract: contract}, nil
}

func NewVRFV2RevertingExampleFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFV2RevertingExampleFilterer, error) {
	contract, err := bindVRFV2RevertingExample(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFV2RevertingExampleFilterer{contract: contract}, nil
}

func bindVRFV2RevertingExample(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFV2RevertingExampleABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_VRFV2RevertingExample *VRFV2RevertingExampleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2RevertingExample.Contract.VRFV2RevertingExampleCaller.contract.Call(opts, result, method, params...)
}

func (_VRFV2RevertingExample *VRFV2RevertingExampleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2RevertingExample.Contract.VRFV2RevertingExampleTransactor.contract.Transfer(opts)
}

func (_VRFV2RevertingExample *VRFV2RevertingExampleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2RevertingExample.Contract.VRFV2RevertingExampleTransactor.contract.Transact(opts, method, params...)
}

func (_VRFV2RevertingExample *VRFV2RevertingExampleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2RevertingExample.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFV2RevertingExample *VRFV2RevertingExampleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2RevertingExample.Contract.contract.Transfer(opts)
}

func (_VRFV2RevertingExample *VRFV2RevertingExampleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2RevertingExample.Contract.contract.Transact(opts, method, params...)
}

func (_VRFV2RevertingExample *VRFV2RevertingExampleCaller) SGasAvailable(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2RevertingExample.contract.Call(opts, &out, "s_gasAvailable")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2RevertingExample *VRFV2RevertingExampleSession) SGasAvailable() (*big.Int, error) {
	return _VRFV2RevertingExample.Contract.SGasAvailable(&_VRFV2RevertingExample.CallOpts)
}

func (_VRFV2RevertingExample *VRFV2RevertingExampleCallerSession) SGasAvailable() (*big.Int, error) {
	return _VRFV2RevertingExample.Contract.SGasAvailable(&_VRFV2RevertingExample.CallOpts)
}

func (_VRFV2RevertingExample *VRFV2RevertingExampleCaller) SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2RevertingExample.contract.Call(opts, &out, "s_randomWords", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2RevertingExample *VRFV2RevertingExampleSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFV2RevertingExample.Contract.SRandomWords(&_VRFV2RevertingExample.CallOpts, arg0)
}

func (_VRFV2RevertingExample *VRFV2RevertingExampleCallerSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFV2RevertingExample.Contract.SRandomWords(&_VRFV2RevertingExample.CallOpts, arg0)
}

func (_VRFV2RevertingExample *VRFV2RevertingExampleCaller) SRequestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2RevertingExample.contract.Call(opts, &out, "s_requestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2RevertingExample *VRFV2RevertingExampleSession) SRequestId() (*big.Int, error) {
	return _VRFV2RevertingExample.Contract.SRequestId(&_VRFV2RevertingExample.CallOpts)
}

func (_VRFV2RevertingExample *VRFV2RevertingExampleCallerSession) SRequestId() (*big.Int, error) {
	return _VRFV2RevertingExample.Contract.SRequestId(&_VRFV2RevertingExample.CallOpts)
}

func (_VRFV2RevertingExample *VRFV2RevertingExampleCaller) SSubId(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _VRFV2RevertingExample.contract.Call(opts, &out, "s_subId")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_VRFV2RevertingExample *VRFV2RevertingExampleSession) SSubId() (uint64, error) {
	return _VRFV2RevertingExample.Contract.SSubId(&_VRFV2RevertingExample.CallOpts)
}

func (_VRFV2RevertingExample *VRFV2RevertingExampleCallerSession) SSubId() (uint64, error) {
	return _VRFV2RevertingExample.Contract.SSubId(&_VRFV2RevertingExample.CallOpts)
}

func (_VRFV2RevertingExample *VRFV2RevertingExampleTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2RevertingExample.contract.Transact(opts, "rawFulfillRandomWords", requestId, randomWords)
}

func (_VRFV2RevertingExample *VRFV2RevertingExampleSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2RevertingExample.Contract.RawFulfillRandomWords(&_VRFV2RevertingExample.TransactOpts, requestId, randomWords)
}

func (_VRFV2RevertingExample *VRFV2RevertingExampleTransactorSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2RevertingExample.Contract.RawFulfillRandomWords(&_VRFV2RevertingExample.TransactOpts, requestId, randomWords)
}

func (_VRFV2RevertingExample *VRFV2RevertingExampleTransactor) TestCreateSubscriptionAndFund(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VRFV2RevertingExample.contract.Transact(opts, "testCreateSubscriptionAndFund", amount)
}

func (_VRFV2RevertingExample *VRFV2RevertingExampleSession) TestCreateSubscriptionAndFund(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2RevertingExample.Contract.TestCreateSubscriptionAndFund(&_VRFV2RevertingExample.TransactOpts, amount)
}

func (_VRFV2RevertingExample *VRFV2RevertingExampleTransactorSession) TestCreateSubscriptionAndFund(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2RevertingExample.Contract.TestCreateSubscriptionAndFund(&_VRFV2RevertingExample.TransactOpts, amount)
}

func (_VRFV2RevertingExample *VRFV2RevertingExampleTransactor) TestRequestRandomness(opts *bind.TransactOpts, keyHash [32]byte, subId uint64, minReqConfs uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _VRFV2RevertingExample.contract.Transact(opts, "testRequestRandomness", keyHash, subId, minReqConfs, callbackGasLimit, numWords)
}

func (_VRFV2RevertingExample *VRFV2RevertingExampleSession) TestRequestRandomness(keyHash [32]byte, subId uint64, minReqConfs uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _VRFV2RevertingExample.Contract.TestRequestRandomness(&_VRFV2RevertingExample.TransactOpts, keyHash, subId, minReqConfs, callbackGasLimit, numWords)
}

func (_VRFV2RevertingExample *VRFV2RevertingExampleTransactorSession) TestRequestRandomness(keyHash [32]byte, subId uint64, minReqConfs uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _VRFV2RevertingExample.Contract.TestRequestRandomness(&_VRFV2RevertingExample.TransactOpts, keyHash, subId, minReqConfs, callbackGasLimit, numWords)
}

func (_VRFV2RevertingExample *VRFV2RevertingExampleTransactor) TopUpSubscription(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VRFV2RevertingExample.contract.Transact(opts, "topUpSubscription", amount)
}

func (_VRFV2RevertingExample *VRFV2RevertingExampleSession) TopUpSubscription(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2RevertingExample.Contract.TopUpSubscription(&_VRFV2RevertingExample.TransactOpts, amount)
}

func (_VRFV2RevertingExample *VRFV2RevertingExampleTransactorSession) TopUpSubscription(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2RevertingExample.Contract.TopUpSubscription(&_VRFV2RevertingExample.TransactOpts, amount)
}

func (_VRFV2RevertingExample *VRFV2RevertingExampleTransactor) UpdateSubscription(opts *bind.TransactOpts, consumers []common.Address) (*types.Transaction, error) {
	return _VRFV2RevertingExample.contract.Transact(opts, "updateSubscription", consumers)
}

func (_VRFV2RevertingExample *VRFV2RevertingExampleSession) UpdateSubscription(consumers []common.Address) (*types.Transaction, error) {
	return _VRFV2RevertingExample.Contract.UpdateSubscription(&_VRFV2RevertingExample.TransactOpts, consumers)
}

func (_VRFV2RevertingExample *VRFV2RevertingExampleTransactorSession) UpdateSubscription(consumers []common.Address) (*types.Transaction, error) {
	return _VRFV2RevertingExample.Contract.UpdateSubscription(&_VRFV2RevertingExample.TransactOpts, consumers)
}

func (_VRFV2RevertingExample *VRFV2RevertingExample) Address() common.Address {
	return _VRFV2RevertingExample.address
}

type VRFV2RevertingExampleInterface interface {
	SGasAvailable(opts *bind.CallOpts) (*big.Int, error)

	SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	SRequestId(opts *bind.CallOpts) (*big.Int, error)

	SSubId(opts *bind.CallOpts) (uint64, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	TestCreateSubscriptionAndFund(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	TestRequestRandomness(opts *bind.TransactOpts, keyHash [32]byte, subId uint64, minReqConfs uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error)

	TopUpSubscription(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	UpdateSubscription(opts *bind.TransactOpts, consumers []common.Address) (*types.Transaction, error)

	Address() common.Address
}
