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
const VRFCoordinatorABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"withdrawableTokens\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"callbacks\",\"outputs\":[{\"name\":\"callbackContract\",\"type\":\"address\"},{\"name\":\"randomnessFee\",\"type\":\"uint256\"},{\"name\":\"seed\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_proof\",\"type\":\"bytes\"}],\"name\":\"fulfillRandomnessRequest\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_fee\",\"type\":\"uint256\"},{\"name\":\"_publicProvingKey\",\"type\":\"uint256[2]\"},{\"name\":\"_jobID\",\"type\":\"bytes32\"}],\"name\":\"registerProvingKey\",\"outputs\":[{\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"name\":\"oracle\",\"type\":\"address\"},{\"name\":\"fee\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"serviceAgreements\",\"outputs\":[{\"name\":\"vRFOracle\",\"type\":\"address\"},{\"name\":\"jobID\",\"type\":\"bytes32\"},{\"name\":\"fee\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_keyHash\",\"type\":\"bytes32\"},{\"name\":\"_seed\",\"type\":\"uint256\"}],\"name\":\"makeRequestId\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_sender\",\"type\":\"address\"},{\"name\":\"_fee\",\"type\":\"uint256\"},{\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_publicKey\",\"type\":\"uint256[2]\"}],\"name\":\"hashOfKey\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_recipient\",\"type\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"proof\",\"type\":\"bytes\"}],\"name\":\"randomValueFromVRFProof\",\"outputs\":[{\"name\":\"output\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_link\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"seed\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"jobID\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"RandomnessRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"NewServiceAgreement\",\"type\":\"event\"}]"

// VRFCoordinatorBin is the compiled bytecode used for deploying new contracts.
var VRFCoordinatorBin = "0x608060405234801561001057600080fd5b50604051602080612bf88339810180604052602081101561003057600080fd5b8101908080519060200190929190505050806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050612b67806100916000396000f3fe6080604052600436106100a3576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680626f6ad0146100a857806321f365091461010d5780635e1c1059146101965780636815851e1461027657806375d3507014610313578063916bf6c71461039c578063a4c0ed36146103f5578063caf70c4a146104e7578063f3fef3a31461056d578063fa8fc6f1146105c8575b600080fd5b3480156100b457600080fd5b506100f7600480360360208110156100cb57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506106a4565b6040518082815260200191505060405180910390f35b34801561011957600080fd5b506101466004803603602081101561013057600080fd5b81019080803590602001909291905050506106bc565b604051808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001838152602001828152602001935050505060405180910390f35b3480156101a257600080fd5b5061025c600480360360208110156101b957600080fd5b81019080803590602001906401000000008111156101d657600080fd5b8201836020820111156101e857600080fd5b8035906020019184600183028401116401000000008311171561020a57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290505050610706565b604051808215151515815260200191505060405180910390f35b34801561028257600080fd5b506102c36004803603608081101561029957600080fd5b81019080803590602001909291908060400190919291929080359060200190929190505050610acf565b604051808481526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001828152602001935050505060405180910390f35b34801561031f57600080fd5b5061034c6004803603602081101561033657600080fd5b8101908080359060200190929190505050610d08565b604051808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001838152602001828152602001935050505060405180910390f35b3480156103a857600080fd5b506103df600480360360408110156103bf57600080fd5b810190808035906020019092919080359060200190929190505050610d52565b6040518082815260200191505060405180910390f35b34801561040157600080fd5b506104e56004803603606081101561041857600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291908035906020019064010000000081111561045f57600080fd5b82018360208201111561047157600080fd5b8035906020019184600183028401116401000000008311171561049357600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290505050610d8b565b005b3480156104f357600080fd5b506105576004803603604081101561050a57600080fd5b8101908080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f8201169050808301925050505050509192919290505050610e99565b6040518082815260200191505060405180910390f35b34801561057957600080fd5b506105c66004803603604081101561059057600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050610ef1565b005b3480156105d457600080fd5b5061068e600480360360208110156105eb57600080fd5b810190808035906020019064010000000081111561060857600080fd5b82018360208201111561061a57600080fd5b8035906020019184600183028401116401000000008311171561063c57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290505050611102565b6040518082815260200191505060405180910390f35b60036020528060005260406000206000915090505481565b60016020528060005260406000206000915090508060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060010154908060020154905083565b6000610710612a78565b600060208401915060e08401519050600061072a83610e99565b905060006107388284610d52565b9050610742612a9a565b60016000838152602001908152602001600020606060405190810160405290816000820160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001600182015481526020016002820154815250509050600073ffffffffffffffffffffffffffffffffffffffff16816000015173ffffffffffffffffffffffffffffffffffffffff1614151515610876576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260188152602001807f6e6f20636f72726573706f6e64696e672072657175657374000000000000000081525060200191505060405180910390fd5b600061088188611102565b9050600160046000868152602001908152602001600020600087815260200190815260200160002060006101000a81548160ff0219169083151502179055508160200151600360006002600088815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282540192505081905550600060608173ffffffffffffffffffffffffffffffffffffffff16631f1f897f90507c01000000000000000000000000000000000000000000000000000000000285846040516024018083815260200182815260200192505050604051602081830303815290604052907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff838183161783525050505090506000846000015173ffffffffffffffffffffffffffffffffffffffff16826040518082805190602001908083835b602083101515610a535780518252602082019150602081019050602083039250610a2e565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d8060008114610ab5576040519150601f19603f3d011682016040523d82523d6000602084013e610aba565b606091505b50509050809950505050505050505050919050565b6000806000610b14856002806020026040519081016040528092919082600260200280828437600081840152601f19601f820116905080830192505050505050610e99565b925060006002600085815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16141515610bf5576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260198152602001807f706c656173652072656769737465722061206e6577206b65790000000000000081525060200191505060405180910390fd5b336002600086815260200190815260200160002060000160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508460026000868152602001908152602001600020600101819055508660026000868152602001908152602001600020600201819055507fae189157e0628c1e62315e9179156e1ea10e90e9c15060002f7021e907dc2cfe8488604051808381526020018281526020019250505060405180910390a1836002600086815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16889350935093505093509350939050565b60026020528060005260406000206000915090508060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060010154908060020154905083565b60008282604051602001808381526020018281526020019250505060405160208183030381529060405280519060200120905092915050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141515610e4f576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260138152602001807f4d75737420757365204c494e4b20746f6b656e0000000000000000000000000081525060200191505060405180910390fd5b600080828060200190516040811015610e6757600080fd5b81019080805190602001909291908051906020019092919050505091509150610e92828286886112d2565b5050505050565b6000816040516020018082600260200280838360005b83811015610eca578082015181840152602081019050610eaf565b50505050905001915050604051602081830303815290604052805190602001209050919050565b8080600360003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205410151515610fa9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260208152602001807f63616e2774207769746864726177206d6f7265207468616e2062616c616e636581525060200191505060405180910390fd5b81600360003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825403925050819055506000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb84846040518363ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b1580156110ba57600080fd5b505af11580156110ce573d6000803e3d6000fd5b505050506040513d60208110156110e457600080fd5b810190808051906020019092919050505015156110fd57fe5b505050565b60006101a0825114151561117e576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260128152602001807f77726f6e672070726f6f66206c656e677468000000000000000000000000000081525060200191505060405180910390fd5b611186612a78565b61118e612a78565b611196612ad2565b60006111a0612a78565b6111a8612a78565b6000888060200190516101a08110156111c057600080fd5b810190809190826040019190826040019190826060018051906020019092919091908260400191908260400180519060200190929190505050869650859550849450839350829250819150809750819850829950839a50849b50859c50869d5050505050505050611271878787600060038110151561123b57fe5b602002015188600160038110151561124f57fe5b602002015189600260038110151561126357fe5b6020020151898989896115b6565b856040516020018082600260200280838360005b838110156112a0578082015181840152602081019050611285565b505050509050019150506040516020818303038152906040528051906020012060019004975050505050505050919050565b818460026000828152602001908152602001600020600201548210151515611362576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f42656c6f7720616772656564207061796d656e7400000000000000000000000081525060200191505060405180910390fd5b858560046000838152602001908152602001600020600082815260200190815260200160002060009054906101000a900460ff1615151561140b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f706c6561736520726571756573742061206e6f76656c2073656564000000000081525060200191505060405180910390fd5b60006114178989610d52565b9050600073ffffffffffffffffffffffffffffffffffffffff166001600083815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614151561148757fe5b856001600083815260200190815260200160002060000160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508660016000838152602001908152602001600020600101819055508760016000838152602001908152602001600020600201819055507fd241d78a52145a5d1d1ff002e32ec15cdc395631bcee66246650c2429dfaccaa8989600260008d815260200190815260200160002060010154898b604051808681526020018581526020018481526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019550505050505060405180910390a1505050505050505050565b6115bf896118e7565b1515611633576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601a8152602001807f7075626c6963206b6579206973206e6f74206f6e20637572766500000000000081525060200191505060405180910390fd5b61163c886118e7565b15156116b0576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260158152602001807f67616d6d61206973206e6f74206f6e206375727665000000000000000000000081525060200191505060405180910390fd5b6116b9836118e7565b151561172d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601d8152602001807f6347616d6d615769746e657373206973206e6f74206f6e20637572766500000081525060200191505060405180910390fd5b611736826118e7565b15156117aa576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601c8152602001807f73486173685769746e657373206973206e6f74206f6e2063757276650000000081525060200191505060405180910390fd5b6117b6878a8887611960565b151561182a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601a8152602001807f6164647228632a706b2b732a6729e289a05f755769746e65737300000000000081525060200191505060405180910390fd5b611832612a78565b61183c8a87611ba8565b9050611846612a78565b611855898b878b868989611dfa565b9050611864828c8c8985612030565b891415156118da576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600d8152602001807f696e76616c69642070726f6f660000000000000000000000000000000000000081525060200191505060405180910390fd5b5050505050505050505050565b60007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80151561191357fe5b82600160028110151561192257fe5b602002015183600160028110151561193657fe5b60200201510961195883600060028110151561194e57fe5b6020020151612163565b149050919050565b60008073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614151515611a06576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600b8152602001807f626164207769746e65737300000000000000000000000000000000000000000081525060200191505060405180910390fd5b60007ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141801515611a3257fe5b84866000600281101515611a4257fe5b6020020151097ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd03641410360010290506000806002876001600281101515611a8357fe5b6020020151811515611a9157fe5b0614611a9e57601c611aa1565b601b5b905060007ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141801515611acf57fe5b876000600281101515611ade57fe5b6020020151890960010290506000600184848a6000600281101515611aff57fe5b60200201516001028560405160008152602001604052604051808581526020018460ff1660ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa158015611b5f573d6000803e3d6000fd5b5050506020604051035190508573ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614945050505050949350505050565b611bb0612a78565b611c1383836040516020018083600260200280838360005b83811015611be3578082015181840152602081019050611bc8565b505050509050018281526020019250505060405160208183030381529060405280519060200120600190046121f7565b816000600281101515611c2257fe5b602002018181525050611c4f611c4a826000600281101515611c4057fe5b6020020151612163565b61225d565b816001600281101515611c5e57fe5b6020020181815250505b611c84816000600281101515611c7a57fe5b6020020151612163565b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f801515611cae57fe5b826001600281101515611cbd57fe5b6020020151836001600281101515611cd157fe5b602002015109141515611d7d57611d24816000600281101515611cf057fe5b60200201516040516020018082815260200191505060405160208183030381529060405280519060200120600190046121f7565b816000600281101515611d3357fe5b602002018181525050611d60611d5b826000600281101515611d5157fe5b6020020151612163565b61225d565b816001600281101515611d6f57fe5b602002018181525050611c68565b60016002826001600281101515611d9057fe5b6020020151811515611d9e57fe5b061415611df457806001600281101515611db457fe5b60200201517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f03816001600281101515611dea57fe5b6020020181815250505b92915050565b611e02612a78565b60007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f846000600281101515611e3457fe5b6020020151886000600281101515611e4857fe5b602002015103811515611e5757fe5b0614151515611ece576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f706f696e747320696e2073756d206d7573742062652064697374696e6374000081525060200191505060405180910390fd5b611ed987898861229f565b1515611f73576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260218152602001807f4669727374206d756c7469706c69636174696f6e20636865636b206661696c6581526020017f640000000000000000000000000000000000000000000000000000000000000081525060400191505060405180910390fd5b611f7e84868561229f565b1515612018576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260228152602001807f5365636f6e64206d756c7469706c69636174696f6e20636865636b206661696c81526020017f656400000000000000000000000000000000000000000000000000000000000081525060400191505060405180910390fd5b612023868484612438565b9050979650505050505050565b600085858584866040516020018086600260200280838360005b8381101561206557808201518184015260208101905061204a565b5050505090500185600260200280838360005b83811015612093578082015181840152602081019050612078565b5050505090500184600260200280838360005b838110156120c15780820151818401526020810190506120a6565b5050505090500183600260200280838360005b838110156120ef5780820151818401526020810190506120d4565b505050509050018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166c01000000000000000000000000028152601401955050505050506040516020818303038152906040528051906020012060019004905095945050505050565b6000807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80151561219057fe5b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8015156121ba57fe5b848509840990507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8015156121eb57fe5b60078208915050919050565b60008190505b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f81101515612258578060405160200180828152602001915050604051602081830303815290604052805190602001206001900490506121fd565b919050565b600061229882600260017ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f01908060020a82049150506125c6565b9050919050565b60008083141515156122b057600080fd5b60007ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd03641418015156122dc57fe5b8560006002811015156122eb57fe5b602002015185096001029050600080600287600160028110151561230b57fe5b602002015181151561231957fe5b06141561232757601b61232a565b601c5b9050836040516020018082600260200280838360005b8381101561235b578082015181840152602081019050612340565b50505050905001915050604051602081830303815290604052805190602001206001900473ffffffffffffffffffffffffffffffffffffffff1660016000600102838960006002811015156123ac57fe5b60200201516001028660405160008152602001604052604051808581526020018460ff1660ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa15801561240c573d6000803e3d6000fd5b5050506020604051035173ffffffffffffffffffffffffffffffffffffffff1614925050509392505050565b612440612a78565b600080600061249d87600060028110151561245757fe5b602002015188600160028110151561246b57fe5b602002015188600060028110151561247f57fe5b602002015189600160028110151561249357fe5b6020020151612736565b80935081945082955050505060017ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8015156124d557fe5b86830914151561254d576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260198152602001807f696e765a206d75737420626520696e7665727365206f66207a0000000000000081525060200191505060405180910390fd5b60408051908101604052807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80151561258257fe5b87860981526020017ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8015156125b457fe5b87850981525093505050509392505050565b6000806125d1612af5565b60208160006006811015156125e257fe5b60200201818152505060208160016006811015156125fc57fe5b602002018181525050602081600260068110151561261657fe5b6020020181815250508481600360068110151561262f57fe5b6020020181815250508381600460068110151561264857fe5b6020020181815250507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f81600560068110151561268157fe5b602002018181525050612692612b18565b60208160c0846005600019fa92506000831415612717576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260128152602001807f6269674d6f64457870206661696c75726521000000000000000000000000000081525060200191505060405180910390fd5b80600060018110151561272657fe5b6020020151935050505092915050565b60008060008060006001809150915060007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80151561277157fe5b897ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f038808905060007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8015156127c457fe5b8b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f038a08905060006127f983838585612916565b809250819950505061280d88828e88612984565b809250819950505061282188828c87612984565b809250819950505060006128378d878b85612984565b809250819950505061284b88828686612916565b809250819950505061285f88828e89612984565b80925081995050508082141515612902577ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80151561289a57fe5b818a0998507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8015156128c957fe5b82890997507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8015156128f857fe5b8183099650612906565b8196505b5050505050509450945094915050565b6000807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80151561294357fe5b8487097ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80151561297057fe5b848709809250819350505094509492505050565b60008060007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8015156129b357fe5b878509905060007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8015156129e457fe5b87877ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f030990507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f801515612a3557fe5b8183087ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f801515612a6257fe5b8689098094508195505050505094509492505050565b6040805190810160405280600290602082028038833980820191505090505090565b606060405190810160405280600073ffffffffffffffffffffffffffffffffffffffff16815260200160008152602001600081525090565b606060405190810160405280600390602082028038833980820191505090505090565b60c060405190810160405280600690602082028038833980820191505090505090565b60206040519081016040528060019060208202803883398082019150509050509056fea165627a7a7230582058985290d4de3a207a671bcbba353cc7dbb2cf58d19a895acdf9b29b17d8d26d0029"

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

// MakeRequestId is a free data retrieval call binding the contract method 0x916bf6c7.
//
// Solidity: function makeRequestId(bytes32 _keyHash, uint256 _seed) constant returns(bytes32)
func (_VRFCoordinator *VRFCoordinatorCaller) MakeRequestId(opts *bind.CallOpts, _keyHash [32]byte, _seed *big.Int) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _VRFCoordinator.contract.Call(opts, out, "makeRequestId", _keyHash, _seed)
	return *ret0, err
}

// MakeRequestId is a free data retrieval call binding the contract method 0x916bf6c7.
//
// Solidity: function makeRequestId(bytes32 _keyHash, uint256 _seed) constant returns(bytes32)
func (_VRFCoordinator *VRFCoordinatorSession) MakeRequestId(_keyHash [32]byte, _seed *big.Int) ([32]byte, error) {
	return _VRFCoordinator.Contract.MakeRequestId(&_VRFCoordinator.CallOpts, _keyHash, _seed)
}

// MakeRequestId is a free data retrieval call binding the contract method 0x916bf6c7.
//
// Solidity: function makeRequestId(bytes32 _keyHash, uint256 _seed) constant returns(bytes32)
func (_VRFCoordinator *VRFCoordinatorCallerSession) MakeRequestId(_keyHash [32]byte, _seed *big.Int) ([32]byte, error) {
	return _VRFCoordinator.Contract.MakeRequestId(&_VRFCoordinator.CallOpts, _keyHash, _seed)
}

// RandomValueFromVRFProof is a free data retrieval call binding the contract method 0xfa8fc6f1.
//
// Solidity: function randomValueFromVRFProof(bytes proof) constant returns(uint256 output)
func (_VRFCoordinator *VRFCoordinatorCaller) RandomValueFromVRFProof(opts *bind.CallOpts, proof []byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _VRFCoordinator.contract.Call(opts, out, "randomValueFromVRFProof", proof)
	return *ret0, err
}

// RandomValueFromVRFProof is a free data retrieval call binding the contract method 0xfa8fc6f1.
//
// Solidity: function randomValueFromVRFProof(bytes proof) constant returns(uint256 output)
func (_VRFCoordinator *VRFCoordinatorSession) RandomValueFromVRFProof(proof []byte) (*big.Int, error) {
	return _VRFCoordinator.Contract.RandomValueFromVRFProof(&_VRFCoordinator.CallOpts, proof)
}

// RandomValueFromVRFProof is a free data retrieval call binding the contract method 0xfa8fc6f1.
//
// Solidity: function randomValueFromVRFProof(bytes proof) constant returns(uint256 output)
func (_VRFCoordinator *VRFCoordinatorCallerSession) RandomValueFromVRFProof(proof []byte) (*big.Int, error) {
	return _VRFCoordinator.Contract.RandomValueFromVRFProof(&_VRFCoordinator.CallOpts, proof)
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
// Solidity: event RandomnessRequest(bytes32 keyHash, uint256 seed, bytes32 jobID, address sender, uint256 fee)
func (_VRFCoordinator *VRFCoordinatorFilterer) FilterRandomnessRequest(opts *bind.FilterOpts) (*VRFCoordinatorRandomnessRequestIterator, error) {

	logs, sub, err := _VRFCoordinator.contract.FilterLogs(opts, "RandomnessRequest")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorRandomnessRequestIterator{contract: _VRFCoordinator.contract, event: "RandomnessRequest", logs: logs, sub: sub}, nil
}

// WatchRandomnessRequest is a free log subscription operation binding the contract event 0xd241d78a52145a5d1d1ff002e32ec15cdc395631bcee66246650c2429dfaccaa.
//
// Solidity: event RandomnessRequest(bytes32 keyHash, uint256 seed, bytes32 jobID, address sender, uint256 fee)
func (_VRFCoordinator *VRFCoordinatorFilterer) WatchRandomnessRequest(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorRandomnessRequest) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinator.contract.WatchLogs(opts, "RandomnessRequest")
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
// Solidity: event RandomnessRequest(bytes32 keyHash, uint256 seed, bytes32 jobID, address sender, uint256 fee)
func (_VRFCoordinator *VRFCoordinatorFilterer) ParseRandomnessRequest(log types.Log) (*VRFCoordinatorRandomnessRequest, error) {
	event := new(VRFCoordinatorRandomnessRequest)
	if err := _VRFCoordinator.contract.UnpackLog(event, "RandomnessRequest", log); err != nil {
		return nil, err
	}
	return event, nil
}
