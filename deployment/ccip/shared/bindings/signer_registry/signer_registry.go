// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package signer_registry

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
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated"
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

type ISignerRegistrySigner struct {
	EvmAddress    common.Address
	NewEVMAddress common.Address
}

var SignerRegistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"maxSigners\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"newSigners\",\"type\":\"tuple[]\",\"internalType\":\"structISignerRegistry.Signer[]\",\"components\":[{\"name\":\"evmAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"newEVMAddress\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addSigners\",\"inputs\":[{\"name\":\"newSigners\",\"type\":\"tuple[]\",\"internalType\":\"structISignerRegistry.Signer[]\",\"components\":[{\"name\":\"evmAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"newEVMAddress\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getMaxSigners\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSignerCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSigners\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structISignerRegistry.Signer[]\",\"components\":[{\"name\":\"evmAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"newEVMAddress\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isSigner\",\"inputs\":[{\"name\":\"signerAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"promoteNewSignerAddresses\",\"inputs\":[{\"name\":\"existingSignerAddresses\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"removeSigners\",\"inputs\":[{\"name\":\"signersToRemove\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setNewSignerAddresses\",\"inputs\":[{\"name\":\"existingSignerAddresses\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"newSignerAddresses\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"typeAndVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"NewSignerAddressSet\",\"inputs\":[{\"name\":\"evmAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"previousNewEVMAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newEVMAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SignerAdded\",\"inputs\":[{\"name\":\"evmAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SignerAddressPromoted\",\"inputs\":[{\"name\":\"evmAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newEVMAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SignerNotRemovedDueToNoMatch\",\"inputs\":[{\"name\":\"signerAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SignerRemoved\",\"inputs\":[{\"name\":\"evmAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"CannotTransferToSelf\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DuplicateSigner\",\"inputs\":[{\"name\":\"signerAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InvalidInputLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSignerAddress\",\"inputs\":[{\"name\":\"signerAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"MissingNewAddress\",\"inputs\":[{\"name\":\"signerAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"MustBeProposedOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoMatchingSignerFound\",\"inputs\":[{\"name\":\"signerAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnerCannotBeZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TooManySigners\",\"inputs\":[{\"name\":\"currentSignerCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"newSignersCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]",
	Bin: "0x60a0604052346102f55761179980380380610019816102f9565b9283398101906040818303126102f5578051602082015190916001600160401b0382116102f5570182601f820112156102f5578051906001600160401b0382116101fc5761006c60208360051b016102f9565b9360208086858152019360061b830101918183116102f557602001925b8284106102a3578585331561029457600180546001600160a01b031916331790556080528051600254906100be908290610332565b6080511061027357505f5b8151811015610254576100dc8183610353565b5180516001600160a01b03168015610242575080516020820180519092916001600160a01b0391821691168114610230575080515f90610124906001600160a01b0316610397565b1215610210575080515f90610141906001600160a01b0316610397565b121561021057506101528183610353565b5190600254680100000000000000008110156101fc57806001610178920160025561037b565b6101e957825181546001600160a01b03199081166001600160a01b039283161783556020909401516001928301805490951690821617909355916101bc8285610353565b5151167f47d1c22a25bb3a5d4e481b9b1e6944c2eade3181a0a20b495ed61d35b5323f245f80a2016100c9565b634e487b7160e01b5f525f60045260245ffd5b634e487b7160e01b5f52604160045260245ffd5b51637010e27960e11b5f9081526001600160a01b03909116600452602490fd5b637010e27960e11b5f5260045260245ffd5b635a1d346d60e01b5f5260045260245ffd5b604051611389908161041082396080518181816101cf0152610cf10152f35b61027e818351610332565b9063118587f960e21b5f5260045260245260445ffd5b639b15e16f60e01b5f5260045ffd5b6040848303126102f5576040805191908201906001600160401b038211838310176101fc5760409260209284526102d98761031e565b81526102e683880161031e565b83820152815201930192610089565b5f80fd5b6040519190601f01601f191682016001600160401b038111838210176101fc57604052565b51906001600160a01b03821682036102f557565b9190820180921161033f57565b634e487b7160e01b5f52601160045260245ffd5b80518210156103675760209160051b010190565b634e487b7160e01b5f52603260045260245ffd5b6002548110156103675760025f5260205f209060011b01905f90565b6001600160a01b0316801561040957600254905f5b8281106103bb575050505f1990565b816103c58261037b565b50546001600160a01b03161480156103ea575b6103e4576001016103ac565b91505090565b50816103f58261037b565b50600101546001600160a01b0316146103d8565b505f199056fe6080806040526004361015610012575f80fd5b5f3560e01c90816301ffc9a71461100e575080631798d73a14610c55578063181f5a7714610ba85780631f40035e1461099657806345218dda146107ae57806379ba5097146106cb5780637df73e27146106845780638d361e43146103b75780638da5cb5b1461036657806394cf795e1461022d578063b715be81146101f2578063d1767c361461019a5763f2fde38b146100ab575f80fd5b346101965760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101965773ffffffffffffffffffffffffffffffffffffffff6100f7611146565b6100ff611256565b1633811461016e57807fffffffffffffffffffffffff00000000000000000000000000000000000000005f5416175f5573ffffffffffffffffffffffffffffffffffffffff600154167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12785f80a3005b7fdad89dca000000000000000000000000000000000000000000000000000000005f5260045ffd5b5f80fd5b34610196575f7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101965760206040517f00000000000000000000000000000000000000000000000000000000000000008152f35b34610196575f7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610196576020600254604051908152f35b34610196575f7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101965760025461026f61026a8261112e565b6110ea565b9080825260208201908160025f527f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace5f915b838310610311578486604051918291602083019060208452518091526040830191905f5b8181106102d3575050500390f35b919350916020604060019273ffffffffffffffffffffffffffffffffffffffff838851828151168452015116838201520194019101918493926102c5565b600260206001926103206110ca565b73ffffffffffffffffffffffffffffffffffffffff865416815273ffffffffffffffffffffffffffffffffffffffff8587015416838201528152019201920191906102a1565b34610196575f7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261019657602073ffffffffffffffffffffffffffffffffffffffff60015416604051908152f35b346101965760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101965760043567ffffffffffffffff81116101965761040690369060040161118a565b61040e611256565b5f5b81811061041957005b61043461042f61042a838587611225565b611235565b61129f565b5f811261062d576002547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff810190811161060057808203610563575b50506002548015610536577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01906104a782611209565b92909261050a575f60018482829655015560025573ffffffffffffffffffffffffffffffffffffffff6104de61042a838688611225565b167f3525e22824a8a7df2c9a6029941c824cf95b6447f1e13d5128fd3826d35afe8b5f80a25b01610410565b7f4e487b71000000000000000000000000000000000000000000000000000000005f525f60045260245ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603160045260245ffd5b61056f61057691611209565b5091611209565b91909161050a5781811461047057600173ffffffffffffffffffffffffffffffffffffffff8183828080965416167fffffffffffffffffffffffff0000000000000000000000000000000000000000875416178655015416920191167fffffffffffffffffffffffff00000000000000000000000000000000000000008254161790558380610470565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b50807fc4c107ee2918eac8fef38e116b14c35374d926cbae170b911dd930bd52050506602061066261042a6001958789611225565b73ffffffffffffffffffffffffffffffffffffffff60405191168152a1610504565b346101965760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101965760205f6106c161042f611146565b1215604051908152f35b34610196575f7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610196575f5473ffffffffffffffffffffffffffffffffffffffff81163303610786577fffffffffffffffffffffffff0000000000000000000000000000000000000000600154913382841617600155165f5573ffffffffffffffffffffffffffffffffffffffff3391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e05f80a3005b7f02b543c6000000000000000000000000000000000000000000000000000000005f5260045ffd5b346101965760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101965760043567ffffffffffffffff8111610196576107fd90369060040161118a565b90610806611256565b5f5b82811061081157005b61082261042f61042a838686611225565b5f81126109475761083290611209565b50600181019073ffffffffffffffffffffffffffffffffffffffff8254169182156108f857602073ffffffffffffffffffffffffffffffffffffffff837f7ed3a7e0dcd2674bb0a11d4d4c383b1933b493f877f36cb3113256a5e46f39b69382806001999897541696167fffffffffffffffffffffffff00000000000000000000000000000000000000008354161782557fffffffffffffffffffffffff000000000000000000000000000000000000000081541690555416604051908152a201610808565b73ffffffffffffffffffffffffffffffffffffffff61091b61042a868989611225565b7f04350e68000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b73ffffffffffffffffffffffffffffffffffffffff61096a61042a848787611225565b7f548fe13f000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b346101965760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101965760043567ffffffffffffffff8111610196576109e590369060040161118a565b60243567ffffffffffffffff811161019657610a0590369060040161118a565b92610a0e611256565b838303610b80575f5b838110610a2057005b610a3161042f61042a838786611225565b5f8112610b5d57610a4190611209565b505f610a5461042f61042a858a89611225565b1215610b0e5790837f8995392911585a9252d0e41882b03a1e008df154f884ed48ba604cd167c9797b604073ffffffffffffffffffffffffffffffffffffffff858a82610aae61042a896001809c0194848654169a611225565b167fffffffffffffffffffffffff000000000000000000000000000000000000000082541617905554169273ffffffffffffffffffffffffffffffffffffffff610afc61042a878d8c611225565b8351928352166020820152a201610a17565b73ffffffffffffffffffffffffffffffffffffffff610b3161042a848988611225565b7fe021c4f2000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b73ffffffffffffffffffffffffffffffffffffffff61096a61042a848887611225565b7f7db491eb000000000000000000000000000000000000000000000000000000005f5260045ffd5b34610196575f7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261019657610bde6110ca565b60148152604060208201917f5369676e6572526567697374727920312e302e3000000000000000000000000083527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8351948593602085525180918160208701528686015e5f85828601015201168101030190f35b346101965760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101965760043567ffffffffffffffff81116101965736602382011215610196578060040135610cb261026a8261112e565b916024602084848152019260061b8201019036821161019657602401915b818310610fd25783610ce0611256565b8051610cef60025480926111bb565b7f000000000000000000000000000000000000000000000000000000000000000010610f9857505f5b8151811015610f9657610d2b81836111c8565b5173ffffffffffffffffffffffffffffffffffffffff8151168015610f6b575073ffffffffffffffffffffffffffffffffffffffff81511690602081019173ffffffffffffffffffffffffffffffffffffffff835116809114610f4057505f610daa73ffffffffffffffffffffffffffffffffffffffff83511661129f565b1215610efd57505f610dd273ffffffffffffffffffffffffffffffffffffffff83511661129f565b1215610efd5750610de381836111c8565b519060025468010000000000000000811015610ed057806001610e099201600255611209565b61050a5773ffffffffffffffffffffffffffffffffffffffff600181602086828085995116167fffffffffffffffffffffffff0000000000000000000000000000000000000000875416178655015116920191167fffffffffffffffffffffffff000000000000000000000000000000000000000082541617905573ffffffffffffffffffffffffffffffffffffffff610ea382856111c8565b5151167f47d1c22a25bb3a5d4e481b9b1e6944c2eade3181a0a20b495ed61d35b5323f245f80a201610d18565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b73ffffffffffffffffffffffffffffffffffffffff9051167fe021c4f2000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b7fe021c4f2000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b7f5a1d346d000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b005b610fa38183516111bb565b907f46161fe4000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b604083360312610196576020604091610fe96110ca565b610ff286611169565b8152610fff838701611169565b83820152815201920191610cd0565b346101965760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261019657600435907fffffffff00000000000000000000000000000000000000000000000000000000821680920361019657817f5e2df9f800000000000000000000000000000000000000000000000000000000602093149081156110a0575b5015158152f35b7f01ffc9a70000000000000000000000000000000000000000000000000000000091501483611099565b604051906040820182811067ffffffffffffffff821117610ed057604052565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f604051930116820182811067ffffffffffffffff821117610ed057604052565b67ffffffffffffffff8111610ed05760051b60200190565b6004359073ffffffffffffffffffffffffffffffffffffffff8216820361019657565b359073ffffffffffffffffffffffffffffffffffffffff8216820361019657565b9181601f840112156101965782359167ffffffffffffffff8311610196576020808501948460051b01011161019657565b9190820180921161060057565b80518210156111dc5760209160051b010190565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b6002548110156111dc5760025f5260205f209060011b01905f90565b91908110156111dc5760051b0190565b3573ffffffffffffffffffffffffffffffffffffffff811681036101965790565b73ffffffffffffffffffffffffffffffffffffffff60015416330361127757565b7f2b5c74de000000000000000000000000000000000000000000000000000000005f5260045ffd5b73ffffffffffffffffffffffffffffffffffffffff16801561135757600254905f5b8281106112ef575050507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90565b8173ffffffffffffffffffffffffffffffffffffffff61130e83611209565b50541614801561132b575b611325576001016112c1565b91505090565b508173ffffffffffffffffffffffffffffffffffffffff600161134d84611209565b5001541614611319565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9056fea164736f6c634300081c000a",
}

var SignerRegistryABI = SignerRegistryMetaData.ABI

var SignerRegistryBin = SignerRegistryMetaData.Bin

func DeploySignerRegistry(auth *bind.TransactOpts, backend bind.ContractBackend, maxSigners *big.Int, newSigners []ISignerRegistrySigner) (common.Address, *types.Transaction, *SignerRegistry, error) {
	parsed, err := SignerRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SignerRegistryBin), backend, maxSigners, newSigners)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SignerRegistry{address: address, abi: *parsed, SignerRegistryCaller: SignerRegistryCaller{contract: contract}, SignerRegistryTransactor: SignerRegistryTransactor{contract: contract}, SignerRegistryFilterer: SignerRegistryFilterer{contract: contract}}, nil
}

type SignerRegistry struct {
	address common.Address
	abi     abi.ABI
	SignerRegistryCaller
	SignerRegistryTransactor
	SignerRegistryFilterer
}

type SignerRegistryCaller struct {
	contract *bind.BoundContract
}

type SignerRegistryTransactor struct {
	contract *bind.BoundContract
}

type SignerRegistryFilterer struct {
	contract *bind.BoundContract
}

type SignerRegistrySession struct {
	Contract     *SignerRegistry
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type SignerRegistryCallerSession struct {
	Contract *SignerRegistryCaller
	CallOpts bind.CallOpts
}

type SignerRegistryTransactorSession struct {
	Contract     *SignerRegistryTransactor
	TransactOpts bind.TransactOpts
}

type SignerRegistryRaw struct {
	Contract *SignerRegistry
}

type SignerRegistryCallerRaw struct {
	Contract *SignerRegistryCaller
}

type SignerRegistryTransactorRaw struct {
	Contract *SignerRegistryTransactor
}

func NewSignerRegistry(address common.Address, backend bind.ContractBackend) (*SignerRegistry, error) {
	abi, err := abi.JSON(strings.NewReader(SignerRegistryABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindSignerRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SignerRegistry{address: address, abi: abi, SignerRegistryCaller: SignerRegistryCaller{contract: contract}, SignerRegistryTransactor: SignerRegistryTransactor{contract: contract}, SignerRegistryFilterer: SignerRegistryFilterer{contract: contract}}, nil
}

func NewSignerRegistryCaller(address common.Address, caller bind.ContractCaller) (*SignerRegistryCaller, error) {
	contract, err := bindSignerRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SignerRegistryCaller{contract: contract}, nil
}

func NewSignerRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*SignerRegistryTransactor, error) {
	contract, err := bindSignerRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SignerRegistryTransactor{contract: contract}, nil
}

func NewSignerRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*SignerRegistryFilterer, error) {
	contract, err := bindSignerRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SignerRegistryFilterer{contract: contract}, nil
}

func bindSignerRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SignerRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_SignerRegistry *SignerRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SignerRegistry.Contract.SignerRegistryCaller.contract.Call(opts, result, method, params...)
}

func (_SignerRegistry *SignerRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SignerRegistry.Contract.SignerRegistryTransactor.contract.Transfer(opts)
}

func (_SignerRegistry *SignerRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SignerRegistry.Contract.SignerRegistryTransactor.contract.Transact(opts, method, params...)
}

func (_SignerRegistry *SignerRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SignerRegistry.Contract.contract.Call(opts, result, method, params...)
}

func (_SignerRegistry *SignerRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SignerRegistry.Contract.contract.Transfer(opts)
}

func (_SignerRegistry *SignerRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SignerRegistry.Contract.contract.Transact(opts, method, params...)
}

func (_SignerRegistry *SignerRegistryCaller) GetMaxSigners(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SignerRegistry.contract.Call(opts, &out, "getMaxSigners")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_SignerRegistry *SignerRegistrySession) GetMaxSigners() (*big.Int, error) {
	return _SignerRegistry.Contract.GetMaxSigners(&_SignerRegistry.CallOpts)
}

func (_SignerRegistry *SignerRegistryCallerSession) GetMaxSigners() (*big.Int, error) {
	return _SignerRegistry.Contract.GetMaxSigners(&_SignerRegistry.CallOpts)
}

func (_SignerRegistry *SignerRegistryCaller) GetSignerCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SignerRegistry.contract.Call(opts, &out, "getSignerCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_SignerRegistry *SignerRegistrySession) GetSignerCount() (*big.Int, error) {
	return _SignerRegistry.Contract.GetSignerCount(&_SignerRegistry.CallOpts)
}

func (_SignerRegistry *SignerRegistryCallerSession) GetSignerCount() (*big.Int, error) {
	return _SignerRegistry.Contract.GetSignerCount(&_SignerRegistry.CallOpts)
}

func (_SignerRegistry *SignerRegistryCaller) GetSigners(opts *bind.CallOpts) ([]ISignerRegistrySigner, error) {
	var out []interface{}
	err := _SignerRegistry.contract.Call(opts, &out, "getSigners")

	if err != nil {
		return *new([]ISignerRegistrySigner), err
	}

	out0 := *abi.ConvertType(out[0], new([]ISignerRegistrySigner)).(*[]ISignerRegistrySigner)

	return out0, err

}

func (_SignerRegistry *SignerRegistrySession) GetSigners() ([]ISignerRegistrySigner, error) {
	return _SignerRegistry.Contract.GetSigners(&_SignerRegistry.CallOpts)
}

func (_SignerRegistry *SignerRegistryCallerSession) GetSigners() ([]ISignerRegistrySigner, error) {
	return _SignerRegistry.Contract.GetSigners(&_SignerRegistry.CallOpts)
}

func (_SignerRegistry *SignerRegistryCaller) IsSigner(opts *bind.CallOpts, signerAddress common.Address) (bool, error) {
	var out []interface{}
	err := _SignerRegistry.contract.Call(opts, &out, "isSigner", signerAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_SignerRegistry *SignerRegistrySession) IsSigner(signerAddress common.Address) (bool, error) {
	return _SignerRegistry.Contract.IsSigner(&_SignerRegistry.CallOpts, signerAddress)
}

func (_SignerRegistry *SignerRegistryCallerSession) IsSigner(signerAddress common.Address) (bool, error) {
	return _SignerRegistry.Contract.IsSigner(&_SignerRegistry.CallOpts, signerAddress)
}

func (_SignerRegistry *SignerRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SignerRegistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_SignerRegistry *SignerRegistrySession) Owner() (common.Address, error) {
	return _SignerRegistry.Contract.Owner(&_SignerRegistry.CallOpts)
}

func (_SignerRegistry *SignerRegistryCallerSession) Owner() (common.Address, error) {
	return _SignerRegistry.Contract.Owner(&_SignerRegistry.CallOpts)
}

func (_SignerRegistry *SignerRegistryCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _SignerRegistry.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_SignerRegistry *SignerRegistrySession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _SignerRegistry.Contract.SupportsInterface(&_SignerRegistry.CallOpts, interfaceId)
}

func (_SignerRegistry *SignerRegistryCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _SignerRegistry.Contract.SupportsInterface(&_SignerRegistry.CallOpts, interfaceId)
}

func (_SignerRegistry *SignerRegistryCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _SignerRegistry.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_SignerRegistry *SignerRegistrySession) TypeAndVersion() (string, error) {
	return _SignerRegistry.Contract.TypeAndVersion(&_SignerRegistry.CallOpts)
}

func (_SignerRegistry *SignerRegistryCallerSession) TypeAndVersion() (string, error) {
	return _SignerRegistry.Contract.TypeAndVersion(&_SignerRegistry.CallOpts)
}

func (_SignerRegistry *SignerRegistryTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SignerRegistry.contract.Transact(opts, "acceptOwnership")
}

func (_SignerRegistry *SignerRegistrySession) AcceptOwnership() (*types.Transaction, error) {
	return _SignerRegistry.Contract.AcceptOwnership(&_SignerRegistry.TransactOpts)
}

func (_SignerRegistry *SignerRegistryTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _SignerRegistry.Contract.AcceptOwnership(&_SignerRegistry.TransactOpts)
}

func (_SignerRegistry *SignerRegistryTransactor) AddSigners(opts *bind.TransactOpts, newSigners []ISignerRegistrySigner) (*types.Transaction, error) {
	return _SignerRegistry.contract.Transact(opts, "addSigners", newSigners)
}

func (_SignerRegistry *SignerRegistrySession) AddSigners(newSigners []ISignerRegistrySigner) (*types.Transaction, error) {
	return _SignerRegistry.Contract.AddSigners(&_SignerRegistry.TransactOpts, newSigners)
}

func (_SignerRegistry *SignerRegistryTransactorSession) AddSigners(newSigners []ISignerRegistrySigner) (*types.Transaction, error) {
	return _SignerRegistry.Contract.AddSigners(&_SignerRegistry.TransactOpts, newSigners)
}

func (_SignerRegistry *SignerRegistryTransactor) PromoteNewSignerAddresses(opts *bind.TransactOpts, existingSignerAddresses []common.Address) (*types.Transaction, error) {
	return _SignerRegistry.contract.Transact(opts, "promoteNewSignerAddresses", existingSignerAddresses)
}

func (_SignerRegistry *SignerRegistrySession) PromoteNewSignerAddresses(existingSignerAddresses []common.Address) (*types.Transaction, error) {
	return _SignerRegistry.Contract.PromoteNewSignerAddresses(&_SignerRegistry.TransactOpts, existingSignerAddresses)
}

func (_SignerRegistry *SignerRegistryTransactorSession) PromoteNewSignerAddresses(existingSignerAddresses []common.Address) (*types.Transaction, error) {
	return _SignerRegistry.Contract.PromoteNewSignerAddresses(&_SignerRegistry.TransactOpts, existingSignerAddresses)
}

func (_SignerRegistry *SignerRegistryTransactor) RemoveSigners(opts *bind.TransactOpts, signersToRemove []common.Address) (*types.Transaction, error) {
	return _SignerRegistry.contract.Transact(opts, "removeSigners", signersToRemove)
}

func (_SignerRegistry *SignerRegistrySession) RemoveSigners(signersToRemove []common.Address) (*types.Transaction, error) {
	return _SignerRegistry.Contract.RemoveSigners(&_SignerRegistry.TransactOpts, signersToRemove)
}

func (_SignerRegistry *SignerRegistryTransactorSession) RemoveSigners(signersToRemove []common.Address) (*types.Transaction, error) {
	return _SignerRegistry.Contract.RemoveSigners(&_SignerRegistry.TransactOpts, signersToRemove)
}

func (_SignerRegistry *SignerRegistryTransactor) SetNewSignerAddresses(opts *bind.TransactOpts, existingSignerAddresses []common.Address, newSignerAddresses []common.Address) (*types.Transaction, error) {
	return _SignerRegistry.contract.Transact(opts, "setNewSignerAddresses", existingSignerAddresses, newSignerAddresses)
}

func (_SignerRegistry *SignerRegistrySession) SetNewSignerAddresses(existingSignerAddresses []common.Address, newSignerAddresses []common.Address) (*types.Transaction, error) {
	return _SignerRegistry.Contract.SetNewSignerAddresses(&_SignerRegistry.TransactOpts, existingSignerAddresses, newSignerAddresses)
}

func (_SignerRegistry *SignerRegistryTransactorSession) SetNewSignerAddresses(existingSignerAddresses []common.Address, newSignerAddresses []common.Address) (*types.Transaction, error) {
	return _SignerRegistry.Contract.SetNewSignerAddresses(&_SignerRegistry.TransactOpts, existingSignerAddresses, newSignerAddresses)
}

func (_SignerRegistry *SignerRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _SignerRegistry.contract.Transact(opts, "transferOwnership", to)
}

func (_SignerRegistry *SignerRegistrySession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _SignerRegistry.Contract.TransferOwnership(&_SignerRegistry.TransactOpts, to)
}

func (_SignerRegistry *SignerRegistryTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _SignerRegistry.Contract.TransferOwnership(&_SignerRegistry.TransactOpts, to)
}

type SignerRegistryNewSignerAddressSetIterator struct {
	Event *SignerRegistryNewSignerAddressSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *SignerRegistryNewSignerAddressSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SignerRegistryNewSignerAddressSet)
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
		it.Event = new(SignerRegistryNewSignerAddressSet)
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

func (it *SignerRegistryNewSignerAddressSetIterator) Error() error {
	return it.fail
}

func (it *SignerRegistryNewSignerAddressSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type SignerRegistryNewSignerAddressSet struct {
	EvmAddress            common.Address
	PreviousNewEVMAddress common.Address
	NewEVMAddress         common.Address
	Raw                   types.Log
}

func (_SignerRegistry *SignerRegistryFilterer) FilterNewSignerAddressSet(opts *bind.FilterOpts, evmAddress []common.Address) (*SignerRegistryNewSignerAddressSetIterator, error) {

	var evmAddressRule []interface{}
	for _, evmAddressItem := range evmAddress {
		evmAddressRule = append(evmAddressRule, evmAddressItem)
	}

	logs, sub, err := _SignerRegistry.contract.FilterLogs(opts, "NewSignerAddressSet", evmAddressRule)
	if err != nil {
		return nil, err
	}
	return &SignerRegistryNewSignerAddressSetIterator{contract: _SignerRegistry.contract, event: "NewSignerAddressSet", logs: logs, sub: sub}, nil
}

func (_SignerRegistry *SignerRegistryFilterer) WatchNewSignerAddressSet(opts *bind.WatchOpts, sink chan<- *SignerRegistryNewSignerAddressSet, evmAddress []common.Address) (event.Subscription, error) {

	var evmAddressRule []interface{}
	for _, evmAddressItem := range evmAddress {
		evmAddressRule = append(evmAddressRule, evmAddressItem)
	}

	logs, sub, err := _SignerRegistry.contract.WatchLogs(opts, "NewSignerAddressSet", evmAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(SignerRegistryNewSignerAddressSet)
				if err := _SignerRegistry.contract.UnpackLog(event, "NewSignerAddressSet", log); err != nil {
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

func (_SignerRegistry *SignerRegistryFilterer) ParseNewSignerAddressSet(log types.Log) (*SignerRegistryNewSignerAddressSet, error) {
	event := new(SignerRegistryNewSignerAddressSet)
	if err := _SignerRegistry.contract.UnpackLog(event, "NewSignerAddressSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type SignerRegistryOwnershipTransferRequestedIterator struct {
	Event *SignerRegistryOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *SignerRegistryOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SignerRegistryOwnershipTransferRequested)
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
		it.Event = new(SignerRegistryOwnershipTransferRequested)
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

func (it *SignerRegistryOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *SignerRegistryOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type SignerRegistryOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_SignerRegistry *SignerRegistryFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*SignerRegistryOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _SignerRegistry.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &SignerRegistryOwnershipTransferRequestedIterator{contract: _SignerRegistry.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_SignerRegistry *SignerRegistryFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *SignerRegistryOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _SignerRegistry.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(SignerRegistryOwnershipTransferRequested)
				if err := _SignerRegistry.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_SignerRegistry *SignerRegistryFilterer) ParseOwnershipTransferRequested(log types.Log) (*SignerRegistryOwnershipTransferRequested, error) {
	event := new(SignerRegistryOwnershipTransferRequested)
	if err := _SignerRegistry.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type SignerRegistryOwnershipTransferredIterator struct {
	Event *SignerRegistryOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *SignerRegistryOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SignerRegistryOwnershipTransferred)
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
		it.Event = new(SignerRegistryOwnershipTransferred)
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

func (it *SignerRegistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *SignerRegistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type SignerRegistryOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_SignerRegistry *SignerRegistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*SignerRegistryOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _SignerRegistry.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &SignerRegistryOwnershipTransferredIterator{contract: _SignerRegistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_SignerRegistry *SignerRegistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *SignerRegistryOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _SignerRegistry.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(SignerRegistryOwnershipTransferred)
				if err := _SignerRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_SignerRegistry *SignerRegistryFilterer) ParseOwnershipTransferred(log types.Log) (*SignerRegistryOwnershipTransferred, error) {
	event := new(SignerRegistryOwnershipTransferred)
	if err := _SignerRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type SignerRegistrySignerAddedIterator struct {
	Event *SignerRegistrySignerAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *SignerRegistrySignerAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SignerRegistrySignerAdded)
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
		it.Event = new(SignerRegistrySignerAdded)
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

func (it *SignerRegistrySignerAddedIterator) Error() error {
	return it.fail
}

func (it *SignerRegistrySignerAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type SignerRegistrySignerAdded struct {
	EvmAddress common.Address
	Raw        types.Log
}

func (_SignerRegistry *SignerRegistryFilterer) FilterSignerAdded(opts *bind.FilterOpts, evmAddress []common.Address) (*SignerRegistrySignerAddedIterator, error) {

	var evmAddressRule []interface{}
	for _, evmAddressItem := range evmAddress {
		evmAddressRule = append(evmAddressRule, evmAddressItem)
	}

	logs, sub, err := _SignerRegistry.contract.FilterLogs(opts, "SignerAdded", evmAddressRule)
	if err != nil {
		return nil, err
	}
	return &SignerRegistrySignerAddedIterator{contract: _SignerRegistry.contract, event: "SignerAdded", logs: logs, sub: sub}, nil
}

func (_SignerRegistry *SignerRegistryFilterer) WatchSignerAdded(opts *bind.WatchOpts, sink chan<- *SignerRegistrySignerAdded, evmAddress []common.Address) (event.Subscription, error) {

	var evmAddressRule []interface{}
	for _, evmAddressItem := range evmAddress {
		evmAddressRule = append(evmAddressRule, evmAddressItem)
	}

	logs, sub, err := _SignerRegistry.contract.WatchLogs(opts, "SignerAdded", evmAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(SignerRegistrySignerAdded)
				if err := _SignerRegistry.contract.UnpackLog(event, "SignerAdded", log); err != nil {
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

func (_SignerRegistry *SignerRegistryFilterer) ParseSignerAdded(log types.Log) (*SignerRegistrySignerAdded, error) {
	event := new(SignerRegistrySignerAdded)
	if err := _SignerRegistry.contract.UnpackLog(event, "SignerAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type SignerRegistrySignerAddressPromotedIterator struct {
	Event *SignerRegistrySignerAddressPromoted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *SignerRegistrySignerAddressPromotedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SignerRegistrySignerAddressPromoted)
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
		it.Event = new(SignerRegistrySignerAddressPromoted)
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

func (it *SignerRegistrySignerAddressPromotedIterator) Error() error {
	return it.fail
}

func (it *SignerRegistrySignerAddressPromotedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type SignerRegistrySignerAddressPromoted struct {
	EvmAddress    common.Address
	NewEVMAddress common.Address
	Raw           types.Log
}

func (_SignerRegistry *SignerRegistryFilterer) FilterSignerAddressPromoted(opts *bind.FilterOpts, evmAddress []common.Address) (*SignerRegistrySignerAddressPromotedIterator, error) {

	var evmAddressRule []interface{}
	for _, evmAddressItem := range evmAddress {
		evmAddressRule = append(evmAddressRule, evmAddressItem)
	}

	logs, sub, err := _SignerRegistry.contract.FilterLogs(opts, "SignerAddressPromoted", evmAddressRule)
	if err != nil {
		return nil, err
	}
	return &SignerRegistrySignerAddressPromotedIterator{contract: _SignerRegistry.contract, event: "SignerAddressPromoted", logs: logs, sub: sub}, nil
}

func (_SignerRegistry *SignerRegistryFilterer) WatchSignerAddressPromoted(opts *bind.WatchOpts, sink chan<- *SignerRegistrySignerAddressPromoted, evmAddress []common.Address) (event.Subscription, error) {

	var evmAddressRule []interface{}
	for _, evmAddressItem := range evmAddress {
		evmAddressRule = append(evmAddressRule, evmAddressItem)
	}

	logs, sub, err := _SignerRegistry.contract.WatchLogs(opts, "SignerAddressPromoted", evmAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(SignerRegistrySignerAddressPromoted)
				if err := _SignerRegistry.contract.UnpackLog(event, "SignerAddressPromoted", log); err != nil {
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

func (_SignerRegistry *SignerRegistryFilterer) ParseSignerAddressPromoted(log types.Log) (*SignerRegistrySignerAddressPromoted, error) {
	event := new(SignerRegistrySignerAddressPromoted)
	if err := _SignerRegistry.contract.UnpackLog(event, "SignerAddressPromoted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type SignerRegistrySignerNotRemovedDueToNoMatchIterator struct {
	Event *SignerRegistrySignerNotRemovedDueToNoMatch

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *SignerRegistrySignerNotRemovedDueToNoMatchIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SignerRegistrySignerNotRemovedDueToNoMatch)
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
		it.Event = new(SignerRegistrySignerNotRemovedDueToNoMatch)
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

func (it *SignerRegistrySignerNotRemovedDueToNoMatchIterator) Error() error {
	return it.fail
}

func (it *SignerRegistrySignerNotRemovedDueToNoMatchIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type SignerRegistrySignerNotRemovedDueToNoMatch struct {
	SignerAddress common.Address
	Raw           types.Log
}

func (_SignerRegistry *SignerRegistryFilterer) FilterSignerNotRemovedDueToNoMatch(opts *bind.FilterOpts) (*SignerRegistrySignerNotRemovedDueToNoMatchIterator, error) {

	logs, sub, err := _SignerRegistry.contract.FilterLogs(opts, "SignerNotRemovedDueToNoMatch")
	if err != nil {
		return nil, err
	}
	return &SignerRegistrySignerNotRemovedDueToNoMatchIterator{contract: _SignerRegistry.contract, event: "SignerNotRemovedDueToNoMatch", logs: logs, sub: sub}, nil
}

func (_SignerRegistry *SignerRegistryFilterer) WatchSignerNotRemovedDueToNoMatch(opts *bind.WatchOpts, sink chan<- *SignerRegistrySignerNotRemovedDueToNoMatch) (event.Subscription, error) {

	logs, sub, err := _SignerRegistry.contract.WatchLogs(opts, "SignerNotRemovedDueToNoMatch")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(SignerRegistrySignerNotRemovedDueToNoMatch)
				if err := _SignerRegistry.contract.UnpackLog(event, "SignerNotRemovedDueToNoMatch", log); err != nil {
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

func (_SignerRegistry *SignerRegistryFilterer) ParseSignerNotRemovedDueToNoMatch(log types.Log) (*SignerRegistrySignerNotRemovedDueToNoMatch, error) {
	event := new(SignerRegistrySignerNotRemovedDueToNoMatch)
	if err := _SignerRegistry.contract.UnpackLog(event, "SignerNotRemovedDueToNoMatch", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type SignerRegistrySignerRemovedIterator struct {
	Event *SignerRegistrySignerRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *SignerRegistrySignerRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SignerRegistrySignerRemoved)
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
		it.Event = new(SignerRegistrySignerRemoved)
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

func (it *SignerRegistrySignerRemovedIterator) Error() error {
	return it.fail
}

func (it *SignerRegistrySignerRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type SignerRegistrySignerRemoved struct {
	EvmAddress common.Address
	Raw        types.Log
}

func (_SignerRegistry *SignerRegistryFilterer) FilterSignerRemoved(opts *bind.FilterOpts, evmAddress []common.Address) (*SignerRegistrySignerRemovedIterator, error) {

	var evmAddressRule []interface{}
	for _, evmAddressItem := range evmAddress {
		evmAddressRule = append(evmAddressRule, evmAddressItem)
	}

	logs, sub, err := _SignerRegistry.contract.FilterLogs(opts, "SignerRemoved", evmAddressRule)
	if err != nil {
		return nil, err
	}
	return &SignerRegistrySignerRemovedIterator{contract: _SignerRegistry.contract, event: "SignerRemoved", logs: logs, sub: sub}, nil
}

func (_SignerRegistry *SignerRegistryFilterer) WatchSignerRemoved(opts *bind.WatchOpts, sink chan<- *SignerRegistrySignerRemoved, evmAddress []common.Address) (event.Subscription, error) {

	var evmAddressRule []interface{}
	for _, evmAddressItem := range evmAddress {
		evmAddressRule = append(evmAddressRule, evmAddressItem)
	}

	logs, sub, err := _SignerRegistry.contract.WatchLogs(opts, "SignerRemoved", evmAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(SignerRegistrySignerRemoved)
				if err := _SignerRegistry.contract.UnpackLog(event, "SignerRemoved", log); err != nil {
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

func (_SignerRegistry *SignerRegistryFilterer) ParseSignerRemoved(log types.Log) (*SignerRegistrySignerRemoved, error) {
	event := new(SignerRegistrySignerRemoved)
	if err := _SignerRegistry.contract.UnpackLog(event, "SignerRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_SignerRegistry *SignerRegistry) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _SignerRegistry.abi.Events["NewSignerAddressSet"].ID:
		return _SignerRegistry.ParseNewSignerAddressSet(log)
	case _SignerRegistry.abi.Events["OwnershipTransferRequested"].ID:
		return _SignerRegistry.ParseOwnershipTransferRequested(log)
	case _SignerRegistry.abi.Events["OwnershipTransferred"].ID:
		return _SignerRegistry.ParseOwnershipTransferred(log)
	case _SignerRegistry.abi.Events["SignerAdded"].ID:
		return _SignerRegistry.ParseSignerAdded(log)
	case _SignerRegistry.abi.Events["SignerAddressPromoted"].ID:
		return _SignerRegistry.ParseSignerAddressPromoted(log)
	case _SignerRegistry.abi.Events["SignerNotRemovedDueToNoMatch"].ID:
		return _SignerRegistry.ParseSignerNotRemovedDueToNoMatch(log)
	case _SignerRegistry.abi.Events["SignerRemoved"].ID:
		return _SignerRegistry.ParseSignerRemoved(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (SignerRegistryNewSignerAddressSet) Topic() common.Hash {
	return common.HexToHash("0x8995392911585a9252d0e41882b03a1e008df154f884ed48ba604cd167c9797b")
}

func (SignerRegistryOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (SignerRegistryOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (SignerRegistrySignerAdded) Topic() common.Hash {
	return common.HexToHash("0x47d1c22a25bb3a5d4e481b9b1e6944c2eade3181a0a20b495ed61d35b5323f24")
}

func (SignerRegistrySignerAddressPromoted) Topic() common.Hash {
	return common.HexToHash("0x7ed3a7e0dcd2674bb0a11d4d4c383b1933b493f877f36cb3113256a5e46f39b6")
}

func (SignerRegistrySignerNotRemovedDueToNoMatch) Topic() common.Hash {
	return common.HexToHash("0xc4c107ee2918eac8fef38e116b14c35374d926cbae170b911dd930bd52050506")
}

func (SignerRegistrySignerRemoved) Topic() common.Hash {
	return common.HexToHash("0x3525e22824a8a7df2c9a6029941c824cf95b6447f1e13d5128fd3826d35afe8b")
}

func (_SignerRegistry *SignerRegistry) Address() common.Address {
	return _SignerRegistry.address
}

type SignerRegistryInterface interface {
	GetMaxSigners(opts *bind.CallOpts) (*big.Int, error)

	GetSignerCount(opts *bind.CallOpts) (*big.Int, error)

	GetSigners(opts *bind.CallOpts) ([]ISignerRegistrySigner, error)

	IsSigner(opts *bind.CallOpts, signerAddress common.Address) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AddSigners(opts *bind.TransactOpts, newSigners []ISignerRegistrySigner) (*types.Transaction, error)

	PromoteNewSignerAddresses(opts *bind.TransactOpts, existingSignerAddresses []common.Address) (*types.Transaction, error)

	RemoveSigners(opts *bind.TransactOpts, signersToRemove []common.Address) (*types.Transaction, error)

	SetNewSignerAddresses(opts *bind.TransactOpts, existingSignerAddresses []common.Address, newSignerAddresses []common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterNewSignerAddressSet(opts *bind.FilterOpts, evmAddress []common.Address) (*SignerRegistryNewSignerAddressSetIterator, error)

	WatchNewSignerAddressSet(opts *bind.WatchOpts, sink chan<- *SignerRegistryNewSignerAddressSet, evmAddress []common.Address) (event.Subscription, error)

	ParseNewSignerAddressSet(log types.Log) (*SignerRegistryNewSignerAddressSet, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*SignerRegistryOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *SignerRegistryOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*SignerRegistryOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*SignerRegistryOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *SignerRegistryOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*SignerRegistryOwnershipTransferred, error)

	FilterSignerAdded(opts *bind.FilterOpts, evmAddress []common.Address) (*SignerRegistrySignerAddedIterator, error)

	WatchSignerAdded(opts *bind.WatchOpts, sink chan<- *SignerRegistrySignerAdded, evmAddress []common.Address) (event.Subscription, error)

	ParseSignerAdded(log types.Log) (*SignerRegistrySignerAdded, error)

	FilterSignerAddressPromoted(opts *bind.FilterOpts, evmAddress []common.Address) (*SignerRegistrySignerAddressPromotedIterator, error)

	WatchSignerAddressPromoted(opts *bind.WatchOpts, sink chan<- *SignerRegistrySignerAddressPromoted, evmAddress []common.Address) (event.Subscription, error)

	ParseSignerAddressPromoted(log types.Log) (*SignerRegistrySignerAddressPromoted, error)

	FilterSignerNotRemovedDueToNoMatch(opts *bind.FilterOpts) (*SignerRegistrySignerNotRemovedDueToNoMatchIterator, error)

	WatchSignerNotRemovedDueToNoMatch(opts *bind.WatchOpts, sink chan<- *SignerRegistrySignerNotRemovedDueToNoMatch) (event.Subscription, error)

	ParseSignerNotRemovedDueToNoMatch(log types.Log) (*SignerRegistrySignerNotRemovedDueToNoMatch, error)

	FilterSignerRemoved(opts *bind.FilterOpts, evmAddress []common.Address) (*SignerRegistrySignerRemovedIterator, error)

	WatchSignerRemoved(opts *bind.WatchOpts, sink chan<- *SignerRegistrySignerRemoved, evmAddress []common.Address) (event.Subscription, error)

	ParseSignerRemoved(log types.Log) (*SignerRegistrySignerRemoved, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
