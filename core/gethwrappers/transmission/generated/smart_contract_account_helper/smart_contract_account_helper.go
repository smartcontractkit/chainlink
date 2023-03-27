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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"entryPoint\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"factory\",\"type\":\"address\"}],\"name\":\"calculateSmartContractAccountAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"topupThreshold\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"topupAmount\",\"type\":\"uint256\"}],\"name\":\"getAbiEncodedDirectRequestData\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"endContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"getFullEndTxEncoding\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"encoding\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"userOpHash\",\"type\":\"bytes32\"}],\"name\":\"getFullHashForSigning\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"factory\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"entryPoint\",\"type\":\"address\"}],\"name\":\"getInitCode\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"entryPoint\",\"type\":\"address\"}],\"name\":\"getSCAInitCodeWithConstructor\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x61134661003a600b82828239805160001a60731461002d57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe730000000000000000000000000000000000000000301460806040526004361061007c5760003560e01c8063e0237bef1161005a578063e0237bef1461020e578063e464b36314610246578063fc59bac31461025957600080fd5b80632c86cb35146100815780633ef1eb79146101005780634b770f56146101fb575b600080fd5b6100ea61008f366004610727565b604080516060808201835273ffffffffffffffffffffffffffffffffffffffff959095168082526020808301958652918301938452825191820152925183820152905182840152805180830390930183526080909101905290565b6040516100f791906107d4565b60405180910390f35b6101ed61010e3660046107ee565b604080517f4750045d47fce615521b32cee713ff8db50147e98aec5ca94926b52651ca3fa0602080830191909152818301939093528151808203830181526060820183528051908401207f190000000000000000000000000000000000000000000000000000000000000060808301527f010000000000000000000000000000000000000000000000000000000000000060818301527f1c7d3b72b37a35523e273aaadd7b4cd66f618bb81429ab053412d51f50ccea61608283015260a2808301919091528251808303909101815260c2909101909152805191012090565b6040519081526020016100f7565b6100ea610209366004610807565b61026c565b61022161021c366004610807565b610410565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100f7565b6100ea61025436600461084a565b61059b565b6100ea6102673660046108ac565b610659565b6040516060907fffffffffffffffffffffffffffffffffffffffff00000000000000000000000084831b16906000906102a7602082016106f1565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe082820381018352601f90910116604081815273ffffffffffffffffffffffffffffffffffffffff8881166020840152871690820152606001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529082905261033c929160200161099f565b60405160208183030381529060405290508560601b630af4926f60e01b838360405160240161036c9291906109ce565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152918152602080830180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090951694909417909352516103f69392016109ef565b604051602081830303815290604052925050509392505050565b6040516000907fffffffffffffffffffffffffffffffffffffffff000000000000000000000000606086901b1690829061044c602082016106f1565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe082820381018352601f90910116604081815273ffffffffffffffffffffffffffffffffffffffff8981166020840152881690820152606001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290526104e1929160200161099f565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001815282825280516020918201207fff000000000000000000000000000000000000000000000000000000000000008285015260609790971b7fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166021840152603583019490945260558083019690965280518083039096018652607590910190525082519201919091209392505050565b6060604051806020016105ad906106f1565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe082820381018352601f90910116604081815273ffffffffffffffffffffffffffffffffffffffff8681166020840152851690820152606001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815290829052610642929160200161099f565b604051602081830303815290604052905092915050565b60607ff3dd22a50000000000000000000000000000000000000000000000000000000085856106888642610a37565b8560405160200161069c9493929190610a76565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290526106d89291602001610abb565b6040516020818303038152906040529050949350505050565b61083680610b0483390190565b803573ffffffffffffffffffffffffffffffffffffffff8116811461072257600080fd5b919050565b60008060006060848603121561073c57600080fd5b610745846106fe565b95602085013595506040909401359392505050565b60005b8381101561077557818101518382015260200161075d565b83811115610784576000848401525b50505050565b600081518084526107a281602086016020860161075a565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6020815260006107e7602083018461078a565b9392505050565b60006020828403121561080057600080fd5b5035919050565b60008060006060848603121561081c57600080fd5b610825846106fe565b9250610833602085016106fe565b9150610841604085016106fe565b90509250925092565b6000806040838503121561085d57600080fd5b610866836106fe565b9150610874602084016106fe565b90509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600080600080608085870312156108c257600080fd5b6108cb856106fe565b93506020850135925060408501359150606085013567ffffffffffffffff808211156108f657600080fd5b818701915087601f83011261090a57600080fd5b81358181111561091c5761091c61087d565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f011681019083821181831017156109625761096261087d565b816040528281528a602084870101111561097b57600080fd5b82602086016020830137600060208483010152809550505050505092959194509250565b600083516109b181846020880161075a565b8351908301906109c581836020880161075a565b01949350505050565b8281526040602082015260006109e7604083018461078a565b949350505050565b7fffffffffffffffffffffffffffffffffffffffff0000000000000000000000008316815260008251610a2981601485016020870161075a565b919091016014019392505050565b60008219821115610a71577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b500190565b73ffffffffffffffffffffffffffffffffffffffff85168152836020820152826040820152608060608201526000610ab1608083018461078a565b9695505050505050565b7fffffffff000000000000000000000000000000000000000000000000000000008316815260008251610af581600485016020870161075a565b91909101600401939250505056fe60c060405234801561001057600080fd5b5060405161083638038061083683398101604081905261002f91610062565b6001600160a01b039182166080521660a052610095565b80516001600160a01b038116811461005d57600080fd5b919050565b6000806040838503121561007557600080fd5b61007e83610046565b915061008c60208401610046565b90509250929050565b60805160a0516107706100c660003960008181607101526103e101526000818160ec01526102de01526107706000f3fe608060405234801561001057600080fd5b50600436106100675760003560e01c80637eccf63e116100505780637eccf63e146100de578063e3978240146100e7578063f3dd22a51461010e57600080fd5b80631a4e75de1461006c5780633a871cdd146100bd575b600080fd5b6100937f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b6100d06100cb366004610529565b610123565b6040519081526020016100b4565b6100d060005481565b6100937f000000000000000000000000000000000000000000000000000000000000000081565b61012161011c36600461057d565b6103c9565b005b60008054846020013514610179576000546040517f7ba633940000000000000000000000000000000000000000000000000000000081526004810191909152602085013560248201526044015b60405180910390fd5b604080517f4750045d47fce615521b32cee713ff8db50147e98aec5ca94926b52651ca3fa06020808301919091528183018690528251808303840181526060830184528051908201207f190000000000000000000000000000000000000000000000000000000000000060808401527f010000000000000000000000000000000000000000000000000000000000000060818401527f1c7d3b72b37a35523e273aaadd7b4cd66f618bb81429ab053412d51f50ccea61608284015260a2808401919091528351808403909101815260c29092019092528051910120600061026461014087018761062c565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092018290525060208501516040808701518751979850919691955091935086925081106102bd576102bd610698565b016020015160f81c905073ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001660018661030c84601b6106f6565b6040805160008152602081018083529390935260ff90911690820152606081018690526080810185905260a0016020604051602081039080840390855afa15801561035b573d6000803e3d6000fd5b5050506020604051035173ffffffffffffffffffffffffffffffffffffffff16146103995761038d60016000806104f1565b955050505050506103c2565b6000805490806103a88361071b565b91905055506103ba60008060006104f1565b955050505050505b9392505050565b3373ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161461043a576040517f4a0bfec1000000000000000000000000000000000000000000000000000000008152336004820152602401610170565b8242111561047d576040517f300249d700000000000000000000000000000000000000000000000000000000815260048101849052426024820152604401610170565b8473ffffffffffffffffffffffffffffffffffffffff168483836040516104a5929190610753565b60006040518083038185875af1925050503d80600081146104e2576040519150601f19603f3d011682016040523d82523d6000602084013e6104e7565b606091505b5050505050505050565b600060d08265ffffffffffff16901b60a08465ffffffffffff16901b8561051957600061051c565b60015b60ff161717949350505050565b60008060006060848603121561053e57600080fd5b833567ffffffffffffffff81111561055557600080fd5b8401610160818703121561056857600080fd5b95602085013595506040909401359392505050565b60008060008060006080868803121561059557600080fd5b853573ffffffffffffffffffffffffffffffffffffffff811681146105b957600080fd5b94506020860135935060408601359250606086013567ffffffffffffffff808211156105e457600080fd5b818801915088601f8301126105f857600080fd5b81358181111561060757600080fd5b89602082850101111561061957600080fd5b9699959850939650602001949392505050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261066157600080fd5b83018035915067ffffffffffffffff82111561067c57600080fd5b60200191503681900382131561069157600080fd5b9250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600060ff821660ff84168060ff03821115610713576107136106c7565b019392505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361074c5761074c6106c7565b5060010190565b818382376000910190815291905056fea164736f6c634300080f000aa164736f6c634300080f000a",
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

func (_SmartContractAccountHelper *SmartContractAccountHelperCaller) GetAbiEncodedDirectRequestData(opts *bind.CallOpts, recipient common.Address, topupThreshold *big.Int, topupAmount *big.Int) ([]byte, error) {
	var out []interface{}
	err := _SmartContractAccountHelper.contract.Call(opts, &out, "getAbiEncodedDirectRequestData", recipient, topupThreshold, topupAmount)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_SmartContractAccountHelper *SmartContractAccountHelperSession) GetAbiEncodedDirectRequestData(recipient common.Address, topupThreshold *big.Int, topupAmount *big.Int) ([]byte, error) {
	return _SmartContractAccountHelper.Contract.GetAbiEncodedDirectRequestData(&_SmartContractAccountHelper.CallOpts, recipient, topupThreshold, topupAmount)
}

func (_SmartContractAccountHelper *SmartContractAccountHelperCallerSession) GetAbiEncodedDirectRequestData(recipient common.Address, topupThreshold *big.Int, topupAmount *big.Int) ([]byte, error) {
	return _SmartContractAccountHelper.Contract.GetAbiEncodedDirectRequestData(&_SmartContractAccountHelper.CallOpts, recipient, topupThreshold, topupAmount)
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

	GetAbiEncodedDirectRequestData(opts *bind.CallOpts, recipient common.Address, topupThreshold *big.Int, topupAmount *big.Int) ([]byte, error)

	GetFullEndTxEncoding(opts *bind.CallOpts, endContract common.Address, value *big.Int, deadline *big.Int, data []byte) ([]byte, error)

	GetFullHashForSigning(opts *bind.CallOpts, userOpHash [32]byte) ([32]byte, error)

	GetInitCode(opts *bind.CallOpts, factory common.Address, owner common.Address, entryPoint common.Address) ([]byte, error)

	GetSCAInitCodeWithConstructor(opts *bind.CallOpts, owner common.Address, entryPoint common.Address) ([]byte, error)

	Address() common.Address
}
