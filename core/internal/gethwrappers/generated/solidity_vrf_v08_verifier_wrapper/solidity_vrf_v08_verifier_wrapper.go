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
	ABI: "[{\"inputs\":[],\"name\":\"HASH_TO_CURVE_HASH_PREFIX\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PROOF_LENGTH\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SCALAR_FROM_CURVE_POINTS_HASH_PREFIX\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"VRF_RANDOM_OUTPUT_HASH_PREFIX\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"p1\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"p2\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"invZ\",\"type\":\"uint256\"}],\"name\":\"affineECAdd_\",\"outputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"\",\"type\":\"uint256[2]\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"base\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"exponent\",\"type\":\"uint256\"}],\"name\":\"bigModExp_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"x\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"scalar\",\"type\":\"uint256\"},{\"internalType\":\"uint256[2]\",\"name\":\"q\",\"type\":\"uint256[2]\"}],\"name\":\"ecmulVerify_\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"b\",\"type\":\"bytes\"}],\"name\":\"fieldHash_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"hashToCurve_\",\"outputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"\",\"type\":\"uint256[2]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"p\",\"type\":\"uint256[2]\"}],\"name\":\"isOnCurve_\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"},{\"internalType\":\"uint256[2]\",\"name\":\"p1\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"cp1Witness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"},{\"internalType\":\"uint256[2]\",\"name\":\"p2\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"sp2Witness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"zInv\",\"type\":\"uint256\"}],\"name\":\"linearCombination_\",\"outputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"\",\"type\":\"uint256[2]\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"px\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"py\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"qx\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"qy\",\"type\":\"uint256\"}],\"name\":\"projectiveECAdd_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"gamma\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"uWitness\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"cGammaWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"sHashWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"zInv\",\"type\":\"uint256\"}],\"internalType\":\"structVRF.Proof\",\"name\":\"proof\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"}],\"name\":\"randomValueFromVRFProof_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"output\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"hash\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"gamma\",\"type\":\"uint256[2]\"},{\"internalType\":\"address\",\"name\":\"uWitness\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"v\",\"type\":\"uint256[2]\"}],\"name\":\"scalarFromCurvePoints_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"squareRoot_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"},{\"internalType\":\"uint256[2]\",\"name\":\"p\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"lcWitness\",\"type\":\"address\"}],\"name\":\"verifyLinearCombinationWithGenerator_\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"gamma\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"uWitness\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"cGammaWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"sHashWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"zInv\",\"type\":\"uint256\"}],\"name\":\"verifyVRFProof_\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"ySquared_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50611b9d806100206000396000f3fe608060405234801561001057600080fd5b50600436106101365760003560e01c806395e6ee92116100b2578063b481e26011610081578063ef3b10ec11610066578063ef3b10ec1461027c578063fd7e4af914610291578063fe54f2a2146102a457600080fd5b8063b481e26014610260578063e911439c1461027357600080fd5b806395e6ee92146101f95780639d6f033714610227578063a5e9508f1461023a578063aa7b2fbb1461024d57600080fd5b806355c82a86116101095780637f8f50a8116100ee5780637f8f50a8146101b05780638af046ea146101c357806391d5f691146101d657600080fd5b806355c82a86146101955780635de600421461019d57600080fd5b8063244f896d1461013b57806327c5728314610164578063354524501461017a57806350e252d51461018d575b600080fd5b61014e6101493660046115b2565b6102b7565b60405161015b9190611a6c565b60405180910390f35b61016c600281565b60405190815260200161015b565b61014e610188366004611682565b6102d4565b61016c600381565b61016c600181565b61016c6101ab36600461194a565b6102ef565b61016c6101be366004611548565b6102fb565b61016c6101d1366004611872565b610314565b6101e96101e4366004611903565b61031f565b604051901515815260200161015b565b61020c61020736600461196c565b610336565b6040805193845260208401929092529082015260600161015b565b61016c610235366004611872565b610357565b61016c6102483660046117ba565b610362565b6101e961025b3660046116ad565b61036e565b61016c61026e3660046116eb565b61037b565b61016c6101a081565b61028f61028a3660046115f0565b610386565b005b6101e961029f36600461152c565b6103a2565b61014e6102b236600461188b565b6103ad565b6102bf611436565b6102ca8484846103d0565b90505b9392505050565b6102dc611436565b6102e68383610504565b90505b92915050565b60006102e68383610568565b600061030a868686868661065c565b9695505050505050565b60006102e9826106ba565b600061032d858585856106f4565b95945050505050565b600080600061034787878787610897565b9250925092509450945094915050565b60006102e982610a2d565b60006102e68383610a85565b60006102ca848484610b0e565b60006102e982610c9b565b610397898989898989898989610cf5565b505050505050505050565b60006102e982610fcc565b6103b5611436565b6103c488888888888888611127565b98975050505050505050565b6103d8611436565b8351602080860151855191860151600093849384936103f993909190610897565b919450925090507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f858209600114610492576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f696e765a206d75737420626520696e7665727365206f66207a0000000000000060448201526064015b60405180910390fd5b60405180604001604052807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f806104cb576104cb611b32565b87860981526020017ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8785099052979650505050505050565b61050c611436565b6105396001848460405160200161052593929190611a4b565b6040516020818303038152906040526112ab565b90505b61054581610fcc565b6102e95780516040805160208101929092526105619101610525565b905061053c565b600080610573611454565b6020808252818101819052604082015260608101859052608081018490527ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f60a08201526105bf611472565b60208160c08460057ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa925082610652576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f6269674d6f64457870206661696c7572652100000000000000000000000000006044820152606401610489565b5195945050505050565b60006002868686858760405160200161067a969594939291906119d9565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101209695505050505050565b60006102e98260026106ed7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f6001611ab8565b901c610568565b600073ffffffffffffffffffffffffffffffffffffffff8216610773576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f626164207769746e6573730000000000000000000000000000000000000000006044820152606401610489565b60208401516000906001161561078a57601c61078d565b601b5b905060007ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd03641418587600060200201510986517ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141918203925060009190890987516040805160008082526020820180845287905260ff88169282019290925260608101929092526080820183905291925060019060a0016020604051602081039080840390855afa158015610844573d6000803e3d6000fd5b50506040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0015173ffffffffffffffffffffffffffffffffffffffff9081169088161495505050505050949350505050565b60008080600180827ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f897ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f038808905060007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f038a089050600061093f83838585611313565b909850905061095088828e8861136b565b909850905061096188828c8761136b565b909850905060006109748d878b8561136b565b909850905061098588828686611313565b909850905061099688828e8961136b565b9098509050818114610a19577ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f818a0998507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f82890997507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8183099650610a1d565b8196505b5050505050509450945094915050565b6000807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80848509840990507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f600782089392505050565b6000610ab98360000151846020015185604001518660600151868860a001518960c001518a60e001518b6101000151610cf5565b60038360200151604051602001610ad1929190611a7a565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101209392505050565b600082610b77576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f7a65726f207363616c61720000000000000000000000000000000000000000006044820152606401610489565b83516020850151600090610b8d90600290611af7565b15610b9957601c610b9c565b601b5b905060007ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd03641418387096040805160008082526020820180845281905260ff86169282019290925260608101869052608081018390529192509060019060a0016020604051602081039080840390855afa158015610c1c573d6000803e3d6000fd5b505050602060405103519050600086604051602001610c3b91906119c7565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152919052805160209091012073ffffffffffffffffffffffffffffffffffffffff92831692169190911498975050505050505050565b805160208201205b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8110610cf057604080516020808201939093528151808203840181529082019091528051910120610ca3565b919050565b610cfe89610fcc565b610d64576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f7075626c6963206b6579206973206e6f74206f6e2063757276650000000000006044820152606401610489565b610d6d88610fcc565b610dd3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f67616d6d61206973206e6f74206f6e20637572766500000000000000000000006044820152606401610489565b610ddc83610fcc565b610e42576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f6347616d6d615769746e657373206973206e6f74206f6e2063757276650000006044820152606401610489565b610e4b82610fcc565b610eb1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f73486173685769746e657373206973206e6f74206f6e206375727665000000006044820152606401610489565b610ebd878a88876106f4565b610f23576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f6164647228632a706b2b732a6729213d5f755769746e657373000000000000006044820152606401610489565b6000610f2f8a87610504565b90506000610f42898b878b868989611127565b90506000610f53838d8d8a8661065c565b9050808a14610fbe576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f696e76616c69642070726f6f66000000000000000000000000000000000000006044820152606401610489565b505050505050505050505050565b80516000907ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f11611059576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e76616c696420782d6f7264696e61746500000000000000000000000000006044820152606401610489565b60208201517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f116110e6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e76616c696420792d6f7264696e61746500000000000000000000000000006044820152606401610489565b60208201517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f9080096111208360005b6020020151610a2d565b1492915050565b61112f611436565b825186517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f919003066111be576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f706f696e747320696e2073756d206d7573742062652064697374696e637400006044820152606401610489565b6111c9878988610b0e565b61122f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4669727374206d756c20636865636b206661696c6564000000000000000000006044820152606401610489565b61123a848685610b0e565b6112a0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f5365636f6e64206d756c20636865636b206661696c65640000000000000000006044820152606401610489565b6103c48684846103d0565b6112b3611436565b6112bc82610c9b565b81526112d16112cc826000611116565b6106ba565b602082018190526002900660011415610cf0576020810180517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f039052919050565b6000807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8487097ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8487099097909650945050505050565b600080807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f878509905060007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f87877ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f030990507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8183087ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f86890990999098509650505050505050565b60405180604001604052806002906020820280368337509192915050565b6040518060c001604052806006906020820280368337509192915050565b60405180602001604052806001906020820280368337509192915050565b803573ffffffffffffffffffffffffffffffffffffffff81168114610cf057600080fd5b600082601f8301126114c557600080fd5b6040516040810181811067ffffffffffffffff821117156114e8576114e8611b61565b80604052508083856040860111156114ff57600080fd5b60005b6002811015611521578135835260209283019290910190600101611502565b509195945050505050565b60006040828403121561153e57600080fd5b6102e683836114b4565b6000806000806000610120868803121561156157600080fd5b61156b87876114b4565b945061157a87604088016114b4565b935061158987608088016114b4565b925061159760c08701611490565b91506115a68760e088016114b4565b90509295509295909350565b600080600060a084860312156115c757600080fd5b6115d185856114b4565b92506115e085604086016114b4565b9150608084013590509250925092565b60008060008060008060008060006101a08a8c03121561160f57600080fd5b6116198b8b6114b4565b98506116288b60408c016114b4565b975060808a0135965060a08a0135955060c08a0135945061164b60e08b01611490565b935061165b8b6101008c016114b4565b925061166b8b6101408c016114b4565b91506101808a013590509295985092959850929598565b6000806060838503121561169557600080fd5b61169f84846114b4565b946040939093013593505050565b600080600060a084860312156116c257600080fd5b6116cc85856114b4565b9250604084013591506116e285606086016114b4565b90509250925092565b6000602082840312156116fd57600080fd5b813567ffffffffffffffff8082111561171557600080fd5b818401915084601f83011261172957600080fd5b81358181111561173b5761173b611b61565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f0116810190838211818310171561178157611781611b61565b8160405282815287602084870101111561179a57600080fd5b826020860160208301376000928101602001929092525095945050505050565b6000808284036101c08112156117cf57600080fd5b6101a0808212156117df57600080fd5b6117e7611a8e565b91506117f386866114b4565b825261180286604087016114b4565b60208301526080850135604083015260a0850135606083015260c0850135608083015261183160e08601611490565b60a0830152610100611845878288016114b4565b60c08401526118588761014088016114b4565b60e084015261018086013590830152909593013593505050565b60006020828403121561188457600080fd5b5035919050565b6000806000806000806000610160888a0312156118a757600080fd5b873596506118b88960208a016114b4565b95506118c78960608a016114b4565b945060a088013593506118dd8960c08a016114b4565b92506118ed896101008a016114b4565b9150610140880135905092959891949750929550565b60008060008060a0858703121561191957600080fd5b8435935061192a86602087016114b4565b92506060850135915061193f60808601611490565b905092959194509250565b6000806040838503121561195d57600080fd5b50508035926020909101359150565b6000806000806080858703121561198257600080fd5b5050823594602084013594506040840135936060013592509050565b8060005b60028110156119c15781518452602093840193909101906001016119a2565b50505050565b6119d1818361199e565b604001919050565b8681526119e9602082018761199e565b6119f6606082018661199e565b611a0360a082018561199e565b611a1060e082018461199e565b60609190911b7fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166101208201526101340195945050505050565b838152611a5b602082018461199e565b606081019190915260800192915050565b604081016102e9828461199e565b828152606081016102cd602083018461199e565b604051610120810167ffffffffffffffff81118282101715611ab257611ab2611b61565b60405290565b60008219821115611af2577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b500190565b600082611b2d577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500690565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
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
	return address, tx, &VRFV08TestHelper{VRFV08TestHelperCaller: VRFV08TestHelperCaller{contract: contract}, VRFV08TestHelperTransactor: VRFV08TestHelperTransactor{contract: contract}, VRFV08TestHelperFilterer: VRFV08TestHelperFilterer{contract: contract}}, nil
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
	parsed, err := abi.JSON(strings.NewReader(VRFV08TestHelperABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
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

func (_VRFV08TestHelper *VRFV08TestHelperCaller) HASHTOCURVEHASHPREFIX(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV08TestHelper.contract.Call(opts, &out, "HASH_TO_CURVE_HASH_PREFIX")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV08TestHelper *VRFV08TestHelperSession) HASHTOCURVEHASHPREFIX() (*big.Int, error) {
	return _VRFV08TestHelper.Contract.HASHTOCURVEHASHPREFIX(&_VRFV08TestHelper.CallOpts)
}

func (_VRFV08TestHelper *VRFV08TestHelperCallerSession) HASHTOCURVEHASHPREFIX() (*big.Int, error) {
	return _VRFV08TestHelper.Contract.HASHTOCURVEHASHPREFIX(&_VRFV08TestHelper.CallOpts)
}

func (_VRFV08TestHelper *VRFV08TestHelperCaller) PROOFLENGTH(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV08TestHelper.contract.Call(opts, &out, "PROOF_LENGTH")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV08TestHelper *VRFV08TestHelperSession) PROOFLENGTH() (*big.Int, error) {
	return _VRFV08TestHelper.Contract.PROOFLENGTH(&_VRFV08TestHelper.CallOpts)
}

func (_VRFV08TestHelper *VRFV08TestHelperCallerSession) PROOFLENGTH() (*big.Int, error) {
	return _VRFV08TestHelper.Contract.PROOFLENGTH(&_VRFV08TestHelper.CallOpts)
}

func (_VRFV08TestHelper *VRFV08TestHelperCaller) SCALARFROMCURVEPOINTSHASHPREFIX(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV08TestHelper.contract.Call(opts, &out, "SCALAR_FROM_CURVE_POINTS_HASH_PREFIX")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV08TestHelper *VRFV08TestHelperSession) SCALARFROMCURVEPOINTSHASHPREFIX() (*big.Int, error) {
	return _VRFV08TestHelper.Contract.SCALARFROMCURVEPOINTSHASHPREFIX(&_VRFV08TestHelper.CallOpts)
}

func (_VRFV08TestHelper *VRFV08TestHelperCallerSession) SCALARFROMCURVEPOINTSHASHPREFIX() (*big.Int, error) {
	return _VRFV08TestHelper.Contract.SCALARFROMCURVEPOINTSHASHPREFIX(&_VRFV08TestHelper.CallOpts)
}

func (_VRFV08TestHelper *VRFV08TestHelperCaller) VRFRANDOMOUTPUTHASHPREFIX(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV08TestHelper.contract.Call(opts, &out, "VRF_RANDOM_OUTPUT_HASH_PREFIX")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV08TestHelper *VRFV08TestHelperSession) VRFRANDOMOUTPUTHASHPREFIX() (*big.Int, error) {
	return _VRFV08TestHelper.Contract.VRFRANDOMOUTPUTHASHPREFIX(&_VRFV08TestHelper.CallOpts)
}

func (_VRFV08TestHelper *VRFV08TestHelperCallerSession) VRFRANDOMOUTPUTHASHPREFIX() (*big.Int, error) {
	return _VRFV08TestHelper.Contract.VRFRANDOMOUTPUTHASHPREFIX(&_VRFV08TestHelper.CallOpts)
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
	HASHTOCURVEHASHPREFIX(opts *bind.CallOpts) (*big.Int, error)

	PROOFLENGTH(opts *bind.CallOpts) (*big.Int, error)

	SCALARFROMCURVEPOINTSHASHPREFIX(opts *bind.CallOpts) (*big.Int, error)

	VRFRANDOMOUTPUTHASHPREFIX(opts *bind.CallOpts) (*big.Int, error)

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
