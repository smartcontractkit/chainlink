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

type InnerDynamicTestStruct struct {
	IntVal int64
	S      string
}

type InnerStaticTestStruct struct {
	IntVal int64
	A      common.Address
}

type MidLevelDynamicTestStruct struct {
	FixedBytes [2]byte
	Inner      InnerDynamicTestStruct
}

type MidLevelStaticTestStruct struct {
	FixedBytes [2]byte
	Inner      InnerStaticTestStruct
}

type TestStruct struct {
	Field               int32
	DifferentField      string
	OracleId            uint8
	OracleIds           [32]uint8
	Account             common.Address
	Accounts            []common.Address
	BigField            *big.Int
	NestedDynamicStruct MidLevelDynamicTestStruct
	NestedStaticStruct  MidLevelStaticTestStruct
}

var ChainReaderTesterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int32\",\"name\":\"field\",\"type\":\"int32\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"oracleId\",\"type\":\"uint8\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerDynamicTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"indexed\":false,\"internalType\":\"structMidLevelDynamicTestStruct\",\"name\":\"nestedDynamicStruct\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"address\",\"name\":\"A\",\"type\":\"address\"}],\"internalType\":\"structInnerStaticTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"indexed\":false,\"internalType\":\"structMidLevelStaticTestStruct\",\"name\":\"nestedStaticStruct\",\"type\":\"tuple\"},{\"indexed\":false,\"internalType\":\"uint8[32]\",\"name\":\"oracleIds\",\"type\":\"uint8[32]\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"Account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"Accounts\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"differentField\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"int192\",\"name\":\"bigField\",\"type\":\"int192\"}],\"name\":\"Triggered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"fieldHash\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"field\",\"type\":\"string\"}],\"name\":\"TriggeredEventWithDynamicTopic\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int32\",\"name\":\"field1\",\"type\":\"int32\"},{\"indexed\":true,\"internalType\":\"int32\",\"name\":\"field2\",\"type\":\"int32\"},{\"indexed\":true,\"internalType\":\"int32\",\"name\":\"field3\",\"type\":\"int32\"}],\"name\":\"TriggeredWithFourTopics\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"field1\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"uint8[32]\",\"name\":\"field2\",\"type\":\"uint8[32]\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"field3\",\"type\":\"bytes32\"}],\"name\":\"TriggeredWithFourTopicsWithHashed\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"int32\",\"name\":\"field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"differentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"oracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"oracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"int192\",\"name\":\"bigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerDynamicTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelDynamicTestStruct\",\"name\":\"nestedDynamicStruct\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"address\",\"name\":\"A\",\"type\":\"address\"}],\"internalType\":\"structInnerStaticTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelStaticTestStruct\",\"name\":\"nestedStaticStruct\",\"type\":\"tuple\"}],\"name\":\"addTestStruct\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAlterablePrimitiveValue\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDifferentPrimitiveValue\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"i\",\"type\":\"uint256\"}],\"name\":\"getElementAtIndex\",\"outputs\":[{\"components\":[{\"internalType\":\"int32\",\"name\":\"Field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"DifferentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"OracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"OracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"address\",\"name\":\"Account\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"Accounts\",\"type\":\"address[]\"},{\"internalType\":\"int192\",\"name\":\"BigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerDynamicTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelDynamicTestStruct\",\"name\":\"NestedDynamicStruct\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"address\",\"name\":\"A\",\"type\":\"address\"}],\"internalType\":\"structInnerStaticTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelStaticTestStruct\",\"name\":\"NestedStaticStruct\",\"type\":\"tuple\"}],\"internalType\":\"structTestStruct\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPrimitiveValue\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSliceValue\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int32\",\"name\":\"field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"differentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"oracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"oracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"int192\",\"name\":\"bigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerDynamicTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelDynamicTestStruct\",\"name\":\"nestedDynamicStruct\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"address\",\"name\":\"A\",\"type\":\"address\"}],\"internalType\":\"structInnerStaticTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelStaticTestStruct\",\"name\":\"nestedStaticStruct\",\"type\":\"tuple\"}],\"name\":\"returnSeen\",\"outputs\":[{\"components\":[{\"internalType\":\"int32\",\"name\":\"Field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"DifferentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"OracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"OracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"address\",\"name\":\"Account\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"Accounts\",\"type\":\"address[]\"},{\"internalType\":\"int192\",\"name\":\"BigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerDynamicTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelDynamicTestStruct\",\"name\":\"NestedDynamicStruct\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"address\",\"name\":\"A\",\"type\":\"address\"}],\"internalType\":\"structInnerStaticTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelStaticTestStruct\",\"name\":\"NestedStaticStruct\",\"type\":\"tuple\"}],\"internalType\":\"structTestStruct\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"value\",\"type\":\"uint64\"}],\"name\":\"setAlterablePrimitiveValue\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int32\",\"name\":\"field\",\"type\":\"int32\"},{\"internalType\":\"uint8\",\"name\":\"oracleId\",\"type\":\"uint8\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerDynamicTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelDynamicTestStruct\",\"name\":\"nestedDynamicStruct\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"address\",\"name\":\"A\",\"type\":\"address\"}],\"internalType\":\"structInnerStaticTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelStaticTestStruct\",\"name\":\"nestedStaticStruct\",\"type\":\"tuple\"},{\"internalType\":\"uint8[32]\",\"name\":\"oracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"string\",\"name\":\"differentField\",\"type\":\"string\"},{\"internalType\":\"int192\",\"name\":\"bigField\",\"type\":\"int192\"}],\"name\":\"triggerEvent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"field\",\"type\":\"string\"}],\"name\":\"triggerEventWithDynamicTopic\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int32\",\"name\":\"field1\",\"type\":\"int32\"},{\"internalType\":\"int32\",\"name\":\"field2\",\"type\":\"int32\"},{\"internalType\":\"int32\",\"name\":\"field3\",\"type\":\"int32\"}],\"name\":\"triggerWithFourTopics\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"field1\",\"type\":\"string\"},{\"internalType\":\"uint8[32]\",\"name\":\"field2\",\"type\":\"uint8[32]\"},{\"internalType\":\"bytes32\",\"name\":\"field3\",\"type\":\"bytes32\"}],\"name\":\"triggerWithFourTopicsWithHashed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50600180548082018255600082905260048082047fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6908101805460086003958616810261010090810a8088026001600160401b0391820219909416939093179093558654808801909755848704909301805496909516909202900a918202910219909216919091179055611d9a806100a96000396000f3fe608060405234801561001057600080fd5b50600436106100d45760003560e01c80636c9a43b611610081578063dbfd73321161005b578063dbfd7332146101de578063ef4e1ced146101f1578063fbe9fbf6146101f857600080fd5b80636c9a43b614610165578063a90e1998146101ae578063ab5e0b38146101c157600080fd5b80634149667f116100b25780634149667f1461012a5780635f7104a21461013d578063679004a41461015057600080fd5b80631b48259e146100d95780632c45576f146100ee5780633272b66c14610117575b600080fd5b6100ec6100e7366004610f5d565b61020a565b005b6101016100fc366004611061565b610264565b60405161010e91906111ca565b60405180910390f35b6100ec610125366004611323565b6105bc565b610101610138366004611365565b610611565b6100ec61014b366004611365565b610731565b610158610af6565b60405161010e9190611458565b6100ec6101733660046114a6565b600280547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001667ffffffffffffffff92909216919091179055565b6100ec6101bc3660046115da565b610b82565b6107c65b60405167ffffffffffffffff909116815260200161010e565b6100ec6101ec36600461168f565b610bdc565b60036101c5565b60025467ffffffffffffffff166101c5565b8a60030b7fae927edae02672fdcce7d7e8cf34c611ed3856914a159df5f2a59307b767c25b8b8b8b8b8b8b8b8b8b8b60405161024f9a99989796959493929190611844565b60405180910390a25050505050505050505050565b61026c610c19565b60006102796001846119c4565b81548110610289576102896119fe565b6000918252602091829020604080516101208101909152600c90920201805460030b825260018101805492939192918401916102c490611a2d565b80601f01602080910402602001604051908101604052809291908181526020018280546102f090611a2d565b801561033d5780601f106103125761010080835404028352916020019161033d565b820191906000526020600020905b81548152906001019060200180831161032057829003601f168201915b5050509183525050600282015460ff166020808301919091526040805161040081018083529190930192916003850191826000855b825461010083900a900460ff1681526020600192830181810494850194909303909202910180841161037257505050928452505050600482015473ffffffffffffffffffffffffffffffffffffffff16602080830191909152600583018054604080518285028101850182528281529401939283018282801561042b57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610400575b5050509183525050600682015460170b6020808301919091526040805180820182526007808601805460f01b7fffff0000000000000000000000000000000000000000000000000000000000001683528351808501855260088801805490930b815260098801805495909701969395919486830194919392840191906104b090611a2d565b80601f01602080910402602001604051908101604052809291908181526020018280546104dc90611a2d565b80156105295780601f106104fe57610100808354040283529160200191610529565b820191906000526020600020905b81548152906001019060200180831161050c57829003601f168201915b505050919092525050509052508152604080518082018252600a84015460f01b7fffff0000000000000000000000000000000000000000000000000000000000001681528151808301909252600b90930154600781900b825268010000000000000000900473ffffffffffffffffffffffffffffffffffffffff1660208083019190915280840191909152015292915050565b81816040516105cc929190611a7a565b60405180910390207f3d969732b1bbbb9f1d7eb9f3f14e4cb50a74d950b3ef916a397b85dfbab93c678383604051610605929190611a8a565b60405180910390a25050565b610619610c19565b6040518061012001604052808d60030b81526020018c8c8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050509082525060ff8b166020808301919091526040805161040081810183529190930192918c9183908390808284376000920191909152505050815273ffffffffffffffffffffffffffffffffffffffff891660208083019190915260408051898302818101840183528a82529190930192918a918a91829190850190849080828437600092019190915250505090825250601786900b602082015260400161070985611a9e565b815260200161071d36859003850185611b3c565b905290505b9b9a5050505050505050505050565b60006040518061012001604052808d60030b81526020018c8c8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050509082525060ff8b166020808301919091526040805161040081810183529190930192918c9183908390808284376000920191909152505050815273ffffffffffffffffffffffffffffffffffffffff891660208083019190915260408051898302818101840183528a82529190930192918a918a91829190850190849080828437600092019190915250505090825250601786900b602082015260400161082385611a9e565b815260200161083736859003850185611b3c565b905281546001808201845560009384526020938490208351600c9093020180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff90931692909217825592820151919290919082019061089d9082611c1b565b5060408201516002820180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff90921691909117905560608201516108eb9060038301906020610c9b565b5060808201516004820180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff90921691909117905560a08201518051610952916005840191602090910190610d2e565b5060c08201516006820180547fffffffffffffffff0000000000000000000000000000000000000000000000001677ffffffffffffffffffffffffffffffffffffffffffffffff90921691909117905560e082015180516007830180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001660f09290921c91909117815560208083015180516008860180547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001667ffffffffffffffff909216919091178155918101519091906009860190610a359082611c1b565b5050505061010092909201518051600a8301805460f09290921c7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00009092169190911790556020908101518051600b9093018054919092015173ffffffffffffffffffffffffffffffffffffffff1668010000000000000000027fffffffff0000000000000000000000000000000000000000000000000000000090911667ffffffffffffffff90931692909217919091179055505050505050505050505050565b60606001805480602002602001604051908101604052809291908181526020018280548015610b7857602002820191906000526020600020906000905b82829054906101000a900467ffffffffffffffff1667ffffffffffffffff1681526020019060080190602082600701049283019260010382029150808411610b335790505b5050505050905090565b8082604051610b919190611d35565b604051809103902084604051610ba79190611d71565b604051908190038120907f7220e4dbe4e9d0ed5f71acd022bc89c26748ac6784f2c548bc17bb8e52af34b090600090a4505050565b8060030b8260030b8460030b7f91c80dc390f3d041b3a04b0099b19634499541ea26972250986ee4b24a12fac560405160405180910390a4505050565b6040805161012081018252600080825260606020830181905292820152908101610c41610da8565b8152600060208201819052606060408301819052820152608001610c63610dc7565b8152602001610c966040805180820182526000808252825180840190935280835260208381019190915290919082015290565b905290565b600183019183908215610d1e5791602002820160005b83821115610cef57835183826101000a81548160ff021916908360ff1602179055509260200192600101602081600001049283019260010302610cb1565b8015610d1c5782816101000a81549060ff0219169055600101602081600001049283019260010302610cef565b505b50610d2a929150610e1a565b5090565b828054828255906000526020600020908101928215610d1e579160200282015b82811115610d1e57825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190610d4e565b6040518061040001604052806020906020820280368337509192915050565b604051806040016040528060007dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168152602001610c966040518060400160405280600060070b8152602001606081525090565b5b80821115610d2a5760008155600101610e1b565b8035600381900b8114610e4157600080fd5b919050565b803560ff81168114610e4157600080fd5b600060408284031215610e6957600080fd5b50919050565b600060608284031215610e6957600080fd5b806104008101831015610e9357600080fd5b92915050565b803573ffffffffffffffffffffffffffffffffffffffff81168114610e4157600080fd5b60008083601f840112610ecf57600080fd5b50813567ffffffffffffffff811115610ee757600080fd5b6020830191508360208260051b8501011115610f0257600080fd5b9250929050565b60008083601f840112610f1b57600080fd5b50813567ffffffffffffffff811115610f3357600080fd5b602083019150836020828501011115610f0257600080fd5b8035601781900b8114610e4157600080fd5b60008060008060008060008060008060006105408c8e031215610f7f57600080fd5b610f888c610e2f565b9a50610f9660208d01610e46565b995067ffffffffffffffff8060408e01351115610fb257600080fd5b610fc28e60408f01358f01610e57565b9950610fd18e60608f01610e6f565b9850610fe08e60c08f01610e81565b9750610fef6104c08e01610e99565b9650806104e08e0135111561100357600080fd5b6110148e6104e08f01358f01610ebd565b90965094506105008d013581101561102b57600080fd5b5061103d8d6105008e01358e01610f09565b909350915061104f6105208d01610f4b565b90509295989b509295989b9093969950565b60006020828403121561107357600080fd5b5035919050565b60005b8381101561109557818101518382015260200161107d565b50506000910152565b600081518084526110b681602086016020860161107a565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b8060005b60208082106110fb5750611112565b825160ff16855293840193909101906001016110ec565b50505050565b600081518084526020808501945080840160005b8381101561115e57815173ffffffffffffffffffffffffffffffffffffffff168752958201959082019060010161112c565b509495945050505050565b7fffff00000000000000000000000000000000000000000000000000000000000081511682526000602082015160406020850152805160070b604085015260208101519050604060608501526111c2608085018261109e565b949350505050565b602081526111de60208201835160030b9052565b6000602083015161054060408401526111fb61056084018261109e565b90506040840151611211606085018260ff169052565b50606084015161122460808501826110e8565b50608084015173ffffffffffffffffffffffffffffffffffffffff1661048084015260a08401517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe084830381016104a08601526112818383611118565b925060c0860151915061129a6104c086018360170b9052565b60e0860151915080858403016104e0860152506112b78282611169565b61010086015180517fffff00000000000000000000000000000000000000000000000000000000000016610500870152602080820151805160070b610520890152015173ffffffffffffffffffffffffffffffffffffffff166105408701529092509050509392505050565b6000806020838503121561133657600080fd5b823567ffffffffffffffff81111561134d57600080fd5b61135985828601610f09565b90969095509350505050565b60008060008060008060008060008060006105408c8e03121561138757600080fd5b6113908c610e2f565b9a5067ffffffffffffffff8060208e013511156113ac57600080fd5b6113bc8e60208f01358f01610f09565b909b5099506113cd60408e01610e46565b98506113dc8e60608f01610e81565b97506113eb6104608e01610e99565b9650806104808e013511156113ff57600080fd5b6114108e6104808f01358f01610ebd565b90965094506114226104a08e01610f4b565b9350806104c08e0135111561143657600080fd5b506114488d6104c08e01358e01610e57565b915061104f8d6104e08e01610e6f565b6020808252825182820181905260009190848201906040850190845b8181101561149a57835167ffffffffffffffff1683529284019291840191600101611474565b50909695505050505050565b6000602082840312156114b857600080fd5b813567ffffffffffffffff811681146114d057600080fd5b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040805190810167ffffffffffffffff81118282101715611529576115296114d7565b60405290565b600082601f83011261154057600080fd5b813567ffffffffffffffff8082111561155b5761155b6114d7565b604051601f83017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f011681019082821181831017156115a1576115a16114d7565b816040528381528660208588010111156115ba57600080fd5b836020870160208301376000602085830101528094505050505092915050565b600080600061044084860312156115f057600080fd5b833567ffffffffffffffff8082111561160857600080fd5b6116148783880161152f565b94506020915086603f87011261162957600080fd5b60405161040081018181108382111715611645576116456114d7565b60405290508061042087018881111561165d57600080fd5b8388015b8181101561167f5761167281610e46565b8452928401928401611661565b5095989097509435955050505050565b6000806000606084860312156116a457600080fd5b6116ad84610e2f565b92506116bb60208501610e2f565b91506116c960408501610e2f565b90509250925092565b80357fffff00000000000000000000000000000000000000000000000000000000000081168114610e4157600080fd5b8035600781900b8114610e4157600080fd5b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b7fffff000000000000000000000000000000000000000000000000000000000000611787826116d2565b16825261179660208201611702565b60070b602083015273ffffffffffffffffffffffffffffffffffffffff6117bf60408301610e99565b1660408301525050565b8060005b60208082106117dc5750611112565b60ff6117e784610e46565b1685529384019391909101906001016117cd565b8183526000602080850194508260005b8581101561115e5773ffffffffffffffffffffffffffffffffffffffff61183183610e99565b168752958201959082019060010161180b565b600061052060ff8d1683528060208401527fffff00000000000000000000000000000000000000000000000000000000000061187f8d6116d2565b16818401525060208b01357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc18c36030181126118ba57600080fd5b60406105408401528b016118cd81611702565b60070b61056084015260208101357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe182360301811261190b57600080fd5b0160208101903567ffffffffffffffff81111561192757600080fd5b80360382131561193657600080fd5b604061058085015261194d6105a085018284611714565b91505061195d604084018c61175d565b61196a60a084018b6117c9565b73ffffffffffffffffffffffffffffffffffffffff89166104a08401528281036104c084015261199b81888a6117fb565b90508281036104e08401526119b1818688611714565b91505061072261050083018460170b9052565b81810381811115610e93577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600181811c90821680611a4157607f821691505b602082108103610e69577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b8183823760009101908152919050565b6020815260006111c2602083018486611714565b600060408236031215611ab057600080fd5b611ab8611506565b611ac1836116d2565b8152602083013567ffffffffffffffff80821115611ade57600080fd5b818501915060408236031215611af357600080fd5b611afb611506565b611b0483611702565b8152602083013582811115611b1857600080fd5b611b243682860161152f565b60208301525080602085015250505080915050919050565b60008183036060811215611b4f57600080fd5b611b57611506565b611b60846116d2565b815260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe083011215611b9257600080fd5b611b9a611506565b9150611ba860208501611702565b8252611bb660408501610e99565b6020830152816020820152809250505092915050565b601f821115611c1657600081815260208120601f850160051c81016020861015611bf35750805b601f850160051c820191505b81811015611c1257828155600101611bff565b5050505b505050565b815167ffffffffffffffff811115611c3557611c356114d7565b611c4981611c438454611a2d565b84611bcc565b602080601f831160018114611c9c5760008415611c665750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555611c12565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015611ce957888601518255948401946001909101908401611cca565b5085821015611d2557878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b60008183825b6020808210611d4a5750611d61565b825160ff1684529283019290910190600101611d3b565b5050506104008201905092915050565b60008251611d8381846020870161107a565b919091019291505056fea164736f6c6343000813000a",
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

func (_ChainReaderTester *ChainReaderTesterCaller) ReturnSeen(opts *bind.CallOpts, field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedDynamicStruct MidLevelDynamicTestStruct, nestedStaticStruct MidLevelStaticTestStruct) (TestStruct, error) {
	var out []interface{}
	err := _ChainReaderTester.contract.Call(opts, &out, "returnSeen", field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedDynamicStruct, nestedStaticStruct)

	if err != nil {
		return *new(TestStruct), err
	}

	out0 := *abi.ConvertType(out[0], new(TestStruct)).(*TestStruct)

	return out0, err

}

func (_ChainReaderTester *ChainReaderTesterSession) ReturnSeen(field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedDynamicStruct MidLevelDynamicTestStruct, nestedStaticStruct MidLevelStaticTestStruct) (TestStruct, error) {
	return _ChainReaderTester.Contract.ReturnSeen(&_ChainReaderTester.CallOpts, field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedDynamicStruct, nestedStaticStruct)
}

func (_ChainReaderTester *ChainReaderTesterCallerSession) ReturnSeen(field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedDynamicStruct MidLevelDynamicTestStruct, nestedStaticStruct MidLevelStaticTestStruct) (TestStruct, error) {
	return _ChainReaderTester.Contract.ReturnSeen(&_ChainReaderTester.CallOpts, field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedDynamicStruct, nestedStaticStruct)
}

func (_ChainReaderTester *ChainReaderTesterTransactor) AddTestStruct(opts *bind.TransactOpts, field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedDynamicStruct MidLevelDynamicTestStruct, nestedStaticStruct MidLevelStaticTestStruct) (*types.Transaction, error) {
	return _ChainReaderTester.contract.Transact(opts, "addTestStruct", field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedDynamicStruct, nestedStaticStruct)
}

func (_ChainReaderTester *ChainReaderTesterSession) AddTestStruct(field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedDynamicStruct MidLevelDynamicTestStruct, nestedStaticStruct MidLevelStaticTestStruct) (*types.Transaction, error) {
	return _ChainReaderTester.Contract.AddTestStruct(&_ChainReaderTester.TransactOpts, field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedDynamicStruct, nestedStaticStruct)
}

func (_ChainReaderTester *ChainReaderTesterTransactorSession) AddTestStruct(field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedDynamicStruct MidLevelDynamicTestStruct, nestedStaticStruct MidLevelStaticTestStruct) (*types.Transaction, error) {
	return _ChainReaderTester.Contract.AddTestStruct(&_ChainReaderTester.TransactOpts, field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedDynamicStruct, nestedStaticStruct)
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

func (_ChainReaderTester *ChainReaderTesterTransactor) TriggerEvent(opts *bind.TransactOpts, field int32, oracleId uint8, nestedDynamicStruct MidLevelDynamicTestStruct, nestedStaticStruct MidLevelStaticTestStruct, oracleIds [32]uint8, account common.Address, accounts []common.Address, differentField string, bigField *big.Int) (*types.Transaction, error) {
	return _ChainReaderTester.contract.Transact(opts, "triggerEvent", field, oracleId, nestedDynamicStruct, nestedStaticStruct, oracleIds, account, accounts, differentField, bigField)
}

func (_ChainReaderTester *ChainReaderTesterSession) TriggerEvent(field int32, oracleId uint8, nestedDynamicStruct MidLevelDynamicTestStruct, nestedStaticStruct MidLevelStaticTestStruct, oracleIds [32]uint8, account common.Address, accounts []common.Address, differentField string, bigField *big.Int) (*types.Transaction, error) {
	return _ChainReaderTester.Contract.TriggerEvent(&_ChainReaderTester.TransactOpts, field, oracleId, nestedDynamicStruct, nestedStaticStruct, oracleIds, account, accounts, differentField, bigField)
}

func (_ChainReaderTester *ChainReaderTesterTransactorSession) TriggerEvent(field int32, oracleId uint8, nestedDynamicStruct MidLevelDynamicTestStruct, nestedStaticStruct MidLevelStaticTestStruct, oracleIds [32]uint8, account common.Address, accounts []common.Address, differentField string, bigField *big.Int) (*types.Transaction, error) {
	return _ChainReaderTester.Contract.TriggerEvent(&_ChainReaderTester.TransactOpts, field, oracleId, nestedDynamicStruct, nestedStaticStruct, oracleIds, account, accounts, differentField, bigField)
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

func (_ChainReaderTester *ChainReaderTesterTransactor) TriggerWithFourTopicsWithHashed(opts *bind.TransactOpts, field1 string, field2 [32]uint8, field3 [32]byte) (*types.Transaction, error) {
	return _ChainReaderTester.contract.Transact(opts, "triggerWithFourTopicsWithHashed", field1, field2, field3)
}

func (_ChainReaderTester *ChainReaderTesterSession) TriggerWithFourTopicsWithHashed(field1 string, field2 [32]uint8, field3 [32]byte) (*types.Transaction, error) {
	return _ChainReaderTester.Contract.TriggerWithFourTopicsWithHashed(&_ChainReaderTester.TransactOpts, field1, field2, field3)
}

func (_ChainReaderTester *ChainReaderTesterTransactorSession) TriggerWithFourTopicsWithHashed(field1 string, field2 [32]uint8, field3 [32]byte) (*types.Transaction, error) {
	return _ChainReaderTester.Contract.TriggerWithFourTopicsWithHashed(&_ChainReaderTester.TransactOpts, field1, field2, field3)
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
	Field               int32
	OracleId            uint8
	NestedDynamicStruct MidLevelDynamicTestStruct
	NestedStaticStruct  MidLevelStaticTestStruct
	OracleIds           [32]uint8
	Account             common.Address
	Accounts            []common.Address
	DifferentField      string
	BigField            *big.Int
	Raw                 types.Log
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

type ChainReaderTesterTriggeredWithFourTopicsWithHashedIterator struct {
	Event *ChainReaderTesterTriggeredWithFourTopicsWithHashed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ChainReaderTesterTriggeredWithFourTopicsWithHashedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChainReaderTesterTriggeredWithFourTopicsWithHashed)
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
		it.Event = new(ChainReaderTesterTriggeredWithFourTopicsWithHashed)
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

func (it *ChainReaderTesterTriggeredWithFourTopicsWithHashedIterator) Error() error {
	return it.fail
}

func (it *ChainReaderTesterTriggeredWithFourTopicsWithHashedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ChainReaderTesterTriggeredWithFourTopicsWithHashed struct {
	Field1 common.Hash
	Field2 [32]uint8
	Field3 [32]byte
	Raw    types.Log
}

func (_ChainReaderTester *ChainReaderTesterFilterer) FilterTriggeredWithFourTopicsWithHashed(opts *bind.FilterOpts, field1 []string, field2 [][32]uint8, field3 [][32]byte) (*ChainReaderTesterTriggeredWithFourTopicsWithHashedIterator, error) {

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

	logs, sub, err := _ChainReaderTester.contract.FilterLogs(opts, "TriggeredWithFourTopicsWithHashed", field1Rule, field2Rule, field3Rule)
	if err != nil {
		return nil, err
	}
	return &ChainReaderTesterTriggeredWithFourTopicsWithHashedIterator{contract: _ChainReaderTester.contract, event: "TriggeredWithFourTopicsWithHashed", logs: logs, sub: sub}, nil
}

func (_ChainReaderTester *ChainReaderTesterFilterer) WatchTriggeredWithFourTopicsWithHashed(opts *bind.WatchOpts, sink chan<- *ChainReaderTesterTriggeredWithFourTopicsWithHashed, field1 []string, field2 [][32]uint8, field3 [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _ChainReaderTester.contract.WatchLogs(opts, "TriggeredWithFourTopicsWithHashed", field1Rule, field2Rule, field3Rule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ChainReaderTesterTriggeredWithFourTopicsWithHashed)
				if err := _ChainReaderTester.contract.UnpackLog(event, "TriggeredWithFourTopicsWithHashed", log); err != nil {
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

func (_ChainReaderTester *ChainReaderTesterFilterer) ParseTriggeredWithFourTopicsWithHashed(log types.Log) (*ChainReaderTesterTriggeredWithFourTopicsWithHashed, error) {
	event := new(ChainReaderTesterTriggeredWithFourTopicsWithHashed)
	if err := _ChainReaderTester.contract.UnpackLog(event, "TriggeredWithFourTopicsWithHashed", log); err != nil {
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
	case _ChainReaderTester.abi.Events["TriggeredWithFourTopicsWithHashed"].ID:
		return _ChainReaderTester.ParseTriggeredWithFourTopicsWithHashed(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (ChainReaderTesterTriggered) Topic() common.Hash {
	return common.HexToHash("0xae927edae02672fdcce7d7e8cf34c611ed3856914a159df5f2a59307b767c25b")
}

func (ChainReaderTesterTriggeredEventWithDynamicTopic) Topic() common.Hash {
	return common.HexToHash("0x3d969732b1bbbb9f1d7eb9f3f14e4cb50a74d950b3ef916a397b85dfbab93c67")
}

func (ChainReaderTesterTriggeredWithFourTopics) Topic() common.Hash {
	return common.HexToHash("0x91c80dc390f3d041b3a04b0099b19634499541ea26972250986ee4b24a12fac5")
}

func (ChainReaderTesterTriggeredWithFourTopicsWithHashed) Topic() common.Hash {
	return common.HexToHash("0x7220e4dbe4e9d0ed5f71acd022bc89c26748ac6784f2c548bc17bb8e52af34b0")
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

	ReturnSeen(opts *bind.CallOpts, field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedDynamicStruct MidLevelDynamicTestStruct, nestedStaticStruct MidLevelStaticTestStruct) (TestStruct, error)

	AddTestStruct(opts *bind.TransactOpts, field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedDynamicStruct MidLevelDynamicTestStruct, nestedStaticStruct MidLevelStaticTestStruct) (*types.Transaction, error)

	SetAlterablePrimitiveValue(opts *bind.TransactOpts, value uint64) (*types.Transaction, error)

	TriggerEvent(opts *bind.TransactOpts, field int32, oracleId uint8, nestedDynamicStruct MidLevelDynamicTestStruct, nestedStaticStruct MidLevelStaticTestStruct, oracleIds [32]uint8, account common.Address, accounts []common.Address, differentField string, bigField *big.Int) (*types.Transaction, error)

	TriggerEventWithDynamicTopic(opts *bind.TransactOpts, field string) (*types.Transaction, error)

	TriggerWithFourTopics(opts *bind.TransactOpts, field1 int32, field2 int32, field3 int32) (*types.Transaction, error)

	TriggerWithFourTopicsWithHashed(opts *bind.TransactOpts, field1 string, field2 [32]uint8, field3 [32]byte) (*types.Transaction, error)

	FilterTriggered(opts *bind.FilterOpts, field []int32) (*ChainReaderTesterTriggeredIterator, error)

	WatchTriggered(opts *bind.WatchOpts, sink chan<- *ChainReaderTesterTriggered, field []int32) (event.Subscription, error)

	ParseTriggered(log types.Log) (*ChainReaderTesterTriggered, error)

	FilterTriggeredEventWithDynamicTopic(opts *bind.FilterOpts, fieldHash []string) (*ChainReaderTesterTriggeredEventWithDynamicTopicIterator, error)

	WatchTriggeredEventWithDynamicTopic(opts *bind.WatchOpts, sink chan<- *ChainReaderTesterTriggeredEventWithDynamicTopic, fieldHash []string) (event.Subscription, error)

	ParseTriggeredEventWithDynamicTopic(log types.Log) (*ChainReaderTesterTriggeredEventWithDynamicTopic, error)

	FilterTriggeredWithFourTopics(opts *bind.FilterOpts, field1 []int32, field2 []int32, field3 []int32) (*ChainReaderTesterTriggeredWithFourTopicsIterator, error)

	WatchTriggeredWithFourTopics(opts *bind.WatchOpts, sink chan<- *ChainReaderTesterTriggeredWithFourTopics, field1 []int32, field2 []int32, field3 []int32) (event.Subscription, error)

	ParseTriggeredWithFourTopics(log types.Log) (*ChainReaderTesterTriggeredWithFourTopics, error)

	FilterTriggeredWithFourTopicsWithHashed(opts *bind.FilterOpts, field1 []string, field2 [][32]uint8, field3 [][32]byte) (*ChainReaderTesterTriggeredWithFourTopicsWithHashedIterator, error)

	WatchTriggeredWithFourTopicsWithHashed(opts *bind.WatchOpts, sink chan<- *ChainReaderTesterTriggeredWithFourTopicsWithHashed, field1 []string, field2 [][32]uint8, field3 [][32]byte) (event.Subscription, error)

	ParseTriggeredWithFourTopicsWithHashed(log types.Log) (*ChainReaderTesterTriggeredWithFourTopicsWithHashed, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
