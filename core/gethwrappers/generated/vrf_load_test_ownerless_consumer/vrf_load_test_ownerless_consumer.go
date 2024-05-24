// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_load_test_ownerless_consumer

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

var VRFLoadTestOwnerlessConsumerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_price\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"PRICE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"randomness\",\"type\":\"uint256\"}],\"name\":\"rawFulfillRandomness\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_responseCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60e060405234801561001057600080fd5b5060405161078d38038061078d83398101604081905261002f9161006e565b6001600160601b0319606093841b811660a0529190921b1660805260c0526100aa565b80516001600160a01b038116811461006957600080fd5b919050565b60008060006060848603121561008357600080fd5b61008c84610052565b925061009a60208501610052565b9150604084015190509250925092565b60805160601c60a05160601c60c05161068d6101006000396000818160560152818161022501528181610255015261027f01526000818160d3015261030c01526000818161018501526102d0015261068d6000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80638d859f3e1461005157806394985ddd1461008a578063a4c0ed361461009f578063dc1670db146100b2575b600080fd5b6100787f000000000000000000000000000000000000000000000000000000000000000081565b60405190815260200160405180910390f35b61009d610098366004610546565b6100bb565b005b61009d6100ad366004610462565b61016d565b61007860015481565b3373ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161461015f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f4f6e6c7920565246436f6f7264696e61746f722063616e2066756c66696c6c0060448201526064015b60405180910390fd5b61016982826102b3565b5050565b3373ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161461020c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f6f6e6c792063616c6c61626c652066726f6d204c494e4b0000000000000000006044820152606401610156565b600061021a8284018461052d565b905060005b8461024a7f000000000000000000000000000000000000000000000000000000000000000083610600565b116102ab57610279827f00000000000000000000000000000000000000000000000000000000000000006102cc565b506102a47f000000000000000000000000000000000000000000000000000000000000000082610600565b905061021f565b505050505050565b600180549060006102c383610618565b91905055505050565b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16634000aea07f000000000000000000000000000000000000000000000000000000000000000084866000604051602001610349929190918252602082015260400190565b6040516020818303038152906040526040518463ffffffff1660e01b815260040161037693929190610568565b602060405180830381600087803b15801561039057600080fd5b505af11580156103a4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103c89190610504565b5060008381526020818152604080832054815180840188905280830185905230606082015260808082018390528351808303909101815260a090910190925281519183019190912086845292909152610422906001610600565b6000858152602081815260409182902092909255805180830187905280820184905281518082038301815260609091019091528051910120949350505050565b6000806000806060858703121561047857600080fd5b843573ffffffffffffffffffffffffffffffffffffffff8116811461049c57600080fd5b935060208501359250604085013567ffffffffffffffff808211156104c057600080fd5b818701915087601f8301126104d457600080fd5b8135818111156104e357600080fd5b8860208285010111156104f557600080fd5b95989497505060200194505050565b60006020828403121561051657600080fd5b8151801515811461052657600080fd5b9392505050565b60006020828403121561053f57600080fd5b5035919050565b6000806040838503121561055957600080fd5b50508035926020909101359150565b73ffffffffffffffffffffffffffffffffffffffff8416815260006020848184015260606040840152835180606085015260005b818110156105b85785810183015185820160800152820161059c565b818111156105ca576000608083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160800195945050505050565b6000821982111561061357610613610651565b500190565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82141561064a5761064a610651565b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fdfea164736f6c6343000806000a",
}

var VRFLoadTestOwnerlessConsumerABI = VRFLoadTestOwnerlessConsumerMetaData.ABI

var VRFLoadTestOwnerlessConsumerBin = VRFLoadTestOwnerlessConsumerMetaData.Bin

func DeployVRFLoadTestOwnerlessConsumer(auth *bind.TransactOpts, backend bind.ContractBackend, _vrfCoordinator common.Address, _link common.Address, _price *big.Int) (common.Address, *types.Transaction, *VRFLoadTestOwnerlessConsumer, error) {
	parsed, err := VRFLoadTestOwnerlessConsumerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFLoadTestOwnerlessConsumerBin), backend, _vrfCoordinator, _link, _price)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFLoadTestOwnerlessConsumer{address: address, abi: *parsed, VRFLoadTestOwnerlessConsumerCaller: VRFLoadTestOwnerlessConsumerCaller{contract: contract}, VRFLoadTestOwnerlessConsumerTransactor: VRFLoadTestOwnerlessConsumerTransactor{contract: contract}, VRFLoadTestOwnerlessConsumerFilterer: VRFLoadTestOwnerlessConsumerFilterer{contract: contract}}, nil
}

type VRFLoadTestOwnerlessConsumer struct {
	address common.Address
	abi     abi.ABI
	VRFLoadTestOwnerlessConsumerCaller
	VRFLoadTestOwnerlessConsumerTransactor
	VRFLoadTestOwnerlessConsumerFilterer
}

type VRFLoadTestOwnerlessConsumerCaller struct {
	contract *bind.BoundContract
}

type VRFLoadTestOwnerlessConsumerTransactor struct {
	contract *bind.BoundContract
}

type VRFLoadTestOwnerlessConsumerFilterer struct {
	contract *bind.BoundContract
}

type VRFLoadTestOwnerlessConsumerSession struct {
	Contract     *VRFLoadTestOwnerlessConsumer
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFLoadTestOwnerlessConsumerCallerSession struct {
	Contract *VRFLoadTestOwnerlessConsumerCaller
	CallOpts bind.CallOpts
}

type VRFLoadTestOwnerlessConsumerTransactorSession struct {
	Contract     *VRFLoadTestOwnerlessConsumerTransactor
	TransactOpts bind.TransactOpts
}

type VRFLoadTestOwnerlessConsumerRaw struct {
	Contract *VRFLoadTestOwnerlessConsumer
}

type VRFLoadTestOwnerlessConsumerCallerRaw struct {
	Contract *VRFLoadTestOwnerlessConsumerCaller
}

type VRFLoadTestOwnerlessConsumerTransactorRaw struct {
	Contract *VRFLoadTestOwnerlessConsumerTransactor
}

func NewVRFLoadTestOwnerlessConsumer(address common.Address, backend bind.ContractBackend) (*VRFLoadTestOwnerlessConsumer, error) {
	abi, err := abi.JSON(strings.NewReader(VRFLoadTestOwnerlessConsumerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFLoadTestOwnerlessConsumer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFLoadTestOwnerlessConsumer{address: address, abi: abi, VRFLoadTestOwnerlessConsumerCaller: VRFLoadTestOwnerlessConsumerCaller{contract: contract}, VRFLoadTestOwnerlessConsumerTransactor: VRFLoadTestOwnerlessConsumerTransactor{contract: contract}, VRFLoadTestOwnerlessConsumerFilterer: VRFLoadTestOwnerlessConsumerFilterer{contract: contract}}, nil
}

func NewVRFLoadTestOwnerlessConsumerCaller(address common.Address, caller bind.ContractCaller) (*VRFLoadTestOwnerlessConsumerCaller, error) {
	contract, err := bindVRFLoadTestOwnerlessConsumer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFLoadTestOwnerlessConsumerCaller{contract: contract}, nil
}

func NewVRFLoadTestOwnerlessConsumerTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFLoadTestOwnerlessConsumerTransactor, error) {
	contract, err := bindVRFLoadTestOwnerlessConsumer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFLoadTestOwnerlessConsumerTransactor{contract: contract}, nil
}

func NewVRFLoadTestOwnerlessConsumerFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFLoadTestOwnerlessConsumerFilterer, error) {
	contract, err := bindVRFLoadTestOwnerlessConsumer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFLoadTestOwnerlessConsumerFilterer{contract: contract}, nil
}

func bindVRFLoadTestOwnerlessConsumer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFLoadTestOwnerlessConsumerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFLoadTestOwnerlessConsumer *VRFLoadTestOwnerlessConsumerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFLoadTestOwnerlessConsumer.Contract.VRFLoadTestOwnerlessConsumerCaller.contract.Call(opts, result, method, params...)
}

func (_VRFLoadTestOwnerlessConsumer *VRFLoadTestOwnerlessConsumerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFLoadTestOwnerlessConsumer.Contract.VRFLoadTestOwnerlessConsumerTransactor.contract.Transfer(opts)
}

func (_VRFLoadTestOwnerlessConsumer *VRFLoadTestOwnerlessConsumerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFLoadTestOwnerlessConsumer.Contract.VRFLoadTestOwnerlessConsumerTransactor.contract.Transact(opts, method, params...)
}

func (_VRFLoadTestOwnerlessConsumer *VRFLoadTestOwnerlessConsumerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFLoadTestOwnerlessConsumer.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFLoadTestOwnerlessConsumer *VRFLoadTestOwnerlessConsumerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFLoadTestOwnerlessConsumer.Contract.contract.Transfer(opts)
}

func (_VRFLoadTestOwnerlessConsumer *VRFLoadTestOwnerlessConsumerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFLoadTestOwnerlessConsumer.Contract.contract.Transact(opts, method, params...)
}

func (_VRFLoadTestOwnerlessConsumer *VRFLoadTestOwnerlessConsumerCaller) PRICE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFLoadTestOwnerlessConsumer.contract.Call(opts, &out, "PRICE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFLoadTestOwnerlessConsumer *VRFLoadTestOwnerlessConsumerSession) PRICE() (*big.Int, error) {
	return _VRFLoadTestOwnerlessConsumer.Contract.PRICE(&_VRFLoadTestOwnerlessConsumer.CallOpts)
}

func (_VRFLoadTestOwnerlessConsumer *VRFLoadTestOwnerlessConsumerCallerSession) PRICE() (*big.Int, error) {
	return _VRFLoadTestOwnerlessConsumer.Contract.PRICE(&_VRFLoadTestOwnerlessConsumer.CallOpts)
}

func (_VRFLoadTestOwnerlessConsumer *VRFLoadTestOwnerlessConsumerCaller) SResponseCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFLoadTestOwnerlessConsumer.contract.Call(opts, &out, "s_responseCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFLoadTestOwnerlessConsumer *VRFLoadTestOwnerlessConsumerSession) SResponseCount() (*big.Int, error) {
	return _VRFLoadTestOwnerlessConsumer.Contract.SResponseCount(&_VRFLoadTestOwnerlessConsumer.CallOpts)
}

func (_VRFLoadTestOwnerlessConsumer *VRFLoadTestOwnerlessConsumerCallerSession) SResponseCount() (*big.Int, error) {
	return _VRFLoadTestOwnerlessConsumer.Contract.SResponseCount(&_VRFLoadTestOwnerlessConsumer.CallOpts)
}

func (_VRFLoadTestOwnerlessConsumer *VRFLoadTestOwnerlessConsumerTransactor) OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _VRFLoadTestOwnerlessConsumer.contract.Transact(opts, "onTokenTransfer", arg0, _amount, _data)
}

func (_VRFLoadTestOwnerlessConsumer *VRFLoadTestOwnerlessConsumerSession) OnTokenTransfer(arg0 common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _VRFLoadTestOwnerlessConsumer.Contract.OnTokenTransfer(&_VRFLoadTestOwnerlessConsumer.TransactOpts, arg0, _amount, _data)
}

func (_VRFLoadTestOwnerlessConsumer *VRFLoadTestOwnerlessConsumerTransactorSession) OnTokenTransfer(arg0 common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _VRFLoadTestOwnerlessConsumer.Contract.OnTokenTransfer(&_VRFLoadTestOwnerlessConsumer.TransactOpts, arg0, _amount, _data)
}

func (_VRFLoadTestOwnerlessConsumer *VRFLoadTestOwnerlessConsumerTransactor) RawFulfillRandomness(opts *bind.TransactOpts, requestId [32]byte, randomness *big.Int) (*types.Transaction, error) {
	return _VRFLoadTestOwnerlessConsumer.contract.Transact(opts, "rawFulfillRandomness", requestId, randomness)
}

func (_VRFLoadTestOwnerlessConsumer *VRFLoadTestOwnerlessConsumerSession) RawFulfillRandomness(requestId [32]byte, randomness *big.Int) (*types.Transaction, error) {
	return _VRFLoadTestOwnerlessConsumer.Contract.RawFulfillRandomness(&_VRFLoadTestOwnerlessConsumer.TransactOpts, requestId, randomness)
}

func (_VRFLoadTestOwnerlessConsumer *VRFLoadTestOwnerlessConsumerTransactorSession) RawFulfillRandomness(requestId [32]byte, randomness *big.Int) (*types.Transaction, error) {
	return _VRFLoadTestOwnerlessConsumer.Contract.RawFulfillRandomness(&_VRFLoadTestOwnerlessConsumer.TransactOpts, requestId, randomness)
}

func (_VRFLoadTestOwnerlessConsumer *VRFLoadTestOwnerlessConsumer) Address() common.Address {
	return _VRFLoadTestOwnerlessConsumer.address
}

type VRFLoadTestOwnerlessConsumerInterface interface {
	PRICE(opts *bind.CallOpts) (*big.Int, error)

	SResponseCount(opts *bind.CallOpts) (*big.Int, error)

	OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error)

	RawFulfillRandomness(opts *bind.TransactOpts, requestId [32]byte, randomness *big.Int) (*types.Transaction, error)

	Address() common.Address
}
