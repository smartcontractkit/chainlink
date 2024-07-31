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
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"expected\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"actual\",\"type\":\"bytes32\"}],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"ForkedChain\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"InvalidConfig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"latestSequenceNumber\",\"type\":\"uint64\"}],\"name\":\"NonIncreasingSequenceNumber\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NonUniqueSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OracleCannotBeZeroAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SignaturesOutOfRegistration\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedTransmitter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"WrongMessageLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WrongNumberOfSignatures\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTransmitters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestSequenceNumber\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setOCR3Config\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561001057600080fd5b5033806000816100675760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b038481169190911790915581161561009757610097816100a3565b5050466080525061014c565b336001600160a01b038216036100fb5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161005e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b608051611c0d6200016f60003960008181610cb40152610d000152611c0d6000f3fe608060405234801561001057600080fd5b50600436106100a35760003560e01c806381ff704811610076578063b1dc65a41161005b578063b1dc65a414610184578063f1c0461614610197578063f2fde38b146101c957600080fd5b806381ff70481461012c5780638da5cb5b1461015c57600080fd5b8063181f5a77146100a8578063666cab8d146100fa5780636a11ee901461010f57806379ba509714610124575b600080fd5b6100e46040518060400160405280600e81526020017f4e6f4f704f43523320312e302e3000000000000000000000000000000000000081525081565b6040516100f1919061153b565b60405180910390f35b6101026101dc565b6040516100f191906115a7565b61012261011d36600461179f565b61024b565b005b610122610a7f565b6004546002546040805163ffffffff808516825264010000000090940490931660208401528201526060016100f1565b60005460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100f1565b6101226101923660046118b8565b610b7c565b60045468010000000000000000900467ffffffffffffffff1660405167ffffffffffffffff90911681526020016100f1565b6101226101d736600461199d565b6111e2565b6060600780548060200260200160405190810160405280929190818152602001828054801561024157602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610216575b5050505050905090565b855185518560ff16601f8311156102c3576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f746f6f206d616e79207369676e6572730000000000000000000000000000000060448201526064015b60405180910390fd5b8060000361032d576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f66206d75737420626520706f736974697665000000000000000000000000000060448201526064016102ba565b8183146103bb576040517f89a61989000000000000000000000000000000000000000000000000000000008152602060048201526024808201527f6f7261636c6520616464726573736573206f7574206f6620726567697374726160448201527f74696f6e0000000000000000000000000000000000000000000000000000000060648201526084016102ba565b6103c68160036119e7565b831161042e576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f6661756c74792d6f7261636c65206620746f6f2068696768000000000000000060448201526064016102ba565b6104366111f6565b60065460005b8181101561052a57600560006006838154811061045b5761045b611a04565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff168352820192909252604001812080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000169055600780546005929190849081106104cb576104cb611a04565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff168352820192909252604001902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000016905560010161043c565b50895160005b818110156108fd5760008c828151811061054c5761054c611a04565b602002602001015190506000600281111561056957610569611a33565b73ffffffffffffffffffffffffffffffffffffffff8216600090815260056020526040902054610100900460ff1660028111156105a8576105a8611a33565b1461060f576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265706561746564207369676e6572206164647265737300000000000000000060448201526064016102ba565b73ffffffffffffffffffffffffffffffffffffffff811661065c576040517fd6c62c9b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805180820190915260ff83168152602081016001905273ffffffffffffffffffffffffffffffffffffffff821660009081526005602090815260409091208251815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082168117835592840151919283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000161761010083600281111561070c5761070c611a33565b021790555090505060008c838151811061072857610728611a04565b602002602001015190506000600281111561074557610745611a33565b73ffffffffffffffffffffffffffffffffffffffff8216600090815260056020526040902054610100900460ff16600281111561078457610784611a33565b146107eb576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f7265706561746564207472616e736d697474657220616464726573730000000060448201526064016102ba565b73ffffffffffffffffffffffffffffffffffffffff8116610838576040517fd6c62c9b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805180820190915260ff84168152602081016002905273ffffffffffffffffffffffffffffffffffffffff821660009081526005602090815260409091208251815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082168117835592840151919283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000016176101008360028111156108e8576108e8611a33565b02179055509050505050806001019050610530565b508a516109119060069060208e0190611419565b5089516109259060079060208d0190611419565b506003805460ff838116610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909216908c1617179055600480546109ab91469130919060009061097d9063ffffffff16611a62565b91906101000a81548163ffffffff021916908363ffffffff160217905563ffffffff168e8e8e8e8e8e611279565b600260000181905550600060048054906101000a900463ffffffff169050436004806101000a81548163ffffffff021916908363ffffffff1602179055506000600460086101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055507f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e0581600260000154600460009054906101000a900463ffffffff168f8f8f8f8f8f604051610a6999989796959493929190611a85565b60405180910390a1505050505050505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610b00576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016102ba565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60045460208901359067ffffffffffffffff68010000000000000000909104811690821611610bff57600480546040517f6e376b6600000000000000000000000000000000000000000000000000000000815267ffffffffffffffff808516938201939093526801000000000000000090910490911660248201526044016102ba565b600480547fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff166801000000000000000067ffffffffffffffff8416021790556040805160608101825260025480825260035460ff808216602085015261010090910416928201929092528a35918214610cb15780516040517f93df584c0000000000000000000000000000000000000000000000000000000081526004810191909152602481018390526044016102ba565b467f000000000000000000000000000000000000000000000000000000000000000014610d32576040517f0f01ce850000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000000060048201524660248201526044016102ba565b6040805183815267ffffffffffffffff851660208201527fe893c2681d327421d89e1cb54fbe64645b4dcea668d6826130b62cf4c6eefea2910160405180910390a16020810151610d84906001611b1b565b60ff168714610dbf576040517f71253a2500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b868514610df8576040517fa75d88af00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3360009081526005602090815260408083208151808301909252805460ff80821684529293919291840191610100909104166002811115610e3b57610e3b611a33565b6002811115610e4c57610e4c611a33565b9052509050600281602001516002811115610e6957610e69611a33565b148015610eb057506007816000015160ff1681548110610e8b57610e8b611a04565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff1633145b610ee6576040517fda0f08e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b506000610ef48660206119e7565b610eff8960206119e7565b610f0b8c610144611b34565b610f159190611b34565b610f1f9190611b34565b9050368114610f63576040517f8e1192e1000000000000000000000000000000000000000000000000000000008152600481018290523660248201526044016102ba565b5060008a8a604051610f76929190611b47565b604051908190038120610f8d918e90602001611b57565b604051602081830303815290604052805190602001209050610fad6114a3565b8860005b818110156111d15760006001858a8460208110610fd057610fd0611a04565b610fdd91901a601b611b1b565b8f8f86818110610fef57610fef611a04565b905060200201358e8e8781811061100857611008611a04565b9050602002013560405160008152602001604052604051611045949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa158015611067573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff8116600090815260056020908152848220848601909552845460ff80821686529397509195509293928401916101009091041660028111156110ea576110ea611a33565b60028111156110fb576110fb611a33565b905250905060018160200151600281111561111857611118611a33565b1461114f576040517fca31867a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8051859060ff16601f811061116657611166611a04565b6020020151156111a2576040517ff67bc7c400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600185826000015160ff16601f81106111bd576111bd611a04565b911515602090920201525050600101610fb1565b505050505050505050505050505050565b6111ea6111f6565b6111f381611324565b50565b60005473ffffffffffffffffffffffffffffffffffffffff163314611277576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016102ba565b565b6000808a8a8a8a8a8a8a8a8a60405160200161129d99989796959493929190611b6b565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179150509998505050505050505050565b3373ffffffffffffffffffffffffffffffffffffffff8216036113a3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016102ba565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b828054828255906000526020600020908101928215611493579160200282015b8281111561149357825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190611439565b5061149f9291506114c2565b5090565b604051806103e00160405280601f906020820280368337509192915050565b5b8082111561149f57600081556001016114c3565b6000815180845260005b818110156114fd576020818501810151868301820152016114e1565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b60208152600061154e60208301846114d7565b9392505050565b60008151808452602080850194506020840160005b8381101561159c57815173ffffffffffffffffffffffffffffffffffffffff168752958201959082019060010161156a565b509495945050505050565b60208152600061154e6020830184611555565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715611630576116306115ba565b604052919050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461165c57600080fd5b919050565b600082601f83011261167257600080fd5b8135602067ffffffffffffffff82111561168e5761168e6115ba565b8160051b61169d8282016115e9565b92835284810182019282810190878511156116b757600080fd5b83870192505b848310156116dd576116ce83611638565b825291830191908301906116bd565b979650505050505050565b803560ff8116811461165c57600080fd5b600082601f83011261170a57600080fd5b813567ffffffffffffffff811115611724576117246115ba565b61175560207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116016115e9565b81815284602083860101111561176a57600080fd5b816020850160208301376000918101602001919091529392505050565b803567ffffffffffffffff8116811461165c57600080fd5b60008060008060008060c087890312156117b857600080fd5b863567ffffffffffffffff808211156117d057600080fd5b6117dc8a838b01611661565b975060208901359150808211156117f257600080fd5b6117fe8a838b01611661565b965061180c60408a016116e8565b9550606089013591508082111561182257600080fd5b61182e8a838b016116f9565b945061183c60808a01611787565b935060a089013591508082111561185257600080fd5b5061185f89828a016116f9565b9150509295509295509295565b60008083601f84011261187e57600080fd5b50813567ffffffffffffffff81111561189657600080fd5b6020830191508360208260051b85010111156118b157600080fd5b9250929050565b60008060008060008060008060e0898b0312156118d457600080fd5b606089018a8111156118e557600080fd5b8998503567ffffffffffffffff808211156118ff57600080fd5b818b0191508b601f83011261191357600080fd5b81358181111561192257600080fd5b8c602082850101111561193457600080fd5b6020830199508098505060808b013591508082111561195257600080fd5b61195e8c838d0161186c565b909750955060a08b013591508082111561197757600080fd5b506119848b828c0161186c565b999c989b50969995989497949560c00135949350505050565b6000602082840312156119af57600080fd5b61154e82611638565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b80820281158282048414176119fe576119fe6119b8565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b600063ffffffff808316818103611a7b57611a7b6119b8565b6001019392505050565b600061012063ffffffff808d1684528b6020850152808b16604085015250806060840152611ab58184018a611555565b90508281036080840152611ac98189611555565b905060ff871660a084015282810360c0840152611ae681876114d7565b905067ffffffffffffffff851660e0840152828103610100840152611b0b81856114d7565b9c9b505050505050505050505050565b60ff81811683821601908111156119fe576119fe6119b8565b808201808211156119fe576119fe6119b8565b8183823760009101908152919050565b828152606082602083013760800192915050565b60006101208b835273ffffffffffffffffffffffffffffffffffffffff8b16602084015267ffffffffffffffff808b166040850152816060850152611bb28285018b611555565b91508382036080850152611bc6828a611555565b915060ff881660a085015283820360c0850152611be382886114d7565b90861660e08501528381036101008501529050611b0b81856114d756fea164736f6c6343000818000a",
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

func (_NoOpOCR3 *NoOpOCR3Caller) LatestSequenceNumber(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _NoOpOCR3.contract.Call(opts, &out, "latestSequenceNumber")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_NoOpOCR3 *NoOpOCR3Session) LatestSequenceNumber() (uint64, error) {
	return _NoOpOCR3.Contract.LatestSequenceNumber(&_NoOpOCR3.CallOpts)
}

func (_NoOpOCR3 *NoOpOCR3CallerSession) LatestSequenceNumber() (uint64, error) {
	return _NoOpOCR3.Contract.LatestSequenceNumber(&_NoOpOCR3.CallOpts)
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

	LatestSequenceNumber(opts *bind.CallOpts) (uint64, error)

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
