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

type VRFBeaconTypesOutputServed struct {
	Height            uint64
	ConfirmationDelay *big.Int
	ProofG1X          *big.Int
	ProofG1Y          *big.Int
}

var LoadTestBeaconVRFConsumerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"shouldFail\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"beaconPeriodBlocks\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"aggregatorRoundId\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"uint40\",\"name\":\"epochAndRound\",\"type\":\"uint40\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint192\",\"name\":\"juelsPerFeeCoin\",\"type\":\"uint192\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"NewTransmission\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"recentBlockHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint192\",\"name\":\"juelsPerFeeCoin\",\"type\":\"uint192\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelay\",\"type\":\"uint24\"},{\"internalType\":\"uint256\",\"name\":\"proofG1X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"proofG1Y\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structVRFBeaconTypes.OutputServed[]\",\"name\":\"outputsServed\",\"type\":\"tuple[]\"}],\"name\":\"OutputsServed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"VRFBeaconTypes.RequestID[]\",\"name\":\"requestIDs\",\"type\":\"uint48[]\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"successfulFulfillment\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes[]\",\"name\":\"truncatedErrorData\",\"type\":\"bytes[]\"}],\"name\":\"RandomWordsFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"nextBeaconOutputHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confDelay\",\"type\":\"uint24\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"subID\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasAllowance\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"weiPerUnitLink\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"}],\"name\":\"RandomnessFulfillmentRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"nextBeaconOutputHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confDelay\",\"type\":\"uint24\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"subID\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"}],\"name\":\"RandomnessRequested\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"NUM_CONF_DELAYS\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"fail\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"reqID\",\"type\":\"uint48\"}],\"name\":\"getFulfillmentDurationByRequestID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_beaconPeriodBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pendingRequests\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID[]\",\"name\":\"\",\"type\":\"uint48[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_ReceivedRandomnessByRequestID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_arguments\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_averageFulfillmentInMillions\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_fastestFulfillment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"}],\"name\":\"s_fulfillmentDurationInBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_gasAvailable\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"}],\"name\":\"s_myBeaconRequests\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.SlotNumber\",\"name\":\"slotNumber\",\"type\":\"uint32\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelay\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_randomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_requestIDs\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"}],\"name\":\"s_requestOutputHeights\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"\",\"type\":\"uint24\"}],\"name\":\"s_requestsIDs\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_resetCounter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_slowestFulfillment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_slowestRequestID\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_subId\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_totalFulfilled\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_totalRequests\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"shouldFail\",\"type\":\"bool\"}],\"name\":\"setFail\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"reqId\",\"type\":\"uint48\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"delay\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"}],\"name\":\"storeBeaconRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"}],\"name\":\"testRedeemRandomness\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"uint64\",\"name\":\"subID\",\"type\":\"uint64\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelayArg\",\"type\":\"uint24\"}],\"name\":\"testRequestRandomness\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subID\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelayArg\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"}],\"name\":\"testRequestRandomnessFulfillment\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subID\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelayArg\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"batchSize\",\"type\":\"uint256\"}],\"name\":\"testRequestRandomnessFulfillmentBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040526000600a556000600b556103e7600c556000600d556000600e556000600f5534801561002f57600080fd5b5060405161181c38038061181c83398101604081905261004e9161008b565b6001600160a01b03929092166080819052600580546001600160a01b03191690911790556008805460ff19169115159190911790556009556100df565b6000806000606084860312156100a057600080fd5b83516001600160a01b03811681146100b757600080fd5b602085015190935080151581146100cd57600080fd5b80925050604084015190509250925092565b6080516117226100fa60003960006108fd01526117226000f3fe608060405234801561001057600080fd5b50600436106101da5760003560e01c806375baff7a116101045780639d769402116100a2578063f08c5daa11610071578063f08c5daa146104ff578063f339c79414610508578063f6eaffc81461051b578063fc7fea371461052e57600080fd5b80639d769402146104b0578063a9cc4718146104d1578063cd0593df146104ee578063d826f88f146104f757600080fd5b80638662aa3e116100de5780638662aa3e1461044c5780638866c6bd146104825780638d0e31651461048b5780639c9cd0151461049d57600080fd5b806375baff7a146103f95780637716cdaa146104245780637ab128ea1461043957600080fd5b80635a1c532b1161017c5780636d162a3e1161014b5780636d162a3e1461038f578063706da1ca146103a2578063737144bc146103e757806374dba124146103f057600080fd5b80635a1c532b1461032d5780635a47dd71146103405780635f15cccc14610353578063601201d31461038657600080fd5b80633d8b70aa116101b85780633d8b70aa146102ae5780633f7d43bd146102c357806345907626146102ee5780634a0aee291461031857600080fd5b80631757f11c146101df57806319a5fa22146101fb5780632f7527cc14610294575b600080fd5b6101e8600b5481565b6040519081526020015b60405180910390f35b610257610209366004610fcf565b60026020526000908152604090205463ffffffff811690640100000000810462ffffff1690670100000000000000810461ffff1690690100000000000000000090046001600160a01b031684565b6040805163ffffffff909516855262ffffff909316602085015261ffff909116918301919091526001600160a01b031660608201526080016101f2565b61029c600881565b60405160ff90911681526020016101f2565b6102c16102bc366004610fcf565b610537565b005b6101e86102d1366004610ff3565b601260209081526000928352604080842090915290825290205481565b6103016102fc366004611130565b6105fd565b60405165ffffffffffff90911681526020016101f2565b6103206106ed565b6040516101f291906111b0565b6102c161033b3660046111fc565b6107f8565b6102c161034e3660046112a8565b6108fb565b61030161036136600461137a565b600160209081526000928352604080842090915290825290205465ffffffffffff1681565b6101e8600e5481565b6102c161039d3660046113a6565b610983565b6006546103ce9074010000000000000000000000000000000000000000900467ffffffffffffffff1681565b60405167ffffffffffffffff90911681526020016101f2565b6101e8600a5481565b6101e8600c5481565b6101e8610407366004610ff3565b601360209081526000928352604080842090915290825290205481565b61042c610ab6565b6040516101f29190611442565b610301610447366004611455565b610b44565b6101e861045a366004610fcf565b600f54600090815260136020908152604080832065ffffffffffff9094168352929052205490565b6101e8600d5481565b6010546103019065ffffffffffff1681565b6103016104ab36600461146e565b610b80565b6102c16104be3660046114b1565b6008805460ff1916911515919091179055565b6008546104de9060ff1681565b60405190151581526020016101f2565b6101e860095481565b6102c1610c80565b6101e860075481565b6101e86105163660046114d3565b610ccd565b6101e8610529366004611455565b610cfe565b6101e8600f5481565b6005546040517f74d8461100000000000000000000000000000000000000000000000000000000815265ffffffffffff831660048201526000916001600160a01b0316906374d84611906024016000604051808303816000875af11580156105a3573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526105cb91908101906114ff565b65ffffffffffff8316600090815260036020908152604090912082519293506105f8929091840190610f33565b505050565b60008060095461060b610d1f565b61061591906115a6565b9050600081600954610625610d1f565b61062f91906115d0565b61063991906115e8565b6005546040517ff645dcb10000000000000000000000000000000000000000000000000000000081529192506000916001600160a01b039091169063f645dcb190610690908c908c908c908c908c906004016115ff565b6020604051808303816000875af11580156106af573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106d39190611643565b90506106e18183898b610983565b98975050505050505050565b60115460609060009067ffffffffffffffff81111561070e5761070e611079565b604051908082528060200260200182016040528015610737578160200160208202803683370190505b5090506000805b6011548110156107f05760006011828154811061075d5761075d611660565b60009182526020808320600580840490910154600f548552601383526040808620929094066006026101000a900465ffffffffffff1680855291529082205490925090036107dd57808484815181106107b8576107b8611660565b65ffffffffffff90921660209283029190910190910152826107d981611676565b9350505b50806107e881611676565b91505061073e565b508152919050565b6000600954610805610d1f565b61080f91906115a6565b905060008160095461081f610d1f565b61082991906115d0565b61083391906115e8565b905060005b838110156108f057600061084f8a8a8a8a8a6105fd565b600d8054919250600061086183611676565b9091555050600f54600090815260126020908152604080832065ffffffffffff948516808552925282208590556011805460018101825592527f31ecc21a745e3968a04e9570e4425bc18fa8019c68028196b546d1669c200c686005808404919091018054919093066006026101000a93840219169202919091179055806108e881611676565b915050610838565b505050505050505050565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031633146109785760405162461bcd60e51b815260206004820152601c60248201527f6f6e6c7920636f6f7264696e61746f722063616e2066756c66696c6c0000000060448201526064015b60405180910390fd5b6105f8838383610da9565b600083815260016020908152604080832062ffffff861684529091528120805465ffffffffffff191665ffffffffffff87161790556009546109c5908561168f565b6040805160808101825263ffffffff928316815262ffffff958616602080830191825261ffff968716838501908152306060850190815265ffffffffffff909b1660009081526002909252939020915182549151935199516001600160a01b03166901000000000000000000027fffffff0000000000000000000000000000000000000000ffffffffffffffffff9a90971667010000000000000002999099167fffffff00000000000000000000000000000000000000000000ffffffffffffff939097166401000000000266ffffffffffffff199091169890931697909717919091171692909217179092555050565b60048054610ac3906116a3565b80601f0160208091040260200160405190810160405280929190818152602001828054610aef906116a3565b8015610b3c5780601f10610b1157610100808354040283529160200191610b3c565b820191906000526020600020905b815481529060010190602001808311610b1f57829003601f168201915b505050505081565b60118181548110610b5457600080fd5b9060005260206000209060059182820401919006600602915054906101000a900465ffffffffffff1681565b600080600954610b8e610d1f565b610b9891906115a6565b9050600081600954610ba8610d1f565b610bb291906115d0565b610bbc91906115e8565b6005546040517fdc92accf00000000000000000000000000000000000000000000000000000000815261ffff8916600482015267ffffffffffffffff8816602482015262ffffff871660448201529192506000916001600160a01b039091169063dc92accf906064016020604051808303816000875af1158015610c44573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c689190611643565b9050610c768183878a610983565b9695505050505050565b6000600a819055600b8190556103e7600c55600d819055600e8190556010805465ffffffffffff19169055600f805491610cb983611676565b90915550610ccb905060116000610f7e565b565b60036020528160005260406000208181548110610ce957600080fd5b90600052602060002001600091509150505481565b60008181548110610d0e57600080fd5b600091825260209091200154905081565b60004661a4b1811480610d34575062066eed81145b15610da25760646001600160a01b031663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610d78573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d9c91906116dd565b91505090565b4391505090565b60085460ff1615610dfc5760405162461bcd60e51b815260206004820152601d60248201527f206661696c656420696e2066756c66696c6c52616e646f6d576f726473000000604482015260640161096f565b65ffffffffffff831660009081526003602090815260409091208351610e2492850190610f33565b50600f54600090815260126020908152604080832065ffffffffffff87168452909152812054610e52610d1f565b610e5c91906115e8565b90506000610e6d82620f42406116f6565b9050600b54821115610e9857600b8290556010805465ffffffffffff191665ffffffffffff87161790555b600c548210610ea957600c54610eab565b815b600c55600e54610ebb5780610eee565b600e54610ec99060016115d0565b81600e54600a54610eda91906116f6565b610ee491906115d0565b610eee919061168f565b600a55600e8054906000610f0183611676565b9091555050600f54600090815260136020908152604080832065ffffffffffff90981683529690529490942055505050565b828054828255906000526020600020908101928215610f6e579160200282015b82811115610f6e578251825591602001919060010190610f53565b50610f7a929150610fa6565b5090565b508054600082556004016005900490600052602060002090810190610fa39190610fa6565b50565b5b80821115610f7a5760008155600101610fa7565b65ffffffffffff81168114610fa357600080fd5b600060208284031215610fe157600080fd5b8135610fec81610fbb565b9392505050565b6000806040838503121561100657600080fd5b82359150602083013561101881610fbb565b809150509250929050565b803567ffffffffffffffff8116811461103b57600080fd5b919050565b803561ffff8116811461103b57600080fd5b803562ffffff8116811461103b57600080fd5b803563ffffffff8116811461103b57600080fd5b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff811182821017156110b8576110b8611079565b604052919050565b600082601f8301126110d157600080fd5b813567ffffffffffffffff8111156110eb576110eb611079565b6110fe601f8201601f191660200161108f565b81815284602083860101111561111357600080fd5b816020850160208301376000918101602001919091529392505050565b600080600080600060a0868803121561114857600080fd5b61115186611023565b945061115f60208701611040565b935061116d60408701611052565b925061117b60608701611065565b9150608086013567ffffffffffffffff81111561119757600080fd5b6111a3888289016110c0565b9150509295509295909350565b6020808252825182820181905260009190848201906040850190845b818110156111f057835165ffffffffffff16835292840192918401916001016111cc565b50909695505050505050565b60008060008060008060c0878903121561121557600080fd5b61121e87611023565b955061122c60208801611040565b945061123a60408801611052565b935061124860608801611065565b9250608087013567ffffffffffffffff81111561126457600080fd5b61127089828a016110c0565b92505060a087013590509295509295509295565b600067ffffffffffffffff82111561129e5761129e611079565b5060051b60200190565b6000806000606084860312156112bd57600080fd5b83356112c881610fbb565b925060208481013567ffffffffffffffff808211156112e657600080fd5b818701915087601f8301126112fa57600080fd5b813561130d61130882611284565b61108f565b81815260059190911b8301840190848101908a83111561132c57600080fd5b938501935b8285101561134a57843582529385019390850190611331565b96505050604087013592508083111561136257600080fd5b5050611370868287016110c0565b9150509250925092565b6000806040838503121561138d57600080fd5b8235915061139d60208401611052565b90509250929050565b600080600080608085870312156113bc57600080fd5b84356113c781610fbb565b9350602085013592506113dc60408601611052565b91506113ea60608601611040565b905092959194509250565b6000815180845260005b8181101561141b576020818501810151868301820152016113ff565b8181111561142d576000602083870101525b50601f01601f19169290920160200192915050565b602081526000610fec60208301846113f5565b60006020828403121561146757600080fd5b5035919050565b60008060006060848603121561148357600080fd5b61148c84611040565b925061149a60208501611023565b91506114a860408501611052565b90509250925092565b6000602082840312156114c357600080fd5b81358015158114610fec57600080fd5b600080604083850312156114e657600080fd5b82356114f181610fbb565b946020939093013593505050565b6000602080838503121561151257600080fd5b825167ffffffffffffffff81111561152957600080fd5b8301601f8101851361153a57600080fd5b805161154861130882611284565b81815260059190911b8201830190838101908783111561156757600080fd5b928401925b828410156115855783518252928401929084019061156c565b979650505050505050565b634e487b7160e01b600052601260045260246000fd5b6000826115b5576115b5611590565b500690565b634e487b7160e01b600052601160045260246000fd5b600082198211156115e3576115e36115ba565b500190565b6000828210156115fa576115fa6115ba565b500390565b67ffffffffffffffff8616815261ffff8516602082015262ffffff8416604082015263ffffffff8316606082015260a06080820152600061158560a08301846113f5565b60006020828403121561165557600080fd5b8151610fec81610fbb565b634e487b7160e01b600052603260045260246000fd5b600060018201611688576116886115ba565b5060010190565b60008261169e5761169e611590565b500490565b600181811c908216806116b757607f821691505b6020821081036116d757634e487b7160e01b600052602260045260246000fd5b50919050565b6000602082840312156116ef57600080fd5b5051919050565b6000816000190483118215151615611710576117106115ba565b50029056fea164736f6c634300080f000a",
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

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCaller) SRequestIDs(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTestBeaconVRFConsumer.contract.Call(opts, &out, "s_requestIDs", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerSession) SRequestIDs(arg0 *big.Int) (*big.Int, error) {
	return _LoadTestBeaconVRFConsumer.Contract.SRequestIDs(&_LoadTestBeaconVRFConsumer.CallOpts, arg0)
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerCallerSession) SRequestIDs(arg0 *big.Int) (*big.Int, error) {
	return _LoadTestBeaconVRFConsumer.Contract.SRequestIDs(&_LoadTestBeaconVRFConsumer.CallOpts, arg0)
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

type LoadTestBeaconVRFConsumerNewTransmissionIterator struct {
	Event *LoadTestBeaconVRFConsumerNewTransmission

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LoadTestBeaconVRFConsumerNewTransmissionIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LoadTestBeaconVRFConsumerNewTransmission)
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
		it.Event = new(LoadTestBeaconVRFConsumerNewTransmission)
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

func (it *LoadTestBeaconVRFConsumerNewTransmissionIterator) Error() error {
	return it.fail
}

func (it *LoadTestBeaconVRFConsumerNewTransmissionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LoadTestBeaconVRFConsumerNewTransmission struct {
	AggregatorRoundId uint32
	EpochAndRound     *big.Int
	Transmitter       common.Address
	JuelsPerFeeCoin   *big.Int
	ConfigDigest      [32]byte
	Raw               types.Log
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerFilterer) FilterNewTransmission(opts *bind.FilterOpts, aggregatorRoundId []uint32, epochAndRound []*big.Int) (*LoadTestBeaconVRFConsumerNewTransmissionIterator, error) {

	var aggregatorRoundIdRule []interface{}
	for _, aggregatorRoundIdItem := range aggregatorRoundId {
		aggregatorRoundIdRule = append(aggregatorRoundIdRule, aggregatorRoundIdItem)
	}
	var epochAndRoundRule []interface{}
	for _, epochAndRoundItem := range epochAndRound {
		epochAndRoundRule = append(epochAndRoundRule, epochAndRoundItem)
	}

	logs, sub, err := _LoadTestBeaconVRFConsumer.contract.FilterLogs(opts, "NewTransmission", aggregatorRoundIdRule, epochAndRoundRule)
	if err != nil {
		return nil, err
	}
	return &LoadTestBeaconVRFConsumerNewTransmissionIterator{contract: _LoadTestBeaconVRFConsumer.contract, event: "NewTransmission", logs: logs, sub: sub}, nil
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerFilterer) WatchNewTransmission(opts *bind.WatchOpts, sink chan<- *LoadTestBeaconVRFConsumerNewTransmission, aggregatorRoundId []uint32, epochAndRound []*big.Int) (event.Subscription, error) {

	var aggregatorRoundIdRule []interface{}
	for _, aggregatorRoundIdItem := range aggregatorRoundId {
		aggregatorRoundIdRule = append(aggregatorRoundIdRule, aggregatorRoundIdItem)
	}
	var epochAndRoundRule []interface{}
	for _, epochAndRoundItem := range epochAndRound {
		epochAndRoundRule = append(epochAndRoundRule, epochAndRoundItem)
	}

	logs, sub, err := _LoadTestBeaconVRFConsumer.contract.WatchLogs(opts, "NewTransmission", aggregatorRoundIdRule, epochAndRoundRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LoadTestBeaconVRFConsumerNewTransmission)
				if err := _LoadTestBeaconVRFConsumer.contract.UnpackLog(event, "NewTransmission", log); err != nil {
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

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerFilterer) ParseNewTransmission(log types.Log) (*LoadTestBeaconVRFConsumerNewTransmission, error) {
	event := new(LoadTestBeaconVRFConsumerNewTransmission)
	if err := _LoadTestBeaconVRFConsumer.contract.UnpackLog(event, "NewTransmission", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LoadTestBeaconVRFConsumerOutputsServedIterator struct {
	Event *LoadTestBeaconVRFConsumerOutputsServed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LoadTestBeaconVRFConsumerOutputsServedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LoadTestBeaconVRFConsumerOutputsServed)
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
		it.Event = new(LoadTestBeaconVRFConsumerOutputsServed)
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

func (it *LoadTestBeaconVRFConsumerOutputsServedIterator) Error() error {
	return it.fail
}

func (it *LoadTestBeaconVRFConsumerOutputsServedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LoadTestBeaconVRFConsumerOutputsServed struct {
	RecentBlockHeight uint64
	Transmitter       common.Address
	JuelsPerFeeCoin   *big.Int
	OutputsServed     []VRFBeaconTypesOutputServed
	Raw               types.Log
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerFilterer) FilterOutputsServed(opts *bind.FilterOpts) (*LoadTestBeaconVRFConsumerOutputsServedIterator, error) {

	logs, sub, err := _LoadTestBeaconVRFConsumer.contract.FilterLogs(opts, "OutputsServed")
	if err != nil {
		return nil, err
	}
	return &LoadTestBeaconVRFConsumerOutputsServedIterator{contract: _LoadTestBeaconVRFConsumer.contract, event: "OutputsServed", logs: logs, sub: sub}, nil
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerFilterer) WatchOutputsServed(opts *bind.WatchOpts, sink chan<- *LoadTestBeaconVRFConsumerOutputsServed) (event.Subscription, error) {

	logs, sub, err := _LoadTestBeaconVRFConsumer.contract.WatchLogs(opts, "OutputsServed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LoadTestBeaconVRFConsumerOutputsServed)
				if err := _LoadTestBeaconVRFConsumer.contract.UnpackLog(event, "OutputsServed", log); err != nil {
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

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerFilterer) ParseOutputsServed(log types.Log) (*LoadTestBeaconVRFConsumerOutputsServed, error) {
	event := new(LoadTestBeaconVRFConsumerOutputsServed)
	if err := _LoadTestBeaconVRFConsumer.contract.UnpackLog(event, "OutputsServed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LoadTestBeaconVRFConsumerRandomWordsFulfilledIterator struct {
	Event *LoadTestBeaconVRFConsumerRandomWordsFulfilled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LoadTestBeaconVRFConsumerRandomWordsFulfilledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LoadTestBeaconVRFConsumerRandomWordsFulfilled)
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
		it.Event = new(LoadTestBeaconVRFConsumerRandomWordsFulfilled)
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

func (it *LoadTestBeaconVRFConsumerRandomWordsFulfilledIterator) Error() error {
	return it.fail
}

func (it *LoadTestBeaconVRFConsumerRandomWordsFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LoadTestBeaconVRFConsumerRandomWordsFulfilled struct {
	RequestIDs            []*big.Int
	SuccessfulFulfillment []byte
	TruncatedErrorData    [][]byte
	Raw                   types.Log
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerFilterer) FilterRandomWordsFulfilled(opts *bind.FilterOpts) (*LoadTestBeaconVRFConsumerRandomWordsFulfilledIterator, error) {

	logs, sub, err := _LoadTestBeaconVRFConsumer.contract.FilterLogs(opts, "RandomWordsFulfilled")
	if err != nil {
		return nil, err
	}
	return &LoadTestBeaconVRFConsumerRandomWordsFulfilledIterator{contract: _LoadTestBeaconVRFConsumer.contract, event: "RandomWordsFulfilled", logs: logs, sub: sub}, nil
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerFilterer) WatchRandomWordsFulfilled(opts *bind.WatchOpts, sink chan<- *LoadTestBeaconVRFConsumerRandomWordsFulfilled) (event.Subscription, error) {

	logs, sub, err := _LoadTestBeaconVRFConsumer.contract.WatchLogs(opts, "RandomWordsFulfilled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LoadTestBeaconVRFConsumerRandomWordsFulfilled)
				if err := _LoadTestBeaconVRFConsumer.contract.UnpackLog(event, "RandomWordsFulfilled", log); err != nil {
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

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerFilterer) ParseRandomWordsFulfilled(log types.Log) (*LoadTestBeaconVRFConsumerRandomWordsFulfilled, error) {
	event := new(LoadTestBeaconVRFConsumerRandomWordsFulfilled)
	if err := _LoadTestBeaconVRFConsumer.contract.UnpackLog(event, "RandomWordsFulfilled", log); err != nil {
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
	RequestID              *big.Int
	Requester              common.Address
	NextBeaconOutputHeight uint64
	ConfDelay              *big.Int
	SubID                  uint64
	NumWords               uint16
	GasAllowance           uint32
	GasPrice               *big.Int
	WeiPerUnitLink         *big.Int
	Arguments              []byte
	Raw                    types.Log
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerFilterer) FilterRandomnessFulfillmentRequested(opts *bind.FilterOpts, requestID []*big.Int, requester []common.Address) (*LoadTestBeaconVRFConsumerRandomnessFulfillmentRequestedIterator, error) {

	var requestIDRule []interface{}
	for _, requestIDItem := range requestID {
		requestIDRule = append(requestIDRule, requestIDItem)
	}
	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _LoadTestBeaconVRFConsumer.contract.FilterLogs(opts, "RandomnessFulfillmentRequested", requestIDRule, requesterRule)
	if err != nil {
		return nil, err
	}
	return &LoadTestBeaconVRFConsumerRandomnessFulfillmentRequestedIterator{contract: _LoadTestBeaconVRFConsumer.contract, event: "RandomnessFulfillmentRequested", logs: logs, sub: sub}, nil
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerFilterer) WatchRandomnessFulfillmentRequested(opts *bind.WatchOpts, sink chan<- *LoadTestBeaconVRFConsumerRandomnessFulfillmentRequested, requestID []*big.Int, requester []common.Address) (event.Subscription, error) {

	var requestIDRule []interface{}
	for _, requestIDItem := range requestID {
		requestIDRule = append(requestIDRule, requestIDItem)
	}
	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _LoadTestBeaconVRFConsumer.contract.WatchLogs(opts, "RandomnessFulfillmentRequested", requestIDRule, requesterRule)
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
	RequestID              *big.Int
	Requester              common.Address
	NextBeaconOutputHeight uint64
	ConfDelay              *big.Int
	SubID                  uint64
	NumWords               uint16
	Raw                    types.Log
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerFilterer) FilterRandomnessRequested(opts *bind.FilterOpts, requestID []*big.Int, requester []common.Address) (*LoadTestBeaconVRFConsumerRandomnessRequestedIterator, error) {

	var requestIDRule []interface{}
	for _, requestIDItem := range requestID {
		requestIDRule = append(requestIDRule, requestIDItem)
	}
	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _LoadTestBeaconVRFConsumer.contract.FilterLogs(opts, "RandomnessRequested", requestIDRule, requesterRule)
	if err != nil {
		return nil, err
	}
	return &LoadTestBeaconVRFConsumerRandomnessRequestedIterator{contract: _LoadTestBeaconVRFConsumer.contract, event: "RandomnessRequested", logs: logs, sub: sub}, nil
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumerFilterer) WatchRandomnessRequested(opts *bind.WatchOpts, sink chan<- *LoadTestBeaconVRFConsumerRandomnessRequested, requestID []*big.Int, requester []common.Address) (event.Subscription, error) {

	var requestIDRule []interface{}
	for _, requestIDItem := range requestID {
		requestIDRule = append(requestIDRule, requestIDItem)
	}
	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _LoadTestBeaconVRFConsumer.contract.WatchLogs(opts, "RandomnessRequested", requestIDRule, requesterRule)
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
	case _LoadTestBeaconVRFConsumer.abi.Events["NewTransmission"].ID:
		return _LoadTestBeaconVRFConsumer.ParseNewTransmission(log)
	case _LoadTestBeaconVRFConsumer.abi.Events["OutputsServed"].ID:
		return _LoadTestBeaconVRFConsumer.ParseOutputsServed(log)
	case _LoadTestBeaconVRFConsumer.abi.Events["RandomWordsFulfilled"].ID:
		return _LoadTestBeaconVRFConsumer.ParseRandomWordsFulfilled(log)
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

func (LoadTestBeaconVRFConsumerNewTransmission) Topic() common.Hash {
	return common.HexToHash("0x0e21a6118d05811033e5ee111f22acf5ffd1d1406aa2200fe054f29458f8fb26")
}

func (LoadTestBeaconVRFConsumerOutputsServed) Topic() common.Hash {
	return common.HexToHash("0xb0ab03ecc205672b29a9a88df77d28c868016f1b9448275c78a245327afa71cc")
}

func (LoadTestBeaconVRFConsumerRandomWordsFulfilled) Topic() common.Hash {
	return common.HexToHash("0x47ddf7bb0cbd94c1b43c5097f1352a80db0ceb3696f029d32b24f32cd631d2b7")
}

func (LoadTestBeaconVRFConsumerRandomnessFulfillmentRequested) Topic() common.Hash {
	return common.HexToHash("0xda7e52094a2e0200a38ad6b72d1ce40dcc840bbc4cf95ccf8065ef521da664dc")
}

func (LoadTestBeaconVRFConsumerRandomnessRequested) Topic() common.Hash {
	return common.HexToHash("0xdcbc2f1d67e18a0e5cf023ad14515e21bfc7e0cccb7ddb1011f46baffc7d9d15")
}

func (_LoadTestBeaconVRFConsumer *LoadTestBeaconVRFConsumer) Address() common.Address {
	return _LoadTestBeaconVRFConsumer.address
}

type LoadTestBeaconVRFConsumerInterface interface {
	NUMCONFDELAYS(opts *bind.CallOpts) (uint8, error)

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

	SMyBeaconRequests(opts *bind.CallOpts, arg0 *big.Int) (SMyBeaconRequests,

		error)

	SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	SRequestIDs(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

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

	TestRedeemRandomness(opts *bind.TransactOpts, requestID *big.Int) (*types.Transaction, error)

	TestRequestRandomness(opts *bind.TransactOpts, numWords uint16, subID uint64, confirmationDelayArg *big.Int) (*types.Transaction, error)

	TestRequestRandomnessFulfillment(opts *bind.TransactOpts, subID uint64, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte) (*types.Transaction, error)

	TestRequestRandomnessFulfillmentBatch(opts *bind.TransactOpts, subID uint64, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte, batchSize *big.Int) (*types.Transaction, error)

	FilterConfigSet(opts *bind.FilterOpts) (*LoadTestBeaconVRFConsumerConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *LoadTestBeaconVRFConsumerConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*LoadTestBeaconVRFConsumerConfigSet, error)

	FilterNewTransmission(opts *bind.FilterOpts, aggregatorRoundId []uint32, epochAndRound []*big.Int) (*LoadTestBeaconVRFConsumerNewTransmissionIterator, error)

	WatchNewTransmission(opts *bind.WatchOpts, sink chan<- *LoadTestBeaconVRFConsumerNewTransmission, aggregatorRoundId []uint32, epochAndRound []*big.Int) (event.Subscription, error)

	ParseNewTransmission(log types.Log) (*LoadTestBeaconVRFConsumerNewTransmission, error)

	FilterOutputsServed(opts *bind.FilterOpts) (*LoadTestBeaconVRFConsumerOutputsServedIterator, error)

	WatchOutputsServed(opts *bind.WatchOpts, sink chan<- *LoadTestBeaconVRFConsumerOutputsServed) (event.Subscription, error)

	ParseOutputsServed(log types.Log) (*LoadTestBeaconVRFConsumerOutputsServed, error)

	FilterRandomWordsFulfilled(opts *bind.FilterOpts) (*LoadTestBeaconVRFConsumerRandomWordsFulfilledIterator, error)

	WatchRandomWordsFulfilled(opts *bind.WatchOpts, sink chan<- *LoadTestBeaconVRFConsumerRandomWordsFulfilled) (event.Subscription, error)

	ParseRandomWordsFulfilled(log types.Log) (*LoadTestBeaconVRFConsumerRandomWordsFulfilled, error)

	FilterRandomnessFulfillmentRequested(opts *bind.FilterOpts, requestID []*big.Int, requester []common.Address) (*LoadTestBeaconVRFConsumerRandomnessFulfillmentRequestedIterator, error)

	WatchRandomnessFulfillmentRequested(opts *bind.WatchOpts, sink chan<- *LoadTestBeaconVRFConsumerRandomnessFulfillmentRequested, requestID []*big.Int, requester []common.Address) (event.Subscription, error)

	ParseRandomnessFulfillmentRequested(log types.Log) (*LoadTestBeaconVRFConsumerRandomnessFulfillmentRequested, error)

	FilterRandomnessRequested(opts *bind.FilterOpts, requestID []*big.Int, requester []common.Address) (*LoadTestBeaconVRFConsumerRandomnessRequestedIterator, error)

	WatchRandomnessRequested(opts *bind.WatchOpts, sink chan<- *LoadTestBeaconVRFConsumerRandomnessRequested, requestID []*big.Int, requester []common.Address) (event.Subscription, error)

	ParseRandomnessRequested(log types.Log) (*LoadTestBeaconVRFConsumerRandomnessRequested, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
