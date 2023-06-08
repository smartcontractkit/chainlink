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

var BeaconVRFConsumerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"shouldFail\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"beaconPeriodBlocks\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"MustBeCoordinator\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeOwnerOrCoordinator\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"}],\"name\":\"CoordinatorUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"fail\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_beaconPeriodBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestID\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_ReceivedRandomnessByRequestID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_arguments\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_gasAvailable\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_mostRecentRequestID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_myBeaconRequests\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.SlotNumber\",\"name\":\"slotNumber\",\"type\":\"uint32\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelay\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_randomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"\",\"type\":\"uint24\"}],\"name\":\"s_requestsIDs\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_subId\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"}],\"name\":\"setCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"shouldFail\",\"type\":\"bool\"}],\"name\":\"setFail\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"reqId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"delay\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"}],\"name\":\"storeBeaconRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requestID\",\"type\":\"uint256\"}],\"name\":\"testRedeemRandomness\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelayArg\",\"type\":\"uint24\"}],\"name\":\"testRequestRandomness\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelayArg\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"}],\"name\":\"testRequestRandomnessFulfillment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162001602380380620016028339810160408190526200003491620001aa565b8233806000816200008c5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000bf57620000bf81620000ff565b5050600280546001600160a01b0319166001600160a01b03939093169290921790915550600b805460ff191692151592909217909155600c555062000201565b336001600160a01b03821603620001595760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000083565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600080600060608486031215620001c057600080fd5b83516001600160a01b0381168114620001d857600080fd5b60208501519093508015158114620001ef57600080fd5b80925050604084015190509250925092565b6113f180620002116000396000f3fe608060405234801561001057600080fd5b506004361061016c5760003560e01c8063a9cc4718116100cd578063ea7502ab11610081578063f2fde38b11610066578063f2fde38b146102f4578063f6eaffc814610307578063ffe97ca41461031a57600080fd5b8063ea7502ab146102d8578063f08c5daa146102eb57600080fd5b8063cd0593df116100b2578063cd0593df146102a9578063d0705f04146102b2578063d21ea8fd146102c557600080fd5b8063a9cc471814610279578063c6d613011461029657600080fd5b80637716cdaa116101245780638da5cb5b116101095780638da5cb5b1461022a5780638ea98117146102455780639d7694021461025857600080fd5b80637716cdaa1461020d57806379ba50971461022257600080fd5b8063689b77ab11610155578063689b77ab146101c45780636df57cc3146101cd578063706da1ca146101e057600080fd5b8063341867a2146101715780635f15cccc14610186575b600080fd5b61018461017f366004610cd3565b6103b3565b005b6101b1610194366004610d0d565b600460209081526000928352604080842090915290825290205481565b6040519081526020015b60405180910390f35b6101b160085481565b6101846101db366004610d4b565b61047d565b6009546101f49067ffffffffffffffff1681565b60405167ffffffffffffffff90911681526020016101bb565b610215610593565b6040516101bb9190610dd7565b610184610621565b6000546040516001600160a01b0390911681526020016101bb565b610184610253366004610df1565b6106e4565b610184610266366004610e1a565b600b805460ff1916911515919091179055565b600b546102869060ff1681565b60405190151581526020016101bb565b6101b16102a4366004610e3c565b610798565b6101b1600c5481565b6101b16102c0366004610cd3565b610895565b6101846102d3366004610f53565b6108c6565b6101b16102e636600461101c565b61091a565b6101b1600a5481565b610184610302366004610df1565b610a23565b6101b16103153660046110a0565b610a37565b6103766103283660046110a0565b60056020526000908152604090205463ffffffff811690640100000000810462ffffff1690670100000000000000810461ffff1690690100000000000000000090046001600160a01b031684565b6040805163ffffffff909516855262ffffff909316602085015261ffff909116918301919091526001600160a01b031660608201526080016101bb565b60025460408051602081018252600080825291517facfc6cdd00000000000000000000000000000000000000000000000000000000815291926001600160a01b03169163acfc6cdd9161040c91879187916004016110b9565b6000604051808303816000875af115801561042b573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f1916820160405261045391908101906110e1565b60008381526006602090815260409091208251929350610477929091840190610c73565b50505050565b600083815260046020908152604080832062ffffff861684529091528120859055600c546104ab908561119e565b6040805160808101825263ffffffff928316815262ffffff958616602080830191825261ffff968716838501908152306060850190815260009b8c52600590925293909920915182549151935199516001600160a01b03166901000000000000000000027fffffff0000000000000000000000000000000000000000ffffffffffffffffff9a90971667010000000000000002999099167fffffff00000000000000000000000000000000000000000000ffffffffffffff939097166401000000000266ffffffffffffff199091169890931697909717919091171692909217179092555050565b600780546105a0906111b2565b80601f01602080910402602001604051908101604052809291908181526020018280546105cc906111b2565b80156106195780601f106105ee57610100808354040283529160200191610619565b820191906000526020600020905b8154815290600101906020018083116105fc57829003601f168201915b505050505081565b6001546001600160a01b031633146106805760405162461bcd60e51b815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b600080543373ffffffffffffffffffffffffffffffffffffffff19808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6000546001600160a01b0316331480159061070a57506002546001600160a01b03163314155b15610741576040517fd4e06fd700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6002805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0383169081179091556040517fc258faa9a17ddfdf4130b4acff63a289202e7d5f9e42f366add65368575486bc90600090a250565b600080600c546107a6610a58565b6107b091906111ec565b9050600081600c546107c0610a58565b6107ca9190611200565b6107d49190611219565b60025460408051602081018252600080825291517f4ffac83a00000000000000000000000000000000000000000000000000000000815293945090926001600160a01b0390921691634ffac83a91610835918a918c918b919060040161122c565b6020604051808303816000875af1158015610854573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108789190611264565b90506108868183878a61047d565b60088190559695505050505050565b600660205281600052604060002081815481106108b157600080fd5b90600052602060002001600091509150505481565b6002546001600160a01b0316331461090a576040517f66bf9c7200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610915838383610ae2565b505050565b6000808490506000600c5461092d610a58565b61093791906111ec565b9050600081600c54610947610a58565b6109519190611200565b61095b9190611219565b60025460408051602081018252600080825291517fdb972c8b00000000000000000000000000000000000000000000000000000000815293945090926001600160a01b039092169163db972c8b916109c0918e918e918a918e918e919060040161127d565b6020604051808303816000875af11580156109df573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a039190611264565b9050610a1181838a8c61047d565b60088190559998505050505050505050565b610a2b610b61565b610a3481610bbd565b50565b60038181548110610a4757600080fd5b600091825260209091200154905081565b60004661a4b1811480610a6d575062066eed81145b15610adb5760646001600160a01b031663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610ab1573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ad59190611264565b91505090565b4391505090565b600b5460ff1615610b355760405162461bcd60e51b815260206004820152601d60248201527f206661696c656420696e2066756c66696c6c52616e646f6d576f7264730000006044820152606401610677565b60008381526006602090815260409091208351610b5492850190610c73565b5060076104778282611324565b6000546001600160a01b03163314610bbb5760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610677565b565b336001600160a01b03821603610c155760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610677565b6001805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b828054828255906000526020600020908101928215610cae579160200282015b82811115610cae578251825591602001919060010190610c93565b50610cba929150610cbe565b5090565b5b80821115610cba5760008155600101610cbf565b60008060408385031215610ce657600080fd5b50508035926020909101359150565b803562ffffff81168114610d0857600080fd5b919050565b60008060408385031215610d2057600080fd5b82359150610d3060208401610cf5565b90509250929050565b803561ffff81168114610d0857600080fd5b60008060008060808587031215610d6157600080fd5b8435935060208501359250610d7860408601610cf5565b9150610d8660608601610d39565b905092959194509250565b6000815180845260005b81811015610db757602081850181015186830182015201610d9b565b506000602082860101526020601f19601f83011685010191505092915050565b602081526000610dea6020830184610d91565b9392505050565b600060208284031215610e0357600080fd5b81356001600160a01b0381168114610dea57600080fd5b600060208284031215610e2c57600080fd5b81358015158114610dea57600080fd5b600080600060608486031215610e5157600080fd5b610e5a84610d39565b925060208401359150610e6f60408501610cf5565b90509250925092565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff81118282101715610eb757610eb7610e78565b604052919050565b600067ffffffffffffffff821115610ed957610ed9610e78565b5060051b60200190565b600082601f830112610ef457600080fd5b813567ffffffffffffffff811115610f0e57610f0e610e78565b610f21601f8201601f1916602001610e8e565b818152846020838601011115610f3657600080fd5b816020850160208301376000918101602001919091529392505050565b600080600060608486031215610f6857600080fd5b8335925060208085013567ffffffffffffffff80821115610f8857600080fd5b818701915087601f830112610f9c57600080fd5b8135610faf610faa82610ebf565b610e8e565b81815260059190911b8301840190848101908a831115610fce57600080fd5b938501935b82851015610fec57843582529385019390850190610fd3565b96505050604087013592508083111561100457600080fd5b505061101286828701610ee3565b9150509250925092565b600080600080600060a0868803121561103457600080fd5b8535945061104460208701610d39565b935061105260408701610cf5565b9250606086013563ffffffff8116811461106b57600080fd5b9150608086013567ffffffffffffffff81111561108757600080fd5b61109388828901610ee3565b9150509295509295909350565b6000602082840312156110b257600080fd5b5035919050565b8381528260208201526060604082015260006110d86060830184610d91565b95945050505050565b600060208083850312156110f457600080fd5b825167ffffffffffffffff81111561110b57600080fd5b8301601f8101851361111c57600080fd5b805161112a610faa82610ebf565b81815260059190911b8201830190838101908783111561114957600080fd5b928401925b828410156111675783518252928401929084019061114e565b979650505050505050565b634e487b7160e01b600052601260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b6000826111ad576111ad611172565b500490565b600181811c908216806111c657607f821691505b6020821081036111e657634e487b7160e01b600052602260045260246000fd5b50919050565b6000826111fb576111fb611172565b500690565b8082018082111561121357611213611188565b92915050565b8181038181111561121357611213611188565b84815261ffff8416602082015262ffffff8316604082015260806060820152600061125a6080830184610d91565b9695505050505050565b60006020828403121561127657600080fd5b5051919050565b86815261ffff8616602082015262ffffff8516604082015263ffffffff8416606082015260c0608082015260006112b760c0830185610d91565b82810360a08401526112c98185610d91565b9998505050505050505050565b601f82111561091557600081815260208120601f850160051c810160208610156112fd5750805b601f850160051c820191505b8181101561131c57828155600101611309565b505050505050565b815167ffffffffffffffff81111561133e5761133e610e78565b6113528161134c84546111b2565b846112d6565b602080601f831160018114611387576000841561136f5750858301515b600019600386901b1c1916600185901b17855561131c565b600085815260208120601f198616915b828110156113b657888601518255948401946001909101908401611397565b50858210156113d45787850151600019600388901b60f8161c191681555b5050505050600190811b0190555056fea164736f6c6343000813000a",
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

func (_BeaconVRFConsumer *BeaconVRFConsumerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BeaconVRFConsumer.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BeaconVRFConsumer *BeaconVRFConsumerSession) Owner() (common.Address, error) {
	return _BeaconVRFConsumer.Contract.Owner(&_BeaconVRFConsumer.CallOpts)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCallerSession) Owner() (common.Address, error) {
	return _BeaconVRFConsumer.Contract.Owner(&_BeaconVRFConsumer.CallOpts)
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

func (_BeaconVRFConsumer *BeaconVRFConsumerCaller) SMostRecentRequestID(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BeaconVRFConsumer.contract.Call(opts, &out, "s_mostRecentRequestID")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_BeaconVRFConsumer *BeaconVRFConsumerSession) SMostRecentRequestID() (*big.Int, error) {
	return _BeaconVRFConsumer.Contract.SMostRecentRequestID(&_BeaconVRFConsumer.CallOpts)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCallerSession) SMostRecentRequestID() (*big.Int, error) {
	return _BeaconVRFConsumer.Contract.SMostRecentRequestID(&_BeaconVRFConsumer.CallOpts)
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

func (_BeaconVRFConsumer *BeaconVRFConsumerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BeaconVRFConsumer.contract.Transact(opts, "acceptOwnership")
}

func (_BeaconVRFConsumer *BeaconVRFConsumerSession) AcceptOwnership() (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.AcceptOwnership(&_BeaconVRFConsumer.TransactOpts)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.AcceptOwnership(&_BeaconVRFConsumer.TransactOpts)
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

func (_BeaconVRFConsumer *BeaconVRFConsumerTransactor) SetCoordinator(opts *bind.TransactOpts, coordinator common.Address) (*types.Transaction, error) {
	return _BeaconVRFConsumer.contract.Transact(opts, "setCoordinator", coordinator)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerSession) SetCoordinator(coordinator common.Address) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.SetCoordinator(&_BeaconVRFConsumer.TransactOpts, coordinator)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerTransactorSession) SetCoordinator(coordinator common.Address) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.SetCoordinator(&_BeaconVRFConsumer.TransactOpts, coordinator)
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

func (_BeaconVRFConsumer *BeaconVRFConsumerTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _BeaconVRFConsumer.contract.Transact(opts, "transferOwnership", to)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.TransferOwnership(&_BeaconVRFConsumer.TransactOpts, to)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.TransferOwnership(&_BeaconVRFConsumer.TransactOpts, to)
}

type BeaconVRFConsumerCoordinatorUpdatedIterator struct {
	Event *BeaconVRFConsumerCoordinatorUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BeaconVRFConsumerCoordinatorUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BeaconVRFConsumerCoordinatorUpdated)
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
		it.Event = new(BeaconVRFConsumerCoordinatorUpdated)
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

func (it *BeaconVRFConsumerCoordinatorUpdatedIterator) Error() error {
	return it.fail
}

func (it *BeaconVRFConsumerCoordinatorUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BeaconVRFConsumerCoordinatorUpdated struct {
	Coordinator common.Address
	Raw         types.Log
}

func (_BeaconVRFConsumer *BeaconVRFConsumerFilterer) FilterCoordinatorUpdated(opts *bind.FilterOpts, coordinator []common.Address) (*BeaconVRFConsumerCoordinatorUpdatedIterator, error) {

	var coordinatorRule []interface{}
	for _, coordinatorItem := range coordinator {
		coordinatorRule = append(coordinatorRule, coordinatorItem)
	}

	logs, sub, err := _BeaconVRFConsumer.contract.FilterLogs(opts, "CoordinatorUpdated", coordinatorRule)
	if err != nil {
		return nil, err
	}
	return &BeaconVRFConsumerCoordinatorUpdatedIterator{contract: _BeaconVRFConsumer.contract, event: "CoordinatorUpdated", logs: logs, sub: sub}, nil
}

func (_BeaconVRFConsumer *BeaconVRFConsumerFilterer) WatchCoordinatorUpdated(opts *bind.WatchOpts, sink chan<- *BeaconVRFConsumerCoordinatorUpdated, coordinator []common.Address) (event.Subscription, error) {

	var coordinatorRule []interface{}
	for _, coordinatorItem := range coordinator {
		coordinatorRule = append(coordinatorRule, coordinatorItem)
	}

	logs, sub, err := _BeaconVRFConsumer.contract.WatchLogs(opts, "CoordinatorUpdated", coordinatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BeaconVRFConsumerCoordinatorUpdated)
				if err := _BeaconVRFConsumer.contract.UnpackLog(event, "CoordinatorUpdated", log); err != nil {
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

func (_BeaconVRFConsumer *BeaconVRFConsumerFilterer) ParseCoordinatorUpdated(log types.Log) (*BeaconVRFConsumerCoordinatorUpdated, error) {
	event := new(BeaconVRFConsumerCoordinatorUpdated)
	if err := _BeaconVRFConsumer.contract.UnpackLog(event, "CoordinatorUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BeaconVRFConsumerOwnershipTransferRequestedIterator struct {
	Event *BeaconVRFConsumerOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BeaconVRFConsumerOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BeaconVRFConsumerOwnershipTransferRequested)
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
		it.Event = new(BeaconVRFConsumerOwnershipTransferRequested)
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

func (it *BeaconVRFConsumerOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *BeaconVRFConsumerOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BeaconVRFConsumerOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_BeaconVRFConsumer *BeaconVRFConsumerFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BeaconVRFConsumerOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BeaconVRFConsumer.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &BeaconVRFConsumerOwnershipTransferRequestedIterator{contract: _BeaconVRFConsumer.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_BeaconVRFConsumer *BeaconVRFConsumerFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *BeaconVRFConsumerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BeaconVRFConsumer.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BeaconVRFConsumerOwnershipTransferRequested)
				if err := _BeaconVRFConsumer.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_BeaconVRFConsumer *BeaconVRFConsumerFilterer) ParseOwnershipTransferRequested(log types.Log) (*BeaconVRFConsumerOwnershipTransferRequested, error) {
	event := new(BeaconVRFConsumerOwnershipTransferRequested)
	if err := _BeaconVRFConsumer.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BeaconVRFConsumerOwnershipTransferredIterator struct {
	Event *BeaconVRFConsumerOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BeaconVRFConsumerOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BeaconVRFConsumerOwnershipTransferred)
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
		it.Event = new(BeaconVRFConsumerOwnershipTransferred)
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

func (it *BeaconVRFConsumerOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *BeaconVRFConsumerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BeaconVRFConsumerOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_BeaconVRFConsumer *BeaconVRFConsumerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BeaconVRFConsumerOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BeaconVRFConsumer.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &BeaconVRFConsumerOwnershipTransferredIterator{contract: _BeaconVRFConsumer.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_BeaconVRFConsumer *BeaconVRFConsumerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BeaconVRFConsumerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BeaconVRFConsumer.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BeaconVRFConsumerOwnershipTransferred)
				if err := _BeaconVRFConsumer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_BeaconVRFConsumer *BeaconVRFConsumerFilterer) ParseOwnershipTransferred(log types.Log) (*BeaconVRFConsumerOwnershipTransferred, error) {
	event := new(BeaconVRFConsumerOwnershipTransferred)
	if err := _BeaconVRFConsumer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
	case _BeaconVRFConsumer.abi.Events["CoordinatorUpdated"].ID:
		return _BeaconVRFConsumer.ParseCoordinatorUpdated(log)
	case _BeaconVRFConsumer.abi.Events["OwnershipTransferRequested"].ID:
		return _BeaconVRFConsumer.ParseOwnershipTransferRequested(log)
	case _BeaconVRFConsumer.abi.Events["OwnershipTransferred"].ID:
		return _BeaconVRFConsumer.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (BeaconVRFConsumerCoordinatorUpdated) Topic() common.Hash {
	return common.HexToHash("0xc258faa9a17ddfdf4130b4acff63a289202e7d5f9e42f366add65368575486bc")
}

func (BeaconVRFConsumerOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (BeaconVRFConsumerOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_BeaconVRFConsumer *BeaconVRFConsumer) Address() common.Address {
	return _BeaconVRFConsumer.address
}

type BeaconVRFConsumerInterface interface {
	Fail(opts *bind.CallOpts) (bool, error)

	IBeaconPeriodBlocks(opts *bind.CallOpts) (*big.Int, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SReceivedRandomnessByRequestID(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error)

	SArguments(opts *bind.CallOpts) ([]byte, error)

	SGasAvailable(opts *bind.CallOpts) (*big.Int, error)

	SMostRecentRequestID(opts *bind.CallOpts) (*big.Int, error)

	SMyBeaconRequests(opts *bind.CallOpts, arg0 *big.Int) (SMyBeaconRequests,

		error)

	SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	SRequestsIDs(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error)

	SSubId(opts *bind.CallOpts) (uint64, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestID *big.Int, randomWords []*big.Int, arguments []byte) (*types.Transaction, error)

	SetCoordinator(opts *bind.TransactOpts, coordinator common.Address) (*types.Transaction, error)

	SetFail(opts *bind.TransactOpts, shouldFail bool) (*types.Transaction, error)

	StoreBeaconRequest(opts *bind.TransactOpts, reqId *big.Int, height *big.Int, delay *big.Int, numWords uint16) (*types.Transaction, error)

	TestRedeemRandomness(opts *bind.TransactOpts, subID *big.Int, requestID *big.Int) (*types.Transaction, error)

	TestRequestRandomness(opts *bind.TransactOpts, numWords uint16, subID *big.Int, confirmationDelayArg *big.Int) (*types.Transaction, error)

	TestRequestRandomnessFulfillment(opts *bind.TransactOpts, subID *big.Int, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterCoordinatorUpdated(opts *bind.FilterOpts, coordinator []common.Address) (*BeaconVRFConsumerCoordinatorUpdatedIterator, error)

	WatchCoordinatorUpdated(opts *bind.WatchOpts, sink chan<- *BeaconVRFConsumerCoordinatorUpdated, coordinator []common.Address) (event.Subscription, error)

	ParseCoordinatorUpdated(log types.Log) (*BeaconVRFConsumerCoordinatorUpdated, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BeaconVRFConsumerOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *BeaconVRFConsumerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*BeaconVRFConsumerOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BeaconVRFConsumerOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BeaconVRFConsumerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*BeaconVRFConsumerOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
