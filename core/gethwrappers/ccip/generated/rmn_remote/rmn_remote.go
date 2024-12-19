// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package rmn_remote

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

type IRMNRemoteSignature struct {
	R [32]byte
	S [32]byte
}

type IRMNTaggedRoot struct {
	CommitStore common.Address
	Root        [32]byte
}

type InternalMerkleRoot struct {
	SourceChainSelector uint64
	OnRampAddress       []byte
	MinSeqNr            uint64
	MaxSeqNr            uint64
	MerkleRoot          [32]byte
}

type RMNRemoteConfig struct {
	RmnHomeContractConfigDigest [32]byte
	Signers                     []RMNRemoteSigner
	F                           uint64
}

type RMNRemoteSigner struct {
	OnchainPublicKey common.Address
	NodeIndex        uint64
}

var RMNRemoteMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"localChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"contractIRMN\",\"name\":\"legacyRMN\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"subject\",\"type\":\"bytes16\"}],\"name\":\"AlreadyCursed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotTransferToSelf\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConfigNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateOnchainPublicKey\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSignature\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSignerOrder\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IsBlessedNotAvailable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeProposedOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"subject\",\"type\":\"bytes16\"}],\"name\":\"NotCursed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotEnoughSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OutOfOrderSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OwnerCannotBeZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ThresholdNotMet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnexpectedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroValueNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"rmnHomeContractConfigDigest\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"onchainPublicKey\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"nodeIndex\",\"type\":\"uint64\"}],\"internalType\":\"structRMNRemote.Signer[]\",\"name\":\"signers\",\"type\":\"tuple[]\"},{\"internalType\":\"uint64\",\"name\":\"f\",\"type\":\"uint64\"}],\"indexed\":false,\"internalType\":\"structRMNRemote.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16[]\",\"name\":\"subjects\",\"type\":\"bytes16[]\"}],\"name\":\"Cursed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16[]\",\"name\":\"subjects\",\"type\":\"bytes16[]\"}],\"name\":\"Uncursed\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"subject\",\"type\":\"bytes16\"}],\"name\":\"curse\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16[]\",\"name\":\"subjects\",\"type\":\"bytes16[]\"}],\"name\":\"curse\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCursedSubjects\",\"outputs\":[{\"internalType\":\"bytes16[]\",\"name\":\"subjects\",\"type\":\"bytes16[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLocalChainSelector\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"localChainSelector\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getReportDigestHeader\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"digestHeader\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getVersionedConfig\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"rmnHomeContractConfigDigest\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"onchainPublicKey\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"nodeIndex\",\"type\":\"uint64\"}],\"internalType\":\"structRMNRemote.Signer[]\",\"name\":\"signers\",\"type\":\"tuple[]\"},{\"internalType\":\"uint64\",\"name\":\"f\",\"type\":\"uint64\"}],\"internalType\":\"structRMNRemote.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"commitStore\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"internalType\":\"structIRMN.TaggedRoot\",\"name\":\"taggedRoot\",\"type\":\"tuple\"}],\"name\":\"isBlessed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"subject\",\"type\":\"bytes16\"}],\"name\":\"isCursed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isCursed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"rmnHomeContractConfigDigest\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"onchainPublicKey\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"nodeIndex\",\"type\":\"uint64\"}],\"internalType\":\"structRMNRemote.Signer[]\",\"name\":\"signers\",\"type\":\"tuple[]\"},{\"internalType\":\"uint64\",\"name\":\"f\",\"type\":\"uint64\"}],\"internalType\":\"structRMNRemote.Config\",\"name\":\"newConfig\",\"type\":\"tuple\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"subject\",\"type\":\"bytes16\"}],\"name\":\"uncurse\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16[]\",\"name\":\"subjects\",\"type\":\"bytes16[]\"}],\"name\":\"uncurse\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"offRampAddress\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRampAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"maxSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"merkleRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structInternal.MerkleRoot[]\",\"name\":\"merkleRoots\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structIRMNRemote.Signature[]\",\"name\":\"signatures\",\"type\":\"tuple[]\"}],\"name\":\"verify\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60c0346100d357601f611ffb38819003918201601f19168301916001600160401b038311848410176100d85780849260409485528339810103126100d35780516001600160401b038116918282036100d35760200151916001600160a01b03831683036100d35733156100c257600180546001600160a01b03191633179055156100b15760805260a052604051611f0c90816100ef82396080518181816102c2015261065e015260a05181610e510152f35b63273e150360e21b60005260046000fd5b639b15e16f60e01b60005260046000fd5b600080fd5b634e487b7160e01b600052604160045260246000fdfe6080604052600436101561001257600080fd5b60003560e01c8063181f5a771461180c578063198f0f77146111595780631add205f14610fac5780632cbc26bb14610f89578063397796f714610f645780634d61677114610e2757806362eed41514610d255780636509a95414610cea5780636d2d399314610be857806370a9089e1461055a57806379ba50971461048f5780638da5cb5b1461045b5780639a19b3291461036d578063d881e092146102e6578063eaa83ddd146102a1578063f2fde38b146101cf5763f8bb876e146100d757600080fd5b346101ca576100e536611991565b6100ed611cc4565b60005b81518110156101955761012e7fffffffffffffffffffffffffffffffff000000000000000000000000000000006101278385611cb0565b5116611ea5565b1561013b576001016100f0565b610166907fffffffffffffffffffffffffffffffff0000000000000000000000000000000092611cb0565b51167f19d5c79b0000000000000000000000000000000000000000000000000000000060005260045260246000fd5b6040517f1716e663a90a76d3b6c7e5f680673d1b051454c19c627e184c8daf28f3104f7490806101c58582611a3a565b0390a1005b600080fd5b346101ca5760206003193601126101ca5773ffffffffffffffffffffffffffffffffffffffff6101fd611956565b610205611cc4565b1633811461027757807fffffffffffffffffffffffff0000000000000000000000000000000000000000600054161760005573ffffffffffffffffffffffffffffffffffffffff600154167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278600080a3005b7fdad89dca0000000000000000000000000000000000000000000000000000000060005260046000fd5b346101ca5760006003193601126101ca57602060405167ffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000168152f35b346101ca5760006003193601126101ca5760405180602060065491828152019060066000527ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d3f9060005b818110610357576103538561034781870382611887565b60405191829182611a3a565b0390f35b8254845260209093019260019283019201610330565b346101ca5761037b36611991565b610383611cc4565b60005b815181101561042b576103c47fffffffffffffffffffffffffffffffff000000000000000000000000000000006103bd8385611cb0565b5116611d0f565b156103d157600101610386565b6103fc907fffffffffffffffffffffffffffffffff0000000000000000000000000000000092611cb0565b51167f73281fa10000000000000000000000000000000000000000000000000000000060005260045260246000fd5b6040517f0676e709c9cc74fa0519fd78f7c33be0f1b2b0bae0507c724aef7229379c6ba190806101c58582611a3a565b346101ca5760006003193601126101ca57602073ffffffffffffffffffffffffffffffffffffffff60015416604051908152f35b346101ca5760006003193601126101ca5760005473ffffffffffffffffffffffffffffffffffffffff81163303610530577fffffffffffffffffffffffff00000000000000000000000000000000000000006001549133828416176001551660005573ffffffffffffffffffffffffffffffffffffffff3391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0600080a3005b7f02b543c60000000000000000000000000000000000000000000000000000000060005260046000fd5b346101ca5760606003193601126101ca57610573611956565b67ffffffffffffffff602435116101ca573660236024350112156101ca57602435600401359067ffffffffffffffff82116101ca573660248360051b81350101116101ca576044359067ffffffffffffffff82116101ca57366023830112156101ca5767ffffffffffffffff8260040135116101ca57366024836004013560061b840101116101ca5763ffffffff6005541615610bbe5767ffffffffffffffff6106208160045416611b3e565b16826004013510610b945760025460405160c0810181811067ffffffffffffffff821117610b6557604052468152602081019267ffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000168452604082019230845273ffffffffffffffffffffffffffffffffffffffff6060840192168252608083019081526106b587611979565b916106c36040519384611887565b8783526000976024803501602085015b60248360051b8135010182106109ad575050509073ffffffffffffffffffffffffffffffffffffffff8095939260a0860193845260405196879567ffffffffffffffff602088019a7f9651943783dbf81935a60e98f218a9d9b5b28823fb2228bbd91320d632facf538c526040808a01526101208901995160608a015251166080880152511660a0860152511660c08401525160e08301525160c0610100830152805180935261014082019260206101408260051b85010192019388905b828210610914575050506107cc9250037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08101835282611887565b5190208290815b836004013583106107e2578480f35b60208560806107f986886004013560248a01611aea565b358361080d888a6004013560248c01611aea565b013560405191878352601b868401526040830152606082015282805260015afa156109095784519073ffffffffffffffffffffffffffffffffffffffff82169081156108e15773ffffffffffffffffffffffffffffffffffffffff82911610156108b9578552600860205260ff6040862054161561089157600192909201916107d3565b6004857faaaa9141000000000000000000000000000000000000000000000000000000008152fd5b6004867fbbe15e7f000000000000000000000000000000000000000000000000000000008152fd5b6004877f8baa579f000000000000000000000000000000000000000000000000000000008152fd5b6040513d86823e3d90fd5b91936020847ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffec08293600195970301855287519067ffffffffffffffff825116815260808061096f8585015160a08786015260a08501906118c8565b9367ffffffffffffffff604082015116604085015267ffffffffffffffff606082015116606085015201519101529601920192018593919492610791565b81359067ffffffffffffffff8211610b615760a07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdc836024350136030112610b61576040519060a0820182811067ffffffffffffffff821117610b3057604052610a1c60248481350101611b97565b82526044836024350101359167ffffffffffffffff8311610b5d5736604360243586018501011215610b5d5767ffffffffffffffff6024848682350101013511610b3057908d9160207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f6024878982350101013501160193610aa36040519586611887565b60248035870182019081013580875236910160440111610b2c57602495602095869560a49387908a90813586018101808301359060440186850137858235010101358301015285840152610afc60648289350101611b97565b6040840152610b1060848289350101611b97565b60608401528635010135608082015281520192019190506106d3565b8380fd5b60248e7f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b8d80fd5b8b80fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7f59fa4a930000000000000000000000000000000000000000000000000000000060005260046000fd5b7face124bc0000000000000000000000000000000000000000000000000000000060005260046000fd5b346101ca5760206003193601126101ca57610c01611927565b604090815190610c118383611887565b600182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe083013660208401377fffffffffffffffffffffffffffffffff00000000000000000000000000000000610c6883611ca3565b91169052610c74611cc4565b60005b8151811015610cbb57610cae7fffffffffffffffffffffffffffffffff000000000000000000000000000000006103bd8385611cb0565b156103d157600101610c77565b82517f0676e709c9cc74fa0519fd78f7c33be0f1b2b0bae0507c724aef7229379c6ba190806101c58582611a3a565b346101ca5760006003193601126101ca5760206040517f9651943783dbf81935a60e98f218a9d9b5b28823fb2228bbd91320d632facf538152f35b346101ca5760206003193601126101ca57610d3e611927565b604090815190610d4e8383611887565b600182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe083013660208401377fffffffffffffffffffffffffffffffff00000000000000000000000000000000610da583611ca3565b91169052610db1611cc4565b60005b8151811015610df857610deb7fffffffffffffffffffffffffffffffff000000000000000000000000000000006101278385611cb0565b1561013b57600101610db4565b82517f1716e663a90a76d3b6c7e5f680673d1b051454c19c627e184c8daf28f3104f7490806101c58582611a3a565b346101ca5760406003193601126101ca57600073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000168015610f3c576020604491604051928380927f4d61677100000000000000000000000000000000000000000000000000000000825273ffffffffffffffffffffffffffffffffffffffff610ec3611956565b16600483015260243560248301525afa908115610f31578291610eee575b6020826040519015158152f35b90506020813d602011610f29575b81610f0960209383611887565b81010312610f2557518015158103610f25576020915082610ee1565b5080fd5b3d9150610efc565b6040513d84823e3d90fd5b6004827f0a7c4edd000000000000000000000000000000000000000000000000000000008152fd5b346101ca5760006003193601126101ca576020610f7f611c46565b6040519015158152f35b346101ca5760206003193601126101ca576020610f7f610fa7611927565b611bac565b346101ca5760006003193601126101ca57600060408051610fcc8161186b565b82815260606020820152015263ffffffff6005541660405190610fee8261186b565b60025482526003549161100083611979565b9261100e6040519485611887565b808452600360009081527fc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85b9190602086015b8282106110fa57505050506020810192835267ffffffffffffffff6004541692604082019384526040519283526040602084015260a083019151604084015251906060808401528151809152602060c0840192019060005b8181106110b85750505067ffffffffffffffff8293511660808301520390f35b8251805173ffffffffffffffffffffffffffffffffffffffff16855260209081015167ffffffffffffffff168186015260409094019390920191600101611098565b6040516040810181811067ffffffffffffffff821117610b6557600192839260209260405267ffffffffffffffff885473ffffffffffffffffffffffffffffffffffffffff8116835260a01c1683820152815201940191019092611040565b346101ca5760206003193601126101ca5760043567ffffffffffffffff81116101ca57806004018136039160606003198401126101ca57611198611cc4565b81359081156117e257909260248201919060015b6111b68486611a96565b9050811015611297576111c98486611a96565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8301918383116112685761120a9260209261120492611aea565b01611b29565b67ffffffffffffffff8061122d6020611204866112278b8d611a96565b90611aea565b169116101561123e576001016111ac565b7f448515170000000000000000000000000000000000000000000000000000000060005260046000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b508390856112a58584611a96565b6044860191506112b482611b29565b60011b6801fffffffffffffffe67fffffffffffffffe821691168103611268576112e667ffffffffffffffff91611b3e565b16116117b857600354805b6116b5575060005b6113038786611a96565b90508110156113dc5773ffffffffffffffffffffffffffffffffffffffff611337611332836112278b8a611a96565b611b76565b16600052600860205260ff604060002054166113b2578073ffffffffffffffffffffffffffffffffffffffff6113766113326001946112278c8b611a96565b1660005260086020526040600020827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00825416179055016112f9565b7f28cae27d0000000000000000000000000000000000000000000000000000000060005260046000fd5b50846113ee8780959685600255611a96565b90680100000000000000008211610b65576003548260035580831061166f575b50600360005260206000206000915b8383106115c1575050505067ffffffffffffffff61143a83611b29565b167fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000060045416176004556005549463ffffffff86169563ffffffff87146112685763ffffffff60017fffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000098011696879116176005557fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdd6040519560208752608087019560208801523591018112156101ca57016024600482013591019267ffffffffffffffff82116101ca578160061b360384136101ca578190606060408701525260a0840192906000905b80821061156e57867f7f22bf988149dbe8de8fb879c6b97a4e56e68b2bd57421ce1a4e79d4ef6b496c87808867ffffffffffffffff6115638a611b97565b1660608301520390a2005b90919384359073ffffffffffffffffffffffffffffffffffffffff82168092036101ca5760408160019382935267ffffffffffffffff6115b060208a01611b97565b166020820152019501920190611525565b600160408273ffffffffffffffffffffffffffffffffffffffff6115e58495611b76565b167fffffffffffffffffffffffff000000000000000000000000000000000000000086541617855561161960208201611b29565b7fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff7bffffffffffffffff000000000000000000000000000000000000000087549260a01b1691161785550192019201919061141d565b60036000527fc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85b9081019083015b8181106116a9575061140e565b6000815560010161169c565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff810181811161126857600090600354111561178b57600390527fc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85a81015473ffffffffffffffffffffffffffffffffffffffff16600090815260086020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690558015611268577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01806112f1565b807f4e487b7100000000000000000000000000000000000000000000000000000000602492526032600452fd5b7f014c50200000000000000000000000000000000000000000000000000000000060005260046000fd5b7f9cf8540c0000000000000000000000000000000000000000000000000000000060005260046000fd5b346101ca5760006003193601126101ca57610353604080519061182f8183611887565b601382527f524d4e52656d6f746520312e362e302d646576000000000000000000000000006020830152519182916020835260208301906118c8565b6060810190811067ffffffffffffffff821117610b6557604052565b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff821117610b6557604052565b919082519283825260005b8481106119125750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8460006020809697860101520116010190565b806020809284010151828286010152016118d3565b600435907fffffffffffffffffffffffffffffffff00000000000000000000000000000000821682036101ca57565b6004359073ffffffffffffffffffffffffffffffffffffffff821682036101ca57565b67ffffffffffffffff8111610b655760051b60200190565b60206003198201126101ca576004359067ffffffffffffffff82116101ca57806023830112156101ca578160040135906119ca82611979565b926119d86040519485611887565b8284526024602085019360051b8201019182116101ca57602401915b818310611a015750505090565b82357fffffffffffffffffffffffffffffffff00000000000000000000000000000000811681036101ca578152602092830192016119f4565b602060408183019282815284518094520192019060005b818110611a5e5750505090565b82517fffffffffffffffffffffffffffffffff0000000000000000000000000000000016845260209384019390920191600101611a51565b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1813603018212156101ca570180359067ffffffffffffffff82116101ca57602001918160061b360383136101ca57565b9190811015611afa5760061b0190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b3567ffffffffffffffff811681036101ca5790565b67ffffffffffffffff60019116019067ffffffffffffffff821161126857565b8054821015611afa5760005260206000200190600090565b3573ffffffffffffffffffffffffffffffffffffffff811681036101ca5790565b359067ffffffffffffffff821682036101ca57565b60065415611c40577fffffffffffffffffffffffffffffffff0000000000000000000000000000000016600052600760205260406000205415801590611bef5790565b507f010000000000000000000000000000010000000000000000000000000000000060005260076020527f70b766b11586b6b505ed3893938b0cc6c6c98bd6f65e969ac311168d34e4f9e254151590565b50600090565b60065415611c9e577f010000000000000000000000000000010000000000000000000000000000000060005260076020527f70b766b11586b6b505ed3893938b0cc6c6c98bd6f65e969ac311168d34e4f9e254151590565b600090565b805115611afa5760200190565b8051821015611afa5760209160051b010190565b73ffffffffffffffffffffffffffffffffffffffff600154163303611ce557565b7f2b5c74de0000000000000000000000000000000000000000000000000000000060005260046000fd5b6000818152600760205260409020548015611e9e577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff810181811161126857600654907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820191821161126857818103611e2f575b5050506006548015611e00577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01611dbd816006611b5e565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82549160031b1b19169055600655600052600760205260006040812055600190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b611e86611e40611e51936006611b5e565b90549060031b1c9283926006611b5e565b81939154907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9060031b92831b921b19161790565b90556000526007602052604060002055388080611d84565b5050600090565b80600052600760205260406000205415600014611c405760065468010000000000000000811015610b6557611ee6611e518260018594016006556006611b5e565b905560065490600052600760205260406000205560019056fea164736f6c634300081a000a",
}

var RMNRemoteABI = RMNRemoteMetaData.ABI

var RMNRemoteBin = RMNRemoteMetaData.Bin

func DeployRMNRemote(auth *bind.TransactOpts, backend bind.ContractBackend, localChainSelector uint64, legacyRMN common.Address) (common.Address, *types.Transaction, *RMNRemote, error) {
	parsed, err := RMNRemoteMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(RMNRemoteBin), backend, localChainSelector, legacyRMN)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &RMNRemote{address: address, abi: *parsed, RMNRemoteCaller: RMNRemoteCaller{contract: contract}, RMNRemoteTransactor: RMNRemoteTransactor{contract: contract}, RMNRemoteFilterer: RMNRemoteFilterer{contract: contract}}, nil
}

type RMNRemote struct {
	address common.Address
	abi     abi.ABI
	RMNRemoteCaller
	RMNRemoteTransactor
	RMNRemoteFilterer
}

type RMNRemoteCaller struct {
	contract *bind.BoundContract
}

type RMNRemoteTransactor struct {
	contract *bind.BoundContract
}

type RMNRemoteFilterer struct {
	contract *bind.BoundContract
}

type RMNRemoteSession struct {
	Contract     *RMNRemote
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type RMNRemoteCallerSession struct {
	Contract *RMNRemoteCaller
	CallOpts bind.CallOpts
}

type RMNRemoteTransactorSession struct {
	Contract     *RMNRemoteTransactor
	TransactOpts bind.TransactOpts
}

type RMNRemoteRaw struct {
	Contract *RMNRemote
}

type RMNRemoteCallerRaw struct {
	Contract *RMNRemoteCaller
}

type RMNRemoteTransactorRaw struct {
	Contract *RMNRemoteTransactor
}

func NewRMNRemote(address common.Address, backend bind.ContractBackend) (*RMNRemote, error) {
	abi, err := abi.JSON(strings.NewReader(RMNRemoteABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindRMNRemote(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RMNRemote{address: address, abi: abi, RMNRemoteCaller: RMNRemoteCaller{contract: contract}, RMNRemoteTransactor: RMNRemoteTransactor{contract: contract}, RMNRemoteFilterer: RMNRemoteFilterer{contract: contract}}, nil
}

func NewRMNRemoteCaller(address common.Address, caller bind.ContractCaller) (*RMNRemoteCaller, error) {
	contract, err := bindRMNRemote(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RMNRemoteCaller{contract: contract}, nil
}

func NewRMNRemoteTransactor(address common.Address, transactor bind.ContractTransactor) (*RMNRemoteTransactor, error) {
	contract, err := bindRMNRemote(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RMNRemoteTransactor{contract: contract}, nil
}

func NewRMNRemoteFilterer(address common.Address, filterer bind.ContractFilterer) (*RMNRemoteFilterer, error) {
	contract, err := bindRMNRemote(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RMNRemoteFilterer{contract: contract}, nil
}

func bindRMNRemote(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RMNRemoteMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_RMNRemote *RMNRemoteRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RMNRemote.Contract.RMNRemoteCaller.contract.Call(opts, result, method, params...)
}

func (_RMNRemote *RMNRemoteRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RMNRemote.Contract.RMNRemoteTransactor.contract.Transfer(opts)
}

func (_RMNRemote *RMNRemoteRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RMNRemote.Contract.RMNRemoteTransactor.contract.Transact(opts, method, params...)
}

func (_RMNRemote *RMNRemoteCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RMNRemote.Contract.contract.Call(opts, result, method, params...)
}

func (_RMNRemote *RMNRemoteTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RMNRemote.Contract.contract.Transfer(opts)
}

func (_RMNRemote *RMNRemoteTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RMNRemote.Contract.contract.Transact(opts, method, params...)
}

func (_RMNRemote *RMNRemoteCaller) GetCursedSubjects(opts *bind.CallOpts) ([][16]byte, error) {
	var out []interface{}
	err := _RMNRemote.contract.Call(opts, &out, "getCursedSubjects")

	if err != nil {
		return *new([][16]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][16]byte)).(*[][16]byte)

	return out0, err

}

func (_RMNRemote *RMNRemoteSession) GetCursedSubjects() ([][16]byte, error) {
	return _RMNRemote.Contract.GetCursedSubjects(&_RMNRemote.CallOpts)
}

func (_RMNRemote *RMNRemoteCallerSession) GetCursedSubjects() ([][16]byte, error) {
	return _RMNRemote.Contract.GetCursedSubjects(&_RMNRemote.CallOpts)
}

func (_RMNRemote *RMNRemoteCaller) GetLocalChainSelector(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _RMNRemote.contract.Call(opts, &out, "getLocalChainSelector")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_RMNRemote *RMNRemoteSession) GetLocalChainSelector() (uint64, error) {
	return _RMNRemote.Contract.GetLocalChainSelector(&_RMNRemote.CallOpts)
}

func (_RMNRemote *RMNRemoteCallerSession) GetLocalChainSelector() (uint64, error) {
	return _RMNRemote.Contract.GetLocalChainSelector(&_RMNRemote.CallOpts)
}

func (_RMNRemote *RMNRemoteCaller) GetReportDigestHeader(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _RMNRemote.contract.Call(opts, &out, "getReportDigestHeader")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_RMNRemote *RMNRemoteSession) GetReportDigestHeader() ([32]byte, error) {
	return _RMNRemote.Contract.GetReportDigestHeader(&_RMNRemote.CallOpts)
}

func (_RMNRemote *RMNRemoteCallerSession) GetReportDigestHeader() ([32]byte, error) {
	return _RMNRemote.Contract.GetReportDigestHeader(&_RMNRemote.CallOpts)
}

func (_RMNRemote *RMNRemoteCaller) GetVersionedConfig(opts *bind.CallOpts) (GetVersionedConfig,

	error) {
	var out []interface{}
	err := _RMNRemote.contract.Call(opts, &out, "getVersionedConfig")

	outstruct := new(GetVersionedConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Version = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.Config = *abi.ConvertType(out[1], new(RMNRemoteConfig)).(*RMNRemoteConfig)

	return *outstruct, err

}

func (_RMNRemote *RMNRemoteSession) GetVersionedConfig() (GetVersionedConfig,

	error) {
	return _RMNRemote.Contract.GetVersionedConfig(&_RMNRemote.CallOpts)
}

func (_RMNRemote *RMNRemoteCallerSession) GetVersionedConfig() (GetVersionedConfig,

	error) {
	return _RMNRemote.Contract.GetVersionedConfig(&_RMNRemote.CallOpts)
}

func (_RMNRemote *RMNRemoteCaller) IsBlessed(opts *bind.CallOpts, taggedRoot IRMNTaggedRoot) (bool, error) {
	var out []interface{}
	err := _RMNRemote.contract.Call(opts, &out, "isBlessed", taggedRoot)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_RMNRemote *RMNRemoteSession) IsBlessed(taggedRoot IRMNTaggedRoot) (bool, error) {
	return _RMNRemote.Contract.IsBlessed(&_RMNRemote.CallOpts, taggedRoot)
}

func (_RMNRemote *RMNRemoteCallerSession) IsBlessed(taggedRoot IRMNTaggedRoot) (bool, error) {
	return _RMNRemote.Contract.IsBlessed(&_RMNRemote.CallOpts, taggedRoot)
}

func (_RMNRemote *RMNRemoteCaller) IsCursed(opts *bind.CallOpts, subject [16]byte) (bool, error) {
	var out []interface{}
	err := _RMNRemote.contract.Call(opts, &out, "isCursed", subject)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_RMNRemote *RMNRemoteSession) IsCursed(subject [16]byte) (bool, error) {
	return _RMNRemote.Contract.IsCursed(&_RMNRemote.CallOpts, subject)
}

func (_RMNRemote *RMNRemoteCallerSession) IsCursed(subject [16]byte) (bool, error) {
	return _RMNRemote.Contract.IsCursed(&_RMNRemote.CallOpts, subject)
}

func (_RMNRemote *RMNRemoteCaller) IsCursed0(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _RMNRemote.contract.Call(opts, &out, "isCursed0")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_RMNRemote *RMNRemoteSession) IsCursed0() (bool, error) {
	return _RMNRemote.Contract.IsCursed0(&_RMNRemote.CallOpts)
}

func (_RMNRemote *RMNRemoteCallerSession) IsCursed0() (bool, error) {
	return _RMNRemote.Contract.IsCursed0(&_RMNRemote.CallOpts)
}

func (_RMNRemote *RMNRemoteCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RMNRemote.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_RMNRemote *RMNRemoteSession) Owner() (common.Address, error) {
	return _RMNRemote.Contract.Owner(&_RMNRemote.CallOpts)
}

func (_RMNRemote *RMNRemoteCallerSession) Owner() (common.Address, error) {
	return _RMNRemote.Contract.Owner(&_RMNRemote.CallOpts)
}

func (_RMNRemote *RMNRemoteCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _RMNRemote.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_RMNRemote *RMNRemoteSession) TypeAndVersion() (string, error) {
	return _RMNRemote.Contract.TypeAndVersion(&_RMNRemote.CallOpts)
}

func (_RMNRemote *RMNRemoteCallerSession) TypeAndVersion() (string, error) {
	return _RMNRemote.Contract.TypeAndVersion(&_RMNRemote.CallOpts)
}

func (_RMNRemote *RMNRemoteCaller) Verify(opts *bind.CallOpts, offRampAddress common.Address, merkleRoots []InternalMerkleRoot, signatures []IRMNRemoteSignature) error {
	var out []interface{}
	err := _RMNRemote.contract.Call(opts, &out, "verify", offRampAddress, merkleRoots, signatures)

	if err != nil {
		return err
	}

	return err

}

func (_RMNRemote *RMNRemoteSession) Verify(offRampAddress common.Address, merkleRoots []InternalMerkleRoot, signatures []IRMNRemoteSignature) error {
	return _RMNRemote.Contract.Verify(&_RMNRemote.CallOpts, offRampAddress, merkleRoots, signatures)
}

func (_RMNRemote *RMNRemoteCallerSession) Verify(offRampAddress common.Address, merkleRoots []InternalMerkleRoot, signatures []IRMNRemoteSignature) error {
	return _RMNRemote.Contract.Verify(&_RMNRemote.CallOpts, offRampAddress, merkleRoots, signatures)
}

func (_RMNRemote *RMNRemoteTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RMNRemote.contract.Transact(opts, "acceptOwnership")
}

func (_RMNRemote *RMNRemoteSession) AcceptOwnership() (*types.Transaction, error) {
	return _RMNRemote.Contract.AcceptOwnership(&_RMNRemote.TransactOpts)
}

func (_RMNRemote *RMNRemoteTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _RMNRemote.Contract.AcceptOwnership(&_RMNRemote.TransactOpts)
}

func (_RMNRemote *RMNRemoteTransactor) Curse(opts *bind.TransactOpts, subject [16]byte) (*types.Transaction, error) {
	return _RMNRemote.contract.Transact(opts, "curse", subject)
}

func (_RMNRemote *RMNRemoteSession) Curse(subject [16]byte) (*types.Transaction, error) {
	return _RMNRemote.Contract.Curse(&_RMNRemote.TransactOpts, subject)
}

func (_RMNRemote *RMNRemoteTransactorSession) Curse(subject [16]byte) (*types.Transaction, error) {
	return _RMNRemote.Contract.Curse(&_RMNRemote.TransactOpts, subject)
}

func (_RMNRemote *RMNRemoteTransactor) Curse0(opts *bind.TransactOpts, subjects [][16]byte) (*types.Transaction, error) {
	return _RMNRemote.contract.Transact(opts, "curse0", subjects)
}

func (_RMNRemote *RMNRemoteSession) Curse0(subjects [][16]byte) (*types.Transaction, error) {
	return _RMNRemote.Contract.Curse0(&_RMNRemote.TransactOpts, subjects)
}

func (_RMNRemote *RMNRemoteTransactorSession) Curse0(subjects [][16]byte) (*types.Transaction, error) {
	return _RMNRemote.Contract.Curse0(&_RMNRemote.TransactOpts, subjects)
}

func (_RMNRemote *RMNRemoteTransactor) SetConfig(opts *bind.TransactOpts, newConfig RMNRemoteConfig) (*types.Transaction, error) {
	return _RMNRemote.contract.Transact(opts, "setConfig", newConfig)
}

func (_RMNRemote *RMNRemoteSession) SetConfig(newConfig RMNRemoteConfig) (*types.Transaction, error) {
	return _RMNRemote.Contract.SetConfig(&_RMNRemote.TransactOpts, newConfig)
}

func (_RMNRemote *RMNRemoteTransactorSession) SetConfig(newConfig RMNRemoteConfig) (*types.Transaction, error) {
	return _RMNRemote.Contract.SetConfig(&_RMNRemote.TransactOpts, newConfig)
}

func (_RMNRemote *RMNRemoteTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _RMNRemote.contract.Transact(opts, "transferOwnership", to)
}

func (_RMNRemote *RMNRemoteSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _RMNRemote.Contract.TransferOwnership(&_RMNRemote.TransactOpts, to)
}

func (_RMNRemote *RMNRemoteTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _RMNRemote.Contract.TransferOwnership(&_RMNRemote.TransactOpts, to)
}

func (_RMNRemote *RMNRemoteTransactor) Uncurse(opts *bind.TransactOpts, subject [16]byte) (*types.Transaction, error) {
	return _RMNRemote.contract.Transact(opts, "uncurse", subject)
}

func (_RMNRemote *RMNRemoteSession) Uncurse(subject [16]byte) (*types.Transaction, error) {
	return _RMNRemote.Contract.Uncurse(&_RMNRemote.TransactOpts, subject)
}

func (_RMNRemote *RMNRemoteTransactorSession) Uncurse(subject [16]byte) (*types.Transaction, error) {
	return _RMNRemote.Contract.Uncurse(&_RMNRemote.TransactOpts, subject)
}

func (_RMNRemote *RMNRemoteTransactor) Uncurse0(opts *bind.TransactOpts, subjects [][16]byte) (*types.Transaction, error) {
	return _RMNRemote.contract.Transact(opts, "uncurse0", subjects)
}

func (_RMNRemote *RMNRemoteSession) Uncurse0(subjects [][16]byte) (*types.Transaction, error) {
	return _RMNRemote.Contract.Uncurse0(&_RMNRemote.TransactOpts, subjects)
}

func (_RMNRemote *RMNRemoteTransactorSession) Uncurse0(subjects [][16]byte) (*types.Transaction, error) {
	return _RMNRemote.Contract.Uncurse0(&_RMNRemote.TransactOpts, subjects)
}

type RMNRemoteConfigSetIterator struct {
	Event *RMNRemoteConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RMNRemoteConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMNRemoteConfigSet)
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
		it.Event = new(RMNRemoteConfigSet)
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

func (it *RMNRemoteConfigSetIterator) Error() error {
	return it.fail
}

func (it *RMNRemoteConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RMNRemoteConfigSet struct {
	Version uint32
	Config  RMNRemoteConfig
	Raw     types.Log
}

func (_RMNRemote *RMNRemoteFilterer) FilterConfigSet(opts *bind.FilterOpts, version []uint32) (*RMNRemoteConfigSetIterator, error) {

	var versionRule []interface{}
	for _, versionItem := range version {
		versionRule = append(versionRule, versionItem)
	}

	logs, sub, err := _RMNRemote.contract.FilterLogs(opts, "ConfigSet", versionRule)
	if err != nil {
		return nil, err
	}
	return &RMNRemoteConfigSetIterator{contract: _RMNRemote.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_RMNRemote *RMNRemoteFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *RMNRemoteConfigSet, version []uint32) (event.Subscription, error) {

	var versionRule []interface{}
	for _, versionItem := range version {
		versionRule = append(versionRule, versionItem)
	}

	logs, sub, err := _RMNRemote.contract.WatchLogs(opts, "ConfigSet", versionRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RMNRemoteConfigSet)
				if err := _RMNRemote.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_RMNRemote *RMNRemoteFilterer) ParseConfigSet(log types.Log) (*RMNRemoteConfigSet, error) {
	event := new(RMNRemoteConfigSet)
	if err := _RMNRemote.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RMNRemoteCursedIterator struct {
	Event *RMNRemoteCursed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RMNRemoteCursedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMNRemoteCursed)
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
		it.Event = new(RMNRemoteCursed)
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

func (it *RMNRemoteCursedIterator) Error() error {
	return it.fail
}

func (it *RMNRemoteCursedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RMNRemoteCursed struct {
	Subjects [][16]byte
	Raw      types.Log
}

func (_RMNRemote *RMNRemoteFilterer) FilterCursed(opts *bind.FilterOpts) (*RMNRemoteCursedIterator, error) {

	logs, sub, err := _RMNRemote.contract.FilterLogs(opts, "Cursed")
	if err != nil {
		return nil, err
	}
	return &RMNRemoteCursedIterator{contract: _RMNRemote.contract, event: "Cursed", logs: logs, sub: sub}, nil
}

func (_RMNRemote *RMNRemoteFilterer) WatchCursed(opts *bind.WatchOpts, sink chan<- *RMNRemoteCursed) (event.Subscription, error) {

	logs, sub, err := _RMNRemote.contract.WatchLogs(opts, "Cursed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RMNRemoteCursed)
				if err := _RMNRemote.contract.UnpackLog(event, "Cursed", log); err != nil {
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

func (_RMNRemote *RMNRemoteFilterer) ParseCursed(log types.Log) (*RMNRemoteCursed, error) {
	event := new(RMNRemoteCursed)
	if err := _RMNRemote.contract.UnpackLog(event, "Cursed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RMNRemoteOwnershipTransferRequestedIterator struct {
	Event *RMNRemoteOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RMNRemoteOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMNRemoteOwnershipTransferRequested)
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
		it.Event = new(RMNRemoteOwnershipTransferRequested)
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

func (it *RMNRemoteOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *RMNRemoteOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RMNRemoteOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_RMNRemote *RMNRemoteFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RMNRemoteOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _RMNRemote.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &RMNRemoteOwnershipTransferRequestedIterator{contract: _RMNRemote.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_RMNRemote *RMNRemoteFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *RMNRemoteOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _RMNRemote.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RMNRemoteOwnershipTransferRequested)
				if err := _RMNRemote.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_RMNRemote *RMNRemoteFilterer) ParseOwnershipTransferRequested(log types.Log) (*RMNRemoteOwnershipTransferRequested, error) {
	event := new(RMNRemoteOwnershipTransferRequested)
	if err := _RMNRemote.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RMNRemoteOwnershipTransferredIterator struct {
	Event *RMNRemoteOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RMNRemoteOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMNRemoteOwnershipTransferred)
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
		it.Event = new(RMNRemoteOwnershipTransferred)
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

func (it *RMNRemoteOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *RMNRemoteOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RMNRemoteOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_RMNRemote *RMNRemoteFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RMNRemoteOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _RMNRemote.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &RMNRemoteOwnershipTransferredIterator{contract: _RMNRemote.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_RMNRemote *RMNRemoteFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *RMNRemoteOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _RMNRemote.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RMNRemoteOwnershipTransferred)
				if err := _RMNRemote.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_RMNRemote *RMNRemoteFilterer) ParseOwnershipTransferred(log types.Log) (*RMNRemoteOwnershipTransferred, error) {
	event := new(RMNRemoteOwnershipTransferred)
	if err := _RMNRemote.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RMNRemoteUncursedIterator struct {
	Event *RMNRemoteUncursed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RMNRemoteUncursedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMNRemoteUncursed)
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
		it.Event = new(RMNRemoteUncursed)
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

func (it *RMNRemoteUncursedIterator) Error() error {
	return it.fail
}

func (it *RMNRemoteUncursedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RMNRemoteUncursed struct {
	Subjects [][16]byte
	Raw      types.Log
}

func (_RMNRemote *RMNRemoteFilterer) FilterUncursed(opts *bind.FilterOpts) (*RMNRemoteUncursedIterator, error) {

	logs, sub, err := _RMNRemote.contract.FilterLogs(opts, "Uncursed")
	if err != nil {
		return nil, err
	}
	return &RMNRemoteUncursedIterator{contract: _RMNRemote.contract, event: "Uncursed", logs: logs, sub: sub}, nil
}

func (_RMNRemote *RMNRemoteFilterer) WatchUncursed(opts *bind.WatchOpts, sink chan<- *RMNRemoteUncursed) (event.Subscription, error) {

	logs, sub, err := _RMNRemote.contract.WatchLogs(opts, "Uncursed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RMNRemoteUncursed)
				if err := _RMNRemote.contract.UnpackLog(event, "Uncursed", log); err != nil {
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

func (_RMNRemote *RMNRemoteFilterer) ParseUncursed(log types.Log) (*RMNRemoteUncursed, error) {
	event := new(RMNRemoteUncursed)
	if err := _RMNRemote.contract.UnpackLog(event, "Uncursed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetVersionedConfig struct {
	Version uint32
	Config  RMNRemoteConfig
}

func (_RMNRemote *RMNRemote) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _RMNRemote.abi.Events["ConfigSet"].ID:
		return _RMNRemote.ParseConfigSet(log)
	case _RMNRemote.abi.Events["Cursed"].ID:
		return _RMNRemote.ParseCursed(log)
	case _RMNRemote.abi.Events["OwnershipTransferRequested"].ID:
		return _RMNRemote.ParseOwnershipTransferRequested(log)
	case _RMNRemote.abi.Events["OwnershipTransferred"].ID:
		return _RMNRemote.ParseOwnershipTransferred(log)
	case _RMNRemote.abi.Events["Uncursed"].ID:
		return _RMNRemote.ParseUncursed(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (RMNRemoteConfigSet) Topic() common.Hash {
	return common.HexToHash("0x7f22bf988149dbe8de8fb879c6b97a4e56e68b2bd57421ce1a4e79d4ef6b496c")
}

func (RMNRemoteCursed) Topic() common.Hash {
	return common.HexToHash("0x1716e663a90a76d3b6c7e5f680673d1b051454c19c627e184c8daf28f3104f74")
}

func (RMNRemoteOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (RMNRemoteOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (RMNRemoteUncursed) Topic() common.Hash {
	return common.HexToHash("0x0676e709c9cc74fa0519fd78f7c33be0f1b2b0bae0507c724aef7229379c6ba1")
}

func (_RMNRemote *RMNRemote) Address() common.Address {
	return _RMNRemote.address
}

type RMNRemoteInterface interface {
	GetCursedSubjects(opts *bind.CallOpts) ([][16]byte, error)

	GetLocalChainSelector(opts *bind.CallOpts) (uint64, error)

	GetReportDigestHeader(opts *bind.CallOpts) ([32]byte, error)

	GetVersionedConfig(opts *bind.CallOpts) (GetVersionedConfig,

		error)

	IsBlessed(opts *bind.CallOpts, taggedRoot IRMNTaggedRoot) (bool, error)

	IsCursed(opts *bind.CallOpts, subject [16]byte) (bool, error)

	IsCursed0(opts *bind.CallOpts) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	Verify(opts *bind.CallOpts, offRampAddress common.Address, merkleRoots []InternalMerkleRoot, signatures []IRMNRemoteSignature) error

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	Curse(opts *bind.TransactOpts, subject [16]byte) (*types.Transaction, error)

	Curse0(opts *bind.TransactOpts, subjects [][16]byte) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, newConfig RMNRemoteConfig) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Uncurse(opts *bind.TransactOpts, subject [16]byte) (*types.Transaction, error)

	Uncurse0(opts *bind.TransactOpts, subjects [][16]byte) (*types.Transaction, error)

	FilterConfigSet(opts *bind.FilterOpts, version []uint32) (*RMNRemoteConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *RMNRemoteConfigSet, version []uint32) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*RMNRemoteConfigSet, error)

	FilterCursed(opts *bind.FilterOpts) (*RMNRemoteCursedIterator, error)

	WatchCursed(opts *bind.WatchOpts, sink chan<- *RMNRemoteCursed) (event.Subscription, error)

	ParseCursed(log types.Log) (*RMNRemoteCursed, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RMNRemoteOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *RMNRemoteOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*RMNRemoteOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RMNRemoteOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *RMNRemoteOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*RMNRemoteOwnershipTransferred, error)

	FilterUncursed(opts *bind.FilterOpts) (*RMNRemoteUncursedIterator, error)

	WatchUncursed(opts *bind.WatchOpts, sink chan<- *RMNRemoteUncursed) (event.Subscription, error)

	ParseUncursed(log types.Log) (*RMNRemoteUncursed, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
