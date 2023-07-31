// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrfv2plus_subscription_manager

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

var VRFV2PlusSubscriptionManagerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"}],\"name\":\"acceptSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"addConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"cancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createSubscription\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"fundSubscriptionWithEth\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"fundSubscriptionWithLink\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newCoordinator\",\"type\":\"address\"}],\"name\":\"migrateToNewCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"removeConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"requestSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_linkToken\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_subId\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_vrfCoordinator\",\"outputs\":[{\"internalType\":\"contractIVRFSubscriptionV2Plus\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"linkToken\",\"type\":\"address\"}],\"name\":\"setLinkToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vrfCoordinator\",\"type\":\"address\"}],\"name\":\"setVRFCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawEth\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawLink\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5033806000816100675760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615610097576100978161009f565b505050610149565b6001600160a01b0381163314156100f85760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161005e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b611a8280620001596000396000f3fe6080604052600436106101445760003560e01c806382359740116100c0578063a0ef91df11610074578063decbca6811610059578063decbca6814610395578063e41cfad61461039d578063f2fde38b146103bd57600080fd5b8063a0ef91df1461036b578063a21a23e41461038057600080fd5b80638dc654a2116100a55780638dc654a2146103095780639c24ea401461031e5780639eccacf61461033e57600080fd5b806382359740146102be5780638da5cb5b146102de57600080fd5b80633d96303511610117578063706da1ca116100fc578063706da1ca146102005780637725135b1461025757806379ba5097146102a957600080fd5b80633d963035146101c057806344ff81ce146101e057600080fd5b80630e27e3df14610149578063112940f91461016b57806324e9edb01461018b57806337ea7367146101a0575b600080fd5b34801561015557600080fd5b506101696101643660046116c2565b6103dd565b005b34801561017757600080fd5b506101696101863660046116c2565b610499565b34801561019757600080fd5b50610169610524565b3480156101ac57600080fd5b506101696101bb3660046116c2565b6105dc565b3480156101cc57600080fd5b506101696101db3660046116c2565b610ba4565b3480156101ec57600080fd5b506101696101fb3660046116c2565b610c2f565b34801561020c57600080fd5b506001546102399074010000000000000000000000000000000000000000900467ffffffffffffffff1681565b60405167ffffffffffffffff90911681526020015b60405180910390f35b34801561026357600080fd5b506003546102849073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161024e565b3480156102b557600080fd5b50610169610cfb565b3480156102ca57600080fd5b506101696102d936600461173a565b610df8565b3480156102ea57600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff16610284565b34801561031557600080fd5b50610169610ee1565b34801561032a57600080fd5b506101696103393660046116c2565b61109f565b34801561034a57600080fd5b506002546102849073ffffffffffffffffffffffffffffffffffffffff1681565b34801561037757600080fd5b5061016961116b565b34801561038c57600080fd5b50610239611228565b610169611324565b3480156103a957600080fd5b506101696103b8366004611708565b6113bd565b3480156103c957600080fd5b506101696103d83660046116c2565b61150e565b6103e561151f565b6002546001546040517f9f87fad70000000000000000000000000000000000000000000000000000000081527401000000000000000000000000000000000000000090910467ffffffffffffffff16600482015273ffffffffffffffffffffffffffffffffffffffff838116602483015290911690639f87fad7906044015b600060405180830381600087803b15801561047e57600080fd5b505af1158015610492573d6000803e3d6000fd5b5050505050565b6104a161151f565b6002546001546040517f7341c10c0000000000000000000000000000000000000000000000000000000081527401000000000000000000000000000000000000000090910467ffffffffffffffff16600482015273ffffffffffffffffffffffffffffffffffffffff838116602483015290911690637341c10c90604401610464565b61052c61151f565b6002546001546040517fd7ae1d300000000000000000000000000000000000000000000000000000000081527401000000000000000000000000000000000000000090910467ffffffffffffffff16600482015230602482015273ffffffffffffffffffffffffffffffffffffffff9091169063d7ae1d3090604401600060405180830381600087803b1580156105c257600080fd5b505af11580156105d6573d6000803e3d6000fd5b50505050565b6105e461151f565b6002546001546040517fa47c76960000000000000000000000000000000000000000000000000000000081527401000000000000000000000000000000000000000090910467ffffffffffffffff16600482015260009182918291829173ffffffffffffffffffffffffffffffffffffffff9091169063a47c76969060240160006040518083038186803b15801561067b57600080fd5b505afa15801561068f573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526106d59190810190611774565b93509350935093503073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614610777576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f4e6f74206f776e6572000000000000000000000000000000000000000000000060448201526064015b60405180910390fd5b61077f610524565b600085905060008173ffffffffffffffffffffffffffffffffffffffff1663a21a23e46040518163ffffffff1660e01b8152600401602060405180830381600087803b1580156107ce57600080fd5b505af11580156107e2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108069190611757565b90506bffffffffffffffffffffffff8616156108d8576003546040805167ffffffffffffffff8416602082015273ffffffffffffffffffffffffffffffffffffffff90921691634000aea0918a918a91016040516020818303038152906040526040518463ffffffff1660e01b81526004016108849392919061193c565b602060405180830381600087803b15801561089e57600080fd5b505af11580156108b2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108d691906116e6565b505b6bffffffffffffffffffffffff851615610986576040517f3697af8b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8216600482015273ffffffffffffffffffffffffffffffffffffffff831690633697af8b906bffffffffffffffffffffffff8816906024016000604051808303818588803b15801561096c57600080fd5b505af1158015610980573d6000803e3d6000fd5b50505050505b60005b8351811015610b05578273ffffffffffffffffffffffffffffffffffffffff16637341c10c838684815181106109c1576109c16119df565b60200260200101516040518363ffffffff1660e01b8152600401610a1192919067ffffffffffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b600060405180830381600087803b158015610a2b57600080fd5b505af1158015610a3f573d6000803e3d6000fd5b50505050838181518110610a5557610a556119df565b60209081029190910101516040517f2d6d99f300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8a8116600483015267ffffffffffffffff8516602483015290911690632d6d99f390604401600060405180830381600087803b158015610ada57600080fd5b505af1158015610aee573d6000803e3d6000fd5b505050508080610afd9061197f565b915050610989565b506001805467ffffffffffffffff90921674010000000000000000000000000000000000000000027fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff9092169190911790556002805473ffffffffffffffffffffffffffffffffffffffff9092167fffffffffffffffffffffffff00000000000000000000000000000000000000009092169190911790555050505050565b610bac61151f565b6002546001546040517f04c357cb0000000000000000000000000000000000000000000000000000000081527401000000000000000000000000000000000000000090910467ffffffffffffffff16600482015273ffffffffffffffffffffffffffffffffffffffff8381166024830152909116906304c357cb90604401610464565b610c3761151f565b73ffffffffffffffffffffffffffffffffffffffff8116610cb4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600f60248201527f496e76616c696420616464726573730000000000000000000000000000000000604482015260640161076e565b600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b60015473ffffffffffffffffffffffffffffffffffffffff163314610d7c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015260640161076e565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610e0061151f565b6002546040517f8235974000000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8316600482015273ffffffffffffffffffffffffffffffffffffffff90911690638235974090602401600060405180830381600087803b158015610e7557600080fd5b505af1158015610e89573d6000803e3d6000fd5b50506001805467ffffffffffffffff90941674010000000000000000000000000000000000000000027fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff909416939093179092555050565b610ee961151f565b6003546040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015273ffffffffffffffffffffffffffffffffffffffff9091169063a9059cbb90339083906370a082319060240160206040518083038186803b158015610f5c57600080fd5b505afa158015610f70573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610f949190611721565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e085901b16815273ffffffffffffffffffffffffffffffffffffffff90921660048301526024820152604401602060405180830381600087803b158015610fff57600080fd5b505af1158015611013573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061103791906116e6565b61109d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f556e61626c6520746f207472616e736665720000000000000000000000000000604482015260640161076e565b565b6110a761151f565b73ffffffffffffffffffffffffffffffffffffffff8116611124576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600f60248201527f496e76616c696420616464726573730000000000000000000000000000000000604482015260640161076e565b600380547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b61117361151f565b604051600090339047908381818185875af1925050503d80600081146111b5576040519150601f19603f3d011682016040523d82523d6000602084013e6111ba565b606091505b5050905080611225576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f556e61626c6520746f207472616e736665720000000000000000000000000000604482015260640161076e565b50565b600061123261151f565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a21a23e46040518163ffffffff1660e01b8152600401602060405180830381600087803b15801561129c57600080fd5b505af11580156112b0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112d49190611757565b600180547fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000067ffffffffffffffff841602179055919050565b61132c61151f565b6002546001546040517f3697af8b0000000000000000000000000000000000000000000000000000000081527401000000000000000000000000000000000000000090910467ffffffffffffffff16600482015273ffffffffffffffffffffffffffffffffffffffff90911690633697af8b9034906024016000604051808303818588803b15801561047e57600080fd5b6113c561151f565b600354600254600154604080517401000000000000000000000000000000000000000090920467ffffffffffffffff16602083015260009373ffffffffffffffffffffffffffffffffffffffff90811693634000aea0939116918691016040516020818303038152906040526040518463ffffffff1660e01b815260040161144f939291906118fe565b602060405180830381600087803b15801561146957600080fd5b505af115801561147d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906114a191906116e6565b90508061150a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600f60248201527f5472616e73666572206661696c65640000000000000000000000000000000000604482015260640161076e565b5050565b61151661151f565b611225816115a0565b60005473ffffffffffffffffffffffffffffffffffffffff16331461109d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161076e565b73ffffffffffffffffffffffffffffffffffffffff8116331415611620576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161076e565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516116a181611a3d565b919050565b80516bffffffffffffffffffffffff811681146116a157600080fd5b6000602082840312156116d457600080fd5b81356116df81611a3d565b9392505050565b6000602082840312156116f857600080fd5b815180151581146116df57600080fd5b60006020828403121561171a57600080fd5b5035919050565b60006020828403121561173357600080fd5b5051919050565b60006020828403121561174c57600080fd5b81356116df81611a5f565b60006020828403121561176957600080fd5b81516116df81611a5f565b6000806000806080858703121561178a57600080fd5b611793856116a6565b935060206117a28187016116a6565b935060408601516117b281611a3d565b606087015190935067ffffffffffffffff808211156117d057600080fd5b818801915088601f8301126117e457600080fd5b8151818111156117f6576117f6611a0e565b8060051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f8301168101818110858211171561183957611839611a0e565b604052828152858101935084860182860187018d101561185857600080fd5b600095505b838610156118825761186e81611696565b85526001959095019493860193860161185d565b50989b979a50959850505050505050565b6000815180845260005b818110156118b95760208185018101518683018201520161189d565b818111156118cb576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b73ffffffffffffffffffffffffffffffffffffffff841681528260208201526060604082015260006119336060830184611893565b95945050505050565b73ffffffffffffffffffffffffffffffffffffffff841681526bffffffffffffffffffffffff831660208201526060604082015260006119336060830184611893565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8214156119d8577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b73ffffffffffffffffffffffffffffffffffffffff8116811461122557600080fd5b67ffffffffffffffff8116811461122557600080fdfea164736f6c6343000806000a",
}

var VRFV2PlusSubscriptionManagerABI = VRFV2PlusSubscriptionManagerMetaData.ABI

var VRFV2PlusSubscriptionManagerBin = VRFV2PlusSubscriptionManagerMetaData.Bin

func DeployVRFV2PlusSubscriptionManager(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *VRFV2PlusSubscriptionManager, error) {
	parsed, err := VRFV2PlusSubscriptionManagerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFV2PlusSubscriptionManagerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFV2PlusSubscriptionManager{VRFV2PlusSubscriptionManagerCaller: VRFV2PlusSubscriptionManagerCaller{contract: contract}, VRFV2PlusSubscriptionManagerTransactor: VRFV2PlusSubscriptionManagerTransactor{contract: contract}, VRFV2PlusSubscriptionManagerFilterer: VRFV2PlusSubscriptionManagerFilterer{contract: contract}}, nil
}

type VRFV2PlusSubscriptionManager struct {
	address common.Address
	abi     abi.ABI
	VRFV2PlusSubscriptionManagerCaller
	VRFV2PlusSubscriptionManagerTransactor
	VRFV2PlusSubscriptionManagerFilterer
}

type VRFV2PlusSubscriptionManagerCaller struct {
	contract *bind.BoundContract
}

type VRFV2PlusSubscriptionManagerTransactor struct {
	contract *bind.BoundContract
}

type VRFV2PlusSubscriptionManagerFilterer struct {
	contract *bind.BoundContract
}

type VRFV2PlusSubscriptionManagerSession struct {
	Contract     *VRFV2PlusSubscriptionManager
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFV2PlusSubscriptionManagerCallerSession struct {
	Contract *VRFV2PlusSubscriptionManagerCaller
	CallOpts bind.CallOpts
}

type VRFV2PlusSubscriptionManagerTransactorSession struct {
	Contract     *VRFV2PlusSubscriptionManagerTransactor
	TransactOpts bind.TransactOpts
}

type VRFV2PlusSubscriptionManagerRaw struct {
	Contract *VRFV2PlusSubscriptionManager
}

type VRFV2PlusSubscriptionManagerCallerRaw struct {
	Contract *VRFV2PlusSubscriptionManagerCaller
}

type VRFV2PlusSubscriptionManagerTransactorRaw struct {
	Contract *VRFV2PlusSubscriptionManagerTransactor
}

func NewVRFV2PlusSubscriptionManager(address common.Address, backend bind.ContractBackend) (*VRFV2PlusSubscriptionManager, error) {
	abi, err := abi.JSON(strings.NewReader(VRFV2PlusSubscriptionManagerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFV2PlusSubscriptionManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusSubscriptionManager{address: address, abi: abi, VRFV2PlusSubscriptionManagerCaller: VRFV2PlusSubscriptionManagerCaller{contract: contract}, VRFV2PlusSubscriptionManagerTransactor: VRFV2PlusSubscriptionManagerTransactor{contract: contract}, VRFV2PlusSubscriptionManagerFilterer: VRFV2PlusSubscriptionManagerFilterer{contract: contract}}, nil
}

func NewVRFV2PlusSubscriptionManagerCaller(address common.Address, caller bind.ContractCaller) (*VRFV2PlusSubscriptionManagerCaller, error) {
	contract, err := bindVRFV2PlusSubscriptionManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusSubscriptionManagerCaller{contract: contract}, nil
}

func NewVRFV2PlusSubscriptionManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFV2PlusSubscriptionManagerTransactor, error) {
	contract, err := bindVRFV2PlusSubscriptionManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusSubscriptionManagerTransactor{contract: contract}, nil
}

func NewVRFV2PlusSubscriptionManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFV2PlusSubscriptionManagerFilterer, error) {
	contract, err := bindVRFV2PlusSubscriptionManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusSubscriptionManagerFilterer{contract: contract}, nil
}

func bindVRFV2PlusSubscriptionManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFV2PlusSubscriptionManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2PlusSubscriptionManager.Contract.VRFV2PlusSubscriptionManagerCaller.contract.Call(opts, result, method, params...)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.Contract.VRFV2PlusSubscriptionManagerTransactor.contract.Transfer(opts)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.Contract.VRFV2PlusSubscriptionManagerTransactor.contract.Transact(opts, method, params...)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2PlusSubscriptionManager.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.Contract.contract.Transfer(opts)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.Contract.contract.Transact(opts, method, params...)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2PlusSubscriptionManager.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerSession) Owner() (common.Address, error) {
	return _VRFV2PlusSubscriptionManager.Contract.Owner(&_VRFV2PlusSubscriptionManager.CallOpts)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerCallerSession) Owner() (common.Address, error) {
	return _VRFV2PlusSubscriptionManager.Contract.Owner(&_VRFV2PlusSubscriptionManager.CallOpts)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerCaller) SLinkToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2PlusSubscriptionManager.contract.Call(opts, &out, "s_linkToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerSession) SLinkToken() (common.Address, error) {
	return _VRFV2PlusSubscriptionManager.Contract.SLinkToken(&_VRFV2PlusSubscriptionManager.CallOpts)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerCallerSession) SLinkToken() (common.Address, error) {
	return _VRFV2PlusSubscriptionManager.Contract.SLinkToken(&_VRFV2PlusSubscriptionManager.CallOpts)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerCaller) SSubId(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _VRFV2PlusSubscriptionManager.contract.Call(opts, &out, "s_subId")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerSession) SSubId() (uint64, error) {
	return _VRFV2PlusSubscriptionManager.Contract.SSubId(&_VRFV2PlusSubscriptionManager.CallOpts)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerCallerSession) SSubId() (uint64, error) {
	return _VRFV2PlusSubscriptionManager.Contract.SSubId(&_VRFV2PlusSubscriptionManager.CallOpts)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerCaller) SVrfCoordinator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2PlusSubscriptionManager.contract.Call(opts, &out, "s_vrfCoordinator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerSession) SVrfCoordinator() (common.Address, error) {
	return _VRFV2PlusSubscriptionManager.Contract.SVrfCoordinator(&_VRFV2PlusSubscriptionManager.CallOpts)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerCallerSession) SVrfCoordinator() (common.Address, error) {
	return _VRFV2PlusSubscriptionManager.Contract.SVrfCoordinator(&_VRFV2PlusSubscriptionManager.CallOpts)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.contract.Transact(opts, "acceptOwnership")
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.Contract.AcceptOwnership(&_VRFV2PlusSubscriptionManager.TransactOpts)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.Contract.AcceptOwnership(&_VRFV2PlusSubscriptionManager.TransactOpts)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerTransactor) AcceptSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId uint64) (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.contract.Transact(opts, "acceptSubscriptionOwnerTransfer", subId)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerSession) AcceptSubscriptionOwnerTransfer(subId uint64) (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.Contract.AcceptSubscriptionOwnerTransfer(&_VRFV2PlusSubscriptionManager.TransactOpts, subId)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerTransactorSession) AcceptSubscriptionOwnerTransfer(subId uint64) (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.Contract.AcceptSubscriptionOwnerTransfer(&_VRFV2PlusSubscriptionManager.TransactOpts, subId)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerTransactor) AddConsumer(opts *bind.TransactOpts, consumer common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.contract.Transact(opts, "addConsumer", consumer)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerSession) AddConsumer(consumer common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.Contract.AddConsumer(&_VRFV2PlusSubscriptionManager.TransactOpts, consumer)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerTransactorSession) AddConsumer(consumer common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.Contract.AddConsumer(&_VRFV2PlusSubscriptionManager.TransactOpts, consumer)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerTransactor) CancelSubscription(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.contract.Transact(opts, "cancelSubscription")
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerSession) CancelSubscription() (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.Contract.CancelSubscription(&_VRFV2PlusSubscriptionManager.TransactOpts)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerTransactorSession) CancelSubscription() (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.Contract.CancelSubscription(&_VRFV2PlusSubscriptionManager.TransactOpts)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerTransactor) CreateSubscription(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.contract.Transact(opts, "createSubscription")
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerSession) CreateSubscription() (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.Contract.CreateSubscription(&_VRFV2PlusSubscriptionManager.TransactOpts)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerTransactorSession) CreateSubscription() (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.Contract.CreateSubscription(&_VRFV2PlusSubscriptionManager.TransactOpts)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerTransactor) FundSubscriptionWithEth(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.contract.Transact(opts, "fundSubscriptionWithEth")
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerSession) FundSubscriptionWithEth() (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.Contract.FundSubscriptionWithEth(&_VRFV2PlusSubscriptionManager.TransactOpts)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerTransactorSession) FundSubscriptionWithEth() (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.Contract.FundSubscriptionWithEth(&_VRFV2PlusSubscriptionManager.TransactOpts)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerTransactor) FundSubscriptionWithLink(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.contract.Transact(opts, "fundSubscriptionWithLink", amount)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerSession) FundSubscriptionWithLink(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.Contract.FundSubscriptionWithLink(&_VRFV2PlusSubscriptionManager.TransactOpts, amount)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerTransactorSession) FundSubscriptionWithLink(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.Contract.FundSubscriptionWithLink(&_VRFV2PlusSubscriptionManager.TransactOpts, amount)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerTransactor) MigrateToNewCoordinator(opts *bind.TransactOpts, newCoordinator common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.contract.Transact(opts, "migrateToNewCoordinator", newCoordinator)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerSession) MigrateToNewCoordinator(newCoordinator common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.Contract.MigrateToNewCoordinator(&_VRFV2PlusSubscriptionManager.TransactOpts, newCoordinator)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerTransactorSession) MigrateToNewCoordinator(newCoordinator common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.Contract.MigrateToNewCoordinator(&_VRFV2PlusSubscriptionManager.TransactOpts, newCoordinator)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerTransactor) RemoveConsumer(opts *bind.TransactOpts, consumer common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.contract.Transact(opts, "removeConsumer", consumer)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerSession) RemoveConsumer(consumer common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.Contract.RemoveConsumer(&_VRFV2PlusSubscriptionManager.TransactOpts, consumer)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerTransactorSession) RemoveConsumer(consumer common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.Contract.RemoveConsumer(&_VRFV2PlusSubscriptionManager.TransactOpts, consumer)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerTransactor) RequestSubscriptionOwnerTransfer(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.contract.Transact(opts, "requestSubscriptionOwnerTransfer", newOwner)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerSession) RequestSubscriptionOwnerTransfer(newOwner common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.Contract.RequestSubscriptionOwnerTransfer(&_VRFV2PlusSubscriptionManager.TransactOpts, newOwner)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerTransactorSession) RequestSubscriptionOwnerTransfer(newOwner common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.Contract.RequestSubscriptionOwnerTransfer(&_VRFV2PlusSubscriptionManager.TransactOpts, newOwner)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerTransactor) SetLinkToken(opts *bind.TransactOpts, linkToken common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.contract.Transact(opts, "setLinkToken", linkToken)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerSession) SetLinkToken(linkToken common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.Contract.SetLinkToken(&_VRFV2PlusSubscriptionManager.TransactOpts, linkToken)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerTransactorSession) SetLinkToken(linkToken common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.Contract.SetLinkToken(&_VRFV2PlusSubscriptionManager.TransactOpts, linkToken)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerTransactor) SetVRFCoordinator(opts *bind.TransactOpts, vrfCoordinator common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.contract.Transact(opts, "setVRFCoordinator", vrfCoordinator)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerSession) SetVRFCoordinator(vrfCoordinator common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.Contract.SetVRFCoordinator(&_VRFV2PlusSubscriptionManager.TransactOpts, vrfCoordinator)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerTransactorSession) SetVRFCoordinator(vrfCoordinator common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.Contract.SetVRFCoordinator(&_VRFV2PlusSubscriptionManager.TransactOpts, vrfCoordinator)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.contract.Transact(opts, "transferOwnership", to)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.Contract.TransferOwnership(&_VRFV2PlusSubscriptionManager.TransactOpts, to)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.Contract.TransferOwnership(&_VRFV2PlusSubscriptionManager.TransactOpts, to)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerTransactor) WithdrawEth(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.contract.Transact(opts, "withdrawEth")
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerSession) WithdrawEth() (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.Contract.WithdrawEth(&_VRFV2PlusSubscriptionManager.TransactOpts)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerTransactorSession) WithdrawEth() (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.Contract.WithdrawEth(&_VRFV2PlusSubscriptionManager.TransactOpts)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerTransactor) WithdrawLink(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.contract.Transact(opts, "withdrawLink")
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerSession) WithdrawLink() (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.Contract.WithdrawLink(&_VRFV2PlusSubscriptionManager.TransactOpts)
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerTransactorSession) WithdrawLink() (*types.Transaction, error) {
	return _VRFV2PlusSubscriptionManager.Contract.WithdrawLink(&_VRFV2PlusSubscriptionManager.TransactOpts)
}

type VRFV2PlusSubscriptionManagerOwnershipTransferRequestedIterator struct {
	Event *VRFV2PlusSubscriptionManagerOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2PlusSubscriptionManagerOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2PlusSubscriptionManagerOwnershipTransferRequested)
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
		it.Event = new(VRFV2PlusSubscriptionManagerOwnershipTransferRequested)
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

func (it *VRFV2PlusSubscriptionManagerOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFV2PlusSubscriptionManagerOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2PlusSubscriptionManagerOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusSubscriptionManagerOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusSubscriptionManager.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusSubscriptionManagerOwnershipTransferRequestedIterator{contract: _VRFV2PlusSubscriptionManager.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFV2PlusSubscriptionManagerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusSubscriptionManager.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2PlusSubscriptionManagerOwnershipTransferRequested)
				if err := _VRFV2PlusSubscriptionManager.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerFilterer) ParseOwnershipTransferRequested(log types.Log) (*VRFV2PlusSubscriptionManagerOwnershipTransferRequested, error) {
	event := new(VRFV2PlusSubscriptionManagerOwnershipTransferRequested)
	if err := _VRFV2PlusSubscriptionManager.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFV2PlusSubscriptionManagerOwnershipTransferredIterator struct {
	Event *VRFV2PlusSubscriptionManagerOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2PlusSubscriptionManagerOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2PlusSubscriptionManagerOwnershipTransferred)
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
		it.Event = new(VRFV2PlusSubscriptionManagerOwnershipTransferred)
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

func (it *VRFV2PlusSubscriptionManagerOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFV2PlusSubscriptionManagerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2PlusSubscriptionManagerOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusSubscriptionManagerOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusSubscriptionManager.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusSubscriptionManagerOwnershipTransferredIterator{contract: _VRFV2PlusSubscriptionManager.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFV2PlusSubscriptionManagerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusSubscriptionManager.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2PlusSubscriptionManagerOwnershipTransferred)
				if err := _VRFV2PlusSubscriptionManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManagerFilterer) ParseOwnershipTransferred(log types.Log) (*VRFV2PlusSubscriptionManagerOwnershipTransferred, error) {
	event := new(VRFV2PlusSubscriptionManagerOwnershipTransferred)
	if err := _VRFV2PlusSubscriptionManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManager) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFV2PlusSubscriptionManager.abi.Events["OwnershipTransferRequested"].ID:
		return _VRFV2PlusSubscriptionManager.ParseOwnershipTransferRequested(log)
	case _VRFV2PlusSubscriptionManager.abi.Events["OwnershipTransferred"].ID:
		return _VRFV2PlusSubscriptionManager.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFV2PlusSubscriptionManagerOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VRFV2PlusSubscriptionManagerOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_VRFV2PlusSubscriptionManager *VRFV2PlusSubscriptionManager) Address() common.Address {
	return _VRFV2PlusSubscriptionManager.address
}

type VRFV2PlusSubscriptionManagerInterface interface {
	Owner(opts *bind.CallOpts) (common.Address, error)

	SLinkToken(opts *bind.CallOpts) (common.Address, error)

	SSubId(opts *bind.CallOpts) (uint64, error)

	SVrfCoordinator(opts *bind.CallOpts) (common.Address, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AcceptSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId uint64) (*types.Transaction, error)

	AddConsumer(opts *bind.TransactOpts, consumer common.Address) (*types.Transaction, error)

	CancelSubscription(opts *bind.TransactOpts) (*types.Transaction, error)

	CreateSubscription(opts *bind.TransactOpts) (*types.Transaction, error)

	FundSubscriptionWithEth(opts *bind.TransactOpts) (*types.Transaction, error)

	FundSubscriptionWithLink(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	MigrateToNewCoordinator(opts *bind.TransactOpts, newCoordinator common.Address) (*types.Transaction, error)

	RemoveConsumer(opts *bind.TransactOpts, consumer common.Address) (*types.Transaction, error)

	RequestSubscriptionOwnerTransfer(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error)

	SetLinkToken(opts *bind.TransactOpts, linkToken common.Address) (*types.Transaction, error)

	SetVRFCoordinator(opts *bind.TransactOpts, vrfCoordinator common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	WithdrawEth(opts *bind.TransactOpts) (*types.Transaction, error)

	WithdrawLink(opts *bind.TransactOpts) (*types.Transaction, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusSubscriptionManagerOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFV2PlusSubscriptionManagerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VRFV2PlusSubscriptionManagerOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusSubscriptionManagerOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFV2PlusSubscriptionManagerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VRFV2PlusSubscriptionManagerOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
