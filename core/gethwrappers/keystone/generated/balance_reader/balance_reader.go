// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package balance_reader

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
	_ = abi.ConvertType
)

var BalanceReaderMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"getNativeBalances\",\"inputs\":[{\"name\":\"addresses\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"typeAndVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610388806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c8063181f5a771461003b5780634c04bf991461008d575b600080fd5b6100776040518060400160405280601381526020017f42616c616e636552656164657220312e302e300000000000000000000000000081525081565b6040516100849190610160565b60405180910390f35b6100a061009b366004610225565b6100ad565b6040516100849190610308565b60606000825167ffffffffffffffff8111156100cb576100cb6101cd565b6040519080825280602002602001820160405280156100f4578160200160208202803683370190505b50905060005b8351811015610159578381815181106101155761011561034c565b602002602001015173ffffffffffffffffffffffffffffffffffffffff16318282815181106101465761014661034c565b60209081029190910101526001016100fa565b5092915050565b60006020808352835180602085015260005b8181101561018e57858101830151858201604001528201610172565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b803573ffffffffffffffffffffffffffffffffffffffff8116811461022057600080fd5b919050565b6000602080838503121561023857600080fd5b823567ffffffffffffffff8082111561025057600080fd5b818501915085601f83011261026457600080fd5b813581811115610276576102766101cd565b8060051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f830116810181811085821117156102b9576102b96101cd565b6040529182528482019250838101850191888311156102d757600080fd5b938501935b828510156102fc576102ed856101fc565b845293850193928501926102dc565b98975050505050505050565b6020808252825182820181905260009190848201906040850190845b8181101561034057835183529284019291840191600101610324565b50909695505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fdfea164736f6c6343000818000a",
}

var BalanceReaderABI = BalanceReaderMetaData.ABI

var BalanceReaderBin = BalanceReaderMetaData.Bin

func DeployBalanceReader(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *BalanceReader, error) {
	parsed, err := BalanceReaderMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BalanceReaderBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BalanceReader{address: address, abi: *parsed, BalanceReaderCaller: BalanceReaderCaller{contract: contract}, BalanceReaderTransactor: BalanceReaderTransactor{contract: contract}, BalanceReaderFilterer: BalanceReaderFilterer{contract: contract}}, nil
}

type BalanceReader struct {
	address common.Address
	abi     abi.ABI
	BalanceReaderCaller
	BalanceReaderTransactor
	BalanceReaderFilterer
}

type BalanceReaderCaller struct {
	contract *bind.BoundContract
}

type BalanceReaderTransactor struct {
	contract *bind.BoundContract
}

type BalanceReaderFilterer struct {
	contract *bind.BoundContract
}

type BalanceReaderSession struct {
	Contract     *BalanceReader
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type BalanceReaderCallerSession struct {
	Contract *BalanceReaderCaller
	CallOpts bind.CallOpts
}

type BalanceReaderTransactorSession struct {
	Contract     *BalanceReaderTransactor
	TransactOpts bind.TransactOpts
}

type BalanceReaderRaw struct {
	Contract *BalanceReader
}

type BalanceReaderCallerRaw struct {
	Contract *BalanceReaderCaller
}

type BalanceReaderTransactorRaw struct {
	Contract *BalanceReaderTransactor
}

func NewBalanceReader(address common.Address, backend bind.ContractBackend) (*BalanceReader, error) {
	abi, err := abi.JSON(strings.NewReader(BalanceReaderABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindBalanceReader(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BalanceReader{address: address, abi: abi, BalanceReaderCaller: BalanceReaderCaller{contract: contract}, BalanceReaderTransactor: BalanceReaderTransactor{contract: contract}, BalanceReaderFilterer: BalanceReaderFilterer{contract: contract}}, nil
}

func NewBalanceReaderCaller(address common.Address, caller bind.ContractCaller) (*BalanceReaderCaller, error) {
	contract, err := bindBalanceReader(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BalanceReaderCaller{contract: contract}, nil
}

func NewBalanceReaderTransactor(address common.Address, transactor bind.ContractTransactor) (*BalanceReaderTransactor, error) {
	contract, err := bindBalanceReader(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BalanceReaderTransactor{contract: contract}, nil
}

func NewBalanceReaderFilterer(address common.Address, filterer bind.ContractFilterer) (*BalanceReaderFilterer, error) {
	contract, err := bindBalanceReader(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BalanceReaderFilterer{contract: contract}, nil
}

func bindBalanceReader(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BalanceReaderMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_BalanceReader *BalanceReaderRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BalanceReader.Contract.BalanceReaderCaller.contract.Call(opts, result, method, params...)
}

func (_BalanceReader *BalanceReaderRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BalanceReader.Contract.BalanceReaderTransactor.contract.Transfer(opts)
}

func (_BalanceReader *BalanceReaderRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BalanceReader.Contract.BalanceReaderTransactor.contract.Transact(opts, method, params...)
}

func (_BalanceReader *BalanceReaderCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BalanceReader.Contract.contract.Call(opts, result, method, params...)
}

func (_BalanceReader *BalanceReaderTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BalanceReader.Contract.contract.Transfer(opts)
}

func (_BalanceReader *BalanceReaderTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BalanceReader.Contract.contract.Transact(opts, method, params...)
}

func (_BalanceReader *BalanceReaderCaller) GetNativeBalances(opts *bind.CallOpts, addresses []common.Address) ([]*big.Int, error) {
	var out []interface{}
	err := _BalanceReader.contract.Call(opts, &out, "getNativeBalances", addresses)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_BalanceReader *BalanceReaderSession) GetNativeBalances(addresses []common.Address) ([]*big.Int, error) {
	return _BalanceReader.Contract.GetNativeBalances(&_BalanceReader.CallOpts, addresses)
}

func (_BalanceReader *BalanceReaderCallerSession) GetNativeBalances(addresses []common.Address) ([]*big.Int, error) {
	return _BalanceReader.Contract.GetNativeBalances(&_BalanceReader.CallOpts, addresses)
}

func (_BalanceReader *BalanceReaderCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _BalanceReader.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_BalanceReader *BalanceReaderSession) TypeAndVersion() (string, error) {
	return _BalanceReader.Contract.TypeAndVersion(&_BalanceReader.CallOpts)
}

func (_BalanceReader *BalanceReaderCallerSession) TypeAndVersion() (string, error) {
	return _BalanceReader.Contract.TypeAndVersion(&_BalanceReader.CallOpts)
}

func (_BalanceReader *BalanceReader) Address() common.Address {
	return _BalanceReader.address
}

type BalanceReaderInterface interface {
	GetNativeBalances(opts *bind.CallOpts, addresses []common.Address) ([]*big.Int, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	Address() common.Address
}
