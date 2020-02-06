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
const VRFTestHelperABI = "[{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"p1\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"p2\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"invZ\",\"type\":\"uint256\"}],\"name\":\"affineECAdd_\",\"outputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"\",\"type\":\"uint256[2]\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"base\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"exponent\",\"type\":\"uint256\"}],\"name\":\"bigModExp_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"x\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"scalar\",\"type\":\"uint256\"},{\"internalType\":\"uint256[2]\",\"name\":\"q\",\"type\":\"uint256[2]\"}],\"name\":\"ecmulVerify_\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"fieldHash_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"hashToCurve_\",\"outputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"\",\"type\":\"uint256[2]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"},{\"internalType\":\"uint256[2]\",\"name\":\"p1\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"cp1Witness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"},{\"internalType\":\"uint256[2]\",\"name\":\"p2\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"sp2Witness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"zInv\",\"type\":\"uint256\"}],\"name\":\"linearCombination_\",\"outputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"\",\"type\":\"uint256[2]\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"px\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"py\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"qx\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"qy\",\"type\":\"uint256\"}],\"name\":\"projectiveECAdd_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"proof\",\"type\":\"bytes\"}],\"name\":\"randomValueFromVRFProof_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"output\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"hash\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"gamma\",\"type\":\"uint256[2]\"},{\"internalType\":\"address\",\"name\":\"uWitness\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"v\",\"type\":\"uint256[2]\"}],\"name\":\"scalarFromCurvePoints_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"squareRoot_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"},{\"internalType\":\"uint256[2]\",\"name\":\"p\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"lcWitness\",\"type\":\"address\"}],\"name\":\"verifyLinearCombinationWithGenerator_\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"gamma\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"uWitness\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"cGammaWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"sHashWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"zInv\",\"type\":\"uint256\"}],\"name\":\"verifyVRFProof_\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"ySquared_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]"

// VRFTestHelperBin is the compiled bytecode used for deploying new contracts.
var VRFTestHelperBin = "0x608060405234801561001057600080fd5b506121b2806100206000396000f3fe608060405234801561001057600080fd5b50600436106100cf5760003560e01c806391d5f6911161008c578063aa7b2fbb11610066578063aa7b2fbb146105f9578063cefda0c5146106c1578063ef3b10ec14610790578063fe54f2a214610901576100cf565b806391d5f6911461049857806395e6ee92146105495780639d6f0337146105b7576100cf565b8063244f896d146100d457806335452450146101c05780635de600421461026b5780637f8f50a8146102b757806380aa7713146104145780638af046ea14610456575b600080fd5b610182600480360360a08110156100ea57600080fd5b8101908080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f820116905080830192505050505050919291929080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f820116905080830192505050505050919291929080359060200190929190505050610a84565b6040518082600260200280838360005b838110156101ad578082015181840152602081019050610192565b5050505090500191505060405180910390f35b61022d600480360360608110156101d657600080fd5b8101908080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f820116905080830192505050505050919291929080359060200190929190505050610aa0565b6040518082600260200280838360005b8381101561025857808201518184015260208101905061023d565b5050505090500191505060405180910390f35b6102a16004803603604081101561028157600080fd5b810190808035906020019092919080359060200190929190505050610aba565b6040518082815260200191505060405180910390f35b6103fe60048036036101208110156102ce57600080fd5b8101908080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f820116905080830192505050505050919291929080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f820116905080830192505050505050919291929080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f8201169050808301925050505050509192919290803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f8201169050808301925050505050509192919290505050610ace565b6040518082815260200191505060405180910390f35b6104406004803603602081101561042a57600080fd5b8101908080359060200190929190505050610ae8565b6040518082815260200191505060405180910390f35b6104826004803603602081101561046c57600080fd5b8101908080359060200190929190505050610afa565b6040518082815260200191505060405180910390f35b61052f600480360360a08110156104ae57600080fd5b810190808035906020019092919080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f820116905080830192505050505050919291929080359060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610b0c565b604051808215151515815260200191505060405180910390f35b6105936004803603608081101561055f57600080fd5b8101908080359060200190929190803590602001909291908035906020019092919080359060200190929190505050610b24565b60405180848152602001838152602001828152602001935050505060405180910390f35b6105e3600480360360208110156105cd57600080fd5b8101908080359060200190929190505050610b45565b6040518082815260200191505060405180910390f35b6106a7600480360360a081101561060f57600080fd5b8101908080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f82011690508083019250505050505091929192908035906020019092919080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f8201169050808301925050505050509192919290505050610b57565b604051808215151515815260200191505060405180910390f35b61077a600480360360208110156106d757600080fd5b81019080803590602001906401000000008111156106f457600080fd5b82018360208201111561070657600080fd5b8035906020019184600183028401116401000000008311171561072857600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290505050610b6d565b6040518082815260200191505060405180910390f35b6108ff60048036036101a08110156107a757600080fd5b8101908080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f820116905080830192505050505050919291929080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f8201169050808301925050505050509192919290803590602001909291908035906020019092919080359060200190929190803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f820116905080830192505050505050919291929080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f820116905080830192505050505050919291929080359060200190929190505050610b7f565b005b610a46600480360361016081101561091857600080fd5b810190808035906020019092919080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f820116905080830192505050505050919291929080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f82011690508083019250505050505091929192908035906020019092919080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f820116905080830192505050505050919291929080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f820116905080830192505050505050919291929080359060200190929190505050610b9b565b6040518082600260200280838360005b83811015610a71578082015181840152602081019050610a56565b5050505090500191505060405180910390f35b610a8c6120b1565b610a97848484610bbf565b90509392505050565b610aa86120b1565b610ab28383610d3d565b905092915050565b6000610ac68383610df3565b905092915050565b6000610add8686868686610f55565b905095945050505050565b6000610af38261107b565b9050919050565b6000610b05826110de565b9050919050565b6000610b1a85858585611118565b9050949350505050565b6000806000610b3587878787611350565b9250925092509450945094915050565b6000610b5082611524565b9050919050565b6000610b648484846115b2565b90509392505050565b6000610b7882611777565b9050919050565b610b9089898989898989898961193e565b505050505050505050565b610ba36120b1565b610bb288888888888888611c69565b9050979650505050505050565b610bc76120b1565b6000806000610c1c87600060028110610bdc57fe5b602002015188600160028110610bee57fe5b602002015188600060028110610c0057fe5b602002015189600160028110610c1257fe5b6020020151611350565b80935081945082955050505060017ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80610c5257fe5b86830914610cc8576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260198152602001807f696e765a206d75737420626520696e7665727365206f66207a0000000000000081525060200191505060405180910390fd5b60405180604001604052807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80610cfb57fe5b87860981526020017ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80610d2b57fe5b87850981525093505050509392505050565b610d456120b1565b610d9d83836040516020018083600260200280838360005b83811015610d78578082015181840152602081019050610d5d565b5050505090500182815260200192505050604051602081830303815290604052611e0d565b90505b610da981611eea565b610ded57610de681600060028110610dbd57fe5b602002015160405160200180828152602001915050604051602081830303815290604052611e0d565b9050610da0565b92915050565b600080610dfe6120d3565b602081600060068110610e0d57fe5b602002018181525050602081600160068110610e2557fe5b602002018181525050602081600260068110610e3d57fe5b6020020181815250508481600360068110610e5457fe5b6020020181815250508381600460068110610e6b57fe5b6020020181815250507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f81600560068110610ea257fe5b602002018181525050610eb36120f5565b60208160c0846005600019fa92506000831415610f38576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260128152602001807f6269674d6f64457870206661696c75726521000000000000000000000000000081525060200191505060405180910390fd5b80600060018110610f4557fe5b6020020151935050505092915050565b600085858584866040516020018086600260200280838360005b83811015610f8a578082015181840152602081019050610f6f565b5050505090500185600260200280838360005b83811015610fb8578082015181840152602081019050610f9d565b5050505090500184600260200280838360005b83811015610fe6578082015181840152602081019050610fcb565b5050505090500183600260200280838360005b83811015611014578082015181840152602081019050610ff9565b505050509050018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1660601b8152601401955050505050506040516020818303038152906040528051906020012060001c905095945050505050565b60008190505b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f81106110d95780604051602001808281526020019150506040516020818303038152906040528051906020012060001c9050611081565b919050565b600061111182600260017ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f01901c610df3565b9050919050565b60008073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614156111bc576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600b8152602001807f626164207769746e65737300000000000000000000000000000000000000000081525060200191505060405180910390fd5b6000806002866001600281106111ce57fe5b6020020151816111da57fe5b06146111e757601c6111ea565b601b5b905060007ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd03641418061121657fe5b858760006002811061122457fe5b6020020151097ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd03641410360001b905060007ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd03641418061127b57fe5b8760006002811061128857fe5b6020020151890960001b90506000600183858a6000600281106112a757fe5b602002015160001b8560405160008152602001604052604051808581526020018460ff1660ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa158015611307573d6000803e3d6000fd5b5050506020604051035190508573ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614945050505050949350505050565b60008060008060006001809150915060007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8061138957fe5b897ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f038808905060007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f806113da57fe5b8b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f038a089050600061140f83838585611f5b565b809250819950505061142388828e88611fc5565b809250819950505061143788828c87611fc5565b8092508199505050600061144d8d878b85611fc5565b809250819950505061146188828686611f5b565b809250819950505061147588828e89611fc5565b8092508199505050808214611510577ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f806114ac57fe5b818a0998507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f806114d957fe5b82890997507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8061150657fe5b8183099650611514565b8196505b5050505050509450945094915050565b6000807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8061154f57fe5b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8061157757fe5b848509840990507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f806115a657fe5b60078208915050919050565b6000808314156115c157600080fd5b6000846000600281106115d057fe5b602002015190506000806002876001600281106115e957fe5b6020020151816115f557fe5b061461160257601c611605565b601b5b905060007ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd03641418061163157fe5b83870960001b9050600060016000801b848660001b8560405160008152602001604052604051808581526020018460ff1660ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa15801561169e573d6000803e3d6000fd5b5050506020604051035190506000866040516020018082600260200280838360005b838110156116db5780820151818401526020810190506116c0565b505050509050019150506040516020818303038152906040528051906020012060001c9050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16141561173a57600080fd5b8073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614955050505050509392505050565b60006101a08251146117f1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260128152602001807f77726f6e672070726f6f66206c656e677468000000000000000000000000000081525060200191505060405180910390fd5b6117f96120b1565b6118016120b1565b611809612117565b60006118136120b1565b61181b6120b1565b6000888060200190516101a081101561183357600080fd5b810190809190826040019190826040019190826060018051906020019092919091908260400191908260400180519060200190929190505050869650859550849450839350829250819150809750819850829950839a50849b50859c50869d50505050505050506118de8787876000600381106118ac57fe5b6020020151886001600381106118be57fe5b6020020151896002600381106118d057fe5b60200201518989898961193e565b856040516020018082600260200280838360005b8381101561190d5780820151818401526020810190506118f2565b505050509050019150506040516020818303038152906040528051906020012060001c975050505050505050919050565b61194789611eea565b6119b9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601a8152602001807f7075626c6963206b6579206973206e6f74206f6e20637572766500000000000081525060200191505060405180910390fd5b6119c288611eea565b611a34576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260158152602001807f67616d6d61206973206e6f74206f6e206375727665000000000000000000000081525060200191505060405180910390fd5b611a3d83611eea565b611aaf576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601d8152602001807f6347616d6d615769746e657373206973206e6f74206f6e20637572766500000081525060200191505060405180910390fd5b611ab882611eea565b611b2a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601c8152602001807f73486173685769746e657373206973206e6f74206f6e2063757276650000000081525060200191505060405180910390fd5b611b36878a8887611118565b611ba8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601a8152602001807f6164647228632a706b2b732a6729e289a05f755769746e65737300000000000081525060200191505060405180910390fd5b611bb06120b1565b611bba8a87610d3d565b9050611bc46120b1565b611bd3898b878b868989611c69565b90506000611be4838d8d8a86610f55565b9050808a14611c5b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600d8152602001807f696e76616c69642070726f6f660000000000000000000000000000000000000081525060200191505060405180910390fd5b505050505050505050505050565b611c716120b1565b60007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f84600060028110611ca157fe5b602002015188600060028110611cb357fe5b60200201510381611cc057fe5b061415611d35576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f706f696e747320696e2073756d206d7573742062652064697374696e6374000081525060200191505060405180910390fd5b611d408789886115b2565b611d95576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602181526020018061213a6021913960400191505060405180910390fd5b611da08486856115b2565b611df5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602281526020018061215b6022913960400191505060405180910390fd5b611e00868484610bbf565b9050979650505050505050565b611e156120b1565b611e28828051906020012060001c61107b565b81600060028110611e3557fe5b602002018181525050611e60611e5b82600060028110611e5157fe5b6020020151611524565b6110de565b81600160028110611e6d57fe5b6020020181815250506001600282600160028110611e8757fe5b602002015181611e9357fe5b061415611ee55780600160028110611ea757fe5b60200201517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f0381600160028110611edb57fe5b6020020181815250505b919050565b60007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80611f1457fe5b82600160028110611f2157fe5b602002015183600160028110611f3357fe5b602002015109611f5383600060028110611f4957fe5b6020020151611524565b149050919050565b6000807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80611f8657fe5b8487097ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80611fb157fe5b848709809250819350505094509492505050565b60008060007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80611ff257fe5b878509905060007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8061202157fe5b87877ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f030990507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8061207057fe5b8183087ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8061209b57fe5b8689098094508195505050505094509492505050565b6040518060400160405280600290602082028038833980820191505090505090565b6040518060c00160405280600690602082028038833980820191505090505090565b6040518060200160405280600190602082028038833980820191505090505090565b604051806060016040528060039060208202803883398082019150509050509056fe4669727374206d756c7469706c69636174696f6e20636865636b206661696c65645365636f6e64206d756c7469706c69636174696f6e20636865636b206661696c6564a26469706673582212209743b07a578f79ffeaec76646ea1a9e7372520005d087b3daa410e4a04ff9c5864736f6c63430006020033"

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

// FieldHash is a free data retrieval call binding the contract method 0x80aa7713.
//
// Solidity: function fieldHash_(uint256 x) constant returns(uint256)
func (_VRFTestHelper *VRFTestHelperCaller) FieldHash(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _VRFTestHelper.contract.Call(opts, out, "fieldHash_", x)
	return *ret0, err
}

// FieldHash is a free data retrieval call binding the contract method 0x80aa7713.
//
// Solidity: function fieldHash_(uint256 x) constant returns(uint256)
func (_VRFTestHelper *VRFTestHelperSession) FieldHash(x *big.Int) (*big.Int, error) {
	return _VRFTestHelper.Contract.FieldHash(&_VRFTestHelper.CallOpts, x)
}

// FieldHash is a free data retrieval call binding the contract method 0x80aa7713.
//
// Solidity: function fieldHash_(uint256 x) constant returns(uint256)
func (_VRFTestHelper *VRFTestHelperCallerSession) FieldHash(x *big.Int) (*big.Int, error) {
	return _VRFTestHelper.Contract.FieldHash(&_VRFTestHelper.CallOpts, x)
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
