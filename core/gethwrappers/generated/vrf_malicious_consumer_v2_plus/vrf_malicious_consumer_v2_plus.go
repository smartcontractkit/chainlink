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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlyCoordinatorCanFulfill\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"}],\"name\":\"OnlyOwnerOrCoordinator\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"createSubscriptionAndFund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSubId\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"requestRandomness\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_gasAvailable\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_randomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"_subId\",\"type\":\"uint64\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"name\":\"updateSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b50604051620013f4380380620013f48339810160408190526200003491620001cc565b8133806000816200008c5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000bf57620000bf8162000103565b5050600280546001600160a01b03199081166001600160a01b0394851617909155600580548216958416959095179094555060068054909316911617905562000204565b6001600160a01b0381163314156200015e5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000083565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b0381168114620001c757600080fd5b919050565b60008060408385031215620001e057600080fd5b620001eb83620001af565b9150620001fb60208401620001af565b90509250929050565b6111e080620002146000396000f3fe608060405234801561001057600080fd5b50600436106100d45760003560e01c8063cf62c8ab11610081578063f2fde38b1161005b578063f2fde38b1461018f578063f6eaffc8146101a2578063fbc6ba6f146101b557600080fd5b8063cf62c8ab1461016a578063e89e106a1461017d578063f08c5daa1461018657600080fd5b80635e3b709f116100b25780635e3b709f1461011457806379ba50971461013a5780638da5cb5b1461014257600080fd5b80631fe543e3146100d95780632d6d99f3146100ee57806336bfffed14610101575b600080fd5b6100ec6100e7366004610e77565b6101f3565b005b6100ec6100fc366004610d48565b610279565b6100ec61010f366004610d7f565b6103a9565b610127610122366004610e45565b610531565b6040519081526020015b60405180910390f35b6100ec610649565b60005460405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610131565b6100ec610178366004610f38565b610746565b61012760045481565b61012760075481565b6100ec61019d366004610d2d565b6109c9565b6101276101b0366004610e45565b6109dd565b60025474010000000000000000000000000000000000000000900467ffffffffffffffff1660405167ffffffffffffffff9091168152602001610131565b60025473ffffffffffffffffffffffffffffffffffffffff16331461026b576002546040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff90911660248201526044015b60405180910390fd5b61027582826109fe565b5050565b60005473ffffffffffffffffffffffffffffffffffffffff1633148015906102b9575060025473ffffffffffffffffffffffffffffffffffffffff163314155b1561033d57336102de60005473ffffffffffffffffffffffffffffffffffffffff1690565b6002546040517f061db9c100000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff93841660048201529183166024830152919091166044820152606401610262565b6002805467ffffffffffffffff90921674010000000000000000000000000000000000000000027fffffffff0000000000000000000000000000000000000000000000000000000090921673ffffffffffffffffffffffffffffffffffffffff90931692909217179055565b60025474010000000000000000000000000000000000000000900467ffffffffffffffff16610434576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f7375624944206e6f7420736574000000000000000000000000000000000000006044820152606401610262565b60005b815181101561027557600554600254835173ffffffffffffffffffffffffffffffffffffffff90921691637341c10c9174010000000000000000000000000000000000000000900467ffffffffffffffff169085908590811061049c5761049c61115f565b60200260200101516040518363ffffffff1660e01b81526004016104ec92919067ffffffffffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b600060405180830381600087803b15801561050657600080fd5b505af115801561051a573d6000803e3d6000fd5b505050508080610529906110ff565b915050610437565b60088190556040805160c08101825282815260025474010000000000000000000000000000000000000000900467ffffffffffffffff1660208083019190915260018284018190526207a1206060840152608083015282519081018352600080825260a083019190915260055492517f596b8b88000000000000000000000000000000000000000000000000000000008152909273ffffffffffffffffffffffffffffffffffffffff169063596b8b88906105f090849060040161101d565b602060405180830381600087803b15801561060a57600080fd5b505af115801561061e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106429190610e5e565b9392505050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146106ca576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610262565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60025474010000000000000000000000000000000000000000900467ffffffffffffffff166108f157600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a21a23e46040518163ffffffff1660e01b8152600401602060405180830381600087803b1580156107d957600080fd5b505af11580156107ed573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108119190610f1b565b600280547fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000067ffffffffffffffff938416810291909117918290556005546040517f7341c10c00000000000000000000000000000000000000000000000000000000815291909204909216600483015230602483015273ffffffffffffffffffffffffffffffffffffffff1690637341c10c90604401600060405180830381600087803b1580156108d857600080fd5b505af11580156108ec573d6000803e3d6000fd5b505050505b600654600554600254604080517401000000000000000000000000000000000000000090920467ffffffffffffffff16602083015273ffffffffffffffffffffffffffffffffffffffff93841693634000aea09316918591016040516020818303038152906040526040518463ffffffff1660e01b815260040161097793929190610fd1565b602060405180830381600087803b15801561099157600080fd5b505af11580156109a5573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102759190610e23565b6109d1610b2b565b6109da81610bae565b50565b600381815481106109ed57600080fd5b600091825260209091200154905081565b5a6007558051610a15906003906020840190610ca4565b5060048281556040805160c081018252600854815260025474010000000000000000000000000000000000000000900467ffffffffffffffff16602080830191909152600182840181905262030d4060608401526080830152825190810183526000815260a082015260055491517f596b8b88000000000000000000000000000000000000000000000000000000008152909273ffffffffffffffffffffffffffffffffffffffff9092169163596b8b8891610ad39185910161101d565b602060405180830381600087803b158015610aed57600080fd5b505af1158015610b01573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b259190610e5e565b50505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610bac576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610262565b565b73ffffffffffffffffffffffffffffffffffffffff8116331415610c2e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610262565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b828054828255906000526020600020908101928215610cdf579160200282015b82811115610cdf578251825591602001919060010190610cc4565b50610ceb929150610cef565b5090565b5b80821115610ceb5760008155600101610cf0565b803573ffffffffffffffffffffffffffffffffffffffff81168114610d2857600080fd5b919050565b600060208284031215610d3f57600080fd5b61064282610d04565b60008060408385031215610d5b57600080fd5b610d6483610d04565b91506020830135610d74816111bd565b809150509250929050565b60006020808385031215610d9257600080fd5b823567ffffffffffffffff811115610da957600080fd5b8301601f81018513610dba57600080fd5b8035610dcd610dc8826110db565b61108c565b80828252848201915084840188868560051b8701011115610ded57600080fd5b600094505b83851015610e1757610e0381610d04565b835260019490940193918501918501610df2565b50979650505050505050565b600060208284031215610e3557600080fd5b8151801515811461064257600080fd5b600060208284031215610e5757600080fd5b5035919050565b600060208284031215610e7057600080fd5b5051919050565b60008060408385031215610e8a57600080fd5b8235915060208084013567ffffffffffffffff811115610ea957600080fd5b8401601f81018613610eba57600080fd5b8035610ec8610dc8826110db565b80828252848201915084840189868560051b8701011115610ee857600080fd5b600094505b83851015610f0b578035835260019490940193918501918501610eed565b5080955050505050509250929050565b600060208284031215610f2d57600080fd5b8151610642816111bd565b600060208284031215610f4a57600080fd5b81356bffffffffffffffffffffffff8116811461064257600080fd5b6000815180845260005b81811015610f8c57602081850181015186830182015201610f70565b81811115610f9e576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b73ffffffffffffffffffffffffffffffffffffffff841681526bffffffffffffffffffffffff831660208201526060604082015260006110146060830184610f66565b95945050505050565b602081528151602082015267ffffffffffffffff602083015116604082015261ffff60408301511660608201526000606083015163ffffffff80821660808501528060808601511660a0850152505060a083015160c08084015261108460e0840182610f66565b949350505050565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156110d3576110d361118e565b604052919050565b600067ffffffffffffffff8211156110f5576110f561118e565b5060051b60200190565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415611158577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b67ffffffffffffffff811681146109da57600080fdfea164736f6c6343000806000a",
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

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusCaller) GetSubId(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _VRFMaliciousConsumerV2Plus.contract.Call(opts, &out, "getSubId")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusSession) GetSubId() (uint64, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.GetSubId(&_VRFMaliciousConsumerV2Plus.CallOpts)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusCallerSession) GetSubId() (uint64, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.GetSubId(&_VRFMaliciousConsumerV2Plus.CallOpts)
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

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusTransactor) SetConfig(opts *bind.TransactOpts, _vrfCoordinator common.Address, _subId uint64) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.contract.Transact(opts, "setConfig", _vrfCoordinator, _subId)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusSession) SetConfig(_vrfCoordinator common.Address, _subId uint64) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.SetConfig(&_VRFMaliciousConsumerV2Plus.TransactOpts, _vrfCoordinator, _subId)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusTransactorSession) SetConfig(_vrfCoordinator common.Address, _subId uint64) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.SetConfig(&_VRFMaliciousConsumerV2Plus.TransactOpts, _vrfCoordinator, _subId)
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
	GetSubId(opts *bind.CallOpts) (uint64, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SGasAvailable(opts *bind.CallOpts) (*big.Int, error)

	SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	SRequestId(opts *bind.CallOpts) (*big.Int, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	CreateSubscriptionAndFund(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	RequestRandomness(opts *bind.TransactOpts, keyHash [32]byte) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, _vrfCoordinator common.Address, _subId uint64) (*types.Transaction, error)

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
