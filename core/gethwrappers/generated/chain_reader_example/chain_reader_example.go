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

var LatestValueHolderMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int32\",\"name\":\"field\",\"type\":\"int32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"differentField\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"oracleId\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint8[32]\",\"name\":\"oracleIds\",\"type\":\"uint8[32]\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"Account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"Accounts\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"int192\",\"name\":\"bigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"indexed\":false,\"internalType\":\"structMidLevelTestStruct\",\"name\":\"nestedStruct\",\"type\":\"tuple\"}],\"name\":\"Triggered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"fieldHash\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"field\",\"type\":\"string\"}],\"name\":\"TriggeredEventWithDynamicTopic\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int32\",\"name\":\"field1\",\"type\":\"int32\"},{\"indexed\":true,\"internalType\":\"int32\",\"name\":\"field2\",\"type\":\"int32\"},{\"indexed\":true,\"internalType\":\"int32\",\"name\":\"field3\",\"type\":\"int32\"}],\"name\":\"TriggeredWithFourTopics\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"int32\",\"name\":\"field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"differentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"oracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"oracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"int192\",\"name\":\"bigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelTestStruct\",\"name\":\"nestedStruct\",\"type\":\"tuple\"}],\"name\":\"addTestStruct\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDifferentPrimitiveValue\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"i\",\"type\":\"uint256\"}],\"name\":\"getElementAtIndex\",\"outputs\":[{\"components\":[{\"internalType\":\"int32\",\"name\":\"Field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"DifferentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"OracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"OracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"address\",\"name\":\"Account\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"Accounts\",\"type\":\"address[]\"},{\"internalType\":\"int192\",\"name\":\"BigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelTestStruct\",\"name\":\"NestedStruct\",\"type\":\"tuple\"}],\"internalType\":\"structTestStruct\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPrimitiveValue\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSliceValue\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int32\",\"name\":\"field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"differentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"oracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"oracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"int192\",\"name\":\"bigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelTestStruct\",\"name\":\"nestedStruct\",\"type\":\"tuple\"}],\"name\":\"returnSeen\",\"outputs\":[{\"components\":[{\"internalType\":\"int32\",\"name\":\"Field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"DifferentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"OracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"OracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"address\",\"name\":\"Account\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"Accounts\",\"type\":\"address[]\"},{\"internalType\":\"int192\",\"name\":\"BigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelTestStruct\",\"name\":\"NestedStruct\",\"type\":\"tuple\"}],\"internalType\":\"structTestStruct\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int32\",\"name\":\"field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"differentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"oracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"oracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"int192\",\"name\":\"bigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"IntVal\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelTestStruct\",\"name\":\"nestedStruct\",\"type\":\"tuple\"}],\"name\":\"triggerEvent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"field\",\"type\":\"string\"}],\"name\":\"triggerEventWithDynamicTopic\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int32\",\"name\":\"field1\",\"type\":\"int32\"},{\"internalType\":\"int32\",\"name\":\"field2\",\"type\":\"int32\"},{\"internalType\":\"int32\",\"name\":\"field3\",\"type\":\"int32\"}],\"name\":\"triggerWithFourTopics\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50600180548082018255600082905260048082047fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6908101805460086003958616810261010090810a8088026001600160401b0391820219909416939093179093558654808801909755848704909301805496909516909202900a91820291021990921691909117905561176c806100a96000396000f3fe608060405234801561001057600080fd5b50600436106100a35760003560e01c80637f002d6711610076578063dbfd73321161005b578063dbfd73321461013e578063ef4e1ced14610151578063f6f871c81461015857600080fd5b80637f002d671461010e578063ab5e0b381461012157600080fd5b80632c45576f146100a85780633272b66c146100d157806349eac2ac146100e6578063679004a4146100f9575b600080fd5b6100bb6100b6366004610baa565b61016b565b6040516100c89190610d09565b60405180910390f35b6100e46100df366004610e48565b610446565b005b6100e46100f4366004610f5d565b61049b565b61010161079e565b6040516100c8919061104f565b6100e461011c366004610f5d565b61082a565b6107c65b60405167ffffffffffffffff90911681526020016100c8565b6100e461014c36600461109d565b610881565b6003610125565b6100bb610166366004610f5d565b6108be565b6101736109c7565b60006101806001846110e0565b815481106101905761019061111a565b6000918252602091829020604080516101008101909152600a90920201805460030b825260018101805492939192918401916101cb90611149565b80601f01602080910402602001604051908101604052809291908181526020018280546101f790611149565b80156102445780601f1061021957610100808354040283529160200191610244565b820191906000526020600020905b81548152906001019060200180831161022757829003601f168201915b5050509183525050600282015460ff166020808301919091526040805161040081018083529190930192916003850191826000855b825461010083900a900460ff1681526020600192830181810494850194909303909202910180841161027957505050928452505050600482015473ffffffffffffffffffffffffffffffffffffffff16602080830191909152600583018054604080518285028101850182528281529401939283018282801561033257602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610307575b5050509183525050600682015460170b6020808301919091526040805180820182526007808601805460f01b7fffff0000000000000000000000000000000000000000000000000000000000001683528351808501855260088801805490930b815260098801805495909701969395919486830194919392840191906103b790611149565b80601f01602080910402602001604051908101604052809291908181526020018280546103e390611149565b80156104305780601f1061040557610100808354040283529160200191610430565b820191906000526020600020905b81548152906001019060200180831161041357829003601f168201915b5050509190925250505090525090525092915050565b8181604051610456929190611196565b60405180910390207f3d969732b1bbbb9f1d7eb9f3f14e4cb50a74d950b3ef916a397b85dfbab93c67838360405161048f9291906111ef565b60405180910390a25050565b60006040518061010001604052808c60030b81526020018b8b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050509082525060ff8a166020808301919091526040805161040081810183529190930192918b9183908390808284376000920191909152505050815273ffffffffffffffffffffffffffffffffffffffff8816602080830191909152604080518883028181018401835289825291909301929189918991829190850190849080828437600092019190915250505090825250601785900b602082015260400161058d846112ec565b905281546001808201845560009384526020938490208351600a9093020180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff9093169290921782559282015191929091908201906105f39082611446565b5060408201516002820180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff90921691909117905560608201516106419060038301906020610a16565b5060808201516004820180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff90921691909117905560a082015180516106a8916005840191602090910190610aa9565b5060c08201516006820180547fffffffffffffffff0000000000000000000000000000000000000000000000001677ffffffffffffffffffffffffffffffffffffffffffffffff90921691909117905560e082015180516007830180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001660f09290921c91909117815560208083015180516008860180547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001667ffffffffffffffff90921691909117815591810151909190600986019061078b9082611446565b5050505050505050505050505050505050565b6060600180548060200260200160405190810160405280929190818152602001828054801561082057602002820191906000526020600020906000905b82829054906101000a900467ffffffffffffffff1667ffffffffffffffff16815260200190600801906020826007010492830192600103820291508084116107db5790505b5050505050905090565b8960030b7f7188419dcd8b51877b71766f075f3626586c0ff190e7d056aa65ce9acb649a3d8a8a8a8a8a8a8a8a8a60405161086d999897969594939291906116a5565b60405180910390a250505050505050505050565b8060030b8260030b8460030b7f91c80dc390f3d041b3a04b0099b19634499541ea26972250986ee4b24a12fac560405160405180910390a4505050565b6108c66109c7565b6040518061010001604052808c60030b81526020018b8b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050509082525060ff8a166020808301919091526040805161040081810183529190930192918b9183908390808284376000920191909152505050815273ffffffffffffffffffffffffffffffffffffffff8816602080830191909152604080518883028181018401835289825291909301929189918991829190850190849080828437600092019190915250505090825250601785900b60208201526040016109b6846112ec565b90529b9a5050505050505050505050565b60408051610100810182526000808252606060208301819052928201529081016109ef610b23565b8152600060208201819052606060408301819052820152608001610a11610b42565b905290565b600183019183908215610a995791602002820160005b83821115610a6a57835183826101000a81548160ff021916908360ff1602179055509260200192600101602081600001049283019260010302610a2c565b8015610a975782816101000a81549060ff0219169055600101602081600001049283019260010302610a6a565b505b50610aa5929150610b95565b5090565b828054828255906000526020600020908101928215610a99579160200282015b82811115610a9957825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190610ac9565b6040518061040001604052806020906020820280368337509192915050565b604051806040016040528060007dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168152602001610a116040518060400160405280600060070b8152602001606081525090565b5b80821115610aa55760008155600101610b96565b600060208284031215610bbc57600080fd5b5035919050565b6000815180845260005b81811015610be957602081850181015186830182015201610bcd565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b8060005b6020808210610c3a5750610c51565b825160ff1685529384019390910190600101610c2b565b50505050565b600081518084526020808501945080840160005b83811015610c9d57815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101610c6b565b509495945050505050565b7fffff00000000000000000000000000000000000000000000000000000000000081511682526000602082015160406020850152805160070b60408501526020810151905060406060850152610d016080850182610bc3565b949350505050565b60208152610d1d60208201835160030b9052565b600060208301516104e0806040850152610d3b610500850183610bc3565b91506040850151610d51606086018260ff169052565b506060850151610d646080860182610c27565b50608085015173ffffffffffffffffffffffffffffffffffffffff1661048085015260a08501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe085840381016104a0870152610dc18483610c57565b935060c08701519150610dda6104c087018360170b9052565b60e0870151915080868503018387015250610df58382610ca8565b9695505050505050565b60008083601f840112610e1157600080fd5b50813567ffffffffffffffff811115610e2957600080fd5b602083019150836020828501011115610e4157600080fd5b9250929050565b60008060208385031215610e5b57600080fd5b823567ffffffffffffffff811115610e7257600080fd5b610e7e85828601610dff565b90969095509350505050565b8035600381900b8114610e9c57600080fd5b919050565b803560ff81168114610e9c57600080fd5b806104008101831015610ec457600080fd5b92915050565b803573ffffffffffffffffffffffffffffffffffffffff81168114610e9c57600080fd5b60008083601f840112610f0057600080fd5b50813567ffffffffffffffff811115610f1857600080fd5b6020830191508360208260051b8501011115610e4157600080fd5b8035601781900b8114610e9c57600080fd5b600060408284031215610f5757600080fd5b50919050565b6000806000806000806000806000806104e08b8d031215610f7d57600080fd5b610f868b610e8a565b995060208b013567ffffffffffffffff80821115610fa357600080fd5b610faf8e838f01610dff565b909b509950899150610fc360408e01610ea1565b9850610fd28e60608f01610eb2565b9750610fe16104608e01610eca565b96506104808d0135915080821115610ff857600080fd5b6110048e838f01610eee565b90965094508491506110196104a08e01610f33565b93506104c08d013591508082111561103057600080fd5b5061103d8d828e01610f45565b9150509295989b9194979a5092959850565b6020808252825182820181905260009190848201906040850190845b8181101561109157835167ffffffffffffffff168352928401929184019160010161106b565b50909695505050505050565b6000806000606084860312156110b257600080fd5b6110bb84610e8a565b92506110c960208501610e8a565b91506110d760408501610e8a565b90509250925092565b81810381811115610ec4577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600181811c9082168061115d57607f821691505b602082108103610f57577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b8183823760009101908152919050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b602081526000610d016020830184866111a6565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040805190810167ffffffffffffffff8111828210171561125557611255611203565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156112a2576112a2611203565b604052919050565b80357fffff00000000000000000000000000000000000000000000000000000000000081168114610e9c57600080fd5b8035600781900b8114610e9c57600080fd5b6000604082360312156112fe57600080fd5b611306611232565b61130f836112aa565b815260208084013567ffffffffffffffff8082111561132d57600080fd5b81860191506040823603121561134257600080fd5b61134a611232565b611353836112da565b8152838301358281111561136657600080fd5b929092019136601f84011261137a57600080fd5b82358281111561138c5761138c611203565b6113bc857fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8401160161125b565b925080835236858286010111156113d257600080fd5b8085850186850137600090830185015280840191909152918301919091525092915050565b601f82111561144157600081815260208120601f850160051c8101602086101561141e5750805b601f850160051c820191505b8181101561143d5782815560010161142a565b5050505b505050565b815167ffffffffffffffff81111561146057611460611203565b6114748161146e8454611149565b846113f7565b602080601f8311600181146114c757600084156114915750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b17855561143d565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015611514578886015182559484019460019091019084016114f5565b508582101561155057878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b8183526000602080850194508260005b85811015610c9d5773ffffffffffffffffffffffffffffffffffffffff61159683610eca565b1687529582019590820190600101611570565b7fffff0000000000000000000000000000000000000000000000000000000000006115d3826112aa565b168252600060208201357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc183360301811261160d57600080fd5b60406020850152820161161f816112da565b60070b604085015260208101357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe182360301811261165c57600080fd5b0160208101903567ffffffffffffffff81111561167857600080fd5b80360382131561168757600080fd5b6040606086015261169c6080860182846111a6565b95945050505050565b60006104c08083526116ba8184018c8e6111a6565b9050602060ff808c1682860152604085018b60005b848110156116f457836116e183610ea1565b16835291840191908401906001016116cf565b505050505073ffffffffffffffffffffffffffffffffffffffff881661044084015282810361046084015261172a818789611560565b905061173c61048084018660170b9052565b8281036104a084015261174f81856115a9565b9c9b50505050505050505050505056fea164736f6c6343000813000a",
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
	err := _LatestValueHolder.contract.Call(opts, &out, "getDifferentPrimitiveValue")

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
	err := _LatestValueHolder.contract.Call(opts, &out, "getElementAtIndex", i)

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
	err := _LatestValueHolder.contract.Call(opts, &out, "getPrimitiveValue")

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
	err := _LatestValueHolder.contract.Call(opts, &out, "getSliceValue")

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
	err := _LatestValueHolder.contract.Call(opts, &out, "returnSeen", field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)

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
	return _LatestValueHolder.contract.Transact(opts, "addTestStruct", field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)
}

func (_LatestValueHolder *LatestValueHolderSession) AddTestStruct(field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedStruct MidLevelTestStruct) (*types.Transaction, error) {
	return _LatestValueHolder.Contract.AddTestStruct(&_LatestValueHolder.TransactOpts, field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)
}

func (_LatestValueHolder *LatestValueHolderTransactorSession) AddTestStruct(field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedStruct MidLevelTestStruct) (*types.Transaction, error) {
	return _LatestValueHolder.Contract.AddTestStruct(&_LatestValueHolder.TransactOpts, field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)
}

func (_LatestValueHolder *LatestValueHolderTransactor) TriggerEvent(opts *bind.TransactOpts, field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedStruct MidLevelTestStruct) (*types.Transaction, error) {
	return _LatestValueHolder.contract.Transact(opts, "triggerEvent", field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)
}

func (_LatestValueHolder *LatestValueHolderSession) TriggerEvent(field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedStruct MidLevelTestStruct) (*types.Transaction, error) {
	return _LatestValueHolder.Contract.TriggerEvent(&_LatestValueHolder.TransactOpts, field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)
}

func (_LatestValueHolder *LatestValueHolderTransactorSession) TriggerEvent(field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedStruct MidLevelTestStruct) (*types.Transaction, error) {
	return _LatestValueHolder.Contract.TriggerEvent(&_LatestValueHolder.TransactOpts, field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)
}

func (_LatestValueHolder *LatestValueHolderTransactor) TriggerEventWithDynamicTopic(opts *bind.TransactOpts, field string) (*types.Transaction, error) {
	return _LatestValueHolder.contract.Transact(opts, "triggerEventWithDynamicTopic", field)
}

func (_LatestValueHolder *LatestValueHolderSession) TriggerEventWithDynamicTopic(field string) (*types.Transaction, error) {
	return _LatestValueHolder.Contract.TriggerEventWithDynamicTopic(&_LatestValueHolder.TransactOpts, field)
}

func (_LatestValueHolder *LatestValueHolderTransactorSession) TriggerEventWithDynamicTopic(field string) (*types.Transaction, error) {
	return _LatestValueHolder.Contract.TriggerEventWithDynamicTopic(&_LatestValueHolder.TransactOpts, field)
}

func (_LatestValueHolder *LatestValueHolderTransactor) TriggerWithFourTopics(opts *bind.TransactOpts, field1 int32, field2 int32, field3 int32) (*types.Transaction, error) {
	return _LatestValueHolder.contract.Transact(opts, "triggerWithFourTopics", field1, field2, field3)
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
