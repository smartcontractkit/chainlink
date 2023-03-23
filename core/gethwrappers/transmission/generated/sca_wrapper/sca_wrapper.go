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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"entryPoint\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"operationHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"InvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"NotAuthorized\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"currentTimestamp\",\"type\":\"uint256\"}],\"name\":\"TransactionExpired\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"executeTransactionFromEntryPoint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_entryPoint\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_nonce\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"callGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"verificationGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"preVerificationGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeePerGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxPriorityFeePerGas\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"paymasterAndData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structUserOperation\",\"name\":\"userOp\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"userOpHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"validateUserOp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"validationData\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60c060405234801561001057600080fd5b5060405161080738038061080783398101604081905261002f91610062565b6001600160a01b039182166080521660a052610095565b80516001600160a01b038116811461005d57600080fd5b919050565b6000806040838503121561007557600080fd5b61007e83610046565b915061008c60208401610046565b90509250929050565b60805160a0516107416100c660003960008181607101526103ad01526000818160ec01526102aa01526107416000f3fe608060405234801561001057600080fd5b50600436106100675760003560e01c80637eccf63e116100505780637eccf63e146100de578063e3978240146100e7578063f3dd22a51461010e57600080fd5b80631a4e75de1461006c5780633a871cdd146100bd575b600080fd5b6100937f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b6100d06100cb3660046104fa565b610123565b6040519081526020016100b4565b6100d060005481565b6100937f000000000000000000000000000000000000000000000000000000000000000081565b61012161011c36600461054e565b610395565b005b600080548460200135146101455761013e60016000806104c2565b905061038e565b604080517f4750045d47fce615521b32cee713ff8db50147e98aec5ca94926b52651ca3fa06020808301919091528183018690528251808303840181526060830184528051908201207f190000000000000000000000000000000000000000000000000000000000000060808401527f010000000000000000000000000000000000000000000000000000000000000060818401527f1c7d3b72b37a35523e273aaadd7b4cd66f618bb81429ab053412d51f50ccea61608284015260a2808401919091528351808403909101815260c2909201909252805191012060006102306101408701876105fd565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920182905250602085015160408087015187519798509196919550919350869250811061028957610289610669565b016020015160f81c905073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000166001866102d884601b6106c7565b6040805160008152602081018083529390935260ff90911690820152606081018690526080810185905260a0016020604051602081039080840390855afa158015610327573d6000803e3d6000fd5b5050506020604051035173ffffffffffffffffffffffffffffffffffffffff16146103655761035960016000806104c2565b9550505050505061038e565b600080549080610374836106ec565b919050555061038660008060006104c2565b955050505050505b9392505050565b3373ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161461040b576040517f4a0bfec10000000000000000000000000000000000000000000000000000000081523360048201526024015b60405180910390fd5b8242111561044e576040517f300249d700000000000000000000000000000000000000000000000000000000815260048101849052426024820152604401610402565b8473ffffffffffffffffffffffffffffffffffffffff16848383604051610476929190610724565b60006040518083038185875af1925050503d80600081146104b3576040519150601f19603f3d011682016040523d82523d6000602084013e6104b8565b606091505b5050505050505050565b600060d08265ffffffffffff16901b60a08465ffffffffffff16901b856104ea5760006104ed565b60015b60ff161717949350505050565b60008060006060848603121561050f57600080fd5b833567ffffffffffffffff81111561052657600080fd5b8401610160818703121561053957600080fd5b95602085013595506040909401359392505050565b60008060008060006080868803121561056657600080fd5b853573ffffffffffffffffffffffffffffffffffffffff8116811461058a57600080fd5b94506020860135935060408601359250606086013567ffffffffffffffff808211156105b557600080fd5b818801915088601f8301126105c957600080fd5b8135818111156105d857600080fd5b8960208285010111156105ea57600080fd5b9699959850939650602001949392505050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261063257600080fd5b83018035915067ffffffffffffffff82111561064d57600080fd5b60200191503681900382131561066257600080fd5b9250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600060ff821660ff84168060ff038211156106e4576106e4610698565b019392505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361071d5761071d610698565b5060010190565b818382376000910190815291905056fea164736f6c634300080f000a",
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

func (_SCA *SCACaller) SEntryPoint(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SCA.contract.Call(opts, &out, "s_entryPoint")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_SCA *SCASession) SEntryPoint() (common.Address, error) {
	return _SCA.Contract.SEntryPoint(&_SCA.CallOpts)
}

func (_SCA *SCACallerSession) SEntryPoint() (common.Address, error) {
	return _SCA.Contract.SEntryPoint(&_SCA.CallOpts)
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

func (_SCA *SCACaller) SOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SCA.contract.Call(opts, &out, "s_owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_SCA *SCASession) SOwner() (common.Address, error) {
	return _SCA.Contract.SOwner(&_SCA.CallOpts)
}

func (_SCA *SCACallerSession) SOwner() (common.Address, error) {
	return _SCA.Contract.SOwner(&_SCA.CallOpts)
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
	SEntryPoint(opts *bind.CallOpts) (common.Address, error)

	SNonce(opts *bind.CallOpts) (*big.Int, error)

	SOwner(opts *bind.CallOpts) (common.Address, error)

	ExecuteTransactionFromEntryPoint(opts *bind.TransactOpts, to common.Address, value *big.Int, deadline *big.Int, data []byte) (*types.Transaction, error)

	ValidateUserOp(opts *bind.TransactOpts, userOp UserOperation, userOpHash [32]byte, arg2 *big.Int) (*types.Transaction, error)

	Address() common.Address
}
