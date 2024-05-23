// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package forwarder

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

var KeystoneForwarderMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"}],\"name\":\"DuplicateSigner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxSigners\",\"type\":\"uint256\"}],\"name\":\"ExcessSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FaultToleranceMustBePositive\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minSigners\",\"type\":\"uint256\"}],\"name\":\"InsufficientSigners\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"}],\"name\":\"InvalidDonId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReport\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"InvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"received\",\"type\":\"uint256\"}],\"name\":\"InvalidSignatureCount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"}],\"name\":\"InvalidSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrantCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"reportId\",\"type\":\"bytes32\"}],\"name\":\"ReportAlreadyProcessed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"workflowExecutionId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"reportName\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"workflowOwner\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"workflowName\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"result\",\"type\":\"bool\"}],\"name\":\"ReportProcessed\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"workflowExecutionId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"reportName\",\"type\":\"bytes32\"}],\"name\":\"getTransmitter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiverAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"rawReport\",\"type\":\"bytes\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"report\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5033806000816100675760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615610097576100978161009f565b505050610148565b336001600160a01b038216036100f75760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161005e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b61165e806101576000396000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c806379ba50971161005b57806379ba5097146101705780638da5cb5b14610178578063c0965dc314610196578063f2fde38b146101a957600080fd5b8063134a46f014610082578063181f5a771461009757806330670580146100df575b600080fd5b610095610090366004611246565b6101bc565b005b604080518082018252601781527f4b657973746f6e65466f7277617264657220312e302e30000000000000000000602082015290516100d6919061131e565b60405180910390f35b61014b6100ed366004611361565b6040805173ffffffffffffffffffffffffffffffffffffffff94851660208083019190915281830194909452606080820193909352815180820390930183526080018152815191830191909120600090815260039092529020541690565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100d6565b6100956105d6565b60005473ffffffffffffffffffffffffffffffffffffffff1661014b565b6100956101a4366004611394565b6106d3565b6100956101b7366004611443565b610ecf565b60015474010000000000000000000000000000000000000000900460ff1615610211576040517f37ed32e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600180547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000017905560ff831660000361028d576040517f0743bae600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601f8111156102d7576040517f61750f4000000000000000000000000000000000000000000000000000000000815260048101829052601f60248201526044015b60405180910390fd5b6102e283600361148d565b60ff16811161034057806102f784600361148d565b6103029060016114b0565b6040517f9dd9e6d8000000000000000000000000000000000000000000000000000000008152600481019290925260ff1660248201526044016102ce565b60005b63ffffffff85166000908152600260205260409020600101548110156103e15763ffffffff8516600090815260026020526040812060010180548390811061038d5761038d6114c9565b600091825260208083209091015463ffffffff891683526002808352604080852073ffffffffffffffffffffffffffffffffffffffff90931685529101909152812055506103da816114f8565b9050610343565b5063ffffffff8416600090815260026020526040902061040590600101838361113e565b5060005b81811015610563576000838383818110610425576104256114c9565b905060200201602081019061043a9190611443565b63ffffffff8716600090815260026020818152604080842073ffffffffffffffffffffffffffffffffffffffff861685529092019052902054909150156104c5576040517fe021c4f200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821660048201526024016102ce565b6104d08260016114b0565b63ffffffff8716600090815260026020818152604080842073ffffffffffffffffffffffffffffffffffffffff909616808552868401835290842060ff959095169094559081526001938401805494850181558252902090910180547fffffffffffffffffffffffff000000000000000000000000000000000000000016909117905561055c816114f8565b9050610409565b50505063ffffffff91909116600090815260026020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff909216919091179055600180547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff169055565b60015473ffffffffffffffffffffffffffffffffffffffff163314610657576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016102ce565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60015474010000000000000000000000000000000000000000900460ff1615610728576040517f37ed32e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600180547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000017905560a48310156107a2576040517fb55ac75400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000806000806000806107ea8a8a8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250610ee392505050565b63ffffffff8516600090815260026020526040812054969c50949a509298509096509450925060ff9091169003610855576040517fea1b312900000000000000000000000000000000000000000000000000000000815263ffffffff861660048201526024016102ce565b6040805173ffffffffffffffffffffffffffffffffffffffff8d81166020808401919091528284018890526060808401869052845180850390910181526080909301845282519281019290922060009081526003909252919020541615610934576040805173ffffffffffffffffffffffffffffffffffffffff8d16602080830191909152818301879052606080830185905283518084039091018152608090920190925280519101206040517f1aac3d290000000000000000000000000000000000000000000000000000000081526004016102ce91815260200190565b63ffffffff8516600090815260026020526040902054879061095a9060ff1660016114b0565b60ff16146109c55763ffffffff85166000908152600260205260409020546109869060ff1660016114b0565b6040517fd6022e8e00000000000000000000000000000000000000000000000000000000815260ff9091166004820152602481018890526044016102ce565b60008a8a6040516109d7929190611530565b604051809103902090506109e96111c6565b6000805b8a811015610bc0576000610a59858e8e85818110610a0d57610a0d6114c9565b9050602002810190610a1f9190611540565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250610f1192505050565b63ffffffff8b16600090815260026020818152604080842073ffffffffffffffffffffffffffffffffffffffff861685529092019052812054945090915060ff84169003610aeb576040517fbf18af4300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821660048201526024016102ce565b610af66001846115a5565b925060008460ff8516601f8110610b0f57610b0f6114c9565b602002015173ffffffffffffffffffffffffffffffffffffffff1614610b79576040517fe021c4f200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821660048201526024016102ce565b80848460ff16601f8110610b8f57610b8f6114c9565b73ffffffffffffffffffffffffffffffffffffffff909216602092909202015250610bb9816114f8565b90506109ed565b505050508a73ffffffffffffffffffffffffffffffffffffffff16630d5c96c7878585858f8f60a4908092610bf7939291906115be565b6040518763ffffffff1660e01b8152600401610c18969594939291906115e8565b600060405180830381600087803b158015610c3257600080fd5b505af1925050508015610c43575060015b610d7a5760408051808201909152338152600060208201819052600390610cb88e88866040805173ffffffffffffffffffffffffffffffffffffffff85166020820152908101839052606081018290526000906080016040516020818303038152906040528051906020012090509392505050565b8152602080820192909252604090810160009081208451815495850151151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00000000000000000000000000000000000000000090961673ffffffffffffffffffffffffffffffffffffffff9182161795909517905581518781529283018690529082015282918691908e16907f5df27a37f78ff317059466ada9bcfe7dfa6f5ebc8b0b1ab36a52eb3f47a7139e9060600160405180910390a4610e9a565b6040805180820182523381526001602080830191909152825173ffffffffffffffffffffffffffffffffffffffff8f168183015280840188905260608082018690528451808303909101815260809091019093528251920191909120600390600090815260208082019290925260409081016000208351815494840151151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00000000000000000000000000000000000000000090951673ffffffffffffffffffffffffffffffffffffffff91821617949094179055805186815291820185905260019082015282918691908e16907f5df27a37f78ff317059466ada9bcfe7dfa6f5ebc8b0b1ab36a52eb3f47a7139e9060600160405180910390a45b5050600180547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff169055505050505050505050565b610ed7610fc6565b610ee081611049565b50565b6020810151604082015160448301516064840151608485015160a490950151939560e09390931c9491939092565b60006041825114610f5057816040517f2adfdc300000000000000000000000000000000000000000000000000000000081526004016102ce919061131e565b602082810151604080850151606080870151835160008082529681018086528a9052951a928501839052840183905260808401819052919260019060a0016020604051602081039080840390855afa158015610fb0573d6000803e3d6000fd5b5050506020604051035193505050505b92915050565b60005473ffffffffffffffffffffffffffffffffffffffff163314611047576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016102ce565b565b3373ffffffffffffffffffffffffffffffffffffffff8216036110c8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016102ce565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b8280548282559060005260206000209081019282156111b6579160200282015b828111156111b65781547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff84351617825560209092019160019091019061115e565b506111c29291506111e5565b5090565b604051806103e00160405280601f906020820280368337509192915050565b5b808211156111c257600081556001016111e6565b60008083601f84011261120c57600080fd5b50813567ffffffffffffffff81111561122457600080fd5b6020830191508360208260051b850101111561123f57600080fd5b9250929050565b6000806000806060858703121561125c57600080fd5b843563ffffffff8116811461127057600080fd5b9350602085013560ff8116811461128657600080fd5b9250604085013567ffffffffffffffff8111156112a257600080fd5b6112ae878288016111fa565b95989497509550505050565b6000815180845260005b818110156112e0576020818501810151868301820152016112c4565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b60208152600061133160208301846112ba565b9392505050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461135c57600080fd5b919050565b60008060006060848603121561137657600080fd5b61137f84611338565b95602085013595506040909401359392505050565b6000806000806000606086880312156113ac57600080fd5b6113b586611338565b9450602086013567ffffffffffffffff808211156113d257600080fd5b818801915088601f8301126113e657600080fd5b8135818111156113f557600080fd5b89602082850101111561140757600080fd5b60208301965080955050604088013591508082111561142557600080fd5b50611432888289016111fa565b969995985093965092949392505050565b60006020828403121561145557600080fd5b61133182611338565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60ff81811683821602908116908181146114a9576114a961145e565b5092915050565b60ff8181168382160190811115610fc057610fc061145e565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036115295761152961145e565b5060010190565b8183823760009101908152919050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261157557600080fd5b83018035915067ffffffffffffffff82111561159057600080fd5b60200191503681900382131561123f57600080fd5b60ff8281168282160390811115610fc057610fc061145e565b600080858511156115ce57600080fd5b838611156115db57600080fd5b5050820193919092039150565b86815285602082015284604082015283606082015260a060808201528160a0820152818360c0830137600081830160c090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01601019594505050505056fea164736f6c6343000813000a",
}

var KeystoneForwarderABI = KeystoneForwarderMetaData.ABI

var KeystoneForwarderBin = KeystoneForwarderMetaData.Bin

func DeployKeystoneForwarder(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *KeystoneForwarder, error) {
	parsed, err := KeystoneForwarderMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KeystoneForwarderBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KeystoneForwarder{address: address, abi: *parsed, KeystoneForwarderCaller: KeystoneForwarderCaller{contract: contract}, KeystoneForwarderTransactor: KeystoneForwarderTransactor{contract: contract}, KeystoneForwarderFilterer: KeystoneForwarderFilterer{contract: contract}}, nil
}

type KeystoneForwarder struct {
	address common.Address
	abi     abi.ABI
	KeystoneForwarderCaller
	KeystoneForwarderTransactor
	KeystoneForwarderFilterer
}

type KeystoneForwarderCaller struct {
	contract *bind.BoundContract
}

type KeystoneForwarderTransactor struct {
	contract *bind.BoundContract
}

type KeystoneForwarderFilterer struct {
	contract *bind.BoundContract
}

type KeystoneForwarderSession struct {
	Contract     *KeystoneForwarder
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type KeystoneForwarderCallerSession struct {
	Contract *KeystoneForwarderCaller
	CallOpts bind.CallOpts
}

type KeystoneForwarderTransactorSession struct {
	Contract     *KeystoneForwarderTransactor
	TransactOpts bind.TransactOpts
}

type KeystoneForwarderRaw struct {
	Contract *KeystoneForwarder
}

type KeystoneForwarderCallerRaw struct {
	Contract *KeystoneForwarderCaller
}

type KeystoneForwarderTransactorRaw struct {
	Contract *KeystoneForwarderTransactor
}

func NewKeystoneForwarder(address common.Address, backend bind.ContractBackend) (*KeystoneForwarder, error) {
	abi, err := abi.JSON(strings.NewReader(KeystoneForwarderABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindKeystoneForwarder(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeystoneForwarder{address: address, abi: abi, KeystoneForwarderCaller: KeystoneForwarderCaller{contract: contract}, KeystoneForwarderTransactor: KeystoneForwarderTransactor{contract: contract}, KeystoneForwarderFilterer: KeystoneForwarderFilterer{contract: contract}}, nil
}

func NewKeystoneForwarderCaller(address common.Address, caller bind.ContractCaller) (*KeystoneForwarderCaller, error) {
	contract, err := bindKeystoneForwarder(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeystoneForwarderCaller{contract: contract}, nil
}

func NewKeystoneForwarderTransactor(address common.Address, transactor bind.ContractTransactor) (*KeystoneForwarderTransactor, error) {
	contract, err := bindKeystoneForwarder(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeystoneForwarderTransactor{contract: contract}, nil
}

func NewKeystoneForwarderFilterer(address common.Address, filterer bind.ContractFilterer) (*KeystoneForwarderFilterer, error) {
	contract, err := bindKeystoneForwarder(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeystoneForwarderFilterer{contract: contract}, nil
}

func bindKeystoneForwarder(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := KeystoneForwarderMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_KeystoneForwarder *KeystoneForwarderRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeystoneForwarder.Contract.KeystoneForwarderCaller.contract.Call(opts, result, method, params...)
}

func (_KeystoneForwarder *KeystoneForwarderRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.KeystoneForwarderTransactor.contract.Transfer(opts)
}

func (_KeystoneForwarder *KeystoneForwarderRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.KeystoneForwarderTransactor.contract.Transact(opts, method, params...)
}

func (_KeystoneForwarder *KeystoneForwarderCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeystoneForwarder.Contract.contract.Call(opts, result, method, params...)
}

func (_KeystoneForwarder *KeystoneForwarderTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.contract.Transfer(opts)
}

func (_KeystoneForwarder *KeystoneForwarderTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.contract.Transact(opts, method, params...)
}

func (_KeystoneForwarder *KeystoneForwarderCaller) GetTransmitter(opts *bind.CallOpts, receiver common.Address, workflowExecutionId [32]byte, reportName [32]byte) (common.Address, error) {
	var out []interface{}
	err := _KeystoneForwarder.contract.Call(opts, &out, "getTransmitter", receiver, workflowExecutionId, reportName)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeystoneForwarder *KeystoneForwarderSession) GetTransmitter(receiver common.Address, workflowExecutionId [32]byte, reportName [32]byte) (common.Address, error) {
	return _KeystoneForwarder.Contract.GetTransmitter(&_KeystoneForwarder.CallOpts, receiver, workflowExecutionId, reportName)
}

func (_KeystoneForwarder *KeystoneForwarderCallerSession) GetTransmitter(receiver common.Address, workflowExecutionId [32]byte, reportName [32]byte) (common.Address, error) {
	return _KeystoneForwarder.Contract.GetTransmitter(&_KeystoneForwarder.CallOpts, receiver, workflowExecutionId, reportName)
}

func (_KeystoneForwarder *KeystoneForwarderCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeystoneForwarder.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeystoneForwarder *KeystoneForwarderSession) Owner() (common.Address, error) {
	return _KeystoneForwarder.Contract.Owner(&_KeystoneForwarder.CallOpts)
}

func (_KeystoneForwarder *KeystoneForwarderCallerSession) Owner() (common.Address, error) {
	return _KeystoneForwarder.Contract.Owner(&_KeystoneForwarder.CallOpts)
}

func (_KeystoneForwarder *KeystoneForwarderCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _KeystoneForwarder.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_KeystoneForwarder *KeystoneForwarderSession) TypeAndVersion() (string, error) {
	return _KeystoneForwarder.Contract.TypeAndVersion(&_KeystoneForwarder.CallOpts)
}

func (_KeystoneForwarder *KeystoneForwarderCallerSession) TypeAndVersion() (string, error) {
	return _KeystoneForwarder.Contract.TypeAndVersion(&_KeystoneForwarder.CallOpts)
}

func (_KeystoneForwarder *KeystoneForwarderTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeystoneForwarder.contract.Transact(opts, "acceptOwnership")
}

func (_KeystoneForwarder *KeystoneForwarderSession) AcceptOwnership() (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.AcceptOwnership(&_KeystoneForwarder.TransactOpts)
}

func (_KeystoneForwarder *KeystoneForwarderTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.AcceptOwnership(&_KeystoneForwarder.TransactOpts)
}

func (_KeystoneForwarder *KeystoneForwarderTransactor) Report(opts *bind.TransactOpts, receiverAddress common.Address, rawReport []byte, signatures [][]byte) (*types.Transaction, error) {
	return _KeystoneForwarder.contract.Transact(opts, "report", receiverAddress, rawReport, signatures)
}

func (_KeystoneForwarder *KeystoneForwarderSession) Report(receiverAddress common.Address, rawReport []byte, signatures [][]byte) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.Report(&_KeystoneForwarder.TransactOpts, receiverAddress, rawReport, signatures)
}

func (_KeystoneForwarder *KeystoneForwarderTransactorSession) Report(receiverAddress common.Address, rawReport []byte, signatures [][]byte) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.Report(&_KeystoneForwarder.TransactOpts, receiverAddress, rawReport, signatures)
}

func (_KeystoneForwarder *KeystoneForwarderTransactor) SetConfig(opts *bind.TransactOpts, donId uint32, f uint8, signers []common.Address) (*types.Transaction, error) {
	return _KeystoneForwarder.contract.Transact(opts, "setConfig", donId, f, signers)
}

func (_KeystoneForwarder *KeystoneForwarderSession) SetConfig(donId uint32, f uint8, signers []common.Address) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.SetConfig(&_KeystoneForwarder.TransactOpts, donId, f, signers)
}

func (_KeystoneForwarder *KeystoneForwarderTransactorSession) SetConfig(donId uint32, f uint8, signers []common.Address) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.SetConfig(&_KeystoneForwarder.TransactOpts, donId, f, signers)
}

func (_KeystoneForwarder *KeystoneForwarderTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _KeystoneForwarder.contract.Transact(opts, "transferOwnership", to)
}

func (_KeystoneForwarder *KeystoneForwarderSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.TransferOwnership(&_KeystoneForwarder.TransactOpts, to)
}

func (_KeystoneForwarder *KeystoneForwarderTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.TransferOwnership(&_KeystoneForwarder.TransactOpts, to)
}

type KeystoneForwarderOwnershipTransferRequestedIterator struct {
	Event *KeystoneForwarderOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeystoneForwarderOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeystoneForwarderOwnershipTransferRequested)
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
		it.Event = new(KeystoneForwarderOwnershipTransferRequested)
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

func (it *KeystoneForwarderOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeystoneForwarderOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeystoneForwarderOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeystoneForwarder *KeystoneForwarderFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeystoneForwarderOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeystoneForwarder.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeystoneForwarderOwnershipTransferRequestedIterator{contract: _KeystoneForwarder.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeystoneForwarder *KeystoneForwarderFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeystoneForwarderOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeystoneForwarder.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeystoneForwarderOwnershipTransferRequested)
				if err := _KeystoneForwarder.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_KeystoneForwarder *KeystoneForwarderFilterer) ParseOwnershipTransferRequested(log types.Log) (*KeystoneForwarderOwnershipTransferRequested, error) {
	event := new(KeystoneForwarderOwnershipTransferRequested)
	if err := _KeystoneForwarder.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeystoneForwarderOwnershipTransferredIterator struct {
	Event *KeystoneForwarderOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeystoneForwarderOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeystoneForwarderOwnershipTransferred)
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
		it.Event = new(KeystoneForwarderOwnershipTransferred)
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

func (it *KeystoneForwarderOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *KeystoneForwarderOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeystoneForwarderOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeystoneForwarder *KeystoneForwarderFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeystoneForwarderOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeystoneForwarder.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeystoneForwarderOwnershipTransferredIterator{contract: _KeystoneForwarder.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_KeystoneForwarder *KeystoneForwarderFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeystoneForwarderOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeystoneForwarder.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeystoneForwarderOwnershipTransferred)
				if err := _KeystoneForwarder.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_KeystoneForwarder *KeystoneForwarderFilterer) ParseOwnershipTransferred(log types.Log) (*KeystoneForwarderOwnershipTransferred, error) {
	event := new(KeystoneForwarderOwnershipTransferred)
	if err := _KeystoneForwarder.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeystoneForwarderReportProcessedIterator struct {
	Event *KeystoneForwarderReportProcessed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeystoneForwarderReportProcessedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeystoneForwarderReportProcessed)
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
		it.Event = new(KeystoneForwarderReportProcessed)
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

func (it *KeystoneForwarderReportProcessedIterator) Error() error {
	return it.fail
}

func (it *KeystoneForwarderReportProcessedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeystoneForwarderReportProcessed struct {
	Receiver            common.Address
	WorkflowExecutionId [32]byte
	ReportName          [32]byte
	WorkflowOwner       [32]byte
	WorkflowName        [32]byte
	Result              bool
	Raw                 types.Log
}

func (_KeystoneForwarder *KeystoneForwarderFilterer) FilterReportProcessed(opts *bind.FilterOpts, receiver []common.Address, workflowExecutionId [][32]byte, reportName [][32]byte) (*KeystoneForwarderReportProcessedIterator, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}
	var workflowExecutionIdRule []interface{}
	for _, workflowExecutionIdItem := range workflowExecutionId {
		workflowExecutionIdRule = append(workflowExecutionIdRule, workflowExecutionIdItem)
	}
	var reportNameRule []interface{}
	for _, reportNameItem := range reportName {
		reportNameRule = append(reportNameRule, reportNameItem)
	}

	logs, sub, err := _KeystoneForwarder.contract.FilterLogs(opts, "ReportProcessed", receiverRule, workflowExecutionIdRule, reportNameRule)
	if err != nil {
		return nil, err
	}
	return &KeystoneForwarderReportProcessedIterator{contract: _KeystoneForwarder.contract, event: "ReportProcessed", logs: logs, sub: sub}, nil
}

func (_KeystoneForwarder *KeystoneForwarderFilterer) WatchReportProcessed(opts *bind.WatchOpts, sink chan<- *KeystoneForwarderReportProcessed, receiver []common.Address, workflowExecutionId [][32]byte, reportName [][32]byte) (event.Subscription, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}
	var workflowExecutionIdRule []interface{}
	for _, workflowExecutionIdItem := range workflowExecutionId {
		workflowExecutionIdRule = append(workflowExecutionIdRule, workflowExecutionIdItem)
	}
	var reportNameRule []interface{}
	for _, reportNameItem := range reportName {
		reportNameRule = append(reportNameRule, reportNameItem)
	}

	logs, sub, err := _KeystoneForwarder.contract.WatchLogs(opts, "ReportProcessed", receiverRule, workflowExecutionIdRule, reportNameRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeystoneForwarderReportProcessed)
				if err := _KeystoneForwarder.contract.UnpackLog(event, "ReportProcessed", log); err != nil {
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

func (_KeystoneForwarder *KeystoneForwarderFilterer) ParseReportProcessed(log types.Log) (*KeystoneForwarderReportProcessed, error) {
	event := new(KeystoneForwarderReportProcessed)
	if err := _KeystoneForwarder.contract.UnpackLog(event, "ReportProcessed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_KeystoneForwarder *KeystoneForwarder) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _KeystoneForwarder.abi.Events["OwnershipTransferRequested"].ID:
		return _KeystoneForwarder.ParseOwnershipTransferRequested(log)
	case _KeystoneForwarder.abi.Events["OwnershipTransferred"].ID:
		return _KeystoneForwarder.ParseOwnershipTransferred(log)
	case _KeystoneForwarder.abi.Events["ReportProcessed"].ID:
		return _KeystoneForwarder.ParseReportProcessed(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (KeystoneForwarderOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (KeystoneForwarderOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (KeystoneForwarderReportProcessed) Topic() common.Hash {
	return common.HexToHash("0x5df27a37f78ff317059466ada9bcfe7dfa6f5ebc8b0b1ab36a52eb3f47a7139e")
}

func (_KeystoneForwarder *KeystoneForwarder) Address() common.Address {
	return _KeystoneForwarder.address
}

type KeystoneForwarderInterface interface {
	GetTransmitter(opts *bind.CallOpts, receiver common.Address, workflowExecutionId [32]byte, reportName [32]byte) (common.Address, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	Report(opts *bind.TransactOpts, receiverAddress common.Address, rawReport []byte, signatures [][]byte) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, donId uint32, f uint8, signers []common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeystoneForwarderOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeystoneForwarderOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*KeystoneForwarderOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeystoneForwarderOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeystoneForwarderOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*KeystoneForwarderOwnershipTransferred, error)

	FilterReportProcessed(opts *bind.FilterOpts, receiver []common.Address, workflowExecutionId [][32]byte, reportName [][32]byte) (*KeystoneForwarderReportProcessedIterator, error)

	WatchReportProcessed(opts *bind.WatchOpts, sink chan<- *KeystoneForwarderReportProcessed, receiver []common.Address, workflowExecutionId [][32]byte, reportName [][32]byte) (event.Subscription, error)

	ParseReportProcessed(log types.Log) (*KeystoneForwarderReportProcessed, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
