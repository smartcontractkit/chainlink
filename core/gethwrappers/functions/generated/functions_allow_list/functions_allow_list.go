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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidProof\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByRouter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByRouterOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RecipientIsBlocked\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RouterMustBeSet\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"acceptor\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"proof\",\"type\":\"bytes\"}],\"name\":\"acceptTermsOfService\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"blockSender\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfigHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"config\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"messageHash\",\"type\":\"bytes32\"}],\"name\":\"getEthSignedMessageHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"acceptor\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"getMessageHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"isAllowedSender\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"isBlockedSender\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"unblockSender\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"name\":\"updateConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b5060405162000fbe38038062000fbe833981016040819052620000349162000156565b81816001600160a01b0382166200005e57604051632530e88560e11b815260040160405180910390fd5b6001600160a01b03821660805262000076816200008b565b80516020909101206000555062000284915050565b60008082806020019051810190620000a4919062000241565b6040805180820182528315158082526001600160a01b0384166020928301819052600380546001600160a81b031916610100600160a81b031984161761010090920291909117905591519182529294509092507f22aa8545955b447cb49ea37e67de742e750839c633ded8c9b5b09614843b229f910160405180910390a1505050565b6001600160a01b03811681146200013d57600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b600080604083850312156200016a57600080fd5b8251620001778162000127565b602084810151919350906001600160401b03808211156200019757600080fd5b818601915086601f830112620001ac57600080fd5b815181811115620001c157620001c162000140565b604051601f8201601f19908116603f01168101908382118183101715620001ec57620001ec62000140565b8160405282815289868487010111156200020557600080fd5b600093505b828410156200022957848401860151818501870152928501926200020a565b60008684830101528096505050505050509250929050565b600080604083850312156200025557600080fd5b825180151581146200026657600080fd5b6020840151909250620002798162000127565b809150509250929050565b608051610d10620002ae600039600081816104a9015281816105ea01526107570152610d106000f3fe608060405234801561001057600080fd5b50600436106100be5760003560e01c80638cc6acce11610076578063a5e1d61d1161005b578063a5e1d61d146101b0578063be8c97b0146101d3578063fa540801146101e657600080fd5b80638cc6acce146101955780639883c10d146101a857600080fd5b806347663acb116100a757806347663acb146100f657806380e8a1511461010957806382184c7b1461018257600080fd5b8063181f5a77146100c35780632179d447146100e1575b600080fd5b6100cb610247565b6040516100d891906109f3565b60405180910390f35b6100f46100ef366004610a84565b610267565b005b6100f4610104366004610b16565b6104a7565b610174610117366004610b3a565b6040517fffffffffffffffffffffffffffffffffffffffff000000000000000000000000606084811b8216602084015283901b16603482015260009060480160405160208183030381529060405280519060200120905092915050565b6040519081526020016100d8565b6100f4610190366004610b16565b6105e8565b6100f46101a3366004610ba2565b61073f565b600054610174565b6101c36101be366004610b16565b6107c4565b60405190151581526020016100d8565b6101c36101e1366004610b16565b610805565b6101746101f4366004610c71565b6040517f19457468657265756d205369676e6564204d6573736167653a0a3332000000006020820152603c8101829052600090605c01604051602081830303815290604052805190602001209050919050565b6060604051806060016040528060288152602001610cdc60289139905090565b73ffffffffffffffffffffffffffffffffffffffff831660009081526002602052604090205460ff16156102c7576040517f62b7a34d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600354610100900473ffffffffffffffffffffffffffffffffffffffff1661038961034d6101f487876040517fffffffffffffffffffffffffffffffffffffffff000000000000000000000000606084811b8216602084015283901b16603482015260009060480160405160208183030381529060405280519060200120905092915050565b84848080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061084692505050565b73ffffffffffffffffffffffffffffffffffffffff16146103d6576040517f09bde33900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff841614158061041b57503373ffffffffffffffffffffffffffffffffffffffff85161480159061041b5750333b155b15610452576040517f09bde33900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b505073ffffffffffffffffffffffffffffffffffffffff16600090815260016020819052604090912080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016909117905550565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16638da5cb5b6040518163ffffffff1660e01b81526004016020604051808303816000875af1158015610514573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105389190610c8a565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461059c576040517fa0f0a44600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff16600090815260026020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16638da5cb5b6040518163ffffffff1660e01b81526004016020604051808303816000875af1158015610655573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106799190610c8a565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146106dd576040517fa0f0a44600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff16600090815260016020818152604080842080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0090811690915560029092529092208054909216179055565b3373ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016146107ae576040517fc41a5b0900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6107b78161091d565b8051602090910120600055565b60035460009060ff166107d957506000919050565b5073ffffffffffffffffffffffffffffffffffffffff1660009081526002602052604090205460ff1690565b60035460009060ff1661081a57506001919050565b5073ffffffffffffffffffffffffffffffffffffffff1660009081526001602052604090205460ff1690565b6000806000808451604114610887576040517f09bde33900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b505050602082810151604080850151606080870151835160008082529681018086528a9052951a928501839052840183905260808401819052919260019060a0016020604051602081039080840390855afa1580156108ea573d6000803e3d6000fd5b50506040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00151979650505050505050565b600080828060200190518101906109349190610ca7565b60408051808201825283151580825273ffffffffffffffffffffffffffffffffffffffff84166020928301819052600380547fffffffffffffffffffffff000000000000000000000000000000000000000000167fffffffffffffffffffffff0000000000000000000000000000000000000000ff84161761010090920291909117905591519182529294509092507f22aa8545955b447cb49ea37e67de742e750839c633ded8c9b5b09614843b229f910160405180910390a1505050565b600060208083528351808285015260005b81811015610a2057858101830151858201604001528201610a04565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b73ffffffffffffffffffffffffffffffffffffffff81168114610a8157600080fd5b50565b60008060008060608587031215610a9a57600080fd5b8435610aa581610a5f565b93506020850135610ab581610a5f565b9250604085013567ffffffffffffffff80821115610ad257600080fd5b818701915087601f830112610ae657600080fd5b813581811115610af557600080fd5b886020828501011115610b0757600080fd5b95989497505060200194505050565b600060208284031215610b2857600080fd5b8135610b3381610a5f565b9392505050565b60008060408385031215610b4d57600080fd5b8235610b5881610a5f565b91506020830135610b6881610a5f565b809150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600060208284031215610bb457600080fd5b813567ffffffffffffffff80821115610bcc57600080fd5b818401915084601f830112610be057600080fd5b813581811115610bf257610bf2610b73565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f01168101908382118183101715610c3857610c38610b73565b81604052828152876020848701011115610c5157600080fd5b826020860160208301376000928101602001929092525095945050505050565b600060208284031215610c8357600080fd5b5035919050565b600060208284031215610c9c57600080fd5b8151610b3381610a5f565b60008060408385031215610cba57600080fd5b82518015158114610cca57600080fd5b6020840151909250610b6881610a5f56fe46756e6374696f6e73205465726d73206f66205365727669636520416c6c6f77204c697374207631a164736f6c6343000813000a",
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

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCaller) GetAllAllowedSenders(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _TermsOfServiceAllowList.contract.Call(opts, &out, "getAllAllowedSenders")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) GetAllAllowedSenders() ([]common.Address, error) {
	return _TermsOfServiceAllowList.Contract.GetAllAllowedSenders(&_TermsOfServiceAllowList.CallOpts)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCallerSession) GetAllAllowedSenders() ([]common.Address, error) {
	return _TermsOfServiceAllowList.Contract.GetAllAllowedSenders(&_TermsOfServiceAllowList.CallOpts)
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

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCaller) HasAccess(opts *bind.CallOpts, user common.Address, arg1 []byte) (bool, error) {
	var out []interface{}
	err := _TermsOfServiceAllowList.contract.Call(opts, &out, "hasAccess", user, arg1)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) HasAccess(user common.Address, arg1 []byte) (bool, error) {
	return _TermsOfServiceAllowList.Contract.HasAccess(&_TermsOfServiceAllowList.CallOpts, user, arg1)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCallerSession) HasAccess(user common.Address, arg1 []byte) (bool, error) {
	return _TermsOfServiceAllowList.Contract.HasAccess(&_TermsOfServiceAllowList.CallOpts, user, arg1)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCaller) HasAccess0(opts *bind.CallOpts, sender common.Address) (bool, error) {
	var out []interface{}
	err := _TermsOfServiceAllowList.contract.Call(opts, &out, "hasAccess0", sender)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) HasAccess0(sender common.Address) (bool, error) {
	return _TermsOfServiceAllowList.Contract.HasAccess0(&_TermsOfServiceAllowList.CallOpts, sender)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCallerSession) HasAccess0(sender common.Address) (bool, error) {
	return _TermsOfServiceAllowList.Contract.HasAccess0(&_TermsOfServiceAllowList.CallOpts, sender)
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

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactor) UnblockSender(opts *bind.TransactOpts, sender common.Address) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.contract.Transact(opts, "unblockSender", sender)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) UnblockSender(sender common.Address) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.UnblockSender(&_TermsOfServiceAllowList.TransactOpts, sender)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactorSession) UnblockSender(sender common.Address) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.UnblockSender(&_TermsOfServiceAllowList.TransactOpts, sender)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactor) UpdateConfig(opts *bind.TransactOpts, config []byte) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.contract.Transact(opts, "updateConfig", config)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) UpdateConfig(config []byte) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.UpdateConfig(&_TermsOfServiceAllowList.TransactOpts, config)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactorSession) UpdateConfig(config []byte) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.UpdateConfig(&_TermsOfServiceAllowList.TransactOpts, config)
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
	GetAllAllowedSenders(opts *bind.CallOpts) ([]common.Address, error)

	GetConfigHash(opts *bind.CallOpts) ([32]byte, error)

	GetEthSignedMessageHash(opts *bind.CallOpts, messageHash [32]byte) ([32]byte, error)

	GetMessageHash(opts *bind.CallOpts, acceptor common.Address, recipient common.Address) ([32]byte, error)

	HasAccess(opts *bind.CallOpts, user common.Address, arg1 []byte) (bool, error)

	HasAccess0(opts *bind.CallOpts, sender common.Address) (bool, error)

	IsBlockedSender(opts *bind.CallOpts, sender common.Address) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptTermsOfService(opts *bind.TransactOpts, acceptor common.Address, recipient common.Address, proof []byte) (*types.Transaction, error)

	BlockSender(opts *bind.TransactOpts, sender common.Address) (*types.Transaction, error)

	UnblockSender(opts *bind.TransactOpts, sender common.Address) (*types.Transaction, error)

	UpdateConfig(opts *bind.TransactOpts, config []byte) (*types.Transaction, error)

	FilterConfigSet(opts *bind.FilterOpts) (*TermsOfServiceAllowListConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *TermsOfServiceAllowListConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*TermsOfServiceAllowListConfigSet, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
