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
	Bin: "0x610100604052346105a1576141988038038061001a816105db565b92833981019080820361016081126105a157608081126105a15761003c6105bc565b9161004681610600565b835260208101516001600160a01b03811681036105a1576020840190815261007060408301610614565b916040850192835260a061008660608301610614565b6060870190815294607f1901126105a15760405160a081016001600160401b038111828210176105a6576040526100bf60808301610614565b81526100cd60a08301610628565b602082019081526100e060c08401610614565b90604083019182526100f460e08501610614565b92606081019384526101096101008601610614565b608082019081526101208601519095906001600160401b0381116105a15781018b601f820112156105a1578051906001600160401b0382116105a6578160051b602001610155906105db565b9c8d838152602001926060028201602001918183116105a157602001925b828410610533575050505061014061018b9101610614565b98331561052257600180546001600160a01b0319163317905580516001600160401b0316158015610510575b80156104fe575b80156104ec575b6104bf57516001600160401b0316608081905295516001600160a01b0390811660a08190529751811660c08190529851811660e081905282519091161580156104da575b80156104d0575b6104bf57815160028054855160ff60a01b90151560a01b166001600160a01b039384166001600160a81b0319909216919091171790558451600380549183166001600160a01b03199283161790558651600480549184169183169190911790558751600580549190931691161790557fc7372d2d886367d7bb1b0e0708a5436f2c91d6963de210eb2dc1ec2ecd6d21f1986101209860606102af6105bc565b8a8152602080820193845260408083019586529290910194855281519a8b5291516001600160a01b03908116928b019290925291518116918901919091529051811660608801529051811660808701529051151560a08601529051811660c08501529051811660e0840152905116610100820152a160005b82518110156104185761033a8184610635565b516001600160401b0361034d8386610635565b5151169081156104035760008281526006602090815260409182902081840151815494840151600160401b600160e81b03198616604883901b600160481b600160e81b031617901515851b68ff000000000000000016179182905583516001600160401b0390951685526001600160a01b031691840191909152811c60ff1615159082015260019291907fd5ad72bc37dc7a80a8b9b9df20500046fd7341adb1be2258a540466fdd7dcef590606090a201610327565b5063c35aa79d60e01b60005260045260246000fd5b506001600160a01b031680156104ae57600780546001600160a01b031916919091179055604051613b38908161066082396080518181816103f901528181610b8c015281816120330152612821015260a05181818161206c0152818161258d015261285a015260c051818181610dfd015281816120a80152612896015260e0518181816120e4015281816128d20152612ed10152f35b6342bcdf7f60e11b60005260046000fd5b6306b7c75960e31b60005260046000fd5b5082511515610210565b5084516001600160a01b031615610209565b5088516001600160a01b0316156101c5565b5087516001600160a01b0316156101be565b5086516001600160a01b0316156101b7565b639b15e16f60e01b60005260046000fd5b6060848303126105a15760405190606082016001600160401b038111838210176105a65760405261056385610600565b82526020850151906001600160a01b03821682036105a1578260209283606095015261059160408801610628565b6040820152815201930192610173565b600080fd5b634e487b7160e01b600052604160045260246000fd5b60405190608082016001600160401b038111838210176105a657604052565b6040519190601f01601f191682016001600160401b038111838210176105a657604052565b51906001600160401b03821682036105a157565b51906001600160a01b03821682036105a157565b519081151582036105a157565b80518210156106495760209160051b010190565b634e487b7160e01b600052603260045260246000fdfe608080604052600436101561001357600080fd5b600090813560e01c90816306285c69146127bb5750806315777ab214612769578063181f5a77146126ea57806320487ded146124b15780632716072b1461220157806327e936f114611df757806348a98aa414611d745780635cb80c5d14611ab757806365b81aab14611a075780636def4ce7146119785780637437ff9f1461185b57806379ba5097146117765780638da5cb5b146117245780639041be3d14611677578063972b4612146115a9578063c9b146b3146111e0578063df0aa9e914610242578063f2fde38b146101555763fbca3b74146100f257600080fd5b346101525760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101525760049061012c612ac7565b507f9e7177c8000000000000000000000000000000000000000000000000000000008152fd5b80fd5b50346101525760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101525773ffffffffffffffffffffffffffffffffffffffff6101a2612b1d565b6101aa6133fd565b1633811461021a57807fffffffffffffffffffffffff000000000000000000000000000000000000000083541617825573ffffffffffffffffffffffffffffffffffffffff600154167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12788380a380f35b6004827fdad89dca000000000000000000000000000000000000000000000000000000008152fd5b50346101525760807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101525761027a612ac7565b67ffffffffffffffff60243511610e615760a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc60243536030112610e61576102c1612b40565b60025460ff8160a01c166111b8577fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff16740100000000000000000000000000000000000000001760025567ffffffffffffffff8216835260066020526040832073ffffffffffffffffffffffffffffffffffffffff82161561119057805460ff8160401c16611122575b60481c73ffffffffffffffffffffffffffffffffffffffff1633036110fa5773ffffffffffffffffffffffffffffffffffffffff6003541680611086575b50805467ffffffffffffffff811667ffffffffffffffff8114611059579067ffffffffffffffff60017fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000009493011692839116179055604051906103eb82612991565b84825267ffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016602083015267ffffffffffffffff84166040830152606082015283608082015261044c6024803501602435600401613069565b61045b60046024350180613069565b610469606460243501612f75565b9361047e6044602435016024356004016130ba565b9490507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe06104c46104ae87612af8565b966104bc60405198896129e6565b808852612af8565b018a5b81811061104257505061051893929161050c91604051986104e78a6129ad565b895273ffffffffffffffffffffffffffffffffffffffff8a1660208a0152369161313a565b6040870152369161313a565b606084015273ffffffffffffffffffffffffffffffffffffffff60209260405161054285826129e6565b88815260808601521660a084015260443560c08401528560e08401526101008301526105786044602435016024356004016130ba565b61058481969296612af8565b9061059260405192836129e6565b8082528382018097368360061b82011161103e5780915b8360061b82018310611003575050505073ffffffffffffffffffffffffffffffffffffffff908782600254166105e3606460243501612f75565b9067ffffffffffffffff8661069d610605608460243501602435600401613069565b9061066d61061860046024350180613069565b9290936040519c8d9a8b998a997f3a49bb49000000000000000000000000000000000000000000000000000000008b521660048a0152166024880152604435604488015260a0606488015260a4870191612c5c565b917ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc858403016084860152612c5c565b03915afa918215610ff8578889918a908b95610f7d575b50608088015260e0870152885b6106d56044602435016024356004016130ba565b9050811015610a2a576106e88184613026565b51906106f261310e565b508682015115610a025773ffffffffffffffffffffffffffffffffffffffff61071d81845116612e72565b16918215801561095f575b61091d57808c878a8a73ffffffffffffffffffffffffffffffffffffffff8f836107d59801518280895116926040519761076189612991565b885267ffffffffffffffff87890196168652816040890191168152606088019283526080880193845267ffffffffffffffff6040519b8c998a997f9a4575b9000000000000000000000000000000000000000000000000000000008b5260048b01525160a060248b015260c48a0190612a84565b965116604488015251166064860152516084850152511660a4830152038183885af1918215610912578d9261086b575b506001938284928b806108649651930151910151916040519361082785612991565b84528c840152604083015260608201528d604051906108468c836129e6565b815260808201526101008b01519061085e8383613026565b52613026565b50016106c1565b91503d808e843e61087c81846129e6565b820191898184031261090a5780519067ffffffffffffffff821161090e57019360408584031261090a57604051916108b3836129ca565b855167ffffffffffffffff811161090657846108d0918801613171565b83528a8601519267ffffffffffffffff8411610906576108f861086495879560019901613171565b8c8201529350915093610805565b8f80fd5b8d80fd5b8e80fd5b6040513d8f823e3d90fd5b517fbf16aab6000000000000000000000000000000000000000000000000000000008c5273ffffffffffffffffffffffffffffffffffffffff1660045260248bfd5b506040517f01ffc9a70000000000000000000000000000000000000000000000000000000081527faff2afbf0000000000000000000000000000000000000000000000000000000060048201528881602481875afa908115610912578d916109c9575b5015610728565b90508881813d83116109fb575b6109e081836129e6565b810103126109f7576109f190612bff565b386109c2565b8c80fd5b503d6109d6565b60048b7f5cf04449000000000000000000000000000000000000000000000000000000008152fd5b868567ffffffffffffffff888d8c878f868885928d88610aa261010073ffffffffffffffffffffffffffffffffffffffff600254169501516040519c8d977f01447eaa0000000000000000000000000000000000000000000000000000000089521660048801526060602488015260648701906131b3565b917ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc8684030160448701525191828152019190855b8a828210610f42575050505082809103915afa948515610f37578395610e77575b5015610d8e5750805b67ffffffffffffffff608087510191169052805b61010086015151811015610b4b5780610b3060019286613026565b516080610b42836101008b0151613026565b51015201610b15565b5083610b5686613478565b91604051848101907f130ac867e79e2789f923760a88743d292acdf7002139a588206e2260f73f7321825267ffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016604082015267ffffffffffffffff8416606082015230608082015260808152610bd660a0826129e6565b51902073ffffffffffffffffffffffffffffffffffffffff8585015116845167ffffffffffffffff6080816060840151169201511673ffffffffffffffffffffffffffffffffffffffff60a08801511660c088015191604051938a850195865260408501526060840152608083015260a082015260a08152610c5960c0826129e6565b51902060608501518681519101206040860151878151910120610100870151604051610cbf81610c938c8201948d865260408301906131b3565b037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081018352826129e6565b51902091608088015189815191012093604051958a870197885260408701526060860152608085015260a084015260c083015260e082015260e08152610d07610100826129e6565b51902082515267ffffffffffffffff60608351015116907f192442a2b2adb6a7948f097023cb6b57d29d3a7a5dd33e6666d33c39cc456f3267ffffffffffffffff60405192169180610d59868261329e565b0390a37fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff600254166002555151604051908152f35b73ffffffffffffffffffffffffffffffffffffffff604051917fea458c0c00000000000000000000000000000000000000000000000000000000835267ffffffffffffffff8716600484015216602482015282816044818573ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165af1908115610e6c578291610e33575b50610b01565b90508281813d8311610e65575b610e4a81836129e6565b81010312610e6157610e5b90613289565b86610e2d565b5080fd5b503d610e40565b6040513d84823e3d90fd5b9094503d8084833e610e8981836129e6565b8101908481830312610f2f5780519067ffffffffffffffff8211610f33570181601f82011215610f2f578051610ebe81612af8565b92610ecc60405194856129e6565b818452868085019260051b84010192818411610f2b57878101925b848410610efa5750505050509387610af8565b835167ffffffffffffffff8111610f27578991610f1c85848094870101613171565b815201930192610ee7565b8880fd5b8680fd5b8380fd5b8480fd5b6040513d85823e3d90fd5b8351805173ffffffffffffffffffffffffffffffffffffffff168652810151818601528a97508c965060409094019390920191600101610ad7565b94505050503d8089843e610f9181846129e6565b8201608083820312610f2757825192610fab868201612bff565b90604081015167ffffffffffffffff8111610ff45783610fcc918301613171565b92606082015167ffffffffffffffff81116109f757610feb9201613171565b939091386106b4565b8b80fd5b6040513d8a823e3d90fd5b60408336031261103a5786604091825161101c816129ca565b61102586612b63565b815282860135838201528152019201916105a9565b8a80fd5b8980fd5b60209061104d61310e565b82828a010152016104c7565b6024867f4e487b710000000000000000000000000000000000000000000000000000000081526011600452fd5b803b15610f33578460405180927fe0a0e5060000000000000000000000000000000000000000000000000000000082528183816110cc6024356004018b60048401612c9b565b03925af180156110ef571561038957846110e8919592956129e6565b9238610389565b6040513d87823e3d90fd5b6004847f1c0a3529000000000000000000000000000000000000000000000000000000008152fd5b73ffffffffffffffffffffffffffffffffffffffff8316600090815260028301602052604090205461034b5760248573ffffffffffffffffffffffffffffffffffffffff857fd0d2597600000000000000000000000000000000000000000000000000000000835216600452fd5b6004847fa4ec7479000000000000000000000000000000000000000000000000000000008152fd5b6004847f3ee5aeb5000000000000000000000000000000000000000000000000000000008152fd5b50346101525760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101525760043567ffffffffffffffff8111610e6157611230903690600401612b84565b73ffffffffffffffffffffffffffffffffffffffff600154163303611561575b919081907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8181360301915b8481101561155d578060051b82013583811215610f3357820191608083360312610f3357604051946112ac86612946565b6112b584612ae3565b86526112c360208501612b10565b9660208701978852604085013567ffffffffffffffff8111611559576112ec9036908701612fc1565b9460408801958652606081013567ffffffffffffffff8111610f2f5761131491369101612fc1565b60608801908152875167ffffffffffffffff1683526006602052604080842099518a547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff169015159182901b68ff000000000000000016178a55909590815151611431575b5095976001019550815b855180518210156113c257906113bb73ffffffffffffffffffffffffffffffffffffffff6113b383600195613026565b5116896138cd565b5001611383565b505095909694506001929193519081516113e2575b50500193929361127b565b61142767ffffffffffffffff7fc237ec1921f855ccd5e9a5af9733f2d58943a5a8501ec5988e305d7a4d42158692511692604051918291602083526020830190612bb5565b0390a238806113d7565b9893959296919094979860001461152257600184019591875b865180518210156114c7576114748273ffffffffffffffffffffffffffffffffffffffff92613026565b5116801561149057906114896001928a61383c565b500161144a565b60248a67ffffffffffffffff8e51167f463258ff000000000000000000000000000000000000000000000000000000008252600452fd5b50509692955090929796937f330939f6eafe8bb516716892fe962ff19770570838686e6579dbc1cc51fc328161151867ffffffffffffffff8a51169251604051918291602083526020830190612bb5565b0390a23880611379565b60248767ffffffffffffffff8b51167f463258ff000000000000000000000000000000000000000000000000000000008252600452fd5b8280fd5b8380f35b73ffffffffffffffffffffffffffffffffffffffff60055416330315611250576004837f905d7d9b000000000000000000000000000000000000000000000000000000008152fd5b50346101525760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101525767ffffffffffffffff6115ea612ac7565b16808252600660205260ff604083205460401c16908252600660205260016040832001916040518093849160208254918281520191845260208420935b81811061165e57505061163c925003836129e6565b61165a60405192839215158352604060208401526040830190612bb5565b0390f35b8454835260019485019487945060209093019201611627565b50346101525760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101525767ffffffffffffffff6116b8612ac7565b1681526006602052600167ffffffffffffffff604083205416019067ffffffffffffffff82116116f75760208267ffffffffffffffff60405191168152f35b807f4e487b7100000000000000000000000000000000000000000000000000000000602492526011600452fd5b503461015257807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261015257602073ffffffffffffffffffffffffffffffffffffffff60015416604051908152f35b503461015257807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261015257805473ffffffffffffffffffffffffffffffffffffffff81163303611833577fffffffffffffffffffffffff000000000000000000000000000000000000000060015491338284161760015516825573ffffffffffffffffffffffffffffffffffffffff3391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08380a380f35b6004827f02b543c6000000000000000000000000000000000000000000000000000000008152fd5b503461015257807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261015257611892612f96565b5060a06040516118a181612991565b60ff60025473ffffffffffffffffffffffffffffffffffffffff81168352831c161515602082015273ffffffffffffffffffffffffffffffffffffffff60035416604082015273ffffffffffffffffffffffffffffffffffffffff60045416606082015273ffffffffffffffffffffffffffffffffffffffff600554166080820152611976604051809273ffffffffffffffffffffffffffffffffffffffff60808092828151168552602081015115156020860152826040820151166040860152826060820151166060860152015116910152565bf35b50346101525760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261015257604060609167ffffffffffffffff6119be612ac7565b1681526006602052205473ffffffffffffffffffffffffffffffffffffffff6040519167ffffffffffffffff8116835260ff8160401c161515602084015260481c166040820152f35b50346101525760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101525773ffffffffffffffffffffffffffffffffffffffff611a54612b1d565b611a5c6133fd565b168015611a8f577fffffffffffffffffffffffff0000000000000000000000000000000000000000600754161760075580f35b6004827f8579befe000000000000000000000000000000000000000000000000000000008152fd5b50346101525760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101525760043567ffffffffffffffff8111610e6157611b07903690600401612b84565b9073ffffffffffffffffffffffffffffffffffffffff6004541690835b83811015611d705773ffffffffffffffffffffffffffffffffffffffff611b4f8260051b8401612f75565b1690604051917f70a08231000000000000000000000000000000000000000000000000000000008352306004840152602083602481845afa928315611d65578793611d32575b5082611ba7575b506001915001611b24565b8460405193611c4360208601957fa9059cbb00000000000000000000000000000000000000000000000000000000875283602482015282604482015260448152611bf26064826129e6565b8a80604098895193611c048b866129e6565b602085527f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65646020860152519082895af1611c3c613448565b9086613a5f565b805180611c7f575b505060207f508d7d183612c18fc339b42618912b9fa3239f631dd7ec0671f950200a0fa66e9160019651908152a338611b9c565b819294959693509060209181010312610f27576020611c9e9101612bff565b15611caf5792919085903880611c4b565b608490517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f742073756363656564000000000000000000000000000000000000000000006064820152fd5b9092506020813d8211611d5d575b81611d4d602093836129e6565b81010312610f2b57519138611b95565b3d9150611d40565b6040513d89823e3d90fd5b8480f35b50346101525760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261015257611dac612ac7565b506024359073ffffffffffffffffffffffffffffffffffffffff82168203610152576020611dd983612e72565b73ffffffffffffffffffffffffffffffffffffffff60405191168152f35b50346101525760a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261015257604051611e3381612991565b611e3b612b1d565b81526024358015158103611559576020820190815260443573ffffffffffffffffffffffffffffffffffffffff81168103610f2f5760408301908152611e7f612b40565b90606084019182526084359273ffffffffffffffffffffffffffffffffffffffff841684036121fd5760808501938452611eb76133fd565b73ffffffffffffffffffffffffffffffffffffffff8551161580156121de575b80156121d4575b6121ac579273ffffffffffffffffffffffffffffffffffffffff859381809461012097827fc7372d2d886367d7bb1b0e0708a5436f2c91d6963de210eb2dc1ec2ecd6d21f19a51167fffffffffffffffffffffffff000000000000000000000000000000000000000060025416176002555115157fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff74ff00000000000000000000000000000000000000006002549260a01b1691161760025551167fffffffffffffffffffffffff0000000000000000000000000000000000000000600354161760035551167fffffffffffffffffffffffff0000000000000000000000000000000000000000600454161760045551167fffffffffffffffffffffffff000000000000000000000000000000000000000060055416176005556121a86040519161202883612946565b67ffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016835273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016602084015273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016604084015273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000166060840152612158604051809473ffffffffffffffffffffffffffffffffffffffff6060809267ffffffffffffffff8151168552826020820151166020860152826040820151166040860152015116910152565b608083019073ffffffffffffffffffffffffffffffffffffffff60808092828151168552602081015115156020860152826040820151166040860152826060820151166060860152015116910152565ba180f35b6004867f35be3ac8000000000000000000000000000000000000000000000000000000008152fd5b5080511515611ede565b5073ffffffffffffffffffffffffffffffffffffffff83511615611ed7565b8580fd5b50346101525760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610152576004359067ffffffffffffffff8211610152573660238301121561015257816004013561225d81612af8565b9261226b60405194856129e6565b8184526024606060208601930282010190368211610f2f57602401915b8183106124095750505061229a6133fd565b805b8251811015612405576122af8184613026565b5167ffffffffffffffff6122c38386613026565b5151169081156123d957907fd5ad72bc37dc7a80a8b9b9df20500046fd7341adb1be2258a540466fdd7dcef5606060019493838752600660205260ff6040882061239a604060208501519483547fffffff0000000000000000000000000000000000000000ffffffffffffffffff7cffffffffffffffffffffffffffffffffffffffff0000000000000000008860481b1691161784550151151582907fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff68ff0000000000000000835492151560401b169116179055565b5473ffffffffffffffffffffffffffffffffffffffff6040519367ffffffffffffffff8316855216602084015260401c1615156040820152a20161229c565b602484837fc35aa79d000000000000000000000000000000000000000000000000000000008252600452fd5b5080f35b606083360312610f2f576040516060810181811067ffffffffffffffff8211176124845760405261243984612ae3565b8152602084013573ffffffffffffffffffffffffffffffffffffffff811681036121fd57918160609360208094015261247460408701612b10565b6040820152815201920191612288565b6024867f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b50346101525760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610152576124e9612ac7565b60243567ffffffffffffffff81116115595760a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc8236030112611559576040517f2cbc26bb00000000000000000000000000000000000000000000000000000000815277ffffffffffffffff000000000000000000000000000000008360801b16600482015260208160248173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165afa9081156126df5784916126a5575b5061266f5761261c9160209173ffffffffffffffffffffffffffffffffffffffff60025416906040518095819482937fd8694ccd0000000000000000000000000000000000000000000000000000000084526004019060048401612c9b565b03915afa908115610e6c578291612639575b602082604051908152f35b90506020813d602011612667575b81612654602093836129e6565b81010312610e615760209150513861262e565b3d9150612647565b60248367ffffffffffffffff847ffdbd6a7200000000000000000000000000000000000000000000000000000000835216600452fd5b90506020813d6020116126d7575b816126c0602093836129e6565b81010312610f2f576126d190612bff565b386125bd565b3d91506126b3565b6040513d86823e3d90fd5b503461015257807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610152575061165a60405161272b6040826129e6565b601081527f4f6e52616d7020312e362e302d646576000000000000000000000000000000006020820152604051918291602083526020830190612a84565b503461015257807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261015257602073ffffffffffffffffffffffffffffffffffffffff60075416604051908152f35b905034610e6157817ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610e6157806127f7606092612946565b8281528260208201528260408201520152608060405161281681612946565b67ffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016815273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016602082015273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016604082015273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000166060820152611976604051809273ffffffffffffffffffffffffffffffffffffffff6060809267ffffffffffffffff8151168552826020820151166020860152826040820151166040860152015116910152565b6080810190811067ffffffffffffffff82111761296257604052565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60a0810190811067ffffffffffffffff82111761296257604052565b610120810190811067ffffffffffffffff82111761296257604052565b6040810190811067ffffffffffffffff82111761296257604052565b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff82111761296257604052565b67ffffffffffffffff811161296257601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b60005b838110612a745750506000910152565b8181015183820152602001612a64565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f602093612ac081518092818752878088019101612a61565b0116010190565b6004359067ffffffffffffffff82168203612ade57565b600080fd5b359067ffffffffffffffff82168203612ade57565b67ffffffffffffffff81116129625760051b60200190565b35908115158203612ade57565b6004359073ffffffffffffffffffffffffffffffffffffffff82168203612ade57565b6064359073ffffffffffffffffffffffffffffffffffffffff82168203612ade57565b359073ffffffffffffffffffffffffffffffffffffffff82168203612ade57565b9181601f84011215612ade5782359167ffffffffffffffff8311612ade576020808501948460051b010111612ade57565b906020808351928381520192019060005b818110612bd35750505090565b825173ffffffffffffffffffffffffffffffffffffffff16845260209384019390920191600101612bc6565b51908115158203612ade57565b90357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe182360301811215612ade57016020813591019167ffffffffffffffff8211612ade578136038313612ade57565b601f82602094937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0938186528686013760008582860101520116010190565b9067ffffffffffffffff9093929316815260406020820152612d11612cd4612cc38580612c0c565b60a0604086015260e0850191612c5c565b612ce16020860186612c0c565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0858403016060860152612c5c565b9060408401357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe185360301811215612ade5784016020813591019267ffffffffffffffff8211612ade578160061b36038413612ade578281037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0016080840152818152602001929060005b818110612e1257505050612ddf8473ffffffffffffffffffffffffffffffffffffffff612dcf6060612e0f979801612b63565b1660a08401526080810190612c0c565b9160c07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc082860301910152612c5c565b90565b90919360408060019273ffffffffffffffffffffffffffffffffffffffff612e3989612b63565b16815260208881013590820152019501929101612d9c565b519073ffffffffffffffffffffffffffffffffffffffff82168203612ade57565b73ffffffffffffffffffffffffffffffffffffffff604051917fbbe4f6db00000000000000000000000000000000000000000000000000000000835216600482015260208160248173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165afa8015612f6957600090612f1c575b73ffffffffffffffffffffffffffffffffffffffff91501690565b506020813d602011612f61575b81612f36602093836129e6565b81010312612ade57612f5c73ffffffffffffffffffffffffffffffffffffffff91612e51565b612f01565b3d9150612f29565b6040513d6000823e3d90fd5b3573ffffffffffffffffffffffffffffffffffffffff81168103612ade5790565b60405190612fa382612991565b60006080838281528260208201528260408201528260608201520152565b9080601f83011215612ade578135612fd881612af8565b92612fe660405194856129e6565b81845260208085019260051b820101928311612ade57602001905b82821061300e5750505090565b6020809161301b84612b63565b815201910190613001565b805182101561303a5760209160051b010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe181360301821215612ade570180359067ffffffffffffffff8211612ade57602001918136038313612ade57565b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe181360301821215612ade570180359067ffffffffffffffff8211612ade57602001918160061b36038313612ade57565b6040519061311b82612991565b6060608083600081528260208201528260408201526000838201520152565b92919261314682612a27565b9161315460405193846129e6565b829481845281830111612ade578281602093846000960137010152565b81601f82011215612ade57805161318781612a27565b9261319560405194856129e6565b81845260208284010111612ade57612e0f9160208085019101612a61565b9080602083519182815201916020808360051b8301019401926000915b8383106131df57505050505090565b909192939460208061327a837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0866001960301875289519073ffffffffffffffffffffffffffffffffffffffff8251168152608061325f61324d8685015160a08886015260a0850190612a84565b60408501518482036040860152612a84565b92606081015160608401520151906080818403910152612a84565b970193019301919392906131d0565b519067ffffffffffffffff82168203612ade57565b90612e0f916020815267ffffffffffffffff6080835180516020850152826020820151166040850152826040820151166060850152826060820151168285015201511660a082015273ffffffffffffffffffffffffffffffffffffffff60208301511660c082015261010061339261335d61332a60408601516101a060e08701526101c0860190612a84565b60608601517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08683030185870152612a84565b60808501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe085830301610120860152612a84565b9273ffffffffffffffffffffffffffffffffffffffff60a08201511661014084015260c081015161016084015260e08101516101808401520151906101a07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0828503019101526131b3565b73ffffffffffffffffffffffffffffffffffffffff60015416330361341e57565b7f2b5c74de0000000000000000000000000000000000000000000000000000000060005260046000fd5b3d15613473573d9061345982612a27565b9161346760405193846129e6565b82523d6000602084013e565b606090565b6000613519819260405161348b816129ad565b613493612f96565b815283602082015260606040820152606080820152606060808201528360a08201528360c08201528360e082015260606101008201525073ffffffffffffffffffffffffffffffffffffffff60075416906040519485809481937f8a06fadb0000000000000000000000000000000000000000000000000000000083526004830161329e565b03925af18091600091613576575b5090612e0f57613572613538613448565b6040519182917f828ebdfb000000000000000000000000000000000000000000000000000000008352602060048401526024830190612a84565b0390fd5b3d8083833e61358581836129e6565b8101906020818303126115595780519067ffffffffffffffff8211610f2f570191828203926101a08412610e615760a0604051946135c2866129ad565b12610e61576040516135d381612991565b815181526135e360208301613289565b60208201526135f460408301613289565b604082015261360560608301613289565b606082015261361660808301613289565b6080820152845261362960a08201612e51565b602085015260c081015167ffffffffffffffff8111611559578361364e918301613171565b604085015260e081015167ffffffffffffffff81116115595783613673918301613171565b606085015261010081015167ffffffffffffffff81116115595783613699918301613171565b60808501526136ab6101208201612e51565b60a085015261014081015160c085015261016081015160e08501526101808101519067ffffffffffffffff8211611559570182601f82011215610e61578051916136f483612af8565b9361370260405195866129e6565b83855260208086019460051b840101928184116115595760208101945b8486106137385750505050505061010082015238613527565b855167ffffffffffffffff8111610f3357820160a07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08286030112610f33576040519061378482612991565b61379060208201612e51565b8252604081015167ffffffffffffffff8111610f2b578560206137b592840101613171565b6020830152606081015167ffffffffffffffff8111610f2b578560206137dd92840101613171565b60408301526080810151606083015260a08101519067ffffffffffffffff8211610f2b579161381486602080969481960101613171565b608082015281520195019461371f565b805482101561303a5760005260206000200190600090565b60008281526001820160205260409020546138c6578054906801000000000000000082101561296257826138af61387a846001809601855584613824565b81939154907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9060031b92831b921b19161790565b905580549260005201602052604060002055600190565b5050600090565b9060018201918160005282602052604060002054801515600014613a56577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8101818111613a27578254907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8201918211613a27578181036139f0575b505050805480156139c1577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01906139828282613824565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82549160031b1b191690555560005260205260006040812055600190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b613a10613a0061387a9386613824565b90549060031b1c92839286613824565b90556000528360205260406000205538808061394a565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b50505050600090565b91929015613ada5750815115613a73575090565b3b15613a7c5790565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152fd5b825190915015613aed5750805190602001fd5b613572906040519182917f08c379a0000000000000000000000000000000000000000000000000000000008352602060048401526024830190612a8456fea164736f6c634300081a000a",
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
