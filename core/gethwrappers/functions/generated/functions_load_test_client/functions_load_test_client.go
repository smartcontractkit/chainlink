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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"EmptyArgs\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptySecrets\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptySource\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoInlineSecrets\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyRouterCanFulfill\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"RequestFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"RequestSent\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"MAX_CALLBACK_GAS\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStats\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"response\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"handleOracleFulfillment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastError\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastRequestID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastResponse\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"resetStats\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"times\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"cborEncodedRequest\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"donId\",\"type\":\"bytes32\"}],\"name\":\"sendEncodedRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"times\",\"type\":\"uint32\"},{\"internalType\":\"string\",\"name\":\"source\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"encryptedSecretsReferences\",\"type\":\"bytes\"},{\"internalType\":\"string[]\",\"name\":\"args\",\"type\":\"string[]\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"donId\",\"type\":\"bytes32\"}],\"name\":\"sendRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"times\",\"type\":\"uint32\"},{\"internalType\":\"string\",\"name\":\"source\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"slotId\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"slotVersion\",\"type\":\"uint64\"},{\"internalType\":\"string[]\",\"name\":\"args\",\"type\":\"string[]\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"donId\",\"type\":\"bytes32\"}],\"name\":\"sendRequestWithDONHostedSecrets\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalEmptyResponses\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalFailedResponses\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalRequests\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSucceededResponses\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b50604051620023ff380380620023ff833981016040819052620000349162000180565b6001600160a01b0381166080523380600081620000985760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000cb57620000cb81620000d5565b50505050620001b2565b336001600160a01b038216036200012f5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016200008f565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000602082840312156200019357600080fd5b81516001600160a01b0381168114620001ab57600080fd5b9392505050565b60805161222a620001d5600039600081816102ac0152610fa0015261222a6000f3fe608060405234801561001057600080fd5b506004361061011b5760003560e01c806379ba5097116100b2578063954491c111610081578063c59d484711610066578063c59d484714610246578063c9429e2a14610261578063f2fde38b1461028157600080fd5b8063954491c114610220578063b2518e0e1461023357600080fd5b806379ba5097146101cd578063887efe94146101d55780638aea61dc146101e85780638da5cb5b146101f857600080fd5b80635c1d92e9116100ee5780635c1d92e91461019b57806362747e42146101b35780636d9809a0146101bb578063724ec8a2146101c557600080fd5b80630ca761751461012057806329f0de3f146101355780632ab424da1461015357806347c0318614610184575b600080fd5b61013361012e3660046118b4565b610294565b005b61013d61033e565b60405161014a9190611985565b60405180910390f35b60055461016f9068010000000000000000900463ffffffff1681565b60405163ffffffff909116815260200161014a565b61018d60025481565b60405190815260200161014a565b60055461016f90640100000000900463ffffffff1681565b61013d6103cc565b61016f6203d09081565b6101336103d9565b61013361044b565b6101336101e3366004611a57565b61054d565b60055461016f9063ffffffff1681565b60005460405173ffffffffffffffffffffffffffffffffffffffff909116815260200161014a565b61013361022e366004611b1b565b6106bf565b610133610241366004611be1565b6107d5565b61024e610854565b60405161014a9796959493929190611c47565b60055461016f906c01000000000000000000000000900463ffffffff1681565b61013361028f366004611ca4565b6109dc565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610303576040517fc6829f8300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61030e8383836109f0565b60405183907f85e1543bf2f84fe80c6badbce3648c8539ad1df4d2b3d822938ca0538be727e690600090a2505050565b6004805461034b90611cda565b80601f016020809104026020016040519081016040528092919081815260200182805461037790611cda565b80156103c45780601f10610399576101008083540402835291602001916103c4565b820191906000526020600020905b8154815290600101906020018083116103a757829003601f168201915b505050505081565b6003805461034b90611cda565b6103e1610afb565b6000600281905560408051602081019091529081526003906104039082611d7b565b506040805160208101909152600081526004906104209082611d7b565b50600580547fffffffffffffffffffffffffffffffff00000000000000000000000000000000169055565b60015473ffffffffffffffffffffffffffffffffffffffff1633146104d1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610555610afb565b6105966040805160e0810190915280600081526020016000815260200160008152602001606081526020016060815260200160608152602001606081525090565b6105d889898080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508593925050610b7e9050565b85156106205761062087878080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508593925050610b8f9050565b831561063a5761063a6106338587611e95565b8290610bd9565b60005b8a63ffffffff168110156106b25761066261065783610c1c565b856203d09086610f9b565b600255600580546001919060009061068190849063ffffffff16611f5c565b92506101000a81548163ffffffff021916908363ffffffff16021790555080806106aa90611f80565b91505061063d565b5050505050505050505050565b6106c7610afb565b6107086040805160e0810190915280600081526020016000815260200160008152602001606081526020016060815260200160608152602001606081525090565b61074a89898080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508593925050610b7e9050565b61075581888861107a565b8315610768576107686106338587611e95565b60005b8a63ffffffff168110156106b25761078561065783610c1c565b60025560058054600191906000906107a490849063ffffffff16611f5c565b92506101000a81548163ffffffff021916908363ffffffff16021790555080806107cd90611f80565b91505061076b565b6107dd610afb565b60005b8463ffffffff1681101561084d576107fd84846203d09085610f9b565b600255600580546001919060009061081c90849063ffffffff16611f5c565b92506101000a81548163ffffffff021916908363ffffffff160217905550808061084590611f80565b9150506107e0565b5050505050565b6000606080600080600080610867610afb565b60025460055460038054909160049163ffffffff808316926801000000000000000081048216926c0100000000000000000000000082048316926401000000009092049091169086906108b990611cda565b80601f01602080910402602001604051908101604052809291908181526020018280546108e590611cda565b80156109325780601f1061090757610100808354040283529160200191610932565b820191906000526020600020905b81548152906001019060200180831161091557829003601f168201915b5050505050955084805461094590611cda565b80601f016020809104026020016040519081016040528092919081815260200182805461097190611cda565b80156109be5780601f10610993576101008083540402835291602001916109be565b820191906000526020600020905b8154815290600101906020018083116109a157829003601f168201915b50505050509450965096509650965096509650965090919293949596565b6109e4610afb565b6109ed8161113d565b50565b60028390556003610a018382611d7b565b506004610a0e8282611d7b565b508151600003610a59576001600560048282829054906101000a900463ffffffff16610a3a9190611f5c565b92506101000a81548163ffffffff021916908363ffffffff1602179055505b805115610aa15760016005600c8282829054906101000a900463ffffffff16610a829190611f5c565b92506101000a81548163ffffffff021916908363ffffffff1602179055505b815115801590610ab057508051155b15610af6576001600560088282829054906101000a900463ffffffff16610ad79190611f5c565b92506101000a81548163ffffffff021916908363ffffffff1602179055505b505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610b7c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016104c8565b565b610b8b8260008084611232565b5050565b8051600003610bca576040517fe889636f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60016020830152608090910152565b8051600003610c14576040517ffe936cb700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60a090910152565b60606000610c2b6101006112c9565b9050610c756040518060400160405280600c81526020017f636f64654c6f636174696f6e0000000000000000000000000000000000000000815250826112ea90919063ffffffff16565b8251610c93906002811115610c8c57610c8c611fb8565b8290611303565b60408051808201909152600881527f6c616e67756167650000000000000000000000000000000000000000000000006020820152610cd29082906112ea565b6040830151610ce9908015610c8c57610c8c611fb8565b60408051808201909152600681527f736f7572636500000000000000000000000000000000000000000000000000006020820152610d289082906112ea565b6060830151610d389082906112ea565b60a08301515115610de55760408051808201909152600481527f61726773000000000000000000000000000000000000000000000000000000006020820152610d829082906112ea565b610d8b8161133c565b60005b8360a0015151811015610ddb57610dcb8460a001518281518110610db457610db4611fe7565b6020026020010151836112ea90919063ffffffff16565b610dd481611f80565b9050610d8e565b50610de581611360565b60808301515115610ee657600083602001516002811115610e0857610e08611fb8565b03610e3f576040517fa80d31f700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60408051808201909152600f81527f736563726574734c6f636174696f6e00000000000000000000000000000000006020820152610e7e9082906112ea565b610e9783602001516002811115610c8c57610c8c611fb8565b60408051808201909152600781527f73656372657473000000000000000000000000000000000000000000000000006020820152610ed69082906112ea565b6080830151610ee690829061137e565b60c08301515115610f935760408051808201909152600981527f62797465734172677300000000000000000000000000000000000000000000006020820152610f309082906112ea565b610f398161133c565b60005b8360c0015151811015610f8957610f798460c001518281518110610f6257610f62611fe7565b60200260200101518361137e90919063ffffffff16565b610f8281611f80565b9050610f3c565b50610f9381611360565b515192915050565b6000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663461d27628688600188886040518663ffffffff1660e01b8152600401611000959493929190612016565b6020604051808303816000875af115801561101f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906110439190612060565b60405190915081907f1131472297a800fee664d1d89cfa8f7676ff07189ecc53f80bbb5f4969099db890600090a295945050505050565b60006110876101006112c9565b90506110d16040518060400160405280600681526020017f736c6f7449440000000000000000000000000000000000000000000000000000815250826112ea90919063ffffffff16565b6110de8160ff851661138b565b60408051808201909152600781527f76657273696f6e00000000000000000000000000000000000000000000000000602082015261111d9082906112ea565b611127818361138b565b6002602085015251516080909301929092525050565b3373ffffffffffffffffffffffffffffffffffffffff8216036111bc576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016104c8565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b805160000361126d576040517f22ce3edd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8383600281111561128057611280611fb8565b9081600281111561129357611293611fb8565b905250604084018280156112a9576112a9611fb8565b908180156112b9576112b9611fb8565b9052506060909301929092525050565b6112d161176b565b80516112dd9083611397565b5060006020820152919050565b6112f78260038351611411565b8151610af69082611538565b81516113109060c2611560565b50610b8b828260405160200161132891815260200190565b60405160208183030381529060405261137e565b6113478160046115c9565b60018160200181815161135a9190612079565b90525050565b61136b8160076115c9565b60018160200181815161135a919061208c565b6112f78260028351611411565b610b8b82600083611411565b6040805180820190915260608152600060208201526113b760208361209f565b156113df576113c760208361209f565b6113d290602061208c565b6113dc9083612079565b91505b60208084018390526040518085526000815290818401018181101561140357600080fd5b604052508290505b92915050565b60178167ffffffffffffffff161161143e5782516114389060e0600585901b168317611560565b50505050565b60ff8167ffffffffffffffff1611611480578251611467906018611fe0600586901b1617611560565b5082516114389067ffffffffffffffff831660016115e0565b61ffff8167ffffffffffffffff16116114c35782516114aa906019611fe0600586901b1617611560565b5082516114389067ffffffffffffffff831660026115e0565b63ffffffff8167ffffffffffffffff16116115085782516114ef90601a611fe0600586901b1617611560565b5082516114389067ffffffffffffffff831660046115e0565b825161151f90601b611fe0600586901b1617611560565b5082516114389067ffffffffffffffff831660086115e0565b60408051808201909152606081526000602082015261155983838451611665565b9392505050565b6040805180820190915260608152600060208201528251516000611585826001612079565b9050846020015182106115a6576115a6856115a18360026120da565b611754565b84516020838201018581535080518211156115bf578181525b5093949350505050565b8151610af690601f611fe0600585901b1617611560565b60408051808201909152606081526000602082015283515160006116048285612079565b9050856020015181111561162157611621866115a18360026120da565b6000600161163186610100612211565b61163b919061208c565b90508651828101878319825116178152508051831115611659578281525b50959695505050505050565b604080518082019091526060815260006020820152825182111561168857600080fd5b83515160006116978483612079565b905085602001518111156116b4576116b4866115a18360026120da565b8551805183820160200191600091808511156116ce578482525b505050602086015b6020861061170e57805182526116ed602083612079565b91506116fa602082612079565b905061170760208761208c565b95506116d6565b5181517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60208890036101000a0190811690199190911617905250849150509392505050565b81516117608383611397565b506114388382611538565b6040518060400160405280611793604051806040016040528060608152602001600081525090565b8152602001600081525090565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715611816576118166117a0565b604052919050565b600067ffffffffffffffff831115611838576118386117a0565b61186960207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f860116016117cf565b905082815283838301111561187d57600080fd5b828260208301376000602084830101529392505050565b600082601f8301126118a557600080fd5b6115598383356020850161181e565b6000806000606084860312156118c957600080fd5b83359250602084013567ffffffffffffffff808211156118e857600080fd5b6118f487838801611894565b9350604086013591508082111561190a57600080fd5b5061191786828701611894565b9150509250925092565b6000815180845260005b818110156119475760208185018101518683018201520161192b565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b6020815260006115596020830184611921565b803563ffffffff811681146119ac57600080fd5b919050565b60008083601f8401126119c357600080fd5b50813567ffffffffffffffff8111156119db57600080fd5b6020830191508360208285010111156119f357600080fd5b9250929050565b60008083601f840112611a0c57600080fd5b50813567ffffffffffffffff811115611a2457600080fd5b6020830191508360208260051b85010111156119f357600080fd5b803567ffffffffffffffff811681146119ac57600080fd5b600080600080600080600080600060c08a8c031215611a7557600080fd5b611a7e8a611998565b985060208a013567ffffffffffffffff80821115611a9b57600080fd5b611aa78d838e016119b1565b909a50985060408c0135915080821115611ac057600080fd5b611acc8d838e016119b1565b909850965060608c0135915080821115611ae557600080fd5b50611af28c828d016119fa565b9095509350611b05905060808b01611a3f565b915060a08a013590509295985092959850929598565b600080600080600080600080600060e08a8c031215611b3957600080fd5b611b428a611998565b985060208a013567ffffffffffffffff80821115611b5f57600080fd5b611b6b8d838e016119b1565b909a50985060408c0135915060ff82168214611b8657600080fd5b819750611b9560608d01611a3f565b965060808c0135915080821115611bab57600080fd5b50611bb88c828d016119fa565b9095509350611bcb905060a08b01611a3f565b915060c08a013590509295985092959850929598565b60008060008060808587031215611bf757600080fd5b611c0085611998565b9350602085013567ffffffffffffffff811115611c1c57600080fd5b611c2887828801611894565b935050611c3760408601611a3f565b9396929550929360600135925050565b87815260e060208201526000611c6060e0830189611921565b8281036040840152611c728189611921565b63ffffffff97881660608501529587166080840152505091841660a083015290921660c0909201919091529392505050565b600060208284031215611cb657600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461155957600080fd5b600181811c90821680611cee57607f821691505b602082108103611d27577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f821115610af657600081815260208120601f850160051c81016020861015611d545750805b601f850160051c820191505b81811015611d7357828155600101611d60565b505050505050565b815167ffffffffffffffff811115611d9557611d956117a0565b611da981611da38454611cda565b84611d2d565b602080601f831160018114611dfc5760008415611dc65750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555611d73565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015611e4957888601518255948401946001909101908401611e2a565b5085821015611e8557878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b600067ffffffffffffffff80841115611eb057611eb06117a0565b8360051b6020611ec18183016117cf565b868152918501918181019036841115611ed957600080fd5b865b84811015611f2157803586811115611ef35760008081fd5b880136601f820112611f055760008081fd5b611f1336823587840161181e565b845250918301918301611edb565b50979650505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b63ffffffff818116838216019080821115611f7957611f79611f2d565b5092915050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203611fb157611fb1611f2d565b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b67ffffffffffffffff8616815260a06020820152600061203960a0830187611921565b61ffff9590951660408301525063ffffffff92909216606083015260809091015292915050565b60006020828403121561207257600080fd5b5051919050565b8082018082111561140b5761140b611f2d565b8181038181111561140b5761140b611f2d565b6000826120d5577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500690565b808202811582820484141761140b5761140b611f2d565b600181815b8085111561214a57817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0482111561213057612130611f2d565b8085161561213d57918102915b93841c93908002906120f6565b509250929050565b6000826121615750600161140b565b8161216e5750600061140b565b8160018114612184576002811461218e576121aa565b600191505061140b565b60ff84111561219f5761219f611f2d565b50506001821b61140b565b5060208310610133831016604e8410600b84101617156121cd575081810a61140b565b6121d783836120f1565b807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0482111561220957612209611f2d565b029392505050565b6000611559838361215256fea164736f6c6343000813000a",
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
	return address, tx, &FunctionsLoadTestClient{address: address, abi: *parsed, FunctionsLoadTestClientCaller: FunctionsLoadTestClientCaller{contract: contract}, FunctionsLoadTestClientTransactor: FunctionsLoadTestClientTransactor{contract: contract}, FunctionsLoadTestClientFilterer: FunctionsLoadTestClientFilterer{contract: contract}}, nil
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

func (_FunctionsLoadTestClient *FunctionsLoadTestClientCaller) GetStats(opts *bind.CallOpts) ([32]byte, []byte, []byte, uint32, uint32, uint32, uint32, error) {
	var out []interface{}
	err := _FunctionsLoadTestClient.contract.Call(opts, &out, "getStats")

	if err != nil {
		return *new([32]byte), *new([]byte), *new([]byte), *new(uint32), *new(uint32), *new(uint32), *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)
	out2 := *abi.ConvertType(out[2], new([]byte)).(*[]byte)
	out3 := *abi.ConvertType(out[3], new(uint32)).(*uint32)
	out4 := *abi.ConvertType(out[4], new(uint32)).(*uint32)
	out5 := *abi.ConvertType(out[5], new(uint32)).(*uint32)
	out6 := *abi.ConvertType(out[6], new(uint32)).(*uint32)

	return out0, out1, out2, out3, out4, out5, out6, err

}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientSession) GetStats() ([32]byte, []byte, []byte, uint32, uint32, uint32, uint32, error) {
	return _FunctionsLoadTestClient.Contract.GetStats(&_FunctionsLoadTestClient.CallOpts)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientCallerSession) GetStats() ([32]byte, []byte, []byte, uint32, uint32, uint32, uint32, error) {
	return _FunctionsLoadTestClient.Contract.GetStats(&_FunctionsLoadTestClient.CallOpts)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientCaller) LastError(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _FunctionsLoadTestClient.contract.Call(opts, &out, "lastError")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientSession) LastError() ([]byte, error) {
	return _FunctionsLoadTestClient.Contract.LastError(&_FunctionsLoadTestClient.CallOpts)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientCallerSession) LastError() ([]byte, error) {
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

func (_FunctionsLoadTestClient *FunctionsLoadTestClientCaller) LastResponse(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _FunctionsLoadTestClient.contract.Call(opts, &out, "lastResponse")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientSession) LastResponse() ([]byte, error) {
	return _FunctionsLoadTestClient.Contract.LastResponse(&_FunctionsLoadTestClient.CallOpts)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientCallerSession) LastResponse() ([]byte, error) {
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

func (_FunctionsLoadTestClient *FunctionsLoadTestClientTransactor) SendEncodedRequest(opts *bind.TransactOpts, times uint32, cborEncodedRequest []byte, subscriptionId uint64, donId [32]byte) (*types.Transaction, error) {
	return _FunctionsLoadTestClient.contract.Transact(opts, "sendEncodedRequest", times, cborEncodedRequest, subscriptionId, donId)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientSession) SendEncodedRequest(times uint32, cborEncodedRequest []byte, subscriptionId uint64, donId [32]byte) (*types.Transaction, error) {
	return _FunctionsLoadTestClient.Contract.SendEncodedRequest(&_FunctionsLoadTestClient.TransactOpts, times, cborEncodedRequest, subscriptionId, donId)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientTransactorSession) SendEncodedRequest(times uint32, cborEncodedRequest []byte, subscriptionId uint64, donId [32]byte) (*types.Transaction, error) {
	return _FunctionsLoadTestClient.Contract.SendEncodedRequest(&_FunctionsLoadTestClient.TransactOpts, times, cborEncodedRequest, subscriptionId, donId)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientTransactor) SendRequest(opts *bind.TransactOpts, times uint32, source string, encryptedSecretsReferences []byte, args []string, subscriptionId uint64, donId [32]byte) (*types.Transaction, error) {
	return _FunctionsLoadTestClient.contract.Transact(opts, "sendRequest", times, source, encryptedSecretsReferences, args, subscriptionId, donId)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientSession) SendRequest(times uint32, source string, encryptedSecretsReferences []byte, args []string, subscriptionId uint64, donId [32]byte) (*types.Transaction, error) {
	return _FunctionsLoadTestClient.Contract.SendRequest(&_FunctionsLoadTestClient.TransactOpts, times, source, encryptedSecretsReferences, args, subscriptionId, donId)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientTransactorSession) SendRequest(times uint32, source string, encryptedSecretsReferences []byte, args []string, subscriptionId uint64, donId [32]byte) (*types.Transaction, error) {
	return _FunctionsLoadTestClient.Contract.SendRequest(&_FunctionsLoadTestClient.TransactOpts, times, source, encryptedSecretsReferences, args, subscriptionId, donId)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientTransactor) SendRequestWithDONHostedSecrets(opts *bind.TransactOpts, times uint32, source string, slotId uint8, slotVersion uint64, args []string, subscriptionId uint64, donId [32]byte) (*types.Transaction, error) {
	return _FunctionsLoadTestClient.contract.Transact(opts, "sendRequestWithDONHostedSecrets", times, source, slotId, slotVersion, args, subscriptionId, donId)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientSession) SendRequestWithDONHostedSecrets(times uint32, source string, slotId uint8, slotVersion uint64, args []string, subscriptionId uint64, donId [32]byte) (*types.Transaction, error) {
	return _FunctionsLoadTestClient.Contract.SendRequestWithDONHostedSecrets(&_FunctionsLoadTestClient.TransactOpts, times, source, slotId, slotVersion, args, subscriptionId, donId)
}

func (_FunctionsLoadTestClient *FunctionsLoadTestClientTransactorSession) SendRequestWithDONHostedSecrets(times uint32, source string, slotId uint8, slotVersion uint64, args []string, subscriptionId uint64, donId [32]byte) (*types.Transaction, error) {
	return _FunctionsLoadTestClient.Contract.SendRequestWithDONHostedSecrets(&_FunctionsLoadTestClient.TransactOpts, times, source, slotId, slotVersion, args, subscriptionId, donId)
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

	GetStats(opts *bind.CallOpts) ([32]byte, []byte, []byte, uint32, uint32, uint32, uint32, error)

	LastError(opts *bind.CallOpts) ([]byte, error)

	LastRequestID(opts *bind.CallOpts) ([32]byte, error)

	LastResponse(opts *bind.CallOpts) ([]byte, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TotalEmptyResponses(opts *bind.CallOpts) (uint32, error)

	TotalFailedResponses(opts *bind.CallOpts) (uint32, error)

	TotalRequests(opts *bind.CallOpts) (uint32, error)

	TotalSucceededResponses(opts *bind.CallOpts) (uint32, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	HandleOracleFulfillment(opts *bind.TransactOpts, requestId [32]byte, response []byte, err []byte) (*types.Transaction, error)

	ResetStats(opts *bind.TransactOpts) (*types.Transaction, error)

	SendEncodedRequest(opts *bind.TransactOpts, times uint32, cborEncodedRequest []byte, subscriptionId uint64, donId [32]byte) (*types.Transaction, error)

	SendRequest(opts *bind.TransactOpts, times uint32, source string, encryptedSecretsReferences []byte, args []string, subscriptionId uint64, donId [32]byte) (*types.Transaction, error)

	SendRequestWithDONHostedSecrets(opts *bind.TransactOpts, times uint32, source string, slotId uint8, slotVersion uint64, args []string, subscriptionId uint64, donId [32]byte) (*types.Transaction, error)

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
