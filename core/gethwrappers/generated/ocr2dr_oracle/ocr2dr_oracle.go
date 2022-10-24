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
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"donPublicKey\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"EmptyRequestData\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InconsistentReportData\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRequestID\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LowGasForConsumer\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"OracleRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"OracleResponse\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDONPublicKey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subscriptionId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"sendRequest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"_f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"_onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"_offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"_offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"transmitters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b5060405162002abb38038062002abb83398101604081905262000034916200023a565b600133806000816200008d5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c057620000c081620000e8565b505050151560f81b6080528051620000e090600890602084019062000194565b505062000369565b6001600160a01b038116331415620001435760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000084565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b828054620001a29062000316565b90600052602060002090601f016020900481019282620001c6576000855562000211565b82601f10620001e157805160ff191683800117855562000211565b8280016001018555821562000211579182015b8281111562000211578251825591602001919060010190620001f4565b506200021f92915062000223565b5090565b5b808211156200021f576000815560010162000224565b600060208083850312156200024e57600080fd5b82516001600160401b03808211156200026657600080fd5b818501915085601f8301126200027b57600080fd5b81518181111562000290576200029062000353565b604051601f8201601f19908116603f01168101908382118183101715620002bb57620002bb62000353565b816040528281528886848701011115620002d457600080fd5b600093505b82841015620002f85784840186015181850187015292850192620002d9565b828411156200030a5760008684830101525b98975050505050505050565b600181811c908216806200032b57607f821691505b602082108114156200034d57634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052604160045260246000fd5b60805160f81c6127336200038860003960006104c001526127336000f3fe608060405234801561001057600080fd5b50600436106100c95760003560e01c8063afcb95d711610081578063d328a91e1161005b578063d328a91e146101e4578063e3d0e712146101ec578063f2fde38b146101ff57600080fd5b8063afcb95d714610190578063b1dc65a4146101b0578063bb9fa3f5146101c357600080fd5b806381411834116100b2578063814118341461012357806381ff7048146101385780638da5cb5b1461016857600080fd5b8063181f5a77146100ce57806379ba509714610119575b600080fd5b60408051808201909152601281527f4f43523244524f7261636c6520302e302e30000000000000000000000000000060208201525b6040516101109190612256565b60405180910390f35b610121610212565b005b61012b610314565b60405161011091906121ba565b6004546002546040805163ffffffff80851682526401000000009094049093166020840152820152606001610110565b60005460405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610110565b604080516001815260006020820181905291810191909152606001610110565b6101216101be366004611f09565b610383565b6101d66101d13660046120a7565b610a2a565b604051908152602001610110565b610103610b99565b6101216101fa366004611e3c565b610c22565b61012161020d366004611e1a565b611605565b60015473ffffffffffffffffffffffffffffffffffffffff163314610298576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6060600780548060200260200160405190810160405280929190818152602001828054801561037957602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff16815260019091019060200180831161034e575b5050505050905090565b60005a604080516020601f8b018190048102820181019092528981529192508a3591818c0135916103d991849163ffffffff851691908e908e908190840183828082843760009201919091525061161992505050565b6040805183815262ffffff600884901c1660208201527fb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62910160405180910390a16040805160608101825260025480825260035460ff808216602085015261010090910416928201929092529083146104ae576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f636f6e666967446967657374206d69736d617463680000000000000000000000604482015260640161028f565b6104bc8b8b8b8b8b8b611702565b60007f000000000000000000000000000000000000000000000000000000000000000015610519576002826020015183604001516104fa919061248d565b61050491906124b2565b61050f90600161248d565b60ff16905061052f565b602082015161052990600161248d565b60ff1690505b888114610598576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f77726f6e67206e756d626572206f66207369676e617475726573000000000000604482015260640161028f565b888714610601576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f7369676e617475726573206f7574206f6620726567697374726174696f6e0000604482015260640161028f565b3360009081526005602090815260408083208151808301909252805460ff808216845292939192918401916101009091041660028111156106445761064461266a565b60028111156106555761065561266a565b90525090506002816020015160028111156106725761067261266a565b1480156106b957506007816000015160ff1681548110610694576106946126c8565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff1633145b61071f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f756e617574686f72697a6564207472616e736d69747465720000000000000000604482015260640161028f565b5050505050600088886040516107369291906121aa565b60405190819003812061074d918c9060200161218e565b60405160208183030381529060405280519060200120905061076d611b93565b604080518082019091526000808252602082015260005b88811015610a085760006001858884602081106107a3576107a36126c8565b6107b091901a601b61248d565b8d8d868181106107c2576107c26126c8565b905060200201358c8c878181106107db576107db6126c8565b9050602002013560405160008152602001604052604051610818949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa15801561083a573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff811660009081526005602090815290849020838501909452835460ff808216855292965092945084019161010090041660028111156108ba576108ba61266a565b60028111156108cb576108cb61266a565b90525092506001836020015160028111156108e8576108e861266a565b1461094f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f61646472657373206e6f7420617574686f72697a656420746f207369676e0000604482015260640161028f565b8251849060ff16601f8110610966576109666126c8565b6020020151156109d2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f6e6f6e2d756e69717565207369676e6174757265000000000000000000000000604482015260640161028f565b600184846000015160ff16601f81106109ed576109ed6126c8565b9115156020909202015250610a01816125d3565b9050610784565b5050505063ffffffff8110610a1f57610a1f61260c565b505050505050505050565b600081610a62576040517ec1cfc000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60098054906000610a72836125d3565b90915550506009546040517fffffffffffffffffffffffffffffffffffffffff0000000000000000000000003360601b1660208201526034810191909152600090605401604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0018152828252805160209182012083830183523384528184018981526000828152600a90935291839020935184547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff9091161784559051600190930192909255519091507f9dc59d1e6d6042f6c2d2af5a3f9f6502bf0abf7cf40f0eb80edf752ca396f2ba90610b89908390879087906121cd565b60405180910390a1949350505050565b606060088054610ba89061257f565b80601f0160208091040260200160405190810160405280929190818152602001828054610bd49061257f565b80156103795780601f10610bf657610100808354040283529160200191610379565b820191906000526020600020905b815481529060010190602001808311610c0457509395945050505050565b855185518560ff16601f831115610c95576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f746f6f206d616e79207369676e65727300000000000000000000000000000000604482015260640161028f565b60008111610cff576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f66206d75737420626520706f7369746976650000000000000000000000000000604482015260640161028f565b818314610d8d576040517f08c379a0000000000000000000000000000000000000000000000000000000008152602060048201526024808201527f6f7261636c6520616464726573736573206f7574206f6620726567697374726160448201527f74696f6e00000000000000000000000000000000000000000000000000000000606482015260840161028f565b610d988160036124fb565b8311610e00576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f6661756c74792d6f7261636c65206620746f6f20686967680000000000000000604482015260640161028f565b610e086117b0565b6040805160c0810182528a8152602081018a905260ff8916918101919091526060810187905267ffffffffffffffff8616608082015260a081018590525b60065415610ffb57600654600090610e6090600190612538565b9050600060068281548110610e7757610e776126c8565b60009182526020822001546007805473ffffffffffffffffffffffffffffffffffffffff90921693509084908110610eb157610eb16126c8565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff85811684526005909252604080842080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000090811690915592909116808452922080549091169055600680549192509080610f3157610f31612699565b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff00000000000000000000000000000000000000001690550190556007805480610f9a57610f9a612699565b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff000000000000000000000000000000000000000016905501905550610e46915050565b60005b8151518110156114605760006005600084600001518481518110611024576110246126c8565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff16825281019190915260400160002054610100900460ff16600281111561106e5761106e61266a565b146110d5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265706561746564207369676e65722061646472657373000000000000000000604482015260640161028f565b6040805180820190915260ff82168152600160208201528251805160059160009185908110611106576111066126c8565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040016000208251815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082168117835592840151919283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000016176101008360028111156111a7576111a761266a565b0217905550600091506111b79050565b60056000846020015184815181106111d1576111d16126c8565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff16825281019190915260400160002054610100900460ff16600281111561121b5761121b61266a565b14611282576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f7265706561746564207472616e736d6974746572206164647265737300000000604482015260640161028f565b6040805180820190915260ff8216815260208101600281525060056000846020015184815181106112b5576112b56126c8565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040016000208251815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082168117835592840151919283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000016176101008360028111156113565761135661266a565b021790555050825180516006925083908110611374576113746126c8565b602090810291909101810151825460018101845560009384529282902090920180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff90931692909217909155820151805160079190839081106113f0576113f06126c8565b60209081029190910181015182546001810184556000938452919092200180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909216919091179055611459816125d3565b9050610ffe565b506040810151600380547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff909216919091179055600480547fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff811664010000000063ffffffff4381168202928317855590830481169360019390926000926114f2928692908216911617612465565b92506101000a81548163ffffffff021916908363ffffffff1602179055506115514630600460009054906101000a900463ffffffff1663ffffffff16856000015186602001518760400151886060015189608001518a60a00151611833565b6002819055825180516003805460ff909216610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff90921691909117905560045460208501516040808701516060880151608089015160a08a015193517f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05986115f0988b98919763ffffffff90921696909591949193919261230e565b60405180910390a15050505050505050505050565b61160d6117b0565b611616816118de565b50565b6060806060838060200190518101906116329190611fc0565b8151835193965091945092501480159061164e57508051835114155b15611685576040517fe915fda500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b83518110156116f9576116e78482815181106116a6576116a66126c8565b60200260200101518483815181106116c0576116c06126c8565b60200260200101518484815181106116da576116da6126c8565b60200260200101516119d4565b806116f1816125d3565b915050611688565b50505050505050565b600061170f8260206124fb565b61171a8560206124fb565b6117268861014461244d565b611730919061244d565b61173a919061244d565b61174590600061244d565b90503681146116f9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f63616c6c64617461206c656e677468206d69736d617463680000000000000000604482015260640161028f565b60005473ffffffffffffffffffffffffffffffffffffffff163314611831576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161028f565b565b6000808a8a8a8a8a8a8a8a8a60405160200161185799989796959493929190612269565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179150509998505050505050505050565b73ffffffffffffffffffffffffffffffffffffffff811633141561195e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161028f565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000838152600a6020526040902054839073ffffffffffffffffffffffffffffffffffffffff16611a31576040517f803ed86300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000848152600a602090815260409182902054915186815273ffffffffffffffffffffffffffffffffffffffff909216917f9e9bc7616d42c2835d05ae617e508454e63b30b934be8aa932ebc125e0e58a64910160405180910390a162061a805a1015611aca576040517f566e3edd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040517f0ca7617500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821690630ca7617590611b2090889088908890600401612221565b600060405180830381600087803b158015611b3a57600080fd5b505af1158015611b4e573d6000803e3d6000fd5b50505060009586525050600a60205250506040822080547fffffffffffffffffffffffff00000000000000000000000000000000000000001681556001019190915550565b604051806103e00160405280601f906020820280368337509192915050565b803573ffffffffffffffffffffffffffffffffffffffff81168114611bd657600080fd5b919050565b600082601f830112611bec57600080fd5b81356020611c01611bfc836123e3565b612394565b80838252828201915082860187848660051b8901011115611c2157600080fd5b60005b85811015611c4757611c3582611bb2565b84529284019290840190600101611c24565b5090979650505050505050565b60008083601f840112611c6657600080fd5b50813567ffffffffffffffff811115611c7e57600080fd5b6020830191508360208260051b8501011115611c9957600080fd5b9250929050565b600082601f830112611cb157600080fd5b81516020611cc1611bfc836123e3565b80838252828201915082860187848660051b8901011115611ce157600080fd5b60005b85811015611c4757815167ffffffffffffffff811115611d0357600080fd5b8801603f81018a13611d1457600080fd5b858101516040611d26611bfc83612407565b8281528c82848601011115611d3a57600080fd5b611d49838a830184870161254f565b87525050509284019290840190600101611ce4565b60008083601f840112611d7057600080fd5b50813567ffffffffffffffff811115611d8857600080fd5b602083019150836020828501011115611c9957600080fd5b600082601f830112611db157600080fd5b8135611dbf611bfc82612407565b818152846020838601011115611dd457600080fd5b816020850160208301376000918101602001919091529392505050565b803567ffffffffffffffff81168114611bd657600080fd5b803560ff81168114611bd657600080fd5b600060208284031215611e2c57600080fd5b611e3582611bb2565b9392505050565b60008060008060008060c08789031215611e5557600080fd5b863567ffffffffffffffff80821115611e6d57600080fd5b611e798a838b01611bdb565b97506020890135915080821115611e8f57600080fd5b611e9b8a838b01611bdb565b9650611ea960408a01611e09565b95506060890135915080821115611ebf57600080fd5b611ecb8a838b01611da0565b9450611ed960808a01611df1565b935060a0890135915080821115611eef57600080fd5b50611efc89828a01611da0565b9150509295509295509295565b60008060008060008060008060e0898b031215611f2557600080fd5b606089018a811115611f3657600080fd5b8998503567ffffffffffffffff80821115611f5057600080fd5b611f5c8c838d01611d5e565b909950975060808b0135915080821115611f7557600080fd5b611f818c838d01611c54565b909750955060a08b0135915080821115611f9a57600080fd5b50611fa78b828c01611c54565b999c989b50969995989497949560c00135949350505050565b600080600060608486031215611fd557600080fd5b835167ffffffffffffffff80821115611fed57600080fd5b818601915086601f83011261200157600080fd5b81516020612011611bfc836123e3565b8083825282820191508286018b848660051b890101111561203157600080fd5b600096505b84871015612054578051835260019690960195918301918301612036565b509189015191975090935050508082111561206e57600080fd5b61207a87838801611ca0565b9350604086015191508082111561209057600080fd5b5061209d86828701611ca0565b9150509250925092565b6000806000604084860312156120bc57600080fd5b83359250602084013567ffffffffffffffff8111156120da57600080fd5b6120e686828701611d5e565b9497909650939450505050565b600081518084526020808501945080840160005b8381101561213957815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101612107565b509495945050505050565b6000815180845261215c81602086016020860161254f565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b8281526060826020830137600060809190910190815292915050565b8183823760009101908152919050565b602081526000611e3560208301846120f3565b83815260406020820152816040820152818360608301376000818301606090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016010192915050565b83815260606020820152600061223a6060830185612144565b828103604084015261224c8185612144565b9695505050505050565b602081526000611e356020830184612144565b60006101208b835273ffffffffffffffffffffffffffffffffffffffff8b16602084015267ffffffffffffffff808b1660408501528160608501526122b08285018b6120f3565b915083820360808501526122c4828a6120f3565b915060ff881660a085015283820360c08501526122e18288612144565b90861660e085015283810361010085015290506122fe8185612144565b9c9b505050505050505050505050565b600061012063ffffffff808d1684528b6020850152808b1660408501525080606084015261233e8184018a6120f3565b9050828103608084015261235281896120f3565b905060ff871660a084015282810360c084015261236f8187612144565b905067ffffffffffffffff851660e08401528281036101008401526122fe8185612144565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156123db576123db6126f7565b604052919050565b600067ffffffffffffffff8211156123fd576123fd6126f7565b5060051b60200190565b600067ffffffffffffffff821115612421576124216126f7565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082198211156124605761246061263b565b500190565b600063ffffffff8083168185168083038211156124845761248461263b565b01949350505050565b600060ff821660ff84168060ff038211156124aa576124aa61263b565b019392505050565b600060ff8316806124ec577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b8060ff84160491505092915050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156125335761253361263b565b500290565b60008282101561254a5761254a61263b565b500390565b60005b8381101561256a578181015183820152602001612552565b83811115612579576000848401525b50505050565b600181811c9082168061259357607f821691505b602082108114156125cd577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8214156126055761260561263b565b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052600160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var OCR2DROracleABI = OCR2DROracleMetaData.ABI

var OCR2DROracleBin = OCR2DROracleMetaData.Bin

func DeployOCR2DROracle(auth *bind.TransactOpts, backend bind.ContractBackend, donPublicKey []byte) (common.Address, *types.Transaction, *OCR2DROracle, error) {
	parsed, err := OCR2DROracleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OCR2DROracleBin), backend, donPublicKey)
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
