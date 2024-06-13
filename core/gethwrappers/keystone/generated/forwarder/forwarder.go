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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"}],\"name\":\"DuplicateSigner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxSigners\",\"type\":\"uint256\"}],\"name\":\"ExcessSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FaultToleranceMustBePositive\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minSigners\",\"type\":\"uint256\"}],\"name\":\"InsufficientSigners\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"configId\",\"type\":\"uint64\"}],\"name\":\"InvalidConfig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReport\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"InvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"received\",\"type\":\"uint256\"}],\"name\":\"InvalidSignatureCount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"}],\"name\":\"InvalidSigner\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"configVersion\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"workflowExecutionId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes2\",\"name\":\"reportId\",\"type\":\"bytes2\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"result\",\"type\":\"bool\"}],\"name\":\"ReportProcessed\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"configVersion\",\"type\":\"uint32\"}],\"name\":\"clearConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRouter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"workflowExecutionId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes2\",\"name\":\"reportId\",\"type\":\"bytes2\"}],\"name\":\"getTransmissionId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"workflowExecutionId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes2\",\"name\":\"reportId\",\"type\":\"bytes2\"}],\"name\":\"getTransmissionState\",\"outputs\":[{\"internalType\":\"enumIRouter.TransmissionState\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"workflowExecutionId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes2\",\"name\":\"reportId\",\"type\":\"bytes2\"}],\"name\":\"getTransmitter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"rawReport\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"reportContext\",\"type\":\"bytes\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"report\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"configVersion\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801562000010575f80fd5b5060405162001a8138038062001a8183398101604081905262000033916200018f565b33805f81620000895760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b5f80546001600160a01b0319166001600160a01b0384811691909117909155811615620000bb57620000bb81620000e5565b5050600280546001600160a01b0319166001600160a01b03939093169290921790915550620001be565b336001600160a01b038216036200013f5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000080565b600180546001600160a01b0319166001600160a01b038381169182179092555f8054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b5f60208284031215620001a0575f80fd5b81516001600160a01b0381168114620001b7575f80fd5b9392505050565b6118b580620001cc5f395ff3fe608060405234801561000f575f80fd5b50600436106100c4575f3560e01c80638864b8641161007d578063ee59d26c11610058578063ee59d26c14610268578063ef6e17a01461027b578063f2fde38b1461028e575f80fd5b80638864b864146101f55780638da5cb5b1461022d578063b0f479a11461024a575f80fd5b8063354bdd66116100ad578063354bdd661461012f57806343c16467146101cd57806379ba5097146101ed575f80fd5b806311289565146100c8578063181f5a77146100dd575b5f80fd5b6100db6100d6366004611299565b6102a1565b005b6101196040518060400160405280601781526020017f4b657973746f6e65466f7277617264657220312e302e3000000000000000000081525081565b604051610126919061133e565b60405180910390f35b6101bf61013d3660046113a8565b6040517fffffffffffffffffffffffffffffffffffffffff000000000000000000000000606085901b166020820152603481018390527fffff000000000000000000000000000000000000000000000000000000000000821660548201525f906056016040516020818303038152906040528051906020012090509392505050565b604051908152602001610126565b6101e06101db3660046113a8565b61086f565b604051610126919061140b565b6100db610981565b6102086102033660046113a8565b610a7d565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610126565b5f5473ffffffffffffffffffffffffffffffffffffffff16610208565b60025473ffffffffffffffffffffffffffffffffffffffff16610208565b6100db610276366004611462565b610b87565b6100db6102893660046114da565b610ef5565b6100db61029c36600461150b565b610f94565b606d8510156102dc576040517fb55ac75400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f805f61031d89898080601f0160208091040260200160405190810160405280939291908181526020018383808284375f92019190915250610fa892505050565b67ffffffffffffffff82165f90815260036020526040812080549497509195509193509160ff169081900361038f576040517fdf3b81ea00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff841660048201526024015b60405180910390fd5b8561039b82600161155a565b60ff16146103ed576103ae81600161155a565b6040517fd6022e8e00000000000000000000000000000000000000000000000000000000815260ff909116600482015260248101879052604401610386565b5f8b8b6040516103fe929190611579565b604051908190038120610417918c908c90602001611588565b604051602081830303815290604052805190602001209050610437611139565b5f5b888110156106b057365f8b8b84818110610455576104556115a1565b905060200281019061046791906115ce565b9092509050604181146104aa5781816040517f2adfdc30000000000000000000000000000000000000000000000000000000008152600401610386929190611676565b5f600186848460408181106104c1576104c16115a1565b6104d392013560f81c9050601b61155a565b6104e060205f8789611689565b6104e9916116b0565b6104f760406020888a611689565b610500916116b0565b604080515f8152602081018083529590955260ff909316928401929092526060830152608082015260a0016020604051602081039080840390855afa15801561054b573d5f803e3d5ffd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0015173ffffffffffffffffffffffffffffffffffffffff81165f90815260028c016020529182205490935091508190036105f0576040517fbf18af4300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff83166004820152602401610386565b5f868260208110610603576106036115a1565b602002015173ffffffffffffffffffffffffffffffffffffffff161461066d576040517fe021c4f200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff83166004820152602401610386565b81868260208110610680576106806115a1565b73ffffffffffffffffffffffffffffffffffffffff90921660209290920201525050600190920191506104399050565b50506002545f945073ffffffffffffffffffffffffffffffffffffffff16925063233fd52d915061076190508c86866040517fffffffffffffffffffffffffffffffffffffffff000000000000000000000000606085901b166020820152603481018390527fffff000000000000000000000000000000000000000000000000000000000000821660548201525f906056016040516020818303038152906040528051906020012090509392505050565b338d8d8d602d90606d9261077793929190611689565b8f8f606d90809261078a93929190611689565b6040518863ffffffff1660e01b81526004016107ac97969594939291906116ec565b6020604051808303815f875af11580156107c8573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906107ec919061174c565b9050817dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916838b73ffffffffffffffffffffffffffffffffffffffff167f3617b009e9785c42daebadb6d3fb553243a4bf586d07ea72d65d80013ce116b58460405161085b911515815260200190565b60405180910390a450505050505050505050565b600254604080517fffffffffffffffffffffffffffffffffffffffff000000000000000000000000606087901b16602080830191909152603482018690527fffff0000000000000000000000000000000000000000000000000000000000008516605483015282518083036036018152605683019384905280519101207f516db40800000000000000000000000000000000000000000000000000000000909252605a8101919091525f9173ffffffffffffffffffffffffffffffffffffffff169063516db40890607a01602060405180830381865afa158015610955573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610979919061176b565b949350505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610a02576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610386565b5f8054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b600254604080517fffffffffffffffffffffffffffffffffffffffff000000000000000000000000606087901b16602080830191909152603482018690527fffff0000000000000000000000000000000000000000000000000000000000008516605483015282518083036036018152605683019384905280519101207fe6b7145800000000000000000000000000000000000000000000000000000000909252605a8101919091525f9173ffffffffffffffffffffffffffffffffffffffff169063e6b7145890607a01602060405180830381865afa158015610b63573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906109799190611789565b610b8f610fc3565b8260ff165f03610bcb576040517f0743bae600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601f811115610c10576040517f61750f4000000000000000000000000000000000000000000000000000000000815260048101829052601f6024820152604401610386565b610c1b8360036117a4565b60ff168111610c795780610c308460036117a4565b610c3b90600161155a565b6040517f9dd9e6d8000000000000000000000000000000000000000000000000000000008152600481019290925260ff166024820152604401610386565b67ffffffff00000000602086901b1663ffffffff8516175f5b67ffffffffffffffff82165f90815260036020526040902060010154811015610d235767ffffffffffffffff82165f9081526003602052604081206001810180546002909201929184908110610cea57610cea6115a1565b5f91825260208083209091015473ffffffffffffffffffffffffffffffffffffffff168352820192909252604001812055600101610c92565b505f5b82811015610e39575f848483818110610d4157610d416115a1565b9050602002016020810190610d56919061150b565b67ffffffffffffffff84165f90815260036020908152604080832073ffffffffffffffffffffffffffffffffffffffff8516845260020190915290205490915015610de5576040517fe021c4f200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82166004820152602401610386565b610df08260016117c7565b67ffffffffffffffff84165f90815260036020908152604080832073ffffffffffffffffffffffffffffffffffffffff9095168352600290940190529190912055600101610d26565b5067ffffffffffffffff81165f908152600360205260409020610e60906001018484611158565b5067ffffffffffffffff81165f908152600360205260409081902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff87161790555163ffffffff86811691908816907f4120bd3b23957dd423555817d55654d4481b438aa15485c21b4180c784f1a45590610ee5908890889088906117da565b60405180910390a3505050505050565b610efd610fc3565b63ffffffff818116602084811b67ffffffff000000001682175f90815260038252604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690558051828152928301905291928516917f4120bd3b23957dd423555817d55654d4481b438aa15485c21b4180c784f1a45591604051610f88929190611841565b60405180910390a35050565b610f9c610fc3565b610fa581611045565b50565b60218101516045820151608b90920151909260c09290921c91565b5f5473ffffffffffffffffffffffffffffffffffffffff163314611043576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610386565b565b3373ffffffffffffffffffffffffffffffffffffffff8216036110c4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610386565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8381169182179092555f8054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6040518061040001604052806020906020820280368337509192915050565b828054828255905f5260205f209081019282156111ce579160200282015b828111156111ce5781547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff843516178255602090920191600190910190611176565b506111da9291506111de565b5090565b5b808211156111da575f81556001016111df565b73ffffffffffffffffffffffffffffffffffffffff81168114610fa5575f80fd5b5f8083601f840112611223575f80fd5b50813567ffffffffffffffff81111561123a575f80fd5b602083019150836020828501011115611251575f80fd5b9250929050565b5f8083601f840112611268575f80fd5b50813567ffffffffffffffff81111561127f575f80fd5b6020830191508360208260051b8501011115611251575f80fd5b5f805f805f805f6080888a0312156112af575f80fd5b87356112ba816111f2565b9650602088013567ffffffffffffffff808211156112d6575f80fd5b6112e28b838c01611213565b909850965060408a01359150808211156112fa575f80fd5b6113068b838c01611213565b909650945060608a013591508082111561131e575f80fd5b5061132b8a828b01611258565b989b979a50959850939692959293505050565b5f602080835283518060208501525f5b8181101561136a5785810183015185820160400152820161134e565b505f6040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b5f805f606084860312156113ba575f80fd5b83356113c5816111f2565b92506020840135915060408401357fffff00000000000000000000000000000000000000000000000000000000000081168114611400575f80fd5b809150509250925092565b6020810160038310611444577f4e487b71000000000000000000000000000000000000000000000000000000005f52602160045260245ffd5b91905290565b803563ffffffff8116811461145d575f80fd5b919050565b5f805f805f60808688031215611476575f80fd5b61147f8661144a565b945061148d6020870161144a565b9350604086013560ff811681146114a2575f80fd5b9250606086013567ffffffffffffffff8111156114bd575f80fd5b6114c988828901611258565b969995985093965092949392505050565b5f80604083850312156114eb575f80fd5b6114f48361144a565b91506115026020840161144a565b90509250929050565b5f6020828403121561151b575f80fd5b8135611526816111f2565b9392505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b60ff81811683821601908111156115735761157361152d565b92915050565b818382375f9101908152919050565b838152818360208301375f910160200190815292915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b5f8083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112611601575f80fd5b83018035915067ffffffffffffffff82111561161b575f80fd5b602001915036819003821315611251575f80fd5b81835281816020850137505f602082840101525f60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b602081525f61097960208301848661162f565b5f8085851115611697575f80fd5b838611156116a3575f80fd5b5050820193919092039150565b80356020831015611573577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff602084900360031b1b1692915050565b8781525f73ffffffffffffffffffffffffffffffffffffffff808916602084015280881660408401525060a0606083015261172b60a08301868861162f565b828103608084015261173e81858761162f565b9a9950505050505050505050565b5f6020828403121561175c575f80fd5b81518015158114611526575f80fd5b5f6020828403121561177b575f80fd5b815160038110611526575f80fd5b5f60208284031215611799575f80fd5b8151611526816111f2565b60ff81811683821602908116908181146117c0576117c061152d565b5092915050565b808201808211156115735761157361152d565b60ff84168152604060208083018290529082018390525f90849060608401835b8681101561183557833561180d816111f2565b73ffffffffffffffffffffffffffffffffffffffff16825292820192908201906001016117fa565b50979650505050505050565b5f6040820160ff851683526020604060208501528185518084526060860191506020870193505f5b8181101561189b57845173ffffffffffffffffffffffffffffffffffffffff1683529383019391830191600101611869565b509097965050505050505056fea164736f6c6343000818000a",
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

func (_KeystoneForwarder *KeystoneForwarderCaller) GetTransmissionId(opts *bind.CallOpts, receiver common.Address, workflowExecutionId [32]byte, reportId [2]byte) ([32]byte, error) {
	var out []interface{}
	err := _KeystoneForwarder.contract.Call(opts, &out, "getTransmissionId", receiver, workflowExecutionId, reportId)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_KeystoneForwarder *KeystoneForwarderSession) GetTransmissionId(receiver common.Address, workflowExecutionId [32]byte, reportId [2]byte) ([32]byte, error) {
	return _KeystoneForwarder.Contract.GetTransmissionId(&_KeystoneForwarder.CallOpts, receiver, workflowExecutionId, reportId)
}

func (_KeystoneForwarder *KeystoneForwarderCallerSession) GetTransmissionId(receiver common.Address, workflowExecutionId [32]byte, reportId [2]byte) ([32]byte, error) {
	return _KeystoneForwarder.Contract.GetTransmissionId(&_KeystoneForwarder.CallOpts, receiver, workflowExecutionId, reportId)
}

func (_KeystoneForwarder *KeystoneForwarderCaller) GetTransmissionState(opts *bind.CallOpts, receiver common.Address, workflowExecutionId [32]byte, reportId [2]byte) (uint8, error) {
	var out []interface{}
	err := _KeystoneForwarder.contract.Call(opts, &out, "getTransmissionState", receiver, workflowExecutionId, reportId)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_KeystoneForwarder *KeystoneForwarderSession) GetTransmissionState(receiver common.Address, workflowExecutionId [32]byte, reportId [2]byte) (uint8, error) {
	return _KeystoneForwarder.Contract.GetTransmissionState(&_KeystoneForwarder.CallOpts, receiver, workflowExecutionId, reportId)
}

func (_KeystoneForwarder *KeystoneForwarderCallerSession) GetTransmissionState(receiver common.Address, workflowExecutionId [32]byte, reportId [2]byte) (uint8, error) {
	return _KeystoneForwarder.Contract.GetTransmissionState(&_KeystoneForwarder.CallOpts, receiver, workflowExecutionId, reportId)
}

func (_KeystoneForwarder *KeystoneForwarderCaller) GetTransmitter(opts *bind.CallOpts, receiver common.Address, workflowExecutionId [32]byte, reportId [2]byte) (common.Address, error) {
	var out []interface{}
	err := _KeystoneForwarder.contract.Call(opts, &out, "getTransmitter", receiver, workflowExecutionId, reportId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeystoneForwarder *KeystoneForwarderSession) GetTransmitter(receiver common.Address, workflowExecutionId [32]byte, reportId [2]byte) (common.Address, error) {
	return _KeystoneForwarder.Contract.GetTransmitter(&_KeystoneForwarder.CallOpts, receiver, workflowExecutionId, reportId)
}

func (_KeystoneForwarder *KeystoneForwarderCallerSession) GetTransmitter(receiver common.Address, workflowExecutionId [32]byte, reportId [2]byte) (common.Address, error) {
	return _KeystoneForwarder.Contract.GetTransmitter(&_KeystoneForwarder.CallOpts, receiver, workflowExecutionId, reportId)
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

func (_KeystoneForwarder *KeystoneForwarderTransactor) Report(opts *bind.TransactOpts, receiver common.Address, rawReport []byte, reportContext []byte, signatures [][]byte) (*types.Transaction, error) {
	return _KeystoneForwarder.contract.Transact(opts, "report", receiver, rawReport, reportContext, signatures)
}

func (_KeystoneForwarder *KeystoneForwarderSession) Report(receiver common.Address, rawReport []byte, reportContext []byte, signatures [][]byte) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.Report(&_KeystoneForwarder.TransactOpts, receiver, rawReport, reportContext, signatures)
}

func (_KeystoneForwarder *KeystoneForwarderTransactorSession) Report(receiver common.Address, rawReport []byte, reportContext []byte, signatures [][]byte) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.Report(&_KeystoneForwarder.TransactOpts, receiver, rawReport, reportContext, signatures)
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

	GetTransmissionId(opts *bind.CallOpts, receiver common.Address, workflowExecutionId [32]byte, reportId [2]byte) ([32]byte, error)

	GetTransmissionState(opts *bind.CallOpts, receiver common.Address, workflowExecutionId [32]byte, reportId [2]byte) (uint8, error)

	GetTransmitter(opts *bind.CallOpts, receiver common.Address, workflowExecutionId [32]byte, reportId [2]byte) (common.Address, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	ClearConfig(opts *bind.TransactOpts, donId uint32, configVersion uint32) (*types.Transaction, error)

	Report(opts *bind.TransactOpts, receiver common.Address, rawReport []byte, reportContext []byte, signatures [][]byte) (*types.Transaction, error)

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
