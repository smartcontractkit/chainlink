// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package erc677

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

var ERC677MetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"allowance\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decreaseAllowance\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"subtractedValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"increaseAllowance\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"addedValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transfer\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferAndCall\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"success\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false}]",
	Bin: "0x60806040523480156200001157600080fd5b50604051620010ae380380620010ae833981016040819052620000349162000123565b818160036200004483826200021c565b5060046200005382826200021c565b5050505050620002e8565b634e487b7160e01b600052604160045260246000fd5b600082601f8301126200008657600080fd5b81516001600160401b0380821115620000a357620000a36200005e565b604051601f8301601f19908116603f01168101908282118183101715620000ce57620000ce6200005e565b81604052838152602092508683858801011115620000eb57600080fd5b600091505b838210156200010f5785820183015181830184015290820190620000f0565b600093810190920192909252949350505050565b600080604083850312156200013757600080fd5b82516001600160401b03808211156200014f57600080fd5b6200015d8683870162000074565b935060208501519150808211156200017457600080fd5b50620001838582860162000074565b9150509250929050565b600181811c90821680620001a257607f821691505b602082108103620001c357634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200021757600081815260208120601f850160051c81016020861015620001f25750805b601f850160051c820191505b818110156200021357828155600101620001fe565b5050505b505050565b81516001600160401b038111156200023857620002386200005e565b62000250816200024984546200018d565b84620001c9565b602080601f8311600181146200028857600084156200026f5750858301515b600019600386901b1c1916600185901b17855562000213565b600085815260208120601f198616915b82811015620002b95788860151825594840194600190910190840162000298565b5085821015620002d85787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b610db680620002f86000396000f3fe608060405234801561001057600080fd5b50600436106100d45760003560e01c80634000aea011610081578063a457c2d71161005b578063a457c2d7146101b2578063a9059cbb146101c5578063dd62ed3e146101d857600080fd5b80634000aea01461016157806370a082311461017457806395d89b41146101aa57600080fd5b806323b872dd116100b257806323b872dd1461012c578063313ce5671461013f578063395093511461014e57600080fd5b806306fdde03146100d9578063095ea7b3146100f757806318160ddd1461011a575b600080fd5b6100e161021e565b6040516100ee9190610aae565b60405180910390f35b61010a610105366004610af1565b6102b0565b60405190151581526020016100ee565b6002545b6040519081526020016100ee565b61010a61013a366004610b1b565b6102ca565b604051601281526020016100ee565b61010a61015c366004610af1565b6102ee565b61010a61016f366004610b86565b61033a565b61011e610182366004610c6f565b73ffffffffffffffffffffffffffffffffffffffff1660009081526020819052604090205490565b6100e161045e565b61010a6101c0366004610af1565b61046d565b61010a6101d3366004610af1565b610543565b61011e6101e6366004610c8a565b73ffffffffffffffffffffffffffffffffffffffff918216600090815260016020908152604080832093909416825291909152205490565b60606003805461022d90610cbd565b80601f016020809104026020016040519081016040528092919081815260200182805461025990610cbd565b80156102a65780601f1061027b576101008083540402835291602001916102a6565b820191906000526020600020905b81548152906001019060200180831161028957829003601f168201915b5050505050905090565b6000336102be818585610551565b60019150505b92915050565b6000336102d8858285610704565b6102e38585856107db565b506001949350505050565b33600081815260016020908152604080832073ffffffffffffffffffffffffffffffffffffffff871684529091528120549091906102be9082908690610335908790610d10565b610551565b60006103468484610543565b508373ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167fe19260aff97b920c7df27010903aeb9c8d2be5d310a2c67824cf3f15396e4c1685856040516103a6929190610d4a565b60405180910390a373ffffffffffffffffffffffffffffffffffffffff84163b15610454576040517fa4c0ed3600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff85169063a4c0ed369061042190339087908790600401610d6b565b600060405180830381600087803b15801561043b57600080fd5b505af115801561044f573d6000803e3d6000fd5b505050505b5060019392505050565b60606004805461022d90610cbd565b33600081815260016020908152604080832073ffffffffffffffffffffffffffffffffffffffff8716845290915281205490919083811015610536576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602560248201527f45524332303a2064656372656173656420616c6c6f77616e63652062656c6f7760448201527f207a65726f00000000000000000000000000000000000000000000000000000060648201526084015b60405180910390fd5b6102e38286868403610551565b6000336102be8185856107db565b73ffffffffffffffffffffffffffffffffffffffff83166105f3576040517f08c379a0000000000000000000000000000000000000000000000000000000008152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f2061646460448201527f7265737300000000000000000000000000000000000000000000000000000000606482015260840161052d565b73ffffffffffffffffffffffffffffffffffffffff8216610696576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f20616464726560448201527f7373000000000000000000000000000000000000000000000000000000000000606482015260840161052d565b73ffffffffffffffffffffffffffffffffffffffff83811660008181526001602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b73ffffffffffffffffffffffffffffffffffffffff8381166000908152600160209081526040808320938616835292905220547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81146107d557818110156107c8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f45524332303a20696e73756666696369656e7420616c6c6f77616e6365000000604482015260640161052d565b6107d58484848403610551565b50505050565b73ffffffffffffffffffffffffffffffffffffffff831661087e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602560248201527f45524332303a207472616e736665722066726f6d20746865207a65726f20616460448201527f6472657373000000000000000000000000000000000000000000000000000000606482015260840161052d565b73ffffffffffffffffffffffffffffffffffffffff8216610921576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602360248201527f45524332303a207472616e7366657220746f20746865207a65726f206164647260448201527f6573730000000000000000000000000000000000000000000000000000000000606482015260840161052d565b73ffffffffffffffffffffffffffffffffffffffff8316600090815260208190526040902054818110156109d7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f45524332303a207472616e7366657220616d6f756e742065786365656473206260448201527f616c616e63650000000000000000000000000000000000000000000000000000606482015260840161052d565b73ffffffffffffffffffffffffffffffffffffffff848116600081815260208181526040808320878703905593871680835291849020805487019055925185815290927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef910160405180910390a36107d5565b6000815180845260005b81811015610a7057602081850181015186830182015201610a54565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b602081526000610ac16020830184610a4a565b9392505050565b803573ffffffffffffffffffffffffffffffffffffffff81168114610aec57600080fd5b919050565b60008060408385031215610b0457600080fd5b610b0d83610ac8565b946020939093013593505050565b600080600060608486031215610b3057600080fd5b610b3984610ac8565b9250610b4760208501610ac8565b9150604084013590509250925092565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600080600060608486031215610b9b57600080fd5b610ba484610ac8565b925060208401359150604084013567ffffffffffffffff80821115610bc857600080fd5b818601915086601f830112610bdc57600080fd5b813581811115610bee57610bee610b57565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f01168101908382118183101715610c3457610c34610b57565b81604052828152896020848701011115610c4d57600080fd5b8260208601602083013760006020848301015280955050505050509250925092565b600060208284031215610c8157600080fd5b610ac182610ac8565b60008060408385031215610c9d57600080fd5b610ca683610ac8565b9150610cb460208401610ac8565b90509250929050565b600181811c90821680610cd157607f821691505b602082108103610d0a577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b808201808211156102c4577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b828152604060208201526000610d636040830184610a4a565b949350505050565b73ffffffffffffffffffffffffffffffffffffffff84168152826020820152606060408201526000610da06060830184610a4a565b9594505050505056fea164736f6c6343000813000a",
}

var ERC677ABI = ERC677MetaData.ABI

var ERC677Bin = ERC677MetaData.Bin

func DeployERC677(auth *bind.TransactOpts, backend bind.ContractBackend, name string, symbol string) (common.Address, *types.Transaction, *ERC677, error) {
	parsed, err := ERC677MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ERC677Bin), backend, name, symbol)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ERC677{address: address, abi: *parsed, ERC677Caller: ERC677Caller{contract: contract}, ERC677Transactor: ERC677Transactor{contract: contract}, ERC677Filterer: ERC677Filterer{contract: contract}}, nil
}

type ERC677 struct {
	address common.Address
	abi     abi.ABI
	ERC677Caller
	ERC677Transactor
	ERC677Filterer
}

type ERC677Caller struct {
	contract *bind.BoundContract
}

type ERC677Transactor struct {
	contract *bind.BoundContract
}

type ERC677Filterer struct {
	contract *bind.BoundContract
}

type ERC677Session struct {
	Contract     *ERC677
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type ERC677CallerSession struct {
	Contract *ERC677Caller
	CallOpts bind.CallOpts
}

type ERC677TransactorSession struct {
	Contract     *ERC677Transactor
	TransactOpts bind.TransactOpts
}

type ERC677Raw struct {
	Contract *ERC677
}

type ERC677CallerRaw struct {
	Contract *ERC677Caller
}

type ERC677TransactorRaw struct {
	Contract *ERC677Transactor
}

func NewERC677(address common.Address, backend bind.ContractBackend) (*ERC677, error) {
	abi, err := abi.JSON(strings.NewReader(ERC677ABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindERC677(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC677{address: address, abi: abi, ERC677Caller: ERC677Caller{contract: contract}, ERC677Transactor: ERC677Transactor{contract: contract}, ERC677Filterer: ERC677Filterer{contract: contract}}, nil
}

func NewERC677Caller(address common.Address, caller bind.ContractCaller) (*ERC677Caller, error) {
	contract, err := bindERC677(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC677Caller{contract: contract}, nil
}

func NewERC677Transactor(address common.Address, transactor bind.ContractTransactor) (*ERC677Transactor, error) {
	contract, err := bindERC677(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC677Transactor{contract: contract}, nil
}

func NewERC677Filterer(address common.Address, filterer bind.ContractFilterer) (*ERC677Filterer, error) {
	contract, err := bindERC677(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC677Filterer{contract: contract}, nil
}

func bindERC677(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ERC677MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_ERC677 *ERC677Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC677.Contract.ERC677Caller.contract.Call(opts, result, method, params...)
}

func (_ERC677 *ERC677Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC677.Contract.ERC677Transactor.contract.Transfer(opts)
}

func (_ERC677 *ERC677Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC677.Contract.ERC677Transactor.contract.Transact(opts, method, params...)
}

func (_ERC677 *ERC677CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC677.Contract.contract.Call(opts, result, method, params...)
}

func (_ERC677 *ERC677TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC677.Contract.contract.Transfer(opts)
}

func (_ERC677 *ERC677TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC677.Contract.contract.Transact(opts, method, params...)
}

func (_ERC677 *ERC677Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ERC677.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_ERC677 *ERC677Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ERC677.Contract.Allowance(&_ERC677.CallOpts, owner, spender)
}

func (_ERC677 *ERC677CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ERC677.Contract.Allowance(&_ERC677.CallOpts, owner, spender)
}

func (_ERC677 *ERC677Caller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ERC677.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_ERC677 *ERC677Session) BalanceOf(account common.Address) (*big.Int, error) {
	return _ERC677.Contract.BalanceOf(&_ERC677.CallOpts, account)
}

func (_ERC677 *ERC677CallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ERC677.Contract.BalanceOf(&_ERC677.CallOpts, account)
}

func (_ERC677 *ERC677Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _ERC677.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_ERC677 *ERC677Session) Decimals() (uint8, error) {
	return _ERC677.Contract.Decimals(&_ERC677.CallOpts)
}

func (_ERC677 *ERC677CallerSession) Decimals() (uint8, error) {
	return _ERC677.Contract.Decimals(&_ERC677.CallOpts)
}

func (_ERC677 *ERC677Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ERC677.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_ERC677 *ERC677Session) Name() (string, error) {
	return _ERC677.Contract.Name(&_ERC677.CallOpts)
}

func (_ERC677 *ERC677CallerSession) Name() (string, error) {
	return _ERC677.Contract.Name(&_ERC677.CallOpts)
}

func (_ERC677 *ERC677Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ERC677.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_ERC677 *ERC677Session) Symbol() (string, error) {
	return _ERC677.Contract.Symbol(&_ERC677.CallOpts)
}

func (_ERC677 *ERC677CallerSession) Symbol() (string, error) {
	return _ERC677.Contract.Symbol(&_ERC677.CallOpts)
}

func (_ERC677 *ERC677Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ERC677.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_ERC677 *ERC677Session) TotalSupply() (*big.Int, error) {
	return _ERC677.Contract.TotalSupply(&_ERC677.CallOpts)
}

func (_ERC677 *ERC677CallerSession) TotalSupply() (*big.Int, error) {
	return _ERC677.Contract.TotalSupply(&_ERC677.CallOpts)
}

func (_ERC677 *ERC677Transactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC677.contract.Transact(opts, "approve", spender, amount)
}

func (_ERC677 *ERC677Session) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC677.Contract.Approve(&_ERC677.TransactOpts, spender, amount)
}

func (_ERC677 *ERC677TransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC677.Contract.Approve(&_ERC677.TransactOpts, spender, amount)
}

func (_ERC677 *ERC677Transactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ERC677.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

func (_ERC677 *ERC677Session) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ERC677.Contract.DecreaseAllowance(&_ERC677.TransactOpts, spender, subtractedValue)
}

func (_ERC677 *ERC677TransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ERC677.Contract.DecreaseAllowance(&_ERC677.TransactOpts, spender, subtractedValue)
}

func (_ERC677 *ERC677Transactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ERC677.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

func (_ERC677 *ERC677Session) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ERC677.Contract.IncreaseAllowance(&_ERC677.TransactOpts, spender, addedValue)
}

func (_ERC677 *ERC677TransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ERC677.Contract.IncreaseAllowance(&_ERC677.TransactOpts, spender, addedValue)
}

func (_ERC677 *ERC677Transactor) Transfer(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC677.contract.Transact(opts, "transfer", to, amount)
}

func (_ERC677 *ERC677Session) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC677.Contract.Transfer(&_ERC677.TransactOpts, to, amount)
}

func (_ERC677 *ERC677TransactorSession) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC677.Contract.Transfer(&_ERC677.TransactOpts, to, amount)
}

func (_ERC677 *ERC677Transactor) TransferAndCall(opts *bind.TransactOpts, to common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _ERC677.contract.Transact(opts, "transferAndCall", to, amount, data)
}

func (_ERC677 *ERC677Session) TransferAndCall(to common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _ERC677.Contract.TransferAndCall(&_ERC677.TransactOpts, to, amount, data)
}

func (_ERC677 *ERC677TransactorSession) TransferAndCall(to common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _ERC677.Contract.TransferAndCall(&_ERC677.TransactOpts, to, amount, data)
}

func (_ERC677 *ERC677Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC677.contract.Transact(opts, "transferFrom", from, to, amount)
}

func (_ERC677 *ERC677Session) TransferFrom(from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC677.Contract.TransferFrom(&_ERC677.TransactOpts, from, to, amount)
}

func (_ERC677 *ERC677TransactorSession) TransferFrom(from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC677.Contract.TransferFrom(&_ERC677.TransactOpts, from, to, amount)
}

type ERC677ApprovalIterator struct {
	Event *ERC677Approval

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ERC677ApprovalIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC677Approval)
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
		it.Event = new(ERC677Approval)
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

func (it *ERC677ApprovalIterator) Error() error {
	return it.fail
}

func (it *ERC677ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ERC677Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log
}

func (_ERC677 *ERC677Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*ERC677ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ERC677.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &ERC677ApprovalIterator{contract: _ERC677.contract, event: "Approval", logs: logs, sub: sub}, nil
}

func (_ERC677 *ERC677Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ERC677Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ERC677.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ERC677Approval)
				if err := _ERC677.contract.UnpackLog(event, "Approval", log); err != nil {
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

func (_ERC677 *ERC677Filterer) ParseApproval(log types.Log) (*ERC677Approval, error) {
	event := new(ERC677Approval)
	if err := _ERC677.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ERC677TransferIterator struct {
	Event *ERC677Transfer

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ERC677TransferIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC677Transfer)
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
		it.Event = new(ERC677Transfer)
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

func (it *ERC677TransferIterator) Error() error {
	return it.fail
}

func (it *ERC677TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ERC677Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log
}

func (_ERC677 *ERC677Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ERC677TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ERC677.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ERC677TransferIterator{contract: _ERC677.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

func (_ERC677 *ERC677Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ERC677Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ERC677.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ERC677Transfer)
				if err := _ERC677.contract.UnpackLog(event, "Transfer", log); err != nil {
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

func (_ERC677 *ERC677Filterer) ParseTransfer(log types.Log) (*ERC677Transfer, error) {
	event := new(ERC677Transfer)
	if err := _ERC677.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ERC677Transfer0Iterator struct {
	Event *ERC677Transfer0

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ERC677Transfer0Iterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC677Transfer0)
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
		it.Event = new(ERC677Transfer0)
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

func (it *ERC677Transfer0Iterator) Error() error {
	return it.fail
}

func (it *ERC677Transfer0Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ERC677Transfer0 struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Data  []byte
	Raw   types.Log
}

func (_ERC677 *ERC677Filterer) FilterTransfer0(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ERC677Transfer0Iterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ERC677.contract.FilterLogs(opts, "Transfer0", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ERC677Transfer0Iterator{contract: _ERC677.contract, event: "Transfer0", logs: logs, sub: sub}, nil
}

func (_ERC677 *ERC677Filterer) WatchTransfer0(opts *bind.WatchOpts, sink chan<- *ERC677Transfer0, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ERC677.contract.WatchLogs(opts, "Transfer0", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ERC677Transfer0)
				if err := _ERC677.contract.UnpackLog(event, "Transfer0", log); err != nil {
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

func (_ERC677 *ERC677Filterer) ParseTransfer0(log types.Log) (*ERC677Transfer0, error) {
	event := new(ERC677Transfer0)
	if err := _ERC677.contract.UnpackLog(event, "Transfer0", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_ERC677 *ERC677) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _ERC677.abi.Events["Approval"].ID:
		return _ERC677.ParseApproval(log)
	case _ERC677.abi.Events["Transfer"].ID:
		return _ERC677.ParseTransfer(log)
	case _ERC677.abi.Events["Transfer0"].ID:
		return _ERC677.ParseTransfer0(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (ERC677Approval) Topic() common.Hash {
	return common.HexToHash("0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925")
}

func (ERC677Transfer) Topic() common.Hash {
	return common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")
}

func (ERC677Transfer0) Topic() common.Hash {
	return common.HexToHash("0xe19260aff97b920c7df27010903aeb9c8d2be5d310a2c67824cf3f15396e4c16")
}

func (_ERC677 *ERC677) Address() common.Address {
	return _ERC677.address
}

type ERC677Interface interface {
	Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error)

	BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error)

	Decimals(opts *bind.CallOpts) (uint8, error)

	Name(opts *bind.CallOpts) (string, error)

	Symbol(opts *bind.CallOpts) (string, error)

	TotalSupply(opts *bind.CallOpts) (*big.Int, error)

	Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error)

	DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error)

	IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error)

	Transfer(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error)

	TransferAndCall(opts *bind.TransactOpts, to common.Address, amount *big.Int, data []byte) (*types.Transaction, error)

	TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error)

	FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*ERC677ApprovalIterator, error)

	WatchApproval(opts *bind.WatchOpts, sink chan<- *ERC677Approval, owner []common.Address, spender []common.Address) (event.Subscription, error)

	ParseApproval(log types.Log) (*ERC677Approval, error)

	FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ERC677TransferIterator, error)

	WatchTransfer(opts *bind.WatchOpts, sink chan<- *ERC677Transfer, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseTransfer(log types.Log) (*ERC677Transfer, error)

	FilterTransfer0(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ERC677Transfer0Iterator, error)

	WatchTransfer0(opts *bind.WatchOpts, sink chan<- *ERC677Transfer0, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseTransfer0(log types.Log) (*ERC677Transfer0, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
