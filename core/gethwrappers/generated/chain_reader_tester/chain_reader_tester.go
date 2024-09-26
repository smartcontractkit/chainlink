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
	AccountStr          common.Address
	Accounts            []common.Address
	BigField            *big.Int
	NestedDynamicStruct MidLevelDynamicTestStruct
	NestedStaticStruct  MidLevelStaticTestStruct
}

var ChainReaderTesterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int32\",\"name\":\"field\",\"type\":\"int32\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"oracleId\",\"type\":\"uint8\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerDynamicTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"indexed\":false,\"internalType\":\"structMidLevelDynamicTestStruct\",\"name\":\"nestedDynamicStruct\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"address\",\"name\":\"A\",\"type\":\"address\"}],\"internalType\":\"structInnerStaticTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"indexed\":false,\"internalType\":\"structMidLevelStaticTestStruct\",\"name\":\"nestedStaticStruct\",\"type\":\"tuple\"},{\"indexed\":false,\"internalType\":\"uint8[32]\",\"name\":\"oracleIds\",\"type\":\"uint8[32]\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"Account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"AccountStr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"Accounts\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"differentField\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"int192\",\"name\":\"bigField\",\"type\":\"int192\"}],\"name\":\"Triggered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"fieldHash\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"field\",\"type\":\"string\"}],\"name\":\"TriggeredEventWithDynamicTopic\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int32\",\"name\":\"field1\",\"type\":\"int32\"},{\"indexed\":true,\"internalType\":\"int32\",\"name\":\"field2\",\"type\":\"int32\"},{\"indexed\":true,\"internalType\":\"int32\",\"name\":\"field3\",\"type\":\"int32\"}],\"name\":\"TriggeredWithFourTopics\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"field1\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"uint8[32]\",\"name\":\"field2\",\"type\":\"uint8[32]\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"field3\",\"type\":\"bytes32\"}],\"name\":\"TriggeredWithFourTopicsWithHashed\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"int32\",\"name\":\"field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"differentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"oracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"oracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"accountStr\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"int192\",\"name\":\"bigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerDynamicTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelDynamicTestStruct\",\"name\":\"nestedDynamicStruct\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"address\",\"name\":\"A\",\"type\":\"address\"}],\"internalType\":\"structInnerStaticTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelStaticTestStruct\",\"name\":\"nestedStaticStruct\",\"type\":\"tuple\"}],\"name\":\"addTestStruct\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAlterablePrimitiveValue\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDifferentPrimitiveValue\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"i\",\"type\":\"uint256\"}],\"name\":\"getElementAtIndex\",\"outputs\":[{\"components\":[{\"internalType\":\"int32\",\"name\":\"Field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"DifferentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"OracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"OracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"address\",\"name\":\"Account\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"AccountStr\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"Accounts\",\"type\":\"address[]\"},{\"internalType\":\"int192\",\"name\":\"BigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerDynamicTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelDynamicTestStruct\",\"name\":\"NestedDynamicStruct\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"address\",\"name\":\"A\",\"type\":\"address\"}],\"internalType\":\"structInnerStaticTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelStaticTestStruct\",\"name\":\"NestedStaticStruct\",\"type\":\"tuple\"}],\"internalType\":\"structTestStruct\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPrimitiveValue\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSliceValue\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int32\",\"name\":\"field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"differentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"oracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"oracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"accountStr\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"int192\",\"name\":\"bigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerDynamicTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelDynamicTestStruct\",\"name\":\"nestedDynamicStruct\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"address\",\"name\":\"A\",\"type\":\"address\"}],\"internalType\":\"structInnerStaticTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelStaticTestStruct\",\"name\":\"nestedStaticStruct\",\"type\":\"tuple\"}],\"name\":\"returnSeen\",\"outputs\":[{\"components\":[{\"internalType\":\"int32\",\"name\":\"Field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"DifferentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"OracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"OracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"address\",\"name\":\"Account\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"AccountStr\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"Accounts\",\"type\":\"address[]\"},{\"internalType\":\"int192\",\"name\":\"BigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerDynamicTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelDynamicTestStruct\",\"name\":\"NestedDynamicStruct\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"address\",\"name\":\"A\",\"type\":\"address\"}],\"internalType\":\"structInnerStaticTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelStaticTestStruct\",\"name\":\"NestedStaticStruct\",\"type\":\"tuple\"}],\"internalType\":\"structTestStruct\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"value\",\"type\":\"uint64\"}],\"name\":\"setAlterablePrimitiveValue\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int32\",\"name\":\"field\",\"type\":\"int32\"},{\"internalType\":\"uint8\",\"name\":\"oracleId\",\"type\":\"uint8\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerDynamicTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelDynamicTestStruct\",\"name\":\"nestedDynamicStruct\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"address\",\"name\":\"A\",\"type\":\"address\"}],\"internalType\":\"structInnerStaticTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelStaticTestStruct\",\"name\":\"nestedStaticStruct\",\"type\":\"tuple\"},{\"internalType\":\"uint8[32]\",\"name\":\"oracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"accountStr\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"string\",\"name\":\"differentField\",\"type\":\"string\"},{\"internalType\":\"int192\",\"name\":\"bigField\",\"type\":\"int192\"}],\"name\":\"triggerEvent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"field\",\"type\":\"string\"}],\"name\":\"triggerEventWithDynamicTopic\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int32\",\"name\":\"field1\",\"type\":\"int32\"},{\"internalType\":\"int32\",\"name\":\"field2\",\"type\":\"int32\"},{\"internalType\":\"int32\",\"name\":\"field3\",\"type\":\"int32\"}],\"name\":\"triggerWithFourTopics\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"field1\",\"type\":\"string\"},{\"internalType\":\"uint8[32]\",\"name\":\"field2\",\"type\":\"uint8[32]\"},{\"internalType\":\"bytes32\",\"name\":\"field3\",\"type\":\"bytes32\"}],\"name\":\"triggerWithFourTopicsWithHashed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50600180548082018255600082905260048082047fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6908101805460086003958616810261010090810a8088026001600160401b0391820219909416939093179093558654808801909755848704909301805496909516909202900a918202910219909216919091179055611e55806100a96000396000f3fe608060405234801561001057600080fd5b50600436106100d45760003560e01c80636c9a43b611610081578063dbfd73321161005b578063dbfd7332146101de578063ef4e1ced146101f1578063fbe9fbf6146101f857600080fd5b80636c9a43b614610165578063a90e1998146101ae578063ab5e0b38146101c157600080fd5b80633272b66c116100b25780633272b66c1461012a57806338e372ba1461013d578063679004a41461015057600080fd5b8063236aadd6146100d95780632c45576f1461010257806331a2c37914610115575b600080fd5b6100ec6100e7366004610faf565b61020a565b6040516100f99190611222565b60405180910390f35b6100ec61011036600461138b565b610338565b6101286101233660046113a4565b6106a8565b005b6101286101383660046114b4565b610705565b61012861014b366004610faf565b61075a565b610158610b42565b6040516100f991906114f6565b610128610173366004611544565b600280547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001667ffffffffffffffff92909216919091179055565b6101286101bc366004611678565b610bce565b6107c65b60405167ffffffffffffffff90911681526020016100f9565b6101286101ec36600461172d565b610c28565b60036101c5565b60025467ffffffffffffffff166101c5565b610212610c65565b6040518061014001604052808e60030b81526020018d8d8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050509082525060ff8c166020808301919091526040805161040081810183529190930192918d9183908390808284376000920191909152505050815273ffffffffffffffffffffffffffffffffffffffff808b16602080840191909152908a1660408084019190915280518983028181018401909252898152606090930192918a918a91829190850190849080828437600092019190915250505090825250601786900b602082015260400161030f856117b2565b815260200161032336859003850185611850565b905290505b9c9b505050505050505050505050565b610340610c65565b600061034d6001846118e0565b8154811061035d5761035d61191a565b6000918252602091829020604080516101408101909152600d90920201805460030b8252600181018054929391929184019161039890611949565b80601f01602080910402602001604051908101604052809291908181526020018280546103c490611949565b80156104115780601f106103e657610100808354040283529160200191610411565b820191906000526020600020905b8154815290600101906020018083116103f457829003601f168201915b5050509183525050600282015460ff166020808301919091526040805161040081018083529190930192916003850191826000855b825461010083900a900460ff1681526020600192830181810494850194909303909202910180841161044657505050928452505050600482015473ffffffffffffffffffffffffffffffffffffffff9081166020808401919091526005840154909116604080840191909152600684018054825181850281018501909352808352606090940193919290919083018282801561051857602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116104ed575b505050918352505060078281015460170b60208084019190915260408051808201825260088601805460f01b7fffff0000000000000000000000000000000000000000000000000000000000001682528251808401845260098801805490960b8152600a880180549490970196929591948681019491939084019161059c90611949565b80601f01602080910402602001604051908101604052809291908181526020018280546105c890611949565b80156106155780601f106105ea57610100808354040283529160200191610615565b820191906000526020600020905b8154815290600101906020018083116105f857829003601f168201915b505050919092525050509052508152604080518082018252600b84015460f01b7fffff0000000000000000000000000000000000000000000000000000000000001681528151808301909252600c90930154600781900b825268010000000000000000900473ffffffffffffffffffffffffffffffffffffffff1660208083019190915280840191909152015292915050565b8b60030b7f855ac250d95b464eea2a7645e23a88fdab21031016175b7dc4d65e8efc72c2ea8c8c8c8c8c8c8c8c8c8c8c6040516106ef9b9a99989796959493929190611ac6565b60405180910390a2505050505050505050505050565b8181604051610715929190611c63565b60405180910390207f3d969732b1bbbb9f1d7eb9f3f14e4cb50a74d950b3ef916a397b85dfbab93c67838360405161074e929190611c73565b60405180910390a25050565b60006040518061014001604052808e60030b81526020018d8d8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050509082525060ff8c166020808301919091526040805161040081810183529190930192918d9183908390808284376000920191909152505050815273ffffffffffffffffffffffffffffffffffffffff808b16602080840191909152908a1660408084019190915280518983028181018401909252898152606090930192918a918a91829190850190849080828437600092019190915250505090825250601786900b6020820152604001610859856117b2565b815260200161086d36859003850185611850565b905281546001808201845560009384526020938490208351600d9093020180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff9093169290921782559282015191929091908201906108d39082611cd6565b5060408201516002820180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff90921691909117905560608201516109219060038301906020610ced565b50608082015160048201805473ffffffffffffffffffffffffffffffffffffffff9283167fffffffffffffffffffffffff00000000000000000000000000000000000000009182161790915560a084015160058401805491909316911617905560c0820151805161099c916006840191602090910190610d80565b5060e08201516007820180547fffffffffffffffff0000000000000000000000000000000000000000000000001677ffffffffffffffffffffffffffffffffffffffffffffffff90921691909117905561010082015180516008830180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001660f09290921c91909117815560208083015180516009860180547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001667ffffffffffffffff90921691909117815591810151909190600a860190610a809082611cd6565b5050505061012092909201518051600b8301805460f09290921c7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00009092169190911790556020908101518051600c9093018054919092015173ffffffffffffffffffffffffffffffffffffffff1668010000000000000000027fffffffff0000000000000000000000000000000000000000000000000000000090911667ffffffffffffffff9093169290921791909117905550505050505050505050505050565b60606001805480602002602001604051908101604052809291908181526020018280548015610bc457602002820191906000526020600020906000905b82829054906101000a900467ffffffffffffffff1667ffffffffffffffff1681526020019060080190602082600701049283019260010382029150808411610b7f5790505b5050505050905090565b8082604051610bdd9190611df0565b604051809103902084604051610bf39190611e2c565b604051908190038120907f7220e4dbe4e9d0ed5f71acd022bc89c26748ac6784f2c548bc17bb8e52af34b090600090a4505050565b8060030b8260030b8460030b7f91c80dc390f3d041b3a04b0099b19634499541ea26972250986ee4b24a12fac560405160405180910390a4505050565b6040805161014081018252600080825260606020830181905292820152908101610c8d610dfa565b815260006020820181905260408201819052606080830152608082015260a001610cb5610e19565b8152602001610ce86040805180820182526000808252825180840190935280835260208381019190915290919082015290565b905290565b600183019183908215610d705791602002820160005b83821115610d4157835183826101000a81548160ff021916908360ff1602179055509260200192600101602081600001049283019260010302610d03565b8015610d6e5782816101000a81549060ff0219169055600101602081600001049283019260010302610d41565b505b50610d7c929150610e6c565b5090565b828054828255906000526020600020908101928215610d70579160200282015b82811115610d7057825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190610da0565b6040518061040001604052806020906020820280368337509192915050565b604051806040016040528060007dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168152602001610ce86040518060400160405280600060070b8152602001606081525090565b5b80821115610d7c5760008155600101610e6d565b8035600381900b8114610e9357600080fd5b919050565b60008083601f840112610eaa57600080fd5b50813567ffffffffffffffff811115610ec257600080fd5b602083019150836020828501011115610eda57600080fd5b9250929050565b803560ff81168114610e9357600080fd5b806104008101831015610f0457600080fd5b92915050565b803573ffffffffffffffffffffffffffffffffffffffff81168114610e9357600080fd5b60008083601f840112610f4057600080fd5b50813567ffffffffffffffff811115610f5857600080fd5b6020830191508360208260051b8501011115610eda57600080fd5b8035601781900b8114610e9357600080fd5b600060408284031215610f9757600080fd5b50919050565b600060608284031215610f9757600080fd5b6000806000806000806000806000806000806105608d8f031215610fd257600080fd5b610fdb8d610e81565b9b5067ffffffffffffffff60208e01351115610ff657600080fd5b6110068e60208f01358f01610e98565b909b50995061101760408e01610ee1565b98506110268e60608f01610ef2565b97506110356104608e01610f0a565b96506110446104808e01610f0a565b955067ffffffffffffffff6104a08e0135111561106057600080fd5b6110718e6104a08f01358f01610f2e565b90955093506110836104c08e01610f73565b925067ffffffffffffffff6104e08e0135111561109f57600080fd5b6110b08e6104e08f01358f01610f85565b91506110c08e6105008f01610f9d565b90509295989b509295989b509295989b565b60005b838110156110ed5781810151838201526020016110d5565b50506000910152565b6000815180845261110e8160208601602086016110d2565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b8060005b6020808210611153575061116a565b825160ff1685529384019390910190600101611144565b50505050565b600081518084526020808501945080840160005b838110156111b657815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101611184565b509495945050505050565b7fffff00000000000000000000000000000000000000000000000000000000000081511682526000602082015160406020850152805160070b6040850152602081015190506040606085015261121a60808501826110f6565b949350505050565b6020815261123660208201835160030b9052565b6000602083015161056060408401526112536105808401826110f6565b90506040840151611269606085018260ff169052565b50606084015161127c6080850182611140565b50608084015173ffffffffffffffffffffffffffffffffffffffff90811661048085015260a0850151166104a084015260c08401518382037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09081016104c08601526112e88383611170565b925060e086015191506113016104e086018360170b9052565b610100860151915080858403016105008601525061131f82826111c1565b61012086015180517fffff00000000000000000000000000000000000000000000000000000000000016610520870152602080820151805160070b610540890152015173ffffffffffffffffffffffffffffffffffffffff166105608701529092509050509392505050565b60006020828403121561139d57600080fd5b5035919050565b6000806000806000806000806000806000806105608d8f0312156113c757600080fd5b6113d08d610e81565b9b506113de60208e01610ee1565b9a5067ffffffffffffffff60408e013511156113f957600080fd5b6114098e60408f01358f01610f85565b99506114188e60608f01610f9d565b98506114278e60c08f01610ef2565b97506114366104c08e01610f0a565b96506114456104e08e01610f0a565b955067ffffffffffffffff6105008e0135111561146157600080fd5b6114728e6105008f01358f01610f2e565b909550935067ffffffffffffffff6105208e0135111561149157600080fd5b6114a28e6105208f01358f01610e98565b90935091506110c06105408e01610f73565b600080602083850312156114c757600080fd5b823567ffffffffffffffff8111156114de57600080fd5b6114ea85828601610e98565b90969095509350505050565b6020808252825182820181905260009190848201906040850190845b8181101561153857835167ffffffffffffffff1683529284019291840191600101611512565b50909695505050505050565b60006020828403121561155657600080fd5b813567ffffffffffffffff8116811461156e57600080fd5b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040805190810167ffffffffffffffff811182821017156115c7576115c7611575565b60405290565b600082601f8301126115de57600080fd5b813567ffffffffffffffff808211156115f9576115f9611575565b604051601f83017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f0116810190828211818310171561163f5761163f611575565b8160405283815286602085880101111561165857600080fd5b836020870160208301376000602085830101528094505050505092915050565b6000806000610440848603121561168e57600080fd5b833567ffffffffffffffff808211156116a657600080fd5b6116b2878388016115cd565b94506020915086603f8701126116c757600080fd5b604051610400810181811083821117156116e3576116e3611575565b6040529050806104208701888111156116fb57600080fd5b8388015b8181101561171d5761171081610ee1565b84529284019284016116ff565b5095989097509435955050505050565b60008060006060848603121561174257600080fd5b61174b84610e81565b925061175960208501610e81565b915061176760408501610e81565b90509250925092565b80357fffff00000000000000000000000000000000000000000000000000000000000081168114610e9357600080fd5b8035600781900b8114610e9357600080fd5b6000604082360312156117c457600080fd5b6117cc6115a4565b6117d583611770565b8152602083013567ffffffffffffffff808211156117f257600080fd5b81850191506040823603121561180757600080fd5b61180f6115a4565b611818836117a0565b815260208301358281111561182c57600080fd5b611838368286016115cd565b60208301525080602085015250505080915050919050565b6000818303606081121561186357600080fd5b61186b6115a4565b61187484611770565b815260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0830112156118a657600080fd5b6118ae6115a4565b91506118bc602085016117a0565b82526118ca60408501610f0a565b6020830152816020820152809250505092915050565b81810381811115610f04577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600181811c9082168061195d57607f821691505b602082108103610f97577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b7fffff000000000000000000000000000000000000000000000000000000000000611a0982611770565b168252611a18602082016117a0565b60070b602083015273ffffffffffffffffffffffffffffffffffffffff611a4160408301610f0a565b1660408301525050565b8060005b6020808210611a5e575061116a565b60ff611a6984610ee1565b168552938401939190910190600101611a4f565b8183526000602080850194508260005b858110156111b65773ffffffffffffffffffffffffffffffffffffffff611ab383610f0a565b1687529582019590820190600101611a8d565b600061054060ff8e1683528060208401527fffff000000000000000000000000000000000000000000000000000000000000611b018e611770565b16818401525060208c01357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc18d3603018112611b3c57600080fd5b60406105608401528c01611b4f816117a0565b60070b61058084015260208101357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1823603018112611b8d57600080fd5b0160208101903567ffffffffffffffff811115611ba957600080fd5b803603821315611bb857600080fd5b60406105a0850152611bcf6105c085018284611996565b915050611bdf604084018d6119df565b611bec60a084018c611a4b565b73ffffffffffffffffffffffffffffffffffffffff8a166104a084015273ffffffffffffffffffffffffffffffffffffffff89166104c08401528281036104e0840152611c3a81888a611a7d565b9050828103610500840152611c50818688611996565b91505061032861052083018460170b9052565b8183823760009101908152919050565b60208152600061121a602083018486611996565b601f821115611cd157600081815260208120601f850160051c81016020861015611cae5750805b601f850160051c820191505b81811015611ccd57828155600101611cba565b5050505b505050565b815167ffffffffffffffff811115611cf057611cf0611575565b611d0481611cfe8454611949565b84611c87565b602080601f831160018114611d575760008415611d215750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555611ccd565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015611da457888601518255948401946001909101908401611d85565b5085821015611de057878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b60008183825b6020808210611e055750611e1c565b825160ff1684529283019290910190600101611df6565b5050506104008201905092915050565b60008251611e3e8184602087016110d2565b919091019291505056fea164736f6c6343000813000a",
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

func (_ChainReaderTester *ChainReaderTesterCaller) ReturnSeen(opts *bind.CallOpts, field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accountStr common.Address, accounts []common.Address, bigField *big.Int, nestedDynamicStruct MidLevelDynamicTestStruct, nestedStaticStruct MidLevelStaticTestStruct) (TestStruct, error) {
	var out []interface{}
	err := _ChainReaderTester.contract.Call(opts, &out, "returnSeen", field, differentField, oracleId, oracleIds, account, accountStr, accounts, bigField, nestedDynamicStruct, nestedStaticStruct)

	if err != nil {
		return *new(TestStruct), err
	}

	out0 := *abi.ConvertType(out[0], new(TestStruct)).(*TestStruct)

	return out0, err

}

func (_ChainReaderTester *ChainReaderTesterSession) ReturnSeen(field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accountStr common.Address, accounts []common.Address, bigField *big.Int, nestedDynamicStruct MidLevelDynamicTestStruct, nestedStaticStruct MidLevelStaticTestStruct) (TestStruct, error) {
	return _ChainReaderTester.Contract.ReturnSeen(&_ChainReaderTester.CallOpts, field, differentField, oracleId, oracleIds, account, accountStr, accounts, bigField, nestedDynamicStruct, nestedStaticStruct)
}

func (_ChainReaderTester *ChainReaderTesterCallerSession) ReturnSeen(field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accountStr common.Address, accounts []common.Address, bigField *big.Int, nestedDynamicStruct MidLevelDynamicTestStruct, nestedStaticStruct MidLevelStaticTestStruct) (TestStruct, error) {
	return _ChainReaderTester.Contract.ReturnSeen(&_ChainReaderTester.CallOpts, field, differentField, oracleId, oracleIds, account, accountStr, accounts, bigField, nestedDynamicStruct, nestedStaticStruct)
}

func (_ChainReaderTester *ChainReaderTesterTransactor) AddTestStruct(opts *bind.TransactOpts, field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accountStr common.Address, accounts []common.Address, bigField *big.Int, nestedDynamicStruct MidLevelDynamicTestStruct, nestedStaticStruct MidLevelStaticTestStruct) (*types.Transaction, error) {
	return _ChainReaderTester.contract.Transact(opts, "addTestStruct", field, differentField, oracleId, oracleIds, account, accountStr, accounts, bigField, nestedDynamicStruct, nestedStaticStruct)
}

func (_ChainReaderTester *ChainReaderTesterSession) AddTestStruct(field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accountStr common.Address, accounts []common.Address, bigField *big.Int, nestedDynamicStruct MidLevelDynamicTestStruct, nestedStaticStruct MidLevelStaticTestStruct) (*types.Transaction, error) {
	return _ChainReaderTester.Contract.AddTestStruct(&_ChainReaderTester.TransactOpts, field, differentField, oracleId, oracleIds, account, accountStr, accounts, bigField, nestedDynamicStruct, nestedStaticStruct)
}

func (_ChainReaderTester *ChainReaderTesterTransactorSession) AddTestStruct(field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accountStr common.Address, accounts []common.Address, bigField *big.Int, nestedDynamicStruct MidLevelDynamicTestStruct, nestedStaticStruct MidLevelStaticTestStruct) (*types.Transaction, error) {
	return _ChainReaderTester.Contract.AddTestStruct(&_ChainReaderTester.TransactOpts, field, differentField, oracleId, oracleIds, account, accountStr, accounts, bigField, nestedDynamicStruct, nestedStaticStruct)
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

func (_ChainReaderTester *ChainReaderTesterTransactor) TriggerEvent(opts *bind.TransactOpts, field int32, oracleId uint8, nestedDynamicStruct MidLevelDynamicTestStruct, nestedStaticStruct MidLevelStaticTestStruct, oracleIds [32]uint8, account common.Address, accountStr common.Address, accounts []common.Address, differentField string, bigField *big.Int) (*types.Transaction, error) {
	return _ChainReaderTester.contract.Transact(opts, "triggerEvent", field, oracleId, nestedDynamicStruct, nestedStaticStruct, oracleIds, account, accountStr, accounts, differentField, bigField)
}

func (_ChainReaderTester *ChainReaderTesterSession) TriggerEvent(field int32, oracleId uint8, nestedDynamicStruct MidLevelDynamicTestStruct, nestedStaticStruct MidLevelStaticTestStruct, oracleIds [32]uint8, account common.Address, accountStr common.Address, accounts []common.Address, differentField string, bigField *big.Int) (*types.Transaction, error) {
	return _ChainReaderTester.Contract.TriggerEvent(&_ChainReaderTester.TransactOpts, field, oracleId, nestedDynamicStruct, nestedStaticStruct, oracleIds, account, accountStr, accounts, differentField, bigField)
}

func (_ChainReaderTester *ChainReaderTesterTransactorSession) TriggerEvent(field int32, oracleId uint8, nestedDynamicStruct MidLevelDynamicTestStruct, nestedStaticStruct MidLevelStaticTestStruct, oracleIds [32]uint8, account common.Address, accountStr common.Address, accounts []common.Address, differentField string, bigField *big.Int) (*types.Transaction, error) {
	return _ChainReaderTester.Contract.TriggerEvent(&_ChainReaderTester.TransactOpts, field, oracleId, nestedDynamicStruct, nestedStaticStruct, oracleIds, account, accountStr, accounts, differentField, bigField)
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
	AccountStr          common.Address
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
	return common.HexToHash("0x855ac250d95b464eea2a7645e23a88fdab21031016175b7dc4d65e8efc72c2ea")
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

	ReturnSeen(opts *bind.CallOpts, field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accountStr common.Address, accounts []common.Address, bigField *big.Int, nestedDynamicStruct MidLevelDynamicTestStruct, nestedStaticStruct MidLevelStaticTestStruct) (TestStruct, error)

	AddTestStruct(opts *bind.TransactOpts, field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accountStr common.Address, accounts []common.Address, bigField *big.Int, nestedDynamicStruct MidLevelDynamicTestStruct, nestedStaticStruct MidLevelStaticTestStruct) (*types.Transaction, error)

	SetAlterablePrimitiveValue(opts *bind.TransactOpts, value uint64) (*types.Transaction, error)

	TriggerEvent(opts *bind.TransactOpts, field int32, oracleId uint8, nestedDynamicStruct MidLevelDynamicTestStruct, nestedStaticStruct MidLevelStaticTestStruct, oracleIds [32]uint8, account common.Address, accountStr common.Address, accounts []common.Address, differentField string, bigField *big.Int) (*types.Transaction, error)

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
