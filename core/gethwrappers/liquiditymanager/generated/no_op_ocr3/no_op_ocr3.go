// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package no_op_ocr3

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

var NoOpOCR3MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"expected\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"actual\",\"type\":\"bytes32\"}],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"ForkedChain\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"InvalidConfig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"latestSequenceNumber\",\"type\":\"uint64\"}],\"name\":\"NonIncreasingSequenceNumber\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NonUniqueSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OracleCannotBeZeroAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SignaturesOutOfRegistration\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedTransmitter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"WrongMessageLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WrongNumberOfSignatures\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTransmitters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setOCR3Config\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561001057600080fd5b5033806000816100675760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b038481169190911790915581161561009757610097816100a3565b5050466080525061014c565b336001600160a01b038216036100fb5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161005e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b608051611bea6200016f60003960008181610c910152610cdd0152611bea6000f3fe608060405234801561001057600080fd5b50600436106100a35760003560e01c806381ff704811610076578063afcb95d71161005b578063afcb95d714610184578063b1dc65a4146101bd578063f2fde38b146101d057600080fd5b806381ff70481461012c5780638da5cb5b1461015c57600080fd5b8063181f5a77146100a8578063666cab8d146100fa5780636a11ee901461010f57806379ba509714610124575b600080fd5b6100e46040518060400160405280600e81526020017f4e6f4f704f43523320312e302e3000000000000000000000000000000000000081525081565b6040516100f19190611518565b60405180910390f35b6101026101e3565b6040516100f19190611584565b61012261011d36600461177c565b610252565b005b610122610a5c565b6004546002546040805163ffffffff808516825264010000000090940490931660208401528201526060016100f1565b60005460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100f1565b600254600454604080516001815260208101939093526801000000000000000090910467ffffffffffffffff16908201526060016100f1565b6101226101cb366004611895565b610b59565b6101226101de36600461197a565b6111bf565b6060600780548060200260200160405190810160405280929190818152602001828054801561024857602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff16815260019091019060200180831161021d575b5050505050905090565b855185518560ff16601f8311156102ca576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f746f6f206d616e79207369676e6572730000000000000000000000000000000060448201526064015b60405180910390fd5b80600003610334576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f66206d75737420626520706f736974697665000000000000000000000000000060448201526064016102c1565b8183146103c2576040517f89a61989000000000000000000000000000000000000000000000000000000008152602060048201526024808201527f6f7261636c6520616464726573736573206f7574206f6620726567697374726160448201527f74696f6e0000000000000000000000000000000000000000000000000000000060648201526084016102c1565b6103cd8160036119c4565b8311610435576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f6661756c74792d6f7261636c65206620746f6f2068696768000000000000000060448201526064016102c1565b61043d6111d3565b60065460005b81811015610531576005600060068381548110610462576104626119e1565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff168352820192909252604001812080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000169055600780546005929190849081106104d2576104d26119e1565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff168352820192909252604001902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000169055600101610443565b50895160005b818110156109045760008c8281518110610553576105536119e1565b602002602001015190506000600281111561057057610570611a10565b73ffffffffffffffffffffffffffffffffffffffff8216600090815260056020526040902054610100900460ff1660028111156105af576105af611a10565b14610616576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265706561746564207369676e6572206164647265737300000000000000000060448201526064016102c1565b73ffffffffffffffffffffffffffffffffffffffff8116610663576040517fd6c62c9b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805180820190915260ff83168152602081016001905273ffffffffffffffffffffffffffffffffffffffff821660009081526005602090815260409091208251815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082168117835592840151919283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000161761010083600281111561071357610713611a10565b021790555090505060008c838151811061072f5761072f6119e1565b602002602001015190506000600281111561074c5761074c611a10565b73ffffffffffffffffffffffffffffffffffffffff8216600090815260056020526040902054610100900460ff16600281111561078b5761078b611a10565b146107f2576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f7265706561746564207472616e736d697474657220616464726573730000000060448201526064016102c1565b73ffffffffffffffffffffffffffffffffffffffff811661083f576040517fd6c62c9b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805180820190915260ff84168152602081016002905273ffffffffffffffffffffffffffffffffffffffff821660009081526005602090815260409091208251815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082168117835592840151919283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000016176101008360028111156108ef576108ef611a10565b02179055509050505050806001019050610537565b508a516109189060069060208e01906113f6565b50895161092c9060079060208d01906113f6565b506003805460ff838116610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909216908c1617179055600480546109b29146913091906000906109849063ffffffff16611a3f565b91906101000a81548163ffffffff021916908363ffffffff160217905563ffffffff168e8e8e8e8e8e611256565b600260000181905550600060048054906101000a900463ffffffff169050436004806101000a81548163ffffffff021916908363ffffffff1602179055507f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e0581600260000154600460009054906101000a900463ffffffff168f8f8f8f8f8f604051610a4699989796959493929190611a62565b60405180910390a1505050505050505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610add576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016102c1565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60045460208901359067ffffffffffffffff68010000000000000000909104811690821611610bdc57600480546040517f6e376b6600000000000000000000000000000000000000000000000000000000815267ffffffffffffffff808516938201939093526801000000000000000090910490911660248201526044016102c1565b600480547fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff166801000000000000000067ffffffffffffffff8416021790556040805160608101825260025480825260035460ff808216602085015261010090910416928201929092528a35918214610c8e5780516040517f93df584c0000000000000000000000000000000000000000000000000000000081526004810191909152602481018390526044016102c1565b467f000000000000000000000000000000000000000000000000000000000000000014610d0f576040517f0f01ce850000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000000060048201524660248201526044016102c1565b6040805183815267ffffffffffffffff851660208201527fe893c2681d327421d89e1cb54fbe64645b4dcea668d6826130b62cf4c6eefea2910160405180910390a16020810151610d61906001611af8565b60ff168714610d9c576040517f71253a2500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b868514610dd5576040517fa75d88af00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3360009081526005602090815260408083208151808301909252805460ff80821684529293919291840191610100909104166002811115610e1857610e18611a10565b6002811115610e2957610e29611a10565b9052509050600281602001516002811115610e4657610e46611a10565b148015610e8d57506007816000015160ff1681548110610e6857610e686119e1565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff1633145b610ec3576040517fda0f08e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b506000610ed18660206119c4565b610edc8960206119c4565b610ee88c610144611b11565b610ef29190611b11565b610efc9190611b11565b9050368114610f40576040517f8e1192e1000000000000000000000000000000000000000000000000000000008152600481018290523660248201526044016102c1565b5060008a8a604051610f53929190611b24565b604051908190038120610f6a918e90602001611b34565b604051602081830303815290604052805190602001209050610f8a611480565b8860005b818110156111ae5760006001858a8460208110610fad57610fad6119e1565b610fba91901a601b611af8565b8f8f86818110610fcc57610fcc6119e1565b905060200201358e8e87818110610fe557610fe56119e1565b9050602002013560405160008152602001604052604051611022949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa158015611044573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff8116600090815260056020908152848220848601909552845460ff80821686529397509195509293928401916101009091041660028111156110c7576110c7611a10565b60028111156110d8576110d8611a10565b90525090506001816020015160028111156110f5576110f5611a10565b1461112c576040517fca31867a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8051859060ff16601f8110611143576111436119e1565b60200201511561117f576040517ff67bc7c400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600185826000015160ff16601f811061119a5761119a6119e1565b911515602090920201525050600101610f8e565b505050505050505050505050505050565b6111c76111d3565b6111d081611301565b50565b60005473ffffffffffffffffffffffffffffffffffffffff163314611254576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016102c1565b565b6000808a8a8a8a8a8a8a8a8a60405160200161127a99989796959493929190611b48565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179150509998505050505050505050565b3373ffffffffffffffffffffffffffffffffffffffff821603611380576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016102c1565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b828054828255906000526020600020908101928215611470579160200282015b8281111561147057825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190611416565b5061147c92915061149f565b5090565b604051806103e00160405280601f906020820280368337509192915050565b5b8082111561147c57600081556001016114a0565b6000815180845260005b818110156114da576020818501810151868301820152016114be565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b60208152600061152b60208301846114b4565b9392505050565b60008151808452602080850194506020840160005b8381101561157957815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101611547565b509495945050505050565b60208152600061152b6020830184611532565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561160d5761160d611597565b604052919050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461163957600080fd5b919050565b600082601f83011261164f57600080fd5b8135602067ffffffffffffffff82111561166b5761166b611597565b8160051b61167a8282016115c6565b928352848101820192828101908785111561169457600080fd5b83870192505b848310156116ba576116ab83611615565b8252918301919083019061169a565b979650505050505050565b803560ff8116811461163957600080fd5b600082601f8301126116e757600080fd5b813567ffffffffffffffff81111561170157611701611597565b61173260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116016115c6565b81815284602083860101111561174757600080fd5b816020850160208301376000918101602001919091529392505050565b803567ffffffffffffffff8116811461163957600080fd5b60008060008060008060c0878903121561179557600080fd5b863567ffffffffffffffff808211156117ad57600080fd5b6117b98a838b0161163e565b975060208901359150808211156117cf57600080fd5b6117db8a838b0161163e565b96506117e960408a016116c5565b955060608901359150808211156117ff57600080fd5b61180b8a838b016116d6565b945061181960808a01611764565b935060a089013591508082111561182f57600080fd5b5061183c89828a016116d6565b9150509295509295509295565b60008083601f84011261185b57600080fd5b50813567ffffffffffffffff81111561187357600080fd5b6020830191508360208260051b850101111561188e57600080fd5b9250929050565b60008060008060008060008060e0898b0312156118b157600080fd5b606089018a8111156118c257600080fd5b8998503567ffffffffffffffff808211156118dc57600080fd5b818b0191508b601f8301126118f057600080fd5b8135818111156118ff57600080fd5b8c602082850101111561191157600080fd5b6020830199508098505060808b013591508082111561192f57600080fd5b61193b8c838d01611849565b909750955060a08b013591508082111561195457600080fd5b506119618b828c01611849565b999c989b50969995989497949560c00135949350505050565b60006020828403121561198c57600080fd5b61152b82611615565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b80820281158282048414176119db576119db611995565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b600063ffffffff808316818103611a5857611a58611995565b6001019392505050565b600061012063ffffffff808d1684528b6020850152808b16604085015250806060840152611a928184018a611532565b90508281036080840152611aa68189611532565b905060ff871660a084015282810360c0840152611ac381876114b4565b905067ffffffffffffffff851660e0840152828103610100840152611ae881856114b4565b9c9b505050505050505050505050565b60ff81811683821601908111156119db576119db611995565b808201808211156119db576119db611995565b8183823760009101908152919050565b828152606082602083013760800192915050565b60006101208b835273ffffffffffffffffffffffffffffffffffffffff8b16602084015267ffffffffffffffff808b166040850152816060850152611b8f8285018b611532565b91508382036080850152611ba3828a611532565b915060ff881660a085015283820360c0850152611bc082886114b4565b90861660e08501528381036101008501529050611ae881856114b456fea164736f6c6343000818000a",
}

var NoOpOCR3ABI = NoOpOCR3MetaData.ABI

var NoOpOCR3Bin = NoOpOCR3MetaData.Bin

func DeployNoOpOCR3(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *NoOpOCR3, error) {
	parsed, err := NoOpOCR3MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(NoOpOCR3Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &NoOpOCR3{address: address, abi: *parsed, NoOpOCR3Caller: NoOpOCR3Caller{contract: contract}, NoOpOCR3Transactor: NoOpOCR3Transactor{contract: contract}, NoOpOCR3Filterer: NoOpOCR3Filterer{contract: contract}}, nil
}

type NoOpOCR3 struct {
	address common.Address
	abi     abi.ABI
	NoOpOCR3Caller
	NoOpOCR3Transactor
	NoOpOCR3Filterer
}

type NoOpOCR3Caller struct {
	contract *bind.BoundContract
}

type NoOpOCR3Transactor struct {
	contract *bind.BoundContract
}

type NoOpOCR3Filterer struct {
	contract *bind.BoundContract
}

type NoOpOCR3Session struct {
	Contract     *NoOpOCR3
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type NoOpOCR3CallerSession struct {
	Contract *NoOpOCR3Caller
	CallOpts bind.CallOpts
}

type NoOpOCR3TransactorSession struct {
	Contract     *NoOpOCR3Transactor
	TransactOpts bind.TransactOpts
}

type NoOpOCR3Raw struct {
	Contract *NoOpOCR3
}

type NoOpOCR3CallerRaw struct {
	Contract *NoOpOCR3Caller
}

type NoOpOCR3TransactorRaw struct {
	Contract *NoOpOCR3Transactor
}

func NewNoOpOCR3(address common.Address, backend bind.ContractBackend) (*NoOpOCR3, error) {
	abi, err := abi.JSON(strings.NewReader(NoOpOCR3ABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindNoOpOCR3(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NoOpOCR3{address: address, abi: abi, NoOpOCR3Caller: NoOpOCR3Caller{contract: contract}, NoOpOCR3Transactor: NoOpOCR3Transactor{contract: contract}, NoOpOCR3Filterer: NoOpOCR3Filterer{contract: contract}}, nil
}

func NewNoOpOCR3Caller(address common.Address, caller bind.ContractCaller) (*NoOpOCR3Caller, error) {
	contract, err := bindNoOpOCR3(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NoOpOCR3Caller{contract: contract}, nil
}

func NewNoOpOCR3Transactor(address common.Address, transactor bind.ContractTransactor) (*NoOpOCR3Transactor, error) {
	contract, err := bindNoOpOCR3(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NoOpOCR3Transactor{contract: contract}, nil
}

func NewNoOpOCR3Filterer(address common.Address, filterer bind.ContractFilterer) (*NoOpOCR3Filterer, error) {
	contract, err := bindNoOpOCR3(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NoOpOCR3Filterer{contract: contract}, nil
}

func bindNoOpOCR3(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := NoOpOCR3MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_NoOpOCR3 *NoOpOCR3Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NoOpOCR3.Contract.NoOpOCR3Caller.contract.Call(opts, result, method, params...)
}

func (_NoOpOCR3 *NoOpOCR3Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NoOpOCR3.Contract.NoOpOCR3Transactor.contract.Transfer(opts)
}

func (_NoOpOCR3 *NoOpOCR3Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NoOpOCR3.Contract.NoOpOCR3Transactor.contract.Transact(opts, method, params...)
}

func (_NoOpOCR3 *NoOpOCR3CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NoOpOCR3.Contract.contract.Call(opts, result, method, params...)
}

func (_NoOpOCR3 *NoOpOCR3TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NoOpOCR3.Contract.contract.Transfer(opts)
}

func (_NoOpOCR3 *NoOpOCR3TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NoOpOCR3.Contract.contract.Transact(opts, method, params...)
}

func (_NoOpOCR3 *NoOpOCR3Caller) GetTransmitters(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _NoOpOCR3.contract.Call(opts, &out, "getTransmitters")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_NoOpOCR3 *NoOpOCR3Session) GetTransmitters() ([]common.Address, error) {
	return _NoOpOCR3.Contract.GetTransmitters(&_NoOpOCR3.CallOpts)
}

func (_NoOpOCR3 *NoOpOCR3CallerSession) GetTransmitters() ([]common.Address, error) {
	return _NoOpOCR3.Contract.GetTransmitters(&_NoOpOCR3.CallOpts)
}

func (_NoOpOCR3 *NoOpOCR3Caller) LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

	error) {
	var out []interface{}
	err := _NoOpOCR3.contract.Call(opts, &out, "latestConfigDetails")

	outstruct := new(LatestConfigDetails)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ConfigDigest = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

func (_NoOpOCR3 *NoOpOCR3Session) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _NoOpOCR3.Contract.LatestConfigDetails(&_NoOpOCR3.CallOpts)
}

func (_NoOpOCR3 *NoOpOCR3CallerSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _NoOpOCR3.Contract.LatestConfigDetails(&_NoOpOCR3.CallOpts)
}

func (_NoOpOCR3 *NoOpOCR3Caller) LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

	error) {
	var out []interface{}
	err := _NoOpOCR3.contract.Call(opts, &out, "latestConfigDigestAndEpoch")

	outstruct := new(LatestConfigDigestAndEpoch)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ScanLogs = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfigDigest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.SequenceNumber = *abi.ConvertType(out[2], new(uint64)).(*uint64)

	return *outstruct, err

}

func (_NoOpOCR3 *NoOpOCR3Session) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _NoOpOCR3.Contract.LatestConfigDigestAndEpoch(&_NoOpOCR3.CallOpts)
}

func (_NoOpOCR3 *NoOpOCR3CallerSession) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _NoOpOCR3.Contract.LatestConfigDigestAndEpoch(&_NoOpOCR3.CallOpts)
}

func (_NoOpOCR3 *NoOpOCR3Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NoOpOCR3.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_NoOpOCR3 *NoOpOCR3Session) Owner() (common.Address, error) {
	return _NoOpOCR3.Contract.Owner(&_NoOpOCR3.CallOpts)
}

func (_NoOpOCR3 *NoOpOCR3CallerSession) Owner() (common.Address, error) {
	return _NoOpOCR3.Contract.Owner(&_NoOpOCR3.CallOpts)
}

func (_NoOpOCR3 *NoOpOCR3Caller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _NoOpOCR3.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_NoOpOCR3 *NoOpOCR3Session) TypeAndVersion() (string, error) {
	return _NoOpOCR3.Contract.TypeAndVersion(&_NoOpOCR3.CallOpts)
}

func (_NoOpOCR3 *NoOpOCR3CallerSession) TypeAndVersion() (string, error) {
	return _NoOpOCR3.Contract.TypeAndVersion(&_NoOpOCR3.CallOpts)
}

func (_NoOpOCR3 *NoOpOCR3Transactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NoOpOCR3.contract.Transact(opts, "acceptOwnership")
}

func (_NoOpOCR3 *NoOpOCR3Session) AcceptOwnership() (*types.Transaction, error) {
	return _NoOpOCR3.Contract.AcceptOwnership(&_NoOpOCR3.TransactOpts)
}

func (_NoOpOCR3 *NoOpOCR3TransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _NoOpOCR3.Contract.AcceptOwnership(&_NoOpOCR3.TransactOpts)
}

func (_NoOpOCR3 *NoOpOCR3Transactor) SetOCR3Config(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _NoOpOCR3.contract.Transact(opts, "setOCR3Config", signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_NoOpOCR3 *NoOpOCR3Session) SetOCR3Config(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _NoOpOCR3.Contract.SetOCR3Config(&_NoOpOCR3.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_NoOpOCR3 *NoOpOCR3TransactorSession) SetOCR3Config(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _NoOpOCR3.Contract.SetOCR3Config(&_NoOpOCR3.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_NoOpOCR3 *NoOpOCR3Transactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _NoOpOCR3.contract.Transact(opts, "transferOwnership", to)
}

func (_NoOpOCR3 *NoOpOCR3Session) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _NoOpOCR3.Contract.TransferOwnership(&_NoOpOCR3.TransactOpts, to)
}

func (_NoOpOCR3 *NoOpOCR3TransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _NoOpOCR3.Contract.TransferOwnership(&_NoOpOCR3.TransactOpts, to)
}

func (_NoOpOCR3 *NoOpOCR3Transactor) Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _NoOpOCR3.contract.Transact(opts, "transmit", reportContext, report, rs, ss, rawVs)
}

func (_NoOpOCR3 *NoOpOCR3Session) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _NoOpOCR3.Contract.Transmit(&_NoOpOCR3.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_NoOpOCR3 *NoOpOCR3TransactorSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _NoOpOCR3.Contract.Transmit(&_NoOpOCR3.TransactOpts, reportContext, report, rs, ss, rawVs)
}

type NoOpOCR3ConfigSetIterator struct {
	Event *NoOpOCR3ConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NoOpOCR3ConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NoOpOCR3ConfigSet)
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
		it.Event = new(NoOpOCR3ConfigSet)
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

func (it *NoOpOCR3ConfigSetIterator) Error() error {
	return it.fail
}

func (it *NoOpOCR3ConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NoOpOCR3ConfigSet struct {
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

func (_NoOpOCR3 *NoOpOCR3Filterer) FilterConfigSet(opts *bind.FilterOpts) (*NoOpOCR3ConfigSetIterator, error) {

	logs, sub, err := _NoOpOCR3.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &NoOpOCR3ConfigSetIterator{contract: _NoOpOCR3.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_NoOpOCR3 *NoOpOCR3Filterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *NoOpOCR3ConfigSet) (event.Subscription, error) {

	logs, sub, err := _NoOpOCR3.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NoOpOCR3ConfigSet)
				if err := _NoOpOCR3.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_NoOpOCR3 *NoOpOCR3Filterer) ParseConfigSet(log types.Log) (*NoOpOCR3ConfigSet, error) {
	event := new(NoOpOCR3ConfigSet)
	if err := _NoOpOCR3.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type NoOpOCR3OwnershipTransferRequestedIterator struct {
	Event *NoOpOCR3OwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NoOpOCR3OwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NoOpOCR3OwnershipTransferRequested)
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
		it.Event = new(NoOpOCR3OwnershipTransferRequested)
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

func (it *NoOpOCR3OwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *NoOpOCR3OwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NoOpOCR3OwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_NoOpOCR3 *NoOpOCR3Filterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*NoOpOCR3OwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _NoOpOCR3.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &NoOpOCR3OwnershipTransferRequestedIterator{contract: _NoOpOCR3.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_NoOpOCR3 *NoOpOCR3Filterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *NoOpOCR3OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _NoOpOCR3.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NoOpOCR3OwnershipTransferRequested)
				if err := _NoOpOCR3.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_NoOpOCR3 *NoOpOCR3Filterer) ParseOwnershipTransferRequested(log types.Log) (*NoOpOCR3OwnershipTransferRequested, error) {
	event := new(NoOpOCR3OwnershipTransferRequested)
	if err := _NoOpOCR3.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type NoOpOCR3OwnershipTransferredIterator struct {
	Event *NoOpOCR3OwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NoOpOCR3OwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NoOpOCR3OwnershipTransferred)
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
		it.Event = new(NoOpOCR3OwnershipTransferred)
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

func (it *NoOpOCR3OwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *NoOpOCR3OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NoOpOCR3OwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_NoOpOCR3 *NoOpOCR3Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*NoOpOCR3OwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _NoOpOCR3.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &NoOpOCR3OwnershipTransferredIterator{contract: _NoOpOCR3.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_NoOpOCR3 *NoOpOCR3Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *NoOpOCR3OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _NoOpOCR3.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NoOpOCR3OwnershipTransferred)
				if err := _NoOpOCR3.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_NoOpOCR3 *NoOpOCR3Filterer) ParseOwnershipTransferred(log types.Log) (*NoOpOCR3OwnershipTransferred, error) {
	event := new(NoOpOCR3OwnershipTransferred)
	if err := _NoOpOCR3.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type NoOpOCR3TransmittedIterator struct {
	Event *NoOpOCR3Transmitted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NoOpOCR3TransmittedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NoOpOCR3Transmitted)
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
		it.Event = new(NoOpOCR3Transmitted)
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

func (it *NoOpOCR3TransmittedIterator) Error() error {
	return it.fail
}

func (it *NoOpOCR3TransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NoOpOCR3Transmitted struct {
	ConfigDigest   [32]byte
	SequenceNumber uint64
	Raw            types.Log
}

func (_NoOpOCR3 *NoOpOCR3Filterer) FilterTransmitted(opts *bind.FilterOpts) (*NoOpOCR3TransmittedIterator, error) {

	logs, sub, err := _NoOpOCR3.contract.FilterLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return &NoOpOCR3TransmittedIterator{contract: _NoOpOCR3.contract, event: "Transmitted", logs: logs, sub: sub}, nil
}

func (_NoOpOCR3 *NoOpOCR3Filterer) WatchTransmitted(opts *bind.WatchOpts, sink chan<- *NoOpOCR3Transmitted) (event.Subscription, error) {

	logs, sub, err := _NoOpOCR3.contract.WatchLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NoOpOCR3Transmitted)
				if err := _NoOpOCR3.contract.UnpackLog(event, "Transmitted", log); err != nil {
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

func (_NoOpOCR3 *NoOpOCR3Filterer) ParseTransmitted(log types.Log) (*NoOpOCR3Transmitted, error) {
	event := new(NoOpOCR3Transmitted)
	if err := _NoOpOCR3.contract.UnpackLog(event, "Transmitted", log); err != nil {
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
	ScanLogs       bool
	ConfigDigest   [32]byte
	SequenceNumber uint64
}

func (_NoOpOCR3 *NoOpOCR3) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _NoOpOCR3.abi.Events["ConfigSet"].ID:
		return _NoOpOCR3.ParseConfigSet(log)
	case _NoOpOCR3.abi.Events["OwnershipTransferRequested"].ID:
		return _NoOpOCR3.ParseOwnershipTransferRequested(log)
	case _NoOpOCR3.abi.Events["OwnershipTransferred"].ID:
		return _NoOpOCR3.ParseOwnershipTransferred(log)
	case _NoOpOCR3.abi.Events["Transmitted"].ID:
		return _NoOpOCR3.ParseTransmitted(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (NoOpOCR3ConfigSet) Topic() common.Hash {
	return common.HexToHash("0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05")
}

func (NoOpOCR3OwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (NoOpOCR3OwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (NoOpOCR3Transmitted) Topic() common.Hash {
	return common.HexToHash("0xe893c2681d327421d89e1cb54fbe64645b4dcea668d6826130b62cf4c6eefea2")
}

func (_NoOpOCR3 *NoOpOCR3) Address() common.Address {
	return _NoOpOCR3.address
}

type NoOpOCR3Interface interface {
	GetTransmitters(opts *bind.CallOpts) ([]common.Address, error)

	LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

		error)

	LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

		error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	SetOCR3Config(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error)

	FilterConfigSet(opts *bind.FilterOpts) (*NoOpOCR3ConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *NoOpOCR3ConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*NoOpOCR3ConfigSet, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*NoOpOCR3OwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *NoOpOCR3OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*NoOpOCR3OwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*NoOpOCR3OwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *NoOpOCR3OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*NoOpOCR3OwnershipTransferred, error)

	FilterTransmitted(opts *bind.FilterOpts) (*NoOpOCR3TransmittedIterator, error)

	WatchTransmitted(opts *bind.WatchOpts, sink chan<- *NoOpOCR3Transmitted) (event.Subscription, error)

	ParseTransmitted(log types.Log) (*NoOpOCR3Transmitted, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
