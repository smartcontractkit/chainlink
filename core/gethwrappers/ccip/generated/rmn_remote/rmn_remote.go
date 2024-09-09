// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package rmn_remote

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

type IRMNV2Signature struct {
	R [32]byte
	S [32]byte
}

type InternalMerkleRoot struct {
	SourceChainSelector uint64
	OnRampAddress       []byte
	MinSeqNr            uint64
	MaxSeqNr            uint64
	MerkleRoot          [32]byte
}

type RMNRemoteConfig struct {
	RmnHomeContractConfigDigest [32]byte
	Signers                     []RMNRemoteSigner
	MinSigners                  uint64
}

type RMNRemoteSigner struct {
	OnchainPublicKey common.Address
	NodeIndex        uint64
}

type RMNRemoteVersionedConfig struct {
	Version uint32
	Config  RMNRemoteConfig
}

var RMNRemoteMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ConfigNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateOnchainPublicKey\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSignature\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSignerOrder\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MinSignersTooHigh\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OutOfOrderSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ThresholdNotMet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnexpectedSigner\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"rmnHomeContractConfigDigest\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"onchainPublicKey\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"nodeIndex\",\"type\":\"uint64\"}],\"internalType\":\"structRMNRemote.Signer[]\",\"name\":\"signers\",\"type\":\"tuple[]\"},{\"internalType\":\"uint64\",\"name\":\"minSigners\",\"type\":\"uint64\"}],\"internalType\":\"structRMNRemote.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"indexed\":false,\"internalType\":\"structRMNRemote.VersionedConfig\",\"name\":\"versionedConfig\",\"type\":\"tuple\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getVersionedConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"rmnHomeContractConfigDigest\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"onchainPublicKey\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"nodeIndex\",\"type\":\"uint64\"}],\"internalType\":\"structRMNRemote.Signer[]\",\"name\":\"signers\",\"type\":\"tuple[]\"},{\"internalType\":\"uint64\",\"name\":\"minSigners\",\"type\":\"uint64\"}],\"internalType\":\"structRMNRemote.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"internalType\":\"structRMNRemote.VersionedConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"subject\",\"type\":\"bytes16\"}],\"name\":\"isCursed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isCursed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"rmnHomeContractConfigDigest\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"onchainPublicKey\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"nodeIndex\",\"type\":\"uint64\"}],\"internalType\":\"structRMNRemote.Signer[]\",\"name\":\"signers\",\"type\":\"tuple[]\"},{\"internalType\":\"uint64\",\"name\":\"minSigners\",\"type\":\"uint64\"}],\"internalType\":\"structRMNRemote.Config\",\"name\":\"newConfig\",\"type\":\"tuple\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"test\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRampAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"maxSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"merkleRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structInternal.MerkleRoot[]\",\"name\":\"destLaneUpdates\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structIRMNV2.Signature[]\",\"name\":\"signatures\",\"type\":\"tuple[]\"}],\"name\":\"verify\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b50604051620018293803806200182983398101604081905262000034916200017e565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000d3565b5050506001600160401b0316608052620001b0565b336001600160a01b038216036200012d5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000602082840312156200019157600080fd5b81516001600160401b0381168114620001a957600080fd5b9392505050565b608051611660620001c9600039600050506116606000f3fe608060405234801561001057600080fd5b50600436106100be5760003560e01c806379ba509711610076578063ccd5af001161005b578063ccd5af001461019a578063f2fde38b146101ac578063f8a8fd6d1461012857600080fd5b806379ba50971461016a5780638da5cb5b1461017257600080fd5b80631add205f116100a75780631add205f1461012a5780632cbc26bb1461013f578063397796f71461016357600080fd5b8063181f5a77146100c3578063198f0f7714610115575b600080fd5b6100ff6040518060400160405280601381526020017f524d4e52656d6f746520312e362e302d6465760000000000000000000000000081525081565b60405161010c9190610c4b565b60405180910390f35b610128610123366004610c65565b6101bf565b005b6101326105f0565b60405161010c9190610ca0565b61015361014d366004610d58565b50600090565b604051901515815260200161010c565b6000610153565b6101286106d6565b60005460405173ffffffffffffffffffffffffffffffffffffffff909116815260200161010c565b6101286101a8366004610f60565b5050565b6101286101ba366004611137565b610a0d565b6101c7610a23565b60015b6101d76020830183611154565b90508110156102a7576101ed6020830183611154565b828181106101fd576101fd6111c3565b905060400201602001602081019061021591906111f2565b67ffffffffffffffff1661022c6020840184611154565b61023760018561123e565b818110610246576102466111c3565b905060400201602001602081019061025e91906111f2565b67ffffffffffffffff161061029f576040517f4485151700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001016101ca565b506102b56020820182611154565b90506102c760608301604084016111f2565b67ffffffffffffffff161115610309576040517ffba0d9e600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60025b60018101541561040557600180820180546006926000929161032e919061123e565b8154811061033e5761033e6111c3565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff168352820192909252604001902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055600181018054806103a8576103a8611257565b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffff0000000000000000000000000000000000000000000000000000000016905501905561030c565b5060005b6104166020830183611154565b905081101561054b57600660006104306020850185611154565b84818110610440576104406111c3565b6104569260206040909202019081019150611137565b73ffffffffffffffffffffffffffffffffffffffff16815260208101919091526040016000205460ff16156104b7576040517f28cae27d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600660006104ca6020860186611154565b858181106104da576104da6111c3565b6104f09260206040909202019081019150611137565b73ffffffffffffffffffffffffffffffffffffffff168152602081019190915260400160002080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055600101610409565b50806002610559828261133f565b5050600580546000919082906105749063ffffffff1661147a565b91906101000a81548163ffffffff021916908363ffffffff160217905590507f6cc65868ae41a007e6c3ed18ce591c123dd4e5864b421888c68ce92dae98cea460405180604001604052808363ffffffff168152602001846105d59061149d565b90526040516105e49190610ca0565b60405180910390a15050565b6105f8610b99565b60408051808201825260055463ffffffff1681528151606081018352600280548252600380548551602082810282018101909752818152949580870195858201939092909160009084015b828210156106b1576000848152602090819020604080518082019091529084015473ffffffffffffffffffffffffffffffffffffffff8116825274010000000000000000000000000000000000000000900467ffffffffffffffff1681830152825260019092019101610643565b505050908252506002919091015467ffffffffffffffff166020909101529052919050565b60015473ffffffffffffffffffffffffffffffffffffffff16331461075c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b84518110156109bf5760008582815181106107f5576107f56111c3565b602002602001015190506000600186601b8460000151856020015160405160008152602001604052604051610846949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa158015610868573d6000803e3d6000fd5b50506040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0015191505073ffffffffffffffffffffffffffffffffffffffff81166108e0576040517f8baa579f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff1610610945576040517fbbe15e7f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff811660009081526006602052604090205460ff166109a4576040517faaaa914100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b9250826109b08561161b565b945050508060010190506107d8565b5060045467ffffffffffffffff16821015610a06576040517f59fa4a9300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050505050565b610a15610a23565b610a1e81610aa4565b50565b565b60005473ffffffffffffffffffffffffffffffffffffffff163314610a21576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610753565b3373ffffffffffffffffffffffffffffffffffffffff821603610b23576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610753565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6040518060400160405280600063ffffffff168152602001610be260405180606001604052806000801916815260200160608152602001600067ffffffffffffffff1681525090565b905290565b6000815180845260005b81811015610c0d57602081850181015186830182015201610bf1565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b602081526000610c5e6020830184610be7565b9392505050565b600060208284031215610c7757600080fd5b813567ffffffffffffffff811115610c8e57600080fd5b820160608185031215610c5e57600080fd5b6000602080835263ffffffff8451168184015280840151604080604086015260c0850182516060870152838301516060608088015281815180845260e0890191508683019350600092505b80831015610d34578351805173ffffffffffffffffffffffffffffffffffffffff16835287015167ffffffffffffffff1687830152928601926001929092019190840190610ceb565b50604085015167ffffffffffffffff811660a08a0152955098975050505050505050565b600060208284031215610d6a57600080fd5b81357fffffffffffffffffffffffffffffffff0000000000000000000000000000000081168114610c5e57600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040805190810167ffffffffffffffff81118282101715610dec57610dec610d9a565b60405290565b60405160a0810167ffffffffffffffff81118282101715610dec57610dec610d9a565b6040516060810167ffffffffffffffff81118282101715610dec57610dec610d9a565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610e7f57610e7f610d9a565b604052919050565b600067ffffffffffffffff821115610ea157610ea1610d9a565b5060051b60200190565b67ffffffffffffffff81168114610a1e57600080fd5b8035610ecc81610eab565b919050565b600082601f830112610ee257600080fd5b81356020610ef7610ef283610e87565b610e38565b82815260069290921b84018101918181019086841115610f1657600080fd5b8286015b84811015610f555760408189031215610f335760008081fd5b610f3b610dc9565b813581528482013585820152835291830191604001610f1a565b509695505050505050565b60008060408385031215610f7357600080fd5b823567ffffffffffffffff80821115610f8b57600080fd5b818501915085601f830112610f9f57600080fd5b81356020610faf610ef283610e87565b82815260059290921b84018101918181019089841115610fce57600080fd5b8286015b848110156110e757803586811115610fe957600080fd5b87017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe060a0828e038201121561101e57600080fd5b611026610df2565b8683013561103381610eab565b815260408301358981111561104757600080fd5b8301603f81018f1361105857600080fd5b878101358a81111561106c5761106c610d9a565b61107c8985601f84011601610e38565b93508084528f604082840101111561109357600080fd5b80604083018a8601376000898286010152505081878201526110b760608401610ec1565b60408201526110c860808401610ec1565b606082015260a092909201356080830152508352918301918301610fd2565b50965050860135925050808211156110fe57600080fd5b5061110b85828601610ed1565b9150509250929050565b73ffffffffffffffffffffffffffffffffffffffff81168114610a1e57600080fd5b60006020828403121561114957600080fd5b8135610c5e81611115565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261118957600080fd5b83018035915067ffffffffffffffff8211156111a457600080fd5b6020019150600681901b36038213156111bc57600080fd5b9250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60006020828403121561120457600080fd5b8135610c5e81610eab565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b818103818111156112515761125161120f565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b6000813561125181610eab565b813561129e81611115565b73ffffffffffffffffffffffffffffffffffffffff811690508154817fffffffffffffffffffffffff0000000000000000000000000000000000000000821617835560208401356112ee81610eab565b7bffffffffffffffff00000000000000000000000000000000000000008160a01b16837fffffffff000000000000000000000000000000000000000000000000000000008416171784555050505050565b81358155600180820160208401357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe185360301811261137d57600080fd5b8401803567ffffffffffffffff81111561139657600080fd5b6020820191508060061b36038213156113ae57600080fd5b680100000000000000008111156113c7576113c7610d9a565b8254818455808210156113fc576000848152602081208381019083015b808210156113f857828255908701906113e4565b5050505b50600092835260208320925b8181101561142c5761141a8385611293565b92840192604092909201918401611408565b50505050506101a861144060408401611286565b6002830167ffffffffffffffff82167fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000008254161781555050565b600063ffffffff8083168181036114935761149361120f565b6001019392505050565b6000606082360312156114af57600080fd5b6114b7610e15565b8235815260208084013567ffffffffffffffff8111156114d657600080fd5b840136601f8201126114e757600080fd5b80356114f5610ef282610e87565b81815260069190911b8201830190838101903683111561151457600080fd5b928401925b8284101561156a57604084360312156115325760008081fd5b61153a610dc9565b843561154581611115565b81528486013561155481610eab565b8187015282526040939093019290840190611519565b8085870152505050505061158060408401610ec1565b604082015292915050565b81811015611608577ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffee08d8a030184528751868151168a528a810151848c8c01526115d7858c0182610be7565b828e015189168c8f01528983015189168a8d0152918701519a87019a909a529850968901969289019260010161158b565b50969d9c50505050505050505050505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361164c5761164c61120f565b506001019056fea164736f6c6343000818000a",
}

var RMNRemoteABI = RMNRemoteMetaData.ABI

var RMNRemoteBin = RMNRemoteMetaData.Bin

func DeployRMNRemote(auth *bind.TransactOpts, backend bind.ContractBackend, chainSelector uint64) (common.Address, *types.Transaction, *RMNRemote, error) {
	parsed, err := RMNRemoteMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(RMNRemoteBin), backend, chainSelector)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &RMNRemote{address: address, abi: *parsed, RMNRemoteCaller: RMNRemoteCaller{contract: contract}, RMNRemoteTransactor: RMNRemoteTransactor{contract: contract}, RMNRemoteFilterer: RMNRemoteFilterer{contract: contract}}, nil
}

type RMNRemote struct {
	address common.Address
	abi     abi.ABI
	RMNRemoteCaller
	RMNRemoteTransactor
	RMNRemoteFilterer
}

type RMNRemoteCaller struct {
	contract *bind.BoundContract
}

type RMNRemoteTransactor struct {
	contract *bind.BoundContract
}

type RMNRemoteFilterer struct {
	contract *bind.BoundContract
}

type RMNRemoteSession struct {
	Contract     *RMNRemote
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type RMNRemoteCallerSession struct {
	Contract *RMNRemoteCaller
	CallOpts bind.CallOpts
}

type RMNRemoteTransactorSession struct {
	Contract     *RMNRemoteTransactor
	TransactOpts bind.TransactOpts
}

type RMNRemoteRaw struct {
	Contract *RMNRemote
}

type RMNRemoteCallerRaw struct {
	Contract *RMNRemoteCaller
}

type RMNRemoteTransactorRaw struct {
	Contract *RMNRemoteTransactor
}

func NewRMNRemote(address common.Address, backend bind.ContractBackend) (*RMNRemote, error) {
	abi, err := abi.JSON(strings.NewReader(RMNRemoteABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindRMNRemote(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RMNRemote{address: address, abi: abi, RMNRemoteCaller: RMNRemoteCaller{contract: contract}, RMNRemoteTransactor: RMNRemoteTransactor{contract: contract}, RMNRemoteFilterer: RMNRemoteFilterer{contract: contract}}, nil
}

func NewRMNRemoteCaller(address common.Address, caller bind.ContractCaller) (*RMNRemoteCaller, error) {
	contract, err := bindRMNRemote(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RMNRemoteCaller{contract: contract}, nil
}

func NewRMNRemoteTransactor(address common.Address, transactor bind.ContractTransactor) (*RMNRemoteTransactor, error) {
	contract, err := bindRMNRemote(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RMNRemoteTransactor{contract: contract}, nil
}

func NewRMNRemoteFilterer(address common.Address, filterer bind.ContractFilterer) (*RMNRemoteFilterer, error) {
	contract, err := bindRMNRemote(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RMNRemoteFilterer{contract: contract}, nil
}

func bindRMNRemote(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RMNRemoteMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_RMNRemote *RMNRemoteRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RMNRemote.Contract.RMNRemoteCaller.contract.Call(opts, result, method, params...)
}

func (_RMNRemote *RMNRemoteRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RMNRemote.Contract.RMNRemoteTransactor.contract.Transfer(opts)
}

func (_RMNRemote *RMNRemoteRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RMNRemote.Contract.RMNRemoteTransactor.contract.Transact(opts, method, params...)
}

func (_RMNRemote *RMNRemoteCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RMNRemote.Contract.contract.Call(opts, result, method, params...)
}

func (_RMNRemote *RMNRemoteTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RMNRemote.Contract.contract.Transfer(opts)
}

func (_RMNRemote *RMNRemoteTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RMNRemote.Contract.contract.Transact(opts, method, params...)
}

func (_RMNRemote *RMNRemoteCaller) GetVersionedConfig(opts *bind.CallOpts) (RMNRemoteVersionedConfig, error) {
	var out []interface{}
	err := _RMNRemote.contract.Call(opts, &out, "getVersionedConfig")

	if err != nil {
		return *new(RMNRemoteVersionedConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(RMNRemoteVersionedConfig)).(*RMNRemoteVersionedConfig)

	return out0, err

}

func (_RMNRemote *RMNRemoteSession) GetVersionedConfig() (RMNRemoteVersionedConfig, error) {
	return _RMNRemote.Contract.GetVersionedConfig(&_RMNRemote.CallOpts)
}

func (_RMNRemote *RMNRemoteCallerSession) GetVersionedConfig() (RMNRemoteVersionedConfig, error) {
	return _RMNRemote.Contract.GetVersionedConfig(&_RMNRemote.CallOpts)
}

func (_RMNRemote *RMNRemoteCaller) IsCursed(opts *bind.CallOpts, subject [16]byte) (bool, error) {
	var out []interface{}
	err := _RMNRemote.contract.Call(opts, &out, "isCursed", subject)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_RMNRemote *RMNRemoteSession) IsCursed(subject [16]byte) (bool, error) {
	return _RMNRemote.Contract.IsCursed(&_RMNRemote.CallOpts, subject)
}

func (_RMNRemote *RMNRemoteCallerSession) IsCursed(subject [16]byte) (bool, error) {
	return _RMNRemote.Contract.IsCursed(&_RMNRemote.CallOpts, subject)
}

func (_RMNRemote *RMNRemoteCaller) IsCursed0(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _RMNRemote.contract.Call(opts, &out, "isCursed0")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_RMNRemote *RMNRemoteSession) IsCursed0() (bool, error) {
	return _RMNRemote.Contract.IsCursed0(&_RMNRemote.CallOpts)
}

func (_RMNRemote *RMNRemoteCallerSession) IsCursed0() (bool, error) {
	return _RMNRemote.Contract.IsCursed0(&_RMNRemote.CallOpts)
}

func (_RMNRemote *RMNRemoteCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RMNRemote.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_RMNRemote *RMNRemoteSession) Owner() (common.Address, error) {
	return _RMNRemote.Contract.Owner(&_RMNRemote.CallOpts)
}

func (_RMNRemote *RMNRemoteCallerSession) Owner() (common.Address, error) {
	return _RMNRemote.Contract.Owner(&_RMNRemote.CallOpts)
}

func (_RMNRemote *RMNRemoteCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _RMNRemote.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_RMNRemote *RMNRemoteSession) TypeAndVersion() (string, error) {
	return _RMNRemote.Contract.TypeAndVersion(&_RMNRemote.CallOpts)
}

func (_RMNRemote *RMNRemoteCallerSession) TypeAndVersion() (string, error) {
	return _RMNRemote.Contract.TypeAndVersion(&_RMNRemote.CallOpts)
}

func (_RMNRemote *RMNRemoteCaller) Verify(opts *bind.CallOpts, destLaneUpdates []InternalMerkleRoot, signatures []IRMNV2Signature) error {
	var out []interface{}
	err := _RMNRemote.contract.Call(opts, &out, "verify", destLaneUpdates, signatures)

	if err != nil {
		return err
	}

	return err

}

func (_RMNRemote *RMNRemoteSession) Verify(destLaneUpdates []InternalMerkleRoot, signatures []IRMNV2Signature) error {
	return _RMNRemote.Contract.Verify(&_RMNRemote.CallOpts, destLaneUpdates, signatures)
}

func (_RMNRemote *RMNRemoteCallerSession) Verify(destLaneUpdates []InternalMerkleRoot, signatures []IRMNV2Signature) error {
	return _RMNRemote.Contract.Verify(&_RMNRemote.CallOpts, destLaneUpdates, signatures)
}

func (_RMNRemote *RMNRemoteTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RMNRemote.contract.Transact(opts, "acceptOwnership")
}

func (_RMNRemote *RMNRemoteSession) AcceptOwnership() (*types.Transaction, error) {
	return _RMNRemote.Contract.AcceptOwnership(&_RMNRemote.TransactOpts)
}

func (_RMNRemote *RMNRemoteTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _RMNRemote.Contract.AcceptOwnership(&_RMNRemote.TransactOpts)
}

func (_RMNRemote *RMNRemoteTransactor) SetConfig(opts *bind.TransactOpts, newConfig RMNRemoteConfig) (*types.Transaction, error) {
	return _RMNRemote.contract.Transact(opts, "setConfig", newConfig)
}

func (_RMNRemote *RMNRemoteSession) SetConfig(newConfig RMNRemoteConfig) (*types.Transaction, error) {
	return _RMNRemote.Contract.SetConfig(&_RMNRemote.TransactOpts, newConfig)
}

func (_RMNRemote *RMNRemoteTransactorSession) SetConfig(newConfig RMNRemoteConfig) (*types.Transaction, error) {
	return _RMNRemote.Contract.SetConfig(&_RMNRemote.TransactOpts, newConfig)
}

func (_RMNRemote *RMNRemoteTransactor) Test(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RMNRemote.contract.Transact(opts, "test")
}

func (_RMNRemote *RMNRemoteSession) Test() (*types.Transaction, error) {
	return _RMNRemote.Contract.Test(&_RMNRemote.TransactOpts)
}

func (_RMNRemote *RMNRemoteTransactorSession) Test() (*types.Transaction, error) {
	return _RMNRemote.Contract.Test(&_RMNRemote.TransactOpts)
}

func (_RMNRemote *RMNRemoteTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _RMNRemote.contract.Transact(opts, "transferOwnership", to)
}

func (_RMNRemote *RMNRemoteSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _RMNRemote.Contract.TransferOwnership(&_RMNRemote.TransactOpts, to)
}

func (_RMNRemote *RMNRemoteTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _RMNRemote.Contract.TransferOwnership(&_RMNRemote.TransactOpts, to)
}

type RMNRemoteConfigSetIterator struct {
	Event *RMNRemoteConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RMNRemoteConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMNRemoteConfigSet)
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
		it.Event = new(RMNRemoteConfigSet)
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

func (it *RMNRemoteConfigSetIterator) Error() error {
	return it.fail
}

func (it *RMNRemoteConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RMNRemoteConfigSet struct {
	VersionedConfig RMNRemoteVersionedConfig
	Raw             types.Log
}

func (_RMNRemote *RMNRemoteFilterer) FilterConfigSet(opts *bind.FilterOpts) (*RMNRemoteConfigSetIterator, error) {

	logs, sub, err := _RMNRemote.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &RMNRemoteConfigSetIterator{contract: _RMNRemote.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_RMNRemote *RMNRemoteFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *RMNRemoteConfigSet) (event.Subscription, error) {

	logs, sub, err := _RMNRemote.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RMNRemoteConfigSet)
				if err := _RMNRemote.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_RMNRemote *RMNRemoteFilterer) ParseConfigSet(log types.Log) (*RMNRemoteConfigSet, error) {
	event := new(RMNRemoteConfigSet)
	if err := _RMNRemote.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RMNRemoteOwnershipTransferRequestedIterator struct {
	Event *RMNRemoteOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RMNRemoteOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMNRemoteOwnershipTransferRequested)
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
		it.Event = new(RMNRemoteOwnershipTransferRequested)
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

func (it *RMNRemoteOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *RMNRemoteOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RMNRemoteOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_RMNRemote *RMNRemoteFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RMNRemoteOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _RMNRemote.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &RMNRemoteOwnershipTransferRequestedIterator{contract: _RMNRemote.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_RMNRemote *RMNRemoteFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *RMNRemoteOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _RMNRemote.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RMNRemoteOwnershipTransferRequested)
				if err := _RMNRemote.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_RMNRemote *RMNRemoteFilterer) ParseOwnershipTransferRequested(log types.Log) (*RMNRemoteOwnershipTransferRequested, error) {
	event := new(RMNRemoteOwnershipTransferRequested)
	if err := _RMNRemote.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RMNRemoteOwnershipTransferredIterator struct {
	Event *RMNRemoteOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RMNRemoteOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMNRemoteOwnershipTransferred)
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
		it.Event = new(RMNRemoteOwnershipTransferred)
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

func (it *RMNRemoteOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *RMNRemoteOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RMNRemoteOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_RMNRemote *RMNRemoteFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RMNRemoteOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _RMNRemote.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &RMNRemoteOwnershipTransferredIterator{contract: _RMNRemote.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_RMNRemote *RMNRemoteFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *RMNRemoteOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _RMNRemote.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RMNRemoteOwnershipTransferred)
				if err := _RMNRemote.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_RMNRemote *RMNRemoteFilterer) ParseOwnershipTransferred(log types.Log) (*RMNRemoteOwnershipTransferred, error) {
	event := new(RMNRemoteOwnershipTransferred)
	if err := _RMNRemote.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_RMNRemote *RMNRemote) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _RMNRemote.abi.Events["ConfigSet"].ID:
		return _RMNRemote.ParseConfigSet(log)
	case _RMNRemote.abi.Events["OwnershipTransferRequested"].ID:
		return _RMNRemote.ParseOwnershipTransferRequested(log)
	case _RMNRemote.abi.Events["OwnershipTransferred"].ID:
		return _RMNRemote.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (RMNRemoteConfigSet) Topic() common.Hash {
	return common.HexToHash("0x6cc65868ae41a007e6c3ed18ce591c123dd4e5864b421888c68ce92dae98cea4")
}

func (RMNRemoteOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (RMNRemoteOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_RMNRemote *RMNRemote) Address() common.Address {
	return _RMNRemote.address
}

type RMNRemoteInterface interface {
	GetVersionedConfig(opts *bind.CallOpts) (RMNRemoteVersionedConfig, error)

	IsCursed(opts *bind.CallOpts, subject [16]byte) (bool, error)

	IsCursed0(opts *bind.CallOpts) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	Verify(opts *bind.CallOpts, destLaneUpdates []InternalMerkleRoot, signatures []IRMNV2Signature) error

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, newConfig RMNRemoteConfig) (*types.Transaction, error)

	Test(opts *bind.TransactOpts) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterConfigSet(opts *bind.FilterOpts) (*RMNRemoteConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *RMNRemoteConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*RMNRemoteConfigSet, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RMNRemoteOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *RMNRemoteOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*RMNRemoteOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RMNRemoteOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *RMNRemoteOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*RMNRemoteOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
