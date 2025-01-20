// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package verifier_proxy_v0_5_0

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

type CommonAddressAndWeight struct {
	Addr   common.Address
	Weight uint64
}

var VerifierProxyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"accessController\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessForbidden\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BadVerification\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"verifier\",\"type\":\"address\"}],\"name\":\"ConfigDigestAlreadySet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FeeManagerInvalid\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"verifier\",\"type\":\"address\"}],\"name\":\"VerifierAlreadyInitialized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"VerifierInvalid\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"VerifierNotFound\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldAccessController\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAccessController\",\"type\":\"address\"}],\"name\":\"AccessControllerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldFeeManager\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newFeeManager\",\"type\":\"address\"}],\"name\":\"FeeManagerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"verifierAddress\",\"type\":\"address\"}],\"name\":\"VerifierInitialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"oldConfigDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"newConfigDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"verifierAddress\",\"type\":\"address\"}],\"name\":\"VerifierSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"verifierAddress\",\"type\":\"address\"}],\"name\":\"VerifierUnset\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"getVerifier\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"verifierAddress\",\"type\":\"address\"}],\"name\":\"initializeVerifier\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_accessController\",\"outputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_feeManager\",\"outputs\":[{\"internalType\":\"contractIVerifierFeeManager\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"accessController\",\"type\":\"address\"}],\"name\":\"setAccessController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIVerifierFeeManager\",\"name\":\"feeManager\",\"type\":\"address\"}],\"name\":\"setFeeManager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"currentConfigDigest\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newConfigDigest\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"weight\",\"type\":\"uint64\"}],\"internalType\":\"structCommon.AddressAndWeight[]\",\"name\":\"addressesAndWeights\",\"type\":\"tuple[]\"}],\"name\":\"setVerifier\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"unsetVerifier\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"parameterPayload\",\"type\":\"bytes\"}],\"name\":\"verify\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"payloads\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"parameterPayload\",\"type\":\"bytes\"}],\"name\":\"verifyBulk\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"verifiedReports\",\"type\":\"bytes[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162001d3638038062001d36833981016040819052620000349162000193565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000e8565b5050600480546001600160a01b0319166001600160a01b03939093169290921790915550620001c5565b336001600160a01b03821603620001425760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600060208284031215620001a657600080fd5b81516001600160a01b0381168114620001be57600080fd5b9392505050565b611b6180620001d56000396000f3fe6080604052600436106100dd5760003560e01c806394ba28461161007f578063f08391d811610059578063f08391d8146102be578063f2fde38b146102de578063f7e83aee146102fe578063f873a61c1461031157600080fd5b806394ba28461461022e578063b011b2471461025b578063eeb7b2481461027b57600080fd5b80636e914094116100bb5780636e914094146101ae57806379ba5097146101ce5780638c2a4d53146101e35780638da5cb5b1461020357600080fd5b8063181f5a77146100e257806338416b5b1461013a578063472d35b91461018c575b600080fd5b3480156100ee57600080fd5b5060408051808201909152601381527f566572696669657250726f787920322e302e300000000000000000000000000060208201525b60405161013191906113c3565b60405180910390f35b34801561014657600080fd5b506005546101679073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610131565b34801561019857600080fd5b506101ac6101a73660046113ff565b610331565b005b3480156101ba57600080fd5b506101ac6101c936600461141c565b6105a9565b3480156101da57600080fd5b506101ac61069a565b3480156101ef57600080fd5b506101ac6101fe3660046113ff565b610797565b34801561020f57600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff16610167565b34801561023a57600080fd5b506004546101679073ffffffffffffffffffffffffffffffffffffffff1681565b34801561026757600080fd5b506101ac610276366004611435565b6109c8565b34801561028757600080fd5b5061016761029636600461141c565b60009081526003602052604090205473ffffffffffffffffffffffffffffffffffffffff1690565b3480156102ca57600080fd5b506101ac6102d93660046113ff565b610bee565b3480156102ea57600080fd5b506101ac6102f93660046113ff565b610c75565b61012461030c366004611501565b610c89565b61032461031f36600461156d565b610e43565b60405161013191906115ee565b6103396110a7565b73ffffffffffffffffffffffffffffffffffffffff8116610386576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040517f01ffc9a70000000000000000000000000000000000000000000000000000000081527fdba45fe000000000000000000000000000000000000000000000000000000000600482015273ffffffffffffffffffffffffffffffffffffffff8216906301ffc9a790602401602060405180830381865afa158015610410573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610434919061166e565b15806104eb57506040517f01ffc9a70000000000000000000000000000000000000000000000000000000081527f6c2f1a1700000000000000000000000000000000000000000000000000000000600482015273ffffffffffffffffffffffffffffffffffffffff8216906301ffc9a790602401602060405180830381865afa1580156104c5573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104e9919061166e565b155b15610522576040517f8238941900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6005805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff000000000000000000000000000000000000000083168117909355604080519190921680825260208201939093527f04628abcaa6b1674651352125cb94b65b289145bc2bc4d67720bb7d966372f0391015b60405180910390a15050565b6105b16110a7565b60008181526003602052604090205473ffffffffffffffffffffffffffffffffffffffff1680610615576040517fb151802b000000000000000000000000000000000000000000000000000000008152600481018390526024015b60405180910390fd5b6000828152600360205260409081902080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169055517f11dc15c4b8ac2b183166cc8427e5385a5ece8308217a4217338c6a7614845c4c9061059d908490849091825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b60015473ffffffffffffffffffffffffffffffffffffffff16331461071b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015260640161060c565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b61079f6110a7565b8073ffffffffffffffffffffffffffffffffffffffff81166107ed576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040517f01ffc9a70000000000000000000000000000000000000000000000000000000081527f3d3ac1b500000000000000000000000000000000000000000000000000000000600482015273ffffffffffffffffffffffffffffffffffffffff8216906301ffc9a790602401602060405180830381865afa158015610877573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061089b919061166e565b6108d1576040517f75b0527a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff821660009081526002602052604090205460ff1615610949576040517f4e01ccfd00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8316600482015260240161060c565b73ffffffffffffffffffffffffffffffffffffffff821660008181526002602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016600117905590519182527f1f2cd7c97f4d801b5efe26cc409617c1fd6c5ef786e79aacb90af40923e4e8e9910161059d565b600083815260036020526040902054839073ffffffffffffffffffffffffffffffffffffffff168015610a46576040517f375d1fe60000000000000000000000000000000000000000000000000000000081526004810183905273ffffffffffffffffffffffffffffffffffffffff8216602482015260440161060c565b3360009081526002602052604090205460ff16610a8f576040517fef67f5d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600085815260036020526040902080547fffffffffffffffffffffffff000000000000000000000000000000000000000016331790558215610ba75760055473ffffffffffffffffffffffffffffffffffffffff16610b1a576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6005546040517ff65df96200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091169063f65df96290610b7490889088908890600401611690565b600060405180830381600087803b158015610b8e57600080fd5b505af1158015610ba2573d6000803e3d6000fd5b505050505b6040805187815260208101879052338183015290517fbeb513e532542a562ac35699e7cd9ae7d198dcd3eee15bada6c857d28ceaddcf9181900360600190a1505050505050565b610bf66110a7565b6004805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff000000000000000000000000000000000000000083168117909355604080519190921680825260208201939093527f953e92b1a6442e9c3242531154a3f6f6eb00b4e9c719ba8118fa6235e4ce89b6910161059d565b610c7d6110a7565b610c868161112a565b50565b60045460609073ffffffffffffffffffffffffffffffffffffffff168015801590610d4957506040517f6b14daf800000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821690636b14daf890610d069033906000903690600401611762565b602060405180830381865afa158015610d23573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d47919061166e565b155b15610d80576040517fef67f5d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60055473ffffffffffffffffffffffffffffffffffffffff168015610e2e576040517fdba45fe000000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82169063dba45fe0903490610dfb908b908b908b908b90339060040161179b565b6000604051808303818588803b158015610e1457600080fd5b505af1158015610e28573d6000803e3d6000fd5b50505050505b610e38878761121f565b979650505050505050565b60045460609073ffffffffffffffffffffffffffffffffffffffff168015801590610f0357506040517f6b14daf800000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821690636b14daf890610ec09033906000903690600401611762565b602060405180830381865afa158015610edd573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610f01919061166e565b155b15610f3a576040517fef67f5d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60055473ffffffffffffffffffffffffffffffffffffffff168015610fe8576040517f6c2f1a1700000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821690636c2f1a17903490610fb5908b908b908b908b9033906004016117eb565b6000604051808303818588803b158015610fce57600080fd5b505af1158015610fe2573d6000803e3d6000fd5b50505050505b8567ffffffffffffffff811115611001576110016118fc565b60405190808252806020026020018201604052801561103457816020015b606081526020019060019003908161101f5790505b50925060005b8681101561109c5761106e8888838181106110575761105761192b565b9050602002810190611069919061195a565b61121f565b8482815181106110805761108061192b565b602002602001018190525080611095906119bf565b905061103a565b505050949350505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314611128576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161060c565b565b3373ffffffffffffffffffffffffffffffffffffffff8216036111a9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161060c565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6060600061122d8385611a1e565b60008181526003602052604090205490915073ffffffffffffffffffffffffffffffffffffffff168061128f576040517fb151802b0000000000000000000000000000000000000000000000000000000081526004810183905260240161060c565b6040517f3d3ac1b500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821690633d3ac1b5906112e590889088903390600401611a5a565b6000604051808303816000875af1158015611304573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405261134a9190810190611a94565b925050505b92915050565b60005b83811015611370578181015183820152602001611358565b50506000910152565b60008151808452611391816020860160208601611355565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6020815260006113d66020830184611379565b9392505050565b73ffffffffffffffffffffffffffffffffffffffff81168114610c8657600080fd5b60006020828403121561141157600080fd5b81356113d6816113dd565b60006020828403121561142e57600080fd5b5035919050565b6000806000806060858703121561144b57600080fd5b8435935060208501359250604085013567ffffffffffffffff8082111561147157600080fd5b818701915087601f83011261148557600080fd5b81358181111561149457600080fd5b8860208260061b85010111156114a957600080fd5b95989497505060200194505050565b60008083601f8401126114ca57600080fd5b50813567ffffffffffffffff8111156114e257600080fd5b6020830191508360208285010111156114fa57600080fd5b9250929050565b6000806000806040858703121561151757600080fd5b843567ffffffffffffffff8082111561152f57600080fd5b61153b888389016114b8565b9096509450602087013591508082111561155457600080fd5b50611561878288016114b8565b95989497509550505050565b6000806000806040858703121561158357600080fd5b843567ffffffffffffffff8082111561159b57600080fd5b818701915087601f8301126115af57600080fd5b8135818111156115be57600080fd5b8860208260051b85010111156115d357600080fd5b60209283019650945090860135908082111561155457600080fd5b6000602080830181845280855180835260408601915060408160051b870101925083870160005b82811015611661577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc088860301845261164f858351611379565b94509285019290850190600101611615565b5092979650505050505050565b60006020828403121561168057600080fd5b815180151581146113d657600080fd5b838152604060208083018290528282018490526000919085906060850184805b8881101561170a5784356116c3816113dd565b73ffffffffffffffffffffffffffffffffffffffff1683528484013567ffffffffffffffff81168082146116f5578384fd5b848601525093850193918501916001016116b0565b50909998505050505050505050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b73ffffffffffffffffffffffffffffffffffffffff84168152604060208201526000611792604083018486611719565b95945050505050565b6060815260006117af606083018789611719565b82810360208401526117c2818688611719565b91505073ffffffffffffffffffffffffffffffffffffffff831660408301529695505050505050565b6060808252810185905260006080600587901b8301810190830188835b898110156118b7577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8086850301835281357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18c360301811261186957600080fd5b8b01602081810191359067ffffffffffffffff82111561188857600080fd5b81360383131561189757600080fd5b6118a2878385611719565b96509485019493909301925050600101611808565b50505082810360208401526118cd818688611719565b9150506118f2604083018473ffffffffffffffffffffffffffffffffffffffff169052565b9695505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261198f57600080fd5b83018035915067ffffffffffffffff8211156119aa57600080fd5b6020019150368190038213156114fa57600080fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203611a17577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b5060010190565b8035602083101561134f577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff602084900360031b1b1692915050565b604081526000611a6e604083018587611719565b905073ffffffffffffffffffffffffffffffffffffffff83166020830152949350505050565b600060208284031215611aa657600080fd5b815167ffffffffffffffff80821115611abe57600080fd5b818401915084601f830112611ad257600080fd5b815181811115611ae457611ae46118fc565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f01168101908382118183101715611b2a57611b2a6118fc565b81604052828152876020848701011115611b4357600080fd5b610e3883602083016020880161135556fea164736f6c6343000813000a",
}

var VerifierProxyABI = VerifierProxyMetaData.ABI

var VerifierProxyBin = VerifierProxyMetaData.Bin

func DeployVerifierProxy(auth *bind.TransactOpts, backend bind.ContractBackend, accessController common.Address) (common.Address, *types.Transaction, *VerifierProxy, error) {
	parsed, err := VerifierProxyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VerifierProxyBin), backend, accessController)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VerifierProxy{address: address, abi: *parsed, VerifierProxyCaller: VerifierProxyCaller{contract: contract}, VerifierProxyTransactor: VerifierProxyTransactor{contract: contract}, VerifierProxyFilterer: VerifierProxyFilterer{contract: contract}}, nil
}

type VerifierProxy struct {
	address common.Address
	abi     abi.ABI
	VerifierProxyCaller
	VerifierProxyTransactor
	VerifierProxyFilterer
}

type VerifierProxyCaller struct {
	contract *bind.BoundContract
}

type VerifierProxyTransactor struct {
	contract *bind.BoundContract
}

type VerifierProxyFilterer struct {
	contract *bind.BoundContract
}

type VerifierProxySession struct {
	Contract     *VerifierProxy
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VerifierProxyCallerSession struct {
	Contract *VerifierProxyCaller
	CallOpts bind.CallOpts
}

type VerifierProxyTransactorSession struct {
	Contract     *VerifierProxyTransactor
	TransactOpts bind.TransactOpts
}

type VerifierProxyRaw struct {
	Contract *VerifierProxy
}

type VerifierProxyCallerRaw struct {
	Contract *VerifierProxyCaller
}

type VerifierProxyTransactorRaw struct {
	Contract *VerifierProxyTransactor
}

func NewVerifierProxy(address common.Address, backend bind.ContractBackend) (*VerifierProxy, error) {
	abi, err := abi.JSON(strings.NewReader(VerifierProxyABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVerifierProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VerifierProxy{address: address, abi: abi, VerifierProxyCaller: VerifierProxyCaller{contract: contract}, VerifierProxyTransactor: VerifierProxyTransactor{contract: contract}, VerifierProxyFilterer: VerifierProxyFilterer{contract: contract}}, nil
}

func NewVerifierProxyCaller(address common.Address, caller bind.ContractCaller) (*VerifierProxyCaller, error) {
	contract, err := bindVerifierProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VerifierProxyCaller{contract: contract}, nil
}

func NewVerifierProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*VerifierProxyTransactor, error) {
	contract, err := bindVerifierProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VerifierProxyTransactor{contract: contract}, nil
}

func NewVerifierProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*VerifierProxyFilterer, error) {
	contract, err := bindVerifierProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VerifierProxyFilterer{contract: contract}, nil
}

func bindVerifierProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VerifierProxyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VerifierProxy *VerifierProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VerifierProxy.Contract.VerifierProxyCaller.contract.Call(opts, result, method, params...)
}

func (_VerifierProxy *VerifierProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VerifierProxy.Contract.VerifierProxyTransactor.contract.Transfer(opts)
}

func (_VerifierProxy *VerifierProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VerifierProxy.Contract.VerifierProxyTransactor.contract.Transact(opts, method, params...)
}

func (_VerifierProxy *VerifierProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VerifierProxy.Contract.contract.Call(opts, result, method, params...)
}

func (_VerifierProxy *VerifierProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VerifierProxy.Contract.contract.Transfer(opts)
}

func (_VerifierProxy *VerifierProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VerifierProxy.Contract.contract.Transact(opts, method, params...)
}

func (_VerifierProxy *VerifierProxyCaller) GetVerifier(opts *bind.CallOpts, configDigest [32]byte) (common.Address, error) {
	var out []interface{}
	err := _VerifierProxy.contract.Call(opts, &out, "getVerifier", configDigest)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VerifierProxy *VerifierProxySession) GetVerifier(configDigest [32]byte) (common.Address, error) {
	return _VerifierProxy.Contract.GetVerifier(&_VerifierProxy.CallOpts, configDigest)
}

func (_VerifierProxy *VerifierProxyCallerSession) GetVerifier(configDigest [32]byte) (common.Address, error) {
	return _VerifierProxy.Contract.GetVerifier(&_VerifierProxy.CallOpts, configDigest)
}

func (_VerifierProxy *VerifierProxyCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VerifierProxy.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VerifierProxy *VerifierProxySession) Owner() (common.Address, error) {
	return _VerifierProxy.Contract.Owner(&_VerifierProxy.CallOpts)
}

func (_VerifierProxy *VerifierProxyCallerSession) Owner() (common.Address, error) {
	return _VerifierProxy.Contract.Owner(&_VerifierProxy.CallOpts)
}

func (_VerifierProxy *VerifierProxyCaller) SAccessController(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VerifierProxy.contract.Call(opts, &out, "s_accessController")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VerifierProxy *VerifierProxySession) SAccessController() (common.Address, error) {
	return _VerifierProxy.Contract.SAccessController(&_VerifierProxy.CallOpts)
}

func (_VerifierProxy *VerifierProxyCallerSession) SAccessController() (common.Address, error) {
	return _VerifierProxy.Contract.SAccessController(&_VerifierProxy.CallOpts)
}

func (_VerifierProxy *VerifierProxyCaller) SFeeManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VerifierProxy.contract.Call(opts, &out, "s_feeManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VerifierProxy *VerifierProxySession) SFeeManager() (common.Address, error) {
	return _VerifierProxy.Contract.SFeeManager(&_VerifierProxy.CallOpts)
}

func (_VerifierProxy *VerifierProxyCallerSession) SFeeManager() (common.Address, error) {
	return _VerifierProxy.Contract.SFeeManager(&_VerifierProxy.CallOpts)
}

func (_VerifierProxy *VerifierProxyCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _VerifierProxy.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_VerifierProxy *VerifierProxySession) TypeAndVersion() (string, error) {
	return _VerifierProxy.Contract.TypeAndVersion(&_VerifierProxy.CallOpts)
}

func (_VerifierProxy *VerifierProxyCallerSession) TypeAndVersion() (string, error) {
	return _VerifierProxy.Contract.TypeAndVersion(&_VerifierProxy.CallOpts)
}

func (_VerifierProxy *VerifierProxyTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VerifierProxy.contract.Transact(opts, "acceptOwnership")
}

func (_VerifierProxy *VerifierProxySession) AcceptOwnership() (*types.Transaction, error) {
	return _VerifierProxy.Contract.AcceptOwnership(&_VerifierProxy.TransactOpts)
}

func (_VerifierProxy *VerifierProxyTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VerifierProxy.Contract.AcceptOwnership(&_VerifierProxy.TransactOpts)
}

func (_VerifierProxy *VerifierProxyTransactor) InitializeVerifier(opts *bind.TransactOpts, verifierAddress common.Address) (*types.Transaction, error) {
	return _VerifierProxy.contract.Transact(opts, "initializeVerifier", verifierAddress)
}

func (_VerifierProxy *VerifierProxySession) InitializeVerifier(verifierAddress common.Address) (*types.Transaction, error) {
	return _VerifierProxy.Contract.InitializeVerifier(&_VerifierProxy.TransactOpts, verifierAddress)
}

func (_VerifierProxy *VerifierProxyTransactorSession) InitializeVerifier(verifierAddress common.Address) (*types.Transaction, error) {
	return _VerifierProxy.Contract.InitializeVerifier(&_VerifierProxy.TransactOpts, verifierAddress)
}

func (_VerifierProxy *VerifierProxyTransactor) SetAccessController(opts *bind.TransactOpts, accessController common.Address) (*types.Transaction, error) {
	return _VerifierProxy.contract.Transact(opts, "setAccessController", accessController)
}

func (_VerifierProxy *VerifierProxySession) SetAccessController(accessController common.Address) (*types.Transaction, error) {
	return _VerifierProxy.Contract.SetAccessController(&_VerifierProxy.TransactOpts, accessController)
}

func (_VerifierProxy *VerifierProxyTransactorSession) SetAccessController(accessController common.Address) (*types.Transaction, error) {
	return _VerifierProxy.Contract.SetAccessController(&_VerifierProxy.TransactOpts, accessController)
}

func (_VerifierProxy *VerifierProxyTransactor) SetFeeManager(opts *bind.TransactOpts, feeManager common.Address) (*types.Transaction, error) {
	return _VerifierProxy.contract.Transact(opts, "setFeeManager", feeManager)
}

func (_VerifierProxy *VerifierProxySession) SetFeeManager(feeManager common.Address) (*types.Transaction, error) {
	return _VerifierProxy.Contract.SetFeeManager(&_VerifierProxy.TransactOpts, feeManager)
}

func (_VerifierProxy *VerifierProxyTransactorSession) SetFeeManager(feeManager common.Address) (*types.Transaction, error) {
	return _VerifierProxy.Contract.SetFeeManager(&_VerifierProxy.TransactOpts, feeManager)
}

func (_VerifierProxy *VerifierProxyTransactor) SetVerifier(opts *bind.TransactOpts, currentConfigDigest [32]byte, newConfigDigest [32]byte, addressesAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _VerifierProxy.contract.Transact(opts, "setVerifier", currentConfigDigest, newConfigDigest, addressesAndWeights)
}

func (_VerifierProxy *VerifierProxySession) SetVerifier(currentConfigDigest [32]byte, newConfigDigest [32]byte, addressesAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _VerifierProxy.Contract.SetVerifier(&_VerifierProxy.TransactOpts, currentConfigDigest, newConfigDigest, addressesAndWeights)
}

func (_VerifierProxy *VerifierProxyTransactorSession) SetVerifier(currentConfigDigest [32]byte, newConfigDigest [32]byte, addressesAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _VerifierProxy.Contract.SetVerifier(&_VerifierProxy.TransactOpts, currentConfigDigest, newConfigDigest, addressesAndWeights)
}

func (_VerifierProxy *VerifierProxyTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VerifierProxy.contract.Transact(opts, "transferOwnership", to)
}

func (_VerifierProxy *VerifierProxySession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VerifierProxy.Contract.TransferOwnership(&_VerifierProxy.TransactOpts, to)
}

func (_VerifierProxy *VerifierProxyTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VerifierProxy.Contract.TransferOwnership(&_VerifierProxy.TransactOpts, to)
}

func (_VerifierProxy *VerifierProxyTransactor) UnsetVerifier(opts *bind.TransactOpts, configDigest [32]byte) (*types.Transaction, error) {
	return _VerifierProxy.contract.Transact(opts, "unsetVerifier", configDigest)
}

func (_VerifierProxy *VerifierProxySession) UnsetVerifier(configDigest [32]byte) (*types.Transaction, error) {
	return _VerifierProxy.Contract.UnsetVerifier(&_VerifierProxy.TransactOpts, configDigest)
}

func (_VerifierProxy *VerifierProxyTransactorSession) UnsetVerifier(configDigest [32]byte) (*types.Transaction, error) {
	return _VerifierProxy.Contract.UnsetVerifier(&_VerifierProxy.TransactOpts, configDigest)
}

func (_VerifierProxy *VerifierProxyTransactor) Verify(opts *bind.TransactOpts, payload []byte, parameterPayload []byte) (*types.Transaction, error) {
	return _VerifierProxy.contract.Transact(opts, "verify", payload, parameterPayload)
}

func (_VerifierProxy *VerifierProxySession) Verify(payload []byte, parameterPayload []byte) (*types.Transaction, error) {
	return _VerifierProxy.Contract.Verify(&_VerifierProxy.TransactOpts, payload, parameterPayload)
}

func (_VerifierProxy *VerifierProxyTransactorSession) Verify(payload []byte, parameterPayload []byte) (*types.Transaction, error) {
	return _VerifierProxy.Contract.Verify(&_VerifierProxy.TransactOpts, payload, parameterPayload)
}

func (_VerifierProxy *VerifierProxyTransactor) VerifyBulk(opts *bind.TransactOpts, payloads [][]byte, parameterPayload []byte) (*types.Transaction, error) {
	return _VerifierProxy.contract.Transact(opts, "verifyBulk", payloads, parameterPayload)
}

func (_VerifierProxy *VerifierProxySession) VerifyBulk(payloads [][]byte, parameterPayload []byte) (*types.Transaction, error) {
	return _VerifierProxy.Contract.VerifyBulk(&_VerifierProxy.TransactOpts, payloads, parameterPayload)
}

func (_VerifierProxy *VerifierProxyTransactorSession) VerifyBulk(payloads [][]byte, parameterPayload []byte) (*types.Transaction, error) {
	return _VerifierProxy.Contract.VerifyBulk(&_VerifierProxy.TransactOpts, payloads, parameterPayload)
}

type VerifierProxyAccessControllerSetIterator struct {
	Event *VerifierProxyAccessControllerSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifierProxyAccessControllerSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifierProxyAccessControllerSet)
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
		it.Event = new(VerifierProxyAccessControllerSet)
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

func (it *VerifierProxyAccessControllerSetIterator) Error() error {
	return it.fail
}

func (it *VerifierProxyAccessControllerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifierProxyAccessControllerSet struct {
	OldAccessController common.Address
	NewAccessController common.Address
	Raw                 types.Log
}

func (_VerifierProxy *VerifierProxyFilterer) FilterAccessControllerSet(opts *bind.FilterOpts) (*VerifierProxyAccessControllerSetIterator, error) {

	logs, sub, err := _VerifierProxy.contract.FilterLogs(opts, "AccessControllerSet")
	if err != nil {
		return nil, err
	}
	return &VerifierProxyAccessControllerSetIterator{contract: _VerifierProxy.contract, event: "AccessControllerSet", logs: logs, sub: sub}, nil
}

func (_VerifierProxy *VerifierProxyFilterer) WatchAccessControllerSet(opts *bind.WatchOpts, sink chan<- *VerifierProxyAccessControllerSet) (event.Subscription, error) {

	logs, sub, err := _VerifierProxy.contract.WatchLogs(opts, "AccessControllerSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifierProxyAccessControllerSet)
				if err := _VerifierProxy.contract.UnpackLog(event, "AccessControllerSet", log); err != nil {
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

func (_VerifierProxy *VerifierProxyFilterer) ParseAccessControllerSet(log types.Log) (*VerifierProxyAccessControllerSet, error) {
	event := new(VerifierProxyAccessControllerSet)
	if err := _VerifierProxy.contract.UnpackLog(event, "AccessControllerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifierProxyFeeManagerSetIterator struct {
	Event *VerifierProxyFeeManagerSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifierProxyFeeManagerSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifierProxyFeeManagerSet)
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
		it.Event = new(VerifierProxyFeeManagerSet)
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

func (it *VerifierProxyFeeManagerSetIterator) Error() error {
	return it.fail
}

func (it *VerifierProxyFeeManagerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifierProxyFeeManagerSet struct {
	OldFeeManager common.Address
	NewFeeManager common.Address
	Raw           types.Log
}

func (_VerifierProxy *VerifierProxyFilterer) FilterFeeManagerSet(opts *bind.FilterOpts) (*VerifierProxyFeeManagerSetIterator, error) {

	logs, sub, err := _VerifierProxy.contract.FilterLogs(opts, "FeeManagerSet")
	if err != nil {
		return nil, err
	}
	return &VerifierProxyFeeManagerSetIterator{contract: _VerifierProxy.contract, event: "FeeManagerSet", logs: logs, sub: sub}, nil
}

func (_VerifierProxy *VerifierProxyFilterer) WatchFeeManagerSet(opts *bind.WatchOpts, sink chan<- *VerifierProxyFeeManagerSet) (event.Subscription, error) {

	logs, sub, err := _VerifierProxy.contract.WatchLogs(opts, "FeeManagerSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifierProxyFeeManagerSet)
				if err := _VerifierProxy.contract.UnpackLog(event, "FeeManagerSet", log); err != nil {
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

func (_VerifierProxy *VerifierProxyFilterer) ParseFeeManagerSet(log types.Log) (*VerifierProxyFeeManagerSet, error) {
	event := new(VerifierProxyFeeManagerSet)
	if err := _VerifierProxy.contract.UnpackLog(event, "FeeManagerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifierProxyOwnershipTransferRequestedIterator struct {
	Event *VerifierProxyOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifierProxyOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifierProxyOwnershipTransferRequested)
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
		it.Event = new(VerifierProxyOwnershipTransferRequested)
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

func (it *VerifierProxyOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VerifierProxyOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifierProxyOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VerifierProxy *VerifierProxyFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VerifierProxyOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VerifierProxy.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VerifierProxyOwnershipTransferRequestedIterator{contract: _VerifierProxy.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VerifierProxy *VerifierProxyFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VerifierProxyOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VerifierProxy.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifierProxyOwnershipTransferRequested)
				if err := _VerifierProxy.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VerifierProxy *VerifierProxyFilterer) ParseOwnershipTransferRequested(log types.Log) (*VerifierProxyOwnershipTransferRequested, error) {
	event := new(VerifierProxyOwnershipTransferRequested)
	if err := _VerifierProxy.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifierProxyOwnershipTransferredIterator struct {
	Event *VerifierProxyOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifierProxyOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifierProxyOwnershipTransferred)
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
		it.Event = new(VerifierProxyOwnershipTransferred)
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

func (it *VerifierProxyOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VerifierProxyOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifierProxyOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VerifierProxy *VerifierProxyFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VerifierProxyOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VerifierProxy.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VerifierProxyOwnershipTransferredIterator{contract: _VerifierProxy.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VerifierProxy *VerifierProxyFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VerifierProxyOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VerifierProxy.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifierProxyOwnershipTransferred)
				if err := _VerifierProxy.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VerifierProxy *VerifierProxyFilterer) ParseOwnershipTransferred(log types.Log) (*VerifierProxyOwnershipTransferred, error) {
	event := new(VerifierProxyOwnershipTransferred)
	if err := _VerifierProxy.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifierProxyVerifierInitializedIterator struct {
	Event *VerifierProxyVerifierInitialized

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifierProxyVerifierInitializedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifierProxyVerifierInitialized)
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
		it.Event = new(VerifierProxyVerifierInitialized)
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

func (it *VerifierProxyVerifierInitializedIterator) Error() error {
	return it.fail
}

func (it *VerifierProxyVerifierInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifierProxyVerifierInitialized struct {
	VerifierAddress common.Address
	Raw             types.Log
}

func (_VerifierProxy *VerifierProxyFilterer) FilterVerifierInitialized(opts *bind.FilterOpts) (*VerifierProxyVerifierInitializedIterator, error) {

	logs, sub, err := _VerifierProxy.contract.FilterLogs(opts, "VerifierInitialized")
	if err != nil {
		return nil, err
	}
	return &VerifierProxyVerifierInitializedIterator{contract: _VerifierProxy.contract, event: "VerifierInitialized", logs: logs, sub: sub}, nil
}

func (_VerifierProxy *VerifierProxyFilterer) WatchVerifierInitialized(opts *bind.WatchOpts, sink chan<- *VerifierProxyVerifierInitialized) (event.Subscription, error) {

	logs, sub, err := _VerifierProxy.contract.WatchLogs(opts, "VerifierInitialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifierProxyVerifierInitialized)
				if err := _VerifierProxy.contract.UnpackLog(event, "VerifierInitialized", log); err != nil {
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

func (_VerifierProxy *VerifierProxyFilterer) ParseVerifierInitialized(log types.Log) (*VerifierProxyVerifierInitialized, error) {
	event := new(VerifierProxyVerifierInitialized)
	if err := _VerifierProxy.contract.UnpackLog(event, "VerifierInitialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifierProxyVerifierSetIterator struct {
	Event *VerifierProxyVerifierSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifierProxyVerifierSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifierProxyVerifierSet)
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
		it.Event = new(VerifierProxyVerifierSet)
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

func (it *VerifierProxyVerifierSetIterator) Error() error {
	return it.fail
}

func (it *VerifierProxyVerifierSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifierProxyVerifierSet struct {
	OldConfigDigest [32]byte
	NewConfigDigest [32]byte
	VerifierAddress common.Address
	Raw             types.Log
}

func (_VerifierProxy *VerifierProxyFilterer) FilterVerifierSet(opts *bind.FilterOpts) (*VerifierProxyVerifierSetIterator, error) {

	logs, sub, err := _VerifierProxy.contract.FilterLogs(opts, "VerifierSet")
	if err != nil {
		return nil, err
	}
	return &VerifierProxyVerifierSetIterator{contract: _VerifierProxy.contract, event: "VerifierSet", logs: logs, sub: sub}, nil
}

func (_VerifierProxy *VerifierProxyFilterer) WatchVerifierSet(opts *bind.WatchOpts, sink chan<- *VerifierProxyVerifierSet) (event.Subscription, error) {

	logs, sub, err := _VerifierProxy.contract.WatchLogs(opts, "VerifierSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifierProxyVerifierSet)
				if err := _VerifierProxy.contract.UnpackLog(event, "VerifierSet", log); err != nil {
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

func (_VerifierProxy *VerifierProxyFilterer) ParseVerifierSet(log types.Log) (*VerifierProxyVerifierSet, error) {
	event := new(VerifierProxyVerifierSet)
	if err := _VerifierProxy.contract.UnpackLog(event, "VerifierSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifierProxyVerifierUnsetIterator struct {
	Event *VerifierProxyVerifierUnset

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifierProxyVerifierUnsetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifierProxyVerifierUnset)
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
		it.Event = new(VerifierProxyVerifierUnset)
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

func (it *VerifierProxyVerifierUnsetIterator) Error() error {
	return it.fail
}

func (it *VerifierProxyVerifierUnsetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifierProxyVerifierUnset struct {
	ConfigDigest    [32]byte
	VerifierAddress common.Address
	Raw             types.Log
}

func (_VerifierProxy *VerifierProxyFilterer) FilterVerifierUnset(opts *bind.FilterOpts) (*VerifierProxyVerifierUnsetIterator, error) {

	logs, sub, err := _VerifierProxy.contract.FilterLogs(opts, "VerifierUnset")
	if err != nil {
		return nil, err
	}
	return &VerifierProxyVerifierUnsetIterator{contract: _VerifierProxy.contract, event: "VerifierUnset", logs: logs, sub: sub}, nil
}

func (_VerifierProxy *VerifierProxyFilterer) WatchVerifierUnset(opts *bind.WatchOpts, sink chan<- *VerifierProxyVerifierUnset) (event.Subscription, error) {

	logs, sub, err := _VerifierProxy.contract.WatchLogs(opts, "VerifierUnset")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifierProxyVerifierUnset)
				if err := _VerifierProxy.contract.UnpackLog(event, "VerifierUnset", log); err != nil {
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

func (_VerifierProxy *VerifierProxyFilterer) ParseVerifierUnset(log types.Log) (*VerifierProxyVerifierUnset, error) {
	event := new(VerifierProxyVerifierUnset)
	if err := _VerifierProxy.contract.UnpackLog(event, "VerifierUnset", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_VerifierProxy *VerifierProxy) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VerifierProxy.abi.Events["AccessControllerSet"].ID:
		return _VerifierProxy.ParseAccessControllerSet(log)
	case _VerifierProxy.abi.Events["FeeManagerSet"].ID:
		return _VerifierProxy.ParseFeeManagerSet(log)
	case _VerifierProxy.abi.Events["OwnershipTransferRequested"].ID:
		return _VerifierProxy.ParseOwnershipTransferRequested(log)
	case _VerifierProxy.abi.Events["OwnershipTransferred"].ID:
		return _VerifierProxy.ParseOwnershipTransferred(log)
	case _VerifierProxy.abi.Events["VerifierInitialized"].ID:
		return _VerifierProxy.ParseVerifierInitialized(log)
	case _VerifierProxy.abi.Events["VerifierSet"].ID:
		return _VerifierProxy.ParseVerifierSet(log)
	case _VerifierProxy.abi.Events["VerifierUnset"].ID:
		return _VerifierProxy.ParseVerifierUnset(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VerifierProxyAccessControllerSet) Topic() common.Hash {
	return common.HexToHash("0x953e92b1a6442e9c3242531154a3f6f6eb00b4e9c719ba8118fa6235e4ce89b6")
}

func (VerifierProxyFeeManagerSet) Topic() common.Hash {
	return common.HexToHash("0x04628abcaa6b1674651352125cb94b65b289145bc2bc4d67720bb7d966372f03")
}

func (VerifierProxyOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VerifierProxyOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (VerifierProxyVerifierInitialized) Topic() common.Hash {
	return common.HexToHash("0x1f2cd7c97f4d801b5efe26cc409617c1fd6c5ef786e79aacb90af40923e4e8e9")
}

func (VerifierProxyVerifierSet) Topic() common.Hash {
	return common.HexToHash("0xbeb513e532542a562ac35699e7cd9ae7d198dcd3eee15bada6c857d28ceaddcf")
}

func (VerifierProxyVerifierUnset) Topic() common.Hash {
	return common.HexToHash("0x11dc15c4b8ac2b183166cc8427e5385a5ece8308217a4217338c6a7614845c4c")
}

func (_VerifierProxy *VerifierProxy) Address() common.Address {
	return _VerifierProxy.address
}

type VerifierProxyInterface interface {
	GetVerifier(opts *bind.CallOpts, configDigest [32]byte) (common.Address, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SAccessController(opts *bind.CallOpts) (common.Address, error)

	SFeeManager(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	InitializeVerifier(opts *bind.TransactOpts, verifierAddress common.Address) (*types.Transaction, error)

	SetAccessController(opts *bind.TransactOpts, accessController common.Address) (*types.Transaction, error)

	SetFeeManager(opts *bind.TransactOpts, feeManager common.Address) (*types.Transaction, error)

	SetVerifier(opts *bind.TransactOpts, currentConfigDigest [32]byte, newConfigDigest [32]byte, addressesAndWeights []CommonAddressAndWeight) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	UnsetVerifier(opts *bind.TransactOpts, configDigest [32]byte) (*types.Transaction, error)

	Verify(opts *bind.TransactOpts, payload []byte, parameterPayload []byte) (*types.Transaction, error)

	VerifyBulk(opts *bind.TransactOpts, payloads [][]byte, parameterPayload []byte) (*types.Transaction, error)

	FilterAccessControllerSet(opts *bind.FilterOpts) (*VerifierProxyAccessControllerSetIterator, error)

	WatchAccessControllerSet(opts *bind.WatchOpts, sink chan<- *VerifierProxyAccessControllerSet) (event.Subscription, error)

	ParseAccessControllerSet(log types.Log) (*VerifierProxyAccessControllerSet, error)

	FilterFeeManagerSet(opts *bind.FilterOpts) (*VerifierProxyFeeManagerSetIterator, error)

	WatchFeeManagerSet(opts *bind.WatchOpts, sink chan<- *VerifierProxyFeeManagerSet) (event.Subscription, error)

	ParseFeeManagerSet(log types.Log) (*VerifierProxyFeeManagerSet, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VerifierProxyOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VerifierProxyOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VerifierProxyOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VerifierProxyOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VerifierProxyOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VerifierProxyOwnershipTransferred, error)

	FilterVerifierInitialized(opts *bind.FilterOpts) (*VerifierProxyVerifierInitializedIterator, error)

	WatchVerifierInitialized(opts *bind.WatchOpts, sink chan<- *VerifierProxyVerifierInitialized) (event.Subscription, error)

	ParseVerifierInitialized(log types.Log) (*VerifierProxyVerifierInitialized, error)

	FilterVerifierSet(opts *bind.FilterOpts) (*VerifierProxyVerifierSetIterator, error)

	WatchVerifierSet(opts *bind.WatchOpts, sink chan<- *VerifierProxyVerifierSet) (event.Subscription, error)

	ParseVerifierSet(log types.Log) (*VerifierProxyVerifierSet, error)

	FilterVerifierUnset(opts *bind.FilterOpts) (*VerifierProxyVerifierUnsetIterator, error)

	WatchVerifierUnset(opts *bind.WatchOpts, sink chan<- *VerifierProxyVerifierUnset) (event.Subscription, error)

	ParseVerifierUnset(log types.Log) (*VerifierProxyVerifierUnset, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
