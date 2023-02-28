// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package exchanger

import (
	"errors"
	"math/big"
	"strings"

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
)

// ExchangerMetaData contains all meta data concerning the Exchanger contract.
var ExchangerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIVerifierProxy\",\"name\":\"verifierProxyAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"lookupURL\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"maxDelay\",\"type\":\"uint8\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"reportFeedID\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"commitmentFeedID\",\"type\":\"bytes32\"}],\"name\":\"FeedIDMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"string[]\",\"name\":\"urls\",\"type\":\"string[]\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"},{\"internalType\":\"bytes4\",\"name\":\"callbackFunction\",\"type\":\"bytes4\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"OffchainLookup\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"reportBlocknumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"commitmentBlocknumber\",\"type\":\"uint256\"}],\"name\":\"ReportBlockMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"blocknumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tradeWindow\",\"type\":\"uint256\"}],\"name\":\"TradeExceedsWindow\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"maxDelay\",\"type\":\"uint8\"}],\"name\":\"SetDelay\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"url\",\"type\":\"string\"}],\"name\":\"SetLookupURL\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractIVerifierProxy\",\"name\":\"verifierProxyAddress\",\"type\":\"address\"}],\"name\":\"SetVerifierProxy\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"commitment\",\"type\":\"bytes32\"}],\"name\":\"TradeCommitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"currencySrc\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"currencyDst\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountSrc\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"minAmountDst\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"int192\",\"name\":\"median\",\"type\":\"int192\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"observationsTimestamp\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"observationsBlocknumber\",\"type\":\"uint64\"}],\"name\":\"TradeExecuted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"commitment\",\"type\":\"bytes32\"}],\"name\":\"commitTrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDelay\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"maxDelay\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLookupURL\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"url\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getVerifierProxyAddress\",\"outputs\":[{\"internalType\":\"contractIVerifierProxy\",\"name\":\"verifierProxyAddress\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedCommitment\",\"type\":\"bytes\"}],\"name\":\"resolveTrade\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"chainlinkBlob\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"encodedCommitment\",\"type\":\"bytes\"}],\"name\":\"resolveTradeWithReport\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"maxDelay\",\"type\":\"uint8\"}],\"name\":\"setDelay\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"url\",\"type\":\"string\"}],\"name\":\"setLookupURL\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIVerifierProxy\",\"name\":\"verifierProxyAddress\",\"type\":\"address\"}],\"name\":\"setVerifierProxyAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b50604051620017b5380380620017b583398101604081905262000034916200029c565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be8162000133565b5050600280546001600160a01b0319166001600160a01b03861617905550604051620000ef90839060200162000389565b6040516020818303038152906040526003908051906020019062000115929190620001df565b506004805460ff191660ff9290921691909117905550620004449050565b6001600160a01b0381163314156200018e5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b828054620001ed90620003f1565b90600052602060002090601f0160209004810192826200021157600085556200025c565b82601f106200022c57805160ff19168380011785556200025c565b828001600101855582156200025c579182015b828111156200025c5782518255916020019190600101906200023f565b506200026a9291506200026e565b5090565b5b808211156200026a57600081556001016200026f565b805160ff811681146200029757600080fd5b919050565b600080600060608486031215620002b257600080fd5b83516001600160a01b0381168114620002ca57600080fd5b60208501519093506001600160401b0380821115620002e857600080fd5b818601915086601f830112620002fd57600080fd5b8151818111156200031257620003126200042e565b604051601f8201601f19908116603f011681019083821181831017156200033d576200033d6200042e565b816040528281528960208487010111156200035757600080fd5b6200036a836020830160208801620003be565b8096505050505050620003806040850162000285565b90509250925092565b6020815260008251806020840152620003aa816040850160208701620003be565b601f01601f19169190910160400192915050565b60005b83811015620003db578181015183820152602001620003c1565b83811115620003eb576000848401525b50505050565b600181811c908216806200040657607f821691505b602082108114156200042857634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052604160045260246000fd5b61136180620004546000396000f3fe608060405234801561001057600080fd5b50600436106100ea5760003560e01c80638da5cb5b1161008c578063d9ded5e011610066578063d9ded5e0146101ff578063ee1e260a14610212578063f2fde38b1461021a578063fa3ea6a31461022d57600080fd5b80638da5cb5b146101b2578063bb0109d3146101d7578063cebc9a82146101ea57600080fd5b80631cfdde7d116100c85780631cfdde7d1461017157806322932b591461018457806331a6ca6b1461019757806379ba5097146101aa57600080fd5b8063017d4892146100ef57806301ffc9a714610104578063181f5a771461013d575b600080fd5b6101026100fd366004610f76565b61023e565b005b610128610112366004610ec6565b6001600160e01b03191663b6f6b1c560e01b1490565b60405190151581526020015b60405180910390f35b60408051808201909152600f81526e45786368616e67657220302e302e3160881b60208201525b6040516101349190611198565b61010261017f366004610e1e565b610470565b610102610192366004610fda565b6104cd565b6101646101a5366004610ef0565b61053a565b610102610584565b6000546001600160a01b03165b6040516001600160a01b039091168152602001610134565b6101026101e5366004610e3b565b61062e565b60045460405160ff9091168152602001610134565b61010261020d3660046110b6565b61067b565b6101646106c5565b610102610228366004610e1e565b61076a565b6002546001600160a01b03166101bf565b6000818060200190518101906102549190611023565b825160208085019190912060045460008281526005909352604090922054929350916102839160ff1690611204565b4311156102d45760045460008281526005602052604090205443916102ad9160ff90911690611204565b60405163647d550d60e01b8152600481019290925260248201526044015b60405180910390fd5b60025460405163473b057f60e11b81526000916001600160a01b031690638e760afe90610305908890600401611198565b600060405180830381600087803b15801561031f57600080fd5b505af1158015610333573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f1916820160405261035b9190810190610f2d565b9050600080600080848060200190518101906103779190610e54565b9350935093509350866000015184146103b0578651604051630841559760e41b81526102cb918691600401918252602082015260400190565b7f67775583a9c58f84357f29dae06d2dacf050cbc4ed5cbeab81d0790d9953c1728760200151886040015189606001518a608001518b60a001518c60c00151878a8a60405161045d999897969594939291909889526020890197909752604088019590955260608701939093526001600160a01b0391821660808701521660a085015260170b60c084015263ffffffff1660e083015267ffffffffffffffff166101008201526101200190565b60405180910390a1505050505050505050565b61047861077e565b600280546001600160a01b0319166001600160a01b0383169081179091556040519081527f8d9f13aae8f2e086b6c478fcb20b85e3f5aab0fcbf26e5d13949ffb23017e539906020015b60405180910390a150565b6104d561077e565b806040516020016104e69190611105565b6040516020818303038152906040526003908051906020019061050a929190610cf7565b507fc849f01f1579431074588d2d77603b0c2754e6776e333aaa728de3a13ad967a7816040516104c29190611198565b80516020808301918220600081815260058352604081205485516060959394919361056b9291880190910190611023565b905061057b8160000151836107d3565b95945050505050565b6001546001600160a01b031633146105d75760405162461bcd60e51b815260206004820152601660248201527526bab9ba10313290383937b837b9b2b21037bbb732b960511b60448201526064016102cb565b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b61063960014361123b565b60008281526005602090815260409182902092909255518281527fa78c3d4234d02c96f4c5223b1956b29eaf57fa2ab28c249e82ee8fecb16fa66b91016104c2565b61068361077e565b6004805460ff191660ff83169081179091556040519081527f89bc7ef46e5099f5518e68e1171b32bc77b5879fa9f1cc154499aeff37182e9b906020016104c2565b6060600380546106d490611299565b80601f016020809104026020016040519081016040528092919081815260200182805461070090611299565b801561074d5780601f106107225761010080835404028352916020019161074d565b820191906000526020600020905b81548152906001019060200180831161073057829003601f168201915b50505050508060200190518101906107659190610f2d565b905090565b61077261077e565b61077b816108ae565b50565b6000546001600160a01b031633146107d15760405162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b60448201526064016102cb565b565b6060600380546107e290611299565b80601f016020809104026020016040519081016040528092919081815260200182805461080e90611299565b801561085b5780601f106108305761010080835404028352916020019161085b565b820191906000526020600020905b81548152906001019060200180831161083e57829003601f168201915b50505050508060200190518101906108739190610f2d565b61087c84610958565b61088584610975565b60405160200161089793929190611121565b604051602081830303815290604052905092915050565b6001600160a01b0381163314156109075760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016102cb565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b606061096f8261096784610a12565b600101610a7c565b92915050565b6060600061098283610c1f565b600101905060008167ffffffffffffffff8111156109a2576109a2611300565b6040519080825280601f01601f1916602001820160405280156109cc576020820181803683370190505b5090508181016020015b600019016f181899199a1a9b1b9c1cb0b131b232b360811b600a86061a8153600a8504945084610a0557610a0a565b6109d6565b509392505050565b600080608083901c15610a2a5760809290921c916010015b604083901c15610a3f5760409290921c916008015b602083901c15610a545760209290921c916004015b601083901c15610a695760109290921c916002015b600883901c1561096f5760010192915050565b60606000610a8b83600261121c565b610a96906002611204565b67ffffffffffffffff811115610aae57610aae611300565b6040519080825280601f01601f191660200182016040528015610ad8576020820181803683370190505b509050600360fc1b81600081518110610af357610af36112ea565b60200101906001600160f81b031916908160001a905350600f60fb1b81600181518110610b2257610b226112ea565b60200101906001600160f81b031916908160001a9053506000610b4684600261121c565b610b51906001611204565b90505b6001811115610bc9576f181899199a1a9b1b9c1cb0b131b232b360811b85600f1660108110610b8557610b856112ea565b1a60f81b828281518110610b9b57610b9b6112ea565b60200101906001600160f81b031916908160001a90535060049490941c93610bc281611282565b9050610b54565b508315610c185760405162461bcd60e51b815260206004820181905260248201527f537472696e67733a20686578206c656e67746820696e73756666696369656e7460448201526064016102cb565b9392505050565b60008072184f03e93ff9f4daa797ed6e38ed64bf6a1f0160401b8310610c5e5772184f03e93ff9f4daa797ed6e38ed64bf6a1f0160401b830492506040015b6d04ee2d6d415b85acef81000000008310610c8a576d04ee2d6d415b85acef8100000000830492506020015b662386f26fc100008310610ca857662386f26fc10000830492506010015b6305f5e1008310610cc0576305f5e100830492506008015b6127108310610cd457612710830492506004015b60648310610ce6576064830492506002015b600a831061096f5760010192915050565b828054610d0390611299565b90600052602060002090601f016020900481019282610d255760008555610d6b565b82601f10610d3e57805160ff1916838001178555610d6b565b82800160010185558215610d6b579182015b82811115610d6b578251825591602001919060010190610d50565b50610d77929150610d7b565b5090565b5b80821115610d775760008155600101610d7c565b6000610da3610d9e846111dc565b6111ab565b9050828152838383011115610db757600080fd5b828260208301376000602084830101529392505050565b6000610ddc610d9e846111dc565b9050828152838383011115610df057600080fd5b610c18836020830184611252565b600082601f830112610e0f57600080fd5b610c1883833560208501610d90565b600060208284031215610e3057600080fd5b8135610c1881611316565b600060208284031215610e4d57600080fd5b5035919050565b60008060008060808587031215610e6a57600080fd5b84519350602085015163ffffffff81168114610e8557600080fd5b604086015190935067ffffffffffffffff81168114610ea357600080fd5b8092505060608501518060170b8114610ebb57600080fd5b939692955090935050565b600060208284031215610ed857600080fd5b81356001600160e01b031981168114610c1857600080fd5b600060208284031215610f0257600080fd5b813567ffffffffffffffff811115610f1957600080fd5b610f2584828501610dfe565b949350505050565b600060208284031215610f3f57600080fd5b815167ffffffffffffffff811115610f5657600080fd5b8201601f81018413610f6757600080fd5b610f2584825160208401610dce565b60008060408385031215610f8957600080fd5b823567ffffffffffffffff80821115610fa157600080fd5b610fad86838701610dfe565b93506020850135915080821115610fc357600080fd5b50610fd085828601610dfe565b9150509250929050565b600060208284031215610fec57600080fd5b813567ffffffffffffffff81111561100357600080fd5b8201601f8101841361101457600080fd5b610f2584823560208401610d90565b600060e0828403121561103557600080fd5b60405160e0810181811067ffffffffffffffff8211171561105857611058611300565b8060405250825181526020830151602082015260408301516040820152606083015160608201526080830151608082015260a083015161109781611316565b60a082015260c08301516110aa81611316565b60c08201529392505050565b6000602082840312156110c857600080fd5b813560ff81168114610c1857600080fd5b600081518084526110f1816020860160208601611252565b601f01601f19169290920160200192915050565b60008251611117818460208701611252565b9190910192915050565b60008451611133818460208901611252565b6a3f6665656449444865783d60a81b908301908152845161115b81600b840160208901611252565b6e264c32426c6f636b6e756d6265723d60881b600b9290910191820152835161118b81601a840160208801611252565b01601a0195945050505050565b602081526000610c1860208301846110d9565b604051601f8201601f1916810167ffffffffffffffff811182821017156111d4576111d4611300565b604052919050565b600067ffffffffffffffff8211156111f6576111f6611300565b50601f01601f191660200190565b60008219821115611217576112176112d4565b500190565b6000816000190483118215151615611236576112366112d4565b500290565b60008282101561124d5761124d6112d4565b500390565b60005b8381101561126d578181015183820152602001611255565b8381111561127c576000848401525b50505050565b600081611291576112916112d4565b506000190190565b600181811c908216806112ad57607f821691505b602082108114156112ce57634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052604160045260246000fd5b6001600160a01b038116811461077b57600080fdfea26469706673582212204aa1094080e003cb8f29cb36ceefd507d1eecdff27b4f75c2a6a2178cd78305d64736f6c63430008060033",
}

// ExchangerABI is the input ABI used to generate the binding from.
// Deprecated: Use ExchangerMetaData.ABI instead.
var ExchangerABI = ExchangerMetaData.ABI

// ExchangerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ExchangerMetaData.Bin instead.
var ExchangerBin = ExchangerMetaData.Bin

// DeployExchanger deploys a new Ethereum contract, binding an instance of Exchanger to it.
func DeployExchanger(auth *bind.TransactOpts, backend bind.ContractBackend, verifierProxyAddress common.Address, lookupURL string, maxDelay uint8) (common.Address, *types.Transaction, *Exchanger, error) {
	parsed, err := ExchangerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ExchangerBin), backend, verifierProxyAddress, lookupURL, maxDelay)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Exchanger{ExchangerCaller: ExchangerCaller{contract: contract}, ExchangerTransactor: ExchangerTransactor{contract: contract}, ExchangerFilterer: ExchangerFilterer{contract: contract}}, nil
}

// Exchanger is an auto generated Go binding around an Ethereum contract.
type Exchanger struct {
	ExchangerCaller     // Read-only binding to the contract
	ExchangerTransactor // Write-only binding to the contract
	ExchangerFilterer   // Log filterer for contract events
}

// ExchangerCaller is an auto generated read-only Go binding around an Ethereum contract.
type ExchangerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExchangerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ExchangerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExchangerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ExchangerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExchangerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ExchangerSession struct {
	Contract     *Exchanger        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ExchangerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ExchangerCallerSession struct {
	Contract *ExchangerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// ExchangerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ExchangerTransactorSession struct {
	Contract     *ExchangerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// ExchangerRaw is an auto generated low-level Go binding around an Ethereum contract.
type ExchangerRaw struct {
	Contract *Exchanger // Generic contract binding to access the raw methods on
}

// ExchangerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ExchangerCallerRaw struct {
	Contract *ExchangerCaller // Generic read-only contract binding to access the raw methods on
}

// ExchangerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ExchangerTransactorRaw struct {
	Contract *ExchangerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewExchanger creates a new instance of Exchanger, bound to a specific deployed contract.
func NewExchanger(address common.Address, backend bind.ContractBackend) (*Exchanger, error) {
	contract, err := bindExchanger(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Exchanger{ExchangerCaller: ExchangerCaller{contract: contract}, ExchangerTransactor: ExchangerTransactor{contract: contract}, ExchangerFilterer: ExchangerFilterer{contract: contract}}, nil
}

// NewExchangerCaller creates a new read-only instance of Exchanger, bound to a specific deployed contract.
func NewExchangerCaller(address common.Address, caller bind.ContractCaller) (*ExchangerCaller, error) {
	contract, err := bindExchanger(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ExchangerCaller{contract: contract}, nil
}

// NewExchangerTransactor creates a new write-only instance of Exchanger, bound to a specific deployed contract.
func NewExchangerTransactor(address common.Address, transactor bind.ContractTransactor) (*ExchangerTransactor, error) {
	contract, err := bindExchanger(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ExchangerTransactor{contract: contract}, nil
}

// NewExchangerFilterer creates a new log filterer instance of Exchanger, bound to a specific deployed contract.
func NewExchangerFilterer(address common.Address, filterer bind.ContractFilterer) (*ExchangerFilterer, error) {
	contract, err := bindExchanger(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ExchangerFilterer{contract: contract}, nil
}

// bindExchanger binds a generic wrapper to an already deployed contract.
func bindExchanger(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ExchangerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Exchanger *ExchangerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Exchanger.Contract.ExchangerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Exchanger *ExchangerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Exchanger.Contract.ExchangerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Exchanger *ExchangerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Exchanger.Contract.ExchangerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Exchanger *ExchangerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Exchanger.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Exchanger *ExchangerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Exchanger.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Exchanger *ExchangerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Exchanger.Contract.contract.Transact(opts, method, params...)
}

// GetDelay is a free data retrieval call binding the contract method 0xcebc9a82.
//
// Solidity: function getDelay() view returns(uint8 maxDelay)
func (_Exchanger *ExchangerCaller) GetDelay(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Exchanger.contract.Call(opts, &out, "getDelay")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetDelay is a free data retrieval call binding the contract method 0xcebc9a82.
//
// Solidity: function getDelay() view returns(uint8 maxDelay)
func (_Exchanger *ExchangerSession) GetDelay() (uint8, error) {
	return _Exchanger.Contract.GetDelay(&_Exchanger.CallOpts)
}

// GetDelay is a free data retrieval call binding the contract method 0xcebc9a82.
//
// Solidity: function getDelay() view returns(uint8 maxDelay)
func (_Exchanger *ExchangerCallerSession) GetDelay() (uint8, error) {
	return _Exchanger.Contract.GetDelay(&_Exchanger.CallOpts)
}

// GetLookupURL is a free data retrieval call binding the contract method 0xee1e260a.
//
// Solidity: function getLookupURL() view returns(string url)
func (_Exchanger *ExchangerCaller) GetLookupURL(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Exchanger.contract.Call(opts, &out, "getLookupURL")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetLookupURL is a free data retrieval call binding the contract method 0xee1e260a.
//
// Solidity: function getLookupURL() view returns(string url)
func (_Exchanger *ExchangerSession) GetLookupURL() (string, error) {
	return _Exchanger.Contract.GetLookupURL(&_Exchanger.CallOpts)
}

// GetLookupURL is a free data retrieval call binding the contract method 0xee1e260a.
//
// Solidity: function getLookupURL() view returns(string url)
func (_Exchanger *ExchangerCallerSession) GetLookupURL() (string, error) {
	return _Exchanger.Contract.GetLookupURL(&_Exchanger.CallOpts)
}

// GetVerifierProxyAddress is a free data retrieval call binding the contract method 0xfa3ea6a3.
//
// Solidity: function getVerifierProxyAddress() view returns(address verifierProxyAddress)
func (_Exchanger *ExchangerCaller) GetVerifierProxyAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Exchanger.contract.Call(opts, &out, "getVerifierProxyAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetVerifierProxyAddress is a free data retrieval call binding the contract method 0xfa3ea6a3.
//
// Solidity: function getVerifierProxyAddress() view returns(address verifierProxyAddress)
func (_Exchanger *ExchangerSession) GetVerifierProxyAddress() (common.Address, error) {
	return _Exchanger.Contract.GetVerifierProxyAddress(&_Exchanger.CallOpts)
}

// GetVerifierProxyAddress is a free data retrieval call binding the contract method 0xfa3ea6a3.
//
// Solidity: function getVerifierProxyAddress() view returns(address verifierProxyAddress)
func (_Exchanger *ExchangerCallerSession) GetVerifierProxyAddress() (common.Address, error) {
	return _Exchanger.Contract.GetVerifierProxyAddress(&_Exchanger.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Exchanger *ExchangerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Exchanger.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Exchanger *ExchangerSession) Owner() (common.Address, error) {
	return _Exchanger.Contract.Owner(&_Exchanger.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Exchanger *ExchangerCallerSession) Owner() (common.Address, error) {
	return _Exchanger.Contract.Owner(&_Exchanger.CallOpts)
}

// ResolveTrade is a free data retrieval call binding the contract method 0x31a6ca6b.
//
// Solidity: function resolveTrade(bytes encodedCommitment) view returns(string)
func (_Exchanger *ExchangerCaller) ResolveTrade(opts *bind.CallOpts, encodedCommitment []byte) (string, error) {
	var out []interface{}
	err := _Exchanger.contract.Call(opts, &out, "resolveTrade", encodedCommitment)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// ResolveTrade is a free data retrieval call binding the contract method 0x31a6ca6b.
//
// Solidity: function resolveTrade(bytes encodedCommitment) view returns(string)
func (_Exchanger *ExchangerSession) ResolveTrade(encodedCommitment []byte) (string, error) {
	return _Exchanger.Contract.ResolveTrade(&_Exchanger.CallOpts, encodedCommitment)
}

// ResolveTrade is a free data retrieval call binding the contract method 0x31a6ca6b.
//
// Solidity: function resolveTrade(bytes encodedCommitment) view returns(string)
func (_Exchanger *ExchangerCallerSession) ResolveTrade(encodedCommitment []byte) (string, error) {
	return _Exchanger.Contract.ResolveTrade(&_Exchanger.CallOpts, encodedCommitment)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) pure returns(bool)
func (_Exchanger *ExchangerCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _Exchanger.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) pure returns(bool)
func (_Exchanger *ExchangerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Exchanger.Contract.SupportsInterface(&_Exchanger.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) pure returns(bool)
func (_Exchanger *ExchangerCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Exchanger.Contract.SupportsInterface(&_Exchanger.CallOpts, interfaceId)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_Exchanger *ExchangerCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Exchanger.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_Exchanger *ExchangerSession) TypeAndVersion() (string, error) {
	return _Exchanger.Contract.TypeAndVersion(&_Exchanger.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_Exchanger *ExchangerCallerSession) TypeAndVersion() (string, error) {
	return _Exchanger.Contract.TypeAndVersion(&_Exchanger.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Exchanger *ExchangerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Exchanger.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Exchanger *ExchangerSession) AcceptOwnership() (*types.Transaction, error) {
	return _Exchanger.Contract.AcceptOwnership(&_Exchanger.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Exchanger *ExchangerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Exchanger.Contract.AcceptOwnership(&_Exchanger.TransactOpts)
}

// CommitTrade is a paid mutator transaction binding the contract method 0xbb0109d3.
//
// Solidity: function commitTrade(bytes32 commitment) returns()
func (_Exchanger *ExchangerTransactor) CommitTrade(opts *bind.TransactOpts, commitment [32]byte) (*types.Transaction, error) {
	return _Exchanger.contract.Transact(opts, "commitTrade", commitment)
}

// CommitTrade is a paid mutator transaction binding the contract method 0xbb0109d3.
//
// Solidity: function commitTrade(bytes32 commitment) returns()
func (_Exchanger *ExchangerSession) CommitTrade(commitment [32]byte) (*types.Transaction, error) {
	return _Exchanger.Contract.CommitTrade(&_Exchanger.TransactOpts, commitment)
}

// CommitTrade is a paid mutator transaction binding the contract method 0xbb0109d3.
//
// Solidity: function commitTrade(bytes32 commitment) returns()
func (_Exchanger *ExchangerTransactorSession) CommitTrade(commitment [32]byte) (*types.Transaction, error) {
	return _Exchanger.Contract.CommitTrade(&_Exchanger.TransactOpts, commitment)
}

// ResolveTradeWithReport is a paid mutator transaction binding the contract method 0x017d4892.
//
// Solidity: function resolveTradeWithReport(bytes chainlinkBlob, bytes encodedCommitment) returns()
func (_Exchanger *ExchangerTransactor) ResolveTradeWithReport(opts *bind.TransactOpts, chainlinkBlob []byte, encodedCommitment []byte) (*types.Transaction, error) {
	return _Exchanger.contract.Transact(opts, "resolveTradeWithReport", chainlinkBlob, encodedCommitment)
}

// ResolveTradeWithReport is a paid mutator transaction binding the contract method 0x017d4892.
//
// Solidity: function resolveTradeWithReport(bytes chainlinkBlob, bytes encodedCommitment) returns()
func (_Exchanger *ExchangerSession) ResolveTradeWithReport(chainlinkBlob []byte, encodedCommitment []byte) (*types.Transaction, error) {
	return _Exchanger.Contract.ResolveTradeWithReport(&_Exchanger.TransactOpts, chainlinkBlob, encodedCommitment)
}

// ResolveTradeWithReport is a paid mutator transaction binding the contract method 0x017d4892.
//
// Solidity: function resolveTradeWithReport(bytes chainlinkBlob, bytes encodedCommitment) returns()
func (_Exchanger *ExchangerTransactorSession) ResolveTradeWithReport(chainlinkBlob []byte, encodedCommitment []byte) (*types.Transaction, error) {
	return _Exchanger.Contract.ResolveTradeWithReport(&_Exchanger.TransactOpts, chainlinkBlob, encodedCommitment)
}

// SetDelay is a paid mutator transaction binding the contract method 0xd9ded5e0.
//
// Solidity: function setDelay(uint8 maxDelay) returns()
func (_Exchanger *ExchangerTransactor) SetDelay(opts *bind.TransactOpts, maxDelay uint8) (*types.Transaction, error) {
	return _Exchanger.contract.Transact(opts, "setDelay", maxDelay)
}

// SetDelay is a paid mutator transaction binding the contract method 0xd9ded5e0.
//
// Solidity: function setDelay(uint8 maxDelay) returns()
func (_Exchanger *ExchangerSession) SetDelay(maxDelay uint8) (*types.Transaction, error) {
	return _Exchanger.Contract.SetDelay(&_Exchanger.TransactOpts, maxDelay)
}

// SetDelay is a paid mutator transaction binding the contract method 0xd9ded5e0.
//
// Solidity: function setDelay(uint8 maxDelay) returns()
func (_Exchanger *ExchangerTransactorSession) SetDelay(maxDelay uint8) (*types.Transaction, error) {
	return _Exchanger.Contract.SetDelay(&_Exchanger.TransactOpts, maxDelay)
}

// SetLookupURL is a paid mutator transaction binding the contract method 0x22932b59.
//
// Solidity: function setLookupURL(string url) returns()
func (_Exchanger *ExchangerTransactor) SetLookupURL(opts *bind.TransactOpts, url string) (*types.Transaction, error) {
	return _Exchanger.contract.Transact(opts, "setLookupURL", url)
}

// SetLookupURL is a paid mutator transaction binding the contract method 0x22932b59.
//
// Solidity: function setLookupURL(string url) returns()
func (_Exchanger *ExchangerSession) SetLookupURL(url string) (*types.Transaction, error) {
	return _Exchanger.Contract.SetLookupURL(&_Exchanger.TransactOpts, url)
}

// SetLookupURL is a paid mutator transaction binding the contract method 0x22932b59.
//
// Solidity: function setLookupURL(string url) returns()
func (_Exchanger *ExchangerTransactorSession) SetLookupURL(url string) (*types.Transaction, error) {
	return _Exchanger.Contract.SetLookupURL(&_Exchanger.TransactOpts, url)
}

// SetVerifierProxyAddress is a paid mutator transaction binding the contract method 0x1cfdde7d.
//
// Solidity: function setVerifierProxyAddress(address verifierProxyAddress) returns()
func (_Exchanger *ExchangerTransactor) SetVerifierProxyAddress(opts *bind.TransactOpts, verifierProxyAddress common.Address) (*types.Transaction, error) {
	return _Exchanger.contract.Transact(opts, "setVerifierProxyAddress", verifierProxyAddress)
}

// SetVerifierProxyAddress is a paid mutator transaction binding the contract method 0x1cfdde7d.
//
// Solidity: function setVerifierProxyAddress(address verifierProxyAddress) returns()
func (_Exchanger *ExchangerSession) SetVerifierProxyAddress(verifierProxyAddress common.Address) (*types.Transaction, error) {
	return _Exchanger.Contract.SetVerifierProxyAddress(&_Exchanger.TransactOpts, verifierProxyAddress)
}

// SetVerifierProxyAddress is a paid mutator transaction binding the contract method 0x1cfdde7d.
//
// Solidity: function setVerifierProxyAddress(address verifierProxyAddress) returns()
func (_Exchanger *ExchangerTransactorSession) SetVerifierProxyAddress(verifierProxyAddress common.Address) (*types.Transaction, error) {
	return _Exchanger.Contract.SetVerifierProxyAddress(&_Exchanger.TransactOpts, verifierProxyAddress)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_Exchanger *ExchangerTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _Exchanger.contract.Transact(opts, "transferOwnership", to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_Exchanger *ExchangerSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _Exchanger.Contract.TransferOwnership(&_Exchanger.TransactOpts, to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_Exchanger *ExchangerTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _Exchanger.Contract.TransferOwnership(&_Exchanger.TransactOpts, to)
}

// ExchangerOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the Exchanger contract.
type ExchangerOwnershipTransferRequestedIterator struct {
	Event *ExchangerOwnershipTransferRequested // Event containing the contract specifics and raw log

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
func (it *ExchangerOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangerOwnershipTransferRequested)
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
		it.Event = new(ExchangerOwnershipTransferRequested)
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
func (it *ExchangerOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangerOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangerOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the Exchanger contract.
type ExchangerOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_Exchanger *ExchangerFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ExchangerOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Exchanger.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ExchangerOwnershipTransferRequestedIterator{contract: _Exchanger.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_Exchanger *ExchangerFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *ExchangerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Exchanger.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangerOwnershipTransferRequested)
				if err := _Exchanger.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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
func (_Exchanger *ExchangerFilterer) ParseOwnershipTransferRequested(log types.Log) (*ExchangerOwnershipTransferRequested, error) {
	event := new(ExchangerOwnershipTransferRequested)
	if err := _Exchanger.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ExchangerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Exchanger contract.
type ExchangerOwnershipTransferredIterator struct {
	Event *ExchangerOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ExchangerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangerOwnershipTransferred)
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
		it.Event = new(ExchangerOwnershipTransferred)
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
func (it *ExchangerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangerOwnershipTransferred represents a OwnershipTransferred event raised by the Exchanger contract.
type ExchangerOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_Exchanger *ExchangerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ExchangerOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Exchanger.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ExchangerOwnershipTransferredIterator{contract: _Exchanger.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_Exchanger *ExchangerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ExchangerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Exchanger.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangerOwnershipTransferred)
				if err := _Exchanger.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Exchanger *ExchangerFilterer) ParseOwnershipTransferred(log types.Log) (*ExchangerOwnershipTransferred, error) {
	event := new(ExchangerOwnershipTransferred)
	if err := _Exchanger.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ExchangerSetDelayIterator is returned from FilterSetDelay and is used to iterate over the raw logs and unpacked data for SetDelay events raised by the Exchanger contract.
type ExchangerSetDelayIterator struct {
	Event *ExchangerSetDelay // Event containing the contract specifics and raw log

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
func (it *ExchangerSetDelayIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangerSetDelay)
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
		it.Event = new(ExchangerSetDelay)
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
func (it *ExchangerSetDelayIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangerSetDelayIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangerSetDelay represents a SetDelay event raised by the Exchanger contract.
type ExchangerSetDelay struct {
	MaxDelay uint8
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSetDelay is a free log retrieval operation binding the contract event 0x89bc7ef46e5099f5518e68e1171b32bc77b5879fa9f1cc154499aeff37182e9b.
//
// Solidity: event SetDelay(uint8 maxDelay)
func (_Exchanger *ExchangerFilterer) FilterSetDelay(opts *bind.FilterOpts) (*ExchangerSetDelayIterator, error) {

	logs, sub, err := _Exchanger.contract.FilterLogs(opts, "SetDelay")
	if err != nil {
		return nil, err
	}
	return &ExchangerSetDelayIterator{contract: _Exchanger.contract, event: "SetDelay", logs: logs, sub: sub}, nil
}

// WatchSetDelay is a free log subscription operation binding the contract event 0x89bc7ef46e5099f5518e68e1171b32bc77b5879fa9f1cc154499aeff37182e9b.
//
// Solidity: event SetDelay(uint8 maxDelay)
func (_Exchanger *ExchangerFilterer) WatchSetDelay(opts *bind.WatchOpts, sink chan<- *ExchangerSetDelay) (event.Subscription, error) {

	logs, sub, err := _Exchanger.contract.WatchLogs(opts, "SetDelay")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangerSetDelay)
				if err := _Exchanger.contract.UnpackLog(event, "SetDelay", log); err != nil {
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

// ParseSetDelay is a log parse operation binding the contract event 0x89bc7ef46e5099f5518e68e1171b32bc77b5879fa9f1cc154499aeff37182e9b.
//
// Solidity: event SetDelay(uint8 maxDelay)
func (_Exchanger *ExchangerFilterer) ParseSetDelay(log types.Log) (*ExchangerSetDelay, error) {
	event := new(ExchangerSetDelay)
	if err := _Exchanger.contract.UnpackLog(event, "SetDelay", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ExchangerSetLookupURLIterator is returned from FilterSetLookupURL and is used to iterate over the raw logs and unpacked data for SetLookupURL events raised by the Exchanger contract.
type ExchangerSetLookupURLIterator struct {
	Event *ExchangerSetLookupURL // Event containing the contract specifics and raw log

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
func (it *ExchangerSetLookupURLIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangerSetLookupURL)
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
		it.Event = new(ExchangerSetLookupURL)
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
func (it *ExchangerSetLookupURLIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangerSetLookupURLIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangerSetLookupURL represents a SetLookupURL event raised by the Exchanger contract.
type ExchangerSetLookupURL struct {
	Url string
	Raw types.Log // Blockchain specific contextual infos
}

// FilterSetLookupURL is a free log retrieval operation binding the contract event 0xc849f01f1579431074588d2d77603b0c2754e6776e333aaa728de3a13ad967a7.
//
// Solidity: event SetLookupURL(string url)
func (_Exchanger *ExchangerFilterer) FilterSetLookupURL(opts *bind.FilterOpts) (*ExchangerSetLookupURLIterator, error) {

	logs, sub, err := _Exchanger.contract.FilterLogs(opts, "SetLookupURL")
	if err != nil {
		return nil, err
	}
	return &ExchangerSetLookupURLIterator{contract: _Exchanger.contract, event: "SetLookupURL", logs: logs, sub: sub}, nil
}

// WatchSetLookupURL is a free log subscription operation binding the contract event 0xc849f01f1579431074588d2d77603b0c2754e6776e333aaa728de3a13ad967a7.
//
// Solidity: event SetLookupURL(string url)
func (_Exchanger *ExchangerFilterer) WatchSetLookupURL(opts *bind.WatchOpts, sink chan<- *ExchangerSetLookupURL) (event.Subscription, error) {

	logs, sub, err := _Exchanger.contract.WatchLogs(opts, "SetLookupURL")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangerSetLookupURL)
				if err := _Exchanger.contract.UnpackLog(event, "SetLookupURL", log); err != nil {
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

// ParseSetLookupURL is a log parse operation binding the contract event 0xc849f01f1579431074588d2d77603b0c2754e6776e333aaa728de3a13ad967a7.
//
// Solidity: event SetLookupURL(string url)
func (_Exchanger *ExchangerFilterer) ParseSetLookupURL(log types.Log) (*ExchangerSetLookupURL, error) {
	event := new(ExchangerSetLookupURL)
	if err := _Exchanger.contract.UnpackLog(event, "SetLookupURL", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ExchangerSetVerifierProxyIterator is returned from FilterSetVerifierProxy and is used to iterate over the raw logs and unpacked data for SetVerifierProxy events raised by the Exchanger contract.
type ExchangerSetVerifierProxyIterator struct {
	Event *ExchangerSetVerifierProxy // Event containing the contract specifics and raw log

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
func (it *ExchangerSetVerifierProxyIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangerSetVerifierProxy)
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
		it.Event = new(ExchangerSetVerifierProxy)
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
func (it *ExchangerSetVerifierProxyIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangerSetVerifierProxyIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangerSetVerifierProxy represents a SetVerifierProxy event raised by the Exchanger contract.
type ExchangerSetVerifierProxy struct {
	VerifierProxyAddress common.Address
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterSetVerifierProxy is a free log retrieval operation binding the contract event 0x8d9f13aae8f2e086b6c478fcb20b85e3f5aab0fcbf26e5d13949ffb23017e539.
//
// Solidity: event SetVerifierProxy(address verifierProxyAddress)
func (_Exchanger *ExchangerFilterer) FilterSetVerifierProxy(opts *bind.FilterOpts) (*ExchangerSetVerifierProxyIterator, error) {

	logs, sub, err := _Exchanger.contract.FilterLogs(opts, "SetVerifierProxy")
	if err != nil {
		return nil, err
	}
	return &ExchangerSetVerifierProxyIterator{contract: _Exchanger.contract, event: "SetVerifierProxy", logs: logs, sub: sub}, nil
}

// WatchSetVerifierProxy is a free log subscription operation binding the contract event 0x8d9f13aae8f2e086b6c478fcb20b85e3f5aab0fcbf26e5d13949ffb23017e539.
//
// Solidity: event SetVerifierProxy(address verifierProxyAddress)
func (_Exchanger *ExchangerFilterer) WatchSetVerifierProxy(opts *bind.WatchOpts, sink chan<- *ExchangerSetVerifierProxy) (event.Subscription, error) {

	logs, sub, err := _Exchanger.contract.WatchLogs(opts, "SetVerifierProxy")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangerSetVerifierProxy)
				if err := _Exchanger.contract.UnpackLog(event, "SetVerifierProxy", log); err != nil {
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

// ParseSetVerifierProxy is a log parse operation binding the contract event 0x8d9f13aae8f2e086b6c478fcb20b85e3f5aab0fcbf26e5d13949ffb23017e539.
//
// Solidity: event SetVerifierProxy(address verifierProxyAddress)
func (_Exchanger *ExchangerFilterer) ParseSetVerifierProxy(log types.Log) (*ExchangerSetVerifierProxy, error) {
	event := new(ExchangerSetVerifierProxy)
	if err := _Exchanger.contract.UnpackLog(event, "SetVerifierProxy", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ExchangerTradeCommittedIterator is returned from FilterTradeCommitted and is used to iterate over the raw logs and unpacked data for TradeCommitted events raised by the Exchanger contract.
type ExchangerTradeCommittedIterator struct {
	Event *ExchangerTradeCommitted // Event containing the contract specifics and raw log

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
func (it *ExchangerTradeCommittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangerTradeCommitted)
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
		it.Event = new(ExchangerTradeCommitted)
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
func (it *ExchangerTradeCommittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangerTradeCommittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangerTradeCommitted represents a TradeCommitted event raised by the Exchanger contract.
type ExchangerTradeCommitted struct {
	Commitment [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterTradeCommitted is a free log retrieval operation binding the contract event 0xa78c3d4234d02c96f4c5223b1956b29eaf57fa2ab28c249e82ee8fecb16fa66b.
//
// Solidity: event TradeCommitted(bytes32 commitment)
func (_Exchanger *ExchangerFilterer) FilterTradeCommitted(opts *bind.FilterOpts) (*ExchangerTradeCommittedIterator, error) {

	logs, sub, err := _Exchanger.contract.FilterLogs(opts, "TradeCommitted")
	if err != nil {
		return nil, err
	}
	return &ExchangerTradeCommittedIterator{contract: _Exchanger.contract, event: "TradeCommitted", logs: logs, sub: sub}, nil
}

// WatchTradeCommitted is a free log subscription operation binding the contract event 0xa78c3d4234d02c96f4c5223b1956b29eaf57fa2ab28c249e82ee8fecb16fa66b.
//
// Solidity: event TradeCommitted(bytes32 commitment)
func (_Exchanger *ExchangerFilterer) WatchTradeCommitted(opts *bind.WatchOpts, sink chan<- *ExchangerTradeCommitted) (event.Subscription, error) {

	logs, sub, err := _Exchanger.contract.WatchLogs(opts, "TradeCommitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangerTradeCommitted)
				if err := _Exchanger.contract.UnpackLog(event, "TradeCommitted", log); err != nil {
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

// ParseTradeCommitted is a log parse operation binding the contract event 0xa78c3d4234d02c96f4c5223b1956b29eaf57fa2ab28c249e82ee8fecb16fa66b.
//
// Solidity: event TradeCommitted(bytes32 commitment)
func (_Exchanger *ExchangerFilterer) ParseTradeCommitted(log types.Log) (*ExchangerTradeCommitted, error) {
	event := new(ExchangerTradeCommitted)
	if err := _Exchanger.contract.UnpackLog(event, "TradeCommitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ExchangerTradeExecutedIterator is returned from FilterTradeExecuted and is used to iterate over the raw logs and unpacked data for TradeExecuted events raised by the Exchanger contract.
type ExchangerTradeExecutedIterator struct {
	Event *ExchangerTradeExecuted // Event containing the contract specifics and raw log

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
func (it *ExchangerTradeExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangerTradeExecuted)
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
		it.Event = new(ExchangerTradeExecuted)
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
func (it *ExchangerTradeExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangerTradeExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangerTradeExecuted represents a TradeExecuted event raised by the Exchanger contract.
type ExchangerTradeExecuted struct {
	CurrencySrc             [32]byte
	CurrencyDst             [32]byte
	AmountSrc               *big.Int
	MinAmountDst            *big.Int
	Sender                  common.Address
	Receiver                common.Address
	Median                  *big.Int
	ObservationsTimestamp   uint32
	ObservationsBlocknumber uint64
	Raw                     types.Log // Blockchain specific contextual infos
}

// FilterTradeExecuted is a free log retrieval operation binding the contract event 0x67775583a9c58f84357f29dae06d2dacf050cbc4ed5cbeab81d0790d9953c172.
//
// Solidity: event TradeExecuted(bytes32 currencySrc, bytes32 currencyDst, uint256 amountSrc, uint256 minAmountDst, address sender, address receiver, int192 median, uint32 observationsTimestamp, uint64 observationsBlocknumber)
func (_Exchanger *ExchangerFilterer) FilterTradeExecuted(opts *bind.FilterOpts) (*ExchangerTradeExecutedIterator, error) {

	logs, sub, err := _Exchanger.contract.FilterLogs(opts, "TradeExecuted")
	if err != nil {
		return nil, err
	}
	return &ExchangerTradeExecutedIterator{contract: _Exchanger.contract, event: "TradeExecuted", logs: logs, sub: sub}, nil
}

// WatchTradeExecuted is a free log subscription operation binding the contract event 0x67775583a9c58f84357f29dae06d2dacf050cbc4ed5cbeab81d0790d9953c172.
//
// Solidity: event TradeExecuted(bytes32 currencySrc, bytes32 currencyDst, uint256 amountSrc, uint256 minAmountDst, address sender, address receiver, int192 median, uint32 observationsTimestamp, uint64 observationsBlocknumber)
func (_Exchanger *ExchangerFilterer) WatchTradeExecuted(opts *bind.WatchOpts, sink chan<- *ExchangerTradeExecuted) (event.Subscription, error) {

	logs, sub, err := _Exchanger.contract.WatchLogs(opts, "TradeExecuted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangerTradeExecuted)
				if err := _Exchanger.contract.UnpackLog(event, "TradeExecuted", log); err != nil {
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

// ParseTradeExecuted is a log parse operation binding the contract event 0x67775583a9c58f84357f29dae06d2dacf050cbc4ed5cbeab81d0790d9953c172.
//
// Solidity: event TradeExecuted(bytes32 currencySrc, bytes32 currencyDst, uint256 amountSrc, uint256 minAmountDst, address sender, address receiver, int192 median, uint32 observationsTimestamp, uint64 observationsBlocknumber)
func (_Exchanger *ExchangerFilterer) ParseTradeExecuted(log types.Log) (*ExchangerTradeExecuted, error) {
	event := new(ExchangerTradeExecuted)
	if err := _Exchanger.contract.UnpackLog(event, "TradeExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
