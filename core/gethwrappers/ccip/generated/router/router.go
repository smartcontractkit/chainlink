// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package router

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

type ClientAny2EVMMessage struct {
	MessageId           [32]byte
	SourceChainSelector uint64
	Sender              []byte
	Data                []byte
	DestTokenAmounts    []ClientEVMTokenAmount
}

type ClientEVM2AnyMessage struct {
	Receiver     []byte
	Data         []byte
	TokenAmounts []ClientEVMTokenAmount
	FeeToken     common.Address
	ExtraArgs    []byte
}

type ClientEVMTokenAmount struct {
	Token  common.Address
	Amount *big.Int
}

type RouterOffRamp struct {
	SourceChainSelector uint64
	OffRamp             common.Address
}

type RouterOnRamp struct {
	DestChainSelector uint64
	OnRamp            common.Address
}

var RouterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"wrappedNative\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"armProxy\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"BadARMSignal\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedToSendValue\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientFeeTokenAmount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidMsgValue\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"InvalidRecipientAddress\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"offRamp\",\"type\":\"address\"}],\"name\":\"OffRampMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyOffRamp\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"}],\"name\":\"UnsupportedDestinationChain\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"offRamp\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"calldataHash\",\"type\":\"bytes32\"}],\"name\":\"MessageExecuted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"offRamp\",\"type\":\"address\"}],\"name\":\"OffRampAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"offRamp\",\"type\":\"address\"}],\"name\":\"OffRampRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"onRamp\",\"type\":\"address\"}],\"name\":\"OnRampSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"MAX_RET_BYTES\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"onRamp\",\"type\":\"address\"}],\"internalType\":\"structRouter.OnRamp[]\",\"name\":\"onRampUpdates\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"offRamp\",\"type\":\"address\"}],\"internalType\":\"structRouter.OffRamp[]\",\"name\":\"offRampRemoves\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"offRamp\",\"type\":\"address\"}],\"internalType\":\"structRouter.OffRamp[]\",\"name\":\"offRampAdds\",\"type\":\"tuple[]\"}],\"name\":\"applyRampUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destinationChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMTokenAmount[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structClient.EVM2AnyMessage\",\"name\":\"message\",\"type\":\"tuple\"}],\"name\":\"ccipSend\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getArmProxy\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destinationChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMTokenAmount[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structClient.EVM2AnyMessage\",\"name\":\"message\",\"type\":\"tuple\"}],\"name\":\"getFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOffRamps\",\"outputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"offRamp\",\"type\":\"address\"}],\"internalType\":\"structRouter.OffRamp[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"}],\"name\":\"getOnRamp\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"}],\"name\":\"getSupportedTokens\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getWrappedNative\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"}],\"name\":\"isChainSupported\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"offRamp\",\"type\":\"address\"}],\"name\":\"isOffRamp\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"recoverTokens\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMTokenAmount[]\",\"name\":\"destTokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structClient.Any2EVMMessage\",\"name\":\"message\",\"type\":\"tuple\"},{\"internalType\":\"uint16\",\"name\":\"gasForCallExactCheck\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"routeMessage\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"retData\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"wrappedNative\",\"type\":\"address\"}],\"name\":\"setWrappedNative\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b5060405162002d7538038062002d758339810160408190526200003491620001af565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000e7565b5050600280546001600160a01b0319166001600160a01b039485161790555016608052620001e7565b336001600160a01b03821603620001415760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b0381168114620001aa57600080fd5b919050565b60008060408385031215620001c357600080fd5b620001ce8362000192565b9150620001de6020840162000192565b90509250929050565b608051612b6462000211600039600081816101f8015281816105de0152610aef0152612b646000f3fe6080604052600436106101295760003560e01c80638da5cb5b116100a5578063a8d87a3b11610074578063e861e90711610059578063e861e90714610408578063f2fde38b14610433578063fbca3b741461045357600080fd5b8063a8d87a3b1461039b578063da5fcac8146103e857600080fd5b80638da5cb5b146102ec57806396f4e9f914610317578063a40e69c71461032a578063a48a90581461034c57600080fd5b806352cb60ca116100fc578063787350e3116100e1578063787350e31461027f57806379ba5097146102a757806383826b2b146102bc57600080fd5b806352cb60ca1461023d5780635f3e849f1461025f57600080fd5b8063181f5a771461012e57806320487ded1461018d5780633cf97983146101bb5780635246492f146101e9575b600080fd5b34801561013a57600080fd5b506101776040518060400160405280600c81526020017f526f7574657220312e322e30000000000000000000000000000000000000000081525081565b6040516101849190611f61565b60405180910390f35b34801561019957600080fd5b506101ad6101a83660046121d2565b610480565b604051908152602001610184565b3480156101c757600080fd5b506101db6101d63660046122cf565b6105d8565b604051610184929190612347565b3480156101f557600080fd5b507f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610184565b34801561024957600080fd5b5061025d610258366004612362565b61083b565b005b34801561026b57600080fd5b5061025d61027a36600461237f565b61088a565b34801561028b57600080fd5b50610294608481565b60405161ffff9091168152602001610184565b3480156102b357600080fd5b5061025d6109d8565b3480156102c857600080fd5b506102dc6102d73660046123c0565b610ad5565b6040519015158152602001610184565b3480156102f857600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff16610218565b6101ad6103253660046121d2565b610aeb565b34801561033657600080fd5b5061033f61108c565b60405161018491906123f7565b34801561035857600080fd5b506102dc610367366004612466565b67ffffffffffffffff1660009081526003602052604090205473ffffffffffffffffffffffffffffffffffffffff16151590565b3480156103a757600080fd5b506102186103b6366004612466565b67ffffffffffffffff1660009081526003602052604090205473ffffffffffffffffffffffffffffffffffffffff1690565b3480156103f457600080fd5b5061025d6104033660046124cd565b611199565b34801561041457600080fd5b5060025473ffffffffffffffffffffffffffffffffffffffff16610218565b34801561043f57600080fd5b5061025d61044e366004612362565b6114b8565b34801561045f57600080fd5b5061047361046e366004612466565b6114cc565b6040516101849190612567565b606081015160009073ffffffffffffffffffffffffffffffffffffffff166104c15760025473ffffffffffffffffffffffffffffffffffffffff1660608301525b67ffffffffffffffff831660009081526003602052604090205473ffffffffffffffffffffffffffffffffffffffff1680610539576040517fae236d9c00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff851660048201526024015b60405180910390fd5b6040517f20487ded00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8216906320487ded9061058d908790879060040161269e565b602060405180830381865afa1580156105aa573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105ce91906126c1565b9150505b92915050565b600060607f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663397796f76040518163ffffffff1660e01b8152600401602060405180830381865afa158015610647573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061066b91906126da565b156106a2576040517fc148371500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6106c66106be6106b86040890160208a01612466565b336115ec565b600490611630565b6106fc576040517fd2316ede00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006385572ffb60e01b876040516024016107179190612809565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009093169290921790915290506107a481858760848a611648565b90935091507f9b877de93ea9895756e337442c657f95a34fc68e7eb988bdfa693d5be83016b687356107dc60408a0160208b01612466565b83516020850120604051610829939291339193845267ffffffffffffffff92909216602084015273ffffffffffffffffffffffffffffffffffffffff166040830152606082015260800190565b60405180910390a15094509492505050565b61084361176b565b600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b61089261176b565b73ffffffffffffffffffffffffffffffffffffffff82166108f7576040517f26a78f8f00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff83166004820152602401610530565b73ffffffffffffffffffffffffffffffffffffffff83166109b25760008273ffffffffffffffffffffffffffffffffffffffff168260405160006040518083038185875af1925050503d806000811461096c576040519150601f19603f3d011682016040523d82523d6000602084013e610971565b606091505b50509050806109ac576040517fe417b80b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b50505050565b6109d373ffffffffffffffffffffffffffffffffffffffff841683836117ee565b505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610a59576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610530565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6000610ae46106be84846115ec565b9392505050565b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663397796f76040518163ffffffff1660e01b8152600401602060405180830381865afa158015610b58573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b7c91906126da565b15610bb3576040517fc148371500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff831660009081526003602052604090205473ffffffffffffffffffffffffffffffffffffffff1680610c26576040517fae236d9c00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff85166004820152602401610530565b606083015160009073ffffffffffffffffffffffffffffffffffffffff16610db85760025473ffffffffffffffffffffffffffffffffffffffff90811660608601526040517f20487ded000000000000000000000000000000000000000000000000000000008152908316906320487ded90610ca8908890889060040161269e565b602060405180830381865afa158015610cc5573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ce991906126c1565b905080341015610d25576040517f07da6ee600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b349050836060015173ffffffffffffffffffffffffffffffffffffffff1663d0e30db0826040518263ffffffff1660e01b81526004016000604051808303818588803b158015610d7457600080fd5b505af1158015610d88573d6000803e3d6000fd5b505050506060850151610db3915073ffffffffffffffffffffffffffffffffffffffff1683836117ee565b610eaf565b3415610df0576040517f1841b4e100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040517f20487ded00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8316906320487ded90610e44908890889060040161269e565b602060405180830381865afa158015610e61573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610e8591906126c1565b6060850151909150610eaf9073ffffffffffffffffffffffffffffffffffffffff163384846118c2565b60005b846040015151811015610fe757600085604001518281518110610ed757610ed7612915565b6020908102919091010151516040517f48a98aa400000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8916600482015273ffffffffffffffffffffffffffffffffffffffff8083166024830152919250610fd6913391908716906348a98aa490604401602060405180830381865afa158015610f69573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610f8d9190612944565b88604001518581518110610fa357610fa3612915565b6020026020010151602001518473ffffffffffffffffffffffffffffffffffffffff166118c2909392919063ffffffff16565b50610fe081612990565b9050610eb2565b506040517fdf0aa9e900000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff83169063df0aa9e9906110409088908890869033906004016129c8565b6020604051808303816000875af115801561105f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061108391906126c1565b95945050505050565b6060600061109a6004611920565b90506000815167ffffffffffffffff8111156110b8576110b8611f91565b6040519080825280602002602001820160405280156110fd57816020015b60408051808201909152600080825260208201528152602001906001900390816110d65790505b50905060005b825181101561119257600083828151811061112057611120612915565b60200260200101519050604051806040016040528060a083901c67ffffffffffffffff1681526020018273ffffffffffffffffffffffffffffffffffffffff1681525083838151811061117557611175612915565b6020026020010181905250508061118b90612990565b9050611103565b5092915050565b6111a161176b565b60005b858110156112855760008787838181106111c0576111c0612915565b9050604002018036038101906111d69190612a18565b60208181018051835167ffffffffffffffff90811660009081526003855260409081902080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff948516179055855193519051921682529394509216917f1f7d0ec248b80e5c0dde0ee531c4fc8fdb6ce9a2b3d90f560c74acd6a7202f23910160405180910390a25061127e81612990565b90506111a4565b5060005b838110156113c65760008585838181106112a5576112a5612915565b6112bb9260206040909202019081019150612466565b905060008686848181106112d1576112d1612915565b90506040020160200160208101906112e99190612362565b90506113006112f883836115ec565b60049061192d565b61135e576040517f4964779000000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8316600482015273ffffffffffffffffffffffffffffffffffffffff82166024820152604401610530565b60405173ffffffffffffffffffffffffffffffffffffffff8216815267ffffffffffffffff8316907fa823809efda3ba66c873364eec120fa0923d9fabda73bc97dd5663341e2d9bcb9060200160405180910390a25050806113bf90612990565b9050611289565b5060005b818110156114af5760008383838181106113e6576113e6612915565b6113fc9260206040909202019081019150612466565b9050600084848481811061141257611412612915565b905060400201602001602081019061142a9190612362565b905061144161143983836115ec565b600490611939565b1561149c5760405173ffffffffffffffffffffffffffffffffffffffff8216815267ffffffffffffffff8316907fa4bdf64ebdf3316320601a081916a75aa144bcef6c4beeb0e9fb1982cacc6b949060200160405180910390a25b5050806114a890612990565b90506113ca565b50505050505050565b6114c061176b565b6114c981611945565b50565b60606115068267ffffffffffffffff1660009081526003602052604090205473ffffffffffffffffffffffffffffffffffffffff16151590565b611520576040805160008082526020820190925290611192565b67ffffffffffffffff8216600081815260036020526040908190205490517ffbca3b74000000000000000000000000000000000000000000000000000000008152600481019290925273ffffffffffffffffffffffffffffffffffffffff169063fbca3b7490602401600060405180830381865afa1580156115a6573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526105d29190810190612a57565b6000610ae473ffffffffffffffffffffffffffffffffffffffff83167bffffffffffffffff000000000000000000000000000000000000000060a086901b16612ae6565b60008181526001830160205260408120541515610ae4565b600060608361ffff1667ffffffffffffffff81111561166957611669611f91565b6040519080825280601f01601f191660200182016040528015611693576020820181803683370190505b5090507f0c3b563c000000000000000000000000000000000000000000000000000000007fafa32a2c000000000000000000000000000000000000000000000000000000007f37c3be2900000000000000000000000000000000000000000000000000000000883b611709578260005260046000fd5b5a8681101561171c578260005260046000fd5b8690036040810481038910611735578160005260046000fd5b506000808b5160208d0160008d8df194503d878111156117525750865b808552806000602087013e505050509550959350505050565b60005473ffffffffffffffffffffffffffffffffffffffff1633146117ec576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610530565b565b60405173ffffffffffffffffffffffffffffffffffffffff83166024820152604481018290526109d39084907fa9059cbb00000000000000000000000000000000000000000000000000000000906064015b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090931692909217909152611a3a565b60405173ffffffffffffffffffffffffffffffffffffffff808516602483015283166044820152606481018290526109ac9085907f23b872dd0000000000000000000000000000000000000000000000000000000090608401611840565b60606000610ae483611b46565b6000610ae48383611ba2565b6000610ae48383611c95565b3373ffffffffffffffffffffffffffffffffffffffff8216036119c4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610530565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000611a9c826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff16611ce49092919063ffffffff16565b8051909150156109d35780806020019051810190611aba91906126da565b6109d3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f742073756363656564000000000000000000000000000000000000000000006064820152608401610530565b606081600001805480602002602001604051908101604052809291908181526020018280548015611b9657602002820191906000526020600020905b815481526020019060010190808311611b82575b50505050509050919050565b60008181526001830160205260408120548015611c8b576000611bc6600183612af9565b8554909150600090611bda90600190612af9565b9050818114611c3f576000866000018281548110611bfa57611bfa612915565b9060005260206000200154905080876000018481548110611c1d57611c1d612915565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080611c5057611c50612b0c565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506105d2565b60009150506105d2565b6000818152600183016020526040812054611cdc575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556105d2565b5060006105d2565b6060611cf38484600085611cfb565b949350505050565b606082471015611d8d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f60448201527f722063616c6c00000000000000000000000000000000000000000000000000006064820152608401610530565b6000808673ffffffffffffffffffffffffffffffffffffffff168587604051611db69190612b3b565b60006040518083038185875af1925050503d8060008114611df3576040519150601f19603f3d011682016040523d82523d6000602084013e611df8565b606091505b5091509150611e0987838387611e14565b979650505050505050565b60608315611eaa578251600003611ea35773ffffffffffffffffffffffffffffffffffffffff85163b611ea3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401610530565b5081611cf3565b611cf38383815115611ebf5781518083602001fd5b806040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105309190611f61565b60005b83811015611f0e578181015183820152602001611ef6565b50506000910152565b60008151808452611f2f816020860160208601611ef3565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000610ae46020830184611f17565b803567ffffffffffffffff81168114611f8c57600080fd5b919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040805190810167ffffffffffffffff81118282101715611fe357611fe3611f91565b60405290565b60405160a0810167ffffffffffffffff81118282101715611fe357611fe3611f91565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561205357612053611f91565b604052919050565b600082601f83011261206c57600080fd5b813567ffffffffffffffff81111561208657612086611f91565b6120b760207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8401160161200c565b8181528460208386010111156120cc57600080fd5b816020850160208301376000918101602001919091529392505050565b600067ffffffffffffffff82111561210357612103611f91565b5060051b60200190565b73ffffffffffffffffffffffffffffffffffffffff811681146114c957600080fd5b8035611f8c8161210d565b600082601f83011261214b57600080fd5b8135602061216061215b836120e9565b61200c565b82815260069290921b8401810191818101908684111561217f57600080fd5b8286015b848110156121c7576040818903121561219c5760008081fd5b6121a4611fc0565b81356121af8161210d565b81528185013585820152835291830191604001612183565b509695505050505050565b600080604083850312156121e557600080fd5b6121ee83611f74565b9150602083013567ffffffffffffffff8082111561220b57600080fd5b9084019060a0828703121561221f57600080fd5b612227611fe9565b82358281111561223657600080fd5b6122428882860161205b565b82525060208301358281111561225757600080fd5b6122638882860161205b565b60208301525060408301358281111561227b57600080fd5b6122878882860161213a565b6040830152506122996060840161212f565b60608201526080830135828111156122b057600080fd5b6122bc8882860161205b565b6080830152508093505050509250929050565b600080600080608085870312156122e557600080fd5b843567ffffffffffffffff8111156122fc57600080fd5b850160a0818803121561230e57600080fd5b9350602085013561ffff8116811461232557600080fd5b925060408501359150606085013561233c8161210d565b939692955090935050565b8215158152604060208201526000611cf36040830184611f17565b60006020828403121561237457600080fd5b8135610ae48161210d565b60008060006060848603121561239457600080fd5b833561239f8161210d565b925060208401356123af8161210d565b929592945050506040919091013590565b600080604083850312156123d357600080fd5b6123dc83611f74565b915060208301356123ec8161210d565b809150509250929050565b602080825282518282018190526000919060409081850190868401855b82811015612459578151805167ffffffffffffffff16855286015173ffffffffffffffffffffffffffffffffffffffff16868501529284019290850190600101612414565b5091979650505050505050565b60006020828403121561247857600080fd5b610ae482611f74565b60008083601f84011261249357600080fd5b50813567ffffffffffffffff8111156124ab57600080fd5b6020830191508360208260061b85010111156124c657600080fd5b9250929050565b600080600080600080606087890312156124e657600080fd5b863567ffffffffffffffff808211156124fe57600080fd5b61250a8a838b01612481565b9098509650602089013591508082111561252357600080fd5b61252f8a838b01612481565b9096509450604089013591508082111561254857600080fd5b5061255589828a01612481565b979a9699509497509295939492505050565b6020808252825182820181905260009190848201906040850190845b818110156125b557835173ffffffffffffffffffffffffffffffffffffffff1683529284019291840191600101612583565b50909695505050505050565b6000815160a084526125d660a0850182611f17565b9050602080840151858303828701526125ef8382611f17565b60408681015188830389830152805180845290850195509092506000918401905b8083101561264f578551805173ffffffffffffffffffffffffffffffffffffffff16835285015185830152948401946001929092019190830190612610565b5060608701519450612679606089018673ffffffffffffffffffffffffffffffffffffffff169052565b6080870151945087810360808901526126928186611f17565b98975050505050505050565b67ffffffffffffffff83168152604060208201526000611cf360408301846125c1565b6000602082840312156126d357600080fd5b5051919050565b6000602082840312156126ec57600080fd5b81518015158114610ae457600080fd5b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261273157600080fd5b830160208101925035905067ffffffffffffffff81111561275157600080fd5b8036038213156124c657600080fd5b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b8183526000602080850194508260005b858110156127fe5781356127cc8161210d565b73ffffffffffffffffffffffffffffffffffffffff1687528183013583880152604096870196909101906001016127b9565b509495945050505050565b6020815281356020820152600061282260208401611f74565b67ffffffffffffffff808216604085015261284060408601866126fc565b925060a0606086015261285760c086018483612760565b92505061286760608601866126fc565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08087860301608088015261289d858385612760565b9450608088013592507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18836030183126128d657600080fd5b602092880192830192359150838211156128ef57600080fd5b8160061b360383131561290157600080fd5b8685030160a0870152611e098482846127a9565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60006020828403121561295657600080fd5b8151610ae48161210d565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036129c1576129c1612961565b5060010190565b67ffffffffffffffff851681526080602082015260006129eb60808301866125c1565b905083604083015273ffffffffffffffffffffffffffffffffffffffff8316606083015295945050505050565b600060408284031215612a2a57600080fd5b612a32611fc0565b612a3b83611f74565b81526020830135612a4b8161210d565b60208201529392505050565b60006020808385031215612a6a57600080fd5b825167ffffffffffffffff811115612a8157600080fd5b8301601f81018513612a9257600080fd5b8051612aa061215b826120e9565b81815260059190911b82018301908381019087831115612abf57600080fd5b928401925b82841015611e09578351612ad78161210d565b82529284019290840190612ac4565b808201808211156105d2576105d2612961565b818103818111156105d2576105d2612961565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b60008251612b4d818460208701611ef3565b919091019291505056fea164736f6c6343000813000a",
}

var RouterABI = RouterMetaData.ABI

var RouterBin = RouterMetaData.Bin

func DeployRouter(auth *bind.TransactOpts, backend bind.ContractBackend, wrappedNative common.Address, armProxy common.Address) (common.Address, *types.Transaction, *Router, error) {
	parsed, err := RouterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(RouterBin), backend, wrappedNative, armProxy)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Router{RouterCaller: RouterCaller{contract: contract}, RouterTransactor: RouterTransactor{contract: contract}, RouterFilterer: RouterFilterer{contract: contract}}, nil
}

type Router struct {
	address common.Address
	abi     abi.ABI
	RouterCaller
	RouterTransactor
	RouterFilterer
}

type RouterCaller struct {
	contract *bind.BoundContract
}

type RouterTransactor struct {
	contract *bind.BoundContract
}

type RouterFilterer struct {
	contract *bind.BoundContract
}

type RouterSession struct {
	Contract     *Router
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type RouterCallerSession struct {
	Contract *RouterCaller
	CallOpts bind.CallOpts
}

type RouterTransactorSession struct {
	Contract     *RouterTransactor
	TransactOpts bind.TransactOpts
}

type RouterRaw struct {
	Contract *Router
}

type RouterCallerRaw struct {
	Contract *RouterCaller
}

type RouterTransactorRaw struct {
	Contract *RouterTransactor
}

func NewRouter(address common.Address, backend bind.ContractBackend) (*Router, error) {
	abi, err := abi.JSON(strings.NewReader(RouterABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindRouter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Router{address: address, abi: abi, RouterCaller: RouterCaller{contract: contract}, RouterTransactor: RouterTransactor{contract: contract}, RouterFilterer: RouterFilterer{contract: contract}}, nil
}

func NewRouterCaller(address common.Address, caller bind.ContractCaller) (*RouterCaller, error) {
	contract, err := bindRouter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RouterCaller{contract: contract}, nil
}

func NewRouterTransactor(address common.Address, transactor bind.ContractTransactor) (*RouterTransactor, error) {
	contract, err := bindRouter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RouterTransactor{contract: contract}, nil
}

func NewRouterFilterer(address common.Address, filterer bind.ContractFilterer) (*RouterFilterer, error) {
	contract, err := bindRouter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RouterFilterer{contract: contract}, nil
}

func bindRouter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RouterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_Router *RouterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Router.Contract.RouterCaller.contract.Call(opts, result, method, params...)
}

func (_Router *RouterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Router.Contract.RouterTransactor.contract.Transfer(opts)
}

func (_Router *RouterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Router.Contract.RouterTransactor.contract.Transact(opts, method, params...)
}

func (_Router *RouterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Router.Contract.contract.Call(opts, result, method, params...)
}

func (_Router *RouterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Router.Contract.contract.Transfer(opts)
}

func (_Router *RouterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Router.Contract.contract.Transact(opts, method, params...)
}

func (_Router *RouterCaller) MAXRETBYTES(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _Router.contract.Call(opts, &out, "MAX_RET_BYTES")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_Router *RouterSession) MAXRETBYTES() (uint16, error) {
	return _Router.Contract.MAXRETBYTES(&_Router.CallOpts)
}

func (_Router *RouterCallerSession) MAXRETBYTES() (uint16, error) {
	return _Router.Contract.MAXRETBYTES(&_Router.CallOpts)
}

func (_Router *RouterCaller) GetArmProxy(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Router.contract.Call(opts, &out, "getArmProxy")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_Router *RouterSession) GetArmProxy() (common.Address, error) {
	return _Router.Contract.GetArmProxy(&_Router.CallOpts)
}

func (_Router *RouterCallerSession) GetArmProxy() (common.Address, error) {
	return _Router.Contract.GetArmProxy(&_Router.CallOpts)
}

func (_Router *RouterCaller) GetFee(opts *bind.CallOpts, destinationChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error) {
	var out []interface{}
	err := _Router.contract.Call(opts, &out, "getFee", destinationChainSelector, message)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_Router *RouterSession) GetFee(destinationChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error) {
	return _Router.Contract.GetFee(&_Router.CallOpts, destinationChainSelector, message)
}

func (_Router *RouterCallerSession) GetFee(destinationChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error) {
	return _Router.Contract.GetFee(&_Router.CallOpts, destinationChainSelector, message)
}

func (_Router *RouterCaller) GetOffRamps(opts *bind.CallOpts) ([]RouterOffRamp, error) {
	var out []interface{}
	err := _Router.contract.Call(opts, &out, "getOffRamps")

	if err != nil {
		return *new([]RouterOffRamp), err
	}

	out0 := *abi.ConvertType(out[0], new([]RouterOffRamp)).(*[]RouterOffRamp)

	return out0, err

}

func (_Router *RouterSession) GetOffRamps() ([]RouterOffRamp, error) {
	return _Router.Contract.GetOffRamps(&_Router.CallOpts)
}

func (_Router *RouterCallerSession) GetOffRamps() ([]RouterOffRamp, error) {
	return _Router.Contract.GetOffRamps(&_Router.CallOpts)
}

func (_Router *RouterCaller) GetOnRamp(opts *bind.CallOpts, destChainSelector uint64) (common.Address, error) {
	var out []interface{}
	err := _Router.contract.Call(opts, &out, "getOnRamp", destChainSelector)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_Router *RouterSession) GetOnRamp(destChainSelector uint64) (common.Address, error) {
	return _Router.Contract.GetOnRamp(&_Router.CallOpts, destChainSelector)
}

func (_Router *RouterCallerSession) GetOnRamp(destChainSelector uint64) (common.Address, error) {
	return _Router.Contract.GetOnRamp(&_Router.CallOpts, destChainSelector)
}

func (_Router *RouterCaller) GetSupportedTokens(opts *bind.CallOpts, chainSelector uint64) ([]common.Address, error) {
	var out []interface{}
	err := _Router.contract.Call(opts, &out, "getSupportedTokens", chainSelector)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_Router *RouterSession) GetSupportedTokens(chainSelector uint64) ([]common.Address, error) {
	return _Router.Contract.GetSupportedTokens(&_Router.CallOpts, chainSelector)
}

func (_Router *RouterCallerSession) GetSupportedTokens(chainSelector uint64) ([]common.Address, error) {
	return _Router.Contract.GetSupportedTokens(&_Router.CallOpts, chainSelector)
}

func (_Router *RouterCaller) GetWrappedNative(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Router.contract.Call(opts, &out, "getWrappedNative")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_Router *RouterSession) GetWrappedNative() (common.Address, error) {
	return _Router.Contract.GetWrappedNative(&_Router.CallOpts)
}

func (_Router *RouterCallerSession) GetWrappedNative() (common.Address, error) {
	return _Router.Contract.GetWrappedNative(&_Router.CallOpts)
}

func (_Router *RouterCaller) IsChainSupported(opts *bind.CallOpts, chainSelector uint64) (bool, error) {
	var out []interface{}
	err := _Router.contract.Call(opts, &out, "isChainSupported", chainSelector)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_Router *RouterSession) IsChainSupported(chainSelector uint64) (bool, error) {
	return _Router.Contract.IsChainSupported(&_Router.CallOpts, chainSelector)
}

func (_Router *RouterCallerSession) IsChainSupported(chainSelector uint64) (bool, error) {
	return _Router.Contract.IsChainSupported(&_Router.CallOpts, chainSelector)
}

func (_Router *RouterCaller) IsOffRamp(opts *bind.CallOpts, sourceChainSelector uint64, offRamp common.Address) (bool, error) {
	var out []interface{}
	err := _Router.contract.Call(opts, &out, "isOffRamp", sourceChainSelector, offRamp)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_Router *RouterSession) IsOffRamp(sourceChainSelector uint64, offRamp common.Address) (bool, error) {
	return _Router.Contract.IsOffRamp(&_Router.CallOpts, sourceChainSelector, offRamp)
}

func (_Router *RouterCallerSession) IsOffRamp(sourceChainSelector uint64, offRamp common.Address) (bool, error) {
	return _Router.Contract.IsOffRamp(&_Router.CallOpts, sourceChainSelector, offRamp)
}

func (_Router *RouterCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Router.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_Router *RouterSession) Owner() (common.Address, error) {
	return _Router.Contract.Owner(&_Router.CallOpts)
}

func (_Router *RouterCallerSession) Owner() (common.Address, error) {
	return _Router.Contract.Owner(&_Router.CallOpts)
}

func (_Router *RouterCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Router.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_Router *RouterSession) TypeAndVersion() (string, error) {
	return _Router.Contract.TypeAndVersion(&_Router.CallOpts)
}

func (_Router *RouterCallerSession) TypeAndVersion() (string, error) {
	return _Router.Contract.TypeAndVersion(&_Router.CallOpts)
}

func (_Router *RouterTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Router.contract.Transact(opts, "acceptOwnership")
}

func (_Router *RouterSession) AcceptOwnership() (*types.Transaction, error) {
	return _Router.Contract.AcceptOwnership(&_Router.TransactOpts)
}

func (_Router *RouterTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Router.Contract.AcceptOwnership(&_Router.TransactOpts)
}

func (_Router *RouterTransactor) ApplyRampUpdates(opts *bind.TransactOpts, onRampUpdates []RouterOnRamp, offRampRemoves []RouterOffRamp, offRampAdds []RouterOffRamp) (*types.Transaction, error) {
	return _Router.contract.Transact(opts, "applyRampUpdates", onRampUpdates, offRampRemoves, offRampAdds)
}

func (_Router *RouterSession) ApplyRampUpdates(onRampUpdates []RouterOnRamp, offRampRemoves []RouterOffRamp, offRampAdds []RouterOffRamp) (*types.Transaction, error) {
	return _Router.Contract.ApplyRampUpdates(&_Router.TransactOpts, onRampUpdates, offRampRemoves, offRampAdds)
}

func (_Router *RouterTransactorSession) ApplyRampUpdates(onRampUpdates []RouterOnRamp, offRampRemoves []RouterOffRamp, offRampAdds []RouterOffRamp) (*types.Transaction, error) {
	return _Router.Contract.ApplyRampUpdates(&_Router.TransactOpts, onRampUpdates, offRampRemoves, offRampAdds)
}

func (_Router *RouterTransactor) CcipSend(opts *bind.TransactOpts, destinationChainSelector uint64, message ClientEVM2AnyMessage) (*types.Transaction, error) {
	return _Router.contract.Transact(opts, "ccipSend", destinationChainSelector, message)
}

func (_Router *RouterSession) CcipSend(destinationChainSelector uint64, message ClientEVM2AnyMessage) (*types.Transaction, error) {
	return _Router.Contract.CcipSend(&_Router.TransactOpts, destinationChainSelector, message)
}

func (_Router *RouterTransactorSession) CcipSend(destinationChainSelector uint64, message ClientEVM2AnyMessage) (*types.Transaction, error) {
	return _Router.Contract.CcipSend(&_Router.TransactOpts, destinationChainSelector, message)
}

func (_Router *RouterTransactor) RecoverTokens(opts *bind.TransactOpts, tokenAddress common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Router.contract.Transact(opts, "recoverTokens", tokenAddress, to, amount)
}

func (_Router *RouterSession) RecoverTokens(tokenAddress common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Router.Contract.RecoverTokens(&_Router.TransactOpts, tokenAddress, to, amount)
}

func (_Router *RouterTransactorSession) RecoverTokens(tokenAddress common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Router.Contract.RecoverTokens(&_Router.TransactOpts, tokenAddress, to, amount)
}

func (_Router *RouterTransactor) RouteMessage(opts *bind.TransactOpts, message ClientAny2EVMMessage, gasForCallExactCheck uint16, gasLimit *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _Router.contract.Transact(opts, "routeMessage", message, gasForCallExactCheck, gasLimit, receiver)
}

func (_Router *RouterSession) RouteMessage(message ClientAny2EVMMessage, gasForCallExactCheck uint16, gasLimit *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _Router.Contract.RouteMessage(&_Router.TransactOpts, message, gasForCallExactCheck, gasLimit, receiver)
}

func (_Router *RouterTransactorSession) RouteMessage(message ClientAny2EVMMessage, gasForCallExactCheck uint16, gasLimit *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _Router.Contract.RouteMessage(&_Router.TransactOpts, message, gasForCallExactCheck, gasLimit, receiver)
}

func (_Router *RouterTransactor) SetWrappedNative(opts *bind.TransactOpts, wrappedNative common.Address) (*types.Transaction, error) {
	return _Router.contract.Transact(opts, "setWrappedNative", wrappedNative)
}

func (_Router *RouterSession) SetWrappedNative(wrappedNative common.Address) (*types.Transaction, error) {
	return _Router.Contract.SetWrappedNative(&_Router.TransactOpts, wrappedNative)
}

func (_Router *RouterTransactorSession) SetWrappedNative(wrappedNative common.Address) (*types.Transaction, error) {
	return _Router.Contract.SetWrappedNative(&_Router.TransactOpts, wrappedNative)
}

func (_Router *RouterTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _Router.contract.Transact(opts, "transferOwnership", to)
}

func (_Router *RouterSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _Router.Contract.TransferOwnership(&_Router.TransactOpts, to)
}

func (_Router *RouterTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _Router.Contract.TransferOwnership(&_Router.TransactOpts, to)
}

type RouterMessageExecutedIterator struct {
	Event *RouterMessageExecuted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RouterMessageExecutedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RouterMessageExecuted)
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
		it.Event = new(RouterMessageExecuted)
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

func (it *RouterMessageExecutedIterator) Error() error {
	return it.fail
}

func (it *RouterMessageExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RouterMessageExecuted struct {
	MessageId           [32]byte
	SourceChainSelector uint64
	OffRamp             common.Address
	CalldataHash        [32]byte
	Raw                 types.Log
}

func (_Router *RouterFilterer) FilterMessageExecuted(opts *bind.FilterOpts) (*RouterMessageExecutedIterator, error) {

	logs, sub, err := _Router.contract.FilterLogs(opts, "MessageExecuted")
	if err != nil {
		return nil, err
	}
	return &RouterMessageExecutedIterator{contract: _Router.contract, event: "MessageExecuted", logs: logs, sub: sub}, nil
}

func (_Router *RouterFilterer) WatchMessageExecuted(opts *bind.WatchOpts, sink chan<- *RouterMessageExecuted) (event.Subscription, error) {

	logs, sub, err := _Router.contract.WatchLogs(opts, "MessageExecuted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RouterMessageExecuted)
				if err := _Router.contract.UnpackLog(event, "MessageExecuted", log); err != nil {
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

func (_Router *RouterFilterer) ParseMessageExecuted(log types.Log) (*RouterMessageExecuted, error) {
	event := new(RouterMessageExecuted)
	if err := _Router.contract.UnpackLog(event, "MessageExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RouterOffRampAddedIterator struct {
	Event *RouterOffRampAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RouterOffRampAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RouterOffRampAdded)
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
		it.Event = new(RouterOffRampAdded)
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

func (it *RouterOffRampAddedIterator) Error() error {
	return it.fail
}

func (it *RouterOffRampAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RouterOffRampAdded struct {
	SourceChainSelector uint64
	OffRamp             common.Address
	Raw                 types.Log
}

func (_Router *RouterFilterer) FilterOffRampAdded(opts *bind.FilterOpts, sourceChainSelector []uint64) (*RouterOffRampAddedIterator, error) {

	var sourceChainSelectorRule []interface{}
	for _, sourceChainSelectorItem := range sourceChainSelector {
		sourceChainSelectorRule = append(sourceChainSelectorRule, sourceChainSelectorItem)
	}

	logs, sub, err := _Router.contract.FilterLogs(opts, "OffRampAdded", sourceChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &RouterOffRampAddedIterator{contract: _Router.contract, event: "OffRampAdded", logs: logs, sub: sub}, nil
}

func (_Router *RouterFilterer) WatchOffRampAdded(opts *bind.WatchOpts, sink chan<- *RouterOffRampAdded, sourceChainSelector []uint64) (event.Subscription, error) {

	var sourceChainSelectorRule []interface{}
	for _, sourceChainSelectorItem := range sourceChainSelector {
		sourceChainSelectorRule = append(sourceChainSelectorRule, sourceChainSelectorItem)
	}

	logs, sub, err := _Router.contract.WatchLogs(opts, "OffRampAdded", sourceChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RouterOffRampAdded)
				if err := _Router.contract.UnpackLog(event, "OffRampAdded", log); err != nil {
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

func (_Router *RouterFilterer) ParseOffRampAdded(log types.Log) (*RouterOffRampAdded, error) {
	event := new(RouterOffRampAdded)
	if err := _Router.contract.UnpackLog(event, "OffRampAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RouterOffRampRemovedIterator struct {
	Event *RouterOffRampRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RouterOffRampRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RouterOffRampRemoved)
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
		it.Event = new(RouterOffRampRemoved)
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

func (it *RouterOffRampRemovedIterator) Error() error {
	return it.fail
}

func (it *RouterOffRampRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RouterOffRampRemoved struct {
	SourceChainSelector uint64
	OffRamp             common.Address
	Raw                 types.Log
}

func (_Router *RouterFilterer) FilterOffRampRemoved(opts *bind.FilterOpts, sourceChainSelector []uint64) (*RouterOffRampRemovedIterator, error) {

	var sourceChainSelectorRule []interface{}
	for _, sourceChainSelectorItem := range sourceChainSelector {
		sourceChainSelectorRule = append(sourceChainSelectorRule, sourceChainSelectorItem)
	}

	logs, sub, err := _Router.contract.FilterLogs(opts, "OffRampRemoved", sourceChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &RouterOffRampRemovedIterator{contract: _Router.contract, event: "OffRampRemoved", logs: logs, sub: sub}, nil
}

func (_Router *RouterFilterer) WatchOffRampRemoved(opts *bind.WatchOpts, sink chan<- *RouterOffRampRemoved, sourceChainSelector []uint64) (event.Subscription, error) {

	var sourceChainSelectorRule []interface{}
	for _, sourceChainSelectorItem := range sourceChainSelector {
		sourceChainSelectorRule = append(sourceChainSelectorRule, sourceChainSelectorItem)
	}

	logs, sub, err := _Router.contract.WatchLogs(opts, "OffRampRemoved", sourceChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RouterOffRampRemoved)
				if err := _Router.contract.UnpackLog(event, "OffRampRemoved", log); err != nil {
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

func (_Router *RouterFilterer) ParseOffRampRemoved(log types.Log) (*RouterOffRampRemoved, error) {
	event := new(RouterOffRampRemoved)
	if err := _Router.contract.UnpackLog(event, "OffRampRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RouterOnRampSetIterator struct {
	Event *RouterOnRampSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RouterOnRampSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RouterOnRampSet)
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
		it.Event = new(RouterOnRampSet)
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

func (it *RouterOnRampSetIterator) Error() error {
	return it.fail
}

func (it *RouterOnRampSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RouterOnRampSet struct {
	DestChainSelector uint64
	OnRamp            common.Address
	Raw               types.Log
}

func (_Router *RouterFilterer) FilterOnRampSet(opts *bind.FilterOpts, destChainSelector []uint64) (*RouterOnRampSetIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _Router.contract.FilterLogs(opts, "OnRampSet", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &RouterOnRampSetIterator{contract: _Router.contract, event: "OnRampSet", logs: logs, sub: sub}, nil
}

func (_Router *RouterFilterer) WatchOnRampSet(opts *bind.WatchOpts, sink chan<- *RouterOnRampSet, destChainSelector []uint64) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _Router.contract.WatchLogs(opts, "OnRampSet", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RouterOnRampSet)
				if err := _Router.contract.UnpackLog(event, "OnRampSet", log); err != nil {
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

func (_Router *RouterFilterer) ParseOnRampSet(log types.Log) (*RouterOnRampSet, error) {
	event := new(RouterOnRampSet)
	if err := _Router.contract.UnpackLog(event, "OnRampSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RouterOwnershipTransferRequestedIterator struct {
	Event *RouterOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RouterOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RouterOwnershipTransferRequested)
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
		it.Event = new(RouterOwnershipTransferRequested)
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

func (it *RouterOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *RouterOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RouterOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_Router *RouterFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RouterOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Router.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &RouterOwnershipTransferRequestedIterator{contract: _Router.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_Router *RouterFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *RouterOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Router.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RouterOwnershipTransferRequested)
				if err := _Router.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_Router *RouterFilterer) ParseOwnershipTransferRequested(log types.Log) (*RouterOwnershipTransferRequested, error) {
	event := new(RouterOwnershipTransferRequested)
	if err := _Router.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RouterOwnershipTransferredIterator struct {
	Event *RouterOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RouterOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RouterOwnershipTransferred)
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
		it.Event = new(RouterOwnershipTransferred)
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

func (it *RouterOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *RouterOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RouterOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_Router *RouterFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RouterOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Router.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &RouterOwnershipTransferredIterator{contract: _Router.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_Router *RouterFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *RouterOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Router.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RouterOwnershipTransferred)
				if err := _Router.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_Router *RouterFilterer) ParseOwnershipTransferred(log types.Log) (*RouterOwnershipTransferred, error) {
	event := new(RouterOwnershipTransferred)
	if err := _Router.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_Router *Router) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _Router.abi.Events["MessageExecuted"].ID:
		return _Router.ParseMessageExecuted(log)
	case _Router.abi.Events["OffRampAdded"].ID:
		return _Router.ParseOffRampAdded(log)
	case _Router.abi.Events["OffRampRemoved"].ID:
		return _Router.ParseOffRampRemoved(log)
	case _Router.abi.Events["OnRampSet"].ID:
		return _Router.ParseOnRampSet(log)
	case _Router.abi.Events["OwnershipTransferRequested"].ID:
		return _Router.ParseOwnershipTransferRequested(log)
	case _Router.abi.Events["OwnershipTransferred"].ID:
		return _Router.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (RouterMessageExecuted) Topic() common.Hash {
	return common.HexToHash("0x9b877de93ea9895756e337442c657f95a34fc68e7eb988bdfa693d5be83016b6")
}

func (RouterOffRampAdded) Topic() common.Hash {
	return common.HexToHash("0xa4bdf64ebdf3316320601a081916a75aa144bcef6c4beeb0e9fb1982cacc6b94")
}

func (RouterOffRampRemoved) Topic() common.Hash {
	return common.HexToHash("0xa823809efda3ba66c873364eec120fa0923d9fabda73bc97dd5663341e2d9bcb")
}

func (RouterOnRampSet) Topic() common.Hash {
	return common.HexToHash("0x1f7d0ec248b80e5c0dde0ee531c4fc8fdb6ce9a2b3d90f560c74acd6a7202f23")
}

func (RouterOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (RouterOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_Router *Router) Address() common.Address {
	return _Router.address
}

type RouterInterface interface {
	MAXRETBYTES(opts *bind.CallOpts) (uint16, error)

	GetArmProxy(opts *bind.CallOpts) (common.Address, error)

	GetFee(opts *bind.CallOpts, destinationChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error)

	GetOffRamps(opts *bind.CallOpts) ([]RouterOffRamp, error)

	GetOnRamp(opts *bind.CallOpts, destChainSelector uint64) (common.Address, error)

	GetSupportedTokens(opts *bind.CallOpts, chainSelector uint64) ([]common.Address, error)

	GetWrappedNative(opts *bind.CallOpts) (common.Address, error)

	IsChainSupported(opts *bind.CallOpts, chainSelector uint64) (bool, error)

	IsOffRamp(opts *bind.CallOpts, sourceChainSelector uint64, offRamp common.Address) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	ApplyRampUpdates(opts *bind.TransactOpts, onRampUpdates []RouterOnRamp, offRampRemoves []RouterOffRamp, offRampAdds []RouterOffRamp) (*types.Transaction, error)

	CcipSend(opts *bind.TransactOpts, destinationChainSelector uint64, message ClientEVM2AnyMessage) (*types.Transaction, error)

	RecoverTokens(opts *bind.TransactOpts, tokenAddress common.Address, to common.Address, amount *big.Int) (*types.Transaction, error)

	RouteMessage(opts *bind.TransactOpts, message ClientAny2EVMMessage, gasForCallExactCheck uint16, gasLimit *big.Int, receiver common.Address) (*types.Transaction, error)

	SetWrappedNative(opts *bind.TransactOpts, wrappedNative common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterMessageExecuted(opts *bind.FilterOpts) (*RouterMessageExecutedIterator, error)

	WatchMessageExecuted(opts *bind.WatchOpts, sink chan<- *RouterMessageExecuted) (event.Subscription, error)

	ParseMessageExecuted(log types.Log) (*RouterMessageExecuted, error)

	FilterOffRampAdded(opts *bind.FilterOpts, sourceChainSelector []uint64) (*RouterOffRampAddedIterator, error)

	WatchOffRampAdded(opts *bind.WatchOpts, sink chan<- *RouterOffRampAdded, sourceChainSelector []uint64) (event.Subscription, error)

	ParseOffRampAdded(log types.Log) (*RouterOffRampAdded, error)

	FilterOffRampRemoved(opts *bind.FilterOpts, sourceChainSelector []uint64) (*RouterOffRampRemovedIterator, error)

	WatchOffRampRemoved(opts *bind.WatchOpts, sink chan<- *RouterOffRampRemoved, sourceChainSelector []uint64) (event.Subscription, error)

	ParseOffRampRemoved(log types.Log) (*RouterOffRampRemoved, error)

	FilterOnRampSet(opts *bind.FilterOpts, destChainSelector []uint64) (*RouterOnRampSetIterator, error)

	WatchOnRampSet(opts *bind.WatchOpts, sink chan<- *RouterOnRampSet, destChainSelector []uint64) (event.Subscription, error)

	ParseOnRampSet(log types.Log) (*RouterOnRampSet, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RouterOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *RouterOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*RouterOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RouterOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *RouterOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*RouterOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
