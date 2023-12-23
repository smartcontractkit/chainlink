// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package chain_reader_example

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

type InnerTestStruct struct {
	I int64
	S string
}

type MidLevelTestStruct struct {
	FixedBytes [2]byte
	Inner      InnerTestStruct
}

type TestStruct struct {
	Field          int32
	DifferentField string
	OracleId       uint8
	OracleIds      [32]uint8
	Account        common.Address
	Accounts       []common.Address
	BigField       *big.Int
	NestedStruct   MidLevelTestStruct
}

var LatestValueHolderMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int32\",\"name\":\"field\",\"type\":\"int32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"differentField\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"oracleId\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint8[32]\",\"name\":\"oracleIds\",\"type\":\"uint8[32]\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"Account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"Accounts\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"int192\",\"name\":\"bigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"I\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"indexed\":false,\"internalType\":\"structMidLevelTestStruct\",\"name\":\"nestedStruct\",\"type\":\"tuple\"}],\"name\":\"Triggered\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"int32\",\"name\":\"field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"differentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"oracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"oracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"int192\",\"name\":\"bigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"I\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelTestStruct\",\"name\":\"nestedStruct\",\"type\":\"tuple\"}],\"name\":\"AddTestStruct\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GetDifferentPrimitiveValue\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"i\",\"type\":\"uint256\"}],\"name\":\"GetElementAtIndex\",\"outputs\":[{\"components\":[{\"internalType\":\"int32\",\"name\":\"Field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"DifferentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"OracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"OracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"address\",\"name\":\"Account\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"Accounts\",\"type\":\"address[]\"},{\"internalType\":\"int192\",\"name\":\"BigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"I\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelTestStruct\",\"name\":\"NestedStruct\",\"type\":\"tuple\"}],\"internalType\":\"structTestStruct\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GetPrimitiveValue\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GetSliceValue\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int32\",\"name\":\"field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"differentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"oracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"oracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"int192\",\"name\":\"bigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"I\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelTestStruct\",\"name\":\"nestedStruct\",\"type\":\"tuple\"}],\"name\":\"ReturnSeen\",\"outputs\":[{\"components\":[{\"internalType\":\"int32\",\"name\":\"Field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"DifferentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"OracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"OracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"address\",\"name\":\"Account\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"Accounts\",\"type\":\"address[]\"},{\"internalType\":\"int192\",\"name\":\"BigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"I\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelTestStruct\",\"name\":\"NestedStruct\",\"type\":\"tuple\"}],\"internalType\":\"structTestStruct\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int32\",\"name\":\"field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"differentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"oracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"oracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"int192\",\"name\":\"bigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"I\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelTestStruct\",\"name\":\"nestedStruct\",\"type\":\"tuple\"}],\"name\":\"TriggerEvent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50600180548082018255600082905260048082047fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6908101805460086003958616810261010090810a8088026001600160401b0391820219909416939093179093558654808801909755848704909301805496909516909202900a9182029102199092169190911790556115ef806100a96000396000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c80639ca04f671161005b5780639ca04f67146100cc578063b9dad6b0146100ec578063bdb37c90146100ff578063da8e7a821461011457600080fd5b8063030d3ca2146100825780636c7cf955146100a45780638b659d6e146100b9575b600080fd5b6107c65b60405167ffffffffffffffff90911681526020015b60405180910390f35b6100b76100b2366004610be2565b61011b565b005b6100b76100c7366004610be2565b61041e565b6100df6100da366004610cd4565b610473565b60405161009b9190610e33565b6100df6100fa366004610be2565b61074e565b610107610857565b60405161009b9190610f29565b6003610086565b60006040518061010001604052808c60030b81526020018b8b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050509082525060ff8a166020808301919091526040805161040081810183529190930192918b9183908390808284376000920191909152505050815273ffffffffffffffffffffffffffffffffffffffff8816602080830191909152604080518883028181018401835289825291909301929189918991829190850190849080828437600092019190915250505090825250601785900b602082015260400161020d84611060565b905281546001808201845560009384526020938490208351600a9093020180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff9093169290921782559282015191929091908201906102739082611207565b5060408201516002820180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff90921691909117905560608201516102c190600383019060206108e3565b5060808201516004820180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff90921691909117905560a08201518051610328916005840191602090910190610976565b5060c08201516006820180547fffffffffffffffff0000000000000000000000000000000000000000000000001677ffffffffffffffffffffffffffffffffffffffffffffffff90921691909117905560e082015180516007830180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001660f09290921c91909117815560208083015180516008860180547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001667ffffffffffffffff90921691909117815591810151909190600986019061040b9082611207565b5050505050505050505050505050505050565b7f7188419dcd8b51877b71766f075f3626586c0ff190e7d056aa65ce9acb649a3d8a8a8a8a8a8a8a8a8a8a60405161045f9a999897969594939291906114af565b60405180910390a150505050505050505050565b61047b6109f0565b6000610488600184611579565b81548110610498576104986115b3565b6000918252602091829020604080516101008101909152600a90920201805460030b825260018101805492939192918401916104d39061116b565b80601f01602080910402602001604051908101604052809291908181526020018280546104ff9061116b565b801561054c5780601f106105215761010080835404028352916020019161054c565b820191906000526020600020905b81548152906001019060200180831161052f57829003601f168201915b5050509183525050600282015460ff166020808301919091526040805161040081018083529190930192916003850191826000855b825461010083900a900460ff1681526020600192830181810494850194909303909202910180841161058157505050928452505050600482015473ffffffffffffffffffffffffffffffffffffffff16602080830191909152600583018054604080518285028101850182528281529401939283018282801561063a57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff16815260019091019060200180831161060f575b5050509183525050600682015460170b6020808301919091526040805180820182526007808601805460f01b7fffff0000000000000000000000000000000000000000000000000000000000001683528351808501855260088801805490930b815260098801805495909701969395919486830194919392840191906106bf9061116b565b80601f01602080910402602001604051908101604052809291908181526020018280546106eb9061116b565b80156107385780601f1061070d57610100808354040283529160200191610738565b820191906000526020600020905b81548152906001019060200180831161071b57829003601f168201915b5050509190925250505090525090525092915050565b6107566109f0565b6040518061010001604052808c60030b81526020018b8b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050509082525060ff8a166020808301919091526040805161040081810183529190930192918b9183908390808284376000920191909152505050815273ffffffffffffffffffffffffffffffffffffffff8816602080830191909152604080518883028181018401835289825291909301929189918991829190850190849080828437600092019190915250505090825250601785900b602082015260400161084684611060565b90529b9a5050505050505050505050565b606060018054806020026020016040519081016040528092919081815260200182805480156108d957602002820191906000526020600020906000905b82829054906101000a900467ffffffffffffffff1667ffffffffffffffff16815260200190600801906020826007010492830192600103820291508084116108945790505b5050505050905090565b6001830191839082156109665791602002820160005b8382111561093757835183826101000a81548160ff021916908360ff16021790555092602001926001016020816000010492830192600103026108f9565b80156109645782816101000a81549060ff0219169055600101602081600001049283019260010302610937565b505b50610972929150610a3f565b5090565b828054828255906000526020600020908101928215610966579160200282015b8281111561096657825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190610996565b6040805161010081018252600080825260606020830181905292820152908101610a18610a54565b8152600060208201819052606060408301819052820152608001610a3a610a73565b905290565b5b808211156109725760008155600101610a40565b6040518061040001604052806020906020820280368337509192915050565b604051806040016040528060007dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168152602001610a3a6040518060400160405280600060070b8152602001606081525090565b8035600381900b8114610ad857600080fd5b919050565b60008083601f840112610aef57600080fd5b50813567ffffffffffffffff811115610b0757600080fd5b602083019150836020828501011115610b1f57600080fd5b9250929050565b803560ff81168114610ad857600080fd5b806104008101831015610b4957600080fd5b92915050565b803573ffffffffffffffffffffffffffffffffffffffff81168114610ad857600080fd5b60008083601f840112610b8557600080fd5b50813567ffffffffffffffff811115610b9d57600080fd5b6020830191508360208260051b8501011115610b1f57600080fd5b8035601781900b8114610ad857600080fd5b600060408284031215610bdc57600080fd5b50919050565b6000806000806000806000806000806104e08b8d031215610c0257600080fd5b610c0b8b610ac6565b995060208b013567ffffffffffffffff80821115610c2857600080fd5b610c348e838f01610add565b909b509950899150610c4860408e01610b26565b9850610c578e60608f01610b37565b9750610c666104608e01610b4f565b96506104808d0135915080821115610c7d57600080fd5b610c898e838f01610b73565b9096509450849150610c9e6104a08e01610bb8565b93506104c08d0135915080821115610cb557600080fd5b50610cc28d828e01610bca565b9150509295989b9194979a5092959850565b600060208284031215610ce657600080fd5b5035919050565b6000815180845260005b81811015610d1357602081850181015186830182015201610cf7565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b8060005b6020808210610d645750610d7b565b825160ff1685529384019390910190600101610d55565b50505050565b600081518084526020808501945080840160005b83811015610dc757815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101610d95565b509495945050505050565b7fffff00000000000000000000000000000000000000000000000000000000000081511682526000602082015160406020850152805160070b60408501526020810151905060406060850152610e2b6080850182610ced565b949350505050565b60208152610e4760208201835160030b9052565b600060208301516104e0806040850152610e65610500850183610ced565b91506040850151610e7b606086018260ff169052565b506060850151610e8e6080860182610d51565b50608085015173ffffffffffffffffffffffffffffffffffffffff1661048085015260a08501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe085840381016104a0870152610eeb8483610d81565b935060c08701519150610f046104c087018360170b9052565b60e0870151915080868503018387015250610f1f8382610dd2565b9695505050505050565b6020808252825182820181905260009190848201906040850190845b81811015610f6b57835167ffffffffffffffff1683529284019291840191600101610f45565b50909695505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040805190810167ffffffffffffffff81118282101715610fc957610fc9610f77565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561101657611016610f77565b604052919050565b80357fffff00000000000000000000000000000000000000000000000000000000000081168114610ad857600080fd5b8035600781900b8114610ad857600080fd5b60006040823603121561107257600080fd5b61107a610fa6565b6110838361101e565b815260208084013567ffffffffffffffff808211156110a157600080fd5b8186019150604082360312156110b657600080fd5b6110be610fa6565b6110c78361104e565b815283830135828111156110da57600080fd5b929092019136601f8401126110ee57600080fd5b82358281111561110057611100610f77565b611130857fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601610fcf565b9250808352368582860101111561114657600080fd5b8085850186850137600090830185015280840191909152918301919091525092915050565b600181811c9082168061117f57607f821691505b602082108103610bdc577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b601f82111561120257600081815260208120601f850160051c810160208610156111df5750805b601f850160051c820191505b818110156111fe578281556001016111eb565b5050505b505050565b815167ffffffffffffffff81111561122157611221610f77565b6112358161122f845461116b565b846111b8565b602080601f83116001811461128857600084156112525750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b1785556111fe565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b828110156112d5578886015182559484019460019091019084016112b6565b508582101561131157878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b8183526000602080850194508260005b85811015610dc75773ffffffffffffffffffffffffffffffffffffffff6113a083610b4f565b168752958201959082019060010161137a565b7fffff0000000000000000000000000000000000000000000000000000000000006113dd8261101e565b168252600060208201357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc183360301811261141757600080fd5b6040602085015282016114298161104e565b60070b604085015260208101357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe182360301811261146657600080fd5b0160208101903567ffffffffffffffff81111561148257600080fd5b80360382131561149157600080fd5b604060608601526114a6608086018284611321565b95945050505050565b60006104e08c60030b8352602081818501526114ce8285018d8f611321565b915060ff808c166040860152606085018b60005b8481101561150757836114f483610b26565b16835291840191908401906001016114e2565b505050505061152f61046084018973ffffffffffffffffffffffffffffffffffffffff169052565b82810361048084015261154381878961136a565b90506115556104a084018660170b9052565b8281036104c084015261156881856113b3565b9d9c50505050505050505050505050565b81810381811115610b49577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fdfea164736f6c6343000813000a",
}

var LatestValueHolderABI = LatestValueHolderMetaData.ABI

var LatestValueHolderBin = LatestValueHolderMetaData.Bin

func DeployLatestValueHolder(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *LatestValueHolder, error) {
	parsed, err := LatestValueHolderMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(LatestValueHolderBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LatestValueHolder{address: address, abi: *parsed, LatestValueHolderCaller: LatestValueHolderCaller{contract: contract}, LatestValueHolderTransactor: LatestValueHolderTransactor{contract: contract}, LatestValueHolderFilterer: LatestValueHolderFilterer{contract: contract}}, nil
}

type LatestValueHolder struct {
	address common.Address
	abi     abi.ABI
	LatestValueHolderCaller
	LatestValueHolderTransactor
	LatestValueHolderFilterer
}

type LatestValueHolderCaller struct {
	contract *bind.BoundContract
}

type LatestValueHolderTransactor struct {
	contract *bind.BoundContract
}

type LatestValueHolderFilterer struct {
	contract *bind.BoundContract
}

type LatestValueHolderSession struct {
	Contract     *LatestValueHolder
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type LatestValueHolderCallerSession struct {
	Contract *LatestValueHolderCaller
	CallOpts bind.CallOpts
}

type LatestValueHolderTransactorSession struct {
	Contract     *LatestValueHolderTransactor
	TransactOpts bind.TransactOpts
}

type LatestValueHolderRaw struct {
	Contract *LatestValueHolder
}

type LatestValueHolderCallerRaw struct {
	Contract *LatestValueHolderCaller
}

type LatestValueHolderTransactorRaw struct {
	Contract *LatestValueHolderTransactor
}

func NewLatestValueHolder(address common.Address, backend bind.ContractBackend) (*LatestValueHolder, error) {
	abi, err := abi.JSON(strings.NewReader(LatestValueHolderABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindLatestValueHolder(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LatestValueHolder{address: address, abi: abi, LatestValueHolderCaller: LatestValueHolderCaller{contract: contract}, LatestValueHolderTransactor: LatestValueHolderTransactor{contract: contract}, LatestValueHolderFilterer: LatestValueHolderFilterer{contract: contract}}, nil
}

func NewLatestValueHolderCaller(address common.Address, caller bind.ContractCaller) (*LatestValueHolderCaller, error) {
	contract, err := bindLatestValueHolder(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LatestValueHolderCaller{contract: contract}, nil
}

func NewLatestValueHolderTransactor(address common.Address, transactor bind.ContractTransactor) (*LatestValueHolderTransactor, error) {
	contract, err := bindLatestValueHolder(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LatestValueHolderTransactor{contract: contract}, nil
}

func NewLatestValueHolderFilterer(address common.Address, filterer bind.ContractFilterer) (*LatestValueHolderFilterer, error) {
	contract, err := bindLatestValueHolder(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LatestValueHolderFilterer{contract: contract}, nil
}

func bindLatestValueHolder(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := LatestValueHolderMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_LatestValueHolder *LatestValueHolderRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LatestValueHolder.Contract.LatestValueHolderCaller.contract.Call(opts, result, method, params...)
}

func (_LatestValueHolder *LatestValueHolderRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LatestValueHolder.Contract.LatestValueHolderTransactor.contract.Transfer(opts)
}

func (_LatestValueHolder *LatestValueHolderRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LatestValueHolder.Contract.LatestValueHolderTransactor.contract.Transact(opts, method, params...)
}

func (_LatestValueHolder *LatestValueHolderCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LatestValueHolder.Contract.contract.Call(opts, result, method, params...)
}

func (_LatestValueHolder *LatestValueHolderTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LatestValueHolder.Contract.contract.Transfer(opts)
}

func (_LatestValueHolder *LatestValueHolderTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LatestValueHolder.Contract.contract.Transact(opts, method, params...)
}

func (_LatestValueHolder *LatestValueHolderCaller) GetDifferentPrimitiveValue(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _LatestValueHolder.contract.Call(opts, &out, "GetDifferentPrimitiveValue")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_LatestValueHolder *LatestValueHolderSession) GetDifferentPrimitiveValue() (uint64, error) {
	return _LatestValueHolder.Contract.GetDifferentPrimitiveValue(&_LatestValueHolder.CallOpts)
}

func (_LatestValueHolder *LatestValueHolderCallerSession) GetDifferentPrimitiveValue() (uint64, error) {
	return _LatestValueHolder.Contract.GetDifferentPrimitiveValue(&_LatestValueHolder.CallOpts)
}

func (_LatestValueHolder *LatestValueHolderCaller) GetElementAtIndex(opts *bind.CallOpts, i *big.Int) (TestStruct, error) {
	var out []interface{}
	err := _LatestValueHolder.contract.Call(opts, &out, "GetElementAtIndex", i)

	if err != nil {
		return *new(TestStruct), err
	}

	out0 := *abi.ConvertType(out[0], new(TestStruct)).(*TestStruct)

	return out0, err

}

func (_LatestValueHolder *LatestValueHolderSession) GetElementAtIndex(i *big.Int) (TestStruct, error) {
	return _LatestValueHolder.Contract.GetElementAtIndex(&_LatestValueHolder.CallOpts, i)
}

func (_LatestValueHolder *LatestValueHolderCallerSession) GetElementAtIndex(i *big.Int) (TestStruct, error) {
	return _LatestValueHolder.Contract.GetElementAtIndex(&_LatestValueHolder.CallOpts, i)
}

func (_LatestValueHolder *LatestValueHolderCaller) GetPrimitiveValue(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _LatestValueHolder.contract.Call(opts, &out, "GetPrimitiveValue")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_LatestValueHolder *LatestValueHolderSession) GetPrimitiveValue() (uint64, error) {
	return _LatestValueHolder.Contract.GetPrimitiveValue(&_LatestValueHolder.CallOpts)
}

func (_LatestValueHolder *LatestValueHolderCallerSession) GetPrimitiveValue() (uint64, error) {
	return _LatestValueHolder.Contract.GetPrimitiveValue(&_LatestValueHolder.CallOpts)
}

func (_LatestValueHolder *LatestValueHolderCaller) GetSliceValue(opts *bind.CallOpts) ([]uint64, error) {
	var out []interface{}
	err := _LatestValueHolder.contract.Call(opts, &out, "GetSliceValue")

	if err != nil {
		return *new([]uint64), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint64)).(*[]uint64)

	return out0, err

}

func (_LatestValueHolder *LatestValueHolderSession) GetSliceValue() ([]uint64, error) {
	return _LatestValueHolder.Contract.GetSliceValue(&_LatestValueHolder.CallOpts)
}

func (_LatestValueHolder *LatestValueHolderCallerSession) GetSliceValue() ([]uint64, error) {
	return _LatestValueHolder.Contract.GetSliceValue(&_LatestValueHolder.CallOpts)
}

func (_LatestValueHolder *LatestValueHolderCaller) ReturnSeen(opts *bind.CallOpts, field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedStruct MidLevelTestStruct) (TestStruct, error) {
	var out []interface{}
	err := _LatestValueHolder.contract.Call(opts, &out, "ReturnSeen", field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)

	if err != nil {
		return *new(TestStruct), err
	}

	out0 := *abi.ConvertType(out[0], new(TestStruct)).(*TestStruct)

	return out0, err

}

func (_LatestValueHolder *LatestValueHolderSession) ReturnSeen(field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedStruct MidLevelTestStruct) (TestStruct, error) {
	return _LatestValueHolder.Contract.ReturnSeen(&_LatestValueHolder.CallOpts, field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)
}

func (_LatestValueHolder *LatestValueHolderCallerSession) ReturnSeen(field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedStruct MidLevelTestStruct) (TestStruct, error) {
	return _LatestValueHolder.Contract.ReturnSeen(&_LatestValueHolder.CallOpts, field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)
}

func (_LatestValueHolder *LatestValueHolderTransactor) AddTestStruct(opts *bind.TransactOpts, field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedStruct MidLevelTestStruct) (*types.Transaction, error) {
	return _LatestValueHolder.contract.Transact(opts, "AddTestStruct", field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)
}

func (_LatestValueHolder *LatestValueHolderSession) AddTestStruct(field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedStruct MidLevelTestStruct) (*types.Transaction, error) {
	return _LatestValueHolder.Contract.AddTestStruct(&_LatestValueHolder.TransactOpts, field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)
}

func (_LatestValueHolder *LatestValueHolderTransactorSession) AddTestStruct(field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedStruct MidLevelTestStruct) (*types.Transaction, error) {
	return _LatestValueHolder.Contract.AddTestStruct(&_LatestValueHolder.TransactOpts, field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)
}

func (_LatestValueHolder *LatestValueHolderTransactor) TriggerEvent(opts *bind.TransactOpts, field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedStruct MidLevelTestStruct) (*types.Transaction, error) {
	return _LatestValueHolder.contract.Transact(opts, "TriggerEvent", field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)
}

func (_LatestValueHolder *LatestValueHolderSession) TriggerEvent(field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedStruct MidLevelTestStruct) (*types.Transaction, error) {
	return _LatestValueHolder.Contract.TriggerEvent(&_LatestValueHolder.TransactOpts, field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)
}

func (_LatestValueHolder *LatestValueHolderTransactorSession) TriggerEvent(field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedStruct MidLevelTestStruct) (*types.Transaction, error) {
	return _LatestValueHolder.Contract.TriggerEvent(&_LatestValueHolder.TransactOpts, field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)
}

type LatestValueHolderTriggeredIterator struct {
	Event *LatestValueHolderTriggered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LatestValueHolderTriggeredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LatestValueHolderTriggered)
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
		it.Event = new(LatestValueHolderTriggered)
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

func (it *LatestValueHolderTriggeredIterator) Error() error {
	return it.fail
}

func (it *LatestValueHolderTriggeredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LatestValueHolderTriggered struct {
	Field          int32
	DifferentField string
	OracleId       uint8
	OracleIds      [32]uint8
	Account        common.Address
	Accounts       []common.Address
	BigField       *big.Int
	NestedStruct   MidLevelTestStruct
	Raw            types.Log
}

func (_LatestValueHolder *LatestValueHolderFilterer) FilterTriggered(opts *bind.FilterOpts) (*LatestValueHolderTriggeredIterator, error) {

	logs, sub, err := _LatestValueHolder.contract.FilterLogs(opts, "Triggered")
	if err != nil {
		return nil, err
	}
	return &LatestValueHolderTriggeredIterator{contract: _LatestValueHolder.contract, event: "Triggered", logs: logs, sub: sub}, nil
}

func (_LatestValueHolder *LatestValueHolderFilterer) WatchTriggered(opts *bind.WatchOpts, sink chan<- *LatestValueHolderTriggered) (event.Subscription, error) {

	logs, sub, err := _LatestValueHolder.contract.WatchLogs(opts, "Triggered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LatestValueHolderTriggered)
				if err := _LatestValueHolder.contract.UnpackLog(event, "Triggered", log); err != nil {
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

func (_LatestValueHolder *LatestValueHolderFilterer) ParseTriggered(log types.Log) (*LatestValueHolderTriggered, error) {
	event := new(LatestValueHolderTriggered)
	if err := _LatestValueHolder.contract.UnpackLog(event, "Triggered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_LatestValueHolder *LatestValueHolder) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _LatestValueHolder.abi.Events["Triggered"].ID:
		return _LatestValueHolder.ParseTriggered(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (LatestValueHolderTriggered) Topic() common.Hash {
	return common.HexToHash("0x7188419dcd8b51877b71766f075f3626586c0ff190e7d056aa65ce9acb649a3d")
}

func (_LatestValueHolder *LatestValueHolder) Address() common.Address {
	return _LatestValueHolder.address
}

type LatestValueHolderInterface interface {
	GetDifferentPrimitiveValue(opts *bind.CallOpts) (uint64, error)

	GetElementAtIndex(opts *bind.CallOpts, i *big.Int) (TestStruct, error)

	GetPrimitiveValue(opts *bind.CallOpts) (uint64, error)

	GetSliceValue(opts *bind.CallOpts) ([]uint64, error)

	ReturnSeen(opts *bind.CallOpts, field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedStruct MidLevelTestStruct) (TestStruct, error)

	AddTestStruct(opts *bind.TransactOpts, field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedStruct MidLevelTestStruct) (*types.Transaction, error)

	TriggerEvent(opts *bind.TransactOpts, field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedStruct MidLevelTestStruct) (*types.Transaction, error)

	FilterTriggered(opts *bind.FilterOpts) (*LatestValueHolderTriggeredIterator, error)

	WatchTriggered(opts *bind.WatchOpts, sink chan<- *LatestValueHolderTriggered) (event.Subscription, error)

	ParseTriggered(log types.Log) (*LatestValueHolderTriggered, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
