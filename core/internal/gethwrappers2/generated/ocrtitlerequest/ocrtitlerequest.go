// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ocrtitlerequest

import (
	"fmt"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated"
)

var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

const OCRTitleRequestABI = "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"threshold\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"encodedConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"encoded\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"title\",\"type\":\"string\"}],\"name\":\"TitleFulfillment\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"url\",\"type\":\"string\"}],\"name\":\"TitleRequest\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"fulfilled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"url\",\"type\":\"string\"}],\"name\":\"request\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"_threshold\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"_onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"_offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"_offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"transmitters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]"

var OCRTitleRequestBin = "0x60a06040523480156200001157600080fd5b50600080546001600160a01b031916331781556080526040805160608101909152602a8082526200004c9190620022d760208301396200009a565b620000706040518060800160405280604e815260200162002301604e91396200009a565b620000946040518060a00160405280606e81526020016200234f606e91396200009a565b6200017d565b805160208083019182206008805460408051808601949094528381018290528051808503820181526060850180835281519190960120600190920190925580845260808301918252855160a0840152855190947f37adadbbe0ac5130611b65b06c5e2cef03817b6563f93855718a80afca1402ef94869488949193919260c09091019180838360005b838110156200013d57818101518382015260200162000123565b50505050905090810190601f1680156200016b5780820380516001836020036101000a031916815260200191505b50935050505060405180910390a15050565b60805160f81c61213c6200019b60003980610a45525061213c6000f3fe608060405234801561001057600080fd5b50600436106100be5760003560e01c806381ff704811610076578063b1dc65a41161005b578063b1dc65a4146102d9578063e3d0e712146103f0578063f2fde38b14610642576100be565b806381ff7048146102795780638da5cb5b146102a8576100be565b80632c199889116100a75780632c1998891461017157806379ba5097146102195780638141183414610221576100be565b8063181f5a77146100c35780632aa91bfd14610140575b600080fd5b6100cb610675565b6040805160208082528351818301528351919283929083019185019080838360005b838110156101055781810151838201526020016100ed565b50505050905090810190601f1680156101325780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b61015d6004803603602081101561015657600080fd5b50356106ac565b604080519115158252519081900360200190f35b6102176004803603602081101561018757600080fd5b8101906020810181356401000000008111156101a257600080fd5b8201836020820111156101b457600080fd5b803590602001918460018302840111640100000000831117156101d657600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295506106c1945050505050565b005b6102176107a1565b6102296108a3565b60408051602080825283518183015283519192839290830191858101910280838360005b8381101561026557818101518382015260200161024d565b505050509050019250505060405180910390f35b610281610912565b6040805163ffffffff94851681529290931660208301528183015290519081900360600190f35b6102b061092e565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b610217600480360360e08110156102ef57600080fd5b81018160808101606082013564010000000081111561030d57600080fd5b82018360208201111561031f57600080fd5b8035906020019184600183028401116401000000008311171561034157600080fd5b91939092909160208101903564010000000081111561035f57600080fd5b82018360208201111561037157600080fd5b8035906020019184602083028401116401000000008311171561039357600080fd5b9193909290916020810190356401000000008111156103b157600080fd5b8201836020820111156103c357600080fd5b803590602001918460208302840111640100000000831117156103e557600080fd5b91935091503561094a565b610217600480360360c081101561040657600080fd5b81019060208101813564010000000081111561042157600080fd5b82018360208201111561043357600080fd5b8035906020019184602083028401116401000000008311171561045557600080fd5b91908080602002602001604051908101604052809392919081815260200183836020028082843760009201919091525092959493602081019350359150506401000000008111156104a557600080fd5b8201836020820111156104b757600080fd5b803590602001918460208302840111640100000000831117156104d957600080fd5b919080806020026020016040519081016040528093929190818152602001838360200280828437600092019190915250929560ff85351695909490935060408101925060200135905064010000000081111561053457600080fd5b82018360208201111561054657600080fd5b8035906020019184600183028401116401000000008311171561056857600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929567ffffffffffffffff8535169590949093506040810192506020013590506401000000008111156105cd57600080fd5b8201836020820111156105df57600080fd5b8035906020019184600183028401116401000000008311171561060157600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610f66945050505050565b6102176004803603602081101561065857600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16611ad5565b60408051808201909152601b81527f4f43525469746c655265717565737420312e302e302d616c7068610000000000602082015290565b60009081526009602052604090205460ff1690565b805160208083019182206008805460408051808601949094528381018290528051808503820181526060850180835281519190960120600190920190925580845260808301918252855160a0840152855190947f37adadbbe0ac5130611b65b06c5e2cef03817b6563f93855718a80afca1402ef94869488949193919260c09091019180838360005b8381101561076257818101518382015260200161074a565b50505050905090810190601f16801561078f5780820380516001836020036101000a031916815260200191505b50935050505060405180910390a15050565b60015473ffffffffffffffffffffffffffffffffffffffff16331461082757604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015290519081900360640190fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6060600780548060200260200160405190810160405280929190818152602001828054801561090857602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116108dd575b5050505050905090565b60045460025463ffffffff808316926401000000009004169192565b60005473ffffffffffffffffffffffffffffffffffffffff1681565b60005a604080516020601f8b018190048102820181019092528981529192508a3591818c01359161099a9184918491908e908e9081908401838280828437600092019190915250611bd192505050565b6040805160608101825260025480825260035460ff80821660208501526101009091041692820192909252908314610a3357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f636f6e666967446967657374206d69736d617463680000000000000000000000604482015290519081900360640190fd5b610a418b8b8b8b8b8b611e07565b60007f000000000000000000000000000000000000000000000000000000000000000015610a8e576002826020015183604001510160ff1681610a8057fe5b0460010160ff169050610a9c565b816020015160010160ff1690505b888114610b0a57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f77726f6e67206e756d626572206f66207369676e617475726573000000000000604482015290519081900360640190fd5b888714610b7857604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f7369676e617475726573206f7574206f6620726567697374726174696f6e0000604482015290519081900360640190fd5b3360009081526005602090815260408083208151808301909252805460ff80821684529293919291840191610100909104166002811115610bb557fe5b6002811115610bc057fe5b9052509050600281602001516002811115610bd757fe5b148015610c1857506007816000015160ff1681548110610bf357fe5b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff1633145b610c8357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f756e617574686f72697a6564207472616e736d69747465720000000000000000604482015290519081900360640190fd5b50505050506000888860405180838380828437808301925050509250505060405180910390208a60405160200180838152602001826003602002808284378083019250505092505050604051602081830303815290604052805190602001209050610cec6120d5565b610cf46120f4565b60005b88811015610f40576000600185888460208110610d1057fe5b1a601b018d8d86818110610d2057fe5b905060200201358c8c87818110610d3357fe5b9050602002013560405160008152602001604052604051808581526020018460ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa158015610d8e573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff811660009081526005602090815290849020838501909452835460ff80821685529296509294508401916101009004166002811115610e0857fe5b6002811115610e1357fe5b9052509250600183602001516002811115610e2a57fe5b14610e9657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f61646472657373206e6f7420617574686f72697a656420746f207369676e0000604482015290519081900360640190fd5b8251849060ff16601f8110610ea757fe5b602002015115610f1857604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f6e6f6e2d756e69717565207369676e6174757265000000000000000000000000604482015290519081900360640190fd5b600184846000015160ff16601f8110610f2d57fe5b9115156020909202015250600101610cf7565b5050505063ffffffff8110610f5157fe5b610f5b8133611e8d565b505050505050505050565b855185518560ff16601f831115610fde57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f746f6f206d616e79207369676e65727300000000000000000000000000000000604482015290519081900360640190fd5b6000811161104d57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f7468726573686f6c64206d75737420626520706f736974697665000000000000604482015290519081900360640190fd5b8183146110a5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602481526020018061210c6024913960400191505060405180910390fd5b80600302831161111657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f6661756c74792d6f7261636c65207468726573686f6c6420746f6f2068696768604482015290519081900360640190fd5b60005473ffffffffffffffffffffffffffffffffffffffff16331461119c57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b6040805160c0810182528a8152602081018a905260ff89169181018290526060810188905267ffffffffffffffff8716608082015260a08101869052906111e39088611e8d565b6006541561139257600680547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8101916000918390811061122057fe5b60009182526020822001546007805473ffffffffffffffffffffffffffffffffffffffff9092169350908490811061125457fe5b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff85811684526005909252604080842080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000908116909155929091168084529220805490911690556006805491925090806112ce57fe5b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff0000000000000000000000000000000000000000169055019055600780548061133157fe5b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff0000000000000000000000000000000000000000169055019055506111e3915050565b60005b8151518110156117bb57600060056000846000015184815181106113b557fe5b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff16825281019190915260400160002054610100900460ff1660028111156113f957fe5b1461146557604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265706561746564207369676e65722061646472657373000000000000000000604482015290519081900360640190fd5b6040805180820190915260ff8216815260016020820152825180516005916000918590811061149057fe5b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff16825281810192909252604001600020825181547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff9091161780825591830151909182907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1661010083600281111561152957fe5b0217905550600091506115399050565b600560008460200151848151811061154d57fe5b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff16825281019190915260400160002054610100900460ff16600281111561159157fe5b146115fd57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f7265706561746564207472616e736d6974746572206164647265737300000000604482015290519081900360640190fd5b6040805180820190915260ff82168152602081016002815250600560008460200151848151811061162a57fe5b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff16825281810192909252604001600020825181547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff9091161780825591830151909182907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101008360028111156116c357fe5b0217905550508251805160069250839081106116db57fe5b602090810291909101810151825460018101845560009384529282902090920180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909316929092179091558201518051600791908390811061175157fe5b60209081029190910181015182546001808201855560009485529290932090920180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff9093169290921790915501611395565b5060408101516003805460ff83167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00909116179055600480544363ffffffff9081166401000000009081027fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff84161780831660010183167fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000009091161793849055855160208701516060880151608089015160a08a015194909604851697469761188d9789973097921695949391611e91565b60026000018190555050816000015151600260010160016101000a81548160ff021916908360ff1602179055507f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e0581600260000154600460009054906101000a900463ffffffff16856000015186602001518760400151886060015189608001518a60a00151604051808a63ffffffff1681526020018981526020018863ffffffff16815260200180602001806020018760ff168152602001806020018667ffffffffffffffff1681526020018060200185810385528b818151815260200191508051906020019060200280838360005b8381101561199657818101518382015260200161197e565b5050505090500185810384528a818151815260200191508051906020019060200280838360005b838110156119d55781810151838201526020016119bd565b50505050905001858103835288818151815260200191508051906020019080838360005b83811015611a115781810151838201526020016119f9565b50505050905090810190601f168015611a3e5780820380516001836020036101000a031916815260200191505b50858103825286518152865160209182019188019080838360005b83811015611a71578181015183820152602001611a59565b50505050905090810190601f168015611a9e5780820380516001836020036101000a031916815260200191505b509d505050505050505050505050505060405180910390a1611ac882604001518360600151611e8d565b5050505050505050505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314611b5b57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600080828060200190516040811015611be957600080fd5b815160208301805160405192949293830192919084640100000000821115611c1057600080fd5b908301906020820185811115611c2557600080fd5b8251640100000000811182820188101715611c3f57600080fd5b82525081516020918201929091019080838360005b83811015611c6c578181015183820152602001611c54565b50505050905090810190601f168015611c995780820380516001836020036101000a031916815260200191505b506040908152600086815260096020522054949650929450505060ff909116159050611d2657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f616c72656164792066756c66696c6c6564000000000000000000000000000000604482015290519081900360640190fd5b600082815260096020908152604080832080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055805185815280830182815285519282019290925284517f7cc5a0960ca99cf39ef66b30fb0dbec840eb2cbbd2ecf40d13c78a10a47bb7639487948794926060850192918601918190849084905b83811015611dc5578181015183820152602001611dad565b50505050905090810190601f168015611df25780820380516001836020036101000a031916815260200191505b50935050505060405180910390a15050505050565b602083810286019082020161014401368114611e8457604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f63616c6c64617461206c656e677468206d69736d617463680000000000000000604482015290519081900360640190fd5b50505050505050565b5050565b6000808a8a8a8a8a8a8a8a8a604051602001808a81526020018973ffffffffffffffffffffffffffffffffffffffff1681526020018867ffffffffffffffff16815260200180602001806020018760ff168152602001806020018667ffffffffffffffff1681526020018060200185810385528b818151815260200191508051906020019060200280838360005b83811015611f37578181015183820152602001611f1f565b5050505090500185810384528a818151815260200191508051906020019060200280838360005b83811015611f76578181015183820152602001611f5e565b50505050905001858103835288818151815260200191508051906020019080838360005b83811015611fb2578181015183820152602001611f9a565b50505050905090810190601f168015611fdf5780820380516001836020036101000a031916815260200191505b50858103825286518152865160209182019188019080838360005b83811015612012578181015183820152602001611ffa565b50505050905090810190601f16801561203f5780820380516001836020036101000a031916815260200191505b50604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179f505050505050505050505050505050509998505050505050505050565b604051806103e00160405280601f906020820280368337509192915050565b60408051808201909152600080825260208201529056fe6f7261636c6520616464726573736573206f7574206f6620726567697374726174696f6ea164736f6c6343000706000a68747470733a2f2f626c6f672e636861696e2e6c696e6b2f776861742d69732d636861696e6c696e6b2f68747470733a2f2f7777772e636f696e6465736b2e636f6d2f6d61726b2d637562616e2d6261636b65642d6e66742d6d61726b6574706c6163652d6d696e7461626c652d7261697365732d31336d68747470733a2f2f7777772e626c6f6f6d626572672e636f6d2f6f70696e696f6e2f61727469636c65732f323032312d30362d32342f666964656c6974792d6d616e616765722d6f776e65642d67616d6573746f702d6275742d6c61636b65642d6469616d6f6e642d68616e6473"

func DeployOCRTitleRequest(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *OCRTitleRequest, error) {
	parsed, err := abi.JSON(strings.NewReader(OCRTitleRequestABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(OCRTitleRequestBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OCRTitleRequest{OCRTitleRequestCaller: OCRTitleRequestCaller{contract: contract}, OCRTitleRequestTransactor: OCRTitleRequestTransactor{contract: contract}, OCRTitleRequestFilterer: OCRTitleRequestFilterer{contract: contract}}, nil
}

type OCRTitleRequest struct {
	address common.Address
	abi     abi.ABI
	OCRTitleRequestCaller
	OCRTitleRequestTransactor
	OCRTitleRequestFilterer
}

type OCRTitleRequestCaller struct {
	contract *bind.BoundContract
}

type OCRTitleRequestTransactor struct {
	contract *bind.BoundContract
}

type OCRTitleRequestFilterer struct {
	contract *bind.BoundContract
}

type OCRTitleRequestSession struct {
	Contract     *OCRTitleRequest
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type OCRTitleRequestCallerSession struct {
	Contract *OCRTitleRequestCaller
	CallOpts bind.CallOpts
}

type OCRTitleRequestTransactorSession struct {
	Contract     *OCRTitleRequestTransactor
	TransactOpts bind.TransactOpts
}

type OCRTitleRequestRaw struct {
	Contract *OCRTitleRequest
}

type OCRTitleRequestCallerRaw struct {
	Contract *OCRTitleRequestCaller
}

type OCRTitleRequestTransactorRaw struct {
	Contract *OCRTitleRequestTransactor
}

func NewOCRTitleRequest(address common.Address, backend bind.ContractBackend) (*OCRTitleRequest, error) {
	abi, err := abi.JSON(strings.NewReader(OCRTitleRequestABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindOCRTitleRequest(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OCRTitleRequest{address: address, abi: abi, OCRTitleRequestCaller: OCRTitleRequestCaller{contract: contract}, OCRTitleRequestTransactor: OCRTitleRequestTransactor{contract: contract}, OCRTitleRequestFilterer: OCRTitleRequestFilterer{contract: contract}}, nil
}

func NewOCRTitleRequestCaller(address common.Address, caller bind.ContractCaller) (*OCRTitleRequestCaller, error) {
	contract, err := bindOCRTitleRequest(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OCRTitleRequestCaller{contract: contract}, nil
}

func NewOCRTitleRequestTransactor(address common.Address, transactor bind.ContractTransactor) (*OCRTitleRequestTransactor, error) {
	contract, err := bindOCRTitleRequest(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OCRTitleRequestTransactor{contract: contract}, nil
}

func NewOCRTitleRequestFilterer(address common.Address, filterer bind.ContractFilterer) (*OCRTitleRequestFilterer, error) {
	contract, err := bindOCRTitleRequest(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OCRTitleRequestFilterer{contract: contract}, nil
}

func bindOCRTitleRequest(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OCRTitleRequestABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_OCRTitleRequest *OCRTitleRequestRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCRTitleRequest.Contract.OCRTitleRequestCaller.contract.Call(opts, result, method, params...)
}

func (_OCRTitleRequest *OCRTitleRequestRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCRTitleRequest.Contract.OCRTitleRequestTransactor.contract.Transfer(opts)
}

func (_OCRTitleRequest *OCRTitleRequestRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCRTitleRequest.Contract.OCRTitleRequestTransactor.contract.Transact(opts, method, params...)
}

func (_OCRTitleRequest *OCRTitleRequestCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCRTitleRequest.Contract.contract.Call(opts, result, method, params...)
}

func (_OCRTitleRequest *OCRTitleRequestTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCRTitleRequest.Contract.contract.Transfer(opts)
}

func (_OCRTitleRequest *OCRTitleRequestTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCRTitleRequest.Contract.contract.Transact(opts, method, params...)
}

func (_OCRTitleRequest *OCRTitleRequestCaller) Fulfilled(opts *bind.CallOpts, requestId [32]byte) (bool, error) {
	var out []interface{}
	err := _OCRTitleRequest.contract.Call(opts, &out, "fulfilled", requestId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_OCRTitleRequest *OCRTitleRequestSession) Fulfilled(requestId [32]byte) (bool, error) {
	return _OCRTitleRequest.Contract.Fulfilled(&_OCRTitleRequest.CallOpts, requestId)
}

func (_OCRTitleRequest *OCRTitleRequestCallerSession) Fulfilled(requestId [32]byte) (bool, error) {
	return _OCRTitleRequest.Contract.Fulfilled(&_OCRTitleRequest.CallOpts, requestId)
}

func (_OCRTitleRequest *OCRTitleRequestCaller) LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

	error) {
	var out []interface{}
	err := _OCRTitleRequest.contract.Call(opts, &out, "latestConfigDetails")

	outstruct := new(LatestConfigDetails)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ConfigDigest = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

func (_OCRTitleRequest *OCRTitleRequestSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _OCRTitleRequest.Contract.LatestConfigDetails(&_OCRTitleRequest.CallOpts)
}

func (_OCRTitleRequest *OCRTitleRequestCallerSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _OCRTitleRequest.Contract.LatestConfigDetails(&_OCRTitleRequest.CallOpts)
}

func (_OCRTitleRequest *OCRTitleRequestCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OCRTitleRequest.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_OCRTitleRequest *OCRTitleRequestSession) Owner() (common.Address, error) {
	return _OCRTitleRequest.Contract.Owner(&_OCRTitleRequest.CallOpts)
}

func (_OCRTitleRequest *OCRTitleRequestCallerSession) Owner() (common.Address, error) {
	return _OCRTitleRequest.Contract.Owner(&_OCRTitleRequest.CallOpts)
}

func (_OCRTitleRequest *OCRTitleRequestCaller) Transmitters(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _OCRTitleRequest.contract.Call(opts, &out, "transmitters")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_OCRTitleRequest *OCRTitleRequestSession) Transmitters() ([]common.Address, error) {
	return _OCRTitleRequest.Contract.Transmitters(&_OCRTitleRequest.CallOpts)
}

func (_OCRTitleRequest *OCRTitleRequestCallerSession) Transmitters() ([]common.Address, error) {
	return _OCRTitleRequest.Contract.Transmitters(&_OCRTitleRequest.CallOpts)
}

func (_OCRTitleRequest *OCRTitleRequestCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _OCRTitleRequest.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_OCRTitleRequest *OCRTitleRequestSession) TypeAndVersion() (string, error) {
	return _OCRTitleRequest.Contract.TypeAndVersion(&_OCRTitleRequest.CallOpts)
}

func (_OCRTitleRequest *OCRTitleRequestCallerSession) TypeAndVersion() (string, error) {
	return _OCRTitleRequest.Contract.TypeAndVersion(&_OCRTitleRequest.CallOpts)
}

func (_OCRTitleRequest *OCRTitleRequestTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCRTitleRequest.contract.Transact(opts, "acceptOwnership")
}

func (_OCRTitleRequest *OCRTitleRequestSession) AcceptOwnership() (*types.Transaction, error) {
	return _OCRTitleRequest.Contract.AcceptOwnership(&_OCRTitleRequest.TransactOpts)
}

func (_OCRTitleRequest *OCRTitleRequestTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _OCRTitleRequest.Contract.AcceptOwnership(&_OCRTitleRequest.TransactOpts)
}

func (_OCRTitleRequest *OCRTitleRequestTransactor) Request(opts *bind.TransactOpts, url string) (*types.Transaction, error) {
	return _OCRTitleRequest.contract.Transact(opts, "request", url)
}

func (_OCRTitleRequest *OCRTitleRequestSession) Request(url string) (*types.Transaction, error) {
	return _OCRTitleRequest.Contract.Request(&_OCRTitleRequest.TransactOpts, url)
}

func (_OCRTitleRequest *OCRTitleRequestTransactorSession) Request(url string) (*types.Transaction, error) {
	return _OCRTitleRequest.Contract.Request(&_OCRTitleRequest.TransactOpts, url)
}

func (_OCRTitleRequest *OCRTitleRequestTransactor) SetConfig(opts *bind.TransactOpts, _signers []common.Address, _transmitters []common.Address, _threshold uint8, _onchainConfig []byte, _offchainConfigVersion uint64, _offchainConfig []byte) (*types.Transaction, error) {
	return _OCRTitleRequest.contract.Transact(opts, "setConfig", _signers, _transmitters, _threshold, _onchainConfig, _offchainConfigVersion, _offchainConfig)
}

func (_OCRTitleRequest *OCRTitleRequestSession) SetConfig(_signers []common.Address, _transmitters []common.Address, _threshold uint8, _onchainConfig []byte, _offchainConfigVersion uint64, _offchainConfig []byte) (*types.Transaction, error) {
	return _OCRTitleRequest.Contract.SetConfig(&_OCRTitleRequest.TransactOpts, _signers, _transmitters, _threshold, _onchainConfig, _offchainConfigVersion, _offchainConfig)
}

func (_OCRTitleRequest *OCRTitleRequestTransactorSession) SetConfig(_signers []common.Address, _transmitters []common.Address, _threshold uint8, _onchainConfig []byte, _offchainConfigVersion uint64, _offchainConfig []byte) (*types.Transaction, error) {
	return _OCRTitleRequest.Contract.SetConfig(&_OCRTitleRequest.TransactOpts, _signers, _transmitters, _threshold, _onchainConfig, _offchainConfigVersion, _offchainConfig)
}

func (_OCRTitleRequest *OCRTitleRequestTransactor) TransferOwnership(opts *bind.TransactOpts, _to common.Address) (*types.Transaction, error) {
	return _OCRTitleRequest.contract.Transact(opts, "transferOwnership", _to)
}

func (_OCRTitleRequest *OCRTitleRequestSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _OCRTitleRequest.Contract.TransferOwnership(&_OCRTitleRequest.TransactOpts, _to)
}

func (_OCRTitleRequest *OCRTitleRequestTransactorSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _OCRTitleRequest.Contract.TransferOwnership(&_OCRTitleRequest.TransactOpts, _to)
}

func (_OCRTitleRequest *OCRTitleRequestTransactor) Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OCRTitleRequest.contract.Transact(opts, "transmit", reportContext, report, rs, ss, rawVs)
}

func (_OCRTitleRequest *OCRTitleRequestSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OCRTitleRequest.Contract.Transmit(&_OCRTitleRequest.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_OCRTitleRequest *OCRTitleRequestTransactorSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OCRTitleRequest.Contract.Transmit(&_OCRTitleRequest.TransactOpts, reportContext, report, rs, ss, rawVs)
}

type OCRTitleRequestConfigSetIterator struct {
	Event *OCRTitleRequestConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCRTitleRequestConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCRTitleRequestConfigSet)
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
		it.Event = new(OCRTitleRequestConfigSet)
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

func (it *OCRTitleRequestConfigSetIterator) Error() error {
	return it.fail
}

func (it *OCRTitleRequestConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCRTitleRequestConfigSet struct {
	PreviousConfigBlockNumber uint32
	ConfigDigest              [32]byte
	ConfigCount               uint64
	Signers                   []common.Address
	Transmitters              []common.Address
	Threshold                 uint8
	OnchainConfig             []byte
	EncodedConfigVersion      uint64
	Encoded                   []byte
	Raw                       types.Log
}

func (_OCRTitleRequest *OCRTitleRequestFilterer) FilterConfigSet(opts *bind.FilterOpts) (*OCRTitleRequestConfigSetIterator, error) {

	logs, sub, err := _OCRTitleRequest.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &OCRTitleRequestConfigSetIterator{contract: _OCRTitleRequest.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_OCRTitleRequest *OCRTitleRequestFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *OCRTitleRequestConfigSet) (event.Subscription, error) {

	logs, sub, err := _OCRTitleRequest.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCRTitleRequestConfigSet)
				if err := _OCRTitleRequest.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_OCRTitleRequest *OCRTitleRequestFilterer) ParseConfigSet(log types.Log) (*OCRTitleRequestConfigSet, error) {
	event := new(OCRTitleRequestConfigSet)
	if err := _OCRTitleRequest.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCRTitleRequestOwnershipTransferRequestedIterator struct {
	Event *OCRTitleRequestOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCRTitleRequestOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCRTitleRequestOwnershipTransferRequested)
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
		it.Event = new(OCRTitleRequestOwnershipTransferRequested)
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

func (it *OCRTitleRequestOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *OCRTitleRequestOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCRTitleRequestOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_OCRTitleRequest *OCRTitleRequestFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OCRTitleRequestOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCRTitleRequest.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OCRTitleRequestOwnershipTransferRequestedIterator{contract: _OCRTitleRequest.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_OCRTitleRequest *OCRTitleRequestFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OCRTitleRequestOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCRTitleRequest.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCRTitleRequestOwnershipTransferRequested)
				if err := _OCRTitleRequest.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_OCRTitleRequest *OCRTitleRequestFilterer) ParseOwnershipTransferRequested(log types.Log) (*OCRTitleRequestOwnershipTransferRequested, error) {
	event := new(OCRTitleRequestOwnershipTransferRequested)
	if err := _OCRTitleRequest.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCRTitleRequestOwnershipTransferredIterator struct {
	Event *OCRTitleRequestOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCRTitleRequestOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCRTitleRequestOwnershipTransferred)
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
		it.Event = new(OCRTitleRequestOwnershipTransferred)
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

func (it *OCRTitleRequestOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *OCRTitleRequestOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCRTitleRequestOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_OCRTitleRequest *OCRTitleRequestFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OCRTitleRequestOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCRTitleRequest.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OCRTitleRequestOwnershipTransferredIterator{contract: _OCRTitleRequest.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_OCRTitleRequest *OCRTitleRequestFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OCRTitleRequestOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCRTitleRequest.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCRTitleRequestOwnershipTransferred)
				if err := _OCRTitleRequest.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_OCRTitleRequest *OCRTitleRequestFilterer) ParseOwnershipTransferred(log types.Log) (*OCRTitleRequestOwnershipTransferred, error) {
	event := new(OCRTitleRequestOwnershipTransferred)
	if err := _OCRTitleRequest.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCRTitleRequestTitleFulfillmentIterator struct {
	Event *OCRTitleRequestTitleFulfillment

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCRTitleRequestTitleFulfillmentIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCRTitleRequestTitleFulfillment)
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
		it.Event = new(OCRTitleRequestTitleFulfillment)
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

func (it *OCRTitleRequestTitleFulfillmentIterator) Error() error {
	return it.fail
}

func (it *OCRTitleRequestTitleFulfillmentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCRTitleRequestTitleFulfillment struct {
	RequestId [32]byte
	Title     string
	Raw       types.Log
}

func (_OCRTitleRequest *OCRTitleRequestFilterer) FilterTitleFulfillment(opts *bind.FilterOpts) (*OCRTitleRequestTitleFulfillmentIterator, error) {

	logs, sub, err := _OCRTitleRequest.contract.FilterLogs(opts, "TitleFulfillment")
	if err != nil {
		return nil, err
	}
	return &OCRTitleRequestTitleFulfillmentIterator{contract: _OCRTitleRequest.contract, event: "TitleFulfillment", logs: logs, sub: sub}, nil
}

func (_OCRTitleRequest *OCRTitleRequestFilterer) WatchTitleFulfillment(opts *bind.WatchOpts, sink chan<- *OCRTitleRequestTitleFulfillment) (event.Subscription, error) {

	logs, sub, err := _OCRTitleRequest.contract.WatchLogs(opts, "TitleFulfillment")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCRTitleRequestTitleFulfillment)
				if err := _OCRTitleRequest.contract.UnpackLog(event, "TitleFulfillment", log); err != nil {
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

func (_OCRTitleRequest *OCRTitleRequestFilterer) ParseTitleFulfillment(log types.Log) (*OCRTitleRequestTitleFulfillment, error) {
	event := new(OCRTitleRequestTitleFulfillment)
	if err := _OCRTitleRequest.contract.UnpackLog(event, "TitleFulfillment", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCRTitleRequestTitleRequestIterator struct {
	Event *OCRTitleRequestTitleRequest

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCRTitleRequestTitleRequestIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCRTitleRequestTitleRequest)
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
		it.Event = new(OCRTitleRequestTitleRequest)
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

func (it *OCRTitleRequestTitleRequestIterator) Error() error {
	return it.fail
}

func (it *OCRTitleRequestTitleRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCRTitleRequestTitleRequest struct {
	RequestId [32]byte
	Url       string
	Raw       types.Log
}

func (_OCRTitleRequest *OCRTitleRequestFilterer) FilterTitleRequest(opts *bind.FilterOpts) (*OCRTitleRequestTitleRequestIterator, error) {

	logs, sub, err := _OCRTitleRequest.contract.FilterLogs(opts, "TitleRequest")
	if err != nil {
		return nil, err
	}
	return &OCRTitleRequestTitleRequestIterator{contract: _OCRTitleRequest.contract, event: "TitleRequest", logs: logs, sub: sub}, nil
}

func (_OCRTitleRequest *OCRTitleRequestFilterer) WatchTitleRequest(opts *bind.WatchOpts, sink chan<- *OCRTitleRequestTitleRequest) (event.Subscription, error) {

	logs, sub, err := _OCRTitleRequest.contract.WatchLogs(opts, "TitleRequest")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCRTitleRequestTitleRequest)
				if err := _OCRTitleRequest.contract.UnpackLog(event, "TitleRequest", log); err != nil {
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

func (_OCRTitleRequest *OCRTitleRequestFilterer) ParseTitleRequest(log types.Log) (*OCRTitleRequestTitleRequest, error) {
	event := new(OCRTitleRequestTitleRequest)
	if err := _OCRTitleRequest.contract.UnpackLog(event, "TitleRequest", log); err != nil {
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

func (_OCRTitleRequest *OCRTitleRequest) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _OCRTitleRequest.abi.Events["ConfigSet"].ID:
		return _OCRTitleRequest.ParseConfigSet(log)
	case _OCRTitleRequest.abi.Events["OwnershipTransferRequested"].ID:
		return _OCRTitleRequest.ParseOwnershipTransferRequested(log)
	case _OCRTitleRequest.abi.Events["OwnershipTransferred"].ID:
		return _OCRTitleRequest.ParseOwnershipTransferred(log)
	case _OCRTitleRequest.abi.Events["TitleFulfillment"].ID:
		return _OCRTitleRequest.ParseTitleFulfillment(log)
	case _OCRTitleRequest.abi.Events["TitleRequest"].ID:
		return _OCRTitleRequest.ParseTitleRequest(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (OCRTitleRequestConfigSet) Topic() common.Hash {
	return common.HexToHash("0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05")
}

func (OCRTitleRequestOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (OCRTitleRequestOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (OCRTitleRequestTitleFulfillment) Topic() common.Hash {
	return common.HexToHash("0x7cc5a0960ca99cf39ef66b30fb0dbec840eb2cbbd2ecf40d13c78a10a47bb763")
}

func (OCRTitleRequestTitleRequest) Topic() common.Hash {
	return common.HexToHash("0x37adadbbe0ac5130611b65b06c5e2cef03817b6563f93855718a80afca1402ef")
}

func (_OCRTitleRequest *OCRTitleRequest) Address() common.Address {
	return _OCRTitleRequest.address
}

type OCRTitleRequestInterface interface {
	Fulfilled(opts *bind.CallOpts, requestId [32]byte) (bool, error)

	LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

		error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	Transmitters(opts *bind.CallOpts) ([]common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	Request(opts *bind.TransactOpts, url string) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, _signers []common.Address, _transmitters []common.Address, _threshold uint8, _onchainConfig []byte, _offchainConfigVersion uint64, _offchainConfig []byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, _to common.Address) (*types.Transaction, error)

	Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error)

	FilterConfigSet(opts *bind.FilterOpts) (*OCRTitleRequestConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *OCRTitleRequestConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*OCRTitleRequestConfigSet, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OCRTitleRequestOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OCRTitleRequestOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*OCRTitleRequestOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OCRTitleRequestOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OCRTitleRequestOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*OCRTitleRequestOwnershipTransferred, error)

	FilterTitleFulfillment(opts *bind.FilterOpts) (*OCRTitleRequestTitleFulfillmentIterator, error)

	WatchTitleFulfillment(opts *bind.WatchOpts, sink chan<- *OCRTitleRequestTitleFulfillment) (event.Subscription, error)

	ParseTitleFulfillment(log types.Log) (*OCRTitleRequestTitleFulfillment, error)

	FilterTitleRequest(opts *bind.FilterOpts) (*OCRTitleRequestTitleRequestIterator, error)

	WatchTitleRequest(opts *bind.WatchOpts, sink chan<- *OCRTitleRequestTitleRequest) (event.Subscription, error)

	ParseTitleRequest(log types.Log) (*OCRTitleRequestTitleRequest, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
