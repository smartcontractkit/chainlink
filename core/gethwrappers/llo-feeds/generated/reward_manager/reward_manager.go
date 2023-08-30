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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"linkAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPoolId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidWeights\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Unauthorized\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newFeeManagerAddress\",\"type\":\"address\"}],\"name\":\"FeeManagerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"internalType\":\"uint192\",\"name\":\"amount\",\"type\":\"uint192\"}],\"indexed\":false,\"internalType\":\"structIRewardManager.FeePayment[]\",\"name\":\"payments\",\"type\":\"tuple[]\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"FeePaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"weight\",\"type\":\"uint64\"}],\"indexed\":false,\"internalType\":\"structCommon.AddressAndWeight[]\",\"name\":\"newRewardRecipients\",\"type\":\"tuple[]\"}],\"name\":\"RewardRecipientsUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint192\",\"name\":\"quantity\",\"type\":\"uint192\"}],\"name\":\"RewardsClaimed\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"poolIds\",\"type\":\"bytes32[]\"}],\"name\":\"claimRewards\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"getAvailableRewardPoolIds\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"internalType\":\"uint192\",\"name\":\"amount\",\"type\":\"uint192\"}],\"internalType\":\"structIRewardManager.FeePayment[]\",\"name\":\"payments\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"onFeePaid\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"internalType\":\"address[]\",\"name\":\"recipients\",\"type\":\"address[]\"}],\"name\":\"payRecipients\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_feeManagerAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_registeredPoolIds\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"s_rewardRecipientWeights\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"s_rewardRecipientWeightsSet\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"s_totalRewardRecipientFees\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"s_totalRewardRecipientFeesLastClaimedAmounts\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newFeeManagerAddress\",\"type\":\"address\"}],\"name\":\"setFeeManager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"weight\",\"type\":\"uint64\"}],\"internalType\":\"structCommon.AddressAndWeight[]\",\"name\":\"rewardRecipientAndWeights\",\"type\":\"tuple[]\"}],\"name\":\"setRewardRecipients\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"weight\",\"type\":\"uint64\"}],\"internalType\":\"structCommon.AddressAndWeight[]\",\"name\":\"newRewardRecipients\",\"type\":\"tuple[]\"}],\"name\":\"updateRewardRecipients\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b5060405162002080380380620020808339810160408190526200003491620001a6565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000fb565b5050506001600160a01b038116620000e95760405163e6c4247b60e01b815260040160405180910390fd5b6001600160a01b0316608052620001d8565b336001600160a01b03821603620001555760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600060208284031215620001b957600080fd5b81516001600160a01b0381168114620001d157600080fd5b9392505050565b608051611e85620001fb60003960008181610ca70152610ee20152611e856000f3fe608060405234801561001057600080fd5b50600436106101365760003560e01c806359256201116100b25780638ac85a5c11610081578063b0d9fa1911610066578063b0d9fa1914610377578063cd5f72921461038a578063f2fde38b1461039d57600080fd5b80638ac85a5c1461032e5780638da5cb5b1461035957600080fd5b806359256201146102b857806360122608146102db5780636992922f1461030657806379ba50971461032657600080fd5b8063276e766011610109578063472d35b9116100ee578063472d35b91461027f5780634944832f146102925780634d322084146102a557600080fd5b8063276e76601461020c57806339ee81e11461025157600080fd5b806301ffc9a71461013b5780630f3c34d1146101a557806314060f23146101ba578063181f5a77146101cd575b600080fd5b6101906101493660046117aa565b7fffffffff00000000000000000000000000000000000000000000000000000000167fb0d9fa19000000000000000000000000000000000000000000000000000000001490565b60405190151581526020015b60405180910390f35b6101b86101b336600461186a565b6103b0565b005b6101b86101c836600461195c565b6103be565b604080518082018252601381527f5265776172644d616e6167657220302e302e31000000000000000000000000006020820152905161019c91906119cc565b60075461022c9073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161019c565b61027161025f366004611a1d565b60026020526000908152604090205481565b60405190815260200161019c565b6101b861028d366004611a5f565b610574565b6101b86102a036600461195c565b610642565b6101b86102b3366004611a81565b6107c4565b6101906102c6366004611a1d565b60056020526000908152604090205460ff1681565b6102716102e9366004611b00565b600360209081526000928352604080842090915290825290205481565b610319610314366004611a5f565b610903565b60405161019c9190611b2c565b6101b8610a34565b61027161033c366004611b00565b600460209081526000928352604080842090915290825290205481565b60005473ffffffffffffffffffffffffffffffffffffffff1661022c565b6101b8610385366004611b70565b610b36565b610271610398366004611a1d565b610d10565b6101b86103ab366004611a5f565b610d31565b6103ba3382610d45565b5050565b60005473ffffffffffffffffffffffffffffffffffffffff1633148015906103fe575060075473ffffffffffffffffffffffffffffffffffffffff163314155b15610435576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000819003610470576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008381526005602052604090205460ff16156104b9576040517f0afa7ee800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6006805460018181019092557ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d3f01849055600084815260056020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169091179055610535838383670de0b6b3a7640000610f12565b827f8f668d6090683f98b3373a8b83d214da45737f7486cb7de554cc07b54e61cfe68383604051610567929190611bdc565b60405180910390a2505050565b61057c611123565b73ffffffffffffffffffffffffffffffffffffffff81166105c9576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600780547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83169081179091556040519081527fe45f5e140399b0a7e12971ab020724b828fbed8ac408c420884dc7d1bbe506b49060200160405180910390a150565b61064a611123565b60408051600180825281830190925260009160208083019080368337019050509050838160008151811061068057610680611c51565b6020026020010181815250506000805b838110156107765760008585838181106106ac576106ac611c51565b6106c29260206040909202019081019150611a5f565b600088815260046020908152604080832073ffffffffffffffffffffffffffffffffffffffff8516845290915281205491925081900361072e576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61075f87878581811061074357610743611c51565b6107599260206040909202019081019150611a5f565b86610d45565b5092909201915061076f81611caf565b9050610690565b5061078385858584610f12565b847f8f668d6090683f98b3373a8b83d214da45737f7486cb7de554cc07b54e61cfe685856040516107b5929190611bdc565b60405180910390a25050505050565b826107e460005473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161415801561083657506000818152600460209081526040808320338452909152902054155b1561086d576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805160018082528183019092526000916020808301908036833701905050905084816000815181106108a3576108a3611c51565b60200260200101818152505060005b838110156108fb576108ea8585838181106108cf576108cf611c51565b90506020020160208101906108e49190611a5f565b83610d45565b506108f481611caf565b90506108b2565b505050505050565b60065460609060008167ffffffffffffffff811115610924576109246117ec565b60405190808252806020026020018201604052801561094d578160200160208202803683370190505b5090506000805b83811015610a2a5760006006828154811061097157610971611c51565b600091825260208083209091015480835260048252604080842073ffffffffffffffffffffffffffffffffffffffff8c16855290925291205490915015610a19576000818152600260209081526040808320546003835281842073ffffffffffffffffffffffffffffffffffffffff8c168552909252909120548114610a175781858580600101965081518110610a0a57610a0a611c51565b6020026020010181815250505b505b50610a2381611caf565b9050610954565b5090949350505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610aba576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60005473ffffffffffffffffffffffffffffffffffffffff163314801590610b76575060075473ffffffffffffffffffffffffffffffffffffffff163314155b15610bad576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000805b83811015610c8c57848482818110610bcb57610bcb611c51565b9050604002016020016020810190610be39190611d0f565b77ffffffffffffffffffffffffffffffffffffffffffffffff1660026000878785818110610c1357610c13611c51565b6040908102929092013583525060208201929092520160002080549091019055848482818110610c4557610c45611c51565b9050604002016020016020810190610c5d9190611d0f565b77ffffffffffffffffffffffffffffffffffffffffffffffff168201915080610c8590611caf565b9050610bb1565b50610ccf73ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000168330846111a6565b7fa1cc025ea76bacce5d740ee4bc331899375dc2c5f2ab33933aaacbd9ba001b66848484604051610d0293929190611d2a565b60405180910390a150505050565b60068181548110610d2057600080fd5b600091825260209091200154905081565b610d39611123565b610d4281611288565b50565b60008060005b8351811015610ec1576000848281518110610d6857610d68611c51565b6020026020010151905060006002600083815260200190815260200160002054905080600003610d99575050610eb1565b600082815260036020908152604080832073ffffffffffffffffffffffffffffffffffffffff8b16808552908352818420548685526004845282852091855292528220549083039190670de0b6b3a764000090830204905080600003610e025750505050610eb1565b600084815260036020908152604080832073ffffffffffffffffffffffffffffffffffffffff8d168085529252909120849055885196820196899087908110610e4d57610e4d611c51565b60200260200101517f989969655bc1d593922527fe85d71347bb8e12fa423cc71f362dd8ef7cb10ef283604051610ea4919077ffffffffffffffffffffffffffffffffffffffffffffffff91909116815260200190565b60405180910390a3505050505b610eba81611caf565b9050610d4b565b508015610f0957610f0973ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016858361137d565b90505b92915050565b610f6d8383808060200260200160405190810160405280939291908181526020016000905b82821015610f6357610f5460408302860136819003810190611db1565b81526020019060010190610f37565b50505050506113d8565b15610fa4576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000805b838110156110e2576000858583818110610fc457610fc4611c51565b9050604002016020016020810190610fdc9190611e0c565b67ffffffffffffffff1690506000868684818110610ffc57610ffc611c51565b6110129260206040909202019081019150611a5f565b905073ffffffffffffffffffffffffffffffffffffffff8116611061576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8160000361109b576040517f84677ce800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600088815260046020908152604080832073ffffffffffffffffffffffffffffffffffffffff90941683529290522081905591909101906110db81611caf565b9050610fa8565b5081811461111c576040517f84677ce800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050505050565b60005473ffffffffffffffffffffffffffffffffffffffff1633146111a4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610ab1565b565b60405173ffffffffffffffffffffffffffffffffffffffff808516602483015283166044820152606481018290526112829085907f23b872dd00000000000000000000000000000000000000000000000000000000906084015b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009093169290921790915261148f565b50505050565b3373ffffffffffffffffffffffffffffffffffffffff821603611307576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610ab1565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60405173ffffffffffffffffffffffffffffffffffffffff83166024820152604481018290526113d39084907fa9059cbb0000000000000000000000000000000000000000000000000000000090606401611200565b505050565b6000805b82518110156114865760006113f2826001611e27565b90505b835181101561147d5783818151811061141057611410611c51565b60200260200101516000015173ffffffffffffffffffffffffffffffffffffffff1684838151811061144457611444611c51565b60200260200101516000015173ffffffffffffffffffffffffffffffffffffffff1603611475575060019392505050565b6001016113f5565b506001016113dc565b50600092915050565b60006114f1826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff1661159b9092919063ffffffff16565b8051909150156113d3578080602001905181019061150f9190611e3a565b6113d3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f742073756363656564000000000000000000000000000000000000000000006064820152608401610ab1565b60606115aa84846000856115b2565b949350505050565b606082471015611644576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f60448201527f722063616c6c00000000000000000000000000000000000000000000000000006064820152608401610ab1565b6000808673ffffffffffffffffffffffffffffffffffffffff16858760405161166d9190611e5c565b60006040518083038185875af1925050503d80600081146116aa576040519150601f19603f3d011682016040523d82523d6000602084013e6116af565b606091505b50915091506116c0878383876116cb565b979650505050505050565b6060831561176157825160000361175a5773ffffffffffffffffffffffffffffffffffffffff85163b61175a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401610ab1565b50816115aa565b6115aa83838151156117765781518083602001fd5b806040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610ab191906119cc565b6000602082840312156117bc57600080fd5b81357fffffffff0000000000000000000000000000000000000000000000000000000081168114610f0957600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715611862576118626117ec565b604052919050565b6000602080838503121561187d57600080fd5b823567ffffffffffffffff8082111561189557600080fd5b818501915085601f8301126118a957600080fd5b8135818111156118bb576118bb6117ec565b8060051b91506118cc84830161181b565b81815291830184019184810190888411156118e657600080fd5b938501935b83851015611904578435825293850193908501906118eb565b98975050505050505050565b60008083601f84011261192257600080fd5b50813567ffffffffffffffff81111561193a57600080fd5b6020830191508360208260061b850101111561195557600080fd5b9250929050565b60008060006040848603121561197157600080fd5b83359250602084013567ffffffffffffffff81111561198f57600080fd5b61199b86828701611910565b9497909650939450505050565b60005b838110156119c35781810151838201526020016119ab565b50506000910152565b60208152600082518060208401526119eb8160408501602087016119a8565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169190910160400192915050565b600060208284031215611a2f57600080fd5b5035919050565b803573ffffffffffffffffffffffffffffffffffffffff81168114611a5a57600080fd5b919050565b600060208284031215611a7157600080fd5b611a7a82611a36565b9392505050565b600080600060408486031215611a9657600080fd5b83359250602084013567ffffffffffffffff80821115611ab557600080fd5b818601915086601f830112611ac957600080fd5b813581811115611ad857600080fd5b8760208260051b8501011115611aed57600080fd5b6020830194508093505050509250925092565b60008060408385031215611b1357600080fd5b82359150611b2360208401611a36565b90509250929050565b6020808252825182820181905260009190848201906040850190845b81811015611b6457835183529284019291840191600101611b48565b50909695505050505050565b600080600060408486031215611b8557600080fd5b833567ffffffffffffffff811115611b9c57600080fd5b611ba886828701611910565b9094509250611bbb905060208501611a36565b90509250925092565b803567ffffffffffffffff81168114611a5a57600080fd5b6020808252818101839052600090604080840186845b87811015611c445773ffffffffffffffffffffffffffffffffffffffff611c1883611a36565b16835267ffffffffffffffff611c2f868401611bc4565b16838601529183019190830190600101611bf2565b5090979650505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203611ce057611ce0611c80565b5060010190565b803577ffffffffffffffffffffffffffffffffffffffffffffffff81168114611a5a57600080fd5b600060208284031215611d2157600080fd5b611a7a82611ce7565b60408082528181018490526000908560608401835b87811015611d865782358252602077ffffffffffffffffffffffffffffffffffffffffffffffff611d71828601611ce7565b16908301529183019190830190600101611d3f565b5080935050505073ffffffffffffffffffffffffffffffffffffffff83166020830152949350505050565b600060408284031215611dc357600080fd5b6040516040810181811067ffffffffffffffff82111715611de657611de66117ec565b604052611df283611a36565b8152611e0060208401611bc4565b60208201529392505050565b600060208284031215611e1e57600080fd5b611a7a82611bc4565b80820180821115610f0c57610f0c611c80565b600060208284031215611e4c57600080fd5b81518015158114610f0957600080fd5b60008251611e6e8184602087016119a8565b919091019291505056fea164736f6c6343000810000a",
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

func (_RewardManager *RewardManagerCaller) GetAvailableRewardPoolIds(opts *bind.CallOpts, recipient common.Address) ([][32]byte, error) {
	var out []interface{}
	err := _RewardManager.contract.Call(opts, &out, "getAvailableRewardPoolIds", recipient)

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

func (_RewardManager *RewardManagerSession) GetAvailableRewardPoolIds(recipient common.Address) ([][32]byte, error) {
	return _RewardManager.Contract.GetAvailableRewardPoolIds(&_RewardManager.CallOpts, recipient)
}

func (_RewardManager *RewardManagerCallerSession) GetAvailableRewardPoolIds(recipient common.Address) ([][32]byte, error) {
	return _RewardManager.Contract.GetAvailableRewardPoolIds(&_RewardManager.CallOpts, recipient)
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

func (_RewardManager *RewardManagerTransactor) OnFeePaid(opts *bind.TransactOpts, payments []IRewardManagerFeePayment, payee common.Address) (*types.Transaction, error) {
	return _RewardManager.contract.Transact(opts, "onFeePaid", payments, payee)
}

func (_RewardManager *RewardManagerSession) OnFeePaid(payments []IRewardManagerFeePayment, payee common.Address) (*types.Transaction, error) {
	return _RewardManager.Contract.OnFeePaid(&_RewardManager.TransactOpts, payments, payee)
}

func (_RewardManager *RewardManagerTransactorSession) OnFeePaid(payments []IRewardManagerFeePayment, payee common.Address) (*types.Transaction, error) {
	return _RewardManager.Contract.OnFeePaid(&_RewardManager.TransactOpts, payments, payee)
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
	Payee    common.Address
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
	GetAvailableRewardPoolIds(opts *bind.CallOpts, recipient common.Address) ([][32]byte, error)

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

	OnFeePaid(opts *bind.TransactOpts, payments []IRewardManagerFeePayment, payee common.Address) (*types.Transaction, error)

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
