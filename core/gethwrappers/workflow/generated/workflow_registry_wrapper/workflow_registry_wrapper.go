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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"AddressNotAuthorized\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"CallerIsNotWorkflowOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotTransferToSelf\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donID\",\"type\":\"uint32\"}],\"name\":\"DONNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidWorkflowID\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeProposedOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OwnerCannotBeZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RegistryLocked\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"providedLength\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"maxAllowedLength\",\"type\":\"uint8\"}],\"name\":\"URLTooLong\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WorkflowAlreadyInDesiredStatus\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WorkflowAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WorkflowContentNotUpdated\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WorkflowDoesNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WorkflowIDAlreadyExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WorkflowIDNotUpdated\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"providedLength\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"maxAllowedLength\",\"type\":\"uint8\"}],\"name\":\"WorkflowNameTooLong\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32[]\",\"name\":\"donIDs\",\"type\":\"uint32[]\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"AllowedDONsUpdatedV1\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"addresses\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"AuthorizedAddressesUpdatedV1\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"lockedBy\",\"type\":\"address\"}],\"name\":\"RegistryLockedV1\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"unlockedBy\",\"type\":\"address\"}],\"name\":\"RegistryUnlockedV1\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"workflowID\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"workflowOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"donID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"workflowName\",\"type\":\"string\"}],\"name\":\"WorkflowActivatedV1\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"workflowID\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"workflowOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"donID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"workflowName\",\"type\":\"string\"}],\"name\":\"WorkflowDeletedV1\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"secretsURLHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"workflowName\",\"type\":\"string\"}],\"name\":\"WorkflowForceUpdateSecretsRequestedV1\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"workflowID\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"workflowOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"donID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"workflowName\",\"type\":\"string\"}],\"name\":\"WorkflowPausedV1\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"workflowID\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"workflowOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"donID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"enumWorkflowRegistry.WorkflowStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"workflowName\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"binaryURL\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"configURL\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"secretsURL\",\"type\":\"string\"}],\"name\":\"WorkflowRegisteredV1\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"oldWorkflowID\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"workflowOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"donID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"newWorkflowID\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"workflowName\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"binaryURL\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"configURL\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"secretsURL\",\"type\":\"string\"}],\"name\":\"WorkflowUpdatedV1\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"workflowKey\",\"type\":\"bytes32\"}],\"name\":\"activateWorkflow\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"field\",\"type\":\"string\"}],\"name\":\"computeHashKey\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"workflowKey\",\"type\":\"bytes32\"}],\"name\":\"deleteWorkflow\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllAllowedDONs\",\"outputs\":[{\"internalType\":\"uint32[]\",\"name\":\"allowedDONs\",\"type\":\"uint32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllAuthorizedAddresses\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"authorizedAddresses\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"workflowOwner\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"workflowName\",\"type\":\"string\"}],\"name\":\"getWorkflowMetadata\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"workflowID\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"donID\",\"type\":\"uint32\"},{\"internalType\":\"enumWorkflowRegistry.WorkflowStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"workflowName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"binaryURL\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"configURL\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"secretsURL\",\"type\":\"string\"}],\"internalType\":\"structWorkflowRegistry.WorkflowMetadata\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donID\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"start\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"limit\",\"type\":\"uint256\"}],\"name\":\"getWorkflowMetadataListByDON\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"workflowID\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"donID\",\"type\":\"uint32\"},{\"internalType\":\"enumWorkflowRegistry.WorkflowStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"workflowName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"binaryURL\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"configURL\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"secretsURL\",\"type\":\"string\"}],\"internalType\":\"structWorkflowRegistry.WorkflowMetadata[]\",\"name\":\"workflowMetadataList\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"workflowOwner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"start\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"limit\",\"type\":\"uint256\"}],\"name\":\"getWorkflowMetadataListByOwner\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"workflowID\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"donID\",\"type\":\"uint32\"},{\"internalType\":\"enumWorkflowRegistry.WorkflowStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"workflowName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"binaryURL\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"configURL\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"secretsURL\",\"type\":\"string\"}],\"internalType\":\"structWorkflowRegistry.WorkflowMetadata[]\",\"name\":\"workflowMetadataList\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isRegistryLocked\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lockRegistry\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"workflowKey\",\"type\":\"bytes32\"}],\"name\":\"pauseWorkflow\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"workflowName\",\"type\":\"string\"},{\"internalType\":\"bytes32\",\"name\":\"workflowID\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"donID\",\"type\":\"uint32\"},{\"internalType\":\"enumWorkflowRegistry.WorkflowStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"binaryURL\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"configURL\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"secretsURL\",\"type\":\"string\"}],\"name\":\"registerWorkflow\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"secretsURL\",\"type\":\"string\"}],\"name\":\"requestForceUpdateSecrets\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unlockRegistry\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32[]\",\"name\":\"donIDs\",\"type\":\"uint32[]\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"updateAllowedDONs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"addresses\",\"type\":\"address[]\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"updateAuthorizedAddresses\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"workflowKey\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newWorkflowID\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"binaryURL\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"configURL\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"secretsURL\",\"type\":\"string\"}],\"name\":\"updateWorkflow\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080806040523461004a57331561003b57600180546001600160a01b03191633179055600b805460ff191690556040516134e290816100508239f35b639b15e16f60e01b8152600490fd5b600080fdfe6080604052600436101561001257600080fd5b60003560e01c806308e7f63a1461210e578063181f5a771461207f5780632303348a14611f425780632b596f6d14611eb45780633ccd14ff14611572578063695e1340146113965780636f351771146112bd578063724c13dd146111c65780637497066b146110ab57806379ba509714610fd55780637ec0846d14610f4a5780638da5cb5b14610ef85780639f4cb53414610ed7578063b87a019414610e81578063d4b89c7414610698578063db800092146105fd578063e3dce080146104d6578063e690f33214610362578063f2fde38b14610284578063f794bdeb146101495763f99ecb6b1461010357600080fd5b346101445760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261014457602060ff600b54166040519015158152f35b600080fd5b346101445760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610144576007805461018581612488565b610192604051918261230f565b81815261019e82612488565b916020937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe060208401940136853760005b82811061023257505050906040519283926020840190602085525180915260408401929160005b82811061020557505050500390f35b835173ffffffffffffffffffffffffffffffffffffffff16855286955093810193928101926001016101f6565b6001908260005273ffffffffffffffffffffffffffffffffffffffff817fa66cc928b5edb82af9bd49922954155ab7b0942694bea4ce44661d9a8736c68801541661027d82876125ba565b52016101cf565b346101445760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610144576102bb6123f6565b6102c3612c53565b73ffffffffffffffffffffffffffffffffffffffff8091169033821461033857817fffffffffffffffffffffffff00000000000000000000000000000000000000006000541617600055600154167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278600080a3005b60046040517fdad89dca000000000000000000000000000000000000000000000000000000008152fd5b346101445760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101445760ff600b54166104ac576103a760043533612eb0565b600181019081549160ff8360c01c16600281101561047d576001146104535778010000000000000000000000000000000000000000000000007fffffffffffffff00ffffffffffffffffffffffffffffffffffffffffffffffff841617905580547f6a0ed88e9cf3cb493ab4028fcb1dc7d18f0130fcdfba096edde0aadbfbf5e99f63ffffffff604051946020865260a01c16938061044e3395600260208401910161265c565b0390a4005b60046040517f6f861db1000000000000000000000000000000000000000000000000000000008152fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b60046040517f78a4e7d9000000000000000000000000000000000000000000000000000000008152fd5b34610144576104e43661237e565b916104ed612c53565b60ff600b54166104ac5760005b828110610589575060405191806040840160408552526060830191906000905b8082106105515785151560208601527f509460cccbb176edde6cac28895a4415a24961b8f3a0bd2617b9bb7b4e166c9b85850386a1005b90919283359073ffffffffffffffffffffffffffffffffffffffff82168092036101445760019181526020809101940192019061051a565b60019084156105cb576105c373ffffffffffffffffffffffffffffffffffffffff6105bd6105b8848888612af3565b612c32565b1661308b565b505b016104fa565b6105f773ffffffffffffffffffffffffffffffffffffffff6105f16105b8848888612af3565b166132bc565b506105c5565b346101445761061d61060e36612419565b916106176124a0565b50612b14565b6000526004602052604060002073ffffffffffffffffffffffffffffffffffffffff6001820154161561066e5761065661066a91612710565b6040519182916020835260208301906121cc565b0390f35b60046040517f871e01b2000000000000000000000000000000000000000000000000000000008152fd5b346101445760a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101445760443567ffffffffffffffff8111610144576106e7903690600401612350565b9060643567ffffffffffffffff811161014457610708903690600401612350565b9160843567ffffffffffffffff811161014457610729903690600401612350565b60ff600b94929454166104ac57610741818688612d48565b61074d60043533612eb0565b9163ffffffff600184015460a01c16956107673388612c9e565b8354946024358614610e57576107a26040516107918161078a8160038b0161265c565b038261230f565b61079c368c856128c4565b90612f1f565b6107c46040516107b98161078a8160048c0161265c565b61079c3686886128c4565b6107e66040516107db8161078a8160058d0161265c565b61079c36898d6128c4565b918080610e50575b80610e49575b610e1f57610803602435612e08565b88600052600660205260406000207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008154169055602435885515610cca575b15610b79575b156108cc575b926108bc7f41161473ce2ed633d9f902aab9702d16a5531da27ec84e1939abeffe54ad7353959361044e936108ae6108a0978d604051998a996024358b5260a060208c0152600260a08c01910161265c565b9189830360408b0152612987565b918683036060880152612987565b9083820360808501523397612987565b6108d96005860154612609565b610b12575b67ffffffffffffffff8411610ae357610907846108fe6005880154612609565b60058801612940565b6000601f85116001146109e3579284926108ae6108bc938a9b9c61098b876108a09b9a61044e9a7f41161473ce2ed633d9f902aab9702d16a5531da27ec84e1939abeffe54ad73539e9f6000926109d8575b50507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b60058a01555b8c87806109ae575b50509c9b9a995093505092949550925061084e565b6109b89133612b14565b60005260056020526109d060043560406000206130dd565b508c87610999565b013590508f80610959565b9860058601600052602060002060005b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe087168110610acb5750926108ae6108bc9361044e969388968c7f41161473ce2ed633d9f902aab9702d16a5531da27ec84e1939abeffe54ad73539c9d9e9f897fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe06108a09e9d1610610a93575b505050600187811b0160058a0155610991565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88b60031b161c199101351690558e8d81610a80565b898c0135825560209b8c019b600190920191016109f3565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040516020810190610b5881610b2c60058a0133866129c6565b037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0810183528261230f565b5190206000526005602052610b736004356040600020613383565b506108de565b67ffffffffffffffff8311610ae357610ba283610b996004890154612609565b60048901612940565b600083601f8111600114610c035780610bee92600091610bf8575b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b6004870155610848565b90508601358d610bbd565b506004870160005260206000209060005b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe086168110610cb25750847fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0811610610c7a575b5050600183811b016004870155610848565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88660031b161c19908601351690558a80610c68565b9091602060018192858a013581550193019101610c14565b67ffffffffffffffff8b11610ae357610cf38b610cea60038a0154612609565b60038a01612940565b60008b601f8111600114610d535780610d3e92600091610d4857507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b6003880155610842565b90508501358e610bbd565b506003880160005260206000209060005b8d7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081168210610e06578091507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0811610610dcd575b905060018092501b016003880155610842565b60f87fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9160031b161c19908501351690558b808c610dba565b5085820135835560019092019160209182019101610d64565b60046040517f6b4a810d000000000000000000000000000000000000000000000000000000008152fd5b50826107f4565b50816107ee565b60046040517f95406722000000000000000000000000000000000000000000000000000000008152fd5b346101445760607ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101445761066a610ecb610ebe6123f6565b6044359060243590612b6f565b60405191829182612270565b34610144576020610ef0610eea36612419565b91612b14565b604051908152f35b346101445760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261014457602073ffffffffffffffffffffffffffffffffffffffff60015416604051908152f35b346101445760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261014457610f81612c53565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00600b5416600b557f11a03e25ee25bf1459f9e1cb293ea03707d84917f54a65e32c9a7be2f2edd68a6020604051338152a1005b346101445760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101445760005473ffffffffffffffffffffffffffffffffffffffff808216330361108157600154917fffffffffffffffffffffffff0000000000000000000000000000000000000000903382851617600155166000553391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0600080a3005b60046040517f02b543c6000000000000000000000000000000000000000000000000000000008152fd5b346101445760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261014457600980546110e781612488565b6110f4604051918261230f565b81815261110082612488565b916020937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe060208401940136853760005b82811061118457505050906040519283926020840190602085525180915260408401929160005b82811061116757505050500390f35b835163ffffffff1685528695509381019392810192600101611158565b6001908260005263ffffffff817f6e1540171b6c0c960b71a7020d9f60077f6af931a8bbf590da0223dacf75c7af0154166111bf82876125ba565b5201611131565b34610144576111d43661237e565b916111dd612c53565b60ff600b54166104ac5760005b828110611269575060405191806040840160408552526060830191906000905b8082106112415785151560208601527fcab63bf31d1e656baa23cebef64e12033ea0ffbd44b1278c3747beec2d2f618c85850386a1005b90919283359063ffffffff82168092036101445760019181526020809101940192019061120a565b600190841561129b5761129363ffffffff61128d611288848888612af3565b612b03565b16612fd2565b505b016111ea565b6112b763ffffffff6112b1611288848888612af3565b16613169565b50611295565b346101445760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101445760ff600b54166104ac5761130260043533612eb0565b600181019081549163ffffffff8360a01c169260ff8160c01c16600281101561047d5715610453577fffffffffffffff00ffffffffffffffffffffffffffffffffffffffffffffffff906113563386612c9e565b16905580547f17b2d730bb5e064df3fbc6165c8aceb3b0d62c524c196c0bc1012209280bc9a6604051602081528061044e3395600260208401910161265c565b34610144576020807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610144576004359060ff600b54166104ac576113de8233612eb0565b916113f6336000526008602052604060002054151590565b156115425782600493546000526006835260406000207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0081541690553360005260028352611448826040600020613383565b50600181019063ffffffff80835460a01c166000526003855261146f846040600020613383565b506005820161147e8154612609565b61150e575b508154925460a01c16917f76ee2dfcae10cb8522e62e713e62660e09ecfaab08db15d9404de19141322571604051868152806114c6339560028a8401910161265c565b0390a46000525261150c6005604060002060008155600060018201556114ee60028201612aaa565b6114fa60038201612aaa565b61150660048201612aaa565b01612aaa565b005b60405161152381610b2c8982019433866129c6565b5190206000526005855261153b846040600020613383565b5086611483565b60246040517f85982a00000000000000000000000000000000000000000000000000000000008152336004820152fd5b346101445760e07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101445760043567ffffffffffffffff8111610144576115c1903690600401612350565b6044359163ffffffff8316830361014457600260643510156101445760843567ffffffffffffffff8111610144576115fd903690600401612350565b91909260a43567ffffffffffffffff811161014457611620903690600401612350565b60c43567ffffffffffffffff811161014457611640903690600401612350565b96909560ff600b54166104ac57611657338a612c9e565b60408511611e7c5761166a888483612d48565b611675858733612b14565b80600052600460205273ffffffffffffffffffffffffffffffffffffffff60016040600020015416611e52576116ac602435612e08565b604051906116b9826122f2565b602435825233602083015263ffffffff8b1660408301526116df606435606084016125fd565b6116ea36888a6128c4565b60808301526116fa3684866128c4565b60a083015261170a3686886128c4565b60c083015261171a368b8b6128c4565b60e0830152806000526004602052604060002091805183556001830173ffffffffffffffffffffffffffffffffffffffff60208301511681549077ffffffff0000000000000000000000000000000000000000604085015160a01b16906060850151600281101561047d5778ff0000000000000000000000000000000000000000000000007fffffffffffffff000000000000000000000000000000000000000000000000009160c01b1693161717179055608081015180519067ffffffffffffffff8211610ae3576117fd826117f46002880154612609565b60028801612940565b602090601f8311600114611d865761184a929160009183611caf5750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b60028401555b60a081015180519067ffffffffffffffff8211610ae357611881826118786003880154612609565b60038801612940565b602090601f8311600114611cba576118ce929160009183611caf5750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b60038401555b60c081015180519067ffffffffffffffff8211610ae357611905826118fc6004880154612609565b60048801612940565b602090601f8311600114611be25791806119569260e09594600092611abd5750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b60048501555b015180519267ffffffffffffffff8411610ae357838d9261198e8e966119856005860154612609565b60058601612940565b602090601f8311600114611ac8579463ffffffff6108ae95819a957fc4399022965bad9b2b468bbd8c758a7e80cdde36ff3088ddbb7f93bdfb5623cb9f9e9d9994611a1a8761044e9f9b98600593611a7e9f9a600092611abd5750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b9101555b336000526002602052611a358360406000206130dd565b50166000526003602052611a4d8160406000206130dd565b508d82611a94575b5050506108a06040519a8b9a611a6d8c606435612161565b60a060208d015260a08c0191612987565b9783890360808501521696339660243596612987565b611ab492611aa29133612b14565b600052600560205260406000206130dd565b508c8f8d611a55565b015190503880610959565b906005840160005260206000209160005b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe085168110611bb857506108ae9563ffffffff9a957fc4399022965bad9b2b468bbd8c758a7e80cdde36ff3088ddbb7f93bdfb5623cb9f9e9d999460018761044e9f9b96928f9693611a7e9f9a94837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe06005971610611b81575b505050811b01910155611a1e565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c19169055388080611b73565b939550918194969750600160209291839285015181550194019201918f9492918f97969492611ad9565b906004860160005260206000209160005b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe085168110611c975750918391600193837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe060e098971610611c60575b505050811b01600485015561195c565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690558f8080611c50565b91926020600181928685015181550194019201611bf3565b015190508f80610959565b9190600386016000526020600020906000935b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe084168510611d6b5760019450837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0811610611d34575b505050811b0160038401556118d4565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690558e8080611d24565b81810151835560209485019460019093019290910190611ccd565b9190600286016000526020600020906000935b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe084168510611e375760019450837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0811610611e00575b505050811b016002840155611850565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690558e8080611df0565b81810151835560209485019460019093019290910190611d99565b60046040517fa0677dd0000000000000000000000000000000000000000000000000000000008152fd5b604485604051907f36a7c503000000000000000000000000000000000000000000000000000000008252600482015260406024820152fd5b346101445760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261014457611eeb612c53565b60017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00600b541617600b557f2789711f6fd67d131ad68378617b5d1d21a2c92b34d7c3745d70b3957c08096c6020604051338152a1005b34610144576020807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101445760043567ffffffffffffffff811161014457611f92903690600401612350565b60ff600b54166104ac57611fa69133612b14565b90816000526005602052604060002091825491821561066e5760005b838110611fcb57005b80611fd860019287612fba565b90549060031b1c60005260048352604060002063ffffffff8382015460a01c16600052600a8452604060002054151580612062575b612019575b5001611fc2565b7f95d94f817db4971aa99ba35d0fe019bd8cc39866fbe02b6d47b5f0f3727fb673604051868152604086820152806120593394600260408401910161265c565b0390a286612012565b5061207a336000526008602052604060002054151590565b61200d565b346101445760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261014457604051604081019080821067ffffffffffffffff831117610ae35761066a91604052601a81527f576f726b666c6f77526567697374727920312e302e302d646576000000000000602082015260405191829160208352602083019061216e565b346101445760607ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101445760043563ffffffff8116810361014457610ecb61066a9160443590602435906127cf565b90600282101561047d5752565b919082519283825260005b8481106121b85750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8460006020809697860101520116010190565b602081830181015184830182015201612179565b61226d9160e061225c61224a6122386101008651865273ffffffffffffffffffffffffffffffffffffffff602088015116602087015263ffffffff604088015116604087015261222460608801516060880190612161565b60808701519080608088015286019061216e565b60a086015185820360a087015261216e565b60c085015184820360c086015261216e565b9201519060e081840391015261216e565b90565b6020808201906020835283518092526040830192602060408460051b8301019501936000915b8483106122a65750505050505090565b90919293949584806122e2837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc086600196030187528a516121cc565b9801930193019194939290612296565b610100810190811067ffffffffffffffff821117610ae357604052565b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff821117610ae357604052565b9181601f840112156101445782359167ffffffffffffffff8311610144576020838186019501011161014457565b9060407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc8301126101445760043567ffffffffffffffff9283821161014457806023830112156101445781600401359384116101445760248460051b8301011161014457602401919060243580151581036101445790565b6004359073ffffffffffffffffffffffffffffffffffffffff8216820361014457565b9060407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc8301126101445760043573ffffffffffffffffffffffffffffffffffffffff8116810361014457916024359067ffffffffffffffff82116101445761248491600401612350565b9091565b67ffffffffffffffff8111610ae35760051b60200190565b604051906124ad826122f2565b606060e0836000815260006020820152600060408201526000838201528260808201528260a08201528260c08201520152565b6040516020810181811067ffffffffffffffff821117610ae3576040526000815290565b9061250e82612488565b61251b604051918261230f565b8281527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe06125498294612488565b019060005b82811061255a57505050565b6020906125656124a0565b8282850101520161254e565b9190820180921161257e57565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b9190820391821161257e57565b80518210156125ce5760209160051b010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600282101561047d5752565b90600182811c92168015612652575b602083101461262357565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b91607f1691612618565b80546000939261266b82612609565b918282526020936001916001811690816000146126d35750600114612692575b5050505050565b90939495506000929192528360002092846000945b8386106126bf5750505050010190388080808061268b565b8054858701830152940193859082016126a7565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00168685015250505090151560051b01019150388080808061268b565b90600560e06040936127cb855191612727836122f2565b6127c483978254855261277160ff600185015473ffffffffffffffffffffffffffffffffffffffff8116602089015263ffffffff8160a01c168489015260c01c16606087016125fd565b80516127848161078a816002880161265c565b6080860152805161279c8161078a816003880161265c565b60a086015280516127b48161078a816004880161265c565b60c086015251809681930161265c565b038461230f565b0152565b63ffffffff1691600083815260036020906003602052604093604084205490818710156128b4576128239181606489931180156128ac575b6128a4575b816128178285612571565b111561289457506125ad565b9461282d86612504565b96845b87811061284257505050505050505090565b60019082875284865261286188882061285b8387612571565b90612fba565b905490861b1c875260048652612878888820612710565b612882828c6125ba565b5261288d818b6125ba565b5001612830565b61289f915082612571565b6125ad565b50606461280c565b508015612807565b505050505050505061226d6124e0565b92919267ffffffffffffffff8211610ae3576040519161290c60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116018461230f565b829481845281830111610144578281602093846000960137010152565b818110612934575050565b60008155600101612929565b9190601f811161294f57505050565b61297b926000526020600020906020601f840160051c8301931061297d575b601f0160051c0190612929565b565b909150819061296e565b601f82602094937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0938186528686013760008582860101520116010190565b91907fffffffffffffffffffffffffffffffffffffffff0000000000000000000000009060601b168252601490600092815492612a0284612609565b92600194600181169081600014612a695750600114612a24575b505050505090565b9091929395945060005260209460206000206000905b858210612a565750505050601492935001013880808080612a1c565b8054858301850152908701908201612a3a565b92505050601494507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0091935016838301528015150201013880808080612a1c565b612ab48154612609565b9081612abe575050565b81601f60009311600114612ad0575055565b908083918252612aef601f60208420940160051c840160018501612929565b5555565b91908110156125ce5760051b0190565b3563ffffffff811681036101445790565b91906034612b6991836040519485927fffffffffffffffffffffffffffffffffffffffff000000000000000000000000602085019860601b16885284840137810160008382015203601481018452018261230f565b51902090565b73ffffffffffffffffffffffffffffffffffffffff169160008381526002926020906002602052604093604084205490818310156128b457612bc69181606485931180156128ac576128a457816128178285612571565b94612bd086612504565b96845b878110612be557505050505050505090565b600190828752838652612bfe88882061285b8388612571565b90549060031b1c875260048652612c16888820612710565b612c20828c6125ba565b52612c2b818b6125ba565b5001612bd3565b3573ffffffffffffffffffffffffffffffffffffffff811681036101445790565b73ffffffffffffffffffffffffffffffffffffffff600154163303612c7457565b60046040517f2b5c74de000000000000000000000000000000000000000000000000000000008152fd5b63ffffffff1680600052600a60205260406000205415612d17575073ffffffffffffffffffffffffffffffffffffffff1680600052600860205260406000205415612ce65750565b602490604051907f85982a000000000000000000000000000000000000000000000000000000000082526004820152fd5b602490604051907f8fe6d7e10000000000000000000000000000000000000000000000000000000082526004820152fd5b9060c891828111612dd25750818111612d9d5750808211612d67575050565b60449250604051917ecd56a800000000000000000000000000000000000000000000000000000000835260048301526024820152fd5b604491604051917ecd56a800000000000000000000000000000000000000000000000000000000835260048301526024820152fd5b60449083604051917ecd56a800000000000000000000000000000000000000000000000000000000835260048301526024820152fd5b8015612e865780600052600660205260ff60406000205416612e5c576000526006602052604060002060017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00825416179055565b60046040517f4cb050e4000000000000000000000000000000000000000000000000000000008152fd5b60046040517f7dc2f4e1000000000000000000000000000000000000000000000000000000008152fd5b90600052600460205260406000209073ffffffffffffffffffffffffffffffffffffffff8060018401541691821561066e5716809103612eee575090565b602490604051907f31ee6dc70000000000000000000000000000000000000000000000000000000082526004820152fd5b9081518151908181149384612f36575b5050505090565b6020929394508201209201201438808080612f2f565b6009548110156125ce5760096000527f6e1540171b6c0c960b71a7020d9f60077f6af931a8bbf590da0223dacf75c7af0190600090565b6007548110156125ce5760076000527fa66cc928b5edb82af9bd49922954155ab7b0942694bea4ce44661d9a8736c6880190600090565b80548210156125ce5760005260206000200190600090565b6000818152600a6020526040812054613086576009546801000000000000000081101561305957908261304561301084600160409601600955612f4c565b81939154907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9060031b92831b921b19161790565b9055600954928152600a6020522055600190565b6024827f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b905090565b60008181526008602052604081205461308657600754680100000000000000008110156130595790826130c961301084600160409601600755612f83565b905560075492815260086020522055600190565b9190600183016000908282528060205260408220541560001461316357845494680100000000000000008610156131365783613126613010886001604098999a01855584612fba565b9055549382526020522055600190565b6024837f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b50925050565b6000818152600a602052604081205490919080156132b7577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9081810181811161328a576009549083820191821161325d57818103613229575b50505060095480156131fc578101906131db82612f4c565b909182549160031b1b191690556009558152600a6020526040812055600190565b6024847f4e487b710000000000000000000000000000000000000000000000000000000081526031600452fd5b61324761323861301093612f4c565b90549060031b1c928392612f4c565b90558452600a60205260408420553880806131c3565b6024867f4e487b710000000000000000000000000000000000000000000000000000000081526011600452fd5b6024857f4e487b710000000000000000000000000000000000000000000000000000000081526011600452fd5b505090565b60008181526008602052604081205490919080156132b7577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9081810181811161328a576007549083820191821161325d5781810361334f575b50505060075480156131fc5781019061332e82612f83565b909182549160031b1b19169055600755815260086020526040812055600190565b61336d61335e61301093612f83565b90549060031b1c928392612f83565b9055845260086020526040842055388080613316565b90600182019060009281845282602052604084205490811515600014612f2f577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff918281018181116134a85782549084820191821161347b57818103613446575b50505080548015613419578201916133fc8383612fba565b909182549160031b1b191690555582526020526040812055600190565b6024867f4e487b710000000000000000000000000000000000000000000000000000000081526031600452fd5b6134666134566130109386612fba565b90549060031b1c92839286612fba565b905586528460205260408620553880806133e4565b6024887f4e487b710000000000000000000000000000000000000000000000000000000081526011600452fd5b6024877f4e487b710000000000000000000000000000000000000000000000000000000081526011600452fdfea164736f6c6343000818000a",
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

func (_WorkflowRegistry *WorkflowRegistryFilterer) FilterRegistryLockedV1(opts *bind.FilterOpts) (*WorkflowRegistryRegistryLockedV1Iterator, error) {

	logs, sub, err := _WorkflowRegistry.contract.FilterLogs(opts, "RegistryLockedV1")
	if err != nil {
		return nil, err
	}
	return &WorkflowRegistryRegistryLockedV1Iterator{contract: _WorkflowRegistry.contract, event: "RegistryLockedV1", logs: logs, sub: sub}, nil
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) WatchRegistryLockedV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryRegistryLockedV1) (event.Subscription, error) {

	logs, sub, err := _WorkflowRegistry.contract.WatchLogs(opts, "RegistryLockedV1")
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

func (_WorkflowRegistry *WorkflowRegistryFilterer) FilterRegistryUnlockedV1(opts *bind.FilterOpts) (*WorkflowRegistryRegistryUnlockedV1Iterator, error) {

	logs, sub, err := _WorkflowRegistry.contract.FilterLogs(opts, "RegistryUnlockedV1")
	if err != nil {
		return nil, err
	}
	return &WorkflowRegistryRegistryUnlockedV1Iterator{contract: _WorkflowRegistry.contract, event: "RegistryUnlockedV1", logs: logs, sub: sub}, nil
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) WatchRegistryUnlockedV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryRegistryUnlockedV1) (event.Subscription, error) {

	logs, sub, err := _WorkflowRegistry.contract.WatchLogs(opts, "RegistryUnlockedV1")
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

	FilterRegistryLockedV1(opts *bind.FilterOpts) (*WorkflowRegistryRegistryLockedV1Iterator, error)

	WatchRegistryLockedV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryRegistryLockedV1) (event.Subscription, error)

	ParseRegistryLockedV1(log types.Log) (*WorkflowRegistryRegistryLockedV1, error)

	FilterRegistryUnlockedV1(opts *bind.FilterOpts) (*WorkflowRegistryRegistryUnlockedV1Iterator, error)

	WatchRegistryUnlockedV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryRegistryUnlockedV1) (event.Subscription, error)

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
