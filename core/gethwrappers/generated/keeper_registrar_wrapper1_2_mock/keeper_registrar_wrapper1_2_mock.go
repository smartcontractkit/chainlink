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

var KeeperRegistrarMockMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"AutoApproveAllowedSenderSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"autoApproveConfigType\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"autoApproveMaxAllowed\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"keeperRegistry\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"minLINKJuels\",\"type\":\"uint96\"}],\"name\":\"ConfigChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"displayName\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"RegistrationApproved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"RegistrationRejected\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"encryptedEmail\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"upkeepContract\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"adminAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"indexed\":true,\"internalType\":\"uint8\",\"name\":\"source\",\"type\":\"uint8\"}],\"name\":\"RegistrationRequested\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"emitAutoApproveAllowedSenderSet\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"autoApproveConfigType\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"autoApproveMaxAllowed\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"keeperRegistry\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"minLINKJuels\",\"type\":\"uint96\"}],\"name\":\"emitConfigChanged\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"emitOwnershipTransferRequested\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"emitOwnershipTransferred\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"displayName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"emitRegistrationApproved\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"emitRegistrationRejected\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"encryptedEmail\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"upkeepContract\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"adminAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"internalType\":\"uint8\",\"name\":\"source\",\"type\":\"uint8\"}],\"name\":\"emitRegistrationRequested\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRegistrationConfig\",\"outputs\":[{\"internalType\":\"enumKeeperRegistrar1_2Mock.AutoApproveType\",\"name\":\"autoApproveConfigType\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"autoApproveMaxAllowed\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"approvedCount\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"keeperRegistry\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"minLINKJuels\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_approvedCount\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_autoApproveConfigType\",\"outputs\":[{\"internalType\":\"enumKeeperRegistrar1_2Mock.AutoApproveType\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_autoApproveMaxAllowed\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_keeperRegistry\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_minLINKJuels\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumKeeperRegistrar1_2Mock.AutoApproveType\",\"name\":\"_autoApproveConfigType\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"_autoApproveMaxAllowed\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_approvedCount\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"_keeperRegistry\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_minLINKJuels\",\"type\":\"uint256\"}],\"name\":\"setRegistrationConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610ba0806100206000396000f3fe608060405234801561001057600080fd5b50600436106100ea5760003560e01c8063aee052f31161008c578063b59d75eb11610066578063b59d75eb14610259578063bb98fe561461026c578063ca40bcd314610285578063f7420bc21461029857600080fd5b8063aee052f314610220578063b019b4e814610233578063b49fd35b1461024657600080fd5b806384638bb6116100c857806384638bb61461014a578063850af0cb146101645780639e105f95146101bb578063adeab0b71461020d57600080fd5b80631701f938146100ef5780634882b5bd1461010457806355e8b24814610120575b600080fd5b6101026100fd36600461093a565b6102ab565b005b61010d60015481565b6040519081526020015b60405180910390f35b60005461013590610100900463ffffffff1681565b60405163ffffffff9091168152602001610117565b6000546101579060ff1681565b6040516101179190610a34565b6000546001546040516101179260ff811692610100820463ffffffff90811693650100000000008404909116926901000000000000000000900473ffffffffffffffffffffffffffffffffffffffff169190610a48565b6000546101e8906901000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610117565b61010261021b366004610789565b610322565b61010261022e366004610888565b610350565b61010261024136600461071a565b61038e565b6101026102543660046108d8565b6103ec565b6101026102673660046107a2565b6104e9565b6000546101359065010000000000900463ffffffff1681565b61010261029336600461074d565b610551565b6101026102a636600461071a565b6105a7565b6040805160ff8616815263ffffffff8516602082015273ffffffffffffffffffffffffffffffffffffffff8416818301526bffffffffffffffffffffffff8316606082015290517f6293a703ec7145dfa23c5cde2e627d6a02e153fc2e9c03b14d1e22cbb4a7e9cd9181900360800190a150505050565b60405181907f3663fb28ebc87645eb972c9dad8521bf665c623f287e79f1c56f1eb374b82a2290600090a250565b80837fb9a292fb7e3edd920cd2d2829a3615a640c43fd7de0a0820aa0668feb4c37d4b846040516103819190610a98565b60405180910390a3505050565b8073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b600080548691907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016600183600281111561042957610429610b35565b0217905550600080547fffffffffffffffffffffffffffffffffffffffffffffff0000000000000000ff1661010063ffffffff968716027fffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffffff1617650100000000009490951693909302939093177fffffff0000000000000000000000000000000000000000ffffffffffffffffff16690100000000000000000073ffffffffffffffffffffffffffffffffffffffff929092169190910217905560015550565b8060ff168673ffffffffffffffffffffffffffffffffffffffff168a7fc3f5df4aefec026f610a3fcb08f19476492d69d2cb78b1c2eba259a8820e6a788b8b8a8a8a8a60405161053e96959493929190610ab2565b60405180910390a4505050505050505050565b8173ffffffffffffffffffffffffffffffffffffffff167f20c6237dac83526a849285a9f79d08a483291bdd3a056a0ef9ae94ecee1ad3568260405161059b911515815260200190565b60405180910390a25050565b8073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127860405160405180910390a35050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461062957600080fd5b919050565b600082601f83011261063f57600080fd5b813567ffffffffffffffff8082111561065a5761065a610b64565b604051601f83017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f011681019082821181831017156106a0576106a0610b64565b816040528381528660208588010111156106b957600080fd5b836020870160208301376000602085830101528094505050505092915050565b803563ffffffff8116811461062957600080fd5b803560ff8116811461062957600080fd5b80356bffffffffffffffffffffffff8116811461062957600080fd5b6000806040838503121561072d57600080fd5b61073683610605565b915061074460208401610605565b90509250929050565b6000806040838503121561076057600080fd5b61076983610605565b91506020830135801515811461077e57600080fd5b809150509250929050565b60006020828403121561079b57600080fd5b5035919050565b60008060008060008060008060006101208a8c0312156107c157600080fd5b8935985060208a013567ffffffffffffffff808211156107e057600080fd5b6107ec8d838e0161062e565b995060408c013591508082111561080257600080fd5b61080e8d838e0161062e565b985061081c60608d01610605565b975061082a60808d016106d9565b965061083860a08d01610605565b955060c08c013591508082111561084e57600080fd5b5061085b8c828d0161062e565b93505061086a60e08b016106fe565b91506108796101008b016106ed565b90509295985092959850929598565b60008060006060848603121561089d57600080fd5b83359250602084013567ffffffffffffffff8111156108bb57600080fd5b6108c78682870161062e565b925050604084013590509250925092565b600080600080600060a086880312156108f057600080fd5b8535600381106108ff57600080fd5b945061090d602087016106d9565b935061091b604087016106d9565b925061092960608701610605565b949793965091946080013592915050565b6000806000806080858703121561095057600080fd5b610959856106ed565b9350610967602086016106d9565b925061097560408601610605565b9150610983606086016106fe565b905092959194509250565b6000815180845260005b818110156109b457602081850181015186830182015201610998565b818111156109c6576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b60038110610a30577f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b9052565b60208101610a4282846109f9565b92915050565b60a08101610a5682886109f9565b63ffffffff808716602084015280861660408401525073ffffffffffffffffffffffffffffffffffffffff841660608301528260808301529695505050505050565b602081526000610aab602083018461098e565b9392505050565b60c081526000610ac560c083018961098e565b8281036020840152610ad7818961098e565b905063ffffffff8716604084015273ffffffffffffffffffffffffffffffffffffffff861660608401528281036080840152610b13818661098e565b9150506bffffffffffffffffffffffff831660a0830152979650505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var KeeperRegistrarMockABI = KeeperRegistrarMockMetaData.ABI

var KeeperRegistrarMockBin = KeeperRegistrarMockMetaData.Bin

func DeployKeeperRegistrarMock(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *KeeperRegistrarMock, error) {
	parsed, err := KeeperRegistrarMockMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KeeperRegistrarMockBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KeeperRegistrarMock{address: address, abi: *parsed, KeeperRegistrarMockCaller: KeeperRegistrarMockCaller{contract: contract}, KeeperRegistrarMockTransactor: KeeperRegistrarMockTransactor{contract: contract}, KeeperRegistrarMockFilterer: KeeperRegistrarMockFilterer{contract: contract}}, nil
}

type KeeperRegistrarMock struct {
	address common.Address
	abi     abi.ABI
	KeeperRegistrarMockCaller
	KeeperRegistrarMockTransactor
	KeeperRegistrarMockFilterer
}

type KeeperRegistrarMockCaller struct {
	contract *bind.BoundContract
}

type KeeperRegistrarMockTransactor struct {
	contract *bind.BoundContract
}

type KeeperRegistrarMockFilterer struct {
	contract *bind.BoundContract
}

type KeeperRegistrarMockSession struct {
	Contract     *KeeperRegistrarMock
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type KeeperRegistrarMockCallerSession struct {
	Contract *KeeperRegistrarMockCaller
	CallOpts bind.CallOpts
}

type KeeperRegistrarMockTransactorSession struct {
	Contract     *KeeperRegistrarMockTransactor
	TransactOpts bind.TransactOpts
}

type KeeperRegistrarMockRaw struct {
	Contract *KeeperRegistrarMock
}

type KeeperRegistrarMockCallerRaw struct {
	Contract *KeeperRegistrarMockCaller
}

type KeeperRegistrarMockTransactorRaw struct {
	Contract *KeeperRegistrarMockTransactor
}

func NewKeeperRegistrarMock(address common.Address, backend bind.ContractBackend) (*KeeperRegistrarMock, error) {
	abi, err := abi.JSON(strings.NewReader(KeeperRegistrarMockABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindKeeperRegistrarMock(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrarMock{address: address, abi: abi, KeeperRegistrarMockCaller: KeeperRegistrarMockCaller{contract: contract}, KeeperRegistrarMockTransactor: KeeperRegistrarMockTransactor{contract: contract}, KeeperRegistrarMockFilterer: KeeperRegistrarMockFilterer{contract: contract}}, nil
}

func NewKeeperRegistrarMockCaller(address common.Address, caller bind.ContractCaller) (*KeeperRegistrarMockCaller, error) {
	contract, err := bindKeeperRegistrarMock(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrarMockCaller{contract: contract}, nil
}

func NewKeeperRegistrarMockTransactor(address common.Address, transactor bind.ContractTransactor) (*KeeperRegistrarMockTransactor, error) {
	contract, err := bindKeeperRegistrarMock(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrarMockTransactor{contract: contract}, nil
}

func NewKeeperRegistrarMockFilterer(address common.Address, filterer bind.ContractFilterer) (*KeeperRegistrarMockFilterer, error) {
	contract, err := bindKeeperRegistrarMock(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrarMockFilterer{contract: contract}, nil
}

func bindKeeperRegistrarMock(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := KeeperRegistrarMockMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_KeeperRegistrarMock *KeeperRegistrarMockRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistrarMock.Contract.KeeperRegistrarMockCaller.contract.Call(opts, result, method, params...)
}

func (_KeeperRegistrarMock *KeeperRegistrarMockRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistrarMock.Contract.KeeperRegistrarMockTransactor.contract.Transfer(opts)
}

func (_KeeperRegistrarMock *KeeperRegistrarMockRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistrarMock.Contract.KeeperRegistrarMockTransactor.contract.Transact(opts, method, params...)
}

func (_KeeperRegistrarMock *KeeperRegistrarMockCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistrarMock.Contract.contract.Call(opts, result, method, params...)
}

func (_KeeperRegistrarMock *KeeperRegistrarMockTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistrarMock.Contract.contract.Transfer(opts)
}

func (_KeeperRegistrarMock *KeeperRegistrarMockTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistrarMock.Contract.contract.Transact(opts, method, params...)
}

func (_KeeperRegistrarMock *KeeperRegistrarMockCaller) GetRegistrationConfig(opts *bind.CallOpts) (GetRegistrationConfig,

	error) {
	var out []interface{}
	err := _KeeperRegistrarMock.contract.Call(opts, &out, "getRegistrationConfig")

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

func (_KeeperRegistrarMock *KeeperRegistrarMockSession) GetRegistrationConfig() (GetRegistrationConfig,

	error) {
	return _KeeperRegistrarMock.Contract.GetRegistrationConfig(&_KeeperRegistrarMock.CallOpts)
}

func (_KeeperRegistrarMock *KeeperRegistrarMockCallerSession) GetRegistrationConfig() (GetRegistrationConfig,

	error) {
	return _KeeperRegistrarMock.Contract.GetRegistrationConfig(&_KeeperRegistrarMock.CallOpts)
}

func (_KeeperRegistrarMock *KeeperRegistrarMockCaller) SApprovedCount(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _KeeperRegistrarMock.contract.Call(opts, &out, "s_approvedCount")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_KeeperRegistrarMock *KeeperRegistrarMockSession) SApprovedCount() (uint32, error) {
	return _KeeperRegistrarMock.Contract.SApprovedCount(&_KeeperRegistrarMock.CallOpts)
}

func (_KeeperRegistrarMock *KeeperRegistrarMockCallerSession) SApprovedCount() (uint32, error) {
	return _KeeperRegistrarMock.Contract.SApprovedCount(&_KeeperRegistrarMock.CallOpts)
}

func (_KeeperRegistrarMock *KeeperRegistrarMockCaller) SAutoApproveConfigType(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistrarMock.contract.Call(opts, &out, "s_autoApproveConfigType")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_KeeperRegistrarMock *KeeperRegistrarMockSession) SAutoApproveConfigType() (uint8, error) {
	return _KeeperRegistrarMock.Contract.SAutoApproveConfigType(&_KeeperRegistrarMock.CallOpts)
}

func (_KeeperRegistrarMock *KeeperRegistrarMockCallerSession) SAutoApproveConfigType() (uint8, error) {
	return _KeeperRegistrarMock.Contract.SAutoApproveConfigType(&_KeeperRegistrarMock.CallOpts)
}

func (_KeeperRegistrarMock *KeeperRegistrarMockCaller) SAutoApproveMaxAllowed(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _KeeperRegistrarMock.contract.Call(opts, &out, "s_autoApproveMaxAllowed")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_KeeperRegistrarMock *KeeperRegistrarMockSession) SAutoApproveMaxAllowed() (uint32, error) {
	return _KeeperRegistrarMock.Contract.SAutoApproveMaxAllowed(&_KeeperRegistrarMock.CallOpts)
}

func (_KeeperRegistrarMock *KeeperRegistrarMockCallerSession) SAutoApproveMaxAllowed() (uint32, error) {
	return _KeeperRegistrarMock.Contract.SAutoApproveMaxAllowed(&_KeeperRegistrarMock.CallOpts)
}

func (_KeeperRegistrarMock *KeeperRegistrarMockCaller) SKeeperRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistrarMock.contract.Call(opts, &out, "s_keeperRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistrarMock *KeeperRegistrarMockSession) SKeeperRegistry() (common.Address, error) {
	return _KeeperRegistrarMock.Contract.SKeeperRegistry(&_KeeperRegistrarMock.CallOpts)
}

func (_KeeperRegistrarMock *KeeperRegistrarMockCallerSession) SKeeperRegistry() (common.Address, error) {
	return _KeeperRegistrarMock.Contract.SKeeperRegistry(&_KeeperRegistrarMock.CallOpts)
}

func (_KeeperRegistrarMock *KeeperRegistrarMockCaller) SMinLINKJuels(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistrarMock.contract.Call(opts, &out, "s_minLINKJuels")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_KeeperRegistrarMock *KeeperRegistrarMockSession) SMinLINKJuels() (*big.Int, error) {
	return _KeeperRegistrarMock.Contract.SMinLINKJuels(&_KeeperRegistrarMock.CallOpts)
}

func (_KeeperRegistrarMock *KeeperRegistrarMockCallerSession) SMinLINKJuels() (*big.Int, error) {
	return _KeeperRegistrarMock.Contract.SMinLINKJuels(&_KeeperRegistrarMock.CallOpts)
}

func (_KeeperRegistrarMock *KeeperRegistrarMockTransactor) EmitAutoApproveAllowedSenderSet(opts *bind.TransactOpts, senderAddress common.Address, allowed bool) (*types.Transaction, error) {
	return _KeeperRegistrarMock.contract.Transact(opts, "emitAutoApproveAllowedSenderSet", senderAddress, allowed)
}

func (_KeeperRegistrarMock *KeeperRegistrarMockSession) EmitAutoApproveAllowedSenderSet(senderAddress common.Address, allowed bool) (*types.Transaction, error) {
	return _KeeperRegistrarMock.Contract.EmitAutoApproveAllowedSenderSet(&_KeeperRegistrarMock.TransactOpts, senderAddress, allowed)
}

func (_KeeperRegistrarMock *KeeperRegistrarMockTransactorSession) EmitAutoApproveAllowedSenderSet(senderAddress common.Address, allowed bool) (*types.Transaction, error) {
	return _KeeperRegistrarMock.Contract.EmitAutoApproveAllowedSenderSet(&_KeeperRegistrarMock.TransactOpts, senderAddress, allowed)
}

func (_KeeperRegistrarMock *KeeperRegistrarMockTransactor) EmitConfigChanged(opts *bind.TransactOpts, autoApproveConfigType uint8, autoApproveMaxAllowed uint32, keeperRegistry common.Address, minLINKJuels *big.Int) (*types.Transaction, error) {
	return _KeeperRegistrarMock.contract.Transact(opts, "emitConfigChanged", autoApproveConfigType, autoApproveMaxAllowed, keeperRegistry, minLINKJuels)
}

func (_KeeperRegistrarMock *KeeperRegistrarMockSession) EmitConfigChanged(autoApproveConfigType uint8, autoApproveMaxAllowed uint32, keeperRegistry common.Address, minLINKJuels *big.Int) (*types.Transaction, error) {
	return _KeeperRegistrarMock.Contract.EmitConfigChanged(&_KeeperRegistrarMock.TransactOpts, autoApproveConfigType, autoApproveMaxAllowed, keeperRegistry, minLINKJuels)
}

func (_KeeperRegistrarMock *KeeperRegistrarMockTransactorSession) EmitConfigChanged(autoApproveConfigType uint8, autoApproveMaxAllowed uint32, keeperRegistry common.Address, minLINKJuels *big.Int) (*types.Transaction, error) {
	return _KeeperRegistrarMock.Contract.EmitConfigChanged(&_KeeperRegistrarMock.TransactOpts, autoApproveConfigType, autoApproveMaxAllowed, keeperRegistry, minLINKJuels)
}

func (_KeeperRegistrarMock *KeeperRegistrarMockTransactor) EmitOwnershipTransferRequested(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistrarMock.contract.Transact(opts, "emitOwnershipTransferRequested", from, to)
}

func (_KeeperRegistrarMock *KeeperRegistrarMockSession) EmitOwnershipTransferRequested(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistrarMock.Contract.EmitOwnershipTransferRequested(&_KeeperRegistrarMock.TransactOpts, from, to)
}

func (_KeeperRegistrarMock *KeeperRegistrarMockTransactorSession) EmitOwnershipTransferRequested(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistrarMock.Contract.EmitOwnershipTransferRequested(&_KeeperRegistrarMock.TransactOpts, from, to)
}

func (_KeeperRegistrarMock *KeeperRegistrarMockTransactor) EmitOwnershipTransferred(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistrarMock.contract.Transact(opts, "emitOwnershipTransferred", from, to)
}

func (_KeeperRegistrarMock *KeeperRegistrarMockSession) EmitOwnershipTransferred(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistrarMock.Contract.EmitOwnershipTransferred(&_KeeperRegistrarMock.TransactOpts, from, to)
}

func (_KeeperRegistrarMock *KeeperRegistrarMockTransactorSession) EmitOwnershipTransferred(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistrarMock.Contract.EmitOwnershipTransferred(&_KeeperRegistrarMock.TransactOpts, from, to)
}

func (_KeeperRegistrarMock *KeeperRegistrarMockTransactor) EmitRegistrationApproved(opts *bind.TransactOpts, hash [32]byte, displayName string, upkeepId *big.Int) (*types.Transaction, error) {
	return _KeeperRegistrarMock.contract.Transact(opts, "emitRegistrationApproved", hash, displayName, upkeepId)
}

func (_KeeperRegistrarMock *KeeperRegistrarMockSession) EmitRegistrationApproved(hash [32]byte, displayName string, upkeepId *big.Int) (*types.Transaction, error) {
	return _KeeperRegistrarMock.Contract.EmitRegistrationApproved(&_KeeperRegistrarMock.TransactOpts, hash, displayName, upkeepId)
}

func (_KeeperRegistrarMock *KeeperRegistrarMockTransactorSession) EmitRegistrationApproved(hash [32]byte, displayName string, upkeepId *big.Int) (*types.Transaction, error) {
	return _KeeperRegistrarMock.Contract.EmitRegistrationApproved(&_KeeperRegistrarMock.TransactOpts, hash, displayName, upkeepId)
}

func (_KeeperRegistrarMock *KeeperRegistrarMockTransactor) EmitRegistrationRejected(opts *bind.TransactOpts, hash [32]byte) (*types.Transaction, error) {
	return _KeeperRegistrarMock.contract.Transact(opts, "emitRegistrationRejected", hash)
}

func (_KeeperRegistrarMock *KeeperRegistrarMockSession) EmitRegistrationRejected(hash [32]byte) (*types.Transaction, error) {
	return _KeeperRegistrarMock.Contract.EmitRegistrationRejected(&_KeeperRegistrarMock.TransactOpts, hash)
}

func (_KeeperRegistrarMock *KeeperRegistrarMockTransactorSession) EmitRegistrationRejected(hash [32]byte) (*types.Transaction, error) {
	return _KeeperRegistrarMock.Contract.EmitRegistrationRejected(&_KeeperRegistrarMock.TransactOpts, hash)
}

func (_KeeperRegistrarMock *KeeperRegistrarMockTransactor) EmitRegistrationRequested(opts *bind.TransactOpts, hash [32]byte, name string, encryptedEmail []byte, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, amount *big.Int, source uint8) (*types.Transaction, error) {
	return _KeeperRegistrarMock.contract.Transact(opts, "emitRegistrationRequested", hash, name, encryptedEmail, upkeepContract, gasLimit, adminAddress, checkData, amount, source)
}

func (_KeeperRegistrarMock *KeeperRegistrarMockSession) EmitRegistrationRequested(hash [32]byte, name string, encryptedEmail []byte, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, amount *big.Int, source uint8) (*types.Transaction, error) {
	return _KeeperRegistrarMock.Contract.EmitRegistrationRequested(&_KeeperRegistrarMock.TransactOpts, hash, name, encryptedEmail, upkeepContract, gasLimit, adminAddress, checkData, amount, source)
}

func (_KeeperRegistrarMock *KeeperRegistrarMockTransactorSession) EmitRegistrationRequested(hash [32]byte, name string, encryptedEmail []byte, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, amount *big.Int, source uint8) (*types.Transaction, error) {
	return _KeeperRegistrarMock.Contract.EmitRegistrationRequested(&_KeeperRegistrarMock.TransactOpts, hash, name, encryptedEmail, upkeepContract, gasLimit, adminAddress, checkData, amount, source)
}

func (_KeeperRegistrarMock *KeeperRegistrarMockTransactor) SetRegistrationConfig(opts *bind.TransactOpts, _autoApproveConfigType uint8, _autoApproveMaxAllowed uint32, _approvedCount uint32, _keeperRegistry common.Address, _minLINKJuels *big.Int) (*types.Transaction, error) {
	return _KeeperRegistrarMock.contract.Transact(opts, "setRegistrationConfig", _autoApproveConfigType, _autoApproveMaxAllowed, _approvedCount, _keeperRegistry, _minLINKJuels)
}

func (_KeeperRegistrarMock *KeeperRegistrarMockSession) SetRegistrationConfig(_autoApproveConfigType uint8, _autoApproveMaxAllowed uint32, _approvedCount uint32, _keeperRegistry common.Address, _minLINKJuels *big.Int) (*types.Transaction, error) {
	return _KeeperRegistrarMock.Contract.SetRegistrationConfig(&_KeeperRegistrarMock.TransactOpts, _autoApproveConfigType, _autoApproveMaxAllowed, _approvedCount, _keeperRegistry, _minLINKJuels)
}

func (_KeeperRegistrarMock *KeeperRegistrarMockTransactorSession) SetRegistrationConfig(_autoApproveConfigType uint8, _autoApproveMaxAllowed uint32, _approvedCount uint32, _keeperRegistry common.Address, _minLINKJuels *big.Int) (*types.Transaction, error) {
	return _KeeperRegistrarMock.Contract.SetRegistrationConfig(&_KeeperRegistrarMock.TransactOpts, _autoApproveConfigType, _autoApproveMaxAllowed, _approvedCount, _keeperRegistry, _minLINKJuels)
}

type KeeperRegistrarMockAutoApproveAllowedSenderSetIterator struct {
	Event *KeeperRegistrarMockAutoApproveAllowedSenderSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistrarMockAutoApproveAllowedSenderSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistrarMockAutoApproveAllowedSenderSet)
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
		it.Event = new(KeeperRegistrarMockAutoApproveAllowedSenderSet)
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

func (it *KeeperRegistrarMockAutoApproveAllowedSenderSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistrarMockAutoApproveAllowedSenderSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistrarMockAutoApproveAllowedSenderSet struct {
	SenderAddress common.Address
	Allowed       bool
	Raw           types.Log
}

func (_KeeperRegistrarMock *KeeperRegistrarMockFilterer) FilterAutoApproveAllowedSenderSet(opts *bind.FilterOpts, senderAddress []common.Address) (*KeeperRegistrarMockAutoApproveAllowedSenderSetIterator, error) {

	var senderAddressRule []interface{}
	for _, senderAddressItem := range senderAddress {
		senderAddressRule = append(senderAddressRule, senderAddressItem)
	}

	logs, sub, err := _KeeperRegistrarMock.contract.FilterLogs(opts, "AutoApproveAllowedSenderSet", senderAddressRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrarMockAutoApproveAllowedSenderSetIterator{contract: _KeeperRegistrarMock.contract, event: "AutoApproveAllowedSenderSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistrarMock *KeeperRegistrarMockFilterer) WatchAutoApproveAllowedSenderSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarMockAutoApproveAllowedSenderSet, senderAddress []common.Address) (event.Subscription, error) {

	var senderAddressRule []interface{}
	for _, senderAddressItem := range senderAddress {
		senderAddressRule = append(senderAddressRule, senderAddressItem)
	}

	logs, sub, err := _KeeperRegistrarMock.contract.WatchLogs(opts, "AutoApproveAllowedSenderSet", senderAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistrarMockAutoApproveAllowedSenderSet)
				if err := _KeeperRegistrarMock.contract.UnpackLog(event, "AutoApproveAllowedSenderSet", log); err != nil {
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

func (_KeeperRegistrarMock *KeeperRegistrarMockFilterer) ParseAutoApproveAllowedSenderSet(log types.Log) (*KeeperRegistrarMockAutoApproveAllowedSenderSet, error) {
	event := new(KeeperRegistrarMockAutoApproveAllowedSenderSet)
	if err := _KeeperRegistrarMock.contract.UnpackLog(event, "AutoApproveAllowedSenderSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistrarMockConfigChangedIterator struct {
	Event *KeeperRegistrarMockConfigChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistrarMockConfigChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistrarMockConfigChanged)
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
		it.Event = new(KeeperRegistrarMockConfigChanged)
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

func (it *KeeperRegistrarMockConfigChangedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistrarMockConfigChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistrarMockConfigChanged struct {
	AutoApproveConfigType uint8
	AutoApproveMaxAllowed uint32
	KeeperRegistry        common.Address
	MinLINKJuels          *big.Int
	Raw                   types.Log
}

func (_KeeperRegistrarMock *KeeperRegistrarMockFilterer) FilterConfigChanged(opts *bind.FilterOpts) (*KeeperRegistrarMockConfigChangedIterator, error) {

	logs, sub, err := _KeeperRegistrarMock.contract.FilterLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrarMockConfigChangedIterator{contract: _KeeperRegistrarMock.contract, event: "ConfigChanged", logs: logs, sub: sub}, nil
}

func (_KeeperRegistrarMock *KeeperRegistrarMockFilterer) WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarMockConfigChanged) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistrarMock.contract.WatchLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistrarMockConfigChanged)
				if err := _KeeperRegistrarMock.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
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

func (_KeeperRegistrarMock *KeeperRegistrarMockFilterer) ParseConfigChanged(log types.Log) (*KeeperRegistrarMockConfigChanged, error) {
	event := new(KeeperRegistrarMockConfigChanged)
	if err := _KeeperRegistrarMock.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistrarMockOwnershipTransferRequestedIterator struct {
	Event *KeeperRegistrarMockOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistrarMockOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistrarMockOwnershipTransferRequested)
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
		it.Event = new(KeeperRegistrarMockOwnershipTransferRequested)
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

func (it *KeeperRegistrarMockOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistrarMockOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistrarMockOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistrarMock *KeeperRegistrarMockFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistrarMockOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistrarMock.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrarMockOwnershipTransferRequestedIterator{contract: _KeeperRegistrarMock.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeeperRegistrarMock *KeeperRegistrarMockFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarMockOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistrarMock.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistrarMockOwnershipTransferRequested)
				if err := _KeeperRegistrarMock.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_KeeperRegistrarMock *KeeperRegistrarMockFilterer) ParseOwnershipTransferRequested(log types.Log) (*KeeperRegistrarMockOwnershipTransferRequested, error) {
	event := new(KeeperRegistrarMockOwnershipTransferRequested)
	if err := _KeeperRegistrarMock.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistrarMockOwnershipTransferredIterator struct {
	Event *KeeperRegistrarMockOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistrarMockOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistrarMockOwnershipTransferred)
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
		it.Event = new(KeeperRegistrarMockOwnershipTransferred)
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

func (it *KeeperRegistrarMockOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistrarMockOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistrarMockOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistrarMock *KeeperRegistrarMockFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistrarMockOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistrarMock.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrarMockOwnershipTransferredIterator{contract: _KeeperRegistrarMock.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_KeeperRegistrarMock *KeeperRegistrarMockFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarMockOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistrarMock.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistrarMockOwnershipTransferred)
				if err := _KeeperRegistrarMock.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_KeeperRegistrarMock *KeeperRegistrarMockFilterer) ParseOwnershipTransferred(log types.Log) (*KeeperRegistrarMockOwnershipTransferred, error) {
	event := new(KeeperRegistrarMockOwnershipTransferred)
	if err := _KeeperRegistrarMock.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistrarMockRegistrationApprovedIterator struct {
	Event *KeeperRegistrarMockRegistrationApproved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistrarMockRegistrationApprovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistrarMockRegistrationApproved)
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
		it.Event = new(KeeperRegistrarMockRegistrationApproved)
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

func (it *KeeperRegistrarMockRegistrationApprovedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistrarMockRegistrationApprovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistrarMockRegistrationApproved struct {
	Hash        [32]byte
	DisplayName string
	UpkeepId    *big.Int
	Raw         types.Log
}

func (_KeeperRegistrarMock *KeeperRegistrarMockFilterer) FilterRegistrationApproved(opts *bind.FilterOpts, hash [][32]byte, upkeepId []*big.Int) (*KeeperRegistrarMockRegistrationApprovedIterator, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}

	var upkeepIdRule []interface{}
	for _, upkeepIdItem := range upkeepId {
		upkeepIdRule = append(upkeepIdRule, upkeepIdItem)
	}

	logs, sub, err := _KeeperRegistrarMock.contract.FilterLogs(opts, "RegistrationApproved", hashRule, upkeepIdRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrarMockRegistrationApprovedIterator{contract: _KeeperRegistrarMock.contract, event: "RegistrationApproved", logs: logs, sub: sub}, nil
}

func (_KeeperRegistrarMock *KeeperRegistrarMockFilterer) WatchRegistrationApproved(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarMockRegistrationApproved, hash [][32]byte, upkeepId []*big.Int) (event.Subscription, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}

	var upkeepIdRule []interface{}
	for _, upkeepIdItem := range upkeepId {
		upkeepIdRule = append(upkeepIdRule, upkeepIdItem)
	}

	logs, sub, err := _KeeperRegistrarMock.contract.WatchLogs(opts, "RegistrationApproved", hashRule, upkeepIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistrarMockRegistrationApproved)
				if err := _KeeperRegistrarMock.contract.UnpackLog(event, "RegistrationApproved", log); err != nil {
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

func (_KeeperRegistrarMock *KeeperRegistrarMockFilterer) ParseRegistrationApproved(log types.Log) (*KeeperRegistrarMockRegistrationApproved, error) {
	event := new(KeeperRegistrarMockRegistrationApproved)
	if err := _KeeperRegistrarMock.contract.UnpackLog(event, "RegistrationApproved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistrarMockRegistrationRejectedIterator struct {
	Event *KeeperRegistrarMockRegistrationRejected

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistrarMockRegistrationRejectedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistrarMockRegistrationRejected)
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
		it.Event = new(KeeperRegistrarMockRegistrationRejected)
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

func (it *KeeperRegistrarMockRegistrationRejectedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistrarMockRegistrationRejectedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistrarMockRegistrationRejected struct {
	Hash [32]byte
	Raw  types.Log
}

func (_KeeperRegistrarMock *KeeperRegistrarMockFilterer) FilterRegistrationRejected(opts *bind.FilterOpts, hash [][32]byte) (*KeeperRegistrarMockRegistrationRejectedIterator, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}

	logs, sub, err := _KeeperRegistrarMock.contract.FilterLogs(opts, "RegistrationRejected", hashRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrarMockRegistrationRejectedIterator{contract: _KeeperRegistrarMock.contract, event: "RegistrationRejected", logs: logs, sub: sub}, nil
}

func (_KeeperRegistrarMock *KeeperRegistrarMockFilterer) WatchRegistrationRejected(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarMockRegistrationRejected, hash [][32]byte) (event.Subscription, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}

	logs, sub, err := _KeeperRegistrarMock.contract.WatchLogs(opts, "RegistrationRejected", hashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistrarMockRegistrationRejected)
				if err := _KeeperRegistrarMock.contract.UnpackLog(event, "RegistrationRejected", log); err != nil {
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

func (_KeeperRegistrarMock *KeeperRegistrarMockFilterer) ParseRegistrationRejected(log types.Log) (*KeeperRegistrarMockRegistrationRejected, error) {
	event := new(KeeperRegistrarMockRegistrationRejected)
	if err := _KeeperRegistrarMock.contract.UnpackLog(event, "RegistrationRejected", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistrarMockRegistrationRequestedIterator struct {
	Event *KeeperRegistrarMockRegistrationRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistrarMockRegistrationRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistrarMockRegistrationRequested)
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
		it.Event = new(KeeperRegistrarMockRegistrationRequested)
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

func (it *KeeperRegistrarMockRegistrationRequestedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistrarMockRegistrationRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistrarMockRegistrationRequested struct {
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

func (_KeeperRegistrarMock *KeeperRegistrarMockFilterer) FilterRegistrationRequested(opts *bind.FilterOpts, hash [][32]byte, upkeepContract []common.Address, source []uint8) (*KeeperRegistrarMockRegistrationRequestedIterator, error) {

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

	logs, sub, err := _KeeperRegistrarMock.contract.FilterLogs(opts, "RegistrationRequested", hashRule, upkeepContractRule, sourceRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrarMockRegistrationRequestedIterator{contract: _KeeperRegistrarMock.contract, event: "RegistrationRequested", logs: logs, sub: sub}, nil
}

func (_KeeperRegistrarMock *KeeperRegistrarMockFilterer) WatchRegistrationRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarMockRegistrationRequested, hash [][32]byte, upkeepContract []common.Address, source []uint8) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistrarMock.contract.WatchLogs(opts, "RegistrationRequested", hashRule, upkeepContractRule, sourceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistrarMockRegistrationRequested)
				if err := _KeeperRegistrarMock.contract.UnpackLog(event, "RegistrationRequested", log); err != nil {
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

func (_KeeperRegistrarMock *KeeperRegistrarMockFilterer) ParseRegistrationRequested(log types.Log) (*KeeperRegistrarMockRegistrationRequested, error) {
	event := new(KeeperRegistrarMockRegistrationRequested)
	if err := _KeeperRegistrarMock.contract.UnpackLog(event, "RegistrationRequested", log); err != nil {
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

func (_KeeperRegistrarMock *KeeperRegistrarMock) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _KeeperRegistrarMock.abi.Events["AutoApproveAllowedSenderSet"].ID:
		return _KeeperRegistrarMock.ParseAutoApproveAllowedSenderSet(log)
	case _KeeperRegistrarMock.abi.Events["ConfigChanged"].ID:
		return _KeeperRegistrarMock.ParseConfigChanged(log)
	case _KeeperRegistrarMock.abi.Events["OwnershipTransferRequested"].ID:
		return _KeeperRegistrarMock.ParseOwnershipTransferRequested(log)
	case _KeeperRegistrarMock.abi.Events["OwnershipTransferred"].ID:
		return _KeeperRegistrarMock.ParseOwnershipTransferred(log)
	case _KeeperRegistrarMock.abi.Events["RegistrationApproved"].ID:
		return _KeeperRegistrarMock.ParseRegistrationApproved(log)
	case _KeeperRegistrarMock.abi.Events["RegistrationRejected"].ID:
		return _KeeperRegistrarMock.ParseRegistrationRejected(log)
	case _KeeperRegistrarMock.abi.Events["RegistrationRequested"].ID:
		return _KeeperRegistrarMock.ParseRegistrationRequested(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (KeeperRegistrarMockAutoApproveAllowedSenderSet) Topic() common.Hash {
	return common.HexToHash("0x20c6237dac83526a849285a9f79d08a483291bdd3a056a0ef9ae94ecee1ad356")
}

func (KeeperRegistrarMockConfigChanged) Topic() common.Hash {
	return common.HexToHash("0x6293a703ec7145dfa23c5cde2e627d6a02e153fc2e9c03b14d1e22cbb4a7e9cd")
}

func (KeeperRegistrarMockOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (KeeperRegistrarMockOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (KeeperRegistrarMockRegistrationApproved) Topic() common.Hash {
	return common.HexToHash("0xb9a292fb7e3edd920cd2d2829a3615a640c43fd7de0a0820aa0668feb4c37d4b")
}

func (KeeperRegistrarMockRegistrationRejected) Topic() common.Hash {
	return common.HexToHash("0x3663fb28ebc87645eb972c9dad8521bf665c623f287e79f1c56f1eb374b82a22")
}

func (KeeperRegistrarMockRegistrationRequested) Topic() common.Hash {
	return common.HexToHash("0xc3f5df4aefec026f610a3fcb08f19476492d69d2cb78b1c2eba259a8820e6a78")
}

func (_KeeperRegistrarMock *KeeperRegistrarMock) Address() common.Address {
	return _KeeperRegistrarMock.address
}

type KeeperRegistrarMockInterface interface {
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

	FilterAutoApproveAllowedSenderSet(opts *bind.FilterOpts, senderAddress []common.Address) (*KeeperRegistrarMockAutoApproveAllowedSenderSetIterator, error)

	WatchAutoApproveAllowedSenderSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarMockAutoApproveAllowedSenderSet, senderAddress []common.Address) (event.Subscription, error)

	ParseAutoApproveAllowedSenderSet(log types.Log) (*KeeperRegistrarMockAutoApproveAllowedSenderSet, error)

	FilterConfigChanged(opts *bind.FilterOpts) (*KeeperRegistrarMockConfigChangedIterator, error)

	WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarMockConfigChanged) (event.Subscription, error)

	ParseConfigChanged(log types.Log) (*KeeperRegistrarMockConfigChanged, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistrarMockOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarMockOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*KeeperRegistrarMockOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistrarMockOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarMockOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*KeeperRegistrarMockOwnershipTransferred, error)

	FilterRegistrationApproved(opts *bind.FilterOpts, hash [][32]byte, upkeepId []*big.Int) (*KeeperRegistrarMockRegistrationApprovedIterator, error)

	WatchRegistrationApproved(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarMockRegistrationApproved, hash [][32]byte, upkeepId []*big.Int) (event.Subscription, error)

	ParseRegistrationApproved(log types.Log) (*KeeperRegistrarMockRegistrationApproved, error)

	FilterRegistrationRejected(opts *bind.FilterOpts, hash [][32]byte) (*KeeperRegistrarMockRegistrationRejectedIterator, error)

	WatchRegistrationRejected(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarMockRegistrationRejected, hash [][32]byte) (event.Subscription, error)

	ParseRegistrationRejected(log types.Log) (*KeeperRegistrarMockRegistrationRejected, error)

	FilterRegistrationRequested(opts *bind.FilterOpts, hash [][32]byte, upkeepContract []common.Address, source []uint8) (*KeeperRegistrarMockRegistrationRequestedIterator, error)

	WatchRegistrationRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarMockRegistrationRequested, hash [][32]byte, upkeepContract []common.Address, source []uint8) (event.Subscription, error)

	ParseRegistrationRequested(log types.Log) (*KeeperRegistrarMockRegistrationRequested, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
