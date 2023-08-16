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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"linkAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPoolId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidWeights\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Unauthorized\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newFeeManagerAddress\",\"type\":\"address\"}],\"name\":\"FeeManagerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"internalType\":\"uint192\",\"name\":\"amount\",\"type\":\"uint192\"}],\"indexed\":false,\"internalType\":\"structIRewardManager.FeePayment[]\",\"name\":\"payments\",\"type\":\"tuple[]\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"FeePaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"weight\",\"type\":\"uint64\"}],\"indexed\":false,\"internalType\":\"structCommon.AddressAndWeight[]\",\"name\":\"newRewardRecipients\",\"type\":\"tuple[]\"}],\"name\":\"RewardRecipientsUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint192\",\"name\":\"quantity\",\"type\":\"uint192\"}],\"name\":\"RewardsClaimed\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"poolIds\",\"type\":\"bytes32[]\"}],\"name\":\"claimRewards\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"getAvailableRewardPoolIds\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"internalType\":\"uint192\",\"name\":\"amount\",\"type\":\"uint192\"}],\"internalType\":\"structIRewardManager.FeePayment[]\",\"name\":\"payments\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"onFeePaid\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"internalType\":\"address[]\",\"name\":\"recipients\",\"type\":\"address[]\"}],\"name\":\"payRecipients\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_feeManagerAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_registeredPoolIds\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"s_rewardRecipientWeights\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"s_totalRewardRecipientFees\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newFeeManagerAddress\",\"type\":\"address\"}],\"name\":\"setFeeManager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"weight\",\"type\":\"uint64\"}],\"internalType\":\"structCommon.AddressAndWeight[]\",\"name\":\"rewardRecipientAndWeights\",\"type\":\"tuple[]\"}],\"name\":\"setRewardRecipients\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"weight\",\"type\":\"uint64\"}],\"internalType\":\"structCommon.AddressAndWeight[]\",\"name\":\"newRewardRecipients\",\"type\":\"tuple[]\"}],\"name\":\"updateRewardRecipients\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b5060405162001c1338038062001c138339810160408190526200003491620001a6565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000fb565b5050506001600160a01b038116620000e95760405163e6c4247b60e01b815260040160405180910390fd5b6001600160a01b0316608052620001d8565b336001600160a01b03821603620001555760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600060208284031215620001b957600080fd5b81516001600160a01b0381168114620001d157600080fd5b9392505050565b608051611a18620001fb60003960008181610c160152610eaa0152611a186000f3fe608060405234801561001057600080fd5b50600436106101005760003560e01c80634d322084116100975780638da5cb5b116100665780638da5cb5b146102d5578063b0d9fa19146102f3578063cd5f729214610306578063f2fde38b1461031957600080fd5b80634d3220841461026f5780636992922f1461028257806379ba5097146102a25780638ac85a5c146102aa57600080fd5b8063276e7660116100d3578063276e7660146101d657806339ee81e11461021b578063472d35b9146102495780634944832f1461025c57600080fd5b806301ffc9a7146101055780630f3c34d11461016f57806314060f2314610184578063181f5a7714610197575b600080fd5b61015a610113366004611362565b7fffffffff00000000000000000000000000000000000000000000000000000000167fb0d9fa19000000000000000000000000000000000000000000000000000000001490565b60405190151581526020015b60405180910390f35b61018261017d366004611422565b61032c565b005b610182610192366004611514565b61033a565b604080518082018252601381527f5265776172644d616e6167657220302e302e3100000000000000000000000000602082015290516101669190611560565b6007546101f69073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610166565b61023b6102293660046115cc565b60026020526000908152604090205481565b604051908152602001610166565b61018261025736600461160e565b6104f0565b61018261026a366004611514565b6105be565b61018261027d366004611630565b610738565b61029561029036600461160e565b61086f565b60405161016691906116af565b61018261096b565b61023b6102b83660046116f3565b600460209081526000928352604080842090915290825290205481565b60005473ffffffffffffffffffffffffffffffffffffffff166101f6565b61018261030136600461171f565b610a6d565b61023b6103143660046115cc565b610cc5565b61018261032736600461160e565b610ce6565b6103363382610cfa565b5050565b60005473ffffffffffffffffffffffffffffffffffffffff16331480159061037a575060075473ffffffffffffffffffffffffffffffffffffffff163314155b156103b1576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008190036103ec576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008381526005602052604090205460ff1615610435576040517f0afa7ee800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6006805460018181019092557ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d3f01849055600084815260056020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690911790556104b1838383670de0b6b3a7640000610f22565b827f8f668d6090683f98b3373a8b83d214da45737f7486cb7de554cc07b54e61cfe683836040516104e392919061178b565b60405180910390a2505050565b6104f8611133565b73ffffffffffffffffffffffffffffffffffffffff8116610545576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600780547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83169081179091556040519081527fe45f5e140399b0a7e12971ab020724b828fbed8ac408c420884dc7d1bbe506b49060200160405180910390a150565b6105c6611133565b6040805160018082528183019092526000916020808301908036833701905050905083816000815181106105fc576105fc611800565b6020026020010181815250506000805b838110156106ea57600085858381811061062857610628611800565b61063e926020604090920201908101915061160e565b600088815260046020908152604080832073ffffffffffffffffffffffffffffffffffffffff851684529091528120549192508190036106aa576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6106db8787858181106106bf576106bf611800565b6106d5926020604090920201908101915061160e565b86610cfa565b5092909201915060010161060c565b506106f785858584610f22565b847f8f668d6090683f98b3373a8b83d214da45737f7486cb7de554cc07b54e61cfe6858560405161072992919061178b565b60405180910390a25050505050565b8261075860005473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141580156107aa57506000818152600460209081526040808320338452909152902054155b156107e1576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60408051600180825281830190925260009160208083019080368337019050509050848160008151811061081757610817611800565b60200260200101818152505060005b838110156108675761085e85858381811061084357610843611800565b9050602002016020810190610858919061160e565b83610cfa565b50600101610826565b505050505050565b60065460609060008167ffffffffffffffff811115610890576108906113a4565b6040519080825280602002602001820160405280156108b9578160200160208202803683370190505b5090506000805b83811015610961576000600682815481106108dd576108dd611800565b600091825260208083209091015480835260048252604080842073ffffffffffffffffffffffffffffffffffffffff8c168552909252912054909150156109585760008181526002602052604090205415610958578084848151811061094557610945611800565b6020026020010181815250508260010192505b506001016108c0565b5090949350505050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146109f1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60005473ffffffffffffffffffffffffffffffffffffffff163314801590610aad575060075473ffffffffffffffffffffffffffffffffffffffff163314155b15610ae4576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000805b83811015610bc357848482818110610b0257610b02611800565b9050604002016020016020810190610b1a9190611857565b77ffffffffffffffffffffffffffffffffffffffffffffffff1660026000878785818110610b4a57610b4a611800565b6040908102929092013583525060208201929092520160002080549091019055848482818110610b7c57610b7c611800565b9050604002016020016020810190610b949190611857565b77ffffffffffffffffffffffffffffffffffffffffffffffff168201915080610bbc906118a1565b9050610ae8565b506040517f23b872dd00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8381166004830152306024830152604482018390527f000000000000000000000000000000000000000000000000000000000000000016906323b872dd906064016020604051808303816000875af1158015610c5f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c8391906118d9565b507fa1cc025ea76bacce5d740ee4bc331899375dc2c5f2ab33933aaacbd9ba001b66848484604051610cb7939291906118fb565b60405180910390a150505050565b60068181548110610cd557600080fd5b600091825260209091200154905081565b610cee611133565b610cf7816111b6565b50565b60008060005b8351811015610e57576000848281518110610d1d57610d1d611800565b6020908102919091018101516000818152600283526040808220546003855281832073ffffffffffffffffffffffffffffffffffffffff8c16808552908652828420548585526004875283852091855295529082205492945092830391670de0b6b3a76400009083020490819003610d985750505050610e47565b600084815260036020908152604080832073ffffffffffffffffffffffffffffffffffffffff8d168085529252909120849055885196820196899087908110610de357610de3611800565b60200260200101517f989969655bc1d593922527fe85d71347bb8e12fa423cc71f362dd8ef7cb10ef283604051610e3a919077ffffffffffffffffffffffffffffffffffffffffffffffff91909116815260200190565b60405180910390a3505050505b610e50816118a1565b9050610d00565b508015610f19576040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8581166004830152602482018390527f0000000000000000000000000000000000000000000000000000000000000000169063a9059cbb906044016020604051808303816000875af1158015610ef3573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610f1791906118d9565b505b90505b92915050565b610f7d8383808060200260200160405190810160405280939291908181526020016000905b82821015610f7357610f6460408302860136819003810190611982565b81526020019060010190610f47565b50505050506112ab565b15610fb4576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000805b838110156110f2576000858583818110610fd457610fd4611800565b9050604002016020016020810190610fec91906119dd565b67ffffffffffffffff169050600086868481811061100c5761100c611800565b611022926020604090920201908101915061160e565b905073ffffffffffffffffffffffffffffffffffffffff8116611071576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b816000036110ab576040517f84677ce800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600088815260046020908152604080832073ffffffffffffffffffffffffffffffffffffffff90941683529290522081905591909101906110eb816118a1565b9050610fb8565b5081811461112c576040517f84677ce800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050505050565b60005473ffffffffffffffffffffffffffffffffffffffff1633146111b4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016109e8565b565b3373ffffffffffffffffffffffffffffffffffffffff821603611235576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016109e8565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000805b82518110156113595760006112c58260016119f8565b90505b8351811015611350578381815181106112e3576112e3611800565b60200260200101516000015173ffffffffffffffffffffffffffffffffffffffff1684838151811061131757611317611800565b60200260200101516000015173ffffffffffffffffffffffffffffffffffffffff1603611348575060019392505050565b6001016112c8565b506001016112af565b50600092915050565b60006020828403121561137457600080fd5b81357fffffffff0000000000000000000000000000000000000000000000000000000081168114610f1957600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561141a5761141a6113a4565b604052919050565b6000602080838503121561143557600080fd5b823567ffffffffffffffff8082111561144d57600080fd5b818501915085601f83011261146157600080fd5b813581811115611473576114736113a4565b8060051b91506114848483016113d3565b818152918301840191848101908884111561149e57600080fd5b938501935b838510156114bc578435825293850193908501906114a3565b98975050505050505050565b60008083601f8401126114da57600080fd5b50813567ffffffffffffffff8111156114f257600080fd5b6020830191508360208260061b850101111561150d57600080fd5b9250929050565b60008060006040848603121561152957600080fd5b83359250602084013567ffffffffffffffff81111561154757600080fd5b611553868287016114c8565b9497909650939450505050565b600060208083528351808285015260005b8181101561158d57858101830151858201604001528201611571565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b6000602082840312156115de57600080fd5b5035919050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461160957600080fd5b919050565b60006020828403121561162057600080fd5b611629826115e5565b9392505050565b60008060006040848603121561164557600080fd5b83359250602084013567ffffffffffffffff8082111561166457600080fd5b818601915086601f83011261167857600080fd5b81358181111561168757600080fd5b8760208260051b850101111561169c57600080fd5b6020830194508093505050509250925092565b6020808252825182820181905260009190848201906040850190845b818110156116e7578351835292840192918401916001016116cb565b50909695505050505050565b6000806040838503121561170657600080fd5b82359150611716602084016115e5565b90509250929050565b60008060006040848603121561173457600080fd5b833567ffffffffffffffff81111561174b57600080fd5b611757868287016114c8565b909450925061176a9050602085016115e5565b90509250925092565b803567ffffffffffffffff8116811461160957600080fd5b6020808252818101839052600090604080840186845b878110156117f35773ffffffffffffffffffffffffffffffffffffffff6117c7836115e5565b16835267ffffffffffffffff6117de868401611773565b168386015291830191908301906001016117a1565b5090979650505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b803577ffffffffffffffffffffffffffffffffffffffffffffffff8116811461160957600080fd5b60006020828403121561186957600080fd5b6116298261182f565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036118d2576118d2611872565b5060010190565b6000602082840312156118eb57600080fd5b81518015158114610f1957600080fd5b60408082528181018490526000908560608401835b878110156119575782358252602077ffffffffffffffffffffffffffffffffffffffffffffffff61194282860161182f565b16908301529183019190830190600101611910565b5080935050505073ffffffffffffffffffffffffffffffffffffffff83166020830152949350505050565b60006040828403121561199457600080fd5b6040516040810181811067ffffffffffffffff821117156119b7576119b76113a4565b6040526119c3836115e5565b81526119d160208401611773565b60208201529392505050565b6000602082840312156119ef57600080fd5b61162982611773565b80820180821115610f1c57610f1c61187256fea164736f6c6343000810000a",
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

	STotalRewardRecipientFees(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error)

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
