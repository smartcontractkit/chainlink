// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ocr2dr_client_example

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
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated"
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

type FunctionsRequest struct {
	CodeLocation    uint8
	SecretsLocation uint8
	Language        uint8
	Source          string
	Secrets         []byte
	Args            []string
}

var OCR2DRClientExampleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"EmptyArgs\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptySecrets\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptySource\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RequestIsAlreadyPending\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RequestIsNotPending\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SenderIsNotRegistry\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"UnexpectedRequestID\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"RequestFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"RequestSent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"source\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"secrets\",\"type\":\"bytes\"},{\"internalType\":\"string[]\",\"name\":\"args\",\"type\":\"string[]\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"}],\"name\":\"SendRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"enumFunctions.Location\",\"name\":\"codeLocation\",\"type\":\"uint8\"},{\"internalType\":\"enumFunctions.Location\",\"name\":\"secretsLocation\",\"type\":\"uint8\"},{\"internalType\":\"enumFunctions.CodeLanguage\",\"name\":\"language\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"source\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"secrets\",\"type\":\"bytes\"},{\"internalType\":\"string[]\",\"name\":\"args\",\"type\":\"string[]\"}],\"internalType\":\"structFunctions.Request\",\"name\":\"req\",\"type\":\"tuple\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"}],\"name\":\"estimateCost\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDONPublicKey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"response\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"handleOracleFulfillment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastError\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastRequestId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastResponse\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162001ed438038062001ed4833981016040819052620000349162000199565b600080546001600160a01b0319166001600160a01b038316178155339081906001600160a01b038216620000af5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600280546001600160a01b0319166001600160a01b0384811691909117909155811615620000e257620000e281620000ec565b50505050620001cb565b6001600160a01b038116331415620001475760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401620000a6565b600380546001600160a01b0319166001600160a01b03838116918217909255600254604051919216907fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127890600090a350565b600060208284031215620001ac57600080fd5b81516001600160a01b0381168114620001c457600080fd5b9392505050565b611cf980620001db6000396000f3fe608060405234801561001057600080fd5b50600436106100be5760003560e01c8063d328a91e11610076578063d769717e1161005b578063d769717e14610166578063f2fde38b14610179578063fc2a88c31461018c57600080fd5b8063d328a91e1461012e578063d4b391751461013657600080fd5b806362747e42116100a757806362747e42146100f657806379ba5097146100fe5780638da5cb5b1461010657600080fd5b80630ca76175146100c357806329f0de3f146100d8575b600080fd5b6100d66100d1366004611533565b6101a3565b005b6100e061026e565b6040516100ed9190611873565b60405180910390f35b6100e06102fc565b6100d6610309565b60025460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100ed565b6100e061040f565b6101496101443660046116e9565b6104d8565b6040516bffffffffffffffffffffffff90911681526020016100ed565b6100d661017436600461160e565b61057b565b6100d66101873660046114e0565b61067c565b61019560045481565b6040519081526020016100ed565b600083815260016020526040902054839073ffffffffffffffffffffffffffffffffffffffff163314610202576040517fa0c5ec6300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008181526001602052604080822080547fffffffffffffffffffffffff00000000000000000000000000000000000000001690555182917f85e1543bf2f84fe80c6badbce3648c8539ad1df4d2b3d822938ca0538be727e691a2610268848484610690565b50505050565b6006805461027b90611b4c565b80601f01602080910402602001604051908101604052809291908181526020018280546102a790611b4c565b80156102f45780601f106102c9576101008083540402835291602001916102f4565b820191906000526020600020905b8154815290600101906020018083116102d757829003601f168201915b505050505081565b6005805461027b90611b4c565b60035473ffffffffffffffffffffffffffffffffffffffff16331461038f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b600280547fffffffffffffffffffffffff00000000000000000000000000000000000000008082163390811790935560038054909116905560405173ffffffffffffffffffffffffffffffffffffffff909116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a350565b60008054604080517fd328a91e000000000000000000000000000000000000000000000000000000008152905160609373ffffffffffffffffffffffffffffffffffffffff9093169263d328a91e9260048082019391829003018186803b15801561047957600080fd5b505afa15801561048d573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526104d391908101906115a0565b905090565b6000805473ffffffffffffffffffffffffffffffffffffffff1663d227d24585610501886106f5565b86866040518563ffffffff1660e01b81526004016105229493929190611886565b60206040518083038186803b15801561053a57600080fd5b505afa15801561054e573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061057291906117fb565b95945050505050565b61058361098b565b6105bd6040805160c08101909152806000815260200160008152602001600081526020016060815260200160608152602001606081525090565b6105ff88888080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508593925050610a0e9050565b84156106475761064786868080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508593925050610a1f9050565b82156106615761066161065a8486611b13565b8290610a66565b61066f8183619c403a610aa6565b6004555050505050505050565b61068461098b565b61068d81610c67565b50565b82600454146106ce576040517fd068bf5b00000000000000000000000000000000000000000000000000000000815260048101849052602401610386565b81516106e19060059060208501906112ad565b5080516102689060069060208401906112ad565b6060610714604051806040016040528060608152602001600081525090565b61072081610100610d5e565b5060408051808201909152600c81527f636f64654c6f636174696f6e00000000000000000000000000000000000000006020820152610760908290610dc9565b825161077e90600181111561077757610777611c3d565b8290610de5565b60408051808201909152600881527f6c616e677561676500000000000000000000000000000000000000000000000060208201526107bd908290610dc9565b60408301516107d490801561077757610777611c3d565b60408051808201909152600681527f736f7572636500000000000000000000000000000000000000000000000000006020820152610813908290610dc9565b6060830151610823908290610dc9565b60a083015151156108d25760408051808201909152600481527f6172677300000000000000000000000000000000000000000000000000000000602082015261086d908290610dc9565b61087681610e0b565b60005b8360a00151518110156108c8576108b68460a00151828151811061089f5761089f611c6c565b602002602001015183610dc990919063ffffffff16565b806108c081611b9a565b915050610879565b506108d281610e16565b608083015151156109845760408051808201909152600f81527f736563726574734c6f636174696f6e0000000000000000000000000000000000602082015261091c908290610dc9565b6109358360200151600181111561077757610777611c3d565b60408051808201909152600781527f73656372657473000000000000000000000000000000000000000000000000006020820152610974908290610dc9565b6080830151610984908290610e21565b5192915050565b60025473ffffffffffffffffffffffffffffffffffffffff163314610a0c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610386565b565b610a1b8260008084610e2e565b5050565b8051610a57576040517fe889636f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006020830152608090910152565b8051610a9e576040517ffe936cb700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60a090910152565b60008054819073ffffffffffffffffffffffffffffffffffffffff1663a98eb20686610ad1896106f5565b87876040518563ffffffff1660e01b8152600401610af29493929190611886565b602060405180830381600087803b158015610b0c57600080fd5b505af1158015610b20573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b44919061151a565b905060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16635ab1bd536040518163ffffffff1660e01b815260040160206040518083038186803b158015610bac57600080fd5b505afa158015610bc0573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610be491906114fd565b60008281526001602052604080822080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff9490941693909317909255905182917f1131472297a800fee664d1d89cfa8f7676ff07189ecc53f80bbb5f4969099db891a295945050505050565b73ffffffffffffffffffffffffffffffffffffffff8116331415610ce7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610386565b600380547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff838116918217909255600254604051919216907fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127890600090a350565b604080518082019091526060815260006020820152610d7e602083611bd3565b15610da657610d8e602083611bd3565b610d99906020611afc565b610da39083611980565b91505b506020808301829052604080518085526000815283019091019052815b92915050565b610dd68260038351610ec2565b610de08282610fd1565b505050565b67ffffffffffffffff811115610dff57610a1b8282610fff565b610a1b82600083610ec2565b61068d816004611036565b61068d816007611036565b610dd68260028351610ec2565b8051610e66576040517f22ce3edd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b83836001811115610e7957610e79611c3d565b90816001811115610e8c57610e8c611c3d565b90525060408401828015610ea257610ea2611c3d565b90818015610eb257610eb2611c3d565b9052506060909301929092525050565b60178167ffffffffffffffff1611610ee7576102688360e0600585901b168317611047565b60ff8167ffffffffffffffff1611610f2557610f0e836018611fe0600586901b1617611047565b506102688367ffffffffffffffff8316600161106c565b61ffff8167ffffffffffffffff1611610f6457610f4d836019611fe0600586901b1617611047565b506102688367ffffffffffffffff8316600261106c565b63ffffffff8167ffffffffffffffff1611610fa557610f8e83601a611fe0600586901b1617611047565b506102688367ffffffffffffffff8316600461106c565b610fba83601b611fe0600586901b1617611047565b506102688367ffffffffffffffff8316600861106c565b604080518082019091526060815260006020820152610ff88384600001515184855161109a565b9392505050565b61100a8260c2611047565b50610a1b828260405160200161102291815260200190565b604051602081830303815290604052610e21565b610de082601f611fe0600585901b16175b604080518082019091526060815260006020820152610ff883846000015151846111a2565b6040805180820190915260608152600060208201526110928485600001515185856111fe565b949350505050565b60408051808201909152606081526000602082015282518211156110bd57600080fd5b60208501516110cc8386611980565b11156110ff576110ff856110ef876020015187866110ea9190611980565b61127f565b6110fa906002611abf565b611296565b60008086518051876020830101935080888701111561111e5787860182525b505050602084015b6020841061115e578051825261113d602083611980565b915061114a602082611980565b9050611157602085611afc565b9350611126565b5181517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60208690036101000a019081169019919091161790525083949350505050565b604080518082019091526060815260006020820152836020015183106111d7576111d784856020015160026110fa9190611abf565b8351805160208583010184815350808514156111f4576001810182525b5093949350505050565b60408051808201909152606081526000602082015260208501516112228584611980565b111561123657611236856110ef8685611980565b60006001611246846101006119f9565b6112509190611afc565b90508551838682010185831982511617815250805184870111156112745783860181525b509495945050505050565b600081831115611290575081610dc3565b50919050565b81516112a28383610d5e565b506102688382610fd1565b8280546112b990611b4c565b90600052602060002090601f0160209004810192826112db5760008555611321565b82601f106112f457805160ff1916838001178555611321565b82800160010185558215611321579182015b82811115611321578251825591602001919060010190611306565b5061132d929150611331565b5090565b5b8082111561132d5760008155600101611332565b600067ffffffffffffffff8084111561136157611361611c9b565b8360051b60206113728183016118eb565b86815293508084018583810189101561138a57600080fd5b60009350835b888110156113c5578135868111156113a6578586fd5b6113b28b828b0161143b565b8452509183019190830190600101611390565b5050505050509392505050565b600082601f8301126113e357600080fd5b610ff883833560208501611346565b60008083601f84011261140457600080fd5b50813567ffffffffffffffff81111561141c57600080fd5b60208301915083602082850101111561143457600080fd5b9250929050565b600082601f83011261144c57600080fd5b813561145f61145a8261193a565b6118eb565b81815284602083860101111561147457600080fd5b816020850160208301376000918101602001919091529392505050565b8035600181106114a057600080fd5b919050565b8035600281106114a057600080fd5b803563ffffffff811681146114a057600080fd5b803567ffffffffffffffff811681146114a057600080fd5b6000602082840312156114f257600080fd5b8135610ff881611cca565b60006020828403121561150f57600080fd5b8151610ff881611cca565b60006020828403121561152c57600080fd5b5051919050565b60008060006060848603121561154857600080fd5b83359250602084013567ffffffffffffffff8082111561156757600080fd5b6115738783880161143b565b9350604086013591508082111561158957600080fd5b506115968682870161143b565b9150509250925092565b6000602082840312156115b257600080fd5b815167ffffffffffffffff8111156115c957600080fd5b8201601f810184136115da57600080fd5b80516115e861145a8261193a565b8181528560208385010111156115fd57600080fd5b610572826020830160208601611b20565b60008060008060008060006080888a03121561162957600080fd5b873567ffffffffffffffff8082111561164157600080fd5b61164d8b838c016113f2565b909950975060208a013591508082111561166657600080fd5b6116728b838c016113f2565b909750955060408a013591508082111561168b57600080fd5b818a0191508a601f83011261169f57600080fd5b8135818111156116ae57600080fd5b8b60208260051b85010111156116c357600080fd5b6020830195508094505050506116db606089016114c8565b905092959891949750929550565b600080600080608085870312156116ff57600080fd5b843567ffffffffffffffff8082111561171757600080fd5b9086019060c0828903121561172b57600080fd5b6117336118c2565b61173c836114a5565b815261174a602084016114a5565b602082015261175b60408401611491565b604082015260608301358281111561177257600080fd5b61177e8a82860161143b565b60608301525060808301358281111561179657600080fd5b6117a28a82860161143b565b60808301525060a0830135828111156117ba57600080fd5b6117c68a8286016113d2565b60a08301525095506117dd915050602086016114c8565b92506117eb604086016114b4565b9396929550929360600135925050565b60006020828403121561180d57600080fd5b81516bffffffffffffffffffffffff81168114610ff857600080fd5b60008151808452611841816020860160208601611b20565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000610ff86020830184611829565b67ffffffffffffffff851681526080602082015260006118a96080830186611829565b63ffffffff949094166040830152506060015292915050565b60405160c0810167ffffffffffffffff811182821017156118e5576118e5611c9b565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561193257611932611c9b565b604052919050565b600067ffffffffffffffff82111561195457611954611c9b565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b6000821982111561199357611993611c0e565b500190565b600181815b808511156119f157817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048211156119d7576119d7611c0e565b808516156119e457918102915b93841c939080029061199d565b509250929050565b6000610ff88383600082611a0f57506001610dc3565b81611a1c57506000610dc3565b8160018114611a325760028114611a3c57611a58565b6001915050610dc3565b60ff841115611a4d57611a4d611c0e565b50506001821b610dc3565b5060208310610133831016604e8410600b8410161715611a7b575081810a610dc3565b611a858383611998565b807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04821115611ab757611ab7611c0e565b029392505050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615611af757611af7611c0e565b500290565b600082821015611b0e57611b0e611c0e565b500390565b6000610ff8368484611346565b60005b83811015611b3b578181015183820152602001611b23565b838111156102685750506000910152565b600181811c90821680611b6057607f821691505b60208210811415611290577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415611bcc57611bcc611c0e565b5060010190565b600082611c09577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500690565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b73ffffffffffffffffffffffffffffffffffffffff8116811461068d57600080fdfea164736f6c6343000806000a",
}

var OCR2DRClientExampleABI = OCR2DRClientExampleMetaData.ABI

var OCR2DRClientExampleBin = OCR2DRClientExampleMetaData.Bin

func DeployOCR2DRClientExample(auth *bind.TransactOpts, backend bind.ContractBackend, oracle common.Address) (common.Address, *types.Transaction, *OCR2DRClientExample, error) {
	parsed, err := OCR2DRClientExampleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OCR2DRClientExampleBin), backend, oracle)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OCR2DRClientExample{OCR2DRClientExampleCaller: OCR2DRClientExampleCaller{contract: contract}, OCR2DRClientExampleTransactor: OCR2DRClientExampleTransactor{contract: contract}, OCR2DRClientExampleFilterer: OCR2DRClientExampleFilterer{contract: contract}}, nil
}

type OCR2DRClientExample struct {
	address common.Address
	abi     abi.ABI
	OCR2DRClientExampleCaller
	OCR2DRClientExampleTransactor
	OCR2DRClientExampleFilterer
}

type OCR2DRClientExampleCaller struct {
	contract *bind.BoundContract
}

type OCR2DRClientExampleTransactor struct {
	contract *bind.BoundContract
}

type OCR2DRClientExampleFilterer struct {
	contract *bind.BoundContract
}

type OCR2DRClientExampleSession struct {
	Contract     *OCR2DRClientExample
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type OCR2DRClientExampleCallerSession struct {
	Contract *OCR2DRClientExampleCaller
	CallOpts bind.CallOpts
}

type OCR2DRClientExampleTransactorSession struct {
	Contract     *OCR2DRClientExampleTransactor
	TransactOpts bind.TransactOpts
}

type OCR2DRClientExampleRaw struct {
	Contract *OCR2DRClientExample
}

type OCR2DRClientExampleCallerRaw struct {
	Contract *OCR2DRClientExampleCaller
}

type OCR2DRClientExampleTransactorRaw struct {
	Contract *OCR2DRClientExampleTransactor
}

func NewOCR2DRClientExample(address common.Address, backend bind.ContractBackend) (*OCR2DRClientExample, error) {
	abi, err := abi.JSON(strings.NewReader(OCR2DRClientExampleABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindOCR2DRClientExample(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OCR2DRClientExample{address: address, abi: abi, OCR2DRClientExampleCaller: OCR2DRClientExampleCaller{contract: contract}, OCR2DRClientExampleTransactor: OCR2DRClientExampleTransactor{contract: contract}, OCR2DRClientExampleFilterer: OCR2DRClientExampleFilterer{contract: contract}}, nil
}

func NewOCR2DRClientExampleCaller(address common.Address, caller bind.ContractCaller) (*OCR2DRClientExampleCaller, error) {
	contract, err := bindOCR2DRClientExample(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OCR2DRClientExampleCaller{contract: contract}, nil
}

func NewOCR2DRClientExampleTransactor(address common.Address, transactor bind.ContractTransactor) (*OCR2DRClientExampleTransactor, error) {
	contract, err := bindOCR2DRClientExample(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OCR2DRClientExampleTransactor{contract: contract}, nil
}

func NewOCR2DRClientExampleFilterer(address common.Address, filterer bind.ContractFilterer) (*OCR2DRClientExampleFilterer, error) {
	contract, err := bindOCR2DRClientExample(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OCR2DRClientExampleFilterer{contract: contract}, nil
}

func bindOCR2DRClientExample(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OCR2DRClientExampleABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_OCR2DRClientExample *OCR2DRClientExampleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR2DRClientExample.Contract.OCR2DRClientExampleCaller.contract.Call(opts, result, method, params...)
}

func (_OCR2DRClientExample *OCR2DRClientExampleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2DRClientExample.Contract.OCR2DRClientExampleTransactor.contract.Transfer(opts)
}

func (_OCR2DRClientExample *OCR2DRClientExampleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR2DRClientExample.Contract.OCR2DRClientExampleTransactor.contract.Transact(opts, method, params...)
}

func (_OCR2DRClientExample *OCR2DRClientExampleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR2DRClientExample.Contract.contract.Call(opts, result, method, params...)
}

func (_OCR2DRClientExample *OCR2DRClientExampleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2DRClientExample.Contract.contract.Transfer(opts)
}

func (_OCR2DRClientExample *OCR2DRClientExampleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR2DRClientExample.Contract.contract.Transact(opts, method, params...)
}

func (_OCR2DRClientExample *OCR2DRClientExampleCaller) EstimateCost(opts *bind.CallOpts, req FunctionsRequest, subscriptionId uint64, gasLimit uint32, gasPrice *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _OCR2DRClientExample.contract.Call(opts, &out, "estimateCost", req, subscriptionId, gasLimit, gasPrice)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_OCR2DRClientExample *OCR2DRClientExampleSession) EstimateCost(req FunctionsRequest, subscriptionId uint64, gasLimit uint32, gasPrice *big.Int) (*big.Int, error) {
	return _OCR2DRClientExample.Contract.EstimateCost(&_OCR2DRClientExample.CallOpts, req, subscriptionId, gasLimit, gasPrice)
}

func (_OCR2DRClientExample *OCR2DRClientExampleCallerSession) EstimateCost(req FunctionsRequest, subscriptionId uint64, gasLimit uint32, gasPrice *big.Int) (*big.Int, error) {
	return _OCR2DRClientExample.Contract.EstimateCost(&_OCR2DRClientExample.CallOpts, req, subscriptionId, gasLimit, gasPrice)
}

func (_OCR2DRClientExample *OCR2DRClientExampleCaller) GetDONPublicKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _OCR2DRClientExample.contract.Call(opts, &out, "getDONPublicKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_OCR2DRClientExample *OCR2DRClientExampleSession) GetDONPublicKey() ([]byte, error) {
	return _OCR2DRClientExample.Contract.GetDONPublicKey(&_OCR2DRClientExample.CallOpts)
}

func (_OCR2DRClientExample *OCR2DRClientExampleCallerSession) GetDONPublicKey() ([]byte, error) {
	return _OCR2DRClientExample.Contract.GetDONPublicKey(&_OCR2DRClientExample.CallOpts)
}

func (_OCR2DRClientExample *OCR2DRClientExampleCaller) LastError(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _OCR2DRClientExample.contract.Call(opts, &out, "lastError")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_OCR2DRClientExample *OCR2DRClientExampleSession) LastError() ([]byte, error) {
	return _OCR2DRClientExample.Contract.LastError(&_OCR2DRClientExample.CallOpts)
}

func (_OCR2DRClientExample *OCR2DRClientExampleCallerSession) LastError() ([]byte, error) {
	return _OCR2DRClientExample.Contract.LastError(&_OCR2DRClientExample.CallOpts)
}

func (_OCR2DRClientExample *OCR2DRClientExampleCaller) LastRequestId(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _OCR2DRClientExample.contract.Call(opts, &out, "lastRequestId")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_OCR2DRClientExample *OCR2DRClientExampleSession) LastRequestId() ([32]byte, error) {
	return _OCR2DRClientExample.Contract.LastRequestId(&_OCR2DRClientExample.CallOpts)
}

func (_OCR2DRClientExample *OCR2DRClientExampleCallerSession) LastRequestId() ([32]byte, error) {
	return _OCR2DRClientExample.Contract.LastRequestId(&_OCR2DRClientExample.CallOpts)
}

func (_OCR2DRClientExample *OCR2DRClientExampleCaller) LastResponse(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _OCR2DRClientExample.contract.Call(opts, &out, "lastResponse")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_OCR2DRClientExample *OCR2DRClientExampleSession) LastResponse() ([]byte, error) {
	return _OCR2DRClientExample.Contract.LastResponse(&_OCR2DRClientExample.CallOpts)
}

func (_OCR2DRClientExample *OCR2DRClientExampleCallerSession) LastResponse() ([]byte, error) {
	return _OCR2DRClientExample.Contract.LastResponse(&_OCR2DRClientExample.CallOpts)
}

func (_OCR2DRClientExample *OCR2DRClientExampleCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OCR2DRClientExample.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_OCR2DRClientExample *OCR2DRClientExampleSession) Owner() (common.Address, error) {
	return _OCR2DRClientExample.Contract.Owner(&_OCR2DRClientExample.CallOpts)
}

func (_OCR2DRClientExample *OCR2DRClientExampleCallerSession) Owner() (common.Address, error) {
	return _OCR2DRClientExample.Contract.Owner(&_OCR2DRClientExample.CallOpts)
}

func (_OCR2DRClientExample *OCR2DRClientExampleTransactor) SendRequest(opts *bind.TransactOpts, source string, secrets []byte, args []string, subscriptionId uint64) (*types.Transaction, error) {
	return _OCR2DRClientExample.contract.Transact(opts, "SendRequest", source, secrets, args, subscriptionId)
}

func (_OCR2DRClientExample *OCR2DRClientExampleSession) SendRequest(source string, secrets []byte, args []string, subscriptionId uint64) (*types.Transaction, error) {
	return _OCR2DRClientExample.Contract.SendRequest(&_OCR2DRClientExample.TransactOpts, source, secrets, args, subscriptionId)
}

func (_OCR2DRClientExample *OCR2DRClientExampleTransactorSession) SendRequest(source string, secrets []byte, args []string, subscriptionId uint64) (*types.Transaction, error) {
	return _OCR2DRClientExample.Contract.SendRequest(&_OCR2DRClientExample.TransactOpts, source, secrets, args, subscriptionId)
}

func (_OCR2DRClientExample *OCR2DRClientExampleTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2DRClientExample.contract.Transact(opts, "acceptOwnership")
}

func (_OCR2DRClientExample *OCR2DRClientExampleSession) AcceptOwnership() (*types.Transaction, error) {
	return _OCR2DRClientExample.Contract.AcceptOwnership(&_OCR2DRClientExample.TransactOpts)
}

func (_OCR2DRClientExample *OCR2DRClientExampleTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _OCR2DRClientExample.Contract.AcceptOwnership(&_OCR2DRClientExample.TransactOpts)
}

func (_OCR2DRClientExample *OCR2DRClientExampleTransactor) HandleOracleFulfillment(opts *bind.TransactOpts, requestId [32]byte, response []byte, err []byte) (*types.Transaction, error) {
	return _OCR2DRClientExample.contract.Transact(opts, "handleOracleFulfillment", requestId, response, err)
}

func (_OCR2DRClientExample *OCR2DRClientExampleSession) HandleOracleFulfillment(requestId [32]byte, response []byte, err []byte) (*types.Transaction, error) {
	return _OCR2DRClientExample.Contract.HandleOracleFulfillment(&_OCR2DRClientExample.TransactOpts, requestId, response, err)
}

func (_OCR2DRClientExample *OCR2DRClientExampleTransactorSession) HandleOracleFulfillment(requestId [32]byte, response []byte, err []byte) (*types.Transaction, error) {
	return _OCR2DRClientExample.Contract.HandleOracleFulfillment(&_OCR2DRClientExample.TransactOpts, requestId, response, err)
}

func (_OCR2DRClientExample *OCR2DRClientExampleTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _OCR2DRClientExample.contract.Transact(opts, "transferOwnership", to)
}

func (_OCR2DRClientExample *OCR2DRClientExampleSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OCR2DRClientExample.Contract.TransferOwnership(&_OCR2DRClientExample.TransactOpts, to)
}

func (_OCR2DRClientExample *OCR2DRClientExampleTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OCR2DRClientExample.Contract.TransferOwnership(&_OCR2DRClientExample.TransactOpts, to)
}

type OCR2DRClientExampleOwnershipTransferRequestedIterator struct {
	Event *OCR2DRClientExampleOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2DRClientExampleOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2DRClientExampleOwnershipTransferRequested)
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
		it.Event = new(OCR2DRClientExampleOwnershipTransferRequested)
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

func (it *OCR2DRClientExampleOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *OCR2DRClientExampleOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2DRClientExampleOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_OCR2DRClientExample *OCR2DRClientExampleFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OCR2DRClientExampleOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR2DRClientExample.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OCR2DRClientExampleOwnershipTransferRequestedIterator{contract: _OCR2DRClientExample.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_OCR2DRClientExample *OCR2DRClientExampleFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OCR2DRClientExampleOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR2DRClientExample.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2DRClientExampleOwnershipTransferRequested)
				if err := _OCR2DRClientExample.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_OCR2DRClientExample *OCR2DRClientExampleFilterer) ParseOwnershipTransferRequested(log types.Log) (*OCR2DRClientExampleOwnershipTransferRequested, error) {
	event := new(OCR2DRClientExampleOwnershipTransferRequested)
	if err := _OCR2DRClientExample.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2DRClientExampleOwnershipTransferredIterator struct {
	Event *OCR2DRClientExampleOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2DRClientExampleOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2DRClientExampleOwnershipTransferred)
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
		it.Event = new(OCR2DRClientExampleOwnershipTransferred)
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

func (it *OCR2DRClientExampleOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *OCR2DRClientExampleOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2DRClientExampleOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_OCR2DRClientExample *OCR2DRClientExampleFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OCR2DRClientExampleOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR2DRClientExample.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OCR2DRClientExampleOwnershipTransferredIterator{contract: _OCR2DRClientExample.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_OCR2DRClientExample *OCR2DRClientExampleFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OCR2DRClientExampleOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR2DRClientExample.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2DRClientExampleOwnershipTransferred)
				if err := _OCR2DRClientExample.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_OCR2DRClientExample *OCR2DRClientExampleFilterer) ParseOwnershipTransferred(log types.Log) (*OCR2DRClientExampleOwnershipTransferred, error) {
	event := new(OCR2DRClientExampleOwnershipTransferred)
	if err := _OCR2DRClientExample.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2DRClientExampleRequestFulfilledIterator struct {
	Event *OCR2DRClientExampleRequestFulfilled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2DRClientExampleRequestFulfilledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2DRClientExampleRequestFulfilled)
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
		it.Event = new(OCR2DRClientExampleRequestFulfilled)
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

func (it *OCR2DRClientExampleRequestFulfilledIterator) Error() error {
	return it.fail
}

func (it *OCR2DRClientExampleRequestFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2DRClientExampleRequestFulfilled struct {
	Id  [32]byte
	Raw types.Log
}

func (_OCR2DRClientExample *OCR2DRClientExampleFilterer) FilterRequestFulfilled(opts *bind.FilterOpts, id [][32]byte) (*OCR2DRClientExampleRequestFulfilledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _OCR2DRClientExample.contract.FilterLogs(opts, "RequestFulfilled", idRule)
	if err != nil {
		return nil, err
	}
	return &OCR2DRClientExampleRequestFulfilledIterator{contract: _OCR2DRClientExample.contract, event: "RequestFulfilled", logs: logs, sub: sub}, nil
}

func (_OCR2DRClientExample *OCR2DRClientExampleFilterer) WatchRequestFulfilled(opts *bind.WatchOpts, sink chan<- *OCR2DRClientExampleRequestFulfilled, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _OCR2DRClientExample.contract.WatchLogs(opts, "RequestFulfilled", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2DRClientExampleRequestFulfilled)
				if err := _OCR2DRClientExample.contract.UnpackLog(event, "RequestFulfilled", log); err != nil {
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

func (_OCR2DRClientExample *OCR2DRClientExampleFilterer) ParseRequestFulfilled(log types.Log) (*OCR2DRClientExampleRequestFulfilled, error) {
	event := new(OCR2DRClientExampleRequestFulfilled)
	if err := _OCR2DRClientExample.contract.UnpackLog(event, "RequestFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2DRClientExampleRequestSentIterator struct {
	Event *OCR2DRClientExampleRequestSent

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2DRClientExampleRequestSentIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2DRClientExampleRequestSent)
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
		it.Event = new(OCR2DRClientExampleRequestSent)
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

func (it *OCR2DRClientExampleRequestSentIterator) Error() error {
	return it.fail
}

func (it *OCR2DRClientExampleRequestSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2DRClientExampleRequestSent struct {
	Id  [32]byte
	Raw types.Log
}

func (_OCR2DRClientExample *OCR2DRClientExampleFilterer) FilterRequestSent(opts *bind.FilterOpts, id [][32]byte) (*OCR2DRClientExampleRequestSentIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _OCR2DRClientExample.contract.FilterLogs(opts, "RequestSent", idRule)
	if err != nil {
		return nil, err
	}
	return &OCR2DRClientExampleRequestSentIterator{contract: _OCR2DRClientExample.contract, event: "RequestSent", logs: logs, sub: sub}, nil
}

func (_OCR2DRClientExample *OCR2DRClientExampleFilterer) WatchRequestSent(opts *bind.WatchOpts, sink chan<- *OCR2DRClientExampleRequestSent, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _OCR2DRClientExample.contract.WatchLogs(opts, "RequestSent", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2DRClientExampleRequestSent)
				if err := _OCR2DRClientExample.contract.UnpackLog(event, "RequestSent", log); err != nil {
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

func (_OCR2DRClientExample *OCR2DRClientExampleFilterer) ParseRequestSent(log types.Log) (*OCR2DRClientExampleRequestSent, error) {
	event := new(OCR2DRClientExampleRequestSent)
	if err := _OCR2DRClientExample.contract.UnpackLog(event, "RequestSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_OCR2DRClientExample *OCR2DRClientExample) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _OCR2DRClientExample.abi.Events["OwnershipTransferRequested"].ID:
		return _OCR2DRClientExample.ParseOwnershipTransferRequested(log)
	case _OCR2DRClientExample.abi.Events["OwnershipTransferred"].ID:
		return _OCR2DRClientExample.ParseOwnershipTransferred(log)
	case _OCR2DRClientExample.abi.Events["RequestFulfilled"].ID:
		return _OCR2DRClientExample.ParseRequestFulfilled(log)
	case _OCR2DRClientExample.abi.Events["RequestSent"].ID:
		return _OCR2DRClientExample.ParseRequestSent(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (OCR2DRClientExampleOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (OCR2DRClientExampleOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (OCR2DRClientExampleRequestFulfilled) Topic() common.Hash {
	return common.HexToHash("0x85e1543bf2f84fe80c6badbce3648c8539ad1df4d2b3d822938ca0538be727e6")
}

func (OCR2DRClientExampleRequestSent) Topic() common.Hash {
	return common.HexToHash("0x1131472297a800fee664d1d89cfa8f7676ff07189ecc53f80bbb5f4969099db8")
}

func (_OCR2DRClientExample *OCR2DRClientExample) Address() common.Address {
	return _OCR2DRClientExample.address
}

type OCR2DRClientExampleInterface interface {
	EstimateCost(opts *bind.CallOpts, req FunctionsRequest, subscriptionId uint64, gasLimit uint32, gasPrice *big.Int) (*big.Int, error)

	GetDONPublicKey(opts *bind.CallOpts) ([]byte, error)

	LastError(opts *bind.CallOpts) ([]byte, error)

	LastRequestId(opts *bind.CallOpts) ([32]byte, error)

	LastResponse(opts *bind.CallOpts) ([]byte, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SendRequest(opts *bind.TransactOpts, source string, secrets []byte, args []string, subscriptionId uint64) (*types.Transaction, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	HandleOracleFulfillment(opts *bind.TransactOpts, requestId [32]byte, response []byte, err []byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OCR2DRClientExampleOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OCR2DRClientExampleOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*OCR2DRClientExampleOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OCR2DRClientExampleOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OCR2DRClientExampleOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*OCR2DRClientExampleOwnershipTransferred, error)

	FilterRequestFulfilled(opts *bind.FilterOpts, id [][32]byte) (*OCR2DRClientExampleRequestFulfilledIterator, error)

	WatchRequestFulfilled(opts *bind.WatchOpts, sink chan<- *OCR2DRClientExampleRequestFulfilled, id [][32]byte) (event.Subscription, error)

	ParseRequestFulfilled(log types.Log) (*OCR2DRClientExampleRequestFulfilled, error)

	FilterRequestSent(opts *bind.FilterOpts, id [][32]byte) (*OCR2DRClientExampleRequestSentIterator, error)

	WatchRequestSent(opts *bind.WatchOpts, sink chan<- *OCR2DRClientExampleRequestSent, id [][32]byte) (event.Subscription, error)

	ParseRequestSent(log types.Log) (*OCR2DRClientExampleRequestSent, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
