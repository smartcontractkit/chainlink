// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package load_test_beacon_consumer

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

var LoadTestBeaconVRFConsumerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"shouldFail\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"beaconPeriodBlocks\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"MustBeCoordinator\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeOwnerOrCoordinator\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"}],\"name\":\"CoordinatorUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"fail\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"reqID\",\"type\":\"uint256\"}],\"name\":\"getFulfillmentDurationByRequestID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_beaconPeriodBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pendingRequests\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestID\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_ReceivedRandomnessByRequestID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_arguments\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_averageFulfillmentInMillions\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_fastestFulfillment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_fulfillmentDurationInBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_gasAvailable\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_mostRecentRequestID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_myBeaconRequests\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"slotNumber\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"confirmationDelay\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_randomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_requestIDs\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_requestOutputHeights\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint24\",\"name\":\"\",\"type\":\"uint24\"}],\"name\":\"s_requestsIDs\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_resetCounter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_slowestFulfillment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_slowestRequestID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_subId\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_totalFulfilled\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_totalRequests\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"}],\"name\":\"setCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"shouldFail\",\"type\":\"bool\"}],\"name\":\"setFail\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"reqId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"},{\"internalType\":\"uint24\",\"name\":\"delay\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"}],\"name\":\"storeBeaconRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requestID\",\"type\":\"uint256\"}],\"name\":\"testRedeemRandomness\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"uint24\",\"name\":\"confirmationDelayArg\",\"type\":\"uint24\"}],\"name\":\"testRequestRandomness\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"uint24\",\"name\":\"confDelay\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"}],\"name\":\"testRequestRandomnessFulfillment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"uint24\",\"name\":\"confirmationDelayArg\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"batchSize\",\"type\":\"uint256\"}],\"name\":\"testRequestRandomnessFulfillmentBatch\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040526000600d556000600e556103e7600f556000601055600060115560006012553480156200003057600080fd5b5060405162001b5d38038062001b5d8339810160408190526200005391620001d0565b828282823380600081620000ae5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000e157620000e18162000125565b5050600280546001600160a01b0319166001600160a01b03939093169290921790915550600b805460ff191692151592909217909155600c55506200022792505050565b336001600160a01b038216036200017f5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401620000a5565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600080600060608486031215620001e657600080fd5b83516001600160a01b0381168114620001fe57600080fd5b602085015190935080151581146200021557600080fd5b80925050604084015190509250925092565b61192680620002376000396000f3fe608060405234801561001057600080fd5b50600436106102265760003560e01c80638866c6bd1161012a578063d0705f04116100bd578063f08c5daa1161008c578063f6eaffc811610071578063f6eaffc8146104bd578063fc7fea37146104d0578063ffe97ca4146104d957600080fd5b8063f08c5daa146104a1578063f2fde38b146104aa57600080fd5b8063d0705f0414610460578063d21ea8fd14610473578063d826f88f14610486578063ea7502ab1461048e57600080fd5b80639d769402116100f95780639d76940214610406578063a9cc471814610427578063c6d6130114610444578063cd0593df1461045757600080fd5b80638866c6bd146103c65780638d0e3165146103cf5780638da5cb5b146103d85780638ea98117146103f357600080fd5b80635f15cccc116101bd578063706da1ca1161018c57806374dba1241161017157806374dba124146103a05780637716cdaa146103a957806379ba5097146103be57600080fd5b8063706da1ca1461036a578063737144bc1461039757600080fd5b80635f15cccc1461031a578063601201d314610345578063689b77ab1461034e5780636df57cc31461035757600080fd5b80632fe8fa31116101f95780632fe8fa31146102b2578063341867a2146102dd5780634a0aee29146102f25780635a9478731461030757600080fd5b80631591950a1461022b5780631757f11c146102695780631e87f20e146102725780632b1a21301461029f575b600080fd5b610256610239366004611202565b601560209081526000928352604080842090915290825290205481565b6040519081526020015b60405180910390f35b610256600e5481565b610256610280366004611224565b6012546000908152601660209081526040808320938352929052205490565b6102566102ad366004611202565b610572565b6102566102c0366004611202565b601660209081526000928352604080842090915290825290205481565b6102f06102eb366004611202565b6105a3565b005b6102fa61066d565b604051610260919061123d565b6102fa610315366004611376565b61077d565b6102566103283660046113f7565b600460209081526000928352604080842090915290825290205481565b61025660115481565b61025660085481565b6102f0610365366004611423565b6108b4565b60095461037e9067ffffffffffffffff1681565b60405167ffffffffffffffff9091168152602001610260565b610256600d5481565b610256600f5481565b6103b16109ca565b60405161026091906114af565b6102f0610a58565b61025660105481565b61025660135481565b6000546040516001600160a01b039091168152602001610260565b6102f06104013660046114c9565b610b1b565b6102f06104143660046114f2565b600b805460ff1916911515919091179055565b600b546104349060ff1681565b6040519015158152602001610260565b610256610452366004611514565b610bcf565b610256600c5481565b61025661046e366004611202565b610ccc565b6102f0610481366004611574565b610ce8565b6102f0610d3c565b61025661049c36600461163d565b610d72565b610256600a5481565b6102f06104b83660046114c9565b610e75565b6102566104cb366004611224565b610e89565b61025660125481565b6105356104e7366004611224565b60056020526000908152604090205463ffffffff811690640100000000810462ffffff1690670100000000000000810461ffff1690690100000000000000000090046001600160a01b031684565b6040805163ffffffff909516855262ffffff909316602085015261ffff909116918301919091526001600160a01b03166060820152608001610260565b6014602052816000526040600020818154811061058e57600080fd5b90600052602060002001600091509150505481565b60025460408051602081018252600080825291517facfc6cdd00000000000000000000000000000000000000000000000000000000815291926001600160a01b03169163acfc6cdd916105fc91879187916004016116b6565b6000604051808303816000875af115801561061b573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f1916820160405261064391908101906116de565b600083815260066020908152604090912082519293506106679290918401906111a2565b50505050565b6012546000908152601460205260408120546060919067ffffffffffffffff81111561069b5761069b6112bf565b6040519080825280602002602001820160405280156106c4578160200160208202803683370190505b5090506000805b6012546000908152601460205260409020548110156107755760125460009081526014602052604081208054839081106107075761070761176f565b6000918252602080832090910154601254835260168252604080842082855290925290822054909250900361076257808484815181106107495761074961176f565b60209081029190910101528261075e8161179b565b9350505b508061076d8161179b565b9150506106cb565b508152919050565b606060008267ffffffffffffffff81111561079a5761079a6112bf565b6040519080825280602002602001820160405280156107c3578160200160208202803683370190505b5090506000600c546107d3610eaa565b6107dd91906117ca565b9050600081600c546107ed610eaa565b6107f791906117de565b61080191906117f7565b905060005b858110156108a557600061081d8c8c8c8c8c610d72565b60108054919250600061082f8361179b565b909155505060128054600090815260156020908152604080832085845282528083208790559254825260148152918120805460018101825590825291902001819055845181908690849081106108875761088761176f565b6020908102919091010152508061089d8161179b565b915050610806565b50919998505050505050505050565b600083815260046020908152604080832062ffffff861684529091528120859055600c546108e2908561180a565b6040805160808101825263ffffffff928316815262ffffff958616602080830191825261ffff968716838501908152306060850190815260009b8c52600590925293909920915182549151935199516001600160a01b03166901000000000000000000027fffffff0000000000000000000000000000000000000000ffffffffffffffffff9a90971667010000000000000002999099167fffffff00000000000000000000000000000000000000000000ffffffffffffff939097166401000000000266ffffffffffffff199091169890931697909717919091171692909217179092555050565b600780546109d79061181e565b80601f0160208091040260200160405190810160405280929190818152602001828054610a039061181e565b8015610a505780601f10610a2557610100808354040283529160200191610a50565b820191906000526020600020905b815481529060010190602001808311610a3357829003601f168201915b505050505081565b6001546001600160a01b03163314610ab75760405162461bcd60e51b815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b600080543373ffffffffffffffffffffffffffffffffffffffff19808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6000546001600160a01b03163314801590610b4157506002546001600160a01b03163314155b15610b78576040517fd4e06fd700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6002805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0383169081179091556040517fc258faa9a17ddfdf4130b4acff63a289202e7d5f9e42f366add65368575486bc90600090a250565b600080600c54610bdd610eaa565b610be791906117ca565b9050600081600c54610bf7610eaa565b610c0191906117de565b610c0b91906117f7565b60025460408051602081018252600080825291517f4ffac83a00000000000000000000000000000000000000000000000000000000815293945090926001600160a01b0390921691634ffac83a91610c6c918a918c918b9190600401611858565b6020604051808303816000875af1158015610c8b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610caf9190611890565b9050610cbd8183878a6108b4565b60088190559695505050505050565b6006602052816000526040600020818154811061058e57600080fd5b6002546001600160a01b03163314610d2c576040517f66bf9c7200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610d37838383610f34565b505050565b6000600d819055600e8190556103e7600f556010819055601181905560138190556012805491610d6b8361179b565b9190505550565b600080600c54610d80610eaa565b610d8a91906117ca565b9050600081600c54610d9a610eaa565b610da491906117de565b610dae91906117f7565b60025460408051602081018252600080825291517fdb972c8b00000000000000000000000000000000000000000000000000000000815293945090926001600160a01b039092169163db972c8b91610e13918d918d918d918d918d91906004016118a9565b6020604051808303816000875af1158015610e32573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610e569190611890565b9050610e648183898b6108b4565b600881905598975050505050505050565b610e7d611090565b610e86816110ec565b50565b60038181548110610e9957600080fd5b600091825260209091200154905081565b60004661a4b1811480610ebf575062066eed81145b15610f2d5760646001600160a01b031663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610f03573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610f279190611890565b91505090565b4391505090565b600b5460ff1615610f875760405162461bcd60e51b815260206004820152601d60248201527f206661696c656420696e2066756c66696c6c52616e646f6d576f7264730000006044820152606401610aae565b60008381526006602090815260409091208351610fa6928501906111a2565b506012546000908152601560209081526040808320868452909152812054610fcc610eaa565b610fd691906117f7565b90506000610fe782620f4240611902565b9050600e54821115610ffe57600e82905560138590555b600f54821061100f57600f54611011565b815b600f556011546110215780611054565b60115461102f9060016117de565b81601154600d546110409190611902565b61104a91906117de565b611054919061180a565b600d55601180549060006110678361179b565b909155505060125460009081526016602090815260408083209783529690529490942055505050565b6000546001600160a01b031633146110ea5760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610aae565b565b336001600160a01b038216036111445760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610aae565b6001805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b8280548282559060005260206000209081019282156111dd579160200282015b828111156111dd5782518255916020019190600101906111c2565b506111e99291506111ed565b5090565b5b808211156111e957600081556001016111ee565b6000806040838503121561121557600080fd5b50508035926020909101359150565b60006020828403121561123657600080fd5b5035919050565b6020808252825182820181905260009190848201906040850190845b8181101561127557835183529284019291840191600101611259565b50909695505050505050565b803561ffff8116811461129357600080fd5b919050565b803562ffffff8116811461129357600080fd5b803563ffffffff8116811461129357600080fd5b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff811182821017156112fe576112fe6112bf565b604052919050565b600082601f83011261131757600080fd5b813567ffffffffffffffff811115611331576113316112bf565b611344601f8201601f19166020016112d5565b81815284602083860101111561135957600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008060008060c0878903121561138f57600080fd5b8635955061139f60208801611281565b94506113ad60408801611298565b93506113bb606088016112ab565b9250608087013567ffffffffffffffff8111156113d757600080fd5b6113e389828a01611306565b92505060a087013590509295509295509295565b6000806040838503121561140a57600080fd5b8235915061141a60208401611298565b90509250929050565b6000806000806080858703121561143957600080fd5b843593506020850135925061145060408601611298565b915061145e60608601611281565b905092959194509250565b6000815180845260005b8181101561148f57602081850181015186830182015201611473565b506000602082860101526020601f19601f83011685010191505092915050565b6020815260006114c26020830184611469565b9392505050565b6000602082840312156114db57600080fd5b81356001600160a01b03811681146114c257600080fd5b60006020828403121561150457600080fd5b813580151581146114c257600080fd5b60008060006060848603121561152957600080fd5b61153284611281565b92506020840135915061154760408501611298565b90509250925092565b600067ffffffffffffffff82111561156a5761156a6112bf565b5060051b60200190565b60008060006060848603121561158957600080fd5b8335925060208085013567ffffffffffffffff808211156115a957600080fd5b818701915087601f8301126115bd57600080fd5b81356115d06115cb82611550565b6112d5565b81815260059190911b8301840190848101908a8311156115ef57600080fd5b938501935b8285101561160d578435825293850193908501906115f4565b96505050604087013592508083111561162557600080fd5b505061163386828701611306565b9150509250925092565b600080600080600060a0868803121561165557600080fd5b8535945061166560208701611281565b935061167360408701611298565b9250611681606087016112ab565b9150608086013567ffffffffffffffff81111561169d57600080fd5b6116a988828901611306565b9150509295509295909350565b8381528260208201526060604082015260006116d56060830184611469565b95945050505050565b600060208083850312156116f157600080fd5b825167ffffffffffffffff81111561170857600080fd5b8301601f8101851361171957600080fd5b80516117276115cb82611550565b81815260059190911b8201830190838101908783111561174657600080fd5b928401925b828410156117645783518252928401929084019061174b565b979650505050505050565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b6000600182016117ad576117ad611785565b5060010190565b634e487b7160e01b600052601260045260246000fd5b6000826117d9576117d96117b4565b500690565b808201808211156117f1576117f1611785565b92915050565b818103818111156117f1576117f1611785565b600082611819576118196117b4565b500490565b600181811c9082168061183257607f821691505b60208210810361185257634e487b7160e01b600052602260045260246000fd5b50919050565b84815261ffff8416602082015262ffffff831660408201526080606082015260006118866080830184611469565b9695505050505050565b6000602082840312156118a257600080fd5b5051919050565b86815261ffff8616602082015262ffffff8516604082015263ffffffff8416606082015260c0608082015260006118e360c0830185611469565b82810360a08401526118f58185611469565b9998505050505050505050565b80820281158282048414176117f1576117f161178556fea164736f6c6343000813000a",
}

var LoadTestBeaconVRFConsumerABI = LoadTestBeaconVRFConsumerMetaData.ABI

var LoadTestBeaconVRFConsumerBin = LoadTestBeaconVRFConsumerMetaData.Bin

func DeployLoadTestBeaconVRFConsumer(auth *bind.TransactOpts, backend bind.ContractBackend, coordinator common.Address, shouldFail bool, beaconPeriodBlocks *big.Int) (common.Address, *types.Transaction, *LoadTestBeaconVRFConsumer, error) {
	parsed, err := LoadTestBeaconVRFConsumerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(LoadTestBeaconVRFConsumerBin), backend, coordinator, shouldFail, beaconPeriodBlocks)
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

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LoadTestBeaconVRFConsumer.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerSession) Owner() (common.Address, error) {
	return _LoadTestBeaconVRFConsumer.Contract.Owner(&_LoadTestBeaconVRFConsumer.CallOpts)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCallerSession) Owner() (common.Address, error) {
	return _LoadTestBeaconVRFConsumer.Contract.Owner(&_LoadTestBeaconVRFConsumer.CallOpts)
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

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.contract.Transact(opts, "acceptOwnership")
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerSession) AcceptOwnership() (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.Contract.AcceptOwnership(&_LoadTestBeaconVRFConsumer.TransactOpts)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.Contract.AcceptOwnership(&_LoadTestBeaconVRFConsumer.TransactOpts)
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

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerTransactor) SetCoordinator(opts *bind.TransactOpts, coordinator common.Address) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.contract.Transact(opts, "setCoordinator", coordinator)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerSession) SetCoordinator(coordinator common.Address) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.Contract.SetCoordinator(&_LoadTestBeaconVRFConsumer.TransactOpts, coordinator)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerTransactorSession) SetCoordinator(coordinator common.Address) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.Contract.SetCoordinator(&_LoadTestBeaconVRFConsumer.TransactOpts, coordinator)
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

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerTransactor) TestRequestRandomnessFulfillment(opts *bind.TransactOpts, subID *big.Int, numWords uint16, confDelay *big.Int, callbackGasLimit uint32, arguments []byte) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.contract.Transact(opts, "testRequestRandomnessFulfillment", subID, numWords, confDelay, callbackGasLimit, arguments)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerSession) TestRequestRandomnessFulfillment(subID *big.Int, numWords uint16, confDelay *big.Int, callbackGasLimit uint32, arguments []byte) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.Contract.TestRequestRandomnessFulfillment(&_LoadTestBeaconVRFConsumer.TransactOpts, subID, numWords, confDelay, callbackGasLimit, arguments)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerTransactorSession) TestRequestRandomnessFulfillment(subID *big.Int, numWords uint16, confDelay *big.Int, callbackGasLimit uint32, arguments []byte) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.Contract.TestRequestRandomnessFulfillment(&_LoadTestBeaconVRFConsumer.TransactOpts, subID, numWords, confDelay, callbackGasLimit, arguments)
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

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.contract.Transact(opts, "transferOwnership", to)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.Contract.TransferOwnership(&_LoadTestBeaconVRFConsumer.TransactOpts, to)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.Contract.TransferOwnership(&_LoadTestBeaconVRFConsumer.TransactOpts, to)
}

type LoadTestBeaconVRFConsumerCoordinatorUpdatedIterator struct {
	Event *LoadTestBeaconVRFConsumerCoordinatorUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LoadTestBeaconVRFConsumerCoordinatorUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LoadTestBeaconVRFConsumerCoordinatorUpdated)
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
		it.Event = new(LoadTestBeaconVRFConsumerCoordinatorUpdated)
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

func (it *LoadTestBeaconVRFConsumerCoordinatorUpdatedIterator) Error() error {
	return it.fail
}

func (it *LoadTestBeaconVRFConsumerCoordinatorUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LoadTestBeaconVRFConsumerCoordinatorUpdated struct {
	Coordinator common.Address
	Raw         types.Log
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerFilterer) FilterCoordinatorUpdated(opts *bind.FilterOpts, coordinator []common.Address) (*LoadTestBeaconVRFConsumerCoordinatorUpdatedIterator, error) {

	var coordinatorRule []interface{}
	for _, coordinatorItem := range coordinator {
		coordinatorRule = append(coordinatorRule, coordinatorItem)
	}

	logs, sub, err := _LoadTestBeaconVRFConsumer.contract.FilterLogs(opts, "CoordinatorUpdated", coordinatorRule)
	if err != nil {
		return nil, err
	}
	return &LoadTestBeaconVRFConsumerCoordinatorUpdatedIterator{contract: _LoadTestBeaconVRFConsumer.contract, event: "CoordinatorUpdated", logs: logs, sub: sub}, nil
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerFilterer) WatchCoordinatorUpdated(opts *bind.WatchOpts, sink chan<- *LoadTestBeaconVRFConsumerCoordinatorUpdated, coordinator []common.Address) (event.Subscription, error) {

	var coordinatorRule []interface{}
	for _, coordinatorItem := range coordinator {
		coordinatorRule = append(coordinatorRule, coordinatorItem)
	}

	logs, sub, err := _LoadTestBeaconVRFConsumer.contract.WatchLogs(opts, "CoordinatorUpdated", coordinatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LoadTestBeaconVRFConsumerCoordinatorUpdated)
				if err := _LoadTestBeaconVRFConsumer.contract.UnpackLog(event, "CoordinatorUpdated", log); err != nil {
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

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerFilterer) ParseCoordinatorUpdated(log types.Log) (*LoadTestBeaconVRFConsumerCoordinatorUpdated, error) {
	event := new(LoadTestBeaconVRFConsumerCoordinatorUpdated)
	if err := _LoadTestBeaconVRFConsumer.contract.UnpackLog(event, "CoordinatorUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LoadTestBeaconVRFConsumerOwnershipTransferRequestedIterator struct {
	Event *LoadTestBeaconVRFConsumerOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LoadTestBeaconVRFConsumerOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LoadTestBeaconVRFConsumerOwnershipTransferRequested)
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
		it.Event = new(LoadTestBeaconVRFConsumerOwnershipTransferRequested)
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

func (it *LoadTestBeaconVRFConsumerOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *LoadTestBeaconVRFConsumerOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LoadTestBeaconVRFConsumerOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LoadTestBeaconVRFConsumerOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LoadTestBeaconVRFConsumer.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &LoadTestBeaconVRFConsumerOwnershipTransferRequestedIterator{contract: _LoadTestBeaconVRFConsumer.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *LoadTestBeaconVRFConsumerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LoadTestBeaconVRFConsumer.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LoadTestBeaconVRFConsumerOwnershipTransferRequested)
				if err := _LoadTestBeaconVRFConsumer.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerFilterer) ParseOwnershipTransferRequested(log types.Log) (*LoadTestBeaconVRFConsumerOwnershipTransferRequested, error) {
	event := new(LoadTestBeaconVRFConsumerOwnershipTransferRequested)
	if err := _LoadTestBeaconVRFConsumer.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LoadTestBeaconVRFConsumerOwnershipTransferredIterator struct {
	Event *LoadTestBeaconVRFConsumerOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LoadTestBeaconVRFConsumerOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LoadTestBeaconVRFConsumerOwnershipTransferred)
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
		it.Event = new(LoadTestBeaconVRFConsumerOwnershipTransferred)
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

func (it *LoadTestBeaconVRFConsumerOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *LoadTestBeaconVRFConsumerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LoadTestBeaconVRFConsumerOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LoadTestBeaconVRFConsumerOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LoadTestBeaconVRFConsumer.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &LoadTestBeaconVRFConsumerOwnershipTransferredIterator{contract: _LoadTestBeaconVRFConsumer.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *LoadTestBeaconVRFConsumerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LoadTestBeaconVRFConsumer.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LoadTestBeaconVRFConsumerOwnershipTransferred)
				if err := _LoadTestBeaconVRFConsumer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerFilterer) ParseOwnershipTransferred(log types.Log) (*LoadTestBeaconVRFConsumerOwnershipTransferred, error) {
	event := new(LoadTestBeaconVRFConsumerOwnershipTransferred)
	if err := _LoadTestBeaconVRFConsumer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type SMyBeaconRequests struct {
	SlotNumber        uint32
	ConfirmationDelay *big.Int
	NumWords          uint16
	Requester         common.Address
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumer) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _LoadTestBeaconVRFConsumer.abi.Events["CoordinatorUpdated"].ID:
		return _LoadTestBeaconVRFConsumer.ParseCoordinatorUpdated(log)
	case _LoadTestBeaconVRFConsumer.abi.Events["OwnershipTransferRequested"].ID:
		return _LoadTestBeaconVRFConsumer.ParseOwnershipTransferRequested(log)
	case _LoadTestBeaconVRFConsumer.abi.Events["OwnershipTransferred"].ID:
		return _LoadTestBeaconVRFConsumer.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (LoadTestBeaconVRFConsumerCoordinatorUpdated) Topic() common.Hash {
	return common.HexToHash("0xc258faa9a17ddfdf4130b4acff63a289202e7d5f9e42f366add65368575486bc")
}

func (LoadTestBeaconVRFConsumerOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (LoadTestBeaconVRFConsumerOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumer) Address() common.Address {
	return _LoadTestBeaconVRFConsumer.address
}

type LoadTestBeaconVRFConsumerInterface interface {
	Fail(opts *bind.CallOpts) (bool, error)

	GetFulfillmentDurationByRequestID(opts *bind.CallOpts, reqID *big.Int) (*big.Int, error)

	IBeaconPeriodBlocks(opts *bind.CallOpts) (*big.Int, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

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

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestID *big.Int, randomWords []*big.Int, arguments []byte) (*types.Transaction, error)

	Reset(opts *bind.TransactOpts) (*types.Transaction, error)

	SetCoordinator(opts *bind.TransactOpts, coordinator common.Address) (*types.Transaction, error)

	SetFail(opts *bind.TransactOpts, shouldFail bool) (*types.Transaction, error)

	StoreBeaconRequest(opts *bind.TransactOpts, reqId *big.Int, height *big.Int, delay *big.Int, numWords uint16) (*types.Transaction, error)

	TestRedeemRandomness(opts *bind.TransactOpts, subID *big.Int, requestID *big.Int) (*types.Transaction, error)

	TestRequestRandomness(opts *bind.TransactOpts, numWords uint16, subID *big.Int, confirmationDelayArg *big.Int) (*types.Transaction, error)

	TestRequestRandomnessFulfillment(opts *bind.TransactOpts, subID *big.Int, numWords uint16, confDelay *big.Int, callbackGasLimit uint32, arguments []byte) (*types.Transaction, error)

	TestRequestRandomnessFulfillmentBatch(opts *bind.TransactOpts, subID *big.Int, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte, batchSize *big.Int) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterCoordinatorUpdated(opts *bind.FilterOpts, coordinator []common.Address) (*LoadTestBeaconVRFConsumerCoordinatorUpdatedIterator, error)

	WatchCoordinatorUpdated(opts *bind.WatchOpts, sink chan<- *LoadTestBeaconVRFConsumerCoordinatorUpdated, coordinator []common.Address) (event.Subscription, error)

	ParseCoordinatorUpdated(log types.Log) (*LoadTestBeaconVRFConsumerCoordinatorUpdated, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LoadTestBeaconVRFConsumerOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *LoadTestBeaconVRFConsumerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*LoadTestBeaconVRFConsumerOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LoadTestBeaconVRFConsumerOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *LoadTestBeaconVRFConsumerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*LoadTestBeaconVRFConsumerOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
