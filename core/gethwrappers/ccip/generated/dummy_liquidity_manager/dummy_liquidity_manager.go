// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package dummy_liquidity_manager

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

type DummyLiquidityManagerCrossChainLiquidityManager struct {
	RemoteLiquidityManager common.Address
	Enabled                bool
}

type ILiquidityManagerCrossChainLiquidityManagerArgs struct {
	RemoteLiquidityManager common.Address
	RemoteChainSelector    uint64
	Enabled                bool
}

var DummyLiquidityManagerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"localChainSelector\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"expected\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"actual\",\"type\":\"bytes32\"}],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"ForkedChain\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"}],\"name\":\"InsufficientLiquidity\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"InvalidConfig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"}],\"name\":\"InvalidRemoteChain\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"latestSequenceNumber\",\"type\":\"uint64\"}],\"name\":\"NonIncreasingSequenceNumber\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NonUniqueSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OracleCannotBeZeroAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SignaturesOutOfRegistration\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedTransmitter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"WrongMessageLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WrongNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroChainSelector\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"LiquidityAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"remover\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"LiquidityRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"ocrSeqNum\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"fromChainSelector\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"toChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"LiquidityTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllCrossChainLiquidityMangers\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"remoteLiquidityManager\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"internalType\":\"structILiquidityManager.CrossChainLiquidityManagerArgs[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"}],\"name\":\"getCrossChainLiquidityManager\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"remoteLiquidityManager\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"internalType\":\"structDummyLiquidityManager.CrossChainLiquidityManager\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSupportedDestChains\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTransmitters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"remoteLiquidityManager\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"internalType\":\"structILiquidityManager.CrossChainLiquidityManagerArgs\",\"name\":\"crossChainLiqManager\",\"type\":\"tuple\"}],\"name\":\"setCrossChainLiquidityManager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"remoteLiquidityManager\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"internalType\":\"structILiquidityManager.CrossChainLiquidityManagerArgs[]\",\"name\":\"crossChainLiquidityManagers\",\"type\":\"tuple[]\"}],\"name\":\"setCrossChainLiquidityManager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setOCR3Config\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60c06040523480156200001157600080fd5b50604051620025af380380620025af8339810160408190526200003491620001ad565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be8162000102565b505046608052506001600160401b038116600003620000f05760405163f89d762960e01b815260040160405180910390fd5b6001600160401b031660a052620001df565b336001600160a01b038216036200015c5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600060208284031215620001c057600080fd5b81516001600160401b0381168114620001d857600080fd5b9392505050565b60805160a0516123a66200020960003960005050600081816110af01526110fb01526123a66000f3fe608060405234801561001057600080fd5b50600436106100ea5760003560e01c806381ff70481161008c578063b1dc65a411610066578063b1dc65a4146102e4578063def51dde146102f7578063f2fde38b1461030c578063f4bb82da1461031f57600080fd5b806381ff7048146102535780638da5cb5b14610283578063afcb95d7146102ab57600080fd5b8063348759c1116100c8578063348759c11461020e578063666cab8d146102235780636a11ee901461023857806379ba50971461024b57600080fd5b8063181f5a77146100ef57806323f641e81461014157806328f18e4b146101f9575b600080fd5b61012b6040518060400160405280601b81526020017f44756d6d794c69717569646974794d616e6167657220312e302e30000000000081525081565b6040516101389190611b10565b60405180910390f35b6101c661014f366004611b47565b60408051808201909152600080825260208201525067ffffffffffffffff1660009081526008602090815260409182902082518084019093525473ffffffffffffffffffffffffffffffffffffffff8116835274010000000000000000000000000000000000000000900460ff1615159082015290565b60408051825173ffffffffffffffffffffffffffffffffffffffff16815260209283015115159281019290925201610138565b61020c610207366004611b62565b610332565b005b610216610568565b6040516101389190611b74565b61022b6105f4565b6040516101389190611c13565b61020c610246366004611dee565b610662565b61020c610e7a565b6004546002546040805163ffffffff80851682526401000000009094049093166020840152820152606001610138565b60005460405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610138565b600254600454604080516001815260208101939093526801000000000000000090910467ffffffffffffffff1690820152606001610138565b61020c6102f2366004611f07565b610f77565b6102ff6115e6565b6040516101389190611fec565b61020c61031a366004612067565b61176f565b61020c61032d366004612082565b611783565b61033a6117cb565b61034a6040820160208301611b47565b67ffffffffffffffff1660000361038d576040517ff89d762900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600061039c6020830183612067565b73ffffffffffffffffffffffffffffffffffffffff16036103e9576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006008816103fe6040850160208601611b47565b67ffffffffffffffff16815260208101919091526040016000205473ffffffffffffffffffffffffffffffffffffffff16036104875760096104466040830160208401611b47565b8154600181018355600092835260209092206004830401805460039093166008026101000a67ffffffffffffffff8181021990941692909316929092021790555b604080518082019091528061049f6020840184612067565b73ffffffffffffffffffffffffffffffffffffffff1681526020016104ca60608401604085016120f7565b15159052600860006104e26040850160208601611b47565b67ffffffffffffffff168152602080820192909252604001600020825181549390920151151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00000000000000000000000000000000000000000090931673ffffffffffffffffffffffffffffffffffffffff9092169190911791909117905550565b606060098054806020026020016040519081016040528092919081815260200182805480156105ea57602002820191906000526020600020906000905b82829054906101000a900467ffffffffffffffff1667ffffffffffffffff16815260200190600801906020826007010492830192600103820291508084116105a55790505b5050505050905090565b606060078054806020026020016040519081016040528092919081815260200182805480156105ea57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff16815260019091019060200180831161062e575050505050905090565b855185518560ff16601f8311156106da576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f746f6f206d616e79207369676e6572730000000000000000000000000000000060448201526064015b60405180910390fd5b80600003610744576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f66206d75737420626520706f736974697665000000000000000000000000000060448201526064016106d1565b8183146107d2576040517f89a61989000000000000000000000000000000000000000000000000000000008152602060048201526024808201527f6f7261636c6520616464726573736573206f7574206f6620726567697374726160448201527f74696f6e0000000000000000000000000000000000000000000000000000000060648201526084016106d1565b6107dd816003612148565b8311610845576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f6661756c74792d6f7261636c65206620746f6f2068696768000000000000000060448201526064016106d1565b61084d6117cb565b60065460005b8181101561094957600560006006838154811061087257610872612165565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff168352820192909252604001812080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000169055600780546005929190849081106108e2576108e2612165565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff168352820192909252604001902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000016905561094281612194565b9050610853565b50895160005b81811015610d225760008c828151811061096b5761096b612165565b6020026020010151905060006002811115610988576109886121cc565b73ffffffffffffffffffffffffffffffffffffffff8216600090815260056020526040902054610100900460ff1660028111156109c7576109c76121cc565b14610a2e576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265706561746564207369676e6572206164647265737300000000000000000060448201526064016106d1565b73ffffffffffffffffffffffffffffffffffffffff8116610a7b576040517fd6c62c9b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805180820190915260ff83168152602081016001905273ffffffffffffffffffffffffffffffffffffffff821660009081526005602090815260409091208251815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082168117835592840151919283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001617610100836002811115610b2b57610b2b6121cc565b021790555090505060008c8381518110610b4757610b47612165565b6020026020010151905060006002811115610b6457610b646121cc565b73ffffffffffffffffffffffffffffffffffffffff8216600090815260056020526040902054610100900460ff166002811115610ba357610ba36121cc565b14610c0a576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f7265706561746564207472616e736d697474657220616464726573730000000060448201526064016106d1565b73ffffffffffffffffffffffffffffffffffffffff8116610c57576040517fd6c62c9b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805180820190915260ff84168152602081016002905273ffffffffffffffffffffffffffffffffffffffff821660009081526005602090815260409091208251815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082168117835592840151919283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001617610100836002811115610d0757610d076121cc565b0217905550905050505080610d1b90612194565b905061094f565b508a51610d369060069060208e01906119ee565b508951610d4a9060079060208d01906119ee565b506003805460ff838116610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909216908c161717905560048054610dd0914691309190600090610da29063ffffffff166121fb565b91906101000a81548163ffffffff021916908363ffffffff160217905563ffffffff168e8e8e8e8e8e61184e565b600260000181905550600060048054906101000a900463ffffffff169050436004806101000a81548163ffffffff021916908363ffffffff1602179055507f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e0581600260000154600460009054906101000a900463ffffffff168f8f8f8f8f8f604051610e649998979695949392919061221e565b60405180910390a1505050505050505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610efb576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016106d1565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60045460208901359067ffffffffffffffff68010000000000000000909104811690821611610ffa57600480546040517f6e376b6600000000000000000000000000000000000000000000000000000000815267ffffffffffffffff808516938201939093526801000000000000000090910490911660248201526044016106d1565b600480547fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff166801000000000000000067ffffffffffffffff8416021790556040805160608101825260025480825260035460ff808216602085015261010090910416928201929092528a359182146110ac5780516040517f93df584c0000000000000000000000000000000000000000000000000000000081526004810191909152602481018390526044016106d1565b467f00000000000000000000000000000000000000000000000000000000000000001461112d576040517f0f01ce850000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000000060048201524660248201526044016106d1565b6040805183815267ffffffffffffffff851660208201527fe893c2681d327421d89e1cb54fbe64645b4dcea668d6826130b62cf4c6eefea2910160405180910390a1602081015161117f9060016122b4565b60ff1687146111ba576040517f71253a2500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8685146111f3576040517fa75d88af00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3360009081526005602090815260408083208151808301909252805460ff80821684529293919291840191610100909104166002811115611236576112366121cc565b6002811115611247576112476121cc565b9052509050600281602001516002811115611264576112646121cc565b1480156112ab57506007816000015160ff168154811061128657611286612165565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff1633145b6112e1576040517fda0f08e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5060006112ef866020612148565b6112fa896020612148565b6113068c6101446122cd565b61131091906122cd565b61131a91906122cd565b905036811461135e576040517f8e1192e1000000000000000000000000000000000000000000000000000000008152600481018290523660248201526044016106d1565b5060008a8a6040516113719291906122e0565b604051908190038120611388918e906020016122f0565b6040516020818303038152906040528051906020012090506113a8611a78565b8860005b818110156115d55760006001858a84602081106113cb576113cb612165565b6113d891901a601b6122b4565b8f8f868181106113ea576113ea612165565b905060200201358e8e8781811061140357611403612165565b9050602002013560405160008152602001604052604051611440949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa158015611462573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff8116600090815260056020908152848220848601909552845460ff80821686529397509195509293928401916101009091041660028111156114e5576114e56121cc565b60028111156114f6576114f66121cc565b9052509050600181602001516002811115611513576115136121cc565b1461154a576040517fca31867a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8051859060ff16601f811061156157611561612165565b60200201511561159d576040517ff67bc7c400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600185826000015160ff16601f81106115b8576115b8612165565b91151560209092020152506115ce905081612194565b90506113ac565b505050505050505050505050505050565b60095460609060009067ffffffffffffffff81111561160757611607611c26565b60405190808252806020026020018201604052801561167057816020015b60408051606081018252600080825260208083018290529282015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9092019101816116255790505b50905060005b6009548110156117695760006009828154811061169557611695612165565b60009182526020808320600483040154600390921660089081026101000a90920467ffffffffffffffff1680845291815260409283902083518085018552905473ffffffffffffffffffffffffffffffffffffffff80821683527401000000000000000000000000000000000000000090910460ff161515828401908152855160608101875283519092168252928101849052915115159382019390935285519193509085908590811061174b5761174b612165565b602002602001018190525050508061176290612194565b9050611676565b50919050565b6117776117cb565b611780816118f9565b50565b61178b6117cb565b60005b818110156117c6576117b68383838181106117ab576117ab612165565b905060600201610332565b6117bf81612194565b905061178e565b505050565b60005473ffffffffffffffffffffffffffffffffffffffff16331461184c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016106d1565b565b6000808a8a8a8a8a8a8a8a8a60405160200161187299989796959493929190612304565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179150509998505050505050505050565b3373ffffffffffffffffffffffffffffffffffffffff821603611978576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016106d1565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b828054828255906000526020600020908101928215611a68579160200282015b82811115611a6857825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190611a0e565b50611a74929150611a97565b5090565b604051806103e00160405280601f906020820280368337509192915050565b5b80821115611a745760008155600101611a98565b6000815180845260005b81811015611ad257602081850181015186830182015201611ab6565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b602081526000611b236020830184611aac565b9392505050565b803567ffffffffffffffff81168114611b4257600080fd5b919050565b600060208284031215611b5957600080fd5b611b2382611b2a565b60006060828403121561176957600080fd5b6020808252825182820181905260009190848201906040850190845b81811015611bb657835167ffffffffffffffff1683529284019291840191600101611b90565b50909695505050505050565b600081518084526020808501945080840160005b83811015611c0857815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101611bd6565b509495945050505050565b602081526000611b236020830184611bc2565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715611c9c57611c9c611c26565b604052919050565b803573ffffffffffffffffffffffffffffffffffffffff81168114611b4257600080fd5b600082601f830112611cd957600080fd5b8135602067ffffffffffffffff821115611cf557611cf5611c26565b8160051b611d04828201611c55565b9283528481018201928281019087851115611d1e57600080fd5b83870192505b84831015611d4457611d3583611ca4565b82529183019190830190611d24565b979650505050505050565b803560ff81168114611b4257600080fd5b600082601f830112611d7157600080fd5b813567ffffffffffffffff811115611d8b57611d8b611c26565b611dbc60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601611c55565b818152846020838601011115611dd157600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008060008060c08789031215611e0757600080fd5b863567ffffffffffffffff80821115611e1f57600080fd5b611e2b8a838b01611cc8565b97506020890135915080821115611e4157600080fd5b611e4d8a838b01611cc8565b9650611e5b60408a01611d4f565b95506060890135915080821115611e7157600080fd5b611e7d8a838b01611d60565b9450611e8b60808a01611b2a565b935060a0890135915080821115611ea157600080fd5b50611eae89828a01611d60565b9150509295509295509295565b60008083601f840112611ecd57600080fd5b50813567ffffffffffffffff811115611ee557600080fd5b6020830191508360208260051b8501011115611f0057600080fd5b9250929050565b60008060008060008060008060e0898b031215611f2357600080fd5b606089018a811115611f3457600080fd5b8998503567ffffffffffffffff80821115611f4e57600080fd5b818b0191508b601f830112611f6257600080fd5b813581811115611f7157600080fd5b8c6020828501011115611f8357600080fd5b6020830199508098505060808b0135915080821115611fa157600080fd5b611fad8c838d01611ebb565b909750955060a08b0135915080821115611fc657600080fd5b50611fd38b828c01611ebb565b999c989b50969995989497949560c00135949350505050565b602080825282518282018190526000919060409081850190868401855b8281101561205a578151805173ffffffffffffffffffffffffffffffffffffffff1685528681015167ffffffffffffffff168786015285015115158585015260609093019290850190600101612009565b5091979650505050505050565b60006020828403121561207957600080fd5b611b2382611ca4565b6000806020838503121561209557600080fd5b823567ffffffffffffffff808211156120ad57600080fd5b818501915085601f8301126120c157600080fd5b8135818111156120d057600080fd5b8660206060830285010111156120e557600080fd5b60209290920196919550909350505050565b60006020828403121561210957600080fd5b81358015158114611b2357600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b808202811582820484141761215f5761215f612119565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036121c5576121c5612119565b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b600063ffffffff80831681810361221457612214612119565b6001019392505050565b600061012063ffffffff808d1684528b6020850152808b1660408501525080606084015261224e8184018a611bc2565b905082810360808401526122628189611bc2565b905060ff871660a084015282810360c084015261227f8187611aac565b905067ffffffffffffffff851660e08401528281036101008401526122a48185611aac565b9c9b505050505050505050505050565b60ff818116838216019081111561215f5761215f612119565b8082018082111561215f5761215f612119565b8183823760009101908152919050565b828152606082602083013760800192915050565b60006101208b835273ffffffffffffffffffffffffffffffffffffffff8b16602084015267ffffffffffffffff808b16604085015281606085015261234b8285018b611bc2565b9150838203608085015261235f828a611bc2565b915060ff881660a085015283820360c085015261237c8288611aac565b90861660e085015283810361010085015290506122a48185611aac56fea164736f6c6343000813000a",
}

var DummyLiquidityManagerABI = DummyLiquidityManagerMetaData.ABI

var DummyLiquidityManagerBin = DummyLiquidityManagerMetaData.Bin

func DeployDummyLiquidityManager(auth *bind.TransactOpts, backend bind.ContractBackend, localChainSelector uint64) (common.Address, *types.Transaction, *DummyLiquidityManager, error) {
	parsed, err := DummyLiquidityManagerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DummyLiquidityManagerBin), backend, localChainSelector)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DummyLiquidityManager{address: address, abi: *parsed, DummyLiquidityManagerCaller: DummyLiquidityManagerCaller{contract: contract}, DummyLiquidityManagerTransactor: DummyLiquidityManagerTransactor{contract: contract}, DummyLiquidityManagerFilterer: DummyLiquidityManagerFilterer{contract: contract}}, nil
}

type DummyLiquidityManager struct {
	address common.Address
	abi     abi.ABI
	DummyLiquidityManagerCaller
	DummyLiquidityManagerTransactor
	DummyLiquidityManagerFilterer
}

type DummyLiquidityManagerCaller struct {
	contract *bind.BoundContract
}

type DummyLiquidityManagerTransactor struct {
	contract *bind.BoundContract
}

type DummyLiquidityManagerFilterer struct {
	contract *bind.BoundContract
}

type DummyLiquidityManagerSession struct {
	Contract     *DummyLiquidityManager
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type DummyLiquidityManagerCallerSession struct {
	Contract *DummyLiquidityManagerCaller
	CallOpts bind.CallOpts
}

type DummyLiquidityManagerTransactorSession struct {
	Contract     *DummyLiquidityManagerTransactor
	TransactOpts bind.TransactOpts
}

type DummyLiquidityManagerRaw struct {
	Contract *DummyLiquidityManager
}

type DummyLiquidityManagerCallerRaw struct {
	Contract *DummyLiquidityManagerCaller
}

type DummyLiquidityManagerTransactorRaw struct {
	Contract *DummyLiquidityManagerTransactor
}

func NewDummyLiquidityManager(address common.Address, backend bind.ContractBackend) (*DummyLiquidityManager, error) {
	abi, err := abi.JSON(strings.NewReader(DummyLiquidityManagerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindDummyLiquidityManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DummyLiquidityManager{address: address, abi: abi, DummyLiquidityManagerCaller: DummyLiquidityManagerCaller{contract: contract}, DummyLiquidityManagerTransactor: DummyLiquidityManagerTransactor{contract: contract}, DummyLiquidityManagerFilterer: DummyLiquidityManagerFilterer{contract: contract}}, nil
}

func NewDummyLiquidityManagerCaller(address common.Address, caller bind.ContractCaller) (*DummyLiquidityManagerCaller, error) {
	contract, err := bindDummyLiquidityManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DummyLiquidityManagerCaller{contract: contract}, nil
}

func NewDummyLiquidityManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*DummyLiquidityManagerTransactor, error) {
	contract, err := bindDummyLiquidityManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DummyLiquidityManagerTransactor{contract: contract}, nil
}

func NewDummyLiquidityManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*DummyLiquidityManagerFilterer, error) {
	contract, err := bindDummyLiquidityManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DummyLiquidityManagerFilterer{contract: contract}, nil
}

func bindDummyLiquidityManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DummyLiquidityManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_DummyLiquidityManager *DummyLiquidityManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DummyLiquidityManager.Contract.DummyLiquidityManagerCaller.contract.Call(opts, result, method, params...)
}

func (_DummyLiquidityManager *DummyLiquidityManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DummyLiquidityManager.Contract.DummyLiquidityManagerTransactor.contract.Transfer(opts)
}

func (_DummyLiquidityManager *DummyLiquidityManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DummyLiquidityManager.Contract.DummyLiquidityManagerTransactor.contract.Transact(opts, method, params...)
}

func (_DummyLiquidityManager *DummyLiquidityManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DummyLiquidityManager.Contract.contract.Call(opts, result, method, params...)
}

func (_DummyLiquidityManager *DummyLiquidityManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DummyLiquidityManager.Contract.contract.Transfer(opts)
}

func (_DummyLiquidityManager *DummyLiquidityManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DummyLiquidityManager.Contract.contract.Transact(opts, method, params...)
}

func (_DummyLiquidityManager *DummyLiquidityManagerCaller) GetAllCrossChainLiquidityMangers(opts *bind.CallOpts) ([]ILiquidityManagerCrossChainLiquidityManagerArgs, error) {
	var out []interface{}
	err := _DummyLiquidityManager.contract.Call(opts, &out, "getAllCrossChainLiquidityMangers")

	if err != nil {
		return *new([]ILiquidityManagerCrossChainLiquidityManagerArgs), err
	}

	out0 := *abi.ConvertType(out[0], new([]ILiquidityManagerCrossChainLiquidityManagerArgs)).(*[]ILiquidityManagerCrossChainLiquidityManagerArgs)

	return out0, err

}

func (_DummyLiquidityManager *DummyLiquidityManagerSession) GetAllCrossChainLiquidityMangers() ([]ILiquidityManagerCrossChainLiquidityManagerArgs, error) {
	return _DummyLiquidityManager.Contract.GetAllCrossChainLiquidityMangers(&_DummyLiquidityManager.CallOpts)
}

func (_DummyLiquidityManager *DummyLiquidityManagerCallerSession) GetAllCrossChainLiquidityMangers() ([]ILiquidityManagerCrossChainLiquidityManagerArgs, error) {
	return _DummyLiquidityManager.Contract.GetAllCrossChainLiquidityMangers(&_DummyLiquidityManager.CallOpts)
}

func (_DummyLiquidityManager *DummyLiquidityManagerCaller) GetCrossChainLiquidityManager(opts *bind.CallOpts, chainSelector uint64) (DummyLiquidityManagerCrossChainLiquidityManager, error) {
	var out []interface{}
	err := _DummyLiquidityManager.contract.Call(opts, &out, "getCrossChainLiquidityManager", chainSelector)

	if err != nil {
		return *new(DummyLiquidityManagerCrossChainLiquidityManager), err
	}

	out0 := *abi.ConvertType(out[0], new(DummyLiquidityManagerCrossChainLiquidityManager)).(*DummyLiquidityManagerCrossChainLiquidityManager)

	return out0, err

}

func (_DummyLiquidityManager *DummyLiquidityManagerSession) GetCrossChainLiquidityManager(chainSelector uint64) (DummyLiquidityManagerCrossChainLiquidityManager, error) {
	return _DummyLiquidityManager.Contract.GetCrossChainLiquidityManager(&_DummyLiquidityManager.CallOpts, chainSelector)
}

func (_DummyLiquidityManager *DummyLiquidityManagerCallerSession) GetCrossChainLiquidityManager(chainSelector uint64) (DummyLiquidityManagerCrossChainLiquidityManager, error) {
	return _DummyLiquidityManager.Contract.GetCrossChainLiquidityManager(&_DummyLiquidityManager.CallOpts, chainSelector)
}

func (_DummyLiquidityManager *DummyLiquidityManagerCaller) GetSupportedDestChains(opts *bind.CallOpts) ([]uint64, error) {
	var out []interface{}
	err := _DummyLiquidityManager.contract.Call(opts, &out, "getSupportedDestChains")

	if err != nil {
		return *new([]uint64), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint64)).(*[]uint64)

	return out0, err

}

func (_DummyLiquidityManager *DummyLiquidityManagerSession) GetSupportedDestChains() ([]uint64, error) {
	return _DummyLiquidityManager.Contract.GetSupportedDestChains(&_DummyLiquidityManager.CallOpts)
}

func (_DummyLiquidityManager *DummyLiquidityManagerCallerSession) GetSupportedDestChains() ([]uint64, error) {
	return _DummyLiquidityManager.Contract.GetSupportedDestChains(&_DummyLiquidityManager.CallOpts)
}

func (_DummyLiquidityManager *DummyLiquidityManagerCaller) GetTransmitters(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _DummyLiquidityManager.contract.Call(opts, &out, "getTransmitters")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_DummyLiquidityManager *DummyLiquidityManagerSession) GetTransmitters() ([]common.Address, error) {
	return _DummyLiquidityManager.Contract.GetTransmitters(&_DummyLiquidityManager.CallOpts)
}

func (_DummyLiquidityManager *DummyLiquidityManagerCallerSession) GetTransmitters() ([]common.Address, error) {
	return _DummyLiquidityManager.Contract.GetTransmitters(&_DummyLiquidityManager.CallOpts)
}

func (_DummyLiquidityManager *DummyLiquidityManagerCaller) LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

	error) {
	var out []interface{}
	err := _DummyLiquidityManager.contract.Call(opts, &out, "latestConfigDetails")

	outstruct := new(LatestConfigDetails)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ConfigDigest = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

func (_DummyLiquidityManager *DummyLiquidityManagerSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _DummyLiquidityManager.Contract.LatestConfigDetails(&_DummyLiquidityManager.CallOpts)
}

func (_DummyLiquidityManager *DummyLiquidityManagerCallerSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _DummyLiquidityManager.Contract.LatestConfigDetails(&_DummyLiquidityManager.CallOpts)
}

func (_DummyLiquidityManager *DummyLiquidityManagerCaller) LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

	error) {
	var out []interface{}
	err := _DummyLiquidityManager.contract.Call(opts, &out, "latestConfigDigestAndEpoch")

	outstruct := new(LatestConfigDigestAndEpoch)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ScanLogs = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfigDigest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.SequenceNumber = *abi.ConvertType(out[2], new(uint64)).(*uint64)

	return *outstruct, err

}

func (_DummyLiquidityManager *DummyLiquidityManagerSession) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _DummyLiquidityManager.Contract.LatestConfigDigestAndEpoch(&_DummyLiquidityManager.CallOpts)
}

func (_DummyLiquidityManager *DummyLiquidityManagerCallerSession) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _DummyLiquidityManager.Contract.LatestConfigDigestAndEpoch(&_DummyLiquidityManager.CallOpts)
}

func (_DummyLiquidityManager *DummyLiquidityManagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DummyLiquidityManager.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_DummyLiquidityManager *DummyLiquidityManagerSession) Owner() (common.Address, error) {
	return _DummyLiquidityManager.Contract.Owner(&_DummyLiquidityManager.CallOpts)
}

func (_DummyLiquidityManager *DummyLiquidityManagerCallerSession) Owner() (common.Address, error) {
	return _DummyLiquidityManager.Contract.Owner(&_DummyLiquidityManager.CallOpts)
}

func (_DummyLiquidityManager *DummyLiquidityManagerCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _DummyLiquidityManager.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_DummyLiquidityManager *DummyLiquidityManagerSession) TypeAndVersion() (string, error) {
	return _DummyLiquidityManager.Contract.TypeAndVersion(&_DummyLiquidityManager.CallOpts)
}

func (_DummyLiquidityManager *DummyLiquidityManagerCallerSession) TypeAndVersion() (string, error) {
	return _DummyLiquidityManager.Contract.TypeAndVersion(&_DummyLiquidityManager.CallOpts)
}

func (_DummyLiquidityManager *DummyLiquidityManagerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DummyLiquidityManager.contract.Transact(opts, "acceptOwnership")
}

func (_DummyLiquidityManager *DummyLiquidityManagerSession) AcceptOwnership() (*types.Transaction, error) {
	return _DummyLiquidityManager.Contract.AcceptOwnership(&_DummyLiquidityManager.TransactOpts)
}

func (_DummyLiquidityManager *DummyLiquidityManagerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _DummyLiquidityManager.Contract.AcceptOwnership(&_DummyLiquidityManager.TransactOpts)
}

func (_DummyLiquidityManager *DummyLiquidityManagerTransactor) SetCrossChainLiquidityManager(opts *bind.TransactOpts, crossChainLiqManager ILiquidityManagerCrossChainLiquidityManagerArgs) (*types.Transaction, error) {
	return _DummyLiquidityManager.contract.Transact(opts, "setCrossChainLiquidityManager", crossChainLiqManager)
}

func (_DummyLiquidityManager *DummyLiquidityManagerSession) SetCrossChainLiquidityManager(crossChainLiqManager ILiquidityManagerCrossChainLiquidityManagerArgs) (*types.Transaction, error) {
	return _DummyLiquidityManager.Contract.SetCrossChainLiquidityManager(&_DummyLiquidityManager.TransactOpts, crossChainLiqManager)
}

func (_DummyLiquidityManager *DummyLiquidityManagerTransactorSession) SetCrossChainLiquidityManager(crossChainLiqManager ILiquidityManagerCrossChainLiquidityManagerArgs) (*types.Transaction, error) {
	return _DummyLiquidityManager.Contract.SetCrossChainLiquidityManager(&_DummyLiquidityManager.TransactOpts, crossChainLiqManager)
}

func (_DummyLiquidityManager *DummyLiquidityManagerTransactor) SetCrossChainLiquidityManager0(opts *bind.TransactOpts, crossChainLiquidityManagers []ILiquidityManagerCrossChainLiquidityManagerArgs) (*types.Transaction, error) {
	return _DummyLiquidityManager.contract.Transact(opts, "setCrossChainLiquidityManager0", crossChainLiquidityManagers)
}

func (_DummyLiquidityManager *DummyLiquidityManagerSession) SetCrossChainLiquidityManager0(crossChainLiquidityManagers []ILiquidityManagerCrossChainLiquidityManagerArgs) (*types.Transaction, error) {
	return _DummyLiquidityManager.Contract.SetCrossChainLiquidityManager0(&_DummyLiquidityManager.TransactOpts, crossChainLiquidityManagers)
}

func (_DummyLiquidityManager *DummyLiquidityManagerTransactorSession) SetCrossChainLiquidityManager0(crossChainLiquidityManagers []ILiquidityManagerCrossChainLiquidityManagerArgs) (*types.Transaction, error) {
	return _DummyLiquidityManager.Contract.SetCrossChainLiquidityManager0(&_DummyLiquidityManager.TransactOpts, crossChainLiquidityManagers)
}

func (_DummyLiquidityManager *DummyLiquidityManagerTransactor) SetOCR3Config(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _DummyLiquidityManager.contract.Transact(opts, "setOCR3Config", signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_DummyLiquidityManager *DummyLiquidityManagerSession) SetOCR3Config(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _DummyLiquidityManager.Contract.SetOCR3Config(&_DummyLiquidityManager.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_DummyLiquidityManager *DummyLiquidityManagerTransactorSession) SetOCR3Config(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _DummyLiquidityManager.Contract.SetOCR3Config(&_DummyLiquidityManager.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_DummyLiquidityManager *DummyLiquidityManagerTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _DummyLiquidityManager.contract.Transact(opts, "transferOwnership", to)
}

func (_DummyLiquidityManager *DummyLiquidityManagerSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _DummyLiquidityManager.Contract.TransferOwnership(&_DummyLiquidityManager.TransactOpts, to)
}

func (_DummyLiquidityManager *DummyLiquidityManagerTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _DummyLiquidityManager.Contract.TransferOwnership(&_DummyLiquidityManager.TransactOpts, to)
}

func (_DummyLiquidityManager *DummyLiquidityManagerTransactor) Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _DummyLiquidityManager.contract.Transact(opts, "transmit", reportContext, report, rs, ss, rawVs)
}

func (_DummyLiquidityManager *DummyLiquidityManagerSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _DummyLiquidityManager.Contract.Transmit(&_DummyLiquidityManager.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_DummyLiquidityManager *DummyLiquidityManagerTransactorSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _DummyLiquidityManager.Contract.Transmit(&_DummyLiquidityManager.TransactOpts, reportContext, report, rs, ss, rawVs)
}

type DummyLiquidityManagerConfigSetIterator struct {
	Event *DummyLiquidityManagerConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DummyLiquidityManagerConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DummyLiquidityManagerConfigSet)
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
		it.Event = new(DummyLiquidityManagerConfigSet)
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

func (it *DummyLiquidityManagerConfigSetIterator) Error() error {
	return it.fail
}

func (it *DummyLiquidityManagerConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DummyLiquidityManagerConfigSet struct {
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

func (_DummyLiquidityManager *DummyLiquidityManagerFilterer) FilterConfigSet(opts *bind.FilterOpts) (*DummyLiquidityManagerConfigSetIterator, error) {

	logs, sub, err := _DummyLiquidityManager.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &DummyLiquidityManagerConfigSetIterator{contract: _DummyLiquidityManager.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_DummyLiquidityManager *DummyLiquidityManagerFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *DummyLiquidityManagerConfigSet) (event.Subscription, error) {

	logs, sub, err := _DummyLiquidityManager.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DummyLiquidityManagerConfigSet)
				if err := _DummyLiquidityManager.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_DummyLiquidityManager *DummyLiquidityManagerFilterer) ParseConfigSet(log types.Log) (*DummyLiquidityManagerConfigSet, error) {
	event := new(DummyLiquidityManagerConfigSet)
	if err := _DummyLiquidityManager.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DummyLiquidityManagerLiquidityAddedIterator struct {
	Event *DummyLiquidityManagerLiquidityAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DummyLiquidityManagerLiquidityAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DummyLiquidityManagerLiquidityAdded)
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
		it.Event = new(DummyLiquidityManagerLiquidityAdded)
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

func (it *DummyLiquidityManagerLiquidityAddedIterator) Error() error {
	return it.fail
}

func (it *DummyLiquidityManagerLiquidityAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DummyLiquidityManagerLiquidityAdded struct {
	Provider common.Address
	Amount   *big.Int
	Raw      types.Log
}

func (_DummyLiquidityManager *DummyLiquidityManagerFilterer) FilterLiquidityAdded(opts *bind.FilterOpts, provider []common.Address, amount []*big.Int) (*DummyLiquidityManagerLiquidityAddedIterator, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _DummyLiquidityManager.contract.FilterLogs(opts, "LiquidityAdded", providerRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &DummyLiquidityManagerLiquidityAddedIterator{contract: _DummyLiquidityManager.contract, event: "LiquidityAdded", logs: logs, sub: sub}, nil
}

func (_DummyLiquidityManager *DummyLiquidityManagerFilterer) WatchLiquidityAdded(opts *bind.WatchOpts, sink chan<- *DummyLiquidityManagerLiquidityAdded, provider []common.Address, amount []*big.Int) (event.Subscription, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _DummyLiquidityManager.contract.WatchLogs(opts, "LiquidityAdded", providerRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DummyLiquidityManagerLiquidityAdded)
				if err := _DummyLiquidityManager.contract.UnpackLog(event, "LiquidityAdded", log); err != nil {
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

func (_DummyLiquidityManager *DummyLiquidityManagerFilterer) ParseLiquidityAdded(log types.Log) (*DummyLiquidityManagerLiquidityAdded, error) {
	event := new(DummyLiquidityManagerLiquidityAdded)
	if err := _DummyLiquidityManager.contract.UnpackLog(event, "LiquidityAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DummyLiquidityManagerLiquidityRemovedIterator struct {
	Event *DummyLiquidityManagerLiquidityRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DummyLiquidityManagerLiquidityRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DummyLiquidityManagerLiquidityRemoved)
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
		it.Event = new(DummyLiquidityManagerLiquidityRemoved)
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

func (it *DummyLiquidityManagerLiquidityRemovedIterator) Error() error {
	return it.fail
}

func (it *DummyLiquidityManagerLiquidityRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DummyLiquidityManagerLiquidityRemoved struct {
	Remover common.Address
	Amount  *big.Int
	Raw     types.Log
}

func (_DummyLiquidityManager *DummyLiquidityManagerFilterer) FilterLiquidityRemoved(opts *bind.FilterOpts, remover []common.Address, amount []*big.Int) (*DummyLiquidityManagerLiquidityRemovedIterator, error) {

	var removerRule []interface{}
	for _, removerItem := range remover {
		removerRule = append(removerRule, removerItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _DummyLiquidityManager.contract.FilterLogs(opts, "LiquidityRemoved", removerRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &DummyLiquidityManagerLiquidityRemovedIterator{contract: _DummyLiquidityManager.contract, event: "LiquidityRemoved", logs: logs, sub: sub}, nil
}

func (_DummyLiquidityManager *DummyLiquidityManagerFilterer) WatchLiquidityRemoved(opts *bind.WatchOpts, sink chan<- *DummyLiquidityManagerLiquidityRemoved, remover []common.Address, amount []*big.Int) (event.Subscription, error) {

	var removerRule []interface{}
	for _, removerItem := range remover {
		removerRule = append(removerRule, removerItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _DummyLiquidityManager.contract.WatchLogs(opts, "LiquidityRemoved", removerRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DummyLiquidityManagerLiquidityRemoved)
				if err := _DummyLiquidityManager.contract.UnpackLog(event, "LiquidityRemoved", log); err != nil {
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

func (_DummyLiquidityManager *DummyLiquidityManagerFilterer) ParseLiquidityRemoved(log types.Log) (*DummyLiquidityManagerLiquidityRemoved, error) {
	event := new(DummyLiquidityManagerLiquidityRemoved)
	if err := _DummyLiquidityManager.contract.UnpackLog(event, "LiquidityRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DummyLiquidityManagerLiquidityTransferredIterator struct {
	Event *DummyLiquidityManagerLiquidityTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DummyLiquidityManagerLiquidityTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DummyLiquidityManagerLiquidityTransferred)
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
		it.Event = new(DummyLiquidityManagerLiquidityTransferred)
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

func (it *DummyLiquidityManagerLiquidityTransferredIterator) Error() error {
	return it.fail
}

func (it *DummyLiquidityManagerLiquidityTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DummyLiquidityManagerLiquidityTransferred struct {
	OcrSeqNum         uint64
	FromChainSelector uint64
	ToChainSelector   uint64
	To                common.Address
	Amount            *big.Int
	Raw               types.Log
}

func (_DummyLiquidityManager *DummyLiquidityManagerFilterer) FilterLiquidityTransferred(opts *bind.FilterOpts, ocrSeqNum []uint64, fromChainSelector []uint64, toChainSelector []uint64) (*DummyLiquidityManagerLiquidityTransferredIterator, error) {

	var ocrSeqNumRule []interface{}
	for _, ocrSeqNumItem := range ocrSeqNum {
		ocrSeqNumRule = append(ocrSeqNumRule, ocrSeqNumItem)
	}
	var fromChainSelectorRule []interface{}
	for _, fromChainSelectorItem := range fromChainSelector {
		fromChainSelectorRule = append(fromChainSelectorRule, fromChainSelectorItem)
	}
	var toChainSelectorRule []interface{}
	for _, toChainSelectorItem := range toChainSelector {
		toChainSelectorRule = append(toChainSelectorRule, toChainSelectorItem)
	}

	logs, sub, err := _DummyLiquidityManager.contract.FilterLogs(opts, "LiquidityTransferred", ocrSeqNumRule, fromChainSelectorRule, toChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &DummyLiquidityManagerLiquidityTransferredIterator{contract: _DummyLiquidityManager.contract, event: "LiquidityTransferred", logs: logs, sub: sub}, nil
}

func (_DummyLiquidityManager *DummyLiquidityManagerFilterer) WatchLiquidityTransferred(opts *bind.WatchOpts, sink chan<- *DummyLiquidityManagerLiquidityTransferred, ocrSeqNum []uint64, fromChainSelector []uint64, toChainSelector []uint64) (event.Subscription, error) {

	var ocrSeqNumRule []interface{}
	for _, ocrSeqNumItem := range ocrSeqNum {
		ocrSeqNumRule = append(ocrSeqNumRule, ocrSeqNumItem)
	}
	var fromChainSelectorRule []interface{}
	for _, fromChainSelectorItem := range fromChainSelector {
		fromChainSelectorRule = append(fromChainSelectorRule, fromChainSelectorItem)
	}
	var toChainSelectorRule []interface{}
	for _, toChainSelectorItem := range toChainSelector {
		toChainSelectorRule = append(toChainSelectorRule, toChainSelectorItem)
	}

	logs, sub, err := _DummyLiquidityManager.contract.WatchLogs(opts, "LiquidityTransferred", ocrSeqNumRule, fromChainSelectorRule, toChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DummyLiquidityManagerLiquidityTransferred)
				if err := _DummyLiquidityManager.contract.UnpackLog(event, "LiquidityTransferred", log); err != nil {
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

func (_DummyLiquidityManager *DummyLiquidityManagerFilterer) ParseLiquidityTransferred(log types.Log) (*DummyLiquidityManagerLiquidityTransferred, error) {
	event := new(DummyLiquidityManagerLiquidityTransferred)
	if err := _DummyLiquidityManager.contract.UnpackLog(event, "LiquidityTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DummyLiquidityManagerOwnershipTransferRequestedIterator struct {
	Event *DummyLiquidityManagerOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DummyLiquidityManagerOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DummyLiquidityManagerOwnershipTransferRequested)
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
		it.Event = new(DummyLiquidityManagerOwnershipTransferRequested)
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

func (it *DummyLiquidityManagerOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *DummyLiquidityManagerOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DummyLiquidityManagerOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_DummyLiquidityManager *DummyLiquidityManagerFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*DummyLiquidityManagerOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _DummyLiquidityManager.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &DummyLiquidityManagerOwnershipTransferRequestedIterator{contract: _DummyLiquidityManager.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_DummyLiquidityManager *DummyLiquidityManagerFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *DummyLiquidityManagerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _DummyLiquidityManager.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DummyLiquidityManagerOwnershipTransferRequested)
				if err := _DummyLiquidityManager.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_DummyLiquidityManager *DummyLiquidityManagerFilterer) ParseOwnershipTransferRequested(log types.Log) (*DummyLiquidityManagerOwnershipTransferRequested, error) {
	event := new(DummyLiquidityManagerOwnershipTransferRequested)
	if err := _DummyLiquidityManager.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DummyLiquidityManagerOwnershipTransferredIterator struct {
	Event *DummyLiquidityManagerOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DummyLiquidityManagerOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DummyLiquidityManagerOwnershipTransferred)
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
		it.Event = new(DummyLiquidityManagerOwnershipTransferred)
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

func (it *DummyLiquidityManagerOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *DummyLiquidityManagerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DummyLiquidityManagerOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_DummyLiquidityManager *DummyLiquidityManagerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*DummyLiquidityManagerOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _DummyLiquidityManager.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &DummyLiquidityManagerOwnershipTransferredIterator{contract: _DummyLiquidityManager.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_DummyLiquidityManager *DummyLiquidityManagerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DummyLiquidityManagerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _DummyLiquidityManager.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DummyLiquidityManagerOwnershipTransferred)
				if err := _DummyLiquidityManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_DummyLiquidityManager *DummyLiquidityManagerFilterer) ParseOwnershipTransferred(log types.Log) (*DummyLiquidityManagerOwnershipTransferred, error) {
	event := new(DummyLiquidityManagerOwnershipTransferred)
	if err := _DummyLiquidityManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DummyLiquidityManagerTransmittedIterator struct {
	Event *DummyLiquidityManagerTransmitted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DummyLiquidityManagerTransmittedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DummyLiquidityManagerTransmitted)
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
		it.Event = new(DummyLiquidityManagerTransmitted)
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

func (it *DummyLiquidityManagerTransmittedIterator) Error() error {
	return it.fail
}

func (it *DummyLiquidityManagerTransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DummyLiquidityManagerTransmitted struct {
	ConfigDigest   [32]byte
	SequenceNumber uint64
	Raw            types.Log
}

func (_DummyLiquidityManager *DummyLiquidityManagerFilterer) FilterTransmitted(opts *bind.FilterOpts) (*DummyLiquidityManagerTransmittedIterator, error) {

	logs, sub, err := _DummyLiquidityManager.contract.FilterLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return &DummyLiquidityManagerTransmittedIterator{contract: _DummyLiquidityManager.contract, event: "Transmitted", logs: logs, sub: sub}, nil
}

func (_DummyLiquidityManager *DummyLiquidityManagerFilterer) WatchTransmitted(opts *bind.WatchOpts, sink chan<- *DummyLiquidityManagerTransmitted) (event.Subscription, error) {

	logs, sub, err := _DummyLiquidityManager.contract.WatchLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DummyLiquidityManagerTransmitted)
				if err := _DummyLiquidityManager.contract.UnpackLog(event, "Transmitted", log); err != nil {
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

func (_DummyLiquidityManager *DummyLiquidityManagerFilterer) ParseTransmitted(log types.Log) (*DummyLiquidityManagerTransmitted, error) {
	event := new(DummyLiquidityManagerTransmitted)
	if err := _DummyLiquidityManager.contract.UnpackLog(event, "Transmitted", log); err != nil {
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

func (_DummyLiquidityManager *DummyLiquidityManager) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _DummyLiquidityManager.abi.Events["ConfigSet"].ID:
		return _DummyLiquidityManager.ParseConfigSet(log)
	case _DummyLiquidityManager.abi.Events["LiquidityAdded"].ID:
		return _DummyLiquidityManager.ParseLiquidityAdded(log)
	case _DummyLiquidityManager.abi.Events["LiquidityRemoved"].ID:
		return _DummyLiquidityManager.ParseLiquidityRemoved(log)
	case _DummyLiquidityManager.abi.Events["LiquidityTransferred"].ID:
		return _DummyLiquidityManager.ParseLiquidityTransferred(log)
	case _DummyLiquidityManager.abi.Events["OwnershipTransferRequested"].ID:
		return _DummyLiquidityManager.ParseOwnershipTransferRequested(log)
	case _DummyLiquidityManager.abi.Events["OwnershipTransferred"].ID:
		return _DummyLiquidityManager.ParseOwnershipTransferred(log)
	case _DummyLiquidityManager.abi.Events["Transmitted"].ID:
		return _DummyLiquidityManager.ParseTransmitted(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (DummyLiquidityManagerConfigSet) Topic() common.Hash {
	return common.HexToHash("0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05")
}

func (DummyLiquidityManagerLiquidityAdded) Topic() common.Hash {
	return common.HexToHash("0xc17cea59c2955cb181b03393209566960365771dbba9dc3d510180e7cb312088")
}

func (DummyLiquidityManagerLiquidityRemoved) Topic() common.Hash {
	return common.HexToHash("0xc2c3f06e49b9f15e7b4af9055e183b0d73362e033ad82a07dec9bf9840171719")
}

func (DummyLiquidityManagerLiquidityTransferred) Topic() common.Hash {
	return common.HexToHash("0xc3699d6ab2762f468855d25ae224adfffe5b2bc4d57ca590cdfdbbbeefbee22f")
}

func (DummyLiquidityManagerOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (DummyLiquidityManagerOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (DummyLiquidityManagerTransmitted) Topic() common.Hash {
	return common.HexToHash("0xe893c2681d327421d89e1cb54fbe64645b4dcea668d6826130b62cf4c6eefea2")
}

func (_DummyLiquidityManager *DummyLiquidityManager) Address() common.Address {
	return _DummyLiquidityManager.address
}

type DummyLiquidityManagerInterface interface {
	GetAllCrossChainLiquidityMangers(opts *bind.CallOpts) ([]ILiquidityManagerCrossChainLiquidityManagerArgs, error)

	GetCrossChainLiquidityManager(opts *bind.CallOpts, chainSelector uint64) (DummyLiquidityManagerCrossChainLiquidityManager, error)

	GetSupportedDestChains(opts *bind.CallOpts) ([]uint64, error)

	GetTransmitters(opts *bind.CallOpts) ([]common.Address, error)

	LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

		error)

	LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

		error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	SetCrossChainLiquidityManager(opts *bind.TransactOpts, crossChainLiqManager ILiquidityManagerCrossChainLiquidityManagerArgs) (*types.Transaction, error)

	SetCrossChainLiquidityManager0(opts *bind.TransactOpts, crossChainLiquidityManagers []ILiquidityManagerCrossChainLiquidityManagerArgs) (*types.Transaction, error)

	SetOCR3Config(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error)

	FilterConfigSet(opts *bind.FilterOpts) (*DummyLiquidityManagerConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *DummyLiquidityManagerConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*DummyLiquidityManagerConfigSet, error)

	FilterLiquidityAdded(opts *bind.FilterOpts, provider []common.Address, amount []*big.Int) (*DummyLiquidityManagerLiquidityAddedIterator, error)

	WatchLiquidityAdded(opts *bind.WatchOpts, sink chan<- *DummyLiquidityManagerLiquidityAdded, provider []common.Address, amount []*big.Int) (event.Subscription, error)

	ParseLiquidityAdded(log types.Log) (*DummyLiquidityManagerLiquidityAdded, error)

	FilterLiquidityRemoved(opts *bind.FilterOpts, remover []common.Address, amount []*big.Int) (*DummyLiquidityManagerLiquidityRemovedIterator, error)

	WatchLiquidityRemoved(opts *bind.WatchOpts, sink chan<- *DummyLiquidityManagerLiquidityRemoved, remover []common.Address, amount []*big.Int) (event.Subscription, error)

	ParseLiquidityRemoved(log types.Log) (*DummyLiquidityManagerLiquidityRemoved, error)

	FilterLiquidityTransferred(opts *bind.FilterOpts, ocrSeqNum []uint64, fromChainSelector []uint64, toChainSelector []uint64) (*DummyLiquidityManagerLiquidityTransferredIterator, error)

	WatchLiquidityTransferred(opts *bind.WatchOpts, sink chan<- *DummyLiquidityManagerLiquidityTransferred, ocrSeqNum []uint64, fromChainSelector []uint64, toChainSelector []uint64) (event.Subscription, error)

	ParseLiquidityTransferred(log types.Log) (*DummyLiquidityManagerLiquidityTransferred, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*DummyLiquidityManagerOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *DummyLiquidityManagerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*DummyLiquidityManagerOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*DummyLiquidityManagerOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DummyLiquidityManagerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*DummyLiquidityManagerOwnershipTransferred, error)

	FilterTransmitted(opts *bind.FilterOpts) (*DummyLiquidityManagerTransmittedIterator, error)

	WatchTransmitted(opts *bind.WatchOpts, sink chan<- *DummyLiquidityManagerTransmitted) (event.Subscription, error)

	ParseTransmitted(log types.Log) (*DummyLiquidityManagerTransmitted, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
