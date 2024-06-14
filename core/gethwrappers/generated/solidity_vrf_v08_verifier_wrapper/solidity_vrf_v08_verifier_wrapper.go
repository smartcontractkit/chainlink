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
	Bin: "0x608060405234801561001057600080fd5b50611b34806100206000396000f3fe608060405234801561001057600080fd5b50600436106100ea5760003560e01c80639d6f03371161008c578063b481e26011610066578063b481e260146101fc578063ef3b10ec1461020f578063fd7e4af914610224578063fe54f2a21461023757600080fd5b80639d6f0337146101c3578063a5e9508f146101d6578063aa7b2fbb146101e957600080fd5b80637f8f50a8116100c85780637f8f50a81461014c5780638af046ea1461015f57806391d5f6911461017257806395e6ee921461019557600080fd5b8063244f896d146100ef57806335452450146101185780635de600421461012b575b600080fd5b6101026100fd366004611549565b61024a565b60405161010f9190611a03565b60405180910390f35b610102610126366004611619565b610267565b61013e6101393660046118e1565b610282565b60405190815260200161010f565b61013e61015a3660046114df565b61028e565b61013e61016d366004611809565b6102a7565b61018561018036600461189a565b6102b2565b604051901515815260200161010f565b6101a86101a3366004611903565b6102c9565b6040805193845260208401929092529082015260600161010f565b61013e6101d1366004611809565b6102ea565b61013e6101e4366004611751565b6102f5565b6101856101f7366004611644565b610301565b61013e61020a366004611682565b61030e565b61022261021d366004611587565b610319565b005b6101856102323660046114c3565b610335565b610102610245366004611822565b610340565b6102526113cd565b61025d848484610363565b90505b9392505050565b61026f6113cd565b6102798383610497565b90505b92915050565b600061027983836104fb565b600061029d86868686866105ef565b9695505050505050565b600061027c8261064d565b60006102c085858585610687565b95945050505050565b60008060006102da8787878761082a565b9250925092509450945094915050565b600061027c826109c0565b60006102798383610a18565b600061025d848484610aa1565b600061027c82610c2e565b61032a898989898989898989610c88565b505050505050505050565b600061027c82610f5f565b6103486113cd565b610357888888888888886110ba565b98975050505050505050565b61036b6113cd565b83516020808601518551918601516000938493849361038c9390919061082a565b919450925090507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f858209600114610425576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f696e765a206d75737420626520696e7665727365206f66207a0000000000000060448201526064015b60405180910390fd5b60405180604001604052807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8061045e5761045e611ac9565b87860981526020017ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8785099052979650505050505050565b61049f6113cd565b6104cc600184846040516020016104b8939291906119e2565b604051602081830303815290604052611242565b90505b6104d881610f5f565b61027c5780516040805160208101929092526104f491016104b8565b90506104cf565b6000806105066113eb565b6020808252818101819052604082015260608101859052608081018490527ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f60a0820152610552611409565b60208160c08460057ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa9250826105e5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f6269674d6f64457870206661696c757265210000000000000000000000000000604482015260640161041c565b5195945050505050565b60006002868686858760405160200161060d96959493929190611970565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101209695505050505050565b600061027c8260026106807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f6001611a4f565b901c6104fb565b600073ffffffffffffffffffffffffffffffffffffffff8216610706576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f626164207769746e657373000000000000000000000000000000000000000000604482015260640161041c565b60208401516000906001161561071d57601c610720565b601b5b905060007ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd03641418587600060200201510986517ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141918203925060009190890987516040805160008082526020820180845287905260ff88169282019290925260608101929092526080820183905291925060019060a0016020604051602081039080840390855afa1580156107d7573d6000803e3d6000fd5b50506040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0015173ffffffffffffffffffffffffffffffffffffffff9081169088161495505050505050949350505050565b60008080600180827ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f897ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f038808905060007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f038a08905060006108d2838385856112aa565b90985090506108e388828e88611302565b90985090506108f488828c87611302565b909850905060006109078d878b85611302565b9098509050610918888286866112aa565b909850905061092988828e89611302565b90985090508181146109ac577ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f818a0998507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f82890997507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f81830996506109b0565b8196505b5050505050509450945094915050565b6000807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80848509840990507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f600782089392505050565b6000610a4c8360000151846020015185604001518660600151868860a001518960c001518a60e001518b6101000151610c88565b60038360200151604051602001610a64929190611a11565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101209392505050565b600082610b0a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f7a65726f207363616c6172000000000000000000000000000000000000000000604482015260640161041c565b83516020850151600090610b2090600290611a8e565b15610b2c57601c610b2f565b601b5b905060007ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd03641418387096040805160008082526020820180845281905260ff86169282019290925260608101869052608081018390529192509060019060a0016020604051602081039080840390855afa158015610baf573d6000803e3d6000fd5b505050602060405103519050600086604051602001610bce919061195e565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152919052805160209091012073ffffffffffffffffffffffffffffffffffffffff92831692169190911498975050505050505050565b805160208201205b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8110610c8357604080516020808201939093528151808203840181529082019091528051910120610c36565b919050565b610c9189610f5f565b610cf7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f7075626c6963206b6579206973206e6f74206f6e206375727665000000000000604482015260640161041c565b610d0088610f5f565b610d66576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f67616d6d61206973206e6f74206f6e2063757276650000000000000000000000604482015260640161041c565b610d6f83610f5f565b610dd5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f6347616d6d615769746e657373206973206e6f74206f6e206375727665000000604482015260640161041c565b610dde82610f5f565b610e44576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f73486173685769746e657373206973206e6f74206f6e20637572766500000000604482015260640161041c565b610e50878a8887610687565b610eb6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f6164647228632a706b2b732a6729213d5f755769746e65737300000000000000604482015260640161041c565b6000610ec28a87610497565b90506000610ed5898b878b8689896110ba565b90506000610ee6838d8d8a866105ef565b9050808a14610f51576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f696e76616c69642070726f6f6600000000000000000000000000000000000000604482015260640161041c565b505050505050505050505050565b80516000907ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f11610fec576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e76616c696420782d6f7264696e6174650000000000000000000000000000604482015260640161041c565b60208201517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f11611079576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e76616c696420792d6f7264696e6174650000000000000000000000000000604482015260640161041c565b60208201517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f9080096110b38360005b60200201516109c0565b1492915050565b6110c26113cd565b825186517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f9081900691061415611155576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f706f696e747320696e2073756d206d7573742062652064697374696e63740000604482015260640161041c565b611160878988610aa1565b6111c6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4669727374206d756c20636865636b206661696c656400000000000000000000604482015260640161041c565b6111d1848685610aa1565b611237576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f5365636f6e64206d756c20636865636b206661696c6564000000000000000000604482015260640161041c565b610357868484610363565b61124a6113cd565b61125382610c2e565b81526112686112638260006110a9565b61064d565b602082018190526002900660011415610c83576020810180517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f039052919050565b6000807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8487097ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8487099097909650945050505050565b600080807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f878509905060007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f87877ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f030990507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8183087ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f86890990999098509650505050505050565b60405180604001604052806002906020820280368337509192915050565b6040518060c001604052806006906020820280368337509192915050565b60405180602001604052806001906020820280368337509192915050565b803573ffffffffffffffffffffffffffffffffffffffff81168114610c8357600080fd5b600082601f83011261145c57600080fd5b6040516040810181811067ffffffffffffffff8211171561147f5761147f611af8565b806040525080838560408601111561149657600080fd5b60005b60028110156114b8578135835260209283019290910190600101611499565b509195945050505050565b6000604082840312156114d557600080fd5b610279838361144b565b600080600080600061012086880312156114f857600080fd5b611502878761144b565b9450611511876040880161144b565b9350611520876080880161144b565b925061152e60c08701611427565b915061153d8760e0880161144b565b90509295509295909350565b600080600060a0848603121561155e57600080fd5b611568858561144b565b9250611577856040860161144b565b9150608084013590509250925092565b60008060008060008060008060006101a08a8c0312156115a657600080fd5b6115b08b8b61144b565b98506115bf8b60408c0161144b565b975060808a0135965060a08a0135955060c08a013594506115e260e08b01611427565b93506115f28b6101008c0161144b565b92506116028b6101408c0161144b565b91506101808a013590509295985092959850929598565b6000806060838503121561162c57600080fd5b611636848461144b565b946040939093013593505050565b600080600060a0848603121561165957600080fd5b611663858561144b565b925060408401359150611679856060860161144b565b90509250925092565b60006020828403121561169457600080fd5b813567ffffffffffffffff808211156116ac57600080fd5b818401915084601f8301126116c057600080fd5b8135818111156116d2576116d2611af8565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f0116810190838211818310171561171857611718611af8565b8160405282815287602084870101111561173157600080fd5b826020860160208301376000928101602001929092525095945050505050565b6000808284036101c081121561176657600080fd5b6101a08082121561177657600080fd5b61177e611a25565b915061178a868661144b565b8252611799866040870161144b565b60208301526080850135604083015260a0850135606083015260c085013560808301526117c860e08601611427565b60a08301526101006117dc8782880161144b565b60c08401526117ef87610140880161144b565b60e084015261018086013590830152909593013593505050565b60006020828403121561181b57600080fd5b5035919050565b6000806000806000806000610160888a03121561183e57600080fd5b8735965061184f8960208a0161144b565b955061185e8960608a0161144b565b945060a088013593506118748960c08a0161144b565b9250611884896101008a0161144b565b9150610140880135905092959891949750929550565b60008060008060a085870312156118b057600080fd5b843593506118c1866020870161144b565b9250606085013591506118d660808601611427565b905092959194509250565b600080604083850312156118f457600080fd5b50508035926020909101359150565b6000806000806080858703121561191957600080fd5b5050823594602084013594506040840135936060013592509050565b8060005b6002811015611958578151845260209384019390910190600101611939565b50505050565b6119688183611935565b604001919050565b8681526119806020820187611935565b61198d6060820186611935565b61199a60a0820185611935565b6119a760e0820184611935565b60609190911b7fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166101208201526101340195945050505050565b8381526119f26020820184611935565b606081019190915260800192915050565b6040810161027c8284611935565b828152606081016102606020830184611935565b604051610120810167ffffffffffffffff81118282101715611a4957611a49611af8565b60405290565b60008219821115611a89577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b500190565b600082611ac4577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500690565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
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
