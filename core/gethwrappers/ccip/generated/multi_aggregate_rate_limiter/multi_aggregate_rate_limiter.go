// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package multi_aggregate_rate_limiter

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

type AuthorizedCallersAuthorizedCallerArgs struct {
	AddedCallers   []common.Address
	RemovedCallers []common.Address
}

type ClientAny2EVMMessage struct {
	MessageId           [32]byte
	SourceChainSelector uint64
	Sender              []byte
	Data                []byte
	DestTokenAmounts    []ClientEVMTokenAmount
}

type ClientEVM2AnyMessage struct {
	Receiver     []byte
	Data         []byte
	TokenAmounts []ClientEVMTokenAmount
	FeeToken     common.Address
	ExtraArgs    []byte
}

type ClientEVMTokenAmount struct {
	Token  common.Address
	Amount *big.Int
}

type MultiAggregateRateLimiterLocalRateLimitToken struct {
	RemoteChainSelector uint64
	LocalToken          common.Address
}

type MultiAggregateRateLimiterRateLimitTokenArgs struct {
	LocalTokenArgs MultiAggregateRateLimiterLocalRateLimitToken
	RemoteToken    []byte
}

type MultiAggregateRateLimiterRateLimiterConfigArgs struct {
	RemoteChainSelector uint64
	IsOutboundLane      bool
	RateLimiterConfig   RateLimiterConfig
}

type RateLimiterConfig struct {
	IsEnabled bool
	Capacity  *big.Int
	Rate      *big.Int
}

type RateLimiterTokenBucket struct {
	Tokens      *big.Int
	LastUpdated uint32
	IsEnabled   bool
	Capacity    *big.Int
	Rate        *big.Int
}

var MultiAggregateRateLimiterMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"feeQuoter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"authorizedCallers\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"applyAuthorizedCallerUpdates\",\"inputs\":[{\"name\":\"authorizedCallerArgs\",\"type\":\"tuple\",\"internalType\":\"structAuthorizedCallers.AuthorizedCallerArgs\",\"components\":[{\"name\":\"addedCallers\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"removedCallers\",\"type\":\"address[]\",\"internalType\":\"address[]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"applyRateLimiterConfigUpdates\",\"inputs\":[{\"name\":\"rateLimiterUpdates\",\"type\":\"tuple[]\",\"internalType\":\"structMultiAggregateRateLimiter.RateLimiterConfigArgs[]\",\"components\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"isOutboundLane\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"rateLimiterConfig\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"currentRateLimiterState\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"isOutboundLane\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.TokenBucket\",\"components\":[{\"name\":\"tokens\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"lastUpdated\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllAuthorizedCallers\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllRateLimitTokens\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"localTokens\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"remoteTokens\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFeeQuoter\",\"inputs\":[],\"outputs\":[{\"name\":\"feeQuoter\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"onInboundMessage\",\"inputs\":[{\"name\":\"message\",\"type\":\"tuple\",\"internalType\":\"structClient.Any2EVMMessage\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destTokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structClient.EVMTokenAmount[]\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"onOutboundMessage\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"message\",\"type\":\"tuple\",\"internalType\":\"structClient.EVM2AnyMessage\",\"components\":[{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"tokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structClient.EVMTokenAmount[]\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"feeToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"extraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setFeeQuoter\",\"inputs\":[{\"name\":\"newFeeQuoter\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"typeAndVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"updateRateLimitTokens\",\"inputs\":[{\"name\":\"removes\",\"type\":\"tuple[]\",\"internalType\":\"structMultiAggregateRateLimiter.LocalRateLimitToken[]\",\"components\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"localToken\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"name\":\"adds\",\"type\":\"tuple[]\",\"internalType\":\"structMultiAggregateRateLimiter.RateLimitTokenArgs[]\",\"components\":[{\"name\":\"localTokenArgs\",\"type\":\"tuple\",\"internalType\":\"structMultiAggregateRateLimiter.LocalRateLimitToken\",\"components\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"localToken\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"name\":\"remoteToken\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AuthorizedCallerAdded\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AuthorizedCallerRemoved\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ConfigChanged\",\"inputs\":[{\"name\":\"config\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FeeQuoterSet\",\"inputs\":[{\"name\":\"newFeeQuoter\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RateLimiterConfigUpdated\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"isOutboundLane\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"},{\"name\":\"config\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TokenAggregateRateLimitAdded\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"remoteToken\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"localToken\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TokenAggregateRateLimitRemoved\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"localToken\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TokensConsumed\",\"inputs\":[{\"name\":\"tokens\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AggregateValueMaxCapacityExceeded\",\"inputs\":[{\"name\":\"capacity\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"requested\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"AggregateValueRateLimitReached\",\"inputs\":[{\"name\":\"minWaitInSeconds\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"available\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"BucketOverfilled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CannotTransferToSelf\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DisabledNonZeroRateLimit\",\"inputs\":[{\"name\":\"config\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}]},{\"type\":\"error\",\"name\":\"InvalidRateLimitRate\",\"inputs\":[{\"name\":\"rateLimiterConfig\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}]},{\"type\":\"error\",\"name\":\"MessageValidationError\",\"inputs\":[{\"name\":\"errorReason\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"MustBeProposedOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnerCannotBeZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PriceNotFoundForToken\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RateLimitMustBeDisabled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TokenMaxCapacityExceeded\",\"inputs\":[{\"name\":\"capacity\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"requested\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"TokenRateLimitReached\",\"inputs\":[{\"name\":\"minWaitInSeconds\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"available\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UnauthorizedCaller\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ZeroAddressNotAllowed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroChainSelectorNotAllowed\",\"inputs\":[]}]",
	Bin: "0x60806040523461026457612f2d8038038061001981610269565b928339810190604081830312610264576100328161028e565b602082015190916001600160401b03821161026457019180601f84011215610264578251926001600160401b038411610225578360051b90602080610078818501610269565b80978152019282010192831161026457602001905b82821061024c57505050331561023b57600180546001600160a01b0319163317905560206100ba81610269565b60008152600036813760408051949085016001600160401b03811186821017610225576040528452808285015260005b8151811015610151576001906001600160a01b0361010882856102a2565b511684610114826102e4565b610121575b5050016100ea565b7fc3803387881faad271c47728894e3e36fac830ffc8602ca6fc07733cbda7758091604051908152a13884610119565b5050915160005b81518110156101c9576001600160a01b0361017382846102a2565b51169081156101b8577feb1b9b92e50b7f88f9ff25d56765095ac6e91540eee214906f4036a908ffbdef85836101aa6001956103e2565b50604051908152a101610158565b6342bcdf7f60e11b60005260046000fd5b50506001600160a01b03169081156101b857600580546001600160a01b031916831790556040519182527f7c737a8eddf62436489aa3600ed26e75e0a58b0f8c0d266bbcee64358c39fdac91a1604051612aea90816104438239f35b634e487b7160e01b600052604160045260246000fd5b639b15e16f60e01b60005260046000fd5b602080916102598461028e565b81520191019061008d565b600080fd5b6040519190601f01601f191682016001600160401b0381118382101761022557604052565b51906001600160a01b038216820361026457565b80518210156102b65760209160051b010190565b634e487b7160e01b600052603260045260246000fd5b80548210156102b65760005260206000200190600090565b60008181526003602052604090205480156103db5760001981018181116103c5576002546000198101919082116103c557808203610374575b505050600254801561035e57600019016103388160026102cc565b8154906000199060031b1b19169055600255600052600360205260006040812055600190565b634e487b7160e01b600052603160045260246000fd5b6103ad6103856103969360026102cc565b90549060031b1c92839260026102cc565b819391549060031b91821b91600019901b19161790565b9055600052600360205260406000205538808061031d565b634e487b7160e01b600052601160045260246000fd5b5050600090565b8060005260036020526040600020541560001461043c57600254680100000000000000008110156102255761042361039682600185940160025560026102cc565b9055600254906000526003602052604060002055600190565b5060009056fe608080604052600436101561001357600080fd5b60003560e01c90816308d450a114611c7c575080630a35bcc414611b4b578063181f5a7714611ace5780631af18b7b1461156c5780632451a627146114bf578063537e304e146111f557806379ba50971461110c5780638da5cb5b146110ba57806391a2749a14610efc578063e0a0e50614610bbb578063e145291614610b69578063e835232b14610a8d578063f2fde38b1461099d5763fe843cd0146100b957600080fd5b346109985760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126109985760043567ffffffffffffffff81116109985736602382011215610998578060040135610113816120b3565b916101216040519384611fd1565b818352602460a06020850193028201019036821161099857602401915b8183106108ea578361014e612416565b6000905b80518210156108e857610165828261231e565b519160408301519267ffffffffffffffff8151169081156108be5760200151151561019081836123e1565b805463ffffffff8160801c16801560001461066757505085516000915015610592576fffffffffffffffffffffffffffffffff6040870151166fffffffffffffffffffffffffffffffff602088015116811090811591610589575b50610526577ff14a5415ce6988a9e870a85fff0b9d7b7dd79bbc228cb63cad610daf6f7b6b979161042960019697608093505b6fffffffffffffffffffffffffffffffff6040820151166fffffffffffffffffffffffffffffffff602083015116825115159160405161025d81611f7d565b828152602081019363ffffffff4216855260408201908152606082019384528882019283528a886000146104315760036103ef966103816fffffffffffffffffffffffffffffffff969587958695600052600660205261033a63ffffffff604060002095888060028901965116167fffffffffffffffffffffffffffffffff00000000000000000000000000000000865416178555511683907fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff73ffffffff0000000000000000000000000000000083549260801b169116179055565b5181547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff1690151560a01b74ff000000000000000000000000000000000000000016179055565b01945116167fffffffffffffffffffffffffffffffff0000000000000000000000000000000084541617835551166fffffffffffffffffffffffffffffffff7fffffffffffffffffffffffffffffffff0000000000000000000000000000000083549260801b169116179055565b60405192835260208301906fffffffffffffffffffffffffffffffff60408092805115158552826020820151166020860152015116910152565ba20190610152565b8d6fffffffffffffffffffffffffffffffff949361038186946104da63ffffffff6105219b8897600052600660205287806040600020975116167fffffffffffffffffffffffffffffffff00000000000000000000000000000000875416178655511684907fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff73ffffffff0000000000000000000000000000000083549260801b169116179055565b5182547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff1690151560a01b74ff000000000000000000000000000000000000000016178255565b6103ef565b606486610587604051917f8020d12400000000000000000000000000000000000000000000000000000000835260048301906fffffffffffffffffffffffffffffffff60408092805115158552826020820151166020860152015116910152565bfd5b905015876101eb565b506fffffffffffffffffffffffffffffffff60408601511615801590610648575b6105e75760807ff14a5415ce6988a9e870a85fff0b9d7b7dd79bbc228cb63cad610daf6f7b6b97916104296001969761021e565b606485610587604051917fd68af9cc00000000000000000000000000000000000000000000000000000000835260048301906fffffffffffffffffffffffffffffffff60408092805115158552826020820151166020860152015116910152565b506fffffffffffffffffffffffffffffffff60208601511615156105b3565b6001969761079c6080947ff14a5415ce6988a9e870a85fff0b9d7b7dd79bbc228cb63cad610daf6f7b6b9796946106a16104299542612409565b9081610804575b50506fffffffffffffffffffffffffffffffff8a8160208601511692828154168085106000146107fc57508280855b16167fffffffffffffffffffffffffffffffff000000000000000000000000000000008254161781556107508651151582907fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff74ff0000000000000000000000000000000000000000835492151560a01b169116179055565b60408601517fffffffffffffffffffffffffffffffff0000000000000000000000000000000060809190911b16939092166fffffffffffffffffffffffffffffffff1692909217910155565b7f9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c1960606040516107f681856fffffffffffffffffffffffffffffffff60408092805115158552826020820151166020860152015116910152565ba16103ef565b8380916106d7565b6fffffffffffffffffffffffffffffffff916108388392838f6108319088015494828616958e1c906126c9565b91166123a5565b808210156108b757505b83547fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff9290911692909216167fffffffffffffffffffffffff0000000000000000000000000000000000000000909116174260801b73ffffffff00000000000000000000000000000000161781558b806106a8565b9050610842565b7fc65608950000000000000000000000000000000000000000000000000000000060005260046000fd5b005b82360360a081126109985760607fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc06040519261092584611fb5565b61092e87612029565b845261093c60208801612185565b602085015201126109985760a09160209160405161095981611fb5565b61096560408801612185565b8152610973606088016122d6565b84820152610983608088016122d6565b6040820152604082015281520192019161013e565b600080fd5b346109985760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126109985773ffffffffffffffffffffffffffffffffffffffff6109e96120cb565b6109f1612416565b16338114610a6357807fffffffffffffffffffffffff0000000000000000000000000000000000000000600054161760005573ffffffffffffffffffffffffffffffffffffffff600154167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278600080a3005b7fdad89dca0000000000000000000000000000000000000000000000000000000060005260046000fd5b346109985760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126109985773ffffffffffffffffffffffffffffffffffffffff610ad96120cb565b610ae1612416565b168015610b3f576020817f7c737a8eddf62436489aa3600ed26e75e0a58b0f8c0d266bbcee64358c39fdac927fffffffffffffffffffffffff00000000000000000000000000000000000000006005541617600555604051908152a1005b7f8579befe0000000000000000000000000000000000000000000000000000000060005260046000fd5b346109985760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261099857602073ffffffffffffffffffffffffffffffffffffffff60055416604051908152f35b346109985760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261099857610bf2612012565b60243567ffffffffffffffff8111610998578036039060a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc83011261099857610c3a612361565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdd60448201359201821215610998570160048101359067ffffffffffffffff8211610998576024018160061b3603811361099857610c9991369161210f565b90610ca56001826123e1565b9160ff835460a01c16610cb457005b6000929167ffffffffffffffff16835b8251811015610ee857816000526004602052610d16604060002073ffffffffffffffffffffffffffffffffffffffff610cfd848761231e565b5151169060019160005201602052604060002054151590565b610d23575b600101610cc4565b93610d2e858461231e565b5173ffffffffffffffffffffffffffffffffffffffff60055416604073ffffffffffffffffffffffffffffffffffffffff83511660248251809481937fd02641a000000000000000000000000000000000000000000000000000000000835260048301525afa8015610edc57600090610e3c575b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff91505116908115610df75791670de0b6b3a7640000610de8610def9360206001960151906126c9565b04906123a5565b949050610d1b565b73ffffffffffffffffffffffffffffffffffffffff9051167f9a655f7b0000000000000000000000000000000000000000000000000000000060005260045260246000fd5b6040823d8211610ed4575b81610e5460409383611fd1565b81010312610ecd5760405191610e6983611f99565b80517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff81168103610ed0578352602001519063ffffffff82168203610ecd575060208201527bffffffffffffffffffffffffffffffffffffffffffffffffffffffff90610da2565b80fd5b8280fd5b3d9150610e47565b6040513d6000823e3d90fd5b5050509080610ef357005b6108e891612461565b346109985760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126109985760043567ffffffffffffffff81116109985760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc82360301126109985760405190610f7682611f99565b806004013567ffffffffffffffff811161099857610f9a9060043691840101612271565b825260248101359067ffffffffffffffff8211610998576004610fc09236920101612271565b60208201908152610fcf612416565b519060005b8251811015611047578073ffffffffffffffffffffffffffffffffffffffff610fff6001938661231e565b511661100a8161275e565b611016575b5001610fd4565b60207fc3803387881faad271c47728894e3e36fac830ffc8602ca6fc07733cbda7758091604051908152a18461100f565b505160005b81518110156108e85773ffffffffffffffffffffffffffffffffffffffff611074828461231e565b5116908115610b3f577feb1b9b92e50b7f88f9ff25d56765095ac6e91540eee214906f4036a908ffbdef6020836110ac600195612a28565b50604051908152a10161104c565b346109985760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261099857602073ffffffffffffffffffffffffffffffffffffffff60015416604051908152f35b346109985760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126109985760005473ffffffffffffffffffffffffffffffffffffffff811633036111cb577fffffffffffffffffffffffff00000000000000000000000000000000000000006001549133828416176001551660005573ffffffffffffffffffffffffffffffffffffffff3391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0600080a3005b7f02b543c60000000000000000000000000000000000000000000000000000000060005260046000fd5b346109985760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126109985767ffffffffffffffff611235612012565b168060005260046020526040600020549061124f826120b3565b9161125d6040519384611fd1565b8083527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe061128a826120b3565b0136602085013761129a816120b3565b916112a86040519384611fd1565b8183527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe06112d5836120b3565b0160005b8181106114ae57505060005b828110611384576113058585604051928392604084526040840190612227565b8281036020840152815180825260208201916020808360051b8301019401926000915b8383106113355786860387f35b919395509193602080611372837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe086600196030187528951612192565b97019301930190928695949293611328565b8160005260046020526040600020600261139e8383612746565b90549060031b1c918260005201602052604060002090604051916000908054906113c7826126dc565b8086529160018116908115611469575060011461142f575b505082916114076001959473ffffffffffffffffffffffffffffffffffffffff930384611fd1565b16611412838961231e565b5261141d828761231e565b52611428818661231e565b50016112e5565b6000908152602081209092505b8183106114535750508201602001816114076113df565b600181602092548386890101520192019161143c565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660208088019190915292151560051b8601909201925083915061140790506113df565b8060606020809388010152016112d9565b346109985760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126109985760405180602060025491828152019060026000527f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace9060005b818110611556576115528561153e81870382611fd1565b604051918291602083526020830190612227565b0390f35b8254845260209093019260019283019201611527565b346109985760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126109985760043567ffffffffffffffff81116109985736602382011215610998578060040135906115c7826120b3565b916115d56040519384611fd1565b8083526024602084019160061b8301019136831161099857602401905b828210611ab4576024358467ffffffffffffffff82116109985736602383011215610998578160040135611625816120b3565b926116336040519485611fd1565b8184526024602085019260051b820101903682116109985760248101925b828410611a21578585611662612416565b60005b815181101561176d578073ffffffffffffffffffffffffffffffffffffffff60206116926001948661231e565b5101511667ffffffffffffffff6116a9838661231e565b515116908160005260046020526116e7816040600020816000526002810160205260406000206116d981546126dc565b908161172b575b50506128f4565b6116f4575b505001611665565b7f530cabd30786b7235e124a6c0db77e0b685ef22813b1fe87554247f404eb8ed69160409182519182526020820152a184806116ec565b81601f600093118a146117425750555b89806116e0565b8183526020832061175d91601f0160051c8101908b0161272f565b808252816020812091555561173b565b8260005b81518110156108e857611784818361231e565b51516020611792838561231e565b51015173ffffffffffffffffffffffffffffffffffffffff6020830151169182158015611a18575b80156119ec575b610b3f5767ffffffffffffffff905116806000526004602052604060002083600052600281016020526040600020835167ffffffffffffffff81116119bd5761180a82546126dc565b601f8111611980575b506020601f82116001146118d5579261186d9282889796959360019a99946000916118ca575b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff828c1b9260031b1c1916179055612a88565b61187b575b50505001611771565b7fad72a792d2a307f400c278be7deaeec6964276783304580cdc4e905436b8d5c5926118b96040519384938452606060208501526060840190612192565b9060408301520390a1838080611872565b90508701518c611839565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe082169083600052806000209160005b818110611968575083899897969361186d969360019c9b968d9410611931575b5050811b019055612a88565b8901517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690558c80611925565b9192602060018192868c015181550194019201611905565b6119ad90836000526020600020601f840160051c810191602085106119b3575b601f0160051c019061272f565b88611813565b90915081906119a0565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b508151602083012060405160208101906000825260208152611a0f604082611fd1565b519020146117c1565b508151156117ba565b833567ffffffffffffffff811161099857820160607fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdc82360301126109985760405191611a6d83611f99565b611a7a36602484016121f1565b835260648201359267ffffffffffffffff841161099857611aa560209493602486953692010161203e565b83820152815201930192611651565b6020604091611ac336856121f1565b8152019101906115f2565b346109985760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610998576115526040805190611b0f8183611fd1565b601f82527f4d756c7469416767726567617465526174654c696d6974657220312e362e3000602083015251918291602083526020830190612192565b346109985760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261099857611b82612012565b6024359081151582036109985760a091611ba491611b9e6122f3565b506123e1565b6fffffffffffffffffffffffffffffffff60405191611bc283611f7d565b8181549181831685526001602086019163ffffffff8560801c16835260ff6040880195891c161515855201549263ffffffff60608701928486168452608088019560801c8652611c106122f3565b508480855116611c3d828b5116611c37611c2d8787511642612409565b858c5116906126c9565b906123a5565b80821015611c7557505b1680985281421681526040519788525116602087015251151560408601525116606084015251166080820152f35b9050611c47565b346109985760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126109985760043567ffffffffffffffff81116109985760a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc823603011261099857611cf282611f7d565b80600401358252611d0560248201612029565b9060208301918252604481013567ffffffffffffffff811161099857611d31906004369184010161203e565b6040840152606481013567ffffffffffffffff811161099857611d5a906004369184010161203e565b606084015260848101359067ffffffffffffffff821161099857019036602383011215610998576080611da067ffffffffffffffff93369060246004820135910161210f565b9301928352611dad612361565b51169051611dbc6000836123e1565b9060ff825460a01c16611dcb57005b60009260005b8251811015610ee857816000526004602052611e0a604060002073ffffffffffffffffffffffffffffffffffffffff610cfd848761231e565b611e17575b600101611dd1565b93611e22858461231e565b5173ffffffffffffffffffffffffffffffffffffffff60055416604073ffffffffffffffffffffffffffffffffffffffff83511660248251809481937fd02641a000000000000000000000000000000000000000000000000000000000835260048301525afa8015610edc57600090611ee4575b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff91505116908115610df75791670de0b6b3a7640000610de8611edc9360206001960151906126c9565b949050611e0f565b6040823d8211611f75575b81611efc60409383611fd1565b81010312610ecd5760405191611f1183611f99565b80517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff81168103610ed0578352602001519063ffffffff82168203610ecd575060208201527bffffffffffffffffffffffffffffffffffffffffffffffffffffffff90611e96565b3d9150611eef565b60a0810190811067ffffffffffffffff8211176119bd57604052565b6040810190811067ffffffffffffffff8211176119bd57604052565b6060810190811067ffffffffffffffff8211176119bd57604052565b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff8211176119bd57604052565b6004359067ffffffffffffffff8216820361099857565b359067ffffffffffffffff8216820361099857565b81601f820112156109985780359067ffffffffffffffff82116119bd5760405192612091601f84017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200185611fd1565b8284526020838301011161099857816000926020809301838601378301015290565b67ffffffffffffffff81116119bd5760051b60200190565b6004359073ffffffffffffffffffffffffffffffffffffffff8216820361099857565b359073ffffffffffffffffffffffffffffffffffffffff8216820361099857565b92919261211b826120b3565b936121296040519586611fd1565b602085848152019260061b82019181831161099857925b82841061214d5750505050565b604084830312610998576020604091825161216781611f99565b612170876120ee565b81528287013583820152815201930192612140565b3590811515820361099857565b919082519283825260005b8481106121dc5750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8460006020809697860101520116010190565b8060208092840101518282860101520161219d565b91908260409103126109985760405161220981611f99565b602061222281839561221a81612029565b8552016120ee565b910152565b906020808351928381520192019060005b8181106122455750505090565b825173ffffffffffffffffffffffffffffffffffffffff16845260209384019390920191600101612238565b9080601f83011215610998578135612288816120b3565b926122966040519485611fd1565b81845260208085019260051b82010192831161099857602001905b8282106122be5750505090565b602080916122cb846120ee565b8152019101906122b1565b35906fffffffffffffffffffffffffffffffff8216820361099857565b6040519061230082611f7d565b60006080838281528260208201528260408201528260608201520152565b80518210156123325760209160051b010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b3360005260036020526040600020541561237757565b7fd86ad9cf000000000000000000000000000000000000000000000000000000006000523360045260246000fd5b919082018092116123b257565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b67ffffffffffffffff1660005260066020526040600020906000146124065760020190565b90565b919082039182116123b257565b73ffffffffffffffffffffffffffffffffffffffff60015416330361243757565b7f2b5c74de0000000000000000000000000000000000000000000000000000000060005260046000fd5b805460ff8160a01c161580156126c1575b6126bc576fffffffffffffffffffffffffffffffff811690600183019081546124b763ffffffff6fffffffffffffffffffffffffffffffff83169360801c1642612409565b908161261e575b50508481106125ec57508382106125435750916020916fffffffffffffffffffffffffffffffff80612511847f1871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690a97612409565b16167fffffffffffffffffffffffffffffffff00000000000000000000000000000000825416179055604051908152a1565b5460801c6125518285612409565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82018281116123b257612584916123a5565b9080156125bd5790047f15279c080000000000000000000000000000000000000000000000000000000060005260045260245260446000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b84907ff94ebcd10000000000000000000000000000000000000000000000000000000060005260045260245260446000fd5b8285929395116126925761263992611c379160801c906126c9565b8083101561268d5750815b83547fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff164260801b73ffffffff00000000000000000000000000000000161784559138806124be565b612644565b7f9725942a0000000000000000000000000000000000000000000000000000000060005260046000fd5b505050565b508215612472565b818102929181159184041417156123b257565b90600182811c92168015612725575b60208310146126f657565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b91607f16916126eb565b81811061273a575050565b6000815560010161272f565b80548210156123325760005260206000200190600090565b60008181526003602052604090205480156128ed577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81018181116123b257600254907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82019182116123b25780820361287e575b505050600254801561284f577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0161280c816002612746565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82549160031b1b19169055600255600052600360205260006040812055600190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b6128d561288f6128a0936002612746565b90549060031b1c9283926002612746565b81939154907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9060031b92831b921b19161790565b905560005260036020526040600020553880806127d3565b5050600090565b9060018201918160005282602052604060002054801515600014612a1f577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81018181116123b2578254907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82019182116123b2578082036129e8575b5050508054801561284f577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01906129a98282612746565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82549160031b1b191690555560005260205260006040812055600190565b612a086129f86128a09386612746565b90549060031b1c92839286612746565b905560005283602052604060002055388080612971565b50505050600090565b80600052600360205260406000205415600014612a8257600254680100000000000000008110156119bd57612a696128a08260018594016002556002612746565b9055600254906000526003602052604060002055600190565b50600090565b60008281526001820160205260409020546128ed57805490680100000000000000008210156119bd5782612ac66128a0846001809601855584612746565b90558054926000520160205260406000205560019056fea164736f6c634300081a000a",
}

var MultiAggregateRateLimiterABI = MultiAggregateRateLimiterMetaData.ABI

var MultiAggregateRateLimiterBin = MultiAggregateRateLimiterMetaData.Bin

func DeployMultiAggregateRateLimiter(auth *bind.TransactOpts, backend bind.ContractBackend, feeQuoter common.Address, authorizedCallers []common.Address) (common.Address, *types.Transaction, *MultiAggregateRateLimiter, error) {
	parsed, err := MultiAggregateRateLimiterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MultiAggregateRateLimiterBin), backend, feeQuoter, authorizedCallers)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MultiAggregateRateLimiter{address: address, abi: *parsed, MultiAggregateRateLimiterCaller: MultiAggregateRateLimiterCaller{contract: contract}, MultiAggregateRateLimiterTransactor: MultiAggregateRateLimiterTransactor{contract: contract}, MultiAggregateRateLimiterFilterer: MultiAggregateRateLimiterFilterer{contract: contract}}, nil
}

type MultiAggregateRateLimiter struct {
	address common.Address
	abi     abi.ABI
	MultiAggregateRateLimiterCaller
	MultiAggregateRateLimiterTransactor
	MultiAggregateRateLimiterFilterer
}

type MultiAggregateRateLimiterCaller struct {
	contract *bind.BoundContract
}

type MultiAggregateRateLimiterTransactor struct {
	contract *bind.BoundContract
}

type MultiAggregateRateLimiterFilterer struct {
	contract *bind.BoundContract
}

type MultiAggregateRateLimiterSession struct {
	Contract     *MultiAggregateRateLimiter
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type MultiAggregateRateLimiterCallerSession struct {
	Contract *MultiAggregateRateLimiterCaller
	CallOpts bind.CallOpts
}

type MultiAggregateRateLimiterTransactorSession struct {
	Contract     *MultiAggregateRateLimiterTransactor
	TransactOpts bind.TransactOpts
}

type MultiAggregateRateLimiterRaw struct {
	Contract *MultiAggregateRateLimiter
}

type MultiAggregateRateLimiterCallerRaw struct {
	Contract *MultiAggregateRateLimiterCaller
}

type MultiAggregateRateLimiterTransactorRaw struct {
	Contract *MultiAggregateRateLimiterTransactor
}

func NewMultiAggregateRateLimiter(address common.Address, backend bind.ContractBackend) (*MultiAggregateRateLimiter, error) {
	abi, err := abi.JSON(strings.NewReader(MultiAggregateRateLimiterABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindMultiAggregateRateLimiter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MultiAggregateRateLimiter{address: address, abi: abi, MultiAggregateRateLimiterCaller: MultiAggregateRateLimiterCaller{contract: contract}, MultiAggregateRateLimiterTransactor: MultiAggregateRateLimiterTransactor{contract: contract}, MultiAggregateRateLimiterFilterer: MultiAggregateRateLimiterFilterer{contract: contract}}, nil
}

func NewMultiAggregateRateLimiterCaller(address common.Address, caller bind.ContractCaller) (*MultiAggregateRateLimiterCaller, error) {
	contract, err := bindMultiAggregateRateLimiter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MultiAggregateRateLimiterCaller{contract: contract}, nil
}

func NewMultiAggregateRateLimiterTransactor(address common.Address, transactor bind.ContractTransactor) (*MultiAggregateRateLimiterTransactor, error) {
	contract, err := bindMultiAggregateRateLimiter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MultiAggregateRateLimiterTransactor{contract: contract}, nil
}

func NewMultiAggregateRateLimiterFilterer(address common.Address, filterer bind.ContractFilterer) (*MultiAggregateRateLimiterFilterer, error) {
	contract, err := bindMultiAggregateRateLimiter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MultiAggregateRateLimiterFilterer{contract: contract}, nil
}

func bindMultiAggregateRateLimiter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MultiAggregateRateLimiterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MultiAggregateRateLimiter.Contract.MultiAggregateRateLimiterCaller.contract.Call(opts, result, method, params...)
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultiAggregateRateLimiter.Contract.MultiAggregateRateLimiterTransactor.contract.Transfer(opts)
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MultiAggregateRateLimiter.Contract.MultiAggregateRateLimiterTransactor.contract.Transact(opts, method, params...)
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MultiAggregateRateLimiter.Contract.contract.Call(opts, result, method, params...)
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultiAggregateRateLimiter.Contract.contract.Transfer(opts)
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MultiAggregateRateLimiter.Contract.contract.Transact(opts, method, params...)
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterCaller) CurrentRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64, isOutboundLane bool) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _MultiAggregateRateLimiter.contract.Call(opts, &out, "currentRateLimiterState", remoteChainSelector, isOutboundLane)

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterSession) CurrentRateLimiterState(remoteChainSelector uint64, isOutboundLane bool) (RateLimiterTokenBucket, error) {
	return _MultiAggregateRateLimiter.Contract.CurrentRateLimiterState(&_MultiAggregateRateLimiter.CallOpts, remoteChainSelector, isOutboundLane)
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterCallerSession) CurrentRateLimiterState(remoteChainSelector uint64, isOutboundLane bool) (RateLimiterTokenBucket, error) {
	return _MultiAggregateRateLimiter.Contract.CurrentRateLimiterState(&_MultiAggregateRateLimiter.CallOpts, remoteChainSelector, isOutboundLane)
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterCaller) GetAllAuthorizedCallers(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _MultiAggregateRateLimiter.contract.Call(opts, &out, "getAllAuthorizedCallers")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterSession) GetAllAuthorizedCallers() ([]common.Address, error) {
	return _MultiAggregateRateLimiter.Contract.GetAllAuthorizedCallers(&_MultiAggregateRateLimiter.CallOpts)
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterCallerSession) GetAllAuthorizedCallers() ([]common.Address, error) {
	return _MultiAggregateRateLimiter.Contract.GetAllAuthorizedCallers(&_MultiAggregateRateLimiter.CallOpts)
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterCaller) GetAllRateLimitTokens(opts *bind.CallOpts, remoteChainSelector uint64) (GetAllRateLimitTokens,

	error) {
	var out []interface{}
	err := _MultiAggregateRateLimiter.contract.Call(opts, &out, "getAllRateLimitTokens", remoteChainSelector)

	outstruct := new(GetAllRateLimitTokens)
	if err != nil {
		return *outstruct, err
	}

	outstruct.LocalTokens = *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	outstruct.RemoteTokens = *abi.ConvertType(out[1], new([][]byte)).(*[][]byte)

	return *outstruct, err

}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterSession) GetAllRateLimitTokens(remoteChainSelector uint64) (GetAllRateLimitTokens,

	error) {
	return _MultiAggregateRateLimiter.Contract.GetAllRateLimitTokens(&_MultiAggregateRateLimiter.CallOpts, remoteChainSelector)
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterCallerSession) GetAllRateLimitTokens(remoteChainSelector uint64) (GetAllRateLimitTokens,

	error) {
	return _MultiAggregateRateLimiter.Contract.GetAllRateLimitTokens(&_MultiAggregateRateLimiter.CallOpts, remoteChainSelector)
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterCaller) GetFeeQuoter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MultiAggregateRateLimiter.contract.Call(opts, &out, "getFeeQuoter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterSession) GetFeeQuoter() (common.Address, error) {
	return _MultiAggregateRateLimiter.Contract.GetFeeQuoter(&_MultiAggregateRateLimiter.CallOpts)
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterCallerSession) GetFeeQuoter() (common.Address, error) {
	return _MultiAggregateRateLimiter.Contract.GetFeeQuoter(&_MultiAggregateRateLimiter.CallOpts)
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MultiAggregateRateLimiter.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterSession) Owner() (common.Address, error) {
	return _MultiAggregateRateLimiter.Contract.Owner(&_MultiAggregateRateLimiter.CallOpts)
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterCallerSession) Owner() (common.Address, error) {
	return _MultiAggregateRateLimiter.Contract.Owner(&_MultiAggregateRateLimiter.CallOpts)
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MultiAggregateRateLimiter.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterSession) TypeAndVersion() (string, error) {
	return _MultiAggregateRateLimiter.Contract.TypeAndVersion(&_MultiAggregateRateLimiter.CallOpts)
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterCallerSession) TypeAndVersion() (string, error) {
	return _MultiAggregateRateLimiter.Contract.TypeAndVersion(&_MultiAggregateRateLimiter.CallOpts)
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultiAggregateRateLimiter.contract.Transact(opts, "acceptOwnership")
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterSession) AcceptOwnership() (*types.Transaction, error) {
	return _MultiAggregateRateLimiter.Contract.AcceptOwnership(&_MultiAggregateRateLimiter.TransactOpts)
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _MultiAggregateRateLimiter.Contract.AcceptOwnership(&_MultiAggregateRateLimiter.TransactOpts)
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterTransactor) ApplyAuthorizedCallerUpdates(opts *bind.TransactOpts, authorizedCallerArgs AuthorizedCallersAuthorizedCallerArgs) (*types.Transaction, error) {
	return _MultiAggregateRateLimiter.contract.Transact(opts, "applyAuthorizedCallerUpdates", authorizedCallerArgs)
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterSession) ApplyAuthorizedCallerUpdates(authorizedCallerArgs AuthorizedCallersAuthorizedCallerArgs) (*types.Transaction, error) {
	return _MultiAggregateRateLimiter.Contract.ApplyAuthorizedCallerUpdates(&_MultiAggregateRateLimiter.TransactOpts, authorizedCallerArgs)
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterTransactorSession) ApplyAuthorizedCallerUpdates(authorizedCallerArgs AuthorizedCallersAuthorizedCallerArgs) (*types.Transaction, error) {
	return _MultiAggregateRateLimiter.Contract.ApplyAuthorizedCallerUpdates(&_MultiAggregateRateLimiter.TransactOpts, authorizedCallerArgs)
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterTransactor) ApplyRateLimiterConfigUpdates(opts *bind.TransactOpts, rateLimiterUpdates []MultiAggregateRateLimiterRateLimiterConfigArgs) (*types.Transaction, error) {
	return _MultiAggregateRateLimiter.contract.Transact(opts, "applyRateLimiterConfigUpdates", rateLimiterUpdates)
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterSession) ApplyRateLimiterConfigUpdates(rateLimiterUpdates []MultiAggregateRateLimiterRateLimiterConfigArgs) (*types.Transaction, error) {
	return _MultiAggregateRateLimiter.Contract.ApplyRateLimiterConfigUpdates(&_MultiAggregateRateLimiter.TransactOpts, rateLimiterUpdates)
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterTransactorSession) ApplyRateLimiterConfigUpdates(rateLimiterUpdates []MultiAggregateRateLimiterRateLimiterConfigArgs) (*types.Transaction, error) {
	return _MultiAggregateRateLimiter.Contract.ApplyRateLimiterConfigUpdates(&_MultiAggregateRateLimiter.TransactOpts, rateLimiterUpdates)
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterTransactor) OnInboundMessage(opts *bind.TransactOpts, message ClientAny2EVMMessage) (*types.Transaction, error) {
	return _MultiAggregateRateLimiter.contract.Transact(opts, "onInboundMessage", message)
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterSession) OnInboundMessage(message ClientAny2EVMMessage) (*types.Transaction, error) {
	return _MultiAggregateRateLimiter.Contract.OnInboundMessage(&_MultiAggregateRateLimiter.TransactOpts, message)
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterTransactorSession) OnInboundMessage(message ClientAny2EVMMessage) (*types.Transaction, error) {
	return _MultiAggregateRateLimiter.Contract.OnInboundMessage(&_MultiAggregateRateLimiter.TransactOpts, message)
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterTransactor) OnOutboundMessage(opts *bind.TransactOpts, destChainSelector uint64, message ClientEVM2AnyMessage) (*types.Transaction, error) {
	return _MultiAggregateRateLimiter.contract.Transact(opts, "onOutboundMessage", destChainSelector, message)
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterSession) OnOutboundMessage(destChainSelector uint64, message ClientEVM2AnyMessage) (*types.Transaction, error) {
	return _MultiAggregateRateLimiter.Contract.OnOutboundMessage(&_MultiAggregateRateLimiter.TransactOpts, destChainSelector, message)
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterTransactorSession) OnOutboundMessage(destChainSelector uint64, message ClientEVM2AnyMessage) (*types.Transaction, error) {
	return _MultiAggregateRateLimiter.Contract.OnOutboundMessage(&_MultiAggregateRateLimiter.TransactOpts, destChainSelector, message)
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterTransactor) SetFeeQuoter(opts *bind.TransactOpts, newFeeQuoter common.Address) (*types.Transaction, error) {
	return _MultiAggregateRateLimiter.contract.Transact(opts, "setFeeQuoter", newFeeQuoter)
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterSession) SetFeeQuoter(newFeeQuoter common.Address) (*types.Transaction, error) {
	return _MultiAggregateRateLimiter.Contract.SetFeeQuoter(&_MultiAggregateRateLimiter.TransactOpts, newFeeQuoter)
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterTransactorSession) SetFeeQuoter(newFeeQuoter common.Address) (*types.Transaction, error) {
	return _MultiAggregateRateLimiter.Contract.SetFeeQuoter(&_MultiAggregateRateLimiter.TransactOpts, newFeeQuoter)
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _MultiAggregateRateLimiter.contract.Transact(opts, "transferOwnership", to)
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _MultiAggregateRateLimiter.Contract.TransferOwnership(&_MultiAggregateRateLimiter.TransactOpts, to)
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _MultiAggregateRateLimiter.Contract.TransferOwnership(&_MultiAggregateRateLimiter.TransactOpts, to)
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterTransactor) UpdateRateLimitTokens(opts *bind.TransactOpts, removes []MultiAggregateRateLimiterLocalRateLimitToken, adds []MultiAggregateRateLimiterRateLimitTokenArgs) (*types.Transaction, error) {
	return _MultiAggregateRateLimiter.contract.Transact(opts, "updateRateLimitTokens", removes, adds)
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterSession) UpdateRateLimitTokens(removes []MultiAggregateRateLimiterLocalRateLimitToken, adds []MultiAggregateRateLimiterRateLimitTokenArgs) (*types.Transaction, error) {
	return _MultiAggregateRateLimiter.Contract.UpdateRateLimitTokens(&_MultiAggregateRateLimiter.TransactOpts, removes, adds)
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterTransactorSession) UpdateRateLimitTokens(removes []MultiAggregateRateLimiterLocalRateLimitToken, adds []MultiAggregateRateLimiterRateLimitTokenArgs) (*types.Transaction, error) {
	return _MultiAggregateRateLimiter.Contract.UpdateRateLimitTokens(&_MultiAggregateRateLimiter.TransactOpts, removes, adds)
}

type MultiAggregateRateLimiterAuthorizedCallerAddedIterator struct {
	Event *MultiAggregateRateLimiterAuthorizedCallerAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MultiAggregateRateLimiterAuthorizedCallerAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiAggregateRateLimiterAuthorizedCallerAdded)
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
		it.Event = new(MultiAggregateRateLimiterAuthorizedCallerAdded)
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

func (it *MultiAggregateRateLimiterAuthorizedCallerAddedIterator) Error() error {
	return it.fail
}

func (it *MultiAggregateRateLimiterAuthorizedCallerAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MultiAggregateRateLimiterAuthorizedCallerAdded struct {
	Caller common.Address
	Raw    types.Log
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterFilterer) FilterAuthorizedCallerAdded(opts *bind.FilterOpts) (*MultiAggregateRateLimiterAuthorizedCallerAddedIterator, error) {

	logs, sub, err := _MultiAggregateRateLimiter.contract.FilterLogs(opts, "AuthorizedCallerAdded")
	if err != nil {
		return nil, err
	}
	return &MultiAggregateRateLimiterAuthorizedCallerAddedIterator{contract: _MultiAggregateRateLimiter.contract, event: "AuthorizedCallerAdded", logs: logs, sub: sub}, nil
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterFilterer) WatchAuthorizedCallerAdded(opts *bind.WatchOpts, sink chan<- *MultiAggregateRateLimiterAuthorizedCallerAdded) (event.Subscription, error) {

	logs, sub, err := _MultiAggregateRateLimiter.contract.WatchLogs(opts, "AuthorizedCallerAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MultiAggregateRateLimiterAuthorizedCallerAdded)
				if err := _MultiAggregateRateLimiter.contract.UnpackLog(event, "AuthorizedCallerAdded", log); err != nil {
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

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterFilterer) ParseAuthorizedCallerAdded(log types.Log) (*MultiAggregateRateLimiterAuthorizedCallerAdded, error) {
	event := new(MultiAggregateRateLimiterAuthorizedCallerAdded)
	if err := _MultiAggregateRateLimiter.contract.UnpackLog(event, "AuthorizedCallerAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MultiAggregateRateLimiterAuthorizedCallerRemovedIterator struct {
	Event *MultiAggregateRateLimiterAuthorizedCallerRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MultiAggregateRateLimiterAuthorizedCallerRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiAggregateRateLimiterAuthorizedCallerRemoved)
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
		it.Event = new(MultiAggregateRateLimiterAuthorizedCallerRemoved)
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

func (it *MultiAggregateRateLimiterAuthorizedCallerRemovedIterator) Error() error {
	return it.fail
}

func (it *MultiAggregateRateLimiterAuthorizedCallerRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MultiAggregateRateLimiterAuthorizedCallerRemoved struct {
	Caller common.Address
	Raw    types.Log
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterFilterer) FilterAuthorizedCallerRemoved(opts *bind.FilterOpts) (*MultiAggregateRateLimiterAuthorizedCallerRemovedIterator, error) {

	logs, sub, err := _MultiAggregateRateLimiter.contract.FilterLogs(opts, "AuthorizedCallerRemoved")
	if err != nil {
		return nil, err
	}
	return &MultiAggregateRateLimiterAuthorizedCallerRemovedIterator{contract: _MultiAggregateRateLimiter.contract, event: "AuthorizedCallerRemoved", logs: logs, sub: sub}, nil
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterFilterer) WatchAuthorizedCallerRemoved(opts *bind.WatchOpts, sink chan<- *MultiAggregateRateLimiterAuthorizedCallerRemoved) (event.Subscription, error) {

	logs, sub, err := _MultiAggregateRateLimiter.contract.WatchLogs(opts, "AuthorizedCallerRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MultiAggregateRateLimiterAuthorizedCallerRemoved)
				if err := _MultiAggregateRateLimiter.contract.UnpackLog(event, "AuthorizedCallerRemoved", log); err != nil {
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

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterFilterer) ParseAuthorizedCallerRemoved(log types.Log) (*MultiAggregateRateLimiterAuthorizedCallerRemoved, error) {
	event := new(MultiAggregateRateLimiterAuthorizedCallerRemoved)
	if err := _MultiAggregateRateLimiter.contract.UnpackLog(event, "AuthorizedCallerRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MultiAggregateRateLimiterConfigChangedIterator struct {
	Event *MultiAggregateRateLimiterConfigChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MultiAggregateRateLimiterConfigChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiAggregateRateLimiterConfigChanged)
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
		it.Event = new(MultiAggregateRateLimiterConfigChanged)
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

func (it *MultiAggregateRateLimiterConfigChangedIterator) Error() error {
	return it.fail
}

func (it *MultiAggregateRateLimiterConfigChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MultiAggregateRateLimiterConfigChanged struct {
	Config RateLimiterConfig
	Raw    types.Log
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterFilterer) FilterConfigChanged(opts *bind.FilterOpts) (*MultiAggregateRateLimiterConfigChangedIterator, error) {

	logs, sub, err := _MultiAggregateRateLimiter.contract.FilterLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return &MultiAggregateRateLimiterConfigChangedIterator{contract: _MultiAggregateRateLimiter.contract, event: "ConfigChanged", logs: logs, sub: sub}, nil
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterFilterer) WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *MultiAggregateRateLimiterConfigChanged) (event.Subscription, error) {

	logs, sub, err := _MultiAggregateRateLimiter.contract.WatchLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MultiAggregateRateLimiterConfigChanged)
				if err := _MultiAggregateRateLimiter.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
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

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterFilterer) ParseConfigChanged(log types.Log) (*MultiAggregateRateLimiterConfigChanged, error) {
	event := new(MultiAggregateRateLimiterConfigChanged)
	if err := _MultiAggregateRateLimiter.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MultiAggregateRateLimiterFeeQuoterSetIterator struct {
	Event *MultiAggregateRateLimiterFeeQuoterSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MultiAggregateRateLimiterFeeQuoterSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiAggregateRateLimiterFeeQuoterSet)
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
		it.Event = new(MultiAggregateRateLimiterFeeQuoterSet)
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

func (it *MultiAggregateRateLimiterFeeQuoterSetIterator) Error() error {
	return it.fail
}

func (it *MultiAggregateRateLimiterFeeQuoterSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MultiAggregateRateLimiterFeeQuoterSet struct {
	NewFeeQuoter common.Address
	Raw          types.Log
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterFilterer) FilterFeeQuoterSet(opts *bind.FilterOpts) (*MultiAggregateRateLimiterFeeQuoterSetIterator, error) {

	logs, sub, err := _MultiAggregateRateLimiter.contract.FilterLogs(opts, "FeeQuoterSet")
	if err != nil {
		return nil, err
	}
	return &MultiAggregateRateLimiterFeeQuoterSetIterator{contract: _MultiAggregateRateLimiter.contract, event: "FeeQuoterSet", logs: logs, sub: sub}, nil
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterFilterer) WatchFeeQuoterSet(opts *bind.WatchOpts, sink chan<- *MultiAggregateRateLimiterFeeQuoterSet) (event.Subscription, error) {

	logs, sub, err := _MultiAggregateRateLimiter.contract.WatchLogs(opts, "FeeQuoterSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MultiAggregateRateLimiterFeeQuoterSet)
				if err := _MultiAggregateRateLimiter.contract.UnpackLog(event, "FeeQuoterSet", log); err != nil {
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

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterFilterer) ParseFeeQuoterSet(log types.Log) (*MultiAggregateRateLimiterFeeQuoterSet, error) {
	event := new(MultiAggregateRateLimiterFeeQuoterSet)
	if err := _MultiAggregateRateLimiter.contract.UnpackLog(event, "FeeQuoterSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MultiAggregateRateLimiterOwnershipTransferRequestedIterator struct {
	Event *MultiAggregateRateLimiterOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MultiAggregateRateLimiterOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiAggregateRateLimiterOwnershipTransferRequested)
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
		it.Event = new(MultiAggregateRateLimiterOwnershipTransferRequested)
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

func (it *MultiAggregateRateLimiterOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *MultiAggregateRateLimiterOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MultiAggregateRateLimiterOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MultiAggregateRateLimiterOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MultiAggregateRateLimiter.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MultiAggregateRateLimiterOwnershipTransferRequestedIterator{contract: _MultiAggregateRateLimiter.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *MultiAggregateRateLimiterOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MultiAggregateRateLimiter.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MultiAggregateRateLimiterOwnershipTransferRequested)
				if err := _MultiAggregateRateLimiter.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterFilterer) ParseOwnershipTransferRequested(log types.Log) (*MultiAggregateRateLimiterOwnershipTransferRequested, error) {
	event := new(MultiAggregateRateLimiterOwnershipTransferRequested)
	if err := _MultiAggregateRateLimiter.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MultiAggregateRateLimiterOwnershipTransferredIterator struct {
	Event *MultiAggregateRateLimiterOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MultiAggregateRateLimiterOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiAggregateRateLimiterOwnershipTransferred)
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
		it.Event = new(MultiAggregateRateLimiterOwnershipTransferred)
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

func (it *MultiAggregateRateLimiterOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *MultiAggregateRateLimiterOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MultiAggregateRateLimiterOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MultiAggregateRateLimiterOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MultiAggregateRateLimiter.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MultiAggregateRateLimiterOwnershipTransferredIterator{contract: _MultiAggregateRateLimiter.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MultiAggregateRateLimiterOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MultiAggregateRateLimiter.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MultiAggregateRateLimiterOwnershipTransferred)
				if err := _MultiAggregateRateLimiter.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterFilterer) ParseOwnershipTransferred(log types.Log) (*MultiAggregateRateLimiterOwnershipTransferred, error) {
	event := new(MultiAggregateRateLimiterOwnershipTransferred)
	if err := _MultiAggregateRateLimiter.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MultiAggregateRateLimiterRateLimiterConfigUpdatedIterator struct {
	Event *MultiAggregateRateLimiterRateLimiterConfigUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MultiAggregateRateLimiterRateLimiterConfigUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiAggregateRateLimiterRateLimiterConfigUpdated)
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
		it.Event = new(MultiAggregateRateLimiterRateLimiterConfigUpdated)
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

func (it *MultiAggregateRateLimiterRateLimiterConfigUpdatedIterator) Error() error {
	return it.fail
}

func (it *MultiAggregateRateLimiterRateLimiterConfigUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MultiAggregateRateLimiterRateLimiterConfigUpdated struct {
	RemoteChainSelector uint64
	IsOutboundLane      bool
	Config              RateLimiterConfig
	Raw                 types.Log
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterFilterer) FilterRateLimiterConfigUpdated(opts *bind.FilterOpts, remoteChainSelector []uint64) (*MultiAggregateRateLimiterRateLimiterConfigUpdatedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _MultiAggregateRateLimiter.contract.FilterLogs(opts, "RateLimiterConfigUpdated", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &MultiAggregateRateLimiterRateLimiterConfigUpdatedIterator{contract: _MultiAggregateRateLimiter.contract, event: "RateLimiterConfigUpdated", logs: logs, sub: sub}, nil
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterFilterer) WatchRateLimiterConfigUpdated(opts *bind.WatchOpts, sink chan<- *MultiAggregateRateLimiterRateLimiterConfigUpdated, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _MultiAggregateRateLimiter.contract.WatchLogs(opts, "RateLimiterConfigUpdated", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MultiAggregateRateLimiterRateLimiterConfigUpdated)
				if err := _MultiAggregateRateLimiter.contract.UnpackLog(event, "RateLimiterConfigUpdated", log); err != nil {
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

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterFilterer) ParseRateLimiterConfigUpdated(log types.Log) (*MultiAggregateRateLimiterRateLimiterConfigUpdated, error) {
	event := new(MultiAggregateRateLimiterRateLimiterConfigUpdated)
	if err := _MultiAggregateRateLimiter.contract.UnpackLog(event, "RateLimiterConfigUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MultiAggregateRateLimiterTokenAggregateRateLimitAddedIterator struct {
	Event *MultiAggregateRateLimiterTokenAggregateRateLimitAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MultiAggregateRateLimiterTokenAggregateRateLimitAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiAggregateRateLimiterTokenAggregateRateLimitAdded)
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
		it.Event = new(MultiAggregateRateLimiterTokenAggregateRateLimitAdded)
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

func (it *MultiAggregateRateLimiterTokenAggregateRateLimitAddedIterator) Error() error {
	return it.fail
}

func (it *MultiAggregateRateLimiterTokenAggregateRateLimitAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MultiAggregateRateLimiterTokenAggregateRateLimitAdded struct {
	RemoteChainSelector uint64
	RemoteToken         []byte
	LocalToken          common.Address
	Raw                 types.Log
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterFilterer) FilterTokenAggregateRateLimitAdded(opts *bind.FilterOpts) (*MultiAggregateRateLimiterTokenAggregateRateLimitAddedIterator, error) {

	logs, sub, err := _MultiAggregateRateLimiter.contract.FilterLogs(opts, "TokenAggregateRateLimitAdded")
	if err != nil {
		return nil, err
	}
	return &MultiAggregateRateLimiterTokenAggregateRateLimitAddedIterator{contract: _MultiAggregateRateLimiter.contract, event: "TokenAggregateRateLimitAdded", logs: logs, sub: sub}, nil
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterFilterer) WatchTokenAggregateRateLimitAdded(opts *bind.WatchOpts, sink chan<- *MultiAggregateRateLimiterTokenAggregateRateLimitAdded) (event.Subscription, error) {

	logs, sub, err := _MultiAggregateRateLimiter.contract.WatchLogs(opts, "TokenAggregateRateLimitAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MultiAggregateRateLimiterTokenAggregateRateLimitAdded)
				if err := _MultiAggregateRateLimiter.contract.UnpackLog(event, "TokenAggregateRateLimitAdded", log); err != nil {
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

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterFilterer) ParseTokenAggregateRateLimitAdded(log types.Log) (*MultiAggregateRateLimiterTokenAggregateRateLimitAdded, error) {
	event := new(MultiAggregateRateLimiterTokenAggregateRateLimitAdded)
	if err := _MultiAggregateRateLimiter.contract.UnpackLog(event, "TokenAggregateRateLimitAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MultiAggregateRateLimiterTokenAggregateRateLimitRemovedIterator struct {
	Event *MultiAggregateRateLimiterTokenAggregateRateLimitRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MultiAggregateRateLimiterTokenAggregateRateLimitRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiAggregateRateLimiterTokenAggregateRateLimitRemoved)
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
		it.Event = new(MultiAggregateRateLimiterTokenAggregateRateLimitRemoved)
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

func (it *MultiAggregateRateLimiterTokenAggregateRateLimitRemovedIterator) Error() error {
	return it.fail
}

func (it *MultiAggregateRateLimiterTokenAggregateRateLimitRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MultiAggregateRateLimiterTokenAggregateRateLimitRemoved struct {
	RemoteChainSelector uint64
	LocalToken          common.Address
	Raw                 types.Log
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterFilterer) FilterTokenAggregateRateLimitRemoved(opts *bind.FilterOpts) (*MultiAggregateRateLimiterTokenAggregateRateLimitRemovedIterator, error) {

	logs, sub, err := _MultiAggregateRateLimiter.contract.FilterLogs(opts, "TokenAggregateRateLimitRemoved")
	if err != nil {
		return nil, err
	}
	return &MultiAggregateRateLimiterTokenAggregateRateLimitRemovedIterator{contract: _MultiAggregateRateLimiter.contract, event: "TokenAggregateRateLimitRemoved", logs: logs, sub: sub}, nil
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterFilterer) WatchTokenAggregateRateLimitRemoved(opts *bind.WatchOpts, sink chan<- *MultiAggregateRateLimiterTokenAggregateRateLimitRemoved) (event.Subscription, error) {

	logs, sub, err := _MultiAggregateRateLimiter.contract.WatchLogs(opts, "TokenAggregateRateLimitRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MultiAggregateRateLimiterTokenAggregateRateLimitRemoved)
				if err := _MultiAggregateRateLimiter.contract.UnpackLog(event, "TokenAggregateRateLimitRemoved", log); err != nil {
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

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterFilterer) ParseTokenAggregateRateLimitRemoved(log types.Log) (*MultiAggregateRateLimiterTokenAggregateRateLimitRemoved, error) {
	event := new(MultiAggregateRateLimiterTokenAggregateRateLimitRemoved)
	if err := _MultiAggregateRateLimiter.contract.UnpackLog(event, "TokenAggregateRateLimitRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MultiAggregateRateLimiterTokensConsumedIterator struct {
	Event *MultiAggregateRateLimiterTokensConsumed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MultiAggregateRateLimiterTokensConsumedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiAggregateRateLimiterTokensConsumed)
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
		it.Event = new(MultiAggregateRateLimiterTokensConsumed)
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

func (it *MultiAggregateRateLimiterTokensConsumedIterator) Error() error {
	return it.fail
}

func (it *MultiAggregateRateLimiterTokensConsumedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MultiAggregateRateLimiterTokensConsumed struct {
	Tokens *big.Int
	Raw    types.Log
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterFilterer) FilterTokensConsumed(opts *bind.FilterOpts) (*MultiAggregateRateLimiterTokensConsumedIterator, error) {

	logs, sub, err := _MultiAggregateRateLimiter.contract.FilterLogs(opts, "TokensConsumed")
	if err != nil {
		return nil, err
	}
	return &MultiAggregateRateLimiterTokensConsumedIterator{contract: _MultiAggregateRateLimiter.contract, event: "TokensConsumed", logs: logs, sub: sub}, nil
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterFilterer) WatchTokensConsumed(opts *bind.WatchOpts, sink chan<- *MultiAggregateRateLimiterTokensConsumed) (event.Subscription, error) {

	logs, sub, err := _MultiAggregateRateLimiter.contract.WatchLogs(opts, "TokensConsumed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MultiAggregateRateLimiterTokensConsumed)
				if err := _MultiAggregateRateLimiter.contract.UnpackLog(event, "TokensConsumed", log); err != nil {
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

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiterFilterer) ParseTokensConsumed(log types.Log) (*MultiAggregateRateLimiterTokensConsumed, error) {
	event := new(MultiAggregateRateLimiterTokensConsumed)
	if err := _MultiAggregateRateLimiter.contract.UnpackLog(event, "TokensConsumed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetAllRateLimitTokens struct {
	LocalTokens  []common.Address
	RemoteTokens [][]byte
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiter) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _MultiAggregateRateLimiter.abi.Events["AuthorizedCallerAdded"].ID:
		return _MultiAggregateRateLimiter.ParseAuthorizedCallerAdded(log)
	case _MultiAggregateRateLimiter.abi.Events["AuthorizedCallerRemoved"].ID:
		return _MultiAggregateRateLimiter.ParseAuthorizedCallerRemoved(log)
	case _MultiAggregateRateLimiter.abi.Events["ConfigChanged"].ID:
		return _MultiAggregateRateLimiter.ParseConfigChanged(log)
	case _MultiAggregateRateLimiter.abi.Events["FeeQuoterSet"].ID:
		return _MultiAggregateRateLimiter.ParseFeeQuoterSet(log)
	case _MultiAggregateRateLimiter.abi.Events["OwnershipTransferRequested"].ID:
		return _MultiAggregateRateLimiter.ParseOwnershipTransferRequested(log)
	case _MultiAggregateRateLimiter.abi.Events["OwnershipTransferred"].ID:
		return _MultiAggregateRateLimiter.ParseOwnershipTransferred(log)
	case _MultiAggregateRateLimiter.abi.Events["RateLimiterConfigUpdated"].ID:
		return _MultiAggregateRateLimiter.ParseRateLimiterConfigUpdated(log)
	case _MultiAggregateRateLimiter.abi.Events["TokenAggregateRateLimitAdded"].ID:
		return _MultiAggregateRateLimiter.ParseTokenAggregateRateLimitAdded(log)
	case _MultiAggregateRateLimiter.abi.Events["TokenAggregateRateLimitRemoved"].ID:
		return _MultiAggregateRateLimiter.ParseTokenAggregateRateLimitRemoved(log)
	case _MultiAggregateRateLimiter.abi.Events["TokensConsumed"].ID:
		return _MultiAggregateRateLimiter.ParseTokensConsumed(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (MultiAggregateRateLimiterAuthorizedCallerAdded) Topic() common.Hash {
	return common.HexToHash("0xeb1b9b92e50b7f88f9ff25d56765095ac6e91540eee214906f4036a908ffbdef")
}

func (MultiAggregateRateLimiterAuthorizedCallerRemoved) Topic() common.Hash {
	return common.HexToHash("0xc3803387881faad271c47728894e3e36fac830ffc8602ca6fc07733cbda77580")
}

func (MultiAggregateRateLimiterConfigChanged) Topic() common.Hash {
	return common.HexToHash("0x9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c19")
}

func (MultiAggregateRateLimiterFeeQuoterSet) Topic() common.Hash {
	return common.HexToHash("0x7c737a8eddf62436489aa3600ed26e75e0a58b0f8c0d266bbcee64358c39fdac")
}

func (MultiAggregateRateLimiterOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (MultiAggregateRateLimiterOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (MultiAggregateRateLimiterRateLimiterConfigUpdated) Topic() common.Hash {
	return common.HexToHash("0xf14a5415ce6988a9e870a85fff0b9d7b7dd79bbc228cb63cad610daf6f7b6b97")
}

func (MultiAggregateRateLimiterTokenAggregateRateLimitAdded) Topic() common.Hash {
	return common.HexToHash("0xad72a792d2a307f400c278be7deaeec6964276783304580cdc4e905436b8d5c5")
}

func (MultiAggregateRateLimiterTokenAggregateRateLimitRemoved) Topic() common.Hash {
	return common.HexToHash("0x530cabd30786b7235e124a6c0db77e0b685ef22813b1fe87554247f404eb8ed6")
}

func (MultiAggregateRateLimiterTokensConsumed) Topic() common.Hash {
	return common.HexToHash("0x1871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690a")
}

func (_MultiAggregateRateLimiter *MultiAggregateRateLimiter) Address() common.Address {
	return _MultiAggregateRateLimiter.address
}

type MultiAggregateRateLimiterInterface interface {
	CurrentRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64, isOutboundLane bool) (RateLimiterTokenBucket, error)

	GetAllAuthorizedCallers(opts *bind.CallOpts) ([]common.Address, error)

	GetAllRateLimitTokens(opts *bind.CallOpts, remoteChainSelector uint64) (GetAllRateLimitTokens,

		error)

	GetFeeQuoter(opts *bind.CallOpts) (common.Address, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	ApplyAuthorizedCallerUpdates(opts *bind.TransactOpts, authorizedCallerArgs AuthorizedCallersAuthorizedCallerArgs) (*types.Transaction, error)

	ApplyRateLimiterConfigUpdates(opts *bind.TransactOpts, rateLimiterUpdates []MultiAggregateRateLimiterRateLimiterConfigArgs) (*types.Transaction, error)

	OnInboundMessage(opts *bind.TransactOpts, message ClientAny2EVMMessage) (*types.Transaction, error)

	OnOutboundMessage(opts *bind.TransactOpts, destChainSelector uint64, message ClientEVM2AnyMessage) (*types.Transaction, error)

	SetFeeQuoter(opts *bind.TransactOpts, newFeeQuoter common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	UpdateRateLimitTokens(opts *bind.TransactOpts, removes []MultiAggregateRateLimiterLocalRateLimitToken, adds []MultiAggregateRateLimiterRateLimitTokenArgs) (*types.Transaction, error)

	FilterAuthorizedCallerAdded(opts *bind.FilterOpts) (*MultiAggregateRateLimiterAuthorizedCallerAddedIterator, error)

	WatchAuthorizedCallerAdded(opts *bind.WatchOpts, sink chan<- *MultiAggregateRateLimiterAuthorizedCallerAdded) (event.Subscription, error)

	ParseAuthorizedCallerAdded(log types.Log) (*MultiAggregateRateLimiterAuthorizedCallerAdded, error)

	FilterAuthorizedCallerRemoved(opts *bind.FilterOpts) (*MultiAggregateRateLimiterAuthorizedCallerRemovedIterator, error)

	WatchAuthorizedCallerRemoved(opts *bind.WatchOpts, sink chan<- *MultiAggregateRateLimiterAuthorizedCallerRemoved) (event.Subscription, error)

	ParseAuthorizedCallerRemoved(log types.Log) (*MultiAggregateRateLimiterAuthorizedCallerRemoved, error)

	FilterConfigChanged(opts *bind.FilterOpts) (*MultiAggregateRateLimiterConfigChangedIterator, error)

	WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *MultiAggregateRateLimiterConfigChanged) (event.Subscription, error)

	ParseConfigChanged(log types.Log) (*MultiAggregateRateLimiterConfigChanged, error)

	FilterFeeQuoterSet(opts *bind.FilterOpts) (*MultiAggregateRateLimiterFeeQuoterSetIterator, error)

	WatchFeeQuoterSet(opts *bind.WatchOpts, sink chan<- *MultiAggregateRateLimiterFeeQuoterSet) (event.Subscription, error)

	ParseFeeQuoterSet(log types.Log) (*MultiAggregateRateLimiterFeeQuoterSet, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MultiAggregateRateLimiterOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *MultiAggregateRateLimiterOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*MultiAggregateRateLimiterOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MultiAggregateRateLimiterOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MultiAggregateRateLimiterOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*MultiAggregateRateLimiterOwnershipTransferred, error)

	FilterRateLimiterConfigUpdated(opts *bind.FilterOpts, remoteChainSelector []uint64) (*MultiAggregateRateLimiterRateLimiterConfigUpdatedIterator, error)

	WatchRateLimiterConfigUpdated(opts *bind.WatchOpts, sink chan<- *MultiAggregateRateLimiterRateLimiterConfigUpdated, remoteChainSelector []uint64) (event.Subscription, error)

	ParseRateLimiterConfigUpdated(log types.Log) (*MultiAggregateRateLimiterRateLimiterConfigUpdated, error)

	FilterTokenAggregateRateLimitAdded(opts *bind.FilterOpts) (*MultiAggregateRateLimiterTokenAggregateRateLimitAddedIterator, error)

	WatchTokenAggregateRateLimitAdded(opts *bind.WatchOpts, sink chan<- *MultiAggregateRateLimiterTokenAggregateRateLimitAdded) (event.Subscription, error)

	ParseTokenAggregateRateLimitAdded(log types.Log) (*MultiAggregateRateLimiterTokenAggregateRateLimitAdded, error)

	FilterTokenAggregateRateLimitRemoved(opts *bind.FilterOpts) (*MultiAggregateRateLimiterTokenAggregateRateLimitRemovedIterator, error)

	WatchTokenAggregateRateLimitRemoved(opts *bind.WatchOpts, sink chan<- *MultiAggregateRateLimiterTokenAggregateRateLimitRemoved) (event.Subscription, error)

	ParseTokenAggregateRateLimitRemoved(log types.Log) (*MultiAggregateRateLimiterTokenAggregateRateLimitRemoved, error)

	FilterTokensConsumed(opts *bind.FilterOpts) (*MultiAggregateRateLimiterTokensConsumedIterator, error)

	WatchTokensConsumed(opts *bind.WatchOpts, sink chan<- *MultiAggregateRateLimiterTokensConsumed) (event.Subscription, error)

	ParseTokensConsumed(log types.Log) (*MultiAggregateRateLimiterTokensConsumed, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
