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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"entryPoint\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"factory\",\"type\":\"address\"}],\"name\":\"calculateSmartContractAccountAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"endContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"getFullEndTxEncoding\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"encoding\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"userOpHash\",\"type\":\"bytes32\"}],\"name\":\"getFullHashForSigning\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"factory\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"entryPoint\",\"type\":\"address\"}],\"name\":\"getInitCode\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"entryPoint\",\"type\":\"address\"}],\"name\":\"getSCAInitCodeWithConstructor\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x61124861003a600b82828239805160001a60731461002d57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600436106100715760003560e01c8063e0237bef1161005a578063e0237bef14610196578063e464b363146101ce578063fc59bac3146101e157600080fd5b80633ef1eb79146100765780634b770f5614610176575b600080fd5b610163610084366004610686565b604080517f4750045d47fce615521b32cee713ff8db50147e98aec5ca94926b52651ca3fa0602080830191909152818301939093528151808203830181526060820183528051908401207f190000000000000000000000000000000000000000000000000000000000000060808301527f010000000000000000000000000000000000000000000000000000000000000060818301527f1c7d3b72b37a35523e273aaadd7b4cd66f618bb81429ab053412d51f50ccea61608283015260a2808301919091528251808303909101815260c2909101909152805191012090565b6040519081526020015b60405180910390f35b6101896101843660046106c8565b6101f4565b60405161016d9190610785565b6101a96101a43660046106c8565b610398565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161016d565b6101896101dc36600461079f565b610523565b6101896101ef366004610801565b6105e1565b6040516060907fffffffffffffffffffffffffffffffffffffffff00000000000000000000000084831b169060009061022f60208201610679565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe082820381018352601f90910116604081815273ffffffffffffffffffffffffffffffffffffffff8881166020840152871690820152606001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290526102c492916020016108f4565b60405160208183030381529060405290508560601b630af4926f60e01b83836040516024016102f4929190610923565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152918152602080830180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909516949094179093525161037e939201610944565b604051602081830303815290604052925050509392505050565b6040516000907fffffffffffffffffffffffffffffffffffffffff000000000000000000000000606086901b169082906103d460208201610679565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe082820381018352601f90910116604081815273ffffffffffffffffffffffffffffffffffffffff8981166020840152881690820152606001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529082905261046992916020016108f4565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001815282825280516020918201207fff000000000000000000000000000000000000000000000000000000000000008285015260609790971b7fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166021840152603583019490945260558083019690965280518083039096018652607590910190525082519201919091209392505050565b60606040518060200161053590610679565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe082820381018352601f90910116604081815273ffffffffffffffffffffffffffffffffffffffff8681166020840152851690820152606001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290526105ca92916020016108f4565b604051602081830303815290604052905092915050565b60607ff3dd22a5000000000000000000000000000000000000000000000000000000008585610610864261098c565b8560405160200161062494939291906109cb565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290526106609291602001610a10565b6040516020818303038152906040529050949350505050565b6107e380610a5983390190565b60006020828403121561069857600080fd5b5035919050565b803573ffffffffffffffffffffffffffffffffffffffff811681146106c357600080fd5b919050565b6000806000606084860312156106dd57600080fd5b6106e68461069f565b92506106f46020850161069f565b91506107026040850161069f565b90509250925092565b60005b8381101561072657818101518382015260200161070e565b83811115610735576000848401525b50505050565b6000815180845261075381602086016020860161070b565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000610798602083018461073b565b9392505050565b600080604083850312156107b257600080fd5b6107bb8361069f565b91506107c96020840161069f565b90509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6000806000806080858703121561081757600080fd5b6108208561069f565b93506020850135925060408501359150606085013567ffffffffffffffff8082111561084b57600080fd5b818701915087601f83011261085f57600080fd5b813581811115610871576108716107d2565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f011681019083821181831017156108b7576108b76107d2565b816040528281528a60208487010111156108d057600080fd5b82602086016020830137600060208483010152809550505050505092959194509250565b6000835161090681846020880161070b565b83519083019061091a81836020880161070b565b01949350505050565b82815260406020820152600061093c604083018461073b565b949350505050565b7fffffffffffffffffffffffffffffffffffffffff000000000000000000000000831681526000825161097e81601485016020870161070b565b919091016014019392505050565b600082198211156109c6577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b500190565b73ffffffffffffffffffffffffffffffffffffffff85168152836020820152826040820152608060608201526000610a06608083018461073b565b9695505050505050565b7fffffffff000000000000000000000000000000000000000000000000000000008316815260008251610a4a81600485016020870161070b565b91909101600401939250505056fe60c060405234801561001057600080fd5b506040516107e33803806107e383398101604081905261002f91610062565b6001600160a01b039182166080521660a052610095565b80516001600160a01b038116811461005d57600080fd5b919050565b6000806040838503121561007557600080fd5b61007e83610046565b915061008c60208401610046565b90509250929050565b60805160a05161071d6100c6600039600081816056015261038901526000818160c80152610286015261071d6000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80631a4e75de146100515780633a871cdd146100a2578063e3978240146100c3578063f3dd22a5146100ea575b600080fd5b6100787f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b6100b56100b03660046104d6565b6100ff565b604051908152602001610099565b6100787f000000000000000000000000000000000000000000000000000000000000000081565b6100fd6100f836600461052a565b610371565b005b600080548460200135146101215761011a600160008061049e565b905061036a565b604080517f4750045d47fce615521b32cee713ff8db50147e98aec5ca94926b52651ca3fa06020808301919091528183018690528251808303840181526060830184528051908201207f190000000000000000000000000000000000000000000000000000000000000060808401527f010000000000000000000000000000000000000000000000000000000000000060818401527f1c7d3b72b37a35523e273aaadd7b4cd66f618bb81429ab053412d51f50ccea61608284015260a2808401919091528351808403909101815260c29092019092528051910120600061020c6101408701876105d9565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920182905250602085015160408087015187519798509196919550919350869250811061026557610265610645565b016020015160f81c905073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000166001866102b484601b6106a3565b6040805160008152602081018083529390935260ff90911690820152606081018690526080810185905260a0016020604051602081039080840390855afa158015610303573d6000803e3d6000fd5b5050506020604051035173ffffffffffffffffffffffffffffffffffffffff161461034157610335600160008061049e565b9550505050505061036a565b600080549080610350836106c8565b9190505550610362600080600061049e565b955050505050505b9392505050565b3373ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016146103e7576040517f4a0bfec10000000000000000000000000000000000000000000000000000000081523360048201526024015b60405180910390fd5b8242111561042a576040517f300249d7000000000000000000000000000000000000000000000000000000008152600481018490524260248201526044016103de565b8473ffffffffffffffffffffffffffffffffffffffff16848383604051610452929190610700565b60006040518083038185875af1925050503d806000811461048f576040519150601f19603f3d011682016040523d82523d6000602084013e610494565b606091505b5050505050505050565b600060d08265ffffffffffff16901b60a08465ffffffffffff16901b856104c65760006104c9565b60015b60ff161717949350505050565b6000806000606084860312156104eb57600080fd5b833567ffffffffffffffff81111561050257600080fd5b8401610160818703121561051557600080fd5b95602085013595506040909401359392505050565b60008060008060006080868803121561054257600080fd5b853573ffffffffffffffffffffffffffffffffffffffff8116811461056657600080fd5b94506020860135935060408601359250606086013567ffffffffffffffff8082111561059157600080fd5b818801915088601f8301126105a557600080fd5b8135818111156105b457600080fd5b8960208285010111156105c657600080fd5b9699959850939650602001949392505050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261060e57600080fd5b83018035915067ffffffffffffffff82111561062957600080fd5b60200191503681900382131561063e57600080fd5b9250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600060ff821660ff84168060ff038211156106c0576106c0610674565b019392505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036106f9576106f9610674565b5060010190565b818382376000910190815291905056fea164736f6c634300080f000aa164736f6c634300080f000a",
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

func (_SmartContractAccountHelper *SmartContractAccountHelperCaller) GetFullHashForSigning(opts *bind.CallOpts, userOpHash [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _SmartContractAccountHelper.contract.Call(opts, &out, "getFullHashForSigning", userOpHash)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_SmartContractAccountHelper *SmartContractAccountHelperSession) GetFullHashForSigning(userOpHash [32]byte) ([32]byte, error) {
	return _SmartContractAccountHelper.Contract.GetFullHashForSigning(&_SmartContractAccountHelper.CallOpts, userOpHash)
}

func (_SmartContractAccountHelper *SmartContractAccountHelperCallerSession) GetFullHashForSigning(userOpHash [32]byte) ([32]byte, error) {
	return _SmartContractAccountHelper.Contract.GetFullHashForSigning(&_SmartContractAccountHelper.CallOpts, userOpHash)
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

	GetFullHashForSigning(opts *bind.CallOpts, userOpHash [32]byte) ([32]byte, error)

	GetInitCode(opts *bind.CallOpts, factory common.Address, owner common.Address, entryPoint common.Address) ([]byte, error)

	GetSCAInitCodeWithConstructor(opts *bind.CallOpts, owner common.Address, entryPoint common.Address) ([]byte, error)

	Address() common.Address
}
