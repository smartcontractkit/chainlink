// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package rmn_home

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

type RMNHomeDynamicConfig struct {
	SourceChains   []RMNHomeSourceChain
	OffchainConfig []byte
}

type RMNHomeNode struct {
	PeerId            [32]byte
	OffchainPublicKey [32]byte
}

type RMNHomeSourceChain struct {
	ChainSelector       uint64
	FObserve            uint64
	ObserverNodesBitmap *big.Int
}

type RMNHomeStaticConfig struct {
	Nodes          []RMNHomeNode
	OffchainConfig []byte
}

type RMNHomeVersionedConfig struct {
	Version       uint32
	ConfigDigest  [32]byte
	StaticConfig  RMNHomeStaticConfig
	DynamicConfig RMNHomeDynamicConfig
}

var RMNHomeMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getActiveDigest\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllConfigs\",\"inputs\":[],\"outputs\":[{\"name\":\"activeConfig\",\"type\":\"tuple\",\"internalType\":\"structRMNHome.VersionedConfig\",\"components\":[{\"name\":\"version\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"configDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"staticConfig\",\"type\":\"tuple\",\"internalType\":\"structRMNHome.StaticConfig\",\"components\":[{\"name\":\"nodes\",\"type\":\"tuple[]\",\"internalType\":\"structRMNHome.Node[]\",\"components\":[{\"name\":\"peerId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"offchainPublicKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"dynamicConfig\",\"type\":\"tuple\",\"internalType\":\"structRMNHome.DynamicConfig\",\"components\":[{\"name\":\"sourceChains\",\"type\":\"tuple[]\",\"internalType\":\"structRMNHome.SourceChain[]\",\"components\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"fObserve\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"observerNodesBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"candidateConfig\",\"type\":\"tuple\",\"internalType\":\"structRMNHome.VersionedConfig\",\"components\":[{\"name\":\"version\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"configDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"staticConfig\",\"type\":\"tuple\",\"internalType\":\"structRMNHome.StaticConfig\",\"components\":[{\"name\":\"nodes\",\"type\":\"tuple[]\",\"internalType\":\"structRMNHome.Node[]\",\"components\":[{\"name\":\"peerId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"offchainPublicKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"dynamicConfig\",\"type\":\"tuple\",\"internalType\":\"structRMNHome.DynamicConfig\",\"components\":[{\"name\":\"sourceChains\",\"type\":\"tuple[]\",\"internalType\":\"structRMNHome.SourceChain[]\",\"components\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"fObserve\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"observerNodesBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCandidateDigest\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getConfig\",\"inputs\":[{\"name\":\"configDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"versionedConfig\",\"type\":\"tuple\",\"internalType\":\"structRMNHome.VersionedConfig\",\"components\":[{\"name\":\"version\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"configDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"staticConfig\",\"type\":\"tuple\",\"internalType\":\"structRMNHome.StaticConfig\",\"components\":[{\"name\":\"nodes\",\"type\":\"tuple[]\",\"internalType\":\"structRMNHome.Node[]\",\"components\":[{\"name\":\"peerId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"offchainPublicKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"dynamicConfig\",\"type\":\"tuple\",\"internalType\":\"structRMNHome.DynamicConfig\",\"components\":[{\"name\":\"sourceChains\",\"type\":\"tuple[]\",\"internalType\":\"structRMNHome.SourceChain[]\",\"components\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"fObserve\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"observerNodesBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"ok\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getConfigDigests\",\"inputs\":[],\"outputs\":[{\"name\":\"activeConfigDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"candidateConfigDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"promoteCandidateAndRevokeActive\",\"inputs\":[{\"name\":\"digestToPromote\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"digestToRevoke\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeCandidate\",\"inputs\":[{\"name\":\"configDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setCandidate\",\"inputs\":[{\"name\":\"staticConfig\",\"type\":\"tuple\",\"internalType\":\"structRMNHome.StaticConfig\",\"components\":[{\"name\":\"nodes\",\"type\":\"tuple[]\",\"internalType\":\"structRMNHome.Node[]\",\"components\":[{\"name\":\"peerId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"offchainPublicKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"dynamicConfig\",\"type\":\"tuple\",\"internalType\":\"structRMNHome.DynamicConfig\",\"components\":[{\"name\":\"sourceChains\",\"type\":\"tuple[]\",\"internalType\":\"structRMNHome.SourceChain[]\",\"components\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"fObserve\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"observerNodesBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"digestToOverwrite\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"newConfigDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setDynamicConfig\",\"inputs\":[{\"name\":\"newDynamicConfig\",\"type\":\"tuple\",\"internalType\":\"structRMNHome.DynamicConfig\",\"components\":[{\"name\":\"sourceChains\",\"type\":\"tuple[]\",\"internalType\":\"structRMNHome.SourceChain[]\",\"components\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"fObserve\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"observerNodesBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"currentDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"typeAndVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"ActiveConfigRevoked\",\"inputs\":[{\"name\":\"configDigest\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CandidateConfigRevoked\",\"inputs\":[{\"name\":\"configDigest\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ConfigPromoted\",\"inputs\":[{\"name\":\"configDigest\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ConfigSet\",\"inputs\":[{\"name\":\"configDigest\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"version\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"staticConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRMNHome.StaticConfig\",\"components\":[{\"name\":\"nodes\",\"type\":\"tuple[]\",\"internalType\":\"structRMNHome.Node[]\",\"components\":[{\"name\":\"peerId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"offchainPublicKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"dynamicConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRMNHome.DynamicConfig\",\"components\":[{\"name\":\"sourceChains\",\"type\":\"tuple[]\",\"internalType\":\"structRMNHome.SourceChain[]\",\"components\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"fObserve\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"observerNodesBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DynamicConfigSet\",\"inputs\":[{\"name\":\"configDigest\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"dynamicConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRMNHome.DynamicConfig\",\"components\":[{\"name\":\"sourceChains\",\"type\":\"tuple[]\",\"internalType\":\"structRMNHome.SourceChain[]\",\"components\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"fObserve\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"observerNodesBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"CannotTransferToSelf\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ConfigDigestMismatch\",\"inputs\":[{\"name\":\"expectedConfigDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"gotConfigDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"DigestNotFound\",\"inputs\":[{\"name\":\"configDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"DuplicateOffchainPublicKey\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DuplicatePeerId\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DuplicateSourceChain\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MustBeProposedOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoOpStateTransitionNotAllowed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotEnoughObservers\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OutOfBoundsNodesLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OutOfBoundsObserverNodeIndex\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnerCannotBeZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RevokingZeroDigestNotAllowed\",\"inputs\":[]}]",
	Bin: "0x60808060405234604d573315603c57600180546001600160a01b03191633179055600e80546001600160401b031916905561262690816100538239f35b639b15e16f60e01b60005260046000fd5b600080fdfe6080604052600436101561001257600080fd5b60003560e01c8063118dbac514610904578063123e65db146108ae578063181f5a77146108315780633567e6b4146107b757806338354c5c1461077657806363507956146106985780636dd5b69d1461063b578063736be802146105b457806379ba5097146104cb5780638c76967f146103255780638da5cb5b146102d3578063f2fde38b146101e05763fb4022d4146100ab57600080fd5b346101db5760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101db576004356100e56123b4565b80156101b15763ffffffff610105600163ffffffff600e5460201c161890565b1660028110159081610152576006600391020191825481810361018157507f53f5d9228f0a4173bea6e5931c9b3afe6eeb6692ede1d182952970f152534e3b600080a26101525760009055005b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f93df584c0000000000000000000000000000000000000000000000000000000060005260045260245260446000fd5b7f0849d8cc0000000000000000000000000000000000000000000000000000000060005260046000fd5b600080fd5b346101db5760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101db5760043573ffffffffffffffffffffffffffffffffffffffff81168091036101db576102386123b4565b3381146102a957807fffffffffffffffffffffffff0000000000000000000000000000000000000000600054161760005573ffffffffffffffffffffffffffffffffffffffff600154167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278600080a3005b7fdad89dca0000000000000000000000000000000000000000000000000000000060005260046000fd5b346101db5760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101db57602073ffffffffffffffffffffffffffffffffffffffff60015416604051908152f35b346101db5760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101db576004356024356103626123b4565b8115806104c3575b6104995763ffffffff610388600163ffffffff600e5460201c161890565b1660028110156101525760060260030154828103610467575063ffffffff600e5460201c166002811015610152576006026003018054828103610467575060009055600e547fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff67ffffffff00000000600163ffffffff8460201c161860201b16911617600e558061043c575b507ffc3e98dbbd47c3fa7c1c05b6ec711caeaf70eca4554192b9ada8fc11a37f298e600080a2005b7f0b31c0055e2d464bef7781994b98c4ff9ef4ae0d05f59feb6a68c42de5e201b8600080a281610414565b90507f93df584c0000000000000000000000000000000000000000000000000000000060005260045260245260446000fd5b7f7b4d1e4f0000000000000000000000000000000000000000000000000000000060005260046000fd5b50801561036a565b346101db5760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101db5760005473ffffffffffffffffffffffffffffffffffffffff8116330361058a577fffffffffffffffffffffffff00000000000000000000000000000000000000006001549133828416176001551660005573ffffffffffffffffffffffffffffffffffffffff3391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0600080a3005b7f02b543c60000000000000000000000000000000000000000000000000000000060005260046000fd5b346101db5760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101db5760043567ffffffffffffffff81116101db5760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc82360301126101db576106399061062d6123b4565b60243590600401611f4a565b005b346101db5760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101db5761068c610678600435611eea565b6040519283926040845260408401906116d2565b90151560208301520390f35b346101db5760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101db576106cf611ca0565b6106d7611ca0565b9063ffffffff600e5460201c1660028110156101525760066106fc9102600201611d92565b602081015161076e575b50600e5460201c63ffffffff166001186002811015610152576107549261073560066107629302600201611d92565b6020810151610766575b506040519384936040855260408501906116d2565b9083820360208501526116d2565b0390f35b90508461073f565b905082610706565b346101db5760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101db5760206107af611c65565b604051908152f35b346101db5760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101db5763ffffffff600e5460201c1660028110156101525760060260030154600e5460201c63ffffffff16600118906002821015610152576003600660409302015482519182526020820152f35b346101db5760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101db576107626040805190610872818361162b565b600d82527f524d4e486f6d6520312e362e300000000000000000000000000000000000000060208301525191829160208352602083019061168f565b346101db5760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101db5763ffffffff600e5460201c1660028110156101525760036006602092020154604051908152f35b346101db5760607ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101db5760043567ffffffffffffffff81116101db5760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc82360301126101db5767ffffffffffffffff602435116101db57602435360360407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc8201126101db576109b96123b4565b6040516000916109c8826115a8565b67ffffffffffffffff8460040135116115a4578360040135840190366023830112156115a0576109fb600483013561180a565b610a08604051918261162b565b6004830135808252602082019060061b84016024013681116112ad5760248501915b81831061156f575050508352602485013567ffffffffffffffff811161156b57610a5a9060043691880101611822565b6020840152610a6e366024356004016118ac565b94610100845151116115435790919484955b845151871015610b765760018701808811610b49575b85518051821015610b3c5788610aab916123ff565b5151610ab88288516123ff565b515114610b14576020610acc8988516123ff565b5101516020610adc8389516123ff565b51015114610aec57600101610a96565b6004877fae00651d000000000000000000000000000000000000000000000000000000008152fd5b6004877f221a8ae8000000000000000000000000000000000000000000000000000000008152fd5b5050600190960195610a80565b6024877f4e487b710000000000000000000000000000000000000000000000000000000081526011600452fd5b949390610b8591515190612413565b610b8d611c65565b6044358103611513576114e7575b600e549363ffffffff85169463ffffffff86146114ba5763ffffffff60017fffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000095969701169384911617600e557dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff604051610ccc610cfa6020830160208152610c5784610c2b604082018a600401611a26565b037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0810186528561162b565b602060405191818301957f45564d000000000000000000000000000000000000000000000000000000000087524660408501523060608501528a608085015260808452610ca560a08561162b565b604051958694610cbd858701998a925192839161166c565b8501915180938584019061166c565b0101037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0810183528261162b565b519020167e0b0000000000000000000000000000000000000000000000000000000000001793610d35600163ffffffff600e5460201c161890565b600281101561148d576006029182600201866003850155857fffffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000825416179055600483017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdd85360301856004013512156112a95767ffffffffffffffff6004830135116112a957600482013560061b360360248301136112a9576801000000000000000060048301351161127c578054600483013580835581116113ff575b508752602087208790602483015b600484013583106113db575050505060058201610e246024850185600401611ad9565b9067ffffffffffffffff82116113ae578190610e408454611b2a565b601f811161135e575b508990601f83116001146112bc578a926112b1575b50507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790555b60068201907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdd6024356004013591018112156112ad57602435019060048201359167ffffffffffffffff83116112a95760240160608302360381136112a95768010000000000000000831161127c5781548383558084106111ca575b509087526020872087915b838310611150575050505060070191610f396024803501602435600401611ad9565b67ffffffffffffffff819592951161112357610f558254611b2a565b601f81116110de575b5090859392918760209890601f831160011461101a57907ff6c6d1be15ba0acc8ee645c1ec613c360ef786d2d3200eb8e695b6dec757dbf096978361100f575b50507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790555b611004610ff160405193849384526060898501526060840190600401611a26565b8281036040840152602435600401611b92565b0390a2604051908152f35b013590508980610f9e565b96907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe083168489528a8920985b8181106110c45750917ff6c6d1be15ba0acc8ee645c1ec613c360ef786d2d3200eb8e695b6dec757dbf097989184600195941061108c575b505050811b019055610fd0565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88560031b161c1991013516905589808061107f565b828401358a556001909901988a9850918b01918b01611047565b82885260208820601f830160051c81019160208410611119575b601f0160051c01905b81811061110e5750610f5e565b888155600101611101565b90915081906110f8565b6024877f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b600260608267ffffffffffffffff611169600195611b7d565b168554907fffffffffffffffffffffffffffffffff000000000000000000000000000000006fffffffffffffffff00000000000000006111ab60208601611b7d565b60401b1692161717855560408101358486015501920192019190610f17565b7f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8116810361124f577f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8416840361124f57828952602089209060011b8101908460011b015b81811061123d5750610f0c565b808a600292558a600182015501611230565b6024897f4e487b710000000000000000000000000000000000000000000000000000000081526011600452fd5b6024887f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b8780fd5b8680fd5b013590508980610e5e565b90917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01691848b5260208b20928b5b818110611346575090846001959493921061130e575b505050811b019055610e90565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88560031b161c19910135169055898080611301565b919360206001819287870135815501950192016112eb565b909150838a5260208a20601f840160051c810191602085106113a4575b90601f859493920160051c01905b8181106113965750610e49565b8b8155849350600101611389565b909150819061137b565b6024897f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b60016002604083600494358655602081013584870155019301930192919050610e01565b7f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8116810361124f5760048301357f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8116900361124f57818952602089209060011b810190600484013560011b015b81811061147b5750610df3565b6002908a81558a60018201550161146e565b6024877f4e487b710000000000000000000000000000000000000000000000000000000081526032600452fd5b6024857f4e487b710000000000000000000000000000000000000000000000000000000081526011600452fd5b6044357f53f5d9228f0a4173bea6e5931c9b3afe6eeb6692ede1d182952970f152534e3b8480a2610b9b565b7f93df584c0000000000000000000000000000000000000000000000000000000084526004526044803560245283fd5b6004857faf26d5e3000000000000000000000000000000000000000000000000000000008152fd5b8480fd5b6040833603126112a95760206040918251611589816115a8565b853581528286013583820152815201920191610a2a565b8380fd5b8280fd5b6040810190811067ffffffffffffffff8211176115c457604052565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6060810190811067ffffffffffffffff8211176115c457604052565b6080810190811067ffffffffffffffff8211176115c457604052565b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff8211176115c457604052565b60005b83811061167f5750506000910152565b818101518382015260200161166f565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f6020936116cb8151809281875287808801910161166c565b0116010190565b91909163ffffffff81511683526020810151602084015260408101516080604085015260c0840190805191604060808701528251809152602060e0870193019060005b8181106117e85750505060609160206117599201517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff808783030160a088015261168f565b910151926060818303910152604081019083519160408252825180915260206060830193019060005b8181106117a65750505060206117a39394015190602081840391015261168f565b90565b909193602060606001926040885167ffffffffffffffff815116835267ffffffffffffffff858201511685840152015160408201520195019101919091611782565b8251805186526020908101518187015260409095019490920191600101611715565b67ffffffffffffffff81116115c45760051b60200190565b81601f820112156101db5780359067ffffffffffffffff82116115c45760405192611875601f84017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0166020018561162b565b828452602083830101116101db57816000926020809301838601378301015290565b359067ffffffffffffffff821682036101db57565b91906040838203126101db57604051906118c5826115a8565b8193803567ffffffffffffffff81116101db57810182601f820112156101db5780356118f08161180a565b916118fe604051938461162b565b818352602060608185019302820101908582116101db57602001915b81831061194d57505050835260208101359167ffffffffffffffff83116101db576020926119489201611822565b910152565b6060838703126101db576020606091604051611968816115f3565b61197186611897565b815261197e838701611897565b838201526040860135604082015281520192019161191a565b90357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1823603018112156101db57016020813591019167ffffffffffffffff82116101db5781360383136101db57565b601f82602094937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0938186528686013760008582860101520116010190565b9190604081019083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1853603018112156101db5784016020813591019267ffffffffffffffff82116101db578160061b360384136101db578190604084525260608201929060005b818110611ab957505050611aab8460206117a395960190611997565b9160208185039101526119e7565b823585526020808401359086015260409485019490920191600101611a8f565b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1813603018212156101db570180359067ffffffffffffffff82116101db576020019181360383136101db57565b90600182811c92168015611b73575b6020831014611b4457565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b91607f1691611b39565b3567ffffffffffffffff811681036101db5790565b9190604081019083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1853603018112156101db5784016020813591019267ffffffffffffffff82116101db5760608202360384136101db578190604084525260608201929060005b818110611c1757505050611aab8460206117a395960190611997565b90919360608060019267ffffffffffffffff611c3289611897565b16815267ffffffffffffffff611c4a60208a01611897565b16602082015260408881013590820152019501929101611bfb565b600e5460201c63ffffffff166001186002811015610152576006026003015490565b60405190611c94826115a8565b60606020838281520152565b60405190611cad8261160f565b816000815260006020820152611cc1611c87565b60408201526060611948611c87565b9060405191826000825492611ce484611b2a565b8084529360018116908115611d525750600114611d0b575b50611d099250038361162b565b565b90506000929192526020600020906000915b818310611d36575050906020611d099282010138611cfc565b6020919350806001915483858901015201910190918492611d1d565b60209350611d099592507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0091501682840152151560051b82010138611cfc565b9060405191611da08361160f565b8263ffffffff8254168152600182015460208201526002820160405190611dc6826115a8565b8054611dd18161180a565b91611ddf604051938461162b565b818352602083019060005260206000206000915b838310611ebd57505050508152611e0c60038401611cd0565b60208201526040820152600482019160405192611e28846115a8565b8054611e338161180a565b91611e41604051938461162b565b818352602083019060005260206000206000915b838310611e7b5750505050600560609392611e7292865201611cd0565b60208401520152565b60026020600192604051611e8e816115f3565b67ffffffffffffffff8654818116835260401c1683820152848601546040820152815201920192019190611e55565b60026020600192604051611ed0816115a8565b855481528486015483820152815201920192019190611df3565b611ef2611ca0565b9060005b6002811015611f4257600060068202908360038301541480611f39575b611f21575050600101611ef6565b91509150611f33925050600201611d92565b90600190565b50831515611f13565b505090600090565b9060005b60028110156123865760006006820290836003830154148061237d575b611f79575050600101611f4e565b90915092919250611f976004820154611f9236856118ac565b612413565b6000906006810183357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1853603018112156115a05784019081359167ffffffffffffffff831161156b57602001606083023603811361156b5768010000000000000000831161235057815483835580841061229e575b509084526020842084915b8383106122245750505050600701906120346020840184611ad9565b919067ffffffffffffffff83116121f75761204f8454611b2a565b601f81116121b2575b5081601f84116001146120ec57926120dc949281927f1f69d1a2edb327babc986b3deb80091f101b9105d42a6c30db4d99c31d7e62949795926120e1575b50507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790555b604051918291602083526020830190611b92565b0390a2565b013590503880612096565b917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0841685845260208420935b81811061219a57509260019285927f1f69d1a2edb327babc986b3deb80091f101b9105d42a6c30db4d99c31d7e629498966120dc989610612162575b505050811b0190556120c8565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88560031b161c19910135169055388080612155565b91936020600181928787013581550195019201612119565b84835260208320601f850160051c810191602086106121ed575b601f0160051c01905b8181106121e25750612058565b8381556001016121d5565b90915081906121cc565b6024827f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b600260608267ffffffffffffffff61223d600195611b7d565b168554907fffffffffffffffffffffffffffffffff000000000000000000000000000000006fffffffffffffffff000000000000000061227f60208601611b7d565b60401b1692161717855560408101358486015501920192019190612018565b7f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81168103612323577f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8416840361232357828652602086209060011b8101908460011b015b818110612311575061200d565b80876002925587600182015501612304565b6024867f4e487b710000000000000000000000000000000000000000000000000000000081526011600452fd5b6024857f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b50831515611f6b565b507fd0b2c0310000000000000000000000000000000000000000000000000000000060005260045260246000fd5b73ffffffffffffffffffffffffffffffffffffffff6001541633036123d557565b7f2b5c74de0000000000000000000000000000000000000000000000000000000060005260046000fd5b80518210156101525760209160051b010190565b908151519160005b8381106124285750505050565b6124338183516123ff565b5160018201808311612508575b8581106125bf575060408101519084610100036101008111612508577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff83911c81160361259557816000925b61253757506020015160011b6801fffffffffffffffe67fffffffffffffffe8216911681036125085760010167ffffffffffffffff81116125085767ffffffffffffffff16116124de5760010161241b565b7fa804bcb30000000000000000000000000000000000000000000000000000000060005260046000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8101908082116125085716917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff811461250857600101918061248c565b7f2847b6060000000000000000000000000000000000000000000000000000000060005260046000fd5b67ffffffffffffffff82511667ffffffffffffffff6125df8387516123ff565b515116146125ef57600101612440565b7f3857f84d0000000000000000000000000000000000000000000000000000000060005260046000fdfea164736f6c634300081a000a",
}

var RMNHomeABI = RMNHomeMetaData.ABI

var RMNHomeBin = RMNHomeMetaData.Bin

func DeployRMNHome(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *RMNHome, error) {
	parsed, err := RMNHomeMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(RMNHomeBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &RMNHome{address: address, abi: *parsed, RMNHomeCaller: RMNHomeCaller{contract: contract}, RMNHomeTransactor: RMNHomeTransactor{contract: contract}, RMNHomeFilterer: RMNHomeFilterer{contract: contract}}, nil
}

type RMNHome struct {
	address common.Address
	abi     abi.ABI
	RMNHomeCaller
	RMNHomeTransactor
	RMNHomeFilterer
}

type RMNHomeCaller struct {
	contract *bind.BoundContract
}

type RMNHomeTransactor struct {
	contract *bind.BoundContract
}

type RMNHomeFilterer struct {
	contract *bind.BoundContract
}

type RMNHomeSession struct {
	Contract     *RMNHome
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type RMNHomeCallerSession struct {
	Contract *RMNHomeCaller
	CallOpts bind.CallOpts
}

type RMNHomeTransactorSession struct {
	Contract     *RMNHomeTransactor
	TransactOpts bind.TransactOpts
}

type RMNHomeRaw struct {
	Contract *RMNHome
}

type RMNHomeCallerRaw struct {
	Contract *RMNHomeCaller
}

type RMNHomeTransactorRaw struct {
	Contract *RMNHomeTransactor
}

func NewRMNHome(address common.Address, backend bind.ContractBackend) (*RMNHome, error) {
	abi, err := abi.JSON(strings.NewReader(RMNHomeABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindRMNHome(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RMNHome{address: address, abi: abi, RMNHomeCaller: RMNHomeCaller{contract: contract}, RMNHomeTransactor: RMNHomeTransactor{contract: contract}, RMNHomeFilterer: RMNHomeFilterer{contract: contract}}, nil
}

func NewRMNHomeCaller(address common.Address, caller bind.ContractCaller) (*RMNHomeCaller, error) {
	contract, err := bindRMNHome(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RMNHomeCaller{contract: contract}, nil
}

func NewRMNHomeTransactor(address common.Address, transactor bind.ContractTransactor) (*RMNHomeTransactor, error) {
	contract, err := bindRMNHome(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RMNHomeTransactor{contract: contract}, nil
}

func NewRMNHomeFilterer(address common.Address, filterer bind.ContractFilterer) (*RMNHomeFilterer, error) {
	contract, err := bindRMNHome(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RMNHomeFilterer{contract: contract}, nil
}

func bindRMNHome(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RMNHomeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_RMNHome *RMNHomeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RMNHome.Contract.RMNHomeCaller.contract.Call(opts, result, method, params...)
}

func (_RMNHome *RMNHomeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RMNHome.Contract.RMNHomeTransactor.contract.Transfer(opts)
}

func (_RMNHome *RMNHomeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RMNHome.Contract.RMNHomeTransactor.contract.Transact(opts, method, params...)
}

func (_RMNHome *RMNHomeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RMNHome.Contract.contract.Call(opts, result, method, params...)
}

func (_RMNHome *RMNHomeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RMNHome.Contract.contract.Transfer(opts)
}

func (_RMNHome *RMNHomeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RMNHome.Contract.contract.Transact(opts, method, params...)
}

func (_RMNHome *RMNHomeCaller) GetActiveDigest(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _RMNHome.contract.Call(opts, &out, "getActiveDigest")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_RMNHome *RMNHomeSession) GetActiveDigest() ([32]byte, error) {
	return _RMNHome.Contract.GetActiveDigest(&_RMNHome.CallOpts)
}

func (_RMNHome *RMNHomeCallerSession) GetActiveDigest() ([32]byte, error) {
	return _RMNHome.Contract.GetActiveDigest(&_RMNHome.CallOpts)
}

func (_RMNHome *RMNHomeCaller) GetAllConfigs(opts *bind.CallOpts) (GetAllConfigs,

	error) {
	var out []interface{}
	err := _RMNHome.contract.Call(opts, &out, "getAllConfigs")

	outstruct := new(GetAllConfigs)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ActiveConfig = *abi.ConvertType(out[0], new(RMNHomeVersionedConfig)).(*RMNHomeVersionedConfig)
	outstruct.CandidateConfig = *abi.ConvertType(out[1], new(RMNHomeVersionedConfig)).(*RMNHomeVersionedConfig)

	return *outstruct, err

}

func (_RMNHome *RMNHomeSession) GetAllConfigs() (GetAllConfigs,

	error) {
	return _RMNHome.Contract.GetAllConfigs(&_RMNHome.CallOpts)
}

func (_RMNHome *RMNHomeCallerSession) GetAllConfigs() (GetAllConfigs,

	error) {
	return _RMNHome.Contract.GetAllConfigs(&_RMNHome.CallOpts)
}

func (_RMNHome *RMNHomeCaller) GetCandidateDigest(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _RMNHome.contract.Call(opts, &out, "getCandidateDigest")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_RMNHome *RMNHomeSession) GetCandidateDigest() ([32]byte, error) {
	return _RMNHome.Contract.GetCandidateDigest(&_RMNHome.CallOpts)
}

func (_RMNHome *RMNHomeCallerSession) GetCandidateDigest() ([32]byte, error) {
	return _RMNHome.Contract.GetCandidateDigest(&_RMNHome.CallOpts)
}

func (_RMNHome *RMNHomeCaller) GetConfig(opts *bind.CallOpts, configDigest [32]byte) (GetConfig,

	error) {
	var out []interface{}
	err := _RMNHome.contract.Call(opts, &out, "getConfig", configDigest)

	outstruct := new(GetConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.VersionedConfig = *abi.ConvertType(out[0], new(RMNHomeVersionedConfig)).(*RMNHomeVersionedConfig)
	outstruct.Ok = *abi.ConvertType(out[1], new(bool)).(*bool)

	return *outstruct, err

}

func (_RMNHome *RMNHomeSession) GetConfig(configDigest [32]byte) (GetConfig,

	error) {
	return _RMNHome.Contract.GetConfig(&_RMNHome.CallOpts, configDigest)
}

func (_RMNHome *RMNHomeCallerSession) GetConfig(configDigest [32]byte) (GetConfig,

	error) {
	return _RMNHome.Contract.GetConfig(&_RMNHome.CallOpts, configDigest)
}

func (_RMNHome *RMNHomeCaller) GetConfigDigests(opts *bind.CallOpts) (GetConfigDigests,

	error) {
	var out []interface{}
	err := _RMNHome.contract.Call(opts, &out, "getConfigDigests")

	outstruct := new(GetConfigDigests)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ActiveConfigDigest = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.CandidateConfigDigest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

func (_RMNHome *RMNHomeSession) GetConfigDigests() (GetConfigDigests,

	error) {
	return _RMNHome.Contract.GetConfigDigests(&_RMNHome.CallOpts)
}

func (_RMNHome *RMNHomeCallerSession) GetConfigDigests() (GetConfigDigests,

	error) {
	return _RMNHome.Contract.GetConfigDigests(&_RMNHome.CallOpts)
}

func (_RMNHome *RMNHomeCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RMNHome.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_RMNHome *RMNHomeSession) Owner() (common.Address, error) {
	return _RMNHome.Contract.Owner(&_RMNHome.CallOpts)
}

func (_RMNHome *RMNHomeCallerSession) Owner() (common.Address, error) {
	return _RMNHome.Contract.Owner(&_RMNHome.CallOpts)
}

func (_RMNHome *RMNHomeCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _RMNHome.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_RMNHome *RMNHomeSession) TypeAndVersion() (string, error) {
	return _RMNHome.Contract.TypeAndVersion(&_RMNHome.CallOpts)
}

func (_RMNHome *RMNHomeCallerSession) TypeAndVersion() (string, error) {
	return _RMNHome.Contract.TypeAndVersion(&_RMNHome.CallOpts)
}

func (_RMNHome *RMNHomeTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RMNHome.contract.Transact(opts, "acceptOwnership")
}

func (_RMNHome *RMNHomeSession) AcceptOwnership() (*types.Transaction, error) {
	return _RMNHome.Contract.AcceptOwnership(&_RMNHome.TransactOpts)
}

func (_RMNHome *RMNHomeTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _RMNHome.Contract.AcceptOwnership(&_RMNHome.TransactOpts)
}

func (_RMNHome *RMNHomeTransactor) PromoteCandidateAndRevokeActive(opts *bind.TransactOpts, digestToPromote [32]byte, digestToRevoke [32]byte) (*types.Transaction, error) {
	return _RMNHome.contract.Transact(opts, "promoteCandidateAndRevokeActive", digestToPromote, digestToRevoke)
}

func (_RMNHome *RMNHomeSession) PromoteCandidateAndRevokeActive(digestToPromote [32]byte, digestToRevoke [32]byte) (*types.Transaction, error) {
	return _RMNHome.Contract.PromoteCandidateAndRevokeActive(&_RMNHome.TransactOpts, digestToPromote, digestToRevoke)
}

func (_RMNHome *RMNHomeTransactorSession) PromoteCandidateAndRevokeActive(digestToPromote [32]byte, digestToRevoke [32]byte) (*types.Transaction, error) {
	return _RMNHome.Contract.PromoteCandidateAndRevokeActive(&_RMNHome.TransactOpts, digestToPromote, digestToRevoke)
}

func (_RMNHome *RMNHomeTransactor) RevokeCandidate(opts *bind.TransactOpts, configDigest [32]byte) (*types.Transaction, error) {
	return _RMNHome.contract.Transact(opts, "revokeCandidate", configDigest)
}

func (_RMNHome *RMNHomeSession) RevokeCandidate(configDigest [32]byte) (*types.Transaction, error) {
	return _RMNHome.Contract.RevokeCandidate(&_RMNHome.TransactOpts, configDigest)
}

func (_RMNHome *RMNHomeTransactorSession) RevokeCandidate(configDigest [32]byte) (*types.Transaction, error) {
	return _RMNHome.Contract.RevokeCandidate(&_RMNHome.TransactOpts, configDigest)
}

func (_RMNHome *RMNHomeTransactor) SetCandidate(opts *bind.TransactOpts, staticConfig RMNHomeStaticConfig, dynamicConfig RMNHomeDynamicConfig, digestToOverwrite [32]byte) (*types.Transaction, error) {
	return _RMNHome.contract.Transact(opts, "setCandidate", staticConfig, dynamicConfig, digestToOverwrite)
}

func (_RMNHome *RMNHomeSession) SetCandidate(staticConfig RMNHomeStaticConfig, dynamicConfig RMNHomeDynamicConfig, digestToOverwrite [32]byte) (*types.Transaction, error) {
	return _RMNHome.Contract.SetCandidate(&_RMNHome.TransactOpts, staticConfig, dynamicConfig, digestToOverwrite)
}

func (_RMNHome *RMNHomeTransactorSession) SetCandidate(staticConfig RMNHomeStaticConfig, dynamicConfig RMNHomeDynamicConfig, digestToOverwrite [32]byte) (*types.Transaction, error) {
	return _RMNHome.Contract.SetCandidate(&_RMNHome.TransactOpts, staticConfig, dynamicConfig, digestToOverwrite)
}

func (_RMNHome *RMNHomeTransactor) SetDynamicConfig(opts *bind.TransactOpts, newDynamicConfig RMNHomeDynamicConfig, currentDigest [32]byte) (*types.Transaction, error) {
	return _RMNHome.contract.Transact(opts, "setDynamicConfig", newDynamicConfig, currentDigest)
}

func (_RMNHome *RMNHomeSession) SetDynamicConfig(newDynamicConfig RMNHomeDynamicConfig, currentDigest [32]byte) (*types.Transaction, error) {
	return _RMNHome.Contract.SetDynamicConfig(&_RMNHome.TransactOpts, newDynamicConfig, currentDigest)
}

func (_RMNHome *RMNHomeTransactorSession) SetDynamicConfig(newDynamicConfig RMNHomeDynamicConfig, currentDigest [32]byte) (*types.Transaction, error) {
	return _RMNHome.Contract.SetDynamicConfig(&_RMNHome.TransactOpts, newDynamicConfig, currentDigest)
}

func (_RMNHome *RMNHomeTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _RMNHome.contract.Transact(opts, "transferOwnership", to)
}

func (_RMNHome *RMNHomeSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _RMNHome.Contract.TransferOwnership(&_RMNHome.TransactOpts, to)
}

func (_RMNHome *RMNHomeTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _RMNHome.Contract.TransferOwnership(&_RMNHome.TransactOpts, to)
}

type RMNHomeActiveConfigRevokedIterator struct {
	Event *RMNHomeActiveConfigRevoked

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RMNHomeActiveConfigRevokedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMNHomeActiveConfigRevoked)
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
		it.Event = new(RMNHomeActiveConfigRevoked)
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

func (it *RMNHomeActiveConfigRevokedIterator) Error() error {
	return it.fail
}

func (it *RMNHomeActiveConfigRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RMNHomeActiveConfigRevoked struct {
	ConfigDigest [32]byte
	Raw          types.Log
}

func (_RMNHome *RMNHomeFilterer) FilterActiveConfigRevoked(opts *bind.FilterOpts, configDigest [][32]byte) (*RMNHomeActiveConfigRevokedIterator, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _RMNHome.contract.FilterLogs(opts, "ActiveConfigRevoked", configDigestRule)
	if err != nil {
		return nil, err
	}
	return &RMNHomeActiveConfigRevokedIterator{contract: _RMNHome.contract, event: "ActiveConfigRevoked", logs: logs, sub: sub}, nil
}

func (_RMNHome *RMNHomeFilterer) WatchActiveConfigRevoked(opts *bind.WatchOpts, sink chan<- *RMNHomeActiveConfigRevoked, configDigest [][32]byte) (event.Subscription, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _RMNHome.contract.WatchLogs(opts, "ActiveConfigRevoked", configDigestRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RMNHomeActiveConfigRevoked)
				if err := _RMNHome.contract.UnpackLog(event, "ActiveConfigRevoked", log); err != nil {
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

func (_RMNHome *RMNHomeFilterer) ParseActiveConfigRevoked(log types.Log) (*RMNHomeActiveConfigRevoked, error) {
	event := new(RMNHomeActiveConfigRevoked)
	if err := _RMNHome.contract.UnpackLog(event, "ActiveConfigRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RMNHomeCandidateConfigRevokedIterator struct {
	Event *RMNHomeCandidateConfigRevoked

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RMNHomeCandidateConfigRevokedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMNHomeCandidateConfigRevoked)
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
		it.Event = new(RMNHomeCandidateConfigRevoked)
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

func (it *RMNHomeCandidateConfigRevokedIterator) Error() error {
	return it.fail
}

func (it *RMNHomeCandidateConfigRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RMNHomeCandidateConfigRevoked struct {
	ConfigDigest [32]byte
	Raw          types.Log
}

func (_RMNHome *RMNHomeFilterer) FilterCandidateConfigRevoked(opts *bind.FilterOpts, configDigest [][32]byte) (*RMNHomeCandidateConfigRevokedIterator, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _RMNHome.contract.FilterLogs(opts, "CandidateConfigRevoked", configDigestRule)
	if err != nil {
		return nil, err
	}
	return &RMNHomeCandidateConfigRevokedIterator{contract: _RMNHome.contract, event: "CandidateConfigRevoked", logs: logs, sub: sub}, nil
}

func (_RMNHome *RMNHomeFilterer) WatchCandidateConfigRevoked(opts *bind.WatchOpts, sink chan<- *RMNHomeCandidateConfigRevoked, configDigest [][32]byte) (event.Subscription, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _RMNHome.contract.WatchLogs(opts, "CandidateConfigRevoked", configDigestRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RMNHomeCandidateConfigRevoked)
				if err := _RMNHome.contract.UnpackLog(event, "CandidateConfigRevoked", log); err != nil {
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

func (_RMNHome *RMNHomeFilterer) ParseCandidateConfigRevoked(log types.Log) (*RMNHomeCandidateConfigRevoked, error) {
	event := new(RMNHomeCandidateConfigRevoked)
	if err := _RMNHome.contract.UnpackLog(event, "CandidateConfigRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RMNHomeConfigPromotedIterator struct {
	Event *RMNHomeConfigPromoted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RMNHomeConfigPromotedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMNHomeConfigPromoted)
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
		it.Event = new(RMNHomeConfigPromoted)
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

func (it *RMNHomeConfigPromotedIterator) Error() error {
	return it.fail
}

func (it *RMNHomeConfigPromotedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RMNHomeConfigPromoted struct {
	ConfigDigest [32]byte
	Raw          types.Log
}

func (_RMNHome *RMNHomeFilterer) FilterConfigPromoted(opts *bind.FilterOpts, configDigest [][32]byte) (*RMNHomeConfigPromotedIterator, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _RMNHome.contract.FilterLogs(opts, "ConfigPromoted", configDigestRule)
	if err != nil {
		return nil, err
	}
	return &RMNHomeConfigPromotedIterator{contract: _RMNHome.contract, event: "ConfigPromoted", logs: logs, sub: sub}, nil
}

func (_RMNHome *RMNHomeFilterer) WatchConfigPromoted(opts *bind.WatchOpts, sink chan<- *RMNHomeConfigPromoted, configDigest [][32]byte) (event.Subscription, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _RMNHome.contract.WatchLogs(opts, "ConfigPromoted", configDigestRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RMNHomeConfigPromoted)
				if err := _RMNHome.contract.UnpackLog(event, "ConfigPromoted", log); err != nil {
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

func (_RMNHome *RMNHomeFilterer) ParseConfigPromoted(log types.Log) (*RMNHomeConfigPromoted, error) {
	event := new(RMNHomeConfigPromoted)
	if err := _RMNHome.contract.UnpackLog(event, "ConfigPromoted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RMNHomeConfigSetIterator struct {
	Event *RMNHomeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RMNHomeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMNHomeConfigSet)
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
		it.Event = new(RMNHomeConfigSet)
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

func (it *RMNHomeConfigSetIterator) Error() error {
	return it.fail
}

func (it *RMNHomeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RMNHomeConfigSet struct {
	ConfigDigest  [32]byte
	Version       uint32
	StaticConfig  RMNHomeStaticConfig
	DynamicConfig RMNHomeDynamicConfig
	Raw           types.Log
}

func (_RMNHome *RMNHomeFilterer) FilterConfigSet(opts *bind.FilterOpts, configDigest [][32]byte) (*RMNHomeConfigSetIterator, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _RMNHome.contract.FilterLogs(opts, "ConfigSet", configDigestRule)
	if err != nil {
		return nil, err
	}
	return &RMNHomeConfigSetIterator{contract: _RMNHome.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_RMNHome *RMNHomeFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *RMNHomeConfigSet, configDigest [][32]byte) (event.Subscription, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _RMNHome.contract.WatchLogs(opts, "ConfigSet", configDigestRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RMNHomeConfigSet)
				if err := _RMNHome.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_RMNHome *RMNHomeFilterer) ParseConfigSet(log types.Log) (*RMNHomeConfigSet, error) {
	event := new(RMNHomeConfigSet)
	if err := _RMNHome.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RMNHomeDynamicConfigSetIterator struct {
	Event *RMNHomeDynamicConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RMNHomeDynamicConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMNHomeDynamicConfigSet)
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
		it.Event = new(RMNHomeDynamicConfigSet)
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

func (it *RMNHomeDynamicConfigSetIterator) Error() error {
	return it.fail
}

func (it *RMNHomeDynamicConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RMNHomeDynamicConfigSet struct {
	ConfigDigest  [32]byte
	DynamicConfig RMNHomeDynamicConfig
	Raw           types.Log
}

func (_RMNHome *RMNHomeFilterer) FilterDynamicConfigSet(opts *bind.FilterOpts, configDigest [][32]byte) (*RMNHomeDynamicConfigSetIterator, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _RMNHome.contract.FilterLogs(opts, "DynamicConfigSet", configDigestRule)
	if err != nil {
		return nil, err
	}
	return &RMNHomeDynamicConfigSetIterator{contract: _RMNHome.contract, event: "DynamicConfigSet", logs: logs, sub: sub}, nil
}

func (_RMNHome *RMNHomeFilterer) WatchDynamicConfigSet(opts *bind.WatchOpts, sink chan<- *RMNHomeDynamicConfigSet, configDigest [][32]byte) (event.Subscription, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _RMNHome.contract.WatchLogs(opts, "DynamicConfigSet", configDigestRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RMNHomeDynamicConfigSet)
				if err := _RMNHome.contract.UnpackLog(event, "DynamicConfigSet", log); err != nil {
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

func (_RMNHome *RMNHomeFilterer) ParseDynamicConfigSet(log types.Log) (*RMNHomeDynamicConfigSet, error) {
	event := new(RMNHomeDynamicConfigSet)
	if err := _RMNHome.contract.UnpackLog(event, "DynamicConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RMNHomeOwnershipTransferRequestedIterator struct {
	Event *RMNHomeOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RMNHomeOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMNHomeOwnershipTransferRequested)
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
		it.Event = new(RMNHomeOwnershipTransferRequested)
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

func (it *RMNHomeOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *RMNHomeOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RMNHomeOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_RMNHome *RMNHomeFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RMNHomeOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _RMNHome.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &RMNHomeOwnershipTransferRequestedIterator{contract: _RMNHome.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_RMNHome *RMNHomeFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *RMNHomeOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _RMNHome.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RMNHomeOwnershipTransferRequested)
				if err := _RMNHome.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_RMNHome *RMNHomeFilterer) ParseOwnershipTransferRequested(log types.Log) (*RMNHomeOwnershipTransferRequested, error) {
	event := new(RMNHomeOwnershipTransferRequested)
	if err := _RMNHome.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RMNHomeOwnershipTransferredIterator struct {
	Event *RMNHomeOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RMNHomeOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMNHomeOwnershipTransferred)
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
		it.Event = new(RMNHomeOwnershipTransferred)
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

func (it *RMNHomeOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *RMNHomeOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RMNHomeOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_RMNHome *RMNHomeFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RMNHomeOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _RMNHome.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &RMNHomeOwnershipTransferredIterator{contract: _RMNHome.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_RMNHome *RMNHomeFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *RMNHomeOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _RMNHome.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RMNHomeOwnershipTransferred)
				if err := _RMNHome.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_RMNHome *RMNHomeFilterer) ParseOwnershipTransferred(log types.Log) (*RMNHomeOwnershipTransferred, error) {
	event := new(RMNHomeOwnershipTransferred)
	if err := _RMNHome.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetAllConfigs struct {
	ActiveConfig    RMNHomeVersionedConfig
	CandidateConfig RMNHomeVersionedConfig
}
type GetConfig struct {
	VersionedConfig RMNHomeVersionedConfig
	Ok              bool
}
type GetConfigDigests struct {
	ActiveConfigDigest    [32]byte
	CandidateConfigDigest [32]byte
}

func (_RMNHome *RMNHome) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _RMNHome.abi.Events["ActiveConfigRevoked"].ID:
		return _RMNHome.ParseActiveConfigRevoked(log)
	case _RMNHome.abi.Events["CandidateConfigRevoked"].ID:
		return _RMNHome.ParseCandidateConfigRevoked(log)
	case _RMNHome.abi.Events["ConfigPromoted"].ID:
		return _RMNHome.ParseConfigPromoted(log)
	case _RMNHome.abi.Events["ConfigSet"].ID:
		return _RMNHome.ParseConfigSet(log)
	case _RMNHome.abi.Events["DynamicConfigSet"].ID:
		return _RMNHome.ParseDynamicConfigSet(log)
	case _RMNHome.abi.Events["OwnershipTransferRequested"].ID:
		return _RMNHome.ParseOwnershipTransferRequested(log)
	case _RMNHome.abi.Events["OwnershipTransferred"].ID:
		return _RMNHome.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (RMNHomeActiveConfigRevoked) Topic() common.Hash {
	return common.HexToHash("0x0b31c0055e2d464bef7781994b98c4ff9ef4ae0d05f59feb6a68c42de5e201b8")
}

func (RMNHomeCandidateConfigRevoked) Topic() common.Hash {
	return common.HexToHash("0x53f5d9228f0a4173bea6e5931c9b3afe6eeb6692ede1d182952970f152534e3b")
}

func (RMNHomeConfigPromoted) Topic() common.Hash {
	return common.HexToHash("0xfc3e98dbbd47c3fa7c1c05b6ec711caeaf70eca4554192b9ada8fc11a37f298e")
}

func (RMNHomeConfigSet) Topic() common.Hash {
	return common.HexToHash("0xf6c6d1be15ba0acc8ee645c1ec613c360ef786d2d3200eb8e695b6dec757dbf0")
}

func (RMNHomeDynamicConfigSet) Topic() common.Hash {
	return common.HexToHash("0x1f69d1a2edb327babc986b3deb80091f101b9105d42a6c30db4d99c31d7e6294")
}

func (RMNHomeOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (RMNHomeOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_RMNHome *RMNHome) Address() common.Address {
	return _RMNHome.address
}

type RMNHomeInterface interface {
	GetActiveDigest(opts *bind.CallOpts) ([32]byte, error)

	GetAllConfigs(opts *bind.CallOpts) (GetAllConfigs,

		error)

	GetCandidateDigest(opts *bind.CallOpts) ([32]byte, error)

	GetConfig(opts *bind.CallOpts, configDigest [32]byte) (GetConfig,

		error)

	GetConfigDigests(opts *bind.CallOpts) (GetConfigDigests,

		error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	PromoteCandidateAndRevokeActive(opts *bind.TransactOpts, digestToPromote [32]byte, digestToRevoke [32]byte) (*types.Transaction, error)

	RevokeCandidate(opts *bind.TransactOpts, configDigest [32]byte) (*types.Transaction, error)

	SetCandidate(opts *bind.TransactOpts, staticConfig RMNHomeStaticConfig, dynamicConfig RMNHomeDynamicConfig, digestToOverwrite [32]byte) (*types.Transaction, error)

	SetDynamicConfig(opts *bind.TransactOpts, newDynamicConfig RMNHomeDynamicConfig, currentDigest [32]byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterActiveConfigRevoked(opts *bind.FilterOpts, configDigest [][32]byte) (*RMNHomeActiveConfigRevokedIterator, error)

	WatchActiveConfigRevoked(opts *bind.WatchOpts, sink chan<- *RMNHomeActiveConfigRevoked, configDigest [][32]byte) (event.Subscription, error)

	ParseActiveConfigRevoked(log types.Log) (*RMNHomeActiveConfigRevoked, error)

	FilterCandidateConfigRevoked(opts *bind.FilterOpts, configDigest [][32]byte) (*RMNHomeCandidateConfigRevokedIterator, error)

	WatchCandidateConfigRevoked(opts *bind.WatchOpts, sink chan<- *RMNHomeCandidateConfigRevoked, configDigest [][32]byte) (event.Subscription, error)

	ParseCandidateConfigRevoked(log types.Log) (*RMNHomeCandidateConfigRevoked, error)

	FilterConfigPromoted(opts *bind.FilterOpts, configDigest [][32]byte) (*RMNHomeConfigPromotedIterator, error)

	WatchConfigPromoted(opts *bind.WatchOpts, sink chan<- *RMNHomeConfigPromoted, configDigest [][32]byte) (event.Subscription, error)

	ParseConfigPromoted(log types.Log) (*RMNHomeConfigPromoted, error)

	FilterConfigSet(opts *bind.FilterOpts, configDigest [][32]byte) (*RMNHomeConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *RMNHomeConfigSet, configDigest [][32]byte) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*RMNHomeConfigSet, error)

	FilterDynamicConfigSet(opts *bind.FilterOpts, configDigest [][32]byte) (*RMNHomeDynamicConfigSetIterator, error)

	WatchDynamicConfigSet(opts *bind.WatchOpts, sink chan<- *RMNHomeDynamicConfigSet, configDigest [][32]byte) (event.Subscription, error)

	ParseDynamicConfigSet(log types.Log) (*RMNHomeDynamicConfigSet, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RMNHomeOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *RMNHomeOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*RMNHomeOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RMNHomeOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *RMNHomeOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*RMNHomeOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
