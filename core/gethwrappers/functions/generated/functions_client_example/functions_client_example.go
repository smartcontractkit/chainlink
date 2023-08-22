// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package functions_client_example

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated"
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

var FunctionsClientExampleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"EmptyArgs\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptySecrets\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptySource\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoInlineSecrets\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyRouterCanFufill\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"UnexpectedRequestID\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"RequestFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"RequestSent\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"MAX_CALLBACK_GAS\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"response\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"handleOracleFulfillment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastError\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastErrorLength\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastRequestId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastResponse\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastResponseLength\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"source\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"encryptedSecretsReferences\",\"type\":\"bytes\"},{\"internalType\":\"string[]\",\"name\":\"args\",\"type\":\"string[]\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"jobId\",\"type\":\"bytes32\"}],\"name\":\"sendRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b50604051620019ed380380620019ed833981016040819052620000349162000198565b600080546001600160a01b0319166001600160a01b038316178155339081906001600160a01b038216620000af5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600180546001600160a01b0319166001600160a01b0384811691909117909155811615620000e257620000e281620000ec565b50505050620001ca565b336001600160a01b03821603620001465760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401620000a6565b600280546001600160a01b0319166001600160a01b03838116918217909255600154604051919216907fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127890600090a350565b600060208284031215620001ab57600080fd5b81516001600160a01b0381168114620001c357600080fd5b9392505050565b61181380620001da6000396000f3fe608060405234801561001057600080fd5b50600436106100c95760003560e01c80636d9809a011610081578063b48cffea1161005b578063b48cffea14610182578063f2fde38b14610192578063fc2a88c3146101a557600080fd5b80636d9809a01461014857806379ba5097146101525780638da5cb5b1461015a57600080fd5b80632c29166b116100b25780632c29166b146100ff5780635fa353e71461012c57806362747e421461013f57600080fd5b80630ca76175146100ce57806329f0de3f146100e3575b600080fd5b6100e16100dc36600461125f565b6101ae565b005b6100ec60055481565b6040519081526020015b60405180910390f35b60065461011790640100000000900463ffffffff1681565b60405163ffffffff90911681526020016100f6565b6100e161013a366004611332565b61020f565b6100ec60045481565b6101176201117081565b6100e1610319565b60015460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100f6565b6006546101179063ffffffff1681565b6100e16101a0366004611416565b61041f565b6100ec60035481565b60005473ffffffffffffffffffffffffffffffffffffffff1633146101ff576040517f5099014100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61020a838383610433565b505050565b610217610501565b6102586040805160e0810190915280600081526020016000815260200160008152602001606081526020016060815260200160608152602001606081525090565b61029a89898080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525085939250506105849050565b85156102e2576102e287878080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525085939250506105959050565b83156102fc576102fc6102f5858761144c565b82906105df565b61030b81846201117085610622565b600355505050505050505050565b60025473ffffffffffffffffffffffffffffffffffffffff16331461039f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000008082163390811790935560028054909116905560405173ffffffffffffffffffffffffffffffffffffffff909116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a350565b610427610501565b61043081610646565b50565b8260035414610471576040517fd068bf5b00000000000000000000000000000000000000000000000000000000815260048101849052602401610396565b61047a8261073c565b6004558151600680547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff9092169190911790556104bc8161073c565b600555516006805463ffffffff909216640100000000027fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff9092169190911790555050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610582576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610396565b565b61059182600080846107c4565b5050565b80516000036105d0576040517fe889636f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60016020830152608090910152565b805160000361061a576040517ffe936cb700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60a090910152565b60008061062e8661085b565b905061063c81868686610bcc565b9695505050505050565b3373ffffffffffffffffffffffffffffffffffffffff8216036106c5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610396565b600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff838116918217909255600154604051919216907fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127890600090a350565b600080600060209050602084511015610753575082515b60005b818110156107bb57610769816008611542565b85828151811061077b5761077b611559565b01602001517fff0000000000000000000000000000000000000000000000000000000000000016901c9290921791806107b381611588565b915050610756565b50909392505050565b80516000036107ff576040517f22ce3edd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b83836002811115610812576108126114e4565b90816002811115610825576108256114e4565b9052506040840182801561083b5761083b6114e4565b9081801561084b5761084b6114e4565b9052506060909301929092525050565b6060610865611116565b805161087390610100610ca5565b506108b3816040518060400160405280600c81526020017f636f64654c6f636174696f6e0000000000000000000000000000000000000000815250610d1f565b6108d281846000015160028111156108cd576108cd6114e4565b610d38565b610911816040518060400160405280600881526020017f6c616e6775616765000000000000000000000000000000000000000000000000815250610d1f565b61092b81846040015160008111156108cd576108cd6114e4565b61096a816040518060400160405280600681526020017f736f757263650000000000000000000000000000000000000000000000000000815250610d1f565b610978818460600151610d1f565b60a08301515115610a1e576109c2816040518060400160405280600481526020017f6172677300000000000000000000000000000000000000000000000000000000815250610d1f565b6109cb81610d71565b60005b8360a0015151811015610a1457610a02828560a0015183815181106109f5576109f5611559565b6020026020010151610d1f565b80610a0c81611588565b9150506109ce565b50610a1e81610d95565b60808301515115610b1e57600083602001516002811115610a4157610a416114e4565b03610a78576040517fa80d31f700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610ab7816040518060400160405280600f81526020017f736563726574734c6f636174696f6e0000000000000000000000000000000000815250610d1f565b610ad181846020015160028111156108cd576108cd6114e4565b610b10816040518060400160405280600781526020017f7365637265747300000000000000000000000000000000000000000000000000815250610d1f565b610b1e818460800151610db3565b60c08301515115610bc457610b68816040518060400160405280600981526020017f6279746573417267730000000000000000000000000000000000000000000000815250610d1f565b610b7181610d71565b60005b8360c0015151811015610bba57610ba8828560c001518381518110610b9b57610b9b611559565b6020026020010151610db3565b80610bb281611588565b915050610b74565b50610bc481610d95565b515192915050565b600080546040517f461d276200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091169063461d276290610c2c9087908990600190899089906004016115c0565b6020604051808303816000875af1158015610c4b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c6f9190611660565b60405190915081907f1131472297a800fee664d1d89cfa8f7676ff07189ecc53f80bbb5f4969099db890600090a2949350505050565b604080518082019091526060815260006020820152610cc5602083611679565b15610ced57610cd5602083611679565b610ce09060206116b4565b610cea90836116c7565b91505b602080840183905260405180855260008152908184010181811015610d1157600080fd5b604052508290505b92915050565b610d2c8260038351610dbc565b815161020a9082610ee3565b8151610d459060c2610f0b565b506105918282604051602001610d5d91815260200190565b604051602081830303815290604052610db3565b610d7c816004610f74565b600181602001818151610d8f91906116c7565b90525050565b610da0816007610f74565b600181602001818151610d8f91906116b4565b610d2c82600283515b60178167ffffffffffffffff1611610de9578251610de39060e0600585901b168317610f0b565b50505050565b60ff8167ffffffffffffffff1611610e2b578251610e12906018611fe0600586901b1617610f0b565b508251610de39067ffffffffffffffff83166001610f8b565b61ffff8167ffffffffffffffff1611610e6e578251610e55906019611fe0600586901b1617610f0b565b508251610de39067ffffffffffffffff83166002610f8b565b63ffffffff8167ffffffffffffffff1611610eb3578251610e9a90601a611fe0600586901b1617610f0b565b508251610de39067ffffffffffffffff83166004610f8b565b8251610eca90601b611fe0600586901b1617610f0b565b508251610de39067ffffffffffffffff83166008610f8b565b604080518082019091526060815260006020820152610f0483838451611010565b9392505050565b6040805180820190915260608152600060208201528251516000610f308260016116c7565b905084602001518210610f5157610f5185610f4c836002611542565b6110ff565b8451602083820101858153508051821115610f6a578181525b5093949350505050565b815161020a90601f611fe0600585901b1617610f0b565b6040805180820190915260608152600060208201528351516000610faf82856116c7565b90508560200151811115610fcc57610fcc86610f4c836002611542565b60006001610fdc866101006117fa565b610fe691906116b4565b90508651828101878319825116178152508051831115611004578281525b50959695505050505050565b604080518082019091526060815260006020820152825182111561103357600080fd5b835151600061104284836116c7565b9050856020015181111561105f5761105f86610f4c836002611542565b855180518382016020019160009180851115611079578482525b505050602086015b602086106110b957805182526110986020836116c7565b91506110a56020826116c7565b90506110b26020876116b4565b9550611081565b5181517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60208890036101000a0190811690199190911617905250849150509392505050565b815161110b8383610ca5565b50610de38382610ee3565b604051806040016040528061113e604051806040016040528060608152602001600081525090565b8152602001600081525090565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156111c1576111c161114b565b604052919050565b600067ffffffffffffffff8311156111e3576111e361114b565b61121460207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8601160161117a565b905082815283838301111561122857600080fd5b828260208301376000602084830101529392505050565b600082601f83011261125057600080fd5b610f04838335602085016111c9565b60008060006060848603121561127457600080fd5b83359250602084013567ffffffffffffffff8082111561129357600080fd5b61129f8783880161123f565b935060408601359150808211156112b557600080fd5b506112c28682870161123f565b9150509250925092565b60008083601f8401126112de57600080fd5b50813567ffffffffffffffff8111156112f657600080fd5b60208301915083602082850101111561130e57600080fd5b9250929050565b803567ffffffffffffffff8116811461132d57600080fd5b919050565b60008060008060008060008060a0898b03121561134e57600080fd5b883567ffffffffffffffff8082111561136657600080fd5b6113728c838d016112cc565b909a50985060208b013591508082111561138b57600080fd5b6113978c838d016112cc565b909850965060408b01359150808211156113b057600080fd5b818b0191508b601f8301126113c457600080fd5b8135818111156113d357600080fd5b8c60208260051b85010111156113e857600080fd5b60208301965080955050505061140060608a01611315565b9150608089013590509295985092959890939650565b60006020828403121561142857600080fd5b813573ffffffffffffffffffffffffffffffffffffffff81168114610f0457600080fd5b600067ffffffffffffffff808411156114675761146761114b565b8360051b602061147881830161117a565b86815291850191818101903684111561149057600080fd5b865b848110156114d8578035868111156114aa5760008081fd5b880136601f8201126114bc5760008081fd5b6114ca3682358784016111c9565b845250918301918301611492565b50979650505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8082028115828204841417610d1957610d19611513565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036115b9576115b9611513565b5060010190565b67ffffffffffffffff861681526000602060a08184015286518060a085015260005b818110156115fe5788810183015185820160c0015282016115e2565b50600060c0828601015260c07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010192505050611647604083018661ffff169052565b63ffffffff939093166060820152608001529392505050565b60006020828403121561167257600080fd5b5051919050565b6000826116af577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500690565b81810381811115610d1957610d19611513565b80820180821115610d1957610d19611513565b600181815b8085111561173357817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0482111561171957611719611513565b8085161561172657918102915b93841c93908002906116df565b509250929050565b60008261174a57506001610d19565b8161175757506000610d19565b816001811461176d576002811461177757611793565b6001915050610d19565b60ff84111561178857611788611513565b50506001821b610d19565b5060208310610133831016604e8410600b84101617156117b6575081810a610d19565b6117c083836116da565b807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048211156117f2576117f2611513565b029392505050565b6000610f04838361173b56fea164736f6c6343000813000a",
}

var FunctionsClientExampleABI = FunctionsClientExampleMetaData.ABI

var FunctionsClientExampleBin = FunctionsClientExampleMetaData.Bin

func DeployFunctionsClientExample(auth *bind.TransactOpts, backend bind.ContractBackend, router common.Address) (common.Address, *types.Transaction, *FunctionsClientExample, error) {
	parsed, err := FunctionsClientExampleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(FunctionsClientExampleBin), backend, router)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &FunctionsClientExample{FunctionsClientExampleCaller: FunctionsClientExampleCaller{contract: contract}, FunctionsClientExampleTransactor: FunctionsClientExampleTransactor{contract: contract}, FunctionsClientExampleFilterer: FunctionsClientExampleFilterer{contract: contract}}, nil
}

type FunctionsClientExample struct {
	address common.Address
	abi     abi.ABI
	FunctionsClientExampleCaller
	FunctionsClientExampleTransactor
	FunctionsClientExampleFilterer
}

type FunctionsClientExampleCaller struct {
	contract *bind.BoundContract
}

type FunctionsClientExampleTransactor struct {
	contract *bind.BoundContract
}

type FunctionsClientExampleFilterer struct {
	contract *bind.BoundContract
}

type FunctionsClientExampleSession struct {
	Contract     *FunctionsClientExample
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type FunctionsClientExampleCallerSession struct {
	Contract *FunctionsClientExampleCaller
	CallOpts bind.CallOpts
}

type FunctionsClientExampleTransactorSession struct {
	Contract     *FunctionsClientExampleTransactor
	TransactOpts bind.TransactOpts
}

type FunctionsClientExampleRaw struct {
	Contract *FunctionsClientExample
}

type FunctionsClientExampleCallerRaw struct {
	Contract *FunctionsClientExampleCaller
}

type FunctionsClientExampleTransactorRaw struct {
	Contract *FunctionsClientExampleTransactor
}

func NewFunctionsClientExample(address common.Address, backend bind.ContractBackend) (*FunctionsClientExample, error) {
	abi, err := abi.JSON(strings.NewReader(FunctionsClientExampleABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindFunctionsClientExample(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FunctionsClientExample{address: address, abi: abi, FunctionsClientExampleCaller: FunctionsClientExampleCaller{contract: contract}, FunctionsClientExampleTransactor: FunctionsClientExampleTransactor{contract: contract}, FunctionsClientExampleFilterer: FunctionsClientExampleFilterer{contract: contract}}, nil
}

func NewFunctionsClientExampleCaller(address common.Address, caller bind.ContractCaller) (*FunctionsClientExampleCaller, error) {
	contract, err := bindFunctionsClientExample(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FunctionsClientExampleCaller{contract: contract}, nil
}

func NewFunctionsClientExampleTransactor(address common.Address, transactor bind.ContractTransactor) (*FunctionsClientExampleTransactor, error) {
	contract, err := bindFunctionsClientExample(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FunctionsClientExampleTransactor{contract: contract}, nil
}

func NewFunctionsClientExampleFilterer(address common.Address, filterer bind.ContractFilterer) (*FunctionsClientExampleFilterer, error) {
	contract, err := bindFunctionsClientExample(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FunctionsClientExampleFilterer{contract: contract}, nil
}

func bindFunctionsClientExample(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FunctionsClientExampleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_FunctionsClientExample *FunctionsClientExampleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FunctionsClientExample.Contract.FunctionsClientExampleCaller.contract.Call(opts, result, method, params...)
}

func (_FunctionsClientExample *FunctionsClientExampleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsClientExample.Contract.FunctionsClientExampleTransactor.contract.Transfer(opts)
}

func (_FunctionsClientExample *FunctionsClientExampleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FunctionsClientExample.Contract.FunctionsClientExampleTransactor.contract.Transact(opts, method, params...)
}

func (_FunctionsClientExample *FunctionsClientExampleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FunctionsClientExample.Contract.contract.Call(opts, result, method, params...)
}

func (_FunctionsClientExample *FunctionsClientExampleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsClientExample.Contract.contract.Transfer(opts)
}

func (_FunctionsClientExample *FunctionsClientExampleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FunctionsClientExample.Contract.contract.Transact(opts, method, params...)
}

func (_FunctionsClientExample *FunctionsClientExampleCaller) MAXCALLBACKGAS(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _FunctionsClientExample.contract.Call(opts, &out, "MAX_CALLBACK_GAS")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_FunctionsClientExample *FunctionsClientExampleSession) MAXCALLBACKGAS() (uint32, error) {
	return _FunctionsClientExample.Contract.MAXCALLBACKGAS(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleCallerSession) MAXCALLBACKGAS() (uint32, error) {
	return _FunctionsClientExample.Contract.MAXCALLBACKGAS(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleCaller) LastError(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _FunctionsClientExample.contract.Call(opts, &out, "lastError")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_FunctionsClientExample *FunctionsClientExampleSession) LastError() ([32]byte, error) {
	return _FunctionsClientExample.Contract.LastError(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleCallerSession) LastError() ([32]byte, error) {
	return _FunctionsClientExample.Contract.LastError(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleCaller) LastErrorLength(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _FunctionsClientExample.contract.Call(opts, &out, "lastErrorLength")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_FunctionsClientExample *FunctionsClientExampleSession) LastErrorLength() (uint32, error) {
	return _FunctionsClientExample.Contract.LastErrorLength(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleCallerSession) LastErrorLength() (uint32, error) {
	return _FunctionsClientExample.Contract.LastErrorLength(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleCaller) LastRequestId(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _FunctionsClientExample.contract.Call(opts, &out, "lastRequestId")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_FunctionsClientExample *FunctionsClientExampleSession) LastRequestId() ([32]byte, error) {
	return _FunctionsClientExample.Contract.LastRequestId(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleCallerSession) LastRequestId() ([32]byte, error) {
	return _FunctionsClientExample.Contract.LastRequestId(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleCaller) LastResponse(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _FunctionsClientExample.contract.Call(opts, &out, "lastResponse")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_FunctionsClientExample *FunctionsClientExampleSession) LastResponse() ([32]byte, error) {
	return _FunctionsClientExample.Contract.LastResponse(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleCallerSession) LastResponse() ([32]byte, error) {
	return _FunctionsClientExample.Contract.LastResponse(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleCaller) LastResponseLength(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _FunctionsClientExample.contract.Call(opts, &out, "lastResponseLength")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_FunctionsClientExample *FunctionsClientExampleSession) LastResponseLength() (uint32, error) {
	return _FunctionsClientExample.Contract.LastResponseLength(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleCallerSession) LastResponseLength() (uint32, error) {
	return _FunctionsClientExample.Contract.LastResponseLength(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FunctionsClientExample.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_FunctionsClientExample *FunctionsClientExampleSession) Owner() (common.Address, error) {
	return _FunctionsClientExample.Contract.Owner(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleCallerSession) Owner() (common.Address, error) {
	return _FunctionsClientExample.Contract.Owner(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsClientExample.contract.Transact(opts, "acceptOwnership")
}

func (_FunctionsClientExample *FunctionsClientExampleSession) AcceptOwnership() (*types.Transaction, error) {
	return _FunctionsClientExample.Contract.AcceptOwnership(&_FunctionsClientExample.TransactOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _FunctionsClientExample.Contract.AcceptOwnership(&_FunctionsClientExample.TransactOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleTransactor) HandleOracleFulfillment(opts *bind.TransactOpts, requestId [32]byte, response []byte, err []byte) (*types.Transaction, error) {
	return _FunctionsClientExample.contract.Transact(opts, "handleOracleFulfillment", requestId, response, err)
}

func (_FunctionsClientExample *FunctionsClientExampleSession) HandleOracleFulfillment(requestId [32]byte, response []byte, err []byte) (*types.Transaction, error) {
	return _FunctionsClientExample.Contract.HandleOracleFulfillment(&_FunctionsClientExample.TransactOpts, requestId, response, err)
}

func (_FunctionsClientExample *FunctionsClientExampleTransactorSession) HandleOracleFulfillment(requestId [32]byte, response []byte, err []byte) (*types.Transaction, error) {
	return _FunctionsClientExample.Contract.HandleOracleFulfillment(&_FunctionsClientExample.TransactOpts, requestId, response, err)
}

func (_FunctionsClientExample *FunctionsClientExampleTransactor) SendRequest(opts *bind.TransactOpts, source string, encryptedSecretsReferences []byte, args []string, subscriptionId uint64, jobId [32]byte) (*types.Transaction, error) {
	return _FunctionsClientExample.contract.Transact(opts, "sendRequest", source, encryptedSecretsReferences, args, subscriptionId, jobId)
}

func (_FunctionsClientExample *FunctionsClientExampleSession) SendRequest(source string, encryptedSecretsReferences []byte, args []string, subscriptionId uint64, jobId [32]byte) (*types.Transaction, error) {
	return _FunctionsClientExample.Contract.SendRequest(&_FunctionsClientExample.TransactOpts, source, encryptedSecretsReferences, args, subscriptionId, jobId)
}

func (_FunctionsClientExample *FunctionsClientExampleTransactorSession) SendRequest(source string, encryptedSecretsReferences []byte, args []string, subscriptionId uint64, jobId [32]byte) (*types.Transaction, error) {
	return _FunctionsClientExample.Contract.SendRequest(&_FunctionsClientExample.TransactOpts, source, encryptedSecretsReferences, args, subscriptionId, jobId)
}

func (_FunctionsClientExample *FunctionsClientExampleTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _FunctionsClientExample.contract.Transact(opts, "transferOwnership", to)
}

func (_FunctionsClientExample *FunctionsClientExampleSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _FunctionsClientExample.Contract.TransferOwnership(&_FunctionsClientExample.TransactOpts, to)
}

func (_FunctionsClientExample *FunctionsClientExampleTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _FunctionsClientExample.Contract.TransferOwnership(&_FunctionsClientExample.TransactOpts, to)
}

type FunctionsClientExampleOwnershipTransferRequestedIterator struct {
	Event *FunctionsClientExampleOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsClientExampleOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsClientExampleOwnershipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(FunctionsClientExampleOwnershipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *FunctionsClientExampleOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *FunctionsClientExampleOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsClientExampleOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_FunctionsClientExample *FunctionsClientExampleFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsClientExampleOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsClientExample.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsClientExampleOwnershipTransferRequestedIterator{contract: _FunctionsClientExample.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_FunctionsClientExample *FunctionsClientExampleFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *FunctionsClientExampleOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsClientExample.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsClientExampleOwnershipTransferRequested)
				if err := _FunctionsClientExample.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_FunctionsClientExample *FunctionsClientExampleFilterer) ParseOwnershipTransferRequested(log types.Log) (*FunctionsClientExampleOwnershipTransferRequested, error) {
	event := new(FunctionsClientExampleOwnershipTransferRequested)
	if err := _FunctionsClientExample.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsClientExampleOwnershipTransferredIterator struct {
	Event *FunctionsClientExampleOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsClientExampleOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsClientExampleOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(FunctionsClientExampleOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *FunctionsClientExampleOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *FunctionsClientExampleOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsClientExampleOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_FunctionsClientExample *FunctionsClientExampleFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsClientExampleOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsClientExample.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsClientExampleOwnershipTransferredIterator{contract: _FunctionsClientExample.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_FunctionsClientExample *FunctionsClientExampleFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FunctionsClientExampleOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsClientExample.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsClientExampleOwnershipTransferred)
				if err := _FunctionsClientExample.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_FunctionsClientExample *FunctionsClientExampleFilterer) ParseOwnershipTransferred(log types.Log) (*FunctionsClientExampleOwnershipTransferred, error) {
	event := new(FunctionsClientExampleOwnershipTransferred)
	if err := _FunctionsClientExample.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsClientExampleRequestFulfilledIterator struct {
	Event *FunctionsClientExampleRequestFulfilled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsClientExampleRequestFulfilledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsClientExampleRequestFulfilled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(FunctionsClientExampleRequestFulfilled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *FunctionsClientExampleRequestFulfilledIterator) Error() error {
	return it.fail
}

func (it *FunctionsClientExampleRequestFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsClientExampleRequestFulfilled struct {
	Id  [32]byte
	Raw types.Log
}

func (_FunctionsClientExample *FunctionsClientExampleFilterer) FilterRequestFulfilled(opts *bind.FilterOpts, id [][32]byte) (*FunctionsClientExampleRequestFulfilledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _FunctionsClientExample.contract.FilterLogs(opts, "RequestFulfilled", idRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsClientExampleRequestFulfilledIterator{contract: _FunctionsClientExample.contract, event: "RequestFulfilled", logs: logs, sub: sub}, nil
}

func (_FunctionsClientExample *FunctionsClientExampleFilterer) WatchRequestFulfilled(opts *bind.WatchOpts, sink chan<- *FunctionsClientExampleRequestFulfilled, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _FunctionsClientExample.contract.WatchLogs(opts, "RequestFulfilled", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsClientExampleRequestFulfilled)
				if err := _FunctionsClientExample.contract.UnpackLog(event, "RequestFulfilled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_FunctionsClientExample *FunctionsClientExampleFilterer) ParseRequestFulfilled(log types.Log) (*FunctionsClientExampleRequestFulfilled, error) {
	event := new(FunctionsClientExampleRequestFulfilled)
	if err := _FunctionsClientExample.contract.UnpackLog(event, "RequestFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsClientExampleRequestSentIterator struct {
	Event *FunctionsClientExampleRequestSent

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsClientExampleRequestSentIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsClientExampleRequestSent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(FunctionsClientExampleRequestSent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *FunctionsClientExampleRequestSentIterator) Error() error {
	return it.fail
}

func (it *FunctionsClientExampleRequestSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsClientExampleRequestSent struct {
	Id  [32]byte
	Raw types.Log
}

func (_FunctionsClientExample *FunctionsClientExampleFilterer) FilterRequestSent(opts *bind.FilterOpts, id [][32]byte) (*FunctionsClientExampleRequestSentIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _FunctionsClientExample.contract.FilterLogs(opts, "RequestSent", idRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsClientExampleRequestSentIterator{contract: _FunctionsClientExample.contract, event: "RequestSent", logs: logs, sub: sub}, nil
}

func (_FunctionsClientExample *FunctionsClientExampleFilterer) WatchRequestSent(opts *bind.WatchOpts, sink chan<- *FunctionsClientExampleRequestSent, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _FunctionsClientExample.contract.WatchLogs(opts, "RequestSent", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsClientExampleRequestSent)
				if err := _FunctionsClientExample.contract.UnpackLog(event, "RequestSent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_FunctionsClientExample *FunctionsClientExampleFilterer) ParseRequestSent(log types.Log) (*FunctionsClientExampleRequestSent, error) {
	event := new(FunctionsClientExampleRequestSent)
	if err := _FunctionsClientExample.contract.UnpackLog(event, "RequestSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_FunctionsClientExample *FunctionsClientExample) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _FunctionsClientExample.abi.Events["OwnershipTransferRequested"].ID:
		return _FunctionsClientExample.ParseOwnershipTransferRequested(log)
	case _FunctionsClientExample.abi.Events["OwnershipTransferred"].ID:
		return _FunctionsClientExample.ParseOwnershipTransferred(log)
	case _FunctionsClientExample.abi.Events["RequestFulfilled"].ID:
		return _FunctionsClientExample.ParseRequestFulfilled(log)
	case _FunctionsClientExample.abi.Events["RequestSent"].ID:
		return _FunctionsClientExample.ParseRequestSent(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (FunctionsClientExampleOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (FunctionsClientExampleOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (FunctionsClientExampleRequestFulfilled) Topic() common.Hash {
	return common.HexToHash("0x85e1543bf2f84fe80c6badbce3648c8539ad1df4d2b3d822938ca0538be727e6")
}

func (FunctionsClientExampleRequestSent) Topic() common.Hash {
	return common.HexToHash("0x1131472297a800fee664d1d89cfa8f7676ff07189ecc53f80bbb5f4969099db8")
}

func (_FunctionsClientExample *FunctionsClientExample) Address() common.Address {
	return _FunctionsClientExample.address
}

type FunctionsClientExampleInterface interface {
	MAXCALLBACKGAS(opts *bind.CallOpts) (uint32, error)

	LastError(opts *bind.CallOpts) ([32]byte, error)

	LastErrorLength(opts *bind.CallOpts) (uint32, error)

	LastRequestId(opts *bind.CallOpts) ([32]byte, error)

	LastResponse(opts *bind.CallOpts) ([32]byte, error)

	LastResponseLength(opts *bind.CallOpts) (uint32, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	HandleOracleFulfillment(opts *bind.TransactOpts, requestId [32]byte, response []byte, err []byte) (*types.Transaction, error)

	SendRequest(opts *bind.TransactOpts, source string, encryptedSecretsReferences []byte, args []string, subscriptionId uint64, jobId [32]byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsClientExampleOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *FunctionsClientExampleOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*FunctionsClientExampleOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsClientExampleOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FunctionsClientExampleOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*FunctionsClientExampleOwnershipTransferred, error)

	FilterRequestFulfilled(opts *bind.FilterOpts, id [][32]byte) (*FunctionsClientExampleRequestFulfilledIterator, error)

	WatchRequestFulfilled(opts *bind.WatchOpts, sink chan<- *FunctionsClientExampleRequestFulfilled, id [][32]byte) (event.Subscription, error)

	ParseRequestFulfilled(log types.Log) (*FunctionsClientExampleRequestFulfilled, error)

	FilterRequestSent(opts *bind.FilterOpts, id [][32]byte) (*FunctionsClientExampleRequestSentIterator, error)

	WatchRequestSent(opts *bind.WatchOpts, sink chan<- *FunctionsClientExampleRequestSent, id [][32]byte) (event.Subscription, error)

	ParseRequestSent(log types.Log) (*FunctionsClientExampleRequestSent, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
