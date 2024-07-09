// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package solidity_vrf_v08_verifier_wrapper

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

type VRFProof struct {
	Pk            [2]*big.Int
	Gamma         [2]*big.Int
	C             *big.Int
	S             *big.Int
	Seed          *big.Int
	UWitness      common.Address
	CGammaWitness [2]*big.Int
	SHashWitness  [2]*big.Int
	ZInv          *big.Int
}

var VRFV08TestHelperMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"p1\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"p2\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"invZ\",\"type\":\"uint256\"}],\"name\":\"affineECAdd_\",\"outputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"\",\"type\":\"uint256[2]\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"base\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"exponent\",\"type\":\"uint256\"}],\"name\":\"bigModExp_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"x\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"scalar\",\"type\":\"uint256\"},{\"internalType\":\"uint256[2]\",\"name\":\"q\",\"type\":\"uint256[2]\"}],\"name\":\"ecmulVerify_\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"b\",\"type\":\"bytes\"}],\"name\":\"fieldHash_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"hashToCurve_\",\"outputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"\",\"type\":\"uint256[2]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"p\",\"type\":\"uint256[2]\"}],\"name\":\"isOnCurve_\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"},{\"internalType\":\"uint256[2]\",\"name\":\"p1\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"cp1Witness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"},{\"internalType\":\"uint256[2]\",\"name\":\"p2\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"sp2Witness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"zInv\",\"type\":\"uint256\"}],\"name\":\"linearCombination_\",\"outputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"\",\"type\":\"uint256[2]\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"px\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"py\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"qx\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"qy\",\"type\":\"uint256\"}],\"name\":\"projectiveECAdd_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"gamma\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"uWitness\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"cGammaWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"sHashWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"zInv\",\"type\":\"uint256\"}],\"internalType\":\"structVRF.Proof\",\"name\":\"proof\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"}],\"name\":\"randomValueFromVRFProof_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"output\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"hash\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"gamma\",\"type\":\"uint256[2]\"},{\"internalType\":\"address\",\"name\":\"uWitness\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"v\",\"type\":\"uint256[2]\"}],\"name\":\"scalarFromCurvePoints_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"squareRoot_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"},{\"internalType\":\"uint256[2]\",\"name\":\"p\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"lcWitness\",\"type\":\"address\"}],\"name\":\"verifyLinearCombinationWithGenerator_\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"gamma\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"uWitness\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"cGammaWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"sHashWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"zInv\",\"type\":\"uint256\"}],\"name\":\"verifyVRFProof_\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"ySquared_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50611b51806100206000396000f3fe608060405234801561001057600080fd5b50600436106100ea5760003560e01c80639d6f03371161008c578063b481e26011610066578063b481e260146101fc578063ef3b10ec1461020f578063fd7e4af914610224578063fe54f2a21461023757600080fd5b80639d6f0337146101c3578063a5e9508f146101d6578063aa7b2fbb146101e957600080fd5b80637f8f50a8116100c85780637f8f50a81461014c5780638af046ea1461015f57806391d5f6911461017257806395e6ee921461019557600080fd5b8063244f896d146100ef57806335452450146101185780635de600421461012b575b600080fd5b6101026100fd3660046115ed565b61024a565b60405161010f9190611a21565b60405180910390f35b6101026101263660046116bd565b610265565b61013e6101393660046118ff565b610280565b60405190815260200161010f565b61013e61015a366004611583565b61028c565b61013e61016d366004611827565b6102a5565b6101856101803660046118b8565b6102b0565b604051901515815260200161010f565b6101a86101a3366004611921565b6102c7565b6040805193845260208401929092529082015260600161010f565b61013e6101d1366004611827565b6102e8565b61013e6101e43660046117f5565b6102f3565b6101856101f73660046116e8565b6102ff565b61013e61020a366004611726565b61030c565b61022261021d36600461162b565b610317565b005b610185610232366004611567565b610333565b610102610245366004611840565b61033e565b610252611456565b61025d848484610361565b949350505050565b61026d611456565b6102778383610495565b90505b92915050565b600061027783836104f9565b600061029b86868686866105ed565b9695505050505050565b600061027a8261064b565b60006102be85858585610685565b95945050505050565b60008060006102d887878787610828565b9250925092509450945094915050565b600061027a826109be565b60006102778383610a16565b600061025d848484610b2a565b600061027a82610cb7565b610328898989898989898989610d11565b505050505050505050565b600061027a82610fe8565b610346611456565b61035588888888888888611143565b98975050505050505050565b610369611456565b83516020808601518551918601516000938493849361038a93909190610828565b919450925090507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f858209600114610423576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f696e765a206d75737420626520696e7665727365206f66207a0000000000000060448201526064015b60405180910390fd5b60405180604001604052807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8061045c5761045c611ae6565b87860981526020017ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8785099052979650505050505050565b61049d611456565b6104ca600184846040516020016104b693929190611a00565b6040516020818303038152906040526112cb565b90505b6104d681610fe8565b61027a5780516040805160208101929092526104f291016104b6565b90506104cd565b600080610504611474565b6020808252818101819052604082015260608101859052608081018490527ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f60a0820152610550611492565b60208160c08460057ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa9250826105e3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f6269674d6f64457870206661696c757265210000000000000000000000000000604482015260640161041a565b5195945050505050565b60006002868686858760405160200161060b9695949392919061198e565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101209695505050505050565b600061027a82600261067e7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f6001611a6c565b901c6104f9565b600073ffffffffffffffffffffffffffffffffffffffff8216610704576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f626164207769746e657373000000000000000000000000000000000000000000604482015260640161041a565b60208401516000906001161561071b57601c61071e565b601b5b905060007ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd03641418587600060200201510986517ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141918203925060009190890987516040805160008082526020820180845287905260ff88169282019290925260608101929092526080820183905291925060019060a0016020604051602081039080840390855afa1580156107d5573d6000803e3d6000fd5b50506040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0015173ffffffffffffffffffffffffffffffffffffffff9081169088161495505050505050949350505050565b60008080600180827ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f897ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f038808905060007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f038a08905060006108d083838585611333565b90985090506108e188828e8861138b565b90985090506108f288828c8761138b565b909850905060006109058d878b8561138b565b909850905061091688828686611333565b909850905061092788828e8961138b565b90985090508181146109aa577ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f818a0998507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f82890997507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f81830996506109ae565b8196505b5050505050509450945094915050565b6000807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80848509840990507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f600782089392505050565b604080518082018252600091610ad69190859060029083908390808284376000920191909152505060408051808201825291508087019060029083908390808284376000920191909152505050608086013560a087013586610a7f6101008a0160e08b0161154c565b604080518082018252906101008c019060029083908390808284376000920191909152505060408051808201825291506101408d0190600290839083908082843760009201919091525050506101808c0135610d11565b600383604001604051602001610aed929190611a52565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101209392505050565b600082610b93576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f7a65726f207363616c6172000000000000000000000000000000000000000000604482015260640161041a565b83516020850151600090610ba990600290611aab565b15610bb557601c610bb8565b601b5b905060007ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd03641418387096040805160008082526020820180845281905260ff86169282019290925260608101869052608081018390529192509060019060a0016020604051602081039080840390855afa158015610c38573d6000803e3d6000fd5b505050602060405103519050600086604051602001610c57919061197c565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152919052805160209091012073ffffffffffffffffffffffffffffffffffffffff92831692169190911498975050505050505050565b805160208201205b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8110610d0c57604080516020808201939093528151808203840181529082019091528051910120610cbf565b919050565b610d1a89610fe8565b610d80576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f7075626c6963206b6579206973206e6f74206f6e206375727665000000000000604482015260640161041a565b610d8988610fe8565b610def576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f67616d6d61206973206e6f74206f6e2063757276650000000000000000000000604482015260640161041a565b610df883610fe8565b610e5e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f6347616d6d615769746e657373206973206e6f74206f6e206375727665000000604482015260640161041a565b610e6782610fe8565b610ecd576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f73486173685769746e657373206973206e6f74206f6e20637572766500000000604482015260640161041a565b610ed9878a8887610685565b610f3f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f6164647228632a706b2b732a6729213d5f755769746e65737300000000000000604482015260640161041a565b6000610f4b8a87610495565b90506000610f5e898b878b868989611143565b90506000610f6f838d8d8a866105ed565b9050808a14610fda576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f696e76616c69642070726f6f6600000000000000000000000000000000000000604482015260640161041a565b505050505050505050505050565b80516000907ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f11611075576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e76616c696420782d6f7264696e6174650000000000000000000000000000604482015260640161041a565b60208201517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f11611102576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e76616c696420792d6f7264696e6174650000000000000000000000000000604482015260640161041a565b60208201517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f90800961113c8360005b60200201516109be565b1492915050565b61114b611456565b825186517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f90819006910614156111de576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f706f696e747320696e2073756d206d7573742062652064697374696e63740000604482015260640161041a565b6111e9878988610b2a565b61124f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4669727374206d756c20636865636b206661696c656400000000000000000000604482015260640161041a565b61125a848685610b2a565b6112c0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f5365636f6e64206d756c20636865636b206661696c6564000000000000000000604482015260640161041a565b610355868484610361565b6112d3611456565b6112dc82610cb7565b81526112f16112ec826000611132565b61064b565b602082018190526002900660011415610d0c576020810180517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f039052919050565b6000807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8487097ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8487099097909650945050505050565b600080807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f878509905060007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f87877ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f030990507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8183087ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f86890990999098509650505050505050565b60405180604001604052806002906020820280368337509192915050565b6040518060c001604052806006906020820280368337509192915050565b60405180602001604052806001906020820280368337509192915050565b803573ffffffffffffffffffffffffffffffffffffffff81168114610d0c57600080fd5b600082601f8301126114e557600080fd5b6040516040810181811067ffffffffffffffff8211171561150857611508611b15565b806040525080838560408601111561151f57600080fd5b60005b6002811015611541578135835260209283019290910190600101611522565b509195945050505050565b60006020828403121561155e57600080fd5b610277826114b0565b60006040828403121561157957600080fd5b61027783836114d4565b6000806000806000610120868803121561159c57600080fd5b6115a687876114d4565b94506115b587604088016114d4565b93506115c487608088016114d4565b92506115d260c087016114b0565b91506115e18760e088016114d4565b90509295509295909350565b600080600060a0848603121561160257600080fd5b61160c85856114d4565b925061161b85604086016114d4565b9150608084013590509250925092565b60008060008060008060008060006101a08a8c03121561164a57600080fd5b6116548b8b6114d4565b98506116638b60408c016114d4565b975060808a0135965060a08a0135955060c08a0135945061168660e08b016114b0565b93506116968b6101008c016114d4565b92506116a68b6101408c016114d4565b91506101808a013590509295985092959850929598565b600080606083850312156116d057600080fd5b6116da84846114d4565b946040939093013593505050565b600080600060a084860312156116fd57600080fd5b61170785856114d4565b92506040840135915061171d85606086016114d4565b90509250925092565b60006020828403121561173857600080fd5b813567ffffffffffffffff8082111561175057600080fd5b818401915084601f83011261176457600080fd5b81358181111561177657611776611b15565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f011681019083821181831017156117bc576117bc611b15565b816040528281528760208487010111156117d557600080fd5b826020860160208301376000928101602001929092525095945050505050565b6000808284036101c081121561180a57600080fd5b6101a08082121561181a57600080fd5b9395938601359450505050565b60006020828403121561183957600080fd5b5035919050565b6000806000806000806000610160888a03121561185c57600080fd5b8735965061186d8960208a016114d4565b955061187c8960608a016114d4565b945060a088013593506118928960c08a016114d4565b92506118a2896101008a016114d4565b9150610140880135905092959891949750929550565b60008060008060a085870312156118ce57600080fd5b843593506118df86602087016114d4565b9250606085013591506118f4608086016114b0565b905092959194509250565b6000806040838503121561191257600080fd5b50508035926020909101359150565b6000806000806080858703121561193757600080fd5b5050823594602084013594506040840135936060013592509050565b8060005b6002811015611976578151845260209384019390910190600101611957565b50505050565b6119868183611953565b604001919050565b86815261199e6020820187611953565b6119ab6060820186611953565b6119b860a0820185611953565b6119c560e0820184611953565b60609190911b7fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166101208201526101340195945050505050565b838152611a106020820184611953565b606081019190915260800192915050565b60408101818360005b6002811015611a49578151835260209283019290910190600101611a2a565b50505092915050565b828152606081016040836020840137600081529392505050565b60008219821115611aa6577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b500190565b600082611ae1577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500690565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var VRFV08TestHelperABI = VRFV08TestHelperMetaData.ABI

var VRFV08TestHelperBin = VRFV08TestHelperMetaData.Bin

func DeployVRFV08TestHelper(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *VRFV08TestHelper, error) {
	parsed, err := VRFV08TestHelperMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFV08TestHelperBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFV08TestHelper{address: address, abi: *parsed, VRFV08TestHelperCaller: VRFV08TestHelperCaller{contract: contract}, VRFV08TestHelperTransactor: VRFV08TestHelperTransactor{contract: contract}, VRFV08TestHelperFilterer: VRFV08TestHelperFilterer{contract: contract}}, nil
}

type VRFV08TestHelper struct {
	address common.Address
	abi     abi.ABI
	VRFV08TestHelperCaller
	VRFV08TestHelperTransactor
	VRFV08TestHelperFilterer
}

type VRFV08TestHelperCaller struct {
	contract *bind.BoundContract
}

type VRFV08TestHelperTransactor struct {
	contract *bind.BoundContract
}

type VRFV08TestHelperFilterer struct {
	contract *bind.BoundContract
}

type VRFV08TestHelperSession struct {
	Contract     *VRFV08TestHelper
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFV08TestHelperCallerSession struct {
	Contract *VRFV08TestHelperCaller
	CallOpts bind.CallOpts
}

type VRFV08TestHelperTransactorSession struct {
	Contract     *VRFV08TestHelperTransactor
	TransactOpts bind.TransactOpts
}

type VRFV08TestHelperRaw struct {
	Contract *VRFV08TestHelper
}

type VRFV08TestHelperCallerRaw struct {
	Contract *VRFV08TestHelperCaller
}

type VRFV08TestHelperTransactorRaw struct {
	Contract *VRFV08TestHelperTransactor
}

func NewVRFV08TestHelper(address common.Address, backend bind.ContractBackend) (*VRFV08TestHelper, error) {
	abi, err := abi.JSON(strings.NewReader(VRFV08TestHelperABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFV08TestHelper(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFV08TestHelper{address: address, abi: abi, VRFV08TestHelperCaller: VRFV08TestHelperCaller{contract: contract}, VRFV08TestHelperTransactor: VRFV08TestHelperTransactor{contract: contract}, VRFV08TestHelperFilterer: VRFV08TestHelperFilterer{contract: contract}}, nil
}

func NewVRFV08TestHelperCaller(address common.Address, caller bind.ContractCaller) (*VRFV08TestHelperCaller, error) {
	contract, err := bindVRFV08TestHelper(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV08TestHelperCaller{contract: contract}, nil
}

func NewVRFV08TestHelperTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFV08TestHelperTransactor, error) {
	contract, err := bindVRFV08TestHelper(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV08TestHelperTransactor{contract: contract}, nil
}

func NewVRFV08TestHelperFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFV08TestHelperFilterer, error) {
	contract, err := bindVRFV08TestHelper(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFV08TestHelperFilterer{contract: contract}, nil
}

func bindVRFV08TestHelper(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFV08TestHelperMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFV08TestHelper *VRFV08TestHelperRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV08TestHelper.Contract.VRFV08TestHelperCaller.contract.Call(opts, result, method, params...)
}

func (_VRFV08TestHelper *VRFV08TestHelperRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV08TestHelper.Contract.VRFV08TestHelperTransactor.contract.Transfer(opts)
}

func (_VRFV08TestHelper *VRFV08TestHelperRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV08TestHelper.Contract.VRFV08TestHelperTransactor.contract.Transact(opts, method, params...)
}

func (_VRFV08TestHelper *VRFV08TestHelperCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV08TestHelper.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFV08TestHelper *VRFV08TestHelperTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV08TestHelper.Contract.contract.Transfer(opts)
}

func (_VRFV08TestHelper *VRFV08TestHelperTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV08TestHelper.Contract.contract.Transact(opts, method, params...)
}

func (_VRFV08TestHelper *VRFV08TestHelperCaller) AffineECAdd(opts *bind.CallOpts, p1 [2]*big.Int, p2 [2]*big.Int, invZ *big.Int) ([2]*big.Int, error) {
	var out []interface{}
	err := _VRFV08TestHelper.contract.Call(opts, &out, "affineECAdd_", p1, p2, invZ)

	if err != nil {
		return *new([2]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([2]*big.Int)).(*[2]*big.Int)

	return out0, err

}

func (_VRFV08TestHelper *VRFV08TestHelperSession) AffineECAdd(p1 [2]*big.Int, p2 [2]*big.Int, invZ *big.Int) ([2]*big.Int, error) {
	return _VRFV08TestHelper.Contract.AffineECAdd(&_VRFV08TestHelper.CallOpts, p1, p2, invZ)
}

func (_VRFV08TestHelper *VRFV08TestHelperCallerSession) AffineECAdd(p1 [2]*big.Int, p2 [2]*big.Int, invZ *big.Int) ([2]*big.Int, error) {
	return _VRFV08TestHelper.Contract.AffineECAdd(&_VRFV08TestHelper.CallOpts, p1, p2, invZ)
}

func (_VRFV08TestHelper *VRFV08TestHelperCaller) BigModExp(opts *bind.CallOpts, base *big.Int, exponent *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VRFV08TestHelper.contract.Call(opts, &out, "bigModExp_", base, exponent)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV08TestHelper *VRFV08TestHelperSession) BigModExp(base *big.Int, exponent *big.Int) (*big.Int, error) {
	return _VRFV08TestHelper.Contract.BigModExp(&_VRFV08TestHelper.CallOpts, base, exponent)
}

func (_VRFV08TestHelper *VRFV08TestHelperCallerSession) BigModExp(base *big.Int, exponent *big.Int) (*big.Int, error) {
	return _VRFV08TestHelper.Contract.BigModExp(&_VRFV08TestHelper.CallOpts, base, exponent)
}

func (_VRFV08TestHelper *VRFV08TestHelperCaller) EcmulVerify(opts *bind.CallOpts, x [2]*big.Int, scalar *big.Int, q [2]*big.Int) (bool, error) {
	var out []interface{}
	err := _VRFV08TestHelper.contract.Call(opts, &out, "ecmulVerify_", x, scalar, q)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_VRFV08TestHelper *VRFV08TestHelperSession) EcmulVerify(x [2]*big.Int, scalar *big.Int, q [2]*big.Int) (bool, error) {
	return _VRFV08TestHelper.Contract.EcmulVerify(&_VRFV08TestHelper.CallOpts, x, scalar, q)
}

func (_VRFV08TestHelper *VRFV08TestHelperCallerSession) EcmulVerify(x [2]*big.Int, scalar *big.Int, q [2]*big.Int) (bool, error) {
	return _VRFV08TestHelper.Contract.EcmulVerify(&_VRFV08TestHelper.CallOpts, x, scalar, q)
}

func (_VRFV08TestHelper *VRFV08TestHelperCaller) FieldHash(opts *bind.CallOpts, b []byte) (*big.Int, error) {
	var out []interface{}
	err := _VRFV08TestHelper.contract.Call(opts, &out, "fieldHash_", b)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV08TestHelper *VRFV08TestHelperSession) FieldHash(b []byte) (*big.Int, error) {
	return _VRFV08TestHelper.Contract.FieldHash(&_VRFV08TestHelper.CallOpts, b)
}

func (_VRFV08TestHelper *VRFV08TestHelperCallerSession) FieldHash(b []byte) (*big.Int, error) {
	return _VRFV08TestHelper.Contract.FieldHash(&_VRFV08TestHelper.CallOpts, b)
}

func (_VRFV08TestHelper *VRFV08TestHelperCaller) HashToCurve(opts *bind.CallOpts, pk [2]*big.Int, x *big.Int) ([2]*big.Int, error) {
	var out []interface{}
	err := _VRFV08TestHelper.contract.Call(opts, &out, "hashToCurve_", pk, x)

	if err != nil {
		return *new([2]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([2]*big.Int)).(*[2]*big.Int)

	return out0, err

}

func (_VRFV08TestHelper *VRFV08TestHelperSession) HashToCurve(pk [2]*big.Int, x *big.Int) ([2]*big.Int, error) {
	return _VRFV08TestHelper.Contract.HashToCurve(&_VRFV08TestHelper.CallOpts, pk, x)
}

func (_VRFV08TestHelper *VRFV08TestHelperCallerSession) HashToCurve(pk [2]*big.Int, x *big.Int) ([2]*big.Int, error) {
	return _VRFV08TestHelper.Contract.HashToCurve(&_VRFV08TestHelper.CallOpts, pk, x)
}

func (_VRFV08TestHelper *VRFV08TestHelperCaller) IsOnCurve(opts *bind.CallOpts, p [2]*big.Int) (bool, error) {
	var out []interface{}
	err := _VRFV08TestHelper.contract.Call(opts, &out, "isOnCurve_", p)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_VRFV08TestHelper *VRFV08TestHelperSession) IsOnCurve(p [2]*big.Int) (bool, error) {
	return _VRFV08TestHelper.Contract.IsOnCurve(&_VRFV08TestHelper.CallOpts, p)
}

func (_VRFV08TestHelper *VRFV08TestHelperCallerSession) IsOnCurve(p [2]*big.Int) (bool, error) {
	return _VRFV08TestHelper.Contract.IsOnCurve(&_VRFV08TestHelper.CallOpts, p)
}

func (_VRFV08TestHelper *VRFV08TestHelperCaller) LinearCombination(opts *bind.CallOpts, c *big.Int, p1 [2]*big.Int, cp1Witness [2]*big.Int, s *big.Int, p2 [2]*big.Int, sp2Witness [2]*big.Int, zInv *big.Int) ([2]*big.Int, error) {
	var out []interface{}
	err := _VRFV08TestHelper.contract.Call(opts, &out, "linearCombination_", c, p1, cp1Witness, s, p2, sp2Witness, zInv)

	if err != nil {
		return *new([2]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([2]*big.Int)).(*[2]*big.Int)

	return out0, err

}

func (_VRFV08TestHelper *VRFV08TestHelperSession) LinearCombination(c *big.Int, p1 [2]*big.Int, cp1Witness [2]*big.Int, s *big.Int, p2 [2]*big.Int, sp2Witness [2]*big.Int, zInv *big.Int) ([2]*big.Int, error) {
	return _VRFV08TestHelper.Contract.LinearCombination(&_VRFV08TestHelper.CallOpts, c, p1, cp1Witness, s, p2, sp2Witness, zInv)
}

func (_VRFV08TestHelper *VRFV08TestHelperCallerSession) LinearCombination(c *big.Int, p1 [2]*big.Int, cp1Witness [2]*big.Int, s *big.Int, p2 [2]*big.Int, sp2Witness [2]*big.Int, zInv *big.Int) ([2]*big.Int, error) {
	return _VRFV08TestHelper.Contract.LinearCombination(&_VRFV08TestHelper.CallOpts, c, p1, cp1Witness, s, p2, sp2Witness, zInv)
}

func (_VRFV08TestHelper *VRFV08TestHelperCaller) ProjectiveECAdd(opts *bind.CallOpts, px *big.Int, py *big.Int, qx *big.Int, qy *big.Int) (*big.Int, *big.Int, *big.Int, error) {
	var out []interface{}
	err := _VRFV08TestHelper.contract.Call(opts, &out, "projectiveECAdd_", px, py, qx, qy)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	out2 := *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return out0, out1, out2, err

}

func (_VRFV08TestHelper *VRFV08TestHelperSession) ProjectiveECAdd(px *big.Int, py *big.Int, qx *big.Int, qy *big.Int) (*big.Int, *big.Int, *big.Int, error) {
	return _VRFV08TestHelper.Contract.ProjectiveECAdd(&_VRFV08TestHelper.CallOpts, px, py, qx, qy)
}

func (_VRFV08TestHelper *VRFV08TestHelperCallerSession) ProjectiveECAdd(px *big.Int, py *big.Int, qx *big.Int, qy *big.Int) (*big.Int, *big.Int, *big.Int, error) {
	return _VRFV08TestHelper.Contract.ProjectiveECAdd(&_VRFV08TestHelper.CallOpts, px, py, qx, qy)
}

func (_VRFV08TestHelper *VRFV08TestHelperCaller) RandomValueFromVRFProof(opts *bind.CallOpts, proof VRFProof, seed *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VRFV08TestHelper.contract.Call(opts, &out, "randomValueFromVRFProof_", proof, seed)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV08TestHelper *VRFV08TestHelperSession) RandomValueFromVRFProof(proof VRFProof, seed *big.Int) (*big.Int, error) {
	return _VRFV08TestHelper.Contract.RandomValueFromVRFProof(&_VRFV08TestHelper.CallOpts, proof, seed)
}

func (_VRFV08TestHelper *VRFV08TestHelperCallerSession) RandomValueFromVRFProof(proof VRFProof, seed *big.Int) (*big.Int, error) {
	return _VRFV08TestHelper.Contract.RandomValueFromVRFProof(&_VRFV08TestHelper.CallOpts, proof, seed)
}

func (_VRFV08TestHelper *VRFV08TestHelperCaller) ScalarFromCurvePoints(opts *bind.CallOpts, hash [2]*big.Int, pk [2]*big.Int, gamma [2]*big.Int, uWitness common.Address, v [2]*big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VRFV08TestHelper.contract.Call(opts, &out, "scalarFromCurvePoints_", hash, pk, gamma, uWitness, v)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV08TestHelper *VRFV08TestHelperSession) ScalarFromCurvePoints(hash [2]*big.Int, pk [2]*big.Int, gamma [2]*big.Int, uWitness common.Address, v [2]*big.Int) (*big.Int, error) {
	return _VRFV08TestHelper.Contract.ScalarFromCurvePoints(&_VRFV08TestHelper.CallOpts, hash, pk, gamma, uWitness, v)
}

func (_VRFV08TestHelper *VRFV08TestHelperCallerSession) ScalarFromCurvePoints(hash [2]*big.Int, pk [2]*big.Int, gamma [2]*big.Int, uWitness common.Address, v [2]*big.Int) (*big.Int, error) {
	return _VRFV08TestHelper.Contract.ScalarFromCurvePoints(&_VRFV08TestHelper.CallOpts, hash, pk, gamma, uWitness, v)
}

func (_VRFV08TestHelper *VRFV08TestHelperCaller) SquareRoot(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VRFV08TestHelper.contract.Call(opts, &out, "squareRoot_", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV08TestHelper *VRFV08TestHelperSession) SquareRoot(x *big.Int) (*big.Int, error) {
	return _VRFV08TestHelper.Contract.SquareRoot(&_VRFV08TestHelper.CallOpts, x)
}

func (_VRFV08TestHelper *VRFV08TestHelperCallerSession) SquareRoot(x *big.Int) (*big.Int, error) {
	return _VRFV08TestHelper.Contract.SquareRoot(&_VRFV08TestHelper.CallOpts, x)
}

func (_VRFV08TestHelper *VRFV08TestHelperCaller) VerifyLinearCombinationWithGenerator(opts *bind.CallOpts, c *big.Int, p [2]*big.Int, s *big.Int, lcWitness common.Address) (bool, error) {
	var out []interface{}
	err := _VRFV08TestHelper.contract.Call(opts, &out, "verifyLinearCombinationWithGenerator_", c, p, s, lcWitness)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_VRFV08TestHelper *VRFV08TestHelperSession) VerifyLinearCombinationWithGenerator(c *big.Int, p [2]*big.Int, s *big.Int, lcWitness common.Address) (bool, error) {
	return _VRFV08TestHelper.Contract.VerifyLinearCombinationWithGenerator(&_VRFV08TestHelper.CallOpts, c, p, s, lcWitness)
}

func (_VRFV08TestHelper *VRFV08TestHelperCallerSession) VerifyLinearCombinationWithGenerator(c *big.Int, p [2]*big.Int, s *big.Int, lcWitness common.Address) (bool, error) {
	return _VRFV08TestHelper.Contract.VerifyLinearCombinationWithGenerator(&_VRFV08TestHelper.CallOpts, c, p, s, lcWitness)
}

func (_VRFV08TestHelper *VRFV08TestHelperCaller) VerifyVRFProof(opts *bind.CallOpts, pk [2]*big.Int, gamma [2]*big.Int, c *big.Int, s *big.Int, seed *big.Int, uWitness common.Address, cGammaWitness [2]*big.Int, sHashWitness [2]*big.Int, zInv *big.Int) error {
	var out []interface{}
	err := _VRFV08TestHelper.contract.Call(opts, &out, "verifyVRFProof_", pk, gamma, c, s, seed, uWitness, cGammaWitness, sHashWitness, zInv)

	if err != nil {
		return err
	}

	return err

}

func (_VRFV08TestHelper *VRFV08TestHelperSession) VerifyVRFProof(pk [2]*big.Int, gamma [2]*big.Int, c *big.Int, s *big.Int, seed *big.Int, uWitness common.Address, cGammaWitness [2]*big.Int, sHashWitness [2]*big.Int, zInv *big.Int) error {
	return _VRFV08TestHelper.Contract.VerifyVRFProof(&_VRFV08TestHelper.CallOpts, pk, gamma, c, s, seed, uWitness, cGammaWitness, sHashWitness, zInv)
}

func (_VRFV08TestHelper *VRFV08TestHelperCallerSession) VerifyVRFProof(pk [2]*big.Int, gamma [2]*big.Int, c *big.Int, s *big.Int, seed *big.Int, uWitness common.Address, cGammaWitness [2]*big.Int, sHashWitness [2]*big.Int, zInv *big.Int) error {
	return _VRFV08TestHelper.Contract.VerifyVRFProof(&_VRFV08TestHelper.CallOpts, pk, gamma, c, s, seed, uWitness, cGammaWitness, sHashWitness, zInv)
}

func (_VRFV08TestHelper *VRFV08TestHelperCaller) YSquared(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VRFV08TestHelper.contract.Call(opts, &out, "ySquared_", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV08TestHelper *VRFV08TestHelperSession) YSquared(x *big.Int) (*big.Int, error) {
	return _VRFV08TestHelper.Contract.YSquared(&_VRFV08TestHelper.CallOpts, x)
}

func (_VRFV08TestHelper *VRFV08TestHelperCallerSession) YSquared(x *big.Int) (*big.Int, error) {
	return _VRFV08TestHelper.Contract.YSquared(&_VRFV08TestHelper.CallOpts, x)
}

func (_VRFV08TestHelper *VRFV08TestHelper) Address() common.Address {
	return _VRFV08TestHelper.address
}

type VRFV08TestHelperInterface interface {
	AffineECAdd(opts *bind.CallOpts, p1 [2]*big.Int, p2 [2]*big.Int, invZ *big.Int) ([2]*big.Int, error)

	BigModExp(opts *bind.CallOpts, base *big.Int, exponent *big.Int) (*big.Int, error)

	EcmulVerify(opts *bind.CallOpts, x [2]*big.Int, scalar *big.Int, q [2]*big.Int) (bool, error)

	FieldHash(opts *bind.CallOpts, b []byte) (*big.Int, error)

	HashToCurve(opts *bind.CallOpts, pk [2]*big.Int, x *big.Int) ([2]*big.Int, error)

	IsOnCurve(opts *bind.CallOpts, p [2]*big.Int) (bool, error)

	LinearCombination(opts *bind.CallOpts, c *big.Int, p1 [2]*big.Int, cp1Witness [2]*big.Int, s *big.Int, p2 [2]*big.Int, sp2Witness [2]*big.Int, zInv *big.Int) ([2]*big.Int, error)

	ProjectiveECAdd(opts *bind.CallOpts, px *big.Int, py *big.Int, qx *big.Int, qy *big.Int) (*big.Int, *big.Int, *big.Int, error)

	RandomValueFromVRFProof(opts *bind.CallOpts, proof VRFProof, seed *big.Int) (*big.Int, error)

	ScalarFromCurvePoints(opts *bind.CallOpts, hash [2]*big.Int, pk [2]*big.Int, gamma [2]*big.Int, uWitness common.Address, v [2]*big.Int) (*big.Int, error)

	SquareRoot(opts *bind.CallOpts, x *big.Int) (*big.Int, error)

	VerifyLinearCombinationWithGenerator(opts *bind.CallOpts, c *big.Int, p [2]*big.Int, s *big.Int, lcWitness common.Address) (bool, error)

	VerifyVRFProof(opts *bind.CallOpts, pk [2]*big.Int, gamma [2]*big.Int, c *big.Int, s *big.Int, seed *big.Int, uWitness common.Address, cGammaWitness [2]*big.Int, sHashWitness [2]*big.Int, zInv *big.Int) error

	YSquared(opts *bind.CallOpts, x *big.Int) (*big.Int, error)

	Address() common.Address
}
