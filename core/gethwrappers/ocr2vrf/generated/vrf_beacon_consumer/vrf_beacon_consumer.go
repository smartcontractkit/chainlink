// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_beacon_consumer

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

var BeaconVRFConsumerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"shouldFail\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"beaconPeriodBlocks\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"aggregatorRoundId\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"uint40\",\"name\":\"epochAndRound\",\"type\":\"uint40\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint192\",\"name\":\"juelsPerFeeCoin\",\"type\":\"uint192\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"reasonableGasPrice\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"NewTransmission\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"recentBlockHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint192\",\"name\":\"juelsPerFeeCoin\",\"type\":\"uint192\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"reasonableGasPrice\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelay\",\"type\":\"uint24\"},{\"internalType\":\"uint256\",\"name\":\"proofG1X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"proofG1Y\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structVRFBeaconTypes.OutputServed[]\",\"name\":\"outputsServed\",\"type\":\"tuple[]\"}],\"name\":\"OutputsServed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"VRFBeaconTypes.RequestID[]\",\"name\":\"requestIDs\",\"type\":\"uint48[]\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"successfulFulfillment\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes[]\",\"name\":\"truncatedErrorData\",\"type\":\"bytes[]\"}],\"name\":\"RandomWordsFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"nextBeaconOutputHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confDelay\",\"type\":\"uint24\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"subID\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasAllowance\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"weiPerUnitLink\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"}],\"name\":\"RandomnessFulfillmentRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"nextBeaconOutputHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confDelay\",\"type\":\"uint24\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"subID\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"}],\"name\":\"RandomnessRequested\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"NUM_CONF_DELAYS\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"fail\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_beaconPeriodBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_ReceivedRandomnessByRequestID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_arguments\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_gasAvailable\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"}],\"name\":\"s_myBeaconRequests\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.SlotNumber\",\"name\":\"slotNumber\",\"type\":\"uint32\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelay\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_randomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"\",\"type\":\"uint24\"}],\"name\":\"s_requestsIDs\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_subId\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"shouldFail\",\"type\":\"bool\"}],\"name\":\"setFail\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"reqId\",\"type\":\"uint48\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"delay\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"}],\"name\":\"storeBeaconRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"}],\"name\":\"testRedeemRandomness\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"uint64\",\"name\":\"subID\",\"type\":\"uint64\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelayArg\",\"type\":\"uint24\"}],\"name\":\"testRequestRandomness\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subID\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelayArg\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"}],\"name\":\"testRequestRandomnessFulfillment\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561001057600080fd5b5060405161120738038061120783398101604081905261002f9161006c565b6001600160a01b03929092166080819052600580546001600160a01b03191690911790556008805460ff19169115159190911790556009556100c0565b60008060006060848603121561008157600080fd5b83516001600160a01b038116811461009857600080fd5b602085015190935080151581146100ae57600080fd5b80925050604084015190509250925092565b60805161112c6100db60003960006104fe015261112c6000f3fe608060405234801561001057600080fd5b50600436106101005760003560e01c80637716cdaa11610097578063cd0593df11610066578063cd0593df14610300578063f08c5daa14610317578063f339c79414610320578063f6eaffc81461033357600080fd5b80637716cdaa1461029a5780639c9cd015146102af5780639d769402146102c2578063a9cc4718146102e357600080fd5b80635a47dd71116100d35780635a47dd71146101fc5780635f15cccc1461020f5780636d162a3e14610242578063706da1ca1461025557600080fd5b806319a5fa22146101055780632f7527cc146101a35780633d8b70aa146101bd57806345907626146101d2575b600080fd5b610161610113366004610a26565b60026020526000908152604090205463ffffffff811690640100000000810462ffffff1690670100000000000000810461ffff1690690100000000000000000090046001600160a01b031684565b6040805163ffffffff909516855262ffffff909316602085015261ffff909116918301919091526001600160a01b031660608201526080015b60405180910390f35b6101ab600881565b60405160ff909116815260200161019a565b6101d06101cb366004610a26565b610346565b005b6101e56101e0366004610b43565b61040c565b60405165ffffffffffff909116815260200161019a565b6101d061020a366004610bf2565b6104fc565b6101e561021d366004610cc4565b600160209081526000928352604080842090915290825290205465ffffffffffff1681565b6101d0610250366004610cf0565b610584565b6006546102819074010000000000000000000000000000000000000000900467ffffffffffffffff1681565b60405167ffffffffffffffff909116815260200161019a565b6102a26106b7565b60405161019a9190610d8c565b6101e56102bd366004610d9f565b610745565b6101d06102d0366004610de2565b6008805460ff1916911515919091179055565b6008546102f09060ff1681565b604051901515815260200161019a565b61030960095481565b60405190815260200161019a565b61030960075481565b61030961032e366004610e04565b610845565b610309610341366004610e30565b610876565b6005546040517f74d8461100000000000000000000000000000000000000000000000000000000815265ffffffffffff831660048201526000916001600160a01b0316906374d84611906024016000604051808303816000875af11580156103b2573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526103da9190810190610e49565b65ffffffffffff8316600090815260036020908152604090912082519293506104079290918401906109af565b505050565b60008060095461041a610897565b6104249190610ef0565b9050600081600954610434610897565b61043e9190610f1a565b6104489190610f32565b6005546040517ff645dcb10000000000000000000000000000000000000000000000000000000081529192506000916001600160a01b039091169063f645dcb19061049f908c908c908c908c908c90600401610f49565b6020604051808303816000875af11580156104be573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104e29190610f8d565b90506104f08183898b610584565b98975050505050505050565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031633146105795760405162461bcd60e51b815260206004820152601c60248201527f6f6e6c7920636f6f7264696e61746f722063616e2066756c66696c6c0000000060448201526064015b60405180910390fd5b610407838383610921565b600083815260016020908152604080832062ffffff861684529091528120805465ffffffffffff191665ffffffffffff87161790556009546105c69085610faa565b6040805160808101825263ffffffff928316815262ffffff958616602080830191825261ffff968716838501908152306060850190815265ffffffffffff909b1660009081526002909252939020915182549151935199516001600160a01b03166901000000000000000000027fffffff0000000000000000000000000000000000000000ffffffffffffffffff9a90971667010000000000000002999099167fffffff00000000000000000000000000000000000000000000ffffffffffffff939097166401000000000266ffffffffffffff199091169890931697909717919091171692909217179092555050565b600480546106c490610fbe565b80601f01602080910402602001604051908101604052809291908181526020018280546106f090610fbe565b801561073d5780601f106107125761010080835404028352916020019161073d565b820191906000526020600020905b81548152906001019060200180831161072057829003601f168201915b505050505081565b600080600954610753610897565b61075d9190610ef0565b905060008160095461076d610897565b6107779190610f1a565b6107819190610f32565b6005546040517fdc92accf00000000000000000000000000000000000000000000000000000000815261ffff8916600482015267ffffffffffffffff8816602482015262ffffff871660448201529192506000916001600160a01b039091169063dc92accf906064016020604051808303816000875af1158015610809573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061082d9190610f8d565b905061083b8183878a610584565b9695505050505050565b6003602052816000526040600020818154811061086157600080fd5b90600052602060002001600091509150505481565b6000818154811061088657600080fd5b600091825260209091200154905081565b60004661a4b18114806108ac575062066eed81145b1561091a5760646001600160a01b031663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156108f0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906109149190610ff8565b91505090565b4391505090565b60085460ff16156109745760405162461bcd60e51b815260206004820152601d60248201527f206661696c656420696e2066756c66696c6c52616e646f6d576f7264730000006044820152606401610570565b65ffffffffffff83166000908152600360209081526040909120835161099c928501906109af565b5060046109a9828261105f565b50505050565b8280548282559060005260206000209081019282156109ea579160200282015b828111156109ea5782518255916020019190600101906109cf565b506109f69291506109fa565b5090565b5b808211156109f657600081556001016109fb565b65ffffffffffff81168114610a2357600080fd5b50565b600060208284031215610a3857600080fd5b8135610a4381610a0f565b9392505050565b803567ffffffffffffffff81168114610a6257600080fd5b919050565b803561ffff81168114610a6257600080fd5b803562ffffff81168114610a6257600080fd5b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff81118282101715610acb57610acb610a8c565b604052919050565b600082601f830112610ae457600080fd5b813567ffffffffffffffff811115610afe57610afe610a8c565b610b11601f8201601f1916602001610aa2565b818152846020838601011115610b2657600080fd5b816020850160208301376000918101602001919091529392505050565b600080600080600060a08688031215610b5b57600080fd5b610b6486610a4a565b9450610b7260208701610a67565b9350610b8060408701610a79565b9250606086013563ffffffff81168114610b9957600080fd5b9150608086013567ffffffffffffffff811115610bb557600080fd5b610bc188828901610ad3565b9150509295509295909350565b600067ffffffffffffffff821115610be857610be8610a8c565b5060051b60200190565b600080600060608486031215610c0757600080fd5b8335610c1281610a0f565b925060208481013567ffffffffffffffff80821115610c3057600080fd5b818701915087601f830112610c4457600080fd5b8135610c57610c5282610bce565b610aa2565b81815260059190911b8301840190848101908a831115610c7657600080fd5b938501935b82851015610c9457843582529385019390850190610c7b565b965050506040870135925080831115610cac57600080fd5b5050610cba86828701610ad3565b9150509250925092565b60008060408385031215610cd757600080fd5b82359150610ce760208401610a79565b90509250929050565b60008060008060808587031215610d0657600080fd5b8435610d1181610a0f565b935060208501359250610d2660408601610a79565b9150610d3460608601610a67565b905092959194509250565b6000815180845260005b81811015610d6557602081850181015186830182015201610d49565b81811115610d77576000602083870101525b50601f01601f19169290920160200192915050565b602081526000610a436020830184610d3f565b600080600060608486031215610db457600080fd5b610dbd84610a67565b9250610dcb60208501610a4a565b9150610dd960408501610a79565b90509250925092565b600060208284031215610df457600080fd5b81358015158114610a4357600080fd5b60008060408385031215610e1757600080fd5b8235610e2281610a0f565b946020939093013593505050565b600060208284031215610e4257600080fd5b5035919050565b60006020808385031215610e5c57600080fd5b825167ffffffffffffffff811115610e7357600080fd5b8301601f81018513610e8457600080fd5b8051610e92610c5282610bce565b81815260059190911b82018301908381019087831115610eb157600080fd5b928401925b82841015610ecf57835182529284019290840190610eb6565b979650505050505050565b634e487b7160e01b600052601260045260246000fd5b600082610eff57610eff610eda565b500690565b634e487b7160e01b600052601160045260246000fd5b60008219821115610f2d57610f2d610f04565b500190565b600082821015610f4457610f44610f04565b500390565b67ffffffffffffffff8616815261ffff8516602082015262ffffff8416604082015263ffffffff8316606082015260a060808201526000610ecf60a0830184610d3f565b600060208284031215610f9f57600080fd5b8151610a4381610a0f565b600082610fb957610fb9610eda565b500490565b600181811c90821680610fd257607f821691505b602082108103610ff257634e487b7160e01b600052602260045260246000fd5b50919050565b60006020828403121561100a57600080fd5b5051919050565b601f82111561040757600081815260208120601f850160051c810160208610156110385750805b601f850160051c820191505b8181101561105757828155600101611044565b505050505050565b815167ffffffffffffffff81111561107957611079610a8c565b61108d816110878454610fbe565b84611011565b602080601f8311600181146110c257600084156110aa5750858301515b600019600386901b1c1916600185901b178555611057565b600085815260208120601f198616915b828110156110f1578886015182559484019460019091019084016110d2565b508582101561110f5787850151600019600388901b60f8161c191681555b5050505050600190811b0190555056fea164736f6c634300080f000a",
}

var BeaconVRFConsumerABI = BeaconVRFConsumerMetaData.ABI

var BeaconVRFConsumerBin = BeaconVRFConsumerMetaData.Bin

func DeployBeaconVRFConsumer(auth *bind.TransactOpts, backend bind.ContractBackend, coordinator common.Address, shouldFail bool, beaconPeriodBlocks *big.Int) (common.Address, *types.Transaction, *BeaconVRFConsumer, error) {
	parsed, err := BeaconVRFConsumerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BeaconVRFConsumerBin), backend, coordinator, shouldFail, beaconPeriodBlocks)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BeaconVRFConsumer{BeaconVRFConsumerCaller: BeaconVRFConsumerCaller{contract: contract}, BeaconVRFConsumerTransactor: BeaconVRFConsumerTransactor{contract: contract}, BeaconVRFConsumerFilterer: BeaconVRFConsumerFilterer{contract: contract}}, nil
}

type BeaconVRFConsumer struct {
	address common.Address
	abi     abi.ABI
	BeaconVRFConsumerCaller
	BeaconVRFConsumerTransactor
	BeaconVRFConsumerFilterer
}

type BeaconVRFConsumerCaller struct {
	contract *bind.BoundContract
}

type BeaconVRFConsumerTransactor struct {
	contract *bind.BoundContract
}

type BeaconVRFConsumerFilterer struct {
	contract *bind.BoundContract
}

type BeaconVRFConsumerSession struct {
	Contract     *BeaconVRFConsumer
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type BeaconVRFConsumerCallerSession struct {
	Contract *BeaconVRFConsumerCaller
	CallOpts bind.CallOpts
}

type BeaconVRFConsumerTransactorSession struct {
	Contract     *BeaconVRFConsumerTransactor
	TransactOpts bind.TransactOpts
}

type BeaconVRFConsumerRaw struct {
	Contract *BeaconVRFConsumer
}

type BeaconVRFConsumerCallerRaw struct {
	Contract *BeaconVRFConsumerCaller
}

type BeaconVRFConsumerTransactorRaw struct {
	Contract *BeaconVRFConsumerTransactor
}

func NewBeaconVRFConsumer(address common.Address, backend bind.ContractBackend) (*BeaconVRFConsumer, error) {
	abi, err := abi.JSON(strings.NewReader(BeaconVRFConsumerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindBeaconVRFConsumer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BeaconVRFConsumer{address: address, abi: abi, BeaconVRFConsumerCaller: BeaconVRFConsumerCaller{contract: contract}, BeaconVRFConsumerTransactor: BeaconVRFConsumerTransactor{contract: contract}, BeaconVRFConsumerFilterer: BeaconVRFConsumerFilterer{contract: contract}}, nil
}

func NewBeaconVRFConsumerCaller(address common.Address, caller bind.ContractCaller) (*BeaconVRFConsumerCaller, error) {
	contract, err := bindBeaconVRFConsumer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BeaconVRFConsumerCaller{contract: contract}, nil
}

func NewBeaconVRFConsumerTransactor(address common.Address, transactor bind.ContractTransactor) (*BeaconVRFConsumerTransactor, error) {
	contract, err := bindBeaconVRFConsumer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BeaconVRFConsumerTransactor{contract: contract}, nil
}

func NewBeaconVRFConsumerFilterer(address common.Address, filterer bind.ContractFilterer) (*BeaconVRFConsumerFilterer, error) {
	contract, err := bindBeaconVRFConsumer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BeaconVRFConsumerFilterer{contract: contract}, nil
}

func bindBeaconVRFConsumer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BeaconVRFConsumerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_BeaconVRFConsumer *BeaconVRFConsumerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BeaconVRFConsumer.Contract.BeaconVRFConsumerCaller.contract.Call(opts, result, method, params...)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.BeaconVRFConsumerTransactor.contract.Transfer(opts)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.BeaconVRFConsumerTransactor.contract.Transact(opts, method, params...)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BeaconVRFConsumer.Contract.contract.Call(opts, result, method, params...)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.contract.Transfer(opts)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.contract.Transact(opts, method, params...)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCaller) NUMCONFDELAYS(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _BeaconVRFConsumer.contract.Call(opts, &out, "NUM_CONF_DELAYS")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_BeaconVRFConsumer *BeaconVRFConsumerSession) NUMCONFDELAYS() (uint8, error) {
	return _BeaconVRFConsumer.Contract.NUMCONFDELAYS(&_BeaconVRFConsumer.CallOpts)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCallerSession) NUMCONFDELAYS() (uint8, error) {
	return _BeaconVRFConsumer.Contract.NUMCONFDELAYS(&_BeaconVRFConsumer.CallOpts)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCaller) Fail(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _BeaconVRFConsumer.contract.Call(opts, &out, "fail")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BeaconVRFConsumer *BeaconVRFConsumerSession) Fail() (bool, error) {
	return _BeaconVRFConsumer.Contract.Fail(&_BeaconVRFConsumer.CallOpts)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCallerSession) Fail() (bool, error) {
	return _BeaconVRFConsumer.Contract.Fail(&_BeaconVRFConsumer.CallOpts)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCaller) IBeaconPeriodBlocks(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BeaconVRFConsumer.contract.Call(opts, &out, "i_beaconPeriodBlocks")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_BeaconVRFConsumer *BeaconVRFConsumerSession) IBeaconPeriodBlocks() (*big.Int, error) {
	return _BeaconVRFConsumer.Contract.IBeaconPeriodBlocks(&_BeaconVRFConsumer.CallOpts)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCallerSession) IBeaconPeriodBlocks() (*big.Int, error) {
	return _BeaconVRFConsumer.Contract.IBeaconPeriodBlocks(&_BeaconVRFConsumer.CallOpts)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCaller) SReceivedRandomnessByRequestID(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _BeaconVRFConsumer.contract.Call(opts, &out, "s_ReceivedRandomnessByRequestID", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_BeaconVRFConsumer *BeaconVRFConsumerSession) SReceivedRandomnessByRequestID(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _BeaconVRFConsumer.Contract.SReceivedRandomnessByRequestID(&_BeaconVRFConsumer.CallOpts, arg0, arg1)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCallerSession) SReceivedRandomnessByRequestID(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _BeaconVRFConsumer.Contract.SReceivedRandomnessByRequestID(&_BeaconVRFConsumer.CallOpts, arg0, arg1)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCaller) SArguments(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _BeaconVRFConsumer.contract.Call(opts, &out, "s_arguments")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_BeaconVRFConsumer *BeaconVRFConsumerSession) SArguments() ([]byte, error) {
	return _BeaconVRFConsumer.Contract.SArguments(&_BeaconVRFConsumer.CallOpts)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCallerSession) SArguments() ([]byte, error) {
	return _BeaconVRFConsumer.Contract.SArguments(&_BeaconVRFConsumer.CallOpts)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCaller) SGasAvailable(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BeaconVRFConsumer.contract.Call(opts, &out, "s_gasAvailable")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_BeaconVRFConsumer *BeaconVRFConsumerSession) SGasAvailable() (*big.Int, error) {
	return _BeaconVRFConsumer.Contract.SGasAvailable(&_BeaconVRFConsumer.CallOpts)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCallerSession) SGasAvailable() (*big.Int, error) {
	return _BeaconVRFConsumer.Contract.SGasAvailable(&_BeaconVRFConsumer.CallOpts)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCaller) SMyBeaconRequests(opts *bind.CallOpts, arg0 *big.Int) (SMyBeaconRequests,

	error) {
	var out []interface{}
	err := _BeaconVRFConsumer.contract.Call(opts, &out, "s_myBeaconRequests", arg0)

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

func (_BeaconVRFConsumer *BeaconVRFConsumerSession) SMyBeaconRequests(arg0 *big.Int) (SMyBeaconRequests,

	error) {
	return _BeaconVRFConsumer.Contract.SMyBeaconRequests(&_BeaconVRFConsumer.CallOpts, arg0)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCallerSession) SMyBeaconRequests(arg0 *big.Int) (SMyBeaconRequests,

	error) {
	return _BeaconVRFConsumer.Contract.SMyBeaconRequests(&_BeaconVRFConsumer.CallOpts, arg0)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCaller) SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _BeaconVRFConsumer.contract.Call(opts, &out, "s_randomWords", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_BeaconVRFConsumer *BeaconVRFConsumerSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _BeaconVRFConsumer.Contract.SRandomWords(&_BeaconVRFConsumer.CallOpts, arg0)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCallerSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _BeaconVRFConsumer.Contract.SRandomWords(&_BeaconVRFConsumer.CallOpts, arg0)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCaller) SRequestsIDs(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _BeaconVRFConsumer.contract.Call(opts, &out, "s_requestsIDs", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_BeaconVRFConsumer *BeaconVRFConsumerSession) SRequestsIDs(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _BeaconVRFConsumer.Contract.SRequestsIDs(&_BeaconVRFConsumer.CallOpts, arg0, arg1)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCallerSession) SRequestsIDs(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _BeaconVRFConsumer.Contract.SRequestsIDs(&_BeaconVRFConsumer.CallOpts, arg0, arg1)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCaller) SSubId(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _BeaconVRFConsumer.contract.Call(opts, &out, "s_subId")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_BeaconVRFConsumer *BeaconVRFConsumerSession) SSubId() (uint64, error) {
	return _BeaconVRFConsumer.Contract.SSubId(&_BeaconVRFConsumer.CallOpts)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCallerSession) SSubId() (uint64, error) {
	return _BeaconVRFConsumer.Contract.SSubId(&_BeaconVRFConsumer.CallOpts)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestID *big.Int, randomWords []*big.Int, arguments []byte) (*types.Transaction, error) {
	return _BeaconVRFConsumer.contract.Transact(opts, "rawFulfillRandomWords", requestID, randomWords, arguments)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerSession) RawFulfillRandomWords(requestID *big.Int, randomWords []*big.Int, arguments []byte) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.RawFulfillRandomWords(&_BeaconVRFConsumer.TransactOpts, requestID, randomWords, arguments)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerTransactorSession) RawFulfillRandomWords(requestID *big.Int, randomWords []*big.Int, arguments []byte) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.RawFulfillRandomWords(&_BeaconVRFConsumer.TransactOpts, requestID, randomWords, arguments)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerTransactor) SetFail(opts *bind.TransactOpts, shouldFail bool) (*types.Transaction, error) {
	return _BeaconVRFConsumer.contract.Transact(opts, "setFail", shouldFail)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerSession) SetFail(shouldFail bool) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.SetFail(&_BeaconVRFConsumer.TransactOpts, shouldFail)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerTransactorSession) SetFail(shouldFail bool) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.SetFail(&_BeaconVRFConsumer.TransactOpts, shouldFail)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerTransactor) StoreBeaconRequest(opts *bind.TransactOpts, reqId *big.Int, height *big.Int, delay *big.Int, numWords uint16) (*types.Transaction, error) {
	return _BeaconVRFConsumer.contract.Transact(opts, "storeBeaconRequest", reqId, height, delay, numWords)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerSession) StoreBeaconRequest(reqId *big.Int, height *big.Int, delay *big.Int, numWords uint16) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.StoreBeaconRequest(&_BeaconVRFConsumer.TransactOpts, reqId, height, delay, numWords)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerTransactorSession) StoreBeaconRequest(reqId *big.Int, height *big.Int, delay *big.Int, numWords uint16) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.StoreBeaconRequest(&_BeaconVRFConsumer.TransactOpts, reqId, height, delay, numWords)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerTransactor) TestRedeemRandomness(opts *bind.TransactOpts, requestID *big.Int) (*types.Transaction, error) {
	return _BeaconVRFConsumer.contract.Transact(opts, "testRedeemRandomness", requestID)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerSession) TestRedeemRandomness(requestID *big.Int) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.TestRedeemRandomness(&_BeaconVRFConsumer.TransactOpts, requestID)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerTransactorSession) TestRedeemRandomness(requestID *big.Int) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.TestRedeemRandomness(&_BeaconVRFConsumer.TransactOpts, requestID)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerTransactor) TestRequestRandomness(opts *bind.TransactOpts, numWords uint16, subID uint64, confirmationDelayArg *big.Int) (*types.Transaction, error) {
	return _BeaconVRFConsumer.contract.Transact(opts, "testRequestRandomness", numWords, subID, confirmationDelayArg)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerSession) TestRequestRandomness(numWords uint16, subID uint64, confirmationDelayArg *big.Int) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.TestRequestRandomness(&_BeaconVRFConsumer.TransactOpts, numWords, subID, confirmationDelayArg)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerTransactorSession) TestRequestRandomness(numWords uint16, subID uint64, confirmationDelayArg *big.Int) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.TestRequestRandomness(&_BeaconVRFConsumer.TransactOpts, numWords, subID, confirmationDelayArg)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerTransactor) TestRequestRandomnessFulfillment(opts *bind.TransactOpts, subID uint64, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte) (*types.Transaction, error) {
	return _BeaconVRFConsumer.contract.Transact(opts, "testRequestRandomnessFulfillment", subID, numWords, confirmationDelayArg, callbackGasLimit, arguments)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerSession) TestRequestRandomnessFulfillment(subID uint64, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.TestRequestRandomnessFulfillment(&_BeaconVRFConsumer.TransactOpts, subID, numWords, confirmationDelayArg, callbackGasLimit, arguments)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerTransactorSession) TestRequestRandomnessFulfillment(subID uint64, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.TestRequestRandomnessFulfillment(&_BeaconVRFConsumer.TransactOpts, subID, numWords, confirmationDelayArg, callbackGasLimit, arguments)
}

type BeaconVRFConsumerConfigSetIterator struct {
	Event *BeaconVRFConsumerConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BeaconVRFConsumerConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BeaconVRFConsumerConfigSet)
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
		it.Event = new(BeaconVRFConsumerConfigSet)
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

func (it *BeaconVRFConsumerConfigSetIterator) Error() error {
	return it.fail
}

func (it *BeaconVRFConsumerConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BeaconVRFConsumerConfigSet struct {
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

func (_BeaconVRFConsumer *BeaconVRFConsumerFilterer) FilterConfigSet(opts *bind.FilterOpts) (*BeaconVRFConsumerConfigSetIterator, error) {

	logs, sub, err := _BeaconVRFConsumer.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &BeaconVRFConsumerConfigSetIterator{contract: _BeaconVRFConsumer.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_BeaconVRFConsumer *BeaconVRFConsumerFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *BeaconVRFConsumerConfigSet) (event.Subscription, error) {

	logs, sub, err := _BeaconVRFConsumer.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BeaconVRFConsumerConfigSet)
				if err := _BeaconVRFConsumer.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_BeaconVRFConsumer *BeaconVRFConsumerFilterer) ParseConfigSet(log types.Log) (*BeaconVRFConsumerConfigSet, error) {
	event := new(BeaconVRFConsumerConfigSet)
	if err := _BeaconVRFConsumer.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BeaconVRFConsumerNewTransmissionIterator struct {
	Event *BeaconVRFConsumerNewTransmission

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BeaconVRFConsumerNewTransmissionIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BeaconVRFConsumerNewTransmission)
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
		it.Event = new(BeaconVRFConsumerNewTransmission)
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

func (it *BeaconVRFConsumerNewTransmissionIterator) Error() error {
	return it.fail
}

func (it *BeaconVRFConsumerNewTransmissionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BeaconVRFConsumerNewTransmission struct {
	AggregatorRoundId  uint32
	EpochAndRound      *big.Int
	Transmitter        common.Address
	JuelsPerFeeCoin    *big.Int
	ReasonableGasPrice uint64
	ConfigDigest       [32]byte
	Raw                types.Log
}

func (_BeaconVRFConsumer *BeaconVRFConsumerFilterer) FilterNewTransmission(opts *bind.FilterOpts, aggregatorRoundId []uint32, epochAndRound []*big.Int) (*BeaconVRFConsumerNewTransmissionIterator, error) {

	var aggregatorRoundIdRule []interface{}
	for _, aggregatorRoundIdItem := range aggregatorRoundId {
		aggregatorRoundIdRule = append(aggregatorRoundIdRule, aggregatorRoundIdItem)
	}
	var epochAndRoundRule []interface{}
	for _, epochAndRoundItem := range epochAndRound {
		epochAndRoundRule = append(epochAndRoundRule, epochAndRoundItem)
	}

	logs, sub, err := _BeaconVRFConsumer.contract.FilterLogs(opts, "NewTransmission", aggregatorRoundIdRule, epochAndRoundRule)
	if err != nil {
		return nil, err
	}
	return &BeaconVRFConsumerNewTransmissionIterator{contract: _BeaconVRFConsumer.contract, event: "NewTransmission", logs: logs, sub: sub}, nil
}

func (_BeaconVRFConsumer *BeaconVRFConsumerFilterer) WatchNewTransmission(opts *bind.WatchOpts, sink chan<- *BeaconVRFConsumerNewTransmission, aggregatorRoundId []uint32, epochAndRound []*big.Int) (event.Subscription, error) {

	var aggregatorRoundIdRule []interface{}
	for _, aggregatorRoundIdItem := range aggregatorRoundId {
		aggregatorRoundIdRule = append(aggregatorRoundIdRule, aggregatorRoundIdItem)
	}
	var epochAndRoundRule []interface{}
	for _, epochAndRoundItem := range epochAndRound {
		epochAndRoundRule = append(epochAndRoundRule, epochAndRoundItem)
	}

	logs, sub, err := _BeaconVRFConsumer.contract.WatchLogs(opts, "NewTransmission", aggregatorRoundIdRule, epochAndRoundRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BeaconVRFConsumerNewTransmission)
				if err := _BeaconVRFConsumer.contract.UnpackLog(event, "NewTransmission", log); err != nil {
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

func (_BeaconVRFConsumer *BeaconVRFConsumerFilterer) ParseNewTransmission(log types.Log) (*BeaconVRFConsumerNewTransmission, error) {
	event := new(BeaconVRFConsumerNewTransmission)
	if err := _BeaconVRFConsumer.contract.UnpackLog(event, "NewTransmission", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BeaconVRFConsumerOutputsServedIterator struct {
	Event *BeaconVRFConsumerOutputsServed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BeaconVRFConsumerOutputsServedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BeaconVRFConsumerOutputsServed)
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
		it.Event = new(BeaconVRFConsumerOutputsServed)
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

func (it *BeaconVRFConsumerOutputsServedIterator) Error() error {
	return it.fail
}

func (it *BeaconVRFConsumerOutputsServedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BeaconVRFConsumerOutputsServed struct {
	RecentBlockHeight  uint64
	Transmitter        common.Address
	JuelsPerFeeCoin    *big.Int
	ReasonableGasPrice uint64
	OutputsServed      []VRFBeaconTypesOutputServed
	Raw                types.Log
}

func (_BeaconVRFConsumer *BeaconVRFConsumerFilterer) FilterOutputsServed(opts *bind.FilterOpts) (*BeaconVRFConsumerOutputsServedIterator, error) {

	logs, sub, err := _BeaconVRFConsumer.contract.FilterLogs(opts, "OutputsServed")
	if err != nil {
		return nil, err
	}
	return &BeaconVRFConsumerOutputsServedIterator{contract: _BeaconVRFConsumer.contract, event: "OutputsServed", logs: logs, sub: sub}, nil
}

func (_BeaconVRFConsumer *BeaconVRFConsumerFilterer) WatchOutputsServed(opts *bind.WatchOpts, sink chan<- *BeaconVRFConsumerOutputsServed) (event.Subscription, error) {

	logs, sub, err := _BeaconVRFConsumer.contract.WatchLogs(opts, "OutputsServed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BeaconVRFConsumerOutputsServed)
				if err := _BeaconVRFConsumer.contract.UnpackLog(event, "OutputsServed", log); err != nil {
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

func (_BeaconVRFConsumer *BeaconVRFConsumerFilterer) ParseOutputsServed(log types.Log) (*BeaconVRFConsumerOutputsServed, error) {
	event := new(BeaconVRFConsumerOutputsServed)
	if err := _BeaconVRFConsumer.contract.UnpackLog(event, "OutputsServed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BeaconVRFConsumerRandomWordsFulfilledIterator struct {
	Event *BeaconVRFConsumerRandomWordsFulfilled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BeaconVRFConsumerRandomWordsFulfilledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BeaconVRFConsumerRandomWordsFulfilled)
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
		it.Event = new(BeaconVRFConsumerRandomWordsFulfilled)
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

func (it *BeaconVRFConsumerRandomWordsFulfilledIterator) Error() error {
	return it.fail
}

func (it *BeaconVRFConsumerRandomWordsFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BeaconVRFConsumerRandomWordsFulfilled struct {
	RequestIDs            []*big.Int
	SuccessfulFulfillment []byte
	TruncatedErrorData    [][]byte
	Raw                   types.Log
}

func (_BeaconVRFConsumer *BeaconVRFConsumerFilterer) FilterRandomWordsFulfilled(opts *bind.FilterOpts) (*BeaconVRFConsumerRandomWordsFulfilledIterator, error) {

	logs, sub, err := _BeaconVRFConsumer.contract.FilterLogs(opts, "RandomWordsFulfilled")
	if err != nil {
		return nil, err
	}
	return &BeaconVRFConsumerRandomWordsFulfilledIterator{contract: _BeaconVRFConsumer.contract, event: "RandomWordsFulfilled", logs: logs, sub: sub}, nil
}

func (_BeaconVRFConsumer *BeaconVRFConsumerFilterer) WatchRandomWordsFulfilled(opts *bind.WatchOpts, sink chan<- *BeaconVRFConsumerRandomWordsFulfilled) (event.Subscription, error) {

	logs, sub, err := _BeaconVRFConsumer.contract.WatchLogs(opts, "RandomWordsFulfilled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BeaconVRFConsumerRandomWordsFulfilled)
				if err := _BeaconVRFConsumer.contract.UnpackLog(event, "RandomWordsFulfilled", log); err != nil {
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

func (_BeaconVRFConsumer *BeaconVRFConsumerFilterer) ParseRandomWordsFulfilled(log types.Log) (*BeaconVRFConsumerRandomWordsFulfilled, error) {
	event := new(BeaconVRFConsumerRandomWordsFulfilled)
	if err := _BeaconVRFConsumer.contract.UnpackLog(event, "RandomWordsFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BeaconVRFConsumerRandomnessFulfillmentRequestedIterator struct {
	Event *BeaconVRFConsumerRandomnessFulfillmentRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BeaconVRFConsumerRandomnessFulfillmentRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BeaconVRFConsumerRandomnessFulfillmentRequested)
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
		it.Event = new(BeaconVRFConsumerRandomnessFulfillmentRequested)
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

func (it *BeaconVRFConsumerRandomnessFulfillmentRequestedIterator) Error() error {
	return it.fail
}

func (it *BeaconVRFConsumerRandomnessFulfillmentRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BeaconVRFConsumerRandomnessFulfillmentRequested struct {
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

func (_BeaconVRFConsumer *BeaconVRFConsumerFilterer) FilterRandomnessFulfillmentRequested(opts *bind.FilterOpts, requestID []*big.Int, requester []common.Address) (*BeaconVRFConsumerRandomnessFulfillmentRequestedIterator, error) {

	var requestIDRule []interface{}
	for _, requestIDItem := range requestID {
		requestIDRule = append(requestIDRule, requestIDItem)
	}
	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _BeaconVRFConsumer.contract.FilterLogs(opts, "RandomnessFulfillmentRequested", requestIDRule, requesterRule)
	if err != nil {
		return nil, err
	}
	return &BeaconVRFConsumerRandomnessFulfillmentRequestedIterator{contract: _BeaconVRFConsumer.contract, event: "RandomnessFulfillmentRequested", logs: logs, sub: sub}, nil
}

func (_BeaconVRFConsumer *BeaconVRFConsumerFilterer) WatchRandomnessFulfillmentRequested(opts *bind.WatchOpts, sink chan<- *BeaconVRFConsumerRandomnessFulfillmentRequested, requestID []*big.Int, requester []common.Address) (event.Subscription, error) {

	var requestIDRule []interface{}
	for _, requestIDItem := range requestID {
		requestIDRule = append(requestIDRule, requestIDItem)
	}
	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _BeaconVRFConsumer.contract.WatchLogs(opts, "RandomnessFulfillmentRequested", requestIDRule, requesterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BeaconVRFConsumerRandomnessFulfillmentRequested)
				if err := _BeaconVRFConsumer.contract.UnpackLog(event, "RandomnessFulfillmentRequested", log); err != nil {
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

func (_BeaconVRFConsumer *BeaconVRFConsumerFilterer) ParseRandomnessFulfillmentRequested(log types.Log) (*BeaconVRFConsumerRandomnessFulfillmentRequested, error) {
	event := new(BeaconVRFConsumerRandomnessFulfillmentRequested)
	if err := _BeaconVRFConsumer.contract.UnpackLog(event, "RandomnessFulfillmentRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BeaconVRFConsumerRandomnessRequestedIterator struct {
	Event *BeaconVRFConsumerRandomnessRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BeaconVRFConsumerRandomnessRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BeaconVRFConsumerRandomnessRequested)
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
		it.Event = new(BeaconVRFConsumerRandomnessRequested)
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

func (it *BeaconVRFConsumerRandomnessRequestedIterator) Error() error {
	return it.fail
}

func (it *BeaconVRFConsumerRandomnessRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BeaconVRFConsumerRandomnessRequested struct {
	RequestID              *big.Int
	Requester              common.Address
	NextBeaconOutputHeight uint64
	ConfDelay              *big.Int
	SubID                  uint64
	NumWords               uint16
	Raw                    types.Log
}

func (_BeaconVRFConsumer *BeaconVRFConsumerFilterer) FilterRandomnessRequested(opts *bind.FilterOpts, requestID []*big.Int, requester []common.Address) (*BeaconVRFConsumerRandomnessRequestedIterator, error) {

	var requestIDRule []interface{}
	for _, requestIDItem := range requestID {
		requestIDRule = append(requestIDRule, requestIDItem)
	}
	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _BeaconVRFConsumer.contract.FilterLogs(opts, "RandomnessRequested", requestIDRule, requesterRule)
	if err != nil {
		return nil, err
	}
	return &BeaconVRFConsumerRandomnessRequestedIterator{contract: _BeaconVRFConsumer.contract, event: "RandomnessRequested", logs: logs, sub: sub}, nil
}

func (_BeaconVRFConsumer *BeaconVRFConsumerFilterer) WatchRandomnessRequested(opts *bind.WatchOpts, sink chan<- *BeaconVRFConsumerRandomnessRequested, requestID []*big.Int, requester []common.Address) (event.Subscription, error) {

	var requestIDRule []interface{}
	for _, requestIDItem := range requestID {
		requestIDRule = append(requestIDRule, requestIDItem)
	}
	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _BeaconVRFConsumer.contract.WatchLogs(opts, "RandomnessRequested", requestIDRule, requesterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BeaconVRFConsumerRandomnessRequested)
				if err := _BeaconVRFConsumer.contract.UnpackLog(event, "RandomnessRequested", log); err != nil {
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

func (_BeaconVRFConsumer *BeaconVRFConsumerFilterer) ParseRandomnessRequested(log types.Log) (*BeaconVRFConsumerRandomnessRequested, error) {
	event := new(BeaconVRFConsumerRandomnessRequested)
	if err := _BeaconVRFConsumer.contract.UnpackLog(event, "RandomnessRequested", log); err != nil {
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

func (_BeaconVRFConsumer *BeaconVRFConsumer) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _BeaconVRFConsumer.abi.Events["ConfigSet"].ID:
		return _BeaconVRFConsumer.ParseConfigSet(log)
	case _BeaconVRFConsumer.abi.Events["NewTransmission"].ID:
		return _BeaconVRFConsumer.ParseNewTransmission(log)
	case _BeaconVRFConsumer.abi.Events["OutputsServed"].ID:
		return _BeaconVRFConsumer.ParseOutputsServed(log)
	case _BeaconVRFConsumer.abi.Events["RandomWordsFulfilled"].ID:
		return _BeaconVRFConsumer.ParseRandomWordsFulfilled(log)
	case _BeaconVRFConsumer.abi.Events["RandomnessFulfillmentRequested"].ID:
		return _BeaconVRFConsumer.ParseRandomnessFulfillmentRequested(log)
	case _BeaconVRFConsumer.abi.Events["RandomnessRequested"].ID:
		return _BeaconVRFConsumer.ParseRandomnessRequested(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (BeaconVRFConsumerConfigSet) Topic() common.Hash {
	return common.HexToHash("0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05")
}

func (BeaconVRFConsumerNewTransmission) Topic() common.Hash {
	return common.HexToHash("0x27bf3f1077f091da6885751ba10f5775d06657fd59e47a6ab1f7635e5a115afe")
}

func (BeaconVRFConsumerOutputsServed) Topic() common.Hash {
	return common.HexToHash("0xe1d18855b43b829b66a7f664301b0733507d67e5d8e163e3f0b778717e884ee0")
}

func (BeaconVRFConsumerRandomWordsFulfilled) Topic() common.Hash {
	return common.HexToHash("0x47ddf7bb0cbd94c1b43c5097f1352a80db0ceb3696f029d32b24f32cd631d2b7")
}

func (BeaconVRFConsumerRandomnessFulfillmentRequested) Topic() common.Hash {
	return common.HexToHash("0xda7e52094a2e0200a38ad6b72d1ce40dcc840bbc4cf95ccf8065ef521da664dc")
}

func (BeaconVRFConsumerRandomnessRequested) Topic() common.Hash {
	return common.HexToHash("0xdcbc2f1d67e18a0e5cf023ad14515e21bfc7e0cccb7ddb1011f46baffc7d9d15")
}

func (_BeaconVRFConsumer *BeaconVRFConsumer) Address() common.Address {
	return _BeaconVRFConsumer.address
}

type BeaconVRFConsumerInterface interface {
	NUMCONFDELAYS(opts *bind.CallOpts) (uint8, error)

	Fail(opts *bind.CallOpts) (bool, error)

	IBeaconPeriodBlocks(opts *bind.CallOpts) (*big.Int, error)

	SReceivedRandomnessByRequestID(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error)

	SArguments(opts *bind.CallOpts) ([]byte, error)

	SGasAvailable(opts *bind.CallOpts) (*big.Int, error)

	SMyBeaconRequests(opts *bind.CallOpts, arg0 *big.Int) (SMyBeaconRequests,

		error)

	SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	SRequestsIDs(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error)

	SSubId(opts *bind.CallOpts) (uint64, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestID *big.Int, randomWords []*big.Int, arguments []byte) (*types.Transaction, error)

	SetFail(opts *bind.TransactOpts, shouldFail bool) (*types.Transaction, error)

	StoreBeaconRequest(opts *bind.TransactOpts, reqId *big.Int, height *big.Int, delay *big.Int, numWords uint16) (*types.Transaction, error)

	TestRedeemRandomness(opts *bind.TransactOpts, requestID *big.Int) (*types.Transaction, error)

	TestRequestRandomness(opts *bind.TransactOpts, numWords uint16, subID uint64, confirmationDelayArg *big.Int) (*types.Transaction, error)

	TestRequestRandomnessFulfillment(opts *bind.TransactOpts, subID uint64, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte) (*types.Transaction, error)

	FilterConfigSet(opts *bind.FilterOpts) (*BeaconVRFConsumerConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *BeaconVRFConsumerConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*BeaconVRFConsumerConfigSet, error)

	FilterNewTransmission(opts *bind.FilterOpts, aggregatorRoundId []uint32, epochAndRound []*big.Int) (*BeaconVRFConsumerNewTransmissionIterator, error)

	WatchNewTransmission(opts *bind.WatchOpts, sink chan<- *BeaconVRFConsumerNewTransmission, aggregatorRoundId []uint32, epochAndRound []*big.Int) (event.Subscription, error)

	ParseNewTransmission(log types.Log) (*BeaconVRFConsumerNewTransmission, error)

	FilterOutputsServed(opts *bind.FilterOpts) (*BeaconVRFConsumerOutputsServedIterator, error)

	WatchOutputsServed(opts *bind.WatchOpts, sink chan<- *BeaconVRFConsumerOutputsServed) (event.Subscription, error)

	ParseOutputsServed(log types.Log) (*BeaconVRFConsumerOutputsServed, error)

	FilterRandomWordsFulfilled(opts *bind.FilterOpts) (*BeaconVRFConsumerRandomWordsFulfilledIterator, error)

	WatchRandomWordsFulfilled(opts *bind.WatchOpts, sink chan<- *BeaconVRFConsumerRandomWordsFulfilled) (event.Subscription, error)

	ParseRandomWordsFulfilled(log types.Log) (*BeaconVRFConsumerRandomWordsFulfilled, error)

	FilterRandomnessFulfillmentRequested(opts *bind.FilterOpts, requestID []*big.Int, requester []common.Address) (*BeaconVRFConsumerRandomnessFulfillmentRequestedIterator, error)

	WatchRandomnessFulfillmentRequested(opts *bind.WatchOpts, sink chan<- *BeaconVRFConsumerRandomnessFulfillmentRequested, requestID []*big.Int, requester []common.Address) (event.Subscription, error)

	ParseRandomnessFulfillmentRequested(log types.Log) (*BeaconVRFConsumerRandomnessFulfillmentRequested, error)

	FilterRandomnessRequested(opts *bind.FilterOpts, requestID []*big.Int, requester []common.Address) (*BeaconVRFConsumerRandomnessRequestedIterator, error)

	WatchRandomnessRequested(opts *bind.WatchOpts, sink chan<- *BeaconVRFConsumerRandomnessRequested, requestID []*big.Int, requester []common.Address) (event.Subscription, error)

	ParseRandomnessRequested(log types.Log) (*BeaconVRFConsumerRandomnessRequested, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
