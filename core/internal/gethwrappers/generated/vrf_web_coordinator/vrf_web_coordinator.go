// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_web_coordinator

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
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated"
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

var VRFWebCoordinatorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vrfCoordinatorV2\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"linkToken\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ApiKeyAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"apiKeyHash\",\"type\":\"bytes32\"}],\"name\":\"InvalidAPISubscription\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidMaxGasLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRequestConfirmations\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"name\":\"InvalidVRFRequestId\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"}],\"name\":\"InvalidVRFSubscription\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotEnoughRequestAllowance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableFromLink\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlyCoordinatorCanFulfill\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TopUpTooSmall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TransferAndCallFailed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"apiKeyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestCap\",\"type\":\"uint256\"}],\"name\":\"ApiKeyHashRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"apiKeyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldCap\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newCap\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"topUpAmountJuels\",\"type\":\"uint256\"}],\"name\":\"RequestAllowanceUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"vrfWebRequestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"apiKeyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"vrfRequestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"VRFWebRandomnessFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"vrfWebRequestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"apiKeyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"vrfRequestId\",\"type\":\"uint256\"}],\"name\":\"VRFWebRandomnessRequested\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"COORDINATOR\",\"outputs\":[{\"internalType\":\"contractVRFCoordinatorV2\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"REGISTER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"REQUESTER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"apiKeyRegisterer\",\"type\":\"address\"}],\"name\":\"addApiKeyRegisterer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"randomnessRequester\",\"type\":\"address\"}],\"name\":\"addRandomnessRequester\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountJuels\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"apiKeyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"requestCap\",\"type\":\"uint256\"}],\"name\":\"registerApiKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"apiKeyRegisterer\",\"type\":\"address\"}],\"name\":\"removeApiKeyRegisterer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"randomnessRequester\",\"type\":\"address\"}],\"name\":\"removeRandomnessRequester\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountJuels\",\"type\":\"uint256\"}],\"name\":\"requestAllowanceFromJuels\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"apiKeyHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"vrfWebRequestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"}],\"name\":\"requestRandomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_apiRequests\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"apiKeyHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"vrfWebRequestId\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"s_apiSubscriptions\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"apiKeyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"requestCap\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requestCount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_config\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60e06040523480156200001157600080fd5b506040516200257e3803806200257e833981016040819052620000349162000271565b8133806000816200008c5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000bf57620000bf81620000f5565b50506001600355506001600160a01b0390811660805282811660a052811660c052620000ed600033620001a0565b5050620002a9565b336001600160a01b038216036200014f5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000083565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b620001ac8282620001b0565b5050565b60008281526002602090815260408083206001600160a01b038516845290915290205460ff16620001ac5760008281526002602090815260408083206001600160a01b03851684529091529020805460ff19166001179055620002103390565b6001600160a01b0316816001600160a01b0316837f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45050565b80516001600160a01b03811681146200026c57600080fd5b919050565b600080604083850312156200028557600080fd5b620002908362000254565b9150620002a06020840162000254565b90509250929050565b60805160a05160c0516122746200030a600039600081816102b70152818161107801526111a901526000818161036d015281816106250152818161077a01528181610eed01526111e701526000818161098c01526109f401526122746000f3fe608060405234801561001057600080fd5b50600436106101b95760003560e01c80635492d919116100f957806391d1485411610097578063c080bbc811610071578063c080bbc8146104fb578063d547741f1461050e578063f2fde38b14610521578063fc9ae6421461053457600080fd5b806391d148541461049a578063a217fddf146104e0578063a4c0ed36146104e857600080fd5b806379ba5097116100d357806379ba50971461044e5780638064f4de146104565780638891cb8f146104695780638da5cb5b1461047c57600080fd5b80635492d919146103c95780636f865728146103dd57806374f533171461042757600080fd5b80631fe543e31161016657806336568abe1161014057806336568abe146103555780633b2bcbf1146103685780634430db7e1461038f5780634d5d526e146103b657600080fd5b80631fe543e3146102fe578063248a9ca3146103115780632f2ff15d1461034257600080fd5b80630bf483f7116101975780630bf483f7146102635780630cc93cee146102765780631b6b6d23146102b257600080fd5b806301ffc9a7146101be578063078d4482146101e6578063088070f5146101fb575b600080fd5b6101d16101cc366004611ab3565b610547565b60405190151581526020015b60405180910390f35b6101f96101f4366004611b17565b6105e0565b005b6004546005546006546102299267ffffffffffffffff16919061ffff81169062010000900463ffffffff1684565b6040805167ffffffffffffffff9095168552602085019390935261ffff9091169183019190915263ffffffff1660608201526080016101dd565b6101f9610271366004611b6c565b610615565b61029d610284366004611bbf565b6008602052600090815260409020805460019091015482565b604080519283526020830191909152016101dd565b6102d97f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101dd565b6101f961030c366004611c7a565b610974565b61033461031f366004611bbf565b60009081526002602052604090206001015490565b6040519081526020016101dd565b6101f9610350366004611d1c565b610a2f565b6101f9610363366004611d1c565b610a5a565b6102d97f000000000000000000000000000000000000000000000000000000000000000081565b6103347f61a3517f153a09154844ed8be639dabc6e78dc22315c2d9a91f7eddf9398c00281565b6101f96103c4366004611b17565b610b09565b6103346103d7366004611bbf565b50606490565b61040c6103eb366004611bbf565b60076020526000908152604090208054600182015460029092015490919083565b604080519384526020840192909252908201526060016101dd565b6103347fd1f21ec03a6eb050fba156f5316dad461735df521fb446dd42c5a4728e9c70fe81565b6101f9610b3b565b6101f9610464366004611d4c565b610c38565b610334610477366004611d6e565b610d8d565b60005473ffffffffffffffffffffffffffffffffffffffff166102d9565b6101d16104a8366004611d1c565b600091825260026020908152604080842073ffffffffffffffffffffffffffffffffffffffff93909316845291905290205460ff1690565b610334600081565b6101f96104f6366004611da7565b610fef565b6101f9610509366004611b17565b61130e565b6101f961051c366004611d1c565b611340565b6101f961052f366004611b17565b611366565b6101f9610542366004611b17565b611377565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167f7965db0b0000000000000000000000000000000000000000000000000000000014806105da57507f01ffc9a7000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000008316145b92915050565b6105e86113a9565b6106127f61a3517f153a09154844ed8be639dabc6e78dc22315c2d9a91f7eddf9398c00282610a2f565b50565b61061d6113a9565b6000806000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663c3f909d46040518163ffffffff1660e01b8152600401608060405180830381865afa15801561068e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106b29190611e30565b93509350935093508361ffff168661ffff1610156106fc576040517ffebcaec200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8263ffffffff168563ffffffff161115610742576040517feaf856a100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040517fa47c769600000000000000000000000000000000000000000000000000000000815267ffffffffffffffff891660048201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff169063a47c769690602401600060405180830381865afa92505050801561081557506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526108129190810190611e84565b60015b61085c576040517fcb714a3800000000000000000000000000000000000000000000000000000000815267ffffffffffffffff891660048201526024015b60405180910390fd5b505050506040805160808101825267ffffffffffffffff8a16808252602082018a905261ffff891682840181905263ffffffff89166060909301839052600480547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001690921790915560058a905560068054620100009093027fffffffffffffffffffffffffffffffffffffffffffffffffffff000000000000909316909117919091179055517f1f7982a35d1886471f903c17cb5a274444ada7e41b33c4a206ab2b07b9507a5a90610962908a908a908a908a9067ffffffffffffffff949094168452602084019290925261ffff16604083015263ffffffff16606082015260800190565b60405180910390a15050505050505050565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610a21576040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000166024820152604401610853565b610a2b828261142c565b5050565b600082815260026020526040902060010154610a4b81336114f3565b610a5583836115c5565b505050565b73ffffffffffffffffffffffffffffffffffffffff81163314610aff576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602f60248201527f416363657373436f6e74726f6c3a2063616e206f6e6c792072656e6f756e636560448201527f20726f6c657320666f722073656c6600000000000000000000000000000000006064820152608401610853565b610a2b82826116b9565b610b116113a9565b6106127fd1f21ec03a6eb050fba156f5316dad461735df521fb446dd42c5a4728e9c70fe82611340565b60015473ffffffffffffffffffffffffffffffffffffffff163314610bbc576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610853565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b3360009081527fd0d9ca06c1ce74fee3fe742f93ff33111e9d81342f54d8aab16738421513a870602052604090205460ff16610cd0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f43616c6c6572206973206e6f74206120726567697374657265720000000000006044820152606401610853565b60008281526007602052604090205415610d16576040517fecd909ff00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805160608101825283815260208082018481526000838501818152878252600790935284902092518355516001830155516002909101555182907f37c03f5a6bcce23c850a1533b2f65bf446d1e1e85eef3a5ec7cdac2b3f03dbdc90610d819084815260200190565b60405180910390a25050565b60008381526007602090815260408083208151606081018352815480825260018301549482019490945260029091015491810191909152908203610e00576040517f6e92778000000000000000000000000000000000000000000000000000000000815260048101869052602401610853565b8060200151816040015110610e41576040517f926f196800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b604080516080810182526004805467ffffffffffffffff168083526005546020840181905260065461ffff811685870181905263ffffffff6201000090920482166060870181905296517f5d3b1d300000000000000000000000000000000000000000000000000000000081529485019290925260248401929092526044830152606482019390935291851660848301529060009073ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001690635d3b1d309060a4016020604051808303816000875af1158015610f36573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610f5a9190611f6b565b6000888152600760205260408120600201805492935090610f7a83611fb3565b909155505060408051808201825288815260208082018981526000858152600883528490209251835551600190920191909155815188815290810183905288917f6793f230b126f5741f1db8929b92f3c8af7d9ccb3285a0adb73f3850911ab17b910160405180910390a29695505050505050565b60026003540361105b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f5265656e7472616e637947756172643a207265656e7472616e742063616c6c006044820152606401610853565b60026003553373ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016146110cf576040517f44b0e3c300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610100811461110a576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600061111882840184611bbf565b60008181526007602052604081205491925003611164576040517f6e92778000000000000000000000000000000000000000000000000000000000815260048101829052602401610853565b606460008281526007602052604081206001018054918391906111878385611feb565b90915550506004546040805167ffffffffffffffff90921660208301526000917f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1691634000aea0917f0000000000000000000000000000000000000000000000000000000000000000918b91016040516020818303038152906040526040518463ffffffff1660e01b81526004016112379392919061207d565b6020604051808303816000875af1158015611256573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061127a91906120bb565b9050806112b3576040517f9702d1a700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b837fefd6c78135a3ce77075968f4357492249c55571cd4e4b9d3959f69632547d12d836112e08682611feb565b60408051928352602083019190915281018a905260600160405180910390a250506001600355505050505050565b6113166113a9565b6106127fd1f21ec03a6eb050fba156f5316dad461735df521fb446dd42c5a4728e9c70fe82610a2f565b60008281526002602052604090206001015461135c81336114f3565b610a5583836116b9565b61136e6113a9565b61061281611774565b61137f6113a9565b6106127f61a3517f153a09154844ed8be639dabc6e78dc22315c2d9a91f7eddf9398c00282611340565b60005473ffffffffffffffffffffffffffffffffffffffff16331461142a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610853565b565b600082815260086020908152604080832081518083019092528054808352600190910154928201929092529103611492576040517f1251720400000000000000000000000000000000000000000000000000000000815260048101849052602401610853565b6000838152600860209081526040808320838155600101929092558201518251915190917f466a0a1234bf0c72637d254d1a2db9a652b5ecb3721de618a60a3389de8dd12e916114e69190879087906120dd565b60405180910390a2505050565b600082815260026020908152604080832073ffffffffffffffffffffffffffffffffffffffff8516845290915290205460ff16610a2b5761154b8173ffffffffffffffffffffffffffffffffffffffff166014611869565b611556836020611869565b604051602001611567929190612132565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290527f08c379a0000000000000000000000000000000000000000000000000000000008252610853916004016121b3565b600082815260026020908152604080832073ffffffffffffffffffffffffffffffffffffffff8516845290915290205460ff16610a2b57600082815260026020908152604080832073ffffffffffffffffffffffffffffffffffffffff85168452909152902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016600117905561165b3390565b73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16837f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45050565b600082815260026020908152604080832073ffffffffffffffffffffffffffffffffffffffff8516845290915290205460ff1615610a2b57600082815260026020908152604080832073ffffffffffffffffffffffffffffffffffffffff8516808552925280832080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905551339285917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a45050565b3373ffffffffffffffffffffffffffffffffffffffff8216036117f3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610853565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b606060006118788360026121c6565b611883906002611feb565b67ffffffffffffffff81111561189b5761189b611bd8565b6040519080825280601f01601f1916602001820160405280156118c5576020820181803683370190505b5090507f3000000000000000000000000000000000000000000000000000000000000000816000815181106118fc576118fc612203565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053507f78000000000000000000000000000000000000000000000000000000000000008160018151811061195f5761195f612203565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350600061199b8460026121c6565b6119a6906001611feb565b90505b6001811115611a43577f303132333435363738396162636465660000000000000000000000000000000085600f16601081106119e7576119e7612203565b1a60f81b8282815181106119fd576119fd612203565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535060049490941c93611a3c81612232565b90506119a9565b508315611aac576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f537472696e67733a20686578206c656e67746820696e73756666696369656e746044820152606401610853565b9392505050565b600060208284031215611ac557600080fd5b81357fffffffff0000000000000000000000000000000000000000000000000000000081168114611aac57600080fd5b73ffffffffffffffffffffffffffffffffffffffff8116811461061257600080fd5b600060208284031215611b2957600080fd5b8135611aac81611af5565b67ffffffffffffffff8116811461061257600080fd5b61ffff8116811461061257600080fd5b63ffffffff8116811461061257600080fd5b60008060008060808587031215611b8257600080fd5b8435611b8d81611b34565b9350602085013592506040850135611ba481611b4a565b91506060850135611bb481611b5a565b939692955090935050565b600060208284031215611bd157600080fd5b5035919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715611c4e57611c4e611bd8565b604052919050565b600067ffffffffffffffff821115611c7057611c70611bd8565b5060051b60200190565b60008060408385031215611c8d57600080fd5b8235915060208084013567ffffffffffffffff811115611cac57600080fd5b8401601f81018613611cbd57600080fd5b8035611cd0611ccb82611c56565b611c07565b81815260059190911b82018301908381019088831115611cef57600080fd5b928401925b82841015611d0d57833582529284019290840190611cf4565b80955050505050509250929050565b60008060408385031215611d2f57600080fd5b823591506020830135611d4181611af5565b809150509250929050565b60008060408385031215611d5f57600080fd5b50508035926020909101359150565b600080600060608486031215611d8357600080fd5b83359250602084013591506040840135611d9c81611b5a565b809150509250925092565b60008060008060608587031215611dbd57600080fd5b8435611dc881611af5565b935060208501359250604085013567ffffffffffffffff80821115611dec57600080fd5b818701915087601f830112611e0057600080fd5b813581811115611e0f57600080fd5b886020828501011115611e2157600080fd5b95989497505060200194505050565b60008060008060808587031215611e4657600080fd5b8451611e5181611b4a565b6020860151909450611e6281611b5a565b6040860151909350611e7381611b5a565b6060860151909250611bb481611b5a565b60008060008060808587031215611e9a57600080fd5b84516bffffffffffffffffffffffff81168114611eb657600080fd5b80945050602080860151611ec981611b34565b6040870151909450611eda81611af5565b606087015190935067ffffffffffffffff811115611ef757600080fd5b8601601f81018813611f0857600080fd5b8051611f16611ccb82611c56565b81815260059190911b8201830190838101908a831115611f3557600080fd5b928401925b82841015611f5c578351611f4d81611af5565b82529284019290840190611f3a565b979a9699509497505050505050565b600060208284031215611f7d57600080fd5b5051919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203611fe457611fe4611f84565b5060010190565b60008219821115611ffe57611ffe611f84565b500190565b60005b8381101561201e578181015183820152602001612006565b8381111561202d576000848401525b50505050565b6000815180845261204b816020860160208601612003565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b73ffffffffffffffffffffffffffffffffffffffff841681528260208201526060604082015260006120b26060830184612033565b95945050505050565b6000602082840312156120cd57600080fd5b81518015158114611aac57600080fd5b6000606082018583526020858185015260606040850152818551808452608086019150828701935060005b8181101561212457845183529383019391830191600101612108565b509098975050505050505050565b7f416363657373436f6e74726f6c3a206163636f756e742000000000000000000081526000835161216a816017850160208801612003565b7f206973206d697373696e6720726f6c652000000000000000000000000000000060179184019182015283516121a7816028840160208801612003565b01602801949350505050565b602081526000611aac6020830184612033565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156121fe576121fe611f84565b500290565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60008161224157612241611f84565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff019056fea164736f6c634300080d000a",
}

var VRFWebCoordinatorABI = VRFWebCoordinatorMetaData.ABI

var VRFWebCoordinatorBin = VRFWebCoordinatorMetaData.Bin

func DeployVRFWebCoordinator(auth *bind.TransactOpts, backend bind.ContractBackend, vrfCoordinatorV2 common.Address, linkToken common.Address) (common.Address, *types.Transaction, *VRFWebCoordinator, error) {
	parsed, err := VRFWebCoordinatorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFWebCoordinatorBin), backend, vrfCoordinatorV2, linkToken)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFWebCoordinator{VRFWebCoordinatorCaller: VRFWebCoordinatorCaller{contract: contract}, VRFWebCoordinatorTransactor: VRFWebCoordinatorTransactor{contract: contract}, VRFWebCoordinatorFilterer: VRFWebCoordinatorFilterer{contract: contract}}, nil
}

type VRFWebCoordinator struct {
	address common.Address
	abi     abi.ABI
	VRFWebCoordinatorCaller
	VRFWebCoordinatorTransactor
	VRFWebCoordinatorFilterer
}

type VRFWebCoordinatorCaller struct {
	contract *bind.BoundContract
}

type VRFWebCoordinatorTransactor struct {
	contract *bind.BoundContract
}

type VRFWebCoordinatorFilterer struct {
	contract *bind.BoundContract
}

type VRFWebCoordinatorSession struct {
	Contract     *VRFWebCoordinator
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFWebCoordinatorCallerSession struct {
	Contract *VRFWebCoordinatorCaller
	CallOpts bind.CallOpts
}

type VRFWebCoordinatorTransactorSession struct {
	Contract     *VRFWebCoordinatorTransactor
	TransactOpts bind.TransactOpts
}

type VRFWebCoordinatorRaw struct {
	Contract *VRFWebCoordinator
}

type VRFWebCoordinatorCallerRaw struct {
	Contract *VRFWebCoordinatorCaller
}

type VRFWebCoordinatorTransactorRaw struct {
	Contract *VRFWebCoordinatorTransactor
}

func NewVRFWebCoordinator(address common.Address, backend bind.ContractBackend) (*VRFWebCoordinator, error) {
	abi, err := abi.JSON(strings.NewReader(VRFWebCoordinatorABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFWebCoordinator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFWebCoordinator{address: address, abi: abi, VRFWebCoordinatorCaller: VRFWebCoordinatorCaller{contract: contract}, VRFWebCoordinatorTransactor: VRFWebCoordinatorTransactor{contract: contract}, VRFWebCoordinatorFilterer: VRFWebCoordinatorFilterer{contract: contract}}, nil
}

func NewVRFWebCoordinatorCaller(address common.Address, caller bind.ContractCaller) (*VRFWebCoordinatorCaller, error) {
	contract, err := bindVRFWebCoordinator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFWebCoordinatorCaller{contract: contract}, nil
}

func NewVRFWebCoordinatorTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFWebCoordinatorTransactor, error) {
	contract, err := bindVRFWebCoordinator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFWebCoordinatorTransactor{contract: contract}, nil
}

func NewVRFWebCoordinatorFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFWebCoordinatorFilterer, error) {
	contract, err := bindVRFWebCoordinator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFWebCoordinatorFilterer{contract: contract}, nil
}

func bindVRFWebCoordinator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFWebCoordinatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_VRFWebCoordinator *VRFWebCoordinatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFWebCoordinator.Contract.VRFWebCoordinatorCaller.contract.Call(opts, result, method, params...)
}

func (_VRFWebCoordinator *VRFWebCoordinatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFWebCoordinator.Contract.VRFWebCoordinatorTransactor.contract.Transfer(opts)
}

func (_VRFWebCoordinator *VRFWebCoordinatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFWebCoordinator.Contract.VRFWebCoordinatorTransactor.contract.Transact(opts, method, params...)
}

func (_VRFWebCoordinator *VRFWebCoordinatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFWebCoordinator.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFWebCoordinator *VRFWebCoordinatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFWebCoordinator.Contract.contract.Transfer(opts)
}

func (_VRFWebCoordinator *VRFWebCoordinatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFWebCoordinator.Contract.contract.Transact(opts, method, params...)
}

func (_VRFWebCoordinator *VRFWebCoordinatorCaller) COORDINATOR(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFWebCoordinator.contract.Call(opts, &out, "COORDINATOR")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFWebCoordinator *VRFWebCoordinatorSession) COORDINATOR() (common.Address, error) {
	return _VRFWebCoordinator.Contract.COORDINATOR(&_VRFWebCoordinator.CallOpts)
}

func (_VRFWebCoordinator *VRFWebCoordinatorCallerSession) COORDINATOR() (common.Address, error) {
	return _VRFWebCoordinator.Contract.COORDINATOR(&_VRFWebCoordinator.CallOpts)
}

func (_VRFWebCoordinator *VRFWebCoordinatorCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _VRFWebCoordinator.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VRFWebCoordinator *VRFWebCoordinatorSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _VRFWebCoordinator.Contract.DEFAULTADMINROLE(&_VRFWebCoordinator.CallOpts)
}

func (_VRFWebCoordinator *VRFWebCoordinatorCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _VRFWebCoordinator.Contract.DEFAULTADMINROLE(&_VRFWebCoordinator.CallOpts)
}

func (_VRFWebCoordinator *VRFWebCoordinatorCaller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFWebCoordinator.contract.Call(opts, &out, "LINK")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFWebCoordinator *VRFWebCoordinatorSession) LINK() (common.Address, error) {
	return _VRFWebCoordinator.Contract.LINK(&_VRFWebCoordinator.CallOpts)
}

func (_VRFWebCoordinator *VRFWebCoordinatorCallerSession) LINK() (common.Address, error) {
	return _VRFWebCoordinator.Contract.LINK(&_VRFWebCoordinator.CallOpts)
}

func (_VRFWebCoordinator *VRFWebCoordinatorCaller) REGISTERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _VRFWebCoordinator.contract.Call(opts, &out, "REGISTER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VRFWebCoordinator *VRFWebCoordinatorSession) REGISTERROLE() ([32]byte, error) {
	return _VRFWebCoordinator.Contract.REGISTERROLE(&_VRFWebCoordinator.CallOpts)
}

func (_VRFWebCoordinator *VRFWebCoordinatorCallerSession) REGISTERROLE() ([32]byte, error) {
	return _VRFWebCoordinator.Contract.REGISTERROLE(&_VRFWebCoordinator.CallOpts)
}

func (_VRFWebCoordinator *VRFWebCoordinatorCaller) REQUESTERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _VRFWebCoordinator.contract.Call(opts, &out, "REQUESTER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VRFWebCoordinator *VRFWebCoordinatorSession) REQUESTERROLE() ([32]byte, error) {
	return _VRFWebCoordinator.Contract.REQUESTERROLE(&_VRFWebCoordinator.CallOpts)
}

func (_VRFWebCoordinator *VRFWebCoordinatorCallerSession) REQUESTERROLE() ([32]byte, error) {
	return _VRFWebCoordinator.Contract.REQUESTERROLE(&_VRFWebCoordinator.CallOpts)
}

func (_VRFWebCoordinator *VRFWebCoordinatorCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _VRFWebCoordinator.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VRFWebCoordinator *VRFWebCoordinatorSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _VRFWebCoordinator.Contract.GetRoleAdmin(&_VRFWebCoordinator.CallOpts, role)
}

func (_VRFWebCoordinator *VRFWebCoordinatorCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _VRFWebCoordinator.Contract.GetRoleAdmin(&_VRFWebCoordinator.CallOpts, role)
}

func (_VRFWebCoordinator *VRFWebCoordinatorCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _VRFWebCoordinator.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_VRFWebCoordinator *VRFWebCoordinatorSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _VRFWebCoordinator.Contract.HasRole(&_VRFWebCoordinator.CallOpts, role, account)
}

func (_VRFWebCoordinator *VRFWebCoordinatorCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _VRFWebCoordinator.Contract.HasRole(&_VRFWebCoordinator.CallOpts, role, account)
}

func (_VRFWebCoordinator *VRFWebCoordinatorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFWebCoordinator.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFWebCoordinator *VRFWebCoordinatorSession) Owner() (common.Address, error) {
	return _VRFWebCoordinator.Contract.Owner(&_VRFWebCoordinator.CallOpts)
}

func (_VRFWebCoordinator *VRFWebCoordinatorCallerSession) Owner() (common.Address, error) {
	return _VRFWebCoordinator.Contract.Owner(&_VRFWebCoordinator.CallOpts)
}

func (_VRFWebCoordinator *VRFWebCoordinatorCaller) SApiRequests(opts *bind.CallOpts, arg0 *big.Int) (SApiRequests,

	error) {
	var out []interface{}
	err := _VRFWebCoordinator.contract.Call(opts, &out, "s_apiRequests", arg0)

	outstruct := new(SApiRequests)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ApiKeyHash = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.VrfWebRequestId = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

func (_VRFWebCoordinator *VRFWebCoordinatorSession) SApiRequests(arg0 *big.Int) (SApiRequests,

	error) {
	return _VRFWebCoordinator.Contract.SApiRequests(&_VRFWebCoordinator.CallOpts, arg0)
}

func (_VRFWebCoordinator *VRFWebCoordinatorCallerSession) SApiRequests(arg0 *big.Int) (SApiRequests,

	error) {
	return _VRFWebCoordinator.Contract.SApiRequests(&_VRFWebCoordinator.CallOpts, arg0)
}

func (_VRFWebCoordinator *VRFWebCoordinatorCaller) SApiSubscriptions(opts *bind.CallOpts, arg0 [32]byte) (SApiSubscriptions,

	error) {
	var out []interface{}
	err := _VRFWebCoordinator.contract.Call(opts, &out, "s_apiSubscriptions", arg0)

	outstruct := new(SApiSubscriptions)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ApiKeyHash = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.RequestCap = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.RequestCount = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_VRFWebCoordinator *VRFWebCoordinatorSession) SApiSubscriptions(arg0 [32]byte) (SApiSubscriptions,

	error) {
	return _VRFWebCoordinator.Contract.SApiSubscriptions(&_VRFWebCoordinator.CallOpts, arg0)
}

func (_VRFWebCoordinator *VRFWebCoordinatorCallerSession) SApiSubscriptions(arg0 [32]byte) (SApiSubscriptions,

	error) {
	return _VRFWebCoordinator.Contract.SApiSubscriptions(&_VRFWebCoordinator.CallOpts, arg0)
}

func (_VRFWebCoordinator *VRFWebCoordinatorCaller) SConfig(opts *bind.CallOpts) (SConfig,

	error) {
	var out []interface{}
	err := _VRFWebCoordinator.contract.Call(opts, &out, "s_config")

	outstruct := new(SConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.SubscriptionId = *abi.ConvertType(out[0], new(uint64)).(*uint64)
	outstruct.KeyHash = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.RequestConfirmations = *abi.ConvertType(out[2], new(uint16)).(*uint16)
	outstruct.CallbackGasLimit = *abi.ConvertType(out[3], new(uint32)).(*uint32)

	return *outstruct, err

}

func (_VRFWebCoordinator *VRFWebCoordinatorSession) SConfig() (SConfig,

	error) {
	return _VRFWebCoordinator.Contract.SConfig(&_VRFWebCoordinator.CallOpts)
}

func (_VRFWebCoordinator *VRFWebCoordinatorCallerSession) SConfig() (SConfig,

	error) {
	return _VRFWebCoordinator.Contract.SConfig(&_VRFWebCoordinator.CallOpts)
}

func (_VRFWebCoordinator *VRFWebCoordinatorCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _VRFWebCoordinator.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_VRFWebCoordinator *VRFWebCoordinatorSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _VRFWebCoordinator.Contract.SupportsInterface(&_VRFWebCoordinator.CallOpts, interfaceId)
}

func (_VRFWebCoordinator *VRFWebCoordinatorCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _VRFWebCoordinator.Contract.SupportsInterface(&_VRFWebCoordinator.CallOpts, interfaceId)
}

func (_VRFWebCoordinator *VRFWebCoordinatorTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFWebCoordinator.contract.Transact(opts, "acceptOwnership")
}

func (_VRFWebCoordinator *VRFWebCoordinatorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFWebCoordinator.Contract.AcceptOwnership(&_VRFWebCoordinator.TransactOpts)
}

func (_VRFWebCoordinator *VRFWebCoordinatorTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFWebCoordinator.Contract.AcceptOwnership(&_VRFWebCoordinator.TransactOpts)
}

func (_VRFWebCoordinator *VRFWebCoordinatorTransactor) AddApiKeyRegisterer(opts *bind.TransactOpts, apiKeyRegisterer common.Address) (*types.Transaction, error) {
	return _VRFWebCoordinator.contract.Transact(opts, "addApiKeyRegisterer", apiKeyRegisterer)
}

func (_VRFWebCoordinator *VRFWebCoordinatorSession) AddApiKeyRegisterer(apiKeyRegisterer common.Address) (*types.Transaction, error) {
	return _VRFWebCoordinator.Contract.AddApiKeyRegisterer(&_VRFWebCoordinator.TransactOpts, apiKeyRegisterer)
}

func (_VRFWebCoordinator *VRFWebCoordinatorTransactorSession) AddApiKeyRegisterer(apiKeyRegisterer common.Address) (*types.Transaction, error) {
	return _VRFWebCoordinator.Contract.AddApiKeyRegisterer(&_VRFWebCoordinator.TransactOpts, apiKeyRegisterer)
}

func (_VRFWebCoordinator *VRFWebCoordinatorTransactor) AddRandomnessRequester(opts *bind.TransactOpts, randomnessRequester common.Address) (*types.Transaction, error) {
	return _VRFWebCoordinator.contract.Transact(opts, "addRandomnessRequester", randomnessRequester)
}

func (_VRFWebCoordinator *VRFWebCoordinatorSession) AddRandomnessRequester(randomnessRequester common.Address) (*types.Transaction, error) {
	return _VRFWebCoordinator.Contract.AddRandomnessRequester(&_VRFWebCoordinator.TransactOpts, randomnessRequester)
}

func (_VRFWebCoordinator *VRFWebCoordinatorTransactorSession) AddRandomnessRequester(randomnessRequester common.Address) (*types.Transaction, error) {
	return _VRFWebCoordinator.Contract.AddRandomnessRequester(&_VRFWebCoordinator.TransactOpts, randomnessRequester)
}

func (_VRFWebCoordinator *VRFWebCoordinatorTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _VRFWebCoordinator.contract.Transact(opts, "grantRole", role, account)
}

func (_VRFWebCoordinator *VRFWebCoordinatorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _VRFWebCoordinator.Contract.GrantRole(&_VRFWebCoordinator.TransactOpts, role, account)
}

func (_VRFWebCoordinator *VRFWebCoordinatorTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _VRFWebCoordinator.Contract.GrantRole(&_VRFWebCoordinator.TransactOpts, role, account)
}

func (_VRFWebCoordinator *VRFWebCoordinatorTransactor) OnTokenTransfer(opts *bind.TransactOpts, sender common.Address, amountJuels *big.Int, data []byte) (*types.Transaction, error) {
	return _VRFWebCoordinator.contract.Transact(opts, "onTokenTransfer", sender, amountJuels, data)
}

func (_VRFWebCoordinator *VRFWebCoordinatorSession) OnTokenTransfer(sender common.Address, amountJuels *big.Int, data []byte) (*types.Transaction, error) {
	return _VRFWebCoordinator.Contract.OnTokenTransfer(&_VRFWebCoordinator.TransactOpts, sender, amountJuels, data)
}

func (_VRFWebCoordinator *VRFWebCoordinatorTransactorSession) OnTokenTransfer(sender common.Address, amountJuels *big.Int, data []byte) (*types.Transaction, error) {
	return _VRFWebCoordinator.Contract.OnTokenTransfer(&_VRFWebCoordinator.TransactOpts, sender, amountJuels, data)
}

func (_VRFWebCoordinator *VRFWebCoordinatorTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFWebCoordinator.contract.Transact(opts, "rawFulfillRandomWords", requestId, randomWords)
}

func (_VRFWebCoordinator *VRFWebCoordinatorSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFWebCoordinator.Contract.RawFulfillRandomWords(&_VRFWebCoordinator.TransactOpts, requestId, randomWords)
}

func (_VRFWebCoordinator *VRFWebCoordinatorTransactorSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFWebCoordinator.Contract.RawFulfillRandomWords(&_VRFWebCoordinator.TransactOpts, requestId, randomWords)
}

func (_VRFWebCoordinator *VRFWebCoordinatorTransactor) RegisterApiKey(opts *bind.TransactOpts, apiKeyHash [32]byte, requestCap *big.Int) (*types.Transaction, error) {
	return _VRFWebCoordinator.contract.Transact(opts, "registerApiKey", apiKeyHash, requestCap)
}

func (_VRFWebCoordinator *VRFWebCoordinatorSession) RegisterApiKey(apiKeyHash [32]byte, requestCap *big.Int) (*types.Transaction, error) {
	return _VRFWebCoordinator.Contract.RegisterApiKey(&_VRFWebCoordinator.TransactOpts, apiKeyHash, requestCap)
}

func (_VRFWebCoordinator *VRFWebCoordinatorTransactorSession) RegisterApiKey(apiKeyHash [32]byte, requestCap *big.Int) (*types.Transaction, error) {
	return _VRFWebCoordinator.Contract.RegisterApiKey(&_VRFWebCoordinator.TransactOpts, apiKeyHash, requestCap)
}

func (_VRFWebCoordinator *VRFWebCoordinatorTransactor) RemoveApiKeyRegisterer(opts *bind.TransactOpts, apiKeyRegisterer common.Address) (*types.Transaction, error) {
	return _VRFWebCoordinator.contract.Transact(opts, "removeApiKeyRegisterer", apiKeyRegisterer)
}

func (_VRFWebCoordinator *VRFWebCoordinatorSession) RemoveApiKeyRegisterer(apiKeyRegisterer common.Address) (*types.Transaction, error) {
	return _VRFWebCoordinator.Contract.RemoveApiKeyRegisterer(&_VRFWebCoordinator.TransactOpts, apiKeyRegisterer)
}

func (_VRFWebCoordinator *VRFWebCoordinatorTransactorSession) RemoveApiKeyRegisterer(apiKeyRegisterer common.Address) (*types.Transaction, error) {
	return _VRFWebCoordinator.Contract.RemoveApiKeyRegisterer(&_VRFWebCoordinator.TransactOpts, apiKeyRegisterer)
}

func (_VRFWebCoordinator *VRFWebCoordinatorTransactor) RemoveRandomnessRequester(opts *bind.TransactOpts, randomnessRequester common.Address) (*types.Transaction, error) {
	return _VRFWebCoordinator.contract.Transact(opts, "removeRandomnessRequester", randomnessRequester)
}

func (_VRFWebCoordinator *VRFWebCoordinatorSession) RemoveRandomnessRequester(randomnessRequester common.Address) (*types.Transaction, error) {
	return _VRFWebCoordinator.Contract.RemoveRandomnessRequester(&_VRFWebCoordinator.TransactOpts, randomnessRequester)
}

func (_VRFWebCoordinator *VRFWebCoordinatorTransactorSession) RemoveRandomnessRequester(randomnessRequester common.Address) (*types.Transaction, error) {
	return _VRFWebCoordinator.Contract.RemoveRandomnessRequester(&_VRFWebCoordinator.TransactOpts, randomnessRequester)
}

func (_VRFWebCoordinator *VRFWebCoordinatorTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _VRFWebCoordinator.contract.Transact(opts, "renounceRole", role, account)
}

func (_VRFWebCoordinator *VRFWebCoordinatorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _VRFWebCoordinator.Contract.RenounceRole(&_VRFWebCoordinator.TransactOpts, role, account)
}

func (_VRFWebCoordinator *VRFWebCoordinatorTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _VRFWebCoordinator.Contract.RenounceRole(&_VRFWebCoordinator.TransactOpts, role, account)
}

func (_VRFWebCoordinator *VRFWebCoordinatorTransactor) RequestAllowanceFromJuels(opts *bind.TransactOpts, amountJuels *big.Int) (*types.Transaction, error) {
	return _VRFWebCoordinator.contract.Transact(opts, "requestAllowanceFromJuels", amountJuels)
}

func (_VRFWebCoordinator *VRFWebCoordinatorSession) RequestAllowanceFromJuels(amountJuels *big.Int) (*types.Transaction, error) {
	return _VRFWebCoordinator.Contract.RequestAllowanceFromJuels(&_VRFWebCoordinator.TransactOpts, amountJuels)
}

func (_VRFWebCoordinator *VRFWebCoordinatorTransactorSession) RequestAllowanceFromJuels(amountJuels *big.Int) (*types.Transaction, error) {
	return _VRFWebCoordinator.Contract.RequestAllowanceFromJuels(&_VRFWebCoordinator.TransactOpts, amountJuels)
}

func (_VRFWebCoordinator *VRFWebCoordinatorTransactor) RequestRandomWords(opts *bind.TransactOpts, apiKeyHash [32]byte, vrfWebRequestId [32]byte, numWords uint32) (*types.Transaction, error) {
	return _VRFWebCoordinator.contract.Transact(opts, "requestRandomWords", apiKeyHash, vrfWebRequestId, numWords)
}

func (_VRFWebCoordinator *VRFWebCoordinatorSession) RequestRandomWords(apiKeyHash [32]byte, vrfWebRequestId [32]byte, numWords uint32) (*types.Transaction, error) {
	return _VRFWebCoordinator.Contract.RequestRandomWords(&_VRFWebCoordinator.TransactOpts, apiKeyHash, vrfWebRequestId, numWords)
}

func (_VRFWebCoordinator *VRFWebCoordinatorTransactorSession) RequestRandomWords(apiKeyHash [32]byte, vrfWebRequestId [32]byte, numWords uint32) (*types.Transaction, error) {
	return _VRFWebCoordinator.Contract.RequestRandomWords(&_VRFWebCoordinator.TransactOpts, apiKeyHash, vrfWebRequestId, numWords)
}

func (_VRFWebCoordinator *VRFWebCoordinatorTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _VRFWebCoordinator.contract.Transact(opts, "revokeRole", role, account)
}

func (_VRFWebCoordinator *VRFWebCoordinatorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _VRFWebCoordinator.Contract.RevokeRole(&_VRFWebCoordinator.TransactOpts, role, account)
}

func (_VRFWebCoordinator *VRFWebCoordinatorTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _VRFWebCoordinator.Contract.RevokeRole(&_VRFWebCoordinator.TransactOpts, role, account)
}

func (_VRFWebCoordinator *VRFWebCoordinatorTransactor) SetConfig(opts *bind.TransactOpts, subscriptionId uint64, keyHash [32]byte, requestConfirmations uint16, callbackGasLimit uint32) (*types.Transaction, error) {
	return _VRFWebCoordinator.contract.Transact(opts, "setConfig", subscriptionId, keyHash, requestConfirmations, callbackGasLimit)
}

func (_VRFWebCoordinator *VRFWebCoordinatorSession) SetConfig(subscriptionId uint64, keyHash [32]byte, requestConfirmations uint16, callbackGasLimit uint32) (*types.Transaction, error) {
	return _VRFWebCoordinator.Contract.SetConfig(&_VRFWebCoordinator.TransactOpts, subscriptionId, keyHash, requestConfirmations, callbackGasLimit)
}

func (_VRFWebCoordinator *VRFWebCoordinatorTransactorSession) SetConfig(subscriptionId uint64, keyHash [32]byte, requestConfirmations uint16, callbackGasLimit uint32) (*types.Transaction, error) {
	return _VRFWebCoordinator.Contract.SetConfig(&_VRFWebCoordinator.TransactOpts, subscriptionId, keyHash, requestConfirmations, callbackGasLimit)
}

func (_VRFWebCoordinator *VRFWebCoordinatorTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFWebCoordinator.contract.Transact(opts, "transferOwnership", to)
}

func (_VRFWebCoordinator *VRFWebCoordinatorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFWebCoordinator.Contract.TransferOwnership(&_VRFWebCoordinator.TransactOpts, to)
}

func (_VRFWebCoordinator *VRFWebCoordinatorTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFWebCoordinator.Contract.TransferOwnership(&_VRFWebCoordinator.TransactOpts, to)
}

type VRFWebCoordinatorApiKeyHashRegisteredIterator struct {
	Event *VRFWebCoordinatorApiKeyHashRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFWebCoordinatorApiKeyHashRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFWebCoordinatorApiKeyHashRegistered)
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
		it.Event = new(VRFWebCoordinatorApiKeyHashRegistered)
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

func (it *VRFWebCoordinatorApiKeyHashRegisteredIterator) Error() error {
	return it.fail
}

func (it *VRFWebCoordinatorApiKeyHashRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFWebCoordinatorApiKeyHashRegistered struct {
	ApiKeyHash [32]byte
	RequestCap *big.Int
	Raw        types.Log
}

func (_VRFWebCoordinator *VRFWebCoordinatorFilterer) FilterApiKeyHashRegistered(opts *bind.FilterOpts, apiKeyHash [][32]byte) (*VRFWebCoordinatorApiKeyHashRegisteredIterator, error) {

	var apiKeyHashRule []interface{}
	for _, apiKeyHashItem := range apiKeyHash {
		apiKeyHashRule = append(apiKeyHashRule, apiKeyHashItem)
	}

	logs, sub, err := _VRFWebCoordinator.contract.FilterLogs(opts, "ApiKeyHashRegistered", apiKeyHashRule)
	if err != nil {
		return nil, err
	}
	return &VRFWebCoordinatorApiKeyHashRegisteredIterator{contract: _VRFWebCoordinator.contract, event: "ApiKeyHashRegistered", logs: logs, sub: sub}, nil
}

func (_VRFWebCoordinator *VRFWebCoordinatorFilterer) WatchApiKeyHashRegistered(opts *bind.WatchOpts, sink chan<- *VRFWebCoordinatorApiKeyHashRegistered, apiKeyHash [][32]byte) (event.Subscription, error) {

	var apiKeyHashRule []interface{}
	for _, apiKeyHashItem := range apiKeyHash {
		apiKeyHashRule = append(apiKeyHashRule, apiKeyHashItem)
	}

	logs, sub, err := _VRFWebCoordinator.contract.WatchLogs(opts, "ApiKeyHashRegistered", apiKeyHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFWebCoordinatorApiKeyHashRegistered)
				if err := _VRFWebCoordinator.contract.UnpackLog(event, "ApiKeyHashRegistered", log); err != nil {
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

func (_VRFWebCoordinator *VRFWebCoordinatorFilterer) ParseApiKeyHashRegistered(log types.Log) (*VRFWebCoordinatorApiKeyHashRegistered, error) {
	event := new(VRFWebCoordinatorApiKeyHashRegistered)
	if err := _VRFWebCoordinator.contract.UnpackLog(event, "ApiKeyHashRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFWebCoordinatorConfigSetIterator struct {
	Event *VRFWebCoordinatorConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFWebCoordinatorConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFWebCoordinatorConfigSet)
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
		it.Event = new(VRFWebCoordinatorConfigSet)
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

func (it *VRFWebCoordinatorConfigSetIterator) Error() error {
	return it.fail
}

func (it *VRFWebCoordinatorConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFWebCoordinatorConfigSet struct {
	SubscriptionId       uint64
	KeyHash              [32]byte
	RequestConfirmations uint16
	CallbackGasLimit     uint32
	Raw                  types.Log
}

func (_VRFWebCoordinator *VRFWebCoordinatorFilterer) FilterConfigSet(opts *bind.FilterOpts) (*VRFWebCoordinatorConfigSetIterator, error) {

	logs, sub, err := _VRFWebCoordinator.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &VRFWebCoordinatorConfigSetIterator{contract: _VRFWebCoordinator.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_VRFWebCoordinator *VRFWebCoordinatorFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *VRFWebCoordinatorConfigSet) (event.Subscription, error) {

	logs, sub, err := _VRFWebCoordinator.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFWebCoordinatorConfigSet)
				if err := _VRFWebCoordinator.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_VRFWebCoordinator *VRFWebCoordinatorFilterer) ParseConfigSet(log types.Log) (*VRFWebCoordinatorConfigSet, error) {
	event := new(VRFWebCoordinatorConfigSet)
	if err := _VRFWebCoordinator.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFWebCoordinatorOwnershipTransferRequestedIterator struct {
	Event *VRFWebCoordinatorOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFWebCoordinatorOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFWebCoordinatorOwnershipTransferRequested)
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
		it.Event = new(VRFWebCoordinatorOwnershipTransferRequested)
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

func (it *VRFWebCoordinatorOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFWebCoordinatorOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFWebCoordinatorOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFWebCoordinator *VRFWebCoordinatorFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFWebCoordinatorOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFWebCoordinator.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFWebCoordinatorOwnershipTransferRequestedIterator{contract: _VRFWebCoordinator.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFWebCoordinator *VRFWebCoordinatorFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFWebCoordinatorOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFWebCoordinator.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFWebCoordinatorOwnershipTransferRequested)
				if err := _VRFWebCoordinator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VRFWebCoordinator *VRFWebCoordinatorFilterer) ParseOwnershipTransferRequested(log types.Log) (*VRFWebCoordinatorOwnershipTransferRequested, error) {
	event := new(VRFWebCoordinatorOwnershipTransferRequested)
	if err := _VRFWebCoordinator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFWebCoordinatorOwnershipTransferredIterator struct {
	Event *VRFWebCoordinatorOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFWebCoordinatorOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFWebCoordinatorOwnershipTransferred)
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
		it.Event = new(VRFWebCoordinatorOwnershipTransferred)
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

func (it *VRFWebCoordinatorOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFWebCoordinatorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFWebCoordinatorOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFWebCoordinator *VRFWebCoordinatorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFWebCoordinatorOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFWebCoordinator.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFWebCoordinatorOwnershipTransferredIterator{contract: _VRFWebCoordinator.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VRFWebCoordinator *VRFWebCoordinatorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFWebCoordinatorOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFWebCoordinator.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFWebCoordinatorOwnershipTransferred)
				if err := _VRFWebCoordinator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VRFWebCoordinator *VRFWebCoordinatorFilterer) ParseOwnershipTransferred(log types.Log) (*VRFWebCoordinatorOwnershipTransferred, error) {
	event := new(VRFWebCoordinatorOwnershipTransferred)
	if err := _VRFWebCoordinator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFWebCoordinatorRequestAllowanceUpdatedIterator struct {
	Event *VRFWebCoordinatorRequestAllowanceUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFWebCoordinatorRequestAllowanceUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFWebCoordinatorRequestAllowanceUpdated)
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
		it.Event = new(VRFWebCoordinatorRequestAllowanceUpdated)
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

func (it *VRFWebCoordinatorRequestAllowanceUpdatedIterator) Error() error {
	return it.fail
}

func (it *VRFWebCoordinatorRequestAllowanceUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFWebCoordinatorRequestAllowanceUpdated struct {
	ApiKeyHash       [32]byte
	OldCap           *big.Int
	NewCap           *big.Int
	TopUpAmountJuels *big.Int
	Raw              types.Log
}

func (_VRFWebCoordinator *VRFWebCoordinatorFilterer) FilterRequestAllowanceUpdated(opts *bind.FilterOpts, apiKeyHash [][32]byte) (*VRFWebCoordinatorRequestAllowanceUpdatedIterator, error) {

	var apiKeyHashRule []interface{}
	for _, apiKeyHashItem := range apiKeyHash {
		apiKeyHashRule = append(apiKeyHashRule, apiKeyHashItem)
	}

	logs, sub, err := _VRFWebCoordinator.contract.FilterLogs(opts, "RequestAllowanceUpdated", apiKeyHashRule)
	if err != nil {
		return nil, err
	}
	return &VRFWebCoordinatorRequestAllowanceUpdatedIterator{contract: _VRFWebCoordinator.contract, event: "RequestAllowanceUpdated", logs: logs, sub: sub}, nil
}

func (_VRFWebCoordinator *VRFWebCoordinatorFilterer) WatchRequestAllowanceUpdated(opts *bind.WatchOpts, sink chan<- *VRFWebCoordinatorRequestAllowanceUpdated, apiKeyHash [][32]byte) (event.Subscription, error) {

	var apiKeyHashRule []interface{}
	for _, apiKeyHashItem := range apiKeyHash {
		apiKeyHashRule = append(apiKeyHashRule, apiKeyHashItem)
	}

	logs, sub, err := _VRFWebCoordinator.contract.WatchLogs(opts, "RequestAllowanceUpdated", apiKeyHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFWebCoordinatorRequestAllowanceUpdated)
				if err := _VRFWebCoordinator.contract.UnpackLog(event, "RequestAllowanceUpdated", log); err != nil {
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

func (_VRFWebCoordinator *VRFWebCoordinatorFilterer) ParseRequestAllowanceUpdated(log types.Log) (*VRFWebCoordinatorRequestAllowanceUpdated, error) {
	event := new(VRFWebCoordinatorRequestAllowanceUpdated)
	if err := _VRFWebCoordinator.contract.UnpackLog(event, "RequestAllowanceUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFWebCoordinatorRoleAdminChangedIterator struct {
	Event *VRFWebCoordinatorRoleAdminChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFWebCoordinatorRoleAdminChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFWebCoordinatorRoleAdminChanged)
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
		it.Event = new(VRFWebCoordinatorRoleAdminChanged)
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

func (it *VRFWebCoordinatorRoleAdminChangedIterator) Error() error {
	return it.fail
}

func (it *VRFWebCoordinatorRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFWebCoordinatorRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log
}

func (_VRFWebCoordinator *VRFWebCoordinatorFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*VRFWebCoordinatorRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _VRFWebCoordinator.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &VRFWebCoordinatorRoleAdminChangedIterator{contract: _VRFWebCoordinator.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

func (_VRFWebCoordinator *VRFWebCoordinatorFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *VRFWebCoordinatorRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _VRFWebCoordinator.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFWebCoordinatorRoleAdminChanged)
				if err := _VRFWebCoordinator.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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

func (_VRFWebCoordinator *VRFWebCoordinatorFilterer) ParseRoleAdminChanged(log types.Log) (*VRFWebCoordinatorRoleAdminChanged, error) {
	event := new(VRFWebCoordinatorRoleAdminChanged)
	if err := _VRFWebCoordinator.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFWebCoordinatorRoleGrantedIterator struct {
	Event *VRFWebCoordinatorRoleGranted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFWebCoordinatorRoleGrantedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFWebCoordinatorRoleGranted)
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
		it.Event = new(VRFWebCoordinatorRoleGranted)
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

func (it *VRFWebCoordinatorRoleGrantedIterator) Error() error {
	return it.fail
}

func (it *VRFWebCoordinatorRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFWebCoordinatorRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log
}

func (_VRFWebCoordinator *VRFWebCoordinatorFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*VRFWebCoordinatorRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _VRFWebCoordinator.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &VRFWebCoordinatorRoleGrantedIterator{contract: _VRFWebCoordinator.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

func (_VRFWebCoordinator *VRFWebCoordinatorFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *VRFWebCoordinatorRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _VRFWebCoordinator.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFWebCoordinatorRoleGranted)
				if err := _VRFWebCoordinator.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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

func (_VRFWebCoordinator *VRFWebCoordinatorFilterer) ParseRoleGranted(log types.Log) (*VRFWebCoordinatorRoleGranted, error) {
	event := new(VRFWebCoordinatorRoleGranted)
	if err := _VRFWebCoordinator.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFWebCoordinatorRoleRevokedIterator struct {
	Event *VRFWebCoordinatorRoleRevoked

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFWebCoordinatorRoleRevokedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFWebCoordinatorRoleRevoked)
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
		it.Event = new(VRFWebCoordinatorRoleRevoked)
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

func (it *VRFWebCoordinatorRoleRevokedIterator) Error() error {
	return it.fail
}

func (it *VRFWebCoordinatorRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFWebCoordinatorRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log
}

func (_VRFWebCoordinator *VRFWebCoordinatorFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*VRFWebCoordinatorRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _VRFWebCoordinator.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &VRFWebCoordinatorRoleRevokedIterator{contract: _VRFWebCoordinator.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

func (_VRFWebCoordinator *VRFWebCoordinatorFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *VRFWebCoordinatorRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _VRFWebCoordinator.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFWebCoordinatorRoleRevoked)
				if err := _VRFWebCoordinator.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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

func (_VRFWebCoordinator *VRFWebCoordinatorFilterer) ParseRoleRevoked(log types.Log) (*VRFWebCoordinatorRoleRevoked, error) {
	event := new(VRFWebCoordinatorRoleRevoked)
	if err := _VRFWebCoordinator.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFWebCoordinatorVRFWebRandomnessFulfilledIterator struct {
	Event *VRFWebCoordinatorVRFWebRandomnessFulfilled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFWebCoordinatorVRFWebRandomnessFulfilledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFWebCoordinatorVRFWebRandomnessFulfilled)
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
		it.Event = new(VRFWebCoordinatorVRFWebRandomnessFulfilled)
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

func (it *VRFWebCoordinatorVRFWebRandomnessFulfilledIterator) Error() error {
	return it.fail
}

func (it *VRFWebCoordinatorVRFWebRandomnessFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFWebCoordinatorVRFWebRandomnessFulfilled struct {
	VrfWebRequestId [32]byte
	ApiKeyHash      [32]byte
	VrfRequestId    *big.Int
	RandomWords     []*big.Int
	Raw             types.Log
}

func (_VRFWebCoordinator *VRFWebCoordinatorFilterer) FilterVRFWebRandomnessFulfilled(opts *bind.FilterOpts, apiKeyHash [][32]byte) (*VRFWebCoordinatorVRFWebRandomnessFulfilledIterator, error) {

	var apiKeyHashRule []interface{}
	for _, apiKeyHashItem := range apiKeyHash {
		apiKeyHashRule = append(apiKeyHashRule, apiKeyHashItem)
	}

	logs, sub, err := _VRFWebCoordinator.contract.FilterLogs(opts, "VRFWebRandomnessFulfilled", apiKeyHashRule)
	if err != nil {
		return nil, err
	}
	return &VRFWebCoordinatorVRFWebRandomnessFulfilledIterator{contract: _VRFWebCoordinator.contract, event: "VRFWebRandomnessFulfilled", logs: logs, sub: sub}, nil
}

func (_VRFWebCoordinator *VRFWebCoordinatorFilterer) WatchVRFWebRandomnessFulfilled(opts *bind.WatchOpts, sink chan<- *VRFWebCoordinatorVRFWebRandomnessFulfilled, apiKeyHash [][32]byte) (event.Subscription, error) {

	var apiKeyHashRule []interface{}
	for _, apiKeyHashItem := range apiKeyHash {
		apiKeyHashRule = append(apiKeyHashRule, apiKeyHashItem)
	}

	logs, sub, err := _VRFWebCoordinator.contract.WatchLogs(opts, "VRFWebRandomnessFulfilled", apiKeyHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFWebCoordinatorVRFWebRandomnessFulfilled)
				if err := _VRFWebCoordinator.contract.UnpackLog(event, "VRFWebRandomnessFulfilled", log); err != nil {
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

func (_VRFWebCoordinator *VRFWebCoordinatorFilterer) ParseVRFWebRandomnessFulfilled(log types.Log) (*VRFWebCoordinatorVRFWebRandomnessFulfilled, error) {
	event := new(VRFWebCoordinatorVRFWebRandomnessFulfilled)
	if err := _VRFWebCoordinator.contract.UnpackLog(event, "VRFWebRandomnessFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFWebCoordinatorVRFWebRandomnessRequestedIterator struct {
	Event *VRFWebCoordinatorVRFWebRandomnessRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFWebCoordinatorVRFWebRandomnessRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFWebCoordinatorVRFWebRandomnessRequested)
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
		it.Event = new(VRFWebCoordinatorVRFWebRandomnessRequested)
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

func (it *VRFWebCoordinatorVRFWebRandomnessRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFWebCoordinatorVRFWebRandomnessRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFWebCoordinatorVRFWebRandomnessRequested struct {
	VrfWebRequestId [32]byte
	ApiKeyHash      [32]byte
	VrfRequestId    *big.Int
	Raw             types.Log
}

func (_VRFWebCoordinator *VRFWebCoordinatorFilterer) FilterVRFWebRandomnessRequested(opts *bind.FilterOpts, apiKeyHash [][32]byte) (*VRFWebCoordinatorVRFWebRandomnessRequestedIterator, error) {

	var apiKeyHashRule []interface{}
	for _, apiKeyHashItem := range apiKeyHash {
		apiKeyHashRule = append(apiKeyHashRule, apiKeyHashItem)
	}

	logs, sub, err := _VRFWebCoordinator.contract.FilterLogs(opts, "VRFWebRandomnessRequested", apiKeyHashRule)
	if err != nil {
		return nil, err
	}
	return &VRFWebCoordinatorVRFWebRandomnessRequestedIterator{contract: _VRFWebCoordinator.contract, event: "VRFWebRandomnessRequested", logs: logs, sub: sub}, nil
}

func (_VRFWebCoordinator *VRFWebCoordinatorFilterer) WatchVRFWebRandomnessRequested(opts *bind.WatchOpts, sink chan<- *VRFWebCoordinatorVRFWebRandomnessRequested, apiKeyHash [][32]byte) (event.Subscription, error) {

	var apiKeyHashRule []interface{}
	for _, apiKeyHashItem := range apiKeyHash {
		apiKeyHashRule = append(apiKeyHashRule, apiKeyHashItem)
	}

	logs, sub, err := _VRFWebCoordinator.contract.WatchLogs(opts, "VRFWebRandomnessRequested", apiKeyHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFWebCoordinatorVRFWebRandomnessRequested)
				if err := _VRFWebCoordinator.contract.UnpackLog(event, "VRFWebRandomnessRequested", log); err != nil {
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

func (_VRFWebCoordinator *VRFWebCoordinatorFilterer) ParseVRFWebRandomnessRequested(log types.Log) (*VRFWebCoordinatorVRFWebRandomnessRequested, error) {
	event := new(VRFWebCoordinatorVRFWebRandomnessRequested)
	if err := _VRFWebCoordinator.contract.UnpackLog(event, "VRFWebRandomnessRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type SApiRequests struct {
	ApiKeyHash      [32]byte
	VrfWebRequestId [32]byte
}
type SApiSubscriptions struct {
	ApiKeyHash   [32]byte
	RequestCap   *big.Int
	RequestCount *big.Int
}
type SConfig struct {
	SubscriptionId       uint64
	KeyHash              [32]byte
	RequestConfirmations uint16
	CallbackGasLimit     uint32
}

func (_VRFWebCoordinator *VRFWebCoordinator) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFWebCoordinator.abi.Events["ApiKeyHashRegistered"].ID:
		return _VRFWebCoordinator.ParseApiKeyHashRegistered(log)
	case _VRFWebCoordinator.abi.Events["ConfigSet"].ID:
		return _VRFWebCoordinator.ParseConfigSet(log)
	case _VRFWebCoordinator.abi.Events["OwnershipTransferRequested"].ID:
		return _VRFWebCoordinator.ParseOwnershipTransferRequested(log)
	case _VRFWebCoordinator.abi.Events["OwnershipTransferred"].ID:
		return _VRFWebCoordinator.ParseOwnershipTransferred(log)
	case _VRFWebCoordinator.abi.Events["RequestAllowanceUpdated"].ID:
		return _VRFWebCoordinator.ParseRequestAllowanceUpdated(log)
	case _VRFWebCoordinator.abi.Events["RoleAdminChanged"].ID:
		return _VRFWebCoordinator.ParseRoleAdminChanged(log)
	case _VRFWebCoordinator.abi.Events["RoleGranted"].ID:
		return _VRFWebCoordinator.ParseRoleGranted(log)
	case _VRFWebCoordinator.abi.Events["RoleRevoked"].ID:
		return _VRFWebCoordinator.ParseRoleRevoked(log)
	case _VRFWebCoordinator.abi.Events["VRFWebRandomnessFulfilled"].ID:
		return _VRFWebCoordinator.ParseVRFWebRandomnessFulfilled(log)
	case _VRFWebCoordinator.abi.Events["VRFWebRandomnessRequested"].ID:
		return _VRFWebCoordinator.ParseVRFWebRandomnessRequested(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFWebCoordinatorApiKeyHashRegistered) Topic() common.Hash {
	return common.HexToHash("0x37c03f5a6bcce23c850a1533b2f65bf446d1e1e85eef3a5ec7cdac2b3f03dbdc")
}

func (VRFWebCoordinatorConfigSet) Topic() common.Hash {
	return common.HexToHash("0x1f7982a35d1886471f903c17cb5a274444ada7e41b33c4a206ab2b07b9507a5a")
}

func (VRFWebCoordinatorOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VRFWebCoordinatorOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (VRFWebCoordinatorRequestAllowanceUpdated) Topic() common.Hash {
	return common.HexToHash("0xefd6c78135a3ce77075968f4357492249c55571cd4e4b9d3959f69632547d12d")
}

func (VRFWebCoordinatorRoleAdminChanged) Topic() common.Hash {
	return common.HexToHash("0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff")
}

func (VRFWebCoordinatorRoleGranted) Topic() common.Hash {
	return common.HexToHash("0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d")
}

func (VRFWebCoordinatorRoleRevoked) Topic() common.Hash {
	return common.HexToHash("0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b")
}

func (VRFWebCoordinatorVRFWebRandomnessFulfilled) Topic() common.Hash {
	return common.HexToHash("0x466a0a1234bf0c72637d254d1a2db9a652b5ecb3721de618a60a3389de8dd12e")
}

func (VRFWebCoordinatorVRFWebRandomnessRequested) Topic() common.Hash {
	return common.HexToHash("0x6793f230b126f5741f1db8929b92f3c8af7d9ccb3285a0adb73f3850911ab17b")
}

func (_VRFWebCoordinator *VRFWebCoordinator) Address() common.Address {
	return _VRFWebCoordinator.address
}

type VRFWebCoordinatorInterface interface {
	COORDINATOR(opts *bind.CallOpts) (common.Address, error)

	DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error)

	LINK(opts *bind.CallOpts) (common.Address, error)

	REGISTERROLE(opts *bind.CallOpts) ([32]byte, error)

	REQUESTERROLE(opts *bind.CallOpts) ([32]byte, error)

	GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error)

	HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SApiRequests(opts *bind.CallOpts, arg0 *big.Int) (SApiRequests,

		error)

	SApiSubscriptions(opts *bind.CallOpts, arg0 [32]byte) (SApiSubscriptions,

		error)

	SConfig(opts *bind.CallOpts) (SConfig,

		error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AddApiKeyRegisterer(opts *bind.TransactOpts, apiKeyRegisterer common.Address) (*types.Transaction, error)

	AddRandomnessRequester(opts *bind.TransactOpts, randomnessRequester common.Address) (*types.Transaction, error)

	GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error)

	OnTokenTransfer(opts *bind.TransactOpts, sender common.Address, amountJuels *big.Int, data []byte) (*types.Transaction, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	RegisterApiKey(opts *bind.TransactOpts, apiKeyHash [32]byte, requestCap *big.Int) (*types.Transaction, error)

	RemoveApiKeyRegisterer(opts *bind.TransactOpts, apiKeyRegisterer common.Address) (*types.Transaction, error)

	RemoveRandomnessRequester(opts *bind.TransactOpts, randomnessRequester common.Address) (*types.Transaction, error)

	RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error)

	RequestAllowanceFromJuels(opts *bind.TransactOpts, amountJuels *big.Int) (*types.Transaction, error)

	RequestRandomWords(opts *bind.TransactOpts, apiKeyHash [32]byte, vrfWebRequestId [32]byte, numWords uint32) (*types.Transaction, error)

	RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, subscriptionId uint64, keyHash [32]byte, requestConfirmations uint16, callbackGasLimit uint32) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterApiKeyHashRegistered(opts *bind.FilterOpts, apiKeyHash [][32]byte) (*VRFWebCoordinatorApiKeyHashRegisteredIterator, error)

	WatchApiKeyHashRegistered(opts *bind.WatchOpts, sink chan<- *VRFWebCoordinatorApiKeyHashRegistered, apiKeyHash [][32]byte) (event.Subscription, error)

	ParseApiKeyHashRegistered(log types.Log) (*VRFWebCoordinatorApiKeyHashRegistered, error)

	FilterConfigSet(opts *bind.FilterOpts) (*VRFWebCoordinatorConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *VRFWebCoordinatorConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*VRFWebCoordinatorConfigSet, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFWebCoordinatorOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFWebCoordinatorOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VRFWebCoordinatorOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFWebCoordinatorOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFWebCoordinatorOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VRFWebCoordinatorOwnershipTransferred, error)

	FilterRequestAllowanceUpdated(opts *bind.FilterOpts, apiKeyHash [][32]byte) (*VRFWebCoordinatorRequestAllowanceUpdatedIterator, error)

	WatchRequestAllowanceUpdated(opts *bind.WatchOpts, sink chan<- *VRFWebCoordinatorRequestAllowanceUpdated, apiKeyHash [][32]byte) (event.Subscription, error)

	ParseRequestAllowanceUpdated(log types.Log) (*VRFWebCoordinatorRequestAllowanceUpdated, error)

	FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*VRFWebCoordinatorRoleAdminChangedIterator, error)

	WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *VRFWebCoordinatorRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error)

	ParseRoleAdminChanged(log types.Log) (*VRFWebCoordinatorRoleAdminChanged, error)

	FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*VRFWebCoordinatorRoleGrantedIterator, error)

	WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *VRFWebCoordinatorRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error)

	ParseRoleGranted(log types.Log) (*VRFWebCoordinatorRoleGranted, error)

	FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*VRFWebCoordinatorRoleRevokedIterator, error)

	WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *VRFWebCoordinatorRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error)

	ParseRoleRevoked(log types.Log) (*VRFWebCoordinatorRoleRevoked, error)

	FilterVRFWebRandomnessFulfilled(opts *bind.FilterOpts, apiKeyHash [][32]byte) (*VRFWebCoordinatorVRFWebRandomnessFulfilledIterator, error)

	WatchVRFWebRandomnessFulfilled(opts *bind.WatchOpts, sink chan<- *VRFWebCoordinatorVRFWebRandomnessFulfilled, apiKeyHash [][32]byte) (event.Subscription, error)

	ParseVRFWebRandomnessFulfilled(log types.Log) (*VRFWebCoordinatorVRFWebRandomnessFulfilled, error)

	FilterVRFWebRandomnessRequested(opts *bind.FilterOpts, apiKeyHash [][32]byte) (*VRFWebCoordinatorVRFWebRandomnessRequestedIterator, error)

	WatchVRFWebRandomnessRequested(opts *bind.WatchOpts, sink chan<- *VRFWebCoordinatorVRFWebRandomnessRequested, apiKeyHash [][32]byte) (event.Subscription, error)

	ParseVRFWebRandomnessRequested(log types.Log) (*VRFWebCoordinatorVRFWebRandomnessRequested, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
