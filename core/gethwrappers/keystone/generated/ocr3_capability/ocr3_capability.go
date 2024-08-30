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
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"InvalidConfig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReportingUnsupported\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes[]\",\"name\":\"signers\",\"type\":\"bytes[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"_signers\",\"type\":\"bytes[]\"},{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"_f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"_onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"_offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"_offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5033806000816100675760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615610097576100978161009f565b505050610148565b336001600160a01b038216036100f75760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161005e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6115fd806101576000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c80638da5cb5b1161005b5780638da5cb5b1461015f578063afcb95d714610187578063b1dc65a4146101a7578063f2fde38b146101ba57600080fd5b8063181f5a771461008d57806379ba5097146100d55780637f3c87d3146100df57806381ff7048146100f2575b600080fd5b604080518082018252600e81527f4b657973746f6e6520312e302e30000000000000000000000000000000000000602082015290516100cc9190610ca8565b60405180910390f35b6100dd6101cd565b005b6100dd6100ed366004610e48565b6102cf565b61013c60015460025463ffffffff74010000000000000000000000000000000000000000830481169378010000000000000000000000000000000000000000000000009093041691565b6040805163ffffffff9485168152939092166020840152908201526060016100cc565b60005460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100cc565b6040805160018152600060208201819052918101919091526060016100cc565b6100dd6101b5366004610f24565b61098d565b6100dd6101c8366004611009565b6109bf565b60015473ffffffffffffffffffffffffffffffffffffffff163314610253576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b868560ff8616601f831115610340576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f746f6f206d616e79207369676e65727300000000000000000000000000000000604482015260640161024a565b806000036103aa576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f66206d75737420626520706f7369746976650000000000000000000000000000604482015260640161024a565b818314610438576040517f89a61989000000000000000000000000000000000000000000000000000000008152602060048201526024808201527f6f7261636c6520616464726573736573206f7574206f6620726567697374726160448201527f74696f6e00000000000000000000000000000000000000000000000000000000606482015260840161024a565b61044381600361106e565b83116104ab576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f6661756c74792d6f7261636c65206620746f6f20686967680000000000000000604482015260640161024a565b6104b36109d3565b60006040518060c001604052808d8d906104cd919061108b565b81526020018b8b8080602002602001604051908101604052809392919081815260200183836020028082843760009201919091525050509082525060ff8a1660208201526040810189905267ffffffffffffffff8816606082015260800186905290505b6004541561056d57600480548061054a5761054a611111565b6001900381819060005260206000200160006105669190610bf6565b9055610531565b60005b8151518110156107a857600073ffffffffffffffffffffffffffffffffffffffff16826020015182815181106105a8576105a8611140565b602002602001015173ffffffffffffffffffffffffffffffffffffffff160361062d576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f7472616e736d6974746572206d757374206e6f7420626520656d707479000000604482015260640161024a565b3660008e8e8481811061064257610642611140565b9050602002810190610654919061116f565b90925090506000815b8061ffff168261ffff16101561075657600084848461ffff1681811061068557610685611140565b919091013560f81c915060009050600886866106a28760026111d4565b61ffff168181106106b5576106b5611140565b919091013560f81c90911b905086866106cf8760016111d4565b61ffff168181106106e2576106e2611140565b6106f39392013560f81c90506111d4565b905036600087876107058860036111d4565b61ffff1690856107168a60036111d4565b61072091906111d4565b61ffff1692610731939291906111f6565b90925090506107418360036111d4565b61074b90876111d4565b95505050505061065d565b60048660000151868151811061076e5761076e611140565b6020908102919091018101518254600181018455600093845291909220019061079790826112c4565b505060019093019250610570915050565b506040810151600380547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff909216919091179055600180547fffffffff00000000ffffffffffffffffffffffffffffffffffffffffffffffff8116780100000000000000000000000000000000000000000000000063ffffffff4381168202929092178085559204811692918291601491610860918491740100000000000000000000000000000000000000009004166113de565b92506101000a81548163ffffffff021916908363ffffffff1602179055506108bf4630600160149054906101000a900463ffffffff1663ffffffff16856000015186602001518760400151886060015189608001518a60a00151610a56565b600281905582518051600380547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1661010060ff9093169290920291909117905560015460208501516040808701516060880151608089015160a08a015193517f36257c6e8d535293ad661e377c0baac536289be6707b8a488ac175ddaa4055c898610976988b9891977401000000000000000000000000000000000000000090920463ffffffff169690959194919391926114c5565b60405180910390a150505050505050505050505050565b6040517f0750181900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6109c76109d3565b6109d081610b01565b50565b60005473ffffffffffffffffffffffffffffffffffffffff163314610a54576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161024a565b565b6000808a8a8a8a8a8a8a8a8a604051602001610a7a9998979695949392919061155b565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179150509998505050505050505050565b3373ffffffffffffffffffffffffffffffffffffffff821603610b80576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161024a565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b508054610c0290611220565b6000825580601f10610c12575050565b601f0160209004906000526020600020908101906109d091905b80821115610c405760008155600101610c2c565b5090565b6000815180845260005b81811015610c6a57602081850181015186830182015201610c4e565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b602081526000610cbb6020830184610c44565b9392505050565b60008083601f840112610cd457600080fd5b50813567ffffffffffffffff811115610cec57600080fd5b6020830191508360208260051b8501011115610d0757600080fd5b9250929050565b803560ff81168114610d1f57600080fd5b919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610d9a57610d9a610d24565b604052919050565b600082601f830112610db357600080fd5b813567ffffffffffffffff811115610dcd57610dcd610d24565b610dfe60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601610d53565b818152846020838601011115610e1357600080fd5b816020850160208301376000918101602001919091529392505050565b803567ffffffffffffffff81168114610d1f57600080fd5b60008060008060008060008060c0898b031215610e6457600080fd5b883567ffffffffffffffff80821115610e7c57600080fd5b610e888c838d01610cc2565b909a50985060208b0135915080821115610ea157600080fd5b610ead8c838d01610cc2565b9098509650869150610ec160408c01610d0e565b955060608b0135915080821115610ed757600080fd5b610ee38c838d01610da2565b9450610ef160808c01610e30565b935060a08b0135915080821115610f0757600080fd5b50610f148b828c01610da2565b9150509295985092959890939650565b60008060008060008060008060e0898b031215610f4057600080fd5b606089018a811115610f5157600080fd5b8998503567ffffffffffffffff80821115610f6b57600080fd5b818b0191508b601f830112610f7f57600080fd5b813581811115610f8e57600080fd5b8c6020828501011115610fa057600080fd5b6020830199508098505060808b0135915080821115610fbe57600080fd5b610fca8c838d01610cc2565b909750955060a08b0135915080821115610fe357600080fd5b50610ff08b828c01610cc2565b999c989b50969995989497949560c00135949350505050565b60006020828403121561101b57600080fd5b813573ffffffffffffffffffffffffffffffffffffffff81168114610cbb57600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b80820281158282048414176110855761108561103f565b92915050565b600067ffffffffffffffff808411156110a6576110a6610d24565b8360051b60206110b860208301610d53565b868152918501916020810190368411156110d157600080fd5b865b84811015611105578035868111156110eb5760008081fd5b6110f736828b01610da2565b8452509183019183016110d3565b50979650505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18436030181126111a457600080fd5b83018035915067ffffffffffffffff8211156111bf57600080fd5b602001915036819003821315610d0757600080fd5b61ffff8181168382160190808211156111ef576111ef61103f565b5092915050565b6000808585111561120657600080fd5b8386111561121357600080fd5b5050820193919092039150565b600181811c9082168061123457607f821691505b60208210810361126d577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f8211156112bf576000816000526020600020601f850160051c8101602086101561129c5750805b601f850160051c820191505b818110156112bb578281556001016112a8565b5050505b505050565b815167ffffffffffffffff8111156112de576112de610d24565b6112f2816112ec8454611220565b84611273565b602080601f831160018114611345576000841561130f5750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b1785556112bb565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b8281101561139257888601518255948401946001909101908401611373565b50858210156113ce57878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b63ffffffff8181168382160190808211156111ef576111ef61103f565b60008282518085526020808601955060208260051b8401016020860160005b84811015611466577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0868403018952611454838351610c44565b9884019892509083019060010161141a565b5090979650505050505050565b60008151808452602080850194506020840160005b838110156114ba57815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101611488565b509495945050505050565b600061012063ffffffff808d1684528b6020850152808b166040850152508060608401526114f58184018a6113fb565b905082810360808401526115098189611473565b905060ff871660a084015282810360c08401526115268187610c44565b905067ffffffffffffffff851660e084015282810361010084015261154b8185610c44565b9c9b505050505050505050505050565b60006101208b835273ffffffffffffffffffffffffffffffffffffffff8b16602084015267ffffffffffffffff808b1660408501528160608501526115a28285018b6113fb565b915083820360808501526115b6828a611473565b915060ff881660a085015283820360c08501526115d38288610c44565b90861660e0850152838103610100850152905061154b8185610c4456fea164736f6c6343000818000a",
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

func (_OCR3Capability *OCR3CapabilityCaller) Transmit(opts *bind.CallOpts, arg0 [3][32]byte, arg1 []byte, arg2 [][32]byte, arg3 [][32]byte, arg4 [32]byte) error {
	var out []interface{}
	err := _OCR3Capability.contract.Call(opts, &out, "transmit", arg0, arg1, arg2, arg3, arg4)

	if err != nil {
		return err
	}

	return err

}

func (_OCR3Capability *OCR3CapabilitySession) Transmit(arg0 [3][32]byte, arg1 []byte, arg2 [][32]byte, arg3 [][32]byte, arg4 [32]byte) error {
	return _OCR3Capability.Contract.Transmit(&_OCR3Capability.CallOpts, arg0, arg1, arg2, arg3, arg4)
}

func (_OCR3Capability *OCR3CapabilityCallerSession) Transmit(arg0 [3][32]byte, arg1 []byte, arg2 [][32]byte, arg3 [][32]byte, arg4 [32]byte) error {
	return _OCR3Capability.Contract.Transmit(&_OCR3Capability.CallOpts, arg0, arg1, arg2, arg3, arg4)
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

func (_OCR3Capability *OCR3CapabilityTransactor) SetConfig(opts *bind.TransactOpts, _signers [][]byte, _transmitters []common.Address, _f uint8, _onchainConfig []byte, _offchainConfigVersion uint64, _offchainConfig []byte) (*types.Transaction, error) {
	return _OCR3Capability.contract.Transact(opts, "setConfig", _signers, _transmitters, _f, _onchainConfig, _offchainConfigVersion, _offchainConfig)
}

func (_OCR3Capability *OCR3CapabilitySession) SetConfig(_signers [][]byte, _transmitters []common.Address, _f uint8, _onchainConfig []byte, _offchainConfigVersion uint64, _offchainConfig []byte) (*types.Transaction, error) {
	return _OCR3Capability.Contract.SetConfig(&_OCR3Capability.TransactOpts, _signers, _transmitters, _f, _onchainConfig, _offchainConfigVersion, _offchainConfig)
}

func (_OCR3Capability *OCR3CapabilityTransactorSession) SetConfig(_signers [][]byte, _transmitters []common.Address, _f uint8, _onchainConfig []byte, _offchainConfigVersion uint64, _offchainConfig []byte) (*types.Transaction, error) {
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
	Signers                   [][]byte
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
	return common.HexToHash("0x36257c6e8d535293ad661e377c0baac536289be6707b8a488ac175ddaa4055c8")
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

	Transmit(opts *bind.CallOpts, arg0 [3][32]byte, arg1 []byte, arg2 [][32]byte, arg3 [][32]byte, arg4 [32]byte) error

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, _signers [][]byte, _transmitters []common.Address, _f uint8, _onchainConfig []byte, _offchainConfigVersion uint64, _offchainConfig []byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

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
