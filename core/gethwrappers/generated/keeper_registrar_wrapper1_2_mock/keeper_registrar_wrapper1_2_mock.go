// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package keeper_registrar_wrapper1_2_mock

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

var KeeperRegistrarMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"emitAutoApproveAllowedSenderSet\",\"inputs\":[{\"name\":\"senderAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowed\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"emitConfigChanged\",\"inputs\":[{\"name\":\"autoApproveConfigType\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"autoApproveMaxAllowed\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"keeperRegistry\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"minLINKJuels\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"emitOwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"emitOwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"emitRegistrationApproved\",\"inputs\":[{\"name\":\"hash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"displayName\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"upkeepId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"emitRegistrationRejected\",\"inputs\":[{\"name\":\"hash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"emitRegistrationRequested\",\"inputs\":[{\"name\":\"hash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encryptedEmail\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"upkeepContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"gasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"adminAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"checkData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"source\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getRegistrationConfig\",\"inputs\":[],\"outputs\":[{\"name\":\"autoApproveConfigType\",\"type\":\"uint8\",\"internalType\":\"enumKeeperRegistrar1_2Mock.AutoApproveType\"},{\"name\":\"autoApproveMaxAllowed\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"approvedCount\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"keeperRegistry\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"minLINKJuels\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_approvedCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_autoApproveConfigType\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumKeeperRegistrar1_2Mock.AutoApproveType\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_autoApproveMaxAllowed\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_keeperRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"s_minLINKJuels\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setRegistrationConfig\",\"inputs\":[{\"name\":\"_autoApproveConfigType\",\"type\":\"uint8\",\"internalType\":\"enumKeeperRegistrar1_2Mock.AutoApproveType\"},{\"name\":\"_autoApproveMaxAllowed\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"_approvedCount\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"_keeperRegistry\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_minLINKJuels\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AutoApproveAllowedSenderSet\",\"inputs\":[{\"name\":\"senderAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"allowed\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ConfigChanged\",\"inputs\":[{\"name\":\"autoApproveConfigType\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"},{\"name\":\"autoApproveMaxAllowed\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"keeperRegistry\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"minLINKJuels\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RegistrationApproved\",\"inputs\":[{\"name\":\"hash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"displayName\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"upkeepId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RegistrationRejected\",\"inputs\":[{\"name\":\"hash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RegistrationRequested\",\"inputs\":[{\"name\":\"hash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"name\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"encryptedEmail\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"upkeepContract\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"gasLimit\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"adminAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"checkData\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"},{\"name\":\"source\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"uint8\"}],\"anonymous\":false}]",
	Bin: "0x608060405234801561001057600080fd5b50610ba0806100206000396000f3fe608060405234801561001057600080fd5b50600436106100ea5760003560e01c8063aee052f31161008c578063b59d75eb11610066578063b59d75eb14610259578063bb98fe561461026c578063ca40bcd314610285578063f7420bc21461029857600080fd5b8063aee052f314610220578063b019b4e814610233578063b49fd35b1461024657600080fd5b806384638bb6116100c857806384638bb61461014a578063850af0cb146101645780639e105f95146101bb578063adeab0b71461020d57600080fd5b80631701f938146100ef5780634882b5bd1461010457806355e8b24814610120575b600080fd5b6101026100fd36600461093a565b6102ab565b005b61010d60015481565b6040519081526020015b60405180910390f35b60005461013590610100900463ffffffff1681565b60405163ffffffff9091168152602001610117565b6000546101579060ff1681565b6040516101179190610a34565b6000546001546040516101179260ff811692610100820463ffffffff90811693650100000000008404909116926901000000000000000000900473ffffffffffffffffffffffffffffffffffffffff169190610a48565b6000546101e8906901000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610117565b61010261021b366004610789565b610322565b61010261022e366004610888565b610350565b61010261024136600461071a565b61038e565b6101026102543660046108d8565b6103ec565b6101026102673660046107a2565b6104e9565b6000546101359065010000000000900463ffffffff1681565b61010261029336600461074d565b610551565b6101026102a636600461071a565b6105a7565b6040805160ff8616815263ffffffff8516602082015273ffffffffffffffffffffffffffffffffffffffff8416818301526bffffffffffffffffffffffff8316606082015290517f6293a703ec7145dfa23c5cde2e627d6a02e153fc2e9c03b14d1e22cbb4a7e9cd9181900360800190a150505050565b60405181907f3663fb28ebc87645eb972c9dad8521bf665c623f287e79f1c56f1eb374b82a2290600090a250565b80837fb9a292fb7e3edd920cd2d2829a3615a640c43fd7de0a0820aa0668feb4c37d4b846040516103819190610a98565b60405180910390a3505050565b8073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b600080548691907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016600183600281111561042957610429610b35565b0217905550600080547fffffffffffffffffffffffffffffffffffffffffffffff0000000000000000ff1661010063ffffffff968716027fffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffffff1617650100000000009490951693909302939093177fffffff0000000000000000000000000000000000000000ffffffffffffffffff16690100000000000000000073ffffffffffffffffffffffffffffffffffffffff929092169190910217905560015550565b8060ff168673ffffffffffffffffffffffffffffffffffffffff168a7fc3f5df4aefec026f610a3fcb08f19476492d69d2cb78b1c2eba259a8820e6a788b8b8a8a8a8a60405161053e96959493929190610ab2565b60405180910390a4505050505050505050565b8173ffffffffffffffffffffffffffffffffffffffff167f20c6237dac83526a849285a9f79d08a483291bdd3a056a0ef9ae94ecee1ad3568260405161059b911515815260200190565b60405180910390a25050565b8073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127860405160405180910390a35050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461062957600080fd5b919050565b600082601f83011261063f57600080fd5b813567ffffffffffffffff8082111561065a5761065a610b64565b604051601f83017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f011681019082821181831017156106a0576106a0610b64565b816040528381528660208588010111156106b957600080fd5b836020870160208301376000602085830101528094505050505092915050565b803563ffffffff8116811461062957600080fd5b803560ff8116811461062957600080fd5b80356bffffffffffffffffffffffff8116811461062957600080fd5b6000806040838503121561072d57600080fd5b61073683610605565b915061074460208401610605565b90509250929050565b6000806040838503121561076057600080fd5b61076983610605565b91506020830135801515811461077e57600080fd5b809150509250929050565b60006020828403121561079b57600080fd5b5035919050565b60008060008060008060008060006101208a8c0312156107c157600080fd5b8935985060208a013567ffffffffffffffff808211156107e057600080fd5b6107ec8d838e0161062e565b995060408c013591508082111561080257600080fd5b61080e8d838e0161062e565b985061081c60608d01610605565b975061082a60808d016106d9565b965061083860a08d01610605565b955060c08c013591508082111561084e57600080fd5b5061085b8c828d0161062e565b93505061086a60e08b016106fe565b91506108796101008b016106ed565b90509295985092959850929598565b60008060006060848603121561089d57600080fd5b83359250602084013567ffffffffffffffff8111156108bb57600080fd5b6108c78682870161062e565b925050604084013590509250925092565b600080600080600060a086880312156108f057600080fd5b8535600381106108ff57600080fd5b945061090d602087016106d9565b935061091b604087016106d9565b925061092960608701610605565b949793965091946080013592915050565b6000806000806080858703121561095057600080fd5b610959856106ed565b9350610967602086016106d9565b925061097560408601610605565b9150610983606086016106fe565b905092959194509250565b6000815180845260005b818110156109b457602081850181015186830182015201610998565b818111156109c6576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b60038110610a30577f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b9052565b60208101610a4282846109f9565b92915050565b60a08101610a5682886109f9565b63ffffffff808716602084015280861660408401525073ffffffffffffffffffffffffffffffffffffffff841660608301528260808301529695505050505050565b602081526000610aab602083018461098e565b9392505050565b60c081526000610ac560c083018961098e565b8281036020840152610ad7818961098e565b905063ffffffff8716604084015273ffffffffffffffffffffffffffffffffffffffff861660608401528281036080840152610b13818661098e565b9150506bffffffffffffffffffffffff831660a0830152979650505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var KeeperRegistrarABI = KeeperRegistrarMetaData.ABI

var KeeperRegistrarBin = KeeperRegistrarMetaData.Bin

func DeployKeeperRegistrar(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *KeeperRegistrar, error) {
	parsed, err := KeeperRegistrarMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KeeperRegistrarBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KeeperRegistrar{address: address, abi: *parsed, KeeperRegistrarCaller: KeeperRegistrarCaller{contract: contract}, KeeperRegistrarTransactor: KeeperRegistrarTransactor{contract: contract}, KeeperRegistrarFilterer: KeeperRegistrarFilterer{contract: contract}}, nil
}

type KeeperRegistrar struct {
	address common.Address
	abi     abi.ABI
	KeeperRegistrarCaller
	KeeperRegistrarTransactor
	KeeperRegistrarFilterer
}

type KeeperRegistrarCaller struct {
	contract *bind.BoundContract
}

type KeeperRegistrarTransactor struct {
	contract *bind.BoundContract
}

type KeeperRegistrarFilterer struct {
	contract *bind.BoundContract
}

type KeeperRegistrarSession struct {
	Contract     *KeeperRegistrar
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type KeeperRegistrarCallerSession struct {
	Contract *KeeperRegistrarCaller
	CallOpts bind.CallOpts
}

type KeeperRegistrarTransactorSession struct {
	Contract     *KeeperRegistrarTransactor
	TransactOpts bind.TransactOpts
}

type KeeperRegistrarRaw struct {
	Contract *KeeperRegistrar
}

type KeeperRegistrarCallerRaw struct {
	Contract *KeeperRegistrarCaller
}

type KeeperRegistrarTransactorRaw struct {
	Contract *KeeperRegistrarTransactor
}

func NewKeeperRegistrar(address common.Address, backend bind.ContractBackend) (*KeeperRegistrar, error) {
	abi, err := abi.JSON(strings.NewReader(KeeperRegistrarABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindKeeperRegistrar(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrar{address: address, abi: abi, KeeperRegistrarCaller: KeeperRegistrarCaller{contract: contract}, KeeperRegistrarTransactor: KeeperRegistrarTransactor{contract: contract}, KeeperRegistrarFilterer: KeeperRegistrarFilterer{contract: contract}}, nil
}

func NewKeeperRegistrarCaller(address common.Address, caller bind.ContractCaller) (*KeeperRegistrarCaller, error) {
	contract, err := bindKeeperRegistrar(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrarCaller{contract: contract}, nil
}

func NewKeeperRegistrarTransactor(address common.Address, transactor bind.ContractTransactor) (*KeeperRegistrarTransactor, error) {
	contract, err := bindKeeperRegistrar(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrarTransactor{contract: contract}, nil
}

func NewKeeperRegistrarFilterer(address common.Address, filterer bind.ContractFilterer) (*KeeperRegistrarFilterer, error) {
	contract, err := bindKeeperRegistrar(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrarFilterer{contract: contract}, nil
}

func bindKeeperRegistrar(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := KeeperRegistrarMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_KeeperRegistrar *KeeperRegistrarRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistrar.Contract.KeeperRegistrarCaller.contract.Call(opts, result, method, params...)
}

func (_KeeperRegistrar *KeeperRegistrarRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.KeeperRegistrarTransactor.contract.Transfer(opts)
}

func (_KeeperRegistrar *KeeperRegistrarRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.KeeperRegistrarTransactor.contract.Transact(opts, method, params...)
}

func (_KeeperRegistrar *KeeperRegistrarCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistrar.Contract.contract.Call(opts, result, method, params...)
}

func (_KeeperRegistrar *KeeperRegistrarTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.contract.Transfer(opts)
}

func (_KeeperRegistrar *KeeperRegistrarTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.contract.Transact(opts, method, params...)
}

func (_KeeperRegistrar *KeeperRegistrarCaller) GetRegistrationConfig(opts *bind.CallOpts) (GetRegistrationConfig,

	error) {
	var out []interface{}
	err := _KeeperRegistrar.contract.Call(opts, &out, "getRegistrationConfig")

	outstruct := new(GetRegistrationConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.AutoApproveConfigType = *abi.ConvertType(out[0], new(uint8)).(*uint8)
	outstruct.AutoApproveMaxAllowed = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ApprovedCount = *abi.ConvertType(out[2], new(uint32)).(*uint32)
	outstruct.KeeperRegistry = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)
	outstruct.MinLINKJuels = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_KeeperRegistrar *KeeperRegistrarSession) GetRegistrationConfig() (GetRegistrationConfig,

	error) {
	return _KeeperRegistrar.Contract.GetRegistrationConfig(&_KeeperRegistrar.CallOpts)
}

func (_KeeperRegistrar *KeeperRegistrarCallerSession) GetRegistrationConfig() (GetRegistrationConfig,

	error) {
	return _KeeperRegistrar.Contract.GetRegistrationConfig(&_KeeperRegistrar.CallOpts)
}

func (_KeeperRegistrar *KeeperRegistrarCaller) SApprovedCount(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _KeeperRegistrar.contract.Call(opts, &out, "s_approvedCount")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_KeeperRegistrar *KeeperRegistrarSession) SApprovedCount() (uint32, error) {
	return _KeeperRegistrar.Contract.SApprovedCount(&_KeeperRegistrar.CallOpts)
}

func (_KeeperRegistrar *KeeperRegistrarCallerSession) SApprovedCount() (uint32, error) {
	return _KeeperRegistrar.Contract.SApprovedCount(&_KeeperRegistrar.CallOpts)
}

func (_KeeperRegistrar *KeeperRegistrarCaller) SAutoApproveConfigType(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistrar.contract.Call(opts, &out, "s_autoApproveConfigType")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_KeeperRegistrar *KeeperRegistrarSession) SAutoApproveConfigType() (uint8, error) {
	return _KeeperRegistrar.Contract.SAutoApproveConfigType(&_KeeperRegistrar.CallOpts)
}

func (_KeeperRegistrar *KeeperRegistrarCallerSession) SAutoApproveConfigType() (uint8, error) {
	return _KeeperRegistrar.Contract.SAutoApproveConfigType(&_KeeperRegistrar.CallOpts)
}

func (_KeeperRegistrar *KeeperRegistrarCaller) SAutoApproveMaxAllowed(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _KeeperRegistrar.contract.Call(opts, &out, "s_autoApproveMaxAllowed")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_KeeperRegistrar *KeeperRegistrarSession) SAutoApproveMaxAllowed() (uint32, error) {
	return _KeeperRegistrar.Contract.SAutoApproveMaxAllowed(&_KeeperRegistrar.CallOpts)
}

func (_KeeperRegistrar *KeeperRegistrarCallerSession) SAutoApproveMaxAllowed() (uint32, error) {
	return _KeeperRegistrar.Contract.SAutoApproveMaxAllowed(&_KeeperRegistrar.CallOpts)
}

func (_KeeperRegistrar *KeeperRegistrarCaller) SKeeperRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistrar.contract.Call(opts, &out, "s_keeperRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistrar *KeeperRegistrarSession) SKeeperRegistry() (common.Address, error) {
	return _KeeperRegistrar.Contract.SKeeperRegistry(&_KeeperRegistrar.CallOpts)
}

func (_KeeperRegistrar *KeeperRegistrarCallerSession) SKeeperRegistry() (common.Address, error) {
	return _KeeperRegistrar.Contract.SKeeperRegistry(&_KeeperRegistrar.CallOpts)
}

func (_KeeperRegistrar *KeeperRegistrarCaller) SMinLINKJuels(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistrar.contract.Call(opts, &out, "s_minLINKJuels")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_KeeperRegistrar *KeeperRegistrarSession) SMinLINKJuels() (*big.Int, error) {
	return _KeeperRegistrar.Contract.SMinLINKJuels(&_KeeperRegistrar.CallOpts)
}

func (_KeeperRegistrar *KeeperRegistrarCallerSession) SMinLINKJuels() (*big.Int, error) {
	return _KeeperRegistrar.Contract.SMinLINKJuels(&_KeeperRegistrar.CallOpts)
}

func (_KeeperRegistrar *KeeperRegistrarTransactor) EmitAutoApproveAllowedSenderSet(opts *bind.TransactOpts, senderAddress common.Address, allowed bool) (*types.Transaction, error) {
	return _KeeperRegistrar.contract.Transact(opts, "emitAutoApproveAllowedSenderSet", senderAddress, allowed)
}

func (_KeeperRegistrar *KeeperRegistrarSession) EmitAutoApproveAllowedSenderSet(senderAddress common.Address, allowed bool) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.EmitAutoApproveAllowedSenderSet(&_KeeperRegistrar.TransactOpts, senderAddress, allowed)
}

func (_KeeperRegistrar *KeeperRegistrarTransactorSession) EmitAutoApproveAllowedSenderSet(senderAddress common.Address, allowed bool) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.EmitAutoApproveAllowedSenderSet(&_KeeperRegistrar.TransactOpts, senderAddress, allowed)
}

func (_KeeperRegistrar *KeeperRegistrarTransactor) EmitConfigChanged(opts *bind.TransactOpts, autoApproveConfigType uint8, autoApproveMaxAllowed uint32, keeperRegistry common.Address, minLINKJuels *big.Int) (*types.Transaction, error) {
	return _KeeperRegistrar.contract.Transact(opts, "emitConfigChanged", autoApproveConfigType, autoApproveMaxAllowed, keeperRegistry, minLINKJuels)
}

func (_KeeperRegistrar *KeeperRegistrarSession) EmitConfigChanged(autoApproveConfigType uint8, autoApproveMaxAllowed uint32, keeperRegistry common.Address, minLINKJuels *big.Int) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.EmitConfigChanged(&_KeeperRegistrar.TransactOpts, autoApproveConfigType, autoApproveMaxAllowed, keeperRegistry, minLINKJuels)
}

func (_KeeperRegistrar *KeeperRegistrarTransactorSession) EmitConfigChanged(autoApproveConfigType uint8, autoApproveMaxAllowed uint32, keeperRegistry common.Address, minLINKJuels *big.Int) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.EmitConfigChanged(&_KeeperRegistrar.TransactOpts, autoApproveConfigType, autoApproveMaxAllowed, keeperRegistry, minLINKJuels)
}

func (_KeeperRegistrar *KeeperRegistrarTransactor) EmitOwnershipTransferRequested(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistrar.contract.Transact(opts, "emitOwnershipTransferRequested", from, to)
}

func (_KeeperRegistrar *KeeperRegistrarSession) EmitOwnershipTransferRequested(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.EmitOwnershipTransferRequested(&_KeeperRegistrar.TransactOpts, from, to)
}

func (_KeeperRegistrar *KeeperRegistrarTransactorSession) EmitOwnershipTransferRequested(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.EmitOwnershipTransferRequested(&_KeeperRegistrar.TransactOpts, from, to)
}

func (_KeeperRegistrar *KeeperRegistrarTransactor) EmitOwnershipTransferred(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistrar.contract.Transact(opts, "emitOwnershipTransferred", from, to)
}

func (_KeeperRegistrar *KeeperRegistrarSession) EmitOwnershipTransferred(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.EmitOwnershipTransferred(&_KeeperRegistrar.TransactOpts, from, to)
}

func (_KeeperRegistrar *KeeperRegistrarTransactorSession) EmitOwnershipTransferred(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.EmitOwnershipTransferred(&_KeeperRegistrar.TransactOpts, from, to)
}

func (_KeeperRegistrar *KeeperRegistrarTransactor) EmitRegistrationApproved(opts *bind.TransactOpts, hash [32]byte, displayName string, upkeepId *big.Int) (*types.Transaction, error) {
	return _KeeperRegistrar.contract.Transact(opts, "emitRegistrationApproved", hash, displayName, upkeepId)
}

func (_KeeperRegistrar *KeeperRegistrarSession) EmitRegistrationApproved(hash [32]byte, displayName string, upkeepId *big.Int) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.EmitRegistrationApproved(&_KeeperRegistrar.TransactOpts, hash, displayName, upkeepId)
}

func (_KeeperRegistrar *KeeperRegistrarTransactorSession) EmitRegistrationApproved(hash [32]byte, displayName string, upkeepId *big.Int) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.EmitRegistrationApproved(&_KeeperRegistrar.TransactOpts, hash, displayName, upkeepId)
}

func (_KeeperRegistrar *KeeperRegistrarTransactor) EmitRegistrationRejected(opts *bind.TransactOpts, hash [32]byte) (*types.Transaction, error) {
	return _KeeperRegistrar.contract.Transact(opts, "emitRegistrationRejected", hash)
}

func (_KeeperRegistrar *KeeperRegistrarSession) EmitRegistrationRejected(hash [32]byte) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.EmitRegistrationRejected(&_KeeperRegistrar.TransactOpts, hash)
}

func (_KeeperRegistrar *KeeperRegistrarTransactorSession) EmitRegistrationRejected(hash [32]byte) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.EmitRegistrationRejected(&_KeeperRegistrar.TransactOpts, hash)
}

func (_KeeperRegistrar *KeeperRegistrarTransactor) EmitRegistrationRequested(opts *bind.TransactOpts, hash [32]byte, name string, encryptedEmail []byte, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, amount *big.Int, source uint8) (*types.Transaction, error) {
	return _KeeperRegistrar.contract.Transact(opts, "emitRegistrationRequested", hash, name, encryptedEmail, upkeepContract, gasLimit, adminAddress, checkData, amount, source)
}

func (_KeeperRegistrar *KeeperRegistrarSession) EmitRegistrationRequested(hash [32]byte, name string, encryptedEmail []byte, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, amount *big.Int, source uint8) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.EmitRegistrationRequested(&_KeeperRegistrar.TransactOpts, hash, name, encryptedEmail, upkeepContract, gasLimit, adminAddress, checkData, amount, source)
}

func (_KeeperRegistrar *KeeperRegistrarTransactorSession) EmitRegistrationRequested(hash [32]byte, name string, encryptedEmail []byte, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, amount *big.Int, source uint8) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.EmitRegistrationRequested(&_KeeperRegistrar.TransactOpts, hash, name, encryptedEmail, upkeepContract, gasLimit, adminAddress, checkData, amount, source)
}

func (_KeeperRegistrar *KeeperRegistrarTransactor) SetRegistrationConfig(opts *bind.TransactOpts, _autoApproveConfigType uint8, _autoApproveMaxAllowed uint32, _approvedCount uint32, _keeperRegistry common.Address, _minLINKJuels *big.Int) (*types.Transaction, error) {
	return _KeeperRegistrar.contract.Transact(opts, "setRegistrationConfig", _autoApproveConfigType, _autoApproveMaxAllowed, _approvedCount, _keeperRegistry, _minLINKJuels)
}

func (_KeeperRegistrar *KeeperRegistrarSession) SetRegistrationConfig(_autoApproveConfigType uint8, _autoApproveMaxAllowed uint32, _approvedCount uint32, _keeperRegistry common.Address, _minLINKJuels *big.Int) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.SetRegistrationConfig(&_KeeperRegistrar.TransactOpts, _autoApproveConfigType, _autoApproveMaxAllowed, _approvedCount, _keeperRegistry, _minLINKJuels)
}

func (_KeeperRegistrar *KeeperRegistrarTransactorSession) SetRegistrationConfig(_autoApproveConfigType uint8, _autoApproveMaxAllowed uint32, _approvedCount uint32, _keeperRegistry common.Address, _minLINKJuels *big.Int) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.SetRegistrationConfig(&_KeeperRegistrar.TransactOpts, _autoApproveConfigType, _autoApproveMaxAllowed, _approvedCount, _keeperRegistry, _minLINKJuels)
}

type KeeperRegistrarAutoApproveAllowedSenderSetIterator struct {
	Event *KeeperRegistrarAutoApproveAllowedSenderSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistrarAutoApproveAllowedSenderSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistrarAutoApproveAllowedSenderSet)
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
		it.Event = new(KeeperRegistrarAutoApproveAllowedSenderSet)
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

func (it *KeeperRegistrarAutoApproveAllowedSenderSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistrarAutoApproveAllowedSenderSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistrarAutoApproveAllowedSenderSet struct {
	SenderAddress common.Address
	Allowed       bool
	Raw           types.Log
}

func (_KeeperRegistrar *KeeperRegistrarFilterer) FilterAutoApproveAllowedSenderSet(opts *bind.FilterOpts, senderAddress []common.Address) (*KeeperRegistrarAutoApproveAllowedSenderSetIterator, error) {

	var senderAddressRule []interface{}
	for _, senderAddressItem := range senderAddress {
		senderAddressRule = append(senderAddressRule, senderAddressItem)
	}

	logs, sub, err := _KeeperRegistrar.contract.FilterLogs(opts, "AutoApproveAllowedSenderSet", senderAddressRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrarAutoApproveAllowedSenderSetIterator{contract: _KeeperRegistrar.contract, event: "AutoApproveAllowedSenderSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistrar *KeeperRegistrarFilterer) WatchAutoApproveAllowedSenderSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarAutoApproveAllowedSenderSet, senderAddress []common.Address) (event.Subscription, error) {

	var senderAddressRule []interface{}
	for _, senderAddressItem := range senderAddress {
		senderAddressRule = append(senderAddressRule, senderAddressItem)
	}

	logs, sub, err := _KeeperRegistrar.contract.WatchLogs(opts, "AutoApproveAllowedSenderSet", senderAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistrarAutoApproveAllowedSenderSet)
				if err := _KeeperRegistrar.contract.UnpackLog(event, "AutoApproveAllowedSenderSet", log); err != nil {
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

func (_KeeperRegistrar *KeeperRegistrarFilterer) ParseAutoApproveAllowedSenderSet(log types.Log) (*KeeperRegistrarAutoApproveAllowedSenderSet, error) {
	event := new(KeeperRegistrarAutoApproveAllowedSenderSet)
	if err := _KeeperRegistrar.contract.UnpackLog(event, "AutoApproveAllowedSenderSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistrarConfigChangedIterator struct {
	Event *KeeperRegistrarConfigChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistrarConfigChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistrarConfigChanged)
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
		it.Event = new(KeeperRegistrarConfigChanged)
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

func (it *KeeperRegistrarConfigChangedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistrarConfigChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistrarConfigChanged struct {
	AutoApproveConfigType uint8
	AutoApproveMaxAllowed uint32
	KeeperRegistry        common.Address
	MinLINKJuels          *big.Int
	Raw                   types.Log
}

func (_KeeperRegistrar *KeeperRegistrarFilterer) FilterConfigChanged(opts *bind.FilterOpts) (*KeeperRegistrarConfigChangedIterator, error) {

	logs, sub, err := _KeeperRegistrar.contract.FilterLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrarConfigChangedIterator{contract: _KeeperRegistrar.contract, event: "ConfigChanged", logs: logs, sub: sub}, nil
}

func (_KeeperRegistrar *KeeperRegistrarFilterer) WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarConfigChanged) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistrar.contract.WatchLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistrarConfigChanged)
				if err := _KeeperRegistrar.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
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

func (_KeeperRegistrar *KeeperRegistrarFilterer) ParseConfigChanged(log types.Log) (*KeeperRegistrarConfigChanged, error) {
	event := new(KeeperRegistrarConfigChanged)
	if err := _KeeperRegistrar.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistrarOwnershipTransferRequestedIterator struct {
	Event *KeeperRegistrarOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistrarOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistrarOwnershipTransferRequested)
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
		it.Event = new(KeeperRegistrarOwnershipTransferRequested)
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

func (it *KeeperRegistrarOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistrarOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistrarOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistrar *KeeperRegistrarFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistrarOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistrar.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrarOwnershipTransferRequestedIterator{contract: _KeeperRegistrar.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeeperRegistrar *KeeperRegistrarFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistrar.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistrarOwnershipTransferRequested)
				if err := _KeeperRegistrar.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_KeeperRegistrar *KeeperRegistrarFilterer) ParseOwnershipTransferRequested(log types.Log) (*KeeperRegistrarOwnershipTransferRequested, error) {
	event := new(KeeperRegistrarOwnershipTransferRequested)
	if err := _KeeperRegistrar.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistrarOwnershipTransferredIterator struct {
	Event *KeeperRegistrarOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistrarOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistrarOwnershipTransferred)
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
		it.Event = new(KeeperRegistrarOwnershipTransferred)
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

func (it *KeeperRegistrarOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistrarOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistrarOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistrar *KeeperRegistrarFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistrarOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistrar.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrarOwnershipTransferredIterator{contract: _KeeperRegistrar.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_KeeperRegistrar *KeeperRegistrarFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistrar.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistrarOwnershipTransferred)
				if err := _KeeperRegistrar.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_KeeperRegistrar *KeeperRegistrarFilterer) ParseOwnershipTransferred(log types.Log) (*KeeperRegistrarOwnershipTransferred, error) {
	event := new(KeeperRegistrarOwnershipTransferred)
	if err := _KeeperRegistrar.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistrarRegistrationApprovedIterator struct {
	Event *KeeperRegistrarRegistrationApproved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistrarRegistrationApprovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistrarRegistrationApproved)
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
		it.Event = new(KeeperRegistrarRegistrationApproved)
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

func (it *KeeperRegistrarRegistrationApprovedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistrarRegistrationApprovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistrarRegistrationApproved struct {
	Hash        [32]byte
	DisplayName string
	UpkeepId    *big.Int
	Raw         types.Log
}

func (_KeeperRegistrar *KeeperRegistrarFilterer) FilterRegistrationApproved(opts *bind.FilterOpts, hash [][32]byte, upkeepId []*big.Int) (*KeeperRegistrarRegistrationApprovedIterator, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}

	var upkeepIdRule []interface{}
	for _, upkeepIdItem := range upkeepId {
		upkeepIdRule = append(upkeepIdRule, upkeepIdItem)
	}

	logs, sub, err := _KeeperRegistrar.contract.FilterLogs(opts, "RegistrationApproved", hashRule, upkeepIdRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrarRegistrationApprovedIterator{contract: _KeeperRegistrar.contract, event: "RegistrationApproved", logs: logs, sub: sub}, nil
}

func (_KeeperRegistrar *KeeperRegistrarFilterer) WatchRegistrationApproved(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarRegistrationApproved, hash [][32]byte, upkeepId []*big.Int) (event.Subscription, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}

	var upkeepIdRule []interface{}
	for _, upkeepIdItem := range upkeepId {
		upkeepIdRule = append(upkeepIdRule, upkeepIdItem)
	}

	logs, sub, err := _KeeperRegistrar.contract.WatchLogs(opts, "RegistrationApproved", hashRule, upkeepIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistrarRegistrationApproved)
				if err := _KeeperRegistrar.contract.UnpackLog(event, "RegistrationApproved", log); err != nil {
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

func (_KeeperRegistrar *KeeperRegistrarFilterer) ParseRegistrationApproved(log types.Log) (*KeeperRegistrarRegistrationApproved, error) {
	event := new(KeeperRegistrarRegistrationApproved)
	if err := _KeeperRegistrar.contract.UnpackLog(event, "RegistrationApproved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistrarRegistrationRejectedIterator struct {
	Event *KeeperRegistrarRegistrationRejected

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistrarRegistrationRejectedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistrarRegistrationRejected)
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
		it.Event = new(KeeperRegistrarRegistrationRejected)
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

func (it *KeeperRegistrarRegistrationRejectedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistrarRegistrationRejectedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistrarRegistrationRejected struct {
	Hash [32]byte
	Raw  types.Log
}

func (_KeeperRegistrar *KeeperRegistrarFilterer) FilterRegistrationRejected(opts *bind.FilterOpts, hash [][32]byte) (*KeeperRegistrarRegistrationRejectedIterator, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}

	logs, sub, err := _KeeperRegistrar.contract.FilterLogs(opts, "RegistrationRejected", hashRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrarRegistrationRejectedIterator{contract: _KeeperRegistrar.contract, event: "RegistrationRejected", logs: logs, sub: sub}, nil
}

func (_KeeperRegistrar *KeeperRegistrarFilterer) WatchRegistrationRejected(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarRegistrationRejected, hash [][32]byte) (event.Subscription, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}

	logs, sub, err := _KeeperRegistrar.contract.WatchLogs(opts, "RegistrationRejected", hashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistrarRegistrationRejected)
				if err := _KeeperRegistrar.contract.UnpackLog(event, "RegistrationRejected", log); err != nil {
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

func (_KeeperRegistrar *KeeperRegistrarFilterer) ParseRegistrationRejected(log types.Log) (*KeeperRegistrarRegistrationRejected, error) {
	event := new(KeeperRegistrarRegistrationRejected)
	if err := _KeeperRegistrar.contract.UnpackLog(event, "RegistrationRejected", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistrarRegistrationRequestedIterator struct {
	Event *KeeperRegistrarRegistrationRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistrarRegistrationRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistrarRegistrationRequested)
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
		it.Event = new(KeeperRegistrarRegistrationRequested)
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

func (it *KeeperRegistrarRegistrationRequestedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistrarRegistrationRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistrarRegistrationRequested struct {
	Hash           [32]byte
	Name           string
	EncryptedEmail []byte
	UpkeepContract common.Address
	GasLimit       uint32
	AdminAddress   common.Address
	CheckData      []byte
	Amount         *big.Int
	Source         uint8
	Raw            types.Log
}

func (_KeeperRegistrar *KeeperRegistrarFilterer) FilterRegistrationRequested(opts *bind.FilterOpts, hash [][32]byte, upkeepContract []common.Address, source []uint8) (*KeeperRegistrarRegistrationRequestedIterator, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}

	var upkeepContractRule []interface{}
	for _, upkeepContractItem := range upkeepContract {
		upkeepContractRule = append(upkeepContractRule, upkeepContractItem)
	}

	var sourceRule []interface{}
	for _, sourceItem := range source {
		sourceRule = append(sourceRule, sourceItem)
	}

	logs, sub, err := _KeeperRegistrar.contract.FilterLogs(opts, "RegistrationRequested", hashRule, upkeepContractRule, sourceRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrarRegistrationRequestedIterator{contract: _KeeperRegistrar.contract, event: "RegistrationRequested", logs: logs, sub: sub}, nil
}

func (_KeeperRegistrar *KeeperRegistrarFilterer) WatchRegistrationRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarRegistrationRequested, hash [][32]byte, upkeepContract []common.Address, source []uint8) (event.Subscription, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}

	var upkeepContractRule []interface{}
	for _, upkeepContractItem := range upkeepContract {
		upkeepContractRule = append(upkeepContractRule, upkeepContractItem)
	}

	var sourceRule []interface{}
	for _, sourceItem := range source {
		sourceRule = append(sourceRule, sourceItem)
	}

	logs, sub, err := _KeeperRegistrar.contract.WatchLogs(opts, "RegistrationRequested", hashRule, upkeepContractRule, sourceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistrarRegistrationRequested)
				if err := _KeeperRegistrar.contract.UnpackLog(event, "RegistrationRequested", log); err != nil {
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

func (_KeeperRegistrar *KeeperRegistrarFilterer) ParseRegistrationRequested(log types.Log) (*KeeperRegistrarRegistrationRequested, error) {
	event := new(KeeperRegistrarRegistrationRequested)
	if err := _KeeperRegistrar.contract.UnpackLog(event, "RegistrationRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetRegistrationConfig struct {
	AutoApproveConfigType uint8
	AutoApproveMaxAllowed uint32
	ApprovedCount         uint32
	KeeperRegistry        common.Address
	MinLINKJuels          *big.Int
}

func (_KeeperRegistrar *KeeperRegistrar) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _KeeperRegistrar.abi.Events["AutoApproveAllowedSenderSet"].ID:
		return _KeeperRegistrar.ParseAutoApproveAllowedSenderSet(log)
	case _KeeperRegistrar.abi.Events["ConfigChanged"].ID:
		return _KeeperRegistrar.ParseConfigChanged(log)
	case _KeeperRegistrar.abi.Events["OwnershipTransferRequested"].ID:
		return _KeeperRegistrar.ParseOwnershipTransferRequested(log)
	case _KeeperRegistrar.abi.Events["OwnershipTransferred"].ID:
		return _KeeperRegistrar.ParseOwnershipTransferred(log)
	case _KeeperRegistrar.abi.Events["RegistrationApproved"].ID:
		return _KeeperRegistrar.ParseRegistrationApproved(log)
	case _KeeperRegistrar.abi.Events["RegistrationRejected"].ID:
		return _KeeperRegistrar.ParseRegistrationRejected(log)
	case _KeeperRegistrar.abi.Events["RegistrationRequested"].ID:
		return _KeeperRegistrar.ParseRegistrationRequested(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (KeeperRegistrarAutoApproveAllowedSenderSet) Topic() common.Hash {
	return common.HexToHash("0x20c6237dac83526a849285a9f79d08a483291bdd3a056a0ef9ae94ecee1ad356")
}

func (KeeperRegistrarConfigChanged) Topic() common.Hash {
	return common.HexToHash("0x6293a703ec7145dfa23c5cde2e627d6a02e153fc2e9c03b14d1e22cbb4a7e9cd")
}

func (KeeperRegistrarOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (KeeperRegistrarOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (KeeperRegistrarRegistrationApproved) Topic() common.Hash {
	return common.HexToHash("0xb9a292fb7e3edd920cd2d2829a3615a640c43fd7de0a0820aa0668feb4c37d4b")
}

func (KeeperRegistrarRegistrationRejected) Topic() common.Hash {
	return common.HexToHash("0x3663fb28ebc87645eb972c9dad8521bf665c623f287e79f1c56f1eb374b82a22")
}

func (KeeperRegistrarRegistrationRequested) Topic() common.Hash {
	return common.HexToHash("0xc3f5df4aefec026f610a3fcb08f19476492d69d2cb78b1c2eba259a8820e6a78")
}

func (_KeeperRegistrar *KeeperRegistrar) Address() common.Address {
	return _KeeperRegistrar.address
}

type KeeperRegistrarInterface interface {
	GetRegistrationConfig(opts *bind.CallOpts) (GetRegistrationConfig,

		error)

	SApprovedCount(opts *bind.CallOpts) (uint32, error)

	SAutoApproveConfigType(opts *bind.CallOpts) (uint8, error)

	SAutoApproveMaxAllowed(opts *bind.CallOpts) (uint32, error)

	SKeeperRegistry(opts *bind.CallOpts) (common.Address, error)

	SMinLINKJuels(opts *bind.CallOpts) (*big.Int, error)

	EmitAutoApproveAllowedSenderSet(opts *bind.TransactOpts, senderAddress common.Address, allowed bool) (*types.Transaction, error)

	EmitConfigChanged(opts *bind.TransactOpts, autoApproveConfigType uint8, autoApproveMaxAllowed uint32, keeperRegistry common.Address, minLINKJuels *big.Int) (*types.Transaction, error)

	EmitOwnershipTransferRequested(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error)

	EmitOwnershipTransferred(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error)

	EmitRegistrationApproved(opts *bind.TransactOpts, hash [32]byte, displayName string, upkeepId *big.Int) (*types.Transaction, error)

	EmitRegistrationRejected(opts *bind.TransactOpts, hash [32]byte) (*types.Transaction, error)

	EmitRegistrationRequested(opts *bind.TransactOpts, hash [32]byte, name string, encryptedEmail []byte, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, amount *big.Int, source uint8) (*types.Transaction, error)

	SetRegistrationConfig(opts *bind.TransactOpts, _autoApproveConfigType uint8, _autoApproveMaxAllowed uint32, _approvedCount uint32, _keeperRegistry common.Address, _minLINKJuels *big.Int) (*types.Transaction, error)

	FilterAutoApproveAllowedSenderSet(opts *bind.FilterOpts, senderAddress []common.Address) (*KeeperRegistrarAutoApproveAllowedSenderSetIterator, error)

	WatchAutoApproveAllowedSenderSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarAutoApproveAllowedSenderSet, senderAddress []common.Address) (event.Subscription, error)

	ParseAutoApproveAllowedSenderSet(log types.Log) (*KeeperRegistrarAutoApproveAllowedSenderSet, error)

	FilterConfigChanged(opts *bind.FilterOpts) (*KeeperRegistrarConfigChangedIterator, error)

	WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarConfigChanged) (event.Subscription, error)

	ParseConfigChanged(log types.Log) (*KeeperRegistrarConfigChanged, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistrarOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*KeeperRegistrarOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistrarOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*KeeperRegistrarOwnershipTransferred, error)

	FilterRegistrationApproved(opts *bind.FilterOpts, hash [][32]byte, upkeepId []*big.Int) (*KeeperRegistrarRegistrationApprovedIterator, error)

	WatchRegistrationApproved(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarRegistrationApproved, hash [][32]byte, upkeepId []*big.Int) (event.Subscription, error)

	ParseRegistrationApproved(log types.Log) (*KeeperRegistrarRegistrationApproved, error)

	FilterRegistrationRejected(opts *bind.FilterOpts, hash [][32]byte) (*KeeperRegistrarRegistrationRejectedIterator, error)

	WatchRegistrationRejected(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarRegistrationRejected, hash [][32]byte) (event.Subscription, error)

	ParseRegistrationRejected(log types.Log) (*KeeperRegistrarRegistrationRejected, error)

	FilterRegistrationRequested(opts *bind.FilterOpts, hash [][32]byte, upkeepContract []common.Address, source []uint8) (*KeeperRegistrarRegistrationRequestedIterator, error)

	WatchRegistrationRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarRegistrationRequested, hash [][32]byte, upkeepContract []common.Address, source []uint8) (event.Subscription, error)

	ParseRegistrationRequested(log types.Log) (*KeeperRegistrarRegistrationRequested, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
