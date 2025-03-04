// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package weth9

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

var WETH9MetaData = &bind.MetaData{
	ABI: "[{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"allowance\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"guy\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"wad\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transfer\",\"inputs\":[{\"name\":\"dst\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"wad\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"src\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"dst\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"wad\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"wad\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"src\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"guy\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"wad\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Deposit\",\"inputs\":[{\"name\":\"dst\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"wad\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"src\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"dst\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"wad\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdrawal\",\"inputs\":[{\"name\":\"src\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"wad\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false}]",
	Bin: "0x60c0604052600d60809081526c2bb930b83832b21022ba3432b960991b60a05260009061002c9082610114565b506040805180820190915260048152630ae8aa8960e31b60208201526001906100559082610114565b506002805460ff1916601217905534801561006f57600080fd5b506101d3565b634e487b7160e01b600052604160045260246000fd5b600181811c9082168061009f57607f821691505b6020821081036100bf57634e487b7160e01b600052602260045260246000fd5b50919050565b601f82111561010f57600081815260208120601f850160051c810160208610156100ec5750805b601f850160051c820191505b8181101561010b578281556001016100f8565b5050505b505050565b81516001600160401b0381111561012d5761012d610075565b6101418161013b845461008b565b846100c5565b602080601f831160018114610176576000841561015e5750858301515b600019600386901b1c1916600185901b17855561010b565b600085815260208120601f198616915b828110156101a557888601518255948401946001909101908401610186565b50858210156101c35787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b6108ac806101e26000396000f3fe6080604052600436106100c05760003560e01c8063313ce56711610074578063a9059cbb1161004e578063a9059cbb146101fa578063d0e30db01461021a578063dd62ed3e1461022257600080fd5b8063313ce5671461018c57806370a08231146101b857806395d89b41146101e557600080fd5b806318160ddd116100a557806318160ddd1461012f57806323b872dd1461014c5780632e1a7d4d1461016c57600080fd5b806306fdde03146100d4578063095ea7b3146100ff57600080fd5b366100cf576100cd61025a565b005b600080fd5b3480156100e057600080fd5b506100e96102b5565b6040516100f69190610695565b60405180910390f35b34801561010b57600080fd5b5061011f61011a36600461072a565b610343565b60405190151581526020016100f6565b34801561013b57600080fd5b50475b6040519081526020016100f6565b34801561015857600080fd5b5061011f610167366004610754565b6103bd565b34801561017857600080fd5b506100cd610187366004610790565b6105c4565b34801561019857600080fd5b506002546101a69060ff1681565b60405160ff90911681526020016100f6565b3480156101c457600080fd5b5061013e6101d33660046107a9565b60036020526000908152604090205481565b3480156101f157600080fd5b506100e961066a565b34801561020657600080fd5b5061011f61021536600461072a565b610677565b6100cd61068b565b34801561022e57600080fd5b5061013e61023d3660046107c4565b600460209081526000928352604080842090915290825290205481565b3360009081526003602052604081208054349290610279908490610826565b909155505060405134815233907fe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c9060200160405180910390a2565b600080546102c290610839565b80601f01602080910402602001604051908101604052809291908181526020018280546102ee90610839565b801561033b5780601f106103105761010080835404028352916020019161033b565b820191906000526020600020905b81548152906001019060200180831161031e57829003601f168201915b505050505081565b33600081815260046020908152604080832073ffffffffffffffffffffffffffffffffffffffff8716808552925280832085905551919290917f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925906103ab9086815260200190565b60405180910390a35060015b92915050565b73ffffffffffffffffffffffffffffffffffffffff83166000908152600360205260408120548211156103ef57600080fd5b73ffffffffffffffffffffffffffffffffffffffff84163314801590610455575073ffffffffffffffffffffffffffffffffffffffff841660009081526004602090815260408083203384529091529020546fffffffffffffffffffffffffffffffff14155b156104dd5773ffffffffffffffffffffffffffffffffffffffff8416600090815260046020908152604080832033845290915290205482111561049757600080fd5b73ffffffffffffffffffffffffffffffffffffffff84166000908152600460209081526040808320338452909152812080548492906104d790849061088c565b90915550505b73ffffffffffffffffffffffffffffffffffffffff84166000908152600360205260408120805484929061051290849061088c565b909155505073ffffffffffffffffffffffffffffffffffffffff83166000908152600360205260408120805484929061054c908490610826565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040516105b291815260200190565b60405180910390a35060019392505050565b336000908152600360205260409020548111156105e057600080fd5b33600090815260036020526040812080548392906105ff90849061088c565b9091555050604051339082156108fc029083906000818181858888f19350505050158015610631573d6000803e3d6000fd5b5060405181815233907f7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b659060200160405180910390a250565b600180546102c290610839565b60006106843384846103bd565b9392505050565b61069361025a565b565b600060208083528351808285015260005b818110156106c2578581018301518582016040015282016106a6565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461072557600080fd5b919050565b6000806040838503121561073d57600080fd5b61074683610701565b946020939093013593505050565b60008060006060848603121561076957600080fd5b61077284610701565b925061078060208501610701565b9150604084013590509250925092565b6000602082840312156107a257600080fd5b5035919050565b6000602082840312156107bb57600080fd5b61068482610701565b600080604083850312156107d757600080fd5b6107e083610701565b91506107ee60208401610701565b90509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b808201808211156103b7576103b76107f7565b600181811c9082168061084d57607f821691505b602082108103610886577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b818103818111156103b7576103b76107f756fea164736f6c6343000813000a",
}

var WETH9ABI = WETH9MetaData.ABI

var WETH9Bin = WETH9MetaData.Bin

func DeployWETH9(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *WETH9, error) {
	parsed, err := WETH9MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(WETH9Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &WETH9{address: address, abi: *parsed, WETH9Caller: WETH9Caller{contract: contract}, WETH9Transactor: WETH9Transactor{contract: contract}, WETH9Filterer: WETH9Filterer{contract: contract}}, nil
}

type WETH9 struct {
	address common.Address
	abi     abi.ABI
	WETH9Caller
	WETH9Transactor
	WETH9Filterer
}

type WETH9Caller struct {
	contract *bind.BoundContract
}

type WETH9Transactor struct {
	contract *bind.BoundContract
}

type WETH9Filterer struct {
	contract *bind.BoundContract
}

type WETH9Session struct {
	Contract     *WETH9
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type WETH9CallerSession struct {
	Contract *WETH9Caller
	CallOpts bind.CallOpts
}

type WETH9TransactorSession struct {
	Contract     *WETH9Transactor
	TransactOpts bind.TransactOpts
}

type WETH9Raw struct {
	Contract *WETH9
}

type WETH9CallerRaw struct {
	Contract *WETH9Caller
}

type WETH9TransactorRaw struct {
	Contract *WETH9Transactor
}

func NewWETH9(address common.Address, backend bind.ContractBackend) (*WETH9, error) {
	abi, err := abi.JSON(strings.NewReader(WETH9ABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindWETH9(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &WETH9{address: address, abi: abi, WETH9Caller: WETH9Caller{contract: contract}, WETH9Transactor: WETH9Transactor{contract: contract}, WETH9Filterer: WETH9Filterer{contract: contract}}, nil
}

func NewWETH9Caller(address common.Address, caller bind.ContractCaller) (*WETH9Caller, error) {
	contract, err := bindWETH9(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &WETH9Caller{contract: contract}, nil
}

func NewWETH9Transactor(address common.Address, transactor bind.ContractTransactor) (*WETH9Transactor, error) {
	contract, err := bindWETH9(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &WETH9Transactor{contract: contract}, nil
}

func NewWETH9Filterer(address common.Address, filterer bind.ContractFilterer) (*WETH9Filterer, error) {
	contract, err := bindWETH9(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &WETH9Filterer{contract: contract}, nil
}

func bindWETH9(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := WETH9MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_WETH9 *WETH9Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WETH9.Contract.WETH9Caller.contract.Call(opts, result, method, params...)
}

func (_WETH9 *WETH9Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WETH9.Contract.WETH9Transactor.contract.Transfer(opts)
}

func (_WETH9 *WETH9Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WETH9.Contract.WETH9Transactor.contract.Transact(opts, method, params...)
}

func (_WETH9 *WETH9CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WETH9.Contract.contract.Call(opts, result, method, params...)
}

func (_WETH9 *WETH9TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WETH9.Contract.contract.Transfer(opts)
}

func (_WETH9 *WETH9TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WETH9.Contract.contract.Transact(opts, method, params...)
}

func (_WETH9 *WETH9Caller) Allowance(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _WETH9.contract.Call(opts, &out, "allowance", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_WETH9 *WETH9Session) Allowance(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _WETH9.Contract.Allowance(&_WETH9.CallOpts, arg0, arg1)
}

func (_WETH9 *WETH9CallerSession) Allowance(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _WETH9.Contract.Allowance(&_WETH9.CallOpts, arg0, arg1)
}

func (_WETH9 *WETH9Caller) BalanceOf(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _WETH9.contract.Call(opts, &out, "balanceOf", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_WETH9 *WETH9Session) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _WETH9.Contract.BalanceOf(&_WETH9.CallOpts, arg0)
}

func (_WETH9 *WETH9CallerSession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _WETH9.Contract.BalanceOf(&_WETH9.CallOpts, arg0)
}

func (_WETH9 *WETH9Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _WETH9.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_WETH9 *WETH9Session) Decimals() (uint8, error) {
	return _WETH9.Contract.Decimals(&_WETH9.CallOpts)
}

func (_WETH9 *WETH9CallerSession) Decimals() (uint8, error) {
	return _WETH9.Contract.Decimals(&_WETH9.CallOpts)
}

func (_WETH9 *WETH9Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _WETH9.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_WETH9 *WETH9Session) Name() (string, error) {
	return _WETH9.Contract.Name(&_WETH9.CallOpts)
}

func (_WETH9 *WETH9CallerSession) Name() (string, error) {
	return _WETH9.Contract.Name(&_WETH9.CallOpts)
}

func (_WETH9 *WETH9Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _WETH9.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_WETH9 *WETH9Session) Symbol() (string, error) {
	return _WETH9.Contract.Symbol(&_WETH9.CallOpts)
}

func (_WETH9 *WETH9CallerSession) Symbol() (string, error) {
	return _WETH9.Contract.Symbol(&_WETH9.CallOpts)
}

func (_WETH9 *WETH9Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _WETH9.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_WETH9 *WETH9Session) TotalSupply() (*big.Int, error) {
	return _WETH9.Contract.TotalSupply(&_WETH9.CallOpts)
}

func (_WETH9 *WETH9CallerSession) TotalSupply() (*big.Int, error) {
	return _WETH9.Contract.TotalSupply(&_WETH9.CallOpts)
}

func (_WETH9 *WETH9Transactor) Approve(opts *bind.TransactOpts, guy common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WETH9.contract.Transact(opts, "approve", guy, wad)
}

func (_WETH9 *WETH9Session) Approve(guy common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WETH9.Contract.Approve(&_WETH9.TransactOpts, guy, wad)
}

func (_WETH9 *WETH9TransactorSession) Approve(guy common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WETH9.Contract.Approve(&_WETH9.TransactOpts, guy, wad)
}

func (_WETH9 *WETH9Transactor) Deposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WETH9.contract.Transact(opts, "deposit")
}

func (_WETH9 *WETH9Session) Deposit() (*types.Transaction, error) {
	return _WETH9.Contract.Deposit(&_WETH9.TransactOpts)
}

func (_WETH9 *WETH9TransactorSession) Deposit() (*types.Transaction, error) {
	return _WETH9.Contract.Deposit(&_WETH9.TransactOpts)
}

func (_WETH9 *WETH9Transactor) Transfer(opts *bind.TransactOpts, dst common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WETH9.contract.Transact(opts, "transfer", dst, wad)
}

func (_WETH9 *WETH9Session) Transfer(dst common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WETH9.Contract.Transfer(&_WETH9.TransactOpts, dst, wad)
}

func (_WETH9 *WETH9TransactorSession) Transfer(dst common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WETH9.Contract.Transfer(&_WETH9.TransactOpts, dst, wad)
}

func (_WETH9 *WETH9Transactor) TransferFrom(opts *bind.TransactOpts, src common.Address, dst common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WETH9.contract.Transact(opts, "transferFrom", src, dst, wad)
}

func (_WETH9 *WETH9Session) TransferFrom(src common.Address, dst common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WETH9.Contract.TransferFrom(&_WETH9.TransactOpts, src, dst, wad)
}

func (_WETH9 *WETH9TransactorSession) TransferFrom(src common.Address, dst common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WETH9.Contract.TransferFrom(&_WETH9.TransactOpts, src, dst, wad)
}

func (_WETH9 *WETH9Transactor) Withdraw(opts *bind.TransactOpts, wad *big.Int) (*types.Transaction, error) {
	return _WETH9.contract.Transact(opts, "withdraw", wad)
}

func (_WETH9 *WETH9Session) Withdraw(wad *big.Int) (*types.Transaction, error) {
	return _WETH9.Contract.Withdraw(&_WETH9.TransactOpts, wad)
}

func (_WETH9 *WETH9TransactorSession) Withdraw(wad *big.Int) (*types.Transaction, error) {
	return _WETH9.Contract.Withdraw(&_WETH9.TransactOpts, wad)
}

func (_WETH9 *WETH9Transactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WETH9.contract.RawTransact(opts, nil)
}

func (_WETH9 *WETH9Session) Receive() (*types.Transaction, error) {
	return _WETH9.Contract.Receive(&_WETH9.TransactOpts)
}

func (_WETH9 *WETH9TransactorSession) Receive() (*types.Transaction, error) {
	return _WETH9.Contract.Receive(&_WETH9.TransactOpts)
}

type WETH9ApprovalIterator struct {
	Event *WETH9Approval

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WETH9ApprovalIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WETH9Approval)
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
		it.Event = new(WETH9Approval)
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

func (it *WETH9ApprovalIterator) Error() error {
	return it.fail
}

func (it *WETH9ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WETH9Approval struct {
	Src common.Address
	Guy common.Address
	Wad *big.Int
	Raw types.Log
}

func (_WETH9 *WETH9Filterer) FilterApproval(opts *bind.FilterOpts, src []common.Address, guy []common.Address) (*WETH9ApprovalIterator, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}
	var guyRule []interface{}
	for _, guyItem := range guy {
		guyRule = append(guyRule, guyItem)
	}

	logs, sub, err := _WETH9.contract.FilterLogs(opts, "Approval", srcRule, guyRule)
	if err != nil {
		return nil, err
	}
	return &WETH9ApprovalIterator{contract: _WETH9.contract, event: "Approval", logs: logs, sub: sub}, nil
}

func (_WETH9 *WETH9Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *WETH9Approval, src []common.Address, guy []common.Address) (event.Subscription, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}
	var guyRule []interface{}
	for _, guyItem := range guy {
		guyRule = append(guyRule, guyItem)
	}

	logs, sub, err := _WETH9.contract.WatchLogs(opts, "Approval", srcRule, guyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WETH9Approval)
				if err := _WETH9.contract.UnpackLog(event, "Approval", log); err != nil {
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

func (_WETH9 *WETH9Filterer) ParseApproval(log types.Log) (*WETH9Approval, error) {
	event := new(WETH9Approval)
	if err := _WETH9.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type WETH9DepositIterator struct {
	Event *WETH9Deposit

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WETH9DepositIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WETH9Deposit)
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
		it.Event = new(WETH9Deposit)
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

func (it *WETH9DepositIterator) Error() error {
	return it.fail
}

func (it *WETH9DepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WETH9Deposit struct {
	Dst common.Address
	Wad *big.Int
	Raw types.Log
}

func (_WETH9 *WETH9Filterer) FilterDeposit(opts *bind.FilterOpts, dst []common.Address) (*WETH9DepositIterator, error) {

	var dstRule []interface{}
	for _, dstItem := range dst {
		dstRule = append(dstRule, dstItem)
	}

	logs, sub, err := _WETH9.contract.FilterLogs(opts, "Deposit", dstRule)
	if err != nil {
		return nil, err
	}
	return &WETH9DepositIterator{contract: _WETH9.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

func (_WETH9 *WETH9Filterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *WETH9Deposit, dst []common.Address) (event.Subscription, error) {

	var dstRule []interface{}
	for _, dstItem := range dst {
		dstRule = append(dstRule, dstItem)
	}

	logs, sub, err := _WETH9.contract.WatchLogs(opts, "Deposit", dstRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WETH9Deposit)
				if err := _WETH9.contract.UnpackLog(event, "Deposit", log); err != nil {
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

func (_WETH9 *WETH9Filterer) ParseDeposit(log types.Log) (*WETH9Deposit, error) {
	event := new(WETH9Deposit)
	if err := _WETH9.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type WETH9TransferIterator struct {
	Event *WETH9Transfer

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WETH9TransferIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WETH9Transfer)
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
		it.Event = new(WETH9Transfer)
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

func (it *WETH9TransferIterator) Error() error {
	return it.fail
}

func (it *WETH9TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WETH9Transfer struct {
	Src common.Address
	Dst common.Address
	Wad *big.Int
	Raw types.Log
}

func (_WETH9 *WETH9Filterer) FilterTransfer(opts *bind.FilterOpts, src []common.Address, dst []common.Address) (*WETH9TransferIterator, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}
	var dstRule []interface{}
	for _, dstItem := range dst {
		dstRule = append(dstRule, dstItem)
	}

	logs, sub, err := _WETH9.contract.FilterLogs(opts, "Transfer", srcRule, dstRule)
	if err != nil {
		return nil, err
	}
	return &WETH9TransferIterator{contract: _WETH9.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

func (_WETH9 *WETH9Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *WETH9Transfer, src []common.Address, dst []common.Address) (event.Subscription, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}
	var dstRule []interface{}
	for _, dstItem := range dst {
		dstRule = append(dstRule, dstItem)
	}

	logs, sub, err := _WETH9.contract.WatchLogs(opts, "Transfer", srcRule, dstRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WETH9Transfer)
				if err := _WETH9.contract.UnpackLog(event, "Transfer", log); err != nil {
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

func (_WETH9 *WETH9Filterer) ParseTransfer(log types.Log) (*WETH9Transfer, error) {
	event := new(WETH9Transfer)
	if err := _WETH9.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type WETH9WithdrawalIterator struct {
	Event *WETH9Withdrawal

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WETH9WithdrawalIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WETH9Withdrawal)
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
		it.Event = new(WETH9Withdrawal)
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

func (it *WETH9WithdrawalIterator) Error() error {
	return it.fail
}

func (it *WETH9WithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WETH9Withdrawal struct {
	Src common.Address
	Wad *big.Int
	Raw types.Log
}

func (_WETH9 *WETH9Filterer) FilterWithdrawal(opts *bind.FilterOpts, src []common.Address) (*WETH9WithdrawalIterator, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}

	logs, sub, err := _WETH9.contract.FilterLogs(opts, "Withdrawal", srcRule)
	if err != nil {
		return nil, err
	}
	return &WETH9WithdrawalIterator{contract: _WETH9.contract, event: "Withdrawal", logs: logs, sub: sub}, nil
}

func (_WETH9 *WETH9Filterer) WatchWithdrawal(opts *bind.WatchOpts, sink chan<- *WETH9Withdrawal, src []common.Address) (event.Subscription, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}

	logs, sub, err := _WETH9.contract.WatchLogs(opts, "Withdrawal", srcRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WETH9Withdrawal)
				if err := _WETH9.contract.UnpackLog(event, "Withdrawal", log); err != nil {
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

func (_WETH9 *WETH9Filterer) ParseWithdrawal(log types.Log) (*WETH9Withdrawal, error) {
	event := new(WETH9Withdrawal)
	if err := _WETH9.contract.UnpackLog(event, "Withdrawal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_WETH9 *WETH9) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _WETH9.abi.Events["Approval"].ID:
		return _WETH9.ParseApproval(log)
	case _WETH9.abi.Events["Deposit"].ID:
		return _WETH9.ParseDeposit(log)
	case _WETH9.abi.Events["Transfer"].ID:
		return _WETH9.ParseTransfer(log)
	case _WETH9.abi.Events["Withdrawal"].ID:
		return _WETH9.ParseWithdrawal(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (WETH9Approval) Topic() common.Hash {
	return common.HexToHash("0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925")
}

func (WETH9Deposit) Topic() common.Hash {
	return common.HexToHash("0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c")
}

func (WETH9Transfer) Topic() common.Hash {
	return common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")
}

func (WETH9Withdrawal) Topic() common.Hash {
	return common.HexToHash("0x7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65")
}

func (_WETH9 *WETH9) Address() common.Address {
	return _WETH9.address
}

type WETH9Interface interface {
	Allowance(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error)

	BalanceOf(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error)

	Decimals(opts *bind.CallOpts) (uint8, error)

	Name(opts *bind.CallOpts) (string, error)

	Symbol(opts *bind.CallOpts) (string, error)

	TotalSupply(opts *bind.CallOpts) (*big.Int, error)

	Approve(opts *bind.TransactOpts, guy common.Address, wad *big.Int) (*types.Transaction, error)

	Deposit(opts *bind.TransactOpts) (*types.Transaction, error)

	Transfer(opts *bind.TransactOpts, dst common.Address, wad *big.Int) (*types.Transaction, error)

	TransferFrom(opts *bind.TransactOpts, src common.Address, dst common.Address, wad *big.Int) (*types.Transaction, error)

	Withdraw(opts *bind.TransactOpts, wad *big.Int) (*types.Transaction, error)

	Receive(opts *bind.TransactOpts) (*types.Transaction, error)

	FilterApproval(opts *bind.FilterOpts, src []common.Address, guy []common.Address) (*WETH9ApprovalIterator, error)

	WatchApproval(opts *bind.WatchOpts, sink chan<- *WETH9Approval, src []common.Address, guy []common.Address) (event.Subscription, error)

	ParseApproval(log types.Log) (*WETH9Approval, error)

	FilterDeposit(opts *bind.FilterOpts, dst []common.Address) (*WETH9DepositIterator, error)

	WatchDeposit(opts *bind.WatchOpts, sink chan<- *WETH9Deposit, dst []common.Address) (event.Subscription, error)

	ParseDeposit(log types.Log) (*WETH9Deposit, error)

	FilterTransfer(opts *bind.FilterOpts, src []common.Address, dst []common.Address) (*WETH9TransferIterator, error)

	WatchTransfer(opts *bind.WatchOpts, sink chan<- *WETH9Transfer, src []common.Address, dst []common.Address) (event.Subscription, error)

	ParseTransfer(log types.Log) (*WETH9Transfer, error)

	FilterWithdrawal(opts *bind.FilterOpts, src []common.Address) (*WETH9WithdrawalIterator, error)

	WatchWithdrawal(opts *bind.WatchOpts, sink chan<- *WETH9Withdrawal, src []common.Address) (event.Subscription, error)

	ParseWithdrawal(log types.Log) (*WETH9Withdrawal, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
