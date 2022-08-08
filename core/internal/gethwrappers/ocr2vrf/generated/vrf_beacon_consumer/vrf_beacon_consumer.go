// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_beacon_consumer

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

type ECCArithmeticG1Point struct {
	P [2]*big.Int
}

var BeaconVRFConsumerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"shouldFail\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"beaconPeriodBlocks\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelayArg\",\"type\":\"uint24\"}],\"name\":\"computeAndStoreExpectedOutput\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"fail\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_beaconPeriodBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_ExpectedRandomnessByBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"\",\"type\":\"uint24\"}],\"name\":\"s_ExpectedSeeds\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_ReceivedRandomnessByRequestID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_gasAvailable\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"}],\"name\":\"s_myBeaconRequests\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.SlotNumber\",\"name\":\"slotNumber\",\"type\":\"uint32\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelay\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_randomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"\",\"type\":\"uint24\"}],\"name\":\"s_requestsIDs\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_subId\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelayArg\",\"type\":\"uint24\"},{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"p\",\"type\":\"uint256[2]\"}],\"internalType\":\"structECCArithmetic.G1Point\",\"name\":\"vrfOutput\",\"type\":\"tuple\"}],\"name\":\"setExpectedSeed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"shouldFail\",\"type\":\"bool\"}],\"name\":\"setFail\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"}],\"name\":\"testRedeemRandomness\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"uint64\",\"name\":\"subID\",\"type\":\"uint64\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelayArg\",\"type\":\"uint24\"}],\"name\":\"testRequestRandomness\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subID\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelayArg\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"}],\"name\":\"testRequestRandomnessFulfillment\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561001057600080fd5b5060405161139838038061139883398101604081905261002f9161006c565b6001600160a01b03929092166080819052600680546001600160a01b03191690911790556009805460ff1916911515919091179055600a556100c0565b60008060006060848603121561008157600080fd5b83516001600160a01b038116811461009857600080fd5b602085015190935080151581146100ae57600080fd5b80925050604084015190509250925092565b6080516112bd6100db600039600061064d01526112bd6000f3fe608060405234801561001057600080fd5b506004361061010b5760003560e01c8063706da1ca116100a2578063cd0593df11610071578063cd0593df1461030b578063e4d47bfa14610314578063f08c5daa14610327578063f339c79414610330578063f6eaffc81461034357600080fd5b8063706da1ca146102885780639c9cd015146102ba5780639d769402146102cd578063a9cc4718146102ee57600080fd5b8063563db24c116100de578063563db24c146101f65780635a47dd71146102175780635f15cccc1461022a578063678d38f71461025d57600080fd5b806319a5fa22146101105780633d8b70aa146101a457806340e836ab146101b957806345907626146101cc575b600080fd5b61016261011e366004610bde565b60026020526000908152604090205463ffffffff811690640100000000810462ffffff1690600160381b810461ffff1690600160481b90046001600160a01b031684565b6040805163ffffffff909516855262ffffff909316602085015261ffff909116918301919091526001600160a01b031660608201526080015b60405180910390f35b6101b76101b2366004610bde565b610356565b005b6101b76101c7366004610cc1565b610403565b6101df6101da366004610de5565b61045d565b60405165ffffffffffff909116815260200161019b565b610209610204366004610e6f565b61061a565b60405190815260200161019b565b6101b7610225366004610ebe565b61064b565b6101df610238366004610f8f565b600160209081526000928352604080842090915290825290205465ffffffffffff1681565b61020961026b366004610f8f565b600560209081526000928352604080842090915290825290205481565b6007546102a290600160a01b90046001600160401b031681565b6040516001600160401b03909116815260200161019b565b6101df6102c8366004610fbb565b6106d3565b6101b76102db366004610ffe565b6009805460ff1916911515919091179055565b6009546102fb9060ff1681565b604051901515815260200161019b565b610209600a5481565b6101b7610322366004611020565b61089f565b61020960085481565b61020961033e366004610e6f565b610aa9565b61020961035136600461104a565b610ac5565b6006546040516374d8461160e01b815265ffffffffffff831660048201526000916001600160a01b0316906374d84611906024016000604051808303816000875af11580156103a9573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526103d19190810190611063565b65ffffffffffff8316600090815260036020908152604090912082519293506103fe929091840190610b67565b505050565b60008160405160200161041691906110f3565b60408051601f1981840301815291815281516020928301206001600160401b0390961660009081526005835281812062ffffff909616815294909152909220929092555050565b600080600a544361046e919061113c565b9050600081600a54436104819190611166565b61048b919061117e565b60065460405163f645dcb160e01b81529192506000916001600160a01b039091169063f645dcb1906104c9908c908c908c908c908c90600401611195565b6020604051808303816000875af11580156104e8573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061050c919061121f565b600083815260016020908152604080832062ffffff8c1684529091528120805465ffffffffffff191665ffffffffffff8416179055600a5491925090610552908461123c565b6040805160808101825263ffffffff928316815262ffffff9a8b16602080830191825261ffff9d8e16838501908152306060850190815265ffffffffffff8916600090815260029093529490912092518354925191519451951666ffffffffffffff199092169190911764010000000091909c16029a909a17600160381b600160e81b031916600160381b91909b1602600160481b600160e81b03191699909917600160481b6001600160a01b03909a16999099029890981790965550939695505050505050565b6004602052816000526040600020818154811061063657600080fd5b90600052602060002001600091509150505481565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031633146106c85760405162461bcd60e51b815260206004820152601c60248201527f6f6e6c7920636f6f7264696e61746f722063616e2066756c66696c6c0000000060448201526064015b60405180910390fd5b6103fe838383610ae6565b600080600a54436106e4919061113c565b9050600081600a54436106f79190611166565b610701919061117e565b60065460405163dc92accf60e01b815261ffff891660048201526001600160401b038816602482015262ffffff871660448201529192506000916001600160a01b039091169063dc92accf906064016020604051808303816000875af115801561076f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610793919061121f565b600083815260016020908152604080832062ffffff8a1684529091528120805465ffffffffffff191665ffffffffffff8416179055600a54919250906107d9908461123c565b6040805160808101825263ffffffff928316815262ffffff988916602080830191825261ffff9c8d16838501908152306060850190815265ffffffffffff8916600090815260029093529490912092518354925191519451951666ffffffffffffff199092169190911764010000000091909a160298909817600160381b600160e81b031916600160381b91909a1602600160481b600160e81b03191698909817600160481b6001600160a01b0390991698909802979097179094555091949350505050565b6001600160401b038216600081815260016020908152604080832062ffffff8681168086529184528285205465ffffffffffff1680865260028552838620845160808082018752915463ffffffff808216835264010000000082048616838a0190815261ffff600160381b84048116858b019081526001600160a01b03600160481b909504851660608088019182529e8e5260058d528b8e209a8e52998c528a8d20548b519c8d0189905286519094169a8c019a909a5290519096169a89019a909a52955190931690860152915190921660a084015260c08301859052939092909160e0016040516020818303038152906040528051906020012090506000836040015161ffff166001600160401b038111156109be576109be610c31565b6040519080825280602002602001820160405280156109e7578160200160208202803683370190505b50905060005b846040015161ffff168161ffff161015610a76578281604051602001610a2a92919091825260f01b6001600160f01b031916602082015260220190565b6040516020818303038152906040528051906020012060001c828261ffff1681518110610a5957610a59611250565b602090810291909101015280610a6e81611266565b9150506109ed565b5065ffffffffffff851660009081526004602090815260409091208251610a9f92840190610b67565b5050505050505050565b6003602052816000526040600020818154811061063657600080fd5b60008181548110610ad557600080fd5b600091825260209091200154905081565b60095460ff1615610b395760405162461bcd60e51b815260206004820152601d60248201527f206661696c656420696e2066756c66696c6c52616e646f6d576f72647300000060448201526064016106bf565b65ffffffffffff831660009081526003602090815260409091208351610b6192850190610b67565b50505050565b828054828255906000526020600020908101928215610ba2579160200282015b82811115610ba2578251825591602001919060010190610b87565b50610bae929150610bb2565b5090565b5b80821115610bae5760008155600101610bb3565b65ffffffffffff81168114610bdb57600080fd5b50565b600060208284031215610bf057600080fd5b8135610bfb81610bc7565b9392505050565b80356001600160401b0381168114610c1957600080fd5b919050565b803562ffffff81168114610c1957600080fd5b634e487b7160e01b600052604160045260246000fd5b604051602081016001600160401b0381118282101715610c6957610c69610c31565b60405290565b604080519081016001600160401b0381118282101715610c6957610c69610c31565b604051601f8201601f191681016001600160401b0381118282101715610cb957610cb9610c31565b604052919050565b60008060008385036080811215610cd757600080fd5b610ce085610c02565b93506020610cef818701610c1e565b93506040603f1983011215610d0357600080fd5b610d0b610c47565b915086605f870112610d1c57600080fd5b610d24610c6f565b806080880189811115610d3657600080fd5b604089015b81811015610d525780358452928401928401610d3b565b50508352509396929550935090915050565b803561ffff81168114610c1957600080fd5b600082601f830112610d8757600080fd5b81356001600160401b03811115610da057610da0610c31565b610db3601f8201601f1916602001610c91565b818152846020838601011115610dc857600080fd5b816020850160208301376000918101602001919091529392505050565b600080600080600060a08688031215610dfd57600080fd5b610e0686610c02565b9450610e1460208701610d64565b9350610e2260408701610c1e565b9250606086013563ffffffff81168114610e3b57600080fd5b915060808601356001600160401b03811115610e5657600080fd5b610e6288828901610d76565b9150509295509295909350565b60008060408385031215610e8257600080fd5b8235610e8d81610bc7565b946020939093013593505050565b60006001600160401b03821115610eb457610eb4610c31565b5060051b60200190565b600080600060608486031215610ed357600080fd5b8335610ede81610bc7565b92506020848101356001600160401b0380821115610efb57600080fd5b818701915087601f830112610f0f57600080fd5b8135610f22610f1d82610e9b565b610c91565b81815260059190911b8301840190848101908a831115610f4157600080fd5b938501935b82851015610f5f57843582529385019390850190610f46565b965050506040870135925080831115610f7757600080fd5b5050610f8586828701610d76565b9150509250925092565b60008060408385031215610fa257600080fd5b82359150610fb260208401610c1e565b90509250929050565b600080600060608486031215610fd057600080fd5b610fd984610d64565b9250610fe760208501610c02565b9150610ff560408501610c1e565b90509250925092565b60006020828403121561101057600080fd5b81358015158114610bfb57600080fd5b6000806040838503121561103357600080fd5b61103c83610c02565b9150610fb260208401610c1e565b60006020828403121561105c57600080fd5b5035919050565b6000602080838503121561107657600080fd5b82516001600160401b0381111561108c57600080fd5b8301601f8101851361109d57600080fd5b80516110ab610f1d82610e9b565b81815260059190911b820183019083810190878311156110ca57600080fd5b928401925b828410156110e8578351825292840192908401906110cf565b979650505050505050565b815160408201908260005b600281101561111d5782518252602092830192909101906001016110fe565b50505092915050565b634e487b7160e01b600052601260045260246000fd5b60008261114b5761114b611126565b500690565b634e487b7160e01b600052601160045260246000fd5b6000821982111561117957611179611150565b500190565b60008282101561119057611190611150565b500390565b6001600160401b03861681526000602061ffff87168184015262ffffff8616604084015263ffffffff8516606084015260a0608084015283518060a085015260005b818110156111f35785810183015185820160c0015282016111d7565b8181111561120557600060c083870101525b50601f01601f19169290920160c001979650505050505050565b60006020828403121561123157600080fd5b8151610bfb81610bc7565b60008261124b5761124b611126565b500490565b634e487b7160e01b600052603260045260246000fd5b600061ffff80831681810361127d5761127d611150565b600101939250505056fea26469706673582212205d99391252c2c91358d7cf65659893f205b174315109c679dca1d670e2b54c9b64736f6c634300080f0033",
}

var BeaconVRFConsumerABI = BeaconVRFConsumerMetaData.ABI

var BeaconVRFConsumerBin = BeaconVRFConsumerMetaData.Bin

func DeployBeaconVRFConsumer(auth *bind.TransactOpts, backend bind.ContractBackend, coordinator common.Address, shouldFail bool, beaconPeriodBlocks *big.Int) (common.Address, *types.Transaction, *BeaconVRFConsumer, error) {
	parsed, err := BeaconVRFConsumerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BeaconVRFConsumerBin), backend, coordinator, shouldFail, beaconPeriodBlocks)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BeaconVRFConsumer{BeaconVRFConsumerCaller: BeaconVRFConsumerCaller{contract: contract}, BeaconVRFConsumerTransactor: BeaconVRFConsumerTransactor{contract: contract}, BeaconVRFConsumerFilterer: BeaconVRFConsumerFilterer{contract: contract}}, nil
}

type BeaconVRFConsumer struct {
	address common.Address
	abi     abi.ABI
	BeaconVRFConsumerCaller
	BeaconVRFConsumerTransactor
	BeaconVRFConsumerFilterer
}

type BeaconVRFConsumerCaller struct {
	contract *bind.BoundContract
}

type BeaconVRFConsumerTransactor struct {
	contract *bind.BoundContract
}

type BeaconVRFConsumerFilterer struct {
	contract *bind.BoundContract
}

type BeaconVRFConsumerSession struct {
	Contract     *BeaconVRFConsumer
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type BeaconVRFConsumerCallerSession struct {
	Contract *BeaconVRFConsumerCaller
	CallOpts bind.CallOpts
}

type BeaconVRFConsumerTransactorSession struct {
	Contract     *BeaconVRFConsumerTransactor
	TransactOpts bind.TransactOpts
}

type BeaconVRFConsumerRaw struct {
	Contract *BeaconVRFConsumer
}

type BeaconVRFConsumerCallerRaw struct {
	Contract *BeaconVRFConsumerCaller
}

type BeaconVRFConsumerTransactorRaw struct {
	Contract *BeaconVRFConsumerTransactor
}

func NewBeaconVRFConsumer(address common.Address, backend bind.ContractBackend) (*BeaconVRFConsumer, error) {
	abi, err := abi.JSON(strings.NewReader(BeaconVRFConsumerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindBeaconVRFConsumer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BeaconVRFConsumer{address: address, abi: abi, BeaconVRFConsumerCaller: BeaconVRFConsumerCaller{contract: contract}, BeaconVRFConsumerTransactor: BeaconVRFConsumerTransactor{contract: contract}, BeaconVRFConsumerFilterer: BeaconVRFConsumerFilterer{contract: contract}}, nil
}

func NewBeaconVRFConsumerCaller(address common.Address, caller bind.ContractCaller) (*BeaconVRFConsumerCaller, error) {
	contract, err := bindBeaconVRFConsumer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BeaconVRFConsumerCaller{contract: contract}, nil
}

func NewBeaconVRFConsumerTransactor(address common.Address, transactor bind.ContractTransactor) (*BeaconVRFConsumerTransactor, error) {
	contract, err := bindBeaconVRFConsumer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BeaconVRFConsumerTransactor{contract: contract}, nil
}

func NewBeaconVRFConsumerFilterer(address common.Address, filterer bind.ContractFilterer) (*BeaconVRFConsumerFilterer, error) {
	contract, err := bindBeaconVRFConsumer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BeaconVRFConsumerFilterer{contract: contract}, nil
}

func bindBeaconVRFConsumer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BeaconVRFConsumerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_BeaconVRFConsumer *BeaconVRFConsumerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BeaconVRFConsumer.Contract.BeaconVRFConsumerCaller.contract.Call(opts, result, method, params...)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.BeaconVRFConsumerTransactor.contract.Transfer(opts)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.BeaconVRFConsumerTransactor.contract.Transact(opts, method, params...)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BeaconVRFConsumer.Contract.contract.Call(opts, result, method, params...)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.contract.Transfer(opts)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.contract.Transact(opts, method, params...)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCaller) Fail(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _BeaconVRFConsumer.contract.Call(opts, &out, "fail")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BeaconVRFConsumer *BeaconVRFConsumerSession) Fail() (bool, error) {
	return _BeaconVRFConsumer.Contract.Fail(&_BeaconVRFConsumer.CallOpts)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCallerSession) Fail() (bool, error) {
	return _BeaconVRFConsumer.Contract.Fail(&_BeaconVRFConsumer.CallOpts)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCaller) IBeaconPeriodBlocks(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BeaconVRFConsumer.contract.Call(opts, &out, "i_beaconPeriodBlocks")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_BeaconVRFConsumer *BeaconVRFConsumerSession) IBeaconPeriodBlocks() (*big.Int, error) {
	return _BeaconVRFConsumer.Contract.IBeaconPeriodBlocks(&_BeaconVRFConsumer.CallOpts)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCallerSession) IBeaconPeriodBlocks() (*big.Int, error) {
	return _BeaconVRFConsumer.Contract.IBeaconPeriodBlocks(&_BeaconVRFConsumer.CallOpts)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCaller) SExpectedRandomnessByBlock(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _BeaconVRFConsumer.contract.Call(opts, &out, "s_ExpectedRandomnessByBlock", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_BeaconVRFConsumer *BeaconVRFConsumerSession) SExpectedRandomnessByBlock(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _BeaconVRFConsumer.Contract.SExpectedRandomnessByBlock(&_BeaconVRFConsumer.CallOpts, arg0, arg1)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCallerSession) SExpectedRandomnessByBlock(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _BeaconVRFConsumer.Contract.SExpectedRandomnessByBlock(&_BeaconVRFConsumer.CallOpts, arg0, arg1)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCaller) SExpectedSeeds(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _BeaconVRFConsumer.contract.Call(opts, &out, "s_ExpectedSeeds", arg0, arg1)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_BeaconVRFConsumer *BeaconVRFConsumerSession) SExpectedSeeds(arg0 *big.Int, arg1 *big.Int) ([32]byte, error) {
	return _BeaconVRFConsumer.Contract.SExpectedSeeds(&_BeaconVRFConsumer.CallOpts, arg0, arg1)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCallerSession) SExpectedSeeds(arg0 *big.Int, arg1 *big.Int) ([32]byte, error) {
	return _BeaconVRFConsumer.Contract.SExpectedSeeds(&_BeaconVRFConsumer.CallOpts, arg0, arg1)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCaller) SReceivedRandomnessByRequestID(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _BeaconVRFConsumer.contract.Call(opts, &out, "s_ReceivedRandomnessByRequestID", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_BeaconVRFConsumer *BeaconVRFConsumerSession) SReceivedRandomnessByRequestID(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _BeaconVRFConsumer.Contract.SReceivedRandomnessByRequestID(&_BeaconVRFConsumer.CallOpts, arg0, arg1)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCallerSession) SReceivedRandomnessByRequestID(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _BeaconVRFConsumer.Contract.SReceivedRandomnessByRequestID(&_BeaconVRFConsumer.CallOpts, arg0, arg1)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCaller) SGasAvailable(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BeaconVRFConsumer.contract.Call(opts, &out, "s_gasAvailable")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_BeaconVRFConsumer *BeaconVRFConsumerSession) SGasAvailable() (*big.Int, error) {
	return _BeaconVRFConsumer.Contract.SGasAvailable(&_BeaconVRFConsumer.CallOpts)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCallerSession) SGasAvailable() (*big.Int, error) {
	return _BeaconVRFConsumer.Contract.SGasAvailable(&_BeaconVRFConsumer.CallOpts)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCaller) SMyBeaconRequests(opts *bind.CallOpts, arg0 *big.Int) (SMyBeaconRequests,

	error) {
	var out []interface{}
	err := _BeaconVRFConsumer.contract.Call(opts, &out, "s_myBeaconRequests", arg0)

	outstruct := new(SMyBeaconRequests)
	if err != nil {
		return *outstruct, err
	}

	outstruct.SlotNumber = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.ConfirmationDelay = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.NumWords = *abi.ConvertType(out[2], new(uint16)).(*uint16)
	outstruct.Requester = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)

	return *outstruct, err

}

func (_BeaconVRFConsumer *BeaconVRFConsumerSession) SMyBeaconRequests(arg0 *big.Int) (SMyBeaconRequests,

	error) {
	return _BeaconVRFConsumer.Contract.SMyBeaconRequests(&_BeaconVRFConsumer.CallOpts, arg0)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCallerSession) SMyBeaconRequests(arg0 *big.Int) (SMyBeaconRequests,

	error) {
	return _BeaconVRFConsumer.Contract.SMyBeaconRequests(&_BeaconVRFConsumer.CallOpts, arg0)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCaller) SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _BeaconVRFConsumer.contract.Call(opts, &out, "s_randomWords", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_BeaconVRFConsumer *BeaconVRFConsumerSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _BeaconVRFConsumer.Contract.SRandomWords(&_BeaconVRFConsumer.CallOpts, arg0)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCallerSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _BeaconVRFConsumer.Contract.SRandomWords(&_BeaconVRFConsumer.CallOpts, arg0)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCaller) SRequestsIDs(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _BeaconVRFConsumer.contract.Call(opts, &out, "s_requestsIDs", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_BeaconVRFConsumer *BeaconVRFConsumerSession) SRequestsIDs(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _BeaconVRFConsumer.Contract.SRequestsIDs(&_BeaconVRFConsumer.CallOpts, arg0, arg1)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCallerSession) SRequestsIDs(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _BeaconVRFConsumer.Contract.SRequestsIDs(&_BeaconVRFConsumer.CallOpts, arg0, arg1)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCaller) SSubId(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _BeaconVRFConsumer.contract.Call(opts, &out, "s_subId")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_BeaconVRFConsumer *BeaconVRFConsumerSession) SSubId() (uint64, error) {
	return _BeaconVRFConsumer.Contract.SSubId(&_BeaconVRFConsumer.CallOpts)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerCallerSession) SSubId() (uint64, error) {
	return _BeaconVRFConsumer.Contract.SSubId(&_BeaconVRFConsumer.CallOpts)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerTransactor) ComputeAndStoreExpectedOutput(opts *bind.TransactOpts, height uint64, confirmationDelayArg *big.Int) (*types.Transaction, error) {
	return _BeaconVRFConsumer.contract.Transact(opts, "computeAndStoreExpectedOutput", height, confirmationDelayArg)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerSession) ComputeAndStoreExpectedOutput(height uint64, confirmationDelayArg *big.Int) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.ComputeAndStoreExpectedOutput(&_BeaconVRFConsumer.TransactOpts, height, confirmationDelayArg)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerTransactorSession) ComputeAndStoreExpectedOutput(height uint64, confirmationDelayArg *big.Int) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.ComputeAndStoreExpectedOutput(&_BeaconVRFConsumer.TransactOpts, height, confirmationDelayArg)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestID *big.Int, randomWords []*big.Int, arguments []byte) (*types.Transaction, error) {
	return _BeaconVRFConsumer.contract.Transact(opts, "rawFulfillRandomWords", requestID, randomWords, arguments)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerSession) RawFulfillRandomWords(requestID *big.Int, randomWords []*big.Int, arguments []byte) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.RawFulfillRandomWords(&_BeaconVRFConsumer.TransactOpts, requestID, randomWords, arguments)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerTransactorSession) RawFulfillRandomWords(requestID *big.Int, randomWords []*big.Int, arguments []byte) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.RawFulfillRandomWords(&_BeaconVRFConsumer.TransactOpts, requestID, randomWords, arguments)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerTransactor) SetExpectedSeed(opts *bind.TransactOpts, height uint64, confirmationDelayArg *big.Int, vrfOutput ECCArithmeticG1Point) (*types.Transaction, error) {
	return _BeaconVRFConsumer.contract.Transact(opts, "setExpectedSeed", height, confirmationDelayArg, vrfOutput)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerSession) SetExpectedSeed(height uint64, confirmationDelayArg *big.Int, vrfOutput ECCArithmeticG1Point) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.SetExpectedSeed(&_BeaconVRFConsumer.TransactOpts, height, confirmationDelayArg, vrfOutput)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerTransactorSession) SetExpectedSeed(height uint64, confirmationDelayArg *big.Int, vrfOutput ECCArithmeticG1Point) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.SetExpectedSeed(&_BeaconVRFConsumer.TransactOpts, height, confirmationDelayArg, vrfOutput)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerTransactor) SetFail(opts *bind.TransactOpts, shouldFail bool) (*types.Transaction, error) {
	return _BeaconVRFConsumer.contract.Transact(opts, "setFail", shouldFail)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerSession) SetFail(shouldFail bool) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.SetFail(&_BeaconVRFConsumer.TransactOpts, shouldFail)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerTransactorSession) SetFail(shouldFail bool) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.SetFail(&_BeaconVRFConsumer.TransactOpts, shouldFail)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerTransactor) TestRedeemRandomness(opts *bind.TransactOpts, requestID *big.Int) (*types.Transaction, error) {
	return _BeaconVRFConsumer.contract.Transact(opts, "testRedeemRandomness", requestID)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerSession) TestRedeemRandomness(requestID *big.Int) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.TestRedeemRandomness(&_BeaconVRFConsumer.TransactOpts, requestID)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerTransactorSession) TestRedeemRandomness(requestID *big.Int) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.TestRedeemRandomness(&_BeaconVRFConsumer.TransactOpts, requestID)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerTransactor) TestRequestRandomness(opts *bind.TransactOpts, numWords uint16, subID uint64, confirmationDelayArg *big.Int) (*types.Transaction, error) {
	return _BeaconVRFConsumer.contract.Transact(opts, "testRequestRandomness", numWords, subID, confirmationDelayArg)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerSession) TestRequestRandomness(numWords uint16, subID uint64, confirmationDelayArg *big.Int) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.TestRequestRandomness(&_BeaconVRFConsumer.TransactOpts, numWords, subID, confirmationDelayArg)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerTransactorSession) TestRequestRandomness(numWords uint16, subID uint64, confirmationDelayArg *big.Int) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.TestRequestRandomness(&_BeaconVRFConsumer.TransactOpts, numWords, subID, confirmationDelayArg)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerTransactor) TestRequestRandomnessFulfillment(opts *bind.TransactOpts, subID uint64, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte) (*types.Transaction, error) {
	return _BeaconVRFConsumer.contract.Transact(opts, "testRequestRandomnessFulfillment", subID, numWords, confirmationDelayArg, callbackGasLimit, arguments)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerSession) TestRequestRandomnessFulfillment(subID uint64, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.TestRequestRandomnessFulfillment(&_BeaconVRFConsumer.TransactOpts, subID, numWords, confirmationDelayArg, callbackGasLimit, arguments)
}

func (_BeaconVRFConsumer *BeaconVRFConsumerTransactorSession) TestRequestRandomnessFulfillment(subID uint64, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte) (*types.Transaction, error) {
	return _BeaconVRFConsumer.Contract.TestRequestRandomnessFulfillment(&_BeaconVRFConsumer.TransactOpts, subID, numWords, confirmationDelayArg, callbackGasLimit, arguments)
}

type SMyBeaconRequests struct {
	SlotNumber        uint32
	ConfirmationDelay *big.Int
	NumWords          uint16
	Requester         common.Address
}

func (_BeaconVRFConsumer *BeaconVRFConsumer) Address() common.Address {
	return _BeaconVRFConsumer.address
}

type BeaconVRFConsumerInterface interface {
	Fail(opts *bind.CallOpts) (bool, error)

	IBeaconPeriodBlocks(opts *bind.CallOpts) (*big.Int, error)

	SExpectedRandomnessByBlock(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error)

	SExpectedSeeds(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) ([32]byte, error)

	SReceivedRandomnessByRequestID(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error)

	SGasAvailable(opts *bind.CallOpts) (*big.Int, error)

	SMyBeaconRequests(opts *bind.CallOpts, arg0 *big.Int) (SMyBeaconRequests,

		error)

	SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	SRequestsIDs(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error)

	SSubId(opts *bind.CallOpts) (uint64, error)

	ComputeAndStoreExpectedOutput(opts *bind.TransactOpts, height uint64, confirmationDelayArg *big.Int) (*types.Transaction, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestID *big.Int, randomWords []*big.Int, arguments []byte) (*types.Transaction, error)

	SetExpectedSeed(opts *bind.TransactOpts, height uint64, confirmationDelayArg *big.Int, vrfOutput ECCArithmeticG1Point) (*types.Transaction, error)

	SetFail(opts *bind.TransactOpts, shouldFail bool) (*types.Transaction, error)

	TestRedeemRandomness(opts *bind.TransactOpts, requestID *big.Int) (*types.Transaction, error)

	TestRequestRandomness(opts *bind.TransactOpts, numWords uint16, subID uint64, confirmationDelayArg *big.Int) (*types.Transaction, error)

	TestRequestRandomnessFulfillment(opts *bind.TransactOpts, subID uint64, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte) (*types.Transaction, error)

	Address() common.Address
}
