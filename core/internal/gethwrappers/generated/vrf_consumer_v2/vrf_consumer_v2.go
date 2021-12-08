// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_consumer_v2

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
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated"
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

var VRFConsumerV2MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"numWords\",\"type\":\"uint256\"}],\"name\":\"FulfillCallbackCalled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"expectedRequestId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"actualRequestId\",\"type\":\"uint256\"}],\"name\":\"IncorrectRequestId\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_gasAvailable\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_randomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_subId\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"testCreateSubscriptionAndFund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"minReqConfs\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"}],\"name\":\"testRequestRandomness\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"topUpSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"name\":\"updateSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561001057600080fd5b50604051610e88380380610e8883398101604081905261002f9161008e565b6001600160601b0319606083901b16608052600280546001600160a01b03199081166001600160a01b0394851617909155600380549290931691161790556100c1565b80516001600160a01b038116811461008957600080fd5b919050565b600080604083850312156100a157600080fd5b6100aa83610072565b91506100b860208401610072565b90509250929050565b60805160601c610da26100e66000396000818161019e01526102060152610da26000f3fe608060405234801561001057600080fd5b50600436106100a35760003560e01c80636802f72611610076578063e89e106a1161005b578063e89e106a14610161578063f08c5daa1461016a578063f6eaffc81461017357600080fd5b80636802f72614610109578063706da1ca1461011c57600080fd5b80631fe543e3146100a857806327784fad146100bd5780632fa4e442146100e357806336bfffed146100f6575b600080fd5b6100bb6100b6366004610ab6565b610186565b005b6100d06100cb366004610a1b565b610246565b6040519081526020015b60405180910390f35b6100bb6100f1366004610b77565b610328565b6100bb610104366004610933565b610488565b6100bb610117366004610b77565b610610565b6003546101489074010000000000000000000000000000000000000000900467ffffffffffffffff1681565b60405167ffffffffffffffff90911681526020016100da565b6100d060015481565b6100d060045481565b6100d0610181366004610a84565b610817565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610238576040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001660248201526044015b60405180910390fd5b6102428282610838565b5050565b6002546040517f5d3b1d300000000000000000000000000000000000000000000000000000000081526004810187905267ffffffffffffffff8616602482015261ffff8516604482015263ffffffff80851660648301528316608482015260009173ffffffffffffffffffffffffffffffffffffffff1690635d3b1d309060a401602060405180830381600087803b1580156102e157600080fd5b505af11580156102f5573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103199190610a9d565b60018190559695505050505050565b60035474010000000000000000000000000000000000000000900467ffffffffffffffff166103b3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f737562206e6f7420736574000000000000000000000000000000000000000000604482015260640161022f565b6003546002546040805174010000000000000000000000000000000000000000840467ffffffffffffffff16602082015273ffffffffffffffffffffffffffffffffffffffff93841693634000aea09316918591015b6040516020818303038152906040526040518463ffffffff1660e01b815260040161043693929190610ba5565b602060405180830381600087803b15801561045057600080fd5b505af1158015610464573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061024291906109f2565b60035474010000000000000000000000000000000000000000900467ffffffffffffffff16610513576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f7375624944206e6f742073657400000000000000000000000000000000000000604482015260640161022f565b60005b815181101561024257600254600354835173ffffffffffffffffffffffffffffffffffffffff90921691637341c10c9174010000000000000000000000000000000000000000900467ffffffffffffffff169085908590811061057b5761057b610d1e565b60200260200101516040518363ffffffff1660e01b81526004016105cb92919067ffffffffffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b600060405180830381600087803b1580156105e557600080fd5b505af11580156105f9573d6000803e3d6000fd5b50505050808061060890610cbe565b915050610516565b60035474010000000000000000000000000000000000000000900467ffffffffffffffff166103b357600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a21a23e46040518163ffffffff1660e01b8152600401602060405180830381600087803b1580156106a357600080fd5b505af11580156106b7573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106db9190610b5a565b600380547fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000067ffffffffffffffff938416810291909117918290556002546040517f7341c10c00000000000000000000000000000000000000000000000000000000815291909204909216600483015230602483015273ffffffffffffffffffffffffffffffffffffffff1690637341c10c90604401600060405180830381600087803b1580156107a257600080fd5b505af11580156107b6573d6000803e3d6000fd5b50506003546002546040805174010000000000000000000000000000000000000000840467ffffffffffffffff16602082015273ffffffffffffffffffffffffffffffffffffffff9384169550634000aea094509290911691859101610409565b6000818154811061082757600080fd5b600091825260209091200154905081565b805160405183907f1cb9f4763b3c7742e8319c07038cb222b705bf70dd94b262c8d80188d114e48d90600090a3600154821461089e576001546040518391907fd8b88f426d1d9840537c232eced585ba09902343e1fb72ccd5ee7c1f6d92ab1090600090a35b5a60045580516108b59060009060208401906108ba565b505050565b8280548282559060005260206000209081019282156108f5579160200282015b828111156108f55782518255916020019190600101906108da565b50610901929150610905565b5090565b5b808211156109015760008155600101610906565b803563ffffffff8116811461092e57600080fd5b919050565b6000602080838503121561094657600080fd5b823567ffffffffffffffff81111561095d57600080fd5b8301601f8101851361096e57600080fd5b803561098161097c82610c9a565b610c4b565b80828252848201915084840188868560051b87010111156109a157600080fd5b60009450845b848110156109e457813573ffffffffffffffffffffffffffffffffffffffff811681146109d2578687fd5b845292860192908601906001016109a7565b509098975050505050505050565b600060208284031215610a0457600080fd5b81518015158114610a1457600080fd5b9392505050565b600080600080600060a08688031215610a3357600080fd5b853594506020860135610a4581610d7c565b9350604086013561ffff81168114610a5c57600080fd5b9250610a6a6060870161091a565b9150610a786080870161091a565b90509295509295909350565b600060208284031215610a9657600080fd5b5035919050565b600060208284031215610aaf57600080fd5b5051919050565b60008060408385031215610ac957600080fd5b8235915060208084013567ffffffffffffffff811115610ae857600080fd5b8401601f81018613610af957600080fd5b8035610b0761097c82610c9a565b80828252848201915084840189868560051b8701011115610b2757600080fd5b600094505b83851015610b4a578035835260019490940193918501918501610b2c565b5080955050505050509250929050565b600060208284031215610b6c57600080fd5b8151610a1481610d7c565b600060208284031215610b8957600080fd5b81356bffffffffffffffffffffffff81168114610a1457600080fd5b73ffffffffffffffffffffffffffffffffffffffff84168152600060206bffffffffffffffffffffffff85168184015260606040840152835180606085015260005b81811015610c0357858101830151858201608001528201610be7565b81811115610c15576000608083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160800195945050505050565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610c9257610c92610d4d565b604052919050565b600067ffffffffffffffff821115610cb457610cb4610d4d565b5060051b60200190565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415610d17577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b67ffffffffffffffff81168114610d9257600080fd5b5056fea164736f6c6343000806000a",
}

var VRFConsumerV2ABI = VRFConsumerV2MetaData.ABI

var VRFConsumerV2Bin = VRFConsumerV2MetaData.Bin

func DeployVRFConsumerV2(auth *bind.TransactOpts, backend bind.ContractBackend, vrfCoordinator common.Address, link common.Address) (common.Address, *types.Transaction, *VRFConsumerV2, error) {
	parsed, err := VRFConsumerV2MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFConsumerV2Bin), backend, vrfCoordinator, link)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFConsumerV2{VRFConsumerV2Caller: VRFConsumerV2Caller{contract: contract}, VRFConsumerV2Transactor: VRFConsumerV2Transactor{contract: contract}, VRFConsumerV2Filterer: VRFConsumerV2Filterer{contract: contract}}, nil
}

type VRFConsumerV2 struct {
	address common.Address
	abi     abi.ABI
	VRFConsumerV2Caller
	VRFConsumerV2Transactor
	VRFConsumerV2Filterer
}

type VRFConsumerV2Caller struct {
	contract *bind.BoundContract
}

type VRFConsumerV2Transactor struct {
	contract *bind.BoundContract
}

type VRFConsumerV2Filterer struct {
	contract *bind.BoundContract
}

type VRFConsumerV2Session struct {
	Contract     *VRFConsumerV2
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFConsumerV2CallerSession struct {
	Contract *VRFConsumerV2Caller
	CallOpts bind.CallOpts
}

type VRFConsumerV2TransactorSession struct {
	Contract     *VRFConsumerV2Transactor
	TransactOpts bind.TransactOpts
}

type VRFConsumerV2Raw struct {
	Contract *VRFConsumerV2
}

type VRFConsumerV2CallerRaw struct {
	Contract *VRFConsumerV2Caller
}

type VRFConsumerV2TransactorRaw struct {
	Contract *VRFConsumerV2Transactor
}

func NewVRFConsumerV2(address common.Address, backend bind.ContractBackend) (*VRFConsumerV2, error) {
	abi, err := abi.JSON(strings.NewReader(VRFConsumerV2ABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFConsumerV2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerV2{address: address, abi: abi, VRFConsumerV2Caller: VRFConsumerV2Caller{contract: contract}, VRFConsumerV2Transactor: VRFConsumerV2Transactor{contract: contract}, VRFConsumerV2Filterer: VRFConsumerV2Filterer{contract: contract}}, nil
}

func NewVRFConsumerV2Caller(address common.Address, caller bind.ContractCaller) (*VRFConsumerV2Caller, error) {
	contract, err := bindVRFConsumerV2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerV2Caller{contract: contract}, nil
}

func NewVRFConsumerV2Transactor(address common.Address, transactor bind.ContractTransactor) (*VRFConsumerV2Transactor, error) {
	contract, err := bindVRFConsumerV2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerV2Transactor{contract: contract}, nil
}

func NewVRFConsumerV2Filterer(address common.Address, filterer bind.ContractFilterer) (*VRFConsumerV2Filterer, error) {
	contract, err := bindVRFConsumerV2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerV2Filterer{contract: contract}, nil
}

func bindVRFConsumerV2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFConsumerV2ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_VRFConsumerV2 *VRFConsumerV2Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFConsumerV2.Contract.VRFConsumerV2Caller.contract.Call(opts, result, method, params...)
}

func (_VRFConsumerV2 *VRFConsumerV2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.VRFConsumerV2Transactor.contract.Transfer(opts)
}

func (_VRFConsumerV2 *VRFConsumerV2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.VRFConsumerV2Transactor.contract.Transact(opts, method, params...)
}

func (_VRFConsumerV2 *VRFConsumerV2CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFConsumerV2.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFConsumerV2 *VRFConsumerV2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.contract.Transfer(opts)
}

func (_VRFConsumerV2 *VRFConsumerV2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.contract.Transact(opts, method, params...)
}

func (_VRFConsumerV2 *VRFConsumerV2Caller) SGasAvailable(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFConsumerV2.contract.Call(opts, &out, "s_gasAvailable")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFConsumerV2 *VRFConsumerV2Session) SGasAvailable() (*big.Int, error) {
	return _VRFConsumerV2.Contract.SGasAvailable(&_VRFConsumerV2.CallOpts)
}

func (_VRFConsumerV2 *VRFConsumerV2CallerSession) SGasAvailable() (*big.Int, error) {
	return _VRFConsumerV2.Contract.SGasAvailable(&_VRFConsumerV2.CallOpts)
}

func (_VRFConsumerV2 *VRFConsumerV2Caller) SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VRFConsumerV2.contract.Call(opts, &out, "s_randomWords", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFConsumerV2 *VRFConsumerV2Session) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFConsumerV2.Contract.SRandomWords(&_VRFConsumerV2.CallOpts, arg0)
}

func (_VRFConsumerV2 *VRFConsumerV2CallerSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFConsumerV2.Contract.SRandomWords(&_VRFConsumerV2.CallOpts, arg0)
}

func (_VRFConsumerV2 *VRFConsumerV2Caller) SRequestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFConsumerV2.contract.Call(opts, &out, "s_requestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFConsumerV2 *VRFConsumerV2Session) SRequestId() (*big.Int, error) {
	return _VRFConsumerV2.Contract.SRequestId(&_VRFConsumerV2.CallOpts)
}

func (_VRFConsumerV2 *VRFConsumerV2CallerSession) SRequestId() (*big.Int, error) {
	return _VRFConsumerV2.Contract.SRequestId(&_VRFConsumerV2.CallOpts)
}

func (_VRFConsumerV2 *VRFConsumerV2Caller) SSubId(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _VRFConsumerV2.contract.Call(opts, &out, "s_subId")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_VRFConsumerV2 *VRFConsumerV2Session) SSubId() (uint64, error) {
	return _VRFConsumerV2.Contract.SSubId(&_VRFConsumerV2.CallOpts)
}

func (_VRFConsumerV2 *VRFConsumerV2CallerSession) SSubId() (uint64, error) {
	return _VRFConsumerV2.Contract.SSubId(&_VRFConsumerV2.CallOpts)
}

func (_VRFConsumerV2 *VRFConsumerV2Transactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2.contract.Transact(opts, "rawFulfillRandomWords", requestId, randomWords)
}

func (_VRFConsumerV2 *VRFConsumerV2Session) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.RawFulfillRandomWords(&_VRFConsumerV2.TransactOpts, requestId, randomWords)
}

func (_VRFConsumerV2 *VRFConsumerV2TransactorSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.RawFulfillRandomWords(&_VRFConsumerV2.TransactOpts, requestId, randomWords)
}

func (_VRFConsumerV2 *VRFConsumerV2Transactor) TestCreateSubscriptionAndFund(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2.contract.Transact(opts, "testCreateSubscriptionAndFund", amount)
}

func (_VRFConsumerV2 *VRFConsumerV2Session) TestCreateSubscriptionAndFund(amount *big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.TestCreateSubscriptionAndFund(&_VRFConsumerV2.TransactOpts, amount)
}

func (_VRFConsumerV2 *VRFConsumerV2TransactorSession) TestCreateSubscriptionAndFund(amount *big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.TestCreateSubscriptionAndFund(&_VRFConsumerV2.TransactOpts, amount)
}

func (_VRFConsumerV2 *VRFConsumerV2Transactor) TestRequestRandomness(opts *bind.TransactOpts, keyHash [32]byte, subId uint64, minReqConfs uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _VRFConsumerV2.contract.Transact(opts, "testRequestRandomness", keyHash, subId, minReqConfs, callbackGasLimit, numWords)
}

func (_VRFConsumerV2 *VRFConsumerV2Session) TestRequestRandomness(keyHash [32]byte, subId uint64, minReqConfs uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.TestRequestRandomness(&_VRFConsumerV2.TransactOpts, keyHash, subId, minReqConfs, callbackGasLimit, numWords)
}

func (_VRFConsumerV2 *VRFConsumerV2TransactorSession) TestRequestRandomness(keyHash [32]byte, subId uint64, minReqConfs uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.TestRequestRandomness(&_VRFConsumerV2.TransactOpts, keyHash, subId, minReqConfs, callbackGasLimit, numWords)
}

func (_VRFConsumerV2 *VRFConsumerV2Transactor) TopUpSubscription(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2.contract.Transact(opts, "topUpSubscription", amount)
}

func (_VRFConsumerV2 *VRFConsumerV2Session) TopUpSubscription(amount *big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.TopUpSubscription(&_VRFConsumerV2.TransactOpts, amount)
}

func (_VRFConsumerV2 *VRFConsumerV2TransactorSession) TopUpSubscription(amount *big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.TopUpSubscription(&_VRFConsumerV2.TransactOpts, amount)
}

func (_VRFConsumerV2 *VRFConsumerV2Transactor) UpdateSubscription(opts *bind.TransactOpts, consumers []common.Address) (*types.Transaction, error) {
	return _VRFConsumerV2.contract.Transact(opts, "updateSubscription", consumers)
}

func (_VRFConsumerV2 *VRFConsumerV2Session) UpdateSubscription(consumers []common.Address) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.UpdateSubscription(&_VRFConsumerV2.TransactOpts, consumers)
}

func (_VRFConsumerV2 *VRFConsumerV2TransactorSession) UpdateSubscription(consumers []common.Address) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.UpdateSubscription(&_VRFConsumerV2.TransactOpts, consumers)
}

type VRFConsumerV2FulfillCallbackCalledIterator struct {
	Event *VRFConsumerV2FulfillCallbackCalled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFConsumerV2FulfillCallbackCalledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFConsumerV2FulfillCallbackCalled)
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
		it.Event = new(VRFConsumerV2FulfillCallbackCalled)
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

func (it *VRFConsumerV2FulfillCallbackCalledIterator) Error() error {
	return it.fail
}

func (it *VRFConsumerV2FulfillCallbackCalledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFConsumerV2FulfillCallbackCalled struct {
	RequestId *big.Int
	NumWords  *big.Int
	Raw       types.Log
}

func (_VRFConsumerV2 *VRFConsumerV2Filterer) FilterFulfillCallbackCalled(opts *bind.FilterOpts, requestId []*big.Int, numWords []*big.Int) (*VRFConsumerV2FulfillCallbackCalledIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var numWordsRule []interface{}
	for _, numWordsItem := range numWords {
		numWordsRule = append(numWordsRule, numWordsItem)
	}

	logs, sub, err := _VRFConsumerV2.contract.FilterLogs(opts, "FulfillCallbackCalled", requestIdRule, numWordsRule)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerV2FulfillCallbackCalledIterator{contract: _VRFConsumerV2.contract, event: "FulfillCallbackCalled", logs: logs, sub: sub}, nil
}

func (_VRFConsumerV2 *VRFConsumerV2Filterer) WatchFulfillCallbackCalled(opts *bind.WatchOpts, sink chan<- *VRFConsumerV2FulfillCallbackCalled, requestId []*big.Int, numWords []*big.Int) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var numWordsRule []interface{}
	for _, numWordsItem := range numWords {
		numWordsRule = append(numWordsRule, numWordsItem)
	}

	logs, sub, err := _VRFConsumerV2.contract.WatchLogs(opts, "FulfillCallbackCalled", requestIdRule, numWordsRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFConsumerV2FulfillCallbackCalled)
				if err := _VRFConsumerV2.contract.UnpackLog(event, "FulfillCallbackCalled", log); err != nil {
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

func (_VRFConsumerV2 *VRFConsumerV2Filterer) ParseFulfillCallbackCalled(log types.Log) (*VRFConsumerV2FulfillCallbackCalled, error) {
	event := new(VRFConsumerV2FulfillCallbackCalled)
	if err := _VRFConsumerV2.contract.UnpackLog(event, "FulfillCallbackCalled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFConsumerV2IncorrectRequestIdIterator struct {
	Event *VRFConsumerV2IncorrectRequestId

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFConsumerV2IncorrectRequestIdIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFConsumerV2IncorrectRequestId)
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
		it.Event = new(VRFConsumerV2IncorrectRequestId)
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

func (it *VRFConsumerV2IncorrectRequestIdIterator) Error() error {
	return it.fail
}

func (it *VRFConsumerV2IncorrectRequestIdIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFConsumerV2IncorrectRequestId struct {
	ExpectedRequestId *big.Int
	ActualRequestId   *big.Int
	Raw               types.Log
}

func (_VRFConsumerV2 *VRFConsumerV2Filterer) FilterIncorrectRequestId(opts *bind.FilterOpts, expectedRequestId []*big.Int, actualRequestId []*big.Int) (*VRFConsumerV2IncorrectRequestIdIterator, error) {

	var expectedRequestIdRule []interface{}
	for _, expectedRequestIdItem := range expectedRequestId {
		expectedRequestIdRule = append(expectedRequestIdRule, expectedRequestIdItem)
	}
	var actualRequestIdRule []interface{}
	for _, actualRequestIdItem := range actualRequestId {
		actualRequestIdRule = append(actualRequestIdRule, actualRequestIdItem)
	}

	logs, sub, err := _VRFConsumerV2.contract.FilterLogs(opts, "IncorrectRequestId", expectedRequestIdRule, actualRequestIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerV2IncorrectRequestIdIterator{contract: _VRFConsumerV2.contract, event: "IncorrectRequestId", logs: logs, sub: sub}, nil
}

func (_VRFConsumerV2 *VRFConsumerV2Filterer) WatchIncorrectRequestId(opts *bind.WatchOpts, sink chan<- *VRFConsumerV2IncorrectRequestId, expectedRequestId []*big.Int, actualRequestId []*big.Int) (event.Subscription, error) {

	var expectedRequestIdRule []interface{}
	for _, expectedRequestIdItem := range expectedRequestId {
		expectedRequestIdRule = append(expectedRequestIdRule, expectedRequestIdItem)
	}
	var actualRequestIdRule []interface{}
	for _, actualRequestIdItem := range actualRequestId {
		actualRequestIdRule = append(actualRequestIdRule, actualRequestIdItem)
	}

	logs, sub, err := _VRFConsumerV2.contract.WatchLogs(opts, "IncorrectRequestId", expectedRequestIdRule, actualRequestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFConsumerV2IncorrectRequestId)
				if err := _VRFConsumerV2.contract.UnpackLog(event, "IncorrectRequestId", log); err != nil {
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

func (_VRFConsumerV2 *VRFConsumerV2Filterer) ParseIncorrectRequestId(log types.Log) (*VRFConsumerV2IncorrectRequestId, error) {
	event := new(VRFConsumerV2IncorrectRequestId)
	if err := _VRFConsumerV2.contract.UnpackLog(event, "IncorrectRequestId", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_VRFConsumerV2 *VRFConsumerV2) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFConsumerV2.abi.Events["FulfillCallbackCalled"].ID:
		return _VRFConsumerV2.ParseFulfillCallbackCalled(log)
	case _VRFConsumerV2.abi.Events["IncorrectRequestId"].ID:
		return _VRFConsumerV2.ParseIncorrectRequestId(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFConsumerV2FulfillCallbackCalled) Topic() common.Hash {
	return common.HexToHash("0x1cb9f4763b3c7742e8319c07038cb222b705bf70dd94b262c8d80188d114e48d")
}

func (VRFConsumerV2IncorrectRequestId) Topic() common.Hash {
	return common.HexToHash("0xd8b88f426d1d9840537c232eced585ba09902343e1fb72ccd5ee7c1f6d92ab10")
}

func (_VRFConsumerV2 *VRFConsumerV2) Address() common.Address {
	return _VRFConsumerV2.address
}

type VRFConsumerV2Interface interface {
	SGasAvailable(opts *bind.CallOpts) (*big.Int, error)

	SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	SRequestId(opts *bind.CallOpts) (*big.Int, error)

	SSubId(opts *bind.CallOpts) (uint64, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	TestCreateSubscriptionAndFund(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	TestRequestRandomness(opts *bind.TransactOpts, keyHash [32]byte, subId uint64, minReqConfs uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error)

	TopUpSubscription(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	UpdateSubscription(opts *bind.TransactOpts, consumers []common.Address) (*types.Transaction, error)

	FilterFulfillCallbackCalled(opts *bind.FilterOpts, requestId []*big.Int, numWords []*big.Int) (*VRFConsumerV2FulfillCallbackCalledIterator, error)

	WatchFulfillCallbackCalled(opts *bind.WatchOpts, sink chan<- *VRFConsumerV2FulfillCallbackCalled, requestId []*big.Int, numWords []*big.Int) (event.Subscription, error)

	ParseFulfillCallbackCalled(log types.Log) (*VRFConsumerV2FulfillCallbackCalled, error)

	FilterIncorrectRequestId(opts *bind.FilterOpts, expectedRequestId []*big.Int, actualRequestId []*big.Int) (*VRFConsumerV2IncorrectRequestIdIterator, error)

	WatchIncorrectRequestId(opts *bind.WatchOpts, sink chan<- *VRFConsumerV2IncorrectRequestId, expectedRequestId []*big.Int, actualRequestId []*big.Int) (event.Subscription, error)

	ParseIncorrectRequestId(log types.Log) (*VRFConsumerV2IncorrectRequestId, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
