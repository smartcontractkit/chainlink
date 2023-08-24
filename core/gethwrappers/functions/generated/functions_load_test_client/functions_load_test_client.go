// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package functions_load_test_client

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

var FunctionsLoadTestClientMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"EmptyArgs\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptySecrets\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptySource\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoInlineSecrets\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyRouterCanFulfill\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"RequestFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"RequestSent\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"MAX_CALLBACK_GAS\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"response\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"handleOracleFulfillment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_lastError\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_lastErrorLength\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_lastRequestId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_lastResponse\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_lastResponseLength\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"source\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"encryptedSecretsReferences\",\"type\":\"bytes\"},{\"internalType\":\"string[]\",\"name\":\"args\",\"type\":\"string[]\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"jobId\",\"type\":\"bytes32\"}],\"name\":\"sendRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b5060405162001a3d38038062001a3d833981016040819052620000349162000180565b6001600160a01b0381166080523380600081620000985760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000cb57620000cb81620000d5565b50505050620001b2565b336001600160a01b038216036200012f5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016200008f565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000602082840312156200019357600080fd5b81516001600160a01b0381168114620001ab57600080fd5b9392505050565b608051611868620001d5600039600081816101c601526109f101526118686000f3fe608060405234801561001057600080fd5b50600436106100c95760003560e01c80636d9809a011610081578063b1e217491161005b578063b1e2174914610182578063f2fde38b1461018b578063f7b4c06f1461019e57600080fd5b80636d9809a01461014857806379ba5097146101525780638da5cb5b1461015a57600080fd5b806342748b2a116100b257806342748b2a146100ff5780634b0795a81461012c5780635fa353e71461013557600080fd5b80630ca76175146100ce5780633944ea3a146100e3575b600080fd5b6100e16100dc3660046112b4565b6101ae565b005b6100ec60035481565b6040519081526020015b60405180910390f35b60055461011790640100000000900463ffffffff1681565b60405163ffffffff90911681526020016100f6565b6100ec60045481565b6100e1610143366004611387565b61022d565b6101176201117081565b6100e1610347565b60005460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100f6565b6100ec60025481565b6100e161019936600461146b565b610449565b6005546101179063ffffffff1681565b3373ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161461021d576040517fc6829f8300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61022883838361045d565b505050565b6102356104f2565b61027e6040805161010081019091528060008152602001600081526020016000815260200160608152602001606081526020016060815260200160608152602001606081525090565b6102c089898080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525085939250506105759050565b85156103085761030887878080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525085939250506105869050565b83156103225761032261031b85876114a1565b82906105d0565b61033961032e82610613565b8462011170856109ec565b600255505050505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146103cd576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6104516104f2565b61045a81610acb565b50565b600283905561046b82610bc0565b6003558151600580547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff9092169190911790556104ad81610bc0565b600455516005805463ffffffff909216640100000000027fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff9092169190911790555050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610573576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016103c4565b565b6105828260008084610c42565b5050565b80516000036105c1576040517fe889636f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60016020830152608090910152565b805160000361060b576040517ffe936cb700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60c090910152565b60606000610622610100610cd9565b905061066c6040518060400160405280600c81526020017f636f64654c6f636174696f6e000000000000000000000000000000000000000081525082610cfa90919063ffffffff16565b825161068a90600281111561068357610683611539565b8290610d13565b60408051808201909152600881527f6c616e677561676500000000000000000000000000000000000000000000000060208201526106c9908290610cfa565b60408301516106e090801561068357610683611539565b60408051808201909152600681527f736f757263650000000000000000000000000000000000000000000000000000602082015261071f908290610cfa565b606083015161072f908290610cfa565b60a083015151156107895760408051808201909152601081527f726571756573745369676e6174757265000000000000000000000000000000006020820152610779908290610cfa565b60a0830151610789908290610d48565b60c083015151156108365760408051808201909152600481527f617267730000000000000000000000000000000000000000000000000000000060208201526107d3908290610cfa565b6107dc81610d55565b60005b8360c001515181101561082c5761081c8460c00151828151811061080557610805611568565b602002602001015183610cfa90919063ffffffff16565b610825816115c6565b90506107df565b5061083681610d79565b608083015151156109375760008360200151600281111561085957610859611539565b03610890576040517fa80d31f700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60408051808201909152600f81527f736563726574734c6f636174696f6e000000000000000000000000000000000060208201526108cf908290610cfa565b6108e88360200151600281111561068357610683611539565b60408051808201909152600781527f73656372657473000000000000000000000000000000000000000000000000006020820152610927908290610cfa565b6080830151610937908290610d48565b60e083015151156109e45760408051808201909152600981527f62797465734172677300000000000000000000000000000000000000000000006020820152610981908290610cfa565b61098a81610d55565b60005b8360e00151518110156109da576109ca8460e0015182815181106109b3576109b3611568565b602002602001015183610d4890919063ffffffff16565b6109d3816115c6565b905061098d565b506109e481610d79565b515192915050565b6000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663461d27628688600188886040518663ffffffff1660e01b8152600401610a519594939291906115fe565b6020604051808303816000875af1158015610a70573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a94919061169e565b60405190915081907f1131472297a800fee664d1d89cfa8f7676ff07189ecc53f80bbb5f4969099db890600090a295945050505050565b3373ffffffffffffffffffffffffffffffffffffffff821603610b4a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016103c4565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60008060209050602083511015610bd5575081515b60005b81811015610c3b57610beb8160086116b7565b848281518110610bfd57610bfd611568565b01602001517fff0000000000000000000000000000000000000000000000000000000000000016901c9290921791610c34816115c6565b9050610bd8565b5050919050565b8051600003610c7d576040517f22ce3edd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b83836002811115610c9057610c90611539565b90816002811115610ca357610ca3611539565b90525060408401828015610cb957610cb9611539565b90818015610cc957610cc9611539565b9052506060909301929092525050565b610ce161116b565b8051610ced9083610d97565b5060006020820152919050565b610d078260038351610e11565b81516102289082610f38565b8151610d209060c2610f60565b506105828282604051602001610d3891815260200190565b6040516020818303038152906040525b610d078260028351610e11565b610d60816004610fc9565b600181602001818151610d7391906116ce565b90525050565b610d84816007610fc9565b600181602001818151610d7391906116e1565b604080518082019091526060815260006020820152610db76020836116f4565b15610ddf57610dc76020836116f4565b610dd29060206116e1565b610ddc90836116ce565b91505b602080840183905260405180855260008152908184010181811015610e0357600080fd5b604052508290505b92915050565b60178167ffffffffffffffff1611610e3e578251610e389060e0600585901b168317610f60565b50505050565b60ff8167ffffffffffffffff1611610e80578251610e67906018611fe0600586901b1617610f60565b508251610e389067ffffffffffffffff83166001610fe0565b61ffff8167ffffffffffffffff1611610ec3578251610eaa906019611fe0600586901b1617610f60565b508251610e389067ffffffffffffffff83166002610fe0565b63ffffffff8167ffffffffffffffff1611610f08578251610eef90601a611fe0600586901b1617610f60565b508251610e389067ffffffffffffffff83166004610fe0565b8251610f1f90601b611fe0600586901b1617610f60565b508251610e389067ffffffffffffffff83166008610fe0565b604080518082019091526060815260006020820152610f5983838451611065565b9392505050565b6040805180820190915260608152600060208201528251516000610f858260016116ce565b905084602001518210610fa657610fa685610fa18360026116b7565b611154565b8451602083820101858153508051821115610fbf578181525b5093949350505050565b815161022890601f611fe0600585901b1617610f60565b604080518082019091526060815260006020820152835151600061100482856116ce565b905085602001518111156110215761102186610fa18360026116b7565b600060016110318661010061184f565b61103b91906116e1565b90508651828101878319825116178152508051831115611059578281525b50959695505050505050565b604080518082019091526060815260006020820152825182111561108857600080fd5b835151600061109784836116ce565b905085602001518111156110b4576110b486610fa18360026116b7565b8551805183820160200191600091808511156110ce578482525b505050602086015b6020861061110e57805182526110ed6020836116ce565b91506110fa6020826116ce565b90506111076020876116e1565b95506110d6565b5181517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60208890036101000a0190811690199190911617905250849150509392505050565b81516111608383610d97565b50610e388382610f38565b6040518060400160405280611193604051806040016040528060608152602001600081525090565b8152602001600081525090565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715611216576112166111a0565b604052919050565b600067ffffffffffffffff831115611238576112386111a0565b61126960207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f860116016111cf565b905082815283838301111561127d57600080fd5b828260208301376000602084830101529392505050565b600082601f8301126112a557600080fd5b610f598383356020850161121e565b6000806000606084860312156112c957600080fd5b83359250602084013567ffffffffffffffff808211156112e857600080fd5b6112f487838801611294565b9350604086013591508082111561130a57600080fd5b5061131786828701611294565b9150509250925092565b60008083601f84011261133357600080fd5b50813567ffffffffffffffff81111561134b57600080fd5b60208301915083602082850101111561136357600080fd5b9250929050565b803567ffffffffffffffff8116811461138257600080fd5b919050565b60008060008060008060008060a0898b0312156113a357600080fd5b883567ffffffffffffffff808211156113bb57600080fd5b6113c78c838d01611321565b909a50985060208b01359150808211156113e057600080fd5b6113ec8c838d01611321565b909850965060408b013591508082111561140557600080fd5b818b0191508b601f83011261141957600080fd5b81358181111561142857600080fd5b8c60208260051b850101111561143d57600080fd5b60208301965080955050505061145560608a0161136a565b9150608089013590509295985092959890939650565b60006020828403121561147d57600080fd5b813573ffffffffffffffffffffffffffffffffffffffff81168114610f5957600080fd5b600067ffffffffffffffff808411156114bc576114bc6111a0565b8360051b60206114cd8183016111cf565b8681529185019181810190368411156114e557600080fd5b865b8481101561152d578035868111156114ff5760008081fd5b880136601f8201126115115760008081fd5b61151f36823587840161121e565b8452509183019183016114e7565b50979650505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036115f7576115f7611597565b5060010190565b67ffffffffffffffff861681526000602060a08184015286518060a085015260005b8181101561163c5788810183015185820160c001528201611620565b50600060c0828601015260c07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010192505050611685604083018661ffff169052565b63ffffffff939093166060820152608001529392505050565b6000602082840312156116b057600080fd5b5051919050565b8082028115828204841417610e0b57610e0b611597565b80820180821115610e0b57610e0b611597565b81810381811115610e0b57610e0b611597565b60008261172a577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500690565b600181815b8085111561178857817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0482111561176e5761176e611597565b8085161561177b57918102915b93841c9390800290611734565b509250929050565b60008261179f57506001610e0b565b816117ac57506000610e0b565b81600181146117c257600281146117cc576117e8565b6001915050610e0b565b60ff8411156117dd576117dd611597565b50506001821b610e0b565b5060208310610133831016604e8410600b841016171561180b575081810a610e0b565b611815838361172f565b807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0482111561184757611847611597565b029392505050565b6000610f59838361179056fea164736f6c6343000813000a",
}

var FunctionsLoadTestClientABI = FunctionsLoadTestClientMetaData.ABI

var FunctionsLoadTestClientBin = FunctionsLoadTestClientMetaData.Bin

func DeployFunctionsLoadTestClient(auth *bind.TransactOpts, backend bind.ContractBackend, router common.Address) (common.Address, *types.Transaction, *FunctionsLoadTestClient, error) {
	parsed, err := FunctionsLoadTestClientMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(FunctionsLoadTestClientBin), backend, router)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &FunctionsLoadTestClient{FunctionsLoadTestClientCaller: FunctionsLoadTestClientCaller{contract: contract}, FunctionsLoadTestClientTransactor: FunctionsLoadTestClientTransactor{contract: contract}, FunctionsLoadTestClientFilterer: FunctionsLoadTestClientFilterer{contract: contract}}, nil
}

type FunctionsLoadTestClient struct {
	address common.Address
	abi     abi.ABI
	FunctionsLoadTestClientCaller
	FunctionsLoadTestClientTransactor
	FunctionsLoadTestClientFilterer
}

type FunctionsLoadTestClientCaller struct {
	contract *bind.BoundContract
}

type FunctionsLoadTestClientTransactor struct {
	contract *bind.BoundContract
}

type FunctionsLoadTestClientFilterer struct {
	contract *bind.BoundContract
}

type FunctionsLoadTestClientSession struct {
	Contract     *FunctionsLoadTestClient
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type FunctionsLoadTestClientCallerSession struct {
	Contract *FunctionsLoadTestClientCaller
	CallOpts bind.CallOpts
}

type FunctionsLoadTestClientTransactorSession struct {
	Contract     *FunctionsLoadTestClientTransactor
	TransactOpts bind.TransactOpts
}

type FunctionsLoadTestClientRaw struct {
	Contract *FunctionsLoadTestClient
}

type FunctionsLoadTestClientCallerRaw struct {
	Contract *FunctionsLoadTestClientCaller
}

type FunctionsLoadTestClientTransactorRaw struct {
	Contract *FunctionsLoadTestClientTransactor
}

func NewFunctionsLoadTestClient(address common.Address, backend bind.ContractBackend) (*FunctionsLoadTestClient, error) {
	abi, err := abi.JSON(strings.NewReader(FunctionsLoadTestClientABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindFunctionsLoadTestClient(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FunctionsLoadTestClient{address: address, abi: abi, FunctionsLoadTestClientCaller: FunctionsLoadTestClientCaller{contract: contract}, FunctionsLoadTestClientTransactor: FunctionsLoadTestClientTransactor{contract: contract}, FunctionsLoadTestClientFilterer: FunctionsLoadTestClientFilterer{contract: contract}}, nil
}

func NewFunctionsLoadTestClientCaller(address common.Address, caller bind.ContractCaller) (*FunctionsLoadTestClientCaller, error) {
	contract, err := bindFunctionsLoadTestClient(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FunctionsLoadTestClientCaller{contract: contract}, nil
}

func NewFunctionsLoadTestClientTransactor(address common.Address, transactor bind.ContractTransactor) (*FunctionsLoadTestClientTransactor, error) {
	contract, err := bindFunctionsLoadTestClient(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FunctionsLoadTestClientTransactor{contract: contract}, nil
}

func NewFunctionsLoadTestClientFilterer(address common.Address, filterer bind.ContractFilterer) (*FunctionsLoadTestClientFilterer, error) {
	contract, err := bindFunctionsLoadTestClient(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FunctionsLoadTestClientFilterer{contract: contract}, nil
}

func bindFunctionsLoadTestClient(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FunctionsLoadTestClientMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FunctionsLoadTestClient.Contract.FunctionsLoadTestClientCaller.contract.Call(opts, result, method, params...)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsLoadTestClient.Contract.FunctionsLoadTestClientTransactor.contract.Transfer(opts)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FunctionsLoadTestClient.Contract.FunctionsLoadTestClientTransactor.contract.Transact(opts, method, params...)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FunctionsLoadTestClient.Contract.contract.Call(opts, result, method, params...)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsLoadTestClient.Contract.contract.Transfer(opts)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FunctionsLoadTestClient.Contract.contract.Transact(opts, method, params...)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientCaller) MAXCALLBACKGAS(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _FunctionsLoadTestClient.contract.Call(opts, &out, "MAX_CALLBACK_GAS")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientSession) MAXCALLBACKGAS() (uint32, error) {
	return _FunctionsLoadTestClient.Contract.MAXCALLBACKGAS(&_FunctionsLoadTestClient.CallOpts)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientCallerSession) MAXCALLBACKGAS() (uint32, error) {
	return _FunctionsLoadTestClient.Contract.MAXCALLBACKGAS(&_FunctionsLoadTestClient.CallOpts)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FunctionsLoadTestClient.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientSession) Owner() (common.Address, error) {
	return _FunctionsLoadTestClient.Contract.Owner(&_FunctionsLoadTestClient.CallOpts)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientCallerSession) Owner() (common.Address, error) {
	return _FunctionsLoadTestClient.Contract.Owner(&_FunctionsLoadTestClient.CallOpts)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientCaller) SLastError(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _FunctionsLoadTestClient.contract.Call(opts, &out, "s_lastError")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientSession) SLastError() ([32]byte, error) {
	return _FunctionsLoadTestClient.Contract.SLastError(&_FunctionsLoadTestClient.CallOpts)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientCallerSession) SLastError() ([32]byte, error) {
	return _FunctionsLoadTestClient.Contract.SLastError(&_FunctionsLoadTestClient.CallOpts)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientCaller) SLastErrorLength(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _FunctionsLoadTestClient.contract.Call(opts, &out, "s_lastErrorLength")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientSession) SLastErrorLength() (uint32, error) {
	return _FunctionsLoadTestClient.Contract.SLastErrorLength(&_FunctionsLoadTestClient.CallOpts)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientCallerSession) SLastErrorLength() (uint32, error) {
	return _FunctionsLoadTestClient.Contract.SLastErrorLength(&_FunctionsLoadTestClient.CallOpts)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientCaller) SLastRequestId(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _FunctionsLoadTestClient.contract.Call(opts, &out, "s_lastRequestId")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientSession) SLastRequestId() ([32]byte, error) {
	return _FunctionsLoadTestClient.Contract.SLastRequestId(&_FunctionsLoadTestClient.CallOpts)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientCallerSession) SLastRequestId() ([32]byte, error) {
	return _FunctionsLoadTestClient.Contract.SLastRequestId(&_FunctionsLoadTestClient.CallOpts)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientCaller) SLastResponse(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _FunctionsLoadTestClient.contract.Call(opts, &out, "s_lastResponse")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientSession) SLastResponse() ([32]byte, error) {
	return _FunctionsLoadTestClient.Contract.SLastResponse(&_FunctionsLoadTestClient.CallOpts)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientCallerSession) SLastResponse() ([32]byte, error) {
	return _FunctionsLoadTestClient.Contract.SLastResponse(&_FunctionsLoadTestClient.CallOpts)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientCaller) SLastResponseLength(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _FunctionsLoadTestClient.contract.Call(opts, &out, "s_lastResponseLength")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientSession) SLastResponseLength() (uint32, error) {
	return _FunctionsLoadTestClient.Contract.SLastResponseLength(&_FunctionsLoadTestClient.CallOpts)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientCallerSession) SLastResponseLength() (uint32, error) {
	return _FunctionsLoadTestClient.Contract.SLastResponseLength(&_FunctionsLoadTestClient.CallOpts)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsLoadTestClient.contract.Transact(opts, "acceptOwnership")
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientSession) AcceptOwnership() (*types.Transaction, error) {
	return _FunctionsLoadTestClient.Contract.AcceptOwnership(&_FunctionsLoadTestClient.TransactOpts)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _FunctionsLoadTestClient.Contract.AcceptOwnership(&_FunctionsLoadTestClient.TransactOpts)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientTransactor) HandleOracleFulfillment(opts *bind.TransactOpts, requestId [32]byte, response []byte, err []byte) (*types.Transaction, error) {
	return _FunctionsLoadTestClient.contract.Transact(opts, "handleOracleFulfillment", requestId, response, err)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientSession) HandleOracleFulfillment(requestId [32]byte, response []byte, err []byte) (*types.Transaction, error) {
	return _FunctionsLoadTestClient.Contract.HandleOracleFulfillment(&_FunctionsLoadTestClient.TransactOpts, requestId, response, err)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientTransactorSession) HandleOracleFulfillment(requestId [32]byte, response []byte, err []byte) (*types.Transaction, error) {
	return _FunctionsLoadTestClient.Contract.HandleOracleFulfillment(&_FunctionsLoadTestClient.TransactOpts, requestId, response, err)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientTransactor) SendRequest(opts *bind.TransactOpts, source string, encryptedSecretsReferences []byte, args []string, subscriptionId uint64, jobId [32]byte) (*types.Transaction, error) {
	return _FunctionsLoadTestClient.contract.Transact(opts, "sendRequest", source, encryptedSecretsReferences, args, subscriptionId, jobId)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientSession) SendRequest(source string, encryptedSecretsReferences []byte, args []string, subscriptionId uint64, jobId [32]byte) (*types.Transaction, error) {
	return _FunctionsLoadTestClient.Contract.SendRequest(&_FunctionsLoadTestClient.TransactOpts, source, encryptedSecretsReferences, args, subscriptionId, jobId)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientTransactorSession) SendRequest(source string, encryptedSecretsReferences []byte, args []string, subscriptionId uint64, jobId [32]byte) (*types.Transaction, error) {
	return _FunctionsLoadTestClient.Contract.SendRequest(&_FunctionsLoadTestClient.TransactOpts, source, encryptedSecretsReferences, args, subscriptionId, jobId)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _FunctionsLoadTestClient.contract.Transact(opts, "transferOwnership", to)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _FunctionsLoadTestClient.Contract.TransferOwnership(&_FunctionsLoadTestClient.TransactOpts, to)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _FunctionsLoadTestClient.Contract.TransferOwnership(&_FunctionsLoadTestClient.TransactOpts, to)
}

type FunctionsLoadTestClientOwnershipTransferRequestedIterator struct {
	Event *FunctionsLoadTestClientOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsLoadTestClientOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsLoadTestClientOwnershipTransferRequested)
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
		it.Event = new(FunctionsLoadTestClientOwnershipTransferRequested)
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

func (it *FunctionsLoadTestClientOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *FunctionsLoadTestClientOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsLoadTestClientOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsLoadTestClientOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsLoadTestClient.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsLoadTestClientOwnershipTransferRequestedIterator{contract: _FunctionsLoadTestClient.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *FunctionsLoadTestClientOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsLoadTestClient.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsLoadTestClientOwnershipTransferRequested)
				if err := _FunctionsLoadTestClient.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_FunctionsLoadTestClient *FunctionsLoadTestClientFilterer) ParseOwnershipTransferRequested(log types.Log) (*FunctionsLoadTestClientOwnershipTransferRequested, error) {
	event := new(FunctionsLoadTestClientOwnershipTransferRequested)
	if err := _FunctionsLoadTestClient.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsLoadTestClientOwnershipTransferredIterator struct {
	Event *FunctionsLoadTestClientOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsLoadTestClientOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsLoadTestClientOwnershipTransferred)
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
		it.Event = new(FunctionsLoadTestClientOwnershipTransferred)
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

func (it *FunctionsLoadTestClientOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *FunctionsLoadTestClientOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsLoadTestClientOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsLoadTestClientOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsLoadTestClient.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsLoadTestClientOwnershipTransferredIterator{contract: _FunctionsLoadTestClient.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FunctionsLoadTestClientOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsLoadTestClient.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsLoadTestClientOwnershipTransferred)
				if err := _FunctionsLoadTestClient.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_FunctionsLoadTestClient *FunctionsLoadTestClientFilterer) ParseOwnershipTransferred(log types.Log) (*FunctionsLoadTestClientOwnershipTransferred, error) {
	event := new(FunctionsLoadTestClientOwnershipTransferred)
	if err := _FunctionsLoadTestClient.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsLoadTestClientRequestFulfilledIterator struct {
	Event *FunctionsLoadTestClientRequestFulfilled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsLoadTestClientRequestFulfilledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsLoadTestClientRequestFulfilled)
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
		it.Event = new(FunctionsLoadTestClientRequestFulfilled)
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

func (it *FunctionsLoadTestClientRequestFulfilledIterator) Error() error {
	return it.fail
}

func (it *FunctionsLoadTestClientRequestFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsLoadTestClientRequestFulfilled struct {
	Id  [32]byte
	Raw types.Log
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientFilterer) FilterRequestFulfilled(opts *bind.FilterOpts, id [][32]byte) (*FunctionsLoadTestClientRequestFulfilledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _FunctionsLoadTestClient.contract.FilterLogs(opts, "RequestFulfilled", idRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsLoadTestClientRequestFulfilledIterator{contract: _FunctionsLoadTestClient.contract, event: "RequestFulfilled", logs: logs, sub: sub}, nil
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientFilterer) WatchRequestFulfilled(opts *bind.WatchOpts, sink chan<- *FunctionsLoadTestClientRequestFulfilled, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _FunctionsLoadTestClient.contract.WatchLogs(opts, "RequestFulfilled", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsLoadTestClientRequestFulfilled)
				if err := _FunctionsLoadTestClient.contract.UnpackLog(event, "RequestFulfilled", log); err != nil {
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

func (_FunctionsLoadTestClient *FunctionsLoadTestClientFilterer) ParseRequestFulfilled(log types.Log) (*FunctionsLoadTestClientRequestFulfilled, error) {
	event := new(FunctionsLoadTestClientRequestFulfilled)
	if err := _FunctionsLoadTestClient.contract.UnpackLog(event, "RequestFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsLoadTestClientRequestSentIterator struct {
	Event *FunctionsLoadTestClientRequestSent

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsLoadTestClientRequestSentIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsLoadTestClientRequestSent)
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
		it.Event = new(FunctionsLoadTestClientRequestSent)
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

func (it *FunctionsLoadTestClientRequestSentIterator) Error() error {
	return it.fail
}

func (it *FunctionsLoadTestClientRequestSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsLoadTestClientRequestSent struct {
	Id  [32]byte
	Raw types.Log
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientFilterer) FilterRequestSent(opts *bind.FilterOpts, id [][32]byte) (*FunctionsLoadTestClientRequestSentIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _FunctionsLoadTestClient.contract.FilterLogs(opts, "RequestSent", idRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsLoadTestClientRequestSentIterator{contract: _FunctionsLoadTestClient.contract, event: "RequestSent", logs: logs, sub: sub}, nil
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientFilterer) WatchRequestSent(opts *bind.WatchOpts, sink chan<- *FunctionsLoadTestClientRequestSent, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _FunctionsLoadTestClient.contract.WatchLogs(opts, "RequestSent", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsLoadTestClientRequestSent)
				if err := _FunctionsLoadTestClient.contract.UnpackLog(event, "RequestSent", log); err != nil {
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

func (_FunctionsLoadTestClient *FunctionsLoadTestClientFilterer) ParseRequestSent(log types.Log) (*FunctionsLoadTestClientRequestSent, error) {
	event := new(FunctionsLoadTestClientRequestSent)
	if err := _FunctionsLoadTestClient.contract.UnpackLog(event, "RequestSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClient) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _FunctionsLoadTestClient.abi.Events["OwnershipTransferRequested"].ID:
		return _FunctionsLoadTestClient.ParseOwnershipTransferRequested(log)
	case _FunctionsLoadTestClient.abi.Events["OwnershipTransferred"].ID:
		return _FunctionsLoadTestClient.ParseOwnershipTransferred(log)
	case _FunctionsLoadTestClient.abi.Events["RequestFulfilled"].ID:
		return _FunctionsLoadTestClient.ParseRequestFulfilled(log)
	case _FunctionsLoadTestClient.abi.Events["RequestSent"].ID:
		return _FunctionsLoadTestClient.ParseRequestSent(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (FunctionsLoadTestClientOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (FunctionsLoadTestClientOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (FunctionsLoadTestClientRequestFulfilled) Topic() common.Hash {
	return common.HexToHash("0x85e1543bf2f84fe80c6badbce3648c8539ad1df4d2b3d822938ca0538be727e6")
}

func (FunctionsLoadTestClientRequestSent) Topic() common.Hash {
	return common.HexToHash("0x1131472297a800fee664d1d89cfa8f7676ff07189ecc53f80bbb5f4969099db8")
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClient) Address() common.Address {
	return _FunctionsLoadTestClient.address
}

type FunctionsLoadTestClientInterface interface {
	MAXCALLBACKGAS(opts *bind.CallOpts) (uint32, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SLastError(opts *bind.CallOpts) ([32]byte, error)

	SLastErrorLength(opts *bind.CallOpts) (uint32, error)

	SLastRequestId(opts *bind.CallOpts) ([32]byte, error)

	SLastResponse(opts *bind.CallOpts) ([32]byte, error)

	SLastResponseLength(opts *bind.CallOpts) (uint32, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	HandleOracleFulfillment(opts *bind.TransactOpts, requestId [32]byte, response []byte, err []byte) (*types.Transaction, error)

	SendRequest(opts *bind.TransactOpts, source string, encryptedSecretsReferences []byte, args []string, subscriptionId uint64, jobId [32]byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsLoadTestClientOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *FunctionsLoadTestClientOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*FunctionsLoadTestClientOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsLoadTestClientOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FunctionsLoadTestClientOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*FunctionsLoadTestClientOwnershipTransferred, error)

	FilterRequestFulfilled(opts *bind.FilterOpts, id [][32]byte) (*FunctionsLoadTestClientRequestFulfilledIterator, error)

	WatchRequestFulfilled(opts *bind.WatchOpts, sink chan<- *FunctionsLoadTestClientRequestFulfilled, id [][32]byte) (event.Subscription, error)

	ParseRequestFulfilled(log types.Log) (*FunctionsLoadTestClientRequestFulfilled, error)

	FilterRequestSent(opts *bind.FilterOpts, id [][32]byte) (*FunctionsLoadTestClientRequestSentIterator, error)

	WatchRequestSent(opts *bind.WatchOpts, sink chan<- *FunctionsLoadTestClientRequestSent, id [][32]byte) (event.Subscription, error)

	ParseRequestSent(log types.Log) (*FunctionsLoadTestClientRequestSent, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
