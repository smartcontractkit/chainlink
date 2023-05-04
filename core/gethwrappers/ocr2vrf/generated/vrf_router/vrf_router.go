// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_router

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

type VRFBeaconTypesOutputServed struct {
	Height            uint64
	ConfirmationDelay *big.Int
	ProofG1X          *big.Int
	ProofG1Y          *big.Int
}

var VRFRouterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"CoordinatorAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CoordinatorNotRegistered\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"route\",\"type\":\"address\"}],\"name\":\"RouteNotFound\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnexpectedMigrationVersion\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorDeregistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"aggregatorRoundId\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"uint40\",\"name\":\"epochAndRound\",\"type\":\"uint40\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint192\",\"name\":\"juelsPerFeeCoin\",\"type\":\"uint192\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"reasonableGasPrice\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"NewTransmission\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"recentBlockHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint192\",\"name\":\"juelsPerFeeCoin\",\"type\":\"uint192\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"reasonableGasPrice\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelay\",\"type\":\"uint24\"},{\"internalType\":\"uint256\",\"name\":\"proofG1X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"proofG1Y\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structVRFBeaconTypes.OutputServed[]\",\"name\":\"outputsServed\",\"type\":\"tuple[]\"}],\"name\":\"OutputsServed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"VRFBeaconTypes.RequestID[]\",\"name\":\"requestIDs\",\"type\":\"uint48[]\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"successfulFulfillment\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes[]\",\"name\":\"truncatedErrorData\",\"type\":\"bytes[]\"}],\"name\":\"RandomWordsFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"nextBeaconOutputHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confDelay\",\"type\":\"uint24\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasAllowance\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"weiPerUnitLink\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"}],\"name\":\"RandomnessFulfillmentRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"nextBeaconOutputHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confDelay\",\"type\":\"uint24\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"}],\"name\":\"RandomnessRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"RouteSet\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"NUM_CONF_DELAYS\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"gasAmount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"callWithExactGasEvenIfTargetIsNoContract\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"sufficientGas\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"deregisterCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCoordinators\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"}],\"name\":\"getRoute\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"isCoordinatorRegistered\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"externalRequestID\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"name\":\"redeemRandomness\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"randomness\",\"type\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"registerCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confDelay\",\"type\":\"uint24\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"name\":\"requestRandomness\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confDelay\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"name\":\"requestRandomnessFulfillment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"}],\"name\":\"resetRoute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_redemptionRoutes\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"}],\"name\":\"setRoute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5033806000816100675760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615610097576100978161009f565b505050610148565b336001600160a01b038216036100f75760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161005e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b611282806101576000396000f3fe608060405234801561001057600080fd5b50600436106100eb5760003560e01c80638da5cb5b116100925780638da5cb5b146102095780639479b74e14610223578063acfc6cdd14610236578063bb9b2b3414610256578063cef40aa614610269578063d9c4a44d1461028c578063db972c8b1461029f578063f2fde38b146102b2578063fa66358a146102c557600080fd5b8063181f5a77146100f05780632f7527cc1461012a5780634b2407d4146101445780634ffac83a1461015957806355fe97631461017a5780635d5d8d19146101a45780637612e884146101b957806379ba509714610201575b600080fd5b604080518082018252600f81526e0565246526f7574657220312e302e3608c1b602082015290516101219190610cbf565b60405180910390f35b610132600881565b60405160ff9091168152602001610121565b610157610152366004610cee565b6102d8565b005b61016c610167366004610de3565b6103da565b604051908152602001610121565b61018d610188366004610e4a565b6104f5565b604080519215158352901515602083015201610121565b6101ac61055e565b6040516101219190610ea0565b6101f36101c7366004610eed565b60056020526000908152604090205465ffffffffffff811690600160301b90046001600160a01b031682565b604051610121929190610f06565b61015761056f565b6000546001600160a01b03165b6040516101219190610f27565b610157610231366004610eed565b61061e565b610249610244366004610f3b565b6106a6565b6040516101219190610f74565b610157610264366004610eed565b61077d565b61027c610277366004610cee565b6107f2565b6040519015158152602001610121565b61015761029a366004610cee565b610805565b61016c6102ad366004610fac565b610871565b6101576102c0366004610cee565b610911565b6102166102d3366004610eed565b610925565b6102e0610987565b6102eb6003826109dc565b156103095760405163dcecb7bf60e01b815260040160405180910390fd5b6000819050806001600160a01b031663294daa496040518163ffffffff1660e01b8152600401602060405180830381865afa15801561034c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103709190611054565b60ff1660000361039357604051630dd02d0960e31b815260040160405180910390fd5b61039e600383610a01565b507fb7cabbfc11e66731fc77de0444614282023bcbd41d16781c753a431d0af01625826040516103ce9190610f27565b60405180910390a15050565b6000806103e686610925565b90506000816001600160a01b03166362f8b62033898989896040518663ffffffff1660e01b815260040161041e959493929190611077565b6020604051808303816000875af115801561043d573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061046191906110be565b905060008183604051602001610478929190610f06565b60408051601f198184030181528282528051602091820120838301835265ffffffffffff95861684526001600160a01b039687168483019081526000828152600590935292909120925183549251909616600160301b026001600160d01b0319909216959094169490941793909317909255915050949350505050565b600080336105046003826109dc565b610521576040516301fd70a160e51b815260040160405180910390fd5b5a611388811061055457611388810390508660408204820311156105545760008086516020880160008a8cf19350600192505b5050935093915050565b606061056a6003610a16565b905090565b6001546001600160a01b031633146105c75760405162461bcd60e51b815260206004820152601660248201527526bab9ba10313290383937b837b9b2b21037bbb732b960511b60448201526064015b60405180910390fd5b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b3361062a6003826109dc565b610647576040516301fd70a160e51b815260040160405180910390fd5b6000828152600260205260409081902080546001600160a01b03191633908117909155905183917fd8bb190f51a471c0cc88b5df464ee10c29138cce0d70fb36c472d60b414f3b3a9161069a9190610f27565b60405180910390a25050565b600082815260056020818152604080842081518083018352815465ffffffffffff811682526001600160a01b03600160301b82048116838701908152978a9052959094526001600160d01b0319909316905592518151935163abbf1c9b60e01b8152606094929391929183169163abbf1c9b9161072c9133918b91908a906004016110e6565b6000604051808303816000875af115801561074b573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610773919081019061111b565b9695505050505050565b336107896003826109dc565b6107a6576040516301fd70a160e51b815260040160405180910390fd5b60008281526002602052604080822080546001600160a01b03191690555183917fd8bb190f51a471c0cc88b5df464ee10c29138cce0d70fb36c472d60b414f3b3a9161069a9190610f27565b60006107ff6003836109dc565b92915050565b61080d610987565b806108196003826109dc565b610836576040516301fd70a160e51b815260040160405180910390fd5b610841600383610a23565b507ff80a1a97fd42251f3c33cda98635e7399253033a6774fe37cd3f650b5282af37826040516103ce9190610f27565b60008061087d88610925565b90506000816001600160a01b03166395009f08338b8b8b8b8b8b6040518863ffffffff1660e01b81526004016108b997969594939291906111c0565b6020604051808303816000875af11580156108d8573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108fc91906110be565b65ffffffffffff169998505050505050505050565b610919610987565b61092281610a38565b50565b6000818152600260205260408120546001600160a01b03168061095d57806040516339910a6760e01b81526004016105be9190610f27565b6109686003826109dc565b6107ff57806040516339910a6760e01b81526004016105be9190610f27565b6000546001600160a01b031633146109da5760405162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b60448201526064016105be565b565b6001600160a01b038116600090815260018301602052604081205415155b9392505050565b60006109fa836001600160a01b038416610adb565b606060006109fa83610b2a565b60006109fa836001600160a01b038416610b86565b336001600160a01b03821603610a8a5760405162461bcd60e51b815260206004820152601760248201527621b0b73737ba103a3930b739b332b9103a379039b2b63360491b60448201526064016105be565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000818152600183016020526040812054610b22575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556107ff565b5060006107ff565b606081600001805480602002602001604051908101604052809291908181526020018280548015610b7a57602002820191906000526020600020905b815481526020019060010190808311610b66575b50505050509050919050565b60008181526001830160205260408120548015610c6f576000610baa600183611228565b8554909150600090610bbe90600190611228565b9050818114610c23576000866000018281548110610bde57610bde611249565b9060005260206000200154905080876000018481548110610c0157610c01611249565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080610c3457610c3461125f565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506107ff565b60009150506107ff565b6000815180845260005b81811015610c9f57602081850181015186830182015201610c83565b506000602082860101526020601f19601f83011685010191505092915050565b6020815260006109fa6020830184610c79565b80356001600160a01b0381168114610ce957600080fd5b919050565b600060208284031215610d0057600080fd5b6109fa82610cd2565b803561ffff81168114610ce957600080fd5b803562ffffff81168114610ce957600080fd5b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f191681016001600160401b0381118282101715610d6c57610d6c610d2e565b604052919050565b600082601f830112610d8557600080fd5b81356001600160401b03811115610d9e57610d9e610d2e565b610db1601f8201601f1916602001610d44565b818152846020838601011115610dc657600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008060808587031215610df957600080fd5b84359350610e0960208601610d09565b9250610e1760408601610d1b565b915060608501356001600160401b03811115610e3257600080fd5b610e3e87828801610d74565b91505092959194509250565b600080600060608486031215610e5f57600080fd5b83359250610e6f60208501610cd2565b915060408401356001600160401b03811115610e8a57600080fd5b610e9686828701610d74565b9150509250925092565b6020808252825182820181905260009190848201906040850190845b81811015610ee15783516001600160a01b031683529284019291840191600101610ebc565b50909695505050505050565b600060208284031215610eff57600080fd5b5035919050565b65ffffffffffff9290921682526001600160a01b0316602082015260400190565b6001600160a01b0391909116815260200190565b600080600060608486031215610f5057600080fd5b833592506020840135915060408401356001600160401b03811115610e8a57600080fd5b6020808252825182820181905260009190848201906040850190845b81811015610ee157835183529284019291840191600101610f90565b60008060008060008060c08789031215610fc557600080fd5b86359550610fd560208801610d09565b9450610fe360408801610d1b565b9350606087013563ffffffff81168114610ffc57600080fd5b925060808701356001600160401b038082111561101857600080fd5b6110248a838b01610d74565b935060a089013591508082111561103a57600080fd5b5061104789828a01610d74565b9150509295509295509295565b60006020828403121561106657600080fd5b815160ff811681146109fa57600080fd5b60018060a01b038616815284602082015261ffff8416604082015262ffffff8316606082015260a0608082015260006110b360a0830184610c79565b979650505050505050565b6000602082840312156110d057600080fd5b815165ffffffffffff811681146109fa57600080fd5b60018060a01b038516815283602082015265ffffffffffff831660408201526080606082015260006107736080830184610c79565b6000602080838503121561112e57600080fd5b82516001600160401b038082111561114557600080fd5b818501915085601f83011261115957600080fd5b81518181111561116b5761116b610d2e565b8060051b915061117c848301610d44565b818152918301840191848101908884111561119657600080fd5b938501935b838510156111b45784518252938501939085019061119b565b98975050505050505050565b60018060a01b038816815286602082015261ffff8616604082015262ffffff8516606082015263ffffffff8416608082015260e060a0820152600061120860e0830185610c79565b82810360c084015261121a8185610c79565b9a9950505050505050505050565b818103818111156107ff57634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052603160045260246000fdfea164736f6c6343000813000a",
}

var VRFRouterABI = VRFRouterMetaData.ABI

var VRFRouterBin = VRFRouterMetaData.Bin

func DeployVRFRouter(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *VRFRouter, error) {
	parsed, err := VRFRouterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFRouterBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFRouter{VRFRouterCaller: VRFRouterCaller{contract: contract}, VRFRouterTransactor: VRFRouterTransactor{contract: contract}, VRFRouterFilterer: VRFRouterFilterer{contract: contract}}, nil
}

type VRFRouter struct {
	address common.Address
	abi     abi.ABI
	VRFRouterCaller
	VRFRouterTransactor
	VRFRouterFilterer
}

type VRFRouterCaller struct {
	contract *bind.BoundContract
}

type VRFRouterTransactor struct {
	contract *bind.BoundContract
}

type VRFRouterFilterer struct {
	contract *bind.BoundContract
}

type VRFRouterSession struct {
	Contract     *VRFRouter
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFRouterCallerSession struct {
	Contract *VRFRouterCaller
	CallOpts bind.CallOpts
}

type VRFRouterTransactorSession struct {
	Contract     *VRFRouterTransactor
	TransactOpts bind.TransactOpts
}

type VRFRouterRaw struct {
	Contract *VRFRouter
}

type VRFRouterCallerRaw struct {
	Contract *VRFRouterCaller
}

type VRFRouterTransactorRaw struct {
	Contract *VRFRouterTransactor
}

func NewVRFRouter(address common.Address, backend bind.ContractBackend) (*VRFRouter, error) {
	abi, err := abi.JSON(strings.NewReader(VRFRouterABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFRouter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFRouter{address: address, abi: abi, VRFRouterCaller: VRFRouterCaller{contract: contract}, VRFRouterTransactor: VRFRouterTransactor{contract: contract}, VRFRouterFilterer: VRFRouterFilterer{contract: contract}}, nil
}

func NewVRFRouterCaller(address common.Address, caller bind.ContractCaller) (*VRFRouterCaller, error) {
	contract, err := bindVRFRouter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFRouterCaller{contract: contract}, nil
}

func NewVRFRouterTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFRouterTransactor, error) {
	contract, err := bindVRFRouter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFRouterTransactor{contract: contract}, nil
}

func NewVRFRouterFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFRouterFilterer, error) {
	contract, err := bindVRFRouter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFRouterFilterer{contract: contract}, nil
}

func bindVRFRouter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFRouterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFRouter *VRFRouterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFRouter.Contract.VRFRouterCaller.contract.Call(opts, result, method, params...)
}

func (_VRFRouter *VRFRouterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFRouter.Contract.VRFRouterTransactor.contract.Transfer(opts)
}

func (_VRFRouter *VRFRouterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFRouter.Contract.VRFRouterTransactor.contract.Transact(opts, method, params...)
}

func (_VRFRouter *VRFRouterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFRouter.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFRouter *VRFRouterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFRouter.Contract.contract.Transfer(opts)
}

func (_VRFRouter *VRFRouterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFRouter.Contract.contract.Transact(opts, method, params...)
}

func (_VRFRouter *VRFRouterCaller) NUMCONFDELAYS(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _VRFRouter.contract.Call(opts, &out, "NUM_CONF_DELAYS")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_VRFRouter *VRFRouterSession) NUMCONFDELAYS() (uint8, error) {
	return _VRFRouter.Contract.NUMCONFDELAYS(&_VRFRouter.CallOpts)
}

func (_VRFRouter *VRFRouterCallerSession) NUMCONFDELAYS() (uint8, error) {
	return _VRFRouter.Contract.NUMCONFDELAYS(&_VRFRouter.CallOpts)
}

func (_VRFRouter *VRFRouterCaller) GetCoordinators(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _VRFRouter.contract.Call(opts, &out, "getCoordinators")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_VRFRouter *VRFRouterSession) GetCoordinators() ([]common.Address, error) {
	return _VRFRouter.Contract.GetCoordinators(&_VRFRouter.CallOpts)
}

func (_VRFRouter *VRFRouterCallerSession) GetCoordinators() ([]common.Address, error) {
	return _VRFRouter.Contract.GetCoordinators(&_VRFRouter.CallOpts)
}

func (_VRFRouter *VRFRouterCaller) GetRoute(opts *bind.CallOpts, subID *big.Int) (common.Address, error) {
	var out []interface{}
	err := _VRFRouter.contract.Call(opts, &out, "getRoute", subID)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFRouter *VRFRouterSession) GetRoute(subID *big.Int) (common.Address, error) {
	return _VRFRouter.Contract.GetRoute(&_VRFRouter.CallOpts, subID)
}

func (_VRFRouter *VRFRouterCallerSession) GetRoute(subID *big.Int) (common.Address, error) {
	return _VRFRouter.Contract.GetRoute(&_VRFRouter.CallOpts, subID)
}

func (_VRFRouter *VRFRouterCaller) IsCoordinatorRegistered(opts *bind.CallOpts, coordinatorAddress common.Address) (bool, error) {
	var out []interface{}
	err := _VRFRouter.contract.Call(opts, &out, "isCoordinatorRegistered", coordinatorAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_VRFRouter *VRFRouterSession) IsCoordinatorRegistered(coordinatorAddress common.Address) (bool, error) {
	return _VRFRouter.Contract.IsCoordinatorRegistered(&_VRFRouter.CallOpts, coordinatorAddress)
}

func (_VRFRouter *VRFRouterCallerSession) IsCoordinatorRegistered(coordinatorAddress common.Address) (bool, error) {
	return _VRFRouter.Contract.IsCoordinatorRegistered(&_VRFRouter.CallOpts, coordinatorAddress)
}

func (_VRFRouter *VRFRouterCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFRouter.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFRouter *VRFRouterSession) Owner() (common.Address, error) {
	return _VRFRouter.Contract.Owner(&_VRFRouter.CallOpts)
}

func (_VRFRouter *VRFRouterCallerSession) Owner() (common.Address, error) {
	return _VRFRouter.Contract.Owner(&_VRFRouter.CallOpts)
}

func (_VRFRouter *VRFRouterCaller) SRedemptionRoutes(opts *bind.CallOpts, arg0 *big.Int) (SRedemptionRoutes,

	error) {
	var out []interface{}
	err := _VRFRouter.contract.Call(opts, &out, "s_redemptionRoutes", arg0)

	outstruct := new(SRedemptionRoutes)
	if err != nil {
		return *outstruct, err
	}

	outstruct.RequestID = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.CoordinatorAddress = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)

	return *outstruct, err

}

func (_VRFRouter *VRFRouterSession) SRedemptionRoutes(arg0 *big.Int) (SRedemptionRoutes,

	error) {
	return _VRFRouter.Contract.SRedemptionRoutes(&_VRFRouter.CallOpts, arg0)
}

func (_VRFRouter *VRFRouterCallerSession) SRedemptionRoutes(arg0 *big.Int) (SRedemptionRoutes,

	error) {
	return _VRFRouter.Contract.SRedemptionRoutes(&_VRFRouter.CallOpts, arg0)
}

func (_VRFRouter *VRFRouterCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _VRFRouter.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_VRFRouter *VRFRouterSession) TypeAndVersion() (string, error) {
	return _VRFRouter.Contract.TypeAndVersion(&_VRFRouter.CallOpts)
}

func (_VRFRouter *VRFRouterCallerSession) TypeAndVersion() (string, error) {
	return _VRFRouter.Contract.TypeAndVersion(&_VRFRouter.CallOpts)
}

func (_VRFRouter *VRFRouterTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFRouter.contract.Transact(opts, "acceptOwnership")
}

func (_VRFRouter *VRFRouterSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFRouter.Contract.AcceptOwnership(&_VRFRouter.TransactOpts)
}

func (_VRFRouter *VRFRouterTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFRouter.Contract.AcceptOwnership(&_VRFRouter.TransactOpts)
}

func (_VRFRouter *VRFRouterTransactor) CallWithExactGasEvenIfTargetIsNoContract(opts *bind.TransactOpts, gasAmount *big.Int, target common.Address, data []byte) (*types.Transaction, error) {
	return _VRFRouter.contract.Transact(opts, "callWithExactGasEvenIfTargetIsNoContract", gasAmount, target, data)
}

func (_VRFRouter *VRFRouterSession) CallWithExactGasEvenIfTargetIsNoContract(gasAmount *big.Int, target common.Address, data []byte) (*types.Transaction, error) {
	return _VRFRouter.Contract.CallWithExactGasEvenIfTargetIsNoContract(&_VRFRouter.TransactOpts, gasAmount, target, data)
}

func (_VRFRouter *VRFRouterTransactorSession) CallWithExactGasEvenIfTargetIsNoContract(gasAmount *big.Int, target common.Address, data []byte) (*types.Transaction, error) {
	return _VRFRouter.Contract.CallWithExactGasEvenIfTargetIsNoContract(&_VRFRouter.TransactOpts, gasAmount, target, data)
}

func (_VRFRouter *VRFRouterTransactor) DeregisterCoordinator(opts *bind.TransactOpts, coordinatorAddress common.Address) (*types.Transaction, error) {
	return _VRFRouter.contract.Transact(opts, "deregisterCoordinator", coordinatorAddress)
}

func (_VRFRouter *VRFRouterSession) DeregisterCoordinator(coordinatorAddress common.Address) (*types.Transaction, error) {
	return _VRFRouter.Contract.DeregisterCoordinator(&_VRFRouter.TransactOpts, coordinatorAddress)
}

func (_VRFRouter *VRFRouterTransactorSession) DeregisterCoordinator(coordinatorAddress common.Address) (*types.Transaction, error) {
	return _VRFRouter.Contract.DeregisterCoordinator(&_VRFRouter.TransactOpts, coordinatorAddress)
}

func (_VRFRouter *VRFRouterTransactor) RedeemRandomness(opts *bind.TransactOpts, subID *big.Int, externalRequestID *big.Int, extraArgs []byte) (*types.Transaction, error) {
	return _VRFRouter.contract.Transact(opts, "redeemRandomness", subID, externalRequestID, extraArgs)
}

func (_VRFRouter *VRFRouterSession) RedeemRandomness(subID *big.Int, externalRequestID *big.Int, extraArgs []byte) (*types.Transaction, error) {
	return _VRFRouter.Contract.RedeemRandomness(&_VRFRouter.TransactOpts, subID, externalRequestID, extraArgs)
}

func (_VRFRouter *VRFRouterTransactorSession) RedeemRandomness(subID *big.Int, externalRequestID *big.Int, extraArgs []byte) (*types.Transaction, error) {
	return _VRFRouter.Contract.RedeemRandomness(&_VRFRouter.TransactOpts, subID, externalRequestID, extraArgs)
}

func (_VRFRouter *VRFRouterTransactor) RegisterCoordinator(opts *bind.TransactOpts, coordinatorAddress common.Address) (*types.Transaction, error) {
	return _VRFRouter.contract.Transact(opts, "registerCoordinator", coordinatorAddress)
}

func (_VRFRouter *VRFRouterSession) RegisterCoordinator(coordinatorAddress common.Address) (*types.Transaction, error) {
	return _VRFRouter.Contract.RegisterCoordinator(&_VRFRouter.TransactOpts, coordinatorAddress)
}

func (_VRFRouter *VRFRouterTransactorSession) RegisterCoordinator(coordinatorAddress common.Address) (*types.Transaction, error) {
	return _VRFRouter.Contract.RegisterCoordinator(&_VRFRouter.TransactOpts, coordinatorAddress)
}

func (_VRFRouter *VRFRouterTransactor) RequestRandomness(opts *bind.TransactOpts, subID *big.Int, numWords uint16, confDelay *big.Int, extraArgs []byte) (*types.Transaction, error) {
	return _VRFRouter.contract.Transact(opts, "requestRandomness", subID, numWords, confDelay, extraArgs)
}

func (_VRFRouter *VRFRouterSession) RequestRandomness(subID *big.Int, numWords uint16, confDelay *big.Int, extraArgs []byte) (*types.Transaction, error) {
	return _VRFRouter.Contract.RequestRandomness(&_VRFRouter.TransactOpts, subID, numWords, confDelay, extraArgs)
}

func (_VRFRouter *VRFRouterTransactorSession) RequestRandomness(subID *big.Int, numWords uint16, confDelay *big.Int, extraArgs []byte) (*types.Transaction, error) {
	return _VRFRouter.Contract.RequestRandomness(&_VRFRouter.TransactOpts, subID, numWords, confDelay, extraArgs)
}

func (_VRFRouter *VRFRouterTransactor) RequestRandomnessFulfillment(opts *bind.TransactOpts, subID *big.Int, numWords uint16, confDelay *big.Int, callbackGasLimit uint32, arguments []byte, extraArgs []byte) (*types.Transaction, error) {
	return _VRFRouter.contract.Transact(opts, "requestRandomnessFulfillment", subID, numWords, confDelay, callbackGasLimit, arguments, extraArgs)
}

func (_VRFRouter *VRFRouterSession) RequestRandomnessFulfillment(subID *big.Int, numWords uint16, confDelay *big.Int, callbackGasLimit uint32, arguments []byte, extraArgs []byte) (*types.Transaction, error) {
	return _VRFRouter.Contract.RequestRandomnessFulfillment(&_VRFRouter.TransactOpts, subID, numWords, confDelay, callbackGasLimit, arguments, extraArgs)
}

func (_VRFRouter *VRFRouterTransactorSession) RequestRandomnessFulfillment(subID *big.Int, numWords uint16, confDelay *big.Int, callbackGasLimit uint32, arguments []byte, extraArgs []byte) (*types.Transaction, error) {
	return _VRFRouter.Contract.RequestRandomnessFulfillment(&_VRFRouter.TransactOpts, subID, numWords, confDelay, callbackGasLimit, arguments, extraArgs)
}

func (_VRFRouter *VRFRouterTransactor) ResetRoute(opts *bind.TransactOpts, subID *big.Int) (*types.Transaction, error) {
	return _VRFRouter.contract.Transact(opts, "resetRoute", subID)
}

func (_VRFRouter *VRFRouterSession) ResetRoute(subID *big.Int) (*types.Transaction, error) {
	return _VRFRouter.Contract.ResetRoute(&_VRFRouter.TransactOpts, subID)
}

func (_VRFRouter *VRFRouterTransactorSession) ResetRoute(subID *big.Int) (*types.Transaction, error) {
	return _VRFRouter.Contract.ResetRoute(&_VRFRouter.TransactOpts, subID)
}

func (_VRFRouter *VRFRouterTransactor) SetRoute(opts *bind.TransactOpts, subID *big.Int) (*types.Transaction, error) {
	return _VRFRouter.contract.Transact(opts, "setRoute", subID)
}

func (_VRFRouter *VRFRouterSession) SetRoute(subID *big.Int) (*types.Transaction, error) {
	return _VRFRouter.Contract.SetRoute(&_VRFRouter.TransactOpts, subID)
}

func (_VRFRouter *VRFRouterTransactorSession) SetRoute(subID *big.Int) (*types.Transaction, error) {
	return _VRFRouter.Contract.SetRoute(&_VRFRouter.TransactOpts, subID)
}

func (_VRFRouter *VRFRouterTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFRouter.contract.Transact(opts, "transferOwnership", to)
}

func (_VRFRouter *VRFRouterSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFRouter.Contract.TransferOwnership(&_VRFRouter.TransactOpts, to)
}

func (_VRFRouter *VRFRouterTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFRouter.Contract.TransferOwnership(&_VRFRouter.TransactOpts, to)
}

type VRFRouterConfigSetIterator struct {
	Event *VRFRouterConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFRouterConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFRouterConfigSet)
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
		it.Event = new(VRFRouterConfigSet)
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

func (it *VRFRouterConfigSetIterator) Error() error {
	return it.fail
}

func (it *VRFRouterConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFRouterConfigSet struct {
	PreviousConfigBlockNumber uint32
	ConfigDigest              [32]byte
	ConfigCount               uint64
	Signers                   []common.Address
	Transmitters              []common.Address
	F                         uint8
	OnchainConfig             []byte
	OffchainConfigVersion     uint64
	OffchainConfig            []byte
	Raw                       types.Log
}

func (_VRFRouter *VRFRouterFilterer) FilterConfigSet(opts *bind.FilterOpts) (*VRFRouterConfigSetIterator, error) {

	logs, sub, err := _VRFRouter.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &VRFRouterConfigSetIterator{contract: _VRFRouter.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_VRFRouter *VRFRouterFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *VRFRouterConfigSet) (event.Subscription, error) {

	logs, sub, err := _VRFRouter.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFRouterConfigSet)
				if err := _VRFRouter.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_VRFRouter *VRFRouterFilterer) ParseConfigSet(log types.Log) (*VRFRouterConfigSet, error) {
	event := new(VRFRouterConfigSet)
	if err := _VRFRouter.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFRouterCoordinatorDeregisteredIterator struct {
	Event *VRFRouterCoordinatorDeregistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFRouterCoordinatorDeregisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFRouterCoordinatorDeregistered)
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
		it.Event = new(VRFRouterCoordinatorDeregistered)
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

func (it *VRFRouterCoordinatorDeregisteredIterator) Error() error {
	return it.fail
}

func (it *VRFRouterCoordinatorDeregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFRouterCoordinatorDeregistered struct {
	CoordinatorAddress common.Address
	Raw                types.Log
}

func (_VRFRouter *VRFRouterFilterer) FilterCoordinatorDeregistered(opts *bind.FilterOpts) (*VRFRouterCoordinatorDeregisteredIterator, error) {

	logs, sub, err := _VRFRouter.contract.FilterLogs(opts, "CoordinatorDeregistered")
	if err != nil {
		return nil, err
	}
	return &VRFRouterCoordinatorDeregisteredIterator{contract: _VRFRouter.contract, event: "CoordinatorDeregistered", logs: logs, sub: sub}, nil
}

func (_VRFRouter *VRFRouterFilterer) WatchCoordinatorDeregistered(opts *bind.WatchOpts, sink chan<- *VRFRouterCoordinatorDeregistered) (event.Subscription, error) {

	logs, sub, err := _VRFRouter.contract.WatchLogs(opts, "CoordinatorDeregistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFRouterCoordinatorDeregistered)
				if err := _VRFRouter.contract.UnpackLog(event, "CoordinatorDeregistered", log); err != nil {
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

func (_VRFRouter *VRFRouterFilterer) ParseCoordinatorDeregistered(log types.Log) (*VRFRouterCoordinatorDeregistered, error) {
	event := new(VRFRouterCoordinatorDeregistered)
	if err := _VRFRouter.contract.UnpackLog(event, "CoordinatorDeregistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFRouterCoordinatorRegisteredIterator struct {
	Event *VRFRouterCoordinatorRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFRouterCoordinatorRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFRouterCoordinatorRegistered)
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
		it.Event = new(VRFRouterCoordinatorRegistered)
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

func (it *VRFRouterCoordinatorRegisteredIterator) Error() error {
	return it.fail
}

func (it *VRFRouterCoordinatorRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFRouterCoordinatorRegistered struct {
	CoordinatorAddress common.Address
	Raw                types.Log
}

func (_VRFRouter *VRFRouterFilterer) FilterCoordinatorRegistered(opts *bind.FilterOpts) (*VRFRouterCoordinatorRegisteredIterator, error) {

	logs, sub, err := _VRFRouter.contract.FilterLogs(opts, "CoordinatorRegistered")
	if err != nil {
		return nil, err
	}
	return &VRFRouterCoordinatorRegisteredIterator{contract: _VRFRouter.contract, event: "CoordinatorRegistered", logs: logs, sub: sub}, nil
}

func (_VRFRouter *VRFRouterFilterer) WatchCoordinatorRegistered(opts *bind.WatchOpts, sink chan<- *VRFRouterCoordinatorRegistered) (event.Subscription, error) {

	logs, sub, err := _VRFRouter.contract.WatchLogs(opts, "CoordinatorRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFRouterCoordinatorRegistered)
				if err := _VRFRouter.contract.UnpackLog(event, "CoordinatorRegistered", log); err != nil {
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

func (_VRFRouter *VRFRouterFilterer) ParseCoordinatorRegistered(log types.Log) (*VRFRouterCoordinatorRegistered, error) {
	event := new(VRFRouterCoordinatorRegistered)
	if err := _VRFRouter.contract.UnpackLog(event, "CoordinatorRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFRouterNewTransmissionIterator struct {
	Event *VRFRouterNewTransmission

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFRouterNewTransmissionIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFRouterNewTransmission)
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
		it.Event = new(VRFRouterNewTransmission)
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

func (it *VRFRouterNewTransmissionIterator) Error() error {
	return it.fail
}

func (it *VRFRouterNewTransmissionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFRouterNewTransmission struct {
	AggregatorRoundId  uint32
	EpochAndRound      *big.Int
	Transmitter        common.Address
	JuelsPerFeeCoin    *big.Int
	ReasonableGasPrice uint64
	ConfigDigest       [32]byte
	Raw                types.Log
}

func (_VRFRouter *VRFRouterFilterer) FilterNewTransmission(opts *bind.FilterOpts, aggregatorRoundId []uint32, epochAndRound []*big.Int) (*VRFRouterNewTransmissionIterator, error) {

	var aggregatorRoundIdRule []interface{}
	for _, aggregatorRoundIdItem := range aggregatorRoundId {
		aggregatorRoundIdRule = append(aggregatorRoundIdRule, aggregatorRoundIdItem)
	}
	var epochAndRoundRule []interface{}
	for _, epochAndRoundItem := range epochAndRound {
		epochAndRoundRule = append(epochAndRoundRule, epochAndRoundItem)
	}

	logs, sub, err := _VRFRouter.contract.FilterLogs(opts, "NewTransmission", aggregatorRoundIdRule, epochAndRoundRule)
	if err != nil {
		return nil, err
	}
	return &VRFRouterNewTransmissionIterator{contract: _VRFRouter.contract, event: "NewTransmission", logs: logs, sub: sub}, nil
}

func (_VRFRouter *VRFRouterFilterer) WatchNewTransmission(opts *bind.WatchOpts, sink chan<- *VRFRouterNewTransmission, aggregatorRoundId []uint32, epochAndRound []*big.Int) (event.Subscription, error) {

	var aggregatorRoundIdRule []interface{}
	for _, aggregatorRoundIdItem := range aggregatorRoundId {
		aggregatorRoundIdRule = append(aggregatorRoundIdRule, aggregatorRoundIdItem)
	}
	var epochAndRoundRule []interface{}
	for _, epochAndRoundItem := range epochAndRound {
		epochAndRoundRule = append(epochAndRoundRule, epochAndRoundItem)
	}

	logs, sub, err := _VRFRouter.contract.WatchLogs(opts, "NewTransmission", aggregatorRoundIdRule, epochAndRoundRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFRouterNewTransmission)
				if err := _VRFRouter.contract.UnpackLog(event, "NewTransmission", log); err != nil {
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

func (_VRFRouter *VRFRouterFilterer) ParseNewTransmission(log types.Log) (*VRFRouterNewTransmission, error) {
	event := new(VRFRouterNewTransmission)
	if err := _VRFRouter.contract.UnpackLog(event, "NewTransmission", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFRouterOutputsServedIterator struct {
	Event *VRFRouterOutputsServed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFRouterOutputsServedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFRouterOutputsServed)
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
		it.Event = new(VRFRouterOutputsServed)
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

func (it *VRFRouterOutputsServedIterator) Error() error {
	return it.fail
}

func (it *VRFRouterOutputsServedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFRouterOutputsServed struct {
	RecentBlockHeight  uint64
	JuelsPerFeeCoin    *big.Int
	ReasonableGasPrice uint64
	OutputsServed      []VRFBeaconTypesOutputServed
	Raw                types.Log
}

func (_VRFRouter *VRFRouterFilterer) FilterOutputsServed(opts *bind.FilterOpts) (*VRFRouterOutputsServedIterator, error) {

	logs, sub, err := _VRFRouter.contract.FilterLogs(opts, "OutputsServed")
	if err != nil {
		return nil, err
	}
	return &VRFRouterOutputsServedIterator{contract: _VRFRouter.contract, event: "OutputsServed", logs: logs, sub: sub}, nil
}

func (_VRFRouter *VRFRouterFilterer) WatchOutputsServed(opts *bind.WatchOpts, sink chan<- *VRFRouterOutputsServed) (event.Subscription, error) {

	logs, sub, err := _VRFRouter.contract.WatchLogs(opts, "OutputsServed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFRouterOutputsServed)
				if err := _VRFRouter.contract.UnpackLog(event, "OutputsServed", log); err != nil {
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

func (_VRFRouter *VRFRouterFilterer) ParseOutputsServed(log types.Log) (*VRFRouterOutputsServed, error) {
	event := new(VRFRouterOutputsServed)
	if err := _VRFRouter.contract.UnpackLog(event, "OutputsServed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFRouterOwnershipTransferRequestedIterator struct {
	Event *VRFRouterOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFRouterOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFRouterOwnershipTransferRequested)
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
		it.Event = new(VRFRouterOwnershipTransferRequested)
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

func (it *VRFRouterOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFRouterOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFRouterOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFRouter *VRFRouterFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFRouterOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFRouter.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFRouterOwnershipTransferRequestedIterator{contract: _VRFRouter.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFRouter *VRFRouterFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFRouterOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFRouter.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFRouterOwnershipTransferRequested)
				if err := _VRFRouter.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VRFRouter *VRFRouterFilterer) ParseOwnershipTransferRequested(log types.Log) (*VRFRouterOwnershipTransferRequested, error) {
	event := new(VRFRouterOwnershipTransferRequested)
	if err := _VRFRouter.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFRouterOwnershipTransferredIterator struct {
	Event *VRFRouterOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFRouterOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFRouterOwnershipTransferred)
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
		it.Event = new(VRFRouterOwnershipTransferred)
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

func (it *VRFRouterOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFRouterOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFRouterOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFRouter *VRFRouterFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFRouterOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFRouter.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFRouterOwnershipTransferredIterator{contract: _VRFRouter.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VRFRouter *VRFRouterFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFRouterOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFRouter.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFRouterOwnershipTransferred)
				if err := _VRFRouter.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VRFRouter *VRFRouterFilterer) ParseOwnershipTransferred(log types.Log) (*VRFRouterOwnershipTransferred, error) {
	event := new(VRFRouterOwnershipTransferred)
	if err := _VRFRouter.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFRouterRandomWordsFulfilledIterator struct {
	Event *VRFRouterRandomWordsFulfilled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFRouterRandomWordsFulfilledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFRouterRandomWordsFulfilled)
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
		it.Event = new(VRFRouterRandomWordsFulfilled)
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

func (it *VRFRouterRandomWordsFulfilledIterator) Error() error {
	return it.fail
}

func (it *VRFRouterRandomWordsFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFRouterRandomWordsFulfilled struct {
	RequestIDs            []*big.Int
	SuccessfulFulfillment []byte
	TruncatedErrorData    [][]byte
	Raw                   types.Log
}

func (_VRFRouter *VRFRouterFilterer) FilterRandomWordsFulfilled(opts *bind.FilterOpts) (*VRFRouterRandomWordsFulfilledIterator, error) {

	logs, sub, err := _VRFRouter.contract.FilterLogs(opts, "RandomWordsFulfilled")
	if err != nil {
		return nil, err
	}
	return &VRFRouterRandomWordsFulfilledIterator{contract: _VRFRouter.contract, event: "RandomWordsFulfilled", logs: logs, sub: sub}, nil
}

func (_VRFRouter *VRFRouterFilterer) WatchRandomWordsFulfilled(opts *bind.WatchOpts, sink chan<- *VRFRouterRandomWordsFulfilled) (event.Subscription, error) {

	logs, sub, err := _VRFRouter.contract.WatchLogs(opts, "RandomWordsFulfilled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFRouterRandomWordsFulfilled)
				if err := _VRFRouter.contract.UnpackLog(event, "RandomWordsFulfilled", log); err != nil {
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

func (_VRFRouter *VRFRouterFilterer) ParseRandomWordsFulfilled(log types.Log) (*VRFRouterRandomWordsFulfilled, error) {
	event := new(VRFRouterRandomWordsFulfilled)
	if err := _VRFRouter.contract.UnpackLog(event, "RandomWordsFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFRouterRandomnessFulfillmentRequestedIterator struct {
	Event *VRFRouterRandomnessFulfillmentRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFRouterRandomnessFulfillmentRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFRouterRandomnessFulfillmentRequested)
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
		it.Event = new(VRFRouterRandomnessFulfillmentRequested)
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

func (it *VRFRouterRandomnessFulfillmentRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFRouterRandomnessFulfillmentRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFRouterRandomnessFulfillmentRequested struct {
	RequestID              *big.Int
	Requester              common.Address
	NextBeaconOutputHeight uint64
	ConfDelay              *big.Int
	SubID                  *big.Int
	NumWords               uint16
	GasAllowance           uint32
	GasPrice               *big.Int
	WeiPerUnitLink         *big.Int
	Arguments              []byte
	Raw                    types.Log
}

func (_VRFRouter *VRFRouterFilterer) FilterRandomnessFulfillmentRequested(opts *bind.FilterOpts, requestID []*big.Int, requester []common.Address) (*VRFRouterRandomnessFulfillmentRequestedIterator, error) {

	var requestIDRule []interface{}
	for _, requestIDItem := range requestID {
		requestIDRule = append(requestIDRule, requestIDItem)
	}
	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _VRFRouter.contract.FilterLogs(opts, "RandomnessFulfillmentRequested", requestIDRule, requesterRule)
	if err != nil {
		return nil, err
	}
	return &VRFRouterRandomnessFulfillmentRequestedIterator{contract: _VRFRouter.contract, event: "RandomnessFulfillmentRequested", logs: logs, sub: sub}, nil
}

func (_VRFRouter *VRFRouterFilterer) WatchRandomnessFulfillmentRequested(opts *bind.WatchOpts, sink chan<- *VRFRouterRandomnessFulfillmentRequested, requestID []*big.Int, requester []common.Address) (event.Subscription, error) {

	var requestIDRule []interface{}
	for _, requestIDItem := range requestID {
		requestIDRule = append(requestIDRule, requestIDItem)
	}
	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _VRFRouter.contract.WatchLogs(opts, "RandomnessFulfillmentRequested", requestIDRule, requesterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFRouterRandomnessFulfillmentRequested)
				if err := _VRFRouter.contract.UnpackLog(event, "RandomnessFulfillmentRequested", log); err != nil {
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

func (_VRFRouter *VRFRouterFilterer) ParseRandomnessFulfillmentRequested(log types.Log) (*VRFRouterRandomnessFulfillmentRequested, error) {
	event := new(VRFRouterRandomnessFulfillmentRequested)
	if err := _VRFRouter.contract.UnpackLog(event, "RandomnessFulfillmentRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFRouterRandomnessRequestedIterator struct {
	Event *VRFRouterRandomnessRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFRouterRandomnessRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFRouterRandomnessRequested)
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
		it.Event = new(VRFRouterRandomnessRequested)
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

func (it *VRFRouterRandomnessRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFRouterRandomnessRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFRouterRandomnessRequested struct {
	RequestID              *big.Int
	Requester              common.Address
	NextBeaconOutputHeight uint64
	ConfDelay              *big.Int
	SubID                  *big.Int
	NumWords               uint16
	Raw                    types.Log
}

func (_VRFRouter *VRFRouterFilterer) FilterRandomnessRequested(opts *bind.FilterOpts, requestID []*big.Int, requester []common.Address) (*VRFRouterRandomnessRequestedIterator, error) {

	var requestIDRule []interface{}
	for _, requestIDItem := range requestID {
		requestIDRule = append(requestIDRule, requestIDItem)
	}
	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _VRFRouter.contract.FilterLogs(opts, "RandomnessRequested", requestIDRule, requesterRule)
	if err != nil {
		return nil, err
	}
	return &VRFRouterRandomnessRequestedIterator{contract: _VRFRouter.contract, event: "RandomnessRequested", logs: logs, sub: sub}, nil
}

func (_VRFRouter *VRFRouterFilterer) WatchRandomnessRequested(opts *bind.WatchOpts, sink chan<- *VRFRouterRandomnessRequested, requestID []*big.Int, requester []common.Address) (event.Subscription, error) {

	var requestIDRule []interface{}
	for _, requestIDItem := range requestID {
		requestIDRule = append(requestIDRule, requestIDItem)
	}
	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _VRFRouter.contract.WatchLogs(opts, "RandomnessRequested", requestIDRule, requesterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFRouterRandomnessRequested)
				if err := _VRFRouter.contract.UnpackLog(event, "RandomnessRequested", log); err != nil {
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

func (_VRFRouter *VRFRouterFilterer) ParseRandomnessRequested(log types.Log) (*VRFRouterRandomnessRequested, error) {
	event := new(VRFRouterRandomnessRequested)
	if err := _VRFRouter.contract.UnpackLog(event, "RandomnessRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFRouterRouteSetIterator struct {
	Event *VRFRouterRouteSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFRouterRouteSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFRouterRouteSet)
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
		it.Event = new(VRFRouterRouteSet)
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

func (it *VRFRouterRouteSetIterator) Error() error {
	return it.fail
}

func (it *VRFRouterRouteSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFRouterRouteSet struct {
	SubID              *big.Int
	CoordinatorAddress common.Address
	Raw                types.Log
}

func (_VRFRouter *VRFRouterFilterer) FilterRouteSet(opts *bind.FilterOpts, subID []*big.Int) (*VRFRouterRouteSetIterator, error) {

	var subIDRule []interface{}
	for _, subIDItem := range subID {
		subIDRule = append(subIDRule, subIDItem)
	}

	logs, sub, err := _VRFRouter.contract.FilterLogs(opts, "RouteSet", subIDRule)
	if err != nil {
		return nil, err
	}
	return &VRFRouterRouteSetIterator{contract: _VRFRouter.contract, event: "RouteSet", logs: logs, sub: sub}, nil
}

func (_VRFRouter *VRFRouterFilterer) WatchRouteSet(opts *bind.WatchOpts, sink chan<- *VRFRouterRouteSet, subID []*big.Int) (event.Subscription, error) {

	var subIDRule []interface{}
	for _, subIDItem := range subID {
		subIDRule = append(subIDRule, subIDItem)
	}

	logs, sub, err := _VRFRouter.contract.WatchLogs(opts, "RouteSet", subIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFRouterRouteSet)
				if err := _VRFRouter.contract.UnpackLog(event, "RouteSet", log); err != nil {
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

func (_VRFRouter *VRFRouterFilterer) ParseRouteSet(log types.Log) (*VRFRouterRouteSet, error) {
	event := new(VRFRouterRouteSet)
	if err := _VRFRouter.contract.UnpackLog(event, "RouteSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type SRedemptionRoutes struct {
	RequestID          *big.Int
	CoordinatorAddress common.Address
}

func (_VRFRouter *VRFRouter) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFRouter.abi.Events["ConfigSet"].ID:
		return _VRFRouter.ParseConfigSet(log)
	case _VRFRouter.abi.Events["CoordinatorDeregistered"].ID:
		return _VRFRouter.ParseCoordinatorDeregistered(log)
	case _VRFRouter.abi.Events["CoordinatorRegistered"].ID:
		return _VRFRouter.ParseCoordinatorRegistered(log)
	case _VRFRouter.abi.Events["NewTransmission"].ID:
		return _VRFRouter.ParseNewTransmission(log)
	case _VRFRouter.abi.Events["OutputsServed"].ID:
		return _VRFRouter.ParseOutputsServed(log)
	case _VRFRouter.abi.Events["OwnershipTransferRequested"].ID:
		return _VRFRouter.ParseOwnershipTransferRequested(log)
	case _VRFRouter.abi.Events["OwnershipTransferred"].ID:
		return _VRFRouter.ParseOwnershipTransferred(log)
	case _VRFRouter.abi.Events["RandomWordsFulfilled"].ID:
		return _VRFRouter.ParseRandomWordsFulfilled(log)
	case _VRFRouter.abi.Events["RandomnessFulfillmentRequested"].ID:
		return _VRFRouter.ParseRandomnessFulfillmentRequested(log)
	case _VRFRouter.abi.Events["RandomnessRequested"].ID:
		return _VRFRouter.ParseRandomnessRequested(log)
	case _VRFRouter.abi.Events["RouteSet"].ID:
		return _VRFRouter.ParseRouteSet(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFRouterConfigSet) Topic() common.Hash {
	return common.HexToHash("0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05")
}

func (VRFRouterCoordinatorDeregistered) Topic() common.Hash {
	return common.HexToHash("0xf80a1a97fd42251f3c33cda98635e7399253033a6774fe37cd3f650b5282af37")
}

func (VRFRouterCoordinatorRegistered) Topic() common.Hash {
	return common.HexToHash("0xb7cabbfc11e66731fc77de0444614282023bcbd41d16781c753a431d0af01625")
}

func (VRFRouterNewTransmission) Topic() common.Hash {
	return common.HexToHash("0x27bf3f1077f091da6885751ba10f5775d06657fd59e47a6ab1f7635e5a115afe")
}

func (VRFRouterOutputsServed) Topic() common.Hash {
	return common.HexToHash("0xf10ea936d00579b4c52035ee33bf46929646b3aa87554c565d8fb2c7aa549c44")
}

func (VRFRouterOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VRFRouterOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (VRFRouterRandomWordsFulfilled) Topic() common.Hash {
	return common.HexToHash("0x47ddf7bb0cbd94c1b43c5097f1352a80db0ceb3696f029d32b24f32cd631d2b7")
}

func (VRFRouterRandomnessFulfillmentRequested) Topic() common.Hash {
	return common.HexToHash("0x24f0e469e0097d1e8d9975137f9f4dd17d2c1481b3a2f25f2382f51287eda1dc")
}

func (VRFRouterRandomnessRequested) Topic() common.Hash {
	return common.HexToHash("0xc3b31df4232b05afd212fc28027dae6fd6a81618c2a3116182cb57c7f0a3fd0a")
}

func (VRFRouterRouteSet) Topic() common.Hash {
	return common.HexToHash("0xd8bb190f51a471c0cc88b5df464ee10c29138cce0d70fb36c472d60b414f3b3a")
}

func (_VRFRouter *VRFRouter) Address() common.Address {
	return _VRFRouter.address
}

type VRFRouterInterface interface {
	NUMCONFDELAYS(opts *bind.CallOpts) (uint8, error)

	GetCoordinators(opts *bind.CallOpts) ([]common.Address, error)

	GetRoute(opts *bind.CallOpts, subID *big.Int) (common.Address, error)

	IsCoordinatorRegistered(opts *bind.CallOpts, coordinatorAddress common.Address) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SRedemptionRoutes(opts *bind.CallOpts, arg0 *big.Int) (SRedemptionRoutes,

		error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	CallWithExactGasEvenIfTargetIsNoContract(opts *bind.TransactOpts, gasAmount *big.Int, target common.Address, data []byte) (*types.Transaction, error)

	DeregisterCoordinator(opts *bind.TransactOpts, coordinatorAddress common.Address) (*types.Transaction, error)

	RedeemRandomness(opts *bind.TransactOpts, subID *big.Int, externalRequestID *big.Int, extraArgs []byte) (*types.Transaction, error)

	RegisterCoordinator(opts *bind.TransactOpts, coordinatorAddress common.Address) (*types.Transaction, error)

	RequestRandomness(opts *bind.TransactOpts, subID *big.Int, numWords uint16, confDelay *big.Int, extraArgs []byte) (*types.Transaction, error)

	RequestRandomnessFulfillment(opts *bind.TransactOpts, subID *big.Int, numWords uint16, confDelay *big.Int, callbackGasLimit uint32, arguments []byte, extraArgs []byte) (*types.Transaction, error)

	ResetRoute(opts *bind.TransactOpts, subID *big.Int) (*types.Transaction, error)

	SetRoute(opts *bind.TransactOpts, subID *big.Int) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterConfigSet(opts *bind.FilterOpts) (*VRFRouterConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *VRFRouterConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*VRFRouterConfigSet, error)

	FilterCoordinatorDeregistered(opts *bind.FilterOpts) (*VRFRouterCoordinatorDeregisteredIterator, error)

	WatchCoordinatorDeregistered(opts *bind.WatchOpts, sink chan<- *VRFRouterCoordinatorDeregistered) (event.Subscription, error)

	ParseCoordinatorDeregistered(log types.Log) (*VRFRouterCoordinatorDeregistered, error)

	FilterCoordinatorRegistered(opts *bind.FilterOpts) (*VRFRouterCoordinatorRegisteredIterator, error)

	WatchCoordinatorRegistered(opts *bind.WatchOpts, sink chan<- *VRFRouterCoordinatorRegistered) (event.Subscription, error)

	ParseCoordinatorRegistered(log types.Log) (*VRFRouterCoordinatorRegistered, error)

	FilterNewTransmission(opts *bind.FilterOpts, aggregatorRoundId []uint32, epochAndRound []*big.Int) (*VRFRouterNewTransmissionIterator, error)

	WatchNewTransmission(opts *bind.WatchOpts, sink chan<- *VRFRouterNewTransmission, aggregatorRoundId []uint32, epochAndRound []*big.Int) (event.Subscription, error)

	ParseNewTransmission(log types.Log) (*VRFRouterNewTransmission, error)

	FilterOutputsServed(opts *bind.FilterOpts) (*VRFRouterOutputsServedIterator, error)

	WatchOutputsServed(opts *bind.WatchOpts, sink chan<- *VRFRouterOutputsServed) (event.Subscription, error)

	ParseOutputsServed(log types.Log) (*VRFRouterOutputsServed, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFRouterOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFRouterOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VRFRouterOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFRouterOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFRouterOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VRFRouterOwnershipTransferred, error)

	FilterRandomWordsFulfilled(opts *bind.FilterOpts) (*VRFRouterRandomWordsFulfilledIterator, error)

	WatchRandomWordsFulfilled(opts *bind.WatchOpts, sink chan<- *VRFRouterRandomWordsFulfilled) (event.Subscription, error)

	ParseRandomWordsFulfilled(log types.Log) (*VRFRouterRandomWordsFulfilled, error)

	FilterRandomnessFulfillmentRequested(opts *bind.FilterOpts, requestID []*big.Int, requester []common.Address) (*VRFRouterRandomnessFulfillmentRequestedIterator, error)

	WatchRandomnessFulfillmentRequested(opts *bind.WatchOpts, sink chan<- *VRFRouterRandomnessFulfillmentRequested, requestID []*big.Int, requester []common.Address) (event.Subscription, error)

	ParseRandomnessFulfillmentRequested(log types.Log) (*VRFRouterRandomnessFulfillmentRequested, error)

	FilterRandomnessRequested(opts *bind.FilterOpts, requestID []*big.Int, requester []common.Address) (*VRFRouterRandomnessRequestedIterator, error)

	WatchRandomnessRequested(opts *bind.WatchOpts, sink chan<- *VRFRouterRandomnessRequested, requestID []*big.Int, requester []common.Address) (event.Subscription, error)

	ParseRandomnessRequested(log types.Log) (*VRFRouterRandomnessRequested, error)

	FilterRouteSet(opts *bind.FilterOpts, subID []*big.Int) (*VRFRouterRouteSetIterator, error)

	WatchRouteSet(opts *bind.WatchOpts, sink chan<- *VRFRouterRouteSet, subID []*big.Int) (event.Subscription, error)

	ParseRouteSet(log types.Log) (*VRFRouterRouteSet, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
