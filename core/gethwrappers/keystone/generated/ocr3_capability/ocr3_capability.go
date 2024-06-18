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
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"InvalidConfig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"ReportInvalid\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReportingUnsupported\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"_f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"_onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"_offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"_offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"transmitters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b503380600081620000695760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156200009c576200009c81620000a5565b50505062000150565b336001600160a01b03821603620000ff5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000060565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b61202980620001606000396000f3fe608060405234801561001057600080fd5b50600436106100a35760003560e01c80638da5cb5b11610076578063b1dc65a41161005b578063b1dc65a4146101c4578063e3d0e712146101d7578063f2fde38b146101ea57600080fd5b80638da5cb5b1461017c578063afcb95d7146101a457600080fd5b8063181f5a77146100a857806379ba5097146100f057806381411834146100fa57806381ff70481461010f575b600080fd5b604080518082018252600e81527f4b657973746f6e6520312e302e30000000000000000000000000000000000000602082015290516100e79190611894565b60405180910390f35b6100f86101fd565b005b6101026102ff565b6040516100e791906118ff565b61015960015460025463ffffffff74010000000000000000000000000000000000000000830481169378010000000000000000000000000000000000000000000000009093041691565b6040805163ffffffff9485168152939092166020840152908201526060016100e7565b60005460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100e7565b6040805160018152600060208201819052918101919091526060016100e7565b6100f86101d236600461195e565b61036e565b6100f86101e5366004611c28565b61097e565b6100f86101f8366004611cf5565b6114f1565b60015473ffffffffffffffffffffffffffffffffffffffff163314610283576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6060600680548060200260200160405190810160405280929190818152602001828054801561036457602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610339575b5050505050905090565b60005a604080518b3580825262ffffff6020808f0135600881901c929092169084015293945092917fb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62910160405180910390a16103cf8a8a8a8a8a8a611505565b6003546000906002906103ed9060ff80821691610100900416611d6e565b6103f79190611d8d565b610402906001611d6e565b60ff169050878114610470576040517f660bd4ba00000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f77726f6e67206e756d626572206f66207369676e617475726573000000000000604482015260640161027a565b8786146104ff576040517f660bd4ba00000000000000000000000000000000000000000000000000000000815260206004820152602860248201527f7265706f727420727320616e64207373206d757374206265206f66206571756160448201527f6c206c656e677468000000000000000000000000000000000000000000000000606482015260840161027a565b3360009081526004602090815260408083208151808301909252805460ff8082168452929391929184019161010090910416600281111561054257610542611dd6565b600281111561055357610553611dd6565b905250905060028160200151600281111561057057610570611dd6565b141580156105b957506006816000015160ff168154811061059357610593611d10565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff163314155b15610620576040517f660bd4ba00000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f756e617574686f72697a6564207472616e736d69747465720000000000000000604482015260640161027a565b5050505061062c611811565b6000808a8a60405161063f929190611e05565b604051908190038120610656918e90602001611e15565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181528282528051602091820120838301909252600080845290830152915060005b898110156109605760006001848984602081106106bf576106bf611d10565b6106cc91901a601b611d6e565b8e8e868181106106de576106de611d10565b905060200201358d8d878181106106f7576106f7611d10565b9050602002013560405160008152602001604052604051610734949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa158015610756573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff811660009081526004602090815290849020838501909452835460ff808216855292965092945084019161010090041660028111156107d6576107d6611dd6565b60028111156107e7576107e7611dd6565b905250925060018360200151600281111561080457610804611dd6565b1461086b576040517f660bd4ba00000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f61646472657373206e6f7420617574686f72697a656420746f207369676e0000604482015260640161027a565b8251600090879060ff16601f811061088557610885611d10565b602002015173ffffffffffffffffffffffffffffffffffffffff1614610907576040517f660bd4ba00000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f6e6f6e2d756e69717565207369676e6174757265000000000000000000000000604482015260640161027a565b8086846000015160ff16601f811061092157610921611d10565b73ffffffffffffffffffffffffffffffffffffffff909216602092909202015261094c600186611d6e565b9450508061095990611e29565b90506106a0565b505050610971833383858e8e6115bc565b5050505050505050505050565b855185518560ff16601f8311156109f1576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f746f6f206d616e79207369676e65727300000000000000000000000000000000604482015260640161027a565b80600003610a5b576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f66206d75737420626520706f7369746976650000000000000000000000000000604482015260640161027a565b818314610ae9576040517f89a61989000000000000000000000000000000000000000000000000000000008152602060048201526024808201527f6f7261636c6520616464726573736573206f7574206f6620726567697374726160448201527f74696f6e00000000000000000000000000000000000000000000000000000000606482015260840161027a565b610af4816003611e61565b8311610b5c576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f6661756c74792d6f7261636c65206620746f6f20686967680000000000000000604482015260640161027a565b610b646115ee565b6040805160c0810182528a8152602081018a905260ff8916918101919091526060810187905267ffffffffffffffff8616608082015260a081018590525b60055415610d5757600554600090610bbc90600190611e78565b9050600060058281548110610bd357610bd3611d10565b60009182526020822001546006805473ffffffffffffffffffffffffffffffffffffffff90921693509084908110610c0d57610c0d611d10565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff85811684526004909252604080842080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000090811690915592909116808452922080549091169055600580549192509080610c8d57610c8d611e8b565b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff00000000000000000000000000000000000000001690550190556006805480610cf657610cf6611e8b565b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff000000000000000000000000000000000000000016905501905550610ba2915050565b60005b81515181101561130e57815180516000919083908110610d7c57610d7c611d10565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1603610e01576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f7369676e6572206d757374206e6f7420626520656d7074790000000000000000604482015260640161027a565b600073ffffffffffffffffffffffffffffffffffffffff1682602001518281518110610e2f57610e2f611d10565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1603610eb4576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f7472616e736d6974746572206d757374206e6f7420626520656d707479000000604482015260640161027a565b60006004600084600001518481518110610ed057610ed0611d10565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff16825281019190915260400160002054610100900460ff166002811115610f1a57610f1a611dd6565b14610f81576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265706561746564207369676e65722061646472657373000000000000000000604482015260640161027a565b6040805180820190915260ff82168152600160208201528251805160049160009185908110610fb257610fb2611d10565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040016000208251815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082168117835592840151919283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000161761010083600281111561105357611053611dd6565b0217905550600091506110639050565b600460008460200151848151811061107d5761107d611d10565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff16825281019190915260400160002054610100900460ff1660028111156110c7576110c7611dd6565b1461112e576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f7265706561746564207472616e736d6974746572206164647265737300000000604482015260640161027a565b6040805180820190915260ff82168152602081016002815250600460008460200151848151811061116157611161611d10565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040016000208251815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082168117835592840151919283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000161761010083600281111561120257611202611dd6565b02179055505082518051600592508390811061122057611220611d10565b602090810291909101810151825460018101845560009384529282902090920180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909316929092179091558201518051600691908390811061129c5761129c611d10565b60209081029190910181015182546001810184556000938452919092200180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff9092169190911790558061130681611e29565b915050610d5a565b506040810151600380547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff909216919091179055600180547fffffffff00000000ffffffffffffffffffffffffffffffffffffffffffffffff8116780100000000000000000000000000000000000000000000000063ffffffff43811682029290921780855592048116929182916014916113c691849174010000000000000000000000000000000000000000900416611eba565b92506101000a81548163ffffffff021916908363ffffffff1602179055506114254630600160149054906101000a900463ffffffff1663ffffffff16856000015186602001518760400151886060015189608001518a60a00151611671565b600281905582518051600380547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1661010060ff9093169290920291909117905560015460208501516040808701516060880151608089015160a08a015193517f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05986114dc988b9891977401000000000000000000000000000000000000000090920463ffffffff16969095919491939192611ede565b60405180910390a15050505050505050505050565b6114f96115ee565b6115028161171c565b50565b6000611512826020611e61565b61151d856020611e61565b61152988610144611f74565b6115339190611f74565b61153d9190611f74565b611548906000611f74565b90503681146115b3576040517f660bd4ba00000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f63616c6c64617461206c656e677468206d69736d617463680000000000000000604482015260640161027a565b50505050505050565b6040517f0750181900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005473ffffffffffffffffffffffffffffffffffffffff16331461166f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161027a565b565b6000808a8a8a8a8a8a8a8a8a60405160200161169599989796959493929190611f87565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179150509998505050505050505050565b3373ffffffffffffffffffffffffffffffffffffffff82160361179b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161027a565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b604051806103e00160405280601f906020820280368337509192915050565b6000815180845260005b818110156118565760208185018101518683018201520161183a565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b6020815260006118a76020830184611830565b9392505050565b600081518084526020808501945080840160005b838110156118f457815173ffffffffffffffffffffffffffffffffffffffff16875295820195908201906001016118c2565b509495945050505050565b6020815260006118a760208301846118ae565b60008083601f84011261192457600080fd5b50813567ffffffffffffffff81111561193c57600080fd5b6020830191508360208260051b850101111561195757600080fd5b9250929050565b60008060008060008060008060e0898b03121561197a57600080fd5b606089018a81111561198b57600080fd5b8998503567ffffffffffffffff808211156119a557600080fd5b818b0191508b601f8301126119b957600080fd5b8135818111156119c857600080fd5b8c60208285010111156119da57600080fd5b6020830199508098505060808b01359150808211156119f857600080fd5b611a048c838d01611912565b909750955060a08b0135915080821115611a1d57600080fd5b50611a2a8b828c01611912565b999c989b50969995989497949560c00135949350505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715611ab957611ab9611a43565b604052919050565b803573ffffffffffffffffffffffffffffffffffffffff81168114611ae557600080fd5b919050565b600082601f830112611afb57600080fd5b8135602067ffffffffffffffff821115611b1757611b17611a43565b8160051b611b26828201611a72565b9283528481018201928281019087851115611b4057600080fd5b83870192505b84831015611b6657611b5783611ac1565b82529183019190830190611b46565b979650505050505050565b803560ff81168114611ae557600080fd5b600082601f830112611b9357600080fd5b813567ffffffffffffffff811115611bad57611bad611a43565b611bde60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601611a72565b818152846020838601011115611bf357600080fd5b816020850160208301376000918101602001919091529392505050565b803567ffffffffffffffff81168114611ae557600080fd5b60008060008060008060c08789031215611c4157600080fd5b863567ffffffffffffffff80821115611c5957600080fd5b611c658a838b01611aea565b97506020890135915080821115611c7b57600080fd5b611c878a838b01611aea565b9650611c9560408a01611b71565b95506060890135915080821115611cab57600080fd5b611cb78a838b01611b82565b9450611cc560808a01611c10565b935060a0890135915080821115611cdb57600080fd5b50611ce889828a01611b82565b9150509295509295509295565b600060208284031215611d0757600080fd5b6118a782611ac1565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60ff8181168382160190811115611d8757611d87611d3f565b92915050565b600060ff831680611dc7577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b8060ff84160491505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b8183823760009101908152919050565b828152606082602083013760800192915050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203611e5a57611e5a611d3f565b5060010190565b8082028115828204841417611d8757611d87611d3f565b81810381811115611d8757611d87611d3f565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b63ffffffff818116838216019080821115611ed757611ed7611d3f565b5092915050565b600061012063ffffffff808d1684528b6020850152808b16604085015250806060840152611f0e8184018a6118ae565b90508281036080840152611f2281896118ae565b905060ff871660a084015282810360c0840152611f3f8187611830565b905067ffffffffffffffff851660e0840152828103610100840152611f648185611830565b9c9b505050505050505050505050565b80820180821115611d8757611d87611d3f565b60006101208b835273ffffffffffffffffffffffffffffffffffffffff8b16602084015267ffffffffffffffff808b166040850152816060850152611fce8285018b6118ae565b91508382036080850152611fe2828a6118ae565b915060ff881660a085015283820360c0850152611fff8288611830565b90861660e08501528381036101008501529050611f64818561183056fea164736f6c6343000813000a",
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
