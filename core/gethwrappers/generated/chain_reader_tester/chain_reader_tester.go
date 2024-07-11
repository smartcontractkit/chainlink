// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package chain_reader_tester

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
	IntVal int64
	S      string
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

var ChainReaderTesterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int32\",\"name\":\"field\",\"type\":\"int32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"differentField\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"oracleId\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint8[32]\",\"name\":\"oracleIds\",\"type\":\"uint8[32]\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"Account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"Accounts\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"int192\",\"name\":\"bigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"indexed\":false,\"internalType\":\"structMidLevelTestStruct\",\"name\":\"nestedStruct\",\"type\":\"tuple\"}],\"name\":\"Triggered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"fieldHash\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"field\",\"type\":\"string\"}],\"name\":\"TriggeredEventWithDynamicTopic\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int32\",\"name\":\"field1\",\"type\":\"int32\"},{\"indexed\":true,\"internalType\":\"int32\",\"name\":\"field2\",\"type\":\"int32\"},{\"indexed\":true,\"internalType\":\"int32\",\"name\":\"field3\",\"type\":\"int32\"}],\"name\":\"TriggeredWithFourTopics\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"int32\",\"name\":\"field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"differentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"oracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"oracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"int192\",\"name\":\"bigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelTestStruct\",\"name\":\"nestedStruct\",\"type\":\"tuple\"}],\"name\":\"addTestStruct\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAlterablePrimitiveValue\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDifferentPrimitiveValue\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"i\",\"type\":\"uint256\"}],\"name\":\"getElementAtIndex\",\"outputs\":[{\"components\":[{\"internalType\":\"int32\",\"name\":\"Field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"DifferentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"OracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"OracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"address\",\"name\":\"Account\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"Accounts\",\"type\":\"address[]\"},{\"internalType\":\"int192\",\"name\":\"BigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelTestStruct\",\"name\":\"NestedStruct\",\"type\":\"tuple\"}],\"internalType\":\"structTestStruct\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPrimitiveValue\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSliceValue\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int32\",\"name\":\"field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"differentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"oracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"oracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"int192\",\"name\":\"bigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelTestStruct\",\"name\":\"nestedStruct\",\"type\":\"tuple\"}],\"name\":\"returnSeen\",\"outputs\":[{\"components\":[{\"internalType\":\"int32\",\"name\":\"Field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"DifferentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"OracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"OracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"address\",\"name\":\"Account\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"Accounts\",\"type\":\"address[]\"},{\"internalType\":\"int192\",\"name\":\"BigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelTestStruct\",\"name\":\"NestedStruct\",\"type\":\"tuple\"}],\"internalType\":\"structTestStruct\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"value\",\"type\":\"uint64\"}],\"name\":\"setAlterablePrimitiveValue\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int32\",\"name\":\"field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"differentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"oracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"oracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"int192\",\"name\":\"bigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelTestStruct\",\"name\":\"nestedStruct\",\"type\":\"tuple\"}],\"name\":\"triggerEvent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"field\",\"type\":\"string\"}],\"name\":\"triggerEventWithDynamicTopic\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int32\",\"name\":\"field1\",\"type\":\"int32\"},{\"internalType\":\"int32\",\"name\":\"field2\",\"type\":\"int32\"},{\"internalType\":\"int32\",\"name\":\"field3\",\"type\":\"int32\"}],\"name\":\"triggerWithFourTopics\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50600180548082018255600082905260048082047fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6908101805460086003958616810261010090810a8088026001600160401b0391820219909416939093179093558654808801909755848704909301805496909516909202900a91820291021990921691909117905561181e806100a96000396000f3fe608060405234801561001057600080fd5b50600436106100c95760003560e01c80637f002d6711610081578063ef4e1ced1161005b578063ef4e1ced146101c0578063f6f871c8146101c7578063fbe9fbf6146101da57600080fd5b80637f002d671461017d578063ab5e0b3814610190578063dbfd7332146101ad57600080fd5b806349eac2ac116100b257806349eac2ac1461010c578063679004a41461011f5780636c9a43b61461013457600080fd5b80632c45576f146100ce5780633272b66c146100f7575b600080fd5b6100e16100dc366004610c2b565b6101ec565b6040516100ee9190610d8a565b60405180910390f35b61010a610105366004610ec9565b6104c7565b005b61010a61011a366004610fde565b61051c565b61012761081f565b6040516100ee91906110d0565b61010a61014236600461111e565b600280547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001667ffffffffffffffff92909216919091179055565b61010a61018b366004610fde565b6108ab565b6107c65b60405167ffffffffffffffff90911681526020016100ee565b61010a6101bb36600461114f565b610902565b6003610194565b6100e16101d5366004610fde565b61093f565b60025467ffffffffffffffff16610194565b6101f4610a48565b6000610201600184611192565b81548110610211576102116111cc565b6000918252602091829020604080516101008101909152600a90920201805460030b8252600181018054929391929184019161024c906111fb565b80601f0160208091040260200160405190810160405280929190818152602001828054610278906111fb565b80156102c55780601f1061029a576101008083540402835291602001916102c5565b820191906000526020600020905b8154815290600101906020018083116102a857829003601f168201915b5050509183525050600282015460ff166020808301919091526040805161040081018083529190930192916003850191826000855b825461010083900a900460ff168152602060019283018181049485019490930390920291018084116102fa57505050928452505050600482015473ffffffffffffffffffffffffffffffffffffffff1660208083019190915260058301805460408051828502810185018252828152940193928301828280156103b357602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610388575b5050509183525050600682015460170b6020808301919091526040805180820182526007808601805460f01b7fffff0000000000000000000000000000000000000000000000000000000000001683528351808501855260088801805490930b81526009880180549590970196939591948683019491939284019190610438906111fb565b80601f0160208091040260200160405190810160405280929190818152602001828054610464906111fb565b80156104b15780601f10610486576101008083540402835291602001916104b1565b820191906000526020600020905b81548152906001019060200180831161049457829003601f168201915b5050509190925250505090525090525092915050565b81816040516104d7929190611248565b60405180910390207f3d969732b1bbbb9f1d7eb9f3f14e4cb50a74d950b3ef916a397b85dfbab93c6783836040516105109291906112a1565b60405180910390a25050565b60006040518061010001604052808c60030b81526020018b8b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050509082525060ff8a166020808301919091526040805161040081810183529190930192918b9183908390808284376000920191909152505050815273ffffffffffffffffffffffffffffffffffffffff8816602080830191909152604080518883028181018401835289825291909301929189918991829190850190849080828437600092019190915250505090825250601785900b602082015260400161060e8461139e565b905281546001808201845560009384526020938490208351600a9093020180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff90931692909217825592820151919290919082019061067490826114f8565b5060408201516002820180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff90921691909117905560608201516106c29060038301906020610a97565b5060808201516004820180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff90921691909117905560a08201518051610729916005840191602090910190610b2a565b5060c08201516006820180547fffffffffffffffff0000000000000000000000000000000000000000000000001677ffffffffffffffffffffffffffffffffffffffffffffffff90921691909117905560e082015180516007830180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001660f09290921c91909117815560208083015180516008860180547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001667ffffffffffffffff90921691909117815591810151909190600986019061080c90826114f8565b5050505050505050505050505050505050565b606060018054806020026020016040519081016040528092919081815260200182805480156108a157602002820191906000526020600020906000905b82829054906101000a900467ffffffffffffffff1667ffffffffffffffff168152602001906008019060208260070104928301926001038202915080841161085c5790505b5050505050905090565b8960030b7f7188419dcd8b51877b71766f075f3626586c0ff190e7d056aa65ce9acb649a3d8a8a8a8a8a8a8a8a8a6040516108ee99989796959493929190611757565b60405180910390a250505050505050505050565b8060030b8260030b8460030b7f91c80dc390f3d041b3a04b0099b19634499541ea26972250986ee4b24a12fac560405160405180910390a4505050565b610947610a48565b6040518061010001604052808c60030b81526020018b8b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050509082525060ff8a166020808301919091526040805161040081810183529190930192918b9183908390808284376000920191909152505050815273ffffffffffffffffffffffffffffffffffffffff8816602080830191909152604080518883028181018401835289825291909301929189918991829190850190849080828437600092019190915250505090825250601785900b6020820152604001610a378461139e565b90529b9a5050505050505050505050565b6040805161010081018252600080825260606020830181905292820152908101610a70610ba4565b8152600060208201819052606060408301819052820152608001610a92610bc3565b905290565b600183019183908215610b1a5791602002820160005b83821115610aeb57835183826101000a81548160ff021916908360ff1602179055509260200192600101602081600001049283019260010302610aad565b8015610b185782816101000a81549060ff0219169055600101602081600001049283019260010302610aeb565b505b50610b26929150610c16565b5090565b828054828255906000526020600020908101928215610b1a579160200282015b82811115610b1a57825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190610b4a565b6040518061040001604052806020906020820280368337509192915050565b604051806040016040528060007dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168152602001610a926040518060400160405280600060070b8152602001606081525090565b5b80821115610b265760008155600101610c17565b600060208284031215610c3d57600080fd5b5035919050565b6000815180845260005b81811015610c6a57602081850181015186830182015201610c4e565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b8060005b6020808210610cbb5750610cd2565b825160ff1685529384019390910190600101610cac565b50505050565b600081518084526020808501945080840160005b83811015610d1e57815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101610cec565b509495945050505050565b7fffff00000000000000000000000000000000000000000000000000000000000081511682526000602082015160406020850152805160070b60408501526020810151905060406060850152610d826080850182610c44565b949350505050565b60208152610d9e60208201835160030b9052565b600060208301516104e0806040850152610dbc610500850183610c44565b91506040850151610dd2606086018260ff169052565b506060850151610de56080860182610ca8565b50608085015173ffffffffffffffffffffffffffffffffffffffff1661048085015260a08501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe085840381016104a0870152610e428483610cd8565b935060c08701519150610e5b6104c087018360170b9052565b60e0870151915080868503018387015250610e768382610d29565b9695505050505050565b60008083601f840112610e9257600080fd5b50813567ffffffffffffffff811115610eaa57600080fd5b602083019150836020828501011115610ec257600080fd5b9250929050565b60008060208385031215610edc57600080fd5b823567ffffffffffffffff811115610ef357600080fd5b610eff85828601610e80565b90969095509350505050565b8035600381900b8114610f1d57600080fd5b919050565b803560ff81168114610f1d57600080fd5b806104008101831015610f4557600080fd5b92915050565b803573ffffffffffffffffffffffffffffffffffffffff81168114610f1d57600080fd5b60008083601f840112610f8157600080fd5b50813567ffffffffffffffff811115610f9957600080fd5b6020830191508360208260051b8501011115610ec257600080fd5b8035601781900b8114610f1d57600080fd5b600060408284031215610fd857600080fd5b50919050565b6000806000806000806000806000806104e08b8d031215610ffe57600080fd5b6110078b610f0b565b995060208b013567ffffffffffffffff8082111561102457600080fd5b6110308e838f01610e80565b909b50995089915061104460408e01610f22565b98506110538e60608f01610f33565b97506110626104608e01610f4b565b96506104808d013591508082111561107957600080fd5b6110858e838f01610f6f565b909650945084915061109a6104a08e01610fb4565b93506104c08d01359150808211156110b157600080fd5b506110be8d828e01610fc6565b9150509295989b9194979a5092959850565b6020808252825182820181905260009190848201906040850190845b8181101561111257835167ffffffffffffffff16835292840192918401916001016110ec565b50909695505050505050565b60006020828403121561113057600080fd5b813567ffffffffffffffff8116811461114857600080fd5b9392505050565b60008060006060848603121561116457600080fd5b61116d84610f0b565b925061117b60208501610f0b565b915061118960408501610f0b565b90509250925092565b81810381811115610f45577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600181811c9082168061120f57607f821691505b602082108103610fd8577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b8183823760009101908152919050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b602081526000610d82602083018486611258565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040805190810167ffffffffffffffff81118282101715611307576113076112b5565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715611354576113546112b5565b604052919050565b80357fffff00000000000000000000000000000000000000000000000000000000000081168114610f1d57600080fd5b8035600781900b8114610f1d57600080fd5b6000604082360312156113b057600080fd5b6113b86112e4565b6113c18361135c565b815260208084013567ffffffffffffffff808211156113df57600080fd5b8186019150604082360312156113f457600080fd5b6113fc6112e4565b6114058361138c565b8152838301358281111561141857600080fd5b929092019136601f84011261142c57600080fd5b82358281111561143e5761143e6112b5565b61146e857fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8401160161130d565b9250808352368582860101111561148457600080fd5b8085850186850137600090830185015280840191909152918301919091525092915050565b601f8211156114f357600081815260208120601f850160051c810160208610156114d05750805b601f850160051c820191505b818110156114ef578281556001016114dc565b5050505b505050565b815167ffffffffffffffff811115611512576115126112b5565b6115268161152084546111fb565b846114a9565b602080601f83116001811461157957600084156115435750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b1785556114ef565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b828110156115c6578886015182559484019460019091019084016115a7565b508582101561160257878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b8183526000602080850194508260005b85811015610d1e5773ffffffffffffffffffffffffffffffffffffffff61164883610f4b565b1687529582019590820190600101611622565b7fffff0000000000000000000000000000000000000000000000000000000000006116858261135c565b168252600060208201357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc18336030181126116bf57600080fd5b6040602085015282016116d18161138c565b60070b604085015260208101357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe182360301811261170e57600080fd5b0160208101903567ffffffffffffffff81111561172a57600080fd5b80360382131561173957600080fd5b6040606086015261174e608086018284611258565b95945050505050565b60006104c080835261176c8184018c8e611258565b9050602060ff808c1682860152604085018b60005b848110156117a6578361179383610f22565b1683529184019190840190600101611781565b505050505073ffffffffffffffffffffffffffffffffffffffff88166104408401528281036104608401526117dc818789611612565b90506117ee61048084018660170b9052565b8281036104a0840152611801818561165b565b9c9b50505050505050505050505056fea164736f6c6343000813000a",
}

var ChainReaderTesterABI = ChainReaderTesterMetaData.ABI

var ChainReaderTesterBin = ChainReaderTesterMetaData.Bin

func DeployChainReaderTester(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ChainReaderTester, error) {
	parsed, err := ChainReaderTesterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ChainReaderTesterBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ChainReaderTester{address: address, abi: *parsed, ChainReaderTesterCaller: ChainReaderTesterCaller{contract: contract}, ChainReaderTesterTransactor: ChainReaderTesterTransactor{contract: contract}, ChainReaderTesterFilterer: ChainReaderTesterFilterer{contract: contract}}, nil
}

type ChainReaderTester struct {
	address common.Address
	abi     abi.ABI
	ChainReaderTesterCaller
	ChainReaderTesterTransactor
	ChainReaderTesterFilterer
}

type ChainReaderTesterCaller struct {
	contract *bind.BoundContract
}

type ChainReaderTesterTransactor struct {
	contract *bind.BoundContract
}

type ChainReaderTesterFilterer struct {
	contract *bind.BoundContract
}

type ChainReaderTesterSession struct {
	Contract     *ChainReaderTester
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type ChainReaderTesterCallerSession struct {
	Contract *ChainReaderTesterCaller
	CallOpts bind.CallOpts
}

type ChainReaderTesterTransactorSession struct {
	Contract     *ChainReaderTesterTransactor
	TransactOpts bind.TransactOpts
}

type ChainReaderTesterRaw struct {
	Contract *ChainReaderTester
}

type ChainReaderTesterCallerRaw struct {
	Contract *ChainReaderTesterCaller
}

type ChainReaderTesterTransactorRaw struct {
	Contract *ChainReaderTesterTransactor
}

func NewChainReaderTester(address common.Address, backend bind.ContractBackend) (*ChainReaderTester, error) {
	abi, err := abi.JSON(strings.NewReader(ChainReaderTesterABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindChainReaderTester(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ChainReaderTester{address: address, abi: abi, ChainReaderTesterCaller: ChainReaderTesterCaller{contract: contract}, ChainReaderTesterTransactor: ChainReaderTesterTransactor{contract: contract}, ChainReaderTesterFilterer: ChainReaderTesterFilterer{contract: contract}}, nil
}

func NewChainReaderTesterCaller(address common.Address, caller bind.ContractCaller) (*ChainReaderTesterCaller, error) {
	contract, err := bindChainReaderTester(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ChainReaderTesterCaller{contract: contract}, nil
}

func NewChainReaderTesterTransactor(address common.Address, transactor bind.ContractTransactor) (*ChainReaderTesterTransactor, error) {
	contract, err := bindChainReaderTester(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ChainReaderTesterTransactor{contract: contract}, nil
}

func NewChainReaderTesterFilterer(address common.Address, filterer bind.ContractFilterer) (*ChainReaderTesterFilterer, error) {
	contract, err := bindChainReaderTester(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ChainReaderTesterFilterer{contract: contract}, nil
}

func bindChainReaderTester(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ChainReaderTesterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_ChainReaderTester *ChainReaderTesterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ChainReaderTester.Contract.ChainReaderTesterCaller.contract.Call(opts, result, method, params...)
}

func (_ChainReaderTester *ChainReaderTesterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChainReaderTester.Contract.ChainReaderTesterTransactor.contract.Transfer(opts)
}

func (_ChainReaderTester *ChainReaderTesterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ChainReaderTester.Contract.ChainReaderTesterTransactor.contract.Transact(opts, method, params...)
}

func (_ChainReaderTester *ChainReaderTesterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ChainReaderTester.Contract.contract.Call(opts, result, method, params...)
}

func (_ChainReaderTester *ChainReaderTesterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChainReaderTester.Contract.contract.Transfer(opts)
}

func (_ChainReaderTester *ChainReaderTesterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ChainReaderTester.Contract.contract.Transact(opts, method, params...)
}

func (_ChainReaderTester *ChainReaderTesterCaller) GetAlterablePrimitiveValue(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _ChainReaderTester.contract.Call(opts, &out, "getAlterablePrimitiveValue")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_ChainReaderTester *ChainReaderTesterSession) GetAlterablePrimitiveValue() (uint64, error) {
	return _ChainReaderTester.Contract.GetAlterablePrimitiveValue(&_ChainReaderTester.CallOpts)
}

func (_ChainReaderTester *ChainReaderTesterCallerSession) GetAlterablePrimitiveValue() (uint64, error) {
	return _ChainReaderTester.Contract.GetAlterablePrimitiveValue(&_ChainReaderTester.CallOpts)
}

func (_ChainReaderTester *ChainReaderTesterCaller) GetDifferentPrimitiveValue(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _ChainReaderTester.contract.Call(opts, &out, "getDifferentPrimitiveValue")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_ChainReaderTester *ChainReaderTesterSession) GetDifferentPrimitiveValue() (uint64, error) {
	return _ChainReaderTester.Contract.GetDifferentPrimitiveValue(&_ChainReaderTester.CallOpts)
}

func (_ChainReaderTester *ChainReaderTesterCallerSession) GetDifferentPrimitiveValue() (uint64, error) {
	return _ChainReaderTester.Contract.GetDifferentPrimitiveValue(&_ChainReaderTester.CallOpts)
}

func (_ChainReaderTester *ChainReaderTesterCaller) GetElementAtIndex(opts *bind.CallOpts, i *big.Int) (TestStruct, error) {
	var out []interface{}
	err := _ChainReaderTester.contract.Call(opts, &out, "getElementAtIndex", i)

	if err != nil {
		return *new(TestStruct), err
	}

	out0 := *abi.ConvertType(out[0], new(TestStruct)).(*TestStruct)

	return out0, err

}

func (_ChainReaderTester *ChainReaderTesterSession) GetElementAtIndex(i *big.Int) (TestStruct, error) {
	return _ChainReaderTester.Contract.GetElementAtIndex(&_ChainReaderTester.CallOpts, i)
}

func (_ChainReaderTester *ChainReaderTesterCallerSession) GetElementAtIndex(i *big.Int) (TestStruct, error) {
	return _ChainReaderTester.Contract.GetElementAtIndex(&_ChainReaderTester.CallOpts, i)
}

func (_ChainReaderTester *ChainReaderTesterCaller) GetPrimitiveValue(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _ChainReaderTester.contract.Call(opts, &out, "getPrimitiveValue")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_ChainReaderTester *ChainReaderTesterSession) GetPrimitiveValue() (uint64, error) {
	return _ChainReaderTester.Contract.GetPrimitiveValue(&_ChainReaderTester.CallOpts)
}

func (_ChainReaderTester *ChainReaderTesterCallerSession) GetPrimitiveValue() (uint64, error) {
	return _ChainReaderTester.Contract.GetPrimitiveValue(&_ChainReaderTester.CallOpts)
}

func (_ChainReaderTester *ChainReaderTesterCaller) GetSliceValue(opts *bind.CallOpts) ([]uint64, error) {
	var out []interface{}
	err := _ChainReaderTester.contract.Call(opts, &out, "getSliceValue")

	if err != nil {
		return *new([]uint64), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint64)).(*[]uint64)

	return out0, err

}

func (_ChainReaderTester *ChainReaderTesterSession) GetSliceValue() ([]uint64, error) {
	return _ChainReaderTester.Contract.GetSliceValue(&_ChainReaderTester.CallOpts)
}

func (_ChainReaderTester *ChainReaderTesterCallerSession) GetSliceValue() ([]uint64, error) {
	return _ChainReaderTester.Contract.GetSliceValue(&_ChainReaderTester.CallOpts)
}

func (_ChainReaderTester *ChainReaderTesterCaller) ReturnSeen(opts *bind.CallOpts, field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedStruct MidLevelTestStruct) (TestStruct, error) {
	var out []interface{}
	err := _ChainReaderTester.contract.Call(opts, &out, "returnSeen", field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)

	if err != nil {
		return *new(TestStruct), err
	}

	out0 := *abi.ConvertType(out[0], new(TestStruct)).(*TestStruct)

	return out0, err

}

func (_ChainReaderTester *ChainReaderTesterSession) ReturnSeen(field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedStruct MidLevelTestStruct) (TestStruct, error) {
	return _ChainReaderTester.Contract.ReturnSeen(&_ChainReaderTester.CallOpts, field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)
}

func (_ChainReaderTester *ChainReaderTesterCallerSession) ReturnSeen(field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedStruct MidLevelTestStruct) (TestStruct, error) {
	return _ChainReaderTester.Contract.ReturnSeen(&_ChainReaderTester.CallOpts, field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)
}

func (_ChainReaderTester *ChainReaderTesterTransactor) AddTestStruct(opts *bind.TransactOpts, field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedStruct MidLevelTestStruct) (*types.Transaction, error) {
	return _ChainReaderTester.contract.Transact(opts, "addTestStruct", field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)
}

func (_ChainReaderTester *ChainReaderTesterSession) AddTestStruct(field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedStruct MidLevelTestStruct) (*types.Transaction, error) {
	return _ChainReaderTester.Contract.AddTestStruct(&_ChainReaderTester.TransactOpts, field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)
}

func (_ChainReaderTester *ChainReaderTesterTransactorSession) AddTestStruct(field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedStruct MidLevelTestStruct) (*types.Transaction, error) {
	return _ChainReaderTester.Contract.AddTestStruct(&_ChainReaderTester.TransactOpts, field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)
}

func (_ChainReaderTester *ChainReaderTesterTransactor) SetAlterablePrimitiveValue(opts *bind.TransactOpts, value uint64) (*types.Transaction, error) {
	return _ChainReaderTester.contract.Transact(opts, "setAlterablePrimitiveValue", value)
}

func (_ChainReaderTester *ChainReaderTesterSession) SetAlterablePrimitiveValue(value uint64) (*types.Transaction, error) {
	return _ChainReaderTester.Contract.SetAlterablePrimitiveValue(&_ChainReaderTester.TransactOpts, value)
}

func (_ChainReaderTester *ChainReaderTesterTransactorSession) SetAlterablePrimitiveValue(value uint64) (*types.Transaction, error) {
	return _ChainReaderTester.Contract.SetAlterablePrimitiveValue(&_ChainReaderTester.TransactOpts, value)
}

func (_ChainReaderTester *ChainReaderTesterTransactor) TriggerEvent(opts *bind.TransactOpts, field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedStruct MidLevelTestStruct) (*types.Transaction, error) {
	return _ChainReaderTester.contract.Transact(opts, "triggerEvent", field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)
}

func (_ChainReaderTester *ChainReaderTesterSession) TriggerEvent(field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedStruct MidLevelTestStruct) (*types.Transaction, error) {
	return _ChainReaderTester.Contract.TriggerEvent(&_ChainReaderTester.TransactOpts, field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)
}

func (_ChainReaderTester *ChainReaderTesterTransactorSession) TriggerEvent(field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedStruct MidLevelTestStruct) (*types.Transaction, error) {
	return _ChainReaderTester.Contract.TriggerEvent(&_ChainReaderTester.TransactOpts, field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)
}

func (_ChainReaderTester *ChainReaderTesterTransactor) TriggerEventWithDynamicTopic(opts *bind.TransactOpts, field string) (*types.Transaction, error) {
	return _ChainReaderTester.contract.Transact(opts, "triggerEventWithDynamicTopic", field)
}

func (_ChainReaderTester *ChainReaderTesterSession) TriggerEventWithDynamicTopic(field string) (*types.Transaction, error) {
	return _ChainReaderTester.Contract.TriggerEventWithDynamicTopic(&_ChainReaderTester.TransactOpts, field)
}

func (_ChainReaderTester *ChainReaderTesterTransactorSession) TriggerEventWithDynamicTopic(field string) (*types.Transaction, error) {
	return _ChainReaderTester.Contract.TriggerEventWithDynamicTopic(&_ChainReaderTester.TransactOpts, field)
}

func (_ChainReaderTester *ChainReaderTesterTransactor) TriggerWithFourTopics(opts *bind.TransactOpts, field1 int32, field2 int32, field3 int32) (*types.Transaction, error) {
	return _ChainReaderTester.contract.Transact(opts, "triggerWithFourTopics", field1, field2, field3)
}

func (_ChainReaderTester *ChainReaderTesterSession) TriggerWithFourTopics(field1 int32, field2 int32, field3 int32) (*types.Transaction, error) {
	return _ChainReaderTester.Contract.TriggerWithFourTopics(&_ChainReaderTester.TransactOpts, field1, field2, field3)
}

func (_ChainReaderTester *ChainReaderTesterTransactorSession) TriggerWithFourTopics(field1 int32, field2 int32, field3 int32) (*types.Transaction, error) {
	return _ChainReaderTester.Contract.TriggerWithFourTopics(&_ChainReaderTester.TransactOpts, field1, field2, field3)
}

type ChainReaderTesterTriggeredIterator struct {
	Event *ChainReaderTesterTriggered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ChainReaderTesterTriggeredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChainReaderTesterTriggered)
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
		it.Event = new(ChainReaderTesterTriggered)
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

func (it *ChainReaderTesterTriggeredIterator) Error() error {
	return it.fail
}

func (it *ChainReaderTesterTriggeredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ChainReaderTesterTriggered struct {
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

func (_ChainReaderTester *ChainReaderTesterFilterer) FilterTriggered(opts *bind.FilterOpts, field []int32) (*ChainReaderTesterTriggeredIterator, error) {

	var fieldRule []interface{}
	for _, fieldItem := range field {
		fieldRule = append(fieldRule, fieldItem)
	}

	logs, sub, err := _ChainReaderTester.contract.FilterLogs(opts, "Triggered", fieldRule)
	if err != nil {
		return nil, err
	}
	return &ChainReaderTesterTriggeredIterator{contract: _ChainReaderTester.contract, event: "Triggered", logs: logs, sub: sub}, nil
}

func (_ChainReaderTester *ChainReaderTesterFilterer) WatchTriggered(opts *bind.WatchOpts, sink chan<- *ChainReaderTesterTriggered, field []int32) (event.Subscription, error) {

	var fieldRule []interface{}
	for _, fieldItem := range field {
		fieldRule = append(fieldRule, fieldItem)
	}

	logs, sub, err := _ChainReaderTester.contract.WatchLogs(opts, "Triggered", fieldRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ChainReaderTesterTriggered)
				if err := _ChainReaderTester.contract.UnpackLog(event, "Triggered", log); err != nil {
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

func (_ChainReaderTester *ChainReaderTesterFilterer) ParseTriggered(log types.Log) (*ChainReaderTesterTriggered, error) {
	event := new(ChainReaderTesterTriggered)
	if err := _ChainReaderTester.contract.UnpackLog(event, "Triggered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ChainReaderTesterTriggeredEventWithDynamicTopicIterator struct {
	Event *ChainReaderTesterTriggeredEventWithDynamicTopic

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ChainReaderTesterTriggeredEventWithDynamicTopicIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChainReaderTesterTriggeredEventWithDynamicTopic)
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
		it.Event = new(ChainReaderTesterTriggeredEventWithDynamicTopic)
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

func (it *ChainReaderTesterTriggeredEventWithDynamicTopicIterator) Error() error {
	return it.fail
}

func (it *ChainReaderTesterTriggeredEventWithDynamicTopicIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ChainReaderTesterTriggeredEventWithDynamicTopic struct {
	FieldHash common.Hash
	Field     string
	Raw       types.Log
}

func (_ChainReaderTester *ChainReaderTesterFilterer) FilterTriggeredEventWithDynamicTopic(opts *bind.FilterOpts, fieldHash []string) (*ChainReaderTesterTriggeredEventWithDynamicTopicIterator, error) {

	var fieldHashRule []interface{}
	for _, fieldHashItem := range fieldHash {
		fieldHashRule = append(fieldHashRule, fieldHashItem)
	}

	logs, sub, err := _ChainReaderTester.contract.FilterLogs(opts, "TriggeredEventWithDynamicTopic", fieldHashRule)
	if err != nil {
		return nil, err
	}
	return &ChainReaderTesterTriggeredEventWithDynamicTopicIterator{contract: _ChainReaderTester.contract, event: "TriggeredEventWithDynamicTopic", logs: logs, sub: sub}, nil
}

func (_ChainReaderTester *ChainReaderTesterFilterer) WatchTriggeredEventWithDynamicTopic(opts *bind.WatchOpts, sink chan<- *ChainReaderTesterTriggeredEventWithDynamicTopic, fieldHash []string) (event.Subscription, error) {

	var fieldHashRule []interface{}
	for _, fieldHashItem := range fieldHash {
		fieldHashRule = append(fieldHashRule, fieldHashItem)
	}

	logs, sub, err := _ChainReaderTester.contract.WatchLogs(opts, "TriggeredEventWithDynamicTopic", fieldHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ChainReaderTesterTriggeredEventWithDynamicTopic)
				if err := _ChainReaderTester.contract.UnpackLog(event, "TriggeredEventWithDynamicTopic", log); err != nil {
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

func (_ChainReaderTester *ChainReaderTesterFilterer) ParseTriggeredEventWithDynamicTopic(log types.Log) (*ChainReaderTesterTriggeredEventWithDynamicTopic, error) {
	event := new(ChainReaderTesterTriggeredEventWithDynamicTopic)
	if err := _ChainReaderTester.contract.UnpackLog(event, "TriggeredEventWithDynamicTopic", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ChainReaderTesterTriggeredWithFourTopicsIterator struct {
	Event *ChainReaderTesterTriggeredWithFourTopics

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ChainReaderTesterTriggeredWithFourTopicsIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChainReaderTesterTriggeredWithFourTopics)
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
		it.Event = new(ChainReaderTesterTriggeredWithFourTopics)
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

func (it *ChainReaderTesterTriggeredWithFourTopicsIterator) Error() error {
	return it.fail
}

func (it *ChainReaderTesterTriggeredWithFourTopicsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ChainReaderTesterTriggeredWithFourTopics struct {
	Field1 int32
	Field2 int32
	Field3 int32
	Raw    types.Log
}

func (_ChainReaderTester *ChainReaderTesterFilterer) FilterTriggeredWithFourTopics(opts *bind.FilterOpts, field1 []int32, field2 []int32, field3 []int32) (*ChainReaderTesterTriggeredWithFourTopicsIterator, error) {

	var field1Rule []interface{}
	for _, field1Item := range field1 {
		field1Rule = append(field1Rule, field1Item)
	}
	var field2Rule []interface{}
	for _, field2Item := range field2 {
		field2Rule = append(field2Rule, field2Item)
	}
	var field3Rule []interface{}
	for _, field3Item := range field3 {
		field3Rule = append(field3Rule, field3Item)
	}

	logs, sub, err := _ChainReaderTester.contract.FilterLogs(opts, "TriggeredWithFourTopics", field1Rule, field2Rule, field3Rule)
	if err != nil {
		return nil, err
	}
	return &ChainReaderTesterTriggeredWithFourTopicsIterator{contract: _ChainReaderTester.contract, event: "TriggeredWithFourTopics", logs: logs, sub: sub}, nil
}

func (_ChainReaderTester *ChainReaderTesterFilterer) WatchTriggeredWithFourTopics(opts *bind.WatchOpts, sink chan<- *ChainReaderTesterTriggeredWithFourTopics, field1 []int32, field2 []int32, field3 []int32) (event.Subscription, error) {

	var field1Rule []interface{}
	for _, field1Item := range field1 {
		field1Rule = append(field1Rule, field1Item)
	}
	var field2Rule []interface{}
	for _, field2Item := range field2 {
		field2Rule = append(field2Rule, field2Item)
	}
	var field3Rule []interface{}
	for _, field3Item := range field3 {
		field3Rule = append(field3Rule, field3Item)
	}

	logs, sub, err := _ChainReaderTester.contract.WatchLogs(opts, "TriggeredWithFourTopics", field1Rule, field2Rule, field3Rule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ChainReaderTesterTriggeredWithFourTopics)
				if err := _ChainReaderTester.contract.UnpackLog(event, "TriggeredWithFourTopics", log); err != nil {
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

func (_ChainReaderTester *ChainReaderTesterFilterer) ParseTriggeredWithFourTopics(log types.Log) (*ChainReaderTesterTriggeredWithFourTopics, error) {
	event := new(ChainReaderTesterTriggeredWithFourTopics)
	if err := _ChainReaderTester.contract.UnpackLog(event, "TriggeredWithFourTopics", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_ChainReaderTester *ChainReaderTester) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _ChainReaderTester.abi.Events["Triggered"].ID:
		return _ChainReaderTester.ParseTriggered(log)
	case _ChainReaderTester.abi.Events["TriggeredEventWithDynamicTopic"].ID:
		return _ChainReaderTester.ParseTriggeredEventWithDynamicTopic(log)
	case _ChainReaderTester.abi.Events["TriggeredWithFourTopics"].ID:
		return _ChainReaderTester.ParseTriggeredWithFourTopics(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (ChainReaderTesterTriggered) Topic() common.Hash {
	return common.HexToHash("0x7188419dcd8b51877b71766f075f3626586c0ff190e7d056aa65ce9acb649a3d")
}

func (ChainReaderTesterTriggeredEventWithDynamicTopic) Topic() common.Hash {
	return common.HexToHash("0x3d969732b1bbbb9f1d7eb9f3f14e4cb50a74d950b3ef916a397b85dfbab93c67")
}

func (ChainReaderTesterTriggeredWithFourTopics) Topic() common.Hash {
	return common.HexToHash("0x91c80dc390f3d041b3a04b0099b19634499541ea26972250986ee4b24a12fac5")
}

func (_ChainReaderTester *ChainReaderTester) Address() common.Address {
	return _ChainReaderTester.address
}

type ChainReaderTesterInterface interface {
	GetAlterablePrimitiveValue(opts *bind.CallOpts) (uint64, error)

	GetDifferentPrimitiveValue(opts *bind.CallOpts) (uint64, error)

	GetElementAtIndex(opts *bind.CallOpts, i *big.Int) (TestStruct, error)

	GetPrimitiveValue(opts *bind.CallOpts) (uint64, error)

	GetSliceValue(opts *bind.CallOpts) ([]uint64, error)

	ReturnSeen(opts *bind.CallOpts, field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedStruct MidLevelTestStruct) (TestStruct, error)

	AddTestStruct(opts *bind.TransactOpts, field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedStruct MidLevelTestStruct) (*types.Transaction, error)

	SetAlterablePrimitiveValue(opts *bind.TransactOpts, value uint64) (*types.Transaction, error)

	TriggerEvent(opts *bind.TransactOpts, field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedStruct MidLevelTestStruct) (*types.Transaction, error)

	TriggerEventWithDynamicTopic(opts *bind.TransactOpts, field string) (*types.Transaction, error)

	TriggerWithFourTopics(opts *bind.TransactOpts, field1 int32, field2 int32, field3 int32) (*types.Transaction, error)

	FilterTriggered(opts *bind.FilterOpts, field []int32) (*ChainReaderTesterTriggeredIterator, error)

	WatchTriggered(opts *bind.WatchOpts, sink chan<- *ChainReaderTesterTriggered, field []int32) (event.Subscription, error)

	ParseTriggered(log types.Log) (*ChainReaderTesterTriggered, error)

	FilterTriggeredEventWithDynamicTopic(opts *bind.FilterOpts, fieldHash []string) (*ChainReaderTesterTriggeredEventWithDynamicTopicIterator, error)

	WatchTriggeredEventWithDynamicTopic(opts *bind.WatchOpts, sink chan<- *ChainReaderTesterTriggeredEventWithDynamicTopic, fieldHash []string) (event.Subscription, error)

	ParseTriggeredEventWithDynamicTopic(log types.Log) (*ChainReaderTesterTriggeredEventWithDynamicTopic, error)

	FilterTriggeredWithFourTopics(opts *bind.FilterOpts, field1 []int32, field2 []int32, field3 []int32) (*ChainReaderTesterTriggeredWithFourTopicsIterator, error)

	WatchTriggeredWithFourTopics(opts *bind.WatchOpts, sink chan<- *ChainReaderTesterTriggeredWithFourTopics, field1 []int32, field2 []int32, field3 []int32) (event.Subscription, error)

	ParseTriggeredWithFourTopics(log types.Log) (*ChainReaderTesterTriggeredWithFourTopics, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
