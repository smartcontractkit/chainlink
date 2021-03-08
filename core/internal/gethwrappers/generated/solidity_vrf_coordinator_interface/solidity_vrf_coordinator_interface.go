// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package solidity_vrf_coordinator_interface

import (
	"fmt"
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
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

const VRFCoordinatorABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_blockHashStore\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"NewServiceAgreement\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"jobID\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"requestID\",\"type\":\"bytes32\"}],\"name\":\"RandomnessRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"output\",\"type\":\"uint256\"}],\"name\":\"RandomnessRequestFulfilled\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"PRESEED_OFFSET\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PROOF_LENGTH\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PUBLIC_KEY_OFFSET\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"callbacks\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"callbackContract\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"randomnessFee\",\"type\":\"uint96\"},{\"internalType\":\"bytes32\",\"name\":\"seedAndBlockNum\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_proof\",\"type\":\"bytes\"}],\"name\":\"fulfillRandomnessRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"_publicKey\",\"type\":\"uint256[2]\"}],\"name\":\"hashOfKey\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"_publicProvingKey\",\"type\":\"uint256[2]\"},{\"internalType\":\"bytes32\",\"name\":\"_jobID\",\"type\":\"bytes32\"}],\"name\":\"registerProvingKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"serviceAgreements\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"vRFOracle\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"fee\",\"type\":\"uint96\"},{\"internalType\":\"bytes32\",\"name\":\"jobID\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"withdrawableTokens\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

var VRFCoordinatorBin = "0x608060405234801561001057600080fd5b506040516127353803806127358339818101604052604081101561003357600080fd5b508051602090910151600080546001600160a01b039384166001600160a01b031991821617909155600180549390921692169190911790556126bb8061007a6000396000f3fe608060405234801561001057600080fd5b50600436106100c85760003560e01c8063a4c0ed3611610081578063d83402091161005b578063d834020914610359578063e911439c1461039d578063f3fef3a3146103a5576100c8565b8063a4c0ed361461023e578063b415f4f514610306578063caf70c4a1461030e576100c8565b80635e1c1059116100b25780635e1c10591461017157806375d35070146102195780638aa7927b14610236576100c8565b80626f6ad0146100cd57806321f3650914610112575b600080fd5b610100600480360360208110156100e357600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166103de565b60408051918252519081900360200190f35b61012f6004803603602081101561012857600080fd5b50356103f0565b6040805173ffffffffffffffffffffffffffffffffffffffff90941684526bffffffffffffffffffffffff909216602084015282820152519081900360600190f35b6102176004803603602081101561018757600080fd5b8101906020810181356401000000008111156101a257600080fd5b8201836020820111156101b457600080fd5b803590602001918460018302840111640100000000831117156101d657600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610448945050505050565b005b61012f6004803603602081101561022f57600080fd5b5035610550565b6101006105a8565b6102176004803603606081101561025457600080fd5b73ffffffffffffffffffffffffffffffffffffffff8235169160208101359181019060608101604082013564010000000081111561029157600080fd5b8201836020820111156102a357600080fd5b803590602001918460018302840111640100000000831117156102c557600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295506105ad945050505050565b61010061066c565b6101006004803603604081101561032457600080fd5b604080518082018252918301929181830191839060029083908390808284376000920191909152509194506106719350505050565b610217600480360360a081101561036f57600080fd5b5080359073ffffffffffffffffffffffffffffffffffffffff602082013516906040810190608001356106c7565b610100610957565b610217600480360360408110156103bb57600080fd5b5073ffffffffffffffffffffffffffffffffffffffff813516906020013561095d565b60046020526000908152604090205481565b6002602052600090815260409020805460019091015473ffffffffffffffffffffffffffffffffffffffff8216917401000000000000000000000000000000000000000090046bffffffffffffffffffffffff169083565b6000610452612597565b60008061045e85610abf565b6000848152600360209081526040808320548287015173ffffffffffffffffffffffffffffffffffffffff9091168085526004909352922054959950939750919550935090916104c1916bffffffffffffffffffffffff1663ffffffff610e4616565b73ffffffffffffffffffffffffffffffffffffffff82166000908152600460209081526040808320939093558582526002905290812081815560010155835161050d9084908490610ec3565b604080518481526020810184905281517fa2e7a402243ebda4a69ceeb3dfb682943b7a9b3ac66d6eefa8db65894009611c929181900390910190a1505050505050565b6003602052600090815260409020805460019091015473ffffffffffffffffffffffffffffffffffffffff8216917401000000000000000000000000000000000000000090046bffffffffffffffffffffffff169083565b602081565b60005473ffffffffffffffffffffffffffffffffffffffff16331461063357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f4d75737420757365204c494e4b20746f6b656e00000000000000000000000000604482015290519081900360640190fd5b60008082806020019051604081101561064b57600080fd5b508051602090910151909250905061066582828688611080565b5050505050565b60e081565b6000816040516020018082600260200280838360005b8381101561069f578181015183820152602001610687565b505050509050019150506040516020818303038152906040528051906020012090505b919050565b6040805180820182526000916106f6919085906002908390839080828437600092019190915250610671915050565b60008181526003602052604090205490915073ffffffffffffffffffffffffffffffffffffffff16801561078b57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f706c656173652072656769737465722061206e6577206b657900000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff851661080d57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f5f6f7261636c65206d757374206e6f7420626520307830000000000000000000604482015290519081900360640190fd5b600082815260036020526040902080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff87161781556001018390556b033b2e3c9fd0803ce80000008611156108c4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252603c815260200180612630603c913960400191505060405180910390fd5b600082815260036020908152604091829020805473ffffffffffffffffffffffffffffffffffffffff16740100000000000000000000000000000000000000006bffffffffffffffffffffffff8b1602179055815184815290810188905281517fae189157e0628c1e62315e9179156e1ea10e90e9c15060002f7021e907dc2cfe929181900390910190a1505050505050565b6101a081565b3360009081526004602052604090205481908111156109dd57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f63616e2774207769746864726177206d6f7265207468616e2062616c616e6365604482015290519081900360640190fd5b336000908152600460205260409020546109fd908363ffffffff61136116565b33600090815260046020818152604080842094909455825484517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8981169482019490945260248101889052945192169363a9059cbb93604480830194928390030190829087803b158015610a8857600080fd5b505af1158015610a9c573d6000803e3d6000fd5b505050506040513d6020811015610ab257600080fd5b5051610aba57fe5b505050565b6000610ac9612597565b825160009081906101c0908114610b4157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f77726f6e672070726f6f66206c656e6774680000000000000000000000000000604482015290519081900360640190fd5b610b496125b7565b5060e086015181870151602088019190610b6283610671565b9750610b6e88836113d8565b6000818152600260209081526040918290208251606081018452815473ffffffffffffffffffffffffffffffffffffffff8116808352740100000000000000000000000000000000000000009091046bffffffffffffffffffffffff169382019390935260019091015492810192909252909850909650610c5057604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f6e6f20636f72726573706f6e64696e6720726571756573740000000000000000604482015290519081900360640190fd5b6040805160208082018590528183018490528251808303840181526060909201835281519101209088015114610ce757604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f77726f6e672070726553656564206f7220626c6f636b206e756d000000000000604482015290519081900360640190fd5b804080610dfa57600154604080517fe9413d3800000000000000000000000000000000000000000000000000000000815260048101859052905173ffffffffffffffffffffffffffffffffffffffff9092169163e9413d3891602480820192602092909190829003018186803b158015610d6057600080fd5b505afa158015610d74573d6000803e3d6000fd5b505050506040513d6020811015610d8a57600080fd5b5051905080610dfa57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f706c656173652070726f766520626c6f636b6861736800000000000000000000604482015290519081900360640190fd5b6040805160208082018690528183018490528251808303840181526060909201909252805191012060e08b018190526101a08b52610e378b611404565b96505050505050509193509193565b600082820183811015610eba57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b90505b92915050565b604080516024810185905260448082018590528251808303909101815260649091019091526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f94985ddd00000000000000000000000000000000000000000000000000000000179052600090620324b0805a1015610fa657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f6e6f7420656e6f7567682067617320666f7220636f6e73756d65720000000000604482015290519081900360640190fd5b60008473ffffffffffffffffffffffffffffffffffffffff16836040518082805190602001908083835b6020831061100d57805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09092019160209182019101610fd0565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d806000811461106f576040519150601f19603f3d011682016040523d82523d6000602084013e611074565b606091505b50505050505050505050565b600084815260036020526040902054829085907401000000000000000000000000000000000000000090046bffffffffffffffffffffffff1682101561112757604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f42656c6f7720616772656564207061796d656e74000000000000000000000000604482015290519081900360640190fd5b600086815260056020908152604080832073ffffffffffffffffffffffffffffffffffffffff871684529091528120549061116488888785611572565b9050600061117289836113d8565b60008181526002602052604090205490915073ffffffffffffffffffffffffffffffffffffffff16156111a157fe5b600081815260026020526040902080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff88161790556b033b2e3c9fd0803ce8000000871061120257fe5b600081815260026020908152604080832080546bffffffffffffffffffffffff8c16740100000000000000000000000000000000000000000273ffffffffffffffffffffffffffffffffffffffff91821617825582518085018890524381850152835180820385018152606082018086528151918701919091206001948501558f875260039095529483902090910154928d905260808401869052891660a084015260c083018a905260e083018490525190917f56bd374744a66d531874338def36c906e3a6cf31176eb1e9afd9f1de69725d5191908190036101000190a2600089815260056020908152604080832073ffffffffffffffffffffffffffffffffffffffff8a16845290915290205461132290600163ffffffff610e4616565b6000998a52600560209081526040808c2073ffffffffffffffffffffffffffffffffffffffff9099168c52979052959098209490945550505050505050565b6000828211156113d257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f536166654d6174683a207375627472616374696f6e206f766572666c6f770000604482015290519081900360640190fd5b50900390565b604080516020808201949094528082019290925280518083038201815260609092019052805191012090565b60006101a082511461147757604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f77726f6e672070726f6f66206c656e6774680000000000000000000000000000604482015290519081900360640190fd5b61147f6125b7565b6114876125b7565b61148f6125d5565b60006114996125b7565b6114a16125b7565b6000888060200190516101a08110156114b957600080fd5b5060e081015161018082015191985060408901975060808901965094506101008801935061014088019250905061150c8787876000602002015188600160200201518960026020020151898989896115c6565b6003866040516020018083815260200182600260200280838360005b83811015611540578181015183820152602001611528565b50505050905001925050506040516020818303038152906040528051906020012060001c975050505050505050919050565b604080516020808201969096528082019490945273ffffffffffffffffffffffffffffffffffffffff9290921660608401526080808401919091528151808403909101815260a09092019052805191012090565b6115cf896118c7565b61163a57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f7075626c6963206b6579206973206e6f74206f6e206375727665000000000000604482015290519081900360640190fd5b611643886118c7565b6116ae57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f67616d6d61206973206e6f74206f6e2063757276650000000000000000000000604482015290519081900360640190fd5b6116b7836118c7565b61172257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f6347616d6d615769746e657373206973206e6f74206f6e206375727665000000604482015290519081900360640190fd5b61172b826118c7565b61179657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f73486173685769746e657373206973206e6f74206f6e20637572766500000000604482015290519081900360640190fd5b6117a2878a888761190b565b61180d57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f6164647228632a706b2b732a6729e289a05f755769746e657373000000000000604482015290519081900360640190fd5b6118156125b7565b61181f8a87611ad7565b90506118296125b7565b611838898b878b868989611b7a565b90506000611849838d8d8a86611ced565b9050808a146118b957604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f696e76616c69642070726f6f6600000000000000000000000000000000000000604482015290519081900360640190fd5b505050505050505050505050565b60208101516000907ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f9080096119048360005b6020020151611e10565b1492915050565b600073ffffffffffffffffffffffffffffffffffffffff821661198f57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f626164207769746e657373000000000000000000000000000000000000000000604482015290519081900360640190fd5b6020840151600090600116156119a657601c6119a9565b601b5b905060007ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd03641418587600060200201510986517ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd036414191820392506000919089098751604080516000808252602082810180855288905260ff8916838501526060830194909452608082018590529151939450909260019260a08084019391927fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081019281900390910190855afa158015611a84573d6000803e3d6000fd5b50506040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0015173ffffffffffffffffffffffffffffffffffffffff9081169088161495505050505050949350505050565b611adf6125b7565b611b3d600184846040516020018084815260200183600260200280838360005b83811015611b17578181015183820152602001611aff565b505050509050018281526020019350505050604051602081830303815290604052611e68565b90505b611b49816118c7565b610ebd578051604080516020818101939093528151808203909301835281019052611b7390611e68565b9050611b40565b611b826125b7565b825186517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f91900306611c1657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f706f696e747320696e2073756d206d7573742062652064697374696e63740000604482015290519081900360640190fd5b611c21878988611ed0565b611c76576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602181526020018061266c6021913960400191505060405180910390fd5b611c81848685611ed0565b611cd6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602281526020018061268d6022913960400191505060405180910390fd5b611ce1868484612036565b98975050505050505050565b6000600286868685876040516020018087815260200186600260200280838360005b83811015611d27578181015183820152602001611d0f565b5050505090500185600260200280838360005b83811015611d52578181015183820152602001611d3a565b5050505090500184600260200280838360005b83811015611d7d578181015183820152602001611d65565b5050505090500183600260200280838360005b83811015611da8578181015183820152602001611d90565b505050509050018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1660601b815260140196505050505050506040516020818303038152906040528051906020012060001c905095945050505050565b6000807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80848509840990507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f600782089392505050565b611e706125b7565b611e7982612164565b8152611e8e611e898260006118fa565b6121b9565b6020820181905260029006600114156106c2576020810180517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f039052919050565b600082611edc57600080fd5b8351602085015160009060011615611ef557601c611ef8565b601b5b905060007ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141838709604080516000808252602080830180855282905260ff871683850152606083018890526080830185905292519394509260019260a08084019391927fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081019281900390910190855afa158015611f9a573d6000803e3d6000fd5b5050506020604051035190506000866040516020018082600260200280838360005b83811015611fd4578181015183820152602001611fbc565b505050509050019150506040516020818303038152906040528051906020012060001c90508073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614955050505050509392505050565b61203e6125b7565b83516020808601518551918601516000938493849361205f939091906121e5565b919450925090507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8582096001146120f857604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f696e765a206d75737420626520696e7665727365206f66207a00000000000000604482015290519081900360640190fd5b60405180604001604052807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8061212b57fe5b87860981526020017ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8785099052979650505050505050565b805160208201205b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f81106106c25760408051602080820193909352815180820384018152908201909152805191012061216c565b6000610ebd827f3fffffffffffffffffffffffffffffffffffffffffffffffffffffffbfffff0c61237b565b60008080600180827ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f897ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f038808905060007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f038a089050600061228d83838585612474565b909850905061229e88828e886124cc565b90985090506122af88828c876124cc565b909850905060006122c28d878b856124cc565b90985090506122d388828686612474565b90985090506122e488828e896124cc565b9098509050818114612367577ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f818a0998507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f82890997507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f818309965061236b565b8196505b5050505050509450945094915050565b6000806123866125f3565b6020808252818101819052604082015260608101859052608081018490527ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f60a08201526123d2612611565b60208160c08460057ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa92508261246a57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f6269674d6f64457870206661696c757265210000000000000000000000000000604482015290519081900360640190fd5b5195945050505050565b6000807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8487097ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8487099097909650945050505050565b600080807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f878509905060007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f87877ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f030990507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8183087ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f86890990999098509650505050505050565b604080516060810182526000808252602082018190529181019190915290565b60405180604001604052806002906020820280368337509192915050565b60405180606001604052806003906020820280368337509192915050565b6040518060c001604052806006906020820280368337509192915050565b6040518060200160405280600190602082028036833750919291505056fe796f752063616e277420636861726765206d6f7265207468616e20616c6c20746865204c494e4b20696e2074686520776f726c642c206772656564794669727374206d756c7469706c69636174696f6e20636865636b206661696c65645365636f6e64206d756c7469706c69636174696f6e20636865636b206661696c6564a164736f6c6343000606000a"

func DeployVRFCoordinator(auth *bind.TransactOpts, backend bind.ContractBackend, _link common.Address, _blockHashStore common.Address) (common.Address, *types.Transaction, *VRFCoordinator, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFCoordinatorABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(VRFCoordinatorBin), backend, _link, _blockHashStore)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFCoordinator{VRFCoordinatorCaller: VRFCoordinatorCaller{contract: contract}, VRFCoordinatorTransactor: VRFCoordinatorTransactor{contract: contract}, VRFCoordinatorFilterer: VRFCoordinatorFilterer{contract: contract}}, nil
}

type VRFCoordinator struct {
	address common.Address
	VRFCoordinatorCaller
	VRFCoordinatorTransactor
	VRFCoordinatorFilterer
}

type VRFCoordinatorCaller struct {
	contract *bind.BoundContract
}

type VRFCoordinatorTransactor struct {
	contract *bind.BoundContract
}

type VRFCoordinatorFilterer struct {
	contract *bind.BoundContract
}

type VRFCoordinatorSession struct {
	Contract     *VRFCoordinator
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFCoordinatorCallerSession struct {
	Contract *VRFCoordinatorCaller
	CallOpts bind.CallOpts
}

type VRFCoordinatorTransactorSession struct {
	Contract     *VRFCoordinatorTransactor
	TransactOpts bind.TransactOpts
}

type VRFCoordinatorRaw struct {
	Contract *VRFCoordinator
}

type VRFCoordinatorCallerRaw struct {
	Contract *VRFCoordinatorCaller
}

type VRFCoordinatorTransactorRaw struct {
	Contract *VRFCoordinatorTransactor
}

func NewVRFCoordinator(address common.Address, backend bind.ContractBackend) (*VRFCoordinator, error) {
	contract, err := bindVRFCoordinator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinator{address: address, VRFCoordinatorCaller: VRFCoordinatorCaller{contract: contract}, VRFCoordinatorTransactor: VRFCoordinatorTransactor{contract: contract}, VRFCoordinatorFilterer: VRFCoordinatorFilterer{contract: contract}}, nil
}

func NewVRFCoordinatorCaller(address common.Address, caller bind.ContractCaller) (*VRFCoordinatorCaller, error) {
	contract, err := bindVRFCoordinator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorCaller{contract: contract}, nil
}

func NewVRFCoordinatorTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFCoordinatorTransactor, error) {
	contract, err := bindVRFCoordinator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTransactor{contract: contract}, nil
}

func NewVRFCoordinatorFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFCoordinatorFilterer, error) {
	contract, err := bindVRFCoordinator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorFilterer{contract: contract}, nil
}

func bindVRFCoordinator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFCoordinatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_VRFCoordinator *VRFCoordinatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFCoordinator.Contract.VRFCoordinatorCaller.contract.Call(opts, result, method, params...)
}

func (_VRFCoordinator *VRFCoordinatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.VRFCoordinatorTransactor.contract.Transfer(opts)
}

func (_VRFCoordinator *VRFCoordinatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.VRFCoordinatorTransactor.contract.Transact(opts, method, params...)
}

func (_VRFCoordinator *VRFCoordinatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFCoordinator.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFCoordinator *VRFCoordinatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.contract.Transfer(opts)
}

func (_VRFCoordinator *VRFCoordinatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.contract.Transact(opts, method, params...)
}

func (_VRFCoordinator *VRFCoordinatorCaller) PRESEEDOFFSET(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "PRESEED_OFFSET")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinator *VRFCoordinatorSession) PRESEEDOFFSET() (*big.Int, error) {
	return _VRFCoordinator.Contract.PRESEEDOFFSET(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) PRESEEDOFFSET() (*big.Int, error) {
	return _VRFCoordinator.Contract.PRESEEDOFFSET(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCaller) PROOFLENGTH(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "PROOF_LENGTH")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinator *VRFCoordinatorSession) PROOFLENGTH() (*big.Int, error) {
	return _VRFCoordinator.Contract.PROOFLENGTH(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) PROOFLENGTH() (*big.Int, error) {
	return _VRFCoordinator.Contract.PROOFLENGTH(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCaller) PUBLICKEYOFFSET(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "PUBLIC_KEY_OFFSET")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinator *VRFCoordinatorSession) PUBLICKEYOFFSET() (*big.Int, error) {
	return _VRFCoordinator.Contract.PUBLICKEYOFFSET(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) PUBLICKEYOFFSET() (*big.Int, error) {
	return _VRFCoordinator.Contract.PUBLICKEYOFFSET(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCaller) Callbacks(opts *bind.CallOpts, arg0 [32]byte) (Callbacks,

	error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "callbacks", arg0)

	outstruct := new(Callbacks)
	if err != nil {
		return *outstruct, err
	}

	outstruct.CallbackContract = out[0].(common.Address)
	outstruct.RandomnessFee = out[1].(*big.Int)
	outstruct.SeedAndBlockNum = out[2].([32]byte)

	return *outstruct, err

}

func (_VRFCoordinator *VRFCoordinatorSession) Callbacks(arg0 [32]byte) (Callbacks,

	error) {
	return _VRFCoordinator.Contract.Callbacks(&_VRFCoordinator.CallOpts, arg0)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) Callbacks(arg0 [32]byte) (Callbacks,

	error) {
	return _VRFCoordinator.Contract.Callbacks(&_VRFCoordinator.CallOpts, arg0)
}

func (_VRFCoordinator *VRFCoordinatorCaller) HashOfKey(opts *bind.CallOpts, _publicKey [2]*big.Int) ([32]byte, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "hashOfKey", _publicKey)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VRFCoordinator *VRFCoordinatorSession) HashOfKey(_publicKey [2]*big.Int) ([32]byte, error) {
	return _VRFCoordinator.Contract.HashOfKey(&_VRFCoordinator.CallOpts, _publicKey)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) HashOfKey(_publicKey [2]*big.Int) ([32]byte, error) {
	return _VRFCoordinator.Contract.HashOfKey(&_VRFCoordinator.CallOpts, _publicKey)
}

func (_VRFCoordinator *VRFCoordinatorCaller) ServiceAgreements(opts *bind.CallOpts, arg0 [32]byte) (ServiceAgreements,

	error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "serviceAgreements", arg0)

	outstruct := new(ServiceAgreements)
	if err != nil {
		return *outstruct, err
	}

	outstruct.VRFOracle = out[0].(common.Address)
	outstruct.Fee = out[1].(*big.Int)
	outstruct.JobID = out[2].([32]byte)

	return *outstruct, err

}

func (_VRFCoordinator *VRFCoordinatorSession) ServiceAgreements(arg0 [32]byte) (ServiceAgreements,

	error) {
	return _VRFCoordinator.Contract.ServiceAgreements(&_VRFCoordinator.CallOpts, arg0)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) ServiceAgreements(arg0 [32]byte) (ServiceAgreements,

	error) {
	return _VRFCoordinator.Contract.ServiceAgreements(&_VRFCoordinator.CallOpts, arg0)
}

func (_VRFCoordinator *VRFCoordinatorCaller) WithdrawableTokens(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "withdrawableTokens", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinator *VRFCoordinatorSession) WithdrawableTokens(arg0 common.Address) (*big.Int, error) {
	return _VRFCoordinator.Contract.WithdrawableTokens(&_VRFCoordinator.CallOpts, arg0)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) WithdrawableTokens(arg0 common.Address) (*big.Int, error) {
	return _VRFCoordinator.Contract.WithdrawableTokens(&_VRFCoordinator.CallOpts, arg0)
}

func (_VRFCoordinator *VRFCoordinatorTransactor) FulfillRandomnessRequest(opts *bind.TransactOpts, _proof []byte) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "fulfillRandomnessRequest", _proof)
}

func (_VRFCoordinator *VRFCoordinatorSession) FulfillRandomnessRequest(_proof []byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.FulfillRandomnessRequest(&_VRFCoordinator.TransactOpts, _proof)
}

func (_VRFCoordinator *VRFCoordinatorTransactorSession) FulfillRandomnessRequest(_proof []byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.FulfillRandomnessRequest(&_VRFCoordinator.TransactOpts, _proof)
}

func (_VRFCoordinator *VRFCoordinatorTransactor) OnTokenTransfer(opts *bind.TransactOpts, _sender common.Address, _fee *big.Int, _data []byte) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "onTokenTransfer", _sender, _fee, _data)
}

func (_VRFCoordinator *VRFCoordinatorSession) OnTokenTransfer(_sender common.Address, _fee *big.Int, _data []byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.OnTokenTransfer(&_VRFCoordinator.TransactOpts, _sender, _fee, _data)
}

func (_VRFCoordinator *VRFCoordinatorTransactorSession) OnTokenTransfer(_sender common.Address, _fee *big.Int, _data []byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.OnTokenTransfer(&_VRFCoordinator.TransactOpts, _sender, _fee, _data)
}

func (_VRFCoordinator *VRFCoordinatorTransactor) RegisterProvingKey(opts *bind.TransactOpts, _fee *big.Int, _oracle common.Address, _publicProvingKey [2]*big.Int, _jobID [32]byte) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "registerProvingKey", _fee, _oracle, _publicProvingKey, _jobID)
}

func (_VRFCoordinator *VRFCoordinatorSession) RegisterProvingKey(_fee *big.Int, _oracle common.Address, _publicProvingKey [2]*big.Int, _jobID [32]byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.RegisterProvingKey(&_VRFCoordinator.TransactOpts, _fee, _oracle, _publicProvingKey, _jobID)
}

func (_VRFCoordinator *VRFCoordinatorTransactorSession) RegisterProvingKey(_fee *big.Int, _oracle common.Address, _publicProvingKey [2]*big.Int, _jobID [32]byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.RegisterProvingKey(&_VRFCoordinator.TransactOpts, _fee, _oracle, _publicProvingKey, _jobID)
}

func (_VRFCoordinator *VRFCoordinatorTransactor) Withdraw(opts *bind.TransactOpts, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "withdraw", _recipient, _amount)
}

func (_VRFCoordinator *VRFCoordinatorSession) Withdraw(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.Withdraw(&_VRFCoordinator.TransactOpts, _recipient, _amount)
}

func (_VRFCoordinator *VRFCoordinatorTransactorSession) Withdraw(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.Withdraw(&_VRFCoordinator.TransactOpts, _recipient, _amount)
}

type VRFCoordinatorNewServiceAgreementIterator struct {
	Event *VRFCoordinatorNewServiceAgreement

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorNewServiceAgreementIterator) Next() bool {

	if it.fail != nil {
		return false
	}

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

func (it *VRFCoordinatorNewServiceAgreementIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorNewServiceAgreementIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorNewServiceAgreement struct {
	KeyHash [32]byte
	Fee     *big.Int
	Raw     types.Log
}

func (_VRFCoordinator *VRFCoordinatorFilterer) FilterNewServiceAgreement(opts *bind.FilterOpts) (*VRFCoordinatorNewServiceAgreementIterator, error) {

	logs, sub, err := _VRFCoordinator.contract.FilterLogs(opts, "NewServiceAgreement")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorNewServiceAgreementIterator{contract: _VRFCoordinator.contract, event: "NewServiceAgreement", logs: logs, sub: sub}, nil
}

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

func (_VRFCoordinator *VRFCoordinatorFilterer) ParseNewServiceAgreement(log types.Log) (*VRFCoordinatorNewServiceAgreement, error) {
	event := new(VRFCoordinatorNewServiceAgreement)
	if err := _VRFCoordinator.contract.UnpackLog(event, "NewServiceAgreement", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorRandomnessRequestIterator struct {
	Event *VRFCoordinatorRandomnessRequest

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorRandomnessRequestIterator) Next() bool {

	if it.fail != nil {
		return false
	}

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

func (it *VRFCoordinatorRandomnessRequestIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorRandomnessRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorRandomnessRequest struct {
	KeyHash   [32]byte
	Seed      *big.Int
	JobID     [32]byte
	Sender    common.Address
	Fee       *big.Int
	RequestID [32]byte
	Raw       types.Log
}

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

func (_VRFCoordinator *VRFCoordinatorFilterer) ParseRandomnessRequest(log types.Log) (*VRFCoordinatorRandomnessRequest, error) {
	event := new(VRFCoordinatorRandomnessRequest)
	if err := _VRFCoordinator.contract.UnpackLog(event, "RandomnessRequest", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorRandomnessRequestFulfilledIterator struct {
	Event *VRFCoordinatorRandomnessRequestFulfilled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorRandomnessRequestFulfilledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorRandomnessRequestFulfilled)
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
		it.Event = new(VRFCoordinatorRandomnessRequestFulfilled)
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

func (it *VRFCoordinatorRandomnessRequestFulfilledIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorRandomnessRequestFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorRandomnessRequestFulfilled struct {
	RequestId [32]byte
	Output    *big.Int
	Raw       types.Log
}

func (_VRFCoordinator *VRFCoordinatorFilterer) FilterRandomnessRequestFulfilled(opts *bind.FilterOpts) (*VRFCoordinatorRandomnessRequestFulfilledIterator, error) {

	logs, sub, err := _VRFCoordinator.contract.FilterLogs(opts, "RandomnessRequestFulfilled")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorRandomnessRequestFulfilledIterator{contract: _VRFCoordinator.contract, event: "RandomnessRequestFulfilled", logs: logs, sub: sub}, nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) WatchRandomnessRequestFulfilled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorRandomnessRequestFulfilled) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinator.contract.WatchLogs(opts, "RandomnessRequestFulfilled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorRandomnessRequestFulfilled)
				if err := _VRFCoordinator.contract.UnpackLog(event, "RandomnessRequestFulfilled", log); err != nil {
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

func (_VRFCoordinator *VRFCoordinatorFilterer) ParseRandomnessRequestFulfilled(log types.Log) (*VRFCoordinatorRandomnessRequestFulfilled, error) {
	event := new(VRFCoordinatorRandomnessRequestFulfilled)
	if err := _VRFCoordinator.contract.UnpackLog(event, "RandomnessRequestFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type Callbacks struct {
	CallbackContract common.Address
	RandomnessFee    *big.Int
	SeedAndBlockNum  [32]byte
}
type ServiceAgreements struct {
	VRFOracle common.Address
	Fee       *big.Int
	JobID     [32]byte
}

func (_VRFCoordinator *VRFCoordinator) UnpackLog(out interface{}, event string, log types.Log) error {
	return _VRFCoordinator.VRFCoordinatorFilterer.contract.UnpackLog(out, event, log)
}

func (_VRFCoordinator *VRFCoordinator) ParseLog(log types.Log) (interface{}, error) {
	abi, err := abi.JSON(strings.NewReader(VRFCoordinatorABI))
	if err != nil {
		return nil, fmt.Errorf("could not parse ABI: " + err.Error())
	}
	switch log.Topics[0] {
	case abi.Events["NewServiceAgreement"].ID:
		return _VRFCoordinator.ParseNewServiceAgreement(log)
	case abi.Events["RandomnessRequest"].ID:
		return _VRFCoordinator.ParseRandomnessRequest(log)
	case abi.Events["RandomnessRequestFulfilled"].ID:
		return _VRFCoordinator.ParseRandomnessRequestFulfilled(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (_VRFCoordinator *VRFCoordinator) Address() common.Address {
	return _VRFCoordinator.address
}

type VRFCoordinatorInterface interface {
	PRESEEDOFFSET(opts *bind.CallOpts) (*big.Int, error)

	PROOFLENGTH(opts *bind.CallOpts) (*big.Int, error)

	PUBLICKEYOFFSET(opts *bind.CallOpts) (*big.Int, error)

	Callbacks(opts *bind.CallOpts, arg0 [32]byte) (Callbacks,

		error)

	HashOfKey(opts *bind.CallOpts, _publicKey [2]*big.Int) ([32]byte, error)

	ServiceAgreements(opts *bind.CallOpts, arg0 [32]byte) (ServiceAgreements,

		error)

	WithdrawableTokens(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error)

	FulfillRandomnessRequest(opts *bind.TransactOpts, _proof []byte) (*types.Transaction, error)

	OnTokenTransfer(opts *bind.TransactOpts, _sender common.Address, _fee *big.Int, _data []byte) (*types.Transaction, error)

	RegisterProvingKey(opts *bind.TransactOpts, _fee *big.Int, _oracle common.Address, _publicProvingKey [2]*big.Int, _jobID [32]byte) (*types.Transaction, error)

	Withdraw(opts *bind.TransactOpts, _recipient common.Address, _amount *big.Int) (*types.Transaction, error)

	FilterNewServiceAgreement(opts *bind.FilterOpts) (*VRFCoordinatorNewServiceAgreementIterator, error)

	WatchNewServiceAgreement(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorNewServiceAgreement) (event.Subscription, error)

	ParseNewServiceAgreement(log types.Log) (*VRFCoordinatorNewServiceAgreement, error)

	FilterRandomnessRequest(opts *bind.FilterOpts, jobID [][32]byte) (*VRFCoordinatorRandomnessRequestIterator, error)

	WatchRandomnessRequest(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorRandomnessRequest, jobID [][32]byte) (event.Subscription, error)

	ParseRandomnessRequest(log types.Log) (*VRFCoordinatorRandomnessRequest, error)

	FilterRandomnessRequestFulfilled(opts *bind.FilterOpts) (*VRFCoordinatorRandomnessRequestFulfilledIterator, error)

	WatchRandomnessRequestFulfilled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorRandomnessRequestFulfilled) (event.Subscription, error)

	ParseRandomnessRequestFulfilled(log types.Log) (*VRFCoordinatorRandomnessRequestFulfilled, error)

	UnpackLog(out interface{}, event string, log types.Log) error

	ParseLog(log types.Log) (interface{}, error)

	Address() common.Address
}
