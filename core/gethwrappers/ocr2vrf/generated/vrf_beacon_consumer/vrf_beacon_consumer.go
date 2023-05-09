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

var BeaconVRFConsumerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"shouldFail\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"beaconPeriodBlocks\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"MustBeRouter\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"aggregatorRoundId\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"uint40\",\"name\":\"epochAndRound\",\"type\":\"uint40\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint192\",\"name\":\"juelsPerFeeCoin\",\"type\":\"uint192\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"reasonableGasPrice\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"NewTransmission\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"recentBlockHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint192\",\"name\":\"juelsPerFeeCoin\",\"type\":\"uint192\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"reasonableGasPrice\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelay\",\"type\":\"uint24\"},{\"internalType\":\"uint256\",\"name\":\"proofG1X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"proofG1Y\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structVRFBeaconTypes.OutputServed[]\",\"name\":\"outputsServed\",\"type\":\"tuple[]\"}],\"name\":\"OutputsServed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"VRFBeaconTypes.RequestID[]\",\"name\":\"requestIDs\",\"type\":\"uint48[]\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"successfulFulfillment\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes[]\",\"name\":\"truncatedErrorData\",\"type\":\"bytes[]\"}],\"name\":\"RandomWordsFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"nextBeaconOutputHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confDelay\",\"type\":\"uint24\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasAllowance\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"weiPerUnitLink\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"}],\"name\":\"RandomnessFulfillmentRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"nextBeaconOutputHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confDelay\",\"type\":\"uint24\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"}],\"name\":\"RandomnessRequested\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"NUM_CONF_DELAYS\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"fail\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_beaconPeriodBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_ReceivedRandomnessByRequestID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_arguments\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_gasAvailable\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"}],\"name\":\"s_myBeaconRequests\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.SlotNumber\",\"name\":\"slotNumber\",\"type\":\"uint32\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelay\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_randomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"\",\"type\":\"uint24\"}],\"name\":\"s_requestsIDs\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"}],\"name\":\"s_reverseRequestIDMapping\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_subId\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"shouldFail\",\"type\":\"bool\"}],\"name\":\"setFail\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"reqId\",\"type\":\"uint48\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"delay\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"}],\"name\":\"storeBeaconRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"}],\"name\":\"testRedeemRandomness\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelayArg\",\"type\":\"uint24\"}],\"name\":\"testRequestRandomness\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelayArg\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"}],\"name\":\"testRequestRandomnessFulfillment\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561001057600080fd5b5060405161139a38038061139a83398101604081905261002f9161006c565b6001600160a01b03929092166080819052600680546001600160a01b03191690911790556009805460ff1916911515919091179055600a556100c0565b60008060006060848603121561008157600080fd5b83516001600160a01b038116811461009857600080fd5b602085015190935080151581146100ae57600080fd5b80925050604084015190509250925092565b6080516112bf6100db600039600061038301526112bf6000f3fe608060405234801561001057600080fd5b506004361061011b5760003560e01c80639d769402116100b2578063d4fab62c11610081578063f08c5daa11610066578063f08c5daa14610352578063f339c7941461035b578063f6eaffc81461036e57600080fd5b8063d4fab62c1461032c578063ea7502ab1461033f57600080fd5b80639d769402146102d2578063a9cc4718146102f3578063c6d6130114610310578063cd0593df1461032357600080fd5b80635f15cccc116100ee5780635f15cccc1461021b5780636d162a3e14610265578063706da1ca146102785780637716cdaa146102bd57600080fd5b806319a5fa22146101205780632f7527cc146101be5780634e2981b1146101d85780635a47dd7114610206575b600080fd5b61017c61012e366004610b34565b60026020526000908152604090205463ffffffff811690640100000000810462ffffff1690670100000000000000810461ffff1690690100000000000000000090046001600160a01b031684565b6040805163ffffffff909516855262ffffff909316602085015261ffff909116918301919091526001600160a01b031660608201526080015b60405180910390f35b6101c6600881565b60405160ff90911681526020016101b5565b6101f86101e6366004610b34565b60046020526000908152604090205481565b6040519081526020016101b5565b610219610214366004610c33565b610381565b005b61024e610229366004610d1d565b600160209081526000928352604080842090915290825290205465ffffffffffff1681565b60405165ffffffffffff90911681526020016101b5565b610219610273366004610d5b565b6103f3565b6007546102a49074010000000000000000000000000000000000000000900467ffffffffffffffff1681565b60405167ffffffffffffffff90911681526020016101b5565b6102c5610526565b6040516101b59190610df7565b6102196102e0366004610e0a565b6009805460ff1916911515919091179055565b6009546103009060ff1681565b60405190151581526020016101b5565b61024e61031e366004610e2c565b6105b4565b6101f8600a5481565b61021961033a366004610e68565b610753565b61024e61034d366004610e98565b61083b565b6101f860085481565b6101f8610369366004610f1c565b61093b565b6101f861037c366004610f48565b61096c565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031633146103e3576040517ff74c318f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6103ee83838361098d565b505050565b600083815260016020908152604080832062ffffff861684529091528120805465ffffffffffff191665ffffffffffff8716179055600a546104359085610f8d565b6040805160808101825263ffffffff928316815262ffffff958616602080830191825261ffff968716838501908152306060850190815265ffffffffffff909b1660009081526002909252939020915182549151935199516001600160a01b03166901000000000000000000027fffffff0000000000000000000000000000000000000000ffffffffffffffffff9a90971667010000000000000002999099167fffffff00000000000000000000000000000000000000000000ffffffffffffff939097166401000000000266ffffffffffffff199091169890931697909717919091171692909217179092555050565b6005805461053390610fa1565b80601f016020809104026020016040519081016040528092919081815260200182805461055f90610fa1565b80156105ac5780601f10610581576101008083540402835291602001916105ac565b820191906000526020600020905b81548152906001019060200180831161058f57829003601f168201915b505050505081565b600080600a546105c2610a33565b6105cc9190610fdb565b9050600081600a546105dc610a33565b6105e69190610fef565b6105f09190611007565b60065460408051602081018252600080825291517f4ffac83a00000000000000000000000000000000000000000000000000000000815293945090926001600160a01b0390921691634ffac83a91610651918a918c918b919060040161101e565b6020604051808303816000875af1158015610670573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106949190611056565b6006546040517f7612e884000000000000000000000000000000000000000000000000000000008152600481018390529192506000916001600160a01b0390911690637612e884906024016040805180830381865afa1580156106fb573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061071f919061106f565b5065ffffffffffff8116600090815260046020526040902083905590506107488184888b6103f3565b979650505050505050565b60065465ffffffffffff82166000908152600460208181526040808420548151928301825284835290517facfc6cdd00000000000000000000000000000000000000000000000000000000815293946001600160a01b03169363acfc6cdd936107c1938993929091016110aa565b6000604051808303816000875af11580156107e0573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f1916820160405261080891908101906110d2565b65ffffffffffff831660009081526003602090815260409091208251929350610835929091840190610abd565b50505050565b600080600a54610849610a33565b6108539190610fdb565b9050600081600a54610863610a33565b61086d9190610fef565b6108779190611007565b60065460408051602081018252600080825291517fdb972c8b00000000000000000000000000000000000000000000000000000000815293945090926001600160a01b039092169163db972c8b916108dc918d918d918d918d918d9190600401611158565b6020604051808303816000875af11580156108fb573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061091f9190611056565b90508061092e81848a8c6103f3565b9998505050505050505050565b6003602052816000526040600020818154811061095757600080fd5b90600052602060002001600091509150505481565b6000818154811061097c57600080fd5b600091825260209091200154905081565b60095460ff16156109fe576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f206661696c656420696e2066756c66696c6c52616e646f6d576f726473000000604482015260640160405180910390fd5b65ffffffffffff831660009081526003602090815260409091208351610a2692850190610abd565b50600561083582826111f2565b60004661a4b1811480610a48575062066eed81145b15610ab65760646001600160a01b031663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610a8c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ab09190611056565b91505090565b4391505090565b828054828255906000526020600020908101928215610af8579160200282015b82811115610af8578251825591602001919060010190610add565b50610b04929150610b08565b5090565b5b80821115610b045760008155600101610b09565b65ffffffffffff81168114610b3157600080fd5b50565b600060208284031215610b4657600080fd5b8135610b5181610b1d565b9392505050565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff81118282101715610b9757610b97610b58565b604052919050565b600067ffffffffffffffff821115610bb957610bb9610b58565b5060051b60200190565b600082601f830112610bd457600080fd5b813567ffffffffffffffff811115610bee57610bee610b58565b610c01601f8201601f1916602001610b6e565b818152846020838601011115610c1657600080fd5b816020850160208301376000918101602001919091529392505050565b600080600060608486031215610c4857600080fd5b8335610c5381610b1d565b925060208481013567ffffffffffffffff80821115610c7157600080fd5b818701915087601f830112610c8557600080fd5b8135610c98610c9382610b9f565b610b6e565b81815260059190911b8301840190848101908a831115610cb757600080fd5b938501935b82851015610cd557843582529385019390850190610cbc565b965050506040870135925080831115610ced57600080fd5b5050610cfb86828701610bc3565b9150509250925092565b803562ffffff81168114610d1857600080fd5b919050565b60008060408385031215610d3057600080fd5b82359150610d4060208401610d05565b90509250929050565b803561ffff81168114610d1857600080fd5b60008060008060808587031215610d7157600080fd5b8435610d7c81610b1d565b935060208501359250610d9160408601610d05565b9150610d9f60608601610d49565b905092959194509250565b6000815180845260005b81811015610dd057602081850181015186830182015201610db4565b81811115610de2576000602083870101525b50601f01601f19169290920160200192915050565b602081526000610b516020830184610daa565b600060208284031215610e1c57600080fd5b81358015158114610b5157600080fd5b600080600060608486031215610e4157600080fd5b610e4a84610d49565b925060208401359150610e5f60408501610d05565b90509250925092565b60008060408385031215610e7b57600080fd5b823591506020830135610e8d81610b1d565b809150509250929050565b600080600080600060a08688031215610eb057600080fd5b85359450610ec060208701610d49565b9350610ece60408701610d05565b9250606086013563ffffffff81168114610ee757600080fd5b9150608086013567ffffffffffffffff811115610f0357600080fd5b610f0f88828901610bc3565b9150509295509295909350565b60008060408385031215610f2f57600080fd5b8235610f3a81610b1d565b946020939093013593505050565b600060208284031215610f5a57600080fd5b5035919050565b634e487b7160e01b600052601260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b600082610f9c57610f9c610f61565b500490565b600181811c90821680610fb557607f821691505b602082108103610fd557634e487b7160e01b600052602260045260246000fd5b50919050565b600082610fea57610fea610f61565b500690565b6000821982111561100257611002610f77565b500190565b60008282101561101957611019610f77565b500390565b84815261ffff8416602082015262ffffff8316604082015260806060820152600061104c6080830184610daa565b9695505050505050565b60006020828403121561106857600080fd5b5051919050565b6000806040838503121561108257600080fd5b825161108d81610b1d565b60208401519092506001600160a01b0381168114610e8d57600080fd5b8381528260208201526060604082015260006110c96060830184610daa565b95945050505050565b600060208083850312156110e557600080fd5b825167ffffffffffffffff8111156110fc57600080fd5b8301601f8101851361110d57600080fd5b805161111b610c9382610b9f565b81815260059190911b8201830190838101908783111561113a57600080fd5b928401925b828410156107485783518252928401929084019061113f565b86815261ffff8616602082015262ffffff8516604082015263ffffffff8416606082015260c06080820152600061119260c0830185610daa565b82810360a084015261092e8185610daa565b601f8211156103ee57600081815260208120601f850160051c810160208610156111cb5750805b601f850160051c820191505b818110156111ea578281556001016111d7565b505050505050565b815167ffffffffffffffff81111561120c5761120c610b58565b6112208161121a8454610fa1565b846111a4565b602080601f831160018114611255576000841561123d5750858301515b600019600386901b1c1916600185901b1785556111ea565b600085815260208120601f198616915b8281101561128457888601518255948401946001909101908401611265565b50858210156112a25787850151600019600388901b60f8161c191681555b5050505050600190811b0190555056fea164736f6c634300080f000a",
}

var BeaconVRFConsumerABI = BeaconVRFConsumerMetaData.ABI

var BeaconVRFConsumerBin = BeaconVRFConsumerMetaData.Bin

func DeployBeaconVRFConsumer(auth *bind.TransactOpts, backend bind.ContractBackend, router common.Address, shouldFail bool, beaconPeriodBlocks *big.Int) (common.Address, *types.Transaction, *BeaconVRFConsumer, error) {
	parsed, err := BeaconVRFConsumerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BeaconVRFConsumerBin), backend, router, shouldFail, beaconPeriodBlocks)
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
	parsed, err := BeaconVRFConsumerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
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

func (_BeaconVRFConsumer *BeaconVRFConsumerCaller) SReverseRequestIDMapping(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _BeaconVRFConsumer.contract.Call(opts, &out, "s_reverseRequestIDMapping", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_BeaconVRFConsumer *BeaconVRFConsumerSession) SReverseRequestIDMapping(arg0 *big.Int) (*big.Int, error) {
	return _BeaconVRFConsumer.Contract.SReverseRequestIDMapping(&_BeaconVRFConsumer.CallOpts, arg0)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCallerSession) SReverseRequestIDMapping(arg0 *big.Int) (*big.Int, error) {
	return _BeaconVRFConsumer.Contract.SReverseRequestIDMapping(&_BeaconVRFConsumer.CallOpts, arg0)
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

func (_BeaconVRFConsumer *BeaconVRFConsumerTransactor) TestRedeemRandomness(opts *bind.TransactOpts, subID *big.Int, requestID *big.Int) (*types.Transaction, error) {
	return _BeaconVRFConsumer.contract.Transact(opts, "testRedeemRandomness", subID, requestID)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerSession) TestRedeemRandomness(subID *big.Int, requestID *big.Int) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.TestRedeemRandomness(&_BeaconVRFConsumer.TransactOpts, subID, requestID)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerTransactorSession) TestRedeemRandomness(subID *big.Int, requestID *big.Int) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.TestRedeemRandomness(&_BeaconVRFConsumer.TransactOpts, subID, requestID)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerTransactor) TestRequestRandomness(opts *bind.TransactOpts, numWords uint16, subID *big.Int, confirmationDelayArg *big.Int) (*types.Transaction, error) {
	return _BeaconVRFConsumer.contract.Transact(opts, "testRequestRandomness", numWords, subID, confirmationDelayArg)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerSession) TestRequestRandomness(numWords uint16, subID *big.Int, confirmationDelayArg *big.Int) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.TestRequestRandomness(&_BeaconVRFConsumer.TransactOpts, numWords, subID, confirmationDelayArg)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerTransactorSession) TestRequestRandomness(numWords uint16, subID *big.Int, confirmationDelayArg *big.Int) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.TestRequestRandomness(&_BeaconVRFConsumer.TransactOpts, numWords, subID, confirmationDelayArg)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerTransactor) TestRequestRandomnessFulfillment(opts *bind.TransactOpts, subID *big.Int, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte) (*types.Transaction, error) {
	return _BeaconVRFConsumer.contract.Transact(opts, "testRequestRandomnessFulfillment", subID, numWords, confirmationDelayArg, callbackGasLimit, arguments)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerSession) TestRequestRandomnessFulfillment(subID *big.Int, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.TestRequestRandomnessFulfillment(&_BeaconVRFConsumer.TransactOpts, subID, numWords, confirmationDelayArg, callbackGasLimit, arguments)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerTransactorSession) TestRequestRandomnessFulfillment(subID *big.Int, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte) (*types.Transaction, error) {
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
	SubID                  *big.Int
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
	SubID                  *big.Int
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
	return common.HexToHash("0xf10ea936d00579b4c52035ee33bf46929646b3aa87554c565d8fb2c7aa549c44")
}

func (BeaconVRFConsumerRandomWordsFulfilled) Topic() common.Hash {
	return common.HexToHash("0x47ddf7bb0cbd94c1b43c5097f1352a80db0ceb3696f029d32b24f32cd631d2b7")
}

func (BeaconVRFConsumerRandomnessFulfillmentRequested) Topic() common.Hash {
	return common.HexToHash("0x24f0e469e0097d1e8d9975137f9f4dd17d2c1481b3a2f25f2382f51287eda1dc")
}

func (BeaconVRFConsumerRandomnessRequested) Topic() common.Hash {
	return common.HexToHash("0xc3b31df4232b05afd212fc28027dae6fd6a81618c2a3116182cb57c7f0a3fd0a")
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

	SReverseRequestIDMapping(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	SSubId(opts *bind.CallOpts) (uint64, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestID *big.Int, randomWords []*big.Int, arguments []byte) (*types.Transaction, error)

	SetFail(opts *bind.TransactOpts, shouldFail bool) (*types.Transaction, error)

	StoreBeaconRequest(opts *bind.TransactOpts, reqId *big.Int, height *big.Int, delay *big.Int, numWords uint16) (*types.Transaction, error)

	TestRedeemRandomness(opts *bind.TransactOpts, subID *big.Int, requestID *big.Int) (*types.Transaction, error)

	TestRequestRandomness(opts *bind.TransactOpts, numWords uint16, subID *big.Int, confirmationDelayArg *big.Int) (*types.Transaction, error)

	TestRequestRandomnessFulfillment(opts *bind.TransactOpts, subID *big.Int, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte) (*types.Transaction, error)

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
