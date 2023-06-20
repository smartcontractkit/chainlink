// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_subscription_api_v2plus

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

var SubscriptionAPIMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"internalBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"externalBalance\",\"type\":\"uint256\"}],\"name\":\"BalanceInvariantViolated\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedToSendEther\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"InvalidConsumer\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSubscription\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LinkAlreadySet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposedOwner\",\"type\":\"address\"}],\"name\":\"MustBeRequestedOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"MustBeSubOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableFromLink\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingRequestExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyConsumers\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"EthFundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountLink\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountEth\",\"type\":\"uint256\"}],\"name\":\"SubscriptionCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"SubscriptionCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldEthBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newEthBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFundedWithEth\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_CONSUMERS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"}],\"name\":\"acceptSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"addConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createSubscription\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"}],\"name\":\"fundSubscriptionWithEth\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"}],\"name\":\"getSubscription\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"ethBalance\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"oracleWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"oracleWithdrawEth\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"}],\"name\":\"ownerCancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"recoverEthFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"requestSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_currentSubId\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_totalBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_totalEthBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"}],\"name\":\"setLINK\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b503380600081620000695760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156200009c576200009c81620000aa565b505060016002555062000156565b6001600160a01b038116331415620001055760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000060565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b61298b80620001666000396000f3fe60806040526004361061016a5760003560e01c806386fe91c7116100cb578063a47c76961161007f578063e72f6e3011610059578063e72f6e3014610430578063e95704bd14610450578063f2fde38b1461048d57600080fd5b8063a47c7696146103c0578063a4c0ed36146103f0578063a8cb447b1461041057600080fd5b80638da5cb5b116100b05780638da5cb5b14610360578063a02e06161461038b578063a21a23e4146103ab57600080fd5b806386fe91c7146102d85780638741e3ee1461032657600080fd5b806364d51a2a116101225780637341c10c116101075780637341c10c1461028357806379ba5097146102a357806382359740146102b857600080fd5b806364d51a2a1461023b57806366316d8d1461026357600080fd5b80631b6b6d23116101535780631b6b6d23146101b15780633697af8b1461020857806346d8d4861461021b57600080fd5b806302bcc5b61461016f57806304c357cb14610191575b600080fd5b34801561017b57600080fd5b5061018f61018a366004612736565b6104ad565b005b34801561019d57600080fd5b5061018f6101ac366004612751565b610559565b3480156101bd57600080fd5b506003546101de9073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b61018f610216366004612736565b61077c565b34801561022757600080fd5b5061018f61023636600461263d565b61097a565b34801561024757600080fd5b50610250606481565b60405161ffff90911681526020016101ff565b34801561026f57600080fd5b5061018f61027e36600461263d565b610ba8565b34801561028f57600080fd5b5061018f61029e366004612751565b610e07565b3480156102af57600080fd5b5061018f6110bc565b3480156102c457600080fd5b5061018f6102d3366004612736565b6111b9565b3480156102e457600080fd5b50600754610309906801000000000000000090046bffffffffffffffffffffffff1681565b6040516bffffffffffffffffffffffff90911681526020016101ff565b34801561033257600080fd5b506007546103479067ffffffffffffffff1681565b60405167ffffffffffffffff90911681526020016101ff565b34801561036c57600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff166101de565b34801561039757600080fd5b5061018f6103a6366004612619565b6113d5565b3480156103b757600080fd5b50610347611474565b3480156103cc57600080fd5b506103e06103db366004612736565b6116b2565b6040516101ff9493929190612788565b3480156103fc57600080fd5b5061018f61040b366004612672565b6117f9565b34801561041c57600080fd5b5061018f61042b366004612619565b611a85565b34801561043c57600080fd5b5061018f61044b366004612619565b611c03565b34801561045c57600080fd5b50600754610309907401000000000000000000000000000000000000000090046bffffffffffffffffffffffff1681565b34801561049957600080fd5b5061018f6104a8366004612619565b611e29565b6104b5611e3a565b67ffffffffffffffff811660009081526005602052604090205473ffffffffffffffffffffffffffffffffffffffff1661051b576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff811660009081526005602052604090205461055690829073ffffffffffffffffffffffffffffffffffffffff16611ebd565b50565b67ffffffffffffffff8216600090815260056020526040902054829073ffffffffffffffffffffffffffffffffffffffff16806105c2576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff82161461062e576040517fd8a3fb5200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821660048201526024015b60405180910390fd5b60028054141561069a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f5265656e7472616e637947756172643a207265656e7472616e742063616c6c006044820152606401610625565b6002805567ffffffffffffffff841660009081526005602052604090206001015473ffffffffffffffffffffffffffffffffffffffff8481169116146107715767ffffffffffffffff841660008181526005602090815260409182902060010180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff88169081179091558251338152918201527f69436ea6df009049404f564eff6622cd00522b0bd6a89efd9e52a355c4a879be91015b60405180910390a25b505060016002555050565b6002805414156107e8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f5265656e7472616e637947756172643a207265656e7472616e742063616c6c006044820152606401610625565b6002805567ffffffffffffffff811660009081526005602052604090205473ffffffffffffffffffffffffffffffffffffffff16610852576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff8116600090815260066020526040902080546c0100000000000000000000000090046bffffffffffffffffffffffff16903490600c61089a8385612829565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555034600760148282829054906101000a90046bffffffffffffffffffffffff166108f19190612829565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055508167ffffffffffffffff167f4e4421a74ab9b76003de9d9527153d002f1338c36c1bbc180069671631cf09588234846109589190612811565b604080519283526020830191909152015b60405180910390a250506001600255565b6002805414156109e6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f5265656e7472616e637947756172643a207265656e7472616e742063616c6c006044820152606401610625565b60028055336000908152600960205260409020546bffffffffffffffffffffffff80831691161015610a44576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3360009081526009602052604081208054839290610a719084906bffffffffffffffffffffffff16612870565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555080600760148282829054906101000a90046bffffffffffffffffffffffff16610ac89190612870565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060008273ffffffffffffffffffffffffffffffffffffffff16826bffffffffffffffffffffffff1660405160006040518083038185875af1925050503d8060008114610b5e576040519150601f19603f3d011682016040523d82523d6000602084013e610b63565b606091505b5050905080610b9e576040517fdcf35db000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050600160025550565b600280541415610c14576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f5265656e7472616e637947756172643a207265656e7472616e742063616c6c006044820152606401610625565b60028055336000908152600860205260409020546bffffffffffffffffffffffff80831691161015610c72576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3360009081526008602052604081208054839290610c9f9084906bffffffffffffffffffffffff16612870565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555080600760088282829054906101000a90046bffffffffffffffffffffffff16610cf69190612870565b82546101009290920a6bffffffffffffffffffffffff8181021990931691831602179091556003546040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff868116600483015292851660248201529116915063a9059cbb90604401602060405180830381600087803b158015610d9057600080fd5b505af1158015610da4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610dc891906126fb565b610dfe576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b50506001600255565b67ffffffffffffffff8216600090815260056020526040902054829073ffffffffffffffffffffffffffffffffffffffff1680610e70576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff821614610ed7576040517fd8a3fb5200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82166004820152602401610625565b600280541415610f43576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f5265656e7472616e637947756172643a207265656e7472616e742063616c6c006044820152606401610625565b600280805567ffffffffffffffff85166000908152600560205260409020015460641415610f9d576040517f05a48e0f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8316600090815260046020908152604080832067ffffffffffffffff80891685529252909120541615610fe457610771565b73ffffffffffffffffffffffffffffffffffffffff8316600081815260046020908152604080832067ffffffffffffffff891680855290835281842080547fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000166001908117909155600584528285206002018054918201815585529383902090930180547fffffffffffffffffffffffff000000000000000000000000000000000000000016851790555192835290917f43dc749a04ac8fb825cbd514f7c0e13f13bc6f2ee66043b76629d51776cff8e09101610768565b60015473ffffffffffffffffffffffffffffffffffffffff16331461113d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610625565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b600280541415611225576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f5265656e7472616e637947756172643a207265656e7472616e742063616c6c006044820152606401610625565b6002805567ffffffffffffffff811660009081526005602052604090205473ffffffffffffffffffffffffffffffffffffffff1661128f576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff811660009081526005602052604090206001015473ffffffffffffffffffffffffffffffffffffffff1633146113315767ffffffffffffffff8116600090815260056020526040908190206001015490517fd084e97500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602401610625565b67ffffffffffffffff81166000818152600560209081526040918290208054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560019093018054909316909255835173ffffffffffffffffffffffffffffffffffffffff909116808252928101919091529092917f6f1dc65165ffffedfd8e507b4a0f1fcfdada045ed11f6c26ba27cedfe87802f09101610969565b6113dd611e3a565b60035473ffffffffffffffffffffffffffffffffffffffff161561142d576040517f2d118a6e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600380547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b60006002805414156114e2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f5265656e7472616e637947756172643a207265656e7472616e742063616c6c006044820152606401610625565b600280556007805467ffffffffffffffff16906000611500836128d6565b82546101009290920a67ffffffffffffffff818102199093169183160217909155600754169050600080604051908082528060200260200182016040528015611553578160200160208202803683370190505b506040805180820182526000808252602080830182815267ffffffffffffffff8816808452600683528584209451855492516bffffffffffffffffffffffff9081166c01000000000000000000000000027fffffffffffffffff00000000000000000000000000000000000000000000000090941691161791909117909355835160608101855233815280820183815281860187815294845260058352949092208251815473ffffffffffffffffffffffffffffffffffffffff9182167fffffffffffffffffffffffff00000000000000000000000000000000000000009182161783559551600183018054919092169616959095179094559151805194955090936116659260028501920190612527565b505060405133815267ffffffffffffffff841691507f464722b4166576d3dcbba877b999bc35cf911f4eaf434b7eba68fa113951d0bf9060200160405180910390a2509050600160025590565b67ffffffffffffffff81166000908152600560205260408120548190819060609073ffffffffffffffffffffffffffffffffffffffff1661171f576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff851660009081526006602090815260408083205460058352928190208054600290910180548351818602810186019094528084526bffffffffffffffffffffffff808716966c010000000000000000000000009004169473ffffffffffffffffffffffffffffffffffffffff9093169391928391908301828280156117e357602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116117b8575b5050505050905093509350935093509193509193565b600280541415611865576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f5265656e7472616e637947756172643a207265656e7472616e742063616c6c006044820152606401610625565b6002805560035473ffffffffffffffffffffffffffffffffffffffff1633146118ba576040517f44b0e3c300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b602081146118f4576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600061190282840184612736565b67ffffffffffffffff811660009081526005602052604090205490915073ffffffffffffffffffffffffffffffffffffffff1661196b576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff8116600090815260066020526040812080546bffffffffffffffffffffffff16918691906119a28385612829565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555084600760088282829054906101000a90046bffffffffffffffffffffffff166119f99190612829565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055508167ffffffffffffffff167fd39ec07f4e209f627a4c427971473820dc129761ba28de8906bd56f57101d4f8828784611a609190612811565b6040805192835260208301919091520160405180910390a25050600160025550505050565b611a8d611e3a565b60075447907401000000000000000000000000000000000000000090046bffffffffffffffffffffffff1681811115611afc576040517fa99da3020000000000000000000000000000000000000000000000000000000081526004810182905260248101839052604401610625565b81811015611bfe576000611b108284612859565b905060008473ffffffffffffffffffffffffffffffffffffffff168260405160006040518083038185875af1925050503d8060008114611b6c576040519150601f19603f3d011682016040523d82523d6000602084013e611b71565b606091505b5050905080611bac576040517fdcf35db000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805173ffffffffffffffffffffffffffffffffffffffff87168152602081018490527f879c9ea2b9d5345b84ccd12610b032602808517cebdb795007f3dcb4df377317910160405180910390a150505b505050565b611c0b611e3a565b6003546040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015260009173ffffffffffffffffffffffffffffffffffffffff16906370a082319060240160206040518083038186803b158015611c7557600080fd5b505afa158015611c89573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611cad919061271d565b6007549091506801000000000000000090046bffffffffffffffffffffffff1681811115611d11576040517fa99da3020000000000000000000000000000000000000000000000000000000081526004810182905260248101839052604401610625565b81811015611bfe576000611d258284612859565b6003546040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff87811660048301526024820184905292935091169063a9059cbb90604401602060405180830381600087803b158015611d9b57600080fd5b505af1158015611daf573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611dd391906126fb565b506040805173ffffffffffffffffffffffffffffffffffffffff86168152602081018390527f59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b436600910160405180910390a150505050565b611e31611e3a565b61055681612431565b60005473ffffffffffffffffffffffffffffffffffffffff163314611ebb576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610625565b565b600280541415611f29576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f5265656e7472616e637947756172643a207265656e7472616e742063616c6c006044820152606401610625565b600280805567ffffffffffffffff831660009081526005602090815260408083208151606081018352815473ffffffffffffffffffffffffffffffffffffffff90811682526001830154168185015294810180548351818602810186018552818152959695929493860193830182828015611fda57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311611faf575b5050509190925250505067ffffffffffffffff841660009081526006602090815260408083208151808301909252546bffffffffffffffffffffffff8082168084526c010000000000000000000000009092041692820183905293945092915b8460400151518110156120e05760046000866040015183815181106120615761206161292d565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040908101600090812067ffffffffffffffff8b168252909252902080547fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000169055806120d88161289d565b91505061203a565b5067ffffffffffffffff8616600090815260056020526040812080547fffffffffffffffffffffffff0000000000000000000000000000000000000000908116825560018201805490911690559061213b60028301826125b1565b505067ffffffffffffffff8616600090815260066020526040902080547fffffffffffffffff000000000000000000000000000000000000000000000000169055600780548391906008906121ab9084906801000000000000000090046bffffffffffffffffffffffff16612870565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555080600760148282829054906101000a90046bffffffffffffffffffffffff166122029190612870565b82546101009290920a6bffffffffffffffffffffffff8181021990931691831602179091556003546040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff898116600483015292861660248201529116915063a9059cbb90604401602060405180830381600087803b15801561229c57600080fd5b505af11580156122b0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906122d491906126fb565b61230a576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008573ffffffffffffffffffffffffffffffffffffffff16826bffffffffffffffffffffffff1660405160006040518083038185875af1925050503d8060008114612372576040519150601f19603f3d011682016040523d82523d6000602084013e612377565b606091505b50509050806123b2576040517fdcf35db000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805173ffffffffffffffffffffffffffffffffffffffff881681526bffffffffffffffffffffffff8581166020830152841681830152905167ffffffffffffffff8916917fb3a76e2c52bb7a484e325118e4cadc4346fdaa16a7df0621c5525e1f7ab4d495919081900360600190a2505060016002555050505050565b73ffffffffffffffffffffffffffffffffffffffff81163314156124b1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610625565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b8280548282559060005260206000209081019282156125a1579160200282015b828111156125a157825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190612547565b506125ad9291506125cb565b5090565b508054600082559060005260206000209081019061055691905b5b808211156125ad57600081556001016125cc565b803567ffffffffffffffff811681146125f857600080fd5b919050565b80356bffffffffffffffffffffffff811681146125f857600080fd5b60006020828403121561262b57600080fd5b81356126368161295c565b9392505050565b6000806040838503121561265057600080fd5b823561265b8161295c565b9150612669602084016125fd565b90509250929050565b6000806000806060858703121561268857600080fd5b84356126938161295c565b935060208501359250604085013567ffffffffffffffff808211156126b757600080fd5b818701915087601f8301126126cb57600080fd5b8135818111156126da57600080fd5b8860208285010111156126ec57600080fd5b95989497505060200194505050565b60006020828403121561270d57600080fd5b8151801515811461263657600080fd5b60006020828403121561272f57600080fd5b5051919050565b60006020828403121561274857600080fd5b612636826125e0565b6000806040838503121561276457600080fd5b61276d836125e0565b9150602083013561277d8161295c565b809150509250929050565b6000608082016bffffffffffffffffffffffff808816845260208188168186015273ffffffffffffffffffffffffffffffffffffffff915081871660408601526080606086015282865180855260a087019150828801945060005b818110156128015785518516835294830194918301916001016127e3565b50909a9950505050505050505050565b60008219821115612824576128246128fe565b500190565b60006bffffffffffffffffffffffff808316818516808303821115612850576128506128fe565b01949350505050565b60008282101561286b5761286b6128fe565b500390565b60006bffffffffffffffffffffffff83811690831681811015612895576128956128fe565b039392505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8214156128cf576128cf6128fe565b5060010190565b600067ffffffffffffffff808316818114156128f4576128f46128fe565b6001019392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b73ffffffffffffffffffffffffffffffffffffffff8116811461055657600080fdfea164736f6c6343000806000a",
}

var SubscriptionAPIABI = SubscriptionAPIMetaData.ABI

var SubscriptionAPIBin = SubscriptionAPIMetaData.Bin

func DeploySubscriptionAPI(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SubscriptionAPI, error) {
	parsed, err := SubscriptionAPIMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SubscriptionAPIBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SubscriptionAPI{SubscriptionAPICaller: SubscriptionAPICaller{contract: contract}, SubscriptionAPITransactor: SubscriptionAPITransactor{contract: contract}, SubscriptionAPIFilterer: SubscriptionAPIFilterer{contract: contract}}, nil
}

type SubscriptionAPI struct {
	address common.Address
	abi     abi.ABI
	SubscriptionAPICaller
	SubscriptionAPITransactor
	SubscriptionAPIFilterer
}

type SubscriptionAPICaller struct {
	contract *bind.BoundContract
}

type SubscriptionAPITransactor struct {
	contract *bind.BoundContract
}

type SubscriptionAPIFilterer struct {
	contract *bind.BoundContract
}

type SubscriptionAPISession struct {
	Contract     *SubscriptionAPI
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type SubscriptionAPICallerSession struct {
	Contract *SubscriptionAPICaller
	CallOpts bind.CallOpts
}

type SubscriptionAPITransactorSession struct {
	Contract     *SubscriptionAPITransactor
	TransactOpts bind.TransactOpts
}

type SubscriptionAPIRaw struct {
	Contract *SubscriptionAPI
}

type SubscriptionAPICallerRaw struct {
	Contract *SubscriptionAPICaller
}

type SubscriptionAPITransactorRaw struct {
	Contract *SubscriptionAPITransactor
}

func NewSubscriptionAPI(address common.Address, backend bind.ContractBackend) (*SubscriptionAPI, error) {
	abi, err := abi.JSON(strings.NewReader(SubscriptionAPIABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindSubscriptionAPI(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SubscriptionAPI{address: address, abi: abi, SubscriptionAPICaller: SubscriptionAPICaller{contract: contract}, SubscriptionAPITransactor: SubscriptionAPITransactor{contract: contract}, SubscriptionAPIFilterer: SubscriptionAPIFilterer{contract: contract}}, nil
}

func NewSubscriptionAPICaller(address common.Address, caller bind.ContractCaller) (*SubscriptionAPICaller, error) {
	contract, err := bindSubscriptionAPI(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SubscriptionAPICaller{contract: contract}, nil
}

func NewSubscriptionAPITransactor(address common.Address, transactor bind.ContractTransactor) (*SubscriptionAPITransactor, error) {
	contract, err := bindSubscriptionAPI(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SubscriptionAPITransactor{contract: contract}, nil
}

func NewSubscriptionAPIFilterer(address common.Address, filterer bind.ContractFilterer) (*SubscriptionAPIFilterer, error) {
	contract, err := bindSubscriptionAPI(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SubscriptionAPIFilterer{contract: contract}, nil
}

func bindSubscriptionAPI(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SubscriptionAPIMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_SubscriptionAPI *SubscriptionAPIRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SubscriptionAPI.Contract.SubscriptionAPICaller.contract.Call(opts, result, method, params...)
}

func (_SubscriptionAPI *SubscriptionAPIRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SubscriptionAPI.Contract.SubscriptionAPITransactor.contract.Transfer(opts)
}

func (_SubscriptionAPI *SubscriptionAPIRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SubscriptionAPI.Contract.SubscriptionAPITransactor.contract.Transact(opts, method, params...)
}

func (_SubscriptionAPI *SubscriptionAPICallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SubscriptionAPI.Contract.contract.Call(opts, result, method, params...)
}

func (_SubscriptionAPI *SubscriptionAPITransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SubscriptionAPI.Contract.contract.Transfer(opts)
}

func (_SubscriptionAPI *SubscriptionAPITransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SubscriptionAPI.Contract.contract.Transact(opts, method, params...)
}

func (_SubscriptionAPI *SubscriptionAPICaller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SubscriptionAPI.contract.Call(opts, &out, "LINK")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_SubscriptionAPI *SubscriptionAPISession) LINK() (common.Address, error) {
	return _SubscriptionAPI.Contract.LINK(&_SubscriptionAPI.CallOpts)
}

func (_SubscriptionAPI *SubscriptionAPICallerSession) LINK() (common.Address, error) {
	return _SubscriptionAPI.Contract.LINK(&_SubscriptionAPI.CallOpts)
}

func (_SubscriptionAPI *SubscriptionAPICaller) MAXCONSUMERS(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _SubscriptionAPI.contract.Call(opts, &out, "MAX_CONSUMERS")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_SubscriptionAPI *SubscriptionAPISession) MAXCONSUMERS() (uint16, error) {
	return _SubscriptionAPI.Contract.MAXCONSUMERS(&_SubscriptionAPI.CallOpts)
}

func (_SubscriptionAPI *SubscriptionAPICallerSession) MAXCONSUMERS() (uint16, error) {
	return _SubscriptionAPI.Contract.MAXCONSUMERS(&_SubscriptionAPI.CallOpts)
}

func (_SubscriptionAPI *SubscriptionAPICaller) GetSubscription(opts *bind.CallOpts, subId uint64) (GetSubscription,

	error) {
	var out []interface{}
	err := _SubscriptionAPI.contract.Call(opts, &out, "getSubscription", subId)

	outstruct := new(GetSubscription)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Balance = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.EthBalance = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Owner = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.Consumers = *abi.ConvertType(out[3], new([]common.Address)).(*[]common.Address)

	return *outstruct, err

}

func (_SubscriptionAPI *SubscriptionAPISession) GetSubscription(subId uint64) (GetSubscription,

	error) {
	return _SubscriptionAPI.Contract.GetSubscription(&_SubscriptionAPI.CallOpts, subId)
}

func (_SubscriptionAPI *SubscriptionAPICallerSession) GetSubscription(subId uint64) (GetSubscription,

	error) {
	return _SubscriptionAPI.Contract.GetSubscription(&_SubscriptionAPI.CallOpts, subId)
}

func (_SubscriptionAPI *SubscriptionAPICaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SubscriptionAPI.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_SubscriptionAPI *SubscriptionAPISession) Owner() (common.Address, error) {
	return _SubscriptionAPI.Contract.Owner(&_SubscriptionAPI.CallOpts)
}

func (_SubscriptionAPI *SubscriptionAPICallerSession) Owner() (common.Address, error) {
	return _SubscriptionAPI.Contract.Owner(&_SubscriptionAPI.CallOpts)
}

func (_SubscriptionAPI *SubscriptionAPICaller) SCurrentSubId(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _SubscriptionAPI.contract.Call(opts, &out, "s_currentSubId")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_SubscriptionAPI *SubscriptionAPISession) SCurrentSubId() (uint64, error) {
	return _SubscriptionAPI.Contract.SCurrentSubId(&_SubscriptionAPI.CallOpts)
}

func (_SubscriptionAPI *SubscriptionAPICallerSession) SCurrentSubId() (uint64, error) {
	return _SubscriptionAPI.Contract.SCurrentSubId(&_SubscriptionAPI.CallOpts)
}

func (_SubscriptionAPI *SubscriptionAPICaller) STotalBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SubscriptionAPI.contract.Call(opts, &out, "s_totalBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_SubscriptionAPI *SubscriptionAPISession) STotalBalance() (*big.Int, error) {
	return _SubscriptionAPI.Contract.STotalBalance(&_SubscriptionAPI.CallOpts)
}

func (_SubscriptionAPI *SubscriptionAPICallerSession) STotalBalance() (*big.Int, error) {
	return _SubscriptionAPI.Contract.STotalBalance(&_SubscriptionAPI.CallOpts)
}

func (_SubscriptionAPI *SubscriptionAPICaller) STotalEthBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SubscriptionAPI.contract.Call(opts, &out, "s_totalEthBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_SubscriptionAPI *SubscriptionAPISession) STotalEthBalance() (*big.Int, error) {
	return _SubscriptionAPI.Contract.STotalEthBalance(&_SubscriptionAPI.CallOpts)
}

func (_SubscriptionAPI *SubscriptionAPICallerSession) STotalEthBalance() (*big.Int, error) {
	return _SubscriptionAPI.Contract.STotalEthBalance(&_SubscriptionAPI.CallOpts)
}

func (_SubscriptionAPI *SubscriptionAPITransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SubscriptionAPI.contract.Transact(opts, "acceptOwnership")
}

func (_SubscriptionAPI *SubscriptionAPISession) AcceptOwnership() (*types.Transaction, error) {
	return _SubscriptionAPI.Contract.AcceptOwnership(&_SubscriptionAPI.TransactOpts)
}

func (_SubscriptionAPI *SubscriptionAPITransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _SubscriptionAPI.Contract.AcceptOwnership(&_SubscriptionAPI.TransactOpts)
}

func (_SubscriptionAPI *SubscriptionAPITransactor) AcceptSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId uint64) (*types.Transaction, error) {
	return _SubscriptionAPI.contract.Transact(opts, "acceptSubscriptionOwnerTransfer", subId)
}

func (_SubscriptionAPI *SubscriptionAPISession) AcceptSubscriptionOwnerTransfer(subId uint64) (*types.Transaction, error) {
	return _SubscriptionAPI.Contract.AcceptSubscriptionOwnerTransfer(&_SubscriptionAPI.TransactOpts, subId)
}

func (_SubscriptionAPI *SubscriptionAPITransactorSession) AcceptSubscriptionOwnerTransfer(subId uint64) (*types.Transaction, error) {
	return _SubscriptionAPI.Contract.AcceptSubscriptionOwnerTransfer(&_SubscriptionAPI.TransactOpts, subId)
}

func (_SubscriptionAPI *SubscriptionAPITransactor) AddConsumer(opts *bind.TransactOpts, subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _SubscriptionAPI.contract.Transact(opts, "addConsumer", subId, consumer)
}

func (_SubscriptionAPI *SubscriptionAPISession) AddConsumer(subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _SubscriptionAPI.Contract.AddConsumer(&_SubscriptionAPI.TransactOpts, subId, consumer)
}

func (_SubscriptionAPI *SubscriptionAPITransactorSession) AddConsumer(subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _SubscriptionAPI.Contract.AddConsumer(&_SubscriptionAPI.TransactOpts, subId, consumer)
}

func (_SubscriptionAPI *SubscriptionAPITransactor) CreateSubscription(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SubscriptionAPI.contract.Transact(opts, "createSubscription")
}

func (_SubscriptionAPI *SubscriptionAPISession) CreateSubscription() (*types.Transaction, error) {
	return _SubscriptionAPI.Contract.CreateSubscription(&_SubscriptionAPI.TransactOpts)
}

func (_SubscriptionAPI *SubscriptionAPITransactorSession) CreateSubscription() (*types.Transaction, error) {
	return _SubscriptionAPI.Contract.CreateSubscription(&_SubscriptionAPI.TransactOpts)
}

func (_SubscriptionAPI *SubscriptionAPITransactor) FundSubscriptionWithEth(opts *bind.TransactOpts, subId uint64) (*types.Transaction, error) {
	return _SubscriptionAPI.contract.Transact(opts, "fundSubscriptionWithEth", subId)
}

func (_SubscriptionAPI *SubscriptionAPISession) FundSubscriptionWithEth(subId uint64) (*types.Transaction, error) {
	return _SubscriptionAPI.Contract.FundSubscriptionWithEth(&_SubscriptionAPI.TransactOpts, subId)
}

func (_SubscriptionAPI *SubscriptionAPITransactorSession) FundSubscriptionWithEth(subId uint64) (*types.Transaction, error) {
	return _SubscriptionAPI.Contract.FundSubscriptionWithEth(&_SubscriptionAPI.TransactOpts, subId)
}

func (_SubscriptionAPI *SubscriptionAPITransactor) OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _SubscriptionAPI.contract.Transact(opts, "onTokenTransfer", arg0, amount, data)
}

func (_SubscriptionAPI *SubscriptionAPISession) OnTokenTransfer(arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _SubscriptionAPI.Contract.OnTokenTransfer(&_SubscriptionAPI.TransactOpts, arg0, amount, data)
}

func (_SubscriptionAPI *SubscriptionAPITransactorSession) OnTokenTransfer(arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _SubscriptionAPI.Contract.OnTokenTransfer(&_SubscriptionAPI.TransactOpts, arg0, amount, data)
}

func (_SubscriptionAPI *SubscriptionAPITransactor) OracleWithdraw(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _SubscriptionAPI.contract.Transact(opts, "oracleWithdraw", recipient, amount)
}

func (_SubscriptionAPI *SubscriptionAPISession) OracleWithdraw(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _SubscriptionAPI.Contract.OracleWithdraw(&_SubscriptionAPI.TransactOpts, recipient, amount)
}

func (_SubscriptionAPI *SubscriptionAPITransactorSession) OracleWithdraw(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _SubscriptionAPI.Contract.OracleWithdraw(&_SubscriptionAPI.TransactOpts, recipient, amount)
}

func (_SubscriptionAPI *SubscriptionAPITransactor) OracleWithdrawEth(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _SubscriptionAPI.contract.Transact(opts, "oracleWithdrawEth", recipient, amount)
}

func (_SubscriptionAPI *SubscriptionAPISession) OracleWithdrawEth(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _SubscriptionAPI.Contract.OracleWithdrawEth(&_SubscriptionAPI.TransactOpts, recipient, amount)
}

func (_SubscriptionAPI *SubscriptionAPITransactorSession) OracleWithdrawEth(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _SubscriptionAPI.Contract.OracleWithdrawEth(&_SubscriptionAPI.TransactOpts, recipient, amount)
}

func (_SubscriptionAPI *SubscriptionAPITransactor) OwnerCancelSubscription(opts *bind.TransactOpts, subId uint64) (*types.Transaction, error) {
	return _SubscriptionAPI.contract.Transact(opts, "ownerCancelSubscription", subId)
}

func (_SubscriptionAPI *SubscriptionAPISession) OwnerCancelSubscription(subId uint64) (*types.Transaction, error) {
	return _SubscriptionAPI.Contract.OwnerCancelSubscription(&_SubscriptionAPI.TransactOpts, subId)
}

func (_SubscriptionAPI *SubscriptionAPITransactorSession) OwnerCancelSubscription(subId uint64) (*types.Transaction, error) {
	return _SubscriptionAPI.Contract.OwnerCancelSubscription(&_SubscriptionAPI.TransactOpts, subId)
}

func (_SubscriptionAPI *SubscriptionAPITransactor) RecoverEthFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _SubscriptionAPI.contract.Transact(opts, "recoverEthFunds", to)
}

func (_SubscriptionAPI *SubscriptionAPISession) RecoverEthFunds(to common.Address) (*types.Transaction, error) {
	return _SubscriptionAPI.Contract.RecoverEthFunds(&_SubscriptionAPI.TransactOpts, to)
}

func (_SubscriptionAPI *SubscriptionAPITransactorSession) RecoverEthFunds(to common.Address) (*types.Transaction, error) {
	return _SubscriptionAPI.Contract.RecoverEthFunds(&_SubscriptionAPI.TransactOpts, to)
}

func (_SubscriptionAPI *SubscriptionAPITransactor) RecoverFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _SubscriptionAPI.contract.Transact(opts, "recoverFunds", to)
}

func (_SubscriptionAPI *SubscriptionAPISession) RecoverFunds(to common.Address) (*types.Transaction, error) {
	return _SubscriptionAPI.Contract.RecoverFunds(&_SubscriptionAPI.TransactOpts, to)
}

func (_SubscriptionAPI *SubscriptionAPITransactorSession) RecoverFunds(to common.Address) (*types.Transaction, error) {
	return _SubscriptionAPI.Contract.RecoverFunds(&_SubscriptionAPI.TransactOpts, to)
}

func (_SubscriptionAPI *SubscriptionAPITransactor) RequestSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId uint64, newOwner common.Address) (*types.Transaction, error) {
	return _SubscriptionAPI.contract.Transact(opts, "requestSubscriptionOwnerTransfer", subId, newOwner)
}

func (_SubscriptionAPI *SubscriptionAPISession) RequestSubscriptionOwnerTransfer(subId uint64, newOwner common.Address) (*types.Transaction, error) {
	return _SubscriptionAPI.Contract.RequestSubscriptionOwnerTransfer(&_SubscriptionAPI.TransactOpts, subId, newOwner)
}

func (_SubscriptionAPI *SubscriptionAPITransactorSession) RequestSubscriptionOwnerTransfer(subId uint64, newOwner common.Address) (*types.Transaction, error) {
	return _SubscriptionAPI.Contract.RequestSubscriptionOwnerTransfer(&_SubscriptionAPI.TransactOpts, subId, newOwner)
}

func (_SubscriptionAPI *SubscriptionAPITransactor) SetLINK(opts *bind.TransactOpts, link common.Address) (*types.Transaction, error) {
	return _SubscriptionAPI.contract.Transact(opts, "setLINK", link)
}

func (_SubscriptionAPI *SubscriptionAPISession) SetLINK(link common.Address) (*types.Transaction, error) {
	return _SubscriptionAPI.Contract.SetLINK(&_SubscriptionAPI.TransactOpts, link)
}

func (_SubscriptionAPI *SubscriptionAPITransactorSession) SetLINK(link common.Address) (*types.Transaction, error) {
	return _SubscriptionAPI.Contract.SetLINK(&_SubscriptionAPI.TransactOpts, link)
}

func (_SubscriptionAPI *SubscriptionAPITransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _SubscriptionAPI.contract.Transact(opts, "transferOwnership", to)
}

func (_SubscriptionAPI *SubscriptionAPISession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _SubscriptionAPI.Contract.TransferOwnership(&_SubscriptionAPI.TransactOpts, to)
}

func (_SubscriptionAPI *SubscriptionAPITransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _SubscriptionAPI.Contract.TransferOwnership(&_SubscriptionAPI.TransactOpts, to)
}

type SubscriptionAPIEthFundsRecoveredIterator struct {
	Event *SubscriptionAPIEthFundsRecovered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *SubscriptionAPIEthFundsRecoveredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SubscriptionAPIEthFundsRecovered)
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
		it.Event = new(SubscriptionAPIEthFundsRecovered)
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

func (it *SubscriptionAPIEthFundsRecoveredIterator) Error() error {
	return it.fail
}

func (it *SubscriptionAPIEthFundsRecoveredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type SubscriptionAPIEthFundsRecovered struct {
	To     common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_SubscriptionAPI *SubscriptionAPIFilterer) FilterEthFundsRecovered(opts *bind.FilterOpts) (*SubscriptionAPIEthFundsRecoveredIterator, error) {

	logs, sub, err := _SubscriptionAPI.contract.FilterLogs(opts, "EthFundsRecovered")
	if err != nil {
		return nil, err
	}
	return &SubscriptionAPIEthFundsRecoveredIterator{contract: _SubscriptionAPI.contract, event: "EthFundsRecovered", logs: logs, sub: sub}, nil
}

func (_SubscriptionAPI *SubscriptionAPIFilterer) WatchEthFundsRecovered(opts *bind.WatchOpts, sink chan<- *SubscriptionAPIEthFundsRecovered) (event.Subscription, error) {

	logs, sub, err := _SubscriptionAPI.contract.WatchLogs(opts, "EthFundsRecovered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(SubscriptionAPIEthFundsRecovered)
				if err := _SubscriptionAPI.contract.UnpackLog(event, "EthFundsRecovered", log); err != nil {
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

func (_SubscriptionAPI *SubscriptionAPIFilterer) ParseEthFundsRecovered(log types.Log) (*SubscriptionAPIEthFundsRecovered, error) {
	event := new(SubscriptionAPIEthFundsRecovered)
	if err := _SubscriptionAPI.contract.UnpackLog(event, "EthFundsRecovered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type SubscriptionAPIFundsRecoveredIterator struct {
	Event *SubscriptionAPIFundsRecovered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *SubscriptionAPIFundsRecoveredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SubscriptionAPIFundsRecovered)
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
		it.Event = new(SubscriptionAPIFundsRecovered)
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

func (it *SubscriptionAPIFundsRecoveredIterator) Error() error {
	return it.fail
}

func (it *SubscriptionAPIFundsRecoveredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type SubscriptionAPIFundsRecovered struct {
	To     common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_SubscriptionAPI *SubscriptionAPIFilterer) FilterFundsRecovered(opts *bind.FilterOpts) (*SubscriptionAPIFundsRecoveredIterator, error) {

	logs, sub, err := _SubscriptionAPI.contract.FilterLogs(opts, "FundsRecovered")
	if err != nil {
		return nil, err
	}
	return &SubscriptionAPIFundsRecoveredIterator{contract: _SubscriptionAPI.contract, event: "FundsRecovered", logs: logs, sub: sub}, nil
}

func (_SubscriptionAPI *SubscriptionAPIFilterer) WatchFundsRecovered(opts *bind.WatchOpts, sink chan<- *SubscriptionAPIFundsRecovered) (event.Subscription, error) {

	logs, sub, err := _SubscriptionAPI.contract.WatchLogs(opts, "FundsRecovered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(SubscriptionAPIFundsRecovered)
				if err := _SubscriptionAPI.contract.UnpackLog(event, "FundsRecovered", log); err != nil {
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

func (_SubscriptionAPI *SubscriptionAPIFilterer) ParseFundsRecovered(log types.Log) (*SubscriptionAPIFundsRecovered, error) {
	event := new(SubscriptionAPIFundsRecovered)
	if err := _SubscriptionAPI.contract.UnpackLog(event, "FundsRecovered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type SubscriptionAPIOwnershipTransferRequestedIterator struct {
	Event *SubscriptionAPIOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *SubscriptionAPIOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SubscriptionAPIOwnershipTransferRequested)
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
		it.Event = new(SubscriptionAPIOwnershipTransferRequested)
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

func (it *SubscriptionAPIOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *SubscriptionAPIOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type SubscriptionAPIOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_SubscriptionAPI *SubscriptionAPIFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*SubscriptionAPIOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _SubscriptionAPI.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &SubscriptionAPIOwnershipTransferRequestedIterator{contract: _SubscriptionAPI.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_SubscriptionAPI *SubscriptionAPIFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *SubscriptionAPIOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _SubscriptionAPI.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(SubscriptionAPIOwnershipTransferRequested)
				if err := _SubscriptionAPI.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_SubscriptionAPI *SubscriptionAPIFilterer) ParseOwnershipTransferRequested(log types.Log) (*SubscriptionAPIOwnershipTransferRequested, error) {
	event := new(SubscriptionAPIOwnershipTransferRequested)
	if err := _SubscriptionAPI.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type SubscriptionAPIOwnershipTransferredIterator struct {
	Event *SubscriptionAPIOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *SubscriptionAPIOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SubscriptionAPIOwnershipTransferred)
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
		it.Event = new(SubscriptionAPIOwnershipTransferred)
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

func (it *SubscriptionAPIOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *SubscriptionAPIOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type SubscriptionAPIOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_SubscriptionAPI *SubscriptionAPIFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*SubscriptionAPIOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _SubscriptionAPI.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &SubscriptionAPIOwnershipTransferredIterator{contract: _SubscriptionAPI.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_SubscriptionAPI *SubscriptionAPIFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *SubscriptionAPIOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _SubscriptionAPI.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(SubscriptionAPIOwnershipTransferred)
				if err := _SubscriptionAPI.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_SubscriptionAPI *SubscriptionAPIFilterer) ParseOwnershipTransferred(log types.Log) (*SubscriptionAPIOwnershipTransferred, error) {
	event := new(SubscriptionAPIOwnershipTransferred)
	if err := _SubscriptionAPI.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type SubscriptionAPISubscriptionCanceledIterator struct {
	Event *SubscriptionAPISubscriptionCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *SubscriptionAPISubscriptionCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SubscriptionAPISubscriptionCanceled)
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
		it.Event = new(SubscriptionAPISubscriptionCanceled)
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

func (it *SubscriptionAPISubscriptionCanceledIterator) Error() error {
	return it.fail
}

func (it *SubscriptionAPISubscriptionCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type SubscriptionAPISubscriptionCanceled struct {
	SubId      uint64
	To         common.Address
	AmountLink *big.Int
	AmountEth  *big.Int
	Raw        types.Log
}

func (_SubscriptionAPI *SubscriptionAPIFilterer) FilterSubscriptionCanceled(opts *bind.FilterOpts, subId []uint64) (*SubscriptionAPISubscriptionCanceledIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _SubscriptionAPI.contract.FilterLogs(opts, "SubscriptionCanceled", subIdRule)
	if err != nil {
		return nil, err
	}
	return &SubscriptionAPISubscriptionCanceledIterator{contract: _SubscriptionAPI.contract, event: "SubscriptionCanceled", logs: logs, sub: sub}, nil
}

func (_SubscriptionAPI *SubscriptionAPIFilterer) WatchSubscriptionCanceled(opts *bind.WatchOpts, sink chan<- *SubscriptionAPISubscriptionCanceled, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _SubscriptionAPI.contract.WatchLogs(opts, "SubscriptionCanceled", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(SubscriptionAPISubscriptionCanceled)
				if err := _SubscriptionAPI.contract.UnpackLog(event, "SubscriptionCanceled", log); err != nil {
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

func (_SubscriptionAPI *SubscriptionAPIFilterer) ParseSubscriptionCanceled(log types.Log) (*SubscriptionAPISubscriptionCanceled, error) {
	event := new(SubscriptionAPISubscriptionCanceled)
	if err := _SubscriptionAPI.contract.UnpackLog(event, "SubscriptionCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type SubscriptionAPISubscriptionConsumerAddedIterator struct {
	Event *SubscriptionAPISubscriptionConsumerAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *SubscriptionAPISubscriptionConsumerAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SubscriptionAPISubscriptionConsumerAdded)
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
		it.Event = new(SubscriptionAPISubscriptionConsumerAdded)
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

func (it *SubscriptionAPISubscriptionConsumerAddedIterator) Error() error {
	return it.fail
}

func (it *SubscriptionAPISubscriptionConsumerAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type SubscriptionAPISubscriptionConsumerAdded struct {
	SubId    uint64
	Consumer common.Address
	Raw      types.Log
}

func (_SubscriptionAPI *SubscriptionAPIFilterer) FilterSubscriptionConsumerAdded(opts *bind.FilterOpts, subId []uint64) (*SubscriptionAPISubscriptionConsumerAddedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _SubscriptionAPI.contract.FilterLogs(opts, "SubscriptionConsumerAdded", subIdRule)
	if err != nil {
		return nil, err
	}
	return &SubscriptionAPISubscriptionConsumerAddedIterator{contract: _SubscriptionAPI.contract, event: "SubscriptionConsumerAdded", logs: logs, sub: sub}, nil
}

func (_SubscriptionAPI *SubscriptionAPIFilterer) WatchSubscriptionConsumerAdded(opts *bind.WatchOpts, sink chan<- *SubscriptionAPISubscriptionConsumerAdded, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _SubscriptionAPI.contract.WatchLogs(opts, "SubscriptionConsumerAdded", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(SubscriptionAPISubscriptionConsumerAdded)
				if err := _SubscriptionAPI.contract.UnpackLog(event, "SubscriptionConsumerAdded", log); err != nil {
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

func (_SubscriptionAPI *SubscriptionAPIFilterer) ParseSubscriptionConsumerAdded(log types.Log) (*SubscriptionAPISubscriptionConsumerAdded, error) {
	event := new(SubscriptionAPISubscriptionConsumerAdded)
	if err := _SubscriptionAPI.contract.UnpackLog(event, "SubscriptionConsumerAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type SubscriptionAPISubscriptionConsumerRemovedIterator struct {
	Event *SubscriptionAPISubscriptionConsumerRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *SubscriptionAPISubscriptionConsumerRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SubscriptionAPISubscriptionConsumerRemoved)
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
		it.Event = new(SubscriptionAPISubscriptionConsumerRemoved)
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

func (it *SubscriptionAPISubscriptionConsumerRemovedIterator) Error() error {
	return it.fail
}

func (it *SubscriptionAPISubscriptionConsumerRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type SubscriptionAPISubscriptionConsumerRemoved struct {
	SubId    uint64
	Consumer common.Address
	Raw      types.Log
}

func (_SubscriptionAPI *SubscriptionAPIFilterer) FilterSubscriptionConsumerRemoved(opts *bind.FilterOpts, subId []uint64) (*SubscriptionAPISubscriptionConsumerRemovedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _SubscriptionAPI.contract.FilterLogs(opts, "SubscriptionConsumerRemoved", subIdRule)
	if err != nil {
		return nil, err
	}
	return &SubscriptionAPISubscriptionConsumerRemovedIterator{contract: _SubscriptionAPI.contract, event: "SubscriptionConsumerRemoved", logs: logs, sub: sub}, nil
}

func (_SubscriptionAPI *SubscriptionAPIFilterer) WatchSubscriptionConsumerRemoved(opts *bind.WatchOpts, sink chan<- *SubscriptionAPISubscriptionConsumerRemoved, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _SubscriptionAPI.contract.WatchLogs(opts, "SubscriptionConsumerRemoved", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(SubscriptionAPISubscriptionConsumerRemoved)
				if err := _SubscriptionAPI.contract.UnpackLog(event, "SubscriptionConsumerRemoved", log); err != nil {
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

func (_SubscriptionAPI *SubscriptionAPIFilterer) ParseSubscriptionConsumerRemoved(log types.Log) (*SubscriptionAPISubscriptionConsumerRemoved, error) {
	event := new(SubscriptionAPISubscriptionConsumerRemoved)
	if err := _SubscriptionAPI.contract.UnpackLog(event, "SubscriptionConsumerRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type SubscriptionAPISubscriptionCreatedIterator struct {
	Event *SubscriptionAPISubscriptionCreated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *SubscriptionAPISubscriptionCreatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SubscriptionAPISubscriptionCreated)
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
		it.Event = new(SubscriptionAPISubscriptionCreated)
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

func (it *SubscriptionAPISubscriptionCreatedIterator) Error() error {
	return it.fail
}

func (it *SubscriptionAPISubscriptionCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type SubscriptionAPISubscriptionCreated struct {
	SubId uint64
	Owner common.Address
	Raw   types.Log
}

func (_SubscriptionAPI *SubscriptionAPIFilterer) FilterSubscriptionCreated(opts *bind.FilterOpts, subId []uint64) (*SubscriptionAPISubscriptionCreatedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _SubscriptionAPI.contract.FilterLogs(opts, "SubscriptionCreated", subIdRule)
	if err != nil {
		return nil, err
	}
	return &SubscriptionAPISubscriptionCreatedIterator{contract: _SubscriptionAPI.contract, event: "SubscriptionCreated", logs: logs, sub: sub}, nil
}

func (_SubscriptionAPI *SubscriptionAPIFilterer) WatchSubscriptionCreated(opts *bind.WatchOpts, sink chan<- *SubscriptionAPISubscriptionCreated, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _SubscriptionAPI.contract.WatchLogs(opts, "SubscriptionCreated", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(SubscriptionAPISubscriptionCreated)
				if err := _SubscriptionAPI.contract.UnpackLog(event, "SubscriptionCreated", log); err != nil {
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

func (_SubscriptionAPI *SubscriptionAPIFilterer) ParseSubscriptionCreated(log types.Log) (*SubscriptionAPISubscriptionCreated, error) {
	event := new(SubscriptionAPISubscriptionCreated)
	if err := _SubscriptionAPI.contract.UnpackLog(event, "SubscriptionCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type SubscriptionAPISubscriptionFundedIterator struct {
	Event *SubscriptionAPISubscriptionFunded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *SubscriptionAPISubscriptionFundedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SubscriptionAPISubscriptionFunded)
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
		it.Event = new(SubscriptionAPISubscriptionFunded)
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

func (it *SubscriptionAPISubscriptionFundedIterator) Error() error {
	return it.fail
}

func (it *SubscriptionAPISubscriptionFundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type SubscriptionAPISubscriptionFunded struct {
	SubId      uint64
	OldBalance *big.Int
	NewBalance *big.Int
	Raw        types.Log
}

func (_SubscriptionAPI *SubscriptionAPIFilterer) FilterSubscriptionFunded(opts *bind.FilterOpts, subId []uint64) (*SubscriptionAPISubscriptionFundedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _SubscriptionAPI.contract.FilterLogs(opts, "SubscriptionFunded", subIdRule)
	if err != nil {
		return nil, err
	}
	return &SubscriptionAPISubscriptionFundedIterator{contract: _SubscriptionAPI.contract, event: "SubscriptionFunded", logs: logs, sub: sub}, nil
}

func (_SubscriptionAPI *SubscriptionAPIFilterer) WatchSubscriptionFunded(opts *bind.WatchOpts, sink chan<- *SubscriptionAPISubscriptionFunded, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _SubscriptionAPI.contract.WatchLogs(opts, "SubscriptionFunded", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(SubscriptionAPISubscriptionFunded)
				if err := _SubscriptionAPI.contract.UnpackLog(event, "SubscriptionFunded", log); err != nil {
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

func (_SubscriptionAPI *SubscriptionAPIFilterer) ParseSubscriptionFunded(log types.Log) (*SubscriptionAPISubscriptionFunded, error) {
	event := new(SubscriptionAPISubscriptionFunded)
	if err := _SubscriptionAPI.contract.UnpackLog(event, "SubscriptionFunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type SubscriptionAPISubscriptionFundedWithEthIterator struct {
	Event *SubscriptionAPISubscriptionFundedWithEth

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *SubscriptionAPISubscriptionFundedWithEthIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SubscriptionAPISubscriptionFundedWithEth)
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
		it.Event = new(SubscriptionAPISubscriptionFundedWithEth)
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

func (it *SubscriptionAPISubscriptionFundedWithEthIterator) Error() error {
	return it.fail
}

func (it *SubscriptionAPISubscriptionFundedWithEthIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type SubscriptionAPISubscriptionFundedWithEth struct {
	SubId         uint64
	OldEthBalance *big.Int
	NewEthBalance *big.Int
	Raw           types.Log
}

func (_SubscriptionAPI *SubscriptionAPIFilterer) FilterSubscriptionFundedWithEth(opts *bind.FilterOpts, subId []uint64) (*SubscriptionAPISubscriptionFundedWithEthIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _SubscriptionAPI.contract.FilterLogs(opts, "SubscriptionFundedWithEth", subIdRule)
	if err != nil {
		return nil, err
	}
	return &SubscriptionAPISubscriptionFundedWithEthIterator{contract: _SubscriptionAPI.contract, event: "SubscriptionFundedWithEth", logs: logs, sub: sub}, nil
}

func (_SubscriptionAPI *SubscriptionAPIFilterer) WatchSubscriptionFundedWithEth(opts *bind.WatchOpts, sink chan<- *SubscriptionAPISubscriptionFundedWithEth, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _SubscriptionAPI.contract.WatchLogs(opts, "SubscriptionFundedWithEth", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(SubscriptionAPISubscriptionFundedWithEth)
				if err := _SubscriptionAPI.contract.UnpackLog(event, "SubscriptionFundedWithEth", log); err != nil {
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

func (_SubscriptionAPI *SubscriptionAPIFilterer) ParseSubscriptionFundedWithEth(log types.Log) (*SubscriptionAPISubscriptionFundedWithEth, error) {
	event := new(SubscriptionAPISubscriptionFundedWithEth)
	if err := _SubscriptionAPI.contract.UnpackLog(event, "SubscriptionFundedWithEth", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type SubscriptionAPISubscriptionOwnerTransferRequestedIterator struct {
	Event *SubscriptionAPISubscriptionOwnerTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *SubscriptionAPISubscriptionOwnerTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SubscriptionAPISubscriptionOwnerTransferRequested)
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
		it.Event = new(SubscriptionAPISubscriptionOwnerTransferRequested)
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

func (it *SubscriptionAPISubscriptionOwnerTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *SubscriptionAPISubscriptionOwnerTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type SubscriptionAPISubscriptionOwnerTransferRequested struct {
	SubId uint64
	From  common.Address
	To    common.Address
	Raw   types.Log
}

func (_SubscriptionAPI *SubscriptionAPIFilterer) FilterSubscriptionOwnerTransferRequested(opts *bind.FilterOpts, subId []uint64) (*SubscriptionAPISubscriptionOwnerTransferRequestedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _SubscriptionAPI.contract.FilterLogs(opts, "SubscriptionOwnerTransferRequested", subIdRule)
	if err != nil {
		return nil, err
	}
	return &SubscriptionAPISubscriptionOwnerTransferRequestedIterator{contract: _SubscriptionAPI.contract, event: "SubscriptionOwnerTransferRequested", logs: logs, sub: sub}, nil
}

func (_SubscriptionAPI *SubscriptionAPIFilterer) WatchSubscriptionOwnerTransferRequested(opts *bind.WatchOpts, sink chan<- *SubscriptionAPISubscriptionOwnerTransferRequested, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _SubscriptionAPI.contract.WatchLogs(opts, "SubscriptionOwnerTransferRequested", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(SubscriptionAPISubscriptionOwnerTransferRequested)
				if err := _SubscriptionAPI.contract.UnpackLog(event, "SubscriptionOwnerTransferRequested", log); err != nil {
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

func (_SubscriptionAPI *SubscriptionAPIFilterer) ParseSubscriptionOwnerTransferRequested(log types.Log) (*SubscriptionAPISubscriptionOwnerTransferRequested, error) {
	event := new(SubscriptionAPISubscriptionOwnerTransferRequested)
	if err := _SubscriptionAPI.contract.UnpackLog(event, "SubscriptionOwnerTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type SubscriptionAPISubscriptionOwnerTransferredIterator struct {
	Event *SubscriptionAPISubscriptionOwnerTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *SubscriptionAPISubscriptionOwnerTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SubscriptionAPISubscriptionOwnerTransferred)
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
		it.Event = new(SubscriptionAPISubscriptionOwnerTransferred)
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

func (it *SubscriptionAPISubscriptionOwnerTransferredIterator) Error() error {
	return it.fail
}

func (it *SubscriptionAPISubscriptionOwnerTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type SubscriptionAPISubscriptionOwnerTransferred struct {
	SubId uint64
	From  common.Address
	To    common.Address
	Raw   types.Log
}

func (_SubscriptionAPI *SubscriptionAPIFilterer) FilterSubscriptionOwnerTransferred(opts *bind.FilterOpts, subId []uint64) (*SubscriptionAPISubscriptionOwnerTransferredIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _SubscriptionAPI.contract.FilterLogs(opts, "SubscriptionOwnerTransferred", subIdRule)
	if err != nil {
		return nil, err
	}
	return &SubscriptionAPISubscriptionOwnerTransferredIterator{contract: _SubscriptionAPI.contract, event: "SubscriptionOwnerTransferred", logs: logs, sub: sub}, nil
}

func (_SubscriptionAPI *SubscriptionAPIFilterer) WatchSubscriptionOwnerTransferred(opts *bind.WatchOpts, sink chan<- *SubscriptionAPISubscriptionOwnerTransferred, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _SubscriptionAPI.contract.WatchLogs(opts, "SubscriptionOwnerTransferred", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(SubscriptionAPISubscriptionOwnerTransferred)
				if err := _SubscriptionAPI.contract.UnpackLog(event, "SubscriptionOwnerTransferred", log); err != nil {
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

func (_SubscriptionAPI *SubscriptionAPIFilterer) ParseSubscriptionOwnerTransferred(log types.Log) (*SubscriptionAPISubscriptionOwnerTransferred, error) {
	event := new(SubscriptionAPISubscriptionOwnerTransferred)
	if err := _SubscriptionAPI.contract.UnpackLog(event, "SubscriptionOwnerTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetSubscription struct {
	Balance    *big.Int
	EthBalance *big.Int
	Owner      common.Address
	Consumers  []common.Address
}

func (_SubscriptionAPI *SubscriptionAPI) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _SubscriptionAPI.abi.Events["EthFundsRecovered"].ID:
		return _SubscriptionAPI.ParseEthFundsRecovered(log)
	case _SubscriptionAPI.abi.Events["FundsRecovered"].ID:
		return _SubscriptionAPI.ParseFundsRecovered(log)
	case _SubscriptionAPI.abi.Events["OwnershipTransferRequested"].ID:
		return _SubscriptionAPI.ParseOwnershipTransferRequested(log)
	case _SubscriptionAPI.abi.Events["OwnershipTransferred"].ID:
		return _SubscriptionAPI.ParseOwnershipTransferred(log)
	case _SubscriptionAPI.abi.Events["SubscriptionCanceled"].ID:
		return _SubscriptionAPI.ParseSubscriptionCanceled(log)
	case _SubscriptionAPI.abi.Events["SubscriptionConsumerAdded"].ID:
		return _SubscriptionAPI.ParseSubscriptionConsumerAdded(log)
	case _SubscriptionAPI.abi.Events["SubscriptionConsumerRemoved"].ID:
		return _SubscriptionAPI.ParseSubscriptionConsumerRemoved(log)
	case _SubscriptionAPI.abi.Events["SubscriptionCreated"].ID:
		return _SubscriptionAPI.ParseSubscriptionCreated(log)
	case _SubscriptionAPI.abi.Events["SubscriptionFunded"].ID:
		return _SubscriptionAPI.ParseSubscriptionFunded(log)
	case _SubscriptionAPI.abi.Events["SubscriptionFundedWithEth"].ID:
		return _SubscriptionAPI.ParseSubscriptionFundedWithEth(log)
	case _SubscriptionAPI.abi.Events["SubscriptionOwnerTransferRequested"].ID:
		return _SubscriptionAPI.ParseSubscriptionOwnerTransferRequested(log)
	case _SubscriptionAPI.abi.Events["SubscriptionOwnerTransferred"].ID:
		return _SubscriptionAPI.ParseSubscriptionOwnerTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (SubscriptionAPIEthFundsRecovered) Topic() common.Hash {
	return common.HexToHash("0x879c9ea2b9d5345b84ccd12610b032602808517cebdb795007f3dcb4df377317")
}

func (SubscriptionAPIFundsRecovered) Topic() common.Hash {
	return common.HexToHash("0x59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b436600")
}

func (SubscriptionAPIOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (SubscriptionAPIOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (SubscriptionAPISubscriptionCanceled) Topic() common.Hash {
	return common.HexToHash("0xb3a76e2c52bb7a484e325118e4cadc4346fdaa16a7df0621c5525e1f7ab4d495")
}

func (SubscriptionAPISubscriptionConsumerAdded) Topic() common.Hash {
	return common.HexToHash("0x43dc749a04ac8fb825cbd514f7c0e13f13bc6f2ee66043b76629d51776cff8e0")
}

func (SubscriptionAPISubscriptionConsumerRemoved) Topic() common.Hash {
	return common.HexToHash("0x182bff9831466789164ca77075fffd84916d35a8180ba73c27e45634549b445b")
}

func (SubscriptionAPISubscriptionCreated) Topic() common.Hash {
	return common.HexToHash("0x464722b4166576d3dcbba877b999bc35cf911f4eaf434b7eba68fa113951d0bf")
}

func (SubscriptionAPISubscriptionFunded) Topic() common.Hash {
	return common.HexToHash("0xd39ec07f4e209f627a4c427971473820dc129761ba28de8906bd56f57101d4f8")
}

func (SubscriptionAPISubscriptionFundedWithEth) Topic() common.Hash {
	return common.HexToHash("0x4e4421a74ab9b76003de9d9527153d002f1338c36c1bbc180069671631cf0958")
}

func (SubscriptionAPISubscriptionOwnerTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x69436ea6df009049404f564eff6622cd00522b0bd6a89efd9e52a355c4a879be")
}

func (SubscriptionAPISubscriptionOwnerTransferred) Topic() common.Hash {
	return common.HexToHash("0x6f1dc65165ffffedfd8e507b4a0f1fcfdada045ed11f6c26ba27cedfe87802f0")
}

func (_SubscriptionAPI *SubscriptionAPI) Address() common.Address {
	return _SubscriptionAPI.address
}

type SubscriptionAPIInterface interface {
	LINK(opts *bind.CallOpts) (common.Address, error)

	MAXCONSUMERS(opts *bind.CallOpts) (uint16, error)

	GetSubscription(opts *bind.CallOpts, subId uint64) (GetSubscription,

		error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SCurrentSubId(opts *bind.CallOpts) (uint64, error)

	STotalBalance(opts *bind.CallOpts) (*big.Int, error)

	STotalEthBalance(opts *bind.CallOpts) (*big.Int, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AcceptSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId uint64) (*types.Transaction, error)

	AddConsumer(opts *bind.TransactOpts, subId uint64, consumer common.Address) (*types.Transaction, error)

	CreateSubscription(opts *bind.TransactOpts) (*types.Transaction, error)

	FundSubscriptionWithEth(opts *bind.TransactOpts, subId uint64) (*types.Transaction, error)

	OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error)

	OracleWithdraw(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error)

	OracleWithdrawEth(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error)

	OwnerCancelSubscription(opts *bind.TransactOpts, subId uint64) (*types.Transaction, error)

	RecoverEthFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	RecoverFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	RequestSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId uint64, newOwner common.Address) (*types.Transaction, error)

	SetLINK(opts *bind.TransactOpts, link common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterEthFundsRecovered(opts *bind.FilterOpts) (*SubscriptionAPIEthFundsRecoveredIterator, error)

	WatchEthFundsRecovered(opts *bind.WatchOpts, sink chan<- *SubscriptionAPIEthFundsRecovered) (event.Subscription, error)

	ParseEthFundsRecovered(log types.Log) (*SubscriptionAPIEthFundsRecovered, error)

	FilterFundsRecovered(opts *bind.FilterOpts) (*SubscriptionAPIFundsRecoveredIterator, error)

	WatchFundsRecovered(opts *bind.WatchOpts, sink chan<- *SubscriptionAPIFundsRecovered) (event.Subscription, error)

	ParseFundsRecovered(log types.Log) (*SubscriptionAPIFundsRecovered, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*SubscriptionAPIOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *SubscriptionAPIOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*SubscriptionAPIOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*SubscriptionAPIOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *SubscriptionAPIOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*SubscriptionAPIOwnershipTransferred, error)

	FilterSubscriptionCanceled(opts *bind.FilterOpts, subId []uint64) (*SubscriptionAPISubscriptionCanceledIterator, error)

	WatchSubscriptionCanceled(opts *bind.WatchOpts, sink chan<- *SubscriptionAPISubscriptionCanceled, subId []uint64) (event.Subscription, error)

	ParseSubscriptionCanceled(log types.Log) (*SubscriptionAPISubscriptionCanceled, error)

	FilterSubscriptionConsumerAdded(opts *bind.FilterOpts, subId []uint64) (*SubscriptionAPISubscriptionConsumerAddedIterator, error)

	WatchSubscriptionConsumerAdded(opts *bind.WatchOpts, sink chan<- *SubscriptionAPISubscriptionConsumerAdded, subId []uint64) (event.Subscription, error)

	ParseSubscriptionConsumerAdded(log types.Log) (*SubscriptionAPISubscriptionConsumerAdded, error)

	FilterSubscriptionConsumerRemoved(opts *bind.FilterOpts, subId []uint64) (*SubscriptionAPISubscriptionConsumerRemovedIterator, error)

	WatchSubscriptionConsumerRemoved(opts *bind.WatchOpts, sink chan<- *SubscriptionAPISubscriptionConsumerRemoved, subId []uint64) (event.Subscription, error)

	ParseSubscriptionConsumerRemoved(log types.Log) (*SubscriptionAPISubscriptionConsumerRemoved, error)

	FilterSubscriptionCreated(opts *bind.FilterOpts, subId []uint64) (*SubscriptionAPISubscriptionCreatedIterator, error)

	WatchSubscriptionCreated(opts *bind.WatchOpts, sink chan<- *SubscriptionAPISubscriptionCreated, subId []uint64) (event.Subscription, error)

	ParseSubscriptionCreated(log types.Log) (*SubscriptionAPISubscriptionCreated, error)

	FilterSubscriptionFunded(opts *bind.FilterOpts, subId []uint64) (*SubscriptionAPISubscriptionFundedIterator, error)

	WatchSubscriptionFunded(opts *bind.WatchOpts, sink chan<- *SubscriptionAPISubscriptionFunded, subId []uint64) (event.Subscription, error)

	ParseSubscriptionFunded(log types.Log) (*SubscriptionAPISubscriptionFunded, error)

	FilterSubscriptionFundedWithEth(opts *bind.FilterOpts, subId []uint64) (*SubscriptionAPISubscriptionFundedWithEthIterator, error)

	WatchSubscriptionFundedWithEth(opts *bind.WatchOpts, sink chan<- *SubscriptionAPISubscriptionFundedWithEth, subId []uint64) (event.Subscription, error)

	ParseSubscriptionFundedWithEth(log types.Log) (*SubscriptionAPISubscriptionFundedWithEth, error)

	FilterSubscriptionOwnerTransferRequested(opts *bind.FilterOpts, subId []uint64) (*SubscriptionAPISubscriptionOwnerTransferRequestedIterator, error)

	WatchSubscriptionOwnerTransferRequested(opts *bind.WatchOpts, sink chan<- *SubscriptionAPISubscriptionOwnerTransferRequested, subId []uint64) (event.Subscription, error)

	ParseSubscriptionOwnerTransferRequested(log types.Log) (*SubscriptionAPISubscriptionOwnerTransferRequested, error)

	FilterSubscriptionOwnerTransferred(opts *bind.FilterOpts, subId []uint64) (*SubscriptionAPISubscriptionOwnerTransferredIterator, error)

	WatchSubscriptionOwnerTransferred(opts *bind.WatchOpts, sink chan<- *SubscriptionAPISubscriptionOwnerTransferred, subId []uint64) (event.Subscription, error)

	ParseSubscriptionOwnerTransferred(log types.Log) (*SubscriptionAPISubscriptionOwnerTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
