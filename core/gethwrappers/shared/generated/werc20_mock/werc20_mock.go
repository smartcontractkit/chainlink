// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package werc20_mock

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

var WERC20MockMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"allowance\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"burn\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decreaseAllowance\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"subtractedValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"increaseAllowance\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"addedValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"mint\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transfer\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"wad\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Deposit\",\"inputs\":[{\"name\":\"dst\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"wad\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdrawal\",\"inputs\":[{\"name\":\"src\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"wad\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false}]",
	Bin: "0x60806040523480156200001157600080fd5b506040518060400160405280600a8152602001695745524332304d6f636b60b01b815250604051806040016040528060048152602001635745524360e01b815250816003908162000063919062000120565b50600462000072828262000120565b505050620001ec565b634e487b7160e01b600052604160045260246000fd5b600181811c90821680620000a657607f821691505b602082108103620000c757634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200011b57600081815260208120601f850160051c81016020861015620000f65750805b601f850160051c820191505b81811015620001175782815560010162000102565b5050505b505050565b81516001600160401b038111156200013c576200013c6200007b565b62000154816200014d845462000091565b84620000cd565b602080601f8311600181146200018c5760008415620001735750858301515b600019600386901b1c1916600185901b17855562000117565b600085815260208120601f198616915b82811015620001bd578886015182559484019460019091019084016200019c565b5085821015620001dc5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b610fc980620001fc6000396000f3fe6080604052600436106100ec5760003560e01c806340c10f191161008a578063a457c2d711610059578063a457c2d71461028e578063a9059cbb146102ae578063d0e30db0146102ce578063dd62ed3e146102d657600080fd5b806340c10f19146101f657806370a082311461021657806395d89b41146102595780639dc29fac1461026e57600080fd5b806323b872dd116100c657806323b872dd1461017a5780632e1a7d4d1461019a578063313ce567146101ba57806339509351146101d657600080fd5b806306fdde0314610100578063095ea7b31461012b57806318160ddd1461015b57600080fd5b366100fb576100f9610329565b005b600080fd5b34801561010c57600080fd5b5061011561036a565b6040516101229190610dc6565b60405180910390f35b34801561013757600080fd5b5061014b610146366004610e5b565b6103fc565b6040519015158152602001610122565b34801561016757600080fd5b506002545b604051908152602001610122565b34801561018657600080fd5b5061014b610195366004610e85565b610416565b3480156101a657600080fd5b506100f96101b5366004610ec1565b61043a565b3480156101c657600080fd5b5060405160128152602001610122565b3480156101e257600080fd5b5061014b6101f1366004610e5b565b6104c6565b34801561020257600080fd5b506100f9610211366004610e5b565b610512565b34801561022257600080fd5b5061016c610231366004610eda565b73ffffffffffffffffffffffffffffffffffffffff1660009081526020819052604090205490565b34801561026557600080fd5b50610115610520565b34801561027a57600080fd5b506100f9610289366004610e5b565b61052f565b34801561029a57600080fd5b5061014b6102a9366004610e5b565b610539565b3480156102ba57600080fd5b5061014b6102c9366004610e5b565b61060f565b6100f9610329565b3480156102e257600080fd5b5061016c6102f1366004610efc565b73ffffffffffffffffffffffffffffffffffffffff918216600090815260016020908152604080832093909416825291909152205490565b610333333461061d565b60405134815233907fe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c9060200160405180910390a2565b60606003805461037990610f2f565b80601f01602080910402602001604051908101604052809291908181526020018280546103a590610f2f565b80156103f25780601f106103c7576101008083540402835291602001916103f2565b820191906000526020600020905b8154815290600101906020018083116103d557829003601f168201915b5050505050905090565b60003361040a818585610710565b60019150505b92915050565b6000336104248582856108c4565b61042f85858561099b565b506001949350505050565b3360009081526020819052604090205481111561045657600080fd5b6104603382610c0a565b604051339082156108fc029083906000818181858888f1935050505015801561048d573d6000803e3d6000fd5b5060405181815233907f7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b659060200160405180910390a250565b33600081815260016020908152604080832073ffffffffffffffffffffffffffffffffffffffff8716845290915281205490919061040a908290869061050d908790610f82565b610710565b61051c828261061d565b5050565b60606004805461037990610f2f565b61051c8282610c0a565b33600081815260016020908152604080832073ffffffffffffffffffffffffffffffffffffffff8716845290915281205490919083811015610602576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602560248201527f45524332303a2064656372656173656420616c6c6f77616e63652062656c6f7760448201527f207a65726f00000000000000000000000000000000000000000000000000000060648201526084015b60405180910390fd5b61042f8286868403610710565b60003361040a81858561099b565b73ffffffffffffffffffffffffffffffffffffffff821661069a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f45524332303a206d696e7420746f20746865207a65726f20616464726573730060448201526064016105f9565b80600260008282546106ac9190610f82565b909155505073ffffffffffffffffffffffffffffffffffffffff8216600081815260208181526040808320805486019055518481527fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef910160405180910390a35050565b73ffffffffffffffffffffffffffffffffffffffff83166107b2576040517f08c379a0000000000000000000000000000000000000000000000000000000008152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f2061646460448201527f726573730000000000000000000000000000000000000000000000000000000060648201526084016105f9565b73ffffffffffffffffffffffffffffffffffffffff8216610855576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f20616464726560448201527f737300000000000000000000000000000000000000000000000000000000000060648201526084016105f9565b73ffffffffffffffffffffffffffffffffffffffff83811660008181526001602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92591015b60405180910390a3505050565b73ffffffffffffffffffffffffffffffffffffffff8381166000908152600160209081526040808320938616835292905220547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81146109955781811015610988576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f45524332303a20696e73756666696369656e7420616c6c6f77616e636500000060448201526064016105f9565b6109958484848403610710565b50505050565b73ffffffffffffffffffffffffffffffffffffffff8316610a3e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602560248201527f45524332303a207472616e736665722066726f6d20746865207a65726f20616460448201527f647265737300000000000000000000000000000000000000000000000000000060648201526084016105f9565b73ffffffffffffffffffffffffffffffffffffffff8216610ae1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602360248201527f45524332303a207472616e7366657220746f20746865207a65726f206164647260448201527f657373000000000000000000000000000000000000000000000000000000000060648201526084016105f9565b73ffffffffffffffffffffffffffffffffffffffff831660009081526020819052604090205481811015610b97576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f45524332303a207472616e7366657220616d6f756e742065786365656473206260448201527f616c616e6365000000000000000000000000000000000000000000000000000060648201526084016105f9565b73ffffffffffffffffffffffffffffffffffffffff848116600081815260208181526040808320878703905593871680835291849020805487019055925185815290927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef910160405180910390a3610995565b73ffffffffffffffffffffffffffffffffffffffff8216610cad576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602160248201527f45524332303a206275726e2066726f6d20746865207a65726f2061646472657360448201527f730000000000000000000000000000000000000000000000000000000000000060648201526084016105f9565b73ffffffffffffffffffffffffffffffffffffffff821660009081526020819052604090205481811015610d63576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602260248201527f45524332303a206275726e20616d6f756e7420657863656564732062616c616e60448201527f636500000000000000000000000000000000000000000000000000000000000060648201526084016105f9565b73ffffffffffffffffffffffffffffffffffffffff83166000818152602081815260408083208686039055600280548790039055518581529192917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef91016108b7565b600060208083528351808285015260005b81811015610df357858101830151858201604001528201610dd7565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b803573ffffffffffffffffffffffffffffffffffffffff81168114610e5657600080fd5b919050565b60008060408385031215610e6e57600080fd5b610e7783610e32565b946020939093013593505050565b600080600060608486031215610e9a57600080fd5b610ea384610e32565b9250610eb160208501610e32565b9150604084013590509250925092565b600060208284031215610ed357600080fd5b5035919050565b600060208284031215610eec57600080fd5b610ef582610e32565b9392505050565b60008060408385031215610f0f57600080fd5b610f1883610e32565b9150610f2660208401610e32565b90509250929050565b600181811c90821680610f4357607f821691505b602082108103610f7c577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b80820180821115610410577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fdfea164736f6c6343000813000a",
}

var WERC20MockABI = WERC20MockMetaData.ABI

var WERC20MockBin = WERC20MockMetaData.Bin

func DeployWERC20Mock(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *WERC20Mock, error) {
	parsed, err := WERC20MockMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(WERC20MockBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &WERC20Mock{address: address, abi: *parsed, WERC20MockCaller: WERC20MockCaller{contract: contract}, WERC20MockTransactor: WERC20MockTransactor{contract: contract}, WERC20MockFilterer: WERC20MockFilterer{contract: contract}}, nil
}

type WERC20Mock struct {
	address common.Address
	abi     abi.ABI
	WERC20MockCaller
	WERC20MockTransactor
	WERC20MockFilterer
}

type WERC20MockCaller struct {
	contract *bind.BoundContract
}

type WERC20MockTransactor struct {
	contract *bind.BoundContract
}

type WERC20MockFilterer struct {
	contract *bind.BoundContract
}

type WERC20MockSession struct {
	Contract     *WERC20Mock
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type WERC20MockCallerSession struct {
	Contract *WERC20MockCaller
	CallOpts bind.CallOpts
}

type WERC20MockTransactorSession struct {
	Contract     *WERC20MockTransactor
	TransactOpts bind.TransactOpts
}

type WERC20MockRaw struct {
	Contract *WERC20Mock
}

type WERC20MockCallerRaw struct {
	Contract *WERC20MockCaller
}

type WERC20MockTransactorRaw struct {
	Contract *WERC20MockTransactor
}

func NewWERC20Mock(address common.Address, backend bind.ContractBackend) (*WERC20Mock, error) {
	abi, err := abi.JSON(strings.NewReader(WERC20MockABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindWERC20Mock(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &WERC20Mock{address: address, abi: abi, WERC20MockCaller: WERC20MockCaller{contract: contract}, WERC20MockTransactor: WERC20MockTransactor{contract: contract}, WERC20MockFilterer: WERC20MockFilterer{contract: contract}}, nil
}

func NewWERC20MockCaller(address common.Address, caller bind.ContractCaller) (*WERC20MockCaller, error) {
	contract, err := bindWERC20Mock(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &WERC20MockCaller{contract: contract}, nil
}

func NewWERC20MockTransactor(address common.Address, transactor bind.ContractTransactor) (*WERC20MockTransactor, error) {
	contract, err := bindWERC20Mock(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &WERC20MockTransactor{contract: contract}, nil
}

func NewWERC20MockFilterer(address common.Address, filterer bind.ContractFilterer) (*WERC20MockFilterer, error) {
	contract, err := bindWERC20Mock(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &WERC20MockFilterer{contract: contract}, nil
}

func bindWERC20Mock(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := WERC20MockMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_WERC20Mock *WERC20MockRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WERC20Mock.Contract.WERC20MockCaller.contract.Call(opts, result, method, params...)
}

func (_WERC20Mock *WERC20MockRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WERC20Mock.Contract.WERC20MockTransactor.contract.Transfer(opts)
}

func (_WERC20Mock *WERC20MockRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WERC20Mock.Contract.WERC20MockTransactor.contract.Transact(opts, method, params...)
}

func (_WERC20Mock *WERC20MockCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WERC20Mock.Contract.contract.Call(opts, result, method, params...)
}

func (_WERC20Mock *WERC20MockTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WERC20Mock.Contract.contract.Transfer(opts)
}

func (_WERC20Mock *WERC20MockTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WERC20Mock.Contract.contract.Transact(opts, method, params...)
}

func (_WERC20Mock *WERC20MockCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _WERC20Mock.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_WERC20Mock *WERC20MockSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _WERC20Mock.Contract.Allowance(&_WERC20Mock.CallOpts, owner, spender)
}

func (_WERC20Mock *WERC20MockCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _WERC20Mock.Contract.Allowance(&_WERC20Mock.CallOpts, owner, spender)
}

func (_WERC20Mock *WERC20MockCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _WERC20Mock.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_WERC20Mock *WERC20MockSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _WERC20Mock.Contract.BalanceOf(&_WERC20Mock.CallOpts, account)
}

func (_WERC20Mock *WERC20MockCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _WERC20Mock.Contract.BalanceOf(&_WERC20Mock.CallOpts, account)
}

func (_WERC20Mock *WERC20MockCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _WERC20Mock.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_WERC20Mock *WERC20MockSession) Decimals() (uint8, error) {
	return _WERC20Mock.Contract.Decimals(&_WERC20Mock.CallOpts)
}

func (_WERC20Mock *WERC20MockCallerSession) Decimals() (uint8, error) {
	return _WERC20Mock.Contract.Decimals(&_WERC20Mock.CallOpts)
}

func (_WERC20Mock *WERC20MockCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _WERC20Mock.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_WERC20Mock *WERC20MockSession) Name() (string, error) {
	return _WERC20Mock.Contract.Name(&_WERC20Mock.CallOpts)
}

func (_WERC20Mock *WERC20MockCallerSession) Name() (string, error) {
	return _WERC20Mock.Contract.Name(&_WERC20Mock.CallOpts)
}

func (_WERC20Mock *WERC20MockCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _WERC20Mock.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_WERC20Mock *WERC20MockSession) Symbol() (string, error) {
	return _WERC20Mock.Contract.Symbol(&_WERC20Mock.CallOpts)
}

func (_WERC20Mock *WERC20MockCallerSession) Symbol() (string, error) {
	return _WERC20Mock.Contract.Symbol(&_WERC20Mock.CallOpts)
}

func (_WERC20Mock *WERC20MockCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _WERC20Mock.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_WERC20Mock *WERC20MockSession) TotalSupply() (*big.Int, error) {
	return _WERC20Mock.Contract.TotalSupply(&_WERC20Mock.CallOpts)
}

func (_WERC20Mock *WERC20MockCallerSession) TotalSupply() (*big.Int, error) {
	return _WERC20Mock.Contract.TotalSupply(&_WERC20Mock.CallOpts)
}

func (_WERC20Mock *WERC20MockTransactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WERC20Mock.contract.Transact(opts, "approve", spender, amount)
}

func (_WERC20Mock *WERC20MockSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WERC20Mock.Contract.Approve(&_WERC20Mock.TransactOpts, spender, amount)
}

func (_WERC20Mock *WERC20MockTransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WERC20Mock.Contract.Approve(&_WERC20Mock.TransactOpts, spender, amount)
}

func (_WERC20Mock *WERC20MockTransactor) Burn(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WERC20Mock.contract.Transact(opts, "burn", account, amount)
}

func (_WERC20Mock *WERC20MockSession) Burn(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WERC20Mock.Contract.Burn(&_WERC20Mock.TransactOpts, account, amount)
}

func (_WERC20Mock *WERC20MockTransactorSession) Burn(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WERC20Mock.Contract.Burn(&_WERC20Mock.TransactOpts, account, amount)
}

func (_WERC20Mock *WERC20MockTransactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _WERC20Mock.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

func (_WERC20Mock *WERC20MockSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _WERC20Mock.Contract.DecreaseAllowance(&_WERC20Mock.TransactOpts, spender, subtractedValue)
}

func (_WERC20Mock *WERC20MockTransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _WERC20Mock.Contract.DecreaseAllowance(&_WERC20Mock.TransactOpts, spender, subtractedValue)
}

func (_WERC20Mock *WERC20MockTransactor) Deposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WERC20Mock.contract.Transact(opts, "deposit")
}

func (_WERC20Mock *WERC20MockSession) Deposit() (*types.Transaction, error) {
	return _WERC20Mock.Contract.Deposit(&_WERC20Mock.TransactOpts)
}

func (_WERC20Mock *WERC20MockTransactorSession) Deposit() (*types.Transaction, error) {
	return _WERC20Mock.Contract.Deposit(&_WERC20Mock.TransactOpts)
}

func (_WERC20Mock *WERC20MockTransactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _WERC20Mock.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

func (_WERC20Mock *WERC20MockSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _WERC20Mock.Contract.IncreaseAllowance(&_WERC20Mock.TransactOpts, spender, addedValue)
}

func (_WERC20Mock *WERC20MockTransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _WERC20Mock.Contract.IncreaseAllowance(&_WERC20Mock.TransactOpts, spender, addedValue)
}

func (_WERC20Mock *WERC20MockTransactor) Mint(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WERC20Mock.contract.Transact(opts, "mint", account, amount)
}

func (_WERC20Mock *WERC20MockSession) Mint(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WERC20Mock.Contract.Mint(&_WERC20Mock.TransactOpts, account, amount)
}

func (_WERC20Mock *WERC20MockTransactorSession) Mint(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WERC20Mock.Contract.Mint(&_WERC20Mock.TransactOpts, account, amount)
}

func (_WERC20Mock *WERC20MockTransactor) Transfer(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WERC20Mock.contract.Transact(opts, "transfer", to, amount)
}

func (_WERC20Mock *WERC20MockSession) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WERC20Mock.Contract.Transfer(&_WERC20Mock.TransactOpts, to, amount)
}

func (_WERC20Mock *WERC20MockTransactorSession) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WERC20Mock.Contract.Transfer(&_WERC20Mock.TransactOpts, to, amount)
}

func (_WERC20Mock *WERC20MockTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WERC20Mock.contract.Transact(opts, "transferFrom", from, to, amount)
}

func (_WERC20Mock *WERC20MockSession) TransferFrom(from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WERC20Mock.Contract.TransferFrom(&_WERC20Mock.TransactOpts, from, to, amount)
}

func (_WERC20Mock *WERC20MockTransactorSession) TransferFrom(from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WERC20Mock.Contract.TransferFrom(&_WERC20Mock.TransactOpts, from, to, amount)
}

func (_WERC20Mock *WERC20MockTransactor) Withdraw(opts *bind.TransactOpts, wad *big.Int) (*types.Transaction, error) {
	return _WERC20Mock.contract.Transact(opts, "withdraw", wad)
}

func (_WERC20Mock *WERC20MockSession) Withdraw(wad *big.Int) (*types.Transaction, error) {
	return _WERC20Mock.Contract.Withdraw(&_WERC20Mock.TransactOpts, wad)
}

func (_WERC20Mock *WERC20MockTransactorSession) Withdraw(wad *big.Int) (*types.Transaction, error) {
	return _WERC20Mock.Contract.Withdraw(&_WERC20Mock.TransactOpts, wad)
}

func (_WERC20Mock *WERC20MockTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WERC20Mock.contract.RawTransact(opts, nil)
}

func (_WERC20Mock *WERC20MockSession) Receive() (*types.Transaction, error) {
	return _WERC20Mock.Contract.Receive(&_WERC20Mock.TransactOpts)
}

func (_WERC20Mock *WERC20MockTransactorSession) Receive() (*types.Transaction, error) {
	return _WERC20Mock.Contract.Receive(&_WERC20Mock.TransactOpts)
}

type WERC20MockApprovalIterator struct {
	Event *WERC20MockApproval

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WERC20MockApprovalIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WERC20MockApproval)
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
		it.Event = new(WERC20MockApproval)
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

func (it *WERC20MockApprovalIterator) Error() error {
	return it.fail
}

func (it *WERC20MockApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WERC20MockApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log
}

func (_WERC20Mock *WERC20MockFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*WERC20MockApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _WERC20Mock.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &WERC20MockApprovalIterator{contract: _WERC20Mock.contract, event: "Approval", logs: logs, sub: sub}, nil
}

func (_WERC20Mock *WERC20MockFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *WERC20MockApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _WERC20Mock.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WERC20MockApproval)
				if err := _WERC20Mock.contract.UnpackLog(event, "Approval", log); err != nil {
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

func (_WERC20Mock *WERC20MockFilterer) ParseApproval(log types.Log) (*WERC20MockApproval, error) {
	event := new(WERC20MockApproval)
	if err := _WERC20Mock.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type WERC20MockDepositIterator struct {
	Event *WERC20MockDeposit

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WERC20MockDepositIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WERC20MockDeposit)
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
		it.Event = new(WERC20MockDeposit)
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

func (it *WERC20MockDepositIterator) Error() error {
	return it.fail
}

func (it *WERC20MockDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WERC20MockDeposit struct {
	Dst common.Address
	Wad *big.Int
	Raw types.Log
}

func (_WERC20Mock *WERC20MockFilterer) FilterDeposit(opts *bind.FilterOpts, dst []common.Address) (*WERC20MockDepositIterator, error) {

	var dstRule []interface{}
	for _, dstItem := range dst {
		dstRule = append(dstRule, dstItem)
	}

	logs, sub, err := _WERC20Mock.contract.FilterLogs(opts, "Deposit", dstRule)
	if err != nil {
		return nil, err
	}
	return &WERC20MockDepositIterator{contract: _WERC20Mock.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

func (_WERC20Mock *WERC20MockFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *WERC20MockDeposit, dst []common.Address) (event.Subscription, error) {

	var dstRule []interface{}
	for _, dstItem := range dst {
		dstRule = append(dstRule, dstItem)
	}

	logs, sub, err := _WERC20Mock.contract.WatchLogs(opts, "Deposit", dstRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WERC20MockDeposit)
				if err := _WERC20Mock.contract.UnpackLog(event, "Deposit", log); err != nil {
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

func (_WERC20Mock *WERC20MockFilterer) ParseDeposit(log types.Log) (*WERC20MockDeposit, error) {
	event := new(WERC20MockDeposit)
	if err := _WERC20Mock.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type WERC20MockTransferIterator struct {
	Event *WERC20MockTransfer

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WERC20MockTransferIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WERC20MockTransfer)
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
		it.Event = new(WERC20MockTransfer)
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

func (it *WERC20MockTransferIterator) Error() error {
	return it.fail
}

func (it *WERC20MockTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WERC20MockTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log
}

func (_WERC20Mock *WERC20MockFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*WERC20MockTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _WERC20Mock.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &WERC20MockTransferIterator{contract: _WERC20Mock.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

func (_WERC20Mock *WERC20MockFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *WERC20MockTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _WERC20Mock.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WERC20MockTransfer)
				if err := _WERC20Mock.contract.UnpackLog(event, "Transfer", log); err != nil {
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

func (_WERC20Mock *WERC20MockFilterer) ParseTransfer(log types.Log) (*WERC20MockTransfer, error) {
	event := new(WERC20MockTransfer)
	if err := _WERC20Mock.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type WERC20MockWithdrawalIterator struct {
	Event *WERC20MockWithdrawal

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WERC20MockWithdrawalIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WERC20MockWithdrawal)
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
		it.Event = new(WERC20MockWithdrawal)
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

func (it *WERC20MockWithdrawalIterator) Error() error {
	return it.fail
}

func (it *WERC20MockWithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WERC20MockWithdrawal struct {
	Src common.Address
	Wad *big.Int
	Raw types.Log
}

func (_WERC20Mock *WERC20MockFilterer) FilterWithdrawal(opts *bind.FilterOpts, src []common.Address) (*WERC20MockWithdrawalIterator, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}

	logs, sub, err := _WERC20Mock.contract.FilterLogs(opts, "Withdrawal", srcRule)
	if err != nil {
		return nil, err
	}
	return &WERC20MockWithdrawalIterator{contract: _WERC20Mock.contract, event: "Withdrawal", logs: logs, sub: sub}, nil
}

func (_WERC20Mock *WERC20MockFilterer) WatchWithdrawal(opts *bind.WatchOpts, sink chan<- *WERC20MockWithdrawal, src []common.Address) (event.Subscription, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}

	logs, sub, err := _WERC20Mock.contract.WatchLogs(opts, "Withdrawal", srcRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WERC20MockWithdrawal)
				if err := _WERC20Mock.contract.UnpackLog(event, "Withdrawal", log); err != nil {
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

func (_WERC20Mock *WERC20MockFilterer) ParseWithdrawal(log types.Log) (*WERC20MockWithdrawal, error) {
	event := new(WERC20MockWithdrawal)
	if err := _WERC20Mock.contract.UnpackLog(event, "Withdrawal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_WERC20Mock *WERC20Mock) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _WERC20Mock.abi.Events["Approval"].ID:
		return _WERC20Mock.ParseApproval(log)
	case _WERC20Mock.abi.Events["Deposit"].ID:
		return _WERC20Mock.ParseDeposit(log)
	case _WERC20Mock.abi.Events["Transfer"].ID:
		return _WERC20Mock.ParseTransfer(log)
	case _WERC20Mock.abi.Events["Withdrawal"].ID:
		return _WERC20Mock.ParseWithdrawal(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (WERC20MockApproval) Topic() common.Hash {
	return common.HexToHash("0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925")
}

func (WERC20MockDeposit) Topic() common.Hash {
	return common.HexToHash("0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c")
}

func (WERC20MockTransfer) Topic() common.Hash {
	return common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")
}

func (WERC20MockWithdrawal) Topic() common.Hash {
	return common.HexToHash("0x7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65")
}

func (_WERC20Mock *WERC20Mock) Address() common.Address {
	return _WERC20Mock.address
}

type WERC20MockInterface interface {
	Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error)

	BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error)

	Decimals(opts *bind.CallOpts) (uint8, error)

	Name(opts *bind.CallOpts) (string, error)

	Symbol(opts *bind.CallOpts) (string, error)

	TotalSupply(opts *bind.CallOpts) (*big.Int, error)

	Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error)

	Burn(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error)

	DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error)

	Deposit(opts *bind.TransactOpts) (*types.Transaction, error)

	IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error)

	Mint(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error)

	Transfer(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error)

	TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error)

	Withdraw(opts *bind.TransactOpts, wad *big.Int) (*types.Transaction, error)

	Receive(opts *bind.TransactOpts) (*types.Transaction, error)

	FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*WERC20MockApprovalIterator, error)

	WatchApproval(opts *bind.WatchOpts, sink chan<- *WERC20MockApproval, owner []common.Address, spender []common.Address) (event.Subscription, error)

	ParseApproval(log types.Log) (*WERC20MockApproval, error)

	FilterDeposit(opts *bind.FilterOpts, dst []common.Address) (*WERC20MockDepositIterator, error)

	WatchDeposit(opts *bind.WatchOpts, sink chan<- *WERC20MockDeposit, dst []common.Address) (event.Subscription, error)

	ParseDeposit(log types.Log) (*WERC20MockDeposit, error)

	FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*WERC20MockTransferIterator, error)

	WatchTransfer(opts *bind.WatchOpts, sink chan<- *WERC20MockTransfer, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseTransfer(log types.Log) (*WERC20MockTransfer, error)

	FilterWithdrawal(opts *bind.FilterOpts, src []common.Address) (*WERC20MockWithdrawalIterator, error)

	WatchWithdrawal(opts *bind.WatchOpts, sink chan<- *WERC20MockWithdrawal, src []common.Address) (event.Subscription, error)

	ParseWithdrawal(log types.Log) (*WERC20MockWithdrawal, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
