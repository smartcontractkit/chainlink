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
	Bin: "0x608060405234801561001057600080fd5b506106a9806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c80630ebd702314610030575b600080fd5b61004361003e366004610361565b610055565b60405190815260200160405180910390f35b60006100a18c8c8c8c8c8c8c8c8c8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508e92508d91506100b19050565b9c9b505050505050505050505050565b6000808b8b8b8b8b8b8b8b8b8b6040516020016100d79a99989796959493929190610518565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179150509a9950505050505050505050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461018357600080fd5b919050565b600082601f83011261019957600080fd5b813560206101ae6101a983610649565b6105fa565b80838252828201915082860187848660051b89010111156101ce57600080fd5b60005b858110156101f4576101e28261015f565b845292840192908401906001016101d1565b5090979650505050505050565b600082601f83011261021257600080fd5b813560206102226101a983610649565b80838252828201915082860187848660051b890101111561024257600080fd5b60005b858110156101f457813584529284019290840190600101610245565b60008083601f84011261027357600080fd5b50813567ffffffffffffffff81111561028b57600080fd5b6020830191508360208285010111156102a357600080fd5b9250929050565b600082601f8301126102bb57600080fd5b813567ffffffffffffffff8111156102d5576102d561066d565b61030660207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116016105fa565b81815284602083860101111561031b57600080fd5b816020850160208301376000918101602001919091529392505050565b803567ffffffffffffffff8116811461018357600080fd5b803560ff8116811461018357600080fd5b60008060008060008060008060008060006101408c8e03121561038357600080fd5b8b359a5060208c0135995061039a60408d0161015f565b98506103a860608d01610338565b975067ffffffffffffffff8060808e013511156103c457600080fd5b6103d48e60808f01358f01610188565b97508060a08e013511156103e757600080fd5b6103f78e60a08f01358f01610201565b965061040560c08e01610350565b95508060e08e0135111561041857600080fd5b6104288e60e08f01358f01610261565b909550935061043a6101008e01610338565b9250806101208e0135111561044e57600080fd5b506104608d6101208e01358e016102aa565b90509295989b509295989b9093969950565b600081518084526020808501945080840160005b838110156104a257815187529582019590820190600101610486565b509495945050505050565b6000815180845260005b818110156104d3576020818501810151868301820152016104b7565b818111156104e5576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b8a815260208082018b905273ffffffffffffffffffffffffffffffffffffffff8a8116604084015267ffffffffffffffff8a1660608401526101406080840181905289519084018190526000926101608501928b820192855b8181101561058f578451831686529483019493830193600101610571565b505050505082810360a08401526105a68189610472565b60ff881660c0850152905082810360e08401526105c381876104ad565b67ffffffffffffffff861661010085015290508281036101208401526105e981856104ad565b9d9c50505050505050505050505050565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156106415761064161066d565b604052919050565b600067ffffffffffffffff8211156106635761066361066d565b5060051b60200190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
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
