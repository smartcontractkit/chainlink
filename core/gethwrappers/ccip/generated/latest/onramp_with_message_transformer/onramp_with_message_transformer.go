// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package onramp_with_message_transformer

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

var OnRampWithMessageTransformerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"staticConfig\",\"type\":\"tuple\",\"internalType\":\"structOnRamp.StaticConfig\",\"components\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"rmnRemote\",\"type\":\"address\",\"internalType\":\"contractIRMNRemote\"},{\"name\":\"nonceManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenAdminRegistry\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"name\":\"dynamicConfig\",\"type\":\"tuple\",\"internalType\":\"structOnRamp.DynamicConfig\",\"components\":[{\"name\":\"feeQuoter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"reentrancyGuardEntered\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"messageInterceptor\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"feeAggregator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowlistAdmin\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"name\":\"destChainConfigs\",\"type\":\"tuple[]\",\"internalType\":\"structOnRamp.DestChainConfigArgs[]\",\"components\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"contractIRouter\"},{\"name\":\"allowlistEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"}]},{\"name\":\"messageTransformerAddr\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"applyAllowlistUpdates\",\"inputs\":[{\"name\":\"allowlistConfigArgsItems\",\"type\":\"tuple[]\",\"internalType\":\"structOnRamp.AllowlistConfigArgs[]\",\"components\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"allowlistEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"addedAllowlistedSenders\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"removedAllowlistedSenders\",\"type\":\"address[]\",\"internalType\":\"address[]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"applyDestChainConfigUpdates\",\"inputs\":[{\"name\":\"destChainConfigArgs\",\"type\":\"tuple[]\",\"internalType\":\"structOnRamp.DestChainConfigArgs[]\",\"components\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"contractIRouter\"},{\"name\":\"allowlistEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"forwardFromRouter\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"message\",\"type\":\"tuple\",\"internalType\":\"structClient.EVM2AnyMessage\",\"components\":[{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"tokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structClient.EVMTokenAmount[]\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"feeToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"extraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"feeTokenAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"originalSender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAllowedSendersList\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"configuredAddresses\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDestChainConfig\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"allowlistEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDynamicConfig\",\"inputs\":[],\"outputs\":[{\"name\":\"dynamicConfig\",\"type\":\"tuple\",\"internalType\":\"structOnRamp.DynamicConfig\",\"components\":[{\"name\":\"feeQuoter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"reentrancyGuardEntered\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"messageInterceptor\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"feeAggregator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowlistAdmin\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getExpectedNextSequenceNumber\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFee\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"message\",\"type\":\"tuple\",\"internalType\":\"structClient.EVM2AnyMessage\",\"components\":[{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"tokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structClient.EVMTokenAmount[]\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"feeToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"extraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"feeTokenAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMessageTransformer\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPoolBySourceToken\",\"inputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sourceToken\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIPoolV1\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStaticConfig\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structOnRamp.StaticConfig\",\"components\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"rmnRemote\",\"type\":\"address\",\"internalType\":\"contractIRMNRemote\"},{\"name\":\"nonceManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenAdminRegistry\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSupportedTokens\",\"inputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setDynamicConfig\",\"inputs\":[{\"name\":\"dynamicConfig\",\"type\":\"tuple\",\"internalType\":\"structOnRamp.DynamicConfig\",\"components\":[{\"name\":\"feeQuoter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"reentrancyGuardEntered\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"messageInterceptor\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"feeAggregator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowlistAdmin\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMessageTransformer\",\"inputs\":[{\"name\":\"messageTransformerAddr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"typeAndVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdrawFeeTokens\",\"inputs\":[{\"name\":\"feeTokens\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AllowListAdminSet\",\"inputs\":[{\"name\":\"allowlistAdmin\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AllowListSendersAdded\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"senders\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AllowListSendersRemoved\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"senders\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CCIPMessageSent\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"message\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structInternal.EVM2AnyRampMessage\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structInternal.RampMessageHeader\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"nonce\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"extraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"feeToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"feeTokenAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"feeValueJuels\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.EVM2AnyTokenTransfer[]\",\"components\":[{\"name\":\"sourcePoolAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"destTokenAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"extraData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"destExecData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ConfigSet\",\"inputs\":[{\"name\":\"staticConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structOnRamp.StaticConfig\",\"components\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"rmnRemote\",\"type\":\"address\",\"internalType\":\"contractIRMNRemote\"},{\"name\":\"nonceManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenAdminRegistry\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"name\":\"dynamicConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structOnRamp.DynamicConfig\",\"components\":[{\"name\":\"feeQuoter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"reentrancyGuardEntered\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"messageInterceptor\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"feeAggregator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowlistAdmin\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DestChainConfigSet\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"router\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"contractIRouter\"},{\"name\":\"allowlistEnabled\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FeeTokenWithdrawn\",\"inputs\":[{\"name\":\"feeAggregator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"feeToken\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"CannotSendZeroTokens\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CannotTransferToSelf\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CursedByRMN\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"GetSupportedTokensFunctionalityRemovedCheckAdminRegistry\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidAllowListRequest\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidConfig\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidDestChainConfig\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MessageTransformError\",\"inputs\":[{\"name\":\"errorReason\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"MustBeCalledByRouter\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MustBeProposedOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwnerOrAllowlistAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnerCannotBeZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RouterMustSetOriginalSender\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SenderNotAllowed\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UnsupportedToken\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ZeroAddressNotAllowed\",\"inputs\":[]}]",
	Bin: "0x610100604052346105a15761418b8038038061001a816105db565b92833981019080820361016081126105a157608081126105a15761003c6105bc565b9161004681610600565b835260208101516001600160a01b03811681036105a1576020840190815261007060408301610614565b916040850192835260a061008660608301610614565b6060870190815294607f1901126105a15760405160a081016001600160401b038111828210176105a6576040526100bf60808301610614565b81526100cd60a08301610628565b602082019081526100e060c08401610614565b90604083019182526100f460e08501610614565b92606081019384526101096101008601610614565b608082019081526101208601519095906001600160401b0381116105a15781018b601f820112156105a1578051906001600160401b0382116105a6578160051b602001610155906105db565b9c8d838152602001926060028201602001918183116105a157602001925b828410610533575050505061014061018b9101610614565b98331561052257600180546001600160a01b0319163317905580516001600160401b0316158015610510575b80156104fe575b80156104ec575b6104bf57516001600160401b0316608081905295516001600160a01b0390811660a08190529751811660c08190529851811660e081905282519091161580156104da575b80156104d0575b6104bf57815160028054855160ff60a01b90151560a01b166001600160a01b039384166001600160a81b0319909216919091171790558451600380549183166001600160a01b03199283161790558651600480549184169183169190911790558751600580549190931691161790557fc7372d2d886367d7bb1b0e0708a5436f2c91d6963de210eb2dc1ec2ecd6d21f1986101209860606102af6105bc565b8a8152602080820193845260408083019586529290910194855281519a8b5291516001600160a01b03908116928b019290925291518116918901919091529051811660608801529051811660808701529051151560a08601529051811660c08501529051811660e0840152905116610100820152a160005b82518110156104185761033a8184610635565b516001600160401b0361034d8386610635565b5151169081156104035760008281526006602090815260409182902081840151815494840151600160401b600160e81b03198616604883901b600160481b600160e81b031617901515851b68ff000000000000000016179182905583516001600160401b0390951685526001600160a01b031691840191909152811c60ff1615159082015260019291907fd5ad72bc37dc7a80a8b9b9df20500046fd7341adb1be2258a540466fdd7dcef590606090a201610327565b5063c35aa79d60e01b60005260045260246000fd5b506001600160a01b031680156104ae57600780546001600160a01b031916919091179055604051613b2b908161066082396080518181816103f901528181610b83015281816120260152612814015260a05181818161205f01528181612580015261284d015260c051818181610df40152818161209b0152612889015260e0518181816120d7015281816128c50152612ec40152f35b6342bcdf7f60e11b60005260046000fd5b6306b7c75960e31b60005260046000fd5b5082511515610210565b5084516001600160a01b031615610209565b5088516001600160a01b0316156101c5565b5087516001600160a01b0316156101be565b5086516001600160a01b0316156101b7565b639b15e16f60e01b60005260046000fd5b6060848303126105a15760405190606082016001600160401b038111838210176105a65760405261056385610600565b82526020850151906001600160a01b03821682036105a1578260209283606095015261059160408801610628565b6040820152815201930192610173565b600080fd5b634e487b7160e01b600052604160045260246000fd5b60405190608082016001600160401b038111838210176105a657604052565b6040519190601f01601f191682016001600160401b038111838210176105a657604052565b51906001600160401b03821682036105a157565b51906001600160a01b03821682036105a157565b519081151582036105a157565b80518210156106495760209160051b010190565b634e487b7160e01b600052603260045260246000fdfe608080604052600436101561001357600080fd5b600090813560e01c90816306285c69146127ae5750806315777ab21461275c578063181f5a77146126dd57806320487ded146124a45780632716072b146121f457806327e936f114611dea57806348a98aa414611d675780635cb80c5d14611aaa57806365b81aab146119fa5780636def4ce71461196b5780637437ff9f1461184e57806379ba5097146117695780638da5cb5b146117175780639041be3d1461166a578063972b46121461159c578063c9b146b3146111d3578063df0aa9e914610242578063f2fde38b146101555763fbca3b74146100f257600080fd5b346101525760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101525760049061012c612aba565b507f9e7177c8000000000000000000000000000000000000000000000000000000008152fd5b80fd5b50346101525760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101525773ffffffffffffffffffffffffffffffffffffffff6101a2612b10565b6101aa6133f0565b1633811461021a57807fffffffffffffffffffffffff000000000000000000000000000000000000000083541617825573ffffffffffffffffffffffffffffffffffffffff600154167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12788380a380f35b6004827fdad89dca000000000000000000000000000000000000000000000000000000008152fd5b50346101525760807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101525761027a612aba565b67ffffffffffffffff60243511610e585760a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc60243536030112610e58576102c1612b33565b60025460ff8160a01c166111ab577fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff16740100000000000000000000000000000000000000001760025567ffffffffffffffff8216835260066020526040832073ffffffffffffffffffffffffffffffffffffffff82161561118357805460ff8160401c16611115575b60481c73ffffffffffffffffffffffffffffffffffffffff1633036110ed5773ffffffffffffffffffffffffffffffffffffffff6003541680611079575b50805467ffffffffffffffff811667ffffffffffffffff811461104c579067ffffffffffffffff60017fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000009493011692839116179055604051906103eb82612984565b84825267ffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016602083015267ffffffffffffffff84166040830152606082015283608082015261044c602480350160243560040161305c565b909361045d6004602435018061305c565b61046b606460243501612f68565b966104806044602435016024356004016130ad565b9590507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe06104c66104b088612aeb565b976104be604051998a6129d9565b808952612aeb565b018a5b81811061103557505061051a93929161050e91604051986104e98a6129a0565b895273ffffffffffffffffffffffffffffffffffffffff8a1660208a0152369161312d565b6040870152369161312d565b606084015260209173ffffffffffffffffffffffffffffffffffffffff6040519661054585896129d9565b888852608086019788521660a085015260443560c085015260e084019087825261010085015273ffffffffffffffffffffffffffffffffffffffff878160025416610594606460243501612f68565b9067ffffffffffffffff8661064e6105b660846024350160243560040161305c565b9061061e6105c96004602435018061305c565b9290936040519b8c9a8b998a997f3a49bb49000000000000000000000000000000000000000000000000000000008b521660048a0152166024880152604435604488015260a0606488015260a4870191612c4f565b917ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc858403016084860152612c4f565b03915afa91821561102a57889189988a918b95610faf575b50525261067d6044602435016024356004016130ad565b9661068788612aeb565b9161069560405193846129d9565b888352858301809960061b820191368311610fab57905b828210610f7457505050885b6106cc6044602435016024356004016130ad565b9050811015610a21576106df8184613019565b51906106e9613101565b5086820151156109f95773ffffffffffffffffffffffffffffffffffffffff61071481845116612e65565b169182158015610956575b61091457808c878a8a73ffffffffffffffffffffffffffffffffffffffff8f836107cc9801518280895116926040519761075889612984565b885267ffffffffffffffff87890196168652816040890191168152606088019283526080880193845267ffffffffffffffff6040519b8c998a997f9a4575b9000000000000000000000000000000000000000000000000000000008b5260048b01525160a060248b015260c48a0190612a77565b965116604488015251166064860152516084850152511660a4830152038183885af1918215610909578d92610862575b506001938284928b8061085b9651930151910151916040519361081e85612984565b84528c840152604083015260608201528d6040519061083d8c836129d9565b815260808201526101008b0151906108558383613019565b52613019565b50016106b8565b91503d808e843e61087381846129d9565b82019189818403126109015780519067ffffffffffffffff821161090557019360408584031261090157604051916108aa836129bd565b855167ffffffffffffffff81116108fd57846108c7918801613164565b83528a8601519267ffffffffffffffff84116108fd576108ef61085b95879560019901613164565b8c82015293509150936107fc565b8f80fd5b8d80fd5b8e80fd5b6040513d8f823e3d90fd5b517fbf16aab6000000000000000000000000000000000000000000000000000000008c5273ffffffffffffffffffffffffffffffffffffffff1660045260248bfd5b506040517f01ffc9a70000000000000000000000000000000000000000000000000000000081527faff2afbf0000000000000000000000000000000000000000000000000000000060048201528881602481875afa908115610909578d916109c0575b501561071f565b90508881813d83116109f2575b6109d781836129d9565b810103126109ee576109e890612bf2565b386109b9565b8c80fd5b503d6109cd565b60048b7f5cf04449000000000000000000000000000000000000000000000000000000008152fd5b868567ffffffffffffffff888d8c878f868885928d88610a9961010073ffffffffffffffffffffffffffffffffffffffff600254169501516040519c8d977f01447eaa0000000000000000000000000000000000000000000000000000000089521660048801526060602488015260648701906131a6565b917ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc8684030160448701525191828152019190855b8a828210610f39575050505082809103915afa948515610f2e578395610e6e575b5015610d855750805b67ffffffffffffffff608087510191169052805b61010086015151811015610b425780610b2760019286613019565b516080610b39836101008b0151613019565b51015201610b0c565b5083610b4d8661346b565b91604051848101907f130ac867e79e2789f923760a88743d292acdf7002139a588206e2260f73f7321825267ffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016604082015267ffffffffffffffff8416606082015230608082015260808152610bcd60a0826129d9565b51902073ffffffffffffffffffffffffffffffffffffffff8585015116845167ffffffffffffffff6080816060840151169201511673ffffffffffffffffffffffffffffffffffffffff60a08801511660c088015191604051938a850195865260408501526060840152608083015260a082015260a08152610c5060c0826129d9565b51902060608501518681519101206040860151878151910120610100870151604051610cb681610c8a8c8201948d865260408301906131a6565b037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081018352826129d9565b51902091608088015189815191012093604051958a870197885260408701526060860152608085015260a084015260c083015260e082015260e08152610cfe610100826129d9565b51902082515267ffffffffffffffff60608351015116907f192442a2b2adb6a7948f097023cb6b57d29d3a7a5dd33e6666d33c39cc456f3267ffffffffffffffff60405192169180610d508682613291565b0390a37fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff600254166002555151604051908152f35b73ffffffffffffffffffffffffffffffffffffffff604051917fea458c0c00000000000000000000000000000000000000000000000000000000835267ffffffffffffffff8716600484015216602482015282816044818573ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165af1908115610e63578291610e2a575b50610af8565b90508281813d8311610e5c575b610e4181836129d9565b81010312610e5857610e529061327c565b86610e24565b5080fd5b503d610e37565b6040513d84823e3d90fd5b9094503d8084833e610e8081836129d9565b8101908481830312610f265780519067ffffffffffffffff8211610f2a570181601f82011215610f26578051610eb581612aeb565b92610ec360405194856129d9565b818452868085019260051b84010192818411610f2257878101925b848410610ef15750505050509387610aef565b835167ffffffffffffffff8111610f1e578991610f1385848094870101613164565b815201930192610ede565b8880fd5b8680fd5b8380fd5b8480fd5b6040513d85823e3d90fd5b8351805173ffffffffffffffffffffffffffffffffffffffff168652810151818601528a97508c965060409094019390920191600101610ace565b604082360312610fab57876040918251610f8d816129bd565b610f9685612b56565b815282850135838201528152019101906106ac565b8b80fd5b935093505096503d8089833e610fc581836129d9565b810196608082890312610f1e578151610fdf868401612bf2565b92604081015167ffffffffffffffff8111610fab578a611000918301613164565b99606082015167ffffffffffffffff81116109ee5761101f9201613164565b909298909338610666565b6040513d8a823e3d90fd5b602090611040613101565b82828b010152016104c9565b6024867f4e487b710000000000000000000000000000000000000000000000000000000081526011600452fd5b803b15610f2a578460405180927fe0a0e5060000000000000000000000000000000000000000000000000000000082528183816110bf6024356004018b60048401612c8e565b03925af180156110e2571561038957846110db919592956129d9565b9238610389565b6040513d87823e3d90fd5b6004847f1c0a3529000000000000000000000000000000000000000000000000000000008152fd5b73ffffffffffffffffffffffffffffffffffffffff8316600090815260028301602052604090205461034b5760248573ffffffffffffffffffffffffffffffffffffffff857fd0d2597600000000000000000000000000000000000000000000000000000000835216600452fd5b6004847fa4ec7479000000000000000000000000000000000000000000000000000000008152fd5b6004847f3ee5aeb5000000000000000000000000000000000000000000000000000000008152fd5b50346101525760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101525760043567ffffffffffffffff8111610e5857611223903690600401612b77565b73ffffffffffffffffffffffffffffffffffffffff600154163303611554575b919081907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8181360301915b84811015611550578060051b82013583811215610f2a57820191608083360312610f2a576040519461129f86612939565b6112a884612ad6565b86526112b660208501612b03565b9660208701978852604085013567ffffffffffffffff811161154c576112df9036908701612fb4565b9460408801958652606081013567ffffffffffffffff8111610f265761130791369101612fb4565b60608801908152875167ffffffffffffffff1683526006602052604080842099518a547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff169015159182901b68ff000000000000000016178a55909590815151611424575b5095976001019550815b855180518210156113b557906113ae73ffffffffffffffffffffffffffffffffffffffff6113a683600195613019565b5116896138c0565b5001611376565b505095909694506001929193519081516113d5575b50500193929361126e565b61141a67ffffffffffffffff7fc237ec1921f855ccd5e9a5af9733f2d58943a5a8501ec5988e305d7a4d42158692511692604051918291602083526020830190612ba8565b0390a238806113ca565b9893959296919094979860001461151557600184019591875b865180518210156114ba576114678273ffffffffffffffffffffffffffffffffffffffff92613019565b51168015611483579061147c6001928a61382f565b500161143d565b60248a67ffffffffffffffff8e51167f463258ff000000000000000000000000000000000000000000000000000000008252600452fd5b50509692955090929796937f330939f6eafe8bb516716892fe962ff19770570838686e6579dbc1cc51fc328161150b67ffffffffffffffff8a51169251604051918291602083526020830190612ba8565b0390a2388061136c565b60248767ffffffffffffffff8b51167f463258ff000000000000000000000000000000000000000000000000000000008252600452fd5b8280fd5b8380f35b73ffffffffffffffffffffffffffffffffffffffff60055416330315611243576004837f905d7d9b000000000000000000000000000000000000000000000000000000008152fd5b50346101525760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101525767ffffffffffffffff6115dd612aba565b16808252600660205260ff604083205460401c16908252600660205260016040832001916040518093849160208254918281520191845260208420935b81811061165157505061162f925003836129d9565b61164d60405192839215158352604060208401526040830190612ba8565b0390f35b845483526001948501948794506020909301920161161a565b50346101525760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101525767ffffffffffffffff6116ab612aba565b1681526006602052600167ffffffffffffffff604083205416019067ffffffffffffffff82116116ea5760208267ffffffffffffffff60405191168152f35b807f4e487b7100000000000000000000000000000000000000000000000000000000602492526011600452fd5b503461015257807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261015257602073ffffffffffffffffffffffffffffffffffffffff60015416604051908152f35b503461015257807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261015257805473ffffffffffffffffffffffffffffffffffffffff81163303611826577fffffffffffffffffffffffff000000000000000000000000000000000000000060015491338284161760015516825573ffffffffffffffffffffffffffffffffffffffff3391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08380a380f35b6004827f02b543c6000000000000000000000000000000000000000000000000000000008152fd5b503461015257807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261015257611885612f89565b5060a060405161189481612984565b60ff60025473ffffffffffffffffffffffffffffffffffffffff81168352831c161515602082015273ffffffffffffffffffffffffffffffffffffffff60035416604082015273ffffffffffffffffffffffffffffffffffffffff60045416606082015273ffffffffffffffffffffffffffffffffffffffff600554166080820152611969604051809273ffffffffffffffffffffffffffffffffffffffff60808092828151168552602081015115156020860152826040820151166040860152826060820151166060860152015116910152565bf35b50346101525760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261015257604060609167ffffffffffffffff6119b1612aba565b1681526006602052205473ffffffffffffffffffffffffffffffffffffffff6040519167ffffffffffffffff8116835260ff8160401c161515602084015260481c166040820152f35b50346101525760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101525773ffffffffffffffffffffffffffffffffffffffff611a47612b10565b611a4f6133f0565b168015611a82577fffffffffffffffffffffffff0000000000000000000000000000000000000000600754161760075580f35b6004827f8579befe000000000000000000000000000000000000000000000000000000008152fd5b50346101525760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101525760043567ffffffffffffffff8111610e5857611afa903690600401612b77565b9073ffffffffffffffffffffffffffffffffffffffff6004541690835b83811015611d635773ffffffffffffffffffffffffffffffffffffffff611b428260051b8401612f68565b1690604051917f70a08231000000000000000000000000000000000000000000000000000000008352306004840152602083602481845afa928315611d58578793611d25575b5082611b9a575b506001915001611b17565b8460405193611c3660208601957fa9059cbb00000000000000000000000000000000000000000000000000000000875283602482015282604482015260448152611be56064826129d9565b8a80604098895193611bf78b866129d9565b602085527f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65646020860152519082895af1611c2f61343b565b9086613a52565b805180611c72575b505060207f508d7d183612c18fc339b42618912b9fa3239f631dd7ec0671f950200a0fa66e9160019651908152a338611b8f565b819294959693509060209181010312610f1e576020611c919101612bf2565b15611ca25792919085903880611c3e565b608490517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f742073756363656564000000000000000000000000000000000000000000006064820152fd5b9092506020813d8211611d50575b81611d40602093836129d9565b81010312610f2257519138611b88565b3d9150611d33565b6040513d89823e3d90fd5b8480f35b50346101525760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261015257611d9f612aba565b506024359073ffffffffffffffffffffffffffffffffffffffff82168203610152576020611dcc83612e65565b73ffffffffffffffffffffffffffffffffffffffff60405191168152f35b50346101525760a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261015257604051611e2681612984565b611e2e612b10565b8152602435801515810361154c576020820190815260443573ffffffffffffffffffffffffffffffffffffffff81168103610f265760408301908152611e72612b33565b90606084019182526084359273ffffffffffffffffffffffffffffffffffffffff841684036121f05760808501938452611eaa6133f0565b73ffffffffffffffffffffffffffffffffffffffff8551161580156121d1575b80156121c7575b61219f579273ffffffffffffffffffffffffffffffffffffffff859381809461012097827fc7372d2d886367d7bb1b0e0708a5436f2c91d6963de210eb2dc1ec2ecd6d21f19a51167fffffffffffffffffffffffff000000000000000000000000000000000000000060025416176002555115157fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff74ff00000000000000000000000000000000000000006002549260a01b1691161760025551167fffffffffffffffffffffffff0000000000000000000000000000000000000000600354161760035551167fffffffffffffffffffffffff0000000000000000000000000000000000000000600454161760045551167fffffffffffffffffffffffff0000000000000000000000000000000000000000600554161760055561219b6040519161201b83612939565b67ffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016835273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016602084015273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016604084015273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016606084015261214b604051809473ffffffffffffffffffffffffffffffffffffffff6060809267ffffffffffffffff8151168552826020820151166020860152826040820151166040860152015116910152565b608083019073ffffffffffffffffffffffffffffffffffffffff60808092828151168552602081015115156020860152826040820151166040860152826060820151166060860152015116910152565ba180f35b6004867f35be3ac8000000000000000000000000000000000000000000000000000000008152fd5b5080511515611ed1565b5073ffffffffffffffffffffffffffffffffffffffff83511615611eca565b8580fd5b50346101525760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610152576004359067ffffffffffffffff8211610152573660238301121561015257816004013561225081612aeb565b9261225e60405194856129d9565b8184526024606060208601930282010190368211610f2657602401915b8183106123fc5750505061228d6133f0565b805b82518110156123f8576122a28184613019565b5167ffffffffffffffff6122b68386613019565b5151169081156123cc57907fd5ad72bc37dc7a80a8b9b9df20500046fd7341adb1be2258a540466fdd7dcef5606060019493838752600660205260ff6040882061238d604060208501519483547fffffff0000000000000000000000000000000000000000ffffffffffffffffff7cffffffffffffffffffffffffffffffffffffffff0000000000000000008860481b1691161784550151151582907fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff68ff0000000000000000835492151560401b169116179055565b5473ffffffffffffffffffffffffffffffffffffffff6040519367ffffffffffffffff8316855216602084015260401c1615156040820152a20161228f565b602484837fc35aa79d000000000000000000000000000000000000000000000000000000008252600452fd5b5080f35b606083360312610f26576040516060810181811067ffffffffffffffff8211176124775760405261242c84612ad6565b8152602084013573ffffffffffffffffffffffffffffffffffffffff811681036121f057918160609360208094015261246760408701612b03565b604082015281520192019161227b565b6024867f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b50346101525760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610152576124dc612aba565b60243567ffffffffffffffff811161154c5760a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc823603011261154c576040517f2cbc26bb00000000000000000000000000000000000000000000000000000000815277ffffffffffffffff000000000000000000000000000000008360801b16600482015260208160248173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165afa9081156126d2578491612698575b506126625761260f9160209173ffffffffffffffffffffffffffffffffffffffff60025416906040518095819482937fd8694ccd0000000000000000000000000000000000000000000000000000000084526004019060048401612c8e565b03915afa908115610e6357829161262c575b602082604051908152f35b90506020813d60201161265a575b81612647602093836129d9565b81010312610e5857602091505138612621565b3d915061263a565b60248367ffffffffffffffff847ffdbd6a7200000000000000000000000000000000000000000000000000000000835216600452fd5b90506020813d6020116126ca575b816126b3602093836129d9565b81010312610f26576126c490612bf2565b386125b0565b3d91506126a6565b6040513d86823e3d90fd5b503461015257807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610152575061164d60405161271e6040826129d9565b600c81527f4f6e52616d7020312e362e3000000000000000000000000000000000000000006020820152604051918291602083526020830190612a77565b503461015257807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261015257602073ffffffffffffffffffffffffffffffffffffffff60075416604051908152f35b905034610e5857817ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610e5857806127ea606092612939565b8281528260208201528260408201520152608060405161280981612939565b67ffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016815273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016602082015273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016604082015273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000166060820152611969604051809273ffffffffffffffffffffffffffffffffffffffff6060809267ffffffffffffffff8151168552826020820151166020860152826040820151166040860152015116910152565b6080810190811067ffffffffffffffff82111761295557604052565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60a0810190811067ffffffffffffffff82111761295557604052565b610120810190811067ffffffffffffffff82111761295557604052565b6040810190811067ffffffffffffffff82111761295557604052565b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff82111761295557604052565b67ffffffffffffffff811161295557601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b60005b838110612a675750506000910152565b8181015183820152602001612a57565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f602093612ab381518092818752878088019101612a54565b0116010190565b6004359067ffffffffffffffff82168203612ad157565b600080fd5b359067ffffffffffffffff82168203612ad157565b67ffffffffffffffff81116129555760051b60200190565b35908115158203612ad157565b6004359073ffffffffffffffffffffffffffffffffffffffff82168203612ad157565b6064359073ffffffffffffffffffffffffffffffffffffffff82168203612ad157565b359073ffffffffffffffffffffffffffffffffffffffff82168203612ad157565b9181601f84011215612ad15782359167ffffffffffffffff8311612ad1576020808501948460051b010111612ad157565b906020808351928381520192019060005b818110612bc65750505090565b825173ffffffffffffffffffffffffffffffffffffffff16845260209384019390920191600101612bb9565b51908115158203612ad157565b90357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe182360301811215612ad157016020813591019167ffffffffffffffff8211612ad1578136038313612ad157565b601f82602094937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0938186528686013760008582860101520116010190565b9067ffffffffffffffff9093929316815260406020820152612d04612cc7612cb68580612bff565b60a0604086015260e0850191612c4f565b612cd46020860186612bff565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0858403016060860152612c4f565b9060408401357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe185360301811215612ad15784016020813591019267ffffffffffffffff8211612ad1578160061b36038413612ad1578281037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0016080840152818152602001929060005b818110612e0557505050612dd28473ffffffffffffffffffffffffffffffffffffffff612dc26060612e02979801612b56565b1660a08401526080810190612bff565b9160c07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc082860301910152612c4f565b90565b90919360408060019273ffffffffffffffffffffffffffffffffffffffff612e2c89612b56565b16815260208881013590820152019501929101612d8f565b519073ffffffffffffffffffffffffffffffffffffffff82168203612ad157565b73ffffffffffffffffffffffffffffffffffffffff604051917fbbe4f6db00000000000000000000000000000000000000000000000000000000835216600482015260208160248173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165afa8015612f5c57600090612f0f575b73ffffffffffffffffffffffffffffffffffffffff91501690565b506020813d602011612f54575b81612f29602093836129d9565b81010312612ad157612f4f73ffffffffffffffffffffffffffffffffffffffff91612e44565b612ef4565b3d9150612f1c565b6040513d6000823e3d90fd5b3573ffffffffffffffffffffffffffffffffffffffff81168103612ad15790565b60405190612f9682612984565b60006080838281528260208201528260408201528260608201520152565b9080601f83011215612ad1578135612fcb81612aeb565b92612fd960405194856129d9565b81845260208085019260051b820101928311612ad157602001905b8282106130015750505090565b6020809161300e84612b56565b815201910190612ff4565b805182101561302d5760209160051b010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe181360301821215612ad1570180359067ffffffffffffffff8211612ad157602001918136038313612ad157565b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe181360301821215612ad1570180359067ffffffffffffffff8211612ad157602001918160061b36038313612ad157565b6040519061310e82612984565b6060608083600081528260208201528260408201526000838201520152565b92919261313982612a1a565b9161314760405193846129d9565b829481845281830111612ad1578281602093846000960137010152565b81601f82011215612ad157805161317a81612a1a565b9261318860405194856129d9565b81845260208284010111612ad157612e029160208085019101612a54565b9080602083519182815201916020808360051b8301019401926000915b8383106131d257505050505090565b909192939460208061326d837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0866001960301875289519073ffffffffffffffffffffffffffffffffffffffff825116815260806132526132408685015160a08886015260a0850190612a77565b60408501518482036040860152612a77565b92606081015160608401520151906080818403910152612a77565b970193019301919392906131c3565b519067ffffffffffffffff82168203612ad157565b90612e02916020815267ffffffffffffffff6080835180516020850152826020820151166040850152826040820151166060850152826060820151168285015201511660a082015273ffffffffffffffffffffffffffffffffffffffff60208301511660c082015261010061338561335061331d60408601516101a060e08701526101c0860190612a77565b60608601517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08683030185870152612a77565b60808501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe085830301610120860152612a77565b9273ffffffffffffffffffffffffffffffffffffffff60a08201511661014084015260c081015161016084015260e08101516101808401520151906101a07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0828503019101526131a6565b73ffffffffffffffffffffffffffffffffffffffff60015416330361341157565b7f2b5c74de0000000000000000000000000000000000000000000000000000000060005260046000fd5b3d15613466573d9061344c82612a1a565b9161345a60405193846129d9565b82523d6000602084013e565b606090565b600061350c819260405161347e816129a0565b613486612f89565b815283602082015260606040820152606080820152606060808201528360a08201528360c08201528360e082015260606101008201525073ffffffffffffffffffffffffffffffffffffffff60075416906040519485809481937f8a06fadb00000000000000000000000000000000000000000000000000000000835260048301613291565b03925af18091600091613569575b5090612e025761356561352b61343b565b6040519182917f828ebdfb000000000000000000000000000000000000000000000000000000008352602060048401526024830190612a77565b0390fd5b3d8083833e61357881836129d9565b81019060208183031261154c5780519067ffffffffffffffff8211610f26570191828203926101a08412610e585760a0604051946135b5866129a0565b12610e58576040516135c681612984565b815181526135d66020830161327c565b60208201526135e76040830161327c565b60408201526135f86060830161327c565b60608201526136096080830161327c565b6080820152845261361c60a08201612e44565b602085015260c081015167ffffffffffffffff811161154c5783613641918301613164565b604085015260e081015167ffffffffffffffff811161154c5783613666918301613164565b606085015261010081015167ffffffffffffffff811161154c578361368c918301613164565b608085015261369e6101208201612e44565b60a085015261014081015160c085015261016081015160e08501526101808101519067ffffffffffffffff821161154c570182601f82011215610e58578051916136e783612aeb565b936136f560405195866129d9565b83855260208086019460051b8401019281841161154c5760208101945b84861061372b575050505050506101008201523861351a565b855167ffffffffffffffff8111610f2a57820160a07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08286030112610f2a576040519061377782612984565b61378360208201612e44565b8252604081015167ffffffffffffffff8111610f22578560206137a892840101613164565b6020830152606081015167ffffffffffffffff8111610f22578560206137d092840101613164565b60408301526080810151606083015260a08101519067ffffffffffffffff8211610f22579161380786602080969481960101613164565b6080820152815201950194613712565b805482101561302d5760005260206000200190600090565b60008281526001820160205260409020546138b9578054906801000000000000000082101561295557826138a261386d846001809601855584613817565b81939154907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9060031b92831b921b19161790565b905580549260005201602052604060002055600190565b5050600090565b9060018201918160005282602052604060002054801515600014613a49577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8101818111613a1a578254907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8201918211613a1a578181036139e3575b505050805480156139b4577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01906139758282613817565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82549160031b1b191690555560005260205260006040812055600190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b613a036139f361386d9386613817565b90549060031b1c92839286613817565b90556000528360205260406000205538808061393d565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b50505050600090565b91929015613acd5750815115613a66575090565b3b15613a6f5790565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152fd5b825190915015613ae05750805190602001fd5b613565906040519182917f08c379a0000000000000000000000000000000000000000000000000000000008352602060048401526024830190612a7756fea164736f6c634300081a000a",
}

var OnRampWithMessageTransformerABI = OnRampWithMessageTransformerMetaData.ABI

var OnRampWithMessageTransformerBin = OnRampWithMessageTransformerMetaData.Bin

func DeployOnRampWithMessageTransformer(auth *bind.TransactOpts, backend bind.ContractBackend, staticConfig OnRampStaticConfig, dynamicConfig OnRampDynamicConfig, destChainConfigs []OnRampDestChainConfigArgs, messageTransformerAddr common.Address) (common.Address, *types.Transaction, *OnRampWithMessageTransformer, error) {
	parsed, err := OnRampWithMessageTransformerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OnRampWithMessageTransformerBin), backend, staticConfig, dynamicConfig, destChainConfigs, messageTransformerAddr)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OnRampWithMessageTransformer{address: address, abi: *parsed, OnRampWithMessageTransformerCaller: OnRampWithMessageTransformerCaller{contract: contract}, OnRampWithMessageTransformerTransactor: OnRampWithMessageTransformerTransactor{contract: contract}, OnRampWithMessageTransformerFilterer: OnRampWithMessageTransformerFilterer{contract: contract}}, nil
}

type OnRampWithMessageTransformer struct {
	address common.Address
	abi     abi.ABI
	OnRampWithMessageTransformerCaller
	OnRampWithMessageTransformerTransactor
	OnRampWithMessageTransformerFilterer
}

type OnRampWithMessageTransformerCaller struct {
	contract *bind.BoundContract
}

type OnRampWithMessageTransformerTransactor struct {
	contract *bind.BoundContract
}

type OnRampWithMessageTransformerFilterer struct {
	contract *bind.BoundContract
}

type OnRampWithMessageTransformerSession struct {
	Contract     *OnRampWithMessageTransformer
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type OnRampWithMessageTransformerCallerSession struct {
	Contract *OnRampWithMessageTransformerCaller
	CallOpts bind.CallOpts
}

type OnRampWithMessageTransformerTransactorSession struct {
	Contract     *OnRampWithMessageTransformerTransactor
	TransactOpts bind.TransactOpts
}

type OnRampWithMessageTransformerRaw struct {
	Contract *OnRampWithMessageTransformer
}

type OnRampWithMessageTransformerCallerRaw struct {
	Contract *OnRampWithMessageTransformerCaller
}

type OnRampWithMessageTransformerTransactorRaw struct {
	Contract *OnRampWithMessageTransformerTransactor
}

func NewOnRampWithMessageTransformer(address common.Address, backend bind.ContractBackend) (*OnRampWithMessageTransformer, error) {
	abi, err := abi.JSON(strings.NewReader(OnRampWithMessageTransformerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindOnRampWithMessageTransformer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OnRampWithMessageTransformer{address: address, abi: abi, OnRampWithMessageTransformerCaller: OnRampWithMessageTransformerCaller{contract: contract}, OnRampWithMessageTransformerTransactor: OnRampWithMessageTransformerTransactor{contract: contract}, OnRampWithMessageTransformerFilterer: OnRampWithMessageTransformerFilterer{contract: contract}}, nil
}

func NewOnRampWithMessageTransformerCaller(address common.Address, caller bind.ContractCaller) (*OnRampWithMessageTransformerCaller, error) {
	contract, err := bindOnRampWithMessageTransformer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OnRampWithMessageTransformerCaller{contract: contract}, nil
}

func NewOnRampWithMessageTransformerTransactor(address common.Address, transactor bind.ContractTransactor) (*OnRampWithMessageTransformerTransactor, error) {
	contract, err := bindOnRampWithMessageTransformer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OnRampWithMessageTransformerTransactor{contract: contract}, nil
}

func NewOnRampWithMessageTransformerFilterer(address common.Address, filterer bind.ContractFilterer) (*OnRampWithMessageTransformerFilterer, error) {
	contract, err := bindOnRampWithMessageTransformer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OnRampWithMessageTransformerFilterer{contract: contract}, nil
}

func bindOnRampWithMessageTransformer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OnRampWithMessageTransformerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OnRampWithMessageTransformer.Contract.OnRampWithMessageTransformerCaller.contract.Call(opts, result, method, params...)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.Contract.OnRampWithMessageTransformerTransactor.contract.Transfer(opts)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.Contract.OnRampWithMessageTransformerTransactor.contract.Transact(opts, method, params...)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OnRampWithMessageTransformer.Contract.contract.Call(opts, result, method, params...)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.Contract.contract.Transfer(opts)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.Contract.contract.Transact(opts, method, params...)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCaller) GetAllowedSendersList(opts *bind.CallOpts, destChainSelector uint64) (GetAllowedSendersList,

	error) {
	var out []interface{}
	err := _OnRampWithMessageTransformer.contract.Call(opts, &out, "getAllowedSendersList", destChainSelector)

	outstruct := new(GetAllowedSendersList)
	if err != nil {
		return *outstruct, err
	}

	outstruct.IsEnabled = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfiguredAddresses = *abi.ConvertType(out[1], new([]common.Address)).(*[]common.Address)

	return *outstruct, err

}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerSession) GetAllowedSendersList(destChainSelector uint64) (GetAllowedSendersList,

	error) {
	return _OnRampWithMessageTransformer.Contract.GetAllowedSendersList(&_OnRampWithMessageTransformer.CallOpts, destChainSelector)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCallerSession) GetAllowedSendersList(destChainSelector uint64) (GetAllowedSendersList,

	error) {
	return _OnRampWithMessageTransformer.Contract.GetAllowedSendersList(&_OnRampWithMessageTransformer.CallOpts, destChainSelector)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCaller) GetDestChainConfig(opts *bind.CallOpts, destChainSelector uint64) (GetDestChainConfig,

	error) {
	var out []interface{}
	err := _OnRampWithMessageTransformer.contract.Call(opts, &out, "getDestChainConfig", destChainSelector)

	outstruct := new(GetDestChainConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.SequenceNumber = *abi.ConvertType(out[0], new(uint64)).(*uint64)
	outstruct.AllowlistEnabled = *abi.ConvertType(out[1], new(bool)).(*bool)
	outstruct.Router = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)

	return *outstruct, err

}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerSession) GetDestChainConfig(destChainSelector uint64) (GetDestChainConfig,

	error) {
	return _OnRampWithMessageTransformer.Contract.GetDestChainConfig(&_OnRampWithMessageTransformer.CallOpts, destChainSelector)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCallerSession) GetDestChainConfig(destChainSelector uint64) (GetDestChainConfig,

	error) {
	return _OnRampWithMessageTransformer.Contract.GetDestChainConfig(&_OnRampWithMessageTransformer.CallOpts, destChainSelector)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCaller) GetDynamicConfig(opts *bind.CallOpts) (OnRampDynamicConfig, error) {
	var out []interface{}
	err := _OnRampWithMessageTransformer.contract.Call(opts, &out, "getDynamicConfig")

	if err != nil {
		return *new(OnRampDynamicConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(OnRampDynamicConfig)).(*OnRampDynamicConfig)

	return out0, err

}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerSession) GetDynamicConfig() (OnRampDynamicConfig, error) {
	return _OnRampWithMessageTransformer.Contract.GetDynamicConfig(&_OnRampWithMessageTransformer.CallOpts)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCallerSession) GetDynamicConfig() (OnRampDynamicConfig, error) {
	return _OnRampWithMessageTransformer.Contract.GetDynamicConfig(&_OnRampWithMessageTransformer.CallOpts)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCaller) GetExpectedNextSequenceNumber(opts *bind.CallOpts, destChainSelector uint64) (uint64, error) {
	var out []interface{}
	err := _OnRampWithMessageTransformer.contract.Call(opts, &out, "getExpectedNextSequenceNumber", destChainSelector)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerSession) GetExpectedNextSequenceNumber(destChainSelector uint64) (uint64, error) {
	return _OnRampWithMessageTransformer.Contract.GetExpectedNextSequenceNumber(&_OnRampWithMessageTransformer.CallOpts, destChainSelector)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCallerSession) GetExpectedNextSequenceNumber(destChainSelector uint64) (uint64, error) {
	return _OnRampWithMessageTransformer.Contract.GetExpectedNextSequenceNumber(&_OnRampWithMessageTransformer.CallOpts, destChainSelector)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCaller) GetFee(opts *bind.CallOpts, destChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error) {
	var out []interface{}
	err := _OnRampWithMessageTransformer.contract.Call(opts, &out, "getFee", destChainSelector, message)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerSession) GetFee(destChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error) {
	return _OnRampWithMessageTransformer.Contract.GetFee(&_OnRampWithMessageTransformer.CallOpts, destChainSelector, message)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCallerSession) GetFee(destChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error) {
	return _OnRampWithMessageTransformer.Contract.GetFee(&_OnRampWithMessageTransformer.CallOpts, destChainSelector, message)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCaller) GetMessageTransformer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OnRampWithMessageTransformer.contract.Call(opts, &out, "getMessageTransformer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerSession) GetMessageTransformer() (common.Address, error) {
	return _OnRampWithMessageTransformer.Contract.GetMessageTransformer(&_OnRampWithMessageTransformer.CallOpts)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCallerSession) GetMessageTransformer() (common.Address, error) {
	return _OnRampWithMessageTransformer.Contract.GetMessageTransformer(&_OnRampWithMessageTransformer.CallOpts)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCaller) GetPoolBySourceToken(opts *bind.CallOpts, arg0 uint64, sourceToken common.Address) (common.Address, error) {
	var out []interface{}
	err := _OnRampWithMessageTransformer.contract.Call(opts, &out, "getPoolBySourceToken", arg0, sourceToken)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerSession) GetPoolBySourceToken(arg0 uint64, sourceToken common.Address) (common.Address, error) {
	return _OnRampWithMessageTransformer.Contract.GetPoolBySourceToken(&_OnRampWithMessageTransformer.CallOpts, arg0, sourceToken)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCallerSession) GetPoolBySourceToken(arg0 uint64, sourceToken common.Address) (common.Address, error) {
	return _OnRampWithMessageTransformer.Contract.GetPoolBySourceToken(&_OnRampWithMessageTransformer.CallOpts, arg0, sourceToken)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCaller) GetStaticConfig(opts *bind.CallOpts) (OnRampStaticConfig, error) {
	var out []interface{}
	err := _OnRampWithMessageTransformer.contract.Call(opts, &out, "getStaticConfig")

	if err != nil {
		return *new(OnRampStaticConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(OnRampStaticConfig)).(*OnRampStaticConfig)

	return out0, err

}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerSession) GetStaticConfig() (OnRampStaticConfig, error) {
	return _OnRampWithMessageTransformer.Contract.GetStaticConfig(&_OnRampWithMessageTransformer.CallOpts)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCallerSession) GetStaticConfig() (OnRampStaticConfig, error) {
	return _OnRampWithMessageTransformer.Contract.GetStaticConfig(&_OnRampWithMessageTransformer.CallOpts)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCaller) GetSupportedTokens(opts *bind.CallOpts, arg0 uint64) ([]common.Address, error) {
	var out []interface{}
	err := _OnRampWithMessageTransformer.contract.Call(opts, &out, "getSupportedTokens", arg0)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerSession) GetSupportedTokens(arg0 uint64) ([]common.Address, error) {
	return _OnRampWithMessageTransformer.Contract.GetSupportedTokens(&_OnRampWithMessageTransformer.CallOpts, arg0)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCallerSession) GetSupportedTokens(arg0 uint64) ([]common.Address, error) {
	return _OnRampWithMessageTransformer.Contract.GetSupportedTokens(&_OnRampWithMessageTransformer.CallOpts, arg0)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OnRampWithMessageTransformer.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerSession) Owner() (common.Address, error) {
	return _OnRampWithMessageTransformer.Contract.Owner(&_OnRampWithMessageTransformer.CallOpts)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCallerSession) Owner() (common.Address, error) {
	return _OnRampWithMessageTransformer.Contract.Owner(&_OnRampWithMessageTransformer.CallOpts)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _OnRampWithMessageTransformer.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerSession) TypeAndVersion() (string, error) {
	return _OnRampWithMessageTransformer.Contract.TypeAndVersion(&_OnRampWithMessageTransformer.CallOpts)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCallerSession) TypeAndVersion() (string, error) {
	return _OnRampWithMessageTransformer.Contract.TypeAndVersion(&_OnRampWithMessageTransformer.CallOpts)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.contract.Transact(opts, "acceptOwnership")
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerSession) AcceptOwnership() (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.Contract.AcceptOwnership(&_OnRampWithMessageTransformer.TransactOpts)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.Contract.AcceptOwnership(&_OnRampWithMessageTransformer.TransactOpts)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerTransactor) ApplyAllowlistUpdates(opts *bind.TransactOpts, allowlistConfigArgsItems []OnRampAllowlistConfigArgs) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.contract.Transact(opts, "applyAllowlistUpdates", allowlistConfigArgsItems)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerSession) ApplyAllowlistUpdates(allowlistConfigArgsItems []OnRampAllowlistConfigArgs) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.Contract.ApplyAllowlistUpdates(&_OnRampWithMessageTransformer.TransactOpts, allowlistConfigArgsItems)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerTransactorSession) ApplyAllowlistUpdates(allowlistConfigArgsItems []OnRampAllowlistConfigArgs) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.Contract.ApplyAllowlistUpdates(&_OnRampWithMessageTransformer.TransactOpts, allowlistConfigArgsItems)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerTransactor) ApplyDestChainConfigUpdates(opts *bind.TransactOpts, destChainConfigArgs []OnRampDestChainConfigArgs) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.contract.Transact(opts, "applyDestChainConfigUpdates", destChainConfigArgs)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerSession) ApplyDestChainConfigUpdates(destChainConfigArgs []OnRampDestChainConfigArgs) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.Contract.ApplyDestChainConfigUpdates(&_OnRampWithMessageTransformer.TransactOpts, destChainConfigArgs)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerTransactorSession) ApplyDestChainConfigUpdates(destChainConfigArgs []OnRampDestChainConfigArgs) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.Contract.ApplyDestChainConfigUpdates(&_OnRampWithMessageTransformer.TransactOpts, destChainConfigArgs)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerTransactor) ForwardFromRouter(opts *bind.TransactOpts, destChainSelector uint64, message ClientEVM2AnyMessage, feeTokenAmount *big.Int, originalSender common.Address) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.contract.Transact(opts, "forwardFromRouter", destChainSelector, message, feeTokenAmount, originalSender)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerSession) ForwardFromRouter(destChainSelector uint64, message ClientEVM2AnyMessage, feeTokenAmount *big.Int, originalSender common.Address) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.Contract.ForwardFromRouter(&_OnRampWithMessageTransformer.TransactOpts, destChainSelector, message, feeTokenAmount, originalSender)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerTransactorSession) ForwardFromRouter(destChainSelector uint64, message ClientEVM2AnyMessage, feeTokenAmount *big.Int, originalSender common.Address) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.Contract.ForwardFromRouter(&_OnRampWithMessageTransformer.TransactOpts, destChainSelector, message, feeTokenAmount, originalSender)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerTransactor) SetDynamicConfig(opts *bind.TransactOpts, dynamicConfig OnRampDynamicConfig) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.contract.Transact(opts, "setDynamicConfig", dynamicConfig)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerSession) SetDynamicConfig(dynamicConfig OnRampDynamicConfig) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.Contract.SetDynamicConfig(&_OnRampWithMessageTransformer.TransactOpts, dynamicConfig)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerTransactorSession) SetDynamicConfig(dynamicConfig OnRampDynamicConfig) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.Contract.SetDynamicConfig(&_OnRampWithMessageTransformer.TransactOpts, dynamicConfig)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerTransactor) SetMessageTransformer(opts *bind.TransactOpts, messageTransformerAddr common.Address) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.contract.Transact(opts, "setMessageTransformer", messageTransformerAddr)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerSession) SetMessageTransformer(messageTransformerAddr common.Address) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.Contract.SetMessageTransformer(&_OnRampWithMessageTransformer.TransactOpts, messageTransformerAddr)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerTransactorSession) SetMessageTransformer(messageTransformerAddr common.Address) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.Contract.SetMessageTransformer(&_OnRampWithMessageTransformer.TransactOpts, messageTransformerAddr)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.contract.Transact(opts, "transferOwnership", to)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.Contract.TransferOwnership(&_OnRampWithMessageTransformer.TransactOpts, to)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.Contract.TransferOwnership(&_OnRampWithMessageTransformer.TransactOpts, to)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerTransactor) WithdrawFeeTokens(opts *bind.TransactOpts, feeTokens []common.Address) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.contract.Transact(opts, "withdrawFeeTokens", feeTokens)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerSession) WithdrawFeeTokens(feeTokens []common.Address) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.Contract.WithdrawFeeTokens(&_OnRampWithMessageTransformer.TransactOpts, feeTokens)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerTransactorSession) WithdrawFeeTokens(feeTokens []common.Address) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.Contract.WithdrawFeeTokens(&_OnRampWithMessageTransformer.TransactOpts, feeTokens)
}

type OnRampWithMessageTransformerAllowListAdminSetIterator struct {
	Event *OnRampWithMessageTransformerAllowListAdminSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OnRampWithMessageTransformerAllowListAdminSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnRampWithMessageTransformerAllowListAdminSet)
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
		it.Event = new(OnRampWithMessageTransformerAllowListAdminSet)
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

func (it *OnRampWithMessageTransformerAllowListAdminSetIterator) Error() error {
	return it.fail
}

func (it *OnRampWithMessageTransformerAllowListAdminSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OnRampWithMessageTransformerAllowListAdminSet struct {
	AllowlistAdmin common.Address
	Raw            types.Log
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) FilterAllowListAdminSet(opts *bind.FilterOpts, allowlistAdmin []common.Address) (*OnRampWithMessageTransformerAllowListAdminSetIterator, error) {

	var allowlistAdminRule []interface{}
	for _, allowlistAdminItem := range allowlistAdmin {
		allowlistAdminRule = append(allowlistAdminRule, allowlistAdminItem)
	}

	logs, sub, err := _OnRampWithMessageTransformer.contract.FilterLogs(opts, "AllowListAdminSet", allowlistAdminRule)
	if err != nil {
		return nil, err
	}
	return &OnRampWithMessageTransformerAllowListAdminSetIterator{contract: _OnRampWithMessageTransformer.contract, event: "AllowListAdminSet", logs: logs, sub: sub}, nil
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) WatchAllowListAdminSet(opts *bind.WatchOpts, sink chan<- *OnRampWithMessageTransformerAllowListAdminSet, allowlistAdmin []common.Address) (event.Subscription, error) {

	var allowlistAdminRule []interface{}
	for _, allowlistAdminItem := range allowlistAdmin {
		allowlistAdminRule = append(allowlistAdminRule, allowlistAdminItem)
	}

	logs, sub, err := _OnRampWithMessageTransformer.contract.WatchLogs(opts, "AllowListAdminSet", allowlistAdminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OnRampWithMessageTransformerAllowListAdminSet)
				if err := _OnRampWithMessageTransformer.contract.UnpackLog(event, "AllowListAdminSet", log); err != nil {
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

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) ParseAllowListAdminSet(log types.Log) (*OnRampWithMessageTransformerAllowListAdminSet, error) {
	event := new(OnRampWithMessageTransformerAllowListAdminSet)
	if err := _OnRampWithMessageTransformer.contract.UnpackLog(event, "AllowListAdminSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OnRampWithMessageTransformerAllowListSendersAddedIterator struct {
	Event *OnRampWithMessageTransformerAllowListSendersAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OnRampWithMessageTransformerAllowListSendersAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnRampWithMessageTransformerAllowListSendersAdded)
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
		it.Event = new(OnRampWithMessageTransformerAllowListSendersAdded)
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

func (it *OnRampWithMessageTransformerAllowListSendersAddedIterator) Error() error {
	return it.fail
}

func (it *OnRampWithMessageTransformerAllowListSendersAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OnRampWithMessageTransformerAllowListSendersAdded struct {
	DestChainSelector uint64
	Senders           []common.Address
	Raw               types.Log
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) FilterAllowListSendersAdded(opts *bind.FilterOpts, destChainSelector []uint64) (*OnRampWithMessageTransformerAllowListSendersAddedIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _OnRampWithMessageTransformer.contract.FilterLogs(opts, "AllowListSendersAdded", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &OnRampWithMessageTransformerAllowListSendersAddedIterator{contract: _OnRampWithMessageTransformer.contract, event: "AllowListSendersAdded", logs: logs, sub: sub}, nil
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) WatchAllowListSendersAdded(opts *bind.WatchOpts, sink chan<- *OnRampWithMessageTransformerAllowListSendersAdded, destChainSelector []uint64) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _OnRampWithMessageTransformer.contract.WatchLogs(opts, "AllowListSendersAdded", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OnRampWithMessageTransformerAllowListSendersAdded)
				if err := _OnRampWithMessageTransformer.contract.UnpackLog(event, "AllowListSendersAdded", log); err != nil {
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

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) ParseAllowListSendersAdded(log types.Log) (*OnRampWithMessageTransformerAllowListSendersAdded, error) {
	event := new(OnRampWithMessageTransformerAllowListSendersAdded)
	if err := _OnRampWithMessageTransformer.contract.UnpackLog(event, "AllowListSendersAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OnRampWithMessageTransformerAllowListSendersRemovedIterator struct {
	Event *OnRampWithMessageTransformerAllowListSendersRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OnRampWithMessageTransformerAllowListSendersRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnRampWithMessageTransformerAllowListSendersRemoved)
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
		it.Event = new(OnRampWithMessageTransformerAllowListSendersRemoved)
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

func (it *OnRampWithMessageTransformerAllowListSendersRemovedIterator) Error() error {
	return it.fail
}

func (it *OnRampWithMessageTransformerAllowListSendersRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OnRampWithMessageTransformerAllowListSendersRemoved struct {
	DestChainSelector uint64
	Senders           []common.Address
	Raw               types.Log
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) FilterAllowListSendersRemoved(opts *bind.FilterOpts, destChainSelector []uint64) (*OnRampWithMessageTransformerAllowListSendersRemovedIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _OnRampWithMessageTransformer.contract.FilterLogs(opts, "AllowListSendersRemoved", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &OnRampWithMessageTransformerAllowListSendersRemovedIterator{contract: _OnRampWithMessageTransformer.contract, event: "AllowListSendersRemoved", logs: logs, sub: sub}, nil
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) WatchAllowListSendersRemoved(opts *bind.WatchOpts, sink chan<- *OnRampWithMessageTransformerAllowListSendersRemoved, destChainSelector []uint64) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _OnRampWithMessageTransformer.contract.WatchLogs(opts, "AllowListSendersRemoved", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OnRampWithMessageTransformerAllowListSendersRemoved)
				if err := _OnRampWithMessageTransformer.contract.UnpackLog(event, "AllowListSendersRemoved", log); err != nil {
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

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) ParseAllowListSendersRemoved(log types.Log) (*OnRampWithMessageTransformerAllowListSendersRemoved, error) {
	event := new(OnRampWithMessageTransformerAllowListSendersRemoved)
	if err := _OnRampWithMessageTransformer.contract.UnpackLog(event, "AllowListSendersRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OnRampWithMessageTransformerCCIPMessageSentIterator struct {
	Event *OnRampWithMessageTransformerCCIPMessageSent

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OnRampWithMessageTransformerCCIPMessageSentIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnRampWithMessageTransformerCCIPMessageSent)
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
		it.Event = new(OnRampWithMessageTransformerCCIPMessageSent)
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

func (it *OnRampWithMessageTransformerCCIPMessageSentIterator) Error() error {
	return it.fail
}

func (it *OnRampWithMessageTransformerCCIPMessageSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OnRampWithMessageTransformerCCIPMessageSent struct {
	DestChainSelector uint64
	SequenceNumber    uint64
	Message           InternalEVM2AnyRampMessage
	Raw               types.Log
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) FilterCCIPMessageSent(opts *bind.FilterOpts, destChainSelector []uint64, sequenceNumber []uint64) (*OnRampWithMessageTransformerCCIPMessageSentIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}
	var sequenceNumberRule []interface{}
	for _, sequenceNumberItem := range sequenceNumber {
		sequenceNumberRule = append(sequenceNumberRule, sequenceNumberItem)
	}

	logs, sub, err := _OnRampWithMessageTransformer.contract.FilterLogs(opts, "CCIPMessageSent", destChainSelectorRule, sequenceNumberRule)
	if err != nil {
		return nil, err
	}
	return &OnRampWithMessageTransformerCCIPMessageSentIterator{contract: _OnRampWithMessageTransformer.contract, event: "CCIPMessageSent", logs: logs, sub: sub}, nil
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) WatchCCIPMessageSent(opts *bind.WatchOpts, sink chan<- *OnRampWithMessageTransformerCCIPMessageSent, destChainSelector []uint64, sequenceNumber []uint64) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}
	var sequenceNumberRule []interface{}
	for _, sequenceNumberItem := range sequenceNumber {
		sequenceNumberRule = append(sequenceNumberRule, sequenceNumberItem)
	}

	logs, sub, err := _OnRampWithMessageTransformer.contract.WatchLogs(opts, "CCIPMessageSent", destChainSelectorRule, sequenceNumberRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OnRampWithMessageTransformerCCIPMessageSent)
				if err := _OnRampWithMessageTransformer.contract.UnpackLog(event, "CCIPMessageSent", log); err != nil {
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

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) ParseCCIPMessageSent(log types.Log) (*OnRampWithMessageTransformerCCIPMessageSent, error) {
	event := new(OnRampWithMessageTransformerCCIPMessageSent)
	if err := _OnRampWithMessageTransformer.contract.UnpackLog(event, "CCIPMessageSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OnRampWithMessageTransformerConfigSetIterator struct {
	Event *OnRampWithMessageTransformerConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OnRampWithMessageTransformerConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnRampWithMessageTransformerConfigSet)
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
		it.Event = new(OnRampWithMessageTransformerConfigSet)
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

func (it *OnRampWithMessageTransformerConfigSetIterator) Error() error {
	return it.fail
}

func (it *OnRampWithMessageTransformerConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OnRampWithMessageTransformerConfigSet struct {
	StaticConfig  OnRampStaticConfig
	DynamicConfig OnRampDynamicConfig
	Raw           types.Log
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) FilterConfigSet(opts *bind.FilterOpts) (*OnRampWithMessageTransformerConfigSetIterator, error) {

	logs, sub, err := _OnRampWithMessageTransformer.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &OnRampWithMessageTransformerConfigSetIterator{contract: _OnRampWithMessageTransformer.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *OnRampWithMessageTransformerConfigSet) (event.Subscription, error) {

	logs, sub, err := _OnRampWithMessageTransformer.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OnRampWithMessageTransformerConfigSet)
				if err := _OnRampWithMessageTransformer.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) ParseConfigSet(log types.Log) (*OnRampWithMessageTransformerConfigSet, error) {
	event := new(OnRampWithMessageTransformerConfigSet)
	if err := _OnRampWithMessageTransformer.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OnRampWithMessageTransformerDestChainConfigSetIterator struct {
	Event *OnRampWithMessageTransformerDestChainConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OnRampWithMessageTransformerDestChainConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnRampWithMessageTransformerDestChainConfigSet)
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
		it.Event = new(OnRampWithMessageTransformerDestChainConfigSet)
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

func (it *OnRampWithMessageTransformerDestChainConfigSetIterator) Error() error {
	return it.fail
}

func (it *OnRampWithMessageTransformerDestChainConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OnRampWithMessageTransformerDestChainConfigSet struct {
	DestChainSelector uint64
	SequenceNumber    uint64
	Router            common.Address
	AllowlistEnabled  bool
	Raw               types.Log
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) FilterDestChainConfigSet(opts *bind.FilterOpts, destChainSelector []uint64) (*OnRampWithMessageTransformerDestChainConfigSetIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _OnRampWithMessageTransformer.contract.FilterLogs(opts, "DestChainConfigSet", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &OnRampWithMessageTransformerDestChainConfigSetIterator{contract: _OnRampWithMessageTransformer.contract, event: "DestChainConfigSet", logs: logs, sub: sub}, nil
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) WatchDestChainConfigSet(opts *bind.WatchOpts, sink chan<- *OnRampWithMessageTransformerDestChainConfigSet, destChainSelector []uint64) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _OnRampWithMessageTransformer.contract.WatchLogs(opts, "DestChainConfigSet", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OnRampWithMessageTransformerDestChainConfigSet)
				if err := _OnRampWithMessageTransformer.contract.UnpackLog(event, "DestChainConfigSet", log); err != nil {
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

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) ParseDestChainConfigSet(log types.Log) (*OnRampWithMessageTransformerDestChainConfigSet, error) {
	event := new(OnRampWithMessageTransformerDestChainConfigSet)
	if err := _OnRampWithMessageTransformer.contract.UnpackLog(event, "DestChainConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OnRampWithMessageTransformerFeeTokenWithdrawnIterator struct {
	Event *OnRampWithMessageTransformerFeeTokenWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OnRampWithMessageTransformerFeeTokenWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnRampWithMessageTransformerFeeTokenWithdrawn)
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
		it.Event = new(OnRampWithMessageTransformerFeeTokenWithdrawn)
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

func (it *OnRampWithMessageTransformerFeeTokenWithdrawnIterator) Error() error {
	return it.fail
}

func (it *OnRampWithMessageTransformerFeeTokenWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OnRampWithMessageTransformerFeeTokenWithdrawn struct {
	FeeAggregator common.Address
	FeeToken      common.Address
	Amount        *big.Int
	Raw           types.Log
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) FilterFeeTokenWithdrawn(opts *bind.FilterOpts, feeAggregator []common.Address, feeToken []common.Address) (*OnRampWithMessageTransformerFeeTokenWithdrawnIterator, error) {

	var feeAggregatorRule []interface{}
	for _, feeAggregatorItem := range feeAggregator {
		feeAggregatorRule = append(feeAggregatorRule, feeAggregatorItem)
	}
	var feeTokenRule []interface{}
	for _, feeTokenItem := range feeToken {
		feeTokenRule = append(feeTokenRule, feeTokenItem)
	}

	logs, sub, err := _OnRampWithMessageTransformer.contract.FilterLogs(opts, "FeeTokenWithdrawn", feeAggregatorRule, feeTokenRule)
	if err != nil {
		return nil, err
	}
	return &OnRampWithMessageTransformerFeeTokenWithdrawnIterator{contract: _OnRampWithMessageTransformer.contract, event: "FeeTokenWithdrawn", logs: logs, sub: sub}, nil
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) WatchFeeTokenWithdrawn(opts *bind.WatchOpts, sink chan<- *OnRampWithMessageTransformerFeeTokenWithdrawn, feeAggregator []common.Address, feeToken []common.Address) (event.Subscription, error) {

	var feeAggregatorRule []interface{}
	for _, feeAggregatorItem := range feeAggregator {
		feeAggregatorRule = append(feeAggregatorRule, feeAggregatorItem)
	}
	var feeTokenRule []interface{}
	for _, feeTokenItem := range feeToken {
		feeTokenRule = append(feeTokenRule, feeTokenItem)
	}

	logs, sub, err := _OnRampWithMessageTransformer.contract.WatchLogs(opts, "FeeTokenWithdrawn", feeAggregatorRule, feeTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OnRampWithMessageTransformerFeeTokenWithdrawn)
				if err := _OnRampWithMessageTransformer.contract.UnpackLog(event, "FeeTokenWithdrawn", log); err != nil {
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

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) ParseFeeTokenWithdrawn(log types.Log) (*OnRampWithMessageTransformerFeeTokenWithdrawn, error) {
	event := new(OnRampWithMessageTransformerFeeTokenWithdrawn)
	if err := _OnRampWithMessageTransformer.contract.UnpackLog(event, "FeeTokenWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OnRampWithMessageTransformerOwnershipTransferRequestedIterator struct {
	Event *OnRampWithMessageTransformerOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OnRampWithMessageTransformerOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnRampWithMessageTransformerOwnershipTransferRequested)
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
		it.Event = new(OnRampWithMessageTransformerOwnershipTransferRequested)
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

func (it *OnRampWithMessageTransformerOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *OnRampWithMessageTransformerOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OnRampWithMessageTransformerOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OnRampWithMessageTransformerOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OnRampWithMessageTransformer.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OnRampWithMessageTransformerOwnershipTransferRequestedIterator{contract: _OnRampWithMessageTransformer.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OnRampWithMessageTransformerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OnRampWithMessageTransformer.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OnRampWithMessageTransformerOwnershipTransferRequested)
				if err := _OnRampWithMessageTransformer.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) ParseOwnershipTransferRequested(log types.Log) (*OnRampWithMessageTransformerOwnershipTransferRequested, error) {
	event := new(OnRampWithMessageTransformerOwnershipTransferRequested)
	if err := _OnRampWithMessageTransformer.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OnRampWithMessageTransformerOwnershipTransferredIterator struct {
	Event *OnRampWithMessageTransformerOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OnRampWithMessageTransformerOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnRampWithMessageTransformerOwnershipTransferred)
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
		it.Event = new(OnRampWithMessageTransformerOwnershipTransferred)
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

func (it *OnRampWithMessageTransformerOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *OnRampWithMessageTransformerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OnRampWithMessageTransformerOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OnRampWithMessageTransformerOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OnRampWithMessageTransformer.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OnRampWithMessageTransformerOwnershipTransferredIterator{contract: _OnRampWithMessageTransformer.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OnRampWithMessageTransformerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OnRampWithMessageTransformer.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OnRampWithMessageTransformerOwnershipTransferred)
				if err := _OnRampWithMessageTransformer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) ParseOwnershipTransferred(log types.Log) (*OnRampWithMessageTransformerOwnershipTransferred, error) {
	event := new(OnRampWithMessageTransformerOwnershipTransferred)
	if err := _OnRampWithMessageTransformer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformer) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _OnRampWithMessageTransformer.abi.Events["AllowListAdminSet"].ID:
		return _OnRampWithMessageTransformer.ParseAllowListAdminSet(log)
	case _OnRampWithMessageTransformer.abi.Events["AllowListSendersAdded"].ID:
		return _OnRampWithMessageTransformer.ParseAllowListSendersAdded(log)
	case _OnRampWithMessageTransformer.abi.Events["AllowListSendersRemoved"].ID:
		return _OnRampWithMessageTransformer.ParseAllowListSendersRemoved(log)
	case _OnRampWithMessageTransformer.abi.Events["CCIPMessageSent"].ID:
		return _OnRampWithMessageTransformer.ParseCCIPMessageSent(log)
	case _OnRampWithMessageTransformer.abi.Events["ConfigSet"].ID:
		return _OnRampWithMessageTransformer.ParseConfigSet(log)
	case _OnRampWithMessageTransformer.abi.Events["DestChainConfigSet"].ID:
		return _OnRampWithMessageTransformer.ParseDestChainConfigSet(log)
	case _OnRampWithMessageTransformer.abi.Events["FeeTokenWithdrawn"].ID:
		return _OnRampWithMessageTransformer.ParseFeeTokenWithdrawn(log)
	case _OnRampWithMessageTransformer.abi.Events["OwnershipTransferRequested"].ID:
		return _OnRampWithMessageTransformer.ParseOwnershipTransferRequested(log)
	case _OnRampWithMessageTransformer.abi.Events["OwnershipTransferred"].ID:
		return _OnRampWithMessageTransformer.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (OnRampWithMessageTransformerAllowListAdminSet) Topic() common.Hash {
	return common.HexToHash("0xb8c9b44ae5b5e3afb195f67391d9ff50cb904f9c0fa5fd520e497a97c1aa5a1e")
}

func (OnRampWithMessageTransformerAllowListSendersAdded) Topic() common.Hash {
	return common.HexToHash("0x330939f6eafe8bb516716892fe962ff19770570838686e6579dbc1cc51fc3281")
}

func (OnRampWithMessageTransformerAllowListSendersRemoved) Topic() common.Hash {
	return common.HexToHash("0xc237ec1921f855ccd5e9a5af9733f2d58943a5a8501ec5988e305d7a4d421586")
}

func (OnRampWithMessageTransformerCCIPMessageSent) Topic() common.Hash {
	return common.HexToHash("0x192442a2b2adb6a7948f097023cb6b57d29d3a7a5dd33e6666d33c39cc456f32")
}

func (OnRampWithMessageTransformerConfigSet) Topic() common.Hash {
	return common.HexToHash("0xc7372d2d886367d7bb1b0e0708a5436f2c91d6963de210eb2dc1ec2ecd6d21f1")
}

func (OnRampWithMessageTransformerDestChainConfigSet) Topic() common.Hash {
	return common.HexToHash("0xd5ad72bc37dc7a80a8b9b9df20500046fd7341adb1be2258a540466fdd7dcef5")
}

func (OnRampWithMessageTransformerFeeTokenWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x508d7d183612c18fc339b42618912b9fa3239f631dd7ec0671f950200a0fa66e")
}

func (OnRampWithMessageTransformerOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (OnRampWithMessageTransformerOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformer) Address() common.Address {
	return _OnRampWithMessageTransformer.address
}

type OnRampWithMessageTransformerInterface interface {
	GetAllowedSendersList(opts *bind.CallOpts, destChainSelector uint64) (GetAllowedSendersList,

		error)

	GetDestChainConfig(opts *bind.CallOpts, destChainSelector uint64) (GetDestChainConfig,

		error)

	GetDynamicConfig(opts *bind.CallOpts) (OnRampDynamicConfig, error)

	GetExpectedNextSequenceNumber(opts *bind.CallOpts, destChainSelector uint64) (uint64, error)

	GetFee(opts *bind.CallOpts, destChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error)

	GetMessageTransformer(opts *bind.CallOpts) (common.Address, error)

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

	SetMessageTransformer(opts *bind.TransactOpts, messageTransformerAddr common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	WithdrawFeeTokens(opts *bind.TransactOpts, feeTokens []common.Address) (*types.Transaction, error)

	FilterAllowListAdminSet(opts *bind.FilterOpts, allowlistAdmin []common.Address) (*OnRampWithMessageTransformerAllowListAdminSetIterator, error)

	WatchAllowListAdminSet(opts *bind.WatchOpts, sink chan<- *OnRampWithMessageTransformerAllowListAdminSet, allowlistAdmin []common.Address) (event.Subscription, error)

	ParseAllowListAdminSet(log types.Log) (*OnRampWithMessageTransformerAllowListAdminSet, error)

	FilterAllowListSendersAdded(opts *bind.FilterOpts, destChainSelector []uint64) (*OnRampWithMessageTransformerAllowListSendersAddedIterator, error)

	WatchAllowListSendersAdded(opts *bind.WatchOpts, sink chan<- *OnRampWithMessageTransformerAllowListSendersAdded, destChainSelector []uint64) (event.Subscription, error)

	ParseAllowListSendersAdded(log types.Log) (*OnRampWithMessageTransformerAllowListSendersAdded, error)

	FilterAllowListSendersRemoved(opts *bind.FilterOpts, destChainSelector []uint64) (*OnRampWithMessageTransformerAllowListSendersRemovedIterator, error)

	WatchAllowListSendersRemoved(opts *bind.WatchOpts, sink chan<- *OnRampWithMessageTransformerAllowListSendersRemoved, destChainSelector []uint64) (event.Subscription, error)

	ParseAllowListSendersRemoved(log types.Log) (*OnRampWithMessageTransformerAllowListSendersRemoved, error)

	FilterCCIPMessageSent(opts *bind.FilterOpts, destChainSelector []uint64, sequenceNumber []uint64) (*OnRampWithMessageTransformerCCIPMessageSentIterator, error)

	WatchCCIPMessageSent(opts *bind.WatchOpts, sink chan<- *OnRampWithMessageTransformerCCIPMessageSent, destChainSelector []uint64, sequenceNumber []uint64) (event.Subscription, error)

	ParseCCIPMessageSent(log types.Log) (*OnRampWithMessageTransformerCCIPMessageSent, error)

	FilterConfigSet(opts *bind.FilterOpts) (*OnRampWithMessageTransformerConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *OnRampWithMessageTransformerConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*OnRampWithMessageTransformerConfigSet, error)

	FilterDestChainConfigSet(opts *bind.FilterOpts, destChainSelector []uint64) (*OnRampWithMessageTransformerDestChainConfigSetIterator, error)

	WatchDestChainConfigSet(opts *bind.WatchOpts, sink chan<- *OnRampWithMessageTransformerDestChainConfigSet, destChainSelector []uint64) (event.Subscription, error)

	ParseDestChainConfigSet(log types.Log) (*OnRampWithMessageTransformerDestChainConfigSet, error)

	FilterFeeTokenWithdrawn(opts *bind.FilterOpts, feeAggregator []common.Address, feeToken []common.Address) (*OnRampWithMessageTransformerFeeTokenWithdrawnIterator, error)

	WatchFeeTokenWithdrawn(opts *bind.WatchOpts, sink chan<- *OnRampWithMessageTransformerFeeTokenWithdrawn, feeAggregator []common.Address, feeToken []common.Address) (event.Subscription, error)

	ParseFeeTokenWithdrawn(log types.Log) (*OnRampWithMessageTransformerFeeTokenWithdrawn, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OnRampWithMessageTransformerOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OnRampWithMessageTransformerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*OnRampWithMessageTransformerOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OnRampWithMessageTransformerOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OnRampWithMessageTransformerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*OnRampWithMessageTransformerOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
