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
const VRFTestHelperABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"p1\",\"type\":\"uint256[2]\"},{\"name\":\"p2\",\"type\":\"uint256[2]\"},{\"name\":\"invZ\",\"type\":\"uint256\"}],\"name\":\"affineECAdd_\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256[2]\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"zqHash_\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"hashToCurve_\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256[2]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"hash\",\"type\":\"uint256[2]\"},{\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"name\":\"gamma\",\"type\":\"uint256[2]\"},{\"name\":\"uWitness\",\"type\":\"address\"},{\"name\":\"v\",\"type\":\"uint256[2]\"}],\"name\":\"scalarFromCurve_\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"base\",\"type\":\"uint256\"},{\"name\":\"exponent\",\"type\":\"uint256\"}],\"name\":\"bigModExp_\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"squareRoot_\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"c\",\"type\":\"uint256\"},{\"name\":\"p\",\"type\":\"uint256[2]\"},{\"name\":\"s\",\"type\":\"uint256\"},{\"name\":\"lcWitness\",\"type\":\"address\"}],\"name\":\"verifyLinearCombinationWithGenerator_\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"px\",\"type\":\"uint256\"},{\"name\":\"py\",\"type\":\"uint256\"},{\"name\":\"qx\",\"type\":\"uint256\"},{\"name\":\"qy\",\"type\":\"uint256\"}],\"name\":\"projectiveECAdd_\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"ySquared_\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"x\",\"type\":\"uint256[2]\"},{\"name\":\"scalar\",\"type\":\"uint256\"},{\"name\":\"q\",\"type\":\"uint256[2]\"}],\"name\":\"ecmulVerify_\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"name\":\"gamma\",\"type\":\"uint256[2]\"},{\"name\":\"c\",\"type\":\"uint256\"},{\"name\":\"s\",\"type\":\"uint256\"},{\"name\":\"seed\",\"type\":\"uint256\"},{\"name\":\"uWitness\",\"type\":\"address\"},{\"name\":\"cGammaWitness\",\"type\":\"uint256[2]\"},{\"name\":\"sHashWitness\",\"type\":\"uint256[2]\"},{\"name\":\"zInv\",\"type\":\"uint256\"}],\"name\":\"verifyVRFProof_\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"proof\",\"type\":\"bytes\"}],\"name\":\"randomValueFromVRFProof\",\"outputs\":[{\"name\":\"output\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"c\",\"type\":\"uint256\"},{\"name\":\"p1\",\"type\":\"uint256[2]\"},{\"name\":\"cp1Witness\",\"type\":\"uint256[2]\"},{\"name\":\"s\",\"type\":\"uint256\"},{\"name\":\"p2\",\"type\":\"uint256[2]\"},{\"name\":\"sp2Witness\",\"type\":\"uint256[2]\"},{\"name\":\"zInv\",\"type\":\"uint256\"}],\"name\":\"linearCombination_\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256[2]\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"}]"

// VRFTestHelperBin is the compiled bytecode used for deploying new contracts.
var VRFTestHelperBin = "0x608060405234801561001057600080fd5b50612395806100206000396000f3fe6080604052600436106100c5576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063244f896d146100ca57806324d72ea9146101c35780633545245014610212578063525413cf146102ca5780635de60042146104345780638af046ea1461048d57806391d5f691146104dc57806395e6ee921461059a5780639d6f033714610615578063aa7b2fbb14610664578063ef3b10ec14610739578063fa8fc6f1146108b7578063fe54f2a214610993575b600080fd5b3480156100d657600080fd5b50610185600480360360a08110156100ed57600080fd5b8101908080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f820116905080830192505050505050919291929080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f820116905080830192505050505050919291929080359060200190929190505050610b23565b6040518082600260200280838360005b838110156101b0578082015181840152602081019050610195565b5050505090500191505060405180910390f35b3480156101cf57600080fd5b506101fc600480360360208110156101e657600080fd5b8101908080359060200190929190505050610b3f565b6040518082815260200191505060405180910390f35b34801561021e57600080fd5b5061028c6004803603606081101561023557600080fd5b8101908080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f820116905080830192505050505050919291929080359060200190929190505050610b51565b6040518082600260200280838360005b838110156102b757808201518184015260208101905061029c565b5050505090500191505060405180910390f35b3480156102d657600080fd5b5061041e60048036036101208110156102ee57600080fd5b8101908080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f820116905080830192505050505050919291929080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f820116905080830192505050505050919291929080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f8201169050808301925050505050509192919290803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f8201169050808301925050505050509192919290505050610b6b565b6040518082815260200191505060405180910390f35b34801561044057600080fd5b506104776004803603604081101561045757600080fd5b810190808035906020019092919080359060200190929190505050610b85565b6040518082815260200191505060405180910390f35b34801561049957600080fd5b506104c6600480360360208110156104b057600080fd5b8101908080359060200190929190505050610b99565b6040518082815260200191505060405180910390f35b3480156104e857600080fd5b50610580600480360360a08110156104ff57600080fd5b810190808035906020019092919080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f820116905080830192505050505050919291929080359060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610bab565b604051808215151515815260200191505060405180910390f35b3480156105a657600080fd5b506105f1600480360360808110156105bd57600080fd5b8101908080359060200190929190803590602001909291908035906020019092919080359060200190929190505050610bc3565b60405180848152602001838152602001828152602001935050505060405180910390f35b34801561062157600080fd5b5061064e6004803603602081101561063857600080fd5b8101908080359060200190929190505050610be4565b6040518082815260200191505060405180910390f35b34801561067057600080fd5b5061071f600480360360a081101561068757600080fd5b8101908080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f82011690508083019250505050505091929192908035906020019092919080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f8201169050808301925050505050509192919290505050610bf6565b604051808215151515815260200191505060405180910390f35b34801561074557600080fd5b506108b560048036036101a081101561075d57600080fd5b8101908080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f820116905080830192505050505050919291929080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f8201169050808301925050505050509192919290803590602001909291908035906020019092919080359060200190929190803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f820116905080830192505050505050919291929080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f820116905080830192505050505050919291929080359060200190929190505050610c0c565b005b3480156108c357600080fd5b5061097d600480360360208110156108da57600080fd5b81019080803590602001906401000000008111156108f757600080fd5b82018360208201111561090957600080fd5b8035906020019184600183028401116401000000008311171561092b57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290505050610c28565b6040518082815260200191505060405180910390f35b34801561099f57600080fd5b50610ae560048036036101608110156109b757600080fd5b810190808035906020019092919080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f820116905080830192505050505050919291929080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f82011690508083019250505050505091929192908035906020019092919080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f820116905080830192505050505050919291929080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f820116905080830192505050505050919291929080359060200190929190505050610df8565b6040518082600260200280838360005b83811015610b10578082015181840152602081019050610af5565b5050505090500191505060405180910390f35b610b2b6122de565b610b36848484610e1c565b90509392505050565b6000610b4a82610faa565b9050919050565b610b596122de565b610b638383611010565b905092915050565b6000610b7a8686868686611262565b905095945050505050565b6000610b918383611395565b905092915050565b6000610ba482611505565b9050919050565b6000610bb985858585611547565b9050949350505050565b6000806000610bd48787878761178f565b9250925092509450945094915050565b6000610bef8261196f565b9050919050565b6000610c03848484611a03565b90509392505050565b610c1d898989898989898989611b9c565b505050505050505050565b60006101a08251141515610ca4576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260128152602001807f77726f6e672070726f6f66206c656e677468000000000000000000000000000081525060200191505060405180910390fd5b610cac6122de565b610cb46122de565b610cbc612300565b6000610cc66122de565b610cce6122de565b6000888060200190516101a0811015610ce657600080fd5b810190809190826040019190826040019190826060018051906020019092919091908260400191908260400180519060200190929190505050869650859550849450839350829250819150809750819850829950839a50849b50859c50869d5050505050505050610d978787876000600381101515610d6157fe5b6020020151886001600381101515610d7557fe5b6020020151896002600381101515610d8957fe5b602002015189898989611b9c565b856040516020018082600260200280838360005b83811015610dc6578082015181840152602081019050610dab565b505050509050019150506040516020818303038152906040528051906020012060019004975050505050505050919050565b610e006122de565b610e0f88888888888888611ecd565b9050979650505050505050565b610e246122de565b6000806000610e81876000600281101515610e3b57fe5b6020020151886001600281101515610e4f57fe5b6020020151886000600281101515610e6357fe5b6020020151896001600281101515610e7757fe5b602002015161178f565b80935081945082955050505060017ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f801515610eb957fe5b868309141515610f31576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260198152602001807f696e765a206d75737420626520696e7665727365206f66207a0000000000000081525060200191505060405180910390fd5b60408051908101604052807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f801515610f6657fe5b87860981526020017ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f801515610f9857fe5b87850981525093505050509392505050565b60008190505b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8110151561100b57806040516020018082815260200191505060405160208183030381529060405280519060200120600190049050610fb0565b919050565b6110186122de565b61107b83836040516020018083600260200280838360005b8381101561104b578082015181840152602081019050611030565b50505050905001828152602001925050506040516020818303038152906040528051906020012060019004610faa565b81600060028110151561108a57fe5b6020020181815250506110b76110b28260006002811015156110a857fe5b602002015161196f565b611505565b8160016002811015156110c657fe5b6020020181815250505b6110ec8160006002811015156110e257fe5b602002015161196f565b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80151561111657fe5b82600160028110151561112557fe5b602002015183600160028110151561113957fe5b6020020151091415156111e55761118c81600060028110151561115857fe5b6020020151604051602001808281526020019150506040516020818303038152906040528051906020012060019004610faa565b81600060028110151561119b57fe5b6020020181815250506111c86111c38260006002811015156111b957fe5b602002015161196f565b611505565b8160016002811015156111d757fe5b6020020181815250506110d0565b600160028260016002811015156111f857fe5b602002015181151561120657fe5b06141561125c5780600160028110151561121c57fe5b60200201517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f0381600160028110151561125257fe5b6020020181815250505b92915050565b600085858584866040516020018086600260200280838360005b8381101561129757808201518184015260208101905061127c565b5050505090500185600260200280838360005b838110156112c55780820151818401526020810190506112aa565b5050505090500184600260200280838360005b838110156112f35780820151818401526020810190506112d8565b5050505090500183600260200280838360005b83811015611321578082015181840152602081019050611306565b505050509050018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166c01000000000000000000000000028152601401955050505050506040516020818303038152906040528051906020012060019004905095945050505050565b6000806113a0612323565b60208160006006811015156113b157fe5b60200201818152505060208160016006811015156113cb57fe5b60200201818152505060208160026006811015156113e557fe5b602002018181525050848160036006811015156113fe57fe5b6020020181815250508381600460068110151561141757fe5b6020020181815250507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f81600560068110151561145057fe5b602002018181525050611461612346565b60208160c0846005600019fa925060008314156114e6576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260128152602001807f6269674d6f64457870206661696c75726521000000000000000000000000000081525060200191505060405180910390fd5b8060006001811015156114f557fe5b6020020151935050505092915050565b600061154082600260017ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f01908060020a8204915050611395565b9050919050565b60008073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16141515156115ed576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600b8152602001807f626164207769746e65737300000000000000000000000000000000000000000081525060200191505060405180910390fd5b60007ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd036414180151561161957fe5b8486600060028110151561162957fe5b6020020151097ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141036001029050600080600287600160028110151561166a57fe5b602002015181151561167857fe5b061461168557601c611688565b601b5b905060007ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd03641418015156116b657fe5b8760006002811015156116c557fe5b6020020151890960010290506000600184848a60006002811015156116e657fe5b60200201516001028560405160008152602001604052604051808581526020018460ff1660ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa158015611746573d6000803e3d6000fd5b5050506020604051035190508573ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614945050505050949350505050565b60008060008060006001809150915060007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8015156117ca57fe5b897ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f038808905060007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80151561181d57fe5b8b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f038a089050600061185283838585612103565b809250819950505061186688828e88612171565b809250819950505061187a88828c87612171565b809250819950505060006118908d878b85612171565b80925081995050506118a488828686612103565b80925081995050506118b888828e89612171565b8092508199505050808214151561195b577ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8015156118f357fe5b818a0998507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80151561192257fe5b82890997507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80151561195157fe5b818309965061195f565b8196505b5050505050509450945094915050565b6000807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80151561199c57fe5b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8015156119c657fe5b848509840990507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8015156119f757fe5b60078208915050919050565b6000808314151515611a1457600080fd5b60007ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141801515611a4057fe5b856000600281101515611a4f57fe5b6020020151850960010290506000806002876001600281101515611a6f57fe5b6020020151811515611a7d57fe5b061415611a8b57601b611a8e565b601c5b9050836040516020018082600260200280838360005b83811015611abf578082015181840152602081019050611aa4565b50505050905001915050604051602081830303815290604052805190602001206001900473ffffffffffffffffffffffffffffffffffffffff166001600060010283896000600281101515611b1057fe5b60200201516001028660405160008152602001604052604051808581526020018460ff1660ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa158015611b70573d6000803e3d6000fd5b5050506020604051035173ffffffffffffffffffffffffffffffffffffffff1614925050509392505050565b611ba589612265565b1515611c19576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601a8152602001807f7075626c6963206b6579206973206e6f74206f6e20637572766500000000000081525060200191505060405180910390fd5b611c2288612265565b1515611c96576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260158152602001807f67616d6d61206973206e6f74206f6e206375727665000000000000000000000081525060200191505060405180910390fd5b611c9f83612265565b1515611d13576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601d8152602001807f6347616d6d615769746e657373206973206e6f74206f6e20637572766500000081525060200191505060405180910390fd5b611d1c82612265565b1515611d90576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601c8152602001807f73486173685769746e657373206973206e6f74206f6e2063757276650000000081525060200191505060405180910390fd5b611d9c878a8887611547565b1515611e10576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601a8152602001807f6164647228632a706b2b732a6729e289a05f755769746e65737300000000000081525060200191505060405180910390fd5b611e186122de565b611e228a87611010565b9050611e2c6122de565b611e3b898b878b868989611ecd565b9050611e4a828c8c8985611262565b89141515611ec0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600d8152602001807f696e76616c69642070726f6f660000000000000000000000000000000000000081525060200191505060405180910390fd5b5050505050505050505050565b611ed56122de565b60007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f846000600281101515611f0757fe5b6020020151886000600281101515611f1b57fe5b602002015103811515611f2a57fe5b0614151515611fa1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f706f696e747320696e2073756d206d7573742062652064697374696e6374000081525060200191505060405180910390fd5b611fac878988611a03565b1515612046576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260218152602001807f4669727374206d756c7469706c69636174696f6e20636865636b206661696c6581526020017f640000000000000000000000000000000000000000000000000000000000000081525060400191505060405180910390fd5b612051848685611a03565b15156120eb576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260228152602001807f5365636f6e64206d756c7469706c69636174696f6e20636865636b206661696c81526020017f656400000000000000000000000000000000000000000000000000000000000081525060400191505060405180910390fd5b6120f6868484610e1c565b9050979650505050505050565b6000807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80151561213057fe5b8487097ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80151561215d57fe5b848709809250819350505094509492505050565b60008060007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8015156121a057fe5b878509905060007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8015156121d157fe5b87877ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f030990507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80151561222257fe5b8183087ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80151561224f57fe5b8689098094508195505050505094509492505050565b60007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80151561229157fe5b8260016002811015156122a057fe5b60200201518360016002811015156122b457fe5b6020020151096122d68360006002811015156122cc57fe5b602002015161196f565b149050919050565b6040805190810160405280600290602082028038833980820191505090505090565b606060405190810160405280600390602082028038833980820191505090505090565b60c060405190810160405280600690602082028038833980820191505090505090565b60206040519081016040528060019060208202803883398082019150509050509056fea165627a7a7230582038c2d8647f813f0661f0ba200b3daad4067cc48b51894cf9cd86747833d4901f0029"

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

// RandomValueFromVRFProof is a free data retrieval call binding the contract method 0xfa8fc6f1.
//
// Solidity: function randomValueFromVRFProof(bytes proof) constant returns(uint256 output)
func (_VRFTestHelper *VRFTestHelperCaller) RandomValueFromVRFProof(opts *bind.CallOpts, proof []byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _VRFTestHelper.contract.Call(opts, out, "randomValueFromVRFProof", proof)
	return *ret0, err
}

// RandomValueFromVRFProof is a free data retrieval call binding the contract method 0xfa8fc6f1.
//
// Solidity: function randomValueFromVRFProof(bytes proof) constant returns(uint256 output)
func (_VRFTestHelper *VRFTestHelperSession) RandomValueFromVRFProof(proof []byte) (*big.Int, error) {
	return _VRFTestHelper.Contract.RandomValueFromVRFProof(&_VRFTestHelper.CallOpts, proof)
}

// RandomValueFromVRFProof is a free data retrieval call binding the contract method 0xfa8fc6f1.
//
// Solidity: function randomValueFromVRFProof(bytes proof) constant returns(uint256 output)
func (_VRFTestHelper *VRFTestHelperCallerSession) RandomValueFromVRFProof(proof []byte) (*big.Int, error) {
	return _VRFTestHelper.Contract.RandomValueFromVRFProof(&_VRFTestHelper.CallOpts, proof)
}

// ScalarFromCurve is a free data retrieval call binding the contract method 0x525413cf.
//
// Solidity: function scalarFromCurve_(uint256[2] hash, uint256[2] pk, uint256[2] gamma, address uWitness, uint256[2] v) constant returns(uint256)
func (_VRFTestHelper *VRFTestHelperCaller) ScalarFromCurve(opts *bind.CallOpts, hash [2]*big.Int, pk [2]*big.Int, gamma [2]*big.Int, uWitness common.Address, v [2]*big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _VRFTestHelper.contract.Call(opts, out, "scalarFromCurve_", hash, pk, gamma, uWitness, v)
	return *ret0, err
}

// ScalarFromCurve is a free data retrieval call binding the contract method 0x525413cf.
//
// Solidity: function scalarFromCurve_(uint256[2] hash, uint256[2] pk, uint256[2] gamma, address uWitness, uint256[2] v) constant returns(uint256)
func (_VRFTestHelper *VRFTestHelperSession) ScalarFromCurve(hash [2]*big.Int, pk [2]*big.Int, gamma [2]*big.Int, uWitness common.Address, v [2]*big.Int) (*big.Int, error) {
	return _VRFTestHelper.Contract.ScalarFromCurve(&_VRFTestHelper.CallOpts, hash, pk, gamma, uWitness, v)
}

// ScalarFromCurve is a free data retrieval call binding the contract method 0x525413cf.
//
// Solidity: function scalarFromCurve_(uint256[2] hash, uint256[2] pk, uint256[2] gamma, address uWitness, uint256[2] v) constant returns(uint256)
func (_VRFTestHelper *VRFTestHelperCallerSession) ScalarFromCurve(hash [2]*big.Int, pk [2]*big.Int, gamma [2]*big.Int, uWitness common.Address, v [2]*big.Int) (*big.Int, error) {
	return _VRFTestHelper.Contract.ScalarFromCurve(&_VRFTestHelper.CallOpts, hash, pk, gamma, uWitness, v)
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

// ZqHash is a free data retrieval call binding the contract method 0x24d72ea9.
//
// Solidity: function zqHash_(uint256 x) constant returns(uint256)
func (_VRFTestHelper *VRFTestHelperCaller) ZqHash(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _VRFTestHelper.contract.Call(opts, out, "zqHash_", x)
	return *ret0, err
}

// ZqHash is a free data retrieval call binding the contract method 0x24d72ea9.
//
// Solidity: function zqHash_(uint256 x) constant returns(uint256)
func (_VRFTestHelper *VRFTestHelperSession) ZqHash(x *big.Int) (*big.Int, error) {
	return _VRFTestHelper.Contract.ZqHash(&_VRFTestHelper.CallOpts, x)
}

// ZqHash is a free data retrieval call binding the contract method 0x24d72ea9.
//
// Solidity: function zqHash_(uint256 x) constant returns(uint256)
func (_VRFTestHelper *VRFTestHelperCallerSession) ZqHash(x *big.Int) (*big.Int, error) {
	return _VRFTestHelper.Contract.ZqHash(&_VRFTestHelper.CallOpts, x)
}
