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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidSignature\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidUsage\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByRouter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByRouterOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RecipientIsBlocked\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RouterMustBeSet\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"AddedAccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"BlockedAccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"signerPublicKey\",\"type\":\"address\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"UnblockedAccess\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"acceptor\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"}],\"name\":\"acceptTermsOfService\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"blockSender\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllAllowedSenders\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"acceptor\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"getMessage\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"hasAccess\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"isBlockedSender\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"unblockSender\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"name\":\"updateConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b50604051620012443803806200124483398101604081905262000034916200014e565b81816001600160a01b0382166200005e57604051632530e88560e11b815260040160405180910390fd5b6001600160a01b038216608052620000768162000080565b505050506200027c565b6000808280602001905181019062000099919062000239565b6040805180820182528315158082526001600160a01b0384166020928301819052600380546001600160a81b031916610100600160a81b031984161761010083021790558351918252918101919091529294509092507f232bb132911bf7acb3a5ab4dd8e603080819db6a26b68da198599a7c83e2494c910160405180910390a1505050565b6001600160a01b03811681146200013557600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b600080604083850312156200016257600080fd5b82516200016f816200011f565b602084810151919350906001600160401b03808211156200018f57600080fd5b818601915086601f830112620001a457600080fd5b815181811115620001b957620001b962000138565b604051601f8201601f19908116603f01168101908382118183101715620001e457620001e462000138565b816040528281528986848701011115620001fd57600080fd5b600093505b8284101562000221578484018601518185018701529285019262000202565b60008684830101528096505050505050509250929050565b600080604083850312156200024d57600080fd5b825180151581146200025e57600080fd5b602084015190925062000271816200011f565b809150509250929050565b608051610f9e620002a660003960008181610414015281816105c7015261075d0152610f9e6000f3fe608060405234801561001057600080fd5b50600436106100a35760003560e01c8063817ef62e116100765780638cc6acce1161005b5780638cc6acce14610139578063a39b06e31461014c578063a5e1d61d1461016d57600080fd5b8063817ef62e1461011157806382184c7b1461012657600080fd5b8063181f5a77146100a85780633908c4d4146100c657806347663acb146100db5780636b14daf8146100ee575b600080fd5b6100b0610180565b6040516100bd9190610b56565b60405180910390f35b6100d96100d4366004610be4565b61019c565b005b6100d96100e9366004610c49565b610412565b6101016100fc366004610c66565b61058a565b60405190151581526020016100bd565b6101196105b4565b6040516100bd9190610ceb565b6100d9610134366004610c49565b6105c5565b6100d9610147366004610d74565b610745565b61015f61015a366004610e43565b6107c0565b6040519081526020016100bd565b61010161017b366004610c49565b61081e565b6040518060600160405280602c8152602001610f66602c913981565b73ffffffffffffffffffffffffffffffffffffffff841660009081526002602052604090205460ff16156101fc576040517f62b7a34d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600061020886866107c0565b6040517f19457468657265756d205369676e6564204d6573736167653a0a3332000000006020820152603c810191909152605c01604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001815282825280516020918201206003546000855291840180845281905260ff8616928401929092526060830187905260808301869052909250610100900473ffffffffffffffffffffffffffffffffffffffff169060019060a0016020604051602081039080840390855afa1580156102e2573d6000803e3d6000fd5b5050506020604051035173ffffffffffffffffffffffffffffffffffffffff1614610339576040517f8baa579f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff861614158061037e57503373ffffffffffffffffffffffffffffffffffffffff87161480159061037e5750333b155b156103b5576040517f381cfcbd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6103c060008661085f565b5060405173ffffffffffffffffffffffffffffffffffffffff861681527f87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db49060200160405180910390a1505050505050565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16638da5cb5b6040518163ffffffff1660e01b81526004016020604051808303816000875af115801561047f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104a39190610e7c565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610507576040517fa0f0a44600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff811660008181526002602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905590519182527f28bbd0761309a99e8fb5e5d02ada0b7b2db2e5357531ff5dbfc205c3f5b6592b91015b60405180910390a150565b60035460009060ff1661059f575060016105ad565b6105aa600085610881565b90505b9392505050565b60606105c060006108b0565b905090565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16638da5cb5b6040518163ffffffff1660e01b81526004016020604051808303816000875af1158015610632573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106569190610e7c565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146106ba576040517fa0f0a44600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6106c56000826108bd565b5073ffffffffffffffffffffffffffffffffffffffff811660008181526002602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016600117905590519182527f337cd0f3f594112b6d830afb510072d3b08556b446514f73b8109162fd1151e1910161057f565b3373ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016146107b4576040517fc41a5b0900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6107bd816108df565b50565b6040517fffffffffffffffffffffffffffffffffffffffff000000000000000000000000606084811b8216602084015283901b1660348201526000906048016040516020818303038152906040528051906020012090505b92915050565b60035460009060ff1661083357506000919050565b5073ffffffffffffffffffffffffffffffffffffffff1660009081526002602052604090205460ff1690565b60006105ad8373ffffffffffffffffffffffffffffffffffffffff84166109b8565b73ffffffffffffffffffffffffffffffffffffffff8116600090815260018301602052604081205415156105ad565b606060006105ad83610a07565b60006105ad8373ffffffffffffffffffffffffffffffffffffffff8416610a63565b600080828060200190518101906108f69190610e99565b60408051808201825283151580825273ffffffffffffffffffffffffffffffffffffffff84166020928301819052600380547fffffffffffffffffffffff000000000000000000000000000000000000000000167fffffffffffffffffffffff0000000000000000000000000000000000000000ff84161761010083021790558351918252918101919091529294509092507f232bb132911bf7acb3a5ab4dd8e603080819db6a26b68da198599a7c83e2494c910160405180910390a1505050565b60008181526001830160205260408120546109ff57508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610818565b506000610818565b606081600001805480602002602001604051908101604052809291908181526020018280548015610a5757602002820191906000526020600020905b815481526020019060010190808311610a43575b50505050509050919050565b60008181526001830160205260408120548015610b4c576000610a87600183610ecd565b8554909150600090610a9b90600190610ecd565b9050818114610b00576000866000018281548110610abb57610abb610f07565b9060005260206000200154905080876000018481548110610ade57610ade610f07565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080610b1157610b11610f36565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050610818565b6000915050610818565b600060208083528351808285015260005b81811015610b8357858101830151858201604001528201610b67565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b73ffffffffffffffffffffffffffffffffffffffff811681146107bd57600080fd5b600080600080600060a08688031215610bfc57600080fd5b8535610c0781610bc2565b94506020860135610c1781610bc2565b93506040860135925060608601359150608086013560ff81168114610c3b57600080fd5b809150509295509295909350565b600060208284031215610c5b57600080fd5b81356105ad81610bc2565b600080600060408486031215610c7b57600080fd5b8335610c8681610bc2565b9250602084013567ffffffffffffffff80821115610ca357600080fd5b818601915086601f830112610cb757600080fd5b813581811115610cc657600080fd5b876020828501011115610cd857600080fd5b6020830194508093505050509250925092565b6020808252825182820181905260009190848201906040850190845b81811015610d3957835173ffffffffffffffffffffffffffffffffffffffff1683529284019291840191600101610d07565b50909695505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600060208284031215610d8657600080fd5b813567ffffffffffffffff80821115610d9e57600080fd5b818401915084601f830112610db257600080fd5b813581811115610dc457610dc4610d45565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f01168101908382118183101715610e0a57610e0a610d45565b81604052828152876020848701011115610e2357600080fd5b826020860160208301376000928101602001929092525095945050505050565b60008060408385031215610e5657600080fd5b8235610e6181610bc2565b91506020830135610e7181610bc2565b809150509250929050565b600060208284031215610e8e57600080fd5b81516105ad81610bc2565b60008060408385031215610eac57600080fd5b82518015158114610ebc57600080fd5b6020840151909250610e7181610bc2565b81810381811115610818577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfe46756e6374696f6e73205465726d73206f66205365727669636520416c6c6f77204c6973742076312e302e30a164736f6c6343000813000a",
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

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCaller) GetMessage(opts *bind.CallOpts, acceptor common.Address, recipient common.Address) ([32]byte, error) {
	var out []interface{}
	err := _TermsOfServiceAllowList.contract.Call(opts, &out, "getMessage", acceptor, recipient)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) GetMessage(acceptor common.Address, recipient common.Address) ([32]byte, error) {
	return _TermsOfServiceAllowList.Contract.GetMessage(&_TermsOfServiceAllowList.CallOpts, acceptor, recipient)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCallerSession) GetMessage(acceptor common.Address, recipient common.Address) ([32]byte, error) {
	return _TermsOfServiceAllowList.Contract.GetMessage(&_TermsOfServiceAllowList.CallOpts, acceptor, recipient)
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

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactor) AcceptTermsOfService(opts *bind.TransactOpts, acceptor common.Address, recipient common.Address, r [32]byte, s [32]byte, v uint8) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.contract.Transact(opts, "acceptTermsOfService", acceptor, recipient, r, s, v)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) AcceptTermsOfService(acceptor common.Address, recipient common.Address, r [32]byte, s [32]byte, v uint8) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.AcceptTermsOfService(&_TermsOfServiceAllowList.TransactOpts, acceptor, recipient, r, s, v)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactorSession) AcceptTermsOfService(acceptor common.Address, recipient common.Address, r [32]byte, s [32]byte, v uint8) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.AcceptTermsOfService(&_TermsOfServiceAllowList.TransactOpts, acceptor, recipient, r, s, v)
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

type TermsOfServiceAllowListAddedAccessIterator struct {
	Event *TermsOfServiceAllowListAddedAccess

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TermsOfServiceAllowListAddedAccessIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TermsOfServiceAllowListAddedAccess)
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
		it.Event = new(TermsOfServiceAllowListAddedAccess)
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

func (it *TermsOfServiceAllowListAddedAccessIterator) Error() error {
	return it.fail
}

func (it *TermsOfServiceAllowListAddedAccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TermsOfServiceAllowListAddedAccess struct {
	User common.Address
	Raw  types.Log
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListFilterer) FilterAddedAccess(opts *bind.FilterOpts) (*TermsOfServiceAllowListAddedAccessIterator, error) {

	logs, sub, err := _TermsOfServiceAllowList.contract.FilterLogs(opts, "AddedAccess")
	if err != nil {
		return nil, err
	}
	return &TermsOfServiceAllowListAddedAccessIterator{contract: _TermsOfServiceAllowList.contract, event: "AddedAccess", logs: logs, sub: sub}, nil
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListFilterer) WatchAddedAccess(opts *bind.WatchOpts, sink chan<- *TermsOfServiceAllowListAddedAccess) (event.Subscription, error) {

	logs, sub, err := _TermsOfServiceAllowList.contract.WatchLogs(opts, "AddedAccess")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TermsOfServiceAllowListAddedAccess)
				if err := _TermsOfServiceAllowList.contract.UnpackLog(event, "AddedAccess", log); err != nil {
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

func (_TermsOfServiceAllowList *TermsOfServiceAllowListFilterer) ParseAddedAccess(log types.Log) (*TermsOfServiceAllowListAddedAccess, error) {
	event := new(TermsOfServiceAllowListAddedAccess)
	if err := _TermsOfServiceAllowList.contract.UnpackLog(event, "AddedAccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TermsOfServiceAllowListBlockedAccessIterator struct {
	Event *TermsOfServiceAllowListBlockedAccess

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TermsOfServiceAllowListBlockedAccessIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TermsOfServiceAllowListBlockedAccess)
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
		it.Event = new(TermsOfServiceAllowListBlockedAccess)
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

func (it *TermsOfServiceAllowListBlockedAccessIterator) Error() error {
	return it.fail
}

func (it *TermsOfServiceAllowListBlockedAccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TermsOfServiceAllowListBlockedAccess struct {
	User common.Address
	Raw  types.Log
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListFilterer) FilterBlockedAccess(opts *bind.FilterOpts) (*TermsOfServiceAllowListBlockedAccessIterator, error) {

	logs, sub, err := _TermsOfServiceAllowList.contract.FilterLogs(opts, "BlockedAccess")
	if err != nil {
		return nil, err
	}
	return &TermsOfServiceAllowListBlockedAccessIterator{contract: _TermsOfServiceAllowList.contract, event: "BlockedAccess", logs: logs, sub: sub}, nil
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListFilterer) WatchBlockedAccess(opts *bind.WatchOpts, sink chan<- *TermsOfServiceAllowListBlockedAccess) (event.Subscription, error) {

	logs, sub, err := _TermsOfServiceAllowList.contract.WatchLogs(opts, "BlockedAccess")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TermsOfServiceAllowListBlockedAccess)
				if err := _TermsOfServiceAllowList.contract.UnpackLog(event, "BlockedAccess", log); err != nil {
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

func (_TermsOfServiceAllowList *TermsOfServiceAllowListFilterer) ParseBlockedAccess(log types.Log) (*TermsOfServiceAllowListBlockedAccess, error) {
	event := new(TermsOfServiceAllowListBlockedAccess)
	if err := _TermsOfServiceAllowList.contract.UnpackLog(event, "BlockedAccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
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
	Enabled         bool
	SignerPublicKey common.Address
	Raw             types.Log
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

type TermsOfServiceAllowListUnblockedAccessIterator struct {
	Event *TermsOfServiceAllowListUnblockedAccess

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TermsOfServiceAllowListUnblockedAccessIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TermsOfServiceAllowListUnblockedAccess)
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
		it.Event = new(TermsOfServiceAllowListUnblockedAccess)
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

func (it *TermsOfServiceAllowListUnblockedAccessIterator) Error() error {
	return it.fail
}

func (it *TermsOfServiceAllowListUnblockedAccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TermsOfServiceAllowListUnblockedAccess struct {
	User common.Address
	Raw  types.Log
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListFilterer) FilterUnblockedAccess(opts *bind.FilterOpts) (*TermsOfServiceAllowListUnblockedAccessIterator, error) {

	logs, sub, err := _TermsOfServiceAllowList.contract.FilterLogs(opts, "UnblockedAccess")
	if err != nil {
		return nil, err
	}
	return &TermsOfServiceAllowListUnblockedAccessIterator{contract: _TermsOfServiceAllowList.contract, event: "UnblockedAccess", logs: logs, sub: sub}, nil
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListFilterer) WatchUnblockedAccess(opts *bind.WatchOpts, sink chan<- *TermsOfServiceAllowListUnblockedAccess) (event.Subscription, error) {

	logs, sub, err := _TermsOfServiceAllowList.contract.WatchLogs(opts, "UnblockedAccess")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TermsOfServiceAllowListUnblockedAccess)
				if err := _TermsOfServiceAllowList.contract.UnpackLog(event, "UnblockedAccess", log); err != nil {
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

func (_TermsOfServiceAllowList *TermsOfServiceAllowListFilterer) ParseUnblockedAccess(log types.Log) (*TermsOfServiceAllowListUnblockedAccess, error) {
	event := new(TermsOfServiceAllowListUnblockedAccess)
	if err := _TermsOfServiceAllowList.contract.UnpackLog(event, "UnblockedAccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowList) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _TermsOfServiceAllowList.abi.Events["AddedAccess"].ID:
		return _TermsOfServiceAllowList.ParseAddedAccess(log)
	case _TermsOfServiceAllowList.abi.Events["BlockedAccess"].ID:
		return _TermsOfServiceAllowList.ParseBlockedAccess(log)
	case _TermsOfServiceAllowList.abi.Events["ConfigSet"].ID:
		return _TermsOfServiceAllowList.ParseConfigSet(log)
	case _TermsOfServiceAllowList.abi.Events["UnblockedAccess"].ID:
		return _TermsOfServiceAllowList.ParseUnblockedAccess(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (TermsOfServiceAllowListAddedAccess) Topic() common.Hash {
	return common.HexToHash("0x87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db4")
}

func (TermsOfServiceAllowListBlockedAccess) Topic() common.Hash {
	return common.HexToHash("0x337cd0f3f594112b6d830afb510072d3b08556b446514f73b8109162fd1151e1")
}

func (TermsOfServiceAllowListConfigSet) Topic() common.Hash {
	return common.HexToHash("0x232bb132911bf7acb3a5ab4dd8e603080819db6a26b68da198599a7c83e2494c")
}

func (TermsOfServiceAllowListUnblockedAccess) Topic() common.Hash {
	return common.HexToHash("0x28bbd0761309a99e8fb5e5d02ada0b7b2db2e5357531ff5dbfc205c3f5b6592b")
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowList) Address() common.Address {
	return _TermsOfServiceAllowList.address
}

type TermsOfServiceAllowListInterface interface {
	GetAllAllowedSenders(opts *bind.CallOpts) ([]common.Address, error)

	GetMessage(opts *bind.CallOpts, acceptor common.Address, recipient common.Address) ([32]byte, error)

	HasAccess(opts *bind.CallOpts, user common.Address, arg1 []byte) (bool, error)

	IsBlockedSender(opts *bind.CallOpts, sender common.Address) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptTermsOfService(opts *bind.TransactOpts, acceptor common.Address, recipient common.Address, r [32]byte, s [32]byte, v uint8) (*types.Transaction, error)

	BlockSender(opts *bind.TransactOpts, sender common.Address) (*types.Transaction, error)

	UnblockSender(opts *bind.TransactOpts, sender common.Address) (*types.Transaction, error)

	UpdateConfig(opts *bind.TransactOpts, config []byte) (*types.Transaction, error)

	FilterAddedAccess(opts *bind.FilterOpts) (*TermsOfServiceAllowListAddedAccessIterator, error)

	WatchAddedAccess(opts *bind.WatchOpts, sink chan<- *TermsOfServiceAllowListAddedAccess) (event.Subscription, error)

	ParseAddedAccess(log types.Log) (*TermsOfServiceAllowListAddedAccess, error)

	FilterBlockedAccess(opts *bind.FilterOpts) (*TermsOfServiceAllowListBlockedAccessIterator, error)

	WatchBlockedAccess(opts *bind.WatchOpts, sink chan<- *TermsOfServiceAllowListBlockedAccess) (event.Subscription, error)

	ParseBlockedAccess(log types.Log) (*TermsOfServiceAllowListBlockedAccess, error)

	FilterConfigSet(opts *bind.FilterOpts) (*TermsOfServiceAllowListConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *TermsOfServiceAllowListConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*TermsOfServiceAllowListConfigSet, error)

	FilterUnblockedAccess(opts *bind.FilterOpts) (*TermsOfServiceAllowListUnblockedAccessIterator, error)

	WatchUnblockedAccess(opts *bind.WatchOpts, sink chan<- *TermsOfServiceAllowListUnblockedAccess) (event.Subscription, error)

	ParseUnblockedAccess(log types.Log) (*TermsOfServiceAllowListUnblockedAccess, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
