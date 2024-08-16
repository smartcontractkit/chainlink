// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package self_funded_ping_pong

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

type ClientAny2EVMMessage struct {
	MessageId           [32]byte
	SourceChainSelector uint64
	Sender              []byte
	Data                []byte
	DestTokenAmounts    []ClientEVMTokenAmount
}

type ClientEVMTokenAmount struct {
	Token  common.Address
	Amount *big.Int
}

var SelfFundedPingPongMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"contractIERC20\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"roundTripsBeforeFunding\",\"type\":\"uint8\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"name\":\"InvalidRouter\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"countIncrBeforeFunding\",\"type\":\"uint8\"}],\"name\":\"CountIncrBeforeFundingSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Funded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"pingPongCount\",\"type\":\"uint256\"}],\"name\":\"Ping\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"pingPongCount\",\"type\":\"uint256\"}],\"name\":\"Pong\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMTokenAmount[]\",\"name\":\"destTokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structClient.Any2EVMMessage\",\"name\":\"message\",\"type\":\"tuple\"}],\"name\":\"ccipReceive\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"pingPongCount\",\"type\":\"uint256\"}],\"name\":\"fundPingPong\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCountIncrBeforeFunding\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCounterpartAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCounterpartChainSelector\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFeeToken\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRouter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isPaused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"countIncrBeforeFunding\",\"type\":\"uint8\"}],\"name\":\"setCountIncrBeforeFunding\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"counterpartChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"counterpartAddress\",\"type\":\"address\"}],\"name\":\"setCounterpart\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"setCounterpartAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"}],\"name\":\"setCounterpartChainSelector\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"pause\",\"type\":\"bool\"}],\"name\":\"setPaused\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"startPingPong\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b506040516200182238038062001822833981016040819052620000349162000291565b828233806000846001600160a01b0381166200006b576040516335fdcccd60e21b8152600060048201526024015b60405180910390fd5b6001600160a01b039081166080528216620000c95760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f0000000000000000604482015260640162000062565b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000fc57620000fc81620001cd565b50506002805460ff60a01b1916905550600380546001600160a01b0319166001600160a01b0383811691821790925560405163095ea7b360e01b8152918416600483015260001960248301529063095ea7b3906044016020604051808303816000875af115801562000172573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001989190620002ea565b505050806002620001aa919062000315565b600360146101000a81548160ff021916908360ff16021790555050505062000347565b336001600160a01b03821603620002275760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000062565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6001600160a01b03811681146200028e57600080fd5b50565b600080600060608486031215620002a757600080fd5b8351620002b48162000278565b6020850151909350620002c78162000278565b604085015190925060ff81168114620002df57600080fd5b809150509250925092565b600060208284031215620002fd57600080fd5b815180151581146200030e57600080fd5b9392505050565b60ff81811683821602908116908181146200034057634e487b7160e01b600052601160045260246000fd5b5092915050565b6080516114aa62000378600039600081816102970152818161063f0152818161072c0152610c2201526114aa6000f3fe608060405234801561001057600080fd5b50600436106101515760003560e01c80638f491cba116100cd578063bee518a411610081578063e6c725f511610066578063e6c725f51461034d578063ef686d8e1461037d578063f2fde38b1461039057600080fd5b8063bee518a4146102f1578063ca709a251461032f57600080fd5b8063b0f479a1116100b2578063b0f479a114610295578063b187bd26146102bb578063b5a11011146102de57600080fd5b80638f491cba1461026f5780639d2aede51461028257600080fd5b80632874d8bf1161012457806379ba50971161010957806379ba50971461023657806385572ffb1461023e5780638da5cb5b1461025157600080fd5b80632874d8bf146101ef5780632b6e5d63146101f757600080fd5b806301ffc9a71461015657806316c38b3c1461017e578063181f5a77146101935780631892b906146101dc575b600080fd5b610169610164366004610e24565b6103a3565b60405190151581526020015b60405180910390f35b61019161018c366004610e6d565b61043c565b005b6101cf6040518060400160405280601881526020017f53656c6646756e64656450696e67506f6e6720312e322e30000000000000000081525081565b6040516101759190610ef3565b6101916101ea366004610f23565b61048e565b6101916104e9565b60025473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610175565b610191610525565b61019161024c366004610f3e565b610627565b60005473ffffffffffffffffffffffffffffffffffffffff16610211565b61019161027d366004610f79565b6106ac565b610191610290366004610fb4565b61088b565b7f0000000000000000000000000000000000000000000000000000000000000000610211565b60025474010000000000000000000000000000000000000000900460ff16610169565b6101916102ec366004610fd1565b6108da565b60015474010000000000000000000000000000000000000000900467ffffffffffffffff1660405167ffffffffffffffff9091168152602001610175565b60035473ffffffffffffffffffffffffffffffffffffffff16610211565b60035474010000000000000000000000000000000000000000900460ff1660405160ff9091168152602001610175565b61019161038b366004611008565b61097c565b61019161039e366004610fb4565b610a04565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167f85572ffb00000000000000000000000000000000000000000000000000000000148061043657507fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a700000000000000000000000000000000000000000000000000000000145b92915050565b610444610a15565b6002805491151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff909216919091179055565b610496610a15565b6001805467ffffffffffffffff90921674010000000000000000000000000000000000000000027fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff909216919091179055565b6104f1610a15565b600280547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff1690556105236001610a96565b565b60015473ffffffffffffffffffffffffffffffffffffffff1633146105ab576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610698576040517fd7f733340000000000000000000000000000000000000000000000000000000081523360048201526024016105a2565b6106a96106a482611230565b610cd9565b50565b60035474010000000000000000000000000000000000000000900460ff1615806106f2575060035474010000000000000000000000000000000000000000900460ff1681105b156106fa5750565b6003546001906107259074010000000000000000000000000000000000000000900460ff16836112dd565b116106a9577f00000000000000000000000000000000000000000000000000000000000000006001546040517fa8d87a3b0000000000000000000000000000000000000000000000000000000081527401000000000000000000000000000000000000000090910467ffffffffffffffff16600482015273ffffffffffffffffffffffffffffffffffffffff919091169063a8d87a3b90602401602060405180830381865afa1580156107dc573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108009190611318565b73ffffffffffffffffffffffffffffffffffffffff1663eff7cc486040518163ffffffff1660e01b8152600401600060405180830381600087803b15801561084757600080fd5b505af115801561085b573d6000803e3d6000fd5b50506040517f302777af5d26fab9dd5120c5f1307c65193ebc51daf33244ada4365fab10602c925060009150a150565b610893610a15565b600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b6108e2610a15565b6001805467ffffffffffffffff90931674010000000000000000000000000000000000000000027fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff909316929092179091556002805473ffffffffffffffffffffffffffffffffffffffff9092167fffffffffffffffffffffffff0000000000000000000000000000000000000000909216919091179055565b610984610a15565b600380547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000060ff8416908102919091179091556040519081527f4768dbf8645b24c54f2887651545d24f748c0d0d1d4c689eb810fb19f0befcf39060200160405180910390a150565b610a0c610a15565b6106a981610d2f565b60005473ffffffffffffffffffffffffffffffffffffffff163314610523576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016105a2565b80600116600103610ad9576040518181527f48257dc961b6f792c2b78a080dacfed693b660960a702de21cee364e20270e2f9060200160405180910390a1610b0d565b6040518181527f58b69f57828e6962d216502094c54f6562f3bf082ba758966c3454f9e37b15259060200160405180910390a15b610b16816106ac565b6040805160a0810190915260025473ffffffffffffffffffffffffffffffffffffffff1660c08201526000908060e08101604051602081830303815290604052815260200183604051602001610b6e91815260200190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905281526020016000604051908082528060200260200182016040528015610be857816020015b6040805180820190915260008082526020820152815260200190600190039081610bc15790505b50815260035473ffffffffffffffffffffffffffffffffffffffff16602080830191909152604080519182018152600082529091015290507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166396f4e9f9600160149054906101000a900467ffffffffffffffff16836040518363ffffffff1660e01b8152600401610c91929190611335565b6020604051808303816000875af1158015610cb0573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610cd4919061144a565b505050565b60008160600151806020019051810190610cf3919061144a565b60025490915074010000000000000000000000000000000000000000900460ff16610d2b57610d2b610d26826001611463565b610a96565b5050565b3373ffffffffffffffffffffffffffffffffffffffff821603610dae576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016105a2565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600060208284031215610e3657600080fd5b81357fffffffff0000000000000000000000000000000000000000000000000000000081168114610e6657600080fd5b9392505050565b600060208284031215610e7f57600080fd5b81358015158114610e6657600080fd5b6000815180845260005b81811015610eb557602081850181015186830182015201610e99565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b602081526000610e666020830184610e8f565b803567ffffffffffffffff81168114610f1e57600080fd5b919050565b600060208284031215610f3557600080fd5b610e6682610f06565b600060208284031215610f5057600080fd5b813567ffffffffffffffff811115610f6757600080fd5b820160a08185031215610e6657600080fd5b600060208284031215610f8b57600080fd5b5035919050565b73ffffffffffffffffffffffffffffffffffffffff811681146106a957600080fd5b600060208284031215610fc657600080fd5b8135610e6681610f92565b60008060408385031215610fe457600080fd5b610fed83610f06565b91506020830135610ffd81610f92565b809150509250929050565b60006020828403121561101a57600080fd5b813560ff81168114610e6657600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040805190810167ffffffffffffffff8111828210171561107d5761107d61102b565b60405290565b60405160a0810167ffffffffffffffff8111828210171561107d5761107d61102b565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156110ed576110ed61102b565b604052919050565b600082601f83011261110657600080fd5b813567ffffffffffffffff8111156111205761112061102b565b61115160207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116016110a6565b81815284602083860101111561116657600080fd5b816020850160208301376000918101602001919091529392505050565b600082601f83011261119457600080fd5b8135602067ffffffffffffffff8211156111b0576111b061102b565b6111be818360051b016110a6565b82815260069290921b840181019181810190868411156111dd57600080fd5b8286015b8481101561122557604081890312156111fa5760008081fd5b61120261105a565b813561120d81610f92565b815281850135858201528352918301916040016111e1565b509695505050505050565b600060a0823603121561124257600080fd5b61124a611083565b8235815261125a60208401610f06565b6020820152604083013567ffffffffffffffff8082111561127a57600080fd5b611286368387016110f5565b6040840152606085013591508082111561129f57600080fd5b6112ab368387016110f5565b606084015260808501359150808211156112c457600080fd5b506112d136828601611183565b60808301525092915050565b600082611313577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500690565b60006020828403121561132a57600080fd5b8151610e6681610f92565b6000604067ffffffffffffffff851683526020604081850152845160a0604086015261136460e0860182610e8f565b9050818601517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc08087840301606088015261139f8383610e8f565b6040890151888203830160808a01528051808352908601945060009350908501905b80841015611400578451805173ffffffffffffffffffffffffffffffffffffffff168352860151868301529385019360019390930192908601906113c1565b50606089015173ffffffffffffffffffffffffffffffffffffffff1660a08901526080890151888203830160c08a0152955061143c8187610e8f565b9a9950505050505050505050565b60006020828403121561145c57600080fd5b5051919050565b80820180821115610436577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fdfea164736f6c6343000818000a",
}

var SelfFundedPingPongABI = SelfFundedPingPongMetaData.ABI

var SelfFundedPingPongBin = SelfFundedPingPongMetaData.Bin

func DeploySelfFundedPingPong(auth *bind.TransactOpts, backend bind.ContractBackend, router common.Address, feeToken common.Address, roundTripsBeforeFunding uint8) (common.Address, *types.Transaction, *SelfFundedPingPong, error) {
	parsed, err := SelfFundedPingPongMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SelfFundedPingPongBin), backend, router, feeToken, roundTripsBeforeFunding)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SelfFundedPingPong{address: address, abi: *parsed, SelfFundedPingPongCaller: SelfFundedPingPongCaller{contract: contract}, SelfFundedPingPongTransactor: SelfFundedPingPongTransactor{contract: contract}, SelfFundedPingPongFilterer: SelfFundedPingPongFilterer{contract: contract}}, nil
}

type SelfFundedPingPong struct {
	address common.Address
	abi     abi.ABI
	SelfFundedPingPongCaller
	SelfFundedPingPongTransactor
	SelfFundedPingPongFilterer
}

type SelfFundedPingPongCaller struct {
	contract *bind.BoundContract
}

type SelfFundedPingPongTransactor struct {
	contract *bind.BoundContract
}

type SelfFundedPingPongFilterer struct {
	contract *bind.BoundContract
}

type SelfFundedPingPongSession struct {
	Contract     *SelfFundedPingPong
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type SelfFundedPingPongCallerSession struct {
	Contract *SelfFundedPingPongCaller
	CallOpts bind.CallOpts
}

type SelfFundedPingPongTransactorSession struct {
	Contract     *SelfFundedPingPongTransactor
	TransactOpts bind.TransactOpts
}

type SelfFundedPingPongRaw struct {
	Contract *SelfFundedPingPong
}

type SelfFundedPingPongCallerRaw struct {
	Contract *SelfFundedPingPongCaller
}

type SelfFundedPingPongTransactorRaw struct {
	Contract *SelfFundedPingPongTransactor
}

func NewSelfFundedPingPong(address common.Address, backend bind.ContractBackend) (*SelfFundedPingPong, error) {
	abi, err := abi.JSON(strings.NewReader(SelfFundedPingPongABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindSelfFundedPingPong(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SelfFundedPingPong{address: address, abi: abi, SelfFundedPingPongCaller: SelfFundedPingPongCaller{contract: contract}, SelfFundedPingPongTransactor: SelfFundedPingPongTransactor{contract: contract}, SelfFundedPingPongFilterer: SelfFundedPingPongFilterer{contract: contract}}, nil
}

func NewSelfFundedPingPongCaller(address common.Address, caller bind.ContractCaller) (*SelfFundedPingPongCaller, error) {
	contract, err := bindSelfFundedPingPong(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SelfFundedPingPongCaller{contract: contract}, nil
}

func NewSelfFundedPingPongTransactor(address common.Address, transactor bind.ContractTransactor) (*SelfFundedPingPongTransactor, error) {
	contract, err := bindSelfFundedPingPong(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SelfFundedPingPongTransactor{contract: contract}, nil
}

func NewSelfFundedPingPongFilterer(address common.Address, filterer bind.ContractFilterer) (*SelfFundedPingPongFilterer, error) {
	contract, err := bindSelfFundedPingPong(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SelfFundedPingPongFilterer{contract: contract}, nil
}

func bindSelfFundedPingPong(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SelfFundedPingPongMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_SelfFundedPingPong *SelfFundedPingPongRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SelfFundedPingPong.Contract.SelfFundedPingPongCaller.contract.Call(opts, result, method, params...)
}

func (_SelfFundedPingPong *SelfFundedPingPongRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SelfFundedPingPong.Contract.SelfFundedPingPongTransactor.contract.Transfer(opts)
}

func (_SelfFundedPingPong *SelfFundedPingPongRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SelfFundedPingPong.Contract.SelfFundedPingPongTransactor.contract.Transact(opts, method, params...)
}

func (_SelfFundedPingPong *SelfFundedPingPongCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SelfFundedPingPong.Contract.contract.Call(opts, result, method, params...)
}

func (_SelfFundedPingPong *SelfFundedPingPongTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SelfFundedPingPong.Contract.contract.Transfer(opts)
}

func (_SelfFundedPingPong *SelfFundedPingPongTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SelfFundedPingPong.Contract.contract.Transact(opts, method, params...)
}

func (_SelfFundedPingPong *SelfFundedPingPongCaller) GetCountIncrBeforeFunding(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _SelfFundedPingPong.contract.Call(opts, &out, "getCountIncrBeforeFunding")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_SelfFundedPingPong *SelfFundedPingPongSession) GetCountIncrBeforeFunding() (uint8, error) {
	return _SelfFundedPingPong.Contract.GetCountIncrBeforeFunding(&_SelfFundedPingPong.CallOpts)
}

func (_SelfFundedPingPong *SelfFundedPingPongCallerSession) GetCountIncrBeforeFunding() (uint8, error) {
	return _SelfFundedPingPong.Contract.GetCountIncrBeforeFunding(&_SelfFundedPingPong.CallOpts)
}

func (_SelfFundedPingPong *SelfFundedPingPongCaller) GetCounterpartAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SelfFundedPingPong.contract.Call(opts, &out, "getCounterpartAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_SelfFundedPingPong *SelfFundedPingPongSession) GetCounterpartAddress() (common.Address, error) {
	return _SelfFundedPingPong.Contract.GetCounterpartAddress(&_SelfFundedPingPong.CallOpts)
}

func (_SelfFundedPingPong *SelfFundedPingPongCallerSession) GetCounterpartAddress() (common.Address, error) {
	return _SelfFundedPingPong.Contract.GetCounterpartAddress(&_SelfFundedPingPong.CallOpts)
}

func (_SelfFundedPingPong *SelfFundedPingPongCaller) GetCounterpartChainSelector(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _SelfFundedPingPong.contract.Call(opts, &out, "getCounterpartChainSelector")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_SelfFundedPingPong *SelfFundedPingPongSession) GetCounterpartChainSelector() (uint64, error) {
	return _SelfFundedPingPong.Contract.GetCounterpartChainSelector(&_SelfFundedPingPong.CallOpts)
}

func (_SelfFundedPingPong *SelfFundedPingPongCallerSession) GetCounterpartChainSelector() (uint64, error) {
	return _SelfFundedPingPong.Contract.GetCounterpartChainSelector(&_SelfFundedPingPong.CallOpts)
}

func (_SelfFundedPingPong *SelfFundedPingPongCaller) GetFeeToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SelfFundedPingPong.contract.Call(opts, &out, "getFeeToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_SelfFundedPingPong *SelfFundedPingPongSession) GetFeeToken() (common.Address, error) {
	return _SelfFundedPingPong.Contract.GetFeeToken(&_SelfFundedPingPong.CallOpts)
}

func (_SelfFundedPingPong *SelfFundedPingPongCallerSession) GetFeeToken() (common.Address, error) {
	return _SelfFundedPingPong.Contract.GetFeeToken(&_SelfFundedPingPong.CallOpts)
}

func (_SelfFundedPingPong *SelfFundedPingPongCaller) GetRouter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SelfFundedPingPong.contract.Call(opts, &out, "getRouter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_SelfFundedPingPong *SelfFundedPingPongSession) GetRouter() (common.Address, error) {
	return _SelfFundedPingPong.Contract.GetRouter(&_SelfFundedPingPong.CallOpts)
}

func (_SelfFundedPingPong *SelfFundedPingPongCallerSession) GetRouter() (common.Address, error) {
	return _SelfFundedPingPong.Contract.GetRouter(&_SelfFundedPingPong.CallOpts)
}

func (_SelfFundedPingPong *SelfFundedPingPongCaller) IsPaused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _SelfFundedPingPong.contract.Call(opts, &out, "isPaused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_SelfFundedPingPong *SelfFundedPingPongSession) IsPaused() (bool, error) {
	return _SelfFundedPingPong.Contract.IsPaused(&_SelfFundedPingPong.CallOpts)
}

func (_SelfFundedPingPong *SelfFundedPingPongCallerSession) IsPaused() (bool, error) {
	return _SelfFundedPingPong.Contract.IsPaused(&_SelfFundedPingPong.CallOpts)
}

func (_SelfFundedPingPong *SelfFundedPingPongCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SelfFundedPingPong.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_SelfFundedPingPong *SelfFundedPingPongSession) Owner() (common.Address, error) {
	return _SelfFundedPingPong.Contract.Owner(&_SelfFundedPingPong.CallOpts)
}

func (_SelfFundedPingPong *SelfFundedPingPongCallerSession) Owner() (common.Address, error) {
	return _SelfFundedPingPong.Contract.Owner(&_SelfFundedPingPong.CallOpts)
}

func (_SelfFundedPingPong *SelfFundedPingPongCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _SelfFundedPingPong.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_SelfFundedPingPong *SelfFundedPingPongSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _SelfFundedPingPong.Contract.SupportsInterface(&_SelfFundedPingPong.CallOpts, interfaceId)
}

func (_SelfFundedPingPong *SelfFundedPingPongCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _SelfFundedPingPong.Contract.SupportsInterface(&_SelfFundedPingPong.CallOpts, interfaceId)
}

func (_SelfFundedPingPong *SelfFundedPingPongCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _SelfFundedPingPong.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_SelfFundedPingPong *SelfFundedPingPongSession) TypeAndVersion() (string, error) {
	return _SelfFundedPingPong.Contract.TypeAndVersion(&_SelfFundedPingPong.CallOpts)
}

func (_SelfFundedPingPong *SelfFundedPingPongCallerSession) TypeAndVersion() (string, error) {
	return _SelfFundedPingPong.Contract.TypeAndVersion(&_SelfFundedPingPong.CallOpts)
}

func (_SelfFundedPingPong *SelfFundedPingPongTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SelfFundedPingPong.contract.Transact(opts, "acceptOwnership")
}

func (_SelfFundedPingPong *SelfFundedPingPongSession) AcceptOwnership() (*types.Transaction, error) {
	return _SelfFundedPingPong.Contract.AcceptOwnership(&_SelfFundedPingPong.TransactOpts)
}

func (_SelfFundedPingPong *SelfFundedPingPongTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _SelfFundedPingPong.Contract.AcceptOwnership(&_SelfFundedPingPong.TransactOpts)
}

func (_SelfFundedPingPong *SelfFundedPingPongTransactor) CcipReceive(opts *bind.TransactOpts, message ClientAny2EVMMessage) (*types.Transaction, error) {
	return _SelfFundedPingPong.contract.Transact(opts, "ccipReceive", message)
}

func (_SelfFundedPingPong *SelfFundedPingPongSession) CcipReceive(message ClientAny2EVMMessage) (*types.Transaction, error) {
	return _SelfFundedPingPong.Contract.CcipReceive(&_SelfFundedPingPong.TransactOpts, message)
}

func (_SelfFundedPingPong *SelfFundedPingPongTransactorSession) CcipReceive(message ClientAny2EVMMessage) (*types.Transaction, error) {
	return _SelfFundedPingPong.Contract.CcipReceive(&_SelfFundedPingPong.TransactOpts, message)
}

func (_SelfFundedPingPong *SelfFundedPingPongTransactor) FundPingPong(opts *bind.TransactOpts, pingPongCount *big.Int) (*types.Transaction, error) {
	return _SelfFundedPingPong.contract.Transact(opts, "fundPingPong", pingPongCount)
}

func (_SelfFundedPingPong *SelfFundedPingPongSession) FundPingPong(pingPongCount *big.Int) (*types.Transaction, error) {
	return _SelfFundedPingPong.Contract.FundPingPong(&_SelfFundedPingPong.TransactOpts, pingPongCount)
}

func (_SelfFundedPingPong *SelfFundedPingPongTransactorSession) FundPingPong(pingPongCount *big.Int) (*types.Transaction, error) {
	return _SelfFundedPingPong.Contract.FundPingPong(&_SelfFundedPingPong.TransactOpts, pingPongCount)
}

func (_SelfFundedPingPong *SelfFundedPingPongTransactor) SetCountIncrBeforeFunding(opts *bind.TransactOpts, countIncrBeforeFunding uint8) (*types.Transaction, error) {
	return _SelfFundedPingPong.contract.Transact(opts, "setCountIncrBeforeFunding", countIncrBeforeFunding)
}

func (_SelfFundedPingPong *SelfFundedPingPongSession) SetCountIncrBeforeFunding(countIncrBeforeFunding uint8) (*types.Transaction, error) {
	return _SelfFundedPingPong.Contract.SetCountIncrBeforeFunding(&_SelfFundedPingPong.TransactOpts, countIncrBeforeFunding)
}

func (_SelfFundedPingPong *SelfFundedPingPongTransactorSession) SetCountIncrBeforeFunding(countIncrBeforeFunding uint8) (*types.Transaction, error) {
	return _SelfFundedPingPong.Contract.SetCountIncrBeforeFunding(&_SelfFundedPingPong.TransactOpts, countIncrBeforeFunding)
}

func (_SelfFundedPingPong *SelfFundedPingPongTransactor) SetCounterpart(opts *bind.TransactOpts, counterpartChainSelector uint64, counterpartAddress common.Address) (*types.Transaction, error) {
	return _SelfFundedPingPong.contract.Transact(opts, "setCounterpart", counterpartChainSelector, counterpartAddress)
}

func (_SelfFundedPingPong *SelfFundedPingPongSession) SetCounterpart(counterpartChainSelector uint64, counterpartAddress common.Address) (*types.Transaction, error) {
	return _SelfFundedPingPong.Contract.SetCounterpart(&_SelfFundedPingPong.TransactOpts, counterpartChainSelector, counterpartAddress)
}

func (_SelfFundedPingPong *SelfFundedPingPongTransactorSession) SetCounterpart(counterpartChainSelector uint64, counterpartAddress common.Address) (*types.Transaction, error) {
	return _SelfFundedPingPong.Contract.SetCounterpart(&_SelfFundedPingPong.TransactOpts, counterpartChainSelector, counterpartAddress)
}

func (_SelfFundedPingPong *SelfFundedPingPongTransactor) SetCounterpartAddress(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _SelfFundedPingPong.contract.Transact(opts, "setCounterpartAddress", addr)
}

func (_SelfFundedPingPong *SelfFundedPingPongSession) SetCounterpartAddress(addr common.Address) (*types.Transaction, error) {
	return _SelfFundedPingPong.Contract.SetCounterpartAddress(&_SelfFundedPingPong.TransactOpts, addr)
}

func (_SelfFundedPingPong *SelfFundedPingPongTransactorSession) SetCounterpartAddress(addr common.Address) (*types.Transaction, error) {
	return _SelfFundedPingPong.Contract.SetCounterpartAddress(&_SelfFundedPingPong.TransactOpts, addr)
}

func (_SelfFundedPingPong *SelfFundedPingPongTransactor) SetCounterpartChainSelector(opts *bind.TransactOpts, chainSelector uint64) (*types.Transaction, error) {
	return _SelfFundedPingPong.contract.Transact(opts, "setCounterpartChainSelector", chainSelector)
}

func (_SelfFundedPingPong *SelfFundedPingPongSession) SetCounterpartChainSelector(chainSelector uint64) (*types.Transaction, error) {
	return _SelfFundedPingPong.Contract.SetCounterpartChainSelector(&_SelfFundedPingPong.TransactOpts, chainSelector)
}

func (_SelfFundedPingPong *SelfFundedPingPongTransactorSession) SetCounterpartChainSelector(chainSelector uint64) (*types.Transaction, error) {
	return _SelfFundedPingPong.Contract.SetCounterpartChainSelector(&_SelfFundedPingPong.TransactOpts, chainSelector)
}

func (_SelfFundedPingPong *SelfFundedPingPongTransactor) SetPaused(opts *bind.TransactOpts, pause bool) (*types.Transaction, error) {
	return _SelfFundedPingPong.contract.Transact(opts, "setPaused", pause)
}

func (_SelfFundedPingPong *SelfFundedPingPongSession) SetPaused(pause bool) (*types.Transaction, error) {
	return _SelfFundedPingPong.Contract.SetPaused(&_SelfFundedPingPong.TransactOpts, pause)
}

func (_SelfFundedPingPong *SelfFundedPingPongTransactorSession) SetPaused(pause bool) (*types.Transaction, error) {
	return _SelfFundedPingPong.Contract.SetPaused(&_SelfFundedPingPong.TransactOpts, pause)
}

func (_SelfFundedPingPong *SelfFundedPingPongTransactor) StartPingPong(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SelfFundedPingPong.contract.Transact(opts, "startPingPong")
}

func (_SelfFundedPingPong *SelfFundedPingPongSession) StartPingPong() (*types.Transaction, error) {
	return _SelfFundedPingPong.Contract.StartPingPong(&_SelfFundedPingPong.TransactOpts)
}

func (_SelfFundedPingPong *SelfFundedPingPongTransactorSession) StartPingPong() (*types.Transaction, error) {
	return _SelfFundedPingPong.Contract.StartPingPong(&_SelfFundedPingPong.TransactOpts)
}

func (_SelfFundedPingPong *SelfFundedPingPongTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _SelfFundedPingPong.contract.Transact(opts, "transferOwnership", to)
}

func (_SelfFundedPingPong *SelfFundedPingPongSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _SelfFundedPingPong.Contract.TransferOwnership(&_SelfFundedPingPong.TransactOpts, to)
}

func (_SelfFundedPingPong *SelfFundedPingPongTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _SelfFundedPingPong.Contract.TransferOwnership(&_SelfFundedPingPong.TransactOpts, to)
}

type SelfFundedPingPongCountIncrBeforeFundingSetIterator struct {
	Event *SelfFundedPingPongCountIncrBeforeFundingSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *SelfFundedPingPongCountIncrBeforeFundingSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SelfFundedPingPongCountIncrBeforeFundingSet)
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
		it.Event = new(SelfFundedPingPongCountIncrBeforeFundingSet)
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

func (it *SelfFundedPingPongCountIncrBeforeFundingSetIterator) Error() error {
	return it.fail
}

func (it *SelfFundedPingPongCountIncrBeforeFundingSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type SelfFundedPingPongCountIncrBeforeFundingSet struct {
	CountIncrBeforeFunding uint8
	Raw                    types.Log
}

func (_SelfFundedPingPong *SelfFundedPingPongFilterer) FilterCountIncrBeforeFundingSet(opts *bind.FilterOpts) (*SelfFundedPingPongCountIncrBeforeFundingSetIterator, error) {

	logs, sub, err := _SelfFundedPingPong.contract.FilterLogs(opts, "CountIncrBeforeFundingSet")
	if err != nil {
		return nil, err
	}
	return &SelfFundedPingPongCountIncrBeforeFundingSetIterator{contract: _SelfFundedPingPong.contract, event: "CountIncrBeforeFundingSet", logs: logs, sub: sub}, nil
}

func (_SelfFundedPingPong *SelfFundedPingPongFilterer) WatchCountIncrBeforeFundingSet(opts *bind.WatchOpts, sink chan<- *SelfFundedPingPongCountIncrBeforeFundingSet) (event.Subscription, error) {

	logs, sub, err := _SelfFundedPingPong.contract.WatchLogs(opts, "CountIncrBeforeFundingSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(SelfFundedPingPongCountIncrBeforeFundingSet)
				if err := _SelfFundedPingPong.contract.UnpackLog(event, "CountIncrBeforeFundingSet", log); err != nil {
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

func (_SelfFundedPingPong *SelfFundedPingPongFilterer) ParseCountIncrBeforeFundingSet(log types.Log) (*SelfFundedPingPongCountIncrBeforeFundingSet, error) {
	event := new(SelfFundedPingPongCountIncrBeforeFundingSet)
	if err := _SelfFundedPingPong.contract.UnpackLog(event, "CountIncrBeforeFundingSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type SelfFundedPingPongFundedIterator struct {
	Event *SelfFundedPingPongFunded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *SelfFundedPingPongFundedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SelfFundedPingPongFunded)
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
		it.Event = new(SelfFundedPingPongFunded)
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

func (it *SelfFundedPingPongFundedIterator) Error() error {
	return it.fail
}

func (it *SelfFundedPingPongFundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type SelfFundedPingPongFunded struct {
	Raw types.Log
}

func (_SelfFundedPingPong *SelfFundedPingPongFilterer) FilterFunded(opts *bind.FilterOpts) (*SelfFundedPingPongFundedIterator, error) {

	logs, sub, err := _SelfFundedPingPong.contract.FilterLogs(opts, "Funded")
	if err != nil {
		return nil, err
	}
	return &SelfFundedPingPongFundedIterator{contract: _SelfFundedPingPong.contract, event: "Funded", logs: logs, sub: sub}, nil
}

func (_SelfFundedPingPong *SelfFundedPingPongFilterer) WatchFunded(opts *bind.WatchOpts, sink chan<- *SelfFundedPingPongFunded) (event.Subscription, error) {

	logs, sub, err := _SelfFundedPingPong.contract.WatchLogs(opts, "Funded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(SelfFundedPingPongFunded)
				if err := _SelfFundedPingPong.contract.UnpackLog(event, "Funded", log); err != nil {
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

func (_SelfFundedPingPong *SelfFundedPingPongFilterer) ParseFunded(log types.Log) (*SelfFundedPingPongFunded, error) {
	event := new(SelfFundedPingPongFunded)
	if err := _SelfFundedPingPong.contract.UnpackLog(event, "Funded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type SelfFundedPingPongOwnershipTransferRequestedIterator struct {
	Event *SelfFundedPingPongOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *SelfFundedPingPongOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SelfFundedPingPongOwnershipTransferRequested)
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
		it.Event = new(SelfFundedPingPongOwnershipTransferRequested)
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

func (it *SelfFundedPingPongOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *SelfFundedPingPongOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type SelfFundedPingPongOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_SelfFundedPingPong *SelfFundedPingPongFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*SelfFundedPingPongOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _SelfFundedPingPong.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &SelfFundedPingPongOwnershipTransferRequestedIterator{contract: _SelfFundedPingPong.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_SelfFundedPingPong *SelfFundedPingPongFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *SelfFundedPingPongOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _SelfFundedPingPong.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(SelfFundedPingPongOwnershipTransferRequested)
				if err := _SelfFundedPingPong.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_SelfFundedPingPong *SelfFundedPingPongFilterer) ParseOwnershipTransferRequested(log types.Log) (*SelfFundedPingPongOwnershipTransferRequested, error) {
	event := new(SelfFundedPingPongOwnershipTransferRequested)
	if err := _SelfFundedPingPong.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type SelfFundedPingPongOwnershipTransferredIterator struct {
	Event *SelfFundedPingPongOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *SelfFundedPingPongOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SelfFundedPingPongOwnershipTransferred)
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
		it.Event = new(SelfFundedPingPongOwnershipTransferred)
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

func (it *SelfFundedPingPongOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *SelfFundedPingPongOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type SelfFundedPingPongOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_SelfFundedPingPong *SelfFundedPingPongFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*SelfFundedPingPongOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _SelfFundedPingPong.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &SelfFundedPingPongOwnershipTransferredIterator{contract: _SelfFundedPingPong.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_SelfFundedPingPong *SelfFundedPingPongFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *SelfFundedPingPongOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _SelfFundedPingPong.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(SelfFundedPingPongOwnershipTransferred)
				if err := _SelfFundedPingPong.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_SelfFundedPingPong *SelfFundedPingPongFilterer) ParseOwnershipTransferred(log types.Log) (*SelfFundedPingPongOwnershipTransferred, error) {
	event := new(SelfFundedPingPongOwnershipTransferred)
	if err := _SelfFundedPingPong.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type SelfFundedPingPongPingIterator struct {
	Event *SelfFundedPingPongPing

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *SelfFundedPingPongPingIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SelfFundedPingPongPing)
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
		it.Event = new(SelfFundedPingPongPing)
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

func (it *SelfFundedPingPongPingIterator) Error() error {
	return it.fail
}

func (it *SelfFundedPingPongPingIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type SelfFundedPingPongPing struct {
	PingPongCount *big.Int
	Raw           types.Log
}

func (_SelfFundedPingPong *SelfFundedPingPongFilterer) FilterPing(opts *bind.FilterOpts) (*SelfFundedPingPongPingIterator, error) {

	logs, sub, err := _SelfFundedPingPong.contract.FilterLogs(opts, "Ping")
	if err != nil {
		return nil, err
	}
	return &SelfFundedPingPongPingIterator{contract: _SelfFundedPingPong.contract, event: "Ping", logs: logs, sub: sub}, nil
}

func (_SelfFundedPingPong *SelfFundedPingPongFilterer) WatchPing(opts *bind.WatchOpts, sink chan<- *SelfFundedPingPongPing) (event.Subscription, error) {

	logs, sub, err := _SelfFundedPingPong.contract.WatchLogs(opts, "Ping")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(SelfFundedPingPongPing)
				if err := _SelfFundedPingPong.contract.UnpackLog(event, "Ping", log); err != nil {
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

func (_SelfFundedPingPong *SelfFundedPingPongFilterer) ParsePing(log types.Log) (*SelfFundedPingPongPing, error) {
	event := new(SelfFundedPingPongPing)
	if err := _SelfFundedPingPong.contract.UnpackLog(event, "Ping", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type SelfFundedPingPongPongIterator struct {
	Event *SelfFundedPingPongPong

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *SelfFundedPingPongPongIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SelfFundedPingPongPong)
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
		it.Event = new(SelfFundedPingPongPong)
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

func (it *SelfFundedPingPongPongIterator) Error() error {
	return it.fail
}

func (it *SelfFundedPingPongPongIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type SelfFundedPingPongPong struct {
	PingPongCount *big.Int
	Raw           types.Log
}

func (_SelfFundedPingPong *SelfFundedPingPongFilterer) FilterPong(opts *bind.FilterOpts) (*SelfFundedPingPongPongIterator, error) {

	logs, sub, err := _SelfFundedPingPong.contract.FilterLogs(opts, "Pong")
	if err != nil {
		return nil, err
	}
	return &SelfFundedPingPongPongIterator{contract: _SelfFundedPingPong.contract, event: "Pong", logs: logs, sub: sub}, nil
}

func (_SelfFundedPingPong *SelfFundedPingPongFilterer) WatchPong(opts *bind.WatchOpts, sink chan<- *SelfFundedPingPongPong) (event.Subscription, error) {

	logs, sub, err := _SelfFundedPingPong.contract.WatchLogs(opts, "Pong")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(SelfFundedPingPongPong)
				if err := _SelfFundedPingPong.contract.UnpackLog(event, "Pong", log); err != nil {
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

func (_SelfFundedPingPong *SelfFundedPingPongFilterer) ParsePong(log types.Log) (*SelfFundedPingPongPong, error) {
	event := new(SelfFundedPingPongPong)
	if err := _SelfFundedPingPong.contract.UnpackLog(event, "Pong", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_SelfFundedPingPong *SelfFundedPingPong) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _SelfFundedPingPong.abi.Events["CountIncrBeforeFundingSet"].ID:
		return _SelfFundedPingPong.ParseCountIncrBeforeFundingSet(log)
	case _SelfFundedPingPong.abi.Events["Funded"].ID:
		return _SelfFundedPingPong.ParseFunded(log)
	case _SelfFundedPingPong.abi.Events["OwnershipTransferRequested"].ID:
		return _SelfFundedPingPong.ParseOwnershipTransferRequested(log)
	case _SelfFundedPingPong.abi.Events["OwnershipTransferred"].ID:
		return _SelfFundedPingPong.ParseOwnershipTransferred(log)
	case _SelfFundedPingPong.abi.Events["Ping"].ID:
		return _SelfFundedPingPong.ParsePing(log)
	case _SelfFundedPingPong.abi.Events["Pong"].ID:
		return _SelfFundedPingPong.ParsePong(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (SelfFundedPingPongCountIncrBeforeFundingSet) Topic() common.Hash {
	return common.HexToHash("0x4768dbf8645b24c54f2887651545d24f748c0d0d1d4c689eb810fb19f0befcf3")
}

func (SelfFundedPingPongFunded) Topic() common.Hash {
	return common.HexToHash("0x302777af5d26fab9dd5120c5f1307c65193ebc51daf33244ada4365fab10602c")
}

func (SelfFundedPingPongOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (SelfFundedPingPongOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (SelfFundedPingPongPing) Topic() common.Hash {
	return common.HexToHash("0x48257dc961b6f792c2b78a080dacfed693b660960a702de21cee364e20270e2f")
}

func (SelfFundedPingPongPong) Topic() common.Hash {
	return common.HexToHash("0x58b69f57828e6962d216502094c54f6562f3bf082ba758966c3454f9e37b1525")
}

func (_SelfFundedPingPong *SelfFundedPingPong) Address() common.Address {
	return _SelfFundedPingPong.address
}

type SelfFundedPingPongInterface interface {
	GetCountIncrBeforeFunding(opts *bind.CallOpts) (uint8, error)

	GetCounterpartAddress(opts *bind.CallOpts) (common.Address, error)

	GetCounterpartChainSelector(opts *bind.CallOpts) (uint64, error)

	GetFeeToken(opts *bind.CallOpts) (common.Address, error)

	GetRouter(opts *bind.CallOpts) (common.Address, error)

	IsPaused(opts *bind.CallOpts) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	CcipReceive(opts *bind.TransactOpts, message ClientAny2EVMMessage) (*types.Transaction, error)

	FundPingPong(opts *bind.TransactOpts, pingPongCount *big.Int) (*types.Transaction, error)

	SetCountIncrBeforeFunding(opts *bind.TransactOpts, countIncrBeforeFunding uint8) (*types.Transaction, error)

	SetCounterpart(opts *bind.TransactOpts, counterpartChainSelector uint64, counterpartAddress common.Address) (*types.Transaction, error)

	SetCounterpartAddress(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error)

	SetCounterpartChainSelector(opts *bind.TransactOpts, chainSelector uint64) (*types.Transaction, error)

	SetPaused(opts *bind.TransactOpts, pause bool) (*types.Transaction, error)

	StartPingPong(opts *bind.TransactOpts) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterCountIncrBeforeFundingSet(opts *bind.FilterOpts) (*SelfFundedPingPongCountIncrBeforeFundingSetIterator, error)

	WatchCountIncrBeforeFundingSet(opts *bind.WatchOpts, sink chan<- *SelfFundedPingPongCountIncrBeforeFundingSet) (event.Subscription, error)

	ParseCountIncrBeforeFundingSet(log types.Log) (*SelfFundedPingPongCountIncrBeforeFundingSet, error)

	FilterFunded(opts *bind.FilterOpts) (*SelfFundedPingPongFundedIterator, error)

	WatchFunded(opts *bind.WatchOpts, sink chan<- *SelfFundedPingPongFunded) (event.Subscription, error)

	ParseFunded(log types.Log) (*SelfFundedPingPongFunded, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*SelfFundedPingPongOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *SelfFundedPingPongOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*SelfFundedPingPongOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*SelfFundedPingPongOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *SelfFundedPingPongOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*SelfFundedPingPongOwnershipTransferred, error)

	FilterPing(opts *bind.FilterOpts) (*SelfFundedPingPongPingIterator, error)

	WatchPing(opts *bind.WatchOpts, sink chan<- *SelfFundedPingPongPing) (event.Subscription, error)

	ParsePing(log types.Log) (*SelfFundedPingPongPing, error)

	FilterPong(opts *bind.FilterOpts) (*SelfFundedPingPongPongIterator, error)

	WatchPong(opts *bind.WatchOpts, sink chan<- *SelfFundedPingPongPong) (event.Subscription, error)

	ParsePong(log types.Log) (*SelfFundedPingPongPong, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
