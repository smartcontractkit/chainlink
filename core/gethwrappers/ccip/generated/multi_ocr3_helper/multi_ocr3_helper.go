// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package multi_ocr3_helper

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

type MultiOCR3BaseConfigInfo struct {
	ConfigDigest                   [32]byte
	F                              uint8
	N                              uint8
	IsSignatureVerificationEnabled bool
}

type MultiOCR3BaseOCRConfig struct {
	ConfigInfo   MultiOCR3BaseConfigInfo
	Signers      []common.Address
	Transmitters []common.Address
}

type MultiOCR3BaseOCRConfigArgs struct {
	ConfigDigest                   [32]byte
	OcrPluginType                  uint8
	F                              uint8
	IsSignatureVerificationEnabled bool
	Signers                        []common.Address
	Transmitters                   []common.Address
}

type MultiOCR3BaseOracle struct {
	Index uint8
	Role  uint8
}

var MultiOCR3HelperMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"expected\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"actual\",\"type\":\"bytes32\"}],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"ForkedChain\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"enumMultiOCR3Base.InvalidConfigErrorType\",\"name\":\"errorType\",\"type\":\"uint8\"}],\"name\":\"InvalidConfig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NonUniqueSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OracleCannotBeZeroAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SignaturesOutOfRegistration\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"}],\"name\":\"StaticConfigCannotBeChanged\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedTransmitter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"WrongMessageLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WrongNumberOfSignatures\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"}],\"name\":\"AfterConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"F\",\"type\":\"uint8\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"oracleAddress\",\"type\":\"address\"}],\"name\":\"getOracle\",\"outputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"index\",\"type\":\"uint8\"},{\"internalType\":\"enumMultiOCR3Base.Role\",\"name\":\"role\",\"type\":\"uint8\"}],\"internalType\":\"structMultiOCR3Base.Oracle\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"}],\"name\":\"latestConfigDetails\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"F\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"n\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"isSignatureVerificationEnabled\",\"type\":\"bool\"}],\"internalType\":\"structMultiOCR3Base.ConfigInfo\",\"name\":\"configInfo\",\"type\":\"tuple\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"}],\"internalType\":\"structMultiOCR3Base.OCRConfig\",\"name\":\"ocrConfig\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"F\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"isSignatureVerificationEnabled\",\"type\":\"bool\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"}],\"internalType\":\"structMultiOCR3Base.OCRConfigArgs[]\",\"name\":\"ocrConfigArgs\",\"type\":\"tuple[]\"}],\"name\":\"setOCR3Configs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"}],\"name\":\"setTransmitOcrPluginType\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmitWithSignatures\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"}],\"name\":\"transmitWithoutSignatures\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b503380600081620000695760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156200009c576200009c81620000a9565b5050466080525062000154565b336001600160a01b03821603620001035760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000060565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b608051611db56200017760003960008181610f6d0152610fb90152611db56000f3fe608060405234801561001057600080fd5b50600436106100be5760003560e01c80637ac0aa1a11610076578063c673e5841161005b578063c673e584146101c5578063f2fde38b146101e5578063f716f99f146101f857600080fd5b80637ac0aa1a1461015b5780638da5cb5b1461019d57600080fd5b806334a9c92e116100a757806334a9c92e1461012057806344e65e551461014057806379ba50971461015357600080fd5b8063181f5a77146100c357806326bf9d261461010b575b600080fd5b604080518082018252601981527f4d756c74694f4352334261736548656c70657220312e302e300000000000000060208201529051610102919061155a565b60405180910390f35b61011e610119366004611621565b61020b565b005b61013361012e3660046116af565b61023a565b6040516101029190611711565b61011e61014e366004611784565b6102ca565b61011e61034d565b61011e610169366004611837565b600480547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff92909216919091179055565b60005460405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610102565b6101d86101d3366004611837565b61044f565b60405161010291906118ab565b61011e6101f336600461193e565b6105c7565b61011e610206366004611aaa565b6105db565b604080516000808252602082019092526004549091506102349060ff168585858580600061061d565b50505050565b6040805180820182526000808252602080830182905260ff86811683526003825284832073ffffffffffffffffffffffffffffffffffffffff871684528252918490208451808601909552805480841686529394939092918401916101009091041660028111156102ad576102ad6116e2565b60028111156102be576102be6116e2565b90525090505b92915050565b60045460408051602080880282810182019093528782526103439360ff16928c928c928c928c918c91829185019084908082843760009201919091525050604080516020808d0282810182019093528c82529093508c92508b9182918501908490808284376000920191909152508a925061061d915050565b5050505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146103d3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6104926040805160e081019091526000606082018181526080830182905260a0830182905260c08301919091528190815260200160608152602001606081525090565b60ff808316600090815260026020818152604092839020835160e081018552815460608201908152600183015480881660808401526101008104881660a0840152620100009004909616151560c08201529485529182018054845181840281018401909552808552929385830193909283018282801561054857602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff16815260019091019060200180831161051d575b50505050508152602001600382018054806020026020016040519081016040528092919081815260200182805480156105b757602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff16815260019091019060200180831161058c575b5050505050815250509050919050565b6105cf6109a1565b6105d881610a24565b50565b6105e36109a1565b60005b81518110156106195761061182828151811061060457610604611c13565b6020026020010151610b19565b6001016105e6565b5050565b60ff8781166000908152600260209081526040808320815160808101835281548152600190910154808616938201939093526101008304851691810191909152620100009091049092161515606083015287359061067c8760a4611c71565b90508260600151156106c4578451610695906020611c84565b86516106a2906020611c84565b6106ad9060a0611c71565b6106b79190611c71565b6106c19082611c71565b90505b368114610706576040517f8e1192e1000000000000000000000000000000000000000000000000000000008152600481018290523660248201526044016103ca565b508151811461074e5781516040517f93df584c0000000000000000000000000000000000000000000000000000000081526004810191909152602481018290526044016103ca565b610756610f6a565b60ff808a16600090815260036020908152604080832033845282528083208151808301909252805480861683529394919390928401916101009091041660028111156107a4576107a46116e2565b60028111156107b5576107b56116e2565b90525090506002816020015160028111156107d2576107d26116e2565b1480156108335750600260008b60ff1660ff168152602001908152602001600020600301816000015160ff168154811061080e5761080e611c13565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff1633145b610869576040517fda0f08e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5081606001511561094b576020820151610884906001611c9b565b60ff168551146108c0576040517f71253a2500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b83518551146108fb576040517fa75d88af00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000878760405161090d929190611cb4565b604051908190038120610924918b90602001611cc4565b6040516020818303038152906040528051906020012090506109498a82888888610feb565b505b6040805182815260208a81013567ffffffffffffffff169082015260ff8b16917f198d6990ef96613a9026203077e422916918b03ff47f0be6bee7b02d8e139ef0910160405180910390a2505050505050505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610a22576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016103ca565b565b3373ffffffffffffffffffffffffffffffffffffffff821603610aa3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016103ca565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b806040015160ff16600003610b5d5760006040517f367f56a20000000000000000000000000000000000000000000000000000000081526004016103ca9190611cd8565b60208082015160ff80821660009081526002909352604083206001810154929390928392169003610bca57606084015160018201805491151562010000027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffff909216919091179055610c1f565b6060840151600182015460ff6201000090910416151590151514610c1f576040517f87f6037c00000000000000000000000000000000000000000000000000000000815260ff841660048201526024016103ca565b60a084015180516101001015610c645760016040517f367f56a20000000000000000000000000000000000000000000000000000000081526004016103ca9190611cd8565b8051600003610ca25760056040517f367f56a20000000000000000000000000000000000000000000000000000000081526004016103ca9190611cd8565b610d158484600301805480602002602001604051908101604052809291908181526020018280548015610d0b57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610ce0575b50505050506111fb565b846060015115610eba57610d908484600201805480602002602001604051908101604052809291908181526020018280548015610d0b5760200282019190600052602060002090815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610ce05750505050506111fb565b608085015180516101001015610dd55760026040517f367f56a20000000000000000000000000000000000000000000000000000000081526004016103ca9190611cd8565b6040860151610de5906003611cf2565b60ff16815111610e245760036040517f367f56a20000000000000000000000000000000000000000000000000000000081526004016103ca9190611cd8565b815181511015610e635760016040517f367f56a20000000000000000000000000000000000000000000000000000000081526004016103ca9190611cd8565b80516001840180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1661010060ff841602179055610eab90600286019060208401906114bb565b50610eb885826001611293565b505b610ec684826002611293565b8051610edb90600385019060208401906114bb565b506040858101516001840180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff8316179055865180855560a088015192517fab8b1b57514019638d7b5ce9c638fe71366fe8e2be1c40a7a80f1733d0e9f54793610f529389939260028a01929190611d15565b60405180910390a1610f6384611482565b5050505050565b467f000000000000000000000000000000000000000000000000000000000000000014610a22576040517f0f01ce850000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000000060048201524660248201526044016103ca565b8251600090815b8181101561034357600060018886846020811061101157611011611c13565b61101e91901a601b611c9b565b89858151811061103057611030611c13565b602002602001015189868151811061104a5761104a611c13565b602002602001015160405160008152602001604052604051611088949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa1580156110aa573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015160ff808e1660009081526003602090815285822073ffffffffffffffffffffffffffffffffffffffff851683528152858220858701909652855480841686529397509095509293928401916101009004166002811115611136576111366116e2565b6002811115611147576111476116e2565b9052509050600181602001516002811115611164576111646116e2565b1461119b576040517fca31867a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8051600160ff9091161b8516156111de576040517ff67bc7c400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b806000015160ff166001901b851794505050806001019050610ff2565b60005b815181101561128e5760ff83166000908152600360205260408120835190919084908490811061123057611230611c13565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff16825281019190915260400160002080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001690556001016111fe565b505050565b60005b82518110156102345760008382815181106112b3576112b3611c13565b60200260200101519050600060028111156112d0576112d06116e2565b60ff808716600090815260036020908152604080832073ffffffffffffffffffffffffffffffffffffffff87168452909152902054610100900416600281111561131c5761131c6116e2565b146113565760046040517f367f56a20000000000000000000000000000000000000000000000000000000081526004016103ca9190611cd8565b73ffffffffffffffffffffffffffffffffffffffff81166113a3576040517fd6c62c9b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60405180604001604052808360ff1681526020018460028111156113c9576113c96116e2565b905260ff808716600090815260036020908152604080832073ffffffffffffffffffffffffffffffffffffffff8716845282529091208351815493167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00841681178255918401519092909183917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000161761010083600281111561146e5761146e6116e2565b021790555090505050806001019050611296565b60405160ff821681527f897ac1b2c12867721b284f3eb147bd4ab046d4eef1cf31c1d8988bfcfb962b539060200160405180910390a150565b828054828255906000526020600020908101928215611535579160200282015b8281111561153557825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff9091161782556020909201916001909101906114db565b50611541929150611545565b5090565b5b808211156115415760008155600101611546565b60006020808352835180602085015260005b818110156115885785810183015185820160400152820161156c565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b80606081018310156102c457600080fd5b60008083601f8401126115ea57600080fd5b50813567ffffffffffffffff81111561160257600080fd5b60208301915083602082850101111561161a57600080fd5b9250929050565b60008060006080848603121561163657600080fd5b61164085856115c7565b9250606084013567ffffffffffffffff81111561165c57600080fd5b611668868287016115d8565b9497909650939450505050565b803560ff8116811461168657600080fd5b919050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461168657600080fd5b600080604083850312156116c257600080fd5b6116cb83611675565b91506116d96020840161168b565b90509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b815160ff1681526020820151604082019060038110611732576117326116e2565b8060208401525092915050565b60008083601f84011261175157600080fd5b50813567ffffffffffffffff81111561176957600080fd5b6020830191508360208260051b850101111561161a57600080fd5b60008060008060008060008060e0898b0312156117a057600080fd5b6117aa8a8a6115c7565b9750606089013567ffffffffffffffff808211156117c757600080fd5b6117d38c838d016115d8565b909950975060808b01359150808211156117ec57600080fd5b6117f88c838d0161173f565b909750955060a08b013591508082111561181157600080fd5b5061181e8b828c0161173f565b999c989b50969995989497949560c00135949350505050565b60006020828403121561184957600080fd5b61185282611675565b9392505050565b60008151808452602080850194506020840160005b838110156118a057815173ffffffffffffffffffffffffffffffffffffffff168752958201959082019060010161186e565b509495945050505050565b60208152600082518051602084015260ff602082015116604084015260ff604082015116606084015260608101511515608084015250602083015160c060a08401526118fa60e0840182611859565b905060408401517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08483030160c08501526119358282611859565b95945050505050565b60006020828403121561195057600080fd5b6118528261168b565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160c0810167ffffffffffffffff811182821017156119ab576119ab611959565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156119f8576119f8611959565b604052919050565b600067ffffffffffffffff821115611a1a57611a1a611959565b5060051b60200190565b8035801515811461168657600080fd5b600082601f830112611a4557600080fd5b81356020611a5a611a5583611a00565b6119b1565b8083825260208201915060208460051b870101935086841115611a7c57600080fd5b602086015b84811015611a9f57611a928161168b565b8352918301918301611a81565b509695505050505050565b60006020808385031215611abd57600080fd5b823567ffffffffffffffff80821115611ad557600080fd5b818501915085601f830112611ae957600080fd5b8135611af7611a5582611a00565b81815260059190911b83018401908481019088831115611b1657600080fd5b8585015b83811015611c0657803585811115611b3157600080fd5b860160c0818c037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0011215611b665760008081fd5b611b6e611988565b8882013581526040611b81818401611675565b8a8301526060611b92818501611675565b8284015260809150611ba5828501611a24565b9083015260a08381013589811115611bbd5760008081fd5b611bcb8f8d83880101611a34565b838501525060c0840135915088821115611be55760008081fd5b611bf38e8c84870101611a34565b9083015250845250918601918601611b1a565b5098975050505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b808201808211156102c4576102c4611c42565b80820281158282048414176102c4576102c4611c42565b60ff81811683821601908111156102c4576102c4611c42565b8183823760009101908152919050565b828152606082602083013760800192915050565b6020810160068310611cec57611cec6116e2565b91905290565b60ff8181168382160290811690818114611d0e57611d0e611c42565b5092915050565b600060a0820160ff88168352602087602085015260a0604085015281875480845260c086019150886000526020600020935060005b81811015611d7c57845473ffffffffffffffffffffffffffffffffffffffff1683526001948501949284019201611d4a565b50508481036060860152611d908188611859565b935050505060ff83166080830152969550505050505056fea164736f6c6343000818000a",
}

var MultiOCR3HelperABI = MultiOCR3HelperMetaData.ABI

var MultiOCR3HelperBin = MultiOCR3HelperMetaData.Bin

func DeployMultiOCR3Helper(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *MultiOCR3Helper, error) {
	parsed, err := MultiOCR3HelperMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MultiOCR3HelperBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MultiOCR3Helper{address: address, abi: *parsed, MultiOCR3HelperCaller: MultiOCR3HelperCaller{contract: contract}, MultiOCR3HelperTransactor: MultiOCR3HelperTransactor{contract: contract}, MultiOCR3HelperFilterer: MultiOCR3HelperFilterer{contract: contract}}, nil
}

type MultiOCR3Helper struct {
	address common.Address
	abi     abi.ABI
	MultiOCR3HelperCaller
	MultiOCR3HelperTransactor
	MultiOCR3HelperFilterer
}

type MultiOCR3HelperCaller struct {
	contract *bind.BoundContract
}

type MultiOCR3HelperTransactor struct {
	contract *bind.BoundContract
}

type MultiOCR3HelperFilterer struct {
	contract *bind.BoundContract
}

type MultiOCR3HelperSession struct {
	Contract     *MultiOCR3Helper
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type MultiOCR3HelperCallerSession struct {
	Contract *MultiOCR3HelperCaller
	CallOpts bind.CallOpts
}

type MultiOCR3HelperTransactorSession struct {
	Contract     *MultiOCR3HelperTransactor
	TransactOpts bind.TransactOpts
}

type MultiOCR3HelperRaw struct {
	Contract *MultiOCR3Helper
}

type MultiOCR3HelperCallerRaw struct {
	Contract *MultiOCR3HelperCaller
}

type MultiOCR3HelperTransactorRaw struct {
	Contract *MultiOCR3HelperTransactor
}

func NewMultiOCR3Helper(address common.Address, backend bind.ContractBackend) (*MultiOCR3Helper, error) {
	abi, err := abi.JSON(strings.NewReader(MultiOCR3HelperABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindMultiOCR3Helper(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MultiOCR3Helper{address: address, abi: abi, MultiOCR3HelperCaller: MultiOCR3HelperCaller{contract: contract}, MultiOCR3HelperTransactor: MultiOCR3HelperTransactor{contract: contract}, MultiOCR3HelperFilterer: MultiOCR3HelperFilterer{contract: contract}}, nil
}

func NewMultiOCR3HelperCaller(address common.Address, caller bind.ContractCaller) (*MultiOCR3HelperCaller, error) {
	contract, err := bindMultiOCR3Helper(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MultiOCR3HelperCaller{contract: contract}, nil
}

func NewMultiOCR3HelperTransactor(address common.Address, transactor bind.ContractTransactor) (*MultiOCR3HelperTransactor, error) {
	contract, err := bindMultiOCR3Helper(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MultiOCR3HelperTransactor{contract: contract}, nil
}

func NewMultiOCR3HelperFilterer(address common.Address, filterer bind.ContractFilterer) (*MultiOCR3HelperFilterer, error) {
	contract, err := bindMultiOCR3Helper(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MultiOCR3HelperFilterer{contract: contract}, nil
}

func bindMultiOCR3Helper(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MultiOCR3HelperMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_MultiOCR3Helper *MultiOCR3HelperRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MultiOCR3Helper.Contract.MultiOCR3HelperCaller.contract.Call(opts, result, method, params...)
}

func (_MultiOCR3Helper *MultiOCR3HelperRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultiOCR3Helper.Contract.MultiOCR3HelperTransactor.contract.Transfer(opts)
}

func (_MultiOCR3Helper *MultiOCR3HelperRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MultiOCR3Helper.Contract.MultiOCR3HelperTransactor.contract.Transact(opts, method, params...)
}

func (_MultiOCR3Helper *MultiOCR3HelperCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MultiOCR3Helper.Contract.contract.Call(opts, result, method, params...)
}

func (_MultiOCR3Helper *MultiOCR3HelperTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultiOCR3Helper.Contract.contract.Transfer(opts)
}

func (_MultiOCR3Helper *MultiOCR3HelperTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MultiOCR3Helper.Contract.contract.Transact(opts, method, params...)
}

func (_MultiOCR3Helper *MultiOCR3HelperCaller) GetOracle(opts *bind.CallOpts, ocrPluginType uint8, oracleAddress common.Address) (MultiOCR3BaseOracle, error) {
	var out []interface{}
	err := _MultiOCR3Helper.contract.Call(opts, &out, "getOracle", ocrPluginType, oracleAddress)

	if err != nil {
		return *new(MultiOCR3BaseOracle), err
	}

	out0 := *abi.ConvertType(out[0], new(MultiOCR3BaseOracle)).(*MultiOCR3BaseOracle)

	return out0, err

}

func (_MultiOCR3Helper *MultiOCR3HelperSession) GetOracle(ocrPluginType uint8, oracleAddress common.Address) (MultiOCR3BaseOracle, error) {
	return _MultiOCR3Helper.Contract.GetOracle(&_MultiOCR3Helper.CallOpts, ocrPluginType, oracleAddress)
}

func (_MultiOCR3Helper *MultiOCR3HelperCallerSession) GetOracle(ocrPluginType uint8, oracleAddress common.Address) (MultiOCR3BaseOracle, error) {
	return _MultiOCR3Helper.Contract.GetOracle(&_MultiOCR3Helper.CallOpts, ocrPluginType, oracleAddress)
}

func (_MultiOCR3Helper *MultiOCR3HelperCaller) LatestConfigDetails(opts *bind.CallOpts, ocrPluginType uint8) (MultiOCR3BaseOCRConfig, error) {
	var out []interface{}
	err := _MultiOCR3Helper.contract.Call(opts, &out, "latestConfigDetails", ocrPluginType)

	if err != nil {
		return *new(MultiOCR3BaseOCRConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(MultiOCR3BaseOCRConfig)).(*MultiOCR3BaseOCRConfig)

	return out0, err

}

func (_MultiOCR3Helper *MultiOCR3HelperSession) LatestConfigDetails(ocrPluginType uint8) (MultiOCR3BaseOCRConfig, error) {
	return _MultiOCR3Helper.Contract.LatestConfigDetails(&_MultiOCR3Helper.CallOpts, ocrPluginType)
}

func (_MultiOCR3Helper *MultiOCR3HelperCallerSession) LatestConfigDetails(ocrPluginType uint8) (MultiOCR3BaseOCRConfig, error) {
	return _MultiOCR3Helper.Contract.LatestConfigDetails(&_MultiOCR3Helper.CallOpts, ocrPluginType)
}

func (_MultiOCR3Helper *MultiOCR3HelperCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MultiOCR3Helper.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_MultiOCR3Helper *MultiOCR3HelperSession) Owner() (common.Address, error) {
	return _MultiOCR3Helper.Contract.Owner(&_MultiOCR3Helper.CallOpts)
}

func (_MultiOCR3Helper *MultiOCR3HelperCallerSession) Owner() (common.Address, error) {
	return _MultiOCR3Helper.Contract.Owner(&_MultiOCR3Helper.CallOpts)
}

func (_MultiOCR3Helper *MultiOCR3HelperCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MultiOCR3Helper.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_MultiOCR3Helper *MultiOCR3HelperSession) TypeAndVersion() (string, error) {
	return _MultiOCR3Helper.Contract.TypeAndVersion(&_MultiOCR3Helper.CallOpts)
}

func (_MultiOCR3Helper *MultiOCR3HelperCallerSession) TypeAndVersion() (string, error) {
	return _MultiOCR3Helper.Contract.TypeAndVersion(&_MultiOCR3Helper.CallOpts)
}

func (_MultiOCR3Helper *MultiOCR3HelperTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultiOCR3Helper.contract.Transact(opts, "acceptOwnership")
}

func (_MultiOCR3Helper *MultiOCR3HelperSession) AcceptOwnership() (*types.Transaction, error) {
	return _MultiOCR3Helper.Contract.AcceptOwnership(&_MultiOCR3Helper.TransactOpts)
}

func (_MultiOCR3Helper *MultiOCR3HelperTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _MultiOCR3Helper.Contract.AcceptOwnership(&_MultiOCR3Helper.TransactOpts)
}

func (_MultiOCR3Helper *MultiOCR3HelperTransactor) SetOCR3Configs(opts *bind.TransactOpts, ocrConfigArgs []MultiOCR3BaseOCRConfigArgs) (*types.Transaction, error) {
	return _MultiOCR3Helper.contract.Transact(opts, "setOCR3Configs", ocrConfigArgs)
}

func (_MultiOCR3Helper *MultiOCR3HelperSession) SetOCR3Configs(ocrConfigArgs []MultiOCR3BaseOCRConfigArgs) (*types.Transaction, error) {
	return _MultiOCR3Helper.Contract.SetOCR3Configs(&_MultiOCR3Helper.TransactOpts, ocrConfigArgs)
}

func (_MultiOCR3Helper *MultiOCR3HelperTransactorSession) SetOCR3Configs(ocrConfigArgs []MultiOCR3BaseOCRConfigArgs) (*types.Transaction, error) {
	return _MultiOCR3Helper.Contract.SetOCR3Configs(&_MultiOCR3Helper.TransactOpts, ocrConfigArgs)
}

func (_MultiOCR3Helper *MultiOCR3HelperTransactor) SetTransmitOcrPluginType(opts *bind.TransactOpts, ocrPluginType uint8) (*types.Transaction, error) {
	return _MultiOCR3Helper.contract.Transact(opts, "setTransmitOcrPluginType", ocrPluginType)
}

func (_MultiOCR3Helper *MultiOCR3HelperSession) SetTransmitOcrPluginType(ocrPluginType uint8) (*types.Transaction, error) {
	return _MultiOCR3Helper.Contract.SetTransmitOcrPluginType(&_MultiOCR3Helper.TransactOpts, ocrPluginType)
}

func (_MultiOCR3Helper *MultiOCR3HelperTransactorSession) SetTransmitOcrPluginType(ocrPluginType uint8) (*types.Transaction, error) {
	return _MultiOCR3Helper.Contract.SetTransmitOcrPluginType(&_MultiOCR3Helper.TransactOpts, ocrPluginType)
}

func (_MultiOCR3Helper *MultiOCR3HelperTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _MultiOCR3Helper.contract.Transact(opts, "transferOwnership", to)
}

func (_MultiOCR3Helper *MultiOCR3HelperSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _MultiOCR3Helper.Contract.TransferOwnership(&_MultiOCR3Helper.TransactOpts, to)
}

func (_MultiOCR3Helper *MultiOCR3HelperTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _MultiOCR3Helper.Contract.TransferOwnership(&_MultiOCR3Helper.TransactOpts, to)
}

func (_MultiOCR3Helper *MultiOCR3HelperTransactor) TransmitWithSignatures(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _MultiOCR3Helper.contract.Transact(opts, "transmitWithSignatures", reportContext, report, rs, ss, rawVs)
}

func (_MultiOCR3Helper *MultiOCR3HelperSession) TransmitWithSignatures(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _MultiOCR3Helper.Contract.TransmitWithSignatures(&_MultiOCR3Helper.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_MultiOCR3Helper *MultiOCR3HelperTransactorSession) TransmitWithSignatures(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _MultiOCR3Helper.Contract.TransmitWithSignatures(&_MultiOCR3Helper.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_MultiOCR3Helper *MultiOCR3HelperTransactor) TransmitWithoutSignatures(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte) (*types.Transaction, error) {
	return _MultiOCR3Helper.contract.Transact(opts, "transmitWithoutSignatures", reportContext, report)
}

func (_MultiOCR3Helper *MultiOCR3HelperSession) TransmitWithoutSignatures(reportContext [3][32]byte, report []byte) (*types.Transaction, error) {
	return _MultiOCR3Helper.Contract.TransmitWithoutSignatures(&_MultiOCR3Helper.TransactOpts, reportContext, report)
}

func (_MultiOCR3Helper *MultiOCR3HelperTransactorSession) TransmitWithoutSignatures(reportContext [3][32]byte, report []byte) (*types.Transaction, error) {
	return _MultiOCR3Helper.Contract.TransmitWithoutSignatures(&_MultiOCR3Helper.TransactOpts, reportContext, report)
}

type MultiOCR3HelperAfterConfigSetIterator struct {
	Event *MultiOCR3HelperAfterConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MultiOCR3HelperAfterConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiOCR3HelperAfterConfigSet)
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
		it.Event = new(MultiOCR3HelperAfterConfigSet)
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

func (it *MultiOCR3HelperAfterConfigSetIterator) Error() error {
	return it.fail
}

func (it *MultiOCR3HelperAfterConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MultiOCR3HelperAfterConfigSet struct {
	OcrPluginType uint8
	Raw           types.Log
}

func (_MultiOCR3Helper *MultiOCR3HelperFilterer) FilterAfterConfigSet(opts *bind.FilterOpts) (*MultiOCR3HelperAfterConfigSetIterator, error) {

	logs, sub, err := _MultiOCR3Helper.contract.FilterLogs(opts, "AfterConfigSet")
	if err != nil {
		return nil, err
	}
	return &MultiOCR3HelperAfterConfigSetIterator{contract: _MultiOCR3Helper.contract, event: "AfterConfigSet", logs: logs, sub: sub}, nil
}

func (_MultiOCR3Helper *MultiOCR3HelperFilterer) WatchAfterConfigSet(opts *bind.WatchOpts, sink chan<- *MultiOCR3HelperAfterConfigSet) (event.Subscription, error) {

	logs, sub, err := _MultiOCR3Helper.contract.WatchLogs(opts, "AfterConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MultiOCR3HelperAfterConfigSet)
				if err := _MultiOCR3Helper.contract.UnpackLog(event, "AfterConfigSet", log); err != nil {
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

func (_MultiOCR3Helper *MultiOCR3HelperFilterer) ParseAfterConfigSet(log types.Log) (*MultiOCR3HelperAfterConfigSet, error) {
	event := new(MultiOCR3HelperAfterConfigSet)
	if err := _MultiOCR3Helper.contract.UnpackLog(event, "AfterConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MultiOCR3HelperConfigSetIterator struct {
	Event *MultiOCR3HelperConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MultiOCR3HelperConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiOCR3HelperConfigSet)
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
		it.Event = new(MultiOCR3HelperConfigSet)
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

func (it *MultiOCR3HelperConfigSetIterator) Error() error {
	return it.fail
}

func (it *MultiOCR3HelperConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MultiOCR3HelperConfigSet struct {
	OcrPluginType uint8
	ConfigDigest  [32]byte
	Signers       []common.Address
	Transmitters  []common.Address
	F             uint8
	Raw           types.Log
}

func (_MultiOCR3Helper *MultiOCR3HelperFilterer) FilterConfigSet(opts *bind.FilterOpts) (*MultiOCR3HelperConfigSetIterator, error) {

	logs, sub, err := _MultiOCR3Helper.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &MultiOCR3HelperConfigSetIterator{contract: _MultiOCR3Helper.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_MultiOCR3Helper *MultiOCR3HelperFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *MultiOCR3HelperConfigSet) (event.Subscription, error) {

	logs, sub, err := _MultiOCR3Helper.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MultiOCR3HelperConfigSet)
				if err := _MultiOCR3Helper.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_MultiOCR3Helper *MultiOCR3HelperFilterer) ParseConfigSet(log types.Log) (*MultiOCR3HelperConfigSet, error) {
	event := new(MultiOCR3HelperConfigSet)
	if err := _MultiOCR3Helper.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MultiOCR3HelperOwnershipTransferRequestedIterator struct {
	Event *MultiOCR3HelperOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MultiOCR3HelperOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiOCR3HelperOwnershipTransferRequested)
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
		it.Event = new(MultiOCR3HelperOwnershipTransferRequested)
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

func (it *MultiOCR3HelperOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *MultiOCR3HelperOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MultiOCR3HelperOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_MultiOCR3Helper *MultiOCR3HelperFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MultiOCR3HelperOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MultiOCR3Helper.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MultiOCR3HelperOwnershipTransferRequestedIterator{contract: _MultiOCR3Helper.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_MultiOCR3Helper *MultiOCR3HelperFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *MultiOCR3HelperOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MultiOCR3Helper.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MultiOCR3HelperOwnershipTransferRequested)
				if err := _MultiOCR3Helper.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_MultiOCR3Helper *MultiOCR3HelperFilterer) ParseOwnershipTransferRequested(log types.Log) (*MultiOCR3HelperOwnershipTransferRequested, error) {
	event := new(MultiOCR3HelperOwnershipTransferRequested)
	if err := _MultiOCR3Helper.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MultiOCR3HelperOwnershipTransferredIterator struct {
	Event *MultiOCR3HelperOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MultiOCR3HelperOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiOCR3HelperOwnershipTransferred)
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
		it.Event = new(MultiOCR3HelperOwnershipTransferred)
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

func (it *MultiOCR3HelperOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *MultiOCR3HelperOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MultiOCR3HelperOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_MultiOCR3Helper *MultiOCR3HelperFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MultiOCR3HelperOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MultiOCR3Helper.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MultiOCR3HelperOwnershipTransferredIterator{contract: _MultiOCR3Helper.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_MultiOCR3Helper *MultiOCR3HelperFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MultiOCR3HelperOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MultiOCR3Helper.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MultiOCR3HelperOwnershipTransferred)
				if err := _MultiOCR3Helper.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_MultiOCR3Helper *MultiOCR3HelperFilterer) ParseOwnershipTransferred(log types.Log) (*MultiOCR3HelperOwnershipTransferred, error) {
	event := new(MultiOCR3HelperOwnershipTransferred)
	if err := _MultiOCR3Helper.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MultiOCR3HelperTransmittedIterator struct {
	Event *MultiOCR3HelperTransmitted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MultiOCR3HelperTransmittedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiOCR3HelperTransmitted)
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
		it.Event = new(MultiOCR3HelperTransmitted)
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

func (it *MultiOCR3HelperTransmittedIterator) Error() error {
	return it.fail
}

func (it *MultiOCR3HelperTransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MultiOCR3HelperTransmitted struct {
	OcrPluginType  uint8
	ConfigDigest   [32]byte
	SequenceNumber uint64
	Raw            types.Log
}

func (_MultiOCR3Helper *MultiOCR3HelperFilterer) FilterTransmitted(opts *bind.FilterOpts, ocrPluginType []uint8) (*MultiOCR3HelperTransmittedIterator, error) {

	var ocrPluginTypeRule []interface{}
	for _, ocrPluginTypeItem := range ocrPluginType {
		ocrPluginTypeRule = append(ocrPluginTypeRule, ocrPluginTypeItem)
	}

	logs, sub, err := _MultiOCR3Helper.contract.FilterLogs(opts, "Transmitted", ocrPluginTypeRule)
	if err != nil {
		return nil, err
	}
	return &MultiOCR3HelperTransmittedIterator{contract: _MultiOCR3Helper.contract, event: "Transmitted", logs: logs, sub: sub}, nil
}

func (_MultiOCR3Helper *MultiOCR3HelperFilterer) WatchTransmitted(opts *bind.WatchOpts, sink chan<- *MultiOCR3HelperTransmitted, ocrPluginType []uint8) (event.Subscription, error) {

	var ocrPluginTypeRule []interface{}
	for _, ocrPluginTypeItem := range ocrPluginType {
		ocrPluginTypeRule = append(ocrPluginTypeRule, ocrPluginTypeItem)
	}

	logs, sub, err := _MultiOCR3Helper.contract.WatchLogs(opts, "Transmitted", ocrPluginTypeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MultiOCR3HelperTransmitted)
				if err := _MultiOCR3Helper.contract.UnpackLog(event, "Transmitted", log); err != nil {
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

func (_MultiOCR3Helper *MultiOCR3HelperFilterer) ParseTransmitted(log types.Log) (*MultiOCR3HelperTransmitted, error) {
	event := new(MultiOCR3HelperTransmitted)
	if err := _MultiOCR3Helper.contract.UnpackLog(event, "Transmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_MultiOCR3Helper *MultiOCR3Helper) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _MultiOCR3Helper.abi.Events["AfterConfigSet"].ID:
		return _MultiOCR3Helper.ParseAfterConfigSet(log)
	case _MultiOCR3Helper.abi.Events["ConfigSet"].ID:
		return _MultiOCR3Helper.ParseConfigSet(log)
	case _MultiOCR3Helper.abi.Events["OwnershipTransferRequested"].ID:
		return _MultiOCR3Helper.ParseOwnershipTransferRequested(log)
	case _MultiOCR3Helper.abi.Events["OwnershipTransferred"].ID:
		return _MultiOCR3Helper.ParseOwnershipTransferred(log)
	case _MultiOCR3Helper.abi.Events["Transmitted"].ID:
		return _MultiOCR3Helper.ParseTransmitted(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (MultiOCR3HelperAfterConfigSet) Topic() common.Hash {
	return common.HexToHash("0x897ac1b2c12867721b284f3eb147bd4ab046d4eef1cf31c1d8988bfcfb962b53")
}

func (MultiOCR3HelperConfigSet) Topic() common.Hash {
	return common.HexToHash("0xab8b1b57514019638d7b5ce9c638fe71366fe8e2be1c40a7a80f1733d0e9f547")
}

func (MultiOCR3HelperOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (MultiOCR3HelperOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (MultiOCR3HelperTransmitted) Topic() common.Hash {
	return common.HexToHash("0x198d6990ef96613a9026203077e422916918b03ff47f0be6bee7b02d8e139ef0")
}

func (_MultiOCR3Helper *MultiOCR3Helper) Address() common.Address {
	return _MultiOCR3Helper.address
}

type MultiOCR3HelperInterface interface {
	GetOracle(opts *bind.CallOpts, ocrPluginType uint8, oracleAddress common.Address) (MultiOCR3BaseOracle, error)

	LatestConfigDetails(opts *bind.CallOpts, ocrPluginType uint8) (MultiOCR3BaseOCRConfig, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	SetOCR3Configs(opts *bind.TransactOpts, ocrConfigArgs []MultiOCR3BaseOCRConfigArgs) (*types.Transaction, error)

	SetTransmitOcrPluginType(opts *bind.TransactOpts, ocrPluginType uint8) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	TransmitWithSignatures(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error)

	TransmitWithoutSignatures(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte) (*types.Transaction, error)

	FilterAfterConfigSet(opts *bind.FilterOpts) (*MultiOCR3HelperAfterConfigSetIterator, error)

	WatchAfterConfigSet(opts *bind.WatchOpts, sink chan<- *MultiOCR3HelperAfterConfigSet) (event.Subscription, error)

	ParseAfterConfigSet(log types.Log) (*MultiOCR3HelperAfterConfigSet, error)

	FilterConfigSet(opts *bind.FilterOpts) (*MultiOCR3HelperConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *MultiOCR3HelperConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*MultiOCR3HelperConfigSet, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MultiOCR3HelperOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *MultiOCR3HelperOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*MultiOCR3HelperOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MultiOCR3HelperOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MultiOCR3HelperOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*MultiOCR3HelperOwnershipTransferred, error)

	FilterTransmitted(opts *bind.FilterOpts, ocrPluginType []uint8) (*MultiOCR3HelperTransmittedIterator, error)

	WatchTransmitted(opts *bind.WatchOpts, sink chan<- *MultiOCR3HelperTransmitted, ocrPluginType []uint8) (event.Subscription, error)

	ParseTransmitted(log types.Log) (*MultiOCR3HelperTransmitted, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
