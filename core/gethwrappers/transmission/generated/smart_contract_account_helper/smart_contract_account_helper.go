// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package smart_contract_account_helper

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

var SmartContractAccountHelperMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"entryPoint\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"factory\",\"type\":\"address\"}],\"name\":\"calculateSmartContractAccountAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"endContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"getFullEndTxEncoding\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"encoding\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"fullEndTxEncoding\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"getFullHashForSigning\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"fullHash\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"factory\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"entryPoint\",\"type\":\"address\"}],\"name\":\"getInitCode\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"entryPoint\",\"type\":\"address\"}],\"name\":\"getSCAInitCodeWithConstructor\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x6113fe61003a600b82828239805160001a60731461002d57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600436106100715760003560e01c8063e0237bef1161005a578063e0237bef146100c0578063e464b363146100f8578063fc59bac31461010b57600080fd5b80634b770f5614610076578063a2c370b01461009f575b600080fd5b6100896100843660046106e5565b61011e565b60405161009691906107a2565b60405180910390f35b6100b26100ad366004610896565b6102c2565b604051908152602001610096565b6100d36100ce3660046106e5565b6103ce565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610096565b6100896101063660046108ed565b610559565b610089610119366004610920565b610617565b6040516060907fffffffffffffffffffffffffffffffffffffffff00000000000000000000000084831b1690600090610159602082016106af565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe082820381018352601f90910116604081815273ffffffffffffffffffffffffffffffffffffffff8881166020840152871690820152606001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290526101ee9291602001610981565b60405160208183030381529060405290508560601b630af4926f60e01b838360405160240161021e9291906109b0565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152918152602080830180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090951694909417909352516102a89392016109d1565b604051602081830303815290604052925050509392505050565b6000807f4750045d47fce615521b32cee713ff8db50147e98aec5ca94926b52651ca3fa0858585466040516020016102fe959493929190610a19565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815282825280516020918201207f1900000000000000000000000000000000000000000000000000000000000000828501527f010000000000000000000000000000000000000000000000000000000000000060218501527f1c7d3b72b37a35523e273aaadd7b4cd66f618bb81429ab053412d51f50ccea616022850152604280850191909152825180850390910181526062909301909152815191012095945050505050565b6040516000907fffffffffffffffffffffffffffffffffffffffff000000000000000000000000606086901b1690829061040a602082016106af565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe082820381018352601f90910116604081815273ffffffffffffffffffffffffffffffffffffffff8981166020840152881690820152606001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529082905261049f9291602001610981565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001815282825280516020918201207fff000000000000000000000000000000000000000000000000000000000000008285015260609790971b7fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166021840152603583019490945260558083019690965280518083039096018652607590910190525082519201919091209392505050565b60606040518060200161056b906106af565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe082820381018352601f90910116604081815273ffffffffffffffffffffffffffffffffffffffff8681166020840152851690820152606001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290526106009291602001610981565b604051602081830303815290604052905092915050565b60607ff3dd22a50000000000000000000000000000000000000000000000000000000085856106468642610a65565b8560405160200161065a9493929190610aa4565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290526106969291602001610ae9565b6040516020818303038152906040529050949350505050565b6108c080610b3283390190565b803573ffffffffffffffffffffffffffffffffffffffff811681146106e057600080fd5b919050565b6000806000606084860312156106fa57600080fd5b610703846106bc565b9250610711602085016106bc565b915061071f604085016106bc565b90509250925092565b60005b8381101561074357818101518382015260200161072b565b83811115610752576000848401525b50505050565b60008151808452610770816020860160208601610728565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6020815260006107b56020830184610758565b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600082601f8301126107fc57600080fd5b813567ffffffffffffffff80821115610817576108176107bc565b604051601f83017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f0116810190828211818310171561085d5761085d6107bc565b8160405283815286602085880101111561087657600080fd5b836020870160208301376000602085830101528094505050505092915050565b6000806000606084860312156108ab57600080fd5b833567ffffffffffffffff8111156108c257600080fd5b6108ce868287016107eb565b9350506108dd602085016106bc565b9150604084013590509250925092565b6000806040838503121561090057600080fd5b610909836106bc565b9150610917602084016106bc565b90509250929050565b6000806000806080858703121561093657600080fd5b61093f856106bc565b93506020850135925060408501359150606085013567ffffffffffffffff81111561096957600080fd5b610975878288016107eb565b91505092959194509250565b60008351610993818460208801610728565b8351908301906109a7818360208801610728565b01949350505050565b8281526040602082015260006109c96040830184610758565b949350505050565b7fffffffffffffffffffffffffffffffffffffffff0000000000000000000000008316815260008251610a0b816014850160208701610728565b919091016014019392505050565b85815260a060208201526000610a3260a0830187610758565b73ffffffffffffffffffffffffffffffffffffffff95909516604083015250606081019290925260809091015292915050565b60008219821115610a9f577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b500190565b73ffffffffffffffffffffffffffffffffffffffff85168152836020820152826040820152608060608201526000610adf6080830184610758565b9695505050505050565b7fffffffff000000000000000000000000000000000000000000000000000000008316815260008251610b23816004850160208701610728565b91909101600401939250505056fe60c060405234801561001057600080fd5b506040516108c03803806108c083398101604081905261002f91610062565b6001600160a01b039182166080521660a052610095565b80516001600160a01b038116811461005d57600080fd5b919050565b6000806040838503121561007557600080fd5b61007e83610046565b915061008c60208401610046565b90509250929050565b60805160a0516107ec6100d4600039600081816056015261041301526000818160c80152818161013f015281816102bf01526103a501526107ec6000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80631a4e75de146100515780633a871cdd146100a2578063e3978240146100c3578063f3dd22a5146100ea575b600080fd5b6100787f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b6100b56100b0366004610523565b6100ff565b604051908152602001610099565b6100787f000000000000000000000000000000000000000000000000000000000000000081565b6100fd6100f8366004610577565b6103fb565b005b6000807f4750045d47fce615521b32cee713ff8db50147e98aec5ca94926b52651ca3fa06101306060870187610626565b600054604051610169949392917f0000000000000000000000000000000000000000000000000000000000000000914690602001610692565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815282825280516020918201207f1900000000000000000000000000000000000000000000000000000000000000828501527f010000000000000000000000000000000000000000000000000000000000000060218501527f1c7d3b72b37a35523e273aaadd7b4cd66f618bb81429ab053412d51f50ccea616022850152604280850182905283518086039091018152606290940190925282519201919091209091506000610245610140880188610626565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920182905250602085015160408087015187519798509196919550919350869250811061029e5761029e610714565b016020015160f81c905073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000166001866102ed84601b610772565b6040805160008152602081018083529390935260ff90911690820152606081018690526080810185905260a0016020604051602081039080840390855afa15801561033c573d6000803e3d6000fd5b5050506020604051035173ffffffffffffffffffffffffffffffffffffffff16146103d7576040517f78d50c7f0000000000000000000000000000000000000000000000000000000081526004810186905273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001660248201526044015b60405180910390fd5b6000805490806103e683610797565b9091555060009b9a5050505050505050505050565b3373ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161461046c576040517f4a0bfec10000000000000000000000000000000000000000000000000000000081523360048201526024016103ce565b824211156104af576040517f300249d7000000000000000000000000000000000000000000000000000000008152600481018490524260248201526044016103ce565b8473ffffffffffffffffffffffffffffffffffffffff168483836040516104d79291906107cf565b60006040518083038185875af1925050503d8060008114610514576040519150601f19603f3d011682016040523d82523d6000602084013e610519565b606091505b5050505050505050565b60008060006060848603121561053857600080fd5b833567ffffffffffffffff81111561054f57600080fd5b8401610160818703121561056257600080fd5b95602085013595506040909401359392505050565b60008060008060006080868803121561058f57600080fd5b853573ffffffffffffffffffffffffffffffffffffffff811681146105b357600080fd5b94506020860135935060408601359250606086013567ffffffffffffffff808211156105de57600080fd5b818801915088601f8301126105f257600080fd5b81358181111561060157600080fd5b89602082850101111561061357600080fd5b9699959850939650602001949392505050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261065b57600080fd5b83018035915067ffffffffffffffff82111561067657600080fd5b60200191503681900382131561068b57600080fd5b9250929050565b86815260a060208201528460a0820152848660c0830137600060c08683010152600060c07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f880116830101905073ffffffffffffffffffffffffffffffffffffffff85166040830152836060830152826080830152979650505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600060ff821660ff84168060ff0382111561078f5761078f610743565b019392505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036107c8576107c8610743565b5060010190565b818382376000910190815291905056fea164736f6c634300080f000aa164736f6c634300080f000a",
}

var SmartContractAccountHelperABI = SmartContractAccountHelperMetaData.ABI

var SmartContractAccountHelperBin = SmartContractAccountHelperMetaData.Bin

func DeploySmartContractAccountHelper(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SmartContractAccountHelper, error) {
	parsed, err := SmartContractAccountHelperMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SmartContractAccountHelperBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SmartContractAccountHelper{SmartContractAccountHelperCaller: SmartContractAccountHelperCaller{contract: contract}, SmartContractAccountHelperTransactor: SmartContractAccountHelperTransactor{contract: contract}, SmartContractAccountHelperFilterer: SmartContractAccountHelperFilterer{contract: contract}}, nil
}

type SmartContractAccountHelper struct {
	address common.Address
	abi     abi.ABI
	SmartContractAccountHelperCaller
	SmartContractAccountHelperTransactor
	SmartContractAccountHelperFilterer
}

type SmartContractAccountHelperCaller struct {
	contract *bind.BoundContract
}

type SmartContractAccountHelperTransactor struct {
	contract *bind.BoundContract
}

type SmartContractAccountHelperFilterer struct {
	contract *bind.BoundContract
}

type SmartContractAccountHelperSession struct {
	Contract     *SmartContractAccountHelper
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type SmartContractAccountHelperCallerSession struct {
	Contract *SmartContractAccountHelperCaller
	CallOpts bind.CallOpts
}

type SmartContractAccountHelperTransactorSession struct {
	Contract     *SmartContractAccountHelperTransactor
	TransactOpts bind.TransactOpts
}

type SmartContractAccountHelperRaw struct {
	Contract *SmartContractAccountHelper
}

type SmartContractAccountHelperCallerRaw struct {
	Contract *SmartContractAccountHelperCaller
}

type SmartContractAccountHelperTransactorRaw struct {
	Contract *SmartContractAccountHelperTransactor
}

func NewSmartContractAccountHelper(address common.Address, backend bind.ContractBackend) (*SmartContractAccountHelper, error) {
	abi, err := abi.JSON(strings.NewReader(SmartContractAccountHelperABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindSmartContractAccountHelper(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SmartContractAccountHelper{address: address, abi: abi, SmartContractAccountHelperCaller: SmartContractAccountHelperCaller{contract: contract}, SmartContractAccountHelperTransactor: SmartContractAccountHelperTransactor{contract: contract}, SmartContractAccountHelperFilterer: SmartContractAccountHelperFilterer{contract: contract}}, nil
}

func NewSmartContractAccountHelperCaller(address common.Address, caller bind.ContractCaller) (*SmartContractAccountHelperCaller, error) {
	contract, err := bindSmartContractAccountHelper(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SmartContractAccountHelperCaller{contract: contract}, nil
}

func NewSmartContractAccountHelperTransactor(address common.Address, transactor bind.ContractTransactor) (*SmartContractAccountHelperTransactor, error) {
	contract, err := bindSmartContractAccountHelper(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SmartContractAccountHelperTransactor{contract: contract}, nil
}

func NewSmartContractAccountHelperFilterer(address common.Address, filterer bind.ContractFilterer) (*SmartContractAccountHelperFilterer, error) {
	contract, err := bindSmartContractAccountHelper(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SmartContractAccountHelperFilterer{contract: contract}, nil
}

func bindSmartContractAccountHelper(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SmartContractAccountHelperABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_SmartContractAccountHelper *SmartContractAccountHelperRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SmartContractAccountHelper.Contract.SmartContractAccountHelperCaller.contract.Call(opts, result, method, params...)
}

func (_SmartContractAccountHelper *SmartContractAccountHelperRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SmartContractAccountHelper.Contract.SmartContractAccountHelperTransactor.contract.Transfer(opts)
}

func (_SmartContractAccountHelper *SmartContractAccountHelperRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SmartContractAccountHelper.Contract.SmartContractAccountHelperTransactor.contract.Transact(opts, method, params...)
}

func (_SmartContractAccountHelper *SmartContractAccountHelperCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SmartContractAccountHelper.Contract.contract.Call(opts, result, method, params...)
}

func (_SmartContractAccountHelper *SmartContractAccountHelperTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SmartContractAccountHelper.Contract.contract.Transfer(opts)
}

func (_SmartContractAccountHelper *SmartContractAccountHelperTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SmartContractAccountHelper.Contract.contract.Transact(opts, method, params...)
}

func (_SmartContractAccountHelper *SmartContractAccountHelperCaller) CalculateSmartContractAccountAddress(opts *bind.CallOpts, owner common.Address, entryPoint common.Address, factory common.Address) (common.Address, error) {
	var out []interface{}
	err := _SmartContractAccountHelper.contract.Call(opts, &out, "calculateSmartContractAccountAddress", owner, entryPoint, factory)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_SmartContractAccountHelper *SmartContractAccountHelperSession) CalculateSmartContractAccountAddress(owner common.Address, entryPoint common.Address, factory common.Address) (common.Address, error) {
	return _SmartContractAccountHelper.Contract.CalculateSmartContractAccountAddress(&_SmartContractAccountHelper.CallOpts, owner, entryPoint, factory)
}

func (_SmartContractAccountHelper *SmartContractAccountHelperCallerSession) CalculateSmartContractAccountAddress(owner common.Address, entryPoint common.Address, factory common.Address) (common.Address, error) {
	return _SmartContractAccountHelper.Contract.CalculateSmartContractAccountAddress(&_SmartContractAccountHelper.CallOpts, owner, entryPoint, factory)
}

func (_SmartContractAccountHelper *SmartContractAccountHelperCaller) GetFullEndTxEncoding(opts *bind.CallOpts, endContract common.Address, value *big.Int, deadline *big.Int, data []byte) ([]byte, error) {
	var out []interface{}
	err := _SmartContractAccountHelper.contract.Call(opts, &out, "getFullEndTxEncoding", endContract, value, deadline, data)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_SmartContractAccountHelper *SmartContractAccountHelperSession) GetFullEndTxEncoding(endContract common.Address, value *big.Int, deadline *big.Int, data []byte) ([]byte, error) {
	return _SmartContractAccountHelper.Contract.GetFullEndTxEncoding(&_SmartContractAccountHelper.CallOpts, endContract, value, deadline, data)
}

func (_SmartContractAccountHelper *SmartContractAccountHelperCallerSession) GetFullEndTxEncoding(endContract common.Address, value *big.Int, deadline *big.Int, data []byte) ([]byte, error) {
	return _SmartContractAccountHelper.Contract.GetFullEndTxEncoding(&_SmartContractAccountHelper.CallOpts, endContract, value, deadline, data)
}

func (_SmartContractAccountHelper *SmartContractAccountHelperCaller) GetFullHashForSigning(opts *bind.CallOpts, fullEndTxEncoding []byte, owner common.Address, nonce *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _SmartContractAccountHelper.contract.Call(opts, &out, "getFullHashForSigning", fullEndTxEncoding, owner, nonce)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_SmartContractAccountHelper *SmartContractAccountHelperSession) GetFullHashForSigning(fullEndTxEncoding []byte, owner common.Address, nonce *big.Int) ([32]byte, error) {
	return _SmartContractAccountHelper.Contract.GetFullHashForSigning(&_SmartContractAccountHelper.CallOpts, fullEndTxEncoding, owner, nonce)
}

func (_SmartContractAccountHelper *SmartContractAccountHelperCallerSession) GetFullHashForSigning(fullEndTxEncoding []byte, owner common.Address, nonce *big.Int) ([32]byte, error) {
	return _SmartContractAccountHelper.Contract.GetFullHashForSigning(&_SmartContractAccountHelper.CallOpts, fullEndTxEncoding, owner, nonce)
}

func (_SmartContractAccountHelper *SmartContractAccountHelperCaller) GetInitCode(opts *bind.CallOpts, factory common.Address, owner common.Address, entryPoint common.Address) ([]byte, error) {
	var out []interface{}
	err := _SmartContractAccountHelper.contract.Call(opts, &out, "getInitCode", factory, owner, entryPoint)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_SmartContractAccountHelper *SmartContractAccountHelperSession) GetInitCode(factory common.Address, owner common.Address, entryPoint common.Address) ([]byte, error) {
	return _SmartContractAccountHelper.Contract.GetInitCode(&_SmartContractAccountHelper.CallOpts, factory, owner, entryPoint)
}

func (_SmartContractAccountHelper *SmartContractAccountHelperCallerSession) GetInitCode(factory common.Address, owner common.Address, entryPoint common.Address) ([]byte, error) {
	return _SmartContractAccountHelper.Contract.GetInitCode(&_SmartContractAccountHelper.CallOpts, factory, owner, entryPoint)
}

func (_SmartContractAccountHelper *SmartContractAccountHelperCaller) GetSCAInitCodeWithConstructor(opts *bind.CallOpts, owner common.Address, entryPoint common.Address) ([]byte, error) {
	var out []interface{}
	err := _SmartContractAccountHelper.contract.Call(opts, &out, "getSCAInitCodeWithConstructor", owner, entryPoint)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_SmartContractAccountHelper *SmartContractAccountHelperSession) GetSCAInitCodeWithConstructor(owner common.Address, entryPoint common.Address) ([]byte, error) {
	return _SmartContractAccountHelper.Contract.GetSCAInitCodeWithConstructor(&_SmartContractAccountHelper.CallOpts, owner, entryPoint)
}

func (_SmartContractAccountHelper *SmartContractAccountHelperCallerSession) GetSCAInitCodeWithConstructor(owner common.Address, entryPoint common.Address) ([]byte, error) {
	return _SmartContractAccountHelper.Contract.GetSCAInitCodeWithConstructor(&_SmartContractAccountHelper.CallOpts, owner, entryPoint)
}

func (_SmartContractAccountHelper *SmartContractAccountHelper) Address() common.Address {
	return _SmartContractAccountHelper.address
}

type SmartContractAccountHelperInterface interface {
	CalculateSmartContractAccountAddress(opts *bind.CallOpts, owner common.Address, entryPoint common.Address, factory common.Address) (common.Address, error)

	GetFullEndTxEncoding(opts *bind.CallOpts, endContract common.Address, value *big.Int, deadline *big.Int, data []byte) ([]byte, error)

	GetFullHashForSigning(opts *bind.CallOpts, fullEndTxEncoding []byte, owner common.Address, nonce *big.Int) ([32]byte, error)

	GetInitCode(opts *bind.CallOpts, factory common.Address, owner common.Address, entryPoint common.Address) ([]byte, error)

	GetSCAInitCodeWithConstructor(opts *bind.CallOpts, owner common.Address, entryPoint common.Address) ([]byte, error)

	Address() common.Address
}
