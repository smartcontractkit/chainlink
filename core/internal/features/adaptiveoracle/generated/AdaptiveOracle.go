// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package adaptiveoracle

import (
	"context"
	"errors"
	"math/big"
	"strings"
	"time"

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
	_ = abi.ConvertType
	_ = time.Tick
	_ = context.Background
)

// AdaptiveOracleMetaData contains all meta data concerning the AdaptiveOracle contract.
var AdaptiveOracleMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"decimals_\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"description_\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"description\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAdaptiveRateLogic\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAggregator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAnswer\",\"inputs\":[{\"name\":\"roundId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLastAdaptiveRate\",\"inputs\":[],\"outputs\":[{\"name\":\"rate\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLastSafeReferenceRate\",\"inputs\":[],\"outputs\":[{\"name\":\"rate\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getReferenceRateAdapter\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoundData\",\"inputs\":[{\"name\":\"_roundId\",\"type\":\"uint80\",\"internalType\":\"uint80\"}],\"outputs\":[{\"name\":\"roundId\",\"type\":\"uint80\",\"internalType\":\"uint80\"},{\"name\":\"answer\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"startedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"updatedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"answeredInRound\",\"type\":\"uint80\",\"internalType\":\"uint80\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTimestamp\",\"inputs\":[{\"name\":\"roundId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getUnderlyingRateFeed\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestAnswer\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestRound\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestRoundData\",\"inputs\":[],\"outputs\":[{\"name\":\"roundId\",\"type\":\"uint80\",\"internalType\":\"uint80\"},{\"name\":\"answer\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"startedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"updatedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"answeredInRound\",\"type\":\"uint80\",\"internalType\":\"uint80\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestTimestamp\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"resetAnchors\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setAdaptiveRateLogic\",\"inputs\":[{\"name\":\"adaptiveRateLogic\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setAggregator\",\"inputs\":[{\"name\":\"aggregator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setReferenceRateAdapter\",\"inputs\":[{\"name\":\"referenceRateAdapter\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUnderlyingRateFeed\",\"inputs\":[{\"name\":\"underlyingRateFeed\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transformAnswer\",\"inputs\":[{\"name\":\"median\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"outputs\":[{\"name\":\"transformedAnswer\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"version\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"event\",\"name\":\"AdaptiveRateLogicSet\",\"inputs\":[{\"name\":\"old\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"current\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AdaptiveRateUpdated\",\"inputs\":[{\"name\":\"adaptiveRate\",\"type\":\"int256\",\"indexed\":false,\"internalType\":\"int256\"},{\"name\":\"referenceRate\",\"type\":\"int256\",\"indexed\":false,\"internalType\":\"int256\"},{\"name\":\"marketRate\",\"type\":\"int256\",\"indexed\":false,\"internalType\":\"int256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AggregatorSet\",\"inputs\":[{\"name\":\"old\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"current\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AnchorsReset\",\"inputs\":[{\"name\":\"referenceRate\",\"type\":\"int256\",\"indexed\":false,\"internalType\":\"int256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AnswerUpdated\",\"inputs\":[{\"name\":\"current\",\"type\":\"int256\",\"indexed\":true,\"internalType\":\"int256\"},{\"name\":\"roundId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"updatedAt\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NewRound\",\"inputs\":[{\"name\":\"roundId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"startedBy\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"startedAt\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReferenceRateAdapterSet\",\"inputs\":[{\"name\":\"old\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"current\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UnderlyingRateFeedSet\",\"inputs\":[{\"name\":\"old\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"current\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"CannotTransferToSelf\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DecimalsMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MustBeProposedOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyAggregatorCanCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnerCannotBeZero\",\"inputs\":[]}]",
	Bin: "0x60a060405234801561001057600080fd5b5060405161179038038061179083398101604081905261002f91610136565b8260006001600160a01b03821661005957604051639b15e16f60e01b815260040160405180910390fd5b600180546001600160a01b0319166001600160a01b038481169190911790915581161561008957610089816100a7565b505060ff8216608052600261009e82826102c4565b50505050610382565b336001600160a01b038216036100d057604051636d6c4ee560e11b815260040160405180910390fd5b600080546001600160a01b0319166001600160a01b03838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b634e487b7160e01b600052604160045260246000fd5b60008060006060848603121561014b57600080fd5b83516001600160a01b038116811461016257600080fd5b602085015190935060ff8116811461017957600080fd5b60408501519092506001600160401b0381111561019557600080fd5b8401601f810186136101a657600080fd5b80516001600160401b038111156101bf576101bf610120565b604051601f8201601f19908116603f011681016001600160401b03811182821017156101ed576101ed610120565b60405281815282820160200188101561020557600080fd5b60005b8281101561022457602081850181015183830182015201610208565b506000602083830101528093505050509250925092565b600181811c9082168061024f57607f821691505b60208210810361026f57634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156102bf57806000526020600020601f840160051c8101602085101561029c5750805b601f840160051c820191505b818110156102bc57600081556001016102a8565b50505b505050565b81516001600160401b038111156102dd576102dd610120565b6102f1816102eb845461023b565b84610275565b6020601f821160018114610325576000831561030d5750848201515b600019600385901b1c1916600184901b1784556102bc565b600084815260208120601f198516915b828110156103555787850151825560209485019460019092019101610335565b50848210156103735786840151600019600387901b60f8161c191681555b50505050600190811b01905550565b6080516113e56103ab600039600081816101d90152818161044a01526104f901526113e56000f3fe608060405234801561001057600080fd5b50600436106101735760003560e01c806379ba5097116100de578063b5ab58dc11610097578063f2fde38b11610071578063f2fde38b1461033b578063f316a2191461034e578063f9120af614610361578063feaf968c1461037457600080fd5b8063b5ab58dc14610304578063b633620c14610317578063f1e4c49e1461032a57600080fd5b806379ba5097146102755780638205bf6a1461027d5780638da5cb5b146102855780639a6fc8f5146102965780639ff3189f146102e0578063a185c46c146102f157600080fd5b806350d25bcd1161013057806350d25bcd1461021457806354fd4d501461022a5780635f45d21714610231578063668a0f021461024d5780637284e4161461025557806372ba8a051461026a57600080fd5b8063049ae55b146101785780630dc7456814610182578063186de40e146101955780631ebabb43146101bf578063313ce567146101d25780633ad59dbc14610203575b600080fd5b61018061037c565b005b61018061019036600461106e565b610431565b6005546001600160a01b03165b6040516001600160a01b0390911681526020015b60405180910390f35b6101806101cd36600461106e565b610596565b60405160ff7f00000000000000000000000000000000000000000000000000000000000000001681526020016101b6565b6003546001600160a01b03166101a2565b61021c6105f0565b6040519081526020016101b6565b600161021c565b6007546008545b604080519283526020830191909152016101b6565b61021c6107c0565b61025d610833565b6040516101b6919061109e565b600954600a54610238565b6101806108c5565b61021c610948565b6001546001600160a01b03166101a2565b6102a96102a4366004611104565b610992565b6040805169ffffffffffffffffffff968716815260208101959095528401929092526060830152909116608082015260a0016101b6565b6004546001600160a01b03166101a2565b61021c6102ff366004611121565b610a25565b61021c610312366004611143565b610b4e565b61021c610325366004611143565b610bbd565b6006546001600160a01b03166101a2565b61018061034936600461106e565b610bef565b61018061035c36600461106e565b610c03565b61018061036f36600461106e565b610d0c565b6102a9610e8d565b610384610f2a565b600061038e610f57565b600981905542600a5560035460405163128a4c0960e31b8152600481018390529192506001600160a01b031690639452604890602401600060405180830381600087803b1580156103de57600080fd5b505af11580156103f2573d6000803e3d6000fd5b505050507e2436516982c1a9948ca2de900eed964a2768a7788465ee1ff188d6d7365a578160405161042691815260200190565b60405180910390a150565b610439610f2a565b6001600160a01b038116156104f4577f000000000000000000000000000000000000000000000000000000000000000060ff16816001600160a01b031663313ce5676040518163ffffffff1660e01b8152600401602060405180830381865afa1580156104aa573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104ce919061115c565b60ff16146104ef57604051635a8dbaed60e01b815260040160405180910390fd5b610544565b6003547f000000000000000000000000000000000000000000000000000000000000000060ff908116600160a01b909204161461054457604051635a8dbaed60e01b815260040160405180910390fd5b600680546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f587234defe9e4ba3e8d3b8b6ad1523bfb568b8f573392aec0a7a92a1b298407f90600090a35050565b61059e610f2a565b600580546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f03fff9b1913deb37b91e6a1a7b9ff287ca87cd10ccf15eff9bf3e56c2bc44b6b90600090a35050565b6000806000600460009054906101000a90046001600160a01b03166001600160a01b0316637f70b3a16040518163ffffffff1660e01b81526004016040805180830381865afa158015610647573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061066b919061117f565b915091508061067a5760075491505b600354604080516350d25bcd60e01b815290516000926001600160a01b0316916350d25bcd9160048083019260209291908290030181865afa1580156106c4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106e891906111b5565b905060008382126106f957836106fb565b815b6006549091506001600160a01b0316806107185750949350505050565b600654604080516350d25bcd60e01b815290516000926001600160a01b0316916350d25bcd9160048083019260209291908290030181865afa158015610762573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061078691906111b5565b6003549091506107a190600160a01b900460ff16600a6112cb565b6107ab82856112da565b6107b5919061130a565b965050505050505090565b60035460408051633345078160e11b815290516000926001600160a01b03169163668a0f029160048083019260209291908290030181865afa15801561080a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061082e91906111b5565b905090565b60606002805461084290611346565b80601f016020809104026020016040519081016040528092919081815260200182805461086e90611346565b80156108bb5780601f10610890576101008083540402835291602001916108bb565b820191906000526020600020905b81548152906001019060200180831161089e57829003601f168201915b5050505050905090565b6000546001600160a01b031633146108f05760405163015aa1e360e11b815260040160405180910390fd5b600180546001600160a01b0319808216339081179093556000805490911681556040516001600160a01b03909216929183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60035460408051634102dfb560e11b815290516000926001600160a01b031691638205bf6a9160048083019260209291908290030181865afa15801561080a573d6000803e3d6000fd5b600354604051639a6fc8f560e01b815269ffffffffffffffffffff8316600482015260009182918291829182916001600160a01b031690639a6fc8f59060240160a060405180830381865afa1580156109ef573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a139190611380565b939a9299509097509550909350915050565b6003546000906001600160a01b03163314610a5357604051630b86ae6b60e41b815260040160405180910390fd5b6000610a5d610f57565b600954600554600a54600854604051633e5a5e3d60e21b81526004810185905260248101869052604481018a90526064810192909252608482015292935090916000916001600160a01b03169063f96978f49060a4016020604051808303816000875af1158015610ad2573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610af691906111b5565b600981905542600a5560408051828152602081018690529081018890529091507f9b7c8bf3ef88bb8390434f301f308adfe8f8fa16942fa34ff2e8e9c94c259b1a9060600160405180910390a1925050505b92915050565b600354604051632d6ad63760e21b8152600481018390526000916001600160a01b03169063b5ab58dc906024015b602060405180830381865afa158015610b99573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b4891906111b5565b600354604051632d8cd88360e21b8152600481018390526000916001600160a01b03169063b633620c90602401610b7c565b610bf7610f2a565b610c0081610ff5565b50565b610c0b610f2a565b6000816001600160a01b031663313ce5676040518163ffffffff1660e01b8152600401602060405180830381865afa158015610c4b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c6f919061115c565b6003549091506001600160a01b031615801590610c9b575060035460ff828116600160a01b9092041614155b15610cb957604051635a8dbaed60e01b815260040160405180910390fd5b600480546001600160a01b038481166001600160a01b0319831681179093556040519116919082907f0b6d8ab269d7033edb82557c5afddef95bb8db22870775aec5686a0f56db6b1d90600090a3505050565b610d14610f2a565b6000816001600160a01b031663313ce5676040518163ffffffff1660e01b8152600401602060405180830381865afa158015610d54573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d78919061115c565b6004549091506001600160a01b031615801590610e1157508060ff16600460009054906101000a90046001600160a01b03166001600160a01b031663313ce5676040518163ffffffff1660e01b8152600401602060405180830381865afa158015610de7573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610e0b919061115c565b60ff1614155b15610e2f57604051635a8dbaed60e01b815260040160405180910390fd5b600380546001600160a01b038481166001600160a81b031983168117600160a01b60ff871602179093556040519116919082907f84e1bba7d2e21dccf1934a4188807abc546c4632d794d2c097bfc4f599230cac90600090a3505050565b6000806000806000600360009054906101000a90046001600160a01b03166001600160a01b031663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa158015610ee8573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610f0c9190611380565b9398509095509350909150610f2190506105f0565b93509091929394565b6001546001600160a01b03163314610f55576040516315ae3a6f60e11b815260040160405180910390fd5b565b6000806000600460009054906101000a90046001600160a01b03166001600160a01b0316637f70b3a16040518163ffffffff1660e01b81526004016040805180830381865afa158015610fae573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610fd2919061117f565b915091508015610feb5750600781905542600855919050565b6007549250505090565b336001600160a01b0382160361101e57604051636d6c4ee560e11b815260040160405180910390fd5b600080546001600160a01b0319166001600160a01b03838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60006020828403121561108057600080fd5b81356001600160a01b038116811461109757600080fd5b9392505050565b602081526000825180602084015260005b818110156110cc57602081860181015160408684010152016110af565b506000604082850101526040601f19601f83011684010191505092915050565b69ffffffffffffffffffff81168114610c0057600080fd5b60006020828403121561111657600080fd5b8135611097816110ec565b6000806040838503121561113457600080fd5b50508035926020909101359150565b60006020828403121561115557600080fd5b5035919050565b60006020828403121561116e57600080fd5b815160ff8116811461109757600080fd5b6000806040838503121561119257600080fd5b8251602084015190925080151581146111aa57600080fd5b809150509250929050565b6000602082840312156111c757600080fd5b5051919050565b634e487b7160e01b600052601160045260246000fd5b6001815b600184111561121f57808504811115611203576112036111ce565b600184161561121157908102905b60019390931c9280026111e8565b935093915050565b60008261123657506001610b48565b8161124357506000610b48565b816001811461125957600281146112635761127f565b6001915050610b48565b60ff841115611274576112746111ce565b50506001821b610b48565b5060208310610133831016604e8410600b84101617156112a2575081810a610b48565b6112af60001984846111e4565b80600019048211156112c3576112c36111ce565b029392505050565b600061109760ff841683611227565b80820260008212600160ff1b841416156112f6576112f66111ce565b8181058314821517610b4857610b486111ce565b60008261132757634e487b7160e01b600052601260045260246000fd5b600160ff1b821460001984141615611341576113416111ce565b500590565b600181811c9082168061135a57607f821691505b60208210810361137a57634e487b7160e01b600052602260045260246000fd5b50919050565b600080600080600060a0868803121561139857600080fd5b85516113a3816110ec565b60208701516040880151606089015160808a0151939850919650945092506113ca816110ec565b80915050929550929590935056fea164736f6c634300081a000a",
}

// AdaptiveOracleABI is the input ABI used to generate the binding from.
// Deprecated: Use AdaptiveOracleMetaData.ABI instead.
var AdaptiveOracleABI = AdaptiveOracleMetaData.ABI

// AdaptiveOracleBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use AdaptiveOracleMetaData.Bin instead.
var AdaptiveOracleBin = AdaptiveOracleMetaData.Bin

// DeployAdaptiveOracle deploys a new Ethereum contract, binding an instance of AdaptiveOracle to it.
func DeployAdaptiveOracle(auth *bind.TransactOpts, backend bind.ContractBackend, newOwner common.Address, decimals_ uint8, description_ string) (common.Address, *types.Transaction, *AdaptiveOracle, error) {
	parsed, err := AdaptiveOracleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AdaptiveOracleBin), backend, newOwner, decimals_, description_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AdaptiveOracle{AdaptiveOracleCaller: AdaptiveOracleCaller{contract: contract}, AdaptiveOracleTransactor: AdaptiveOracleTransactor{contract: contract}, AdaptiveOracleFilterer: AdaptiveOracleFilterer{contract: contract}}, nil
}

// AdaptiveOracle is an auto generated Go binding around an Ethereum contract.
type AdaptiveOracle struct {
	AdaptiveOracleCaller     // Read-only binding to the contract
	AdaptiveOracleTransactor // Write-only binding to the contract
	AdaptiveOracleFilterer   // Log filterer for contract events
}

// AdaptiveOracleCaller is an auto generated read-only Go binding around an Ethereum contract.
type AdaptiveOracleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AdaptiveOracleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AdaptiveOracleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AdaptiveOracleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AdaptiveOracleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AdaptiveOracleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AdaptiveOracleSession struct {
	Contract     *AdaptiveOracle   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AdaptiveOracleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AdaptiveOracleCallerSession struct {
	Contract *AdaptiveOracleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// AdaptiveOracleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AdaptiveOracleTransactorSession struct {
	Contract     *AdaptiveOracleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// AdaptiveOracleRaw is an auto generated low-level Go binding around an Ethereum contract.
type AdaptiveOracleRaw struct {
	Contract *AdaptiveOracle // Generic contract binding to access the raw methods on
}

// AdaptiveOracleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AdaptiveOracleCallerRaw struct {
	Contract *AdaptiveOracleCaller // Generic read-only contract binding to access the raw methods on
}

// AdaptiveOracleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AdaptiveOracleTransactorRaw struct {
	Contract *AdaptiveOracleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAdaptiveOracle creates a new instance of AdaptiveOracle, bound to a specific deployed contract.
func NewAdaptiveOracle(address common.Address, backend bind.ContractBackend) (*AdaptiveOracle, error) {
	contract, err := bindAdaptiveOracle(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AdaptiveOracle{AdaptiveOracleCaller: AdaptiveOracleCaller{contract: contract}, AdaptiveOracleTransactor: AdaptiveOracleTransactor{contract: contract}, AdaptiveOracleFilterer: AdaptiveOracleFilterer{contract: contract}}, nil
}

// NewAdaptiveOracleCaller creates a new read-only instance of AdaptiveOracle, bound to a specific deployed contract.
func NewAdaptiveOracleCaller(address common.Address, caller bind.ContractCaller) (*AdaptiveOracleCaller, error) {
	contract, err := bindAdaptiveOracle(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AdaptiveOracleCaller{contract: contract}, nil
}

// NewAdaptiveOracleTransactor creates a new write-only instance of AdaptiveOracle, bound to a specific deployed contract.
func NewAdaptiveOracleTransactor(address common.Address, transactor bind.ContractTransactor) (*AdaptiveOracleTransactor, error) {
	contract, err := bindAdaptiveOracle(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AdaptiveOracleTransactor{contract: contract}, nil
}

// NewAdaptiveOracleFilterer creates a new log filterer instance of AdaptiveOracle, bound to a specific deployed contract.
func NewAdaptiveOracleFilterer(address common.Address, filterer bind.ContractFilterer) (*AdaptiveOracleFilterer, error) {
	contract, err := bindAdaptiveOracle(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AdaptiveOracleFilterer{contract: contract}, nil
}

// bindAdaptiveOracle binds a generic wrapper to an already deployed contract.
func bindAdaptiveOracle(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AdaptiveOracleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AdaptiveOracle *AdaptiveOracleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AdaptiveOracle.Contract.AdaptiveOracleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AdaptiveOracle *AdaptiveOracleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AdaptiveOracle.Contract.AdaptiveOracleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AdaptiveOracle *AdaptiveOracleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AdaptiveOracle.Contract.AdaptiveOracleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AdaptiveOracle *AdaptiveOracleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AdaptiveOracle.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AdaptiveOracle *AdaptiveOracleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AdaptiveOracle.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AdaptiveOracle *AdaptiveOracleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AdaptiveOracle.Contract.contract.Transact(opts, method, params...)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_AdaptiveOracle *AdaptiveOracleCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _AdaptiveOracle.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_AdaptiveOracle *AdaptiveOracleSession) Decimals() (uint8, error) {
	return _AdaptiveOracle.Contract.Decimals(&_AdaptiveOracle.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_AdaptiveOracle *AdaptiveOracleCallerSession) Decimals() (uint8, error) {
	return _AdaptiveOracle.Contract.Decimals(&_AdaptiveOracle.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_AdaptiveOracle *AdaptiveOracleCaller) Description(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AdaptiveOracle.contract.Call(opts, &out, "description")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_AdaptiveOracle *AdaptiveOracleSession) Description() (string, error) {
	return _AdaptiveOracle.Contract.Description(&_AdaptiveOracle.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_AdaptiveOracle *AdaptiveOracleCallerSession) Description() (string, error) {
	return _AdaptiveOracle.Contract.Description(&_AdaptiveOracle.CallOpts)
}

// GetAdaptiveRateLogic is a free data retrieval call binding the contract method 0x186de40e.
//
// Solidity: function getAdaptiveRateLogic() view returns(address)
func (_AdaptiveOracle *AdaptiveOracleCaller) GetAdaptiveRateLogic(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AdaptiveOracle.contract.Call(opts, &out, "getAdaptiveRateLogic")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAdaptiveRateLogic is a free data retrieval call binding the contract method 0x186de40e.
//
// Solidity: function getAdaptiveRateLogic() view returns(address)
func (_AdaptiveOracle *AdaptiveOracleSession) GetAdaptiveRateLogic() (common.Address, error) {
	return _AdaptiveOracle.Contract.GetAdaptiveRateLogic(&_AdaptiveOracle.CallOpts)
}

// GetAdaptiveRateLogic is a free data retrieval call binding the contract method 0x186de40e.
//
// Solidity: function getAdaptiveRateLogic() view returns(address)
func (_AdaptiveOracle *AdaptiveOracleCallerSession) GetAdaptiveRateLogic() (common.Address, error) {
	return _AdaptiveOracle.Contract.GetAdaptiveRateLogic(&_AdaptiveOracle.CallOpts)
}

// GetAggregator is a free data retrieval call binding the contract method 0x3ad59dbc.
//
// Solidity: function getAggregator() view returns(address)
func (_AdaptiveOracle *AdaptiveOracleCaller) GetAggregator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AdaptiveOracle.contract.Call(opts, &out, "getAggregator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAggregator is a free data retrieval call binding the contract method 0x3ad59dbc.
//
// Solidity: function getAggregator() view returns(address)
func (_AdaptiveOracle *AdaptiveOracleSession) GetAggregator() (common.Address, error) {
	return _AdaptiveOracle.Contract.GetAggregator(&_AdaptiveOracle.CallOpts)
}

// GetAggregator is a free data retrieval call binding the contract method 0x3ad59dbc.
//
// Solidity: function getAggregator() view returns(address)
func (_AdaptiveOracle *AdaptiveOracleCallerSession) GetAggregator() (common.Address, error) {
	return _AdaptiveOracle.Contract.GetAggregator(&_AdaptiveOracle.CallOpts)
}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 roundId) view returns(int256)
func (_AdaptiveOracle *AdaptiveOracleCaller) GetAnswer(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AdaptiveOracle.contract.Call(opts, &out, "getAnswer", roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 roundId) view returns(int256)
func (_AdaptiveOracle *AdaptiveOracleSession) GetAnswer(roundId *big.Int) (*big.Int, error) {
	return _AdaptiveOracle.Contract.GetAnswer(&_AdaptiveOracle.CallOpts, roundId)
}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 roundId) view returns(int256)
func (_AdaptiveOracle *AdaptiveOracleCallerSession) GetAnswer(roundId *big.Int) (*big.Int, error) {
	return _AdaptiveOracle.Contract.GetAnswer(&_AdaptiveOracle.CallOpts, roundId)
}

// GetLastAdaptiveRate is a free data retrieval call binding the contract method 0x72ba8a05.
//
// Solidity: function getLastAdaptiveRate() view returns(int256 rate, uint256 timestamp)
func (_AdaptiveOracle *AdaptiveOracleCaller) GetLastAdaptiveRate(opts *bind.CallOpts) (struct {
	Rate      *big.Int
	Timestamp *big.Int
}, error) {
	var out []interface{}
	err := _AdaptiveOracle.contract.Call(opts, &out, "getLastAdaptiveRate")

	outstruct := new(struct {
		Rate      *big.Int
		Timestamp *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Rate = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Timestamp = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetLastAdaptiveRate is a free data retrieval call binding the contract method 0x72ba8a05.
//
// Solidity: function getLastAdaptiveRate() view returns(int256 rate, uint256 timestamp)
func (_AdaptiveOracle *AdaptiveOracleSession) GetLastAdaptiveRate() (struct {
	Rate      *big.Int
	Timestamp *big.Int
}, error) {
	return _AdaptiveOracle.Contract.GetLastAdaptiveRate(&_AdaptiveOracle.CallOpts)
}

// GetLastAdaptiveRate is a free data retrieval call binding the contract method 0x72ba8a05.
//
// Solidity: function getLastAdaptiveRate() view returns(int256 rate, uint256 timestamp)
func (_AdaptiveOracle *AdaptiveOracleCallerSession) GetLastAdaptiveRate() (struct {
	Rate      *big.Int
	Timestamp *big.Int
}, error) {
	return _AdaptiveOracle.Contract.GetLastAdaptiveRate(&_AdaptiveOracle.CallOpts)
}

// GetLastSafeReferenceRate is a free data retrieval call binding the contract method 0x5f45d217.
//
// Solidity: function getLastSafeReferenceRate() view returns(int256 rate, uint256 timestamp)
func (_AdaptiveOracle *AdaptiveOracleCaller) GetLastSafeReferenceRate(opts *bind.CallOpts) (struct {
	Rate      *big.Int
	Timestamp *big.Int
}, error) {
	var out []interface{}
	err := _AdaptiveOracle.contract.Call(opts, &out, "getLastSafeReferenceRate")

	outstruct := new(struct {
		Rate      *big.Int
		Timestamp *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Rate = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Timestamp = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetLastSafeReferenceRate is a free data retrieval call binding the contract method 0x5f45d217.
//
// Solidity: function getLastSafeReferenceRate() view returns(int256 rate, uint256 timestamp)
func (_AdaptiveOracle *AdaptiveOracleSession) GetLastSafeReferenceRate() (struct {
	Rate      *big.Int
	Timestamp *big.Int
}, error) {
	return _AdaptiveOracle.Contract.GetLastSafeReferenceRate(&_AdaptiveOracle.CallOpts)
}

// GetLastSafeReferenceRate is a free data retrieval call binding the contract method 0x5f45d217.
//
// Solidity: function getLastSafeReferenceRate() view returns(int256 rate, uint256 timestamp)
func (_AdaptiveOracle *AdaptiveOracleCallerSession) GetLastSafeReferenceRate() (struct {
	Rate      *big.Int
	Timestamp *big.Int
}, error) {
	return _AdaptiveOracle.Contract.GetLastSafeReferenceRate(&_AdaptiveOracle.CallOpts)
}

// GetReferenceRateAdapter is a free data retrieval call binding the contract method 0x9ff3189f.
//
// Solidity: function getReferenceRateAdapter() view returns(address)
func (_AdaptiveOracle *AdaptiveOracleCaller) GetReferenceRateAdapter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AdaptiveOracle.contract.Call(opts, &out, "getReferenceRateAdapter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetReferenceRateAdapter is a free data retrieval call binding the contract method 0x9ff3189f.
//
// Solidity: function getReferenceRateAdapter() view returns(address)
func (_AdaptiveOracle *AdaptiveOracleSession) GetReferenceRateAdapter() (common.Address, error) {
	return _AdaptiveOracle.Contract.GetReferenceRateAdapter(&_AdaptiveOracle.CallOpts)
}

// GetReferenceRateAdapter is a free data retrieval call binding the contract method 0x9ff3189f.
//
// Solidity: function getReferenceRateAdapter() view returns(address)
func (_AdaptiveOracle *AdaptiveOracleCallerSession) GetReferenceRateAdapter() (common.Address, error) {
	return _AdaptiveOracle.Contract.GetReferenceRateAdapter(&_AdaptiveOracle.CallOpts)
}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AdaptiveOracle *AdaptiveOracleCaller) GetRoundData(opts *bind.CallOpts, _roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	var out []interface{}
	err := _AdaptiveOracle.contract.Call(opts, &out, "getRoundData", _roundId)

	outstruct := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RoundId = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Answer = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.StartedAt = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.UpdatedAt = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.AnsweredInRound = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AdaptiveOracle *AdaptiveOracleSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AdaptiveOracle.Contract.GetRoundData(&_AdaptiveOracle.CallOpts, _roundId)
}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AdaptiveOracle *AdaptiveOracleCallerSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AdaptiveOracle.Contract.GetRoundData(&_AdaptiveOracle.CallOpts, _roundId)
}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 roundId) view returns(uint256)
func (_AdaptiveOracle *AdaptiveOracleCaller) GetTimestamp(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AdaptiveOracle.contract.Call(opts, &out, "getTimestamp", roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 roundId) view returns(uint256)
func (_AdaptiveOracle *AdaptiveOracleSession) GetTimestamp(roundId *big.Int) (*big.Int, error) {
	return _AdaptiveOracle.Contract.GetTimestamp(&_AdaptiveOracle.CallOpts, roundId)
}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 roundId) view returns(uint256)
func (_AdaptiveOracle *AdaptiveOracleCallerSession) GetTimestamp(roundId *big.Int) (*big.Int, error) {
	return _AdaptiveOracle.Contract.GetTimestamp(&_AdaptiveOracle.CallOpts, roundId)
}

// GetUnderlyingRateFeed is a free data retrieval call binding the contract method 0xf1e4c49e.
//
// Solidity: function getUnderlyingRateFeed() view returns(address)
func (_AdaptiveOracle *AdaptiveOracleCaller) GetUnderlyingRateFeed(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AdaptiveOracle.contract.Call(opts, &out, "getUnderlyingRateFeed")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetUnderlyingRateFeed is a free data retrieval call binding the contract method 0xf1e4c49e.
//
// Solidity: function getUnderlyingRateFeed() view returns(address)
func (_AdaptiveOracle *AdaptiveOracleSession) GetUnderlyingRateFeed() (common.Address, error) {
	return _AdaptiveOracle.Contract.GetUnderlyingRateFeed(&_AdaptiveOracle.CallOpts)
}

// GetUnderlyingRateFeed is a free data retrieval call binding the contract method 0xf1e4c49e.
//
// Solidity: function getUnderlyingRateFeed() view returns(address)
func (_AdaptiveOracle *AdaptiveOracleCallerSession) GetUnderlyingRateFeed() (common.Address, error) {
	return _AdaptiveOracle.Contract.GetUnderlyingRateFeed(&_AdaptiveOracle.CallOpts)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_AdaptiveOracle *AdaptiveOracleCaller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AdaptiveOracle.contract.Call(opts, &out, "latestAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_AdaptiveOracle *AdaptiveOracleSession) LatestAnswer() (*big.Int, error) {
	return _AdaptiveOracle.Contract.LatestAnswer(&_AdaptiveOracle.CallOpts)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_AdaptiveOracle *AdaptiveOracleCallerSession) LatestAnswer() (*big.Int, error) {
	return _AdaptiveOracle.Contract.LatestAnswer(&_AdaptiveOracle.CallOpts)
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_AdaptiveOracle *AdaptiveOracleCaller) LatestRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AdaptiveOracle.contract.Call(opts, &out, "latestRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_AdaptiveOracle *AdaptiveOracleSession) LatestRound() (*big.Int, error) {
	return _AdaptiveOracle.Contract.LatestRound(&_AdaptiveOracle.CallOpts)
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_AdaptiveOracle *AdaptiveOracleCallerSession) LatestRound() (*big.Int, error) {
	return _AdaptiveOracle.Contract.LatestRound(&_AdaptiveOracle.CallOpts)
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AdaptiveOracle *AdaptiveOracleCaller) LatestRoundData(opts *bind.CallOpts) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	var out []interface{}
	err := _AdaptiveOracle.contract.Call(opts, &out, "latestRoundData")

	outstruct := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RoundId = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Answer = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.StartedAt = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.UpdatedAt = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.AnsweredInRound = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AdaptiveOracle *AdaptiveOracleSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AdaptiveOracle.Contract.LatestRoundData(&_AdaptiveOracle.CallOpts)
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AdaptiveOracle *AdaptiveOracleCallerSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AdaptiveOracle.Contract.LatestRoundData(&_AdaptiveOracle.CallOpts)
}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256)
func (_AdaptiveOracle *AdaptiveOracleCaller) LatestTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AdaptiveOracle.contract.Call(opts, &out, "latestTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256)
func (_AdaptiveOracle *AdaptiveOracleSession) LatestTimestamp() (*big.Int, error) {
	return _AdaptiveOracle.Contract.LatestTimestamp(&_AdaptiveOracle.CallOpts)
}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256)
func (_AdaptiveOracle *AdaptiveOracleCallerSession) LatestTimestamp() (*big.Int, error) {
	return _AdaptiveOracle.Contract.LatestTimestamp(&_AdaptiveOracle.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AdaptiveOracle *AdaptiveOracleCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AdaptiveOracle.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AdaptiveOracle *AdaptiveOracleSession) Owner() (common.Address, error) {
	return _AdaptiveOracle.Contract.Owner(&_AdaptiveOracle.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AdaptiveOracle *AdaptiveOracleCallerSession) Owner() (common.Address, error) {
	return _AdaptiveOracle.Contract.Owner(&_AdaptiveOracle.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() pure returns(uint256)
func (_AdaptiveOracle *AdaptiveOracleCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AdaptiveOracle.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() pure returns(uint256)
func (_AdaptiveOracle *AdaptiveOracleSession) Version() (*big.Int, error) {
	return _AdaptiveOracle.Contract.Version(&_AdaptiveOracle.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() pure returns(uint256)
func (_AdaptiveOracle *AdaptiveOracleCallerSession) Version() (*big.Int, error) {
	return _AdaptiveOracle.Contract.Version(&_AdaptiveOracle.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_AdaptiveOracle *AdaptiveOracleTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AdaptiveOracle.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_AdaptiveOracle *AdaptiveOracleSession) AcceptOwnership() (*types.Transaction, error) {
	return _AdaptiveOracle.Contract.AcceptOwnership(&_AdaptiveOracle.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_AdaptiveOracle *AdaptiveOracleTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _AdaptiveOracle.Contract.AcceptOwnership(&_AdaptiveOracle.TransactOpts)
}

// ResetAnchors is a paid mutator transaction binding the contract method 0x049ae55b.
//
// Solidity: function resetAnchors() returns()
func (_AdaptiveOracle *AdaptiveOracleTransactor) ResetAnchors(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AdaptiveOracle.contract.Transact(opts, "resetAnchors")
}

// ResetAnchors is a paid mutator transaction binding the contract method 0x049ae55b.
//
// Solidity: function resetAnchors() returns()
func (_AdaptiveOracle *AdaptiveOracleSession) ResetAnchors() (*types.Transaction, error) {
	return _AdaptiveOracle.Contract.ResetAnchors(&_AdaptiveOracle.TransactOpts)
}

// ResetAnchors is a paid mutator transaction binding the contract method 0x049ae55b.
//
// Solidity: function resetAnchors() returns()
func (_AdaptiveOracle *AdaptiveOracleTransactorSession) ResetAnchors() (*types.Transaction, error) {
	return _AdaptiveOracle.Contract.ResetAnchors(&_AdaptiveOracle.TransactOpts)
}

// SetAdaptiveRateLogic is a paid mutator transaction binding the contract method 0x1ebabb43.
//
// Solidity: function setAdaptiveRateLogic(address adaptiveRateLogic) returns()
func (_AdaptiveOracle *AdaptiveOracleTransactor) SetAdaptiveRateLogic(opts *bind.TransactOpts, adaptiveRateLogic common.Address) (*types.Transaction, error) {
	return _AdaptiveOracle.contract.Transact(opts, "setAdaptiveRateLogic", adaptiveRateLogic)
}

// SetAdaptiveRateLogic is a paid mutator transaction binding the contract method 0x1ebabb43.
//
// Solidity: function setAdaptiveRateLogic(address adaptiveRateLogic) returns()
func (_AdaptiveOracle *AdaptiveOracleSession) SetAdaptiveRateLogic(adaptiveRateLogic common.Address) (*types.Transaction, error) {
	return _AdaptiveOracle.Contract.SetAdaptiveRateLogic(&_AdaptiveOracle.TransactOpts, adaptiveRateLogic)
}

// SetAdaptiveRateLogic is a paid mutator transaction binding the contract method 0x1ebabb43.
//
// Solidity: function setAdaptiveRateLogic(address adaptiveRateLogic) returns()
func (_AdaptiveOracle *AdaptiveOracleTransactorSession) SetAdaptiveRateLogic(adaptiveRateLogic common.Address) (*types.Transaction, error) {
	return _AdaptiveOracle.Contract.SetAdaptiveRateLogic(&_AdaptiveOracle.TransactOpts, adaptiveRateLogic)
}

// SetAggregator is a paid mutator transaction binding the contract method 0xf9120af6.
//
// Solidity: function setAggregator(address aggregator) returns()
func (_AdaptiveOracle *AdaptiveOracleTransactor) SetAggregator(opts *bind.TransactOpts, aggregator common.Address) (*types.Transaction, error) {
	return _AdaptiveOracle.contract.Transact(opts, "setAggregator", aggregator)
}

// SetAggregator is a paid mutator transaction binding the contract method 0xf9120af6.
//
// Solidity: function setAggregator(address aggregator) returns()
func (_AdaptiveOracle *AdaptiveOracleSession) SetAggregator(aggregator common.Address) (*types.Transaction, error) {
	return _AdaptiveOracle.Contract.SetAggregator(&_AdaptiveOracle.TransactOpts, aggregator)
}

// SetAggregator is a paid mutator transaction binding the contract method 0xf9120af6.
//
// Solidity: function setAggregator(address aggregator) returns()
func (_AdaptiveOracle *AdaptiveOracleTransactorSession) SetAggregator(aggregator common.Address) (*types.Transaction, error) {
	return _AdaptiveOracle.Contract.SetAggregator(&_AdaptiveOracle.TransactOpts, aggregator)
}

// SetReferenceRateAdapter is a paid mutator transaction binding the contract method 0xf316a219.
//
// Solidity: function setReferenceRateAdapter(address referenceRateAdapter) returns()
func (_AdaptiveOracle *AdaptiveOracleTransactor) SetReferenceRateAdapter(opts *bind.TransactOpts, referenceRateAdapter common.Address) (*types.Transaction, error) {
	return _AdaptiveOracle.contract.Transact(opts, "setReferenceRateAdapter", referenceRateAdapter)
}

// SetReferenceRateAdapter is a paid mutator transaction binding the contract method 0xf316a219.
//
// Solidity: function setReferenceRateAdapter(address referenceRateAdapter) returns()
func (_AdaptiveOracle *AdaptiveOracleSession) SetReferenceRateAdapter(referenceRateAdapter common.Address) (*types.Transaction, error) {
	return _AdaptiveOracle.Contract.SetReferenceRateAdapter(&_AdaptiveOracle.TransactOpts, referenceRateAdapter)
}

// SetReferenceRateAdapter is a paid mutator transaction binding the contract method 0xf316a219.
//
// Solidity: function setReferenceRateAdapter(address referenceRateAdapter) returns()
func (_AdaptiveOracle *AdaptiveOracleTransactorSession) SetReferenceRateAdapter(referenceRateAdapter common.Address) (*types.Transaction, error) {
	return _AdaptiveOracle.Contract.SetReferenceRateAdapter(&_AdaptiveOracle.TransactOpts, referenceRateAdapter)
}

// SetUnderlyingRateFeed is a paid mutator transaction binding the contract method 0x0dc74568.
//
// Solidity: function setUnderlyingRateFeed(address underlyingRateFeed) returns()
func (_AdaptiveOracle *AdaptiveOracleTransactor) SetUnderlyingRateFeed(opts *bind.TransactOpts, underlyingRateFeed common.Address) (*types.Transaction, error) {
	return _AdaptiveOracle.contract.Transact(opts, "setUnderlyingRateFeed", underlyingRateFeed)
}

// SetUnderlyingRateFeed is a paid mutator transaction binding the contract method 0x0dc74568.
//
// Solidity: function setUnderlyingRateFeed(address underlyingRateFeed) returns()
func (_AdaptiveOracle *AdaptiveOracleSession) SetUnderlyingRateFeed(underlyingRateFeed common.Address) (*types.Transaction, error) {
	return _AdaptiveOracle.Contract.SetUnderlyingRateFeed(&_AdaptiveOracle.TransactOpts, underlyingRateFeed)
}

// SetUnderlyingRateFeed is a paid mutator transaction binding the contract method 0x0dc74568.
//
// Solidity: function setUnderlyingRateFeed(address underlyingRateFeed) returns()
func (_AdaptiveOracle *AdaptiveOracleTransactorSession) SetUnderlyingRateFeed(underlyingRateFeed common.Address) (*types.Transaction, error) {
	return _AdaptiveOracle.Contract.SetUnderlyingRateFeed(&_AdaptiveOracle.TransactOpts, underlyingRateFeed)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_AdaptiveOracle *AdaptiveOracleTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _AdaptiveOracle.contract.Transact(opts, "transferOwnership", to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_AdaptiveOracle *AdaptiveOracleSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _AdaptiveOracle.Contract.TransferOwnership(&_AdaptiveOracle.TransactOpts, to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_AdaptiveOracle *AdaptiveOracleTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _AdaptiveOracle.Contract.TransferOwnership(&_AdaptiveOracle.TransactOpts, to)
}

// TransformAnswer is a paid mutator transaction binding the contract method 0xa185c46c.
//
// Solidity: function transformAnswer(int256 median, int256 ) returns(int256 transformedAnswer)
func (_AdaptiveOracle *AdaptiveOracleTransactor) TransformAnswer(opts *bind.TransactOpts, median *big.Int, arg1 *big.Int) (*types.Transaction, error) {
	return _AdaptiveOracle.contract.Transact(opts, "transformAnswer", median, arg1)
}

// TransformAnswer is a paid mutator transaction binding the contract method 0xa185c46c.
//
// Solidity: function transformAnswer(int256 median, int256 ) returns(int256 transformedAnswer)
func (_AdaptiveOracle *AdaptiveOracleSession) TransformAnswer(median *big.Int, arg1 *big.Int) (*types.Transaction, error) {
	return _AdaptiveOracle.Contract.TransformAnswer(&_AdaptiveOracle.TransactOpts, median, arg1)
}

// TransformAnswer is a paid mutator transaction binding the contract method 0xa185c46c.
//
// Solidity: function transformAnswer(int256 median, int256 ) returns(int256 transformedAnswer)
func (_AdaptiveOracle *AdaptiveOracleTransactorSession) TransformAnswer(median *big.Int, arg1 *big.Int) (*types.Transaction, error) {
	return _AdaptiveOracle.Contract.TransformAnswer(&_AdaptiveOracle.TransactOpts, median, arg1)
}

// AdaptiveOracleAdaptiveRateLogicSetIterator is returned from FilterAdaptiveRateLogicSet and is used to iterate over the raw logs and unpacked data for AdaptiveRateLogicSet events raised by the AdaptiveOracle contract.
type AdaptiveOracleAdaptiveRateLogicSetIterator struct {
	Event *AdaptiveOracleAdaptiveRateLogicSet // Event containing the contract specifics and raw log

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
func (it *AdaptiveOracleAdaptiveRateLogicSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AdaptiveOracleAdaptiveRateLogicSet)
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
		it.Event = new(AdaptiveOracleAdaptiveRateLogicSet)
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
func (it *AdaptiveOracleAdaptiveRateLogicSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AdaptiveOracleAdaptiveRateLogicSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AdaptiveOracleAdaptiveRateLogicSet represents a AdaptiveRateLogicSet event raised by the AdaptiveOracle contract.
type AdaptiveOracleAdaptiveRateLogicSet struct {
	Old     common.Address
	Current common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAdaptiveRateLogicSet is a free log retrieval operation binding the contract event 0x03fff9b1913deb37b91e6a1a7b9ff287ca87cd10ccf15eff9bf3e56c2bc44b6b.
//
// Solidity: event AdaptiveRateLogicSet(address indexed old, address indexed current)
func (_AdaptiveOracle *AdaptiveOracleFilterer) FilterAdaptiveRateLogicSet(opts *bind.FilterOpts, old []common.Address, current []common.Address) (*AdaptiveOracleAdaptiveRateLogicSetIterator, error) {

	var oldRule []interface{}
	for _, oldItem := range old {
		oldRule = append(oldRule, oldItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _AdaptiveOracle.contract.FilterLogs(opts, "AdaptiveRateLogicSet", oldRule, currentRule)
	if err != nil {
		return nil, err
	}
	return &AdaptiveOracleAdaptiveRateLogicSetIterator{contract: _AdaptiveOracle.contract, event: "AdaptiveRateLogicSet", logs: logs, sub: sub}, nil
}

// WatchAdaptiveRateLogicSet is a free log subscription operation binding the contract event 0x03fff9b1913deb37b91e6a1a7b9ff287ca87cd10ccf15eff9bf3e56c2bc44b6b.
//
// Solidity: event AdaptiveRateLogicSet(address indexed old, address indexed current)
func (_AdaptiveOracle *AdaptiveOracleFilterer) WatchAdaptiveRateLogicSet(opts *bind.WatchOpts, sink chan<- *AdaptiveOracleAdaptiveRateLogicSet, old []common.Address, current []common.Address) (event.Subscription, error) {

	var oldRule []interface{}
	for _, oldItem := range old {
		oldRule = append(oldRule, oldItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _AdaptiveOracle.contract.WatchLogs(opts, "AdaptiveRateLogicSet", oldRule, currentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AdaptiveOracleAdaptiveRateLogicSet)
				if err := _AdaptiveOracle.contract.UnpackLog(event, "AdaptiveRateLogicSet", log); err != nil {
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

// ParseAdaptiveRateLogicSet is a log parse operation binding the contract event 0x03fff9b1913deb37b91e6a1a7b9ff287ca87cd10ccf15eff9bf3e56c2bc44b6b.
//
// Solidity: event AdaptiveRateLogicSet(address indexed old, address indexed current)
func (_AdaptiveOracle *AdaptiveOracleFilterer) ParseAdaptiveRateLogicSet(log types.Log) (*AdaptiveOracleAdaptiveRateLogicSet, error) {
	event := new(AdaptiveOracleAdaptiveRateLogicSet)
	if err := _AdaptiveOracle.contract.UnpackLog(event, "AdaptiveRateLogicSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AdaptiveOracleAdaptiveRateUpdatedIterator is returned from FilterAdaptiveRateUpdated and is used to iterate over the raw logs and unpacked data for AdaptiveRateUpdated events raised by the AdaptiveOracle contract.
type AdaptiveOracleAdaptiveRateUpdatedIterator struct {
	Event *AdaptiveOracleAdaptiveRateUpdated // Event containing the contract specifics and raw log

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
func (it *AdaptiveOracleAdaptiveRateUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AdaptiveOracleAdaptiveRateUpdated)
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
		it.Event = new(AdaptiveOracleAdaptiveRateUpdated)
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
func (it *AdaptiveOracleAdaptiveRateUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AdaptiveOracleAdaptiveRateUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AdaptiveOracleAdaptiveRateUpdated represents a AdaptiveRateUpdated event raised by the AdaptiveOracle contract.
type AdaptiveOracleAdaptiveRateUpdated struct {
	AdaptiveRate  *big.Int
	ReferenceRate *big.Int
	MarketRate    *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdaptiveRateUpdated is a free log retrieval operation binding the contract event 0x9b7c8bf3ef88bb8390434f301f308adfe8f8fa16942fa34ff2e8e9c94c259b1a.
//
// Solidity: event AdaptiveRateUpdated(int256 adaptiveRate, int256 referenceRate, int256 marketRate)
func (_AdaptiveOracle *AdaptiveOracleFilterer) FilterAdaptiveRateUpdated(opts *bind.FilterOpts) (*AdaptiveOracleAdaptiveRateUpdatedIterator, error) {

	logs, sub, err := _AdaptiveOracle.contract.FilterLogs(opts, "AdaptiveRateUpdated")
	if err != nil {
		return nil, err
	}
	return &AdaptiveOracleAdaptiveRateUpdatedIterator{contract: _AdaptiveOracle.contract, event: "AdaptiveRateUpdated", logs: logs, sub: sub}, nil
}

// WatchAdaptiveRateUpdated is a free log subscription operation binding the contract event 0x9b7c8bf3ef88bb8390434f301f308adfe8f8fa16942fa34ff2e8e9c94c259b1a.
//
// Solidity: event AdaptiveRateUpdated(int256 adaptiveRate, int256 referenceRate, int256 marketRate)
func (_AdaptiveOracle *AdaptiveOracleFilterer) WatchAdaptiveRateUpdated(opts *bind.WatchOpts, sink chan<- *AdaptiveOracleAdaptiveRateUpdated) (event.Subscription, error) {

	logs, sub, err := _AdaptiveOracle.contract.WatchLogs(opts, "AdaptiveRateUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AdaptiveOracleAdaptiveRateUpdated)
				if err := _AdaptiveOracle.contract.UnpackLog(event, "AdaptiveRateUpdated", log); err != nil {
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

// ParseAdaptiveRateUpdated is a log parse operation binding the contract event 0x9b7c8bf3ef88bb8390434f301f308adfe8f8fa16942fa34ff2e8e9c94c259b1a.
//
// Solidity: event AdaptiveRateUpdated(int256 adaptiveRate, int256 referenceRate, int256 marketRate)
func (_AdaptiveOracle *AdaptiveOracleFilterer) ParseAdaptiveRateUpdated(log types.Log) (*AdaptiveOracleAdaptiveRateUpdated, error) {
	event := new(AdaptiveOracleAdaptiveRateUpdated)
	if err := _AdaptiveOracle.contract.UnpackLog(event, "AdaptiveRateUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AdaptiveOracleAggregatorSetIterator is returned from FilterAggregatorSet and is used to iterate over the raw logs and unpacked data for AggregatorSet events raised by the AdaptiveOracle contract.
type AdaptiveOracleAggregatorSetIterator struct {
	Event *AdaptiveOracleAggregatorSet // Event containing the contract specifics and raw log

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
func (it *AdaptiveOracleAggregatorSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AdaptiveOracleAggregatorSet)
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
		it.Event = new(AdaptiveOracleAggregatorSet)
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
func (it *AdaptiveOracleAggregatorSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AdaptiveOracleAggregatorSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AdaptiveOracleAggregatorSet represents a AggregatorSet event raised by the AdaptiveOracle contract.
type AdaptiveOracleAggregatorSet struct {
	Old     common.Address
	Current common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAggregatorSet is a free log retrieval operation binding the contract event 0x84e1bba7d2e21dccf1934a4188807abc546c4632d794d2c097bfc4f599230cac.
//
// Solidity: event AggregatorSet(address indexed old, address indexed current)
func (_AdaptiveOracle *AdaptiveOracleFilterer) FilterAggregatorSet(opts *bind.FilterOpts, old []common.Address, current []common.Address) (*AdaptiveOracleAggregatorSetIterator, error) {

	var oldRule []interface{}
	for _, oldItem := range old {
		oldRule = append(oldRule, oldItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _AdaptiveOracle.contract.FilterLogs(opts, "AggregatorSet", oldRule, currentRule)
	if err != nil {
		return nil, err
	}
	return &AdaptiveOracleAggregatorSetIterator{contract: _AdaptiveOracle.contract, event: "AggregatorSet", logs: logs, sub: sub}, nil
}

// WatchAggregatorSet is a free log subscription operation binding the contract event 0x84e1bba7d2e21dccf1934a4188807abc546c4632d794d2c097bfc4f599230cac.
//
// Solidity: event AggregatorSet(address indexed old, address indexed current)
func (_AdaptiveOracle *AdaptiveOracleFilterer) WatchAggregatorSet(opts *bind.WatchOpts, sink chan<- *AdaptiveOracleAggregatorSet, old []common.Address, current []common.Address) (event.Subscription, error) {

	var oldRule []interface{}
	for _, oldItem := range old {
		oldRule = append(oldRule, oldItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _AdaptiveOracle.contract.WatchLogs(opts, "AggregatorSet", oldRule, currentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AdaptiveOracleAggregatorSet)
				if err := _AdaptiveOracle.contract.UnpackLog(event, "AggregatorSet", log); err != nil {
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

// ParseAggregatorSet is a log parse operation binding the contract event 0x84e1bba7d2e21dccf1934a4188807abc546c4632d794d2c097bfc4f599230cac.
//
// Solidity: event AggregatorSet(address indexed old, address indexed current)
func (_AdaptiveOracle *AdaptiveOracleFilterer) ParseAggregatorSet(log types.Log) (*AdaptiveOracleAggregatorSet, error) {
	event := new(AdaptiveOracleAggregatorSet)
	if err := _AdaptiveOracle.contract.UnpackLog(event, "AggregatorSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AdaptiveOracleAnchorsResetIterator is returned from FilterAnchorsReset and is used to iterate over the raw logs and unpacked data for AnchorsReset events raised by the AdaptiveOracle contract.
type AdaptiveOracleAnchorsResetIterator struct {
	Event *AdaptiveOracleAnchorsReset // Event containing the contract specifics and raw log

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
func (it *AdaptiveOracleAnchorsResetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AdaptiveOracleAnchorsReset)
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
		it.Event = new(AdaptiveOracleAnchorsReset)
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
func (it *AdaptiveOracleAnchorsResetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AdaptiveOracleAnchorsResetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AdaptiveOracleAnchorsReset represents a AnchorsReset event raised by the AdaptiveOracle contract.
type AdaptiveOracleAnchorsReset struct {
	ReferenceRate *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAnchorsReset is a free log retrieval operation binding the contract event 0x002436516982c1a9948ca2de900eed964a2768a7788465ee1ff188d6d7365a57.
//
// Solidity: event AnchorsReset(int256 referenceRate)
func (_AdaptiveOracle *AdaptiveOracleFilterer) FilterAnchorsReset(opts *bind.FilterOpts) (*AdaptiveOracleAnchorsResetIterator, error) {

	logs, sub, err := _AdaptiveOracle.contract.FilterLogs(opts, "AnchorsReset")
	if err != nil {
		return nil, err
	}
	return &AdaptiveOracleAnchorsResetIterator{contract: _AdaptiveOracle.contract, event: "AnchorsReset", logs: logs, sub: sub}, nil
}

// WatchAnchorsReset is a free log subscription operation binding the contract event 0x002436516982c1a9948ca2de900eed964a2768a7788465ee1ff188d6d7365a57.
//
// Solidity: event AnchorsReset(int256 referenceRate)
func (_AdaptiveOracle *AdaptiveOracleFilterer) WatchAnchorsReset(opts *bind.WatchOpts, sink chan<- *AdaptiveOracleAnchorsReset) (event.Subscription, error) {

	logs, sub, err := _AdaptiveOracle.contract.WatchLogs(opts, "AnchorsReset")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AdaptiveOracleAnchorsReset)
				if err := _AdaptiveOracle.contract.UnpackLog(event, "AnchorsReset", log); err != nil {
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

// ParseAnchorsReset is a log parse operation binding the contract event 0x002436516982c1a9948ca2de900eed964a2768a7788465ee1ff188d6d7365a57.
//
// Solidity: event AnchorsReset(int256 referenceRate)
func (_AdaptiveOracle *AdaptiveOracleFilterer) ParseAnchorsReset(log types.Log) (*AdaptiveOracleAnchorsReset, error) {
	event := new(AdaptiveOracleAnchorsReset)
	if err := _AdaptiveOracle.contract.UnpackLog(event, "AnchorsReset", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AdaptiveOracleAnswerUpdatedIterator is returned from FilterAnswerUpdated and is used to iterate over the raw logs and unpacked data for AnswerUpdated events raised by the AdaptiveOracle contract.
type AdaptiveOracleAnswerUpdatedIterator struct {
	Event *AdaptiveOracleAnswerUpdated // Event containing the contract specifics and raw log

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
func (it *AdaptiveOracleAnswerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AdaptiveOracleAnswerUpdated)
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
		it.Event = new(AdaptiveOracleAnswerUpdated)
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
func (it *AdaptiveOracleAnswerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AdaptiveOracleAnswerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AdaptiveOracleAnswerUpdated represents a AnswerUpdated event raised by the AdaptiveOracle contract.
type AdaptiveOracleAnswerUpdated struct {
	Current   *big.Int
	RoundId   *big.Int
	UpdatedAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAnswerUpdated is a free log retrieval operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
func (_AdaptiveOracle *AdaptiveOracleFilterer) FilterAnswerUpdated(opts *bind.FilterOpts, current []*big.Int, roundId []*big.Int) (*AdaptiveOracleAnswerUpdatedIterator, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _AdaptiveOracle.contract.FilterLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return &AdaptiveOracleAnswerUpdatedIterator{contract: _AdaptiveOracle.contract, event: "AnswerUpdated", logs: logs, sub: sub}, nil
}

// WatchAnswerUpdated is a free log subscription operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
func (_AdaptiveOracle *AdaptiveOracleFilterer) WatchAnswerUpdated(opts *bind.WatchOpts, sink chan<- *AdaptiveOracleAnswerUpdated, current []*big.Int, roundId []*big.Int) (event.Subscription, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _AdaptiveOracle.contract.WatchLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AdaptiveOracleAnswerUpdated)
				if err := _AdaptiveOracle.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
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

// ParseAnswerUpdated is a log parse operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
func (_AdaptiveOracle *AdaptiveOracleFilterer) ParseAnswerUpdated(log types.Log) (*AdaptiveOracleAnswerUpdated, error) {
	event := new(AdaptiveOracleAnswerUpdated)
	if err := _AdaptiveOracle.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AdaptiveOracleNewRoundIterator is returned from FilterNewRound and is used to iterate over the raw logs and unpacked data for NewRound events raised by the AdaptiveOracle contract.
type AdaptiveOracleNewRoundIterator struct {
	Event *AdaptiveOracleNewRound // Event containing the contract specifics and raw log

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
func (it *AdaptiveOracleNewRoundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AdaptiveOracleNewRound)
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
		it.Event = new(AdaptiveOracleNewRound)
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
func (it *AdaptiveOracleNewRoundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AdaptiveOracleNewRoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AdaptiveOracleNewRound represents a NewRound event raised by the AdaptiveOracle contract.
type AdaptiveOracleNewRound struct {
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterNewRound is a free log retrieval operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_AdaptiveOracle *AdaptiveOracleFilterer) FilterNewRound(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*AdaptiveOracleNewRoundIterator, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _AdaptiveOracle.contract.FilterLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return &AdaptiveOracleNewRoundIterator{contract: _AdaptiveOracle.contract, event: "NewRound", logs: logs, sub: sub}, nil
}

// WatchNewRound is a free log subscription operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_AdaptiveOracle *AdaptiveOracleFilterer) WatchNewRound(opts *bind.WatchOpts, sink chan<- *AdaptiveOracleNewRound, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _AdaptiveOracle.contract.WatchLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AdaptiveOracleNewRound)
				if err := _AdaptiveOracle.contract.UnpackLog(event, "NewRound", log); err != nil {
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

// ParseNewRound is a log parse operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_AdaptiveOracle *AdaptiveOracleFilterer) ParseNewRound(log types.Log) (*AdaptiveOracleNewRound, error) {
	event := new(AdaptiveOracleNewRound)
	if err := _AdaptiveOracle.contract.UnpackLog(event, "NewRound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AdaptiveOracleOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the AdaptiveOracle contract.
type AdaptiveOracleOwnershipTransferRequestedIterator struct {
	Event *AdaptiveOracleOwnershipTransferRequested // Event containing the contract specifics and raw log

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
func (it *AdaptiveOracleOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AdaptiveOracleOwnershipTransferRequested)
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
		it.Event = new(AdaptiveOracleOwnershipTransferRequested)
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
func (it *AdaptiveOracleOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AdaptiveOracleOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AdaptiveOracleOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the AdaptiveOracle contract.
type AdaptiveOracleOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_AdaptiveOracle *AdaptiveOracleFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AdaptiveOracleOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AdaptiveOracle.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AdaptiveOracleOwnershipTransferRequestedIterator{contract: _AdaptiveOracle.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_AdaptiveOracle *AdaptiveOracleFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *AdaptiveOracleOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AdaptiveOracle.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AdaptiveOracleOwnershipTransferRequested)
				if err := _AdaptiveOracle.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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
func (_AdaptiveOracle *AdaptiveOracleFilterer) ParseOwnershipTransferRequested(log types.Log) (*AdaptiveOracleOwnershipTransferRequested, error) {
	event := new(AdaptiveOracleOwnershipTransferRequested)
	if err := _AdaptiveOracle.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AdaptiveOracleOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the AdaptiveOracle contract.
type AdaptiveOracleOwnershipTransferredIterator struct {
	Event *AdaptiveOracleOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *AdaptiveOracleOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AdaptiveOracleOwnershipTransferred)
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
		it.Event = new(AdaptiveOracleOwnershipTransferred)
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
func (it *AdaptiveOracleOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AdaptiveOracleOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AdaptiveOracleOwnershipTransferred represents a OwnershipTransferred event raised by the AdaptiveOracle contract.
type AdaptiveOracleOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_AdaptiveOracle *AdaptiveOracleFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AdaptiveOracleOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AdaptiveOracle.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AdaptiveOracleOwnershipTransferredIterator{contract: _AdaptiveOracle.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_AdaptiveOracle *AdaptiveOracleFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AdaptiveOracleOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AdaptiveOracle.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AdaptiveOracleOwnershipTransferred)
				if err := _AdaptiveOracle.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_AdaptiveOracle *AdaptiveOracleFilterer) ParseOwnershipTransferred(log types.Log) (*AdaptiveOracleOwnershipTransferred, error) {
	event := new(AdaptiveOracleOwnershipTransferred)
	if err := _AdaptiveOracle.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AdaptiveOracleReferenceRateAdapterSetIterator is returned from FilterReferenceRateAdapterSet and is used to iterate over the raw logs and unpacked data for ReferenceRateAdapterSet events raised by the AdaptiveOracle contract.
type AdaptiveOracleReferenceRateAdapterSetIterator struct {
	Event *AdaptiveOracleReferenceRateAdapterSet // Event containing the contract specifics and raw log

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
func (it *AdaptiveOracleReferenceRateAdapterSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AdaptiveOracleReferenceRateAdapterSet)
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
		it.Event = new(AdaptiveOracleReferenceRateAdapterSet)
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
func (it *AdaptiveOracleReferenceRateAdapterSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AdaptiveOracleReferenceRateAdapterSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AdaptiveOracleReferenceRateAdapterSet represents a ReferenceRateAdapterSet event raised by the AdaptiveOracle contract.
type AdaptiveOracleReferenceRateAdapterSet struct {
	Old     common.Address
	Current common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterReferenceRateAdapterSet is a free log retrieval operation binding the contract event 0x0b6d8ab269d7033edb82557c5afddef95bb8db22870775aec5686a0f56db6b1d.
//
// Solidity: event ReferenceRateAdapterSet(address indexed old, address indexed current)
func (_AdaptiveOracle *AdaptiveOracleFilterer) FilterReferenceRateAdapterSet(opts *bind.FilterOpts, old []common.Address, current []common.Address) (*AdaptiveOracleReferenceRateAdapterSetIterator, error) {

	var oldRule []interface{}
	for _, oldItem := range old {
		oldRule = append(oldRule, oldItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _AdaptiveOracle.contract.FilterLogs(opts, "ReferenceRateAdapterSet", oldRule, currentRule)
	if err != nil {
		return nil, err
	}
	return &AdaptiveOracleReferenceRateAdapterSetIterator{contract: _AdaptiveOracle.contract, event: "ReferenceRateAdapterSet", logs: logs, sub: sub}, nil
}

// WatchReferenceRateAdapterSet is a free log subscription operation binding the contract event 0x0b6d8ab269d7033edb82557c5afddef95bb8db22870775aec5686a0f56db6b1d.
//
// Solidity: event ReferenceRateAdapterSet(address indexed old, address indexed current)
func (_AdaptiveOracle *AdaptiveOracleFilterer) WatchReferenceRateAdapterSet(opts *bind.WatchOpts, sink chan<- *AdaptiveOracleReferenceRateAdapterSet, old []common.Address, current []common.Address) (event.Subscription, error) {

	var oldRule []interface{}
	for _, oldItem := range old {
		oldRule = append(oldRule, oldItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _AdaptiveOracle.contract.WatchLogs(opts, "ReferenceRateAdapterSet", oldRule, currentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AdaptiveOracleReferenceRateAdapterSet)
				if err := _AdaptiveOracle.contract.UnpackLog(event, "ReferenceRateAdapterSet", log); err != nil {
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

// ParseReferenceRateAdapterSet is a log parse operation binding the contract event 0x0b6d8ab269d7033edb82557c5afddef95bb8db22870775aec5686a0f56db6b1d.
//
// Solidity: event ReferenceRateAdapterSet(address indexed old, address indexed current)
func (_AdaptiveOracle *AdaptiveOracleFilterer) ParseReferenceRateAdapterSet(log types.Log) (*AdaptiveOracleReferenceRateAdapterSet, error) {
	event := new(AdaptiveOracleReferenceRateAdapterSet)
	if err := _AdaptiveOracle.contract.UnpackLog(event, "ReferenceRateAdapterSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AdaptiveOracleUnderlyingRateFeedSetIterator is returned from FilterUnderlyingRateFeedSet and is used to iterate over the raw logs and unpacked data for UnderlyingRateFeedSet events raised by the AdaptiveOracle contract.
type AdaptiveOracleUnderlyingRateFeedSetIterator struct {
	Event *AdaptiveOracleUnderlyingRateFeedSet // Event containing the contract specifics and raw log

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
func (it *AdaptiveOracleUnderlyingRateFeedSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AdaptiveOracleUnderlyingRateFeedSet)
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
		it.Event = new(AdaptiveOracleUnderlyingRateFeedSet)
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
func (it *AdaptiveOracleUnderlyingRateFeedSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AdaptiveOracleUnderlyingRateFeedSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AdaptiveOracleUnderlyingRateFeedSet represents a UnderlyingRateFeedSet event raised by the AdaptiveOracle contract.
type AdaptiveOracleUnderlyingRateFeedSet struct {
	Old     common.Address
	Current common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnderlyingRateFeedSet is a free log retrieval operation binding the contract event 0x587234defe9e4ba3e8d3b8b6ad1523bfb568b8f573392aec0a7a92a1b298407f.
//
// Solidity: event UnderlyingRateFeedSet(address indexed old, address indexed current)
func (_AdaptiveOracle *AdaptiveOracleFilterer) FilterUnderlyingRateFeedSet(opts *bind.FilterOpts, old []common.Address, current []common.Address) (*AdaptiveOracleUnderlyingRateFeedSetIterator, error) {

	var oldRule []interface{}
	for _, oldItem := range old {
		oldRule = append(oldRule, oldItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _AdaptiveOracle.contract.FilterLogs(opts, "UnderlyingRateFeedSet", oldRule, currentRule)
	if err != nil {
		return nil, err
	}
	return &AdaptiveOracleUnderlyingRateFeedSetIterator{contract: _AdaptiveOracle.contract, event: "UnderlyingRateFeedSet", logs: logs, sub: sub}, nil
}

// WatchUnderlyingRateFeedSet is a free log subscription operation binding the contract event 0x587234defe9e4ba3e8d3b8b6ad1523bfb568b8f573392aec0a7a92a1b298407f.
//
// Solidity: event UnderlyingRateFeedSet(address indexed old, address indexed current)
func (_AdaptiveOracle *AdaptiveOracleFilterer) WatchUnderlyingRateFeedSet(opts *bind.WatchOpts, sink chan<- *AdaptiveOracleUnderlyingRateFeedSet, old []common.Address, current []common.Address) (event.Subscription, error) {

	var oldRule []interface{}
	for _, oldItem := range old {
		oldRule = append(oldRule, oldItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _AdaptiveOracle.contract.WatchLogs(opts, "UnderlyingRateFeedSet", oldRule, currentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AdaptiveOracleUnderlyingRateFeedSet)
				if err := _AdaptiveOracle.contract.UnpackLog(event, "UnderlyingRateFeedSet", log); err != nil {
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

// ParseUnderlyingRateFeedSet is a log parse operation binding the contract event 0x587234defe9e4ba3e8d3b8b6ad1523bfb568b8f573392aec0a7a92a1b298407f.
//
// Solidity: event UnderlyingRateFeedSet(address indexed old, address indexed current)
func (_AdaptiveOracle *AdaptiveOracleFilterer) ParseUnderlyingRateFeedSet(log types.Log) (*AdaptiveOracleUnderlyingRateFeedSet, error) {
	event := new(AdaptiveOracleUnderlyingRateFeedSet)
	if err := _AdaptiveOracle.contract.UnpackLog(event, "UnderlyingRateFeedSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
