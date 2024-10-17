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
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"StaticBytes\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int32\",\"name\":\"field\",\"type\":\"int32\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"oracleId\",\"type\":\"uint8\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerDynamicTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"indexed\":false,\"internalType\":\"structMidLevelDynamicTestStruct\",\"name\":\"nestedDynamicStruct\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"address\",\"name\":\"A\",\"type\":\"address\"}],\"internalType\":\"structInnerStaticTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"indexed\":false,\"internalType\":\"structMidLevelStaticTestStruct\",\"name\":\"nestedStaticStruct\",\"type\":\"tuple\"},{\"indexed\":false,\"internalType\":\"uint8[32]\",\"name\":\"oracleIds\",\"type\":\"uint8[32]\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"Account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"AccountStr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"Accounts\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"differentField\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"int192\",\"name\":\"bigField\",\"type\":\"int192\"}],\"name\":\"Triggered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"fieldHash\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"field\",\"type\":\"string\"}],\"name\":\"TriggeredEventWithDynamicTopic\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int32\",\"name\":\"field1\",\"type\":\"int32\"},{\"indexed\":true,\"internalType\":\"int32\",\"name\":\"field2\",\"type\":\"int32\"},{\"indexed\":true,\"internalType\":\"int32\",\"name\":\"field3\",\"type\":\"int32\"}],\"name\":\"TriggeredWithFourTopics\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"field1\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"uint8[32]\",\"name\":\"field2\",\"type\":\"uint8[32]\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"field3\",\"type\":\"bytes32\"}],\"name\":\"TriggeredWithFourTopicsWithHashed\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"int32\",\"name\":\"field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"differentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"oracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"oracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"accountStr\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"int192\",\"name\":\"bigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerDynamicTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelDynamicTestStruct\",\"name\":\"nestedDynamicStruct\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"address\",\"name\":\"A\",\"type\":\"address\"}],\"internalType\":\"structInnerStaticTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelStaticTestStruct\",\"name\":\"nestedStaticStruct\",\"type\":\"tuple\"}],\"name\":\"addTestStruct\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAlterablePrimitiveValue\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDifferentPrimitiveValue\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"i\",\"type\":\"uint256\"}],\"name\":\"getElementAtIndex\",\"outputs\":[{\"components\":[{\"internalType\":\"int32\",\"name\":\"Field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"DifferentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"OracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"OracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"address\",\"name\":\"Account\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"AccountStr\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"Accounts\",\"type\":\"address[]\"},{\"internalType\":\"int192\",\"name\":\"BigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerDynamicTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelDynamicTestStruct\",\"name\":\"NestedDynamicStruct\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"address\",\"name\":\"A\",\"type\":\"address\"}],\"internalType\":\"structInnerStaticTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelStaticTestStruct\",\"name\":\"NestedStaticStruct\",\"type\":\"tuple\"}],\"internalType\":\"structTestStruct\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPrimitiveValue\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSliceValue\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int32\",\"name\":\"field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"differentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"oracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"oracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"accountStr\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"int192\",\"name\":\"bigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerDynamicTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelDynamicTestStruct\",\"name\":\"nestedDynamicStruct\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"address\",\"name\":\"A\",\"type\":\"address\"}],\"internalType\":\"structInnerStaticTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelStaticTestStruct\",\"name\":\"nestedStaticStruct\",\"type\":\"tuple\"}],\"name\":\"returnSeen\",\"outputs\":[{\"components\":[{\"internalType\":\"int32\",\"name\":\"Field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"DifferentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"OracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"OracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"address\",\"name\":\"Account\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"AccountStr\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"Accounts\",\"type\":\"address[]\"},{\"internalType\":\"int192\",\"name\":\"BigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerDynamicTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelDynamicTestStruct\",\"name\":\"NestedDynamicStruct\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"address\",\"name\":\"A\",\"type\":\"address\"}],\"internalType\":\"structInnerStaticTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelStaticTestStruct\",\"name\":\"NestedStaticStruct\",\"type\":\"tuple\"}],\"internalType\":\"structTestStruct\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"value\",\"type\":\"uint64\"}],\"name\":\"setAlterablePrimitiveValue\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int32\",\"name\":\"field\",\"type\":\"int32\"},{\"internalType\":\"uint8\",\"name\":\"oracleId\",\"type\":\"uint8\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerDynamicTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelDynamicTestStruct\",\"name\":\"nestedDynamicStruct\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"address\",\"name\":\"A\",\"type\":\"address\"}],\"internalType\":\"structInnerStaticTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelStaticTestStruct\",\"name\":\"nestedStaticStruct\",\"type\":\"tuple\"},{\"internalType\":\"uint8[32]\",\"name\":\"oracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"accountStr\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"string\",\"name\":\"differentField\",\"type\":\"string\"},{\"internalType\":\"int192\",\"name\":\"bigField\",\"type\":\"int192\"}],\"name\":\"triggerEvent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"field\",\"type\":\"string\"}],\"name\":\"triggerEventWithDynamicTopic\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"val1\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"val2\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"val3\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"val4\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"val5\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"val6\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"val7\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"raw\",\"type\":\"bytes\"}],\"name\":\"triggerStaticBytes\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int32\",\"name\":\"field1\",\"type\":\"int32\"},{\"internalType\":\"int32\",\"name\":\"field2\",\"type\":\"int32\"},{\"internalType\":\"int32\",\"name\":\"field3\",\"type\":\"int32\"}],\"name\":\"triggerWithFourTopics\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"field1\",\"type\":\"string\"},{\"internalType\":\"uint8[32]\",\"name\":\"field2\",\"type\":\"uint8[32]\"},{\"internalType\":\"bytes32\",\"name\":\"field3\",\"type\":\"bytes32\"}],\"name\":\"triggerWithFourTopicsWithHashed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50600180548082018255600082905260048082047fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6908101805460086003958616810261010090810a8088026001600160401b0391820219909416939093179093558654808801909755848704909301805496909516909202900a918202910219909216919091179055612077806100a96000396000f3fe608060405234801561001057600080fd5b50600436106100df5760003560e01c8063679004a41161008c578063ab5e0b3811610066578063ab5e0b38146101df578063dbfd7332146101fc578063ef4e1ced1461020f578063fbe9fbf61461021657600080fd5b8063679004a41461016e5780636c9a43b614610183578063a90e1998146101cc57600080fd5b80633272b66c116100bd5780633272b66c1461013557806338e372ba1461014857806351f3f54d1461015b57600080fd5b8063236aadd6146100e45780632c45576f1461010d57806331a2c37914610120575b600080fd5b6100f76100f2366004611041565b610228565b60405161010491906112b4565b60405180910390f35b6100f761011b36600461141d565b610356565b61013361012e366004611436565b6106c6565b005b610133610143366004611546565b610723565b610133610156366004611041565b610778565b6101336101693660046116a0565b610b60565b610176610bd4565b6040516101049190611751565b61013361019136600461179f565b600280547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001667ffffffffffffffff92909216919091179055565b6101336101da3660046117e1565b610c60565b6107c65b60405167ffffffffffffffff9091168152602001610104565b61013361020a366004611896565b610cba565b60036101e3565b60025467ffffffffffffffff166101e3565b610230610cf7565b6040518061014001604052808e60030b81526020018d8d8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050509082525060ff8c166020808301919091526040805161040081810183529190930192918d9183908390808284376000920191909152505050815273ffffffffffffffffffffffffffffffffffffffff808b16602080840191909152908a1660408084019190915280518983028181018401909252898152606090930192918a918a91829190850190849080828437600092019190915250505090825250601786900b602082015260400161032d8561191b565b8152602001610341368590038501856119b9565b905290505b9c9b505050505050505050505050565b61035e610cf7565b600061036b600184611a49565b8154811061037b5761037b611a83565b6000918252602091829020604080516101408101909152600d90920201805460030b825260018101805492939192918401916103b690611ab2565b80601f01602080910402602001604051908101604052809291908181526020018280546103e290611ab2565b801561042f5780601f106104045761010080835404028352916020019161042f565b820191906000526020600020905b81548152906001019060200180831161041257829003601f168201915b5050509183525050600282015460ff166020808301919091526040805161040081018083529190930192916003850191826000855b825461010083900a900460ff1681526020600192830181810494850194909303909202910180841161046457505050928452505050600482015473ffffffffffffffffffffffffffffffffffffffff9081166020808401919091526005840154909116604080840191909152600684018054825181850281018501909352808352606090940193919290919083018282801561053657602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff16815260019091019060200180831161050b575b505050918352505060078281015460170b60208084019190915260408051808201825260088601805460f01b7fffff0000000000000000000000000000000000000000000000000000000000001682528251808401845260098801805490960b8152600a88018054949097019692959194868101949193908401916105ba90611ab2565b80601f01602080910402602001604051908101604052809291908181526020018280546105e690611ab2565b80156106335780601f1061060857610100808354040283529160200191610633565b820191906000526020600020905b81548152906001019060200180831161061657829003601f168201915b505050919092525050509052508152604080518082018252600b84015460f01b7fffff0000000000000000000000000000000000000000000000000000000000001681528151808301909252600c90930154600781900b825268010000000000000000900473ffffffffffffffffffffffffffffffffffffffff1660208083019190915280840191909152015292915050565b8b60030b7f855ac250d95b464eea2a7645e23a88fdab21031016175b7dc4d65e8efc72c2ea8c8c8c8c8c8c8c8c8c8c8c60405161070d9b9a99989796959493929190611c2f565b60405180910390a2505050505050505050505050565b8181604051610733929190611dcc565b60405180910390207f3d969732b1bbbb9f1d7eb9f3f14e4cb50a74d950b3ef916a397b85dfbab93c67838360405161076c929190611ddc565b60405180910390a25050565b60006040518061014001604052808e60030b81526020018d8d8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050509082525060ff8c166020808301919091526040805161040081810183529190930192918d9183908390808284376000920191909152505050815273ffffffffffffffffffffffffffffffffffffffff808b16602080840191909152908a1660408084019190915280518983028181018401909252898152606090930192918a918a91829190850190849080828437600092019190915250505090825250601786900b60208201526040016108778561191b565b815260200161088b368590038501856119b9565b905281546001808201845560009384526020938490208351600d9093020180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff9093169290921782559282015191929091908201906108f19082611e3f565b5060408201516002820180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff909216919091179055606082015161093f9060038301906020610d7f565b50608082015160048201805473ffffffffffffffffffffffffffffffffffffffff9283167fffffffffffffffffffffffff00000000000000000000000000000000000000009182161790915560a084015160058401805491909316911617905560c082015180516109ba916006840191602090910190610e12565b5060e08201516007820180547fffffffffffffffff0000000000000000000000000000000000000000000000001677ffffffffffffffffffffffffffffffffffffffffffffffff90921691909117905561010082015180516008830180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001660f09290921c91909117815560208083015180516009860180547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001667ffffffffffffffff90921691909117815591810151909190600a860190610a9e9082611e3f565b5050505061012092909201518051600b8301805460f09290921c7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00009092169190911790556020908101518051600c9093018054919092015173ffffffffffffffffffffffffffffffffffffffff1668010000000000000000027fffffffff0000000000000000000000000000000000000000000000000000000090911667ffffffffffffffff9093169290921791909117905550505050505050505050505050565b60008888888888888888604051602001610b81989796959493929190611f59565b60405160208183030381529060405290507f1e40927ec0bdc7319f09a53452590433ec395dec3b70b982eba779c740685bfe81604051610bc19190611fff565b60405180910390a1505050505050505050565b60606001805480602002602001604051908101604052809291908181526020018280548015610c5657602002820191906000526020600020906000905b82829054906101000a900467ffffffffffffffff1667ffffffffffffffff1681526020019060080190602082600701049283019260010382029150808411610c115790505b5050505050905090565b8082604051610c6f9190612012565b604051809103902084604051610c85919061204e565b604051908190038120907f7220e4dbe4e9d0ed5f71acd022bc89c26748ac6784f2c548bc17bb8e52af34b090600090a4505050565b8060030b8260030b8460030b7f91c80dc390f3d041b3a04b0099b19634499541ea26972250986ee4b24a12fac560405160405180910390a4505050565b6040805161014081018252600080825260606020830181905292820152908101610d1f610e8c565b815260006020820181905260408201819052606080830152608082015260a001610d47610eab565b8152602001610d7a6040805180820182526000808252825180840190935280835260208381019190915290919082015290565b905290565b600183019183908215610e025791602002820160005b83821115610dd357835183826101000a81548160ff021916908360ff1602179055509260200192600101602081600001049283019260010302610d95565b8015610e005782816101000a81549060ff0219169055600101602081600001049283019260010302610dd3565b505b50610e0e929150610efe565b5090565b828054828255906000526020600020908101928215610e02579160200282015b82811115610e0257825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190610e32565b6040518061040001604052806020906020820280368337509192915050565b604051806040016040528060007dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168152602001610d7a6040518060400160405280600060070b8152602001606081525090565b5b80821115610e0e5760008155600101610eff565b8035600381900b8114610f2557600080fd5b919050565b60008083601f840112610f3c57600080fd5b50813567ffffffffffffffff811115610f5457600080fd5b602083019150836020828501011115610f6c57600080fd5b9250929050565b803560ff81168114610f2557600080fd5b806104008101831015610f9657600080fd5b92915050565b803573ffffffffffffffffffffffffffffffffffffffff81168114610f2557600080fd5b60008083601f840112610fd257600080fd5b50813567ffffffffffffffff811115610fea57600080fd5b6020830191508360208260051b8501011115610f6c57600080fd5b8035601781900b8114610f2557600080fd5b60006040828403121561102957600080fd5b50919050565b60006060828403121561102957600080fd5b6000806000806000806000806000806000806105608d8f03121561106457600080fd5b61106d8d610f13565b9b5067ffffffffffffffff60208e0135111561108857600080fd5b6110988e60208f01358f01610f2a565b909b5099506110a960408e01610f73565b98506110b88e60608f01610f84565b97506110c76104608e01610f9c565b96506110d66104808e01610f9c565b955067ffffffffffffffff6104a08e013511156110f257600080fd5b6111038e6104a08f01358f01610fc0565b90955093506111156104c08e01611005565b925067ffffffffffffffff6104e08e0135111561113157600080fd5b6111428e6104e08f01358f01611017565b91506111528e6105008f0161102f565b90509295989b509295989b509295989b565b60005b8381101561117f578181015183820152602001611167565b50506000910152565b600081518084526111a0816020860160208601611164565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b8060005b60208082106111e557506111fc565b825160ff16855293840193909101906001016111d6565b50505050565b600081518084526020808501945080840160005b8381101561124857815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101611216565b509495945050505050565b7fffff00000000000000000000000000000000000000000000000000000000000081511682526000602082015160406020850152805160070b604085015260208101519050604060608501526112ac6080850182611188565b949350505050565b602081526112c860208201835160030b9052565b6000602083015161056060408401526112e5610580840182611188565b905060408401516112fb606085018260ff169052565b50606084015161130e60808501826111d2565b50608084015173ffffffffffffffffffffffffffffffffffffffff90811661048085015260a0850151166104a084015260c08401518382037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09081016104c086015261137a8383611202565b925060e086015191506113936104e086018360170b9052565b61010086015191508085840301610500860152506113b18282611253565b61012086015180517fffff00000000000000000000000000000000000000000000000000000000000016610520870152602080820151805160070b610540890152015173ffffffffffffffffffffffffffffffffffffffff166105608701529092509050509392505050565b60006020828403121561142f57600080fd5b5035919050565b6000806000806000806000806000806000806105608d8f03121561145957600080fd5b6114628d610f13565b9b5061147060208e01610f73565b9a5067ffffffffffffffff60408e0135111561148b57600080fd5b61149b8e60408f01358f01611017565b99506114aa8e60608f0161102f565b98506114b98e60c08f01610f84565b97506114c86104c08e01610f9c565b96506114d76104e08e01610f9c565b955067ffffffffffffffff6105008e013511156114f357600080fd5b6115048e6105008f01358f01610fc0565b909550935067ffffffffffffffff6105208e0135111561152357600080fd5b6115348e6105208f01358f01610f2a565b90935091506111526105408e01611005565b6000806020838503121561155957600080fd5b823567ffffffffffffffff81111561157057600080fd5b61157c85828601610f2a565b90969095509350505050565b803563ffffffff81168114610f2557600080fd5b803567ffffffffffffffff81168114610f2557600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040805190810167ffffffffffffffff81118282101715611606576116066115b4565b60405290565b600067ffffffffffffffff80841115611627576116276115b4565b604051601f85017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f0116810190828211818310171561166d5761166d6115b4565b8160405280935085815286868601111561168657600080fd5b858560208301376000602087830101525050509392505050565b600080600080600080600080610100898b0312156116bd57600080fd5b6116c689611588565b97506116d460208a01611588565b96506116e260408a01611588565b95506116f060608a0161159c565b94506080890135935060a0890135925060c0890135915060e089013567ffffffffffffffff81111561172157600080fd5b8901601f81018b1361173257600080fd5b6117418b82356020840161160c565b9150509295985092959890939650565b6020808252825182820181905260009190848201906040850190845b8181101561179357835167ffffffffffffffff168352928401929184019160010161176d565b50909695505050505050565b6000602082840312156117b157600080fd5b6117ba8261159c565b9392505050565b600082601f8301126117d257600080fd5b6117ba8383356020850161160c565b600080600061044084860312156117f757600080fd5b833567ffffffffffffffff8082111561180f57600080fd5b61181b878388016117c1565b94506020915086603f87011261183057600080fd5b6040516104008101818110838211171561184c5761184c6115b4565b60405290508061042087018881111561186457600080fd5b8388015b818110156118865761187981610f73565b8452928401928401611868565b5095989097509435955050505050565b6000806000606084860312156118ab57600080fd5b6118b484610f13565b92506118c260208501610f13565b91506118d060408501610f13565b90509250925092565b80357fffff00000000000000000000000000000000000000000000000000000000000081168114610f2557600080fd5b8035600781900b8114610f2557600080fd5b60006040823603121561192d57600080fd5b6119356115e3565b61193e836118d9565b8152602083013567ffffffffffffffff8082111561195b57600080fd5b81850191506040823603121561197057600080fd5b6119786115e3565b61198183611909565b815260208301358281111561199557600080fd5b6119a1368286016117c1565b60208301525080602085015250505080915050919050565b600081830360608112156119cc57600080fd5b6119d46115e3565b6119dd846118d9565b815260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe083011215611a0f57600080fd5b611a176115e3565b9150611a2560208501611909565b8252611a3360408501610f9c565b6020830152816020820152809250505092915050565b81810381811115610f96577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600181811c90821680611ac657607f821691505b602082108103611029577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b7fffff000000000000000000000000000000000000000000000000000000000000611b72826118d9565b168252611b8160208201611909565b60070b602083015273ffffffffffffffffffffffffffffffffffffffff611baa60408301610f9c565b1660408301525050565b8060005b6020808210611bc757506111fc565b60ff611bd284610f73565b168552938401939190910190600101611bb8565b8183526000602080850194508260005b858110156112485773ffffffffffffffffffffffffffffffffffffffff611c1c83610f9c565b1687529582019590820190600101611bf6565b600061054060ff8e1683528060208401527fffff000000000000000000000000000000000000000000000000000000000000611c6a8e6118d9565b16818401525060208c01357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc18d3603018112611ca557600080fd5b60406105608401528c01611cb881611909565b60070b61058084015260208101357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1823603018112611cf657600080fd5b0160208101903567ffffffffffffffff811115611d1257600080fd5b803603821315611d2157600080fd5b60406105a0850152611d386105c085018284611aff565b915050611d48604084018d611b48565b611d5560a084018c611bb4565b73ffffffffffffffffffffffffffffffffffffffff8a166104a084015273ffffffffffffffffffffffffffffffffffffffff89166104c08401528281036104e0840152611da381888a611be6565b9050828103610500840152611db9818688611aff565b91505061034661052083018460170b9052565b8183823760009101908152919050565b6020815260006112ac602083018486611aff565b601f821115611e3a57600081815260208120601f850160051c81016020861015611e175750805b601f850160051c820191505b81811015611e3657828155600101611e23565b5050505b505050565b815167ffffffffffffffff811115611e5957611e596115b4565b611e6d81611e678454611ab2565b84611df0565b602080601f831160018114611ec05760008415611e8a5750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555611e36565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015611f0d57888601518255948401946001909101908401611eee565b5085821015611f4957878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b60007fffffffff00000000000000000000000000000000000000000000000000000000808b60e01b168352808a60e01b166004840152808960e01b166008840152507fffffffffffffffff0000000000000000000000000000000000000000000000008760c01b16600c8301528560148301528460348301528360548301528251611feb816074850160208701611164565b919091016074019998505050505050505050565b6020815260006117ba6020830184611188565b60008183825b6020808210612027575061203e565b825160ff1684529283019290910190600101612018565b5050506104008201905092915050565b60008251612060818460208701611164565b919091019291505056fea164736f6c6343000813000a",
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
