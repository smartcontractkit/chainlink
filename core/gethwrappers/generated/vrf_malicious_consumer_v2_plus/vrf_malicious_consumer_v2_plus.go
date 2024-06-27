// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_malicious_consumer_v2_plus

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

var VRFMaliciousConsumerV2PlusMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlyCoordinatorCanFulfill\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"}],\"name\":\"OnlyOwnerOrCoordinator\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"createSubscriptionAndFund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"requestRandomness\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_gasAvailable\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_randomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_vrfCoordinator\",\"outputs\":[{\"internalType\":\"contractIVRFCoordinatorV2Plus\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vrfCoordinator\",\"type\":\"address\"}],\"name\":\"setCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"name\":\"updateSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162001239380380620012398339810160408190526200003491620001c2565b8133806000816200008c5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000bf57620000bf81620000f9565b5050600280546001600160a01b039384166001600160a01b0319918216179091556005805494909316931692909217905550620001fa9050565b6001600160a01b038116331415620001545760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000083565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b0381168114620001bd57600080fd5b919050565b60008060408385031215620001d657600080fd5b620001e183620001a5565b9150620001f160208401620001a5565b90509250929050565b61102f806200020a6000396000f3fe608060405234801561001057600080fd5b50600436106100d45760003560e01c80639eccacf611610081578063f08c5daa1161005b578063f08c5daa146101bd578063f2fde38b146101c6578063f6eaffc8146101d957600080fd5b80639eccacf614610181578063cf62c8ab146101a1578063e89e106a146101b457600080fd5b806379ba5097116100b257806379ba5097146101275780638da5cb5b1461012f5780638ea981171461016e57600080fd5b80631fe543e3146100d957806336bfffed146100ee5780635e3b709f14610101575b600080fd5b6100ec6100e7366004610d03565b6101ec565b005b6100ec6100fc366004610c0b565b610272565b61011461010f366004610cd1565b6103aa565b6040519081526020015b60405180910390f35b6100ec6104a0565b60005473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161011e565b6100ec61017c366004610bf0565b61059d565b6002546101499073ffffffffffffffffffffffffffffffffffffffff1681565b6100ec6101af366004610da7565b6106a8565b61011460045481565b61011460065481565b6100ec6101d4366004610bf0565b6108ae565b6101146101e7366004610cd1565b6108c2565b60025473ffffffffffffffffffffffffffffffffffffffff163314610264576002546040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff90911660248201526044015b60405180910390fd5b61026e82826108e3565b5050565b6007546102db576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f7375624944206e6f742073657400000000000000000000000000000000000000604482015260640161025b565b60005b815181101561026e57600254600754835173ffffffffffffffffffffffffffffffffffffffff9092169163bec4c08c919085908590811061032157610321610fc4565b60200260200101516040518363ffffffff1660e01b815260040161036592919091825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b600060405180830381600087803b15801561037f57600080fd5b505af1158015610393573d6000803e3d6000fd5b5050505080806103a290610f64565b9150506102de565b60088190556040805160c08101825282815260075460208083019190915260018284018190526207a1206060840152608083015282519081018352600080825260a083019190915260025492517f9b1c385e000000000000000000000000000000000000000000000000000000008152909273ffffffffffffffffffffffffffffffffffffffff1690639b1c385e90610447908490600401610e8c565b602060405180830381600087803b15801561046157600080fd5b505af1158015610475573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104999190610cea565b9392505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610521576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015260640161025b565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60005473ffffffffffffffffffffffffffffffffffffffff1633148015906105dd575060025473ffffffffffffffffffffffffffffffffffffffff163314155b15610661573361060260005473ffffffffffffffffffffffffffffffffffffffff1690565b6002546040517f061db9c100000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9384166004820152918316602483015291909116604482015260640161025b565b600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b6007546107e057600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a21a23e46040518163ffffffff1660e01b8152600401602060405180830381600087803b15801561071957600080fd5b505af115801561072d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107519190610cea565b60078190556002546040517fbec4c08c000000000000000000000000000000000000000000000000000000008152600481019290925230602483015273ffffffffffffffffffffffffffffffffffffffff169063bec4c08c90604401600060405180830381600087803b1580156107c757600080fd5b505af11580156107db573d6000803e3d6000fd5b505050505b6005546002546007546040805160208082019390935281518082039093018352808201918290527f4000aea00000000000000000000000000000000000000000000000000000000090915273ffffffffffffffffffffffffffffffffffffffff93841693634000aea09361085c93911691869190604401610e40565b602060405180830381600087803b15801561087657600080fd5b505af115801561088a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061026e9190610caf565b6108b66109ee565b6108bf81610a71565b50565b600381815481106108d257600080fd5b600091825260209091200154905081565b5a60065580516108fa906003906020840190610b67565b5060048281556040805160c0810182526008548152600754602080830191909152600182840181905262030d4060608401526080830152825190810183526000815260a082015260025491517f9b1c385e000000000000000000000000000000000000000000000000000000008152909273ffffffffffffffffffffffffffffffffffffffff90921691639b1c385e9161099691859101610e8c565b602060405180830381600087803b1580156109b057600080fd5b505af11580156109c4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906109e89190610cea565b50505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610a6f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161025b565b565b73ffffffffffffffffffffffffffffffffffffffff8116331415610af1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161025b565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b828054828255906000526020600020908101928215610ba2579160200282015b82811115610ba2578251825591602001919060010190610b87565b50610bae929150610bb2565b5090565b5b80821115610bae5760008155600101610bb3565b803573ffffffffffffffffffffffffffffffffffffffff81168114610beb57600080fd5b919050565b600060208284031215610c0257600080fd5b61049982610bc7565b60006020808385031215610c1e57600080fd5b823567ffffffffffffffff811115610c3557600080fd5b8301601f81018513610c4657600080fd5b8035610c59610c5482610f40565b610ef1565b80828252848201915084840188868560051b8701011115610c7957600080fd5b600094505b83851015610ca357610c8f81610bc7565b835260019490940193918501918501610c7e565b50979650505050505050565b600060208284031215610cc157600080fd5b8151801515811461049957600080fd5b600060208284031215610ce357600080fd5b5035919050565b600060208284031215610cfc57600080fd5b5051919050565b60008060408385031215610d1657600080fd5b8235915060208084013567ffffffffffffffff811115610d3557600080fd5b8401601f81018613610d4657600080fd5b8035610d54610c5482610f40565b80828252848201915084840189868560051b8701011115610d7457600080fd5b600094505b83851015610d97578035835260019490940193918501918501610d79565b5080955050505050509250929050565b600060208284031215610db957600080fd5b81356bffffffffffffffffffffffff8116811461049957600080fd5b6000815180845260005b81811015610dfb57602081850181015186830182015201610ddf565b81811115610e0d576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b73ffffffffffffffffffffffffffffffffffffffff841681526bffffffffffffffffffffffff83166020820152606060408201526000610e836060830184610dd5565b95945050505050565b60208152815160208201526020820151604082015261ffff60408301511660608201526000606083015163ffffffff80821660808501528060808601511660a0850152505060a083015160c080840152610ee960e0840182610dd5565b949350505050565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610f3857610f38610ff3565b604052919050565b600067ffffffffffffffff821115610f5a57610f5a610ff3565b5060051b60200190565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415610fbd577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var VRFMaliciousConsumerV2PlusABI = VRFMaliciousConsumerV2PlusMetaData.ABI

var VRFMaliciousConsumerV2PlusBin = VRFMaliciousConsumerV2PlusMetaData.Bin

func DeployVRFMaliciousConsumerV2Plus(auth *bind.TransactOpts, backend bind.ContractBackend, vrfCoordinator common.Address, link common.Address) (common.Address, *types.Transaction, *VRFMaliciousConsumerV2Plus, error) {
	parsed, err := VRFMaliciousConsumerV2PlusMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFMaliciousConsumerV2PlusBin), backend, vrfCoordinator, link)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFMaliciousConsumerV2Plus{VRFMaliciousConsumerV2PlusCaller: VRFMaliciousConsumerV2PlusCaller{contract: contract}, VRFMaliciousConsumerV2PlusTransactor: VRFMaliciousConsumerV2PlusTransactor{contract: contract}, VRFMaliciousConsumerV2PlusFilterer: VRFMaliciousConsumerV2PlusFilterer{contract: contract}}, nil
}

type VRFMaliciousConsumerV2Plus struct {
	address common.Address
	abi     abi.ABI
	VRFMaliciousConsumerV2PlusCaller
	VRFMaliciousConsumerV2PlusTransactor
	VRFMaliciousConsumerV2PlusFilterer
}

type VRFMaliciousConsumerV2PlusCaller struct {
	contract *bind.BoundContract
}

type VRFMaliciousConsumerV2PlusTransactor struct {
	contract *bind.BoundContract
}

type VRFMaliciousConsumerV2PlusFilterer struct {
	contract *bind.BoundContract
}

type VRFMaliciousConsumerV2PlusSession struct {
	Contract     *VRFMaliciousConsumerV2Plus
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFMaliciousConsumerV2PlusCallerSession struct {
	Contract *VRFMaliciousConsumerV2PlusCaller
	CallOpts bind.CallOpts
}

type VRFMaliciousConsumerV2PlusTransactorSession struct {
	Contract     *VRFMaliciousConsumerV2PlusTransactor
	TransactOpts bind.TransactOpts
}

type VRFMaliciousConsumerV2PlusRaw struct {
	Contract *VRFMaliciousConsumerV2Plus
}

type VRFMaliciousConsumerV2PlusCallerRaw struct {
	Contract *VRFMaliciousConsumerV2PlusCaller
}

type VRFMaliciousConsumerV2PlusTransactorRaw struct {
	Contract *VRFMaliciousConsumerV2PlusTransactor
}

func NewVRFMaliciousConsumerV2Plus(address common.Address, backend bind.ContractBackend) (*VRFMaliciousConsumerV2Plus, error) {
	abi, err := abi.JSON(strings.NewReader(VRFMaliciousConsumerV2PlusABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFMaliciousConsumerV2Plus(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFMaliciousConsumerV2Plus{address: address, abi: abi, VRFMaliciousConsumerV2PlusCaller: VRFMaliciousConsumerV2PlusCaller{contract: contract}, VRFMaliciousConsumerV2PlusTransactor: VRFMaliciousConsumerV2PlusTransactor{contract: contract}, VRFMaliciousConsumerV2PlusFilterer: VRFMaliciousConsumerV2PlusFilterer{contract: contract}}, nil
}

func NewVRFMaliciousConsumerV2PlusCaller(address common.Address, caller bind.ContractCaller) (*VRFMaliciousConsumerV2PlusCaller, error) {
	contract, err := bindVRFMaliciousConsumerV2Plus(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFMaliciousConsumerV2PlusCaller{contract: contract}, nil
}

func NewVRFMaliciousConsumerV2PlusTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFMaliciousConsumerV2PlusTransactor, error) {
	contract, err := bindVRFMaliciousConsumerV2Plus(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFMaliciousConsumerV2PlusTransactor{contract: contract}, nil
}

func NewVRFMaliciousConsumerV2PlusFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFMaliciousConsumerV2PlusFilterer, error) {
	contract, err := bindVRFMaliciousConsumerV2Plus(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFMaliciousConsumerV2PlusFilterer{contract: contract}, nil
}

func bindVRFMaliciousConsumerV2Plus(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFMaliciousConsumerV2PlusMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFMaliciousConsumerV2Plus.Contract.VRFMaliciousConsumerV2PlusCaller.contract.Call(opts, result, method, params...)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.VRFMaliciousConsumerV2PlusTransactor.contract.Transfer(opts)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.VRFMaliciousConsumerV2PlusTransactor.contract.Transact(opts, method, params...)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFMaliciousConsumerV2Plus.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.contract.Transfer(opts)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.contract.Transact(opts, method, params...)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFMaliciousConsumerV2Plus.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusSession) Owner() (common.Address, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.Owner(&_VRFMaliciousConsumerV2Plus.CallOpts)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusCallerSession) Owner() (common.Address, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.Owner(&_VRFMaliciousConsumerV2Plus.CallOpts)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusCaller) SGasAvailable(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFMaliciousConsumerV2Plus.contract.Call(opts, &out, "s_gasAvailable")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusSession) SGasAvailable() (*big.Int, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.SGasAvailable(&_VRFMaliciousConsumerV2Plus.CallOpts)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusCallerSession) SGasAvailable() (*big.Int, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.SGasAvailable(&_VRFMaliciousConsumerV2Plus.CallOpts)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusCaller) SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VRFMaliciousConsumerV2Plus.contract.Call(opts, &out, "s_randomWords", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.SRandomWords(&_VRFMaliciousConsumerV2Plus.CallOpts, arg0)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusCallerSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.SRandomWords(&_VRFMaliciousConsumerV2Plus.CallOpts, arg0)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusCaller) SRequestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFMaliciousConsumerV2Plus.contract.Call(opts, &out, "s_requestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusSession) SRequestId() (*big.Int, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.SRequestId(&_VRFMaliciousConsumerV2Plus.CallOpts)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusCallerSession) SRequestId() (*big.Int, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.SRequestId(&_VRFMaliciousConsumerV2Plus.CallOpts)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusCaller) SVrfCoordinator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFMaliciousConsumerV2Plus.contract.Call(opts, &out, "s_vrfCoordinator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusSession) SVrfCoordinator() (common.Address, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.SVrfCoordinator(&_VRFMaliciousConsumerV2Plus.CallOpts)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusCallerSession) SVrfCoordinator() (common.Address, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.SVrfCoordinator(&_VRFMaliciousConsumerV2Plus.CallOpts)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.contract.Transact(opts, "acceptOwnership")
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.AcceptOwnership(&_VRFMaliciousConsumerV2Plus.TransactOpts)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.AcceptOwnership(&_VRFMaliciousConsumerV2Plus.TransactOpts)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusTransactor) CreateSubscriptionAndFund(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.contract.Transact(opts, "createSubscriptionAndFund", amount)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusSession) CreateSubscriptionAndFund(amount *big.Int) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.CreateSubscriptionAndFund(&_VRFMaliciousConsumerV2Plus.TransactOpts, amount)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusTransactorSession) CreateSubscriptionAndFund(amount *big.Int) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.CreateSubscriptionAndFund(&_VRFMaliciousConsumerV2Plus.TransactOpts, amount)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.contract.Transact(opts, "rawFulfillRandomWords", requestId, randomWords)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.RawFulfillRandomWords(&_VRFMaliciousConsumerV2Plus.TransactOpts, requestId, randomWords)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusTransactorSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.RawFulfillRandomWords(&_VRFMaliciousConsumerV2Plus.TransactOpts, requestId, randomWords)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusTransactor) RequestRandomness(opts *bind.TransactOpts, keyHash [32]byte) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.contract.Transact(opts, "requestRandomness", keyHash)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusSession) RequestRandomness(keyHash [32]byte) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.RequestRandomness(&_VRFMaliciousConsumerV2Plus.TransactOpts, keyHash)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusTransactorSession) RequestRandomness(keyHash [32]byte) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.RequestRandomness(&_VRFMaliciousConsumerV2Plus.TransactOpts, keyHash)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusTransactor) SetCoordinator(opts *bind.TransactOpts, _vrfCoordinator common.Address) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.contract.Transact(opts, "setCoordinator", _vrfCoordinator)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusSession) SetCoordinator(_vrfCoordinator common.Address) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.SetCoordinator(&_VRFMaliciousConsumerV2Plus.TransactOpts, _vrfCoordinator)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusTransactorSession) SetCoordinator(_vrfCoordinator common.Address) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.SetCoordinator(&_VRFMaliciousConsumerV2Plus.TransactOpts, _vrfCoordinator)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.contract.Transact(opts, "transferOwnership", to)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.TransferOwnership(&_VRFMaliciousConsumerV2Plus.TransactOpts, to)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.TransferOwnership(&_VRFMaliciousConsumerV2Plus.TransactOpts, to)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusTransactor) UpdateSubscription(opts *bind.TransactOpts, consumers []common.Address) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.contract.Transact(opts, "updateSubscription", consumers)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusSession) UpdateSubscription(consumers []common.Address) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.UpdateSubscription(&_VRFMaliciousConsumerV2Plus.TransactOpts, consumers)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusTransactorSession) UpdateSubscription(consumers []common.Address) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.UpdateSubscription(&_VRFMaliciousConsumerV2Plus.TransactOpts, consumers)
}

type VRFMaliciousConsumerV2PlusOwnershipTransferRequestedIterator struct {
	Event *VRFMaliciousConsumerV2PlusOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFMaliciousConsumerV2PlusOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFMaliciousConsumerV2PlusOwnershipTransferRequested)
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
		it.Event = new(VRFMaliciousConsumerV2PlusOwnershipTransferRequested)
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

func (it *VRFMaliciousConsumerV2PlusOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFMaliciousConsumerV2PlusOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFMaliciousConsumerV2PlusOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFMaliciousConsumerV2PlusOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFMaliciousConsumerV2Plus.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFMaliciousConsumerV2PlusOwnershipTransferRequestedIterator{contract: _VRFMaliciousConsumerV2Plus.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFMaliciousConsumerV2PlusOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFMaliciousConsumerV2Plus.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFMaliciousConsumerV2PlusOwnershipTransferRequested)
				if err := _VRFMaliciousConsumerV2Plus.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusFilterer) ParseOwnershipTransferRequested(log types.Log) (*VRFMaliciousConsumerV2PlusOwnershipTransferRequested, error) {
	event := new(VRFMaliciousConsumerV2PlusOwnershipTransferRequested)
	if err := _VRFMaliciousConsumerV2Plus.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFMaliciousConsumerV2PlusOwnershipTransferredIterator struct {
	Event *VRFMaliciousConsumerV2PlusOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFMaliciousConsumerV2PlusOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFMaliciousConsumerV2PlusOwnershipTransferred)
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
		it.Event = new(VRFMaliciousConsumerV2PlusOwnershipTransferred)
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

func (it *VRFMaliciousConsumerV2PlusOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFMaliciousConsumerV2PlusOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFMaliciousConsumerV2PlusOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFMaliciousConsumerV2PlusOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFMaliciousConsumerV2Plus.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFMaliciousConsumerV2PlusOwnershipTransferredIterator{contract: _VRFMaliciousConsumerV2Plus.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFMaliciousConsumerV2PlusOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFMaliciousConsumerV2Plus.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFMaliciousConsumerV2PlusOwnershipTransferred)
				if err := _VRFMaliciousConsumerV2Plus.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusFilterer) ParseOwnershipTransferred(log types.Log) (*VRFMaliciousConsumerV2PlusOwnershipTransferred, error) {
	event := new(VRFMaliciousConsumerV2PlusOwnershipTransferred)
	if err := _VRFMaliciousConsumerV2Plus.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2Plus) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFMaliciousConsumerV2Plus.abi.Events["OwnershipTransferRequested"].ID:
		return _VRFMaliciousConsumerV2Plus.ParseOwnershipTransferRequested(log)
	case _VRFMaliciousConsumerV2Plus.abi.Events["OwnershipTransferred"].ID:
		return _VRFMaliciousConsumerV2Plus.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFMaliciousConsumerV2PlusOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VRFMaliciousConsumerV2PlusOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2Plus) Address() common.Address {
	return _VRFMaliciousConsumerV2Plus.address
}

type VRFMaliciousConsumerV2PlusInterface interface {
	Owner(opts *bind.CallOpts) (common.Address, error)

	SGasAvailable(opts *bind.CallOpts) (*big.Int, error)

	SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	SRequestId(opts *bind.CallOpts) (*big.Int, error)

	SVrfCoordinator(opts *bind.CallOpts) (common.Address, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	CreateSubscriptionAndFund(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	RequestRandomness(opts *bind.TransactOpts, keyHash [32]byte) (*types.Transaction, error)

	SetCoordinator(opts *bind.TransactOpts, _vrfCoordinator common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	UpdateSubscription(opts *bind.TransactOpts, consumers []common.Address) (*types.Transaction, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFMaliciousConsumerV2PlusOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFMaliciousConsumerV2PlusOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VRFMaliciousConsumerV2PlusOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFMaliciousConsumerV2PlusOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFMaliciousConsumerV2PlusOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VRFMaliciousConsumerV2PlusOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
