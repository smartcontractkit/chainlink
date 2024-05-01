// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package optimism_l1_bridge_adapter_encoder

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

type OptimismL1BridgeAdapterFinalizeWithdrawERC20Payload struct {
	Action uint8
	Data   []byte
}

type OptimismL1BridgeAdapterOptimismFinalizationPayload struct {
	WithdrawalTransaction TypesWithdrawalTransaction
}

type OptimismL1BridgeAdapterOptimismProveWithdrawalPayload struct {
	WithdrawalTransaction TypesWithdrawalTransaction
	L2OutputIndex         *big.Int
	OutputRootProof       TypesOutputRootProof
	WithdrawalProof       [][]byte
}

type TypesOutputRootProof struct {
	Version                  [32]byte
	StateRoot                [32]byte
	MessagePasserStorageRoot [32]byte
	LatestBlockhash          [32]byte
}

type TypesWithdrawalTransaction struct {
	Nonce    *big.Int
	Sender   common.Address
	Target   common.Address
	Value    *big.Int
	GasLimit *big.Int
	Data     []byte
}

var OptimismL1BridgeAdapterEncoderMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"enumOptimismL1BridgeAdapter.FinalizationAction\",\"name\":\"action\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structOptimismL1BridgeAdapter.FinalizeWithdrawERC20Payload\",\"name\":\"payload\",\"type\":\"tuple\"}],\"name\":\"encodeFinalizeWithdrawalERC20Payload\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structTypes.WithdrawalTransaction\",\"name\":\"withdrawalTransaction\",\"type\":\"tuple\"}],\"internalType\":\"structOptimismL1BridgeAdapter.OptimismFinalizationPayload\",\"name\":\"payload\",\"type\":\"tuple\"}],\"name\":\"encodeOptimismFinalizationPayload\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structTypes.WithdrawalTransaction\",\"name\":\"withdrawalTransaction\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"l2OutputIndex\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"version\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"messagePasserStorageRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"latestBlockhash\",\"type\":\"bytes32\"}],\"internalType\":\"structTypes.OutputRootProof\",\"name\":\"outputRootProof\",\"type\":\"tuple\"},{\"internalType\":\"bytes[]\",\"name\":\"withdrawalProof\",\"type\":\"bytes[]\"}],\"internalType\":\"structOptimismL1BridgeAdapter.OptimismProveWithdrawalPayload\",\"name\":\"payload\",\"type\":\"tuple\"}],\"name\":\"encodeOptimismProveWithdrawalPayload\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
}

var OptimismL1BridgeAdapterEncoderABI = OptimismL1BridgeAdapterEncoderMetaData.ABI

type OptimismL1BridgeAdapterEncoder struct {
	address common.Address
	abi     abi.ABI
	OptimismL1BridgeAdapterEncoderCaller
	OptimismL1BridgeAdapterEncoderTransactor
	OptimismL1BridgeAdapterEncoderFilterer
}

type OptimismL1BridgeAdapterEncoderCaller struct {
	contract *bind.BoundContract
}

type OptimismL1BridgeAdapterEncoderTransactor struct {
	contract *bind.BoundContract
}

type OptimismL1BridgeAdapterEncoderFilterer struct {
	contract *bind.BoundContract
}

type OptimismL1BridgeAdapterEncoderSession struct {
	Contract     *OptimismL1BridgeAdapterEncoder
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type OptimismL1BridgeAdapterEncoderCallerSession struct {
	Contract *OptimismL1BridgeAdapterEncoderCaller
	CallOpts bind.CallOpts
}

type OptimismL1BridgeAdapterEncoderTransactorSession struct {
	Contract     *OptimismL1BridgeAdapterEncoderTransactor
	TransactOpts bind.TransactOpts
}

type OptimismL1BridgeAdapterEncoderRaw struct {
	Contract *OptimismL1BridgeAdapterEncoder
}

type OptimismL1BridgeAdapterEncoderCallerRaw struct {
	Contract *OptimismL1BridgeAdapterEncoderCaller
}

type OptimismL1BridgeAdapterEncoderTransactorRaw struct {
	Contract *OptimismL1BridgeAdapterEncoderTransactor
}

func NewOptimismL1BridgeAdapterEncoder(address common.Address, backend bind.ContractBackend) (*OptimismL1BridgeAdapterEncoder, error) {
	abi, err := abi.JSON(strings.NewReader(OptimismL1BridgeAdapterEncoderABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindOptimismL1BridgeAdapterEncoder(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OptimismL1BridgeAdapterEncoder{address: address, abi: abi, OptimismL1BridgeAdapterEncoderCaller: OptimismL1BridgeAdapterEncoderCaller{contract: contract}, OptimismL1BridgeAdapterEncoderTransactor: OptimismL1BridgeAdapterEncoderTransactor{contract: contract}, OptimismL1BridgeAdapterEncoderFilterer: OptimismL1BridgeAdapterEncoderFilterer{contract: contract}}, nil
}

func NewOptimismL1BridgeAdapterEncoderCaller(address common.Address, caller bind.ContractCaller) (*OptimismL1BridgeAdapterEncoderCaller, error) {
	contract, err := bindOptimismL1BridgeAdapterEncoder(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OptimismL1BridgeAdapterEncoderCaller{contract: contract}, nil
}

func NewOptimismL1BridgeAdapterEncoderTransactor(address common.Address, transactor bind.ContractTransactor) (*OptimismL1BridgeAdapterEncoderTransactor, error) {
	contract, err := bindOptimismL1BridgeAdapterEncoder(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OptimismL1BridgeAdapterEncoderTransactor{contract: contract}, nil
}

func NewOptimismL1BridgeAdapterEncoderFilterer(address common.Address, filterer bind.ContractFilterer) (*OptimismL1BridgeAdapterEncoderFilterer, error) {
	contract, err := bindOptimismL1BridgeAdapterEncoder(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OptimismL1BridgeAdapterEncoderFilterer{contract: contract}, nil
}

func bindOptimismL1BridgeAdapterEncoder(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OptimismL1BridgeAdapterEncoderMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_OptimismL1BridgeAdapterEncoder *OptimismL1BridgeAdapterEncoderRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OptimismL1BridgeAdapterEncoder.Contract.OptimismL1BridgeAdapterEncoderCaller.contract.Call(opts, result, method, params...)
}

func (_OptimismL1BridgeAdapterEncoder *OptimismL1BridgeAdapterEncoderRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OptimismL1BridgeAdapterEncoder.Contract.OptimismL1BridgeAdapterEncoderTransactor.contract.Transfer(opts)
}

func (_OptimismL1BridgeAdapterEncoder *OptimismL1BridgeAdapterEncoderRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OptimismL1BridgeAdapterEncoder.Contract.OptimismL1BridgeAdapterEncoderTransactor.contract.Transact(opts, method, params...)
}

func (_OptimismL1BridgeAdapterEncoder *OptimismL1BridgeAdapterEncoderCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OptimismL1BridgeAdapterEncoder.Contract.contract.Call(opts, result, method, params...)
}

func (_OptimismL1BridgeAdapterEncoder *OptimismL1BridgeAdapterEncoderTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OptimismL1BridgeAdapterEncoder.Contract.contract.Transfer(opts)
}

func (_OptimismL1BridgeAdapterEncoder *OptimismL1BridgeAdapterEncoderTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OptimismL1BridgeAdapterEncoder.Contract.contract.Transact(opts, method, params...)
}

func (_OptimismL1BridgeAdapterEncoder *OptimismL1BridgeAdapterEncoderCaller) EncodeFinalizeWithdrawalERC20Payload(opts *bind.CallOpts, payload OptimismL1BridgeAdapterFinalizeWithdrawERC20Payload) error {
	var out []interface{}
	err := _OptimismL1BridgeAdapterEncoder.contract.Call(opts, &out, "encodeFinalizeWithdrawalERC20Payload", payload)

	if err != nil {
		return err
	}

	return err

}

func (_OptimismL1BridgeAdapterEncoder *OptimismL1BridgeAdapterEncoderSession) EncodeFinalizeWithdrawalERC20Payload(payload OptimismL1BridgeAdapterFinalizeWithdrawERC20Payload) error {
	return _OptimismL1BridgeAdapterEncoder.Contract.EncodeFinalizeWithdrawalERC20Payload(&_OptimismL1BridgeAdapterEncoder.CallOpts, payload)
}

func (_OptimismL1BridgeAdapterEncoder *OptimismL1BridgeAdapterEncoderCallerSession) EncodeFinalizeWithdrawalERC20Payload(payload OptimismL1BridgeAdapterFinalizeWithdrawERC20Payload) error {
	return _OptimismL1BridgeAdapterEncoder.Contract.EncodeFinalizeWithdrawalERC20Payload(&_OptimismL1BridgeAdapterEncoder.CallOpts, payload)
}

func (_OptimismL1BridgeAdapterEncoder *OptimismL1BridgeAdapterEncoderCaller) EncodeOptimismFinalizationPayload(opts *bind.CallOpts, payload OptimismL1BridgeAdapterOptimismFinalizationPayload) error {
	var out []interface{}
	err := _OptimismL1BridgeAdapterEncoder.contract.Call(opts, &out, "encodeOptimismFinalizationPayload", payload)

	if err != nil {
		return err
	}

	return err

}

func (_OptimismL1BridgeAdapterEncoder *OptimismL1BridgeAdapterEncoderSession) EncodeOptimismFinalizationPayload(payload OptimismL1BridgeAdapterOptimismFinalizationPayload) error {
	return _OptimismL1BridgeAdapterEncoder.Contract.EncodeOptimismFinalizationPayload(&_OptimismL1BridgeAdapterEncoder.CallOpts, payload)
}

func (_OptimismL1BridgeAdapterEncoder *OptimismL1BridgeAdapterEncoderCallerSession) EncodeOptimismFinalizationPayload(payload OptimismL1BridgeAdapterOptimismFinalizationPayload) error {
	return _OptimismL1BridgeAdapterEncoder.Contract.EncodeOptimismFinalizationPayload(&_OptimismL1BridgeAdapterEncoder.CallOpts, payload)
}

func (_OptimismL1BridgeAdapterEncoder *OptimismL1BridgeAdapterEncoderCaller) EncodeOptimismProveWithdrawalPayload(opts *bind.CallOpts, payload OptimismL1BridgeAdapterOptimismProveWithdrawalPayload) error {
	var out []interface{}
	err := _OptimismL1BridgeAdapterEncoder.contract.Call(opts, &out, "encodeOptimismProveWithdrawalPayload", payload)

	if err != nil {
		return err
	}

	return err

}

func (_OptimismL1BridgeAdapterEncoder *OptimismL1BridgeAdapterEncoderSession) EncodeOptimismProveWithdrawalPayload(payload OptimismL1BridgeAdapterOptimismProveWithdrawalPayload) error {
	return _OptimismL1BridgeAdapterEncoder.Contract.EncodeOptimismProveWithdrawalPayload(&_OptimismL1BridgeAdapterEncoder.CallOpts, payload)
}

func (_OptimismL1BridgeAdapterEncoder *OptimismL1BridgeAdapterEncoderCallerSession) EncodeOptimismProveWithdrawalPayload(payload OptimismL1BridgeAdapterOptimismProveWithdrawalPayload) error {
	return _OptimismL1BridgeAdapterEncoder.Contract.EncodeOptimismProveWithdrawalPayload(&_OptimismL1BridgeAdapterEncoder.CallOpts, payload)
}

func (_OptimismL1BridgeAdapterEncoder *OptimismL1BridgeAdapterEncoder) Address() common.Address {
	return _OptimismL1BridgeAdapterEncoder.address
}

type OptimismL1BridgeAdapterEncoderInterface interface {
	EncodeFinalizeWithdrawalERC20Payload(opts *bind.CallOpts, payload OptimismL1BridgeAdapterFinalizeWithdrawERC20Payload) error

	EncodeOptimismFinalizationPayload(opts *bind.CallOpts, payload OptimismL1BridgeAdapterOptimismFinalizationPayload) error

	EncodeOptimismProveWithdrawalPayload(opts *bind.CallOpts, payload OptimismL1BridgeAdapterOptimismProveWithdrawalPayload) error

	Address() common.Address
}
