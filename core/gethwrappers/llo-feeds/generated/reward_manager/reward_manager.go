// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package reward_manager

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

type IRewardManagerFeePayment struct {
	PoolId [32]byte
	Amount *big.Int
}

var RewardManagerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"linkAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPoolId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPoolLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidWeights\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Unauthorized\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newFeeManagerAddress\",\"type\":\"address\"}],\"name\":\"FeeManagerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"internalType\":\"uint192\",\"name\":\"amount\",\"type\":\"uint192\"}],\"indexed\":false,\"internalType\":\"structIRewardManager.FeePayment[]\",\"name\":\"payments\",\"type\":\"tuple[]\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payer\",\"type\":\"address\"}],\"name\":\"FeePaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"weight\",\"type\":\"uint64\"}],\"indexed\":false,\"internalType\":\"structCommon.AddressAndWeight[]\",\"name\":\"newRewardRecipients\",\"type\":\"tuple[]\"}],\"name\":\"RewardRecipientsUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint192\",\"name\":\"quantity\",\"type\":\"uint192\"}],\"name\":\"RewardsClaimed\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"poolIds\",\"type\":\"bytes32[]\"}],\"name\":\"claimRewards\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endIndex\",\"type\":\"uint256\"}],\"name\":\"getAvailableRewardPoolIds\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"internalType\":\"uint192\",\"name\":\"amount\",\"type\":\"uint192\"}],\"internalType\":\"structIRewardManager.FeePayment[]\",\"name\":\"payments\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"payer\",\"type\":\"address\"}],\"name\":\"onFeePaid\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"internalType\":\"address[]\",\"name\":\"recipients\",\"type\":\"address[]\"}],\"name\":\"payRecipients\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_feeManagerAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_registeredPoolIds\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"s_rewardRecipientWeights\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"s_rewardRecipientWeightsSet\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"s_totalRewardRecipientFees\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"s_totalRewardRecipientFeesLastClaimedAmounts\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newFeeManagerAddress\",\"type\":\"address\"}],\"name\":\"setFeeManager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"weight\",\"type\":\"uint64\"}],\"internalType\":\"structCommon.AddressAndWeight[]\",\"name\":\"rewardRecipientAndWeights\",\"type\":\"tuple[]\"}],\"name\":\"setRewardRecipients\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"weight\",\"type\":\"uint64\"}],\"internalType\":\"structCommon.AddressAndWeight[]\",\"name\":\"newRewardRecipients\",\"type\":\"tuple[]\"}],\"name\":\"updateRewardRecipients\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b5060405162002079380380620020798339810160408190526200003491620001a6565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000fb565b5050506001600160a01b038116620000e95760405163e6c4247b60e01b815260040160405180910390fd5b6001600160a01b0316608052620001d8565b336001600160a01b03821603620001555760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600060208284031215620001b957600080fd5b81516001600160a01b0381168114620001d157600080fd5b9392505050565b608051611e7e620001fb60003960008181610ca10152610edc0152611e7e6000f3fe608060405234801561001057600080fd5b50600436106101365760003560e01c80634d322084116100b25780638ac85a5c11610081578063b0d9fa1911610066578063b0d9fa1914610377578063cd5f72921461038a578063f2fde38b1461039d57600080fd5b80638ac85a5c1461032e5780638da5cb5b1461035957600080fd5b80634d322084146102c557806359256201146102d857806360122608146102fb57806379ba50971461032657600080fd5b8063276e76601161010957806347226475116100ee578063472264751461027f578063472d35b91461029f5780634944832f146102b257600080fd5b8063276e76601461020c57806339ee81e11461025157600080fd5b806301ffc9a71461013b5780630f3c34d1146101a557806314060f23146101ba578063181f5a77146101cd575b600080fd5b61019061014936600461176a565b7fffffffff00000000000000000000000000000000000000000000000000000000167fb0d9fa19000000000000000000000000000000000000000000000000000000001490565b60405190151581526020015b60405180910390f35b6101b86101b336600461182a565b6103b0565b005b6101b86101c836600461191c565b6103be565b604080518082018252601381527f5265776172644d616e6167657220312e302e30000000000000000000000000006020820152905161019c919061198c565b60075461022c9073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161019c565b61027161025f3660046119dd565b60026020526000908152604090205481565b60405190815260200161019c565b61029261028d366004611a1f565b610574565b60405161019c9190611a52565b6101b86102ad366004611a96565b6106fe565b6101b86102c036600461191c565b6107cc565b6101b86102d3366004611ab8565b610915565b6101906102e63660046119dd565b60056020526000908152604090205460ff1681565b610271610309366004611b37565b600360209081526000928352604080842090915290825290205481565b6101b8610a54565b61027161033c366004611b37565b600460209081526000928352604080842090915290825290205481565b60005473ffffffffffffffffffffffffffffffffffffffff1661022c565b6101b8610385366004611b63565b610b56565b6102716103983660046119dd565b610d0a565b6101b86103ab366004611a96565b610d2b565b6103ba3382610d3f565b5050565b60005473ffffffffffffffffffffffffffffffffffffffff1633148015906103fe575060075473ffffffffffffffffffffffffffffffffffffffff163314155b15610435576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000819003610470576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008381526005602052604090205460ff16156104b9576040517f0afa7ee800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6006805460018181019092557ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d3f01849055600084815260056020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169091179055610535838383670de0b6b3a7640000610f0c565b827f8f668d6090683f98b3373a8b83d214da45737f7486cb7de554cc07b54e61cfe68383604051610567929190611bcf565b60405180910390a2505050565b6006546060906000818411610589578361058b565b815b9050808511156105c7576040517fa22caccc00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006105d38683611c66565b67ffffffffffffffff8111156105eb576105eb6117ac565b604051908082528060200260200182016040528015610614578160200160208202803683370190505b5090506000865b838110156106f15760006006828154811061063857610638611c79565b600091825260208083209091015480835260048252604080842073ffffffffffffffffffffffffffffffffffffffff8f168552909252912054909150156106e0576000818152600260209081526040808320546003835281842073ffffffffffffffffffffffffffffffffffffffff8f1685529092529091205481146106de57818585806001019650815181106106d1576106d1611c79565b6020026020010181815250505b505b506106ea81611ca8565b905061061b565b5090979650505050505050565b6107066110e3565b73ffffffffffffffffffffffffffffffffffffffff8116610753576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600780547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83169081179091556040519081527fe45f5e140399b0a7e12971ab020724b828fbed8ac408c420884dc7d1bbe506b49060200160405180910390a150565b6107d46110e3565b60408051600180825281830190925260009160208083019080368337019050509050838160008151811061080a5761080a611c79565b6020026020010181815250506000805b838110156108c757600085858381811061083657610836611c79565b61084c9260206040909202019081019150611a96565b600088815260046020908152604080832073ffffffffffffffffffffffffffffffffffffffff851684529091529020549091506108b087878581811061089457610894611c79565b6108aa9260206040909202019081019150611a96565b86610d3f565b509290920191506108c081611ca8565b905061081a565b506108d485858584610f0c565b847f8f668d6090683f98b3373a8b83d214da45737f7486cb7de554cc07b54e61cfe68585604051610906929190611bcf565b60405180910390a25050505050565b8261093560005473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161415801561098757506000818152600460209081526040808320338452909152902054155b156109be576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805160018082528183019092526000916020808301908036833701905050905084816000815181106109f4576109f4611c79565b60200260200101818152505060005b83811015610a4c57610a3b858583818110610a2057610a20611c79565b9050602002016020810190610a359190611a96565b83610d3f565b50610a4581611ca8565b9050610a03565b505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610ada576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60075473ffffffffffffffffffffffffffffffffffffffff163314610ba7576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000805b83811015610c8657848482818110610bc557610bc5611c79565b9050604002016020016020810190610bdd9190611d08565b77ffffffffffffffffffffffffffffffffffffffffffffffff1660026000878785818110610c0d57610c0d611c79565b6040908102929092013583525060208201929092520160002080549091019055848482818110610c3f57610c3f611c79565b9050604002016020016020810190610c579190611d08565b77ffffffffffffffffffffffffffffffffffffffffffffffff168201915080610c7f90611ca8565b9050610bab565b50610cc973ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016833084611166565b7fa1cc025ea76bacce5d740ee4bc331899375dc2c5f2ab33933aaacbd9ba001b66848484604051610cfc93929190611d23565b60405180910390a150505050565b60068181548110610d1a57600080fd5b600091825260209091200154905081565b610d336110e3565b610d3c81611248565b50565b60008060005b8351811015610ebb576000848281518110610d6257610d62611c79565b6020026020010151905060006002600083815260200190815260200160002054905080600003610d93575050610eab565b600082815260036020908152604080832073ffffffffffffffffffffffffffffffffffffffff8b16808552908352818420548685526004845282852091855292528220549083039190670de0b6b3a764000090830204905080600003610dfc5750505050610eab565b600084815260036020908152604080832073ffffffffffffffffffffffffffffffffffffffff8d168085529252909120849055885196820196899087908110610e4757610e47611c79565b60200260200101517f989969655bc1d593922527fe85d71347bb8e12fa423cc71f362dd8ef7cb10ef283604051610e9e919077ffffffffffffffffffffffffffffffffffffffffffffffff91909116815260200190565b60405180910390a3505050505b610eb481611ca8565b9050610d45565b508015610f0357610f0373ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016858361133d565b90505b92915050565b610f678383808060200260200160405190810160405280939291908181526020016000905b82821015610f5d57610f4e60408302860136819003810190611daa565b81526020019060010190610f31565b5050505050611398565b15610f9e576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000805b838110156110a2576000858583818110610fbe57610fbe611c79565b9050604002016020016020810190610fd69190611e05565b67ffffffffffffffff1690506000868684818110610ff657610ff6611c79565b61100c9260206040909202019081019150611a96565b905073ffffffffffffffffffffffffffffffffffffffff811661105b576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600088815260046020908152604080832073ffffffffffffffffffffffffffffffffffffffff909416835292905220819055919091019061109b81611ca8565b9050610fa2565b508181146110dc576040517f84677ce800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314611164576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610ad1565b565b60405173ffffffffffffffffffffffffffffffffffffffff808516602483015283166044820152606481018290526112429085907f23b872dd00000000000000000000000000000000000000000000000000000000906084015b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009093169290921790915261144f565b50505050565b3373ffffffffffffffffffffffffffffffffffffffff8216036112c7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610ad1565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60405173ffffffffffffffffffffffffffffffffffffffff83166024820152604481018290526113939084907fa9059cbb00000000000000000000000000000000000000000000000000000000906064016111c0565b505050565b6000805b82518110156114465760006113b2826001611e20565b90505b835181101561143d578381815181106113d0576113d0611c79565b60200260200101516000015173ffffffffffffffffffffffffffffffffffffffff1684838151811061140457611404611c79565b60200260200101516000015173ffffffffffffffffffffffffffffffffffffffff1603611435575060019392505050565b6001016113b5565b5060010161139c565b50600092915050565b60006114b1826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff1661155b9092919063ffffffff16565b80519091501561139357808060200190518101906114cf9190611e33565b611393576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f742073756363656564000000000000000000000000000000000000000000006064820152608401610ad1565b606061156a8484600085611572565b949350505050565b606082471015611604576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f60448201527f722063616c6c00000000000000000000000000000000000000000000000000006064820152608401610ad1565b6000808673ffffffffffffffffffffffffffffffffffffffff16858760405161162d9190611e55565b60006040518083038185875af1925050503d806000811461166a576040519150601f19603f3d011682016040523d82523d6000602084013e61166f565b606091505b50915091506116808783838761168b565b979650505050505050565b6060831561172157825160000361171a5773ffffffffffffffffffffffffffffffffffffffff85163b61171a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401610ad1565b508161156a565b61156a83838151156117365781518083602001fd5b806040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610ad1919061198c565b60006020828403121561177c57600080fd5b81357fffffffff0000000000000000000000000000000000000000000000000000000081168114610f0357600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715611822576118226117ac565b604052919050565b6000602080838503121561183d57600080fd5b823567ffffffffffffffff8082111561185557600080fd5b818501915085601f83011261186957600080fd5b81358181111561187b5761187b6117ac565b8060051b915061188c8483016117db565b81815291830184019184810190888411156118a657600080fd5b938501935b838510156118c4578435825293850193908501906118ab565b98975050505050505050565b60008083601f8401126118e257600080fd5b50813567ffffffffffffffff8111156118fa57600080fd5b6020830191508360208260061b850101111561191557600080fd5b9250929050565b60008060006040848603121561193157600080fd5b83359250602084013567ffffffffffffffff81111561194f57600080fd5b61195b868287016118d0565b9497909650939450505050565b60005b8381101561198357818101518382015260200161196b565b50506000910152565b60208152600082518060208401526119ab816040850160208701611968565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169190910160400192915050565b6000602082840312156119ef57600080fd5b5035919050565b803573ffffffffffffffffffffffffffffffffffffffff81168114611a1a57600080fd5b919050565b600080600060608486031215611a3457600080fd5b611a3d846119f6565b95602085013595506040909401359392505050565b6020808252825182820181905260009190848201906040850190845b81811015611a8a57835183529284019291840191600101611a6e565b50909695505050505050565b600060208284031215611aa857600080fd5b611ab1826119f6565b9392505050565b600080600060408486031215611acd57600080fd5b83359250602084013567ffffffffffffffff80821115611aec57600080fd5b818601915086601f830112611b0057600080fd5b813581811115611b0f57600080fd5b8760208260051b8501011115611b2457600080fd5b6020830194508093505050509250925092565b60008060408385031215611b4a57600080fd5b82359150611b5a602084016119f6565b90509250929050565b600080600060408486031215611b7857600080fd5b833567ffffffffffffffff811115611b8f57600080fd5b611b9b868287016118d0565b9094509250611bae9050602085016119f6565b90509250925092565b803567ffffffffffffffff81168114611a1a57600080fd5b6020808252818101839052600090604080840186845b878110156106f15773ffffffffffffffffffffffffffffffffffffffff611c0b836119f6565b16835267ffffffffffffffff611c22868401611bb7565b16838601529183019190830190600101611be5565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b81810381811115610f0657610f06611c37565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203611cd957611cd9611c37565b5060010190565b803577ffffffffffffffffffffffffffffffffffffffffffffffff81168114611a1a57600080fd5b600060208284031215611d1a57600080fd5b611ab182611ce0565b60408082528181018490526000908560608401835b87811015611d7f5782358252602077ffffffffffffffffffffffffffffffffffffffffffffffff611d6a828601611ce0565b16908301529183019190830190600101611d38565b5080935050505073ffffffffffffffffffffffffffffffffffffffff83166020830152949350505050565b600060408284031215611dbc57600080fd5b6040516040810181811067ffffffffffffffff82111715611ddf57611ddf6117ac565b604052611deb836119f6565b8152611df960208401611bb7565b60208201529392505050565b600060208284031215611e1757600080fd5b611ab182611bb7565b80820180821115610f0657610f06611c37565b600060208284031215611e4557600080fd5b81518015158114610f0357600080fd5b60008251611e67818460208701611968565b919091019291505056fea164736f6c6343000810000a",
}

var RewardManagerABI = RewardManagerMetaData.ABI

var RewardManagerBin = RewardManagerMetaData.Bin

func DeployRewardManager(auth *bind.TransactOpts, backend bind.ContractBackend, linkAddress common.Address) (common.Address, *types.Transaction, *RewardManager, error) {
	parsed, err := RewardManagerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(RewardManagerBin), backend, linkAddress)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &RewardManager{RewardManagerCaller: RewardManagerCaller{contract: contract}, RewardManagerTransactor: RewardManagerTransactor{contract: contract}, RewardManagerFilterer: RewardManagerFilterer{contract: contract}}, nil
}

type RewardManager struct {
	address common.Address
	abi     abi.ABI
	RewardManagerCaller
	RewardManagerTransactor
	RewardManagerFilterer
}

type RewardManagerCaller struct {
	contract *bind.BoundContract
}

type RewardManagerTransactor struct {
	contract *bind.BoundContract
}

type RewardManagerFilterer struct {
	contract *bind.BoundContract
}

type RewardManagerSession struct {
	Contract     *RewardManager
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type RewardManagerCallerSession struct {
	Contract *RewardManagerCaller
	CallOpts bind.CallOpts
}

type RewardManagerTransactorSession struct {
	Contract     *RewardManagerTransactor
	TransactOpts bind.TransactOpts
}

type RewardManagerRaw struct {
	Contract *RewardManager
}

type RewardManagerCallerRaw struct {
	Contract *RewardManagerCaller
}

type RewardManagerTransactorRaw struct {
	Contract *RewardManagerTransactor
}

func NewRewardManager(address common.Address, backend bind.ContractBackend) (*RewardManager, error) {
	abi, err := abi.JSON(strings.NewReader(RewardManagerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindRewardManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RewardManager{address: address, abi: abi, RewardManagerCaller: RewardManagerCaller{contract: contract}, RewardManagerTransactor: RewardManagerTransactor{contract: contract}, RewardManagerFilterer: RewardManagerFilterer{contract: contract}}, nil
}

func NewRewardManagerCaller(address common.Address, caller bind.ContractCaller) (*RewardManagerCaller, error) {
	contract, err := bindRewardManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RewardManagerCaller{contract: contract}, nil
}

func NewRewardManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*RewardManagerTransactor, error) {
	contract, err := bindRewardManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RewardManagerTransactor{contract: contract}, nil
}

func NewRewardManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*RewardManagerFilterer, error) {
	contract, err := bindRewardManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RewardManagerFilterer{contract: contract}, nil
}

func bindRewardManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RewardManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_RewardManager *RewardManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RewardManager.Contract.RewardManagerCaller.contract.Call(opts, result, method, params...)
}

func (_RewardManager *RewardManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RewardManager.Contract.RewardManagerTransactor.contract.Transfer(opts)
}

func (_RewardManager *RewardManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RewardManager.Contract.RewardManagerTransactor.contract.Transact(opts, method, params...)
}

func (_RewardManager *RewardManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RewardManager.Contract.contract.Call(opts, result, method, params...)
}

func (_RewardManager *RewardManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RewardManager.Contract.contract.Transfer(opts)
}

func (_RewardManager *RewardManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RewardManager.Contract.contract.Transact(opts, method, params...)
}

func (_RewardManager *RewardManagerCaller) GetAvailableRewardPoolIds(opts *bind.CallOpts, recipient common.Address, startIndex *big.Int, endIndex *big.Int) ([][32]byte, error) {
	var out []interface{}
	err := _RewardManager.contract.Call(opts, &out, "getAvailableRewardPoolIds", recipient, startIndex, endIndex)

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

func (_RewardManager *RewardManagerSession) GetAvailableRewardPoolIds(recipient common.Address, startIndex *big.Int, endIndex *big.Int) ([][32]byte, error) {
	return _RewardManager.Contract.GetAvailableRewardPoolIds(&_RewardManager.CallOpts, recipient, startIndex, endIndex)
}

func (_RewardManager *RewardManagerCallerSession) GetAvailableRewardPoolIds(recipient common.Address, startIndex *big.Int, endIndex *big.Int) ([][32]byte, error) {
	return _RewardManager.Contract.GetAvailableRewardPoolIds(&_RewardManager.CallOpts, recipient, startIndex, endIndex)
}

func (_RewardManager *RewardManagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RewardManager.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_RewardManager *RewardManagerSession) Owner() (common.Address, error) {
	return _RewardManager.Contract.Owner(&_RewardManager.CallOpts)
}

func (_RewardManager *RewardManagerCallerSession) Owner() (common.Address, error) {
	return _RewardManager.Contract.Owner(&_RewardManager.CallOpts)
}

func (_RewardManager *RewardManagerCaller) SFeeManagerAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RewardManager.contract.Call(opts, &out, "s_feeManagerAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_RewardManager *RewardManagerSession) SFeeManagerAddress() (common.Address, error) {
	return _RewardManager.Contract.SFeeManagerAddress(&_RewardManager.CallOpts)
}

func (_RewardManager *RewardManagerCallerSession) SFeeManagerAddress() (common.Address, error) {
	return _RewardManager.Contract.SFeeManagerAddress(&_RewardManager.CallOpts)
}

func (_RewardManager *RewardManagerCaller) SRegisteredPoolIds(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _RewardManager.contract.Call(opts, &out, "s_registeredPoolIds", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_RewardManager *RewardManagerSession) SRegisteredPoolIds(arg0 *big.Int) ([32]byte, error) {
	return _RewardManager.Contract.SRegisteredPoolIds(&_RewardManager.CallOpts, arg0)
}

func (_RewardManager *RewardManagerCallerSession) SRegisteredPoolIds(arg0 *big.Int) ([32]byte, error) {
	return _RewardManager.Contract.SRegisteredPoolIds(&_RewardManager.CallOpts, arg0)
}

func (_RewardManager *RewardManagerCaller) SRewardRecipientWeights(opts *bind.CallOpts, arg0 [32]byte, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _RewardManager.contract.Call(opts, &out, "s_rewardRecipientWeights", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_RewardManager *RewardManagerSession) SRewardRecipientWeights(arg0 [32]byte, arg1 common.Address) (*big.Int, error) {
	return _RewardManager.Contract.SRewardRecipientWeights(&_RewardManager.CallOpts, arg0, arg1)
}

func (_RewardManager *RewardManagerCallerSession) SRewardRecipientWeights(arg0 [32]byte, arg1 common.Address) (*big.Int, error) {
	return _RewardManager.Contract.SRewardRecipientWeights(&_RewardManager.CallOpts, arg0, arg1)
}

func (_RewardManager *RewardManagerCaller) SRewardRecipientWeightsSet(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _RewardManager.contract.Call(opts, &out, "s_rewardRecipientWeightsSet", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_RewardManager *RewardManagerSession) SRewardRecipientWeightsSet(arg0 [32]byte) (bool, error) {
	return _RewardManager.Contract.SRewardRecipientWeightsSet(&_RewardManager.CallOpts, arg0)
}

func (_RewardManager *RewardManagerCallerSession) SRewardRecipientWeightsSet(arg0 [32]byte) (bool, error) {
	return _RewardManager.Contract.SRewardRecipientWeightsSet(&_RewardManager.CallOpts, arg0)
}

func (_RewardManager *RewardManagerCaller) STotalRewardRecipientFees(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _RewardManager.contract.Call(opts, &out, "s_totalRewardRecipientFees", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_RewardManager *RewardManagerSession) STotalRewardRecipientFees(arg0 [32]byte) (*big.Int, error) {
	return _RewardManager.Contract.STotalRewardRecipientFees(&_RewardManager.CallOpts, arg0)
}

func (_RewardManager *RewardManagerCallerSession) STotalRewardRecipientFees(arg0 [32]byte) (*big.Int, error) {
	return _RewardManager.Contract.STotalRewardRecipientFees(&_RewardManager.CallOpts, arg0)
}

func (_RewardManager *RewardManagerCaller) STotalRewardRecipientFeesLastClaimedAmounts(opts *bind.CallOpts, arg0 [32]byte, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _RewardManager.contract.Call(opts, &out, "s_totalRewardRecipientFeesLastClaimedAmounts", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_RewardManager *RewardManagerSession) STotalRewardRecipientFeesLastClaimedAmounts(arg0 [32]byte, arg1 common.Address) (*big.Int, error) {
	return _RewardManager.Contract.STotalRewardRecipientFeesLastClaimedAmounts(&_RewardManager.CallOpts, arg0, arg1)
}

func (_RewardManager *RewardManagerCallerSession) STotalRewardRecipientFeesLastClaimedAmounts(arg0 [32]byte, arg1 common.Address) (*big.Int, error) {
	return _RewardManager.Contract.STotalRewardRecipientFeesLastClaimedAmounts(&_RewardManager.CallOpts, arg0, arg1)
}

func (_RewardManager *RewardManagerCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _RewardManager.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_RewardManager *RewardManagerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _RewardManager.Contract.SupportsInterface(&_RewardManager.CallOpts, interfaceId)
}

func (_RewardManager *RewardManagerCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _RewardManager.Contract.SupportsInterface(&_RewardManager.CallOpts, interfaceId)
}

func (_RewardManager *RewardManagerCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _RewardManager.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_RewardManager *RewardManagerSession) TypeAndVersion() (string, error) {
	return _RewardManager.Contract.TypeAndVersion(&_RewardManager.CallOpts)
}

func (_RewardManager *RewardManagerCallerSession) TypeAndVersion() (string, error) {
	return _RewardManager.Contract.TypeAndVersion(&_RewardManager.CallOpts)
}

func (_RewardManager *RewardManagerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RewardManager.contract.Transact(opts, "acceptOwnership")
}

func (_RewardManager *RewardManagerSession) AcceptOwnership() (*types.Transaction, error) {
	return _RewardManager.Contract.AcceptOwnership(&_RewardManager.TransactOpts)
}

func (_RewardManager *RewardManagerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _RewardManager.Contract.AcceptOwnership(&_RewardManager.TransactOpts)
}

func (_RewardManager *RewardManagerTransactor) ClaimRewards(opts *bind.TransactOpts, poolIds [][32]byte) (*types.Transaction, error) {
	return _RewardManager.contract.Transact(opts, "claimRewards", poolIds)
}

func (_RewardManager *RewardManagerSession) ClaimRewards(poolIds [][32]byte) (*types.Transaction, error) {
	return _RewardManager.Contract.ClaimRewards(&_RewardManager.TransactOpts, poolIds)
}

func (_RewardManager *RewardManagerTransactorSession) ClaimRewards(poolIds [][32]byte) (*types.Transaction, error) {
	return _RewardManager.Contract.ClaimRewards(&_RewardManager.TransactOpts, poolIds)
}

func (_RewardManager *RewardManagerTransactor) OnFeePaid(opts *bind.TransactOpts, payments []IRewardManagerFeePayment, payer common.Address) (*types.Transaction, error) {
	return _RewardManager.contract.Transact(opts, "onFeePaid", payments, payer)
}

func (_RewardManager *RewardManagerSession) OnFeePaid(payments []IRewardManagerFeePayment, payer common.Address) (*types.Transaction, error) {
	return _RewardManager.Contract.OnFeePaid(&_RewardManager.TransactOpts, payments, payer)
}

func (_RewardManager *RewardManagerTransactorSession) OnFeePaid(payments []IRewardManagerFeePayment, payer common.Address) (*types.Transaction, error) {
	return _RewardManager.Contract.OnFeePaid(&_RewardManager.TransactOpts, payments, payer)
}

func (_RewardManager *RewardManagerTransactor) PayRecipients(opts *bind.TransactOpts, poolId [32]byte, recipients []common.Address) (*types.Transaction, error) {
	return _RewardManager.contract.Transact(opts, "payRecipients", poolId, recipients)
}

func (_RewardManager *RewardManagerSession) PayRecipients(poolId [32]byte, recipients []common.Address) (*types.Transaction, error) {
	return _RewardManager.Contract.PayRecipients(&_RewardManager.TransactOpts, poolId, recipients)
}

func (_RewardManager *RewardManagerTransactorSession) PayRecipients(poolId [32]byte, recipients []common.Address) (*types.Transaction, error) {
	return _RewardManager.Contract.PayRecipients(&_RewardManager.TransactOpts, poolId, recipients)
}

func (_RewardManager *RewardManagerTransactor) SetFeeManager(opts *bind.TransactOpts, newFeeManagerAddress common.Address) (*types.Transaction, error) {
	return _RewardManager.contract.Transact(opts, "setFeeManager", newFeeManagerAddress)
}

func (_RewardManager *RewardManagerSession) SetFeeManager(newFeeManagerAddress common.Address) (*types.Transaction, error) {
	return _RewardManager.Contract.SetFeeManager(&_RewardManager.TransactOpts, newFeeManagerAddress)
}

func (_RewardManager *RewardManagerTransactorSession) SetFeeManager(newFeeManagerAddress common.Address) (*types.Transaction, error) {
	return _RewardManager.Contract.SetFeeManager(&_RewardManager.TransactOpts, newFeeManagerAddress)
}

func (_RewardManager *RewardManagerTransactor) SetRewardRecipients(opts *bind.TransactOpts, poolId [32]byte, rewardRecipientAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _RewardManager.contract.Transact(opts, "setRewardRecipients", poolId, rewardRecipientAndWeights)
}

func (_RewardManager *RewardManagerSession) SetRewardRecipients(poolId [32]byte, rewardRecipientAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _RewardManager.Contract.SetRewardRecipients(&_RewardManager.TransactOpts, poolId, rewardRecipientAndWeights)
}

func (_RewardManager *RewardManagerTransactorSession) SetRewardRecipients(poolId [32]byte, rewardRecipientAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _RewardManager.Contract.SetRewardRecipients(&_RewardManager.TransactOpts, poolId, rewardRecipientAndWeights)
}

func (_RewardManager *RewardManagerTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _RewardManager.contract.Transact(opts, "transferOwnership", to)
}

func (_RewardManager *RewardManagerSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _RewardManager.Contract.TransferOwnership(&_RewardManager.TransactOpts, to)
}

func (_RewardManager *RewardManagerTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _RewardManager.Contract.TransferOwnership(&_RewardManager.TransactOpts, to)
}

func (_RewardManager *RewardManagerTransactor) UpdateRewardRecipients(opts *bind.TransactOpts, poolId [32]byte, newRewardRecipients []CommonAddressAndWeight) (*types.Transaction, error) {
	return _RewardManager.contract.Transact(opts, "updateRewardRecipients", poolId, newRewardRecipients)
}

func (_RewardManager *RewardManagerSession) UpdateRewardRecipients(poolId [32]byte, newRewardRecipients []CommonAddressAndWeight) (*types.Transaction, error) {
	return _RewardManager.Contract.UpdateRewardRecipients(&_RewardManager.TransactOpts, poolId, newRewardRecipients)
}

func (_RewardManager *RewardManagerTransactorSession) UpdateRewardRecipients(poolId [32]byte, newRewardRecipients []CommonAddressAndWeight) (*types.Transaction, error) {
	return _RewardManager.Contract.UpdateRewardRecipients(&_RewardManager.TransactOpts, poolId, newRewardRecipients)
}

type RewardManagerFeeManagerUpdatedIterator struct {
	Event *RewardManagerFeeManagerUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RewardManagerFeeManagerUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardManagerFeeManagerUpdated)
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
		it.Event = new(RewardManagerFeeManagerUpdated)
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

func (it *RewardManagerFeeManagerUpdatedIterator) Error() error {
	return it.fail
}

func (it *RewardManagerFeeManagerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RewardManagerFeeManagerUpdated struct {
	NewFeeManagerAddress common.Address
	Raw                  types.Log
}

func (_RewardManager *RewardManagerFilterer) FilterFeeManagerUpdated(opts *bind.FilterOpts) (*RewardManagerFeeManagerUpdatedIterator, error) {

	logs, sub, err := _RewardManager.contract.FilterLogs(opts, "FeeManagerUpdated")
	if err != nil {
		return nil, err
	}
	return &RewardManagerFeeManagerUpdatedIterator{contract: _RewardManager.contract, event: "FeeManagerUpdated", logs: logs, sub: sub}, nil
}

func (_RewardManager *RewardManagerFilterer) WatchFeeManagerUpdated(opts *bind.WatchOpts, sink chan<- *RewardManagerFeeManagerUpdated) (event.Subscription, error) {

	logs, sub, err := _RewardManager.contract.WatchLogs(opts, "FeeManagerUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RewardManagerFeeManagerUpdated)
				if err := _RewardManager.contract.UnpackLog(event, "FeeManagerUpdated", log); err != nil {
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

func (_RewardManager *RewardManagerFilterer) ParseFeeManagerUpdated(log types.Log) (*RewardManagerFeeManagerUpdated, error) {
	event := new(RewardManagerFeeManagerUpdated)
	if err := _RewardManager.contract.UnpackLog(event, "FeeManagerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RewardManagerFeePaidIterator struct {
	Event *RewardManagerFeePaid

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RewardManagerFeePaidIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardManagerFeePaid)
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
		it.Event = new(RewardManagerFeePaid)
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

func (it *RewardManagerFeePaidIterator) Error() error {
	return it.fail
}

func (it *RewardManagerFeePaidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RewardManagerFeePaid struct {
	Payments []IRewardManagerFeePayment
	Payer    common.Address
	Raw      types.Log
}

func (_RewardManager *RewardManagerFilterer) FilterFeePaid(opts *bind.FilterOpts) (*RewardManagerFeePaidIterator, error) {

	logs, sub, err := _RewardManager.contract.FilterLogs(opts, "FeePaid")
	if err != nil {
		return nil, err
	}
	return &RewardManagerFeePaidIterator{contract: _RewardManager.contract, event: "FeePaid", logs: logs, sub: sub}, nil
}

func (_RewardManager *RewardManagerFilterer) WatchFeePaid(opts *bind.WatchOpts, sink chan<- *RewardManagerFeePaid) (event.Subscription, error) {

	logs, sub, err := _RewardManager.contract.WatchLogs(opts, "FeePaid")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RewardManagerFeePaid)
				if err := _RewardManager.contract.UnpackLog(event, "FeePaid", log); err != nil {
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

func (_RewardManager *RewardManagerFilterer) ParseFeePaid(log types.Log) (*RewardManagerFeePaid, error) {
	event := new(RewardManagerFeePaid)
	if err := _RewardManager.contract.UnpackLog(event, "FeePaid", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RewardManagerOwnershipTransferRequestedIterator struct {
	Event *RewardManagerOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RewardManagerOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardManagerOwnershipTransferRequested)
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
		it.Event = new(RewardManagerOwnershipTransferRequested)
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

func (it *RewardManagerOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *RewardManagerOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RewardManagerOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_RewardManager *RewardManagerFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RewardManagerOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _RewardManager.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &RewardManagerOwnershipTransferRequestedIterator{contract: _RewardManager.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_RewardManager *RewardManagerFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *RewardManagerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _RewardManager.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RewardManagerOwnershipTransferRequested)
				if err := _RewardManager.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_RewardManager *RewardManagerFilterer) ParseOwnershipTransferRequested(log types.Log) (*RewardManagerOwnershipTransferRequested, error) {
	event := new(RewardManagerOwnershipTransferRequested)
	if err := _RewardManager.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RewardManagerOwnershipTransferredIterator struct {
	Event *RewardManagerOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RewardManagerOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardManagerOwnershipTransferred)
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
		it.Event = new(RewardManagerOwnershipTransferred)
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

func (it *RewardManagerOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *RewardManagerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RewardManagerOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_RewardManager *RewardManagerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RewardManagerOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _RewardManager.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &RewardManagerOwnershipTransferredIterator{contract: _RewardManager.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_RewardManager *RewardManagerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *RewardManagerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _RewardManager.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RewardManagerOwnershipTransferred)
				if err := _RewardManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_RewardManager *RewardManagerFilterer) ParseOwnershipTransferred(log types.Log) (*RewardManagerOwnershipTransferred, error) {
	event := new(RewardManagerOwnershipTransferred)
	if err := _RewardManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RewardManagerRewardRecipientsUpdatedIterator struct {
	Event *RewardManagerRewardRecipientsUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RewardManagerRewardRecipientsUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardManagerRewardRecipientsUpdated)
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
		it.Event = new(RewardManagerRewardRecipientsUpdated)
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

func (it *RewardManagerRewardRecipientsUpdatedIterator) Error() error {
	return it.fail
}

func (it *RewardManagerRewardRecipientsUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RewardManagerRewardRecipientsUpdated struct {
	PoolId              [32]byte
	NewRewardRecipients []CommonAddressAndWeight
	Raw                 types.Log
}

func (_RewardManager *RewardManagerFilterer) FilterRewardRecipientsUpdated(opts *bind.FilterOpts, poolId [][32]byte) (*RewardManagerRewardRecipientsUpdatedIterator, error) {

	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}

	logs, sub, err := _RewardManager.contract.FilterLogs(opts, "RewardRecipientsUpdated", poolIdRule)
	if err != nil {
		return nil, err
	}
	return &RewardManagerRewardRecipientsUpdatedIterator{contract: _RewardManager.contract, event: "RewardRecipientsUpdated", logs: logs, sub: sub}, nil
}

func (_RewardManager *RewardManagerFilterer) WatchRewardRecipientsUpdated(opts *bind.WatchOpts, sink chan<- *RewardManagerRewardRecipientsUpdated, poolId [][32]byte) (event.Subscription, error) {

	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}

	logs, sub, err := _RewardManager.contract.WatchLogs(opts, "RewardRecipientsUpdated", poolIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RewardManagerRewardRecipientsUpdated)
				if err := _RewardManager.contract.UnpackLog(event, "RewardRecipientsUpdated", log); err != nil {
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

func (_RewardManager *RewardManagerFilterer) ParseRewardRecipientsUpdated(log types.Log) (*RewardManagerRewardRecipientsUpdated, error) {
	event := new(RewardManagerRewardRecipientsUpdated)
	if err := _RewardManager.contract.UnpackLog(event, "RewardRecipientsUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RewardManagerRewardsClaimedIterator struct {
	Event *RewardManagerRewardsClaimed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RewardManagerRewardsClaimedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardManagerRewardsClaimed)
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
		it.Event = new(RewardManagerRewardsClaimed)
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

func (it *RewardManagerRewardsClaimedIterator) Error() error {
	return it.fail
}

func (it *RewardManagerRewardsClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RewardManagerRewardsClaimed struct {
	PoolId    [32]byte
	Recipient common.Address
	Quantity  *big.Int
	Raw       types.Log
}

func (_RewardManager *RewardManagerFilterer) FilterRewardsClaimed(opts *bind.FilterOpts, poolId [][32]byte, recipient []common.Address) (*RewardManagerRewardsClaimedIterator, error) {

	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _RewardManager.contract.FilterLogs(opts, "RewardsClaimed", poolIdRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &RewardManagerRewardsClaimedIterator{contract: _RewardManager.contract, event: "RewardsClaimed", logs: logs, sub: sub}, nil
}

func (_RewardManager *RewardManagerFilterer) WatchRewardsClaimed(opts *bind.WatchOpts, sink chan<- *RewardManagerRewardsClaimed, poolId [][32]byte, recipient []common.Address) (event.Subscription, error) {

	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _RewardManager.contract.WatchLogs(opts, "RewardsClaimed", poolIdRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RewardManagerRewardsClaimed)
				if err := _RewardManager.contract.UnpackLog(event, "RewardsClaimed", log); err != nil {
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

func (_RewardManager *RewardManagerFilterer) ParseRewardsClaimed(log types.Log) (*RewardManagerRewardsClaimed, error) {
	event := new(RewardManagerRewardsClaimed)
	if err := _RewardManager.contract.UnpackLog(event, "RewardsClaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_RewardManager *RewardManager) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _RewardManager.abi.Events["FeeManagerUpdated"].ID:
		return _RewardManager.ParseFeeManagerUpdated(log)
	case _RewardManager.abi.Events["FeePaid"].ID:
		return _RewardManager.ParseFeePaid(log)
	case _RewardManager.abi.Events["OwnershipTransferRequested"].ID:
		return _RewardManager.ParseOwnershipTransferRequested(log)
	case _RewardManager.abi.Events["OwnershipTransferred"].ID:
		return _RewardManager.ParseOwnershipTransferred(log)
	case _RewardManager.abi.Events["RewardRecipientsUpdated"].ID:
		return _RewardManager.ParseRewardRecipientsUpdated(log)
	case _RewardManager.abi.Events["RewardsClaimed"].ID:
		return _RewardManager.ParseRewardsClaimed(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (RewardManagerFeeManagerUpdated) Topic() common.Hash {
	return common.HexToHash("0xe45f5e140399b0a7e12971ab020724b828fbed8ac408c420884dc7d1bbe506b4")
}

func (RewardManagerFeePaid) Topic() common.Hash {
	return common.HexToHash("0xa1cc025ea76bacce5d740ee4bc331899375dc2c5f2ab33933aaacbd9ba001b66")
}

func (RewardManagerOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (RewardManagerOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (RewardManagerRewardRecipientsUpdated) Topic() common.Hash {
	return common.HexToHash("0x8f668d6090683f98b3373a8b83d214da45737f7486cb7de554cc07b54e61cfe6")
}

func (RewardManagerRewardsClaimed) Topic() common.Hash {
	return common.HexToHash("0x989969655bc1d593922527fe85d71347bb8e12fa423cc71f362dd8ef7cb10ef2")
}

func (_RewardManager *RewardManager) Address() common.Address {
	return _RewardManager.address
}

type RewardManagerInterface interface {
	GetAvailableRewardPoolIds(opts *bind.CallOpts, recipient common.Address, startIndex *big.Int, endIndex *big.Int) ([][32]byte, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SFeeManagerAddress(opts *bind.CallOpts) (common.Address, error)

	SRegisteredPoolIds(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error)

	SRewardRecipientWeights(opts *bind.CallOpts, arg0 [32]byte, arg1 common.Address) (*big.Int, error)

	SRewardRecipientWeightsSet(opts *bind.CallOpts, arg0 [32]byte) (bool, error)

	STotalRewardRecipientFees(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error)

	STotalRewardRecipientFeesLastClaimedAmounts(opts *bind.CallOpts, arg0 [32]byte, arg1 common.Address) (*big.Int, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	ClaimRewards(opts *bind.TransactOpts, poolIds [][32]byte) (*types.Transaction, error)

	OnFeePaid(opts *bind.TransactOpts, payments []IRewardManagerFeePayment, payer common.Address) (*types.Transaction, error)

	PayRecipients(opts *bind.TransactOpts, poolId [32]byte, recipients []common.Address) (*types.Transaction, error)

	SetFeeManager(opts *bind.TransactOpts, newFeeManagerAddress common.Address) (*types.Transaction, error)

	SetRewardRecipients(opts *bind.TransactOpts, poolId [32]byte, rewardRecipientAndWeights []CommonAddressAndWeight) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	UpdateRewardRecipients(opts *bind.TransactOpts, poolId [32]byte, newRewardRecipients []CommonAddressAndWeight) (*types.Transaction, error)

	FilterFeeManagerUpdated(opts *bind.FilterOpts) (*RewardManagerFeeManagerUpdatedIterator, error)

	WatchFeeManagerUpdated(opts *bind.WatchOpts, sink chan<- *RewardManagerFeeManagerUpdated) (event.Subscription, error)

	ParseFeeManagerUpdated(log types.Log) (*RewardManagerFeeManagerUpdated, error)

	FilterFeePaid(opts *bind.FilterOpts) (*RewardManagerFeePaidIterator, error)

	WatchFeePaid(opts *bind.WatchOpts, sink chan<- *RewardManagerFeePaid) (event.Subscription, error)

	ParseFeePaid(log types.Log) (*RewardManagerFeePaid, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RewardManagerOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *RewardManagerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*RewardManagerOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RewardManagerOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *RewardManagerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*RewardManagerOwnershipTransferred, error)

	FilterRewardRecipientsUpdated(opts *bind.FilterOpts, poolId [][32]byte) (*RewardManagerRewardRecipientsUpdatedIterator, error)

	WatchRewardRecipientsUpdated(opts *bind.WatchOpts, sink chan<- *RewardManagerRewardRecipientsUpdated, poolId [][32]byte) (event.Subscription, error)

	ParseRewardRecipientsUpdated(log types.Log) (*RewardManagerRewardRecipientsUpdated, error)

	FilterRewardsClaimed(opts *bind.FilterOpts, poolId [][32]byte, recipient []common.Address) (*RewardManagerRewardsClaimedIterator, error)

	WatchRewardsClaimed(opts *bind.WatchOpts, sink chan<- *RewardManagerRewardsClaimed, poolId [][32]byte, recipient []common.Address) (event.Subscription, error)

	ParseRewardsClaimed(log types.Log) (*RewardManagerRewardsClaimed, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
