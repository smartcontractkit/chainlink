// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrfv2plus_reverting_example

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

var VRFV2PlusRevertingExampleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlyCoordinatorCanFulfill\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"createSubscriptionAndFund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"minReqConfs\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"}],\"name\":\"requestRandomness\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_gasAvailable\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_randomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_subId\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"topUpSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"name\":\"updateSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561001057600080fd5b50604051610dad380380610dad83398101604081905261002f9161008e565b6001600160601b0319606083901b16608052600280546001600160a01b03199081166001600160a01b0394851617909155600380549290931691161790556100c1565b80516001600160a01b038116811461008957600080fd5b919050565b600080604083850312156100a157600080fd5b6100aa83610072565b91506100b860208401610072565b90509250929050565b60805160601c610cc76100e66000396000818161028701526102ef0152610cc76000f3fe608060405234801561001057600080fd5b50600436106100a35760003560e01c8063706da1ca11610076578063e89e106a1161005b578063e89e106a14610161578063f08c5daa1461016a578063f6eaffc81461017357600080fd5b8063706da1ca14610109578063cf62c8ab1461014e57600080fd5b8063177b9692146100a85780631fe543e3146100ce5780632fa4e442146100e357806336bfffed146100f6575b600080fd5b6100bb6100b6366004610940565b610186565b6040519081526020015b60405180910390f35b6100e16100dc3660046109db565b61026f565b005b6100e16100f1366004610a9c565b61032f565b6100e1610104366004610858565b61048f565b6003546101359074010000000000000000000000000000000000000000900467ffffffffffffffff1681565b60405167ffffffffffffffff90911681526020016100c5565b6100e161015c366004610a9c565b610617565b6100bb60015481565b6100bb60045481565b6100bb6101813660046109a9565b61081e565b6002546040517fefcf1d940000000000000000000000000000000000000000000000000000000081526004810187905267ffffffffffffffff8616602482015261ffff8516604482015263ffffffff808516606483015283166084820152600060a482018190529173ffffffffffffffffffffffffffffffffffffffff169063efcf1d949060c401602060405180830381600087803b15801561022857600080fd5b505af115801561023c573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061026091906109c2565b60018190559695505050505050565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610321576040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001660248201526044015b60405180910390fd5b61032b8282600080fd5b5050565b60035474010000000000000000000000000000000000000000900467ffffffffffffffff166103ba576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f737562206e6f74207365740000000000000000000000000000000000000000006044820152606401610318565b6003546002546040805174010000000000000000000000000000000000000000840467ffffffffffffffff16602082015273ffffffffffffffffffffffffffffffffffffffff93841693634000aea09316918591015b6040516020818303038152906040526040518463ffffffff1660e01b815260040161043d93929190610aca565b602060405180830381600087803b15801561045757600080fd5b505af115801561046b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061032b9190610917565b60035474010000000000000000000000000000000000000000900467ffffffffffffffff1661051a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f7375624944206e6f7420736574000000000000000000000000000000000000006044820152606401610318565b60005b815181101561032b57600254600354835173ffffffffffffffffffffffffffffffffffffffff90921691637341c10c9174010000000000000000000000000000000000000000900467ffffffffffffffff169085908590811061058257610582610c43565b60200260200101516040518363ffffffff1660e01b81526004016105d292919067ffffffffffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b600060405180830381600087803b1580156105ec57600080fd5b505af1158015610600573d6000803e3d6000fd5b50505050808061060f90610be3565b91505061051d565b60035474010000000000000000000000000000000000000000900467ffffffffffffffff166103ba57600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a21a23e46040518163ffffffff1660e01b8152600401602060405180830381600087803b1580156106aa57600080fd5b505af11580156106be573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106e29190610a7f565b600380547fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000067ffffffffffffffff938416810291909117918290556002546040517f7341c10c00000000000000000000000000000000000000000000000000000000815291909204909216600483015230602483015273ffffffffffffffffffffffffffffffffffffffff1690637341c10c90604401600060405180830381600087803b1580156107a957600080fd5b505af11580156107bd573d6000803e3d6000fd5b50506003546002546040805174010000000000000000000000000000000000000000840467ffffffffffffffff16602082015273ffffffffffffffffffffffffffffffffffffffff9384169550634000aea094509290911691859101610410565b6000818154811061082e57600080fd5b600091825260209091200154905081565b803563ffffffff8116811461085357600080fd5b919050565b6000602080838503121561086b57600080fd5b823567ffffffffffffffff81111561088257600080fd5b8301601f8101851361089357600080fd5b80356108a66108a182610bbf565b610b70565b80828252848201915084840188868560051b87010111156108c657600080fd5b60009450845b8481101561090957813573ffffffffffffffffffffffffffffffffffffffff811681146108f7578687fd5b845292860192908601906001016108cc565b509098975050505050505050565b60006020828403121561092957600080fd5b8151801515811461093957600080fd5b9392505050565b600080600080600060a0868803121561095857600080fd5b85359450602086013561096a81610ca1565b9350604086013561ffff8116811461098157600080fd5b925061098f6060870161083f565b915061099d6080870161083f565b90509295509295909350565b6000602082840312156109bb57600080fd5b5035919050565b6000602082840312156109d457600080fd5b5051919050565b600080604083850312156109ee57600080fd5b8235915060208084013567ffffffffffffffff811115610a0d57600080fd5b8401601f81018613610a1e57600080fd5b8035610a2c6108a182610bbf565b80828252848201915084840189868560051b8701011115610a4c57600080fd5b600094505b83851015610a6f578035835260019490940193918501918501610a51565b5080955050505050509250929050565b600060208284031215610a9157600080fd5b815161093981610ca1565b600060208284031215610aae57600080fd5b81356bffffffffffffffffffffffff8116811461093957600080fd5b73ffffffffffffffffffffffffffffffffffffffff84168152600060206bffffffffffffffffffffffff85168184015260606040840152835180606085015260005b81811015610b2857858101830151858201608001528201610b0c565b81811115610b3a576000608083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160800195945050505050565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610bb757610bb7610c72565b604052919050565b600067ffffffffffffffff821115610bd957610bd9610c72565b5060051b60200190565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415610c3c577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b67ffffffffffffffff81168114610cb757600080fd5b5056fea164736f6c6343000806000a",
}

var VRFV2PlusRevertingExampleABI = VRFV2PlusRevertingExampleMetaData.ABI

var VRFV2PlusRevertingExampleBin = VRFV2PlusRevertingExampleMetaData.Bin

func DeployVRFV2PlusRevertingExample(auth *bind.TransactOpts, backend bind.ContractBackend, vrfCoordinator common.Address, link common.Address) (common.Address, *types.Transaction, *VRFV2PlusRevertingExample, error) {
	parsed, err := VRFV2PlusRevertingExampleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFV2PlusRevertingExampleBin), backend, vrfCoordinator, link)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFV2PlusRevertingExample{VRFV2PlusRevertingExampleCaller: VRFV2PlusRevertingExampleCaller{contract: contract}, VRFV2PlusRevertingExampleTransactor: VRFV2PlusRevertingExampleTransactor{contract: contract}, VRFV2PlusRevertingExampleFilterer: VRFV2PlusRevertingExampleFilterer{contract: contract}}, nil
}

type VRFV2PlusRevertingExample struct {
	address common.Address
	abi     abi.ABI
	VRFV2PlusRevertingExampleCaller
	VRFV2PlusRevertingExampleTransactor
	VRFV2PlusRevertingExampleFilterer
}

type VRFV2PlusRevertingExampleCaller struct {
	contract *bind.BoundContract
}

type VRFV2PlusRevertingExampleTransactor struct {
	contract *bind.BoundContract
}

type VRFV2PlusRevertingExampleFilterer struct {
	contract *bind.BoundContract
}

type VRFV2PlusRevertingExampleSession struct {
	Contract     *VRFV2PlusRevertingExample
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFV2PlusRevertingExampleCallerSession struct {
	Contract *VRFV2PlusRevertingExampleCaller
	CallOpts bind.CallOpts
}

type VRFV2PlusRevertingExampleTransactorSession struct {
	Contract     *VRFV2PlusRevertingExampleTransactor
	TransactOpts bind.TransactOpts
}

type VRFV2PlusRevertingExampleRaw struct {
	Contract *VRFV2PlusRevertingExample
}

type VRFV2PlusRevertingExampleCallerRaw struct {
	Contract *VRFV2PlusRevertingExampleCaller
}

type VRFV2PlusRevertingExampleTransactorRaw struct {
	Contract *VRFV2PlusRevertingExampleTransactor
}

func NewVRFV2PlusRevertingExample(address common.Address, backend bind.ContractBackend) (*VRFV2PlusRevertingExample, error) {
	abi, err := abi.JSON(strings.NewReader(VRFV2PlusRevertingExampleABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFV2PlusRevertingExample(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusRevertingExample{address: address, abi: abi, VRFV2PlusRevertingExampleCaller: VRFV2PlusRevertingExampleCaller{contract: contract}, VRFV2PlusRevertingExampleTransactor: VRFV2PlusRevertingExampleTransactor{contract: contract}, VRFV2PlusRevertingExampleFilterer: VRFV2PlusRevertingExampleFilterer{contract: contract}}, nil
}

func NewVRFV2PlusRevertingExampleCaller(address common.Address, caller bind.ContractCaller) (*VRFV2PlusRevertingExampleCaller, error) {
	contract, err := bindVRFV2PlusRevertingExample(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusRevertingExampleCaller{contract: contract}, nil
}

func NewVRFV2PlusRevertingExampleTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFV2PlusRevertingExampleTransactor, error) {
	contract, err := bindVRFV2PlusRevertingExample(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusRevertingExampleTransactor{contract: contract}, nil
}

func NewVRFV2PlusRevertingExampleFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFV2PlusRevertingExampleFilterer, error) {
	contract, err := bindVRFV2PlusRevertingExample(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusRevertingExampleFilterer{contract: contract}, nil
}

func bindVRFV2PlusRevertingExample(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFV2PlusRevertingExampleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2PlusRevertingExample.Contract.VRFV2PlusRevertingExampleCaller.contract.Call(opts, result, method, params...)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.VRFV2PlusRevertingExampleTransactor.contract.Transfer(opts)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.VRFV2PlusRevertingExampleTransactor.contract.Transact(opts, method, params...)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2PlusRevertingExample.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.contract.Transfer(opts)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.contract.Transact(opts, method, params...)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleCaller) SGasAvailable(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusRevertingExample.contract.Call(opts, &out, "s_gasAvailable")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) SGasAvailable() (*big.Int, error) {
	return _VRFV2PlusRevertingExample.Contract.SGasAvailable(&_VRFV2PlusRevertingExample.CallOpts)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleCallerSession) SGasAvailable() (*big.Int, error) {
	return _VRFV2PlusRevertingExample.Contract.SGasAvailable(&_VRFV2PlusRevertingExample.CallOpts)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleCaller) SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusRevertingExample.contract.Call(opts, &out, "s_randomWords", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFV2PlusRevertingExample.Contract.SRandomWords(&_VRFV2PlusRevertingExample.CallOpts, arg0)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleCallerSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFV2PlusRevertingExample.Contract.SRandomWords(&_VRFV2PlusRevertingExample.CallOpts, arg0)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleCaller) SRequestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusRevertingExample.contract.Call(opts, &out, "s_requestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) SRequestId() (*big.Int, error) {
	return _VRFV2PlusRevertingExample.Contract.SRequestId(&_VRFV2PlusRevertingExample.CallOpts)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleCallerSession) SRequestId() (*big.Int, error) {
	return _VRFV2PlusRevertingExample.Contract.SRequestId(&_VRFV2PlusRevertingExample.CallOpts)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleCaller) SSubId(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _VRFV2PlusRevertingExample.contract.Call(opts, &out, "s_subId")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) SSubId() (uint64, error) {
	return _VRFV2PlusRevertingExample.Contract.SSubId(&_VRFV2PlusRevertingExample.CallOpts)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleCallerSession) SSubId() (uint64, error) {
	return _VRFV2PlusRevertingExample.Contract.SSubId(&_VRFV2PlusRevertingExample.CallOpts)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactor) CreateSubscriptionAndFund(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.contract.Transact(opts, "createSubscriptionAndFund", amount)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) CreateSubscriptionAndFund(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.CreateSubscriptionAndFund(&_VRFV2PlusRevertingExample.TransactOpts, amount)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactorSession) CreateSubscriptionAndFund(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.CreateSubscriptionAndFund(&_VRFV2PlusRevertingExample.TransactOpts, amount)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.contract.Transact(opts, "rawFulfillRandomWords", requestId, randomWords)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.RawFulfillRandomWords(&_VRFV2PlusRevertingExample.TransactOpts, requestId, randomWords)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactorSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.RawFulfillRandomWords(&_VRFV2PlusRevertingExample.TransactOpts, requestId, randomWords)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactor) RequestRandomness(opts *bind.TransactOpts, keyHash [32]byte, subId uint64, minReqConfs uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.contract.Transact(opts, "requestRandomness", keyHash, subId, minReqConfs, callbackGasLimit, numWords)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) RequestRandomness(keyHash [32]byte, subId uint64, minReqConfs uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.RequestRandomness(&_VRFV2PlusRevertingExample.TransactOpts, keyHash, subId, minReqConfs, callbackGasLimit, numWords)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactorSession) RequestRandomness(keyHash [32]byte, subId uint64, minReqConfs uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.RequestRandomness(&_VRFV2PlusRevertingExample.TransactOpts, keyHash, subId, minReqConfs, callbackGasLimit, numWords)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactor) TopUpSubscription(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.contract.Transact(opts, "topUpSubscription", amount)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) TopUpSubscription(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.TopUpSubscription(&_VRFV2PlusRevertingExample.TransactOpts, amount)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactorSession) TopUpSubscription(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.TopUpSubscription(&_VRFV2PlusRevertingExample.TransactOpts, amount)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactor) UpdateSubscription(opts *bind.TransactOpts, consumers []common.Address) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.contract.Transact(opts, "updateSubscription", consumers)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) UpdateSubscription(consumers []common.Address) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.UpdateSubscription(&_VRFV2PlusRevertingExample.TransactOpts, consumers)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactorSession) UpdateSubscription(consumers []common.Address) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.UpdateSubscription(&_VRFV2PlusRevertingExample.TransactOpts, consumers)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExample) Address() common.Address {
	return _VRFV2PlusRevertingExample.address
}

type VRFV2PlusRevertingExampleInterface interface {
	SGasAvailable(opts *bind.CallOpts) (*big.Int, error)

	SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	SRequestId(opts *bind.CallOpts) (*big.Int, error)

	SSubId(opts *bind.CallOpts) (uint64, error)

	CreateSubscriptionAndFund(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	RequestRandomness(opts *bind.TransactOpts, keyHash [32]byte, subId uint64, minReqConfs uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error)

	TopUpSubscription(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	UpdateSubscription(opts *bind.TransactOpts, consumers []common.Address) (*types.Transaction, error)

	Address() common.Address
}
