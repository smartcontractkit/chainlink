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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlyCoordinatorCanFulfill\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"}],\"name\":\"OnlyOwnerOrCoordinator\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"createSubscriptionAndFund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createSubscriptionAndFundNative\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"idx\",\"type\":\"uint256\"}],\"name\":\"getRandomness\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"randomWord\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"nativePayment\",\"type\":\"bool\"}],\"name\":\"requestRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_linkToken\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_recentRequestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_requests\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"fulfilled\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_subId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_vrfCoordinator\",\"outputs\":[{\"internalType\":\"contractIVRFCoordinatorV2Plus\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_vrfCoordinatorApiV1\",\"outputs\":[{\"internalType\":\"contractIVRFCoordinatorV2Plus\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vrfCoordinator\",\"type\":\"address\"}],\"name\":\"setCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"setSubId\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"topUpSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"topUpSubscriptionNative\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"name\":\"updateSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b50604051620019c6380380620019c68339810160408190526200003491620001cc565b8133806000816200008c5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000bf57620000bf8162000103565b5050600280546001600160a01b03199081166001600160a01b0394851617909155600580548216958416959095179094555060038054909316911617905562000204565b6001600160a01b0381163314156200015e5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000083565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b0381168114620001c757600080fd5b919050565b60008060408385031215620001e057600080fd5b620001eb83620001af565b9150620001fb60208401620001af565b90509250929050565b6117b280620002146000396000f3fe6080604052600436106101445760003560e01c806380980043116100c0578063b96dbba711610074578063de367c8e11610059578063de367c8e146103c0578063eff27017146103ed578063f2fde38b1461040d57600080fd5b8063b96dbba714610398578063cf62c8ab146103a057600080fd5b80638ea98117116100a55780638ea98117146102c45780639eccacf6146102e4578063a168fa891461031157600080fd5b806380980043146102795780638da5cb5b1461029957600080fd5b806336bfffed11610117578063706da1ca116100fc578063706da1ca146101fc5780637725135b1461021257806379ba50971461026457600080fd5b806336bfffed146101c65780635d7d53e3146101e657600080fd5b80631d2b2afd146101495780631fe543e31461015357806329e5d831146101735780632fa4e442146101a6575b600080fd5b61015161042d565b005b34801561015f57600080fd5b5061015161016e3660046113eb565b610528565b34801561017f57600080fd5b5061019361018e36600461148f565b6105a9565b6040519081526020015b60405180910390f35b3480156101b257600080fd5b506101516101c136600461151c565b6106e6565b3480156101d257600080fd5b506101516101e13660046112f8565b610808565b3480156101f257600080fd5b5061019360045481565b34801561020857600080fd5b5061019360065481565b34801561021e57600080fd5b5060035461023f9073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161019d565b34801561027057600080fd5b50610151610940565b34801561028557600080fd5b506101516102943660046113b9565b600655565b3480156102a557600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff1661023f565b3480156102d057600080fd5b506101516102df3660046112d6565b610a3d565b3480156102f057600080fd5b5060025461023f9073ffffffffffffffffffffffffffffffffffffffff1681565b34801561031d57600080fd5b5061036661032c3660046113b9565b6007602052600090815260409020805460019091015460ff821691610100900473ffffffffffffffffffffffffffffffffffffffff169083565b60408051931515845273ffffffffffffffffffffffffffffffffffffffff90921660208401529082015260600161019d565b610151610b48565b3480156103ac57600080fd5b506101516103bb36600461151c565b610bae565b3480156103cc57600080fd5b5060055461023f9073ffffffffffffffffffffffffffffffffffffffff1681565b3480156103f957600080fd5b506101516104083660046114b1565b610bf5565b34801561041957600080fd5b506101516104283660046112d6565b610de0565b60065461049b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f737562206e6f742073657400000000000000000000000000000000000000000060448201526064015b60405180910390fd5b6005546006546040517f95b55cfc000000000000000000000000000000000000000000000000000000008152600481019190915273ffffffffffffffffffffffffffffffffffffffff909116906395b55cfc9034906024015b6000604051808303818588803b15801561050d57600080fd5b505af1158015610521573d6000803e3d6000fd5b5050505050565b60025473ffffffffffffffffffffffffffffffffffffffff16331461059b576002546040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff9091166024820152604401610492565b6105a58282610df4565b5050565b60008281526007602090815260408083208151608081018352815460ff811615158252610100900473ffffffffffffffffffffffffffffffffffffffff16818501526001820154818401526002820180548451818702810187019095528085528695929460608601939092919083018282801561064557602002820191906000526020600020905b815481526020019060010190808311610631575b50505050508152505090508060400151600014156106bf576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265717565737420494420697320696e636f72726563740000000000000000006044820152606401610492565b806060015183815181106106d5576106d5611739565b602002602001015191505092915050565b60065461074f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f737562206e6f74207365740000000000000000000000000000000000000000006044820152606401610492565b60035460025460065460408051602081019290925273ffffffffffffffffffffffffffffffffffffffff93841693634000aea09316918591015b6040516020818303038152906040526040518463ffffffff1660e01b81526004016107b6939291906115b5565b602060405180830381600087803b1580156107d057600080fd5b505af11580156107e4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105a5919061139c565b600654610871576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f7375624944206e6f7420736574000000000000000000000000000000000000006044820152606401610492565b60005b81518110156105a557600554600654835173ffffffffffffffffffffffffffffffffffffffff9092169163bec4c08c91908590859081106108b7576108b7611739565b60200260200101516040518363ffffffff1660e01b81526004016108fb92919091825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b600060405180830381600087803b15801561091557600080fd5b505af1158015610929573d6000803e3d6000fd5b505050508080610938906116d9565b915050610874565b60015473ffffffffffffffffffffffffffffffffffffffff1633146109c1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610492565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60005473ffffffffffffffffffffffffffffffffffffffff163314801590610a7d575060025473ffffffffffffffffffffffffffffffffffffffff163314155b15610b015733610aa260005473ffffffffffffffffffffffffffffffffffffffff1690565b6002546040517f061db9c100000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff93841660048201529183166024830152919091166044820152606401610492565b600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b610b50610ebf565b506005546006546040517f95b55cfc000000000000000000000000000000000000000000000000000000008152600481019190915273ffffffffffffffffffffffffffffffffffffffff909116906395b55cfc9034906024016104f4565b610bb6610ebf565b5060035460025460065460408051602081019290925273ffffffffffffffffffffffffffffffffffffffff93841693634000aea0931691859101610789565b60006040518060c0016040528084815260200160065481526020018661ffff1681526020018763ffffffff1681526020018563ffffffff168152602001610c4b6040518060200160405280861515815250611004565b90526002546040517f9b1c385e00000000000000000000000000000000000000000000000000000000815291925060009173ffffffffffffffffffffffffffffffffffffffff90911690639b1c385e90610ca9908590600401611601565b602060405180830381600087803b158015610cc357600080fd5b505af1158015610cd7573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610cfb91906113d2565b604080516080810182526000808252336020808401918252838501868152855184815280830187526060860190815287855260078352959093208451815493517fffffffffffffffffffffff0000000000000000000000000000000000000000009094169015157fffffffffffffffffffffff0000000000000000000000000000000000000000ff161761010073ffffffffffffffffffffffffffffffffffffffff9094169390930292909217825591516001820155925180519495509193849392610dce926002850192910190611239565b50505060049190915550505050505050565b610de86110c0565b610df181611143565b50565b6004548214610e5f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265717565737420494420697320696e636f72726563740000000000000000006044820152606401610492565b60008281526007602090815260409091208251610e8492600290920191840190611239565b5050600090815260076020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055565b600060065460001415610ffd57600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a21a23e46040518163ffffffff1660e01b8152600401602060405180830381600087803b158015610f3657600080fd5b505af1158015610f4a573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610f6e91906113d2565b60068190556005546040517fbec4c08c000000000000000000000000000000000000000000000000000000008152600481019290925230602483015273ffffffffffffffffffffffffffffffffffffffff169063bec4c08c90604401600060405180830381600087803b158015610fe457600080fd5b505af1158015610ff8573d6000803e3d6000fd5b505050505b5060065490565b60607f92fd13387c7fe7befbc38d303d6468778fb9731bc4583f17d92989c6fcfdeaaa8260405160240161103d91511515815260200190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009093169290921790915292915050565b60005473ffffffffffffffffffffffffffffffffffffffff163314611141576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610492565b565b73ffffffffffffffffffffffffffffffffffffffff81163314156111c3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610492565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b828054828255906000526020600020908101928215611274579160200282015b82811115611274578251825591602001919060010190611259565b50611280929150611284565b5090565b5b808211156112805760008155600101611285565b803573ffffffffffffffffffffffffffffffffffffffff811681146112bd57600080fd5b919050565b803563ffffffff811681146112bd57600080fd5b6000602082840312156112e857600080fd5b6112f182611299565b9392505050565b6000602080838503121561130b57600080fd5b823567ffffffffffffffff81111561132257600080fd5b8301601f8101851361133357600080fd5b8035611346611341826116b5565b611666565b80828252848201915084840188868560051b870101111561136657600080fd5b600094505b838510156113905761137c81611299565b83526001949094019391850191850161136b565b50979650505050505050565b6000602082840312156113ae57600080fd5b81516112f181611797565b6000602082840312156113cb57600080fd5b5035919050565b6000602082840312156113e457600080fd5b5051919050565b600080604083850312156113fe57600080fd5b8235915060208084013567ffffffffffffffff81111561141d57600080fd5b8401601f8101861361142e57600080fd5b803561143c611341826116b5565b80828252848201915084840189868560051b870101111561145c57600080fd5b600094505b8385101561147f578035835260019490940193918501918501611461565b5080955050505050509250929050565b600080604083850312156114a257600080fd5b50508035926020909101359150565b600080600080600060a086880312156114c957600080fd5b6114d2866112c2565b9450602086013561ffff811681146114e957600080fd5b93506114f7604087016112c2565b925060608601359150608086013561150e81611797565b809150509295509295909350565b60006020828403121561152e57600080fd5b81356bffffffffffffffffffffffff811681146112f157600080fd5b6000815180845260005b8181101561157057602081850181015186830182015201611554565b81811115611582576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b73ffffffffffffffffffffffffffffffffffffffff841681526bffffffffffffffffffffffff831660208201526060604082015260006115f8606083018461154a565b95945050505050565b60208152815160208201526020820151604082015261ffff60408301511660608201526000606083015163ffffffff80821660808501528060808601511660a0850152505060a083015160c08084015261165e60e084018261154a565b949350505050565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156116ad576116ad611768565b604052919050565b600067ffffffffffffffff8211156116cf576116cf611768565b5060051b60200190565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415611732577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b8015158114610df157600080fdfea164736f6c6343000806000a",
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
	return address, tx, &VRFV2PlusConsumerExample{address: address, abi: *parsed, VRFV2PlusConsumerExampleCaller: VRFV2PlusConsumerExampleCaller{contract: contract}, VRFV2PlusConsumerExampleTransactor: VRFV2PlusConsumerExampleTransactor{contract: contract}, VRFV2PlusConsumerExampleFilterer: VRFV2PlusConsumerExampleFilterer{contract: contract}}, nil
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

	outstruct.Fulfilled = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Requester = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.RequestId = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

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

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleCaller) SSubId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusConsumerExample.contract.Call(opts, &out, "s_subId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleSession) SSubId() (*big.Int, error) {
	return _VRFV2PlusConsumerExample.Contract.SSubId(&_VRFV2PlusConsumerExample.CallOpts)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleCallerSession) SSubId() (*big.Int, error) {
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

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleCaller) SVrfCoordinatorApiV1(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2PlusConsumerExample.contract.Call(opts, &out, "s_vrfCoordinatorApiV1")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleSession) SVrfCoordinatorApiV1() (common.Address, error) {
	return _VRFV2PlusConsumerExample.Contract.SVrfCoordinatorApiV1(&_VRFV2PlusConsumerExample.CallOpts)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleCallerSession) SVrfCoordinatorApiV1() (common.Address, error) {
	return _VRFV2PlusConsumerExample.Contract.SVrfCoordinatorApiV1(&_VRFV2PlusConsumerExample.CallOpts)
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

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactor) CreateSubscriptionAndFundNative(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.contract.Transact(opts, "createSubscriptionAndFundNative")
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleSession) CreateSubscriptionAndFundNative() (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.CreateSubscriptionAndFundNative(&_VRFV2PlusConsumerExample.TransactOpts)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactorSession) CreateSubscriptionAndFundNative() (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.CreateSubscriptionAndFundNative(&_VRFV2PlusConsumerExample.TransactOpts)
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

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactor) RequestRandomWords(opts *bind.TransactOpts, callbackGasLimit uint32, requestConfirmations uint16, numWords uint32, keyHash [32]byte, nativePayment bool) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.contract.Transact(opts, "requestRandomWords", callbackGasLimit, requestConfirmations, numWords, keyHash, nativePayment)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleSession) RequestRandomWords(callbackGasLimit uint32, requestConfirmations uint16, numWords uint32, keyHash [32]byte, nativePayment bool) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.RequestRandomWords(&_VRFV2PlusConsumerExample.TransactOpts, callbackGasLimit, requestConfirmations, numWords, keyHash, nativePayment)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactorSession) RequestRandomWords(callbackGasLimit uint32, requestConfirmations uint16, numWords uint32, keyHash [32]byte, nativePayment bool) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.RequestRandomWords(&_VRFV2PlusConsumerExample.TransactOpts, callbackGasLimit, requestConfirmations, numWords, keyHash, nativePayment)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactor) SetCoordinator(opts *bind.TransactOpts, _vrfCoordinator common.Address) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.contract.Transact(opts, "setCoordinator", _vrfCoordinator)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleSession) SetCoordinator(_vrfCoordinator common.Address) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.SetCoordinator(&_VRFV2PlusConsumerExample.TransactOpts, _vrfCoordinator)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactorSession) SetCoordinator(_vrfCoordinator common.Address) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.SetCoordinator(&_VRFV2PlusConsumerExample.TransactOpts, _vrfCoordinator)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactor) SetSubId(opts *bind.TransactOpts, subId *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.contract.Transact(opts, "setSubId", subId)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleSession) SetSubId(subId *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.SetSubId(&_VRFV2PlusConsumerExample.TransactOpts, subId)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactorSession) SetSubId(subId *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.SetSubId(&_VRFV2PlusConsumerExample.TransactOpts, subId)
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

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactor) TopUpSubscriptionNative(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.contract.Transact(opts, "topUpSubscriptionNative")
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleSession) TopUpSubscriptionNative() (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.TopUpSubscriptionNative(&_VRFV2PlusConsumerExample.TransactOpts)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactorSession) TopUpSubscriptionNative() (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.TopUpSubscriptionNative(&_VRFV2PlusConsumerExample.TransactOpts)
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
	Fulfilled bool
	Requester common.Address
	RequestId *big.Int
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

	SSubId(opts *bind.CallOpts) (*big.Int, error)

	SVrfCoordinator(opts *bind.CallOpts) (common.Address, error)

	SVrfCoordinatorApiV1(opts *bind.CallOpts) (common.Address, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	CreateSubscriptionAndFund(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	CreateSubscriptionAndFundNative(opts *bind.TransactOpts) (*types.Transaction, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	RequestRandomWords(opts *bind.TransactOpts, callbackGasLimit uint32, requestConfirmations uint16, numWords uint32, keyHash [32]byte, nativePayment bool) (*types.Transaction, error)

	SetCoordinator(opts *bind.TransactOpts, _vrfCoordinator common.Address) (*types.Transaction, error)

	SetSubId(opts *bind.TransactOpts, subId *big.Int) (*types.Transaction, error)

	TopUpSubscription(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	TopUpSubscriptionNative(opts *bind.TransactOpts) (*types.Transaction, error)

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
