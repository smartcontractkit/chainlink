// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ethereum

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
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

// VRFv2ConsumerMetaData contains all meta data concerning the VRFv2Consumer contract.
var VRFv2ConsumerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vrfCoordinator\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlyCoordinatorCanFulfill\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"RequestFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"}],\"name\":\"RequestSent\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_requestId\",\"type\":\"uint256\"}],\"name\":\"getRequestStatus\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"fulfilled\",\"type\":\"bool\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastRequestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"requestIds\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"requestRandomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_requests\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"fulfilled\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"exists\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b50604051620018a4380380620018a483398181016040528101906200003791906200034e565b33806000838073ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff168152505050600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603620000e3576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401620000da90620003e1565b60405180910390fd5b816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16146200016a576200016981620001b560201b60201c565b5b50505080600360006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505062000475565b3373ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff160362000226576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016200021d9062000453565b60405180910390fd5b80600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508073ffffffffffffffffffffffffffffffffffffffff1660008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127860405160405180910390a350565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006200031682620002e9565b9050919050565b620003288162000309565b81146200033457600080fd5b50565b60008151905062000348816200031d565b92915050565b600060208284031215620003675762000366620002e4565b5b6000620003778482850162000337565b91505092915050565b600082825260208201905092915050565b7f43616e6e6f7420736574206f776e657220746f207a65726f0000000000000000600082015250565b6000620003c960188362000380565b9150620003d68262000391565b602082019050919050565b60006020820190508181036000830152620003fc81620003ba565b9050919050565b7f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000600082015250565b60006200043b60178362000380565b9150620004488262000403565b602082019050919050565b600060208201905081810360008301526200046e816200042c565b9050919050565b60805161140c62000498600039600081816101da015261022e015261140c6000f3fe608060405234801561001057600080fd5b50600436106100935760003560e01c80639561f023116100665780639561f0231461010c578063a168fa891461013c578063d8a4676f1461016d578063f2fde38b1461019e578063fc2a88c3146101ba57610093565b80631fe543e31461009857806379ba5097146100b45780638796ba8c146100be5780638da5cb5b146100ee575b600080fd5b6100b260048036038101906100ad9190610cc1565b6101d8565b005b6100bc610298565b005b6100d860048036038101906100d39190610d1d565b61042d565b6040516100e59190610d59565b60405180910390f35b6100f6610451565b6040516101039190610db5565b60405180910390f35b61012660048036038101906101219190610ebc565b61047a565b6040516101339190610d59565b60405180910390f35b61015660048036038101906101519190610d1d565b61067b565b604051610164929190610f52565b60405180910390f35b61018760048036038101906101829190610d1d565b6106b9565b604051610195929190611039565b60405180910390f35b6101b860048036038101906101b39190611095565b6107e4565b005b6101c26107f8565b6040516101cf9190610d59565b60405180910390f35b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461028a57337f00000000000000000000000000000000000000000000000000000000000000006040517f1cf993f40000000000000000000000000000000000000000000000000000000081526004016102819291906110c2565b60405180910390fd5b61029482826107fe565b5050565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610328576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161031f90611148565b60405180910390fd5b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506000600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055503373ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a350565b6004818154811061043d57600080fd5b906000526020600020016000915090505481565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b60006104846108f8565b600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16635d3b1d3083888789886040518663ffffffff1660e01b81526004016104e79594939291906111a4565b6020604051808303816000875af1158015610506573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061052a919061120c565b90506040518060600160405280600015158152602001600115158152602001600067ffffffffffffffff81111561056457610563610b7e565b5b6040519080825280602002602001820160405280156105925781602001602082028036833780820191505090505b508152506002600083815260200190815260200160002060008201518160000160006101000a81548160ff02191690831515021790555060208201518160000160016101000a81548160ff0219169083151502179055506040820151816001019080519060200190610605929190610ab4565b509050506004819080600181540180825580915050600190039060005260206000200160009091909190915055806005819055507fcc58b13ad3eab50626c6a6300b1d139cd6ebb1688a7cced9461c2f7e762665ee818460405161066a929190611239565b60405180910390a195945050505050565b60026020528060005260406000206000915090508060000160009054906101000a900460ff16908060000160019054906101000a900460ff16905082565b600060606002600084815260200190815260200160002060000160019054906101000a900460ff16610720576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610717906112ae565b60405180910390fd5b6000600260008581526020019081526020016000206040518060600160405290816000820160009054906101000a900460ff161515151581526020016000820160019054906101000a900460ff16151515158152602001600182018054806020026020016040519081016040528092919081815260200182805480156107c557602002820191906000526020600020905b8154815260200190600101908083116107b1575b5050505050815250509050806000015181604001519250925050915091565b6107ec6108f8565b6107f581610988565b50565b60055481565b6002600083815260200190815260200160002060000160019054906101000a900460ff16610861576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610858906112ae565b60405180910390fd5b60016002600084815260200190815260200160002060000160006101000a81548160ff021916908315150217905550806002600084815260200190815260200160002060010190805190602001906108ba929190610ab4565b507ffe2e2d779dba245964d4e3ef9b994be63856fd568bf7d3ca9e224755cb1bd54d82826040516108ec9291906112ce565b60405180910390a15050565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610986576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161097d9061134a565b60405180910390fd5b565b3373ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036109f6576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016109ed906113b6565b60405180910390fd5b80600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508073ffffffffffffffffffffffffffffffffffffffff1660008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127860405160405180910390a350565b828054828255906000526020600020908101928215610af0579160200282015b82811115610aef578251825591602001919060010190610ad4565b5b509050610afd9190610b01565b5090565b5b80821115610b1a576000816000905550600101610b02565b5090565b6000604051905090565b600080fd5b600080fd5b6000819050919050565b610b4581610b32565b8114610b5057600080fd5b50565b600081359050610b6281610b3c565b92915050565b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b610bb682610b6d565b810181811067ffffffffffffffff82111715610bd557610bd4610b7e565b5b80604052505050565b6000610be8610b1e565b9050610bf48282610bad565b919050565b600067ffffffffffffffff821115610c1457610c13610b7e565b5b602082029050602081019050919050565b600080fd5b6000610c3d610c3884610bf9565b610bde565b90508083825260208201905060208402830185811115610c6057610c5f610c25565b5b835b81811015610c895780610c758882610b53565b845260208401935050602081019050610c62565b5050509392505050565b600082601f830112610ca857610ca7610b68565b5b8135610cb8848260208601610c2a565b91505092915050565b60008060408385031215610cd857610cd7610b28565b5b6000610ce685828601610b53565b925050602083013567ffffffffffffffff811115610d0757610d06610b2d565b5b610d1385828601610c93565b9150509250929050565b600060208284031215610d3357610d32610b28565b5b6000610d4184828501610b53565b91505092915050565b610d5381610b32565b82525050565b6000602082019050610d6e6000830184610d4a565b92915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610d9f82610d74565b9050919050565b610daf81610d94565b82525050565b6000602082019050610dca6000830184610da6565b92915050565b600067ffffffffffffffff82169050919050565b610ded81610dd0565b8114610df857600080fd5b50565b600081359050610e0a81610de4565b92915050565b600063ffffffff82169050919050565b610e2981610e10565b8114610e3457600080fd5b50565b600081359050610e4681610e20565b92915050565b600061ffff82169050919050565b610e6381610e4c565b8114610e6e57600080fd5b50565b600081359050610e8081610e5a565b92915050565b6000819050919050565b610e9981610e86565b8114610ea457600080fd5b50565b600081359050610eb681610e90565b92915050565b600080600080600060a08688031215610ed857610ed7610b28565b5b6000610ee688828901610dfb565b9550506020610ef788828901610e37565b9450506040610f0888828901610e71565b9350506060610f1988828901610e37565b9250506080610f2a88828901610ea7565b9150509295509295909350565b60008115159050919050565b610f4c81610f37565b82525050565b6000604082019050610f676000830185610f43565b610f746020830184610f43565b9392505050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b610fb081610b32565b82525050565b6000610fc28383610fa7565b60208301905092915050565b6000602082019050919050565b6000610fe682610f7b565b610ff08185610f86565b9350610ffb83610f97565b8060005b8381101561102c5781516110138882610fb6565b975061101e83610fce565b925050600181019050610fff565b5085935050505092915050565b600060408201905061104e6000830185610f43565b81810360208301526110608184610fdb565b90509392505050565b61107281610d94565b811461107d57600080fd5b50565b60008135905061108f81611069565b92915050565b6000602082840312156110ab576110aa610b28565b5b60006110b984828501611080565b91505092915050565b60006040820190506110d76000830185610da6565b6110e46020830184610da6565b9392505050565b600082825260208201905092915050565b7f4d7573742062652070726f706f736564206f776e657200000000000000000000600082015250565b60006111326016836110eb565b915061113d826110fc565b602082019050919050565b6000602082019050818103600083015261116181611125565b9050919050565b61117181610e86565b82525050565b61118081610dd0565b82525050565b61118f81610e4c565b82525050565b61119e81610e10565b82525050565b600060a0820190506111b96000830188611168565b6111c66020830187611177565b6111d36040830186611186565b6111e06060830185611195565b6111ed6080830184611195565b9695505050505050565b60008151905061120681610b3c565b92915050565b60006020828403121561122257611221610b28565b5b6000611230848285016111f7565b91505092915050565b600060408201905061124e6000830185610d4a565b61125b6020830184611195565b9392505050565b7f72657175657374206e6f7420666f756e64000000000000000000000000000000600082015250565b60006112986011836110eb565b91506112a382611262565b602082019050919050565b600060208201905081810360008301526112c78161128b565b9050919050565b60006040820190506112e36000830185610d4a565b81810360208301526112f58184610fdb565b90509392505050565b7f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000600082015250565b60006113346016836110eb565b915061133f826112fe565b602082019050919050565b6000602082019050818103600083015261136381611327565b9050919050565b7f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000600082015250565b60006113a06017836110eb565b91506113ab8261136a565b602082019050919050565b600060208201905081810360008301526113cf81611393565b905091905056fea2646970667358221220e3ffc5c5c5c9efece6a6dd560edba3d094c9853bc2c4400e1d3ea7e3d950ecfa64736f6c63430008110033",
}

// VRFv2ConsumerABI is the input ABI used to generate the binding from.
// Deprecated: Use VRFv2ConsumerMetaData.ABI instead.
var VRFv2ConsumerABI = VRFv2ConsumerMetaData.ABI

// VRFv2ConsumerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use VRFv2ConsumerMetaData.Bin instead.
var VRFv2ConsumerBin = VRFv2ConsumerMetaData.Bin

// DeployVRFv2Consumer deploys a new Ethereum contract, binding an instance of VRFv2Consumer to it.
func DeployVRFv2Consumer(auth *bind.TransactOpts, backend bind.ContractBackend, vrfCoordinator common.Address) (common.Address, *types.Transaction, *VRFv2Consumer, error) {
	parsed, err := VRFv2ConsumerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFv2ConsumerBin), backend, vrfCoordinator)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFv2Consumer{VRFv2ConsumerCaller: VRFv2ConsumerCaller{contract: contract}, VRFv2ConsumerTransactor: VRFv2ConsumerTransactor{contract: contract}, VRFv2ConsumerFilterer: VRFv2ConsumerFilterer{contract: contract}}, nil
}

// VRFv2Consumer is an auto generated Go binding around an Ethereum contract.
type VRFv2Consumer struct {
	VRFv2ConsumerCaller     // Read-only binding to the contract
	VRFv2ConsumerTransactor // Write-only binding to the contract
	VRFv2ConsumerFilterer   // Log filterer for contract events
}

// VRFv2ConsumerCaller is an auto generated read-only Go binding around an Ethereum contract.
type VRFv2ConsumerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFv2ConsumerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VRFv2ConsumerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFv2ConsumerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VRFv2ConsumerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFv2ConsumerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VRFv2ConsumerSession struct {
	Contract     *VRFv2Consumer    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VRFv2ConsumerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VRFv2ConsumerCallerSession struct {
	Contract *VRFv2ConsumerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// VRFv2ConsumerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VRFv2ConsumerTransactorSession struct {
	Contract     *VRFv2ConsumerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// VRFv2ConsumerRaw is an auto generated low-level Go binding around an Ethereum contract.
type VRFv2ConsumerRaw struct {
	Contract *VRFv2Consumer // Generic contract binding to access the raw methods on
}

// VRFv2ConsumerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VRFv2ConsumerCallerRaw struct {
	Contract *VRFv2ConsumerCaller // Generic read-only contract binding to access the raw methods on
}

// VRFv2ConsumerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VRFv2ConsumerTransactorRaw struct {
	Contract *VRFv2ConsumerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVRFv2Consumer creates a new instance of VRFv2Consumer, bound to a specific deployed contract.
func NewVRFv2Consumer(address common.Address, backend bind.ContractBackend) (*VRFv2Consumer, error) {
	contract, err := bindVRFv2Consumer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFv2Consumer{VRFv2ConsumerCaller: VRFv2ConsumerCaller{contract: contract}, VRFv2ConsumerTransactor: VRFv2ConsumerTransactor{contract: contract}, VRFv2ConsumerFilterer: VRFv2ConsumerFilterer{contract: contract}}, nil
}

// NewVRFv2ConsumerCaller creates a new read-only instance of VRFv2Consumer, bound to a specific deployed contract.
func NewVRFv2ConsumerCaller(address common.Address, caller bind.ContractCaller) (*VRFv2ConsumerCaller, error) {
	contract, err := bindVRFv2Consumer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFv2ConsumerCaller{contract: contract}, nil
}

// NewVRFv2ConsumerTransactor creates a new write-only instance of VRFv2Consumer, bound to a specific deployed contract.
func NewVRFv2ConsumerTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFv2ConsumerTransactor, error) {
	contract, err := bindVRFv2Consumer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFv2ConsumerTransactor{contract: contract}, nil
}

// NewVRFv2ConsumerFilterer creates a new log filterer instance of VRFv2Consumer, bound to a specific deployed contract.
func NewVRFv2ConsumerFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFv2ConsumerFilterer, error) {
	contract, err := bindVRFv2Consumer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFv2ConsumerFilterer{contract: contract}, nil
}

// bindVRFv2Consumer binds a generic wrapper to an already deployed contract.
func bindVRFv2Consumer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFv2ConsumerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VRFv2Consumer *VRFv2ConsumerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFv2Consumer.Contract.VRFv2ConsumerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VRFv2Consumer *VRFv2ConsumerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFv2Consumer.Contract.VRFv2ConsumerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VRFv2Consumer *VRFv2ConsumerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFv2Consumer.Contract.VRFv2ConsumerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VRFv2Consumer *VRFv2ConsumerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFv2Consumer.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VRFv2Consumer *VRFv2ConsumerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFv2Consumer.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VRFv2Consumer *VRFv2ConsumerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFv2Consumer.Contract.contract.Transact(opts, method, params...)
}

// GetRequestStatus is a free data retrieval call binding the contract method 0xd8a4676f.
//
// Solidity: function getRequestStatus(uint256 _requestId) view returns(bool fulfilled, uint256[] randomWords)
func (_VRFv2Consumer *VRFv2ConsumerCaller) GetRequestStatus(opts *bind.CallOpts, _requestId *big.Int) (struct {
	Fulfilled   bool
	RandomWords []*big.Int
}, error) {
	var out []interface{}
	err := _VRFv2Consumer.contract.Call(opts, &out, "getRequestStatus", _requestId)

	outstruct := new(struct {
		Fulfilled   bool
		RandomWords []*big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Fulfilled = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.RandomWords = *abi.ConvertType(out[1], new([]*big.Int)).(*[]*big.Int)

	return *outstruct, err

}

// GetRequestStatus is a free data retrieval call binding the contract method 0xd8a4676f.
//
// Solidity: function getRequestStatus(uint256 _requestId) view returns(bool fulfilled, uint256[] randomWords)
func (_VRFv2Consumer *VRFv2ConsumerSession) GetRequestStatus(_requestId *big.Int) (struct {
	Fulfilled   bool
	RandomWords []*big.Int
}, error) {
	return _VRFv2Consumer.Contract.GetRequestStatus(&_VRFv2Consumer.CallOpts, _requestId)
}

// GetRequestStatus is a free data retrieval call binding the contract method 0xd8a4676f.
//
// Solidity: function getRequestStatus(uint256 _requestId) view returns(bool fulfilled, uint256[] randomWords)
func (_VRFv2Consumer *VRFv2ConsumerCallerSession) GetRequestStatus(_requestId *big.Int) (struct {
	Fulfilled   bool
	RandomWords []*big.Int
}, error) {
	return _VRFv2Consumer.Contract.GetRequestStatus(&_VRFv2Consumer.CallOpts, _requestId)
}

// LastRequestId is a free data retrieval call binding the contract method 0xfc2a88c3.
//
// Solidity: function lastRequestId() view returns(uint256)
func (_VRFv2Consumer *VRFv2ConsumerCaller) LastRequestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFv2Consumer.contract.Call(opts, &out, "lastRequestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastRequestId is a free data retrieval call binding the contract method 0xfc2a88c3.
//
// Solidity: function lastRequestId() view returns(uint256)
func (_VRFv2Consumer *VRFv2ConsumerSession) LastRequestId() (*big.Int, error) {
	return _VRFv2Consumer.Contract.LastRequestId(&_VRFv2Consumer.CallOpts)
}

// LastRequestId is a free data retrieval call binding the contract method 0xfc2a88c3.
//
// Solidity: function lastRequestId() view returns(uint256)
func (_VRFv2Consumer *VRFv2ConsumerCallerSession) LastRequestId() (*big.Int, error) {
	return _VRFv2Consumer.Contract.LastRequestId(&_VRFv2Consumer.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_VRFv2Consumer *VRFv2ConsumerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFv2Consumer.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_VRFv2Consumer *VRFv2ConsumerSession) Owner() (common.Address, error) {
	return _VRFv2Consumer.Contract.Owner(&_VRFv2Consumer.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_VRFv2Consumer *VRFv2ConsumerCallerSession) Owner() (common.Address, error) {
	return _VRFv2Consumer.Contract.Owner(&_VRFv2Consumer.CallOpts)
}

// RequestIds is a free data retrieval call binding the contract method 0x8796ba8c.
//
// Solidity: function requestIds(uint256 ) view returns(uint256)
func (_VRFv2Consumer *VRFv2ConsumerCaller) RequestIds(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VRFv2Consumer.contract.Call(opts, &out, "requestIds", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RequestIds is a free data retrieval call binding the contract method 0x8796ba8c.
//
// Solidity: function requestIds(uint256 ) view returns(uint256)
func (_VRFv2Consumer *VRFv2ConsumerSession) RequestIds(arg0 *big.Int) (*big.Int, error) {
	return _VRFv2Consumer.Contract.RequestIds(&_VRFv2Consumer.CallOpts, arg0)
}

// RequestIds is a free data retrieval call binding the contract method 0x8796ba8c.
//
// Solidity: function requestIds(uint256 ) view returns(uint256)
func (_VRFv2Consumer *VRFv2ConsumerCallerSession) RequestIds(arg0 *big.Int) (*big.Int, error) {
	return _VRFv2Consumer.Contract.RequestIds(&_VRFv2Consumer.CallOpts, arg0)
}

// SRequests is a free data retrieval call binding the contract method 0xa168fa89.
//
// Solidity: function s_requests(uint256 ) view returns(bool fulfilled, bool exists)
func (_VRFv2Consumer *VRFv2ConsumerCaller) SRequests(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Fulfilled bool
	Exists    bool
}, error) {
	var out []interface{}
	err := _VRFv2Consumer.contract.Call(opts, &out, "s_requests", arg0)

	outstruct := new(struct {
		Fulfilled bool
		Exists    bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Fulfilled = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Exists = *abi.ConvertType(out[1], new(bool)).(*bool)

	return *outstruct, err

}

// SRequests is a free data retrieval call binding the contract method 0xa168fa89.
//
// Solidity: function s_requests(uint256 ) view returns(bool fulfilled, bool exists)
func (_VRFv2Consumer *VRFv2ConsumerSession) SRequests(arg0 *big.Int) (struct {
	Fulfilled bool
	Exists    bool
}, error) {
	return _VRFv2Consumer.Contract.SRequests(&_VRFv2Consumer.CallOpts, arg0)
}

// SRequests is a free data retrieval call binding the contract method 0xa168fa89.
//
// Solidity: function s_requests(uint256 ) view returns(bool fulfilled, bool exists)
func (_VRFv2Consumer *VRFv2ConsumerCallerSession) SRequests(arg0 *big.Int) (struct {
	Fulfilled bool
	Exists    bool
}, error) {
	return _VRFv2Consumer.Contract.SRequests(&_VRFv2Consumer.CallOpts, arg0)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_VRFv2Consumer *VRFv2ConsumerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFv2Consumer.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_VRFv2Consumer *VRFv2ConsumerSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFv2Consumer.Contract.AcceptOwnership(&_VRFv2Consumer.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_VRFv2Consumer *VRFv2ConsumerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFv2Consumer.Contract.AcceptOwnership(&_VRFv2Consumer.TransactOpts)
}

// RawFulfillRandomWords is a paid mutator transaction binding the contract method 0x1fe543e3.
//
// Solidity: function rawFulfillRandomWords(uint256 requestId, uint256[] randomWords) returns()
func (_VRFv2Consumer *VRFv2ConsumerTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFv2Consumer.contract.Transact(opts, "rawFulfillRandomWords", requestId, randomWords)
}

// RawFulfillRandomWords is a paid mutator transaction binding the contract method 0x1fe543e3.
//
// Solidity: function rawFulfillRandomWords(uint256 requestId, uint256[] randomWords) returns()
func (_VRFv2Consumer *VRFv2ConsumerSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFv2Consumer.Contract.RawFulfillRandomWords(&_VRFv2Consumer.TransactOpts, requestId, randomWords)
}

// RawFulfillRandomWords is a paid mutator transaction binding the contract method 0x1fe543e3.
//
// Solidity: function rawFulfillRandomWords(uint256 requestId, uint256[] randomWords) returns()
func (_VRFv2Consumer *VRFv2ConsumerTransactorSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFv2Consumer.Contract.RawFulfillRandomWords(&_VRFv2Consumer.TransactOpts, requestId, randomWords)
}

// RequestRandomWords is a paid mutator transaction binding the contract method 0x9561f023.
//
// Solidity: function requestRandomWords(uint64 subId, uint32 callbackGasLimit, uint16 requestConfirmations, uint32 numWords, bytes32 keyHash) returns(uint256 requestId)
func (_VRFv2Consumer *VRFv2ConsumerTransactor) RequestRandomWords(opts *bind.TransactOpts, subId uint64, callbackGasLimit uint32, requestConfirmations uint16, numWords uint32, keyHash [32]byte) (*types.Transaction, error) {
	return _VRFv2Consumer.contract.Transact(opts, "requestRandomWords", subId, callbackGasLimit, requestConfirmations, numWords, keyHash)
}

// RequestRandomWords is a paid mutator transaction binding the contract method 0x9561f023.
//
// Solidity: function requestRandomWords(uint64 subId, uint32 callbackGasLimit, uint16 requestConfirmations, uint32 numWords, bytes32 keyHash) returns(uint256 requestId)
func (_VRFv2Consumer *VRFv2ConsumerSession) RequestRandomWords(subId uint64, callbackGasLimit uint32, requestConfirmations uint16, numWords uint32, keyHash [32]byte) (*types.Transaction, error) {
	return _VRFv2Consumer.Contract.RequestRandomWords(&_VRFv2Consumer.TransactOpts, subId, callbackGasLimit, requestConfirmations, numWords, keyHash)
}

// RequestRandomWords is a paid mutator transaction binding the contract method 0x9561f023.
//
// Solidity: function requestRandomWords(uint64 subId, uint32 callbackGasLimit, uint16 requestConfirmations, uint32 numWords, bytes32 keyHash) returns(uint256 requestId)
func (_VRFv2Consumer *VRFv2ConsumerTransactorSession) RequestRandomWords(subId uint64, callbackGasLimit uint32, requestConfirmations uint16, numWords uint32, keyHash [32]byte) (*types.Transaction, error) {
	return _VRFv2Consumer.Contract.RequestRandomWords(&_VRFv2Consumer.TransactOpts, subId, callbackGasLimit, requestConfirmations, numWords, keyHash)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_VRFv2Consumer *VRFv2ConsumerTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFv2Consumer.contract.Transact(opts, "transferOwnership", to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_VRFv2Consumer *VRFv2ConsumerSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFv2Consumer.Contract.TransferOwnership(&_VRFv2Consumer.TransactOpts, to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_VRFv2Consumer *VRFv2ConsumerTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFv2Consumer.Contract.TransferOwnership(&_VRFv2Consumer.TransactOpts, to)
}

// VRFv2ConsumerOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the VRFv2Consumer contract.
type VRFv2ConsumerOwnershipTransferRequestedIterator struct {
	Event *VRFv2ConsumerOwnershipTransferRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *VRFv2ConsumerOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFv2ConsumerOwnershipTransferRequested)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(VRFv2ConsumerOwnershipTransferRequested)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *VRFv2ConsumerOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VRFv2ConsumerOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VRFv2ConsumerOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the VRFv2Consumer contract.
type VRFv2ConsumerOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_VRFv2Consumer *VRFv2ConsumerFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFv2ConsumerOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFv2Consumer.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFv2ConsumerOwnershipTransferRequestedIterator{contract: _VRFv2Consumer.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_VRFv2Consumer *VRFv2ConsumerFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFv2ConsumerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFv2Consumer.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VRFv2ConsumerOwnershipTransferRequested)
				if err := _VRFv2Consumer.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

// ParseOwnershipTransferRequested is a log parse operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_VRFv2Consumer *VRFv2ConsumerFilterer) ParseOwnershipTransferRequested(log types.Log) (*VRFv2ConsumerOwnershipTransferRequested, error) {
	event := new(VRFv2ConsumerOwnershipTransferRequested)
	if err := _VRFv2Consumer.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VRFv2ConsumerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the VRFv2Consumer contract.
type VRFv2ConsumerOwnershipTransferredIterator struct {
	Event *VRFv2ConsumerOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *VRFv2ConsumerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFv2ConsumerOwnershipTransferred)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(VRFv2ConsumerOwnershipTransferred)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *VRFv2ConsumerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VRFv2ConsumerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VRFv2ConsumerOwnershipTransferred represents a OwnershipTransferred event raised by the VRFv2Consumer contract.
type VRFv2ConsumerOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_VRFv2Consumer *VRFv2ConsumerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFv2ConsumerOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFv2Consumer.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFv2ConsumerOwnershipTransferredIterator{contract: _VRFv2Consumer.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_VRFv2Consumer *VRFv2ConsumerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFv2ConsumerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFv2Consumer.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VRFv2ConsumerOwnershipTransferred)
				if err := _VRFv2Consumer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_VRFv2Consumer *VRFv2ConsumerFilterer) ParseOwnershipTransferred(log types.Log) (*VRFv2ConsumerOwnershipTransferred, error) {
	event := new(VRFv2ConsumerOwnershipTransferred)
	if err := _VRFv2Consumer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VRFv2ConsumerRequestFulfilledIterator is returned from FilterRequestFulfilled and is used to iterate over the raw logs and unpacked data for RequestFulfilled events raised by the VRFv2Consumer contract.
type VRFv2ConsumerRequestFulfilledIterator struct {
	Event *VRFv2ConsumerRequestFulfilled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *VRFv2ConsumerRequestFulfilledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFv2ConsumerRequestFulfilled)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(VRFv2ConsumerRequestFulfilled)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *VRFv2ConsumerRequestFulfilledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VRFv2ConsumerRequestFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VRFv2ConsumerRequestFulfilled represents a RequestFulfilled event raised by the VRFv2Consumer contract.
type VRFv2ConsumerRequestFulfilled struct {
	RequestId   *big.Int
	RandomWords []*big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterRequestFulfilled is a free log retrieval operation binding the contract event 0xfe2e2d779dba245964d4e3ef9b994be63856fd568bf7d3ca9e224755cb1bd54d.
//
// Solidity: event RequestFulfilled(uint256 requestId, uint256[] randomWords)
func (_VRFv2Consumer *VRFv2ConsumerFilterer) FilterRequestFulfilled(opts *bind.FilterOpts) (*VRFv2ConsumerRequestFulfilledIterator, error) {

	logs, sub, err := _VRFv2Consumer.contract.FilterLogs(opts, "RequestFulfilled")
	if err != nil {
		return nil, err
	}
	return &VRFv2ConsumerRequestFulfilledIterator{contract: _VRFv2Consumer.contract, event: "RequestFulfilled", logs: logs, sub: sub}, nil
}

// WatchRequestFulfilled is a free log subscription operation binding the contract event 0xfe2e2d779dba245964d4e3ef9b994be63856fd568bf7d3ca9e224755cb1bd54d.
//
// Solidity: event RequestFulfilled(uint256 requestId, uint256[] randomWords)
func (_VRFv2Consumer *VRFv2ConsumerFilterer) WatchRequestFulfilled(opts *bind.WatchOpts, sink chan<- *VRFv2ConsumerRequestFulfilled) (event.Subscription, error) {

	logs, sub, err := _VRFv2Consumer.contract.WatchLogs(opts, "RequestFulfilled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VRFv2ConsumerRequestFulfilled)
				if err := _VRFv2Consumer.contract.UnpackLog(event, "RequestFulfilled", log); err != nil {
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

// ParseRequestFulfilled is a log parse operation binding the contract event 0xfe2e2d779dba245964d4e3ef9b994be63856fd568bf7d3ca9e224755cb1bd54d.
//
// Solidity: event RequestFulfilled(uint256 requestId, uint256[] randomWords)
func (_VRFv2Consumer *VRFv2ConsumerFilterer) ParseRequestFulfilled(log types.Log) (*VRFv2ConsumerRequestFulfilled, error) {
	event := new(VRFv2ConsumerRequestFulfilled)
	if err := _VRFv2Consumer.contract.UnpackLog(event, "RequestFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VRFv2ConsumerRequestSentIterator is returned from FilterRequestSent and is used to iterate over the raw logs and unpacked data for RequestSent events raised by the VRFv2Consumer contract.
type VRFv2ConsumerRequestSentIterator struct {
	Event *VRFv2ConsumerRequestSent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *VRFv2ConsumerRequestSentIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFv2ConsumerRequestSent)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(VRFv2ConsumerRequestSent)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *VRFv2ConsumerRequestSentIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VRFv2ConsumerRequestSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VRFv2ConsumerRequestSent represents a RequestSent event raised by the VRFv2Consumer contract.
type VRFv2ConsumerRequestSent struct {
	RequestId *big.Int
	NumWords  uint32
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRequestSent is a free log retrieval operation binding the contract event 0xcc58b13ad3eab50626c6a6300b1d139cd6ebb1688a7cced9461c2f7e762665ee.
//
// Solidity: event RequestSent(uint256 requestId, uint32 numWords)
func (_VRFv2Consumer *VRFv2ConsumerFilterer) FilterRequestSent(opts *bind.FilterOpts) (*VRFv2ConsumerRequestSentIterator, error) {

	logs, sub, err := _VRFv2Consumer.contract.FilterLogs(opts, "RequestSent")
	if err != nil {
		return nil, err
	}
	return &VRFv2ConsumerRequestSentIterator{contract: _VRFv2Consumer.contract, event: "RequestSent", logs: logs, sub: sub}, nil
}

// WatchRequestSent is a free log subscription operation binding the contract event 0xcc58b13ad3eab50626c6a6300b1d139cd6ebb1688a7cced9461c2f7e762665ee.
//
// Solidity: event RequestSent(uint256 requestId, uint32 numWords)
func (_VRFv2Consumer *VRFv2ConsumerFilterer) WatchRequestSent(opts *bind.WatchOpts, sink chan<- *VRFv2ConsumerRequestSent) (event.Subscription, error) {

	logs, sub, err := _VRFv2Consumer.contract.WatchLogs(opts, "RequestSent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VRFv2ConsumerRequestSent)
				if err := _VRFv2Consumer.contract.UnpackLog(event, "RequestSent", log); err != nil {
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

// ParseRequestSent is a log parse operation binding the contract event 0xcc58b13ad3eab50626c6a6300b1d139cd6ebb1688a7cced9461c2f7e762665ee.
//
// Solidity: event RequestSent(uint256 requestId, uint32 numWords)
func (_VRFv2Consumer *VRFv2ConsumerFilterer) ParseRequestSent(log types.Log) (*VRFv2ConsumerRequestSent, error) {
	event := new(VRFv2ConsumerRequestSent)
	if err := _VRFv2Consumer.contract.UnpackLog(event, "RequestSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
