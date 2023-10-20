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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"senders\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"changedBy\",\"type\":\"address\"}],\"name\":\"AuthorizedSendersChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"OwnershipTransferRequestedWithMessage\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"forward\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAuthorizedSenders\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"isAuthorizedSender\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkToken\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"tos\",\"type\":\"address[]\"},{\"internalType\":\"bytes[]\",\"name\":\"datas\",\"type\":\"bytes[]\"}],\"name\":\"multiForward\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"ownerForward\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"senders\",\"type\":\"address[]\"}],\"name\":\"setAuthorizedSenders\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"transferOwnershipWithMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b506040516200169d3803806200169d83398101604081905262000034916200029d565b82826001600160a01b038216620000925760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c557620000c58162000199565b50506001600160a01b0384166200012b5760405162461bcd60e51b815260206004820152602360248201527f4c696e6b20746f6b656e2063616e6e6f742062652061207a65726f206164647260448201526265737360e81b606482015260840162000089565b6001600160a01b038085166080528216156200018f57816001600160a01b0316836001600160a01b03167f4e1e878dc28d5f040db5969163ff1acd75c44c3f655da2dde9c70bbd8e56dc7e836040516200018691906200038e565b60405180910390a35b50505050620003c3565b336001600160a01b03821603620001f35760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000089565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b03811681146200025c57600080fd5b919050565b634e487b7160e01b600052604160045260246000fd5b60005b83811015620002945781810151838201526020016200027a565b50506000910152565b60008060008060808587031215620002b457600080fd5b620002bf8562000244565b9350620002cf6020860162000244565b9250620002df6040860162000244565b60608601519092506001600160401b0380821115620002fd57600080fd5b818701915087601f8301126200031257600080fd5b81518181111562000327576200032762000261565b604051601f8201601f19908116603f0116810190838211818310171562000352576200035262000261565b816040528281528a60208487010111156200036c57600080fd5b6200037f83602083016020880162000277565b979a9699509497505050505050565b6020815260008251806020840152620003af81604085016020870162000277565b601f01601f19169190910160400192915050565b6080516112b0620003ed6000396000818161016d0152818161037501526105d301526112b06000f3fe608060405234801561001057600080fd5b50600436106100d45760003560e01c806379ba509711610081578063ee56997b1161005b578063ee56997b14610200578063f2fde38b14610213578063fa00763a1461022657600080fd5b806379ba5097146101c75780638da5cb5b146101cf578063b64fa9e6146101ed57600080fd5b80634d3e2323116100b25780634d3e23231461015557806357970e93146101685780636fadcf72146101b457600080fd5b8063033f49f7146100d9578063181f5a77146100ee5780632408afaa14610140575b600080fd5b6100ec6100e7366004610e76565b61026f565b005b61012a6040518060400160405280601981526020017f417574686f72697a6564466f7277617264657220312e302e300000000000000081525081565b6040516101379190610ef9565b60405180910390f35b610148610287565b6040516101379190610f65565b6100ec610163366004610e76565b6102f6565b61018f7f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610137565b6100ec6101c2366004610e76565b61036b565b6100ec61042d565b60005473ffffffffffffffffffffffffffffffffffffffff1661018f565b6100ec6101fb36600461100b565b61052a565b6100ec61020e366004611077565b6106cb565b6100ec6102213660046110b9565b6109e0565b61025f6102343660046110b9565b73ffffffffffffffffffffffffffffffffffffffff1660009081526002602052604090205460ff1690565b6040519015158152602001610137565b6102776109f4565b610282838383610a77565b505050565b606060038054806020026020016040519081016040528092919081815260200182805480156102ec57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116102c1575b5050505050905090565b6102ff836109e0565b8273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f4e1e878dc28d5f040db5969163ff1acd75c44c3f655da2dde9c70bbd8e56dc7e848460405161035e9291906110db565b60405180910390a3505050565b610373610c04565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610277576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f43616e6e6f7420666f727761726420746f204c696e6b20746f6b656e0000000060448201526064015b60405180910390fd5b60015473ffffffffffffffffffffffffffffffffffffffff1633146104ae576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610424565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610532610c04565b82811461059b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f417272617973206d7573742068617665207468652073616d65206c656e6774686044820152606401610424565b60005b838110156106c45760008585838181106105ba576105ba611128565b90506020020160208101906105cf91906110b9565b90507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610686576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f43616e6e6f7420666f727761726420746f204c696e6b20746f6b656e000000006044820152606401610424565b6106b38185858581811061069c5761069c611128565b90506020028101906106ae9190611157565b610a77565b506106bd816111bc565b905061059e565b5050505050565b6106d3610c7d565b610739576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f43616e6e6f742073657420617574686f72697a65642073656e646572730000006044820152606401610424565b806107a0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f4d7573742068617665206174206c6561737420312073656e64657200000000006044820152606401610424565b60035460005b8181101561083857600060026000600384815481106107c7576107c7611128565b60009182526020808320919091015473ffffffffffffffffffffffffffffffffffffffff168352820192909252604001902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001691151591909117905580610830816111bc565b9150506107a6565b5060005b82811015610992576002600085858481811061085a5761085a611128565b905060200201602081019061086f91906110b9565b73ffffffffffffffffffffffffffffffffffffffff16815260208101919091526040016000205460ff1615610900576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f4d757374206e6f742068617665206475706c69636174652073656e64657273006044820152606401610424565b60016002600086868581811061091857610918611128565b905060200201602081019061092d91906110b9565b73ffffffffffffffffffffffffffffffffffffffff168152602081019190915260400160002080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169115159190911790558061098a816111bc565b91505061083c565b5061099f60038484610db0565b507ff263cfb3e4298332e776194610cf9fdc09ccb3ada8b9aa39764d882e11fbf0a08383336040516109d39392919061121b565b60405180910390a1505050565b6109e86109f4565b6109f181610cbb565b50565b60005473ffffffffffffffffffffffffffffffffffffffff163314610a75576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610424565b565b73ffffffffffffffffffffffffffffffffffffffff83163b610af5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f4d75737420666f727761726420746f206120636f6e74726163740000000000006044820152606401610424565b6000808473ffffffffffffffffffffffffffffffffffffffff168484604051610b1f929190611293565b6000604051808303816000865af19150503d8060008114610b5c576040519150601f19603f3d011682016040523d82523d6000602084013e610b61565b606091505b5091509150816106c4578051600003610bfc576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f466f727761726465642063616c6c20726576657274656420776974686f75742060448201527f726561736f6e00000000000000000000000000000000000000000000000000006064820152608401610424565b805181602001fd5b3360009081526002602052604090205460ff16610a75576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f4e6f7420617574686f72697a65642073656e64657200000000000000000000006044820152606401610424565b600033610c9f60005473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff1614905090565b3373ffffffffffffffffffffffffffffffffffffffff821603610d3a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610424565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b828054828255906000526020600020908101928215610e28579160200282015b82811115610e285781547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff843516178255602090920191600190910190610dd0565b50610e34929150610e38565b5090565b5b80821115610e345760008155600101610e39565b803573ffffffffffffffffffffffffffffffffffffffff81168114610e7157600080fd5b919050565b600080600060408486031215610e8b57600080fd5b610e9484610e4d565b9250602084013567ffffffffffffffff80821115610eb157600080fd5b818601915086601f830112610ec557600080fd5b813581811115610ed457600080fd5b876020828501011115610ee657600080fd5b6020830194508093505050509250925092565b600060208083528351808285015260005b81811015610f2657858101830151858201604001528201610f0a565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b6020808252825182820181905260009190848201906040850190845b81811015610fb357835173ffffffffffffffffffffffffffffffffffffffff1683529284019291840191600101610f81565b50909695505050505050565b60008083601f840112610fd157600080fd5b50813567ffffffffffffffff811115610fe957600080fd5b6020830191508360208260051b850101111561100457600080fd5b9250929050565b6000806000806040858703121561102157600080fd5b843567ffffffffffffffff8082111561103957600080fd5b61104588838901610fbf565b9096509450602087013591508082111561105e57600080fd5b5061106b87828801610fbf565b95989497509550505050565b6000806020838503121561108a57600080fd5b823567ffffffffffffffff8111156110a157600080fd5b6110ad85828601610fbf565b90969095509350505050565b6000602082840312156110cb57600080fd5b6110d482610e4d565b9392505050565b60208152816020820152818360408301376000818301604090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0160101919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261118c57600080fd5b83018035915067ffffffffffffffff8211156111a757600080fd5b60200191503681900382131561100457600080fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203611214577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b5060010190565b6040808252810183905260008460608301825b868110156112695773ffffffffffffffffffffffffffffffffffffffff61125484610e4d565b1682526020928301929091019060010161122e565b50809250505073ffffffffffffffffffffffffffffffffffffffff83166020830152949350505050565b818382376000910190815291905056fea164736f6c6343000813000a",
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
	return address, tx, &AuthorizedForwarder{AuthorizedForwarderCaller: AuthorizedForwarderCaller{contract: contract}, AuthorizedForwarderTransactor: AuthorizedForwarderTransactor{contract: contract}, AuthorizedForwarderFilterer: AuthorizedForwarderFilterer{contract: contract}}, nil
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

func (_AuthorizedForwarder *AuthorizedForwarder) Address() common.Address {
	return _AuthorizedForwarder.address
}

type AuthorizedForwarderInterface interface {
	GetAuthorizedSenders(opts *bind.CallOpts) ([]common.Address, error)

	IsAuthorizedSender(opts *bind.CallOpts, sender common.Address) (bool, error)

	LinkToken(opts *bind.CallOpts) (common.Address, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	Forward(opts *bind.TransactOpts, to common.Address, data []byte) (*types.Transaction, error)

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

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
