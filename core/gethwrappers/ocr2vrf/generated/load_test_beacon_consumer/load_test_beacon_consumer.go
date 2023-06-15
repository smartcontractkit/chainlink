// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package load_test_beacon_consumer

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
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

var LoadTestBeaconVRFConsumerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"shouldFail\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"beaconPeriodBlocks\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"MustBeRouter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"fail\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"reqID\",\"type\":\"uint256\"}],\"name\":\"getFulfillmentDurationByRequestID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_beaconPeriodBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pendingRequests\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestID\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_ReceivedRandomnessByRequestID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_arguments\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_averageFulfillmentInMillions\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_fastestFulfillment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_fulfillmentDurationInBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_gasAvailable\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_mostRecentRequestID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_myBeaconRequests\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.SlotNumber\",\"name\":\"slotNumber\",\"type\":\"uint32\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelay\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_randomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_requestIDs\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_requestOutputHeights\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"\",\"type\":\"uint24\"}],\"name\":\"s_requestsIDs\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_resetCounter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_slowestFulfillment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_slowestRequestID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_subId\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_totalFulfilled\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_totalRequests\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"shouldFail\",\"type\":\"bool\"}],\"name\":\"setFail\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"reqId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"delay\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"}],\"name\":\"storeBeaconRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requestID\",\"type\":\"uint256\"}],\"name\":\"testRedeemRandomness\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelayArg\",\"type\":\"uint24\"}],\"name\":\"testRequestRandomness\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelayArg\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"}],\"name\":\"testRequestRandomnessFulfillment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelayArg\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"batchSize\",\"type\":\"uint256\"}],\"name\":\"testRequestRandomnessFulfillmentBatch\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040526000600b556000600c556103e7600d556000600e556000600f55600060105534801561002f57600080fd5b5060405161171f38038061171f83398101604081905261004e9161008b565b6001600160a01b03929092166080819052600680546001600160a01b03191690911790556009805460ff1916911515919091179055600a556100df565b6000806000606084860312156100a057600080fd5b83516001600160a01b03811681146100b757600080fd5b602085015190935080151581146100cd57600080fd5b80925050604084015190509250925092565b6080516116256100fa6000396000610af601526116256000f3fe608060405234801561001057600080fd5b50600436106101da5760003560e01c80637716cdaa11610104578063d0705f04116100a2578063f08c5daa11610071578063f08c5daa14610437578063f6eaffc814610440578063fc7fea3714610453578063ffe97ca41461045c57600080fd5b8063d0705f04146103f6578063d21ea8fd14610409578063d826f88f1461041c578063ea7502ab1461042457600080fd5b80639d769402116100de5780639d7694021461039c578063a9cc4718146103bd578063c6d61301146103da578063cd0593df146103ed57600080fd5b80637716cdaa146103755780638866c6bd1461038a5780638d0e31651461039357600080fd5b80635a9478731161017c5780636df57cc31161014b5780636df57cc31461030b578063706da1ca1461031e578063737144bc1461036357806374dba1241461036c57600080fd5b80635a947873146102bb5780635f15cccc146102ce578063601201d3146102f9578063689b77ab1461030257600080fd5b80632b1a2130116101b85780632b1a2130146102535780632fe8fa3114610266578063341867a2146102915780634a0aee29146102a657600080fd5b80631591950a146101df5780631757f11c1461021d5780631e87f20e14610226575b600080fd5b61020a6101ed366004610f2a565b601360209081526000928352604080842090915290825290205481565b6040519081526020015b60405180910390f35b61020a600c5481565b61020a610234366004610f4c565b6010546000908152601460209081526040808320938352929052205490565b61020a610261366004610f2a565b6104f5565b61020a610274366004610f2a565b601460209081526000928352604080842090915290825290205481565b6102a461029f366004610f2a565b610526565b005b6102ae6105f0565b6040516102149190610f65565b6102ae6102c936600461109e565b610700565b61020a6102dc36600461111f565b600160209081526000928352604080842090915290825290205481565b61020a600f5481565b61020a60055481565b6102a461031936600461114b565b610837565b60075461034a9074010000000000000000000000000000000000000000900467ffffffffffffffff1681565b60405167ffffffffffffffff9091168152602001610214565b61020a600b5481565b61020a600d5481565b61037d61094d565b60405161021491906111d7565b61020a600e5481565b61020a60115481565b6102a46103aa3660046111f1565b6009805460ff1916911515919091179055565b6009546103ca9060ff1681565b6040519015158152602001610214565b61020a6103e8366004611213565b6109db565b61020a600a5481565b61020a610404366004610f2a565b610ad8565b6102a4610417366004611273565b610af4565b6102a4610b66565b61020a61043236600461133c565b610b9c565b61020a60085481565b61020a61044e366004610f4c565b610ca5565b61020a60105481565b6104b861046a366004610f4c565b60026020526000908152604090205463ffffffff811690640100000000810462ffffff1690670100000000000000810461ffff1690690100000000000000000090046001600160a01b031684565b6040805163ffffffff909516855262ffffff909316602085015261ffff909116918301919091526001600160a01b03166060820152608001610214565b6012602052816000526040600020818154811061051157600080fd5b90600052602060002001600091509150505481565b60065460408051602081018252600080825291517facfc6cdd00000000000000000000000000000000000000000000000000000000815291926001600160a01b03169163acfc6cdd9161057f91879187916004016113b5565b6000604051808303816000875af115801561059e573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526105c691908101906113dd565b600083815260036020908152604090912082519293506105ea929091840190610eca565b50505050565b6010546000908152601260205260408120546060919067ffffffffffffffff81111561061e5761061e610fe7565b604051908082528060200260200182016040528015610647578160200160208202803683370190505b5090506000805b6010546000908152601260205260409020548110156106f857601054600090815260126020526040812080548390811061068a5761068a61146e565b600091825260208083209091015460105483526014825260408084208285529092529082205490925090036106e557808484815181106106cc576106cc61146e565b6020908102919091010152826106e18161149a565b9350505b50806106f08161149a565b91505061064e565b508152919050565b606060008267ffffffffffffffff81111561071d5761071d610fe7565b604051908082528060200260200182016040528015610746578160200160208202803683370190505b5090506000600a54610756610cc6565b61076091906114c9565b9050600081600a54610770610cc6565b61077a91906114dd565b61078491906114f6565b905060005b858110156108285760006107a08c8c8c8c8c610b9c565b600e805491925060006107b28361149a565b9091555050601080546000908152601360209081526040808320858452825280832087905592548252601281529181208054600181018255908252919020018190558451819086908490811061080a5761080a61146e565b602090810291909101015250806108208161149a565b915050610789565b50919998505050505050505050565b600083815260016020908152604080832062ffffff861684529091528120859055600a546108659085611509565b6040805160808101825263ffffffff928316815262ffffff958616602080830191825261ffff968716838501908152306060850190815260009b8c52600290925293909920915182549151935199516001600160a01b03166901000000000000000000027fffffff0000000000000000000000000000000000000000ffffffffffffffffff9a90971667010000000000000002999099167fffffff00000000000000000000000000000000000000000000ffffffffffffff939097166401000000000266ffffffffffffff199091169890931697909717919091171692909217179092555050565b6004805461095a9061151d565b80601f01602080910402602001604051908101604052809291908181526020018280546109869061151d565b80156109d35780601f106109a8576101008083540402835291602001916109d3565b820191906000526020600020905b8154815290600101906020018083116109b657829003601f168201915b505050505081565b600080600a546109e9610cc6565b6109f391906114c9565b9050600081600a54610a03610cc6565b610a0d91906114dd565b610a1791906114f6565b60065460408051602081018252600080825291517f4ffac83a00000000000000000000000000000000000000000000000000000000815293945090926001600160a01b0390921691634ffac83a91610a78918a918c918b9190600401611557565b6020604051808303816000875af1158015610a97573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610abb919061158f565b9050610ac98183878a610837565b60058190559695505050505050565b6003602052816000526040600020818154811061051157600080fd5b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03163314610b56576040517ff74c318f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610b61838383610d50565b505050565b6000600b819055600c8190556103e7600d55600e819055600f81905560118190556010805491610b958361149a565b9190505550565b6000808490506000600a54610baf610cc6565b610bb991906114c9565b9050600081600a54610bc9610cc6565b610bd391906114dd565b610bdd91906114f6565b60065460408051602081018252600080825291517fdb972c8b00000000000000000000000000000000000000000000000000000000815293945090926001600160a01b039092169163db972c8b91610c42918e918e918a918e918e91906004016115a8565b6020604051808303816000875af1158015610c61573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c85919061158f565b9050610c9381838a8c610837565b60058190559998505050505050505050565b60008181548110610cb557600080fd5b600091825260209091200154905081565b60004661a4b1811480610cdb575062066eed81145b15610d495760646001600160a01b031663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610d1f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d43919061158f565b91505090565b4391505090565b60095460ff1615610dc1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f206661696c656420696e2066756c66696c6c52616e646f6d576f726473000000604482015260640160405180910390fd5b60008381526003602090815260409091208351610de092850190610eca565b506010546000908152601360209081526040808320868452909152812054610e06610cc6565b610e1091906114f6565b90506000610e2182620f4240611601565b9050600c54821115610e3857600c82905560118590555b600d548210610e4957600d54610e4b565b815b600d55600f54610e5b5780610e8e565b600f54610e699060016114dd565b81600f54600b54610e7a9190611601565b610e8491906114dd565b610e8e9190611509565b600b55600f8054906000610ea18361149a565b909155505060105460009081526014602090815260408083209783529690529490942055505050565b828054828255906000526020600020908101928215610f05579160200282015b82811115610f05578251825591602001919060010190610eea565b50610f11929150610f15565b5090565b5b80821115610f115760008155600101610f16565b60008060408385031215610f3d57600080fd5b50508035926020909101359150565b600060208284031215610f5e57600080fd5b5035919050565b6020808252825182820181905260009190848201906040850190845b81811015610f9d57835183529284019291840191600101610f81565b50909695505050505050565b803561ffff81168114610fbb57600080fd5b919050565b803562ffffff81168114610fbb57600080fd5b803563ffffffff81168114610fbb57600080fd5b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff8111828210171561102657611026610fe7565b604052919050565b600082601f83011261103f57600080fd5b813567ffffffffffffffff81111561105957611059610fe7565b61106c601f8201601f1916602001610ffd565b81815284602083860101111561108157600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008060008060c087890312156110b757600080fd5b863595506110c760208801610fa9565b94506110d560408801610fc0565b93506110e360608801610fd3565b9250608087013567ffffffffffffffff8111156110ff57600080fd5b61110b89828a0161102e565b92505060a087013590509295509295509295565b6000806040838503121561113257600080fd5b8235915061114260208401610fc0565b90509250929050565b6000806000806080858703121561116157600080fd5b843593506020850135925061117860408601610fc0565b915061118660608601610fa9565b905092959194509250565b6000815180845260005b818110156111b75760208185018101518683018201520161119b565b506000602082860101526020601f19601f83011685010191505092915050565b6020815260006111ea6020830184611191565b9392505050565b60006020828403121561120357600080fd5b813580151581146111ea57600080fd5b60008060006060848603121561122857600080fd5b61123184610fa9565b92506020840135915061124660408501610fc0565b90509250925092565b600067ffffffffffffffff82111561126957611269610fe7565b5060051b60200190565b60008060006060848603121561128857600080fd5b8335925060208085013567ffffffffffffffff808211156112a857600080fd5b818701915087601f8301126112bc57600080fd5b81356112cf6112ca8261124f565b610ffd565b81815260059190911b8301840190848101908a8311156112ee57600080fd5b938501935b8285101561130c578435825293850193908501906112f3565b96505050604087013592508083111561132457600080fd5b50506113328682870161102e565b9150509250925092565b600080600080600060a0868803121561135457600080fd5b8535945061136460208701610fa9565b935061137260408701610fc0565b925061138060608701610fd3565b9150608086013567ffffffffffffffff81111561139c57600080fd5b6113a88882890161102e565b9150509295509295909350565b8381528260208201526060604082015260006113d46060830184611191565b95945050505050565b600060208083850312156113f057600080fd5b825167ffffffffffffffff81111561140757600080fd5b8301601f8101851361141857600080fd5b80516114266112ca8261124f565b81815260059190911b8201830190838101908783111561144557600080fd5b928401925b828410156114635783518252928401929084019061144a565b979650505050505050565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b6000600182016114ac576114ac611484565b5060010190565b634e487b7160e01b600052601260045260246000fd5b6000826114d8576114d86114b3565b500690565b808201808211156114f0576114f0611484565b92915050565b818103818111156114f0576114f0611484565b600082611518576115186114b3565b500490565b600181811c9082168061153157607f821691505b60208210810361155157634e487b7160e01b600052602260045260246000fd5b50919050565b84815261ffff8416602082015262ffffff831660408201526080606082015260006115856080830184611191565b9695505050505050565b6000602082840312156115a157600080fd5b5051919050565b86815261ffff8616602082015262ffffff8516604082015263ffffffff8416606082015260c0608082015260006115e260c0830185611191565b82810360a08401526115f48185611191565b9998505050505050505050565b80820281158282048414176114f0576114f061148456fea164736f6c6343000813000a",
}

var LoadTestBeaconVRFConsumerABI = LoadTestBeaconVRFConsumerMetaData.ABI

var LoadTestBeaconVRFConsumerBin = LoadTestBeaconVRFConsumerMetaData.Bin

func DeployLoadTestBeaconVRFConsumer(auth *bind.TransactOpts, backend bind.ContractBackend, router common.Address, shouldFail bool, beaconPeriodBlocks *big.Int) (common.Address, *types.Transaction, *LoadTestBeaconVRFConsumer, error) {
	parsed, err := LoadTestBeaconVRFConsumerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(LoadTestBeaconVRFConsumerBin), backend, router, shouldFail, beaconPeriodBlocks)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LoadTestBeaconVRFConsumer{LoadTestBeaconVRFConsumerCaller: LoadTestBeaconVRFConsumerCaller{contract: contract}, LoadTestBeaconVRFConsumerTransactor: LoadTestBeaconVRFConsumerTransactor{contract: contract}, LoadTestBeaconVRFConsumerFilterer: LoadTestBeaconVRFConsumerFilterer{contract: contract}}, nil
}

type LoadTestBeaconVRFConsumer struct {
	address common.Address
	abi     abi.ABI
	LoadTestBeaconVRFConsumerCaller
	LoadTestBeaconVRFConsumerTransactor
	LoadTestBeaconVRFConsumerFilterer
}

type LoadTestBeaconVRFConsumerCaller struct {
	contract *bind.BoundContract
}

type LoadTestBeaconVRFConsumerTransactor struct {
	contract *bind.BoundContract
}

type LoadTestBeaconVRFConsumerFilterer struct {
	contract *bind.BoundContract
}

type LoadTestBeaconVRFConsumerSession struct {
	Contract     *LoadTestBeaconVRFConsumer
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type LoadTestBeaconVRFConsumerCallerSession struct {
	Contract *LoadTestBeaconVRFConsumerCaller
	CallOpts bind.CallOpts
}

type LoadTestBeaconVRFConsumerTransactorSession struct {
	Contract     *LoadTestBeaconVRFConsumerTransactor
	TransactOpts bind.TransactOpts
}

type LoadTestBeaconVRFConsumerRaw struct {
	Contract *LoadTestBeaconVRFConsumer
}

type LoadTestBeaconVRFConsumerCallerRaw struct {
	Contract *LoadTestBeaconVRFConsumerCaller
}

type LoadTestBeaconVRFConsumerTransactorRaw struct {
	Contract *LoadTestBeaconVRFConsumerTransactor
}

func NewLoadTestBeaconVRFConsumer(address common.Address, backend bind.ContractBackend) (*LoadTestBeaconVRFConsumer, error) {
	abi, err := abi.JSON(strings.NewReader(LoadTestBeaconVRFConsumerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindLoadTestBeaconVRFConsumer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LoadTestBeaconVRFConsumer{address: address, abi: abi, LoadTestBeaconVRFConsumerCaller: LoadTestBeaconVRFConsumerCaller{contract: contract}, LoadTestBeaconVRFConsumerTransactor: LoadTestBeaconVRFConsumerTransactor{contract: contract}, LoadTestBeaconVRFConsumerFilterer: LoadTestBeaconVRFConsumerFilterer{contract: contract}}, nil
}

func NewLoadTestBeaconVRFConsumerCaller(address common.Address, caller bind.ContractCaller) (*LoadTestBeaconVRFConsumerCaller, error) {
	contract, err := bindLoadTestBeaconVRFConsumer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LoadTestBeaconVRFConsumerCaller{contract: contract}, nil
}

func NewLoadTestBeaconVRFConsumerTransactor(address common.Address, transactor bind.ContractTransactor) (*LoadTestBeaconVRFConsumerTransactor, error) {
	contract, err := bindLoadTestBeaconVRFConsumer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LoadTestBeaconVRFConsumerTransactor{contract: contract}, nil
}

func NewLoadTestBeaconVRFConsumerFilterer(address common.Address, filterer bind.ContractFilterer) (*LoadTestBeaconVRFConsumerFilterer, error) {
	contract, err := bindLoadTestBeaconVRFConsumer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LoadTestBeaconVRFConsumerFilterer{contract: contract}, nil
}

func bindLoadTestBeaconVRFConsumer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := LoadTestBeaconVRFConsumerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LoadTestBeaconVRFConsumer.Contract.LoadTestBeaconVRFConsumerCaller.contract.Call(opts, result, method, params...)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.Contract.LoadTestBeaconVRFConsumerTransactor.contract.Transfer(opts)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.Contract.LoadTestBeaconVRFConsumerTransactor.contract.Transact(opts, method, params...)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LoadTestBeaconVRFConsumer.Contract.contract.Call(opts, result, method, params...)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.Contract.contract.Transfer(opts)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.Contract.contract.Transact(opts, method, params...)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCaller) Fail(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _LoadTestBeaconVRFConsumer.contract.Call(opts, &out, "fail")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerSession) Fail() (bool, error) {
	return _LoadTestBeaconVRFConsumer.Contract.Fail(&_LoadTestBeaconVRFConsumer.CallOpts)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCallerSession) Fail() (bool, error) {
	return _LoadTestBeaconVRFConsumer.Contract.Fail(&_LoadTestBeaconVRFConsumer.CallOpts)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCaller) GetFulfillmentDurationByRequestID(opts *bind.CallOpts, reqID *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTestBeaconVRFConsumer.contract.Call(opts, &out, "getFulfillmentDurationByRequestID", reqID)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerSession) GetFulfillmentDurationByRequestID(reqID *big.Int) (*big.Int, error) {
	return _LoadTestBeaconVRFConsumer.Contract.GetFulfillmentDurationByRequestID(&_LoadTestBeaconVRFConsumer.CallOpts, reqID)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCallerSession) GetFulfillmentDurationByRequestID(reqID *big.Int) (*big.Int, error) {
	return _LoadTestBeaconVRFConsumer.Contract.GetFulfillmentDurationByRequestID(&_LoadTestBeaconVRFConsumer.CallOpts, reqID)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCaller) IBeaconPeriodBlocks(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LoadTestBeaconVRFConsumer.contract.Call(opts, &out, "i_beaconPeriodBlocks")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerSession) IBeaconPeriodBlocks() (*big.Int, error) {
	return _LoadTestBeaconVRFConsumer.Contract.IBeaconPeriodBlocks(&_LoadTestBeaconVRFConsumer.CallOpts)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCallerSession) IBeaconPeriodBlocks() (*big.Int, error) {
	return _LoadTestBeaconVRFConsumer.Contract.IBeaconPeriodBlocks(&_LoadTestBeaconVRFConsumer.CallOpts)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCaller) PendingRequests(opts *bind.CallOpts) ([]*big.Int, error) {
	var out []interface{}
	err := _LoadTestBeaconVRFConsumer.contract.Call(opts, &out, "pendingRequests")

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerSession) PendingRequests() ([]*big.Int, error) {
	return _LoadTestBeaconVRFConsumer.Contract.PendingRequests(&_LoadTestBeaconVRFConsumer.CallOpts)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCallerSession) PendingRequests() ([]*big.Int, error) {
	return _LoadTestBeaconVRFConsumer.Contract.PendingRequests(&_LoadTestBeaconVRFConsumer.CallOpts)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCaller) SReceivedRandomnessByRequestID(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTestBeaconVRFConsumer.contract.Call(opts, &out, "s_ReceivedRandomnessByRequestID", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerSession) SReceivedRandomnessByRequestID(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _LoadTestBeaconVRFConsumer.Contract.SReceivedRandomnessByRequestID(&_LoadTestBeaconVRFConsumer.CallOpts, arg0, arg1)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCallerSession) SReceivedRandomnessByRequestID(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _LoadTestBeaconVRFConsumer.Contract.SReceivedRandomnessByRequestID(&_LoadTestBeaconVRFConsumer.CallOpts, arg0, arg1)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCaller) SArguments(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _LoadTestBeaconVRFConsumer.contract.Call(opts, &out, "s_arguments")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerSession) SArguments() ([]byte, error) {
	return _LoadTestBeaconVRFConsumer.Contract.SArguments(&_LoadTestBeaconVRFConsumer.CallOpts)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCallerSession) SArguments() ([]byte, error) {
	return _LoadTestBeaconVRFConsumer.Contract.SArguments(&_LoadTestBeaconVRFConsumer.CallOpts)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCaller) SAverageFulfillmentInMillions(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LoadTestBeaconVRFConsumer.contract.Call(opts, &out, "s_averageFulfillmentInMillions")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerSession) SAverageFulfillmentInMillions() (*big.Int, error) {
	return _LoadTestBeaconVRFConsumer.Contract.SAverageFulfillmentInMillions(&_LoadTestBeaconVRFConsumer.CallOpts)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCallerSession) SAverageFulfillmentInMillions() (*big.Int, error) {
	return _LoadTestBeaconVRFConsumer.Contract.SAverageFulfillmentInMillions(&_LoadTestBeaconVRFConsumer.CallOpts)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCaller) SFastestFulfillment(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LoadTestBeaconVRFConsumer.contract.Call(opts, &out, "s_fastestFulfillment")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerSession) SFastestFulfillment() (*big.Int, error) {
	return _LoadTestBeaconVRFConsumer.Contract.SFastestFulfillment(&_LoadTestBeaconVRFConsumer.CallOpts)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCallerSession) SFastestFulfillment() (*big.Int, error) {
	return _LoadTestBeaconVRFConsumer.Contract.SFastestFulfillment(&_LoadTestBeaconVRFConsumer.CallOpts)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCaller) SFulfillmentDurationInBlocks(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTestBeaconVRFConsumer.contract.Call(opts, &out, "s_fulfillmentDurationInBlocks", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerSession) SFulfillmentDurationInBlocks(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _LoadTestBeaconVRFConsumer.Contract.SFulfillmentDurationInBlocks(&_LoadTestBeaconVRFConsumer.CallOpts, arg0, arg1)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCallerSession) SFulfillmentDurationInBlocks(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _LoadTestBeaconVRFConsumer.Contract.SFulfillmentDurationInBlocks(&_LoadTestBeaconVRFConsumer.CallOpts, arg0, arg1)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCaller) SGasAvailable(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LoadTestBeaconVRFConsumer.contract.Call(opts, &out, "s_gasAvailable")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerSession) SGasAvailable() (*big.Int, error) {
	return _LoadTestBeaconVRFConsumer.Contract.SGasAvailable(&_LoadTestBeaconVRFConsumer.CallOpts)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCallerSession) SGasAvailable() (*big.Int, error) {
	return _LoadTestBeaconVRFConsumer.Contract.SGasAvailable(&_LoadTestBeaconVRFConsumer.CallOpts)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCaller) SMostRecentRequestID(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LoadTestBeaconVRFConsumer.contract.Call(opts, &out, "s_mostRecentRequestID")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerSession) SMostRecentRequestID() (*big.Int, error) {
	return _LoadTestBeaconVRFConsumer.Contract.SMostRecentRequestID(&_LoadTestBeaconVRFConsumer.CallOpts)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCallerSession) SMostRecentRequestID() (*big.Int, error) {
	return _LoadTestBeaconVRFConsumer.Contract.SMostRecentRequestID(&_LoadTestBeaconVRFConsumer.CallOpts)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCaller) SMyBeaconRequests(opts *bind.CallOpts, arg0 *big.Int) (SMyBeaconRequests,

	error) {
	var out []interface{}
	err := _LoadTestBeaconVRFConsumer.contract.Call(opts, &out, "s_myBeaconRequests", arg0)

	outstruct := new(SMyBeaconRequests)
	if err != nil {
		return *outstruct, err
	}

	outstruct.SlotNumber = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.ConfirmationDelay = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.NumWords = *abi.ConvertType(out[2], new(uint16)).(*uint16)
	outstruct.Requester = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)

	return *outstruct, err

}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerSession) SMyBeaconRequests(arg0 *big.Int) (SMyBeaconRequests,

	error) {
	return _LoadTestBeaconVRFConsumer.Contract.SMyBeaconRequests(&_LoadTestBeaconVRFConsumer.CallOpts, arg0)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCallerSession) SMyBeaconRequests(arg0 *big.Int) (SMyBeaconRequests,

	error) {
	return _LoadTestBeaconVRFConsumer.Contract.SMyBeaconRequests(&_LoadTestBeaconVRFConsumer.CallOpts, arg0)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCaller) SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTestBeaconVRFConsumer.contract.Call(opts, &out, "s_randomWords", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _LoadTestBeaconVRFConsumer.Contract.SRandomWords(&_LoadTestBeaconVRFConsumer.CallOpts, arg0)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCallerSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _LoadTestBeaconVRFConsumer.Contract.SRandomWords(&_LoadTestBeaconVRFConsumer.CallOpts, arg0)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCaller) SRequestIDs(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTestBeaconVRFConsumer.contract.Call(opts, &out, "s_requestIDs", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerSession) SRequestIDs(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _LoadTestBeaconVRFConsumer.Contract.SRequestIDs(&_LoadTestBeaconVRFConsumer.CallOpts, arg0, arg1)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCallerSession) SRequestIDs(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _LoadTestBeaconVRFConsumer.Contract.SRequestIDs(&_LoadTestBeaconVRFConsumer.CallOpts, arg0, arg1)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCaller) SRequestOutputHeights(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTestBeaconVRFConsumer.contract.Call(opts, &out, "s_requestOutputHeights", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerSession) SRequestOutputHeights(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _LoadTestBeaconVRFConsumer.Contract.SRequestOutputHeights(&_LoadTestBeaconVRFConsumer.CallOpts, arg0, arg1)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCallerSession) SRequestOutputHeights(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _LoadTestBeaconVRFConsumer.Contract.SRequestOutputHeights(&_LoadTestBeaconVRFConsumer.CallOpts, arg0, arg1)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCaller) SRequestsIDs(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTestBeaconVRFConsumer.contract.Call(opts, &out, "s_requestsIDs", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerSession) SRequestsIDs(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _LoadTestBeaconVRFConsumer.Contract.SRequestsIDs(&_LoadTestBeaconVRFConsumer.CallOpts, arg0, arg1)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCallerSession) SRequestsIDs(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _LoadTestBeaconVRFConsumer.Contract.SRequestsIDs(&_LoadTestBeaconVRFConsumer.CallOpts, arg0, arg1)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCaller) SResetCounter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LoadTestBeaconVRFConsumer.contract.Call(opts, &out, "s_resetCounter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerSession) SResetCounter() (*big.Int, error) {
	return _LoadTestBeaconVRFConsumer.Contract.SResetCounter(&_LoadTestBeaconVRFConsumer.CallOpts)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCallerSession) SResetCounter() (*big.Int, error) {
	return _LoadTestBeaconVRFConsumer.Contract.SResetCounter(&_LoadTestBeaconVRFConsumer.CallOpts)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCaller) SSlowestFulfillment(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LoadTestBeaconVRFConsumer.contract.Call(opts, &out, "s_slowestFulfillment")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerSession) SSlowestFulfillment() (*big.Int, error) {
	return _LoadTestBeaconVRFConsumer.Contract.SSlowestFulfillment(&_LoadTestBeaconVRFConsumer.CallOpts)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCallerSession) SSlowestFulfillment() (*big.Int, error) {
	return _LoadTestBeaconVRFConsumer.Contract.SSlowestFulfillment(&_LoadTestBeaconVRFConsumer.CallOpts)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCaller) SSlowestRequestID(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LoadTestBeaconVRFConsumer.contract.Call(opts, &out, "s_slowestRequestID")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerSession) SSlowestRequestID() (*big.Int, error) {
	return _LoadTestBeaconVRFConsumer.Contract.SSlowestRequestID(&_LoadTestBeaconVRFConsumer.CallOpts)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCallerSession) SSlowestRequestID() (*big.Int, error) {
	return _LoadTestBeaconVRFConsumer.Contract.SSlowestRequestID(&_LoadTestBeaconVRFConsumer.CallOpts)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCaller) SSubId(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _LoadTestBeaconVRFConsumer.contract.Call(opts, &out, "s_subId")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerSession) SSubId() (uint64, error) {
	return _LoadTestBeaconVRFConsumer.Contract.SSubId(&_LoadTestBeaconVRFConsumer.CallOpts)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCallerSession) SSubId() (uint64, error) {
	return _LoadTestBeaconVRFConsumer.Contract.SSubId(&_LoadTestBeaconVRFConsumer.CallOpts)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCaller) STotalFulfilled(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LoadTestBeaconVRFConsumer.contract.Call(opts, &out, "s_totalFulfilled")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerSession) STotalFulfilled() (*big.Int, error) {
	return _LoadTestBeaconVRFConsumer.Contract.STotalFulfilled(&_LoadTestBeaconVRFConsumer.CallOpts)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCallerSession) STotalFulfilled() (*big.Int, error) {
	return _LoadTestBeaconVRFConsumer.Contract.STotalFulfilled(&_LoadTestBeaconVRFConsumer.CallOpts)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCaller) STotalRequests(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LoadTestBeaconVRFConsumer.contract.Call(opts, &out, "s_totalRequests")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerSession) STotalRequests() (*big.Int, error) {
	return _LoadTestBeaconVRFConsumer.Contract.STotalRequests(&_LoadTestBeaconVRFConsumer.CallOpts)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCallerSession) STotalRequests() (*big.Int, error) {
	return _LoadTestBeaconVRFConsumer.Contract.STotalRequests(&_LoadTestBeaconVRFConsumer.CallOpts)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestID *big.Int, randomWords []*big.Int, arguments []byte) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.contract.Transact(opts, "rawFulfillRandomWords", requestID, randomWords, arguments)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerSession) RawFulfillRandomWords(requestID *big.Int, randomWords []*big.Int, arguments []byte) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.Contract.RawFulfillRandomWords(&_LoadTestBeaconVRFConsumer.TransactOpts, requestID, randomWords, arguments)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerTransactorSession) RawFulfillRandomWords(requestID *big.Int, randomWords []*big.Int, arguments []byte) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.Contract.RawFulfillRandomWords(&_LoadTestBeaconVRFConsumer.TransactOpts, requestID, randomWords, arguments)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerTransactor) Reset(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.contract.Transact(opts, "reset")
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerSession) Reset() (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.Contract.Reset(&_LoadTestBeaconVRFConsumer.TransactOpts)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerTransactorSession) Reset() (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.Contract.Reset(&_LoadTestBeaconVRFConsumer.TransactOpts)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerTransactor) SetFail(opts *bind.TransactOpts, shouldFail bool) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.contract.Transact(opts, "setFail", shouldFail)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerSession) SetFail(shouldFail bool) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.Contract.SetFail(&_LoadTestBeaconVRFConsumer.TransactOpts, shouldFail)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerTransactorSession) SetFail(shouldFail bool) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.Contract.SetFail(&_LoadTestBeaconVRFConsumer.TransactOpts, shouldFail)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerTransactor) StoreBeaconRequest(opts *bind.TransactOpts, reqId *big.Int, height *big.Int, delay *big.Int, numWords uint16) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.contract.Transact(opts, "storeBeaconRequest", reqId, height, delay, numWords)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerSession) StoreBeaconRequest(reqId *big.Int, height *big.Int, delay *big.Int, numWords uint16) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.Contract.StoreBeaconRequest(&_LoadTestBeaconVRFConsumer.TransactOpts, reqId, height, delay, numWords)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerTransactorSession) StoreBeaconRequest(reqId *big.Int, height *big.Int, delay *big.Int, numWords uint16) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.Contract.StoreBeaconRequest(&_LoadTestBeaconVRFConsumer.TransactOpts, reqId, height, delay, numWords)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerTransactor) TestRedeemRandomness(opts *bind.TransactOpts, subID *big.Int, requestID *big.Int) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.contract.Transact(opts, "testRedeemRandomness", subID, requestID)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerSession) TestRedeemRandomness(subID *big.Int, requestID *big.Int) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.Contract.TestRedeemRandomness(&_LoadTestBeaconVRFConsumer.TransactOpts, subID, requestID)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerTransactorSession) TestRedeemRandomness(subID *big.Int, requestID *big.Int) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.Contract.TestRedeemRandomness(&_LoadTestBeaconVRFConsumer.TransactOpts, subID, requestID)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerTransactor) TestRequestRandomness(opts *bind.TransactOpts, numWords uint16, subID *big.Int, confirmationDelayArg *big.Int) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.contract.Transact(opts, "testRequestRandomness", numWords, subID, confirmationDelayArg)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerSession) TestRequestRandomness(numWords uint16, subID *big.Int, confirmationDelayArg *big.Int) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.Contract.TestRequestRandomness(&_LoadTestBeaconVRFConsumer.TransactOpts, numWords, subID, confirmationDelayArg)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerTransactorSession) TestRequestRandomness(numWords uint16, subID *big.Int, confirmationDelayArg *big.Int) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.Contract.TestRequestRandomness(&_LoadTestBeaconVRFConsumer.TransactOpts, numWords, subID, confirmationDelayArg)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerTransactor) TestRequestRandomnessFulfillment(opts *bind.TransactOpts, subID *big.Int, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.contract.Transact(opts, "testRequestRandomnessFulfillment", subID, numWords, confirmationDelayArg, callbackGasLimit, arguments)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerSession) TestRequestRandomnessFulfillment(subID *big.Int, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.Contract.TestRequestRandomnessFulfillment(&_LoadTestBeaconVRFConsumer.TransactOpts, subID, numWords, confirmationDelayArg, callbackGasLimit, arguments)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerTransactorSession) TestRequestRandomnessFulfillment(subID *big.Int, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.Contract.TestRequestRandomnessFulfillment(&_LoadTestBeaconVRFConsumer.TransactOpts, subID, numWords, confirmationDelayArg, callbackGasLimit, arguments)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerTransactor) TestRequestRandomnessFulfillmentBatch(opts *bind.TransactOpts, subID *big.Int, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte, batchSize *big.Int) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.contract.Transact(opts, "testRequestRandomnessFulfillmentBatch", subID, numWords, confirmationDelayArg, callbackGasLimit, arguments, batchSize)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerSession) TestRequestRandomnessFulfillmentBatch(subID *big.Int, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte, batchSize *big.Int) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.Contract.TestRequestRandomnessFulfillmentBatch(&_LoadTestBeaconVRFConsumer.TransactOpts, subID, numWords, confirmationDelayArg, callbackGasLimit, arguments, batchSize)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerTransactorSession) TestRequestRandomnessFulfillmentBatch(subID *big.Int, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte, batchSize *big.Int) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.Contract.TestRequestRandomnessFulfillmentBatch(&_LoadTestBeaconVRFConsumer.TransactOpts, subID, numWords, confirmationDelayArg, callbackGasLimit, arguments, batchSize)
}

type SMyBeaconRequests struct {
	SlotNumber        uint32
	ConfirmationDelay *big.Int
	NumWords          uint16
	Requester         common.Address
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumer) Address() common.Address {
	return _LoadTestBeaconVRFConsumer.address
}

type LoadTestBeaconVRFConsumerInterface interface {
	Fail(opts *bind.CallOpts) (bool, error)

	GetFulfillmentDurationByRequestID(opts *bind.CallOpts, reqID *big.Int) (*big.Int, error)

	IBeaconPeriodBlocks(opts *bind.CallOpts) (*big.Int, error)

	PendingRequests(opts *bind.CallOpts) ([]*big.Int, error)

	SReceivedRandomnessByRequestID(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error)

	SArguments(opts *bind.CallOpts) ([]byte, error)

	SAverageFulfillmentInMillions(opts *bind.CallOpts) (*big.Int, error)

	SFastestFulfillment(opts *bind.CallOpts) (*big.Int, error)

	SFulfillmentDurationInBlocks(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error)

	SGasAvailable(opts *bind.CallOpts) (*big.Int, error)

	SMostRecentRequestID(opts *bind.CallOpts) (*big.Int, error)

	SMyBeaconRequests(opts *bind.CallOpts, arg0 *big.Int) (SMyBeaconRequests,

		error)

	SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	SRequestIDs(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error)

	SRequestOutputHeights(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error)

	SRequestsIDs(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error)

	SResetCounter(opts *bind.CallOpts) (*big.Int, error)

	SSlowestFulfillment(opts *bind.CallOpts) (*big.Int, error)

	SSlowestRequestID(opts *bind.CallOpts) (*big.Int, error)

	SSubId(opts *bind.CallOpts) (uint64, error)

	STotalFulfilled(opts *bind.CallOpts) (*big.Int, error)

	STotalRequests(opts *bind.CallOpts) (*big.Int, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestID *big.Int, randomWords []*big.Int, arguments []byte) (*types.Transaction, error)

	Reset(opts *bind.TransactOpts) (*types.Transaction, error)

	SetFail(opts *bind.TransactOpts, shouldFail bool) (*types.Transaction, error)

	StoreBeaconRequest(opts *bind.TransactOpts, reqId *big.Int, height *big.Int, delay *big.Int, numWords uint16) (*types.Transaction, error)

	TestRedeemRandomness(opts *bind.TransactOpts, subID *big.Int, requestID *big.Int) (*types.Transaction, error)

	TestRequestRandomness(opts *bind.TransactOpts, numWords uint16, subID *big.Int, confirmationDelayArg *big.Int) (*types.Transaction, error)

	TestRequestRandomnessFulfillment(opts *bind.TransactOpts, subID *big.Int, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte) (*types.Transaction, error)

	TestRequestRandomnessFulfillmentBatch(opts *bind.TransactOpts, subID *big.Int, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte, batchSize *big.Int) (*types.Transaction, error)

	Address() common.Address
}
