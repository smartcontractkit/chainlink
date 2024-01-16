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
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int32\",\"name\":\"field\",\"type\":\"int32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"differentField\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"oracleId\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint8[32]\",\"name\":\"oracleIds\",\"type\":\"uint8[32]\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"Account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"Accounts\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"int192\",\"name\":\"bigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"I\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"indexed\":false,\"internalType\":\"structMidLevelTestStruct\",\"name\":\"nestedStruct\",\"type\":\"tuple\"}],\"name\":\"Triggered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"fieldHash\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"field\",\"type\":\"string\"}],\"name\":\"TriggeredEventWithDynamicTopic\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int32\",\"name\":\"field1\",\"type\":\"int32\"},{\"indexed\":true,\"internalType\":\"int32\",\"name\":\"field2\",\"type\":\"int32\"},{\"indexed\":true,\"internalType\":\"int32\",\"name\":\"field3\",\"type\":\"int32\"}],\"name\":\"TriggeredWithFourTopics\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"int32\",\"name\":\"field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"differentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"oracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"oracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"int192\",\"name\":\"bigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"I\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelTestStruct\",\"name\":\"nestedStruct\",\"type\":\"tuple\"}],\"name\":\"AddTestStruct\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GetDifferentPrimitiveValue\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"i\",\"type\":\"uint256\"}],\"name\":\"GetElementAtIndex\",\"outputs\":[{\"components\":[{\"internalType\":\"int32\",\"name\":\"Field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"DifferentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"OracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"OracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"address\",\"name\":\"Account\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"Accounts\",\"type\":\"address[]\"},{\"internalType\":\"int192\",\"name\":\"BigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"I\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelTestStruct\",\"name\":\"NestedStruct\",\"type\":\"tuple\"}],\"internalType\":\"structTestStruct\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GetPrimitiveValue\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GetSliceValue\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int32\",\"name\":\"field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"differentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"oracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"oracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"int192\",\"name\":\"bigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"I\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelTestStruct\",\"name\":\"nestedStruct\",\"type\":\"tuple\"}],\"name\":\"ReturnSeen\",\"outputs\":[{\"components\":[{\"internalType\":\"int32\",\"name\":\"Field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"DifferentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"OracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"OracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"address\",\"name\":\"Account\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"Accounts\",\"type\":\"address[]\"},{\"internalType\":\"int192\",\"name\":\"BigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"I\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelTestStruct\",\"name\":\"NestedStruct\",\"type\":\"tuple\"}],\"internalType\":\"structTestStruct\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int32\",\"name\":\"field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"differentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"oracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"oracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"int192\",\"name\":\"bigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"I\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelTestStruct\",\"name\":\"nestedStruct\",\"type\":\"tuple\"}],\"name\":\"TriggerEvent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"field\",\"type\":\"string\"}],\"name\":\"TriggerEventWithDynamicTopic\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int32\",\"name\":\"field1\",\"type\":\"int32\"},{\"internalType\":\"int32\",\"name\":\"field2\",\"type\":\"int32\"},{\"internalType\":\"int32\",\"name\":\"field3\",\"type\":\"int32\"}],\"name\":\"TriggerWithFourTopics\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50600180548082018255600082905260048082047fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6908101805460086003958616810261010090810a8088026001600160401b0391820219909416939093179093558654808801909755848704909301805496909516909202900a918202910219909216919091179055611768806100a96000396000f3fe608060405234801561001057600080fd5b50600436106100a35760003560e01c80638b659d6e11610076578063b9dad6b01161005b578063b9dad6b014610138578063bdb37c901461014b578063da8e7a821461016057600080fd5b80638b659d6e146101055780639ca04f671461011857600080fd5b8063030d3ca2146100a85780633205fb6b146100ca5780636c7cf955146100df57806383b0abee146100f2575b600080fd5b6107c65b60405167ffffffffffffffff90911681526020015b60405180910390f35b6100dd6100d8366004610bef565b610167565b005b6100dd6100ed366004610d04565b6101bc565b6100dd610100366004610df6565b6104bf565b6100dd610113366004610d04565b6104fc565b61012b610126366004610e39565b610553565b6040516100c19190610f98565b61012b610146366004610d04565b61082e565b610153610937565b6040516100c1919061108e565b60036100ac565b81816040516101779291906110dc565b60405180910390207f3d969732b1bbbb9f1d7eb9f3f14e4cb50a74d950b3ef916a397b85dfbab93c6783836040516101b0929190611135565b60405180910390a25050565b60006040518061010001604052808c60030b81526020018b8b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050509082525060ff8a166020808301919091526040805161040081810183529190930192918b9183908390808284376000920191909152505050815273ffffffffffffffffffffffffffffffffffffffff8816602080830191909152604080518883028181018401835289825291909301929189918991829190850190849080828437600092019190915250505090825250601785900b60208201526040016102ae84611232565b905281546001808201845560009384526020938490208351600a9093020180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff90931692909217825592820151919290919082019061031490826113d9565b5060408201516002820180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff909216919091179055606082015161036290600383019060206109c3565b5060808201516004820180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff90921691909117905560a082015180516103c9916005840191602090910190610a56565b5060c08201516006820180547fffffffffffffffff0000000000000000000000000000000000000000000000001677ffffffffffffffffffffffffffffffffffffffffffffffff90921691909117905560e082015180516007830180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001660f09290921c91909117815560208083015180516008860180547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001667ffffffffffffffff9092169190911781559181015190919060098601906104ac90826113d9565b5050505050505050505050505050505050565b8060030b8260030b8460030b7f91c80dc390f3d041b3a04b0099b19634499541ea26972250986ee4b24a12fac560405160405180910390a4505050565b8960030b7f7188419dcd8b51877b71766f075f3626586c0ff190e7d056aa65ce9acb649a3d8a8a8a8a8a8a8a8a8a60405161053f99989796959493929190611638565b60405180910390a250505050505050505050565b61055b610ad0565b60006105686001846116f2565b815481106105785761057861172c565b6000918252602091829020604080516101008101909152600a90920201805460030b825260018101805492939192918401916105b39061133d565b80601f01602080910402602001604051908101604052809291908181526020018280546105df9061133d565b801561062c5780601f106106015761010080835404028352916020019161062c565b820191906000526020600020905b81548152906001019060200180831161060f57829003601f168201915b5050509183525050600282015460ff166020808301919091526040805161040081018083529190930192916003850191826000855b825461010083900a900460ff1681526020600192830181810494850194909303909202910180841161066157505050928452505050600482015473ffffffffffffffffffffffffffffffffffffffff16602080830191909152600583018054604080518285028101850182528281529401939283018282801561071a57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116106ef575b5050509183525050600682015460170b6020808301919091526040805180820182526007808601805460f01b7fffff0000000000000000000000000000000000000000000000000000000000001683528351808501855260088801805490930b8152600988018054959097019693959194868301949193928401919061079f9061133d565b80601f01602080910402602001604051908101604052809291908181526020018280546107cb9061133d565b80156108185780601f106107ed57610100808354040283529160200191610818565b820191906000526020600020905b8154815290600101906020018083116107fb57829003601f168201915b5050509190925250505090525090525092915050565b610836610ad0565b6040518061010001604052808c60030b81526020018b8b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050509082525060ff8a166020808301919091526040805161040081810183529190930192918b9183908390808284376000920191909152505050815273ffffffffffffffffffffffffffffffffffffffff8816602080830191909152604080518883028181018401835289825291909301929189918991829190850190849080828437600092019190915250505090825250601785900b602082015260400161092684611232565b90529b9a5050505050505050505050565b606060018054806020026020016040519081016040528092919081815260200182805480156109b957602002820191906000526020600020906000905b82829054906101000a900467ffffffffffffffff1667ffffffffffffffff16815260200190600801906020826007010492830192600103820291508084116109745790505b5050505050905090565b600183019183908215610a465791602002820160005b83821115610a1757835183826101000a81548160ff021916908360ff16021790555092602001926001016020816000010492830192600103026109d9565b8015610a445782816101000a81549060ff0219169055600101602081600001049283019260010302610a17565b505b50610a52929150610b1f565b5090565b828054828255906000526020600020908101928215610a46579160200282015b82811115610a4657825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190610a76565b6040805161010081018252600080825260606020830181905292820152908101610af8610b34565b8152600060208201819052606060408301819052820152608001610b1a610b53565b905290565b5b80821115610a525760008155600101610b20565b6040518061040001604052806020906020820280368337509192915050565b604051806040016040528060007dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168152602001610b1a6040518060400160405280600060070b8152602001606081525090565b60008083601f840112610bb857600080fd5b50813567ffffffffffffffff811115610bd057600080fd5b602083019150836020828501011115610be857600080fd5b9250929050565b60008060208385031215610c0257600080fd5b823567ffffffffffffffff811115610c1957600080fd5b610c2585828601610ba6565b90969095509350505050565b8035600381900b8114610c4357600080fd5b919050565b803560ff81168114610c4357600080fd5b806104008101831015610c6b57600080fd5b92915050565b803573ffffffffffffffffffffffffffffffffffffffff81168114610c4357600080fd5b60008083601f840112610ca757600080fd5b50813567ffffffffffffffff811115610cbf57600080fd5b6020830191508360208260051b8501011115610be857600080fd5b8035601781900b8114610c4357600080fd5b600060408284031215610cfe57600080fd5b50919050565b6000806000806000806000806000806104e08b8d031215610d2457600080fd5b610d2d8b610c31565b995060208b013567ffffffffffffffff80821115610d4a57600080fd5b610d568e838f01610ba6565b909b509950899150610d6a60408e01610c48565b9850610d798e60608f01610c59565b9750610d886104608e01610c71565b96506104808d0135915080821115610d9f57600080fd5b610dab8e838f01610c95565b9096509450849150610dc06104a08e01610cda565b93506104c08d0135915080821115610dd757600080fd5b50610de48d828e01610cec565b9150509295989b9194979a5092959850565b600080600060608486031215610e0b57600080fd5b610e1484610c31565b9250610e2260208501610c31565b9150610e3060408501610c31565b90509250925092565b600060208284031215610e4b57600080fd5b5035919050565b6000815180845260005b81811015610e7857602081850181015186830182015201610e5c565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b8060005b6020808210610ec95750610ee0565b825160ff1685529384019390910190600101610eba565b50505050565b600081518084526020808501945080840160005b83811015610f2c57815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101610efa565b509495945050505050565b7fffff00000000000000000000000000000000000000000000000000000000000081511682526000602082015160406020850152805160070b60408501526020810151905060406060850152610f906080850182610e52565b949350505050565b60208152610fac60208201835160030b9052565b600060208301516104e0806040850152610fca610500850183610e52565b91506040850151610fe0606086018260ff169052565b506060850151610ff36080860182610eb6565b50608085015173ffffffffffffffffffffffffffffffffffffffff1661048085015260a08501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe085840381016104a08701526110508483610ee6565b935060c087015191506110696104c087018360170b9052565b60e08701519150808685030183870152506110848382610f37565b9695505050505050565b6020808252825182820181905260009190848201906040850190845b818110156110d057835167ffffffffffffffff16835292840192918401916001016110aa565b50909695505050505050565b8183823760009101908152919050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b602081526000610f906020830184866110ec565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040805190810167ffffffffffffffff8111828210171561119b5761119b611149565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156111e8576111e8611149565b604052919050565b80357fffff00000000000000000000000000000000000000000000000000000000000081168114610c4357600080fd5b8035600781900b8114610c4357600080fd5b60006040823603121561124457600080fd5b61124c611178565b611255836111f0565b815260208084013567ffffffffffffffff8082111561127357600080fd5b81860191506040823603121561128857600080fd5b611290611178565b61129983611220565b815283830135828111156112ac57600080fd5b929092019136601f8401126112c057600080fd5b8235828111156112d2576112d2611149565b611302857fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116016111a1565b9250808352368582860101111561131857600080fd5b8085850186850137600090830185015280840191909152918301919091525092915050565b600181811c9082168061135157607f821691505b602082108103610cfe577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b601f8211156113d457600081815260208120601f850160051c810160208610156113b15750805b601f850160051c820191505b818110156113d0578281556001016113bd565b5050505b505050565b815167ffffffffffffffff8111156113f3576113f3611149565b61140781611401845461133d565b8461138a565b602080601f83116001811461145a57600084156114245750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b1785556113d0565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b828110156114a757888601518255948401946001909101908401611488565b50858210156114e357878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b8183526000602080850194508260005b85811015610f2c5773ffffffffffffffffffffffffffffffffffffffff61152983610c71565b1687529582019590820190600101611503565b7fffff000000000000000000000000000000000000000000000000000000000000611566826111f0565b168252600060208201357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc18336030181126115a057600080fd5b6040602085015282016115b281611220565b60070b604085015260208101357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18236030181126115ef57600080fd5b0160208101903567ffffffffffffffff81111561160b57600080fd5b80360382131561161a57600080fd5b6040606086015261162f6080860182846110ec565b95945050505050565b60006104c080835261164d8184018c8e6110ec565b9050602060ff808c1682860152604085018b60005b84811015611687578361167483610c48565b1683529184019190840190600101611662565b505050505073ffffffffffffffffffffffffffffffffffffffff88166104408401528281036104608401526116bd8187896114f3565b90506116cf61048084018660170b9052565b8281036104a08401526116e2818561153c565b9c9b505050505050505050505050565b81810381811115610c6b577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fdfea164736f6c6343000813000a",
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

func (_LatestValueHolder *LatestValueHolderTransactor) TriggerEventWithDynamicTopic(opts *bind.TransactOpts, field string) (*types.Transaction, error) {
	return _LatestValueHolder.contract.Transact(opts, "TriggerEventWithDynamicTopic", field)
}

func (_LatestValueHolder *LatestValueHolderSession) TriggerEventWithDynamicTopic(field string) (*types.Transaction, error) {
	return _LatestValueHolder.Contract.TriggerEventWithDynamicTopic(&_LatestValueHolder.TransactOpts, field)
}

func (_LatestValueHolder *LatestValueHolderTransactorSession) TriggerEventWithDynamicTopic(field string) (*types.Transaction, error) {
	return _LatestValueHolder.Contract.TriggerEventWithDynamicTopic(&_LatestValueHolder.TransactOpts, field)
}

func (_LatestValueHolder *LatestValueHolderTransactor) TriggerWithFourTopics(opts *bind.TransactOpts, field1 int32, field2 int32, field3 int32) (*types.Transaction, error) {
	return _LatestValueHolder.contract.Transact(opts, "TriggerWithFourTopics", field1, field2, field3)
}

func (_LatestValueHolder *LatestValueHolderSession) TriggerWithFourTopics(field1 int32, field2 int32, field3 int32) (*types.Transaction, error) {
	return _LatestValueHolder.Contract.TriggerWithFourTopics(&_LatestValueHolder.TransactOpts, field1, field2, field3)
}

func (_LatestValueHolder *LatestValueHolderTransactorSession) TriggerWithFourTopics(field1 int32, field2 int32, field3 int32) (*types.Transaction, error) {
	return _LatestValueHolder.Contract.TriggerWithFourTopics(&_LatestValueHolder.TransactOpts, field1, field2, field3)
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

func (_LatestValueHolder *LatestValueHolderFilterer) FilterTriggered(opts *bind.FilterOpts, field []int32) (*LatestValueHolderTriggeredIterator, error) {

	var fieldRule []interface{}
	for _, fieldItem := range field {
		fieldRule = append(fieldRule, fieldItem)
	}

	logs, sub, err := _LatestValueHolder.contract.FilterLogs(opts, "Triggered", fieldRule)
	if err != nil {
		return nil, err
	}
	return &LatestValueHolderTriggeredIterator{contract: _LatestValueHolder.contract, event: "Triggered", logs: logs, sub: sub}, nil
}

func (_LatestValueHolder *LatestValueHolderFilterer) WatchTriggered(opts *bind.WatchOpts, sink chan<- *LatestValueHolderTriggered, field []int32) (event.Subscription, error) {

	var fieldRule []interface{}
	for _, fieldItem := range field {
		fieldRule = append(fieldRule, fieldItem)
	}

	logs, sub, err := _LatestValueHolder.contract.WatchLogs(opts, "Triggered", fieldRule)
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

type LatestValueHolderTriggeredEventWithDynamicTopicIterator struct {
	Event *LatestValueHolderTriggeredEventWithDynamicTopic

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LatestValueHolderTriggeredEventWithDynamicTopicIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LatestValueHolderTriggeredEventWithDynamicTopic)
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
		it.Event = new(LatestValueHolderTriggeredEventWithDynamicTopic)
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

func (it *LatestValueHolderTriggeredEventWithDynamicTopicIterator) Error() error {
	return it.fail
}

func (it *LatestValueHolderTriggeredEventWithDynamicTopicIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LatestValueHolderTriggeredEventWithDynamicTopic struct {
	FieldHash common.Hash
	Field     string
	Raw       types.Log
}

func (_LatestValueHolder *LatestValueHolderFilterer) FilterTriggeredEventWithDynamicTopic(opts *bind.FilterOpts, fieldHash []string) (*LatestValueHolderTriggeredEventWithDynamicTopicIterator, error) {

	var fieldHashRule []interface{}
	for _, fieldHashItem := range fieldHash {
		fieldHashRule = append(fieldHashRule, fieldHashItem)
	}

	logs, sub, err := _LatestValueHolder.contract.FilterLogs(opts, "TriggeredEventWithDynamicTopic", fieldHashRule)
	if err != nil {
		return nil, err
	}
	return &LatestValueHolderTriggeredEventWithDynamicTopicIterator{contract: _LatestValueHolder.contract, event: "TriggeredEventWithDynamicTopic", logs: logs, sub: sub}, nil
}

func (_LatestValueHolder *LatestValueHolderFilterer) WatchTriggeredEventWithDynamicTopic(opts *bind.WatchOpts, sink chan<- *LatestValueHolderTriggeredEventWithDynamicTopic, fieldHash []string) (event.Subscription, error) {

	var fieldHashRule []interface{}
	for _, fieldHashItem := range fieldHash {
		fieldHashRule = append(fieldHashRule, fieldHashItem)
	}

	logs, sub, err := _LatestValueHolder.contract.WatchLogs(opts, "TriggeredEventWithDynamicTopic", fieldHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LatestValueHolderTriggeredEventWithDynamicTopic)
				if err := _LatestValueHolder.contract.UnpackLog(event, "TriggeredEventWithDynamicTopic", log); err != nil {
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

func (_LatestValueHolder *LatestValueHolderFilterer) ParseTriggeredEventWithDynamicTopic(log types.Log) (*LatestValueHolderTriggeredEventWithDynamicTopic, error) {
	event := new(LatestValueHolderTriggeredEventWithDynamicTopic)
	if err := _LatestValueHolder.contract.UnpackLog(event, "TriggeredEventWithDynamicTopic", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LatestValueHolderTriggeredWithFourTopicsIterator struct {
	Event *LatestValueHolderTriggeredWithFourTopics

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LatestValueHolderTriggeredWithFourTopicsIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LatestValueHolderTriggeredWithFourTopics)
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
		it.Event = new(LatestValueHolderTriggeredWithFourTopics)
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

func (it *LatestValueHolderTriggeredWithFourTopicsIterator) Error() error {
	return it.fail
}

func (it *LatestValueHolderTriggeredWithFourTopicsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LatestValueHolderTriggeredWithFourTopics struct {
	Field1 int32
	Field2 int32
	Field3 int32
	Raw    types.Log
}

func (_LatestValueHolder *LatestValueHolderFilterer) FilterTriggeredWithFourTopics(opts *bind.FilterOpts, field1 []int32, field2 []int32, field3 []int32) (*LatestValueHolderTriggeredWithFourTopicsIterator, error) {

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

	logs, sub, err := _LatestValueHolder.contract.FilterLogs(opts, "TriggeredWithFourTopics", field1Rule, field2Rule, field3Rule)
	if err != nil {
		return nil, err
	}
	return &LatestValueHolderTriggeredWithFourTopicsIterator{contract: _LatestValueHolder.contract, event: "TriggeredWithFourTopics", logs: logs, sub: sub}, nil
}

func (_LatestValueHolder *LatestValueHolderFilterer) WatchTriggeredWithFourTopics(opts *bind.WatchOpts, sink chan<- *LatestValueHolderTriggeredWithFourTopics, field1 []int32, field2 []int32, field3 []int32) (event.Subscription, error) {

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

	logs, sub, err := _LatestValueHolder.contract.WatchLogs(opts, "TriggeredWithFourTopics", field1Rule, field2Rule, field3Rule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LatestValueHolderTriggeredWithFourTopics)
				if err := _LatestValueHolder.contract.UnpackLog(event, "TriggeredWithFourTopics", log); err != nil {
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

func (_LatestValueHolder *LatestValueHolderFilterer) ParseTriggeredWithFourTopics(log types.Log) (*LatestValueHolderTriggeredWithFourTopics, error) {
	event := new(LatestValueHolderTriggeredWithFourTopics)
	if err := _LatestValueHolder.contract.UnpackLog(event, "TriggeredWithFourTopics", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_LatestValueHolder *LatestValueHolder) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _LatestValueHolder.abi.Events["Triggered"].ID:
		return _LatestValueHolder.ParseTriggered(log)
	case _LatestValueHolder.abi.Events["TriggeredEventWithDynamicTopic"].ID:
		return _LatestValueHolder.ParseTriggeredEventWithDynamicTopic(log)
	case _LatestValueHolder.abi.Events["TriggeredWithFourTopics"].ID:
		return _LatestValueHolder.ParseTriggeredWithFourTopics(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (LatestValueHolderTriggered) Topic() common.Hash {
	return common.HexToHash("0x7188419dcd8b51877b71766f075f3626586c0ff190e7d056aa65ce9acb649a3d")
}

func (LatestValueHolderTriggeredEventWithDynamicTopic) Topic() common.Hash {
	return common.HexToHash("0x3d969732b1bbbb9f1d7eb9f3f14e4cb50a74d950b3ef916a397b85dfbab93c67")
}

func (LatestValueHolderTriggeredWithFourTopics) Topic() common.Hash {
	return common.HexToHash("0x91c80dc390f3d041b3a04b0099b19634499541ea26972250986ee4b24a12fac5")
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

	TriggerEventWithDynamicTopic(opts *bind.TransactOpts, field string) (*types.Transaction, error)

	TriggerWithFourTopics(opts *bind.TransactOpts, field1 int32, field2 int32, field3 int32) (*types.Transaction, error)

	FilterTriggered(opts *bind.FilterOpts, field []int32) (*LatestValueHolderTriggeredIterator, error)

	WatchTriggered(opts *bind.WatchOpts, sink chan<- *LatestValueHolderTriggered, field []int32) (event.Subscription, error)

	ParseTriggered(log types.Log) (*LatestValueHolderTriggered, error)

	FilterTriggeredEventWithDynamicTopic(opts *bind.FilterOpts, fieldHash []string) (*LatestValueHolderTriggeredEventWithDynamicTopicIterator, error)

	WatchTriggeredEventWithDynamicTopic(opts *bind.WatchOpts, sink chan<- *LatestValueHolderTriggeredEventWithDynamicTopic, fieldHash []string) (event.Subscription, error)

	ParseTriggeredEventWithDynamicTopic(log types.Log) (*LatestValueHolderTriggeredEventWithDynamicTopic, error)

	FilterTriggeredWithFourTopics(opts *bind.FilterOpts, field1 []int32, field2 []int32, field3 []int32) (*LatestValueHolderTriggeredWithFourTopicsIterator, error)

	WatchTriggeredWithFourTopics(opts *bind.WatchOpts, sink chan<- *LatestValueHolderTriggeredWithFourTopics, field1 []int32, field2 []int32, field3 []int32) (event.Subscription, error)

	ParseTriggeredWithFourTopics(log types.Log) (*LatestValueHolderTriggeredWithFourTopics, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
