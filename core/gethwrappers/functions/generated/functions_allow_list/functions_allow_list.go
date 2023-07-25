// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package functions_allow_list

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

var TermsOfServiceAllowListMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidProof\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByRouter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByRouterOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RecipientIsBlocked\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RouterMustBeSet\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"acceptor\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"proof\",\"type\":\"bytes\"}],\"name\":\"acceptTermsOfService\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"blockSender\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfigHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"config\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"messageHash\",\"type\":\"bytes32\"}],\"name\":\"getEthSignedMessageHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"acceptor\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"getMessageHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"isAllowedSender\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"isBlockedSender\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"unblockSender\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b506040516200114238038062001142833981016040819052620000349162000164565b81816001600160a01b0382166200005e57604051632530e88560e11b815260040160405180910390fd5b600180546001600160a01b0319166001600160a01b038416179055620000848162000099565b80516020909101206000555062000292915050565b60008082806020019051810190620000b291906200024f565b6040805180820182528315158082526001600160a01b0384166020928301819052600680546001600160a81b031916610100600160a81b031984161761010090920291909117905591519182529294509092507f22aa8545955b447cb49ea37e67de742e750839c633ded8c9b5b09614843b229f910160405180910390a1505050565b6001600160a01b03811681146200014b57600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b600080604083850312156200017857600080fd5b8251620001858162000135565b602084810151919350906001600160401b0380821115620001a557600080fd5b818601915086601f830112620001ba57600080fd5b815181811115620001cf57620001cf6200014e565b604051601f8201601f19908116603f01168101908382118183101715620001fa57620001fa6200014e565b8160405282815289868487010111156200021357600080fd5b600093505b8284101562000237578484018601518185018701529285019262000218565b60008684830101528096505050505050509250929050565b600080604083850312156200026357600080fd5b825180151581146200027457600080fd5b6020840151909250620002878162000135565b809150509250929050565b610ea080620002a26000396000f3fe608060405234801561001057600080fd5b50600436106100be5760003560e01c806382184c7b11610076578063a5e1d61d1161005b578063a5e1d61d14610158578063be8c97b01461017b578063fa5408011461018e57600080fd5b806382184c7b1461013d5780639883c10d1461015057600080fd5b8063354497aa116100a7578063354497aa146100f657806347663acb1461010957806380e8a1511461011c57600080fd5b8063181f5a77146100c35780632179d447146100e1575b600080fd5b6100cb6101ef565b6040516100d89190610af2565b60405180910390f35b6100f46100ef366004610b83565b61020f565b005b6100f4610104366004610c44565b61044c565b6100f4610117366004610d13565b6104b3565b61012f61012a366004610d30565b6105b9565b6040519081526020016100d8565b6100f461014b366004610d13565b610617565b60005461012f565b61016b610166366004610d13565b610725565b60405190151581526020016100d8565b61016b610189366004610d13565b610749565b61012f61019c366004610d69565b6040517f19457468657265756d205369676e6564204d6573736167653a0a3332000000006020820152603c8101829052600090605c01604051602081830303815290604052805190602001209050919050565b6060604051806060016040528060288152602001610e6c60289139905090565b61021a600484610769565b15610251576040517f62b7a34d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600061025d85856105b9565b905060006102b8826040517f19457468657265756d205369676e6564204d6573736167653a0a3332000000006020820152603c8101829052600090605c01604051602081830303815290604052805190602001209050919050565b905060006102fc8286868080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061079b92505050565b60065490915073ffffffffffffffffffffffffffffffffffffffff8083166101009092041614610358576040517f09bde33900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b63ffffffff333b1615801590816103a557503373ffffffffffffffffffffffffffffffffffffffff89161415806103a557503373ffffffffffffffffffffffffffffffffffffffff881614155b156103dc576040517f09bde33900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8080156103ff57503373ffffffffffffffffffffffffffffffffffffffff881614155b15610436576040517f09bde33900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610441600288610838565b505050505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff16331461049d576040517fc41a5b0900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6104a68161085a565b8051602090910120600055565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16638da5cb5b6040518163ffffffff1660e01b81526004016020604051808303816000875af1158015610522573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105469190610d82565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146105aa576040517fa0f0a44600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6105b5600482610930565b5050565b6040517fffffffffffffffffffffffffffffffffffffffff000000000000000000000000606084811b8216602084015283901b1660348201526000906048016040516020818303038152906040528051906020012090505b92915050565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16638da5cb5b6040518163ffffffff1660e01b81526004016020604051808303816000875af1158015610686573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106aa9190610d82565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461070e576040517fa0f0a44600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610719600282610930565b506105b5600482610838565b60065460009060ff161515810361073e57506000919050565b610611600483610769565b60065460009060ff161515810361076257506001919050565b6106116002835b73ffffffffffffffffffffffffffffffffffffffff8116600090815260018301602052604081205415155b9392505050565b6000806000806107aa85610952565b6040805160008152602081018083528b905260ff8316918101919091526060810184905260808101839052929550909350915060019060a0016020604051602081039080840390855afa158015610805573d6000803e3d6000fd5b50506040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00151979650505050505050565b60006107948373ffffffffffffffffffffffffffffffffffffffff84166109b0565b600080828060200190518101906108719190610d9f565b60408051808201825283151580825273ffffffffffffffffffffffffffffffffffffffff84166020928301819052600680547fffffffffffffffffffffff000000000000000000000000000000000000000000167fffffffffffffffffffffff0000000000000000000000000000000000000000ff84161761010090920291909117905591519182529294509092507f22aa8545955b447cb49ea37e67de742e750839c633ded8c9b5b09614843b229f910160405180910390a1505050565b60006107948373ffffffffffffffffffffffffffffffffffffffff84166109ff565b60008060008351604114610992576040517f09bde33900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b50505060208101516040820151606090920151909260009190911a90565b60008181526001830160205260408120546109f757508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610611565b506000610611565b60008181526001830160205260408120548015610ae8576000610a23600183610dd3565b8554909150600090610a3790600190610dd3565b9050818114610a9c576000866000018281548110610a5757610a57610e0d565b9060005260206000200154905080876000018481548110610a7a57610a7a610e0d565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080610aad57610aad610e3c565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050610611565b6000915050610611565b600060208083528351808285015260005b81811015610b1f57858101830151858201604001528201610b03565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b73ffffffffffffffffffffffffffffffffffffffff81168114610b8057600080fd5b50565b60008060008060608587031215610b9957600080fd5b8435610ba481610b5e565b93506020850135610bb481610b5e565b9250604085013567ffffffffffffffff80821115610bd157600080fd5b818701915087601f830112610be557600080fd5b813581811115610bf457600080fd5b886020828501011115610c0657600080fd5b95989497505060200194505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600060208284031215610c5657600080fd5b813567ffffffffffffffff80821115610c6e57600080fd5b818401915084601f830112610c8257600080fd5b813581811115610c9457610c94610c15565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f01168101908382118183101715610cda57610cda610c15565b81604052828152876020848701011115610cf357600080fd5b826020860160208301376000928101602001929092525095945050505050565b600060208284031215610d2557600080fd5b813561079481610b5e565b60008060408385031215610d4357600080fd5b8235610d4e81610b5e565b91506020830135610d5e81610b5e565b809150509250929050565b600060208284031215610d7b57600080fd5b5035919050565b600060208284031215610d9457600080fd5b815161079481610b5e565b60008060408385031215610db257600080fd5b82518015158114610dc257600080fd5b6020840151909250610d5e81610b5e565b81810381811115610611577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfe46756e6374696f6e73205465726d73206f66205365727669636520416c6c6f77204c697374207631a164736f6c6343000813000a",
}

var TermsOfServiceAllowListABI = TermsOfServiceAllowListMetaData.ABI

var TermsOfServiceAllowListBin = TermsOfServiceAllowListMetaData.Bin

func DeployTermsOfServiceAllowList(auth *bind.TransactOpts, backend bind.ContractBackend, router common.Address, config []byte) (common.Address, *types.Transaction, *TermsOfServiceAllowList, error) {
	parsed, err := TermsOfServiceAllowListMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TermsOfServiceAllowListBin), backend, router, config)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TermsOfServiceAllowList{TermsOfServiceAllowListCaller: TermsOfServiceAllowListCaller{contract: contract}, TermsOfServiceAllowListTransactor: TermsOfServiceAllowListTransactor{contract: contract}, TermsOfServiceAllowListFilterer: TermsOfServiceAllowListFilterer{contract: contract}}, nil
}

type TermsOfServiceAllowList struct {
	address common.Address
	abi     abi.ABI
	TermsOfServiceAllowListCaller
	TermsOfServiceAllowListTransactor
	TermsOfServiceAllowListFilterer
}

type TermsOfServiceAllowListCaller struct {
	contract *bind.BoundContract
}

type TermsOfServiceAllowListTransactor struct {
	contract *bind.BoundContract
}

type TermsOfServiceAllowListFilterer struct {
	contract *bind.BoundContract
}

type TermsOfServiceAllowListSession struct {
	Contract     *TermsOfServiceAllowList
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type TermsOfServiceAllowListCallerSession struct {
	Contract *TermsOfServiceAllowListCaller
	CallOpts bind.CallOpts
}

type TermsOfServiceAllowListTransactorSession struct {
	Contract     *TermsOfServiceAllowListTransactor
	TransactOpts bind.TransactOpts
}

type TermsOfServiceAllowListRaw struct {
	Contract *TermsOfServiceAllowList
}

type TermsOfServiceAllowListCallerRaw struct {
	Contract *TermsOfServiceAllowListCaller
}

type TermsOfServiceAllowListTransactorRaw struct {
	Contract *TermsOfServiceAllowListTransactor
}

func NewTermsOfServiceAllowList(address common.Address, backend bind.ContractBackend) (*TermsOfServiceAllowList, error) {
	abi, err := abi.JSON(strings.NewReader(TermsOfServiceAllowListABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindTermsOfServiceAllowList(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TermsOfServiceAllowList{address: address, abi: abi, TermsOfServiceAllowListCaller: TermsOfServiceAllowListCaller{contract: contract}, TermsOfServiceAllowListTransactor: TermsOfServiceAllowListTransactor{contract: contract}, TermsOfServiceAllowListFilterer: TermsOfServiceAllowListFilterer{contract: contract}}, nil
}

func NewTermsOfServiceAllowListCaller(address common.Address, caller bind.ContractCaller) (*TermsOfServiceAllowListCaller, error) {
	contract, err := bindTermsOfServiceAllowList(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TermsOfServiceAllowListCaller{contract: contract}, nil
}

func NewTermsOfServiceAllowListTransactor(address common.Address, transactor bind.ContractTransactor) (*TermsOfServiceAllowListTransactor, error) {
	contract, err := bindTermsOfServiceAllowList(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TermsOfServiceAllowListTransactor{contract: contract}, nil
}

func NewTermsOfServiceAllowListFilterer(address common.Address, filterer bind.ContractFilterer) (*TermsOfServiceAllowListFilterer, error) {
	contract, err := bindTermsOfServiceAllowList(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TermsOfServiceAllowListFilterer{contract: contract}, nil
}

func bindTermsOfServiceAllowList(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TermsOfServiceAllowListMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TermsOfServiceAllowList.Contract.TermsOfServiceAllowListCaller.contract.Call(opts, result, method, params...)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.TermsOfServiceAllowListTransactor.contract.Transfer(opts)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.TermsOfServiceAllowListTransactor.contract.Transact(opts, method, params...)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TermsOfServiceAllowList.Contract.contract.Call(opts, result, method, params...)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.contract.Transfer(opts)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.contract.Transact(opts, method, params...)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCaller) GetConfigHash(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TermsOfServiceAllowList.contract.Call(opts, &out, "getConfigHash")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) GetConfigHash() ([32]byte, error) {
	return _TermsOfServiceAllowList.Contract.GetConfigHash(&_TermsOfServiceAllowList.CallOpts)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCallerSession) GetConfigHash() ([32]byte, error) {
	return _TermsOfServiceAllowList.Contract.GetConfigHash(&_TermsOfServiceAllowList.CallOpts)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCaller) GetEthSignedMessageHash(opts *bind.CallOpts, messageHash [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _TermsOfServiceAllowList.contract.Call(opts, &out, "getEthSignedMessageHash", messageHash)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) GetEthSignedMessageHash(messageHash [32]byte) ([32]byte, error) {
	return _TermsOfServiceAllowList.Contract.GetEthSignedMessageHash(&_TermsOfServiceAllowList.CallOpts, messageHash)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCallerSession) GetEthSignedMessageHash(messageHash [32]byte) ([32]byte, error) {
	return _TermsOfServiceAllowList.Contract.GetEthSignedMessageHash(&_TermsOfServiceAllowList.CallOpts, messageHash)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCaller) GetMessageHash(opts *bind.CallOpts, acceptor common.Address, recipient common.Address) ([32]byte, error) {
	var out []interface{}
	err := _TermsOfServiceAllowList.contract.Call(opts, &out, "getMessageHash", acceptor, recipient)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) GetMessageHash(acceptor common.Address, recipient common.Address) ([32]byte, error) {
	return _TermsOfServiceAllowList.Contract.GetMessageHash(&_TermsOfServiceAllowList.CallOpts, acceptor, recipient)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCallerSession) GetMessageHash(acceptor common.Address, recipient common.Address) ([32]byte, error) {
	return _TermsOfServiceAllowList.Contract.GetMessageHash(&_TermsOfServiceAllowList.CallOpts, acceptor, recipient)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCaller) IsAllowedSender(opts *bind.CallOpts, sender common.Address) (bool, error) {
	var out []interface{}
	err := _TermsOfServiceAllowList.contract.Call(opts, &out, "isAllowedSender", sender)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) IsAllowedSender(sender common.Address) (bool, error) {
	return _TermsOfServiceAllowList.Contract.IsAllowedSender(&_TermsOfServiceAllowList.CallOpts, sender)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCallerSession) IsAllowedSender(sender common.Address) (bool, error) {
	return _TermsOfServiceAllowList.Contract.IsAllowedSender(&_TermsOfServiceAllowList.CallOpts, sender)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCaller) IsBlockedSender(opts *bind.CallOpts, sender common.Address) (bool, error) {
	var out []interface{}
	err := _TermsOfServiceAllowList.contract.Call(opts, &out, "isBlockedSender", sender)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) IsBlockedSender(sender common.Address) (bool, error) {
	return _TermsOfServiceAllowList.Contract.IsBlockedSender(&_TermsOfServiceAllowList.CallOpts, sender)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCallerSession) IsBlockedSender(sender common.Address) (bool, error) {
	return _TermsOfServiceAllowList.Contract.IsBlockedSender(&_TermsOfServiceAllowList.CallOpts, sender)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TermsOfServiceAllowList.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) TypeAndVersion() (string, error) {
	return _TermsOfServiceAllowList.Contract.TypeAndVersion(&_TermsOfServiceAllowList.CallOpts)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCallerSession) TypeAndVersion() (string, error) {
	return _TermsOfServiceAllowList.Contract.TypeAndVersion(&_TermsOfServiceAllowList.CallOpts)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactor) AcceptTermsOfService(opts *bind.TransactOpts, acceptor common.Address, recipient common.Address, proof []byte) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.contract.Transact(opts, "acceptTermsOfService", acceptor, recipient, proof)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) AcceptTermsOfService(acceptor common.Address, recipient common.Address, proof []byte) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.AcceptTermsOfService(&_TermsOfServiceAllowList.TransactOpts, acceptor, recipient, proof)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactorSession) AcceptTermsOfService(acceptor common.Address, recipient common.Address, proof []byte) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.AcceptTermsOfService(&_TermsOfServiceAllowList.TransactOpts, acceptor, recipient, proof)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactor) BlockSender(opts *bind.TransactOpts, sender common.Address) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.contract.Transact(opts, "blockSender", sender)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) BlockSender(sender common.Address) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.BlockSender(&_TermsOfServiceAllowList.TransactOpts, sender)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactorSession) BlockSender(sender common.Address) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.BlockSender(&_TermsOfServiceAllowList.TransactOpts, sender)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactor) SetConfig(opts *bind.TransactOpts, config []byte) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.contract.Transact(opts, "setConfig", config)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) SetConfig(config []byte) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.SetConfig(&_TermsOfServiceAllowList.TransactOpts, config)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactorSession) SetConfig(config []byte) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.SetConfig(&_TermsOfServiceAllowList.TransactOpts, config)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactor) UnblockSender(opts *bind.TransactOpts, sender common.Address) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.contract.Transact(opts, "unblockSender", sender)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) UnblockSender(sender common.Address) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.UnblockSender(&_TermsOfServiceAllowList.TransactOpts, sender)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactorSession) UnblockSender(sender common.Address) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.UnblockSender(&_TermsOfServiceAllowList.TransactOpts, sender)
}

type TermsOfServiceAllowListConfigSetIterator struct {
	Event *TermsOfServiceAllowListConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TermsOfServiceAllowListConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TermsOfServiceAllowListConfigSet)
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
		it.Event = new(TermsOfServiceAllowListConfigSet)
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

func (it *TermsOfServiceAllowListConfigSetIterator) Error() error {
	return it.fail
}

func (it *TermsOfServiceAllowListConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TermsOfServiceAllowListConfigSet struct {
	Enabled bool
	Raw     types.Log
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListFilterer) FilterConfigSet(opts *bind.FilterOpts) (*TermsOfServiceAllowListConfigSetIterator, error) {

	logs, sub, err := _TermsOfServiceAllowList.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &TermsOfServiceAllowListConfigSetIterator{contract: _TermsOfServiceAllowList.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *TermsOfServiceAllowListConfigSet) (event.Subscription, error) {

	logs, sub, err := _TermsOfServiceAllowList.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TermsOfServiceAllowListConfigSet)
				if err := _TermsOfServiceAllowList.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_TermsOfServiceAllowList *TermsOfServiceAllowListFilterer) ParseConfigSet(log types.Log) (*TermsOfServiceAllowListConfigSet, error) {
	event := new(TermsOfServiceAllowListConfigSet)
	if err := _TermsOfServiceAllowList.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowList) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _TermsOfServiceAllowList.abi.Events["ConfigSet"].ID:
		return _TermsOfServiceAllowList.ParseConfigSet(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (TermsOfServiceAllowListConfigSet) Topic() common.Hash {
	return common.HexToHash("0x22aa8545955b447cb49ea37e67de742e750839c633ded8c9b5b09614843b229f")
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowList) Address() common.Address {
	return _TermsOfServiceAllowList.address
}

type TermsOfServiceAllowListInterface interface {
	GetConfigHash(opts *bind.CallOpts) ([32]byte, error)

	GetEthSignedMessageHash(opts *bind.CallOpts, messageHash [32]byte) ([32]byte, error)

	GetMessageHash(opts *bind.CallOpts, acceptor common.Address, recipient common.Address) ([32]byte, error)

	IsAllowedSender(opts *bind.CallOpts, sender common.Address) (bool, error)

	IsBlockedSender(opts *bind.CallOpts, sender common.Address) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptTermsOfService(opts *bind.TransactOpts, acceptor common.Address, recipient common.Address, proof []byte) (*types.Transaction, error)

	BlockSender(opts *bind.TransactOpts, sender common.Address) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, config []byte) (*types.Transaction, error)

	UnblockSender(opts *bind.TransactOpts, sender common.Address) (*types.Transaction, error)

	FilterConfigSet(opts *bind.FilterOpts) (*TermsOfServiceAllowListConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *TermsOfServiceAllowListConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*TermsOfServiceAllowListConfigSet, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
