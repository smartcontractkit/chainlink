// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package destination_verifier

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

var DestinationVerifierMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"verifierProxy\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessForbidden\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BadVerification\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes24\",\"name\":\"DONConfigID\",\"type\":\"bytes24\"}],\"name\":\"DONConfigAlreadyExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DONConfigDoesNotExist\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxSigners\",\"type\":\"uint256\"}],\"name\":\"ExcessSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FaultToleranceMustBePositive\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FeeManagerInvalid\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minSigners\",\"type\":\"uint256\"}],\"name\":\"InsufficientSigners\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"rsLength\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ssLength\",\"type\":\"uint256\"}],\"name\":\"MismatchedSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NonUniqueSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldAccessController\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAccessController\",\"type\":\"address\"}],\"name\":\"AccessControllerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes24\",\"name\":\"DONConfigID\",\"type\":\"bytes24\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isActive\",\"type\":\"bool\"}],\"name\":\"ConfigActivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes24\",\"name\":\"DONConfigID\",\"type\":\"bytes24\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"weight\",\"type\":\"uint64\"}],\"indexed\":false,\"internalType\":\"structCommon.AddressAndWeight[]\",\"name\":\"recipientAddressesAndWeights\",\"type\":\"tuple[]\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldFeeManager\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newFeeManager\",\"type\":\"address\"}],\"name\":\"FeeManagerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"}],\"name\":\"ReportVerified\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAccessController\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFeeManager\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_verifierProxy\",\"outputs\":[{\"internalType\":\"contractIDestinationVerifierProxy\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"accessController\",\"type\":\"address\"}],\"name\":\"setAccessController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"weight\",\"type\":\"uint64\"}],\"internalType\":\"structCommon.AddressAndWeight[]\",\"name\":\"recipientAddressesAndWeights\",\"type\":\"tuple[]\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"DONConfigIndex\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isActive\",\"type\":\"bool\"}],\"name\":\"setConfigActive\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"feeManager\",\"type\":\"address\"}],\"name\":\"setFeeManager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isVerifier\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"signedReport\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"parameterPayload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"verify\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"signedReports\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"parameterPayload\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"verifyBulk\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"\",\"type\":\"bytes[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b5060405162002d9438038062002d948339810160408190526200003491620001a6565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000fb565b5050506001600160a01b038116620000e95760405163d92e233d60e01b815260040160405180910390fd5b6001600160a01b0316608052620001d8565b336001600160a01b03821603620001555760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600060208284031215620001b957600080fd5b81516001600160a01b0381168114620001d157600080fd5b9392505050565b608051612b9262000202600039600081816102bb0152818161036c0152610ea00152612b926000f3fe6080604052600436106100dd5760003560e01c806379ba50971161007f578063d7c72e4e11610059578063d7c72e4e146102dd578063f08391d8146102fd578063f2d638261461031d578063f2fde38b1461034857600080fd5b806379ba5097146102695780638da5cb5b1461027e578063b97455c7146102a957600080fd5b8063294d2bb1116100bb578063294d2bb1146101f4578063453ec61b14610207578063472d35b9146102295780635ad72fae1461024957600080fd5b806301ffc9a7146100e257806316d6b5f614610159578063181f5a77146101a5575b600080fd5b3480156100ee57600080fd5b506101446100fd366004611cd1565b7fffffffff00000000000000000000000000000000000000000000000000000000167f294d2bb1000000000000000000000000000000000000000000000000000000001490565b60405190151581526020015b60405180910390f35b34801561016557600080fd5b5060055473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610150565b3480156101b157600080fd5b5060408051808201909152601981527f44657374696e6174696f6e566572696669657220312e302e300000000000000060208201525b6040516101509190611d7e565b6101e7610202366004611e03565b610368565b34801561021357600080fd5b50610227610222366004612030565b6105cd565b005b34801561023557600080fd5b50610227610244366004612103565b610bfa565b34801561025557600080fd5b5061022761026436600461212c565b610c89565b34801561027557600080fd5b50610227610d9f565b34801561028a57600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff16610180565b3480156102b557600080fd5b506101807f000000000000000000000000000000000000000000000000000000000000000081565b6102f06102eb36600461215c565b610e9c565b60405161015091906121df565b34801561030957600080fd5b50610227610318366004612103565b6111a8565b34801561032957600080fd5b5060045473ffffffffffffffffffffffffffffffffffffffff16610180565b34801561035457600080fd5b50610227610363366004612103565b61122f565b60607f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1633146103d9576040517fef67f5d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600554829073ffffffffffffffffffffffffffffffffffffffff16801580159061049857506040517f6b14daf800000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821690636b14daf89061045590859060009036906004016122a8565b602060405180830381865afa158015610472573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061049691906122e1565b155b156104cf576040517fef67f5d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000806104dd8a8a88611243565b600454919350915073ffffffffffffffffffffffffffffffffffffffff16156105c057600480546040517f86968cfd00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff909116916386968cfd9161055f9185918f918f918f918f918f91016122fe565b600060405180830381600087803b15801561057957600080fd5b505af192505050801561058a575060015b6105c0576040517f4df18f0700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5098975050505050505050565b82518260ff168060000361060d576040517f0743bae600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601f821115610657576040517f61750f4000000000000000000000000000000000000000000000000000000000815260048101839052601f60248201526044015b60405180910390fd5b610662816003612384565b82116106ba5781610674826003612384565b61067f9060016123a1565b6040517f9dd9e6d80000000000000000000000000000000000000000000000000000000081526004810192909252602482015260440161064e565b6106c2611748565b6106cb856117cb565b15610702576040517ff67bc7c400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61071b8560006001885161071691906123b4565b61187a565b600085856040516020016107309291906123c7565b60405160208183030381529060405280519060200120905060005b86518110156108d757600073ffffffffffffffffffffffffffffffffffffffff1687828151811061077e5761077e61243f565b602002602001015173ffffffffffffffffffffffffffffffffffffffff16036107d3576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600260008984815181106107eb576107eb61243f565b60200260200101518560405160200161085792919060609290921b7fffffffffffffffffffffffffffffffffffffffff0000000000000000000000001682527fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000166014820152602c0190565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001815291815281516020928301208352908201929092520160002080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169115159190911790556108d08161246e565b905061074b565b506003541580159061095c5750600380547fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000083169190610919906001906123b4565b815481106109295761092961243f565b60009182526020909120015460401b7fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000016145b156109b7576040517f9a92755c0000000000000000000000000000000000000000000000000000000081527fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000008216600482015260240161064e565b835115610a4857600480546040517ff65df96200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091169163f65df96291610a1591859189910161250c565b600060405180830381600087803b158015610a2f57600080fd5b505af1158015610a43573d6000803e3d6000fd5b505050505b604080516080810182527fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000831680825260ff80891660208401908152600184860181815263ffffffff4281166060880190815260038054948501815560005296517fc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85b90930180549451925197519091167a010000000000000000000000000000000000000000000000000000027fffff00000000ffffffffffffffffffffffffffffffffffffffffffffffffffff97151579010000000000000000000000000000000000000000000000000002979097167fffff0000000000ffffffffffffffffffffffffffffffffffffffffffffffffff929095167801000000000000000000000000000000000000000000000000027fffffffffffffff0000000000000000000000000000000000000000000000000090941692881c929092179290921791909116919091179290921790915590517f2d763a674a99583454a287d792819ffb9ff7e791c23e7745a082701136ce336c90610bea9089908990899061254f565b60405180910390a2505050505050565b610c02611748565b6004805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff000000000000000000000000000000000000000083168117909355604080519190921680825260208201939093527f04628abcaa6b1674651352125cb94b65b289145bc2bc4d67720bb7d966372f0391015b60405180910390a15050565b610c91611748565b6003548210610ccc576040517f59d7257e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600060038381548110610ce157610ce161243f565b600091825260209182902001805484151579010000000000000000000000000000000000000000000000000081027fffffffffffff00ffffffffffffffffffffffffffffffffffffffffffffffffff909216919091178083556040805191811b7fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000168252938101919091529092507f90186a1e77b498ec417ea88bd026cae00d7043c357cc45221777623bda582dd4910160405180910390a1505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610e20576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015260640161064e565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60607f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff163314610f0d576040517fef67f5d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600554829073ffffffffffffffffffffffffffffffffffffffff168015801590610fcc57506040517f6b14daf800000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821690636b14daf890610f8990859060009036906004016122a8565b602060405180830381865afa158015610fa6573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610fca91906122e1565b155b15611003576040517fef67f5d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008767ffffffffffffffff81111561101e5761101e611e84565b60405190808252806020026020018201604052801561105157816020015b606081526020019060019003908161103c5790505b50905060008867ffffffffffffffff81111561106f5761106f611e84565b604051908082528060200260200182016040528015611098578160200160208202803683370190505b50905060005b8981101561112a576000806110d68d8d858181106110be576110be61243f565b90506020028101906110d091906125c6565b8b611243565b91509150818584815181106110ed576110ed61243f565b60200260200101819052508084848151811061110b5761110b61243f565b6020026020010181815250505050806111239061246e565b905061109e565b5060045473ffffffffffffffffffffffffffffffffffffffff16156105c057600480546040517f3690750900000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091169163369075099161055f9185918f918f918f918f918f910161262b565b6111b0611748565b6005805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff000000000000000000000000000000000000000083168117909355604080519190921680825260208201939093527f953e92b1a6442e9c3242531154a3f6f6eb00b4e9c719ba8118fa6235e4ce89b69101610c7d565b611237611748565b61124081611a6a565b50565b606060008080808080611258898b018b612852565b9450945094509450945081518351146112aa57825182516040517ff0d314080000000000000000000000000000000000000000000000000000000081526004810192909252602482015260440161064e565b82516000036112e5576040517fc7af40f000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600084805190602001208660405160200161130192919061292d565b6040516020818303038152906040528051906020012090506000845167ffffffffffffffff81111561133557611335611e84565b60405190808252806020026020018201604052801561135e578160200160208202803683370190505b50905060005b855181101561146c576001838583602081106113825761138261243f565b61138f91901a601b612969565b8884815181106113a1576113a161243f565b60200260200101518885815181106113bb576113bb61243f565b6020026020010151604051600081526020016040526040516113f9949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa15801561141b573d6000803e3d6000fd5b505050602060405103518282815181106114375761143761243f565b73ffffffffffffffffffffffffffffffffffffffff909216602092830291909101909101526114658161246e565b9050611364565b50611476816117cb565b156114ad576040517f4df18f0700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006114b887611b5f565b905060006114c582611b85565b80519091507fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000016611522576040517f4df18f0700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b806040015161155d576040517f4df18f0700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b806020015160ff1683511161159e576040517f4df18f0700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000805b84518110156116bc578481815181106115bd576115bd61243f565b6020026020010151836000015160405160200161162d92919060609290921b7fffffffffffffffffffffffffffffffffffffffff0000000000000000000000001682527fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000166014820152602c0190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291815281516020928301206000818152600290935291205490925060ff166116ac576040517f4df18f0700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6116b58161246e565b90506115a2565b506116c689612982565b60405173ffffffffffffffffffffffffffffffffffffffff8f1681527f58ca9502e98a536e06e72d680fcc251e5d10b72291a281665a2c2dc0ac30fcc59060200160405180910390a25051969d7fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000009097169c50959a5050505050505050505050565b60005473ffffffffffffffffffffffffffffffffffffffff1633146117c9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161064e565b565b6000805b82518110156118715760006117e58260016123a1565b90505b8351811015611868578381815181106118035761180361243f565b602002602001015173ffffffffffffffffffffffffffffffffffffffff168483815181106118335761183361243f565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1603611860575060019392505050565b6001016117e8565b506001016117cf565b50600092915050565b818180820361188a575050505050565b600085600261189987876129c7565b6118a391906129e7565b6118ad9087612a76565b815181106118bd576118bd61243f565b602002602001015190505b818313611a3c575b8073ffffffffffffffffffffffffffffffffffffffff168684815181106118f9576118f961243f565b602002602001015173ffffffffffffffffffffffffffffffffffffffff16101561192f578261192781612a9e565b9350506118d0565b8582815181106119415761194161243f565b602002602001015173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16101561198e578161198681612acf565b92505061192f565b818313611a37578582815181106119a7576119a761243f565b60200260200101518684815181106119c1576119c161243f565b60200260200101518785815181106119db576119db61243f565b602002602001018885815181106119f4576119f461243f565b73ffffffffffffffffffffffffffffffffffffffff93841660209182029290920101529116905282611a2581612a9e565b9350508180611a3390612acf565b9250505b6118c8565b81851215611a4f57611a4f86868461187a565b83831215611a6257611a6286848661187a565b505050505050565b3373ffffffffffffffffffffffffffffffffffffffff821603611ae9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161064e565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60008082806020019051810190611b769190612b3a565b63ffffffff1695945050505050565b604080516080810182526000808252602082018190529181018290526060810191909152604080516080810182526000808252602082018190529181018290526060810191909152600354600090611bdf906001906123b4565b90505b60038181548110611bf557611bf561243f565b600091825260209182902060408051608081018252929091015480821b7fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001683527801000000000000000000000000000000000000000000000000810460ff9081169484019490945279010000000000000000000000000000000000000000000000000081049093161515908201527a01000000000000000000000000000000000000000000000000000090910463ffffffff1660608201819052909250841015611cca57611cc381612b76565b9050611be2565b5092915050565b600060208284031215611ce357600080fd5b81357fffffffff0000000000000000000000000000000000000000000000000000000081168114611d1357600080fd5b9392505050565b6000815180845260005b81811015611d4057602081850181015186830182015201611d24565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b602081526000611d136020830184611d1a565b60008083601f840112611da357600080fd5b50813567ffffffffffffffff811115611dbb57600080fd5b602083019150836020828501011115611dd357600080fd5b9250929050565b803573ffffffffffffffffffffffffffffffffffffffff81168114611dfe57600080fd5b919050565b600080600080600060608688031215611e1b57600080fd5b853567ffffffffffffffff80821115611e3357600080fd5b611e3f89838a01611d91565b90975095506020880135915080821115611e5857600080fd5b50611e6588828901611d91565b9094509250611e78905060408701611dda565b90509295509295909350565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040805190810167ffffffffffffffff81118282101715611ed657611ed6611e84565b60405290565b6040516060810167ffffffffffffffff81118282101715611ed657611ed6611e84565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715611f4657611f46611e84565b604052919050565b600067ffffffffffffffff821115611f6857611f68611e84565b5060051b60200190565b803560ff81168114611dfe57600080fd5b600082601f830112611f9457600080fd5b81356020611fa9611fa483611f4e565b611eff565b82815260069290921b84018101918181019086841115611fc857600080fd5b8286015b848110156120255760408189031215611fe55760008081fd5b611fed611eb3565b611ff682611dda565b81528482013567ffffffffffffffff811681146120135760008081fd5b81860152835291830191604001611fcc565b509695505050505050565b60008060006060848603121561204557600080fd5b833567ffffffffffffffff8082111561205d57600080fd5b818601915086601f83011261207157600080fd5b81356020612081611fa483611f4e565b82815260059290921b8401810191818101908a8411156120a057600080fd5b948201945b838610156120c5576120b686611dda565b825294820194908201906120a5565b97506120d49050888201611f72565b9550505060408601359150808211156120ec57600080fd5b506120f986828701611f83565b9150509250925092565b60006020828403121561211557600080fd5b611d1382611dda565b801515811461124057600080fd5b6000806040838503121561213f57600080fd5b8235915060208301356121518161211e565b809150509250929050565b60008060008060006060868803121561217457600080fd5b853567ffffffffffffffff8082111561218c57600080fd5b818801915088601f8301126121a057600080fd5b8135818111156121af57600080fd5b8960208260051b85010111156121c457600080fd5b602092830197509550908701359080821115611e5857600080fd5b6000602080830181845280855180835260408601915060408160051b870101925083870160005b82811015612252577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0888603018452612240858351611d1a565b94509285019290850190600101612206565b5092979650505050505050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b73ffffffffffffffffffffffffffffffffffffffff841681526040602082015260006122d860408301848661225f565b95945050505050565b6000602082840312156122f357600080fd5b8151611d138161211e565b86815260806020820152600061231860808301878961225f565b828103604084015261232b81868861225f565b91505073ffffffffffffffffffffffffffffffffffffffff83166060830152979650505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b808202811582820484141761239b5761239b612355565b92915050565b8082018082111561239b5761239b612355565b8181038181111561239b5761239b612355565b825160009082906020808701845b8381101561240757815173ffffffffffffffffffffffffffffffffffffffff16855293820193908201906001016123d5565b5050505060f89390931b7fff000000000000000000000000000000000000000000000000000000000000001683525050600101919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361249f5761249f612355565b5060010190565b600081518084526020808501945080840160005b83811015612501578151805173ffffffffffffffffffffffffffffffffffffffff16885283015167ffffffffffffffff1683880152604090960195908201906001016124ba565b509495945050505050565b7fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000008316815260406020820152600061254760408301846124a6565b949350505050565b606080825284519082018190526000906020906080840190828801845b8281101561259e57815173ffffffffffffffffffffffffffffffffffffffff168452928401929084019060010161256c565b50505060ff86168285015283810360408501526125bb81866124a6565b979650505050505050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18436030181126125fb57600080fd5b83018035915067ffffffffffffffff82111561261657600080fd5b602001915036819003821315611dd357600080fd5b6080808252875190820181905260009060209060a0840190828b01845b8281101561266457815184529284019290840190600101612648565b50505083810382850152878152818101600589901b820183018a60005b8b81101561272c577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe085840301845281357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18e36030181126126e257600080fd5b8d01868101903567ffffffffffffffff8111156126fe57600080fd5b80360382131561270d57600080fd5b61271885828461225f565b958801959450505090850190600101612681565b5050858103604087015261274181898b61225f565b9450505050506125bb606083018473ffffffffffffffffffffffffffffffffffffffff169052565b600082601f83011261277a57600080fd5b813567ffffffffffffffff81111561279457612794611e84565b6127c560207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601611eff565b8181528460208386010111156127da57600080fd5b816020850160208301376000918101602001919091529392505050565b600082601f83011261280857600080fd5b81356020612818611fa483611f4e565b82815260059290921b8401810191818101908684111561283757600080fd5b8286015b84811015612025578035835291830191830161283b565b600080600080600060e0868803121561286a57600080fd5b86601f87011261287957600080fd5b612881611edc565b80606088018981111561289357600080fd5b885b818110156128ad578035845260209384019301612895565b5090965035905067ffffffffffffffff808211156128ca57600080fd5b6128d689838a01612769565b955060808801359150808211156128ec57600080fd5b6128f889838a016127f7565b945060a088013591508082111561290e57600080fd5b5061291b888289016127f7565b9598949750929560c001359392505050565b828152600060208083018460005b60038110156129585781518352918301919083019060010161293b565b505050506080820190509392505050565b60ff818116838216019081111561239b5761239b612355565b805160208083015191908110156129c1577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8160200360031b1b821691505b50919050565b8181036000831280158383131683831282161715611cca57611cca612355565b600082612a1d577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff83147f800000000000000000000000000000000000000000000000000000000000000083141615612a7157612a71612355565b500590565b8082018281126000831280158216821582161715612a9657612a96612355565b505092915050565b60007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361249f5761249f612355565b60007f80000000000000000000000000000000000000000000000000000000000000008203612b0057612b00612355565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190565b805163ffffffff81168114611dfe57600080fd5b600080600060608486031215612b4f57600080fd5b83519250612b5f60208501612b26565b9150612b6d60408501612b26565b90509250925092565b600081612b0057612b0061235556fea164736f6c6343000813000a",
}

var DestinationVerifierABI = DestinationVerifierMetaData.ABI

var DestinationVerifierBin = DestinationVerifierMetaData.Bin

func DeployDestinationVerifier(auth *bind.TransactOpts, backend bind.ContractBackend, verifierProxy common.Address) (common.Address, *types.Transaction, *DestinationVerifier, error) {
	parsed, err := DestinationVerifierMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DestinationVerifierBin), backend, verifierProxy)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DestinationVerifier{address: address, abi: *parsed, DestinationVerifierCaller: DestinationVerifierCaller{contract: contract}, DestinationVerifierTransactor: DestinationVerifierTransactor{contract: contract}, DestinationVerifierFilterer: DestinationVerifierFilterer{contract: contract}}, nil
}

type DestinationVerifier struct {
	address common.Address
	abi     abi.ABI
	DestinationVerifierCaller
	DestinationVerifierTransactor
	DestinationVerifierFilterer
}

type DestinationVerifierCaller struct {
	contract *bind.BoundContract
}

type DestinationVerifierTransactor struct {
	contract *bind.BoundContract
}

type DestinationVerifierFilterer struct {
	contract *bind.BoundContract
}

type DestinationVerifierSession struct {
	Contract     *DestinationVerifier
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type DestinationVerifierCallerSession struct {
	Contract *DestinationVerifierCaller
	CallOpts bind.CallOpts
}

type DestinationVerifierTransactorSession struct {
	Contract     *DestinationVerifierTransactor
	TransactOpts bind.TransactOpts
}

type DestinationVerifierRaw struct {
	Contract *DestinationVerifier
}

type DestinationVerifierCallerRaw struct {
	Contract *DestinationVerifierCaller
}

type DestinationVerifierTransactorRaw struct {
	Contract *DestinationVerifierTransactor
}

func NewDestinationVerifier(address common.Address, backend bind.ContractBackend) (*DestinationVerifier, error) {
	abi, err := abi.JSON(strings.NewReader(DestinationVerifierABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindDestinationVerifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DestinationVerifier{address: address, abi: abi, DestinationVerifierCaller: DestinationVerifierCaller{contract: contract}, DestinationVerifierTransactor: DestinationVerifierTransactor{contract: contract}, DestinationVerifierFilterer: DestinationVerifierFilterer{contract: contract}}, nil
}

func NewDestinationVerifierCaller(address common.Address, caller bind.ContractCaller) (*DestinationVerifierCaller, error) {
	contract, err := bindDestinationVerifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DestinationVerifierCaller{contract: contract}, nil
}

func NewDestinationVerifierTransactor(address common.Address, transactor bind.ContractTransactor) (*DestinationVerifierTransactor, error) {
	contract, err := bindDestinationVerifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DestinationVerifierTransactor{contract: contract}, nil
}

func NewDestinationVerifierFilterer(address common.Address, filterer bind.ContractFilterer) (*DestinationVerifierFilterer, error) {
	contract, err := bindDestinationVerifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DestinationVerifierFilterer{contract: contract}, nil
}

func bindDestinationVerifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DestinationVerifierMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_DestinationVerifier *DestinationVerifierRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DestinationVerifier.Contract.DestinationVerifierCaller.contract.Call(opts, result, method, params...)
}

func (_DestinationVerifier *DestinationVerifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DestinationVerifier.Contract.DestinationVerifierTransactor.contract.Transfer(opts)
}

func (_DestinationVerifier *DestinationVerifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DestinationVerifier.Contract.DestinationVerifierTransactor.contract.Transact(opts, method, params...)
}

func (_DestinationVerifier *DestinationVerifierCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DestinationVerifier.Contract.contract.Call(opts, result, method, params...)
}

func (_DestinationVerifier *DestinationVerifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DestinationVerifier.Contract.contract.Transfer(opts)
}

func (_DestinationVerifier *DestinationVerifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DestinationVerifier.Contract.contract.Transact(opts, method, params...)
}

func (_DestinationVerifier *DestinationVerifierCaller) GetAccessController(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DestinationVerifier.contract.Call(opts, &out, "getAccessController")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_DestinationVerifier *DestinationVerifierSession) GetAccessController() (common.Address, error) {
	return _DestinationVerifier.Contract.GetAccessController(&_DestinationVerifier.CallOpts)
}

func (_DestinationVerifier *DestinationVerifierCallerSession) GetAccessController() (common.Address, error) {
	return _DestinationVerifier.Contract.GetAccessController(&_DestinationVerifier.CallOpts)
}

func (_DestinationVerifier *DestinationVerifierCaller) GetFeeManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DestinationVerifier.contract.Call(opts, &out, "getFeeManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_DestinationVerifier *DestinationVerifierSession) GetFeeManager() (common.Address, error) {
	return _DestinationVerifier.Contract.GetFeeManager(&_DestinationVerifier.CallOpts)
}

func (_DestinationVerifier *DestinationVerifierCallerSession) GetFeeManager() (common.Address, error) {
	return _DestinationVerifier.Contract.GetFeeManager(&_DestinationVerifier.CallOpts)
}

func (_DestinationVerifier *DestinationVerifierCaller) IVerifierProxy(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DestinationVerifier.contract.Call(opts, &out, "i_verifierProxy")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_DestinationVerifier *DestinationVerifierSession) IVerifierProxy() (common.Address, error) {
	return _DestinationVerifier.Contract.IVerifierProxy(&_DestinationVerifier.CallOpts)
}

func (_DestinationVerifier *DestinationVerifierCallerSession) IVerifierProxy() (common.Address, error) {
	return _DestinationVerifier.Contract.IVerifierProxy(&_DestinationVerifier.CallOpts)
}

func (_DestinationVerifier *DestinationVerifierCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DestinationVerifier.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_DestinationVerifier *DestinationVerifierSession) Owner() (common.Address, error) {
	return _DestinationVerifier.Contract.Owner(&_DestinationVerifier.CallOpts)
}

func (_DestinationVerifier *DestinationVerifierCallerSession) Owner() (common.Address, error) {
	return _DestinationVerifier.Contract.Owner(&_DestinationVerifier.CallOpts)
}

func (_DestinationVerifier *DestinationVerifierCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _DestinationVerifier.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_DestinationVerifier *DestinationVerifierSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _DestinationVerifier.Contract.SupportsInterface(&_DestinationVerifier.CallOpts, interfaceId)
}

func (_DestinationVerifier *DestinationVerifierCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _DestinationVerifier.Contract.SupportsInterface(&_DestinationVerifier.CallOpts, interfaceId)
}

func (_DestinationVerifier *DestinationVerifierCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _DestinationVerifier.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_DestinationVerifier *DestinationVerifierSession) TypeAndVersion() (string, error) {
	return _DestinationVerifier.Contract.TypeAndVersion(&_DestinationVerifier.CallOpts)
}

func (_DestinationVerifier *DestinationVerifierCallerSession) TypeAndVersion() (string, error) {
	return _DestinationVerifier.Contract.TypeAndVersion(&_DestinationVerifier.CallOpts)
}

func (_DestinationVerifier *DestinationVerifierTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DestinationVerifier.contract.Transact(opts, "acceptOwnership")
}

func (_DestinationVerifier *DestinationVerifierSession) AcceptOwnership() (*types.Transaction, error) {
	return _DestinationVerifier.Contract.AcceptOwnership(&_DestinationVerifier.TransactOpts)
}

func (_DestinationVerifier *DestinationVerifierTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _DestinationVerifier.Contract.AcceptOwnership(&_DestinationVerifier.TransactOpts)
}

func (_DestinationVerifier *DestinationVerifierTransactor) SetAccessController(opts *bind.TransactOpts, accessController common.Address) (*types.Transaction, error) {
	return _DestinationVerifier.contract.Transact(opts, "setAccessController", accessController)
}

func (_DestinationVerifier *DestinationVerifierSession) SetAccessController(accessController common.Address) (*types.Transaction, error) {
	return _DestinationVerifier.Contract.SetAccessController(&_DestinationVerifier.TransactOpts, accessController)
}

func (_DestinationVerifier *DestinationVerifierTransactorSession) SetAccessController(accessController common.Address) (*types.Transaction, error) {
	return _DestinationVerifier.Contract.SetAccessController(&_DestinationVerifier.TransactOpts, accessController)
}

func (_DestinationVerifier *DestinationVerifierTransactor) SetConfig(opts *bind.TransactOpts, signers []common.Address, f uint8, recipientAddressesAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _DestinationVerifier.contract.Transact(opts, "setConfig", signers, f, recipientAddressesAndWeights)
}

func (_DestinationVerifier *DestinationVerifierSession) SetConfig(signers []common.Address, f uint8, recipientAddressesAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _DestinationVerifier.Contract.SetConfig(&_DestinationVerifier.TransactOpts, signers, f, recipientAddressesAndWeights)
}

func (_DestinationVerifier *DestinationVerifierTransactorSession) SetConfig(signers []common.Address, f uint8, recipientAddressesAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _DestinationVerifier.Contract.SetConfig(&_DestinationVerifier.TransactOpts, signers, f, recipientAddressesAndWeights)
}

func (_DestinationVerifier *DestinationVerifierTransactor) SetConfigActive(opts *bind.TransactOpts, DONConfigIndex *big.Int, isActive bool) (*types.Transaction, error) {
	return _DestinationVerifier.contract.Transact(opts, "setConfigActive", DONConfigIndex, isActive)
}

func (_DestinationVerifier *DestinationVerifierSession) SetConfigActive(DONConfigIndex *big.Int, isActive bool) (*types.Transaction, error) {
	return _DestinationVerifier.Contract.SetConfigActive(&_DestinationVerifier.TransactOpts, DONConfigIndex, isActive)
}

func (_DestinationVerifier *DestinationVerifierTransactorSession) SetConfigActive(DONConfigIndex *big.Int, isActive bool) (*types.Transaction, error) {
	return _DestinationVerifier.Contract.SetConfigActive(&_DestinationVerifier.TransactOpts, DONConfigIndex, isActive)
}

func (_DestinationVerifier *DestinationVerifierTransactor) SetFeeManager(opts *bind.TransactOpts, feeManager common.Address) (*types.Transaction, error) {
	return _DestinationVerifier.contract.Transact(opts, "setFeeManager", feeManager)
}

func (_DestinationVerifier *DestinationVerifierSession) SetFeeManager(feeManager common.Address) (*types.Transaction, error) {
	return _DestinationVerifier.Contract.SetFeeManager(&_DestinationVerifier.TransactOpts, feeManager)
}

func (_DestinationVerifier *DestinationVerifierTransactorSession) SetFeeManager(feeManager common.Address) (*types.Transaction, error) {
	return _DestinationVerifier.Contract.SetFeeManager(&_DestinationVerifier.TransactOpts, feeManager)
}

func (_DestinationVerifier *DestinationVerifierTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _DestinationVerifier.contract.Transact(opts, "transferOwnership", to)
}

func (_DestinationVerifier *DestinationVerifierSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _DestinationVerifier.Contract.TransferOwnership(&_DestinationVerifier.TransactOpts, to)
}

func (_DestinationVerifier *DestinationVerifierTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _DestinationVerifier.Contract.TransferOwnership(&_DestinationVerifier.TransactOpts, to)
}

func (_DestinationVerifier *DestinationVerifierTransactor) Verify(opts *bind.TransactOpts, signedReport []byte, parameterPayload []byte, sender common.Address) (*types.Transaction, error) {
	return _DestinationVerifier.contract.Transact(opts, "verify", signedReport, parameterPayload, sender)
}

func (_DestinationVerifier *DestinationVerifierSession) Verify(signedReport []byte, parameterPayload []byte, sender common.Address) (*types.Transaction, error) {
	return _DestinationVerifier.Contract.Verify(&_DestinationVerifier.TransactOpts, signedReport, parameterPayload, sender)
}

func (_DestinationVerifier *DestinationVerifierTransactorSession) Verify(signedReport []byte, parameterPayload []byte, sender common.Address) (*types.Transaction, error) {
	return _DestinationVerifier.Contract.Verify(&_DestinationVerifier.TransactOpts, signedReport, parameterPayload, sender)
}

func (_DestinationVerifier *DestinationVerifierTransactor) VerifyBulk(opts *bind.TransactOpts, signedReports [][]byte, parameterPayload []byte, sender common.Address) (*types.Transaction, error) {
	return _DestinationVerifier.contract.Transact(opts, "verifyBulk", signedReports, parameterPayload, sender)
}

func (_DestinationVerifier *DestinationVerifierSession) VerifyBulk(signedReports [][]byte, parameterPayload []byte, sender common.Address) (*types.Transaction, error) {
	return _DestinationVerifier.Contract.VerifyBulk(&_DestinationVerifier.TransactOpts, signedReports, parameterPayload, sender)
}

func (_DestinationVerifier *DestinationVerifierTransactorSession) VerifyBulk(signedReports [][]byte, parameterPayload []byte, sender common.Address) (*types.Transaction, error) {
	return _DestinationVerifier.Contract.VerifyBulk(&_DestinationVerifier.TransactOpts, signedReports, parameterPayload, sender)
}

type DestinationVerifierAccessControllerSetIterator struct {
	Event *DestinationVerifierAccessControllerSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DestinationVerifierAccessControllerSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DestinationVerifierAccessControllerSet)
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
		it.Event = new(DestinationVerifierAccessControllerSet)
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

func (it *DestinationVerifierAccessControllerSetIterator) Error() error {
	return it.fail
}

func (it *DestinationVerifierAccessControllerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DestinationVerifierAccessControllerSet struct {
	OldAccessController common.Address
	NewAccessController common.Address
	Raw                 types.Log
}

func (_DestinationVerifier *DestinationVerifierFilterer) FilterAccessControllerSet(opts *bind.FilterOpts) (*DestinationVerifierAccessControllerSetIterator, error) {

	logs, sub, err := _DestinationVerifier.contract.FilterLogs(opts, "AccessControllerSet")
	if err != nil {
		return nil, err
	}
	return &DestinationVerifierAccessControllerSetIterator{contract: _DestinationVerifier.contract, event: "AccessControllerSet", logs: logs, sub: sub}, nil
}

func (_DestinationVerifier *DestinationVerifierFilterer) WatchAccessControllerSet(opts *bind.WatchOpts, sink chan<- *DestinationVerifierAccessControllerSet) (event.Subscription, error) {

	logs, sub, err := _DestinationVerifier.contract.WatchLogs(opts, "AccessControllerSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DestinationVerifierAccessControllerSet)
				if err := _DestinationVerifier.contract.UnpackLog(event, "AccessControllerSet", log); err != nil {
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

func (_DestinationVerifier *DestinationVerifierFilterer) ParseAccessControllerSet(log types.Log) (*DestinationVerifierAccessControllerSet, error) {
	event := new(DestinationVerifierAccessControllerSet)
	if err := _DestinationVerifier.contract.UnpackLog(event, "AccessControllerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DestinationVerifierConfigActivatedIterator struct {
	Event *DestinationVerifierConfigActivated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DestinationVerifierConfigActivatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DestinationVerifierConfigActivated)
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
		it.Event = new(DestinationVerifierConfigActivated)
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

func (it *DestinationVerifierConfigActivatedIterator) Error() error {
	return it.fail
}

func (it *DestinationVerifierConfigActivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DestinationVerifierConfigActivated struct {
	DONConfigID [24]byte
	IsActive    bool
	Raw         types.Log
}

func (_DestinationVerifier *DestinationVerifierFilterer) FilterConfigActivated(opts *bind.FilterOpts) (*DestinationVerifierConfigActivatedIterator, error) {

	logs, sub, err := _DestinationVerifier.contract.FilterLogs(opts, "ConfigActivated")
	if err != nil {
		return nil, err
	}
	return &DestinationVerifierConfigActivatedIterator{contract: _DestinationVerifier.contract, event: "ConfigActivated", logs: logs, sub: sub}, nil
}

func (_DestinationVerifier *DestinationVerifierFilterer) WatchConfigActivated(opts *bind.WatchOpts, sink chan<- *DestinationVerifierConfigActivated) (event.Subscription, error) {

	logs, sub, err := _DestinationVerifier.contract.WatchLogs(opts, "ConfigActivated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DestinationVerifierConfigActivated)
				if err := _DestinationVerifier.contract.UnpackLog(event, "ConfigActivated", log); err != nil {
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

func (_DestinationVerifier *DestinationVerifierFilterer) ParseConfigActivated(log types.Log) (*DestinationVerifierConfigActivated, error) {
	event := new(DestinationVerifierConfigActivated)
	if err := _DestinationVerifier.contract.UnpackLog(event, "ConfigActivated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DestinationVerifierConfigSetIterator struct {
	Event *DestinationVerifierConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DestinationVerifierConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DestinationVerifierConfigSet)
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
		it.Event = new(DestinationVerifierConfigSet)
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

func (it *DestinationVerifierConfigSetIterator) Error() error {
	return it.fail
}

func (it *DestinationVerifierConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DestinationVerifierConfigSet struct {
	DONConfigID                  [24]byte
	Signers                      []common.Address
	F                            uint8
	RecipientAddressesAndWeights []CommonAddressAndWeight
	Raw                          types.Log
}

func (_DestinationVerifier *DestinationVerifierFilterer) FilterConfigSet(opts *bind.FilterOpts, DONConfigID [][24]byte) (*DestinationVerifierConfigSetIterator, error) {

	var DONConfigIDRule []interface{}
	for _, DONConfigIDItem := range DONConfigID {
		DONConfigIDRule = append(DONConfigIDRule, DONConfigIDItem)
	}

	logs, sub, err := _DestinationVerifier.contract.FilterLogs(opts, "ConfigSet", DONConfigIDRule)
	if err != nil {
		return nil, err
	}
	return &DestinationVerifierConfigSetIterator{contract: _DestinationVerifier.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_DestinationVerifier *DestinationVerifierFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *DestinationVerifierConfigSet, DONConfigID [][24]byte) (event.Subscription, error) {

	var DONConfigIDRule []interface{}
	for _, DONConfigIDItem := range DONConfigID {
		DONConfigIDRule = append(DONConfigIDRule, DONConfigIDItem)
	}

	logs, sub, err := _DestinationVerifier.contract.WatchLogs(opts, "ConfigSet", DONConfigIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DestinationVerifierConfigSet)
				if err := _DestinationVerifier.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_DestinationVerifier *DestinationVerifierFilterer) ParseConfigSet(log types.Log) (*DestinationVerifierConfigSet, error) {
	event := new(DestinationVerifierConfigSet)
	if err := _DestinationVerifier.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DestinationVerifierFeeManagerSetIterator struct {
	Event *DestinationVerifierFeeManagerSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DestinationVerifierFeeManagerSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DestinationVerifierFeeManagerSet)
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
		it.Event = new(DestinationVerifierFeeManagerSet)
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

func (it *DestinationVerifierFeeManagerSetIterator) Error() error {
	return it.fail
}

func (it *DestinationVerifierFeeManagerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DestinationVerifierFeeManagerSet struct {
	OldFeeManager common.Address
	NewFeeManager common.Address
	Raw           types.Log
}

func (_DestinationVerifier *DestinationVerifierFilterer) FilterFeeManagerSet(opts *bind.FilterOpts) (*DestinationVerifierFeeManagerSetIterator, error) {

	logs, sub, err := _DestinationVerifier.contract.FilterLogs(opts, "FeeManagerSet")
	if err != nil {
		return nil, err
	}
	return &DestinationVerifierFeeManagerSetIterator{contract: _DestinationVerifier.contract, event: "FeeManagerSet", logs: logs, sub: sub}, nil
}

func (_DestinationVerifier *DestinationVerifierFilterer) WatchFeeManagerSet(opts *bind.WatchOpts, sink chan<- *DestinationVerifierFeeManagerSet) (event.Subscription, error) {

	logs, sub, err := _DestinationVerifier.contract.WatchLogs(opts, "FeeManagerSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DestinationVerifierFeeManagerSet)
				if err := _DestinationVerifier.contract.UnpackLog(event, "FeeManagerSet", log); err != nil {
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

func (_DestinationVerifier *DestinationVerifierFilterer) ParseFeeManagerSet(log types.Log) (*DestinationVerifierFeeManagerSet, error) {
	event := new(DestinationVerifierFeeManagerSet)
	if err := _DestinationVerifier.contract.UnpackLog(event, "FeeManagerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DestinationVerifierOwnershipTransferRequestedIterator struct {
	Event *DestinationVerifierOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DestinationVerifierOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DestinationVerifierOwnershipTransferRequested)
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
		it.Event = new(DestinationVerifierOwnershipTransferRequested)
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

func (it *DestinationVerifierOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *DestinationVerifierOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DestinationVerifierOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_DestinationVerifier *DestinationVerifierFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*DestinationVerifierOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _DestinationVerifier.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &DestinationVerifierOwnershipTransferRequestedIterator{contract: _DestinationVerifier.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_DestinationVerifier *DestinationVerifierFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *DestinationVerifierOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _DestinationVerifier.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DestinationVerifierOwnershipTransferRequested)
				if err := _DestinationVerifier.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_DestinationVerifier *DestinationVerifierFilterer) ParseOwnershipTransferRequested(log types.Log) (*DestinationVerifierOwnershipTransferRequested, error) {
	event := new(DestinationVerifierOwnershipTransferRequested)
	if err := _DestinationVerifier.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DestinationVerifierOwnershipTransferredIterator struct {
	Event *DestinationVerifierOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DestinationVerifierOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DestinationVerifierOwnershipTransferred)
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
		it.Event = new(DestinationVerifierOwnershipTransferred)
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

func (it *DestinationVerifierOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *DestinationVerifierOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DestinationVerifierOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_DestinationVerifier *DestinationVerifierFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*DestinationVerifierOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _DestinationVerifier.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &DestinationVerifierOwnershipTransferredIterator{contract: _DestinationVerifier.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_DestinationVerifier *DestinationVerifierFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DestinationVerifierOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _DestinationVerifier.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DestinationVerifierOwnershipTransferred)
				if err := _DestinationVerifier.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_DestinationVerifier *DestinationVerifierFilterer) ParseOwnershipTransferred(log types.Log) (*DestinationVerifierOwnershipTransferred, error) {
	event := new(DestinationVerifierOwnershipTransferred)
	if err := _DestinationVerifier.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DestinationVerifierReportVerifiedIterator struct {
	Event *DestinationVerifierReportVerified

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DestinationVerifierReportVerifiedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DestinationVerifierReportVerified)
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
		it.Event = new(DestinationVerifierReportVerified)
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

func (it *DestinationVerifierReportVerifiedIterator) Error() error {
	return it.fail
}

func (it *DestinationVerifierReportVerifiedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DestinationVerifierReportVerified struct {
	FeedId    [32]byte
	Requester common.Address
	Raw       types.Log
}

func (_DestinationVerifier *DestinationVerifierFilterer) FilterReportVerified(opts *bind.FilterOpts, feedId [][32]byte) (*DestinationVerifierReportVerifiedIterator, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _DestinationVerifier.contract.FilterLogs(opts, "ReportVerified", feedIdRule)
	if err != nil {
		return nil, err
	}
	return &DestinationVerifierReportVerifiedIterator{contract: _DestinationVerifier.contract, event: "ReportVerified", logs: logs, sub: sub}, nil
}

func (_DestinationVerifier *DestinationVerifierFilterer) WatchReportVerified(opts *bind.WatchOpts, sink chan<- *DestinationVerifierReportVerified, feedId [][32]byte) (event.Subscription, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _DestinationVerifier.contract.WatchLogs(opts, "ReportVerified", feedIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DestinationVerifierReportVerified)
				if err := _DestinationVerifier.contract.UnpackLog(event, "ReportVerified", log); err != nil {
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

func (_DestinationVerifier *DestinationVerifierFilterer) ParseReportVerified(log types.Log) (*DestinationVerifierReportVerified, error) {
	event := new(DestinationVerifierReportVerified)
	if err := _DestinationVerifier.contract.UnpackLog(event, "ReportVerified", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_DestinationVerifier *DestinationVerifier) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _DestinationVerifier.abi.Events["AccessControllerSet"].ID:
		return _DestinationVerifier.ParseAccessControllerSet(log)
	case _DestinationVerifier.abi.Events["ConfigActivated"].ID:
		return _DestinationVerifier.ParseConfigActivated(log)
	case _DestinationVerifier.abi.Events["ConfigSet"].ID:
		return _DestinationVerifier.ParseConfigSet(log)
	case _DestinationVerifier.abi.Events["FeeManagerSet"].ID:
		return _DestinationVerifier.ParseFeeManagerSet(log)
	case _DestinationVerifier.abi.Events["OwnershipTransferRequested"].ID:
		return _DestinationVerifier.ParseOwnershipTransferRequested(log)
	case _DestinationVerifier.abi.Events["OwnershipTransferred"].ID:
		return _DestinationVerifier.ParseOwnershipTransferred(log)
	case _DestinationVerifier.abi.Events["ReportVerified"].ID:
		return _DestinationVerifier.ParseReportVerified(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (DestinationVerifierAccessControllerSet) Topic() common.Hash {
	return common.HexToHash("0x953e92b1a6442e9c3242531154a3f6f6eb00b4e9c719ba8118fa6235e4ce89b6")
}

func (DestinationVerifierConfigActivated) Topic() common.Hash {
	return common.HexToHash("0x90186a1e77b498ec417ea88bd026cae00d7043c357cc45221777623bda582dd4")
}

func (DestinationVerifierConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2d763a674a99583454a287d792819ffb9ff7e791c23e7745a082701136ce336c")
}

func (DestinationVerifierFeeManagerSet) Topic() common.Hash {
	return common.HexToHash("0x04628abcaa6b1674651352125cb94b65b289145bc2bc4d67720bb7d966372f03")
}

func (DestinationVerifierOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (DestinationVerifierOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (DestinationVerifierReportVerified) Topic() common.Hash {
	return common.HexToHash("0x58ca9502e98a536e06e72d680fcc251e5d10b72291a281665a2c2dc0ac30fcc5")
}

func (_DestinationVerifier *DestinationVerifier) Address() common.Address {
	return _DestinationVerifier.address
}

type DestinationVerifierInterface interface {
	GetAccessController(opts *bind.CallOpts) (common.Address, error)

	GetFeeManager(opts *bind.CallOpts) (common.Address, error)

	IVerifierProxy(opts *bind.CallOpts) (common.Address, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	SetAccessController(opts *bind.TransactOpts, accessController common.Address) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, signers []common.Address, f uint8, recipientAddressesAndWeights []CommonAddressAndWeight) (*types.Transaction, error)

	SetConfigActive(opts *bind.TransactOpts, DONConfigIndex *big.Int, isActive bool) (*types.Transaction, error)

	SetFeeManager(opts *bind.TransactOpts, feeManager common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Verify(opts *bind.TransactOpts, signedReport []byte, parameterPayload []byte, sender common.Address) (*types.Transaction, error)

	VerifyBulk(opts *bind.TransactOpts, signedReports [][]byte, parameterPayload []byte, sender common.Address) (*types.Transaction, error)

	FilterAccessControllerSet(opts *bind.FilterOpts) (*DestinationVerifierAccessControllerSetIterator, error)

	WatchAccessControllerSet(opts *bind.WatchOpts, sink chan<- *DestinationVerifierAccessControllerSet) (event.Subscription, error)

	ParseAccessControllerSet(log types.Log) (*DestinationVerifierAccessControllerSet, error)

	FilterConfigActivated(opts *bind.FilterOpts) (*DestinationVerifierConfigActivatedIterator, error)

	WatchConfigActivated(opts *bind.WatchOpts, sink chan<- *DestinationVerifierConfigActivated) (event.Subscription, error)

	ParseConfigActivated(log types.Log) (*DestinationVerifierConfigActivated, error)

	FilterConfigSet(opts *bind.FilterOpts, DONConfigID [][24]byte) (*DestinationVerifierConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *DestinationVerifierConfigSet, DONConfigID [][24]byte) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*DestinationVerifierConfigSet, error)

	FilterFeeManagerSet(opts *bind.FilterOpts) (*DestinationVerifierFeeManagerSetIterator, error)

	WatchFeeManagerSet(opts *bind.WatchOpts, sink chan<- *DestinationVerifierFeeManagerSet) (event.Subscription, error)

	ParseFeeManagerSet(log types.Log) (*DestinationVerifierFeeManagerSet, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*DestinationVerifierOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *DestinationVerifierOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*DestinationVerifierOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*DestinationVerifierOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DestinationVerifierOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*DestinationVerifierOwnershipTransferred, error)

	FilterReportVerified(opts *bind.FilterOpts, feedId [][32]byte) (*DestinationVerifierReportVerifiedIterator, error)

	WatchReportVerified(opts *bind.WatchOpts, sink chan<- *DestinationVerifierReportVerified, feedId [][32]byte) (event.Subscription, error)

	ParseReportVerified(log types.Log) (*DestinationVerifierReportVerified, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
