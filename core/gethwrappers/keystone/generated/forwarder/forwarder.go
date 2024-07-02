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
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"transmissionId\",\"type\":\"bytes32\"}],\"name\":\"AlreadyAttempted\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"}],\"name\":\"DuplicateSigner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxSigners\",\"type\":\"uint256\"}],\"name\":\"ExcessSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FaultToleranceMustBePositive\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minSigners\",\"type\":\"uint256\"}],\"name\":\"InsufficientSigners\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"configId\",\"type\":\"uint64\"}],\"name\":\"InvalidConfig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReport\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"InvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"received\",\"type\":\"uint256\"}],\"name\":\"InvalidSignatureCount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"}],\"name\":\"InvalidSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedForwarder\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"configVersion\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"}],\"name\":\"ForwarderAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"}],\"name\":\"ForwarderRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"workflowExecutionId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes2\",\"name\":\"reportId\",\"type\":\"bytes2\"}],\"name\":\"Log2\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"workflowExecutionId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes2\",\"name\":\"reportId\",\"type\":\"bytes2\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"result\",\"type\":\"bool\"}],\"name\":\"ReportProcessed\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"}],\"name\":\"addForwarder\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"configVersion\",\"type\":\"uint32\"}],\"name\":\"clearConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"workflowExecutionId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes2\",\"name\":\"reportId\",\"type\":\"bytes2\"}],\"name\":\"getTransmissionId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"workflowExecutionId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes2\",\"name\":\"reportId\",\"type\":\"bytes2\"}],\"name\":\"getTransmissionState\",\"outputs\":[{\"internalType\":\"enumIRouter.TransmissionState\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"workflowExecutionId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes2\",\"name\":\"reportId\",\"type\":\"bytes2\"}],\"name\":\"getTransmitter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"}],\"name\":\"isForwarder\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"}],\"name\":\"removeForwarder\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"rawReport\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"reportContext\",\"type\":\"bytes\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"report\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"transmissionId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"validatedReport\",\"type\":\"bytes\"}],\"name\":\"route\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donId\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"configVersion\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5033806000816100675760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b038481169190911790915581161561009757610097816100b9565b5050306000908152600360205260409020805460ff1916600117905550610162565b336001600160a01b038216036101115760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161005e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b611535806101716000396000f3fe608060405234801561001057600080fd5b50600436106100ea5760003560e01c806379ba50971161008c578063abcef55411610066578063abcef5541461023e578063ee59d26c14610277578063ef6e17a01461028a578063f2fde38b1461029d57600080fd5b806379ba5097146101e05780638864b864146101e85780638da5cb5b1461022057600080fd5b8063354bdd66116100c8578063354bdd661461017957806343c164671461019a5780634d93172d146101ba5780635c41d2fe146101cd57600080fd5b806311289565146100ef578063181f5a7714610104578063233fd52d14610156575b600080fd5b6101026100fd366004611010565b6102b0565b005b6101406040518060400160405280601a81526020017f466f7277617264657220616e6420526f7574657220312e302e3000000000000081525081565b60405161014d91906110bb565b60405180910390f35b610169610164366004611128565b6103c3565b604051901515815260200161014d565b61018c6101873660046111b0565b6105bb565b60405190815260200161014d565b6101ad6101a83660046111b0565b61063f565b60405161014d9190611215565b6101026101c8366004611256565b6106c4565b6101026101db366004611256565b610740565b6101026107bf565b6101fb6101f63660046111b0565b6108bc565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161014d565b60005473ffffffffffffffffffffffffffffffffffffffff166101fb565b61016961024c366004611256565b73ffffffffffffffffffffffffffffffffffffffff1660009081526003602052604090205460ff1690565b610102610285366004611285565b6108fc565b610102610298366004611303565b610c75565b6101026102ab366004611256565b610d15565b606d8510156102eb576040517fb55ac75400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080600061032f89898080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250610d2992505050565b67ffffffffffffffff82166000818152600260209081526040918290208251878152918201939093527fffff000000000000000000000000000000000000000000000000000000000000841681830152905194975091955091935090917fa0ef515c75e6fba31047e81f41e37e079b3bbdebec13915c30110fcad6b4fdec9181900360600190a15050505050505050505050565b3360009081526003602052604081205460ff1661040c576040517fd79e123d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008881526004602052604090205473ffffffffffffffffffffffffffffffffffffffff1615610470576040517fa53dc8ca000000000000000000000000000000000000000000000000000000008152600481018990526024015b60405180910390fd5b600088815260046020526040812080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8a81169190911790915587163b90036104d2575060006105b0565b6040517f805f213200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff87169063805f21329061052a90889088908890889060040161137f565b600060405180830381600087803b15801561054457600080fd5b505af1925050508015610555575060015b610561575060006105b0565b50600087815260046020526040902080547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000017905560015b979650505050505050565b6040517fffffffffffffffffffffffffffffffffffffffff000000000000000000000000606085901b166020820152603481018390527fffff000000000000000000000000000000000000000000000000000000000000821660548201526000906056016040516020818303038152906040528051906020012090505b9392505050565b60008061064d8585856105bb565b60008181526004602052604090205490915073ffffffffffffffffffffffffffffffffffffffff16610683576000915050610638565b60008181526004602052604090205474010000000000000000000000000000000000000000900460ff166106b85760026106bb565b60015b95945050505050565b6106cc610d44565b73ffffffffffffffffffffffffffffffffffffffff811660008181526003602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055517fb96d15bf9258c7b8df062753a6a262864611fc7b060a5ee2e57e79b85f898d389190a250565b610748610d44565b73ffffffffffffffffffffffffffffffffffffffff811660008181526003602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055517f0ea0ce2c048ff45a4a95f2947879de3fb94abec2f152190400cab2d1272a68e79190a250565b60015473ffffffffffffffffffffffffffffffffffffffff163314610840576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610467565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6000600460006108cd8686866105bb565b815260208101919091526040016000205473ffffffffffffffffffffffffffffffffffffffff16949350505050565b610904610d44565b8260ff16600003610941576040517f0743bae600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601f811115610986576040517f61750f4000000000000000000000000000000000000000000000000000000000815260048101829052601f6024820152604401610467565b6109918360036113d5565b60ff1681116109ef57806109a68460036113d5565b6109b19060016113f8565b6040517f9dd9e6d8000000000000000000000000000000000000000000000000000000008152600481019290925260ff166024820152604401610467565b67ffffffff00000000602086901b1663ffffffff85161760005b67ffffffffffffffff8216600090815260026020526040902060010154811015610a9f5767ffffffffffffffff8216600090815260026020819052604082206001810180549190920192919084908110610a6557610a65611417565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff168352820192909252604001812055600101610a09565b5060005b82811015610bb7576000848483818110610abf57610abf611417565b9050602002016020810190610ad49190611256565b67ffffffffffffffff8416600090815260026020818152604080842073ffffffffffffffffffffffffffffffffffffffff86168552909201905290205490915015610b63576040517fe021c4f200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82166004820152602401610467565b610b6e826001611446565b67ffffffffffffffff8416600090815260026020818152604080842073ffffffffffffffffffffffffffffffffffffffff90961684529490910190529190912055600101610aa3565b5067ffffffffffffffff81166000908152600260205260409020610bdf906001018484610ebc565b5067ffffffffffffffff81166000908152600260205260409081902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff87161790555163ffffffff86811691908816907f4120bd3b23957dd423555817d55654d4481b438aa15485c21b4180c784f1a45590610c6590889088908890611459565b60405180910390a3505050505050565b610c7d610d44565b63ffffffff818116602084811b67ffffffff00000000168217600090815260028252604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690558051828152928301905291928516917f4120bd3b23957dd423555817d55654d4481b438aa15485c21b4180c784f1a45591604051610d099291906114bf565b60405180910390a35050565b610d1d610d44565b610d2681610dc7565b50565b60218101516045820151608b90920151909260c09290921c91565b60005473ffffffffffffffffffffffffffffffffffffffff163314610dc5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610467565b565b3373ffffffffffffffffffffffffffffffffffffffff821603610e46576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610467565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b828054828255906000526020600020908101928215610f34579160200282015b82811115610f345781547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff843516178255602090920191600190910190610edc565b50610f40929150610f44565b5090565b5b80821115610f405760008155600101610f45565b803573ffffffffffffffffffffffffffffffffffffffff81168114610f7d57600080fd5b919050565b60008083601f840112610f9457600080fd5b50813567ffffffffffffffff811115610fac57600080fd5b602083019150836020828501011115610fc457600080fd5b9250929050565b60008083601f840112610fdd57600080fd5b50813567ffffffffffffffff811115610ff557600080fd5b6020830191508360208260051b8501011115610fc457600080fd5b60008060008060008060006080888a03121561102b57600080fd5b61103488610f59565b9650602088013567ffffffffffffffff8082111561105157600080fd5b61105d8b838c01610f82565b909850965060408a013591508082111561107657600080fd5b6110828b838c01610f82565b909650945060608a013591508082111561109b57600080fd5b506110a88a828b01610fcb565b989b979a50959850939692959293505050565b60006020808352835180602085015260005b818110156110e9578581018301518582016040015282016110cd565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b600080600080600080600060a0888a03121561114357600080fd5b8735965061115360208901610f59565b955061116160408901610f59565b9450606088013567ffffffffffffffff8082111561117e57600080fd5b61118a8b838c01610f82565b909650945060808a01359150808211156111a357600080fd5b506110a88a828b01610f82565b6000806000606084860312156111c557600080fd5b6111ce84610f59565b92506020840135915060408401357fffff0000000000000000000000000000000000000000000000000000000000008116811461120a57600080fd5b809150509250925092565b6020810160038310611250577f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b91905290565b60006020828403121561126857600080fd5b61063882610f59565b803563ffffffff81168114610f7d57600080fd5b60008060008060006080868803121561129d57600080fd5b6112a686611271565b94506112b460208701611271565b9350604086013560ff811681146112ca57600080fd5b9250606086013567ffffffffffffffff8111156112e657600080fd5b6112f288828901610fcb565b969995985093965092949392505050565b6000806040838503121561131657600080fd5b61131f83611271565b915061132d60208401611271565b90509250929050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b604081526000611393604083018688611336565b82810360208401526105b0818587611336565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60ff81811683821602908116908181146113f1576113f16113a6565b5092915050565b60ff8181168382160190811115611411576114116113a6565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b80820180821115611411576114116113a6565b60ff8416815260406020808301829052908201839052600090849060608401835b868110156114b35773ffffffffffffffffffffffffffffffffffffffff6114a085610f59565b168252928201929082019060010161147a565b50979650505050505050565b60006040820160ff8516835260206040602085015281855180845260608601915060208701935060005b8181101561151b57845173ffffffffffffffffffffffffffffffffffffffff16835293830193918301916001016114e9565b509097965050505050505056fea164736f6c6343000818000a",
}

var KeystoneForwarderABI = KeystoneForwarderMetaData.ABI

var KeystoneForwarderBin = KeystoneForwarderMetaData.Bin

func DeployKeystoneForwarder(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *KeystoneForwarder, error) {
	parsed, err := KeystoneForwarderMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KeystoneForwarderBin), backend)
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

func (_KeystoneForwarder *KeystoneForwarderCaller) IsForwarder(opts *bind.CallOpts, forwarder common.Address) (bool, error) {
	var out []interface{}
	err := _KeystoneForwarder.contract.Call(opts, &out, "isForwarder", forwarder)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_KeystoneForwarder *KeystoneForwarderSession) IsForwarder(forwarder common.Address) (bool, error) {
	return _KeystoneForwarder.Contract.IsForwarder(&_KeystoneForwarder.CallOpts, forwarder)
}

func (_KeystoneForwarder *KeystoneForwarderCallerSession) IsForwarder(forwarder common.Address) (bool, error) {
	return _KeystoneForwarder.Contract.IsForwarder(&_KeystoneForwarder.CallOpts, forwarder)
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

func (_KeystoneForwarder *KeystoneForwarderTransactor) AddForwarder(opts *bind.TransactOpts, forwarder common.Address) (*types.Transaction, error) {
	return _KeystoneForwarder.contract.Transact(opts, "addForwarder", forwarder)
}

func (_KeystoneForwarder *KeystoneForwarderSession) AddForwarder(forwarder common.Address) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.AddForwarder(&_KeystoneForwarder.TransactOpts, forwarder)
}

func (_KeystoneForwarder *KeystoneForwarderTransactorSession) AddForwarder(forwarder common.Address) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.AddForwarder(&_KeystoneForwarder.TransactOpts, forwarder)
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

func (_KeystoneForwarder *KeystoneForwarderTransactor) RemoveForwarder(opts *bind.TransactOpts, forwarder common.Address) (*types.Transaction, error) {
	return _KeystoneForwarder.contract.Transact(opts, "removeForwarder", forwarder)
}

func (_KeystoneForwarder *KeystoneForwarderSession) RemoveForwarder(forwarder common.Address) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.RemoveForwarder(&_KeystoneForwarder.TransactOpts, forwarder)
}

func (_KeystoneForwarder *KeystoneForwarderTransactorSession) RemoveForwarder(forwarder common.Address) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.RemoveForwarder(&_KeystoneForwarder.TransactOpts, forwarder)
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

func (_KeystoneForwarder *KeystoneForwarderTransactor) Route(opts *bind.TransactOpts, transmissionId [32]byte, transmitter common.Address, receiver common.Address, metadata []byte, validatedReport []byte) (*types.Transaction, error) {
	return _KeystoneForwarder.contract.Transact(opts, "route", transmissionId, transmitter, receiver, metadata, validatedReport)
}

func (_KeystoneForwarder *KeystoneForwarderSession) Route(transmissionId [32]byte, transmitter common.Address, receiver common.Address, metadata []byte, validatedReport []byte) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.Route(&_KeystoneForwarder.TransactOpts, transmissionId, transmitter, receiver, metadata, validatedReport)
}

func (_KeystoneForwarder *KeystoneForwarderTransactorSession) Route(transmissionId [32]byte, transmitter common.Address, receiver common.Address, metadata []byte, validatedReport []byte) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.Route(&_KeystoneForwarder.TransactOpts, transmissionId, transmitter, receiver, metadata, validatedReport)
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

type KeystoneForwarderForwarderAddedIterator struct {
	Event *KeystoneForwarderForwarderAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeystoneForwarderForwarderAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeystoneForwarderForwarderAdded)
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
		it.Event = new(KeystoneForwarderForwarderAdded)
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

func (it *KeystoneForwarderForwarderAddedIterator) Error() error {
	return it.fail
}

func (it *KeystoneForwarderForwarderAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeystoneForwarderForwarderAdded struct {
	Forwarder common.Address
	Raw       types.Log
}

func (_KeystoneForwarder *KeystoneForwarderFilterer) FilterForwarderAdded(opts *bind.FilterOpts, forwarder []common.Address) (*KeystoneForwarderForwarderAddedIterator, error) {

	var forwarderRule []interface{}
	for _, forwarderItem := range forwarder {
		forwarderRule = append(forwarderRule, forwarderItem)
	}

	logs, sub, err := _KeystoneForwarder.contract.FilterLogs(opts, "ForwarderAdded", forwarderRule)
	if err != nil {
		return nil, err
	}
	return &KeystoneForwarderForwarderAddedIterator{contract: _KeystoneForwarder.contract, event: "ForwarderAdded", logs: logs, sub: sub}, nil
}

func (_KeystoneForwarder *KeystoneForwarderFilterer) WatchForwarderAdded(opts *bind.WatchOpts, sink chan<- *KeystoneForwarderForwarderAdded, forwarder []common.Address) (event.Subscription, error) {

	var forwarderRule []interface{}
	for _, forwarderItem := range forwarder {
		forwarderRule = append(forwarderRule, forwarderItem)
	}

	logs, sub, err := _KeystoneForwarder.contract.WatchLogs(opts, "ForwarderAdded", forwarderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeystoneForwarderForwarderAdded)
				if err := _KeystoneForwarder.contract.UnpackLog(event, "ForwarderAdded", log); err != nil {
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

func (_KeystoneForwarder *KeystoneForwarderFilterer) ParseForwarderAdded(log types.Log) (*KeystoneForwarderForwarderAdded, error) {
	event := new(KeystoneForwarderForwarderAdded)
	if err := _KeystoneForwarder.contract.UnpackLog(event, "ForwarderAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeystoneForwarderForwarderRemovedIterator struct {
	Event *KeystoneForwarderForwarderRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeystoneForwarderForwarderRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeystoneForwarderForwarderRemoved)
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
		it.Event = new(KeystoneForwarderForwarderRemoved)
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

func (it *KeystoneForwarderForwarderRemovedIterator) Error() error {
	return it.fail
}

func (it *KeystoneForwarderForwarderRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeystoneForwarderForwarderRemoved struct {
	Forwarder common.Address
	Raw       types.Log
}

func (_KeystoneForwarder *KeystoneForwarderFilterer) FilterForwarderRemoved(opts *bind.FilterOpts, forwarder []common.Address) (*KeystoneForwarderForwarderRemovedIterator, error) {

	var forwarderRule []interface{}
	for _, forwarderItem := range forwarder {
		forwarderRule = append(forwarderRule, forwarderItem)
	}

	logs, sub, err := _KeystoneForwarder.contract.FilterLogs(opts, "ForwarderRemoved", forwarderRule)
	if err != nil {
		return nil, err
	}
	return &KeystoneForwarderForwarderRemovedIterator{contract: _KeystoneForwarder.contract, event: "ForwarderRemoved", logs: logs, sub: sub}, nil
}

func (_KeystoneForwarder *KeystoneForwarderFilterer) WatchForwarderRemoved(opts *bind.WatchOpts, sink chan<- *KeystoneForwarderForwarderRemoved, forwarder []common.Address) (event.Subscription, error) {

	var forwarderRule []interface{}
	for _, forwarderItem := range forwarder {
		forwarderRule = append(forwarderRule, forwarderItem)
	}

	logs, sub, err := _KeystoneForwarder.contract.WatchLogs(opts, "ForwarderRemoved", forwarderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeystoneForwarderForwarderRemoved)
				if err := _KeystoneForwarder.contract.UnpackLog(event, "ForwarderRemoved", log); err != nil {
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

func (_KeystoneForwarder *KeystoneForwarderFilterer) ParseForwarderRemoved(log types.Log) (*KeystoneForwarderForwarderRemoved, error) {
	event := new(KeystoneForwarderForwarderRemoved)
	if err := _KeystoneForwarder.contract.UnpackLog(event, "ForwarderRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeystoneForwarderLog2Iterator struct {
	Event *KeystoneForwarderLog2

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeystoneForwarderLog2Iterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeystoneForwarderLog2)
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
		it.Event = new(KeystoneForwarderLog2)
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

func (it *KeystoneForwarderLog2Iterator) Error() error {
	return it.fail
}

func (it *KeystoneForwarderLog2Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeystoneForwarderLog2 struct {
	WorkflowExecutionId [32]byte
	ConfigId            uint64
	ReportId            [2]byte
	Raw                 types.Log
}

func (_KeystoneForwarder *KeystoneForwarderFilterer) FilterLog2(opts *bind.FilterOpts) (*KeystoneForwarderLog2Iterator, error) {

	logs, sub, err := _KeystoneForwarder.contract.FilterLogs(opts, "Log2")
	if err != nil {
		return nil, err
	}
	return &KeystoneForwarderLog2Iterator{contract: _KeystoneForwarder.contract, event: "Log2", logs: logs, sub: sub}, nil
}

func (_KeystoneForwarder *KeystoneForwarderFilterer) WatchLog2(opts *bind.WatchOpts, sink chan<- *KeystoneForwarderLog2) (event.Subscription, error) {

	logs, sub, err := _KeystoneForwarder.contract.WatchLogs(opts, "Log2")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeystoneForwarderLog2)
				if err := _KeystoneForwarder.contract.UnpackLog(event, "Log2", log); err != nil {
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

func (_KeystoneForwarder *KeystoneForwarderFilterer) ParseLog2(log types.Log) (*KeystoneForwarderLog2, error) {
	event := new(KeystoneForwarderLog2)
	if err := _KeystoneForwarder.contract.UnpackLog(event, "Log2", log); err != nil {
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
	case _KeystoneForwarder.abi.Events["ForwarderAdded"].ID:
		return _KeystoneForwarder.ParseForwarderAdded(log)
	case _KeystoneForwarder.abi.Events["ForwarderRemoved"].ID:
		return _KeystoneForwarder.ParseForwarderRemoved(log)
	case _KeystoneForwarder.abi.Events["Log2"].ID:
		return _KeystoneForwarder.ParseLog2(log)
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

func (KeystoneForwarderForwarderAdded) Topic() common.Hash {
	return common.HexToHash("0x0ea0ce2c048ff45a4a95f2947879de3fb94abec2f152190400cab2d1272a68e7")
}

func (KeystoneForwarderForwarderRemoved) Topic() common.Hash {
	return common.HexToHash("0xb96d15bf9258c7b8df062753a6a262864611fc7b060a5ee2e57e79b85f898d38")
}

func (KeystoneForwarderLog2) Topic() common.Hash {
	return common.HexToHash("0xa0ef515c75e6fba31047e81f41e37e079b3bbdebec13915c30110fcad6b4fdec")
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
	GetTransmissionId(opts *bind.CallOpts, receiver common.Address, workflowExecutionId [32]byte, reportId [2]byte) ([32]byte, error)

	GetTransmissionState(opts *bind.CallOpts, receiver common.Address, workflowExecutionId [32]byte, reportId [2]byte) (uint8, error)

	GetTransmitter(opts *bind.CallOpts, receiver common.Address, workflowExecutionId [32]byte, reportId [2]byte) (common.Address, error)

	IsForwarder(opts *bind.CallOpts, forwarder common.Address) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AddForwarder(opts *bind.TransactOpts, forwarder common.Address) (*types.Transaction, error)

	ClearConfig(opts *bind.TransactOpts, donId uint32, configVersion uint32) (*types.Transaction, error)

	RemoveForwarder(opts *bind.TransactOpts, forwarder common.Address) (*types.Transaction, error)

	Report(opts *bind.TransactOpts, receiver common.Address, rawReport []byte, reportContext []byte, signatures [][]byte) (*types.Transaction, error)

	Route(opts *bind.TransactOpts, transmissionId [32]byte, transmitter common.Address, receiver common.Address, metadata []byte, validatedReport []byte) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, donId uint32, configVersion uint32, f uint8, signers []common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterConfigSet(opts *bind.FilterOpts, donId []uint32, configVersion []uint32) (*KeystoneForwarderConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *KeystoneForwarderConfigSet, donId []uint32, configVersion []uint32) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*KeystoneForwarderConfigSet, error)

	FilterForwarderAdded(opts *bind.FilterOpts, forwarder []common.Address) (*KeystoneForwarderForwarderAddedIterator, error)

	WatchForwarderAdded(opts *bind.WatchOpts, sink chan<- *KeystoneForwarderForwarderAdded, forwarder []common.Address) (event.Subscription, error)

	ParseForwarderAdded(log types.Log) (*KeystoneForwarderForwarderAdded, error)

	FilterForwarderRemoved(opts *bind.FilterOpts, forwarder []common.Address) (*KeystoneForwarderForwarderRemovedIterator, error)

	WatchForwarderRemoved(opts *bind.WatchOpts, sink chan<- *KeystoneForwarderForwarderRemoved, forwarder []common.Address) (event.Subscription, error)

	ParseForwarderRemoved(log types.Log) (*KeystoneForwarderForwarderRemoved, error)

	FilterLog2(opts *bind.FilterOpts) (*KeystoneForwarderLog2Iterator, error)

	WatchLog2(opts *bind.WatchOpts, sink chan<- *KeystoneForwarderLog2) (event.Subscription, error)

	ParseLog2(log types.Log) (*KeystoneForwarderLog2, error)

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
