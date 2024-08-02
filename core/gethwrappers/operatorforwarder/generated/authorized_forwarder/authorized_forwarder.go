// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package authorized_forwarder

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

var AuthorizedForwarderMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"senders\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"changedBy\",\"type\":\"address\"}],\"name\":\"AuthorizedSendersChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"OwnershipTransferRequestedWithMessage\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"UnusedEvent\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"forward\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAuthorizedSenders\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"infiniteLoop\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"isAuthorizedSender\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkToken\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"tos\",\"type\":\"address[]\"},{\"internalType\":\"bytes[]\",\"name\":\"datas\",\"type\":\"bytes[]\"}],\"name\":\"multiForward\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"ownerForward\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"senders\",\"type\":\"address[]\"}],\"name\":\"setAuthorizedSenders\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"someNumber\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"transferOwnershipWithMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b50604051620017303803806200173083398101604081905262000034916200029d565b82826001600160a01b038216620000925760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c557620000c58162000199565b50506001600160a01b0384166200012b5760405162461bcd60e51b815260206004820152602360248201527f4c696e6b20746f6b656e2063616e6e6f742062652061207a65726f206164647260448201526265737360e81b606482015260840162000089565b6001600160a01b038085166080528216156200018f57816001600160a01b0316836001600160a01b03167f4e1e878dc28d5f040db5969163ff1acd75c44c3f655da2dde9c70bbd8e56dc7e836040516200018691906200038e565b60405180910390a35b50505050620003c3565b336001600160a01b03821603620001f35760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000089565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b03811681146200025c57600080fd5b919050565b634e487b7160e01b600052604160045260246000fd5b60005b83811015620002945781810151838201526020016200027a565b50506000910152565b60008060008060808587031215620002b457600080fd5b620002bf8562000244565b9350620002cf6020860162000244565b9250620002df6040860162000244565b60608601519092506001600160401b0380821115620002fd57600080fd5b818701915087601f8301126200031257600080fd5b81518181111562000327576200032762000261565b604051601f8201601f19908116603f0116810190838211818310171562000352576200035262000261565b816040528281528a60208487010111156200036c57600080fd5b6200037f83602083016020880162000277565b979a9699509497505050505050565b6020815260008251806020840152620003af81604085016020870162000277565b601f01601f19169190910160400192915050565b608051611343620003ed6000396000818161018b0152818161040c015261066a01526113436000f3fe608060405234801561001057600080fd5b50600436106100ea5760003560e01c806379ba50971161008c578063b64fa9e611610066578063b64fa9e614610222578063ee56997b14610235578063f2fde38b14610248578063fa00763a1461025b57600080fd5b806379ba5097146101e55780638da5cb5b146101ed578063931c9b991461020b57600080fd5b80632408afaa116100c85780632408afaa1461015e5780634d3e23231461017357806357970e93146101865780636fadcf72146101d257600080fd5b8063033f49f7146100ef578063181f5a77146101045780631dbf353d14610156575b600080fd5b6101026100fd366004610f09565b6102a4565b005b6101406040518060400160405280601981526020017f417574686f72697a6564466f7277617264657220312e312e300000000000000081525081565b60405161014d9190610f8c565b60405180910390f35b6101026102bc565b61016661031e565b60405161014d9190610ff8565b610102610181366004610f09565b61038d565b6101ad7f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161014d565b6101026101e0366004610f09565b610402565b6101026104c4565b60005473ffffffffffffffffffffffffffffffffffffffff166101ad565b61021460045481565b60405190815260200161014d565b61010261023036600461109e565b6105c1565b61010261024336600461110a565b610762565b61010261025636600461114c565b610a73565b61029461026936600461114c565b73ffffffffffffffffffffffffffffffffffffffff1660009081526002602052604090205460ff1690565b604051901515815260200161014d565b6102ac610a87565b6102b7838383610b0a565b505050565b3073ffffffffffffffffffffffffffffffffffffffff16631dbf353d6040518163ffffffff1660e01b8152600401600060405180830381600087803b15801561030457600080fd5b505af1158015610318573d6000803e3d6000fd5b50505050565b6060600380548060200260200160405190810160405280929190818152602001828054801561038357602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610358575b5050505050905090565b61039683610a73565b8273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f4e1e878dc28d5f040db5969163ff1acd75c44c3f655da2dde9c70bbd8e56dc7e84846040516103f592919061116e565b60405180910390a3505050565b61040a610c97565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16036102ac576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f43616e6e6f7420666f727761726420746f204c696e6b20746f6b656e0000000060448201526064015b60405180910390fd5b60015473ffffffffffffffffffffffffffffffffffffffff163314610545576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016104bb565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6105c9610c97565b828114610632576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f417272617973206d7573742068617665207468652073616d65206c656e67746860448201526064016104bb565b60005b8381101561075b576000858583818110610651576106516111bb565b9050602002016020810190610666919061114c565b90507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff160361071d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f43616e6e6f7420666f727761726420746f204c696e6b20746f6b656e0000000060448201526064016104bb565b61074a81858585818110610733576107336111bb565b905060200281019061074591906111ea565b610b0a565b506107548161124f565b9050610635565b5050505050565b61076a610d10565b6107d0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f43616e6e6f742073657420617574686f72697a65642073656e6465727300000060448201526064016104bb565b80610837576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f4d7573742068617665206174206c6561737420312073656e646572000000000060448201526064016104bb565b60035460005b818110156108cd576000600260006003848154811061085e5761085e6111bb565b60009182526020808320919091015473ffffffffffffffffffffffffffffffffffffffff168352820192909252604001902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169115159190911790556108c68161124f565b905061083d565b5060005b82811015610a2557600260008585848181106108ef576108ef6111bb565b9050602002016020810190610904919061114c565b73ffffffffffffffffffffffffffffffffffffffff16815260208101919091526040016000205460ff1615610995576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f4d757374206e6f742068617665206475706c69636174652073656e646572730060448201526064016104bb565b6001600260008686858181106109ad576109ad6111bb565b90506020020160208101906109c2919061114c565b73ffffffffffffffffffffffffffffffffffffffff168152602081019190915260400160002080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055610a1e8161124f565b90506108d1565b50610a3260038484610e43565b507ff263cfb3e4298332e776194610cf9fdc09ccb3ada8b9aa39764d882e11fbf0a0838333604051610a66939291906112ae565b60405180910390a1505050565b610a7b610a87565b610a8481610d4e565b50565b60005473ffffffffffffffffffffffffffffffffffffffff163314610b08576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016104bb565b565b73ffffffffffffffffffffffffffffffffffffffff83163b610b88576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f4d75737420666f727761726420746f206120636f6e747261637400000000000060448201526064016104bb565b6000808473ffffffffffffffffffffffffffffffffffffffff168484604051610bb2929190611326565b6000604051808303816000865af19150503d8060008114610bef576040519150601f19603f3d011682016040523d82523d6000602084013e610bf4565b606091505b50915091508161075b578051600003610c8f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f466f727761726465642063616c6c20726576657274656420776974686f75742060448201527f726561736f6e000000000000000000000000000000000000000000000000000060648201526084016104bb565b805181602001fd5b3360009081526002602052604090205460ff16610b08576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f4e6f7420617574686f72697a65642073656e646572000000000000000000000060448201526064016104bb565b600033610d3260005473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff1614905090565b3373ffffffffffffffffffffffffffffffffffffffff821603610dcd576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016104bb565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b828054828255906000526020600020908101928215610ebb579160200282015b82811115610ebb5781547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff843516178255602090920191600190910190610e63565b50610ec7929150610ecb565b5090565b5b80821115610ec75760008155600101610ecc565b803573ffffffffffffffffffffffffffffffffffffffff81168114610f0457600080fd5b919050565b600080600060408486031215610f1e57600080fd5b610f2784610ee0565b9250602084013567ffffffffffffffff80821115610f4457600080fd5b818601915086601f830112610f5857600080fd5b813581811115610f6757600080fd5b876020828501011115610f7957600080fd5b6020830194508093505050509250925092565b600060208083528351808285015260005b81811015610fb957858101830151858201604001528201610f9d565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b6020808252825182820181905260009190848201906040850190845b8181101561104657835173ffffffffffffffffffffffffffffffffffffffff1683529284019291840191600101611014565b50909695505050505050565b60008083601f84011261106457600080fd5b50813567ffffffffffffffff81111561107c57600080fd5b6020830191508360208260051b850101111561109757600080fd5b9250929050565b600080600080604085870312156110b457600080fd5b843567ffffffffffffffff808211156110cc57600080fd5b6110d888838901611052565b909650945060208701359150808211156110f157600080fd5b506110fe87828801611052565b95989497509550505050565b6000806020838503121561111d57600080fd5b823567ffffffffffffffff81111561113457600080fd5b61114085828601611052565b90969095509350505050565b60006020828403121561115e57600080fd5b61116782610ee0565b9392505050565b60208152816020820152818360408301376000818301604090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0160101919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261121f57600080fd5b83018035915067ffffffffffffffff82111561123a57600080fd5b60200191503681900382131561109757600080fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036112a7577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b5060010190565b6040808252810183905260008460608301825b868110156112fc5773ffffffffffffffffffffffffffffffffffffffff6112e784610ee0565b168252602092830192909101906001016112c1565b50809250505073ffffffffffffffffffffffffffffffffffffffff83166020830152949350505050565b818382376000910190815291905056fea164736f6c6343000813000a",
}

var AuthorizedForwarderABI = AuthorizedForwarderMetaData.ABI

var AuthorizedForwarderBin = AuthorizedForwarderMetaData.Bin

func DeployAuthorizedForwarder(auth *bind.TransactOpts, backend bind.ContractBackend, link common.Address, owner common.Address, recipient common.Address, message []byte) (common.Address, *types.Transaction, *AuthorizedForwarder, error) {
	parsed, err := AuthorizedForwarderMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AuthorizedForwarderBin), backend, link, owner, recipient, message)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AuthorizedForwarder{address: address, abi: *parsed, AuthorizedForwarderCaller: AuthorizedForwarderCaller{contract: contract}, AuthorizedForwarderTransactor: AuthorizedForwarderTransactor{contract: contract}, AuthorizedForwarderFilterer: AuthorizedForwarderFilterer{contract: contract}}, nil
}

type AuthorizedForwarder struct {
	address common.Address
	abi     abi.ABI
	AuthorizedForwarderCaller
	AuthorizedForwarderTransactor
	AuthorizedForwarderFilterer
}

type AuthorizedForwarderCaller struct {
	contract *bind.BoundContract
}

type AuthorizedForwarderTransactor struct {
	contract *bind.BoundContract
}

type AuthorizedForwarderFilterer struct {
	contract *bind.BoundContract
}

type AuthorizedForwarderSession struct {
	Contract     *AuthorizedForwarder
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type AuthorizedForwarderCallerSession struct {
	Contract *AuthorizedForwarderCaller
	CallOpts bind.CallOpts
}

type AuthorizedForwarderTransactorSession struct {
	Contract     *AuthorizedForwarderTransactor
	TransactOpts bind.TransactOpts
}

type AuthorizedForwarderRaw struct {
	Contract *AuthorizedForwarder
}

type AuthorizedForwarderCallerRaw struct {
	Contract *AuthorizedForwarderCaller
}

type AuthorizedForwarderTransactorRaw struct {
	Contract *AuthorizedForwarderTransactor
}

func NewAuthorizedForwarder(address common.Address, backend bind.ContractBackend) (*AuthorizedForwarder, error) {
	abi, err := abi.JSON(strings.NewReader(AuthorizedForwarderABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindAuthorizedForwarder(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AuthorizedForwarder{address: address, abi: abi, AuthorizedForwarderCaller: AuthorizedForwarderCaller{contract: contract}, AuthorizedForwarderTransactor: AuthorizedForwarderTransactor{contract: contract}, AuthorizedForwarderFilterer: AuthorizedForwarderFilterer{contract: contract}}, nil
}

func NewAuthorizedForwarderCaller(address common.Address, caller bind.ContractCaller) (*AuthorizedForwarderCaller, error) {
	contract, err := bindAuthorizedForwarder(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AuthorizedForwarderCaller{contract: contract}, nil
}

func NewAuthorizedForwarderTransactor(address common.Address, transactor bind.ContractTransactor) (*AuthorizedForwarderTransactor, error) {
	contract, err := bindAuthorizedForwarder(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AuthorizedForwarderTransactor{contract: contract}, nil
}

func NewAuthorizedForwarderFilterer(address common.Address, filterer bind.ContractFilterer) (*AuthorizedForwarderFilterer, error) {
	contract, err := bindAuthorizedForwarder(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AuthorizedForwarderFilterer{contract: contract}, nil
}

func bindAuthorizedForwarder(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AuthorizedForwarderMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_AuthorizedForwarder *AuthorizedForwarderRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AuthorizedForwarder.Contract.AuthorizedForwarderCaller.contract.Call(opts, result, method, params...)
}

func (_AuthorizedForwarder *AuthorizedForwarderRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AuthorizedForwarder.Contract.AuthorizedForwarderTransactor.contract.Transfer(opts)
}

func (_AuthorizedForwarder *AuthorizedForwarderRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AuthorizedForwarder.Contract.AuthorizedForwarderTransactor.contract.Transact(opts, method, params...)
}

func (_AuthorizedForwarder *AuthorizedForwarderCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AuthorizedForwarder.Contract.contract.Call(opts, result, method, params...)
}

func (_AuthorizedForwarder *AuthorizedForwarderTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AuthorizedForwarder.Contract.contract.Transfer(opts)
}

func (_AuthorizedForwarder *AuthorizedForwarderTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AuthorizedForwarder.Contract.contract.Transact(opts, method, params...)
}

func (_AuthorizedForwarder *AuthorizedForwarderCaller) GetAuthorizedSenders(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _AuthorizedForwarder.contract.Call(opts, &out, "getAuthorizedSenders")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_AuthorizedForwarder *AuthorizedForwarderSession) GetAuthorizedSenders() ([]common.Address, error) {
	return _AuthorizedForwarder.Contract.GetAuthorizedSenders(&_AuthorizedForwarder.CallOpts)
}

func (_AuthorizedForwarder *AuthorizedForwarderCallerSession) GetAuthorizedSenders() ([]common.Address, error) {
	return _AuthorizedForwarder.Contract.GetAuthorizedSenders(&_AuthorizedForwarder.CallOpts)
}

func (_AuthorizedForwarder *AuthorizedForwarderCaller) IsAuthorizedSender(opts *bind.CallOpts, sender common.Address) (bool, error) {
	var out []interface{}
	err := _AuthorizedForwarder.contract.Call(opts, &out, "isAuthorizedSender", sender)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_AuthorizedForwarder *AuthorizedForwarderSession) IsAuthorizedSender(sender common.Address) (bool, error) {
	return _AuthorizedForwarder.Contract.IsAuthorizedSender(&_AuthorizedForwarder.CallOpts, sender)
}

func (_AuthorizedForwarder *AuthorizedForwarderCallerSession) IsAuthorizedSender(sender common.Address) (bool, error) {
	return _AuthorizedForwarder.Contract.IsAuthorizedSender(&_AuthorizedForwarder.CallOpts, sender)
}

func (_AuthorizedForwarder *AuthorizedForwarderCaller) LinkToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AuthorizedForwarder.contract.Call(opts, &out, "linkToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AuthorizedForwarder *AuthorizedForwarderSession) LinkToken() (common.Address, error) {
	return _AuthorizedForwarder.Contract.LinkToken(&_AuthorizedForwarder.CallOpts)
}

func (_AuthorizedForwarder *AuthorizedForwarderCallerSession) LinkToken() (common.Address, error) {
	return _AuthorizedForwarder.Contract.LinkToken(&_AuthorizedForwarder.CallOpts)
}

func (_AuthorizedForwarder *AuthorizedForwarderCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AuthorizedForwarder.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AuthorizedForwarder *AuthorizedForwarderSession) Owner() (common.Address, error) {
	return _AuthorizedForwarder.Contract.Owner(&_AuthorizedForwarder.CallOpts)
}

func (_AuthorizedForwarder *AuthorizedForwarderCallerSession) Owner() (common.Address, error) {
	return _AuthorizedForwarder.Contract.Owner(&_AuthorizedForwarder.CallOpts)
}

func (_AuthorizedForwarder *AuthorizedForwarderCaller) SomeNumber(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AuthorizedForwarder.contract.Call(opts, &out, "someNumber")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AuthorizedForwarder *AuthorizedForwarderSession) SomeNumber() (*big.Int, error) {
	return _AuthorizedForwarder.Contract.SomeNumber(&_AuthorizedForwarder.CallOpts)
}

func (_AuthorizedForwarder *AuthorizedForwarderCallerSession) SomeNumber() (*big.Int, error) {
	return _AuthorizedForwarder.Contract.SomeNumber(&_AuthorizedForwarder.CallOpts)
}

func (_AuthorizedForwarder *AuthorizedForwarderCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AuthorizedForwarder.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_AuthorizedForwarder *AuthorizedForwarderSession) TypeAndVersion() (string, error) {
	return _AuthorizedForwarder.Contract.TypeAndVersion(&_AuthorizedForwarder.CallOpts)
}

func (_AuthorizedForwarder *AuthorizedForwarderCallerSession) TypeAndVersion() (string, error) {
	return _AuthorizedForwarder.Contract.TypeAndVersion(&_AuthorizedForwarder.CallOpts)
}

func (_AuthorizedForwarder *AuthorizedForwarderTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AuthorizedForwarder.contract.Transact(opts, "acceptOwnership")
}

func (_AuthorizedForwarder *AuthorizedForwarderSession) AcceptOwnership() (*types.Transaction, error) {
	return _AuthorizedForwarder.Contract.AcceptOwnership(&_AuthorizedForwarder.TransactOpts)
}

func (_AuthorizedForwarder *AuthorizedForwarderTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _AuthorizedForwarder.Contract.AcceptOwnership(&_AuthorizedForwarder.TransactOpts)
}

func (_AuthorizedForwarder *AuthorizedForwarderTransactor) Forward(opts *bind.TransactOpts, to common.Address, data []byte) (*types.Transaction, error) {
	return _AuthorizedForwarder.contract.Transact(opts, "forward", to, data)
}

func (_AuthorizedForwarder *AuthorizedForwarderSession) Forward(to common.Address, data []byte) (*types.Transaction, error) {
	return _AuthorizedForwarder.Contract.Forward(&_AuthorizedForwarder.TransactOpts, to, data)
}

func (_AuthorizedForwarder *AuthorizedForwarderTransactorSession) Forward(to common.Address, data []byte) (*types.Transaction, error) {
	return _AuthorizedForwarder.Contract.Forward(&_AuthorizedForwarder.TransactOpts, to, data)
}

func (_AuthorizedForwarder *AuthorizedForwarderTransactor) InfiniteLoop(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AuthorizedForwarder.contract.Transact(opts, "infiniteLoop")
}

func (_AuthorizedForwarder *AuthorizedForwarderSession) InfiniteLoop() (*types.Transaction, error) {
	return _AuthorizedForwarder.Contract.InfiniteLoop(&_AuthorizedForwarder.TransactOpts)
}

func (_AuthorizedForwarder *AuthorizedForwarderTransactorSession) InfiniteLoop() (*types.Transaction, error) {
	return _AuthorizedForwarder.Contract.InfiniteLoop(&_AuthorizedForwarder.TransactOpts)
}

func (_AuthorizedForwarder *AuthorizedForwarderTransactor) MultiForward(opts *bind.TransactOpts, tos []common.Address, datas [][]byte) (*types.Transaction, error) {
	return _AuthorizedForwarder.contract.Transact(opts, "multiForward", tos, datas)
}

func (_AuthorizedForwarder *AuthorizedForwarderSession) MultiForward(tos []common.Address, datas [][]byte) (*types.Transaction, error) {
	return _AuthorizedForwarder.Contract.MultiForward(&_AuthorizedForwarder.TransactOpts, tos, datas)
}

func (_AuthorizedForwarder *AuthorizedForwarderTransactorSession) MultiForward(tos []common.Address, datas [][]byte) (*types.Transaction, error) {
	return _AuthorizedForwarder.Contract.MultiForward(&_AuthorizedForwarder.TransactOpts, tos, datas)
}

func (_AuthorizedForwarder *AuthorizedForwarderTransactor) OwnerForward(opts *bind.TransactOpts, to common.Address, data []byte) (*types.Transaction, error) {
	return _AuthorizedForwarder.contract.Transact(opts, "ownerForward", to, data)
}

func (_AuthorizedForwarder *AuthorizedForwarderSession) OwnerForward(to common.Address, data []byte) (*types.Transaction, error) {
	return _AuthorizedForwarder.Contract.OwnerForward(&_AuthorizedForwarder.TransactOpts, to, data)
}

func (_AuthorizedForwarder *AuthorizedForwarderTransactorSession) OwnerForward(to common.Address, data []byte) (*types.Transaction, error) {
	return _AuthorizedForwarder.Contract.OwnerForward(&_AuthorizedForwarder.TransactOpts, to, data)
}

func (_AuthorizedForwarder *AuthorizedForwarderTransactor) SetAuthorizedSenders(opts *bind.TransactOpts, senders []common.Address) (*types.Transaction, error) {
	return _AuthorizedForwarder.contract.Transact(opts, "setAuthorizedSenders", senders)
}

func (_AuthorizedForwarder *AuthorizedForwarderSession) SetAuthorizedSenders(senders []common.Address) (*types.Transaction, error) {
	return _AuthorizedForwarder.Contract.SetAuthorizedSenders(&_AuthorizedForwarder.TransactOpts, senders)
}

func (_AuthorizedForwarder *AuthorizedForwarderTransactorSession) SetAuthorizedSenders(senders []common.Address) (*types.Transaction, error) {
	return _AuthorizedForwarder.Contract.SetAuthorizedSenders(&_AuthorizedForwarder.TransactOpts, senders)
}

func (_AuthorizedForwarder *AuthorizedForwarderTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _AuthorizedForwarder.contract.Transact(opts, "transferOwnership", to)
}

func (_AuthorizedForwarder *AuthorizedForwarderSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _AuthorizedForwarder.Contract.TransferOwnership(&_AuthorizedForwarder.TransactOpts, to)
}

func (_AuthorizedForwarder *AuthorizedForwarderTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _AuthorizedForwarder.Contract.TransferOwnership(&_AuthorizedForwarder.TransactOpts, to)
}

func (_AuthorizedForwarder *AuthorizedForwarderTransactor) TransferOwnershipWithMessage(opts *bind.TransactOpts, to common.Address, message []byte) (*types.Transaction, error) {
	return _AuthorizedForwarder.contract.Transact(opts, "transferOwnershipWithMessage", to, message)
}

func (_AuthorizedForwarder *AuthorizedForwarderSession) TransferOwnershipWithMessage(to common.Address, message []byte) (*types.Transaction, error) {
	return _AuthorizedForwarder.Contract.TransferOwnershipWithMessage(&_AuthorizedForwarder.TransactOpts, to, message)
}

func (_AuthorizedForwarder *AuthorizedForwarderTransactorSession) TransferOwnershipWithMessage(to common.Address, message []byte) (*types.Transaction, error) {
	return _AuthorizedForwarder.Contract.TransferOwnershipWithMessage(&_AuthorizedForwarder.TransactOpts, to, message)
}

type AuthorizedForwarderAuthorizedSendersChangedIterator struct {
	Event *AuthorizedForwarderAuthorizedSendersChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AuthorizedForwarderAuthorizedSendersChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AuthorizedForwarderAuthorizedSendersChanged)
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
		it.Event = new(AuthorizedForwarderAuthorizedSendersChanged)
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

func (it *AuthorizedForwarderAuthorizedSendersChangedIterator) Error() error {
	return it.fail
}

func (it *AuthorizedForwarderAuthorizedSendersChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AuthorizedForwarderAuthorizedSendersChanged struct {
	Senders   []common.Address
	ChangedBy common.Address
	Raw       types.Log
}

func (_AuthorizedForwarder *AuthorizedForwarderFilterer) FilterAuthorizedSendersChanged(opts *bind.FilterOpts) (*AuthorizedForwarderAuthorizedSendersChangedIterator, error) {

	logs, sub, err := _AuthorizedForwarder.contract.FilterLogs(opts, "AuthorizedSendersChanged")
	if err != nil {
		return nil, err
	}
	return &AuthorizedForwarderAuthorizedSendersChangedIterator{contract: _AuthorizedForwarder.contract, event: "AuthorizedSendersChanged", logs: logs, sub: sub}, nil
}

func (_AuthorizedForwarder *AuthorizedForwarderFilterer) WatchAuthorizedSendersChanged(opts *bind.WatchOpts, sink chan<- *AuthorizedForwarderAuthorizedSendersChanged) (event.Subscription, error) {

	logs, sub, err := _AuthorizedForwarder.contract.WatchLogs(opts, "AuthorizedSendersChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AuthorizedForwarderAuthorizedSendersChanged)
				if err := _AuthorizedForwarder.contract.UnpackLog(event, "AuthorizedSendersChanged", log); err != nil {
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

func (_AuthorizedForwarder *AuthorizedForwarderFilterer) ParseAuthorizedSendersChanged(log types.Log) (*AuthorizedForwarderAuthorizedSendersChanged, error) {
	event := new(AuthorizedForwarderAuthorizedSendersChanged)
	if err := _AuthorizedForwarder.contract.UnpackLog(event, "AuthorizedSendersChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AuthorizedForwarderOwnershipTransferRequestedIterator struct {
	Event *AuthorizedForwarderOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AuthorizedForwarderOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AuthorizedForwarderOwnershipTransferRequested)
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
		it.Event = new(AuthorizedForwarderOwnershipTransferRequested)
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

func (it *AuthorizedForwarderOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *AuthorizedForwarderOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AuthorizedForwarderOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AuthorizedForwarder *AuthorizedForwarderFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AuthorizedForwarderOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AuthorizedForwarder.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AuthorizedForwarderOwnershipTransferRequestedIterator{contract: _AuthorizedForwarder.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_AuthorizedForwarder *AuthorizedForwarderFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *AuthorizedForwarderOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AuthorizedForwarder.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AuthorizedForwarderOwnershipTransferRequested)
				if err := _AuthorizedForwarder.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_AuthorizedForwarder *AuthorizedForwarderFilterer) ParseOwnershipTransferRequested(log types.Log) (*AuthorizedForwarderOwnershipTransferRequested, error) {
	event := new(AuthorizedForwarderOwnershipTransferRequested)
	if err := _AuthorizedForwarder.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AuthorizedForwarderOwnershipTransferRequestedWithMessageIterator struct {
	Event *AuthorizedForwarderOwnershipTransferRequestedWithMessage

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AuthorizedForwarderOwnershipTransferRequestedWithMessageIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AuthorizedForwarderOwnershipTransferRequestedWithMessage)
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
		it.Event = new(AuthorizedForwarderOwnershipTransferRequestedWithMessage)
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

func (it *AuthorizedForwarderOwnershipTransferRequestedWithMessageIterator) Error() error {
	return it.fail
}

func (it *AuthorizedForwarderOwnershipTransferRequestedWithMessageIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AuthorizedForwarderOwnershipTransferRequestedWithMessage struct {
	From    common.Address
	To      common.Address
	Message []byte
	Raw     types.Log
}

func (_AuthorizedForwarder *AuthorizedForwarderFilterer) FilterOwnershipTransferRequestedWithMessage(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AuthorizedForwarderOwnershipTransferRequestedWithMessageIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AuthorizedForwarder.contract.FilterLogs(opts, "OwnershipTransferRequestedWithMessage", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AuthorizedForwarderOwnershipTransferRequestedWithMessageIterator{contract: _AuthorizedForwarder.contract, event: "OwnershipTransferRequestedWithMessage", logs: logs, sub: sub}, nil
}

func (_AuthorizedForwarder *AuthorizedForwarderFilterer) WatchOwnershipTransferRequestedWithMessage(opts *bind.WatchOpts, sink chan<- *AuthorizedForwarderOwnershipTransferRequestedWithMessage, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AuthorizedForwarder.contract.WatchLogs(opts, "OwnershipTransferRequestedWithMessage", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AuthorizedForwarderOwnershipTransferRequestedWithMessage)
				if err := _AuthorizedForwarder.contract.UnpackLog(event, "OwnershipTransferRequestedWithMessage", log); err != nil {
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

func (_AuthorizedForwarder *AuthorizedForwarderFilterer) ParseOwnershipTransferRequestedWithMessage(log types.Log) (*AuthorizedForwarderOwnershipTransferRequestedWithMessage, error) {
	event := new(AuthorizedForwarderOwnershipTransferRequestedWithMessage)
	if err := _AuthorizedForwarder.contract.UnpackLog(event, "OwnershipTransferRequestedWithMessage", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AuthorizedForwarderOwnershipTransferredIterator struct {
	Event *AuthorizedForwarderOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AuthorizedForwarderOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AuthorizedForwarderOwnershipTransferred)
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
		it.Event = new(AuthorizedForwarderOwnershipTransferred)
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

func (it *AuthorizedForwarderOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *AuthorizedForwarderOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AuthorizedForwarderOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AuthorizedForwarder *AuthorizedForwarderFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AuthorizedForwarderOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AuthorizedForwarder.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AuthorizedForwarderOwnershipTransferredIterator{contract: _AuthorizedForwarder.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_AuthorizedForwarder *AuthorizedForwarderFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AuthorizedForwarderOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AuthorizedForwarder.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AuthorizedForwarderOwnershipTransferred)
				if err := _AuthorizedForwarder.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_AuthorizedForwarder *AuthorizedForwarderFilterer) ParseOwnershipTransferred(log types.Log) (*AuthorizedForwarderOwnershipTransferred, error) {
	event := new(AuthorizedForwarderOwnershipTransferred)
	if err := _AuthorizedForwarder.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AuthorizedForwarderUnusedEventIterator struct {
	Event *AuthorizedForwarderUnusedEvent

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AuthorizedForwarderUnusedEventIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AuthorizedForwarderUnusedEvent)
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
		it.Event = new(AuthorizedForwarderUnusedEvent)
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

func (it *AuthorizedForwarderUnusedEventIterator) Error() error {
	return it.fail
}

func (it *AuthorizedForwarderUnusedEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AuthorizedForwarderUnusedEvent struct {
	Raw types.Log
}

func (_AuthorizedForwarder *AuthorizedForwarderFilterer) FilterUnusedEvent(opts *bind.FilterOpts) (*AuthorizedForwarderUnusedEventIterator, error) {

	logs, sub, err := _AuthorizedForwarder.contract.FilterLogs(opts, "UnusedEvent")
	if err != nil {
		return nil, err
	}
	return &AuthorizedForwarderUnusedEventIterator{contract: _AuthorizedForwarder.contract, event: "UnusedEvent", logs: logs, sub: sub}, nil
}

func (_AuthorizedForwarder *AuthorizedForwarderFilterer) WatchUnusedEvent(opts *bind.WatchOpts, sink chan<- *AuthorizedForwarderUnusedEvent) (event.Subscription, error) {

	logs, sub, err := _AuthorizedForwarder.contract.WatchLogs(opts, "UnusedEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AuthorizedForwarderUnusedEvent)
				if err := _AuthorizedForwarder.contract.UnpackLog(event, "UnusedEvent", log); err != nil {
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

func (_AuthorizedForwarder *AuthorizedForwarderFilterer) ParseUnusedEvent(log types.Log) (*AuthorizedForwarderUnusedEvent, error) {
	event := new(AuthorizedForwarderUnusedEvent)
	if err := _AuthorizedForwarder.contract.UnpackLog(event, "UnusedEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_AuthorizedForwarder *AuthorizedForwarder) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _AuthorizedForwarder.abi.Events["AuthorizedSendersChanged"].ID:
		return _AuthorizedForwarder.ParseAuthorizedSendersChanged(log)
	case _AuthorizedForwarder.abi.Events["OwnershipTransferRequested"].ID:
		return _AuthorizedForwarder.ParseOwnershipTransferRequested(log)
	case _AuthorizedForwarder.abi.Events["OwnershipTransferRequestedWithMessage"].ID:
		return _AuthorizedForwarder.ParseOwnershipTransferRequestedWithMessage(log)
	case _AuthorizedForwarder.abi.Events["OwnershipTransferred"].ID:
		return _AuthorizedForwarder.ParseOwnershipTransferred(log)
	case _AuthorizedForwarder.abi.Events["UnusedEvent"].ID:
		return _AuthorizedForwarder.ParseUnusedEvent(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (AuthorizedForwarderAuthorizedSendersChanged) Topic() common.Hash {
	return common.HexToHash("0xf263cfb3e4298332e776194610cf9fdc09ccb3ada8b9aa39764d882e11fbf0a0")
}

func (AuthorizedForwarderOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (AuthorizedForwarderOwnershipTransferRequestedWithMessage) Topic() common.Hash {
	return common.HexToHash("0x4e1e878dc28d5f040db5969163ff1acd75c44c3f655da2dde9c70bbd8e56dc7e")
}

func (AuthorizedForwarderOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (AuthorizedForwarderUnusedEvent) Topic() common.Hash {
	return common.HexToHash("0x1ae8ed1ad86196069877a857260241bcd39074e708bc22e358e74aa71f85f687")
}

func (_AuthorizedForwarder *AuthorizedForwarder) Address() common.Address {
	return _AuthorizedForwarder.address
}

type AuthorizedForwarderInterface interface {
	GetAuthorizedSenders(opts *bind.CallOpts) ([]common.Address, error)

	IsAuthorizedSender(opts *bind.CallOpts, sender common.Address) (bool, error)

	LinkToken(opts *bind.CallOpts) (common.Address, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SomeNumber(opts *bind.CallOpts) (*big.Int, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	Forward(opts *bind.TransactOpts, to common.Address, data []byte) (*types.Transaction, error)

	InfiniteLoop(opts *bind.TransactOpts) (*types.Transaction, error)

	MultiForward(opts *bind.TransactOpts, tos []common.Address, datas [][]byte) (*types.Transaction, error)

	OwnerForward(opts *bind.TransactOpts, to common.Address, data []byte) (*types.Transaction, error)

	SetAuthorizedSenders(opts *bind.TransactOpts, senders []common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	TransferOwnershipWithMessage(opts *bind.TransactOpts, to common.Address, message []byte) (*types.Transaction, error)

	FilterAuthorizedSendersChanged(opts *bind.FilterOpts) (*AuthorizedForwarderAuthorizedSendersChangedIterator, error)

	WatchAuthorizedSendersChanged(opts *bind.WatchOpts, sink chan<- *AuthorizedForwarderAuthorizedSendersChanged) (event.Subscription, error)

	ParseAuthorizedSendersChanged(log types.Log) (*AuthorizedForwarderAuthorizedSendersChanged, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AuthorizedForwarderOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *AuthorizedForwarderOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*AuthorizedForwarderOwnershipTransferRequested, error)

	FilterOwnershipTransferRequestedWithMessage(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AuthorizedForwarderOwnershipTransferRequestedWithMessageIterator, error)

	WatchOwnershipTransferRequestedWithMessage(opts *bind.WatchOpts, sink chan<- *AuthorizedForwarderOwnershipTransferRequestedWithMessage, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequestedWithMessage(log types.Log) (*AuthorizedForwarderOwnershipTransferRequestedWithMessage, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AuthorizedForwarderOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AuthorizedForwarderOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*AuthorizedForwarderOwnershipTransferred, error)

	FilterUnusedEvent(opts *bind.FilterOpts) (*AuthorizedForwarderUnusedEventIterator, error)

	WatchUnusedEvent(opts *bind.WatchOpts, sink chan<- *AuthorizedForwarderUnusedEvent) (event.Subscription, error)

	ParseUnusedEvent(log types.Log) (*AuthorizedForwarderUnusedEvent, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
