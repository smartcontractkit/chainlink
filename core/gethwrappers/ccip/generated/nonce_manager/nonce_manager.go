// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package nonce_manager

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

type AuthorizedCallersAuthorizedCallerArgs struct {
	AddedCallers   []common.Address
	RemovedCallers []common.Address
}

type NonceManagerPreviousRamps struct {
	PrevOnRamp common.Address
}

type NonceManagerPreviousRampsArgs struct {
	RemotChainSelector uint64
	PrevRamps          NonceManagerPreviousRamps
}

var NonceManagerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"authorizedCallers\",\"type\":\"address[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"PreviousRampAlreadySet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"UnauthorizedCaller\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"AuthorizedCallerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"AuthorizedCallerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"prevOnRamp\",\"type\":\"address\"}],\"name\":\"PreviousOnRampUpdated\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address[]\",\"name\":\"addedCallers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"removedCallers\",\"type\":\"address[]\"}],\"internalType\":\"structAuthorizedCallers.AuthorizedCallerArgs\",\"name\":\"authorizedCallerArgs\",\"type\":\"tuple\"}],\"name\":\"applyAuthorizedCallerUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"remotChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"prevOnRamp\",\"type\":\"address\"}],\"internalType\":\"structNonceManager.PreviousRamps\",\"name\":\"prevRamps\",\"type\":\"tuple\"}],\"internalType\":\"structNonceManager.PreviousRampsArgs[]\",\"name\":\"previousRampsArgs\",\"type\":\"tuple[]\"}],\"name\":\"applyPreviousRampsUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllAuthorizedCallers\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"}],\"name\":\"getIncrementedOutboundNonce\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"}],\"name\":\"getOutboundNonce\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"}],\"name\":\"getPreviousRamps\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"prevOnRamp\",\"type\":\"address\"}],\"internalType\":\"structNonceManager.PreviousRamps\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162001659380380620016598339810160408190526200003491620004b0565b8033806000816200008c5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000bf57620000bf81620000f6565b5050604080518082018252838152815160008152602080820190935291810191909152620000ee9150620001a1565b5050620005d0565b336001600160a01b03821603620001505760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000083565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b602081015160005b815181101562000231576000828281518110620001ca57620001ca62000582565b60209081029190910101519050620001e4600282620002f0565b1562000227576040516001600160a01b03821681527fc3803387881faad271c47728894e3e36fac830ffc8602ca6fc07733cbda775809060200160405180910390a15b50600101620001a9565b50815160005b8151811015620002ea57600082828151811062000258576200025862000582565b6020026020010151905060006001600160a01b0316816001600160a01b03160362000296576040516342bcdf7f60e11b815260040160405180910390fd5b620002a360028262000310565b506040516001600160a01b03821681527feb1b9b92e50b7f88f9ff25d56765095ac6e91540eee214906f4036a908ffbdef9060200160405180910390a15060010162000237565b50505050565b600062000307836001600160a01b03841662000327565b90505b92915050565b600062000307836001600160a01b0384166200042b565b60008181526001830160205260408120548015620004205760006200034e60018362000598565b8554909150600090620003649060019062000598565b9050818114620003d057600086600001828154811062000388576200038862000582565b9060005260206000200154905080876000018481548110620003ae57620003ae62000582565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080620003e457620003e4620005ba565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506200030a565b60009150506200030a565b600081815260018301602052604081205462000474575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556200030a565b5060006200030a565b634e487b7160e01b600052604160045260246000fd5b80516001600160a01b0381168114620004ab57600080fd5b919050565b60006020808385031215620004c457600080fd5b82516001600160401b0380821115620004dc57600080fd5b818501915085601f830112620004f157600080fd5b8151818111156200050657620005066200047d565b8060051b604051601f19603f830116810181811085821117156200052e576200052e6200047d565b6040529182528482019250838101850191888311156200054d57600080fd5b938501935b828510156200057657620005668562000493565b8452938501939285019262000552565b98975050505050505050565b634e487b7160e01b600052603260045260246000fd5b818103818111156200030a57634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052603160045260246000fd5b61107980620005e06000396000f3fe608060405234801561001057600080fd5b50600436106100a35760003560e01c806379ba50971161007657806391a2749a1161005b57806391a2749a146101b5578063d18be31b146101c8578063f2fde38b146101db57600080fd5b806379ba5097146101855780638da5cb5b1461018d57600080fd5b80631ce2b142146100a85780632451a627146100bd578063294b5630146100db57806331b89ff314610159575b600080fd5b6100bb6100b6366004610c1d565b6101ee565b005b6100c5610363565b6040516100d29190610c92565b60405180910390f35b6101346100e9366004610d02565b604080516020808201835260009182905267ffffffffffffffff9390931681526004835281902081519283019091525473ffffffffffffffffffffffffffffffffffffffff16815290565b604051905173ffffffffffffffffffffffffffffffffffffffff1681526020016100d2565b61016c610167366004610d1f565b610374565b60405167ffffffffffffffff90911681526020016100d2565b6100bb61038b565b60005460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100d2565b6100bb6101c3366004610eba565b61048d565b61016c6101d6366004610d1f565b6104a1565b6100bb6101e9366004610f61565b610543565b6101f6610554565b60005b8181101561035e573683838381811061021457610214610f7e565b604002919091019150600090506004816102316020850185610d02565b67ffffffffffffffff1681526020810191909152604001600020805490915073ffffffffffffffffffffffffffffffffffffffff161561029d576040517fc6117ae200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6102ad6040830160208401610f61565b81547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff919091161781556102fa6020830183610d02565b815460405173ffffffffffffffffffffffffffffffffffffffff909116815267ffffffffffffffff91909116907f89d2355e2829b1e15855fec87fb400638aebc9f03728949d702d3b5d4ea999549060200160405180910390a250506001016101f9565b505050565b606061036f60026105d7565b905090565b60006103818484846105e4565b90505b9392505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610411576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610495610554565b61049e81610732565b50565b60006104ab6108c4565b60006104b88585856105e4565b6104c3906001610fdc565b67ffffffffffffffff861660009081526005602052604090819020905191925082916104f29087908790610ffd565b908152604051908190036020019020805467ffffffffffffffff929092167fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000090921691909117905590509392505050565b61054b610554565b61049e81610907565b60005473ffffffffffffffffffffffffffffffffffffffff1633146105d5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610408565b565b60606000610384836109fc565b67ffffffffffffffff831660009081526005602052604080822090518291906106109086908690610ffd565b9081526040519081900360200190205467ffffffffffffffff16905060008190036103815767ffffffffffffffff851660009081526004602052604090205473ffffffffffffffffffffffffffffffffffffffff1680156107295773ffffffffffffffffffffffffffffffffffffffff811663856c824761069386880188610f61565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e084901b16815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602401602060405180830381865afa1580156106fc573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610720919061100d565b92505050610384565b50949350505050565b602081015160005b81518110156107cd57600082828151811061075757610757610f7e565b60200260200101519050610775816002610a5890919063ffffffff16565b156107c45760405173ffffffffffffffffffffffffffffffffffffffff821681527fc3803387881faad271c47728894e3e36fac830ffc8602ca6fc07733cbda775809060200160405180910390a15b5060010161073a565b50815160005b81518110156108be5760008282815181106107f0576107f0610f7e565b60200260200101519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610860576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61086b600282610a83565b5060405173ffffffffffffffffffffffffffffffffffffffff821681527feb1b9b92e50b7f88f9ff25d56765095ac6e91540eee214906f4036a908ffbdef9060200160405180910390a1506001016107d3565b50505050565b6108cf600233610aa5565b6105d5576040517fd86ad9cf000000000000000000000000000000000000000000000000000000008152336004820152602401610408565b3373ffffffffffffffffffffffffffffffffffffffff821603610986576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610408565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b606081600001805480602002602001604051908101604052809291908181526020018280548015610a4c57602002820191906000526020600020905b815481526020019060010190808311610a38575b50505050509050919050565b6000610a7a8373ffffffffffffffffffffffffffffffffffffffff8416610ad4565b90505b92915050565b6000610a7a8373ffffffffffffffffffffffffffffffffffffffff8416610bce565b73ffffffffffffffffffffffffffffffffffffffff811660009081526001830160205260408120541515610a7a565b60008181526001830160205260408120548015610bbd576000610af860018361102a565b8554909150600090610b0c9060019061102a565b9050818114610b71576000866000018281548110610b2c57610b2c610f7e565b9060005260206000200154905080876000018481548110610b4f57610b4f610f7e565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080610b8257610b8261103d565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050610a7d565b6000915050610a7d565b5092915050565b6000818152600183016020526040812054610c1557508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610a7d565b506000610a7d565b60008060208385031215610c3057600080fd5b823567ffffffffffffffff80821115610c4857600080fd5b818501915085601f830112610c5c57600080fd5b813581811115610c6b57600080fd5b8660208260061b8501011115610c8057600080fd5b60209290920196919550909350505050565b6020808252825182820181905260009190848201906040850190845b81811015610ce057835173ffffffffffffffffffffffffffffffffffffffff1683529284019291840191600101610cae565b50909695505050505050565b67ffffffffffffffff8116811461049e57600080fd5b600060208284031215610d1457600080fd5b813561038481610cec565b600080600060408486031215610d3457600080fd5b8335610d3f81610cec565b9250602084013567ffffffffffffffff80821115610d5c57600080fd5b818601915086601f830112610d7057600080fd5b813581811115610d7f57600080fd5b876020828501011115610d9157600080fd5b6020830194508093505050509250925092565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b73ffffffffffffffffffffffffffffffffffffffff8116811461049e57600080fd5b600082601f830112610e0657600080fd5b8135602067ffffffffffffffff80831115610e2357610e23610da4565b8260051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f83011681018181108482111715610e6657610e66610da4565b6040529384526020818701810194908101925087851115610e8657600080fd5b6020870191505b84821015610eaf578135610ea081610dd3565b83529183019190830190610e8d565b979650505050505050565b600060208284031215610ecc57600080fd5b813567ffffffffffffffff80821115610ee457600080fd5b9083019060408286031215610ef857600080fd5b604051604081018181108382111715610f1357610f13610da4565b604052823582811115610f2557600080fd5b610f3187828601610df5565b825250602083013582811115610f4657600080fd5b610f5287828601610df5565b60208301525095945050505050565b600060208284031215610f7357600080fd5b813561038481610dd3565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b67ffffffffffffffff818116838216019080821115610bc757610bc7610fad565b8183823760009101908152919050565b60006020828403121561101f57600080fd5b815161038481610cec565b81810381811115610a7d57610a7d610fad565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea164736f6c6343000818000a",
}

var NonceManagerABI = NonceManagerMetaData.ABI

var NonceManagerBin = NonceManagerMetaData.Bin

func DeployNonceManager(auth *bind.TransactOpts, backend bind.ContractBackend, authorizedCallers []common.Address) (common.Address, *types.Transaction, *NonceManager, error) {
	parsed, err := NonceManagerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(NonceManagerBin), backend, authorizedCallers)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &NonceManager{address: address, abi: *parsed, NonceManagerCaller: NonceManagerCaller{contract: contract}, NonceManagerTransactor: NonceManagerTransactor{contract: contract}, NonceManagerFilterer: NonceManagerFilterer{contract: contract}}, nil
}

type NonceManager struct {
	address common.Address
	abi     abi.ABI
	NonceManagerCaller
	NonceManagerTransactor
	NonceManagerFilterer
}

type NonceManagerCaller struct {
	contract *bind.BoundContract
}

type NonceManagerTransactor struct {
	contract *bind.BoundContract
}

type NonceManagerFilterer struct {
	contract *bind.BoundContract
}

type NonceManagerSession struct {
	Contract     *NonceManager
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type NonceManagerCallerSession struct {
	Contract *NonceManagerCaller
	CallOpts bind.CallOpts
}

type NonceManagerTransactorSession struct {
	Contract     *NonceManagerTransactor
	TransactOpts bind.TransactOpts
}

type NonceManagerRaw struct {
	Contract *NonceManager
}

type NonceManagerCallerRaw struct {
	Contract *NonceManagerCaller
}

type NonceManagerTransactorRaw struct {
	Contract *NonceManagerTransactor
}

func NewNonceManager(address common.Address, backend bind.ContractBackend) (*NonceManager, error) {
	abi, err := abi.JSON(strings.NewReader(NonceManagerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindNonceManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NonceManager{address: address, abi: abi, NonceManagerCaller: NonceManagerCaller{contract: contract}, NonceManagerTransactor: NonceManagerTransactor{contract: contract}, NonceManagerFilterer: NonceManagerFilterer{contract: contract}}, nil
}

func NewNonceManagerCaller(address common.Address, caller bind.ContractCaller) (*NonceManagerCaller, error) {
	contract, err := bindNonceManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NonceManagerCaller{contract: contract}, nil
}

func NewNonceManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*NonceManagerTransactor, error) {
	contract, err := bindNonceManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NonceManagerTransactor{contract: contract}, nil
}

func NewNonceManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*NonceManagerFilterer, error) {
	contract, err := bindNonceManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NonceManagerFilterer{contract: contract}, nil
}

func bindNonceManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := NonceManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_NonceManager *NonceManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NonceManager.Contract.NonceManagerCaller.contract.Call(opts, result, method, params...)
}

func (_NonceManager *NonceManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NonceManager.Contract.NonceManagerTransactor.contract.Transfer(opts)
}

func (_NonceManager *NonceManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NonceManager.Contract.NonceManagerTransactor.contract.Transact(opts, method, params...)
}

func (_NonceManager *NonceManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NonceManager.Contract.contract.Call(opts, result, method, params...)
}

func (_NonceManager *NonceManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NonceManager.Contract.contract.Transfer(opts)
}

func (_NonceManager *NonceManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NonceManager.Contract.contract.Transact(opts, method, params...)
}

func (_NonceManager *NonceManagerCaller) GetAllAuthorizedCallers(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _NonceManager.contract.Call(opts, &out, "getAllAuthorizedCallers")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_NonceManager *NonceManagerSession) GetAllAuthorizedCallers() ([]common.Address, error) {
	return _NonceManager.Contract.GetAllAuthorizedCallers(&_NonceManager.CallOpts)
}

func (_NonceManager *NonceManagerCallerSession) GetAllAuthorizedCallers() ([]common.Address, error) {
	return _NonceManager.Contract.GetAllAuthorizedCallers(&_NonceManager.CallOpts)
}

func (_NonceManager *NonceManagerCaller) GetOutboundNonce(opts *bind.CallOpts, destChainSelector uint64, sender []byte) (uint64, error) {
	var out []interface{}
	err := _NonceManager.contract.Call(opts, &out, "getOutboundNonce", destChainSelector, sender)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_NonceManager *NonceManagerSession) GetOutboundNonce(destChainSelector uint64, sender []byte) (uint64, error) {
	return _NonceManager.Contract.GetOutboundNonce(&_NonceManager.CallOpts, destChainSelector, sender)
}

func (_NonceManager *NonceManagerCallerSession) GetOutboundNonce(destChainSelector uint64, sender []byte) (uint64, error) {
	return _NonceManager.Contract.GetOutboundNonce(&_NonceManager.CallOpts, destChainSelector, sender)
}

func (_NonceManager *NonceManagerCaller) GetPreviousRamps(opts *bind.CallOpts, chainSelector uint64) (NonceManagerPreviousRamps, error) {
	var out []interface{}
	err := _NonceManager.contract.Call(opts, &out, "getPreviousRamps", chainSelector)

	if err != nil {
		return *new(NonceManagerPreviousRamps), err
	}

	out0 := *abi.ConvertType(out[0], new(NonceManagerPreviousRamps)).(*NonceManagerPreviousRamps)

	return out0, err

}

func (_NonceManager *NonceManagerSession) GetPreviousRamps(chainSelector uint64) (NonceManagerPreviousRamps, error) {
	return _NonceManager.Contract.GetPreviousRamps(&_NonceManager.CallOpts, chainSelector)
}

func (_NonceManager *NonceManagerCallerSession) GetPreviousRamps(chainSelector uint64) (NonceManagerPreviousRamps, error) {
	return _NonceManager.Contract.GetPreviousRamps(&_NonceManager.CallOpts, chainSelector)
}

func (_NonceManager *NonceManagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NonceManager.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_NonceManager *NonceManagerSession) Owner() (common.Address, error) {
	return _NonceManager.Contract.Owner(&_NonceManager.CallOpts)
}

func (_NonceManager *NonceManagerCallerSession) Owner() (common.Address, error) {
	return _NonceManager.Contract.Owner(&_NonceManager.CallOpts)
}

func (_NonceManager *NonceManagerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NonceManager.contract.Transact(opts, "acceptOwnership")
}

func (_NonceManager *NonceManagerSession) AcceptOwnership() (*types.Transaction, error) {
	return _NonceManager.Contract.AcceptOwnership(&_NonceManager.TransactOpts)
}

func (_NonceManager *NonceManagerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _NonceManager.Contract.AcceptOwnership(&_NonceManager.TransactOpts)
}

func (_NonceManager *NonceManagerTransactor) ApplyAuthorizedCallerUpdates(opts *bind.TransactOpts, authorizedCallerArgs AuthorizedCallersAuthorizedCallerArgs) (*types.Transaction, error) {
	return _NonceManager.contract.Transact(opts, "applyAuthorizedCallerUpdates", authorizedCallerArgs)
}

func (_NonceManager *NonceManagerSession) ApplyAuthorizedCallerUpdates(authorizedCallerArgs AuthorizedCallersAuthorizedCallerArgs) (*types.Transaction, error) {
	return _NonceManager.Contract.ApplyAuthorizedCallerUpdates(&_NonceManager.TransactOpts, authorizedCallerArgs)
}

func (_NonceManager *NonceManagerTransactorSession) ApplyAuthorizedCallerUpdates(authorizedCallerArgs AuthorizedCallersAuthorizedCallerArgs) (*types.Transaction, error) {
	return _NonceManager.Contract.ApplyAuthorizedCallerUpdates(&_NonceManager.TransactOpts, authorizedCallerArgs)
}

func (_NonceManager *NonceManagerTransactor) ApplyPreviousRampsUpdates(opts *bind.TransactOpts, previousRampsArgs []NonceManagerPreviousRampsArgs) (*types.Transaction, error) {
	return _NonceManager.contract.Transact(opts, "applyPreviousRampsUpdates", previousRampsArgs)
}

func (_NonceManager *NonceManagerSession) ApplyPreviousRampsUpdates(previousRampsArgs []NonceManagerPreviousRampsArgs) (*types.Transaction, error) {
	return _NonceManager.Contract.ApplyPreviousRampsUpdates(&_NonceManager.TransactOpts, previousRampsArgs)
}

func (_NonceManager *NonceManagerTransactorSession) ApplyPreviousRampsUpdates(previousRampsArgs []NonceManagerPreviousRampsArgs) (*types.Transaction, error) {
	return _NonceManager.Contract.ApplyPreviousRampsUpdates(&_NonceManager.TransactOpts, previousRampsArgs)
}

func (_NonceManager *NonceManagerTransactor) GetIncrementedOutboundNonce(opts *bind.TransactOpts, destChainSelector uint64, sender []byte) (*types.Transaction, error) {
	return _NonceManager.contract.Transact(opts, "getIncrementedOutboundNonce", destChainSelector, sender)
}

func (_NonceManager *NonceManagerSession) GetIncrementedOutboundNonce(destChainSelector uint64, sender []byte) (*types.Transaction, error) {
	return _NonceManager.Contract.GetIncrementedOutboundNonce(&_NonceManager.TransactOpts, destChainSelector, sender)
}

func (_NonceManager *NonceManagerTransactorSession) GetIncrementedOutboundNonce(destChainSelector uint64, sender []byte) (*types.Transaction, error) {
	return _NonceManager.Contract.GetIncrementedOutboundNonce(&_NonceManager.TransactOpts, destChainSelector, sender)
}

func (_NonceManager *NonceManagerTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _NonceManager.contract.Transact(opts, "transferOwnership", to)
}

func (_NonceManager *NonceManagerSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _NonceManager.Contract.TransferOwnership(&_NonceManager.TransactOpts, to)
}

func (_NonceManager *NonceManagerTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _NonceManager.Contract.TransferOwnership(&_NonceManager.TransactOpts, to)
}

type NonceManagerAuthorizedCallerAddedIterator struct {
	Event *NonceManagerAuthorizedCallerAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NonceManagerAuthorizedCallerAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NonceManagerAuthorizedCallerAdded)
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
		it.Event = new(NonceManagerAuthorizedCallerAdded)
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

func (it *NonceManagerAuthorizedCallerAddedIterator) Error() error {
	return it.fail
}

func (it *NonceManagerAuthorizedCallerAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NonceManagerAuthorizedCallerAdded struct {
	Caller common.Address
	Raw    types.Log
}

func (_NonceManager *NonceManagerFilterer) FilterAuthorizedCallerAdded(opts *bind.FilterOpts) (*NonceManagerAuthorizedCallerAddedIterator, error) {

	logs, sub, err := _NonceManager.contract.FilterLogs(opts, "AuthorizedCallerAdded")
	if err != nil {
		return nil, err
	}
	return &NonceManagerAuthorizedCallerAddedIterator{contract: _NonceManager.contract, event: "AuthorizedCallerAdded", logs: logs, sub: sub}, nil
}

func (_NonceManager *NonceManagerFilterer) WatchAuthorizedCallerAdded(opts *bind.WatchOpts, sink chan<- *NonceManagerAuthorizedCallerAdded) (event.Subscription, error) {

	logs, sub, err := _NonceManager.contract.WatchLogs(opts, "AuthorizedCallerAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NonceManagerAuthorizedCallerAdded)
				if err := _NonceManager.contract.UnpackLog(event, "AuthorizedCallerAdded", log); err != nil {
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

func (_NonceManager *NonceManagerFilterer) ParseAuthorizedCallerAdded(log types.Log) (*NonceManagerAuthorizedCallerAdded, error) {
	event := new(NonceManagerAuthorizedCallerAdded)
	if err := _NonceManager.contract.UnpackLog(event, "AuthorizedCallerAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type NonceManagerAuthorizedCallerRemovedIterator struct {
	Event *NonceManagerAuthorizedCallerRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NonceManagerAuthorizedCallerRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NonceManagerAuthorizedCallerRemoved)
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
		it.Event = new(NonceManagerAuthorizedCallerRemoved)
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

func (it *NonceManagerAuthorizedCallerRemovedIterator) Error() error {
	return it.fail
}

func (it *NonceManagerAuthorizedCallerRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NonceManagerAuthorizedCallerRemoved struct {
	Caller common.Address
	Raw    types.Log
}

func (_NonceManager *NonceManagerFilterer) FilterAuthorizedCallerRemoved(opts *bind.FilterOpts) (*NonceManagerAuthorizedCallerRemovedIterator, error) {

	logs, sub, err := _NonceManager.contract.FilterLogs(opts, "AuthorizedCallerRemoved")
	if err != nil {
		return nil, err
	}
	return &NonceManagerAuthorizedCallerRemovedIterator{contract: _NonceManager.contract, event: "AuthorizedCallerRemoved", logs: logs, sub: sub}, nil
}

func (_NonceManager *NonceManagerFilterer) WatchAuthorizedCallerRemoved(opts *bind.WatchOpts, sink chan<- *NonceManagerAuthorizedCallerRemoved) (event.Subscription, error) {

	logs, sub, err := _NonceManager.contract.WatchLogs(opts, "AuthorizedCallerRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NonceManagerAuthorizedCallerRemoved)
				if err := _NonceManager.contract.UnpackLog(event, "AuthorizedCallerRemoved", log); err != nil {
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

func (_NonceManager *NonceManagerFilterer) ParseAuthorizedCallerRemoved(log types.Log) (*NonceManagerAuthorizedCallerRemoved, error) {
	event := new(NonceManagerAuthorizedCallerRemoved)
	if err := _NonceManager.contract.UnpackLog(event, "AuthorizedCallerRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type NonceManagerOwnershipTransferRequestedIterator struct {
	Event *NonceManagerOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NonceManagerOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NonceManagerOwnershipTransferRequested)
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
		it.Event = new(NonceManagerOwnershipTransferRequested)
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

func (it *NonceManagerOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *NonceManagerOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NonceManagerOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_NonceManager *NonceManagerFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*NonceManagerOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _NonceManager.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &NonceManagerOwnershipTransferRequestedIterator{contract: _NonceManager.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_NonceManager *NonceManagerFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *NonceManagerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _NonceManager.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NonceManagerOwnershipTransferRequested)
				if err := _NonceManager.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_NonceManager *NonceManagerFilterer) ParseOwnershipTransferRequested(log types.Log) (*NonceManagerOwnershipTransferRequested, error) {
	event := new(NonceManagerOwnershipTransferRequested)
	if err := _NonceManager.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type NonceManagerOwnershipTransferredIterator struct {
	Event *NonceManagerOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NonceManagerOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NonceManagerOwnershipTransferred)
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
		it.Event = new(NonceManagerOwnershipTransferred)
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

func (it *NonceManagerOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *NonceManagerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NonceManagerOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_NonceManager *NonceManagerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*NonceManagerOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _NonceManager.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &NonceManagerOwnershipTransferredIterator{contract: _NonceManager.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_NonceManager *NonceManagerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *NonceManagerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _NonceManager.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NonceManagerOwnershipTransferred)
				if err := _NonceManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_NonceManager *NonceManagerFilterer) ParseOwnershipTransferred(log types.Log) (*NonceManagerOwnershipTransferred, error) {
	event := new(NonceManagerOwnershipTransferred)
	if err := _NonceManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type NonceManagerPreviousOnRampUpdatedIterator struct {
	Event *NonceManagerPreviousOnRampUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NonceManagerPreviousOnRampUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NonceManagerPreviousOnRampUpdated)
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
		it.Event = new(NonceManagerPreviousOnRampUpdated)
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

func (it *NonceManagerPreviousOnRampUpdatedIterator) Error() error {
	return it.fail
}

func (it *NonceManagerPreviousOnRampUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NonceManagerPreviousOnRampUpdated struct {
	DestChainSelector uint64
	PrevOnRamp        common.Address
	Raw               types.Log
}

func (_NonceManager *NonceManagerFilterer) FilterPreviousOnRampUpdated(opts *bind.FilterOpts, destChainSelector []uint64) (*NonceManagerPreviousOnRampUpdatedIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _NonceManager.contract.FilterLogs(opts, "PreviousOnRampUpdated", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &NonceManagerPreviousOnRampUpdatedIterator{contract: _NonceManager.contract, event: "PreviousOnRampUpdated", logs: logs, sub: sub}, nil
}

func (_NonceManager *NonceManagerFilterer) WatchPreviousOnRampUpdated(opts *bind.WatchOpts, sink chan<- *NonceManagerPreviousOnRampUpdated, destChainSelector []uint64) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _NonceManager.contract.WatchLogs(opts, "PreviousOnRampUpdated", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NonceManagerPreviousOnRampUpdated)
				if err := _NonceManager.contract.UnpackLog(event, "PreviousOnRampUpdated", log); err != nil {
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

func (_NonceManager *NonceManagerFilterer) ParsePreviousOnRampUpdated(log types.Log) (*NonceManagerPreviousOnRampUpdated, error) {
	event := new(NonceManagerPreviousOnRampUpdated)
	if err := _NonceManager.contract.UnpackLog(event, "PreviousOnRampUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_NonceManager *NonceManager) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _NonceManager.abi.Events["AuthorizedCallerAdded"].ID:
		return _NonceManager.ParseAuthorizedCallerAdded(log)
	case _NonceManager.abi.Events["AuthorizedCallerRemoved"].ID:
		return _NonceManager.ParseAuthorizedCallerRemoved(log)
	case _NonceManager.abi.Events["OwnershipTransferRequested"].ID:
		return _NonceManager.ParseOwnershipTransferRequested(log)
	case _NonceManager.abi.Events["OwnershipTransferred"].ID:
		return _NonceManager.ParseOwnershipTransferred(log)
	case _NonceManager.abi.Events["PreviousOnRampUpdated"].ID:
		return _NonceManager.ParsePreviousOnRampUpdated(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (NonceManagerAuthorizedCallerAdded) Topic() common.Hash {
	return common.HexToHash("0xeb1b9b92e50b7f88f9ff25d56765095ac6e91540eee214906f4036a908ffbdef")
}

func (NonceManagerAuthorizedCallerRemoved) Topic() common.Hash {
	return common.HexToHash("0xc3803387881faad271c47728894e3e36fac830ffc8602ca6fc07733cbda77580")
}

func (NonceManagerOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (NonceManagerOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (NonceManagerPreviousOnRampUpdated) Topic() common.Hash {
	return common.HexToHash("0x89d2355e2829b1e15855fec87fb400638aebc9f03728949d702d3b5d4ea99954")
}

func (_NonceManager *NonceManager) Address() common.Address {
	return _NonceManager.address
}

type NonceManagerInterface interface {
	GetAllAuthorizedCallers(opts *bind.CallOpts) ([]common.Address, error)

	GetOutboundNonce(opts *bind.CallOpts, destChainSelector uint64, sender []byte) (uint64, error)

	GetPreviousRamps(opts *bind.CallOpts, chainSelector uint64) (NonceManagerPreviousRamps, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	ApplyAuthorizedCallerUpdates(opts *bind.TransactOpts, authorizedCallerArgs AuthorizedCallersAuthorizedCallerArgs) (*types.Transaction, error)

	ApplyPreviousRampsUpdates(opts *bind.TransactOpts, previousRampsArgs []NonceManagerPreviousRampsArgs) (*types.Transaction, error)

	GetIncrementedOutboundNonce(opts *bind.TransactOpts, destChainSelector uint64, sender []byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterAuthorizedCallerAdded(opts *bind.FilterOpts) (*NonceManagerAuthorizedCallerAddedIterator, error)

	WatchAuthorizedCallerAdded(opts *bind.WatchOpts, sink chan<- *NonceManagerAuthorizedCallerAdded) (event.Subscription, error)

	ParseAuthorizedCallerAdded(log types.Log) (*NonceManagerAuthorizedCallerAdded, error)

	FilterAuthorizedCallerRemoved(opts *bind.FilterOpts) (*NonceManagerAuthorizedCallerRemovedIterator, error)

	WatchAuthorizedCallerRemoved(opts *bind.WatchOpts, sink chan<- *NonceManagerAuthorizedCallerRemoved) (event.Subscription, error)

	ParseAuthorizedCallerRemoved(log types.Log) (*NonceManagerAuthorizedCallerRemoved, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*NonceManagerOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *NonceManagerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*NonceManagerOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*NonceManagerOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *NonceManagerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*NonceManagerOwnershipTransferred, error)

	FilterPreviousOnRampUpdated(opts *bind.FilterOpts, destChainSelector []uint64) (*NonceManagerPreviousOnRampUpdatedIterator, error)

	WatchPreviousOnRampUpdated(opts *bind.WatchOpts, sink chan<- *NonceManagerPreviousOnRampUpdated, destChainSelector []uint64) (event.Subscription, error)

	ParsePreviousOnRampUpdated(log types.Log) (*NonceManagerPreviousOnRampUpdated, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
