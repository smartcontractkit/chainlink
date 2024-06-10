// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package forwarder

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

var KeystoneForwarderMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"}],\"name\":\"DuplicateSigner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxSigners\",\"type\":\"uint256\"}],\"name\":\"ExcessSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FaultToleranceMustBePositive\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minSigners\",\"type\":\"uint256\"}],\"name\":\"InsufficientSigners\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"configId\",\"type\":\"uint64\"}],\"name\":\"InvalidConfig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReport\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"InvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"received\",\"type\":\"uint256\"}],\"name\":\"InvalidSignatureCount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"}],\"name\":\"InvalidSigner\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"configVersion\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"workflowExecutionId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes2\",\"name\":\"reportId\",\"type\":\"bytes2\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"result\",\"type\":\"bool\"}],\"name\":\"ReportProcessed\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"configVersion\",\"type\":\"uint32\"}],\"name\":\"clearConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRouter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiverAddress\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"workflowExecutionId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes2\",\"name\":\"reportId\",\"type\":\"bytes2\"}],\"name\":\"getTransmissionId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiverAddress\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"workflowExecutionId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes2\",\"name\":\"reportId\",\"type\":\"bytes2\"}],\"name\":\"getTransmissionState\",\"outputs\":[{\"internalType\":\"enumIRouter.TransmissionState\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiverAddress\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"workflowExecutionId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes2\",\"name\":\"reportId\",\"type\":\"bytes2\"}],\"name\":\"getTransmitter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiverAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"rawReport\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"reportContext\",\"type\":\"bytes\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"report\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"configVersion\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162001b4a38038062001b4a833981016040819052620000349162000193565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000e8565b5050600280546001600160a01b0319166001600160a01b03939093169290921790915550620001c5565b336001600160a01b03821603620001425760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600060208284031215620001a657600080fd5b81516001600160a01b0381168114620001be57600080fd5b9392505050565b61197580620001d56000396000f3fe608060405234801561001057600080fd5b50600436106100c95760003560e01c80638864b86411610081578063ee59d26c1161005b578063ee59d26c14610270578063ef6e17a014610283578063f2fde38b1461029657600080fd5b80638864b864146101fc5780638da5cb5b14610234578063b0f479a11461025257600080fd5b8063354bdd66116100b2578063354bdd661461013557806343c16467146101d457806379ba5097146101f457600080fd5b806311289565146100ce578063181f5a77146100e3575b600080fd5b6100e16100dc3660046112ec565b6102a9565b005b61011f6040518060400160405280601781526020017f4b657973746f6e65466f7277617264657220312e302e3000000000000000000081525081565b60405161012c9190611399565b60405180910390f35b6101c6610143366004611405565b6040517fffffffffffffffffffffffffffffffffffffffff000000000000000000000000606085901b166020820152603481018390527fffff000000000000000000000000000000000000000000000000000000000000821660548201526000906056016040516020818303038152906040528051906020012090509392505050565b60405190815260200161012c565b6101e76101e2366004611405565b610891565b60405161012c919061146c565b6100e16109a6565b61020f61020a366004611405565b610aa3565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161012c565b60005473ffffffffffffffffffffffffffffffffffffffff1661020f565b60025473ffffffffffffffffffffffffffffffffffffffff1661020f565b6100e161027e3660046114c6565b610bb0565b6100e1610291366004611544565b610f39565b6100e16102a4366004611577565b610fd9565b606d8510156102e4576040517fb55ac75400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080600061032889898080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250610fed92505050565b67ffffffffffffffff8216600090815260036020526040812080549497509195509193509160ff169081900361039b576040517fdf3b81ea00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff841660048201526024015b60405180910390fd5b856103a78260016115ca565b60ff16146103f9576103ba8160016115ca565b6040517fd6022e8e00000000000000000000000000000000000000000000000000000000815260ff909116600482015260248101879052604401610392565b60008b8b60405161040b9291906115e9565b604051908190038120610424918c908c906020016115f9565b604051602081830303815290604052805190602001209050610444611180565b60005b888110156106cd573660008b8b8481811061046457610464611613565b90506020028101906104769190611642565b9092509050604181146104b95781816040517f2adfdc300000000000000000000000000000000000000000000000000000000081526004016103929291906116f0565b6000600186848460408181106104d1576104d1611613565b6104e392013560f81c9050601b6115ca565b6104f1602060008789611704565b6104fa9161172e565b61050860406020888a611704565b6105119161172e565b6040805160008152602081018083529590955260ff909316928401929092526060830152608082015260a0016020604051602081039080840390855afa15801561055f573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0015173ffffffffffffffffffffffffffffffffffffffff8116600090815260028c01602052918220549093509150819003610605576040517fbf18af4300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff83166004820152602401610392565b600086826020811061061957610619611613565b602002015173ffffffffffffffffffffffffffffffffffffffff1614610683576040517fe021c4f200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff83166004820152602401610392565b8186826020811061069657610696611613565b73ffffffffffffffffffffffffffffffffffffffff9092166020929092020152506106c6925083915061176a9050565b9050610447565b50506002546000945073ffffffffffffffffffffffffffffffffffffffff16925063233fd52d915061078090508c86866040517fffffffffffffffffffffffffffffffffffffffff000000000000000000000000606085901b166020820152603481018390527fffff000000000000000000000000000000000000000000000000000000000000821660548201526000906056016040516020818303038152906040528051906020012090509392505050565b338d8d8d602d90606d9261079693929190611704565b8f8f606d9080926107a993929190611704565b6040518863ffffffff1660e01b81526004016107cb97969594939291906117a2565b6020604051808303816000875af11580156107ea573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061080e9190611803565b9050817dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916838b73ffffffffffffffffffffffffffffffffffffffff167f3617b009e9785c42daebadb6d3fb553243a4bf586d07ea72d65d80013ce116b58460405161087d911515815260200190565b60405180910390a450505050505050505050565b600254604080517fffffffffffffffffffffffffffffffffffffffff000000000000000000000000606087901b16602080830191909152603482018690527fffff0000000000000000000000000000000000000000000000000000000000008516605483015282518083036036018152605683019384905280519101207f516db40800000000000000000000000000000000000000000000000000000000909252605a81019190915260009173ffffffffffffffffffffffffffffffffffffffff169063516db40890607a01602060405180830381865afa15801561097a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061099e9190611825565b949350505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610a27576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610392565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b600254604080517fffffffffffffffffffffffffffffffffffffffff000000000000000000000000606087901b16602080830191909152603482018690527fffff0000000000000000000000000000000000000000000000000000000000008516605483015282518083036036018152605683019384905280519101207fe6b7145800000000000000000000000000000000000000000000000000000000909252605a81019190915260009173ffffffffffffffffffffffffffffffffffffffff169063e6b7145890607a01602060405180830381865afa158015610b8c573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061099e9190611846565b610bb8611008565b8260ff16600003610bf5576040517f0743bae600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601f811115610c3a576040517f61750f4000000000000000000000000000000000000000000000000000000000815260048101829052601f6024820152604401610392565b610c45836003611863565b60ff168111610ca35780610c5a846003611863565b610c659060016115ca565b6040517f9dd9e6d8000000000000000000000000000000000000000000000000000000008152600481019290925260ff166024820152604401610392565b67ffffffff00000000602086901b1663ffffffff85161760005b67ffffffffffffffff8216600090815260036020526040902060010154811015610d595767ffffffffffffffff821660009081526003602052604081206001810180546002909201929184908110610d1757610d17611613565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff168352820192909252604001812055610d528161176a565b9050610cbd565b5060005b82811015610e7b576000848483818110610d7957610d79611613565b9050602002016020810190610d8e9190611577565b67ffffffffffffffff8416600090815260036020908152604080832073ffffffffffffffffffffffffffffffffffffffff8516845260020190915290205490915015610e1e576040517fe021c4f200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82166004820152602401610392565b610e29826001611886565b67ffffffffffffffff8416600090815260036020908152604080832073ffffffffffffffffffffffffffffffffffffffff9095168352600290940190529190912055610e748161176a565b9050610d5d565b5067ffffffffffffffff81166000908152600360205260409020610ea390600101848461119f565b5067ffffffffffffffff81166000908152600360205260409081902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff87161790555163ffffffff86811691908816907f4120bd3b23957dd423555817d55654d4481b438aa15485c21b4180c784f1a45590610f2990889088908890611899565b60405180910390a3505050505050565b610f41611008565b63ffffffff818116602084811b67ffffffff00000000168217600090815260038252604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690558051828152928301905291928516917f4120bd3b23957dd423555817d55654d4481b438aa15485c21b4180c784f1a45591604051610fcd929190611901565b60405180910390a35050565b610fe1611008565b610fea8161108b565b50565b60218101516045820151608b90920151909260c09290921c91565b60005473ffffffffffffffffffffffffffffffffffffffff163314611089576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610392565b565b3373ffffffffffffffffffffffffffffffffffffffff82160361110a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610392565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6040518061040001604052806020906020820280368337509192915050565b828054828255906000526020600020908101928215611217579160200282015b828111156112175781547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8435161782556020909201916001909101906111bf565b50611223929150611227565b5090565b5b808211156112235760008155600101611228565b73ffffffffffffffffffffffffffffffffffffffff81168114610fea57600080fd5b60008083601f84011261127057600080fd5b50813567ffffffffffffffff81111561128857600080fd5b6020830191508360208285010111156112a057600080fd5b9250929050565b60008083601f8401126112b957600080fd5b50813567ffffffffffffffff8111156112d157600080fd5b6020830191508360208260051b85010111156112a057600080fd5b60008060008060008060006080888a03121561130757600080fd5b87356113128161123c565b9650602088013567ffffffffffffffff8082111561132f57600080fd5b61133b8b838c0161125e565b909850965060408a013591508082111561135457600080fd5b6113608b838c0161125e565b909650945060608a013591508082111561137957600080fd5b506113868a828b016112a7565b989b979a50959850939692959293505050565b600060208083528351808285015260005b818110156113c6578581018301518582016040015282016113aa565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b60008060006060848603121561141a57600080fd5b83356114258161123c565b92506020840135915060408401357fffff0000000000000000000000000000000000000000000000000000000000008116811461146157600080fd5b809150509250925092565b60208101600383106114a7577f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b91905290565b803563ffffffff811681146114c157600080fd5b919050565b6000806000806000608086880312156114de57600080fd5b6114e7866114ad565b94506114f5602087016114ad565b9350604086013560ff8116811461150b57600080fd5b9250606086013567ffffffffffffffff81111561152757600080fd5b611533888289016112a7565b969995985093965092949392505050565b6000806040838503121561155757600080fd5b611560836114ad565b915061156e602084016114ad565b90509250929050565b60006020828403121561158957600080fd5b81356115948161123c565b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60ff81811683821601908111156115e3576115e361159b565b92915050565b8183823760009101908152919050565b838152818360208301376000910160200190815292915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261167757600080fd5b83018035915067ffffffffffffffff82111561169257600080fd5b6020019150368190038213156112a057600080fd5b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b60208152600061099e6020830184866116a7565b6000808585111561171457600080fd5b8386111561172157600080fd5b5050820193919092039150565b803560208310156115e3577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff602084900360031b1b1692915050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361179b5761179b61159b565b5060010190565b878152600073ffffffffffffffffffffffffffffffffffffffff808916602084015280881660408401525060a060608301526117e260a0830186886116a7565b82810360808401526117f58185876116a7565b9a9950505050505050505050565b60006020828403121561181557600080fd5b8151801515811461159457600080fd5b60006020828403121561183757600080fd5b81516003811061159457600080fd5b60006020828403121561185857600080fd5b81516115948161123c565b60ff818116838216029081169081811461187f5761187f61159b565b5092915050565b808201808211156115e3576115e361159b565b60ff8416815260406020808301829052908201839052600090849060608401835b868110156118f55783356118cd8161123c565b73ffffffffffffffffffffffffffffffffffffffff16825292820192908201906001016118ba565b50979650505050505050565b60006040820160ff851683526020604081850152818551808452606086019150828701935060005b8181101561195b57845173ffffffffffffffffffffffffffffffffffffffff1683529383019391830191600101611929565b509097965050505050505056fea164736f6c6343000813000a",
}

var KeystoneForwarderABI = KeystoneForwarderMetaData.ABI

var KeystoneForwarderBin = KeystoneForwarderMetaData.Bin

func DeployKeystoneForwarder(auth *bind.TransactOpts, backend bind.ContractBackend, router common.Address) (common.Address, *types.Transaction, *KeystoneForwarder, error) {
	parsed, err := KeystoneForwarderMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KeystoneForwarderBin), backend, router)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KeystoneForwarder{address: address, abi: *parsed, KeystoneForwarderCaller: KeystoneForwarderCaller{contract: contract}, KeystoneForwarderTransactor: KeystoneForwarderTransactor{contract: contract}, KeystoneForwarderFilterer: KeystoneForwarderFilterer{contract: contract}}, nil
}

type KeystoneForwarder struct {
	address common.Address
	abi     abi.ABI
	KeystoneForwarderCaller
	KeystoneForwarderTransactor
	KeystoneForwarderFilterer
}

type KeystoneForwarderCaller struct {
	contract *bind.BoundContract
}

type KeystoneForwarderTransactor struct {
	contract *bind.BoundContract
}

type KeystoneForwarderFilterer struct {
	contract *bind.BoundContract
}

type KeystoneForwarderSession struct {
	Contract     *KeystoneForwarder
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type KeystoneForwarderCallerSession struct {
	Contract *KeystoneForwarderCaller
	CallOpts bind.CallOpts
}

type KeystoneForwarderTransactorSession struct {
	Contract     *KeystoneForwarderTransactor
	TransactOpts bind.TransactOpts
}

type KeystoneForwarderRaw struct {
	Contract *KeystoneForwarder
}

type KeystoneForwarderCallerRaw struct {
	Contract *KeystoneForwarderCaller
}

type KeystoneForwarderTransactorRaw struct {
	Contract *KeystoneForwarderTransactor
}

func NewKeystoneForwarder(address common.Address, backend bind.ContractBackend) (*KeystoneForwarder, error) {
	abi, err := abi.JSON(strings.NewReader(KeystoneForwarderABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindKeystoneForwarder(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeystoneForwarder{address: address, abi: abi, KeystoneForwarderCaller: KeystoneForwarderCaller{contract: contract}, KeystoneForwarderTransactor: KeystoneForwarderTransactor{contract: contract}, KeystoneForwarderFilterer: KeystoneForwarderFilterer{contract: contract}}, nil
}

func NewKeystoneForwarderCaller(address common.Address, caller bind.ContractCaller) (*KeystoneForwarderCaller, error) {
	contract, err := bindKeystoneForwarder(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeystoneForwarderCaller{contract: contract}, nil
}

func NewKeystoneForwarderTransactor(address common.Address, transactor bind.ContractTransactor) (*KeystoneForwarderTransactor, error) {
	contract, err := bindKeystoneForwarder(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeystoneForwarderTransactor{contract: contract}, nil
}

func NewKeystoneForwarderFilterer(address common.Address, filterer bind.ContractFilterer) (*KeystoneForwarderFilterer, error) {
	contract, err := bindKeystoneForwarder(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeystoneForwarderFilterer{contract: contract}, nil
}

func bindKeystoneForwarder(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := KeystoneForwarderMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_KeystoneForwarder *KeystoneForwarderRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeystoneForwarder.Contract.KeystoneForwarderCaller.contract.Call(opts, result, method, params...)
}

func (_KeystoneForwarder *KeystoneForwarderRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.KeystoneForwarderTransactor.contract.Transfer(opts)
}

func (_KeystoneForwarder *KeystoneForwarderRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.KeystoneForwarderTransactor.contract.Transact(opts, method, params...)
}

func (_KeystoneForwarder *KeystoneForwarderCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeystoneForwarder.Contract.contract.Call(opts, result, method, params...)
}

func (_KeystoneForwarder *KeystoneForwarderTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.contract.Transfer(opts)
}

func (_KeystoneForwarder *KeystoneForwarderTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.contract.Transact(opts, method, params...)
}

func (_KeystoneForwarder *KeystoneForwarderCaller) GetRouter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeystoneForwarder.contract.Call(opts, &out, "getRouter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeystoneForwarder *KeystoneForwarderSession) GetRouter() (common.Address, error) {
	return _KeystoneForwarder.Contract.GetRouter(&_KeystoneForwarder.CallOpts)
}

func (_KeystoneForwarder *KeystoneForwarderCallerSession) GetRouter() (common.Address, error) {
	return _KeystoneForwarder.Contract.GetRouter(&_KeystoneForwarder.CallOpts)
}

func (_KeystoneForwarder *KeystoneForwarderCaller) GetTransmissionId(opts *bind.CallOpts, receiverAddress common.Address, workflowExecutionId [32]byte, reportId [2]byte) ([32]byte, error) {
	var out []interface{}
	err := _KeystoneForwarder.contract.Call(opts, &out, "getTransmissionId", receiverAddress, workflowExecutionId, reportId)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_KeystoneForwarder *KeystoneForwarderSession) GetTransmissionId(receiverAddress common.Address, workflowExecutionId [32]byte, reportId [2]byte) ([32]byte, error) {
	return _KeystoneForwarder.Contract.GetTransmissionId(&_KeystoneForwarder.CallOpts, receiverAddress, workflowExecutionId, reportId)
}

func (_KeystoneForwarder *KeystoneForwarderCallerSession) GetTransmissionId(receiverAddress common.Address, workflowExecutionId [32]byte, reportId [2]byte) ([32]byte, error) {
	return _KeystoneForwarder.Contract.GetTransmissionId(&_KeystoneForwarder.CallOpts, receiverAddress, workflowExecutionId, reportId)
}

func (_KeystoneForwarder *KeystoneForwarderCaller) GetTransmissionState(opts *bind.CallOpts, receiverAddress common.Address, workflowExecutionId [32]byte, reportId [2]byte) (uint8, error) {
	var out []interface{}
	err := _KeystoneForwarder.contract.Call(opts, &out, "getTransmissionState", receiverAddress, workflowExecutionId, reportId)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_KeystoneForwarder *KeystoneForwarderSession) GetTransmissionState(receiverAddress common.Address, workflowExecutionId [32]byte, reportId [2]byte) (uint8, error) {
	return _KeystoneForwarder.Contract.GetTransmissionState(&_KeystoneForwarder.CallOpts, receiverAddress, workflowExecutionId, reportId)
}

func (_KeystoneForwarder *KeystoneForwarderCallerSession) GetTransmissionState(receiverAddress common.Address, workflowExecutionId [32]byte, reportId [2]byte) (uint8, error) {
	return _KeystoneForwarder.Contract.GetTransmissionState(&_KeystoneForwarder.CallOpts, receiverAddress, workflowExecutionId, reportId)
}

func (_KeystoneForwarder *KeystoneForwarderCaller) GetTransmitter(opts *bind.CallOpts, receiverAddress common.Address, workflowExecutionId [32]byte, reportId [2]byte) (common.Address, error) {
	var out []interface{}
	err := _KeystoneForwarder.contract.Call(opts, &out, "getTransmitter", receiverAddress, workflowExecutionId, reportId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeystoneForwarder *KeystoneForwarderSession) GetTransmitter(receiverAddress common.Address, workflowExecutionId [32]byte, reportId [2]byte) (common.Address, error) {
	return _KeystoneForwarder.Contract.GetTransmitter(&_KeystoneForwarder.CallOpts, receiverAddress, workflowExecutionId, reportId)
}

func (_KeystoneForwarder *KeystoneForwarderCallerSession) GetTransmitter(receiverAddress common.Address, workflowExecutionId [32]byte, reportId [2]byte) (common.Address, error) {
	return _KeystoneForwarder.Contract.GetTransmitter(&_KeystoneForwarder.CallOpts, receiverAddress, workflowExecutionId, reportId)
}

func (_KeystoneForwarder *KeystoneForwarderCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeystoneForwarder.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeystoneForwarder *KeystoneForwarderSession) Owner() (common.Address, error) {
	return _KeystoneForwarder.Contract.Owner(&_KeystoneForwarder.CallOpts)
}

func (_KeystoneForwarder *KeystoneForwarderCallerSession) Owner() (common.Address, error) {
	return _KeystoneForwarder.Contract.Owner(&_KeystoneForwarder.CallOpts)
}

func (_KeystoneForwarder *KeystoneForwarderCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _KeystoneForwarder.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_KeystoneForwarder *KeystoneForwarderSession) TypeAndVersion() (string, error) {
	return _KeystoneForwarder.Contract.TypeAndVersion(&_KeystoneForwarder.CallOpts)
}

func (_KeystoneForwarder *KeystoneForwarderCallerSession) TypeAndVersion() (string, error) {
	return _KeystoneForwarder.Contract.TypeAndVersion(&_KeystoneForwarder.CallOpts)
}

func (_KeystoneForwarder *KeystoneForwarderTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeystoneForwarder.contract.Transact(opts, "acceptOwnership")
}

func (_KeystoneForwarder *KeystoneForwarderSession) AcceptOwnership() (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.AcceptOwnership(&_KeystoneForwarder.TransactOpts)
}

func (_KeystoneForwarder *KeystoneForwarderTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.AcceptOwnership(&_KeystoneForwarder.TransactOpts)
}

func (_KeystoneForwarder *KeystoneForwarderTransactor) ClearConfig(opts *bind.TransactOpts, donId uint32, configVersion uint32) (*types.Transaction, error) {
	return _KeystoneForwarder.contract.Transact(opts, "clearConfig", donId, configVersion)
}

func (_KeystoneForwarder *KeystoneForwarderSession) ClearConfig(donId uint32, configVersion uint32) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.ClearConfig(&_KeystoneForwarder.TransactOpts, donId, configVersion)
}

func (_KeystoneForwarder *KeystoneForwarderTransactorSession) ClearConfig(donId uint32, configVersion uint32) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.ClearConfig(&_KeystoneForwarder.TransactOpts, donId, configVersion)
}

func (_KeystoneForwarder *KeystoneForwarderTransactor) Report(opts *bind.TransactOpts, receiverAddress common.Address, rawReport []byte, reportContext []byte, signatures [][]byte) (*types.Transaction, error) {
	return _KeystoneForwarder.contract.Transact(opts, "report", receiverAddress, rawReport, reportContext, signatures)
}

func (_KeystoneForwarder *KeystoneForwarderSession) Report(receiverAddress common.Address, rawReport []byte, reportContext []byte, signatures [][]byte) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.Report(&_KeystoneForwarder.TransactOpts, receiverAddress, rawReport, reportContext, signatures)
}

func (_KeystoneForwarder *KeystoneForwarderTransactorSession) Report(receiverAddress common.Address, rawReport []byte, reportContext []byte, signatures [][]byte) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.Report(&_KeystoneForwarder.TransactOpts, receiverAddress, rawReport, reportContext, signatures)
}

func (_KeystoneForwarder *KeystoneForwarderTransactor) SetConfig(opts *bind.TransactOpts, donId uint32, configVersion uint32, f uint8, signers []common.Address) (*types.Transaction, error) {
	return _KeystoneForwarder.contract.Transact(opts, "setConfig", donId, configVersion, f, signers)
}

func (_KeystoneForwarder *KeystoneForwarderSession) SetConfig(donId uint32, configVersion uint32, f uint8, signers []common.Address) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.SetConfig(&_KeystoneForwarder.TransactOpts, donId, configVersion, f, signers)
}

func (_KeystoneForwarder *KeystoneForwarderTransactorSession) SetConfig(donId uint32, configVersion uint32, f uint8, signers []common.Address) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.SetConfig(&_KeystoneForwarder.TransactOpts, donId, configVersion, f, signers)
}

func (_KeystoneForwarder *KeystoneForwarderTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _KeystoneForwarder.contract.Transact(opts, "transferOwnership", to)
}

func (_KeystoneForwarder *KeystoneForwarderSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.TransferOwnership(&_KeystoneForwarder.TransactOpts, to)
}

func (_KeystoneForwarder *KeystoneForwarderTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.TransferOwnership(&_KeystoneForwarder.TransactOpts, to)
}

type KeystoneForwarderConfigSetIterator struct {
	Event *KeystoneForwarderConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeystoneForwarderConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeystoneForwarderConfigSet)
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
		it.Event = new(KeystoneForwarderConfigSet)
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

func (it *KeystoneForwarderConfigSetIterator) Error() error {
	return it.fail
}

func (it *KeystoneForwarderConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeystoneForwarderConfigSet struct {
	DonId         uint32
	ConfigVersion uint32
	F             uint8
	Signers       []common.Address
	Raw           types.Log
}

func (_KeystoneForwarder *KeystoneForwarderFilterer) FilterConfigSet(opts *bind.FilterOpts, donId []uint32, configVersion []uint32) (*KeystoneForwarderConfigSetIterator, error) {

	var donIdRule []interface{}
	for _, donIdItem := range donId {
		donIdRule = append(donIdRule, donIdItem)
	}
	var configVersionRule []interface{}
	for _, configVersionItem := range configVersion {
		configVersionRule = append(configVersionRule, configVersionItem)
	}

	logs, sub, err := _KeystoneForwarder.contract.FilterLogs(opts, "ConfigSet", donIdRule, configVersionRule)
	if err != nil {
		return nil, err
	}
	return &KeystoneForwarderConfigSetIterator{contract: _KeystoneForwarder.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_KeystoneForwarder *KeystoneForwarderFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *KeystoneForwarderConfigSet, donId []uint32, configVersion []uint32) (event.Subscription, error) {

	var donIdRule []interface{}
	for _, donIdItem := range donId {
		donIdRule = append(donIdRule, donIdItem)
	}
	var configVersionRule []interface{}
	for _, configVersionItem := range configVersion {
		configVersionRule = append(configVersionRule, configVersionItem)
	}

	logs, sub, err := _KeystoneForwarder.contract.WatchLogs(opts, "ConfigSet", donIdRule, configVersionRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeystoneForwarderConfigSet)
				if err := _KeystoneForwarder.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_KeystoneForwarder *KeystoneForwarderFilterer) ParseConfigSet(log types.Log) (*KeystoneForwarderConfigSet, error) {
	event := new(KeystoneForwarderConfigSet)
	if err := _KeystoneForwarder.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeystoneForwarderOwnershipTransferRequestedIterator struct {
	Event *KeystoneForwarderOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeystoneForwarderOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeystoneForwarderOwnershipTransferRequested)
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
		it.Event = new(KeystoneForwarderOwnershipTransferRequested)
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

func (it *KeystoneForwarderOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeystoneForwarderOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeystoneForwarderOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeystoneForwarder *KeystoneForwarderFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeystoneForwarderOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeystoneForwarder.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeystoneForwarderOwnershipTransferRequestedIterator{contract: _KeystoneForwarder.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeystoneForwarder *KeystoneForwarderFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeystoneForwarderOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeystoneForwarder.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeystoneForwarderOwnershipTransferRequested)
				if err := _KeystoneForwarder.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_KeystoneForwarder *KeystoneForwarderFilterer) ParseOwnershipTransferRequested(log types.Log) (*KeystoneForwarderOwnershipTransferRequested, error) {
	event := new(KeystoneForwarderOwnershipTransferRequested)
	if err := _KeystoneForwarder.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeystoneForwarderOwnershipTransferredIterator struct {
	Event *KeystoneForwarderOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeystoneForwarderOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeystoneForwarderOwnershipTransferred)
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
		it.Event = new(KeystoneForwarderOwnershipTransferred)
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

func (it *KeystoneForwarderOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *KeystoneForwarderOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeystoneForwarderOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeystoneForwarder *KeystoneForwarderFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeystoneForwarderOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeystoneForwarder.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeystoneForwarderOwnershipTransferredIterator{contract: _KeystoneForwarder.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_KeystoneForwarder *KeystoneForwarderFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeystoneForwarderOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeystoneForwarder.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeystoneForwarderOwnershipTransferred)
				if err := _KeystoneForwarder.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_KeystoneForwarder *KeystoneForwarderFilterer) ParseOwnershipTransferred(log types.Log) (*KeystoneForwarderOwnershipTransferred, error) {
	event := new(KeystoneForwarderOwnershipTransferred)
	if err := _KeystoneForwarder.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeystoneForwarderReportProcessedIterator struct {
	Event *KeystoneForwarderReportProcessed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeystoneForwarderReportProcessedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeystoneForwarderReportProcessed)
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
		it.Event = new(KeystoneForwarderReportProcessed)
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

func (it *KeystoneForwarderReportProcessedIterator) Error() error {
	return it.fail
}

func (it *KeystoneForwarderReportProcessedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeystoneForwarderReportProcessed struct {
	Receiver            common.Address
	WorkflowExecutionId [32]byte
	ReportId            [2]byte
	Result              bool
	Raw                 types.Log
}

func (_KeystoneForwarder *KeystoneForwarderFilterer) FilterReportProcessed(opts *bind.FilterOpts, receiver []common.Address, workflowExecutionId [][32]byte, reportId [][2]byte) (*KeystoneForwarderReportProcessedIterator, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}
	var workflowExecutionIdRule []interface{}
	for _, workflowExecutionIdItem := range workflowExecutionId {
		workflowExecutionIdRule = append(workflowExecutionIdRule, workflowExecutionIdItem)
	}
	var reportIdRule []interface{}
	for _, reportIdItem := range reportId {
		reportIdRule = append(reportIdRule, reportIdItem)
	}

	logs, sub, err := _KeystoneForwarder.contract.FilterLogs(opts, "ReportProcessed", receiverRule, workflowExecutionIdRule, reportIdRule)
	if err != nil {
		return nil, err
	}
	return &KeystoneForwarderReportProcessedIterator{contract: _KeystoneForwarder.contract, event: "ReportProcessed", logs: logs, sub: sub}, nil
}

func (_KeystoneForwarder *KeystoneForwarderFilterer) WatchReportProcessed(opts *bind.WatchOpts, sink chan<- *KeystoneForwarderReportProcessed, receiver []common.Address, workflowExecutionId [][32]byte, reportId [][2]byte) (event.Subscription, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}
	var workflowExecutionIdRule []interface{}
	for _, workflowExecutionIdItem := range workflowExecutionId {
		workflowExecutionIdRule = append(workflowExecutionIdRule, workflowExecutionIdItem)
	}
	var reportIdRule []interface{}
	for _, reportIdItem := range reportId {
		reportIdRule = append(reportIdRule, reportIdItem)
	}

	logs, sub, err := _KeystoneForwarder.contract.WatchLogs(opts, "ReportProcessed", receiverRule, workflowExecutionIdRule, reportIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeystoneForwarderReportProcessed)
				if err := _KeystoneForwarder.contract.UnpackLog(event, "ReportProcessed", log); err != nil {
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

func (_KeystoneForwarder *KeystoneForwarderFilterer) ParseReportProcessed(log types.Log) (*KeystoneForwarderReportProcessed, error) {
	event := new(KeystoneForwarderReportProcessed)
	if err := _KeystoneForwarder.contract.UnpackLog(event, "ReportProcessed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_KeystoneForwarder *KeystoneForwarder) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _KeystoneForwarder.abi.Events["ConfigSet"].ID:
		return _KeystoneForwarder.ParseConfigSet(log)
	case _KeystoneForwarder.abi.Events["OwnershipTransferRequested"].ID:
		return _KeystoneForwarder.ParseOwnershipTransferRequested(log)
	case _KeystoneForwarder.abi.Events["OwnershipTransferred"].ID:
		return _KeystoneForwarder.ParseOwnershipTransferred(log)
	case _KeystoneForwarder.abi.Events["ReportProcessed"].ID:
		return _KeystoneForwarder.ParseReportProcessed(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (KeystoneForwarderConfigSet) Topic() common.Hash {
	return common.HexToHash("0x4120bd3b23957dd423555817d55654d4481b438aa15485c21b4180c784f1a455")
}

func (KeystoneForwarderOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (KeystoneForwarderOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (KeystoneForwarderReportProcessed) Topic() common.Hash {
	return common.HexToHash("0x3617b009e9785c42daebadb6d3fb553243a4bf586d07ea72d65d80013ce116b5")
}

func (_KeystoneForwarder *KeystoneForwarder) Address() common.Address {
	return _KeystoneForwarder.address
}

type KeystoneForwarderInterface interface {
	GetRouter(opts *bind.CallOpts) (common.Address, error)

	GetTransmissionId(opts *bind.CallOpts, receiverAddress common.Address, workflowExecutionId [32]byte, reportId [2]byte) ([32]byte, error)

	GetTransmissionState(opts *bind.CallOpts, receiverAddress common.Address, workflowExecutionId [32]byte, reportId [2]byte) (uint8, error)

	GetTransmitter(opts *bind.CallOpts, receiverAddress common.Address, workflowExecutionId [32]byte, reportId [2]byte) (common.Address, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	ClearConfig(opts *bind.TransactOpts, donId uint32, configVersion uint32) (*types.Transaction, error)

	Report(opts *bind.TransactOpts, receiverAddress common.Address, rawReport []byte, reportContext []byte, signatures [][]byte) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, donId uint32, configVersion uint32, f uint8, signers []common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterConfigSet(opts *bind.FilterOpts, donId []uint32, configVersion []uint32) (*KeystoneForwarderConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *KeystoneForwarderConfigSet, donId []uint32, configVersion []uint32) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*KeystoneForwarderConfigSet, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeystoneForwarderOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeystoneForwarderOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*KeystoneForwarderOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeystoneForwarderOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeystoneForwarderOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*KeystoneForwarderOwnershipTransferred, error)

	FilterReportProcessed(opts *bind.FilterOpts, receiver []common.Address, workflowExecutionId [][32]byte, reportId [][2]byte) (*KeystoneForwarderReportProcessedIterator, error)

	WatchReportProcessed(opts *bind.WatchOpts, sink chan<- *KeystoneForwarderReportProcessed, receiver []common.Address, workflowExecutionId [][32]byte, reportId [][2]byte) (event.Subscription, error)

	ParseReportProcessed(log types.Log) (*KeystoneForwarderReportProcessed, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
