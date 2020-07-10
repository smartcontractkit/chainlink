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
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// VRFCoordinatorABI is the input ABI used to generate the binding from.
const VRFCoordinatorABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"NewServiceAgreement\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"jobID\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"RandomnessRequest\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"callbacks\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"callbackContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"randomnessFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_proof\",\"type\":\"bytes\"}],\"name\":\"fulfillRandomnessRequest\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"_publicKey\",\"type\":\"uint256[2]\"}],\"name\":\"hashOfKey\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"_publicProvingKey\",\"type\":\"uint256[2]\"},{\"internalType\":\"bytes32\",\"name\":\"_jobID\",\"type\":\"bytes32\"}],\"name\":\"registerProvingKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"serviceAgreements\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"vRFOracle\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"jobID\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"withdrawableTokens\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// VRFCoordinatorBin is the compiled bytecode used for deploying new contracts.
var VRFCoordinatorBin = "0x608060405234801561001057600080fd5b506040516121cc3803806121cc8339818101604052602081101561003357600080fd5b5051600080546001600160a01b039092166001600160a01b0319909216919091179055612167806100656000396000f3fe608060405234801561001057600080fd5b50600436106100875760003560e01c8063a4c0ed361161005b578063a4c0ed36146101fa578063caf70c4a146102c4578063d83402091461030f578063f3fef3a31461035357610087565b80626f6ad01461008c57806321f36509146100d15780635e1c10591461012357806375d35070146101dd575b600080fd5b6100bf600480360360208110156100a257600080fd5b503573ffffffffffffffffffffffffffffffffffffffff1661038c565b60408051918252519081900360200190f35b6100ee600480360360208110156100e757600080fd5b503561039e565b6040805173ffffffffffffffffffffffffffffffffffffffff9094168452602084019290925282820152519081900360600190f35b6101c96004803603602081101561013957600080fd5b81019060208101813564010000000081111561015457600080fd5b82018360208201111561016657600080fd5b8035906020019184600183028401116401000000008311171561018857600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295506103d7945050505050565b604080519115158252519081900360200190f35b6100ee600480360360208110156101f357600080fd5b50356106c5565b6102c26004803603606081101561021057600080fd5b73ffffffffffffffffffffffffffffffffffffffff8235169160208101359181019060608101604082013564010000000081111561024d57600080fd5b82018360208201111561025f57600080fd5b8035906020019184600183028401116401000000008311171561028157600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295506106fe945050505050565b005b6100bf600480360360408110156102da57600080fd5b604080518082018252918301929181830191839060029083908390808284376000920191909152509194506107bd9350505050565b6102c2600480360360a081101561032557600080fd5b5080359073ffffffffffffffffffffffffffffffffffffffff60208201351690604081019060800135610812565b6102c26004803603604081101561036957600080fd5b5073ffffffffffffffffffffffffffffffffffffffff81351690602001356109f5565b60036020526000908152604090205481565b600160208190526000918252604090912080549181015460029091015473ffffffffffffffffffffffffffffffffffffffff9092169183565b60006103e161203f565b5060e0820151602083019060006103f7836107bd565b905060006104058284610b57565b905061040f61205d565b506000818152600160208181526040928390208351606081018552815473ffffffffffffffffffffffffffffffffffffffff16808252938201549281019290925260020154928101929092526104c657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f6e6f20636f72726573706f6e64696e6720726571756573740000000000000000604482015290519081900360640190fd5b60006104d188610b86565b6000858152600260209081526040808320548683015173ffffffffffffffffffffffffffffffffffffffff90911680855260039093529220549293509161051d9163ffffffff610cf416565b73ffffffffffffffffffffffffffffffffffffffff8083166000908152600360209081526040808320949094558351602481018990526044808201889052855180830390910181526064909101855290810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f94985ddd00000000000000000000000000000000000000000000000000000000178152875194518251939592948694931692859282918083835b6020831061060757805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe090920191602091820191016105ca565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d8060008114610669576040519150601f19603f3d011682016040523d82523d6000602084013e61066e565b606091505b50506000978852600160208190526040892080547fffffffffffffffffffffffff0000000000000000000000000000000000000000168155908101899055600201979097555094985050505050505050505b919050565b600260208190526000918252604090912080546001820154919092015473ffffffffffffffffffffffffffffffffffffffff9092169183565b60005473ffffffffffffffffffffffffffffffffffffffff16331461078457604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f4d75737420757365204c494e4b20746f6b656e00000000000000000000000000604482015290519081900360640190fd5b60008082806020019051604081101561079c57600080fd5b50805160209091015190925090506107b682828688610d6f565b5050505050565b6000816040516020018082600260200280838360005b838110156107eb5781810151838201526020016107d3565b50505050905001915050604051602081830303815290604052805190602001209050919050565b6040805180820182526000916108419190859060029083908390808284376000920191909152506107bd915050565b60008181526002602052604090205490915073ffffffffffffffffffffffffffffffffffffffff1680156108d657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f706c656173652072656769737465722061206e6577206b657900000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff851661095857604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f5f6f7261636c65206d757374206e6f7420626520307830000000000000000000604482015290519081900360640190fd5b60008281526002602081815260409283902080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8a1617815560018101879055909101889055815184815290810188905281517fae189157e0628c1e62315e9179156e1ea10e90e9c15060002f7021e907dc2cfe929181900390910190a1505050505050565b336000908152600360205260409020548190811115610a7557604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f63616e2774207769746864726177206d6f7265207468616e2062616c616e6365604482015290519081900360640190fd5b33600090815260036020526040902054610a95908363ffffffff610fa316565b33600090815260036020908152604080832093909355815483517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8881166004830152602482018890529451949091169363a9059cbb93604480840194938390030190829087803b158015610b2057600080fd5b505af1158015610b34573d6000803e3d6000fd5b505050506040513d6020811015610b4a57600080fd5b5051610b5257fe5b505050565b604080516020808201859052818301849052825180830384018152606090920190925280519101205b92915050565b60006101a0825114610bf957604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f77726f6e672070726f6f66206c656e6774680000000000000000000000000000604482015290519081900360640190fd5b610c0161203f565b610c0961203f565b610c11612094565b6000610c1b61203f565b610c2361203f565b6000888060200190516101a0811015610c3b57600080fd5b5060e0810151610180820151919850604089019750608089019650945061010088019350610140880192509050610c8e87878760006020020151886001602002015189600260200201518989898961101a565b6003866040516020018083815260200182600260200280838360005b83811015610cc2578181015183820152602001610caa565b50505050905001925050506040516020818303038152906040528051906020012060001c975050505050505050919050565b600082820183811015610d6857604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b9392505050565b81846002600082815260200190815260200160002060020154821015610df657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f42656c6f7720616772656564207061796d656e74000000000000000000000000604482015290519081900360640190fd5b600086815260046020908152604080832073ffffffffffffffffffffffffffffffffffffffff8716845290915281205490610e338888878561131b565b90506000610e418983610b57565b60008181526001602052604090205490915073ffffffffffffffffffffffffffffffffffffffff1615610e7057fe5b6000818152600160208181526040808420805473ffffffffffffffffffffffffffffffffffffffff8c167fffffffffffffffffffffffff0000000000000000000000000000000000000000909116811782558185018d905560029182018890558e8652908352938190209092015482518d815291820186905281830193909352606081018a905290517fd241d78a52145a5d1d1ff002e32ec15cdc395631bcee66246650c2429dfaccaa9181900360800190a2600089815260046020908152604080832073ffffffffffffffffffffffffffffffffffffffff8a168452909152902054610f6490600163ffffffff610cf416565b6000998a52600460209081526040808c2073ffffffffffffffffffffffffffffffffffffffff9099168c52979052959098209490945550505050505050565b60008282111561101457604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f536166654d6174683a207375627472616374696f6e206f766572666c6f770000604482015290519081900360640190fd5b50900390565b6110238961136f565b61108e57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f7075626c6963206b6579206973206e6f74206f6e206375727665000000000000604482015290519081900360640190fd5b6110978861136f565b61110257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f67616d6d61206973206e6f74206f6e2063757276650000000000000000000000604482015290519081900360640190fd5b61110b8361136f565b61117657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f6347616d6d615769746e657373206973206e6f74206f6e206375727665000000604482015290519081900360640190fd5b61117f8261136f565b6111ea57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f73486173685769746e657373206973206e6f74206f6e20637572766500000000604482015290519081900360640190fd5b6111f6878a88876113b3565b61126157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f6164647228632a706b2b732a6729e289a05f755769746e657373000000000000604482015290519081900360640190fd5b61126961203f565b6112738a8761157f565b905061127d61203f565b61128c898b878b868989611622565b9050600061129d838d8d8a86611795565b9050808a1461130d57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f696e76616c69642070726f6f6600000000000000000000000000000000000000604482015290519081900360640190fd5b505050505050505050505050565b604080516020808201969096528082019490945273ffffffffffffffffffffffffffffffffffffffff9290921660608401526080808401919091528151808403909101815260a09092019052805191012090565b60208101516000907ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f9080096113ac8360005b60200201516118b8565b1492915050565b600073ffffffffffffffffffffffffffffffffffffffff821661143757604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f626164207769746e657373000000000000000000000000000000000000000000604482015290519081900360640190fd5b60208401516000906001161561144e57601c611451565b601b5b905060007ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd03641418587600060200201510986517ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd036414191820392506000919089098751604080516000808252602082810180855288905260ff8916838501526060830194909452608082018590529151939450909260019260a08084019391927fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081019281900390910190855afa15801561152c573d6000803e3d6000fd5b50506040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0015173ffffffffffffffffffffffffffffffffffffffff9081169088161495505050505050949350505050565b61158761203f565b6115e5600184846040516020018084815260200183600260200280838360005b838110156115bf5781810151838201526020016115a7565b505050509050018281526020019350505050604051602081830303815290604052611910565b90505b6115f18161136f565b610b8057805160408051602081810193909352815180820390930183528101905261161b90611910565b90506115e8565b61162a61203f565b825186517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f919003066116be57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f706f696e747320696e2073756d206d7573742062652064697374696e63740000604482015290519081900360640190fd5b6116c9878988611978565b61171e576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260218152602001806120ef6021913960400191505060405180910390fd5b611729848685611978565b61177e576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260228152602001806121106022913960400191505060405180910390fd5b611789868484611ade565b98975050505050505050565b6000600286868685876040516020018087815260200186600260200280838360005b838110156117cf5781810151838201526020016117b7565b5050505090500185600260200280838360005b838110156117fa5781810151838201526020016117e2565b5050505090500184600260200280838360005b8381101561182557818101518382015260200161180d565b5050505090500183600260200280838360005b83811015611850578181015183820152602001611838565b505050509050018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1660601b815260140196505050505050506040516020818303038152906040528051906020012060001c905095945050505050565b6000807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80848509840990507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f600782089392505050565b61191861203f565b61192182611c0c565b81526119366119318260006113a2565b611c61565b6020820181905260029006600114156106c0576020810180517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f039052919050565b60008261198457600080fd5b835160208501516000906001161561199d57601c6119a0565b601b5b905060007ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141838709604080516000808252602080830180855282905260ff871683850152606083018890526080830185905292519394509260019260a08084019391927fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081019281900390910190855afa158015611a42573d6000803e3d6000fd5b5050506020604051035190506000866040516020018082600260200280838360005b83811015611a7c578181015183820152602001611a64565b505050509050019150506040516020818303038152906040528051906020012060001c90508073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614955050505050509392505050565b611ae661203f565b835160208086015185519186015160009384938493611b0793909190611c8d565b919450925090507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f858209600114611ba057604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f696e765a206d75737420626520696e7665727365206f66207a00000000000000604482015290519081900360640190fd5b60405180604001604052807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80611bd357fe5b87860981526020017ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8785099052979650505050505050565b805160208201205b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f81106106c057604080516020808201939093528151808203840181529082019091528051910120611c14565b6000610b80827f3fffffffffffffffffffffffffffffffffffffffffffffffffffffffbfffff0c611e23565b60008080600180827ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f897ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f038808905060007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f038a0890506000611d3583838585611f1c565b9098509050611d4688828e88611f74565b9098509050611d5788828c87611f74565b90985090506000611d6a8d878b85611f74565b9098509050611d7b88828686611f1c565b9098509050611d8c88828e89611f74565b9098509050818114611e0f577ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f818a0998507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f82890997507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8183099650611e13565b8196505b5050505050509450945094915050565b600080611e2e6120b2565b6020808252818101819052604082015260608101859052608081018490527ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f60a0820152611e7a6120d0565b60208160c08460057ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa925082611f1257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f6269674d6f64457870206661696c757265210000000000000000000000000000604482015290519081900360640190fd5b5195945050505050565b6000807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8487097ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8487099097909650945050505050565b600080807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f878509905060007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f87877ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f030990507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8183087ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f86890990999098509650505050505050565b60405180604001604052806002906020820280368337509192915050565b6040518060600160405280600073ffffffffffffffffffffffffffffffffffffffff16815260200160008152602001600081525090565b60405180606001604052806003906020820280368337509192915050565b6040518060c001604052806006906020820280368337509192915050565b6040518060200160405280600190602082028036833750919291505056fe4669727374206d756c7469706c69636174696f6e20636865636b206661696c65645365636f6e64206d756c7469706c69636174696f6e20636865636b206661696c6564a264697066735822beefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeef64736f6c6343decafe0033"

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
// Solidity: function callbacks(bytes32 ) view returns(address callbackContract, uint256 randomnessFee, uint256 seed)
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
// Solidity: function callbacks(bytes32 ) view returns(address callbackContract, uint256 randomnessFee, uint256 seed)
func (_VRFCoordinator *VRFCoordinatorSession) Callbacks(arg0 [32]byte) (struct {
	CallbackContract common.Address
	RandomnessFee    *big.Int
	Seed             *big.Int
}, error) {
	return _VRFCoordinator.Contract.Callbacks(&_VRFCoordinator.CallOpts, arg0)
}

// Callbacks is a free data retrieval call binding the contract method 0x21f36509.
//
// Solidity: function callbacks(bytes32 ) view returns(address callbackContract, uint256 randomnessFee, uint256 seed)
func (_VRFCoordinator *VRFCoordinatorCallerSession) Callbacks(arg0 [32]byte) (struct {
	CallbackContract common.Address
	RandomnessFee    *big.Int
	Seed             *big.Int
}, error) {
	return _VRFCoordinator.Contract.Callbacks(&_VRFCoordinator.CallOpts, arg0)
}

// HashOfKey is a free data retrieval call binding the contract method 0xcaf70c4a.
//
// Solidity: function hashOfKey(uint256[2] _publicKey) pure returns(bytes32)
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
// Solidity: function hashOfKey(uint256[2] _publicKey) pure returns(bytes32)
func (_VRFCoordinator *VRFCoordinatorSession) HashOfKey(_publicKey [2]*big.Int) ([32]byte, error) {
	return _VRFCoordinator.Contract.HashOfKey(&_VRFCoordinator.CallOpts, _publicKey)
}

// HashOfKey is a free data retrieval call binding the contract method 0xcaf70c4a.
//
// Solidity: function hashOfKey(uint256[2] _publicKey) pure returns(bytes32)
func (_VRFCoordinator *VRFCoordinatorCallerSession) HashOfKey(_publicKey [2]*big.Int) ([32]byte, error) {
	return _VRFCoordinator.Contract.HashOfKey(&_VRFCoordinator.CallOpts, _publicKey)
}

// ServiceAgreements is a free data retrieval call binding the contract method 0x75d35070.
//
// Solidity: function serviceAgreements(bytes32 ) view returns(address vRFOracle, bytes32 jobID, uint256 fee)
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
// Solidity: function serviceAgreements(bytes32 ) view returns(address vRFOracle, bytes32 jobID, uint256 fee)
func (_VRFCoordinator *VRFCoordinatorSession) ServiceAgreements(arg0 [32]byte) (struct {
	VRFOracle common.Address
	JobID     [32]byte
	Fee       *big.Int
}, error) {
	return _VRFCoordinator.Contract.ServiceAgreements(&_VRFCoordinator.CallOpts, arg0)
}

// ServiceAgreements is a free data retrieval call binding the contract method 0x75d35070.
//
// Solidity: function serviceAgreements(bytes32 ) view returns(address vRFOracle, bytes32 jobID, uint256 fee)
func (_VRFCoordinator *VRFCoordinatorCallerSession) ServiceAgreements(arg0 [32]byte) (struct {
	VRFOracle common.Address
	JobID     [32]byte
	Fee       *big.Int
}, error) {
	return _VRFCoordinator.Contract.ServiceAgreements(&_VRFCoordinator.CallOpts, arg0)
}

// WithdrawableTokens is a free data retrieval call binding the contract method 0x006f6ad0.
//
// Solidity: function withdrawableTokens(address ) view returns(uint256)
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
// Solidity: function withdrawableTokens(address ) view returns(uint256)
func (_VRFCoordinator *VRFCoordinatorSession) WithdrawableTokens(arg0 common.Address) (*big.Int, error) {
	return _VRFCoordinator.Contract.WithdrawableTokens(&_VRFCoordinator.CallOpts, arg0)
}

// WithdrawableTokens is a free data retrieval call binding the contract method 0x006f6ad0.
//
// Solidity: function withdrawableTokens(address ) view returns(uint256)
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

// RegisterProvingKey is a paid mutator transaction binding the contract method 0xd8340209.
//
// Solidity: function registerProvingKey(uint256 _fee, address _oracle, uint256[2] _publicProvingKey, bytes32 _jobID) returns()
func (_VRFCoordinator *VRFCoordinatorTransactor) RegisterProvingKey(opts *bind.TransactOpts, _fee *big.Int, _oracle common.Address, _publicProvingKey [2]*big.Int, _jobID [32]byte) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "registerProvingKey", _fee, _oracle, _publicProvingKey, _jobID)
}

// RegisterProvingKey is a paid mutator transaction binding the contract method 0xd8340209.
//
// Solidity: function registerProvingKey(uint256 _fee, address _oracle, uint256[2] _publicProvingKey, bytes32 _jobID) returns()
func (_VRFCoordinator *VRFCoordinatorSession) RegisterProvingKey(_fee *big.Int, _oracle common.Address, _publicProvingKey [2]*big.Int, _jobID [32]byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.RegisterProvingKey(&_VRFCoordinator.TransactOpts, _fee, _oracle, _publicProvingKey, _jobID)
}

// RegisterProvingKey is a paid mutator transaction binding the contract method 0xd8340209.
//
// Solidity: function registerProvingKey(uint256 _fee, address _oracle, uint256[2] _publicProvingKey, bytes32 _jobID) returns()
func (_VRFCoordinator *VRFCoordinatorTransactorSession) RegisterProvingKey(_fee *big.Int, _oracle common.Address, _publicProvingKey [2]*big.Int, _jobID [32]byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.RegisterProvingKey(&_VRFCoordinator.TransactOpts, _fee, _oracle, _publicProvingKey, _jobID)
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
