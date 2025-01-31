// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package forwarder

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

type IRouterTransmissionInfo struct {
	TransmissionId  [32]byte
	State           uint8
	Transmitter     common.Address
	InvalidReceiver bool
	Success         bool
	GasLimit        *big.Int
}

var KeystoneForwarderMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"transmissionId\",\"type\":\"bytes32\"}],\"name\":\"AlreadyAttempted\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"}],\"name\":\"DuplicateSigner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxSigners\",\"type\":\"uint256\"}],\"name\":\"ExcessSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FaultToleranceMustBePositive\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"transmissionId\",\"type\":\"bytes32\"}],\"name\":\"InsufficientGasForRouting\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minSigners\",\"type\":\"uint256\"}],\"name\":\"InsufficientSigners\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"configId\",\"type\":\"uint64\"}],\"name\":\"InvalidConfig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReport\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"InvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"received\",\"type\":\"uint256\"}],\"name\":\"InvalidSignatureCount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"}],\"name\":\"InvalidSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedForwarder\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"configVersion\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"}],\"name\":\"ForwarderAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"}],\"name\":\"ForwarderRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"workflowExecutionId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes2\",\"name\":\"reportId\",\"type\":\"bytes2\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"result\",\"type\":\"bool\"}],\"name\":\"ReportProcessed\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"}],\"name\":\"addForwarder\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"configVersion\",\"type\":\"uint32\"}],\"name\":\"clearConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"workflowExecutionId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes2\",\"name\":\"reportId\",\"type\":\"bytes2\"}],\"name\":\"getTransmissionId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"workflowExecutionId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes2\",\"name\":\"reportId\",\"type\":\"bytes2\"}],\"name\":\"getTransmissionInfo\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"transmissionId\",\"type\":\"bytes32\"},{\"internalType\":\"enumIRouter.TransmissionState\",\"name\":\"state\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"invalidReceiver\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"uint80\",\"name\":\"gasLimit\",\"type\":\"uint80\"}],\"internalType\":\"structIRouter.TransmissionInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"workflowExecutionId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes2\",\"name\":\"reportId\",\"type\":\"bytes2\"}],\"name\":\"getTransmitter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"}],\"name\":\"isForwarder\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"}],\"name\":\"removeForwarder\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"rawReport\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"reportContext\",\"type\":\"bytes\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"report\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"transmissionId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"validatedReport\",\"type\":\"bytes\"}],\"name\":\"route\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"configVersion\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b503380600081620000695760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156200009c576200009c81620000bf565b5050306000908152600360205260409020805460ff19166001179055506200016a565b336001600160a01b03821603620001195760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000060565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b61210c806200017a6000396000f3fe608060405234801561001057600080fd5b50600436106100ea5760003560e01c806379ba50971161008c578063abcef55411610066578063abcef5541461035d578063ee59d26c14610396578063ef6e17a0146103a9578063f2fde38b146103bc57600080fd5b806379ba50971461025e5780638864b864146102665780638da5cb5b1461033f57600080fd5b8063272cbd93116100c8578063272cbd9314610179578063354bdd66146101995780634d93172d146102385780635c41d2fe1461024b57600080fd5b806311289565146100ef578063181f5a7714610104578063233fd52d14610156575b600080fd5b6101026100fd3660046119b0565b6103cf565b005b6101406040518060400160405280601781526020017f4b657973746f6e65466f7277617264657220312e302e3000000000000000000081525081565b60405161014d9190611a5b565b60405180910390f35b610169610164366004611ac8565b610943565b604051901515815260200161014d565b61018c610187366004611b50565b610d04565b60405161014d9190611be4565b61022a6101a7366004611b50565b6040517fffffffffffffffffffffffffffffffffffffffff000000000000000000000000606085901b166020820152603481018390527fffff000000000000000000000000000000000000000000000000000000000000821660548201526000906056016040516020818303038152906040528051906020012090509392505050565b60405190815260200161014d565b610102610246366004611c8c565b610f0a565b610102610259366004611c8c565b610f86565b610102611005565b61031a610274366004611b50565b6040805160609490941b7fffffffffffffffffffffffffffffffffffffffff0000000000000000000000001660208086019190915260348501939093527fffff000000000000000000000000000000000000000000000000000000000000919091166054840152805160368185030181526056909301815282519282019290922060009081526004909152205473ffffffffffffffffffffffffffffffffffffffff1690565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161014d565b60005473ffffffffffffffffffffffffffffffffffffffff1661031a565b61016961036b366004611c8c565b73ffffffffffffffffffffffffffffffffffffffff1660009081526003602052604090205460ff1690565b6101026103a4366004611cbb565b611102565b6101026103b7366004611d39565b6114a2565b6101026103ca366004611c8c565b611542565b606d85101561040a576040517fb55ac75400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080600061044e89898080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061155692505050565b67ffffffffffffffff8216600090815260026020526040902080549396509094509092509060ff1685610482826001611d9b565b60ff16146104d957610495816001611d9b565b6040517fd6022e8e00000000000000000000000000000000000000000000000000000000815260ff9091166004820152602481018790526044015b60405180910390fd5b60008b8b6040516104eb929190611db4565b604051908190038120610504918c908c90602001611dc4565b60405160208183030381529060405280519060200120905061052461183d565b60005b888110156107a6573660008b8b8481811061054457610544611dde565b90506020028101906105569190611e0d565b9092509050604181146105995781816040517f2adfdc300000000000000000000000000000000000000000000000000000000081526004016104d0929190611ebb565b6000600186848460408181106105b1576105b1611dde565b6105c392013560f81c9050601b611d9b565b6105d1602060008789611ed7565b6105da91611f01565b6105e860406020888a611ed7565b6105f191611f01565b6040805160008152602081018083529590955260ff909316928401929092526060830152608082015260a0016020604051602081039080840390855afa15801561063f573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0015173ffffffffffffffffffffffffffffffffffffffff8116600090815260028c016020529182205490935091508190036106e5576040517fbf18af4300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff831660048201526024016104d0565b60008682602081106106f9576106f9611dde565b602002015173ffffffffffffffffffffffffffffffffffffffff1614610763576040517fe021c4f200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff831660048201526024016104d0565b8186826020811061077657610776611dde565b73ffffffffffffffffffffffffffffffffffffffff90921660209290920201525050600190920191506105279050565b50506040805160608f901b7fffffffffffffffffffffffffffffffffffffffff00000000000000000000000016602080830191909152603482018990527fffff0000000000000000000000000000000000000000000000000000000000008816605483015282516036818403018152605690920190925280519101206000945030935063233fd52d92509050338d8d8d602d90606d9261084893929190611ed7565b8f8f606d90809261085b93929190611ed7565b6040518863ffffffff1660e01b815260040161087d9796959493929190611f3d565b6020604051808303816000875af115801561089c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108c09190611f9e565b9050817dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916838b73ffffffffffffffffffffffffffffffffffffffff167f3617b009e9785c42daebadb6d3fb553243a4bf586d07ea72d65d80013ce116b58460405161092f911515815260200190565b60405180910390a450505050505050505050565b3360009081526003602052604081205460ff1661098c576040517fd79e123d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600061099c6113886161a8611fc7565b5a6109a79190611fda565b90506109b76113886161a8611fc7565b6109c49062015f90611fc7565b6109d090612710611fc7565b811015610a0c576040517f0bfecd63000000000000000000000000000000000000000000000000000000008152600481018a90526024016104d0565b6000898152600460209081526040918290208251608081018452905473ffffffffffffffffffffffffffffffffffffffff8116825274010000000000000000000000000000000000000000810460ff90811615159383019390935275010000000000000000000000000000000000000000008104909216151592810183905276010000000000000000000000000000000000000000000090910469ffffffffffffffffffff1660608201529080610ac4575080602001515b15610afe576040517fa53dc8ca000000000000000000000000000000000000000000000000000000008152600481018b90526024016104d0565b60008a8152600460205260409020805469ffffffffffffffffffff84167601000000000000000000000000000000000000000000000275ffff000000000000000000000000000000000000000090911673ffffffffffffffffffffffffffffffffffffffff8c1617179055610b93887f805f213200000000000000000000000000000000000000000000000000000000611571565b610bea57505050600087815260046020526040812080547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff1674010000000000000000000000000000000000000000179055610cf9565b60008088888888604051602401610c049493929190611fed565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f805f213200000000000000000000000000000000000000000000000000000000179052905060006113885a610c8c9190611fda565b905060008083516020850160008f86f192508215610cf15760008d815260046020526040902080547fffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffff1675010000000000000000000000000000000000000000001790555b509093505050505b979650505050505050565b6040805160c0810182526000808252602080830182905282840182905260608084018390526080840183905260a0840183905284519088901b7fffffffffffffffffffffffffffffffffffffffff0000000000000000000000001681830152603481018790527fffff000000000000000000000000000000000000000000000000000000000000861660548201528451603681830301815260568201808752815191840191909120808552600490935285842060d68301909652945473ffffffffffffffffffffffffffffffffffffffff811680875274010000000000000000000000000000000000000000820460ff9081161515607685015275010000000000000000000000000000000000000000008304161515609684015276010000000000000000000000000000000000000000000090910469ffffffffffffffffffff1660b69092019190915292939092909190610e6257506000610e8a565b816020015115610e7457506002610e8a565b8160400151610e84576003610e87565b60015b90505b6040518060c00160405280848152602001826003811115610ead57610ead611bb5565b8152602001836000015173ffffffffffffffffffffffffffffffffffffffff168152602001836020015115158152602001836040015115158152602001836060015169ffffffffffffffffffff1681525093505050509392505050565b610f12611596565b73ffffffffffffffffffffffffffffffffffffffff811660008181526003602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055517fb96d15bf9258c7b8df062753a6a262864611fc7b060a5ee2e57e79b85f898d389190a250565b610f8e611596565b73ffffffffffffffffffffffffffffffffffffffff811660008181526003602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055517f0ea0ce2c048ff45a4a95f2947879de3fb94abec2f152190400cab2d1272a68e79190a250565b60015473ffffffffffffffffffffffffffffffffffffffff163314611086576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016104d0565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b61110a611596565b601f81111561114f576040517f61750f4000000000000000000000000000000000000000000000000000000000815260048101829052601f60248201526044016104d0565b61115a836003612014565b60ff1681116111b8578061116f846003612014565b61117a906001611d9b565b6040517f9dd9e6d8000000000000000000000000000000000000000000000000000000008152600481019290925260ff1660248201526044016104d0565b67ffffffff00000000602086901b1663ffffffff85161760005b67ffffffffffffffff82166000908152600260205260409020600101548110156112685767ffffffffffffffff821660009081526002602081905260408220600181018054919092019291908490811061122e5761122e611dde565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff1683528201929092526040018120556001016111d2565b5060005b828110156113e457600084848381811061128857611288611dde565b905060200201602081019061129d9190611c8c565b905073ffffffffffffffffffffffffffffffffffffffff8116611304576040517fbf18af4300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821660048201526024016104d0565b67ffffffffffffffff8316600090815260026020818152604080842073ffffffffffffffffffffffffffffffffffffffff86168552909201905290205415611390576040517fe021c4f200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821660048201526024016104d0565b61139b826001611fc7565b67ffffffffffffffff8416600090815260026020818152604080842073ffffffffffffffffffffffffffffffffffffffff9096168452949091019052919091205560010161126c565b5067ffffffffffffffff8116600090815260026020526040902061140c90600101848461185c565b5067ffffffffffffffff81166000908152600260205260409081902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff87161790555163ffffffff86811691908816907f4120bd3b23957dd423555817d55654d4481b438aa15485c21b4180c784f1a4559061149290889088908890612030565b60405180910390a3505050505050565b6114aa611596565b63ffffffff818116602084811b67ffffffff00000000168217600090815260028252604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690558051828152928301905291928516917f4120bd3b23957dd423555817d55654d4481b438aa15485c21b4180c784f1a45591604051611536929190612096565b60405180910390a35050565b61154a611596565b61155381611619565b50565b60218101516045820151608b90920151909260c09290921c91565b600061157c8361170e565b801561158d575061158d8383611772565b90505b92915050565b60005473ffffffffffffffffffffffffffffffffffffffff163314611617576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016104d0565b565b3373ffffffffffffffffffffffffffffffffffffffff821603611698576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016104d0565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600061173a827f01ffc9a700000000000000000000000000000000000000000000000000000000611772565b8015611590575061176b827fffffffff00000000000000000000000000000000000000000000000000000000611772565b1592915050565b604080517fffffffff000000000000000000000000000000000000000000000000000000008316602480830191909152825180830390910181526044909101909152602080820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f01ffc9a700000000000000000000000000000000000000000000000000000000178152825160009392849283928392918391908a617530fa92503d9150600051905082801561182a575060208210155b8015610cf9575015159695505050505050565b6040518061040001604052806020906020820280368337509192915050565b8280548282559060005260206000209081019282156118d4579160200282015b828111156118d45781547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff84351617825560209092019160019091019061187c565b506118e09291506118e4565b5090565b5b808211156118e057600081556001016118e5565b803573ffffffffffffffffffffffffffffffffffffffff8116811461191d57600080fd5b919050565b60008083601f84011261193457600080fd5b50813567ffffffffffffffff81111561194c57600080fd5b60208301915083602082850101111561196457600080fd5b9250929050565b60008083601f84011261197d57600080fd5b50813567ffffffffffffffff81111561199557600080fd5b6020830191508360208260051b850101111561196457600080fd5b60008060008060008060006080888a0312156119cb57600080fd5b6119d4886118f9565b9650602088013567ffffffffffffffff808211156119f157600080fd5b6119fd8b838c01611922565b909850965060408a0135915080821115611a1657600080fd5b611a228b838c01611922565b909650945060608a0135915080821115611a3b57600080fd5b50611a488a828b0161196b565b989b979a50959850939692959293505050565b60006020808352835180602085015260005b81811015611a8957858101830151858201604001528201611a6d565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b600080600080600080600060a0888a031215611ae357600080fd5b87359650611af3602089016118f9565b9550611b01604089016118f9565b9450606088013567ffffffffffffffff80821115611b1e57600080fd5b611b2a8b838c01611922565b909650945060808a0135915080821115611b4357600080fd5b50611a488a828b01611922565b600080600060608486031215611b6557600080fd5b611b6e846118f9565b92506020840135915060408401357fffff00000000000000000000000000000000000000000000000000000000000081168114611baa57600080fd5b809150509250925092565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b81518152602082015160c082019060048110611c29577f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b8060208401525073ffffffffffffffffffffffffffffffffffffffff604084015116604083015260608301511515606083015260808301511515608083015260a0830151611c8560a084018269ffffffffffffffffffff169052565b5092915050565b600060208284031215611c9e57600080fd5b61158d826118f9565b803563ffffffff8116811461191d57600080fd5b600080600080600060808688031215611cd357600080fd5b611cdc86611ca7565b9450611cea60208701611ca7565b9350604086013560ff81168114611d0057600080fd5b9250606086013567ffffffffffffffff811115611d1c57600080fd5b611d288882890161196b565b969995985093965092949392505050565b60008060408385031215611d4c57600080fd5b611d5583611ca7565b9150611d6360208401611ca7565b90509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60ff818116838216019081111561159057611590611d6c565b8183823760009101908152919050565b838152818360208301376000910160200190815292915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112611e4257600080fd5b83018035915067ffffffffffffffff821115611e5d57600080fd5b60200191503681900382131561196457600080fd5b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b602081526000611ecf602083018486611e72565b949350505050565b60008085851115611ee757600080fd5b83861115611ef457600080fd5b5050820193919092039150565b80356020831015611590577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff602084900360031b1b1692915050565b878152600073ffffffffffffffffffffffffffffffffffffffff808916602084015280881660408401525060a06060830152611f7d60a083018688611e72565b8281036080840152611f90818587611e72565b9a9950505050505050505050565b600060208284031215611fb057600080fd5b81518015158114611fc057600080fd5b9392505050565b8082018082111561159057611590611d6c565b8181038181111561159057611590611d6c565b604081526000612001604083018688611e72565b8281036020840152610cf9818587611e72565b60ff8181168382160290811690818114611c8557611c85611d6c565b60ff8416815260406020808301829052908201839052600090849060608401835b8681101561208a5773ffffffffffffffffffffffffffffffffffffffff612077856118f9565b1682529282019290820190600101612051565b50979650505050505050565b60006040820160ff8516835260206040602085015281855180845260608601915060208701935060005b818110156120f257845173ffffffffffffffffffffffffffffffffffffffff16835293830193918301916001016120c0565b509097965050505050505056fea164736f6c6343000818000a",
}

var KeystoneForwarderABI = KeystoneForwarderMetaData.ABI

var KeystoneForwarderBin = KeystoneForwarderMetaData.Bin

func DeployKeystoneForwarder(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *KeystoneForwarder, error) {
	parsed, err := KeystoneForwarderMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KeystoneForwarderBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KeystoneForwarder{address: address, abi: *parsed, KeystoneForwarderCaller: KeystoneForwarderCaller{contract: contract}, KeystoneForwarderTransactor: KeystoneForwarderTransactor{contract: contract}, KeystoneForwarderFilterer: KeystoneForwarderFilterer{contract: contract}}, nil
}

type KeystoneForwarder struct {
	address common.Address
	abi     abi.ABI
	KeystoneForwarderCaller
	KeystoneForwarderTransactor
	KeystoneForwarderFilterer
}

type KeystoneForwarderCaller struct {
	contract *bind.BoundContract
}

type KeystoneForwarderTransactor struct {
	contract *bind.BoundContract
}

type KeystoneForwarderFilterer struct {
	contract *bind.BoundContract
}

type KeystoneForwarderSession struct {
	Contract     *KeystoneForwarder
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type KeystoneForwarderCallerSession struct {
	Contract *KeystoneForwarderCaller
	CallOpts bind.CallOpts
}

type KeystoneForwarderTransactorSession struct {
	Contract     *KeystoneForwarderTransactor
	TransactOpts bind.TransactOpts
}

type KeystoneForwarderRaw struct {
	Contract *KeystoneForwarder
}

type KeystoneForwarderCallerRaw struct {
	Contract *KeystoneForwarderCaller
}

type KeystoneForwarderTransactorRaw struct {
	Contract *KeystoneForwarderTransactor
}

func NewKeystoneForwarder(address common.Address, backend bind.ContractBackend) (*KeystoneForwarder, error) {
	abi, err := abi.JSON(strings.NewReader(KeystoneForwarderABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindKeystoneForwarder(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeystoneForwarder{address: address, abi: abi, KeystoneForwarderCaller: KeystoneForwarderCaller{contract: contract}, KeystoneForwarderTransactor: KeystoneForwarderTransactor{contract: contract}, KeystoneForwarderFilterer: KeystoneForwarderFilterer{contract: contract}}, nil
}

func NewKeystoneForwarderCaller(address common.Address, caller bind.ContractCaller) (*KeystoneForwarderCaller, error) {
	contract, err := bindKeystoneForwarder(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeystoneForwarderCaller{contract: contract}, nil
}

func NewKeystoneForwarderTransactor(address common.Address, transactor bind.ContractTransactor) (*KeystoneForwarderTransactor, error) {
	contract, err := bindKeystoneForwarder(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeystoneForwarderTransactor{contract: contract}, nil
}

func NewKeystoneForwarderFilterer(address common.Address, filterer bind.ContractFilterer) (*KeystoneForwarderFilterer, error) {
	contract, err := bindKeystoneForwarder(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeystoneForwarderFilterer{contract: contract}, nil
}

func bindKeystoneForwarder(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := KeystoneForwarderMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_KeystoneForwarder *KeystoneForwarderRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeystoneForwarder.Contract.KeystoneForwarderCaller.contract.Call(opts, result, method, params...)
}

func (_KeystoneForwarder *KeystoneForwarderRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.KeystoneForwarderTransactor.contract.Transfer(opts)
}

func (_KeystoneForwarder *KeystoneForwarderRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.KeystoneForwarderTransactor.contract.Transact(opts, method, params...)
}

func (_KeystoneForwarder *KeystoneForwarderCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeystoneForwarder.Contract.contract.Call(opts, result, method, params...)
}

func (_KeystoneForwarder *KeystoneForwarderTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.contract.Transfer(opts)
}

func (_KeystoneForwarder *KeystoneForwarderTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.contract.Transact(opts, method, params...)
}

func (_KeystoneForwarder *KeystoneForwarderCaller) GetTransmissionId(opts *bind.CallOpts, receiver common.Address, workflowExecutionId [32]byte, reportId [2]byte) ([32]byte, error) {
	var out []interface{}
	err := _KeystoneForwarder.contract.Call(opts, &out, "getTransmissionId", receiver, workflowExecutionId, reportId)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_KeystoneForwarder *KeystoneForwarderSession) GetTransmissionId(receiver common.Address, workflowExecutionId [32]byte, reportId [2]byte) ([32]byte, error) {
	return _KeystoneForwarder.Contract.GetTransmissionId(&_KeystoneForwarder.CallOpts, receiver, workflowExecutionId, reportId)
}

func (_KeystoneForwarder *KeystoneForwarderCallerSession) GetTransmissionId(receiver common.Address, workflowExecutionId [32]byte, reportId [2]byte) ([32]byte, error) {
	return _KeystoneForwarder.Contract.GetTransmissionId(&_KeystoneForwarder.CallOpts, receiver, workflowExecutionId, reportId)
}

func (_KeystoneForwarder *KeystoneForwarderCaller) GetTransmissionInfo(opts *bind.CallOpts, receiver common.Address, workflowExecutionId [32]byte, reportId [2]byte) (IRouterTransmissionInfo, error) {
	var out []interface{}
	err := _KeystoneForwarder.contract.Call(opts, &out, "getTransmissionInfo", receiver, workflowExecutionId, reportId)

	if err != nil {
		return *new(IRouterTransmissionInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(IRouterTransmissionInfo)).(*IRouterTransmissionInfo)

	return out0, err

}

func (_KeystoneForwarder *KeystoneForwarderSession) GetTransmissionInfo(receiver common.Address, workflowExecutionId [32]byte, reportId [2]byte) (IRouterTransmissionInfo, error) {
	return _KeystoneForwarder.Contract.GetTransmissionInfo(&_KeystoneForwarder.CallOpts, receiver, workflowExecutionId, reportId)
}

func (_KeystoneForwarder *KeystoneForwarderCallerSession) GetTransmissionInfo(receiver common.Address, workflowExecutionId [32]byte, reportId [2]byte) (IRouterTransmissionInfo, error) {
	return _KeystoneForwarder.Contract.GetTransmissionInfo(&_KeystoneForwarder.CallOpts, receiver, workflowExecutionId, reportId)
}

func (_KeystoneForwarder *KeystoneForwarderCaller) GetTransmitter(opts *bind.CallOpts, receiver common.Address, workflowExecutionId [32]byte, reportId [2]byte) (common.Address, error) {
	var out []interface{}
	err := _KeystoneForwarder.contract.Call(opts, &out, "getTransmitter", receiver, workflowExecutionId, reportId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeystoneForwarder *KeystoneForwarderSession) GetTransmitter(receiver common.Address, workflowExecutionId [32]byte, reportId [2]byte) (common.Address, error) {
	return _KeystoneForwarder.Contract.GetTransmitter(&_KeystoneForwarder.CallOpts, receiver, workflowExecutionId, reportId)
}

func (_KeystoneForwarder *KeystoneForwarderCallerSession) GetTransmitter(receiver common.Address, workflowExecutionId [32]byte, reportId [2]byte) (common.Address, error) {
	return _KeystoneForwarder.Contract.GetTransmitter(&_KeystoneForwarder.CallOpts, receiver, workflowExecutionId, reportId)
}

func (_KeystoneForwarder *KeystoneForwarderCaller) IsForwarder(opts *bind.CallOpts, forwarder common.Address) (bool, error) {
	var out []interface{}
	err := _KeystoneForwarder.contract.Call(opts, &out, "isForwarder", forwarder)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_KeystoneForwarder *KeystoneForwarderSession) IsForwarder(forwarder common.Address) (bool, error) {
	return _KeystoneForwarder.Contract.IsForwarder(&_KeystoneForwarder.CallOpts, forwarder)
}

func (_KeystoneForwarder *KeystoneForwarderCallerSession) IsForwarder(forwarder common.Address) (bool, error) {
	return _KeystoneForwarder.Contract.IsForwarder(&_KeystoneForwarder.CallOpts, forwarder)
}

func (_KeystoneForwarder *KeystoneForwarderCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeystoneForwarder.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeystoneForwarder *KeystoneForwarderSession) Owner() (common.Address, error) {
	return _KeystoneForwarder.Contract.Owner(&_KeystoneForwarder.CallOpts)
}

func (_KeystoneForwarder *KeystoneForwarderCallerSession) Owner() (common.Address, error) {
	return _KeystoneForwarder.Contract.Owner(&_KeystoneForwarder.CallOpts)
}

func (_KeystoneForwarder *KeystoneForwarderCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _KeystoneForwarder.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_KeystoneForwarder *KeystoneForwarderSession) TypeAndVersion() (string, error) {
	return _KeystoneForwarder.Contract.TypeAndVersion(&_KeystoneForwarder.CallOpts)
}

func (_KeystoneForwarder *KeystoneForwarderCallerSession) TypeAndVersion() (string, error) {
	return _KeystoneForwarder.Contract.TypeAndVersion(&_KeystoneForwarder.CallOpts)
}

func (_KeystoneForwarder *KeystoneForwarderTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeystoneForwarder.contract.Transact(opts, "acceptOwnership")
}

func (_KeystoneForwarder *KeystoneForwarderSession) AcceptOwnership() (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.AcceptOwnership(&_KeystoneForwarder.TransactOpts)
}

func (_KeystoneForwarder *KeystoneForwarderTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.AcceptOwnership(&_KeystoneForwarder.TransactOpts)
}

func (_KeystoneForwarder *KeystoneForwarderTransactor) AddForwarder(opts *bind.TransactOpts, forwarder common.Address) (*types.Transaction, error) {
	return _KeystoneForwarder.contract.Transact(opts, "addForwarder", forwarder)
}

func (_KeystoneForwarder *KeystoneForwarderSession) AddForwarder(forwarder common.Address) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.AddForwarder(&_KeystoneForwarder.TransactOpts, forwarder)
}

func (_KeystoneForwarder *KeystoneForwarderTransactorSession) AddForwarder(forwarder common.Address) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.AddForwarder(&_KeystoneForwarder.TransactOpts, forwarder)
}

func (_KeystoneForwarder *KeystoneForwarderTransactor) ClearConfig(opts *bind.TransactOpts, donId uint32, configVersion uint32) (*types.Transaction, error) {
	return _KeystoneForwarder.contract.Transact(opts, "clearConfig", donId, configVersion)
}

func (_KeystoneForwarder *KeystoneForwarderSession) ClearConfig(donId uint32, configVersion uint32) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.ClearConfig(&_KeystoneForwarder.TransactOpts, donId, configVersion)
}

func (_KeystoneForwarder *KeystoneForwarderTransactorSession) ClearConfig(donId uint32, configVersion uint32) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.ClearConfig(&_KeystoneForwarder.TransactOpts, donId, configVersion)
}

func (_KeystoneForwarder *KeystoneForwarderTransactor) RemoveForwarder(opts *bind.TransactOpts, forwarder common.Address) (*types.Transaction, error) {
	return _KeystoneForwarder.contract.Transact(opts, "removeForwarder", forwarder)
}

func (_KeystoneForwarder *KeystoneForwarderSession) RemoveForwarder(forwarder common.Address) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.RemoveForwarder(&_KeystoneForwarder.TransactOpts, forwarder)
}

func (_KeystoneForwarder *KeystoneForwarderTransactorSession) RemoveForwarder(forwarder common.Address) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.RemoveForwarder(&_KeystoneForwarder.TransactOpts, forwarder)
}

func (_KeystoneForwarder *KeystoneForwarderTransactor) Report(opts *bind.TransactOpts, receiver common.Address, rawReport []byte, reportContext []byte, signatures [][]byte) (*types.Transaction, error) {
	return _KeystoneForwarder.contract.Transact(opts, "report", receiver, rawReport, reportContext, signatures)
}

func (_KeystoneForwarder *KeystoneForwarderSession) Report(receiver common.Address, rawReport []byte, reportContext []byte, signatures [][]byte) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.Report(&_KeystoneForwarder.TransactOpts, receiver, rawReport, reportContext, signatures)
}

func (_KeystoneForwarder *KeystoneForwarderTransactorSession) Report(receiver common.Address, rawReport []byte, reportContext []byte, signatures [][]byte) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.Report(&_KeystoneForwarder.TransactOpts, receiver, rawReport, reportContext, signatures)
}

func (_KeystoneForwarder *KeystoneForwarderTransactor) Route(opts *bind.TransactOpts, transmissionId [32]byte, transmitter common.Address, receiver common.Address, metadata []byte, validatedReport []byte) (*types.Transaction, error) {
	return _KeystoneForwarder.contract.Transact(opts, "route", transmissionId, transmitter, receiver, metadata, validatedReport)
}

func (_KeystoneForwarder *KeystoneForwarderSession) Route(transmissionId [32]byte, transmitter common.Address, receiver common.Address, metadata []byte, validatedReport []byte) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.Route(&_KeystoneForwarder.TransactOpts, transmissionId, transmitter, receiver, metadata, validatedReport)
}

func (_KeystoneForwarder *KeystoneForwarderTransactorSession) Route(transmissionId [32]byte, transmitter common.Address, receiver common.Address, metadata []byte, validatedReport []byte) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.Route(&_KeystoneForwarder.TransactOpts, transmissionId, transmitter, receiver, metadata, validatedReport)
}

func (_KeystoneForwarder *KeystoneForwarderTransactor) SetConfig(opts *bind.TransactOpts, donId uint32, configVersion uint32, f uint8, signers []common.Address) (*types.Transaction, error) {
	return _KeystoneForwarder.contract.Transact(opts, "setConfig", donId, configVersion, f, signers)
}

func (_KeystoneForwarder *KeystoneForwarderSession) SetConfig(donId uint32, configVersion uint32, f uint8, signers []common.Address) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.SetConfig(&_KeystoneForwarder.TransactOpts, donId, configVersion, f, signers)
}

func (_KeystoneForwarder *KeystoneForwarderTransactorSession) SetConfig(donId uint32, configVersion uint32, f uint8, signers []common.Address) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.SetConfig(&_KeystoneForwarder.TransactOpts, donId, configVersion, f, signers)
}

func (_KeystoneForwarder *KeystoneForwarderTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _KeystoneForwarder.contract.Transact(opts, "transferOwnership", to)
}

func (_KeystoneForwarder *KeystoneForwarderSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.TransferOwnership(&_KeystoneForwarder.TransactOpts, to)
}

func (_KeystoneForwarder *KeystoneForwarderTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.TransferOwnership(&_KeystoneForwarder.TransactOpts, to)
}

type KeystoneForwarderConfigSetIterator struct {
	Event *KeystoneForwarderConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeystoneForwarderConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeystoneForwarderConfigSet)
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
		it.Event = new(KeystoneForwarderConfigSet)
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

func (it *KeystoneForwarderConfigSetIterator) Error() error {
	return it.fail
}

func (it *KeystoneForwarderConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeystoneForwarderConfigSet struct {
	DonId         uint32
	ConfigVersion uint32
	F             uint8
	Signers       []common.Address
	Raw           types.Log
}

func (_KeystoneForwarder *KeystoneForwarderFilterer) FilterConfigSet(opts *bind.FilterOpts, donId []uint32, configVersion []uint32) (*KeystoneForwarderConfigSetIterator, error) {

	var donIdRule []interface{}
	for _, donIdItem := range donId {
		donIdRule = append(donIdRule, donIdItem)
	}
	var configVersionRule []interface{}
	for _, configVersionItem := range configVersion {
		configVersionRule = append(configVersionRule, configVersionItem)
	}

	logs, sub, err := _KeystoneForwarder.contract.FilterLogs(opts, "ConfigSet", donIdRule, configVersionRule)
	if err != nil {
		return nil, err
	}
	return &KeystoneForwarderConfigSetIterator{contract: _KeystoneForwarder.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_KeystoneForwarder *KeystoneForwarderFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *KeystoneForwarderConfigSet, donId []uint32, configVersion []uint32) (event.Subscription, error) {

	var donIdRule []interface{}
	for _, donIdItem := range donId {
		donIdRule = append(donIdRule, donIdItem)
	}
	var configVersionRule []interface{}
	for _, configVersionItem := range configVersion {
		configVersionRule = append(configVersionRule, configVersionItem)
	}

	logs, sub, err := _KeystoneForwarder.contract.WatchLogs(opts, "ConfigSet", donIdRule, configVersionRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeystoneForwarderConfigSet)
				if err := _KeystoneForwarder.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_KeystoneForwarder *KeystoneForwarderFilterer) ParseConfigSet(log types.Log) (*KeystoneForwarderConfigSet, error) {
	event := new(KeystoneForwarderConfigSet)
	if err := _KeystoneForwarder.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeystoneForwarderForwarderAddedIterator struct {
	Event *KeystoneForwarderForwarderAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeystoneForwarderForwarderAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeystoneForwarderForwarderAdded)
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
		it.Event = new(KeystoneForwarderForwarderAdded)
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

func (it *KeystoneForwarderForwarderAddedIterator) Error() error {
	return it.fail
}

func (it *KeystoneForwarderForwarderAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeystoneForwarderForwarderAdded struct {
	Forwarder common.Address
	Raw       types.Log
}

func (_KeystoneForwarder *KeystoneForwarderFilterer) FilterForwarderAdded(opts *bind.FilterOpts, forwarder []common.Address) (*KeystoneForwarderForwarderAddedIterator, error) {

	var forwarderRule []interface{}
	for _, forwarderItem := range forwarder {
		forwarderRule = append(forwarderRule, forwarderItem)
	}

	logs, sub, err := _KeystoneForwarder.contract.FilterLogs(opts, "ForwarderAdded", forwarderRule)
	if err != nil {
		return nil, err
	}
	return &KeystoneForwarderForwarderAddedIterator{contract: _KeystoneForwarder.contract, event: "ForwarderAdded", logs: logs, sub: sub}, nil
}

func (_KeystoneForwarder *KeystoneForwarderFilterer) WatchForwarderAdded(opts *bind.WatchOpts, sink chan<- *KeystoneForwarderForwarderAdded, forwarder []common.Address) (event.Subscription, error) {

	var forwarderRule []interface{}
	for _, forwarderItem := range forwarder {
		forwarderRule = append(forwarderRule, forwarderItem)
	}

	logs, sub, err := _KeystoneForwarder.contract.WatchLogs(opts, "ForwarderAdded", forwarderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeystoneForwarderForwarderAdded)
				if err := _KeystoneForwarder.contract.UnpackLog(event, "ForwarderAdded", log); err != nil {
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

func (_KeystoneForwarder *KeystoneForwarderFilterer) ParseForwarderAdded(log types.Log) (*KeystoneForwarderForwarderAdded, error) {
	event := new(KeystoneForwarderForwarderAdded)
	if err := _KeystoneForwarder.contract.UnpackLog(event, "ForwarderAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeystoneForwarderForwarderRemovedIterator struct {
	Event *KeystoneForwarderForwarderRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeystoneForwarderForwarderRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeystoneForwarderForwarderRemoved)
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
		it.Event = new(KeystoneForwarderForwarderRemoved)
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

func (it *KeystoneForwarderForwarderRemovedIterator) Error() error {
	return it.fail
}

func (it *KeystoneForwarderForwarderRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeystoneForwarderForwarderRemoved struct {
	Forwarder common.Address
	Raw       types.Log
}

func (_KeystoneForwarder *KeystoneForwarderFilterer) FilterForwarderRemoved(opts *bind.FilterOpts, forwarder []common.Address) (*KeystoneForwarderForwarderRemovedIterator, error) {

	var forwarderRule []interface{}
	for _, forwarderItem := range forwarder {
		forwarderRule = append(forwarderRule, forwarderItem)
	}

	logs, sub, err := _KeystoneForwarder.contract.FilterLogs(opts, "ForwarderRemoved", forwarderRule)
	if err != nil {
		return nil, err
	}
	return &KeystoneForwarderForwarderRemovedIterator{contract: _KeystoneForwarder.contract, event: "ForwarderRemoved", logs: logs, sub: sub}, nil
}

func (_KeystoneForwarder *KeystoneForwarderFilterer) WatchForwarderRemoved(opts *bind.WatchOpts, sink chan<- *KeystoneForwarderForwarderRemoved, forwarder []common.Address) (event.Subscription, error) {

	var forwarderRule []interface{}
	for _, forwarderItem := range forwarder {
		forwarderRule = append(forwarderRule, forwarderItem)
	}

	logs, sub, err := _KeystoneForwarder.contract.WatchLogs(opts, "ForwarderRemoved", forwarderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeystoneForwarderForwarderRemoved)
				if err := _KeystoneForwarder.contract.UnpackLog(event, "ForwarderRemoved", log); err != nil {
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

func (_KeystoneForwarder *KeystoneForwarderFilterer) ParseForwarderRemoved(log types.Log) (*KeystoneForwarderForwarderRemoved, error) {
	event := new(KeystoneForwarderForwarderRemoved)
	if err := _KeystoneForwarder.contract.UnpackLog(event, "ForwarderRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeystoneForwarderOwnershipTransferRequestedIterator struct {
	Event *KeystoneForwarderOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeystoneForwarderOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeystoneForwarderOwnershipTransferRequested)
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
		it.Event = new(KeystoneForwarderOwnershipTransferRequested)
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

func (it *KeystoneForwarderOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeystoneForwarderOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeystoneForwarderOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeystoneForwarder *KeystoneForwarderFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeystoneForwarderOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeystoneForwarder.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeystoneForwarderOwnershipTransferRequestedIterator{contract: _KeystoneForwarder.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeystoneForwarder *KeystoneForwarderFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeystoneForwarderOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeystoneForwarder.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeystoneForwarderOwnershipTransferRequested)
				if err := _KeystoneForwarder.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_KeystoneForwarder *KeystoneForwarderFilterer) ParseOwnershipTransferRequested(log types.Log) (*KeystoneForwarderOwnershipTransferRequested, error) {
	event := new(KeystoneForwarderOwnershipTransferRequested)
	if err := _KeystoneForwarder.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeystoneForwarderOwnershipTransferredIterator struct {
	Event *KeystoneForwarderOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeystoneForwarderOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeystoneForwarderOwnershipTransferred)
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
		it.Event = new(KeystoneForwarderOwnershipTransferred)
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

func (it *KeystoneForwarderOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *KeystoneForwarderOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeystoneForwarderOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeystoneForwarder *KeystoneForwarderFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeystoneForwarderOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeystoneForwarder.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeystoneForwarderOwnershipTransferredIterator{contract: _KeystoneForwarder.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_KeystoneForwarder *KeystoneForwarderFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeystoneForwarderOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeystoneForwarder.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeystoneForwarderOwnershipTransferred)
				if err := _KeystoneForwarder.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_KeystoneForwarder *KeystoneForwarderFilterer) ParseOwnershipTransferred(log types.Log) (*KeystoneForwarderOwnershipTransferred, error) {
	event := new(KeystoneForwarderOwnershipTransferred)
	if err := _KeystoneForwarder.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeystoneForwarderReportProcessedIterator struct {
	Event *KeystoneForwarderReportProcessed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeystoneForwarderReportProcessedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeystoneForwarderReportProcessed)
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
		it.Event = new(KeystoneForwarderReportProcessed)
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

func (it *KeystoneForwarderReportProcessedIterator) Error() error {
	return it.fail
}

func (it *KeystoneForwarderReportProcessedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeystoneForwarderReportProcessed struct {
	Receiver            common.Address
	WorkflowExecutionId [32]byte
	ReportId            [2]byte
	Result              bool
	Raw                 types.Log
}

func (_KeystoneForwarder *KeystoneForwarderFilterer) FilterReportProcessed(opts *bind.FilterOpts, receiver []common.Address, workflowExecutionId [][32]byte, reportId [][2]byte) (*KeystoneForwarderReportProcessedIterator, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}
	var workflowExecutionIdRule []interface{}
	for _, workflowExecutionIdItem := range workflowExecutionId {
		workflowExecutionIdRule = append(workflowExecutionIdRule, workflowExecutionIdItem)
	}
	var reportIdRule []interface{}
	for _, reportIdItem := range reportId {
		reportIdRule = append(reportIdRule, reportIdItem)
	}

	logs, sub, err := _KeystoneForwarder.contract.FilterLogs(opts, "ReportProcessed", receiverRule, workflowExecutionIdRule, reportIdRule)
	if err != nil {
		return nil, err
	}
	return &KeystoneForwarderReportProcessedIterator{contract: _KeystoneForwarder.contract, event: "ReportProcessed", logs: logs, sub: sub}, nil
}

func (_KeystoneForwarder *KeystoneForwarderFilterer) WatchReportProcessed(opts *bind.WatchOpts, sink chan<- *KeystoneForwarderReportProcessed, receiver []common.Address, workflowExecutionId [][32]byte, reportId [][2]byte) (event.Subscription, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}
	var workflowExecutionIdRule []interface{}
	for _, workflowExecutionIdItem := range workflowExecutionId {
		workflowExecutionIdRule = append(workflowExecutionIdRule, workflowExecutionIdItem)
	}
	var reportIdRule []interface{}
	for _, reportIdItem := range reportId {
		reportIdRule = append(reportIdRule, reportIdItem)
	}

	logs, sub, err := _KeystoneForwarder.contract.WatchLogs(opts, "ReportProcessed", receiverRule, workflowExecutionIdRule, reportIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeystoneForwarderReportProcessed)
				if err := _KeystoneForwarder.contract.UnpackLog(event, "ReportProcessed", log); err != nil {
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

func (_KeystoneForwarder *KeystoneForwarderFilterer) ParseReportProcessed(log types.Log) (*KeystoneForwarderReportProcessed, error) {
	event := new(KeystoneForwarderReportProcessed)
	if err := _KeystoneForwarder.contract.UnpackLog(event, "ReportProcessed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_KeystoneForwarder *KeystoneForwarder) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _KeystoneForwarder.abi.Events["ConfigSet"].ID:
		return _KeystoneForwarder.ParseConfigSet(log)
	case _KeystoneForwarder.abi.Events["ForwarderAdded"].ID:
		return _KeystoneForwarder.ParseForwarderAdded(log)
	case _KeystoneForwarder.abi.Events["ForwarderRemoved"].ID:
		return _KeystoneForwarder.ParseForwarderRemoved(log)
	case _KeystoneForwarder.abi.Events["OwnershipTransferRequested"].ID:
		return _KeystoneForwarder.ParseOwnershipTransferRequested(log)
	case _KeystoneForwarder.abi.Events["OwnershipTransferred"].ID:
		return _KeystoneForwarder.ParseOwnershipTransferred(log)
	case _KeystoneForwarder.abi.Events["ReportProcessed"].ID:
		return _KeystoneForwarder.ParseReportProcessed(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (KeystoneForwarderConfigSet) Topic() common.Hash {
	return common.HexToHash("0x4120bd3b23957dd423555817d55654d4481b438aa15485c21b4180c784f1a455")
}

func (KeystoneForwarderForwarderAdded) Topic() common.Hash {
	return common.HexToHash("0x0ea0ce2c048ff45a4a95f2947879de3fb94abec2f152190400cab2d1272a68e7")
}

func (KeystoneForwarderForwarderRemoved) Topic() common.Hash {
	return common.HexToHash("0xb96d15bf9258c7b8df062753a6a262864611fc7b060a5ee2e57e79b85f898d38")
}

func (KeystoneForwarderOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (KeystoneForwarderOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (KeystoneForwarderReportProcessed) Topic() common.Hash {
	return common.HexToHash("0x3617b009e9785c42daebadb6d3fb553243a4bf586d07ea72d65d80013ce116b5")
}

func (_KeystoneForwarder *KeystoneForwarder) Address() common.Address {
	return _KeystoneForwarder.address
}

type KeystoneForwarderInterface interface {
	GetTransmissionId(opts *bind.CallOpts, receiver common.Address, workflowExecutionId [32]byte, reportId [2]byte) ([32]byte, error)

	GetTransmissionInfo(opts *bind.CallOpts, receiver common.Address, workflowExecutionId [32]byte, reportId [2]byte) (IRouterTransmissionInfo, error)

	GetTransmitter(opts *bind.CallOpts, receiver common.Address, workflowExecutionId [32]byte, reportId [2]byte) (common.Address, error)

	IsForwarder(opts *bind.CallOpts, forwarder common.Address) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AddForwarder(opts *bind.TransactOpts, forwarder common.Address) (*types.Transaction, error)

	ClearConfig(opts *bind.TransactOpts, donId uint32, configVersion uint32) (*types.Transaction, error)

	RemoveForwarder(opts *bind.TransactOpts, forwarder common.Address) (*types.Transaction, error)

	Report(opts *bind.TransactOpts, receiver common.Address, rawReport []byte, reportContext []byte, signatures [][]byte) (*types.Transaction, error)

	Route(opts *bind.TransactOpts, transmissionId [32]byte, transmitter common.Address, receiver common.Address, metadata []byte, validatedReport []byte) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, donId uint32, configVersion uint32, f uint8, signers []common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterConfigSet(opts *bind.FilterOpts, donId []uint32, configVersion []uint32) (*KeystoneForwarderConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *KeystoneForwarderConfigSet, donId []uint32, configVersion []uint32) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*KeystoneForwarderConfigSet, error)

	FilterForwarderAdded(opts *bind.FilterOpts, forwarder []common.Address) (*KeystoneForwarderForwarderAddedIterator, error)

	WatchForwarderAdded(opts *bind.WatchOpts, sink chan<- *KeystoneForwarderForwarderAdded, forwarder []common.Address) (event.Subscription, error)

	ParseForwarderAdded(log types.Log) (*KeystoneForwarderForwarderAdded, error)

	FilterForwarderRemoved(opts *bind.FilterOpts, forwarder []common.Address) (*KeystoneForwarderForwarderRemovedIterator, error)

	WatchForwarderRemoved(opts *bind.WatchOpts, sink chan<- *KeystoneForwarderForwarderRemoved, forwarder []common.Address) (event.Subscription, error)

	ParseForwarderRemoved(log types.Log) (*KeystoneForwarderForwarderRemoved, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeystoneForwarderOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeystoneForwarderOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*KeystoneForwarderOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeystoneForwarderOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeystoneForwarderOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*KeystoneForwarderOwnershipTransferred, error)

	FilterReportProcessed(opts *bind.FilterOpts, receiver []common.Address, workflowExecutionId [][32]byte, reportId [][2]byte) (*KeystoneForwarderReportProcessedIterator, error)

	WatchReportProcessed(opts *bind.WatchOpts, sink chan<- *KeystoneForwarderReportProcessed, receiver []common.Address, workflowExecutionId [][32]byte, reportId [][2]byte) (event.Subscription, error)

	ParseReportProcessed(log types.Log) (*KeystoneForwarderReportProcessed, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
