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
	F                   uint64
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
	ABI: "[{\"inputs\":[],\"name\":\"CannotTransferToSelf\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"expectedConfigDigest\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"gotConfigDigest\",\"type\":\"bytes32\"}],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"DigestNotFound\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateOffchainPublicKey\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicatePeerId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateSourceChain\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeProposedOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoOpStateTransitionNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotEnoughObservers\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OutOfBoundsNodesLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OutOfBoundsObserverNodeIndex\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OwnerCannotBeZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RevokingZeroDigestNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"ActiveConfigRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"CandidateConfigRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"ConfigPromoted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"peerId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"offchainPublicKey\",\"type\":\"bytes32\"}],\"internalType\":\"structRMNHome.Node[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"indexed\":false,\"internalType\":\"structRMNHome.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"f\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"observerNodesBitmap\",\"type\":\"uint256\"}],\"internalType\":\"structRMNHome.SourceChain[]\",\"name\":\"sourceChains\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"indexed\":false,\"internalType\":\"structRMNHome.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"components\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"f\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"observerNodesBitmap\",\"type\":\"uint256\"}],\"internalType\":\"structRMNHome.SourceChain[]\",\"name\":\"sourceChains\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"indexed\":false,\"internalType\":\"structRMNHome.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"name\":\"DynamicConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getActiveDigest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllConfigs\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"peerId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"offchainPublicKey\",\"type\":\"bytes32\"}],\"internalType\":\"structRMNHome.Node[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structRMNHome.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"f\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"observerNodesBitmap\",\"type\":\"uint256\"}],\"internalType\":\"structRMNHome.SourceChain[]\",\"name\":\"sourceChains\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structRMNHome.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"internalType\":\"structRMNHome.VersionedConfig\",\"name\":\"activeConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"peerId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"offchainPublicKey\",\"type\":\"bytes32\"}],\"internalType\":\"structRMNHome.Node[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structRMNHome.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"f\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"observerNodesBitmap\",\"type\":\"uint256\"}],\"internalType\":\"structRMNHome.SourceChain[]\",\"name\":\"sourceChains\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structRMNHome.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"internalType\":\"structRMNHome.VersionedConfig\",\"name\":\"candidateConfig\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCandidateDigest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"getConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"peerId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"offchainPublicKey\",\"type\":\"bytes32\"}],\"internalType\":\"structRMNHome.Node[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structRMNHome.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"f\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"observerNodesBitmap\",\"type\":\"uint256\"}],\"internalType\":\"structRMNHome.SourceChain[]\",\"name\":\"sourceChains\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structRMNHome.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"internalType\":\"structRMNHome.VersionedConfig\",\"name\":\"versionedConfig\",\"type\":\"tuple\"},{\"internalType\":\"bool\",\"name\":\"ok\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfigDigests\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"activeConfigDigest\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"candidateConfigDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"digestToPromote\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"digestToRevoke\",\"type\":\"bytes32\"}],\"name\":\"promoteCandidateAndRevokeActive\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"revokeCandidate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"peerId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"offchainPublicKey\",\"type\":\"bytes32\"}],\"internalType\":\"structRMNHome.Node[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structRMNHome.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"f\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"observerNodesBitmap\",\"type\":\"uint256\"}],\"internalType\":\"structRMNHome.SourceChain[]\",\"name\":\"sourceChains\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structRMNHome.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"digestToOverwrite\",\"type\":\"bytes32\"}],\"name\":\"setCandidate\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"newConfigDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"f\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"observerNodesBitmap\",\"type\":\"uint256\"}],\"internalType\":\"structRMNHome.SourceChain[]\",\"name\":\"sourceChains\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structRMNHome.DynamicConfig\",\"name\":\"newDynamicConfig\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"currentDigest\",\"type\":\"bytes32\"}],\"name\":\"setDynamicConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60808060405234604d573315603c57600180546001600160a01b03191633179055600e80546001600160401b031916905561244690816100538239f35b639b15e16f60e01b60005260046000fd5b600080fdfe6080604052600436101561001257600080fd5b60003560e01c8063118dbac51461077e578063123e65db14610746578063181f5a77146106e75780633567e6b41461068b57806338354c5c1461066857806363507956146105a85780636dd5b69d14610569578063736be8021461051e57806379ba5097146104535780638c76967f146102cb5780638da5cb5b14610297578063f2fde38b146101c25763fb4022d4146100ab57600080fd5b346101bd5760206003193601126101bd576004356100c76121d4565b80156101935763ffffffff6100e7600163ffffffff600e5460201c161890565b1660028110159081610134576006600391020191825481810361016357507f53f5d9228f0a4173bea6e5931c9b3afe6eeb6692ede1d182952970f152534e3b600080a26101345760009055005b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f93df584c0000000000000000000000000000000000000000000000000000000060005260045260245260446000fd5b7f0849d8cc0000000000000000000000000000000000000000000000000000000060005260046000fd5b600080fd5b346101bd5760206003193601126101bd5760043573ffffffffffffffffffffffffffffffffffffffff81168091036101bd576101fc6121d4565b33811461026d57807fffffffffffffffffffffffff0000000000000000000000000000000000000000600054161760005573ffffffffffffffffffffffffffffffffffffffff600154167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278600080a3005b7fdad89dca0000000000000000000000000000000000000000000000000000000060005260046000fd5b346101bd5760006003193601126101bd57602073ffffffffffffffffffffffffffffffffffffffff60015416604051908152f35b346101bd5760406003193601126101bd576004356024356102ea6121d4565b81158061044b575b6104215763ffffffff610310600163ffffffff600e5460201c161890565b16600281101561013457600602600301548281036103ef575063ffffffff600e5460201c1660028110156101345760060260030180548281036103ef575060009055600e547fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff67ffffffff00000000600163ffffffff8460201c161860201b16911617600e55806103c4575b507ffc3e98dbbd47c3fa7c1c05b6ec711caeaf70eca4554192b9ada8fc11a37f298e600080a2005b7f0b31c0055e2d464bef7781994b98c4ff9ef4ae0d05f59feb6a68c42de5e201b8600080a28161039c565b90507f93df584c0000000000000000000000000000000000000000000000000000000060005260045260245260446000fd5b7f7b4d1e4f0000000000000000000000000000000000000000000000000000000060005260046000fd5b5080156102f2565b346101bd5760006003193601126101bd5760005473ffffffffffffffffffffffffffffffffffffffff811633036104f4577fffffffffffffffffffffffff00000000000000000000000000000000000000006001549133828416176001551660005573ffffffffffffffffffffffffffffffffffffffff3391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0600080a3005b7f02b543c60000000000000000000000000000000000000000000000000000000060005260046000fd5b346101bd5760406003193601126101bd5760043567ffffffffffffffff81116101bd57604060031982360301126101bd576105679061055b6121d4565b60243590600401611d6a565b005b346101bd5760206003193601126101bd5761059c610588600435611d0a565b6040519283926040845260408401906114f2565b90151560208301520390f35b346101bd5760006003193601126101bd576105c1611ac0565b6105c9611ac0565b9063ffffffff600e5460201c1660028110156101345760066105ee9102600201611bb2565b6020810151610660575b50600e5460201c63ffffffff166001186002811015610134576106469261062760066106549302600201611bb2565b6020810151610658575b506040519384936040855260408501906114f2565b9083820360208501526114f2565b0390f35b905084610631565b9050826105f8565b346101bd5760006003193601126101bd576020610683611a85565b604051908152f35b346101bd5760006003193601126101bd5763ffffffff600e5460201c1660028110156101345760060260030154600e5460201c63ffffffff16600118906002821015610134576003600660409302015482519182526020820152f35b346101bd5760006003193601126101bd57610654604080519061070a818361144b565b601182527f524d4e486f6d6520312e362e302d6465760000000000000000000000000000006020830152519182916020835260208301906114af565b346101bd5760006003193601126101bd5763ffffffff600e5460201c1660028110156101345760036006602092020154604051908152f35b346101bd5760606003193601126101bd5760043567ffffffffffffffff81116101bd57604060031982360301126101bd5767ffffffffffffffff602435116101bd57602435360360406003198201126101bd576107d96121d4565b6040516000916107e8826113c8565b67ffffffffffffffff8460040135116113c4578360040135840190366023830112156113c05761081b600483013561162a565b610828604051918261144b565b6004830135808252602082019060061b84016024013681116110cd5760248501915b81831061138f575050508352602485013567ffffffffffffffff811161138b5761087a9060043691880101611642565b602084015261088e366024356004016116cc565b94610100845151116113635790919484955b8451518710156109965760018701808811610969575b8551805182101561095c57886108cb9161221f565b51516108d882885161221f565b5151146109345760206108ec89885161221f565b51015160206108fc83895161221f565b5101511461090c576001016108b6565b6004877fae00651d000000000000000000000000000000000000000000000000000000008152fd5b6004877f221a8ae8000000000000000000000000000000000000000000000000000000008152fd5b50506001909601956108a0565b6024877f4e487b710000000000000000000000000000000000000000000000000000000081526011600452fd5b9493906109a591515190612233565b6109ad611a85565b604435810361133357611307575b600e549363ffffffff85169463ffffffff86146112da5763ffffffff60017fffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000095969701169384911617600e557dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff604051610aec610b1a6020830160208152610a7784610a4b604082018a600401611846565b037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0810186528561144b565b602060405191818301957f45564d000000000000000000000000000000000000000000000000000000000087524660408501523060608501528a608085015260808452610ac560a08561144b565b604051958694610add858701998a925192839161148c565b8501915180938584019061148c565b0101037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0810183528261144b565b519020167e0b0000000000000000000000000000000000000000000000000000000000001793610b55600163ffffffff600e5460201c161890565b60028110156112ad576006029182600201866003850155857fffffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000825416179055600483017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdd85360301856004013512156110c95767ffffffffffffffff6004830135116110c957600482013560061b360360248301136110c9576801000000000000000060048301351161109c5780546004830135808355811161121f575b508752602087208790602483015b600484013583106111fb575050505060058201610c4460248501856004016118f9565b9067ffffffffffffffff82116111ce578190610c60845461194a565b601f811161117e575b508990601f83116001146110dc578a926110d1575b50507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790555b60068201907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdd6024356004013591018112156110cd57602435019060048201359167ffffffffffffffff83116110c95760240160608302360381136110c95768010000000000000000831161109c578154838355808410610fea575b509087526020872087915b838310610f70575050505060070191610d5960248035016024356004016118f9565b67ffffffffffffffff8195929511610f4357610d75825461194a565b601f8111610efe575b5090859392918760209890601f8311600114610e3a57907ff6c6d1be15ba0acc8ee645c1ec613c360ef786d2d3200eb8e695b6dec757dbf0969783610e2f575b50507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790555b610e24610e1160405193849384526060898501526060840190600401611846565b82810360408401526024356004016119b2565b0390a2604051908152f35b013590508980610dbe565b96907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe083168489528a8920985b818110610ee45750917ff6c6d1be15ba0acc8ee645c1ec613c360ef786d2d3200eb8e695b6dec757dbf0979891846001959410610eac575b505050811b019055610df0565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88560031b161c19910135169055898080610e9f565b828401358a556001909901988a9850918b01918b01610e67565b82885260208820601f830160051c81019160208410610f39575b601f0160051c01905b818110610f2e5750610d7e565b888155600101610f21565b9091508190610f18565b6024877f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b600260608267ffffffffffffffff610f8960019561199d565b168554907fffffffffffffffffffffffffffffffff000000000000000000000000000000006fffffffffffffffff0000000000000000610fcb6020860161199d565b60401b1692161717855560408101358486015501920192019190610d37565b7f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8116810361106f577f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8416840361106f57828952602089209060011b8101908460011b015b81811061105d5750610d2c565b808a600292558a600182015501611050565b6024897f4e487b710000000000000000000000000000000000000000000000000000000081526011600452fd5b6024887f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b8780fd5b8680fd5b013590508980610c7e565b90917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01691848b5260208b20928b5b818110611166575090846001959493921061112e575b505050811b019055610cb0565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88560031b161c19910135169055898080611121565b9193602060018192878701358155019501920161110b565b909150838a5260208a20601f840160051c810191602085106111c4575b90601f859493920160051c01905b8181106111b65750610c69565b8b81558493506001016111a9565b909150819061119b565b6024897f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b60016002604083600494358655602081013584870155019301930192919050610c21565b7f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8116810361106f5760048301357f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8116900361106f57818952602089209060011b810190600484013560011b015b81811061129b5750610c13565b6002908a81558a60018201550161128e565b6024877f4e487b710000000000000000000000000000000000000000000000000000000081526032600452fd5b6024857f4e487b710000000000000000000000000000000000000000000000000000000081526011600452fd5b6044357f53f5d9228f0a4173bea6e5931c9b3afe6eeb6692ede1d182952970f152534e3b8480a26109bb565b7f93df584c0000000000000000000000000000000000000000000000000000000084526004526044803560245283fd5b6004857faf26d5e3000000000000000000000000000000000000000000000000000000008152fd5b8480fd5b6040833603126110c957602060409182516113a9816113c8565b85358152828601358382015281520192019161084a565b8380fd5b8280fd5b6040810190811067ffffffffffffffff8211176113e457604052565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6060810190811067ffffffffffffffff8211176113e457604052565b6080810190811067ffffffffffffffff8211176113e457604052565b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff8211176113e457604052565b60005b83811061149f5750506000910152565b818101518382015260200161148f565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f6020936114eb8151809281875287808801910161148c565b0116010190565b91909163ffffffff81511683526020810151602084015260408101516080604085015260c0840190805191604060808701528251809152602060e0870193019060005b8181106116085750505060609160206115799201517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff808783030160a08801526114af565b910151926060818303910152604081019083519160408252825180915260206060830193019060005b8181106115c65750505060206115c3939401519060208184039101526114af565b90565b909193602060606001926040885167ffffffffffffffff815116835267ffffffffffffffff8582015116858401520151604082015201950191019190916115a2565b8251805186526020908101518187015260409095019490920191600101611535565b67ffffffffffffffff81116113e45760051b60200190565b81601f820112156101bd5780359067ffffffffffffffff82116113e45760405192611695601f84017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0166020018561144b565b828452602083830101116101bd57816000926020809301838601378301015290565b359067ffffffffffffffff821682036101bd57565b91906040838203126101bd57604051906116e5826113c8565b8193803567ffffffffffffffff81116101bd57810182601f820112156101bd5780356117108161162a565b9161171e604051938461144b565b818352602060608185019302820101908582116101bd57602001915b81831061176d57505050835260208101359167ffffffffffffffff83116101bd576020926117689201611642565b910152565b6060838703126101bd57602060609160405161178881611413565b611791866116b7565b815261179e8387016116b7565b838201526040860135604082015281520192019161173a565b90357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1823603018112156101bd57016020813591019167ffffffffffffffff82116101bd5781360383136101bd57565b601f82602094937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0938186528686013760008582860101520116010190565b9190604081019083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1853603018112156101bd5784016020813591019267ffffffffffffffff82116101bd578160061b360384136101bd578190604084525260608201929060005b8181106118d9575050506118cb8460206115c3959601906117b7565b916020818503910152611807565b8235855260208084013590860152604094850194909201916001016118af565b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1813603018212156101bd570180359067ffffffffffffffff82116101bd576020019181360383136101bd57565b90600182811c92168015611993575b602083101461196457565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b91607f1691611959565b3567ffffffffffffffff811681036101bd5790565b9190604081019083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1853603018112156101bd5784016020813591019267ffffffffffffffff82116101bd5760608202360384136101bd578190604084525260608201929060005b818110611a37575050506118cb8460206115c3959601906117b7565b90919360608060019267ffffffffffffffff611a52896116b7565b16815267ffffffffffffffff611a6a60208a016116b7565b16602082015260408881013590820152019501929101611a1b565b600e5460201c63ffffffff166001186002811015610134576006026003015490565b60405190611ab4826113c8565b60606020838281520152565b60405190611acd8261142f565b816000815260006020820152611ae1611aa7565b60408201526060611768611aa7565b9060405191826000825492611b048461194a565b8084529360018116908115611b725750600114611b2b575b50611b299250038361144b565b565b90506000929192526020600020906000915b818310611b56575050906020611b299282010138611b1c565b6020919350806001915483858901015201910190918492611b3d565b60209350611b299592507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0091501682840152151560051b82010138611b1c565b9060405191611bc08361142f565b8263ffffffff8254168152600182015460208201526002820160405190611be6826113c8565b8054611bf18161162a565b91611bff604051938461144b565b818352602083019060005260206000206000915b838310611cdd57505050508152611c2c60038401611af0565b60208201526040820152600482019160405192611c48846113c8565b8054611c538161162a565b91611c61604051938461144b565b818352602083019060005260206000206000915b838310611c9b5750505050600560609392611c9292865201611af0565b60208401520152565b60026020600192604051611cae81611413565b67ffffffffffffffff8654818116835260401c1683820152848601546040820152815201920192019190611c75565b60026020600192604051611cf0816113c8565b855481528486015483820152815201920192019190611c13565b611d12611ac0565b9060005b6002811015611d6257600060068202908360038301541480611d59575b611d41575050600101611d16565b91509150611d53925050600201611bb2565b90600190565b50831515611d33565b505090600090565b9060005b60028110156121a65760006006820290836003830154148061219d575b611d99575050600101611d6e565b90915092919250611db76004820154611db236856116cc565b612233565b6000906006810183357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1853603018112156113c05784019081359167ffffffffffffffff831161138b57602001606083023603811361138b576801000000000000000083116121705781548383558084106120be575b509084526020842084915b838310612044575050505060070190611e5460208401846118f9565b919067ffffffffffffffff831161201757611e6f845461194a565b601f8111611fd2575b5081601f8411600114611f0c5792611efc949281927f1f69d1a2edb327babc986b3deb80091f101b9105d42a6c30db4d99c31d7e6294979592611f01575b50507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790555b6040519182916020835260208301906119b2565b0390a2565b013590503880611eb6565b917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0841685845260208420935b818110611fba57509260019285927f1f69d1a2edb327babc986b3deb80091f101b9105d42a6c30db4d99c31d7e62949896611efc989610611f82575b505050811b019055611ee8565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88560031b161c19910135169055388080611f75565b91936020600181928787013581550195019201611f39565b84835260208320601f850160051c8101916020861061200d575b601f0160051c01905b8181106120025750611e78565b838155600101611ff5565b9091508190611fec565b6024827f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b600260608267ffffffffffffffff61205d60019561199d565b168554907fffffffffffffffffffffffffffffffff000000000000000000000000000000006fffffffffffffffff000000000000000061209f6020860161199d565b60401b1692161717855560408101358486015501920192019190611e38565b7f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81168103612143577f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8416840361214357828652602086209060011b8101908460011b015b8181106121315750611e2d565b80876002925587600182015501612124565b6024867f4e487b710000000000000000000000000000000000000000000000000000000081526011600452fd5b6024857f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b50831515611d8b565b507fd0b2c0310000000000000000000000000000000000000000000000000000000060005260045260246000fd5b73ffffffffffffffffffffffffffffffffffffffff6001541633036121f557565b7f2b5c74de0000000000000000000000000000000000000000000000000000000060005260046000fd5b80518210156101345760209160051b010190565b908151519160005b8381106122485750505050565b61225381835161221f565b5160018201808311612328575b8581106123df575060408101519084610100036101008111612328577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff83911c8116036123b557816000925b61235757506020015160011b6801fffffffffffffffe67fffffffffffffffe8216911681036123285760010167ffffffffffffffff81116123285767ffffffffffffffff16116122fe5760010161223b565b7fa804bcb30000000000000000000000000000000000000000000000000000000060005260046000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8101908082116123285716917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81146123285760010191806122ac565b7f2847b6060000000000000000000000000000000000000000000000000000000060005260046000fd5b67ffffffffffffffff82511667ffffffffffffffff6123ff83875161221f565b5151161461240f57600101612260565b7f3857f84d0000000000000000000000000000000000000000000000000000000060005260046000fdfea164736f6c634300081a000a",
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
