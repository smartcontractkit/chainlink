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

type IRMNV2Signature struct {
	R [32]byte
	S [32]byte
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
	MinSigners                  uint64
}

type RMNRemoteSigner struct {
	OnchainPublicKey common.Address
	NodeIndex        uint64
}

var RMNRemoteMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"localChainSelector\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"subject\",\"type\":\"bytes16\"}],\"name\":\"AlreadyCursed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConfigNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateOnchainPublicKey\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSignature\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSignerOrder\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MinSignersTooHigh\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"subject\",\"type\":\"bytes16\"}],\"name\":\"NotCursed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OutOfOrderSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ThresholdNotMet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnexpectedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroValueNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"rmnHomeContractConfigDigest\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"onchainPublicKey\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"nodeIndex\",\"type\":\"uint64\"}],\"internalType\":\"structRMNRemote.Signer[]\",\"name\":\"signers\",\"type\":\"tuple[]\"},{\"internalType\":\"uint64\",\"name\":\"minSigners\",\"type\":\"uint64\"}],\"indexed\":false,\"internalType\":\"structRMNRemote.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16[]\",\"name\":\"subjects\",\"type\":\"bytes16[]\"}],\"name\":\"Cursed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16[]\",\"name\":\"subjects\",\"type\":\"bytes16[]\"}],\"name\":\"Uncursed\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"subject\",\"type\":\"bytes16\"}],\"name\":\"curse\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16[]\",\"name\":\"subjects\",\"type\":\"bytes16[]\"}],\"name\":\"curse\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCursedSubjects\",\"outputs\":[{\"internalType\":\"bytes16[]\",\"name\":\"subjects\",\"type\":\"bytes16[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLocalChainSelector\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"localChainSelector\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getVersionedConfig\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"rmnHomeContractConfigDigest\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"onchainPublicKey\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"nodeIndex\",\"type\":\"uint64\"}],\"internalType\":\"structRMNRemote.Signer[]\",\"name\":\"signers\",\"type\":\"tuple[]\"},{\"internalType\":\"uint64\",\"name\":\"minSigners\",\"type\":\"uint64\"}],\"internalType\":\"structRMNRemote.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"subject\",\"type\":\"bytes16\"}],\"name\":\"isCursed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isCursed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"rmnHomeContractConfigDigest\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"onchainPublicKey\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"nodeIndex\",\"type\":\"uint64\"}],\"internalType\":\"structRMNRemote.Signer[]\",\"name\":\"signers\",\"type\":\"tuple[]\"},{\"internalType\":\"uint64\",\"name\":\"minSigners\",\"type\":\"uint64\"}],\"internalType\":\"structRMNRemote.Config\",\"name\":\"newConfig\",\"type\":\"tuple\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"subject\",\"type\":\"bytes16\"}],\"name\":\"uncurse\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16[]\",\"name\":\"subjects\",\"type\":\"bytes16[]\"}],\"name\":\"uncurse\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"offrampAddress\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRampAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"maxSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"merkleRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structInternal.MerkleRoot[]\",\"name\":\"merkleRoots\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structIRMNV2.Signature[]\",\"name\":\"signatures\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"rawVs\",\"type\":\"uint256\"}],\"name\":\"verify\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b5060405162002131380380620021318339810160408190526200003491620001a9565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000fe565b505050806001600160401b0316600003620000ec5760405163273e150360e21b815260040160405180910390fd5b6001600160401b0316608052620001db565b336001600160a01b03821603620001585760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600060208284031215620001bc57600080fd5b81516001600160401b0381168114620001d457600080fd5b9392505050565b608051611f33620001fe600039600081816102410152610a920152611f336000f3fe608060405234801561001057600080fd5b50600436106100f55760003560e01c806379ba509711610097578063d881e09211610066578063d881e0921461021e578063eaa83ddd14610233578063f2fde38b1461026b578063f8bb876e1461027e57600080fd5b806379ba5097146101c85780638d8741cb146101d05780638da5cb5b146101e35780639a19b3291461020b57600080fd5b80632cbc26bb116100d35780632cbc26bb14610177578063397796f71461019a57806362eed415146101a25780636d2d3993146101b557600080fd5b8063181f5a77146100fa578063198f0f771461014c5780631add205f14610161575b600080fd5b6101366040518060400160405280601381526020017f524d4e52656d6f746520312e362e302d6465760000000000000000000000000081525081565b6040516101439190611380565b60405180910390f35b61015f61015a366004611393565b610291565b005b61016961063d565b6040516101439291906113ce565b61018a6101853660046114ac565b610735565b6040519015158152602001610143565b61018a610792565b61015f6101b03660046114ac565b61080c565b61015f6101c33660046114ac565b610880565b61015f6108f0565b61015f6101de366004611535565b6109f2565b60005460405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610143565b61015f6102193660046116bc565b610d46565b610226610e43565b6040516101439190611759565b60405167ffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000168152602001610143565b61015f6102793660046117bf565b610e4f565b61015f61028c3660046116bc565b610e63565b610299610f55565b60015b6102a960208301836117dc565b9050811015610379576102bf60208301836117dc565b828181106102cf576102cf611844565b90506040020160200160208101906102e79190611894565b67ffffffffffffffff166102fe60208401846117dc565b6103096001856118e0565b81811061031857610318611844565b90506040020160200160208101906103309190611894565b67ffffffffffffffff1610610371576040517f4485151700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60010161029c565b5061038760208201826117dc565b90506103996060830160408401611894565b67ffffffffffffffff1611156103db576040517ffba0d9e600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6003545b801561046d576008600060036103f66001856118e0565b8154811061040657610406611844565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff168352820192909252604001902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055610466816118f3565b90506103df565b5060005b61047e60208301836117dc565b90508110156105b3576008600061049860208501856117dc565b848181106104a8576104a8611844565b6104be92602060409092020190810191506117bf565b73ffffffffffffffffffffffffffffffffffffffff16815260208101919091526040016000205460ff161561051f576040517f28cae27d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60016008600061053260208601866117dc565b8581811061054257610542611844565b61055892602060409092020190810191506117bf565b73ffffffffffffffffffffffffffffffffffffffff168152602081019190915260400160002080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055600101610471565b508060026105c182826119e1565b5050600580546000919082906105dc9063ffffffff16611b1c565b91906101000a81548163ffffffff021916908363ffffffff160217905590508063ffffffff167f7f22bf988149dbe8de8fb879c6b97a4e56e68b2bd57421ce1a4e79d4ef6b496c836040516106319190611b3f565b60405180910390a25050565b6040805160608082018352600080835260208301919091529181018290526005546040805160608101825260028054825260038054845160208281028201810190965281815263ffffffff9096169592948593818601939092909160009084015b8282101561070c576000848152602090819020604080518082019091529084015473ffffffffffffffffffffffffffffffffffffffff8116825274010000000000000000000000000000000000000000900467ffffffffffffffff168183015282526001909201910161069e565b505050908252506002919091015467ffffffffffffffff16602090910152919491935090915050565b60006107416006610fd8565b60000361075057506000919050565b61075b600683610fe2565b8061078c575061078c60067f0100000000000000000000000000000100000000000000000000000000000000610fe2565b92915050565b600061079e6006610fd8565b6000036107ab5750600090565b6107d660067f0100000000000000000000000000000000000000000000000000000000000000610fe2565b80610807575061080760067f0100000000000000000000000000000100000000000000000000000000000000610fe2565b905090565b60408051600180825281830190925260009160208083019080368337019050509050818160008151811061084257610842611844565b7fffffffffffffffffffffffffffffffff000000000000000000000000000000009092166020928302919091019091015261087c81610e63565b5050565b6040805160018082528183019092526000916020808301908036833701905050905081816000815181106108b6576108b6611844565b7fffffffffffffffffffffffffffffffff000000000000000000000000000000009092166020928302919091019091015261087c81610d46565b60015473ffffffffffffffffffffffffffffffffffffffff163314610976576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60055463ffffffff16600003610a34576040517face124bc00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60045467ffffffffffffffff16821015610a7a576040517f59fa4a9300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805160c08101825246815267ffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000166020820152309181019190915273ffffffffffffffffffffffffffffffffffffffff8716606082015260025460808201526000907f9651943783dbf81935a60e98f218a9d9b5b28823fb2228bbd91320d632facf539060a08101610b16888a611c49565b9052604051610b29929190602001611da9565b60405160208183030381529060405280519060200120905060008060005b85811015610d3a57600184610b6182841b8816601b611ede565b898985818110610b7357610b73611844565b905060400201600001358a8a86818110610b8f57610b8f611844565b9050604002016020013560405160008152602001604052604051610bcf949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa158015610bf1573d6000803e3d6000fd5b50506040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0015192505073ffffffffffffffffffffffffffffffffffffffff8216610c69576040517f8baa579f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1610610cce576040517fbbe15e7f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff821660009081526008602052604090205460ff16610d2d576040517faaaa914100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b9091508190600101610b47565b50505050505050505050565b610d4e610f55565b60005b8151811015610e0857610d87828281518110610d6f57610d6f611844565b6020026020010151600661102090919063ffffffff16565b610e0057818181518110610d9d57610d9d611844565b60200260200101516040517f73281fa100000000000000000000000000000000000000000000000000000000815260040161096d91907fffffffffffffffffffffffffffffffff0000000000000000000000000000000091909116815260200190565b600101610d51565b507f0676e709c9cc74fa0519fd78f7c33be0f1b2b0bae0507c724aef7229379c6ba181604051610e389190611759565b60405180910390a150565b6060610807600661104e565b610e57610f55565b610e608161105b565b50565b610e6b610f55565b60005b8151811015610f2557610ea4828281518110610e8c57610e8c611844565b6020026020010151600661115090919063ffffffff16565b610f1d57818181518110610eba57610eba611844565b60200260200101516040517f19d5c79b00000000000000000000000000000000000000000000000000000000815260040161096d91907fffffffffffffffffffffffffffffffff0000000000000000000000000000000091909116815260200190565b600101610e6e565b507f1716e663a90a76d3b6c7e5f680673d1b051454c19c627e184c8daf28f3104f7481604051610e389190611759565b60005473ffffffffffffffffffffffffffffffffffffffff163314610fd6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161096d565b565b600061078c825490565b7fffffffffffffffffffffffffffffffff000000000000000000000000000000008116600090815260018301602052604081205415155b9392505050565b6000611019837fffffffffffffffffffffffffffffffff00000000000000000000000000000000841661117e565b6060600061101983611271565b3373ffffffffffffffffffffffffffffffffffffffff8216036110da576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161096d565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000611019837fffffffffffffffffffffffffffffffff0000000000000000000000000000000084166112cd565b600081815260018301602052604081205480156112675760006111a26001836118e0565b85549091506000906111b6906001906118e0565b905080821461121b5760008660000182815481106111d6576111d6611844565b90600052602060002001549050808760000184815481106111f9576111f9611844565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061122c5761122c611ef7565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505061078c565b600091505061078c565b6060816000018054806020026020016040519081016040528092919081815260200182805480156112c157602002820191906000526020600020905b8154815260200190600101908083116112ad575b50505050509050919050565b60008181526001830160205260408120546113145750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915561078c565b50600061078c565b6000815180845260005b8181101561134257602081850181015186830182015201611326565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b602081526000611019602083018461131c565b6000602082840312156113a557600080fd5b813567ffffffffffffffff8111156113bc57600080fd5b82016060818503121561101957600080fd5b63ffffffff831681526040602080830182905283518383015283810151606080850152805160a085018190526000939291820190849060c08701905b80831015611453578351805173ffffffffffffffffffffffffffffffffffffffff16835285015167ffffffffffffffff168583015292840192600192909201919085019061140a565b50604088015167ffffffffffffffff81166080890152945098975050505050505050565b80357fffffffffffffffffffffffffffffffff00000000000000000000000000000000811681146114a757600080fd5b919050565b6000602082840312156114be57600080fd5b61101982611477565b73ffffffffffffffffffffffffffffffffffffffff81168114610e6057600080fd5b60008083601f8401126114fb57600080fd5b50813567ffffffffffffffff81111561151357600080fd5b6020830191508360208260061b850101111561152e57600080fd5b9250929050565b6000806000806000806080878903121561154e57600080fd5b8635611559816114c7565b9550602087013567ffffffffffffffff8082111561157657600080fd5b818901915089601f83011261158a57600080fd5b81358181111561159957600080fd5b8a60208260051b85010111156115ae57600080fd5b6020830197508096505060408901359150808211156115cc57600080fd5b506115d989828a016114e9565b979a9699509497949695606090950135949350505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160a0810167ffffffffffffffff81118282101715611643576116436115f1565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715611690576116906115f1565b604052919050565b600067ffffffffffffffff8211156116b2576116b26115f1565b5060051b60200190565b600060208083850312156116cf57600080fd5b823567ffffffffffffffff8111156116e657600080fd5b8301601f810185136116f757600080fd5b803561170a61170582611698565b611649565b81815260059190911b8201830190838101908783111561172957600080fd5b928401925b8284101561174e5761173f84611477565b8252928401929084019061172e565b979650505050505050565b6020808252825182820181905260009190848201906040850190845b818110156117b35783517fffffffffffffffffffffffffffffffff000000000000000000000000000000001683529284019291840191600101611775565b50909695505050505050565b6000602082840312156117d157600080fd5b8135611019816114c7565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261181157600080fd5b83018035915067ffffffffffffffff82111561182c57600080fd5b6020019150600681901b360382131561152e57600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b67ffffffffffffffff81168114610e6057600080fd5b80356114a781611873565b6000602082840312156118a657600080fd5b813561101981611873565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8181038181111561078c5761078c6118b1565b600081611902576119026118b1565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190565b6000813561078c81611873565b8135611940816114c7565b73ffffffffffffffffffffffffffffffffffffffff811690508154817fffffffffffffffffffffffff00000000000000000000000000000000000000008216178355602084013561199081611873565b7bffffffffffffffff00000000000000000000000000000000000000008160a01b16837fffffffff000000000000000000000000000000000000000000000000000000008416171784555050505050565b81358155600180820160208401357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1853603018112611a1f57600080fd5b8401803567ffffffffffffffff811115611a3857600080fd5b6020820191508060061b3603821315611a5057600080fd5b68010000000000000000811115611a6957611a696115f1565b825481845580821015611a9e576000848152602081208381019083015b80821015611a9a5782825590870190611a86565b5050505b50600092835260208320925b81811015611ace57611abc8385611935565b92840192604092909201918401611aaa565b505050505061087c611ae260408401611928565b6002830167ffffffffffffffff82167fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000008254161781555050565b600063ffffffff808316818103611b3557611b356118b1565b6001019392505050565b6000602080835260808301843582850152818501357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1863603018112611b8457600080fd5b8501828101903567ffffffffffffffff80821115611ba157600080fd5b8160061b3603831315611bb357600080fd5b6040606060408901528483865260a089019050849550600094505b83851015611c1e578535611be1816114c7565b73ffffffffffffffffffffffffffffffffffffffff16815285870135611c0681611873565b83168188015294810194600194909401938101611bce565b611c2a60408b01611889565b67ffffffffffffffff811660608b015296509998505050505050505050565b6000611c5761170584611698565b80848252602080830192508560051b850136811115611c7557600080fd5b855b81811015611d9d57803567ffffffffffffffff80821115611c985760008081fd5b818901915060a08236031215611cae5760008081fd5b611cb6611620565b8235611cc181611873565b81528286013582811115611cd55760008081fd5b8301601f3681830112611ce85760008081fd5b813584811115611cfa57611cfa6115f1565b611d29897fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08484011601611649565b94508085523689828501011115611d4257600091508182fd5b808984018a8701376000898287010152505050818682015260409150611d69828401611889565b8282015260609150611d7c828401611889565b91810191909152608091820135918101919091528552938201938201611c77565b50919695505050505050565b60006040848352602060408185015261010084018551604086015281860151606067ffffffffffffffff808316606089015260408901519250608073ffffffffffffffffffffffffffffffffffffffff80851660808b015260608b0151945060a081861660a08c015260808c015160c08c015260a08c0151955060c060e08c015286915085518088526101209750878c019250878160051b8d01019750888701965060005b81811015611ecb577ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffee08d8a030184528751868151168a528a810151848c8c0152611e9a858c018261131c565b828e015189168c8f01528983015189168a8d0152918701519a87019a909a5298509689019692890192600101611e4e565b50969d9c50505050505050505050505050565b60ff818116838216019081111561078c5761078c6118b1565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea164736f6c6343000818000a",
}

var RMNRemoteABI = RMNRemoteMetaData.ABI

var RMNRemoteBin = RMNRemoteMetaData.Bin

func DeployRMNRemote(auth *bind.TransactOpts, backend bind.ContractBackend, localChainSelector uint64) (common.Address, *types.Transaction, *RMNRemote, error) {
	parsed, err := RMNRemoteMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(RMNRemoteBin), backend, localChainSelector)
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

func (_RMNRemote *RMNRemoteCaller) Verify(opts *bind.CallOpts, offrampAddress common.Address, merkleRoots []InternalMerkleRoot, signatures []IRMNV2Signature, rawVs *big.Int) error {
	var out []interface{}
	err := _RMNRemote.contract.Call(opts, &out, "verify", offrampAddress, merkleRoots, signatures, rawVs)

	if err != nil {
		return err
	}

	return err

}

func (_RMNRemote *RMNRemoteSession) Verify(offrampAddress common.Address, merkleRoots []InternalMerkleRoot, signatures []IRMNV2Signature, rawVs *big.Int) error {
	return _RMNRemote.Contract.Verify(&_RMNRemote.CallOpts, offrampAddress, merkleRoots, signatures, rawVs)
}

func (_RMNRemote *RMNRemoteCallerSession) Verify(offrampAddress common.Address, merkleRoots []InternalMerkleRoot, signatures []IRMNV2Signature, rawVs *big.Int) error {
	return _RMNRemote.Contract.Verify(&_RMNRemote.CallOpts, offrampAddress, merkleRoots, signatures, rawVs)
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

	GetVersionedConfig(opts *bind.CallOpts) (GetVersionedConfig,

		error)

	IsCursed(opts *bind.CallOpts, subject [16]byte) (bool, error)

	IsCursed0(opts *bind.CallOpts) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	Verify(opts *bind.CallOpts, offrampAddress common.Address, merkleRoots []InternalMerkleRoot, signatures []IRMNV2Signature, rawVs *big.Int) error

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
