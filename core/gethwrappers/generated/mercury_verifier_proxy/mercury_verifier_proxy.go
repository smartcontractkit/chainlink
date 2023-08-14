// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package mercury_verifier_proxy

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

type CommonAddressAndWeight struct {
	Addr   common.Address
	Weight *big.Int
}

var MercuryVerifierProxyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"accessController\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessForbidden\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BadVerification\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"verifier\",\"type\":\"address\"}],\"name\":\"ConfigDigestAlreadySet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"verifier\",\"type\":\"address\"}],\"name\":\"VerifierAlreadyInitialized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"VerifierInvalid\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"VerifierNotFound\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldAccessController\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAccessController\",\"type\":\"address\"}],\"name\":\"AccessControllerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldFeeManager\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newFeeManager\",\"type\":\"address\"}],\"name\":\"FeeManagerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"verifierAddress\",\"type\":\"address\"}],\"name\":\"VerifierInitialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"oldConfigDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"newConfigDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"verifierAddress\",\"type\":\"address\"}],\"name\":\"VerifierSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"verifierAddress\",\"type\":\"address\"}],\"name\":\"VerifierUnset\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"getVerifier\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"verifierAddress\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"verifierAddress\",\"type\":\"address\"}],\"name\":\"initializeVerifier\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_accessController\",\"outputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_feeManager\",\"outputs\":[{\"internalType\":\"contractIVerifierFeeManager\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"accessController\",\"type\":\"address\"}],\"name\":\"setAccessController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIVerifierFeeManager\",\"name\":\"feeManager\",\"type\":\"address\"}],\"name\":\"setFeeManager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"currentConfigDigest\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newConfigDigest\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"weight\",\"type\":\"uint256\"}],\"internalType\":\"structCommon.AddressAndWeight[]\",\"name\":\"addressesAndWeights\",\"type\":\"tuple[]\"}],\"name\":\"setVerifier\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"unsetVerifier\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"verify\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"verifierResponse\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50604051620015333803806200153383398101604081905261003191610189565b33806000816100875760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156100b7576100b7816100e0565b5050600480546001600160a01b0319166001600160a01b039390931692909217909155506101b9565b336001600160a01b038216036101385760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161007e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60006020828403121561019b57600080fd5b81516001600160a01b03811681146101b257600080fd5b9392505050565b61136a80620001c96000396000f3fe6080604052600436106100d25760003560e01c80638c2a4d531161007f57806394ba28461161005957806394ba284614610256578063eeb7b24814610283578063f08391d8146102c6578063f2fde38b146102e657600080fd5b80638c2a4d53146101f85780638da5cb5b146102185780638e760afe1461024357600080fd5b8063589ede28116100b0578063589ede28146101a35780636e914094146101c357806379ba5097146101e357600080fd5b8063181f5a77146100d757806338416b5b1461012f578063472d35b914610181575b600080fd5b3480156100e357600080fd5b5060408051808201909152601381527f566572696669657250726f787920312e312e300000000000000000000000000060208201525b6040516101269190610f7d565b60405180910390f35b34801561013b57600080fd5b5060055461015c9073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610126565b34801561018d57600080fd5b506101a161019c366004610fb9565b610306565b005b3480156101af57600080fd5b506101a16101be366004610fd6565b6103e2565b3480156101cf57600080fd5b506101a16101de366004611059565b61060d565b3480156101ef57600080fd5b506101a16106f9565b34801561020457600080fd5b506101a1610213366004610fb9565b6107f6565b34801561022457600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff1661015c565b610119610251366004611072565b610a27565b34801561026257600080fd5b5060045461015c9073ffffffffffffffffffffffffffffffffffffffff1681565b34801561028f57600080fd5b5061015c61029e366004611059565b60009081526003602052604090205473ffffffffffffffffffffffffffffffffffffffff1690565b3480156102d257600080fd5b506101a16102e1366004610fb9565b610cfc565b3480156102f257600080fd5b506101a1610301366004610fb9565b610d83565b61030e610d97565b73ffffffffffffffffffffffffffffffffffffffff811661035b576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6005805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff000000000000000000000000000000000000000083168117909355604080519190921680825260208201939093527f04628abcaa6b1674651352125cb94b65b289145bc2bc4d67720bb7d966372f0391015b60405180910390a15050565b600083815260036020526040902054839073ffffffffffffffffffffffffffffffffffffffff168015610465576040517f375d1fe60000000000000000000000000000000000000000000000000000000081526004810183905273ffffffffffffffffffffffffffffffffffffffff821660248201526044015b60405180910390fd5b3360009081526002602052604090205460ff166104ae576040517fef67f5d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600085815260036020526040902080547fffffffffffffffffffffffff0000000000000000000000000000000000000000163317905582156105c65760055473ffffffffffffffffffffffffffffffffffffffff16610539576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6005546040517f69fd2b3400000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff909116906369fd2b3490610593908890889088906004016110e4565b600060405180830381600087803b1580156105ad57600080fd5b505af11580156105c1573d6000803e3d6000fd5b505050505b6040805187815260208101879052338183015290517fbeb513e532542a562ac35699e7cd9ae7d198dcd3eee15bada6c857d28ceaddcf9181900360600190a1505050505050565b610615610d97565b60008181526003602052604090205473ffffffffffffffffffffffffffffffffffffffff1680610674576040517fb151802b0000000000000000000000000000000000000000000000000000000081526004810183905260240161045c565b6000828152600360205260409081902080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169055517f11dc15c4b8ac2b183166cc8427e5385a5ece8308217a4217338c6a7614845c4c906103d6908490849091825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b60015473ffffffffffffffffffffffffffffffffffffffff16331461077a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015260640161045c565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6107fe610d97565b8073ffffffffffffffffffffffffffffffffffffffff811661084c576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040517f01ffc9a70000000000000000000000000000000000000000000000000000000081527f3d3ac1b500000000000000000000000000000000000000000000000000000000600482015273ffffffffffffffffffffffffffffffffffffffff8216906301ffc9a790602401602060405180830381865afa1580156108d6573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108fa9190611153565b610930576040517f75b0527a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff821660009081526002602052604090205460ff16156109a8576040517f4e01ccfd00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8316600482015260240161045c565b73ffffffffffffffffffffffffffffffffffffffff821660008181526002602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016600117905590519182527f1f2cd7c97f4d801b5efe26cc409617c1fd6c5ef786e79aacb90af40923e4e8e991016103d6565b60045460609073ffffffffffffffffffffffffffffffffffffffff168015801590610ae757506040517f6b14daf800000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821690636b14daf890610aa490339060009036906004016111be565b602060405180830381865afa158015610ac1573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ae59190611153565b155b15610b1e576040517fef67f5d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000610b2a84866111f7565b60008181526003602052604090205490915073ffffffffffffffffffffffffffffffffffffffff1680610b8c576040517fb151802b0000000000000000000000000000000000000000000000000000000081526004810183905260240161045c565b60055473ffffffffffffffffffffffffffffffffffffffff168015610c36576040517ff1387e1600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82169063f1387e16903490610c03908b908b903390600401611234565b6000604051808303818588803b158015610c1c57600080fd5b505af1158015610c30573d6000803e3d6000fd5b50505050505b6040517f3d3ac1b500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff831690633d3ac1b590610c8c908a908a903390600401611234565b6000604051808303816000875af1158015610cab573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052610cf1919081019061129d565b979650505050505050565b610d04610d97565b6004805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff000000000000000000000000000000000000000083168117909355604080519190921680825260208201939093527f953e92b1a6442e9c3242531154a3f6f6eb00b4e9c719ba8118fa6235e4ce89b691016103d6565b610d8b610d97565b610d9481610e1a565b50565b60005473ffffffffffffffffffffffffffffffffffffffff163314610e18576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161045c565b565b3373ffffffffffffffffffffffffffffffffffffffff821603610e99576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161045c565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60005b83811015610f2a578181015183820152602001610f12565b50506000910152565b60008151808452610f4b816020860160208601610f0f565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000610f906020830184610f33565b9392505050565b73ffffffffffffffffffffffffffffffffffffffff81168114610d9457600080fd5b600060208284031215610fcb57600080fd5b8135610f9081610f97565b60008060008060608587031215610fec57600080fd5b8435935060208501359250604085013567ffffffffffffffff8082111561101257600080fd5b818701915087601f83011261102657600080fd5b81358181111561103557600080fd5b8860208260061b850101111561104a57600080fd5b95989497505060200194505050565b60006020828403121561106b57600080fd5b5035919050565b6000806020838503121561108557600080fd5b823567ffffffffffffffff8082111561109d57600080fd5b818501915085601f8301126110b157600080fd5b8135818111156110c057600080fd5b8660208285010111156110d257600080fd5b60209290920196919550909350505050565b8381526040602080830182905282820184905260009190859060608501845b8781101561114657833561111681610f97565b73ffffffffffffffffffffffffffffffffffffffff16825283830135838301529284019290840190600101611103565b5098975050505050505050565b60006020828403121561116557600080fd5b81518015158114610f9057600080fd5b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b73ffffffffffffffffffffffffffffffffffffffff841681526040602082015260006111ee604083018486611175565b95945050505050565b8035602083101561122e577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff602084900360031b1b165b92915050565b604081526000611248604083018587611175565b905073ffffffffffffffffffffffffffffffffffffffff83166020830152949350505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6000602082840312156112af57600080fd5b815167ffffffffffffffff808211156112c757600080fd5b818401915084601f8301126112db57600080fd5b8151818111156112ed576112ed61126e565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f011681019083821181831017156113335761133361126e565b8160405282815287602084870101111561134c57600080fd5b610cf1836020830160208801610f0f56fea164736f6c6343000810000a",
}

var MercuryVerifierProxyABI = MercuryVerifierProxyMetaData.ABI

var MercuryVerifierProxyBin = MercuryVerifierProxyMetaData.Bin

func DeployMercuryVerifierProxy(auth *bind.TransactOpts, backend bind.ContractBackend, accessController common.Address) (common.Address, *types.Transaction, *MercuryVerifierProxy, error) {
	parsed, err := MercuryVerifierProxyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MercuryVerifierProxyBin), backend, accessController)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MercuryVerifierProxy{MercuryVerifierProxyCaller: MercuryVerifierProxyCaller{contract: contract}, MercuryVerifierProxyTransactor: MercuryVerifierProxyTransactor{contract: contract}, MercuryVerifierProxyFilterer: MercuryVerifierProxyFilterer{contract: contract}}, nil
}

type MercuryVerifierProxy struct {
	address common.Address
	abi     abi.ABI
	MercuryVerifierProxyCaller
	MercuryVerifierProxyTransactor
	MercuryVerifierProxyFilterer
}

type MercuryVerifierProxyCaller struct {
	contract *bind.BoundContract
}

type MercuryVerifierProxyTransactor struct {
	contract *bind.BoundContract
}

type MercuryVerifierProxyFilterer struct {
	contract *bind.BoundContract
}

type MercuryVerifierProxySession struct {
	Contract     *MercuryVerifierProxy
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type MercuryVerifierProxyCallerSession struct {
	Contract *MercuryVerifierProxyCaller
	CallOpts bind.CallOpts
}

type MercuryVerifierProxyTransactorSession struct {
	Contract     *MercuryVerifierProxyTransactor
	TransactOpts bind.TransactOpts
}

type MercuryVerifierProxyRaw struct {
	Contract *MercuryVerifierProxy
}

type MercuryVerifierProxyCallerRaw struct {
	Contract *MercuryVerifierProxyCaller
}

type MercuryVerifierProxyTransactorRaw struct {
	Contract *MercuryVerifierProxyTransactor
}

func NewMercuryVerifierProxy(address common.Address, backend bind.ContractBackend) (*MercuryVerifierProxy, error) {
	abi, err := abi.JSON(strings.NewReader(MercuryVerifierProxyABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindMercuryVerifierProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MercuryVerifierProxy{address: address, abi: abi, MercuryVerifierProxyCaller: MercuryVerifierProxyCaller{contract: contract}, MercuryVerifierProxyTransactor: MercuryVerifierProxyTransactor{contract: contract}, MercuryVerifierProxyFilterer: MercuryVerifierProxyFilterer{contract: contract}}, nil
}

func NewMercuryVerifierProxyCaller(address common.Address, caller bind.ContractCaller) (*MercuryVerifierProxyCaller, error) {
	contract, err := bindMercuryVerifierProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MercuryVerifierProxyCaller{contract: contract}, nil
}

func NewMercuryVerifierProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*MercuryVerifierProxyTransactor, error) {
	contract, err := bindMercuryVerifierProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MercuryVerifierProxyTransactor{contract: contract}, nil
}

func NewMercuryVerifierProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*MercuryVerifierProxyFilterer, error) {
	contract, err := bindMercuryVerifierProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MercuryVerifierProxyFilterer{contract: contract}, nil
}

func bindMercuryVerifierProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MercuryVerifierProxyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_MercuryVerifierProxy *MercuryVerifierProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MercuryVerifierProxy.Contract.MercuryVerifierProxyCaller.contract.Call(opts, result, method, params...)
}

func (_MercuryVerifierProxy *MercuryVerifierProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MercuryVerifierProxy.Contract.MercuryVerifierProxyTransactor.contract.Transfer(opts)
}

func (_MercuryVerifierProxy *MercuryVerifierProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MercuryVerifierProxy.Contract.MercuryVerifierProxyTransactor.contract.Transact(opts, method, params...)
}

func (_MercuryVerifierProxy *MercuryVerifierProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MercuryVerifierProxy.Contract.contract.Call(opts, result, method, params...)
}

func (_MercuryVerifierProxy *MercuryVerifierProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MercuryVerifierProxy.Contract.contract.Transfer(opts)
}

func (_MercuryVerifierProxy *MercuryVerifierProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MercuryVerifierProxy.Contract.contract.Transact(opts, method, params...)
}

func (_MercuryVerifierProxy *MercuryVerifierProxyCaller) GetVerifier(opts *bind.CallOpts, configDigest [32]byte) (common.Address, error) {
	var out []interface{}
	err := _MercuryVerifierProxy.contract.Call(opts, &out, "getVerifier", configDigest)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_MercuryVerifierProxy *MercuryVerifierProxySession) GetVerifier(configDigest [32]byte) (common.Address, error) {
	return _MercuryVerifierProxy.Contract.GetVerifier(&_MercuryVerifierProxy.CallOpts, configDigest)
}

func (_MercuryVerifierProxy *MercuryVerifierProxyCallerSession) GetVerifier(configDigest [32]byte) (common.Address, error) {
	return _MercuryVerifierProxy.Contract.GetVerifier(&_MercuryVerifierProxy.CallOpts, configDigest)
}

func (_MercuryVerifierProxy *MercuryVerifierProxyCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MercuryVerifierProxy.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_MercuryVerifierProxy *MercuryVerifierProxySession) Owner() (common.Address, error) {
	return _MercuryVerifierProxy.Contract.Owner(&_MercuryVerifierProxy.CallOpts)
}

func (_MercuryVerifierProxy *MercuryVerifierProxyCallerSession) Owner() (common.Address, error) {
	return _MercuryVerifierProxy.Contract.Owner(&_MercuryVerifierProxy.CallOpts)
}

func (_MercuryVerifierProxy *MercuryVerifierProxyCaller) SAccessController(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MercuryVerifierProxy.contract.Call(opts, &out, "s_accessController")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_MercuryVerifierProxy *MercuryVerifierProxySession) SAccessController() (common.Address, error) {
	return _MercuryVerifierProxy.Contract.SAccessController(&_MercuryVerifierProxy.CallOpts)
}

func (_MercuryVerifierProxy *MercuryVerifierProxyCallerSession) SAccessController() (common.Address, error) {
	return _MercuryVerifierProxy.Contract.SAccessController(&_MercuryVerifierProxy.CallOpts)
}

func (_MercuryVerifierProxy *MercuryVerifierProxyCaller) SFeeManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MercuryVerifierProxy.contract.Call(opts, &out, "s_feeManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_MercuryVerifierProxy *MercuryVerifierProxySession) SFeeManager() (common.Address, error) {
	return _MercuryVerifierProxy.Contract.SFeeManager(&_MercuryVerifierProxy.CallOpts)
}

func (_MercuryVerifierProxy *MercuryVerifierProxyCallerSession) SFeeManager() (common.Address, error) {
	return _MercuryVerifierProxy.Contract.SFeeManager(&_MercuryVerifierProxy.CallOpts)
}

func (_MercuryVerifierProxy *MercuryVerifierProxyCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MercuryVerifierProxy.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_MercuryVerifierProxy *MercuryVerifierProxySession) TypeAndVersion() (string, error) {
	return _MercuryVerifierProxy.Contract.TypeAndVersion(&_MercuryVerifierProxy.CallOpts)
}

func (_MercuryVerifierProxy *MercuryVerifierProxyCallerSession) TypeAndVersion() (string, error) {
	return _MercuryVerifierProxy.Contract.TypeAndVersion(&_MercuryVerifierProxy.CallOpts)
}

func (_MercuryVerifierProxy *MercuryVerifierProxyTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MercuryVerifierProxy.contract.Transact(opts, "acceptOwnership")
}

func (_MercuryVerifierProxy *MercuryVerifierProxySession) AcceptOwnership() (*types.Transaction, error) {
	return _MercuryVerifierProxy.Contract.AcceptOwnership(&_MercuryVerifierProxy.TransactOpts)
}

func (_MercuryVerifierProxy *MercuryVerifierProxyTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _MercuryVerifierProxy.Contract.AcceptOwnership(&_MercuryVerifierProxy.TransactOpts)
}

func (_MercuryVerifierProxy *MercuryVerifierProxyTransactor) InitializeVerifier(opts *bind.TransactOpts, verifierAddress common.Address) (*types.Transaction, error) {
	return _MercuryVerifierProxy.contract.Transact(opts, "initializeVerifier", verifierAddress)
}

func (_MercuryVerifierProxy *MercuryVerifierProxySession) InitializeVerifier(verifierAddress common.Address) (*types.Transaction, error) {
	return _MercuryVerifierProxy.Contract.InitializeVerifier(&_MercuryVerifierProxy.TransactOpts, verifierAddress)
}

func (_MercuryVerifierProxy *MercuryVerifierProxyTransactorSession) InitializeVerifier(verifierAddress common.Address) (*types.Transaction, error) {
	return _MercuryVerifierProxy.Contract.InitializeVerifier(&_MercuryVerifierProxy.TransactOpts, verifierAddress)
}

func (_MercuryVerifierProxy *MercuryVerifierProxyTransactor) SetAccessController(opts *bind.TransactOpts, accessController common.Address) (*types.Transaction, error) {
	return _MercuryVerifierProxy.contract.Transact(opts, "setAccessController", accessController)
}

func (_MercuryVerifierProxy *MercuryVerifierProxySession) SetAccessController(accessController common.Address) (*types.Transaction, error) {
	return _MercuryVerifierProxy.Contract.SetAccessController(&_MercuryVerifierProxy.TransactOpts, accessController)
}

func (_MercuryVerifierProxy *MercuryVerifierProxyTransactorSession) SetAccessController(accessController common.Address) (*types.Transaction, error) {
	return _MercuryVerifierProxy.Contract.SetAccessController(&_MercuryVerifierProxy.TransactOpts, accessController)
}

func (_MercuryVerifierProxy *MercuryVerifierProxyTransactor) SetFeeManager(opts *bind.TransactOpts, feeManager common.Address) (*types.Transaction, error) {
	return _MercuryVerifierProxy.contract.Transact(opts, "setFeeManager", feeManager)
}

func (_MercuryVerifierProxy *MercuryVerifierProxySession) SetFeeManager(feeManager common.Address) (*types.Transaction, error) {
	return _MercuryVerifierProxy.Contract.SetFeeManager(&_MercuryVerifierProxy.TransactOpts, feeManager)
}

func (_MercuryVerifierProxy *MercuryVerifierProxyTransactorSession) SetFeeManager(feeManager common.Address) (*types.Transaction, error) {
	return _MercuryVerifierProxy.Contract.SetFeeManager(&_MercuryVerifierProxy.TransactOpts, feeManager)
}

func (_MercuryVerifierProxy *MercuryVerifierProxyTransactor) SetVerifier(opts *bind.TransactOpts, currentConfigDigest [32]byte, newConfigDigest [32]byte, addressesAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _MercuryVerifierProxy.contract.Transact(opts, "setVerifier", currentConfigDigest, newConfigDigest, addressesAndWeights)
}

func (_MercuryVerifierProxy *MercuryVerifierProxySession) SetVerifier(currentConfigDigest [32]byte, newConfigDigest [32]byte, addressesAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _MercuryVerifierProxy.Contract.SetVerifier(&_MercuryVerifierProxy.TransactOpts, currentConfigDigest, newConfigDigest, addressesAndWeights)
}

func (_MercuryVerifierProxy *MercuryVerifierProxyTransactorSession) SetVerifier(currentConfigDigest [32]byte, newConfigDigest [32]byte, addressesAndWeights []CommonAddressAndWeight) (*types.Transaction, error) {
	return _MercuryVerifierProxy.Contract.SetVerifier(&_MercuryVerifierProxy.TransactOpts, currentConfigDigest, newConfigDigest, addressesAndWeights)
}

func (_MercuryVerifierProxy *MercuryVerifierProxyTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _MercuryVerifierProxy.contract.Transact(opts, "transferOwnership", to)
}

func (_MercuryVerifierProxy *MercuryVerifierProxySession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _MercuryVerifierProxy.Contract.TransferOwnership(&_MercuryVerifierProxy.TransactOpts, to)
}

func (_MercuryVerifierProxy *MercuryVerifierProxyTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _MercuryVerifierProxy.Contract.TransferOwnership(&_MercuryVerifierProxy.TransactOpts, to)
}

func (_MercuryVerifierProxy *MercuryVerifierProxyTransactor) UnsetVerifier(opts *bind.TransactOpts, configDigest [32]byte) (*types.Transaction, error) {
	return _MercuryVerifierProxy.contract.Transact(opts, "unsetVerifier", configDigest)
}

func (_MercuryVerifierProxy *MercuryVerifierProxySession) UnsetVerifier(configDigest [32]byte) (*types.Transaction, error) {
	return _MercuryVerifierProxy.Contract.UnsetVerifier(&_MercuryVerifierProxy.TransactOpts, configDigest)
}

func (_MercuryVerifierProxy *MercuryVerifierProxyTransactorSession) UnsetVerifier(configDigest [32]byte) (*types.Transaction, error) {
	return _MercuryVerifierProxy.Contract.UnsetVerifier(&_MercuryVerifierProxy.TransactOpts, configDigest)
}

func (_MercuryVerifierProxy *MercuryVerifierProxyTransactor) Verify(opts *bind.TransactOpts, payload []byte) (*types.Transaction, error) {
	return _MercuryVerifierProxy.contract.Transact(opts, "verify", payload)
}

func (_MercuryVerifierProxy *MercuryVerifierProxySession) Verify(payload []byte) (*types.Transaction, error) {
	return _MercuryVerifierProxy.Contract.Verify(&_MercuryVerifierProxy.TransactOpts, payload)
}

func (_MercuryVerifierProxy *MercuryVerifierProxyTransactorSession) Verify(payload []byte) (*types.Transaction, error) {
	return _MercuryVerifierProxy.Contract.Verify(&_MercuryVerifierProxy.TransactOpts, payload)
}

type MercuryVerifierProxyAccessControllerSetIterator struct {
	Event *MercuryVerifierProxyAccessControllerSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MercuryVerifierProxyAccessControllerSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MercuryVerifierProxyAccessControllerSet)
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
		it.Event = new(MercuryVerifierProxyAccessControllerSet)
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

func (it *MercuryVerifierProxyAccessControllerSetIterator) Error() error {
	return it.fail
}

func (it *MercuryVerifierProxyAccessControllerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MercuryVerifierProxyAccessControllerSet struct {
	OldAccessController common.Address
	NewAccessController common.Address
	Raw                 types.Log
}

func (_MercuryVerifierProxy *MercuryVerifierProxyFilterer) FilterAccessControllerSet(opts *bind.FilterOpts) (*MercuryVerifierProxyAccessControllerSetIterator, error) {

	logs, sub, err := _MercuryVerifierProxy.contract.FilterLogs(opts, "AccessControllerSet")
	if err != nil {
		return nil, err
	}
	return &MercuryVerifierProxyAccessControllerSetIterator{contract: _MercuryVerifierProxy.contract, event: "AccessControllerSet", logs: logs, sub: sub}, nil
}

func (_MercuryVerifierProxy *MercuryVerifierProxyFilterer) WatchAccessControllerSet(opts *bind.WatchOpts, sink chan<- *MercuryVerifierProxyAccessControllerSet) (event.Subscription, error) {

	logs, sub, err := _MercuryVerifierProxy.contract.WatchLogs(opts, "AccessControllerSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MercuryVerifierProxyAccessControllerSet)
				if err := _MercuryVerifierProxy.contract.UnpackLog(event, "AccessControllerSet", log); err != nil {
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

func (_MercuryVerifierProxy *MercuryVerifierProxyFilterer) ParseAccessControllerSet(log types.Log) (*MercuryVerifierProxyAccessControllerSet, error) {
	event := new(MercuryVerifierProxyAccessControllerSet)
	if err := _MercuryVerifierProxy.contract.UnpackLog(event, "AccessControllerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MercuryVerifierProxyFeeManagerSetIterator struct {
	Event *MercuryVerifierProxyFeeManagerSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MercuryVerifierProxyFeeManagerSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MercuryVerifierProxyFeeManagerSet)
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
		it.Event = new(MercuryVerifierProxyFeeManagerSet)
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

func (it *MercuryVerifierProxyFeeManagerSetIterator) Error() error {
	return it.fail
}

func (it *MercuryVerifierProxyFeeManagerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MercuryVerifierProxyFeeManagerSet struct {
	OldFeeManager common.Address
	NewFeeManager common.Address
	Raw           types.Log
}

func (_MercuryVerifierProxy *MercuryVerifierProxyFilterer) FilterFeeManagerSet(opts *bind.FilterOpts) (*MercuryVerifierProxyFeeManagerSetIterator, error) {

	logs, sub, err := _MercuryVerifierProxy.contract.FilterLogs(opts, "FeeManagerSet")
	if err != nil {
		return nil, err
	}
	return &MercuryVerifierProxyFeeManagerSetIterator{contract: _MercuryVerifierProxy.contract, event: "FeeManagerSet", logs: logs, sub: sub}, nil
}

func (_MercuryVerifierProxy *MercuryVerifierProxyFilterer) WatchFeeManagerSet(opts *bind.WatchOpts, sink chan<- *MercuryVerifierProxyFeeManagerSet) (event.Subscription, error) {

	logs, sub, err := _MercuryVerifierProxy.contract.WatchLogs(opts, "FeeManagerSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MercuryVerifierProxyFeeManagerSet)
				if err := _MercuryVerifierProxy.contract.UnpackLog(event, "FeeManagerSet", log); err != nil {
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

func (_MercuryVerifierProxy *MercuryVerifierProxyFilterer) ParseFeeManagerSet(log types.Log) (*MercuryVerifierProxyFeeManagerSet, error) {
	event := new(MercuryVerifierProxyFeeManagerSet)
	if err := _MercuryVerifierProxy.contract.UnpackLog(event, "FeeManagerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MercuryVerifierProxyOwnershipTransferRequestedIterator struct {
	Event *MercuryVerifierProxyOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MercuryVerifierProxyOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MercuryVerifierProxyOwnershipTransferRequested)
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
		it.Event = new(MercuryVerifierProxyOwnershipTransferRequested)
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

func (it *MercuryVerifierProxyOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *MercuryVerifierProxyOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MercuryVerifierProxyOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_MercuryVerifierProxy *MercuryVerifierProxyFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MercuryVerifierProxyOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MercuryVerifierProxy.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MercuryVerifierProxyOwnershipTransferRequestedIterator{contract: _MercuryVerifierProxy.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_MercuryVerifierProxy *MercuryVerifierProxyFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *MercuryVerifierProxyOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MercuryVerifierProxy.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MercuryVerifierProxyOwnershipTransferRequested)
				if err := _MercuryVerifierProxy.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_MercuryVerifierProxy *MercuryVerifierProxyFilterer) ParseOwnershipTransferRequested(log types.Log) (*MercuryVerifierProxyOwnershipTransferRequested, error) {
	event := new(MercuryVerifierProxyOwnershipTransferRequested)
	if err := _MercuryVerifierProxy.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MercuryVerifierProxyOwnershipTransferredIterator struct {
	Event *MercuryVerifierProxyOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MercuryVerifierProxyOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MercuryVerifierProxyOwnershipTransferred)
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
		it.Event = new(MercuryVerifierProxyOwnershipTransferred)
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

func (it *MercuryVerifierProxyOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *MercuryVerifierProxyOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MercuryVerifierProxyOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_MercuryVerifierProxy *MercuryVerifierProxyFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MercuryVerifierProxyOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MercuryVerifierProxy.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MercuryVerifierProxyOwnershipTransferredIterator{contract: _MercuryVerifierProxy.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_MercuryVerifierProxy *MercuryVerifierProxyFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MercuryVerifierProxyOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MercuryVerifierProxy.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MercuryVerifierProxyOwnershipTransferred)
				if err := _MercuryVerifierProxy.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_MercuryVerifierProxy *MercuryVerifierProxyFilterer) ParseOwnershipTransferred(log types.Log) (*MercuryVerifierProxyOwnershipTransferred, error) {
	event := new(MercuryVerifierProxyOwnershipTransferred)
	if err := _MercuryVerifierProxy.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MercuryVerifierProxyVerifierInitializedIterator struct {
	Event *MercuryVerifierProxyVerifierInitialized

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MercuryVerifierProxyVerifierInitializedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MercuryVerifierProxyVerifierInitialized)
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
		it.Event = new(MercuryVerifierProxyVerifierInitialized)
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

func (it *MercuryVerifierProxyVerifierInitializedIterator) Error() error {
	return it.fail
}

func (it *MercuryVerifierProxyVerifierInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MercuryVerifierProxyVerifierInitialized struct {
	VerifierAddress common.Address
	Raw             types.Log
}

func (_MercuryVerifierProxy *MercuryVerifierProxyFilterer) FilterVerifierInitialized(opts *bind.FilterOpts) (*MercuryVerifierProxyVerifierInitializedIterator, error) {

	logs, sub, err := _MercuryVerifierProxy.contract.FilterLogs(opts, "VerifierInitialized")
	if err != nil {
		return nil, err
	}
	return &MercuryVerifierProxyVerifierInitializedIterator{contract: _MercuryVerifierProxy.contract, event: "VerifierInitialized", logs: logs, sub: sub}, nil
}

func (_MercuryVerifierProxy *MercuryVerifierProxyFilterer) WatchVerifierInitialized(opts *bind.WatchOpts, sink chan<- *MercuryVerifierProxyVerifierInitialized) (event.Subscription, error) {

	logs, sub, err := _MercuryVerifierProxy.contract.WatchLogs(opts, "VerifierInitialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MercuryVerifierProxyVerifierInitialized)
				if err := _MercuryVerifierProxy.contract.UnpackLog(event, "VerifierInitialized", log); err != nil {
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

func (_MercuryVerifierProxy *MercuryVerifierProxyFilterer) ParseVerifierInitialized(log types.Log) (*MercuryVerifierProxyVerifierInitialized, error) {
	event := new(MercuryVerifierProxyVerifierInitialized)
	if err := _MercuryVerifierProxy.contract.UnpackLog(event, "VerifierInitialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MercuryVerifierProxyVerifierSetIterator struct {
	Event *MercuryVerifierProxyVerifierSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MercuryVerifierProxyVerifierSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MercuryVerifierProxyVerifierSet)
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
		it.Event = new(MercuryVerifierProxyVerifierSet)
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

func (it *MercuryVerifierProxyVerifierSetIterator) Error() error {
	return it.fail
}

func (it *MercuryVerifierProxyVerifierSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MercuryVerifierProxyVerifierSet struct {
	OldConfigDigest [32]byte
	NewConfigDigest [32]byte
	VerifierAddress common.Address
	Raw             types.Log
}

func (_MercuryVerifierProxy *MercuryVerifierProxyFilterer) FilterVerifierSet(opts *bind.FilterOpts) (*MercuryVerifierProxyVerifierSetIterator, error) {

	logs, sub, err := _MercuryVerifierProxy.contract.FilterLogs(opts, "VerifierSet")
	if err != nil {
		return nil, err
	}
	return &MercuryVerifierProxyVerifierSetIterator{contract: _MercuryVerifierProxy.contract, event: "VerifierSet", logs: logs, sub: sub}, nil
}

func (_MercuryVerifierProxy *MercuryVerifierProxyFilterer) WatchVerifierSet(opts *bind.WatchOpts, sink chan<- *MercuryVerifierProxyVerifierSet) (event.Subscription, error) {

	logs, sub, err := _MercuryVerifierProxy.contract.WatchLogs(opts, "VerifierSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MercuryVerifierProxyVerifierSet)
				if err := _MercuryVerifierProxy.contract.UnpackLog(event, "VerifierSet", log); err != nil {
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

func (_MercuryVerifierProxy *MercuryVerifierProxyFilterer) ParseVerifierSet(log types.Log) (*MercuryVerifierProxyVerifierSet, error) {
	event := new(MercuryVerifierProxyVerifierSet)
	if err := _MercuryVerifierProxy.contract.UnpackLog(event, "VerifierSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MercuryVerifierProxyVerifierUnsetIterator struct {
	Event *MercuryVerifierProxyVerifierUnset

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MercuryVerifierProxyVerifierUnsetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MercuryVerifierProxyVerifierUnset)
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
		it.Event = new(MercuryVerifierProxyVerifierUnset)
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

func (it *MercuryVerifierProxyVerifierUnsetIterator) Error() error {
	return it.fail
}

func (it *MercuryVerifierProxyVerifierUnsetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MercuryVerifierProxyVerifierUnset struct {
	ConfigDigest    [32]byte
	VerifierAddress common.Address
	Raw             types.Log
}

func (_MercuryVerifierProxy *MercuryVerifierProxyFilterer) FilterVerifierUnset(opts *bind.FilterOpts) (*MercuryVerifierProxyVerifierUnsetIterator, error) {

	logs, sub, err := _MercuryVerifierProxy.contract.FilterLogs(opts, "VerifierUnset")
	if err != nil {
		return nil, err
	}
	return &MercuryVerifierProxyVerifierUnsetIterator{contract: _MercuryVerifierProxy.contract, event: "VerifierUnset", logs: logs, sub: sub}, nil
}

func (_MercuryVerifierProxy *MercuryVerifierProxyFilterer) WatchVerifierUnset(opts *bind.WatchOpts, sink chan<- *MercuryVerifierProxyVerifierUnset) (event.Subscription, error) {

	logs, sub, err := _MercuryVerifierProxy.contract.WatchLogs(opts, "VerifierUnset")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MercuryVerifierProxyVerifierUnset)
				if err := _MercuryVerifierProxy.contract.UnpackLog(event, "VerifierUnset", log); err != nil {
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

func (_MercuryVerifierProxy *MercuryVerifierProxyFilterer) ParseVerifierUnset(log types.Log) (*MercuryVerifierProxyVerifierUnset, error) {
	event := new(MercuryVerifierProxyVerifierUnset)
	if err := _MercuryVerifierProxy.contract.UnpackLog(event, "VerifierUnset", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_MercuryVerifierProxy *MercuryVerifierProxy) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _MercuryVerifierProxy.abi.Events["AccessControllerSet"].ID:
		return _MercuryVerifierProxy.ParseAccessControllerSet(log)
	case _MercuryVerifierProxy.abi.Events["FeeManagerSet"].ID:
		return _MercuryVerifierProxy.ParseFeeManagerSet(log)
	case _MercuryVerifierProxy.abi.Events["OwnershipTransferRequested"].ID:
		return _MercuryVerifierProxy.ParseOwnershipTransferRequested(log)
	case _MercuryVerifierProxy.abi.Events["OwnershipTransferred"].ID:
		return _MercuryVerifierProxy.ParseOwnershipTransferred(log)
	case _MercuryVerifierProxy.abi.Events["VerifierInitialized"].ID:
		return _MercuryVerifierProxy.ParseVerifierInitialized(log)
	case _MercuryVerifierProxy.abi.Events["VerifierSet"].ID:
		return _MercuryVerifierProxy.ParseVerifierSet(log)
	case _MercuryVerifierProxy.abi.Events["VerifierUnset"].ID:
		return _MercuryVerifierProxy.ParseVerifierUnset(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (MercuryVerifierProxyAccessControllerSet) Topic() common.Hash {
	return common.HexToHash("0x953e92b1a6442e9c3242531154a3f6f6eb00b4e9c719ba8118fa6235e4ce89b6")
}

func (MercuryVerifierProxyFeeManagerSet) Topic() common.Hash {
	return common.HexToHash("0x04628abcaa6b1674651352125cb94b65b289145bc2bc4d67720bb7d966372f03")
}

func (MercuryVerifierProxyOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (MercuryVerifierProxyOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (MercuryVerifierProxyVerifierInitialized) Topic() common.Hash {
	return common.HexToHash("0x1f2cd7c97f4d801b5efe26cc409617c1fd6c5ef786e79aacb90af40923e4e8e9")
}

func (MercuryVerifierProxyVerifierSet) Topic() common.Hash {
	return common.HexToHash("0xbeb513e532542a562ac35699e7cd9ae7d198dcd3eee15bada6c857d28ceaddcf")
}

func (MercuryVerifierProxyVerifierUnset) Topic() common.Hash {
	return common.HexToHash("0x11dc15c4b8ac2b183166cc8427e5385a5ece8308217a4217338c6a7614845c4c")
}

func (_MercuryVerifierProxy *MercuryVerifierProxy) Address() common.Address {
	return _MercuryVerifierProxy.address
}

type MercuryVerifierProxyInterface interface {
	GetVerifier(opts *bind.CallOpts, configDigest [32]byte) (common.Address, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SAccessController(opts *bind.CallOpts) (common.Address, error)

	SFeeManager(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	InitializeVerifier(opts *bind.TransactOpts, verifierAddress common.Address) (*types.Transaction, error)

	SetAccessController(opts *bind.TransactOpts, accessController common.Address) (*types.Transaction, error)

	SetFeeManager(opts *bind.TransactOpts, feeManager common.Address) (*types.Transaction, error)

	SetVerifier(opts *bind.TransactOpts, currentConfigDigest [32]byte, newConfigDigest [32]byte, addressesAndWeights []CommonAddressAndWeight) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	UnsetVerifier(opts *bind.TransactOpts, configDigest [32]byte) (*types.Transaction, error)

	Verify(opts *bind.TransactOpts, payload []byte) (*types.Transaction, error)

	FilterAccessControllerSet(opts *bind.FilterOpts) (*MercuryVerifierProxyAccessControllerSetIterator, error)

	WatchAccessControllerSet(opts *bind.WatchOpts, sink chan<- *MercuryVerifierProxyAccessControllerSet) (event.Subscription, error)

	ParseAccessControllerSet(log types.Log) (*MercuryVerifierProxyAccessControllerSet, error)

	FilterFeeManagerSet(opts *bind.FilterOpts) (*MercuryVerifierProxyFeeManagerSetIterator, error)

	WatchFeeManagerSet(opts *bind.WatchOpts, sink chan<- *MercuryVerifierProxyFeeManagerSet) (event.Subscription, error)

	ParseFeeManagerSet(log types.Log) (*MercuryVerifierProxyFeeManagerSet, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MercuryVerifierProxyOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *MercuryVerifierProxyOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*MercuryVerifierProxyOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MercuryVerifierProxyOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MercuryVerifierProxyOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*MercuryVerifierProxyOwnershipTransferred, error)

	FilterVerifierInitialized(opts *bind.FilterOpts) (*MercuryVerifierProxyVerifierInitializedIterator, error)

	WatchVerifierInitialized(opts *bind.WatchOpts, sink chan<- *MercuryVerifierProxyVerifierInitialized) (event.Subscription, error)

	ParseVerifierInitialized(log types.Log) (*MercuryVerifierProxyVerifierInitialized, error)

	FilterVerifierSet(opts *bind.FilterOpts) (*MercuryVerifierProxyVerifierSetIterator, error)

	WatchVerifierSet(opts *bind.WatchOpts, sink chan<- *MercuryVerifierProxyVerifierSet) (event.Subscription, error)

	ParseVerifierSet(log types.Log) (*MercuryVerifierProxyVerifierSet, error)

	FilterVerifierUnset(opts *bind.FilterOpts) (*MercuryVerifierProxyVerifierUnsetIterator, error)

	WatchVerifierUnset(opts *bind.WatchOpts, sink chan<- *MercuryVerifierProxyVerifierUnset) (event.Subscription, error)

	ParseVerifierUnset(log types.Log) (*MercuryVerifierProxyVerifierUnset, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
