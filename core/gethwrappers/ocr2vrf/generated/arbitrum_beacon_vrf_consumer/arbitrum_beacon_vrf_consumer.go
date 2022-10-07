// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package arbitrum_beacon_vrf_consumer

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

var ArbitrumBeaconVRFConsumerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"shouldFail\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"beaconPeriodBlocks\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ARBSYS\",\"outputs\":[{\"internalType\":\"contractArbSys\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"NUM_CONF_DELAYS\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"fail\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_beaconPeriodBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_ReceivedRandomnessByRequestID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_gasAvailable\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"}],\"name\":\"s_myBeaconRequests\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.SlotNumber\",\"name\":\"slotNumber\",\"type\":\"uint32\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelay\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_randomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"\",\"type\":\"uint24\"}],\"name\":\"s_requestsIDs\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_subId\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"shouldFail\",\"type\":\"bool\"}],\"name\":\"setFail\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"reqId\",\"type\":\"uint48\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"delay\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"}],\"name\":\"storeBeaconRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"}],\"name\":\"testRedeemRandomness\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"uint64\",\"name\":\"subID\",\"type\":\"uint64\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelayArg\",\"type\":\"uint24\"}],\"name\":\"testRequestRandomness\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subID\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelayArg\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"}],\"name\":\"testRequestRandomnessFulfillment\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561001057600080fd5b5060405161108038038061108083398101604081905261002f9161007f565b6001600160a01b039092166080819052600480546001600160a01b031990811690921790556007805492151560ff19909316929092179091556008919091556009805490911660641790556100d3565b60008060006060848603121561009457600080fd5b83516001600160a01b03811681146100ab57600080fd5b602085015190935080151581146100c157600080fd5b80925050604084015190509250925092565b608051610f926100ee60003960006105860152610f926000f3fe608060405234801561001057600080fd5b50600436106101005760003560e01c80639c9cd01511610097578063cd0593df11610066578063cd0593df14610316578063f08c5daa1461032d578063f339c79414610336578063f6eaffc81461034957600080fd5b80639c9cd0151461029a5780639d769402146102ad578063a9cc4718146102ce578063bf0a12cf146102eb57600080fd5b80635a47dd71116100d35780635a47dd71146101fc5780635f15cccc1461020f5780636d162a3e14610242578063706da1ca1461025557600080fd5b806319a5fa22146101055780632f7527cc146101a35780633d8b70aa146101bd57806345907626146101d2575b600080fd5b6101616101133660046109f8565b60026020526000908152604090205463ffffffff811690640100000000810462ffffff1690670100000000000000810461ffff1690690100000000000000000090046001600160a01b031684565b6040805163ffffffff909516855262ffffff909316602085015261ffff909116918301919091526001600160a01b031660608201526080015b60405180910390f35b6101ab600881565b60405160ff909116815260200161019a565b6101d06101cb3660046109f8565b61035c565b005b6101e56101e0366004610b15565b610427565b60405165ffffffffffff909116815260200161019a565b6101d061020a366004610bc4565b610584565b6101e561021d366004610c96565b600160209081526000928352604080842090915290825290205465ffffffffffff1681565b6101d0610250366004610cc2565b61060c565b6005546102819074010000000000000000000000000000000000000000900467ffffffffffffffff1681565b60405167ffffffffffffffff909116815260200161019a565b6101e56102a8366004610d11565b61073f565b6101d06102bb366004610d54565b6007805460ff1916911515919091179055565b6007546102db9060ff1681565b604051901515815260200161019a565b6009546102fe906001600160a01b031681565b6040516001600160a01b03909116815260200161019a565b61031f60085481565b60405190815260200161019a565b61031f60065481565b61031f610344366004610d76565b6108ae565b61031f610357366004610da2565b6108df565b600480546040517f74d8461100000000000000000000000000000000000000000000000000000000815265ffffffffffff8416928101929092526000916001600160a01b03909116906374d84611906024016000604051808303816000875af11580156103cd573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526103f59190810190610dbb565b65ffffffffffff831660009081526003602090815260409091208251929350610422929091840190610981565b505050565b600080600960009054906101000a90046001600160a01b03166001600160a01b031663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa15801561047d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104a19190610e41565b90506000600854826104b39190610e70565b9050600081600854846104c69190610e9a565b6104d09190610eb2565b600480546040517ff645dcb10000000000000000000000000000000000000000000000000000000081529293506000926001600160a01b039091169163f645dcb191610526918e918e918e918e918e9101610ec9565b6020604051808303816000875af1158015610545573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105699190610f54565b905061057781838a8c61060c565b9998505050505050505050565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031633146106015760405162461bcd60e51b815260206004820152601c60248201527f6f6e6c7920636f6f7264696e61746f722063616e2066756c66696c6c0000000060448201526064015b60405180910390fd5b610422838383610900565b600083815260016020908152604080832062ffffff861684529091528120805465ffffffffffff191665ffffffffffff871617905560085461064e9085610f71565b6040805160808101825263ffffffff928316815262ffffff958616602080830191825261ffff968716838501908152306060850190815265ffffffffffff909b1660009081526002909252939020915182549151935199516001600160a01b03166901000000000000000000027fffffff0000000000000000000000000000000000000000ffffffffffffffffff9a90971667010000000000000002999099167fffffff00000000000000000000000000000000000000000000ffffffffffffff939097166401000000000266ffffffffffffff199091169890931697909717919091171692909217179092555050565b600080600960009054906101000a90046001600160a01b03166001600160a01b031663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610795573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107b99190610e41565b90506000600854826107cb9190610e70565b9050600081600854846107de9190610e9a565b6107e89190610eb2565b600480546040517fdc92accf00000000000000000000000000000000000000000000000000000000815261ffff8b169281019290925267ffffffffffffffff8916602483015262ffffff881660448301529192506000916001600160a01b03169063dc92accf906064016020604051808303816000875af1158015610871573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108959190610f54565b90506108a38183888b61060c565b979650505050505050565b600360205281600052604060002081815481106108ca57600080fd5b90600052602060002001600091509150505481565b600081815481106108ef57600080fd5b600091825260209091200154905081565b60075460ff16156109535760405162461bcd60e51b815260206004820152601d60248201527f206661696c656420696e2066756c66696c6c52616e646f6d576f72647300000060448201526064016105f8565b65ffffffffffff83166000908152600360209081526040909120835161097b92850190610981565b50505050565b8280548282559060005260206000209081019282156109bc579160200282015b828111156109bc5782518255916020019190600101906109a1565b506109c89291506109cc565b5090565b5b808211156109c857600081556001016109cd565b65ffffffffffff811681146109f557600080fd5b50565b600060208284031215610a0a57600080fd5b8135610a15816109e1565b9392505050565b803567ffffffffffffffff81168114610a3457600080fd5b919050565b803561ffff81168114610a3457600080fd5b803562ffffff81168114610a3457600080fd5b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff81118282101715610a9d57610a9d610a5e565b604052919050565b600082601f830112610ab657600080fd5b813567ffffffffffffffff811115610ad057610ad0610a5e565b610ae3601f8201601f1916602001610a74565b818152846020838601011115610af857600080fd5b816020850160208301376000918101602001919091529392505050565b600080600080600060a08688031215610b2d57600080fd5b610b3686610a1c565b9450610b4460208701610a39565b9350610b5260408701610a4b565b9250606086013563ffffffff81168114610b6b57600080fd5b9150608086013567ffffffffffffffff811115610b8757600080fd5b610b9388828901610aa5565b9150509295509295909350565b600067ffffffffffffffff821115610bba57610bba610a5e565b5060051b60200190565b600080600060608486031215610bd957600080fd5b8335610be4816109e1565b925060208481013567ffffffffffffffff80821115610c0257600080fd5b818701915087601f830112610c1657600080fd5b8135610c29610c2482610ba0565b610a74565b81815260059190911b8301840190848101908a831115610c4857600080fd5b938501935b82851015610c6657843582529385019390850190610c4d565b965050506040870135925080831115610c7e57600080fd5b5050610c8c86828701610aa5565b9150509250925092565b60008060408385031215610ca957600080fd5b82359150610cb960208401610a4b565b90509250929050565b60008060008060808587031215610cd857600080fd5b8435610ce3816109e1565b935060208501359250610cf860408601610a4b565b9150610d0660608601610a39565b905092959194509250565b600080600060608486031215610d2657600080fd5b610d2f84610a39565b9250610d3d60208501610a1c565b9150610d4b60408501610a4b565b90509250925092565b600060208284031215610d6657600080fd5b81358015158114610a1557600080fd5b60008060408385031215610d8957600080fd5b8235610d94816109e1565b946020939093013593505050565b600060208284031215610db457600080fd5b5035919050565b60006020808385031215610dce57600080fd5b825167ffffffffffffffff811115610de557600080fd5b8301601f81018513610df657600080fd5b8051610e04610c2482610ba0565b81815260059190911b82018301908381019087831115610e2357600080fd5b928401925b828410156108a357835182529284019290840190610e28565b600060208284031215610e5357600080fd5b5051919050565b634e487b7160e01b600052601260045260246000fd5b600082610e7f57610e7f610e5a565b500690565b634e487b7160e01b600052601160045260246000fd5b60008219821115610ead57610ead610e84565b500190565b600082821015610ec457610ec4610e84565b500390565b67ffffffffffffffff861681526000602061ffff87168184015262ffffff8616604084015263ffffffff8516606084015260a0608084015283518060a085015260005b81811015610f285785810183015185820160c001528201610f0c565b81811115610f3a57600060c083870101525b50601f01601f19169290920160c001979650505050505050565b600060208284031215610f6657600080fd5b8151610a15816109e1565b600082610f8057610f80610e5a565b50049056fea164736f6c634300080f000a",
}

var ArbitrumBeaconVRFConsumerABI = ArbitrumBeaconVRFConsumerMetaData.ABI

var ArbitrumBeaconVRFConsumerBin = ArbitrumBeaconVRFConsumerMetaData.Bin

func DeployArbitrumBeaconVRFConsumer(auth *bind.TransactOpts, backend bind.ContractBackend, coordinator common.Address, shouldFail bool, beaconPeriodBlocks *big.Int) (common.Address, *types.Transaction, *ArbitrumBeaconVRFConsumer, error) {
	parsed, err := ArbitrumBeaconVRFConsumerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ArbitrumBeaconVRFConsumerBin), backend, coordinator, shouldFail, beaconPeriodBlocks)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ArbitrumBeaconVRFConsumer{ArbitrumBeaconVRFConsumerCaller: ArbitrumBeaconVRFConsumerCaller{contract: contract}, ArbitrumBeaconVRFConsumerTransactor: ArbitrumBeaconVRFConsumerTransactor{contract: contract}, ArbitrumBeaconVRFConsumerFilterer: ArbitrumBeaconVRFConsumerFilterer{contract: contract}}, nil
}

type ArbitrumBeaconVRFConsumer struct {
	address common.Address
	abi     abi.ABI
	ArbitrumBeaconVRFConsumerCaller
	ArbitrumBeaconVRFConsumerTransactor
	ArbitrumBeaconVRFConsumerFilterer
}

type ArbitrumBeaconVRFConsumerCaller struct {
	contract *bind.BoundContract
}

type ArbitrumBeaconVRFConsumerTransactor struct {
	contract *bind.BoundContract
}

type ArbitrumBeaconVRFConsumerFilterer struct {
	contract *bind.BoundContract
}

type ArbitrumBeaconVRFConsumerSession struct {
	Contract     *ArbitrumBeaconVRFConsumer
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type ArbitrumBeaconVRFConsumerCallerSession struct {
	Contract *ArbitrumBeaconVRFConsumerCaller
	CallOpts bind.CallOpts
}

type ArbitrumBeaconVRFConsumerTransactorSession struct {
	Contract     *ArbitrumBeaconVRFConsumerTransactor
	TransactOpts bind.TransactOpts
}

type ArbitrumBeaconVRFConsumerRaw struct {
	Contract *ArbitrumBeaconVRFConsumer
}

type ArbitrumBeaconVRFConsumerCallerRaw struct {
	Contract *ArbitrumBeaconVRFConsumerCaller
}

type ArbitrumBeaconVRFConsumerTransactorRaw struct {
	Contract *ArbitrumBeaconVRFConsumerTransactor
}

func NewArbitrumBeaconVRFConsumer(address common.Address, backend bind.ContractBackend) (*ArbitrumBeaconVRFConsumer, error) {
	abi, err := abi.JSON(strings.NewReader(ArbitrumBeaconVRFConsumerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindArbitrumBeaconVRFConsumer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ArbitrumBeaconVRFConsumer{address: address, abi: abi, ArbitrumBeaconVRFConsumerCaller: ArbitrumBeaconVRFConsumerCaller{contract: contract}, ArbitrumBeaconVRFConsumerTransactor: ArbitrumBeaconVRFConsumerTransactor{contract: contract}, ArbitrumBeaconVRFConsumerFilterer: ArbitrumBeaconVRFConsumerFilterer{contract: contract}}, nil
}

func NewArbitrumBeaconVRFConsumerCaller(address common.Address, caller bind.ContractCaller) (*ArbitrumBeaconVRFConsumerCaller, error) {
	contract, err := bindArbitrumBeaconVRFConsumer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ArbitrumBeaconVRFConsumerCaller{contract: contract}, nil
}

func NewArbitrumBeaconVRFConsumerTransactor(address common.Address, transactor bind.ContractTransactor) (*ArbitrumBeaconVRFConsumerTransactor, error) {
	contract, err := bindArbitrumBeaconVRFConsumer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ArbitrumBeaconVRFConsumerTransactor{contract: contract}, nil
}

func NewArbitrumBeaconVRFConsumerFilterer(address common.Address, filterer bind.ContractFilterer) (*ArbitrumBeaconVRFConsumerFilterer, error) {
	contract, err := bindArbitrumBeaconVRFConsumer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ArbitrumBeaconVRFConsumerFilterer{contract: contract}, nil
}

func bindArbitrumBeaconVRFConsumer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ArbitrumBeaconVRFConsumerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ArbitrumBeaconVRFConsumer.Contract.ArbitrumBeaconVRFConsumerCaller.contract.Call(opts, result, method, params...)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ArbitrumBeaconVRFConsumer.Contract.ArbitrumBeaconVRFConsumerTransactor.contract.Transfer(opts)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ArbitrumBeaconVRFConsumer.Contract.ArbitrumBeaconVRFConsumerTransactor.contract.Transact(opts, method, params...)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ArbitrumBeaconVRFConsumer.Contract.contract.Call(opts, result, method, params...)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ArbitrumBeaconVRFConsumer.Contract.contract.Transfer(opts)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ArbitrumBeaconVRFConsumer.Contract.contract.Transact(opts, method, params...)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerCaller) ARBSYS(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ArbitrumBeaconVRFConsumer.contract.Call(opts, &out, "ARBSYS")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerSession) ARBSYS() (common.Address, error) {
	return _ArbitrumBeaconVRFConsumer.Contract.ARBSYS(&_ArbitrumBeaconVRFConsumer.CallOpts)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerCallerSession) ARBSYS() (common.Address, error) {
	return _ArbitrumBeaconVRFConsumer.Contract.ARBSYS(&_ArbitrumBeaconVRFConsumer.CallOpts)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerCaller) NUMCONFDELAYS(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _ArbitrumBeaconVRFConsumer.contract.Call(opts, &out, "NUM_CONF_DELAYS")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerSession) NUMCONFDELAYS() (uint8, error) {
	return _ArbitrumBeaconVRFConsumer.Contract.NUMCONFDELAYS(&_ArbitrumBeaconVRFConsumer.CallOpts)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerCallerSession) NUMCONFDELAYS() (uint8, error) {
	return _ArbitrumBeaconVRFConsumer.Contract.NUMCONFDELAYS(&_ArbitrumBeaconVRFConsumer.CallOpts)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerCaller) Fail(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _ArbitrumBeaconVRFConsumer.contract.Call(opts, &out, "fail")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerSession) Fail() (bool, error) {
	return _ArbitrumBeaconVRFConsumer.Contract.Fail(&_ArbitrumBeaconVRFConsumer.CallOpts)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerCallerSession) Fail() (bool, error) {
	return _ArbitrumBeaconVRFConsumer.Contract.Fail(&_ArbitrumBeaconVRFConsumer.CallOpts)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerCaller) IBeaconPeriodBlocks(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ArbitrumBeaconVRFConsumer.contract.Call(opts, &out, "i_beaconPeriodBlocks")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerSession) IBeaconPeriodBlocks() (*big.Int, error) {
	return _ArbitrumBeaconVRFConsumer.Contract.IBeaconPeriodBlocks(&_ArbitrumBeaconVRFConsumer.CallOpts)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerCallerSession) IBeaconPeriodBlocks() (*big.Int, error) {
	return _ArbitrumBeaconVRFConsumer.Contract.IBeaconPeriodBlocks(&_ArbitrumBeaconVRFConsumer.CallOpts)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerCaller) SReceivedRandomnessByRequestID(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _ArbitrumBeaconVRFConsumer.contract.Call(opts, &out, "s_ReceivedRandomnessByRequestID", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerSession) SReceivedRandomnessByRequestID(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _ArbitrumBeaconVRFConsumer.Contract.SReceivedRandomnessByRequestID(&_ArbitrumBeaconVRFConsumer.CallOpts, arg0, arg1)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerCallerSession) SReceivedRandomnessByRequestID(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _ArbitrumBeaconVRFConsumer.Contract.SReceivedRandomnessByRequestID(&_ArbitrumBeaconVRFConsumer.CallOpts, arg0, arg1)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerCaller) SGasAvailable(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ArbitrumBeaconVRFConsumer.contract.Call(opts, &out, "s_gasAvailable")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerSession) SGasAvailable() (*big.Int, error) {
	return _ArbitrumBeaconVRFConsumer.Contract.SGasAvailable(&_ArbitrumBeaconVRFConsumer.CallOpts)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerCallerSession) SGasAvailable() (*big.Int, error) {
	return _ArbitrumBeaconVRFConsumer.Contract.SGasAvailable(&_ArbitrumBeaconVRFConsumer.CallOpts)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerCaller) SMyBeaconRequests(opts *bind.CallOpts, arg0 *big.Int) (SMyBeaconRequests,

	error) {
	var out []interface{}
	err := _ArbitrumBeaconVRFConsumer.contract.Call(opts, &out, "s_myBeaconRequests", arg0)

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

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerSession) SMyBeaconRequests(arg0 *big.Int) (SMyBeaconRequests,

	error) {
	return _ArbitrumBeaconVRFConsumer.Contract.SMyBeaconRequests(&_ArbitrumBeaconVRFConsumer.CallOpts, arg0)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerCallerSession) SMyBeaconRequests(arg0 *big.Int) (SMyBeaconRequests,

	error) {
	return _ArbitrumBeaconVRFConsumer.Contract.SMyBeaconRequests(&_ArbitrumBeaconVRFConsumer.CallOpts, arg0)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerCaller) SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _ArbitrumBeaconVRFConsumer.contract.Call(opts, &out, "s_randomWords", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _ArbitrumBeaconVRFConsumer.Contract.SRandomWords(&_ArbitrumBeaconVRFConsumer.CallOpts, arg0)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerCallerSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _ArbitrumBeaconVRFConsumer.Contract.SRandomWords(&_ArbitrumBeaconVRFConsumer.CallOpts, arg0)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerCaller) SRequestsIDs(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _ArbitrumBeaconVRFConsumer.contract.Call(opts, &out, "s_requestsIDs", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerSession) SRequestsIDs(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _ArbitrumBeaconVRFConsumer.Contract.SRequestsIDs(&_ArbitrumBeaconVRFConsumer.CallOpts, arg0, arg1)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerCallerSession) SRequestsIDs(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _ArbitrumBeaconVRFConsumer.Contract.SRequestsIDs(&_ArbitrumBeaconVRFConsumer.CallOpts, arg0, arg1)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerCaller) SSubId(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _ArbitrumBeaconVRFConsumer.contract.Call(opts, &out, "s_subId")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerSession) SSubId() (uint64, error) {
	return _ArbitrumBeaconVRFConsumer.Contract.SSubId(&_ArbitrumBeaconVRFConsumer.CallOpts)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerCallerSession) SSubId() (uint64, error) {
	return _ArbitrumBeaconVRFConsumer.Contract.SSubId(&_ArbitrumBeaconVRFConsumer.CallOpts)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestID *big.Int, randomWords []*big.Int, arguments []byte) (*types.Transaction, error) {
	return _ArbitrumBeaconVRFConsumer.contract.Transact(opts, "rawFulfillRandomWords", requestID, randomWords, arguments)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerSession) RawFulfillRandomWords(requestID *big.Int, randomWords []*big.Int, arguments []byte) (*types.Transaction, error) {
	return _ArbitrumBeaconVRFConsumer.Contract.RawFulfillRandomWords(&_ArbitrumBeaconVRFConsumer.TransactOpts, requestID, randomWords, arguments)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerTransactorSession) RawFulfillRandomWords(requestID *big.Int, randomWords []*big.Int, arguments []byte) (*types.Transaction, error) {
	return _ArbitrumBeaconVRFConsumer.Contract.RawFulfillRandomWords(&_ArbitrumBeaconVRFConsumer.TransactOpts, requestID, randomWords, arguments)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerTransactor) SetFail(opts *bind.TransactOpts, shouldFail bool) (*types.Transaction, error) {
	return _ArbitrumBeaconVRFConsumer.contract.Transact(opts, "setFail", shouldFail)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerSession) SetFail(shouldFail bool) (*types.Transaction, error) {
	return _ArbitrumBeaconVRFConsumer.Contract.SetFail(&_ArbitrumBeaconVRFConsumer.TransactOpts, shouldFail)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerTransactorSession) SetFail(shouldFail bool) (*types.Transaction, error) {
	return _ArbitrumBeaconVRFConsumer.Contract.SetFail(&_ArbitrumBeaconVRFConsumer.TransactOpts, shouldFail)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerTransactor) StoreBeaconRequest(opts *bind.TransactOpts, reqId *big.Int, height *big.Int, delay *big.Int, numWords uint16) (*types.Transaction, error) {
	return _ArbitrumBeaconVRFConsumer.contract.Transact(opts, "storeBeaconRequest", reqId, height, delay, numWords)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerSession) StoreBeaconRequest(reqId *big.Int, height *big.Int, delay *big.Int, numWords uint16) (*types.Transaction, error) {
	return _ArbitrumBeaconVRFConsumer.Contract.StoreBeaconRequest(&_ArbitrumBeaconVRFConsumer.TransactOpts, reqId, height, delay, numWords)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerTransactorSession) StoreBeaconRequest(reqId *big.Int, height *big.Int, delay *big.Int, numWords uint16) (*types.Transaction, error) {
	return _ArbitrumBeaconVRFConsumer.Contract.StoreBeaconRequest(&_ArbitrumBeaconVRFConsumer.TransactOpts, reqId, height, delay, numWords)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerTransactor) TestRedeemRandomness(opts *bind.TransactOpts, requestID *big.Int) (*types.Transaction, error) {
	return _ArbitrumBeaconVRFConsumer.contract.Transact(opts, "testRedeemRandomness", requestID)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerSession) TestRedeemRandomness(requestID *big.Int) (*types.Transaction, error) {
	return _ArbitrumBeaconVRFConsumer.Contract.TestRedeemRandomness(&_ArbitrumBeaconVRFConsumer.TransactOpts, requestID)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerTransactorSession) TestRedeemRandomness(requestID *big.Int) (*types.Transaction, error) {
	return _ArbitrumBeaconVRFConsumer.Contract.TestRedeemRandomness(&_ArbitrumBeaconVRFConsumer.TransactOpts, requestID)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerTransactor) TestRequestRandomness(opts *bind.TransactOpts, numWords uint16, subID uint64, confirmationDelayArg *big.Int) (*types.Transaction, error) {
	return _ArbitrumBeaconVRFConsumer.contract.Transact(opts, "testRequestRandomness", numWords, subID, confirmationDelayArg)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerSession) TestRequestRandomness(numWords uint16, subID uint64, confirmationDelayArg *big.Int) (*types.Transaction, error) {
	return _ArbitrumBeaconVRFConsumer.Contract.TestRequestRandomness(&_ArbitrumBeaconVRFConsumer.TransactOpts, numWords, subID, confirmationDelayArg)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerTransactorSession) TestRequestRandomness(numWords uint16, subID uint64, confirmationDelayArg *big.Int) (*types.Transaction, error) {
	return _ArbitrumBeaconVRFConsumer.Contract.TestRequestRandomness(&_ArbitrumBeaconVRFConsumer.TransactOpts, numWords, subID, confirmationDelayArg)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerTransactor) TestRequestRandomnessFulfillment(opts *bind.TransactOpts, subID uint64, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte) (*types.Transaction, error) {
	return _ArbitrumBeaconVRFConsumer.contract.Transact(opts, "testRequestRandomnessFulfillment", subID, numWords, confirmationDelayArg, callbackGasLimit, arguments)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerSession) TestRequestRandomnessFulfillment(subID uint64, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte) (*types.Transaction, error) {
	return _ArbitrumBeaconVRFConsumer.Contract.TestRequestRandomnessFulfillment(&_ArbitrumBeaconVRFConsumer.TransactOpts, subID, numWords, confirmationDelayArg, callbackGasLimit, arguments)
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumerTransactorSession) TestRequestRandomnessFulfillment(subID uint64, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte) (*types.Transaction, error) {
	return _ArbitrumBeaconVRFConsumer.Contract.TestRequestRandomnessFulfillment(&_ArbitrumBeaconVRFConsumer.TransactOpts, subID, numWords, confirmationDelayArg, callbackGasLimit, arguments)
}

type SMyBeaconRequests struct {
	SlotNumber        uint32
	ConfirmationDelay *big.Int
	NumWords          uint16
	Requester         common.Address
}

func (_ArbitrumBeaconVRFConsumer *ArbitrumBeaconVRFConsumer) Address() common.Address {
	return _ArbitrumBeaconVRFConsumer.address
}

type ArbitrumBeaconVRFConsumerInterface interface {
	ARBSYS(opts *bind.CallOpts) (common.Address, error)

	NUMCONFDELAYS(opts *bind.CallOpts) (uint8, error)

	Fail(opts *bind.CallOpts) (bool, error)

	IBeaconPeriodBlocks(opts *bind.CallOpts) (*big.Int, error)

	SReceivedRandomnessByRequestID(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error)

	SGasAvailable(opts *bind.CallOpts) (*big.Int, error)

	SMyBeaconRequests(opts *bind.CallOpts, arg0 *big.Int) (SMyBeaconRequests,

		error)

	SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	SRequestsIDs(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error)

	SSubId(opts *bind.CallOpts) (uint64, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestID *big.Int, randomWords []*big.Int, arguments []byte) (*types.Transaction, error)

	SetFail(opts *bind.TransactOpts, shouldFail bool) (*types.Transaction, error)

	StoreBeaconRequest(opts *bind.TransactOpts, reqId *big.Int, height *big.Int, delay *big.Int, numWords uint16) (*types.Transaction, error)

	TestRedeemRandomness(opts *bind.TransactOpts, requestID *big.Int) (*types.Transaction, error)

	TestRequestRandomness(opts *bind.TransactOpts, numWords uint16, subID uint64, confirmationDelayArg *big.Int) (*types.Transaction, error)

	TestRequestRandomnessFulfillment(opts *bind.TransactOpts, subID uint64, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte) (*types.Transaction, error)

	Address() common.Address
}
