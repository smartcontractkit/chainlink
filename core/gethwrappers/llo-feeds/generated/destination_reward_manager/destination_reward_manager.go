// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package destination_reward_manager

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

type IDestinationRewardManagerFeePayment struct {
	PoolId [32]byte
	Amount *big.Int
}

var DestinationRewardManagerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"linkAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPoolId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPoolLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidWeights\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Unauthorized\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newFeeManagerAddress\",\"type\":\"address\"}],\"name\":\"FeeManagerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"internalType\":\"uint192\",\"name\":\"amount\",\"type\":\"uint192\"}],\"indexed\":false,\"internalType\":\"structIDestinationRewardManager.FeePayment[]\",\"name\":\"payments\",\"type\":\"tuple[]\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payer\",\"type\":\"address\"}],\"name\":\"FeePaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"weight\",\"type\":\"uint64\"}],\"indexed\":false,\"internalType\":\"structCommon.AddressAndWeight[]\",\"name\":\"newRewardRecipients\",\"type\":\"tuple[]\"}],\"name\":\"RewardRecipientsUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint192\",\"name\":\"quantity\",\"type\":\"uint192\"}],\"name\":\"RewardsClaimed\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newFeeManagerAddress\",\"type\":\"address\"}],\"name\":\"addFeeManager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"poolIds\",\"type\":\"bytes32[]\"}],\"name\":\"claimRewards\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endIndex\",\"type\":\"uint256\"}],\"name\":\"getAvailableRewardPoolIds\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_linkAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"internalType\":\"uint192\",\"name\":\"amount\",\"type\":\"uint192\"}],\"internalType\":\"structIDestinationRewardManager.FeePayment[]\",\"name\":\"payments\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"payer\",\"type\":\"address\"}],\"name\":\"onFeePaid\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"internalType\":\"address[]\",\"name\":\"recipients\",\"type\":\"address[]\"}],\"name\":\"payRecipients\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"feeManagerAddress\",\"type\":\"address\"}],\"name\":\"removeFeeManager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"s_feeManagerAddressList\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_registeredPoolIds\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"s_rewardRecipientWeights\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"s_rewardRecipientWeightsSet\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"s_totalRewardRecipientFees\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"s_totalRewardRecipientFeesLastClaimedAmounts\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"weight\",\"type\":\"uint64\"}],\"internalType\":\"structCommon.AddressAndWeight[]\",\"name\":\"rewardRecipientAndWeights\",\"type\":\"tuple[]\"}],\"name\":\"setRewardRecipients\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"weight\",\"type\":\"uint64\"}],\"internalType\":\"structCommon.AddressAndWeight[]\",\"name\":\"newRewardRecipients\",\"type\":\"tuple[]\"}],\"name\":\"updateRewardRecipients\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b5060405162002205380380620022058339810160408190526200003491620001a6565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000fb565b5050506001600160a01b038116620000e95760405163e6c4247b60e01b815260040160405180910390fd5b6001600160a01b0316608052620001d8565b336001600160a01b03821603620001555760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600060208284031215620001b957600080fd5b81516001600160a01b0381168114620001d157600080fd5b9392505050565b60805161200362000202600039600081816103fd01528181610e26015261106101526120036000f3fe608060405234801561001057600080fd5b506004361061016c5760003560e01c806360122608116100cd5780639604d05f11610081578063cd5f729211610066578063cd5f7292146103e5578063ea4b861b146103f8578063f2fde38b1461041f57600080fd5b80639604d05f1461039c578063b0d9fa19146103d257600080fd5b80638115c9cc116100b25780638115c9cc1461031f5780638ac85a5c146103325780638da5cb5b1461035d57600080fd5b806360122608146102ec57806379ba50971461031757600080fd5b806339ee81e1116101245780634944832f116101095780634944832f146102a35780634d322084146102b657806359256201146102c957600080fd5b806339ee81e114610255578063472264751461028357600080fd5b806314060f231161015557806314060f23146101f0578063181f5a77146102035780631f2d32c31461024257600080fd5b806301ffc9a7146101715780630f3c34d1146101db575b600080fd5b6101c661017f3660046118ef565b7fffffffff00000000000000000000000000000000000000000000000000000000167f768ffd3a000000000000000000000000000000000000000000000000000000001490565b60405190151581526020015b60405180910390f35b6101ee6101e93660046119af565b610432565b005b6101ee6101fe366004611aa1565b610440565b604080518082018252601e81527f44657374696e6174696f6e5265776172644d616e6167657220302e342e300000602082015290516101d29190611b11565b6101ee610250366004611b8b565b610602565b610275610263366004611bad565b60026020526000908152604090205481565b6040519081526020016101d2565b610296610291366004611bc6565b61073a565b6040516101d29190611bf9565b6101ee6102b1366004611aa1565b6108c4565b6101ee6102c4366004611c3d565b610a0d565b6101c66102d7366004611bad565b60056020526000908152604090205460ff1681565b6102756102fa366004611cbc565b600360209081526000928352604080842090915290825290205481565b6101ee610b1b565b6101ee61032d366004611b8b565b610c1d565b610275610340366004611cbc565b600460209081526000928352604080842090915290825290205481565b60005473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101d2565b6103776103aa366004611b8b565b60076020526000908152604090205473ffffffffffffffffffffffffffffffffffffffff1681565b6101ee6103e0366004611ce8565b610ccf565b6102756103f3366004611bad565b610e8f565b6103777f000000000000000000000000000000000000000000000000000000000000000081565b6101ee61042d366004611b8b565b610eb0565b61043c3382610ec4565b5050565b3360008181526007602052604090205473ffffffffffffffffffffffffffffffffffffffff161480159061048c575060005473ffffffffffffffffffffffffffffffffffffffff163314155b156104c3576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008190036104fe576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008381526005602052604090205460ff1615610547576040517f0afa7ee800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6006805460018181019092557ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d3f01849055600084815260056020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690911790556105c3838383670de0b6b3a7640000611091565b827f8f668d6090683f98b3373a8b83d214da45737f7486cb7de554cc07b54e61cfe683836040516105f5929190611d54565b60405180910390a2505050565b61060a611268565b73ffffffffffffffffffffffffffffffffffffffff8116610657576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff81811660009081526007602052604090205416156106b6576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff811660008181526007602090815260409182902080547fffffffffffffffffffffffff0000000000000000000000000000000000000000168417905590519182527fe45f5e140399b0a7e12971ab020724b828fbed8ac408c420884dc7d1bbe506b4910160405180910390a150565b600654606090600081841161074f5783610751565b815b90508085111561078d576040517fa22caccc00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006107998683611deb565b67ffffffffffffffff8111156107b1576107b1611931565b6040519080825280602002602001820160405280156107da578160200160208202803683370190505b5090506000865b838110156108b7576000600682815481106107fe576107fe611dfe565b600091825260208083209091015480835260048252604080842073ffffffffffffffffffffffffffffffffffffffff8f168552909252912054909150156108a6576000818152600260209081526040808320546003835281842073ffffffffffffffffffffffffffffffffffffffff8f1685529092529091205481146108a4578185858060010196508151811061089757610897611dfe565b6020026020010181815250505b505b506108b081611e2d565b90506107e1565b5090979650505050505050565b6108cc611268565b60408051600180825281830190925260009160208083019080368337019050509050838160008151811061090257610902611dfe565b6020026020010181815250506000805b838110156109bf57600085858381811061092e5761092e611dfe565b6109449260206040909202019081019150611b8b565b600088815260046020908152604080832073ffffffffffffffffffffffffffffffffffffffff851684529091529020549091506109a887878581811061098c5761098c611dfe565b6109a29260206040909202019081019150611b8b565b86610ec4565b509290920191506109b881611e2d565b9050610912565b506109cc85858584611091565b847f8f668d6090683f98b3373a8b83d214da45737f7486cb7de554cc07b54e61cfe685856040516109fe929190611d54565b60405180910390a25050505050565b60008381526004602090815260408083203384529091529020548390158015610a4e575060005473ffffffffffffffffffffffffffffffffffffffff163314155b15610a85576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b604080516001808252818301909252600091602080830190803683370190505090508481600081518110610abb57610abb611dfe565b60200260200101818152505060005b83811015610b1357610b02858583818110610ae757610ae7611dfe565b9050602002016020810190610afc9190611b8b565b83610ec4565b50610b0c81611e2d565b9050610aca565b505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610ba1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610c25611268565b73ffffffffffffffffffffffffffffffffffffffff81811660009081526007602052604090205416610c83576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff16600090815260076020526040902080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169055565b3360008181526007602052604090205473ffffffffffffffffffffffffffffffffffffffff1614610d2c576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000805b83811015610e0b57848482818110610d4a57610d4a611dfe565b9050604002016020016020810190610d629190611e8d565b77ffffffffffffffffffffffffffffffffffffffffffffffff1660026000878785818110610d9257610d92611dfe565b6040908102929092013583525060208201929092520160002080549091019055848482818110610dc457610dc4611dfe565b9050604002016020016020810190610ddc9190611e8d565b77ffffffffffffffffffffffffffffffffffffffffffffffff168201915080610e0490611e2d565b9050610d30565b50610e4e73ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000168330846112eb565b7fa1cc025ea76bacce5d740ee4bc331899375dc2c5f2ab33933aaacbd9ba001b66848484604051610e8193929190611ea8565b60405180910390a150505050565b60068181548110610e9f57600080fd5b600091825260209091200154905081565b610eb8611268565b610ec1816113cd565b50565b60008060005b8351811015611040576000848281518110610ee757610ee7611dfe565b6020026020010151905060006002600083815260200190815260200160002054905080600003610f18575050611030565b600082815260036020908152604080832073ffffffffffffffffffffffffffffffffffffffff8b16808552908352818420548685526004845282852091855292528220549083039190670de0b6b3a764000090830204905080600003610f815750505050611030565b600084815260036020908152604080832073ffffffffffffffffffffffffffffffffffffffff8d168085529252909120849055885196820196899087908110610fcc57610fcc611dfe565b60200260200101517f989969655bc1d593922527fe85d71347bb8e12fa423cc71f362dd8ef7cb10ef283604051611023919077ffffffffffffffffffffffffffffffffffffffffffffffff91909116815260200190565b60405180910390a3505050505b61103981611e2d565b9050610eca565b5080156110885761108873ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001685836114c2565b90505b92915050565b6110ec8383808060200260200160405190810160405280939291908181526020016000905b828210156110e2576110d360408302860136819003810190611f2f565b815260200190600101906110b6565b505050505061151d565b15611123576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000805b8381101561122757600085858381811061114357611143611dfe565b905060400201602001602081019061115b9190611f8a565b67ffffffffffffffff169050600086868481811061117b5761117b611dfe565b6111919260206040909202019081019150611b8b565b905073ffffffffffffffffffffffffffffffffffffffff81166111e0576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600088815260046020908152604080832073ffffffffffffffffffffffffffffffffffffffff909416835292905220819055919091019061122081611e2d565b9050611127565b50818114611261576040517f84677ce800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050505050565b60005473ffffffffffffffffffffffffffffffffffffffff1633146112e9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610b98565b565b60405173ffffffffffffffffffffffffffffffffffffffff808516602483015283166044820152606481018290526113c79085907f23b872dd00000000000000000000000000000000000000000000000000000000906084015b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909316929092179091526115d4565b50505050565b3373ffffffffffffffffffffffffffffffffffffffff82160361144c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610b98565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60405173ffffffffffffffffffffffffffffffffffffffff83166024820152604481018290526115189084907fa9059cbb0000000000000000000000000000000000000000000000000000000090606401611345565b505050565b6000805b82518110156115cb576000611537826001611fa5565b90505b83518110156115c25783818151811061155557611555611dfe565b60200260200101516000015173ffffffffffffffffffffffffffffffffffffffff1684838151811061158957611589611dfe565b60200260200101516000015173ffffffffffffffffffffffffffffffffffffffff16036115ba575060019392505050565b60010161153a565b50600101611521565b50600092915050565b6000611636826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff166116e09092919063ffffffff16565b80519091501561151857808060200190518101906116549190611fb8565b611518576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f742073756363656564000000000000000000000000000000000000000000006064820152608401610b98565b60606116ef84846000856116f7565b949350505050565b606082471015611789576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f60448201527f722063616c6c00000000000000000000000000000000000000000000000000006064820152608401610b98565b6000808673ffffffffffffffffffffffffffffffffffffffff1685876040516117b29190611fda565b60006040518083038185875af1925050503d80600081146117ef576040519150601f19603f3d011682016040523d82523d6000602084013e6117f4565b606091505b509150915061180587838387611810565b979650505050505050565b606083156118a657825160000361189f5773ffffffffffffffffffffffffffffffffffffffff85163b61189f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401610b98565b50816116ef565b6116ef83838151156118bb5781518083602001fd5b806040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610b989190611b11565b60006020828403121561190157600080fd5b81357fffffffff000000000000000000000000000000000000000000000000000000008116811461108857600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156119a7576119a7611931565b604052919050565b600060208083850312156119c257600080fd5b823567ffffffffffffffff808211156119da57600080fd5b818501915085601f8301126119ee57600080fd5b813581811115611a0057611a00611931565b8060051b9150611a11848301611960565b8181529183018401918481019088841115611a2b57600080fd5b938501935b83851015611a4957843582529385019390850190611a30565b98975050505050505050565b60008083601f840112611a6757600080fd5b50813567ffffffffffffffff811115611a7f57600080fd5b6020830191508360208260061b8501011115611a9a57600080fd5b9250929050565b600080600060408486031215611ab657600080fd5b83359250602084013567ffffffffffffffff811115611ad457600080fd5b611ae086828701611a55565b9497909650939450505050565b60005b83811015611b08578181015183820152602001611af0565b50506000910152565b6020815260008251806020840152611b30816040850160208701611aed565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169190910160400192915050565b803573ffffffffffffffffffffffffffffffffffffffff81168114611b8657600080fd5b919050565b600060208284031215611b9d57600080fd5b611ba682611b62565b9392505050565b600060208284031215611bbf57600080fd5b5035919050565b600080600060608486031215611bdb57600080fd5b611be484611b62565b95602085013595506040909401359392505050565b6020808252825182820181905260009190848201906040850190845b81811015611c3157835183529284019291840191600101611c15565b50909695505050505050565b600080600060408486031215611c5257600080fd5b83359250602084013567ffffffffffffffff80821115611c7157600080fd5b818601915086601f830112611c8557600080fd5b813581811115611c9457600080fd5b8760208260051b8501011115611ca957600080fd5b6020830194508093505050509250925092565b60008060408385031215611ccf57600080fd5b82359150611cdf60208401611b62565b90509250929050565b600080600060408486031215611cfd57600080fd5b833567ffffffffffffffff811115611d1457600080fd5b611d2086828701611a55565b9094509250611d33905060208501611b62565b90509250925092565b803567ffffffffffffffff81168114611b8657600080fd5b6020808252818101839052600090604080840186845b878110156108b75773ffffffffffffffffffffffffffffffffffffffff611d9083611b62565b16835267ffffffffffffffff611da7868401611d3c565b16838601529183019190830190600101611d6a565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8181038181111561108b5761108b611dbc565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203611e5e57611e5e611dbc565b5060010190565b803577ffffffffffffffffffffffffffffffffffffffffffffffff81168114611b8657600080fd5b600060208284031215611e9f57600080fd5b611ba682611e65565b60408082528181018490526000908560608401835b87811015611f045782358252602077ffffffffffffffffffffffffffffffffffffffffffffffff611eef828601611e65565b16908301529183019190830190600101611ebd565b5080935050505073ffffffffffffffffffffffffffffffffffffffff83166020830152949350505050565b600060408284031215611f4157600080fd5b6040516040810181811067ffffffffffffffff82111715611f6457611f64611931565b604052611f7083611b62565b8152611f7e60208401611d3c565b60208201529392505050565b600060208284031215611f9c57600080fd5b611ba682611d3c565b8082018082111561108b5761108b611dbc565b600060208284031215611fca57600080fd5b8151801515811461108857600080fd5b60008251611fec818460208701611aed565b919091019291505056fea164736f6c6343000813000a",
}

var DestinationRewardManagerABI = DestinationRewardManagerMetaData.ABI

var DestinationRewardManagerBin = DestinationRewardManagerMetaData.Bin

func DeployDestinationRewardManager(auth *bind.TransactOpts, backend bind.ContractBackend, linkAddress common.Address) (common.Address, *types.Transaction, *DestinationRewardManager, error) {
	parsed, err := DestinationRewardManagerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DestinationRewardManagerBin), backend, linkAddress)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DestinationRewardManager{address: address, abi: *parsed, DestinationRewardManagerCaller: DestinationRewardManagerCaller{contract: contract}, DestinationRewardManagerTransactor: DestinationRewardManagerTransactor{contract: contract}, DestinationRewardManagerFilterer: DestinationRewardManagerFilterer{contract: contract}}, nil
}

type DestinationRewardManager struct {
	address common.Address
	abi     abi.ABI
	DestinationRewardManagerCaller
	DestinationRewardManagerTransactor
	DestinationRewardManagerFilterer
}

type DestinationRewardManagerCaller struct {
	contract *bind.BoundContract
}

type DestinationRewardManagerTransactor struct {
	contract *bind.BoundContract
}

type DestinationRewardManagerFilterer struct {
	contract *bind.BoundContract
}

type DestinationRewardManagerSession struct {
	Contract     *DestinationRewardManager
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type DestinationRewardManagerCallerSession struct {
	Contract *DestinationRewardManagerCaller
	CallOpts bind.CallOpts
}

type DestinationRewardManagerTransactorSession struct {
	Contract     *DestinationRewardManagerTransactor
	TransactOpts bind.TransactOpts
}

type DestinationRewardManagerRaw struct {
	Contract *DestinationRewardManager
}

type DestinationRewardManagerCallerRaw struct {
	Contract *DestinationRewardManagerCaller
}

type DestinationRewardManagerTransactorRaw struct {
	Contract *DestinationRewardManagerTransactor
}

func NewDestinationRewardManager(address common.Address, backend bind.ContractBackend) (*DestinationRewardManager, error) {
	abi, err := abi.JSON(strings.NewReader(DestinationRewardManagerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindDestinationRewardManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DestinationRewardManager{address: address, abi: abi, DestinationRewardManagerCaller: DestinationRewardManagerCaller{contract: contract}, DestinationRewardManagerTransactor: DestinationRewardManagerTransactor{contract: contract}, DestinationRewardManagerFilterer: DestinationRewardManagerFilterer{contract: contract}}, nil
}

func NewDestinationRewardManagerCaller(address common.Address, caller bind.ContractCaller) (*DestinationRewardManagerCaller, error) {
	contract, err := bindDestinationRewardManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DestinationRewardManagerCaller{contract: contract}, nil
}

func NewDestinationRewardManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*DestinationRewardManagerTransactor, error) {
	contract, err := bindDestinationRewardManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DestinationRewardManagerTransactor{contract: contract}, nil
}

func NewDestinationRewardManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*DestinationRewardManagerFilterer, error) {
	contract, err := bindDestinationRewardManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DestinationRewardManagerFilterer{contract: contract}, nil
}

func bindDestinationRewardManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DestinationRewardManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_DestinationRewardManager *DestinationRewardManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DestinationRewardManager.Contract.DestinationRewardManagerCaller.contract.Call(opts, result, method, params...)
}

func (_DestinationRewardManager *DestinationRewardManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DestinationRewardManager.Contract.DestinationRewardManagerTransactor.contract.Transfer(opts)
}

func (_DestinationRewardManager *DestinationRewardManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DestinationRewardManager.Contract.DestinationRewardManagerTransactor.contract.Transact(opts, method, params...)
}

func (_DestinationRewardManager *DestinationRewardManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DestinationRewardManager.Contract.contract.Call(opts, result, method, params...)
}

func (_DestinationRewardManager *DestinationRewardManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DestinationRewardManager.Contract.contract.Transfer(opts)
}

func (_DestinationRewardManager *DestinationRewardManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DestinationRewardManager.Contract.contract.Transact(opts, method, params...)
}

func (_DestinationRewardManager *DestinationRewardManagerCaller) GetAvailableRewardPoolIds(opts *bind.CallOpts, recipient common.Address, startIndex *big.Int, endIndex *big.Int) ([][32]byte, error) {
	var out []interface{}
	err := _DestinationRewardManager.contract.Call(opts, &out, "getAvailableRewardPoolIds", recipient, startIndex, endIndex)

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

func (_DestinationRewardManager *DestinationRewardManagerSession) GetAvailableRewardPoolIds(recipient common.Address, startIndex *big.Int, endIndex *big.Int) ([][32]byte, error) {
	return _DestinationRewardManager.Contract.GetAvailableRewardPoolIds(&_DestinationRewardManager.CallOpts, recipient, startIndex, endIndex)
}

func (_DestinationRewardManager *DestinationRewardManagerCallerSession) GetAvailableRewardPoolIds(recipient common.Address, startIndex *big.Int, endIndex *big.Int) ([][32]byte, error) {
	return _DestinationRewardManager.Contract.GetAvailableRewardPoolIds(&_DestinationRewardManager.CallOpts, recipient, startIndex, endIndex)
}

func (_DestinationRewardManager *DestinationRewardManagerCaller) ILinkAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DestinationRewardManager.contract.Call(opts, &out, "i_linkAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_DestinationRewardManager *DestinationRewardManagerSession) ILinkAddress() (common.Address, error) {
	return _DestinationRewardManager.Contract.ILinkAddress(&_DestinationRewardManager.CallOpts)
}

func (_DestinationRewardManager *DestinationRewardManagerCallerSession) ILinkAddress() (common.Address, error) {
	return _DestinationRewardManager.Contract.ILinkAddress(&_DestinationRewardManager.CallOpts)
}

func (_DestinationRewardManager *DestinationRewardManagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DestinationRewardManager.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_DestinationRewardManager *DestinationRewardManagerSession) Owner() (common.Address, error) {
	return _DestinationRewardManager.Contract.Owner(&_DestinationRewardManager.CallOpts)
}

func (_DestinationRewardManager *DestinationRewardManagerCallerSession) Owner() (common.Address, error) {
	return _DestinationRewardManager.Contract.Owner(&_DestinationRewardManager.CallOpts)
}

func (_DestinationRewardManager *DestinationRewardManagerCaller) SFeeManagerAddressList(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var out []interface{}
	err := _DestinationRewardManager.contract.Call(opts, &out, "s_feeManagerAddressList", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_DestinationRewardManager *DestinationRewardManagerSession) SFeeManagerAddressList(arg0 common.Address) (common.Address, error) {
	return _DestinationRewardManager.Contract.SFeeManagerAddressList(&_DestinationRewardManager.CallOpts, arg0)
}

func (_DestinationRewardManager *DestinationRewardManagerCallerSession) SFeeManagerAddressList(arg0 common.Address) (common.Address, error) {
	return _DestinationRewardManager.Contract.SFeeManagerAddressList(&_DestinationRewardManager.CallOpts, arg0)
}

func (_DestinationRewardManager *DestinationRewardManagerCaller) SRegisteredPoolIds(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _DestinationRewardManager.contract.Call(opts, &out, "s_registeredPoolIds", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_DestinationRewardManager *DestinationRewardManagerSession) SRegisteredPoolIds(arg0 *big.Int) ([32]byte, error) {
	return _DestinationRewardManager.Contract.SRegisteredPoolIds(&_DestinationRewardManager.CallOpts, arg0)
}

func (_DestinationRewardManager *DestinationRewardManagerCallerSession) SRegisteredPoolIds(arg0 *big.Int) ([32]byte, error) {
	return _DestinationRewardManager.Contract.SRegisteredPoolIds(&_DestinationRewardManager.CallOpts, arg0)
}

func (_DestinationRewardManager *DestinationRewardManagerCaller) SRewardRecipientWeights(opts *bind.CallOpts, arg0 [32]byte, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _DestinationRewardManager.contract.Call(opts, &out, "s_rewardRecipientWeights", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_DestinationRewardManager *DestinationRewardManagerSession) SRewardRecipientWeights(arg0 [32]byte, arg1 common.Address) (*big.Int, error) {
	return _DestinationRewardManager.Contract.SRewardRecipientWeights(&_DestinationRewardManager.CallOpts, arg0, arg1)
}

func (_DestinationRewardManager *DestinationRewardManagerCallerSession) SRewardRecipientWeights(arg0 [32]byte, arg1 common.Address) (*big.Int, error) {
	return _DestinationRewardManager.Contract.SRewardRecipientWeights(&_DestinationRewardManager.CallOpts, arg0, arg1)
}

func (_DestinationRewardManager *DestinationRewardManagerCaller) SRewardRecipientWeightsSet(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _DestinationRewardManager.contract.Call(opts, &out, "s_rewardRecipientWeightsSet", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_DestinationRewardManager *DestinationRewardManagerSession) SRewardRecipientWeightsSet(arg0 [32]byte) (bool, error) {
	return _DestinationRewardManager.Contract.SRewardRecipientWeightsSet(&_DestinationRewardManager.CallOpts, arg0)
}

func (_DestinationRewardManager *DestinationRewardManagerCallerSession) SRewardRecipientWeightsSet(arg0 [32]byte) (bool, error) {
	return _DestinationRewardManager.Contract.SRewardRecipientWeightsSet(&_DestinationRewardManager.CallOpts, arg0)
}

func (_DestinationRewardManager *DestinationRewardManagerCaller) STotalRewardRecipientFees(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _DestinationRewardManager.contract.Call(opts, &out, "s_totalRewardRecipientFees", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_DestinationRewardManager *DestinationRewardManagerSession) STotalRewardRecipientFees(arg0 [32]byte) (*big.Int, error) {
	return _DestinationRewardManager.Contract.STotalRewardRecipientFees(&_DestinationRewardManager.CallOpts, arg0)
}

func (_DestinationRewardManager *DestinationRewardManagerCallerSession) STotalRewardRecipientFees(arg0 [32]byte) (*big.Int, error) {
	return _DestinationRewardManager.Contract.STotalRewardRecipientFees(&_DestinationRewardManager.CallOpts, arg0)
}

func (_DestinationRewardManager *DestinationRewardManagerCaller) STotalRewardRecipientFeesLastClaimedAmounts(opts *bind.CallOpts, arg0 [32]byte, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _DestinationRewardManager.contract.Call(opts, &out, "s_totalRewardRecipientFeesLastClaimedAmounts", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_DestinationRewardManager *DestinationRewardManagerSession) STotalRewardRecipientFeesLastClaimedAmounts(arg0 [32]byte, arg1 common.Address) (*big.Int, error) {
	return _DestinationRewardManager.Contract.STotalRewardRecipientFeesLastClaimedAmounts(&_DestinationRewardManager.CallOpts, arg0, arg1)
}

func (_DestinationRewardManager *DestinationRewardManagerCallerSession) STotalRewardRecipientFeesLastClaimedAmounts(arg0 [32]byte, arg1 common.Address) (*big.Int, error) {
	return _DestinationRewardManager.Contract.STotalRewardRecipientFeesLastClaimedAmounts(&_DestinationRewardManager.CallOpts, arg0, arg1)
}

func (_DestinationRewardManager *DestinationRewardManagerCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _DestinationRewardManager.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_DestinationRewardManager *DestinationRewardManagerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _DestinationRewardManager.Contract.SupportsInterface(&_DestinationRewardManager.CallOpts, interfaceId)
}

func (_DestinationRewardManager *DestinationRewardManagerCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _DestinationRewardManager.Contract.SupportsInterface(&_DestinationRewardManager.CallOpts, interfaceId)
}

func (_DestinationRewardManager *DestinationRewardManagerCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _DestinationRewardManager.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_DestinationRewardManager *DestinationRewardManagerSession) TypeAndVersion() (string, error) {
	return _DestinationRewardManager.Contract.TypeAndVersion(&_DestinationRewardManager.CallOpts)
}

func (_DestinationRewardManager *DestinationRewardManagerCallerSession) TypeAndVersion() (string, error) {
	return _DestinationRewardManager.Contract.TypeAndVersion(&_DestinationRewardManager.CallOpts)
}

func (_DestinationRewardManager *DestinationRewardManagerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DestinationRewardManager.contract.Transact(opts, "acceptOwnership")
}

func (_DestinationRewardManager *DestinationRewardManagerSession) AcceptOwnership() (*types.Transaction, error) {
	return _DestinationRewardManager.Contract.AcceptOwnership(&_DestinationRewardManager.TransactOpts)
}

func (_DestinationRewardManager *DestinationRewardManagerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _DestinationRewardManager.Contract.AcceptOwnership(&_DestinationRewardManager.TransactOpts)
}

func (_DestinationRewardManager *DestinationRewardManagerTransactor) AddFeeManager(opts *bind.TransactOpts, newFeeManagerAddress common.Address) (*types.Transaction, error) {
	return _DestinationRewardManager.contract.Transact(opts, "addFeeManager", newFeeManagerAddress)
}

func (_DestinationRewardManager *DestinationRewardManagerSession) AddFeeManager(newFeeManagerAddress common.Address) (*types.Transaction, error) {
	return _DestinationRewardManager.Contract.AddFeeManager(&_DestinationRewardManager.TransactOpts, newFeeManagerAddress)
}

func (_DestinationRewardManager *DestinationRewardManagerTransactorSession) AddFeeManager(newFeeManagerAddress common.Address) (*types.Transaction, error) {
	return _DestinationRewardManager.Contract.AddFeeManager(&_DestinationRewardManager.TransactOpts, newFeeManagerAddress)
}

func (_DestinationRewardManager *DestinationRewardManagerTransactor) ClaimRewards(opts *bind.TransactOpts, poolIds [][32]byte) (*types.Transaction, error) {
	return _DestinationRewardManager.contract.Transact(opts, "claimRewards", poolIds)
}

func (_DestinationRewardManager *DestinationRewardManagerSession) ClaimRewards(poolIds [][32]byte) (*types.Transaction, error) {
	return _DestinationRewardManager.Contract.ClaimRewards(&_DestinationRewardManager.TransactOpts, poolIds)
}

func (_DestinationRewardManager *DestinationRewardManagerTransactorSession) ClaimRewards(poolIds [][32]byte) (*types.Transaction, error) {
	return _DestinationRewardManager.Contract.ClaimRewards(&_DestinationRewardManager.TransactOpts, poolIds)
}

func (_DestinationRewardManager *DestinationRewardManagerTransactor) OnFeePaid(opts *bind.TransactOpts, payments []IDestinationRewardManagerFeePayment, payer common.Address) (*types.Transaction, error) {
	return _DestinationRewardManager.contract.Transact(opts, "onFeePaid", payments, payer)
}

func (_DestinationRewardManager *DestinationRewardManagerSession) OnFeePaid(payments []IDestinationRewardManagerFeePayment, payer common.Address) (*types.Transaction, error) {
	return _DestinationRewardManager.Contract.OnFeePaid(&_DestinationRewardManager.TransactOpts, payments, payer)
}

func (_DestinationRewardManager *DestinationRewardManagerTransactorSession) OnFeePaid(payments []IDestinationRewardManagerFeePayment, payer common.Address) (*types.Transaction, error) {
	return _DestinationRewardManager.Contract.OnFeePaid(&_DestinationRewardManager.TransactOpts, payments, payer)
}

func (_DestinationRewardManager *DestinationRewardManagerTransactor) PayRecipients(opts *bind.TransactOpts, poolId [32]byte, recipients []common.Address) (*types.Transaction, error) {
	return _DestinationRewardManager.contract.Transact(opts, "payRecipients", poolId, recipients)
}

func (_DestinationRewardManager *DestinationRewardManagerSession) PayRecipients(poolId [32]byte, recipients []common.Address) (*types.Transaction, error) {
	return _DestinationRewardManager.Contract.PayRecipients(&_DestinationRewardManager.TransactOpts, poolId, recipients)
}

func (_DestinationRewardManager *DestinationRewardManagerTransactorSession) PayRecipients(poolId [32]byte, recipients []common.Address) (*types.Transaction, error) {
	return _DestinationRewardManager.Contract.PayRecipients(&_DestinationRewardManager.TransactOpts, poolId, recipients)
}

func (_DestinationRewardManager *DestinationRewardManagerTransactor) RemoveFeeManager(opts *bind.TransactOpts, feeManagerAddress common.Address) (*types.Transaction, error) {
	return _DestinationRewardManager.contract.Transact(opts, "removeFeeManager", feeManagerAddress)
}

func (_DestinationRewardManager *DestinationRewardManagerSession) RemoveFeeManager(feeManagerAddress common.Address) (*types.Transaction, error) {
	return _DestinationRewardManager.Contract.RemoveFeeManager(&_DestinationRewardManager.TransactOpts, feeManagerAddress)
}

func (_DestinationRewardManager *DestinationRewardManagerTransactorSession) RemoveFeeManager(feeManagerAddress common.Address) (*types.Transaction, error) {
	return _DestinationRewardManager.Contract.RemoveFeeManager(&_DestinationRewardManager.TransactOpts, feeManagerAddress)
}

func (_DestinationRewardManager *DestinationRewardManagerTransactor) SetRewardRecipients(opts *bind.TransactOpts, poolId [32]byte, rewardRecipientAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _DestinationRewardManager.contract.Transact(opts, "setRewardRecipients", poolId, rewardRecipientAndWeights)
}

func (_DestinationRewardManager *DestinationRewardManagerSession) SetRewardRecipients(poolId [32]byte, rewardRecipientAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _DestinationRewardManager.Contract.SetRewardRecipients(&_DestinationRewardManager.TransactOpts, poolId, rewardRecipientAndWeights)
}

func (_DestinationRewardManager *DestinationRewardManagerTransactorSession) SetRewardRecipients(poolId [32]byte, rewardRecipientAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _DestinationRewardManager.Contract.SetRewardRecipients(&_DestinationRewardManager.TransactOpts, poolId, rewardRecipientAndWeights)
}

func (_DestinationRewardManager *DestinationRewardManagerTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _DestinationRewardManager.contract.Transact(opts, "transferOwnership", to)
}

func (_DestinationRewardManager *DestinationRewardManagerSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _DestinationRewardManager.Contract.TransferOwnership(&_DestinationRewardManager.TransactOpts, to)
}

func (_DestinationRewardManager *DestinationRewardManagerTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _DestinationRewardManager.Contract.TransferOwnership(&_DestinationRewardManager.TransactOpts, to)
}

func (_DestinationRewardManager *DestinationRewardManagerTransactor) UpdateRewardRecipients(opts *bind.TransactOpts, poolId [32]byte, newRewardRecipients []CommonAddressAndWeight) (*types.Transaction, error) {
	return _DestinationRewardManager.contract.Transact(opts, "updateRewardRecipients", poolId, newRewardRecipients)
}

func (_DestinationRewardManager *DestinationRewardManagerSession) UpdateRewardRecipients(poolId [32]byte, newRewardRecipients []CommonAddressAndWeight) (*types.Transaction, error) {
	return _DestinationRewardManager.Contract.UpdateRewardRecipients(&_DestinationRewardManager.TransactOpts, poolId, newRewardRecipients)
}

func (_DestinationRewardManager *DestinationRewardManagerTransactorSession) UpdateRewardRecipients(poolId [32]byte, newRewardRecipients []CommonAddressAndWeight) (*types.Transaction, error) {
	return _DestinationRewardManager.Contract.UpdateRewardRecipients(&_DestinationRewardManager.TransactOpts, poolId, newRewardRecipients)
}

type DestinationRewardManagerFeeManagerUpdatedIterator struct {
	Event *DestinationRewardManagerFeeManagerUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DestinationRewardManagerFeeManagerUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DestinationRewardManagerFeeManagerUpdated)
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
		it.Event = new(DestinationRewardManagerFeeManagerUpdated)
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

func (it *DestinationRewardManagerFeeManagerUpdatedIterator) Error() error {
	return it.fail
}

func (it *DestinationRewardManagerFeeManagerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DestinationRewardManagerFeeManagerUpdated struct {
	NewFeeManagerAddress common.Address
	Raw                  types.Log
}

func (_DestinationRewardManager *DestinationRewardManagerFilterer) FilterFeeManagerUpdated(opts *bind.FilterOpts) (*DestinationRewardManagerFeeManagerUpdatedIterator, error) {

	logs, sub, err := _DestinationRewardManager.contract.FilterLogs(opts, "FeeManagerUpdated")
	if err != nil {
		return nil, err
	}
	return &DestinationRewardManagerFeeManagerUpdatedIterator{contract: _DestinationRewardManager.contract, event: "FeeManagerUpdated", logs: logs, sub: sub}, nil
}

func (_DestinationRewardManager *DestinationRewardManagerFilterer) WatchFeeManagerUpdated(opts *bind.WatchOpts, sink chan<- *DestinationRewardManagerFeeManagerUpdated) (event.Subscription, error) {

	logs, sub, err := _DestinationRewardManager.contract.WatchLogs(opts, "FeeManagerUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DestinationRewardManagerFeeManagerUpdated)
				if err := _DestinationRewardManager.contract.UnpackLog(event, "FeeManagerUpdated", log); err != nil {
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

func (_DestinationRewardManager *DestinationRewardManagerFilterer) ParseFeeManagerUpdated(log types.Log) (*DestinationRewardManagerFeeManagerUpdated, error) {
	event := new(DestinationRewardManagerFeeManagerUpdated)
	if err := _DestinationRewardManager.contract.UnpackLog(event, "FeeManagerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DestinationRewardManagerFeePaidIterator struct {
	Event *DestinationRewardManagerFeePaid

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DestinationRewardManagerFeePaidIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DestinationRewardManagerFeePaid)
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
		it.Event = new(DestinationRewardManagerFeePaid)
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

func (it *DestinationRewardManagerFeePaidIterator) Error() error {
	return it.fail
}

func (it *DestinationRewardManagerFeePaidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DestinationRewardManagerFeePaid struct {
	Payments []IDestinationRewardManagerFeePayment
	Payer    common.Address
	Raw      types.Log
}

func (_DestinationRewardManager *DestinationRewardManagerFilterer) FilterFeePaid(opts *bind.FilterOpts) (*DestinationRewardManagerFeePaidIterator, error) {

	logs, sub, err := _DestinationRewardManager.contract.FilterLogs(opts, "FeePaid")
	if err != nil {
		return nil, err
	}
	return &DestinationRewardManagerFeePaidIterator{contract: _DestinationRewardManager.contract, event: "FeePaid", logs: logs, sub: sub}, nil
}

func (_DestinationRewardManager *DestinationRewardManagerFilterer) WatchFeePaid(opts *bind.WatchOpts, sink chan<- *DestinationRewardManagerFeePaid) (event.Subscription, error) {

	logs, sub, err := _DestinationRewardManager.contract.WatchLogs(opts, "FeePaid")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DestinationRewardManagerFeePaid)
				if err := _DestinationRewardManager.contract.UnpackLog(event, "FeePaid", log); err != nil {
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

func (_DestinationRewardManager *DestinationRewardManagerFilterer) ParseFeePaid(log types.Log) (*DestinationRewardManagerFeePaid, error) {
	event := new(DestinationRewardManagerFeePaid)
	if err := _DestinationRewardManager.contract.UnpackLog(event, "FeePaid", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DestinationRewardManagerOwnershipTransferRequestedIterator struct {
	Event *DestinationRewardManagerOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DestinationRewardManagerOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DestinationRewardManagerOwnershipTransferRequested)
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
		it.Event = new(DestinationRewardManagerOwnershipTransferRequested)
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

func (it *DestinationRewardManagerOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *DestinationRewardManagerOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DestinationRewardManagerOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_DestinationRewardManager *DestinationRewardManagerFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*DestinationRewardManagerOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _DestinationRewardManager.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &DestinationRewardManagerOwnershipTransferRequestedIterator{contract: _DestinationRewardManager.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_DestinationRewardManager *DestinationRewardManagerFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *DestinationRewardManagerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _DestinationRewardManager.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DestinationRewardManagerOwnershipTransferRequested)
				if err := _DestinationRewardManager.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_DestinationRewardManager *DestinationRewardManagerFilterer) ParseOwnershipTransferRequested(log types.Log) (*DestinationRewardManagerOwnershipTransferRequested, error) {
	event := new(DestinationRewardManagerOwnershipTransferRequested)
	if err := _DestinationRewardManager.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DestinationRewardManagerOwnershipTransferredIterator struct {
	Event *DestinationRewardManagerOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DestinationRewardManagerOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DestinationRewardManagerOwnershipTransferred)
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
		it.Event = new(DestinationRewardManagerOwnershipTransferred)
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

func (it *DestinationRewardManagerOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *DestinationRewardManagerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DestinationRewardManagerOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_DestinationRewardManager *DestinationRewardManagerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*DestinationRewardManagerOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _DestinationRewardManager.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &DestinationRewardManagerOwnershipTransferredIterator{contract: _DestinationRewardManager.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_DestinationRewardManager *DestinationRewardManagerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DestinationRewardManagerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _DestinationRewardManager.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DestinationRewardManagerOwnershipTransferred)
				if err := _DestinationRewardManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_DestinationRewardManager *DestinationRewardManagerFilterer) ParseOwnershipTransferred(log types.Log) (*DestinationRewardManagerOwnershipTransferred, error) {
	event := new(DestinationRewardManagerOwnershipTransferred)
	if err := _DestinationRewardManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DestinationRewardManagerRewardRecipientsUpdatedIterator struct {
	Event *DestinationRewardManagerRewardRecipientsUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DestinationRewardManagerRewardRecipientsUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DestinationRewardManagerRewardRecipientsUpdated)
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
		it.Event = new(DestinationRewardManagerRewardRecipientsUpdated)
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

func (it *DestinationRewardManagerRewardRecipientsUpdatedIterator) Error() error {
	return it.fail
}

func (it *DestinationRewardManagerRewardRecipientsUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DestinationRewardManagerRewardRecipientsUpdated struct {
	PoolId              [32]byte
	NewRewardRecipients []CommonAddressAndWeight
	Raw                 types.Log
}

func (_DestinationRewardManager *DestinationRewardManagerFilterer) FilterRewardRecipientsUpdated(opts *bind.FilterOpts, poolId [][32]byte) (*DestinationRewardManagerRewardRecipientsUpdatedIterator, error) {

	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}

	logs, sub, err := _DestinationRewardManager.contract.FilterLogs(opts, "RewardRecipientsUpdated", poolIdRule)
	if err != nil {
		return nil, err
	}
	return &DestinationRewardManagerRewardRecipientsUpdatedIterator{contract: _DestinationRewardManager.contract, event: "RewardRecipientsUpdated", logs: logs, sub: sub}, nil
}

func (_DestinationRewardManager *DestinationRewardManagerFilterer) WatchRewardRecipientsUpdated(opts *bind.WatchOpts, sink chan<- *DestinationRewardManagerRewardRecipientsUpdated, poolId [][32]byte) (event.Subscription, error) {

	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}

	logs, sub, err := _DestinationRewardManager.contract.WatchLogs(opts, "RewardRecipientsUpdated", poolIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DestinationRewardManagerRewardRecipientsUpdated)
				if err := _DestinationRewardManager.contract.UnpackLog(event, "RewardRecipientsUpdated", log); err != nil {
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

func (_DestinationRewardManager *DestinationRewardManagerFilterer) ParseRewardRecipientsUpdated(log types.Log) (*DestinationRewardManagerRewardRecipientsUpdated, error) {
	event := new(DestinationRewardManagerRewardRecipientsUpdated)
	if err := _DestinationRewardManager.contract.UnpackLog(event, "RewardRecipientsUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DestinationRewardManagerRewardsClaimedIterator struct {
	Event *DestinationRewardManagerRewardsClaimed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DestinationRewardManagerRewardsClaimedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DestinationRewardManagerRewardsClaimed)
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
		it.Event = new(DestinationRewardManagerRewardsClaimed)
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

func (it *DestinationRewardManagerRewardsClaimedIterator) Error() error {
	return it.fail
}

func (it *DestinationRewardManagerRewardsClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DestinationRewardManagerRewardsClaimed struct {
	PoolId    [32]byte
	Recipient common.Address
	Quantity  *big.Int
	Raw       types.Log
}

func (_DestinationRewardManager *DestinationRewardManagerFilterer) FilterRewardsClaimed(opts *bind.FilterOpts, poolId [][32]byte, recipient []common.Address) (*DestinationRewardManagerRewardsClaimedIterator, error) {

	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _DestinationRewardManager.contract.FilterLogs(opts, "RewardsClaimed", poolIdRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &DestinationRewardManagerRewardsClaimedIterator{contract: _DestinationRewardManager.contract, event: "RewardsClaimed", logs: logs, sub: sub}, nil
}

func (_DestinationRewardManager *DestinationRewardManagerFilterer) WatchRewardsClaimed(opts *bind.WatchOpts, sink chan<- *DestinationRewardManagerRewardsClaimed, poolId [][32]byte, recipient []common.Address) (event.Subscription, error) {

	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _DestinationRewardManager.contract.WatchLogs(opts, "RewardsClaimed", poolIdRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DestinationRewardManagerRewardsClaimed)
				if err := _DestinationRewardManager.contract.UnpackLog(event, "RewardsClaimed", log); err != nil {
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

func (_DestinationRewardManager *DestinationRewardManagerFilterer) ParseRewardsClaimed(log types.Log) (*DestinationRewardManagerRewardsClaimed, error) {
	event := new(DestinationRewardManagerRewardsClaimed)
	if err := _DestinationRewardManager.contract.UnpackLog(event, "RewardsClaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_DestinationRewardManager *DestinationRewardManager) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _DestinationRewardManager.abi.Events["FeeManagerUpdated"].ID:
		return _DestinationRewardManager.ParseFeeManagerUpdated(log)
	case _DestinationRewardManager.abi.Events["FeePaid"].ID:
		return _DestinationRewardManager.ParseFeePaid(log)
	case _DestinationRewardManager.abi.Events["OwnershipTransferRequested"].ID:
		return _DestinationRewardManager.ParseOwnershipTransferRequested(log)
	case _DestinationRewardManager.abi.Events["OwnershipTransferred"].ID:
		return _DestinationRewardManager.ParseOwnershipTransferred(log)
	case _DestinationRewardManager.abi.Events["RewardRecipientsUpdated"].ID:
		return _DestinationRewardManager.ParseRewardRecipientsUpdated(log)
	case _DestinationRewardManager.abi.Events["RewardsClaimed"].ID:
		return _DestinationRewardManager.ParseRewardsClaimed(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (DestinationRewardManagerFeeManagerUpdated) Topic() common.Hash {
	return common.HexToHash("0xe45f5e140399b0a7e12971ab020724b828fbed8ac408c420884dc7d1bbe506b4")
}

func (DestinationRewardManagerFeePaid) Topic() common.Hash {
	return common.HexToHash("0xa1cc025ea76bacce5d740ee4bc331899375dc2c5f2ab33933aaacbd9ba001b66")
}

func (DestinationRewardManagerOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (DestinationRewardManagerOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (DestinationRewardManagerRewardRecipientsUpdated) Topic() common.Hash {
	return common.HexToHash("0x8f668d6090683f98b3373a8b83d214da45737f7486cb7de554cc07b54e61cfe6")
}

func (DestinationRewardManagerRewardsClaimed) Topic() common.Hash {
	return common.HexToHash("0x989969655bc1d593922527fe85d71347bb8e12fa423cc71f362dd8ef7cb10ef2")
}

func (_DestinationRewardManager *DestinationRewardManager) Address() common.Address {
	return _DestinationRewardManager.address
}

type DestinationRewardManagerInterface interface {
	GetAvailableRewardPoolIds(opts *bind.CallOpts, recipient common.Address, startIndex *big.Int, endIndex *big.Int) ([][32]byte, error)

	ILinkAddress(opts *bind.CallOpts) (common.Address, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SFeeManagerAddressList(opts *bind.CallOpts, arg0 common.Address) (common.Address, error)

	SRegisteredPoolIds(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error)

	SRewardRecipientWeights(opts *bind.CallOpts, arg0 [32]byte, arg1 common.Address) (*big.Int, error)

	SRewardRecipientWeightsSet(opts *bind.CallOpts, arg0 [32]byte) (bool, error)

	STotalRewardRecipientFees(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error)

	STotalRewardRecipientFeesLastClaimedAmounts(opts *bind.CallOpts, arg0 [32]byte, arg1 common.Address) (*big.Int, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AddFeeManager(opts *bind.TransactOpts, newFeeManagerAddress common.Address) (*types.Transaction, error)

	ClaimRewards(opts *bind.TransactOpts, poolIds [][32]byte) (*types.Transaction, error)

	OnFeePaid(opts *bind.TransactOpts, payments []IDestinationRewardManagerFeePayment, payer common.Address) (*types.Transaction, error)

	PayRecipients(opts *bind.TransactOpts, poolId [32]byte, recipients []common.Address) (*types.Transaction, error)

	RemoveFeeManager(opts *bind.TransactOpts, feeManagerAddress common.Address) (*types.Transaction, error)

	SetRewardRecipients(opts *bind.TransactOpts, poolId [32]byte, rewardRecipientAndWeights []CommonAddressAndWeight) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	UpdateRewardRecipients(opts *bind.TransactOpts, poolId [32]byte, newRewardRecipients []CommonAddressAndWeight) (*types.Transaction, error)

	FilterFeeManagerUpdated(opts *bind.FilterOpts) (*DestinationRewardManagerFeeManagerUpdatedIterator, error)

	WatchFeeManagerUpdated(opts *bind.WatchOpts, sink chan<- *DestinationRewardManagerFeeManagerUpdated) (event.Subscription, error)

	ParseFeeManagerUpdated(log types.Log) (*DestinationRewardManagerFeeManagerUpdated, error)

	FilterFeePaid(opts *bind.FilterOpts) (*DestinationRewardManagerFeePaidIterator, error)

	WatchFeePaid(opts *bind.WatchOpts, sink chan<- *DestinationRewardManagerFeePaid) (event.Subscription, error)

	ParseFeePaid(log types.Log) (*DestinationRewardManagerFeePaid, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*DestinationRewardManagerOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *DestinationRewardManagerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*DestinationRewardManagerOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*DestinationRewardManagerOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DestinationRewardManagerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*DestinationRewardManagerOwnershipTransferred, error)

	FilterRewardRecipientsUpdated(opts *bind.FilterOpts, poolId [][32]byte) (*DestinationRewardManagerRewardRecipientsUpdatedIterator, error)

	WatchRewardRecipientsUpdated(opts *bind.WatchOpts, sink chan<- *DestinationRewardManagerRewardRecipientsUpdated, poolId [][32]byte) (event.Subscription, error)

	ParseRewardRecipientsUpdated(log types.Log) (*DestinationRewardManagerRewardRecipientsUpdated, error)

	FilterRewardsClaimed(opts *bind.FilterOpts, poolId [][32]byte, recipient []common.Address) (*DestinationRewardManagerRewardsClaimedIterator, error)

	WatchRewardsClaimed(opts *bind.WatchOpts, sink chan<- *DestinationRewardManagerRewardsClaimed, poolId [][32]byte, recipient []common.Address) (event.Subscription, error)

	ParseRewardsClaimed(log types.Log) (*DestinationRewardManagerRewardsClaimed, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
