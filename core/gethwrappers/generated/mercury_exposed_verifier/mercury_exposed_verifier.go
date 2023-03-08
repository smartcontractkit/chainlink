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
)

var MercuryExposedVerifierMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_chainId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_contractAddress\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"_configCount\",\"type\":\"uint64\"},{\"internalType\":\"address[]\",\"name\":\"_signers\",\"type\":\"address[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"_offchainTransmitters\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint8\",\"name\":\"_f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"_onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"_encodedConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"_encodedConfig\",\"type\":\"bytes\"}],\"name\":\"exposedConfigDigestFromConfigData\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610691806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c8063b05a355014610030575b600080fd5b61004361003e36600461035c565b610055565b60405190815260200160405180910390f35b60006100a08b8b8b8b8b8b8b8b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508d92508c91506100af9050565b9b9a5050505050505050505050565b6000808a8a8a8a8a8a8a8a8a6040516020016100d399989796959493929190610505565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179150509998505050505050505050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461017e57600080fd5b919050565b600082601f83011261019457600080fd5b813560206101a96101a483610631565b6105e2565b80838252828201915082860187848660051b89010111156101c957600080fd5b60005b858110156101ef576101dd8261015a565b845292840192908401906001016101cc565b5090979650505050505050565b600082601f83011261020d57600080fd5b8135602061021d6101a483610631565b80838252828201915082860187848660051b890101111561023d57600080fd5b60005b858110156101ef57813584529284019290840190600101610240565b60008083601f84011261026e57600080fd5b50813567ffffffffffffffff81111561028657600080fd5b60208301915083602082850101111561029e57600080fd5b9250929050565b600082601f8301126102b657600080fd5b813567ffffffffffffffff8111156102d0576102d0610655565b61030160207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116016105e2565b81815284602083860101111561031657600080fd5b816020850160208301376000918101602001919091529392505050565b803567ffffffffffffffff8116811461017e57600080fd5b803560ff8116811461017e57600080fd5b6000806000806000806000806000806101208b8d03121561037c57600080fd5b8a35995061038c60208c0161015a565b985061039a60408c01610333565b975060608b013567ffffffffffffffff808211156103b757600080fd5b6103c38e838f01610183565b985060808d01359150808211156103d957600080fd5b6103e58e838f016101fc565b97506103f360a08e0161034b565b965060c08d013591508082111561040957600080fd5b6104158e838f0161025c565b909650945084915061042960e08e01610333565b93506101008d013591508082111561044057600080fd5b5061044d8d828e016102a5565b9150509295989b9194979a5092959850565b600081518084526020808501945080840160005b8381101561048f57815187529582019590820190600101610473565b509495945050505050565b6000815180845260005b818110156104c0576020818501810151868301820152016104a4565b818111156104d2576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b60006101208083018c8452602073ffffffffffffffffffffffffffffffffffffffff808e168287015267ffffffffffffffff8d1660408701528360608701528293508b5180845261014087019450828d01935060005b8181101561057957845183168652948301949383019360010161055b565b50505050508281036080840152610590818961045f565b60ff881660a0850152905082810360c08401526105ad818761049a565b67ffffffffffffffff861660e085015290508281036101008401526105d2818561049a565b9c9b505050505050505050505050565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561062957610629610655565b604052919050565b600067ffffffffffffffff82111561064b5761064b610655565b5060051b60200190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
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
	parsed, err := abi.JSON(strings.NewReader(MercuryExposedVerifierABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
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

func (_MercuryExposedVerifier *MercuryExposedVerifierCaller) ExposedConfigDigestFromConfigData(opts *bind.CallOpts, _chainId *big.Int, _contractAddress common.Address, _configCount uint64, _signers []common.Address, _offchainTransmitters [][32]byte, _f uint8, _onchainConfig []byte, _encodedConfigVersion uint64, _encodedConfig []byte) ([32]byte, error) {
	var out []interface{}
	err := _MercuryExposedVerifier.contract.Call(opts, &out, "exposedConfigDigestFromConfigData", _chainId, _contractAddress, _configCount, _signers, _offchainTransmitters, _f, _onchainConfig, _encodedConfigVersion, _encodedConfig)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_MercuryExposedVerifier *MercuryExposedVerifierSession) ExposedConfigDigestFromConfigData(_chainId *big.Int, _contractAddress common.Address, _configCount uint64, _signers []common.Address, _offchainTransmitters [][32]byte, _f uint8, _onchainConfig []byte, _encodedConfigVersion uint64, _encodedConfig []byte) ([32]byte, error) {
	return _MercuryExposedVerifier.Contract.ExposedConfigDigestFromConfigData(&_MercuryExposedVerifier.CallOpts, _chainId, _contractAddress, _configCount, _signers, _offchainTransmitters, _f, _onchainConfig, _encodedConfigVersion, _encodedConfig)
}

func (_MercuryExposedVerifier *MercuryExposedVerifierCallerSession) ExposedConfigDigestFromConfigData(_chainId *big.Int, _contractAddress common.Address, _configCount uint64, _signers []common.Address, _offchainTransmitters [][32]byte, _f uint8, _onchainConfig []byte, _encodedConfigVersion uint64, _encodedConfig []byte) ([32]byte, error) {
	return _MercuryExposedVerifier.Contract.ExposedConfigDigestFromConfigData(&_MercuryExposedVerifier.CallOpts, _chainId, _contractAddress, _configCount, _signers, _offchainTransmitters, _f, _onchainConfig, _encodedConfigVersion, _encodedConfig)
}

func (_MercuryExposedVerifier *MercuryExposedVerifier) Address() common.Address {
	return _MercuryExposedVerifier.address
}

type MercuryExposedVerifierInterface interface {
	ExposedConfigDigestFromConfigData(opts *bind.CallOpts, _chainId *big.Int, _contractAddress common.Address, _configCount uint64, _signers []common.Address, _offchainTransmitters [][32]byte, _f uint8, _onchainConfig []byte, _encodedConfigVersion uint64, _encodedConfig []byte) ([32]byte, error)

	Address() common.Address
}
