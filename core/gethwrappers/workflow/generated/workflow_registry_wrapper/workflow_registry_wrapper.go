// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package workflow_registry_wrapper

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

type WorkflowRegistryWorkflowMetadata struct {
	WorkflowID   [32]byte
	Owner        common.Address
	DonID        uint32
	Status       uint8
	WorkflowName string
	BinaryURL    string
	ConfigURL    string
	SecretsURL   string
}

var WorkflowRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"AddressNotAuthorized\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"CallerIsNotWorkflowOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotTransferToSelf\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donID\",\"type\":\"uint32\"}],\"name\":\"DONNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidWorkflowID\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeProposedOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OwnerCannotBeZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RegistryLocked\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"providedLength\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"maxAllowedLength\",\"type\":\"uint8\"}],\"name\":\"URLTooLong\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WorkflowAlreadyInDesiredStatus\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WorkflowAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WorkflowContentNotUpdated\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WorkflowDoesNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WorkflowIDAlreadyExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WorkflowIDNotUpdated\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"providedLength\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"maxAllowedLength\",\"type\":\"uint8\"}],\"name\":\"WorkflowNameTooLong\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32[]\",\"name\":\"donIDs\",\"type\":\"uint32[]\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"AllowedDONsUpdatedV1\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"addresses\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"AuthorizedAddressesUpdatedV1\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"lockedBy\",\"type\":\"address\"}],\"name\":\"RegistryLockedV1\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"unlockedBy\",\"type\":\"address\"}],\"name\":\"RegistryUnlockedV1\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"workflowID\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"workflowOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"donID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"workflowName\",\"type\":\"string\"}],\"name\":\"WorkflowActivatedV1\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"workflowID\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"workflowOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"donID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"workflowName\",\"type\":\"string\"}],\"name\":\"WorkflowDeletedV1\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"secretsURLHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"workflowName\",\"type\":\"string\"}],\"name\":\"WorkflowForceUpdateSecretsRequestedV1\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"workflowID\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"workflowOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"donID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"workflowName\",\"type\":\"string\"}],\"name\":\"WorkflowPausedV1\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"workflowID\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"workflowOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"donID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"enumWorkflowRegistry.WorkflowStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"workflowName\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"binaryURL\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"configURL\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"secretsURL\",\"type\":\"string\"}],\"name\":\"WorkflowRegisteredV1\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"oldWorkflowID\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"workflowOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"donID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"newWorkflowID\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"workflowName\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"binaryURL\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"configURL\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"secretsURL\",\"type\":\"string\"}],\"name\":\"WorkflowUpdatedV1\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"workflowKey\",\"type\":\"bytes32\"}],\"name\":\"activateWorkflow\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"field\",\"type\":\"string\"}],\"name\":\"computeHashKey\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"workflowKey\",\"type\":\"bytes32\"}],\"name\":\"deleteWorkflow\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllAllowedDONs\",\"outputs\":[{\"internalType\":\"uint32[]\",\"name\":\"allowedDONs\",\"type\":\"uint32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllAuthorizedAddresses\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"authorizedAddresses\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"workflowOwner\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"workflowName\",\"type\":\"string\"}],\"name\":\"getWorkflowMetadata\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"workflowID\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"donID\",\"type\":\"uint32\"},{\"internalType\":\"enumWorkflowRegistry.WorkflowStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"workflowName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"binaryURL\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"configURL\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"secretsURL\",\"type\":\"string\"}],\"internalType\":\"structWorkflowRegistry.WorkflowMetadata\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donID\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"start\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"limit\",\"type\":\"uint256\"}],\"name\":\"getWorkflowMetadataListByDON\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"workflowID\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"donID\",\"type\":\"uint32\"},{\"internalType\":\"enumWorkflowRegistry.WorkflowStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"workflowName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"binaryURL\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"configURL\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"secretsURL\",\"type\":\"string\"}],\"internalType\":\"structWorkflowRegistry.WorkflowMetadata[]\",\"name\":\"workflowMetadataList\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"workflowOwner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"start\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"limit\",\"type\":\"uint256\"}],\"name\":\"getWorkflowMetadataListByOwner\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"workflowID\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"donID\",\"type\":\"uint32\"},{\"internalType\":\"enumWorkflowRegistry.WorkflowStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"workflowName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"binaryURL\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"configURL\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"secretsURL\",\"type\":\"string\"}],\"internalType\":\"structWorkflowRegistry.WorkflowMetadata[]\",\"name\":\"workflowMetadataList\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isRegistryLocked\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lockRegistry\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"workflowKey\",\"type\":\"bytes32\"}],\"name\":\"pauseWorkflow\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"workflowName\",\"type\":\"string\"},{\"internalType\":\"bytes32\",\"name\":\"workflowID\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"donID\",\"type\":\"uint32\"},{\"internalType\":\"enumWorkflowRegistry.WorkflowStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"binaryURL\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"configURL\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"secretsURL\",\"type\":\"string\"}],\"name\":\"registerWorkflow\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"secretsURL\",\"type\":\"string\"}],\"name\":\"requestForceUpdateSecrets\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unlockRegistry\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32[]\",\"name\":\"donIDs\",\"type\":\"uint32[]\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"updateAllowedDONs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"addresses\",\"type\":\"address[]\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"updateAuthorizedAddresses\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"workflowKey\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newWorkflowID\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"binaryURL\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"configURL\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"secretsURL\",\"type\":\"string\"}],\"name\":\"updateWorkflow\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080806040523461004a57331561003b57600180546001600160a01b03191633179055600a805460ff191690556040516133df90816100508239f35b639b15e16f60e01b8152600490fd5b600080fdfe6080604052600436101561001257600080fd5b60003560e01c806308e7f63a14612082578063181f5a7714611ff35780632303348a14611eb65780632b596f6d14611e2c5780633ccd14ff146114f2578063695e1340146113565780636f3517711461127d578063724c13dd146111865780637497066b1461106b57806379ba509714610f955780637ec0846d14610f0e5780638da5cb5b14610ebc5780639f4cb53414610e9b578063b87a019414610e45578063d4b89c7414610698578063db800092146105fd578063e3dce080146104d6578063e690f33214610362578063f2fde38b14610284578063f794bdeb146101495763f99ecb6b1461010357600080fd5b346101445760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261014457602060ff600a54166040519015158152f35b600080fd5b346101445760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101445760068054610185816123fc565b6101926040519182612283565b81815261019e826123fc565b916020937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe060208401940136853760005b82811061023257505050906040519283926020840190602085525180915260408401929160005b82811061020557505050500390f35b835173ffffffffffffffffffffffffffffffffffffffff16855286955093810193928101926001016101f6565b6001908260005273ffffffffffffffffffffffffffffffffffffffff817ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d3f01541661027d828761252e565b52016101cf565b346101445760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610144576102bb61236a565b6102c3612bc7565b73ffffffffffffffffffffffffffffffffffffffff8091169033821461033857817fffffffffffffffffffffffff00000000000000000000000000000000000000006000541617600055600154167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278600080a3005b60046040517fdad89dca000000000000000000000000000000000000000000000000000000008152fd5b346101445760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101445760ff600a54166104ac576103a760043533612dad565b600181019081549160ff8360c01c16600281101561047d576001146104535778010000000000000000000000000000000000000000000000007fffffffffffffff00ffffffffffffffffffffffffffffffffffffffffffffffff841617905580547f6a0ed88e9cf3cb493ab4028fcb1dc7d18f0130fcdfba096edde0aadbfbf5e99f63ffffffff604051946020865260a01c16938061044e339560026020840191016125d0565b0390a4005b60046040517f6f861db1000000000000000000000000000000000000000000000000000000008152fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b60046040517f78a4e7d9000000000000000000000000000000000000000000000000000000008152fd5b34610144576104e4366122f2565b916104ed612bc7565b60ff600a54166104ac5760005b828110610589575060405191806040840160408552526060830191906000905b8082106105515785151560208601527f509460cccbb176edde6cac28895a4415a24961b8f3a0bd2617b9bb7b4e166c9b85850386a1005b90919283359073ffffffffffffffffffffffffffffffffffffffff82168092036101445760019181526020809101940192019061051a565b60019084156105cb576105c373ffffffffffffffffffffffffffffffffffffffff6105bd6105b8848888612a67565b612ba6565b16612f88565b505b016104fa565b6105f773ffffffffffffffffffffffffffffffffffffffff6105f16105b8848888612a67565b166131b9565b506105c5565b346101445761061d61060e3661238d565b91610617612414565b50612a88565b6000526004602052604060002073ffffffffffffffffffffffffffffffffffffffff6001820154161561066e5761065661066a91612684565b604051918291602083526020830190612140565b0390f35b60046040517f871e01b2000000000000000000000000000000000000000000000000000000008152fd5b346101445760a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101445760443567ffffffffffffffff8111610144576106e79036906004016122c4565b9060643567ffffffffffffffff8111610144576107089036906004016122c4565b9160843567ffffffffffffffff8111610144576107299036906004016122c4565b60ff600a94929454166104ac57610744818688602435612cbc565b61075060043533612dad565b9163ffffffff600184015460a01c169561076a3388612c12565b8354946024358614610e1b576107a56040516107948161078d8160038b016125d0565b0382612283565b61079f368c85612838565b90612e1c565b6107c76040516107bc8161078d8160048c016125d0565b61079f368688612838565b6107e96040516107de8161078d8160058d016125d0565b61079f36898d612838565b918080610e14575b80610e0d575b610de357602435885515610c8e575b15610b3d575b15610890575b926108807f41161473ce2ed633d9f902aab9702d16a5531da27ec84e1939abeffe54ad7353959361044e93610872610864978d604051998a996024358b5260a060208c0152600260a08c0191016125d0565b9189830360408b01526128fb565b9186830360608801526128fb565b90838203608085015233976128fb565b61089d600586015461257d565b610ad6575b67ffffffffffffffff8411610aa7576108cb846108c2600588015461257d565b600588016128b4565b6000601f85116001146109a757928492610872610880938a9b9c61094f876108649b9a61044e9a7f41161473ce2ed633d9f902aab9702d16a5531da27ec84e1939abeffe54ad73539e9f60009261099c575b50507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b60058a01555b8c8780610972575b50509c9b9a9950935050929495509250610812565b61097c9133612a88565b60005260056020526109946004356040600020612fda565b508c8761095d565b013590508f8061091d565b9860058601600052602060002060005b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe087168110610a8f5750926108726108809361044e969388968c7f41161473ce2ed633d9f902aab9702d16a5531da27ec84e1939abeffe54ad73539c9d9e9f897fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe06108649e9d1610610a57575b505050600187811b0160058a0155610955565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88b60031b161c199101351690558e8d81610a44565b898c0135825560209b8c019b600190920191016109b7565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040516020810190610b1c81610af060058a01338661293a565b037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08101835282612283565b5190206000526005602052610b376004356040600020613280565b506108a2565b67ffffffffffffffff8311610aa757610b6683610b5d600489015461257d565b600489016128b4565b600083601f8111600114610bc75780610bb292600091610bbc575b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b600487015561080c565b90508601358d610b81565b506004870160005260206000209060005b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe086168110610c765750847fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0811610610c3e575b5050600183811b01600487015561080c565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88660031b161c19908601351690558a80610c2c565b9091602060018192858a013581550193019101610bd8565b67ffffffffffffffff8b11610aa757610cb78b610cae60038a015461257d565b60038a016128b4565b60008b601f8111600114610d175780610d0292600091610d0c57507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b6003880155610806565b90508501358e610b81565b506003880160005260206000209060005b8d7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081168210610dca578091507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0811610610d91575b905060018092501b016003880155610806565b60f87fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9160031b161c19908501351690558b808c610d7e565b5085820135835560019092019160209182019101610d28565b60046040517f6b4a810d000000000000000000000000000000000000000000000000000000008152fd5b50826107f7565b50816107f1565b60046040517f95406722000000000000000000000000000000000000000000000000000000008152fd5b346101445760607ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101445761066a610e8f610e8261236a565b6044359060243590612ae3565b604051918291826121e4565b34610144576020610eb4610eae3661238d565b91612a88565b604051908152f35b346101445760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261014457602073ffffffffffffffffffffffffffffffffffffffff60015416604051908152f35b346101445760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261014457610f45612bc7565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00600a5416600a55337f11a03e25ee25bf1459f9e1cb293ea03707d84917f54a65e32c9a7be2f2edd68a600080a2005b346101445760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101445760005473ffffffffffffffffffffffffffffffffffffffff808216330361104157600154917fffffffffffffffffffffffff0000000000000000000000000000000000000000903382851617600155166000553391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0600080a3005b60046040517f02b543c6000000000000000000000000000000000000000000000000000000008152fd5b346101445760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261014457600880546110a7816123fc565b6110b46040519182612283565b8181526110c0826123fc565b916020937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe060208401940136853760005b82811061114457505050906040519283926020840190602085525180915260408401929160005b82811061112757505050500390f35b835163ffffffff1685528695509381019392810192600101611118565b6001908260005263ffffffff817ff3f7a9fe364faab93b216da50a3214154f22a0a2b415b23a84c8169e8b636ee301541661117f828761252e565b52016110f1565b3461014457611194366122f2565b9161119d612bc7565b60ff600a54166104ac5760005b828110611229575060405191806040840160408552526060830191906000905b8082106112015785151560208601527fcab63bf31d1e656baa23cebef64e12033ea0ffbd44b1278c3747beec2d2f618c85850386a1005b90919283359063ffffffff8216809203610144576001918152602080910194019201906111ca565b600190841561125b5761125363ffffffff61124d611248848888612a67565b612a77565b16612ecf565b505b016111aa565b61127763ffffffff611271611248848888612a67565b16613066565b50611255565b346101445760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101445760ff600a54166104ac576112c260043533612dad565b600181019081549163ffffffff8360a01c169260ff8160c01c16600281101561047d5715610453577fffffffffffffff00ffffffffffffffffffffffffffffffffffffffffffffffff906113163386612c12565b16905580547f17b2d730bb5e064df3fbc6165c8aceb3b0d62c524c196c0bc1012209280bc9a6604051602081528061044e339560026020840191016125d0565b34610144576020807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610144576004359060ff600a54166104ac5761139e8233612dad565b63ffffffff600182015460a01c16926113c4336000526007602052604060002054151590565b156114c25733600052600283526113df816040600020613280565b5083600052600383526113f6816040600020613280565b5060058201611405815461257d565b61148e575b506000526004825261145160056040600020600081556000600182015561143360028201612a1e565b61143f60038201612a1e565b61144b60048201612a1e565b01612a1e565b7f76ee2dfcae10cb8522e62e713e62660e09ecfaab08db15d9404de1914132257161044e82549260405191829186835260023397840191016125d0565b6040516114a381610af087820194338661293a565b519020600052600583526114bb816040600020613280565b508461140a565b60246040517f85982a00000000000000000000000000000000000000000000000000000000008152336004820152fd5b346101445760e07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101445760043567ffffffffffffffff8111610144576115419036906004016122c4565b6044359163ffffffff8316830361014457600260643510156101445760843567ffffffffffffffff81116101445761157d9036906004016122c4565b91909260a43567ffffffffffffffff8111610144576115a09036906004016122c4565b60c43567ffffffffffffffff8111610144576115c09036906004016122c4565b96909560ff600a54166104ac576115d7338a612c12565b60408511611df4576115ed888483602435612cbc565b6115f8858733612a88565b80600052600460205273ffffffffffffffffffffffffffffffffffffffff60016040600020015416611dca576040519061163182612266565b602435825233602083015263ffffffff8b16604083015261165760643560608401612571565b61166236888a612838565b6080830152611672368486612838565b60a0830152611682368688612838565b60c0830152611692368b8b612838565b60e0830152806000526004602052604060002091805183556001830173ffffffffffffffffffffffffffffffffffffffff60208301511681549077ffffffff0000000000000000000000000000000000000000604085015160a01b16906060850151600281101561047d5778ff0000000000000000000000000000000000000000000000007fffffffffffffff000000000000000000000000000000000000000000000000009160c01b1693161717179055608081015180519067ffffffffffffffff8211610aa7576117758261176c600288015461257d565b600288016128b4565b602090601f8311600114611cfe576117c2929160009183611c275750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b60028401555b60a081015180519067ffffffffffffffff8211610aa7576117f9826117f0600388015461257d565b600388016128b4565b602090601f8311600114611c3257611846929160009183611c275750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b60038401555b60c081015180519067ffffffffffffffff8211610aa75761187d82611874600488015461257d565b600488016128b4565b602090601f8311600114611b5a5791806118ce9260e09594600092611a355750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b60048501555b015180519267ffffffffffffffff8411610aa757838d926119068e966118fd600586015461257d565b600586016128b4565b602090601f8311600114611a40579463ffffffff61087295819a957fc4399022965bad9b2b468bbd8c758a7e80cdde36ff3088ddbb7f93bdfb5623cb9f9e9d99946119928761044e9f9b986005936119f69f9a600092611a355750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b9101555b3360005260026020526119ad836040600020612fda565b501660005260036020526119c5816040600020612fda565b508d82611a0c575b5050506108646040519a8b9a6119e58c6064356120d5565b60a060208d015260a08c01916128fb565b97838903608085015216963396602435966128fb565b611a2c92611a1a9133612a88565b60005260056020526040600020612fda565b508c8f8d6119cd565b01519050388061091d565b906005840160005260206000209160005b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe085168110611b3057506108729563ffffffff9a957fc4399022965bad9b2b468bbd8c758a7e80cdde36ff3088ddbb7f93bdfb5623cb9f9e9d999460018761044e9f9b96928f96936119f69f9a94837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe06005971610611af9575b505050811b01910155611996565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c19169055388080611aeb565b939550918194969750600160209291839285015181550194019201918f9492918f97969492611a51565b906004860160005260206000209160005b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe085168110611c0f5750918391600193837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe060e098971610611bd8575b505050811b0160048501556118d4565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690558f8080611bc8565b91926020600181928685015181550194019201611b6b565b015190508f8061091d565b9190600386016000526020600020906000935b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe084168510611ce35760019450837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0811610611cac575b505050811b01600384015561184c565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690558e8080611c9c565b81810151835560209485019460019093019290910190611c45565b9190600286016000526020600020906000935b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe084168510611daf5760019450837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0811610611d78575b505050811b0160028401556117c8565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690558e8080611d68565b81810151835560209485019460019093019290910190611d11565b60046040517fa0677dd0000000000000000000000000000000000000000000000000000000008152fd5b604485604051907f36a7c503000000000000000000000000000000000000000000000000000000008252600482015260406024820152fd5b346101445760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261014457611e63612bc7565b60017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00600a541617600a55337f2789711f6fd67d131ad68378617b5d1d21a2c92b34d7c3745d70b3957c08096c600080a2005b34610144576020807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101445760043567ffffffffffffffff811161014457611f069036906004016122c4565b60ff600a54166104ac57611f1a9133612a88565b90816000526005602052604060002091825491821561066e5760005b838110611f3f57005b80611f4c60019287612eb7565b90549060031b1c60005260048352604060002063ffffffff8382015460a01c1660005260098452604060002054151580611fd6575b611f8d575b5001611f36565b7f95d94f817db4971aa99ba35d0fe019bd8cc39866fbe02b6d47b5f0f3727fb67360405186815260408682015280611fcd339460026040840191016125d0565b0390a286611f86565b50611fee336000526007602052604060002054151590565b611f81565b346101445760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261014457604051604081019080821067ffffffffffffffff831117610aa75761066a91604052601a81527f576f726b666c6f77526567697374727920312e302e302d64657600000000000060208201526040519182916020835260208301906120e2565b346101445760607ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101445760043563ffffffff8116810361014457610e8f61066a916044359060243590612743565b90600282101561047d5752565b919082519283825260005b84811061212c5750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8460006020809697860101520116010190565b6020818301810151848301820152016120ed565b6121e19160e06121d06121be6121ac6101008651865273ffffffffffffffffffffffffffffffffffffffff602088015116602087015263ffffffff6040880151166040870152612198606088015160608801906120d5565b6080870151908060808801528601906120e2565b60a086015185820360a08701526120e2565b60c085015184820360c08601526120e2565b9201519060e08184039101526120e2565b90565b6020808201906020835283518092526040830192602060408460051b8301019501936000915b84831061221a5750505050505090565b9091929394958480612256837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc086600196030187528a51612140565b980193019301919493929061220a565b610100810190811067ffffffffffffffff821117610aa757604052565b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff821117610aa757604052565b9181601f840112156101445782359167ffffffffffffffff8311610144576020838186019501011161014457565b9060407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc8301126101445760043567ffffffffffffffff9283821161014457806023830112156101445781600401359384116101445760248460051b8301011161014457602401919060243580151581036101445790565b6004359073ffffffffffffffffffffffffffffffffffffffff8216820361014457565b9060407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc8301126101445760043573ffffffffffffffffffffffffffffffffffffffff8116810361014457916024359067ffffffffffffffff8211610144576123f8916004016122c4565b9091565b67ffffffffffffffff8111610aa75760051b60200190565b6040519061242182612266565b606060e0836000815260006020820152600060408201526000838201528260808201528260a08201528260c08201520152565b6040516020810181811067ffffffffffffffff821117610aa7576040526000815290565b90612482826123fc565b61248f6040519182612283565b8281527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe06124bd82946123fc565b019060005b8281106124ce57505050565b6020906124d9612414565b828285010152016124c2565b919082018092116124f257565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b919082039182116124f257565b80518210156125425760209160051b010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600282101561047d5752565b90600182811c921680156125c6575b602083101461259757565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b91607f169161258c565b8054600093926125df8261257d565b918282526020936001916001811690816000146126475750600114612606575b5050505050565b90939495506000929192528360002092846000945b838610612633575050505001019038808080806125ff565b80548587018301529401938590820161261b565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00168685015250505090151560051b0101915038808080806125ff565b90600560e060409361273f85519161269b83612266565b6127388397825485526126e560ff600185015473ffffffffffffffffffffffffffffffffffffffff8116602089015263ffffffff8160a01c168489015260c01c1660608701612571565b80516126f88161078d81600288016125d0565b608086015280516127108161078d81600388016125d0565b60a086015280516127288161078d81600488016125d0565b60c08601525180968193016125d0565b0384612283565b0152565b63ffffffff16916000838152600360209060036020526040936040842054908187101561282857612797918160648993118015612820575b612818575b8161278b82856124e5565b11156128085750612521565b946127a186612478565b96845b8781106127b657505050505050505090565b6001908287528486526127d58888206127cf83876124e5565b90612eb7565b905490861b1c8752600486526127ec888820612684565b6127f6828c61252e565b52612801818b61252e565b50016127a4565b6128139150826124e5565b612521565b506064612780565b50801561277b565b50505050505050506121e1612454565b92919267ffffffffffffffff8211610aa7576040519161288060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8401160184612283565b829481845281830111610144578281602093846000960137010152565b8181106128a8575050565b6000815560010161289d565b9190601f81116128c357505050565b6128ef926000526020600020906020601f840160051c830193106128f1575b601f0160051c019061289d565b565b90915081906128e2565b601f82602094937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0938186528686013760008582860101520116010190565b91907fffffffffffffffffffffffffffffffffffffffff0000000000000000000000009060601b1682526014906000928154926129768461257d565b926001946001811690816000146129dd5750600114612998575b505050505090565b9091929395945060005260209460206000206000905b8582106129ca5750505050601492935001013880808080612990565b80548583018501529087019082016129ae565b92505050601494507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0091935016838301528015150201013880808080612990565b612a28815461257d565b9081612a32575050565b81601f60009311600114612a44575055565b908083918252612a63601f60208420940160051c84016001850161289d565b5555565b91908110156125425760051b0190565b3563ffffffff811681036101445790565b91906034612add91836040519485927fffffffffffffffffffffffffffffffffffffffff000000000000000000000000602085019860601b168852848401378101600083820152036014810184520182612283565b51902090565b73ffffffffffffffffffffffffffffffffffffffff1691600083815260029260209060026020526040936040842054908183101561282857612b3a91816064859311801561282057612818578161278b82856124e5565b94612b4486612478565b96845b878110612b5957505050505050505090565b600190828752838652612b728888206127cf83886124e5565b90549060031b1c875260048652612b8a888820612684565b612b94828c61252e565b52612b9f818b61252e565b5001612b47565b3573ffffffffffffffffffffffffffffffffffffffff811681036101445790565b73ffffffffffffffffffffffffffffffffffffffff600154163303612be857565b60046040517f2b5c74de000000000000000000000000000000000000000000000000000000008152fd5b63ffffffff1680600052600960205260406000205415612c8b575073ffffffffffffffffffffffffffffffffffffffff1680600052600760205260406000205415612c5a5750565b602490604051907f85982a000000000000000000000000000000000000000000000000000000000082526004820152fd5b602490604051907f8fe6d7e10000000000000000000000000000000000000000000000000000000082526004820152fd5b91909115612d835760c891828111612d4d5750818111612d185750808211612ce2575050565b60449250604051917ecd56a800000000000000000000000000000000000000000000000000000000835260048301526024820152fd5b604491604051917ecd56a800000000000000000000000000000000000000000000000000000000835260048301526024820152fd5b60449083604051917ecd56a800000000000000000000000000000000000000000000000000000000835260048301526024820152fd5b60046040517f7dc2f4e1000000000000000000000000000000000000000000000000000000008152fd5b90600052600460205260406000209073ffffffffffffffffffffffffffffffffffffffff8060018401541691821561066e5716809103612deb575090565b602490604051907f31ee6dc70000000000000000000000000000000000000000000000000000000082526004820152fd5b9081518151908181149384612e33575b5050505090565b6020929394508201209201201438808080612e2c565b6008548110156125425760086000527ff3f7a9fe364faab93b216da50a3214154f22a0a2b415b23a84c8169e8b636ee30190600090565b6006548110156125425760066000527ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d3f0190600090565b80548210156125425760005260206000200190600090565b600081815260096020526040812054612f835760085468010000000000000000811015612f56579082612f42612f0d84600160409601600855612e49565b81939154907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9060031b92831b921b19161790565b905560085492815260096020522055600190565b6024827f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b905090565b600081815260076020526040812054612f835760065468010000000000000000811015612f56579082612fc6612f0d84600160409601600655612e80565b905560065492815260076020522055600190565b9190600183016000908282528060205260408220541560001461306057845494680100000000000000008610156130335783613023612f0d886001604098999a01855584612eb7565b9055549382526020522055600190565b6024837f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b50925050565b60008181526009602052604081205490919080156131b4577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90818101818111613187576008549083820191821161315a57818103613126575b50505060085480156130f9578101906130d882612e49565b909182549160031b1b19169055600855815260096020526040812055600190565b6024847f4e487b710000000000000000000000000000000000000000000000000000000081526031600452fd5b613144613135612f0d93612e49565b90549060031b1c928392612e49565b90558452600960205260408420553880806130c0565b6024867f4e487b710000000000000000000000000000000000000000000000000000000081526011600452fd5b6024857f4e487b710000000000000000000000000000000000000000000000000000000081526011600452fd5b505090565b60008181526007602052604081205490919080156131b4577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90818101818111613187576006549083820191821161315a5781810361324c575b50505060065480156130f95781019061322b82612e80565b909182549160031b1b19169055600655815260076020526040812055600190565b61326a61325b612f0d93612e80565b90549060031b1c928392612e80565b9055845260076020526040842055388080613213565b90600182019060009281845282602052604084205490811515600014612e2c577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff918281018181116133a55782549084820191821161337857818103613343575b50505080548015613316578201916132f98383612eb7565b909182549160031b1b191690555582526020526040812055600190565b6024867f4e487b710000000000000000000000000000000000000000000000000000000081526031600452fd5b613363613353612f0d9386612eb7565b90549060031b1c92839286612eb7565b905586528460205260408620553880806132e1565b6024887f4e487b710000000000000000000000000000000000000000000000000000000081526011600452fd5b6024877f4e487b710000000000000000000000000000000000000000000000000000000081526011600452fdfea164736f6c6343000818000a",
}

var WorkflowRegistryABI = WorkflowRegistryMetaData.ABI

var WorkflowRegistryBin = WorkflowRegistryMetaData.Bin

func DeployWorkflowRegistry(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *WorkflowRegistry, error) {
	parsed, err := WorkflowRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(WorkflowRegistryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &WorkflowRegistry{address: address, abi: *parsed, WorkflowRegistryCaller: WorkflowRegistryCaller{contract: contract}, WorkflowRegistryTransactor: WorkflowRegistryTransactor{contract: contract}, WorkflowRegistryFilterer: WorkflowRegistryFilterer{contract: contract}}, nil
}

type WorkflowRegistry struct {
	address common.Address
	abi     abi.ABI
	WorkflowRegistryCaller
	WorkflowRegistryTransactor
	WorkflowRegistryFilterer
}

type WorkflowRegistryCaller struct {
	contract *bind.BoundContract
}

type WorkflowRegistryTransactor struct {
	contract *bind.BoundContract
}

type WorkflowRegistryFilterer struct {
	contract *bind.BoundContract
}

type WorkflowRegistrySession struct {
	Contract     *WorkflowRegistry
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type WorkflowRegistryCallerSession struct {
	Contract *WorkflowRegistryCaller
	CallOpts bind.CallOpts
}

type WorkflowRegistryTransactorSession struct {
	Contract     *WorkflowRegistryTransactor
	TransactOpts bind.TransactOpts
}

type WorkflowRegistryRaw struct {
	Contract *WorkflowRegistry
}

type WorkflowRegistryCallerRaw struct {
	Contract *WorkflowRegistryCaller
}

type WorkflowRegistryTransactorRaw struct {
	Contract *WorkflowRegistryTransactor
}

func NewWorkflowRegistry(address common.Address, backend bind.ContractBackend) (*WorkflowRegistry, error) {
	abi, err := abi.JSON(strings.NewReader(WorkflowRegistryABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindWorkflowRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &WorkflowRegistry{address: address, abi: abi, WorkflowRegistryCaller: WorkflowRegistryCaller{contract: contract}, WorkflowRegistryTransactor: WorkflowRegistryTransactor{contract: contract}, WorkflowRegistryFilterer: WorkflowRegistryFilterer{contract: contract}}, nil
}

func NewWorkflowRegistryCaller(address common.Address, caller bind.ContractCaller) (*WorkflowRegistryCaller, error) {
	contract, err := bindWorkflowRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &WorkflowRegistryCaller{contract: contract}, nil
}

func NewWorkflowRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*WorkflowRegistryTransactor, error) {
	contract, err := bindWorkflowRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &WorkflowRegistryTransactor{contract: contract}, nil
}

func NewWorkflowRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*WorkflowRegistryFilterer, error) {
	contract, err := bindWorkflowRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &WorkflowRegistryFilterer{contract: contract}, nil
}

func bindWorkflowRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := WorkflowRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_WorkflowRegistry *WorkflowRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WorkflowRegistry.Contract.WorkflowRegistryCaller.contract.Call(opts, result, method, params...)
}

func (_WorkflowRegistry *WorkflowRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.WorkflowRegistryTransactor.contract.Transfer(opts)
}

func (_WorkflowRegistry *WorkflowRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.WorkflowRegistryTransactor.contract.Transact(opts, method, params...)
}

func (_WorkflowRegistry *WorkflowRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WorkflowRegistry.Contract.contract.Call(opts, result, method, params...)
}

func (_WorkflowRegistry *WorkflowRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.contract.Transfer(opts)
}

func (_WorkflowRegistry *WorkflowRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.contract.Transact(opts, method, params...)
}

func (_WorkflowRegistry *WorkflowRegistryCaller) ComputeHashKey(opts *bind.CallOpts, owner common.Address, field string) ([32]byte, error) {
	var out []interface{}
	err := _WorkflowRegistry.contract.Call(opts, &out, "computeHashKey", owner, field)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_WorkflowRegistry *WorkflowRegistrySession) ComputeHashKey(owner common.Address, field string) ([32]byte, error) {
	return _WorkflowRegistry.Contract.ComputeHashKey(&_WorkflowRegistry.CallOpts, owner, field)
}

func (_WorkflowRegistry *WorkflowRegistryCallerSession) ComputeHashKey(owner common.Address, field string) ([32]byte, error) {
	return _WorkflowRegistry.Contract.ComputeHashKey(&_WorkflowRegistry.CallOpts, owner, field)
}

func (_WorkflowRegistry *WorkflowRegistryCaller) GetAllAllowedDONs(opts *bind.CallOpts) ([]uint32, error) {
	var out []interface{}
	err := _WorkflowRegistry.contract.Call(opts, &out, "getAllAllowedDONs")

	if err != nil {
		return *new([]uint32), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint32)).(*[]uint32)

	return out0, err

}

func (_WorkflowRegistry *WorkflowRegistrySession) GetAllAllowedDONs() ([]uint32, error) {
	return _WorkflowRegistry.Contract.GetAllAllowedDONs(&_WorkflowRegistry.CallOpts)
}

func (_WorkflowRegistry *WorkflowRegistryCallerSession) GetAllAllowedDONs() ([]uint32, error) {
	return _WorkflowRegistry.Contract.GetAllAllowedDONs(&_WorkflowRegistry.CallOpts)
}

func (_WorkflowRegistry *WorkflowRegistryCaller) GetAllAuthorizedAddresses(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _WorkflowRegistry.contract.Call(opts, &out, "getAllAuthorizedAddresses")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_WorkflowRegistry *WorkflowRegistrySession) GetAllAuthorizedAddresses() ([]common.Address, error) {
	return _WorkflowRegistry.Contract.GetAllAuthorizedAddresses(&_WorkflowRegistry.CallOpts)
}

func (_WorkflowRegistry *WorkflowRegistryCallerSession) GetAllAuthorizedAddresses() ([]common.Address, error) {
	return _WorkflowRegistry.Contract.GetAllAuthorizedAddresses(&_WorkflowRegistry.CallOpts)
}

func (_WorkflowRegistry *WorkflowRegistryCaller) GetWorkflowMetadata(opts *bind.CallOpts, workflowOwner common.Address, workflowName string) (WorkflowRegistryWorkflowMetadata, error) {
	var out []interface{}
	err := _WorkflowRegistry.contract.Call(opts, &out, "getWorkflowMetadata", workflowOwner, workflowName)

	if err != nil {
		return *new(WorkflowRegistryWorkflowMetadata), err
	}

	out0 := *abi.ConvertType(out[0], new(WorkflowRegistryWorkflowMetadata)).(*WorkflowRegistryWorkflowMetadata)

	return out0, err

}

func (_WorkflowRegistry *WorkflowRegistrySession) GetWorkflowMetadata(workflowOwner common.Address, workflowName string) (WorkflowRegistryWorkflowMetadata, error) {
	return _WorkflowRegistry.Contract.GetWorkflowMetadata(&_WorkflowRegistry.CallOpts, workflowOwner, workflowName)
}

func (_WorkflowRegistry *WorkflowRegistryCallerSession) GetWorkflowMetadata(workflowOwner common.Address, workflowName string) (WorkflowRegistryWorkflowMetadata, error) {
	return _WorkflowRegistry.Contract.GetWorkflowMetadata(&_WorkflowRegistry.CallOpts, workflowOwner, workflowName)
}

func (_WorkflowRegistry *WorkflowRegistryCaller) GetWorkflowMetadataListByDON(opts *bind.CallOpts, donID uint32, start *big.Int, limit *big.Int) ([]WorkflowRegistryWorkflowMetadata, error) {
	var out []interface{}
	err := _WorkflowRegistry.contract.Call(opts, &out, "getWorkflowMetadataListByDON", donID, start, limit)

	if err != nil {
		return *new([]WorkflowRegistryWorkflowMetadata), err
	}

	out0 := *abi.ConvertType(out[0], new([]WorkflowRegistryWorkflowMetadata)).(*[]WorkflowRegistryWorkflowMetadata)

	return out0, err

}

func (_WorkflowRegistry *WorkflowRegistrySession) GetWorkflowMetadataListByDON(donID uint32, start *big.Int, limit *big.Int) ([]WorkflowRegistryWorkflowMetadata, error) {
	return _WorkflowRegistry.Contract.GetWorkflowMetadataListByDON(&_WorkflowRegistry.CallOpts, donID, start, limit)
}

func (_WorkflowRegistry *WorkflowRegistryCallerSession) GetWorkflowMetadataListByDON(donID uint32, start *big.Int, limit *big.Int) ([]WorkflowRegistryWorkflowMetadata, error) {
	return _WorkflowRegistry.Contract.GetWorkflowMetadataListByDON(&_WorkflowRegistry.CallOpts, donID, start, limit)
}

func (_WorkflowRegistry *WorkflowRegistryCaller) GetWorkflowMetadataListByOwner(opts *bind.CallOpts, workflowOwner common.Address, start *big.Int, limit *big.Int) ([]WorkflowRegistryWorkflowMetadata, error) {
	var out []interface{}
	err := _WorkflowRegistry.contract.Call(opts, &out, "getWorkflowMetadataListByOwner", workflowOwner, start, limit)

	if err != nil {
		return *new([]WorkflowRegistryWorkflowMetadata), err
	}

	out0 := *abi.ConvertType(out[0], new([]WorkflowRegistryWorkflowMetadata)).(*[]WorkflowRegistryWorkflowMetadata)

	return out0, err

}

func (_WorkflowRegistry *WorkflowRegistrySession) GetWorkflowMetadataListByOwner(workflowOwner common.Address, start *big.Int, limit *big.Int) ([]WorkflowRegistryWorkflowMetadata, error) {
	return _WorkflowRegistry.Contract.GetWorkflowMetadataListByOwner(&_WorkflowRegistry.CallOpts, workflowOwner, start, limit)
}

func (_WorkflowRegistry *WorkflowRegistryCallerSession) GetWorkflowMetadataListByOwner(workflowOwner common.Address, start *big.Int, limit *big.Int) ([]WorkflowRegistryWorkflowMetadata, error) {
	return _WorkflowRegistry.Contract.GetWorkflowMetadataListByOwner(&_WorkflowRegistry.CallOpts, workflowOwner, start, limit)
}

func (_WorkflowRegistry *WorkflowRegistryCaller) IsRegistryLocked(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _WorkflowRegistry.contract.Call(opts, &out, "isRegistryLocked")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_WorkflowRegistry *WorkflowRegistrySession) IsRegistryLocked() (bool, error) {
	return _WorkflowRegistry.Contract.IsRegistryLocked(&_WorkflowRegistry.CallOpts)
}

func (_WorkflowRegistry *WorkflowRegistryCallerSession) IsRegistryLocked() (bool, error) {
	return _WorkflowRegistry.Contract.IsRegistryLocked(&_WorkflowRegistry.CallOpts)
}

func (_WorkflowRegistry *WorkflowRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _WorkflowRegistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_WorkflowRegistry *WorkflowRegistrySession) Owner() (common.Address, error) {
	return _WorkflowRegistry.Contract.Owner(&_WorkflowRegistry.CallOpts)
}

func (_WorkflowRegistry *WorkflowRegistryCallerSession) Owner() (common.Address, error) {
	return _WorkflowRegistry.Contract.Owner(&_WorkflowRegistry.CallOpts)
}

func (_WorkflowRegistry *WorkflowRegistryCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _WorkflowRegistry.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_WorkflowRegistry *WorkflowRegistrySession) TypeAndVersion() (string, error) {
	return _WorkflowRegistry.Contract.TypeAndVersion(&_WorkflowRegistry.CallOpts)
}

func (_WorkflowRegistry *WorkflowRegistryCallerSession) TypeAndVersion() (string, error) {
	return _WorkflowRegistry.Contract.TypeAndVersion(&_WorkflowRegistry.CallOpts)
}

func (_WorkflowRegistry *WorkflowRegistryTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WorkflowRegistry.contract.Transact(opts, "acceptOwnership")
}

func (_WorkflowRegistry *WorkflowRegistrySession) AcceptOwnership() (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.AcceptOwnership(&_WorkflowRegistry.TransactOpts)
}

func (_WorkflowRegistry *WorkflowRegistryTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.AcceptOwnership(&_WorkflowRegistry.TransactOpts)
}

func (_WorkflowRegistry *WorkflowRegistryTransactor) ActivateWorkflow(opts *bind.TransactOpts, workflowKey [32]byte) (*types.Transaction, error) {
	return _WorkflowRegistry.contract.Transact(opts, "activateWorkflow", workflowKey)
}

func (_WorkflowRegistry *WorkflowRegistrySession) ActivateWorkflow(workflowKey [32]byte) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.ActivateWorkflow(&_WorkflowRegistry.TransactOpts, workflowKey)
}

func (_WorkflowRegistry *WorkflowRegistryTransactorSession) ActivateWorkflow(workflowKey [32]byte) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.ActivateWorkflow(&_WorkflowRegistry.TransactOpts, workflowKey)
}

func (_WorkflowRegistry *WorkflowRegistryTransactor) DeleteWorkflow(opts *bind.TransactOpts, workflowKey [32]byte) (*types.Transaction, error) {
	return _WorkflowRegistry.contract.Transact(opts, "deleteWorkflow", workflowKey)
}

func (_WorkflowRegistry *WorkflowRegistrySession) DeleteWorkflow(workflowKey [32]byte) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.DeleteWorkflow(&_WorkflowRegistry.TransactOpts, workflowKey)
}

func (_WorkflowRegistry *WorkflowRegistryTransactorSession) DeleteWorkflow(workflowKey [32]byte) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.DeleteWorkflow(&_WorkflowRegistry.TransactOpts, workflowKey)
}

func (_WorkflowRegistry *WorkflowRegistryTransactor) LockRegistry(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WorkflowRegistry.contract.Transact(opts, "lockRegistry")
}

func (_WorkflowRegistry *WorkflowRegistrySession) LockRegistry() (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.LockRegistry(&_WorkflowRegistry.TransactOpts)
}

func (_WorkflowRegistry *WorkflowRegistryTransactorSession) LockRegistry() (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.LockRegistry(&_WorkflowRegistry.TransactOpts)
}

func (_WorkflowRegistry *WorkflowRegistryTransactor) PauseWorkflow(opts *bind.TransactOpts, workflowKey [32]byte) (*types.Transaction, error) {
	return _WorkflowRegistry.contract.Transact(opts, "pauseWorkflow", workflowKey)
}

func (_WorkflowRegistry *WorkflowRegistrySession) PauseWorkflow(workflowKey [32]byte) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.PauseWorkflow(&_WorkflowRegistry.TransactOpts, workflowKey)
}

func (_WorkflowRegistry *WorkflowRegistryTransactorSession) PauseWorkflow(workflowKey [32]byte) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.PauseWorkflow(&_WorkflowRegistry.TransactOpts, workflowKey)
}

func (_WorkflowRegistry *WorkflowRegistryTransactor) RegisterWorkflow(opts *bind.TransactOpts, workflowName string, workflowID [32]byte, donID uint32, status uint8, binaryURL string, configURL string, secretsURL string) (*types.Transaction, error) {
	return _WorkflowRegistry.contract.Transact(opts, "registerWorkflow", workflowName, workflowID, donID, status, binaryURL, configURL, secretsURL)
}

func (_WorkflowRegistry *WorkflowRegistrySession) RegisterWorkflow(workflowName string, workflowID [32]byte, donID uint32, status uint8, binaryURL string, configURL string, secretsURL string) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.RegisterWorkflow(&_WorkflowRegistry.TransactOpts, workflowName, workflowID, donID, status, binaryURL, configURL, secretsURL)
}

func (_WorkflowRegistry *WorkflowRegistryTransactorSession) RegisterWorkflow(workflowName string, workflowID [32]byte, donID uint32, status uint8, binaryURL string, configURL string, secretsURL string) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.RegisterWorkflow(&_WorkflowRegistry.TransactOpts, workflowName, workflowID, donID, status, binaryURL, configURL, secretsURL)
}

func (_WorkflowRegistry *WorkflowRegistryTransactor) RequestForceUpdateSecrets(opts *bind.TransactOpts, secretsURL string) (*types.Transaction, error) {
	return _WorkflowRegistry.contract.Transact(opts, "requestForceUpdateSecrets", secretsURL)
}

func (_WorkflowRegistry *WorkflowRegistrySession) RequestForceUpdateSecrets(secretsURL string) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.RequestForceUpdateSecrets(&_WorkflowRegistry.TransactOpts, secretsURL)
}

func (_WorkflowRegistry *WorkflowRegistryTransactorSession) RequestForceUpdateSecrets(secretsURL string) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.RequestForceUpdateSecrets(&_WorkflowRegistry.TransactOpts, secretsURL)
}

func (_WorkflowRegistry *WorkflowRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _WorkflowRegistry.contract.Transact(opts, "transferOwnership", to)
}

func (_WorkflowRegistry *WorkflowRegistrySession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.TransferOwnership(&_WorkflowRegistry.TransactOpts, to)
}

func (_WorkflowRegistry *WorkflowRegistryTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.TransferOwnership(&_WorkflowRegistry.TransactOpts, to)
}

func (_WorkflowRegistry *WorkflowRegistryTransactor) UnlockRegistry(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WorkflowRegistry.contract.Transact(opts, "unlockRegistry")
}

func (_WorkflowRegistry *WorkflowRegistrySession) UnlockRegistry() (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.UnlockRegistry(&_WorkflowRegistry.TransactOpts)
}

func (_WorkflowRegistry *WorkflowRegistryTransactorSession) UnlockRegistry() (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.UnlockRegistry(&_WorkflowRegistry.TransactOpts)
}

func (_WorkflowRegistry *WorkflowRegistryTransactor) UpdateAllowedDONs(opts *bind.TransactOpts, donIDs []uint32, allowed bool) (*types.Transaction, error) {
	return _WorkflowRegistry.contract.Transact(opts, "updateAllowedDONs", donIDs, allowed)
}

func (_WorkflowRegistry *WorkflowRegistrySession) UpdateAllowedDONs(donIDs []uint32, allowed bool) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.UpdateAllowedDONs(&_WorkflowRegistry.TransactOpts, donIDs, allowed)
}

func (_WorkflowRegistry *WorkflowRegistryTransactorSession) UpdateAllowedDONs(donIDs []uint32, allowed bool) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.UpdateAllowedDONs(&_WorkflowRegistry.TransactOpts, donIDs, allowed)
}

func (_WorkflowRegistry *WorkflowRegistryTransactor) UpdateAuthorizedAddresses(opts *bind.TransactOpts, addresses []common.Address, allowed bool) (*types.Transaction, error) {
	return _WorkflowRegistry.contract.Transact(opts, "updateAuthorizedAddresses", addresses, allowed)
}

func (_WorkflowRegistry *WorkflowRegistrySession) UpdateAuthorizedAddresses(addresses []common.Address, allowed bool) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.UpdateAuthorizedAddresses(&_WorkflowRegistry.TransactOpts, addresses, allowed)
}

func (_WorkflowRegistry *WorkflowRegistryTransactorSession) UpdateAuthorizedAddresses(addresses []common.Address, allowed bool) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.UpdateAuthorizedAddresses(&_WorkflowRegistry.TransactOpts, addresses, allowed)
}

func (_WorkflowRegistry *WorkflowRegistryTransactor) UpdateWorkflow(opts *bind.TransactOpts, workflowKey [32]byte, newWorkflowID [32]byte, binaryURL string, configURL string, secretsURL string) (*types.Transaction, error) {
	return _WorkflowRegistry.contract.Transact(opts, "updateWorkflow", workflowKey, newWorkflowID, binaryURL, configURL, secretsURL)
}

func (_WorkflowRegistry *WorkflowRegistrySession) UpdateWorkflow(workflowKey [32]byte, newWorkflowID [32]byte, binaryURL string, configURL string, secretsURL string) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.UpdateWorkflow(&_WorkflowRegistry.TransactOpts, workflowKey, newWorkflowID, binaryURL, configURL, secretsURL)
}

func (_WorkflowRegistry *WorkflowRegistryTransactorSession) UpdateWorkflow(workflowKey [32]byte, newWorkflowID [32]byte, binaryURL string, configURL string, secretsURL string) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.UpdateWorkflow(&_WorkflowRegistry.TransactOpts, workflowKey, newWorkflowID, binaryURL, configURL, secretsURL)
}

type WorkflowRegistryAllowedDONsUpdatedV1Iterator struct {
	Event *WorkflowRegistryAllowedDONsUpdatedV1

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WorkflowRegistryAllowedDONsUpdatedV1Iterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WorkflowRegistryAllowedDONsUpdatedV1)
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
		it.Event = new(WorkflowRegistryAllowedDONsUpdatedV1)
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

func (it *WorkflowRegistryAllowedDONsUpdatedV1Iterator) Error() error {
	return it.fail
}

func (it *WorkflowRegistryAllowedDONsUpdatedV1Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WorkflowRegistryAllowedDONsUpdatedV1 struct {
	DonIDs  []uint32
	Allowed bool
	Raw     types.Log
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) FilterAllowedDONsUpdatedV1(opts *bind.FilterOpts) (*WorkflowRegistryAllowedDONsUpdatedV1Iterator, error) {

	logs, sub, err := _WorkflowRegistry.contract.FilterLogs(opts, "AllowedDONsUpdatedV1")
	if err != nil {
		return nil, err
	}
	return &WorkflowRegistryAllowedDONsUpdatedV1Iterator{contract: _WorkflowRegistry.contract, event: "AllowedDONsUpdatedV1", logs: logs, sub: sub}, nil
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) WatchAllowedDONsUpdatedV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryAllowedDONsUpdatedV1) (event.Subscription, error) {

	logs, sub, err := _WorkflowRegistry.contract.WatchLogs(opts, "AllowedDONsUpdatedV1")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WorkflowRegistryAllowedDONsUpdatedV1)
				if err := _WorkflowRegistry.contract.UnpackLog(event, "AllowedDONsUpdatedV1", log); err != nil {
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

func (_WorkflowRegistry *WorkflowRegistryFilterer) ParseAllowedDONsUpdatedV1(log types.Log) (*WorkflowRegistryAllowedDONsUpdatedV1, error) {
	event := new(WorkflowRegistryAllowedDONsUpdatedV1)
	if err := _WorkflowRegistry.contract.UnpackLog(event, "AllowedDONsUpdatedV1", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type WorkflowRegistryAuthorizedAddressesUpdatedV1Iterator struct {
	Event *WorkflowRegistryAuthorizedAddressesUpdatedV1

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WorkflowRegistryAuthorizedAddressesUpdatedV1Iterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WorkflowRegistryAuthorizedAddressesUpdatedV1)
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
		it.Event = new(WorkflowRegistryAuthorizedAddressesUpdatedV1)
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

func (it *WorkflowRegistryAuthorizedAddressesUpdatedV1Iterator) Error() error {
	return it.fail
}

func (it *WorkflowRegistryAuthorizedAddressesUpdatedV1Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WorkflowRegistryAuthorizedAddressesUpdatedV1 struct {
	Addresses []common.Address
	Allowed   bool
	Raw       types.Log
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) FilterAuthorizedAddressesUpdatedV1(opts *bind.FilterOpts) (*WorkflowRegistryAuthorizedAddressesUpdatedV1Iterator, error) {

	logs, sub, err := _WorkflowRegistry.contract.FilterLogs(opts, "AuthorizedAddressesUpdatedV1")
	if err != nil {
		return nil, err
	}
	return &WorkflowRegistryAuthorizedAddressesUpdatedV1Iterator{contract: _WorkflowRegistry.contract, event: "AuthorizedAddressesUpdatedV1", logs: logs, sub: sub}, nil
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) WatchAuthorizedAddressesUpdatedV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryAuthorizedAddressesUpdatedV1) (event.Subscription, error) {

	logs, sub, err := _WorkflowRegistry.contract.WatchLogs(opts, "AuthorizedAddressesUpdatedV1")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WorkflowRegistryAuthorizedAddressesUpdatedV1)
				if err := _WorkflowRegistry.contract.UnpackLog(event, "AuthorizedAddressesUpdatedV1", log); err != nil {
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

func (_WorkflowRegistry *WorkflowRegistryFilterer) ParseAuthorizedAddressesUpdatedV1(log types.Log) (*WorkflowRegistryAuthorizedAddressesUpdatedV1, error) {
	event := new(WorkflowRegistryAuthorizedAddressesUpdatedV1)
	if err := _WorkflowRegistry.contract.UnpackLog(event, "AuthorizedAddressesUpdatedV1", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type WorkflowRegistryOwnershipTransferRequestedIterator struct {
	Event *WorkflowRegistryOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WorkflowRegistryOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WorkflowRegistryOwnershipTransferRequested)
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
		it.Event = new(WorkflowRegistryOwnershipTransferRequested)
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

func (it *WorkflowRegistryOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *WorkflowRegistryOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WorkflowRegistryOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*WorkflowRegistryOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _WorkflowRegistry.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &WorkflowRegistryOwnershipTransferRequestedIterator{contract: _WorkflowRegistry.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _WorkflowRegistry.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WorkflowRegistryOwnershipTransferRequested)
				if err := _WorkflowRegistry.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_WorkflowRegistry *WorkflowRegistryFilterer) ParseOwnershipTransferRequested(log types.Log) (*WorkflowRegistryOwnershipTransferRequested, error) {
	event := new(WorkflowRegistryOwnershipTransferRequested)
	if err := _WorkflowRegistry.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type WorkflowRegistryOwnershipTransferredIterator struct {
	Event *WorkflowRegistryOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WorkflowRegistryOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WorkflowRegistryOwnershipTransferred)
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
		it.Event = new(WorkflowRegistryOwnershipTransferred)
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

func (it *WorkflowRegistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *WorkflowRegistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WorkflowRegistryOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*WorkflowRegistryOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _WorkflowRegistry.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &WorkflowRegistryOwnershipTransferredIterator{contract: _WorkflowRegistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _WorkflowRegistry.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WorkflowRegistryOwnershipTransferred)
				if err := _WorkflowRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_WorkflowRegistry *WorkflowRegistryFilterer) ParseOwnershipTransferred(log types.Log) (*WorkflowRegistryOwnershipTransferred, error) {
	event := new(WorkflowRegistryOwnershipTransferred)
	if err := _WorkflowRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type WorkflowRegistryRegistryLockedV1Iterator struct {
	Event *WorkflowRegistryRegistryLockedV1

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WorkflowRegistryRegistryLockedV1Iterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WorkflowRegistryRegistryLockedV1)
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
		it.Event = new(WorkflowRegistryRegistryLockedV1)
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

func (it *WorkflowRegistryRegistryLockedV1Iterator) Error() error {
	return it.fail
}

func (it *WorkflowRegistryRegistryLockedV1Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WorkflowRegistryRegistryLockedV1 struct {
	LockedBy common.Address
	Raw      types.Log
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) FilterRegistryLockedV1(opts *bind.FilterOpts, lockedBy []common.Address) (*WorkflowRegistryRegistryLockedV1Iterator, error) {

	var lockedByRule []interface{}
	for _, lockedByItem := range lockedBy {
		lockedByRule = append(lockedByRule, lockedByItem)
	}

	logs, sub, err := _WorkflowRegistry.contract.FilterLogs(opts, "RegistryLockedV1", lockedByRule)
	if err != nil {
		return nil, err
	}
	return &WorkflowRegistryRegistryLockedV1Iterator{contract: _WorkflowRegistry.contract, event: "RegistryLockedV1", logs: logs, sub: sub}, nil
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) WatchRegistryLockedV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryRegistryLockedV1, lockedBy []common.Address) (event.Subscription, error) {

	var lockedByRule []interface{}
	for _, lockedByItem := range lockedBy {
		lockedByRule = append(lockedByRule, lockedByItem)
	}

	logs, sub, err := _WorkflowRegistry.contract.WatchLogs(opts, "RegistryLockedV1", lockedByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WorkflowRegistryRegistryLockedV1)
				if err := _WorkflowRegistry.contract.UnpackLog(event, "RegistryLockedV1", log); err != nil {
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

func (_WorkflowRegistry *WorkflowRegistryFilterer) ParseRegistryLockedV1(log types.Log) (*WorkflowRegistryRegistryLockedV1, error) {
	event := new(WorkflowRegistryRegistryLockedV1)
	if err := _WorkflowRegistry.contract.UnpackLog(event, "RegistryLockedV1", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type WorkflowRegistryRegistryUnlockedV1Iterator struct {
	Event *WorkflowRegistryRegistryUnlockedV1

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WorkflowRegistryRegistryUnlockedV1Iterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WorkflowRegistryRegistryUnlockedV1)
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
		it.Event = new(WorkflowRegistryRegistryUnlockedV1)
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

func (it *WorkflowRegistryRegistryUnlockedV1Iterator) Error() error {
	return it.fail
}

func (it *WorkflowRegistryRegistryUnlockedV1Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WorkflowRegistryRegistryUnlockedV1 struct {
	UnlockedBy common.Address
	Raw        types.Log
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) FilterRegistryUnlockedV1(opts *bind.FilterOpts, unlockedBy []common.Address) (*WorkflowRegistryRegistryUnlockedV1Iterator, error) {

	var unlockedByRule []interface{}
	for _, unlockedByItem := range unlockedBy {
		unlockedByRule = append(unlockedByRule, unlockedByItem)
	}

	logs, sub, err := _WorkflowRegistry.contract.FilterLogs(opts, "RegistryUnlockedV1", unlockedByRule)
	if err != nil {
		return nil, err
	}
	return &WorkflowRegistryRegistryUnlockedV1Iterator{contract: _WorkflowRegistry.contract, event: "RegistryUnlockedV1", logs: logs, sub: sub}, nil
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) WatchRegistryUnlockedV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryRegistryUnlockedV1, unlockedBy []common.Address) (event.Subscription, error) {

	var unlockedByRule []interface{}
	for _, unlockedByItem := range unlockedBy {
		unlockedByRule = append(unlockedByRule, unlockedByItem)
	}

	logs, sub, err := _WorkflowRegistry.contract.WatchLogs(opts, "RegistryUnlockedV1", unlockedByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WorkflowRegistryRegistryUnlockedV1)
				if err := _WorkflowRegistry.contract.UnpackLog(event, "RegistryUnlockedV1", log); err != nil {
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

func (_WorkflowRegistry *WorkflowRegistryFilterer) ParseRegistryUnlockedV1(log types.Log) (*WorkflowRegistryRegistryUnlockedV1, error) {
	event := new(WorkflowRegistryRegistryUnlockedV1)
	if err := _WorkflowRegistry.contract.UnpackLog(event, "RegistryUnlockedV1", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type WorkflowRegistryWorkflowActivatedV1Iterator struct {
	Event *WorkflowRegistryWorkflowActivatedV1

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WorkflowRegistryWorkflowActivatedV1Iterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WorkflowRegistryWorkflowActivatedV1)
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
		it.Event = new(WorkflowRegistryWorkflowActivatedV1)
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

func (it *WorkflowRegistryWorkflowActivatedV1Iterator) Error() error {
	return it.fail
}

func (it *WorkflowRegistryWorkflowActivatedV1Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WorkflowRegistryWorkflowActivatedV1 struct {
	WorkflowID    [32]byte
	WorkflowOwner common.Address
	DonID         uint32
	WorkflowName  string
	Raw           types.Log
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) FilterWorkflowActivatedV1(opts *bind.FilterOpts, workflowID [][32]byte, workflowOwner []common.Address, donID []uint32) (*WorkflowRegistryWorkflowActivatedV1Iterator, error) {

	var workflowIDRule []interface{}
	for _, workflowIDItem := range workflowID {
		workflowIDRule = append(workflowIDRule, workflowIDItem)
	}
	var workflowOwnerRule []interface{}
	for _, workflowOwnerItem := range workflowOwner {
		workflowOwnerRule = append(workflowOwnerRule, workflowOwnerItem)
	}
	var donIDRule []interface{}
	for _, donIDItem := range donID {
		donIDRule = append(donIDRule, donIDItem)
	}

	logs, sub, err := _WorkflowRegistry.contract.FilterLogs(opts, "WorkflowActivatedV1", workflowIDRule, workflowOwnerRule, donIDRule)
	if err != nil {
		return nil, err
	}
	return &WorkflowRegistryWorkflowActivatedV1Iterator{contract: _WorkflowRegistry.contract, event: "WorkflowActivatedV1", logs: logs, sub: sub}, nil
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) WatchWorkflowActivatedV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryWorkflowActivatedV1, workflowID [][32]byte, workflowOwner []common.Address, donID []uint32) (event.Subscription, error) {

	var workflowIDRule []interface{}
	for _, workflowIDItem := range workflowID {
		workflowIDRule = append(workflowIDRule, workflowIDItem)
	}
	var workflowOwnerRule []interface{}
	for _, workflowOwnerItem := range workflowOwner {
		workflowOwnerRule = append(workflowOwnerRule, workflowOwnerItem)
	}
	var donIDRule []interface{}
	for _, donIDItem := range donID {
		donIDRule = append(donIDRule, donIDItem)
	}

	logs, sub, err := _WorkflowRegistry.contract.WatchLogs(opts, "WorkflowActivatedV1", workflowIDRule, workflowOwnerRule, donIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WorkflowRegistryWorkflowActivatedV1)
				if err := _WorkflowRegistry.contract.UnpackLog(event, "WorkflowActivatedV1", log); err != nil {
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

func (_WorkflowRegistry *WorkflowRegistryFilterer) ParseWorkflowActivatedV1(log types.Log) (*WorkflowRegistryWorkflowActivatedV1, error) {
	event := new(WorkflowRegistryWorkflowActivatedV1)
	if err := _WorkflowRegistry.contract.UnpackLog(event, "WorkflowActivatedV1", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type WorkflowRegistryWorkflowDeletedV1Iterator struct {
	Event *WorkflowRegistryWorkflowDeletedV1

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WorkflowRegistryWorkflowDeletedV1Iterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WorkflowRegistryWorkflowDeletedV1)
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
		it.Event = new(WorkflowRegistryWorkflowDeletedV1)
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

func (it *WorkflowRegistryWorkflowDeletedV1Iterator) Error() error {
	return it.fail
}

func (it *WorkflowRegistryWorkflowDeletedV1Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WorkflowRegistryWorkflowDeletedV1 struct {
	WorkflowID    [32]byte
	WorkflowOwner common.Address
	DonID         uint32
	WorkflowName  string
	Raw           types.Log
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) FilterWorkflowDeletedV1(opts *bind.FilterOpts, workflowID [][32]byte, workflowOwner []common.Address, donID []uint32) (*WorkflowRegistryWorkflowDeletedV1Iterator, error) {

	var workflowIDRule []interface{}
	for _, workflowIDItem := range workflowID {
		workflowIDRule = append(workflowIDRule, workflowIDItem)
	}
	var workflowOwnerRule []interface{}
	for _, workflowOwnerItem := range workflowOwner {
		workflowOwnerRule = append(workflowOwnerRule, workflowOwnerItem)
	}
	var donIDRule []interface{}
	for _, donIDItem := range donID {
		donIDRule = append(donIDRule, donIDItem)
	}

	logs, sub, err := _WorkflowRegistry.contract.FilterLogs(opts, "WorkflowDeletedV1", workflowIDRule, workflowOwnerRule, donIDRule)
	if err != nil {
		return nil, err
	}
	return &WorkflowRegistryWorkflowDeletedV1Iterator{contract: _WorkflowRegistry.contract, event: "WorkflowDeletedV1", logs: logs, sub: sub}, nil
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) WatchWorkflowDeletedV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryWorkflowDeletedV1, workflowID [][32]byte, workflowOwner []common.Address, donID []uint32) (event.Subscription, error) {

	var workflowIDRule []interface{}
	for _, workflowIDItem := range workflowID {
		workflowIDRule = append(workflowIDRule, workflowIDItem)
	}
	var workflowOwnerRule []interface{}
	for _, workflowOwnerItem := range workflowOwner {
		workflowOwnerRule = append(workflowOwnerRule, workflowOwnerItem)
	}
	var donIDRule []interface{}
	for _, donIDItem := range donID {
		donIDRule = append(donIDRule, donIDItem)
	}

	logs, sub, err := _WorkflowRegistry.contract.WatchLogs(opts, "WorkflowDeletedV1", workflowIDRule, workflowOwnerRule, donIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WorkflowRegistryWorkflowDeletedV1)
				if err := _WorkflowRegistry.contract.UnpackLog(event, "WorkflowDeletedV1", log); err != nil {
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

func (_WorkflowRegistry *WorkflowRegistryFilterer) ParseWorkflowDeletedV1(log types.Log) (*WorkflowRegistryWorkflowDeletedV1, error) {
	event := new(WorkflowRegistryWorkflowDeletedV1)
	if err := _WorkflowRegistry.contract.UnpackLog(event, "WorkflowDeletedV1", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1Iterator struct {
	Event *WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1Iterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1)
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
		it.Event = new(WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1)
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

func (it *WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1Iterator) Error() error {
	return it.fail
}

func (it *WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1 struct {
	Owner          common.Address
	SecretsURLHash [32]byte
	WorkflowName   string
	Raw            types.Log
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) FilterWorkflowForceUpdateSecretsRequestedV1(opts *bind.FilterOpts, owner []common.Address) (*WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1Iterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _WorkflowRegistry.contract.FilterLogs(opts, "WorkflowForceUpdateSecretsRequestedV1", ownerRule)
	if err != nil {
		return nil, err
	}
	return &WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1Iterator{contract: _WorkflowRegistry.contract, event: "WorkflowForceUpdateSecretsRequestedV1", logs: logs, sub: sub}, nil
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) WatchWorkflowForceUpdateSecretsRequestedV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1, owner []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _WorkflowRegistry.contract.WatchLogs(opts, "WorkflowForceUpdateSecretsRequestedV1", ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1)
				if err := _WorkflowRegistry.contract.UnpackLog(event, "WorkflowForceUpdateSecretsRequestedV1", log); err != nil {
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

func (_WorkflowRegistry *WorkflowRegistryFilterer) ParseWorkflowForceUpdateSecretsRequestedV1(log types.Log) (*WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1, error) {
	event := new(WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1)
	if err := _WorkflowRegistry.contract.UnpackLog(event, "WorkflowForceUpdateSecretsRequestedV1", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type WorkflowRegistryWorkflowPausedV1Iterator struct {
	Event *WorkflowRegistryWorkflowPausedV1

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WorkflowRegistryWorkflowPausedV1Iterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WorkflowRegistryWorkflowPausedV1)
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
		it.Event = new(WorkflowRegistryWorkflowPausedV1)
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

func (it *WorkflowRegistryWorkflowPausedV1Iterator) Error() error {
	return it.fail
}

func (it *WorkflowRegistryWorkflowPausedV1Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WorkflowRegistryWorkflowPausedV1 struct {
	WorkflowID    [32]byte
	WorkflowOwner common.Address
	DonID         uint32
	WorkflowName  string
	Raw           types.Log
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) FilterWorkflowPausedV1(opts *bind.FilterOpts, workflowID [][32]byte, workflowOwner []common.Address, donID []uint32) (*WorkflowRegistryWorkflowPausedV1Iterator, error) {

	var workflowIDRule []interface{}
	for _, workflowIDItem := range workflowID {
		workflowIDRule = append(workflowIDRule, workflowIDItem)
	}
	var workflowOwnerRule []interface{}
	for _, workflowOwnerItem := range workflowOwner {
		workflowOwnerRule = append(workflowOwnerRule, workflowOwnerItem)
	}
	var donIDRule []interface{}
	for _, donIDItem := range donID {
		donIDRule = append(donIDRule, donIDItem)
	}

	logs, sub, err := _WorkflowRegistry.contract.FilterLogs(opts, "WorkflowPausedV1", workflowIDRule, workflowOwnerRule, donIDRule)
	if err != nil {
		return nil, err
	}
	return &WorkflowRegistryWorkflowPausedV1Iterator{contract: _WorkflowRegistry.contract, event: "WorkflowPausedV1", logs: logs, sub: sub}, nil
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) WatchWorkflowPausedV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryWorkflowPausedV1, workflowID [][32]byte, workflowOwner []common.Address, donID []uint32) (event.Subscription, error) {

	var workflowIDRule []interface{}
	for _, workflowIDItem := range workflowID {
		workflowIDRule = append(workflowIDRule, workflowIDItem)
	}
	var workflowOwnerRule []interface{}
	for _, workflowOwnerItem := range workflowOwner {
		workflowOwnerRule = append(workflowOwnerRule, workflowOwnerItem)
	}
	var donIDRule []interface{}
	for _, donIDItem := range donID {
		donIDRule = append(donIDRule, donIDItem)
	}

	logs, sub, err := _WorkflowRegistry.contract.WatchLogs(opts, "WorkflowPausedV1", workflowIDRule, workflowOwnerRule, donIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WorkflowRegistryWorkflowPausedV1)
				if err := _WorkflowRegistry.contract.UnpackLog(event, "WorkflowPausedV1", log); err != nil {
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

func (_WorkflowRegistry *WorkflowRegistryFilterer) ParseWorkflowPausedV1(log types.Log) (*WorkflowRegistryWorkflowPausedV1, error) {
	event := new(WorkflowRegistryWorkflowPausedV1)
	if err := _WorkflowRegistry.contract.UnpackLog(event, "WorkflowPausedV1", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type WorkflowRegistryWorkflowRegisteredV1Iterator struct {
	Event *WorkflowRegistryWorkflowRegisteredV1

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WorkflowRegistryWorkflowRegisteredV1Iterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WorkflowRegistryWorkflowRegisteredV1)
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
		it.Event = new(WorkflowRegistryWorkflowRegisteredV1)
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

func (it *WorkflowRegistryWorkflowRegisteredV1Iterator) Error() error {
	return it.fail
}

func (it *WorkflowRegistryWorkflowRegisteredV1Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WorkflowRegistryWorkflowRegisteredV1 struct {
	WorkflowID    [32]byte
	WorkflowOwner common.Address
	DonID         uint32
	Status        uint8
	WorkflowName  string
	BinaryURL     string
	ConfigURL     string
	SecretsURL    string
	Raw           types.Log
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) FilterWorkflowRegisteredV1(opts *bind.FilterOpts, workflowID [][32]byte, workflowOwner []common.Address, donID []uint32) (*WorkflowRegistryWorkflowRegisteredV1Iterator, error) {

	var workflowIDRule []interface{}
	for _, workflowIDItem := range workflowID {
		workflowIDRule = append(workflowIDRule, workflowIDItem)
	}
	var workflowOwnerRule []interface{}
	for _, workflowOwnerItem := range workflowOwner {
		workflowOwnerRule = append(workflowOwnerRule, workflowOwnerItem)
	}
	var donIDRule []interface{}
	for _, donIDItem := range donID {
		donIDRule = append(donIDRule, donIDItem)
	}

	logs, sub, err := _WorkflowRegistry.contract.FilterLogs(opts, "WorkflowRegisteredV1", workflowIDRule, workflowOwnerRule, donIDRule)
	if err != nil {
		return nil, err
	}
	return &WorkflowRegistryWorkflowRegisteredV1Iterator{contract: _WorkflowRegistry.contract, event: "WorkflowRegisteredV1", logs: logs, sub: sub}, nil
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) WatchWorkflowRegisteredV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryWorkflowRegisteredV1, workflowID [][32]byte, workflowOwner []common.Address, donID []uint32) (event.Subscription, error) {

	var workflowIDRule []interface{}
	for _, workflowIDItem := range workflowID {
		workflowIDRule = append(workflowIDRule, workflowIDItem)
	}
	var workflowOwnerRule []interface{}
	for _, workflowOwnerItem := range workflowOwner {
		workflowOwnerRule = append(workflowOwnerRule, workflowOwnerItem)
	}
	var donIDRule []interface{}
	for _, donIDItem := range donID {
		donIDRule = append(donIDRule, donIDItem)
	}

	logs, sub, err := _WorkflowRegistry.contract.WatchLogs(opts, "WorkflowRegisteredV1", workflowIDRule, workflowOwnerRule, donIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WorkflowRegistryWorkflowRegisteredV1)
				if err := _WorkflowRegistry.contract.UnpackLog(event, "WorkflowRegisteredV1", log); err != nil {
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

func (_WorkflowRegistry *WorkflowRegistryFilterer) ParseWorkflowRegisteredV1(log types.Log) (*WorkflowRegistryWorkflowRegisteredV1, error) {
	event := new(WorkflowRegistryWorkflowRegisteredV1)
	if err := _WorkflowRegistry.contract.UnpackLog(event, "WorkflowRegisteredV1", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type WorkflowRegistryWorkflowUpdatedV1Iterator struct {
	Event *WorkflowRegistryWorkflowUpdatedV1

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WorkflowRegistryWorkflowUpdatedV1Iterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WorkflowRegistryWorkflowUpdatedV1)
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
		it.Event = new(WorkflowRegistryWorkflowUpdatedV1)
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

func (it *WorkflowRegistryWorkflowUpdatedV1Iterator) Error() error {
	return it.fail
}

func (it *WorkflowRegistryWorkflowUpdatedV1Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WorkflowRegistryWorkflowUpdatedV1 struct {
	OldWorkflowID [32]byte
	WorkflowOwner common.Address
	DonID         uint32
	NewWorkflowID [32]byte
	WorkflowName  string
	BinaryURL     string
	ConfigURL     string
	SecretsURL    string
	Raw           types.Log
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) FilterWorkflowUpdatedV1(opts *bind.FilterOpts, oldWorkflowID [][32]byte, workflowOwner []common.Address, donID []uint32) (*WorkflowRegistryWorkflowUpdatedV1Iterator, error) {

	var oldWorkflowIDRule []interface{}
	for _, oldWorkflowIDItem := range oldWorkflowID {
		oldWorkflowIDRule = append(oldWorkflowIDRule, oldWorkflowIDItem)
	}
	var workflowOwnerRule []interface{}
	for _, workflowOwnerItem := range workflowOwner {
		workflowOwnerRule = append(workflowOwnerRule, workflowOwnerItem)
	}
	var donIDRule []interface{}
	for _, donIDItem := range donID {
		donIDRule = append(donIDRule, donIDItem)
	}

	logs, sub, err := _WorkflowRegistry.contract.FilterLogs(opts, "WorkflowUpdatedV1", oldWorkflowIDRule, workflowOwnerRule, donIDRule)
	if err != nil {
		return nil, err
	}
	return &WorkflowRegistryWorkflowUpdatedV1Iterator{contract: _WorkflowRegistry.contract, event: "WorkflowUpdatedV1", logs: logs, sub: sub}, nil
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) WatchWorkflowUpdatedV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryWorkflowUpdatedV1, oldWorkflowID [][32]byte, workflowOwner []common.Address, donID []uint32) (event.Subscription, error) {

	var oldWorkflowIDRule []interface{}
	for _, oldWorkflowIDItem := range oldWorkflowID {
		oldWorkflowIDRule = append(oldWorkflowIDRule, oldWorkflowIDItem)
	}
	var workflowOwnerRule []interface{}
	for _, workflowOwnerItem := range workflowOwner {
		workflowOwnerRule = append(workflowOwnerRule, workflowOwnerItem)
	}
	var donIDRule []interface{}
	for _, donIDItem := range donID {
		donIDRule = append(donIDRule, donIDItem)
	}

	logs, sub, err := _WorkflowRegistry.contract.WatchLogs(opts, "WorkflowUpdatedV1", oldWorkflowIDRule, workflowOwnerRule, donIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WorkflowRegistryWorkflowUpdatedV1)
				if err := _WorkflowRegistry.contract.UnpackLog(event, "WorkflowUpdatedV1", log); err != nil {
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

func (_WorkflowRegistry *WorkflowRegistryFilterer) ParseWorkflowUpdatedV1(log types.Log) (*WorkflowRegistryWorkflowUpdatedV1, error) {
	event := new(WorkflowRegistryWorkflowUpdatedV1)
	if err := _WorkflowRegistry.contract.UnpackLog(event, "WorkflowUpdatedV1", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_WorkflowRegistry *WorkflowRegistry) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _WorkflowRegistry.abi.Events["AllowedDONsUpdatedV1"].ID:
		return _WorkflowRegistry.ParseAllowedDONsUpdatedV1(log)
	case _WorkflowRegistry.abi.Events["AuthorizedAddressesUpdatedV1"].ID:
		return _WorkflowRegistry.ParseAuthorizedAddressesUpdatedV1(log)
	case _WorkflowRegistry.abi.Events["OwnershipTransferRequested"].ID:
		return _WorkflowRegistry.ParseOwnershipTransferRequested(log)
	case _WorkflowRegistry.abi.Events["OwnershipTransferred"].ID:
		return _WorkflowRegistry.ParseOwnershipTransferred(log)
	case _WorkflowRegistry.abi.Events["RegistryLockedV1"].ID:
		return _WorkflowRegistry.ParseRegistryLockedV1(log)
	case _WorkflowRegistry.abi.Events["RegistryUnlockedV1"].ID:
		return _WorkflowRegistry.ParseRegistryUnlockedV1(log)
	case _WorkflowRegistry.abi.Events["WorkflowActivatedV1"].ID:
		return _WorkflowRegistry.ParseWorkflowActivatedV1(log)
	case _WorkflowRegistry.abi.Events["WorkflowDeletedV1"].ID:
		return _WorkflowRegistry.ParseWorkflowDeletedV1(log)
	case _WorkflowRegistry.abi.Events["WorkflowForceUpdateSecretsRequestedV1"].ID:
		return _WorkflowRegistry.ParseWorkflowForceUpdateSecretsRequestedV1(log)
	case _WorkflowRegistry.abi.Events["WorkflowPausedV1"].ID:
		return _WorkflowRegistry.ParseWorkflowPausedV1(log)
	case _WorkflowRegistry.abi.Events["WorkflowRegisteredV1"].ID:
		return _WorkflowRegistry.ParseWorkflowRegisteredV1(log)
	case _WorkflowRegistry.abi.Events["WorkflowUpdatedV1"].ID:
		return _WorkflowRegistry.ParseWorkflowUpdatedV1(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (WorkflowRegistryAllowedDONsUpdatedV1) Topic() common.Hash {
	return common.HexToHash("0xcab63bf31d1e656baa23cebef64e12033ea0ffbd44b1278c3747beec2d2f618c")
}

func (WorkflowRegistryAuthorizedAddressesUpdatedV1) Topic() common.Hash {
	return common.HexToHash("0x509460cccbb176edde6cac28895a4415a24961b8f3a0bd2617b9bb7b4e166c9b")
}

func (WorkflowRegistryOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (WorkflowRegistryOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (WorkflowRegistryRegistryLockedV1) Topic() common.Hash {
	return common.HexToHash("0x2789711f6fd67d131ad68378617b5d1d21a2c92b34d7c3745d70b3957c08096c")
}

func (WorkflowRegistryRegistryUnlockedV1) Topic() common.Hash {
	return common.HexToHash("0x11a03e25ee25bf1459f9e1cb293ea03707d84917f54a65e32c9a7be2f2edd68a")
}

func (WorkflowRegistryWorkflowActivatedV1) Topic() common.Hash {
	return common.HexToHash("0x17b2d730bb5e064df3fbc6165c8aceb3b0d62c524c196c0bc1012209280bc9a6")
}

func (WorkflowRegistryWorkflowDeletedV1) Topic() common.Hash {
	return common.HexToHash("0x76ee2dfcae10cb8522e62e713e62660e09ecfaab08db15d9404de19141322571")
}

func (WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1) Topic() common.Hash {
	return common.HexToHash("0x95d94f817db4971aa99ba35d0fe019bd8cc39866fbe02b6d47b5f0f3727fb673")
}

func (WorkflowRegistryWorkflowPausedV1) Topic() common.Hash {
	return common.HexToHash("0x6a0ed88e9cf3cb493ab4028fcb1dc7d18f0130fcdfba096edde0aadbfbf5e99f")
}

func (WorkflowRegistryWorkflowRegisteredV1) Topic() common.Hash {
	return common.HexToHash("0xc4399022965bad9b2b468bbd8c758a7e80cdde36ff3088ddbb7f93bdfb5623cb")
}

func (WorkflowRegistryWorkflowUpdatedV1) Topic() common.Hash {
	return common.HexToHash("0x41161473ce2ed633d9f902aab9702d16a5531da27ec84e1939abeffe54ad7353")
}

func (_WorkflowRegistry *WorkflowRegistry) Address() common.Address {
	return _WorkflowRegistry.address
}

type WorkflowRegistryInterface interface {
	ComputeHashKey(opts *bind.CallOpts, owner common.Address, field string) ([32]byte, error)

	GetAllAllowedDONs(opts *bind.CallOpts) ([]uint32, error)

	GetAllAuthorizedAddresses(opts *bind.CallOpts) ([]common.Address, error)

	GetWorkflowMetadata(opts *bind.CallOpts, workflowOwner common.Address, workflowName string) (WorkflowRegistryWorkflowMetadata, error)

	GetWorkflowMetadataListByDON(opts *bind.CallOpts, donID uint32, start *big.Int, limit *big.Int) ([]WorkflowRegistryWorkflowMetadata, error)

	GetWorkflowMetadataListByOwner(opts *bind.CallOpts, workflowOwner common.Address, start *big.Int, limit *big.Int) ([]WorkflowRegistryWorkflowMetadata, error)

	IsRegistryLocked(opts *bind.CallOpts) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	ActivateWorkflow(opts *bind.TransactOpts, workflowKey [32]byte) (*types.Transaction, error)

	DeleteWorkflow(opts *bind.TransactOpts, workflowKey [32]byte) (*types.Transaction, error)

	LockRegistry(opts *bind.TransactOpts) (*types.Transaction, error)

	PauseWorkflow(opts *bind.TransactOpts, workflowKey [32]byte) (*types.Transaction, error)

	RegisterWorkflow(opts *bind.TransactOpts, workflowName string, workflowID [32]byte, donID uint32, status uint8, binaryURL string, configURL string, secretsURL string) (*types.Transaction, error)

	RequestForceUpdateSecrets(opts *bind.TransactOpts, secretsURL string) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	UnlockRegistry(opts *bind.TransactOpts) (*types.Transaction, error)

	UpdateAllowedDONs(opts *bind.TransactOpts, donIDs []uint32, allowed bool) (*types.Transaction, error)

	UpdateAuthorizedAddresses(opts *bind.TransactOpts, addresses []common.Address, allowed bool) (*types.Transaction, error)

	UpdateWorkflow(opts *bind.TransactOpts, workflowKey [32]byte, newWorkflowID [32]byte, binaryURL string, configURL string, secretsURL string) (*types.Transaction, error)

	FilterAllowedDONsUpdatedV1(opts *bind.FilterOpts) (*WorkflowRegistryAllowedDONsUpdatedV1Iterator, error)

	WatchAllowedDONsUpdatedV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryAllowedDONsUpdatedV1) (event.Subscription, error)

	ParseAllowedDONsUpdatedV1(log types.Log) (*WorkflowRegistryAllowedDONsUpdatedV1, error)

	FilterAuthorizedAddressesUpdatedV1(opts *bind.FilterOpts) (*WorkflowRegistryAuthorizedAddressesUpdatedV1Iterator, error)

	WatchAuthorizedAddressesUpdatedV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryAuthorizedAddressesUpdatedV1) (event.Subscription, error)

	ParseAuthorizedAddressesUpdatedV1(log types.Log) (*WorkflowRegistryAuthorizedAddressesUpdatedV1, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*WorkflowRegistryOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*WorkflowRegistryOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*WorkflowRegistryOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*WorkflowRegistryOwnershipTransferred, error)

	FilterRegistryLockedV1(opts *bind.FilterOpts, lockedBy []common.Address) (*WorkflowRegistryRegistryLockedV1Iterator, error)

	WatchRegistryLockedV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryRegistryLockedV1, lockedBy []common.Address) (event.Subscription, error)

	ParseRegistryLockedV1(log types.Log) (*WorkflowRegistryRegistryLockedV1, error)

	FilterRegistryUnlockedV1(opts *bind.FilterOpts, unlockedBy []common.Address) (*WorkflowRegistryRegistryUnlockedV1Iterator, error)

	WatchRegistryUnlockedV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryRegistryUnlockedV1, unlockedBy []common.Address) (event.Subscription, error)

	ParseRegistryUnlockedV1(log types.Log) (*WorkflowRegistryRegistryUnlockedV1, error)

	FilterWorkflowActivatedV1(opts *bind.FilterOpts, workflowID [][32]byte, workflowOwner []common.Address, donID []uint32) (*WorkflowRegistryWorkflowActivatedV1Iterator, error)

	WatchWorkflowActivatedV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryWorkflowActivatedV1, workflowID [][32]byte, workflowOwner []common.Address, donID []uint32) (event.Subscription, error)

	ParseWorkflowActivatedV1(log types.Log) (*WorkflowRegistryWorkflowActivatedV1, error)

	FilterWorkflowDeletedV1(opts *bind.FilterOpts, workflowID [][32]byte, workflowOwner []common.Address, donID []uint32) (*WorkflowRegistryWorkflowDeletedV1Iterator, error)

	WatchWorkflowDeletedV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryWorkflowDeletedV1, workflowID [][32]byte, workflowOwner []common.Address, donID []uint32) (event.Subscription, error)

	ParseWorkflowDeletedV1(log types.Log) (*WorkflowRegistryWorkflowDeletedV1, error)

	FilterWorkflowForceUpdateSecretsRequestedV1(opts *bind.FilterOpts, owner []common.Address) (*WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1Iterator, error)

	WatchWorkflowForceUpdateSecretsRequestedV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1, owner []common.Address) (event.Subscription, error)

	ParseWorkflowForceUpdateSecretsRequestedV1(log types.Log) (*WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1, error)

	FilterWorkflowPausedV1(opts *bind.FilterOpts, workflowID [][32]byte, workflowOwner []common.Address, donID []uint32) (*WorkflowRegistryWorkflowPausedV1Iterator, error)

	WatchWorkflowPausedV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryWorkflowPausedV1, workflowID [][32]byte, workflowOwner []common.Address, donID []uint32) (event.Subscription, error)

	ParseWorkflowPausedV1(log types.Log) (*WorkflowRegistryWorkflowPausedV1, error)

	FilterWorkflowRegisteredV1(opts *bind.FilterOpts, workflowID [][32]byte, workflowOwner []common.Address, donID []uint32) (*WorkflowRegistryWorkflowRegisteredV1Iterator, error)

	WatchWorkflowRegisteredV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryWorkflowRegisteredV1, workflowID [][32]byte, workflowOwner []common.Address, donID []uint32) (event.Subscription, error)

	ParseWorkflowRegisteredV1(log types.Log) (*WorkflowRegistryWorkflowRegisteredV1, error)

	FilterWorkflowUpdatedV1(opts *bind.FilterOpts, oldWorkflowID [][32]byte, workflowOwner []common.Address, donID []uint32) (*WorkflowRegistryWorkflowUpdatedV1Iterator, error)

	WatchWorkflowUpdatedV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryWorkflowUpdatedV1, oldWorkflowID [][32]byte, workflowOwner []common.Address, donID []uint32) (event.Subscription, error)

	ParseWorkflowUpdatedV1(log types.Log) (*WorkflowRegistryWorkflowUpdatedV1, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
