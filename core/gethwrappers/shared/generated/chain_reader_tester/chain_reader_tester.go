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

type AccountStruct struct {
	Account    common.Address
	AccountStr common.Address
}

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
	AccountStruct       AccountStruct
	Accounts            []common.Address
	BigField            *big.Int
	NestedDynamicStruct MidLevelDynamicTestStruct
	NestedStaticStruct  MidLevelStaticTestStruct
}

var ChainReaderTesterMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addTestStruct\",\"inputs\":[{\"name\":\"field\",\"type\":\"int32\",\"internalType\":\"int32\"},{\"name\":\"differentField\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"oracleId\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"oracleIds\",\"type\":\"uint8[32]\",\"internalType\":\"uint8[32]\"},{\"name\":\"accountStruct\",\"type\":\"tuple\",\"internalType\":\"structAccountStruct\",\"components\":[{\"name\":\"Account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"AccountStr\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"name\":\"accounts\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"bigField\",\"type\":\"int192\",\"internalType\":\"int192\"},{\"name\":\"nestedDynamicStruct\",\"type\":\"tuple\",\"internalType\":\"structMidLevelDynamicTestStruct\",\"components\":[{\"name\":\"FixedBytes\",\"type\":\"bytes2\",\"internalType\":\"bytes2\"},{\"name\":\"Inner\",\"type\":\"tuple\",\"internalType\":\"structInnerDynamicTestStruct\",\"components\":[{\"name\":\"IntVal\",\"type\":\"int64\",\"internalType\":\"int64\"},{\"name\":\"S\",\"type\":\"string\",\"internalType\":\"string\"}]}]},{\"name\":\"nestedStaticStruct\",\"type\":\"tuple\",\"internalType\":\"structMidLevelStaticTestStruct\",\"components\":[{\"name\":\"FixedBytes\",\"type\":\"bytes2\",\"internalType\":\"bytes2\"},{\"name\":\"Inner\",\"type\":\"tuple\",\"internalType\":\"structInnerStaticTestStruct\",\"components\":[{\"name\":\"IntVal\",\"type\":\"int64\",\"internalType\":\"int64\"},{\"name\":\"A\",\"type\":\"address\",\"internalType\":\"address\"}]}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAlterablePrimitiveValue\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDifferentPrimitiveValue\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getElementAtIndex\",\"inputs\":[{\"name\":\"i\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structTestStruct\",\"components\":[{\"name\":\"Field\",\"type\":\"int32\",\"internalType\":\"int32\"},{\"name\":\"DifferentField\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"OracleId\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"OracleIds\",\"type\":\"uint8[32]\",\"internalType\":\"uint8[32]\"},{\"name\":\"AccountStruct\",\"type\":\"tuple\",\"internalType\":\"structAccountStruct\",\"components\":[{\"name\":\"Account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"AccountStr\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"name\":\"Accounts\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"BigField\",\"type\":\"int192\",\"internalType\":\"int192\"},{\"name\":\"NestedDynamicStruct\",\"type\":\"tuple\",\"internalType\":\"structMidLevelDynamicTestStruct\",\"components\":[{\"name\":\"FixedBytes\",\"type\":\"bytes2\",\"internalType\":\"bytes2\"},{\"name\":\"Inner\",\"type\":\"tuple\",\"internalType\":\"structInnerDynamicTestStruct\",\"components\":[{\"name\":\"IntVal\",\"type\":\"int64\",\"internalType\":\"int64\"},{\"name\":\"S\",\"type\":\"string\",\"internalType\":\"string\"}]}]},{\"name\":\"NestedStaticStruct\",\"type\":\"tuple\",\"internalType\":\"structMidLevelStaticTestStruct\",\"components\":[{\"name\":\"FixedBytes\",\"type\":\"bytes2\",\"internalType\":\"bytes2\"},{\"name\":\"Inner\",\"type\":\"tuple\",\"internalType\":\"structInnerStaticTestStruct\",\"components\":[{\"name\":\"IntVal\",\"type\":\"int64\",\"internalType\":\"int64\"},{\"name\":\"A\",\"type\":\"address\",\"internalType\":\"address\"}]}]}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPrimitiveValue\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getSliceValue\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"returnSeen\",\"inputs\":[{\"name\":\"field\",\"type\":\"int32\",\"internalType\":\"int32\"},{\"name\":\"differentField\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"oracleId\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"oracleIds\",\"type\":\"uint8[32]\",\"internalType\":\"uint8[32]\"},{\"name\":\"accountStruct\",\"type\":\"tuple\",\"internalType\":\"structAccountStruct\",\"components\":[{\"name\":\"Account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"AccountStr\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"name\":\"accounts\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"bigField\",\"type\":\"int192\",\"internalType\":\"int192\"},{\"name\":\"nestedDynamicStruct\",\"type\":\"tuple\",\"internalType\":\"structMidLevelDynamicTestStruct\",\"components\":[{\"name\":\"FixedBytes\",\"type\":\"bytes2\",\"internalType\":\"bytes2\"},{\"name\":\"Inner\",\"type\":\"tuple\",\"internalType\":\"structInnerDynamicTestStruct\",\"components\":[{\"name\":\"IntVal\",\"type\":\"int64\",\"internalType\":\"int64\"},{\"name\":\"S\",\"type\":\"string\",\"internalType\":\"string\"}]}]},{\"name\":\"nestedStaticStruct\",\"type\":\"tuple\",\"internalType\":\"structMidLevelStaticTestStruct\",\"components\":[{\"name\":\"FixedBytes\",\"type\":\"bytes2\",\"internalType\":\"bytes2\"},{\"name\":\"Inner\",\"type\":\"tuple\",\"internalType\":\"structInnerStaticTestStruct\",\"components\":[{\"name\":\"IntVal\",\"type\":\"int64\",\"internalType\":\"int64\"},{\"name\":\"A\",\"type\":\"address\",\"internalType\":\"address\"}]}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structTestStruct\",\"components\":[{\"name\":\"Field\",\"type\":\"int32\",\"internalType\":\"int32\"},{\"name\":\"DifferentField\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"OracleId\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"OracleIds\",\"type\":\"uint8[32]\",\"internalType\":\"uint8[32]\"},{\"name\":\"AccountStruct\",\"type\":\"tuple\",\"internalType\":\"structAccountStruct\",\"components\":[{\"name\":\"Account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"AccountStr\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"name\":\"Accounts\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"BigField\",\"type\":\"int192\",\"internalType\":\"int192\"},{\"name\":\"NestedDynamicStruct\",\"type\":\"tuple\",\"internalType\":\"structMidLevelDynamicTestStruct\",\"components\":[{\"name\":\"FixedBytes\",\"type\":\"bytes2\",\"internalType\":\"bytes2\"},{\"name\":\"Inner\",\"type\":\"tuple\",\"internalType\":\"structInnerDynamicTestStruct\",\"components\":[{\"name\":\"IntVal\",\"type\":\"int64\",\"internalType\":\"int64\"},{\"name\":\"S\",\"type\":\"string\",\"internalType\":\"string\"}]}]},{\"name\":\"NestedStaticStruct\",\"type\":\"tuple\",\"internalType\":\"structMidLevelStaticTestStruct\",\"components\":[{\"name\":\"FixedBytes\",\"type\":\"bytes2\",\"internalType\":\"bytes2\"},{\"name\":\"Inner\",\"type\":\"tuple\",\"internalType\":\"structInnerStaticTestStruct\",\"components\":[{\"name\":\"IntVal\",\"type\":\"int64\",\"internalType\":\"int64\"},{\"name\":\"A\",\"type\":\"address\",\"internalType\":\"address\"}]}]}]}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"setAlterablePrimitiveValue\",\"inputs\":[{\"name\":\"value\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"triggerEvent\",\"inputs\":[{\"name\":\"field\",\"type\":\"int32\",\"internalType\":\"int32\"},{\"name\":\"oracleId\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"nestedDynamicStruct\",\"type\":\"tuple\",\"internalType\":\"structMidLevelDynamicTestStruct\",\"components\":[{\"name\":\"FixedBytes\",\"type\":\"bytes2\",\"internalType\":\"bytes2\"},{\"name\":\"Inner\",\"type\":\"tuple\",\"internalType\":\"structInnerDynamicTestStruct\",\"components\":[{\"name\":\"IntVal\",\"type\":\"int64\",\"internalType\":\"int64\"},{\"name\":\"S\",\"type\":\"string\",\"internalType\":\"string\"}]}]},{\"name\":\"nestedStaticStruct\",\"type\":\"tuple\",\"internalType\":\"structMidLevelStaticTestStruct\",\"components\":[{\"name\":\"FixedBytes\",\"type\":\"bytes2\",\"internalType\":\"bytes2\"},{\"name\":\"Inner\",\"type\":\"tuple\",\"internalType\":\"structInnerStaticTestStruct\",\"components\":[{\"name\":\"IntVal\",\"type\":\"int64\",\"internalType\":\"int64\"},{\"name\":\"A\",\"type\":\"address\",\"internalType\":\"address\"}]}]},{\"name\":\"oracleIds\",\"type\":\"uint8[32]\",\"internalType\":\"uint8[32]\"},{\"name\":\"accountStruct\",\"type\":\"tuple\",\"internalType\":\"structAccountStruct\",\"components\":[{\"name\":\"Account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"AccountStr\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"name\":\"accounts\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"differentField\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"bigField\",\"type\":\"int192\",\"internalType\":\"int192\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"triggerEventWithDynamicTopic\",\"inputs\":[{\"name\":\"field\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"triggerStaticBytes\",\"inputs\":[{\"name\":\"val1\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"val2\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"val3\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"val4\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"val5\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"val6\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"val7\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"raw\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"triggerWithFourTopics\",\"inputs\":[{\"name\":\"field1\",\"type\":\"int32\",\"internalType\":\"int32\"},{\"name\":\"field2\",\"type\":\"int32\",\"internalType\":\"int32\"},{\"name\":\"field3\",\"type\":\"int32\",\"internalType\":\"int32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"triggerWithFourTopicsWithHashed\",\"inputs\":[{\"name\":\"field1\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"field2\",\"type\":\"uint8[32]\",\"internalType\":\"uint8[32]\"},{\"name\":\"field3\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"StaticBytes\",\"inputs\":[{\"name\":\"message\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Triggered\",\"inputs\":[{\"name\":\"field\",\"type\":\"int32\",\"indexed\":true,\"internalType\":\"int32\"},{\"name\":\"oracleId\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"},{\"name\":\"nestedDynamicStruct\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structMidLevelDynamicTestStruct\",\"components\":[{\"name\":\"FixedBytes\",\"type\":\"bytes2\",\"internalType\":\"bytes2\"},{\"name\":\"Inner\",\"type\":\"tuple\",\"internalType\":\"structInnerDynamicTestStruct\",\"components\":[{\"name\":\"IntVal\",\"type\":\"int64\",\"internalType\":\"int64\"},{\"name\":\"S\",\"type\":\"string\",\"internalType\":\"string\"}]}]},{\"name\":\"nestedStaticStruct\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structMidLevelStaticTestStruct\",\"components\":[{\"name\":\"FixedBytes\",\"type\":\"bytes2\",\"internalType\":\"bytes2\"},{\"name\":\"Inner\",\"type\":\"tuple\",\"internalType\":\"structInnerStaticTestStruct\",\"components\":[{\"name\":\"IntVal\",\"type\":\"int64\",\"internalType\":\"int64\"},{\"name\":\"A\",\"type\":\"address\",\"internalType\":\"address\"}]}]},{\"name\":\"oracleIds\",\"type\":\"uint8[32]\",\"indexed\":false,\"internalType\":\"uint8[32]\"},{\"name\":\"accountStruct\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structAccountStruct\",\"components\":[{\"name\":\"Account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"AccountStr\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"name\":\"Accounts\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"differentField\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"bigField\",\"type\":\"int192\",\"indexed\":false,\"internalType\":\"int192\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TriggeredEventWithDynamicTopic\",\"inputs\":[{\"name\":\"fieldHash\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"field\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TriggeredWithFourTopics\",\"inputs\":[{\"name\":\"field1\",\"type\":\"int32\",\"indexed\":true,\"internalType\":\"int32\"},{\"name\":\"field2\",\"type\":\"int32\",\"indexed\":true,\"internalType\":\"int32\"},{\"name\":\"field3\",\"type\":\"int32\",\"indexed\":true,\"internalType\":\"int32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TriggeredWithFourTopicsWithHashed\",\"inputs\":[{\"name\":\"field1\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"field2\",\"type\":\"uint8[32]\",\"indexed\":true,\"internalType\":\"uint8[32]\"},{\"name\":\"field3\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false}]",
	Bin: "0x608060405234801561001057600080fd5b50600180548082018255600082905260048082047fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6908101805460086003958616810261010090810a8088026001600160401b0391820219909416939093179093558654808801909755848704909301805496909516909202900a91820291021990921691909117905561209a806100a96000396000f3fe608060405234801561001057600080fd5b50600436106100df5760003560e01c80636cf01bfc1161008c578063ab5e0b3811610066578063ab5e0b38146101df578063dbfd7332146101fc578063ef4e1ced1461020f578063fbe9fbf61461021657600080fd5b80636cf01bfc146101a657806380d77303146101b9578063a90e1998146101cc57600080fd5b8063660b5962116100bd578063660b596214610135578063679004a4146101485780636c9a43b61461015d57600080fd5b80632c45576f146100e45780633272b66c1461010d57806351f3f54d14610122575b600080fd5b6100f76100f2366004610f24565b610228565b604051610104919061108d565b60405180910390f35b61012061011b366004611240565b61059b565b005b61012061013036600461139f565b6105f0565b61012061014336600461150c565b610664565b610150610a33565b6040516101049190611612565b61012061016b366004611660565b600280547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001667ffffffffffffffff92909216919091179055565b6100f76101b436600461150c565b610abf565b6101206101c7366004611682565b610bd6565b6101206101da366004611795565b610c30565b6107c65b60405167ffffffffffffffff9091168152602001610104565b61012061020a36600461184a565b610c8a565b60036101e3565b60025467ffffffffffffffff166101e3565b610230610cc7565b600061023d60018461188d565b8154811061024d5761024d6118c7565b6000918252602091829020604080516101208101909152600d90920201805460030b82526001810180549293919291840191610288906118f6565b80601f01602080910402602001604051908101604052809291908181526020018280546102b4906118f6565b80156103015780601f106102d657610100808354040283529160200191610301565b820191906000526020600020905b8154815290600101906020018083116102e457829003601f168201915b5050509183525050600282015460ff166020808301919091526040805161040081018083529190930192916003850191826000855b825461010083900a900460ff168152602060019283018181049485019490930390920291018084116103365790505050509183525050604080518082018252600484015473ffffffffffffffffffffffffffffffffffffffff90811682526005850154166020828101919091528084019190915260068401805483518184028101840185528181529390940193909183018282801561040b57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116103e0575b505050918352505060078281015460170b60208084019190915260408051808201825260088601805460f01b7fffff0000000000000000000000000000000000000000000000000000000000001682528251808401845260098801805490960b8152600a880180549490970196929591948681019491939084019161048f906118f6565b80601f01602080910402602001604051908101604052809291908181526020018280546104bb906118f6565b80156105085780601f106104dd57610100808354040283529160200191610508565b820191906000526020600020905b8154815290600101906020018083116104eb57829003601f168201915b505050919092525050509052508152604080518082018252600b84015460f01b7fffff0000000000000000000000000000000000000000000000000000000000001681528151808301909252600c90930154600781900b825268010000000000000000900473ffffffffffffffffffffffffffffffffffffffff1660208083019190915280840191909152015292915050565b81816040516105ab929190611943565b60405180910390207f3d969732b1bbbb9f1d7eb9f3f14e4cb50a74d950b3ef916a397b85dfbab93c6783836040516105e492919061199c565b60405180910390a25050565b600088888888888888886040516020016106119897969594939291906119b0565b60405160208183030381529060405290507f1e40927ec0bdc7319f09a53452590433ec395dec3b70b982eba779c740685bfe816040516106519190611a56565b60405180910390a1505050505050505050565b60006040518061012001604052808d60030b81526020018c8c8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050509082525060ff8b166020808301919091526040805161040081810183529190930192918c918390839080828437600092019190915250505081526020016106fc368a90038a018a611a8d565b8152602001878780806020026020016040519081016040528093929190818152602001838360200280828437600092019190915250505090825250601786900b602082015260400161074d85611b0c565b815260200161076136859003850185611baa565b905281546001808201845560009384526020938490208351600d9093020180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff9093169290921782559282015191929091908201906107c79082611c89565b5060408201516002820180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff90921691909117905560608201516108159060038301906020610d90565b506080820151805160048301805473ffffffffffffffffffffffffffffffffffffffff9283167fffffffffffffffffffffffff00000000000000000000000000000000000000009182161790915560209283015160058501805491909316911617905560a0830151805161088f9260068501920190610e23565b5060c08201516007820180547fffffffffffffffff0000000000000000000000000000000000000000000000001677ffffffffffffffffffffffffffffffffffffffffffffffff90921691909117905560e082015180516008830180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001660f09290921c91909117815560208083015180516009860180547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001667ffffffffffffffff90921691909117815591810151909190600a8601906109729082611c89565b5050505061010092909201518051600b8301805460f09290921c7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00009092169190911790556020908101518051600c9093018054919092015173ffffffffffffffffffffffffffffffffffffffff1668010000000000000000027fffffffff0000000000000000000000000000000000000000000000000000000090911667ffffffffffffffff90931692909217919091179055505050505050505050505050565b60606001805480602002602001604051908101604052809291908181526020018280548015610ab557602002820191906000526020600020906000905b82829054906101000a900467ffffffffffffffff1667ffffffffffffffff1681526020019060080190602082600701049283019260010382029150808411610a705790505b5050505050905090565b610ac7610cc7565b6040518061012001604052808d60030b81526020018c8c8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050509082525060ff8b166020808301919091526040805161040081810183529190930192918c91839083908082843760009201919091525050508152602001610b5d368a90038a018a611a8d565b8152602001878780806020026020016040519081016040528093929190818152602001838360200280828437600092019190915250505090825250601786900b6020820152604001610bae85611b0c565b8152602001610bc236859003850185611baa565b905290505b9b9a5050505050505050505050565b8a60030b7f6a1c0d59d12839c1d2bd0a8cd5e73032dc60569f7f448b99b572b61c442b8add8b8b8b8b8b8b8b8b8b8b604051610c1b9a99989796959493929190611ec4565b60405180910390a25050505050505050505050565b8082604051610c3f9190612035565b604051809103902084604051610c559190612071565b604051908190038120907f7220e4dbe4e9d0ed5f71acd022bc89c26748ac6784f2c548bc17bb8e52af34b090600090a4505050565b8060030b8260030b8460030b7f91c80dc390f3d041b3a04b0099b19634499541ea26972250986ee4b24a12fac560405160405180910390a4505050565b6040805161012081018252600080825260606020830181905292820152908101610cef610e9d565b8152602001610d3d6040518060400160405280600073ffffffffffffffffffffffffffffffffffffffff168152602001600073ffffffffffffffffffffffffffffffffffffffff1681525090565b81526060602082018190526000604083015201610d58610ebc565b8152602001610d8b6040805180820182526000808252825180840190935280835260208381019190915290919082015290565b905290565b600183019183908215610e135791602002820160005b83821115610de457835183826101000a81548160ff021916908360ff1602179055509260200192600101602081600001049283019260010302610da6565b8015610e115782816101000a81549060ff0219169055600101602081600001049283019260010302610de4565b505b50610e1f929150610f0f565b5090565b828054828255906000526020600020908101928215610e13579160200282015b82811115610e1357825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190610e43565b6040518061040001604052806020906020820280368337509192915050565b604051806040016040528060007dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168152602001610d8b6040518060400160405280600060070b8152602001606081525090565b5b80821115610e1f5760008155600101610f10565b600060208284031215610f3657600080fd5b5035919050565b60005b83811015610f58578181015183820152602001610f40565b50506000910152565b60008151808452610f79816020860160208601610f3d565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b8060005b6020808210610fbe5750610fd5565b825160ff1685529384019390910190600101610faf565b50505050565b600081518084526020808501945080840160005b8381101561102157815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101610fef565b509495945050505050565b7fffff00000000000000000000000000000000000000000000000000000000000081511682526000602082015160406020850152805160070b604085015260208101519050604060608501526110856080850182610f61565b949350505050565b602081526110a160208201835160030b9052565b6000602083015161056060408401526110be610580840182610f61565b905060408401516110d4606085018260ff169052565b5060608401516110e76080850182610fab565b506080840151805173ffffffffffffffffffffffffffffffffffffffff9081166104808601526020820151166104a08501525060a08401517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe080858403016104c08601526111558383610fdb565b925060c0860151915061116e6104e086018360170b9052565b60e0860151915080858403016105008601525061118b828261102c565b61010086015180517fffff00000000000000000000000000000000000000000000000000000000000016610520870152602080820151805160070b610540890152015173ffffffffffffffffffffffffffffffffffffffff166105608701529092509050509392505050565b60008083601f84011261120957600080fd5b50813567ffffffffffffffff81111561122157600080fd5b60208301915083602082850101111561123957600080fd5b9250929050565b6000806020838503121561125357600080fd5b823567ffffffffffffffff81111561126a57600080fd5b611276858286016111f7565b90969095509350505050565b803563ffffffff8116811461129657600080fd5b919050565b803567ffffffffffffffff8116811461129657600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040805190810167ffffffffffffffff81118282101715611305576113056112b3565b60405290565b600067ffffffffffffffff80841115611326576113266112b3565b604051601f85017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f0116810190828211818310171561136c5761136c6112b3565b8160405280935085815286868601111561138557600080fd5b858560208301376000602087830101525050509392505050565b600080600080600080600080610100898b0312156113bc57600080fd5b6113c589611282565b97506113d360208a01611282565b96506113e160408a01611282565b95506113ef60608a0161129b565b94506080890135935060a0890135925060c0890135915060e089013567ffffffffffffffff81111561142057600080fd5b8901601f81018b1361143157600080fd5b6114408b82356020840161130b565b9150509295985092959890939650565b8035600381900b811461129657600080fd5b803560ff8116811461129657600080fd5b80610400810183101561148557600080fd5b92915050565b60006040828403121561149d57600080fd5b50919050565b60008083601f8401126114b557600080fd5b50813567ffffffffffffffff8111156114cd57600080fd5b6020830191508360208260051b850101111561123957600080fd5b8035601781900b811461129657600080fd5b60006060828403121561149d57600080fd5b60008060008060008060008060008060006105608c8e03121561152e57600080fd5b6115378c611450565b9a5067ffffffffffffffff8060208e0135111561155357600080fd5b6115638e60208f01358f016111f7565b909b50995061157460408e01611462565b98506115838e60608f01611473565b97506115938e6104608f0161148b565b9650806104a08e013511156115a757600080fd5b6115b88e6104a08f01358f016114a3565b90965094506115ca6104c08e016114e8565b9350806104e08e013511156115de57600080fd5b506115f08d6104e08e01358e0161148b565b91506116008d6105008e016114fa565b90509295989b509295989b9093969950565b6020808252825182820181905260009190848201906040850190845b8181101561165457835167ffffffffffffffff168352928401929184019160010161162e565b50909695505050505050565b60006020828403121561167257600080fd5b61167b8261129b565b9392505050565b60008060008060008060008060008060006105608c8e0312156116a457600080fd5b6116ad8c611450565b9a506116bb60208d01611462565b995067ffffffffffffffff8060408e013511156116d757600080fd5b6116e78e60408f01358f0161148b565b99506116f68e60608f016114fa565b98506117058e60c08f01611473565b97506117158e6104c08f0161148b565b9650806105008e0135111561172957600080fd5b61173a8e6105008f01358f016114a3565b90965094506105208d013581101561175157600080fd5b506117638d6105208e01358e016111f7565b90935091506116006105408d016114e8565b600082601f83011261178657600080fd5b61167b8383356020850161130b565b600080600061044084860312156117ab57600080fd5b833567ffffffffffffffff808211156117c357600080fd5b6117cf87838801611775565b94506020915086603f8701126117e457600080fd5b60405161040081018181108382111715611800576118006112b3565b60405290508061042087018881111561181857600080fd5b8388015b8181101561183a5761182d81611462565b845292840192840161181c565b5095989097509435955050505050565b60008060006060848603121561185f57600080fd5b61186884611450565b925061187660208501611450565b915061188460408501611450565b90509250925092565b81810381811115611485577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600181811c9082168061190a57607f821691505b60208210810361149d577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b8183823760009101908152919050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b602081526000611085602083018486611953565b60007fffffffff00000000000000000000000000000000000000000000000000000000808b60e01b168352808a60e01b166004840152808960e01b166008840152507fffffffffffffffff0000000000000000000000000000000000000000000000008760c01b16600c8301528560148301528460348301528360548301528251611a42816074850160208701610f3d565b919091016074019998505050505050505050565b60208152600061167b6020830184610f61565b803573ffffffffffffffffffffffffffffffffffffffff8116811461129657600080fd5b600060408284031215611a9f57600080fd5b611aa76112e2565b611ab083611a69565b8152611abe60208401611a69565b60208201529392505050565b80357fffff0000000000000000000000000000000000000000000000000000000000008116811461129657600080fd5b8035600781900b811461129657600080fd5b600060408236031215611b1e57600080fd5b611b266112e2565b611b2f83611aca565b8152602083013567ffffffffffffffff80821115611b4c57600080fd5b818501915060408236031215611b6157600080fd5b611b696112e2565b611b7283611afa565b8152602083013582811115611b8657600080fd5b611b9236828601611775565b60208301525080602085015250505080915050919050565b60008183036060811215611bbd57600080fd5b611bc56112e2565b611bce84611aca565b815260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe083011215611c0057600080fd5b611c086112e2565b9150611c1660208501611afa565b8252611c2460408501611a69565b6020830152816020820152809250505092915050565b601f821115611c8457600081815260208120601f850160051c81016020861015611c615750805b601f850160051c820191505b81811015611c8057828155600101611c6d565b5050505b505050565b815167ffffffffffffffff811115611ca357611ca36112b3565b611cb781611cb184546118f6565b84611c3a565b602080601f831160018114611d0a5760008415611cd45750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555611c80565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015611d5757888601518255948401946001909101908401611d38565b5085821015611d9357878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b7fffff000000000000000000000000000000000000000000000000000000000000611dcd82611aca565b168252611ddc60208201611afa565b60070b602083015273ffffffffffffffffffffffffffffffffffffffff611e0560408301611a69565b1660408301525050565b8060005b6020808210611e225750610fd5565b60ff611e2d84611462565b168552938401939190910190600101611e13565b73ffffffffffffffffffffffffffffffffffffffff80611e6083611a69565b16835280611e7060208401611a69565b166020840152505050565b8183526000602080850194508260005b858110156110215773ffffffffffffffffffffffffffffffffffffffff611eb183611a69565b1687529582019590820190600101611e8b565b600061054060ff8d1683528060208401527fffff000000000000000000000000000000000000000000000000000000000000611eff8d611aca565b16818401525060208b01357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc18c3603018112611f3a57600080fd5b60406105608401528b01611f4d81611afa565b60070b61058084015260208101357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1823603018112611f8b57600080fd5b0160208101903567ffffffffffffffff811115611fa757600080fd5b803603821315611fb657600080fd5b60406105a0850152611fcd6105c085018284611953565b915050611fdd604084018c611da3565b611fea60a084018b611e0f565b611ff86104a084018a611e41565b8281036104e084015261200c81888a611e7b565b9050828103610500840152612022818688611953565b915050610bc761052083018460170b9052565b60008183825b602080821061204a5750612061565b825160ff168452928301929091019060010161203b565b5050506104008201905092915050565b60008251612083818460208701610f3d565b919091019291505056fea164736f6c6343000813000a",
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

func (_ChainReaderTester *ChainReaderTesterCaller) ReturnSeen(opts *bind.CallOpts, field int32, differentField string, oracleId uint8, oracleIds [32]uint8, accountStruct AccountStruct, accounts []common.Address, bigField *big.Int, nestedDynamicStruct MidLevelDynamicTestStruct, nestedStaticStruct MidLevelStaticTestStruct) (TestStruct, error) {
	var out []interface{}
	err := _ChainReaderTester.contract.Call(opts, &out, "returnSeen", field, differentField, oracleId, oracleIds, accountStruct, accounts, bigField, nestedDynamicStruct, nestedStaticStruct)

	if err != nil {
		return *new(TestStruct), err
	}

	out0 := *abi.ConvertType(out[0], new(TestStruct)).(*TestStruct)

	return out0, err

}

func (_ChainReaderTester *ChainReaderTesterSession) ReturnSeen(field int32, differentField string, oracleId uint8, oracleIds [32]uint8, accountStruct AccountStruct, accounts []common.Address, bigField *big.Int, nestedDynamicStruct MidLevelDynamicTestStruct, nestedStaticStruct MidLevelStaticTestStruct) (TestStruct, error) {
	return _ChainReaderTester.Contract.ReturnSeen(&_ChainReaderTester.CallOpts, field, differentField, oracleId, oracleIds, accountStruct, accounts, bigField, nestedDynamicStruct, nestedStaticStruct)
}

func (_ChainReaderTester *ChainReaderTesterCallerSession) ReturnSeen(field int32, differentField string, oracleId uint8, oracleIds [32]uint8, accountStruct AccountStruct, accounts []common.Address, bigField *big.Int, nestedDynamicStruct MidLevelDynamicTestStruct, nestedStaticStruct MidLevelStaticTestStruct) (TestStruct, error) {
	return _ChainReaderTester.Contract.ReturnSeen(&_ChainReaderTester.CallOpts, field, differentField, oracleId, oracleIds, accountStruct, accounts, bigField, nestedDynamicStruct, nestedStaticStruct)
}

func (_ChainReaderTester *ChainReaderTesterTransactor) AddTestStruct(opts *bind.TransactOpts, field int32, differentField string, oracleId uint8, oracleIds [32]uint8, accountStruct AccountStruct, accounts []common.Address, bigField *big.Int, nestedDynamicStruct MidLevelDynamicTestStruct, nestedStaticStruct MidLevelStaticTestStruct) (*types.Transaction, error) {
	return _ChainReaderTester.contract.Transact(opts, "addTestStruct", field, differentField, oracleId, oracleIds, accountStruct, accounts, bigField, nestedDynamicStruct, nestedStaticStruct)
}

func (_ChainReaderTester *ChainReaderTesterSession) AddTestStruct(field int32, differentField string, oracleId uint8, oracleIds [32]uint8, accountStruct AccountStruct, accounts []common.Address, bigField *big.Int, nestedDynamicStruct MidLevelDynamicTestStruct, nestedStaticStruct MidLevelStaticTestStruct) (*types.Transaction, error) {
	return _ChainReaderTester.Contract.AddTestStruct(&_ChainReaderTester.TransactOpts, field, differentField, oracleId, oracleIds, accountStruct, accounts, bigField, nestedDynamicStruct, nestedStaticStruct)
}

func (_ChainReaderTester *ChainReaderTesterTransactorSession) AddTestStruct(field int32, differentField string, oracleId uint8, oracleIds [32]uint8, accountStruct AccountStruct, accounts []common.Address, bigField *big.Int, nestedDynamicStruct MidLevelDynamicTestStruct, nestedStaticStruct MidLevelStaticTestStruct) (*types.Transaction, error) {
	return _ChainReaderTester.Contract.AddTestStruct(&_ChainReaderTester.TransactOpts, field, differentField, oracleId, oracleIds, accountStruct, accounts, bigField, nestedDynamicStruct, nestedStaticStruct)
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

func (_ChainReaderTester *ChainReaderTesterTransactor) TriggerEvent(opts *bind.TransactOpts, field int32, oracleId uint8, nestedDynamicStruct MidLevelDynamicTestStruct, nestedStaticStruct MidLevelStaticTestStruct, oracleIds [32]uint8, accountStruct AccountStruct, accounts []common.Address, differentField string, bigField *big.Int) (*types.Transaction, error) {
	return _ChainReaderTester.contract.Transact(opts, "triggerEvent", field, oracleId, nestedDynamicStruct, nestedStaticStruct, oracleIds, accountStruct, accounts, differentField, bigField)
}

func (_ChainReaderTester *ChainReaderTesterSession) TriggerEvent(field int32, oracleId uint8, nestedDynamicStruct MidLevelDynamicTestStruct, nestedStaticStruct MidLevelStaticTestStruct, oracleIds [32]uint8, accountStruct AccountStruct, accounts []common.Address, differentField string, bigField *big.Int) (*types.Transaction, error) {
	return _ChainReaderTester.Contract.TriggerEvent(&_ChainReaderTester.TransactOpts, field, oracleId, nestedDynamicStruct, nestedStaticStruct, oracleIds, accountStruct, accounts, differentField, bigField)
}

func (_ChainReaderTester *ChainReaderTesterTransactorSession) TriggerEvent(field int32, oracleId uint8, nestedDynamicStruct MidLevelDynamicTestStruct, nestedStaticStruct MidLevelStaticTestStruct, oracleIds [32]uint8, accountStruct AccountStruct, accounts []common.Address, differentField string, bigField *big.Int) (*types.Transaction, error) {
	return _ChainReaderTester.Contract.TriggerEvent(&_ChainReaderTester.TransactOpts, field, oracleId, nestedDynamicStruct, nestedStaticStruct, oracleIds, accountStruct, accounts, differentField, bigField)
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

func (_ChainReaderTester *ChainReaderTesterTransactor) TriggerStaticBytes(opts *bind.TransactOpts, val1 uint32, val2 uint32, val3 uint32, val4 uint64, val5 [32]byte, val6 [32]byte, val7 [32]byte, raw []byte) (*types.Transaction, error) {
	return _ChainReaderTester.contract.Transact(opts, "triggerStaticBytes", val1, val2, val3, val4, val5, val6, val7, raw)
}

func (_ChainReaderTester *ChainReaderTesterSession) TriggerStaticBytes(val1 uint32, val2 uint32, val3 uint32, val4 uint64, val5 [32]byte, val6 [32]byte, val7 [32]byte, raw []byte) (*types.Transaction, error) {
	return _ChainReaderTester.Contract.TriggerStaticBytes(&_ChainReaderTester.TransactOpts, val1, val2, val3, val4, val5, val6, val7, raw)
}

func (_ChainReaderTester *ChainReaderTesterTransactorSession) TriggerStaticBytes(val1 uint32, val2 uint32, val3 uint32, val4 uint64, val5 [32]byte, val6 [32]byte, val7 [32]byte, raw []byte) (*types.Transaction, error) {
	return _ChainReaderTester.Contract.TriggerStaticBytes(&_ChainReaderTester.TransactOpts, val1, val2, val3, val4, val5, val6, val7, raw)
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

type ChainReaderTesterStaticBytesIterator struct {
	Event *ChainReaderTesterStaticBytes

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ChainReaderTesterStaticBytesIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChainReaderTesterStaticBytes)
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
		it.Event = new(ChainReaderTesterStaticBytes)
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

func (it *ChainReaderTesterStaticBytesIterator) Error() error {
	return it.fail
}

func (it *ChainReaderTesterStaticBytesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ChainReaderTesterStaticBytes struct {
	Message []byte
	Raw     types.Log
}

func (_ChainReaderTester *ChainReaderTesterFilterer) FilterStaticBytes(opts *bind.FilterOpts) (*ChainReaderTesterStaticBytesIterator, error) {

	logs, sub, err := _ChainReaderTester.contract.FilterLogs(opts, "StaticBytes")
	if err != nil {
		return nil, err
	}
	return &ChainReaderTesterStaticBytesIterator{contract: _ChainReaderTester.contract, event: "StaticBytes", logs: logs, sub: sub}, nil
}

func (_ChainReaderTester *ChainReaderTesterFilterer) WatchStaticBytes(opts *bind.WatchOpts, sink chan<- *ChainReaderTesterStaticBytes) (event.Subscription, error) {

	logs, sub, err := _ChainReaderTester.contract.WatchLogs(opts, "StaticBytes")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ChainReaderTesterStaticBytes)
				if err := _ChainReaderTester.contract.UnpackLog(event, "StaticBytes", log); err != nil {
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

func (_ChainReaderTester *ChainReaderTesterFilterer) ParseStaticBytes(log types.Log) (*ChainReaderTesterStaticBytes, error) {
	event := new(ChainReaderTesterStaticBytes)
	if err := _ChainReaderTester.contract.UnpackLog(event, "StaticBytes", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
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
	AccountStruct       AccountStruct
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
	case _ChainReaderTester.abi.Events["StaticBytes"].ID:
		return _ChainReaderTester.ParseStaticBytes(log)
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

func (ChainReaderTesterStaticBytes) Topic() common.Hash {
	return common.HexToHash("0x1e40927ec0bdc7319f09a53452590433ec395dec3b70b982eba779c740685bfe")
}

func (ChainReaderTesterTriggered) Topic() common.Hash {
	return common.HexToHash("0x6a1c0d59d12839c1d2bd0a8cd5e73032dc60569f7f448b99b572b61c442b8add")
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

	ReturnSeen(opts *bind.CallOpts, field int32, differentField string, oracleId uint8, oracleIds [32]uint8, accountStruct AccountStruct, accounts []common.Address, bigField *big.Int, nestedDynamicStruct MidLevelDynamicTestStruct, nestedStaticStruct MidLevelStaticTestStruct) (TestStruct, error)

	AddTestStruct(opts *bind.TransactOpts, field int32, differentField string, oracleId uint8, oracleIds [32]uint8, accountStruct AccountStruct, accounts []common.Address, bigField *big.Int, nestedDynamicStruct MidLevelDynamicTestStruct, nestedStaticStruct MidLevelStaticTestStruct) (*types.Transaction, error)

	SetAlterablePrimitiveValue(opts *bind.TransactOpts, value uint64) (*types.Transaction, error)

	TriggerEvent(opts *bind.TransactOpts, field int32, oracleId uint8, nestedDynamicStruct MidLevelDynamicTestStruct, nestedStaticStruct MidLevelStaticTestStruct, oracleIds [32]uint8, accountStruct AccountStruct, accounts []common.Address, differentField string, bigField *big.Int) (*types.Transaction, error)

	TriggerEventWithDynamicTopic(opts *bind.TransactOpts, field string) (*types.Transaction, error)

	TriggerStaticBytes(opts *bind.TransactOpts, val1 uint32, val2 uint32, val3 uint32, val4 uint64, val5 [32]byte, val6 [32]byte, val7 [32]byte, raw []byte) (*types.Transaction, error)

	TriggerWithFourTopics(opts *bind.TransactOpts, field1 int32, field2 int32, field3 int32) (*types.Transaction, error)

	TriggerWithFourTopicsWithHashed(opts *bind.TransactOpts, field1 string, field2 [32]uint8, field3 [32]byte) (*types.Transaction, error)

	FilterStaticBytes(opts *bind.FilterOpts) (*ChainReaderTesterStaticBytesIterator, error)

	WatchStaticBytes(opts *bind.WatchOpts, sink chan<- *ChainReaderTesterStaticBytes) (event.Subscription, error)

	ParseStaticBytes(log types.Log) (*ChainReaderTesterStaticBytes, error)

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
