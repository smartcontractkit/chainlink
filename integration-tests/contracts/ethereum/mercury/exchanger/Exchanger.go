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
	_ = abi.ConvertType
)

// ExchangerMetaData contains all meta data concerning the Exchanger contract.
var ExchangerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIVerifierProxy\",\"name\":\"verifierProxyAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"lookupURL\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"maxDelay\",\"type\":\"uint8\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"reportBlockhash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"upperBoundBlockhash\",\"type\":\"bytes32\"}],\"name\":\"BlockhashMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"reportFeedID\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"commitmentFeedID\",\"type\":\"bytes32\"}],\"name\":\"FeedIDMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"string[]\",\"name\":\"urls\",\"type\":\"string[]\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"},{\"internalType\":\"bytes4\",\"name\":\"callbackFunction\",\"type\":\"bytes4\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"OffchainLookup\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"blocknumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tradeWindow\",\"type\":\"uint256\"}],\"name\":\"TradeExceedsWindow\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"maxDelay\",\"type\":\"uint8\"}],\"name\":\"SetDelay\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"url\",\"type\":\"string\"}],\"name\":\"SetLookupURL\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractIVerifierProxy\",\"name\":\"verifierProxyAddress\",\"type\":\"address\"}],\"name\":\"SetVerifierProxy\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"commitment\",\"type\":\"bytes32\"}],\"name\":\"TradeCommitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"currencySrc\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"currencyDst\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountSrc\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"minAmountDst\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"observationsTimestamp\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"blocknumberLowerBound\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"blocknumberUpperBound\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"upperBlockhash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"int192\",\"name\":\"median\",\"type\":\"int192\"},{\"indexed\":false,\"internalType\":\"int192\",\"name\":\"bid\",\"type\":\"int192\"},{\"indexed\":false,\"internalType\":\"int192\",\"name\":\"ask\",\"type\":\"int192\"}],\"name\":\"TradeExecuted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"commitment\",\"type\":\"bytes32\"}],\"name\":\"commitTrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDelay\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"maxDelay\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLookupURL\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"url\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getVerifierProxyAddress\",\"outputs\":[{\"internalType\":\"contractIVerifierProxy\",\"name\":\"verifierProxyAddress\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedCommitment\",\"type\":\"bytes\"}],\"name\":\"resolveTrade\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"chainlinkBlob\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"encodedCommitment\",\"type\":\"bytes\"}],\"name\":\"resolveTradeWithReport\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"maxDelay\",\"type\":\"uint8\"}],\"name\":\"setDelay\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"url\",\"type\":\"string\"}],\"name\":\"setLookupURL\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIVerifierProxy\",\"name\":\"verifierProxyAddress\",\"type\":\"address\"}],\"name\":\"setVerifierProxyAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b50604051620019513803806200195183398101604081905262000034916200022a565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be816200012c565b5050600280546001600160a01b0319166001600160a01b03861617905550604051620000ef90839060200162000317565b604051602081830303815290604052600390816200010e9190620003db565b506004805460ff191660ff9290921691909117905550620004a79050565b336001600160a01b03821603620001865760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b634e487b7160e01b600052604160045260246000fd5b60005b838110156200020a578181015183820152602001620001f0565b50506000910152565b805160ff811681146200022557600080fd5b919050565b6000806000606084860312156200024057600080fd5b83516001600160a01b03811681146200025857600080fd5b60208501519093506001600160401b03808211156200027657600080fd5b818601915086601f8301126200028b57600080fd5b815181811115620002a057620002a0620001d7565b604051601f8201601f19908116603f01168101908382118183101715620002cb57620002cb620001d7565b81604052828152896020848701011115620002e557600080fd5b620002f8836020830160208801620001ed565b80965050505050506200030e6040850162000213565b90509250925092565b602081526000825180602084015262000338816040850160208701620001ed565b601f01601f19169190910160400192915050565b600181811c908216806200036157607f821691505b6020821081036200038257634e487b7160e01b600052602260045260246000fd5b50919050565b601f821115620003d657600081815260208120601f850160051c81016020861015620003b15750805b601f850160051c820191505b81811015620003d257828155600101620003bd565b5050505b505050565b81516001600160401b03811115620003f757620003f7620001d7565b6200040f816200040884546200034c565b8462000388565b602080601f8311600181146200044757600084156200042e5750858301515b600019600386901b1c1916600185901b178555620003d2565b600085815260208120601f198616915b82811015620004785788860151825594840194600190910190840162000457565b5085821015620004975787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b61149a80620004b76000396000f3fe608060405234801561001057600080fd5b50600436106100ea5760003560e01c80638da5cb5b1161008c578063d9ded5e011610066578063d9ded5e0146101ff578063ee1e260a14610212578063f2fde38b1461021a578063fa3ea6a31461022d57600080fd5b80638da5cb5b146101b2578063bb0109d3146101d7578063cebc9a82146101ea57600080fd5b80631cfdde7d116100c85780631cfdde7d1461017157806322932b591461018457806331a6ca6b1461019757806379ba5097146101aa57600080fd5b8063017d4892146100ef57806301ffc9a714610104578063181f5a771461013d575b600080fd5b6101026100fd366004610e48565b61023e565b005b610128610112366004610eab565b6001600160e01b03191663b6f6b1c560e01b1490565b60405190151581526020015b60405180910390f35b60408051808201909152600f81526e45786368616e67657220302e302e3160881b60208201525b6040516101349190610f25565b61010261017f366004610f4d565b6104e6565b610102610192366004610f6a565b610543565b6101646101a5366004610fba565b6105a9565b6101026105f3565b6000546001600160a01b03165b6040516001600160a01b039091168152602001610134565b6101026101e5366004610fee565b61069d565b60045460405160ff9091168152602001610134565b61010261020d366004611007565b6106ea565b610164610734565b610102610228366004610f4d565b6107d9565b6002546001600160a01b03166101bf565b600081806020019051810190610254919061102a565b825160208085019190912060045460008281526005909352604090922054929350916102839160ff16906110d2565b4311156102d45760045460008281526005602052604090205443916102ad9160ff909116906110d2565b60405163647d550d60e01b8152600481019290925260248201526044015b60405180910390fd5b60025460405163473b057f60e11b81526000916001600160a01b031690638e760afe90610305908890600401610f25565b6000604051808303816000875af1158015610324573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f1916820160405261034c9190810190611115565b9050600081806020019051810190610364919061118b565b84518151919250146103965780518451604051630841559760e41b8152600481019290925260248201526044016102cb565b8060a001516001600160401b0316408160c00151146103e55760c081015160a08201516040516313ffdc7d60e01b815260048101929092526001600160401b03164060248201526044016102cb565b8360c001516001600160a01b03168460a001516001600160a01b031685600001517f1031fb49f3ccdf415485b2e6652f2cdffadb9a0d6374515a949cd64a2a98d7438760200151886040015189606001518a6080015188602001518960e001518a60a001518b60c001518c604001518d606001518e608001516040516104d69b9a999897969594939291909a8b5260208b019990995260408a0197909752606089019590955263ffffffff9390931660808801526001600160401b0391821660a08801521660c086015260e0850152601790810b61010085015290810b6101208401520b6101408201526101600190565b60405180910390a4505050505050565b6104ee6107ed565b600280546001600160a01b0319166001600160a01b0383169081179091556040519081527f8d9f13aae8f2e086b6c478fcb20b85e3f5aab0fcbf26e5d13949ffb23017e539906020015b60405180910390a150565b61054b6107ed565b8060405160200161055c919061122c565b6040516020818303038152906040526003908161057991906112d1565b507fc849f01f1579431074588d2d77603b0c2754e6776e333aaa728de3a13ad967a7816040516105389190610f25565b8051602080830191822060008181526005835260408120548551606095939491936105da929188019091019061102a565b90506105ea816000015183610842565b95945050505050565b6001546001600160a01b031633146106465760405162461bcd60e51b815260206004820152601660248201527526bab9ba10313290383937b837b9b2b21037bbb732b960511b60448201526064016102cb565b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6106a8600143611390565b60008281526005602090815260409182902092909255518281527fa78c3d4234d02c96f4c5223b1956b29eaf57fa2ab28c249e82ee8fecb16fa66b9101610538565b6106f26107ed565b6004805460ff191660ff83169081179091556040519081527f89bc7ef46e5099f5518e68e1171b32bc77b5879fa9f1cc154499aeff37182e9b90602001610538565b60606003805461074390611248565b80601f016020809104026020016040519081016040528092919081815260200182805461076f90611248565b80156107bc5780601f10610791576101008083540402835291602001916107bc565b820191906000526020600020905b81548152906001019060200180831161079f57829003601f168201915b50505050508060200190518101906107d49190611115565b905090565b6107e16107ed565b6107ea8161091e565b50565b6000546001600160a01b031633146108405760405162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b60448201526064016102cb565b565b60606003805461085190611248565b80601f016020809104026020016040519081016040528092919081815260200182805461087d90611248565b80156108ca5780601f1061089f576101008083540402835291602001916108ca565b820191906000526020600020905b8154815290600101906020018083116108ad57829003601f168201915b50505050508060200190518101906108e29190611115565b6108eb846109c7565b6108f4846109de565b604051602001610906939291906113a3565b60405160208183030381529060405290505b92915050565b336001600160a01b038216036109765760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016102cb565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6060610918826109d684610a70565b600101610ada565b606060006109eb83610c7c565b60010190506000816001600160401b03811115610a0a57610a0a610d54565b6040519080825280601f01601f191660200182016040528015610a34576020820181803683370190505b5090508181016020015b600019016f181899199a1a9b1b9c1cb0b131b232b360811b600a86061a8153600a8504945084610a3e57509392505050565b600080608083901c15610a885760809290921c916010015b604083901c15610a9d5760409290921c916008015b602083901c15610ab25760209290921c916004015b601083901c15610ac75760109290921c916002015b600883901c156109185760010192915050565b60606000610ae9836002611418565b610af49060026110d2565b6001600160401b03811115610b0b57610b0b610d54565b6040519080825280601f01601f191660200182016040528015610b35576020820181803683370190505b509050600360fc1b81600081518110610b5057610b50611437565b60200101906001600160f81b031916908160001a905350600f60fb1b81600181518110610b7f57610b7f611437565b60200101906001600160f81b031916908160001a9053506000610ba3846002611418565b610bae9060016110d2565b90505b6001811115610c26576f181899199a1a9b1b9c1cb0b131b232b360811b85600f1660108110610be257610be2611437565b1a60f81b828281518110610bf857610bf8611437565b60200101906001600160f81b031916908160001a90535060049490941c93610c1f8161144d565b9050610bb1565b508315610c755760405162461bcd60e51b815260206004820181905260248201527f537472696e67733a20686578206c656e67746820696e73756666696369656e7460448201526064016102cb565b9392505050565b60008072184f03e93ff9f4daa797ed6e38ed64bf6a1f0160401b8310610cbb5772184f03e93ff9f4daa797ed6e38ed64bf6a1f0160401b830492506040015b6d04ee2d6d415b85acef81000000008310610ce7576d04ee2d6d415b85acef8100000000830492506020015b662386f26fc100008310610d0557662386f26fc10000830492506010015b6305f5e1008310610d1d576305f5e100830492506008015b6127108310610d3157612710830492506004015b60648310610d43576064830492506002015b600a83106109185760010192915050565b634e487b7160e01b600052604160045260246000fd5b60405161010081016001600160401b0381118282101715610d8d57610d8d610d54565b60405290565b604051601f8201601f191681016001600160401b0381118282101715610dbb57610dbb610d54565b604052919050565b60006001600160401b03821115610ddc57610ddc610d54565b50601f01601f191660200190565b6000610dfd610df884610dc3565b610d93565b9050828152838383011115610e1157600080fd5b828260208301376000602084830101529392505050565b600082601f830112610e3957600080fd5b610c7583833560208501610dea565b60008060408385031215610e5b57600080fd5b82356001600160401b0380821115610e7257600080fd5b610e7e86838701610e28565b93506020850135915080821115610e9457600080fd5b50610ea185828601610e28565b9150509250929050565b600060208284031215610ebd57600080fd5b81356001600160e01b031981168114610c7557600080fd5b60005b83811015610ef0578181015183820152602001610ed8565b50506000910152565b60008151808452610f11816020860160208601610ed5565b601f01601f19169290920160200192915050565b602081526000610c756020830184610ef9565b6001600160a01b03811681146107ea57600080fd5b600060208284031215610f5f57600080fd5b8135610c7581610f38565b600060208284031215610f7c57600080fd5b81356001600160401b03811115610f9257600080fd5b8201601f81018413610fa357600080fd5b610fb284823560208401610dea565b949350505050565b600060208284031215610fcc57600080fd5b81356001600160401b03811115610fe257600080fd5b610fb284828501610e28565b60006020828403121561100057600080fd5b5035919050565b60006020828403121561101957600080fd5b813560ff81168114610c7557600080fd5b600060e0828403121561103c57600080fd5b60405160e081018181106001600160401b038211171561105e5761105e610d54565b8060405250825181526020830151602082015260408301516040820152606083015160608201526080830151608082015260a083015161109d81610f38565b60a082015260c08301516110b081610f38565b60c08201529392505050565b634e487b7160e01b600052601160045260246000fd5b80820180821115610918576109186110bc565b60006110f3610df884610dc3565b905082815283838301111561110757600080fd5b610c75836020830184610ed5565b60006020828403121561112757600080fd5b81516001600160401b0381111561113d57600080fd5b8201601f8101841361114e57600080fd5b610fb2848251602084016110e5565b8051601781900b811461116f57600080fd5b919050565b80516001600160401b038116811461116f57600080fd5b6000610100828403121561119e57600080fd5b6111a6610d6a565b82518152602083015163ffffffff811681146111c157600080fd5b60208201526111d26040840161115d565b60408201526111e36060840161115d565b60608201526111f46080840161115d565b608082015261120560a08401611174565b60a082015260c083015160c082015261122060e08401611174565b60e08201529392505050565b6000825161123e818460208701610ed5565b9190910192915050565b600181811c9082168061125c57607f821691505b60208210810361127c57634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156112cc57600081815260208120601f850160051c810160208610156112a95750805b601f850160051c820191505b818110156112c8578281556001016112b5565b5050505b505050565b81516001600160401b038111156112ea576112ea610d54565b6112fe816112f88454611248565b84611282565b602080601f831160018114611333576000841561131b5750858301515b600019600386901b1c1916600185901b1785556112c8565b600085815260208120601f198616915b8281101561136257888601518255948401946001909101908401611343565b50858210156113805787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b81810381811115610918576109186110bc565b600084516113b5818460208901610ed5565b6a3f6665656449444865783d60a81b90830190815284516113dd81600b840160208901610ed5565b6c26626c6f636b4e756d6265723d60981b600b9290910191820152835161140b816018840160208801610ed5565b0160180195945050505050565b6000816000190483118215151615611432576114326110bc565b500290565b634e487b7160e01b600052603260045260246000fd5b60008161145c5761145c6110bc565b50600019019056fea264697066735822122007ec37859fd4e499efe622b90d9cfc4f35a6fabe2bf50fd12f76b8448d1ad34c64736f6c63430008100033",
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
	parsed, err := ExchangerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
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
	FeedId                [32]byte
	CurrencySrc           [32]byte
	CurrencyDst           [32]byte
	AmountSrc             *big.Int
	MinAmountDst          *big.Int
	Sender                common.Address
	Receiver              common.Address
	ObservationsTimestamp uint32
	BlocknumberLowerBound uint64
	BlocknumberUpperBound uint64
	UpperBlockhash        [32]byte
	Median                *big.Int
	Bid                   *big.Int
	Ask                   *big.Int
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterTradeExecuted is a free log retrieval operation binding the contract event 0x1031fb49f3ccdf415485b2e6652f2cdffadb9a0d6374515a949cd64a2a98d743.
//
// Solidity: event TradeExecuted(bytes32 indexed feedId, bytes32 currencySrc, bytes32 currencyDst, uint256 amountSrc, uint256 minAmountDst, address indexed sender, address indexed receiver, uint32 observationsTimestamp, uint64 blocknumberLowerBound, uint64 blocknumberUpperBound, bytes32 upperBlockhash, int192 median, int192 bid, int192 ask)
func (_Exchanger *ExchangerFilterer) FilterTradeExecuted(opts *bind.FilterOpts, feedId [][32]byte, sender []common.Address, receiver []common.Address) (*ExchangerTradeExecutedIterator, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _Exchanger.contract.FilterLogs(opts, "TradeExecuted", feedIdRule, senderRule, receiverRule)
	if err != nil {
		return nil, err
	}
	return &ExchangerTradeExecutedIterator{contract: _Exchanger.contract, event: "TradeExecuted", logs: logs, sub: sub}, nil
}

// WatchTradeExecuted is a free log subscription operation binding the contract event 0x1031fb49f3ccdf415485b2e6652f2cdffadb9a0d6374515a949cd64a2a98d743.
//
// Solidity: event TradeExecuted(bytes32 indexed feedId, bytes32 currencySrc, bytes32 currencyDst, uint256 amountSrc, uint256 minAmountDst, address indexed sender, address indexed receiver, uint32 observationsTimestamp, uint64 blocknumberLowerBound, uint64 blocknumberUpperBound, bytes32 upperBlockhash, int192 median, int192 bid, int192 ask)
func (_Exchanger *ExchangerFilterer) WatchTradeExecuted(opts *bind.WatchOpts, sink chan<- *ExchangerTradeExecuted, feedId [][32]byte, sender []common.Address, receiver []common.Address) (event.Subscription, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _Exchanger.contract.WatchLogs(opts, "TradeExecuted", feedIdRule, senderRule, receiverRule)
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

// ParseTradeExecuted is a log parse operation binding the contract event 0x1031fb49f3ccdf415485b2e6652f2cdffadb9a0d6374515a949cd64a2a98d743.
//
// Solidity: event TradeExecuted(bytes32 indexed feedId, bytes32 currencySrc, bytes32 currencyDst, uint256 amountSrc, uint256 minAmountDst, address indexed sender, address indexed receiver, uint32 observationsTimestamp, uint64 blocknumberLowerBound, uint64 blocknumberUpperBound, bytes32 upperBlockhash, int192 median, int192 bid, int192 ask)
func (_Exchanger *ExchangerFilterer) ParseTradeExecuted(log types.Log) (*ExchangerTradeExecuted, error) {
	event := new(ExchangerTradeExecuted)
	if err := _Exchanger.contract.UnpackLog(event, "TradeExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
