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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addForwarder\",\"inputs\":[{\"name\":\"forwarder\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"clearConfig\",\"inputs\":[{\"name\":\"donId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"configVersion\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getTransmissionId\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"workflowExecutionId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"reportId\",\"type\":\"bytes2\",\"internalType\":\"bytes2\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getTransmissionInfo\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"workflowExecutionId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"reportId\",\"type\":\"bytes2\",\"internalType\":\"bytes2\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIRouter.TransmissionInfo\",\"components\":[{\"name\":\"transmissionId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"state\",\"type\":\"uint8\",\"internalType\":\"enumIRouter.TransmissionState\"},{\"name\":\"transmitter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"invalidReceiver\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"success\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"gasLimit\",\"type\":\"uint80\",\"internalType\":\"uint80\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTransmitter\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"workflowExecutionId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"reportId\",\"type\":\"bytes2\",\"internalType\":\"bytes2\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isForwarder\",\"inputs\":[{\"name\":\"forwarder\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"removeForwarder\",\"inputs\":[{\"name\":\"forwarder\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"report\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"rawReport\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"reportContext\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signatures\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"route\",\"inputs\":[{\"name\":\"transmissionId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"transmitter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"metadata\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"validatedReport\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setConfig\",\"inputs\":[{\"name\":\"donId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"configVersion\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"f\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"signers\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"typeAndVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"ConfigSet\",\"inputs\":[{\"name\":\"donId\",\"type\":\"uint32\",\"indexed\":true,\"internalType\":\"uint32\"},{\"name\":\"configVersion\",\"type\":\"uint32\",\"indexed\":true,\"internalType\":\"uint32\"},{\"name\":\"f\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"},{\"name\":\"signers\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ForwarderAdded\",\"inputs\":[{\"name\":\"forwarder\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ForwarderRemoved\",\"inputs\":[{\"name\":\"forwarder\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReportProcessed\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"workflowExecutionId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"reportId\",\"type\":\"bytes2\",\"indexed\":true,\"internalType\":\"bytes2\"},{\"name\":\"result\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AlreadyAttempted\",\"inputs\":[{\"name\":\"transmissionId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"DuplicateSigner\",\"inputs\":[{\"name\":\"signer\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ExcessSigners\",\"inputs\":[{\"name\":\"numSigners\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxSigners\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"FaultToleranceMustBePositive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientGasForRouting\",\"inputs\":[{\"name\":\"transmissionId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"InsufficientSigners\",\"inputs\":[{\"name\":\"numSigners\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minSigners\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidConfig\",\"inputs\":[{\"name\":\"configId\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidReport\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSignature\",\"inputs\":[{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"InvalidSignatureCount\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"received\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidSigner\",\"inputs\":[{\"name\":\"signer\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UnauthorizedForwarder\",\"inputs\":[]}]",
	Bin: "0x608060405234801561000f575f80fd5b5033805f816100655760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b5f80546001600160a01b0319166001600160a01b038481169190911790915581161561009457610094816100b5565b5050305f908152600360205260409020805460ff191660011790555061015d565b336001600160a01b0382160361010d5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161005c565b600180546001600160a01b0319166001600160a01b038381169182179092555f8054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6121188061016a5f395ff3fe608060405234801561000f575f80fd5b50600436106100e5575f3560e01c806379ba509711610088578063abcef55411610063578063abcef55414610354578063ee59d26c1461038c578063ef6e17a01461039f578063f2fde38b146103b2575f80fd5b806379ba5097146102575780638864b8641461025f5780638da5cb5b14610337575f80fd5b8063272cbd93116100c3578063272cbd9314610173578063354bdd66146101935780634d93172d146102315780635c41d2fe14610244575f80fd5b806311289565146100e9578063181f5a77146100fe578063233fd52d14610150575b5f80fd5b6100fc6100f73660046119e9565b6103c5565b005b61013a6040518060400160405280601781526020017f4b657973746f6e65466f7277617264657220312e302e3000000000000000000081525081565b6040516101479190611a98565b60405180910390f35b61016361015e366004611aeb565b61096d565b6040519015158152602001610147565b610186610181366004611b71565b610d24565b6040516101479190611bff565b6102236101a1366004611b71565b6040517fffffffffffffffffffffffffffffffffffffffff000000000000000000000000606085901b166020820152603481018390527fffff000000000000000000000000000000000000000000000000000000000000821660548201525f906056016040516020818303038152906040528051906020012090509392505050565b604051908152602001610147565b6100fc61023f366004611cad565b610f27565b6100fc610252366004611cad565b610fa2565b6100fc611020565b61031261026d366004611b71565b6040805160609490941b7fffffffffffffffffffffffffffffffffffffffff0000000000000000000000001660208086019190915260348501939093527fffff00000000000000000000000000000000000000000000000000000000000091909116605484015280516036818503018152605690930181528251928201929092205f9081526004909152205473ffffffffffffffffffffffffffffffffffffffff1690565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610147565b5f5473ffffffffffffffffffffffffffffffffffffffff16610312565b610163610362366004611cad565b73ffffffffffffffffffffffffffffffffffffffff165f9081526003602052604090205460ff1690565b6100fc61039a366004611cd9565b61111c565b6100fc6103ad366004611d51565b6114ee565b6100fc6103c0366004611cad565b61158d565b606d851015610400576040517fb55ac75400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f805f61044189898080601f0160208091040260200160405190810160405280939291908181526020018383808284375f920191909152506115a192505050565b67ffffffffffffffff82165f90815260026020526040812080549497509195509193509160ff16908190036104b3576040517fdf3b81ea00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff841660048201526024015b60405180910390fd5b856104bf826001611daf565b60ff1614610511576104d2816001611daf565b6040517fd6022e8e00000000000000000000000000000000000000000000000000000000815260ff9091166004820152602481018790526044016104aa565b5f8b8b604051610522929190611dc8565b60405190819003812061053b918c908c90602001611dd7565b60405160208183030381529060405280519060200120905061055b611882565b5f5b888110156107d457365f8b8b8481811061057957610579611df0565b905060200281019061058b9190611e1d565b9092509050604181146105ce5781816040517f2adfdc300000000000000000000000000000000000000000000000000000000081526004016104aa929190611ec5565b5f600186848460408181106105e5576105e5611df0565b6105f792013560f81c9050601b611daf565b61060460205f8789611ee0565b61060d91611f07565b61061b60406020888a611ee0565b61062491611f07565b604080515f8152602081018083529590955260ff909316928401929092526060830152608082015260a0016020604051602081039080840390855afa15801561066f573d5f803e3d5ffd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0015173ffffffffffffffffffffffffffffffffffffffff81165f90815260028c01602052918220549093509150819003610714576040517fbf18af4300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff831660048201526024016104aa565b5f86826020811061072757610727611df0565b602002015173ffffffffffffffffffffffffffffffffffffffff1614610791576040517fe021c4f200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff831660048201526024016104aa565b818682602081106107a4576107a4611df0565b73ffffffffffffffffffffffffffffffffffffffff909216602092909202015250506001909201915061055d9050565b50506040805160608f901b7fffffffffffffffffffffffffffffffffffffffff00000000000000000000000016602080830191909152603482018990527fffff0000000000000000000000000000000000000000000000000000000000008816605483015282516036818403018152605690920190925280519101205f945030935063233fd52d92509050338d8d8d602d90606d9261087593929190611ee0565b8f8f606d90809261088893929190611ee0565b6040518863ffffffff1660e01b81526004016108aa9796959493929190611f43565b6020604051808303815f875af11580156108c6573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906108ea9190611fb5565b9050817dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916838b73ffffffffffffffffffffffffffffffffffffffff167f3617b009e9785c42daebadb6d3fb553243a4bf586d07ea72d65d80013ce116b584604051610959911515815260200190565b60405180910390a450505050505050505050565b335f9081526003602052604081205460ff166109b5576040517fd79e123d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f6109c46113886161a8611fdb565b5a6109cf9190611fee565b90506109df6113886161a8611fdb565b6109ec9062015f90611fdb565b6109f890612710611fdb565b811015610a34576040517f0bfecd63000000000000000000000000000000000000000000000000000000008152600481018a90526024016104aa565b5f898152600460209081526040918290208251608081018452905473ffffffffffffffffffffffffffffffffffffffff8116825274010000000000000000000000000000000000000000810460ff90811615159383019390935275010000000000000000000000000000000000000000008104909216151592810183905276010000000000000000000000000000000000000000000090910469ffffffffffffffffffff1660608201529080610aeb575080602001515b15610b25576040517fa53dc8ca000000000000000000000000000000000000000000000000000000008152600481018b90526024016104aa565b5f8a8152600460205260409020805469ffffffffffffffffffff84167601000000000000000000000000000000000000000000000275ffff000000000000000000000000000000000000000090911673ffffffffffffffffffffffffffffffffffffffff8c1617179055610bb9887f805f2132000000000000000000000000000000000000000000000000000000006115bc565b610c0f575050505f87815260046020526040812080547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff1674010000000000000000000000000000000000000000179055610d19565b5f8088888888604051602401610c289493929190612001565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f805f21320000000000000000000000000000000000000000000000000000000017905290505f6113885a610caf9190611fee565b90505f808351602085015f8f86f192508215610d11575f8d815260046020526040902080547fffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffff1675010000000000000000000000000000000000000000001790555b509093505050505b979650505050505050565b6040805160c0810182525f808252602080830182905282840182905260608084018390526080840183905260a0840183905284519088901b7fffffffffffffffffffffffffffffffffffffffff0000000000000000000000001681830152603481018790527fffff000000000000000000000000000000000000000000000000000000000000861660548201528451603681830301815260568201808752815191840191909120808552600490935285842060d68301909652945473ffffffffffffffffffffffffffffffffffffffff811680875274010000000000000000000000000000000000000000820460ff9081161515607685015275010000000000000000000000000000000000000000008304161515609684015276010000000000000000000000000000000000000000000090910469ffffffffffffffffffff1660b69092019190915292939092909190610e8057505f610ea8565b816020015115610e9257506002610ea8565b8160400151610ea2576003610ea5565b60015b90505b6040518060c00160405280848152602001826003811115610ecb57610ecb611bd2565b8152602001835f015173ffffffffffffffffffffffffffffffffffffffff168152602001836020015115158152602001836040015115158152602001836060015169ffffffffffffffffffff1681525093505050509392505050565b610f2f6115e0565b73ffffffffffffffffffffffffffffffffffffffff81165f8181526003602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055517fb96d15bf9258c7b8df062753a6a262864611fc7b060a5ee2e57e79b85f898d389190a250565b610faa6115e0565b73ffffffffffffffffffffffffffffffffffffffff81165f8181526003602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055517f0ea0ce2c048ff45a4a95f2947879de3fb94abec2f152190400cab2d1272a68e79190a250565b60015473ffffffffffffffffffffffffffffffffffffffff1633146110a1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016104aa565b5f8054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6111246115e0565b8260ff165f03611160576040517f0743bae600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601f8111156111a5576040517f61750f4000000000000000000000000000000000000000000000000000000000815260048101829052601f60248201526044016104aa565b6111b0836003612027565b60ff16811161120e57806111c5846003612027565b6111d0906001611daf565b6040517f9dd9e6d8000000000000000000000000000000000000000000000000000000008152600481019290925260ff1660248201526044016104aa565b67ffffffff00000000602086901b1663ffffffff8516175f5b67ffffffffffffffff82165f908152600260205260409020600101548110156112ba5767ffffffffffffffff82165f9081526002602081905260408220600181018054919092019291908490811061128157611281611df0565b5f91825260208083209091015473ffffffffffffffffffffffffffffffffffffffff168352820192909252604001812055600101611227565b505f5b82811015611432575f8484838181106112d8576112d8611df0565b90506020020160208101906112ed9190611cad565b905073ffffffffffffffffffffffffffffffffffffffff8116611354576040517fbf18af4300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821660048201526024016104aa565b67ffffffffffffffff83165f90815260026020818152604080842073ffffffffffffffffffffffffffffffffffffffff861685529092019052902054156113df576040517fe021c4f200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821660048201526024016104aa565b6113ea826001611fdb565b67ffffffffffffffff84165f90815260026020818152604080842073ffffffffffffffffffffffffffffffffffffffff909616845294909101905291909120556001016112bd565b5067ffffffffffffffff81165f9081526002602052604090206114599060010184846118a1565b5067ffffffffffffffff81165f908152600260205260409081902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff87161790555163ffffffff86811691908816907f4120bd3b23957dd423555817d55654d4481b438aa15485c21b4180c784f1a455906114de90889088908890612043565b60405180910390a3505050505050565b6114f66115e0565b63ffffffff818116602084811b67ffffffff000000001682175f90815260028252604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690558051828152928301905291928516917f4120bd3b23957dd423555817d55654d4481b438aa15485c21b4180c784f1a455916040516115819291906120a5565b60405180910390a35050565b6115956115e0565b61159e81611662565b50565b60218101516045820151608b90920151909260c09290921c91565b5f6115c683611756565b80156115d757506115d783836117b9565b90505b92915050565b5f5473ffffffffffffffffffffffffffffffffffffffff163314611660576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016104aa565b565b3373ffffffffffffffffffffffffffffffffffffffff8216036116e1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016104aa565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8381169182179092555f8054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b5f611781827f01ffc9a7000000000000000000000000000000000000000000000000000000006117b9565b80156115da57506117b2827fffffffff000000000000000000000000000000000000000000000000000000006117b9565b1592915050565b604080517fffffffff000000000000000000000000000000000000000000000000000000008316602480830191909152825180830390910181526044909101909152602080820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f01ffc9a70000000000000000000000000000000000000000000000000000000017815282515f9392849283928392918391908a617530fa92503d91505f51905082801561186f575060208210155b8015610d19575015159695505050505050565b6040518061040001604052806020906020820280368337509192915050565b828054828255905f5260205f20908101928215611917579160200282015b828111156119175781547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8435161782556020909201916001909101906118bf565b50611923929150611927565b5090565b5b80821115611923575f8155600101611928565b803573ffffffffffffffffffffffffffffffffffffffff8116811461195e575f80fd5b919050565b5f8083601f840112611973575f80fd5b50813567ffffffffffffffff81111561198a575f80fd5b6020830191508360208285010111156119a1575f80fd5b9250929050565b5f8083601f8401126119b8575f80fd5b50813567ffffffffffffffff8111156119cf575f80fd5b6020830191508360208260051b85010111156119a1575f80fd5b5f805f805f805f6080888a0312156119ff575f80fd5b611a088861193b565b9650602088013567ffffffffffffffff811115611a23575f80fd5b611a2f8a828b01611963565b909750955050604088013567ffffffffffffffff811115611a4e575f80fd5b611a5a8a828b01611963565b909550935050606088013567ffffffffffffffff811115611a79575f80fd5b611a858a828b016119a8565b989b979a50959850939692959293505050565b602081525f82518060208401528060208501604085015e5f6040828501015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011684010191505092915050565b5f805f805f805f60a0888a031215611b01575f80fd5b87359650611b116020890161193b565b9550611b1f6040890161193b565b9450606088013567ffffffffffffffff811115611b3a575f80fd5b611b468a828b01611963565b909550935050608088013567ffffffffffffffff811115611b65575f80fd5b611a858a828b01611963565b5f805f60608486031215611b83575f80fd5b611b8c8461193b565b92506020840135915060408401357fffff00000000000000000000000000000000000000000000000000000000000081168114611bc7575f80fd5b809150509250925092565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602160045260245ffd5b81518152602082015160c082019060048110611c42577f4e487b71000000000000000000000000000000000000000000000000000000005f52602160045260245ffd5b8060208401525073ffffffffffffffffffffffffffffffffffffffff60408401511660408301526060830151151560608301526080830151611c88608084018215159052565b5060a0830151611ca660a084018269ffffffffffffffffffff169052565b5092915050565b5f60208284031215611cbd575f80fd5b6115d78261193b565b803563ffffffff8116811461195e575f80fd5b5f805f805f60808688031215611ced575f80fd5b611cf686611cc6565b9450611d0460208701611cc6565b9350604086013560ff81168114611d19575f80fd5b9250606086013567ffffffffffffffff811115611d34575f80fd5b611d40888289016119a8565b969995985093965092949392505050565b5f8060408385031215611d62575f80fd5b611d6b83611cc6565b9150611d7960208401611cc6565b90509250929050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b60ff81811683821601908111156115da576115da611d82565b818382375f9101908152919050565b838152818360208301375f910160200190815292915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b5f8083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112611e50575f80fd5b83018035915067ffffffffffffffff821115611e6a575f80fd5b6020019150368190038213156119a1575f80fd5b81835281816020850137505f602082840101525f60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b602081525f611ed8602083018486611e7e565b949350505050565b5f8085851115611eee575f80fd5b83861115611efa575f80fd5b5050820193919092039150565b803560208310156115da577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff602084900360031b1b1692915050565b87815273ffffffffffffffffffffffffffffffffffffffff8716602082015273ffffffffffffffffffffffffffffffffffffffff8616604082015260a060608201525f611f9460a083018688611e7e565b8281036080840152611fa7818587611e7e565b9a9950505050505050505050565b5f60208284031215611fc5575f80fd5b81518015158114611fd4575f80fd5b9392505050565b808201808211156115da576115da611d82565b818103818111156115da576115da611d82565b604081525f612014604083018688611e7e565b8281036020840152610d19818587611e7e565b60ff8181168382160290811690818114611ca657611ca6611d82565b60ff8416815260406020820181905281018290525f8360608301825b8581101561209a5773ffffffffffffffffffffffffffffffffffffffff6120858461193b565b1682526020928301929091019060010161205f565b509695505050505050565b5f6040820160ff85168352604060208401528084518083526060850191506020860192505f5b818110156120ff57835173ffffffffffffffffffffffffffffffffffffffff168352602093840193909201916001016120cb565b5090969550505050505056fea164736f6c634300081a000a",
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
