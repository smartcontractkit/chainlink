// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package onramp

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

type InternalEVM2AnyRampMessage struct {
	Header         InternalRampMessageHeader
	Sender         common.Address
	Data           []byte
	Receiver       []byte
	ExtraArgs      []byte
	FeeToken       common.Address
	FeeTokenAmount *big.Int
	FeeValueJuels  *big.Int
	TokenAmounts   []InternalEVM2AnyTokenTransfer
}

type InternalEVM2AnyTokenTransfer struct {
	SourcePoolAddress common.Address
	DestTokenAddress  []byte
	ExtraData         []byte
	Amount            *big.Int
	DestExecData      []byte
}

type InternalRampMessageHeader struct {
	MessageId           [32]byte
	SourceChainSelector uint64
	DestChainSelector   uint64
	SequenceNumber      uint64
	Nonce               uint64
}

type OnRampAllowlistConfigArgs struct {
	DestChainSelector         uint64
	AllowlistEnabled          bool
	AddedAllowlistedSenders   []common.Address
	RemovedAllowlistedSenders []common.Address
}

type OnRampDestChainConfigArgs struct {
	DestChainSelector uint64
	Router            common.Address
	AllowlistEnabled  bool
}

type OnRampDynamicConfig struct {
	FeeQuoter              common.Address
	ReentrancyGuardEntered bool
	MessageInterceptor     common.Address
	FeeAggregator          common.Address
	AllowlistAdmin         common.Address
}

type OnRampStaticConfig struct {
	ChainSelector      uint64
	RmnRemote          common.Address
	NonceManager       common.Address
	TokenAdminRegistry common.Address
}

var OnRampMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"staticConfig\",\"type\":\"tuple\",\"internalType\":\"structOnRamp.StaticConfig\",\"components\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"rmnRemote\",\"type\":\"address\",\"internalType\":\"contractIRMNRemote\"},{\"name\":\"nonceManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenAdminRegistry\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"name\":\"dynamicConfig\",\"type\":\"tuple\",\"internalType\":\"structOnRamp.DynamicConfig\",\"components\":[{\"name\":\"feeQuoter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"reentrancyGuardEntered\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"messageInterceptor\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"feeAggregator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowlistAdmin\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"name\":\"destChainConfigArgs\",\"type\":\"tuple[]\",\"internalType\":\"structOnRamp.DestChainConfigArgs[]\",\"components\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"contractIRouter\"},{\"name\":\"allowlistEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"applyAllowlistUpdates\",\"inputs\":[{\"name\":\"allowlistConfigArgsItems\",\"type\":\"tuple[]\",\"internalType\":\"structOnRamp.AllowlistConfigArgs[]\",\"components\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"allowlistEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"addedAllowlistedSenders\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"removedAllowlistedSenders\",\"type\":\"address[]\",\"internalType\":\"address[]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"applyDestChainConfigUpdates\",\"inputs\":[{\"name\":\"destChainConfigArgs\",\"type\":\"tuple[]\",\"internalType\":\"structOnRamp.DestChainConfigArgs[]\",\"components\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"contractIRouter\"},{\"name\":\"allowlistEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"forwardFromRouter\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"message\",\"type\":\"tuple\",\"internalType\":\"structClient.EVM2AnyMessage\",\"components\":[{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"tokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structClient.EVMTokenAmount[]\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"feeToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"extraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"feeTokenAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"originalSender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAllowedSendersList\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"configuredAddresses\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDestChainConfig\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"allowlistEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDynamicConfig\",\"inputs\":[],\"outputs\":[{\"name\":\"dynamicConfig\",\"type\":\"tuple\",\"internalType\":\"structOnRamp.DynamicConfig\",\"components\":[{\"name\":\"feeQuoter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"reentrancyGuardEntered\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"messageInterceptor\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"feeAggregator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowlistAdmin\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getExpectedNextSequenceNumber\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFee\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"message\",\"type\":\"tuple\",\"internalType\":\"structClient.EVM2AnyMessage\",\"components\":[{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"tokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structClient.EVMTokenAmount[]\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"feeToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"extraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"feeTokenAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPoolBySourceToken\",\"inputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sourceToken\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIPoolV1\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStaticConfig\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structOnRamp.StaticConfig\",\"components\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"rmnRemote\",\"type\":\"address\",\"internalType\":\"contractIRMNRemote\"},{\"name\":\"nonceManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenAdminRegistry\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSupportedTokens\",\"inputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setDynamicConfig\",\"inputs\":[{\"name\":\"dynamicConfig\",\"type\":\"tuple\",\"internalType\":\"structOnRamp.DynamicConfig\",\"components\":[{\"name\":\"feeQuoter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"reentrancyGuardEntered\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"messageInterceptor\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"feeAggregator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowlistAdmin\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"typeAndVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdrawFeeTokens\",\"inputs\":[{\"name\":\"feeTokens\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AllowListAdminSet\",\"inputs\":[{\"name\":\"allowlistAdmin\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AllowListSendersAdded\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"senders\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AllowListSendersRemoved\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"senders\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CCIPMessageSent\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"message\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structInternal.EVM2AnyRampMessage\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structInternal.RampMessageHeader\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"nonce\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"extraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"feeToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"feeTokenAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"feeValueJuels\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.EVM2AnyTokenTransfer[]\",\"components\":[{\"name\":\"sourcePoolAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"destTokenAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"extraData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"destExecData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ConfigSet\",\"inputs\":[{\"name\":\"staticConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structOnRamp.StaticConfig\",\"components\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"rmnRemote\",\"type\":\"address\",\"internalType\":\"contractIRMNRemote\"},{\"name\":\"nonceManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenAdminRegistry\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"name\":\"dynamicConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structOnRamp.DynamicConfig\",\"components\":[{\"name\":\"feeQuoter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"reentrancyGuardEntered\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"messageInterceptor\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"feeAggregator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowlistAdmin\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DestChainConfigSet\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"router\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"contractIRouter\"},{\"name\":\"allowlistEnabled\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FeeTokenWithdrawn\",\"inputs\":[{\"name\":\"feeAggregator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"feeToken\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"CannotSendZeroTokens\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CannotTransferToSelf\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CursedByRMN\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"GetSupportedTokensFunctionalityRemovedCheckAdminRegistry\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidAllowListRequest\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidConfig\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidDestChainConfig\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MustBeCalledByRouter\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MustBeProposedOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwnerOrAllowlistAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnerCannotBeZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RouterMustSetOriginalSender\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SenderNotAllowed\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UnsupportedToken\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	Bin: "0x6101006040523461055f57613c018038038061001a81610599565b92833981019080820390610140821261055f576080821261055f5761003d61057a565b90610047816105be565b82526020810151906001600160a01b038216820361055f5760208301918252610072604082016105d2565b926040810193845260a0610088606084016105d2565b6060830190815295607f19011261055f5760405160a081016001600160401b03811182821017610564576040526100c1608084016105d2565b81526100cf60a084016105e6565b602082019081526100e260c085016105d2565b91604081019283526100f660e086016105d2565b936060820194855261010b61010087016105d2565b6080830190815261012087015190966001600160401b03821161055f57018a601f8201121561055f578051906001600160401b0382116105645760209b8c6060610159828660051b01610599565b9e8f8681520194028301019181831161055f57602001925b8284106104f1575050505033156104e057600180546001600160a01b0319163317905580516001600160401b03161580156104ce575b80156104bc575b80156104aa575b61047d57516001600160401b0316608081905295516001600160a01b0390811660a08190529751811660c08190529851811660e08190528251909116158015610498575b801561048e575b61047d57815160028054855160ff60a01b90151560a01b166001600160a01b039384166001600160a81b0319909216919091171790558451600380549183166001600160a01b03199283161790558651600480549184169183169190911790558751600580549190931691161790557fc7372d2d886367d7bb1b0e0708a5436f2c91d6963de210eb2dc1ec2ecd6d21f19861012098606061029f61057a565b8a8152602080820193845260408083019586529290910194855281519a8b5291516001600160a01b03908116928b019290925291518116918901919091529051811660608801529051811660808701529051151560a08601529051811660c08501529051811660e0840152905116610100820152a16000905b805182101561040b5761032b82826105f3565b51916001600160401b0361033f82846105f3565b5151169283156103f65760008481526006602090815260409182902081840151815494840151600160401b600160e81b03198616604883901b600160481b600160e81b031617901515851b68ff000000000000000016179182905583516001600160401b0390951685526001600160a01b031691840191909152811c60ff1615159082015291926001927fd5ad72bc37dc7a80a8b9b9df20500046fd7341adb1be2258a540466fdd7dcef590606090a20190610318565b8363c35aa79d60e01b60005260045260246000fd5b6040516135e3908161061e82396080518181816103e301528181610b39015281816120c8015261286b015260a0518181816121010152818161261e01526128a4015260c051818181610efc0152818161213d01526128e0015260e0518181816121790152818161291c0152612edd0152f35b6306b7c75960e31b60005260046000fd5b5082511515610200565b5084516001600160a01b0316156101f9565b5088516001600160a01b0316156101b5565b5087516001600160a01b0316156101ae565b5086516001600160a01b0316156101a7565b639b15e16f60e01b60005260046000fd5b60608483031261055f5760405190606082016001600160401b0381118382101761056457604052610521856105be565b82526020850151906001600160a01b038216820361055f578260209283606095015261054f604088016105e6565b6040820152815201930192610171565b600080fd5b634e487b7160e01b600052604160045260246000fd5b60405190608082016001600160401b0381118382101761056457604052565b6040519190601f01601f191682016001600160401b0381118382101761056457604052565b51906001600160401b038216820361055f57565b51906001600160a01b038216820361055f57565b5190811515820361055f57565b80518210156106075760209160051b010190565b634e487b7160e01b600052603260045260246000fdfe608080604052600436101561001357600080fd5b600090813560e01c90816306285c691461280557508063181f5a771461278657806320487ded146125425780632716072b1461229257806327e936f114611e8c57806348a98aa414611e095780635cb80c5d14611b215780636def4ce714611a925780637437ff9f1461195957806379ba5097146118745780638da5cb5b146118225780639041be3d14611775578063972b4612146116a7578063c9b146b3146112da578063df0aa9e91461022c578063f2fde38b1461013f5763fbca3b74146100dc57600080fd5b3461013c5760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261013c57600490610116612af4565b507f9e7177c8000000000000000000000000000000000000000000000000000000008152fd5b80fd5b503461013c5760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261013c5773ffffffffffffffffffffffffffffffffffffffff61018c612b4a565b610194613280565b1633811461020457807fffffffffffffffffffffffff000000000000000000000000000000000000000083541617825573ffffffffffffffffffffffffffffffffffffffff600154167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12788380a380f35b6004827fdad89dca000000000000000000000000000000000000000000000000000000008152fd5b503461013c5760807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261013c57610264612af4565b67ffffffffffffffff602435116112d65760a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc602435360301126112d6576102ab612b6d565b60025460ff8160a01c166112ae577fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff16740100000000000000000000000000000000000000001760025567ffffffffffffffff8216835260066020526040832073ffffffffffffffffffffffffffffffffffffffff82161561128657805460ff8160401c16611218575b60481c73ffffffffffffffffffffffffffffffffffffffff1633036111f05773ffffffffffffffffffffffffffffffffffffffff6003541680611183575b50805467ffffffffffffffff811667ffffffffffffffff8114611156579067ffffffffffffffff60017fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000009493011692839116179055604051906103d5826129db565b84825267ffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016602083015267ffffffffffffffff8416604083015260608201528360808201526104366024803501602435600401613060565b61044560046024350180613060565b610453606460243501612f97565b936104686044602435016024356004016130b1565b9490507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe06104ae61049887612b25565b966104a66040519889612a13565b808852612b25565b018a5b81811061113f575050604051968761012081011067ffffffffffffffff6101208a011117611112576101208801604052875273ffffffffffffffffffffffffffffffffffffffff8816602088015261051c939291610510913691613131565b60408701523691613131565b606084015273ffffffffffffffffffffffffffffffffffffffff6020926040516105468582612a13565b88815260808601521660a084015260443560c08401528560e084015261010083015261057c6044602435016024356004016130b1565b61058881969296612b25565b906105966040519283612a13565b8082528382018097368360061b8201116110895780915b8360061b820183106110d75750505050865b6105d36044602435016024356004016130b1565b9050811015610959576105e6818361301d565b516106006105f960046024350180613060565b3691613131565b90610609613105565b5085810151156109315773ffffffffffffffffffffffffffffffffffffffff61063481835116612e7e565b169182158015610887575b61084457908a6106ee949392888873ffffffffffffffffffffffffffffffffffffffff8d838701518280895116926040519761067a896129db565b885267ffffffffffffffff87890196168652816040890191168152606088019283526080880193845267ffffffffffffffff6040519d8e998a997f9a4575b9000000000000000000000000000000000000000000000000000000008b5260048b01525160a060248b015260c48a0190612ab1565b965116604488015251166064860152516084850152511660a4830152038183865af1938415610839578b94610783575b5061077c918484928980600198519301519101519160405193610740856129db565b84528a8401526040830152606082015260405161075d8982612a13565b8c8152608082015261010089015190610776838361301d565b5261301d565b50016105bf565b93503d91828c863e6107958386612a13565b878584810103126108355784519167ffffffffffffffff83116108315760408387018588010312610831576040516107cc816129f7565b8387015167ffffffffffffffff811161082d576107f09086890190868a0101613168565b81528984880101519467ffffffffffffffff861161082d578761077c96889661081f9360019b01920101613168565b8a820152955091509161071e565b8e80fd5b8c80fd5b8b80fd5b6040513d8d823e3d90fd5b60248b73ffffffffffffffffffffffffffffffffffffffff8451167fbf16aab6000000000000000000000000000000000000000000000000000000008252600452fd5b506040517f01ffc9a70000000000000000000000000000000000000000000000000000000081527faff2afbf0000000000000000000000000000000000000000000000000000000060048201528781602481875afa908115610926578c916108f1575b501561063f565b90508781813d831161091f575b6109088183612a13565b810103126108355761091990612c2c565b386108ea565b503d6108fe565b6040513d8e823e3d90fd5b60048a7f5cf04449000000000000000000000000000000000000000000000000000000008152fd5b50867f430d138c00000000000000000000000000000000000000000000000000000000839282979873ffffffffffffffffffffffffffffffffffffffff600254169187610a496109ad606460243501612f97565b610a198c61010073ffffffffffffffffffffffffffffffffffffffff6109dd608460243501602435600401613060565b92909301519467ffffffffffffffff6040519e8f9d8e521660048d01521660248b015260443560448b015260c060648b015260c48a0191612c89565b907ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc8883030160848901526131aa565b917ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc8684030160a48701525191828152019190855b8982821061109c575050505082809103915afa9586156110915785869087938899610f7f575b5060e087015215610e8d5750845b67ffffffffffffffff6080865101911690526080840152835b61010084015151811015610b015780610ae66001928861301d565b516080610af88361010089015161301d565b51015201610acb565b509092604051848101907f130ac867e79e2789f923760a88743d292acdf7002139a588206e2260f73f7321825267ffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016604082015267ffffffffffffffff8416606082015230608082015260808152610b8360a082612a13565b51902073ffffffffffffffffffffffffffffffffffffffff602085015116845167ffffffffffffffff6080816060840151169201511673ffffffffffffffffffffffffffffffffffffffff60a08801511660c088015191604051938a850195865260408501526060840152608083015260a082015260a08152610c0760c082612a13565b51902060608501518681519101206040860151878151910120610100870151604051610c6d81610c418c8201948d865260408301906131aa565b037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08101835282612a13565b51902091608088015189815191012093604051958a870197885260408701526060860152608085015260a084015260c083015260e082015260e08152610cb561010082612a13565b51902082515267ffffffffffffffff60608351015116907f192442a2b2adb6a7948f097023cb6b57d29d3a7a5dd33e6666d33c39cc456f3260405185815267ffffffffffffffff608086518051898501528289820151166040850152826040820151166060850152826060820151168285015201511660a082015273ffffffffffffffffffffffffffffffffffffffff60208601511660c082015280610e58610ddf610daa610d7560408a01516101a060e08701526101c0860190612ab1565b60608a01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe086830301610100870152612ab1565b60808901517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe085830301610120860152612ab1565b73ffffffffffffffffffffffffffffffffffffffff60a08901511661014084015260c088015161016084015260e088015161018084015267ffffffffffffffff610100890151967fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0858403016101a086015216956131aa565b0390a37fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff600254166002555151604051908152f35b73ffffffffffffffffffffffffffffffffffffffff604051917fea458c0c00000000000000000000000000000000000000000000000000000000835267ffffffffffffffff8516600484015216602482015283816044818973ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165af1908115610f74578691610f32575b50610ab2565b90508381813d8311610f6d575b610f498183612a13565b81010312610f69575167ffffffffffffffff81168103610f695787610f2c565b8580fd5b503d610f3f565b6040513d88823e3d90fd5b9850505090503d908186883e610f958288612a13565b6080878381010312610f6957865196610faf858201612c2c565b97604082015167ffffffffffffffff811161108d57610fd390858401908401613168565b9160608101519467ffffffffffffffff861161108957808201601f87840101121561108957858201519161100683612b25565b966110146040519889612a13565b838852898801928083018b8660051b848601010111610831578a82840101935b8b8660051b84860101018510611054575050505050509790929789610aa4565b845167ffffffffffffffff811161082d578c809261107c829383878a0191898b010101613168565b8152019501949050611034565b8980fd5b8880fd5b6040513d87823e3d90fd5b8351805173ffffffffffffffffffffffffffffffffffffffff168652810151818601528c975088965060409094019390920191600101610a7e565b60408336031261110e578660409182516110f0816129f7565b6110f986612b90565b815282860135838201528152019201916105ad565b8a80fd5b60248b7f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b60209061114a613105565b82828a010152016104b1565b6024867f4e487b710000000000000000000000000000000000000000000000000000000081526011600452fd5b803b156111ec578460405180927fe0a0e5060000000000000000000000000000000000000000000000000000000082528183816111c96024356004018b60048401612cc8565b03925af18015611091571561037357846111e591959295612a13565b9238610373565b8480fd5b6004847f1c0a3529000000000000000000000000000000000000000000000000000000008152fd5b73ffffffffffffffffffffffffffffffffffffffff831660009081526002830160205260409020546103355760248573ffffffffffffffffffffffffffffffffffffffff857fd0d2597600000000000000000000000000000000000000000000000000000000835216600452fd5b6004847fa4ec7479000000000000000000000000000000000000000000000000000000008152fd5b6004847f3ee5aeb5000000000000000000000000000000000000000000000000000000008152fd5b5080fd5b503461013c5760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261013c5760043567ffffffffffffffff81116112d65761132a903690600401612bb1565b73ffffffffffffffffffffffffffffffffffffffff60015416330361165f575b919081907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8181360301915b8481101561165b578060051b820135838112156111ec578201916080833603126111ec57604051946113a686612990565b6113af84612b10565b86526113bd60208501612b3d565b9660208701978852604085013567ffffffffffffffff8111611657576113e69036908701612fb8565b9460408801958652606081013567ffffffffffffffff81116116535761140e91369101612fb8565b60608801908152875167ffffffffffffffff1683526006602052604080842099518a547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff169015159182901b68ff000000000000000016178a5590959081515161152b575b5095976001019550815b855180518210156114bc57906114b573ffffffffffffffffffffffffffffffffffffffff6114ad8360019561301d565b511689613374565b500161147d565b505095909694506001929193519081516114dc575b505001939293611375565b61152167ffffffffffffffff7fc237ec1921f855ccd5e9a5af9733f2d58943a5a8501ec5988e305d7a4d42158692511692604051918291602083526020830190612be2565b0390a238806114d1565b9893959296919094979860001461161c57600184019591875b865180518210156115c15761156e8273ffffffffffffffffffffffffffffffffffffffff9261301d565b5116801561158a57906115836001928a6132e3565b5001611544565b60248a67ffffffffffffffff8e51167f463258ff000000000000000000000000000000000000000000000000000000008252600452fd5b50509692955090929796937f330939f6eafe8bb516716892fe962ff19770570838686e6579dbc1cc51fc328161161267ffffffffffffffff8a51169251604051918291602083526020830190612be2565b0390a23880611473565b60248767ffffffffffffffff8b51167f463258ff000000000000000000000000000000000000000000000000000000008252600452fd5b8380fd5b8280fd5b8380f35b73ffffffffffffffffffffffffffffffffffffffff6005541633031561134a576004837f905d7d9b000000000000000000000000000000000000000000000000000000008152fd5b503461013c5760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261013c5767ffffffffffffffff6116e8612af4565b16808252600660205260ff604083205460401c16908252600660205260016040832001916040518093849160208254918281520191845260208420935b81811061175c57505061173a92500383612a13565b61175860405192839215158352604060208401526040830190612be2565b0390f35b8454835260019485019487945060209093019201611725565b503461013c5760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261013c5767ffffffffffffffff6117b6612af4565b1681526006602052600167ffffffffffffffff604083205416019067ffffffffffffffff82116117f55760208267ffffffffffffffff60405191168152f35b807f4e487b7100000000000000000000000000000000000000000000000000000000602492526011600452fd5b503461013c57807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261013c57602073ffffffffffffffffffffffffffffffffffffffff60015416604051908152f35b503461013c57807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261013c57805473ffffffffffffffffffffffffffffffffffffffff81163303611931577fffffffffffffffffffffffff000000000000000000000000000000000000000060015491338284161760015516825573ffffffffffffffffffffffffffffffffffffffff3391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08380a380f35b6004827f02b543c6000000000000000000000000000000000000000000000000000000008152fd5b503461013c57807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261013c576080604051611996816129db565b828152826020820152826040820152826060820152015260a06040516119bb816129db565b60ff60025473ffffffffffffffffffffffffffffffffffffffff81168352831c161515602082015273ffffffffffffffffffffffffffffffffffffffff60035416604082015273ffffffffffffffffffffffffffffffffffffffff60045416606082015273ffffffffffffffffffffffffffffffffffffffff600554166080820152611a90604051809273ffffffffffffffffffffffffffffffffffffffff60808092828151168552602081015115156020860152826040820151166040860152826060820151166060860152015116910152565bf35b503461013c5760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261013c57604060609167ffffffffffffffff611ad8612af4565b1681526006602052205473ffffffffffffffffffffffffffffffffffffffff6040519167ffffffffffffffff8116835260ff8160401c161515602084015260481c166040820152f35b503461013c5760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261013c5760043567ffffffffffffffff81116112d657611b71903690600401612bb1565b9073ffffffffffffffffffffffffffffffffffffffff6004541690835b83811015611e055773ffffffffffffffffffffffffffffffffffffffff611bb98260051b8401612f97565b1690604051917f70a08231000000000000000000000000000000000000000000000000000000008352306004840152602083602481845afa928315611dfa578793611dc3575b5082611c11575b506001915001611b8e565b8460405193611ccc60208601957fa9059cbb00000000000000000000000000000000000000000000000000000000875283602482015282604482015260448152611c5c606482612a13565b8a80604098895193611c6e8b86612a13565b602085527f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65646020860152519082895af13d15611dbb573d90611caf82612a54565b91611cbc8a519384612a13565b82523d8d602084013e5b86613506565b805180611d08575b505060207f508d7d183612c18fc339b42618912b9fa3239f631dd7ec0671f950200a0fa66e9160019651908152a338611c06565b81929495969350906020918101031261108d576020611d279101612c2c565b15611d385792919085903880611cd4565b608490517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f742073756363656564000000000000000000000000000000000000000000006064820152fd5b606090611cc6565b9092506020813d8211611df2575b81611dde60209383612a13565b81010312611dee57519138611bff565b8680fd5b3d9150611dd1565b6040513d89823e3d90fd5b8480f35b503461013c5760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261013c57611e41612af4565b506024359073ffffffffffffffffffffffffffffffffffffffff8216820361013c576020611e6e83612e7e565b73ffffffffffffffffffffffffffffffffffffffff60405191168152f35b503461013c5760a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261013c57604051611ec8816129db565b611ed0612b4a565b81526024358015158103611657576020820190815260443573ffffffffffffffffffffffffffffffffffffffff811681036116535760408301908152611f14612b6d565b90606084019182526084359273ffffffffffffffffffffffffffffffffffffffff84168403610f695760808501938452611f4c613280565b73ffffffffffffffffffffffffffffffffffffffff855116158015612273575b8015612269575b612241579273ffffffffffffffffffffffffffffffffffffffff859381809461012097827fc7372d2d886367d7bb1b0e0708a5436f2c91d6963de210eb2dc1ec2ecd6d21f19a51167fffffffffffffffffffffffff000000000000000000000000000000000000000060025416176002555115157fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff74ff00000000000000000000000000000000000000006002549260a01b1691161760025551167fffffffffffffffffffffffff0000000000000000000000000000000000000000600354161760035551167fffffffffffffffffffffffff0000000000000000000000000000000000000000600454161760045551167fffffffffffffffffffffffff0000000000000000000000000000000000000000600554161760055561223d604051916120bd83612990565b67ffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016835273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016602084015273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016604084015273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001660608401526121ed604051809473ffffffffffffffffffffffffffffffffffffffff6060809267ffffffffffffffff8151168552826020820151166020860152826040820151166040860152015116910152565b608083019073ffffffffffffffffffffffffffffffffffffffff60808092828151168552602081015115156020860152826040820151166040860152826060820151166060860152015116910152565ba180f35b6004867f35be3ac8000000000000000000000000000000000000000000000000000000008152fd5b5080511515611f73565b5073ffffffffffffffffffffffffffffffffffffffff83511615611f6c565b503461013c5760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261013c576004359067ffffffffffffffff821161013c573660238301121561013c5781600401356122ee81612b25565b926122fc6040519485612a13565b818452602460606020860193028201019036821161165357602401915b81831061249a5750505061232b613280565b805b825181101561249657612340818461301d565b5167ffffffffffffffff612354838661301d565b51511690811561246a57907fd5ad72bc37dc7a80a8b9b9df20500046fd7341adb1be2258a540466fdd7dcef5606060019493838752600660205260ff6040882061242b604060208501519483547fffffff0000000000000000000000000000000000000000ffffffffffffffffff7cffffffffffffffffffffffffffffffffffffffff0000000000000000008860481b1691161784550151151582907fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff68ff0000000000000000835492151560401b169116179055565b5473ffffffffffffffffffffffffffffffffffffffff6040519367ffffffffffffffff8316855216602084015260401c1615156040820152a20161232d565b602484837fc35aa79d000000000000000000000000000000000000000000000000000000008252600452fd5b5080f35b606083360312611653576040516060810181811067ffffffffffffffff821117612515576040526124ca84612b10565b8152602084013573ffffffffffffffffffffffffffffffffffffffff81168103610f6957918160609360208094015261250560408701612b3d565b6040820152815201920191612319565b6024867f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b503461013c5760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261013c5761257a612af4565b60243567ffffffffffffffff81116116575760a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc8236030112611657576040517f2cbc26bb00000000000000000000000000000000000000000000000000000000815277ffffffffffffffff000000000000000000000000000000008360801b16600482015260208160248173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165afa90811561277b578491612741575b5061270b576126ad9160209173ffffffffffffffffffffffffffffffffffffffff60025416906040518095819482937fd8694ccd0000000000000000000000000000000000000000000000000000000084526004019060048401612cc8565b03915afa9081156127005782916126ca575b602082604051908152f35b90506020813d6020116126f8575b816126e560209383612a13565b810103126112d6576020915051386126bf565b3d91506126d8565b6040513d84823e3d90fd5b60248367ffffffffffffffff847ffdbd6a7200000000000000000000000000000000000000000000000000000000835216600452fd5b90506020813d602011612773575b8161275c60209383612a13565b810103126116535761276d90612c2c565b3861264e565b3d915061274f565b6040513d86823e3d90fd5b503461013c57807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261013c57506117586040516127c7604082612a13565b601081527f4f6e52616d7020312e362e302d646576000000000000000000000000000000006020820152604051918291602083526020830190612ab1565b9050346112d657817ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126112d65780612841606092612990565b8281528260208201528260408201520152608060405161286081612990565b67ffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016815273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016602082015273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016604082015273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000166060820152611a90604051809273ffffffffffffffffffffffffffffffffffffffff6060809267ffffffffffffffff8151168552826020820151166020860152826040820151166040860152015116910152565b6080810190811067ffffffffffffffff8211176129ac57604052565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60a0810190811067ffffffffffffffff8211176129ac57604052565b6040810190811067ffffffffffffffff8211176129ac57604052565b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff8211176129ac57604052565b67ffffffffffffffff81116129ac57601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b60005b838110612aa15750506000910152565b8181015183820152602001612a91565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f602093612aed81518092818752878088019101612a8e565b0116010190565b6004359067ffffffffffffffff82168203612b0b57565b600080fd5b359067ffffffffffffffff82168203612b0b57565b67ffffffffffffffff81116129ac5760051b60200190565b35908115158203612b0b57565b6004359073ffffffffffffffffffffffffffffffffffffffff82168203612b0b57565b6064359073ffffffffffffffffffffffffffffffffffffffff82168203612b0b57565b359073ffffffffffffffffffffffffffffffffffffffff82168203612b0b57565b9181601f84011215612b0b5782359167ffffffffffffffff8311612b0b576020808501948460051b010111612b0b57565b906020808351928381520192019060005b818110612c005750505090565b825173ffffffffffffffffffffffffffffffffffffffff16845260209384019390920191600101612bf3565b51908115158203612b0b57565b90357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe182360301811215612b0b57016020813591019167ffffffffffffffff8211612b0b578136038313612b0b57565b601f82602094937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0938186528686013760008582860101520116010190565b9067ffffffffffffffff9093929316815260406020820152612d3e612d01612cf08580612c39565b60a0604086015260e0850191612c89565b612d0e6020860186612c39565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0858403016060860152612c89565b9060408401357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe185360301811215612b0b5784016020813591019267ffffffffffffffff8211612b0b578160061b36038413612b0b578281037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0016080840152818152602001929060005b818110612e3f57505050612e0c8473ffffffffffffffffffffffffffffffffffffffff612dfc6060612e3c979801612b90565b1660a08401526080810190612c39565b9160c07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc082860301910152612c89565b90565b90919360408060019273ffffffffffffffffffffffffffffffffffffffff612e6689612b90565b16815260208881013590820152019501929101612dc9565b73ffffffffffffffffffffffffffffffffffffffff604051917fbbe4f6db00000000000000000000000000000000000000000000000000000000835216600482015260208160248173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165afa908115612f8b57600091612f28575b5073ffffffffffffffffffffffffffffffffffffffff1690565b6020813d602011612f83575b81612f4160209383612a13565b810103126112d657519073ffffffffffffffffffffffffffffffffffffffff8216820361013c575073ffffffffffffffffffffffffffffffffffffffff612f0e565b3d9150612f34565b6040513d6000823e3d90fd5b3573ffffffffffffffffffffffffffffffffffffffff81168103612b0b5790565b9080601f83011215612b0b578135612fcf81612b25565b92612fdd6040519485612a13565b81845260208085019260051b820101928311612b0b57602001905b8282106130055750505090565b6020809161301284612b90565b815201910190612ff8565b80518210156130315760209160051b010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe181360301821215612b0b570180359067ffffffffffffffff8211612b0b57602001918136038313612b0b57565b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe181360301821215612b0b570180359067ffffffffffffffff8211612b0b57602001918160061b36038313612b0b57565b60405190613112826129db565b6060608083600081528260208201528260408201526000838201520152565b92919261313d82612a54565b9161314b6040519384612a13565b829481845281830111612b0b578281602093846000960137010152565b81601f82011215612b0b57805161317e81612a54565b9261318c6040519485612a13565b81845260208284010111612b0b57612e3c9160208085019101612a8e565b9080602083519182815201916020808360051b8301019401926000915b8383106131d657505050505090565b9091929394602080613271837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0866001960301875289519073ffffffffffffffffffffffffffffffffffffffff825116815260806132566132448685015160a08886015260a0850190612ab1565b60408501518482036040860152612ab1565b92606081015160608401520151906080818403910152612ab1565b970193019301919392906131c7565b73ffffffffffffffffffffffffffffffffffffffff6001541633036132a157565b7f2b5c74de0000000000000000000000000000000000000000000000000000000060005260046000fd5b80548210156130315760005260206000200190600090565b600082815260018201602052604090205461336d57805490680100000000000000008210156129ac57826133566133218460018096018555846132cb565b81939154907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9060031b92831b921b19161790565b905580549260005201602052604060002055600190565b5050600090565b90600182019181600052826020526040600020548015156000146134fd577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81018181116134ce578254907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82019182116134ce57818103613497575b50505080548015613468577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff019061342982826132cb565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82549160031b1b191690555560005260205260006040812055600190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b6134b76134a761332193866132cb565b90549060031b1c928392866132cb565b9055600052836020526040600020553880806133f1565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b50505050600090565b91929015613581575081511561351a575090565b3b156135235790565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152fd5b8251909150156135945750805190602001fd5b6135d2906040519182917f08c379a0000000000000000000000000000000000000000000000000000000008352602060048401526024830190612ab1565b0390fdfea164736f6c634300081a000a",
}

var OnRampABI = OnRampMetaData.ABI

var OnRampBin = OnRampMetaData.Bin

func DeployOnRamp(auth *bind.TransactOpts, backend bind.ContractBackend, staticConfig OnRampStaticConfig, dynamicConfig OnRampDynamicConfig, destChainConfigArgs []OnRampDestChainConfigArgs) (common.Address, *types.Transaction, *OnRamp, error) {
	parsed, err := OnRampMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OnRampBin), backend, staticConfig, dynamicConfig, destChainConfigArgs)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OnRamp{address: address, abi: *parsed, OnRampCaller: OnRampCaller{contract: contract}, OnRampTransactor: OnRampTransactor{contract: contract}, OnRampFilterer: OnRampFilterer{contract: contract}}, nil
}

type OnRamp struct {
	address common.Address
	abi     abi.ABI
	OnRampCaller
	OnRampTransactor
	OnRampFilterer
}

type OnRampCaller struct {
	contract *bind.BoundContract
}

type OnRampTransactor struct {
	contract *bind.BoundContract
}

type OnRampFilterer struct {
	contract *bind.BoundContract
}

type OnRampSession struct {
	Contract     *OnRamp
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type OnRampCallerSession struct {
	Contract *OnRampCaller
	CallOpts bind.CallOpts
}

type OnRampTransactorSession struct {
	Contract     *OnRampTransactor
	TransactOpts bind.TransactOpts
}

type OnRampRaw struct {
	Contract *OnRamp
}

type OnRampCallerRaw struct {
	Contract *OnRampCaller
}

type OnRampTransactorRaw struct {
	Contract *OnRampTransactor
}

func NewOnRamp(address common.Address, backend bind.ContractBackend) (*OnRamp, error) {
	abi, err := abi.JSON(strings.NewReader(OnRampABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindOnRamp(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OnRamp{address: address, abi: abi, OnRampCaller: OnRampCaller{contract: contract}, OnRampTransactor: OnRampTransactor{contract: contract}, OnRampFilterer: OnRampFilterer{contract: contract}}, nil
}

func NewOnRampCaller(address common.Address, caller bind.ContractCaller) (*OnRampCaller, error) {
	contract, err := bindOnRamp(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OnRampCaller{contract: contract}, nil
}

func NewOnRampTransactor(address common.Address, transactor bind.ContractTransactor) (*OnRampTransactor, error) {
	contract, err := bindOnRamp(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OnRampTransactor{contract: contract}, nil
}

func NewOnRampFilterer(address common.Address, filterer bind.ContractFilterer) (*OnRampFilterer, error) {
	contract, err := bindOnRamp(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OnRampFilterer{contract: contract}, nil
}

func bindOnRamp(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OnRampMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_OnRamp *OnRampRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OnRamp.Contract.OnRampCaller.contract.Call(opts, result, method, params...)
}

func (_OnRamp *OnRampRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OnRamp.Contract.OnRampTransactor.contract.Transfer(opts)
}

func (_OnRamp *OnRampRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OnRamp.Contract.OnRampTransactor.contract.Transact(opts, method, params...)
}

func (_OnRamp *OnRampCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OnRamp.Contract.contract.Call(opts, result, method, params...)
}

func (_OnRamp *OnRampTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OnRamp.Contract.contract.Transfer(opts)
}

func (_OnRamp *OnRampTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OnRamp.Contract.contract.Transact(opts, method, params...)
}

func (_OnRamp *OnRampCaller) GetAllowedSendersList(opts *bind.CallOpts, destChainSelector uint64) (GetAllowedSendersList,

	error) {
	var out []interface{}
	err := _OnRamp.contract.Call(opts, &out, "getAllowedSendersList", destChainSelector)

	outstruct := new(GetAllowedSendersList)
	if err != nil {
		return *outstruct, err
	}

	outstruct.IsEnabled = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfiguredAddresses = *abi.ConvertType(out[1], new([]common.Address)).(*[]common.Address)

	return *outstruct, err

}

func (_OnRamp *OnRampSession) GetAllowedSendersList(destChainSelector uint64) (GetAllowedSendersList,

	error) {
	return _OnRamp.Contract.GetAllowedSendersList(&_OnRamp.CallOpts, destChainSelector)
}

func (_OnRamp *OnRampCallerSession) GetAllowedSendersList(destChainSelector uint64) (GetAllowedSendersList,

	error) {
	return _OnRamp.Contract.GetAllowedSendersList(&_OnRamp.CallOpts, destChainSelector)
}

func (_OnRamp *OnRampCaller) GetDestChainConfig(opts *bind.CallOpts, destChainSelector uint64) (GetDestChainConfig,

	error) {
	var out []interface{}
	err := _OnRamp.contract.Call(opts, &out, "getDestChainConfig", destChainSelector)

	outstruct := new(GetDestChainConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.SequenceNumber = *abi.ConvertType(out[0], new(uint64)).(*uint64)
	outstruct.AllowlistEnabled = *abi.ConvertType(out[1], new(bool)).(*bool)
	outstruct.Router = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)

	return *outstruct, err

}

func (_OnRamp *OnRampSession) GetDestChainConfig(destChainSelector uint64) (GetDestChainConfig,

	error) {
	return _OnRamp.Contract.GetDestChainConfig(&_OnRamp.CallOpts, destChainSelector)
}

func (_OnRamp *OnRampCallerSession) GetDestChainConfig(destChainSelector uint64) (GetDestChainConfig,

	error) {
	return _OnRamp.Contract.GetDestChainConfig(&_OnRamp.CallOpts, destChainSelector)
}

func (_OnRamp *OnRampCaller) GetDynamicConfig(opts *bind.CallOpts) (OnRampDynamicConfig, error) {
	var out []interface{}
	err := _OnRamp.contract.Call(opts, &out, "getDynamicConfig")

	if err != nil {
		return *new(OnRampDynamicConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(OnRampDynamicConfig)).(*OnRampDynamicConfig)

	return out0, err

}

func (_OnRamp *OnRampSession) GetDynamicConfig() (OnRampDynamicConfig, error) {
	return _OnRamp.Contract.GetDynamicConfig(&_OnRamp.CallOpts)
}

func (_OnRamp *OnRampCallerSession) GetDynamicConfig() (OnRampDynamicConfig, error) {
	return _OnRamp.Contract.GetDynamicConfig(&_OnRamp.CallOpts)
}

func (_OnRamp *OnRampCaller) GetExpectedNextSequenceNumber(opts *bind.CallOpts, destChainSelector uint64) (uint64, error) {
	var out []interface{}
	err := _OnRamp.contract.Call(opts, &out, "getExpectedNextSequenceNumber", destChainSelector)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_OnRamp *OnRampSession) GetExpectedNextSequenceNumber(destChainSelector uint64) (uint64, error) {
	return _OnRamp.Contract.GetExpectedNextSequenceNumber(&_OnRamp.CallOpts, destChainSelector)
}

func (_OnRamp *OnRampCallerSession) GetExpectedNextSequenceNumber(destChainSelector uint64) (uint64, error) {
	return _OnRamp.Contract.GetExpectedNextSequenceNumber(&_OnRamp.CallOpts, destChainSelector)
}

func (_OnRamp *OnRampCaller) GetFee(opts *bind.CallOpts, destChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error) {
	var out []interface{}
	err := _OnRamp.contract.Call(opts, &out, "getFee", destChainSelector, message)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_OnRamp *OnRampSession) GetFee(destChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error) {
	return _OnRamp.Contract.GetFee(&_OnRamp.CallOpts, destChainSelector, message)
}

func (_OnRamp *OnRampCallerSession) GetFee(destChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error) {
	return _OnRamp.Contract.GetFee(&_OnRamp.CallOpts, destChainSelector, message)
}

func (_OnRamp *OnRampCaller) GetPoolBySourceToken(opts *bind.CallOpts, arg0 uint64, sourceToken common.Address) (common.Address, error) {
	var out []interface{}
	err := _OnRamp.contract.Call(opts, &out, "getPoolBySourceToken", arg0, sourceToken)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_OnRamp *OnRampSession) GetPoolBySourceToken(arg0 uint64, sourceToken common.Address) (common.Address, error) {
	return _OnRamp.Contract.GetPoolBySourceToken(&_OnRamp.CallOpts, arg0, sourceToken)
}

func (_OnRamp *OnRampCallerSession) GetPoolBySourceToken(arg0 uint64, sourceToken common.Address) (common.Address, error) {
	return _OnRamp.Contract.GetPoolBySourceToken(&_OnRamp.CallOpts, arg0, sourceToken)
}

func (_OnRamp *OnRampCaller) GetStaticConfig(opts *bind.CallOpts) (OnRampStaticConfig, error) {
	var out []interface{}
	err := _OnRamp.contract.Call(opts, &out, "getStaticConfig")

	if err != nil {
		return *new(OnRampStaticConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(OnRampStaticConfig)).(*OnRampStaticConfig)

	return out0, err

}

func (_OnRamp *OnRampSession) GetStaticConfig() (OnRampStaticConfig, error) {
	return _OnRamp.Contract.GetStaticConfig(&_OnRamp.CallOpts)
}

func (_OnRamp *OnRampCallerSession) GetStaticConfig() (OnRampStaticConfig, error) {
	return _OnRamp.Contract.GetStaticConfig(&_OnRamp.CallOpts)
}

func (_OnRamp *OnRampCaller) GetSupportedTokens(opts *bind.CallOpts, arg0 uint64) ([]common.Address, error) {
	var out []interface{}
	err := _OnRamp.contract.Call(opts, &out, "getSupportedTokens", arg0)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_OnRamp *OnRampSession) GetSupportedTokens(arg0 uint64) ([]common.Address, error) {
	return _OnRamp.Contract.GetSupportedTokens(&_OnRamp.CallOpts, arg0)
}

func (_OnRamp *OnRampCallerSession) GetSupportedTokens(arg0 uint64) ([]common.Address, error) {
	return _OnRamp.Contract.GetSupportedTokens(&_OnRamp.CallOpts, arg0)
}

func (_OnRamp *OnRampCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OnRamp.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_OnRamp *OnRampSession) Owner() (common.Address, error) {
	return _OnRamp.Contract.Owner(&_OnRamp.CallOpts)
}

func (_OnRamp *OnRampCallerSession) Owner() (common.Address, error) {
	return _OnRamp.Contract.Owner(&_OnRamp.CallOpts)
}

func (_OnRamp *OnRampCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _OnRamp.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_OnRamp *OnRampSession) TypeAndVersion() (string, error) {
	return _OnRamp.Contract.TypeAndVersion(&_OnRamp.CallOpts)
}

func (_OnRamp *OnRampCallerSession) TypeAndVersion() (string, error) {
	return _OnRamp.Contract.TypeAndVersion(&_OnRamp.CallOpts)
}

func (_OnRamp *OnRampTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OnRamp.contract.Transact(opts, "acceptOwnership")
}

func (_OnRamp *OnRampSession) AcceptOwnership() (*types.Transaction, error) {
	return _OnRamp.Contract.AcceptOwnership(&_OnRamp.TransactOpts)
}

func (_OnRamp *OnRampTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _OnRamp.Contract.AcceptOwnership(&_OnRamp.TransactOpts)
}

func (_OnRamp *OnRampTransactor) ApplyAllowlistUpdates(opts *bind.TransactOpts, allowlistConfigArgsItems []OnRampAllowlistConfigArgs) (*types.Transaction, error) {
	return _OnRamp.contract.Transact(opts, "applyAllowlistUpdates", allowlistConfigArgsItems)
}

func (_OnRamp *OnRampSession) ApplyAllowlistUpdates(allowlistConfigArgsItems []OnRampAllowlistConfigArgs) (*types.Transaction, error) {
	return _OnRamp.Contract.ApplyAllowlistUpdates(&_OnRamp.TransactOpts, allowlistConfigArgsItems)
}

func (_OnRamp *OnRampTransactorSession) ApplyAllowlistUpdates(allowlistConfigArgsItems []OnRampAllowlistConfigArgs) (*types.Transaction, error) {
	return _OnRamp.Contract.ApplyAllowlistUpdates(&_OnRamp.TransactOpts, allowlistConfigArgsItems)
}

func (_OnRamp *OnRampTransactor) ApplyDestChainConfigUpdates(opts *bind.TransactOpts, destChainConfigArgs []OnRampDestChainConfigArgs) (*types.Transaction, error) {
	return _OnRamp.contract.Transact(opts, "applyDestChainConfigUpdates", destChainConfigArgs)
}

func (_OnRamp *OnRampSession) ApplyDestChainConfigUpdates(destChainConfigArgs []OnRampDestChainConfigArgs) (*types.Transaction, error) {
	return _OnRamp.Contract.ApplyDestChainConfigUpdates(&_OnRamp.TransactOpts, destChainConfigArgs)
}

func (_OnRamp *OnRampTransactorSession) ApplyDestChainConfigUpdates(destChainConfigArgs []OnRampDestChainConfigArgs) (*types.Transaction, error) {
	return _OnRamp.Contract.ApplyDestChainConfigUpdates(&_OnRamp.TransactOpts, destChainConfigArgs)
}

func (_OnRamp *OnRampTransactor) ForwardFromRouter(opts *bind.TransactOpts, destChainSelector uint64, message ClientEVM2AnyMessage, feeTokenAmount *big.Int, originalSender common.Address) (*types.Transaction, error) {
	return _OnRamp.contract.Transact(opts, "forwardFromRouter", destChainSelector, message, feeTokenAmount, originalSender)
}

func (_OnRamp *OnRampSession) ForwardFromRouter(destChainSelector uint64, message ClientEVM2AnyMessage, feeTokenAmount *big.Int, originalSender common.Address) (*types.Transaction, error) {
	return _OnRamp.Contract.ForwardFromRouter(&_OnRamp.TransactOpts, destChainSelector, message, feeTokenAmount, originalSender)
}

func (_OnRamp *OnRampTransactorSession) ForwardFromRouter(destChainSelector uint64, message ClientEVM2AnyMessage, feeTokenAmount *big.Int, originalSender common.Address) (*types.Transaction, error) {
	return _OnRamp.Contract.ForwardFromRouter(&_OnRamp.TransactOpts, destChainSelector, message, feeTokenAmount, originalSender)
}

func (_OnRamp *OnRampTransactor) SetDynamicConfig(opts *bind.TransactOpts, dynamicConfig OnRampDynamicConfig) (*types.Transaction, error) {
	return _OnRamp.contract.Transact(opts, "setDynamicConfig", dynamicConfig)
}

func (_OnRamp *OnRampSession) SetDynamicConfig(dynamicConfig OnRampDynamicConfig) (*types.Transaction, error) {
	return _OnRamp.Contract.SetDynamicConfig(&_OnRamp.TransactOpts, dynamicConfig)
}

func (_OnRamp *OnRampTransactorSession) SetDynamicConfig(dynamicConfig OnRampDynamicConfig) (*types.Transaction, error) {
	return _OnRamp.Contract.SetDynamicConfig(&_OnRamp.TransactOpts, dynamicConfig)
}

func (_OnRamp *OnRampTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _OnRamp.contract.Transact(opts, "transferOwnership", to)
}

func (_OnRamp *OnRampSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OnRamp.Contract.TransferOwnership(&_OnRamp.TransactOpts, to)
}

func (_OnRamp *OnRampTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OnRamp.Contract.TransferOwnership(&_OnRamp.TransactOpts, to)
}

func (_OnRamp *OnRampTransactor) WithdrawFeeTokens(opts *bind.TransactOpts, feeTokens []common.Address) (*types.Transaction, error) {
	return _OnRamp.contract.Transact(opts, "withdrawFeeTokens", feeTokens)
}

func (_OnRamp *OnRampSession) WithdrawFeeTokens(feeTokens []common.Address) (*types.Transaction, error) {
	return _OnRamp.Contract.WithdrawFeeTokens(&_OnRamp.TransactOpts, feeTokens)
}

func (_OnRamp *OnRampTransactorSession) WithdrawFeeTokens(feeTokens []common.Address) (*types.Transaction, error) {
	return _OnRamp.Contract.WithdrawFeeTokens(&_OnRamp.TransactOpts, feeTokens)
}

type OnRampAllowListAdminSetIterator struct {
	Event *OnRampAllowListAdminSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OnRampAllowListAdminSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnRampAllowListAdminSet)
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
		it.Event = new(OnRampAllowListAdminSet)
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

func (it *OnRampAllowListAdminSetIterator) Error() error {
	return it.fail
}

func (it *OnRampAllowListAdminSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OnRampAllowListAdminSet struct {
	AllowlistAdmin common.Address
	Raw            types.Log
}

func (_OnRamp *OnRampFilterer) FilterAllowListAdminSet(opts *bind.FilterOpts, allowlistAdmin []common.Address) (*OnRampAllowListAdminSetIterator, error) {

	var allowlistAdminRule []interface{}
	for _, allowlistAdminItem := range allowlistAdmin {
		allowlistAdminRule = append(allowlistAdminRule, allowlistAdminItem)
	}

	logs, sub, err := _OnRamp.contract.FilterLogs(opts, "AllowListAdminSet", allowlistAdminRule)
	if err != nil {
		return nil, err
	}
	return &OnRampAllowListAdminSetIterator{contract: _OnRamp.contract, event: "AllowListAdminSet", logs: logs, sub: sub}, nil
}

func (_OnRamp *OnRampFilterer) WatchAllowListAdminSet(opts *bind.WatchOpts, sink chan<- *OnRampAllowListAdminSet, allowlistAdmin []common.Address) (event.Subscription, error) {

	var allowlistAdminRule []interface{}
	for _, allowlistAdminItem := range allowlistAdmin {
		allowlistAdminRule = append(allowlistAdminRule, allowlistAdminItem)
	}

	logs, sub, err := _OnRamp.contract.WatchLogs(opts, "AllowListAdminSet", allowlistAdminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OnRampAllowListAdminSet)
				if err := _OnRamp.contract.UnpackLog(event, "AllowListAdminSet", log); err != nil {
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

func (_OnRamp *OnRampFilterer) ParseAllowListAdminSet(log types.Log) (*OnRampAllowListAdminSet, error) {
	event := new(OnRampAllowListAdminSet)
	if err := _OnRamp.contract.UnpackLog(event, "AllowListAdminSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OnRampAllowListSendersAddedIterator struct {
	Event *OnRampAllowListSendersAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OnRampAllowListSendersAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnRampAllowListSendersAdded)
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
		it.Event = new(OnRampAllowListSendersAdded)
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

func (it *OnRampAllowListSendersAddedIterator) Error() error {
	return it.fail
}

func (it *OnRampAllowListSendersAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OnRampAllowListSendersAdded struct {
	DestChainSelector uint64
	Senders           []common.Address
	Raw               types.Log
}

func (_OnRamp *OnRampFilterer) FilterAllowListSendersAdded(opts *bind.FilterOpts, destChainSelector []uint64) (*OnRampAllowListSendersAddedIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _OnRamp.contract.FilterLogs(opts, "AllowListSendersAdded", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &OnRampAllowListSendersAddedIterator{contract: _OnRamp.contract, event: "AllowListSendersAdded", logs: logs, sub: sub}, nil
}

func (_OnRamp *OnRampFilterer) WatchAllowListSendersAdded(opts *bind.WatchOpts, sink chan<- *OnRampAllowListSendersAdded, destChainSelector []uint64) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _OnRamp.contract.WatchLogs(opts, "AllowListSendersAdded", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OnRampAllowListSendersAdded)
				if err := _OnRamp.contract.UnpackLog(event, "AllowListSendersAdded", log); err != nil {
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

func (_OnRamp *OnRampFilterer) ParseAllowListSendersAdded(log types.Log) (*OnRampAllowListSendersAdded, error) {
	event := new(OnRampAllowListSendersAdded)
	if err := _OnRamp.contract.UnpackLog(event, "AllowListSendersAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OnRampAllowListSendersRemovedIterator struct {
	Event *OnRampAllowListSendersRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OnRampAllowListSendersRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnRampAllowListSendersRemoved)
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
		it.Event = new(OnRampAllowListSendersRemoved)
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

func (it *OnRampAllowListSendersRemovedIterator) Error() error {
	return it.fail
}

func (it *OnRampAllowListSendersRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OnRampAllowListSendersRemoved struct {
	DestChainSelector uint64
	Senders           []common.Address
	Raw               types.Log
}

func (_OnRamp *OnRampFilterer) FilterAllowListSendersRemoved(opts *bind.FilterOpts, destChainSelector []uint64) (*OnRampAllowListSendersRemovedIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _OnRamp.contract.FilterLogs(opts, "AllowListSendersRemoved", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &OnRampAllowListSendersRemovedIterator{contract: _OnRamp.contract, event: "AllowListSendersRemoved", logs: logs, sub: sub}, nil
}

func (_OnRamp *OnRampFilterer) WatchAllowListSendersRemoved(opts *bind.WatchOpts, sink chan<- *OnRampAllowListSendersRemoved, destChainSelector []uint64) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _OnRamp.contract.WatchLogs(opts, "AllowListSendersRemoved", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OnRampAllowListSendersRemoved)
				if err := _OnRamp.contract.UnpackLog(event, "AllowListSendersRemoved", log); err != nil {
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

func (_OnRamp *OnRampFilterer) ParseAllowListSendersRemoved(log types.Log) (*OnRampAllowListSendersRemoved, error) {
	event := new(OnRampAllowListSendersRemoved)
	if err := _OnRamp.contract.UnpackLog(event, "AllowListSendersRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OnRampCCIPMessageSentIterator struct {
	Event *OnRampCCIPMessageSent

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OnRampCCIPMessageSentIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnRampCCIPMessageSent)
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
		it.Event = new(OnRampCCIPMessageSent)
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

func (it *OnRampCCIPMessageSentIterator) Error() error {
	return it.fail
}

func (it *OnRampCCIPMessageSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OnRampCCIPMessageSent struct {
	DestChainSelector uint64
	SequenceNumber    uint64
	Message           InternalEVM2AnyRampMessage
	Raw               types.Log
}

func (_OnRamp *OnRampFilterer) FilterCCIPMessageSent(opts *bind.FilterOpts, destChainSelector []uint64, sequenceNumber []uint64) (*OnRampCCIPMessageSentIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}
	var sequenceNumberRule []interface{}
	for _, sequenceNumberItem := range sequenceNumber {
		sequenceNumberRule = append(sequenceNumberRule, sequenceNumberItem)
	}

	logs, sub, err := _OnRamp.contract.FilterLogs(opts, "CCIPMessageSent", destChainSelectorRule, sequenceNumberRule)
	if err != nil {
		return nil, err
	}
	return &OnRampCCIPMessageSentIterator{contract: _OnRamp.contract, event: "CCIPMessageSent", logs: logs, sub: sub}, nil
}

func (_OnRamp *OnRampFilterer) WatchCCIPMessageSent(opts *bind.WatchOpts, sink chan<- *OnRampCCIPMessageSent, destChainSelector []uint64, sequenceNumber []uint64) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}
	var sequenceNumberRule []interface{}
	for _, sequenceNumberItem := range sequenceNumber {
		sequenceNumberRule = append(sequenceNumberRule, sequenceNumberItem)
	}

	logs, sub, err := _OnRamp.contract.WatchLogs(opts, "CCIPMessageSent", destChainSelectorRule, sequenceNumberRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OnRampCCIPMessageSent)
				if err := _OnRamp.contract.UnpackLog(event, "CCIPMessageSent", log); err != nil {
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

func (_OnRamp *OnRampFilterer) ParseCCIPMessageSent(log types.Log) (*OnRampCCIPMessageSent, error) {
	event := new(OnRampCCIPMessageSent)
	if err := _OnRamp.contract.UnpackLog(event, "CCIPMessageSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OnRampConfigSetIterator struct {
	Event *OnRampConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OnRampConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnRampConfigSet)
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
		it.Event = new(OnRampConfigSet)
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

func (it *OnRampConfigSetIterator) Error() error {
	return it.fail
}

func (it *OnRampConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OnRampConfigSet struct {
	StaticConfig  OnRampStaticConfig
	DynamicConfig OnRampDynamicConfig
	Raw           types.Log
}

func (_OnRamp *OnRampFilterer) FilterConfigSet(opts *bind.FilterOpts) (*OnRampConfigSetIterator, error) {

	logs, sub, err := _OnRamp.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &OnRampConfigSetIterator{contract: _OnRamp.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_OnRamp *OnRampFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *OnRampConfigSet) (event.Subscription, error) {

	logs, sub, err := _OnRamp.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OnRampConfigSet)
				if err := _OnRamp.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_OnRamp *OnRampFilterer) ParseConfigSet(log types.Log) (*OnRampConfigSet, error) {
	event := new(OnRampConfigSet)
	if err := _OnRamp.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OnRampDestChainConfigSetIterator struct {
	Event *OnRampDestChainConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OnRampDestChainConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnRampDestChainConfigSet)
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
		it.Event = new(OnRampDestChainConfigSet)
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

func (it *OnRampDestChainConfigSetIterator) Error() error {
	return it.fail
}

func (it *OnRampDestChainConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OnRampDestChainConfigSet struct {
	DestChainSelector uint64
	SequenceNumber    uint64
	Router            common.Address
	AllowlistEnabled  bool
	Raw               types.Log
}

func (_OnRamp *OnRampFilterer) FilterDestChainConfigSet(opts *bind.FilterOpts, destChainSelector []uint64) (*OnRampDestChainConfigSetIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _OnRamp.contract.FilterLogs(opts, "DestChainConfigSet", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &OnRampDestChainConfigSetIterator{contract: _OnRamp.contract, event: "DestChainConfigSet", logs: logs, sub: sub}, nil
}

func (_OnRamp *OnRampFilterer) WatchDestChainConfigSet(opts *bind.WatchOpts, sink chan<- *OnRampDestChainConfigSet, destChainSelector []uint64) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _OnRamp.contract.WatchLogs(opts, "DestChainConfigSet", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OnRampDestChainConfigSet)
				if err := _OnRamp.contract.UnpackLog(event, "DestChainConfigSet", log); err != nil {
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

func (_OnRamp *OnRampFilterer) ParseDestChainConfigSet(log types.Log) (*OnRampDestChainConfigSet, error) {
	event := new(OnRampDestChainConfigSet)
	if err := _OnRamp.contract.UnpackLog(event, "DestChainConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OnRampFeeTokenWithdrawnIterator struct {
	Event *OnRampFeeTokenWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OnRampFeeTokenWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnRampFeeTokenWithdrawn)
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
		it.Event = new(OnRampFeeTokenWithdrawn)
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

func (it *OnRampFeeTokenWithdrawnIterator) Error() error {
	return it.fail
}

func (it *OnRampFeeTokenWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OnRampFeeTokenWithdrawn struct {
	FeeAggregator common.Address
	FeeToken      common.Address
	Amount        *big.Int
	Raw           types.Log
}

func (_OnRamp *OnRampFilterer) FilterFeeTokenWithdrawn(opts *bind.FilterOpts, feeAggregator []common.Address, feeToken []common.Address) (*OnRampFeeTokenWithdrawnIterator, error) {

	var feeAggregatorRule []interface{}
	for _, feeAggregatorItem := range feeAggregator {
		feeAggregatorRule = append(feeAggregatorRule, feeAggregatorItem)
	}
	var feeTokenRule []interface{}
	for _, feeTokenItem := range feeToken {
		feeTokenRule = append(feeTokenRule, feeTokenItem)
	}

	logs, sub, err := _OnRamp.contract.FilterLogs(opts, "FeeTokenWithdrawn", feeAggregatorRule, feeTokenRule)
	if err != nil {
		return nil, err
	}
	return &OnRampFeeTokenWithdrawnIterator{contract: _OnRamp.contract, event: "FeeTokenWithdrawn", logs: logs, sub: sub}, nil
}

func (_OnRamp *OnRampFilterer) WatchFeeTokenWithdrawn(opts *bind.WatchOpts, sink chan<- *OnRampFeeTokenWithdrawn, feeAggregator []common.Address, feeToken []common.Address) (event.Subscription, error) {

	var feeAggregatorRule []interface{}
	for _, feeAggregatorItem := range feeAggregator {
		feeAggregatorRule = append(feeAggregatorRule, feeAggregatorItem)
	}
	var feeTokenRule []interface{}
	for _, feeTokenItem := range feeToken {
		feeTokenRule = append(feeTokenRule, feeTokenItem)
	}

	logs, sub, err := _OnRamp.contract.WatchLogs(opts, "FeeTokenWithdrawn", feeAggregatorRule, feeTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OnRampFeeTokenWithdrawn)
				if err := _OnRamp.contract.UnpackLog(event, "FeeTokenWithdrawn", log); err != nil {
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

func (_OnRamp *OnRampFilterer) ParseFeeTokenWithdrawn(log types.Log) (*OnRampFeeTokenWithdrawn, error) {
	event := new(OnRampFeeTokenWithdrawn)
	if err := _OnRamp.contract.UnpackLog(event, "FeeTokenWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OnRampOwnershipTransferRequestedIterator struct {
	Event *OnRampOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OnRampOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnRampOwnershipTransferRequested)
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
		it.Event = new(OnRampOwnershipTransferRequested)
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

func (it *OnRampOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *OnRampOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OnRampOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_OnRamp *OnRampFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OnRampOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OnRamp.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OnRampOwnershipTransferRequestedIterator{contract: _OnRamp.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_OnRamp *OnRampFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OnRampOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OnRamp.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OnRampOwnershipTransferRequested)
				if err := _OnRamp.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_OnRamp *OnRampFilterer) ParseOwnershipTransferRequested(log types.Log) (*OnRampOwnershipTransferRequested, error) {
	event := new(OnRampOwnershipTransferRequested)
	if err := _OnRamp.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OnRampOwnershipTransferredIterator struct {
	Event *OnRampOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OnRampOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnRampOwnershipTransferred)
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
		it.Event = new(OnRampOwnershipTransferred)
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

func (it *OnRampOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *OnRampOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OnRampOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_OnRamp *OnRampFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OnRampOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OnRamp.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OnRampOwnershipTransferredIterator{contract: _OnRamp.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_OnRamp *OnRampFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OnRampOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OnRamp.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OnRampOwnershipTransferred)
				if err := _OnRamp.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_OnRamp *OnRampFilterer) ParseOwnershipTransferred(log types.Log) (*OnRampOwnershipTransferred, error) {
	event := new(OnRampOwnershipTransferred)
	if err := _OnRamp.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetAllowedSendersList struct {
	IsEnabled           bool
	ConfiguredAddresses []common.Address
}
type GetDestChainConfig struct {
	SequenceNumber   uint64
	AllowlistEnabled bool
	Router           common.Address
}

func (_OnRamp *OnRamp) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _OnRamp.abi.Events["AllowListAdminSet"].ID:
		return _OnRamp.ParseAllowListAdminSet(log)
	case _OnRamp.abi.Events["AllowListSendersAdded"].ID:
		return _OnRamp.ParseAllowListSendersAdded(log)
	case _OnRamp.abi.Events["AllowListSendersRemoved"].ID:
		return _OnRamp.ParseAllowListSendersRemoved(log)
	case _OnRamp.abi.Events["CCIPMessageSent"].ID:
		return _OnRamp.ParseCCIPMessageSent(log)
	case _OnRamp.abi.Events["ConfigSet"].ID:
		return _OnRamp.ParseConfigSet(log)
	case _OnRamp.abi.Events["DestChainConfigSet"].ID:
		return _OnRamp.ParseDestChainConfigSet(log)
	case _OnRamp.abi.Events["FeeTokenWithdrawn"].ID:
		return _OnRamp.ParseFeeTokenWithdrawn(log)
	case _OnRamp.abi.Events["OwnershipTransferRequested"].ID:
		return _OnRamp.ParseOwnershipTransferRequested(log)
	case _OnRamp.abi.Events["OwnershipTransferred"].ID:
		return _OnRamp.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (OnRampAllowListAdminSet) Topic() common.Hash {
	return common.HexToHash("0xb8c9b44ae5b5e3afb195f67391d9ff50cb904f9c0fa5fd520e497a97c1aa5a1e")
}

func (OnRampAllowListSendersAdded) Topic() common.Hash {
	return common.HexToHash("0x330939f6eafe8bb516716892fe962ff19770570838686e6579dbc1cc51fc3281")
}

func (OnRampAllowListSendersRemoved) Topic() common.Hash {
	return common.HexToHash("0xc237ec1921f855ccd5e9a5af9733f2d58943a5a8501ec5988e305d7a4d421586")
}

func (OnRampCCIPMessageSent) Topic() common.Hash {
	return common.HexToHash("0x192442a2b2adb6a7948f097023cb6b57d29d3a7a5dd33e6666d33c39cc456f32")
}

func (OnRampConfigSet) Topic() common.Hash {
	return common.HexToHash("0xc7372d2d886367d7bb1b0e0708a5436f2c91d6963de210eb2dc1ec2ecd6d21f1")
}

func (OnRampDestChainConfigSet) Topic() common.Hash {
	return common.HexToHash("0xd5ad72bc37dc7a80a8b9b9df20500046fd7341adb1be2258a540466fdd7dcef5")
}

func (OnRampFeeTokenWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x508d7d183612c18fc339b42618912b9fa3239f631dd7ec0671f950200a0fa66e")
}

func (OnRampOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (OnRampOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_OnRamp *OnRamp) Address() common.Address {
	return _OnRamp.address
}

type OnRampInterface interface {
	GetAllowedSendersList(opts *bind.CallOpts, destChainSelector uint64) (GetAllowedSendersList,

		error)

	GetDestChainConfig(opts *bind.CallOpts, destChainSelector uint64) (GetDestChainConfig,

		error)

	GetDynamicConfig(opts *bind.CallOpts) (OnRampDynamicConfig, error)

	GetExpectedNextSequenceNumber(opts *bind.CallOpts, destChainSelector uint64) (uint64, error)

	GetFee(opts *bind.CallOpts, destChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error)

	GetPoolBySourceToken(opts *bind.CallOpts, arg0 uint64, sourceToken common.Address) (common.Address, error)

	GetStaticConfig(opts *bind.CallOpts) (OnRampStaticConfig, error)

	GetSupportedTokens(opts *bind.CallOpts, arg0 uint64) ([]common.Address, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	ApplyAllowlistUpdates(opts *bind.TransactOpts, allowlistConfigArgsItems []OnRampAllowlistConfigArgs) (*types.Transaction, error)

	ApplyDestChainConfigUpdates(opts *bind.TransactOpts, destChainConfigArgs []OnRampDestChainConfigArgs) (*types.Transaction, error)

	ForwardFromRouter(opts *bind.TransactOpts, destChainSelector uint64, message ClientEVM2AnyMessage, feeTokenAmount *big.Int, originalSender common.Address) (*types.Transaction, error)

	SetDynamicConfig(opts *bind.TransactOpts, dynamicConfig OnRampDynamicConfig) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	WithdrawFeeTokens(opts *bind.TransactOpts, feeTokens []common.Address) (*types.Transaction, error)

	FilterAllowListAdminSet(opts *bind.FilterOpts, allowlistAdmin []common.Address) (*OnRampAllowListAdminSetIterator, error)

	WatchAllowListAdminSet(opts *bind.WatchOpts, sink chan<- *OnRampAllowListAdminSet, allowlistAdmin []common.Address) (event.Subscription, error)

	ParseAllowListAdminSet(log types.Log) (*OnRampAllowListAdminSet, error)

	FilterAllowListSendersAdded(opts *bind.FilterOpts, destChainSelector []uint64) (*OnRampAllowListSendersAddedIterator, error)

	WatchAllowListSendersAdded(opts *bind.WatchOpts, sink chan<- *OnRampAllowListSendersAdded, destChainSelector []uint64) (event.Subscription, error)

	ParseAllowListSendersAdded(log types.Log) (*OnRampAllowListSendersAdded, error)

	FilterAllowListSendersRemoved(opts *bind.FilterOpts, destChainSelector []uint64) (*OnRampAllowListSendersRemovedIterator, error)

	WatchAllowListSendersRemoved(opts *bind.WatchOpts, sink chan<- *OnRampAllowListSendersRemoved, destChainSelector []uint64) (event.Subscription, error)

	ParseAllowListSendersRemoved(log types.Log) (*OnRampAllowListSendersRemoved, error)

	FilterCCIPMessageSent(opts *bind.FilterOpts, destChainSelector []uint64, sequenceNumber []uint64) (*OnRampCCIPMessageSentIterator, error)

	WatchCCIPMessageSent(opts *bind.WatchOpts, sink chan<- *OnRampCCIPMessageSent, destChainSelector []uint64, sequenceNumber []uint64) (event.Subscription, error)

	ParseCCIPMessageSent(log types.Log) (*OnRampCCIPMessageSent, error)

	FilterConfigSet(opts *bind.FilterOpts) (*OnRampConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *OnRampConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*OnRampConfigSet, error)

	FilterDestChainConfigSet(opts *bind.FilterOpts, destChainSelector []uint64) (*OnRampDestChainConfigSetIterator, error)

	WatchDestChainConfigSet(opts *bind.WatchOpts, sink chan<- *OnRampDestChainConfigSet, destChainSelector []uint64) (event.Subscription, error)

	ParseDestChainConfigSet(log types.Log) (*OnRampDestChainConfigSet, error)

	FilterFeeTokenWithdrawn(opts *bind.FilterOpts, feeAggregator []common.Address, feeToken []common.Address) (*OnRampFeeTokenWithdrawnIterator, error)

	WatchFeeTokenWithdrawn(opts *bind.WatchOpts, sink chan<- *OnRampFeeTokenWithdrawn, feeAggregator []common.Address, feeToken []common.Address) (event.Subscription, error)

	ParseFeeTokenWithdrawn(log types.Log) (*OnRampFeeTokenWithdrawn, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OnRampOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OnRampOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*OnRampOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OnRampOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OnRampOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*OnRampOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
