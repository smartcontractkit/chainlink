// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package mercury_exposed_verifier

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

var MercuryExposedVerifierMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_feedId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_chainId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_contractAddress\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"_configCount\",\"type\":\"uint64\"},{\"internalType\":\"address[]\",\"name\":\"_signers\",\"type\":\"address[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"_offchainTransmitters\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint8\",\"name\":\"_f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"_onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"_encodedConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"_encodedConfig\",\"type\":\"bytes\"}],\"name\":\"exposedConfigDigestFromConfigData\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610696806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c80630ebd702314610030575b600080fd5b61004361003e3660046103f7565b610055565b60405190815260200160405180910390f35b60006100a18c8c8c8c8c8c8c8c8c8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508e92508d91506100b19050565b9c9b505050505050505050505050565b6000808b8b8b8b8b8b8b8b8b8b6040516020016100d79a999897969594939291906105a7565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179150509a9950505050505050505050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461018357600080fd5b919050565b803567ffffffffffffffff8116811461018357600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610216576102166101a0565b604052919050565b600067ffffffffffffffff821115610238576102386101a0565b5060051b60200190565b600082601f83011261025357600080fd5b813560206102686102638361021e565b6101cf565b82815260059290921b8401810191818101908684111561028757600080fd5b8286015b848110156102a95761029c8161015f565b835291830191830161028b565b509695505050505050565b600082601f8301126102c557600080fd5b813560206102d56102638361021e565b82815260059290921b840181019181810190868411156102f457600080fd5b8286015b848110156102a957803583529183019183016102f8565b803560ff8116811461018357600080fd5b60008083601f84011261033257600080fd5b50813567ffffffffffffffff81111561034a57600080fd5b60208301915083602082850101111561036257600080fd5b9250929050565b600082601f83011261037a57600080fd5b813567ffffffffffffffff811115610394576103946101a0565b6103c560207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116016101cf565b8181528460208386010111156103da57600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008060008060008060008060006101408c8e03121561041957600080fd5b8b359a5060208c0135995061043060408d0161015f565b985061043e60608d01610188565b975067ffffffffffffffff8060808e0135111561045a57600080fd5b61046a8e60808f01358f01610242565b97508060a08e0135111561047d57600080fd5b61048d8e60a08f01358f016102b4565b965061049b60c08e0161030f565b95508060e08e013511156104ae57600080fd5b6104be8e60e08f01358f01610320565b90955093506104d06101008e01610188565b9250806101208e013511156104e457600080fd5b506104f68d6101208e01358e01610369565b90509295989b509295989b9093969950565b600081518084526020808501945080840160005b838110156105385781518752958201959082019060010161051c565b509495945050505050565b6000815180845260005b818110156105695760208185018101518683018201520161054d565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b8a815260208082018b905273ffffffffffffffffffffffffffffffffffffffff8a8116604084015267ffffffffffffffff8a1660608401526101406080840181905289519084018190526000926101608501928b820192855b8181101561061e578451831686529483019493830193600101610600565b505050505082810360a08401526106358189610508565b60ff881660c0850152905082810360e08401526106528187610543565b67ffffffffffffffff861661010085015290508281036101208401526106788185610543565b9d9c5050505050505050505050505056fea164736f6c6343000810000a",
}

var MercuryExposedVerifierABI = MercuryExposedVerifierMetaData.ABI

var MercuryExposedVerifierBin = MercuryExposedVerifierMetaData.Bin

func DeployMercuryExposedVerifier(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *MercuryExposedVerifier, error) {
	parsed, err := MercuryExposedVerifierMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MercuryExposedVerifierBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MercuryExposedVerifier{MercuryExposedVerifierCaller: MercuryExposedVerifierCaller{contract: contract}, MercuryExposedVerifierTransactor: MercuryExposedVerifierTransactor{contract: contract}, MercuryExposedVerifierFilterer: MercuryExposedVerifierFilterer{contract: contract}}, nil
}

type MercuryExposedVerifier struct {
	address common.Address
	abi     abi.ABI
	MercuryExposedVerifierCaller
	MercuryExposedVerifierTransactor
	MercuryExposedVerifierFilterer
}

type MercuryExposedVerifierCaller struct {
	contract *bind.BoundContract
}

type MercuryExposedVerifierTransactor struct {
	contract *bind.BoundContract
}

type MercuryExposedVerifierFilterer struct {
	contract *bind.BoundContract
}

type MercuryExposedVerifierSession struct {
	Contract     *MercuryExposedVerifier
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type MercuryExposedVerifierCallerSession struct {
	Contract *MercuryExposedVerifierCaller
	CallOpts bind.CallOpts
}

type MercuryExposedVerifierTransactorSession struct {
	Contract     *MercuryExposedVerifierTransactor
	TransactOpts bind.TransactOpts
}

type MercuryExposedVerifierRaw struct {
	Contract *MercuryExposedVerifier
}

type MercuryExposedVerifierCallerRaw struct {
	Contract *MercuryExposedVerifierCaller
}

type MercuryExposedVerifierTransactorRaw struct {
	Contract *MercuryExposedVerifierTransactor
}

func NewMercuryExposedVerifier(address common.Address, backend bind.ContractBackend) (*MercuryExposedVerifier, error) {
	abi, err := abi.JSON(strings.NewReader(MercuryExposedVerifierABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindMercuryExposedVerifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MercuryExposedVerifier{address: address, abi: abi, MercuryExposedVerifierCaller: MercuryExposedVerifierCaller{contract: contract}, MercuryExposedVerifierTransactor: MercuryExposedVerifierTransactor{contract: contract}, MercuryExposedVerifierFilterer: MercuryExposedVerifierFilterer{contract: contract}}, nil
}

func NewMercuryExposedVerifierCaller(address common.Address, caller bind.ContractCaller) (*MercuryExposedVerifierCaller, error) {
	contract, err := bindMercuryExposedVerifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MercuryExposedVerifierCaller{contract: contract}, nil
}

func NewMercuryExposedVerifierTransactor(address common.Address, transactor bind.ContractTransactor) (*MercuryExposedVerifierTransactor, error) {
	contract, err := bindMercuryExposedVerifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MercuryExposedVerifierTransactor{contract: contract}, nil
}

func NewMercuryExposedVerifierFilterer(address common.Address, filterer bind.ContractFilterer) (*MercuryExposedVerifierFilterer, error) {
	contract, err := bindMercuryExposedVerifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MercuryExposedVerifierFilterer{contract: contract}, nil
}

func bindMercuryExposedVerifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MercuryExposedVerifierMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_MercuryExposedVerifier *MercuryExposedVerifierRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MercuryExposedVerifier.Contract.MercuryExposedVerifierCaller.contract.Call(opts, result, method, params...)
}

func (_MercuryExposedVerifier *MercuryExposedVerifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MercuryExposedVerifier.Contract.MercuryExposedVerifierTransactor.contract.Transfer(opts)
}

func (_MercuryExposedVerifier *MercuryExposedVerifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MercuryExposedVerifier.Contract.MercuryExposedVerifierTransactor.contract.Transact(opts, method, params...)
}

func (_MercuryExposedVerifier *MercuryExposedVerifierCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MercuryExposedVerifier.Contract.contract.Call(opts, result, method, params...)
}

func (_MercuryExposedVerifier *MercuryExposedVerifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MercuryExposedVerifier.Contract.contract.Transfer(opts)
}

func (_MercuryExposedVerifier *MercuryExposedVerifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MercuryExposedVerifier.Contract.contract.Transact(opts, method, params...)
}

func (_MercuryExposedVerifier *MercuryExposedVerifierCaller) ExposedConfigDigestFromConfigData(opts *bind.CallOpts, _feedId [32]byte, _chainId *big.Int, _contractAddress common.Address, _configCount uint64, _signers []common.Address, _offchainTransmitters [][32]byte, _f uint8, _onchainConfig []byte, _encodedConfigVersion uint64, _encodedConfig []byte) ([32]byte, error) {
	var out []interface{}
	err := _MercuryExposedVerifier.contract.Call(opts, &out, "exposedConfigDigestFromConfigData", _feedId, _chainId, _contractAddress, _configCount, _signers, _offchainTransmitters, _f, _onchainConfig, _encodedConfigVersion, _encodedConfig)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_MercuryExposedVerifier *MercuryExposedVerifierSession) ExposedConfigDigestFromConfigData(_feedId [32]byte, _chainId *big.Int, _contractAddress common.Address, _configCount uint64, _signers []common.Address, _offchainTransmitters [][32]byte, _f uint8, _onchainConfig []byte, _encodedConfigVersion uint64, _encodedConfig []byte) ([32]byte, error) {
	return _MercuryExposedVerifier.Contract.ExposedConfigDigestFromConfigData(&_MercuryExposedVerifier.CallOpts, _feedId, _chainId, _contractAddress, _configCount, _signers, _offchainTransmitters, _f, _onchainConfig, _encodedConfigVersion, _encodedConfig)
}

func (_MercuryExposedVerifier *MercuryExposedVerifierCallerSession) ExposedConfigDigestFromConfigData(_feedId [32]byte, _chainId *big.Int, _contractAddress common.Address, _configCount uint64, _signers []common.Address, _offchainTransmitters [][32]byte, _f uint8, _onchainConfig []byte, _encodedConfigVersion uint64, _encodedConfig []byte) ([32]byte, error) {
	return _MercuryExposedVerifier.Contract.ExposedConfigDigestFromConfigData(&_MercuryExposedVerifier.CallOpts, _feedId, _chainId, _contractAddress, _configCount, _signers, _offchainTransmitters, _f, _onchainConfig, _encodedConfigVersion, _encodedConfig)
}

func (_MercuryExposedVerifier *MercuryExposedVerifier) Address() common.Address {
	return _MercuryExposedVerifier.address
}

type MercuryExposedVerifierInterface interface {
	ExposedConfigDigestFromConfigData(opts *bind.CallOpts, _feedId [32]byte, _chainId *big.Int, _contractAddress common.Address, _configCount uint64, _signers []common.Address, _offchainTransmitters [][32]byte, _f uint8, _onchainConfig []byte, _encodedConfigVersion uint64, _encodedConfig []byte) ([32]byte, error)

	Address() common.Address
}
