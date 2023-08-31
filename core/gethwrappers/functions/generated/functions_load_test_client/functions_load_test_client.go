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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"EmptyArgs\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptySecrets\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptySource\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoInlineSecrets\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyRouterCanFulfill\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"RequestFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"RequestSent\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"MAX_CALLBACK_GAS\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStats\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"response\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"handleOracleFulfillment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastError\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastRequestID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastResponse\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"resetCounters\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"source\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"encryptedSecretsReferences\",\"type\":\"bytes\"},{\"internalType\":\"string[]\",\"name\":\"args\",\"type\":\"string[]\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"jobId\",\"type\":\"bytes32\"}],\"name\":\"sendRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalEmptyResponses\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalFailedResponses\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalRequests\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSucceededResponses\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b5060405162001c3538038062001c35833981016040819052620000349162000180565b6001600160a01b0381166080523380600081620000985760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000cb57620000cb81620000d5565b50505050620001b2565b336001600160a01b038216036200012f5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016200008f565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000602082840312156200019357600080fd5b81516001600160a01b0381168114620001ab57600080fd5b9392505050565b608051611a60620001d56000396000818161026e0152610bc50152611a606000f3fe608060405234801561001057600080fd5b50600436106100f55760003560e01c80636d9809a011610097578063c59d484711610066578063c59d4847146101e3578063c9429e2a1461021b578063e4ed53a91461023b578063f2fde38b1461024357600080fd5b80636d9809a01461019957806379ba5097146101a35780638aea61dc146101ab5780638da5cb5b146101bb57600080fd5b806347c03186116100d357806347c031861461015c5780635c1d92e9146101655780635fa353e71461017d57806362747e421461019057600080fd5b80630ca76175146100fa57806329f0de3f1461010f5780632ab424da1461012b575b600080fd5b61010d610108366004611488565b610256565b005b61011860045481565b6040519081526020015b60405180910390f35b6005546101479068010000000000000000900463ffffffff1681565b60405163ffffffff9091168152602001610122565b61011860025481565b60055461014790640100000000900463ffffffff1681565b61010d61018b36600461155b565b6102d5565b61011860035481565b6101476201117081565b61010d610429565b6005546101479063ffffffff1681565b60005460405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610122565b6101eb61052b565b6040805163ffffffff95861681529385166020850152918416918301919091529091166060820152608001610122565b600554610147906c01000000000000000000000000900463ffffffff1681565b61010d610578565b61010d61025136600461163f565b6105aa565b3373ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016146102c5576040517fc6829f8300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6102d08383836105be565b505050565b6102dd6106c6565b6103266040805161010081019091528060008152602001600081526020016000815260200160608152602001606081526020016060815260200160608152602001606081525090565b61036889898080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525085939250506107499050565b85156103b0576103b087878080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250859392505061075a9050565b83156103ca576103ca6103c38587611675565b82906107a4565b6103e16103d6826107e7565b846201117085610bc0565b600255600580546001919060009061040090849063ffffffff1661173c565b92506101000a81548163ffffffff021916908363ffffffff160217905550505050505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146104af576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6000806000806105396106c6565b505060055463ffffffff80821694680100000000000000008304821694506c010000000000000000000000008304821693506401000000009092041690565b6105806106c6565b600580547fffffffffffffffffffffffffffffffff00000000000000000000000000000000169055565b6105b26106c6565b6105bb81610c9f565b50565b60028390556105cc82610d94565b6003556105d881610d94565b6004558151600003610625576001600560048282829054906101000a900463ffffffff16610606919061173c565b92506101000a81548163ffffffff021916908363ffffffff1602179055505b80511561066d5760016005600c8282829054906101000a900463ffffffff1661064e919061173c565b92506101000a81548163ffffffff021916908363ffffffff1602179055505b81511580159061067c57508051155b156102d0576001600560088282829054906101000a900463ffffffff166106a3919061173c565b92506101000a81548163ffffffff021916908363ffffffff160217905550505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610747576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016104a6565b565b6107568260008084610e16565b5050565b8051600003610795576040517fe889636f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60016020830152608090910152565b80516000036107df576040517ffe936cb700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60c090910152565b606060006107f6610100610ead565b90506108406040518060400160405280600c81526020017f636f64654c6f636174696f6e000000000000000000000000000000000000000081525082610ece90919063ffffffff16565b825161085e90600281111561085757610857611760565b8290610ee7565b60408051808201909152600881527f6c616e6775616765000000000000000000000000000000000000000000000000602082015261089d908290610ece565b60408301516108b490801561085757610857611760565b60408051808201909152600681527f736f75726365000000000000000000000000000000000000000000000000000060208201526108f3908290610ece565b6060830151610903908290610ece565b60a0830151511561095d5760408051808201909152601081527f726571756573745369676e617475726500000000000000000000000000000000602082015261094d908290610ece565b60a083015161095d908290610f1c565b60c08301515115610a0a5760408051808201909152600481527f617267730000000000000000000000000000000000000000000000000000000060208201526109a7908290610ece565b6109b081610f29565b60005b8360c0015151811015610a00576109f08460c0015182815181106109d9576109d961178f565b602002602001015183610ece90919063ffffffff16565b6109f9816117be565b90506109b3565b50610a0a81610f4d565b60808301515115610b0b57600083602001516002811115610a2d57610a2d611760565b03610a64576040517fa80d31f700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60408051808201909152600f81527f736563726574734c6f636174696f6e00000000000000000000000000000000006020820152610aa3908290610ece565b610abc8360200151600281111561085757610857611760565b60408051808201909152600781527f73656372657473000000000000000000000000000000000000000000000000006020820152610afb908290610ece565b6080830151610b0b908290610f1c565b60e08301515115610bb85760408051808201909152600981527f62797465734172677300000000000000000000000000000000000000000000006020820152610b55908290610ece565b610b5e81610f29565b60005b8360e0015151811015610bae57610b9e8460e001518281518110610b8757610b8761178f565b602002602001015183610f1c90919063ffffffff16565b610ba7816117be565b9050610b61565b50610bb881610f4d565b515192915050565b6000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663461d27628688600188886040518663ffffffff1660e01b8152600401610c259594939291906117f6565b6020604051808303816000875af1158015610c44573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c689190611896565b60405190915081907f1131472297a800fee664d1d89cfa8f7676ff07189ecc53f80bbb5f4969099db890600090a295945050505050565b3373ffffffffffffffffffffffffffffffffffffffff821603610d1e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016104a6565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60008060209050602083511015610da9575081515b60005b81811015610e0f57610dbf8160086118af565b848281518110610dd157610dd161178f565b01602001517fff0000000000000000000000000000000000000000000000000000000000000016901c9290921791610e08816117be565b9050610dac565b5050919050565b8051600003610e51576040517f22ce3edd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b83836002811115610e6457610e64611760565b90816002811115610e7757610e77611760565b90525060408401828015610e8d57610e8d611760565b90818015610e9d57610e9d611760565b9052506060909301929092525050565b610eb561133f565b8051610ec19083610f6b565b5060006020820152919050565b610edb8260038351610fe5565b81516102d0908261110c565b8151610ef49060c2611134565b506107568282604051602001610f0c91815260200190565b6040516020818303038152906040525b610edb8260028351610fe5565b610f3481600461119d565b600181602001818151610f4791906118c6565b90525050565b610f5881600761119d565b600181602001818151610f4791906118d9565b604080518082019091526060815260006020820152610f8b6020836118ec565b15610fb357610f9b6020836118ec565b610fa69060206118d9565b610fb090836118c6565b91505b602080840183905260405180855260008152908184010181811015610fd757600080fd5b604052508290505b92915050565b60178167ffffffffffffffff161161101257825161100c9060e0600585901b168317611134565b50505050565b60ff8167ffffffffffffffff161161105457825161103b906018611fe0600586901b1617611134565b50825161100c9067ffffffffffffffff831660016111b4565b61ffff8167ffffffffffffffff161161109757825161107e906019611fe0600586901b1617611134565b50825161100c9067ffffffffffffffff831660026111b4565b63ffffffff8167ffffffffffffffff16116110dc5782516110c390601a611fe0600586901b1617611134565b50825161100c9067ffffffffffffffff831660046111b4565b82516110f390601b611fe0600586901b1617611134565b50825161100c9067ffffffffffffffff831660086111b4565b60408051808201909152606081526000602082015261112d83838451611239565b9392505050565b60408051808201909152606081526000602082015282515160006111598260016118c6565b90508460200151821061117a5761117a856111758360026118af565b611328565b8451602083820101858153508051821115611193578181525b5093949350505050565b81516102d090601f611fe0600585901b1617611134565b60408051808201909152606081526000602082015283515160006111d882856118c6565b905085602001518111156111f5576111f5866111758360026118af565b6000600161120586610100611a47565b61120f91906118d9565b9050865182810187831982511617815250805183111561122d578281525b50959695505050505050565b604080518082019091526060815260006020820152825182111561125c57600080fd5b835151600061126b84836118c6565b9050856020015181111561128857611288866111758360026118af565b8551805183820160200191600091808511156112a2578482525b505050602086015b602086106112e257805182526112c16020836118c6565b91506112ce6020826118c6565b90506112db6020876118d9565b95506112aa565b5181517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60208890036101000a0190811690199190911617905250849150509392505050565b81516113348383610f6b565b5061100c838261110c565b6040518060400160405280611367604051806040016040528060608152602001600081525090565b8152602001600081525090565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156113ea576113ea611374565b604052919050565b600067ffffffffffffffff83111561140c5761140c611374565b61143d60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f860116016113a3565b905082815283838301111561145157600080fd5b828260208301376000602084830101529392505050565b600082601f83011261147957600080fd5b61112d838335602085016113f2565b60008060006060848603121561149d57600080fd5b83359250602084013567ffffffffffffffff808211156114bc57600080fd5b6114c887838801611468565b935060408601359150808211156114de57600080fd5b506114eb86828701611468565b9150509250925092565b60008083601f84011261150757600080fd5b50813567ffffffffffffffff81111561151f57600080fd5b60208301915083602082850101111561153757600080fd5b9250929050565b803567ffffffffffffffff8116811461155657600080fd5b919050565b60008060008060008060008060a0898b03121561157757600080fd5b883567ffffffffffffffff8082111561158f57600080fd5b61159b8c838d016114f5565b909a50985060208b01359150808211156115b457600080fd5b6115c08c838d016114f5565b909850965060408b01359150808211156115d957600080fd5b818b0191508b601f8301126115ed57600080fd5b8135818111156115fc57600080fd5b8c60208260051b850101111561161157600080fd5b60208301965080955050505061162960608a0161153e565b9150608089013590509295985092959890939650565b60006020828403121561165157600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461112d57600080fd5b600067ffffffffffffffff8084111561169057611690611374565b8360051b60206116a18183016113a3565b8681529185019181810190368411156116b957600080fd5b865b84811015611701578035868111156116d35760008081fd5b880136601f8201126116e55760008081fd5b6116f33682358784016113f2565b8452509183019183016116bb565b50979650505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b63ffffffff8181168382160190808211156117595761175961170d565b5092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036117ef576117ef61170d565b5060010190565b67ffffffffffffffff861681526000602060a08184015286518060a085015260005b818110156118345788810183015185820160c001528201611818565b50600060c0828601015260c07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505061187d604083018661ffff169052565b63ffffffff939093166060820152608001529392505050565b6000602082840312156118a857600080fd5b5051919050565b8082028115828204841417610fdf57610fdf61170d565b80820180821115610fdf57610fdf61170d565b81810381811115610fdf57610fdf61170d565b600082611922577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500690565b600181815b8085111561198057817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048211156119665761196661170d565b8085161561197357918102915b93841c939080029061192c565b509250929050565b60008261199757506001610fdf565b816119a457506000610fdf565b81600181146119ba57600281146119c4576119e0565b6001915050610fdf565b60ff8411156119d5576119d561170d565b50506001821b610fdf565b5060208310610133831016604e8410600b8410161715611a03575081810a610fdf565b611a0d8383611927565b807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04821115611a3f57611a3f61170d565b029392505050565b600061112d838361198856fea164736f6c6343000813000a",
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

func (_FunctionsLoadTestClient *FunctionsLoadTestClientCaller) GetStats(opts *bind.CallOpts) (uint32, uint32, uint32, uint32, error) {
	var out []interface{}
	err := _FunctionsLoadTestClient.contract.Call(opts, &out, "getStats")

	if err != nil {
		return *new(uint32), *new(uint32), *new(uint32), *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)
	out1 := *abi.ConvertType(out[1], new(uint32)).(*uint32)
	out2 := *abi.ConvertType(out[2], new(uint32)).(*uint32)
	out3 := *abi.ConvertType(out[3], new(uint32)).(*uint32)

	return out0, out1, out2, out3, err

}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientSession) GetStats() (uint32, uint32, uint32, uint32, error) {
	return _FunctionsLoadTestClient.Contract.GetStats(&_FunctionsLoadTestClient.CallOpts)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientCallerSession) GetStats() (uint32, uint32, uint32, uint32, error) {
	return _FunctionsLoadTestClient.Contract.GetStats(&_FunctionsLoadTestClient.CallOpts)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientCaller) LastError(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _FunctionsLoadTestClient.contract.Call(opts, &out, "lastError")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientSession) LastError() ([32]byte, error) {
	return _FunctionsLoadTestClient.Contract.LastError(&_FunctionsLoadTestClient.CallOpts)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientCallerSession) LastError() ([32]byte, error) {
	return _FunctionsLoadTestClient.Contract.LastError(&_FunctionsLoadTestClient.CallOpts)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientCaller) LastRequestID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _FunctionsLoadTestClient.contract.Call(opts, &out, "lastRequestID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientSession) LastRequestID() ([32]byte, error) {
	return _FunctionsLoadTestClient.Contract.LastRequestID(&_FunctionsLoadTestClient.CallOpts)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientCallerSession) LastRequestID() ([32]byte, error) {
	return _FunctionsLoadTestClient.Contract.LastRequestID(&_FunctionsLoadTestClient.CallOpts)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientCaller) LastResponse(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _FunctionsLoadTestClient.contract.Call(opts, &out, "lastResponse")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientSession) LastResponse() ([32]byte, error) {
	return _FunctionsLoadTestClient.Contract.LastResponse(&_FunctionsLoadTestClient.CallOpts)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientCallerSession) LastResponse() ([32]byte, error) {
	return _FunctionsLoadTestClient.Contract.LastResponse(&_FunctionsLoadTestClient.CallOpts)
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

func (_FunctionsLoadTestClient *FunctionsLoadTestClientCaller) TotalEmptyResponses(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _FunctionsLoadTestClient.contract.Call(opts, &out, "totalEmptyResponses")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientSession) TotalEmptyResponses() (uint32, error) {
	return _FunctionsLoadTestClient.Contract.TotalEmptyResponses(&_FunctionsLoadTestClient.CallOpts)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientCallerSession) TotalEmptyResponses() (uint32, error) {
	return _FunctionsLoadTestClient.Contract.TotalEmptyResponses(&_FunctionsLoadTestClient.CallOpts)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientCaller) TotalFailedResponses(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _FunctionsLoadTestClient.contract.Call(opts, &out, "totalFailedResponses")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientSession) TotalFailedResponses() (uint32, error) {
	return _FunctionsLoadTestClient.Contract.TotalFailedResponses(&_FunctionsLoadTestClient.CallOpts)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientCallerSession) TotalFailedResponses() (uint32, error) {
	return _FunctionsLoadTestClient.Contract.TotalFailedResponses(&_FunctionsLoadTestClient.CallOpts)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientCaller) TotalRequests(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _FunctionsLoadTestClient.contract.Call(opts, &out, "totalRequests")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientSession) TotalRequests() (uint32, error) {
	return _FunctionsLoadTestClient.Contract.TotalRequests(&_FunctionsLoadTestClient.CallOpts)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientCallerSession) TotalRequests() (uint32, error) {
	return _FunctionsLoadTestClient.Contract.TotalRequests(&_FunctionsLoadTestClient.CallOpts)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientCaller) TotalSucceededResponses(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _FunctionsLoadTestClient.contract.Call(opts, &out, "totalSucceededResponses")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientSession) TotalSucceededResponses() (uint32, error) {
	return _FunctionsLoadTestClient.Contract.TotalSucceededResponses(&_FunctionsLoadTestClient.CallOpts)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientCallerSession) TotalSucceededResponses() (uint32, error) {
	return _FunctionsLoadTestClient.Contract.TotalSucceededResponses(&_FunctionsLoadTestClient.CallOpts)
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

func (_FunctionsLoadTestClient *FunctionsLoadTestClientTransactor) ResetCounters(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsLoadTestClient.contract.Transact(opts, "resetCounters")
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientSession) ResetCounters() (*types.Transaction, error) {
	return _FunctionsLoadTestClient.Contract.ResetCounters(&_FunctionsLoadTestClient.TransactOpts)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientTransactorSession) ResetCounters() (*types.Transaction, error) {
	return _FunctionsLoadTestClient.Contract.ResetCounters(&_FunctionsLoadTestClient.TransactOpts)
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

	GetStats(opts *bind.CallOpts) (uint32, uint32, uint32, uint32, error)

	LastError(opts *bind.CallOpts) ([32]byte, error)

	LastRequestID(opts *bind.CallOpts) ([32]byte, error)

	LastResponse(opts *bind.CallOpts) ([32]byte, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TotalEmptyResponses(opts *bind.CallOpts) (uint32, error)

	TotalFailedResponses(opts *bind.CallOpts) (uint32, error)

	TotalRequests(opts *bind.CallOpts) (uint32, error)

	TotalSucceededResponses(opts *bind.CallOpts) (uint32, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	HandleOracleFulfillment(opts *bind.TransactOpts, requestId [32]byte, response []byte, err []byte) (*types.Transaction, error)

	ResetCounters(opts *bind.TransactOpts) (*types.Transaction, error)

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
