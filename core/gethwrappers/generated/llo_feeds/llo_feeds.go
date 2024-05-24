// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package llo_feeds

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

var LLOVerifierProxyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"accessController\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessForbidden\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"verifier\",\"type\":\"address\"}],\"name\":\"ConfigDigestAlreadySet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"verifier\",\"type\":\"address\"}],\"name\":\"VerifierAlreadyInitialized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"VerifierInvalid\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"VerifierNotFound\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldAccessController\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAccessController\",\"type\":\"address\"}],\"name\":\"AccessControllerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"verifierAddress\",\"type\":\"address\"}],\"name\":\"VerifierInitialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"oldConfigDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"newConfigDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"verifierAddress\",\"type\":\"address\"}],\"name\":\"VerifierSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"verifierAddress\",\"type\":\"address\"}],\"name\":\"VerifierUnset\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAccessController\",\"outputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"accessController\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"getVerifier\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"verifierAddress\",\"type\":\"address\"}],\"name\":\"initializeVerifier\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"accessController\",\"type\":\"address\"}],\"name\":\"setAccessController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"currentConfigDigest\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newConfigDigest\",\"type\":\"bytes32\"}],\"name\":\"setVerifier\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"unsetVerifier\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"signedReport\",\"type\":\"bytes\"}],\"name\":\"verify\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"verifierResponse\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060405161112e38038061112e83398101604081905261002f91610187565b33806000816100855760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156100b5576100b5816100de565b5050600480546001600160a01b0319166001600160a01b039390931692909217909155506101b7565b336001600160a01b038216036101365760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161007c565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60006020828403121561019957600080fd5b81516001600160a01b03811681146101b057600080fd5b9392505050565b610f68806101c66000396000f3fe608060405234801561001057600080fd5b50600436106100c95760003560e01c80638c2a4d5311610081578063eeb7b2481161005b578063eeb7b248146101c8578063f08391d8146101fe578063f2fde38b1461021157600080fd5b80638c2a4d53146101845780638da5cb5b146101975780638e760afe146101b557600080fd5b80632cc99477116100b25780632cc99477146101545780636e9140941461016957806379ba50971461017c57600080fd5b806316d6b5f6146100ce578063181f5a7714610112575b600080fd5b60045473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b60408051808201909152601381527f566572696669657250726f787920312e302e300000000000000000000000000060208201525b6040516101099190610c40565b610167610162366004610c5a565b610224565b005b610167610177366004610c7c565b61036f565b610167610467565b610167610192366004610cb7565b610564565b60005473ffffffffffffffffffffffffffffffffffffffff166100e8565b6101476101c3366004610cd4565b610795565b6100e86101d6366004610c7c565b60009081526003602052604090205473ffffffffffffffffffffffffffffffffffffffff1690565b61016761020c366004610cb7565b6109bf565b61016761021f366004610cb7565b610a46565b600081815260036020526040902054819073ffffffffffffffffffffffffffffffffffffffff1680156102a7576040517f375d1fe60000000000000000000000000000000000000000000000000000000081526004810183905273ffffffffffffffffffffffffffffffffffffffff821660248201526044015b60405180910390fd5b3360009081526002602052604090205460ff166102f0576040517fef67f5d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000838152600360209081526040918290208054337fffffffffffffffffffffffff0000000000000000000000000000000000000000909116811790915582518781529182018690528183015290517fbeb513e532542a562ac35699e7cd9ae7d198dcd3eee15bada6c857d28ceaddcf9181900360600190a150505050565b610377610a5a565b60008181526003602052604090205473ffffffffffffffffffffffffffffffffffffffff16806103d6576040517fb151802b0000000000000000000000000000000000000000000000000000000081526004810183905260240161029e565b6000828152600360205260409081902080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169055517f11dc15c4b8ac2b183166cc8427e5385a5ece8308217a4217338c6a7614845c4c9061045b908490849091825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b60405180910390a15050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146104e8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015260640161029e565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b61056c610a5a565b8073ffffffffffffffffffffffffffffffffffffffff81166105ba576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040517f01ffc9a70000000000000000000000000000000000000000000000000000000081527f3d3ac1b500000000000000000000000000000000000000000000000000000000600482015273ffffffffffffffffffffffffffffffffffffffff8216906301ffc9a790602401602060405180830381865afa158015610644573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106689190610d46565b61069e576040517f75b0527a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff821660009081526002602052604090205460ff1615610716576040517f4e01ccfd00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8316600482015260240161029e565b73ffffffffffffffffffffffffffffffffffffffff821660008181526002602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016600117905590519182527f1f2cd7c97f4d801b5efe26cc409617c1fd6c5ef786e79aacb90af40923e4e8e9910161045b565b60045460609073ffffffffffffffffffffffffffffffffffffffff16801580159061085557506040517f6b14daf800000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821690636b14daf8906108129033906000903690600401610db1565b602060405180830381865afa15801561082f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108539190610d46565b155b1561088c576040517fef67f5d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006108988486610dea565b60008181526003602052604090205490915073ffffffffffffffffffffffffffffffffffffffff16806108fa576040517fb151802b0000000000000000000000000000000000000000000000000000000081526004810183905260240161029e565b6040517f3d3ac1b500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821690633d3ac1b59061095090899089903390600401610e27565b6000604051808303816000875af115801561096f573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526109b59190810190610e90565b9695505050505050565b6109c7610a5a565b6004805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff000000000000000000000000000000000000000083168117909355604080519190921680825260208201939093527f953e92b1a6442e9c3242531154a3f6f6eb00b4e9c719ba8118fa6235e4ce89b6910161045b565b610a4e610a5a565b610a5781610add565b50565b60005473ffffffffffffffffffffffffffffffffffffffff163314610adb576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161029e565b565b3373ffffffffffffffffffffffffffffffffffffffff821603610b5c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161029e565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60005b83811015610bed578181015183820152602001610bd5565b50506000910152565b60008151808452610c0e816020860160208601610bd2565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000610c536020830184610bf6565b9392505050565b60008060408385031215610c6d57600080fd5b50508035926020909101359150565b600060208284031215610c8e57600080fd5b5035919050565b73ffffffffffffffffffffffffffffffffffffffff81168114610a5757600080fd5b600060208284031215610cc957600080fd5b8135610c5381610c95565b60008060208385031215610ce757600080fd5b823567ffffffffffffffff80821115610cff57600080fd5b818501915085601f830112610d1357600080fd5b813581811115610d2257600080fd5b866020828501011115610d3457600080fd5b60209290920196919550909350505050565b600060208284031215610d5857600080fd5b81518015158114610c5357600080fd5b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b73ffffffffffffffffffffffffffffffffffffffff84168152604060208201526000610de1604083018486610d68565b95945050505050565b80356020831015610e21577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff602084900360031b1b165b92915050565b604081526000610e3b604083018587610d68565b905073ffffffffffffffffffffffffffffffffffffffff83166020830152949350505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600060208284031215610ea257600080fd5b815167ffffffffffffffff80821115610eba57600080fd5b818401915084601f830112610ece57600080fd5b815181811115610ee057610ee0610e61565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f01168101908382118183101715610f2657610f26610e61565b81604052828152876020848701011115610f3f57600080fd5b610f50836020830160208801610bd2565b97965050505050505056fea164736f6c6343000810000a",
}

var LLOVerifierProxyABI = LLOVerifierProxyMetaData.ABI

var LLOVerifierProxyBin = LLOVerifierProxyMetaData.Bin

func DeployLLOVerifierProxy(auth *bind.TransactOpts, backend bind.ContractBackend, accessController common.Address) (common.Address, *types.Transaction, *LLOVerifierProxy, error) {
	parsed, err := LLOVerifierProxyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(LLOVerifierProxyBin), backend, accessController)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LLOVerifierProxy{LLOVerifierProxyCaller: LLOVerifierProxyCaller{contract: contract}, LLOVerifierProxyTransactor: LLOVerifierProxyTransactor{contract: contract}, LLOVerifierProxyFilterer: LLOVerifierProxyFilterer{contract: contract}}, nil
}

type LLOVerifierProxy struct {
	address common.Address
	abi     abi.ABI
	LLOVerifierProxyCaller
	LLOVerifierProxyTransactor
	LLOVerifierProxyFilterer
}

type LLOVerifierProxyCaller struct {
	contract *bind.BoundContract
}

type LLOVerifierProxyTransactor struct {
	contract *bind.BoundContract
}

type LLOVerifierProxyFilterer struct {
	contract *bind.BoundContract
}

type LLOVerifierProxySession struct {
	Contract     *LLOVerifierProxy
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type LLOVerifierProxyCallerSession struct {
	Contract *LLOVerifierProxyCaller
	CallOpts bind.CallOpts
}

type LLOVerifierProxyTransactorSession struct {
	Contract     *LLOVerifierProxyTransactor
	TransactOpts bind.TransactOpts
}

type LLOVerifierProxyRaw struct {
	Contract *LLOVerifierProxy
}

type LLOVerifierProxyCallerRaw struct {
	Contract *LLOVerifierProxyCaller
}

type LLOVerifierProxyTransactorRaw struct {
	Contract *LLOVerifierProxyTransactor
}

func NewLLOVerifierProxy(address common.Address, backend bind.ContractBackend) (*LLOVerifierProxy, error) {
	abi, err := abi.JSON(strings.NewReader(LLOVerifierProxyABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindLLOVerifierProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LLOVerifierProxy{address: address, abi: abi, LLOVerifierProxyCaller: LLOVerifierProxyCaller{contract: contract}, LLOVerifierProxyTransactor: LLOVerifierProxyTransactor{contract: contract}, LLOVerifierProxyFilterer: LLOVerifierProxyFilterer{contract: contract}}, nil
}

func NewLLOVerifierProxyCaller(address common.Address, caller bind.ContractCaller) (*LLOVerifierProxyCaller, error) {
	contract, err := bindLLOVerifierProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LLOVerifierProxyCaller{contract: contract}, nil
}

func NewLLOVerifierProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*LLOVerifierProxyTransactor, error) {
	contract, err := bindLLOVerifierProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LLOVerifierProxyTransactor{contract: contract}, nil
}

func NewLLOVerifierProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*LLOVerifierProxyFilterer, error) {
	contract, err := bindLLOVerifierProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LLOVerifierProxyFilterer{contract: contract}, nil
}

func bindLLOVerifierProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := LLOVerifierProxyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_LLOVerifierProxy *LLOVerifierProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LLOVerifierProxy.Contract.LLOVerifierProxyCaller.contract.Call(opts, result, method, params...)
}

func (_LLOVerifierProxy *LLOVerifierProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LLOVerifierProxy.Contract.LLOVerifierProxyTransactor.contract.Transfer(opts)
}

func (_LLOVerifierProxy *LLOVerifierProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LLOVerifierProxy.Contract.LLOVerifierProxyTransactor.contract.Transact(opts, method, params...)
}

func (_LLOVerifierProxy *LLOVerifierProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LLOVerifierProxy.Contract.contract.Call(opts, result, method, params...)
}

func (_LLOVerifierProxy *LLOVerifierProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LLOVerifierProxy.Contract.contract.Transfer(opts)
}

func (_LLOVerifierProxy *LLOVerifierProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LLOVerifierProxy.Contract.contract.Transact(opts, method, params...)
}

func (_LLOVerifierProxy *LLOVerifierProxyCaller) GetAccessController(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LLOVerifierProxy.contract.Call(opts, &out, "getAccessController")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_LLOVerifierProxy *LLOVerifierProxySession) GetAccessController() (common.Address, error) {
	return _LLOVerifierProxy.Contract.GetAccessController(&_LLOVerifierProxy.CallOpts)
}

func (_LLOVerifierProxy *LLOVerifierProxyCallerSession) GetAccessController() (common.Address, error) {
	return _LLOVerifierProxy.Contract.GetAccessController(&_LLOVerifierProxy.CallOpts)
}

func (_LLOVerifierProxy *LLOVerifierProxyCaller) GetVerifier(opts *bind.CallOpts, configDigest [32]byte) (common.Address, error) {
	var out []interface{}
	err := _LLOVerifierProxy.contract.Call(opts, &out, "getVerifier", configDigest)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_LLOVerifierProxy *LLOVerifierProxySession) GetVerifier(configDigest [32]byte) (common.Address, error) {
	return _LLOVerifierProxy.Contract.GetVerifier(&_LLOVerifierProxy.CallOpts, configDigest)
}

func (_LLOVerifierProxy *LLOVerifierProxyCallerSession) GetVerifier(configDigest [32]byte) (common.Address, error) {
	return _LLOVerifierProxy.Contract.GetVerifier(&_LLOVerifierProxy.CallOpts, configDigest)
}

func (_LLOVerifierProxy *LLOVerifierProxyCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LLOVerifierProxy.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_LLOVerifierProxy *LLOVerifierProxySession) Owner() (common.Address, error) {
	return _LLOVerifierProxy.Contract.Owner(&_LLOVerifierProxy.CallOpts)
}

func (_LLOVerifierProxy *LLOVerifierProxyCallerSession) Owner() (common.Address, error) {
	return _LLOVerifierProxy.Contract.Owner(&_LLOVerifierProxy.CallOpts)
}

func (_LLOVerifierProxy *LLOVerifierProxyCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _LLOVerifierProxy.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_LLOVerifierProxy *LLOVerifierProxySession) TypeAndVersion() (string, error) {
	return _LLOVerifierProxy.Contract.TypeAndVersion(&_LLOVerifierProxy.CallOpts)
}

func (_LLOVerifierProxy *LLOVerifierProxyCallerSession) TypeAndVersion() (string, error) {
	return _LLOVerifierProxy.Contract.TypeAndVersion(&_LLOVerifierProxy.CallOpts)
}

func (_LLOVerifierProxy *LLOVerifierProxyTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LLOVerifierProxy.contract.Transact(opts, "acceptOwnership")
}

func (_LLOVerifierProxy *LLOVerifierProxySession) AcceptOwnership() (*types.Transaction, error) {
	return _LLOVerifierProxy.Contract.AcceptOwnership(&_LLOVerifierProxy.TransactOpts)
}

func (_LLOVerifierProxy *LLOVerifierProxyTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _LLOVerifierProxy.Contract.AcceptOwnership(&_LLOVerifierProxy.TransactOpts)
}

func (_LLOVerifierProxy *LLOVerifierProxyTransactor) InitializeVerifier(opts *bind.TransactOpts, verifierAddress common.Address) (*types.Transaction, error) {
	return _LLOVerifierProxy.contract.Transact(opts, "initializeVerifier", verifierAddress)
}

func (_LLOVerifierProxy *LLOVerifierProxySession) InitializeVerifier(verifierAddress common.Address) (*types.Transaction, error) {
	return _LLOVerifierProxy.Contract.InitializeVerifier(&_LLOVerifierProxy.TransactOpts, verifierAddress)
}

func (_LLOVerifierProxy *LLOVerifierProxyTransactorSession) InitializeVerifier(verifierAddress common.Address) (*types.Transaction, error) {
	return _LLOVerifierProxy.Contract.InitializeVerifier(&_LLOVerifierProxy.TransactOpts, verifierAddress)
}

func (_LLOVerifierProxy *LLOVerifierProxyTransactor) SetAccessController(opts *bind.TransactOpts, accessController common.Address) (*types.Transaction, error) {
	return _LLOVerifierProxy.contract.Transact(opts, "setAccessController", accessController)
}

func (_LLOVerifierProxy *LLOVerifierProxySession) SetAccessController(accessController common.Address) (*types.Transaction, error) {
	return _LLOVerifierProxy.Contract.SetAccessController(&_LLOVerifierProxy.TransactOpts, accessController)
}

func (_LLOVerifierProxy *LLOVerifierProxyTransactorSession) SetAccessController(accessController common.Address) (*types.Transaction, error) {
	return _LLOVerifierProxy.Contract.SetAccessController(&_LLOVerifierProxy.TransactOpts, accessController)
}

func (_LLOVerifierProxy *LLOVerifierProxyTransactor) SetVerifier(opts *bind.TransactOpts, currentConfigDigest [32]byte, newConfigDigest [32]byte) (*types.Transaction, error) {
	return _LLOVerifierProxy.contract.Transact(opts, "setVerifier", currentConfigDigest, newConfigDigest)
}

func (_LLOVerifierProxy *LLOVerifierProxySession) SetVerifier(currentConfigDigest [32]byte, newConfigDigest [32]byte) (*types.Transaction, error) {
	return _LLOVerifierProxy.Contract.SetVerifier(&_LLOVerifierProxy.TransactOpts, currentConfigDigest, newConfigDigest)
}

func (_LLOVerifierProxy *LLOVerifierProxyTransactorSession) SetVerifier(currentConfigDigest [32]byte, newConfigDigest [32]byte) (*types.Transaction, error) {
	return _LLOVerifierProxy.Contract.SetVerifier(&_LLOVerifierProxy.TransactOpts, currentConfigDigest, newConfigDigest)
}

func (_LLOVerifierProxy *LLOVerifierProxyTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _LLOVerifierProxy.contract.Transact(opts, "transferOwnership", to)
}

func (_LLOVerifierProxy *LLOVerifierProxySession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _LLOVerifierProxy.Contract.TransferOwnership(&_LLOVerifierProxy.TransactOpts, to)
}

func (_LLOVerifierProxy *LLOVerifierProxyTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _LLOVerifierProxy.Contract.TransferOwnership(&_LLOVerifierProxy.TransactOpts, to)
}

func (_LLOVerifierProxy *LLOVerifierProxyTransactor) UnsetVerifier(opts *bind.TransactOpts, configDigest [32]byte) (*types.Transaction, error) {
	return _LLOVerifierProxy.contract.Transact(opts, "unsetVerifier", configDigest)
}

func (_LLOVerifierProxy *LLOVerifierProxySession) UnsetVerifier(configDigest [32]byte) (*types.Transaction, error) {
	return _LLOVerifierProxy.Contract.UnsetVerifier(&_LLOVerifierProxy.TransactOpts, configDigest)
}

func (_LLOVerifierProxy *LLOVerifierProxyTransactorSession) UnsetVerifier(configDigest [32]byte) (*types.Transaction, error) {
	return _LLOVerifierProxy.Contract.UnsetVerifier(&_LLOVerifierProxy.TransactOpts, configDigest)
}

func (_LLOVerifierProxy *LLOVerifierProxyTransactor) Verify(opts *bind.TransactOpts, signedReport []byte) (*types.Transaction, error) {
	return _LLOVerifierProxy.contract.Transact(opts, "verify", signedReport)
}

func (_LLOVerifierProxy *LLOVerifierProxySession) Verify(signedReport []byte) (*types.Transaction, error) {
	return _LLOVerifierProxy.Contract.Verify(&_LLOVerifierProxy.TransactOpts, signedReport)
}

func (_LLOVerifierProxy *LLOVerifierProxyTransactorSession) Verify(signedReport []byte) (*types.Transaction, error) {
	return _LLOVerifierProxy.Contract.Verify(&_LLOVerifierProxy.TransactOpts, signedReport)
}

type LLOVerifierProxyAccessControllerSetIterator struct {
	Event *LLOVerifierProxyAccessControllerSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LLOVerifierProxyAccessControllerSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LLOVerifierProxyAccessControllerSet)
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
		it.Event = new(LLOVerifierProxyAccessControllerSet)
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

func (it *LLOVerifierProxyAccessControllerSetIterator) Error() error {
	return it.fail
}

func (it *LLOVerifierProxyAccessControllerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LLOVerifierProxyAccessControllerSet struct {
	OldAccessController common.Address
	NewAccessController common.Address
	Raw                 types.Log
}

func (_LLOVerifierProxy *LLOVerifierProxyFilterer) FilterAccessControllerSet(opts *bind.FilterOpts) (*LLOVerifierProxyAccessControllerSetIterator, error) {

	logs, sub, err := _LLOVerifierProxy.contract.FilterLogs(opts, "AccessControllerSet")
	if err != nil {
		return nil, err
	}
	return &LLOVerifierProxyAccessControllerSetIterator{contract: _LLOVerifierProxy.contract, event: "AccessControllerSet", logs: logs, sub: sub}, nil
}

func (_LLOVerifierProxy *LLOVerifierProxyFilterer) WatchAccessControllerSet(opts *bind.WatchOpts, sink chan<- *LLOVerifierProxyAccessControllerSet) (event.Subscription, error) {

	logs, sub, err := _LLOVerifierProxy.contract.WatchLogs(opts, "AccessControllerSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LLOVerifierProxyAccessControllerSet)
				if err := _LLOVerifierProxy.contract.UnpackLog(event, "AccessControllerSet", log); err != nil {
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

func (_LLOVerifierProxy *LLOVerifierProxyFilterer) ParseAccessControllerSet(log types.Log) (*LLOVerifierProxyAccessControllerSet, error) {
	event := new(LLOVerifierProxyAccessControllerSet)
	if err := _LLOVerifierProxy.contract.UnpackLog(event, "AccessControllerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LLOVerifierProxyOwnershipTransferRequestedIterator struct {
	Event *LLOVerifierProxyOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LLOVerifierProxyOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LLOVerifierProxyOwnershipTransferRequested)
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
		it.Event = new(LLOVerifierProxyOwnershipTransferRequested)
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

func (it *LLOVerifierProxyOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *LLOVerifierProxyOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LLOVerifierProxyOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_LLOVerifierProxy *LLOVerifierProxyFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LLOVerifierProxyOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LLOVerifierProxy.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &LLOVerifierProxyOwnershipTransferRequestedIterator{contract: _LLOVerifierProxy.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_LLOVerifierProxy *LLOVerifierProxyFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *LLOVerifierProxyOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LLOVerifierProxy.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LLOVerifierProxyOwnershipTransferRequested)
				if err := _LLOVerifierProxy.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_LLOVerifierProxy *LLOVerifierProxyFilterer) ParseOwnershipTransferRequested(log types.Log) (*LLOVerifierProxyOwnershipTransferRequested, error) {
	event := new(LLOVerifierProxyOwnershipTransferRequested)
	if err := _LLOVerifierProxy.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LLOVerifierProxyOwnershipTransferredIterator struct {
	Event *LLOVerifierProxyOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LLOVerifierProxyOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LLOVerifierProxyOwnershipTransferred)
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
		it.Event = new(LLOVerifierProxyOwnershipTransferred)
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

func (it *LLOVerifierProxyOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *LLOVerifierProxyOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LLOVerifierProxyOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_LLOVerifierProxy *LLOVerifierProxyFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LLOVerifierProxyOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LLOVerifierProxy.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &LLOVerifierProxyOwnershipTransferredIterator{contract: _LLOVerifierProxy.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_LLOVerifierProxy *LLOVerifierProxyFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *LLOVerifierProxyOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LLOVerifierProxy.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LLOVerifierProxyOwnershipTransferred)
				if err := _LLOVerifierProxy.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_LLOVerifierProxy *LLOVerifierProxyFilterer) ParseOwnershipTransferred(log types.Log) (*LLOVerifierProxyOwnershipTransferred, error) {
	event := new(LLOVerifierProxyOwnershipTransferred)
	if err := _LLOVerifierProxy.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LLOVerifierProxyVerifierInitializedIterator struct {
	Event *LLOVerifierProxyVerifierInitialized

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LLOVerifierProxyVerifierInitializedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LLOVerifierProxyVerifierInitialized)
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
		it.Event = new(LLOVerifierProxyVerifierInitialized)
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

func (it *LLOVerifierProxyVerifierInitializedIterator) Error() error {
	return it.fail
}

func (it *LLOVerifierProxyVerifierInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LLOVerifierProxyVerifierInitialized struct {
	VerifierAddress common.Address
	Raw             types.Log
}

func (_LLOVerifierProxy *LLOVerifierProxyFilterer) FilterVerifierInitialized(opts *bind.FilterOpts) (*LLOVerifierProxyVerifierInitializedIterator, error) {

	logs, sub, err := _LLOVerifierProxy.contract.FilterLogs(opts, "VerifierInitialized")
	if err != nil {
		return nil, err
	}
	return &LLOVerifierProxyVerifierInitializedIterator{contract: _LLOVerifierProxy.contract, event: "VerifierInitialized", logs: logs, sub: sub}, nil
}

func (_LLOVerifierProxy *LLOVerifierProxyFilterer) WatchVerifierInitialized(opts *bind.WatchOpts, sink chan<- *LLOVerifierProxyVerifierInitialized) (event.Subscription, error) {

	logs, sub, err := _LLOVerifierProxy.contract.WatchLogs(opts, "VerifierInitialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LLOVerifierProxyVerifierInitialized)
				if err := _LLOVerifierProxy.contract.UnpackLog(event, "VerifierInitialized", log); err != nil {
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

func (_LLOVerifierProxy *LLOVerifierProxyFilterer) ParseVerifierInitialized(log types.Log) (*LLOVerifierProxyVerifierInitialized, error) {
	event := new(LLOVerifierProxyVerifierInitialized)
	if err := _LLOVerifierProxy.contract.UnpackLog(event, "VerifierInitialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LLOVerifierProxyVerifierSetIterator struct {
	Event *LLOVerifierProxyVerifierSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LLOVerifierProxyVerifierSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LLOVerifierProxyVerifierSet)
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
		it.Event = new(LLOVerifierProxyVerifierSet)
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

func (it *LLOVerifierProxyVerifierSetIterator) Error() error {
	return it.fail
}

func (it *LLOVerifierProxyVerifierSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LLOVerifierProxyVerifierSet struct {
	OldConfigDigest [32]byte
	NewConfigDigest [32]byte
	VerifierAddress common.Address
	Raw             types.Log
}

func (_LLOVerifierProxy *LLOVerifierProxyFilterer) FilterVerifierSet(opts *bind.FilterOpts) (*LLOVerifierProxyVerifierSetIterator, error) {

	logs, sub, err := _LLOVerifierProxy.contract.FilterLogs(opts, "VerifierSet")
	if err != nil {
		return nil, err
	}
	return &LLOVerifierProxyVerifierSetIterator{contract: _LLOVerifierProxy.contract, event: "VerifierSet", logs: logs, sub: sub}, nil
}

func (_LLOVerifierProxy *LLOVerifierProxyFilterer) WatchVerifierSet(opts *bind.WatchOpts, sink chan<- *LLOVerifierProxyVerifierSet) (event.Subscription, error) {

	logs, sub, err := _LLOVerifierProxy.contract.WatchLogs(opts, "VerifierSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LLOVerifierProxyVerifierSet)
				if err := _LLOVerifierProxy.contract.UnpackLog(event, "VerifierSet", log); err != nil {
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

func (_LLOVerifierProxy *LLOVerifierProxyFilterer) ParseVerifierSet(log types.Log) (*LLOVerifierProxyVerifierSet, error) {
	event := new(LLOVerifierProxyVerifierSet)
	if err := _LLOVerifierProxy.contract.UnpackLog(event, "VerifierSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LLOVerifierProxyVerifierUnsetIterator struct {
	Event *LLOVerifierProxyVerifierUnset

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LLOVerifierProxyVerifierUnsetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LLOVerifierProxyVerifierUnset)
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
		it.Event = new(LLOVerifierProxyVerifierUnset)
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

func (it *LLOVerifierProxyVerifierUnsetIterator) Error() error {
	return it.fail
}

func (it *LLOVerifierProxyVerifierUnsetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LLOVerifierProxyVerifierUnset struct {
	ConfigDigest    [32]byte
	VerifierAddress common.Address
	Raw             types.Log
}

func (_LLOVerifierProxy *LLOVerifierProxyFilterer) FilterVerifierUnset(opts *bind.FilterOpts) (*LLOVerifierProxyVerifierUnsetIterator, error) {

	logs, sub, err := _LLOVerifierProxy.contract.FilterLogs(opts, "VerifierUnset")
	if err != nil {
		return nil, err
	}
	return &LLOVerifierProxyVerifierUnsetIterator{contract: _LLOVerifierProxy.contract, event: "VerifierUnset", logs: logs, sub: sub}, nil
}

func (_LLOVerifierProxy *LLOVerifierProxyFilterer) WatchVerifierUnset(opts *bind.WatchOpts, sink chan<- *LLOVerifierProxyVerifierUnset) (event.Subscription, error) {

	logs, sub, err := _LLOVerifierProxy.contract.WatchLogs(opts, "VerifierUnset")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LLOVerifierProxyVerifierUnset)
				if err := _LLOVerifierProxy.contract.UnpackLog(event, "VerifierUnset", log); err != nil {
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

func (_LLOVerifierProxy *LLOVerifierProxyFilterer) ParseVerifierUnset(log types.Log) (*LLOVerifierProxyVerifierUnset, error) {
	event := new(LLOVerifierProxyVerifierUnset)
	if err := _LLOVerifierProxy.contract.UnpackLog(event, "VerifierUnset", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_LLOVerifierProxy *LLOVerifierProxy) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _LLOVerifierProxy.abi.Events["AccessControllerSet"].ID:
		return _LLOVerifierProxy.ParseAccessControllerSet(log)
	case _LLOVerifierProxy.abi.Events["OwnershipTransferRequested"].ID:
		return _LLOVerifierProxy.ParseOwnershipTransferRequested(log)
	case _LLOVerifierProxy.abi.Events["OwnershipTransferred"].ID:
		return _LLOVerifierProxy.ParseOwnershipTransferred(log)
	case _LLOVerifierProxy.abi.Events["VerifierInitialized"].ID:
		return _LLOVerifierProxy.ParseVerifierInitialized(log)
	case _LLOVerifierProxy.abi.Events["VerifierSet"].ID:
		return _LLOVerifierProxy.ParseVerifierSet(log)
	case _LLOVerifierProxy.abi.Events["VerifierUnset"].ID:
		return _LLOVerifierProxy.ParseVerifierUnset(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (LLOVerifierProxyAccessControllerSet) Topic() common.Hash {
	return common.HexToHash("0x953e92b1a6442e9c3242531154a3f6f6eb00b4e9c719ba8118fa6235e4ce89b6")
}

func (LLOVerifierProxyOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (LLOVerifierProxyOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (LLOVerifierProxyVerifierInitialized) Topic() common.Hash {
	return common.HexToHash("0x1f2cd7c97f4d801b5efe26cc409617c1fd6c5ef786e79aacb90af40923e4e8e9")
}

func (LLOVerifierProxyVerifierSet) Topic() common.Hash {
	return common.HexToHash("0xbeb513e532542a562ac35699e7cd9ae7d198dcd3eee15bada6c857d28ceaddcf")
}

func (LLOVerifierProxyVerifierUnset) Topic() common.Hash {
	return common.HexToHash("0x11dc15c4b8ac2b183166cc8427e5385a5ece8308217a4217338c6a7614845c4c")
}

func (_LLOVerifierProxy *LLOVerifierProxy) Address() common.Address {
	return _LLOVerifierProxy.address
}

type LLOVerifierProxyInterface interface {
	GetAccessController(opts *bind.CallOpts) (common.Address, error)

	GetVerifier(opts *bind.CallOpts, configDigest [32]byte) (common.Address, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	InitializeVerifier(opts *bind.TransactOpts, verifierAddress common.Address) (*types.Transaction, error)

	SetAccessController(opts *bind.TransactOpts, accessController common.Address) (*types.Transaction, error)

	SetVerifier(opts *bind.TransactOpts, currentConfigDigest [32]byte, newConfigDigest [32]byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	UnsetVerifier(opts *bind.TransactOpts, configDigest [32]byte) (*types.Transaction, error)

	Verify(opts *bind.TransactOpts, signedReport []byte) (*types.Transaction, error)

	FilterAccessControllerSet(opts *bind.FilterOpts) (*LLOVerifierProxyAccessControllerSetIterator, error)

	WatchAccessControllerSet(opts *bind.WatchOpts, sink chan<- *LLOVerifierProxyAccessControllerSet) (event.Subscription, error)

	ParseAccessControllerSet(log types.Log) (*LLOVerifierProxyAccessControllerSet, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LLOVerifierProxyOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *LLOVerifierProxyOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*LLOVerifierProxyOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LLOVerifierProxyOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *LLOVerifierProxyOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*LLOVerifierProxyOwnershipTransferred, error)

	FilterVerifierInitialized(opts *bind.FilterOpts) (*LLOVerifierProxyVerifierInitializedIterator, error)

	WatchVerifierInitialized(opts *bind.WatchOpts, sink chan<- *LLOVerifierProxyVerifierInitialized) (event.Subscription, error)

	ParseVerifierInitialized(log types.Log) (*LLOVerifierProxyVerifierInitialized, error)

	FilterVerifierSet(opts *bind.FilterOpts) (*LLOVerifierProxyVerifierSetIterator, error)

	WatchVerifierSet(opts *bind.WatchOpts, sink chan<- *LLOVerifierProxyVerifierSet) (event.Subscription, error)

	ParseVerifierSet(log types.Log) (*LLOVerifierProxyVerifierSet, error)

	FilterVerifierUnset(opts *bind.FilterOpts) (*LLOVerifierProxyVerifierUnsetIterator, error)

	WatchVerifierUnset(opts *bind.WatchOpts, sink chan<- *LLOVerifierProxyVerifierUnset) (event.Subscription, error)

	ParseVerifierUnset(log types.Log) (*LLOVerifierProxyVerifierUnset, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
