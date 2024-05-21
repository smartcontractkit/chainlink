// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ocr3_capability

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated"
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

var OCR3CapabilityMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ReportInvalid\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReportingUnsupported\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"_f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"_onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"_offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"_offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"transmitters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b50600133806000816200006b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156200009e576200009e81620000ac565b505050151560805262000157565b336001600160a01b03821603620001065760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000062565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b608051611f7d6200017360003960006104a40152611f7d6000f3fe608060405234801561001057600080fd5b50600436106100a35760003560e01c80638da5cb5b11610076578063b1dc65a41161005b578063b1dc65a414610187578063e3d0e7121461019a578063f2fde38b146101ad57600080fd5b80638da5cb5b1461013f578063afcb95d71461016757600080fd5b8063181f5a77146100a857806379ba5097146100f057806381411834146100fa57806381ff70481461010f575b600080fd5b604080518082018252600e81527f4b657973746f6e6520302e302e30000000000000000000000000000000000000602082015290516100e791906117e8565b60405180910390f35b6100f86101c0565b005b6101026102c2565b6040516100e79190611853565b6004546002546040805163ffffffff808516825264010000000090940490931660208401528201526060016100e7565b60005460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100e7565b6040805160018152600060208201819052918101919091526060016100e7565b6100f86101953660046118b2565b610331565b6100f86101a8366004611b7c565b610a62565b6100f86101bb366004611c49565b61143d565b60015473ffffffffffffffffffffffffffffffffffffffff163314610246576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6060600780548060200260200160405190810160405280929190818152602001828054801561032757602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116102fc575b5050505050905090565b60005a604080516020601f8b018190048102820181019092528981529192508a3591818c01359161038791849163ffffffff851691908e908e908190840183828082843760009201919091525061145192505050565b6103bd576040517f0be3632800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805183815262ffffff600884901c1660208201527fb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62910160405180910390a16040805160608101825260025480825260035460ff80821660208501526101009091041692820192909252908314610492576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f636f6e666967446967657374206d69736d617463680000000000000000000000604482015260640161023d565b6104a08b8b8b8b8b8b61145a565b60007f0000000000000000000000000000000000000000000000000000000000000000156104fd576002826020015183604001516104de9190611cc2565b6104e89190611ce1565b6104f3906001611cc2565b60ff169050610513565b602082015161050d906001611cc2565b60ff1690505b88811461057c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f77726f6e67206e756d626572206f66207369676e617475726573000000000000604482015260640161023d565b8887146105e5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f7369676e617475726573206f7574206f6620726567697374726174696f6e0000604482015260640161023d565b3360009081526005602090815260408083208151808301909252805460ff8082168452929391929184019161010090910416600281111561062857610628611d2a565b600281111561063957610639611d2a565b905250905060028160200151600281111561065657610656611d2a565b14801561069d57506007816000015160ff168154811061067857610678611c64565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff1633145b610703576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f756e617574686f72697a6564207472616e736d69747465720000000000000000604482015260640161023d565b5050505050610710611765565b6000808a8a604051610723929190611d59565b60405190819003812061073a918e90602001611d69565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181528282528051602091820120838301909252600080845290830152915060005b89811015610a445760006001848984602081106107a3576107a3611c64565b6107b091901a601b611cc2565b8e8e868181106107c2576107c2611c64565b905060200201358d8d878181106107db576107db611c64565b9050602002013560405160008152602001604052604051610818949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa15801561083a573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff811660009081526005602090815290849020838501909452835460ff808216855292965092945084019161010090041660028111156108ba576108ba611d2a565b60028111156108cb576108cb611d2a565b90525092506001836020015160028111156108e8576108e8611d2a565b1461094f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f61646472657373206e6f7420617574686f72697a656420746f207369676e0000604482015260640161023d565b8251600090879060ff16601f811061096957610969611c64565b602002015173ffffffffffffffffffffffffffffffffffffffff16146109eb576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f6e6f6e2d756e69717565207369676e6174757265000000000000000000000000604482015260640161023d565b8086846000015160ff16601f8110610a0557610a05611c64565b73ffffffffffffffffffffffffffffffffffffffff9092166020929092020152610a30600186611cc2565b94505080610a3d90611d7d565b9050610784565b505050610a55833383858e8e611511565b5050505050505050505050565b855185518560ff16601f831115610ad5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f746f6f206d616e79207369676e65727300000000000000000000000000000000604482015260640161023d565b60008111610b3f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f66206d75737420626520706f7369746976650000000000000000000000000000604482015260640161023d565b818314610bcd576040517f08c379a0000000000000000000000000000000000000000000000000000000008152602060048201526024808201527f6f7261636c6520616464726573736573206f7574206f6620726567697374726160448201527f74696f6e00000000000000000000000000000000000000000000000000000000606482015260840161023d565b610bd8816003611db5565b8311610c40576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f6661756c74792d6f7261636c65206620746f6f20686967680000000000000000604482015260640161023d565b610c48611543565b6040805160c0810182528a8152602081018a905260ff8916918101919091526060810187905267ffffffffffffffff8616608082015260a081018590525b60065415610e3b57600654600090610ca090600190611dcc565b9050600060068281548110610cb757610cb7611c64565b60009182526020822001546007805473ffffffffffffffffffffffffffffffffffffffff90921693509084908110610cf157610cf1611c64565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff85811684526005909252604080842080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000090811690915592909116808452922080549091169055600680549192509080610d7157610d71611ddf565b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff00000000000000000000000000000000000000001690550190556007805480610dda57610dda611ddf565b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff000000000000000000000000000000000000000016905501905550610c86915050565b60005b8151518110156112a05760006005600084600001518481518110610e6457610e64611c64565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff16825281019190915260400160002054610100900460ff166002811115610eae57610eae611d2a565b14610f15576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265706561746564207369676e65722061646472657373000000000000000000604482015260640161023d565b6040805180820190915260ff82168152600160208201528251805160059160009185908110610f4657610f46611c64565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040016000208251815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082168117835592840151919283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001617610100836002811115610fe757610fe7611d2a565b021790555060009150610ff79050565b600560008460200151848151811061101157611011611c64565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff16825281019190915260400160002054610100900460ff16600281111561105b5761105b611d2a565b146110c2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f7265706561746564207472616e736d6974746572206164647265737300000000604482015260640161023d565b6040805180820190915260ff8216815260208101600281525060056000846020015184815181106110f5576110f5611c64565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040016000208251815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082168117835592840151919283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000161761010083600281111561119657611196611d2a565b0217905550508251805160069250839081106111b4576111b4611c64565b602090810291909101810151825460018101845560009384529282902090920180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909316929092179091558201518051600791908390811061123057611230611c64565b60209081029190910181015182546001810184556000938452919092200180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff90921691909117905561129981611d7d565b9050610e3e565b506040810151600380547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff909216919091179055600480547fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff811664010000000063ffffffff438116820292831785559083048116936001939092600092611332928692908216911617611e0e565b92506101000a81548163ffffffff021916908363ffffffff1602179055506113914630600460009054906101000a900463ffffffff1663ffffffff16856000015186602001518760400151886060015189608001518a60a001516115c6565b6002819055825180516003805460ff909216610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff90921691909117905560045460208501516040808701516060880151608089015160a08a015193517f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e0598611430988b98919763ffffffff909216969095919491939192611e32565b60405180910390a1610a55565b611445611543565b61144e81611670565b50565b60019392505050565b6000611467826020611db5565b611472856020611db5565b61147e88610144611ec8565b6114889190611ec8565b6114929190611ec8565b61149d906000611ec8565b9050368114611508576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f63616c6c64617461206c656e677468206d69736d617463680000000000000000604482015260640161023d565b50505050505050565b6040517f0750181900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005473ffffffffffffffffffffffffffffffffffffffff1633146115c4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161023d565b565b6000808a8a8a8a8a8a8a8a8a6040516020016115ea99989796959493929190611edb565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179b9a5050505050505050505050565b3373ffffffffffffffffffffffffffffffffffffffff8216036116ef576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161023d565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b604051806103e00160405280601f906020820280368337509192915050565b6000815180845260005b818110156117aa5760208185018101518683018201520161178e565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b6020815260006117fb6020830184611784565b9392505050565b600081518084526020808501945080840160005b8381101561184857815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101611816565b509495945050505050565b6020815260006117fb6020830184611802565b60008083601f84011261187857600080fd5b50813567ffffffffffffffff81111561189057600080fd5b6020830191508360208260051b85010111156118ab57600080fd5b9250929050565b60008060008060008060008060e0898b0312156118ce57600080fd5b606089018a8111156118df57600080fd5b8998503567ffffffffffffffff808211156118f957600080fd5b818b0191508b601f83011261190d57600080fd5b81358181111561191c57600080fd5b8c602082850101111561192e57600080fd5b6020830199508098505060808b013591508082111561194c57600080fd5b6119588c838d01611866565b909750955060a08b013591508082111561197157600080fd5b5061197e8b828c01611866565b999c989b50969995989497949560c00135949350505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715611a0d57611a0d611997565b604052919050565b803573ffffffffffffffffffffffffffffffffffffffff81168114611a3957600080fd5b919050565b600082601f830112611a4f57600080fd5b8135602067ffffffffffffffff821115611a6b57611a6b611997565b8160051b611a7a8282016119c6565b9283528481018201928281019087851115611a9457600080fd5b83870192505b84831015611aba57611aab83611a15565b82529183019190830190611a9a565b979650505050505050565b803560ff81168114611a3957600080fd5b600082601f830112611ae757600080fd5b813567ffffffffffffffff811115611b0157611b01611997565b611b3260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116016119c6565b818152846020838601011115611b4757600080fd5b816020850160208301376000918101602001919091529392505050565b803567ffffffffffffffff81168114611a3957600080fd5b60008060008060008060c08789031215611b9557600080fd5b863567ffffffffffffffff80821115611bad57600080fd5b611bb98a838b01611a3e565b97506020890135915080821115611bcf57600080fd5b611bdb8a838b01611a3e565b9650611be960408a01611ac5565b95506060890135915080821115611bff57600080fd5b611c0b8a838b01611ad6565b9450611c1960808a01611b64565b935060a0890135915080821115611c2f57600080fd5b50611c3c89828a01611ad6565b9150509295509295509295565b600060208284031215611c5b57600080fd5b6117fb82611a15565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60ff8181168382160190811115611cdb57611cdb611c93565b92915050565b600060ff831680611d1b577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b8060ff84160491505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b8183823760009101908152919050565b828152606082602083013760800192915050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203611dae57611dae611c93565b5060010190565b8082028115828204841417611cdb57611cdb611c93565b81810381811115611cdb57611cdb611c93565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b63ffffffff818116838216019080821115611e2b57611e2b611c93565b5092915050565b600061012063ffffffff808d1684528b6020850152808b16604085015250806060840152611e628184018a611802565b90508281036080840152611e768189611802565b905060ff871660a084015282810360c0840152611e938187611784565b905067ffffffffffffffff851660e0840152828103610100840152611eb88185611784565b9c9b505050505050505050505050565b80820180821115611cdb57611cdb611c93565b60006101208b835273ffffffffffffffffffffffffffffffffffffffff8b16602084015267ffffffffffffffff808b166040850152816060850152611f228285018b611802565b91508382036080850152611f36828a611802565b915060ff881660a085015283820360c0850152611f538288611784565b90861660e08501528381036101008501529050611eb8818561178456fea164736f6c6343000813000a",
}

var OCR3CapabilityABI = OCR3CapabilityMetaData.ABI

var OCR3CapabilityBin = OCR3CapabilityMetaData.Bin

func DeployOCR3Capability(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *OCR3Capability, error) {
	parsed, err := OCR3CapabilityMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OCR3CapabilityBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OCR3Capability{address: address, abi: *parsed, OCR3CapabilityCaller: OCR3CapabilityCaller{contract: contract}, OCR3CapabilityTransactor: OCR3CapabilityTransactor{contract: contract}, OCR3CapabilityFilterer: OCR3CapabilityFilterer{contract: contract}}, nil
}

type OCR3Capability struct {
	address common.Address
	abi     abi.ABI
	OCR3CapabilityCaller
	OCR3CapabilityTransactor
	OCR3CapabilityFilterer
}

type OCR3CapabilityCaller struct {
	contract *bind.BoundContract
}

type OCR3CapabilityTransactor struct {
	contract *bind.BoundContract
}

type OCR3CapabilityFilterer struct {
	contract *bind.BoundContract
}

type OCR3CapabilitySession struct {
	Contract     *OCR3Capability
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type OCR3CapabilityCallerSession struct {
	Contract *OCR3CapabilityCaller
	CallOpts bind.CallOpts
}

type OCR3CapabilityTransactorSession struct {
	Contract     *OCR3CapabilityTransactor
	TransactOpts bind.TransactOpts
}

type OCR3CapabilityRaw struct {
	Contract *OCR3Capability
}

type OCR3CapabilityCallerRaw struct {
	Contract *OCR3CapabilityCaller
}

type OCR3CapabilityTransactorRaw struct {
	Contract *OCR3CapabilityTransactor
}

func NewOCR3Capability(address common.Address, backend bind.ContractBackend) (*OCR3Capability, error) {
	abi, err := abi.JSON(strings.NewReader(OCR3CapabilityABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindOCR3Capability(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OCR3Capability{address: address, abi: abi, OCR3CapabilityCaller: OCR3CapabilityCaller{contract: contract}, OCR3CapabilityTransactor: OCR3CapabilityTransactor{contract: contract}, OCR3CapabilityFilterer: OCR3CapabilityFilterer{contract: contract}}, nil
}

func NewOCR3CapabilityCaller(address common.Address, caller bind.ContractCaller) (*OCR3CapabilityCaller, error) {
	contract, err := bindOCR3Capability(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OCR3CapabilityCaller{contract: contract}, nil
}

func NewOCR3CapabilityTransactor(address common.Address, transactor bind.ContractTransactor) (*OCR3CapabilityTransactor, error) {
	contract, err := bindOCR3Capability(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OCR3CapabilityTransactor{contract: contract}, nil
}

func NewOCR3CapabilityFilterer(address common.Address, filterer bind.ContractFilterer) (*OCR3CapabilityFilterer, error) {
	contract, err := bindOCR3Capability(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OCR3CapabilityFilterer{contract: contract}, nil
}

func bindOCR3Capability(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OCR3CapabilityMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_OCR3Capability *OCR3CapabilityRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR3Capability.Contract.OCR3CapabilityCaller.contract.Call(opts, result, method, params...)
}

func (_OCR3Capability *OCR3CapabilityRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR3Capability.Contract.OCR3CapabilityTransactor.contract.Transfer(opts)
}

func (_OCR3Capability *OCR3CapabilityRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR3Capability.Contract.OCR3CapabilityTransactor.contract.Transact(opts, method, params...)
}

func (_OCR3Capability *OCR3CapabilityCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR3Capability.Contract.contract.Call(opts, result, method, params...)
}

func (_OCR3Capability *OCR3CapabilityTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR3Capability.Contract.contract.Transfer(opts)
}

func (_OCR3Capability *OCR3CapabilityTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR3Capability.Contract.contract.Transact(opts, method, params...)
}

func (_OCR3Capability *OCR3CapabilityCaller) LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

	error) {
	var out []interface{}
	err := _OCR3Capability.contract.Call(opts, &out, "latestConfigDetails")

	outstruct := new(LatestConfigDetails)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ConfigDigest = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

func (_OCR3Capability *OCR3CapabilitySession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _OCR3Capability.Contract.LatestConfigDetails(&_OCR3Capability.CallOpts)
}

func (_OCR3Capability *OCR3CapabilityCallerSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _OCR3Capability.Contract.LatestConfigDetails(&_OCR3Capability.CallOpts)
}

func (_OCR3Capability *OCR3CapabilityCaller) LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

	error) {
	var out []interface{}
	err := _OCR3Capability.contract.Call(opts, &out, "latestConfigDigestAndEpoch")

	outstruct := new(LatestConfigDigestAndEpoch)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ScanLogs = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfigDigest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Epoch = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

func (_OCR3Capability *OCR3CapabilitySession) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _OCR3Capability.Contract.LatestConfigDigestAndEpoch(&_OCR3Capability.CallOpts)
}

func (_OCR3Capability *OCR3CapabilityCallerSession) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _OCR3Capability.Contract.LatestConfigDigestAndEpoch(&_OCR3Capability.CallOpts)
}

func (_OCR3Capability *OCR3CapabilityCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OCR3Capability.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_OCR3Capability *OCR3CapabilitySession) Owner() (common.Address, error) {
	return _OCR3Capability.Contract.Owner(&_OCR3Capability.CallOpts)
}

func (_OCR3Capability *OCR3CapabilityCallerSession) Owner() (common.Address, error) {
	return _OCR3Capability.Contract.Owner(&_OCR3Capability.CallOpts)
}

func (_OCR3Capability *OCR3CapabilityCaller) Transmitters(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _OCR3Capability.contract.Call(opts, &out, "transmitters")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_OCR3Capability *OCR3CapabilitySession) Transmitters() ([]common.Address, error) {
	return _OCR3Capability.Contract.Transmitters(&_OCR3Capability.CallOpts)
}

func (_OCR3Capability *OCR3CapabilityCallerSession) Transmitters() ([]common.Address, error) {
	return _OCR3Capability.Contract.Transmitters(&_OCR3Capability.CallOpts)
}

func (_OCR3Capability *OCR3CapabilityCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _OCR3Capability.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_OCR3Capability *OCR3CapabilitySession) TypeAndVersion() (string, error) {
	return _OCR3Capability.Contract.TypeAndVersion(&_OCR3Capability.CallOpts)
}

func (_OCR3Capability *OCR3CapabilityCallerSession) TypeAndVersion() (string, error) {
	return _OCR3Capability.Contract.TypeAndVersion(&_OCR3Capability.CallOpts)
}

func (_OCR3Capability *OCR3CapabilityTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR3Capability.contract.Transact(opts, "acceptOwnership")
}

func (_OCR3Capability *OCR3CapabilitySession) AcceptOwnership() (*types.Transaction, error) {
	return _OCR3Capability.Contract.AcceptOwnership(&_OCR3Capability.TransactOpts)
}

func (_OCR3Capability *OCR3CapabilityTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _OCR3Capability.Contract.AcceptOwnership(&_OCR3Capability.TransactOpts)
}

func (_OCR3Capability *OCR3CapabilityTransactor) SetConfig(opts *bind.TransactOpts, _signers []common.Address, _transmitters []common.Address, _f uint8, _onchainConfig []byte, _offchainConfigVersion uint64, _offchainConfig []byte) (*types.Transaction, error) {
	return _OCR3Capability.contract.Transact(opts, "setConfig", _signers, _transmitters, _f, _onchainConfig, _offchainConfigVersion, _offchainConfig)
}

func (_OCR3Capability *OCR3CapabilitySession) SetConfig(_signers []common.Address, _transmitters []common.Address, _f uint8, _onchainConfig []byte, _offchainConfigVersion uint64, _offchainConfig []byte) (*types.Transaction, error) {
	return _OCR3Capability.Contract.SetConfig(&_OCR3Capability.TransactOpts, _signers, _transmitters, _f, _onchainConfig, _offchainConfigVersion, _offchainConfig)
}

func (_OCR3Capability *OCR3CapabilityTransactorSession) SetConfig(_signers []common.Address, _transmitters []common.Address, _f uint8, _onchainConfig []byte, _offchainConfigVersion uint64, _offchainConfig []byte) (*types.Transaction, error) {
	return _OCR3Capability.Contract.SetConfig(&_OCR3Capability.TransactOpts, _signers, _transmitters, _f, _onchainConfig, _offchainConfigVersion, _offchainConfig)
}

func (_OCR3Capability *OCR3CapabilityTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _OCR3Capability.contract.Transact(opts, "transferOwnership", to)
}

func (_OCR3Capability *OCR3CapabilitySession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OCR3Capability.Contract.TransferOwnership(&_OCR3Capability.TransactOpts, to)
}

func (_OCR3Capability *OCR3CapabilityTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OCR3Capability.Contract.TransferOwnership(&_OCR3Capability.TransactOpts, to)
}

func (_OCR3Capability *OCR3CapabilityTransactor) Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OCR3Capability.contract.Transact(opts, "transmit", reportContext, report, rs, ss, rawVs)
}

func (_OCR3Capability *OCR3CapabilitySession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OCR3Capability.Contract.Transmit(&_OCR3Capability.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_OCR3Capability *OCR3CapabilityTransactorSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OCR3Capability.Contract.Transmit(&_OCR3Capability.TransactOpts, reportContext, report, rs, ss, rawVs)
}

type OCR3CapabilityConfigSetIterator struct {
	Event *OCR3CapabilityConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR3CapabilityConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR3CapabilityConfigSet)
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
		it.Event = new(OCR3CapabilityConfigSet)
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

func (it *OCR3CapabilityConfigSetIterator) Error() error {
	return it.fail
}

func (it *OCR3CapabilityConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR3CapabilityConfigSet struct {
	PreviousConfigBlockNumber uint32
	ConfigDigest              [32]byte
	ConfigCount               uint64
	Signers                   []common.Address
	Transmitters              []common.Address
	F                         uint8
	OnchainConfig             []byte
	OffchainConfigVersion     uint64
	OffchainConfig            []byte
	Raw                       types.Log
}

func (_OCR3Capability *OCR3CapabilityFilterer) FilterConfigSet(opts *bind.FilterOpts) (*OCR3CapabilityConfigSetIterator, error) {

	logs, sub, err := _OCR3Capability.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &OCR3CapabilityConfigSetIterator{contract: _OCR3Capability.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_OCR3Capability *OCR3CapabilityFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *OCR3CapabilityConfigSet) (event.Subscription, error) {

	logs, sub, err := _OCR3Capability.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR3CapabilityConfigSet)
				if err := _OCR3Capability.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_OCR3Capability *OCR3CapabilityFilterer) ParseConfigSet(log types.Log) (*OCR3CapabilityConfigSet, error) {
	event := new(OCR3CapabilityConfigSet)
	if err := _OCR3Capability.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR3CapabilityOwnershipTransferRequestedIterator struct {
	Event *OCR3CapabilityOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR3CapabilityOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR3CapabilityOwnershipTransferRequested)
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
		it.Event = new(OCR3CapabilityOwnershipTransferRequested)
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

func (it *OCR3CapabilityOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *OCR3CapabilityOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR3CapabilityOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_OCR3Capability *OCR3CapabilityFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OCR3CapabilityOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR3Capability.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OCR3CapabilityOwnershipTransferRequestedIterator{contract: _OCR3Capability.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_OCR3Capability *OCR3CapabilityFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OCR3CapabilityOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR3Capability.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR3CapabilityOwnershipTransferRequested)
				if err := _OCR3Capability.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_OCR3Capability *OCR3CapabilityFilterer) ParseOwnershipTransferRequested(log types.Log) (*OCR3CapabilityOwnershipTransferRequested, error) {
	event := new(OCR3CapabilityOwnershipTransferRequested)
	if err := _OCR3Capability.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR3CapabilityOwnershipTransferredIterator struct {
	Event *OCR3CapabilityOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR3CapabilityOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR3CapabilityOwnershipTransferred)
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
		it.Event = new(OCR3CapabilityOwnershipTransferred)
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

func (it *OCR3CapabilityOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *OCR3CapabilityOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR3CapabilityOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_OCR3Capability *OCR3CapabilityFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OCR3CapabilityOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR3Capability.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OCR3CapabilityOwnershipTransferredIterator{contract: _OCR3Capability.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_OCR3Capability *OCR3CapabilityFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OCR3CapabilityOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR3Capability.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR3CapabilityOwnershipTransferred)
				if err := _OCR3Capability.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_OCR3Capability *OCR3CapabilityFilterer) ParseOwnershipTransferred(log types.Log) (*OCR3CapabilityOwnershipTransferred, error) {
	event := new(OCR3CapabilityOwnershipTransferred)
	if err := _OCR3Capability.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR3CapabilityTransmittedIterator struct {
	Event *OCR3CapabilityTransmitted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR3CapabilityTransmittedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR3CapabilityTransmitted)
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
		it.Event = new(OCR3CapabilityTransmitted)
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

func (it *OCR3CapabilityTransmittedIterator) Error() error {
	return it.fail
}

func (it *OCR3CapabilityTransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR3CapabilityTransmitted struct {
	ConfigDigest [32]byte
	Epoch        uint32
	Raw          types.Log
}

func (_OCR3Capability *OCR3CapabilityFilterer) FilterTransmitted(opts *bind.FilterOpts) (*OCR3CapabilityTransmittedIterator, error) {

	logs, sub, err := _OCR3Capability.contract.FilterLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return &OCR3CapabilityTransmittedIterator{contract: _OCR3Capability.contract, event: "Transmitted", logs: logs, sub: sub}, nil
}

func (_OCR3Capability *OCR3CapabilityFilterer) WatchTransmitted(opts *bind.WatchOpts, sink chan<- *OCR3CapabilityTransmitted) (event.Subscription, error) {

	logs, sub, err := _OCR3Capability.contract.WatchLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR3CapabilityTransmitted)
				if err := _OCR3Capability.contract.UnpackLog(event, "Transmitted", log); err != nil {
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

func (_OCR3Capability *OCR3CapabilityFilterer) ParseTransmitted(log types.Log) (*OCR3CapabilityTransmitted, error) {
	event := new(OCR3CapabilityTransmitted)
	if err := _OCR3Capability.contract.UnpackLog(event, "Transmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LatestConfigDetails struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}
type LatestConfigDigestAndEpoch struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}

func (_OCR3Capability *OCR3Capability) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _OCR3Capability.abi.Events["ConfigSet"].ID:
		return _OCR3Capability.ParseConfigSet(log)
	case _OCR3Capability.abi.Events["OwnershipTransferRequested"].ID:
		return _OCR3Capability.ParseOwnershipTransferRequested(log)
	case _OCR3Capability.abi.Events["OwnershipTransferred"].ID:
		return _OCR3Capability.ParseOwnershipTransferred(log)
	case _OCR3Capability.abi.Events["Transmitted"].ID:
		return _OCR3Capability.ParseTransmitted(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (OCR3CapabilityConfigSet) Topic() common.Hash {
	return common.HexToHash("0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05")
}

func (OCR3CapabilityOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (OCR3CapabilityOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (OCR3CapabilityTransmitted) Topic() common.Hash {
	return common.HexToHash("0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62")
}

func (_OCR3Capability *OCR3Capability) Address() common.Address {
	return _OCR3Capability.address
}

type OCR3CapabilityInterface interface {
	LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

		error)

	LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

		error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	Transmitters(opts *bind.CallOpts) ([]common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, _signers []common.Address, _transmitters []common.Address, _f uint8, _onchainConfig []byte, _offchainConfigVersion uint64, _offchainConfig []byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error)

	FilterConfigSet(opts *bind.FilterOpts) (*OCR3CapabilityConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *OCR3CapabilityConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*OCR3CapabilityConfigSet, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OCR3CapabilityOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OCR3CapabilityOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*OCR3CapabilityOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OCR3CapabilityOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OCR3CapabilityOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*OCR3CapabilityOwnershipTransferred, error)

	FilterTransmitted(opts *bind.FilterOpts) (*OCR3CapabilityTransmittedIterator, error)

	WatchTransmitted(opts *bind.WatchOpts, sink chan<- *OCR3CapabilityTransmitted) (event.Subscription, error)

	ParseTransmitted(log types.Log) (*OCR3CapabilityTransmitted, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
