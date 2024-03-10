// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_owner_test_consumer

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

var VRFV2OwnerTestConsumerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlyCoordinatorCanFulfill\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"SubscriptionCreatedFundedAndConsumerAdded\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"COORDINATOR\",\"outputs\":[{\"internalType\":\"contractVRFCoordinatorV2Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINKTOKEN\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_requestId\",\"type\":\"uint256\"}],\"name\":\"getRequestStatus\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"fulfilled\",\"type\":\"bool\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"requestTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fulfilmentTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requestBlockNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fulfilmentBlockNumber\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"_requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bytes32\",\"name\":\"_keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"_callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_numWords\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"_requestCount\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"_subTopUpAmount\",\"type\":\"uint256\"}],\"name\":\"requestRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_averageFulfillmentInMillions\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_fastestFulfillment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_lastRequestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_requests\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"fulfilled\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"requestTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fulfilmentTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requestBlockNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fulfilmentBlockNumber\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_responseCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_slowestFulfillment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"subId\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"topUpSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a0604052600060065560006007556103e76008553480156200002157600080fd5b50604051620016ee380380620016ee8339810160408190526200004491620001e2565b6001600160601b0319606083901b166080523380600081620000ad5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000e057620000e08162000119565b5050600280546001600160a01b039485166001600160a01b0319918216179091556003805493909416921691909117909155506200021a565b6001600160a01b038116331415620001745760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401620000a4565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b0381168114620001dd57600080fd5b919050565b60008060408385031215620001f657600080fd5b6200020183620001c5565b91506200021160208401620001c5565b90509250929050565b60805160601c6114ae620002406000396000818161036901526103d101526114ae6000f3fe608060405234801561001057600080fd5b50600436106101365760003560e01c80638da5cb5b116100b2578063d8a4676f11610081578063eb1d28bb11610066578063eb1d28bb146102e6578063f2fde38b1461032b578063f82d24381461033e57600080fd5b8063d8a4676f146102b8578063dc1670db146102dd57600080fd5b80638da5cb5b14610207578063a168fa8914610225578063b1e2174914610290578063d826f88f1461029957600080fd5b8063557d2e921161010957806374dba124116100ee57806374dba124146101e357806379ba5097146101ec57806386850e93146101f457600080fd5b8063557d2e92146101d1578063737144bc146101da57600080fd5b80631757f11c1461013b5780631fe543e3146101575780633b2bcbf11461016c57806355380dfb146101b1575b600080fd5b61014460075481565b6040519081526020015b60405180910390f35b61016a610165366004611124565b610351565b005b60025461018c9073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161014e565b60035461018c9073ffffffffffffffffffffffffffffffffffffffff1681565b61014460055481565b61014460065481565b61014460085481565b61016a610411565b61016a6102023660046110f2565b61050e565b60005473ffffffffffffffffffffffffffffffffffffffff1661018c565b6102666102333660046110f2565b600b602052600090815260409020805460028201546003830154600484015460059094015460ff90931693919290919085565b6040805195151586526020860194909452928401919091526060830152608082015260a00161014e565b61014460095481565b61016a6000600681905560078190556103e76008556005819055600455565b6102cb6102c63660046110f2565b6105ea565b60405161014e969594939291906112d5565b61014460045481565b6003546103129074010000000000000000000000000000000000000000900467ffffffffffffffff1681565b60405167ffffffffffffffff909116815260200161014e565b61016a61033936600461102d565b6106cf565b61016a61034c36600461108c565b6106e3565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610403576040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001660248201526044015b60405180910390fd5b61040d8282610c3e565b5050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610492576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016103fa565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610516610d65565b6003546002546040805174010000000000000000000000000000000000000000840467ffffffffffffffff16602082015273ffffffffffffffffffffffffffffffffffffffff93841693634000aea09316918591016040516020818303038152906040526040518463ffffffff1660e01b81526004016105989392919061123d565b602060405180830381600087803b1580156105b257600080fd5b505af11580156105c6573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061040d919061106a565b6000818152600b60209081526040808320815160c081018352815460ff161515815260018201805484518187028101870190955280855260609587958695869586958695919492938584019390929083018282801561066857602002820191906000526020600020905b815481526020019060010190808311610654575b505050505081526020016002820154815260200160038201548152602001600482015481526020016005820154815250509050806000015181602001518260400151836060015184608001518560a001519650965096509650965096505091939550919395565b6106d7610d65565b6106e081610de8565b50565b6106eb610d65565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a21a23e46040518163ffffffff1660e01b8152600401602060405180830381600087803b15801561075557600080fd5b505af1158015610769573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061078d9190611213565b600380547fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000067ffffffffffffffff938416810291909117918290556002546040517f7341c10c00000000000000000000000000000000000000000000000000000000815291909204909216600483015230602483015273ffffffffffffffffffffffffffffffffffffffff1690637341c10c90604401600060405180830381600087803b15801561085457600080fd5b505af1158015610868573d6000803e3d6000fd5b505050506108758161050e565b600354604080517401000000000000000000000000000000000000000090920467ffffffffffffffff16825230602083015281018290527f56c142509574e8340ca0190b029c74464b84037d2876278ea0ade3ffb1f0042c9060600160405180910390a160005b8261ffff168161ffff161015610ad9576002546003546040517f5d3b1d30000000000000000000000000000000000000000000000000000000008152600481018990527401000000000000000000000000000000000000000090910467ffffffffffffffff16602482015261ffff8916604482015263ffffffff80881660648301528616608482015260009173ffffffffffffffffffffffffffffffffffffffff1690635d3b1d309060a401602060405180830381600087803b1580156109a257600080fd5b505af11580156109b6573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906109da919061110b565b6009819055905060006109eb610ede565b6040805160c08101825260008082528251818152602080820185528084019182524284860152606084018390526080840186905260a08401839052878352600b815293909120825181547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169015151781559051805194955091939092610a79926001850192910190610fa2565b5060408201516002820155606082015160038201556080820151600482015560a0909101516005918201558054906000610ab28361140a565b90915550506000918252600a60205260409091205580610ad1816113e8565b9150506108dc565b506002546003546040517f9f87fad70000000000000000000000000000000000000000000000000000000081527401000000000000000000000000000000000000000090910467ffffffffffffffff16600482015230602482015273ffffffffffffffffffffffffffffffffffffffff90911690639f87fad790604401600060405180830381600087803b158015610b7057600080fd5b505af1158015610b84573d6000803e3d6000fd5b50506002546003546040517fd7ae1d300000000000000000000000000000000000000000000000000000000081527401000000000000000000000000000000000000000090910467ffffffffffffffff16600482015233602482015273ffffffffffffffffffffffffffffffffffffffff909116925063d7ae1d309150604401600060405180830381600087803b158015610c1e57600080fd5b505af1158015610c32573d6000803e3d6000fd5b50505050505050505050565b6000610c48610ede565b6000848152600a602052604081205491925090610c6590836113d1565b90506000610c7682620f4240611394565b9050600754821115610c885760078290555b6008548210610c9957600854610c9b565b815b600855600454610cab5780610cde565b600454610cb9906001611341565b81600454600654610cca9190611394565b610cd49190611341565b610cde9190611359565b6006556000858152600b6020908152604090912080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016600190811782558651610d30939290910191870190610fa2565b506000858152600b602052604081204260038201556005018490556004805491610d598361140a565b91905055505050505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610de6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016103fa565b565b73ffffffffffffffffffffffffffffffffffffffff8116331415610e68576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016103fa565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600046610eea81610f7b565b15610f7457606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b815260040160206040518083038186803b158015610f3657600080fd5b505afa158015610f4a573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610f6e919061110b565b91505090565b4391505090565b600061a4b1821480610f8f575062066eed82145b80610f9c575062066eee82145b92915050565b828054828255906000526020600020908101928215610fdd579160200282015b82811115610fdd578251825591602001919060010190610fc2565b50610fe9929150610fed565b5090565b5b80821115610fe95760008155600101610fee565b803561ffff8116811461101457600080fd5b919050565b803563ffffffff8116811461101457600080fd5b60006020828403121561103f57600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461106357600080fd5b9392505050565b60006020828403121561107c57600080fd5b8151801515811461106357600080fd5b60008060008060008060c087890312156110a557600080fd5b6110ae87611002565b9550602087013594506110c360408801611019565b93506110d160608801611019565b92506110df60808801611002565b915060a087013590509295509295509295565b60006020828403121561110457600080fd5b5035919050565b60006020828403121561111d57600080fd5b5051919050565b6000806040838503121561113757600080fd5b8235915060208084013567ffffffffffffffff8082111561115757600080fd5b818601915086601f83011261116b57600080fd5b81358181111561117d5761117d611472565b8060051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f830116810181811085821117156111c0576111c0611472565b604052828152858101935084860182860187018b10156111df57600080fd5b600095505b838610156112025780358552600195909501949386019386016111e4565b508096505050505050509250929050565b60006020828403121561122557600080fd5b815167ffffffffffffffff8116811461106357600080fd5b73ffffffffffffffffffffffffffffffffffffffff8416815260006020848184015260606040840152835180606085015260005b8181101561128d57858101830151858201608001528201611271565b8181111561129f576000608083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160800195945050505050565b600060c082018815158352602060c08185015281895180845260e086019150828b01935060005b81811015611318578451835293830193918301916001016112fc565b505060408501989098525050506060810193909352608083019190915260a09091015292915050565b6000821982111561135457611354611443565b500190565b60008261138f577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156113cc576113cc611443565b500290565b6000828210156113e3576113e3611443565b500390565b600061ffff8083168181141561140057611400611443565b6001019392505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82141561143c5761143c611443565b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var VRFV2OwnerTestConsumerABI = VRFV2OwnerTestConsumerMetaData.ABI

var VRFV2OwnerTestConsumerBin = VRFV2OwnerTestConsumerMetaData.Bin

func DeployVRFV2OwnerTestConsumer(auth *bind.TransactOpts, backend bind.ContractBackend, _vrfCoordinator common.Address, _link common.Address) (common.Address, *types.Transaction, *VRFV2OwnerTestConsumer, error) {
	parsed, err := VRFV2OwnerTestConsumerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFV2OwnerTestConsumerBin), backend, _vrfCoordinator, _link)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFV2OwnerTestConsumer{address: address, abi: *parsed, VRFV2OwnerTestConsumerCaller: VRFV2OwnerTestConsumerCaller{contract: contract}, VRFV2OwnerTestConsumerTransactor: VRFV2OwnerTestConsumerTransactor{contract: contract}, VRFV2OwnerTestConsumerFilterer: VRFV2OwnerTestConsumerFilterer{contract: contract}}, nil
}

type VRFV2OwnerTestConsumer struct {
	address common.Address
	abi     abi.ABI
	VRFV2OwnerTestConsumerCaller
	VRFV2OwnerTestConsumerTransactor
	VRFV2OwnerTestConsumerFilterer
}

type VRFV2OwnerTestConsumerCaller struct {
	contract *bind.BoundContract
}

type VRFV2OwnerTestConsumerTransactor struct {
	contract *bind.BoundContract
}

type VRFV2OwnerTestConsumerFilterer struct {
	contract *bind.BoundContract
}

type VRFV2OwnerTestConsumerSession struct {
	Contract     *VRFV2OwnerTestConsumer
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFV2OwnerTestConsumerCallerSession struct {
	Contract *VRFV2OwnerTestConsumerCaller
	CallOpts bind.CallOpts
}

type VRFV2OwnerTestConsumerTransactorSession struct {
	Contract     *VRFV2OwnerTestConsumerTransactor
	TransactOpts bind.TransactOpts
}

type VRFV2OwnerTestConsumerRaw struct {
	Contract *VRFV2OwnerTestConsumer
}

type VRFV2OwnerTestConsumerCallerRaw struct {
	Contract *VRFV2OwnerTestConsumerCaller
}

type VRFV2OwnerTestConsumerTransactorRaw struct {
	Contract *VRFV2OwnerTestConsumerTransactor
}

func NewVRFV2OwnerTestConsumer(address common.Address, backend bind.ContractBackend) (*VRFV2OwnerTestConsumer, error) {
	abi, err := abi.JSON(strings.NewReader(VRFV2OwnerTestConsumerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFV2OwnerTestConsumer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFV2OwnerTestConsumer{address: address, abi: abi, VRFV2OwnerTestConsumerCaller: VRFV2OwnerTestConsumerCaller{contract: contract}, VRFV2OwnerTestConsumerTransactor: VRFV2OwnerTestConsumerTransactor{contract: contract}, VRFV2OwnerTestConsumerFilterer: VRFV2OwnerTestConsumerFilterer{contract: contract}}, nil
}

func NewVRFV2OwnerTestConsumerCaller(address common.Address, caller bind.ContractCaller) (*VRFV2OwnerTestConsumerCaller, error) {
	contract, err := bindVRFV2OwnerTestConsumer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2OwnerTestConsumerCaller{contract: contract}, nil
}

func NewVRFV2OwnerTestConsumerTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFV2OwnerTestConsumerTransactor, error) {
	contract, err := bindVRFV2OwnerTestConsumer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2OwnerTestConsumerTransactor{contract: contract}, nil
}

func NewVRFV2OwnerTestConsumerFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFV2OwnerTestConsumerFilterer, error) {
	contract, err := bindVRFV2OwnerTestConsumer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFV2OwnerTestConsumerFilterer{contract: contract}, nil
}

func bindVRFV2OwnerTestConsumer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFV2OwnerTestConsumerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2OwnerTestConsumer.Contract.VRFV2OwnerTestConsumerCaller.contract.Call(opts, result, method, params...)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2OwnerTestConsumer.Contract.VRFV2OwnerTestConsumerTransactor.contract.Transfer(opts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2OwnerTestConsumer.Contract.VRFV2OwnerTestConsumerTransactor.contract.Transact(opts, method, params...)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2OwnerTestConsumer.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2OwnerTestConsumer.Contract.contract.Transfer(opts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2OwnerTestConsumer.Contract.contract.Transact(opts, method, params...)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCaller) COORDINATOR(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2OwnerTestConsumer.contract.Call(opts, &out, "COORDINATOR")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerSession) COORDINATOR() (common.Address, error) {
	return _VRFV2OwnerTestConsumer.Contract.COORDINATOR(&_VRFV2OwnerTestConsumer.CallOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCallerSession) COORDINATOR() (common.Address, error) {
	return _VRFV2OwnerTestConsumer.Contract.COORDINATOR(&_VRFV2OwnerTestConsumer.CallOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCaller) LINKTOKEN(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2OwnerTestConsumer.contract.Call(opts, &out, "LINKTOKEN")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerSession) LINKTOKEN() (common.Address, error) {
	return _VRFV2OwnerTestConsumer.Contract.LINKTOKEN(&_VRFV2OwnerTestConsumer.CallOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCallerSession) LINKTOKEN() (common.Address, error) {
	return _VRFV2OwnerTestConsumer.Contract.LINKTOKEN(&_VRFV2OwnerTestConsumer.CallOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCaller) GetRequestStatus(opts *bind.CallOpts, _requestId *big.Int) (GetRequestStatus,

	error) {
	var out []interface{}
	err := _VRFV2OwnerTestConsumer.contract.Call(opts, &out, "getRequestStatus", _requestId)

	outstruct := new(GetRequestStatus)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Fulfilled = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.RandomWords = *abi.ConvertType(out[1], new([]*big.Int)).(*[]*big.Int)
	outstruct.RequestTimestamp = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.FulfilmentTimestamp = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.RequestBlockNumber = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.FulfilmentBlockNumber = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerSession) GetRequestStatus(_requestId *big.Int) (GetRequestStatus,

	error) {
	return _VRFV2OwnerTestConsumer.Contract.GetRequestStatus(&_VRFV2OwnerTestConsumer.CallOpts, _requestId)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCallerSession) GetRequestStatus(_requestId *big.Int) (GetRequestStatus,

	error) {
	return _VRFV2OwnerTestConsumer.Contract.GetRequestStatus(&_VRFV2OwnerTestConsumer.CallOpts, _requestId)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2OwnerTestConsumer.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerSession) Owner() (common.Address, error) {
	return _VRFV2OwnerTestConsumer.Contract.Owner(&_VRFV2OwnerTestConsumer.CallOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCallerSession) Owner() (common.Address, error) {
	return _VRFV2OwnerTestConsumer.Contract.Owner(&_VRFV2OwnerTestConsumer.CallOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCaller) SAverageFulfillmentInMillions(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2OwnerTestConsumer.contract.Call(opts, &out, "s_averageFulfillmentInMillions")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerSession) SAverageFulfillmentInMillions() (*big.Int, error) {
	return _VRFV2OwnerTestConsumer.Contract.SAverageFulfillmentInMillions(&_VRFV2OwnerTestConsumer.CallOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCallerSession) SAverageFulfillmentInMillions() (*big.Int, error) {
	return _VRFV2OwnerTestConsumer.Contract.SAverageFulfillmentInMillions(&_VRFV2OwnerTestConsumer.CallOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCaller) SFastestFulfillment(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2OwnerTestConsumer.contract.Call(opts, &out, "s_fastestFulfillment")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerSession) SFastestFulfillment() (*big.Int, error) {
	return _VRFV2OwnerTestConsumer.Contract.SFastestFulfillment(&_VRFV2OwnerTestConsumer.CallOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCallerSession) SFastestFulfillment() (*big.Int, error) {
	return _VRFV2OwnerTestConsumer.Contract.SFastestFulfillment(&_VRFV2OwnerTestConsumer.CallOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCaller) SLastRequestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2OwnerTestConsumer.contract.Call(opts, &out, "s_lastRequestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerSession) SLastRequestId() (*big.Int, error) {
	return _VRFV2OwnerTestConsumer.Contract.SLastRequestId(&_VRFV2OwnerTestConsumer.CallOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCallerSession) SLastRequestId() (*big.Int, error) {
	return _VRFV2OwnerTestConsumer.Contract.SLastRequestId(&_VRFV2OwnerTestConsumer.CallOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCaller) SRequestCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2OwnerTestConsumer.contract.Call(opts, &out, "s_requestCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerSession) SRequestCount() (*big.Int, error) {
	return _VRFV2OwnerTestConsumer.Contract.SRequestCount(&_VRFV2OwnerTestConsumer.CallOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCallerSession) SRequestCount() (*big.Int, error) {
	return _VRFV2OwnerTestConsumer.Contract.SRequestCount(&_VRFV2OwnerTestConsumer.CallOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCaller) SRequests(opts *bind.CallOpts, arg0 *big.Int) (SRequests,

	error) {
	var out []interface{}
	err := _VRFV2OwnerTestConsumer.contract.Call(opts, &out, "s_requests", arg0)

	outstruct := new(SRequests)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Fulfilled = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.RequestTimestamp = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.FulfilmentTimestamp = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.RequestBlockNumber = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.FulfilmentBlockNumber = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerSession) SRequests(arg0 *big.Int) (SRequests,

	error) {
	return _VRFV2OwnerTestConsumer.Contract.SRequests(&_VRFV2OwnerTestConsumer.CallOpts, arg0)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCallerSession) SRequests(arg0 *big.Int) (SRequests,

	error) {
	return _VRFV2OwnerTestConsumer.Contract.SRequests(&_VRFV2OwnerTestConsumer.CallOpts, arg0)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCaller) SResponseCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2OwnerTestConsumer.contract.Call(opts, &out, "s_responseCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerSession) SResponseCount() (*big.Int, error) {
	return _VRFV2OwnerTestConsumer.Contract.SResponseCount(&_VRFV2OwnerTestConsumer.CallOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCallerSession) SResponseCount() (*big.Int, error) {
	return _VRFV2OwnerTestConsumer.Contract.SResponseCount(&_VRFV2OwnerTestConsumer.CallOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCaller) SSlowestFulfillment(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2OwnerTestConsumer.contract.Call(opts, &out, "s_slowestFulfillment")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerSession) SSlowestFulfillment() (*big.Int, error) {
	return _VRFV2OwnerTestConsumer.Contract.SSlowestFulfillment(&_VRFV2OwnerTestConsumer.CallOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCallerSession) SSlowestFulfillment() (*big.Int, error) {
	return _VRFV2OwnerTestConsumer.Contract.SSlowestFulfillment(&_VRFV2OwnerTestConsumer.CallOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCaller) SubId(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _VRFV2OwnerTestConsumer.contract.Call(opts, &out, "subId")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerSession) SubId() (uint64, error) {
	return _VRFV2OwnerTestConsumer.Contract.SubId(&_VRFV2OwnerTestConsumer.CallOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCallerSession) SubId() (uint64, error) {
	return _VRFV2OwnerTestConsumer.Contract.SubId(&_VRFV2OwnerTestConsumer.CallOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2OwnerTestConsumer.contract.Transact(opts, "acceptOwnership")
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFV2OwnerTestConsumer.Contract.AcceptOwnership(&_VRFV2OwnerTestConsumer.TransactOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFV2OwnerTestConsumer.Contract.AcceptOwnership(&_VRFV2OwnerTestConsumer.TransactOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2OwnerTestConsumer.contract.Transact(opts, "rawFulfillRandomWords", requestId, randomWords)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2OwnerTestConsumer.Contract.RawFulfillRandomWords(&_VRFV2OwnerTestConsumer.TransactOpts, requestId, randomWords)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerTransactorSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2OwnerTestConsumer.Contract.RawFulfillRandomWords(&_VRFV2OwnerTestConsumer.TransactOpts, requestId, randomWords)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerTransactor) RequestRandomWords(opts *bind.TransactOpts, _requestConfirmations uint16, _keyHash [32]byte, _callbackGasLimit uint32, _numWords uint32, _requestCount uint16, _subTopUpAmount *big.Int) (*types.Transaction, error) {
	return _VRFV2OwnerTestConsumer.contract.Transact(opts, "requestRandomWords", _requestConfirmations, _keyHash, _callbackGasLimit, _numWords, _requestCount, _subTopUpAmount)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerSession) RequestRandomWords(_requestConfirmations uint16, _keyHash [32]byte, _callbackGasLimit uint32, _numWords uint32, _requestCount uint16, _subTopUpAmount *big.Int) (*types.Transaction, error) {
	return _VRFV2OwnerTestConsumer.Contract.RequestRandomWords(&_VRFV2OwnerTestConsumer.TransactOpts, _requestConfirmations, _keyHash, _callbackGasLimit, _numWords, _requestCount, _subTopUpAmount)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerTransactorSession) RequestRandomWords(_requestConfirmations uint16, _keyHash [32]byte, _callbackGasLimit uint32, _numWords uint32, _requestCount uint16, _subTopUpAmount *big.Int) (*types.Transaction, error) {
	return _VRFV2OwnerTestConsumer.Contract.RequestRandomWords(&_VRFV2OwnerTestConsumer.TransactOpts, _requestConfirmations, _keyHash, _callbackGasLimit, _numWords, _requestCount, _subTopUpAmount)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerTransactor) Reset(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2OwnerTestConsumer.contract.Transact(opts, "reset")
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerSession) Reset() (*types.Transaction, error) {
	return _VRFV2OwnerTestConsumer.Contract.Reset(&_VRFV2OwnerTestConsumer.TransactOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerTransactorSession) Reset() (*types.Transaction, error) {
	return _VRFV2OwnerTestConsumer.Contract.Reset(&_VRFV2OwnerTestConsumer.TransactOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerTransactor) TopUpSubscription(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VRFV2OwnerTestConsumer.contract.Transact(opts, "topUpSubscription", amount)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerSession) TopUpSubscription(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2OwnerTestConsumer.Contract.TopUpSubscription(&_VRFV2OwnerTestConsumer.TransactOpts, amount)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerTransactorSession) TopUpSubscription(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2OwnerTestConsumer.Contract.TopUpSubscription(&_VRFV2OwnerTestConsumer.TransactOpts, amount)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFV2OwnerTestConsumer.contract.Transact(opts, "transferOwnership", to)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFV2OwnerTestConsumer.Contract.TransferOwnership(&_VRFV2OwnerTestConsumer.TransactOpts, to)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFV2OwnerTestConsumer.Contract.TransferOwnership(&_VRFV2OwnerTestConsumer.TransactOpts, to)
}

type VRFV2OwnerTestConsumerOwnershipTransferRequestedIterator struct {
	Event *VRFV2OwnerTestConsumerOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2OwnerTestConsumerOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2OwnerTestConsumerOwnershipTransferRequested)
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
		it.Event = new(VRFV2OwnerTestConsumerOwnershipTransferRequested)
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

func (it *VRFV2OwnerTestConsumerOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFV2OwnerTestConsumerOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2OwnerTestConsumerOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2OwnerTestConsumerOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2OwnerTestConsumer.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2OwnerTestConsumerOwnershipTransferRequestedIterator{contract: _VRFV2OwnerTestConsumer.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFV2OwnerTestConsumerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2OwnerTestConsumer.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2OwnerTestConsumerOwnershipTransferRequested)
				if err := _VRFV2OwnerTestConsumer.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerFilterer) ParseOwnershipTransferRequested(log types.Log) (*VRFV2OwnerTestConsumerOwnershipTransferRequested, error) {
	event := new(VRFV2OwnerTestConsumerOwnershipTransferRequested)
	if err := _VRFV2OwnerTestConsumer.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFV2OwnerTestConsumerOwnershipTransferredIterator struct {
	Event *VRFV2OwnerTestConsumerOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2OwnerTestConsumerOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2OwnerTestConsumerOwnershipTransferred)
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
		it.Event = new(VRFV2OwnerTestConsumerOwnershipTransferred)
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

func (it *VRFV2OwnerTestConsumerOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFV2OwnerTestConsumerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2OwnerTestConsumerOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2OwnerTestConsumerOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2OwnerTestConsumer.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2OwnerTestConsumerOwnershipTransferredIterator{contract: _VRFV2OwnerTestConsumer.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFV2OwnerTestConsumerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2OwnerTestConsumer.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2OwnerTestConsumerOwnershipTransferred)
				if err := _VRFV2OwnerTestConsumer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerFilterer) ParseOwnershipTransferred(log types.Log) (*VRFV2OwnerTestConsumerOwnershipTransferred, error) {
	event := new(VRFV2OwnerTestConsumerOwnershipTransferred)
	if err := _VRFV2OwnerTestConsumer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFV2OwnerTestConsumerSubscriptionCreatedFundedAndConsumerAddedIterator struct {
	Event *VRFV2OwnerTestConsumerSubscriptionCreatedFundedAndConsumerAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2OwnerTestConsumerSubscriptionCreatedFundedAndConsumerAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2OwnerTestConsumerSubscriptionCreatedFundedAndConsumerAdded)
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
		it.Event = new(VRFV2OwnerTestConsumerSubscriptionCreatedFundedAndConsumerAdded)
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

func (it *VRFV2OwnerTestConsumerSubscriptionCreatedFundedAndConsumerAddedIterator) Error() error {
	return it.fail
}

func (it *VRFV2OwnerTestConsumerSubscriptionCreatedFundedAndConsumerAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2OwnerTestConsumerSubscriptionCreatedFundedAndConsumerAdded struct {
	SubId    uint64
	Consumer common.Address
	Amount   *big.Int
	Raw      types.Log
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerFilterer) FilterSubscriptionCreatedFundedAndConsumerAdded(opts *bind.FilterOpts) (*VRFV2OwnerTestConsumerSubscriptionCreatedFundedAndConsumerAddedIterator, error) {

	logs, sub, err := _VRFV2OwnerTestConsumer.contract.FilterLogs(opts, "SubscriptionCreatedFundedAndConsumerAdded")
	if err != nil {
		return nil, err
	}
	return &VRFV2OwnerTestConsumerSubscriptionCreatedFundedAndConsumerAddedIterator{contract: _VRFV2OwnerTestConsumer.contract, event: "SubscriptionCreatedFundedAndConsumerAdded", logs: logs, sub: sub}, nil
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerFilterer) WatchSubscriptionCreatedFundedAndConsumerAdded(opts *bind.WatchOpts, sink chan<- *VRFV2OwnerTestConsumerSubscriptionCreatedFundedAndConsumerAdded) (event.Subscription, error) {

	logs, sub, err := _VRFV2OwnerTestConsumer.contract.WatchLogs(opts, "SubscriptionCreatedFundedAndConsumerAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2OwnerTestConsumerSubscriptionCreatedFundedAndConsumerAdded)
				if err := _VRFV2OwnerTestConsumer.contract.UnpackLog(event, "SubscriptionCreatedFundedAndConsumerAdded", log); err != nil {
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

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerFilterer) ParseSubscriptionCreatedFundedAndConsumerAdded(log types.Log) (*VRFV2OwnerTestConsumerSubscriptionCreatedFundedAndConsumerAdded, error) {
	event := new(VRFV2OwnerTestConsumerSubscriptionCreatedFundedAndConsumerAdded)
	if err := _VRFV2OwnerTestConsumer.contract.UnpackLog(event, "SubscriptionCreatedFundedAndConsumerAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetRequestStatus struct {
	Fulfilled             bool
	RandomWords           []*big.Int
	RequestTimestamp      *big.Int
	FulfilmentTimestamp   *big.Int
	RequestBlockNumber    *big.Int
	FulfilmentBlockNumber *big.Int
}
type SRequests struct {
	Fulfilled             bool
	RequestTimestamp      *big.Int
	FulfilmentTimestamp   *big.Int
	RequestBlockNumber    *big.Int
	FulfilmentBlockNumber *big.Int
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumer) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFV2OwnerTestConsumer.abi.Events["OwnershipTransferRequested"].ID:
		return _VRFV2OwnerTestConsumer.ParseOwnershipTransferRequested(log)
	case _VRFV2OwnerTestConsumer.abi.Events["OwnershipTransferred"].ID:
		return _VRFV2OwnerTestConsumer.ParseOwnershipTransferred(log)
	case _VRFV2OwnerTestConsumer.abi.Events["SubscriptionCreatedFundedAndConsumerAdded"].ID:
		return _VRFV2OwnerTestConsumer.ParseSubscriptionCreatedFundedAndConsumerAdded(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFV2OwnerTestConsumerOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VRFV2OwnerTestConsumerOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (VRFV2OwnerTestConsumerSubscriptionCreatedFundedAndConsumerAdded) Topic() common.Hash {
	return common.HexToHash("0x56c142509574e8340ca0190b029c74464b84037d2876278ea0ade3ffb1f0042c")
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumer) Address() common.Address {
	return _VRFV2OwnerTestConsumer.address
}

type VRFV2OwnerTestConsumerInterface interface {
	COORDINATOR(opts *bind.CallOpts) (common.Address, error)

	LINKTOKEN(opts *bind.CallOpts) (common.Address, error)

	GetRequestStatus(opts *bind.CallOpts, _requestId *big.Int) (GetRequestStatus,

		error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SAverageFulfillmentInMillions(opts *bind.CallOpts) (*big.Int, error)

	SFastestFulfillment(opts *bind.CallOpts) (*big.Int, error)

	SLastRequestId(opts *bind.CallOpts) (*big.Int, error)

	SRequestCount(opts *bind.CallOpts) (*big.Int, error)

	SRequests(opts *bind.CallOpts, arg0 *big.Int) (SRequests,

		error)

	SResponseCount(opts *bind.CallOpts) (*big.Int, error)

	SSlowestFulfillment(opts *bind.CallOpts) (*big.Int, error)

	SubId(opts *bind.CallOpts) (uint64, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	RequestRandomWords(opts *bind.TransactOpts, _requestConfirmations uint16, _keyHash [32]byte, _callbackGasLimit uint32, _numWords uint32, _requestCount uint16, _subTopUpAmount *big.Int) (*types.Transaction, error)

	Reset(opts *bind.TransactOpts) (*types.Transaction, error)

	TopUpSubscription(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2OwnerTestConsumerOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFV2OwnerTestConsumerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VRFV2OwnerTestConsumerOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2OwnerTestConsumerOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFV2OwnerTestConsumerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VRFV2OwnerTestConsumerOwnershipTransferred, error)

	FilterSubscriptionCreatedFundedAndConsumerAdded(opts *bind.FilterOpts) (*VRFV2OwnerTestConsumerSubscriptionCreatedFundedAndConsumerAddedIterator, error)

	WatchSubscriptionCreatedFundedAndConsumerAdded(opts *bind.WatchOpts, sink chan<- *VRFV2OwnerTestConsumerSubscriptionCreatedFundedAndConsumerAdded) (event.Subscription, error)

	ParseSubscriptionCreatedFundedAndConsumerAdded(log types.Log) (*VRFV2OwnerTestConsumerSubscriptionCreatedFundedAndConsumerAdded, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
