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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"EmptyArgs\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptySecrets\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptySource\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoInlineSecrets\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyRouterCanFulfill\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"RequestFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"RequestSent\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"MAX_CALLBACK_GAS\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStats\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"response\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"handleOracleFulfillment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastError\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastRequestID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastResponse\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"resetStats\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"times\",\"type\":\"uint32\"},{\"internalType\":\"string\",\"name\":\"source\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"encryptedSecretsReferences\",\"type\":\"bytes\"},{\"internalType\":\"string[]\",\"name\":\"args\",\"type\":\"string[]\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"jobId\",\"type\":\"bytes32\"}],\"name\":\"sendRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalEmptyResponses\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalFailedResponses\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalRequests\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSucceededResponses\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b5060405162001cc238038062001cc2833981016040819052620000349162000180565b6001600160a01b0381166080523380600081620000985760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000cb57620000cb81620000d5565b50505050620001b2565b336001600160a01b038216036200012f5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016200008f565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000602082840312156200019357600080fd5b81516001600160a01b0381168114620001ab57600080fd5b9392505050565b608051611aed620001d56000396000818161027e0152610c220152611aed6000f3fe608060405234801561001057600080fd5b50600436106100f55760003560e01c8063724ec8a2116100975780638da5cb5b116100665780638da5cb5b146101c3578063c59d4847146101eb578063c9429e2a14610233578063f2fde38b1461025357600080fd5b8063724ec8a21461019057806379ba509714610198578063887efe94146101a05780638aea61dc146101b357600080fd5b806347c03186116100d357806347c031861461015c5780635c1d92e91461016557806362747e421461017d5780636d9809a01461018657600080fd5b80630ca76175146100fa57806329f0de3f1461010f5780632ab424da1461012b575b600080fd5b61010d6101083660046114e5565b610266565b005b61011860045481565b6040519081526020015b60405180910390f35b6005546101479068010000000000000000900463ffffffff1681565b60405163ffffffff9091168152602001610122565b61011860025481565b60055461014790640100000000900463ffffffff1681565b61011860035481565b6101476201117081565b61010d6102e5565b61010d610326565b61010d6101ae3660046115fd565b610428565b6005546101479063ffffffff1681565b60005460405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610122565b6101f36105a2565b6040805197885260208801969096529486019390935263ffffffff91821660608601528116608085015290811660a08401521660c082015260e001610122565b600554610147906c01000000000000000000000000900463ffffffff1681565b61010d6102613660046116cc565b610607565b3373ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016146102d5576040517fc6829f8300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6102e083838361061b565b505050565b6102ed610723565b600060028190556003819055600455600580547fffffffffffffffffffffffffffffffff00000000000000000000000000000000169055565b60015473ffffffffffffffffffffffffffffffffffffffff1633146103ac576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610430610723565b6104796040805161010081019091528060008152602001600081526020016000815260200160608152602001606081526020016060815260200160608152602001606081525090565b6104bb89898080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525085939250506107a69050565b85156105035761050387878080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525085939250506107b79050565b831561051d5761051d6105168587611702565b8290610801565b60005b8a63ffffffff168110156105955761054561053a83610844565b856201117086610c1d565b600255600580546001919060009061056490849063ffffffff166117c9565b92506101000a81548163ffffffff021916908363ffffffff160217905550808061058d906117ed565b915050610520565b5050505050505050505050565b60008060008060008060006105b5610723565b50506002546003546004546005549298919750955063ffffffff8083169550680100000000000000008304811694506c0100000000000000000000000083048116935064010000000090920490911690565b61060f610723565b61061881610cfc565b50565b600283905561062982610df1565b60035561063581610df1565b6004558151600003610682576001600560048282829054906101000a900463ffffffff1661066391906117c9565b92506101000a81548163ffffffff021916908363ffffffff1602179055505b8051156106ca5760016005600c8282829054906101000a900463ffffffff166106ab91906117c9565b92506101000a81548163ffffffff021916908363ffffffff1602179055505b8151158015906106d957508051155b156102e0576001600560088282829054906101000a900463ffffffff1661070091906117c9565b92506101000a81548163ffffffff021916908363ffffffff160217905550505050565b60005473ffffffffffffffffffffffffffffffffffffffff1633146107a4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016103a3565b565b6107b38260008084610e73565b5050565b80516000036107f2576040517fe889636f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60016020830152608090910152565b805160000361083c576040517ffe936cb700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60c090910152565b60606000610853610100610f0a565b905061089d6040518060400160405280600c81526020017f636f64654c6f636174696f6e000000000000000000000000000000000000000081525082610f2b90919063ffffffff16565b82516108bb9060028111156108b4576108b4611825565b8290610f44565b60408051808201909152600881527f6c616e677561676500000000000000000000000000000000000000000000000060208201526108fa908290610f2b565b60408301516109119080156108b4576108b4611825565b60408051808201909152600681527f736f7572636500000000000000000000000000000000000000000000000000006020820152610950908290610f2b565b6060830151610960908290610f2b565b60a083015151156109ba5760408051808201909152601081527f726571756573745369676e61747572650000000000000000000000000000000060208201526109aa908290610f2b565b60a08301516109ba908290610f79565b60c08301515115610a675760408051808201909152600481527f61726773000000000000000000000000000000000000000000000000000000006020820152610a04908290610f2b565b610a0d81610f86565b60005b8360c0015151811015610a5d57610a4d8460c001518281518110610a3657610a36611854565b602002602001015183610f2b90919063ffffffff16565b610a56816117ed565b9050610a10565b50610a6781610faa565b60808301515115610b6857600083602001516002811115610a8a57610a8a611825565b03610ac1576040517fa80d31f700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60408051808201909152600f81527f736563726574734c6f636174696f6e00000000000000000000000000000000006020820152610b00908290610f2b565b610b19836020015160028111156108b4576108b4611825565b60408051808201909152600781527f73656372657473000000000000000000000000000000000000000000000000006020820152610b58908290610f2b565b6080830151610b68908290610f79565b60e08301515115610c155760408051808201909152600981527f62797465734172677300000000000000000000000000000000000000000000006020820152610bb2908290610f2b565b610bbb81610f86565b60005b8360e0015151811015610c0b57610bfb8460e001518281518110610be457610be4611854565b602002602001015183610f7990919063ffffffff16565b610c04816117ed565b9050610bbe565b50610c1581610faa565b515192915050565b6000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663461d27628688600188886040518663ffffffff1660e01b8152600401610c82959493929190611883565b6020604051808303816000875af1158015610ca1573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610cc59190611923565b60405190915081907f1131472297a800fee664d1d89cfa8f7676ff07189ecc53f80bbb5f4969099db890600090a295945050505050565b3373ffffffffffffffffffffffffffffffffffffffff821603610d7b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016103a3565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60008060209050602083511015610e06575081515b60005b81811015610e6c57610e1c81600861193c565b848281518110610e2e57610e2e611854565b01602001517fff0000000000000000000000000000000000000000000000000000000000000016901c9290921791610e65816117ed565b9050610e09565b5050919050565b8051600003610eae576040517f22ce3edd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b83836002811115610ec157610ec1611825565b90816002811115610ed457610ed4611825565b90525060408401828015610eea57610eea611825565b90818015610efa57610efa611825565b9052506060909301929092525050565b610f1261139c565b8051610f1e9083610fc8565b5060006020820152919050565b610f388260038351611042565b81516102e09082611169565b8151610f519060c2611191565b506107b38282604051602001610f6991815260200190565b6040516020818303038152906040525b610f388260028351611042565b610f918160046111fa565b600181602001818151610fa49190611953565b90525050565b610fb58160076111fa565b600181602001818151610fa49190611966565b604080518082019091526060815260006020820152610fe8602083611979565b1561101057610ff8602083611979565b611003906020611966565b61100d9083611953565b91505b60208084018390526040518085526000815290818401018181101561103457600080fd5b604052508290505b92915050565b60178167ffffffffffffffff161161106f5782516110699060e0600585901b168317611191565b50505050565b60ff8167ffffffffffffffff16116110b1578251611098906018611fe0600586901b1617611191565b5082516110699067ffffffffffffffff83166001611211565b61ffff8167ffffffffffffffff16116110f45782516110db906019611fe0600586901b1617611191565b5082516110699067ffffffffffffffff83166002611211565b63ffffffff8167ffffffffffffffff161161113957825161112090601a611fe0600586901b1617611191565b5082516110699067ffffffffffffffff83166004611211565b825161115090601b611fe0600586901b1617611191565b5082516110699067ffffffffffffffff83166008611211565b60408051808201909152606081526000602082015261118a83838451611296565b9392505050565b60408051808201909152606081526000602082015282515160006111b6826001611953565b9050846020015182106111d7576111d7856111d283600261193c565b611385565b84516020838201018581535080518211156111f0578181525b5093949350505050565b81516102e090601f611fe0600585901b1617611191565b60408051808201909152606081526000602082015283515160006112358285611953565b9050856020015181111561125257611252866111d283600261193c565b6000600161126286610100611ad4565b61126c9190611966565b9050865182810187831982511617815250805183111561128a578281525b50959695505050505050565b60408051808201909152606081526000602082015282518211156112b957600080fd5b83515160006112c88483611953565b905085602001518111156112e5576112e5866111d283600261193c565b8551805183820160200191600091808511156112ff578482525b505050602086015b6020861061133f578051825261131e602083611953565b915061132b602082611953565b9050611338602087611966565b9550611307565b5181517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60208890036101000a0190811690199190911617905250849150509392505050565b81516113918383610fc8565b506110698382611169565b60405180604001604052806113c4604051806040016040528060608152602001600081525090565b8152602001600081525090565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715611447576114476113d1565b604052919050565b600067ffffffffffffffff831115611469576114696113d1565b61149a60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f86011601611400565b90508281528383830111156114ae57600080fd5b828260208301376000602084830101529392505050565b600082601f8301126114d657600080fd5b61118a8383356020850161144f565b6000806000606084860312156114fa57600080fd5b83359250602084013567ffffffffffffffff8082111561151957600080fd5b611525878388016114c5565b9350604086013591508082111561153b57600080fd5b50611548868287016114c5565b9150509250925092565b60008083601f84011261156457600080fd5b50813567ffffffffffffffff81111561157c57600080fd5b60208301915083602082850101111561159457600080fd5b9250929050565b60008083601f8401126115ad57600080fd5b50813567ffffffffffffffff8111156115c557600080fd5b6020830191508360208260051b850101111561159457600080fd5b803567ffffffffffffffff811681146115f857600080fd5b919050565b600080600080600080600080600060c08a8c03121561161b57600080fd5b893563ffffffff8116811461162f57600080fd5b985060208a013567ffffffffffffffff8082111561164c57600080fd5b6116588d838e01611552565b909a50985060408c013591508082111561167157600080fd5b61167d8d838e01611552565b909850965060608c013591508082111561169657600080fd5b506116a38c828d0161159b565b90955093506116b6905060808b016115e0565b915060a08a013590509295985092959850929598565b6000602082840312156116de57600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461118a57600080fd5b600067ffffffffffffffff8084111561171d5761171d6113d1565b8360051b602061172e818301611400565b86815291850191818101903684111561174657600080fd5b865b8481101561178e578035868111156117605760008081fd5b880136601f8201126117725760008081fd5b61178036823587840161144f565b845250918301918301611748565b50979650505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b63ffffffff8181168382160190808211156117e6576117e661179a565b5092915050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361181e5761181e61179a565b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b67ffffffffffffffff861681526000602060a08184015286518060a085015260005b818110156118c15788810183015185820160c0015282016118a5565b50600060c0828601015260c07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505061190a604083018661ffff169052565b63ffffffff939093166060820152608001529392505050565b60006020828403121561193557600080fd5b5051919050565b808202811582820484141761103c5761103c61179a565b8082018082111561103c5761103c61179a565b8181038181111561103c5761103c61179a565b6000826119af577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500690565b600181815b80851115611a0d57817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048211156119f3576119f361179a565b80851615611a0057918102915b93841c93908002906119b9565b509250929050565b600082611a245750600161103c565b81611a315750600061103c565b8160018114611a475760028114611a5157611a6d565b600191505061103c565b60ff841115611a6257611a6261179a565b50506001821b61103c565b5060208310610133831016604e8410600b8410161715611a90575081810a61103c565b611a9a83836119b4565b807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04821115611acc57611acc61179a565b029392505050565b600061118a8383611a1556fea164736f6c6343000813000a",
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

func (_FunctionsLoadTestClient *FunctionsLoadTestClientCaller) GetStats(opts *bind.CallOpts) ([32]byte, [32]byte, [32]byte, uint32, uint32, uint32, uint32, error) {
	var out []interface{}
	err := _FunctionsLoadTestClient.contract.Call(opts, &out, "getStats")

	if err != nil {
		return *new([32]byte), *new([32]byte), *new([32]byte), *new(uint32), *new(uint32), *new(uint32), *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	out1 := *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	out2 := *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)
	out3 := *abi.ConvertType(out[3], new(uint32)).(*uint32)
	out4 := *abi.ConvertType(out[4], new(uint32)).(*uint32)
	out5 := *abi.ConvertType(out[5], new(uint32)).(*uint32)
	out6 := *abi.ConvertType(out[6], new(uint32)).(*uint32)

	return out0, out1, out2, out3, out4, out5, out6, err

}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientSession) GetStats() ([32]byte, [32]byte, [32]byte, uint32, uint32, uint32, uint32, error) {
	return _FunctionsLoadTestClient.Contract.GetStats(&_FunctionsLoadTestClient.CallOpts)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientCallerSession) GetStats() ([32]byte, [32]byte, [32]byte, uint32, uint32, uint32, uint32, error) {
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

func (_FunctionsLoadTestClient *FunctionsLoadTestClientTransactor) ResetStats(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsLoadTestClient.contract.Transact(opts, "resetStats")
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientSession) ResetStats() (*types.Transaction, error) {
	return _FunctionsLoadTestClient.Contract.ResetStats(&_FunctionsLoadTestClient.TransactOpts)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientTransactorSession) ResetStats() (*types.Transaction, error) {
	return _FunctionsLoadTestClient.Contract.ResetStats(&_FunctionsLoadTestClient.TransactOpts)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientTransactor) SendRequest(opts *bind.TransactOpts, times uint32, source string, encryptedSecretsReferences []byte, args []string, subscriptionId uint64, jobId [32]byte) (*types.Transaction, error) {
	return _FunctionsLoadTestClient.contract.Transact(opts, "sendRequest", times, source, encryptedSecretsReferences, args, subscriptionId, jobId)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientSession) SendRequest(times uint32, source string, encryptedSecretsReferences []byte, args []string, subscriptionId uint64, jobId [32]byte) (*types.Transaction, error) {
	return _FunctionsLoadTestClient.Contract.SendRequest(&_FunctionsLoadTestClient.TransactOpts, times, source, encryptedSecretsReferences, args, subscriptionId, jobId)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientTransactorSession) SendRequest(times uint32, source string, encryptedSecretsReferences []byte, args []string, subscriptionId uint64, jobId [32]byte) (*types.Transaction, error) {
	return _FunctionsLoadTestClient.Contract.SendRequest(&_FunctionsLoadTestClient.TransactOpts, times, source, encryptedSecretsReferences, args, subscriptionId, jobId)
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

	GetStats(opts *bind.CallOpts) ([32]byte, [32]byte, [32]byte, uint32, uint32, uint32, uint32, error)

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

	ResetStats(opts *bind.TransactOpts) (*types.Transaction, error)

	SendRequest(opts *bind.TransactOpts, times uint32, source string, encryptedSecretsReferences []byte, args []string, subscriptionId uint64, jobId [32]byte) (*types.Transaction, error)

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
