// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package exposed_channel_verifier

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

var ExposedChannelVerifierMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_chainId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_contractAddress\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"_configCount\",\"type\":\"uint64\"},{\"internalType\":\"address[]\",\"name\":\"_signers\",\"type\":\"address[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"_offchainTransmitters\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint8\",\"name\":\"_f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"_onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"_encodedConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"_encodedConfig\",\"type\":\"bytes\"}],\"name\":\"exposedConfigDigestFromConfigData\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5061067e806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c8063b05a355014610030575b600080fd5b61004361003e3660046103f2565b610055565b60405190815260200160405180910390f35b60006100a08b8b8b8b8b8b8b8b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508d92508c91506100af9050565b9b9a5050505050505050505050565b6000808a8a8a8a8a8a8a8a8a6040516020016100d399989796959493929190610594565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e09000000000000000000000000000000000000000000000000000000000000179150509998505050505050505050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461017e57600080fd5b919050565b803567ffffffffffffffff8116811461017e57600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156102115761021161019b565b604052919050565b600067ffffffffffffffff8211156102335761023361019b565b5060051b60200190565b600082601f83011261024e57600080fd5b8135602061026361025e83610219565b6101ca565b82815260059290921b8401810191818101908684111561028257600080fd5b8286015b848110156102a4576102978161015a565b8352918301918301610286565b509695505050505050565b600082601f8301126102c057600080fd5b813560206102d061025e83610219565b82815260059290921b840181019181810190868411156102ef57600080fd5b8286015b848110156102a457803583529183019183016102f3565b803560ff8116811461017e57600080fd5b60008083601f84011261032d57600080fd5b50813567ffffffffffffffff81111561034557600080fd5b60208301915083602082850101111561035d57600080fd5b9250929050565b600082601f83011261037557600080fd5b813567ffffffffffffffff81111561038f5761038f61019b565b6103c060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116016101ca565b8181528460208386010111156103d557600080fd5b816020850160208301376000918101602001919091529392505050565b6000806000806000806000806000806101208b8d03121561041257600080fd5b8a35995061042260208c0161015a565b985061043060408c01610183565b975060608b013567ffffffffffffffff8082111561044d57600080fd5b6104598e838f0161023d565b985060808d013591508082111561046f57600080fd5b61047b8e838f016102af565b975061048960a08e0161030a565b965060c08d013591508082111561049f57600080fd5b6104ab8e838f0161031b565b90965094508491506104bf60e08e01610183565b93506101008d01359150808211156104d657600080fd5b506104e38d828e01610364565b9150509295989b9194979a5092959850565b600081518084526020808501945080840160005b8381101561052557815187529582019590820190600101610509565b509495945050505050565b6000815180845260005b818110156105565760208185018101518683018201520161053a565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b60006101208083018c8452602073ffffffffffffffffffffffffffffffffffffffff808e168287015267ffffffffffffffff8d1660408701528360608701528293508b5180845261014087019450828d01935060005b818110156106085784518316865294830194938301936001016105ea565b5050505050828103608084015261061f81896104f5565b60ff881660a0850152905082810360c084015261063c8187610530565b67ffffffffffffffff861660e085015290508281036101008401526106618185610530565b9c9b50505050505050505050505056fea164736f6c6343000813000a",
}

var ExposedChannelVerifierABI = ExposedChannelVerifierMetaData.ABI

var ExposedChannelVerifierBin = ExposedChannelVerifierMetaData.Bin

func DeployExposedChannelVerifier(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ExposedChannelVerifier, error) {
	parsed, err := ExposedChannelVerifierMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ExposedChannelVerifierBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ExposedChannelVerifier{address: address, abi: *parsed, ExposedChannelVerifierCaller: ExposedChannelVerifierCaller{contract: contract}, ExposedChannelVerifierTransactor: ExposedChannelVerifierTransactor{contract: contract}, ExposedChannelVerifierFilterer: ExposedChannelVerifierFilterer{contract: contract}}, nil
}

type ExposedChannelVerifier struct {
	address common.Address
	abi     abi.ABI
	ExposedChannelVerifierCaller
	ExposedChannelVerifierTransactor
	ExposedChannelVerifierFilterer
}

type ExposedChannelVerifierCaller struct {
	contract *bind.BoundContract
}

type ExposedChannelVerifierTransactor struct {
	contract *bind.BoundContract
}

type ExposedChannelVerifierFilterer struct {
	contract *bind.BoundContract
}

type ExposedChannelVerifierSession struct {
	Contract     *ExposedChannelVerifier
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type ExposedChannelVerifierCallerSession struct {
	Contract *ExposedChannelVerifierCaller
	CallOpts bind.CallOpts
}

type ExposedChannelVerifierTransactorSession struct {
	Contract     *ExposedChannelVerifierTransactor
	TransactOpts bind.TransactOpts
}

type ExposedChannelVerifierRaw struct {
	Contract *ExposedChannelVerifier
}

type ExposedChannelVerifierCallerRaw struct {
	Contract *ExposedChannelVerifierCaller
}

type ExposedChannelVerifierTransactorRaw struct {
	Contract *ExposedChannelVerifierTransactor
}

func NewExposedChannelVerifier(address common.Address, backend bind.ContractBackend) (*ExposedChannelVerifier, error) {
	abi, err := abi.JSON(strings.NewReader(ExposedChannelVerifierABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindExposedChannelVerifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ExposedChannelVerifier{address: address, abi: abi, ExposedChannelVerifierCaller: ExposedChannelVerifierCaller{contract: contract}, ExposedChannelVerifierTransactor: ExposedChannelVerifierTransactor{contract: contract}, ExposedChannelVerifierFilterer: ExposedChannelVerifierFilterer{contract: contract}}, nil
}

func NewExposedChannelVerifierCaller(address common.Address, caller bind.ContractCaller) (*ExposedChannelVerifierCaller, error) {
	contract, err := bindExposedChannelVerifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ExposedChannelVerifierCaller{contract: contract}, nil
}

func NewExposedChannelVerifierTransactor(address common.Address, transactor bind.ContractTransactor) (*ExposedChannelVerifierTransactor, error) {
	contract, err := bindExposedChannelVerifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ExposedChannelVerifierTransactor{contract: contract}, nil
}

func NewExposedChannelVerifierFilterer(address common.Address, filterer bind.ContractFilterer) (*ExposedChannelVerifierFilterer, error) {
	contract, err := bindExposedChannelVerifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ExposedChannelVerifierFilterer{contract: contract}, nil
}

func bindExposedChannelVerifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ExposedChannelVerifierMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_ExposedChannelVerifier *ExposedChannelVerifierRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ExposedChannelVerifier.Contract.ExposedChannelVerifierCaller.contract.Call(opts, result, method, params...)
}

func (_ExposedChannelVerifier *ExposedChannelVerifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ExposedChannelVerifier.Contract.ExposedChannelVerifierTransactor.contract.Transfer(opts)
}

func (_ExposedChannelVerifier *ExposedChannelVerifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ExposedChannelVerifier.Contract.ExposedChannelVerifierTransactor.contract.Transact(opts, method, params...)
}

func (_ExposedChannelVerifier *ExposedChannelVerifierCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ExposedChannelVerifier.Contract.contract.Call(opts, result, method, params...)
}

func (_ExposedChannelVerifier *ExposedChannelVerifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ExposedChannelVerifier.Contract.contract.Transfer(opts)
}

func (_ExposedChannelVerifier *ExposedChannelVerifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ExposedChannelVerifier.Contract.contract.Transact(opts, method, params...)
}

func (_ExposedChannelVerifier *ExposedChannelVerifierCaller) ExposedConfigDigestFromConfigData(opts *bind.CallOpts, _chainId *big.Int, _contractAddress common.Address, _configCount uint64, _signers []common.Address, _offchainTransmitters [][32]byte, _f uint8, _onchainConfig []byte, _encodedConfigVersion uint64, _encodedConfig []byte) ([32]byte, error) {
	var out []interface{}
	err := _ExposedChannelVerifier.contract.Call(opts, &out, "exposedConfigDigestFromConfigData", _chainId, _contractAddress, _configCount, _signers, _offchainTransmitters, _f, _onchainConfig, _encodedConfigVersion, _encodedConfig)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_ExposedChannelVerifier *ExposedChannelVerifierSession) ExposedConfigDigestFromConfigData(_chainId *big.Int, _contractAddress common.Address, _configCount uint64, _signers []common.Address, _offchainTransmitters [][32]byte, _f uint8, _onchainConfig []byte, _encodedConfigVersion uint64, _encodedConfig []byte) ([32]byte, error) {
	return _ExposedChannelVerifier.Contract.ExposedConfigDigestFromConfigData(&_ExposedChannelVerifier.CallOpts, _chainId, _contractAddress, _configCount, _signers, _offchainTransmitters, _f, _onchainConfig, _encodedConfigVersion, _encodedConfig)
}

func (_ExposedChannelVerifier *ExposedChannelVerifierCallerSession) ExposedConfigDigestFromConfigData(_chainId *big.Int, _contractAddress common.Address, _configCount uint64, _signers []common.Address, _offchainTransmitters [][32]byte, _f uint8, _onchainConfig []byte, _encodedConfigVersion uint64, _encodedConfig []byte) ([32]byte, error) {
	return _ExposedChannelVerifier.Contract.ExposedConfigDigestFromConfigData(&_ExposedChannelVerifier.CallOpts, _chainId, _contractAddress, _configCount, _signers, _offchainTransmitters, _f, _onchainConfig, _encodedConfigVersion, _encodedConfig)
}

func (_ExposedChannelVerifier *ExposedChannelVerifier) Address() common.Address {
	return _ExposedChannelVerifier.address
}

type ExposedChannelVerifierInterface interface {
	ExposedConfigDigestFromConfigData(opts *bind.CallOpts, _chainId *big.Int, _contractAddress common.Address, _configCount uint64, _signers []common.Address, _offchainTransmitters [][32]byte, _f uint8, _onchainConfig []byte, _encodedConfigVersion uint64, _encodedConfig []byte) ([32]byte, error)

	Address() common.Address
}
