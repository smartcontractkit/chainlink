// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package solidity_vrf_coordinator_interface

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

// VRFCoordinatorABI is the input ABI used to generate the binding from.
const VRFCoordinatorABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"NewServiceAgreement\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"jobID\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"RandomnessRequest\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"callbacks\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"callbackContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"randomnessFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_proof\",\"type\":\"bytes\"}],\"name\":\"fulfillRandomnessRequest\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"_publicKey\",\"type\":\"uint256[2]\"}],\"name\":\"hashOfKey\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"},{\"internalType\":\"uint256[2]\",\"name\":\"_publicProvingKey\",\"type\":\"uint256[2]\"},{\"internalType\":\"bytes32\",\"name\":\"_jobID\",\"type\":\"bytes32\"}],\"name\":\"registerProvingKey\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"serviceAgreements\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"vRFOracle\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"jobID\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"withdrawableTokens\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// VRFCoordinatorBin is the compiled bytecode used for deploying new contracts.
var VRFCoordinatorBin = "0x608060405234801561001057600080fd5b506040516129273803806129278339818101604052602081101561003357600080fd5b8101908080519060200190929190505050806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050612893806100946000396000f3fe608060405234801561001057600080fd5b50600436106100875760003560e01c806375d350701161005b57806375d35070146102c3578063a4c0ed361461033f578063caf70c4a14610424578063f3fef3a31461049d57610087565b80626f6ad01461008c57806321f36509146100e45780635e1c1059146101605780636815851e14610233575b600080fd5b6100ce600480360360208110156100a257600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506104eb565b6040518082815260200191505060405180910390f35b610110600480360360208110156100fa57600080fd5b8101908080359060200190929190505050610503565b604051808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001838152602001828152602001935050505060405180910390f35b6102196004803603602081101561017657600080fd5b810190808035906020019064010000000081111561019357600080fd5b8201836020820111156101a557600080fd5b803590602001918460018302840111640100000000831117156101c757600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050919291929050505061054d565b604051808215151515815260200191505060405180910390f35b6102736004803603608081101561024957600080fd5b810190808035906020019092919080604001909192919290803590602001909291905050506108ec565b604051808481526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001828152602001935050505060405180910390f35b6102ef600480360360208110156102d957600080fd5b8101908080359060200190929190505050610b23565b604051808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001838152602001828152602001935050505060405180910390f35b6104226004803603606081101561035557600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291908035906020019064010000000081111561039c57600080fd5b8201836020820111156103ae57600080fd5b803590602001918460018302840111640100000000831117156103d057600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290505050610b6d565b005b6104876004803603604081101561043a57600080fd5b8101908080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f8201169050808301925050505050509192919290505050610c79565b6040518082815260200191505060405180910390f35b6104e9600480360360408110156104b357600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050610cd1565b005b60036020528060005260406000206000915090505481565b60016020528060005260406000206000915090508060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060010154908060020154905083565b600061055761275b565b600060208401915060e08401519050600061057183610c79565b9050600061057f8284610ec2565b905061058961277d565b600160008381526020019081526020016000206040518060600160405290816000820160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001600182015481526020016002820154815250509050600073ffffffffffffffffffffffffffffffffffffffff16816000015173ffffffffffffffffffffffffffffffffffffffff1614156106ba576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260188152602001807f6e6f20636f72726573706f6e64696e672072657175657374000000000000000081525060200191505060405180910390fd5b60006106c588610efb565b90508160200151600360006002600088815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254019250508190555060006060631f1f897f60e01b85846040516024018083815260200182815260200192505050604051602081830303815290604052907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff838183161783525050505090506000846000015173ffffffffffffffffffffffffffffffffffffffff16826040518082805190602001908083835b602083106108235780518252602082019150602081019050602083039250610800565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d8060008114610885576040519150601f19603f3d011682016040523d82523d6000602084013e61088a565b606091505b5050905060016000878152602001908152602001600020600080820160006101000a81549073ffffffffffffffffffffffffffffffffffffffff0219169055600182016000905560028201600090555050809950505050505050505050919050565b6000806000610931856002806020026040519081016040528092919082600260200280828437600081840152601f19601f820116905080830192505050505050610c79565b925060006002600085815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614610a10576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260198152602001807f706c656173652072656769737465722061206e6577206b65790000000000000081525060200191505060405180910390fd5b336002600086815260200190815260200160002060000160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508460026000868152602001908152602001600020600101819055508660026000868152602001908152602001600020600201819055507fae189157e0628c1e62315e9179156e1ea10e90e9c15060002f7021e907dc2cfe8488604051808381526020018281526020019250505060405180910390a1836002600086815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16889350935093505093509350939050565b60026020528060005260406000206000915090508060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060010154908060020154905083565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610c2f576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260138152602001807f4d75737420757365204c494e4b20746f6b656e0000000000000000000000000081525060200191505060405180910390fd5b600080828060200190516040811015610c4757600080fd5b81019080805190602001909291908051906020019092919050505091509150610c72828286886110c2565b5050505050565b6000816040516020018082600260200280838360005b83811015610caa578082015181840152602081019050610c8f565b50505050905001915050604051602081830303815290604052805190602001209050919050565b8080600360003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020541015610d87576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260208152602001807f63616e2774207769746864726177206d6f7265207468616e2062616c616e636581525060200191505060405180910390fd5b81600360003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825403925050819055506000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb84846040518363ffffffff1660e01b8152600401808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b158015610e7c57600080fd5b505af1158015610e90573d6000803e3d6000fd5b505050506040513d6020811015610ea657600080fd5b8101908080519060200190929190505050610ebd57fe5b505050565b60008282604051602001808381526020018281526020019250505060405160208183030381529060405280519060200120905092915050565b60006101a0825114610f75576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260128152602001807f77726f6e672070726f6f66206c656e677468000000000000000000000000000081525060200191505060405180910390fd5b610f7d61275b565b610f8561275b565b610f8d6127b4565b6000610f9761275b565b610f9f61275b565b6000888060200190516101a0811015610fb757600080fd5b810190809190826040019190826040019190826060018051906020019092919091908260400191908260400180519060200190929190505050869650859550849450839350829250819150809750819850829950839a50849b50859c50869d505050505050505061106287878760006003811061103057fe5b60200201518860016003811061104257fe5b60200201518960026003811061105457fe5b6020020151898989896113b6565b856040516020018082600260200280838360005b83811015611091578082015181840152602081019050611076565b505050509050019150506040516020818303038152906040528051906020012060001c975050505050505050919050565b81846002600082815260200190815260200160002060020154821015611150576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f42656c6f7720616772656564207061796d656e7400000000000000000000000081525060200191505060405180910390fd5b60006004600088815260200190815260200160002060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905060006111b3888887856116e1565b905060006111c18983610ec2565b9050600073ffffffffffffffffffffffffffffffffffffffff166001600083815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff161461122f57fe5b856001600083815260200190815260200160002060000160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550866001600083815260200190815260200160002060010181905550816001600083815260200190815260200160002060020181905550600260008a8152602001908152602001600020600101547fd241d78a52145a5d1d1ff002e32ec15cdc395631bcee66246650c2429dfaccaa8a84898b604051808581526020018481526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200182815260200194505050505060405180910390a26001600460008b815260200190815260200160002060008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282540192505081905550505050505050505050565b6113bf8961175b565b611431576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601a8152602001807f7075626c6963206b6579206973206e6f74206f6e20637572766500000000000081525060200191505060405180910390fd5b61143a8861175b565b6114ac576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260158152602001807f67616d6d61206973206e6f74206f6e206375727665000000000000000000000081525060200191505060405180910390fd5b6114b58361175b565b611527576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601d8152602001807f6347616d6d615769746e657373206973206e6f74206f6e20637572766500000081525060200191505060405180910390fd5b6115308261175b565b6115a2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601c8152602001807f73486173685769746e657373206973206e6f74206f6e2063757276650000000081525060200191505060405180910390fd5b6115ae878a88876117cc565b611620576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601a8152602001807f6164647228632a706b2b732a6729e289a05f755769746e65737300000000000081525060200191505060405180910390fd5b61162861275b565b6116328a87611a04565b905061163c61275b565b61164b898b878b868989611aba565b9050600061165c838d8d8a86611c5e565b9050808a146116d3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600d8152602001807f696e76616c69642070726f6f660000000000000000000000000000000000000081525060200191505060405180910390fd5b505050505050505050505050565b600084848484604051602001808581526020018481526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019450505050506040516020818303038152906040528051906020012060001c9050949350505050565b60007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8061178557fe5b8260016002811061179257fe5b6020020151836001600281106117a457fe5b6020020151096117c4836000600281106117ba57fe5b6020020151611d84565b149050919050565b60008073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415611870576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600b8152602001807f626164207769746e65737300000000000000000000000000000000000000000081525060200191505060405180910390fd5b60008060028660016002811061188257fe5b60200201518161188e57fe5b061461189b57601c61189e565b601b5b905060007ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141806118ca57fe5b85876000600281106118d857fe5b6020020151097ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd03641410360001b905060007ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd03641418061192f57fe5b8760006002811061193c57fe5b6020020151890960001b90506000600183858a60006002811061195b57fe5b602002015160001b8560405160008152602001604052604051808581526020018460ff1660ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa1580156119bb573d6000803e3d6000fd5b5050506020604051035190508573ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614945050505050949350505050565b611a0c61275b565b611a6483836040516020018083600260200280838360005b83811015611a3f578082015181840152602081019050611a24565b5050505090500182815260200192505050604051602081830303815290604052611e12565b90505b611a708161175b565b611ab457611aad81600060028110611a8457fe5b602002015160405160200180828152602001915050604051602081830303815290604052611e12565b9050611a67565b92915050565b611ac261275b565b60007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f84600060028110611af257fe5b602002015188600060028110611b0457fe5b60200201510381611b1157fe5b061415611b86576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f706f696e747320696e2073756d206d7573742062652064697374696e6374000081525060200191505060405180910390fd5b611b91878988611eef565b611be6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602181526020018061281b6021913960400191505060405180910390fd5b611bf1848685611eef565b611c46576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602281526020018061283c6022913960400191505060405180910390fd5b611c518684846120b4565b9050979650505050505050565b600085858584866040516020018086600260200280838360005b83811015611c93578082015181840152602081019050611c78565b5050505090500185600260200280838360005b83811015611cc1578082015181840152602081019050611ca6565b5050505090500184600260200280838360005b83811015611cef578082015181840152602081019050611cd4565b5050505090500183600260200280838360005b83811015611d1d578082015181840152602081019050611d02565b505050509050018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1660601b8152601401955050505050506040516020818303038152906040528051906020012060001c905095945050505050565b6000807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80611daf57fe5b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80611dd757fe5b848509840990507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80611e0657fe5b60078208915050919050565b611e1a61275b565b611e2d828051906020012060001c612232565b81600060028110611e3a57fe5b602002018181525050611e65611e6082600060028110611e5657fe5b6020020151611d84565b612295565b81600160028110611e7257fe5b6020020181815250506001600282600160028110611e8c57fe5b602002015181611e9857fe5b061415611eea5780600160028110611eac57fe5b60200201517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f0381600160028110611ee057fe5b6020020181815250505b919050565b600080831415611efe57600080fd5b600084600060028110611f0d57fe5b60200201519050600080600287600160028110611f2657fe5b602002015181611f3257fe5b0614611f3f57601c611f42565b601b5b905060007ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd036414180611f6e57fe5b83870960001b9050600060016000801b848660001b8560405160008152602001604052604051808581526020018460ff1660ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa158015611fdb573d6000803e3d6000fd5b5050506020604051035190506000866040516020018082600260200280838360005b83811015612018578082015181840152602081019050611ffd565b505050509050019150506040516020818303038152906040528051906020012060001c9050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16141561207757600080fd5b8073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614955050505050509392505050565b6120bc61275b565b6000806000612111876000600281106120d157fe5b6020020151886001600281106120e357fe5b6020020151886000600281106120f557fe5b60200201518960016002811061210757fe5b60200201516122cf565b80935081945082955050505060017ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8061214757fe5b868309146121bd576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260198152602001807f696e765a206d75737420626520696e7665727365206f66207a0000000000000081525060200191505060405180910390fd5b60405180604001604052807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f806121f057fe5b87860981526020017ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8061222057fe5b87850981525093505050509392505050565b60008190505b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f81106122905780604051602001808281526020019150506040516020818303038152906040528051906020012060001c9050612238565b919050565b60006122c882600260017ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f01901c6124a3565b9050919050565b60008060008060006001809150915060007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8061230857fe5b897ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f038808905060007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8061235957fe5b8b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f038a089050600061238e83838585612605565b80925081995050506123a288828e8861266f565b80925081995050506123b688828c8761266f565b809250819950505060006123cc8d878b8561266f565b80925081995050506123e088828686612605565b80925081995050506123f488828e8961266f565b809250819950505080821461248f577ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8061242b57fe5b818a0998507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8061245857fe5b82890997507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8061248557fe5b8183099650612493565b8196505b5050505050509450945094915050565b6000806124ae6127d6565b6020816000600681106124bd57fe5b6020020181815250506020816001600681106124d557fe5b6020020181815250506020816002600681106124ed57fe5b602002018181525050848160036006811061250457fe5b602002018181525050838160046006811061251b57fe5b6020020181815250507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8160056006811061255257fe5b6020020181815250506125636127f8565b60208160c0846005600019fa925060008314156125e8576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260128152602001807f6269674d6f64457870206661696c75726521000000000000000000000000000081525060200191505060405180910390fd5b806000600181106125f557fe5b6020020151935050505092915050565b6000807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8061263057fe5b8487097ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8061265b57fe5b848709809250819350505094509492505050565b60008060007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8061269c57fe5b878509905060007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f806126cb57fe5b87877ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f030990507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8061271a57fe5b8183087ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8061274557fe5b8689098094508195505050505094509492505050565b6040518060400160405280600290602082028038833980820191505090505090565b6040518060600160405280600073ffffffffffffffffffffffffffffffffffffffff16815260200160008152602001600081525090565b6040518060600160405280600390602082028038833980820191505090505090565b6040518060c00160405280600690602082028038833980820191505090505090565b604051806020016040528060019060208202803883398082019150509050509056fe4669727374206d756c7469706c69636174696f6e20636865636b206661696c65645365636f6e64206d756c7469706c69636174696f6e20636865636b206661696c6564a264697066735822122007131c4d804cfd9f73e5c0a2ae0f8118cbf1e3e948873d22ce90f89357cd0b6464736f6c63430006020033"

// DeployVRFCoordinator deploys a new Ethereum contract, binding an instance of VRFCoordinator to it.
func DeployVRFCoordinator(auth *bind.TransactOpts, backend bind.ContractBackend, _link common.Address) (common.Address, *types.Transaction, *VRFCoordinator, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFCoordinatorABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(VRFCoordinatorBin), backend, _link)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFCoordinator{VRFCoordinatorCaller: VRFCoordinatorCaller{contract: contract}, VRFCoordinatorTransactor: VRFCoordinatorTransactor{contract: contract}, VRFCoordinatorFilterer: VRFCoordinatorFilterer{contract: contract}}, nil
}

// VRFCoordinator is an auto generated Go binding around an Ethereum contract.
type VRFCoordinator struct {
	VRFCoordinatorCaller     // Read-only binding to the contract
	VRFCoordinatorTransactor // Write-only binding to the contract
	VRFCoordinatorFilterer   // Log filterer for contract events
}

// VRFCoordinatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type VRFCoordinatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFCoordinatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VRFCoordinatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFCoordinatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VRFCoordinatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFCoordinatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VRFCoordinatorSession struct {
	Contract     *VRFCoordinator   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VRFCoordinatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VRFCoordinatorCallerSession struct {
	Contract *VRFCoordinatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// VRFCoordinatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VRFCoordinatorTransactorSession struct {
	Contract     *VRFCoordinatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// VRFCoordinatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type VRFCoordinatorRaw struct {
	Contract *VRFCoordinator // Generic contract binding to access the raw methods on
}

// VRFCoordinatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VRFCoordinatorCallerRaw struct {
	Contract *VRFCoordinatorCaller // Generic read-only contract binding to access the raw methods on
}

// VRFCoordinatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VRFCoordinatorTransactorRaw struct {
	Contract *VRFCoordinatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVRFCoordinator creates a new instance of VRFCoordinator, bound to a specific deployed contract.
func NewVRFCoordinator(address common.Address, backend bind.ContractBackend) (*VRFCoordinator, error) {
	contract, err := bindVRFCoordinator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinator{VRFCoordinatorCaller: VRFCoordinatorCaller{contract: contract}, VRFCoordinatorTransactor: VRFCoordinatorTransactor{contract: contract}, VRFCoordinatorFilterer: VRFCoordinatorFilterer{contract: contract}}, nil
}

// NewVRFCoordinatorCaller creates a new read-only instance of VRFCoordinator, bound to a specific deployed contract.
func NewVRFCoordinatorCaller(address common.Address, caller bind.ContractCaller) (*VRFCoordinatorCaller, error) {
	contract, err := bindVRFCoordinator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorCaller{contract: contract}, nil
}

// NewVRFCoordinatorTransactor creates a new write-only instance of VRFCoordinator, bound to a specific deployed contract.
func NewVRFCoordinatorTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFCoordinatorTransactor, error) {
	contract, err := bindVRFCoordinator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTransactor{contract: contract}, nil
}

// NewVRFCoordinatorFilterer creates a new log filterer instance of VRFCoordinator, bound to a specific deployed contract.
func NewVRFCoordinatorFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFCoordinatorFilterer, error) {
	contract, err := bindVRFCoordinator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorFilterer{contract: contract}, nil
}

// bindVRFCoordinator binds a generic wrapper to an already deployed contract.
func bindVRFCoordinator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFCoordinatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VRFCoordinator *VRFCoordinatorRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _VRFCoordinator.Contract.VRFCoordinatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VRFCoordinator *VRFCoordinatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.VRFCoordinatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VRFCoordinator *VRFCoordinatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.VRFCoordinatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VRFCoordinator *VRFCoordinatorCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _VRFCoordinator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VRFCoordinator *VRFCoordinatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VRFCoordinator *VRFCoordinatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.contract.Transact(opts, method, params...)
}

// Callbacks is a free data retrieval call binding the contract method 0x21f36509.
//
// Solidity: function callbacks(bytes32 ) constant returns(address callbackContract, uint256 randomnessFee, uint256 seed)
func (_VRFCoordinator *VRFCoordinatorCaller) Callbacks(opts *bind.CallOpts, arg0 [32]byte) (struct {
	CallbackContract common.Address
	RandomnessFee    *big.Int
	Seed             *big.Int
}, error) {
	ret := new(struct {
		CallbackContract common.Address
		RandomnessFee    *big.Int
		Seed             *big.Int
	})
	out := ret
	err := _VRFCoordinator.contract.Call(opts, out, "callbacks", arg0)
	return *ret, err
}

// Callbacks is a free data retrieval call binding the contract method 0x21f36509.
//
// Solidity: function callbacks(bytes32 ) constant returns(address callbackContract, uint256 randomnessFee, uint256 seed)
func (_VRFCoordinator *VRFCoordinatorSession) Callbacks(arg0 [32]byte) (struct {
	CallbackContract common.Address
	RandomnessFee    *big.Int
	Seed             *big.Int
}, error) {
	return _VRFCoordinator.Contract.Callbacks(&_VRFCoordinator.CallOpts, arg0)
}

// Callbacks is a free data retrieval call binding the contract method 0x21f36509.
//
// Solidity: function callbacks(bytes32 ) constant returns(address callbackContract, uint256 randomnessFee, uint256 seed)
func (_VRFCoordinator *VRFCoordinatorCallerSession) Callbacks(arg0 [32]byte) (struct {
	CallbackContract common.Address
	RandomnessFee    *big.Int
	Seed             *big.Int
}, error) {
	return _VRFCoordinator.Contract.Callbacks(&_VRFCoordinator.CallOpts, arg0)
}

// HashOfKey is a free data retrieval call binding the contract method 0xcaf70c4a.
//
// Solidity: function hashOfKey(uint256[2] _publicKey) constant returns(bytes32)
func (_VRFCoordinator *VRFCoordinatorCaller) HashOfKey(opts *bind.CallOpts, _publicKey [2]*big.Int) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _VRFCoordinator.contract.Call(opts, out, "hashOfKey", _publicKey)
	return *ret0, err
}

// HashOfKey is a free data retrieval call binding the contract method 0xcaf70c4a.
//
// Solidity: function hashOfKey(uint256[2] _publicKey) constant returns(bytes32)
func (_VRFCoordinator *VRFCoordinatorSession) HashOfKey(_publicKey [2]*big.Int) ([32]byte, error) {
	return _VRFCoordinator.Contract.HashOfKey(&_VRFCoordinator.CallOpts, _publicKey)
}

// HashOfKey is a free data retrieval call binding the contract method 0xcaf70c4a.
//
// Solidity: function hashOfKey(uint256[2] _publicKey) constant returns(bytes32)
func (_VRFCoordinator *VRFCoordinatorCallerSession) HashOfKey(_publicKey [2]*big.Int) ([32]byte, error) {
	return _VRFCoordinator.Contract.HashOfKey(&_VRFCoordinator.CallOpts, _publicKey)
}

// ServiceAgreements is a free data retrieval call binding the contract method 0x75d35070.
//
// Solidity: function serviceAgreements(bytes32 ) constant returns(address vRFOracle, bytes32 jobID, uint256 fee)
func (_VRFCoordinator *VRFCoordinatorCaller) ServiceAgreements(opts *bind.CallOpts, arg0 [32]byte) (struct {
	VRFOracle common.Address
	JobID     [32]byte
	Fee       *big.Int
}, error) {
	ret := new(struct {
		VRFOracle common.Address
		JobID     [32]byte
		Fee       *big.Int
	})
	out := ret
	err := _VRFCoordinator.contract.Call(opts, out, "serviceAgreements", arg0)
	return *ret, err
}

// ServiceAgreements is a free data retrieval call binding the contract method 0x75d35070.
//
// Solidity: function serviceAgreements(bytes32 ) constant returns(address vRFOracle, bytes32 jobID, uint256 fee)
func (_VRFCoordinator *VRFCoordinatorSession) ServiceAgreements(arg0 [32]byte) (struct {
	VRFOracle common.Address
	JobID     [32]byte
	Fee       *big.Int
}, error) {
	return _VRFCoordinator.Contract.ServiceAgreements(&_VRFCoordinator.CallOpts, arg0)
}

// ServiceAgreements is a free data retrieval call binding the contract method 0x75d35070.
//
// Solidity: function serviceAgreements(bytes32 ) constant returns(address vRFOracle, bytes32 jobID, uint256 fee)
func (_VRFCoordinator *VRFCoordinatorCallerSession) ServiceAgreements(arg0 [32]byte) (struct {
	VRFOracle common.Address
	JobID     [32]byte
	Fee       *big.Int
}, error) {
	return _VRFCoordinator.Contract.ServiceAgreements(&_VRFCoordinator.CallOpts, arg0)
}

// WithdrawableTokens is a free data retrieval call binding the contract method 0x006f6ad0.
//
// Solidity: function withdrawableTokens(address ) constant returns(uint256)
func (_VRFCoordinator *VRFCoordinatorCaller) WithdrawableTokens(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _VRFCoordinator.contract.Call(opts, out, "withdrawableTokens", arg0)
	return *ret0, err
}

// WithdrawableTokens is a free data retrieval call binding the contract method 0x006f6ad0.
//
// Solidity: function withdrawableTokens(address ) constant returns(uint256)
func (_VRFCoordinator *VRFCoordinatorSession) WithdrawableTokens(arg0 common.Address) (*big.Int, error) {
	return _VRFCoordinator.Contract.WithdrawableTokens(&_VRFCoordinator.CallOpts, arg0)
}

// WithdrawableTokens is a free data retrieval call binding the contract method 0x006f6ad0.
//
// Solidity: function withdrawableTokens(address ) constant returns(uint256)
func (_VRFCoordinator *VRFCoordinatorCallerSession) WithdrawableTokens(arg0 common.Address) (*big.Int, error) {
	return _VRFCoordinator.Contract.WithdrawableTokens(&_VRFCoordinator.CallOpts, arg0)
}

// FulfillRandomnessRequest is a paid mutator transaction binding the contract method 0x5e1c1059.
//
// Solidity: function fulfillRandomnessRequest(bytes _proof) returns(bool)
func (_VRFCoordinator *VRFCoordinatorTransactor) FulfillRandomnessRequest(opts *bind.TransactOpts, _proof []byte) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "fulfillRandomnessRequest", _proof)
}

// FulfillRandomnessRequest is a paid mutator transaction binding the contract method 0x5e1c1059.
//
// Solidity: function fulfillRandomnessRequest(bytes _proof) returns(bool)
func (_VRFCoordinator *VRFCoordinatorSession) FulfillRandomnessRequest(_proof []byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.FulfillRandomnessRequest(&_VRFCoordinator.TransactOpts, _proof)
}

// FulfillRandomnessRequest is a paid mutator transaction binding the contract method 0x5e1c1059.
//
// Solidity: function fulfillRandomnessRequest(bytes _proof) returns(bool)
func (_VRFCoordinator *VRFCoordinatorTransactorSession) FulfillRandomnessRequest(_proof []byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.FulfillRandomnessRequest(&_VRFCoordinator.TransactOpts, _proof)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address _sender, uint256 _fee, bytes _data) returns()
func (_VRFCoordinator *VRFCoordinatorTransactor) OnTokenTransfer(opts *bind.TransactOpts, _sender common.Address, _fee *big.Int, _data []byte) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "onTokenTransfer", _sender, _fee, _data)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address _sender, uint256 _fee, bytes _data) returns()
func (_VRFCoordinator *VRFCoordinatorSession) OnTokenTransfer(_sender common.Address, _fee *big.Int, _data []byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.OnTokenTransfer(&_VRFCoordinator.TransactOpts, _sender, _fee, _data)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address _sender, uint256 _fee, bytes _data) returns()
func (_VRFCoordinator *VRFCoordinatorTransactorSession) OnTokenTransfer(_sender common.Address, _fee *big.Int, _data []byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.OnTokenTransfer(&_VRFCoordinator.TransactOpts, _sender, _fee, _data)
}

// RegisterProvingKey is a paid mutator transaction binding the contract method 0x6815851e.
//
// Solidity: function registerProvingKey(uint256 _fee, uint256[2] _publicProvingKey, bytes32 _jobID) returns(bytes32 keyHash, address oracle, uint256 fee)
func (_VRFCoordinator *VRFCoordinatorTransactor) RegisterProvingKey(opts *bind.TransactOpts, _fee *big.Int, _publicProvingKey [2]*big.Int, _jobID [32]byte) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "registerProvingKey", _fee, _publicProvingKey, _jobID)
}

// RegisterProvingKey is a paid mutator transaction binding the contract method 0x6815851e.
//
// Solidity: function registerProvingKey(uint256 _fee, uint256[2] _publicProvingKey, bytes32 _jobID) returns(bytes32 keyHash, address oracle, uint256 fee)
func (_VRFCoordinator *VRFCoordinatorSession) RegisterProvingKey(_fee *big.Int, _publicProvingKey [2]*big.Int, _jobID [32]byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.RegisterProvingKey(&_VRFCoordinator.TransactOpts, _fee, _publicProvingKey, _jobID)
}

// RegisterProvingKey is a paid mutator transaction binding the contract method 0x6815851e.
//
// Solidity: function registerProvingKey(uint256 _fee, uint256[2] _publicProvingKey, bytes32 _jobID) returns(bytes32 keyHash, address oracle, uint256 fee)
func (_VRFCoordinator *VRFCoordinatorTransactorSession) RegisterProvingKey(_fee *big.Int, _publicProvingKey [2]*big.Int, _jobID [32]byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.RegisterProvingKey(&_VRFCoordinator.TransactOpts, _fee, _publicProvingKey, _jobID)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address _recipient, uint256 _amount) returns()
func (_VRFCoordinator *VRFCoordinatorTransactor) Withdraw(opts *bind.TransactOpts, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "withdraw", _recipient, _amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address _recipient, uint256 _amount) returns()
func (_VRFCoordinator *VRFCoordinatorSession) Withdraw(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.Withdraw(&_VRFCoordinator.TransactOpts, _recipient, _amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address _recipient, uint256 _amount) returns()
func (_VRFCoordinator *VRFCoordinatorTransactorSession) Withdraw(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.Withdraw(&_VRFCoordinator.TransactOpts, _recipient, _amount)
}

// VRFCoordinatorNewServiceAgreementIterator is returned from FilterNewServiceAgreement and is used to iterate over the raw logs and unpacked data for NewServiceAgreement events raised by the VRFCoordinator contract.
type VRFCoordinatorNewServiceAgreementIterator struct {
	Event *VRFCoordinatorNewServiceAgreement // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *VRFCoordinatorNewServiceAgreementIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorNewServiceAgreement)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(VRFCoordinatorNewServiceAgreement)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *VRFCoordinatorNewServiceAgreementIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VRFCoordinatorNewServiceAgreementIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VRFCoordinatorNewServiceAgreement represents a NewServiceAgreement event raised by the VRFCoordinator contract.
type VRFCoordinatorNewServiceAgreement struct {
	KeyHash [32]byte
	Fee     *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterNewServiceAgreement is a free log retrieval operation binding the contract event 0xae189157e0628c1e62315e9179156e1ea10e90e9c15060002f7021e907dc2cfe.
//
// Solidity: event NewServiceAgreement(bytes32 keyHash, uint256 fee)
func (_VRFCoordinator *VRFCoordinatorFilterer) FilterNewServiceAgreement(opts *bind.FilterOpts) (*VRFCoordinatorNewServiceAgreementIterator, error) {

	logs, sub, err := _VRFCoordinator.contract.FilterLogs(opts, "NewServiceAgreement")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorNewServiceAgreementIterator{contract: _VRFCoordinator.contract, event: "NewServiceAgreement", logs: logs, sub: sub}, nil
}

// WatchNewServiceAgreement is a free log subscription operation binding the contract event 0xae189157e0628c1e62315e9179156e1ea10e90e9c15060002f7021e907dc2cfe.
//
// Solidity: event NewServiceAgreement(bytes32 keyHash, uint256 fee)
func (_VRFCoordinator *VRFCoordinatorFilterer) WatchNewServiceAgreement(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorNewServiceAgreement) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinator.contract.WatchLogs(opts, "NewServiceAgreement")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VRFCoordinatorNewServiceAgreement)
				if err := _VRFCoordinator.contract.UnpackLog(event, "NewServiceAgreement", log); err != nil {
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

// ParseNewServiceAgreement is a log parse operation binding the contract event 0xae189157e0628c1e62315e9179156e1ea10e90e9c15060002f7021e907dc2cfe.
//
// Solidity: event NewServiceAgreement(bytes32 keyHash, uint256 fee)
func (_VRFCoordinator *VRFCoordinatorFilterer) ParseNewServiceAgreement(log types.Log) (*VRFCoordinatorNewServiceAgreement, error) {
	event := new(VRFCoordinatorNewServiceAgreement)
	if err := _VRFCoordinator.contract.UnpackLog(event, "NewServiceAgreement", log); err != nil {
		return nil, err
	}
	return event, nil
}

// VRFCoordinatorRandomnessRequestIterator is returned from FilterRandomnessRequest and is used to iterate over the raw logs and unpacked data for RandomnessRequest events raised by the VRFCoordinator contract.
type VRFCoordinatorRandomnessRequestIterator struct {
	Event *VRFCoordinatorRandomnessRequest // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *VRFCoordinatorRandomnessRequestIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorRandomnessRequest)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(VRFCoordinatorRandomnessRequest)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *VRFCoordinatorRandomnessRequestIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VRFCoordinatorRandomnessRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VRFCoordinatorRandomnessRequest represents a RandomnessRequest event raised by the VRFCoordinator contract.
type VRFCoordinatorRandomnessRequest struct {
	KeyHash [32]byte
	Seed    *big.Int
	JobID   [32]byte
	Sender  common.Address
	Fee     *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRandomnessRequest is a free log retrieval operation binding the contract event 0xd241d78a52145a5d1d1ff002e32ec15cdc395631bcee66246650c2429dfaccaa.
//
// Solidity: event RandomnessRequest(bytes32 keyHash, uint256 seed, bytes32 indexed jobID, address sender, uint256 fee)
func (_VRFCoordinator *VRFCoordinatorFilterer) FilterRandomnessRequest(opts *bind.FilterOpts, jobID [][32]byte) (*VRFCoordinatorRandomnessRequestIterator, error) {

	var jobIDRule []interface{}
	for _, jobIDItem := range jobID {
		jobIDRule = append(jobIDRule, jobIDItem)
	}

	logs, sub, err := _VRFCoordinator.contract.FilterLogs(opts, "RandomnessRequest", jobIDRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorRandomnessRequestIterator{contract: _VRFCoordinator.contract, event: "RandomnessRequest", logs: logs, sub: sub}, nil
}

// WatchRandomnessRequest is a free log subscription operation binding the contract event 0xd241d78a52145a5d1d1ff002e32ec15cdc395631bcee66246650c2429dfaccaa.
//
// Solidity: event RandomnessRequest(bytes32 keyHash, uint256 seed, bytes32 indexed jobID, address sender, uint256 fee)
func (_VRFCoordinator *VRFCoordinatorFilterer) WatchRandomnessRequest(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorRandomnessRequest, jobID [][32]byte) (event.Subscription, error) {

	var jobIDRule []interface{}
	for _, jobIDItem := range jobID {
		jobIDRule = append(jobIDRule, jobIDItem)
	}

	logs, sub, err := _VRFCoordinator.contract.WatchLogs(opts, "RandomnessRequest", jobIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VRFCoordinatorRandomnessRequest)
				if err := _VRFCoordinator.contract.UnpackLog(event, "RandomnessRequest", log); err != nil {
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

// ParseRandomnessRequest is a log parse operation binding the contract event 0xd241d78a52145a5d1d1ff002e32ec15cdc395631bcee66246650c2429dfaccaa.
//
// Solidity: event RandomnessRequest(bytes32 keyHash, uint256 seed, bytes32 indexed jobID, address sender, uint256 fee)
func (_VRFCoordinator *VRFCoordinatorFilterer) ParseRandomnessRequest(log types.Log) (*VRFCoordinatorRandomnessRequest, error) {
	event := new(VRFCoordinatorRandomnessRequest)
	if err := _VRFCoordinator.contract.UnpackLog(event, "RandomnessRequest", log); err != nil {
		return nil, err
	}
	return event, nil
}
