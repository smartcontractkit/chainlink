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
	_ = abi.ConvertType
)

var SmartContractAccountHelperMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"entryPoint\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"factory\",\"type\":\"address\"}],\"name\":\"calculateSmartContractAccountAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"topupThreshold\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"topupAmount\",\"type\":\"uint256\"}],\"name\":\"getAbiEncodedDirectRequestData\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"endContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"getFullEndTxEncoding\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"encoding\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"userOpHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"scaAddress\",\"type\":\"address\"}],\"name\":\"getFullHashForSigning\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"factory\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"entryPoint\",\"type\":\"address\"}],\"name\":\"getInitCode\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"entryPoint\",\"type\":\"address\"}],\"name\":\"getSCAInitCodeWithConstructor\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x61162261003a600b82828239805160001a60731461002d57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe730000000000000000000000000000000000000000301460806040526004361061007c5760003560e01c8063e0237bef1161005a578063e0237bef14610134578063e464b3631461016c578063fc59bac31461017f57600080fd5b80632c86cb35146100815780634b770f561461010057806382311e3314610113575b600080fd5b6100ea61008f36600461076b565b604080516060808201835273ffffffffffffffffffffffffffffffffffffffff959095168082526020808301958652918301938452825191820152925183820152905182840152805180830390930183526080909101905290565b6040516100f79190610818565b60405180910390f35b6100ea61010e36600461082b565b610192565b61012661012136600461086e565b610336565b6040519081526020016100f7565b61014761014236600461082b565b610454565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100f7565b6100ea61017a36600461089a565b6105df565b6100ea61018d3660046108f3565b61069d565b6040516060907fffffffffffffffffffffffffffffffffffffffff00000000000000000000000084831b16906000906101cd60208201610735565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe082820381018352601f90910116604081815273ffffffffffffffffffffffffffffffffffffffff8881166020840152871690820152606001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529082905261026292916020016109e6565b60405160208183030381529060405290508560601b630af4926f60e01b8383604051602401610292929190610a15565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152918152602080830180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909516949094179093525161031c939201610a36565b604051602081830303815290604052925050509392505050565b600061044d8383604080517f4750045d47fce615521b32cee713ff8db50147e98aec5ca94926b52651ca3fa060208083019190915281830194909452815180820383018152606080830184528151918601919091207f190000000000000000000000000000000000000000000000000000000000000060808401527f010000000000000000000000000000000000000000000000000000000000000060818401527f1c7d3b72b37a35523e273aaadd7b4cd66f618bb81429ab053412d51f50ccea6160828401524660a284015293901b7fffffffffffffffffffffffffffffffffffffffff0000000000000000000000001660c282015260d6808201939093528151808203909301835260f6019052805191012090565b9392505050565b6040516000907fffffffffffffffffffffffffffffffffffffffff000000000000000000000000606086901b1690829061049060208201610735565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe082820381018352601f90910116604081815273ffffffffffffffffffffffffffffffffffffffff8981166020840152881690820152606001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529082905261052592916020016109e6565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001815282825280516020918201207fff000000000000000000000000000000000000000000000000000000000000008285015260609790971b7fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166021840152603583019490945260558083019690965280518083039096018652607590910190525082519201919091209392505050565b6060604051806020016105f190610735565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe082820381018352601f90910116604081815273ffffffffffffffffffffffffffffffffffffffff8681166020840152851690820152606001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529082905261068692916020016109e6565b604051602081830303815290604052905092915050565b60607f89553be40000000000000000000000000000000000000000000000000000000085856106cc8642610a7e565b856040516020016106e09493929190610abd565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529082905261071c9291602001610b02565b6040516020818303038152906040529050949350505050565b610acb80610b4b83390190565b803573ffffffffffffffffffffffffffffffffffffffff8116811461076657600080fd5b919050565b60008060006060848603121561078057600080fd5b61078984610742565b95602085013595506040909401359392505050565b60005b838110156107b95781810151838201526020016107a1565b838111156107c8576000848401525b50505050565b600081518084526107e681602086016020860161079e565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b60208152600061044d60208301846107ce565b60008060006060848603121561084057600080fd5b61084984610742565b925061085760208501610742565b915061086560408501610742565b90509250925092565b6000806040838503121561088157600080fd5b8235915061089160208401610742565b90509250929050565b600080604083850312156108ad57600080fd5b6108b683610742565b915061089160208401610742565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6000806000806080858703121561090957600080fd5b61091285610742565b93506020850135925060408501359150606085013567ffffffffffffffff8082111561093d57600080fd5b818701915087601f83011261095157600080fd5b813581811115610963576109636108c4565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f011681019083821181831017156109a9576109a96108c4565b816040528281528a60208487010111156109c257600080fd5b82602086016020830137600060208483010152809550505050505092959194509250565b600083516109f881846020880161079e565b835190830190610a0c81836020880161079e565b01949350505050565b828152604060208201526000610a2e60408301846107ce565b949350505050565b7fffffffffffffffffffffffffffffffffffffffff0000000000000000000000008316815260008251610a7081601485016020870161079e565b919091016014019392505050565b60008219821115610ab8577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b500190565b73ffffffffffffffffffffffffffffffffffffffff85168152836020820152826040820152608060608201526000610af860808301846107ce565b9695505050505050565b7fffffffff000000000000000000000000000000000000000000000000000000008316815260008251610b3c81600485016020870161079e565b91909101600401939250505056fe60c060405234801561001057600080fd5b50604051610acb380380610acb83398101604081905261002f91610062565b6001600160a01b039182166080521660a052610095565b80516001600160a01b038116811461005d57600080fd5b919050565b6000806040838503121561007557600080fd5b61007e83610046565b915061008c60208401610046565b90509250929050565b60805160a051610a046100c760003960008181607101526103c301526000818161010101526102ee0152610a046000f3fe608060405234801561001057600080fd5b50600436106100675760003560e01c80637eccf63e116100505780637eccf63e146100de57806389553be4146100e7578063dba6335f146100fc57600080fd5b8063140fcfb11461006c5780633a871cdd146100bd575b600080fd5b6100937f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b6100d06100cb36600461063a565b610123565b6040519081526020016100b4565b6100d060005481565b6100fa6100f53660046106ce565b6103ab565b005b6100937f000000000000000000000000000000000000000000000000000000000000000081565b60008054846020013514610179576000546040517f7ba633940000000000000000000000000000000000000000000000000000000081526004810191909152602085013560248201526044015b60405180910390fd5b60006102908430604080517f4750045d47fce615521b32cee713ff8db50147e98aec5ca94926b52651ca3fa060208083019190915281830194909452815180820383018152606080830184528151918601919091207f190000000000000000000000000000000000000000000000000000000000000060808401527f010000000000000000000000000000000000000000000000000000000000000060818401527f1c7d3b72b37a35523e273aaadd7b4cd66f618bb81429ab053412d51f50ccea6160828401524660a284015293901b7fffffffffffffffffffffffffffffffffffffffff0000000000000000000000001660c282015260d6808201939093528151808203909301835260f6019052805191012090565b905060006102a261014087018761076b565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509293505073ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016915061031c90508284610544565b73ffffffffffffffffffffffffffffffffffffffff161461034d576103446001600080610602565b925050506103a4565b60008054908061035c83610806565b9091555060009050610371606088018861076b565b61037f91600490829061083e565b81019061038c9190610897565b509250505061039e6000826000610602565b93505050505b9392505050565b3373ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161461041c576040517f4a0bfec1000000000000000000000000000000000000000000000000000000008152336004820152602401610170565b65ffffffffffff83161580159061043a57508265ffffffffffff1642115b15610481576040517f300249d700000000000000000000000000000000000000000000000000000000815265ffffffffffff84166004820152426024820152604401610170565b6000808673ffffffffffffffffffffffffffffffffffffffff168685856040516104ac929190610993565b60006040518083038185875af1925050503d80600081146104e9576040519150601f19603f3d011682016040523d82523d6000602084013e6104ee565b606091505b50915091508161053b578051600003610533576040517f20e9b5d200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b805181602001fd5b50505050505050565b602082015160408084015184516000939284918791908110610568576105686109a3565b016020015160f81c905060018561058083601b6109d2565b6040805160008152602081018083529390935260ff90911690820152606081018590526080810184905260a0016020604051602081039080840390855afa1580156105cf573d6000803e3d6000fd5b50506040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00151979650505050505050565b600060d08265ffffffffffff16901b60a08465ffffffffffff16901b8561062a57600061062d565b60015b60ff161717949350505050565b60008060006060848603121561064f57600080fd5b833567ffffffffffffffff81111561066657600080fd5b8401610160818703121561067957600080fd5b95602085013595506040909401359392505050565b73ffffffffffffffffffffffffffffffffffffffff811681146106b057600080fd5b50565b803565ffffffffffff811681146106c957600080fd5b919050565b6000806000806000608086880312156106e657600080fd5b85356106f18161068e565b945060208601359350610706604087016106b3565b9250606086013567ffffffffffffffff8082111561072357600080fd5b818801915088601f83011261073757600080fd5b81358181111561074657600080fd5b89602082850101111561075857600080fd5b9699959850939650602001949392505050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18436030181126107a057600080fd5b83018035915067ffffffffffffffff8211156107bb57600080fd5b6020019150368190038213156107d057600080fd5b9250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203610837576108376107d7565b5060010190565b6000808585111561084e57600080fd5b8386111561085b57600080fd5b5050820193919092039150565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600080600080608085870312156108ad57600080fd5b84356108b88161068e565b9350602085013592506108cd604086016106b3565b9150606085013567ffffffffffffffff808211156108ea57600080fd5b818701915087601f8301126108fe57600080fd5b81358181111561091057610910610868565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f0116810190838211818310171561095657610956610868565b816040528281528a602084870101111561096f57600080fd5b82602086016020830137600060208483010152809550505050505092959194509250565b8183823760009101908152919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600060ff821660ff84168060ff038211156109ef576109ef6107d7565b01939250505056fea164736f6c634300080f000aa164736f6c634300080f000a",
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
	parsed, err := SmartContractAccountHelperMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
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

func (_SmartContractAccountHelper *SmartContractAccountHelperCaller) GetFullHashForSigning(opts *bind.CallOpts, userOpHash [32]byte, scaAddress common.Address) ([32]byte, error) {
	var out []interface{}
	err := _SmartContractAccountHelper.contract.Call(opts, &out, "getFullHashForSigning", userOpHash, scaAddress)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_SmartContractAccountHelper *SmartContractAccountHelperSession) GetFullHashForSigning(userOpHash [32]byte, scaAddress common.Address) ([32]byte, error) {
	return _SmartContractAccountHelper.Contract.GetFullHashForSigning(&_SmartContractAccountHelper.CallOpts, userOpHash, scaAddress)
}

func (_SmartContractAccountHelper *SmartContractAccountHelperCallerSession) GetFullHashForSigning(userOpHash [32]byte, scaAddress common.Address) ([32]byte, error) {
	return _SmartContractAccountHelper.Contract.GetFullHashForSigning(&_SmartContractAccountHelper.CallOpts, userOpHash, scaAddress)
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

	GetFullHashForSigning(opts *bind.CallOpts, userOpHash [32]byte, scaAddress common.Address) ([32]byte, error)

	GetInitCode(opts *bind.CallOpts, factory common.Address, owner common.Address, entryPoint common.Address) ([]byte, error)

	GetSCAInitCodeWithConstructor(opts *bind.CallOpts, owner common.Address, entryPoint common.Address) ([]byte, error)

	Address() common.Address
}
