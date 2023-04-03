// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package sca_wrapper

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

type UserOperation struct {
	Sender               common.Address
	Nonce                *big.Int
	InitCode             []byte
	CallData             []byte
	CallGasLimit         *big.Int
	VerificationGasLimit *big.Int
	PreVerificationGas   *big.Int
	MaxFeePerGas         *big.Int
	MaxPriorityFeePerGas *big.Int
	PaymasterAndData     []byte
	Signature            []byte
}

var SCAMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"entryPoint\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"currentNonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonceGiven\",\"type\":\"uint256\"}],\"name\":\"IncorrectNonce\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"operationHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"InvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"NotAuthorized\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"currentTimestamp\",\"type\":\"uint256\"}],\"name\":\"TransactionExpired\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint48\",\"name\":\"deadline\",\"type\":\"uint48\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"executeTransactionFromEntryPoint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_entryPoint\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_nonce\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"callGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"verificationGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"preVerificationGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeePerGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxPriorityFeePerGas\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"paymasterAndData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structUserOperation\",\"name\":\"userOp\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"userOpHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"validateUserOp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"validationData\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60c060405234801561001057600080fd5b50604051610a4b380380610a4b83398101604081905261002f91610062565b6001600160a01b039182166080521660a052610095565b80516001600160a01b038116811461005d57600080fd5b919050565b6000806040838503121561007557600080fd5b61007e83610046565b915061008c60208401610046565b90509250929050565b60805160a0516109846100c76000396000818160710152610450015260008181610101015261031c01526109846000f3fe608060405234801561001057600080fd5b50600436106100675760003560e01c80637eccf63e116100505780637eccf63e146100de57806389553be4146100e7578063dba6335f146100fc57600080fd5b8063140fcfb11461006c5780633a871cdd146100bd575b600080fd5b6100937f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b6100d06100cb3660046105ba565b610123565b6040519081526020016100b4565b6100d060005481565b6100fa6100f536600461064e565b610438565b005b6100937f000000000000000000000000000000000000000000000000000000000000000081565b60008054846020013514610179576000546040517f7ba633940000000000000000000000000000000000000000000000000000000081526004810191909152602085013560248201526044015b60405180910390fd5b60006102908430604080517f4750045d47fce615521b32cee713ff8db50147e98aec5ca94926b52651ca3fa060208083019190915281830194909452815180820383018152606080830184528151918601919091207f190000000000000000000000000000000000000000000000000000000000000060808401527f010000000000000000000000000000000000000000000000000000000000000060818401527f1c7d3b72b37a35523e273aaadd7b4cd66f618bb81429ab053412d51f50ccea6160828401524660a284015293901b7fffffffffffffffffffffffffffffffffffffffff0000000000000000000000001660c282015260d6808201939093528151808203909301835260f6019052805191012090565b905060006102a26101408701876106eb565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092018290525060208501516040808701518751979850919691955091935086925081106102fb576102fb610757565b016020015160f81c905073ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001660018661034a84601b6107b5565b6040805160008152602081018083529390935260ff90911690820152606081018690526080810185905260a0016020604051602081039080840390855afa158015610399573d6000803e3d6000fd5b5050506020604051035173ffffffffffffffffffffffffffffffffffffffff16146103d7576103cb6001600080610582565b95505050505050610431565b6000805490806103e6836107da565b90915550600090506103fb60608b018b6106eb565b610409916004908290610812565b810190610416919061086b565b50925050506104286000826000610582565b96505050505050505b9392505050565b3373ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016146104a9576040517f4a0bfec1000000000000000000000000000000000000000000000000000000008152336004820152602401610170565b65ffffffffffff8316158015906104c757508265ffffffffffff1642115b1561050e576040517f300249d700000000000000000000000000000000000000000000000000000000815265ffffffffffff84166004820152426024820152604401610170565b8473ffffffffffffffffffffffffffffffffffffffff16848383604051610536929190610967565b60006040518083038185875af1925050503d8060008114610573576040519150601f19603f3d011682016040523d82523d6000602084013e610578565b606091505b5050505050505050565b600060d08265ffffffffffff16901b60a08465ffffffffffff16901b856105aa5760006105ad565b60015b60ff161717949350505050565b6000806000606084860312156105cf57600080fd5b833567ffffffffffffffff8111156105e657600080fd5b840161016081870312156105f957600080fd5b95602085013595506040909401359392505050565b73ffffffffffffffffffffffffffffffffffffffff8116811461063057600080fd5b50565b803565ffffffffffff8116811461064957600080fd5b919050565b60008060008060006080868803121561066657600080fd5b85356106718161060e565b94506020860135935061068660408701610633565b9250606086013567ffffffffffffffff808211156106a357600080fd5b818801915088601f8301126106b757600080fd5b8135818111156106c657600080fd5b8960208285010111156106d857600080fd5b9699959850939650602001949392505050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261072057600080fd5b83018035915067ffffffffffffffff82111561073b57600080fd5b60200191503681900382131561075057600080fd5b9250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600060ff821660ff84168060ff038211156107d2576107d2610786565b019392505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361080b5761080b610786565b5060010190565b6000808585111561082257600080fd5b8386111561082f57600080fd5b5050820193919092039150565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6000806000806080858703121561088157600080fd5b843561088c8161060e565b9350602085013592506108a160408601610633565b9150606085013567ffffffffffffffff808211156108be57600080fd5b818701915087601f8301126108d257600080fd5b8135818111156108e4576108e461083c565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f0116810190838211818310171561092a5761092a61083c565b816040528281528a602084870101111561094357600080fd5b82602086016020830137600060208483010152809550505050505092959194509250565b818382376000910190815291905056fea164736f6c634300080f000a",
}

var SCAABI = SCAMetaData.ABI

var SCABin = SCAMetaData.Bin

func DeploySCA(auth *bind.TransactOpts, backend bind.ContractBackend, owner common.Address, entryPoint common.Address) (common.Address, *types.Transaction, *SCA, error) {
	parsed, err := SCAMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SCABin), backend, owner, entryPoint)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SCA{SCACaller: SCACaller{contract: contract}, SCATransactor: SCATransactor{contract: contract}, SCAFilterer: SCAFilterer{contract: contract}}, nil
}

type SCA struct {
	address common.Address
	abi     abi.ABI
	SCACaller
	SCATransactor
	SCAFilterer
}

type SCACaller struct {
	contract *bind.BoundContract
}

type SCATransactor struct {
	contract *bind.BoundContract
}

type SCAFilterer struct {
	contract *bind.BoundContract
}

type SCASession struct {
	Contract     *SCA
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type SCACallerSession struct {
	Contract *SCACaller
	CallOpts bind.CallOpts
}

type SCATransactorSession struct {
	Contract     *SCATransactor
	TransactOpts bind.TransactOpts
}

type SCARaw struct {
	Contract *SCA
}

type SCACallerRaw struct {
	Contract *SCACaller
}

type SCATransactorRaw struct {
	Contract *SCATransactor
}

func NewSCA(address common.Address, backend bind.ContractBackend) (*SCA, error) {
	abi, err := abi.JSON(strings.NewReader(SCAABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindSCA(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SCA{address: address, abi: abi, SCACaller: SCACaller{contract: contract}, SCATransactor: SCATransactor{contract: contract}, SCAFilterer: SCAFilterer{contract: contract}}, nil
}

func NewSCACaller(address common.Address, caller bind.ContractCaller) (*SCACaller, error) {
	contract, err := bindSCA(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SCACaller{contract: contract}, nil
}

func NewSCATransactor(address common.Address, transactor bind.ContractTransactor) (*SCATransactor, error) {
	contract, err := bindSCA(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SCATransactor{contract: contract}, nil
}

func NewSCAFilterer(address common.Address, filterer bind.ContractFilterer) (*SCAFilterer, error) {
	contract, err := bindSCA(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SCAFilterer{contract: contract}, nil
}

func bindSCA(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SCAABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_SCA *SCARaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SCA.Contract.SCACaller.contract.Call(opts, result, method, params...)
}

func (_SCA *SCARaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SCA.Contract.SCATransactor.contract.Transfer(opts)
}

func (_SCA *SCARaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SCA.Contract.SCATransactor.contract.Transact(opts, method, params...)
}

func (_SCA *SCACallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SCA.Contract.contract.Call(opts, result, method, params...)
}

func (_SCA *SCATransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SCA.Contract.contract.Transfer(opts)
}

func (_SCA *SCATransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SCA.Contract.contract.Transact(opts, method, params...)
}

func (_SCA *SCACaller) IEntryPoint(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SCA.contract.Call(opts, &out, "i_entryPoint")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_SCA *SCASession) IEntryPoint() (common.Address, error) {
	return _SCA.Contract.IEntryPoint(&_SCA.CallOpts)
}

func (_SCA *SCACallerSession) IEntryPoint() (common.Address, error) {
	return _SCA.Contract.IEntryPoint(&_SCA.CallOpts)
}

func (_SCA *SCACaller) IOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SCA.contract.Call(opts, &out, "i_owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_SCA *SCASession) IOwner() (common.Address, error) {
	return _SCA.Contract.IOwner(&_SCA.CallOpts)
}

func (_SCA *SCACallerSession) IOwner() (common.Address, error) {
	return _SCA.Contract.IOwner(&_SCA.CallOpts)
}

func (_SCA *SCACaller) SNonce(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SCA.contract.Call(opts, &out, "s_nonce")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_SCA *SCASession) SNonce() (*big.Int, error) {
	return _SCA.Contract.SNonce(&_SCA.CallOpts)
}

func (_SCA *SCACallerSession) SNonce() (*big.Int, error) {
	return _SCA.Contract.SNonce(&_SCA.CallOpts)
}

func (_SCA *SCATransactor) ExecuteTransactionFromEntryPoint(opts *bind.TransactOpts, to common.Address, value *big.Int, deadline *big.Int, data []byte) (*types.Transaction, error) {
	return _SCA.contract.Transact(opts, "executeTransactionFromEntryPoint", to, value, deadline, data)
}

func (_SCA *SCASession) ExecuteTransactionFromEntryPoint(to common.Address, value *big.Int, deadline *big.Int, data []byte) (*types.Transaction, error) {
	return _SCA.Contract.ExecuteTransactionFromEntryPoint(&_SCA.TransactOpts, to, value, deadline, data)
}

func (_SCA *SCATransactorSession) ExecuteTransactionFromEntryPoint(to common.Address, value *big.Int, deadline *big.Int, data []byte) (*types.Transaction, error) {
	return _SCA.Contract.ExecuteTransactionFromEntryPoint(&_SCA.TransactOpts, to, value, deadline, data)
}

func (_SCA *SCATransactor) ValidateUserOp(opts *bind.TransactOpts, userOp UserOperation, userOpHash [32]byte, arg2 *big.Int) (*types.Transaction, error) {
	return _SCA.contract.Transact(opts, "validateUserOp", userOp, userOpHash, arg2)
}

func (_SCA *SCASession) ValidateUserOp(userOp UserOperation, userOpHash [32]byte, arg2 *big.Int) (*types.Transaction, error) {
	return _SCA.Contract.ValidateUserOp(&_SCA.TransactOpts, userOp, userOpHash, arg2)
}

func (_SCA *SCATransactorSession) ValidateUserOp(userOp UserOperation, userOpHash [32]byte, arg2 *big.Int) (*types.Transaction, error) {
	return _SCA.Contract.ValidateUserOp(&_SCA.TransactOpts, userOp, userOpHash, arg2)
}

func (_SCA *SCA) Address() common.Address {
	return _SCA.address
}

type SCAInterface interface {
	IEntryPoint(opts *bind.CallOpts) (common.Address, error)

	IOwner(opts *bind.CallOpts) (common.Address, error)

	SNonce(opts *bind.CallOpts) (*big.Int, error)

	ExecuteTransactionFromEntryPoint(opts *bind.TransactOpts, to common.Address, value *big.Int, deadline *big.Int, data []byte) (*types.Transaction, error)

	ValidateUserOp(opts *bind.TransactOpts, userOp UserOperation, userOpHash [32]byte, arg2 *big.Int) (*types.Transaction, error)

	Address() common.Address
}
