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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_vrfV2Wrapper\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"LINKAlreadySet\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"payment\",\"type\":\"uint256\"}],\"name\":\"WrappedRequestFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"paid\",\"type\":\"uint256\"}],\"name\":\"WrapperRequestMade\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_requestId\",\"type\":\"uint256\"}],\"name\":\"getRequestStatus\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"paid\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"fulfilled\",\"type\":\"bool\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getWrapper\",\"outputs\":[{\"internalType\":\"contractIVRFV2PlusWrapper\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"_requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"_numWords\",\"type\":\"uint32\"}],\"name\":\"makeRequest\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"_requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"_numWords\",\"type\":\"uint32\"}],\"name\":\"makeRequestNative\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"_randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_requests\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"paid\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"fulfilled\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"native\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"}],\"name\":\"setLinkToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawLink\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawNative\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162001862380380620018628339810160408190526200003491620001e2565b3380600084846001600160a01b038216156200006657600080546001600160a01b0319166001600160a01b0384161790555b600180546001600160a01b0319166001600160a01b03928316179055831615159050620000da5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600280546001600160a01b0319166001600160a01b03848116919091179091558116156200010d576200010d8162000118565b50505050506200021a565b6001600160a01b038116331415620001735760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401620000d1565b600380546001600160a01b0319166001600160a01b03838116918217909255600254604051919216907fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127890600090a350565b80516001600160a01b0381168114620001dd57600080fd5b919050565b60008060408385031215620001f657600080fd5b6200020183620001c5565b91506200021160208401620001c5565b90509250929050565b611638806200022a6000396000f3fe6080604052600436106100d65760003560e01c806384276d811161007f5780639c24ea40116100595780639c24ea4014610236578063a168fa8914610256578063d8a4676f146102b5578063f2fde38b146102e457600080fd5b806384276d811461019f5780638da5cb5b146101bf5780638f6b70701461020b57600080fd5b80631fe543e3116100b05780631fe543e31461014857806379ba50971461016a5780637a8042bd1461017f57600080fd5b80630c09b832146100e257806312065fe0146101155780631e1a34991461012857600080fd5b366100dd57005b600080fd5b3480156100ee57600080fd5b506101026100fd366004611458565b610304565b6040519081526020015b60405180910390f35b34801561012157600080fd5b5047610102565b34801561013457600080fd5b50610102610143366004611458565b6104e9565b34801561015457600080fd5b50610168610163366004611369565b610655565b005b34801561017657600080fd5b5061016861070e565b34801561018b57600080fd5b5061016861019a366004611337565b61080f565b3480156101ab57600080fd5b506101686101ba366004611337565b6108f9565b3480156101cb57600080fd5b5060025473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161010c565b34801561021757600080fd5b5060015473ffffffffffffffffffffffffffffffffffffffff166101e6565b34801561024257600080fd5b506101686102513660046112d8565b6109e9565b34801561026257600080fd5b50610298610271366004611337565b600460205260009081526040902080546001820154600390920154909160ff908116911683565b60408051938452911515602084015215159082015260600161010c565b3480156102c157600080fd5b506102d56102d0366004611337565b610a80565b60405161010c939291906115a8565b3480156102f057600080fd5b506101686102ff3660046112d8565b610ba2565b600061030e610bb6565b600061032a604051806020016040528060001515815250610c39565b905061033885858584610cf5565b6001546040517f4306d35400000000000000000000000000000000000000000000000000000000815263ffffffff8816600482015291935060009173ffffffffffffffffffffffffffffffffffffffff90911690634306d3549060240160206040518083038186803b1580156103ad57600080fd5b505afa1580156103c1573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103e59190611350565b6040805160808101825282815260006020808301828152845183815280830186528486019081526060850184905289845260048352949092208351815591516001830180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055925180519495509193909261047192600285019291019061125f565b5060609190910151600390910180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001691151591909117905560405181815283907f5f56b4c20db9f5b294cbf6f681368de4a992a27e2de2ee702dcf2cbbfa791ec49060200160405180910390a250509392505050565b60006104f3610bb6565b600061050f604051806020016040528060011515815250610c39565b905061051d85858584610eea565b6001546040517f4b16093500000000000000000000000000000000000000000000000000000000815263ffffffff8816600482015291935060009173ffffffffffffffffffffffffffffffffffffffff90911690634b1609359060240160206040518083038186803b15801561059257600080fd5b505afa1580156105a6573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105ca9190611350565b6040805160808101825282815260006020808301828152845183815280830186528486019081526001606086018190528a8552600484529590932084518155905194810180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001695151595909517909455905180519495509193610471926002850192019061125f565b60015473ffffffffffffffffffffffffffffffffffffffff163314610700576040517f08c379a0000000000000000000000000000000000000000000000000000000008152602060048201526024808201527f6f6e6c792056524620563220506c757320777261707065722063616e2066756c60448201527f66696c6c0000000000000000000000000000000000000000000000000000000060648201526084015b60405180910390fd5b61070a8282611051565b5050565b60035473ffffffffffffffffffffffffffffffffffffffff16331461078f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016106f7565b600280547fffffffffffffffffffffffff00000000000000000000000000000000000000008082163390811790935560038054909116905560405173ffffffffffffffffffffffffffffffffffffffff909116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a350565b610817610bb6565b60005473ffffffffffffffffffffffffffffffffffffffff1663a9059cbb61085460025473ffffffffffffffffffffffffffffffffffffffff1690565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e084901b16815273ffffffffffffffffffffffffffffffffffffffff909116600482015260248101849052604401602060405180830381600087803b1580156108c157600080fd5b505af11580156108d5573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061070a9190611315565b610901610bb6565b600061092260025473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff168260405160006040518083038185875af1925050503d8060008114610979576040519150601f19603f3d011682016040523d82523d6000602084013e61097e565b606091505b505090508061070a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f77697468647261774e6174697665206661696c6564000000000000000000000060448201526064016106f7565b60005473ffffffffffffffffffffffffffffffffffffffff1615610a39576040517f64f778ae00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b6000818152600460205260408120548190606090610afa576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f72657175657374206e6f7420666f756e6400000000000000000000000000000060448201526064016106f7565b6000848152600460209081526040808320815160808101835281548152600182015460ff16151581850152600282018054845181870281018701865281815292959394860193830182828015610b6f57602002820191906000526020600020905b815481526020019060010190808311610b5b575b50505091835250506003919091015460ff1615156020918201528151908201516040909201519097919650945092505050565b610baa610bb6565b610bb381611168565b50565b60025473ffffffffffffffffffffffffffffffffffffffff163314610c37576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016106f7565b565b60607f92fd13387c7fe7befbc38d303d6468778fb9731bc4583f17d92989c6fcfdeaaa82604051602401610c7291511515815260200190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009093169290921790915292915050565b600080546001546040517f4306d35400000000000000000000000000000000000000000000000000000000815263ffffffff8816600482015273ffffffffffffffffffffffffffffffffffffffff92831692634000aea09216908190634306d3549060240160206040518083038186803b158015610d7257600080fd5b505afa158015610d86573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610daa9190611350565b88888888604051602001610dc194939291906115c9565b6040516020818303038152906040526040518463ffffffff1660e01b8152600401610dee9392919061154a565b602060405180830381600087803b158015610e0857600080fd5b505af1158015610e1c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610e409190611315565b50600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663fc2a88c36040518163ffffffff1660e01b815260040160206040518083038186803b158015610ea957600080fd5b505afa158015610ebd573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ee19190611350565b95945050505050565b6001546040517f4b16093500000000000000000000000000000000000000000000000000000000815263ffffffff86166004820152600091829173ffffffffffffffffffffffffffffffffffffffff90911690634b1609359060240160206040518083038186803b158015610f5e57600080fd5b505afa158015610f72573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610f969190611350565b6001546040517f9cfc058e00000000000000000000000000000000000000000000000000000000815291925073ffffffffffffffffffffffffffffffffffffffff1690639cfc058e908390610ff5908a908a908a908a906004016115c9565b6020604051808303818588803b15801561100e57600080fd5b505af1158015611022573d6000803e3d6000fd5b50505050506040513d601f19601f820116820180604052508101906110479190611350565b9695505050505050565b6000828152600460205260409020546110c6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f72657175657374206e6f7420666f756e6400000000000000000000000000000060448201526064016106f7565b6000828152600460209081526040909120600181810180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016909117905582516111199260029092019184019061125f565b50600082815260046020526040908190205490517f6c84e12b4c188e61f1b4727024a5cf05c025fa58467e5eedf763c0744c89da7b9161115c918591859161157f565b60405180910390a15050565b73ffffffffffffffffffffffffffffffffffffffff81163314156111e8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016106f7565b600380547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff838116918217909255600254604051919216907fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127890600090a350565b82805482825590600052602060002090810192821561129a579160200282015b8281111561129a57825182559160200191906001019061127f565b506112a69291506112aa565b5090565b5b808211156112a657600081556001016112ab565b803563ffffffff811681146112d357600080fd5b919050565b6000602082840312156112ea57600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461130e57600080fd5b9392505050565b60006020828403121561132757600080fd5b8151801515811461130e57600080fd5b60006020828403121561134957600080fd5b5035919050565b60006020828403121561136257600080fd5b5051919050565b6000806040838503121561137c57600080fd5b8235915060208084013567ffffffffffffffff8082111561139c57600080fd5b818601915086601f8301126113b057600080fd5b8135818111156113c2576113c26115fc565b8060051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f83011681018181108582111715611405576114056115fc565b604052828152858101935084860182860187018b101561142457600080fd5b600095505b83861015611447578035855260019590950194938601938601611429565b508096505050505050509250929050565b60008060006060848603121561146d57600080fd5b611476846112bf565b9250602084013561ffff8116811461148d57600080fd5b915061149b604085016112bf565b90509250925092565b600081518084526020808501945080840160005b838110156114d4578151875295820195908201906001016114b8565b509495945050505050565b6000815180845260005b81811015611505576020818501810151868301820152016114e9565b81811115611517576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b73ffffffffffffffffffffffffffffffffffffffff84168152826020820152606060408201526000610ee160608301846114df565b83815260606020820152600061159860608301856114a4565b9050826040830152949350505050565b8381528215156020820152606060408201526000610ee160608301846114a4565b600063ffffffff808716835261ffff861660208401528085166040840152506080606083015261104760808301846114df565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var VRFV2PlusWrapperConsumerExampleABI = VRFV2PlusWrapperConsumerExampleMetaData.ABI

var VRFV2PlusWrapperConsumerExampleBin = VRFV2PlusWrapperConsumerExampleMetaData.Bin

func DeployVRFV2PlusWrapperConsumerExample(auth *bind.TransactOpts, backend bind.ContractBackend, _link common.Address, _vrfV2Wrapper common.Address) (common.Address, *types.Transaction, *VRFV2PlusWrapperConsumerExample, error) {
	parsed, err := VRFV2PlusWrapperConsumerExampleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFV2PlusWrapperConsumerExampleBin), backend, _link, _vrfV2Wrapper)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFV2PlusWrapperConsumerExample{VRFV2PlusWrapperConsumerExampleCaller: VRFV2PlusWrapperConsumerExampleCaller{contract: contract}, VRFV2PlusWrapperConsumerExampleTransactor: VRFV2PlusWrapperConsumerExampleTransactor{contract: contract}, VRFV2PlusWrapperConsumerExampleFilterer: VRFV2PlusWrapperConsumerExampleFilterer{contract: contract}}, nil
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

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleCaller) GetWrapper(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2PlusWrapperConsumerExample.contract.Call(opts, &out, "getWrapper")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleSession) GetWrapper() (common.Address, error) {
	return _VRFV2PlusWrapperConsumerExample.Contract.GetWrapper(&_VRFV2PlusWrapperConsumerExample.CallOpts)
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleCallerSession) GetWrapper() (common.Address, error) {
	return _VRFV2PlusWrapperConsumerExample.Contract.GetWrapper(&_VRFV2PlusWrapperConsumerExample.CallOpts)
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

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleTransactor) SetLinkToken(opts *bind.TransactOpts, _link common.Address) (*types.Transaction, error) {
	return _VRFV2PlusWrapperConsumerExample.contract.Transact(opts, "setLinkToken", _link)
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleSession) SetLinkToken(_link common.Address) (*types.Transaction, error) {
	return _VRFV2PlusWrapperConsumerExample.Contract.SetLinkToken(&_VRFV2PlusWrapperConsumerExample.TransactOpts, _link)
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleTransactorSession) SetLinkToken(_link common.Address) (*types.Transaction, error) {
	return _VRFV2PlusWrapperConsumerExample.Contract.SetLinkToken(&_VRFV2PlusWrapperConsumerExample.TransactOpts, _link)
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

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusWrapperConsumerExample.contract.RawTransact(opts, nil)
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleSession) Receive() (*types.Transaction, error) {
	return _VRFV2PlusWrapperConsumerExample.Contract.Receive(&_VRFV2PlusWrapperConsumerExample.TransactOpts)
}

func (_VRFV2PlusWrapperConsumerExample *VRFV2PlusWrapperConsumerExampleTransactorSession) Receive() (*types.Transaction, error) {
	return _VRFV2PlusWrapperConsumerExample.Contract.Receive(&_VRFV2PlusWrapperConsumerExample.TransactOpts)
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

	GetRequestStatus(opts *bind.CallOpts, _requestId *big.Int) (GetRequestStatus,

		error)

	GetWrapper(opts *bind.CallOpts) (common.Address, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SRequests(opts *bind.CallOpts, arg0 *big.Int) (SRequests,

		error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	MakeRequest(opts *bind.TransactOpts, _callbackGasLimit uint32, _requestConfirmations uint16, _numWords uint32) (*types.Transaction, error)

	MakeRequestNative(opts *bind.TransactOpts, _callbackGasLimit uint32, _requestConfirmations uint16, _numWords uint32) (*types.Transaction, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, _requestId *big.Int, _randomWords []*big.Int) (*types.Transaction, error)

	SetLinkToken(opts *bind.TransactOpts, _link common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	WithdrawLink(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	WithdrawNative(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	Receive(opts *bind.TransactOpts) (*types.Transaction, error)

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
