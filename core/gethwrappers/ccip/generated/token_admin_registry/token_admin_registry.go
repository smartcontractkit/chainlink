// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package token_admin_registry

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

type TokenAdminRegistryTokenConfig struct {
	Administrator        common.Address
	PendingAdministrator common.Address
	TokenPool            common.Address
}

var TokenAdminRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"AlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"InvalidTokenPoolToken\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"OnlyAdministrator\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"OnlyPendingAdministrator\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"OnlyRegistryModuleOrOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"currentAdmin\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdministratorTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdministratorTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousPool\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newPool\",\"type\":\"address\"}],\"name\":\"PoolSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"module\",\"type\":\"address\"}],\"name\":\"RegistryModuleAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"module\",\"type\":\"address\"}],\"name\":\"RegistryModuleRemoved\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"}],\"name\":\"acceptAdminRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"module\",\"type\":\"address\"}],\"name\":\"addRegistryModule\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"startIndex\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"maxCount\",\"type\":\"uint64\"}],\"name\":\"getAllConfiguredTokens\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"tokens\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getPool\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"tokens\",\"type\":\"address[]\"}],\"name\":\"getPools\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getTokenConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"administrator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"pendingAdministrator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenPool\",\"type\":\"address\"}],\"internalType\":\"structTokenAdminRegistry.TokenConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"administrator\",\"type\":\"address\"}],\"name\":\"isAdministrator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"module\",\"type\":\"address\"}],\"name\":\"isRegistryModule\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"administrator\",\"type\":\"address\"}],\"name\":\"proposeAdministrator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"module\",\"type\":\"address\"}],\"name\":\"removeRegistryModule\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"}],\"name\":\"setPool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"transferAdminRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5033806000816100675760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615610097576100978161009f565b505050610148565b336001600160a01b038216036100f75760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161005e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b611449806101576000396000f3fe608060405234801561001057600080fd5b50600436106101005760003560e01c80637d3f255211610097578063cb67e3b111610066578063cb67e3b1146102bc578063ddadfa8e14610374578063e677ae3714610387578063f2fde38b1461039a57600080fd5b80637d3f2552146101e05780638da5cb5b14610203578063bbe4f6db14610242578063c1af6e031461027f57600080fd5b80634e847fc7116100d35780634e847fc7146101925780635e63547a146101a557806372d64a81146101c557806379ba5097146101d857600080fd5b806310cbcf1814610105578063156194da1461011a578063181f5a771461012d5780633dc457721461017f575b600080fd5b61011861011336600461116c565b6103ad565b005b61011861012836600461116c565b61040a565b6101696040518060400160405280601881526020017f546f6b656e41646d696e526567697374727920312e352e30000000000000000081525081565b6040516101769190611187565b60405180910390f35b61011861018d36600461116c565b61050f565b6101186101a03660046111f4565b610573565b6101b86101b3366004611227565b6107d3565b604051610176919061129c565b6101b86101d336600461130e565b6108cc565b6101186109e2565b6101f36101ee36600461116c565b610adf565b6040519015158152602001610176565b60005473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610176565b61021d61025036600461116c565b73ffffffffffffffffffffffffffffffffffffffff908116600090815260026020819052604090912001541690565b6101f361028d3660046111f4565b73ffffffffffffffffffffffffffffffffffffffff918216600090815260026020526040902054821691161490565b6103356102ca36600461116c565b60408051606080820183526000808352602080840182905292840181905273ffffffffffffffffffffffffffffffffffffffff948516815260028084529084902084519283018552805486168352600181015486169383019390935291909101549092169082015290565b60408051825173ffffffffffffffffffffffffffffffffffffffff90811682526020808501518216908301529282015190921690820152606001610176565b6101186103823660046111f4565b610aec565b6101186103953660046111f4565b610bf6565b6101186103a836600461116c565b610dbe565b6103b5610dcf565b6103c0600582610e52565b156104075760405173ffffffffffffffffffffffffffffffffffffffff8216907f93eaa26dcb9275e56bacb1d33fdbf402262da6f0f4baf2a6e2cd154b73f387f890600090a25b50565b73ffffffffffffffffffffffffffffffffffffffff808216600090815260026020526040902060018101549091163314610493576040517f3edffe7500000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff831660248201526044015b60405180910390fd5b8054337fffffffffffffffffffffffff00000000000000000000000000000000000000009182168117835560018301805490921690915560405173ffffffffffffffffffffffffffffffffffffffff8416907f399b55200f7f639a63d76efe3dcfa9156ce367058d6b673041b84a628885f5a790600090a35050565b610517610dcf565b610522600582610e7b565b156104075760405173ffffffffffffffffffffffffffffffffffffffff821681527f3cabf004338366bfeaeb610ad827cb58d16b588017c509501f2c97c83caae7b29060200160405180910390a150565b73ffffffffffffffffffffffffffffffffffffffff80831660009081526002602052604090205483911633146105f3576040517fed5d85b500000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff8216602482015260440161048a565b73ffffffffffffffffffffffffffffffffffffffff8216158015906106a557506040517f240028e800000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff848116600483015283169063240028e890602401602060405180830381865afa15801561067f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106a39190611338565b155b156106f4576040517f962b60e600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8416600482015260240161048a565b73ffffffffffffffffffffffffffffffffffffffff808416600090815260026020819052604090912090810180548584167fffffffffffffffffffffffff0000000000000000000000000000000000000000821681179092559192919091169081146107cc578373ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff168673ffffffffffffffffffffffffffffffffffffffff167f754449ec3aff3bd528bfce43ae9319c4a381b67fcd1d20097b3b24dacaecc35d60405160405180910390a45b5050505050565b606060008267ffffffffffffffff8111156107f0576107f061135a565b604051908082528060200260200182016040528015610819578160200160208202803683370190505b50905060005b838110156108c2576002600086868481811061083d5761083d611389565b9050602002016020810190610852919061116c565b73ffffffffffffffffffffffffffffffffffffffff9081168252602082019290925260400160002060020154835191169083908390811061089557610895611389565b73ffffffffffffffffffffffffffffffffffffffff9092166020928302919091019091015260010161081f565b5090505b92915050565b606060006108da6003610e9d565b9050808467ffffffffffffffff16106108f357506108c6565b67ffffffffffffffff80841690829061090e908716836113e7565b111561092b5761092867ffffffffffffffff8616836113fa565b90505b8067ffffffffffffffff8111156109445761094461135a565b60405190808252806020026020018201604052801561096d578160200160208202803683370190505b50925060005b818110156109d95761099a6109928267ffffffffffffffff89166113e7565b600390610ea7565b8482815181106109ac576109ac611389565b73ffffffffffffffffffffffffffffffffffffffff90921660209283029190910190910152600101610973565b50505092915050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610a63576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015260640161048a565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60006108c6600583610eb3565b73ffffffffffffffffffffffffffffffffffffffff8083166000908152600260205260409020548391163314610b6c576040517fed5d85b500000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff8216602482015260440161048a565b73ffffffffffffffffffffffffffffffffffffffff8381166000818152600260205260408082206001810180547fffffffffffffffffffffffff00000000000000000000000000000000000000001695881695861790559051909392339290917fc54c3051ff16e63bb9203214432372aca006c589e3653619b577a3265675b7169190a450505050565b610bff33610adf565b158015610c24575060005473ffffffffffffffffffffffffffffffffffffffff163314155b15610c5d576040517f51ca1ec300000000000000000000000000000000000000000000000000000000815233600482015260240161048a565b73ffffffffffffffffffffffffffffffffffffffff8116610caa576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8083166000908152600260205260409020805490911615610d24576040517f45ed80e900000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8416600482015260240161048a565b6001810180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8416179055610d71600384610e7b565b5060405173ffffffffffffffffffffffffffffffffffffffff808416916000918616907fc54c3051ff16e63bb9203214432372aca006c589e3653619b577a3265675b716908390a4505050565b610dc6610dcf565b61040781610ee2565b60005473ffffffffffffffffffffffffffffffffffffffff163314610e50576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161048a565b565b6000610e748373ffffffffffffffffffffffffffffffffffffffff8416610fd7565b9392505050565b6000610e748373ffffffffffffffffffffffffffffffffffffffff84166110ca565b60006108c6825490565b6000610e748383611119565b73ffffffffffffffffffffffffffffffffffffffff811660009081526001830160205260408120541515610e74565b3373ffffffffffffffffffffffffffffffffffffffff821603610f61576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161048a565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600081815260018301602052604081205480156110c0576000610ffb6001836113fa565b855490915060009061100f906001906113fa565b905080821461107457600086600001828154811061102f5761102f611389565b906000526020600020015490508087600001848154811061105257611052611389565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806110855761108561140d565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506108c6565b60009150506108c6565b6000818152600183016020526040812054611111575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556108c6565b5060006108c6565b600082600001828154811061113057611130611389565b9060005260206000200154905092915050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461116757600080fd5b919050565b60006020828403121561117e57600080fd5b610e7482611143565b60006020808352835180602085015260005b818110156111b557858101830151858201604001528201611199565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b6000806040838503121561120757600080fd5b61121083611143565b915061121e60208401611143565b90509250929050565b6000806020838503121561123a57600080fd5b823567ffffffffffffffff8082111561125257600080fd5b818501915085601f83011261126657600080fd5b81358181111561127557600080fd5b8660208260051b850101111561128a57600080fd5b60209290920196919550909350505050565b6020808252825182820181905260009190848201906040850190845b818110156112ea57835173ffffffffffffffffffffffffffffffffffffffff16835292840192918401916001016112b8565b50909695505050505050565b803567ffffffffffffffff8116811461116757600080fd5b6000806040838503121561132157600080fd5b61132a836112f6565b915061121e602084016112f6565b60006020828403121561134a57600080fd5b81518015158114610e7457600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b808201808211156108c6576108c66113b8565b818103818111156108c6576108c66113b8565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea164736f6c6343000818000a",
}

var TokenAdminRegistryABI = TokenAdminRegistryMetaData.ABI

var TokenAdminRegistryBin = TokenAdminRegistryMetaData.Bin

func DeployTokenAdminRegistry(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TokenAdminRegistry, error) {
	parsed, err := TokenAdminRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TokenAdminRegistryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TokenAdminRegistry{address: address, abi: *parsed, TokenAdminRegistryCaller: TokenAdminRegistryCaller{contract: contract}, TokenAdminRegistryTransactor: TokenAdminRegistryTransactor{contract: contract}, TokenAdminRegistryFilterer: TokenAdminRegistryFilterer{contract: contract}}, nil
}

type TokenAdminRegistry struct {
	address common.Address
	abi     abi.ABI
	TokenAdminRegistryCaller
	TokenAdminRegistryTransactor
	TokenAdminRegistryFilterer
}

type TokenAdminRegistryCaller struct {
	contract *bind.BoundContract
}

type TokenAdminRegistryTransactor struct {
	contract *bind.BoundContract
}

type TokenAdminRegistryFilterer struct {
	contract *bind.BoundContract
}

type TokenAdminRegistrySession struct {
	Contract     *TokenAdminRegistry
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type TokenAdminRegistryCallerSession struct {
	Contract *TokenAdminRegistryCaller
	CallOpts bind.CallOpts
}

type TokenAdminRegistryTransactorSession struct {
	Contract     *TokenAdminRegistryTransactor
	TransactOpts bind.TransactOpts
}

type TokenAdminRegistryRaw struct {
	Contract *TokenAdminRegistry
}

type TokenAdminRegistryCallerRaw struct {
	Contract *TokenAdminRegistryCaller
}

type TokenAdminRegistryTransactorRaw struct {
	Contract *TokenAdminRegistryTransactor
}

func NewTokenAdminRegistry(address common.Address, backend bind.ContractBackend) (*TokenAdminRegistry, error) {
	abi, err := abi.JSON(strings.NewReader(TokenAdminRegistryABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindTokenAdminRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TokenAdminRegistry{address: address, abi: abi, TokenAdminRegistryCaller: TokenAdminRegistryCaller{contract: contract}, TokenAdminRegistryTransactor: TokenAdminRegistryTransactor{contract: contract}, TokenAdminRegistryFilterer: TokenAdminRegistryFilterer{contract: contract}}, nil
}

func NewTokenAdminRegistryCaller(address common.Address, caller bind.ContractCaller) (*TokenAdminRegistryCaller, error) {
	contract, err := bindTokenAdminRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TokenAdminRegistryCaller{contract: contract}, nil
}

func NewTokenAdminRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*TokenAdminRegistryTransactor, error) {
	contract, err := bindTokenAdminRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TokenAdminRegistryTransactor{contract: contract}, nil
}

func NewTokenAdminRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*TokenAdminRegistryFilterer, error) {
	contract, err := bindTokenAdminRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TokenAdminRegistryFilterer{contract: contract}, nil
}

func bindTokenAdminRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TokenAdminRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_TokenAdminRegistry *TokenAdminRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TokenAdminRegistry.Contract.TokenAdminRegistryCaller.contract.Call(opts, result, method, params...)
}

func (_TokenAdminRegistry *TokenAdminRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenAdminRegistry.Contract.TokenAdminRegistryTransactor.contract.Transfer(opts)
}

func (_TokenAdminRegistry *TokenAdminRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenAdminRegistry.Contract.TokenAdminRegistryTransactor.contract.Transact(opts, method, params...)
}

func (_TokenAdminRegistry *TokenAdminRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TokenAdminRegistry.Contract.contract.Call(opts, result, method, params...)
}

func (_TokenAdminRegistry *TokenAdminRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenAdminRegistry.Contract.contract.Transfer(opts)
}

func (_TokenAdminRegistry *TokenAdminRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenAdminRegistry.Contract.contract.Transact(opts, method, params...)
}

func (_TokenAdminRegistry *TokenAdminRegistryCaller) GetAllConfiguredTokens(opts *bind.CallOpts, startIndex uint64, maxCount uint64) ([]common.Address, error) {
	var out []interface{}
	err := _TokenAdminRegistry.contract.Call(opts, &out, "getAllConfiguredTokens", startIndex, maxCount)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_TokenAdminRegistry *TokenAdminRegistrySession) GetAllConfiguredTokens(startIndex uint64, maxCount uint64) ([]common.Address, error) {
	return _TokenAdminRegistry.Contract.GetAllConfiguredTokens(&_TokenAdminRegistry.CallOpts, startIndex, maxCount)
}

func (_TokenAdminRegistry *TokenAdminRegistryCallerSession) GetAllConfiguredTokens(startIndex uint64, maxCount uint64) ([]common.Address, error) {
	return _TokenAdminRegistry.Contract.GetAllConfiguredTokens(&_TokenAdminRegistry.CallOpts, startIndex, maxCount)
}

func (_TokenAdminRegistry *TokenAdminRegistryCaller) GetPool(opts *bind.CallOpts, token common.Address) (common.Address, error) {
	var out []interface{}
	err := _TokenAdminRegistry.contract.Call(opts, &out, "getPool", token)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_TokenAdminRegistry *TokenAdminRegistrySession) GetPool(token common.Address) (common.Address, error) {
	return _TokenAdminRegistry.Contract.GetPool(&_TokenAdminRegistry.CallOpts, token)
}

func (_TokenAdminRegistry *TokenAdminRegistryCallerSession) GetPool(token common.Address) (common.Address, error) {
	return _TokenAdminRegistry.Contract.GetPool(&_TokenAdminRegistry.CallOpts, token)
}

func (_TokenAdminRegistry *TokenAdminRegistryCaller) GetPools(opts *bind.CallOpts, tokens []common.Address) ([]common.Address, error) {
	var out []interface{}
	err := _TokenAdminRegistry.contract.Call(opts, &out, "getPools", tokens)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_TokenAdminRegistry *TokenAdminRegistrySession) GetPools(tokens []common.Address) ([]common.Address, error) {
	return _TokenAdminRegistry.Contract.GetPools(&_TokenAdminRegistry.CallOpts, tokens)
}

func (_TokenAdminRegistry *TokenAdminRegistryCallerSession) GetPools(tokens []common.Address) ([]common.Address, error) {
	return _TokenAdminRegistry.Contract.GetPools(&_TokenAdminRegistry.CallOpts, tokens)
}

func (_TokenAdminRegistry *TokenAdminRegistryCaller) GetTokenConfig(opts *bind.CallOpts, token common.Address) (TokenAdminRegistryTokenConfig, error) {
	var out []interface{}
	err := _TokenAdminRegistry.contract.Call(opts, &out, "getTokenConfig", token)

	if err != nil {
		return *new(TokenAdminRegistryTokenConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(TokenAdminRegistryTokenConfig)).(*TokenAdminRegistryTokenConfig)

	return out0, err

}

func (_TokenAdminRegistry *TokenAdminRegistrySession) GetTokenConfig(token common.Address) (TokenAdminRegistryTokenConfig, error) {
	return _TokenAdminRegistry.Contract.GetTokenConfig(&_TokenAdminRegistry.CallOpts, token)
}

func (_TokenAdminRegistry *TokenAdminRegistryCallerSession) GetTokenConfig(token common.Address) (TokenAdminRegistryTokenConfig, error) {
	return _TokenAdminRegistry.Contract.GetTokenConfig(&_TokenAdminRegistry.CallOpts, token)
}

func (_TokenAdminRegistry *TokenAdminRegistryCaller) IsAdministrator(opts *bind.CallOpts, localToken common.Address, administrator common.Address) (bool, error) {
	var out []interface{}
	err := _TokenAdminRegistry.contract.Call(opts, &out, "isAdministrator", localToken, administrator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_TokenAdminRegistry *TokenAdminRegistrySession) IsAdministrator(localToken common.Address, administrator common.Address) (bool, error) {
	return _TokenAdminRegistry.Contract.IsAdministrator(&_TokenAdminRegistry.CallOpts, localToken, administrator)
}

func (_TokenAdminRegistry *TokenAdminRegistryCallerSession) IsAdministrator(localToken common.Address, administrator common.Address) (bool, error) {
	return _TokenAdminRegistry.Contract.IsAdministrator(&_TokenAdminRegistry.CallOpts, localToken, administrator)
}

func (_TokenAdminRegistry *TokenAdminRegistryCaller) IsRegistryModule(opts *bind.CallOpts, module common.Address) (bool, error) {
	var out []interface{}
	err := _TokenAdminRegistry.contract.Call(opts, &out, "isRegistryModule", module)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_TokenAdminRegistry *TokenAdminRegistrySession) IsRegistryModule(module common.Address) (bool, error) {
	return _TokenAdminRegistry.Contract.IsRegistryModule(&_TokenAdminRegistry.CallOpts, module)
}

func (_TokenAdminRegistry *TokenAdminRegistryCallerSession) IsRegistryModule(module common.Address) (bool, error) {
	return _TokenAdminRegistry.Contract.IsRegistryModule(&_TokenAdminRegistry.CallOpts, module)
}

func (_TokenAdminRegistry *TokenAdminRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TokenAdminRegistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_TokenAdminRegistry *TokenAdminRegistrySession) Owner() (common.Address, error) {
	return _TokenAdminRegistry.Contract.Owner(&_TokenAdminRegistry.CallOpts)
}

func (_TokenAdminRegistry *TokenAdminRegistryCallerSession) Owner() (common.Address, error) {
	return _TokenAdminRegistry.Contract.Owner(&_TokenAdminRegistry.CallOpts)
}

func (_TokenAdminRegistry *TokenAdminRegistryCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TokenAdminRegistry.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_TokenAdminRegistry *TokenAdminRegistrySession) TypeAndVersion() (string, error) {
	return _TokenAdminRegistry.Contract.TypeAndVersion(&_TokenAdminRegistry.CallOpts)
}

func (_TokenAdminRegistry *TokenAdminRegistryCallerSession) TypeAndVersion() (string, error) {
	return _TokenAdminRegistry.Contract.TypeAndVersion(&_TokenAdminRegistry.CallOpts)
}

func (_TokenAdminRegistry *TokenAdminRegistryTransactor) AcceptAdminRole(opts *bind.TransactOpts, localToken common.Address) (*types.Transaction, error) {
	return _TokenAdminRegistry.contract.Transact(opts, "acceptAdminRole", localToken)
}

func (_TokenAdminRegistry *TokenAdminRegistrySession) AcceptAdminRole(localToken common.Address) (*types.Transaction, error) {
	return _TokenAdminRegistry.Contract.AcceptAdminRole(&_TokenAdminRegistry.TransactOpts, localToken)
}

func (_TokenAdminRegistry *TokenAdminRegistryTransactorSession) AcceptAdminRole(localToken common.Address) (*types.Transaction, error) {
	return _TokenAdminRegistry.Contract.AcceptAdminRole(&_TokenAdminRegistry.TransactOpts, localToken)
}

func (_TokenAdminRegistry *TokenAdminRegistryTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenAdminRegistry.contract.Transact(opts, "acceptOwnership")
}

func (_TokenAdminRegistry *TokenAdminRegistrySession) AcceptOwnership() (*types.Transaction, error) {
	return _TokenAdminRegistry.Contract.AcceptOwnership(&_TokenAdminRegistry.TransactOpts)
}

func (_TokenAdminRegistry *TokenAdminRegistryTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _TokenAdminRegistry.Contract.AcceptOwnership(&_TokenAdminRegistry.TransactOpts)
}

func (_TokenAdminRegistry *TokenAdminRegistryTransactor) AddRegistryModule(opts *bind.TransactOpts, module common.Address) (*types.Transaction, error) {
	return _TokenAdminRegistry.contract.Transact(opts, "addRegistryModule", module)
}

func (_TokenAdminRegistry *TokenAdminRegistrySession) AddRegistryModule(module common.Address) (*types.Transaction, error) {
	return _TokenAdminRegistry.Contract.AddRegistryModule(&_TokenAdminRegistry.TransactOpts, module)
}

func (_TokenAdminRegistry *TokenAdminRegistryTransactorSession) AddRegistryModule(module common.Address) (*types.Transaction, error) {
	return _TokenAdminRegistry.Contract.AddRegistryModule(&_TokenAdminRegistry.TransactOpts, module)
}

func (_TokenAdminRegistry *TokenAdminRegistryTransactor) ProposeAdministrator(opts *bind.TransactOpts, localToken common.Address, administrator common.Address) (*types.Transaction, error) {
	return _TokenAdminRegistry.contract.Transact(opts, "proposeAdministrator", localToken, administrator)
}

func (_TokenAdminRegistry *TokenAdminRegistrySession) ProposeAdministrator(localToken common.Address, administrator common.Address) (*types.Transaction, error) {
	return _TokenAdminRegistry.Contract.ProposeAdministrator(&_TokenAdminRegistry.TransactOpts, localToken, administrator)
}

func (_TokenAdminRegistry *TokenAdminRegistryTransactorSession) ProposeAdministrator(localToken common.Address, administrator common.Address) (*types.Transaction, error) {
	return _TokenAdminRegistry.Contract.ProposeAdministrator(&_TokenAdminRegistry.TransactOpts, localToken, administrator)
}

func (_TokenAdminRegistry *TokenAdminRegistryTransactor) RemoveRegistryModule(opts *bind.TransactOpts, module common.Address) (*types.Transaction, error) {
	return _TokenAdminRegistry.contract.Transact(opts, "removeRegistryModule", module)
}

func (_TokenAdminRegistry *TokenAdminRegistrySession) RemoveRegistryModule(module common.Address) (*types.Transaction, error) {
	return _TokenAdminRegistry.Contract.RemoveRegistryModule(&_TokenAdminRegistry.TransactOpts, module)
}

func (_TokenAdminRegistry *TokenAdminRegistryTransactorSession) RemoveRegistryModule(module common.Address) (*types.Transaction, error) {
	return _TokenAdminRegistry.Contract.RemoveRegistryModule(&_TokenAdminRegistry.TransactOpts, module)
}

func (_TokenAdminRegistry *TokenAdminRegistryTransactor) SetPool(opts *bind.TransactOpts, localToken common.Address, pool common.Address) (*types.Transaction, error) {
	return _TokenAdminRegistry.contract.Transact(opts, "setPool", localToken, pool)
}

func (_TokenAdminRegistry *TokenAdminRegistrySession) SetPool(localToken common.Address, pool common.Address) (*types.Transaction, error) {
	return _TokenAdminRegistry.Contract.SetPool(&_TokenAdminRegistry.TransactOpts, localToken, pool)
}

func (_TokenAdminRegistry *TokenAdminRegistryTransactorSession) SetPool(localToken common.Address, pool common.Address) (*types.Transaction, error) {
	return _TokenAdminRegistry.Contract.SetPool(&_TokenAdminRegistry.TransactOpts, localToken, pool)
}

func (_TokenAdminRegistry *TokenAdminRegistryTransactor) TransferAdminRole(opts *bind.TransactOpts, localToken common.Address, newAdmin common.Address) (*types.Transaction, error) {
	return _TokenAdminRegistry.contract.Transact(opts, "transferAdminRole", localToken, newAdmin)
}

func (_TokenAdminRegistry *TokenAdminRegistrySession) TransferAdminRole(localToken common.Address, newAdmin common.Address) (*types.Transaction, error) {
	return _TokenAdminRegistry.Contract.TransferAdminRole(&_TokenAdminRegistry.TransactOpts, localToken, newAdmin)
}

func (_TokenAdminRegistry *TokenAdminRegistryTransactorSession) TransferAdminRole(localToken common.Address, newAdmin common.Address) (*types.Transaction, error) {
	return _TokenAdminRegistry.Contract.TransferAdminRole(&_TokenAdminRegistry.TransactOpts, localToken, newAdmin)
}

func (_TokenAdminRegistry *TokenAdminRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _TokenAdminRegistry.contract.Transact(opts, "transferOwnership", to)
}

func (_TokenAdminRegistry *TokenAdminRegistrySession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _TokenAdminRegistry.Contract.TransferOwnership(&_TokenAdminRegistry.TransactOpts, to)
}

func (_TokenAdminRegistry *TokenAdminRegistryTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _TokenAdminRegistry.Contract.TransferOwnership(&_TokenAdminRegistry.TransactOpts, to)
}

type TokenAdminRegistryAdministratorTransferRequestedIterator struct {
	Event *TokenAdminRegistryAdministratorTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TokenAdminRegistryAdministratorTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenAdminRegistryAdministratorTransferRequested)
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
		it.Event = new(TokenAdminRegistryAdministratorTransferRequested)
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

func (it *TokenAdminRegistryAdministratorTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *TokenAdminRegistryAdministratorTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TokenAdminRegistryAdministratorTransferRequested struct {
	Token        common.Address
	CurrentAdmin common.Address
	NewAdmin     common.Address
	Raw          types.Log
}

func (_TokenAdminRegistry *TokenAdminRegistryFilterer) FilterAdministratorTransferRequested(opts *bind.FilterOpts, token []common.Address, currentAdmin []common.Address, newAdmin []common.Address) (*TokenAdminRegistryAdministratorTransferRequestedIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var currentAdminRule []interface{}
	for _, currentAdminItem := range currentAdmin {
		currentAdminRule = append(currentAdminRule, currentAdminItem)
	}
	var newAdminRule []interface{}
	for _, newAdminItem := range newAdmin {
		newAdminRule = append(newAdminRule, newAdminItem)
	}

	logs, sub, err := _TokenAdminRegistry.contract.FilterLogs(opts, "AdministratorTransferRequested", tokenRule, currentAdminRule, newAdminRule)
	if err != nil {
		return nil, err
	}
	return &TokenAdminRegistryAdministratorTransferRequestedIterator{contract: _TokenAdminRegistry.contract, event: "AdministratorTransferRequested", logs: logs, sub: sub}, nil
}

func (_TokenAdminRegistry *TokenAdminRegistryFilterer) WatchAdministratorTransferRequested(opts *bind.WatchOpts, sink chan<- *TokenAdminRegistryAdministratorTransferRequested, token []common.Address, currentAdmin []common.Address, newAdmin []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var currentAdminRule []interface{}
	for _, currentAdminItem := range currentAdmin {
		currentAdminRule = append(currentAdminRule, currentAdminItem)
	}
	var newAdminRule []interface{}
	for _, newAdminItem := range newAdmin {
		newAdminRule = append(newAdminRule, newAdminItem)
	}

	logs, sub, err := _TokenAdminRegistry.contract.WatchLogs(opts, "AdministratorTransferRequested", tokenRule, currentAdminRule, newAdminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TokenAdminRegistryAdministratorTransferRequested)
				if err := _TokenAdminRegistry.contract.UnpackLog(event, "AdministratorTransferRequested", log); err != nil {
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

func (_TokenAdminRegistry *TokenAdminRegistryFilterer) ParseAdministratorTransferRequested(log types.Log) (*TokenAdminRegistryAdministratorTransferRequested, error) {
	event := new(TokenAdminRegistryAdministratorTransferRequested)
	if err := _TokenAdminRegistry.contract.UnpackLog(event, "AdministratorTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TokenAdminRegistryAdministratorTransferredIterator struct {
	Event *TokenAdminRegistryAdministratorTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TokenAdminRegistryAdministratorTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenAdminRegistryAdministratorTransferred)
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
		it.Event = new(TokenAdminRegistryAdministratorTransferred)
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

func (it *TokenAdminRegistryAdministratorTransferredIterator) Error() error {
	return it.fail
}

func (it *TokenAdminRegistryAdministratorTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TokenAdminRegistryAdministratorTransferred struct {
	Token    common.Address
	NewAdmin common.Address
	Raw      types.Log
}

func (_TokenAdminRegistry *TokenAdminRegistryFilterer) FilterAdministratorTransferred(opts *bind.FilterOpts, token []common.Address, newAdmin []common.Address) (*TokenAdminRegistryAdministratorTransferredIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var newAdminRule []interface{}
	for _, newAdminItem := range newAdmin {
		newAdminRule = append(newAdminRule, newAdminItem)
	}

	logs, sub, err := _TokenAdminRegistry.contract.FilterLogs(opts, "AdministratorTransferred", tokenRule, newAdminRule)
	if err != nil {
		return nil, err
	}
	return &TokenAdminRegistryAdministratorTransferredIterator{contract: _TokenAdminRegistry.contract, event: "AdministratorTransferred", logs: logs, sub: sub}, nil
}

func (_TokenAdminRegistry *TokenAdminRegistryFilterer) WatchAdministratorTransferred(opts *bind.WatchOpts, sink chan<- *TokenAdminRegistryAdministratorTransferred, token []common.Address, newAdmin []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var newAdminRule []interface{}
	for _, newAdminItem := range newAdmin {
		newAdminRule = append(newAdminRule, newAdminItem)
	}

	logs, sub, err := _TokenAdminRegistry.contract.WatchLogs(opts, "AdministratorTransferred", tokenRule, newAdminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TokenAdminRegistryAdministratorTransferred)
				if err := _TokenAdminRegistry.contract.UnpackLog(event, "AdministratorTransferred", log); err != nil {
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

func (_TokenAdminRegistry *TokenAdminRegistryFilterer) ParseAdministratorTransferred(log types.Log) (*TokenAdminRegistryAdministratorTransferred, error) {
	event := new(TokenAdminRegistryAdministratorTransferred)
	if err := _TokenAdminRegistry.contract.UnpackLog(event, "AdministratorTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TokenAdminRegistryOwnershipTransferRequestedIterator struct {
	Event *TokenAdminRegistryOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TokenAdminRegistryOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenAdminRegistryOwnershipTransferRequested)
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
		it.Event = new(TokenAdminRegistryOwnershipTransferRequested)
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

func (it *TokenAdminRegistryOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *TokenAdminRegistryOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TokenAdminRegistryOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_TokenAdminRegistry *TokenAdminRegistryFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TokenAdminRegistryOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TokenAdminRegistry.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &TokenAdminRegistryOwnershipTransferRequestedIterator{contract: _TokenAdminRegistry.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_TokenAdminRegistry *TokenAdminRegistryFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *TokenAdminRegistryOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TokenAdminRegistry.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TokenAdminRegistryOwnershipTransferRequested)
				if err := _TokenAdminRegistry.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_TokenAdminRegistry *TokenAdminRegistryFilterer) ParseOwnershipTransferRequested(log types.Log) (*TokenAdminRegistryOwnershipTransferRequested, error) {
	event := new(TokenAdminRegistryOwnershipTransferRequested)
	if err := _TokenAdminRegistry.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TokenAdminRegistryOwnershipTransferredIterator struct {
	Event *TokenAdminRegistryOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TokenAdminRegistryOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenAdminRegistryOwnershipTransferred)
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
		it.Event = new(TokenAdminRegistryOwnershipTransferred)
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

func (it *TokenAdminRegistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *TokenAdminRegistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TokenAdminRegistryOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_TokenAdminRegistry *TokenAdminRegistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TokenAdminRegistryOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TokenAdminRegistry.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &TokenAdminRegistryOwnershipTransferredIterator{contract: _TokenAdminRegistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_TokenAdminRegistry *TokenAdminRegistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TokenAdminRegistryOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TokenAdminRegistry.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TokenAdminRegistryOwnershipTransferred)
				if err := _TokenAdminRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_TokenAdminRegistry *TokenAdminRegistryFilterer) ParseOwnershipTransferred(log types.Log) (*TokenAdminRegistryOwnershipTransferred, error) {
	event := new(TokenAdminRegistryOwnershipTransferred)
	if err := _TokenAdminRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TokenAdminRegistryPoolSetIterator struct {
	Event *TokenAdminRegistryPoolSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TokenAdminRegistryPoolSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenAdminRegistryPoolSet)
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
		it.Event = new(TokenAdminRegistryPoolSet)
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

func (it *TokenAdminRegistryPoolSetIterator) Error() error {
	return it.fail
}

func (it *TokenAdminRegistryPoolSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TokenAdminRegistryPoolSet struct {
	Token        common.Address
	PreviousPool common.Address
	NewPool      common.Address
	Raw          types.Log
}

func (_TokenAdminRegistry *TokenAdminRegistryFilterer) FilterPoolSet(opts *bind.FilterOpts, token []common.Address, previousPool []common.Address, newPool []common.Address) (*TokenAdminRegistryPoolSetIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var previousPoolRule []interface{}
	for _, previousPoolItem := range previousPool {
		previousPoolRule = append(previousPoolRule, previousPoolItem)
	}
	var newPoolRule []interface{}
	for _, newPoolItem := range newPool {
		newPoolRule = append(newPoolRule, newPoolItem)
	}

	logs, sub, err := _TokenAdminRegistry.contract.FilterLogs(opts, "PoolSet", tokenRule, previousPoolRule, newPoolRule)
	if err != nil {
		return nil, err
	}
	return &TokenAdminRegistryPoolSetIterator{contract: _TokenAdminRegistry.contract, event: "PoolSet", logs: logs, sub: sub}, nil
}

func (_TokenAdminRegistry *TokenAdminRegistryFilterer) WatchPoolSet(opts *bind.WatchOpts, sink chan<- *TokenAdminRegistryPoolSet, token []common.Address, previousPool []common.Address, newPool []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var previousPoolRule []interface{}
	for _, previousPoolItem := range previousPool {
		previousPoolRule = append(previousPoolRule, previousPoolItem)
	}
	var newPoolRule []interface{}
	for _, newPoolItem := range newPool {
		newPoolRule = append(newPoolRule, newPoolItem)
	}

	logs, sub, err := _TokenAdminRegistry.contract.WatchLogs(opts, "PoolSet", tokenRule, previousPoolRule, newPoolRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TokenAdminRegistryPoolSet)
				if err := _TokenAdminRegistry.contract.UnpackLog(event, "PoolSet", log); err != nil {
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

func (_TokenAdminRegistry *TokenAdminRegistryFilterer) ParsePoolSet(log types.Log) (*TokenAdminRegistryPoolSet, error) {
	event := new(TokenAdminRegistryPoolSet)
	if err := _TokenAdminRegistry.contract.UnpackLog(event, "PoolSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TokenAdminRegistryRegistryModuleAddedIterator struct {
	Event *TokenAdminRegistryRegistryModuleAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TokenAdminRegistryRegistryModuleAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenAdminRegistryRegistryModuleAdded)
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
		it.Event = new(TokenAdminRegistryRegistryModuleAdded)
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

func (it *TokenAdminRegistryRegistryModuleAddedIterator) Error() error {
	return it.fail
}

func (it *TokenAdminRegistryRegistryModuleAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TokenAdminRegistryRegistryModuleAdded struct {
	Module common.Address
	Raw    types.Log
}

func (_TokenAdminRegistry *TokenAdminRegistryFilterer) FilterRegistryModuleAdded(opts *bind.FilterOpts) (*TokenAdminRegistryRegistryModuleAddedIterator, error) {

	logs, sub, err := _TokenAdminRegistry.contract.FilterLogs(opts, "RegistryModuleAdded")
	if err != nil {
		return nil, err
	}
	return &TokenAdminRegistryRegistryModuleAddedIterator{contract: _TokenAdminRegistry.contract, event: "RegistryModuleAdded", logs: logs, sub: sub}, nil
}

func (_TokenAdminRegistry *TokenAdminRegistryFilterer) WatchRegistryModuleAdded(opts *bind.WatchOpts, sink chan<- *TokenAdminRegistryRegistryModuleAdded) (event.Subscription, error) {

	logs, sub, err := _TokenAdminRegistry.contract.WatchLogs(opts, "RegistryModuleAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TokenAdminRegistryRegistryModuleAdded)
				if err := _TokenAdminRegistry.contract.UnpackLog(event, "RegistryModuleAdded", log); err != nil {
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

func (_TokenAdminRegistry *TokenAdminRegistryFilterer) ParseRegistryModuleAdded(log types.Log) (*TokenAdminRegistryRegistryModuleAdded, error) {
	event := new(TokenAdminRegistryRegistryModuleAdded)
	if err := _TokenAdminRegistry.contract.UnpackLog(event, "RegistryModuleAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TokenAdminRegistryRegistryModuleRemovedIterator struct {
	Event *TokenAdminRegistryRegistryModuleRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TokenAdminRegistryRegistryModuleRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenAdminRegistryRegistryModuleRemoved)
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
		it.Event = new(TokenAdminRegistryRegistryModuleRemoved)
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

func (it *TokenAdminRegistryRegistryModuleRemovedIterator) Error() error {
	return it.fail
}

func (it *TokenAdminRegistryRegistryModuleRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TokenAdminRegistryRegistryModuleRemoved struct {
	Module common.Address
	Raw    types.Log
}

func (_TokenAdminRegistry *TokenAdminRegistryFilterer) FilterRegistryModuleRemoved(opts *bind.FilterOpts, module []common.Address) (*TokenAdminRegistryRegistryModuleRemovedIterator, error) {

	var moduleRule []interface{}
	for _, moduleItem := range module {
		moduleRule = append(moduleRule, moduleItem)
	}

	logs, sub, err := _TokenAdminRegistry.contract.FilterLogs(opts, "RegistryModuleRemoved", moduleRule)
	if err != nil {
		return nil, err
	}
	return &TokenAdminRegistryRegistryModuleRemovedIterator{contract: _TokenAdminRegistry.contract, event: "RegistryModuleRemoved", logs: logs, sub: sub}, nil
}

func (_TokenAdminRegistry *TokenAdminRegistryFilterer) WatchRegistryModuleRemoved(opts *bind.WatchOpts, sink chan<- *TokenAdminRegistryRegistryModuleRemoved, module []common.Address) (event.Subscription, error) {

	var moduleRule []interface{}
	for _, moduleItem := range module {
		moduleRule = append(moduleRule, moduleItem)
	}

	logs, sub, err := _TokenAdminRegistry.contract.WatchLogs(opts, "RegistryModuleRemoved", moduleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TokenAdminRegistryRegistryModuleRemoved)
				if err := _TokenAdminRegistry.contract.UnpackLog(event, "RegistryModuleRemoved", log); err != nil {
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

func (_TokenAdminRegistry *TokenAdminRegistryFilterer) ParseRegistryModuleRemoved(log types.Log) (*TokenAdminRegistryRegistryModuleRemoved, error) {
	event := new(TokenAdminRegistryRegistryModuleRemoved)
	if err := _TokenAdminRegistry.contract.UnpackLog(event, "RegistryModuleRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_TokenAdminRegistry *TokenAdminRegistry) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _TokenAdminRegistry.abi.Events["AdministratorTransferRequested"].ID:
		return _TokenAdminRegistry.ParseAdministratorTransferRequested(log)
	case _TokenAdminRegistry.abi.Events["AdministratorTransferred"].ID:
		return _TokenAdminRegistry.ParseAdministratorTransferred(log)
	case _TokenAdminRegistry.abi.Events["OwnershipTransferRequested"].ID:
		return _TokenAdminRegistry.ParseOwnershipTransferRequested(log)
	case _TokenAdminRegistry.abi.Events["OwnershipTransferred"].ID:
		return _TokenAdminRegistry.ParseOwnershipTransferred(log)
	case _TokenAdminRegistry.abi.Events["PoolSet"].ID:
		return _TokenAdminRegistry.ParsePoolSet(log)
	case _TokenAdminRegistry.abi.Events["RegistryModuleAdded"].ID:
		return _TokenAdminRegistry.ParseRegistryModuleAdded(log)
	case _TokenAdminRegistry.abi.Events["RegistryModuleRemoved"].ID:
		return _TokenAdminRegistry.ParseRegistryModuleRemoved(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (TokenAdminRegistryAdministratorTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xc54c3051ff16e63bb9203214432372aca006c589e3653619b577a3265675b716")
}

func (TokenAdminRegistryAdministratorTransferred) Topic() common.Hash {
	return common.HexToHash("0x399b55200f7f639a63d76efe3dcfa9156ce367058d6b673041b84a628885f5a7")
}

func (TokenAdminRegistryOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (TokenAdminRegistryOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (TokenAdminRegistryPoolSet) Topic() common.Hash {
	return common.HexToHash("0x754449ec3aff3bd528bfce43ae9319c4a381b67fcd1d20097b3b24dacaecc35d")
}

func (TokenAdminRegistryRegistryModuleAdded) Topic() common.Hash {
	return common.HexToHash("0x3cabf004338366bfeaeb610ad827cb58d16b588017c509501f2c97c83caae7b2")
}

func (TokenAdminRegistryRegistryModuleRemoved) Topic() common.Hash {
	return common.HexToHash("0x93eaa26dcb9275e56bacb1d33fdbf402262da6f0f4baf2a6e2cd154b73f387f8")
}

func (_TokenAdminRegistry *TokenAdminRegistry) Address() common.Address {
	return _TokenAdminRegistry.address
}

type TokenAdminRegistryInterface interface {
	GetAllConfiguredTokens(opts *bind.CallOpts, startIndex uint64, maxCount uint64) ([]common.Address, error)

	GetPool(opts *bind.CallOpts, token common.Address) (common.Address, error)

	GetPools(opts *bind.CallOpts, tokens []common.Address) ([]common.Address, error)

	GetTokenConfig(opts *bind.CallOpts, token common.Address) (TokenAdminRegistryTokenConfig, error)

	IsAdministrator(opts *bind.CallOpts, localToken common.Address, administrator common.Address) (bool, error)

	IsRegistryModule(opts *bind.CallOpts, module common.Address) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptAdminRole(opts *bind.TransactOpts, localToken common.Address) (*types.Transaction, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AddRegistryModule(opts *bind.TransactOpts, module common.Address) (*types.Transaction, error)

	ProposeAdministrator(opts *bind.TransactOpts, localToken common.Address, administrator common.Address) (*types.Transaction, error)

	RemoveRegistryModule(opts *bind.TransactOpts, module common.Address) (*types.Transaction, error)

	SetPool(opts *bind.TransactOpts, localToken common.Address, pool common.Address) (*types.Transaction, error)

	TransferAdminRole(opts *bind.TransactOpts, localToken common.Address, newAdmin common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterAdministratorTransferRequested(opts *bind.FilterOpts, token []common.Address, currentAdmin []common.Address, newAdmin []common.Address) (*TokenAdminRegistryAdministratorTransferRequestedIterator, error)

	WatchAdministratorTransferRequested(opts *bind.WatchOpts, sink chan<- *TokenAdminRegistryAdministratorTransferRequested, token []common.Address, currentAdmin []common.Address, newAdmin []common.Address) (event.Subscription, error)

	ParseAdministratorTransferRequested(log types.Log) (*TokenAdminRegistryAdministratorTransferRequested, error)

	FilterAdministratorTransferred(opts *bind.FilterOpts, token []common.Address, newAdmin []common.Address) (*TokenAdminRegistryAdministratorTransferredIterator, error)

	WatchAdministratorTransferred(opts *bind.WatchOpts, sink chan<- *TokenAdminRegistryAdministratorTransferred, token []common.Address, newAdmin []common.Address) (event.Subscription, error)

	ParseAdministratorTransferred(log types.Log) (*TokenAdminRegistryAdministratorTransferred, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TokenAdminRegistryOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *TokenAdminRegistryOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*TokenAdminRegistryOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TokenAdminRegistryOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TokenAdminRegistryOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*TokenAdminRegistryOwnershipTransferred, error)

	FilterPoolSet(opts *bind.FilterOpts, token []common.Address, previousPool []common.Address, newPool []common.Address) (*TokenAdminRegistryPoolSetIterator, error)

	WatchPoolSet(opts *bind.WatchOpts, sink chan<- *TokenAdminRegistryPoolSet, token []common.Address, previousPool []common.Address, newPool []common.Address) (event.Subscription, error)

	ParsePoolSet(log types.Log) (*TokenAdminRegistryPoolSet, error)

	FilterRegistryModuleAdded(opts *bind.FilterOpts) (*TokenAdminRegistryRegistryModuleAddedIterator, error)

	WatchRegistryModuleAdded(opts *bind.WatchOpts, sink chan<- *TokenAdminRegistryRegistryModuleAdded) (event.Subscription, error)

	ParseRegistryModuleAdded(log types.Log) (*TokenAdminRegistryRegistryModuleAdded, error)

	FilterRegistryModuleRemoved(opts *bind.FilterOpts, module []common.Address) (*TokenAdminRegistryRegistryModuleRemovedIterator, error)

	WatchRegistryModuleRemoved(opts *bind.WatchOpts, sink chan<- *TokenAdminRegistryRegistryModuleRemoved, module []common.Address) (event.Subscription, error)

	ParseRegistryModuleRemoved(log types.Log) (*TokenAdminRegistryRegistryModuleRemoved, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
