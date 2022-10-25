// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ocr2dr_oracle

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
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated"
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
)

var OCR2DROracleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"EmptyPublicKey\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptyRequestData\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InconsistentReportData\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRequestID\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LowGasForConsumer\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"OracleRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"OracleResponse\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDONPublicKey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subscriptionId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"sendRequest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"_f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"_onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"_offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"_offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"donPublicKey\",\"type\":\"bytes\"}],\"name\":\"setDONPublicKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"transmitters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b50600133806000816200006b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156200009e576200009e81620000af565b505050151560f81b6080526200015b565b6001600160a01b0381163314156200010a5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000062565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60805160f81c61289a6200017a600039600061052e015261289a6000f3fe608060405234801561001057600080fd5b50600436106100d45760003560e01c8063afcb95d711610081578063d328a91e1161005b578063d328a91e14610202578063e3d0e7121461020a578063f2fde38b1461021d57600080fd5b8063afcb95d7146101ae578063b1dc65a4146101ce578063bb9fa3f5146101e157600080fd5b806381411834116100b2578063814118341461014157806381ff7048146101565780638da5cb5b1461018657600080fd5b8063181f5a77146100d957806379ba5097146101245780637f15e1661461012e575b600080fd5b60408051808201909152601281527f4f43523244524f7261636c6520302e302e30000000000000000000000000000060208201525b60405161011b91906123bd565b60405180910390f35b61012c610230565b005b61012c61013c3660046121cc565b610332565b610149610382565b60405161011b9190612321565b6004546002546040805163ffffffff8085168252640100000000909404909316602084015282015260600161011b565b60005460405173ffffffffffffffffffffffffffffffffffffffff909116815260200161011b565b60408051600181526000602082018190529181019190915260600161011b565b61012c6101dc36600461202e565b6103f1565b6101f46101ef36600461220e565b610a98565b60405190815260200161011b565b61010e610c07565b61012c610218366004611f61565b610c90565b61012c61022b366004611f3f565b611673565b60015473ffffffffffffffffffffffffffffffffffffffff1633146102b6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b61033a611687565b80610371576040517f4f42be3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61037d60088383611c01565b505050565b606060078054806020026020016040519081016040528092919081815260200182805480156103e757602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116103bc575b5050505050905090565b60005a604080516020601f8b018190048102820181019092528981529192508a3591818c01359161044791849163ffffffff851691908e908e908190840183828082843760009201919091525061170a92505050565b6040805183815262ffffff600884901c1660208201527fb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62910160405180910390a16040805160608101825260025480825260035460ff8082166020850152610100909104169282019290925290831461051c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f636f6e666967446967657374206d69736d61746368000000000000000000000060448201526064016102ad565b61052a8b8b8b8b8b8b6117f3565b60007f0000000000000000000000000000000000000000000000000000000000000000156105875760028260200151836040015161056891906125f4565b6105729190612619565b61057d9060016125f4565b60ff16905061059d565b60208201516105979060016125f4565b60ff1690505b888114610606576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f77726f6e67206e756d626572206f66207369676e61747572657300000000000060448201526064016102ad565b88871461066f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f7369676e617475726573206f7574206f6620726567697374726174696f6e000060448201526064016102ad565b3360009081526005602090815260408083208151808301909252805460ff808216845292939192918401916101009091041660028111156106b2576106b26127d1565b60028111156106c3576106c36127d1565b90525090506002816020015160028111156106e0576106e06127d1565b14801561072757506007816000015160ff16815481106107025761070261282f565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff1633145b61078d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f756e617574686f72697a6564207472616e736d6974746572000000000000000060448201526064016102ad565b5050505050600088886040516107a4929190612311565b6040519081900381206107bb918c906020016122f5565b6040516020818303038152906040528051906020012090506107db611ca3565b604080518082019091526000808252602082015260005b88811015610a765760006001858884602081106108115761081161282f565b61081e91901a601b6125f4565b8d8d868181106108305761083061282f565b905060200201358c8c878181106108495761084961282f565b9050602002013560405160008152602001604052604051610886949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa1580156108a8573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff811660009081526005602090815290849020838501909452835460ff80821685529296509294508401916101009004166002811115610928576109286127d1565b6002811115610939576109396127d1565b9052509250600183602001516002811115610956576109566127d1565b146109bd576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f61646472657373206e6f7420617574686f72697a656420746f207369676e000060448201526064016102ad565b8251849060ff16601f81106109d4576109d461282f565b602002015115610a40576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f6e6f6e2d756e69717565207369676e617475726500000000000000000000000060448201526064016102ad565b600184846000015160ff16601f8110610a5b57610a5b61282f565b9115156020909202015250610a6f8161273a565b90506107f2565b5050505063ffffffff8110610a8d57610a8d612773565b505050505050505050565b600081610ad0576040517ec1cfc000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60098054906000610ae08361273a565b90915550506009546040517fffffffffffffffffffffffffffffffffffffffff0000000000000000000000003360601b1660208201526034810191909152600090605401604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0018152828252805160209182012083830183523384528184018981526000828152600a90935291839020935184547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff9091161784559051600190930192909255519091507f9dc59d1e6d6042f6c2d2af5a3f9f6502bf0abf7cf40f0eb80edf752ca396f2ba90610bf790839087908790612334565b60405180910390a1949350505050565b606060088054610c16906126e6565b80601f0160208091040260200160405190810160405280929190818152602001828054610c42906126e6565b80156103e75780601f10610c64576101008083540402835291602001916103e7565b820191906000526020600020905b815481529060010190602001808311610c7257509395945050505050565b855185518560ff16601f831115610d03576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f746f6f206d616e79207369676e6572730000000000000000000000000000000060448201526064016102ad565b60008111610d6d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f66206d75737420626520706f736974697665000000000000000000000000000060448201526064016102ad565b818314610dfb576040517f08c379a0000000000000000000000000000000000000000000000000000000008152602060048201526024808201527f6f7261636c6520616464726573736573206f7574206f6620726567697374726160448201527f74696f6e0000000000000000000000000000000000000000000000000000000060648201526084016102ad565b610e06816003612662565b8311610e6e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f6661756c74792d6f7261636c65206620746f6f2068696768000000000000000060448201526064016102ad565b610e76611687565b6040805160c0810182528a8152602081018a905260ff8916918101919091526060810187905267ffffffffffffffff8616608082015260a081018590525b6006541561106957600654600090610ece9060019061269f565b9050600060068281548110610ee557610ee561282f565b60009182526020822001546007805473ffffffffffffffffffffffffffffffffffffffff90921693509084908110610f1f57610f1f61282f565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff85811684526005909252604080842080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000090811690915592909116808452922080549091169055600680549192509080610f9f57610f9f612800565b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff0000000000000000000000000000000000000000169055019055600780548061100857611008612800565b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff000000000000000000000000000000000000000016905501905550610eb4915050565b60005b8151518110156114ce57600060056000846000015184815181106110925761109261282f565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff16825281019190915260400160002054610100900460ff1660028111156110dc576110dc6127d1565b14611143576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265706561746564207369676e6572206164647265737300000000000000000060448201526064016102ad565b6040805180820190915260ff821681526001602082015282518051600591600091859081106111745761117461282f565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040016000208251815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082168117835592840151919283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001617610100836002811115611215576112156127d1565b0217905550600091506112259050565b600560008460200151848151811061123f5761123f61282f565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff16825281019190915260400160002054610100900460ff166002811115611289576112896127d1565b146112f0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f7265706561746564207472616e736d697474657220616464726573730000000060448201526064016102ad565b6040805180820190915260ff8216815260208101600281525060056000846020015184815181106113235761132361282f565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040016000208251815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082168117835592840151919283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000016176101008360028111156113c4576113c46127d1565b0217905550508251805160069250839081106113e2576113e261282f565b602090810291909101810151825460018101845560009384529282902090920180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909316929092179091558201518051600791908390811061145e5761145e61282f565b60209081029190910181015182546001810184556000938452919092200180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff9092169190911790556114c78161273a565b905061106c565b506040810151600380547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff909216919091179055600480547fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff811664010000000063ffffffff4381168202928317855590830481169360019390926000926115609286929082169116176125cc565b92506101000a81548163ffffffff021916908363ffffffff1602179055506115bf4630600460009054906101000a900463ffffffff1663ffffffff16856000015186602001518760400151886060015189608001518a60a001516118a1565b6002819055825180516003805460ff909216610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff90921691909117905560045460208501516040808701516060880151608089015160a08a015193517f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e059861165e988b98919763ffffffff909216969095919491939192612475565b60405180910390a15050505050505050505050565b61167b611687565b6116848161194c565b50565b60005473ffffffffffffffffffffffffffffffffffffffff163314611708576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016102ad565b565b60608060608380602001905181019061172391906120e5565b8151835193965091945092501480159061173f57508051835114155b15611776576040517fe915fda500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b83518110156117ea576117d88482815181106117975761179761282f565b60200260200101518483815181106117b1576117b161282f565b60200260200101518484815181106117cb576117cb61282f565b6020026020010151611a42565b806117e28161273a565b915050611779565b50505050505050565b6000611800826020612662565b61180b856020612662565b611817886101446125b4565b61182191906125b4565b61182b91906125b4565b6118369060006125b4565b90503681146117ea576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f63616c6c64617461206c656e677468206d69736d61746368000000000000000060448201526064016102ad565b6000808a8a8a8a8a8a8a8a8a6040516020016118c5999897969594939291906123d0565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179150509998505050505050505050565b73ffffffffffffffffffffffffffffffffffffffff81163314156119cc576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016102ad565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000838152600a6020526040902054839073ffffffffffffffffffffffffffffffffffffffff16611a9f576040517f803ed86300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000848152600a602090815260409182902054915186815273ffffffffffffffffffffffffffffffffffffffff909216917f9e9bc7616d42c2835d05ae617e508454e63b30b934be8aa932ebc125e0e58a64910160405180910390a162061a805a1015611b38576040517f566e3edd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040517f0ca7617500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821690630ca7617590611b8e90889088908890600401612388565b600060405180830381600087803b158015611ba857600080fd5b505af1158015611bbc573d6000803e3d6000fd5b50505060009586525050600a60205250506040822080547fffffffffffffffffffffffff00000000000000000000000000000000000000001681556001019190915550565b828054611c0d906126e6565b90600052602060002090601f016020900481019282611c2f5760008555611c93565b82601f10611c66578280017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00823516178555611c93565b82800160010185558215611c93579182015b82811115611c93578235825591602001919060010190611c78565b50611c9f929150611cc2565b5090565b604051806103e00160405280601f906020820280368337509192915050565b5b80821115611c9f5760008155600101611cc3565b803573ffffffffffffffffffffffffffffffffffffffff81168114611cfb57600080fd5b919050565b600082601f830112611d1157600080fd5b81356020611d26611d218361254a565b6124fb565b80838252828201915082860187848660051b8901011115611d4657600080fd5b60005b85811015611d6c57611d5a82611cd7565b84529284019290840190600101611d49565b5090979650505050505050565b60008083601f840112611d8b57600080fd5b50813567ffffffffffffffff811115611da357600080fd5b6020830191508360208260051b8501011115611dbe57600080fd5b9250929050565b600082601f830112611dd657600080fd5b81516020611de6611d218361254a565b80838252828201915082860187848660051b8901011115611e0657600080fd5b60005b85811015611d6c57815167ffffffffffffffff811115611e2857600080fd5b8801603f81018a13611e3957600080fd5b858101516040611e4b611d218361256e565b8281528c82848601011115611e5f57600080fd5b611e6e838a83018487016126b6565b87525050509284019290840190600101611e09565b60008083601f840112611e9557600080fd5b50813567ffffffffffffffff811115611ead57600080fd5b602083019150836020828501011115611dbe57600080fd5b600082601f830112611ed657600080fd5b8135611ee4611d218261256e565b818152846020838601011115611ef957600080fd5b816020850160208301376000918101602001919091529392505050565b803567ffffffffffffffff81168114611cfb57600080fd5b803560ff81168114611cfb57600080fd5b600060208284031215611f5157600080fd5b611f5a82611cd7565b9392505050565b60008060008060008060c08789031215611f7a57600080fd5b863567ffffffffffffffff80821115611f9257600080fd5b611f9e8a838b01611d00565b97506020890135915080821115611fb457600080fd5b611fc08a838b01611d00565b9650611fce60408a01611f2e565b95506060890135915080821115611fe457600080fd5b611ff08a838b01611ec5565b9450611ffe60808a01611f16565b935060a089013591508082111561201457600080fd5b5061202189828a01611ec5565b9150509295509295509295565b60008060008060008060008060e0898b03121561204a57600080fd5b606089018a81111561205b57600080fd5b8998503567ffffffffffffffff8082111561207557600080fd5b6120818c838d01611e83565b909950975060808b013591508082111561209a57600080fd5b6120a68c838d01611d79565b909750955060a08b01359150808211156120bf57600080fd5b506120cc8b828c01611d79565b999c989b50969995989497949560c00135949350505050565b6000806000606084860312156120fa57600080fd5b835167ffffffffffffffff8082111561211257600080fd5b818601915086601f83011261212657600080fd5b81516020612136611d218361254a565b8083825282820191508286018b848660051b890101111561215657600080fd5b600096505b8487101561217957805183526001969096019591830191830161215b565b509189015191975090935050508082111561219357600080fd5b61219f87838801611dc5565b935060408601519150808211156121b557600080fd5b506121c286828701611dc5565b9150509250925092565b600080602083850312156121df57600080fd5b823567ffffffffffffffff8111156121f657600080fd5b61220285828601611e83565b90969095509350505050565b60008060006040848603121561222357600080fd5b83359250602084013567ffffffffffffffff81111561224157600080fd5b61224d86828701611e83565b9497909650939450505050565b600081518084526020808501945080840160005b838110156122a057815173ffffffffffffffffffffffffffffffffffffffff168752958201959082019060010161226e565b509495945050505050565b600081518084526122c38160208601602086016126b6565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b8281526060826020830137600060809190910190815292915050565b8183823760009101908152919050565b602081526000611f5a602083018461225a565b83815260406020820152816040820152818360608301376000818301606090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016010192915050565b8381526060602082015260006123a160608301856122ab565b82810360408401526123b381856122ab565b9695505050505050565b602081526000611f5a60208301846122ab565b60006101208b835273ffffffffffffffffffffffffffffffffffffffff8b16602084015267ffffffffffffffff808b1660408501528160608501526124178285018b61225a565b9150838203608085015261242b828a61225a565b915060ff881660a085015283820360c085015261244882886122ab565b90861660e0850152838103610100850152905061246581856122ab565b9c9b505050505050505050505050565b600061012063ffffffff808d1684528b6020850152808b166040850152508060608401526124a58184018a61225a565b905082810360808401526124b9818961225a565b905060ff871660a084015282810360c08401526124d681876122ab565b905067ffffffffffffffff851660e084015282810361010084015261246581856122ab565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156125425761254261285e565b604052919050565b600067ffffffffffffffff8211156125645761256461285e565b5060051b60200190565b600067ffffffffffffffff8211156125885761258861285e565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082198211156125c7576125c76127a2565b500190565b600063ffffffff8083168185168083038211156125eb576125eb6127a2565b01949350505050565b600060ff821660ff84168060ff03821115612611576126116127a2565b019392505050565b600060ff831680612653577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b8060ff84160491505092915050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048311821515161561269a5761269a6127a2565b500290565b6000828210156126b1576126b16127a2565b500390565b60005b838110156126d15781810151838201526020016126b9565b838111156126e0576000848401525b50505050565b600181811c908216806126fa57607f821691505b60208210811415612734577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82141561276c5761276c6127a2565b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052600160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var OCR2DROracleABI = OCR2DROracleMetaData.ABI

var OCR2DROracleBin = OCR2DROracleMetaData.Bin

func DeployOCR2DROracle(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *OCR2DROracle, error) {
	parsed, err := OCR2DROracleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OCR2DROracleBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OCR2DROracle{OCR2DROracleCaller: OCR2DROracleCaller{contract: contract}, OCR2DROracleTransactor: OCR2DROracleTransactor{contract: contract}, OCR2DROracleFilterer: OCR2DROracleFilterer{contract: contract}}, nil
}

type OCR2DROracle struct {
	address common.Address
	abi     abi.ABI
	OCR2DROracleCaller
	OCR2DROracleTransactor
	OCR2DROracleFilterer
}

type OCR2DROracleCaller struct {
	contract *bind.BoundContract
}

type OCR2DROracleTransactor struct {
	contract *bind.BoundContract
}

type OCR2DROracleFilterer struct {
	contract *bind.BoundContract
}

type OCR2DROracleSession struct {
	Contract     *OCR2DROracle
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type OCR2DROracleCallerSession struct {
	Contract *OCR2DROracleCaller
	CallOpts bind.CallOpts
}

type OCR2DROracleTransactorSession struct {
	Contract     *OCR2DROracleTransactor
	TransactOpts bind.TransactOpts
}

type OCR2DROracleRaw struct {
	Contract *OCR2DROracle
}

type OCR2DROracleCallerRaw struct {
	Contract *OCR2DROracleCaller
}

type OCR2DROracleTransactorRaw struct {
	Contract *OCR2DROracleTransactor
}

func NewOCR2DROracle(address common.Address, backend bind.ContractBackend) (*OCR2DROracle, error) {
	abi, err := abi.JSON(strings.NewReader(OCR2DROracleABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindOCR2DROracle(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OCR2DROracle{address: address, abi: abi, OCR2DROracleCaller: OCR2DROracleCaller{contract: contract}, OCR2DROracleTransactor: OCR2DROracleTransactor{contract: contract}, OCR2DROracleFilterer: OCR2DROracleFilterer{contract: contract}}, nil
}

func NewOCR2DROracleCaller(address common.Address, caller bind.ContractCaller) (*OCR2DROracleCaller, error) {
	contract, err := bindOCR2DROracle(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OCR2DROracleCaller{contract: contract}, nil
}

func NewOCR2DROracleTransactor(address common.Address, transactor bind.ContractTransactor) (*OCR2DROracleTransactor, error) {
	contract, err := bindOCR2DROracle(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OCR2DROracleTransactor{contract: contract}, nil
}

func NewOCR2DROracleFilterer(address common.Address, filterer bind.ContractFilterer) (*OCR2DROracleFilterer, error) {
	contract, err := bindOCR2DROracle(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OCR2DROracleFilterer{contract: contract}, nil
}

func bindOCR2DROracle(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OCR2DROracleABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_OCR2DROracle *OCR2DROracleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR2DROracle.Contract.OCR2DROracleCaller.contract.Call(opts, result, method, params...)
}

func (_OCR2DROracle *OCR2DROracleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2DROracle.Contract.OCR2DROracleTransactor.contract.Transfer(opts)
}

func (_OCR2DROracle *OCR2DROracleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR2DROracle.Contract.OCR2DROracleTransactor.contract.Transact(opts, method, params...)
}

func (_OCR2DROracle *OCR2DROracleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR2DROracle.Contract.contract.Call(opts, result, method, params...)
}

func (_OCR2DROracle *OCR2DROracleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2DROracle.Contract.contract.Transfer(opts)
}

func (_OCR2DROracle *OCR2DROracleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR2DROracle.Contract.contract.Transact(opts, method, params...)
}

func (_OCR2DROracle *OCR2DROracleCaller) GetDONPublicKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _OCR2DROracle.contract.Call(opts, &out, "getDONPublicKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_OCR2DROracle *OCR2DROracleSession) GetDONPublicKey() ([]byte, error) {
	return _OCR2DROracle.Contract.GetDONPublicKey(&_OCR2DROracle.CallOpts)
}

func (_OCR2DROracle *OCR2DROracleCallerSession) GetDONPublicKey() ([]byte, error) {
	return _OCR2DROracle.Contract.GetDONPublicKey(&_OCR2DROracle.CallOpts)
}

func (_OCR2DROracle *OCR2DROracleCaller) LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

	error) {
	var out []interface{}
	err := _OCR2DROracle.contract.Call(opts, &out, "latestConfigDetails")

	outstruct := new(LatestConfigDetails)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ConfigDigest = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

func (_OCR2DROracle *OCR2DROracleSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _OCR2DROracle.Contract.LatestConfigDetails(&_OCR2DROracle.CallOpts)
}

func (_OCR2DROracle *OCR2DROracleCallerSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _OCR2DROracle.Contract.LatestConfigDetails(&_OCR2DROracle.CallOpts)
}

func (_OCR2DROracle *OCR2DROracleCaller) LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

	error) {
	var out []interface{}
	err := _OCR2DROracle.contract.Call(opts, &out, "latestConfigDigestAndEpoch")

	outstruct := new(LatestConfigDigestAndEpoch)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ScanLogs = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfigDigest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Epoch = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

func (_OCR2DROracle *OCR2DROracleSession) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _OCR2DROracle.Contract.LatestConfigDigestAndEpoch(&_OCR2DROracle.CallOpts)
}

func (_OCR2DROracle *OCR2DROracleCallerSession) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _OCR2DROracle.Contract.LatestConfigDigestAndEpoch(&_OCR2DROracle.CallOpts)
}

func (_OCR2DROracle *OCR2DROracleCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OCR2DROracle.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_OCR2DROracle *OCR2DROracleSession) Owner() (common.Address, error) {
	return _OCR2DROracle.Contract.Owner(&_OCR2DROracle.CallOpts)
}

func (_OCR2DROracle *OCR2DROracleCallerSession) Owner() (common.Address, error) {
	return _OCR2DROracle.Contract.Owner(&_OCR2DROracle.CallOpts)
}

func (_OCR2DROracle *OCR2DROracleCaller) Transmitters(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _OCR2DROracle.contract.Call(opts, &out, "transmitters")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_OCR2DROracle *OCR2DROracleSession) Transmitters() ([]common.Address, error) {
	return _OCR2DROracle.Contract.Transmitters(&_OCR2DROracle.CallOpts)
}

func (_OCR2DROracle *OCR2DROracleCallerSession) Transmitters() ([]common.Address, error) {
	return _OCR2DROracle.Contract.Transmitters(&_OCR2DROracle.CallOpts)
}

func (_OCR2DROracle *OCR2DROracleCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _OCR2DROracle.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_OCR2DROracle *OCR2DROracleSession) TypeAndVersion() (string, error) {
	return _OCR2DROracle.Contract.TypeAndVersion(&_OCR2DROracle.CallOpts)
}

func (_OCR2DROracle *OCR2DROracleCallerSession) TypeAndVersion() (string, error) {
	return _OCR2DROracle.Contract.TypeAndVersion(&_OCR2DROracle.CallOpts)
}

func (_OCR2DROracle *OCR2DROracleTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2DROracle.contract.Transact(opts, "acceptOwnership")
}

func (_OCR2DROracle *OCR2DROracleSession) AcceptOwnership() (*types.Transaction, error) {
	return _OCR2DROracle.Contract.AcceptOwnership(&_OCR2DROracle.TransactOpts)
}

func (_OCR2DROracle *OCR2DROracleTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _OCR2DROracle.Contract.AcceptOwnership(&_OCR2DROracle.TransactOpts)
}

func (_OCR2DROracle *OCR2DROracleTransactor) SendRequest(opts *bind.TransactOpts, subscriptionId *big.Int, data []byte) (*types.Transaction, error) {
	return _OCR2DROracle.contract.Transact(opts, "sendRequest", subscriptionId, data)
}

func (_OCR2DROracle *OCR2DROracleSession) SendRequest(subscriptionId *big.Int, data []byte) (*types.Transaction, error) {
	return _OCR2DROracle.Contract.SendRequest(&_OCR2DROracle.TransactOpts, subscriptionId, data)
}

func (_OCR2DROracle *OCR2DROracleTransactorSession) SendRequest(subscriptionId *big.Int, data []byte) (*types.Transaction, error) {
	return _OCR2DROracle.Contract.SendRequest(&_OCR2DROracle.TransactOpts, subscriptionId, data)
}

func (_OCR2DROracle *OCR2DROracleTransactor) SetConfig(opts *bind.TransactOpts, _signers []common.Address, _transmitters []common.Address, _f uint8, _onchainConfig []byte, _offchainConfigVersion uint64, _offchainConfig []byte) (*types.Transaction, error) {
	return _OCR2DROracle.contract.Transact(opts, "setConfig", _signers, _transmitters, _f, _onchainConfig, _offchainConfigVersion, _offchainConfig)
}

func (_OCR2DROracle *OCR2DROracleSession) SetConfig(_signers []common.Address, _transmitters []common.Address, _f uint8, _onchainConfig []byte, _offchainConfigVersion uint64, _offchainConfig []byte) (*types.Transaction, error) {
	return _OCR2DROracle.Contract.SetConfig(&_OCR2DROracle.TransactOpts, _signers, _transmitters, _f, _onchainConfig, _offchainConfigVersion, _offchainConfig)
}

func (_OCR2DROracle *OCR2DROracleTransactorSession) SetConfig(_signers []common.Address, _transmitters []common.Address, _f uint8, _onchainConfig []byte, _offchainConfigVersion uint64, _offchainConfig []byte) (*types.Transaction, error) {
	return _OCR2DROracle.Contract.SetConfig(&_OCR2DROracle.TransactOpts, _signers, _transmitters, _f, _onchainConfig, _offchainConfigVersion, _offchainConfig)
}

func (_OCR2DROracle *OCR2DROracleTransactor) SetDONPublicKey(opts *bind.TransactOpts, donPublicKey []byte) (*types.Transaction, error) {
	return _OCR2DROracle.contract.Transact(opts, "setDONPublicKey", donPublicKey)
}

func (_OCR2DROracle *OCR2DROracleSession) SetDONPublicKey(donPublicKey []byte) (*types.Transaction, error) {
	return _OCR2DROracle.Contract.SetDONPublicKey(&_OCR2DROracle.TransactOpts, donPublicKey)
}

func (_OCR2DROracle *OCR2DROracleTransactorSession) SetDONPublicKey(donPublicKey []byte) (*types.Transaction, error) {
	return _OCR2DROracle.Contract.SetDONPublicKey(&_OCR2DROracle.TransactOpts, donPublicKey)
}

func (_OCR2DROracle *OCR2DROracleTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _OCR2DROracle.contract.Transact(opts, "transferOwnership", to)
}

func (_OCR2DROracle *OCR2DROracleSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OCR2DROracle.Contract.TransferOwnership(&_OCR2DROracle.TransactOpts, to)
}

func (_OCR2DROracle *OCR2DROracleTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OCR2DROracle.Contract.TransferOwnership(&_OCR2DROracle.TransactOpts, to)
}

func (_OCR2DROracle *OCR2DROracleTransactor) Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OCR2DROracle.contract.Transact(opts, "transmit", reportContext, report, rs, ss, rawVs)
}

func (_OCR2DROracle *OCR2DROracleSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OCR2DROracle.Contract.Transmit(&_OCR2DROracle.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_OCR2DROracle *OCR2DROracleTransactorSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OCR2DROracle.Contract.Transmit(&_OCR2DROracle.TransactOpts, reportContext, report, rs, ss, rawVs)
}

type OCR2DROracleConfigSetIterator struct {
	Event *OCR2DROracleConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2DROracleConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2DROracleConfigSet)
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
		it.Event = new(OCR2DROracleConfigSet)
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

func (it *OCR2DROracleConfigSetIterator) Error() error {
	return it.fail
}

func (it *OCR2DROracleConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2DROracleConfigSet struct {
	PreviousConfigBlockNumber uint32
	ConfigDigest              [32]byte
	ConfigCount               uint64
	Signers                   []common.Address
	Transmitters              []common.Address
	F                         uint8
	OnchainConfig             []byte
	OffchainConfigVersion     uint64
	OffchainConfig            []byte
	Raw                       types.Log
}

func (_OCR2DROracle *OCR2DROracleFilterer) FilterConfigSet(opts *bind.FilterOpts) (*OCR2DROracleConfigSetIterator, error) {

	logs, sub, err := _OCR2DROracle.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &OCR2DROracleConfigSetIterator{contract: _OCR2DROracle.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_OCR2DROracle *OCR2DROracleFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *OCR2DROracleConfigSet) (event.Subscription, error) {

	logs, sub, err := _OCR2DROracle.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2DROracleConfigSet)
				if err := _OCR2DROracle.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_OCR2DROracle *OCR2DROracleFilterer) ParseConfigSet(log types.Log) (*OCR2DROracleConfigSet, error) {
	event := new(OCR2DROracleConfigSet)
	if err := _OCR2DROracle.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2DROracleOracleRequestIterator struct {
	Event *OCR2DROracleOracleRequest

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2DROracleOracleRequestIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2DROracleOracleRequest)
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
		it.Event = new(OCR2DROracleOracleRequest)
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

func (it *OCR2DROracleOracleRequestIterator) Error() error {
	return it.fail
}

func (it *OCR2DROracleOracleRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2DROracleOracleRequest struct {
	RequestId [32]byte
	Data      []byte
	Raw       types.Log
}

func (_OCR2DROracle *OCR2DROracleFilterer) FilterOracleRequest(opts *bind.FilterOpts) (*OCR2DROracleOracleRequestIterator, error) {

	logs, sub, err := _OCR2DROracle.contract.FilterLogs(opts, "OracleRequest")
	if err != nil {
		return nil, err
	}
	return &OCR2DROracleOracleRequestIterator{contract: _OCR2DROracle.contract, event: "OracleRequest", logs: logs, sub: sub}, nil
}

func (_OCR2DROracle *OCR2DROracleFilterer) WatchOracleRequest(opts *bind.WatchOpts, sink chan<- *OCR2DROracleOracleRequest) (event.Subscription, error) {

	logs, sub, err := _OCR2DROracle.contract.WatchLogs(opts, "OracleRequest")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2DROracleOracleRequest)
				if err := _OCR2DROracle.contract.UnpackLog(event, "OracleRequest", log); err != nil {
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

func (_OCR2DROracle *OCR2DROracleFilterer) ParseOracleRequest(log types.Log) (*OCR2DROracleOracleRequest, error) {
	event := new(OCR2DROracleOracleRequest)
	if err := _OCR2DROracle.contract.UnpackLog(event, "OracleRequest", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2DROracleOracleResponseIterator struct {
	Event *OCR2DROracleOracleResponse

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2DROracleOracleResponseIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2DROracleOracleResponse)
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
		it.Event = new(OCR2DROracleOracleResponse)
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

func (it *OCR2DROracleOracleResponseIterator) Error() error {
	return it.fail
}

func (it *OCR2DROracleOracleResponseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2DROracleOracleResponse struct {
	RequestId [32]byte
	Raw       types.Log
}

func (_OCR2DROracle *OCR2DROracleFilterer) FilterOracleResponse(opts *bind.FilterOpts) (*OCR2DROracleOracleResponseIterator, error) {

	logs, sub, err := _OCR2DROracle.contract.FilterLogs(opts, "OracleResponse")
	if err != nil {
		return nil, err
	}
	return &OCR2DROracleOracleResponseIterator{contract: _OCR2DROracle.contract, event: "OracleResponse", logs: logs, sub: sub}, nil
}

func (_OCR2DROracle *OCR2DROracleFilterer) WatchOracleResponse(opts *bind.WatchOpts, sink chan<- *OCR2DROracleOracleResponse) (event.Subscription, error) {

	logs, sub, err := _OCR2DROracle.contract.WatchLogs(opts, "OracleResponse")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2DROracleOracleResponse)
				if err := _OCR2DROracle.contract.UnpackLog(event, "OracleResponse", log); err != nil {
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

func (_OCR2DROracle *OCR2DROracleFilterer) ParseOracleResponse(log types.Log) (*OCR2DROracleOracleResponse, error) {
	event := new(OCR2DROracleOracleResponse)
	if err := _OCR2DROracle.contract.UnpackLog(event, "OracleResponse", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2DROracleOwnershipTransferRequestedIterator struct {
	Event *OCR2DROracleOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2DROracleOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2DROracleOwnershipTransferRequested)
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
		it.Event = new(OCR2DROracleOwnershipTransferRequested)
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

func (it *OCR2DROracleOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *OCR2DROracleOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2DROracleOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_OCR2DROracle *OCR2DROracleFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OCR2DROracleOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR2DROracle.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OCR2DROracleOwnershipTransferRequestedIterator{contract: _OCR2DROracle.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_OCR2DROracle *OCR2DROracleFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OCR2DROracleOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR2DROracle.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2DROracleOwnershipTransferRequested)
				if err := _OCR2DROracle.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_OCR2DROracle *OCR2DROracleFilterer) ParseOwnershipTransferRequested(log types.Log) (*OCR2DROracleOwnershipTransferRequested, error) {
	event := new(OCR2DROracleOwnershipTransferRequested)
	if err := _OCR2DROracle.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2DROracleOwnershipTransferredIterator struct {
	Event *OCR2DROracleOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2DROracleOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2DROracleOwnershipTransferred)
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
		it.Event = new(OCR2DROracleOwnershipTransferred)
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

func (it *OCR2DROracleOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *OCR2DROracleOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2DROracleOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_OCR2DROracle *OCR2DROracleFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OCR2DROracleOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR2DROracle.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OCR2DROracleOwnershipTransferredIterator{contract: _OCR2DROracle.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_OCR2DROracle *OCR2DROracleFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OCR2DROracleOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR2DROracle.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2DROracleOwnershipTransferred)
				if err := _OCR2DROracle.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_OCR2DROracle *OCR2DROracleFilterer) ParseOwnershipTransferred(log types.Log) (*OCR2DROracleOwnershipTransferred, error) {
	event := new(OCR2DROracleOwnershipTransferred)
	if err := _OCR2DROracle.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2DROracleTransmittedIterator struct {
	Event *OCR2DROracleTransmitted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2DROracleTransmittedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2DROracleTransmitted)
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
		it.Event = new(OCR2DROracleTransmitted)
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

func (it *OCR2DROracleTransmittedIterator) Error() error {
	return it.fail
}

func (it *OCR2DROracleTransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2DROracleTransmitted struct {
	ConfigDigest [32]byte
	Epoch        uint32
	Raw          types.Log
}

func (_OCR2DROracle *OCR2DROracleFilterer) FilterTransmitted(opts *bind.FilterOpts) (*OCR2DROracleTransmittedIterator, error) {

	logs, sub, err := _OCR2DROracle.contract.FilterLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return &OCR2DROracleTransmittedIterator{contract: _OCR2DROracle.contract, event: "Transmitted", logs: logs, sub: sub}, nil
}

func (_OCR2DROracle *OCR2DROracleFilterer) WatchTransmitted(opts *bind.WatchOpts, sink chan<- *OCR2DROracleTransmitted) (event.Subscription, error) {

	logs, sub, err := _OCR2DROracle.contract.WatchLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2DROracleTransmitted)
				if err := _OCR2DROracle.contract.UnpackLog(event, "Transmitted", log); err != nil {
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

func (_OCR2DROracle *OCR2DROracleFilterer) ParseTransmitted(log types.Log) (*OCR2DROracleTransmitted, error) {
	event := new(OCR2DROracleTransmitted)
	if err := _OCR2DROracle.contract.UnpackLog(event, "Transmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LatestConfigDetails struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}
type LatestConfigDigestAndEpoch struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}

func (_OCR2DROracle *OCR2DROracle) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _OCR2DROracle.abi.Events["ConfigSet"].ID:
		return _OCR2DROracle.ParseConfigSet(log)
	case _OCR2DROracle.abi.Events["OracleRequest"].ID:
		return _OCR2DROracle.ParseOracleRequest(log)
	case _OCR2DROracle.abi.Events["OracleResponse"].ID:
		return _OCR2DROracle.ParseOracleResponse(log)
	case _OCR2DROracle.abi.Events["OwnershipTransferRequested"].ID:
		return _OCR2DROracle.ParseOwnershipTransferRequested(log)
	case _OCR2DROracle.abi.Events["OwnershipTransferred"].ID:
		return _OCR2DROracle.ParseOwnershipTransferred(log)
	case _OCR2DROracle.abi.Events["Transmitted"].ID:
		return _OCR2DROracle.ParseTransmitted(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (OCR2DROracleConfigSet) Topic() common.Hash {
	return common.HexToHash("0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05")
}

func (OCR2DROracleOracleRequest) Topic() common.Hash {
	return common.HexToHash("0x9dc59d1e6d6042f6c2d2af5a3f9f6502bf0abf7cf40f0eb80edf752ca396f2ba")
}

func (OCR2DROracleOracleResponse) Topic() common.Hash {
	return common.HexToHash("0x9e9bc7616d42c2835d05ae617e508454e63b30b934be8aa932ebc125e0e58a64")
}

func (OCR2DROracleOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (OCR2DROracleOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (OCR2DROracleTransmitted) Topic() common.Hash {
	return common.HexToHash("0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62")
}

func (_OCR2DROracle *OCR2DROracle) Address() common.Address {
	return _OCR2DROracle.address
}

type OCR2DROracleInterface interface {
	GetDONPublicKey(opts *bind.CallOpts) ([]byte, error)

	LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

		error)

	LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

		error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	Transmitters(opts *bind.CallOpts) ([]common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	SendRequest(opts *bind.TransactOpts, subscriptionId *big.Int, data []byte) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, _signers []common.Address, _transmitters []common.Address, _f uint8, _onchainConfig []byte, _offchainConfigVersion uint64, _offchainConfig []byte) (*types.Transaction, error)

	SetDONPublicKey(opts *bind.TransactOpts, donPublicKey []byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error)

	FilterConfigSet(opts *bind.FilterOpts) (*OCR2DROracleConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *OCR2DROracleConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*OCR2DROracleConfigSet, error)

	FilterOracleRequest(opts *bind.FilterOpts) (*OCR2DROracleOracleRequestIterator, error)

	WatchOracleRequest(opts *bind.WatchOpts, sink chan<- *OCR2DROracleOracleRequest) (event.Subscription, error)

	ParseOracleRequest(log types.Log) (*OCR2DROracleOracleRequest, error)

	FilterOracleResponse(opts *bind.FilterOpts) (*OCR2DROracleOracleResponseIterator, error)

	WatchOracleResponse(opts *bind.WatchOpts, sink chan<- *OCR2DROracleOracleResponse) (event.Subscription, error)

	ParseOracleResponse(log types.Log) (*OCR2DROracleOracleResponse, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OCR2DROracleOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OCR2DROracleOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*OCR2DROracleOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OCR2DROracleOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OCR2DROracleOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*OCR2DROracleOwnershipTransferred, error)

	FilterTransmitted(opts *bind.FilterOpts) (*OCR2DROracleTransmittedIterator, error)

	WatchTransmitted(opts *bind.WatchOpts, sink chan<- *OCR2DROracleTransmitted) (event.Subscription, error)

	ParseTransmitted(log types.Log) (*OCR2DROracleTransmitted, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
