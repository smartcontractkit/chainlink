// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package solidity_verifier_wrapper

import (
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
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// VRFTestHelperABI is the input ABI used to generate the binding from.
const VRFTestHelperABI = "[{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"p1\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"p2\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"invZ\",\"type\":\"uint256\"}],\"name\":\"affineECAdd_\",\"outputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"\",\"type\":\"uint256[2]\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"base\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"exponent\",\"type\":\"uint256\"}],\"name\":\"bigModExp_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"x\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"scalar\",\"type\":\"uint256\"},{\"internalType\":\"uint256[2]\",\"name\":\"q\",\"type\":\"uint256[2]\"}],\"name\":\"ecmulVerify_\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"b\",\"type\":\"bytes\"}],\"name\":\"fieldHash_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"hashToCurve_\",\"outputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"\",\"type\":\"uint256[2]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"},{\"internalType\":\"uint256[2]\",\"name\":\"p1\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"cp1Witness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"},{\"internalType\":\"uint256[2]\",\"name\":\"p2\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"sp2Witness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"zInv\",\"type\":\"uint256\"}],\"name\":\"linearCombination_\",\"outputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"\",\"type\":\"uint256[2]\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"px\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"py\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"qx\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"qy\",\"type\":\"uint256\"}],\"name\":\"projectiveECAdd_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"proof\",\"type\":\"bytes\"}],\"name\":\"randomValueFromVRFProof_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"output\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"hash\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"gamma\",\"type\":\"uint256[2]\"},{\"internalType\":\"address\",\"name\":\"uWitness\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"v\",\"type\":\"uint256[2]\"}],\"name\":\"scalarFromCurvePoints_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"squareRoot_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"},{\"internalType\":\"uint256[2]\",\"name\":\"p\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"lcWitness\",\"type\":\"address\"}],\"name\":\"verifyLinearCombinationWithGenerator_\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"gamma\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"uWitness\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"cGammaWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"sHashWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"zInv\",\"type\":\"uint256\"}],\"name\":\"verifyVRFProof_\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"ySquared_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]"

// VRFTestHelperBin is the compiled bytecode used for deploying new contracts.
var VRFTestHelperBin = "0x608060405234801561001057600080fd5b50612220806100206000396000f3fe608060405234801561001057600080fd5b50600436106100cf5760003560e01c806395e6ee921161008c578063b481e26011610066578063b481e2601461067f578063cefda0c51461074e578063ef3b10ec1461081d578063fe54f2a21461098e576100cf565b806395e6ee92146105075780639d6f033714610575578063aa7b2fbb146105b7576100cf565b8063244f896d146100d457806335452450146101c05780635de600421461026b5780637f8f50a8146102b75780638af046ea1461041457806391d5f69114610456575b600080fd5b610182600480360360a08110156100ea57600080fd5b8101908080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f820116905080830192505050505050919291929080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f820116905080830192505050505050919291929080359060200190929190505050610b11565b6040518082600260200280838360005b838110156101ad578082015181840152602081019050610192565b5050505090500191505060405180910390f35b61022d600480360360608110156101d657600080fd5b8101908080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f820116905080830192505050505050919291929080359060200190929190505050610b2d565b6040518082600260200280838360005b8381101561025857808201518184015260208101905061023d565b5050505090500191505060405180910390f35b6102a16004803603604081101561028157600080fd5b810190808035906020019092919080359060200190929190505050610b47565b6040518082815260200191505060405180910390f35b6103fe60048036036101208110156102ce57600080fd5b8101908080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f820116905080830192505050505050919291929080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f820116905080830192505050505050919291929080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f8201169050808301925050505050509192919290803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f8201169050808301925050505050509192919290505050610b5b565b6040518082815260200191505060405180910390f35b6104406004803603602081101561042a57600080fd5b8101908080359060200190929190505050610b75565b6040518082815260200191505060405180910390f35b6104ed600480360360a081101561046c57600080fd5b810190808035906020019092919080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f820116905080830192505050505050919291929080359060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610b87565b604051808215151515815260200191505060405180910390f35b6105516004803603608081101561051d57600080fd5b8101908080359060200190929190803590602001909291908035906020019092919080359060200190929190505050610b9f565b60405180848152602001838152602001828152602001935050505060405180910390f35b6105a16004803603602081101561058b57600080fd5b8101908080359060200190929190505050610bc0565b6040518082815260200191505060405180910390f35b610665600480360360a08110156105cd57600080fd5b8101908080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f82011690508083019250505050505091929192908035906020019092919080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f8201169050808301925050505050509192919290505050610bd2565b604051808215151515815260200191505060405180910390f35b6107386004803603602081101561069557600080fd5b81019080803590602001906401000000008111156106b257600080fd5b8201836020820111156106c457600080fd5b803590602001918460018302840111640100000000831117156106e657600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290505050610be8565b6040518082815260200191505060405180910390f35b6108076004803603602081101561076457600080fd5b810190808035906020019064010000000081111561078157600080fd5b82018360208201111561079357600080fd5b803590602001918460018302840111640100000000831117156107b557600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290505050610bfa565b6040518082815260200191505060405180910390f35b61098c60048036036101a081101561083457600080fd5b8101908080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f820116905080830192505050505050919291929080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f8201169050808301925050505050509192919290803590602001909291908035906020019092919080359060200190929190803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f820116905080830192505050505050919291929080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f820116905080830192505050505050919291929080359060200190929190505050610c0c565b005b610ad360048036036101608110156109a557600080fd5b810190808035906020019092919080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f820116905080830192505050505050919291929080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f82011690508083019250505050505091929192908035906020019092919080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f820116905080830192505050505050919291929080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f820116905080830192505050505050919291929080359060200190929190505050610c28565b6040518082600260200280838360005b83811015610afe578082015181840152602081019050610ae3565b5050505090500191505060405180910390f35b610b1961211f565b610b24848484610c4c565b90509392505050565b610b3561211f565b610b3f8383610dca565b905092915050565b6000610b538383610e89565b905092915050565b6000610b6a8686868686610feb565b905095945050505050565b6000610b808261111a565b9050919050565b6000610b9585858585611154565b9050949350505050565b6000806000610bb08787878761138c565b9250925092509450945094915050565b6000610bcb82611560565b9050919050565b6000610bdf8484846115ee565b90509392505050565b6000610bf382611779565b9050919050565b6000610c05826117e6565b9050919050565b610c1d8989898989898989896119b6565b505050505050505050565b610c3061211f565b610c3f88888888888888611ce1565b9050979650505050505050565b610c5461211f565b6000806000610ca987600060028110610c6957fe5b602002015188600160028110610c7b57fe5b602002015188600060028110610c8d57fe5b602002015189600160028110610c9f57fe5b602002015161138c565b80935081945082955050505060017ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80610cdf57fe5b86830914610d55576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260198152602001807f696e765a206d75737420626520696e7665727365206f66207a0000000000000081525060200191505060405180910390fd5b60405180604001604052807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80610d8857fe5b87860981526020017ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80610db857fe5b87850981525093505050509392505050565b610dd261211f565b610e33600184846040516020018084815260200183600260200280838360005b83811015610e0d578082015181840152602081019050610df2565b505050509050018281526020019350505050604051602081830303815290604052611e85565b90505b610e3f81611f58565b610e8357610e7c81600060028110610e5357fe5b602002015160405160200180828152602001915050604051602081830303815290604052611e85565b9050610e36565b92915050565b600080610e94612141565b602081600060068110610ea357fe5b602002018181525050602081600160068110610ebb57fe5b602002018181525050602081600260068110610ed357fe5b6020020181815250508481600360068110610eea57fe5b6020020181815250508381600460068110610f0157fe5b6020020181815250507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f81600560068110610f3857fe5b602002018181525050610f49612163565b60208160c0846005600019fa92506000831415610fce576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260128152602001807f6269674d6f64457870206661696c75726521000000000000000000000000000081525060200191505060405180910390fd5b80600060018110610fdb57fe5b6020020151935050505092915050565b6000600286868685876040516020018087815260200186600260200280838360005b8381101561102857808201518184015260208101905061100d565b5050505090500185600260200280838360005b8381101561105657808201518184015260208101905061103b565b5050505090500184600260200280838360005b83811015611084578082015181840152602081019050611069565b5050505090500183600260200280838360005b838110156110b2578082015181840152602081019050611097565b505050509050018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1660601b815260140196505050505050506040516020818303038152906040528051906020012060001c905095945050505050565b600061114d82600260017ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f01901c610e89565b9050919050565b60008073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614156111f8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600b8152602001807f626164207769746e65737300000000000000000000000000000000000000000081525060200191505060405180910390fd5b60008060028660016002811061120a57fe5b60200201518161121657fe5b061461122357601c611226565b601b5b905060007ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd03641418061125257fe5b858760006002811061126057fe5b6020020151097ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd03641410360001b905060007ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141806112b757fe5b876000600281106112c457fe5b6020020151890960001b90506000600183858a6000600281106112e357fe5b602002015160001b8560405160008152602001604052604051808581526020018460ff1660ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa158015611343573d6000803e3d6000fd5b5050506020604051035190508573ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614945050505050949350505050565b60008060008060006001809150915060007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f806113c557fe5b897ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f038808905060007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8061141657fe5b8b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f038a089050600061144b83838585611fc9565b809250819950505061145f88828e88612033565b809250819950505061147388828c87612033565b809250819950505060006114898d878b85612033565b809250819950505061149d88828686611fc9565b80925081995050506114b188828e89612033565b809250819950505080821461154c577ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f806114e857fe5b818a0998507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8061151557fe5b82890997507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8061154257fe5b8183099650611550565b8196505b5050505050509450945094915050565b6000807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8061158b57fe5b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f806115b357fe5b848509840990507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f806115e257fe5b60078208915050919050565b6000808314156115fd57600080fd5b60008460006002811061160c57fe5b6020020151905060008060028760016002811061162557fe5b60200201518161163157fe5b061461163e57601c611641565b601b5b905060007ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd03641418061166d57fe5b83870960001b9050600060016000801b848660001b8560405160008152602001604052604051808581526020018460ff1660ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa1580156116da573d6000803e3d6000fd5b5050506020604051035190506000866040516020018082600260200280838360005b838110156117175780820151818401526020810190506116fc565b505050509050019150506040516020818303038152906040528051906020012060001c90508073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614955050505050509392505050565b6000818051906020012060001c90505b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f81106117e15780604051602001808281526020019150506040516020818303038152906040528051906020012060001c9050611789565b919050565b60006101a0825114611860576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260128152602001807f77726f6e672070726f6f66206c656e677468000000000000000000000000000081525060200191505060405180910390fd5b61186861211f565b61187061211f565b611878612185565b600061188261211f565b61188a61211f565b6000888060200190516101a08110156118a257600080fd5b810190809190826040019190826040019190826060018051906020019092919091908260400191908260400180519060200190929190505050869650859550849450839350829250819150809750819850829950839a50849b50859c50869d505050505050505061194d87878760006003811061191b57fe5b60200201518860016003811061192d57fe5b60200201518960026003811061193f57fe5b6020020151898989896119b6565b6003866040516020018083815260200182600260200280838360005b83811015611984578082015181840152602081019050611969565b50505050905001925050506040516020818303038152906040528051906020012060001c975050505050505050919050565b6119bf89611f58565b611a31576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601a8152602001807f7075626c6963206b6579206973206e6f74206f6e20637572766500000000000081525060200191505060405180910390fd5b611a3a88611f58565b611aac576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260158152602001807f67616d6d61206973206e6f74206f6e206375727665000000000000000000000081525060200191505060405180910390fd5b611ab583611f58565b611b27576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601d8152602001807f6347616d6d615769746e657373206973206e6f74206f6e20637572766500000081525060200191505060405180910390fd5b611b3082611f58565b611ba2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601c8152602001807f73486173685769746e657373206973206e6f74206f6e2063757276650000000081525060200191505060405180910390fd5b611bae878a8887611154565b611c20576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601a8152602001807f6164647228632a706b2b732a6729e289a05f755769746e65737300000000000081525060200191505060405180910390fd5b611c2861211f565b611c328a87610dca565b9050611c3c61211f565b611c4b898b878b868989611ce1565b90506000611c5c838d8d8a86610feb565b9050808a14611cd3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600d8152602001807f696e76616c69642070726f6f660000000000000000000000000000000000000081525060200191505060405180910390fd5b505050505050505050505050565b611ce961211f565b60007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f84600060028110611d1957fe5b602002015188600060028110611d2b57fe5b60200201510381611d3857fe5b061415611dad576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f706f696e747320696e2073756d206d7573742062652064697374696e6374000081525060200191505060405180910390fd5b611db88789886115ee565b611e0d576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260218152602001806121a86021913960400191505060405180910390fd5b611e188486856115ee565b611e6d576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260228152602001806121c96022913960400191505060405180910390fd5b611e78868484610c4c565b9050979650505050505050565b611e8d61211f565b611e9682611779565b81600060028110611ea357fe5b602002018181525050611ece611ec982600060028110611ebf57fe5b6020020151611560565b61111a565b81600160028110611edb57fe5b6020020181815250506001600282600160028110611ef557fe5b602002015181611f0157fe5b061415611f535780600160028110611f1557fe5b60200201517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f0381600160028110611f4957fe5b6020020181815250505b919050565b60007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80611f8257fe5b82600160028110611f8f57fe5b602002015183600160028110611fa157fe5b602002015109611fc183600060028110611fb757fe5b6020020151611560565b149050919050565b6000807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80611ff457fe5b8487097ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8061201f57fe5b848709809250819350505094509492505050565b60008060007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8061206057fe5b878509905060007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8061208f57fe5b87877ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f030990507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f806120de57fe5b8183087ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8061210957fe5b8689098094508195505050505094509492505050565b6040518060400160405280600290602082028038833980820191505090505090565b6040518060c00160405280600690602082028038833980820191505090505090565b6040518060200160405280600190602082028038833980820191505090505090565b604051806060016040528060039060208202803883398082019150509050509056fe4669727374206d756c7469706c69636174696f6e20636865636b206661696c65645365636f6e64206d756c7469706c69636174696f6e20636865636b206661696c6564a2646970667358220000000000000000000000000000000000000000000000000000000000000000000064736f6c63430000000033"

// DeployVRFTestHelper deploys a new Ethereum contract, binding an instance of VRFTestHelper to it.
func DeployVRFTestHelper(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *VRFTestHelper, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFTestHelperABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(VRFTestHelperBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFTestHelper{VRFTestHelperCaller: VRFTestHelperCaller{contract: contract}, VRFTestHelperTransactor: VRFTestHelperTransactor{contract: contract}, VRFTestHelperFilterer: VRFTestHelperFilterer{contract: contract}}, nil
}

// VRFTestHelper is an auto generated Go binding around an Ethereum contract.
type VRFTestHelper struct {
	VRFTestHelperCaller     // Read-only binding to the contract
	VRFTestHelperTransactor // Write-only binding to the contract
	VRFTestHelperFilterer   // Log filterer for contract events
}

// VRFTestHelperCaller is an auto generated read-only Go binding around an Ethereum contract.
type VRFTestHelperCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFTestHelperTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VRFTestHelperTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFTestHelperFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VRFTestHelperFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFTestHelperSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VRFTestHelperSession struct {
	Contract     *VRFTestHelper    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VRFTestHelperCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VRFTestHelperCallerSession struct {
	Contract *VRFTestHelperCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// VRFTestHelperTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VRFTestHelperTransactorSession struct {
	Contract     *VRFTestHelperTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// VRFTestHelperRaw is an auto generated low-level Go binding around an Ethereum contract.
type VRFTestHelperRaw struct {
	Contract *VRFTestHelper // Generic contract binding to access the raw methods on
}

// VRFTestHelperCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VRFTestHelperCallerRaw struct {
	Contract *VRFTestHelperCaller // Generic read-only contract binding to access the raw methods on
}

// VRFTestHelperTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VRFTestHelperTransactorRaw struct {
	Contract *VRFTestHelperTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVRFTestHelper creates a new instance of VRFTestHelper, bound to a specific deployed contract.
func NewVRFTestHelper(address common.Address, backend bind.ContractBackend) (*VRFTestHelper, error) {
	contract, err := bindVRFTestHelper(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFTestHelper{VRFTestHelperCaller: VRFTestHelperCaller{contract: contract}, VRFTestHelperTransactor: VRFTestHelperTransactor{contract: contract}, VRFTestHelperFilterer: VRFTestHelperFilterer{contract: contract}}, nil
}

// NewVRFTestHelperCaller creates a new read-only instance of VRFTestHelper, bound to a specific deployed contract.
func NewVRFTestHelperCaller(address common.Address, caller bind.ContractCaller) (*VRFTestHelperCaller, error) {
	contract, err := bindVRFTestHelper(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFTestHelperCaller{contract: contract}, nil
}

// NewVRFTestHelperTransactor creates a new write-only instance of VRFTestHelper, bound to a specific deployed contract.
func NewVRFTestHelperTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFTestHelperTransactor, error) {
	contract, err := bindVRFTestHelper(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFTestHelperTransactor{contract: contract}, nil
}

// NewVRFTestHelperFilterer creates a new log filterer instance of VRFTestHelper, bound to a specific deployed contract.
func NewVRFTestHelperFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFTestHelperFilterer, error) {
	contract, err := bindVRFTestHelper(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFTestHelperFilterer{contract: contract}, nil
}

// bindVRFTestHelper binds a generic wrapper to an already deployed contract.
func bindVRFTestHelper(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFTestHelperABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VRFTestHelper *VRFTestHelperRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _VRFTestHelper.Contract.VRFTestHelperCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VRFTestHelper *VRFTestHelperRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFTestHelper.Contract.VRFTestHelperTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VRFTestHelper *VRFTestHelperRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFTestHelper.Contract.VRFTestHelperTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VRFTestHelper *VRFTestHelperCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _VRFTestHelper.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VRFTestHelper *VRFTestHelperTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFTestHelper.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VRFTestHelper *VRFTestHelperTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFTestHelper.Contract.contract.Transact(opts, method, params...)
}

// AffineECAdd is a free data retrieval call binding the contract method 0x244f896d.
//
// Solidity: function affineECAdd_(uint256[2] p1, uint256[2] p2, uint256 invZ) constant returns(uint256[2])
func (_VRFTestHelper *VRFTestHelperCaller) AffineECAdd(opts *bind.CallOpts, p1 [2]*big.Int, p2 [2]*big.Int, invZ *big.Int) ([2]*big.Int, error) {
	var (
		ret0 = new([2]*big.Int)
	)
	out := ret0
	err := _VRFTestHelper.contract.Call(opts, out, "affineECAdd_", p1, p2, invZ)
	return *ret0, err
}

// AffineECAdd is a free data retrieval call binding the contract method 0x244f896d.
//
// Solidity: function affineECAdd_(uint256[2] p1, uint256[2] p2, uint256 invZ) constant returns(uint256[2])
func (_VRFTestHelper *VRFTestHelperSession) AffineECAdd(p1 [2]*big.Int, p2 [2]*big.Int, invZ *big.Int) ([2]*big.Int, error) {
	return _VRFTestHelper.Contract.AffineECAdd(&_VRFTestHelper.CallOpts, p1, p2, invZ)
}

// AffineECAdd is a free data retrieval call binding the contract method 0x244f896d.
//
// Solidity: function affineECAdd_(uint256[2] p1, uint256[2] p2, uint256 invZ) constant returns(uint256[2])
func (_VRFTestHelper *VRFTestHelperCallerSession) AffineECAdd(p1 [2]*big.Int, p2 [2]*big.Int, invZ *big.Int) ([2]*big.Int, error) {
	return _VRFTestHelper.Contract.AffineECAdd(&_VRFTestHelper.CallOpts, p1, p2, invZ)
}

// BigModExp is a free data retrieval call binding the contract method 0x5de60042.
//
// Solidity: function bigModExp_(uint256 base, uint256 exponent) constant returns(uint256)
func (_VRFTestHelper *VRFTestHelperCaller) BigModExp(opts *bind.CallOpts, base *big.Int, exponent *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _VRFTestHelper.contract.Call(opts, out, "bigModExp_", base, exponent)
	return *ret0, err
}

// BigModExp is a free data retrieval call binding the contract method 0x5de60042.
//
// Solidity: function bigModExp_(uint256 base, uint256 exponent) constant returns(uint256)
func (_VRFTestHelper *VRFTestHelperSession) BigModExp(base *big.Int, exponent *big.Int) (*big.Int, error) {
	return _VRFTestHelper.Contract.BigModExp(&_VRFTestHelper.CallOpts, base, exponent)
}

// BigModExp is a free data retrieval call binding the contract method 0x5de60042.
//
// Solidity: function bigModExp_(uint256 base, uint256 exponent) constant returns(uint256)
func (_VRFTestHelper *VRFTestHelperCallerSession) BigModExp(base *big.Int, exponent *big.Int) (*big.Int, error) {
	return _VRFTestHelper.Contract.BigModExp(&_VRFTestHelper.CallOpts, base, exponent)
}

// EcmulVerify is a free data retrieval call binding the contract method 0xaa7b2fbb.
//
// Solidity: function ecmulVerify_(uint256[2] x, uint256 scalar, uint256[2] q) constant returns(bool)
func (_VRFTestHelper *VRFTestHelperCaller) EcmulVerify(opts *bind.CallOpts, x [2]*big.Int, scalar *big.Int, q [2]*big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _VRFTestHelper.contract.Call(opts, out, "ecmulVerify_", x, scalar, q)
	return *ret0, err
}

// EcmulVerify is a free data retrieval call binding the contract method 0xaa7b2fbb.
//
// Solidity: function ecmulVerify_(uint256[2] x, uint256 scalar, uint256[2] q) constant returns(bool)
func (_VRFTestHelper *VRFTestHelperSession) EcmulVerify(x [2]*big.Int, scalar *big.Int, q [2]*big.Int) (bool, error) {
	return _VRFTestHelper.Contract.EcmulVerify(&_VRFTestHelper.CallOpts, x, scalar, q)
}

// EcmulVerify is a free data retrieval call binding the contract method 0xaa7b2fbb.
//
// Solidity: function ecmulVerify_(uint256[2] x, uint256 scalar, uint256[2] q) constant returns(bool)
func (_VRFTestHelper *VRFTestHelperCallerSession) EcmulVerify(x [2]*big.Int, scalar *big.Int, q [2]*big.Int) (bool, error) {
	return _VRFTestHelper.Contract.EcmulVerify(&_VRFTestHelper.CallOpts, x, scalar, q)
}

// FieldHash is a free data retrieval call binding the contract method 0xb481e260.
//
// Solidity: function fieldHash_(bytes b) constant returns(uint256)
func (_VRFTestHelper *VRFTestHelperCaller) FieldHash(opts *bind.CallOpts, b []byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _VRFTestHelper.contract.Call(opts, out, "fieldHash_", b)
	return *ret0, err
}

// FieldHash is a free data retrieval call binding the contract method 0xb481e260.
//
// Solidity: function fieldHash_(bytes b) constant returns(uint256)
func (_VRFTestHelper *VRFTestHelperSession) FieldHash(b []byte) (*big.Int, error) {
	return _VRFTestHelper.Contract.FieldHash(&_VRFTestHelper.CallOpts, b)
}

// FieldHash is a free data retrieval call binding the contract method 0xb481e260.
//
// Solidity: function fieldHash_(bytes b) constant returns(uint256)
func (_VRFTestHelper *VRFTestHelperCallerSession) FieldHash(b []byte) (*big.Int, error) {
	return _VRFTestHelper.Contract.FieldHash(&_VRFTestHelper.CallOpts, b)
}

// HashToCurve is a free data retrieval call binding the contract method 0x35452450.
//
// Solidity: function hashToCurve_(uint256[2] pk, uint256 x) constant returns(uint256[2])
func (_VRFTestHelper *VRFTestHelperCaller) HashToCurve(opts *bind.CallOpts, pk [2]*big.Int, x *big.Int) ([2]*big.Int, error) {
	var (
		ret0 = new([2]*big.Int)
	)
	out := ret0
	err := _VRFTestHelper.contract.Call(opts, out, "hashToCurve_", pk, x)
	return *ret0, err
}

// HashToCurve is a free data retrieval call binding the contract method 0x35452450.
//
// Solidity: function hashToCurve_(uint256[2] pk, uint256 x) constant returns(uint256[2])
func (_VRFTestHelper *VRFTestHelperSession) HashToCurve(pk [2]*big.Int, x *big.Int) ([2]*big.Int, error) {
	return _VRFTestHelper.Contract.HashToCurve(&_VRFTestHelper.CallOpts, pk, x)
}

// HashToCurve is a free data retrieval call binding the contract method 0x35452450.
//
// Solidity: function hashToCurve_(uint256[2] pk, uint256 x) constant returns(uint256[2])
func (_VRFTestHelper *VRFTestHelperCallerSession) HashToCurve(pk [2]*big.Int, x *big.Int) ([2]*big.Int, error) {
	return _VRFTestHelper.Contract.HashToCurve(&_VRFTestHelper.CallOpts, pk, x)
}

// LinearCombination is a free data retrieval call binding the contract method 0xfe54f2a2.
//
// Solidity: function linearCombination_(uint256 c, uint256[2] p1, uint256[2] cp1Witness, uint256 s, uint256[2] p2, uint256[2] sp2Witness, uint256 zInv) constant returns(uint256[2])
func (_VRFTestHelper *VRFTestHelperCaller) LinearCombination(opts *bind.CallOpts, c *big.Int, p1 [2]*big.Int, cp1Witness [2]*big.Int, s *big.Int, p2 [2]*big.Int, sp2Witness [2]*big.Int, zInv *big.Int) ([2]*big.Int, error) {
	var (
		ret0 = new([2]*big.Int)
	)
	out := ret0
	err := _VRFTestHelper.contract.Call(opts, out, "linearCombination_", c, p1, cp1Witness, s, p2, sp2Witness, zInv)
	return *ret0, err
}

// LinearCombination is a free data retrieval call binding the contract method 0xfe54f2a2.
//
// Solidity: function linearCombination_(uint256 c, uint256[2] p1, uint256[2] cp1Witness, uint256 s, uint256[2] p2, uint256[2] sp2Witness, uint256 zInv) constant returns(uint256[2])
func (_VRFTestHelper *VRFTestHelperSession) LinearCombination(c *big.Int, p1 [2]*big.Int, cp1Witness [2]*big.Int, s *big.Int, p2 [2]*big.Int, sp2Witness [2]*big.Int, zInv *big.Int) ([2]*big.Int, error) {
	return _VRFTestHelper.Contract.LinearCombination(&_VRFTestHelper.CallOpts, c, p1, cp1Witness, s, p2, sp2Witness, zInv)
}

// LinearCombination is a free data retrieval call binding the contract method 0xfe54f2a2.
//
// Solidity: function linearCombination_(uint256 c, uint256[2] p1, uint256[2] cp1Witness, uint256 s, uint256[2] p2, uint256[2] sp2Witness, uint256 zInv) constant returns(uint256[2])
func (_VRFTestHelper *VRFTestHelperCallerSession) LinearCombination(c *big.Int, p1 [2]*big.Int, cp1Witness [2]*big.Int, s *big.Int, p2 [2]*big.Int, sp2Witness [2]*big.Int, zInv *big.Int) ([2]*big.Int, error) {
	return _VRFTestHelper.Contract.LinearCombination(&_VRFTestHelper.CallOpts, c, p1, cp1Witness, s, p2, sp2Witness, zInv)
}

// ProjectiveECAdd is a free data retrieval call binding the contract method 0x95e6ee92.
//
// Solidity: function projectiveECAdd_(uint256 px, uint256 py, uint256 qx, uint256 qy) constant returns(uint256, uint256, uint256)
func (_VRFTestHelper *VRFTestHelperCaller) ProjectiveECAdd(opts *bind.CallOpts, px *big.Int, py *big.Int, qx *big.Int, qy *big.Int) (*big.Int, *big.Int, *big.Int, error) {
	var (
		ret0 = new(*big.Int)
		ret1 = new(*big.Int)
		ret2 = new(*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
	}
	err := _VRFTestHelper.contract.Call(opts, out, "projectiveECAdd_", px, py, qx, qy)
	return *ret0, *ret1, *ret2, err
}

// ProjectiveECAdd is a free data retrieval call binding the contract method 0x95e6ee92.
//
// Solidity: function projectiveECAdd_(uint256 px, uint256 py, uint256 qx, uint256 qy) constant returns(uint256, uint256, uint256)
func (_VRFTestHelper *VRFTestHelperSession) ProjectiveECAdd(px *big.Int, py *big.Int, qx *big.Int, qy *big.Int) (*big.Int, *big.Int, *big.Int, error) {
	return _VRFTestHelper.Contract.ProjectiveECAdd(&_VRFTestHelper.CallOpts, px, py, qx, qy)
}

// ProjectiveECAdd is a free data retrieval call binding the contract method 0x95e6ee92.
//
// Solidity: function projectiveECAdd_(uint256 px, uint256 py, uint256 qx, uint256 qy) constant returns(uint256, uint256, uint256)
func (_VRFTestHelper *VRFTestHelperCallerSession) ProjectiveECAdd(px *big.Int, py *big.Int, qx *big.Int, qy *big.Int) (*big.Int, *big.Int, *big.Int, error) {
	return _VRFTestHelper.Contract.ProjectiveECAdd(&_VRFTestHelper.CallOpts, px, py, qx, qy)
}

// RandomValueFromVRFProof is a free data retrieval call binding the contract method 0xcefda0c5.
//
// Solidity: function randomValueFromVRFProof_(bytes proof) constant returns(uint256 output)
func (_VRFTestHelper *VRFTestHelperCaller) RandomValueFromVRFProof(opts *bind.CallOpts, proof []byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _VRFTestHelper.contract.Call(opts, out, "randomValueFromVRFProof_", proof)
	return *ret0, err
}

// RandomValueFromVRFProof is a free data retrieval call binding the contract method 0xcefda0c5.
//
// Solidity: function randomValueFromVRFProof_(bytes proof) constant returns(uint256 output)
func (_VRFTestHelper *VRFTestHelperSession) RandomValueFromVRFProof(proof []byte) (*big.Int, error) {
	return _VRFTestHelper.Contract.RandomValueFromVRFProof(&_VRFTestHelper.CallOpts, proof)
}

// RandomValueFromVRFProof is a free data retrieval call binding the contract method 0xcefda0c5.
//
// Solidity: function randomValueFromVRFProof_(bytes proof) constant returns(uint256 output)
func (_VRFTestHelper *VRFTestHelperCallerSession) RandomValueFromVRFProof(proof []byte) (*big.Int, error) {
	return _VRFTestHelper.Contract.RandomValueFromVRFProof(&_VRFTestHelper.CallOpts, proof)
}

// ScalarFromCurvePoints is a free data retrieval call binding the contract method 0x7f8f50a8.
//
// Solidity: function scalarFromCurvePoints_(uint256[2] hash, uint256[2] pk, uint256[2] gamma, address uWitness, uint256[2] v) constant returns(uint256)
func (_VRFTestHelper *VRFTestHelperCaller) ScalarFromCurvePoints(opts *bind.CallOpts, hash [2]*big.Int, pk [2]*big.Int, gamma [2]*big.Int, uWitness common.Address, v [2]*big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _VRFTestHelper.contract.Call(opts, out, "scalarFromCurvePoints_", hash, pk, gamma, uWitness, v)
	return *ret0, err
}

// ScalarFromCurvePoints is a free data retrieval call binding the contract method 0x7f8f50a8.
//
// Solidity: function scalarFromCurvePoints_(uint256[2] hash, uint256[2] pk, uint256[2] gamma, address uWitness, uint256[2] v) constant returns(uint256)
func (_VRFTestHelper *VRFTestHelperSession) ScalarFromCurvePoints(hash [2]*big.Int, pk [2]*big.Int, gamma [2]*big.Int, uWitness common.Address, v [2]*big.Int) (*big.Int, error) {
	return _VRFTestHelper.Contract.ScalarFromCurvePoints(&_VRFTestHelper.CallOpts, hash, pk, gamma, uWitness, v)
}

// ScalarFromCurvePoints is a free data retrieval call binding the contract method 0x7f8f50a8.
//
// Solidity: function scalarFromCurvePoints_(uint256[2] hash, uint256[2] pk, uint256[2] gamma, address uWitness, uint256[2] v) constant returns(uint256)
func (_VRFTestHelper *VRFTestHelperCallerSession) ScalarFromCurvePoints(hash [2]*big.Int, pk [2]*big.Int, gamma [2]*big.Int, uWitness common.Address, v [2]*big.Int) (*big.Int, error) {
	return _VRFTestHelper.Contract.ScalarFromCurvePoints(&_VRFTestHelper.CallOpts, hash, pk, gamma, uWitness, v)
}

// SquareRoot is a free data retrieval call binding the contract method 0x8af046ea.
//
// Solidity: function squareRoot_(uint256 x) constant returns(uint256)
func (_VRFTestHelper *VRFTestHelperCaller) SquareRoot(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _VRFTestHelper.contract.Call(opts, out, "squareRoot_", x)
	return *ret0, err
}

// SquareRoot is a free data retrieval call binding the contract method 0x8af046ea.
//
// Solidity: function squareRoot_(uint256 x) constant returns(uint256)
func (_VRFTestHelper *VRFTestHelperSession) SquareRoot(x *big.Int) (*big.Int, error) {
	return _VRFTestHelper.Contract.SquareRoot(&_VRFTestHelper.CallOpts, x)
}

// SquareRoot is a free data retrieval call binding the contract method 0x8af046ea.
//
// Solidity: function squareRoot_(uint256 x) constant returns(uint256)
func (_VRFTestHelper *VRFTestHelperCallerSession) SquareRoot(x *big.Int) (*big.Int, error) {
	return _VRFTestHelper.Contract.SquareRoot(&_VRFTestHelper.CallOpts, x)
}

// VerifyLinearCombinationWithGenerator is a free data retrieval call binding the contract method 0x91d5f691.
//
// Solidity: function verifyLinearCombinationWithGenerator_(uint256 c, uint256[2] p, uint256 s, address lcWitness) constant returns(bool)
func (_VRFTestHelper *VRFTestHelperCaller) VerifyLinearCombinationWithGenerator(opts *bind.CallOpts, c *big.Int, p [2]*big.Int, s *big.Int, lcWitness common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _VRFTestHelper.contract.Call(opts, out, "verifyLinearCombinationWithGenerator_", c, p, s, lcWitness)
	return *ret0, err
}

// VerifyLinearCombinationWithGenerator is a free data retrieval call binding the contract method 0x91d5f691.
//
// Solidity: function verifyLinearCombinationWithGenerator_(uint256 c, uint256[2] p, uint256 s, address lcWitness) constant returns(bool)
func (_VRFTestHelper *VRFTestHelperSession) VerifyLinearCombinationWithGenerator(c *big.Int, p [2]*big.Int, s *big.Int, lcWitness common.Address) (bool, error) {
	return _VRFTestHelper.Contract.VerifyLinearCombinationWithGenerator(&_VRFTestHelper.CallOpts, c, p, s, lcWitness)
}

// VerifyLinearCombinationWithGenerator is a free data retrieval call binding the contract method 0x91d5f691.
//
// Solidity: function verifyLinearCombinationWithGenerator_(uint256 c, uint256[2] p, uint256 s, address lcWitness) constant returns(bool)
func (_VRFTestHelper *VRFTestHelperCallerSession) VerifyLinearCombinationWithGenerator(c *big.Int, p [2]*big.Int, s *big.Int, lcWitness common.Address) (bool, error) {
	return _VRFTestHelper.Contract.VerifyLinearCombinationWithGenerator(&_VRFTestHelper.CallOpts, c, p, s, lcWitness)
}

// VerifyVRFProof is a free data retrieval call binding the contract method 0xef3b10ec.
//
// Solidity: function verifyVRFProof_(uint256[2] pk, uint256[2] gamma, uint256 c, uint256 s, uint256 seed, address uWitness, uint256[2] cGammaWitness, uint256[2] sHashWitness, uint256 zInv) constant returns()
func (_VRFTestHelper *VRFTestHelperCaller) VerifyVRFProof(opts *bind.CallOpts, pk [2]*big.Int, gamma [2]*big.Int, c *big.Int, s *big.Int, seed *big.Int, uWitness common.Address, cGammaWitness [2]*big.Int, sHashWitness [2]*big.Int, zInv *big.Int) error {
	var ()
	out := &[]interface{}{}
	err := _VRFTestHelper.contract.Call(opts, out, "verifyVRFProof_", pk, gamma, c, s, seed, uWitness, cGammaWitness, sHashWitness, zInv)
	return err
}

// VerifyVRFProof is a free data retrieval call binding the contract method 0xef3b10ec.
//
// Solidity: function verifyVRFProof_(uint256[2] pk, uint256[2] gamma, uint256 c, uint256 s, uint256 seed, address uWitness, uint256[2] cGammaWitness, uint256[2] sHashWitness, uint256 zInv) constant returns()
func (_VRFTestHelper *VRFTestHelperSession) VerifyVRFProof(pk [2]*big.Int, gamma [2]*big.Int, c *big.Int, s *big.Int, seed *big.Int, uWitness common.Address, cGammaWitness [2]*big.Int, sHashWitness [2]*big.Int, zInv *big.Int) error {
	return _VRFTestHelper.Contract.VerifyVRFProof(&_VRFTestHelper.CallOpts, pk, gamma, c, s, seed, uWitness, cGammaWitness, sHashWitness, zInv)
}

// VerifyVRFProof is a free data retrieval call binding the contract method 0xef3b10ec.
//
// Solidity: function verifyVRFProof_(uint256[2] pk, uint256[2] gamma, uint256 c, uint256 s, uint256 seed, address uWitness, uint256[2] cGammaWitness, uint256[2] sHashWitness, uint256 zInv) constant returns()
func (_VRFTestHelper *VRFTestHelperCallerSession) VerifyVRFProof(pk [2]*big.Int, gamma [2]*big.Int, c *big.Int, s *big.Int, seed *big.Int, uWitness common.Address, cGammaWitness [2]*big.Int, sHashWitness [2]*big.Int, zInv *big.Int) error {
	return _VRFTestHelper.Contract.VerifyVRFProof(&_VRFTestHelper.CallOpts, pk, gamma, c, s, seed, uWitness, cGammaWitness, sHashWitness, zInv)
}

// YSquared is a free data retrieval call binding the contract method 0x9d6f0337.
//
// Solidity: function ySquared_(uint256 x) constant returns(uint256)
func (_VRFTestHelper *VRFTestHelperCaller) YSquared(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _VRFTestHelper.contract.Call(opts, out, "ySquared_", x)
	return *ret0, err
}

// YSquared is a free data retrieval call binding the contract method 0x9d6f0337.
//
// Solidity: function ySquared_(uint256 x) constant returns(uint256)
func (_VRFTestHelper *VRFTestHelperSession) YSquared(x *big.Int) (*big.Int, error) {
	return _VRFTestHelper.Contract.YSquared(&_VRFTestHelper.CallOpts, x)
}

// YSquared is a free data retrieval call binding the contract method 0x9d6f0337.
//
// Solidity: function ySquared_(uint256 x) constant returns(uint256)
func (_VRFTestHelper *VRFTestHelperCallerSession) YSquared(x *big.Int) (*big.Int, error) {
	return _VRFTestHelper.Contract.YSquared(&_VRFTestHelper.CallOpts, x)
}
