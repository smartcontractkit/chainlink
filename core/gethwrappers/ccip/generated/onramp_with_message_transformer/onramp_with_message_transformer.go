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
	Bin: "0x610100604052346105a1576140f08038038061001a816105db565b92833981019080820361016081126105a157608081126105a15761003c6105bc565b9161004681610600565b835260208101516001600160a01b03811681036105a1576020840190815261007060408301610614565b916040850192835260a061008660608301610614565b6060870190815294607f1901126105a15760405160a081016001600160401b038111828210176105a6576040526100bf60808301610614565b81526100cd60a08301610628565b602082019081526100e060c08401610614565b90604083019182526100f460e08501610614565b92606081019384526101096101008601610614565b608082019081526101208601519095906001600160401b0381116105a15781018b601f820112156105a1578051906001600160401b0382116105a6578160051b602001610155906105db565b9c8d838152602001926060028201602001918183116105a157602001925b828410610533575050505061014061018b9101610614565b98331561052257600180546001600160a01b0319163317905580516001600160401b0316158015610510575b80156104fe575b80156104ec575b6104bf57516001600160401b0316608081905295516001600160a01b0390811660a08190529751811660c08190529851811660e081905282519091161580156104da575b80156104d0575b6104bf57815160028054855160ff60a01b90151560a01b166001600160a01b039384166001600160a81b0319909216919091171790558451600380549183166001600160a01b03199283161790558651600480549184169183169190911790558751600580549190931691161790557fc7372d2d886367d7bb1b0e0708a5436f2c91d6963de210eb2dc1ec2ecd6d21f1986101209860606102af6105bc565b8a8152602080820193845260408083019586529290910194855281519a8b5291516001600160a01b03908116928b019290925291518116918901919091529051811660608801529051811660808701529051151560a08601529051811660c08501529051811660e0840152905116610100820152a160005b82518110156104185761033a8184610635565b516001600160401b0361034d8386610635565b5151169081156104035760008281526006602090815260409182902081840151815494840151600160401b600160e81b03198616604883901b600160481b600160e81b031617901515851b68ff000000000000000016179182905583516001600160401b0390951685526001600160a01b031691840191909152811c60ff1615159082015260019291907fd5ad72bc37dc7a80a8b9b9df20500046fd7341adb1be2258a540466fdd7dcef590606090a201610327565b5063c35aa79d60e01b60005260045260246000fd5b506001600160a01b031680156104ae57600780546001600160a01b031916919091179055604051613a90908161066082396080518181816103f901528181610b3001528181611f8f0152612779015260a051818181611fc8015281816124e501526127b2015260c051818181610da10152818161200401526127ee015260e0518181816120400152818161282a0152612e290152f35b6342bcdf7f60e11b60005260046000fd5b6306b7c75960e31b60005260046000fd5b5082511515610210565b5084516001600160a01b031615610209565b5088516001600160a01b0316156101c5565b5087516001600160a01b0316156101be565b5086516001600160a01b0316156101b7565b639b15e16f60e01b60005260046000fd5b6060848303126105a15760405190606082016001600160401b038111838210176105a65760405261056385610600565b82526020850151906001600160a01b03821682036105a1578260209283606095015261059160408801610628565b6040820152815201930192610173565b600080fd5b634e487b7160e01b600052604160045260246000fd5b60405190608082016001600160401b038111838210176105a657604052565b6040519190601f01601f191682016001600160401b038111838210176105a657604052565b51906001600160401b03821682036105a157565b51906001600160a01b03821682036105a157565b519081151582036105a157565b80518210156106495760209160051b010190565b634e487b7160e01b600052603260045260246000fdfe608080604052600436101561001357600080fd5b600090813560e01c90816306285c69146127135750806315777ab2146126c1578063181f5a771461264257806320487ded146124095780632716072b1461215957806327e936f114611d5357806348a98aa414611cd05780635cb80c5d14611a1357806365b81aab146119635780636def4ce7146118d45780637437ff9f146117b757806379ba5097146116d25780638da5cb5b146116805780639041be3d146115d3578063972b461214611505578063c9b146b31461113c578063df0aa9e914610242578063f2fde38b146101555763fbca3b74146100f257600080fd5b346101525760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101525760049061012c612a1f565b507f9e7177c8000000000000000000000000000000000000000000000000000000008152fd5b80fd5b50346101525760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101525773ffffffffffffffffffffffffffffffffffffffff6101a2612a75565b6101aa613355565b1633811461021a57807fffffffffffffffffffffffff000000000000000000000000000000000000000083541617825573ffffffffffffffffffffffffffffffffffffffff600154167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12788380a380f35b6004827fdad89dca000000000000000000000000000000000000000000000000000000008152fd5b50346101525760807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101525761027a612a1f565b67ffffffffffffffff602435116111385760a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc60243536030112611138576102c1612a98565b60025460ff8160a01c16611110577fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff16740100000000000000000000000000000000000000001760025567ffffffffffffffff8216835260066020526040832073ffffffffffffffffffffffffffffffffffffffff8216156110e857805460ff8160401c1661107a575b60481c73ffffffffffffffffffffffffffffffffffffffff1633036110525773ffffffffffffffffffffffffffffffffffffffff6003541680610fda575b50805467ffffffffffffffff811667ffffffffffffffff8114610fad579067ffffffffffffffff60017fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000009493011692839116179055604051906103eb826128e9565b84825267ffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016602083015267ffffffffffffffff84166040830152606082015283608082015261044c6024803501602435600401612fc1565b61045b60046024350180612fc1565b610469606460243501612ecd565b9361047e604460243501602435600401613012565b9490507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe06104c46104ae87612a50565b966104bc604051988961293e565b808852612a50565b018a5b818110610f9657505061051893929161050c91604051986104e78a612905565b895273ffffffffffffffffffffffffffffffffffffffff8a1660208a01523691613092565b60408701523691613092565b606084015273ffffffffffffffffffffffffffffffffffffffff602092604051610542858261293e565b88815260808601521660a084015260443560c08401528560e0840152610100830152610578604460243501602435600401613012565b61058481969296612a50565b90610592604051928361293e565b8082528382018097368360061b820111610f925780915b8360061b82018310610f5b5750505050865b6105cf604460243501602435600401613012565b9050811015610945576105e28183612f7e565b51906105fd6105f660046024350180612fc1565b3691613092565b91610606613066565b50858101511561091d5773ffffffffffffffffffffffffffffffffffffffff61063181835116612dca565b169283158015610873575b61083057908a6106e992888873ffffffffffffffffffffffffffffffffffffffff8d8387015182808951169260405197610675896128e9565b885267ffffffffffffffff87890196168652816040890191168152606088019283526080880193845267ffffffffffffffff6040519b8c998a997f9a4575b9000000000000000000000000000000000000000000000000000000008b5260048b01525160a060248b015260c48a01906129dc565b965116604488015251166064860152516084850152511660a4830152038183885af1918215610825578b9261077e575b5060019382849289806107779651930151910151916040519361073b856128e9565b84528a84015260408301526060820152604051610758898261293e565b8c81526080820152610100890151906107718383612f7e565b52612f7e565b50016105bb565b91503d808c843e61078f818461293e565b820191878184031261081d5780519067ffffffffffffffff821161082157019360408584031261081d57604051916107c683612922565b855167ffffffffffffffff811161081957846107e39188016130c9565b8352888601519267ffffffffffffffff84116108195761080b610777958795600199016130c9565b8a8201529350915093610719565b8d80fd5b8b80fd5b8c80fd5b6040513d8d823e3d90fd5b60248b73ffffffffffffffffffffffffffffffffffffffff8451167fbf16aab6000000000000000000000000000000000000000000000000000000008252600452fd5b506040517f01ffc9a70000000000000000000000000000000000000000000000000000000081527faff2afbf0000000000000000000000000000000000000000000000000000000060048201528781602481885afa908115610912578c916108dd575b501561063c565b90508781813d831161090b575b6108f4818361293e565b8101031261081d5761090590612b57565b386108d6565b503d6108ea565b6040513d8e823e3d90fd5b60048a7f5cf04449000000000000000000000000000000000000000000000000000000008152fd5b50957f430d138c0000000000000000000000000000000000000000000000000000000081939497839773ffffffffffffffffffffffffffffffffffffffff600254169187610a378c610a0761099e606460243501612ecd565b9161010073ffffffffffffffffffffffffffffffffffffffff6109cb608460243501602435600401612fc1565b92909301519467ffffffffffffffff6040519e8f9d8e521660048d01521660248b015260443560448b015260c060648b015260c48a0191612bb4565b907ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc88830301608489015261310b565b917ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc8684030160a48701525191828152019190855b89828210610f20575050505082809103915afa938415610f155782839084938597610e1b575b5060e089015215610d325750815b67ffffffffffffffff6080885101911690526080860152805b61010086015151811015610aef5780610ad460019286612f7e565b516080610ae6836101008b0151612f7e565b51015201610ab9565b5083610afa866133d0565b91604051848101907f130ac867e79e2789f923760a88743d292acdf7002139a588206e2260f73f7321825267ffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016604082015267ffffffffffffffff8416606082015230608082015260808152610b7a60a08261293e565b51902073ffffffffffffffffffffffffffffffffffffffff8585015116845167ffffffffffffffff6080816060840151169201511673ffffffffffffffffffffffffffffffffffffffff60a08801511660c088015191604051938a850195865260408501526060840152608083015260a082015260a08152610bfd60c08261293e565b51902060608501518681519101206040860151878151910120610100870151604051610c6381610c378c8201948d8652604083019061310b565b037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0810183528261293e565b51902091608088015189815191012093604051958a870197885260408701526060860152608085015260a084015260c083015260e082015260e08152610cab6101008261293e565b51902082515267ffffffffffffffff60608351015116907f192442a2b2adb6a7948f097023cb6b57d29d3a7a5dd33e6666d33c39cc456f3267ffffffffffffffff60405192169180610cfd86826131f6565b0390a37fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff600254166002555151604051908152f35b73ffffffffffffffffffffffffffffffffffffffff604051917fea458c0c00000000000000000000000000000000000000000000000000000000835267ffffffffffffffff8816600484015216602482015283816044818673ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165af1908115610e10578391610dd7575b50610aa0565b90508381813d8311610e09575b610dee818361293e565b81010312610e0557610dff906131e1565b38610dd1565b8280fd5b503d610de4565b6040513d85823e3d90fd5b9650505090503d8083863e610e30818661293e565b840190608085830312610e0557845194610e4b858201612b57565b95604082015167ffffffffffffffff8111610f0d5784610e6c9184016130c9565b9160608101519067ffffffffffffffff8211610f11570184601f82011215610f0d578051610e9981612a50565b95610ea7604051978861293e565b818752888088019260051b84010192818411610f0957898101925b848410610ed85750505050509590929538610a92565b835167ffffffffffffffff8111610f05578b91610efa858480948701016130c9565b815201930192610ec2565b8a80fd5b8880fd5b8580fd5b8680fd5b6040513d84823e3d90fd5b8351805173ffffffffffffffffffffffffffffffffffffffff1686528101518186015289975088965060409094019390920191600101610a6c565b604083360312610f0557866040918251610f7481612922565b610f7d86612abb565b815282860135838201528152019201916105a9565b8980fd5b602090610fa1613066565b82828a010152016104c7565b6024867f4e487b710000000000000000000000000000000000000000000000000000000081526011600452fd5b803b1561104e578460405180927fe0a0e5060000000000000000000000000000000000000000000000000000000082528183816110206024356004018b60048401612bf3565b03925af180156110435715610389578461103c9195929561293e565b9238610389565b6040513d87823e3d90fd5b8480fd5b6004847f1c0a3529000000000000000000000000000000000000000000000000000000008152fd5b73ffffffffffffffffffffffffffffffffffffffff8316600090815260028301602052604090205461034b5760248573ffffffffffffffffffffffffffffffffffffffff857fd0d2597600000000000000000000000000000000000000000000000000000000835216600452fd5b6004847fa4ec7479000000000000000000000000000000000000000000000000000000008152fd5b6004847f3ee5aeb5000000000000000000000000000000000000000000000000000000008152fd5b5080fd5b50346101525760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101525760043567ffffffffffffffff81116111385761118c903690600401612adc565b73ffffffffffffffffffffffffffffffffffffffff6001541633036114bd575b919081907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8181360301915b848110156114b9578060051b8201358381121561104e5782019160808336031261104e57604051946112088661289e565b61121184612a3b565b865261121f60208501612a68565b9660208701978852604085013567ffffffffffffffff8111610e05576112489036908701612f19565b9460408801958652606081013567ffffffffffffffff81116114b55761127091369101612f19565b60608801908152875167ffffffffffffffff1683526006602052604080842099518a547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff169015159182901b68ff000000000000000016178a5590959081515161138d575b5095976001019550815b8551805182101561131e579061131773ffffffffffffffffffffffffffffffffffffffff61130f83600195612f7e565b511689613825565b50016112df565b5050959096945060019291935190815161133e575b5050019392936111d7565b61138367ffffffffffffffff7fc237ec1921f855ccd5e9a5af9733f2d58943a5a8501ec5988e305d7a4d42158692511692604051918291602083526020830190612b0d565b0390a23880611333565b9893959296919094979860001461147e57600184019591875b86518051821015611423576113d08273ffffffffffffffffffffffffffffffffffffffff92612f7e565b511680156113ec57906113e56001928a613794565b50016113a6565b60248a67ffffffffffffffff8e51167f463258ff000000000000000000000000000000000000000000000000000000008252600452fd5b50509692955090929796937f330939f6eafe8bb516716892fe962ff19770570838686e6579dbc1cc51fc328161147467ffffffffffffffff8a51169251604051918291602083526020830190612b0d565b0390a238806112d5565b60248767ffffffffffffffff8b51167f463258ff000000000000000000000000000000000000000000000000000000008252600452fd5b8380fd5b8380f35b73ffffffffffffffffffffffffffffffffffffffff600554163303156111ac576004837f905d7d9b000000000000000000000000000000000000000000000000000000008152fd5b50346101525760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101525767ffffffffffffffff611546612a1f565b16808252600660205260ff604083205460401c16908252600660205260016040832001916040518093849160208254918281520191845260208420935b8181106115ba5750506115989250038361293e565b6115b660405192839215158352604060208401526040830190612b0d565b0390f35b8454835260019485019487945060209093019201611583565b50346101525760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101525767ffffffffffffffff611614612a1f565b1681526006602052600167ffffffffffffffff604083205416019067ffffffffffffffff82116116535760208267ffffffffffffffff60405191168152f35b807f4e487b7100000000000000000000000000000000000000000000000000000000602492526011600452fd5b503461015257807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261015257602073ffffffffffffffffffffffffffffffffffffffff60015416604051908152f35b503461015257807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261015257805473ffffffffffffffffffffffffffffffffffffffff8116330361178f577fffffffffffffffffffffffff000000000000000000000000000000000000000060015491338284161760015516825573ffffffffffffffffffffffffffffffffffffffff3391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08380a380f35b6004827f02b543c6000000000000000000000000000000000000000000000000000000008152fd5b503461015257807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610152576117ee612eee565b5060a06040516117fd816128e9565b60ff60025473ffffffffffffffffffffffffffffffffffffffff81168352831c161515602082015273ffffffffffffffffffffffffffffffffffffffff60035416604082015273ffffffffffffffffffffffffffffffffffffffff60045416606082015273ffffffffffffffffffffffffffffffffffffffff6005541660808201526118d2604051809273ffffffffffffffffffffffffffffffffffffffff60808092828151168552602081015115156020860152826040820151166040860152826060820151166060860152015116910152565bf35b50346101525760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261015257604060609167ffffffffffffffff61191a612a1f565b1681526006602052205473ffffffffffffffffffffffffffffffffffffffff6040519167ffffffffffffffff8116835260ff8160401c161515602084015260481c166040820152f35b50346101525760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101525773ffffffffffffffffffffffffffffffffffffffff6119b0612a75565b6119b8613355565b1680156119eb577fffffffffffffffffffffffff0000000000000000000000000000000000000000600754161760075580f35b6004827f8579befe000000000000000000000000000000000000000000000000000000008152fd5b50346101525760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101525760043567ffffffffffffffff811161113857611a63903690600401612adc565b9073ffffffffffffffffffffffffffffffffffffffff6004541690835b83811015611ccc5773ffffffffffffffffffffffffffffffffffffffff611aab8260051b8401612ecd565b1690604051917f70a08231000000000000000000000000000000000000000000000000000000008352306004840152602083602481845afa928315611cc1578793611c8e575b5082611b03575b506001915001611a80565b8460405193611b9f60208601957fa9059cbb00000000000000000000000000000000000000000000000000000000875283602482015282604482015260448152611b4e60648261293e565b8a80604098895193611b608b8661293e565b602085527f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65646020860152519082895af1611b986133a0565b90866139b7565b805180611bdb575b505060207f508d7d183612c18fc339b42618912b9fa3239f631dd7ec0671f950200a0fa66e9160019651908152a338611af8565b819294959693509060209181010312610f09576020611bfa9101612b57565b15611c0b5792919085903880611ba7565b608490517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f742073756363656564000000000000000000000000000000000000000000006064820152fd5b9092506020813d8211611cb9575b81611ca96020938361293e565b81010312610f1157519138611af1565b3d9150611c9c565b6040513d89823e3d90fd5b8480f35b50346101525760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261015257611d08612a1f565b506024359073ffffffffffffffffffffffffffffffffffffffff82168203610152576020611d3583612dca565b73ffffffffffffffffffffffffffffffffffffffff60405191168152f35b50346101525760a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261015257604051611d8f816128e9565b611d97612a75565b81526024358015158103610e05576020820190815260443573ffffffffffffffffffffffffffffffffffffffff811681036114b55760408301908152611ddb612a98565b90606084019182526084359273ffffffffffffffffffffffffffffffffffffffff84168403610f0d5760808501938452611e13613355565b73ffffffffffffffffffffffffffffffffffffffff85511615801561213a575b8015612130575b612108579273ffffffffffffffffffffffffffffffffffffffff859381809461012097827fc7372d2d886367d7bb1b0e0708a5436f2c91d6963de210eb2dc1ec2ecd6d21f19a51167fffffffffffffffffffffffff000000000000000000000000000000000000000060025416176002555115157fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff74ff00000000000000000000000000000000000000006002549260a01b1691161760025551167fffffffffffffffffffffffff0000000000000000000000000000000000000000600354161760035551167fffffffffffffffffffffffff0000000000000000000000000000000000000000600454161760045551167fffffffffffffffffffffffff0000000000000000000000000000000000000000600554161760055561210460405191611f848361289e565b67ffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016835273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016602084015273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016604084015273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001660608401526120b4604051809473ffffffffffffffffffffffffffffffffffffffff6060809267ffffffffffffffff8151168552826020820151166020860152826040820151166040860152015116910152565b608083019073ffffffffffffffffffffffffffffffffffffffff60808092828151168552602081015115156020860152826040820151166040860152826060820151166060860152015116910152565ba180f35b6004867f35be3ac8000000000000000000000000000000000000000000000000000000008152fd5b5080511515611e3a565b5073ffffffffffffffffffffffffffffffffffffffff83511615611e33565b50346101525760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610152576004359067ffffffffffffffff821161015257366023830112156101525781600401356121b581612a50565b926121c3604051948561293e565b81845260246060602086019302820101903682116114b557602401915b818310612361575050506121f2613355565b805b825181101561235d576122078184612f7e565b5167ffffffffffffffff61221b8386612f7e565b51511690811561233157907fd5ad72bc37dc7a80a8b9b9df20500046fd7341adb1be2258a540466fdd7dcef5606060019493838752600660205260ff604088206122f2604060208501519483547fffffff0000000000000000000000000000000000000000ffffffffffffffffff7cffffffffffffffffffffffffffffffffffffffff0000000000000000008860481b1691161784550151151582907fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff68ff0000000000000000835492151560401b169116179055565b5473ffffffffffffffffffffffffffffffffffffffff6040519367ffffffffffffffff8316855216602084015260401c1615156040820152a2016121f4565b602484837fc35aa79d000000000000000000000000000000000000000000000000000000008252600452fd5b5080f35b6060833603126114b5576040516060810181811067ffffffffffffffff8211176123dc5760405261239184612a3b565b8152602084013573ffffffffffffffffffffffffffffffffffffffff81168103610f0d5791816060936020809401526123cc60408701612a68565b60408201528152019201916121e0565b6024867f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b50346101525760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261015257612441612a1f565b60243567ffffffffffffffff8111610e055760a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc8236030112610e05576040517f2cbc26bb00000000000000000000000000000000000000000000000000000000815277ffffffffffffffff000000000000000000000000000000008360801b16600482015260208160248173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165afa9081156126375784916125fd575b506125c7576125749160209173ffffffffffffffffffffffffffffffffffffffff60025416906040518095819482937fd8694ccd0000000000000000000000000000000000000000000000000000000084526004019060048401612bf3565b03915afa908115610f15578291612591575b602082604051908152f35b90506020813d6020116125bf575b816125ac6020938361293e565b8101031261113857602091505138612586565b3d915061259f565b60248367ffffffffffffffff847ffdbd6a7200000000000000000000000000000000000000000000000000000000835216600452fd5b90506020813d60201161262f575b816126186020938361293e565b810103126114b55761262990612b57565b38612515565b3d915061260b565b6040513d86823e3d90fd5b503461015257807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261015257506115b660405161268360408261293e565b601081527f4f6e52616d7020312e362e302d6465760000000000000000000000000000000060208201526040519182916020835260208301906129dc565b503461015257807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261015257602073ffffffffffffffffffffffffffffffffffffffff60075416604051908152f35b90503461113857817ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112611138578061274f60609261289e565b8281528260208201528260408201520152608060405161276e8161289e565b67ffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016815273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016602082015273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016604082015273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001660608201526118d2604051809273ffffffffffffffffffffffffffffffffffffffff6060809267ffffffffffffffff8151168552826020820151166020860152826040820151166040860152015116910152565b6080810190811067ffffffffffffffff8211176128ba57604052565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60a0810190811067ffffffffffffffff8211176128ba57604052565b610120810190811067ffffffffffffffff8211176128ba57604052565b6040810190811067ffffffffffffffff8211176128ba57604052565b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff8211176128ba57604052565b67ffffffffffffffff81116128ba57601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b60005b8381106129cc5750506000910152565b81810151838201526020016129bc565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f602093612a18815180928187528780880191016129b9565b0116010190565b6004359067ffffffffffffffff82168203612a3657565b600080fd5b359067ffffffffffffffff82168203612a3657565b67ffffffffffffffff81116128ba5760051b60200190565b35908115158203612a3657565b6004359073ffffffffffffffffffffffffffffffffffffffff82168203612a3657565b6064359073ffffffffffffffffffffffffffffffffffffffff82168203612a3657565b359073ffffffffffffffffffffffffffffffffffffffff82168203612a3657565b9181601f84011215612a365782359167ffffffffffffffff8311612a36576020808501948460051b010111612a3657565b906020808351928381520192019060005b818110612b2b5750505090565b825173ffffffffffffffffffffffffffffffffffffffff16845260209384019390920191600101612b1e565b51908115158203612a3657565b90357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe182360301811215612a3657016020813591019167ffffffffffffffff8211612a36578136038313612a3657565b601f82602094937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0938186528686013760008582860101520116010190565b9067ffffffffffffffff9093929316815260406020820152612c69612c2c612c1b8580612b64565b60a0604086015260e0850191612bb4565b612c396020860186612b64565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0858403016060860152612bb4565b9060408401357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe185360301811215612a365784016020813591019267ffffffffffffffff8211612a36578160061b36038413612a36578281037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0016080840152818152602001929060005b818110612d6a57505050612d378473ffffffffffffffffffffffffffffffffffffffff612d276060612d67979801612abb565b1660a08401526080810190612b64565b9160c07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc082860301910152612bb4565b90565b90919360408060019273ffffffffffffffffffffffffffffffffffffffff612d9189612abb565b16815260208881013590820152019501929101612cf4565b519073ffffffffffffffffffffffffffffffffffffffff82168203612a3657565b73ffffffffffffffffffffffffffffffffffffffff604051917fbbe4f6db00000000000000000000000000000000000000000000000000000000835216600482015260208160248173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165afa8015612ec157600090612e74575b73ffffffffffffffffffffffffffffffffffffffff91501690565b506020813d602011612eb9575b81612e8e6020938361293e565b81010312612a3657612eb473ffffffffffffffffffffffffffffffffffffffff91612da9565b612e59565b3d9150612e81565b6040513d6000823e3d90fd5b3573ffffffffffffffffffffffffffffffffffffffff81168103612a365790565b60405190612efb826128e9565b60006080838281528260208201528260408201528260608201520152565b9080601f83011215612a36578135612f3081612a50565b92612f3e604051948561293e565b81845260208085019260051b820101928311612a3657602001905b828210612f665750505090565b60208091612f7384612abb565b815201910190612f59565b8051821015612f925760209160051b010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe181360301821215612a36570180359067ffffffffffffffff8211612a3657602001918136038313612a3657565b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe181360301821215612a36570180359067ffffffffffffffff8211612a3657602001918160061b36038313612a3657565b60405190613073826128e9565b6060608083600081528260208201528260408201526000838201520152565b92919261309e8261297f565b916130ac604051938461293e565b829481845281830111612a36578281602093846000960137010152565b81601f82011215612a365780516130df8161297f565b926130ed604051948561293e565b81845260208284010111612a3657612d6791602080850191016129b9565b9080602083519182815201916020808360051b8301019401926000915b83831061313757505050505090565b90919293946020806131d2837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0866001960301875289519073ffffffffffffffffffffffffffffffffffffffff825116815260806131b76131a58685015160a08886015260a08501906129dc565b604085015184820360408601526129dc565b926060810151606084015201519060808184039101526129dc565b97019301930191939290613128565b519067ffffffffffffffff82168203612a3657565b90612d67916020815267ffffffffffffffff6080835180516020850152826020820151166040850152826040820151166060850152826060820151168285015201511660a082015273ffffffffffffffffffffffffffffffffffffffff60208301511660c08201526101006132ea6132b561328260408601516101a060e08701526101c08601906129dc565b60608601517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe086830301858701526129dc565b60808501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0858303016101208601526129dc565b9273ffffffffffffffffffffffffffffffffffffffff60a08201511661014084015260c081015161016084015260e08101516101808401520151906101a07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08285030191015261310b565b73ffffffffffffffffffffffffffffffffffffffff60015416330361337657565b7f2b5c74de0000000000000000000000000000000000000000000000000000000060005260046000fd5b3d156133cb573d906133b18261297f565b916133bf604051938461293e565b82523d6000602084013e565b606090565b600061347181926040516133e381612905565b6133eb612eee565b815283602082015260606040820152606080820152606060808201528360a08201528360c08201528360e082015260606101008201525073ffffffffffffffffffffffffffffffffffffffff60075416906040519485809481937f8a06fadb000000000000000000000000000000000000000000000000000000008352600483016131f6565b03925af180916000916134ce575b5090612d67576134ca6134906133a0565b6040519182917f828ebdfb0000000000000000000000000000000000000000000000000000000083526020600484015260248301906129dc565b0390fd5b3d8083833e6134dd818361293e565b810190602081830312610e055780519067ffffffffffffffff82116114b5570191828203926101a084126111385760a06040519461351a86612905565b126111385760405161352b816128e9565b8151815261353b602083016131e1565b602082015261354c604083016131e1565b604082015261355d606083016131e1565b606082015261356e608083016131e1565b6080820152845261358160a08201612da9565b602085015260c081015167ffffffffffffffff8111610e0557836135a69183016130c9565b604085015260e081015167ffffffffffffffff8111610e0557836135cb9183016130c9565b606085015261010081015167ffffffffffffffff8111610e0557836135f19183016130c9565b60808501526136036101208201612da9565b60a085015261014081015160c085015261016081015160e08501526101808101519067ffffffffffffffff8211610e05570182601f820112156111385780519161364c83612a50565b9361365a604051958661293e565b83855260208086019460051b84010192818411610e055760208101945b848610613690575050505050506101008201523861347f565b855167ffffffffffffffff811161104e57820160a07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0828603011261104e57604051906136dc826128e9565b6136e860208201612da9565b8252604081015167ffffffffffffffff8111610f115785602061370d928401016130c9565b6020830152606081015167ffffffffffffffff8111610f1157856020613735928401016130c9565b60408301526080810151606083015260a08101519067ffffffffffffffff8211610f11579161376c866020809694819601016130c9565b6080820152815201950194613677565b8054821015612f925760005260206000200190600090565b600082815260018201602052604090205461381e57805490680100000000000000008210156128ba57826138076137d284600180960185558461377c565b81939154907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9060031b92831b921b19161790565b905580549260005201602052604060002055600190565b5050600090565b90600182019181600052826020526040600020548015156000146139ae577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff810181811161397f578254907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820191821161397f57818103613948575b50505080548015613919577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01906138da828261377c565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82549160031b1b191690555560005260205260006040812055600190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b6139686139586137d2938661377c565b90549060031b1c9283928661377c565b9055600052836020526040600020553880806138a2565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b50505050600090565b91929015613a3257508151156139cb575090565b3b156139d45790565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152fd5b825190915015613a455750805190602001fd5b6134ca906040519182917f08c379a00000000000000000000000000000000000000000000000000000000083526020600484015260248301906129dc56fea164736f6c634300081a000a",
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
