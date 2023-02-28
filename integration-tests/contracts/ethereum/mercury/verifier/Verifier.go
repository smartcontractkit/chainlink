// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package verifier

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// VerifierMetaData contains all meta data concerning the Verifier contract.
var VerifierMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"verifierProxyAddr\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessForbidden\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"CannotDeactivateLatestConfig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DigestEmpty\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"DigestInactive\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"DigestNotSet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxSigners\",\"type\":\"uint256\"}],\"name\":\"ExcessSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FaultToleranceMustBePositive\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FeedIdEmpty\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"requiredFeedId\",\"type\":\"bytes32\"}],\"name\":\"FeedIdIncorrect\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expectedNumSigners\",\"type\":\"uint256\"}],\"name\":\"IncorrectSignatureCount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minSigners\",\"type\":\"uint256\"}],\"name\":\"InsufficientSigners\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"rsLength\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ssLength\",\"type\":\"uint256\"}],\"name\":\"MismatchedSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NonUniqueSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"ConfigActivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"ConfigDeactivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"reportHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"}],\"name\":\"ReportVerified\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"activateConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"deactivateConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFeedId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isVerifier\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"signedReport\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"verify\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"response\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60c06040523480156200001157600080fd5b50604051620019e5380380620019e58339810160408190526200003491620001d0565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be8162000124565b50505081620000e05760405163bd9f4c2b60e01b815260040160405180910390fd5b6001600160a01b038116620001085760405163d92e233d60e01b815260040160405180910390fd5b60809190915260601b6001600160601b03191660a0526200020f565b6001600160a01b0381163314156200017f5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60008060408385031215620001e457600080fd5b825160208401519092506001600160a01b03811681146200020457600080fd5b809150509250929050565b60805160a05160601c6117986200024d600039600081816103f101526109fc0152600081816101d30152818161045b015261049601526117986000f3fe608060405234801561001057600080fd5b50600436106100cf5760003560e01c806381ff70481161008c578063afcb95d711610066578063afcb95d7146101fc578063b1dc65a41461022c578063e3d0e7121461023f578063f2fde38b1461025257600080fd5b806381ff7048146101835780638da5cb5b146101b357806393e5c6e9146101ce57600080fd5b806301ffc9a7146100d45780630d1d79af1461010d5780630f672ef414610122578063181f5a77146101355780633d3ac1b51461016857806379ba50971461017b575b600080fd5b6100f86100e2366004611332565b6001600160e01b031916633d3ac1b560e01b1490565b60405190151581526020015b60405180910390f35b61012061011b366004611319565b610265565b005b610120610130366004611319565b61031b565b60408051808201909152600e81526d566572696669657220302e302e3160901b60208201525b604051610104919061147c565b61015b61017636600461135c565b6103e4565b61012061065d565b6003546004546040805163ffffffff80851682526401000000009094049093166020840152820152606001610104565b6000546040516001600160a01b039091168152602001610104565b6040517f00000000000000000000000000000000000000000000000000000000000000008152602001610104565b60045460035460408051600081526020810193909352600160401b90910463ffffffff1690820152606001610104565b61012061023a36600461117f565b610707565b61012061024d3660046110b3565b61074f565b610120610260366004611091565b610aec565b61026d610b00565b8061028b5760405163e332262760e01b815260040160405180910390fd5b60008181526002602052604090205460ff166102c2576040516374eb4b9360e01b8152600481018290526024015b60405180910390fd5b60008181526002602081905260409182902001805460ff19166001179055517fa543797a0501218bba8a3daf75a71c8df8d1a7f791f4e44d40e43b6450183cea906103109083815260200190565b60405180910390a150565b610323610b00565b806103415760405163e332262760e01b815260040160405180910390fd5b60008181526002602052604090205460ff16610373576040516374eb4b9360e01b8152600481018290526024016102b9565b600454811415610399576040516319e18fd160e21b8152600481018290526024016102b9565b60008181526002602081905260409182902001805460ff19169055517f5bfaab86edc1b932e3c334327a591c9ded067cb521abae19b95ca927d6076579906103109083815260200190565b6060336001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000161461042f57604051631decfebb60e31b815260040160405180910390fd5b600080808080610441888a018a611235565b9450945094509450945060008461045790611684565b90507f000000000000000000000000000000000000000000000000000000000000000081146104c257604051636782d0b960e11b8152600481018290527f000000000000000000000000000000000000000000000000000000000000000060248201526044016102b9565b8551600081815260026020526040812080549091906104e59060ff166001611640565b60008481526002602052604090205490915060ff1661051a576040516374eb4b9360e01b8152600481018490526024016102b9565b600282015460ff166105425760405163d990d62160e01b8152600481018490526024016102b9565b8060ff168751146105755786516040516329a4514160e11b8152600481019190915260ff821660248201526044016102b9565b85518751146105a45786518651604051631e1a628160e31b8152600481019290925260248201526044016102b9565b50602088015160035463ffffffff600883901c811691600160401b9004168111156105ed57600380546bffffffff00000000000000001916600160401b63ffffffff8416021790555b885160208a0120610602818c8b8b8b89610b55565b60408051878152602081018390526001600160a01b038f168183015290517f5a779ad0380d86bfcdb2748a39a349fe4fda58c933a2df958ba997284db8a2699181900360600190a150979d9c50505050505050505050505050565b6001546001600160a01b031633146106b05760405162461bcd60e51b815260206004820152601660248201527526bab9ba10313290383937b837b9b2b21037bbb732b960511b60448201526064016102b9565b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60405162461bcd60e51b815260206004820152601a60248201527f7472616e736d69742066756e6374696f6e2064697361626c656400000000000060448201526064016102b9565b855160ff851680610773576040516303a1dd7360e11b815260040160405180910390fd5b601f82111561079f57604051630185d43d60e61b815260048101839052601f60248201526044016102b9565b6107aa816003611665565b82116107e957816107bc826003611665565b6107c7906001611628565b6040516313bb3cdb60e31b8152600481019290925260248201526044016102b9565b6107f1610b00565b6003805463ffffffff16906000610807836116c6565b82546101009290920a63ffffffff8181021990931691831602179091556003546000925061083f9146913091168c8c8c8c8c8c610d74565b60008181526002602081905260408220805460ff8c1660ff199182161782559101805490911660011790559091505b89518160ff1610156109ca5760008a8260ff168151811061089157610891611736565b6020026020010151905060006001600160a01b0316816001600160a01b031614156108cf5760405163d92e233d60e01b815260040160405180910390fd5b60008060008581526002602090815260408083206001600160a01b0387168452600190810190925290912054610100900460ff169081111561091357610913611720565b148015915061093557604051633d9ef1f160e21b815260040160405180910390fd5b6040805180820190915260ff84168152602081016001905260008581526002602090815260408083206001600160a01b03871684526001908101835292208351815460ff90911660ff198216811783559285015191939192849261ffff1990921617906101009084908111156109ad576109ad611720565b0217905550905050505080806109c2906116ea565b91505061086e565b50600354600163ffffffff9091161115610a615760048054604051632cc9947760e01b815291820152602481018290527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031690632cc9947790604401600060405180830381600087803b158015610a4857600080fd5b505af1158015610a5c573d6000803e3d6000fd5b505050505b6003546040517f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e0591610ab39163ffffffff640100000000830481169286929116908e908e908e908e908e908e90611528565b60405180910390a1600380546bffffffffffffffff00000000191664010000000063ffffffff4316021790556004555050505050505050565b610af4610b00565b610afd81610dcf565b50565b6000546001600160a01b03163314610b535760405162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b60448201526064016102b9565b565b60008686604051602001610b6a929190611440565b6040516020818303038152906040528051906020012090506000610b9e604080518082019091526000808252602082015290565b8651600090815b81811015610d2557600186898360208110610bc257610bc2611736565b610bcf91901a601b611640565b8c8481518110610be157610be1611736565b60200260200101518c8581518110610bfb57610bfb611736565b602002602001015160405160008152602001604052604051610c39949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa158015610c5b573d6000803e3d6000fd5b505060408051601f198101516001600160a01b03811660009081526001808d01602090815291859020848601909552845460ff808216865293995093955090850192610100900490911690811115610cb557610cb5611720565b6001811115610cc657610cc6611720565b9052509350600184602001516001811115610ce357610ce3611720565b14610d0157604051631decfebb60e31b815260040160405180910390fd5b836000015160080260ff166001901b8501945080610d1e906116ab565b9050610ba5565b50837e01010101010101010101010101010101010101010101010101010101010101851614610d6757604051633d9ef1f160e21b815260040160405180910390fd5b5050505050505050505050565b6000808a8a8a8a8a8a8a8a8a604051602001610d989998979695949392919061148f565b60408051601f1981840301815291905280516020909101206001600160f01b0316600160f01b179150509998505050505050505050565b6001600160a01b038116331415610e285760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016102b9565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80356001600160a01b0381168114610e9057600080fd5b919050565b600082601f830112610ea657600080fd5b81356020610ebb610eb683611605565b6115d5565b80838252828201915082860187848660051b8901011115610edb57600080fd5b60005b85811015610f0157610eef82610e79565b84529284019290840190600101610ede565b5090979650505050505050565b60008083601f840112610f2057600080fd5b5081356001600160401b03811115610f3757600080fd5b6020830191508360208260051b8501011115610f5257600080fd5b9250929050565b600082601f830112610f6a57600080fd5b81356020610f7a610eb683611605565b80838252828201915082860187848660051b8901011115610f9a57600080fd5b60005b85811015610f0157813584529284019290840190600101610f9d565b60008083601f840112610fcb57600080fd5b5081356001600160401b03811115610fe257600080fd5b602083019150836020828501011115610f5257600080fd5b600082601f83011261100b57600080fd5b81356001600160401b038111156110245761102461174c565b611037601f8201601f19166020016115d5565b81815284602083860101111561104c57600080fd5b816020850160208301376000918101602001919091529392505050565b80356001600160401b0381168114610e9057600080fd5b803560ff81168114610e9057600080fd5b6000602082840312156110a357600080fd5b6110ac82610e79565b9392505050565b60008060008060008060c087890312156110cc57600080fd5b86356001600160401b03808211156110e357600080fd5b6110ef8a838b01610e95565b9750602089013591508082111561110557600080fd5b6111118a838b01610e95565b965061111f60408a01611080565b9550606089013591508082111561113557600080fd5b6111418a838b01610ffa565b945061114f60808a01611069565b935060a089013591508082111561116557600080fd5b5061117289828a01610ffa565b9150509295509295509295565b60008060008060008060008060e0898b03121561119b57600080fd5b606089018a8111156111ac57600080fd5b899850356001600160401b03808211156111c557600080fd5b6111d18c838d01610fb9565b909950975060808b01359150808211156111ea57600080fd5b6111f68c838d01610f0e565b909750955060a08b013591508082111561120f57600080fd5b5061121c8b828c01610f0e565b999c989b50969995989497949560c00135949350505050565b600080600080600060e0868803121561124d57600080fd5b86601f87011261125c57600080fd5b6112646115ad565b8087606089018a81111561127757600080fd5b60005b600381101561129957823585526020948501949092019160010161127a565b50919750503590506001600160401b03808211156112b657600080fd5b6112c289838a01610ffa565b955060808801359150808211156112d857600080fd5b6112e489838a01610f59565b945060a08801359150808211156112fa57600080fd5b5061130788828901610f59565b9598949750929560c001359392505050565b60006020828403121561132b57600080fd5b5035919050565b60006020828403121561134457600080fd5b81356001600160e01b0319811681146110ac57600080fd5b60008060006040848603121561137157600080fd5b83356001600160401b0381111561138757600080fd5b61139386828701610fb9565b90945092506113a6905060208501610e79565b90509250925092565b600081518084526020808501945080840160005b838110156113e85781516001600160a01b0316875295820195908201906001016113c3565b509495945050505050565b6000815180845260005b81811015611419576020818501810151868301820152016113fd565b8181111561142b576000602083870101525b50601f01601f19169290920160200192915050565b828152600060208083018460005b600381101561146b5781518352918301919083019060010161144e565b505050506080820190509392505050565b6020815260006110ac60208301846113f3565b8981526001600160a01b03891660208201526001600160401b038881166040830152610120606083018190526000916114ca8483018b6113af565b915083820360808501526114de828a6113af565b915060ff881660a085015283820360c08501526114fb82886113f3565b90861660e0850152838103610100850152905061151881856113f3565b9c9b505050505050505050505050565b600061012063ffffffff808d1684528b6020850152808b166040850152508060608401526115588184018a6113af565b9050828103608084015261156c81896113af565b905060ff871660a084015282810360c084015261158981876113f3565b90506001600160401b03851660e084015282810361010084015261151881856113f3565b604051606081016001600160401b03811182821017156115cf576115cf61174c565b60405290565b604051601f8201601f191681016001600160401b03811182821017156115fd576115fd61174c565b604052919050565b60006001600160401b0382111561161e5761161e61174c565b5060051b60200190565b6000821982111561163b5761163b61170a565b500190565b600060ff821660ff84168060ff0382111561165d5761165d61170a565b019392505050565b600081600019048311821515161561167f5761167f61170a565b500290565b805160208083015191908110156116a5576000198160200360031b1b821691505b50919050565b60006000198214156116bf576116bf61170a565b5060010190565b600063ffffffff808316818114156116e0576116e061170a565b6001019392505050565b600060ff821660ff8114156117015761170161170a565b60010192915050565b634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052602160045260246000fd5b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052604160045260246000fdfea26469706673582212205a3de6471276fbfa4b8bff0ab766089e74c52e65ae0d36e50fb416496099718d64736f6c63430008060033",
}

// VerifierABI is the input ABI used to generate the binding from.
// Deprecated: Use VerifierMetaData.ABI instead.
var VerifierABI = VerifierMetaData.ABI

// VerifierBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use VerifierMetaData.Bin instead.
var VerifierBin = VerifierMetaData.Bin

// DeployVerifier deploys a new Ethereum contract, binding an instance of Verifier to it.
func DeployVerifier(auth *bind.TransactOpts, backend bind.ContractBackend, feedId [32]byte, verifierProxyAddr common.Address) (common.Address, *types.Transaction, *Verifier, error) {
	parsed, err := VerifierMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VerifierBin), backend, feedId, verifierProxyAddr)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Verifier{VerifierCaller: VerifierCaller{contract: contract}, VerifierTransactor: VerifierTransactor{contract: contract}, VerifierFilterer: VerifierFilterer{contract: contract}}, nil
}

// Verifier is an auto generated Go binding around an Ethereum contract.
type Verifier struct {
	VerifierCaller     // Read-only binding to the contract
	VerifierTransactor // Write-only binding to the contract
	VerifierFilterer   // Log filterer for contract events
}

// VerifierCaller is an auto generated read-only Go binding around an Ethereum contract.
type VerifierCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VerifierTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VerifierTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VerifierFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VerifierFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VerifierSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VerifierSession struct {
	Contract     *Verifier         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VerifierCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VerifierCallerSession struct {
	Contract *VerifierCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// VerifierTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VerifierTransactorSession struct {
	Contract     *VerifierTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// VerifierRaw is an auto generated low-level Go binding around an Ethereum contract.
type VerifierRaw struct {
	Contract *Verifier // Generic contract binding to access the raw methods on
}

// VerifierCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VerifierCallerRaw struct {
	Contract *VerifierCaller // Generic read-only contract binding to access the raw methods on
}

// VerifierTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VerifierTransactorRaw struct {
	Contract *VerifierTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVerifier creates a new instance of Verifier, bound to a specific deployed contract.
func NewVerifier(address common.Address, backend bind.ContractBackend) (*Verifier, error) {
	contract, err := bindVerifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Verifier{VerifierCaller: VerifierCaller{contract: contract}, VerifierTransactor: VerifierTransactor{contract: contract}, VerifierFilterer: VerifierFilterer{contract: contract}}, nil
}

// NewVerifierCaller creates a new read-only instance of Verifier, bound to a specific deployed contract.
func NewVerifierCaller(address common.Address, caller bind.ContractCaller) (*VerifierCaller, error) {
	contract, err := bindVerifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VerifierCaller{contract: contract}, nil
}

// NewVerifierTransactor creates a new write-only instance of Verifier, bound to a specific deployed contract.
func NewVerifierTransactor(address common.Address, transactor bind.ContractTransactor) (*VerifierTransactor, error) {
	contract, err := bindVerifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VerifierTransactor{contract: contract}, nil
}

// NewVerifierFilterer creates a new log filterer instance of Verifier, bound to a specific deployed contract.
func NewVerifierFilterer(address common.Address, filterer bind.ContractFilterer) (*VerifierFilterer, error) {
	contract, err := bindVerifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VerifierFilterer{contract: contract}, nil
}

// bindVerifier binds a generic wrapper to an already deployed contract.
func bindVerifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VerifierABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Verifier *VerifierRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Verifier.Contract.VerifierCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Verifier *VerifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Verifier.Contract.VerifierTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Verifier *VerifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Verifier.Contract.VerifierTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Verifier *VerifierCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Verifier.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Verifier *VerifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Verifier.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Verifier *VerifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Verifier.Contract.contract.Transact(opts, method, params...)
}

// GetFeedId is a free data retrieval call binding the contract method 0x93e5c6e9.
//
// Solidity: function getFeedId() view returns(bytes32 feedId)
func (_Verifier *VerifierCaller) GetFeedId(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Verifier.contract.Call(opts, &out, "getFeedId")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetFeedId is a free data retrieval call binding the contract method 0x93e5c6e9.
//
// Solidity: function getFeedId() view returns(bytes32 feedId)
func (_Verifier *VerifierSession) GetFeedId() ([32]byte, error) {
	return _Verifier.Contract.GetFeedId(&_Verifier.CallOpts)
}

// GetFeedId is a free data retrieval call binding the contract method 0x93e5c6e9.
//
// Solidity: function getFeedId() view returns(bytes32 feedId)
func (_Verifier *VerifierCallerSession) GetFeedId() ([32]byte, error) {
	return _Verifier.Contract.GetFeedId(&_Verifier.CallOpts)
}

// LatestConfigDetails is a free data retrieval call binding the contract method 0x81ff7048.
//
// Solidity: function latestConfigDetails() view returns(uint32 configCount, uint32 blockNumber, bytes32 configDigest)
func (_Verifier *VerifierCaller) LatestConfigDetails(opts *bind.CallOpts) (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}, error) {
	var out []interface{}
	err := _Verifier.contract.Call(opts, &out, "latestConfigDetails")

	outstruct := new(struct {
		ConfigCount  uint32
		BlockNumber  uint32
		ConfigDigest [32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ConfigDigest = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

// LatestConfigDetails is a free data retrieval call binding the contract method 0x81ff7048.
//
// Solidity: function latestConfigDetails() view returns(uint32 configCount, uint32 blockNumber, bytes32 configDigest)
func (_Verifier *VerifierSession) LatestConfigDetails() (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}, error) {
	return _Verifier.Contract.LatestConfigDetails(&_Verifier.CallOpts)
}

// LatestConfigDetails is a free data retrieval call binding the contract method 0x81ff7048.
//
// Solidity: function latestConfigDetails() view returns(uint32 configCount, uint32 blockNumber, bytes32 configDigest)
func (_Verifier *VerifierCallerSession) LatestConfigDetails() (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}, error) {
	return _Verifier.Contract.LatestConfigDetails(&_Verifier.CallOpts)
}

// LatestConfigDigestAndEpoch is a free data retrieval call binding the contract method 0xafcb95d7.
//
// Solidity: function latestConfigDigestAndEpoch() view returns(bool scanLogs, bytes32 configDigest, uint32 epoch)
func (_Verifier *VerifierCaller) LatestConfigDigestAndEpoch(opts *bind.CallOpts) (struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}, error) {
	var out []interface{}
	err := _Verifier.contract.Call(opts, &out, "latestConfigDigestAndEpoch")

	outstruct := new(struct {
		ScanLogs     bool
		ConfigDigest [32]byte
		Epoch        uint32
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ScanLogs = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfigDigest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Epoch = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

// LatestConfigDigestAndEpoch is a free data retrieval call binding the contract method 0xafcb95d7.
//
// Solidity: function latestConfigDigestAndEpoch() view returns(bool scanLogs, bytes32 configDigest, uint32 epoch)
func (_Verifier *VerifierSession) LatestConfigDigestAndEpoch() (struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}, error) {
	return _Verifier.Contract.LatestConfigDigestAndEpoch(&_Verifier.CallOpts)
}

// LatestConfigDigestAndEpoch is a free data retrieval call binding the contract method 0xafcb95d7.
//
// Solidity: function latestConfigDigestAndEpoch() view returns(bool scanLogs, bytes32 configDigest, uint32 epoch)
func (_Verifier *VerifierCallerSession) LatestConfigDigestAndEpoch() (struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}, error) {
	return _Verifier.Contract.LatestConfigDigestAndEpoch(&_Verifier.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Verifier *VerifierCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Verifier.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Verifier *VerifierSession) Owner() (common.Address, error) {
	return _Verifier.Contract.Owner(&_Verifier.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Verifier *VerifierCallerSession) Owner() (common.Address, error) {
	return _Verifier.Contract.Owner(&_Verifier.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) pure returns(bool isVerifier)
func (_Verifier *VerifierCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _Verifier.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) pure returns(bool isVerifier)
func (_Verifier *VerifierSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Verifier.Contract.SupportsInterface(&_Verifier.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) pure returns(bool isVerifier)
func (_Verifier *VerifierCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Verifier.Contract.SupportsInterface(&_Verifier.CallOpts, interfaceId)
}

// Transmit is a free data retrieval call binding the contract method 0xb1dc65a4.
//
// Solidity: function transmit(bytes32[3] , bytes , bytes32[] , bytes32[] , bytes32 ) pure returns()
func (_Verifier *VerifierCaller) Transmit(opts *bind.CallOpts, arg0 [3][32]byte, arg1 []byte, arg2 [][32]byte, arg3 [][32]byte, arg4 [32]byte) error {
	var out []interface{}
	err := _Verifier.contract.Call(opts, &out, "transmit", arg0, arg1, arg2, arg3, arg4)

	if err != nil {
		return err
	}

	return err

}

// Transmit is a free data retrieval call binding the contract method 0xb1dc65a4.
//
// Solidity: function transmit(bytes32[3] , bytes , bytes32[] , bytes32[] , bytes32 ) pure returns()
func (_Verifier *VerifierSession) Transmit(arg0 [3][32]byte, arg1 []byte, arg2 [][32]byte, arg3 [][32]byte, arg4 [32]byte) error {
	return _Verifier.Contract.Transmit(&_Verifier.CallOpts, arg0, arg1, arg2, arg3, arg4)
}

// Transmit is a free data retrieval call binding the contract method 0xb1dc65a4.
//
// Solidity: function transmit(bytes32[3] , bytes , bytes32[] , bytes32[] , bytes32 ) pure returns()
func (_Verifier *VerifierCallerSession) Transmit(arg0 [3][32]byte, arg1 []byte, arg2 [][32]byte, arg3 [][32]byte, arg4 [32]byte) error {
	return _Verifier.Contract.Transmit(&_Verifier.CallOpts, arg0, arg1, arg2, arg3, arg4)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_Verifier *VerifierCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Verifier.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_Verifier *VerifierSession) TypeAndVersion() (string, error) {
	return _Verifier.Contract.TypeAndVersion(&_Verifier.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_Verifier *VerifierCallerSession) TypeAndVersion() (string, error) {
	return _Verifier.Contract.TypeAndVersion(&_Verifier.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Verifier *VerifierTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Verifier.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Verifier *VerifierSession) AcceptOwnership() (*types.Transaction, error) {
	return _Verifier.Contract.AcceptOwnership(&_Verifier.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Verifier *VerifierTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Verifier.Contract.AcceptOwnership(&_Verifier.TransactOpts)
}

// ActivateConfig is a paid mutator transaction binding the contract method 0x0d1d79af.
//
// Solidity: function activateConfig(bytes32 configDigest) returns()
func (_Verifier *VerifierTransactor) ActivateConfig(opts *bind.TransactOpts, configDigest [32]byte) (*types.Transaction, error) {
	return _Verifier.contract.Transact(opts, "activateConfig", configDigest)
}

// ActivateConfig is a paid mutator transaction binding the contract method 0x0d1d79af.
//
// Solidity: function activateConfig(bytes32 configDigest) returns()
func (_Verifier *VerifierSession) ActivateConfig(configDigest [32]byte) (*types.Transaction, error) {
	return _Verifier.Contract.ActivateConfig(&_Verifier.TransactOpts, configDigest)
}

// ActivateConfig is a paid mutator transaction binding the contract method 0x0d1d79af.
//
// Solidity: function activateConfig(bytes32 configDigest) returns()
func (_Verifier *VerifierTransactorSession) ActivateConfig(configDigest [32]byte) (*types.Transaction, error) {
	return _Verifier.Contract.ActivateConfig(&_Verifier.TransactOpts, configDigest)
}

// DeactivateConfig is a paid mutator transaction binding the contract method 0x0f672ef4.
//
// Solidity: function deactivateConfig(bytes32 configDigest) returns()
func (_Verifier *VerifierTransactor) DeactivateConfig(opts *bind.TransactOpts, configDigest [32]byte) (*types.Transaction, error) {
	return _Verifier.contract.Transact(opts, "deactivateConfig", configDigest)
}

// DeactivateConfig is a paid mutator transaction binding the contract method 0x0f672ef4.
//
// Solidity: function deactivateConfig(bytes32 configDigest) returns()
func (_Verifier *VerifierSession) DeactivateConfig(configDigest [32]byte) (*types.Transaction, error) {
	return _Verifier.Contract.DeactivateConfig(&_Verifier.TransactOpts, configDigest)
}

// DeactivateConfig is a paid mutator transaction binding the contract method 0x0f672ef4.
//
// Solidity: function deactivateConfig(bytes32 configDigest) returns()
func (_Verifier *VerifierTransactorSession) DeactivateConfig(configDigest [32]byte) (*types.Transaction, error) {
	return _Verifier.Contract.DeactivateConfig(&_Verifier.TransactOpts, configDigest)
}

// SetConfig is a paid mutator transaction binding the contract method 0xe3d0e712.
//
// Solidity: function setConfig(address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig) returns()
func (_Verifier *VerifierTransactor) SetConfig(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _Verifier.contract.Transact(opts, "setConfig", signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

// SetConfig is a paid mutator transaction binding the contract method 0xe3d0e712.
//
// Solidity: function setConfig(address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig) returns()
func (_Verifier *VerifierSession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _Verifier.Contract.SetConfig(&_Verifier.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

// SetConfig is a paid mutator transaction binding the contract method 0xe3d0e712.
//
// Solidity: function setConfig(address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig) returns()
func (_Verifier *VerifierTransactorSession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _Verifier.Contract.SetConfig(&_Verifier.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_Verifier *VerifierTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _Verifier.contract.Transact(opts, "transferOwnership", to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_Verifier *VerifierSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _Verifier.Contract.TransferOwnership(&_Verifier.TransactOpts, to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_Verifier *VerifierTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _Verifier.Contract.TransferOwnership(&_Verifier.TransactOpts, to)
}

// Verify is a paid mutator transaction binding the contract method 0x3d3ac1b5.
//
// Solidity: function verify(bytes signedReport, address sender) returns(bytes response)
func (_Verifier *VerifierTransactor) Verify(opts *bind.TransactOpts, signedReport []byte, sender common.Address) (*types.Transaction, error) {
	return _Verifier.contract.Transact(opts, "verify", signedReport, sender)
}

// Verify is a paid mutator transaction binding the contract method 0x3d3ac1b5.
//
// Solidity: function verify(bytes signedReport, address sender) returns(bytes response)
func (_Verifier *VerifierSession) Verify(signedReport []byte, sender common.Address) (*types.Transaction, error) {
	return _Verifier.Contract.Verify(&_Verifier.TransactOpts, signedReport, sender)
}

// Verify is a paid mutator transaction binding the contract method 0x3d3ac1b5.
//
// Solidity: function verify(bytes signedReport, address sender) returns(bytes response)
func (_Verifier *VerifierTransactorSession) Verify(signedReport []byte, sender common.Address) (*types.Transaction, error) {
	return _Verifier.Contract.Verify(&_Verifier.TransactOpts, signedReport, sender)
}

// VerifierConfigActivatedIterator is returned from FilterConfigActivated and is used to iterate over the raw logs and unpacked data for ConfigActivated events raised by the Verifier contract.
type VerifierConfigActivatedIterator struct {
	Event *VerifierConfigActivated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *VerifierConfigActivatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifierConfigActivated)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(VerifierConfigActivated)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *VerifierConfigActivatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VerifierConfigActivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VerifierConfigActivated represents a ConfigActivated event raised by the Verifier contract.
type VerifierConfigActivated struct {
	ConfigDigest [32]byte
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterConfigActivated is a free log retrieval operation binding the contract event 0xa543797a0501218bba8a3daf75a71c8df8d1a7f791f4e44d40e43b6450183cea.
//
// Solidity: event ConfigActivated(bytes32 configDigest)
func (_Verifier *VerifierFilterer) FilterConfigActivated(opts *bind.FilterOpts) (*VerifierConfigActivatedIterator, error) {

	logs, sub, err := _Verifier.contract.FilterLogs(opts, "ConfigActivated")
	if err != nil {
		return nil, err
	}
	return &VerifierConfigActivatedIterator{contract: _Verifier.contract, event: "ConfigActivated", logs: logs, sub: sub}, nil
}

// WatchConfigActivated is a free log subscription operation binding the contract event 0xa543797a0501218bba8a3daf75a71c8df8d1a7f791f4e44d40e43b6450183cea.
//
// Solidity: event ConfigActivated(bytes32 configDigest)
func (_Verifier *VerifierFilterer) WatchConfigActivated(opts *bind.WatchOpts, sink chan<- *VerifierConfigActivated) (event.Subscription, error) {

	logs, sub, err := _Verifier.contract.WatchLogs(opts, "ConfigActivated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VerifierConfigActivated)
				if err := _Verifier.contract.UnpackLog(event, "ConfigActivated", log); err != nil {
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

// ParseConfigActivated is a log parse operation binding the contract event 0xa543797a0501218bba8a3daf75a71c8df8d1a7f791f4e44d40e43b6450183cea.
//
// Solidity: event ConfigActivated(bytes32 configDigest)
func (_Verifier *VerifierFilterer) ParseConfigActivated(log types.Log) (*VerifierConfigActivated, error) {
	event := new(VerifierConfigActivated)
	if err := _Verifier.contract.UnpackLog(event, "ConfigActivated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VerifierConfigDeactivatedIterator is returned from FilterConfigDeactivated and is used to iterate over the raw logs and unpacked data for ConfigDeactivated events raised by the Verifier contract.
type VerifierConfigDeactivatedIterator struct {
	Event *VerifierConfigDeactivated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *VerifierConfigDeactivatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifierConfigDeactivated)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(VerifierConfigDeactivated)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *VerifierConfigDeactivatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VerifierConfigDeactivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VerifierConfigDeactivated represents a ConfigDeactivated event raised by the Verifier contract.
type VerifierConfigDeactivated struct {
	ConfigDigest [32]byte
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterConfigDeactivated is a free log retrieval operation binding the contract event 0x5bfaab86edc1b932e3c334327a591c9ded067cb521abae19b95ca927d6076579.
//
// Solidity: event ConfigDeactivated(bytes32 configDigest)
func (_Verifier *VerifierFilterer) FilterConfigDeactivated(opts *bind.FilterOpts) (*VerifierConfigDeactivatedIterator, error) {

	logs, sub, err := _Verifier.contract.FilterLogs(opts, "ConfigDeactivated")
	if err != nil {
		return nil, err
	}
	return &VerifierConfigDeactivatedIterator{contract: _Verifier.contract, event: "ConfigDeactivated", logs: logs, sub: sub}, nil
}

// WatchConfigDeactivated is a free log subscription operation binding the contract event 0x5bfaab86edc1b932e3c334327a591c9ded067cb521abae19b95ca927d6076579.
//
// Solidity: event ConfigDeactivated(bytes32 configDigest)
func (_Verifier *VerifierFilterer) WatchConfigDeactivated(opts *bind.WatchOpts, sink chan<- *VerifierConfigDeactivated) (event.Subscription, error) {

	logs, sub, err := _Verifier.contract.WatchLogs(opts, "ConfigDeactivated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VerifierConfigDeactivated)
				if err := _Verifier.contract.UnpackLog(event, "ConfigDeactivated", log); err != nil {
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

// ParseConfigDeactivated is a log parse operation binding the contract event 0x5bfaab86edc1b932e3c334327a591c9ded067cb521abae19b95ca927d6076579.
//
// Solidity: event ConfigDeactivated(bytes32 configDigest)
func (_Verifier *VerifierFilterer) ParseConfigDeactivated(log types.Log) (*VerifierConfigDeactivated, error) {
	event := new(VerifierConfigDeactivated)
	if err := _Verifier.contract.UnpackLog(event, "ConfigDeactivated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VerifierConfigSetIterator is returned from FilterConfigSet and is used to iterate over the raw logs and unpacked data for ConfigSet events raised by the Verifier contract.
type VerifierConfigSetIterator struct {
	Event *VerifierConfigSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *VerifierConfigSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifierConfigSet)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(VerifierConfigSet)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *VerifierConfigSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VerifierConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VerifierConfigSet represents a ConfigSet event raised by the Verifier contract.
type VerifierConfigSet struct {
	PreviousConfigBlockNumber uint32
	ConfigDigest              [32]byte
	ConfigCount               uint64
	Signers                   []common.Address
	Transmitters              []common.Address
	F                         uint8
	OnchainConfig             []byte
	OffchainConfigVersion     uint64
	OffchainConfig            []byte
	Raw                       types.Log // Blockchain specific contextual infos
}

// FilterConfigSet is a free log retrieval operation binding the contract event 0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05.
//
// Solidity: event ConfigSet(uint32 previousConfigBlockNumber, bytes32 configDigest, uint64 configCount, address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig)
func (_Verifier *VerifierFilterer) FilterConfigSet(opts *bind.FilterOpts) (*VerifierConfigSetIterator, error) {

	logs, sub, err := _Verifier.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &VerifierConfigSetIterator{contract: _Verifier.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

// WatchConfigSet is a free log subscription operation binding the contract event 0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05.
//
// Solidity: event ConfigSet(uint32 previousConfigBlockNumber, bytes32 configDigest, uint64 configCount, address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig)
func (_Verifier *VerifierFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *VerifierConfigSet) (event.Subscription, error) {

	logs, sub, err := _Verifier.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VerifierConfigSet)
				if err := _Verifier.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

// ParseConfigSet is a log parse operation binding the contract event 0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05.
//
// Solidity: event ConfigSet(uint32 previousConfigBlockNumber, bytes32 configDigest, uint64 configCount, address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig)
func (_Verifier *VerifierFilterer) ParseConfigSet(log types.Log) (*VerifierConfigSet, error) {
	event := new(VerifierConfigSet)
	if err := _Verifier.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VerifierOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the Verifier contract.
type VerifierOwnershipTransferRequestedIterator struct {
	Event *VerifierOwnershipTransferRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *VerifierOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifierOwnershipTransferRequested)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(VerifierOwnershipTransferRequested)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *VerifierOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VerifierOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VerifierOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the Verifier contract.
type VerifierOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_Verifier *VerifierFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VerifierOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Verifier.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VerifierOwnershipTransferRequestedIterator{contract: _Verifier.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_Verifier *VerifierFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VerifierOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Verifier.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VerifierOwnershipTransferRequested)
				if err := _Verifier.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

// ParseOwnershipTransferRequested is a log parse operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_Verifier *VerifierFilterer) ParseOwnershipTransferRequested(log types.Log) (*VerifierOwnershipTransferRequested, error) {
	event := new(VerifierOwnershipTransferRequested)
	if err := _Verifier.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VerifierOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Verifier contract.
type VerifierOwnershipTransferredIterator struct {
	Event *VerifierOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *VerifierOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifierOwnershipTransferred)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(VerifierOwnershipTransferred)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *VerifierOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VerifierOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VerifierOwnershipTransferred represents a OwnershipTransferred event raised by the Verifier contract.
type VerifierOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_Verifier *VerifierFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VerifierOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Verifier.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VerifierOwnershipTransferredIterator{contract: _Verifier.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_Verifier *VerifierFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VerifierOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Verifier.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VerifierOwnershipTransferred)
				if err := _Verifier.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_Verifier *VerifierFilterer) ParseOwnershipTransferred(log types.Log) (*VerifierOwnershipTransferred, error) {
	event := new(VerifierOwnershipTransferred)
	if err := _Verifier.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VerifierReportVerifiedIterator is returned from FilterReportVerified and is used to iterate over the raw logs and unpacked data for ReportVerified events raised by the Verifier contract.
type VerifierReportVerifiedIterator struct {
	Event *VerifierReportVerified // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *VerifierReportVerifiedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifierReportVerified)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(VerifierReportVerified)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *VerifierReportVerifiedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VerifierReportVerifiedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VerifierReportVerified represents a ReportVerified event raised by the Verifier contract.
type VerifierReportVerified struct {
	FeedId     [32]byte
	ReportHash [32]byte
	Requester  common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterReportVerified is a free log retrieval operation binding the contract event 0x5a779ad0380d86bfcdb2748a39a349fe4fda58c933a2df958ba997284db8a269.
//
// Solidity: event ReportVerified(bytes32 feedId, bytes32 reportHash, address requester)
func (_Verifier *VerifierFilterer) FilterReportVerified(opts *bind.FilterOpts) (*VerifierReportVerifiedIterator, error) {

	logs, sub, err := _Verifier.contract.FilterLogs(opts, "ReportVerified")
	if err != nil {
		return nil, err
	}
	return &VerifierReportVerifiedIterator{contract: _Verifier.contract, event: "ReportVerified", logs: logs, sub: sub}, nil
}

// WatchReportVerified is a free log subscription operation binding the contract event 0x5a779ad0380d86bfcdb2748a39a349fe4fda58c933a2df958ba997284db8a269.
//
// Solidity: event ReportVerified(bytes32 feedId, bytes32 reportHash, address requester)
func (_Verifier *VerifierFilterer) WatchReportVerified(opts *bind.WatchOpts, sink chan<- *VerifierReportVerified) (event.Subscription, error) {

	logs, sub, err := _Verifier.contract.WatchLogs(opts, "ReportVerified")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VerifierReportVerified)
				if err := _Verifier.contract.UnpackLog(event, "ReportVerified", log); err != nil {
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

// ParseReportVerified is a log parse operation binding the contract event 0x5a779ad0380d86bfcdb2748a39a349fe4fda58c933a2df958ba997284db8a269.
//
// Solidity: event ReportVerified(bytes32 feedId, bytes32 reportHash, address requester)
func (_Verifier *VerifierFilterer) ParseReportVerified(log types.Log) (*VerifierReportVerified, error) {
	event := new(VerifierReportVerified)
	if err := _Verifier.contract.UnpackLog(event, "ReportVerified", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VerifierTransmittedIterator is returned from FilterTransmitted and is used to iterate over the raw logs and unpacked data for Transmitted events raised by the Verifier contract.
type VerifierTransmittedIterator struct {
	Event *VerifierTransmitted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *VerifierTransmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifierTransmitted)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(VerifierTransmitted)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *VerifierTransmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VerifierTransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VerifierTransmitted represents a Transmitted event raised by the Verifier contract.
type VerifierTransmitted struct {
	ConfigDigest [32]byte
	Epoch        uint32
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterTransmitted is a free log retrieval operation binding the contract event 0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62.
//
// Solidity: event Transmitted(bytes32 configDigest, uint32 epoch)
func (_Verifier *VerifierFilterer) FilterTransmitted(opts *bind.FilterOpts) (*VerifierTransmittedIterator, error) {

	logs, sub, err := _Verifier.contract.FilterLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return &VerifierTransmittedIterator{contract: _Verifier.contract, event: "Transmitted", logs: logs, sub: sub}, nil
}

// WatchTransmitted is a free log subscription operation binding the contract event 0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62.
//
// Solidity: event Transmitted(bytes32 configDigest, uint32 epoch)
func (_Verifier *VerifierFilterer) WatchTransmitted(opts *bind.WatchOpts, sink chan<- *VerifierTransmitted) (event.Subscription, error) {

	logs, sub, err := _Verifier.contract.WatchLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VerifierTransmitted)
				if err := _Verifier.contract.UnpackLog(event, "Transmitted", log); err != nil {
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

// ParseTransmitted is a log parse operation binding the contract event 0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62.
//
// Solidity: event Transmitted(bytes32 configDigest, uint32 epoch)
func (_Verifier *VerifierFilterer) ParseTransmitted(log types.Log) (*VerifierTransmitted, error) {
	event := new(VerifierTransmitted)
	if err := _Verifier.contract.UnpackLog(event, "Transmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
