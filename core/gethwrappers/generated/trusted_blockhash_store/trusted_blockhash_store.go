// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package trusted_blockhash_store

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

var TrustedBlockhashStoreMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"whitelist\",\"type\":\"address[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidRecentBlockhash\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTrustedBlockhashes\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInWhitelist\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getBlockhash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_whitelist\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"s_whitelistStatus\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"whitelist\",\"type\":\"address[]\"}],\"name\":\"setWhitelist\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"store\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"storeEarliest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"blockNums\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"blockhashes\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256\",\"name\":\"recentBlockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"recentBlockhash\",\"type\":\"bytes32\"}],\"name\":\"storeTrusted\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"header\",\"type\":\"bytes\"}],\"name\":\"storeVerifyHeader\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b50604051620014b8380380620014b88339810160408190526200003491620003fd565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000d9565b505050620000d2816200018460201b60201c565b506200050d565b336001600160a01b03821603620001335760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6200018e620002eb565b60006004805480602002602001604051908101604052809291908181526020018280548015620001e857602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311620001c9575b50508551939450620002069360049350602087019250905062000349565b5060005b815181101562000276576000600360008484815181106200022f576200022f620004cf565b6020908102919091018101516001600160a01b03168252810191909152604001600020805460ff1916911515919091179055806200026d81620004e5565b9150506200020a565b5060005b8251811015620002e6576001600360008584815181106200029f576200029f620004cf565b6020908102919091018101516001600160a01b03168252810191909152604001600020805460ff191691151591909117905580620002dd81620004e5565b9150506200027a565b505050565b6000546001600160a01b03163314620003475760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640162000082565b565b828054828255906000526020600020908101928215620003a1579160200282015b82811115620003a157825182546001600160a01b0319166001600160a01b039091161782556020909201916001909101906200036a565b50620003af929150620003b3565b5090565b5b80821115620003af5760008155600101620003b4565b634e487b7160e01b600052604160045260246000fd5b80516001600160a01b0381168114620003f857600080fd5b919050565b600060208083850312156200041157600080fd5b82516001600160401b03808211156200042957600080fd5b818501915085601f8301126200043e57600080fd5b815181811115620004535762000453620003ca565b8060051b604051601f19603f830116810181811085821117156200047b576200047b620003ca565b6040529182528482019250838101850191888311156200049a57600080fd5b938501935b82851015620004c357620004b385620003e0565b845293850193928501926200049f565b98975050505050505050565b634e487b7160e01b600052603260045260246000fd5b6000600182016200050657634e487b7160e01b600052601160045260246000fd5b5060010190565b610f9b806200051d6000396000f3fe608060405234801561001057600080fd5b50600436106100c95760003560e01c80638da5cb5b11610081578063f2fde38b1161005b578063f2fde38b146101b5578063f4217648146101c8578063fadff0e1146101db57600080fd5b80638da5cb5b14610143578063e9413d3814610161578063e9ecc1541461018257600080fd5b80636057361d116100b25780636057361d1461012057806379ba50971461013357806383b6d6b71461013b57600080fd5b80633b69ad60146100ce5780635c7de309146100e3575b600080fd5b6100e16100dc366004610bf6565b6101ee565b005b6100f66100f1366004610c74565b610326565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b6100e161012e366004610c74565b61035d565b6100e16103ec565b6100e16104e9565b60005473ffffffffffffffffffffffffffffffffffffffff166100f6565b61017461016f366004610c74565b610503565b604051908152602001610117565b6101a5610190366004610cb6565b60036020526000908152604090205460ff1681565b6040519015158152602001610117565b6100e16101c3366004610cb6565b610581565b6100e16101d6366004610d4f565b610595565b6100e16101e9366004610dfc565b61074b565b60006101f9836107ee565b9050818114610234576040517fd2f69c9500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3360009081526003602052604090205460ff1661027d576040517f5b0aa2ba00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8584146102b6576040517fbd75093300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b8681101561031c578585828181106102d3576102d3610eb9565b90506020020135600260008a8a858181106102f0576102f0610eb9565b90506020020135815260200190815260200160002081905550808061031490610f17565b9150506102b9565b5050505050505050565b6004818154811061033657600080fd5b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff16905081565b6000610368826107ee565b905060008190036103da576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f626c6f636b68617368286e29206661696c65640000000000000000000000000060448201526064015b60405180910390fd5b60009182526002602052604090912055565b60015473ffffffffffffffffffffffffffffffffffffffff16331461046d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016103d1565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6105016101006104f76108e4565b61012e9190610f4f565b565b60008181526002602052604081205480820361057b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f626c6f636b68617368206e6f7420666f756e6420696e2073746f72650000000060448201526064016103d1565b92915050565b610589610972565b610592816109f3565b50565b61059d610972565b6000600480548060200260200160405190810160405280929190818152602001828054801561060257602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116105d7575b5050855193945061061e93600493506020870192509050610b0b565b5060005b81518110156106b25760006003600084848151811061064357610643610eb9565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff16825281019190915260400160002080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055806106aa81610f17565b915050610622565b5060005b8251811015610746576001600360008584815181106106d7576106d7610eb9565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff16825281019190915260400160002080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169115159190911790558061073e81610f17565b9150506106b6565b505050565b6002600061075a846001610f62565b8152602001908152602001600020548180519060200120146107d8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f6865616465722068617320756e6b6e6f776e20626c6f636b686173680000000060448201526064016103d1565b6024015160009182526002602052604090912055565b6000466107fa81610ae8565b156108d4576101008367ffffffffffffffff166108156108e4565b61081f9190610f4f565b118061083c575061082e6108e4565b8367ffffffffffffffff1610155b1561084a5750600092915050565b6040517f2b407a8200000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84166004820152606490632b407a8290602401602060405180830381865afa1580156108a9573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108cd9190610f75565b9392505050565b505067ffffffffffffffff164090565b6000466108f081610ae8565b1561096b57606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610941573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906109659190610f75565b91505090565b4391505090565b60005473ffffffffffffffffffffffffffffffffffffffff163314610501576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016103d1565b3373ffffffffffffffffffffffffffffffffffffffff821603610a72576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016103d1565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600061a4b1821480610afc575062066eed82145b8061057b57505062066eee1490565b828054828255906000526020600020908101928215610b85579160200282015b82811115610b8557825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190610b2b565b50610b91929150610b95565b5090565b5b80821115610b915760008155600101610b96565b60008083601f840112610bbc57600080fd5b50813567ffffffffffffffff811115610bd457600080fd5b6020830191508360208260051b8501011115610bef57600080fd5b9250929050565b60008060008060008060808789031215610c0f57600080fd5b863567ffffffffffffffff80821115610c2757600080fd5b610c338a838b01610baa565b90985096506020890135915080821115610c4c57600080fd5b50610c5989828a01610baa565b979a9699509760408101359660609091013595509350505050565b600060208284031215610c8657600080fd5b5035919050565b803573ffffffffffffffffffffffffffffffffffffffff81168114610cb157600080fd5b919050565b600060208284031215610cc857600080fd5b6108cd82610c8d565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610d4757610d47610cd1565b604052919050565b60006020808385031215610d6257600080fd5b823567ffffffffffffffff80821115610d7a57600080fd5b818501915085601f830112610d8e57600080fd5b813581811115610da057610da0610cd1565b8060051b9150610db1848301610d00565b8181529183018401918481019088841115610dcb57600080fd5b938501935b83851015610df057610de185610c8d565b82529385019390850190610dd0565b98975050505050505050565b60008060408385031215610e0f57600080fd5b8235915060208084013567ffffffffffffffff80821115610e2f57600080fd5b818601915086601f830112610e4357600080fd5b813581811115610e5557610e55610cd1565b610e85847fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601610d00565b91508082528784828501011115610e9b57600080fd5b80848401858401376000848284010152508093505050509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203610f4857610f48610ee8565b5060010190565b8181038181111561057b5761057b610ee8565b8082018082111561057b5761057b610ee8565b600060208284031215610f8757600080fd5b505191905056fea164736f6c6343000813000a",
}

var TrustedBlockhashStoreABI = TrustedBlockhashStoreMetaData.ABI

var TrustedBlockhashStoreBin = TrustedBlockhashStoreMetaData.Bin

func DeployTrustedBlockhashStore(auth *bind.TransactOpts, backend bind.ContractBackend, whitelist []common.Address) (common.Address, *types.Transaction, *TrustedBlockhashStore, error) {
	parsed, err := TrustedBlockhashStoreMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TrustedBlockhashStoreBin), backend, whitelist)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TrustedBlockhashStore{address: address, abi: *parsed, TrustedBlockhashStoreCaller: TrustedBlockhashStoreCaller{contract: contract}, TrustedBlockhashStoreTransactor: TrustedBlockhashStoreTransactor{contract: contract}, TrustedBlockhashStoreFilterer: TrustedBlockhashStoreFilterer{contract: contract}}, nil
}

type TrustedBlockhashStore struct {
	address common.Address
	abi     abi.ABI
	TrustedBlockhashStoreCaller
	TrustedBlockhashStoreTransactor
	TrustedBlockhashStoreFilterer
}

type TrustedBlockhashStoreCaller struct {
	contract *bind.BoundContract
}

type TrustedBlockhashStoreTransactor struct {
	contract *bind.BoundContract
}

type TrustedBlockhashStoreFilterer struct {
	contract *bind.BoundContract
}

type TrustedBlockhashStoreSession struct {
	Contract     *TrustedBlockhashStore
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type TrustedBlockhashStoreCallerSession struct {
	Contract *TrustedBlockhashStoreCaller
	CallOpts bind.CallOpts
}

type TrustedBlockhashStoreTransactorSession struct {
	Contract     *TrustedBlockhashStoreTransactor
	TransactOpts bind.TransactOpts
}

type TrustedBlockhashStoreRaw struct {
	Contract *TrustedBlockhashStore
}

type TrustedBlockhashStoreCallerRaw struct {
	Contract *TrustedBlockhashStoreCaller
}

type TrustedBlockhashStoreTransactorRaw struct {
	Contract *TrustedBlockhashStoreTransactor
}

func NewTrustedBlockhashStore(address common.Address, backend bind.ContractBackend) (*TrustedBlockhashStore, error) {
	abi, err := abi.JSON(strings.NewReader(TrustedBlockhashStoreABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindTrustedBlockhashStore(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TrustedBlockhashStore{address: address, abi: abi, TrustedBlockhashStoreCaller: TrustedBlockhashStoreCaller{contract: contract}, TrustedBlockhashStoreTransactor: TrustedBlockhashStoreTransactor{contract: contract}, TrustedBlockhashStoreFilterer: TrustedBlockhashStoreFilterer{contract: contract}}, nil
}

func NewTrustedBlockhashStoreCaller(address common.Address, caller bind.ContractCaller) (*TrustedBlockhashStoreCaller, error) {
	contract, err := bindTrustedBlockhashStore(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TrustedBlockhashStoreCaller{contract: contract}, nil
}

func NewTrustedBlockhashStoreTransactor(address common.Address, transactor bind.ContractTransactor) (*TrustedBlockhashStoreTransactor, error) {
	contract, err := bindTrustedBlockhashStore(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TrustedBlockhashStoreTransactor{contract: contract}, nil
}

func NewTrustedBlockhashStoreFilterer(address common.Address, filterer bind.ContractFilterer) (*TrustedBlockhashStoreFilterer, error) {
	contract, err := bindTrustedBlockhashStore(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TrustedBlockhashStoreFilterer{contract: contract}, nil
}

func bindTrustedBlockhashStore(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TrustedBlockhashStoreMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_TrustedBlockhashStore *TrustedBlockhashStoreRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TrustedBlockhashStore.Contract.TrustedBlockhashStoreCaller.contract.Call(opts, result, method, params...)
}

func (_TrustedBlockhashStore *TrustedBlockhashStoreRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TrustedBlockhashStore.Contract.TrustedBlockhashStoreTransactor.contract.Transfer(opts)
}

func (_TrustedBlockhashStore *TrustedBlockhashStoreRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TrustedBlockhashStore.Contract.TrustedBlockhashStoreTransactor.contract.Transact(opts, method, params...)
}

func (_TrustedBlockhashStore *TrustedBlockhashStoreCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TrustedBlockhashStore.Contract.contract.Call(opts, result, method, params...)
}

func (_TrustedBlockhashStore *TrustedBlockhashStoreTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TrustedBlockhashStore.Contract.contract.Transfer(opts)
}

func (_TrustedBlockhashStore *TrustedBlockhashStoreTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TrustedBlockhashStore.Contract.contract.Transact(opts, method, params...)
}

func (_TrustedBlockhashStore *TrustedBlockhashStoreCaller) GetBlockhash(opts *bind.CallOpts, n *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _TrustedBlockhashStore.contract.Call(opts, &out, "getBlockhash", n)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_TrustedBlockhashStore *TrustedBlockhashStoreSession) GetBlockhash(n *big.Int) ([32]byte, error) {
	return _TrustedBlockhashStore.Contract.GetBlockhash(&_TrustedBlockhashStore.CallOpts, n)
}

func (_TrustedBlockhashStore *TrustedBlockhashStoreCallerSession) GetBlockhash(n *big.Int) ([32]byte, error) {
	return _TrustedBlockhashStore.Contract.GetBlockhash(&_TrustedBlockhashStore.CallOpts, n)
}

func (_TrustedBlockhashStore *TrustedBlockhashStoreCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TrustedBlockhashStore.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_TrustedBlockhashStore *TrustedBlockhashStoreSession) Owner() (common.Address, error) {
	return _TrustedBlockhashStore.Contract.Owner(&_TrustedBlockhashStore.CallOpts)
}

func (_TrustedBlockhashStore *TrustedBlockhashStoreCallerSession) Owner() (common.Address, error) {
	return _TrustedBlockhashStore.Contract.Owner(&_TrustedBlockhashStore.CallOpts)
}

func (_TrustedBlockhashStore *TrustedBlockhashStoreCaller) SWhitelist(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _TrustedBlockhashStore.contract.Call(opts, &out, "s_whitelist", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_TrustedBlockhashStore *TrustedBlockhashStoreSession) SWhitelist(arg0 *big.Int) (common.Address, error) {
	return _TrustedBlockhashStore.Contract.SWhitelist(&_TrustedBlockhashStore.CallOpts, arg0)
}

func (_TrustedBlockhashStore *TrustedBlockhashStoreCallerSession) SWhitelist(arg0 *big.Int) (common.Address, error) {
	return _TrustedBlockhashStore.Contract.SWhitelist(&_TrustedBlockhashStore.CallOpts, arg0)
}

func (_TrustedBlockhashStore *TrustedBlockhashStoreCaller) SWhitelistStatus(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _TrustedBlockhashStore.contract.Call(opts, &out, "s_whitelistStatus", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_TrustedBlockhashStore *TrustedBlockhashStoreSession) SWhitelistStatus(arg0 common.Address) (bool, error) {
	return _TrustedBlockhashStore.Contract.SWhitelistStatus(&_TrustedBlockhashStore.CallOpts, arg0)
}

func (_TrustedBlockhashStore *TrustedBlockhashStoreCallerSession) SWhitelistStatus(arg0 common.Address) (bool, error) {
	return _TrustedBlockhashStore.Contract.SWhitelistStatus(&_TrustedBlockhashStore.CallOpts, arg0)
}

func (_TrustedBlockhashStore *TrustedBlockhashStoreTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TrustedBlockhashStore.contract.Transact(opts, "acceptOwnership")
}

func (_TrustedBlockhashStore *TrustedBlockhashStoreSession) AcceptOwnership() (*types.Transaction, error) {
	return _TrustedBlockhashStore.Contract.AcceptOwnership(&_TrustedBlockhashStore.TransactOpts)
}

func (_TrustedBlockhashStore *TrustedBlockhashStoreTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _TrustedBlockhashStore.Contract.AcceptOwnership(&_TrustedBlockhashStore.TransactOpts)
}

func (_TrustedBlockhashStore *TrustedBlockhashStoreTransactor) SetWhitelist(opts *bind.TransactOpts, whitelist []common.Address) (*types.Transaction, error) {
	return _TrustedBlockhashStore.contract.Transact(opts, "setWhitelist", whitelist)
}

func (_TrustedBlockhashStore *TrustedBlockhashStoreSession) SetWhitelist(whitelist []common.Address) (*types.Transaction, error) {
	return _TrustedBlockhashStore.Contract.SetWhitelist(&_TrustedBlockhashStore.TransactOpts, whitelist)
}

func (_TrustedBlockhashStore *TrustedBlockhashStoreTransactorSession) SetWhitelist(whitelist []common.Address) (*types.Transaction, error) {
	return _TrustedBlockhashStore.Contract.SetWhitelist(&_TrustedBlockhashStore.TransactOpts, whitelist)
}

func (_TrustedBlockhashStore *TrustedBlockhashStoreTransactor) Store(opts *bind.TransactOpts, n *big.Int) (*types.Transaction, error) {
	return _TrustedBlockhashStore.contract.Transact(opts, "store", n)
}

func (_TrustedBlockhashStore *TrustedBlockhashStoreSession) Store(n *big.Int) (*types.Transaction, error) {
	return _TrustedBlockhashStore.Contract.Store(&_TrustedBlockhashStore.TransactOpts, n)
}

func (_TrustedBlockhashStore *TrustedBlockhashStoreTransactorSession) Store(n *big.Int) (*types.Transaction, error) {
	return _TrustedBlockhashStore.Contract.Store(&_TrustedBlockhashStore.TransactOpts, n)
}

func (_TrustedBlockhashStore *TrustedBlockhashStoreTransactor) StoreEarliest(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TrustedBlockhashStore.contract.Transact(opts, "storeEarliest")
}

func (_TrustedBlockhashStore *TrustedBlockhashStoreSession) StoreEarliest() (*types.Transaction, error) {
	return _TrustedBlockhashStore.Contract.StoreEarliest(&_TrustedBlockhashStore.TransactOpts)
}

func (_TrustedBlockhashStore *TrustedBlockhashStoreTransactorSession) StoreEarliest() (*types.Transaction, error) {
	return _TrustedBlockhashStore.Contract.StoreEarliest(&_TrustedBlockhashStore.TransactOpts)
}

func (_TrustedBlockhashStore *TrustedBlockhashStoreTransactor) StoreTrusted(opts *bind.TransactOpts, blockNums []*big.Int, blockhashes [][32]byte, recentBlockNumber *big.Int, recentBlockhash [32]byte) (*types.Transaction, error) {
	return _TrustedBlockhashStore.contract.Transact(opts, "storeTrusted", blockNums, blockhashes, recentBlockNumber, recentBlockhash)
}

func (_TrustedBlockhashStore *TrustedBlockhashStoreSession) StoreTrusted(blockNums []*big.Int, blockhashes [][32]byte, recentBlockNumber *big.Int, recentBlockhash [32]byte) (*types.Transaction, error) {
	return _TrustedBlockhashStore.Contract.StoreTrusted(&_TrustedBlockhashStore.TransactOpts, blockNums, blockhashes, recentBlockNumber, recentBlockhash)
}

func (_TrustedBlockhashStore *TrustedBlockhashStoreTransactorSession) StoreTrusted(blockNums []*big.Int, blockhashes [][32]byte, recentBlockNumber *big.Int, recentBlockhash [32]byte) (*types.Transaction, error) {
	return _TrustedBlockhashStore.Contract.StoreTrusted(&_TrustedBlockhashStore.TransactOpts, blockNums, blockhashes, recentBlockNumber, recentBlockhash)
}

func (_TrustedBlockhashStore *TrustedBlockhashStoreTransactor) StoreVerifyHeader(opts *bind.TransactOpts, n *big.Int, header []byte) (*types.Transaction, error) {
	return _TrustedBlockhashStore.contract.Transact(opts, "storeVerifyHeader", n, header)
}

func (_TrustedBlockhashStore *TrustedBlockhashStoreSession) StoreVerifyHeader(n *big.Int, header []byte) (*types.Transaction, error) {
	return _TrustedBlockhashStore.Contract.StoreVerifyHeader(&_TrustedBlockhashStore.TransactOpts, n, header)
}

func (_TrustedBlockhashStore *TrustedBlockhashStoreTransactorSession) StoreVerifyHeader(n *big.Int, header []byte) (*types.Transaction, error) {
	return _TrustedBlockhashStore.Contract.StoreVerifyHeader(&_TrustedBlockhashStore.TransactOpts, n, header)
}

func (_TrustedBlockhashStore *TrustedBlockhashStoreTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _TrustedBlockhashStore.contract.Transact(opts, "transferOwnership", to)
}

func (_TrustedBlockhashStore *TrustedBlockhashStoreSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _TrustedBlockhashStore.Contract.TransferOwnership(&_TrustedBlockhashStore.TransactOpts, to)
}

func (_TrustedBlockhashStore *TrustedBlockhashStoreTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _TrustedBlockhashStore.Contract.TransferOwnership(&_TrustedBlockhashStore.TransactOpts, to)
}

type TrustedBlockhashStoreOwnershipTransferRequestedIterator struct {
	Event *TrustedBlockhashStoreOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TrustedBlockhashStoreOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TrustedBlockhashStoreOwnershipTransferRequested)
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
		it.Event = new(TrustedBlockhashStoreOwnershipTransferRequested)
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

func (it *TrustedBlockhashStoreOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *TrustedBlockhashStoreOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TrustedBlockhashStoreOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_TrustedBlockhashStore *TrustedBlockhashStoreFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TrustedBlockhashStoreOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TrustedBlockhashStore.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &TrustedBlockhashStoreOwnershipTransferRequestedIterator{contract: _TrustedBlockhashStore.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_TrustedBlockhashStore *TrustedBlockhashStoreFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *TrustedBlockhashStoreOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TrustedBlockhashStore.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TrustedBlockhashStoreOwnershipTransferRequested)
				if err := _TrustedBlockhashStore.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_TrustedBlockhashStore *TrustedBlockhashStoreFilterer) ParseOwnershipTransferRequested(log types.Log) (*TrustedBlockhashStoreOwnershipTransferRequested, error) {
	event := new(TrustedBlockhashStoreOwnershipTransferRequested)
	if err := _TrustedBlockhashStore.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TrustedBlockhashStoreOwnershipTransferredIterator struct {
	Event *TrustedBlockhashStoreOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TrustedBlockhashStoreOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TrustedBlockhashStoreOwnershipTransferred)
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
		it.Event = new(TrustedBlockhashStoreOwnershipTransferred)
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

func (it *TrustedBlockhashStoreOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *TrustedBlockhashStoreOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TrustedBlockhashStoreOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_TrustedBlockhashStore *TrustedBlockhashStoreFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TrustedBlockhashStoreOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TrustedBlockhashStore.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &TrustedBlockhashStoreOwnershipTransferredIterator{contract: _TrustedBlockhashStore.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_TrustedBlockhashStore *TrustedBlockhashStoreFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TrustedBlockhashStoreOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TrustedBlockhashStore.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TrustedBlockhashStoreOwnershipTransferred)
				if err := _TrustedBlockhashStore.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_TrustedBlockhashStore *TrustedBlockhashStoreFilterer) ParseOwnershipTransferred(log types.Log) (*TrustedBlockhashStoreOwnershipTransferred, error) {
	event := new(TrustedBlockhashStoreOwnershipTransferred)
	if err := _TrustedBlockhashStore.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_TrustedBlockhashStore *TrustedBlockhashStore) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _TrustedBlockhashStore.abi.Events["OwnershipTransferRequested"].ID:
		return _TrustedBlockhashStore.ParseOwnershipTransferRequested(log)
	case _TrustedBlockhashStore.abi.Events["OwnershipTransferred"].ID:
		return _TrustedBlockhashStore.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (TrustedBlockhashStoreOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (TrustedBlockhashStoreOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_TrustedBlockhashStore *TrustedBlockhashStore) Address() common.Address {
	return _TrustedBlockhashStore.address
}

type TrustedBlockhashStoreInterface interface {
	GetBlockhash(opts *bind.CallOpts, n *big.Int) ([32]byte, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SWhitelist(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error)

	SWhitelistStatus(opts *bind.CallOpts, arg0 common.Address) (bool, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	SetWhitelist(opts *bind.TransactOpts, whitelist []common.Address) (*types.Transaction, error)

	Store(opts *bind.TransactOpts, n *big.Int) (*types.Transaction, error)

	StoreEarliest(opts *bind.TransactOpts) (*types.Transaction, error)

	StoreTrusted(opts *bind.TransactOpts, blockNums []*big.Int, blockhashes [][32]byte, recentBlockNumber *big.Int, recentBlockhash [32]byte) (*types.Transaction, error)

	StoreVerifyHeader(opts *bind.TransactOpts, n *big.Int, header []byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TrustedBlockhashStoreOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *TrustedBlockhashStoreOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*TrustedBlockhashStoreOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TrustedBlockhashStoreOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TrustedBlockhashStoreOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*TrustedBlockhashStoreOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
