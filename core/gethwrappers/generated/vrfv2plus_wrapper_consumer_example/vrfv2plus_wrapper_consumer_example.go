// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrfv2plus_wrapper_consumer_example

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

var VRFV2PlusWrapperConsumerExampleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vrfV2Wrapper\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlyVRFWrapperCanFulfill\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"payment\",\"type\":\"uint256\"}],\"name\":\"WrappedRequestFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"paid\",\"type\":\"uint256\"}],\"name\":\"WrapperRequestMade\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkToken\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_requestId\",\"type\":\"uint256\"}],\"name\":\"getRequestStatus\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"paid\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"fulfilled\",\"type\":\"bool\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_vrfV2PlusWrapper\",\"outputs\":[{\"internalType\":\"contractIVRFV2PlusWrapper\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"_requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"_numWords\",\"type\":\"uint32\"}],\"name\":\"makeRequest\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"_requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"_numWords\",\"type\":\"uint32\"}],\"name\":\"makeRequestNative\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"_randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_requests\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"paid\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"fulfilled\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"native\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawLink\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawNative\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60c06040523480156200001157600080fd5b5060405162001662380380620016628339810160408190526200003491620001f5565b33806000836000819050806001600160a01b0316631c4695f46040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200007d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620000a39190620001f5565b6001600160a01b0390811660805290811660a052831690506200010d5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620001405762000140816200014a565b5050505062000227565b336001600160a01b03821603620001a45760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000104565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000602082840312156200020857600080fd5b81516001600160a01b03811681146200022057600080fd5b9392505050565b60805160a0516113dd62000285600039600081816101aa0152818161047b01528181610aa801528181610b5301528181610bfd01528181610ce30152610d52015260008181610242015281816106220152610b1701526113dd6000f3fe608060405234801561001057600080fd5b50600436106100df5760003560e01c806384276d811161008c578063a168fa8911610066578063a168fa89146101cc578063d8a4676f1461021e578063e76d516814610240578063f2fde38b1461026657600080fd5b806384276d81146101535780638da5cb5b146101665780639ed0868d146101a557600080fd5b80631fe543e3116100bd5780631fe543e31461012357806379ba5097146101385780637a8042bd1461014057600080fd5b80630c09b832146100e457806312065fe01461010a5780631e1a349914610110575b600080fd5b6100f76100f2366004611078565b610279565b6040519081526020015b60405180910390f35b476100f7565b6100f761011e366004611078565b6103b4565b6101366101313660046110f3565b610479565b005b61013661051b565b61013661014e3660046111db565b610618565b6101366101613660046111db565b610715565b60005473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610101565b6101807f000000000000000000000000000000000000000000000000000000000000000081565b6102016101da3660046111db565b600260205260009081526040902080546001820154600390920154909160ff908116911683565b604080519384529115156020840152151590820152606001610101565b61023161022c3660046111db565b6107e7565b6040516101019392919061122f565b7f0000000000000000000000000000000000000000000000000000000000000000610180565b610136610274366004611259565b610907565b600061028361091b565b600061029f60405180602001604052806000151581525061099e565b905060006102af86868685610a5a565b604080516080810182528281526000602080830182815284518381528083018652848601908152606085018490528784526002808452959093208451815590516001820180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055915180519699509496509194909361033c938501920190610fff565b5060609190910151600390910180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001691151591909117905560405181815283907f5f56b4c20db9f5b294cbf6f681368de4a992a27e2de2ee702dcf2cbbfa791ec49060200160405180910390a250509392505050565b60006103be61091b565b60006103da60405180602001604052806001151581525061099e565b905060006103ea86868685610c95565b60408051608081018252828152600060208083018281528451838152808301865284860190815260016060860181905288855260028085529690942085518155915193820180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001694151594909417909355915180519699509496509194909361033c938501920190610fff565b7f00000000000000000000000000000000000000000000000000000000000000003373ffffffffffffffffffffffffffffffffffffffff82161461050c576040517f8ba9316e00000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff821660248201526044015b60405180910390fd5b6105168383610df3565b505050565b60015473ffffffffffffffffffffffffffffffffffffffff16331461059c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610503565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b61062061091b565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb61067b60005473ffffffffffffffffffffffffffffffffffffffff1690565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e084901b16815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602481018490526044016020604051808303816000875af11580156106ed573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107119190611296565b5050565b61071d61091b565b6000805460405173ffffffffffffffffffffffffffffffffffffffff9091169083908381818185875af1925050503d8060008114610777576040519150601f19603f3d011682016040523d82523d6000602084013e61077c565b606091505b5050905080610711576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f77697468647261774e6174697665206661696c656400000000000000000000006044820152606401610503565b6000818152600260205260408120548190606090610861576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f72657175657374206e6f7420666f756e640000000000000000000000000000006044820152606401610503565b6000848152600260208181526040808420815160808101835281548152600182015460ff16151581850152938101805483518186028101860185528181529294938601938301828280156108d457602002820191906000526020600020905b8154815260200190600101908083116108c0575b50505091835250506003919091015460ff1615156020918201528151908201516040909201519097919650945092505050565b61090f61091b565b61091881610f0a565b50565b60005473ffffffffffffffffffffffffffffffffffffffff16331461099c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610503565b565b60607f92fd13387c7fe7befbc38d303d6468778fb9731bc4583f17d92989c6fcfdeaaa826040516024016109d791511515815260200190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009093169290921790915292915050565b6040517f4306d35400000000000000000000000000000000000000000000000000000000815263ffffffff85166004820152600090819073ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001690634306d35490602401602060405180830381865afa158015610aef573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b1391906112b8565b90507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16634000aea07f00000000000000000000000000000000000000000000000000000000000000008389898989604051602001610b8a9493929190611335565b6040516020818303038152906040526040518463ffffffff1660e01b8152600401610bb793929190611372565b6020604051808303816000875af1158015610bd6573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610bfa9190611296565b507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663fc2a88c36040518163ffffffff1660e01b8152600401602060405180830381865afa158015610c66573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c8a91906112b8565b915094509492505050565b6040517f4b16093500000000000000000000000000000000000000000000000000000000815263ffffffff85166004820152600090819073ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001690634b16093590602401602060405180830381865afa158015610d2a573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d4e91906112b8565b90507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16639cfc058e82888888886040518663ffffffff1660e01b8152600401610db09493929190611335565b60206040518083038185885af1158015610dce573d6000803e3d6000fd5b50505050506040513d601f19601f82011682018060405250810190610c8a91906112b8565b600082815260026020526040902054610e68576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f72657175657374206e6f7420666f756e640000000000000000000000000000006044820152606401610503565b6000828152600260208181526040909220600181810180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690911790558351610ebb93919092019190840190610fff565b50600082815260026020526040908190205490517f6c84e12b4c188e61f1b4727024a5cf05c025fa58467e5eedf763c0744c89da7b91610efe91859185916113a7565b60405180910390a15050565b3373ffffffffffffffffffffffffffffffffffffffff821603610f89576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610503565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b82805482825590600052602060002090810192821561103a579160200282015b8281111561103a57825182559160200191906001019061101f565b5061104692915061104a565b5090565b5b80821115611046576000815560010161104b565b803563ffffffff8116811461107357600080fd5b919050565b60008060006060848603121561108d57600080fd5b6110968461105f565b9250602084013561ffff811681146110ad57600080fd5b91506110bb6040850161105f565b90509250925092565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6000806040838503121561110657600080fd5b8235915060208084013567ffffffffffffffff8082111561112657600080fd5b818601915086601f83011261113a57600080fd5b81358181111561114c5761114c6110c4565b8060051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f8301168101818110858211171561118f5761118f6110c4565b6040529182528482019250838101850191898311156111ad57600080fd5b938501935b828510156111cb578435845293850193928501926111b2565b8096505050505050509250929050565b6000602082840312156111ed57600080fd5b5035919050565b600081518084526020808501945080840160005b8381101561122457815187529582019590820190600101611208565b509495945050505050565b838152821515602082015260606040820152600061125060608301846111f4565b95945050505050565b60006020828403121561126b57600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461128f57600080fd5b9392505050565b6000602082840312156112a857600080fd5b8151801515811461128f57600080fd5b6000602082840312156112ca57600080fd5b5051919050565b6000815180845260005b818110156112f7576020818501810151868301820152016112db565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b600063ffffffff808716835261ffff861660208401528085166040840152506080606083015261136860808301846112d1565b9695505050505050565b73ffffffffffffffffffffffffffffffffffffffff8416815282602082015260606040820152600061125060608301846112d1565b8381526060602082015260006113c060608301856111f4565b905082604083015294935050505056fea164736f6c6343000813000a",
}

var VRFV2PlusWrapperConsumerExampleABI = VRFV2PlusWrapperConsumerExampleMetaData.ABI

var VRFV2PlusWrapperConsumerExampleBin = VRFV2PlusWrapperConsumerExampleMetaData.Bin

func DeployVRFV2PlusWrapperConsumerExample(auth *bind.TransactOpts, backend bind.ContractBackend, _vrfV2Wrapper common.Address) (common.Address, *types.Transaction, *VRFV2PlusWrapperConsumerExample, error) {
	parsed, err := VRFV2PlusWrapperConsumerExampleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFV2PlusWrapperConsumerExampleBin), backend, _vrfV2Wrapper)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFV2PlusWrapperConsumerExample{address: address, abi: *parsed, VRFV2PlusWrapperConsumerExampleCaller: VRFV2PlusWrapperConsumerExampleCaller{contract: contract}, VRFV2PlusWrapperConsumerExampleTransactor: VRFV2PlusWrapperConsumerExampleTransactor{contract: contract}, VRFV2PlusWrapperConsumerExampleFilterer: VRFV2PlusWrapperConsumerExampleFilterer{contract: contract}}, nil
}

type VRFV2PlusWrapperConsumerExample struct {
	address common.Address
	abi     abi.ABI
	VRFV2PlusWrapperConsumerExampleCaller
	VRFV2PlusWrapperConsumerExampleTransactor
	VRFV2PlusWrapperConsumerExampleFilterer
}

type VRFV2PlusWrapperConsumerExampleCaller struct {
	contract *bind.BoundContract
}

type VRFV2PlusWrapperConsumerExampleTransactor struct {
	contract *bind.BoundContract
}

type VRFV2PlusWrapperConsumerExampleFilterer struct {
	contract *bind.BoundContract
}

type VRFV2PlusWrapperConsumerExampleSession struct {
	Contract     *VRFV2PlusWrapperConsumerExample
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFV2PlusWrapperConsumerExampleCallerSession struct {
	Contract *VRFV2PlusWrapperConsumerExampleCaller
	CallOpts bind.CallOpts
}

type VRFV2PlusWrapperConsumerExampleTransactorSession struct {
	Contract     *VRFV2PlusWrapperConsumerExampleTransactor
	TransactOpts bind.TransactOpts
}

type VRFV2PlusWrapperConsumerExampleRaw struct {
	Contract *VRFV2PlusWrapperConsumerExample
}

type VRFV2PlusWrapperConsumerExampleCallerRaw struct {
	Contract *VRFV2PlusWrapperConsumerExampleCaller
}

type VRFV2PlusWrapperConsumerExampleTransactorRaw struct {
	Contract *VRFV2PlusWrapperConsumerExampleTransactor
}

func NewVRFV2PlusWrapperConsumerExample(address common.Address, backend bind.ContractBackend) (*VRFV2PlusWrapperConsumerExample, error) {
	abi, err := abi.JSON(strings.NewReader(VRFV2PlusWrapperConsumerExampleABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFV2PlusWrapperConsumerExample(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusWrapperConsumerExample{address: address, abi: abi, VRFV2PlusWrapperConsumerExampleCaller: VRFV2PlusWrapperConsumerExampleCaller{contract: contract}, VRFV2PlusWrapperConsumerExampleTransactor: VRFV2PlusWrapperConsumerExampleTransactor{contract: contract}, VRFV2PlusWrapperConsumerExampleFilterer: VRFV2PlusWrapperConsumerExampleFilterer{contract: contract}}, nil
}

func NewVRFV2PlusWrapperConsumerExampleCaller(address common.Address, caller bind.ContractCaller) (*VRFV2PlusWrapperConsumerExampleCaller, error) {
	contract, err := bindVRFV2PlusWrapperConsumerExample(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusWrapperConsumerExampleCaller{contract: contract}, nil
}

func NewVRFV2PlusWrapperConsumerExampleTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFV2PlusWrapperConsumerExampleTransactor, error) {
	contract, err := bindVRFV2PlusWrapperConsumerExample(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusWrapperConsumerExampleTransactor{contract: contract}, nil
}

func NewVRFV2PlusWrapperConsumerExampleFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFV2PlusWrapperConsumerExampleFilterer, error) {
	contract, err := bindVRFV2PlusWrapperConsumerExample(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusWrapperConsumerExampleFilterer{contract: contract}, nil
}

func bindVRFV2PlusWrapperConsumerExample(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFV2PlusWrapperConsumerExampleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2PlusWrapperConsumerExample.Contract.VRFV2PlusWrapperConsumerExampleCaller.contract.Call(opts, result, method, params...)
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusWrapperConsumerExample.Contract.VRFV2PlusWrapperConsumerExampleTransactor.contract.Transfer(opts)
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2PlusWrapperConsumerExample.Contract.VRFV2PlusWrapperConsumerExampleTransactor.contract.Transact(opts, method, params...)
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2PlusWrapperConsumerExample.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusWrapperConsumerExample.Contract.contract.Transfer(opts)
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2PlusWrapperConsumerExample.Contract.contract.Transact(opts, method, params...)
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleCaller) GetBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusWrapperConsumerExample.contract.Call(opts, &out, "getBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleSession) GetBalance() (*big.Int, error) {
	return _VRFV2PlusWrapperConsumerExample.Contract.GetBalance(&_VRFV2PlusWrapperConsumerExample.CallOpts)
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleCallerSession) GetBalance() (*big.Int, error) {
	return _VRFV2PlusWrapperConsumerExample.Contract.GetBalance(&_VRFV2PlusWrapperConsumerExample.CallOpts)
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleCaller) GetLinkToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2PlusWrapperConsumerExample.contract.Call(opts, &out, "getLinkToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleSession) GetLinkToken() (common.Address, error) {
	return _VRFV2PlusWrapperConsumerExample.Contract.GetLinkToken(&_VRFV2PlusWrapperConsumerExample.CallOpts)
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleCallerSession) GetLinkToken() (common.Address, error) {
	return _VRFV2PlusWrapperConsumerExample.Contract.GetLinkToken(&_VRFV2PlusWrapperConsumerExample.CallOpts)
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleCaller) GetRequestStatus(opts *bind.CallOpts, _requestId *big.Int) (GetRequestStatus,

	error) {
	var out []interface{}
	err := _VRFV2PlusWrapperConsumerExample.contract.Call(opts, &out, "getRequestStatus", _requestId)

	outstruct := new(GetRequestStatus)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Paid = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Fulfilled = *abi.ConvertType(out[1], new(bool)).(*bool)
	outstruct.RandomWords = *abi.ConvertType(out[2], new([]*big.Int)).(*[]*big.Int)

	return *outstruct, err

}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleSession) GetRequestStatus(_requestId *big.Int) (GetRequestStatus,

	error) {
	return _VRFV2PlusWrapperConsumerExample.Contract.GetRequestStatus(&_VRFV2PlusWrapperConsumerExample.CallOpts, _requestId)
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleCallerSession) GetRequestStatus(_requestId *big.Int) (GetRequestStatus,

	error) {
	return _VRFV2PlusWrapperConsumerExample.Contract.GetRequestStatus(&_VRFV2PlusWrapperConsumerExample.CallOpts, _requestId)
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleCaller) IVrfV2PlusWrapper(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2PlusWrapperConsumerExample.contract.Call(opts, &out, "i_vrfV2PlusWrapper")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleSession) IVrfV2PlusWrapper() (common.Address, error) {
	return _VRFV2PlusWrapperConsumerExample.Contract.IVrfV2PlusWrapper(&_VRFV2PlusWrapperConsumerExample.CallOpts)
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleCallerSession) IVrfV2PlusWrapper() (common.Address, error) {
	return _VRFV2PlusWrapperConsumerExample.Contract.IVrfV2PlusWrapper(&_VRFV2PlusWrapperConsumerExample.CallOpts)
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2PlusWrapperConsumerExample.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleSession) Owner() (common.Address, error) {
	return _VRFV2PlusWrapperConsumerExample.Contract.Owner(&_VRFV2PlusWrapperConsumerExample.CallOpts)
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleCallerSession) Owner() (common.Address, error) {
	return _VRFV2PlusWrapperConsumerExample.Contract.Owner(&_VRFV2PlusWrapperConsumerExample.CallOpts)
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleCaller) SRequests(opts *bind.CallOpts, arg0 *big.Int) (SRequests,

	error) {
	var out []interface{}
	err := _VRFV2PlusWrapperConsumerExample.contract.Call(opts, &out, "s_requests", arg0)

	outstruct := new(SRequests)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Paid = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Fulfilled = *abi.ConvertType(out[1], new(bool)).(*bool)
	outstruct.Native = *abi.ConvertType(out[2], new(bool)).(*bool)

	return *outstruct, err

}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleSession) SRequests(arg0 *big.Int) (SRequests,

	error) {
	return _VRFV2PlusWrapperConsumerExample.Contract.SRequests(&_VRFV2PlusWrapperConsumerExample.CallOpts, arg0)
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleCallerSession) SRequests(arg0 *big.Int) (SRequests,

	error) {
	return _VRFV2PlusWrapperConsumerExample.Contract.SRequests(&_VRFV2PlusWrapperConsumerExample.CallOpts, arg0)
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusWrapperConsumerExample.contract.Transact(opts, "acceptOwnership")
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFV2PlusWrapperConsumerExample.Contract.AcceptOwnership(&_VRFV2PlusWrapperConsumerExample.TransactOpts)
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFV2PlusWrapperConsumerExample.Contract.AcceptOwnership(&_VRFV2PlusWrapperConsumerExample.TransactOpts)
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleTransactor) MakeRequest(opts *bind.TransactOpts, _callbackGasLimit uint32, _requestConfirmations uint16, _numWords uint32) (*types.Transaction, error) {
	return _VRFV2PlusWrapperConsumerExample.contract.Transact(opts, "makeRequest", _callbackGasLimit, _requestConfirmations, _numWords)
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleSession) MakeRequest(_callbackGasLimit uint32, _requestConfirmations uint16, _numWords uint32) (*types.Transaction, error) {
	return _VRFV2PlusWrapperConsumerExample.Contract.MakeRequest(&_VRFV2PlusWrapperConsumerExample.TransactOpts, _callbackGasLimit, _requestConfirmations, _numWords)
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleTransactorSession) MakeRequest(_callbackGasLimit uint32, _requestConfirmations uint16, _numWords uint32) (*types.Transaction, error) {
	return _VRFV2PlusWrapperConsumerExample.Contract.MakeRequest(&_VRFV2PlusWrapperConsumerExample.TransactOpts, _callbackGasLimit, _requestConfirmations, _numWords)
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleTransactor) MakeRequestNative(opts *bind.TransactOpts, _callbackGasLimit uint32, _requestConfirmations uint16, _numWords uint32) (*types.Transaction, error) {
	return _VRFV2PlusWrapperConsumerExample.contract.Transact(opts, "makeRequestNative", _callbackGasLimit, _requestConfirmations, _numWords)
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleSession) MakeRequestNative(_callbackGasLimit uint32, _requestConfirmations uint16, _numWords uint32) (*types.Transaction, error) {
	return _VRFV2PlusWrapperConsumerExample.Contract.MakeRequestNative(&_VRFV2PlusWrapperConsumerExample.TransactOpts, _callbackGasLimit, _requestConfirmations, _numWords)
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleTransactorSession) MakeRequestNative(_callbackGasLimit uint32, _requestConfirmations uint16, _numWords uint32) (*types.Transaction, error) {
	return _VRFV2PlusWrapperConsumerExample.Contract.MakeRequestNative(&_VRFV2PlusWrapperConsumerExample.TransactOpts, _callbackGasLimit, _requestConfirmations, _numWords)
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, _requestId *big.Int, _randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2PlusWrapperConsumerExample.contract.Transact(opts, "rawFulfillRandomWords", _requestId, _randomWords)
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleSession) RawFulfillRandomWords(_requestId *big.Int, _randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2PlusWrapperConsumerExample.Contract.RawFulfillRandomWords(&_VRFV2PlusWrapperConsumerExample.TransactOpts, _requestId, _randomWords)
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleTransactorSession) RawFulfillRandomWords(_requestId *big.Int, _randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2PlusWrapperConsumerExample.Contract.RawFulfillRandomWords(&_VRFV2PlusWrapperConsumerExample.TransactOpts, _requestId, _randomWords)
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusWrapperConsumerExample.contract.Transact(opts, "transferOwnership", to)
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusWrapperConsumerExample.Contract.TransferOwnership(&_VRFV2PlusWrapperConsumerExample.TransactOpts, to)
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusWrapperConsumerExample.Contract.TransferOwnership(&_VRFV2PlusWrapperConsumerExample.TransactOpts, to)
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleTransactor) WithdrawLink(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusWrapperConsumerExample.contract.Transact(opts, "withdrawLink", amount)
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleSession) WithdrawLink(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusWrapperConsumerExample.Contract.WithdrawLink(&_VRFV2PlusWrapperConsumerExample.TransactOpts, amount)
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleTransactorSession) WithdrawLink(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusWrapperConsumerExample.Contract.WithdrawLink(&_VRFV2PlusWrapperConsumerExample.TransactOpts, amount)
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleTransactor) WithdrawNative(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusWrapperConsumerExample.contract.Transact(opts, "withdrawNative", amount)
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleSession) WithdrawNative(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusWrapperConsumerExample.Contract.WithdrawNative(&_VRFV2PlusWrapperConsumerExample.TransactOpts, amount)
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleTransactorSession) WithdrawNative(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusWrapperConsumerExample.Contract.WithdrawNative(&_VRFV2PlusWrapperConsumerExample.TransactOpts, amount)
}

type VRFV2PlusWrapperConsumerExampleOwnershipTransferRequestedIterator struct {
	Event *VRFV2PlusWrapperConsumerExampleOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2PlusWrapperConsumerExampleOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2PlusWrapperConsumerExampleOwnershipTransferRequested)
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
		it.Event = new(VRFV2PlusWrapperConsumerExampleOwnershipTransferRequested)
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

func (it *VRFV2PlusWrapperConsumerExampleOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFV2PlusWrapperConsumerExampleOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2PlusWrapperConsumerExampleOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusWrapperConsumerExampleOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusWrapperConsumerExample.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusWrapperConsumerExampleOwnershipTransferRequestedIterator{contract: _VRFV2PlusWrapperConsumerExample.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFV2PlusWrapperConsumerExampleOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusWrapperConsumerExample.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2PlusWrapperConsumerExampleOwnershipTransferRequested)
				if err := _VRFV2PlusWrapperConsumerExample.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleFilterer) ParseOwnershipTransferRequested(log types.Log) (*VRFV2PlusWrapperConsumerExampleOwnershipTransferRequested, error) {
	event := new(VRFV2PlusWrapperConsumerExampleOwnershipTransferRequested)
	if err := _VRFV2PlusWrapperConsumerExample.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFV2PlusWrapperConsumerExampleOwnershipTransferredIterator struct {
	Event *VRFV2PlusWrapperConsumerExampleOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2PlusWrapperConsumerExampleOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2PlusWrapperConsumerExampleOwnershipTransferred)
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
		it.Event = new(VRFV2PlusWrapperConsumerExampleOwnershipTransferred)
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

func (it *VRFV2PlusWrapperConsumerExampleOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFV2PlusWrapperConsumerExampleOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2PlusWrapperConsumerExampleOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusWrapperConsumerExampleOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusWrapperConsumerExample.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusWrapperConsumerExampleOwnershipTransferredIterator{contract: _VRFV2PlusWrapperConsumerExample.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFV2PlusWrapperConsumerExampleOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusWrapperConsumerExample.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2PlusWrapperConsumerExampleOwnershipTransferred)
				if err := _VRFV2PlusWrapperConsumerExample.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleFilterer) ParseOwnershipTransferred(log types.Log) (*VRFV2PlusWrapperConsumerExampleOwnershipTransferred, error) {
	event := new(VRFV2PlusWrapperConsumerExampleOwnershipTransferred)
	if err := _VRFV2PlusWrapperConsumerExample.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFV2PlusWrapperConsumerExampleWrappedRequestFulfilledIterator struct {
	Event *VRFV2PlusWrapperConsumerExampleWrappedRequestFulfilled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2PlusWrapperConsumerExampleWrappedRequestFulfilledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2PlusWrapperConsumerExampleWrappedRequestFulfilled)
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
		it.Event = new(VRFV2PlusWrapperConsumerExampleWrappedRequestFulfilled)
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

func (it *VRFV2PlusWrapperConsumerExampleWrappedRequestFulfilledIterator) Error() error {
	return it.fail
}

func (it *VRFV2PlusWrapperConsumerExampleWrappedRequestFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2PlusWrapperConsumerExampleWrappedRequestFulfilled struct {
	RequestId   *big.Int
	RandomWords []*big.Int
	Payment     *big.Int
	Raw         types.Log
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleFilterer) FilterWrappedRequestFulfilled(opts *bind.FilterOpts) (*VRFV2PlusWrapperConsumerExampleWrappedRequestFulfilledIterator, error) {

	logs, sub, err := _VRFV2PlusWrapperConsumerExample.contract.FilterLogs(opts, "WrappedRequestFulfilled")
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusWrapperConsumerExampleWrappedRequestFulfilledIterator{contract: _VRFV2PlusWrapperConsumerExample.contract, event: "WrappedRequestFulfilled", logs: logs, sub: sub}, nil
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleFilterer) WatchWrappedRequestFulfilled(opts *bind.WatchOpts, sink chan<- *VRFV2PlusWrapperConsumerExampleWrappedRequestFulfilled) (event.Subscription, error) {

	logs, sub, err := _VRFV2PlusWrapperConsumerExample.contract.WatchLogs(opts, "WrappedRequestFulfilled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2PlusWrapperConsumerExampleWrappedRequestFulfilled)
				if err := _VRFV2PlusWrapperConsumerExample.contract.UnpackLog(event, "WrappedRequestFulfilled", log); err != nil {
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

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleFilterer) ParseWrappedRequestFulfilled(log types.Log) (*VRFV2PlusWrapperConsumerExampleWrappedRequestFulfilled, error) {
	event := new(VRFV2PlusWrapperConsumerExampleWrappedRequestFulfilled)
	if err := _VRFV2PlusWrapperConsumerExample.contract.UnpackLog(event, "WrappedRequestFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFV2PlusWrapperConsumerExampleWrapperRequestMadeIterator struct {
	Event *VRFV2PlusWrapperConsumerExampleWrapperRequestMade

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2PlusWrapperConsumerExampleWrapperRequestMadeIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2PlusWrapperConsumerExampleWrapperRequestMade)
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
		it.Event = new(VRFV2PlusWrapperConsumerExampleWrapperRequestMade)
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

func (it *VRFV2PlusWrapperConsumerExampleWrapperRequestMadeIterator) Error() error {
	return it.fail
}

func (it *VRFV2PlusWrapperConsumerExampleWrapperRequestMadeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2PlusWrapperConsumerExampleWrapperRequestMade struct {
	RequestId *big.Int
	Paid      *big.Int
	Raw       types.Log
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleFilterer) FilterWrapperRequestMade(opts *bind.FilterOpts, requestId []*big.Int) (*VRFV2PlusWrapperConsumerExampleWrapperRequestMadeIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _VRFV2PlusWrapperConsumerExample.contract.FilterLogs(opts, "WrapperRequestMade", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusWrapperConsumerExampleWrapperRequestMadeIterator{contract: _VRFV2PlusWrapperConsumerExample.contract, event: "WrapperRequestMade", logs: logs, sub: sub}, nil
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleFilterer) WatchWrapperRequestMade(opts *bind.WatchOpts, sink chan<- *VRFV2PlusWrapperConsumerExampleWrapperRequestMade, requestId []*big.Int) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _VRFV2PlusWrapperConsumerExample.contract.WatchLogs(opts, "WrapperRequestMade", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2PlusWrapperConsumerExampleWrapperRequestMade)
				if err := _VRFV2PlusWrapperConsumerExample.contract.UnpackLog(event, "WrapperRequestMade", log); err != nil {
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

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleFilterer) ParseWrapperRequestMade(log types.Log) (*VRFV2PlusWrapperConsumerExampleWrapperRequestMade, error) {
	event := new(VRFV2PlusWrapperConsumerExampleWrapperRequestMade)
	if err := _VRFV2PlusWrapperConsumerExample.contract.UnpackLog(event, "WrapperRequestMade", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetRequestStatus struct {
	Paid        *big.Int
	Fulfilled   bool
	RandomWords []*big.Int
}
type SRequests struct {
	Paid      *big.Int
	Fulfilled bool
	Native    bool
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExample) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFV2PlusWrapperConsumerExample.abi.Events["OwnershipTransferRequested"].ID:
		return _VRFV2PlusWrapperConsumerExample.ParseOwnershipTransferRequested(log)
	case _VRFV2PlusWrapperConsumerExample.abi.Events["OwnershipTransferred"].ID:
		return _VRFV2PlusWrapperConsumerExample.ParseOwnershipTransferred(log)
	case _VRFV2PlusWrapperConsumerExample.abi.Events["WrappedRequestFulfilled"].ID:
		return _VRFV2PlusWrapperConsumerExample.ParseWrappedRequestFulfilled(log)
	case _VRFV2PlusWrapperConsumerExample.abi.Events["WrapperRequestMade"].ID:
		return _VRFV2PlusWrapperConsumerExample.ParseWrapperRequestMade(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFV2PlusWrapperConsumerExampleOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VRFV2PlusWrapperConsumerExampleOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (VRFV2PlusWrapperConsumerExampleWrappedRequestFulfilled) Topic() common.Hash {
	return common.HexToHash("0x6c84e12b4c188e61f1b4727024a5cf05c025fa58467e5eedf763c0744c89da7b")
}

func (VRFV2PlusWrapperConsumerExampleWrapperRequestMade) Topic() common.Hash {
	return common.HexToHash("0x5f56b4c20db9f5b294cbf6f681368de4a992a27e2de2ee702dcf2cbbfa791ec4")
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExample) Address() common.Address {
	return _VRFV2PlusWrapperConsumerExample.address
}

type VRFV2PlusWrapperConsumerExampleInterface interface {
	GetBalance(opts *bind.CallOpts) (*big.Int, error)

	GetLinkToken(opts *bind.CallOpts) (common.Address, error)

	GetRequestStatus(opts *bind.CallOpts, _requestId *big.Int) (GetRequestStatus,

		error)

	IVrfV2PlusWrapper(opts *bind.CallOpts) (common.Address, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SRequests(opts *bind.CallOpts, arg0 *big.Int) (SRequests,

		error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	MakeRequest(opts *bind.TransactOpts, _callbackGasLimit uint32, _requestConfirmations uint16, _numWords uint32) (*types.Transaction, error)

	MakeRequestNative(opts *bind.TransactOpts, _callbackGasLimit uint32, _requestConfirmations uint16, _numWords uint32) (*types.Transaction, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, _requestId *big.Int, _randomWords []*big.Int) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	WithdrawLink(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	WithdrawNative(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusWrapperConsumerExampleOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFV2PlusWrapperConsumerExampleOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VRFV2PlusWrapperConsumerExampleOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusWrapperConsumerExampleOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFV2PlusWrapperConsumerExampleOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VRFV2PlusWrapperConsumerExampleOwnershipTransferred, error)

	FilterWrappedRequestFulfilled(opts *bind.FilterOpts) (*VRFV2PlusWrapperConsumerExampleWrappedRequestFulfilledIterator, error)

	WatchWrappedRequestFulfilled(opts *bind.WatchOpts, sink chan<- *VRFV2PlusWrapperConsumerExampleWrappedRequestFulfilled) (event.Subscription, error)

	ParseWrappedRequestFulfilled(log types.Log) (*VRFV2PlusWrapperConsumerExampleWrappedRequestFulfilled, error)

	FilterWrapperRequestMade(opts *bind.FilterOpts, requestId []*big.Int) (*VRFV2PlusWrapperConsumerExampleWrapperRequestMadeIterator, error)

	WatchWrapperRequestMade(opts *bind.WatchOpts, sink chan<- *VRFV2PlusWrapperConsumerExampleWrapperRequestMade, requestId []*big.Int) (event.Subscription, error)

	ParseWrapperRequestMade(log types.Log) (*VRFV2PlusWrapperConsumerExampleWrapperRequestMade, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
