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
	ABI: "[{\"inputs\":[],\"name\":\"CannotTransferToSelf\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"expected\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"actual\",\"type\":\"bytes32\"}],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"ForkedChain\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"enumMultiOCR3Base.InvalidConfigErrorType\",\"name\":\"errorType\",\"type\":\"uint8\"}],\"name\":\"InvalidConfig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeProposedOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NonUniqueSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OracleCannotBeZeroAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OwnerCannotBeZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SignaturesOutOfRegistration\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"}],\"name\":\"StaticConfigCannotBeChanged\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedTransmitter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"WrongMessageLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WrongNumberOfSignatures\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"}],\"name\":\"AfterConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"F\",\"type\":\"uint8\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"oracleAddress\",\"type\":\"address\"}],\"name\":\"getOracle\",\"outputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"index\",\"type\":\"uint8\"},{\"internalType\":\"enumMultiOCR3Base.Role\",\"name\":\"role\",\"type\":\"uint8\"}],\"internalType\":\"structMultiOCR3Base.Oracle\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"}],\"name\":\"latestConfigDetails\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"F\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"n\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"isSignatureVerificationEnabled\",\"type\":\"bool\"}],\"internalType\":\"structMultiOCR3Base.ConfigInfo\",\"name\":\"configInfo\",\"type\":\"tuple\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"}],\"internalType\":\"structMultiOCR3Base.OCRConfig\",\"name\":\"ocrConfig\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"F\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"isSignatureVerificationEnabled\",\"type\":\"bool\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"}],\"internalType\":\"structMultiOCR3Base.OCRConfigArgs[]\",\"name\":\"ocrConfigArgs\",\"type\":\"tuple[]\"}],\"name\":\"setOCR3Configs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"}],\"name\":\"setTransmitOcrPluginType\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[2]\",\"name\":\"reportContext\",\"type\":\"bytes32[2]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmitWithSignatures\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[2]\",\"name\":\"reportContext\",\"type\":\"bytes32[2]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"}],\"name\":\"transmitWithoutSignatures\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x60a080604052346051573315604057600180546001600160a01b0319163317905546608052611f2090816100578239608051818181610f2f015261169b0152f35b639b15e16f60e01b60005260046000fd5b600080fdfe6080604052600436101561001257600080fd5b60003560e01c806310061068146115df578063181f5a771461150f57806334a9c92e146114285780633ecdb95b14610e2d57806379ba509714610d445780637ac0aa1a14610cdc5780638da5cb5b14610c8a578063c673e58414610b2e578063f2fde38b14610a3b5763f716f99f1461008a57600080fd5b34610a365760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610a365760043567ffffffffffffffff8111610a365736602382011215610a365780600401356100e481611c1a565b916100f26040519384611bd9565b8183526024602084019260051b82010190368211610a365760248101925b82841061093b5784610120611e58565b6002906000805b82518110156109395761013a8184611dc8565b5190604082019060ff8251161561090a5760ff60208401511692836000528660205260406000206001810190815460ff8116156000146108c6577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffff62ff00006060860151151560101b1691161782555b60a08301928351956101008751116107dc578651156108975760038301996101eb6101de6101e58d60405192838092611d78565b0382611bd9565b8a611ea3565b6060830151610566575b60005b88518110156103b55773ffffffffffffffffffffffffffffffffffffffff610220828b611dc8565b5116908a600052600360205260ff60408060002060009073ffffffffffffffffffffffffffffffffffffffff86168252602052205460081c16600381101561032e5761038757811561035d578b6040519261027a84611b85565b60ff83168452602084019161032e578f604060ff928f8493865260005260036020528160002073ffffffffffffffffffffffffffffffffffffffff60009216825260205220945116167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008454161783555191600383101561032e576001927fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff61ff0083549260081b169116179055016101f8565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7fd6c62c9b0000000000000000000000000000000000000000000000000000000060005260046000fd5b7f367f56a2000000000000000000000000000000000000000000000000000000006000526004805260246000fd5b509997969091929394959781519167ffffffffffffffff83116105375768010000000000000000831161053757815483835580841061050e575b509060208d989796959493920190600052602060002060005b8381106104e157505050509360019796936104cb7fab8b1b57514019638d7b5ce9c638fe71366fe8e2be1c40a7a80f1733d0e9f547946104bd7f897ac1b2c12867721b284f3eb147bd4ab046d4eef1cf31c1d8988bfcfb962b53999560ff60209a51169460ff86167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00825416179055519485835551916040519687968a88528b88015260a0604088015260a087019101611d78565b908482036060860152611b3b565b9060808301520390a1604051908152a101610127565b825173ffffffffffffffffffffffffffffffffffffffff16818301558e9950602090920191600101610408565b8260005283602060002091820191015b81811061052b57506103ef565b6000815560010161051e565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b8b840161058360405161057d816101de8186611d78565b8b611ea3565b60808401519061010082511161086957815160ff8551166003029060ff821691820361083a57111561080b5781518a51116107dc5781519087547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff61ff008460081b16911617885567ffffffffffffffff8211610537576801000000000000000082116105375780548282558083106107b3575b506020830190600052602060002060005b8381106107895750600193600093508392509050835b61064c575b505050506101f5565b80518310156107845773ffffffffffffffffffffffffffffffffffffffff6106748483611dc8565b5116928d600052600360205260ff60408060002060009073ffffffffffffffffffffffffffffffffffffffff88168252602052205460081c16600381101561032e5761038757831561035d5782604051946106ce86611b85565b60ff83168652602086019161032e578f604060ff9283928a865260005260036020528160002073ffffffffffffffffffffffffffffffffffffffff60009216825260205220965116167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008654161785555190600382101561032e57859485927fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff61ff0083549260081b169116179055019261063e565b610643565b600190602073ffffffffffffffffffffffffffffffffffffffff8551169401938184015501610628565b8160005282602060002091820191015b8181106107d05750610617565b600081556001016107c3565b7f367f56a200000000000000000000000000000000000000000000000000000000600052600160045260246000fd5b7f367f56a200000000000000000000000000000000000000000000000000000000600052600360045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8d7f367f56a20000000000000000000000000000000000000000000000000000000060005260045260246000fd5b7f367f56a200000000000000000000000000000000000000000000000000000000600052600560045260246000fd5b60ff606085015115159160101c16151503156101aa57857f87f6037c0000000000000000000000000000000000000000000000000000000060005260045260246000fd5b7f367f56a200000000000000000000000000000000000000000000000000000000600052600060045260246000fd5b005b833567ffffffffffffffff8111610a3657820160c07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdc8236030112610a36576040519160c0830183811067ffffffffffffffff82111761053757604052602482013583526109ab60448301611afc565b60208401526109bc60648301611afc565b604084015260848201358015158103610a3657606084015260a482013567ffffffffffffffff8111610a36576109f89060243691850101611c32565b608084015260c48201359267ffffffffffffffff8411610a3657610a26602094936024869536920101611c32565b60a0820152815201930192610110565b600080fd5b34610a365760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610a365760043573ffffffffffffffffffffffffffffffffffffffff8116809103610a3657610a93611e58565b338114610b0457807fffffffffffffffffffffffff0000000000000000000000000000000000000000600054161760005573ffffffffffffffffffffffffffffffffffffffff600154167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278600080a3005b7fdad89dca0000000000000000000000000000000000000000000000000000000060005260046000fd5b34610a365760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610a365760ff610b67611aec565b606060408051610b7681611ba1565b8151610b8181611bbd565b60008152600060208201526000838201526000848201528152826020820152015216600052600260205260606040600020610c866003610c5560405193610bc785611ba1565b610bd081611d3f565b8552610c0860405191610bf183610bea8160028501611d78565b0384611bd9565b60208701928352610bea6040518096819301611d78565b6040850192835260405195869560208752518051602088015260ff602082015116604088015260ff604082015116828801520151151560808601525160c060a086015260e0850190611b3b565b90517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08483030160c0850152611b3b565b0390f35b34610a365760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610a3657602073ffffffffffffffffffffffffffffffffffffffff60015416604051908152f35b34610a365760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610a365760ff610d15611aec565b167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff006004541617600455600080f35b34610a365760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610a365760005473ffffffffffffffffffffffffffffffffffffffff81163303610e03577fffffffffffffffffffffffff00000000000000000000000000000000000000006001549133828416176001551660005573ffffffffffffffffffffffffffffffffffffffff3391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0600080a3005b7f02b543c60000000000000000000000000000000000000000000000000000000060005260046000fd5b34610a365760c07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610a365736604411610a365760443567ffffffffffffffff8111610a3657610e84903690600401611abe565b9060643567ffffffffffffffff8111610a3657610ea5903690600401611b0a565b919060843567ffffffffffffffff8111610a3657610eca610ee9913690600401611b0a565b9190610ee160a4359460ff60045416973691611cf3565b923691611cf3565b90846000526002602052610f006040600020611d3f565b9560043594610f0e82611e0b565b97606081019889516113db575b8036036113aa5750805187810361137857507f00000000000000000000000000000000000000000000000000000000000000004681036113475750876000526003602052604060002073ffffffffffffffffffffffffffffffffffffffff331660005260205260406000209860405199610f948b611b85565b5460ff81168b52610faf60ff60208d019260081c1682611ce7565b519960038b101561032e57600260009b1490816112d2575b50156112aa5751611014575b88887f198d6990ef96613a9026203077e422916918b03ff47f0be6bee7b02d8e139ef060408a815190815267ffffffffffffffff602435166020820152a280f35b60ff611027816020875194015116611e34565b160361128257825184510361125a578761104083611cad565b9161104e6040519384611bd9565b8383526020830193368183011161125657806020928637830101525190206040516020810191825260406004818301376060815261108d608082611bd9565b51902090869281519488945b8686106110a65750610fd3565b60208610156112295760208a60806110bf858a1a611e46565b6110c98a89611dc8565b516110d48b89611dc8565b519060ff604051938c855216868401526040830152606082015282805260015afa1561121e578951898b52600360205273ffffffffffffffffffffffffffffffffffffffff60408c2091168b5260205260408a206040519061113582611b85565b5460ff8116825261115060ff602084019260081c1682611ce7565b5160038110156111f1577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff016111c957600160ff8251161b82166111a15790600160ff819351161b17950194611099565b60048b7ff67bc7c4000000000000000000000000000000000000000000000000000000008152fd5b60048b7fca31867a000000000000000000000000000000000000000000000000000000008152fd5b60248c7f4e487b710000000000000000000000000000000000000000000000000000000081526021600452fd5b6040513d8b823e3d90fd5b60248a7f4e487b710000000000000000000000000000000000000000000000000000000081526032600452fd5b8280fd5b6004887fa75d88af000000000000000000000000000000000000000000000000000000008152fd5b6004887f71253a25000000000000000000000000000000000000000000000000000000008152fd5b60048a7fda0f08e8000000000000000000000000000000000000000000000000000000008152fd5b9050898b52600260205260ff600360408d200191511690805482101561131a579073ffffffffffffffffffffffffffffffffffffffff918c5260208c2001541633148b610fc7565b60248c7f4e487b710000000000000000000000000000000000000000000000000000000081526032600452fd5b7f0f01ce85000000000000000000000000000000000000000000000000000000006000526004524660245260446000fd5b87907f93df584c0000000000000000000000000000000000000000000000000000000060005260045260245260446000fd5b7f8e1192e1000000000000000000000000000000000000000000000000000000006000526004523660245260446000fd5b84518060051b908082046020149015171561083a576113f990611e19565b908651918260051b928084046020149015171561083a576114239261141d91611e27565b90611e27565b610f1b565b34610a365760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610a365761145f611aec565b6024359073ffffffffffffffffffffffffffffffffffffffff82168203610a365760ff906000602060405161149381611b85565b828152015216600052600360205273ffffffffffffffffffffffffffffffffffffffff604060002091166000526020526040600020604051906114d582611b85565b5460ff811682526114f060ff602084019260081c1682611ce7565b60ff60405192511682525190600382101561032e576040916020820152f35b34610a365760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610a3657604080519061154d8183611bd9565b601982527f4d756c74694f4352334261736548656c70657220312e302e30000000000000006020830152805180926020825280519081602084015260005b8281106115c85750506000828201840152601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168101030190f35b60208282018101518783018701528694500161158b565b34610a365760607ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610a365736604411610a365760443567ffffffffffffffff8111610a3657611636903690600401611abe565b604051602092916116478483611bd9565b60008252600036813760ff6004541691826000526002855261166c6040600020611d3f565b936004359261167a81611e0b565b9560608101968751611a79575b8036036113aa57508051858103611a4757507f000000000000000000000000000000000000000000000000000000000000000046810361134757508560005260038852604060002073ffffffffffffffffffffffffffffffffffffffff33166000528852604060002096604051976116fe89611b85565b5460ff8116895261171860ff8b8b019260081c1682611ce7565b5197600389101561032e576002600099149081611a01575b50156119d9575161177d575b86867f198d6990ef96613a9026203077e422916918b03ff47f0be6bee7b02d8e139ef06040888c825191825267ffffffffffffffff6024351690820152a280f35b60ff61178f818a875194015116611e34565b16036119b15761179e81611cad565b906117ac6040519283611bd9565b8082528782019236828201116119ad578188928a928637830101525190206040518681019182526040600481830137606081526117ea608082611bd9565b519020849082519286925b848410611802575061173c565b88841015611980578888608061181982881a611e46565b6118238887611dc8565b5161182e8988611dc8565b519060ff604051938a855216868401526040830152606082015282805260015afa1561197557875187895260038a5273ffffffffffffffffffffffffffffffffffffffff60408a20911689528952604088206040519061188d82611b85565b5460ff811682526118a760ff8c84019260081c1682611ce7565b516003811015611948577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0161192057600160ff8251161b82166118f85790600160ff819351161b179301926117f5565b6004897ff67bc7c4000000000000000000000000000000000000000000000000000000008152fd5b6004897fca31867a000000000000000000000000000000000000000000000000000000008152fd5b60248a7f4e487b710000000000000000000000000000000000000000000000000000000081526021600452fd5b6040513d89823e3d90fd5b6024887f4e487b710000000000000000000000000000000000000000000000000000000081526032600452fd5b8780fd5b6004867f71253a25000000000000000000000000000000000000000000000000000000008152fd5b6004887fda0f08e8000000000000000000000000000000000000000000000000000000008152fd5b905087895260028a5260ff600360408b2001915116908054821015611229579073ffffffffffffffffffffffffffffffffffffffff918a528a8a2001541633148a611730565b85907f93df584c0000000000000000000000000000000000000000000000000000000060005260045260245260446000fd5b84518060051b908082048b149015171561083a57611a9690611e19565b908551918260051b928084048c149015171561083a57611ab99261141d91611e27565b611687565b9181601f84011215610a365782359167ffffffffffffffff8311610a365760208381860195010111610a3657565b6004359060ff82168203610a3657565b359060ff82168203610a3657565b9181601f84011215610a365782359167ffffffffffffffff8311610a36576020808501948460051b010111610a3657565b906020808351928381520192019060005b818110611b595750505090565b825173ffffffffffffffffffffffffffffffffffffffff16845260209384019390920191600101611b4c565b6040810190811067ffffffffffffffff82111761053757604052565b6060810190811067ffffffffffffffff82111761053757604052565b6080810190811067ffffffffffffffff82111761053757604052565b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff82111761053757604052565b67ffffffffffffffff81116105375760051b60200190565b9080601f83011215610a3657813590611c4a82611c1a565b92611c586040519485611bd9565b82845260208085019360051b820101918211610a3657602001915b818310611c805750505090565b823573ffffffffffffffffffffffffffffffffffffffff81168103610a3657815260209283019201611c73565b67ffffffffffffffff811161053757601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600382101561032e5752565b929190611cff81611c1a565b93611d0d6040519586611bd9565b602085838152019160051b8101928311610a3657905b828210611d2f57505050565b8135815260209182019101611d23565b90604051611d4c81611bbd565b606060ff600183958054855201548181166020850152818160081c16604085015260101c161515910152565b906020825491828152019160005260206000209060005b818110611d9c5750505090565b825473ffffffffffffffffffffffffffffffffffffffff16845260209093019260019283019201611d8f565b8051821015611ddc5760209160051b010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b608401908160841161083a57565b60a001908160a01161083a57565b9190820180921161083a57565b60ff60019116019060ff821161083a57565b60ff601b9116019060ff821161083a57565b73ffffffffffffffffffffffffffffffffffffffff600154163303611e7957565b7f2b5c74de0000000000000000000000000000000000000000000000000000000060005260046000fd5b91909160005b8351811015611f0d5760019060ff831660005260036020526000604080822073ffffffffffffffffffffffffffffffffffffffff611ee7858a611dc8565b511673ffffffffffffffffffffffffffffffffffffffff16835260205281205501611ea9565b5050905056fea164736f6c634300081a000a",
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

func (_MultiOCR3Helper *MultiOCR3HelperTransactor) TransmitWithSignatures(opts *bind.TransactOpts, reportContext [2][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _MultiOCR3Helper.contract.Transact(opts, "transmitWithSignatures", reportContext, report, rs, ss, rawVs)
}

func (_MultiOCR3Helper *MultiOCR3HelperSession) TransmitWithSignatures(reportContext [2][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _MultiOCR3Helper.Contract.TransmitWithSignatures(&_MultiOCR3Helper.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_MultiOCR3Helper *MultiOCR3HelperTransactorSession) TransmitWithSignatures(reportContext [2][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _MultiOCR3Helper.Contract.TransmitWithSignatures(&_MultiOCR3Helper.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_MultiOCR3Helper *MultiOCR3HelperTransactor) TransmitWithoutSignatures(opts *bind.TransactOpts, reportContext [2][32]byte, report []byte) (*types.Transaction, error) {
	return _MultiOCR3Helper.contract.Transact(opts, "transmitWithoutSignatures", reportContext, report)
}

func (_MultiOCR3Helper *MultiOCR3HelperSession) TransmitWithoutSignatures(reportContext [2][32]byte, report []byte) (*types.Transaction, error) {
	return _MultiOCR3Helper.Contract.TransmitWithoutSignatures(&_MultiOCR3Helper.TransactOpts, reportContext, report)
}

func (_MultiOCR3Helper *MultiOCR3HelperTransactorSession) TransmitWithoutSignatures(reportContext [2][32]byte, report []byte) (*types.Transaction, error) {
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

	TransmitWithSignatures(opts *bind.TransactOpts, reportContext [2][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error)

	TransmitWithoutSignatures(opts *bind.TransactOpts, reportContext [2][32]byte, report []byte) (*types.Transaction, error)

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
