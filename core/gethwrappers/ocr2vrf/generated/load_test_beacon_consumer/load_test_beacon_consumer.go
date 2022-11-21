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

type VRFBeaconTypesCallback struct {
	RequestID      *big.Int
	NumWords       uint16
	Requester      common.Address
	Arguments      []byte
	SubID          uint64
	GasAllowance   *big.Int
	GasPrice       *big.Int
	WeiPerUnitLink *big.Int
}

var LoadTestBeaconVRFConsumerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"shouldFail\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"beaconPeriodBlocks\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"nextBeaconOutputHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confDelay\",\"type\":\"uint24\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"subID\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"subID\",\"type\":\"uint64\"},{\"internalType\":\"uint96\",\"name\":\"gasAllowance\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"weiPerUnitLink\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structVRFBeaconTypes.Callback\",\"name\":\"callback\",\"type\":\"tuple\"}],\"name\":\"RandomnessFulfillmentRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"nextBeaconOutputHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confDelay\",\"type\":\"uint24\"}],\"name\":\"RandomnessRequested\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"NUM_CONF_DELAYS\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"fail\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_beaconPeriodBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_ReceivedRandomnessByRequestID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_averageFulfillmentInMillions\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_fastestFulfillment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_gasAvailable\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"}],\"name\":\"s_myBeaconRequests\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.SlotNumber\",\"name\":\"slotNumber\",\"type\":\"uint32\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelay\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_randomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"}],\"name\":\"s_requestOutputHeights\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"\",\"type\":\"uint24\"}],\"name\":\"s_requestsIDs\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_slowestFulfillment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_subId\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_totalFulfilled\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_totalRequests\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"shouldFail\",\"type\":\"bool\"}],\"name\":\"setFail\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"reqId\",\"type\":\"uint48\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"delay\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"}],\"name\":\"storeBeaconRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"}],\"name\":\"testRedeemRandomness\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"uint64\",\"name\":\"subID\",\"type\":\"uint64\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelayArg\",\"type\":\"uint24\"}],\"name\":\"testRequestRandomness\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subID\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelayArg\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"}],\"name\":\"testRequestRandomnessFulfillment\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subID\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelayArg\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"batchSize\",\"type\":\"uint256\"}],\"name\":\"testRequestRandomnessFulfillmentBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040526000600a556000600b556103e7600c556000600d556000600e5534801561002a57600080fd5b5060405161132238038061132283398101604081905261004991610086565b6001600160a01b03929092166080819052600580546001600160a01b03191690911790556008805460ff19169115159190911790556009556100da565b60008060006060848603121561009b57600080fd5b83516001600160a01b03811681146100b257600080fd5b602085015190935080151581146100c857600080fd5b80925050604084015190509250925092565b60805161122d6100f5600039600061066f015261122d6000f3fe608060405234801561001057600080fd5b50600436106101825760003560e01c8063737144bc116100d8578063a9cc47181161008c578063f08c5daa11610066578063f08c5daa146103e4578063f339c794146103ed578063f6eaffc81461040057600080fd5b8063a9cc47181461039e578063cb4192e3146103bb578063cd0593df146103db57600080fd5b80638866c6bd116100bd5780638866c6bd146103615780639c9cd0151461036a5780639d7694021461037d57600080fd5b8063737144bc1461034f57806374dba1241461035857600080fd5b80635a1c532b1161013a578063601201d311610114578063601201d3146102ee5780636d162a3e146102f7578063706da1ca1461030a57600080fd5b80635a1c532b146102955780635a47dd71146102a85780635f15cccc146102bb57600080fd5b80632f7527cc1161016b5780632f7527cc1461023c5780633d8b70aa14610256578063459076261461026b57600080fd5b80631757f11c1461018757806319a5fa22146101a3575b600080fd5b610190600b5481565b6040519081526020015b60405180910390f35b6101ff6101b1366004610bbf565b60026020526000908152604090205463ffffffff811690640100000000810462ffffff1690670100000000000000810461ffff1690690100000000000000000090046001600160a01b031684565b6040805163ffffffff909516855262ffffff909316602085015261ffff909116918301919091526001600160a01b0316606082015260800161019a565b610244600881565b60405160ff909116815260200161019a565b610269610264366004610bbf565b610413565b005b61027e610279366004610cf0565b6104d9565b60405165ffffffffffff909116815260200161019a565b6102696102a3366004610d70565b6105c9565b6102696102b6366004610e1c565b61066d565b61027e6102c9366004610eee565b600160209081526000928352604080842090915290825290205465ffffffffffff1681565b610190600e5481565b610269610305366004610f1a565b6106f5565b6006546103369074010000000000000000000000000000000000000000900467ffffffffffffffff1681565b60405167ffffffffffffffff909116815260200161019a565b610190600a5481565b610190600c5481565b610190600d5481565b61027e610378366004610f69565b610828565b61026961038b366004610fac565b6008805460ff1916911515919091179055565b6008546103ab9060ff1681565b604051901515815260200161019a565b6101906103c9366004610bbf565b600f6020526000908152604090205481565b61019060095481565b61019060075481565b6101906103fb366004610fce565b610928565b61019061040e366004610ffa565b610959565b6005546040517f74d8461100000000000000000000000000000000000000000000000000000000815265ffffffffffff831660048201526000916001600160a01b0316906374d84611906024016000604051808303816000875af115801561047f573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526104a79190810190611013565b65ffffffffffff8316600090815260036020908152604090912082519293506104d4929091840190610b48565b505050565b6000806009546104e761097a565b6104f191906110ba565b905060008160095461050161097a565b61050b91906110e4565b61051591906110fc565b6005546040517ff645dcb10000000000000000000000000000000000000000000000000000000081529192506000916001600160a01b039091169063f645dcb19061056c908c908c908c908c908c90600401611113565b6020604051808303816000875af115801561058b573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105af919061119e565b90506105bd8183898b6106f5565b98975050505050505050565b60006009546105d661097a565b6105e091906110ba565b90506000816009546105f061097a565b6105fa91906110e4565b61060491906110fc565b905060005b838110156106625760006106208a8a8a8a8a6104d9565b600d80549192506000610632836111bb565b909155505065ffffffffffff166000908152600f602052604090208290558061065a816111bb565b915050610609565b505050505050505050565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031633146106ea5760405162461bcd60e51b815260206004820152601c60248201527f6f6e6c7920636f6f7264696e61746f722063616e2066756c66696c6c0000000060448201526064015b60405180910390fd5b6104d4838383610a04565b600083815260016020908152604080832062ffffff861684529091528120805465ffffffffffff191665ffffffffffff871617905560095461073790856111d4565b6040805160808101825263ffffffff928316815262ffffff958616602080830191825261ffff968716838501908152306060850190815265ffffffffffff909b1660009081526002909252939020915182549151935199516001600160a01b03166901000000000000000000027fffffff0000000000000000000000000000000000000000ffffffffffffffffff9a90971667010000000000000002999099167fffffff00000000000000000000000000000000000000000000ffffffffffffff939097166401000000000266ffffffffffffff199091169890931697909717919091171692909217179092555050565b60008060095461083661097a565b61084091906110ba565b905060008160095461085061097a565b61085a91906110e4565b61086491906110fc565b6005546040517fdc92accf00000000000000000000000000000000000000000000000000000000815261ffff8916600482015267ffffffffffffffff8816602482015262ffffff871660448201529192506000916001600160a01b039091169063dc92accf906064016020604051808303816000875af11580156108ec573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610910919061119e565b905061091e8183878a6106f5565b9695505050505050565b6003602052816000526040600020818154811061094457600080fd5b90600052602060002001600091509150505481565b6000818154811061096957600080fd5b600091825260209091200154905081565b60004661a4b181148061098f575062066eed81145b156109fd5760646001600160a01b031663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156109d3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906109f791906111e8565b91505090565b4391505090565b60085460ff1615610a575760405162461bcd60e51b815260206004820152601d60248201527f206661696c656420696e2066756c66696c6c52616e646f6d576f72647300000060448201526064016106e1565b65ffffffffffff831660009081526003602090815260409091208351610a7f92850190610b48565b5065ffffffffffff83166000908152600f6020526040812054610aa061097a565b610aaa91906110fc565b90506000610abb82620f4240611201565b9050600b548211610ace57600b54610ad0565b815b600b55600c548210610ae457600c54610ae6565b815b600c55600e54610af65780610b29565b600e54610b049060016110e4565b81600e54600a54610b159190611201565b610b1f91906110e4565b610b2991906111d4565b600a55600e8054906000610b3c836111bb565b91905055505050505050565b828054828255906000526020600020908101928215610b83579160200282015b82811115610b83578251825591602001919060010190610b68565b50610b8f929150610b93565b5090565b5b80821115610b8f5760008155600101610b94565b65ffffffffffff81168114610bbc57600080fd5b50565b600060208284031215610bd157600080fd5b8135610bdc81610ba8565b9392505050565b803567ffffffffffffffff81168114610bfb57600080fd5b919050565b803561ffff81168114610bfb57600080fd5b803562ffffff81168114610bfb57600080fd5b803563ffffffff81168114610bfb57600080fd5b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff81118282101715610c7857610c78610c39565b604052919050565b600082601f830112610c9157600080fd5b813567ffffffffffffffff811115610cab57610cab610c39565b610cbe601f8201601f1916602001610c4f565b818152846020838601011115610cd357600080fd5b816020850160208301376000918101602001919091529392505050565b600080600080600060a08688031215610d0857600080fd5b610d1186610be3565b9450610d1f60208701610c00565b9350610d2d60408701610c12565b9250610d3b60608701610c25565b9150608086013567ffffffffffffffff811115610d5757600080fd5b610d6388828901610c80565b9150509295509295909350565b60008060008060008060c08789031215610d8957600080fd5b610d9287610be3565b9550610da060208801610c00565b9450610dae60408801610c12565b9350610dbc60608801610c25565b9250608087013567ffffffffffffffff811115610dd857600080fd5b610de489828a01610c80565b92505060a087013590509295509295509295565b600067ffffffffffffffff821115610e1257610e12610c39565b5060051b60200190565b600080600060608486031215610e3157600080fd5b8335610e3c81610ba8565b925060208481013567ffffffffffffffff80821115610e5a57600080fd5b818701915087601f830112610e6e57600080fd5b8135610e81610e7c82610df8565b610c4f565b81815260059190911b8301840190848101908a831115610ea057600080fd5b938501935b82851015610ebe57843582529385019390850190610ea5565b965050506040870135925080831115610ed657600080fd5b5050610ee486828701610c80565b9150509250925092565b60008060408385031215610f0157600080fd5b82359150610f1160208401610c12565b90509250929050565b60008060008060808587031215610f3057600080fd5b8435610f3b81610ba8565b935060208501359250610f5060408601610c12565b9150610f5e60608601610c00565b905092959194509250565b600080600060608486031215610f7e57600080fd5b610f8784610c00565b9250610f9560208501610be3565b9150610fa360408501610c12565b90509250925092565b600060208284031215610fbe57600080fd5b81358015158114610bdc57600080fd5b60008060408385031215610fe157600080fd5b8235610fec81610ba8565b946020939093013593505050565b60006020828403121561100c57600080fd5b5035919050565b6000602080838503121561102657600080fd5b825167ffffffffffffffff81111561103d57600080fd5b8301601f8101851361104e57600080fd5b805161105c610e7c82610df8565b81815260059190911b8201830190838101908783111561107b57600080fd5b928401925b8284101561109957835182529284019290840190611080565b979650505050505050565b634e487b7160e01b600052601260045260246000fd5b6000826110c9576110c96110a4565b500690565b634e487b7160e01b600052601160045260246000fd5b600082198211156110f7576110f76110ce565b500190565b60008282101561110e5761110e6110ce565b500390565b67ffffffffffffffff861681526000602061ffff87168184015262ffffff8616604084015263ffffffff8516606084015260a0608084015283518060a085015260005b818110156111725785810183015185820160c001528201611156565b8181111561118457600060c083870101525b50601f01601f19169290920160c001979650505050505050565b6000602082840312156111b057600080fd5b8151610bdc81610ba8565b6000600182016111cd576111cd6110ce565b5060010190565b6000826111e3576111e36110a4565b500490565b6000602082840312156111fa57600080fd5b5051919050565b600081600019048311821515161561121b5761121b6110ce565b50029056fea164736f6c634300080f000a",
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
	parsed, err := abi.JSON(strings.NewReader(LoadTestBeaconVRFConsumerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
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

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCaller) NUMCONFDELAYS(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _LoadTestBeaconVRFConsumer.contract.Call(opts, &out, "NUM_CONF_DELAYS")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerSession) NUMCONFDELAYS() (uint8, error) {
	return _LoadTestBeaconVRFConsumer.Contract.NUMCONFDELAYS(&_LoadTestBeaconVRFConsumer.CallOpts)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCallerSession) NUMCONFDELAYS() (uint8, error) {
	return _LoadTestBeaconVRFConsumer.Contract.NUMCONFDELAYS(&_LoadTestBeaconVRFConsumer.CallOpts)
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

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCaller) SRequestOutputHeights(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTestBeaconVRFConsumer.contract.Call(opts, &out, "s_requestOutputHeights", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerSession) SRequestOutputHeights(arg0 *big.Int) (*big.Int, error) {
	return _LoadTestBeaconVRFConsumer.Contract.SRequestOutputHeights(&_LoadTestBeaconVRFConsumer.CallOpts, arg0)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCallerSession) SRequestOutputHeights(arg0 *big.Int) (*big.Int, error) {
	return _LoadTestBeaconVRFConsumer.Contract.SRequestOutputHeights(&_LoadTestBeaconVRFConsumer.CallOpts, arg0)
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

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerTransactor) TestRedeemRandomness(opts *bind.TransactOpts, requestID *big.Int) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.contract.Transact(opts, "testRedeemRandomness", requestID)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerSession) TestRedeemRandomness(requestID *big.Int) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.Contract.TestRedeemRandomness(&_LoadTestBeaconVRFConsumer.TransactOpts, requestID)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerTransactorSession) TestRedeemRandomness(requestID *big.Int) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.Contract.TestRedeemRandomness(&_LoadTestBeaconVRFConsumer.TransactOpts, requestID)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerTransactor) TestRequestRandomness(opts *bind.TransactOpts, numWords uint16, subID uint64, confirmationDelayArg *big.Int) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.contract.Transact(opts, "testRequestRandomness", numWords, subID, confirmationDelayArg)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerSession) TestRequestRandomness(numWords uint16, subID uint64, confirmationDelayArg *big.Int) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.Contract.TestRequestRandomness(&_LoadTestBeaconVRFConsumer.TransactOpts, numWords, subID, confirmationDelayArg)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerTransactorSession) TestRequestRandomness(numWords uint16, subID uint64, confirmationDelayArg *big.Int) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.Contract.TestRequestRandomness(&_LoadTestBeaconVRFConsumer.TransactOpts, numWords, subID, confirmationDelayArg)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerTransactor) TestRequestRandomnessFulfillment(opts *bind.TransactOpts, subID uint64, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.contract.Transact(opts, "testRequestRandomnessFulfillment", subID, numWords, confirmationDelayArg, callbackGasLimit, arguments)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerSession) TestRequestRandomnessFulfillment(subID uint64, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.Contract.TestRequestRandomnessFulfillment(&_LoadTestBeaconVRFConsumer.TransactOpts, subID, numWords, confirmationDelayArg, callbackGasLimit, arguments)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerTransactorSession) TestRequestRandomnessFulfillment(subID uint64, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.Contract.TestRequestRandomnessFulfillment(&_LoadTestBeaconVRFConsumer.TransactOpts, subID, numWords, confirmationDelayArg, callbackGasLimit, arguments)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerTransactor) TestRequestRandomnessFulfillmentBatch(opts *bind.TransactOpts, subID uint64, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte, batchSize *big.Int) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.contract.Transact(opts, "testRequestRandomnessFulfillmentBatch", subID, numWords, confirmationDelayArg, callbackGasLimit, arguments, batchSize)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerSession) TestRequestRandomnessFulfillmentBatch(subID uint64, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte, batchSize *big.Int) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.Contract.TestRequestRandomnessFulfillmentBatch(&_LoadTestBeaconVRFConsumer.TransactOpts, subID, numWords, confirmationDelayArg, callbackGasLimit, arguments, batchSize)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerTransactorSession) TestRequestRandomnessFulfillmentBatch(subID uint64, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte, batchSize *big.Int) (*types.Transaction, error) {
	return _LoadTestBeaconVRFConsumer.Contract.TestRequestRandomnessFulfillmentBatch(&_LoadTestBeaconVRFConsumer.TransactOpts, subID, numWords, confirmationDelayArg, callbackGasLimit, arguments, batchSize)
}

type LoadTestBeaconVRFConsumerConfigSetIterator struct {
	Event *LoadTestBeaconVRFConsumerConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LoadTestBeaconVRFConsumerConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LoadTestBeaconVRFConsumerConfigSet)
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
		it.Event = new(LoadTestBeaconVRFConsumerConfigSet)
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

func (it *LoadTestBeaconVRFConsumerConfigSetIterator) Error() error {
	return it.fail
}

func (it *LoadTestBeaconVRFConsumerConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LoadTestBeaconVRFConsumerConfigSet struct {
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

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerFilterer) FilterConfigSet(opts *bind.FilterOpts) (*LoadTestBeaconVRFConsumerConfigSetIterator, error) {

	logs, sub, err := _LoadTestBeaconVRFConsumer.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &LoadTestBeaconVRFConsumerConfigSetIterator{contract: _LoadTestBeaconVRFConsumer.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *LoadTestBeaconVRFConsumerConfigSet) (event.Subscription, error) {

	logs, sub, err := _LoadTestBeaconVRFConsumer.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LoadTestBeaconVRFConsumerConfigSet)
				if err := _LoadTestBeaconVRFConsumer.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerFilterer) ParseConfigSet(log types.Log) (*LoadTestBeaconVRFConsumerConfigSet, error) {
	event := new(LoadTestBeaconVRFConsumerConfigSet)
	if err := _LoadTestBeaconVRFConsumer.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LoadTestBeaconVRFConsumerRandomnessFulfillmentRequestedIterator struct {
	Event *LoadTestBeaconVRFConsumerRandomnessFulfillmentRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LoadTestBeaconVRFConsumerRandomnessFulfillmentRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LoadTestBeaconVRFConsumerRandomnessFulfillmentRequested)
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
		it.Event = new(LoadTestBeaconVRFConsumerRandomnessFulfillmentRequested)
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

func (it *LoadTestBeaconVRFConsumerRandomnessFulfillmentRequestedIterator) Error() error {
	return it.fail
}

func (it *LoadTestBeaconVRFConsumerRandomnessFulfillmentRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LoadTestBeaconVRFConsumerRandomnessFulfillmentRequested struct {
	NextBeaconOutputHeight uint64
	ConfDelay              *big.Int
	SubID                  uint64
	Callback               VRFBeaconTypesCallback
	Raw                    types.Log
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerFilterer) FilterRandomnessFulfillmentRequested(opts *bind.FilterOpts) (*LoadTestBeaconVRFConsumerRandomnessFulfillmentRequestedIterator, error) {

	logs, sub, err := _LoadTestBeaconVRFConsumer.contract.FilterLogs(opts, "RandomnessFulfillmentRequested")
	if err != nil {
		return nil, err
	}
	return &LoadTestBeaconVRFConsumerRandomnessFulfillmentRequestedIterator{contract: _LoadTestBeaconVRFConsumer.contract, event: "RandomnessFulfillmentRequested", logs: logs, sub: sub}, nil
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerFilterer) WatchRandomnessFulfillmentRequested(opts *bind.WatchOpts, sink chan<- *LoadTestBeaconVRFConsumerRandomnessFulfillmentRequested) (event.Subscription, error) {

	logs, sub, err := _LoadTestBeaconVRFConsumer.contract.WatchLogs(opts, "RandomnessFulfillmentRequested")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LoadTestBeaconVRFConsumerRandomnessFulfillmentRequested)
				if err := _LoadTestBeaconVRFConsumer.contract.UnpackLog(event, "RandomnessFulfillmentRequested", log); err != nil {
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

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerFilterer) ParseRandomnessFulfillmentRequested(log types.Log) (*LoadTestBeaconVRFConsumerRandomnessFulfillmentRequested, error) {
	event := new(LoadTestBeaconVRFConsumerRandomnessFulfillmentRequested)
	if err := _LoadTestBeaconVRFConsumer.contract.UnpackLog(event, "RandomnessFulfillmentRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LoadTestBeaconVRFConsumerRandomnessRequestedIterator struct {
	Event *LoadTestBeaconVRFConsumerRandomnessRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LoadTestBeaconVRFConsumerRandomnessRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LoadTestBeaconVRFConsumerRandomnessRequested)
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
		it.Event = new(LoadTestBeaconVRFConsumerRandomnessRequested)
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

func (it *LoadTestBeaconVRFConsumerRandomnessRequestedIterator) Error() error {
	return it.fail
}

func (it *LoadTestBeaconVRFConsumerRandomnessRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LoadTestBeaconVRFConsumerRandomnessRequested struct {
	NextBeaconOutputHeight uint64
	ConfDelay              *big.Int
	Raw                    types.Log
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerFilterer) FilterRandomnessRequested(opts *bind.FilterOpts, nextBeaconOutputHeight []uint64) (*LoadTestBeaconVRFConsumerRandomnessRequestedIterator, error) {

	var nextBeaconOutputHeightRule []interface{}
	for _, nextBeaconOutputHeightItem := range nextBeaconOutputHeight {
		nextBeaconOutputHeightRule = append(nextBeaconOutputHeightRule, nextBeaconOutputHeightItem)
	}

	logs, sub, err := _LoadTestBeaconVRFConsumer.contract.FilterLogs(opts, "RandomnessRequested", nextBeaconOutputHeightRule)
	if err != nil {
		return nil, err
	}
	return &LoadTestBeaconVRFConsumerRandomnessRequestedIterator{contract: _LoadTestBeaconVRFConsumer.contract, event: "RandomnessRequested", logs: logs, sub: sub}, nil
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerFilterer) WatchRandomnessRequested(opts *bind.WatchOpts, sink chan<- *LoadTestBeaconVRFConsumerRandomnessRequested, nextBeaconOutputHeight []uint64) (event.Subscription, error) {

	var nextBeaconOutputHeightRule []interface{}
	for _, nextBeaconOutputHeightItem := range nextBeaconOutputHeight {
		nextBeaconOutputHeightRule = append(nextBeaconOutputHeightRule, nextBeaconOutputHeightItem)
	}

	logs, sub, err := _LoadTestBeaconVRFConsumer.contract.WatchLogs(opts, "RandomnessRequested", nextBeaconOutputHeightRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LoadTestBeaconVRFConsumerRandomnessRequested)
				if err := _LoadTestBeaconVRFConsumer.contract.UnpackLog(event, "RandomnessRequested", log); err != nil {
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

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerFilterer) ParseRandomnessRequested(log types.Log) (*LoadTestBeaconVRFConsumerRandomnessRequested, error) {
	event := new(LoadTestBeaconVRFConsumerRandomnessRequested)
	if err := _LoadTestBeaconVRFConsumer.contract.UnpackLog(event, "RandomnessRequested", log); err != nil {
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
	case _LoadTestBeaconVRFConsumer.abi.Events["ConfigSet"].ID:
		return _LoadTestBeaconVRFConsumer.ParseConfigSet(log)
	case _LoadTestBeaconVRFConsumer.abi.Events["RandomnessFulfillmentRequested"].ID:
		return _LoadTestBeaconVRFConsumer.ParseRandomnessFulfillmentRequested(log)
	case _LoadTestBeaconVRFConsumer.abi.Events["RandomnessRequested"].ID:
		return _LoadTestBeaconVRFConsumer.ParseRandomnessRequested(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (LoadTestBeaconVRFConsumerConfigSet) Topic() common.Hash {
	return common.HexToHash("0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05")
}

func (LoadTestBeaconVRFConsumerRandomnessFulfillmentRequested) Topic() common.Hash {
	return common.HexToHash("0x230a611a84a1119ffcb8fd157032b598d075e6d783856460329674b80c1c0ac9")
}

func (LoadTestBeaconVRFConsumerRandomnessRequested) Topic() common.Hash {
	return common.HexToHash("0xc334d6f57be304c8192da2e39220c48e35f7e9afa16c541e68a6a859eff4dbc5")
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumer) Address() common.Address {
	return _LoadTestBeaconVRFConsumer.address
}

type LoadTestBeaconVRFConsumerInterface interface {
	NUMCONFDELAYS(opts *bind.CallOpts) (uint8, error)

	Fail(opts *bind.CallOpts) (bool, error)

	IBeaconPeriodBlocks(opts *bind.CallOpts) (*big.Int, error)

	SReceivedRandomnessByRequestID(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error)

	SAverageFulfillmentInMillions(opts *bind.CallOpts) (*big.Int, error)

	SFastestFulfillment(opts *bind.CallOpts) (*big.Int, error)

	SGasAvailable(opts *bind.CallOpts) (*big.Int, error)

	SMyBeaconRequests(opts *bind.CallOpts, arg0 *big.Int) (SMyBeaconRequests,

		error)

	SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	SRequestOutputHeights(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	SRequestsIDs(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error)

	SSlowestFulfillment(opts *bind.CallOpts) (*big.Int, error)

	SSubId(opts *bind.CallOpts) (uint64, error)

	STotalFulfilled(opts *bind.CallOpts) (*big.Int, error)

	STotalRequests(opts *bind.CallOpts) (*big.Int, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestID *big.Int, randomWords []*big.Int, arguments []byte) (*types.Transaction, error)

	SetFail(opts *bind.TransactOpts, shouldFail bool) (*types.Transaction, error)

	StoreBeaconRequest(opts *bind.TransactOpts, reqId *big.Int, height *big.Int, delay *big.Int, numWords uint16) (*types.Transaction, error)

	TestRedeemRandomness(opts *bind.TransactOpts, requestID *big.Int) (*types.Transaction, error)

	TestRequestRandomness(opts *bind.TransactOpts, numWords uint16, subID uint64, confirmationDelayArg *big.Int) (*types.Transaction, error)

	TestRequestRandomnessFulfillment(opts *bind.TransactOpts, subID uint64, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte) (*types.Transaction, error)

	TestRequestRandomnessFulfillmentBatch(opts *bind.TransactOpts, subID uint64, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte, batchSize *big.Int) (*types.Transaction, error)

	FilterConfigSet(opts *bind.FilterOpts) (*LoadTestBeaconVRFConsumerConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *LoadTestBeaconVRFConsumerConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*LoadTestBeaconVRFConsumerConfigSet, error)

	FilterRandomnessFulfillmentRequested(opts *bind.FilterOpts) (*LoadTestBeaconVRFConsumerRandomnessFulfillmentRequestedIterator, error)

	WatchRandomnessFulfillmentRequested(opts *bind.WatchOpts, sink chan<- *LoadTestBeaconVRFConsumerRandomnessFulfillmentRequested) (event.Subscription, error)

	ParseRandomnessFulfillmentRequested(log types.Log) (*LoadTestBeaconVRFConsumerRandomnessFulfillmentRequested, error)

	FilterRandomnessRequested(opts *bind.FilterOpts, nextBeaconOutputHeight []uint64) (*LoadTestBeaconVRFConsumerRandomnessRequestedIterator, error)

	WatchRandomnessRequested(opts *bind.WatchOpts, sink chan<- *LoadTestBeaconVRFConsumerRandomnessRequested, nextBeaconOutputHeight []uint64) (event.Subscription, error)

	ParseRandomnessRequested(log types.Log) (*LoadTestBeaconVRFConsumerRandomnessRequested, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
