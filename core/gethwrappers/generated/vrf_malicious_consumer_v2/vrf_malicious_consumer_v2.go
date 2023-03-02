// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_malicious_consumer_v2

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

var VRFMaliciousConsumerV2MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlyCoordinatorCanFulfill\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"createSubscriptionAndFund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"requestRandomness\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_gasAvailable\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_randomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_subId\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"setKeyHash\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"name\":\"updateSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561001057600080fd5b50604051610d9c380380610d9c83398101604081905261002f9161008e565b6001600160601b0319606083901b16608052600280546001600160a01b03199081166001600160a01b0394851617909155600380549290931691161790556100c1565b80516001600160a01b038116811461008957600080fd5b919050565b600080604083850312156100a157600080fd5b6100aa83610072565b91506100b860208401610072565b90509250929050565b60805160601c610cb66100e66000396000818161019301526101fb0152610cb66000f3fe608060405234801561001057600080fd5b50600436106100a35760003560e01c8063cf62c8ab11610076578063f08c5daa1161005b578063f08c5daa14610157578063f6eaffc814610160578063f8413b071461017357600080fd5b8063cf62c8ab1461012d578063e89e106a1461014057600080fd5b80631fe543e3146100a857806336bfffed146100bd578063706da1ca146100d0578063985447101461011a575b600080fd5b6100bb6100b63660046109d6565b61017b565b005b6100bb6100cb3660046108bc565b61023b565b6003546100fc9074010000000000000000000000000000000000000000900467ffffffffffffffff1681565b60405167ffffffffffffffff90911681526020015b60405180910390f35b6100bb6101283660046109a4565b600555565b6100bb61013b366004610aa4565b6103c3565b61014960015481565b604051908152602001610111565b61014960045481565b61014961016e3660046109a4565b610642565b610149610663565b3373ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161461022d576040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001660248201526044015b60405180910390fd5b6102378282610752565b5050565b60035474010000000000000000000000000000000000000000900467ffffffffffffffff166102c6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f7375624944206e6f7420736574000000000000000000000000000000000000006044820152606401610224565b60005b815181101561023757600254600354835173ffffffffffffffffffffffffffffffffffffffff90921691637341c10c9174010000000000000000000000000000000000000000900467ffffffffffffffff169085908590811061032e5761032e610c4b565b60200260200101516040518363ffffffff1660e01b815260040161037e92919067ffffffffffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b600060405180830381600087803b15801561039857600080fd5b505af11580156103ac573d6000803e3d6000fd5b5050505080806103bb90610beb565b9150506102c9565b60035474010000000000000000000000000000000000000000900467ffffffffffffffff1661056e57600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a21a23e46040518163ffffffff1660e01b8152600401602060405180830381600087803b15801561045657600080fd5b505af115801561046a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061048e9190610a7a565b600380547fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000067ffffffffffffffff938416810291909117918290556002546040517f7341c10c00000000000000000000000000000000000000000000000000000000815291909204909216600483015230602483015273ffffffffffffffffffffffffffffffffffffffff1690637341c10c90604401600060405180830381600087803b15801561055557600080fd5b505af1158015610569573d6000803e3d6000fd5b505050505b6003546002546040805174010000000000000000000000000000000000000000840467ffffffffffffffff16602082015273ffffffffffffffffffffffffffffffffffffffff93841693634000aea09316918591016040516020818303038152906040526040518463ffffffff1660e01b81526004016105f093929190610ad2565b602060405180830381600087803b15801561060a57600080fd5b505af115801561061e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610237919061097b565b6000818154811061065257600080fd5b600091825260209091200154905081565b6002546005546003546040517f5d3b1d30000000000000000000000000000000000000000000000000000000008152600481019290925274010000000000000000000000000000000000000000900467ffffffffffffffff1660248201526001604482018190526207a1206064830152608482015260009173ffffffffffffffffffffffffffffffffffffffff1690635d3b1d309060a401602060405180830381600087803b15801561071557600080fd5b505af1158015610729573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061074d91906109bd565b905090565b5a600455805161076990600090602084019061085c565b5060018281556002546005546003546040517f5d3b1d30000000000000000000000000000000000000000000000000000000008152600481019290925274010000000000000000000000000000000000000000900467ffffffffffffffff1660248201526044810183905262030d406064820152608481019290925273ffffffffffffffffffffffffffffffffffffffff1690635d3b1d309060a401602060405180830381600087803b15801561081f57600080fd5b505af1158015610833573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061085791906109bd565b505050565b828054828255906000526020600020908101928215610897579160200282015b8281111561089757825182559160200191906001019061087c565b506108a39291506108a7565b5090565b5b808211156108a357600081556001016108a8565b600060208083850312156108cf57600080fd5b823567ffffffffffffffff8111156108e657600080fd5b8301601f810185136108f757600080fd5b803561090a61090582610bc7565b610b78565b80828252848201915084840188868560051b870101111561092a57600080fd5b60009450845b8481101561096d57813573ffffffffffffffffffffffffffffffffffffffff8116811461095b578687fd5b84529286019290860190600101610930565b509098975050505050505050565b60006020828403121561098d57600080fd5b8151801515811461099d57600080fd5b9392505050565b6000602082840312156109b657600080fd5b5035919050565b6000602082840312156109cf57600080fd5b5051919050565b600080604083850312156109e957600080fd5b8235915060208084013567ffffffffffffffff811115610a0857600080fd5b8401601f81018613610a1957600080fd5b8035610a2761090582610bc7565b80828252848201915084840189868560051b8701011115610a4757600080fd5b600094505b83851015610a6a578035835260019490940193918501918501610a4c565b5080955050505050509250929050565b600060208284031215610a8c57600080fd5b815167ffffffffffffffff8116811461099d57600080fd5b600060208284031215610ab657600080fd5b81356bffffffffffffffffffffffff8116811461099d57600080fd5b73ffffffffffffffffffffffffffffffffffffffff84168152600060206bffffffffffffffffffffffff85168184015260606040840152835180606085015260005b81811015610b3057858101830151858201608001528201610b14565b81811115610b42576000608083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160800195945050505050565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610bbf57610bbf610c7a565b604052919050565b600067ffffffffffffffff821115610be157610be1610c7a565b5060051b60200190565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415610c44577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var VRFMaliciousConsumerV2ABI = VRFMaliciousConsumerV2MetaData.ABI

var VRFMaliciousConsumerV2Bin = VRFMaliciousConsumerV2MetaData.Bin

func DeployVRFMaliciousConsumerV2(auth *bind.TransactOpts, backend bind.ContractBackend, vrfCoordinator common.Address, link common.Address) (common.Address, *types.Transaction, *VRFMaliciousConsumerV2, error) {
	parsed, err := VRFMaliciousConsumerV2MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFMaliciousConsumerV2Bin), backend, vrfCoordinator, link)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFMaliciousConsumerV2{VRFMaliciousConsumerV2Caller: VRFMaliciousConsumerV2Caller{contract: contract}, VRFMaliciousConsumerV2Transactor: VRFMaliciousConsumerV2Transactor{contract: contract}, VRFMaliciousConsumerV2Filterer: VRFMaliciousConsumerV2Filterer{contract: contract}}, nil
}

type VRFMaliciousConsumerV2 struct {
	address common.Address
	abi     abi.ABI
	VRFMaliciousConsumerV2Caller
	VRFMaliciousConsumerV2Transactor
	VRFMaliciousConsumerV2Filterer
}

type VRFMaliciousConsumerV2Caller struct {
	contract *bind.BoundContract
}

type VRFMaliciousConsumerV2Transactor struct {
	contract *bind.BoundContract
}

type VRFMaliciousConsumerV2Filterer struct {
	contract *bind.BoundContract
}

type VRFMaliciousConsumerV2Session struct {
	Contract     *VRFMaliciousConsumerV2
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFMaliciousConsumerV2CallerSession struct {
	Contract *VRFMaliciousConsumerV2Caller
	CallOpts bind.CallOpts
}

type VRFMaliciousConsumerV2TransactorSession struct {
	Contract     *VRFMaliciousConsumerV2Transactor
	TransactOpts bind.TransactOpts
}

type VRFMaliciousConsumerV2Raw struct {
	Contract *VRFMaliciousConsumerV2
}

type VRFMaliciousConsumerV2CallerRaw struct {
	Contract *VRFMaliciousConsumerV2Caller
}

type VRFMaliciousConsumerV2TransactorRaw struct {
	Contract *VRFMaliciousConsumerV2Transactor
}

func NewVRFMaliciousConsumerV2(address common.Address, backend bind.ContractBackend) (*VRFMaliciousConsumerV2, error) {
	abi, err := abi.JSON(strings.NewReader(VRFMaliciousConsumerV2ABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFMaliciousConsumerV2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFMaliciousConsumerV2{address: address, abi: abi, VRFMaliciousConsumerV2Caller: VRFMaliciousConsumerV2Caller{contract: contract}, VRFMaliciousConsumerV2Transactor: VRFMaliciousConsumerV2Transactor{contract: contract}, VRFMaliciousConsumerV2Filterer: VRFMaliciousConsumerV2Filterer{contract: contract}}, nil
}

func NewVRFMaliciousConsumerV2Caller(address common.Address, caller bind.ContractCaller) (*VRFMaliciousConsumerV2Caller, error) {
	contract, err := bindVRFMaliciousConsumerV2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFMaliciousConsumerV2Caller{contract: contract}, nil
}

func NewVRFMaliciousConsumerV2Transactor(address common.Address, transactor bind.ContractTransactor) (*VRFMaliciousConsumerV2Transactor, error) {
	contract, err := bindVRFMaliciousConsumerV2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFMaliciousConsumerV2Transactor{contract: contract}, nil
}

func NewVRFMaliciousConsumerV2Filterer(address common.Address, filterer bind.ContractFilterer) (*VRFMaliciousConsumerV2Filterer, error) {
	contract, err := bindVRFMaliciousConsumerV2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFMaliciousConsumerV2Filterer{contract: contract}, nil
}

func bindVRFMaliciousConsumerV2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFMaliciousConsumerV2ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFMaliciousConsumerV2.Contract.VRFMaliciousConsumerV2Caller.contract.Call(opts, result, method, params...)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2.Contract.VRFMaliciousConsumerV2Transactor.contract.Transfer(opts)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2.Contract.VRFMaliciousConsumerV2Transactor.contract.Transact(opts, method, params...)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFMaliciousConsumerV2.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2.Contract.contract.Transfer(opts)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2.Contract.contract.Transact(opts, method, params...)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2Caller) SGasAvailable(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFMaliciousConsumerV2.contract.Call(opts, &out, "s_gasAvailable")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2Session) SGasAvailable() (*big.Int, error) {
	return _VRFMaliciousConsumerV2.Contract.SGasAvailable(&_VRFMaliciousConsumerV2.CallOpts)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2CallerSession) SGasAvailable() (*big.Int, error) {
	return _VRFMaliciousConsumerV2.Contract.SGasAvailable(&_VRFMaliciousConsumerV2.CallOpts)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2Caller) SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VRFMaliciousConsumerV2.contract.Call(opts, &out, "s_randomWords", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2Session) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFMaliciousConsumerV2.Contract.SRandomWords(&_VRFMaliciousConsumerV2.CallOpts, arg0)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2CallerSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFMaliciousConsumerV2.Contract.SRandomWords(&_VRFMaliciousConsumerV2.CallOpts, arg0)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2Caller) SRequestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFMaliciousConsumerV2.contract.Call(opts, &out, "s_requestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2Session) SRequestId() (*big.Int, error) {
	return _VRFMaliciousConsumerV2.Contract.SRequestId(&_VRFMaliciousConsumerV2.CallOpts)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2CallerSession) SRequestId() (*big.Int, error) {
	return _VRFMaliciousConsumerV2.Contract.SRequestId(&_VRFMaliciousConsumerV2.CallOpts)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2Caller) SSubId(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _VRFMaliciousConsumerV2.contract.Call(opts, &out, "s_subId")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2Session) SSubId() (uint64, error) {
	return _VRFMaliciousConsumerV2.Contract.SSubId(&_VRFMaliciousConsumerV2.CallOpts)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2CallerSession) SSubId() (uint64, error) {
	return _VRFMaliciousConsumerV2.Contract.SSubId(&_VRFMaliciousConsumerV2.CallOpts)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2Transactor) CreateSubscriptionAndFund(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2.contract.Transact(opts, "createSubscriptionAndFund", amount)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2Session) CreateSubscriptionAndFund(amount *big.Int) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2.Contract.CreateSubscriptionAndFund(&_VRFMaliciousConsumerV2.TransactOpts, amount)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2TransactorSession) CreateSubscriptionAndFund(amount *big.Int) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2.Contract.CreateSubscriptionAndFund(&_VRFMaliciousConsumerV2.TransactOpts, amount)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2Transactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2.contract.Transact(opts, "rawFulfillRandomWords", requestId, randomWords)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2Session) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2.Contract.RawFulfillRandomWords(&_VRFMaliciousConsumerV2.TransactOpts, requestId, randomWords)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2TransactorSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2.Contract.RawFulfillRandomWords(&_VRFMaliciousConsumerV2.TransactOpts, requestId, randomWords)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2Transactor) RequestRandomness(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2.contract.Transact(opts, "requestRandomness")
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2Session) RequestRandomness() (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2.Contract.RequestRandomness(&_VRFMaliciousConsumerV2.TransactOpts)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2TransactorSession) RequestRandomness() (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2.Contract.RequestRandomness(&_VRFMaliciousConsumerV2.TransactOpts)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2Transactor) SetKeyHash(opts *bind.TransactOpts, keyHash [32]byte) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2.contract.Transact(opts, "setKeyHash", keyHash)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2Session) SetKeyHash(keyHash [32]byte) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2.Contract.SetKeyHash(&_VRFMaliciousConsumerV2.TransactOpts, keyHash)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2TransactorSession) SetKeyHash(keyHash [32]byte) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2.Contract.SetKeyHash(&_VRFMaliciousConsumerV2.TransactOpts, keyHash)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2Transactor) UpdateSubscription(opts *bind.TransactOpts, consumers []common.Address) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2.contract.Transact(opts, "updateSubscription", consumers)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2Session) UpdateSubscription(consumers []common.Address) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2.Contract.UpdateSubscription(&_VRFMaliciousConsumerV2.TransactOpts, consumers)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2TransactorSession) UpdateSubscription(consumers []common.Address) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2.Contract.UpdateSubscription(&_VRFMaliciousConsumerV2.TransactOpts, consumers)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2) Address() common.Address {
	return _VRFMaliciousConsumerV2.address
}

type VRFMaliciousConsumerV2Interface interface {
	SGasAvailable(opts *bind.CallOpts) (*big.Int, error)

	SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	SRequestId(opts *bind.CallOpts) (*big.Int, error)

	SSubId(opts *bind.CallOpts) (uint64, error)

	CreateSubscriptionAndFund(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	RequestRandomness(opts *bind.TransactOpts) (*types.Transaction, error)

	SetKeyHash(opts *bind.TransactOpts, keyHash [32]byte) (*types.Transaction, error)

	UpdateSubscription(opts *bind.TransactOpts, consumers []common.Address) (*types.Transaction, error)

	Address() common.Address
}
