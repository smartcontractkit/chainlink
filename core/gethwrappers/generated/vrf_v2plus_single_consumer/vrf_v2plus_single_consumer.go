// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_v2plus_single_consumer

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

var VRFV2PlusSingleConsumerExampleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"nativePayment\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlyCoordinatorCanFulfill\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"}],\"name\":\"OnlyOwnerOrCoordinator\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"vrfCoordinator\",\"type\":\"address\"}],\"name\":\"CoordinatorSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"fundAndRequestRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"requestRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_randomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestConfig\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"nativePayment\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_vrfCoordinator\",\"outputs\":[{\"internalType\":\"contractIVRFCoordinatorV2Plus\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vrfCoordinator\",\"type\":\"address\"}],\"name\":\"setCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"subscribe\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"topUpSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"unsubscribe\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b50604051620017e7380380620017e783398101604081905262000034916200046f565b8633806000816200008c5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000bf57620000bf81620001d0565b5050506001600160a01b038116620000ea5760405163d92e233d60e01b815260040160405180910390fd5b600280546001600160a01b03199081166001600160a01b039384161790915560038054821692891692909217909155600a80543392169190911790556040805160c081018252600080825263ffffffff8881166020840181905261ffff8916948401859052908716606084018190526080840187905285151560a09094018490526004929092556005805465ffffffffffff19169091176401000000009094029390931763ffffffff60301b191666010000000000009091021790915560068390556007805460ff19169091179055620001c36200027b565b505050505050506200053b565b336001600160a01b038216036200022a5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000083565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b62000285620003df565b604080516001808252818301909252600091602080830190803683370190505090503081600081518110620002be57620002be6200050b565b6001600160a01b039283166020918202929092018101919091526002546040805163288688f960e21b81529051919093169263a21a23e492600480830193919282900301816000875af11580156200031a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000340919062000521565b600481905560025482516001600160a01b039091169163bec4c08c9184906000906200037057620003706200050b565b60200260200101516040518363ffffffff1660e01b8152600401620003a89291909182526001600160a01b0316602082015260400190565b600060405180830381600087803b158015620003c357600080fd5b505af1158015620003d8573d6000803e3d6000fd5b5050505050565b6000546001600160a01b031633146200043b5760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640162000083565b565b80516001600160a01b03811681146200045557600080fd5b919050565b805163ffffffff811681146200045557600080fd5b600080600080600080600060e0888a0312156200048b57600080fd5b62000496886200043d565b9650620004a6602089016200043d565b9550620004b6604089016200045a565b9450606088015161ffff81168114620004ce57600080fd5b9350620004de608089016200045a565b925060a0880151915060c08801518015158114620004fb57600080fd5b8091505092959891949750929550565b634e487b7160e01b600052603260045260246000fd5b6000602082840312156200053457600080fd5b5051919050565b61129c806200054b6000396000f3fe608060405234801561001057600080fd5b50600436106100f45760003560e01c80638da5cb5b11610097578063e0c8628911610066578063e0c862891461025c578063e89e106a14610264578063f2fde38b1461027b578063f6eaffc81461028e57600080fd5b80638da5cb5b146101e25780638ea98117146102215780638f449a05146102345780639eccacf61461023c57600080fd5b80637262561c116100d35780637262561c1461013457806379ba5097146101475780637db9263f1461014f57806386850e93146101cf57600080fd5b8062f714ce146100f95780631fe543e31461010e5780636fd700bb14610121575b600080fd5b61010c610107366004611038565b6102a1565b005b61010c61011c366004611064565b61034b565b61010c61012f3660046110e3565b6103ce565b61010c6101423660046110fc565b6105e6565b61010c610683565b60045460055460065460075461018b939263ffffffff8082169361ffff6401000000008404169366010000000000009093049091169160ff1686565b6040805196875263ffffffff958616602088015261ffff90941693860193909352921660608401526080830191909152151560a082015260c0015b60405180910390f35b61010c6101dd3660046110e3565b610780565b60005473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101c6565b61010c61022f3660046110fc565b61084b565b61010c6109d5565b6002546101fc9073ffffffffffffffffffffffffffffffffffffffff1681565b61010c610b6b565b61026d60095481565b6040519081526020016101c6565b61010c6102893660046110fc565b610cc9565b61026d61029c3660046110e3565b610cdd565b6102a9610cfe565b6003546040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8381166004830152602482018590529091169063a9059cbb906044016020604051808303816000875af1158015610322573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610346919061111e565b505050565b60025473ffffffffffffffffffffffffffffffffffffffff1633146103c3576002546040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff90911660248201526044015b60405180910390fd5b610346838383610d81565b6103d6610cfe565b6040805160c08101825260045480825260055463ffffffff808216602080860191909152640100000000830461ffff16858701526601000000000000909204166060840152600654608084015260075460ff16151560a0840152600354600254855192830193909352929373ffffffffffffffffffffffffffffffffffffffff93841693634000aea09316918691016040516020818303038152906040526040518463ffffffff1660e01b8152600401610492939291906111a4565b6020604051808303816000875af11580156104b1573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104d5919061111e565b5060006040518060c001604052808360800151815260200183600001518152602001836040015161ffff168152602001836020015163ffffffff168152602001836060015163ffffffff16815260200161054260405180602001604052808660a001511515815250610dfe565b90526002546040517f9b1c385e00000000000000000000000000000000000000000000000000000000815291925073ffffffffffffffffffffffffffffffffffffffff1690639b1c385e9061059b9084906004016111e2565b6020604051808303816000875af11580156105ba573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105de9190611247565b600955505050565b6105ee610cfe565b600254600480546040517f0ae095400000000000000000000000000000000000000000000000000000000081529182015273ffffffffffffffffffffffffffffffffffffffff838116602483015290911690630ae0954090604401600060405180830381600087803b15801561066357600080fd5b505af1158015610677573d6000803e3d6000fd5b50506000600455505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610704576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016103ba565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610788610cfe565b6003546002546004546040805160208082019390935281518082039093018352808201918290527f4000aea00000000000000000000000000000000000000000000000000000000090915273ffffffffffffffffffffffffffffffffffffffff93841693634000aea093610804939116918691906044016111a4565b6020604051808303816000875af1158015610823573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610847919061111e565b5050565b60005473ffffffffffffffffffffffffffffffffffffffff16331480159061088b575060025473ffffffffffffffffffffffffffffffffffffffff163314155b1561090f57336108b060005473ffffffffffffffffffffffffffffffffffffffff1690565b6002546040517f061db9c100000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff938416600482015291831660248301529190911660448201526064016103ba565b73ffffffffffffffffffffffffffffffffffffffff811661095c576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83169081179091556040519081527fd1a6a14209a385a964d036e404cb5cfb71f4000cdb03c9366292430787261be69060200160405180910390a150565b6109dd610cfe565b604080516001808252818301909252600091602080830190803683370190505090503081600081518110610a1357610a13611260565b73ffffffffffffffffffffffffffffffffffffffff928316602091820292909201810191909152600254604080517fa21a23e40000000000000000000000000000000000000000000000000000000081529051919093169263a21a23e492600480830193919282900301816000875af1158015610a94573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ab89190611247565b6004819055600254825173ffffffffffffffffffffffffffffffffffffffff9091169163bec4c08c918490600090610af257610af2611260565b60200260200101516040518363ffffffff1660e01b8152600401610b3692919091825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b600060405180830381600087803b158015610b5057600080fd5b505af1158015610b64573d6000803e3d6000fd5b5050505050565b610b73610cfe565b6040805160c08082018352600454825260055463ffffffff808216602080860191825261ffff640100000000850481168789019081526601000000000000909504841660608089019182526006546080808b0191825260075460ff16151560a0808d019182528d519b8c018e5292518b528b518b8801529851909416898c0152945186169088015251909316928501929092528551918201909552905115158152919260009290820190610c2690610dfe565b90526002546040517f9b1c385e00000000000000000000000000000000000000000000000000000000815291925073ffffffffffffffffffffffffffffffffffffffff1690639b1c385e90610c7f9084906004016111e2565b6020604051808303816000875af1158015610c9e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610cc29190611247565b6009555050565b610cd1610cfe565b610cda81610eba565b50565b60088181548110610ced57600080fd5b600091825260209091200154905081565b60005473ffffffffffffffffffffffffffffffffffffffff163314610d7f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016103ba565b565b6009548314610dec576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265717565737420494420697320696e636f727265637400000000000000000060448201526064016103ba565b610df860088383610faf565b50505050565b60607f92fd13387c7fe7befbc38d303d6468778fb9731bc4583f17d92989c6fcfdeaaa82604051602401610e3791511515815260200190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009093169290921790915292915050565b3373ffffffffffffffffffffffffffffffffffffffff821603610f39576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016103ba565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b828054828255906000526020600020908101928215610fea579160200282015b82811115610fea578235825591602001919060010190610fcf565b50610ff6929150610ffa565b5090565b5b80821115610ff65760008155600101610ffb565b803573ffffffffffffffffffffffffffffffffffffffff8116811461103357600080fd5b919050565b6000806040838503121561104b57600080fd5b8235915061105b6020840161100f565b90509250929050565b60008060006040848603121561107957600080fd5b83359250602084013567ffffffffffffffff8082111561109857600080fd5b818601915086601f8301126110ac57600080fd5b8135818111156110bb57600080fd5b8760208260051b85010111156110d057600080fd5b6020830194508093505050509250925092565b6000602082840312156110f557600080fd5b5035919050565b60006020828403121561110e57600080fd5b6111178261100f565b9392505050565b60006020828403121561113057600080fd5b8151801515811461111757600080fd5b6000815180845260005b818110156111665760208185018101518683018201520161114a565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b73ffffffffffffffffffffffffffffffffffffffff841681528260208201526060604082015260006111d96060830184611140565b95945050505050565b60208152815160208201526020820151604082015261ffff60408301511660608201526000606083015163ffffffff80821660808501528060808601511660a0850152505060a083015160c08084015261123f60e0840182611140565b949350505050565b60006020828403121561125957600080fd5b5051919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fdfea164736f6c6343000813000a",
}

var VRFV2PlusSingleConsumerExampleABI = VRFV2PlusSingleConsumerExampleMetaData.ABI

var VRFV2PlusSingleConsumerExampleBin = VRFV2PlusSingleConsumerExampleMetaData.Bin

func DeployVRFV2PlusSingleConsumerExample(auth *bind.TransactOpts, backend bind.ContractBackend, vrfCoordinator common.Address, link common.Address, callbackGasLimit uint32, requestConfirmations uint16, numWords uint32, keyHash [32]byte, nativePayment bool) (common.Address, *types.Transaction, *VRFV2PlusSingleConsumerExample, error) {
	parsed, err := VRFV2PlusSingleConsumerExampleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFV2PlusSingleConsumerExampleBin), backend, vrfCoordinator, link, callbackGasLimit, requestConfirmations, numWords, keyHash, nativePayment)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFV2PlusSingleConsumerExample{address: address, abi: *parsed, VRFV2PlusSingleConsumerExampleCaller: VRFV2PlusSingleConsumerExampleCaller{contract: contract}, VRFV2PlusSingleConsumerExampleTransactor: VRFV2PlusSingleConsumerExampleTransactor{contract: contract}, VRFV2PlusSingleConsumerExampleFilterer: VRFV2PlusSingleConsumerExampleFilterer{contract: contract}}, nil
}

type VRFV2PlusSingleConsumerExample struct {
	address common.Address
	abi     abi.ABI
	VRFV2PlusSingleConsumerExampleCaller
	VRFV2PlusSingleConsumerExampleTransactor
	VRFV2PlusSingleConsumerExampleFilterer
}

type VRFV2PlusSingleConsumerExampleCaller struct {
	contract *bind.BoundContract
}

type VRFV2PlusSingleConsumerExampleTransactor struct {
	contract *bind.BoundContract
}

type VRFV2PlusSingleConsumerExampleFilterer struct {
	contract *bind.BoundContract
}

type VRFV2PlusSingleConsumerExampleSession struct {
	Contract     *VRFV2PlusSingleConsumerExample
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFV2PlusSingleConsumerExampleCallerSession struct {
	Contract *VRFV2PlusSingleConsumerExampleCaller
	CallOpts bind.CallOpts
}

type VRFV2PlusSingleConsumerExampleTransactorSession struct {
	Contract     *VRFV2PlusSingleConsumerExampleTransactor
	TransactOpts bind.TransactOpts
}

type VRFV2PlusSingleConsumerExampleRaw struct {
	Contract *VRFV2PlusSingleConsumerExample
}

type VRFV2PlusSingleConsumerExampleCallerRaw struct {
	Contract *VRFV2PlusSingleConsumerExampleCaller
}

type VRFV2PlusSingleConsumerExampleTransactorRaw struct {
	Contract *VRFV2PlusSingleConsumerExampleTransactor
}

func NewVRFV2PlusSingleConsumerExample(address common.Address, backend bind.ContractBackend) (*VRFV2PlusSingleConsumerExample, error) {
	abi, err := abi.JSON(strings.NewReader(VRFV2PlusSingleConsumerExampleABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFV2PlusSingleConsumerExample(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusSingleConsumerExample{address: address, abi: abi, VRFV2PlusSingleConsumerExampleCaller: VRFV2PlusSingleConsumerExampleCaller{contract: contract}, VRFV2PlusSingleConsumerExampleTransactor: VRFV2PlusSingleConsumerExampleTransactor{contract: contract}, VRFV2PlusSingleConsumerExampleFilterer: VRFV2PlusSingleConsumerExampleFilterer{contract: contract}}, nil
}

func NewVRFV2PlusSingleConsumerExampleCaller(address common.Address, caller bind.ContractCaller) (*VRFV2PlusSingleConsumerExampleCaller, error) {
	contract, err := bindVRFV2PlusSingleConsumerExample(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusSingleConsumerExampleCaller{contract: contract}, nil
}

func NewVRFV2PlusSingleConsumerExampleTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFV2PlusSingleConsumerExampleTransactor, error) {
	contract, err := bindVRFV2PlusSingleConsumerExample(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusSingleConsumerExampleTransactor{contract: contract}, nil
}

func NewVRFV2PlusSingleConsumerExampleFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFV2PlusSingleConsumerExampleFilterer, error) {
	contract, err := bindVRFV2PlusSingleConsumerExample(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusSingleConsumerExampleFilterer{contract: contract}, nil
}

func bindVRFV2PlusSingleConsumerExample(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFV2PlusSingleConsumerExampleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2PlusSingleConsumerExample.Contract.VRFV2PlusSingleConsumerExampleCaller.contract.Call(opts, result, method, params...)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.VRFV2PlusSingleConsumerExampleTransactor.contract.Transfer(opts)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.VRFV2PlusSingleConsumerExampleTransactor.contract.Transact(opts, method, params...)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2PlusSingleConsumerExample.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.contract.Transfer(opts)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.contract.Transact(opts, method, params...)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2PlusSingleConsumerExample.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleSession) Owner() (common.Address, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.Owner(&_VRFV2PlusSingleConsumerExample.CallOpts)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleCallerSession) Owner() (common.Address, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.Owner(&_VRFV2PlusSingleConsumerExample.CallOpts)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleCaller) SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusSingleConsumerExample.contract.Call(opts, &out, "s_randomWords", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.SRandomWords(&_VRFV2PlusSingleConsumerExample.CallOpts, arg0)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleCallerSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.SRandomWords(&_VRFV2PlusSingleConsumerExample.CallOpts, arg0)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleCaller) SRequestConfig(opts *bind.CallOpts) (SRequestConfig,

	error) {
	var out []interface{}
	err := _VRFV2PlusSingleConsumerExample.contract.Call(opts, &out, "s_requestConfig")

	outstruct := new(SRequestConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.SubId = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.CallbackGasLimit = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.RequestConfirmations = *abi.ConvertType(out[2], new(uint16)).(*uint16)
	outstruct.NumWords = *abi.ConvertType(out[3], new(uint32)).(*uint32)
	outstruct.KeyHash = *abi.ConvertType(out[4], new([32]byte)).(*[32]byte)
	outstruct.NativePayment = *abi.ConvertType(out[5], new(bool)).(*bool)

	return *outstruct, err

}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleSession) SRequestConfig() (SRequestConfig,

	error) {
	return _VRFV2PlusSingleConsumerExample.Contract.SRequestConfig(&_VRFV2PlusSingleConsumerExample.CallOpts)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleCallerSession) SRequestConfig() (SRequestConfig,

	error) {
	return _VRFV2PlusSingleConsumerExample.Contract.SRequestConfig(&_VRFV2PlusSingleConsumerExample.CallOpts)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleCaller) SRequestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusSingleConsumerExample.contract.Call(opts, &out, "s_requestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleSession) SRequestId() (*big.Int, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.SRequestId(&_VRFV2PlusSingleConsumerExample.CallOpts)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleCallerSession) SRequestId() (*big.Int, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.SRequestId(&_VRFV2PlusSingleConsumerExample.CallOpts)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleCaller) SVrfCoordinator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2PlusSingleConsumerExample.contract.Call(opts, &out, "s_vrfCoordinator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleSession) SVrfCoordinator() (common.Address, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.SVrfCoordinator(&_VRFV2PlusSingleConsumerExample.CallOpts)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleCallerSession) SVrfCoordinator() (common.Address, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.SVrfCoordinator(&_VRFV2PlusSingleConsumerExample.CallOpts)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.contract.Transact(opts, "acceptOwnership")
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.AcceptOwnership(&_VRFV2PlusSingleConsumerExample.TransactOpts)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.AcceptOwnership(&_VRFV2PlusSingleConsumerExample.TransactOpts)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactor) FundAndRequestRandomWords(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.contract.Transact(opts, "fundAndRequestRandomWords", amount)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleSession) FundAndRequestRandomWords(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.FundAndRequestRandomWords(&_VRFV2PlusSingleConsumerExample.TransactOpts, amount)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactorSession) FundAndRequestRandomWords(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.FundAndRequestRandomWords(&_VRFV2PlusSingleConsumerExample.TransactOpts, amount)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.contract.Transact(opts, "rawFulfillRandomWords", requestId, randomWords)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.RawFulfillRandomWords(&_VRFV2PlusSingleConsumerExample.TransactOpts, requestId, randomWords)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactorSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.RawFulfillRandomWords(&_VRFV2PlusSingleConsumerExample.TransactOpts, requestId, randomWords)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactor) RequestRandomWords(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.contract.Transact(opts, "requestRandomWords")
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleSession) RequestRandomWords() (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.RequestRandomWords(&_VRFV2PlusSingleConsumerExample.TransactOpts)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactorSession) RequestRandomWords() (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.RequestRandomWords(&_VRFV2PlusSingleConsumerExample.TransactOpts)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactor) SetCoordinator(opts *bind.TransactOpts, _vrfCoordinator common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.contract.Transact(opts, "setCoordinator", _vrfCoordinator)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleSession) SetCoordinator(_vrfCoordinator common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.SetCoordinator(&_VRFV2PlusSingleConsumerExample.TransactOpts, _vrfCoordinator)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactorSession) SetCoordinator(_vrfCoordinator common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.SetCoordinator(&_VRFV2PlusSingleConsumerExample.TransactOpts, _vrfCoordinator)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactor) Subscribe(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.contract.Transact(opts, "subscribe")
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleSession) Subscribe() (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.Subscribe(&_VRFV2PlusSingleConsumerExample.TransactOpts)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactorSession) Subscribe() (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.Subscribe(&_VRFV2PlusSingleConsumerExample.TransactOpts)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactor) TopUpSubscription(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.contract.Transact(opts, "topUpSubscription", amount)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleSession) TopUpSubscription(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.TopUpSubscription(&_VRFV2PlusSingleConsumerExample.TransactOpts, amount)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactorSession) TopUpSubscription(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.TopUpSubscription(&_VRFV2PlusSingleConsumerExample.TransactOpts, amount)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.contract.Transact(opts, "transferOwnership", to)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.TransferOwnership(&_VRFV2PlusSingleConsumerExample.TransactOpts, to)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.TransferOwnership(&_VRFV2PlusSingleConsumerExample.TransactOpts, to)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactor) Unsubscribe(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.contract.Transact(opts, "unsubscribe", to)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleSession) Unsubscribe(to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.Unsubscribe(&_VRFV2PlusSingleConsumerExample.TransactOpts, to)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactorSession) Unsubscribe(to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.Unsubscribe(&_VRFV2PlusSingleConsumerExample.TransactOpts, to)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactor) Withdraw(opts *bind.TransactOpts, amount *big.Int, to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.contract.Transact(opts, "withdraw", amount, to)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleSession) Withdraw(amount *big.Int, to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.Withdraw(&_VRFV2PlusSingleConsumerExample.TransactOpts, amount, to)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactorSession) Withdraw(amount *big.Int, to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.Withdraw(&_VRFV2PlusSingleConsumerExample.TransactOpts, amount, to)
}

type VRFV2PlusSingleConsumerExampleCoordinatorSetIterator struct {
	Event *VRFV2PlusSingleConsumerExampleCoordinatorSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2PlusSingleConsumerExampleCoordinatorSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2PlusSingleConsumerExampleCoordinatorSet)
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
		it.Event = new(VRFV2PlusSingleConsumerExampleCoordinatorSet)
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

func (it *VRFV2PlusSingleConsumerExampleCoordinatorSetIterator) Error() error {
	return it.fail
}

func (it *VRFV2PlusSingleConsumerExampleCoordinatorSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2PlusSingleConsumerExampleCoordinatorSet struct {
	VrfCoordinator common.Address
	Raw            types.Log
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleFilterer) FilterCoordinatorSet(opts *bind.FilterOpts) (*VRFV2PlusSingleConsumerExampleCoordinatorSetIterator, error) {

	logs, sub, err := _VRFV2PlusSingleConsumerExample.contract.FilterLogs(opts, "CoordinatorSet")
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusSingleConsumerExampleCoordinatorSetIterator{contract: _VRFV2PlusSingleConsumerExample.contract, event: "CoordinatorSet", logs: logs, sub: sub}, nil
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleFilterer) WatchCoordinatorSet(opts *bind.WatchOpts, sink chan<- *VRFV2PlusSingleConsumerExampleCoordinatorSet) (event.Subscription, error) {

	logs, sub, err := _VRFV2PlusSingleConsumerExample.contract.WatchLogs(opts, "CoordinatorSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2PlusSingleConsumerExampleCoordinatorSet)
				if err := _VRFV2PlusSingleConsumerExample.contract.UnpackLog(event, "CoordinatorSet", log); err != nil {
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

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleFilterer) ParseCoordinatorSet(log types.Log) (*VRFV2PlusSingleConsumerExampleCoordinatorSet, error) {
	event := new(VRFV2PlusSingleConsumerExampleCoordinatorSet)
	if err := _VRFV2PlusSingleConsumerExample.contract.UnpackLog(event, "CoordinatorSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFV2PlusSingleConsumerExampleOwnershipTransferRequestedIterator struct {
	Event *VRFV2PlusSingleConsumerExampleOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2PlusSingleConsumerExampleOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2PlusSingleConsumerExampleOwnershipTransferRequested)
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
		it.Event = new(VRFV2PlusSingleConsumerExampleOwnershipTransferRequested)
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

func (it *VRFV2PlusSingleConsumerExampleOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFV2PlusSingleConsumerExampleOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2PlusSingleConsumerExampleOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusSingleConsumerExampleOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusSingleConsumerExample.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusSingleConsumerExampleOwnershipTransferRequestedIterator{contract: _VRFV2PlusSingleConsumerExample.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFV2PlusSingleConsumerExampleOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusSingleConsumerExample.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2PlusSingleConsumerExampleOwnershipTransferRequested)
				if err := _VRFV2PlusSingleConsumerExample.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleFilterer) ParseOwnershipTransferRequested(log types.Log) (*VRFV2PlusSingleConsumerExampleOwnershipTransferRequested, error) {
	event := new(VRFV2PlusSingleConsumerExampleOwnershipTransferRequested)
	if err := _VRFV2PlusSingleConsumerExample.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFV2PlusSingleConsumerExampleOwnershipTransferredIterator struct {
	Event *VRFV2PlusSingleConsumerExampleOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2PlusSingleConsumerExampleOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2PlusSingleConsumerExampleOwnershipTransferred)
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
		it.Event = new(VRFV2PlusSingleConsumerExampleOwnershipTransferred)
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

func (it *VRFV2PlusSingleConsumerExampleOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFV2PlusSingleConsumerExampleOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2PlusSingleConsumerExampleOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusSingleConsumerExampleOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusSingleConsumerExample.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusSingleConsumerExampleOwnershipTransferredIterator{contract: _VRFV2PlusSingleConsumerExample.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFV2PlusSingleConsumerExampleOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusSingleConsumerExample.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2PlusSingleConsumerExampleOwnershipTransferred)
				if err := _VRFV2PlusSingleConsumerExample.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleFilterer) ParseOwnershipTransferred(log types.Log) (*VRFV2PlusSingleConsumerExampleOwnershipTransferred, error) {
	event := new(VRFV2PlusSingleConsumerExampleOwnershipTransferred)
	if err := _VRFV2PlusSingleConsumerExample.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type SRequestConfig struct {
	SubId                *big.Int
	CallbackGasLimit     uint32
	RequestConfirmations uint16
	NumWords             uint32
	KeyHash              [32]byte
	NativePayment        bool
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExample) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFV2PlusSingleConsumerExample.abi.Events["CoordinatorSet"].ID:
		return _VRFV2PlusSingleConsumerExample.ParseCoordinatorSet(log)
	case _VRFV2PlusSingleConsumerExample.abi.Events["OwnershipTransferRequested"].ID:
		return _VRFV2PlusSingleConsumerExample.ParseOwnershipTransferRequested(log)
	case _VRFV2PlusSingleConsumerExample.abi.Events["OwnershipTransferred"].ID:
		return _VRFV2PlusSingleConsumerExample.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFV2PlusSingleConsumerExampleCoordinatorSet) Topic() common.Hash {
	return common.HexToHash("0xd1a6a14209a385a964d036e404cb5cfb71f4000cdb03c9366292430787261be6")
}

func (VRFV2PlusSingleConsumerExampleOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VRFV2PlusSingleConsumerExampleOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExample) Address() common.Address {
	return _VRFV2PlusSingleConsumerExample.address
}

type VRFV2PlusSingleConsumerExampleInterface interface {
	Owner(opts *bind.CallOpts) (common.Address, error)

	SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	SRequestConfig(opts *bind.CallOpts) (SRequestConfig,

		error)

	SRequestId(opts *bind.CallOpts) (*big.Int, error)

	SVrfCoordinator(opts *bind.CallOpts) (common.Address, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	FundAndRequestRandomWords(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	RequestRandomWords(opts *bind.TransactOpts) (*types.Transaction, error)

	SetCoordinator(opts *bind.TransactOpts, _vrfCoordinator common.Address) (*types.Transaction, error)

	Subscribe(opts *bind.TransactOpts) (*types.Transaction, error)

	TopUpSubscription(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Unsubscribe(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Withdraw(opts *bind.TransactOpts, amount *big.Int, to common.Address) (*types.Transaction, error)

	FilterCoordinatorSet(opts *bind.FilterOpts) (*VRFV2PlusSingleConsumerExampleCoordinatorSetIterator, error)

	WatchCoordinatorSet(opts *bind.WatchOpts, sink chan<- *VRFV2PlusSingleConsumerExampleCoordinatorSet) (event.Subscription, error)

	ParseCoordinatorSet(log types.Log) (*VRFV2PlusSingleConsumerExampleCoordinatorSet, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusSingleConsumerExampleOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFV2PlusSingleConsumerExampleOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VRFV2PlusSingleConsumerExampleOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusSingleConsumerExampleOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFV2PlusSingleConsumerExampleOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VRFV2PlusSingleConsumerExampleOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
