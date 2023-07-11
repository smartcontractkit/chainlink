// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrfv2plus_consumer_example

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

var VRFV2PlusConsumerExampleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlyCoordinatorCanFulfill\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlySubOwnerCanSetVRFCoordinator\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"createSubscriptionAndFund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"idx\",\"type\":\"uint256\"}],\"name\":\"getRandomness\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"randomWord\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"name\":\"redeemRandomness\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"nativePayment\",\"type\":\"bool\"}],\"name\":\"requestRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_linkToken\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_recentRequestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_requests\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"fulfilled\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"cost\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_subId\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_vrfCoordinator\",\"outputs\":[{\"internalType\":\"contractIVRFCoordinatorV2Plus\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"subOwner\",\"type\":\"address\"}],\"name\":\"setSubOwner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vrfCoordinator\",\"type\":\"address\"}],\"name\":\"setVRFCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"topUpSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"name\":\"updateSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162001c5138038062001c51833981016040819052620000349162000213565b8133806000816200008c5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000bf57620000bf816200014a565b5050506001600160a01b038116620001095760405162461bcd60e51b815260206004820152600c60248201526b7a65726f206164647265737360a01b604482015260640162000083565b600280546001600160a01b03199081166001600160a01b0393841617909155600480548216948316949094179093556005805490931691161790556200024b565b6001600160a01b038116331415620001a55760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000083565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b03811681146200020e57600080fd5b919050565b600080604083850312156200022757600080fd5b6200023283620001f6565b91506200024260208401620001f6565b90509250929050565b6119f6806200025b6000396000f3fe60806040526004361061010e5760003560e01c806379ba5097116100a55780639eccacf611610074578063cf62c8ab11610059578063cf62c8ab146103c8578063eff27017146103e8578063f2fde38b1461040857600080fd5b80639eccacf614610302578063a168fa891461032f57600080fd5b806379ba5097146102825780637ec8773a146102975780638da5cb5b146102b75780639c7450ff146102e257600080fd5b806344ff81ce116100e157806344ff81ce146101a85780635d7d53e3146101c8578063706da1ca146101de5780637725135b1461023057600080fd5b80631fe543e31461011357806329e5d831146101355780632fa4e4421461016857806336bfffed14610188575b600080fd5b34801561011f57600080fd5b5061013361012e366004611611565b610428565b005b34801561014157600080fd5b506101556101503660046116b5565b6104ae565b6040519081526020015b60405180910390f35b34801561017457600080fd5b5061013361018336600461176c565b6105f1565b34801561019457600080fd5b506101336101a336600461151e565b6106c5565b3480156101b457600080fd5b506101336101c33660046114fc565b61084d565b3480156101d457600080fd5b5061015560065481565b3480156101ea57600080fd5b506005546102179074010000000000000000000000000000000000000000900467ffffffffffffffff1681565b60405167ffffffffffffffff909116815260200161015f565b34801561023c57600080fd5b5060055461025d9073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161015f565b34801561028e57600080fd5b50610133610907565b3480156102a357600080fd5b506101336102b23660046114fc565b610a04565b3480156102c357600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff1661025d565b6102f56102f03660046115df565b610a10565b60405161015f919061184c565b34801561030e57600080fd5b5060045461025d9073ffffffffffffffffffffffffffffffffffffffff1681565b34801561033b57600080fd5b5061038c61034a3660046115df565b600760205260009081526040902080546002820154600390920154909160ff81169161010090910473ffffffffffffffffffffffffffffffffffffffff169084565b60405161015f9493929190938452911515602084015273ffffffffffffffffffffffffffffffffffffffff166040830152606082015260800190565b3480156103d457600080fd5b506101336103e336600461176c565b610cba565b3480156103f457600080fd5b506101336104033660046116d7565b610e70565b34801561041457600080fd5b506101336104233660046114fc565b6110a0565b60025473ffffffffffffffffffffffffffffffffffffffff1633146104a0576002546040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff90911660248201526044015b60405180910390fd5b6104aa82826110b1565b5050565b6000828152600760209081526040808320815160a0810183528154815260018201805484518187028101870190955280855286959294858401939092919083018282801561051b57602002820191906000526020600020905b815481526020019060010190808311610507575b5050509183525050600282015460ff811615156020830152610100900473ffffffffffffffffffffffffffffffffffffffff16604082015260039091015460609091015280519091506105ca576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265717565737420494420697320696e636f72726563740000000000000000006044820152606401610497565b806020015183815181106105e0576105e0611963565b602002602001015191505092915050565b6005546004546040805174010000000000000000000000000000000000000000840467ffffffffffffffff16602082015273ffffffffffffffffffffffffffffffffffffffff93841693634000aea09316918591016040516020818303038152906040526040518463ffffffff1660e01b8152600401610673939291906117a6565b602060405180830381600087803b15801561068d57600080fd5b505af11580156106a1573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104aa91906115c2565b60055474010000000000000000000000000000000000000000900467ffffffffffffffff16610750576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f7375624944206e6f7420736574000000000000000000000000000000000000006044820152606401610497565b60005b81518110156104aa57600454600554835173ffffffffffffffffffffffffffffffffffffffff90921691637341c10c9174010000000000000000000000000000000000000000900467ffffffffffffffff16908590859081106107b8576107b8611963565b60200260200101516040518363ffffffff1660e01b815260040161080892919067ffffffffffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b600060405180830381600087803b15801561082257600080fd5b505af1158015610836573d6000803e3d6000fd5b50505050808061084590611903565b915050610753565b60035473ffffffffffffffffffffffffffffffffffffffff1633146108c0576003546040517f4ae338ff00000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff9091166024820152604401610497565b600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b60015473ffffffffffffffffffffffffffffffffffffffff163314610988576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610497565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610a0d81611252565b50565b6000818152600760209081526040808320815160a081018352815481526001820180548451818702810187019095528085526060969592948584019390929190830182828015610a7f57602002820191906000526020600020905b815481526020019060010190808311610a6b575b5050509183525050600282015460ff811615156020830152610100900473ffffffffffffffffffffffffffffffffffffffff1660408201526003909101546060909101528051909150610b2e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265717565737420494420697320696e636f72726563740000000000000000006044820152606401610497565b8060400151610b99576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f72657175657374206e6f742066756c66696c6c656420796574000000000000006044820152606401610497565b606081015173ffffffffffffffffffffffffffffffffffffffff163314610c42576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602360248201527f6f6e6c792063616c6c61626c652062792072657175657374696e67206164647260448201527f65737300000000000000000000000000000000000000000000000000000000006064820152608401610497565b8060800151341015610cb0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e73756666696369656e742066756e647300000000000000000000000000006044820152606401610497565b6020015192915050565b60055474010000000000000000000000000000000000000000900467ffffffffffffffff166105f15760048054604080517fa21a23e4000000000000000000000000000000000000000000000000000000008152905173ffffffffffffffffffffffffffffffffffffffff9092169263a21a23e49282820192602092908290030181600087803b158015610d4d57600080fd5b505af1158015610d61573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d859190611742565b600580547fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000067ffffffffffffffff93841681029190911791829055600480546040517f7341c10c000000000000000000000000000000000000000000000000000000008152929093049093169281019290925230602483015273ffffffffffffffffffffffffffffffffffffffff1690637341c10c90604401600060405180830381600087803b158015610e4f57600080fd5b505af1158015610e63573d6000803e3d6000fd5b505050506105f130611252565b600480546005546040517fefcf1d9400000000000000000000000000000000000000000000000000000000815292830185905274010000000000000000000000000000000000000000900467ffffffffffffffff16602483015261ffff8616604483015263ffffffff80881660648401528516608483015282151560a483015260009173ffffffffffffffffffffffffffffffffffffffff9091169063efcf1d949060c401602060405180830381600087803b158015610f2f57600080fd5b505af1158015610f43573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610f6791906115f8565b905060006040518060a00160405280838152602001600067ffffffffffffffff811115610f9657610f96611992565b604051908082528060200260200182016040528015610fbf578160200160208202803683370190505b5081526000602080830182905233604080850191909152606090930182905285825260078152919020825181558282015180519394508493919261100b9260018501929091019061145f565b506040820151600282018054606085015173ffffffffffffffffffffffffffffffffffffffff16610100027fffffffffffffffffffffff0000000000000000000000000000000000000000ff931515939093167fffffffffffffffffffffff00000000000000000000000000000000000000000090911617919091179055608090910151600390910155506006555050505050565b6110a86112e6565b610a0d81611369565b600654821461111c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265717565737420494420697320696e636f72726563740000000000000000006044820152606401610497565b600082815260076020908152604090912082516111419260019092019184019061145f565b506000828152600760205260409081902060020180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790556004805491517ffd26ba4b00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9092169163fd26ba4b916111d69186910190815260200190565b60206040518083038186803b1580156111ee57600080fd5b505afa158015611202573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112269190611789565b6bffffffffffffffffffffffff1660076000848152602001908152602001600020600301819055505050565b73ffffffffffffffffffffffffffffffffffffffff811661129f576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600380547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b60005473ffffffffffffffffffffffffffffffffffffffff163314611367576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610497565b565b73ffffffffffffffffffffffffffffffffffffffff81163314156113e9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610497565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b82805482825590600052602060002090810192821561149a579160200282015b8281111561149a57825182559160200191906001019061147f565b506114a69291506114aa565b5090565b5b808211156114a657600081556001016114ab565b803573ffffffffffffffffffffffffffffffffffffffff811681146114e357600080fd5b919050565b803563ffffffff811681146114e357600080fd5b60006020828403121561150e57600080fd5b611517826114bf565b9392505050565b6000602080838503121561153157600080fd5b823567ffffffffffffffff81111561154857600080fd5b8301601f8101851361155957600080fd5b803561156c611567826118df565b611890565b80828252848201915084840188868560051b870101111561158c57600080fd5b600094505b838510156115b6576115a2816114bf565b835260019490940193918501918501611591565b50979650505050505050565b6000602082840312156115d457600080fd5b8151611517816119c1565b6000602082840312156115f157600080fd5b5035919050565b60006020828403121561160a57600080fd5b5051919050565b6000806040838503121561162457600080fd5b8235915060208084013567ffffffffffffffff81111561164357600080fd5b8401601f8101861361165457600080fd5b8035611662611567826118df565b80828252848201915084840189868560051b870101111561168257600080fd5b600094505b838510156116a5578035835260019490940193918501918501611687565b5080955050505050509250929050565b600080604083850312156116c857600080fd5b50508035926020909101359150565b600080600080600060a086880312156116ef57600080fd5b6116f8866114e8565b9450602086013561ffff8116811461170f57600080fd5b935061171d604087016114e8565b9250606086013591506080860135611734816119c1565b809150509295509295909350565b60006020828403121561175457600080fd5b815167ffffffffffffffff8116811461151757600080fd5b60006020828403121561177e57600080fd5b8135611517816119cf565b60006020828403121561179b57600080fd5b8151611517816119cf565b73ffffffffffffffffffffffffffffffffffffffff84168152600060206bffffffffffffffffffffffff85168184015260606040840152835180606085015260005b81811015611804578581018301518582016080015282016117e8565b81811115611816576000608083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160800195945050505050565b6020808252825182820181905260009190848201906040850190845b8181101561188457835183529284019291840191600101611868565b50909695505050505050565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156118d7576118d7611992565b604052919050565b600067ffffffffffffffff8211156118f9576118f9611992565b5060051b60200190565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82141561195c577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b8015158114610a0d57600080fd5b6bffffffffffffffffffffffff81168114610a0d57600080fdfea164736f6c6343000806000a",
}

var VRFV2PlusConsumerExampleABI = VRFV2PlusConsumerExampleMetaData.ABI

var VRFV2PlusConsumerExampleBin = VRFV2PlusConsumerExampleMetaData.Bin

func DeployVRFV2PlusConsumerExample(auth *bind.TransactOpts, backend bind.ContractBackend, vrfCoordinator common.Address, link common.Address) (common.Address, *types.Transaction, *VRFV2PlusConsumerExample, error) {
	parsed, err := VRFV2PlusConsumerExampleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFV2PlusConsumerExampleBin), backend, vrfCoordinator, link)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFV2PlusConsumerExample{VRFV2PlusConsumerExampleCaller: VRFV2PlusConsumerExampleCaller{contract: contract}, VRFV2PlusConsumerExampleTransactor: VRFV2PlusConsumerExampleTransactor{contract: contract}, VRFV2PlusConsumerExampleFilterer: VRFV2PlusConsumerExampleFilterer{contract: contract}}, nil
}

type VRFV2PlusConsumerExample struct {
	address common.Address
	abi     abi.ABI
	VRFV2PlusConsumerExampleCaller
	VRFV2PlusConsumerExampleTransactor
	VRFV2PlusConsumerExampleFilterer
}

type VRFV2PlusConsumerExampleCaller struct {
	contract *bind.BoundContract
}

type VRFV2PlusConsumerExampleTransactor struct {
	contract *bind.BoundContract
}

type VRFV2PlusConsumerExampleFilterer struct {
	contract *bind.BoundContract
}

type VRFV2PlusConsumerExampleSession struct {
	Contract     *VRFV2PlusConsumerExample
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFV2PlusConsumerExampleCallerSession struct {
	Contract *VRFV2PlusConsumerExampleCaller
	CallOpts bind.CallOpts
}

type VRFV2PlusConsumerExampleTransactorSession struct {
	Contract     *VRFV2PlusConsumerExampleTransactor
	TransactOpts bind.TransactOpts
}

type VRFV2PlusConsumerExampleRaw struct {
	Contract *VRFV2PlusConsumerExample
}

type VRFV2PlusConsumerExampleCallerRaw struct {
	Contract *VRFV2PlusConsumerExampleCaller
}

type VRFV2PlusConsumerExampleTransactorRaw struct {
	Contract *VRFV2PlusConsumerExampleTransactor
}

func NewVRFV2PlusConsumerExample(address common.Address, backend bind.ContractBackend) (*VRFV2PlusConsumerExample, error) {
	abi, err := abi.JSON(strings.NewReader(VRFV2PlusConsumerExampleABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFV2PlusConsumerExample(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusConsumerExample{address: address, abi: abi, VRFV2PlusConsumerExampleCaller: VRFV2PlusConsumerExampleCaller{contract: contract}, VRFV2PlusConsumerExampleTransactor: VRFV2PlusConsumerExampleTransactor{contract: contract}, VRFV2PlusConsumerExampleFilterer: VRFV2PlusConsumerExampleFilterer{contract: contract}}, nil
}

func NewVRFV2PlusConsumerExampleCaller(address common.Address, caller bind.ContractCaller) (*VRFV2PlusConsumerExampleCaller, error) {
	contract, err := bindVRFV2PlusConsumerExample(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusConsumerExampleCaller{contract: contract}, nil
}

func NewVRFV2PlusConsumerExampleTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFV2PlusConsumerExampleTransactor, error) {
	contract, err := bindVRFV2PlusConsumerExample(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusConsumerExampleTransactor{contract: contract}, nil
}

func NewVRFV2PlusConsumerExampleFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFV2PlusConsumerExampleFilterer, error) {
	contract, err := bindVRFV2PlusConsumerExample(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusConsumerExampleFilterer{contract: contract}, nil
}

func bindVRFV2PlusConsumerExample(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFV2PlusConsumerExampleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2PlusConsumerExample.Contract.VRFV2PlusConsumerExampleCaller.contract.Call(opts, result, method, params...)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.VRFV2PlusConsumerExampleTransactor.contract.Transfer(opts)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.VRFV2PlusConsumerExampleTransactor.contract.Transact(opts, method, params...)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2PlusConsumerExample.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.contract.Transfer(opts)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.contract.Transact(opts, method, params...)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleCaller) GetRandomness(opts *bind.CallOpts, requestId *big.Int, idx *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusConsumerExample.contract.Call(opts, &out, "getRandomness", requestId, idx)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleSession) GetRandomness(requestId *big.Int, idx *big.Int) (*big.Int, error) {
	return _VRFV2PlusConsumerExample.Contract.GetRandomness(&_VRFV2PlusConsumerExample.CallOpts, requestId, idx)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleCallerSession) GetRandomness(requestId *big.Int, idx *big.Int) (*big.Int, error) {
	return _VRFV2PlusConsumerExample.Contract.GetRandomness(&_VRFV2PlusConsumerExample.CallOpts, requestId, idx)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2PlusConsumerExample.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleSession) Owner() (common.Address, error) {
	return _VRFV2PlusConsumerExample.Contract.Owner(&_VRFV2PlusConsumerExample.CallOpts)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleCallerSession) Owner() (common.Address, error) {
	return _VRFV2PlusConsumerExample.Contract.Owner(&_VRFV2PlusConsumerExample.CallOpts)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleCaller) SLinkToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2PlusConsumerExample.contract.Call(opts, &out, "s_linkToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleSession) SLinkToken() (common.Address, error) {
	return _VRFV2PlusConsumerExample.Contract.SLinkToken(&_VRFV2PlusConsumerExample.CallOpts)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleCallerSession) SLinkToken() (common.Address, error) {
	return _VRFV2PlusConsumerExample.Contract.SLinkToken(&_VRFV2PlusConsumerExample.CallOpts)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleCaller) SRecentRequestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusConsumerExample.contract.Call(opts, &out, "s_recentRequestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleSession) SRecentRequestId() (*big.Int, error) {
	return _VRFV2PlusConsumerExample.Contract.SRecentRequestId(&_VRFV2PlusConsumerExample.CallOpts)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleCallerSession) SRecentRequestId() (*big.Int, error) {
	return _VRFV2PlusConsumerExample.Contract.SRecentRequestId(&_VRFV2PlusConsumerExample.CallOpts)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleCaller) SRequests(opts *bind.CallOpts, arg0 *big.Int) (SRequests,

	error) {
	var out []interface{}
	err := _VRFV2PlusConsumerExample.contract.Call(opts, &out, "s_requests", arg0)

	outstruct := new(SRequests)
	if err != nil {
		return *outstruct, err
	}

	outstruct.RequestId = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Fulfilled = *abi.ConvertType(out[1], new(bool)).(*bool)
	outstruct.Requester = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.Cost = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleSession) SRequests(arg0 *big.Int) (SRequests,

	error) {
	return _VRFV2PlusConsumerExample.Contract.SRequests(&_VRFV2PlusConsumerExample.CallOpts, arg0)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleCallerSession) SRequests(arg0 *big.Int) (SRequests,

	error) {
	return _VRFV2PlusConsumerExample.Contract.SRequests(&_VRFV2PlusConsumerExample.CallOpts, arg0)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleCaller) SSubId(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _VRFV2PlusConsumerExample.contract.Call(opts, &out, "s_subId")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleSession) SSubId() (uint64, error) {
	return _VRFV2PlusConsumerExample.Contract.SSubId(&_VRFV2PlusConsumerExample.CallOpts)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleCallerSession) SSubId() (uint64, error) {
	return _VRFV2PlusConsumerExample.Contract.SSubId(&_VRFV2PlusConsumerExample.CallOpts)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleCaller) SVrfCoordinator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2PlusConsumerExample.contract.Call(opts, &out, "s_vrfCoordinator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleSession) SVrfCoordinator() (common.Address, error) {
	return _VRFV2PlusConsumerExample.Contract.SVrfCoordinator(&_VRFV2PlusConsumerExample.CallOpts)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleCallerSession) SVrfCoordinator() (common.Address, error) {
	return _VRFV2PlusConsumerExample.Contract.SVrfCoordinator(&_VRFV2PlusConsumerExample.CallOpts)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.contract.Transact(opts, "acceptOwnership")
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.AcceptOwnership(&_VRFV2PlusConsumerExample.TransactOpts)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.AcceptOwnership(&_VRFV2PlusConsumerExample.TransactOpts)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactor) CreateSubscriptionAndFund(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.contract.Transact(opts, "createSubscriptionAndFund", amount)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleSession) CreateSubscriptionAndFund(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.CreateSubscriptionAndFund(&_VRFV2PlusConsumerExample.TransactOpts, amount)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactorSession) CreateSubscriptionAndFund(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.CreateSubscriptionAndFund(&_VRFV2PlusConsumerExample.TransactOpts, amount)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.contract.Transact(opts, "rawFulfillRandomWords", requestId, randomWords)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.RawFulfillRandomWords(&_VRFV2PlusConsumerExample.TransactOpts, requestId, randomWords)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactorSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.RawFulfillRandomWords(&_VRFV2PlusConsumerExample.TransactOpts, requestId, randomWords)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactor) RedeemRandomness(opts *bind.TransactOpts, requestId *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.contract.Transact(opts, "redeemRandomness", requestId)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleSession) RedeemRandomness(requestId *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.RedeemRandomness(&_VRFV2PlusConsumerExample.TransactOpts, requestId)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactorSession) RedeemRandomness(requestId *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.RedeemRandomness(&_VRFV2PlusConsumerExample.TransactOpts, requestId)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactor) RequestRandomWords(opts *bind.TransactOpts, callbackGasLimit uint32, requestConfirmations uint16, numWords uint32, keyHash [32]byte, nativePayment bool) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.contract.Transact(opts, "requestRandomWords", callbackGasLimit, requestConfirmations, numWords, keyHash, nativePayment)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleSession) RequestRandomWords(callbackGasLimit uint32, requestConfirmations uint16, numWords uint32, keyHash [32]byte, nativePayment bool) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.RequestRandomWords(&_VRFV2PlusConsumerExample.TransactOpts, callbackGasLimit, requestConfirmations, numWords, keyHash, nativePayment)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactorSession) RequestRandomWords(callbackGasLimit uint32, requestConfirmations uint16, numWords uint32, keyHash [32]byte, nativePayment bool) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.RequestRandomWords(&_VRFV2PlusConsumerExample.TransactOpts, callbackGasLimit, requestConfirmations, numWords, keyHash, nativePayment)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactor) SetSubOwner(opts *bind.TransactOpts, subOwner common.Address) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.contract.Transact(opts, "setSubOwner", subOwner)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleSession) SetSubOwner(subOwner common.Address) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.SetSubOwner(&_VRFV2PlusConsumerExample.TransactOpts, subOwner)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactorSession) SetSubOwner(subOwner common.Address) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.SetSubOwner(&_VRFV2PlusConsumerExample.TransactOpts, subOwner)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactor) SetVRFCoordinator(opts *bind.TransactOpts, _vrfCoordinator common.Address) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.contract.Transact(opts, "setVRFCoordinator", _vrfCoordinator)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleSession) SetVRFCoordinator(_vrfCoordinator common.Address) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.SetVRFCoordinator(&_VRFV2PlusConsumerExample.TransactOpts, _vrfCoordinator)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactorSession) SetVRFCoordinator(_vrfCoordinator common.Address) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.SetVRFCoordinator(&_VRFV2PlusConsumerExample.TransactOpts, _vrfCoordinator)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactor) TopUpSubscription(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.contract.Transact(opts, "topUpSubscription", amount)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleSession) TopUpSubscription(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.TopUpSubscription(&_VRFV2PlusConsumerExample.TransactOpts, amount)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactorSession) TopUpSubscription(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.TopUpSubscription(&_VRFV2PlusConsumerExample.TransactOpts, amount)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.contract.Transact(opts, "transferOwnership", to)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.TransferOwnership(&_VRFV2PlusConsumerExample.TransactOpts, to)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.TransferOwnership(&_VRFV2PlusConsumerExample.TransactOpts, to)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactor) UpdateSubscription(opts *bind.TransactOpts, consumers []common.Address) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.contract.Transact(opts, "updateSubscription", consumers)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleSession) UpdateSubscription(consumers []common.Address) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.UpdateSubscription(&_VRFV2PlusConsumerExample.TransactOpts, consumers)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactorSession) UpdateSubscription(consumers []common.Address) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.UpdateSubscription(&_VRFV2PlusConsumerExample.TransactOpts, consumers)
}

type VRFV2PlusConsumerExampleOwnershipTransferRequestedIterator struct {
	Event *VRFV2PlusConsumerExampleOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2PlusConsumerExampleOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2PlusConsumerExampleOwnershipTransferRequested)
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
		it.Event = new(VRFV2PlusConsumerExampleOwnershipTransferRequested)
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

func (it *VRFV2PlusConsumerExampleOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFV2PlusConsumerExampleOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2PlusConsumerExampleOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusConsumerExampleOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusConsumerExample.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusConsumerExampleOwnershipTransferRequestedIterator{contract: _VRFV2PlusConsumerExample.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFV2PlusConsumerExampleOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusConsumerExample.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2PlusConsumerExampleOwnershipTransferRequested)
				if err := _VRFV2PlusConsumerExample.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleFilterer) ParseOwnershipTransferRequested(log types.Log) (*VRFV2PlusConsumerExampleOwnershipTransferRequested, error) {
	event := new(VRFV2PlusConsumerExampleOwnershipTransferRequested)
	if err := _VRFV2PlusConsumerExample.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFV2PlusConsumerExampleOwnershipTransferredIterator struct {
	Event *VRFV2PlusConsumerExampleOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2PlusConsumerExampleOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2PlusConsumerExampleOwnershipTransferred)
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
		it.Event = new(VRFV2PlusConsumerExampleOwnershipTransferred)
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

func (it *VRFV2PlusConsumerExampleOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFV2PlusConsumerExampleOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2PlusConsumerExampleOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusConsumerExampleOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusConsumerExample.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusConsumerExampleOwnershipTransferredIterator{contract: _VRFV2PlusConsumerExample.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFV2PlusConsumerExampleOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusConsumerExample.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2PlusConsumerExampleOwnershipTransferred)
				if err := _VRFV2PlusConsumerExample.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleFilterer) ParseOwnershipTransferred(log types.Log) (*VRFV2PlusConsumerExampleOwnershipTransferred, error) {
	event := new(VRFV2PlusConsumerExampleOwnershipTransferred)
	if err := _VRFV2PlusConsumerExample.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type SRequests struct {
	RequestId *big.Int
	Fulfilled bool
	Requester common.Address
	Cost      *big.Int
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExample) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFV2PlusConsumerExample.abi.Events["OwnershipTransferRequested"].ID:
		return _VRFV2PlusConsumerExample.ParseOwnershipTransferRequested(log)
	case _VRFV2PlusConsumerExample.abi.Events["OwnershipTransferred"].ID:
		return _VRFV2PlusConsumerExample.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFV2PlusConsumerExampleOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VRFV2PlusConsumerExampleOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExample) Address() common.Address {
	return _VRFV2PlusConsumerExample.address
}

type VRFV2PlusConsumerExampleInterface interface {
	GetRandomness(opts *bind.CallOpts, requestId *big.Int, idx *big.Int) (*big.Int, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SLinkToken(opts *bind.CallOpts) (common.Address, error)

	SRecentRequestId(opts *bind.CallOpts) (*big.Int, error)

	SRequests(opts *bind.CallOpts, arg0 *big.Int) (SRequests,

		error)

	SSubId(opts *bind.CallOpts) (uint64, error)

	SVrfCoordinator(opts *bind.CallOpts) (common.Address, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	CreateSubscriptionAndFund(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	RedeemRandomness(opts *bind.TransactOpts, requestId *big.Int) (*types.Transaction, error)

	RequestRandomWords(opts *bind.TransactOpts, callbackGasLimit uint32, requestConfirmations uint16, numWords uint32, keyHash [32]byte, nativePayment bool) (*types.Transaction, error)

	SetSubOwner(opts *bind.TransactOpts, subOwner common.Address) (*types.Transaction, error)

	SetVRFCoordinator(opts *bind.TransactOpts, _vrfCoordinator common.Address) (*types.Transaction, error)

	TopUpSubscription(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	UpdateSubscription(opts *bind.TransactOpts, consumers []common.Address) (*types.Transaction, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusConsumerExampleOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFV2PlusConsumerExampleOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VRFV2PlusConsumerExampleOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusConsumerExampleOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFV2PlusConsumerExampleOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VRFV2PlusConsumerExampleOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
