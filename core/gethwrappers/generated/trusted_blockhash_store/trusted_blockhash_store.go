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
	ABI: "[{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"whitelist\",\"type\":\"address[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidRecentBlockhash\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTrustedBlockhashes\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInWhitelist\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getBlockhash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_blockhashes\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_whitelist\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"s_whitelistStatus\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"whitelist\",\"type\":\"address[]\"}],\"name\":\"setWhitelist\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"store\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"blockNums\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"blockhashes\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256\",\"name\":\"recentBlockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"recentBlockhash\",\"type\":\"bytes32\"}],\"name\":\"storeTrusted\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b506040516200104c3803806200104c833981016040819052620000349162000228565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000de565b50508151620000d6915060049060208401906200018a565b505062000317565b6001600160a01b038116331415620001395760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b828054828255906000526020600020908101928215620001e2579160200282015b82811115620001e257825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190620001ab565b50620001f0929150620001f4565b5090565b5b80821115620001f05760008155600101620001f5565b80516001600160a01b03811681146200022357600080fd5b919050565b600060208083850312156200023c57600080fd5b82516001600160401b03808211156200025457600080fd5b818501915085601f8301126200026957600080fd5b8151818111156200027e576200027e62000301565b8060051b604051601f19603f83011681018181108582111715620002a657620002a662000301565b604052828152858101935084860182860187018a1015620002c657600080fd5b600095505b83861015620002f457620002df816200020b565b855260019590950194938601938601620002cb565b5098975050505050505050565b634e487b7160e01b600052604160045260246000fd5b610d2580620003276000396000f3fe608060405234801561001057600080fd5b50600436106100be5760003560e01c80638da5cb5b11610076578063e9ecc1541161005b578063e9ecc1541461018f578063f2fde38b146101c2578063f4217648146101d557600080fd5b80638da5cb5b1461015e578063e9413d381461017c57600080fd5b80635c7de309116100a75780635c7de3091461010b5780636057361d1461014357806379ba50971461015657600080fd5b8063352633dd146100c35780633b69ad60146100f6575b600080fd5b6100e36100d1366004610c51565b60026020526000908152604090205481565b6040519081526020015b60405180910390f35b610109610104366004610bba565b6101e8565b005b61011e610119366004610c51565b610320565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100ed565b610109610151366004610c51565b610357565b6101096103e2565b60005473ffffffffffffffffffffffffffffffffffffffff1661011e565b6100e361018a366004610c51565b6104df565b6101b261019d366004610b42565b60036020526000908152604090205460ff1681565b60405190151581526020016100ed565b6101096101d0366004610b42565b61055b565b6101096101e3366004610b78565b61056f565b60006101f38361072c565b905081811461022e576040517fd2f69c9500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3360009081526003602052604090205460ff16610277576040517f5b0aa2ba00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8584146102b0576040517fbd75093300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b86811015610316578585828181106102cd576102cd610ce9565b90506020020135600260008a8a858181106102ea576102ea610ce9565b90506020020135815260200190815260200160002081905550808061030e90610c81565b9150506102b3565b5050505050505050565b6004818154811061033057600080fd5b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff16905081565b60006103628261072c565b9050806103d0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f626c6f636b68617368286e29206661696c65640000000000000000000000000060448201526064015b60405180910390fd5b60009182526002602052604090912055565b60015473ffffffffffffffffffffffffffffffffffffffff163314610463576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016103c7565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60008181526002602052604081205480610555576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f626c6f636b68617368206e6f7420666f756e6420696e2073746f72650000000060448201526064016103c7565b92915050565b61056361083a565b61056c816108bd565b50565b61057761083a565b600060048054806020026020016040519081016040528092919081815260200182805480156105dc57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116105b1575b505050505090508282600491906105f4929190610a59565b5060005b81518110156106885760006003600084848151811061061957610619610ce9565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff16825281019190915260400160002080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169115159190911790558061068081610c81565b9150506105f8565b5060005b82811015610726576001600360008686858181106106ac576106ac610ce9565b90506020020160208101906106c19190610b42565b73ffffffffffffffffffffffffffffffffffffffff168152602081019190915260400160002080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169115159190911790558061071e81610c81565b91505061068c565b50505050565b60004661a4b1811480610741575062066eed81145b1561082a576101008367ffffffffffffffff1661075c6109b3565b6107669190610c6a565b118061078357506107756109b3565b8367ffffffffffffffff1610155b156107915750600092915050565b6040517f2b407a8200000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84166004820152606490632b407a829060240160206040518083038186803b1580156107eb57600080fd5b505afa1580156107ff573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108239190610c38565b9392505050565b505067ffffffffffffffff164090565b60005473ffffffffffffffffffffffffffffffffffffffff1633146108bb576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016103c7565b565b73ffffffffffffffffffffffffffffffffffffffff811633141561093d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016103c7565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60004661a4b18114806109c8575062066eed81145b15610a5257606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b815260040160206040518083038186803b158015610a1457600080fd5b505afa158015610a28573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a4c9190610c38565b91505090565b4391505090565b828054828255906000526020600020908101928215610ad1579160200282015b82811115610ad15781547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff843516178255602090920191600190910190610a79565b50610add929150610ae1565b5090565b5b80821115610add5760008155600101610ae2565b60008083601f840112610b0857600080fd5b50813567ffffffffffffffff811115610b2057600080fd5b6020830191508360208260051b8501011115610b3b57600080fd5b9250929050565b600060208284031215610b5457600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461082357600080fd5b60008060208385031215610b8b57600080fd5b823567ffffffffffffffff811115610ba257600080fd5b610bae85828601610af6565b90969095509350505050565b60008060008060008060808789031215610bd357600080fd5b863567ffffffffffffffff80821115610beb57600080fd5b610bf78a838b01610af6565b90985096506020890135915080821115610c1057600080fd5b50610c1d89828a01610af6565b979a9699509760408101359660609091013595509350505050565b600060208284031215610c4a57600080fd5b5051919050565b600060208284031215610c6357600080fd5b5035919050565b600082821015610c7c57610c7c610cba565b500390565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415610cb357610cb3610cba565b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fdfea164736f6c6343000806000a",
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
	return address, tx, &TrustedBlockhashStore{TrustedBlockhashStoreCaller: TrustedBlockhashStoreCaller{contract: contract}, TrustedBlockhashStoreTransactor: TrustedBlockhashStoreTransactor{contract: contract}, TrustedBlockhashStoreFilterer: TrustedBlockhashStoreFilterer{contract: contract}}, nil
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

func (_TrustedBlockhashStore *TrustedBlockhashStoreCaller) SBlockhashes(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _TrustedBlockhashStore.contract.Call(opts, &out, "s_blockhashes", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_TrustedBlockhashStore *TrustedBlockhashStoreSession) SBlockhashes(arg0 *big.Int) ([32]byte, error) {
	return _TrustedBlockhashStore.Contract.SBlockhashes(&_TrustedBlockhashStore.CallOpts, arg0)
}

func (_TrustedBlockhashStore *TrustedBlockhashStoreCallerSession) SBlockhashes(arg0 *big.Int) ([32]byte, error) {
	return _TrustedBlockhashStore.Contract.SBlockhashes(&_TrustedBlockhashStore.CallOpts, arg0)
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

func (_TrustedBlockhashStore *TrustedBlockhashStoreTransactor) StoreTrusted(opts *bind.TransactOpts, blockNums []*big.Int, blockhashes [][32]byte, recentBlockNumber *big.Int, recentBlockhash [32]byte) (*types.Transaction, error) {
	return _TrustedBlockhashStore.contract.Transact(opts, "storeTrusted", blockNums, blockhashes, recentBlockNumber, recentBlockhash)
}

func (_TrustedBlockhashStore *TrustedBlockhashStoreSession) StoreTrusted(blockNums []*big.Int, blockhashes [][32]byte, recentBlockNumber *big.Int, recentBlockhash [32]byte) (*types.Transaction, error) {
	return _TrustedBlockhashStore.Contract.StoreTrusted(&_TrustedBlockhashStore.TransactOpts, blockNums, blockhashes, recentBlockNumber, recentBlockhash)
}

func (_TrustedBlockhashStore *TrustedBlockhashStoreTransactorSession) StoreTrusted(blockNums []*big.Int, blockhashes [][32]byte, recentBlockNumber *big.Int, recentBlockhash [32]byte) (*types.Transaction, error) {
	return _TrustedBlockhashStore.Contract.StoreTrusted(&_TrustedBlockhashStore.TransactOpts, blockNums, blockhashes, recentBlockNumber, recentBlockhash)
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

	SBlockhashes(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error)

	SWhitelist(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error)

	SWhitelistStatus(opts *bind.CallOpts, arg0 common.Address) (bool, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	SetWhitelist(opts *bind.TransactOpts, whitelist []common.Address) (*types.Transaction, error)

	Store(opts *bind.TransactOpts, n *big.Int) (*types.Transaction, error)

	StoreTrusted(opts *bind.TransactOpts, blockNums []*big.Int, blockhashes [][32]byte, recentBlockNumber *big.Int, recentBlockhash [32]byte) (*types.Transaction, error)

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
